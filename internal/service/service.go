package service

import (
	"context"
	"startup_back/internal/config"
	"startup_back/internal/domain"
	"startup_back/internal/dto"
	"startup_back/internal/repository"
)


type UserService interface{
	CreateUser(ctx context.Context,input dto.CreateUserInput)(*domain.User, error)
	GetUserById(ctx context.Context,id uint)(*domain.User, error)
	GetUserByEmail(ctx context.Context,email string)(*domain.User, error)
	UpdateUser(ctx context.Context,id uint,input dto.CreateUserInput)error
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
	SignUpUser(ctx context.Context, input dto.CreateUserInput)  (response AuthResponse, err error)
}
type TokenService interface{
	GenerateAccessToken(userId uint) (string, error)
	GenerateRefreshToken(userId uint) (string, error)
	ValidateAccessToken(tokenString string) (uint, error)
}

type StartupService interface {
	Create(ctx context.Context, startup *dto.CreateStartupInput, categoryIDs []uint, vacancyRoleIDs []uint) (*domain.Startup, error)
  GetByID(ctx context.Context, id uint) (*domain.Startup, error)
  List(ctx context.Context, limit, offset int, categoryID uint) ([]*domain.Startup, error)
  Delete(ctx context.Context, id uint) error	
}

type Services struct {
    User   UserService
		Password PasswordService
		Token TokenService
		Auth AuthService
		Startup StartupService

}

func NewServices(repos *repository.Repositories,cfg *config.Config) *Services{
	return &Services{
		User: NewUserService(repos.User),
		Password: NewPasswordService(),
		Token: NewTokenService(cfg.JWT.AccessSecret,cfg.JWT.RefreshSecret),
		Auth: NewAuthService(NewUserService(repos.User),NewPasswordService(),NewTokenService(cfg.JWT.AccessSecret,cfg.JWT.RefreshSecret)),
		Startup: NewStartupService(repos.Startup),
	}
}