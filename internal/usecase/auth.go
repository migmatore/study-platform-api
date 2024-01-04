package usecase

import (
	"context"
	"github.com/migmatore/study-platform-api/internal/apperrors"
	"github.com/migmatore/study-platform-api/internal/core"
	"github.com/migmatore/study-platform-api/internal/middleware"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	IsExist(ctx context.Context, email string) (bool, error)
	Create(ctx context.Context, user core.User) (core.UserModel, error)
}

type InstitutionService interface {
	IsExist(ctx context.Context, name string) (bool, error)
	Create(ctx context.Context, inst core.Institution) (core.InstitutionModel, error)
}

type AuthUseCase struct {
	transactionService TransactionService
	userService        UserService
	institutionService InstitutionService
}

func NewAuthUseCase(
	transactionService TransactionService,
	userService UserService,
	institutionService InstitutionService,
) *AuthUseCase {
	return &AuthUseCase{
		transactionService: transactionService,
		userService:        userService,
		institutionService: institutionService,
	}
}

func (uc AuthUseCase) Signin(ctx context.Context, req core.UserSigninRequest) (core.UserAuthResponse, error) {
	return core.UserAuthResponse{}, nil
}

func (uc AuthUseCase) Signup(ctx context.Context, req core.UserSignupRequest) (core.UserAuthResponse, error) {
	instExist, err := uc.institutionService.IsExist(ctx, req.InstitutionName)
	if err != nil {
		return core.UserAuthResponse{}, err
	}

	if instExist {
		return core.UserAuthResponse{}, apperrors.EntityAlreadyExist
	}

	userExist, err := uc.userService.IsExist(ctx, req.Email)
	if err != nil {
		return core.UserAuthResponse{}, err
	}

	if userExist {
		return core.UserAuthResponse{}, apperrors.EntityAlreadyExist
	}

	var user core.UserModel

	if req.Role == core.AdminRole {
		err = uc.transactionService.WithinTransaction(ctx, func(txCtx context.Context) error {
			inst, err := uc.institutionService.Create(txCtx, core.Institution{
				Name: req.InstitutionName,
			})
			if err != nil {
				return err
			}

			hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
			if err != nil {
				return err
			}

			user, err = uc.userService.Create(txCtx, core.User{
				FullName:     req.FullName,
				Email:        req.Email,
				PasswordHash: string(hash),
				Role:         core.AdminRole,
				Institution:  &inst,
			})
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return core.UserAuthResponse{}, err
		}
	} else {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return core.UserAuthResponse{}, err
		}

		user, err = uc.userService.Create(ctx, core.User{
			FullName:     req.FullName,
			Email:        req.Email,
			PasswordHash: string(hash),
			Role:         core.TeacherRole,
		})
		if err != nil {
			return core.UserAuthResponse{}, err
		}
	}

	tokenClaims, err := middleware.GenerateNewAccessToken(user.Id, string(req.Role))
	if err != nil {
		return core.UserAuthResponse{}, err
	}

	return core.UserAuthResponse{
		Token: tokenClaims.Token,
		Role:  tokenClaims.Role,
	}, nil
}
