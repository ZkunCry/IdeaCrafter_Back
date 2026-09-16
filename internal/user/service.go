package user

import (
	"context"
	"startup_back/internal/entity"
)

type Service interface {
	Create(ctx context.Context, input CreateUserInput) (*entity.User, error)
	GetByID(ctx context.Context, id uint) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, id uint, input CreateUserInput) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, input CreateUserInput) (*entity.User, error) {
	user, err := s.repo.Create(ctx, &entity.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: input.Password,
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *service) Update(ctx context.Context, id uint, input CreateUserInput) error {
	return s.repo.Update(ctx, id, &entity.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: input.Password,
	})
}
