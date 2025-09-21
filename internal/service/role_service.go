package service

import (
	"context"
	"startup_back/internal/domain"
	"startup_back/internal/dto"
	"startup_back/internal/repository"
)
type RoleService interface {
	GetAll(ctx context.Context) ([]*domain.Role, error)
	GetByID(ctx context.Context, id uint) (*domain.Role, error)
	Create(ctx context.Context, input *dto.CreateRoleInput) (*domain.Role, error)
	Update(ctx context.Context, id uint, role *domain.Role) error
	Delete(ctx context.Context, id uint) error
}

type roleService struct {
	repo repository.RoleRepository
}
func NewRoleService(repo repository.RoleRepository) RoleService {
	return &roleService{repo: repo}
}
func (r *roleService) GetAll(ctx context.Context) ([]*domain.Role, error) {
	roles, err:= r.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return roles,nil
}

func (r *roleService) GetByID(ctx context.Context, id uint) (*domain.Role, error) {
	role, err := r.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (r *roleService) Create(ctx context.Context, input *dto.CreateRoleInput) (*domain.Role, error) {
	roleInput := &domain.Role{Name: input.Name}
	role, err := r.repo.Create(ctx, roleInput)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (r *roleService) Update(ctx context.Context, id uint, role *domain.Role) error {
	if err := r.repo.Update(ctx, id, role); err != nil {
		return err
	}
	return nil
}

func (r *roleService) Delete(ctx context.Context, id uint) error {
	if err := r.repo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}