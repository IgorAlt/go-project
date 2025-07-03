package repository

import "unrealProject/internal/models"

type UserRepositoryInterface interface {
	Create(user *models.User) (*models.User, error)
	GetById(id int) (*models.User, error)
}
