package repository

import (
	"context"
	"startup_back/internal/domain"

	"gorm.io/gorm"
)

type StageRepository interface {
	GetAll(ctx context.Context) ([]*domain.Stage, error)
	GetByID(ctx context.Context, id uint) (*domain.Stage, error)
	Create(ctx context.Context, role *domain.Stage) (*domain.Stage, error)
	Update(ctx context.Context, id uint, role *domain.Stage) error
	Delete(ctx context.Context, id uint) error
}

type stageRepository struct {
	db *gorm.DB
}

func NewStageRepository(db *gorm.DB) StageRepository {
	return &stageRepository{db: db}
}

func (r *stageRepository) GetAll(ctx context.Context) ([]*domain.Stage, error) {
	var stages []*domain.Stage
	if err := r.db.WithContext(ctx).Find(&stages).Error; err != nil {
		return nil, err
	}
	return stages, nil
}

func (r * stageRepository) GetByID(ctx context.Context, id uint) (*domain.Stage, error) {
	var stage domain.Stage
	if err := r.db.WithContext(ctx).First(&stage, id).Error; err != nil {
		return nil, err
	}
	return &stage, nil
}

func (r *stageRepository) Create(ctx context.Context, stage *domain.Stage) (*domain.Stage, error) {
	if err := r.db.WithContext(ctx).Create(stage).Error; err != nil {
		return nil, err
	}
	return stage, nil
}

func (r *stageRepository) Update(ctx context.Context, id uint, stage *domain.Stage) error {
	var existing domain.Stage
	if err := r.db.WithContext(ctx).First(&existing, id).Error; err != nil {
		return err
	}
	existing.Name = stage.Name
	if err := r.db.WithContext(ctx).Save(&existing).Error; err != nil {
		return err
	}
	return nil
}

func (r *stageRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&domain.Stage{}, id).Error; err != nil {
		return err
	}
	return nil
}