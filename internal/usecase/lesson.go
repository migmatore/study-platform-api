package usecase

import (
	"context"
	"errors"
	"github.com/migmatore/study-platform-api/internal/apperrors"
	"github.com/migmatore/study-platform-api/internal/core"
)

type LessonService interface {
	All(ctx context.Context, classroomId int) ([]core.Lesson, error)
	ById(ctx context.Context, lessonId int) (core.Lesson, error)
	Create(ctx context.Context, lesson core.Lesson) (core.Lesson, error)
	Update(ctx context.Context, lesson core.UpdateLesson) error
	Delete(ctx context.Context, id int) error
	IsBelongs(ctx context.Context, lessonId int, teacherId int) (bool, error)
}

type LessonClassroomService interface {
	ById(ctx context.Context, id int) (core.Classroom, error)
	IsBelongs(ctx context.Context, classroomId int, teacherId int) (bool, error)
	IsIn(ctx context.Context, classroomId, studentId int) (bool, error)
	Students(ctx context.Context, classroomId int) ([]core.Student, error)
}

type LessonTeacherService interface {
	ById(ctx context.Context, id int) (core.User, error)
}

type LessonUseCase struct {
	lessonsService   LessonService
	teacherService   LessonTeacherService
	classroomService LessonClassroomService
}

func NewLessonUseCase(
	lessonsService LessonService,
	classroomService LessonClassroomService,
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
	belongs, err := uc.classroomService.IsBelongs(ctx, classroomId, metadata.UserId)
	if err != nil {
		return nil, err
	}

	if !belongs {
		return nil, apperrors.AccessDenied
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
			Content:     lesson.Content,
			Active:      lesson.Active,
		})
	}

	return lessonsResp, nil
}

func (uc LessonUseCase) ById(
	ctx context.Context,
	metadata core.TokenMetadata,
	lessonId int,
) (core.LessonResponse, error) {
	if core.RoleType(metadata.Role) == core.StudentRole {
		return core.LessonResponse{}, apperrors.AccessDenied
	}

	lesson, err := uc.lessonsService.ById(ctx, lessonId)
	if err != nil {
		return core.LessonResponse{}, err
	}

	// TODO: implement admin access check

	belongs, err := uc.classroomService.IsBelongs(ctx, lesson.ClassroomId, metadata.UserId)
	if err != nil {
		return core.LessonResponse{}, err
	}

	if !belongs {
		return core.LessonResponse{}, apperrors.AccessDenied
	}

	return core.LessonResponse{
		Id:          lesson.Id,
		Title:       lesson.Title,
		ClassroomId: lesson.ClassroomId,
		Content:     lesson.Content,
		Active:      lesson.Active,
	}, nil
}

func (uc LessonUseCase) Current(
	ctx context.Context,
	metadata core.TokenMetadata,
	classroomId int,
) (core.LessonResponse, error) {
	in, err := uc.classroomService.IsIn(ctx, classroomId, metadata.UserId)
	if err != nil {
		return core.LessonResponse{}, err
	}

	if !in {
		return core.LessonResponse{}, apperrors.AccessDenied
	}

	lessons, err := uc.lessonsService.All(ctx, classroomId)
	if err != nil {
		return core.LessonResponse{}, err
	}

	for _, lesson := range lessons {
		if !lesson.Active {
			continue
		}

		return core.LessonResponse{
			Id:          lesson.Id,
			Title:       lesson.Title,
			ClassroomId: lesson.ClassroomId,
			Content:     lesson.Content,
			Active:      lesson.Active,
		}, nil
	}

	return core.LessonResponse{}, apperrors.EntityNotFound
}

func (uc LessonUseCase) Create(
	ctx context.Context,
	metadata core.TokenMetadata,
	classroomId int,
	req core.CreateLessonRequest,
) (core.LessonResponse, error) {
	if core.RoleType(metadata.Role) != core.TeacherRole {
		return core.LessonResponse{}, apperrors.AccessDenied
	}

	belongs, err := uc.classroomService.IsBelongs(ctx, classroomId, metadata.UserId)
	if err != nil {
		return core.LessonResponse{}, err
	}

	if !belongs {
		return core.LessonResponse{}, apperrors.AccessDenied
	}

	lessons, err := uc.lessonsService.All(ctx, classroomId)
	if err != nil {
		return core.LessonResponse{}, err
	}

	if req.Active {
		for _, lesson := range lessons {
			if !lesson.Active {
				continue
			}

			inactive := false

			if err := uc.lessonsService.Update(ctx, core.UpdateLesson{
				Id:     lesson.Id,
				Active: &inactive,
			}); err != nil {
				return core.LessonResponse{}, err
			}
		}
	}

	newLesson, err := uc.lessonsService.Create(ctx, core.Lesson{
		Title:       req.Title,
		ClassroomId: classroomId,
		Active:      req.Active,
	})
	if err != nil {
		return core.LessonResponse{}, err
	}

	return core.LessonResponse{
		Id:          newLesson.Id,
		Title:       newLesson.Title,
		ClassroomId: newLesson.ClassroomId,
		Content:     newLesson.Content,
		Active:      newLesson.Active,
	}, nil
}

func (uc LessonUseCase) Update(
	ctx context.Context,
	metadata core.TokenMetadata,
	req core.UpdateLessonRequest,
) error {
	if core.RoleType(metadata.Role) != core.TeacherRole {
		return apperrors.AccessDenied
	}

	if req.ClassroomId == nil || req.LessonId == nil {
		return errors.New("classroomId or lessonId must be number")
	}

	belongs, err := uc.classroomService.IsBelongs(ctx, *req.ClassroomId, metadata.UserId)
	if err != nil {
		return err
	}

	if !belongs {
		return apperrors.AccessDenied
	}

	lessons, err := uc.lessonsService.All(ctx, *req.ClassroomId)
	if err != nil {
		return err
	}

	if req.Active != nil && *req.Active {
		for _, lesson := range lessons {
			if !lesson.Active {
				continue
			}

			inactive := false

			if err := uc.lessonsService.Update(ctx, core.UpdateLesson{
				Id:     lesson.Id,
				Active: &inactive,
			}); err != nil {
				return err
			}
		}
	}

	if err := uc.lessonsService.Update(ctx, core.UpdateLesson{
		Id:          *req.LessonId,
		Title:       req.Title,
		ClassroomId: req.ClassroomId,
		Content:     req.Content,
		Active:      req.Active,
	}); err != nil {
		return err
	}

	return nil
}

func (uc LessonUseCase) Delete(ctx context.Context, metadata core.TokenMetadata, lessonId int) error {
	switch core.RoleType(metadata.Role) {
	case core.AdminRole:
		return apperrors.AccessDenied
	case core.TeacherRole:
		isBelongs, err := uc.lessonsService.IsBelongs(ctx, lessonId, metadata.UserId)
		if err != nil {
			return err
		}

		if !isBelongs {
			return apperrors.AccessDenied
		}

		if err := uc.lessonsService.Delete(ctx, lessonId); err != nil {
			return err
		}
	case core.StudentRole:
		return apperrors.AccessDenied
	}

	return nil
}
