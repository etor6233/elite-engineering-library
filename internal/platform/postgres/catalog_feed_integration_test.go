package postgres_test

// AUTHORED receiver protocol fixture, not a live Google/Meta/ML API claim.
// Provider-specific SDK mapping remains the explicit T2805 adapter owner.
import (
	"bytes"
	"context"
	"crypto/sha256"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func attachCatalogFeedFixture(t *testing.T, client *catalogHTTPFixture, store *db.CatalogRelease) httpapi.CatalogFeedService {
	t.Helper()
	var mu sync.Mutex
	posts := map[int64]int{}
	gets := map[int64]int{}
	expected := map[int64][]byte{}
	acks := map[string]cr.FeedAck{}
	last := int64(0)
	receiver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer fixture-feed-token" {
			w.WriteHeader(401)
			return
		}
		key := strings.TrimPrefix(r.URL.Path, "/publications/")
		generation, e := strconv.ParseInt(r.Header.Get("X-Catalog-Generation"), 10, 64)
		if e != nil || key == "" || key == r.URL.Path {
			w.WriteHeader(400)
			return
		}
		mu.Lock()
		defer mu.Unlock()
		if r.Method == "GET" {
			gets[generation]++
			ack, ok := acks[key]
			if !ok {
				w.WriteHeader(404)
				return
			}
			_ = json.NewEncoder(w).Encode(ack)
			return
		}
		if r.Method != "POST" {
			w.WriteHeader(405)
			return
		}
		raw, e := io.ReadAll(io.LimitReader(r.Body, 524289))
		if e != nil || len(raw) > 524288 {
			w.WriteHeader(413)
			return
		}
		sum := sha256.Sum256(raw)
		hash := hex.EncodeToString(sum[:])
		if !bytes.Equal(raw, expected[generation]) || r.Header.Get("X-Content-SHA256") != hash || r.Header.Get("Idempotency-Key") != key {
			t.Error("feed bytes/header differ from published source")
			w.WriteHeader(400)
			return
		}
		posts[generation]++
		if posts[generation] != 1 {
			t.Error("duplicate provider POST", generation)
		}
		ack := cr.FeedAck{Generation: generation, PayloadSHA256: hash, RemoteID: fmt.Sprintf("fixture-catalog-%d", generation), AcceptedAt: time.Now().UTC()}
		if generation == 1 || generation < last {
			ack.Status = "rejected"
			acks[key] = ack
			w.WriteHeader(422)
			_ = json.NewEncoder(w).Encode(ack)
			return
		}
		ack.Status = "accepted"
		last = generation
		acks[key] = ack
		if generation == 2 {
			conn, _, e := w.(http.Hijacker).Hijack()
			if e != nil {
				t.Error(e)
				return
			}
			conn.Close()
			return
		}
		w.WriteHeader(201)
		_ = json.NewEncoder(w).Encode(ack)
	}))
	t.Cleanup(receiver.Close)
	sender, e := cr.NewFeedHTTP(cr.FeedProfile{Schema: "elite-catalog-feed/v1", ReceiverID: "fixture-receiver", BaseURL: receiver.URL, Mode: "fixture"}, "fixture-feed-token")
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(sender.Close)
	feeder, e := db.NewCatalogFeed(store, bytes.Repeat([]byte{'f'}, 32), sender)
	if e != nil {
		t.Fatal(e)
	}
	client.feedProbe = func(ctx context.Context, p identity.Principal, generation int64, raw []byte) {
		t.Helper()
		mu.Lock()
		expected[generation] = append([]byte(nil), raw...)
		mu.Unlock()
		path := fmt.Sprintf("/v1/admin/catalog/feed/%d", generation)
		invoke := func(method, suffix string) (cr.FeedStatus, error) {
			var out cr.FeedStatus
			data, _, _, e := client.request(ctx, &p, method, path+suffix, "", nil, "")
			if e == nil {
				e = json.Unmarshal(data, &out)
			}
			return out, e
		}
		_, e := invoke("POST", "")
		if generation == 1 || generation == 2 {
			if e == nil {
				t.Fatal("rejected/unknown feed reported accepted", generation)
			}
		} else if e != nil {
			t.Fatal("feed", e)
		}
		status, e := invoke("GET", "")
		if e != nil {
			t.Fatal(e)
		}
		wanted := "accepted"
		if generation == 1 {
			wanted = "failed_terminal"
		}
		if generation == 2 {
			wanted = "unknown"
		}
		digest := sha256.Sum256(raw)
		if status.State != wanted || status.Generation != generation || status.PayloadSHA256 != hex.EncodeToString(digest[:]) {
			t.Fatal("feed durable status", status, wanted)
		}
		// Repeating a known terminal or uncertain send cannot issue another POST.
		_, again := invoke("POST", "")
		if (generation == 1 || generation == 2) && again == nil {
			t.Fatal("terminal/unknown retry not fenced")
		}
		if generation == 2 {
			status, e = invoke("POST", "/reconcile")
			if e != nil || status.State != "accepted" {
				t.Fatal("GET reconciliation", status, e)
			}
			if status, e = invoke("POST", "/reconcile"); e != nil || status.State != "accepted" {
				t.Fatal("reconciled replay", e)
			}
		}
	}
	t.Cleanup(func() {
		mu.Lock()
		defer mu.Unlock()
		if len(posts) != 4 || posts[1] != 1 || posts[2] != 1 || posts[3] != 1 || posts[5] != 1 || gets[2] != 1 || last != 5 {
			t.Error("feed coverage", posts, gets, last)
		}
		t.Log("CATALOG_FEED_FENCE_PASS four current publications; exact storefront bytes; rejection terminal; accepted response lost then GET reconciled; no repeated provider POST; retained PG intent/fence")
	})
	return feeder
}
