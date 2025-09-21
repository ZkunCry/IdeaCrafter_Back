package dto

type CreateVacancyInput struct {
	StartupID   uint   `json:"startup_id" validate:"required"`
	RoleID      uint   `json:"role_id" validate:"required"`
	Description string `json:"description"`
}

type UpdateVacancyInput struct {
	Description string `json:"description,omitempty"`
	IsOpen      *bool  `json:"is_open,omitempty"`
}

type VacancyResponse struct {
	ID          uint   `json:"id"`
	StartupID   uint   `json:"startup_id"`
	RoleID      uint   `json:"role_id"`
	RoleName    string `json:"role_name,omitempty"`
	Description string `json:"description"`
	IsOpen      bool   `json:"is_open"`
}