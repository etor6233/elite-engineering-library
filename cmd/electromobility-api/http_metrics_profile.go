package main

import (
	"context"
	"errors"
	"net/http"
)

// Optional owner hook: smaller profiles keep the same host without importing
// metrics dependencies. Explicit activation without its owner fails closed.
var httpMetricsFactory func(http.Handler, func(string) string) (http.Handler, func(context.Context) error, func(context.Context) error, error)

func prepareHTTPMetrics(h http.Handler, lookup func(string) string) (http.Handler, func(context.Context) error, func(context.Context) error, error) {
	switch lookup("HTTP_METRICS_ENABLED") {
	case "", "false":
		return h, nil, nil, nil
	case "true":
		if httpMetricsFactory == nil {
			return nil, nil, nil, errors.New("HTTP metrics owner is not selected")
		}
		return httpMetricsFactory(h, lookup)
	default:
		return nil, nil, nil, errors.New("HTTP_METRICS_ENABLED must be true or false")
	}
}
