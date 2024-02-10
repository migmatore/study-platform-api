package service

import (
	"context"
	"github.com/migmatore/study-platform-api/internal/core"
)

type LessonRepo interface {
	All(ctx context.Context, classroomId int) ([]core.LessonModel, error)
}

type LessonService struct {
	lessonRepo LessonRepo
}

func NewLessonService(lessonRepo LessonRepo) *LessonService {
	return &LessonService{lessonRepo: lessonRepo}
}

func (s LessonService) All(ctx context.Context, classroomId int) ([]core.Lesson, error) {
	lessonsModel, err := s.lessonRepo.All(ctx, classroomId)
	if err != nil {
		return nil, err
	}

	lessons := make([]core.Lesson, 0, len(lessonsModel))

	for _, model := range lessonsModel {
		lessons = append(lessons, core.Lesson{
			Id:          model.Id,
			Title:       model.Title,
			ClassroomId: model.ClassroomId,
			Active:      model.Active,
		})
	}

	return lessons, nil
}
