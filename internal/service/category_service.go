package service

import (
	"context"
	"startup_back/internal/domain"
	"startup_back/internal/dto"
	"startup_back/internal/repository"
)

type CategoryService interface {
	Create(ctx context.Context, input dto.CreateCategoryInput) (*dto.Category, error)
	Update(ctx context.Context, id uint, input *dto.UpdateCategoryInput) (*dto.Category, error)
	GetById(ctx context.Context, id uint) (*dto.Category, error)
	GetAll(ctx context.Context, searchString string, limit, offset int) ([]dto.Category,int64, error)
	Delete(ctx context.Context, id uint) error
}

type categoryService struct {
	repo repository.CategoryRepository
}


func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s * categoryService) Create(ctx context.Context, input dto.CreateCategoryInput) (*dto.Category, error){
	category, err := s.repo.Create(ctx, &domain.Category{ Name: input.Name})
	if err != nil {
		return nil, err
	}
	categoryDto := &dto.Category{ID: category.ID, Name: category.Name}
	return categoryDto, nil
}

func (s * categoryService) GetById(ctx context.Context,id uint,) (*dto.Category, error) {
	category, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	categoryDto := &dto.Category{ID: category.ID, Name: category.Name}
	return categoryDto, nil
}

func (s * categoryService) GetAll(ctx context.Context, searchString string, limit, offset int) ([]dto.Category, int64, error) {
	categories, totalCount, err := s.repo.GetAll(ctx, searchString, limit, offset)
	if err != nil {
		return nil, 0,err
	}
	categoriesDto := make([]dto.Category, len(categories))
	for i, category := range categories {
		categoriesDto[i] = dto.Category{ID: category.ID, Name: category.Name}
	}
	return categoriesDto, totalCount, nil
}

func (s * categoryService) Update(ctx context.Context, id uint, input *dto.UpdateCategoryInput) (*dto.Category, error) {
	if err := s.repo.Update(ctx, id, &domain.Category{ Name: input.Name}); err != nil {
		return nil, err
	}
	category, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	categoryDto := &dto.Category{ID: category.ID, Name: category.Name}
	return categoryDto, nil
}

func (s * categoryService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}