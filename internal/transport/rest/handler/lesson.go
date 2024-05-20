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

type LessonUseCase interface {
	All(ctx context.Context, metadata core.TokenMetadata, classroomId int) ([]core.LessonResponse, error)
	ById(ctx context.Context, metadata core.TokenMetadata, lessonId int) (core.LessonResponse, error)
	Current(ctx context.Context, metadata core.TokenMetadata, classroomId int) (core.LessonResponse, error)
	Create(ctx context.Context, metadata core.TokenMetadata, classroomId int, req core.CreateLessonRequest) (core.LessonResponse, error)
	Update(ctx context.Context, metadata core.TokenMetadata, req core.UpdateLessonRequest) error
}

type LessonHandler struct {
	lessonUseCase LessonUseCase
}

func NewLessonHandler(lessonUseCase LessonUseCase) *LessonHandler {
	return &LessonHandler{lessonUseCase: lessonUseCase}
}

func (h LessonHandler) ById(c *fiber.Ctx) error {
	ctx := c.UserContext()
	claims := jwt.ExtractTokenMetadata(c)

	lessonId, err := c.ParamsInt("id")
	if err != nil {
		return utils.FiberError(c, fiber.StatusBadRequest, errors.New("the id must be number"))
	}

	lesson, err := h.lessonUseCase.ById(ctx, claims, lessonId)
	if err != nil {
		if errors.Is(err, apperrors.AccessDenied) {
			return utils.FiberError(c, fiber.StatusForbidden, err)
		}

		if errors.Is(err, apperrors.EntityNotFound) {
			return utils.FiberError(c, fiber.StatusNotFound, err)
		}

		return utils.FiberError(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(lesson)
}
