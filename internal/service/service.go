package service

import "EffectiveMobile/internal/repository"

type User interface {
	Test()
}

type Service struct {
	User
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		User: NewUserService(repo),
	}
}
