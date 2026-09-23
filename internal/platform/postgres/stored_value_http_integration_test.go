package postgres_test

// AUTHORED HTTP composition fixture. Only bearer verification is a principal
// fixture; real HTTP decoding, permission checks, PG owners and source IPC run.
import (
	"context"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	sv "elite.local/enterprise/internal/storedvaluebridge"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type storedValueFixtureVerifier struct{ p identity.Principal }

func (v storedValueFixtureVerifier) Verify(ctx context.Context, raw string) (identity.Principal, error) {
	if raw != "local-bearer-fixture" {
		return identity.Principal{}, fmt.Errorf("invalid fixture bearer")
	}
	return v.p, nil
}
func storedValueHTTPCall(t *testing.T, store *db.StoredValue, p identity.Principal, method, path, key, body string, want int) []byte {
	t.Helper()
	mux := http.NewServeMux()
	httpapi.StoredValueModule{Service: store}.Register(mux, storedValueFixtureVerifier{p})
	server := httptest.NewServer(mux)
	defer server.Close()
	request, e := http.NewRequest(method, server.URL+path, strings.NewReader(body))
	if e != nil {
		t.Fatal(e)
	}
	request.Header.Set("Authorization", "Bearer local-bearer-fixture")
	request.Header.Set("Content-Type", "application/json")
	if key != "" {
		request.Header.Set("Idempotency-Key", key)
	}
	client := http.Client{Timeout: 8 * time.Second}
	response, e := client.Do(request)
	if e != nil {
		t.Fatal(e)
	}
	defer response.Body.Close()
	raw, e := io.ReadAll(io.LimitReader(response.Body, 262145))
	if e != nil || len(raw) > 262144 {
		t.Fatal(e)
	}
	if response.StatusCode != want {
		t.Fatalf("%s %s status%d want%d: %s", method, path, response.StatusCode, want, raw)
	}
	return raw
}
func storedValueHTTPView(t *testing.T, raw []byte) sv.ApprovalView {
	t.Helper()
	var v struct {
		ApprovalID string `json:"approval_id"`
		SHA        string `json:"payload_sha256"`
		State      string `json:"state"`
		Replay     bool   `json:"replay"`
		Payload    string `json:"payload_json"`
		Receipt    string `json:"receipt_json"`
	}
	if json.Unmarshal(raw, &v) != nil {
		t.Fatal("wire view")
	}
	var payload sv.BoundOperation
	if sv.Decode([]byte(v.Payload), &payload) != nil {
		t.Fatal("exact payload")
	}
	_, hash, e := sv.Canonical(payload)
	if e != nil || hash != v.SHA {
		t.Fatal("wire hash")
	}
	return sv.ApprovalView{ApprovalID: v.ApprovalID, PayloadSHA256: v.SHA, State: v.State, Replay: v.Replay, Payload: payload, Receipt: json.RawMessage(v.Receipt)}
}
func storedValueHTTPPropose(t *testing.T, store *db.StoredValue, p identity.Principal, r sv.Request) sv.ApprovalView {
	t.Helper()
	raw, e := json.Marshal(map[string]string{"operation_id": r.OperationID, "operation": r.Operation, "organization_id": r.OrganizationID, "program_id": r.ProgramID, "order_id": r.OrderID, "expected_order_version": sv.Number(r.ExpectedOrderVersion), "account_id": r.AccountID, "original_operation_id": r.OriginalOperationID, "profile_sha256": r.ProfileSHA256})
	if e != nil {
		t.Fatal(e)
	}
	return storedValueHTTPView(t, storedValueHTTPCall(t, store, p, "POST", "/v1/franchise/stored-value/proposals", "", string(raw), 201))
}
func storedValueHTTPApprove(t *testing.T, store *db.StoredValue, p identity.Principal, v sv.ApprovalView) sv.ApprovalView {
	t.Helper()
	raw, _ := json.Marshal(map[string]any{"organization_id": v.Payload.Request.OrganizationID, "payload_sha256": v.PayloadSHA256, "approved": true, "reason": "Exact source result reviewed through HTTP fixture"})
	out := storedValueHTTPView(t, storedValueHTTPCall(t, store, p, "POST", "/v1/franchise/stored-value/approvals/"+v.ApprovalID+"/decision", "", string(raw), 200))
	recovered := storedValueHTTPView(t, storedValueHTTPCall(t, store, p, "GET", "/v1/franchise/stored-value/approvals/"+v.ApprovalID+"?organization_id="+v.Payload.Request.OrganizationID, "", "", 200))
	if string(recovered.Receipt) != string(out.Receipt) || recovered.PayloadSHA256 != out.PayloadSHA256 {
		t.Fatal("HTTP decision recovery")
	}
	return out
}
func TestStoredValueHTTPRejectsScopeAndUnsafeNumbers(t *testing.T) {
	ctx := context.Background()
	pool := connectedPool(t)
	tenant := connectedSeedOrder(t, pool, "stripe")
	store, p, _ := fixtureStoredValue(t, pool, tenant, "gift_card")
	p.Organizations = map[string]struct{}{"store": {}, "other": {}}
	storedValueHTTPCall(t, store, p, "GET", "/v1/franchise/stored-value/orders/order?organization_id=other", "", "", 403)
	body := `{"operation_id":"invalid-number","operation":"issue","organization_id":"store","program_id":"reference-gift_card","order_id":"source","expected_order_version":9007199254740993,"account_id":"","original_operation_id":"","profile_sha256":"` + store.ProfileSHA256() + `"}`
	storedValueHTTPCall(t, store, p, "POST", "/v1/franchise/stored-value/proposals", "", body, 400)
	if _, e := pool.Exec(ctx, `insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'exact-large','store','customer','placed','ARS',9007199254740993,9007199254740993)`, tenant); e != nil {
		t.Fatal(e)
	}
	raw := storedValueHTTPCall(t, store, p, "GET", "/v1/franchise/stored-value/orders/exact-large?organization_id=store", "", "", 200)
	var out map[string]any
	if json.Unmarshal(raw, &out) != nil || out["gross_minor_units"] != "9007199254740993" || out["order_version"] != "9007199254740993" {
		t.Fatal("large amount/version lost", string(raw))
	}
	denied := p
	denied.Permissions = map[string]struct{}{"stored_value:read": {}}
	storedValueHTTPCall(t, store, denied, "POST", "/v1/franchise/stored-value/orders/order/funding", "permission-negative-key", `{"organization_id":"store","expected_order_version":"2"}`, 403)
	t.Log("actual HTTP/PG: large int64 strings exact; wrong configured scope and numeric JSON version rejected")
}
