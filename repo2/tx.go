package repo2

import (
	"context"
	"database/sql"
)

// Interface assertion.
var _ Txer = (*baseRepo)(nil)

// newBaseRepo returns a new baseRepo.
// This is just a helper for cleaner code at the actual repository implementation site.
func newBaseRepo(db *sql.DB) baseRepo {
	return baseRepo{
		conn: db,
		db:   db,
	}
}

// baseRepo needs to be embedded in all repositories.
type baseRepo struct {
	conn *sql.DB
	db   dbtx
}

func (r *baseRepo) beginTx(ctx context.Context) (*sql.Tx, error) {
	return r.conn.BeginTx(ctx, nil)
}

func (r *baseRepo) withTx(tx *sql.Tx) {
	r.db = tx
}

type dbtx interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

// Txer allows external interfaces to use Transaction,
// while also limiting access to the transaction itself.
type Txer interface {
	beginTx(ctx context.Context) (*sql.Tx, error)
	withTx(*sql.Tx)
}

// Transaction is a top level function that executes arbitrary code in an transaction.
func Transaction[T Txer](ctx context.Context, repo T, f func(tx T) error) error {
	tx, err := repo.beginTx(ctx)
	if err != nil {
		return err
	}
	var txRepo T
	txRepo.withTx(tx)

	if err := f(txRepo); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
