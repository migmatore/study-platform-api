package service

import (
	"context"
	"github.com/migmatore/study-platform-api/internal/core"
)

type LessonRepo interface {
	All(ctx context.Context, classroomId int) ([]core.LessonModel, error)
	ById(ctx context.Context, lessonId int) (core.LessonModel, error)
	Insert(ctx context.Context, lesson core.LessonModel) (core.LessonModel, error)
	Update(ctx context.Context, lesson core.UpdateLessonModel) error
	Delete(ctx context.Context, id int) error
}

type LessonClassroomRepo interface {
	ById(ctx context.Context, id int) (core.ClassroomModel, error)
}

type LessonService struct {
	lessonRepo    LessonRepo
	classroomRepo LessonClassroomRepo
}

func NewLessonService(lessonRepo LessonRepo, classroomRepo LessonClassroomRepo) *LessonService {
	return &LessonService{lessonRepo: lessonRepo, classroomRepo: classroomRepo}
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

func (s LessonService) ById(ctx context.Context, lessonId int) (core.Lesson, error) {
	model, err := s.lessonRepo.ById(ctx, lessonId)
	if err != nil {
		return core.Lesson{}, err
	}

	return core.Lesson{
		Id:          model.Id,
		Title:       model.Title,
		ClassroomId: model.ClassroomId,
		Content:     model.Content,
		Active:      model.Active,
	}, nil
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

func (s LessonService) Delete(ctx context.Context, id int) error {
	return s.lessonRepo.Delete(ctx, id)
}
func (s LessonService) IsBelongs(ctx context.Context, lessonId int, teacherId int) (bool, error) {
	lesson, err := s.lessonRepo.ById(ctx, lessonId)
	if err != nil {
		return false, err
	}

	classroom, err := s.classroomRepo.ById(ctx, lesson.ClassroomId)
	if err != nil {
		return false, err
	}

	return classroom.TeacherId == teacherId, nil
}
