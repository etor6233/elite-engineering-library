// AUTHORED loopback HTTP integration fixture; bearer identities are explicit
// fixture principals, not proof of JWT/IdP cryptography.
package postgres_test

import (
	"bytes"
	"context"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	sc "elite.local/enterprise/internal/serialsupply"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"
)

type supplyFixtureVerifier struct {
	mu         sync.Mutex
	next       int
	principals map[string]identity.Principal
}

func (v *supplyFixtureVerifier) token(p identity.Principal) string {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.next++
	token := fmt.Sprintf("supply-fixture-%d", v.next)
	v.principals[token] = p
	return token
}
func (v *supplyFixtureVerifier) Verify(_ context.Context, token string) (identity.Principal, error) {
	v.mu.Lock()
	defer v.mu.Unlock()
	p, ok := v.principals[token]
	if !ok {
		return p, identity.ErrUnauthenticated
	}
	return p, nil
}

type supplyHTTPReference struct {
	client        *http.Client
	base          string
	verifier      *supplyFixtureVerifier
	mu            sync.Mutex
	posts         map[string]int
	acceptedPosts map[string]int
	dropped       map[string]bool
	recoveries    int
}

func supplyFixtureOrg(p identity.Principal) string {
	values := []string{}
	for o := range p.Organizations {
		values = append(values, o)
	}
	sort.Strings(values)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}
