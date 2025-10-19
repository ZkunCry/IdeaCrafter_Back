package service

import (
	"context"
	"startup_back/internal/domain"
	"startup_back/internal/dto"
	"startup_back/internal/repository"
)
type StageService interface {
	GetAll(ctx context.Context) ([]*domain.Stage, error)
	GetByID(ctx context.Context, id uint) (*domain.Stage, error)
	Create(ctx context.Context, input *dto.CreateStageInput) (*domain.Stage, error)
	Update(ctx context.Context, id uint, role *domain.Stage) error
	Delete(ctx context.Context, id uint) error
}

type stageService struct {
	repo repository.StageRepository
}


func NewStageService(repo repository.StageRepository) StageService {
	return &stageService{repo: repo}
}


func (s *stageService) GetAll(ctx context.Context) ([]*domain.Stage, error) {
	stages, err:= s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return stages,nil
}

func (s *stageService) GetByID(ctx context.Context, id uint) (*domain.Stage, error) {
	stage, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return stage, nil
}

func (s *stageService) Create(ctx context.Context, input *dto.CreateStageInput) (*domain.Stage, error) {
	stageInput := &domain.Stage{Name: input.Name}
	stage, err := s.repo.Create(ctx, stageInput)
	if err != nil {
		return nil, err
	}
	return stage, nil
}

func (s *stageService) Update(ctx context.Context, id uint, role *domain.Stage) error {
	if err := s.repo.Update(ctx, id, role); err != nil {
		return err
	}
	return nil
}

func (s *stageService) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

