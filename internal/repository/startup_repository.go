package repository

import (
	"context"
	"fmt"
	"startup_back/internal/domain"

	"gorm.io/gorm"
)

type StartupRepository interface {
	Create(ctx context.Context, startup *domain.Startup, categoryIDs []uint) (*domain.Startup, error)
	GetByID(ctx context.Context, id uint) (*domain.Startup, error)
	GetAll(ctx context.Context, searchString string, limit, offset int) ([]*domain.Startup, int, error)
	Delete(ctx context.Context, id uint) error
	GetUserStartups(ctx context.Context, userID uint) ([]domain.Startup, error)
}
type startupRepository struct {
	db *gorm.DB
}

func (s *startupRepository) GetUserStartups(ctx context.Context, userID uint) ([]domain.Startup, error) {
	var startups [] domain.Startup
	if err := s.db.WithContext(ctx).
		Preload("Categories").
		Preload("Creator").
		Preload("Stage").
		Preload("Files").
		Where("creator_id = ?", userID).
		Find(&startups).Error; err != nil {
		return nil, err
	}
	return startups, nil
}

func NewStartupRepository(db *gorm.DB) StartupRepository {
	return &startupRepository{db: db}
}

func (s *startupRepository) Create(ctx context.Context, startup *domain.Startup, categoryIDs []uint) (*domain.Startup, error) {
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(startup).Error; err != nil {
			return err
		}
		if len(categoryIDs) > 0 {
			var categories []domain.Category
			if err := tx.Where("id IN ?", categoryIDs).Find(&categories).Error; err != nil {
				return err
			}
			if err := tx.Model(startup).Association("Categories").Replace(categories); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	if err := s.db.WithContext(ctx).
		Preload("Categories").
		Preload("Creator").
		Preload("Stage").
		Preload("Files").
		First(startup, startup.ID).Error; err != nil {
		return nil, err
	}
	return startup, nil
}

func (s *startupRepository) GetByID(ctx context.Context, id uint) (*domain.Startup, error) {
	var startup domain.Startup
	err := s.db.Where("id = ?", id).Preload("Categories").Preload("Creator").Preload("Stage").Preload("Files").Preload("Vacancies").First(&startup).Error
	if err != nil {
		return nil, err
	}
	return &startup, nil
}
func (s *startupRepository) GetAll(ctx context.Context, searchString string, limit, offset int) ([]*domain.Startup, int, error) {
	var startups []*domain.Startup
	var totalCount int64
	query := s.db.WithContext(ctx).Model(&domain.Startup{})
	fmt.Printf("SEARCH STRING %v", searchString)
	if searchString != "" {
		fmt.Println("TWESTSD")
		searchPattern := "%" + searchString + "%"
		query = query.Where("LOWER(name) LIKE LOWER(?) OR LOWER(description) LIKE LOWER(?)", searchPattern, searchPattern)

	}
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	sqlOffset := offset * limit
	if err := query.
		Preload("Creator").
		Preload("Categories").
		Preload("Vacancies").
		Preload("Files").
		Preload("Stage").
		Order("created_at DESC").
		Limit(limit).
		Offset(sqlOffset).
		Find(&startups).Error; err != nil {
		return nil, 0, err
	}

	return startups, int(totalCount), nil
}
func (s *startupRepository) Delete(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&domain.Startup{}, id).Error
}
