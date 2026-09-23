// AUTHORED composition over pinned official OpenTelemetry and Prometheus APIs.
// Trusted local qualification only. It does not enable a production endpoint.
package core

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	exporter "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/propagation"
	sdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/exemplar"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace/noop"
)

func Instrument(next http.Handler) (http.Handler, http.Handler, func(context.Context) error, error) {
	registry := prometheus.NewRegistry()
	reader, err := exporter.New(exporter.WithRegisterer(registry), exporter.WithoutScopeInfo(), exporter.WithoutTargetInfo())
	if err != nil {
		return nil, nil, nil, err
	}
	view := func(i sdk.Instrument) (sdk.Stream, bool) {
		if i.Name != "http.server.request.duration" {
			return sdk.Stream{Aggregation: sdk.AggregationDrop{}}, true
		}
		return sdk.Stream{
			Name: i.Name, Unit: "s",
			Aggregation:     sdk.AggregationExplicitBucketHistogram{Boundaries: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 5, 10}},
			AttributeFilter: attribute.NewAllowKeysFilter(attribute.Key("http.request.method"), attribute.Key("http.response.status_code")),
		}, true
	}
	provider := sdk.NewMeterProvider(sdk.WithReader(reader), sdk.WithResource(resource.Empty()), sdk.WithView(view), sdk.WithCardinalityLimit(128), sdk.WithExemplarFilter(exemplar.AlwaysOffFilter))
	wrapped := otelhttp.NewHandler(next, "reference-http", otelhttp.WithMeterProvider(provider), otelhttp.WithTracerProvider(noop.NewTracerProvider()), otelhttp.WithPropagators(propagation.NewCompositeTextMapPropagator()))
	return wrapped, promhttp.HandlerFor(registry, promhttp.HandlerOpts{ErrorHandling: promhttp.HTTPErrorOnError, Timeout: 5e9}), provider.Shutdown, nil
}
