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

type TeacherUseCase interface {
	All(ctx context.Context, metadata core.TokenMetadata) ([]core.TeacherResponse, error)
}

type TeacherHandler struct {
	teacherUseCase TeacherUseCase
}

func NewTeacherHandler(teacherUseCase TeacherUseCase) *TeacherHandler {
	return &TeacherHandler{teacherUseCase: teacherUseCase}
}

func (h TeacherHandler) All(c *fiber.Ctx) error {
	ctx := c.UserContext()
	claims := jwt.ExtractTokenMetadata(c)

	teachers, err := h.teacherUseCase.All(ctx, claims)
	if err != nil {
		if errors.Is(err, apperrors.AccessDenied) {
			return utils.FiberError(c, fiber.StatusForbidden, err)
		}

		return utils.FiberError(c, fiber.StatusInternalServerError, err)
	}

	return c.Status(fiber.StatusOK).JSON(teachers)
}
