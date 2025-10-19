package handler

import (
	"startup_back/internal/service"
)

type Handlers struct {
	Auth *AuthHandler
	Startup *StartupHandler
	Vacancy * VacancyHandler
	Role *RoleHandler
	Application *ApplicationHandler
	Stage *StageHandler	
}

func NewHandlers(services * service.Services) *Handlers {
	
	return &Handlers{
		Auth: NewAuthHandler(services),
		Startup: NewStartupHandler(services),
		Vacancy:  NewVacancyHandler(services),
		Role: NewRoleHandler(services),
		Application: NewApplicationHandler(services),
		Stage: NewStageHandler(services),
	}
}