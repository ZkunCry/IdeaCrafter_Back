package service

import (
	"fmt"
	"startup_back/internal/domain"
	"startup_back/internal/dto"
	"startup_back/internal/repository"

	"context"
)
type StartupService interface {
	Create(ctx context.Context, input *dto.CreateStartupInput) (*domain.Startup, error)
  GetByID(ctx context.Context, id uint) (*domain.Startup, error)
  GetAll(ctx context.Context, limit, offset int) ([]*domain.Startup, int64, error)
  Delete(ctx context.Context, id uint) error	
}
type startupService struct{
	repo repository.StartupRepository

}

func NewStartupService(repo repository.StartupRepository) StartupService{
	return &startupService{repo:repo}
}
func (s *startupService) Create(ctx context.Context, input *dto.CreateStartupInput) (*domain.Startup, error) {
	if input == nil {
		return nil, fmt.Errorf("input is nil")
	}

	startup := &domain.Startup{
		Name:             input.Name,
		Description:      input.Description,
		ShortDescription: input.ShortDescription,
		TargetAudience:   input.TargetAudience,
		Problem:          input.Problem,
		Solution:         input.Solution,
		CreatorID:        input.CreatorID,
		StageID:          input.StageID,
	}

	created, err := s.repo.Create(ctx, startup, input.CategoryIDs)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Startup: %v\n\n",created)
	return created, nil
}

func (s * startupService) GetByID(ctx context.Context, id uint) (*domain.Startup, error) {
	startup, err := s.repo.GetByID(ctx,id)
	if err !=nil{
		return nil,err
	}
	return startup,nil
}

func (s * startupService) GetAll(ctx context.Context, limit, offset int) ([]*domain.Startup,int64, error) {
	startups,totalCount, err := s.repo.GetAll(ctx,limit,offset)
	if err !=nil{
		return nil,0,err
	}
	return startups,totalCount, nil
}

func (s * startupService) Delete(ctx context.Context, id uint) error {
	err := s.repo.Delete(ctx,id)
	if err !=nil{
		return err
	}
	return nil
}