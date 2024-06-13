package usecase

import (
	"context"
	"github.com/migmatore/study-platform-api/internal/apperrors"
	"github.com/migmatore/study-platform-api/internal/core"
)

type ClassroomService interface {
	Create(ctx context.Context, classroom core.Classroom) (core.Classroom, error)
	Delete(ctx context.Context, id int) error
	ById(ctx context.Context, id int) (core.Classroom, error)
	IsBelongs(ctx context.Context, classroomId int, teacherId int) (bool, error)
	IsIn(ctx context.Context, classroomId, studentId int) (bool, error)
	Students(ctx context.Context, classroomId int) ([]core.Student, error)
	AddStudent(ctx context.Context, studentId int, classroomsId []int) error
}

type AdminService interface {
	AllClassrooms(ctx context.Context, studentId int) ([]core.Classroom, error)
}

type ClassroomTeacherService interface {
	ById(ctx context.Context, id int) (core.User, error)
	AllClassrooms(ctx context.Context, teacherId int) ([]core.Classroom, error)
}

type ClassroomStudentService interface {
	AllClassrooms(ctx context.Context, studentId int) ([]core.Classroom, error)
}

type ClassroomUseCase struct {
	classroomService ClassroomService
	teacherService   TeacherService
	studentService   ClassroomStudentService
}

func NewClassroomUseCase(
	classroomService ClassroomService,
	teacherService TeacherService,
	studentService ClassroomStudentService,
) *ClassroomUseCase {
	return &ClassroomUseCase{
		classroomService: classroomService,
		teacherService:   teacherService,
		studentService:   studentService,
	}
}

func (uc ClassroomUseCase) All(ctx context.Context, metadata core.TokenMetadata) ([]core.ClassroomResponse, error) {
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

func (uc ClassroomUseCase) Create(
	ctx context.Context,
	metadata core.TokenMetadata,
	req core.CreateClassroomRequest,
) (core.ClassroomResponse, error) {
	if core.RoleType(metadata.Role) != core.TeacherRole {
		return core.ClassroomResponse{}, apperrors.AccessDenied
	}

	newClassroom, err := uc.classroomService.Create(ctx, core.Classroom{
		Title:       req.Title,
		Description: req.Description,
		TeacherId:   metadata.UserId,
		MaxStudents: req.MaxStudents,
	})
	if err != nil {
		return core.ClassroomResponse{}, err
	}

	return core.ClassroomResponse{
		Id:          newClassroom.Id,
		Title:       newClassroom.Title,
		Description: newClassroom.Description,
		TeacherId:   newClassroom.TeacherId,
		MaxStudents: newClassroom.MaxStudents,
	}, nil
}

func (uc ClassroomUseCase) Delete(ctx context.Context, metadata core.TokenMetadata, id int) error {
	switch core.RoleType(metadata.Role) {
	case core.AdminRole:
		return apperrors.AccessDenied
	case core.TeacherRole:
		belongs, err := uc.classroomService.IsBelongs(ctx, id, metadata.UserId)
		if err != nil {
			return err
		}

		if !belongs {
			return apperrors.AccessDenied
		}

		if err := uc.classroomService.Delete(ctx, id); err != nil {
			return err
		}
	case core.StudentRole:
		return apperrors.AccessDenied
	}

	return nil
}

func (uc ClassroomUseCase) Students(
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
