package connection

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/beanstalkd/go-beanstalk"
)

type TubeConn struct {
	Conn    *beanstalk.Conn
	Tube    *beanstalk.Tube
	TubeSet *beanstalk.TubeSet
	Pri     uint32
	Delay   int
	TTR     int
	TimeOut int
}

var (
	connections map[string]*TubeConn
	mutex       sync.RWMutex
)

func init() {
	connections = make(map[string]*TubeConn)
}

// NewTubeConnection 创建新的 beanstalkd 连接并设置 tube
func NewTubeConnection(host string, port int, tubeName string, pri uint32, delay int, ttr int, timeout int) (*TubeConn, error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := beanstalk.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	tube := &beanstalk.Tube{Conn: conn, Name: tubeName}
	tubeSet := beanstalk.NewTubeSet(conn, tubeName)

	return &TubeConn{
		Conn:    conn,
		Tube:    tube,
		TubeSet: tubeSet,
		Pri:     pri,
		Delay:   delay,
		TTR:     ttr,
		TimeOut: timeout,
	}, nil
}

// SetConnection 设置连接
func SetConnection(name string, tubeConn *TubeConn) {
	mutex.Lock()
	defer mutex.Unlock()
	connections[name] = tubeConn
}

// GetConnection 获取连接
func GetConnection(name string) (*TubeConn, bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	conn, ok := connections[name]
	return conn, ok
}

// CloseConnection 关闭指定连接
func CloseConnection(name string) error {
	mutex.Lock()
	defer mutex.Unlock()
	if conn, ok := connections[name]; ok {
		err := conn.Conn.Close()
		delete(connections, name)
		return err
	}
	return nil
}

// CloseAll 关闭所有连接
func CloseAll() {
	mutex.Lock()
	defer mutex.Unlock()
	for name, conn := range connections {
		_ = conn.Conn.Close()
		delete(connections, name)
	}
}

// GetAllConnections 获取所有连接 (用于内部查找)
func GetAllConnections() map[string]*TubeConn {
	mutex.RLock()
	defer mutex.RUnlock()
	// 返回副本避免并发问题
	copy := make(map[string]*TubeConn)
	for k, v := range connections {
		copy[k] = v
	}
	return copy
}

// 自定义操作函数

// Put
func (tc *TubeConn) Put(ctx context.Context, body interface{}) (uint64, error) {
	bodyStr, err := anyToString(body)
	if err != nil {
		return 0, err
	}
	return tc.Tube.Put([]byte(bodyStr), tc.Pri, time.Duration(tc.Delay), time.Duration(tc.TTR))
}

// Reserve 从队列中获取一个任务（阻塞等待）
func (tc *TubeConn) Reserve(ctx context.Context) (uint64, string, error) {
	id, body, err := tc.Conn.Reserve(time.Duration(tc.TimeOut))
	if err != nil {
		return 0, "", err
	}
	return id, string(body), nil
}

// ReserveWithTimeout 带超时等待获取任务
func (tc *TubeConn) ReserveWithTimeout(ctx context.Context, timeout time.Duration) (uint64, string, error) {
	id, body, err := tc.Conn.Reserve(timeout)
	if err != nil {
		return 0, "", err
	}
	return id, string(body), nil
}

// Delete 删除指定ID的任务
func (tc *TubeConn) Delete(ctx context.Context, id uint64) error {
	return tc.Conn.Delete(id)
}

// Release 将任务重新放回队列（可设置优先级和延迟）
func (tc *TubeConn) Release(ctx context.Context, id uint64, pri uint32, delay time.Duration) error {
	return tc.Conn.Release(id, pri, delay)
}

// anyToString 将任意类型转换为字符串
// 字符串不变，数字、bool值转成字符串，其他类型转成json字符串
func anyToString(v interface{}) (string, error) {
	switch val := v.(type) {
	case string:
		return val, nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return strconv.FormatInt(reflect.ValueOf(v).Int(), 10), nil
	case float32, float64:
		return strconv.FormatFloat(reflect.ValueOf(v).Float(), 'f', -1, 64), nil
	case bool:
		return strconv.FormatBool(val), nil
	default:
		jsonBytes, err := json.Marshal(val)
		if err != nil {
			return "", err
		}
		return string(jsonBytes), nil
	}
}
