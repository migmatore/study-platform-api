package usecase

import (
	"context"
	"github.com/migmatore/study-platform-api/internal/apperrors"
	"github.com/migmatore/study-platform-api/internal/core"
	"golang.org/x/crypto/bcrypt"
)

type TeacherService interface {
	All(ctx context.Context, institutionId int) ([]core.Teacher, error)
	ById(ctx context.Context, id int) (core.User, error)
	AllClassrooms(ctx context.Context, teacherId int) ([]core.Classroom, error)
	Students(ctx context.Context, teacherId int) ([]core.Student, error)
}

type TeacherUserService interface {
	ById(ctx context.Context, id int) (core.User, error)
	Create(ctx context.Context, user core.User) (core.User, error)
	IsExist(ctx context.Context, email string) (bool, error)
	Delete(ctx context.Context, id int) error
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

func (uc TeacherUseCase) Create(
	ctx context.Context,
	metadata core.TokenMetadata,
	req core.CreateTeacherRequest,
) (core.TeacherResponse, error) {
	if core.RoleType(metadata.Role) != core.AdminRole {
		return core.TeacherResponse{}, apperrors.AccessDenied
	}

	exist, err := uc.userService.IsExist(ctx, req.Email)
	if err != nil {
		return core.TeacherResponse{}, err
	}

	if exist {
		return core.TeacherResponse{}, apperrors.EntityAlreadyExist
	}

	admin, err := uc.userService.ById(ctx, metadata.UserId)
	if err != nil {
		return core.TeacherResponse{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return core.TeacherResponse{}, err
	}

	newTeacher, err := uc.userService.Create(ctx, core.User{
		FullName:      req.FullName,
		Phone:         req.Phone,
		Email:         req.Email,
		PasswordHash:  string(hash),
		Role:          core.TeacherRole,
		InstitutionId: admin.InstitutionId,
	})
	if err != nil {
		return core.TeacherResponse{}, err
	}

	return core.TeacherResponse{
		Id:       newTeacher.Id,
		FullName: newTeacher.FullName,
		Phone:    newTeacher.Phone,
		Email:    newTeacher.Email,
	}, nil
}

func (uc TeacherUseCase) Delete(ctx context.Context, metadata core.TokenMetadata, id int) error {
	if core.RoleType(metadata.Role) != core.AdminRole {
		return apperrors.AccessDenied
	}

	if err := uc.userService.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
