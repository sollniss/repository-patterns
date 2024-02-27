package main

import (
	"context"

	"github.com/sollniss/repository-patterns/entity"
)

// Repo1 hides transactions. Implementation is responsible for creating atomic methods.
type Repo1 interface {
	Get(ctx context.Context, id entity.ID) (entity.Entity, error)
	Update(ctx context.Context, e entity.Entity) error
	UpdateName(ctx context.Context, id entity.ID, name string) error
}

func Work1(ctx context.Context, r Repo1) {
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

	// no external transaction
}
