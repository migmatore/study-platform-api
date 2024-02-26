package repository

import (
	"context"
	"github.com/migmatore/study-platform-api/internal/core"
	"github.com/migmatore/study-platform-api/internal/repository/psql"
	"github.com/migmatore/study-platform-api/pkg/logger"
	"github.com/migmatore/study-platform-api/pkg/utils"
)

type LessonRepo struct {
	logger logger.Logger
	pool   psql.AtomicPoolClient
}

func NewLessonRepo(logger logger.Logger, pool psql.AtomicPoolClient) *LessonRepo {
	return &LessonRepo{logger: logger, pool: pool}
}

func (r LessonRepo) Insert(ctx context.Context, lesson core.LessonModel) (core.LessonModel, error) {
	q := `INSERT INTO lessons(title, classroom_id, content, active) VALUES($1, $2, $3, $4) 
			RETURNING id, title, classroom_id, content, active`

	newLesson := core.LessonModel{}

	if err := r.pool.QueryRow(
		ctx,
		q,
		lesson.Title,
		lesson.ClassroomId,
		lesson.Content,
		lesson.Active,
	).Scan(
		&newLesson.Id,
		&newLesson.Title,
		&newLesson.ClassroomId,
		&newLesson.Content,
		&newLesson.Active,
	); err != nil {
		if err := utils.ParsePgError(err); err != nil {
			r.logger.Errorf("Error: %v", err)
			return core.LessonModel{}, err
		}

		r.logger.Errorf("Query error. %v", err)
		return core.LessonModel{}, err
	}

	return newLesson, nil
}

func (r LessonRepo) All(ctx context.Context, classroomId int) ([]core.LessonModel, error) {
	q := `select id, title, classroom_id, content, active from lessons where classroom_id = $1 ORDER BY id`

	lessons := make([]core.LessonModel, 0)

	rows, err := r.pool.Query(ctx, q, classroomId)
	if err != nil {
		r.logger.Errorf("Query error. %v", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		lesson := core.LessonModel{}

		err := rows.Scan(
			&lesson.Id,
			&lesson.Title,
			&lesson.ClassroomId,
			&lesson.Content,
			&lesson.Active,
		)
		if err != nil {
			r.logger.Errorf("Query error. %v", err)
			return nil, err
		}

		lessons = append(lessons, lesson)
	}

	return lessons, nil
}

func (r LessonRepo) Update(ctx context.Context, lesson core.UpdateLessonModel) error {
	updateQuery := psql.NewSQLUpdateBuilder("lessons")

	if lesson.Title != nil {
		updateQuery.AddUpdateColumn("title", lesson.Title)
	}

	if lesson.ClassroomId != nil {
		updateQuery.AddUpdateColumn("classroom_id", lesson.ClassroomId)
	}

	if lesson.Content != nil {
		updateQuery.AddUpdateColumn("content", lesson.Content)
	}

	if lesson.Active != nil {
		updateQuery.AddUpdateColumn("active", lesson.Active)
	}

	updateQuery.AddWhere("id", lesson.Id)

	if _, err := r.pool.Exec(ctx, updateQuery.GetQuery(), updateQuery.GetValues()...); err != nil {
		if err := utils.ParsePgError(err); err != nil {
			r.logger.Errorf("Error: %v", err)
			return err
		}

		r.logger.Errorf("Query error. %v", err)
		return err
	}

	return nil
}
