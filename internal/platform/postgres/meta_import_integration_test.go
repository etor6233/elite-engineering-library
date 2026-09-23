package postgres

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"os"
	"testing"
	"time"

	"elite.local/enterprise/internal/leadstream"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestMetaEvidenceImportDurabilityRejectionAndReplay(t *testing.T) {
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
	const organization = "meta-store"
	if _, err := pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Meta Import Test','Meta Import Test')`, tenant, "meta-"+tenant[:8]); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,$2,$2,'Meta Store','store')`, tenant, organization); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_, _ = pool.Exec(context.Background(), `update platform.outbox_event set published_at=greatest(clock_timestamp(),occurred_at),claimed_by=null,claimed_until=null where tenant_id=$1 and published_at is null`, tenant)
	}()

	provider, batch, receipt := postgresMetaArtifacts(t, tenant, organization)
	repository := NewLeadIngress(pool)
	first, err := leadstream.ImportMetaEvidence(ctx, repository, provider, batch, receipt, time.Now().UTC())
	if err != nil || !first.Complete || first.Normalized != 1 || first.Rejected != 1 {
		t.Fatalf("first=%+v err=%v", first, err)
	}
	replay, err := leadstream.ImportMetaEvidence(ctx, repository, provider, batch, receipt, time.Now().UTC())
	if err != nil || replay.Duplicates != 2 || !replay.Complete {
		t.Fatalf("replay=%+v err=%v", replay, err)
	}

	var rawCount, candidateCount, rejectedCount, outboxCount int
	if err := pool.QueryRow(ctx, `select count(*) from integration.lead_ingress_raw where tenant_id=$1 and provider='meta_lead_ads' and provider_event_id in ('meta-live-1','meta-rejected-1')`, tenant).Scan(&rawCount); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, `select count(*) from integration.lead_candidate where tenant_id=$1 and provider='meta_lead_ads' and provider_event_id='meta-live-1' and contact_eligibility='pending_policy'`, tenant).Scan(&candidateCount); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, `select count(*) from integration.lead_ingress_raw where tenant_id=$1 and provider='meta_lead_ads' and provider_event_id='meta-rejected-1' and state='rejected' and normalization_error_code='UNAPPROVED_FORM_FIELD'`, tenant).Scan(&rejectedCount); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id in ('meta_lead_ads:meta-live-1','meta_lead_ads:meta-rejected-1')`, tenant).Scan(&outboxCount); err != nil {
		t.Fatal(err)
	}
	if rawCount != 2 || candidateCount != 1 || rejectedCount != 1 || outboxCount != 2 {
		t.Fatalf("raw=%d candidate=%d rejected=%d outbox=%d", rawCount, candidateCount, rejectedCount, outboxCount)
	}
}

func postgresMetaArtifacts(t *testing.T, tenant, organization string) ([]byte, []byte, []byte) {
	t.Helper()
	rows := []map[string]any{
		{"id": "meta-live-1", "created_time": "2026-09-04T12:00:00+0000", "form_id": "123456789", "campaign_id": "campaign-1", "adset_id": "adset-1", "ad_id": "ad-1", "platform": "facebook", "field_data": []map[string]any{{"name": "email", "values": []string{"person@example.test"}}}},
		{"id": "meta-rejected-1", "created_time": "2026-09-04T12:01:00+0000", "form_id": "123456789", "campaign_id": "campaign-1", "adset_id": "adset-1", "ad_id": "ad-1", "platform": "facebook", "field_data": []map[string]any{{"name": "unapproved", "values": []string{"retained-only-in-raw"}}}},
	}
	provider := append(mustJSONIndent(t, map[string]any{"graph_api_version": "v26.0", "retrieval_mode": "LIVE", "rows": rows, "truncated_by_local_bound": false}), '\n')
	providerHash := testSHA(provider)
	outcomes := []map[string]any{
		{"status": "CANDIDATE", "normalization_codes": []string{}, "provider_lead_id_sha256": testSHA([]byte("meta-live-1")), "candidate": map[string]any{"tenant_id": tenant, "organization_id": organization, "provider": "meta_lead_ads", "provider_lead_id": "meta-live-1", "form_id": "123456789", "campaign_id": "campaign-1", "ad_group_id": "adset-1", "creative_id": "ad-1", "source_kind": "facebook", "submitted_at": "2026-09-04T12:00:00+0000", "is_test": false, "fields": []map[string]string{{"id": "email", "value": "person@example.test"}}, "contact_eligibility": "pending_policy"}},
		{"status": "REJECTED", "normalization_codes": []string{"UNAPPROVED_FORM_FIELD"}, "provider_lead_id_sha256": testSHA([]byte("meta-rejected-1")), "candidate": nil},
	}
	batch := append(mustJSONIndent(t, map[string]any{"schema": "elite-meta-lead-candidate-batch/v1", "provider": "meta_lead_ads", "source_commit": "788f363d15b1269ab5efb7cd00fb5e3b133cd99b", "tenant_id": tenant, "organization_id": organization, "provider_response_sha256": providerHash, "outcomes": outcomes, "automatic_business_write": false, "automatic_contact_eligibility": false}), '\n')
	receipt := append(mustJSONIndent(t, map[string]any{"schema": "elite-meta-lead-reconciliation-receipt/v1", "sdk_version": "26.0.1", "graph_api_version": "v26.0", "source_commit": "788f363d15b1269ab5efb7cd00fb5e3b133cd99b", "form_id_sha256": testSHA([]byte("123456789")), "provider_response_sha256": providerHash, "candidate_batch_sha256": testSHA(batch), "row_count": 2, "candidate_count": 1, "rejected_count": 1, "truncated_by_local_bound": false, "automatic_business_write": false, "automatic_contact_eligibility": false}), '\n')
	return provider, batch, receipt
}

func mustJSONIndent(t *testing.T, value any) []byte {
	t.Helper()
	encoded, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	return encoded
}

func testSHA(value []byte) string {
	const digits = "0123456789abcdef"
	return sha256Hex(value, digits)
}

func sha256Hex(value []byte, digits string) string {
	h := sha256.Sum256(value)
	encoded := make([]byte, len(h)*2)
	for i, b := range h {
		encoded[i*2], encoded[i*2+1] = digits[b>>4], digits[b&15]
	}
	return string(encoded)
}
