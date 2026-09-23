package postgres_test

// AUTHORED transport fixture. The verifier returns explicitly assigned test
// principals; JWT cryptography is not claimed by this fixture. Real HTTP,
// handlers, repository authorization, durable commands and PG effects execute.
import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/randomid"
	wc "elite.local/enterprise/internal/warrantyclaim"
)

type warrantyClaimCommands interface {
	OpenClaim(context.Context, identity.Principal, wc.OpenClaim) (wc.Step, error)
	Diagnose(context.Context, identity.Principal, wc.Diagnose) (wc.Step, error)
	PlanRepair(context.Context, identity.Principal, wc.Plan) (wc.Step, error)
	DecideRepair(context.Context, identity.Principal, wc.DecideRepair) (wc.Step, error)
	CompleteRepairWork(context.Context, identity.Principal, wc.CompleteWork) (wc.Step, error)
	RecordRepairQuality(context.Context, identity.Principal, wc.Quality) (wc.Step, error)
	AcceptRepair(context.Context, identity.Principal, wc.AcceptRepair) (wc.Step, error)
	ReconcileRepair(context.Context, identity.Principal, wc.ReconcileRepair) (wc.Step, error)
	CancelRepair(context.Context, identity.Principal, wc.CancelRepair) (wc.Step, error)
}
type warrantyHTTPFixture struct {
	mu         sync.Mutex
	server     *httptest.Server
	principals map[string]identity.Principal
	drops      map[string]int
	recoveries int
	posts      map[string]int
}

func (h *warrantyHTTPFixture) Verify(_ context.Context, token string) (identity.Principal, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	p, ok := h.principals[token]
	if !ok {
		return p, identity.ErrUnauthenticated
	}
	return p, nil
}
func newWarrantyHTTPFixture(t *testing.T, service httpapi.WarrantyService, drop bool) *warrantyHTTPFixture {
	t.Helper()
	h := &warrantyHTTPFixture{principals: map[string]identity.Principal{}, drops: map[string]int{}, posts: map[string]int{}}
	mux := http.NewServeMux()
	httpapi.WarrantyModule{Service: service}.Register(mux, h)
	h.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, r)
		shouldDrop := false
		if r.Method == "POST" {
			h.mu.Lock()
			h.posts[r.URL.Path]++
			if drop && recorder.Code >= 200 && recorder.Code < 300 && (strings.HasSuffix(r.URL.Path, "/work") || strings.HasSuffix(r.URL.Path, "/reconciliation")) && h.drops[r.URL.Path] == 0 {
				h.drops[r.URL.Path]++
				shouldDrop = true
			}
			h.mu.Unlock()
		}
		if shouldDrop {
			conn, _, err := w.(http.Hijacker).Hijack()
			if err != nil {
				t.Error(err)
				return
			}
			conn.Close()
			return
		}
		for key, values := range recorder.Header() {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.WriteHeader(recorder.Code)
		_, _ = w.Write(recorder.Body.Bytes())
	}))
	h.server.Client().Timeout = 10 * time.Second
	t.Cleanup(h.server.Close)
	return h
}
func (h *warrantyHTTPFixture) send(ctx context.Context, p identity.Principal, method, path string, body any, out any) error {
	raw, err := json.Marshal(body)
	if err != nil {
		return err
	}
	id := randomid.Generator{}.New()
	copy := identity.Principal{TenantID: p.TenantID, Subject: p.Subject, Permissions: map[string]struct{}{}, Organizations: map[string]struct{}{}}
	for k := range p.Permissions {
		copy.Permissions[k] = struct{}{}
	}
	for k := range p.Organizations {
		copy.Organizations[k] = struct{}{}
	}
	h.mu.Lock()
	h.principals[id] = copy
	h.mu.Unlock()
	req, err := http.NewRequestWithContext(ctx, method, h.server.URL+path, bytes.NewReader(raw))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+id)
	req.Header.Set("Content-Type", "application/json")
	response, err := h.server.Client().Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	data, err := io.ReadAll(io.LimitReader(response.Body, 131073))
	if err != nil {
		return err
	}
	if len(data) > 131072 {
		return fmt.Errorf("fixture response exceeded bound")
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("HTTP%d: %s", response.StatusCode, data)
	}
	return json.Unmarshal(data, out)
}
func (h *warrantyHTTPFixture) step(ctx context.Context, p identity.Principal, c wc.Command, surface, action string, body any) (wc.Step, error) {
	org := "store"
	if surface == "factory" {
		org = "factory"
	}
	path := "/v1/" + surface + "/warranty/claims/" + url.PathEscape(c.CaseID) + "/" + action + "?organization_id=" + org
	var out wc.Step
	err := h.send(ctx, p, "POST", path, body, &out)
	if err != nil {
		var recovered wc.Step
		lookup := "/v1/" + surface + "/warranty/claims/" + url.PathEscape(c.CaseID) + "?organization_id=" + org + "&command_id=" + url.QueryEscape(c.CommandID)
		if recovery := h.send(ctx, p, "GET", lookup, nil, &recovered); recovery == nil {
			_, expected, e := wc.Canonical(body)
			if e != nil || recovered.RequestSHA256 != expected || recovered.Actor != p.Subject {
				return wc.Step{}, fmt.Errorf("recovery did not bind original request")
			}
			h.mu.Lock()
			h.recoveries++
			h.mu.Unlock()
			return recovered, nil
		}
	}
	return out, err
}
func (h *warrantyHTTPFixture) OpenClaim(ctx context.Context, p identity.Principal, r wc.OpenClaim) (wc.Step, error) {
	surface := "franchise"
	if p.Subject == "customer" {
		surface = "customer"
	}
	var out wc.Step
	err := h.send(ctx, p, "POST", "/v1/"+surface+"/warranty/claims?organization_id=store", r, &out)
	return out, err
}

