package usecase

import (
	"context"
	"github.com/migmatore/study-platform-api/internal/core"
)

type TeacherService interface {
	ById(ctx context.Context, id int) (core.User, error)
	AllClassrooms(ctx context.Context, teacherId int) ([]core.Classroom, error)
	Students(ctx context.Context, teacherId int) ([]core.Student, error)
}
