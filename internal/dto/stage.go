package dto

type CreateStageInput struct {
	Name string `json:"name" validate:"required"`
}
type Stage struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
type StageResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}