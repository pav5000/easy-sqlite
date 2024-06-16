//nolint:wrapcheck
package easysqlite

import (
	"context"
	"database/sql"

	"github.com/pav5000/easy-sqlite/internal/txctx"
)

func (s *EasySqlite) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if tx := txctx.Extract(ctx); tx != nil {
		return tx.ExecContext(ctx, query, args...)
	}

	return s.db.ExecContext(ctx, query, args...)
}

func (s *EasySqlite) GetContext(ctx context.Context, dest any, query string, args ...any) error {
	if tx := txctx.Extract(ctx); tx != nil {
		return tx.GetContext(ctx, dest, query, args...)
	}

	return s.db.GetContext(ctx, dest, query, args...)
}

func (s *EasySqlite) SelectContext(ctx context.Context, dest any, query string, args ...any) error {
	if tx := txctx.Extract(ctx); tx != nil {
		return tx.SelectContext(ctx, dest, query, args...)
	}

	return s.db.SelectContext(ctx, dest, query, args...)
}

func (s *EasySqlite) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	if tx := txctx.Extract(ctx); tx != nil {
		return tx.PrepareContext(ctx, query)
	}

	return s.db.PrepareContext(ctx, query)
}

func (s *EasySqlite) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if tx := txctx.Extract(ctx); tx != nil {
		return tx.QueryContext(ctx, query, args...)
	}

	return s.db.QueryContext(ctx, query, args...)
}

func (s *EasySqlite) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	if tx := txctx.Extract(ctx); tx != nil {
		return tx.QueryRowContext(ctx, query, args...)
	}

	return s.db.QueryRowContext(ctx, query, args...)
}

func (s *EasySqlite) MustExecContext(ctx context.Context, query string, args ...any) sql.Result {
	res, err := s.ExecContext(ctx, query, args...)
	if err != nil {
		panic(err)
	}

	return res
}
