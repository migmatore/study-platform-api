package usecase

type Deps struct {
	TransactionService TransactionService
	UserService        UserService
	InstitutionService InstitutionService
	TokenService       TokenService
	TeacherService     TeacherService
	StudentService     StudentService
	LessonService      LessonService
	ClassroomService   ClassroomService
}

type UseCase struct {
	Auth      *AuthUseCase
	User      *UserUseCase
	Classroom *ClassroomUseCase
	Lesson    *LessonUseCase
	Student   *StudentUseCase
}

func New(deps Deps) *UseCase {
	return &UseCase{
		Auth: NewAuthUseCase(
			deps.TransactionService,
			deps.UserService,
			deps.InstitutionService,
			deps.TokenService,
		),
		User:      NewUserUseCase(deps.UserService),
		Classroom: NewClassroomUseCase(deps.ClassroomService, deps.TeacherService, deps.StudentService),
		Lesson:    NewLessonUseCase(deps.LessonService, deps.ClassroomService, deps.TeacherService),
		Student: NewStudentsUseCase(
			deps.TransactionService,
			deps.TeacherService,
			deps.UserService,
			deps.ClassroomService,
		),
	}
}
