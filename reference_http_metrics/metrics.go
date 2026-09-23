// AUTHORED compatibility glue; the original closed instrument lives in core.
package main

import (
	"context"
	"elite.local/enterprise/reference_http_metrics/core"
	"net/http"
)

func instrument(h http.Handler) (http.Handler, http.Handler, func(context.Context) error, error) {
	return core.Instrument(h)
}
