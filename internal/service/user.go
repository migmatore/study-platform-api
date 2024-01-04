package service

import (
	"context"
	"github.com/migmatore/study-platform-api/internal/core"
)

type UserRepo interface {
	IsExist(ctx context.Context, email string) (bool, error)
	Create(ctx context.Context, user core.UserModel) (core.UserModel, error)
}

type UserRoleRepo interface {
	GetByName(ctx context.Context, name string) (core.RoleModel, error)
}

type UserService struct {
	userRepo UserRepo
	roleRepo UserRoleRepo
}

func NewUserService(userRepo UserRepo, roleRepo UserRoleRepo) *UserService {
	return &UserService{userRepo: userRepo, roleRepo: roleRepo}
}

func (s UserService) IsExist(ctx context.Context, email string) (bool, error) {
	return s.userRepo.IsExist(ctx, email)
}

func (s UserService) Create(ctx context.Context, user core.User) (core.UserModel, error) {
	role, err := s.roleRepo.GetByName(ctx, string(user.Role))
	if err != nil {
		return core.UserModel{}, err
	}

	if user.Institution == nil {
		return s.userRepo.Create(ctx, core.UserModel{
			FullName:      user.FullName,
			Phone:         user.Phone,
			Email:         user.Email,
			PasswordHash:  user.PasswordHash,
			RoleId:        role.Id,
			InstitutionId: nil,
		})
	}

	return s.userRepo.Create(ctx, core.UserModel{
		FullName:      user.FullName,
		Phone:         user.Phone,
		Email:         user.Email,
		PasswordHash:  user.PasswordHash,
		RoleId:        role.Id,
		InstitutionId: &user.Institution.Id,
	})
}
