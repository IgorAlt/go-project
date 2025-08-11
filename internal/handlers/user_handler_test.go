package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"unrealProject/internal/handlers/dto"
	"unrealProject/internal/models"
)

type MockUserService struct{ mock.Mock }

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
	mockService.
		On("CreateUser", mock.Anything).
		Return(nil, errors.New("service error"))

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

	user := &models.User{ID: 1, Name: "John", Email: "john@gmail.com"}
	matchReq := mock.MatchedBy(func(r *dto.CreateUserRequest) bool {
		return r.Name == "John" && r.Email == "john@gmail.com" && r.Password == "123456"
	})
	mockService.On("CreateUser", matchReq).Return(user, nil)

	handler := NewUserHandler(mockService)

	reqBody := `{"name":"John","email":"john@gmail.com","password":"123456"}`
	req := httptest.NewRequest(http.MethodPost, "/users/create", strings.NewReader(reqBody))
	rec := httptest.NewRecorder()

	handler.CreateUser(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.User
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err, "response must be valid JSON")
	assert.Equal(t, *user, response)

	mockService.AssertExpectations(t)
}

func withChiParam(req *http.Request, key, val string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, val)
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

func TestGetUserById_InvalidID(t *testing.T) {
	handler := NewUserHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/users/abc", nil)
	req = withChiParam(req, "id", "abc")
	rec := httptest.NewRecorder()

	handler.GetUserById(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid user id")
}

func TestGetUserById_ServiceError(t *testing.T) {
	mockService := new(MockUserService)
	mockService.
		On("GetUserById", 42).
		Return(nil, errors.New("db down"))

	handler := NewUserHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/users/42", nil)
	req = withChiParam(req, "id", "42")
	rec := httptest.NewRecorder()

	handler.GetUserById(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "db down")

	mockService.AssertExpectations(t)
}

func TestGetUserById_Success(t *testing.T) {
	mockService := new(MockUserService)
	user := &models.User{ID: 7, Name: "Alice", Email: "alice@example.com"}
	mockService.On("GetUserById", 7).Return(user, nil)

	handler := NewUserHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/users/7", nil)
	req = withChiParam(req, "id", "7")
	rec := httptest.NewRecorder()

	handler.GetUserById(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var got models.User
	err := json.NewDecoder(rec.Body).Decode(&got)
	assert.NoError(t, err, "response must be valid JSON")
	assert.Equal(t, *user, got)

	mockService.AssertExpectations(t)
}
