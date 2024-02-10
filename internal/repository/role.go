package repository

import (
	"context"
	"github.com/migmatore/study-platform-api/internal/core"
	"github.com/migmatore/study-platform-api/internal/repository/psql"
	"github.com/migmatore/study-platform-api/pkg/logger"
	"github.com/migmatore/study-platform-api/pkg/utils"
)

type RoleRepo struct {
	logger logger.Logger
	pool   psql.AtomicPoolClient
}

func NewRoleRepo(logger logger.Logger, pool psql.AtomicPoolClient) *RoleRepo {
	return &RoleRepo{logger: logger, pool: pool}
}

func (r RoleRepo) ByName(ctx context.Context, name string) (core.RoleModel, error) {
	q := `SELECT id, name FROM roles WHERE name = $1`

	var role core.RoleModel

	if err := r.pool.QueryRow(ctx, q, name).Scan(&role.Id, &role.Name); err != nil {
		if err := utils.ParsePgError(err); err != nil {
			r.logger.Errorf("Error: %v", err)
			return role, err
		}

		r.logger.Errorf("Query error. %v", err)
		return role, err
	}

	return role, nil
}

func (r RoleRepo) ById(ctx context.Context, id int) (core.RoleModel, error) {
	q := `SELECT id, name FROM roles WHERE id = $1`

	var role core.RoleModel

	if err := r.pool.QueryRow(ctx, q, id).Scan(&role.Id, &role.Name); err != nil {
		if err := utils.ParsePgError(err); err != nil {
			r.logger.Errorf("Error: %v", err)
			return role, err
		}

		r.logger.Errorf("Query error. %v", err)
		return role, err
	}

	return role, nil
}
