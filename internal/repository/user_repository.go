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
func (r *userRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
    if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
        return nil, err
    }
    return user, nil
}
func (r *userRepository)	Update(ctx context.Context,id uint, user *domain.User) error{
		return r.db.WithContext(ctx).
        Model(&domain.User{}).
        Where("id = ?", id).
        Updates(user).Error
}
func (r *userRepository) Delete(ctx context.Context, id uint) error{

	return  r.db.WithContext(ctx).Delete(&domain.User{},id).Error

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

func (r * userRepository) GetByEmail(ctx context.Context, email string)(*domain.User, error){
	var user domain.User
	err:= r.db.WithContext(ctx).First(&user,"email = ?",email).Error
	if err != nil{
		return nil,err
	}
	return &user,err
}