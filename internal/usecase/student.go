package usecase

import (
	"context"
	"github.com/migmatore/study-platform-api/internal/apperrors"
	"github.com/migmatore/study-platform-api/internal/core"
)

type StudentService interface {
	AllClassrooms(ctx context.Context, studentId int) ([]core.Classroom, error)
}

type StudentTeacherService interface {
	Students(ctx context.Context, teacherId int) ([]core.Student, error)
}

type StudentUseCase struct {
	studentTeacherService StudentTeacherService
}

func NewStudentsUseCase(studentTeacherService TeacherService) *StudentUseCase {
	return &StudentUseCase{studentTeacherService: studentTeacherService}
}

func (uc StudentUseCase) All(ctx context.Context, metadata core.TokenMetadata) ([]core.StudentResponse, error) {
	switch core.RoleType(metadata.Role) {
	case core.AdminRole:
	case core.TeacherRole:
		students, err := uc.studentTeacherService.Students(ctx, metadata.UserId)
		if err != nil {
			return nil, err
		}

		studentsResponse := make([]core.StudentResponse, 0, len(students))

		for _, student := range students {
			studentsResponse = append(studentsResponse, core.StudentResponse{
				Id:           student.Id,
				FullName:     student.FullName,
				Phone:        student.Phone,
				Email:        student.Email,
				ClassroomsId: student.ClassroomsId,
			})
		}

		return studentsResponse, nil
	case core.StudentRole:
		return nil, apperrors.AccessDenied
	}

	return nil, apperrors.AccessDenied
}
