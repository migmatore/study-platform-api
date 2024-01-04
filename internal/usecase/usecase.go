package usecase

type Deps struct {
	TransactionService TransactionService
	UserService        UserService
	InstitutionService InstitutionService
}

type UseCase struct {
	Auth *AuthUseCase
}

func New(deps Deps) *UseCase {
	return &UseCase{Auth: NewAuthUseCase(
		deps.TransactionService,
		deps.UserService,
		deps.InstitutionService,
	)}
}
