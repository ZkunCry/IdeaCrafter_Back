package service

import (
	"context"
	"errors"
	"fmt"
	"startup_back/internal/domain"
	"startup_back/internal/dto"
	"startup_back/internal/repository"
)
type VacancyService interface{
	Create(ctx context.Context, input *dto.CreateVacancyInput) (*domain.Vacancy, error)
	GetByID(ctx context.Context, id uint) (*domain.Vacancy, error)
	Update(ctx context.Context, id uint, input *dto.UpdateVacancyInput) ( *domain.Vacancy, error)
	Delete(ctx context.Context, id uint) error

	GetByStartupID(ctx context.Context, startupID uint) ([]*domain.Vacancy, error)
	GetAll(ctx context.Context) ([]*domain.Vacancy, error)
}
type vacancyService struct{
	repo repository.VacancyRepository
}
func NewVacancyService(repo  repository.VacancyRepository) VacancyService{
	return &vacancyService{repo: repo}
}

func (v *vacancyService) Create(ctx context.Context, input *dto.CreateVacancyInput) (*domain.Vacancy, error){
    if input.StartupID == 0 || input.RoleID == 0 {
        return nil, errors.New("startup_id and role_id are required")
    }
    vacancy := &domain.Vacancy{
        StartupID:   input.StartupID,
        RoleID:      input.RoleID,
        Description: input.Description,
        IsOpen:      true,
    }

		created, err := v.repo.Create(ctx, vacancy)
		if err != nil {
			fmt.Print("ERROR NOT NIL")
				return nil, err
		}	
		fmt.Print(created)
    return created,nil
}

func (v *vacancyService) GetByID(ctx context.Context, id uint) (*domain.Vacancy, error){
	vacancy, err:= v.repo.GetByID(ctx, id)
	if err != nil{
		return nil, err
	}
	return vacancy, nil
}

func (v *vacancyService) Update(ctx context.Context, id uint, input *dto.UpdateVacancyInput) ( *domain.Vacancy, error){
	 vacancy := &domain.Vacancy{}
    if input.Description != "" {
        vacancy.Description = input.Description
    }
    if input.IsOpen != nil {
        vacancy.IsOpen = *input.IsOpen
    }

    updated, err := v.repo.Update(ctx, id, vacancy)
    if err != nil {
        return nil, err
    }

    return updated, nil
}

func (v *vacancyService) Delete(ctx context.Context, id uint) error{
	err := v.repo.Delete(ctx, id)
	if err != nil{
		return err
	}
	return nil
}

func (v *vacancyService) GetByStartupID(ctx context.Context, startupID uint) ([]*domain.Vacancy, error){
	vacancies, err := v.repo.GetByStartupID(ctx, startupID)
	if err != nil{
		return nil, err
	}
	return vacancies, nil
}

func (v *vacancyService) GetAll(ctx context.Context) ([]*domain.Vacancy, error){
	vacancies, err := v.repo.GetAll(ctx)
	if err != nil{
		return nil, err
	}
	return vacancies, nil
}