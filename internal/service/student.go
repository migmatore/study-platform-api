package service

import (
	"context"
	"github.com/migmatore/study-platform-api/internal/core"
)

type StudentClassroomRepo interface {
	StudentClassrooms(ctx context.Context, studentId int) ([]core.ClassroomModel, error)
}

type StudentUserRepo interface {
	ByInstitutionId(ctx context.Context, institutionId int) ([]core.UserModel, error)
}

type StudentRoleRepo interface {
	ByName(ctx context.Context, name string) (core.RoleModel, error)
}

type StudentService struct {
	classroomRepo StudentClassroomRepo
	userRepo      StudentUserRepo
	roleRepo      StudentRoleRepo
}

func NewStudentService(
	classroomRepo StudentClassroomRepo,
	userRepo StudentUserRepo,
	roleRepo StudentRoleRepo,
) *StudentService {
	return &StudentService{classroomRepo: classroomRepo, userRepo: userRepo, roleRepo: roleRepo}
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

func (s StudentService) ByInstitutionId(ctx context.Context, institutionId int) ([]core.Student, error) {
	usersModel, err := s.userRepo.ByInstitutionId(ctx, institutionId)
	if err != nil {
		return nil, err
	}

	studentRole, err := s.roleRepo.ByName(ctx, string(core.StudentRole))
	if err != nil {
		return nil, err
	}

	students := make([]core.Student, 0, len(usersModel))

	for _, model := range usersModel {
		if model.RoleId != studentRole.Id {
			continue
		}

		students = append(students, core.Student{
			Id:           model.Id,
			FullName:     model.FullName,
			Phone:        model.Phone,
			Email:        model.Email,
			ClassroomsId: nil,
		})
	}

	return students, err
}
