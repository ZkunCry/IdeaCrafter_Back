package favorite

import (
	"context"
	"startup_back/internal/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	Add(ctx context.Context, userID, startupID uint) error
	Remove(ctx context.Context, userID, startupID uint) error
	GetStartupIDs(ctx context.Context, userID uint) ([]uint, error)
	GetStartups(ctx context.Context, userID uint) ([]*entity.Startup, error)
	CountByStartup(ctx context.Context, startupID uint) (int64, error)
	StartupExists(ctx context.Context, startupID uint) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Add(ctx context.Context, userID, startupID uint) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&entity.Favorite{UserID: userID, StartupID: startupID}).Error
}


func (r *repository) Remove(ctx context.Context, userID, startupID uint) error {
	return r.db.WithContext(ctx).Unscoped().
		Where("user_id = ? AND startup_id = ?", userID, startupID).
		Delete(&entity.Favorite{}).Error
}

func (r *repository) GetStartupIDs(ctx context.Context, userID uint) ([]uint, error) {
	ids := make([]uint, 0)
	if err := r.db.WithContext(ctx).Model(&entity.Favorite{}).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Pluck("startup_id", &ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

func (r *repository) GetStartups(ctx context.Context, userID uint) ([]*entity.Startup, error) {
	var startups []*entity.Startup
	if err := r.db.WithContext(ctx).
		Model(&entity.Startup{}).
		Select("startups.*").
		Joins("JOIN favorites ON favorites.startup_id = startups.id AND favorites.deleted_at IS NULL").
		Where("favorites.user_id = ?", userID).
		Preload("Creator").Preload("Categories").Preload("Stage").Preload("Files").
		Preload("Vacancies").Preload("Vacancies.Role").Preload("Vacancies.User").
		Order("favorites.created_at DESC").
		Find(&startups).Error; err != nil {
		return nil, err
	}
	return startups, nil
}

func (r *repository) CountByStartup(ctx context.Context, startupID uint) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.Favorite{}).
		Where("startup_id = ?", startupID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *repository) StartupExists(ctx context.Context, startupID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.Startup{}).
		Where("id = ?", startupID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
