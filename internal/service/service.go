package service

import "github.com/migmatore/study-platform-api/config"

type Deps struct {
	TransactorRepo  TransactionRepo
	UserRepo        UserRepo
	RoleRepo        RoleRepo
	InstitutionRepo InstitutionRepo
}

type Service struct {
	config      *config.Config
	Transaction *TransactionService
	User        *UserService
	Institution *InstitutionService
	Token       *TokenService
}

func New(config *config.Config, deps Deps) *Service {
	return &Service{
		Transaction: NewTransactionService(deps.TransactorRepo),
		User:        NewUserService(deps.UserRepo, deps.RoleRepo),
		Institution: NewInstitutionService(deps.InstitutionRepo),
		Token:       NewTokenService(config),
	}
}
