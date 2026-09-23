package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"elite.local/enterprise/internal/order"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ids struct{}

func (ids) New() string {
	var value [16]byte
	if _, err := rand.Read(value[:]); err != nil {
		panic(err)
	}
	return hex.EncodeToString(value[:])
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		slog.Error("DATABASE_URL is required")
		os.Exit(2)
	}
	issuer, audience := os.Getenv("OIDC_ISSUER"), os.Getenv("OIDC_AUDIENCE")
	if issuer == "" || audience == "" {
		slog.Error("OIDC_ISSUER and OIDC_AUDIENCE are required")
		os.Exit(2)
	}
	verifier, err := identity.NewOIDCVerifier(ctx, issuer, audience)
	if err != nil {
		slog.Error("OIDC discovery failed")
		os.Exit(1)
	}
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		slog.Error("database configuration failed")
		os.Exit(2)
	}
	defer pool.Close()
	if err = pool.Ping(ctx); err != nil {
		slog.Error("database unavailable")
		os.Exit(1)
	}
	handler := httpapi.New(order.NewService(postgres.NewOrders(pool), ids{}), verifier)
	server := &http.Server{Addr: ":8080", Handler: handler, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 1 << 20}
	slog.Info("api starting", "address", server.Addr)
	if err = httpapi.ServeUntilShutdown(ctx, server, 15*time.Second); err != nil && !errors.Is(err, context.Canceled) {
		slog.Error("server failed")
		os.Exit(1)
	}
}
