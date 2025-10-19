package dto

type CreateStageInput struct {
	Name string `json:"name" validate:"required"`
}

type StageResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}