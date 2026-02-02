package svc

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxExecutor struct {
	pool *pgxpool.Pool
}

// Must match: Exec(ctx context.Context, sql string, args ...any) error
func (e *PgxExecutor) Exec(ctx context.Context, sql string, args ...any) error {
	_, err := e.pool.Exec(ctx, sql, args...)
	return err
}

// Must match: Ping(ctx context.Context) error
func (e *PgxExecutor) Ping(ctx context.Context) error {
	return e.pool.Ping(ctx)
}

// Must match: Close() error
func (e *PgxExecutor) Close() error {
	e.pool.Close()
	return nil
}

func NewPgxExecutor(dsn string) (*PgxExecutor, error) {
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	return &PgxExecutor{pool: pool}, nil
}
