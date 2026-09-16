package role

import (
	"context"
	"startup_back/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	GetAll(ctx context.Context) ([]*entity.Role, error)
	GetByID(ctx context.Context, id uint) (*entity.Role, error)
	Create(ctx context.Context, role *entity.Role) (*entity.Role, error)
	Update(ctx context.Context, id uint, role *entity.Role) error
	Delete(ctx context.Context, id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAll(ctx context.Context) ([]*entity.Role, error) {
	var roles []*entity.Role
	if err := r.db.WithContext(ctx).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *repository) GetByID(ctx context.Context, id uint) (*entity.Role, error) {
	var role entity.Role
	if err := r.db.WithContext(ctx).First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *repository) Create(ctx context.Context, role *entity.Role) (*entity.Role, error) {
	if err := r.db.WithContext(ctx).Create(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func (r *repository) Update(ctx context.Context, id uint, role *entity.Role) error {
	var existing entity.Role
	if err := r.db.WithContext(ctx).First(&existing, id).Error; err != nil {
		return err
	}
	existing.Name = role.Name
	return r.db.WithContext(ctx).Save(&existing).Error
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Role{}, id).Error
}
