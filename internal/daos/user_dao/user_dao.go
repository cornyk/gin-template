package user_dao

import (
	"cornyk/gin-template/internal/models"
	"cornyk/gin-template/pkg/global"
	"time"

	"github.com/gin-gonic/gin"
)

// GetAllUsers 获取所有用户
func GetAllUsers(c *gin.Context) ([]models.UserModel, error) {
	redis := global.RedisConn()
	redis.Set(c, "TEST_KEY", "testValue", time.Second*1000)

	redis2 := global.RedisConn("cache")
	redis3 := global.RedisConn("session")
	redis2.Set(c, "TEST_KEY1", "testValue", time.Second*1000)
	redis3.Set(c, "TEST_KEY2", "testValue", time.Second*1000)

	beanstalkd := global.BeanstalkdConn()
	beanstalkd.Conn.Put([]byte("test"), 1024, time.Duration(0), time.Duration(60))

	db := global.DBConn()
	var users []models.UserModel
	if err := db.WithContext(c).Select("id", "name").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
