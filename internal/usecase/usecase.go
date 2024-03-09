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
	Classroom *ClassroomUseCase
	Lesson    *LessonUseCase
}

func New(deps Deps) *UseCase {
	return &UseCase{
		Auth: NewAuthUseCase(
			deps.TransactionService,
			deps.UserService,
			deps.InstitutionService,
			deps.TokenService,
		),
		Classroom: NewClassroomUseCase(deps.ClassroomService, deps.TeacherService, deps.StudentService),
		Lesson:    NewLessonUseCase(deps.LessonService, deps.ClassroomService, deps.TeacherService),
	}
}
