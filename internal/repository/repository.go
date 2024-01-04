package repository

import (
	"github.com/migmatore/study-platform-api/internal/repository/psql"
	"github.com/migmatore/study-platform-api/pkg/logger"
)

type Repository struct {
	Logger      logger.Logger
	Transaction *TransactionRepo
	User        *UserRepo
	Role        *RoleRepo
	Institution *InstitutionRepo
}

func New(logger logger.Logger, pool psql.AtomicPoolClient) *Repository {
	return &Repository{
		Transaction: NewTransactor(pool),
		User:        NewUserRepo(logger, pool),
		Role:        NewRoleRepo(logger, pool),
		Institution: NewInstitutionRepo(logger, pool),
	}
}
