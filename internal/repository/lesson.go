package repository

import (
	"context"
	"github.com/migmatore/study-platform-api/internal/core"
	"github.com/migmatore/study-platform-api/internal/repository/psql"
	"github.com/migmatore/study-platform-api/pkg/logger"
)

type LessonRepo struct {
	logger logger.Logger
	pool   psql.AtomicPoolClient
}

func NewLessonRepo(logger logger.Logger, pool psql.AtomicPoolClient) *LessonRepo {
	return &LessonRepo{logger: logger, pool: pool}
}

func (r LessonRepo) All(ctx context.Context, classroomId int) ([]core.LessonModel, error) {
	q := `select id, title, classroom_id, active from lessons where classroom_id = $1 ORDER BY id`

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
