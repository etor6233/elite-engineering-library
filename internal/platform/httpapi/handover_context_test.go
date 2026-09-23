package httpapi

import (
	"context"
	"elite.local/enterprise/internal/franchisejourney"
	"net/http"
	"net/http/httptest"
	"testing"
)

type handoverContextSpy struct {
	calls              int
	tenant, org, order string
}

func (s *handoverContextSpy) OperatorContext(_ context.Context, t, o, id string) (franchisejourney.HandoverOperatorContext, error) {
	s.calls++
	s.tenant = t
	s.org = o
	s.order = id
	return franchisejourney.HandoverOperatorContext{}, nil
}
func TestHandoverContextAuthorization(t *testing.T) {
	for _, mode := range []string{"allowed", "unauthenticated", "permission", "organization"} {
		t.Run(mode, func(t *testing.T) {
			p := initialHandoverPrincipal()
			want := 200
			auth := "Bearer fixture"
			switch mode {
			case "unauthenticated":
				want = 401
				auth = ""
			case "permission":
				want = 403
				p.Permissions = map[string]struct{}{}
			case "organization":
				want = 403
				p.Organizations = map[string]struct{}{}
			}
			s := &handoverContextSpy{}
			mux := http.NewServeMux()
			HandoverContextModule{Service: s}.Register(mux, journeyVerifier{principal: p})
			r := httptest.NewRequest("GET", "/v1/franchise/orders/order/handover-context?organization_id=store", nil)
			r.Header.Set("Authorization", auth)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, r)
			if w.Code != want {
				t.Fatal(w.Code, want)
			}
			if want != 200 && s.calls != 0 {
				t.Fatal("unauthorized repository call")
			}
			if want == 200 && (s.tenant != "verified-tenant" || s.org != "store" || s.order != "order" || w.Header().Get("Cache-Control") != "no-store") {
				t.Fatal("unbound context")
			}
		})
	}
}
