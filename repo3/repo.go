package repo3

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sollniss/repository-patterns/entity"
)

func New(dsn string) (Repo, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("could not open connection: %w", err)
	}
	return &repo{
		conn: db,
		repoQueries: repoQueries{
			db: db,
		},
	}, nil
}

// repo only contains composite queries. Atomic queries are embedded via repoQueries.
type repo struct {
	conn *sql.DB
	repoQueries
}

// repoTx holds atomic transactions.
// Has to be replicated for every repository.
type repoTx struct {
	tx *sql.Tx
	repoQueries
}

func (tx repoTx) Commit() error {
	return tx.tx.Commit()
}

func (tx repoTx) Rollback() error {
	return tx.tx.Rollback()
}

// BeginTx has to be replicated for every repository.
func (r repo) BeginTx(ctx context.Context) (RepoTx, error) {
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", err)
	}
	return &repoTx{
		tx: tx,
		repoQueries: repoQueries{
			db: tx,
		},
	}, nil
}

// UpdateName is just for demonstration for a non-atomic function.
// Same as in repo1.
func (r repo) UpdateName(ctx context.Context, id entity.ID, name string) (err error) {
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				err = errors.Join(err, rollbackErr)
			}
		}
	}()

	// access our repoQueries.
	queries := repoQueries{tx}

	entity, err := queries.Get(ctx, id)
	if err != nil {
		return err
	}

	entity.Name = name

	err = queries.Update(ctx, entity)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// repoQueries contains atomic queries only.
// Can be split up into reader and writer queries with CQRS.
type repoQueries struct {
	db dbtx
}

func (q repoQueries) Get(ctx context.Context, id entity.ID) (entity.Entity, error) {
	return entity.Entity{}, nil
}

func (q repoQueries) Update(ctx context.Context, e entity.Entity) error {
	return nil
}
