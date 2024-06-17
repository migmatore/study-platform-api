package service

import (
	"context"
	"github.com/migmatore/study-platform-api/internal/core"
)

type TeacherUserRepo interface {
	ById(ctx context.Context, id int) (core.UserModel, error)
	ByInstitutionId(ctx context.Context, institutionId int) ([]core.UserModel, error)
}

type TeacherClassroomRepo interface {
	TeacherClassrooms(ctx context.Context, teacherId int) ([]core.ClassroomModel, error)
	StudentsByClassroomsId(ctx context.Context, ids []int) ([]core.StudentModel, error)
}

type TeacherRoleRepo interface {
	ById(ctx context.Context, id int) (core.RoleModel, error)
	ByName(ctx context.Context, name string) (core.RoleModel, error)
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

func (s TeacherService) All(ctx context.Context, institutionId int) ([]core.Teacher, error) {
	usersModel, err := s.userRepo.ByInstitutionId(ctx, institutionId)
	if err != nil {
		return nil, err
	}

	teacherRole, err := s.roleRepo.ByName(ctx, string(core.TeacherRole))
	if err != nil {
		return nil, err
	}

	teachers := make([]core.Teacher, 0, len(usersModel))

	for _, user := range usersModel {
		if user.RoleId != teacherRole.Id {
			continue
		}

		teachers = append(teachers, core.Teacher{
			Id:       user.Id,
			FullName: user.FullName,
			Phone:    user.Phone,
			Email:    user.Email,
		})
	}

	return teachers, nil
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
		Id:            userModel.Id,
		FullName:      userModel.FullName,
		Phone:         userModel.Phone,
		Email:         userModel.Email,
		PasswordHash:  userModel.PasswordHash,
		Role:          core.RoleType(role.Name),
		InstitutionId: nil,
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

func (s TeacherService) Students(ctx context.Context, teacherId int) ([]core.Student, error) {
	classroomsModel, err := s.classroomRepo.TeacherClassrooms(ctx, teacherId)
	if err != nil {
		return nil, err
	}

	ids := make([]int, 0, len(classroomsModel))

	for _, model := range classroomsModel {
		ids = append(ids, model.Id)
	}

	studentsModel, err := s.classroomRepo.StudentsByClassroomsId(ctx, ids)
	if err != nil {
		return nil, err
	}

	students := make([]core.Student, 0, len(studentsModel))

	for _, student := range studentsModel {
		students = append(students, core.Student{
			Id:           student.Id,
			FullName:     student.FullName,
			Phone:        student.Phone,
			Email:        student.Email,
			ClassroomsId: student.ClassroomsId,
		})
	}

	return students, nil
}
