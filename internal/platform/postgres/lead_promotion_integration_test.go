package postgres

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/leadpromotion"
	"elite.local/enterprise/internal/leadstream"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestLeadPromotionRequiresEvidenceAndIsAtomicAndIdempotent(t *testing.T) {
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
	organization := "promotion-store"
	if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Promotion Test','Promotion Test')`, tenant, "promotion-"+tenant[:8]); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,$2,$2,'Promotion Store','store')`, tenant, organization); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_, _ = pool.Exec(context.Background(), `update platform.outbox_event set published_at=greatest(clock_timestamp(),occurred_at),claimed_by=null,claimed_until=null where tenant_id=$1 and published_at is null`, tenant)
	}()

	repo := NewLeadIngress(pool)
	when := time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC)
	raw, candidate := promotionCandidate(t, tenant, organization, "provider-event-1", false, when)
	if _, err = repo.Record(ctx, raw, candidate, ""); err != nil {
		t.Fatal(err)
	}
	command := leadpromotion.Command{TenantID: tenant, Provider: "google_ads", ProviderEventID: "provider-event-1", LeadID: "crm-lead-1", ConsentID: "consent-1", PurposeCode: "sales-contact", PolicyVersion: "policy-approved-7", EvidenceSHA256: strings.Repeat("a", 64), DecisionAt: when.Add(time.Minute), MappingVersion: "google-form-42-v3", FieldMapping: map[string]string{"EMAIL": "email", "FULL_NAME": "name"}, IdempotencyKey: "promotion-request-0001"}
	promoter := NewLeadPromotion(pool)
	receipt, err := promoter.Promote(ctx, command)
	if err != nil || receipt.Replayed || receipt.LeadID != "crm-lead-1" {
		t.Fatalf("first promotion=%+v err=%v", receipt, err)
	}
	receipt, err = promoter.Promote(ctx, command)
	if err != nil || !receipt.Replayed {
		t.Fatalf("replay=%+v err=%v", receipt, err)
	}
	changed := command
	changed.PolicyVersion = "different"
	if _, err = promoter.Promote(ctx, changed); !errors.Is(err, leadpromotion.ErrConflict) {
		t.Fatalf("expected conflicting replay, got %v", err)
	}
	var leads, consents, promotions, events int
	if err = pool.QueryRow(ctx, `select count(*) from crm.lead where tenant_id=$1 and lead_id='crm-lead-1' and contact_payload->>'email'='person@example.test'`, tenant).Scan(&leads); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from crm.consent_evidence where tenant_id=$1 and consent_id='consent-1' and decision='granted' and policy_version='policy-approved-7'`, tenant).Scan(&consents); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from integration.lead_promotion where tenant_id=$1 and provider_event_id='provider-event-1' and source_payload_sha256=$2`, tenant, raw.SourceSHA256).Scan(&promotions); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and event_type='lead.promoted' and aggregate_id='crm-lead-1'`, tenant).Scan(&events); err != nil {
		t.Fatal(err)
	}
	if leads != 1 || consents != 1 || promotions != 1 || events != 1 {
		t.Fatalf("leads=%d consents=%d promotions=%d events=%d", leads, consents, promotions, events)
	}

	testRaw, testCandidate := promotionCandidate(t, tenant, organization, "provider-event-test", true, when)
	if _, err = repo.Record(ctx, testRaw, testCandidate, ""); err != nil {
		t.Fatal(err)
	}
	testCommand := command
	testCommand.ProviderEventID = "provider-event-test"
	testCommand.LeadID = "must-not-exist"
	testCommand.ConsentID = "must-not-exist"
	testCommand.IdempotencyKey = "promotion-request-test"
	if _, err = promoter.Promote(ctx, testCommand); !errors.Is(err, leadpromotion.ErrTestLead) {
		t.Fatalf("test lead promotion=%v", err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from crm.lead where tenant_id=$1 and lead_id='must-not-exist'`, tenant).Scan(&leads); err != nil || leads != 0 {
		t.Fatalf("test lead persisted=%d err=%v", leads, err)
	}
}

func promotionCandidate(t *testing.T, tenant, organization, eventID string, isTest bool, when time.Time) (leadstream.RawEvent, *leadstream.LeadCandidate) {
	t.Helper()
	payload := []byte(`{"lead_id":"` + eventID + `"}`)
	raw, err := leadstream.NewRawEvent(tenant, organization, "google_ads", eventID, "google.ads.lead.v3", "https://googleads.googleapis.com/lead-form", "3", when, when.Add(time.Second), payload)
	if err != nil {
		t.Fatal(err)
	}
	return raw, &leadstream.LeadCandidate{TenantID: tenant, OrganizationID: organization, Provider: "google_ads", ProviderLeadID: eventID, FormID: "42", SubmittedAt: when, IsTest: isTest, Fields: []leadstream.Field{{ID: "EMAIL", Value: "person@example.test"}, {ID: "FULL_NAME", Value: "Ada Lovelace"}}, ContactEligibility: "pending_policy"}
}
