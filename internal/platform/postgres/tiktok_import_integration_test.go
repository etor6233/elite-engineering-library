package postgres

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"testing"
	"time"

	"elite.local/enterprise/internal/leadstream"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestTikTokImportUsesSharedDurableOwnerAndReplays(t *testing.T) {
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
	tenant := leadstream.StableUUID(t.Name(), time.Now().UTC().Format(time.RFC3339Nano))
	organization := "store-1"
	if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'TikTok Import Test','TikTok Import Test')`, tenant, "tiktok-"+tenant[:8]); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,$2,$2,'TikTok Store','store')`, tenant, organization); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_, _ = pool.Exec(context.Background(), `update platform.outbox_event set published_at=greatest(clock_timestamp(),occurred_at),claimed_by=null,claimed_until=null where tenant_id=$1 and published_at is null`, tenant)
	}()
	provider, batch, receipt := tiktokIntegrationArtifacts(t, tenant, organization)
	store := NewLeadIngress(pool)
	first, err := leadstream.ImportTikTokEvidence(ctx, store, provider, batch, receipt, time.Now().UTC())
	if err != nil || first.State != leadstream.ReceiptNormalized || !first.Complete {
		t.Fatalf("first=%+v err=%v", first, err)
	}
	replay, err := leadstream.ImportTikTokEvidence(ctx, store, provider, batch, receipt, time.Now().UTC())
	if err != nil || replay.State != leadstream.ReceiptDuplicate {
		t.Fatalf("replay=%+v err=%v", replay, err)
	}
	var rawCount, candidateCount, outboxCount int
	if err = pool.QueryRow(ctx, `select count(*) from integration.lead_ingress_raw where tenant_id=$1 and provider='tiktok_lead_generation' and provider_event_id='lead-1'`, tenant).Scan(&rawCount); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from integration.lead_candidate where tenant_id=$1 and provider='tiktok_lead_generation' and provider_event_id='lead-1' and contact_eligibility='pending_policy'`, tenant).Scan(&candidateCount); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id='tiktok_lead_generation:lead-1'`, tenant).Scan(&outboxCount); err != nil {
		t.Fatal(err)
	}
	if rawCount != 1 || candidateCount != 1 || outboxCount != 1 {
		t.Fatalf("raw=%d candidate=%d outbox=%d", rawCount, candidateCount, outboxCount)
	}
}

func tiktokIntegrationArtifacts(t *testing.T, tenant, organization string) ([]byte, []byte, []byte) {
	t.Helper()
	hash := func(value []byte) string {
		sum := sha256.Sum256(value)
		return hex.EncodeToString(sum[:])
	}
	providerObject := map[string]any{"schema": "elite-tiktok-lead-sdk-response/v1", "sdk_version": "1.1.3", "api_version": "v1.3", "request_identity": map[string]any{"mode": "TEST", "lead_source": "INSTANT_FORM", "advertiser_id_sha256": hash([]byte("123")), "library_id_sha256": "", "page_id_sha256": hash([]byte("456"))}, "data": map[string]any{"lead_data": map[string]any{"Email": "person@example.test"}, "meta_data": map[string]any{"lead_source": "INSTANT_FORM", "lead_id": "lead-1", "page_id": "456", "campaign_id": "11", "adgroup_id": "22", "ad_id": "33", "create_time": "2026-09-04 12:00:00"}}, "request_id": "req-1"}
	provider, _ := json.Marshal(providerObject)
	provider = append(provider, '\n')
	candidate := map[string]any{"tenant_id": tenant, "organization_id": organization, "provider": "tiktok_lead_generation", "provider_lead_id": "lead-1", "form_id": "456", "campaign_id": "11", "ad_group_id": "22", "creative_id": "33", "source_kind": "INSTANT_FORM", "submitted_at": "2026-09-04T12:00:00Z", "is_test": true, "fields": []map[string]string{{"id": "Email", "value": "person@example.test"}}, "contact_eligibility": "pending_policy"}
	batchObject := map[string]any{"schema": "elite-tiktok-lead-candidate-batch/v1", "provider": "tiktok_lead_generation", "sdk_version": "1.1.3", "sdk_wheel_sha256": "663b4a1f2585c4b386144d968646f695befd733c464418bbaeeac41befff33b7", "api_version": "v1.3", "contract_observed_at": "2026-09-04", "tenant_id": tenant, "organization_id": organization, "provider_response_sha256": hash(provider), "outcomes": []any{map[string]any{"status": "CANDIDATE", "normalization_codes": []string{}, "provider_lead_id_sha256": hash([]byte("lead-1")), "candidate": candidate}}, "automatic_business_write": false, "automatic_contact_eligibility": false}
	batch, _ := json.Marshal(batchObject)
	batch = append(batch, '\n')
	receiptObject := map[string]any{"schema": "elite-tiktok-lead-retrieval-receipt/v1", "sdk_version": "1.1.3", "sdk_wheel_sha256": "663b4a1f2585c4b386144d968646f695befd733c464418bbaeeac41befff33b7", "api_version": "v1.3", "contract_observed_at": "2026-09-04", "retrieval_mode": "TEST", "request_id_sha256": hash([]byte("req-1")), "provider_response_sha256": hash(provider), "candidate_batch_sha256": hash(batch), "row_count": 1, "candidate_count": 1, "rejected_count": 0, "automatic_business_write": false, "automatic_contact_eligibility": false}
	receipt, _ := json.Marshal(receiptObject)
	return provider, batch, append(receipt, '\n')
}
