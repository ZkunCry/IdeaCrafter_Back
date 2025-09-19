package repository

import (
	"context"
	"startup_back/internal/domain"

	"gorm.io/gorm"
)

type startupRepository struct {
	db *gorm.DB
}

func NewStartupRepository(db *gorm.DB) StartupRepository {
	return &startupRepository{db: db}
}

func (s *startupRepository) Create(ctx context.Context, startup *domain.Startup, categoryIDs []uint, vacancyRoleIDs []uint) (*domain.Startup, error) {
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

		if len(vacancyRoleIDs) > 0 {
			var vacancies []domain.Vacancy
			for _, roleID := range vacancyRoleIDs {
				vacancies = append(vacancies, domain.Vacancy{
					RoleID:    roleID,
					StartupID: startup.ID,
				})
			}
			if err := tx.Create(&vacancies).Error; err != nil {
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
		Preload("Vacancies.Role").
		Preload("Files").
		First(startup, startup.ID).Error; err != nil {
		return nil, err
	}

	return startup, nil
}

func (s *startupRepository) GetByID(ctx context.Context, id uint) (*domain.Startup, error){
	var startup domain.Startup
	err := s.db.Where("id = ?", id).First(&startup).Error
	if err !=nil{
		return nil,err
	}
	return &startup, nil
}
func (s *startupRepository)  List(ctx context.Context, limit, offset int, categoryID uint) ([]*domain.Startup, error){
	var startups []*domain.Startup
	query := s.db.WithContext(ctx).
	Preload("Creator").
	Preload("Categories").
	Preload("Vacancies").
	Preload("Files").
	Limit(limit).
	Offset(offset)
	if categoryID != 0 {
		query = query.Joins(
			"JOIN startup_categories sc ON sc.startup_id = startups.id").
			Where("sc.category_id = ?",categoryID)
	}

	if err := query.Find(&startups).Error; err != nil {
		return nil, err
	}
	return startups, nil
}
func (s *startupRepository)  Delete(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&domain.Startup{},id).Error
}