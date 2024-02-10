package service

import (
	"context"
	"github.com/migmatore/study-platform-api/internal/core"
)

type ClassroomRepo interface {
	TeacherClassrooms(ctx context.Context, teacherId int) ([]core.ClassroomModel, error)
	ById(ctx context.Context, id int) (core.ClassroomModel, error)
}

type ClassroomService struct {
	classroomRepo ClassroomRepo
}

func NewClassroomService(classroomRepo ClassroomRepo) *ClassroomService {
	return &ClassroomService{classroomRepo: classroomRepo}
}

func (s ClassroomService) ById(ctx context.Context, id int) (core.Classroom, error) {
	classroomModel, err := s.classroomRepo.ById(ctx, id)
	if err != nil {
		return core.Classroom{}, err
	}

	return core.Classroom{
		Id:          classroomModel.Id,
		Title:       classroomModel.Title,
		Description: classroomModel.Description,
		TeacherId:   classroomModel.TeacherId,
		MaxStudents: classroomModel.MaxStudents,
	}, nil
}
