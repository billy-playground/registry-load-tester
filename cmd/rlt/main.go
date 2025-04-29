package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/billy-playground/registry-load-tester/cmd/rlt/root"
)

func main() {
	if err := func() error {
		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
		defer cancel()
		return root.New().ExecuteContext(ctx)
	}(); err != nil {
		os.Exit(1)
	}
}
