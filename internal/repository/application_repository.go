package repository

import (
	"context"
	"startup_back/internal/domain"

	"gorm.io/gorm"
)

type ApplicationRepository interface {
	Create(ctx context.Context, app *domain.Application) (*domain.Application, error)
	Update(ctx context.Context, application *domain.Application) (*domain.Application, error)
	GetByID(ctx context.Context, id uint) (*domain.Application, error)
	GetByVacancyID(ctx context.Context, vacancyID uint) ([]*domain.Application, error)
	UpdateStatus(ctx context.Context, id uint, status string) (*domain.Application, error)
	Delete(ctx context.Context, id uint) error
}
type applicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) ApplicationRepository {
	return &applicationRepository{db: db}
}


func (r *applicationRepository) Create(ctx context.Context, app *domain.Application) (*domain.Application, error) {
	if err := r.db.WithContext(ctx).Create(app).Error; err != nil {
		return nil, err
	}
	return app, nil
}

func (r *applicationRepository) Update(ctx context.Context, application *domain.Application) (*domain.Application, error) {
	if err := r.db.WithContext(ctx).Save(application).Error; err != nil {
		return nil, err
	}
	return application, nil
}

func (r *applicationRepository) GetByID(ctx context.Context, id uint) (*domain.Application, error) {
	var app domain.Application
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Vacancy").
		First(&app, id).Error; err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *applicationRepository) GetByVacancyID(ctx context.Context, vacancyID uint) ([]*domain.Application, error) {
	var apps []*domain.Application
	if err := r.db.WithContext(ctx).
		Where("vacancy_id = ?", vacancyID).
		Preload("User").
		Find(&apps).Error; err != nil {
		return nil, err
	}
	return apps, nil
}

func (r *applicationRepository) UpdateStatus(ctx context.Context, id uint, status string) (*domain.Application, error) {
	var app domain.Application
	if err := r.db.WithContext(ctx).First(&app, id).Error; err != nil {
		return nil, err
	}
	app.Status = status
	if err := r.db.WithContext(ctx).Save(&app).Error; err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *applicationRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Application{}, id).Error
}