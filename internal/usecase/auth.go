package usecase

import (
	"context"
	"github.com/migmatore/study-platform-api/internal/core"
)

type UserService interface {
}

type AuthUseCase struct {
	userService UserService
}

func NewAuthUseCase(userService UserService) *AuthUseCase {
	return &AuthUseCase{userService: userService}
}

func (uc AuthUseCase) Signin(ctx context.Context, req core.UserSigninRequest) (string, error) {

}

func (uc AuthUseCase) Signup(ctx context.Context, req core.UserSignupRequest) (string, error) {

}
