package whatsappbridge

// AUTHORED regressions of the new HTTP boundary; no database or provider calls.
import (
	"context"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

type scheduledBoundaryVerifier struct {
	principal identity.Principal
	err       error
}

func (v scheduledBoundaryVerifier) Verify(context.Context, string) (identity.Principal, error) {
	return v.principal, v.err
}

type scheduledDecodeWriter struct{ *httptest.ResponseRecorder }

func (scheduledDecodeWriter) SetReadDeadline(time.Time) error { return nil }

func TestScheduledRequestBoundary(t *testing.T) {
	for _, raw := range []string{
		`{"request_id":"one","REQUEST_ID":"two"}`,
		`{"extra":true}`, `{}{}`, `{"not_before":"invalid"}`,
		strings.Repeat(" ", 32769), `{"request":{"recipient":"one","Recipient":"two"}}`,
	} {
		r := httptest.NewRequest("POST", "/", strings.NewReader(raw))
		r.Header.Set("Content-Type", "application/json")
		if decodeSchedule(scheduledDecodeWriter{httptest.NewRecorder()}, r, new(ScheduleRequest)) == nil {
			t.Fatalf("accepted ambiguous or oversized request %q", raw[:min(len(raw), 80)])
		}
	}
	for _, headers := range []struct{ content, url string }{{"text/plain", "/"}, {"application/json", "/?override=1"}} {
		r := httptest.NewRequest("POST", headers.url, strings.NewReader("{}"))
		r.Header.Set("Content-Type", headers.content)
		if decodeSchedule(scheduledDecodeWriter{httptest.NewRecorder()}, r, new(ScheduleRequest)) == nil {
			t.Fatal("unsupported request metadata")
		}
	}
}

func TestScheduledAuthorizationBoundary(t *testing.T) {
	approved := identity.Principal{Subject: "operator", TenantID: "tenant", Permissions: map[string]struct{}{"notification:request": {}}, Organizations: map[string]struct{}{"store": {}}}
	for _, tc := range []struct {
		name string
		p    identity.Principal
		auth []string
		err  error
		want int
	}{
		{"missing bearer", approved, nil, nil, 401},
		{"duplicate bearer", approved, []string{"Bearer one", "Bearer two"}, nil, 401},
		{"malformed bearer", approved, []string{"Bearer one extra"}, nil, 401},
		{"rejected token", approved, []string{"Bearer one"}, errors.New("fixture"), 401},
		{"other tenant", identity.Principal{Subject: "operator", TenantID: "other", Permissions: approved.Permissions, Organizations: approved.Organizations}, []string{"Bearer one"}, nil, 403},
		{"other org", identity.Principal{Subject: "operator", TenantID: "tenant", Permissions: approved.Permissions, Organizations: map[string]struct{}{"other": {}}}, []string{"Bearer one"}, nil, 403},
		{"no permission", identity.Principal{Subject: "operator", TenantID: "tenant", Organizations: approved.Organizations}, []string{"Bearer one"}, nil, 403},
	} {
		t.Run(tc.name, func(t *testing.T) {
			m := &ScheduledNotifications{approvals: &ScheduleApprovals{tenant: "tenant", org: "store"}}
			mux := http.NewServeMux()
			m.Register(mux, scheduledBoundaryVerifier{tc.p, tc.err})
			r := httptest.NewRequest("POST", "/v1/franchise/notifications/scheduled/prepare", strings.NewReader("{}"))
			r.Header["Authorization"] = tc.auth
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, r)
			if w.Code != tc.want {
				t.Fatalf("status %d want %d", w.Code, tc.want)
			}
		})
	}
}

func FuzzScheduledRequestBoundary(f *testing.F) {
	for _, seed := range []string{`{}`, `{"request_id":"schedule-0000001","appointment_id":"appointment-1","recipient":"5491112345678","template_name":"order_update","language_code":"es_AR","body_parameters":["A-1","recordatorio"],"not_before":"2026-09-15T12:00:00Z","expires_at":"2026-09-15T13:00:00Z"}`, `{"request_id":"one","REQUEST_ID":"two"}`, `{"body_parameters":["\ud800"]}`} {
		f.Add([]byte(seed))
	}
	f.Fuzz(func(t *testing.T, raw []byte) {
		if len(raw) > 32769 {
			t.Skip()
		}
		decode := func(body []byte) (ScheduleRequest, error) {
			r := httptest.NewRequest("POST", "/", strings.NewReader(string(body)))
			r.Header.Set("Content-Type", "application/json")
			var request ScheduleRequest
			err := decodeSchedule(scheduledDecodeWriter{httptest.NewRecorder()}, r, &request)
			return request, err
		}
		first, e := decode(raw)
		if e != nil {
			return
		}
		canonical, e := json.Marshal(first)
		if e != nil {
			t.Fatal(e)
		}
		second, e := decode(canonical)
		if e != nil || !reflect.DeepEqual(first, second) {
			t.Fatal("accepted request cannot round-trip through its exact contract")
		}
	})
}
