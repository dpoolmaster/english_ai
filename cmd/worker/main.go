package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/english-ai/english_ai/internal/platform/config"
	"github.com/english-ai/english_ai/internal/platform/logging"
	"github.com/english-ai/english_ai/internal/platform/worker"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "worker:", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("configuration: %w", err)
	}

	logger := logging.New(os.Stdout, cfg.LogLevel).With("process", "worker")
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger.Info("worker process starting", "env", cfg.Env)
	return worker.New(logger).Run(ctx)
}
