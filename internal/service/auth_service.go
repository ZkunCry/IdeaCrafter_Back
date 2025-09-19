package service

import (
	"context"
	"errors"
	"startup_back/internal/domain"
	"startup_back/internal/dto"
)
type authService struct {
	UserService     UserService
	PasswordService PasswordService
	TokenService    TokenService
}

func NewAuthService(userSvc UserService, passSvc PasswordService,tokenSvc TokenService) AuthService {
	return &authService{
		UserService:     userSvc,
		PasswordService: passSvc,
		TokenService:  tokenSvc,
	}
}

func (s *authService) SignUpUser(ctx context.Context, input dto.CreateUserInput) (response AuthResponse, err error) {
	user,err:= s.UserService.GetUserByEmail(ctx, input.Email)
	if err == nil && user.ID != 0 {
		return AuthResponse{}, errors.New("User already exists")
	}
	hashPassword, err := s.PasswordService.HashPassword(input.Password)
	if err != nil {
		return AuthResponse{}, err
	}
	newUser:= &domain.User{
		Username: input.Username,
		Email: input.Email,
		PasswordHash: hashPassword,
	}
	createUser, err := s.UserService.CreateUser(ctx, dto.CreateUserInput{
		Username: newUser.Username,
		Email: newUser.Email,
		Password: hashPassword,
	})
	if err != nil {
		return AuthResponse{}, err
	}

	accessToken,err := s.TokenService.GenerateAccessToken(createUser.ID)
	if err != nil {
		return AuthResponse{}, err
	}
	refreshToken,err := s.TokenService.GenerateRefreshToken(createUser.ID)
	if err != nil {
		return AuthResponse{}, err
	}
	return AuthResponse{
		User: createUser,
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	},nil
	
}

func (s *authService) SignInUser(ctx context.Context, email, password string) (response AuthResponse, err error) {
	user, err := s.UserService.GetUserByEmail(ctx, email)
	if err != nil || user.ID == 0 {
		return AuthResponse{}, errors.New("invalid email or password")
	}

	if err := s.PasswordService.ComparePassword( user.PasswordHash,password); err != nil {
		return AuthResponse{}, errors.New("invalid email or password")
	}

	accessToken, err := s.TokenService.GenerateAccessToken(user.ID)
	if err != nil {
		return AuthResponse{}, err
	}

	refreshToken, err := s.TokenService.GenerateRefreshToken(user.ID)
	if err != nil {
		return AuthResponse{}, err
	}
	return AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}