package service

import (
	"context"
	"startup_back/internal/domain"
	"startup_back/internal/dto"
	"startup_back/internal/repository"
)
type userService struct{
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService{
	return &userService{repo:repo}
}
func (s *userService) CreateUser(ctx context.Context, input dto.CreateUserInput) (*domain.User, error) {
	user,err := s.repo.Create(ctx,&domain.User{
		Username:input.Username,
		Email:input.Email,
		PasswordHash:input.Password,
	})
	if err!=nil{
		return nil,err
	}
	
	return user,nil
}

func (s *userService)	GetUserById(ctx context.Context,id uint)(*domain.User, error){
	user,err:=  s.repo.GetById(ctx,id)
 if err!=nil{
		return nil,err
	}
 return user,nil
}
func (s *userService)	UpdateUser(ctx context.Context,id uint,input dto.CreateUserInput) error{
	err:= s.repo.Update(ctx,id,&domain.User{
		Username:input.Username,
		Email:input.Email,
		PasswordHash:input.Password,
	})
	if err!=nil{
		return err
	}
	return  nil
}

func (s * userService	)	GetUserByEmail(ctx context.Context,email string)(*domain.User, error){
	user,err:=  s.repo.GetByEmail(ctx,email)
 if err!=nil{
		return nil,err
	}
 return user,nil
}