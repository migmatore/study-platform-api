package service

import (
	"context"
	"github.com/migmatore/study-platform-api/internal/core"
)

type StudentClassroomRepo interface {
	StudentClassrooms(ctx context.Context, studentId int) ([]core.ClassroomModel, error)
}

type StudentService struct {
	classroomRepo StudentClassroomRepo
}

func NewStudentService(classroomRepo StudentClassroomRepo) *StudentService {
	return &StudentService{classroomRepo: classroomRepo}
}

func (s StudentService) AllClassrooms(ctx context.Context, studentId int) ([]core.Classroom, error) {
	classroomsModel, err := s.classroomRepo.StudentClassrooms(ctx, studentId)
	if err != nil {
		return nil, err
	}

	classrooms := make([]core.Classroom, 0, len(classroomsModel))

	for _, model := range classroomsModel {
		classrooms = append(classrooms, core.Classroom{
			Id:          model.Id,
			Title:       model.Title,
			Description: model.Description,
			TeacherId:   model.TeacherId,
			MaxStudents: model.MaxStudents,
		})
	}

	return classrooms, nil
}
