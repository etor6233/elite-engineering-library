package postgres

// AUTHORED delta tests; isolated PostgreSQL fixtures, no external provider.
import (
	"bytes"
	"context"
	"crypto/sha256"
	"elite.local/enterprise/internal/businesspolicy"
	"elite.local/enterprise/internal/commerce"
	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/platform/randomid"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"strings"
	"testing"
	"time"
)

func policyFixture(t *testing.T, old, new string) *businesspolicy.Profile {
	t.Helper()
	raw := bytes.ReplaceAll(businesspolicy.ReferenceJSON(), []byte(old), []byte(new))
	p, err := businesspolicy.Load(raw, fmt.Sprintf("%x", sha256.Sum256(raw)))
	if err != nil {
		t.Fatal(err)
	}
	return p
}
func TestBusinessPolicyPersistenceBindingsAndReplay(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires disposable loopback database")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := randomid.Generator{}.New()
	tenantCode := "policy-" + strings.ReplaceAll(tenant, "-", "")
	for _, q := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'policy-'||replace(($1::uuid)::text,'-',''),'Synthetic','Synthetic')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic','store')`,
		`insert into org.public_location(tenant_id,organization_id,city,region,country,published)values($1,'store','Synthetic','Synthetic','AR',true)`,
		`insert into crm.lead(tenant_id,lead_id,organization_id,lifecycle_state,source_code,contact_payload)values($1,'lead','store','new','fixture','{}')`,
	} {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	p := policyFixture(t, `"lead_time_seconds": 1800`, `"lead_time_seconds": 0`)
	repo, err := NewFranchiseJourneyWithProfile(pool, p)
	if err != nil {
		t.Fatal(err)
	}
	strict := policyFixture(t, `"maximum_slot_capacity": 100`, `"maximum_slot_capacity": 1`)
	strictRepo, _ := NewFranchiseJourneyWithProfile(pool, strict)
	if repo.BusinessPolicySHA256() != p.SHA256() {
		t.Fatal("repository policy binding")
	}
	if _, err = NewFranchiseJourneyWithProfile(pool, nil); err == nil {
		t.Fatal("invalid policy accepted")
	}
	start := time.Now().UTC().Add(10 * time.Minute).Truncate(time.Second)
	end := start.Add(time.Hour)
	if _, err = repo.CreateAvailability(ctx, tenant, "scheduler", franchisejourney.AvailabilityEntry{ID: "working", OrganizationID: "store", EntryType: "working", StartsAt: start, EndsAt: end.Add(4 * time.Hour)}, randomid.Generator{}.New()); err != nil {
		t.Fatal(err)
	}
	slot := franchisejourney.AppointmentSlot{ID: "slot", OrganizationID: "store", Kind: "consultation", StartsAt: start, EndsAt: end, Capacity: 2}
	if _, err = NewFranchiseJourney(pool).CreateAppointmentSlot(ctx, tenant, slot, randomid.Generator{}.New()); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("repository bypassed reference lead time", err)
	}
	tooMany := slot
	tooMany.ID = "too-many"
	tooMany.StartsAt = end
	tooMany.EndsAt = end.Add(time.Hour)
	if _, err = strictRepo.CreateAppointmentSlot(ctx, tenant, tooMany, randomid.Generator{}.New()); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("repository bypassed configured capacity", err)
	}
	short := policyFixture(t, `"maximum_slot_seconds": 28800`, `"maximum_slot_seconds": 3600`)
	shortRepo, _ := NewFranchiseJourneyWithProfile(pool, short)
	long := tooMany
	long.ID = "too-long"
	long.EndsAt = long.StartsAt.Add(2 * time.Hour)
	if _, err = shortRepo.CreateAppointmentSlot(ctx, tenant, long, randomid.Generator{}.New()); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("repository bypassed duration", err)
	}
	hash := strings.Repeat("a", 64)
	event := randomid.Generator{}.New()
	first, replay, err := repo.CreateAppointmentSlotOnce(ctx, tenant, "scheduler", "slot-key", hash, slot, event)
	if err != nil || replay || first.ID != slot.ID {
		t.Fatal(first, replay, err)
	}
	again, replay, err := repo.CreateAppointmentSlotOnce(ctx, tenant, "scheduler", "slot-key", hash, slot, randomid.Generator{}.New())
	if err != nil || !replay || again.ID != first.ID {
		t.Fatal("same policy replay", again, replay, err)
	}
	if _, _, err = strictRepo.CreateAppointmentSlotOnce(ctx, tenant, "scheduler", "slot-key", hash, slot, randomid.Generator{}.New()); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("changed policy replay accepted", err)
	}
	if items, e := repo.PublicAppointmentSlots(ctx, tenantCode, "store", "consultation", start.Add(-time.Minute), end); e != nil || len(items) != 1 {
		t.Fatal("custom policy listing ignored", len(items), e)
	}
	if items, e := NewFranchiseJourney(pool).PublicAppointmentSlots(ctx, tenantCode, "store", "consultation", start.Add(-time.Minute), end); e != nil || len(items) != 0 {
		t.Fatal("unbookable reference slot advertised", len(items), e)
	}
	appointment := franchisejourney.Appointment{ID: "appointment", LeadID: "lead", Kind: "consultation", StartsAt: start}
	if _, _, err = NewFranchiseJourney(pool).RequestAppointment(ctx, tenantCode, "store", "reference-booking", appointment, hash, randomid.Generator{}.New()); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("direct booking bypassed reference lead time", err)
	}
	booked, replay, err := repo.RequestAppointment(ctx, tenantCode, "store", "booking-key", appointment, hash, randomid.Generator{}.New())
	if err != nil || replay || booked.ID != appointment.ID {
		t.Fatal(booked, replay, err)
	}
	booked, replay, err = repo.RequestAppointment(ctx, tenantCode, "store", "booking-key", appointment, hash, randomid.Generator{}.New())
	if err != nil || !replay || booked.ID != appointment.ID {
		t.Fatal("booking replay", booked, replay, err)
	}
	if _, _, err = strictRepo.RequestAppointment(ctx, tenantCode, "store", "booking-key", appointment, hash, randomid.Generator{}.New()); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("booking policy mismatch", err)
	}
	var count int
	var bound, receipt string
	if err = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and event_type in ('appointment-slot.created','appointment.requested')`, tenant).Scan(&count); err != nil || count != 2 {
		t.Fatal("duplicate or partial effects", count, err)
	}
	if err = pool.QueryRow(ctx, `select payload->>'policy_sha256' from platform.outbox_event where tenant_id=$1 and event_id=$2`, tenant, event).Scan(&bound); err != nil || bound != p.SHA256() {
		t.Fatal("event hash binding", bound, err)
	}
	expected, _ := p.BindRequestHash(hash)
	if err = pool.QueryRow(ctx, `select request_sha256_hex from platform.idempotency_record where tenant_id=$1 and scope='franchise-slot' and idempotency_key='slot-key'`, tenant).Scan(&receipt); err != nil || receipt != expected || receipt == hash {
		t.Fatal("receipt hash binding", receipt, err)
	}
	// Existing records are prospective: lowering maximum capacity does not rewrite slots.
	if err = pool.QueryRow(ctx, `select capacity from crm.appointment_slot where tenant_id=$1 and slot_id='slot'`, tenant).Scan(&count); err != nil || count != 2 {
		t.Fatal("configuration rewrote persisted quota")
	}
	// Narrow Commerce delta: exact profile binding in activation audit; no pricing/order suite rerun.
	cr, err := NewCommerceWithProfile(pool, p)
	if err != nil || cr.BusinessPolicySHA256() != p.SHA256() {
		t.Fatal("commerce binding", err)
	}
	if _, err = NewCommerceWithProfile(pool, nil); err == nil {
		t.Fatal("commerce invalid profile")
	}
	for _, q := range []string{
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values($1,'model','model','Synthetic','bicycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'variant','model','variant','Synthetic','{}','active')`,
	} {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	b := commerce.PriceBook{ID: "book", Market: "AR", Currency: "ARS", ValidFrom: time.Now().Add(-time.Hour), Status: "draft", Entries: []commerce.PriceEntry{{VariantID: "variant", AmountMinorUnits: 1000, TaxMode: "inclusive"}}}
	if err = cr.CreatePriceBook(ctx, tenant, randomid.Generator{}.New(), b); err != nil {
		t.Fatal(err)
	}
	activation := randomid.Generator{}.New()
	if err = cr.ActivatePriceBook(ctx, tenant, b.ID, activation); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select payload->>'policy_sha256' from platform.outbox_event where tenant_id=$1 and event_id=$2`, tenant, activation).Scan(&bound); err != nil || bound != p.SHA256() {
		t.Fatal("pricebook policy audit", bound, err)
	}
}
