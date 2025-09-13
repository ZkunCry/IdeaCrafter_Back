package repository

import (
	"context"
	"startup_back/internal/domain"

	"gorm.io/gorm"
)
type userRepository struct{
	db *gorm.DB
}
func NewUserRepository(db *gorm.DB) UserRepository{
	return &userRepository{db:db}
}
func (r *userRepository) Create(ctx context.Context,user *domain.User) error{
	return r.db.WithContext(ctx).Create(user).Error
}
func (r *userRepository)	Update(ctx context.Context, user *domain.User) error{
		
}
func (r *userRepository) Delete(ctx context.Context, id uint) error{

}
func (r *userRepository) GetById(ctx context.Context,id uint)(*domain.User, error){
	var user domain.User
	err:= r.db.WithContext(ctx).
		Preload("Startups").
		Preload("Favorites.Startup").
		Preload("Memberships.Startup").
		First(&user,id).Error
	if err != nil{
		return nil,err
	}
	return &user,err
}