package repository

import (
	"context"
	"github.com/migmatore/study-platform-api/internal/core"
	"github.com/migmatore/study-platform-api/internal/repository/psql"
	"github.com/migmatore/study-platform-api/pkg/logger"
	"github.com/migmatore/study-platform-api/pkg/utils"
)

type InstitutionRepo struct {
	logger logger.Logger
	pool   psql.AtomicPoolClient
}

func NewInstitutionRepo(logger logger.Logger, pool psql.AtomicPoolClient) *InstitutionRepo {
	return &InstitutionRepo{logger: logger, pool: pool}
}

func (r InstitutionRepo) IsExist(ctx context.Context, name string) (bool, error) {
	q := `SELECT EXISTS(SELECT * FROM institutions WHERE name = $1)`

	var exist bool

	if err := r.pool.QueryRow(ctx, q, name).Scan(&exist); err != nil {
		if err := utils.ParsePgError(err); err != nil {
			r.logger.Errorf("Error: %v", err)
			return exist, err
		}

		r.logger.Errorf("Query error. %v", err)
		return exist, err
	}

	return exist, nil
}

func (r InstitutionRepo) Create(ctx context.Context, inst core.InstitutionModel) (core.InstitutionModel, error) {
	q := `INSERT INTO institutions(name) 
		  VALUES ($1)
          RETURNING id, name, description`

	var i core.InstitutionModel

	if err := r.pool.QueryRow(
		ctx,
		q,
		inst.Name,
	).Scan(
		&i.Id,
		&i.Name,
		&i.Description,
	); err != nil {
		if err := utils.ParsePgError(err); err != nil {
			r.logger.Errorf("Error: %v", err)
			return i, err
		}

		r.logger.Errorf("Query error. %v", err)
		return i, err
	}

	return i, nil
}

func (r InstitutionRepo) ById(ctx context.Context, id int) (core.InstitutionModel, error) {
	q := `SELECT id, name, description FROM institutions WHERE id = $1`

	var i core.InstitutionModel

	if err := r.pool.QueryRow(ctx, q, id).Scan(&i.Id, &i.Name, &i.Description); err != nil {
		if err := utils.ParsePgError(err); err != nil {
			r.logger.Errorf("Error: %v", err)
			return core.InstitutionModel{}, err
		}

		r.logger.Errorf("Query error. %v", err)
		return core.InstitutionModel{}, err
	}

	return i, nil
}
