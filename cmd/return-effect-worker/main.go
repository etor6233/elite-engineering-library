package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	"elite.local/enterprise/internal/returneffects"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	databaseURL := os.Getenv("DATABASE_URL")
	workerID := os.Getenv("RETURN_EFFECT_WORKER_ID")
	if databaseURL == "" || workerID == "" {
		slog.Error("DATABASE_URL and RETURN_EFFECT_WORKER_ID are required")
		os.Exit(2)
	}
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		slog.Error("database configuration is invalid")
		os.Exit(2)
	}
	config.MaxConns = 4
	config.MinConns = 0
	config.MaxConnLifetime = 30 * time.Minute
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		slog.Error("database pool creation failed")
		os.Exit(1)
	}
	defer pool.Close()
	if err = pool.Ping(ctx); err != nil {
		slog.Error("database unavailable")
		os.Exit(1)
	}
	processor, err := returneffects.NewInventoryProcessor(postgres.NewReturnEffects(pool), randomid.Generator{}, workerID, 2*time.Minute, 5*time.Second)
	if err != nil {
		slog.Error("return effect worker configuration is invalid")
		os.Exit(2)
	}
	if err = run(ctx, processor); err != nil && !errors.Is(err, context.Canceled) {
		slog.Error("return effect worker stopped unexpectedly")
		os.Exit(1)
	}
}

func run(ctx context.Context, processor *returneffects.Processor) error {
	timer := time.NewTimer(0)
	if !timer.Stop() {
		<-timer.C
	}
	defer timer.Stop()
	for {
		_, err := processor.ProcessOne(ctx)
		if err == nil {
			continue
		}
		if errors.Is(err, context.Canceled) {
			return err
		}
		wait := 2 * time.Second
		if errors.Is(err, returneffects.ErrNoWork) {
			wait = 500 * time.Millisecond
		} else {
			slog.Warn("return inventory effect deferred")
		}
		timer.Reset(wait)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
		}
	}
}
