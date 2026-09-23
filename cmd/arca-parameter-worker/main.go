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

var parameterWorkerPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,63}$`)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	databaseURL, socketPath, workerID := os.Getenv("DATABASE_URL"), os.Getenv("ARCA_WSFE_SOCKET"), os.Getenv("FISCAL_PARAMETER_WORKER_ID")
	if databaseURL == "" || socketPath == "" || !parameterWorkerPattern.MatchString(workerID) {
		slog.Error("DATABASE_URL, ARCA_WSFE_SOCKET and a valid FISCAL_PARAMETER_WORKER_ID are required")
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
	repository, ids := postgres.NewFiscal(pool), randomid.Generator{}
	registry := fiscal.NewParameterRegistry(repository, provider, ids, time.Now)
	processor, err := fiscal.NewParameterRefreshProcessor(repository, registry, ids, workerID, 2*time.Minute, 5*time.Minute)
	if err != nil {
		slog.Error("parameter processor configuration is invalid")
		os.Exit(2)
	}
	if err = runParameterWorker(ctx, processor); err != nil && !errors.Is(err, context.Canceled) {
		slog.Error("parameter worker stopped unexpectedly")
		os.Exit(1)
	}
}

func runParameterWorker(ctx context.Context, processor *fiscal.ParameterRefreshProcessor) error {
	timer := time.NewTimer(0)
	if !timer.Stop() {
		<-timer.C
	}
	defer timer.Stop()
	for {
		err := processor.ProcessOne(ctx)
		if err == nil {
			continue
		}
		if errors.Is(err, context.Canceled) {
			return err
		}
		wait := 2 * time.Second
		if errors.Is(err, fiscal.ErrNoParameterWork) {
			wait = 30 * time.Second
		} else {
			slog.Warn("parameter refresh deferred after recoverable failure")
		}
		timer.Reset(wait)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
		}
	}
}
