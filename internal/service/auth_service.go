package service

import (
	"context"
	"errors"
	"startup_back/internal/domain"
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

func (s *authService) SignUpUser(ctx context.Context, input CreateUserInput) (response AuthResponse, err error) {
	user,err:= s.UserService.GetUserByEmail(ctx, input.Email)
	if err == nil {
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
	createUser, err := s.UserService.CreateUser(ctx, CreateUserInput{
		Username: newUser.Username,
		Email: newUser.Email,
		Password: hashPassword,
	})
	if err != nil {
		return AuthResponse{}, err
	}

	accessToken,err := s.TokenService.GenerateAccessToken(user.ID)
	if err != nil {
		return AuthResponse{}, err
	}
	refreshToken,err := s.TokenService.GenerateRefreshToken(user.ID)
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


	if err != nil {
		return AuthResponse{}, errors.New("user not found")
	}

	if err := s.PasswordService.ComparePassword("fsdfds", password); err != nil {
		return  AuthResponse{}, errors.New("invalid credentials")
	}
	return AuthResponse{}, nil

}