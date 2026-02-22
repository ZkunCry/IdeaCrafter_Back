package dto

type CreateApplicationInput struct {
	VacancyID uint   `json:"vacancy_id" validate:"required"`
	UserID    uint   `json:"user_id" validate:"required"`
	Message   string `json:"message"`
}

type UpdateApplicationStatusInput struct {
	Status string `json:"status" validate:"required,oneof=pending accepted rejected"`
}

type UpdateApplicationInput struct {
	Message string `json:"message"`
}
