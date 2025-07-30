package connection

import (
	"fmt"
	"sync"

	"github.com/beanstalkd/go-beanstalk"
)

type TubeConn struct {
	Conn    *beanstalk.Conn
	Tube    *beanstalk.Tube
	TubeSet *beanstalk.TubeSet
}

var (
	connections map[string]*TubeConn
	mutex       sync.RWMutex
)

func init() {
	connections = make(map[string]*TubeConn)
}

// NewTubeConnection 创建新的 beanstalkd 连接并设置 tube
func NewTubeConnection(host string, port int, tubeName string) (*TubeConn, error) {
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
