package psql

import (
	"context"
	"github.com/jackc/pgx/v4"
)

type txKey struct{}

// InjectTx injects transaction to context
func InjectTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// ExtractTx extracts transaction from context
func ExtractTx(ctx context.Context) pgx.Tx {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return tx
	}
	return nil
}
