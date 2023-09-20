package service

import (
	"EffectiveMobile/entity"
	"EffectiveMobile/internal/repository"
)

type User interface {
	CreateUser(user entity.User) (int, error)
	EnrichUser(fio entity.Fio) (int, error)
	DeleteUser(userId int) error
	UpdateUser(userId int, updatedUser entity.User) error
	GetUserById(userId int) (entity.User, error)
	GetUsers(options entity.Options) ([]entity.User, error)
}

type Service struct {
	User
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		User: NewUserService(repo),
	}
}
