package repository

import (
	"context"
	"fmt"
	"startup_back/internal/domain"

	"gorm.io/gorm"
)

type ApplicationRepository interface {
	Create(ctx context.Context, app *domain.Application) (*domain.Application, error)
	Update(ctx context.Context, application *domain.Application) (*domain.Application, error)
	GetByID(ctx context.Context, id uint) (*domain.Application, error)
	GetByVacancyID(ctx context.Context, vacancyID uint) ([]*domain.Application, error)
	UpdateStatus(ctx context.Context, id uint, status string) (*domain.Application, error)
	UpdateStatusAndAssign(ctx context.Context, id uint, status string) (*domain.Application, error)
	Delete(ctx context.Context, id uint) error
	ExistsByVacancyAndUser(ctx context.Context, vacancyID, userID uint) (bool, error)
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
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Vacancy").
		Preload("Vacancy.Role").
		First(app, app.ID).Error; err != nil {
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
		Preload("Vacancy.Role").
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
		Preload("Vacancy").
		Preload("Vacancy.Role").
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

func (r *applicationRepository) UpdateStatusAndAssign(ctx context.Context, id uint, status string) (*domain.Application, error) {
	var app domain.Application

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&app, id).Error; err != nil {
			return err
		}
		app.Status = status
		if err := tx.Save(&app).Error; err != nil {
			return err
		}

		if status == "accepted" {
			var vacancy domain.Vacancy
			if err := tx.First(&vacancy, app.VacancyID).Error; err != nil {
				return err
			}
			if vacancy.UserID != nil {
				return fmt.Errorf("vacancy already assigned")
			}
			vacancy.UserID = &app.UserID
			vacancy.IsOpen = false
			if err := tx.Save(&vacancy).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Vacancy").
		Preload("Vacancy.Role").
		First(&app, id).Error; err != nil {
		return nil, err
	}

	return &app, nil
}

func (r *applicationRepository) ExistsByVacancyAndUser(ctx context.Context, vacancyID, userID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&domain.Application{}).
		Where("vacancy_id = ? AND user_id = ?", vacancyID, userID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
