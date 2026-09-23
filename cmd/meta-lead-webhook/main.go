package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"elite.local/enterprise/internal/leadstream"
	"elite.local/enterprise/internal/platform/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	required := map[string]string{
		"DATABASE_URL": os.Getenv("DATABASE_URL"), "META_APP_SECRET": os.Getenv("META_APP_SECRET"),
		"META_VERIFY_TOKEN": os.Getenv("META_VERIFY_TOKEN"), "TENANT_ID": os.Getenv("TENANT_ID"),
		"ORGANIZATION_ID": os.Getenv("ORGANIZATION_ID"),
	}
	for name, value := range required {
		if value == "" {
			slog.Error("required configuration is missing", "name", name)
			os.Exit(2)
		}
	}
	pool, err := pgxpool.New(ctx, required["DATABASE_URL"])
	if err != nil {
		slog.Error("database configuration failed")
		os.Exit(2)
	}
	defer pool.Close()
	if err = pool.Ping(ctx); err != nil {
		slog.Error("database unavailable")
		os.Exit(1)
	}
	handler := leadstream.MetaLeadWebhookHandler{
		Store: postgres.NewMetaLeadWebhookStore(pool), TenantID: required["TENANT_ID"], OrganizationID: required["ORGANIZATION_ID"],
		AppSecret: required["META_APP_SECRET"], VerifyToken: required["META_VERIFY_TOKEN"],
	}
	mux := http.NewServeMux()
	mux.Handle("/webhooks/meta-leads", handler)
	address := os.Getenv("LISTEN_ADDR")
	if address == "" {
		address = ":8082"
	}
	server := &http.Server{Addr: address, Handler: mux, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 1 << 20}
	go func() {
		<-ctx.Done()
		shutdown, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdown)
	}()
	slog.Info("Meta lead webhook listening", "address", address)
	if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("Meta lead webhook failed")
		os.Exit(1)
	}
}
