package postgres

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/leadstream"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestMetaLeadWebhookDurabilityReplayDivergenceAndTenantIsolation(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	tenantA := leadstream.StableUUID(t.Name(), "a", time.Now().UTC().Format(time.RFC3339Nano))
	tenantB := leadstream.StableUUID(t.Name(), "b", time.Now().UTC().Format(time.RFC3339Nano))
	for index, tenant := range []string{tenantA, tenantB} {
		code := "meta-" + tenant[:8]
		if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,$3,$3)`, tenant, code, "Meta Test "+string(rune('A'+index))); err != nil {
			t.Fatal(err)
		}
		if _, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'store-1','store-1','Store 1','store')`, tenant); err != nil {
			t.Fatal(err)
		}
	}
	if _, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'store-2','store-2','Store 2','store')`, tenantA); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_, _ = pool.Exec(context.Background(), `update platform.outbox_event set published_at=greatest(clock_timestamp(),occurred_at),claimed_by=null,claimed_until=null where tenant_id=any($1::uuid[]) and published_at is null`, []string{tenantA, tenantB})
	}()

	payload := `{"object":"page","entry":[{"id":"page-1","time":1630087927,"changes":[{"value":{"form_id":"form-1","leadgen_id":"lead-1","created_time":1630087926,"page_id":"page-1"},"field":"leadgen"}]}]}`
	sign := func(value string) string {
		mac := hmac.New(sha256.New, []byte("secret"))
		_, _ = mac.Write([]byte(value))
		return "sha256=" + hex.EncodeToString(mac.Sum(nil))
	}
	decode := func(tenant, value string) leadstream.MetaWebhookBatch {
		batch, decodeErr := leadstream.DecodeMetaLeadWebhook([]byte(value), sign(value), "secret", tenant, "store-1", time.Now().UTC())
		if decodeErr != nil {
			t.Fatal(decodeErr)
		}
		return batch
	}
	store := NewMetaLeadWebhookStore(pool)
	first, err := store.RecordMetaWebhook(ctx, decode(tenantA, payload))
	if err != nil || first.NewSignals != 1 {
		t.Fatalf("first=%+v err=%v", first, err)
	}
	replay, err := store.RecordMetaWebhook(ctx, decode(tenantA, payload))
	if err != nil || replay.DuplicateSignals != 1 {
		t.Fatalf("replay=%+v err=%v", replay, err)
	}
	other, err := store.RecordMetaWebhook(ctx, decode(tenantB, payload))
	if err != nil || other.NewSignals != 1 {
		t.Fatalf("tenant B=%+v err=%v", other, err)
	}
	changed := strings.Replace(payload, `"form_id":"form-1"`, `"form_id":"form-2"`, 1)
	if _, err = store.RecordMetaWebhook(ctx, decode(tenantA, changed)); !errors.Is(err, leadstream.ErrDivergentDuplicate) {
		t.Fatalf("expected divergence, got %v", err)
	}
	crossOrganization := decode(tenantA, payload)
	crossOrganization.OrganizationID = "store-2"
	if _, err = store.RecordMetaWebhook(ctx, crossOrganization); !errors.Is(err, leadstream.ErrDivergentDuplicate) {
		t.Fatalf("expected organization divergence, got %v", err)
	}
	var signals, batches, outbox int
	if err = pool.QueryRow(ctx, `select count(*) from integration.meta_lead_webhook_signal where tenant_id=$1 and leadgen_id='lead-1' and state='pending_retrieval'`, tenantA).Scan(&signals); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from integration.meta_lead_webhook_batch where tenant_id=$1`, tenantA).Scan(&batches); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and event_type='meta-lead.signal-received'`, tenantA).Scan(&outbox); err != nil {
		t.Fatal(err)
	}
	if signals != 1 || batches != 1 || outbox != 1 {
		t.Fatalf("signals=%d batches=%d outbox=%d", signals, batches, outbox)
	}
}
