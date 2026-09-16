package auth

import (
	"context"
	"errors"
	"startup_back/internal/entity"
	"startup_back/internal/user"
)

type Response struct {
	User         *entity.User
	AccessToken  string `json:"access_token"`
	RefreshToken string
}

type Service interface {
	SignIn(ctx context.Context, email, password string) (Response, error)
	SignUp(ctx context.Context, input user.CreateUserInput) (Response, error)
	IdentityMe(ctx context.Context, token string) (Response, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, error)
}

type service struct {
	users     user.Service
	passwords PasswordService
	tokens    TokenService
}

func NewService(userService user.Service, passwordService PasswordService, tokenService TokenService) Service {
	return &service{
		users:     userService,
		passwords: passwordService,
		tokens:    tokenService,
	}
}

func (s *service) SignUp(ctx context.Context, input user.CreateUserInput) (Response, error) {
	existingUser, err := s.users.GetByEmail(ctx, input.Email)
	if err == nil && existingUser.ID != 0 {
		return Response{}, errors.New("User already exists")
	}

	hashPassword, err := s.passwords.HashPassword(input.Password)
	if err != nil {
		return Response{}, err
	}

	createdUser, err := s.users.Create(ctx, user.CreateUserInput{
		Username: input.Username,
		Email:    input.Email,
		Password: hashPassword,
	})
	if err != nil {
		return Response{}, err
	}

	accessToken, err := s.tokens.GenerateAccessToken(createdUser.ID)
	if err != nil {
		return Response{}, err
	}

	refreshToken, err := s.tokens.GenerateRefreshToken(createdUser.ID)
	if err != nil {
		return Response{}, err
	}

	return Response{
		User:         createdUser,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *service) SignIn(ctx context.Context, email, password string) (Response, error) {
	foundUser, err := s.users.GetByEmail(ctx, email)
	if err != nil || foundUser.ID == 0 {
		return Response{}, errors.New("invalid email or password")
	}

	if err := s.passwords.ComparePassword(foundUser.PasswordHash, password); err != nil {
		return Response{}, errors.New("invalid email or password")
	}

	accessToken, err := s.tokens.GenerateAccessToken(foundUser.ID)
	if err != nil {
		return Response{}, err
	}

	refreshToken, err := s.tokens.GenerateRefreshToken(foundUser.ID)
	if err != nil {
		return Response{}, err
	}

	return Response{
		User:         foundUser,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *service) IdentityMe(ctx context.Context, token string) (Response, error) {
	if token == "" {
		return Response{}, errors.New("access token not provided")
	}

	userID, err := s.tokens.ValidateAccessToken(token)
	if err != nil {
		return Response{}, errors.New("invalid or expired token")
	}

	foundUser, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return Response{}, err
	}

	return Response{User: foundUser}, nil
}

func (s *service) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	userID, err := s.tokens.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", errors.New("invalid or expired refresh token")
	}

	return s.tokens.GenerateAccessToken(userID)
}
