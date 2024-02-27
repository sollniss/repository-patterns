package repo4

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sollniss/repository-patterns/entity"
)

func New(dsn string) (*Repo, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("could not open connection: %w", err)
	}
	return &Repo{
		conn: db,
		repoQueries: repoQueries{
			db: db,
		},
	}, nil
}

type Repo struct {
	conn *sql.DB
	repoQueries
}

// GetAndUpdate allows updateFunc to modify the entity after reading
// while also handling automatic commit and rollback on error.
func (r Repo) GetAndUpdate(ctx context.Context, id entity.ID, updateFunc func(e *entity.Entity) error) (err error) {
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
	q := repoQueries{tx}

	// Get can make use of `select ... for update` in SQL.
	entity, err := q.Get(ctx, id)
	if err != nil {
		return err
	}

	err = updateFunc(&entity)
	if err != nil {
		return err
	}

	err = q.Update(ctx, entity)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// UpdateName is just for demonstration for a non-atomic function.
// Same as in repo1.
func (r Repo) UpdateName(ctx context.Context, id entity.ID, name string) (err error) {
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
