package psql

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/migmatore/study-platform-api/config"
	"github.com/migmatore/study-platform-api/pkg/logger"
	"github.com/migmatore/study-platform-api/pkg/utils"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type AtomicPoolClient interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	BeginTxFunc(ctx context.Context, txOptions pgx.TxOptions, f func(pgx.Tx) error) error
	Ping(ctx context.Context) error
	Close()
}

type AtomicPool struct {
	pool *pgxpool.Pool
}

func NewPostgres(ctx context.Context, maxAttempts int, cfg *config.Config, logger logger.Logger) (*AtomicPool, error) {
	var err error
	var pool *pgxpool.Pool

	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Postgres.Username,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DBName,
	)

	err = utils.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		pool, err = pgxpool.Connect(ctx, dsn)
		if err != nil {
			logger.Errorf("DB connection error. %v", err)
			return err
		}

		if err := pool.Ping(ctx); err != nil {
			logger.Errorf("DB ping error. %v\n", err)
			return err
		}

		return nil
	}, maxAttempts, 5*time.Second)

	return &AtomicPool{pool: pool}, err
}

func (p *AtomicPool) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	if tx := ExtractTx(ctx); tx != nil {
		return tx.Exec(ctx, sql, arguments...)
	}

	return p.pool.Exec(ctx, sql, arguments...)
}

func (p *AtomicPool) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	if tx := ExtractTx(ctx); tx != nil {
		return tx.Query(ctx, sql, args...)
	}

	return p.pool.Query(ctx, sql, args...)
}

func (p *AtomicPool) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	if tx := ExtractTx(ctx); tx != nil {
		return tx.QueryRow(ctx, sql, args...)
	}

	return p.pool.QueryRow(ctx, sql, args...)
}

func (p *AtomicPool) BeginTxFunc(ctx context.Context, txOptions pgx.TxOptions, f func(pgx.Tx) error) error {
	return p.pool.BeginTxFunc(ctx, txOptions, f)
}

func (p *AtomicPool) Ping(ctx context.Context) error {
	return p.pool.Ping(ctx)
}

func (p *AtomicPool) Close() {
	p.pool.Close()
}

// Reconnect Auto reconnecting to db
func (p *AtomicPool) Reconnect(ctx context.Context, cfg *config.Config, logger logger.Logger) {
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		cfg.Postgres.Username,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DBName,
	)

	for {
		if err := p.Ping(ctx); err != nil {
			p.Close()
			//ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
			//defer cancel()
			if p != nil {
				pool, err := pgxpool.Connect(ctx, dsn)
				if err != nil {
					logger.Errorf("DB reconnection error. %v", err)
					time.Sleep(1 * time.Second)

					continue
				}
				p.pool = pool

			}

		}

		time.Sleep(1 * time.Second)
	}
}
