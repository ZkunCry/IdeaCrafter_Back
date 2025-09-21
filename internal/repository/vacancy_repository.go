package repository

import (
	"context"
	"startup_back/internal/domain"

	"gorm.io/gorm"
)

type vacancyRepository struct {
	db *gorm.DB
}

func NewVacancyRepository(db *gorm.DB) VacancyRepository {
	return &vacancyRepository{db: db}
}

func (v *vacancyRepository) Create(ctx context.Context, vacancy *domain.Vacancy) (*domain.Vacancy, error) {
	if err := v.db.WithContext(ctx).Create(vacancy).Error; err != nil {
		return nil, err
	}
	if err := v.db.WithContext(ctx).
        Preload("Role").
        First(vacancy, vacancy.ID).Error; err != nil {
        return nil, err
    }
	return vacancy, nil
}

func (v *vacancyRepository) GetByID(ctx context.Context, id uint) (*domain.Vacancy, error) {
	var vacancy domain.Vacancy
	if err := v.db.WithContext(ctx).First(&vacancy, id).Error; err != nil {
		return nil, err
	}
	return &vacancy, nil
}

func (v *vacancyRepository) Update(ctx context.Context, id uint, vacancy *domain.Vacancy) (*domain.Vacancy, error) {
    var existing domain.Vacancy
    if err := v.db.WithContext(ctx).First(&existing, id).Error; err != nil {
        return nil, err
    }
    existing.Description = vacancy.Description
    existing.IsOpen = vacancy.IsOpen
    if err := v.db.WithContext(ctx).Save(&existing).Error; err != nil {
        return nil, err
    }

    return &existing, nil
}

func (v *vacancyRepository) Delete(ctx context.Context, id uint) error {
	if err := v.db.WithContext(ctx).Delete(&domain.Vacancy{}, id).Error; err != nil {
		return err
	}
	return nil
}
 
func (v *vacancyRepository) GetByStartupID(ctx context.Context, startupID uint) ([] *domain.Vacancy, error) {
	var vacancies []*domain.Vacancy
	if err := v.db.WithContext(ctx).Where("startup_id = ?", startupID).Find(&vacancies).Error; err != nil {
		return nil, err
	}
	return vacancies, nil
}

func (v *vacancyRepository) GetAll(ctx context.Context) ([] *domain.Vacancy, error) {
	var vacancies []*domain.Vacancy
	if err := v.db.WithContext(ctx).Find(&vacancies).Error; err != nil {
		return nil, err
	}
	return vacancies, nil
}