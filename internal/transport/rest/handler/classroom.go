package handler

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/migmatore/study-platform-api/internal/core"
	"github.com/migmatore/study-platform-api/pkg/jwt"
	"github.com/migmatore/study-platform-api/pkg/utils"
)

type ClassroomUseCase interface {
	All(ctx context.Context, metadata core.TokenMetadata) ([]core.ClassroomResponse, error)
}

type LessonUseCase interface {
	All(ctx context.Context, metadata core.TokenMetadata, classroomId int) ([]core.LessonResponse, error)
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
