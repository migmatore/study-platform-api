package usecase

import (
	"context"
	"github.com/migmatore/study-platform-api/internal/core"
)

type UserService interface {
	IsExist(ctx context.Context, email string) (bool, error)
	IsExistById(ctx context.Context, id int) (bool, error)
	ByEmail(ctx context.Context, email string) (core.User, error)
	ById(ctx context.Context, id int) (core.User, error)
	Create(ctx context.Context, user core.User) (core.User, error)
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
