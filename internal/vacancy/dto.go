package vacancy

import (
	"startup_back/internal/entity"
	"startup_back/internal/user"
)

type CreateInput struct {
	StartupID   uint   `json:"startup_id" validate:"required"`
	RoleID      uint   `json:"role_id" validate:"required"`
	Description string `json:"description"`
}

type UpdateInput struct {
	Description *string `json:"description,omitempty"`
	IsOpen      *bool   `json:"is_open,omitempty"`
}

type RoleResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Response struct {
	ID          uint           `json:"id"`
	StartupID   uint           `json:"startup_id"`
	RoleID      uint           `json:"role_id"`
	Role        RoleResponse   `json:"role"`
	RoleName    string         `json:"role_name,omitempty"`
	Description string         `json:"description"`
	IsOpen      bool           `json:"is_open"`
	UserID      *uint          `json:"user_id"`
	User        *user.Response `json:"user"`
}

func NewResponse(v *entity.Vacancy) Response {
	if v == nil {
		return Response{}
	}

	response := Response{
		ID:          v.ID,
		StartupID:   v.StartupID,
		RoleID:      v.RoleID,
		Role:        RoleResponse{ID: v.Role.ID, Name: v.Role.Name},
		RoleName:    v.Role.Name,
		Description: v.Description,
		IsOpen:      v.IsOpen,
		UserID:      v.UserID,
	}

	if v.User != nil {
		assignee := user.NewResponse(v.User)
		response.User = &assignee
	}

	return response
}

func NewResponses(vacancies []*entity.Vacancy) []Response {
	responses := make([]Response, 0, len(vacancies))
	for _, vacancy := range vacancies {
		responses = append(responses, NewResponse(vacancy))
	}
	return responses
}
