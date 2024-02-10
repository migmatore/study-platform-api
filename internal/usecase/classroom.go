package usecase

import (
	"context"
	"github.com/migmatore/study-platform-api/internal/core"
)

type AdminService interface {
	AllClassrooms(ctx context.Context, studentId int) ([]core.Classroom, error)
}

type TeacherService interface {
	ById(ctx context.Context, id int) (core.User, error)
	AllClassrooms(ctx context.Context, teacherId int) ([]core.Classroom, error)
}

type StudentService interface {
	AllClassrooms(ctx context.Context, studentId int) ([]core.Classroom, error)
}

type ClassroomUseCase struct {
	teacherService TeacherService
}

func NewClassroomUseCase(teacherService TeacherService) *ClassroomUseCase {
	return &ClassroomUseCase{teacherService: teacherService}
}

func (uc *ClassroomUseCase) All(ctx context.Context, metadata core.TokenMetadata) ([]core.ClassroomResponse, error) {
	switch core.RoleType(metadata.Role) {
	case core.AdminRole:
	case core.TeacherRole:
		classrooms, err := uc.teacherService.AllClassrooms(ctx, metadata.UserId)
		if err != nil {
			return nil, err
		}

		classroomsResp := make([]core.ClassroomResponse, 0, len(classrooms))

		for _, classroom := range classrooms {
			classroomsResp = append(classroomsResp, core.ClassroomResponse{
				Id:          classroom.Id,
				Title:       classroom.Title,
				Description: classroom.Description,
				TeacherId:   classroom.TeacherId,
				MaxStudents: classroom.MaxStudents,
			})
		}

		return classroomsResp, nil

	case core.StudentRole:
	}

	return nil, nil
}
