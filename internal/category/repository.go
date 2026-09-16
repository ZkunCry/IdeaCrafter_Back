package category

import (
	"context"
	"startup_back/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, category *entity.Category) (*entity.Category, error)
	Update(ctx context.Context, id uint, category *entity.Category) error
	Delete(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (*entity.Category, error)
	GetAll(ctx context.Context, searchString string, limit, offset int) ([]entity.Category, int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, category *entity.Category) (*entity.Category, error) {
	if err := r.db.WithContext(ctx).Create(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (r *repository) Update(ctx context.Context, id uint, category *entity.Category) error {
	var existing entity.Category
	if err := r.db.WithContext(ctx).First(&existing, id).Error; err != nil {
		return err
	}
	existing.Name = category.Name
	existing.Slug = category.Slug
	return r.db.WithContext(ctx).Save(&existing).Error
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Category{}, id).Error
}

func (r *repository) GetByID(ctx context.Context, id uint) (*entity.Category, error) {
	var category entity.Category
	query := r.categoryQuery(ctx).Where("categories.id = ?", id)
	if err := query.First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *repository) GetAll(ctx context.Context, searchString string, limit, offset int) ([]entity.Category, int64, error) {
	var categories []entity.Category
	var totalCount int64
	query := r.categoryQuery(ctx)
	if searchString != "" {
		searchPattern := "%" + searchString + "%"
		query = query.Where("LOWER(categories.name) LIKE LOWER(?)", searchPattern)
	}
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	sqlOffset := offset * limit
	if err := query.Order("categories.created_at DESC").Limit(limit).Offset(sqlOffset).Find(&categories).Error; err != nil {
		return nil, 0, err
	}
	return categories, totalCount, nil
}

func (r *repository) categoryQuery(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).
		Model(&entity.Category{}).
		Select("categories.*, COUNT(DISTINCT startups.id) AS startup_count").
		Joins("LEFT JOIN startup_categories ON startup_categories.category_id = categories.id").
		Joins("LEFT JOIN startups ON startups.id = startup_categories.startup_id AND startups.deleted_at IS NULL").
		Group("categories.id")
}
