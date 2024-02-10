package service

import (
	"context"
	"github.com/migmatore/study-platform-api/internal/core"
)

type RoleRepo interface {
	ByName(ctx context.Context, name string) (core.RoleModel, error)
	ById(ctx context.Context, id int) (core.RoleModel, error)
}
