package main

import (
	"context"

	"github.com/sollniss/repository-patterns/entity"
	"github.com/sollniss/repository-patterns/repo2"
)

// Repo2 contains a packet level function that executes arbitrary code in a transaction.
// Does automatic committing and rollbacks.
// Requires repo2.Txer in all atomic repository interfaces.
type Repo2 interface {
	Get(ctx context.Context, id entity.ID) (entity.Entity, error)
	Update(ctx context.Context, e entity.Entity) error
	UpdateName(ctx context.Context, id entity.ID, name string) error
	repo2.Txer
}

func Work2(ctx context.Context, r Repo2) {
	// manual 1
	e, err := r.Get(ctx, 1)
	if err != nil {
		panic(err)
	}
	e.Name = "test"
	err = r.Update(ctx, e)
	if err != nil {
		panic(err)
	}

	// internal transaction
	err = r.UpdateName(ctx, 1, "test")
	if err != nil {
		panic(err)
	}

	// external transaction
	err = repo2.Transaction(ctx, r, func(tx Repo2) error {
		e, err := tx.Get(ctx, 1)
		if err != nil {
			return err
		}

		e.Name = "test"

		err = tx.Update(ctx, e)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}