func supplyFixtureSurface(p identity.Principal) string {
	if p.Allowed("supply:factory-read") {
		return "factory"
	}
	return "franchise"
}
func (c *supplyHTTPReference) call(ctx context.Context, p identity.Principal, method, path string, query url.Values, input, out any) error {
	if query == nil {
		query = url.Values{}
	}
	query.Set("organization_id", supplyFixtureOrg(p))
	var body []byte
	var err error
	if input != nil {
		body, err = json.Marshal(input)
		if err != nil {
			return err
		}
	}
	req, err := http.NewRequestWithContext(ctx, method, c.base+path+"?"+query.Encode(), bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.verifier.token(p))
	req.Header.Set("Content-Type", "application/json")
	response, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(response.Body, 1048577))
	if err != nil {
		return err
	}
	if len(raw) > 1048576 {
		return errors.New("oversized fixture response")
	}
	if response.Header.Get("Cache-Control") != "no-store" {
		return errors.New("private response not no-store")
	}
	if response.StatusCode >= 400 {
		return fmt.Errorf("HTTP %d: %s", response.StatusCode, raw)
	}
	return json.Unmarshal(raw, out)
}
func (c *supplyHTTPReference) BindPlan(ctx context.Context, p identity.Principal, r sc.PlanRequest) (sc.Receipt, error) {
	var out sc.Receipt
	err := c.call(ctx, p, "POST", "/v1/franchise/supply/orders/"+url.PathEscape(r.PurchaseOrderID)+"/plan", nil, r, &out)
	return out, err
}
func (c *supplyHTTPReference) Plan(ctx context.Context, p identity.Principal, id, after string) (sc.Plan, error) {
	var out sc.Plan
	q := url.Values{}
	if after != "" {
		q.Set("after_unit", after)
	}
	err := c.call(ctx, p, "GET", "/v1/"+supplyFixtureSurface(p)+"/supply/orders/"+url.PathEscape(id), q, nil, &out)
	return out, err
}
func (c *supplyHTTPReference) CommandReceipt(ctx context.Context, p identity.Principal, id, key string) (sc.Receipt, error) {
	var out sc.Receipt
	err := c.call(ctx, p, "GET", "/v1/"+supplyFixtureSurface(p)+"/supply/orders/"+url.PathEscape(id), url.Values{"command_id": {key}}, nil, &out)
	return out, err
}
func (c *supplyHTTPReference) Apply(ctx context.Context, p identity.Principal, r sc.Command) (sc.Receipt, error) {
	surface := "franchise"
	switch r.Kind {
	case "confirm", "start", "register", "milestone", "ship":
		surface = "factory"
	}
	var out sc.Receipt
	err := c.call(ctx, p, "POST", "/v1/"+surface+"/supply/orders/"+url.PathEscape(r.PurchaseOrderID)+"/"+r.Kind, nil, r, &out)
	if err == nil {
		return out, nil
	}
	// Only transport loss triggers fixture recovery. A rejected HTTP command is
	// never silently retried, and no second POST is made here.
	var network *url.Error
	if !errors.As(err, &network) {
		return out, err
	}
	recovered, getErr := c.CommandReceipt(ctx, p, r.PurchaseOrderID, r.CommandID)
	if getErr != nil {
		return out, fmt.Errorf("unconfirmed POST %v; recovery %v", err, getErr)
	}
	_, hash, hashErr := sc.Canonical(r)
	if hashErr != nil || recovered.PurchaseOrderID != r.PurchaseOrderID || recovered.CommandID != r.CommandID || recovered.Actor != p.Subject || recovered.RequestSHA256 != hash {
		return out, errors.New("recovery evidence mismatch")
	}
	c.mu.Lock()
	c.recoveries++
	c.mu.Unlock()
	return recovered, nil
}
func TestSerialSupplyHTTPConnectedReference(t *testing.T) {
	var reference *supplyHTTPReference
	testSerialSupplyConnected(t, func(t *testing.T, _ *pgxpool.Pool, store *db.SerialSupply) serialSupplyReference {
		verifier := &supplyFixtureVerifier{principals: map[string]identity.Principal{}}
		mux := http.NewServeMux()
		httpapi.SerialSupplyModule{Service: store}.Register(mux, verifier)
		reference = &supplyHTTPReference{verifier: verifier, posts: map[string]int{}, acceptedPosts: map[string]int{}, dropped: map[string]bool{}}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var command sc.Command
			if r.Method == "POST" {
				raw, err := io.ReadAll(io.LimitReader(r.Body, 32769))
				if err != nil {
					t.Error(err)
				}
				r.Body.Close()
				r.Body = io.NopCloser(bytes.NewReader(raw))
				_ = json.Unmarshal(raw, &command)
				reference.mu.Lock()
				reference.posts[command.CommandID]++
				reference.mu.Unlock()
			}
			result := httptest.NewRecorder()
			mux.ServeHTTP(result, r)
			if r.Method == "POST" && (result.Code == 200 || result.Code == 201) {
				reference.mu.Lock()
				reference.acceptedPosts[command.CommandID]++
				reference.mu.Unlock()
			}
			drop := ""
			if r.Method == "POST" && result.Code == 201 {
				if command.Kind == "receive" && command.ShipmentID == "asn-one" {
					drop = "first-receipt"
				}
				if command.Kind == "ship" && command.ShipmentID == "asn-two" {
					drop = "last-shipment"
				}
			}
			reference.mu.Lock()
			lose := drop != "" && !reference.dropped[drop]
			if lose {
				reference.dropped[drop] = true
				reference.dropped["command:"+command.CommandID] = true
			}
			reference.mu.Unlock()
			if lose {
				connection, _, err := w.(http.Hijacker).Hijack()
				if err != nil {
					t.Error(err)
					return
				}
				connection.Close()
				return
			}
			for key, values := range result.Header() {
				for _, value := range values {
					w.Header().Add(key, value)
				}
			}
			w.WriteHeader(result.Code)
			_, _ = w.Write(result.Body.Bytes())
		}))
		transport := &http.Transport{DisableKeepAlives: true}
		reference.client = &http.Client{Transport: transport, Timeout: 5 * time.Second}
		reference.base = server.URL
		t.Cleanup(func() { transport.CloseIdleConnections(); server.Close() })
		return reference
	})
	reference.mu.Lock()
	defer reference.mu.Unlock()
	if reference.recoveries != 2 || !reference.dropped["first-receipt"] || !reference.dropped["last-shipment"] {
		t.Fatal("lost-response coverage", reference.recoveries, reference.dropped)
	}
	for key := range reference.dropped {
		if strings.HasPrefix(key, "command:") && reference.acceptedPosts[strings.TrimPrefix(key, "command:")] != 1 {
			t.Fatal("hidden POST retry", key, reference.posts)
		}
	}
	t.Log("SERIAL_SUPPLY_HTTP_PASS 30-version connected journey; two committed responses lost and recovered by actor/hash GET; generic approval/stock/purchase bypass refused")
}
