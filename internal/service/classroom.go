package service

import (
	"context"
	"github.com/migmatore/study-platform-api/internal/core"
)

type ClassroomRepo interface {
	TeacherClassrooms(ctx context.Context, teacherId int) ([]core.ClassroomModel, error)
	ById(ctx context.Context, id int) (core.ClassroomModel, error)
}

type ClassroomTeacherUserRepo interface {
	ById(ctx context.Context, id int) (core.UserModel, error)
}

type ClassroomService struct {
	classroomRepo ClassroomRepo
	teacherRepo   ClassroomTeacherUserRepo
}

func NewClassroomService(classroomRepo ClassroomRepo, teacherRepo ClassroomTeacherUserRepo) *ClassroomService {
	return &ClassroomService{classroomRepo: classroomRepo, teacherRepo: teacherRepo}
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

func (s ClassroomService) IsBelongs(ctx context.Context, classroomId int, teacherId int) (bool, error) {
	teacher, err := s.teacherRepo.ById(ctx, teacherId)
	if err != nil {
		return false, err
	}

	classroom, err := s.classroomRepo.ById(ctx, classroomId)
	if err != nil {
		return false, err
	}

	return classroom.TeacherId == teacher.Id, nil
}
