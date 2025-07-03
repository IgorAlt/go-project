package service

import (
	"github.com/stretchr/testify/mock"
	"testing"
	"unrealProject/internal/handlers/dto"
	"unrealProject/internal/models"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(user *models.User) (*models.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepo) GetById(id int) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func TestUserService_CreateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	service := NewUserService(mockRepo)

	req := &dto.CreateUserRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	expectedUser := &models.User{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	mockRepo.On("Create", mock.MatchedBy(func(user *models.User) bool {
		return user.Name == expectedUser.Name && user.Email == expectedUser.Email && user.Password != ""
	})).Return(expectedUser, nil)

	createdUser, err := service.CreateUser(req)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if createdUser.Name != expectedUser.Name || createdUser.Email != expectedUser.Email {
		t.Errorf("expected user %v, got %v", expectedUser, createdUser)
	}

	mockRepo.AssertExpectations(t)
}
