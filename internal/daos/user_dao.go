package daos

import (
	"cornyk/gin-template/internal/models"
	"cornyk/gin-template/pkg/global"
	"github.com/gin-gonic/gin"
)

// GetAllUsers 获取所有用户
func GetAllUsers(c *gin.Context) ([]models.User, error) {
	db := global.MainDB
	var users []models.User
	if err := db.WithContext(c).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
