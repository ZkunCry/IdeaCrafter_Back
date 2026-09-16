package startup

import (
	"context"
	"errors"
	"fmt"
	"startup_back/internal/entity"
)

var ErrForbidden = errors.New("only the startup owner can edit it")

type Service interface {
	Create(ctx context.Context, input CreateInput) (*entity.Startup, error)
	GetByID(ctx context.Context, id uint) (*entity.Startup, error)
	GetAll(ctx context.Context, searchString, categorySlug string, limit, offset int) ([]*entity.Startup, int, error)
	Update(ctx context.Context, id uint, input UpdateInput, actorID uint) (*entity.Startup, error)
	Delete(ctx context.Context, id uint) error
	GetUserStartups(ctx context.Context, userID uint) ([]entity.Startup, error)
	AddCategories(ctx context.Context, startupID uint, input AddCategoriesInput) (*entity.Startup, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, input CreateInput) (*entity.Startup, error) {
	startup := &entity.Startup{
		Name:             input.Name,
		Description:      input.Description,
		ShortDescription: input.ShortDescription,
		TargetAudience:   input.TargetAudience,
		Problem:          input.Problem,
		Solution:         input.Solution,
		CreatorID:        input.CreatorID,
		StageID:          input.StageID,
		LogoURL:          input.LogoFile,
	}
	return s.repo.Create(ctx, startup, input.CategoryIDs)
}

func (s *service) GetByID(ctx context.Context, id uint) (*entity.Startup, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetAll(ctx context.Context, searchString, categorySlug string, limit, offset int) ([]*entity.Startup, int, error) {
	return s.repo.GetAll(ctx, searchString, categorySlug, limit, offset)
}

func (s *service) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) GetUserStartups(ctx context.Context, userID uint) ([]entity.Startup, error) {
	return s.repo.GetUserStartups(ctx, userID)
}

func (s *service) AddCategories(ctx context.Context, startupID uint, input AddCategoriesInput) (*entity.Startup, error) {
	if len(input.CategoryIDs) == 0 {
		return nil, fmt.Errorf("category_ids are required")
	}
	return s.repo.AddCategories(ctx, startupID, input.CategoryIDs)
}

func (s *service) Update(ctx context.Context, id uint, input UpdateInput, actorID uint) (*entity.Startup, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing.CreatorID != actorID {
		return nil, ErrForbidden
	}
	if input.Name == "" || input.ShortDescription == "" || input.Description == "" {
		return nil, fmt.Errorf("name, short_description and description are required")
	}
	return s.repo.Update(ctx, id, input)
}
