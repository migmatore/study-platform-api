package service

import (
	"github.com/migmatore/study-platform-api/config"
)

type Deps struct {
	TransactorRepo  TransactionRepo
	UserRepo        UserRepo
	RoleRepo        RoleRepo
	InstitutionRepo InstitutionRepo
	ClassroomRepo   ClassroomRepo
	LessonRepo      LessonRepo
}

type Service struct {
	config      *config.Config
	Transaction *TransactionService
	User        *UserService
	Institution *InstitutionService
	Token       *TokenService
	Teacher     *TeacherService
	Student     *StudentService
	Classroom   *ClassroomService
	Lesson      *LessonService
}

func New(config *config.Config, deps Deps) *Service {
	return &Service{
		Transaction: NewTransactionService(deps.TransactorRepo),
		User:        NewUserService(deps.UserRepo, deps.RoleRepo),
		Institution: NewInstitutionService(deps.InstitutionRepo),
		Token:       NewTokenService(config),
		Teacher:     NewTeacherService(deps.ClassroomRepo, deps.UserRepo, deps.RoleRepo),
		Student:     NewStudentService(deps.ClassroomRepo),
		Classroom:   NewClassroomService(deps.ClassroomRepo, deps.UserRepo),
		Lesson:      NewLessonService(deps.LessonRepo),
	}
}
