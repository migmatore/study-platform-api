package service

import (
	"context"
	"github.com/migmatore/study-platform-api/internal/core"
)

type UserRepo interface {
	IsExist(ctx context.Context, email string) (bool, error)
	IsExistById(ctx context.Context, id int) (bool, error)
	Create(ctx context.Context, user core.UserModel) (core.UserModel, error)
	ByEmail(ctx context.Context, email string) (core.UserModel, error)
	ById(ctx context.Context, id int) (core.UserModel, error)
}

type UserRoleRepo interface {
	ByName(ctx context.Context, name string) (core.RoleModel, error)
	ById(ctx context.Context, id int) (core.RoleModel, error)
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

func (s UserService) IsExistById(ctx context.Context, id int) (bool, error) {
	return s.userRepo.IsExistById(ctx, id)
}

func (s UserService) Create(ctx context.Context, user core.User) (core.User, error) {
	role, err := s.roleRepo.ByName(ctx, string(user.Role))
	if err != nil {
		return core.User{}, err
	}

	var userModel core.UserModel

	if user.Institution == nil {
		userModel, err = s.userRepo.Create(ctx, core.UserModel{
			FullName:      user.FullName,
			Phone:         user.Phone,
			Email:         user.Email,
			PasswordHash:  user.PasswordHash,
			RoleId:        role.Id,
			InstitutionId: nil,
		})
	} else {
		userModel, err = s.userRepo.Create(ctx, core.UserModel{
			FullName:      user.FullName,
			Phone:         user.Phone,
			Email:         user.Email,
			PasswordHash:  user.PasswordHash,
			RoleId:        role.Id,
			InstitutionId: &user.Institution.Id,
		})
	}

	if err != nil {
		return core.User{}, err
	}

	return core.User{
		Id:           userModel.Id,
		FullName:     userModel.FullName,
		Phone:        userModel.Phone,
		Email:        userModel.Email,
		PasswordHash: userModel.PasswordHash,
		Role:         core.RoleType(role.Name),
		Institution:  nil,
	}, nil
}

func (s UserService) ByEmail(ctx context.Context, email string) (core.User, error) {
	userModel, err := s.userRepo.ByEmail(ctx, email)
	if err != nil {
		return core.User{}, err
	}

	role, err := s.roleRepo.ById(ctx, userModel.RoleId)
	if err != nil {
		return core.User{}, err
	}

	return core.User{
		Id:           userModel.Id,
		FullName:     userModel.FullName,
		Phone:        userModel.Phone,
		Email:        userModel.Email,
		PasswordHash: userModel.PasswordHash,
		Role:         core.RoleType(role.Name),
		Institution:  nil,
	}, nil
}

func (s UserService) ById(ctx context.Context, id int) (core.User, error) {
	userModel, err := s.userRepo.ById(ctx, id)
	if err != nil {
		return core.User{}, err
	}

	role, err := s.roleRepo.ById(ctx, userModel.RoleId)
	if err != nil {
		return core.User{}, err
	}

	return core.User{
		Id:           userModel.Id,
		FullName:     userModel.FullName,
		Phone:        userModel.Phone,
		Email:        userModel.Email,
		PasswordHash: userModel.PasswordHash,
		Role:         core.RoleType(role.Name),
		Institution:  nil,
	}, nil
}
