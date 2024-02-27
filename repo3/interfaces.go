package repo3

import (
	"context"

	"github.com/sollniss/repository-patterns/entity"
)

// This pattern is only possible if we return interfaces from the New and BeginTx functions,
// which goes against the "accept interfaces, return structs" principle.

// Repo contains everything that can be called without starting a transaction.
type Repo interface {
	BeginTx(ctx context.Context) (RepoTx, error)

	UpdateName(ctx context.Context, id entity.ID, name string) error

	atomicRepo
}

// RepoTx represents a single transaction.
type RepoTx interface {
	Commit() error
	Rollback() error

	atomicRepo
}

// atomicRepo holds functionality that can be used with both Repo and RepoTx.
type atomicRepo interface {
	Get(ctx context.Context, id entity.ID) (entity.Entity, error)
	Update(ctx context.Context, e entity.Entity) error
}
