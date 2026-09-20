package main

import (
	"context"
	"log"

	"github.com/jjmrocha/warren/internal/config"
	"github.com/jjmrocha/warren/internal/engine"
	"github.com/jjmrocha/warren/internal/setup"
)

func main() {
	if err := setup.BuildIfNeed(); err != nil {
		log.Fatal(err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	if err := engine.Run(ctx, cfg); err != nil {
		log.Fatal(err)
	}
}
