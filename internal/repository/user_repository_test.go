package repository

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
	"unrealProject/internal/models"
)

func setupTestDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return sqlxDB, mock
}

func setupTestRepository(t *testing.T) (*UserRepository, sqlmock.Sqlmock, func()) {
	db, mock := setupTestDB(t)
	repository := NewUserRepository(db)
	cleanup := func() { _ = db.Close() }
	return repository, mock, cleanup
}

func TestUserRepository_Create(t *testing.T) {
	repository, mock, cleanup := setupTestRepository(t)
	defer cleanup()

	user := &models.User{Name: "John", Email: "john@gmail.com", Password: "hashed-password"}

	mock.ExpectQuery("INSERT INTO users").WithArgs(user.Name, user.Email, user.Password).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	createdUser, err := repository.Create(user)

	assert.NoError(t, err)
	assert.Equal(t, 1, createdUser.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetById_Success(t *testing.T) {
	repository, mock, cleanup := setupTestRepository(t)
	defer cleanup()

	expectedUser := &models.User{ID: 1, Name: "John", Email: "john@gmail.com"}

	rows := sqlmock.NewRows([]string{"id", "name", "email"}).AddRow(expectedUser.ID, expectedUser.Name, expectedUser.Email)

	mock.ExpectQuery("SELECT id, name, email FROM users WHERE id = \\$1").WithArgs(expectedUser.ID).WillReturnRows(rows)

	user, err := repository.GetById(expectedUser.ID)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetById_DBError(t *testing.T) {
	repository, mock, cleanup := setupTestRepository(t)
	defer cleanup()

	mock.ExpectQuery("SELECT id, name, email FROM users WHERE id = \\$1").WithArgs(1).WillReturnError(fmt.Errorf("some error"))

	user, err := repository.GetById(1)

	assert.Nil(t, user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "some error")
	assert.NoError(t, mock.ExpectationsWereMet())
}
