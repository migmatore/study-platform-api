package handler

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/migmatore/study-platform-api/internal/apperrors"
	"github.com/migmatore/study-platform-api/internal/core"
	"github.com/migmatore/study-platform-api/pkg/jwt"
	"github.com/migmatore/study-platform-api/pkg/utils"
)

type UserUseCase interface {
	Profile(ctx context.Context, metadata core.TokenMetadata) (core.ProfileResponse, error)
	UpdateProfile(ctx context.Context, metadata core.TokenMetadata, req core.UpdateProfileRequest) (core.ProfileResponse, error)
}

type UserHandler struct {
	userUseCase UserUseCase
}

func NewUserHandler(userUseCase UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: userUseCase}
}

func (h UserHandler) Profile(c *fiber.Ctx) error {
	ctx := c.UserContext()
	claims := jwt.ExtractTokenMetadata(c)

	profile, err := h.userUseCase.Profile(ctx, claims)
	if err != nil {
		if errors.Is(err, apperrors.EntityNotFound) {
			return utils.FiberError(c, fiber.StatusNotFound, err)
		}

		return utils.FiberError(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(profile)
}

func (h UserHandler) UpdateProfile(c *fiber.Ctx) error {
	ctx := c.UserContext()
	claims := jwt.ExtractTokenMetadata(c)
	req := core.UpdateProfileRequest{}

	if err := c.BodyParser(&req); err != nil {
		return utils.FiberError(c, fiber.StatusBadRequest, err)
	}

	newProfile, err := h.userUseCase.UpdateProfile(ctx, claims, req)
	if err != nil {
		if errors.Is(err, apperrors.EntityNotFound) {
			return utils.FiberError(c, fiber.StatusNotFound, err)
		}

		if errors.Is(err, apperrors.EntityAlreadyExist) {
			return utils.FiberError(c, fiber.StatusConflict, err)
		}

		return utils.FiberError(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(newProfile)
}
