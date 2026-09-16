package vacancy

import (
	"context"
	"errors"
	"startup_back/internal/entity"
)

// ErrForbidden is returned when the acting user does not own the startup the
// vacancy belongs to.
var ErrForbidden = errors.New("only the startup owner can manage its vacancies")

type Service interface {
	Create(ctx context.Context, input *CreateInput, actorID uint) (*entity.Vacancy, error)
	GetByID(ctx context.Context, id uint) (*entity.Vacancy, error)
	Update(ctx context.Context, id uint, input *UpdateInput, actorID uint) (*entity.Vacancy, error)
	Delete(ctx context.Context, id uint, actorID uint) error
	GetByStartupID(ctx context.Context, startupID uint) ([]*entity.Vacancy, error)
	GetAll(ctx context.Context) ([]*entity.Vacancy, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) assertOwner(ctx context.Context, startupID, actorID uint) error {
	creatorID, err := s.repo.StartupCreatorID(ctx, startupID)
	if err != nil {
		return err
	}
	if creatorID != actorID {
		return ErrForbidden
	}
	return nil
}

func (s *service) Create(ctx context.Context, input *CreateInput, actorID uint) (*entity.Vacancy, error) {
	if input.StartupID == 0 || input.RoleID == 0 {
		return nil, errors.New("startup_id and role_id are required")
	}
	if err := s.assertOwner(ctx, input.StartupID, actorID); err != nil {
		return nil, err
	}
	vacancy := &entity.Vacancy{
		StartupID:   input.StartupID,
		RoleID:      input.RoleID,
		Description: input.Description,
		IsOpen:      true,
	}
	return s.repo.Create(ctx, vacancy)
}

func (s *service) GetByID(ctx context.Context, id uint) (*entity.Vacancy, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) Update(ctx context.Context, id uint, input *UpdateInput, actorID uint) (*entity.Vacancy, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := s.assertOwner(ctx, existing.StartupID, actorID); err != nil {
		return nil, err
	}
	if input.IsOpen != nil && *input.IsOpen && existing.UserID != nil {
		return nil, errors.New("vacancy already has an assigned member")
	}
	return s.repo.Update(ctx, id, input.Description, input.IsOpen)
}

func (s *service) Delete(ctx context.Context, id uint, actorID uint) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.assertOwner(ctx, existing.StartupID, actorID); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}

func (s *service) GetByStartupID(ctx context.Context, startupID uint) ([]*entity.Vacancy, error) {
	return s.repo.GetByStartupID(ctx, startupID)
}

func (s *service) GetAll(ctx context.Context) ([]*entity.Vacancy, error) {
	return s.repo.GetAll(ctx)
}
