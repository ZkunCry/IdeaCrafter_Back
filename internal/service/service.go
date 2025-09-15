package service

import (
	"context"
	"startup_back/internal/config"
	"startup_back/internal/domain"
	"startup_back/internal/repository"
)
type UserService interface{
	CreateUser(ctx context.Context,input CreateUserInput)(*domain.User, error)
	GetUserById(ctx context.Context,id uint)(*domain.User, error)
	GetUserByEmail(ctx context.Context,email string)(*domain.User, error)
	UpdateUser(ctx context.Context,id uint,input CreateUserInput)error
}
type PasswordService interface{
	HashPassword(password string) (string, error)
  ComparePassword(hashedPassword, password string) error
}
type AuthResponse struct {
	User *domain.User
	AccessToken string `json:"access_token"`
	RefreshToken string 
}

type AuthService interface{
	SignInUser(ctx context.Context, email, password string) (response AuthResponse, err error)
	SignUpUser(ctx context.Context, input CreateUserInput)  (response AuthResponse, err error)
}
type TokenService interface{
	GenerateAccessToken(userId uint) (string, error)
	GenerateRefreshToken(userId uint) (string, error)
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
		Password PasswordService
		Token TokenService
		Auth AuthService
}

func NewServices(repos *repository.Repositories,cfg *config.Config) *Services{

	return &Services{
		User: NewUserService(repos.User),
		Password: NewPasswordService(),
		Token: NewTokenService(cfg.JWT.AccessSecret,cfg.JWT.RefreshSecret),
		Auth: NewAuthService(NewUserService(repos.User),NewPasswordService(),NewTokenService(cfg.JWT.AccessSecret,cfg.JWT.RefreshSecret)),
	}
}