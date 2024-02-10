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
}

type InstitutionService interface {
	IsExist(ctx context.Context, name string) (bool, error)
	Create(ctx context.Context, inst core.Institution) (core.Institution, error)
}

type TokenService interface {
	Token(userId int, role string) (core.TokenWithClaims, error)
	RefreshToken(userId int) (core.RefreshTokenWithClaims, error)
	ExtractTokenMetadata(tokenString string) (core.TokenMetadata, error)
	ExtractRefreshTokenMetadata(tokenString string) (core.RefreshTokenMetadata, error)
}

type AuthUseCase struct {
	transactionService TransactionService
	userService        UserService
	institutionService InstitutionService
	tokenService       TokenService
}

func NewAuthUseCase(
	transactionService TransactionService,
	userService UserService,
	institutionService InstitutionService,
	tokenService TokenService,
) *AuthUseCase {
	return &AuthUseCase{
		transactionService: transactionService,
		userService:        userService,
		institutionService: institutionService,
		tokenService:       tokenService,
	}
}

func (uc AuthUseCase) Signin(ctx context.Context, req core.UserSigninRequest) (core.UserAuthResponse, error) {
	userExist, err := uc.userService.IsExist(ctx, req.Email)
	if err != nil {
		return core.UserAuthResponse{}, err
	}

	if !userExist {
		return core.UserAuthResponse{}, apperrors.EntityNotFound
	}

	user, err := uc.userService.ByEmail(ctx, req.Email)
	if err != nil {
		return core.UserAuthResponse{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return core.UserAuthResponse{}, apperrors.IncorrectPassword
	}

	tokenClaims, err := uc.tokenService.Token(user.Id, string(user.Role))
	if err != nil {
		return core.UserAuthResponse{}, err
	}

	refreshTokenClaims, err := uc.tokenService.RefreshToken(user.Id)
	if err != nil {
		if err != nil {
			return core.UserAuthResponse{}, err
		}
	}

	return core.UserAuthResponse{
		Token:        tokenClaims.Token,
		RefreshToken: refreshTokenClaims.Token,
		Role:         tokenClaims.Role,
	}, nil
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

	var user core.User

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

	tokenClaims, err := uc.tokenService.Token(user.Id, string(req.Role))
	if err != nil {
		return core.UserAuthResponse{}, err
	}

	refreshTokenClaims, err := uc.tokenService.RefreshToken(user.Id)
	if err != nil {
		if err != nil {
			return core.UserAuthResponse{}, err
		}
	}

	return core.UserAuthResponse{
		Token:        tokenClaims.Token,
		RefreshToken: refreshTokenClaims.Token,
		Role:         tokenClaims.Role,
	}, nil
}

func (uc AuthUseCase) Refresh(ctx context.Context, req core.UserTokenRefreshRequest) (core.UserAuthResponse, error) {
	metadata, err := uc.tokenService.ExtractRefreshTokenMetadata(req.RefreshToken)
	if err != nil {
		return core.UserAuthResponse{}, err
	}

	exist, err := uc.userService.IsExistById(ctx, metadata.UserId)
	if err != nil {
		return core.UserAuthResponse{}, err
	}

	if !exist {
		return core.UserAuthResponse{}, apperrors.EntityNotFound
	}

	user, err := uc.userService.ById(ctx, metadata.UserId)
	if err != nil {
		return core.UserAuthResponse{}, err
	}

	tokenClaims, err := uc.tokenService.Token(user.Id, string(user.Role))
	if err != nil {
		return core.UserAuthResponse{}, err
	}

	refreshTokenClaims, err := uc.tokenService.RefreshToken(user.Id)
	if err != nil {
		if err != nil {
			return core.UserAuthResponse{}, err
		}
	}

	return core.UserAuthResponse{
		Token:        tokenClaims.Token,
		RefreshToken: refreshTokenClaims.Token,
		Role:         tokenClaims.Role,
	}, nil
}
