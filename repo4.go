package main

import (
	"context"

	"github.com/sollniss/repository-patterns/entity"
)

// Repo4 provides functionality that allows modification of the entity between get and update.
// Kind of a mixture between Repo1 and Repo2.
type Repo4 interface {
	Get(ctx context.Context, id entity.ID) (entity.Entity, error)
	Update(ctx context.Context, e entity.Entity) error
	UpdateName(ctx context.Context, id entity.ID, name string) error
	GetAndUpdate(ctx context.Context, id entity.ID, updateFunc func(e *entity.Entity) error) error
}

func Work4(ctx context.Context, r Repo4) {
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
	err = r.GetAndUpdate(ctx, 1, func(e *entity.Entity) error {
		e.Name = "test"
		return nil
	})
	if err != nil {
		panic(err)
	}
}
