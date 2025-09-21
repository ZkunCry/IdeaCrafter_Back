package dto

import (
	"startup_back/internal/domain"
)
type CreateStartupInput struct {
	CreatorId   uint   `json:"creator_id"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	CategoryIDs []uint `json:"category_ids" validate:"required"`
}

type GetStartupList struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type StartupResponse struct{
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatorID   uint   `json:"creator_id"`
	Creator     domain.User   `json:"creator"`
	Categories  []domain.Category `json:"categories"`
	Files       []domain.StartupFile `json:"files"`
	Vacansies   []domain.Vacancy `json:"vacancies"`
	Memberships []domain.Membership `json:"memberships"`
}