package postgres

import (
	"context"
	"errors"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/contactidentity"
	"elite.local/enterprise/internal/leadstream"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestContactIdentityEvidenceVersioningResolutionAndRevocation(t *testing.T) {
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
	seed := time.Now().UTC().Format(time.RFC3339Nano)
	tenant := leadstream.StableUUID(t.Name(), seed, "tenant")
	otherTenant := leadstream.StableUUID(t.Name(), seed, "other")
	for _, id := range []string{tenant, otherTenant} {
		if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Identity Test','Identity Test')`, id, "identity-"+id[len(id)-4:]); err != nil {
			t.Fatal(err)
		}
		if _, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'store-1','store-1','Store 1','store')`, id); err != nil {
			t.Fatal(err)
		}
		if _, err = pool.Exec(ctx, `insert into crm.lead(tenant_id,lead_id,organization_id,lifecycle_state,source_code,contact_payload,consent_required) values($1,'lead-1','store-1','new','meta','{}',true)`, id); err != nil {
			t.Fatal(err)
		}
	}
	key := []byte("0123456789abcdef0123456789abcdef")
	store, err := NewContactIdentityStore(pool, key)
	if err != nil {
		t.Fatal(err)
	}
	when := time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC)
	command := contactidentity.Command{TenantID: tenant, ChannelCode: "whatsapp", ExternalID: "+5491112345678", LeadID: "lead-1", SubjectID: "lead:lead-1", PIIAllowed: true, State: contactidentity.StateActive, PolicyVersion: "sales-contact-v1", EvidenceSHA256: strings.Repeat("a", 64), EffectiveAt: when, RequestID: "binding-1"}
	receipt, err := store.Apply(ctx, command)
	if err != nil || receipt.Version != 1 || receipt.Replayed {
		t.Fatalf("receipt=%+v err=%v", receipt, err)
	}
	replay, err := store.Apply(ctx, command)
	if err != nil || !replay.Replayed || replay.Version != 1 {
		t.Fatalf("replay=%+v err=%v", replay, err)
	}
	changed := command
	changed.PolicyVersion = "changed"
	if _, err = store.Apply(ctx, changed); !errors.Is(err, contactidentity.ErrConflict) {
		t.Fatalf("expected replay conflict, got %v", err)
	}
	message := channels.Message{ChannelCode: "whatsapp", TenantID: tenant, ExternalID: command.ExternalID, ProviderMessageID: "wamid.1", OccurredAt: when, Direction: channels.DirectionIn, Text: "hola"}
	contact, err := store.Resolve(ctx, message)
	if err != nil || contact.Scope.OrganizationID != "store-1" || contact.Scope.LeadID != "lead-1" || contact.SubjectID != "lead:lead-1" || !contact.PIIAllowed {
		t.Fatalf("contact=%+v err=%v", contact, err)
	}
	message.TenantID = otherTenant
	if _, err = store.Resolve(ctx, message); !errors.Is(err, contactidentity.ErrNotFound) {
		t.Fatalf("cross-tenant resolve=%v", err)
	}
	command.ExpectedVersion, command.RequestID, command.PolicyVersion, command.EvidenceSHA256, command.PIIAllowed = 1, "binding-2", "sales-contact-v2", strings.Repeat("b", 64), false
	receipt, err = store.Apply(ctx, command)
	if err != nil || receipt.Version != 2 {
		t.Fatalf("policy update=%+v err=%v", receipt, err)
	}
	message.TenantID = tenant
	contact, err = store.Resolve(ctx, message)
	if err != nil || contact.PIIAllowed {
		t.Fatalf("updated contact=%+v err=%v", contact, err)
	}
	command.ExpectedVersion, command.RequestID, command.State, command.PolicyVersion, command.EvidenceSHA256 = 2, "binding-3", contactidentity.StateRevoked, "sales-contact-revoked-v1", strings.Repeat("c", 64)
	receipt, err = store.Apply(ctx, command)
	if err != nil || receipt.Version != 3 || receipt.State != contactidentity.StateRevoked {
		t.Fatalf("revoke=%+v err=%v", receipt, err)
	}
	if _, err = store.Resolve(ctx, message); !errors.Is(err, contactidentity.ErrNotFound) {
		t.Fatalf("revoked resolve=%v", err)
	}
	var decisions int
	var leaked bool
	if err = pool.QueryRow(ctx, `select count(*),bool_or(external_id_hmac like '%12345678%') from communication.contact_channel_binding_decision where tenant_id=$1`, tenant).Scan(&decisions, &leaked); err != nil {
		t.Fatal(err)
	}
	if decisions != 3 || leaked {
		t.Fatalf("decisions=%d leaked=%v", decisions, leaked)
	}
	concurrent := contactidentity.Command{TenantID: tenant, ChannelCode: "sms", ExternalID: "+5491187654321", LeadID: "lead-1", SubjectID: "lead:lead-1", State: contactidentity.StateActive, PolicyVersion: "sales-contact-v1", EvidenceSHA256: strings.Repeat("d", 64), EffectiveAt: when, ExpectedVersion: 0}
	results := make(chan error, 2)
	var start sync.WaitGroup
	start.Add(2)
	for i := 0; i < 2; i++ {
		go func(i int) {
			defer start.Done()
			copy := concurrent
			copy.RequestID = []string{"concurrent-a", "concurrent-b"}[i]
			_, applyErr := store.Apply(context.Background(), copy)
			results <- applyErr
		}(i)
	}
	start.Wait()
	close(results)
	succeeded, conflicted := 0, 0
	for applyErr := range results {
		switch {
		case applyErr == nil:
			succeeded++
		case errors.Is(applyErr, contactidentity.ErrConflict):
			conflicted++
		default:
			t.Fatalf("unexpected concurrent result: %v", applyErr)
		}
	}
	if succeeded != 1 || conflicted != 1 {
		t.Fatalf("concurrent succeeded=%d conflicted=%d", succeeded, conflicted)
	}
}
