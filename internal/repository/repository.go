package repository

import (
	"EffectiveMobile/entity"
	"github.com/jmoiron/sqlx"
)

const usersTable = "users"

type User interface {
	Create(user entity.User) (int, error)
	Delete(userId int) error
	Update(userId int, user entity.User) error
	GetById(userId int) (entity.User, error)
	GetAll(query string) ([]entity.User, error)
}

type Repository struct {
	User
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User: NewUserPostgres(db),
	}
}
