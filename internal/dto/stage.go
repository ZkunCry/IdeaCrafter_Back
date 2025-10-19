package dto

type CreateStageInput struct {
	Name string `json:"name" validate:"required"`
}