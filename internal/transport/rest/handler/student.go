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

type StudentUseCase interface {
	All(ctx context.Context, metadata core.TokenMetadata) ([]core.StudentResponse, error)
}

type StudentHandler struct {
	studentUseCase StudentUseCase
}

func NewStudentsHandler(studentUseCase StudentUseCase) *StudentHandler {
	return &StudentHandler{studentUseCase: studentUseCase}
}

func (h StudentHandler) Students(c *fiber.Ctx) error {
	ctx := c.UserContext()
	claims := jwt.ExtractTokenMetadata(c)

	students, err := h.studentUseCase.All(ctx, claims)
	if err != nil {
		if errors.Is(err, apperrors.AccessDenied) {
			return utils.FiberError(c, fiber.StatusForbidden, err)
		}

		return utils.FiberError(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(students)
}
