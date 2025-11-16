package service

import (
	"startup_back/internal/config"
	"startup_back/internal/domain"
	"startup_back/internal/repository"

	"github.com/aws/aws-sdk-go-v2/service/s3"
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
		Role RoleService
		Application ApplicationService
		Stage StageService
		Category CategoryService
		S3  *s3.Client

}

func NewServices(repos *repository.Repositories,cfg *config.AppConfig) *Services{
	return &Services{
		User: NewUserService(repos.User),
		Password: NewPasswordService(),
		Token: NewTokenService(cfg.JWT.AccessSecret,cfg.JWT.RefreshSecret),
		Auth: NewAuthService(NewUserService(repos.User),NewPasswordService(),NewTokenService(cfg.JWT.AccessSecret,cfg.JWT.RefreshSecret)),
		Startup: NewStartupService(repos.Startup),
		Vacancy: NewVacancyService(repos.Vacancy),
		Role: NewRoleService(repos.Role),	
		Application: NewApplicationService(repos.Application),
		Stage: NewStageService(repos.Stage),
		Category: NewCategoryService(repos.Category),
		S3: cfg.S3Client,
	}
}