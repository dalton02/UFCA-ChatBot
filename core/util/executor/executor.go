package executor

import (
	"context"
	"database/sql"
	"licor_model/core/server/shared"
)

type Executor interface {
	// Métodos COM Context (para operações assíncronas/timeouts)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)

	// Métodos SEM Context (para operações síncronas simples)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Prepare(query string) (*sql.Stmt, error)
	GetDB() *sql.DB // ✅ Expor de forma controlada

}

// DBExecutor implementação concreta
type DBExecutor struct {
	db *sql.DB
	tx *sql.Tx
}

func NewDBExecutor(tx *sql.Tx) Executor {
	return &DBExecutor{db: shared.DB, tx: tx}
}

func (e *DBExecutor) GetDB() *sql.DB {
	return e.db
}

// Métodos COM Context
func (e *DBExecutor) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if e.tx != nil {
		return e.tx.ExecContext(ctx, query, args...)
	}
	return e.db.ExecContext(ctx, query, args...)
}

func (e *DBExecutor) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if e.tx != nil {
		return e.tx.QueryContext(ctx, query, args...)
	}
	return e.db.QueryContext(ctx, query, args...)
}

func (e *DBExecutor) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	if e.tx != nil {
		return e.tx.QueryRowContext(ctx, query, args...)
	}
	return e.db.QueryRowContext(ctx, query, args...)
}

func (e *DBExecutor) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	if e.tx != nil {
		return e.tx.PrepareContext(ctx, query)
	}
	return e.db.PrepareContext(ctx, query)
}

// Métodos SEM Context (usam context.Background())
func (e *DBExecutor) Exec(query string, args ...interface{}) (sql.Result, error) {
	return e.ExecContext(context.Background(), query, args...)
}

func (e *DBExecutor) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return e.QueryContext(context.Background(), query, args...)
}

func (e *DBExecutor) QueryRow(query string, args ...interface{}) *sql.Row {
	return e.QueryRowContext(context.Background(), query, args...)
}

func (e *DBExecutor) Prepare(query string) (*sql.Stmt, error) {
	return e.PrepareContext(context.Background(), query)
}
