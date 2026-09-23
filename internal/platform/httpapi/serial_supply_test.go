package httpapi

// AUTHORED trust-boundary regressions. The spy asserts no repository invocation
// for malformed commands and that broad identity grants are narrowed per route.
import (
	"context"
	"elite.local/enterprise/internal/platform/identity"
	sc "elite.local/enterprise/internal/serialsupply"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type supplyBoundaryVerifier struct{}

func (supplyBoundaryVerifier) Verify(context.Context, string) (identity.Principal, error) {
	return identity.Principal{Subject: "operator", TenantID: "tenant", Permissions: map[string]struct{}{"*": {}}, Organizations: map[string]struct{}{"*": {}}}, nil
}

type supplyBoundaryStore struct {
	calls   int
	p       identity.Principal
	command sc.Command
}

func (s *supplyBoundaryStore) Plan(context.Context, identity.Principal, string, string) (sc.Plan, error) {
	s.calls++
	return sc.Plan{}, nil
}
func (s *supplyBoundaryStore) CommandReceipt(context.Context, identity.Principal, string, string) (sc.Receipt, error) {
	s.calls++
	return sc.Receipt{}, nil
}
func (s *supplyBoundaryStore) BindPlan(context.Context, identity.Principal, sc.PlanRequest) (sc.Receipt, error) {
	s.calls++
	return sc.Receipt{}, nil
}
func (s *supplyBoundaryStore) Apply(_ context.Context, p identity.Principal, r sc.Command) (sc.Receipt, error) {
	s.calls++
	s.p = p
	s.command = r
	return sc.Receipt{Version: r.ExpectedVersion}, nil
}
func TestSerialSupplyHTTPBoundaries(t *testing.T) {
	raw := `{"purchase_order_id":"po","command_id":"receive-1","kind":"receive","expected_version":"9007199254740993","evidence_sha256":"` + strings.Repeat("a", 64) + `","shipment_id":"asn","units":["u1"]}`
	for _, c := range []struct {
		name, body, query string
		status            int
	}{
		{"exact-int64", raw, "organization_id=store", 201},
		{"number-not-string", strings.Replace(raw, `"9007199254740993"`, `9007199254740993`, 1), "organization_id=store", 400},
		{"duplicate-key", strings.Replace(raw, `"command_id":"receive-1"`, `"command_id":"receive-1","command_id":"changed"`, 1), "organization_id=store", 400},
		{"case-collision", strings.Replace(raw, `"command_id":"receive-1"`, `"command_id":"receive-1","COMMAND_ID":"changed"`, 1), "organization_id=store", 400},
		{"actor-injection", strings.TrimSuffix(raw, "}") + `,"actor":"other"}`, "organization_id=store", 400},
		{"path-body-mismatch", strings.Replace(raw, `"po"`, `"other"`, 1), "organization_id=store", 400},
		{"kind-route-mismatch", strings.Replace(raw, `"receive"`, `"ship"`, 1), "organization_id=store", 400},
		{"trailing-object", raw + "{}", "organization_id=store", 400},
		{"oversize", raw + strings.Repeat(" ", 32769), "organization_id=store", 413},
		{"missing-scope", raw, "", 400},
		{"repeated-scope", raw, "organization_id=store&organization_id=other", 400},
		{"write-cursor", raw, "organization_id=store&after_unit=u", 400},
		{"unknown-query", raw, "organization_id=store&actor=other", 400},
	} {
		t.Run(c.name, func(t *testing.T) {
			store := &supplyBoundaryStore{}
			mux := http.NewServeMux()
			SerialSupplyModule{Service: store}.Register(mux, supplyBoundaryVerifier{})
			req := httptest.NewRequest("POST", "/v1/franchise/supply/orders/po/receive?"+c.query, strings.NewReader(c.body))
			req.Header.Set("Authorization", "Bearer fixture")
			out := httptest.NewRecorder()
			mux.ServeHTTP(out, req)
			if out.Code != c.status || out.Header().Get("Cache-Control") != "no-store" {
				t.Fatal(out.Code, out.Body.String())
			}
			if c.status != 201 {
				if store.calls != 0 {
					t.Fatal("invalid input invoked repository")
				}
				return
			}
			if store.calls != 1 || store.command.ExpectedVersion != 9007199254740993 || len(store.p.Permissions) != 1 || !store.p.Allowed("supply:receive") || store.p.Allowed("supply:plan") || len(store.p.Organizations) != 1 {
				t.Fatal("authority/version changed", store)
			}
			if _, ok := store.p.Organizations["store"]; !ok {
				t.Fatal("route organization not bound")
			}
			var value map[string]any
			if err := json.Unmarshal(out.Body.Bytes(), &value); err != nil || value["version"] != "9007199254740993" {
				t.Fatal("response lost int64 precision", value, err)
			}
		})
	}
}
