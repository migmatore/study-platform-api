package usecase

import (
	"context"
	"github.com/migmatore/study-platform-api/internal/apperrors"
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
	classroomService ClassroomService
	teacherService   TeacherService
	studentService   StudentService
}

func NewClassroomUseCase(
	classroomService ClassroomService,
	teacherService TeacherService,
	studentService StudentService,
) *ClassroomUseCase {
	return &ClassroomUseCase{
		classroomService: classroomService,
		teacherService:   teacherService,
		studentService:   studentService,
	}
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
		classrooms, err := uc.studentService.AllClassrooms(ctx, metadata.UserId)
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
	}

	return nil, nil
}

func (uc *ClassroomUseCase) Students(
	ctx context.Context,
	metadata core.TokenMetadata,
	classroomId int,
) ([]core.StudentResponse, error) {
	if core.RoleType(metadata.Role) == core.StudentRole {
		return nil, apperrors.AccessDenied
	}

	belongs, err := uc.classroomService.IsBelongs(ctx, classroomId, metadata.UserId)
	if err != nil {
		return nil, err
	}

	if !belongs {
		return nil, apperrors.AccessDenied
	}

	students, err := uc.classroomService.Students(ctx, classroomId)
	if err != nil {
		return nil, err
	}

	studentsResp := make([]core.StudentResponse, 0, len(students))

	for _, student := range students {
		studentsResp = append(studentsResp, core.StudentResponse{
			Id:       student.Id,
			FullName: student.FullName,
			Phone:    student.Phone,
			Email:    student.Email,
		})
	}

	return studentsResp, nil
}
