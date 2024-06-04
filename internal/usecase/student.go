package usecase

import (
	"context"
	"github.com/migmatore/study-platform-api/internal/apperrors"
	"github.com/migmatore/study-platform-api/internal/core"
	"golang.org/x/crypto/bcrypt"
)

type StudentService interface {
	AllClassrooms(ctx context.Context, studentId int) ([]core.Classroom, error)
}

type StudentTeacherService interface {
	Students(ctx context.Context, teacherId int) ([]core.Student, error)
}

type StudentUserService interface {
	ById(ctx context.Context, id int) (core.User, error)
	Create(ctx context.Context, user core.User) (core.User, error)
}

type StudentUseCase struct {
	studentTeacherService StudentTeacherService
	studentUserService    StudentUserService
}

func NewStudentsUseCase(studentTeacherService TeacherService, studentUserService StudentUserService) *StudentUseCase {
	return &StudentUseCase{studentTeacherService: studentTeacherService, studentUserService: studentUserService}
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

func (uc StudentUseCase) Create(ctx context.Context, metadata core.TokenMetadata, req core.CreateStudentRequest) (core.StudentResponse, error) {
	if core.RoleType(metadata.Role) == core.StudentRole {
		return core.StudentResponse{}, apperrors.AccessDenied
	}

	user, err := uc.studentUserService.ById(ctx, metadata.UserId)
	if err != nil {
		return core.StudentResponse{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return core.StudentResponse{}, err
	}

	student, err := uc.studentUserService.Create(ctx, core.User{
		FullName:     req.FullName,
		Phone:        req.Phone,
		Email:        req.Email,
		PasswordHash: string(hash),
		Role:         core.StudentRole,
		Institution:  user.Institution,
	})
	if err != nil {
		return core.StudentResponse{}, err
	}

	return core.StudentResponse{
		Id:           student.Id,
		FullName:     student.FullName,
		Phone:        student.Phone,
		Email:        student.Email,
		ClassroomsId: nil,
	}, nil
}
