package httpapi

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/order"
	"elite.local/enterprise/internal/platform/identity"
)

type fixedIDs struct{}

func (fixedIDs) New() string { return "order-1" }

type repository struct{ created order.Order }

type verifier struct {
	principal identity.Principal
	err       error
}

func (v verifier) Verify(context.Context, string) (identity.Principal, error) {
	return v.principal, v.err
}

func allowedVerifier() verifier {
	return verifier{principal: identity.Principal{
		Subject: "customer-1", TenantID: "018f4d4a-7b36-7a21-8d10-2f4c54c28a01",
		Permissions:   map[string]struct{}{"order:create": {}},
		Organizations: map[string]struct{}{"org-1": {}},
	}}
}

func (r *repository) Create(_ context.Context, entity order.Order, _, _ string) (order.Order, bool, error) {
	r.created = entity
	return entity, false, nil
}
func (r *repository) Get(context.Context, string, string, string) (order.Order, error) {
	return order.Order{}, order.ErrNotFound
}
func (r *repository) SaveTransition(context.Context, order.Order, int64) error { return nil }

func TestCreateOrderContract(t *testing.T) {
	repo := &repository{}
	handler := New(order.NewService(repo, fixedIDs{}), allowedVerifier())
	body := `{"OrganizationID":"org-1","Currency":"USD","TotalMinorUnits":100}`
	request := httptest.NewRequest(http.MethodPost, "/v1/orders", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Idempotency-Key", "order-create-000001")
	request.Header.Set("Authorization", "Bearer verified-by-test-double")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", response.Code, response.Body.String())
	}
	if repo.created.State != order.Draft {
		t.Fatalf("unexpected state %s", repo.created.State)
	}
}

func TestProblemUsesProblemJSON(t *testing.T) {
	handler := New(order.NewService(&repository{}, fixedIDs{}), allowedVerifier())
	request := httptest.NewRequest(http.MethodPost, "/v1/orders", strings.NewReader(`{}`))
	request.Header.Set("Authorization", "Bearer verified-by-test-double")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("expected 415, got %d", response.Code)
	}
	if got := response.Header().Get("Content-Type"); got != "application/problem+json" {
		t.Fatalf("unexpected content type %q", got)
	}
}

func TestCreateOrderRequiresAuthentication(t *testing.T) {
	handler := New(order.NewService(&repository{}, fixedIDs{}), allowedVerifier())
	request := httptest.NewRequest(http.MethodPost, "/v1/orders", strings.NewReader(`{}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", response.Code)
	}
}

func TestCreateOrderRejectsUnauthorizedOrganization(t *testing.T) {
	repo := &repository{}
	handler := New(order.NewService(repo, fixedIDs{}), allowedVerifier())
	request := httptest.NewRequest(http.MethodPost, "/v1/orders", strings.NewReader(`{"OrganizationID":"org-other","Currency":"USD","TotalMinorUnits":100}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Idempotency-Key", "order-create-000002")
	request.Header.Set("Authorization", "Bearer verified-by-test-double")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d: %s", response.Code, response.Body.String())
	}
	if repo.created.ID != "" {
		t.Fatal("repository called for unauthorized organization")
	}
}

func TestHTTPShutdownDrainsActiveRequest(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	started, release := make(chan struct{}), make(chan struct{})
	var releaseOnce sync.Once
	defer releaseOnce.Do(func() { close(release) })
	server := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(started)
		<-release
		_, _ = io.WriteString(w, "completed")
	})}
	defer server.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	acceptClosed := make(chan struct{})
	go func() {
		done <- serveUntilShutdown(ctx, server, 2*time.Second, func() error { err := server.Serve(listener); close(acceptClosed); return err })
	}()
	type response struct {
		status int
		body   string
		err    error
	}
	received := make(chan response, 1)
	client := &http.Client{Timeout: 3 * time.Second}
	defer client.CloseIdleConnections()
	go func() {
		r, e := client.Get("http://" + listener.Addr().String())
		if e != nil {
			received <- response{err: e}
			return
		}
		defer r.Body.Close()
		body, e := io.ReadAll(r.Body)
		received <- response{r.StatusCode, string(body), e}
	}()
	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("request did not start")
	}
	cancel()
	select {
	case <-acceptClosed:
	case <-time.After(2 * time.Second):
		t.Fatal("accept loop did not stop")
	}
	connection, err := net.DialTimeout("tcp", listener.Addr().String(), 200*time.Millisecond)
	if err == nil {
		connection.Close()
		t.Error("accepted a connection while draining")
	}
	early := false
	var serveErr error
	select {
	case serveErr = <-done:
		early = true
		t.Error("host returned before the active request drained")
	case <-time.After(150 * time.Millisecond):
	}
	releaseOnce.Do(func() { close(release) })
	if !early {
		select {
		case serveErr = <-done:
		case <-time.After(2 * time.Second):
			t.Fatal("drain did not complete")
		}
	}
	if serveErr != nil {
		t.Error(serveErr)
	}
	select {
	case r := <-received:
		if r.err != nil || r.status != 200 || r.body != "completed" {
			t.Fatal("active response lost", r)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("response missing")
	}
	if !t.Failed() {
		t.Log("HTTP_GRACEFUL_DRAIN_PASS accept_closed=1 request_completed=1 host_waited=1")
	}
}

