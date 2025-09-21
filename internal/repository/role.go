package repository

import (
	"context"
	"startup_back/internal/domain"

	"gorm.io/gorm"
)

type RoleRepository interface {
	GetAll(ctx context.Context) ([]*domain.Role, error)
	GetByID(ctx context.Context, id uint) (*domain.Role, error)
	Create(ctx context.Context, role *domain.Role) (*domain.Role, error)
	Update(ctx context.Context, id uint, role *domain.Role) error
	Delete(ctx context.Context, id uint) error
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) GetAll(ctx context.Context) ([]*domain.Role, error) {
	var roles []*domain.Role
	if err := r.db.WithContext(ctx).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r * roleRepository) Create(ctx context.Context, role *domain.Role) (*domain.Role, error) {
	if err := r.db.WithContext(ctx).Create(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func (r *roleRepository) GetByID(ctx context.Context, id uint) (*domain.Role, error) {
	var role domain.Role
	if err := r.db.WithContext(ctx).First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) Update(ctx context.Context, id uint, role *domain.Role) error {
	var existing domain.Role
	if err := r.db.WithContext(ctx).First(&existing, id).Error; err != nil {
		return err
	}
	existing.Name = role.Name
	if err := r.db.WithContext(ctx).Save(&existing).Error; err != nil {
		return err
	}
	return nil
}

func (r *roleRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&domain.Role{}, id).Error; err != nil {
		return err
	}
	return nil
}