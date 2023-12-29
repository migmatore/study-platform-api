package handler

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/migmatore/study-platform-api/internal/core"
	"github.com/migmatore/study-platform-api/pkg/utils"
)

type AuthUseCase interface {
	Signin(ctx context.Context, req core.UserSigninRequest) (string, error)
	Signup(ctx context.Context, req core.UserSignupRequest) (string, error)
}

type AuthHandler struct {
	authUseCase AuthUseCase
}

func NewUserHandler(authUseCase AuthUseCase) *AuthHandler {
	return &AuthHandler{authUseCase: authUseCase}
}

func (h AuthHandler) Signin(c *fiber.Ctx) error {
	ctx := c.UserContext()
	user := core.UserSigninRequest{}

	if err := c.BodyParser(&user); err != nil {
		return utils.FiberError(c, fiber.StatusBadRequest, err)
	}

	if user.Email == "" || user.Password == "" {
		return utils.FiberError(c, fiber.StatusBadRequest, errors.New("the required parameters cannot be empty"))
	}

	token, err := h.authUseCase.Signin(ctx, user)
	if err != nil {
		return utils.FiberError(c, fiber.StatusInternalServerError, err)
	}

	return c.Status(fiber.StatusOK).JSON(token)
}

func (h AuthHandler) Signup(c *fiber.Ctx) error {
	ctx := c.UserContext()

	user := core.UserSignupRequest{}

	if err := c.BodyParser(&user); err != nil {
		return utils.FiberError(c, fiber.StatusBadRequest, err)
	}

	if user.Email == "" || user.Password == "" {
		return utils.FiberError(c, fiber.StatusBadRequest, errors.New("the required parameters cannot be empty"))
	}

	token, err := h.authUseCase.Signup(ctx, user)
	if err != nil {
		return utils.FiberError(c, fiber.StatusInternalServerError, err)
	}

	return c.Status(fiber.StatusOK).JSON(token)
}
