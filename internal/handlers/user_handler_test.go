package handlers

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"unrealProject/internal/handlers/dto"
	"unrealProject/internal/models"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(r *dto.CreateUserRequest) (*models.User, error) {
	args := m.Called(r)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}
func (m *MockUserService) GetUserById(id int) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func TestCreateUser_BadInput(t *testing.T) {
	handler := NewUserHandler(nil)

	req := httptest.NewRequest(http.MethodPost, "/users/create", strings.NewReader("invalid json"))
	rec := httptest.NewRecorder()

	handler.CreateUser(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid input")
}

func TestCreateUser_ServiceError(t *testing.T) {
	mockService := new(MockUserService)

	mockService.On("CreateUser", mock.Anything).Return(nil, errors.New("service error"))

	handler := NewUserHandler(mockService)

	reqBody := `{"name": "John", "email": "john@gmail.com"}`
	req := httptest.NewRequest(http.MethodPost, "/users/create", strings.NewReader(reqBody))
	rec := httptest.NewRecorder()

	handler.CreateUser(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "service error")
}

func TestCreateUser_Success(t *testing.T) {
	mockService := new(MockUserService)

	user := &models.User{
		ID:    1,
		Name:  "John",
		Email: "john@gmail.com",
	}

	mockService.On("CreateUser", mock.Anything).Return(user, nil)

	handler := NewUserHandler(mockService)

	reqBody := `{"name": "John", "email": "john@gmail.com", "password": "123456"}`
	req := httptest.NewRequest(http.MethodPost, "/users/create", strings.NewReader(reqBody))
	rec := httptest.NewRecorder()

	handler.CreateUser(rec, req)

	var response models.User
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, *user, response)
}
