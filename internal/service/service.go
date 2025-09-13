package service

import (
	"context"
	"startup_back/internal/domain"
	"startup_back/internal/repository"
)
type UserService interface{
	CreateUser(ctx context.Context,input CreateUserInput)(*domain.User, error)
	GetUserById(ctx context.Context,id uint)(*domain.User, error)
	UpdateUser(ctx context.Context,id uint,input CreateUserInput)(*domain.User, error)
}

type CreateUserInput struct {
    Username string `json:"username" validate:"required,min=3"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

type UpdateUserInput struct {
    Username string `json:"username,omitempty" validate:"omitempty,min=3"`
    Email    string `json:"email,omitempty" validate:"omitempty,email"`
}

type Services struct {
    User   UserService
}

func NewServices(repos *repository.Repositories) *Services{
	return &Services{
		User: NewUserService(repos.User),
	}
}