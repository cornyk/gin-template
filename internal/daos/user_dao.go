package daos

import (
	"cornyk/gin-template/internal/models"
	"cornyk/gin-template/pkg/global"
	"github.com/gin-gonic/gin"
	"time"
)

// GetAllUsers 获取所有用户
func GetAllUsers(c *gin.Context) ([]models.User, error) {
	redis := global.RedisConn()
	redis.Set(c, "TEST_KEY", "testValue", time.Second*1000)

	redis2 := global.RedisConn("cache")
	redis3 := global.RedisConn("session")
	redis2.Set(c, "TEST_KEY1", "testValue", time.Second*1000)
	redis3.Set(c, "TEST_KEY2", "testValue", time.Second*1000)

	db := global.DBConn()
	var users []models.User
	if err := db.WithContext(c).Select("id", "name").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
