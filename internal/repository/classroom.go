package repository

import (
	"context"
	"github.com/migmatore/study-platform-api/internal/core"
	"github.com/migmatore/study-platform-api/internal/repository/psql"
	"github.com/migmatore/study-platform-api/pkg/logger"
	"github.com/migmatore/study-platform-api/pkg/utils"
)

type ClassroomRepo struct {
	logger logger.Logger
	pool   psql.AtomicPoolClient
}

func NewClassroomRepo(logger logger.Logger, pool psql.AtomicPoolClient) *ClassroomRepo {
	return &ClassroomRepo{logger: logger, pool: pool}
}

func (r ClassroomRepo) ById(ctx context.Context, id int) (core.ClassroomModel, error) {
	q := `SELECT id, title, description, teacher_id, max_students FROM classrooms WHERE id = $1`

	var classroom core.ClassroomModel

	if err := r.pool.QueryRow(ctx, q, id).Scan(
		&classroom.Id,
		&classroom.Title,
		&classroom.Description,
		&classroom.TeacherId,
		&classroom.MaxStudents,
	); err != nil {
		if err := utils.ParsePgError(err); err != nil {
			r.logger.Errorf("Error: %v", err)
			return classroom, err
		}

		r.logger.Errorf("Query error. %v", err)
		return classroom, err
	}

	return classroom, nil
}

func (r ClassroomRepo) TeacherClassrooms(ctx context.Context, teacherId int) ([]core.ClassroomModel, error) {
	q := `select id, title, description, teacher_id, max_students from classrooms where teacher_id = $1`

	classrooms := make([]core.ClassroomModel, 0)

	rows, err := r.pool.Query(ctx, q, teacherId)
	if err != nil {
		r.logger.Errorf("Query error. %v", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		classroom := core.ClassroomModel{}

		err := rows.Scan(
			&classroom.Id,
			&classroom.Title,
			&classroom.Description,
			&classroom.TeacherId,
			&classroom.MaxStudents,
		)
		if err != nil {
			r.logger.Errorf("Query error. %v", err)
			return nil, err
		}

		classrooms = append(classrooms, classroom)
	}

	return classrooms, nil
}
