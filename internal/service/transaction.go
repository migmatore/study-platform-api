package service

import "context"

type TransactionRepo interface {
	WithinTransaction(ctx context.Context, txFunc func(txCtx context.Context) error) error
}

type TransactionService struct {
	repo TransactionRepo
}

func NewTransactionService(repo TransactionRepo) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s TransactionService) WithinTransaction(ctx context.Context, txFunc func(txCtx context.Context) error) error {
	return s.repo.WithinTransaction(ctx, txFunc)
}
