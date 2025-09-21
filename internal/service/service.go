package service

import (
	"startup_back/internal/config"
	"startup_back/internal/domain"
	"startup_back/internal/repository"
)
type AuthResponse struct {
	User *domain.User
	AccessToken string `json:"access_token"`
	RefreshToken string 
}

type Services struct {
    User   UserService
		Password PasswordService
		Token TokenService
		Auth AuthService
		Startup StartupService
		Vacancy VacancyService

}

func NewServices(repos *repository.Repositories,cfg *config.Config) *Services{
	return &Services{
		User: NewUserService(repos.User),
		Password: NewPasswordService(),
		Token: NewTokenService(cfg.JWT.AccessSecret,cfg.JWT.RefreshSecret),
		Auth: NewAuthService(NewUserService(repos.User),NewPasswordService(),NewTokenService(cfg.JWT.AccessSecret,cfg.JWT.RefreshSecret)),
		Startup: NewStartupService(repos.Startup),
		Vacancy: NewVacancyService(repos.Vacancy),
	}
}