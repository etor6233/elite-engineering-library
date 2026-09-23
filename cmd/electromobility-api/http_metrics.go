package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strconv"
	"time"

	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/reference_http_metrics/core"
)

func init() { httpMetricsFactory = localHTTPMetrics }

func localHTTPMetrics(h http.Handler, lookup func(string) string) (http.Handler, func(context.Context) error, func(context.Context) error, error) {
	address := lookup("HTTP_METRICS_ADDRESS")
	host, port, err := net.SplitHostPort(address)
	n, parseErr := strconv.Atoi(port)
	if err != nil || parseErr != nil || host != "127.0.0.1" || n < 1024 || n > 65535 || strconv.Itoa(n) != port {
		return nil, nil, nil, errors.New("local metrics requires explicit loopback address")
	}
	wrapped, metrics, close, err := core.Instrument(h)
	if err != nil {
		return nil, nil, nil, err
	}
	server := &http.Server{Addr: address, Handler: metrics, ReadHeaderTimeout: 3 * time.Second, ReadTimeout: 5 * time.Second, WriteTimeout: 5 * time.Second, IdleTimeout: 10 * time.Second, MaxHeaderBytes: 4096}
	run := func(ctx context.Context) error { return httpapi.ServeUntilShutdown(ctx, server, 5*time.Second) }
	return wrapped, run, close, nil
}
