package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRealHTTPStatusAndClosedLabels(t *testing.T) {
	handler, metrics, close, err := instrument(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/fail" {
			w.WriteHeader(503)
			return
		}
		w.WriteHeader(204)
	}))
	if err != nil {
		t.Fatal(err)
	}
	defer close(context.Background())
	server := httptest.NewServer(handler)
	defer server.Close()
	for _, path := range []string{"/ok?secret=private@example.invalid", "/fail"} {
		req, _ := http.NewRequest("GET", server.URL+path, strings.NewReader("PRIVATE_BODY"))
		req.Host = "PRIVATE_HOST.invalid"
		req.Header.Set("Authorization", "Bearer PRIVATE_TOKEN")
		req.Header.Set("Cookie", "PRIVATE_COOKIE")
		req.Header.Set("X-Private", "PRIVATE_HEADER")
		res, err := server.Client().Do(req)
		if err != nil {
			t.Fatal(err)
		}
		io.Copy(io.Discard, res.Body)
		res.Body.Close()
		if path == "/fail" && res.StatusCode != 503 {
			t.Fatal(res.StatusCode)
		}
	}
	w := httptest.NewRecorder()
	metrics.ServeHTTP(w, httptest.NewRequest("GET", "http://metrics/", nil))
	body := w.Body.String()
	for _, marker := range []string{"PRIVATE_", "private@example", "/fail", "/ok", "server_address", "url_", "http_route", "target_info", "otel_scope", "http_server_request_body"} {
		if strings.Contains(body, marker) {
			t.Fatalf("unexpected metric content %s", marker)
		}
	}
	for _, line := range []string{
		"http_server_request_duration_seconds_count{http_request_method=\"GET\",http_response_status_code=\"204\"} 1",
		"http_server_request_duration_seconds_count{http_request_method=\"GET\",http_response_status_code=\"503\"} 1",
	} {
		if !strings.Contains(body, line) {
			t.Fatalf("missing %s in %s", line, body)
		}
	}
}

func TestIndependentInstancesDoNotMixRequests(t *testing.T) {
	a, am, ac, err := instrument(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(201) }))
	if err != nil {
		t.Fatal(err)
	}
	defer ac(context.Background())
	b, bm, bc, err := instrument(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(503) }))
	if err != nil {
		t.Fatal(err)
	}
	defer bc(context.Background())
	a.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("POST", "http://reference/order", nil))
	b.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "http://reference/failure", nil))
	for _, c := range []struct {
		h            http.Handler
		want, reject string
	}{
		{am, "http_response_status_code=\"201\"", "http_response_status_code=\"503\""},
		{bm, "http_response_status_code=\"503\"", "http_response_status_code=\"201\""},
	} {
		w := httptest.NewRecorder()
		c.h.ServeHTTP(w, httptest.NewRequest("GET", "http://metrics/", nil))
		if !strings.Contains(w.Body.String(), c.want) || strings.Contains(w.Body.String(), c.reject) {
			t.Fatal("instance metrics mixed")
		}
		for _, forbidden := range []string{"go_gc", "go_goroutines", "process_", "target_info"} {
			if strings.Contains(w.Body.String(), forbidden) {
				t.Fatal("default collector leaked", forbidden)
			}
		}
	}
}
