package stage

import (
	"context"
	"startup_back/internal/entity"
)

type Service interface {
	GetAll(ctx context.Context) ([]*entity.Stage, error)
	GetByID(ctx context.Context, id uint) (*entity.Stage, error)
	Create(ctx context.Context, input *CreateInput) (*entity.Stage, error)
	Update(ctx context.Context, id uint, stage *entity.Stage) error
	Delete(ctx context.Context, id uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAll(ctx context.Context) ([]*entity.Stage, error) {
	return s.repo.GetAll(ctx)
}

func (s *service) GetByID(ctx context.Context, id uint) (*entity.Stage, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) Create(ctx context.Context, input *CreateInput) (*entity.Stage, error) {
	return s.repo.Create(ctx, &entity.Stage{Name: input.Name})
}

func (s *service) Update(ctx context.Context, id uint, stage *entity.Stage) error {
	return s.repo.Update(ctx, id, stage)
}

func (s *service) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
