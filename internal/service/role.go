package service

import (
	"context"
	"github.com/migmatore/study-platform-api/internal/core"
)

type RoleRepo interface {
	GetByName(ctx context.Context, name string) (core.RoleModel, error)
}
