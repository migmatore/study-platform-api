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
	Create(ctx context.Context, metadata core.TokenMetadata, req core.CreateStudentRequest) (core.StudentResponse, error)
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
func (h StudentHandler) Create(c *fiber.Ctx) error {
	ctx := c.UserContext()
	claims := jwt.ExtractTokenMetadata(c)
	req := core.CreateStudentRequest{}

	if err := c.BodyParser(&req); err != nil {
		return utils.FiberError(c, fiber.StatusBadRequest, err)
	}

	student, err := h.studentUseCase.Create(ctx, claims, req)
	if err != nil {
		if errors.Is(err, apperrors.AccessDenied) {
			return utils.FiberError(c, fiber.StatusForbidden, err)
		}

		if errors.Is(err, apperrors.EntityAlreadyExist) {
			return utils.FiberError(c, fiber.StatusConflict, err)
		}

		if errors.Is(err, apperrors.NumberOfStudentsExceeded) {
			return utils.FiberError(c, fiber.StatusBadRequest, err)
		}

		return utils.FiberError(c, fiber.StatusInternalServerError, err)
	}

	return c.Status(fiber.StatusCreated).JSON(student)
}
