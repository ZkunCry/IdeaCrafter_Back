package user

import "startup_back/internal/entity"

type CreateUserInput struct {
	Username string `json:"username" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type Response struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UpdateUserInput struct {
	Username string `json:"username,omitempty" validate:"omitempty,min=3"`
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
}

func NewResponse(foundUser *entity.User) Response {
	if foundUser == nil {
		return Response{}
	}

	return Response{
		ID:       foundUser.ID,
		Username: foundUser.Username,
		Email:    foundUser.Email,
	}
}
