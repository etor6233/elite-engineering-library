package main

import (
	"context"
	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	"elite.local/enterprise/internal/returnfiscal"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	url, id := os.Getenv("DATABASE_URL"), os.Getenv("RETURN_FISCAL_WORKER_ID")
	if url == "" || id == "" {
		os.Exit(2)
	}
	cfg, e := pgxpool.ParseConfig(url)
	if e != nil {
		os.Exit(2)
	}
	cfg.MaxConns = 4
	p, e := pgxpool.NewWithConfig(ctx, cfg)
	if e != nil {
		os.Exit(1)
	}
	defer p.Close()
	if e = p.Ping(ctx); e != nil {
		os.Exit(1)
	}
	processor, e := returnfiscal.NewProcessor(postgres.NewReturnFiscal(p), randomid.Generator{}, id, 2*time.Minute, 30*time.Second)
	if e != nil {
		os.Exit(2)
	}
	if e = run(ctx, processor); e != nil && !errors.Is(e, context.Canceled) {
		os.Exit(1)
	}
}
func run(ctx context.Context, p *returnfiscal.Processor) error {
	timer := time.NewTimer(0)
	if !timer.Stop() {
		<-timer.C
	}
	defer timer.Stop()
	for {
		_, e := p.ProcessOne(ctx)
		if e == nil {
			continue
		}
		if errors.Is(e, context.Canceled) {
			return e
		}
		wait := 2 * time.Second
		if errors.Is(e, returnfiscal.ErrNoWork) {
			wait = 500 * time.Millisecond
		} else if !errors.Is(e, returnfiscal.ErrPending) {
			slog.Warn("return fiscal deferred")
		}
		timer.Reset(wait)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
		}
	}
}
