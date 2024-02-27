package repo1

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sollniss/repository-patterns/entity"
)

// Repo only contains transactional queries, but also embeds all atomic queries.
type Repo struct {
	conn *sql.DB
	repoQueries
}

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

// UpdateName handles the update atomically inside the repository.
// Could also directly use an update query in case of SQL.
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

	// Access our repoQueries.
	queries := repoQueries{tx}

	entity, err := queries.Get(ctx, id)

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
