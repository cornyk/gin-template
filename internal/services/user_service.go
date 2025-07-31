package services

import (
	"cornyk/gin-template/internal/daos"
	"cornyk/gin-template/internal/models"

	"github.com/gin-gonic/gin"
)

type UserService struct {
	userDao *daos.UserDao
}

func (s *UserService) GetAllUsers(ctx *gin.Context) ([]models.UserModel, error) {
	return s.userDao.GetAllUsers(ctx)
}
