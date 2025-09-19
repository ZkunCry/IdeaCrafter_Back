package service

import (
	"startup_back/internal/domain"
	"startup_back/internal/dto"
	"startup_back/internal/repository"

	"context"
)

type startupService struct{
	repo repository.StartupRepository

}

func NewStartupService(repo repository.StartupRepository) StartupService{
	return &startupService{repo:repo}
}
func (s * startupService) Create (ctx context.Context, startup *dto.CreateStartupInput, categoryIDs []uint, vacancyRoleIDs []uint) (*domain.Startup, error) {
 startupCreated,err := s.repo.Create(ctx, &domain.Startup{Name: startup.Name, Description: startup.Description}, categoryIDs, vacancyRoleIDs)
 if err != nil {
	return nil,err
 }
return startupCreated,nil
}

func (s * startupService) GetByID(ctx context.Context, id uint) (*domain.Startup, error) {
	startup, err := s.repo.GetByID(ctx,id)
	if err !=nil{
		return nil,err
	}
	return startup,nil
}

func (s * startupService) List(ctx context.Context, limit, offset int, categoryID uint) ([]*domain.Startup, error) {
	startups, err := s.repo.List(ctx,limit,offset,categoryID)
	if err !=nil{
		return nil,err
	}
	return startups, nil
}

func (s * startupService) Delete(ctx context.Context, id uint) error {
	err := s.repo.Delete(ctx,id)
	if err !=nil{
		return err
	}
	return nil
}