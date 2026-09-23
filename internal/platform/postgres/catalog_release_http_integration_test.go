package postgres_test

// AUTHORED real HTTP/PG fixture. Bearer verifier assigns explicit principals;
// this does not substitute T2803 JWT/IdP evidence.
import (
	"bytes"
	"context"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"
)

type catalogHTTPFixture struct {
	publisher identity.Principal
	feedProbe func(context.Context, identity.Principal, int64, []byte)
	webProbe  func(cr.PublicDocument)

	t            *testing.T
	client       *http.Client
	base         string
	mu           sync.Mutex
	counter      int
	principals   map[string]identity.Principal
	accepted     map[string]int
	dropped      bool
	recoveries   int
	etag         string
	cached       []byte
	cacheChanges int
}

func (c *catalogHTTPFixture) Verify(_ context.Context, token string) (identity.Principal, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	p, ok := c.principals[token]
	if !ok {
		return p, identity.ErrUnauthenticated
	}
	return p, nil
}
func (c *catalogHTTPFixture) request(ctx context.Context, p *identity.Principal, method, path, contentType string, body []byte, etag string) ([]byte, int, string, error) {
	if p != nil {
		org := ""
		for v := range p.Organizations {
			org = v
			break
		}
		path += "?organization_id=" + url.QueryEscape(org)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.base+path, bytes.NewReader(body))
	if err != nil {
		return nil, 0, "", err
	}
	req.Header.Set("Content-Type", contentType)
	if etag != "" {
		req.Header.Set("If-None-Match", etag)
	}
	if p != nil {
		c.mu.Lock()
		c.counter++
		token := fmt.Sprintf("catalog-fixture-%d", c.counter)
		c.principals[token] = *p
		c.mu.Unlock()
		req.Header.Set("Authorization", "Bearer "+token)
	}
	response, err := c.client.Do(req)
	if err != nil {
		return nil, 0, "", err
	}
	defer response.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(response.Body, (4<<20)+1))
	if err != nil {
		return nil, 0, "", err
	}
	if len(raw) > 4<<20 {
		return nil, 0, "", errors.New("oversize fixture response")
	}
	if p != nil && response.Header.Get("Cache-Control") != "no-store" {
		return nil, 0, "", errors.New("private response cacheable")
	}
	if response.StatusCode >= 400 {
		return nil, response.StatusCode, "", fmt.Errorf("HTTP %d: %s", response.StatusCode, raw)
	}
	return raw, response.StatusCode, response.Header.Get("ETag"), nil
}
func (c *catalogHTTPFixture) call(ctx context.Context, p identity.Principal, method, path string, in, out any) error {
	var raw []byte
	var err error
	if in != nil {
		raw, err = json.Marshal(in)
		if err != nil {
			return err
		}
	}
	body, _, _, err := c.request(ctx, &p, method, path, "application/json", raw, "")
	if err != nil {
		return err
	}
	return json.Unmarshal(body, out)
}
func (c *catalogHTTPFixture) UploadPNG(ctx context.Context, p identity.Principal, key string, raw []byte) (cr.Receipt, error) {
	var out cr.Receipt
	body, _, _, err := c.request(ctx, &p, "POST", "/v1/admin/catalog/media/"+url.PathEscape(key), "image/png", raw, "")
	if err == nil {
		err = json.Unmarshal(body, &out)
	}
	return out, err
}
func (c *catalogHTTPFixture) CreateDraft(ctx context.Context, p identity.Principal, in cr.DraftRequest) (cr.Receipt, error) {
	var out cr.Receipt
	err := c.call(ctx, p, "POST", "/v1/admin/catalog/drafts", in, &out)
	return out, err
}
func (c *catalogHTTPFixture) Draft(ctx context.Context, p identity.Principal, id string) (cr.Draft, error) {
	var out cr.Draft
	err := c.call(ctx, p, "GET", "/v1/admin/catalog/drafts/"+url.PathEscape(id), nil, &out)
	return out, err
}
func (c *catalogHTTPFixture) Review(ctx context.Context, p identity.Principal, in cr.ReviewRequest) (cr.Receipt, error) {
	var out cr.Receipt
	err := c.call(ctx, p, "POST", "/v1/admin/catalog/drafts/"+url.PathEscape(in.DraftID)+"/review/"+in.Stage, in, &out)
	return out, err
}
func (c *catalogHTTPFixture) CommandReceipt(ctx context.Context, p identity.Principal, id string) (cr.Receipt, error) {
	var out cr.Receipt
	err := c.call(ctx, p, "GET", "/v1/admin/catalog/commands/"+url.PathEscape(id), nil, &out)
	return out, err
}
func (c *catalogHTTPFixture) Publish(ctx context.Context, p identity.Principal, in cr.PublishRequest) (cr.Receipt, error) {
	c.mu.Lock()
	c.publisher = p
	c.mu.Unlock()
	var out cr.Receipt
	err := c.call(ctx, p, "POST", "/v1/admin/catalog/drafts/"+url.PathEscape(in.DraftID)+"/publish", in, &out)
	if err == nil {
		return out, nil
	}
	var network *url.Error
	if !errors.As(err, &network) {
		return out, err
	}
	got, getErr := c.CommandReceipt(ctx, p, in.CommandID)
	if getErr != nil {
		return out, getErr
	}
	_, hash, hashErr := cr.Canonical(in)
	if hashErr != nil || got.RequestSHA256 != hash || got.Actor != p.Subject || got.CommandID != in.CommandID || got.ResourceID != in.DraftID || got.SnapshotSHA256 != in.SnapshotSHA256 {
		return out, errors.New("unconfirmed publication identity")
	}
	c.mu.Lock()
	c.recoveries++
	c.mu.Unlock()
	return got, nil
}
func (c *catalogHTTPFixture) Public(ctx context.Context) (cr.Publication, error) {
	raw, status, etag, err := c.request(ctx, nil, "GET", "/v1/public/catalog", "", nil, c.etag)
	if err != nil {
		return cr.Publication{}, err
	}
	if status == 304 {
		raw = c.cached
	} else {
		if c.etag != "" {
			if etag == c.etag {
				return cr.Publication{}, errors.New("new version did not invalidate cache")
			}
			c.cacheChanges++
		}
		c.etag = etag
		c.cached = raw
	}
	_, sameStatus, _, err := c.request(ctx, nil, "GET", "/v1/public/catalog", "", nil, etag)
	if err != nil || sameStatus != 304 {
		return cr.Publication{}, fmt.Errorf("unchanged ETag %d %v", sameStatus, err)
	}
	feed, _, _, err := c.request(ctx, nil, "GET", "/v1/public/catalog/feed", "", nil, "")
	if err != nil || !bytes.Equal(raw, feed) {
		return cr.Publication{}, fmt.Errorf("feed differs from storefront: %v", err)
	}
	if bytes.Contains(raw, []byte("tenant_id")) || bytes.Contains(raw, []byte("organization_id")) {
		return cr.Publication{}, errors.New("private stream profile leaked")
	}
	var document cr.PublicDocument
	if err = json.Unmarshal(raw, &document); err != nil {
		return cr.Publication{}, err
	}
	if c.feedProbe != nil {
		c.feedProbe(ctx, c.publisher, document.Generation, raw)
	}
	if c.webProbe != nil {
		c.webProbe(document)
	}
	return cr.Publication{Generation: document.Generation, SHA256: document.SourceSHA256, EffectivePriceBookID: document.EffectivePriceBookID}, nil
}
func (c *catalogHTTPFixture) PublicMedia(ctx context.Context, sha string) ([]byte, error) {
	raw, _, _, err := c.request(ctx, nil, "GET", "/v1/public/catalog/media/"+sha, "", nil, "")
	return raw, err
}
func (c *catalogHTTPFixture) PublicSearch(ctx context.Context, q string) ([]cr.SearchHit, error) {
	raw, _, _, err := c.request(ctx, nil, "GET", "/v1/public/catalog/search?q="+url.QueryEscape(q), "", nil, "")
	if err != nil {
		return nil, err
	}
	var out []cr.SearchHit
	err = json.Unmarshal(raw, &out)
	return out, err
}
func TestCatalogReleaseHTTPReference(t *testing.T) { testCatalogReleaseHTTPReference(t, false) }
func TestCatalogReleaseFeedReference(t *testing.T) { testCatalogReleaseHTTPReference(t, true) }
func testCatalogReleaseHTTPReference(t *testing.T, withFeed bool) {
	var client *catalogHTTPFixture
	testCatalogReleaseConnected(t, func(t *testing.T, store *db.CatalogRelease) catalogReleaseReference {
		client = &catalogHTTPFixture{t: t, principals: map[string]identity.Principal{}, accepted: map[string]int{}}
		mux := http.NewServeMux()
		var feeder httpapi.CatalogFeedService
		if withFeed {
			feeder = attachCatalogFeedFixture(t, client, store)
		}
		httpapi.CatalogReleaseModule{Service: store, Feed: feeder}.Register(mux, client)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var cmd cr.PublishRequest
			if r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/publish") {
				raw, _ := io.ReadAll(io.LimitReader(r.Body, 32769))
				r.Body.Close()
				r.Body = io.NopCloser(bytes.NewReader(raw))
				_ = json.Unmarshal(raw, &cmd)
			}
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, r)
			drop := false
			if cmd.CommandID != "" && (rec.Code == 200 || rec.Code == 201) {
				client.mu.Lock()
				client.accepted[cmd.CommandID]++
				if cmd.CommandID == "publish-two" && rec.Code == 201 && !client.dropped {
					client.dropped = true
					drop = true
				}
				client.mu.Unlock()
			}
			if drop {
				conn, _, err := w.(http.Hijacker).Hijack()
				if err != nil {
					t.Error(err)
					return
				}
				conn.Close()
				return
			}
			for k, values := range rec.Header() {
				for _, value := range values {
					w.Header().Add(k, value)
				}
			}
			w.WriteHeader(rec.Code)
			_, _ = w.Write(rec.Body.Bytes())
		}))
		transport := &http.Transport{DisableKeepAlives: true}
		client.client = &http.Client{Transport: transport, Timeout: 5 * time.Second}
		client.base = server.URL
		t.Cleanup(func() { transport.CloseIdleConnections(); server.Close() })
		if catalogStorefrontFixtureRequested() {
			attachCatalogStorefrontFixture(t, client)
		}
		return client
	})
	if client.recoveries != 1 || !client.dropped || client.accepted["publish-two"] != 1 || client.cacheChanges != 3 {
		t.Fatal("HTTP recovery/cache coverage", client.recoveries, client.accepted, client.cacheChanges)
	}
	t.Log("CATALOG_RELEASE_HTTP_PASS connected role commands; lost committed response recovered by GET without repeat POST; storefront/feed byte-identical; three stale ETags invalidated; normalized media and scoped FTS match current publication")
}
