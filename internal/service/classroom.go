package service

import (
	"context"
	"github.com/migmatore/study-platform-api/internal/core"
)

type ClassroomRepo interface {
	Create(ctx context.Context, classroom core.ClassroomModel) (core.ClassroomModel, error)
	TeacherClassrooms(ctx context.Context, teacherId int) ([]core.ClassroomModel, error)
	StudentClassrooms(ctx context.Context, studentId int) ([]core.ClassroomModel, error)
	ById(ctx context.Context, id int) (core.ClassroomModel, error)
	IsIn(ctx context.Context, classroomId, studentId int) (bool, error)
	Students(ctx context.Context, classroomId int) ([]core.UserModel, error)
	StudentsByClassroomsId(ctx context.Context, ids []int) ([]core.StudentModel, error)
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

func (s ClassroomService) Create(ctx context.Context, classroom core.Classroom) (core.Classroom, error) {
	classroomModel, err := s.classroomRepo.Create(ctx, core.ClassroomModel{
		Id:          classroom.Id,
		Title:       classroom.Title,
		Description: classroom.Description,
		TeacherId:   classroom.TeacherId,
		MaxStudents: classroom.MaxStudents,
	})
	if err != nil {
		return core.Classroom{}, err
	}

	return core.Classroom{
		Id:          classroomModel.Id,
		Title:       classroom.Title,
		Description: classroom.Description,
		TeacherId:   classroom.TeacherId,
		MaxStudents: classroom.MaxStudents,
	}, nil
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

func (s ClassroomService) IsBelongs(ctx context.Context, classroomId, teacherId int) (bool, error) {
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

func (s ClassroomService) IsIn(ctx context.Context, classroomId, studentId int) (bool, error) {
	return s.classroomRepo.IsIn(ctx, classroomId, studentId)
}

func (s ClassroomService) Students(ctx context.Context, classroomId int) ([]core.Student, error) {
	usersModel, err := s.classroomRepo.Students(ctx, classroomId)
	if err != nil {
		return nil, err
	}

	students := make([]core.Student, 0, len(usersModel))

	for _, user := range usersModel {
		students = append(students, core.Student{
			Id:       user.Id,
			FullName: user.FullName,
			Phone:    user.Phone,
			Email:    user.Email,
		})
	}

	return students, nil
}
