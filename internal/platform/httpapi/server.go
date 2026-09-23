package httpapi

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"elite.local/enterprise/internal/order"
	"elite.local/enterprise/internal/platform/identity"
)

type Server struct {
	orders   *order.Service
	verifier identity.Verifier
}

func New(orders *order.Service, verifier identity.Verifier) http.Handler {
	s := &Server{orders: orders, verifier: verifier}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health/live", s.live)
	mux.HandleFunc("POST /v1/orders", s.createOrder)
	return recoverMiddleware(mux)
}

func (s *Server) live(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) createOrder(w http.ResponseWriter, r *http.Request) {
	principal, err := authenticate(r.Context(), r.Header.Get("Authorization"), s.verifier)
	if err != nil {
		writeProblem(w, 401, "UNAUTHENTICATED", "a valid bearer token is required")
		return
	}
	if !principal.Allowed("order:create") {
		writeProblem(w, 403, "FORBIDDEN", "order:create permission is required")
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		writeProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE", "Content-Type must be application/json")
		return
	}
	key := r.Header.Get("Idempotency-Key")
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
	if err != nil {
		writeProblem(w, 400, "INVALID_BODY", "body is invalid or too large")
		return
	}
	var input struct {
		OrganizationID, Currency string
		TotalMinorUnits          int64
	}
	decoder := json.NewDecoder(strings.NewReader(string(body)))
	decoder.DisallowUnknownFields()
	if err = decoder.Decode(&input); err != nil {
		writeProblem(w, 400, "INVALID_BODY", "body does not match the contract")
		return
	}
	if !principal.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	hash := sha256.Sum256(body)
	entity, replayed, err := s.orders.Create(r.Context(), order.Create{TenantID: principal.TenantID, OrganizationID: input.OrganizationID, CustomerPrincipal: principal.Subject, Currency: input.Currency, TotalMinorUnits: input.TotalMinorUnits, IdempotencyKey: key, RequestHash: hex.EncodeToString(hash[:])})
	if err != nil {
		if errors.Is(err, order.ErrConflict) {
			writeProblem(w, 409, "ORDER_CONFLICT", err.Error())
			return
		}
		writeProblem(w, 500, "INTERNAL_ERROR", "request failed")
		return
	}
	if replayed {
		writeJSON(w, 200, entity)
	} else {
		writeJSON(w, 201, entity)
	}
}

func authenticate(ctx context.Context, authorization string, verifier identity.Verifier) (identity.Principal, error) {
	const prefix = "Bearer "
	if !strings.HasPrefix(authorization, prefix) || len(authorization) == len(prefix) {
		return identity.Principal{}, identity.ErrUnauthenticated
	}
	return verifier.Verify(ctx, strings.TrimSpace(strings.TrimPrefix(authorization, prefix)))
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func writeProblem(w http.ResponseWriter, status int, code, detail string) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{"type": "urn:elite:problem:" + code, "status": status, "code": code, "detail": detail})
}
func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recover() != nil {
				writeProblem(w, 500, "INTERNAL_ERROR", "request failed")
			}
		}()
		ctx, cancel := contextWithTimeout(r, 10*time.Second)
		defer cancel()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ServeUntilShutdown stops admission on cancellation and waits for active HTTP
// requests to drain before the caller can close shared resources. A drain timeout
// force-closes connections and is returned to the caller. Hijacked connections
// require a separate owner; net/http Shutdown and Close do not wait for them.
func ServeUntilShutdown(ctx context.Context, server *http.Server, timeout time.Duration) error {
	if server == nil {
		return errors.New("HTTP server is required")
	}
	return serveUntilShutdown(ctx, server, timeout, server.ListenAndServe)
}
func serveUntilShutdown(ctx context.Context, server *http.Server, timeout time.Duration, serve func() error) error {
	if ctx == nil || server == nil || timeout <= 0 || serve == nil {
		return errors.New("HTTP shutdown requires context, server, positive timeout and serve callback")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	done := make(chan error, 1)
	go func() { done <- serve() }()
	var serveErr error
	serveFinished := false
	select {
	case <-ctx.Done():
	case serveErr = <-done:
		serveFinished = true
	}
	// ListenAndServe/Serve returns as soon as listeners close, before requests
	// drain. Shutdown must therefore be awaited on this caller's path, including
	// when an unexpected listener failure initiated the shutdown.
	shutdown, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	shutdownErr := server.Shutdown(shutdown)
	var closeErr error
	if shutdownErr != nil {
		closeErr = server.Close()
	}
	if !serveFinished {
		serveErr = <-done
	}
	if ctx.Err() != nil && errors.Is(serveErr, http.ErrServerClosed) {
		serveErr = nil
	}
	return errors.Join(serveErr, shutdownErr, closeErr)
}
