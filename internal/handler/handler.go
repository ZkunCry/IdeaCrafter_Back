package handler

import (
	"startup_back/internal/service"
)

type Handlers struct {
	Auth *AuthHandler
	Startup *StartupHandler
	Vacancy * VacancyHandler
}

func NewHandlers(services * service.Services) *Handlers {
	
	return &Handlers{
		Auth: NewAuthHandler(services),
		Startup: NewStartupHandler(services),
		Vacancy:  NewVacancyHandler(services),
	}
}