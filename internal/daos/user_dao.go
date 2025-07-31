package daos

import (
	"cornyk/gin-template/internal/models"
	"cornyk/gin-template/pkg/global"
	"time"

	"github.com/gin-gonic/gin"
)

type UserDao struct{}

// GetAllUsers 获取所有用户
func (d *UserDao) GetAllUsers(ctx *gin.Context) ([]models.UserModel, error) {
	redis := global.RedisConn()
	redis.Set(ctx, "TEST_KEY", "testValue", time.Second*1000)

	redis2 := global.RedisConn("cache")
	redis3 := global.RedisConn("session")
	redis2.Set(ctx, "TEST_KEY1", "testValue", time.Second*1000)
	redis3.Set(ctx, "TEST_KEY2", "testValue", time.Second*1000)

	beanstalkd := global.BeanstalkdConn("reporting")
	m := make(map[string]string)
	m["name"] = "Alice"
	m["age"] = "30"
	beanstalkd.Put(ctx, m)

	db := global.DBConn()
	var users []models.UserModel
	if err := db.WithContext(ctx).Select("id", "name").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
