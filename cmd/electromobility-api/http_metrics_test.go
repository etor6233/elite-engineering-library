package main

import (
	"context"
	"net/http"
	"testing"
)

func TestHTTPMetricsExplicitLocalProfile(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(204) })
	for _, address := range []string{"", ":9090", "0.0.0.0:9090", "example.com:9090", "127.0.0.1:80", "127.0.0.1:09090", "[::1]:9090"} {
		if _, _, _, err := localHTTPMetrics(h, func(string) string { return address }); err == nil {
			t.Fatal("unbound listener accepted")
		}
	}
	lookup := func(key string) string {
		if key == "HTTP_METRICS_ENABLED" {
			return "true"
		}
		return "127.0.0.1:19090"
	}
	wrapped, run, close, err := prepareHTTPMetrics(h, lookup)
	if err != nil || wrapped == nil || run == nil || close == nil {
		t.Fatal("owner not assembled", err)
	}
	if err := close(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, run, close, err := prepareHTTPMetrics(h, func(string) string { return "" }); err != nil || run != nil || close != nil {
		t.Fatal("disabled profile has effects")
	}
	if _, _, _, err := prepareHTTPMetrics(h, func(string) string { return "TRUE" }); err == nil {
		t.Fatal("ambiguous activation accepted")
	}
}
