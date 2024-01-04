package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/migmatore/study-platform-api/internal/repository/psql"
)

type TransactionRepo struct {
	pool psql.AtomicPoolClient
}

func NewTransactor(pool psql.AtomicPoolClient) *TransactionRepo {
	return &TransactionRepo{pool: pool}
}

func (s *TransactionRepo) WithinTransaction(ctx context.Context, txFunc func(context context.Context) error) error {
	return s.pool.BeginTxFunc(ctx, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return txFunc(psql.InjectTx(ctx, tx))
	})
}
