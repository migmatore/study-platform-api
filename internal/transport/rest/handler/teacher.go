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
	Create(ctx context.Context, metadata core.TokenMetadata, req core.CreateTeacherRequest) (core.TeacherResponse, error)
	Delete(ctx context.Context, metadata core.TokenMetadata, id int) error
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

func (h TeacherHandler) Create(c *fiber.Ctx) error {
	ctx := c.UserContext()
	claims := jwt.ExtractTokenMetadata(c)
	req := core.CreateTeacherRequest{}

	if err := c.BodyParser(&req); err != nil {
		return utils.FiberError(c, fiber.StatusBadRequest, err)
	}

	teacher, err := h.teacherUseCase.Create(ctx, claims, req)
	if err != nil {
		if errors.Is(err, apperrors.AccessDenied) {
			return utils.FiberError(c, fiber.StatusForbidden, err)
		}

		if errors.Is(err, apperrors.EntityAlreadyExist) {
			return utils.FiberError(c, fiber.StatusConflict, err)
		}

		return utils.FiberError(c, fiber.StatusInternalServerError, err)
	}

	return c.Status(fiber.StatusCreated).JSON(teacher)
}

func (h TeacherHandler) Delete(c *fiber.Ctx) error {
	ctx := c.UserContext()
	claims := jwt.ExtractTokenMetadata(c)

	studentId, err := c.ParamsInt("id")
	if err != nil {
		return utils.FiberError(c, fiber.StatusBadRequest, errors.New("the id must be int"))
	}

	if err := h.teacherUseCase.Delete(ctx, claims, studentId); err != nil {
		if errors.Is(err, apperrors.AccessDenied) {
			return utils.FiberError(c, fiber.StatusForbidden, err)
		}

		return utils.FiberError(c, fiber.StatusInternalServerError, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "teacher successfully deleted",
	})
}
