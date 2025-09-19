package dto

type CreateUserInput struct {
	Username string `json:"username" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UpdateUserInput struct {
	Username string `json:"username,omitempty" validate:"omitempty,min=3"`
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
}