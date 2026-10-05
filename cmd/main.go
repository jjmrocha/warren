package main

import (
	"context"
	"log"
	"os"

	"github.com/jjmrocha/warren/internal/cli"
	"github.com/jjmrocha/warren/internal/config"
	"github.com/jjmrocha/warren/internal/engine"
	"github.com/jjmrocha/warren/internal/setup"
)

func main() {
	args := cli.Parse(os.Args[1:])

	if err := setup.BuildIfNeed(); err != nil {
		log.Fatal(err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	if err := engine.Run(ctx, cfg, args.SessionID); err != nil {
		log.Fatal(err)
	}
}
