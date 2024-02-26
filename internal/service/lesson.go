package service

import (
	"context"
	"github.com/migmatore/study-platform-api/internal/core"
)

type LessonRepo interface {
	All(ctx context.Context, classroomId int) ([]core.LessonModel, error)
	Insert(ctx context.Context, lesson core.LessonModel) (core.LessonModel, error)
	Update(ctx context.Context, lesson core.UpdateLessonModel) error
}

type LessonService struct {
	lessonRepo LessonRepo
}

func NewLessonService(lessonRepo LessonRepo) *LessonService {
	return &LessonService{lessonRepo: lessonRepo}
}

func (s LessonService) Create(ctx context.Context, lesson core.Lesson) (core.Lesson, error) {
	newLesson, err := s.lessonRepo.Insert(ctx, core.LessonModel{
		Title:       lesson.Title,
		ClassroomId: lesson.ClassroomId,
		Active:      lesson.Active,
	})
	if err != nil {
		return core.Lesson{}, err
	}

	return core.Lesson{
		Id:          newLesson.Id,
		Title:       newLesson.Title,
		ClassroomId: newLesson.ClassroomId,
		Content:     newLesson.Content,
		Active:      newLesson.Active,
	}, nil
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
			Content:     model.Content,
			Active:      model.Active,
		})
	}

	return lessons, nil
}

func (s LessonService) Update(ctx context.Context, lesson core.UpdateLesson) error {
	return s.lessonRepo.Update(ctx, core.UpdateLessonModel{
		Id:          lesson.Id,
		Title:       lesson.Title,
		ClassroomId: lesson.ClassroomId,
		Content:     lesson.Content,
		Active:      lesson.Active,
	})
}
