package repository

import (
	"context"
	"startup_back/internal/domain"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uint) error

}

type StartupRepository interface {
    Create(ctx context.Context, startup *domain.Startup, categoryIDs []uint) error
    GetByID(ctx context.Context, id uint) (*domain.Startup, error)
    List(ctx context.Context, limit, offset int, categoryID uint) ([]*domain.Startup, error)
    Delete(ctx context.Context, id uint) error
}

type Repositories struct {
    User          UserRepository
    Startup       StartupRepository

}

func NewRespositories(db *gorm.DB) * Repositories{
	return &Repositories{
		User : NewUserRepository(db),
	}
}