package main

import (
	"context"
	"encoding/hex"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/socialbridge"
	"github.com/jackc/pgx/v5/pgxpool"
)

func run(ctx context.Context) error {
	p, e := socialbridge.LoadFile(os.Getenv("SOCIAL_PROFILE_FILE"), os.Getenv("SOCIAL_PROFILE_SHA256"))
	if e != nil {
		return socialbridge.ErrBinding
	}
	key, e := hex.DecodeString(os.Getenv("SOCIAL_FENCE_HMAC_KEY_HEX"))
	if e != nil || len(key) < 32 {
		return socialbridge.ErrBinding
	}
	adapter := socialbridge.Process{Python: os.Getenv("SOCIAL_PYTHON"), Script: os.Getenv("SOCIAL_BRIDGE_FILE"), ScriptSHA256: os.Getenv("SOCIAL_BRIDGE_SHA256")}
	for _, k := range []string{"META_APP_ID", "META_APP_SECRET", "META_PAGE_ACCESS_TOKEN"} {
		v := os.Getenv(k)
		if v == "" {
			return socialbridge.ErrBinding
		}
		adapter.Environment = append(adapter.Environment, k+"="+v)
	}
	if adapter.Validate() != nil || os.Getenv("DATABASE_URL") == "" || os.Getenv("OIDC_ISSUER") == "" || os.Getenv("OIDC_AUDIENCE") == "" {
		return socialbridge.ErrBinding
	}
	if adapter.Preflight(ctx) != nil {
		return socialbridge.ErrBinding
	}
	verifier, e := identity.NewOIDCVerifier(ctx, os.Getenv("OIDC_ISSUER"), os.Getenv("OIDC_AUDIENCE"))
	if e != nil {
		return e
	}
	pool, e := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if e != nil {
		return e
	}
	defer pool.Close()
	if e = pool.Ping(ctx); e != nil {
		return e
	}
	service, e := postgres.NewSocialPublishing(pool, p, key, adapter)
	if e != nil {
		return e
	}
	workerCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	workerDone := make(chan struct{})
	go func() {
		defer close(workerDone)
		lastState := ""
		report := func(state string) {
			if state != lastState {
				slog.Info("social worker state", "state", state)
				lastState = state
			}
		}
		ticker := time.NewTicker(time.Duration(p.Config().PollSeconds) * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-workerCtx.Done():
				return
			case <-ticker.C:
				jobs, e := service.Claim(workerCtx, "social-publishing")
				if e != nil {
					report("QUEUE_UNAVAILABLE")
					continue
				}
				if len(jobs) == 0 {
					report("IDLE")
				}
				for _, job := range jobs {
					e := service.Process(workerCtx, job, "social-publishing")
					if e == nil {
						report("ACKNOWLEDGED")
					} else if errors.Is(e, socialbridge.ErrUnknown) {
						report("PROVIDER_UNRESOLVED")
					} else {
						report("DISPATCH_NOT_ADMITTED")
					}
				}
			}
		}
	}()
	defer func() { cancel(); <-workerDone }()
	addr := os.Getenv("SOCIAL_HTTP_ADDRESS")
	if addr == "" {
		addr = "127.0.0.1:8097"
	}
	mux := http.NewServeMux()
	mux.Handle("/", httpapi.NewSocialPublishing(service, verifier))
	mux.HandleFunc("GET /health/live", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) })
	mux.HandleFunc("GET /health/ready", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if pool.Ping(ctx) != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
	server := &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 45 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 16384}
	return httpapi.ServeUntilShutdown(ctx, server, 40*time.Second)
}
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if e := run(ctx); e != nil && !errors.Is(e, context.Canceled) {
		slog.Error("social publishing configuration or runtime unavailable")
		os.Exit(2)
	}
}
