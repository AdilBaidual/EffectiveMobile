package service

import (
	"EffectiveMobile/internal/repository"
	"fmt"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Test() {
	fmt.Println("Test done")
}
