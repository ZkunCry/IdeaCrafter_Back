package service

import (
	"fmt"
	"startup_back/internal/domain"
	"startup_back/internal/dto"
	"startup_back/internal/repository"

	"context"
)
type StartupService interface {
	Create(ctx context.Context, startup *dto.CreateStartupInput, categoryIDs []uint) (*domain.Startup, error)
  GetByID(ctx context.Context, id uint) (*domain.Startup, error)
  List(ctx context.Context, limit, offset int) ([]*domain.Startup, error)
  Delete(ctx context.Context, id uint) error	
}
type startupService struct{
	repo repository.StartupRepository

}

func NewStartupService(repo repository.StartupRepository) StartupService{
	return &startupService{repo:repo}
}
func (s * startupService) Create (ctx context.Context, startup *dto.CreateStartupInput, categoryIDs []uint) (*domain.Startup, error) {
	fmt.Println("create")
	if(s.repo == nil){
		return nil,fmt.Errorf("repo is nil")
	}
 startupCreated,err := s.repo.Create(ctx, &domain.Startup{Name: startup.Name, Description: startup.Description, CreatorID: startup.CreatorId}, categoryIDs)
 
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

func (s * startupService) List(ctx context.Context, limit, offset int) ([]*domain.Startup, error) {
	startups, err := s.repo.List(ctx,limit,offset)
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