package dto

type CreateRoleInput struct {
	Name string `json:"name" validate:"required"`
}

type RoleResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}