package service

import (
	"golang.org/x/crypto/bcrypt"
	"unrealProject/internal/handlers/dto"
	"unrealProject/internal/models"
	"unrealProject/internal/repository"
)

type UserService struct {
	repo repository.UserRepositoryInterface
}

func NewUserService(repo repository.UserRepositoryInterface) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(req *dto.CreateUserRequest) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	return s.repo.Create(user)
}

func (s *UserService) GetUserById(id int) (*models.User, error) {
	return s.repo.GetById(id)
}
