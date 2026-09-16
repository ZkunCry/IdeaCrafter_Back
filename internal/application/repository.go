package application

import (
	"context"
	"fmt"
	"startup_back/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, app *entity.Application) (*entity.Application, error)
	Update(ctx context.Context, application *entity.Application) (*entity.Application, error)
	GetByID(ctx context.Context, id uint) (*entity.Application, error)
	GetByVacancyID(ctx context.Context, vacancyID uint) ([]*entity.Application, error)
	GetByStartupID(ctx context.Context, startupID uint) ([]*entity.Application, error)
	GetByUserID(ctx context.Context, userID uint) ([]*entity.Application, error)
	UpdateStatusAndAssign(ctx context.Context, id uint, status string) (*entity.Application, error)
	Delete(ctx context.Context, id uint) error
	ExistsByVacancyAndUser(ctx context.Context, vacancyID, userID uint) (bool, error)
	StartupCreatorID(ctx context.Context, startupID uint) (uint, error)
	StartupCreatorIDByVacancy(ctx context.Context, vacancyID uint) (uint, error)
	GetStartupsByIDs(ctx context.Context, ids []uint) ([]entity.Startup, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) preloaded(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).
		Preload("User").
		Preload("Vacancy", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).
		Preload("Vacancy.Role").
		Preload("Vacancy.User")
}

func (r *repository) Create(ctx context.Context, app *entity.Application) (*entity.Application, error) {
	if err := r.db.WithContext(ctx).Create(app).Error; err != nil {
		return nil, err
	}
	if err := r.preloaded(ctx).First(app, app.ID).Error; err != nil {
		return nil, err
	}
	return app, nil
}

func (r *repository) Update(ctx context.Context, application *entity.Application) (*entity.Application, error) {
	if err := r.db.WithContext(ctx).Save(application).Error; err != nil {
		return nil, err
	}
	if err := r.preloaded(ctx).First(application, application.ID).Error; err != nil {
		return nil, err
	}
	return application, nil
}

func (r *repository) GetByID(ctx context.Context, id uint) (*entity.Application, error) {
	var app entity.Application
	if err := r.preloaded(ctx).First(&app, id).Error; err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *repository) GetByVacancyID(ctx context.Context, vacancyID uint) ([]*entity.Application, error) {
	var apps []*entity.Application
	if err := r.preloaded(ctx).
		Where("vacancy_id = ?", vacancyID).
		Order("created_at DESC").
		Find(&apps).Error; err != nil {
		return nil, err
	}
	return apps, nil
}

func (r *repository) GetByStartupID(ctx context.Context, startupID uint) ([]*entity.Application, error) {
	var apps []*entity.Application
	if err := r.preloaded(ctx).
		Model(&entity.Application{}).
		Select("applications.*").
		Joins("JOIN vacancies ON vacancies.id = applications.vacancy_id").
		Where("vacancies.startup_id = ? AND vacancies.deleted_at IS NULL", startupID).
		Order("applications.created_at DESC").
		Find(&apps).Error; err != nil {
		return nil, err
	}
	return apps, nil
}

func (r *repository) GetByUserID(ctx context.Context, userID uint) ([]*entity.Application, error) {
	var apps []*entity.Application
	if err := r.preloaded(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&apps).Error; err != nil {
		return nil, err
	}
	return apps, nil
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Application{}, id).Error
}


func (r *repository) UpdateStatusAndAssign(ctx context.Context, id uint, status string) (*entity.Application, error) {
	var app entity.Application
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&app, id).Error; err != nil {
			return err
		}
		app.Status = status
		if err := tx.Save(&app).Error; err != nil {
			return err
		}
		if status != StatusAccepted {
			return nil
		}

		var vacancy entity.Vacancy
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

		return tx.Model(&entity.Application{}).
			Where("vacancy_id = ? AND id <> ? AND status = ?", app.VacancyID, app.ID, StatusPending).
			Update("status", StatusRejected).Error
	})
	if err != nil {
		return nil, err
	}
	if err := r.preloaded(ctx).First(&app, id).Error; err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *repository) ExistsByVacancyAndUser(ctx context.Context, vacancyID, userID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.Application{}).
		Where("vacancy_id = ? AND user_id = ?", vacancyID, userID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *repository) StartupCreatorID(ctx context.Context, startupID uint) (uint, error) {
	var startup entity.Startup
	if err := r.db.WithContext(ctx).Select("id", "creator_id").First(&startup, startupID).Error; err != nil {
		return 0, err
	}
	return startup.CreatorID, nil
}

func (r *repository) StartupCreatorIDByVacancy(ctx context.Context, vacancyID uint) (uint, error) {
	var vacancy entity.Vacancy
	if err := r.db.WithContext(ctx).Select("id", "startup_id").First(&vacancy, vacancyID).Error; err != nil {
		return 0, err
	}
	return r.StartupCreatorID(ctx, vacancy.StartupID)
}

// GetStartupsByIDs includes soft-deleted startups so an applicant's history
// still names the project after it was removed.
func (r *repository) GetStartupsByIDs(ctx context.Context, ids []uint) ([]entity.Startup, error) {
	var startups []entity.Startup
	if len(ids) == 0 {
		return startups, nil
	}
	if err := r.db.WithContext(ctx).Unscoped().
		Select("id", "name", "short_description", "logo_url", "deleted_at").
		Where("id IN ?", ids).
		Find(&startups).Error; err != nil {
		return nil, err
	}
	return startups, nil
}
