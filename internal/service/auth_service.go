package service

import (
	"context"
	"errors"
	"startup_back/internal/config"
)
type authService struct {
	UserService     UserService
	PasswordService PasswordService
	AccessSecret    []byte
	RefreshSecret   []byte
}

func NewAuthService(cfg config.Config, userSvc UserService, passSvc PasswordService) AuthService {
	return &authService{
		UserService:     userSvc,
		PasswordService: passSvc,
		AccessSecret:    []byte(cfg.JWT.AccessSecret),
		RefreshSecret:   []byte(cfg.JWT.RefreshSecret),
	}
}

func (s *authService) SignUpUser(ctx context.Context, input CreateUserInput) (string, error) {
	user, err := s.UserService.CreateUser(ctx, input)
	if err != nil {
		return "", err
	}
	//TODO: добавить генерацию аксес токена
	
}

func (s *authService) SignInUser(ctx context.Context, email, password string) (string, error) {
		//TODO: добавить получение юзера по емайл

	if err != nil {
		return "", errors.New("user not found")
	}

	if err := s.PasswordService.ComparePassword(user.PasswordHash, password); err != nil {
		return "", errors.New("invalid credentials")
	}
	//TODO: добавить генерацию аксес токена

}