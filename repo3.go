package main

import (
	"context"

	"github.com/sollniss/repository-patterns/repo3"
)

// Repo3 exports it's own interface for the repository.
// The repository includes all transaction functionality.
type Repo3 = repo3.Repo

func Work3(ctx context.Context, r Repo3) {
	// manual
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
	tx, _ := r.BeginTx(ctx)
	e, _ = tx.Get(ctx, 1)
	e.Name = "test"
	err = tx.Update(ctx, e)
	if err != nil {
		panic(err)
	}
}
