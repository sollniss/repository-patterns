package main

import (
	"context"

	"github.com/sollniss/repository-patterns/repo1"
	"github.com/sollniss/repository-patterns/repo2"
	"github.com/sollniss/repository-patterns/repo3"
	"github.com/sollniss/repository-patterns/repo4"
)

func main() {
	ctx := context.Background()

	r1, _ := repo1.New("")
	Work1(ctx, r1)

	r2, _ := repo2.New("")
	Work2(ctx, r2)

	r3, _ := repo3.New("")
	Work3(ctx, r3)

	r4, _ := repo4.New("")
	Work4(ctx, r4)
}
