package favorite

import (
	"context"
	"errors"
	"startup_back/internal/entity"
)

// ErrStartupNotFound is returned when the favorited startup does not exist.
var ErrStartupNotFound = errors.New("startup not found")

type Service interface {
	Add(ctx context.Context, userID, startupID uint) (int64, error)
	Remove(ctx context.Context, userID, startupID uint) (int64, error)
	GetStartupIDs(ctx context.Context, userID uint) ([]uint, error)
	GetStartups(ctx context.Context, userID uint) ([]*entity.Startup, error)
	CountByStartup(ctx context.Context, startupID uint) (int64, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// Add returns the fresh favorites count so the UI can update without a refetch.
func (s *service) Add(ctx context.Context, userID, startupID uint) (int64, error) {
	exists, err := s.repo.StartupExists(ctx, startupID)
	if err != nil {
		return 0, err
	}
	if !exists {
		return 0, ErrStartupNotFound
	}
	if err := s.repo.Add(ctx, userID, startupID); err != nil {
		return 0, err
	}
	return s.repo.CountByStartup(ctx, startupID)
}

func (s *service) Remove(ctx context.Context, userID, startupID uint) (int64, error) {
	if err := s.repo.Remove(ctx, userID, startupID); err != nil {
		return 0, err
	}
	return s.repo.CountByStartup(ctx, startupID)
}

func (s *service) GetStartupIDs(ctx context.Context, userID uint) ([]uint, error) {
	return s.repo.GetStartupIDs(ctx, userID)
}

func (s *service) GetStartups(ctx context.Context, userID uint) ([]*entity.Startup, error) {
	return s.repo.GetStartups(ctx, userID)
}

func (s *service) CountByStartup(ctx context.Context, startupID uint) (int64, error) {
	return s.repo.CountByStartup(ctx, startupID)
}
