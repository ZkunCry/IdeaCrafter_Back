package handler

import (
	"startup_back/internal/service"
)

type Handlers struct {
	User *UserHandler
}

func NewHandlers(services * service.Services) *Handlers {
	
	return &Handlers{
		User: NewUserHandler(services.User),
	}
}