package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService interface {
	GenerateAccessToken(userID uint) (string, error)
	GenerateRefreshToken(userID uint) (string, error)
	ValidateAccessToken(tokenString string) (uint, error)
	ValidateRefreshToken(tokenString string) (uint, error)
}

type tokenService struct {
	accessSecret  []byte
	refreshSecret []byte
}

func NewTokenService(accessSecret, refreshSecret string) TokenService {
	return &tokenService{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
	}
}

func (s *tokenService) GenerateAccessToken(userID uint) (string, error) {
	claims := &jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(60 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.accessSecret)
}

func (s *tokenService) GenerateRefreshToken(userID uint) (string, error) {
	claims := &jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.refreshSecret)
}

func (s *tokenService) ValidateAccessToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return s.accessSecret, nil
	})
	if err != nil {
		return 0, err
	}
	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims type")
	}

	uidFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("user_id not found in token")
	}

	return uint(uidFloat), nil
}

func (s *tokenService) ValidateRefreshToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return s.refreshSecret, nil
	})
	if err != nil {
		return 0, err
	}
	if !token.Valid {
		return 0, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims type")
	}

	uidFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("user_id not found in token")
	}

	return uint(uidFloat), nil
}
