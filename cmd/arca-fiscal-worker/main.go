package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"elite.local/enterprise/internal/fiscal"
	"elite.local/enterprise/internal/fiscal/wsfeipc"
	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var workerIDPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,63}$`)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	databaseURL := os.Getenv("DATABASE_URL")
	socketPath := os.Getenv("ARCA_WSFE_SOCKET")
	workerID := os.Getenv("FISCAL_WORKER_ID")
	if databaseURL == "" || socketPath == "" || !workerIDPattern.MatchString(workerID) {
		slog.Error("DATABASE_URL, ARCA_WSFE_SOCKET and a valid FISCAL_WORKER_ID are required")
		os.Exit(2)
	}

	poolConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		slog.Error("database configuration is invalid")
		os.Exit(2)
	}
	poolConfig.MaxConns = 4
	poolConfig.MinConns = 0
	poolConfig.MaxConnLifetime = 30 * time.Minute
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		slog.Error("database pool creation failed")
		os.Exit(1)
	}
	defer pool.Close()
	if err = pool.Ping(ctx); err != nil {
		slog.Error("database unavailable")
		os.Exit(1)
	}

	provider, err := wsfeipc.NewUnixProvider(socketPath, 30*time.Second)
	if err != nil {
		slog.Error("ARCA WSFE socket configuration is invalid")
		os.Exit(2)
	}
	processor, err := fiscal.NewProcessor(postgres.NewFiscal(pool), provider, randomid.Generator{}, workerID, 2*time.Minute)
	if err != nil {
		slog.Error("fiscal processor configuration is invalid")
		os.Exit(2)
	}

	if err = run(ctx, processor); err != nil && !errors.Is(err, context.Canceled) {
		slog.Error("fiscal worker stopped unexpectedly")
		os.Exit(1)
	}
}

func run(ctx context.Context, processor *fiscal.Processor) error {
	idle := time.NewTimer(0)
	if !idle.Stop() {
		<-idle.C
	}
	defer idle.Stop()
	for {
		_, err := processor.ProcessOne(ctx)
		if err == nil {
			continue
		}
		if errors.Is(err, context.Canceled) {
			return err
		}
		wait := 2 * time.Second
		if errors.Is(err, fiscal.ErrNoWork) {
			wait = 500 * time.Millisecond
		} else {
			slog.Warn("fiscal work deferred after a recoverable processing failure")
		}
		idle.Reset(wait)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-idle.C:
		}
	}
}
