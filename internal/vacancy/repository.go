package vacancy

import (
	"context"
	"startup_back/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, vacancy *entity.Vacancy) (*entity.Vacancy, error)
	GetByID(ctx context.Context, id uint) (*entity.Vacancy, error)
	Update(ctx context.Context, id uint, description *string, isOpen *bool) (*entity.Vacancy, error)
	Delete(ctx context.Context, id uint) error
	GetByStartupID(ctx context.Context, startupID uint) ([]*entity.Vacancy, error)
	GetAll(ctx context.Context) ([]*entity.Vacancy, error)
	StartupCreatorID(ctx context.Context, startupID uint) (uint, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) preloaded(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Preload("Role").Preload("User")
}

func (r *repository) Create(ctx context.Context, vacancy *entity.Vacancy) (*entity.Vacancy, error) {
	if err := r.db.WithContext(ctx).Create(vacancy).Error; err != nil {
		return nil, err
	}
	if err := r.preloaded(ctx).First(vacancy, vacancy.ID).Error; err != nil {
		return nil, err
	}
	return vacancy, nil
}

func (r *repository) GetByID(ctx context.Context, id uint) (*entity.Vacancy, error) {
	var vacancy entity.Vacancy
	if err := r.preloaded(ctx).First(&vacancy, id).Error; err != nil {
		return nil, err
	}
	return &vacancy, nil
}

func (r *repository) Update(ctx context.Context, id uint, description *string, isOpen *bool) (*entity.Vacancy, error) {
	var existing entity.Vacancy
	if err := r.db.WithContext(ctx).First(&existing, id).Error; err != nil {
		return nil, err
	}
	if description != nil {
		existing.Description = *description
	}
	if isOpen != nil {
		existing.IsOpen = *isOpen
	}
	if err := r.db.WithContext(ctx).Save(&existing).Error; err != nil {
		return nil, err
	}
	if err := r.preloaded(ctx).First(&existing, id).Error; err != nil {
		return nil, err
	}
	return &existing, nil
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Vacancy{}, id).Error
}

func (r *repository) GetByStartupID(ctx context.Context, startupID uint) ([]*entity.Vacancy, error) {
	var vacancies []*entity.Vacancy
	if err := r.preloaded(ctx).Where("startup_id = ?", startupID).Order("created_at ASC").Find(&vacancies).Error; err != nil {
		return nil, err
	}
	return vacancies, nil
}

func (r *repository) GetAll(ctx context.Context) ([]*entity.Vacancy, error) {
	var vacancies []*entity.Vacancy
	if err := r.preloaded(ctx).Find(&vacancies).Error; err != nil {
		return nil, err
	}
	return vacancies, nil
}

func (r *repository) StartupCreatorID(ctx context.Context, startupID uint) (uint, error) {
	var startup entity.Startup
	if err := r.db.WithContext(ctx).Select("id", "creator_id").First(&startup, startupID).Error; err != nil {
		return 0, err
	}
	return startup.CreatorID, nil
}
