package usecase

type Deps struct {
	UserService UserService
}

type UseCase struct {
	Auth *AuthUseCase
}

func New(deps Deps) *UseCase {
	return &UseCase{Auth: NewAuthUseCase(deps.UserService)}
}
