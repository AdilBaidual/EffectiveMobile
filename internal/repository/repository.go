package repository

import (
	"EffectiveMobile/entity"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

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

func NewRepository(db *sqlx.DB, client *redis.Client) *Repository {
	cache := NewRedisRepository(client)
	return &Repository{
		User: NewUserPostgres(db, cache),
	}
}
