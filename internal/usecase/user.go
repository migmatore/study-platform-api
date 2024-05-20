package usecase

import (
	"context"
	"github.com/migmatore/study-platform-api/internal/apperrors"
	"github.com/migmatore/study-platform-api/internal/core"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	IsExist(ctx context.Context, email string) (bool, error)
	IsExistById(ctx context.Context, id int) (bool, error)
	ByEmail(ctx context.Context, email string) (core.User, error)
	ById(ctx context.Context, id int) (core.User, error)
	Create(ctx context.Context, user core.User) (core.User, error)
	UpdateProfile(ctx context.Context, userId int, profile core.UpdateUserProfile) (core.UserProfile, error)
}

type UserUseCase struct {
	userService UserService
}

func NewUserUseCase(userService UserService) *UserUseCase {
	return &UserUseCase{userService: userService}
}

func (uc UserUseCase) Profile(ctx context.Context, metadata core.TokenMetadata) (core.ProfileResponse, error) {
	user, err := uc.userService.ById(ctx, metadata.UserId)
	if err != nil {
		return core.ProfileResponse{}, err
	}

	return core.ProfileResponse{
		FullName: user.FullName,
		Phone:    user.Phone,
		Email:    user.Email,
	}, nil
}

// TODO CRITICAL SECURITY ISSUE
func (uc UserUseCase) UpdateProfile(
	ctx context.Context,
	metadata core.TokenMetadata,
	req core.UpdateProfileRequest,
) (core.ProfileResponse, error) {
	var newProfile core.UpdateUserProfile

	if req.Email != nil && *req.Email != "" {
		isExist, err := uc.userService.IsExist(ctx, *req.Email)
		if err != nil {
			return core.ProfileResponse{}, err
		}

		if isExist {
			return core.ProfileResponse{}, apperrors.EntityAlreadyExist
		}

		newProfile.Email = req.Email
	}

	if req.Password != nil && *req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return core.ProfileResponse{}, err
		}

		hashStr := string(hash)

		newProfile.Password = &hashStr
	}

	if req.FullName != nil && *req.FullName != "" {
		newProfile.FullName = req.FullName
	}

	if req.Phone != nil && *req.Phone != "" {
		newProfile.Phone = req.Phone
	}

	updatedProfile, err := uc.userService.UpdateProfile(ctx, metadata.UserId, newProfile)
	if err != nil {
		return core.ProfileResponse{}, err
	}

	return core.ProfileResponse{
		FullName: updatedProfile.FullName,
		Phone:    updatedProfile.Phone,
		Email:    updatedProfile.Email,
	}, nil
}
