package service

import (
	"context"
	"fmt"
	"startup_back/internal/domain"
	"startup_back/internal/dto"
	"startup_back/internal/repository"
)


type ApplicationService interface {
	Create(ctx context.Context, input *dto.CreateApplicationInput) (*domain.Application, error)
	Update(ctx context.Context, id uint, input *dto.UpdateApplicationInput) (*domain.Application, error)
	UpdateStatus(ctx context.Context, id uint, input *dto.UpdateApplicationStatusInput) (*domain.Application, error)
	GetByVacancyID(ctx context.Context, vacancyID uint) ([]*domain.Application, error)
	GetByID(ctx context.Context, id uint) (*domain.Application, error)
	Delete(ctx context.Context, id uint) error
}

type applicationService struct {
	repo repository.ApplicationRepository
}

func NewApplicationService(repo repository.ApplicationRepository) ApplicationService {
	return &applicationService{repo: repo}
}


func (s *applicationService) Create(ctx context.Context, input *dto.CreateApplicationInput) (*domain.Application, error) {
	exists, err := s.repo.ExistsByVacancyAndUser(ctx, input.VacancyID, input.UserID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("application already exists for this vacancy")
	}
	application := &domain.Application{
		VacancyID: input.VacancyID,
		UserID:    input.UserID,
		Message:   input.Message,
		Status:    "pending",
	}
	return s.repo.Create(ctx, application)
}

func (s *applicationService) Update(ctx context.Context, id uint, input *dto.UpdateApplicationInput) (*domain.Application, error) {
	application, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	application.Message = input.Message

	return s.repo.Update(ctx, application)
}


func (s *applicationService) UpdateStatus(ctx context.Context, id uint, input *dto.UpdateApplicationStatusInput) (*domain.Application, error) {
	return s.repo.UpdateStatusAndAssign(ctx, id, input.Status)
}


func (s *applicationService) GetByVacancyID(ctx context.Context, vacancyID uint) ([]*domain.Application, error) {
	return s.repo.GetByVacancyID(ctx, vacancyID)
}

func (s *applicationService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *applicationService) GetByID(ctx context.Context, id uint) (*domain.Application, error) {
	return s.repo.GetByID(ctx, id)
}
