package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type tokenSerivce struct{
	accessSecret []byte
	refreshSecret []byte
}

func NewTokenService(accessSecret,refreshSecret string) TokenService{
	return &tokenSerivce{
		accessSecret: []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
	}
}

func (s *tokenSerivce) GenerateAccessToken(userId uint) (string, error) {

	claims := &jwt.MapClaims{
		"user_id": userId,
		"exp":time.Now().Add(time.Minute * 15).Unix(),
	}
	token:= jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	return token.SignedString(s.accessSecret)
}

func (s *tokenSerivce) GenerateRefreshToken(userId uint) (string, error) {
	claims := &jwt.MapClaims{
		"user_id": userId,
		"exp":time.Now().Add(time.Hour * 24 * 7).Unix(),
	}
	token:= jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	return token.SignedString(s.refreshSecret)
}