package usecase

import (
	"context"
	"github.com/migmatore/study-platform-api/internal/apperrors"
	"github.com/migmatore/study-platform-api/internal/core"
)

type TeacherService interface {
	All(ctx context.Context, institutionId int) ([]core.Teacher, error)
	ById(ctx context.Context, id int) (core.User, error)
	AllClassrooms(ctx context.Context, teacherId int) ([]core.Classroom, error)
	Students(ctx context.Context, teacherId int) ([]core.Student, error)
}

type TeacherUserService interface {
	ById(ctx context.Context, id int) (core.User, error)
}

type TeacherUseCase struct {
	teacherService TeacherService
	userService    TeacherUserService
}

func NewTeacherUseCase(teacherService TeacherService, userService TeacherUserService) *TeacherUseCase {
	return &TeacherUseCase{teacherService: teacherService, userService: userService}
}

func (uc TeacherUseCase) All(ctx context.Context, metadata core.TokenMetadata) ([]core.TeacherResponse, error) {
	if core.RoleType(metadata.Role) != core.AdminRole {
		return nil, apperrors.AccessDenied
	}

	admin, err := uc.userService.ById(ctx, metadata.UserId)
	if err != nil {
		return nil, err
	}

	teachers, err := uc.teacherService.All(ctx, *admin.InstitutionId)
	if err != nil {
		return nil, err
	}

	teachersResp := make([]core.TeacherResponse, 0, len(teachers))

	for _, teacher := range teachers {
		teachersResp = append(teachersResp, core.TeacherResponse{
			Id:       teacher.Id,
			FullName: teacher.FullName,
			Phone:    teacher.Phone,
			Email:    teacher.Email,
		})
	}

	return teachersResp, nil
}
