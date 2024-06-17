package repository

import (
	"context"
	"github.com/migmatore/study-platform-api/internal/core"
	"github.com/migmatore/study-platform-api/internal/repository/psql"
	"github.com/migmatore/study-platform-api/pkg/logger"
	"github.com/migmatore/study-platform-api/pkg/utils"
)

type UserRepo struct {
	logger logger.Logger
	pool   psql.AtomicPoolClient
}

func NewUserRepo(logger logger.Logger, pool psql.AtomicPoolClient) *UserRepo {
	return &UserRepo{logger: logger, pool: pool}
}

func (r UserRepo) IsExist(ctx context.Context, email string) (bool, error) {
	q := `SELECT EXISTS(SELECT * FROM users WHERE email = $1)`

	var exist bool

	if err := r.pool.QueryRow(ctx, q, email).Scan(&exist); err != nil {
		if err := utils.ParsePgError(err); err != nil {
			r.logger.Errorf("Error: %v", err)
			return exist, err
		}

		r.logger.Errorf("Query error. %v", err)
		return exist, err
	}

	return exist, nil
}

func (r UserRepo) IsExistById(ctx context.Context, id int) (bool, error) {
	q := `SELECT EXISTS(SELECT * FROM users WHERE id = $1)`

	var exist bool

	if err := r.pool.QueryRow(ctx, q, id).Scan(&exist); err != nil {
		if err := utils.ParsePgError(err); err != nil {
			r.logger.Errorf("Error: %v", err)
			return exist, err
		}

		r.logger.Errorf("Query error. %v", err)
		return exist, err
	}

	return exist, nil
}

func (r UserRepo) Create(ctx context.Context, user core.UserModel) (core.UserModel, error) {
	q := `INSERT INTO users(full_name, phone, email, password_hash, role_id, institution_id) 
		  VALUES ($1, $2, $3, $4, $5, $6)
          RETURNING id, full_name, phone, email, password_hash, role_id, institution_id`

	var u core.UserModel

	if err := r.pool.QueryRow(
		ctx,
		q,
		user.FullName,
		user.Phone,
		user.Email,
		user.PasswordHash,
		user.RoleId,
		user.InstitutionId,
	).Scan(
		&u.Id,
		&u.FullName,
		&u.Phone,
		&u.Email,
		&u.PasswordHash,
		&u.RoleId,
		&u.InstitutionId,
	); err != nil {
		if err := utils.ParsePgError(err); err != nil {
			r.logger.Errorf("Error: %v", err)
			return u, err
		}

		r.logger.Errorf("Query error. %v", err)
		return u, err
	}

	return u, nil
}

func (r UserRepo) ByEmail(ctx context.Context, email string) (core.UserModel, error) {
	q := `SELECT id, full_name, phone, email, password_hash, role_id, institution_id FROM users WHERE email = $1`

	var u core.UserModel

	if err := r.pool.QueryRow(ctx, q, email).Scan(
		&u.Id,
		&u.FullName,
		&u.Phone,
		&u.Email,
		&u.PasswordHash,
		&u.RoleId,
		&u.InstitutionId,
	); err != nil {
		if err := utils.ParsePgError(err); err != nil {
			r.logger.Errorf("Error: %v", err)
			return u, err
		}

		r.logger.Errorf("Query error. %v", err)
		return u, err
	}

	return u, nil
}

func (r UserRepo) ById(ctx context.Context, id int) (core.UserModel, error) {
	q := `SELECT id, full_name, phone, email, password_hash, role_id, institution_id FROM users WHERE id = $1`

	var u core.UserModel

	if err := r.pool.QueryRow(ctx, q, id).Scan(
		&u.Id,
		&u.FullName,
		&u.Phone,
		&u.Email,
		&u.PasswordHash,
		&u.RoleId,
		&u.InstitutionId,
	); err != nil {
		if err := utils.ParsePgError(err); err != nil {
			r.logger.Errorf("Error: %v", err)
			return u, err
		}

		r.logger.Errorf("Query error. %v", err)
		return u, err
	}

	return u, nil
}

func (r UserRepo) UpdateProfile(
	ctx context.Context,
	userId int,
	profile core.UpdateUserProfileModel,
) (core.UserProfileModel, error) {
	updateQuery := psql.NewSQLUpdateBuilder("users")

	if profile.FullName != nil {
		updateQuery.AddUpdateColumn("full_name", *profile.FullName)
	}

	if profile.Phone != nil {
		updateQuery.AddUpdateColumn("phone", *profile.Phone)
	}

	if profile.Email != nil {
		updateQuery.AddUpdateColumn("email", *profile.Email)
	}

	if profile.Password != nil {
		updateQuery.AddUpdateColumn("password_hash", *profile.Password)
	}

	updateQuery.AddWhere("id", userId)
	updateQuery.AddReturning("full_name", "phone", "email")

	var p core.UserProfileModel

	if err := r.pool.QueryRow(ctx, updateQuery.GetQuery(), updateQuery.GetValues()...).Scan(
		&p.FullName,
		&p.Phone,
		&p.Email,
	); err != nil {
		if err := utils.ParsePgError(err); err != nil {
			r.logger.Errorf("Error: %v", err)
			return core.UserProfileModel{}, err
		}

		r.logger.Errorf("Query error. %v", err)
		return core.UserProfileModel{}, err
	}

	return p, nil
}

func (r UserRepo) ByInstitutionId(ctx context.Context, institutionId int) ([]core.UserModel, error) {
	q := `SELECT id, full_name, phone, email, password_hash, role_id, institution_id  FROM users WHERE institution_id = $1`

	users := make([]core.UserModel, 0)

	rows, err := r.pool.Query(ctx, q, institutionId)
	if err != nil {
		r.logger.Errorf("Query error. %v", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		user := core.UserModel{}

		err := rows.Scan(
			&user.Id,
			&user.FullName,
			&user.Phone,
			&user.Email,
			&user.PasswordHash,
			&user.RoleId,
			&user.InstitutionId,
		)
		if err != nil {
			r.logger.Errorf("Query error. %v", err)
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (r UserRepo) Delete(ctx context.Context, id int) error {
	q := `DELETE FROM users WHERE id = $1`

	if _, err := r.pool.Exec(ctx, q, id); err != nil {
		if err := utils.ParsePgError(err); err != nil {
			r.logger.Errorf("Error: %v", err)
			return err
		}

		r.logger.Errorf("Query error. %v", err)
		return err
	}

	return nil
}
