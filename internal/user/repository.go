package user

import (
	"context"
	"startup_back/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	Update(ctx context.Context, id uint, user *entity.User) error
	Delete(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *repository) Update(ctx context.Context, id uint, user *entity.User) error {
	return r.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("id = ?", id).
		Updates(user).Error
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.User{}, id).Error
}

func (r *repository) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
