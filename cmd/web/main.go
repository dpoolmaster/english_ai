package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/english-ai/english_ai/internal/platform/config"
	"github.com/english-ai/english_ai/internal/platform/httpserver"
	"github.com/english-ai/english_ai/internal/platform/httpx"
	"github.com/english-ai/english_ai/internal/platform/logging"
	"github.com/english-ai/english_ai/internal/platform/metrics"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "web:", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("configuration: %w", err)
	}

	logger := logging.New(os.Stdout, cfg.LogLevel)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	router := httpx.NewRouter(logger.With("process", "web"))
	metricsHandler := metrics.New().Mux()

	logger.Info("web process starting", "env", cfg.Env, "base_url", cfg.BaseURL)

	group, groupCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		return httpserver.New(cfg.HTTPAddr, router, logger.With("process", "web")).Run(groupCtx, cfg.ShutdownTimeout)
	})
	group.Go(func() error {
		return httpserver.New(cfg.MetricsAddr, metricsHandler, logger.With("process", "metrics")).Run(groupCtx, cfg.ShutdownTimeout)
	})

	return group.Wait()
}
