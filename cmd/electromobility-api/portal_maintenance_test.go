package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type portalFixtureBroker struct{ failure bool }

func (b portalFixtureBroker) AccessToken(context.Context) (string, error) {
	if b.failure {
		return "", errors.New("fixture grant unavailable")
	}
	return "synthetic-service", nil
}
func TestPortalMaintenanceResponseBoundary(t *testing.T) {
	cases := []struct {
		name, body, media string
		status            int
		ok                bool
	}{
		{"empty", `{"claimed":0,"confirmed":0,"pending":0,"purged":0,"unconfirmed_purged":0}`, "application/json", 200, true},
		{"confirmed", `{"claimed":2,"confirmed":2,"pending":0,"purged":100,"unconfirmed_purged":1}`, "application/json", 200, true},
		{"pending", `{"claimed":1,"confirmed":0,"pending":1,"purged":0,"unconfirmed_purged":0}`, "application/json", 503, false},
		{"missing", `{"claimed":0,"confirmed":0,"pending":0,"purged":0}`, "application/json", 200, false},
		{"unknown", `{"claimed":0,"confirmed":0,"pending":0,"purged":0,"other":0}`, "application/json", 200, false},
		{"inconsistent", `{"claimed":2,"confirmed":1,"pending":0,"purged":0,"unconfirmed_purged":0}`, "application/json", 200, false},
		{"negative", `{"claimed":-1,"confirmed":0,"pending":0,"purged":0,"unconfirmed_purged":0}`, "application/json", 200, false},
		{"purge_bound", `{"claimed":0,"confirmed":0,"pending":0,"purged":101,"unconfirmed_purged":0}`, "application/json", 200, false},
		{"revocation_bound", `{"claimed":0,"confirmed":0,"pending":0,"purged":1,"unconfirmed_purged":2}`, "application/json", 200, false},
		{"claim_bound", `{"claimed":3,"confirmed":3,"pending":0,"purged":0,"unconfirmed_purged":0}`, "application/json", 200, false},
		{"wrong_media", `{"claimed":0,"confirmed":0,"pending":0,"purged":0,"unconfirmed_purged":0}`, "text/plain", 200, false},
		{"json_prefix", `{"claimed":0,"confirmed":0,"pending":0,"purged":0,"unconfirmed_purged":0}`, "application/jsonjunk", 200, false},
		{"wrong_status", `{"claimed":0,"confirmed":0,"pending":0,"purged":0,"unconfirmed_purged":0}`, "application/json", 500, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				if r.Method != "POST" || r.Header.Get("Authorization") != "Bearer synthetic-service" {
					t.Error("request contract")
				}
				w.Header().Set("Content-Type", tc.media)
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte(tc.body))
			}))
			defer server.Close()
			p := portalSessionRuntime{broker: portalFixtureBroker{}, target: server.URL, client: server.Client()}
			if err := p.poll(context.Background()); (err == nil) != tc.ok {
				t.Fatalf("unexpected result %v", err)
			}
			if calls != 1 {
				t.Fatal("retried poll")
			}
		})
	}
}
func TestPortalMaintenanceCancellationAndGrant(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls++; w.WriteHeader(500) }))
	defer server.Close()
	p := portalSessionRuntime{broker: portalFixtureBroker{failure: true}, target: server.URL, client: server.Client()}
	if p.poll(context.Background()) == nil || calls != 0 {
		t.Fatal("grant failure must not reach target")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	done := make(chan struct{})
	go func() { defer close(done); p.run(ctx) }()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("shutdown not joined")
	}
	if calls != 0 {
		t.Fatal("cancelled worker ran")
	}
}
func TestPortalActivationRequiresSelectedPack(t *testing.T) {
	old := portalRuntimeFactory
	defer func() { portalRuntimeFactory = old }()
	portalRuntimeFactory = nil
	for _, enabled := range []string{"true", "TRUE", "1"} {
		_, err := selectedPortalHost(context.Background(), nil, func(string) string { return enabled })
		if err == nil {
			t.Fatal("activation accepted", enabled)
		}
	}
	for _, enabled := range []string{"", "false"} {
		runtime, err := selectedPortalHost(context.Background(), nil, func(string) string { return enabled })
		if err != nil || runtime != nil {
			t.Fatal("disabled lifecycle not inert")
		}
	}
}
