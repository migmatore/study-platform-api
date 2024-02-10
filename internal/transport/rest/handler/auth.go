package handler

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/migmatore/study-platform-api/internal/apperrors"
	"github.com/migmatore/study-platform-api/internal/core"
	"github.com/migmatore/study-platform-api/pkg/utils"
)

type AuthUseCase interface {
	Signin(ctx context.Context, req core.UserSigninRequest) (core.UserAuthResponse, error)
	Signup(ctx context.Context, req core.UserSignupRequest) (core.UserAuthResponse, error)
	Refresh(ctx context.Context, req core.UserTokenRefreshRequest) (core.UserAuthResponse, error)
}

type AuthHandler struct {
	authUseCase AuthUseCase
}

func NewUserHandler(authUseCase AuthUseCase) *AuthHandler {
	return &AuthHandler{authUseCase: authUseCase}
}

func (h AuthHandler) Signin(c *fiber.Ctx) error {
	ctx := c.UserContext()
	req := core.UserSigninRequest{}

	if err := c.BodyParser(&req); err != nil {
		return utils.FiberError(c, fiber.StatusBadRequest, err)
	}

	if req.Email == "" || req.Password == "" {
		return utils.FiberError(c, fiber.StatusBadRequest, errors.New("the required parameters cannot be empty"))
	}

	resp, err := h.authUseCase.Signin(ctx, req)
	if err != nil {
		if errors.Is(err, apperrors.EntityNotFound) || errors.Is(err, apperrors.IncorrectPassword) {
			return utils.FiberError(c, fiber.StatusUnauthorized, err)
		}

		return utils.FiberError(c, fiber.StatusInternalServerError, err)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h AuthHandler) Signup(c *fiber.Ctx) error {
	ctx := c.UserContext()
	req := core.UserSignupRequest{}

	if err := c.BodyParser(&req); err != nil {
		return utils.FiberError(c, fiber.StatusBadRequest, err)
	}

	if req.Email == "" || req.Password == "" {
		return utils.FiberError(c, fiber.StatusBadRequest, errors.New("the required parameters cannot be empty"))
	}

	resp, err := h.authUseCase.Signup(ctx, req)
	if err != nil {
		if errors.Is(err, apperrors.EntityAlreadyExist) {
			return utils.FiberError(c, fiber.StatusConflict, err)
		}

		return utils.FiberError(c, fiber.StatusInternalServerError, err)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h AuthHandler) Refresh(c *fiber.Ctx) error {
	ctx := c.UserContext()
	req := core.UserTokenRefreshRequest{}

	if err := c.BodyParser(&req); err != nil {
		return utils.FiberError(c, fiber.StatusBadRequest, err)
	}

	if req.RefreshToken == "" {
		return utils.FiberError(c, fiber.StatusBadRequest, errors.New("the required parameters cannot be empty"))
	}

	resp, err := h.authUseCase.Refresh(ctx, req)
	if err != nil {
		if errors.Is(err, apperrors.EntityNotFound) || errors.Is(err, apperrors.InvalidToken) {
			return utils.FiberError(c, fiber.StatusForbidden, err)
		}

		return utils.FiberError(c, fiber.StatusInternalServerError, err)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
