package application

import (
	"startup_back/internal/entity"
	"startup_back/internal/user"
	"startup_back/internal/vacancy"
	"time"
)

const (
	StatusPending  = "pending"
	StatusAccepted = "accepted"
	StatusRejected = "rejected"
)

type CreateInput struct {
	VacancyID uint   `json:"vacancy_id" validate:"required"`
	UserID    uint   `json:"user_id" validate:"required"`
	Message   string `json:"message"`
}

type UpdateStatusInput struct {
	Status string `json:"status" validate:"required,oneof=pending accepted rejected"`
}

type UpdateInput struct {
	Message string `json:"message"`
}


type StartupBrief struct {
	ID               uint   `json:"id"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	LogoURL          string `json:"logo_url"`
	Deleted          bool   `json:"deleted"`
}

type Response struct {
	ID        uint             `json:"id"`
	VacancyID uint             `json:"vacancy_id"`
	Vacancy   vacancy.Response `json:"vacancy"`
	StartupID uint             `json:"startup_id"`
	Startup   *StartupBrief    `json:"startup,omitempty"`
	UserID    uint             `json:"user_id"`
	User      user.Response    `json:"user"`
	Message   string           `json:"message"`
	Status    string           `json:"status"`
	CreatedAt time.Time        `json:"created_at"`
}

func NewResponse(a *entity.Application) Response {
	if a == nil {
		return Response{}
	}

	return Response{
		ID:        a.ID,
		VacancyID: a.VacancyID,
		Vacancy:   vacancy.NewResponse(&a.Vacancy),
		StartupID: a.Vacancy.StartupID,
		UserID:    a.UserID,
		User:      user.NewResponse(&a.User),
		Message:   a.Message,
		Status:    a.Status,
		CreatedAt: a.CreatedAt,
	}
}

func NewResponses(applications []*entity.Application) []Response {
	responses := make([]Response, 0, len(applications))
	for _, application := range applications {
		responses = append(responses, NewResponse(application))
	}
	return responses
}
