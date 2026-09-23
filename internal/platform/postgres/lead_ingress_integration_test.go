package postgres

import (
	"context"
	"errors"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/leadstream"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestLeadIngressDurabilityReplayDivergenceAndConcurrency(t *testing.T) {
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
	tenantCode := "lead-" + tenant[:8]
	const organization = "lead-store"
	if _, err := pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)
		values($1,$2,'Lead Runtime Test','Lead Runtime Test')`, tenant, tenantCode); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)
		values($1,$2,'lead-store','Lead Store','store')`, tenant, organization); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_, _ = pool.Exec(context.Background(), `update platform.outbox_event
			set published_at=greatest(clock_timestamp(),occurred_at),claimed_by=null,claimed_until=null
			where tenant_id=$1 and published_at is null`, tenant)
	}()

	repo := NewLeadIngress(pool)
	when := time.Date(2020, 9, 4, 18, 0, 0, 0, time.UTC)
	makeLead := func(eventID string, payload []byte) (leadstream.RawEvent, *leadstream.LeadCandidate) {
		raw, makeErr := leadstream.NewRawEvent(tenant, organization, "google_ads", eventID, "google.ads.lead.v3", "https://googleads.googleapis.com/lead-form", "3", when, when.Add(time.Second), payload)
		if makeErr != nil {
			t.Fatal(makeErr)
		}
		return raw, &leadstream.LeadCandidate{
			TenantID: tenant, OrganizationID: organization, Provider: "google_ads", ProviderLeadID: eventID,
			FormID: "9223372036854775806", CampaignID: "42", SubmittedAt: when, IsTest: true,
			Fields: []leadstream.Field{{ID: "EMAIL", Value: "fixture@example.invalid"}}, ContactEligibility: "pending_policy",
		}
	}

	googlePayload := []byte(`{"lead_id":"lead-db-1","user_column_data":[{"column_id":"EMAIL","string_value":"fixture@example.invalid"}],"api_version":"3","google_key":"integration-secret","is_test":true,"lead_submit_time":"2020-09-04T18:00:00Z"}`)
	decoded, err := leadstream.DecodeGoogleWebhook(googlePayload, "integration-secret", tenant, organization, when.Add(time.Second))
	if err != nil || decoded.Candidate == nil {
		t.Fatalf("decode candidate=%+v err=%v", decoded.Candidate, err)
	}
	raw, candidate := decoded.Raw, decoded.Candidate
	receipt, err := repo.Record(ctx, raw, candidate, "")
	if err != nil || receipt.State != leadstream.ReceiptNormalized {
		t.Fatalf("first record state=%q err=%v", receipt.State, err)
	}
	receipt, err = repo.Record(ctx, raw, candidate, "")
	if err != nil || receipt.State != leadstream.ReceiptDuplicate {
		t.Fatalf("exact replay state=%q err=%v", receipt.State, err)
	}
	changedPayload := []byte(`{"lead_id":"lead-db-1","user_column_data":[{"column_id":"EMAIL","string_value":"changed@example.invalid"}],"api_version":"3","google_key":"integration-secret","is_test":true,"lead_submit_time":"2020-09-04T18:00:00Z"}`)
	changed, decodeErr := leadstream.DecodeGoogleWebhook(changedPayload, "integration-secret", tenant, organization, when.Add(time.Second))
	if decodeErr != nil {
		t.Fatal(decodeErr)
	}
	if _, err := repo.Record(ctx, changed.Raw, changed.Candidate, ""); !errors.Is(err, leadstream.ErrDivergentDuplicate) {
		t.Fatalf("expected divergent duplicate, got %v", err)
	}
	var rawCount, candidateCount, outboxCount int
	var sourceHash, storedHash, storedPayload, redactionProfile string
	if err := pool.QueryRow(ctx, `select count(*) from integration.lead_ingress_raw where tenant_id=$1 and provider_event_id='lead-db-1'`, tenant).Scan(&rawCount); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, `select count(*) from integration.lead_candidate where tenant_id=$1 and provider_event_id='lead-db-1' and contact_eligibility='pending_policy'`, tenant).Scan(&candidateCount); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id='google_ads:lead-db-1'`, tenant).Scan(&outboxCount); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, `select source_payload_sha256,stored_payload_sha256,convert_from(payload_redacted,'UTF8'),redaction_profile
		from integration.lead_ingress_raw where tenant_id=$1 and provider='google_ads' and provider_event_id='lead-db-1'`, tenant).Scan(&sourceHash, &storedHash, &storedPayload, &redactionProfile); err != nil {
		t.Fatal(err)
	}
	if rawCount != 1 || candidateCount != 1 || outboxCount != 1 {
		t.Fatalf("raw=%d candidate=%d outbox=%d", rawCount, candidateCount, outboxCount)
	}
	if sourceHash != raw.SourceSHA256 || storedHash != raw.StoredSHA256 || redactionProfile != raw.RedactionProfile ||
		strings.Contains(storedPayload, "integration-secret") || strings.Contains(storedPayload, "google_key") {
		t.Fatalf("unsafe or inconsistent stored evidence source=%s stored=%s profile=%q payload=%q", sourceHash, storedHash, redactionProfile, storedPayload)
	}

	rejectedRaw, makeErr := leadstream.NewRawEvent(tenant, organization, "google_ads", "reject-db-1", "google.ads.lead.v3", "https://googleads.googleapis.com/lead-form", "3", when, when.Add(time.Second), []byte(`{"authenticated":"but-unmappable"}`))
	if makeErr != nil {
		t.Fatal(makeErr)
	}
	receipt, err = repo.Record(ctx, rejectedRaw, nil, "UNMAPPABLE_FIXTURE")
	if err != nil || receipt.State != leadstream.ReceiptRejected {
		t.Fatalf("rejection state=%q err=%v", receipt.State, err)
	}
	if err := pool.QueryRow(ctx, `select count(*) from integration.lead_candidate where tenant_id=$1 and provider_event_id='reject-db-1'`, tenant).Scan(&candidateCount); err != nil || candidateCount != 0 {
		t.Fatalf("rejected candidate count=%d err=%v", candidateCount, err)
	}

	concurrentRaw, concurrentCandidate := makeLead("lead-db-concurrent", []byte(`{"lead_id":"lead-db-concurrent"}`))
	const workers = 16
	states := make(chan leadstream.ReceiptState, workers)
	errs := make(chan error, workers)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			got, recordErr := repo.Record(ctx, concurrentRaw, concurrentCandidate, "")
			states <- got.State
			errs <- recordErr
		}()
	}
	wg.Wait()
	close(states)
	close(errs)
	for recordErr := range errs {
		if recordErr != nil {
			t.Fatal(recordErr)
		}
	}
	normalized, duplicates := 0, 0
	for state := range states {
		switch state {
		case leadstream.ReceiptNormalized:
			normalized++
		case leadstream.ReceiptDuplicate:
			duplicates++
		default:
			t.Fatalf("unexpected concurrent state %q", state)
		}
	}
	if normalized != 1 || duplicates != workers-1 {
		t.Fatalf("normalized=%d duplicates=%d", normalized, duplicates)
	}
	if err := pool.QueryRow(ctx, `select count(*) from integration.lead_ingress_raw where tenant_id=$1 and provider_event_id='lead-db-concurrent'`, tenant).Scan(&rawCount); err != nil || rawCount != 1 {
		t.Fatalf("concurrent raw count=%d err=%v", rawCount, err)
	}
}
