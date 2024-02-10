package usecase

import (
	"context"
	"github.com/migmatore/study-platform-api/internal/apperrors"
	"github.com/migmatore/study-platform-api/internal/core"
)

type LessonService interface {
	All(ctx context.Context, classroomId int) ([]core.Lesson, error)
}

type ClassroomService interface {
	ById(ctx context.Context, id int) (core.Classroom, error)
}

type LessonTeacherService interface {
	ById(ctx context.Context, id int) (core.User, error)
}

type LessonUseCase struct {
	lessonsService   LessonService
	teacherService   LessonTeacherService
	classroomService ClassroomService
}

func NewLessonUseCase(
	lessonsService LessonService,
	classroomService ClassroomService,
	teacherService LessonTeacherService,
) *LessonUseCase {
	return &LessonUseCase{
		lessonsService:   lessonsService,
		classroomService: classroomService,
		teacherService:   teacherService,
	}
}

func (uc LessonUseCase) All(
	ctx context.Context,
	metadata core.TokenMetadata,
	classroomId int,
) ([]core.LessonResponse, error) {
	teacher, err := uc.teacherService.ById(ctx, metadata.UserId)
	if err != nil {
		return nil, err
	}

	classroom, err := uc.classroomService.ById(ctx, classroomId)
	if err != nil {
		return nil, err
	}

	if classroom.TeacherId != teacher.Id {
		return nil, apperrors.EntityNotFound
	}

	lessons, err := uc.lessonsService.All(ctx, classroomId)
	if err != nil {
		return nil, err
	}

	lessonsResp := make([]core.LessonResponse, 0, len(lessons))

	for _, lesson := range lessons {
		lessonsResp = append(lessonsResp, core.LessonResponse{
			Id:          lesson.Id,
			Title:       lesson.Title,
			ClassroomId: lesson.ClassroomId,
			Active:      lesson.Active,
		})
	}

	return lessonsResp, nil
}