func (h *warrantyHTTPFixture) Diagnose(ctx context.Context, p identity.Principal, r wc.Diagnose) (wc.Step, error) {
	return h.step(ctx, p, r.Command, "franchise", "diagnosis", r)
}

func (h *warrantyHTTPFixture) PlanRepair(ctx context.Context, p identity.Principal, r wc.Plan) (wc.Step, error) {
	return h.step(ctx, p, r.Command, "franchise", "plan", r)
}

func (h *warrantyHTTPFixture) DecideRepair(ctx context.Context, p identity.Principal, r wc.DecideRepair) (wc.Step, error) {
	return h.step(ctx, p, r.Command, "franchise", "decision", r)
}

func (h *warrantyHTTPFixture) CompleteRepairWork(ctx context.Context, p identity.Principal, r wc.CompleteWork) (wc.Step, error) {
	return h.step(ctx, p, r.Command, "franchise", "work", r)
}

func (h *warrantyHTTPFixture) RecordRepairQuality(ctx context.Context, p identity.Principal, r wc.Quality) (wc.Step, error) {
	return h.step(ctx, p, r.Command, "franchise", "quality", r)
}

func (h *warrantyHTTPFixture) AcceptRepair(ctx context.Context, p identity.Principal, r wc.AcceptRepair) (wc.Step, error) {
	return h.step(ctx, p, r.Command, "customer", "acceptance", r)
}

func (h *warrantyHTTPFixture) ReconcileRepair(ctx context.Context, p identity.Principal, r wc.ReconcileRepair) (wc.Step, error) {
	return h.step(ctx, p, r.Command, "factory", "reconciliation", r)
}

func (h *warrantyHTTPFixture) CancelRepair(ctx context.Context, p identity.Principal, r wc.CancelRepair) (wc.Step, error) {
	return h.step(ctx, p, r.Command, "franchise", "cancellation", r)
}
