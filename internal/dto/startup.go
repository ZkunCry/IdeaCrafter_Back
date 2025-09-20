package dto

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