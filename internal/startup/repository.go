package startup

import (
	"context"
	"fmt"
	"startup_back/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, startup *entity.Startup, categoryIDs []uint) (*entity.Startup, error)
	GetByID(ctx context.Context, id uint) (*entity.Startup, error)
	GetAll(ctx context.Context, searchString, categorySlug string, limit, offset int) ([]*entity.Startup, int, error)
	Update(ctx context.Context, id uint, input UpdateInput) (*entity.Startup, error)
	Delete(ctx context.Context, id uint) error
	GetUserStartups(ctx context.Context, userID uint) ([]entity.Startup, error)
	AddCategories(ctx context.Context, startupID uint, categoryIDs []uint) (*entity.Startup, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, startup *entity.Startup, categoryIDs []uint) (*entity.Startup, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(startup).Error; err != nil {
			return err
		}
		if len(categoryIDs) > 0 {
			var categories []entity.Category
			if err := tx.Where("id IN ?", categoryIDs).Find(&categories).Error; err != nil {
				return err
			}
			if err := tx.Model(startup).Association("Categories").Replace(categories); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Preload("Categories").Preload("Creator").Preload("Stage").Preload("Files").First(startup, startup.ID).Error; err != nil {
		return nil, err
	}
	return startup, nil
}

func (r *repository) GetByID(ctx context.Context, id uint) (*entity.Startup, error) {
	var startup entity.Startup
	if err := r.db.WithContext(ctx).Where("id = ?", id).Preload("Categories").Preload("Creator").Preload("Stage").Preload("Files").Preload("Vacancies").Preload("Vacancies.Role").Preload("Vacancies.User").First(&startup).Error; err != nil {
		return nil, err
	}
	return &startup, nil
}

func (r *repository) GetAll(ctx context.Context, searchString, categorySlug string, limit, offset int) ([]*entity.Startup, int, error) {
	var startups []*entity.Startup
	var totalCount int64
	query := r.db.WithContext(ctx).Model(&entity.Startup{})
	if searchString != "" {
		searchPattern := "%" + searchString + "%"
		query = query.Where("LOWER(startups.name) LIKE LOWER(?) OR LOWER(startups.description) LIKE LOWER(?)", searchPattern, searchPattern)
	}
	if categorySlug != "" {
		query = query.
			Joins("JOIN startup_categories ON startup_categories.startup_id = startups.id").
			Joins("JOIN categories ON categories.id = startup_categories.category_id").
			Where("categories.slug = ?", categorySlug)
	}
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	sqlOffset := offset * limit
	if err := query.Preload("Creator").Preload("Categories").Preload("Vacancies").Preload("Vacancies.Role").Preload("Vacancies.User").Preload("Files").Preload("Stage").Order("created_at DESC").Limit(limit).Offset(sqlOffset).Find(&startups).Error; err != nil {
		return nil, 0, err
	}
	return startups, int(totalCount), nil
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Startup{}, id).Error
}

func (r *repository) GetUserStartups(ctx context.Context, userID uint) ([]entity.Startup, error) {
	var startups []entity.Startup
	if err := r.db.WithContext(ctx).Preload("Categories").Preload("Creator").Preload("Stage").Preload("Files").Preload("Vacancies").Preload("Vacancies.Role").Preload("Vacancies.User").Where("creator_id = ?", userID).Order("created_at DESC").Find(&startups).Error; err != nil {
		return nil, err
	}
	return startups, nil
}

func (r *repository) AddCategories(ctx context.Context, startupID uint, categoryIDs []uint) (*entity.Startup, error) {
	var startup entity.Startup
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&startup, startupID).Error; err != nil {
			return err
		}
		if len(categoryIDs) == 0 {
			return nil
		}
		var categories []entity.Category
		if err := tx.Where("id IN ?", categoryIDs).Find(&categories).Error; err != nil {
			return err
		}
		if len(categories) == 0 {
			return fmt.Errorf("categories not found")
		}
		if err := tx.Model(&startup).Association("Categories").Append(categories); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Preload("Categories").Preload("Creator").Preload("Stage").Preload("Files").Preload("Vacancies").Preload("Vacancies.Role").Preload("Vacancies.User").First(&startup, startupID).Error; err != nil {
		return nil, err
	}
	return &startup, nil
}

func (r *repository) Update(ctx context.Context, id uint, input UpdateInput) (*entity.Startup, error) {
	var startup entity.Startup
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&startup, id).Error; err != nil {
			return err
		}

		startup.Name = input.Name
		startup.ShortDescription = input.ShortDescription
		startup.Description = input.Description
		startup.TargetAudience = input.TargetAudience
		startup.Problem = input.Problem
		startup.Solution = input.Solution
		if input.StageID != 0 {
			startup.StageID = input.StageID
		}
	
		if input.LogoFile != "" {
			startup.LogoURL = input.LogoFile
		}

		if err := tx.Save(&startup).Error; err != nil {
			return err
		}

		if len(input.CategoryIDs) > 0 {
			var categories []entity.Category
			if err := tx.Where("id IN ?", input.CategoryIDs).Find(&categories).Error; err != nil {
				return err
			}
			if len(categories) == 0 {
				return fmt.Errorf("categories not found")
			}
			if err := tx.Model(&startup).Association("Categories").Replace(categories); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}
