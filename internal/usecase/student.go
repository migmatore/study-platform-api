package usecase

import (
	"context"
	"github.com/migmatore/study-platform-api/internal/apperrors"
	"github.com/migmatore/study-platform-api/internal/core"
	"golang.org/x/crypto/bcrypt"
)

type StudentService interface {
	AllClassrooms(ctx context.Context, studentId int) ([]core.Classroom, error)
	ByInstitutionId(ctx context.Context, institutionId int) ([]core.Student, error)
}

type StudentTeacherService interface {
	Students(ctx context.Context, teacherId int) ([]core.Student, error)
}

type StudentUserService interface {
	ById(ctx context.Context, id int) (core.User, error)
	Create(ctx context.Context, user core.User) (core.User, error)
	IsExist(ctx context.Context, email string) (bool, error)
	Delete(ctx context.Context, id int) error
}

type StudentClassroomService interface {
	AddStudent(ctx context.Context, studentId int, classroomsId []int) error
	ById(ctx context.Context, id int) (core.Classroom, error)
	Students(ctx context.Context, classroomId int) ([]core.Student, error)
}

type StudentUseCase struct {
	transactionService      TransactionService
	studentService          StudentService
	studentTeacherService   StudentTeacherService
	studentUserService      StudentUserService
	studentClassroomService StudentClassroomService
}

func NewStudentsUseCase(
	transactionService TransactionService,
	studentService StudentService,
	studentTeacherService TeacherService,
	studentUserService StudentUserService,
	studentClassroomService StudentClassroomService,
) *StudentUseCase {
	return &StudentUseCase{
		transactionService:      transactionService,
		studentService:          studentService,
		studentTeacherService:   studentTeacherService,
		studentUserService:      studentUserService,
		studentClassroomService: studentClassroomService,
	}
}

func (uc StudentUseCase) All(ctx context.Context, metadata core.TokenMetadata) ([]core.StudentResponse, error) {
	switch core.RoleType(metadata.Role) {
	case core.AdminRole:
		admin, err := uc.studentUserService.ById(ctx, metadata.UserId)
		if err != nil {
			return nil, err
		}

		students, err := uc.studentService.ByInstitutionId(ctx, *admin.InstitutionId)
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
				ClassroomsId: nil,
			})
		}

		return studentsResponse, nil
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

func (uc StudentUseCase) Create(
	ctx context.Context,
	metadata core.TokenMetadata,
	req core.CreateStudentRequest,
) (core.StudentResponse, error) {
	if core.RoleType(metadata.Role) == core.StudentRole {
		return core.StudentResponse{}, apperrors.AccessDenied
	}

	exist, err := uc.studentUserService.IsExist(ctx, req.Email)
	if err != nil {
		return core.StudentResponse{}, err
	}

	if exist {
		return core.StudentResponse{}, apperrors.EntityAlreadyExist
	}

	for _, classroomId := range req.ClassroomsId {
		classroom, err := uc.studentClassroomService.ById(ctx, classroomId)
		if err != nil {
			return core.StudentResponse{}, err
		}

		students, err := uc.studentTeacherService.Students(ctx, metadata.UserId)
		if err != nil {
			return core.StudentResponse{}, err
		}

		if len(students)+1 > classroom.MaxStudents {
			return core.StudentResponse{}, apperrors.NumberOfStudentsExceeded
		}
	}

	user, err := uc.studentUserService.ById(ctx, metadata.UserId)
	if err != nil {
		return core.StudentResponse{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return core.StudentResponse{}, err
	}

	var student core.User

	if err := uc.transactionService.WithinTransaction(ctx, func(txCtx context.Context) error {
		student, err = uc.studentUserService.Create(txCtx, core.User{
			FullName:      req.FullName,
			Phone:         req.Phone,
			Email:         req.Email,
			PasswordHash:  string(hash),
			Role:          core.StudentRole,
			InstitutionId: user.InstitutionId,
		})
		if err != nil {
			return err
		}

		if err := uc.studentClassroomService.AddStudent(txCtx, student.Id, req.ClassroomsId); err != nil {
			return err
		}

		return nil
	}); err != nil {
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

func (uc StudentUseCase) Delete(ctx context.Context, metadata core.TokenMetadata, id int) error {
	switch core.RoleType(metadata.Role) {
	case core.AdminRole:
		return apperrors.AccessDenied
	case core.TeacherRole:
		if err := uc.studentUserService.Delete(ctx, id); err != nil {
			return err
		}
	case core.StudentRole:
		return apperrors.AccessDenied
	}

	return nil
}
