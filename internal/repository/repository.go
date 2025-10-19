package repository

import (
	"gorm.io/gorm"
)

type Repositories struct {
    User          UserRepository
    Startup       StartupRepository
		Vacancy       VacancyRepository
		Role          RoleRepository
		Application   ApplicationRepository
		Stage 				StageRepository
}

func NewRespositories(db *gorm.DB) * Repositories{
	return &Repositories{
		User : NewUserRepository(db),
		Startup: NewStartupRepository(db),
		Vacancy: NewVacancyRepository(db),
		Role: NewRoleRepository(db),
		Application: NewApplicationRepository(db),
		Stage: NewStageRepository(db),
	}
}