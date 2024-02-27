package repo2

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sollniss/repository-patterns/entity"
)

// Repo only contains bare bone repository functionality.
type Repo struct {
	baseRepo
}

func New(dsn string) (*Repo, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("could not open connection: %w", err)
	}
	return &Repo{
		baseRepo: newBaseRepo(db),
	}, nil
}

// UpdateName has to be implemented differently compared to repo1
// since we have no way to access our other functions atomically.
//
// TODO: add queries
func (r Repo) UpdateName(ctx context.Context, id entity.ID, name string) (err error) {
	tx, err := r.baseRepo.beginTx(ctx)
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

	// Manual update.
	// ...

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (q Repo) Get(ctx context.Context, id entity.ID) (entity.Entity, error) {
	return entity.Entity{}, nil
}

func (q Repo) Update(ctx context.Context, e entity.Entity) error {
	return nil
}
