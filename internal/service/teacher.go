package service

import (
	"context"
	"github.com/migmatore/study-platform-api/internal/core"
)

type TeacherUserRepo interface {
	ById(ctx context.Context, id int) (core.UserModel, error)
}

type TeacherClassroomRepo interface {
	TeacherClassrooms(ctx context.Context, teacherId int) ([]core.ClassroomModel, error)
}

type TeacherRoleRepo interface {
	ById(ctx context.Context, id int) (core.RoleModel, error)
}

type TeacherService struct {
	classroomRepo TeacherClassroomRepo
	userRepo      TeacherUserRepo
	roleRepo      TeacherRoleRepo
}

func NewTeacherService(
	classroomRepo TeacherClassroomRepo,
	userRepo TeacherUserRepo,
	roleRepo TeacherRoleRepo,
) *TeacherService {
	return &TeacherService{classroomRepo: classroomRepo, userRepo: userRepo, roleRepo: roleRepo}
}

func (s TeacherService) ById(ctx context.Context, id int) (core.User, error) {
	userModel, err := s.userRepo.ById(ctx, id)
	if err != nil {
		return core.User{}, err
	}

	role, err := s.roleRepo.ById(ctx, userModel.RoleId)
	if err != nil {
		return core.User{}, err
	}

	return core.User{
		Id:           userModel.Id,
		FullName:     userModel.FullName,
		Phone:        userModel.Phone,
		Email:        userModel.Email,
		PasswordHash: userModel.PasswordHash,
		Role:         core.RoleType(role.Name),
		Institution:  nil,
	}, nil
}

func (s TeacherService) AllClassrooms(ctx context.Context, teacherId int) ([]core.Classroom, error) {
	classroomsModel, err := s.classroomRepo.TeacherClassrooms(ctx, teacherId)
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
