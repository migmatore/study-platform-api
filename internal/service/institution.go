package service

import (
	"context"
	"github.com/migmatore/study-platform-api/internal/core"
)

type InstitutionRepo interface {
	IsExist(ctx context.Context, name string) (bool, error)
	Create(ctx context.Context, inst core.InstitutionModel) (core.InstitutionModel, error)
}

type InstitutionService struct {
	institutionRepo InstitutionRepo
}

func NewInstitutionService(institutionRepo InstitutionRepo) *InstitutionService {
	return &InstitutionService{institutionRepo: institutionRepo}
}

func (s InstitutionService) IsExist(ctx context.Context, name string) (bool, error) {
	return s.institutionRepo.IsExist(ctx, name)
}

func (s InstitutionService) Create(ctx context.Context, inst core.Institution) (core.InstitutionModel, error) {
	return s.institutionRepo.Create(ctx, core.InstitutionModel{
		Name: inst.Name,
	})
}