func TestHTTPShutdownDeadlineClosesActiveConnection(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	started, canceled := make(chan struct{}), make(chan struct{})
	server := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(started)
		<-r.Context().Done()
		close(canceled)
	})}
	defer server.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() {
		done <- serveUntilShutdown(ctx, server, 80*time.Millisecond, func() error { return server.Serve(listener) })
	}()
	client := &http.Client{Timeout: 2 * time.Second}
	defer client.CloseIdleConnections()
	requestDone := make(chan error, 1)
	go func() {
		r, e := client.Get("http://" + listener.Addr().String())
		if r != nil {
			r.Body.Close()
		}
		requestDone <- e
	}()
	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("request did not start")
	}
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("missing drain failure: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("shutdown exceeded bound")
	}
	select {
	case <-canceled:
	case <-time.After(time.Second):
		t.Fatal("active connection was not canceled")
	}
	select {
	case err := <-requestDone:
		if err == nil {
			t.Fatal("truncated request reported success")
		}
	case <-time.After(time.Second):
		t.Fatal("client did not terminate")
	}
	t.Log("HTTP_SHUTDOWN_DEADLINE_PASS failure_reported=1 active_connection_closed=1")
}

func TestHTTPShutdownStartupFailure(t *testing.T) {
	occupied, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer occupied.Close()
	server := &http.Server{Addr: occupied.Addr().String()}
	defer server.Close()
	done := make(chan error, 1)
	go func() { done <- ServeUntilShutdown(context.Background(), server, time.Second) }()
	select {
	case err := <-done:
		var op *net.OpError
		if !errors.As(err, &op) || op.Op != "listen" {
			t.Fatalf("startup cause lost: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("startup failure waits for cancellation")
	}
	t.Log("HTTP_STARTUP_FAILURE_PASS cause_preserved=1 returned_without_signal=1")
}

func TestHTTPShutdownListenerFailureDrainsRequest(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	started, release := make(chan struct{}), make(chan struct{})
	var releaseOnce sync.Once
	defer releaseOnce.Do(func() { close(release) })
	server := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(started)
		<-release
		_, _ = io.WriteString(w, "committed")
	})}
	defer server.Close()
	acceptClosed := make(chan struct{})
	done := make(chan error, 1)
	go func() {
		done <- serveUntilShutdown(context.Background(), server, 2*time.Second, func() error {
			err := server.Serve(listener)
			close(acceptClosed)
			return err
		})
	}()
	response := make(chan string, 1)
	client := &http.Client{Timeout: 3 * time.Second}
	defer client.CloseIdleConnections()
	go func() {
		r, err := client.Get("http://" + listener.Addr().String())
		if err != nil {
			response <- "failed"
			return
		}
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			response <- "failed"
			return
		}
		response <- string(body)
	}()
	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("request did not start")
	}
	if err := listener.Close(); err != nil {
		t.Fatal(err)
	}
	select {
	case <-acceptClosed:
	case <-time.After(time.Second):
		t.Fatal("listener failure not observed")
	}
	select {
	case err := <-done:
		t.Fatalf("returned before request drained: %v", err)
	case <-time.After(150 * time.Millisecond):
	}
	releaseOnce.Do(func() { close(release) })
	select {
	case err := <-done:
		if !errors.Is(err, net.ErrClosed) {
			t.Fatalf("listener cause lost: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("drain did not complete")
	}
	select {
	case body := <-response:
		if body != "committed" {
			t.Fatalf("active response lost: %s", body)
		}
	case <-time.After(time.Second):
		t.Fatal("response did not complete")
	}
	t.Log("HTTP_LISTENER_FAILURE_DRAIN_PASS cause_preserved=1 request_completed=1")
}

func TestHTTPShutdownRejectsInvalidLifecycle(t *testing.T) {
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	for _, tc := range []struct {
		name     string
		ctx      context.Context
		server   *http.Server
		timeout  time.Duration
		nilServe bool
	}{
		{"nil_context", nil, &http.Server{}, time.Second, false},
		{"nil_server", context.Background(), nil, time.Second, false},
		{"zero_timeout", context.Background(), &http.Server{}, 0, false},
		{"negative_timeout", context.Background(), &http.Server{}, -time.Second, false},
		{"nil_serve", context.Background(), &http.Server{}, time.Second, true},
		{"already_canceled", canceled, &http.Server{}, time.Second, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			called := false
			var serve func() error = func() error { called = true; return nil }
			if tc.nilServe {
				serve = nil
			}
			err := serveUntilShutdown(tc.ctx, tc.server, tc.timeout, serve)
			if err == nil || called {
				t.Fatalf("invalid lifecycle started: called=%v error=%v", called, err)
			}
			if tc.name == "already_canceled" && !errors.Is(err, context.Canceled) {
				t.Fatal("cancellation cause lost")
			}
		})
	}
	if err := ServeUntilShutdown(context.Background(), nil, time.Second); err == nil {
		t.Fatal("nil public server accepted")
	}
}
