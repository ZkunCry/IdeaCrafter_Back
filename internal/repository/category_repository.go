package repository

import (
	"context"
	"fmt"
	"startup_back/internal/domain"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(ctx context.Context, categoryFields *domain.Category) (*domain.Category, error)
	Update(ctx context.Context, id uint, categoryFields *domain.Category) error
	Delete(ctx context.Context, id uint) error
	GetById(ctx context.Context, id uint) (*domain.Category, error)	
	GetAll(ctx context.Context, searchString string, limit, offset int) ([]domain.Category, int64, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, categoryFields *domain.Category) (*domain.Category, error) {
	if err := r.db.WithContext(ctx).Create(categoryFields).Error; err != nil {
		return nil, err
	}
	return categoryFields, nil
}

func (r *categoryRepository) Update(ctx context.Context, id uint, categoryFields *domain.Category) error {
	var existing domain.Category
	if err := r.db.WithContext(ctx).First(&existing, id).Error; err != nil {
		return err
	}
	existing.Name = categoryFields.Name
	if err := r.db.WithContext(ctx).Save(&existing).Error; err != nil {
		return err
	}
	return nil
}

func (r *categoryRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&domain.Category{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *categoryRepository) GetById(ctx context.Context, id uint) (*domain.Category, error) {
	var category domain.Category
	if err := r.db.WithContext(ctx).First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) GetAll(ctx context.Context, searchString string, limit, offset int) ([]domain.Category, int64, error) {
	var categories []domain.Category
	var totalCount int64;
	query := r.db.WithContext(ctx).Model(&domain.Category{})
	if searchString != ""{
		searchPattern := "%" + searchString + "%"
		query = query.Where("LOWER(name) LIKE LOWER(?) OR LOWER(description) LIKE LOWER(?)", searchPattern, searchPattern)

	}
	if err := query.Count(&totalCount).Error; err !=nil{
		return nil,0,err	
	}
	sqlOffset := offset * limit
	fmt.Println("Count ",totalCount)
	if err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(sqlOffset).
		Find(&categories).Error; err != nil {
		return nil, 0, err
	}
	return categories, totalCount, nil
}