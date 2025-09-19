package dto

type CreateStartupInput struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	CategoryIDs []uint `json:"category_ids" validate:"required"`
	VacancyIDs  []uint `json:"vacancy_ids,omitempty"`
}