package stage

import (
	"context"
	"startup_back/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	GetAll(ctx context.Context) ([]*entity.Stage, error)
	GetByID(ctx context.Context, id uint) (*entity.Stage, error)
	Create(ctx context.Context, stage *entity.Stage) (*entity.Stage, error)
	Update(ctx context.Context, id uint, stage *entity.Stage) error
	Delete(ctx context.Context, id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAll(ctx context.Context) ([]*entity.Stage, error) {
	var stages []*entity.Stage
	if err := r.db.WithContext(ctx).Find(&stages).Error; err != nil {
		return nil, err
	}
	return stages, nil
}

func (r *repository) GetByID(ctx context.Context, id uint) (*entity.Stage, error) {
	var stage entity.Stage
	if err := r.db.WithContext(ctx).First(&stage, id).Error; err != nil {
		return nil, err
	}
	return &stage, nil
}

func (r *repository) Create(ctx context.Context, stage *entity.Stage) (*entity.Stage, error) {
	if err := r.db.WithContext(ctx).Create(stage).Error; err != nil {
		return nil, err
	}
	return stage, nil
}

func (r *repository) Update(ctx context.Context, id uint, stage *entity.Stage) error {
	var existing entity.Stage
	if err := r.db.WithContext(ctx).First(&existing, id).Error; err != nil {
		return err
	}
	existing.Name = stage.Name
	return r.db.WithContext(ctx).Save(&existing).Error
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Stage{}, id).Error
}
