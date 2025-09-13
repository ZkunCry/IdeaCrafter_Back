package service

import (
	"context"
	"startup_back/internal/domain"
	"startup_back/internal/repository"
)
type userService struct{
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService{
	return &userService{repo:repo}
}
func (s *userService) CreateUser(ctx context.Context, input CreateUserInput) (*domain.User, error) {
	
}

func (s *userService)	GetUserById(ctx context.Context,id uint)(*domain.User, error)
func (s *userService)	UpdateUser(ctx context.Context,id uint,input CreateUserInput)(*domain.User, error)