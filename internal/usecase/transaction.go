package usecase

import "context"

type TransactionService interface {
	WithinTransaction(ctx context.Context, txFunc func(txCtx context.Context) error) error
}
