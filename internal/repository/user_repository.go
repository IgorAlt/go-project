package repository

import (
	"github.com/jmoiron/sqlx"
	"unrealProject/internal/models"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) (*models.User, error) {
	err := r.db.QueryRow(
		"INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id",
		user.Name, user.Email, user.Password).Scan(&user.ID)
	return user, err
}

func (r *UserRepository) GetById(id int) (*models.User, error) {
	user := &models.User{}
	err := r.db.Get(user, "SELECT id, name, email FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
