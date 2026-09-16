package category

import (
	"context"
	"startup_back/internal/entity"
)

type Service interface {
	Create(ctx context.Context, input CreateInput) (*Response, error)
	Update(ctx context.Context, id uint, input *UpdateInput) (*Response, error)
	GetByID(ctx context.Context, id uint) (*Response, error)
	GetAll(ctx context.Context, searchString string, limit, offset int) ([]Response, int64, error)
	Delete(ctx context.Context, id uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, input CreateInput) (*Response, error) {
	category, err := s.repo.Create(ctx, &entity.Category{Name: input.Name, Slug: input.Slug})
	if err != nil {
		return nil, err
	}
	return &Response{ID: category.ID, Name: category.Name, Slug: category.Slug, Count: category.StartupCount}, nil
}

func (s *service) Update(ctx context.Context, id uint, input *UpdateInput) (*Response, error) {
	if err := s.repo.Update(ctx, id, &entity.Category{Name: input.Name, Slug: input.Slug}); err != nil {
		return nil, err
	}
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &Response{ID: category.ID, Name: category.Name, Slug: category.Slug, Count: category.StartupCount}, nil
}

func (s *service) GetByID(ctx context.Context, id uint) (*Response, error) {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &Response{ID: category.ID, Name: category.Name, Slug: category.Slug, Count: category.StartupCount}, nil
}

func (s *service) GetAll(ctx context.Context, searchString string, limit, offset int) ([]Response, int64, error) {
	categories, totalCount, err := s.repo.GetAll(ctx, searchString, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	responses := make([]Response, len(categories))
	for i, category := range categories {
		responses[i] = Response{ID: category.ID, Name: category.Name, Slug: category.Slug, Count: category.StartupCount}
	}
	return responses, totalCount, nil
}

func (s *service) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
