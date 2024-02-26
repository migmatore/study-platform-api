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

type ClassroomUseCase interface {
	All(ctx context.Context, metadata core.TokenMetadata) ([]core.ClassroomResponse, error)
}

type LessonUseCase interface {
	All(ctx context.Context, metadata core.TokenMetadata, classroomId int) ([]core.LessonResponse, error)
	Create(ctx context.Context, metadata core.TokenMetadata, classroomId int, req core.CreateLessonRequest) (core.LessonResponse, error)
	Update(ctx context.Context, metadata core.TokenMetadata, req core.UpdateLessonRequest) error
}

type ClassroomHandler struct {
	classroomUseCase ClassroomUseCase
	lessonUseCase    LessonUseCase
}

func NewClassroomHandler(classroomUseCase ClassroomUseCase, lessonUseCase LessonUseCase) *ClassroomHandler {
	return &ClassroomHandler{classroomUseCase: classroomUseCase, lessonUseCase: lessonUseCase}
}

func (h ClassroomHandler) All(c *fiber.Ctx) error {
	ctx := c.UserContext()
	claims := jwt.ExtractTokenMetadata(c)

	classrooms, err := h.classroomUseCase.All(ctx, claims)
	if err != nil {
		return utils.FiberError(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(classrooms)
}

func (h ClassroomHandler) Lessons(c *fiber.Ctx) error {
	ctx := c.UserContext()
	claims := jwt.ExtractTokenMetadata(c)

	classroomId, err := c.ParamsInt("id")
	if err != nil {
		return utils.FiberError(c, fiber.StatusBadRequest, errors.New("the id must be number"))
	}

	lessons, err := h.lessonUseCase.All(ctx, claims, classroomId)
	if err != nil {
		return utils.FiberError(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(lessons)
}

func (h ClassroomHandler) CreateLesson(c *fiber.Ctx) error {
	ctx := c.UserContext()
	claims := jwt.ExtractTokenMetadata(c)

	classroomId, err := c.ParamsInt("id")
	if err != nil {
		return utils.FiberError(c, fiber.StatusBadRequest, errors.New("the id must be number"))
	}

	req := core.CreateLessonRequest{}

	if err := c.BodyParser(&req); err != nil {
		return utils.FiberError(c, fiber.StatusBadRequest, err)
	}

	newLesson, err := h.lessonUseCase.Create(ctx, claims, classroomId, req)
	if err != nil {
		if errors.Is(err, apperrors.AccessDenied) {
			return utils.FiberError(c, fiber.StatusForbidden, err)
		}

		return utils.FiberError(c, fiber.StatusInternalServerError, err)
	}

	return c.Status(fiber.StatusCreated).JSON(newLesson)
}

func (h ClassroomHandler) UpdateLesson(c *fiber.Ctx) error {
	ctx := c.UserContext()
	claims := jwt.ExtractTokenMetadata(c)

	classroomId, err := c.ParamsInt("classroomId")
	if err != nil {
		return utils.FiberError(c, fiber.StatusBadRequest, errors.New("the id must be number"))
	}

	req := core.UpdateLessonRequest{
		ClassroomId: &classroomId,
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.FiberError(c, fiber.StatusBadRequest, err)
	}

	if req.LessonId == nil {
		return utils.FiberError(c, fiber.StatusBadRequest, errors.New("the required parameters cannot be empty"))
	}

	if err := h.lessonUseCase.Update(ctx, claims, req); err != nil {
		if errors.Is(err, apperrors.AccessDenied) {
			return utils.FiberError(c, fiber.StatusForbidden, err)
		}

		return utils.FiberError(c, fiber.StatusInternalServerError, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successful update",
	})
}
