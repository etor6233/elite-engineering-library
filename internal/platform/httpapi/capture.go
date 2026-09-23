package httpapi

import (
	"context"
	cap "elite.local/enterprise/internal/capture"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

type CaptureDecoder interface {
	Decode(context.Context, []byte) cap.Result
}
type CaptureResolver interface {
	ResolveCapture(context.Context, identity.Principal, string, string, string) (postgres.CaptureResolution, error)
}
type CaptureModule struct {
	Decoder CaptureDecoder
	Store   CaptureResolver
}

func (m CaptureModule) principal(w http.ResponseWriter, r *http.Request, v identity.Verifier) (identity.Principal, bool) {
	values := r.Header.Values("Authorization")
	if v == nil || len(values) != 1 || !strings.HasPrefix(values[0], "Bearer ") || len(values[0]) > 16391 || len(strings.Fields(values[0])) != 2 {
		documentJSON(w, 401, map[string]string{"code": "UNAUTHENTICATED"})
		return identity.Principal{}, false
	}
	p, e := v.Verify(r.Context(), strings.TrimPrefix(values[0], "Bearer "))
	if e != nil {
		documentJSON(w, 401, map[string]string{"code": "UNAUTHENTICATED"})
		return p, false
	}
	if !postgres.CaptureAllowed(p) {
		documentJSON(w, 403, map[string]string{"code": "FORBIDDEN"})
		return p, false
	}
	if r.URL.RawQuery != "" {
		documentJSON(w, 400, map[string]string{"code": "INVALID_QUERY"})
		return p, false
	}
	return p, true
}
func (m CaptureModule) Register(mux *http.ServeMux, v identity.Verifier) {
	mux.HandleFunc("POST /v1/capture/decode", func(w http.ResponseWriter, r *http.Request) {
		if _, ok := m.principal(w, r, v); !ok {
			return
		}
		if m.Decoder == nil {
			documentJSON(w, 503, map[string]string{"code": "CAPTURE_UNAVAILABLE"})
			return
		}
		if r.Header.Get("Content-Type") != "application/octet-stream" {
			documentJSON(w, 415, map[string]string{"code": "UNSUPPORTED_MEDIA_TYPE"})
			return
		}
		controller := http.NewResponseController(w)
		deadlineErr := controller.SetReadDeadline(time.Now().Add(5 * time.Second))
		if deadlineErr != nil && !errors.Is(deadlineErr, http.ErrNotSupported) {
			documentJSON(w, 503, map[string]string{"code": "CAPTURE_UNAVAILABLE"})
			return
		}
		body, e := io.ReadAll(http.MaxBytesReader(w, r.Body, cap.MaxFrameBytes))
		if e != nil {
			var ne net.Error
			if errors.As(e, &ne) && ne.Timeout() {
				documentJSON(w, 408, map[string]string{"code": "FRAME_TIMEOUT"})
				return
			}
			documentJSON(w, 413, map[string]string{"code": "FRAME_TOO_LARGE"})
			return
		}
		if !cap.ValidFrame(body) {
			documentJSON(w, 400, map[string]string{"code": "INVALID_FRAME"})
			return
		}
		result := m.Decoder.Decode(r.Context(), body)
		documentJSON(w, 200, result)
	})
	mux.HandleFunc("POST /v1/capture/resolve", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v)
		if !ok {
			return
		}
		if m.Store == nil {
			documentJSON(w, 503, map[string]string{"code": "CAPTURE_UNAVAILABLE"})
			return
		}
		var b struct {
			Text      string `json:"text"`
			Encoding  string `json:"encoding"`
			Operation string `json:"operation"`
		}
		if documentBody(w, r, &b) != nil {
			documentJSON(w, 400, map[string]string{"code": "INVALID_CONTRACT"})
			return
		}
		result, e := m.Store.ResolveCapture(r.Context(), p, b.Text, b.Encoding, b.Operation)
		if e != nil {
			switch {
			case errors.Is(e, cap.ErrContract):
				documentJSON(w, 400, map[string]string{"code": "INVALID_CONTRACT"})
			case errors.Is(e, postgres.ErrCaptureNotFound):
				documentJSON(w, 404, map[string]string{"code": "NOT_FOUND"})
			default:
				documentJSON(w, 503, map[string]string{"code": "CAPTURE_UNAVAILABLE"})
			}
			return
		}
		documentJSON(w, 200, result)
	})
}
