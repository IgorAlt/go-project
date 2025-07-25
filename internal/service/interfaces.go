package service

import (
	"unrealProject/internal/handlers/dto"
	"unrealProject/internal/models"
)

type UserServiceInterface interface {
	CreateUser(req *dto.CreateUserRequest) (*models.User, error)
	GetUserById(id int) (*models.User, error)
}
