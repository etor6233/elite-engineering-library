package postgres

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/platform/randomid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type journeyReadProbeKey struct{}

func TestBookingRequiresCurrentAvailability(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires disposable loopback database")
	}
	for _, mode := range []string{"cancelled", "unavailable", "race-cancel", "race-unavailable"} {
		t.Run(mode, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			pool, err := pgxpool.NewWithConfig(ctx, cfg.Copy())
			if err != nil {
				t.Fatal(err)
			}
			defer pool.Close()
			tenant := randomid.Generator{}.New()
			tenantCode := "booking-" + strings.ReplaceAll(tenant, "-", "")
			for _, q := range []string{
				`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'booking-'||replace(($1::uuid)::text,'-',''),'Synthetic','Synthetic')`,
				`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic','store')`,
				`insert into org.public_location(tenant_id,organization_id,city,region,country,published)values($1,'store','Synthetic','Synthetic','AR',true)`,
				`insert into crm.lead(tenant_id,lead_id,organization_id,lifecycle_state,source_code,contact_payload)values($1,'lead','store','new','fixture','{}')`,
			} {
				if _, err = pool.Exec(ctx, q, tenant); err != nil {
					t.Fatal(err)
				}
			}
			repo := NewFranchiseJourney(pool)
			start := time.Now().UTC().Add(72 * time.Hour).Truncate(time.Second)
			end := start.Add(time.Hour)
			window := franchisejourney.AvailabilityEntry{ID: "working", OrganizationID: "store", EntryType: "working", StartsAt: start, EndsAt: end}
			if _, err = repo.CreateAvailability(ctx, tenant, "scheduler", window, randomid.Generator{}.New()); err != nil {
				t.Fatal(err)
			}
			if _, err = repo.CreateAppointmentSlot(ctx, tenant, franchisejourney.AppointmentSlot{ID: "slot", OrganizationID: "store", Kind: "consultation", StartsAt: start, EndsAt: end, Capacity: 1}, randomid.Generator{}.New()); err != nil {
				t.Fatal(err)
			}
			appointment := franchisejourney.Appointment{ID: "booking", LeadID: "lead", Kind: "consultation", StartsAt: start}
			book := func() error {
				_, _, e := repo.RequestAppointment(ctx, tenantCode, "store", "booking-key-v305", appointment, strings.Repeat("a", 64), randomid.Generator{}.New())
				return e
			}
			change := func() error {
				if mode == "cancelled" || mode == "race-cancel" {
					_, e := repo.CancelAvailability(ctx, tenant, "store", "working", 1, "scheduler", "schedule-correction", randomid.Generator{}.New())
					return e
				}
				window.ID, window.EntryType, window.ReasonCode = "absence", "unavailable", "synthetic-absence"
				_, e := repo.CreateAvailability(ctx, tenant, "scheduler", window, randomid.Generator{}.New())
				return e
			}
			if !strings.HasPrefix(mode, "race-") {
				if err = change(); err != nil {
					t.Fatal(err)
				}
				items, e := repo.PublicAppointmentSlots(ctx, tenantCode, "store", "consultation", start.Add(-time.Minute), end.Add(time.Minute))
				if e != nil || len(items) != 0 {
					t.Errorf("unavailable slot published count=%d err=%v", len(items), e)
				}
				if e = book(); !errors.Is(e, franchisejourney.ErrConflict) {
					t.Errorf("booking outside availability accepted: %v", e)
				}
				var effects int
				if e = pool.QueryRow(ctx, `select (select count(*) from crm.appointment where tenant_id=$1)+(select count(*) from platform.idempotency_record where tenant_id=$1 and scope='public-appointment')+(select count(*) from platform.outbox_event where tenant_id=$1 and event_type='appointment.requested')`, tenant).Scan(&effects); e != nil || effects != 0 {
					t.Errorf("unexpected booking effects=%d err=%v", effects, e)
				}
			} else {
				startTogether := make(chan struct{})
				results := make(chan error, 2)
				go func() { <-startTogether; results <- book() }()
				go func() { <-startTogether; results <- change() }()
				close(startTogether)
				success, conflicts := 0, 0
				for i := 0; i < 2; i++ {
					e := <-results
					if e == nil {
						success++
					} else if errors.Is(e, franchisejourney.ErrConflict) {
						conflicts++
					} else {
						t.Fatal(e)
					}
				}
				if success != 1 || conflicts != 1 {
					t.Fatalf("race success=%d conflicts=%d", success, conflicts)
				}
				var booked, eligible int
				if err = pool.QueryRow(ctx, `select (select count(*) from crm.appointment where tenant_id=$1),(select count(*) from crm.availability_entry where tenant_id=$1 and availability_id='working' and state='active')-(select count(*) from crm.availability_entry where tenant_id=$1 and entry_type='unavailable' and state='active')`, tenant).Scan(&booked, &eligible); err != nil || booked != eligible {
					t.Fatal("contradictory schedule", booked, eligible, err)
				}
				if booked == 1 {
					replayed, wasReplay, e := repo.RequestAppointment(ctx, tenantCode, "store", "booking-key-v305", appointment, strings.Repeat("a", 64), randomid.Generator{}.New())
					if e != nil || !wasReplay || replayed.ID != "booking" {
						t.Fatal("replay changed", replayed, wasReplay, e)
					}
				}
			}
		})
	}
}

func TestAvailabilityCancellationAtomicity(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires disposable loopback database")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := randomid.Generator{}.New()
	for _, q := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'avail-'||replace(($1::uuid)::text,'-',''),'Synthetic','Synthetic')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic','store')`,
		`insert into crm.lead(tenant_id,lead_id,organization_id,lifecycle_state,source_code,contact_payload)values($1,'lead','store','new','fixture','{}')`,
	} {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewFranchiseJourney(pool)
	start := time.Now().UTC().Add(72 * time.Hour).Truncate(time.Second)
	creationEvent := randomid.Generator{}.New()
	if _, err = repo.CreateAvailability(ctx, tenant, "creator", franchisejourney.AvailabilityEntry{ID: "atomic-window", OrganizationID: "store", EntryType: "working", StartsAt: start, EndsAt: start.Add(time.Hour)}, creationEvent); err != nil {
		t.Fatal(err)
	}
	if _, err = repo.CancelAvailability(ctx, tenant, "other", "atomic-window", 1, "operator", "schedule-correction", randomid.Generator{}.New()); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("foreign organization", err)
	}
	if _, err = repo.CancelAvailability(ctx, randomid.Generator{}.New(), "store", "atomic-window", 1, "operator", "schedule-correction", randomid.Generator{}.New()); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("foreign tenant", err)
	}
	if _, err = repo.CancelAvailability(ctx, tenant, "store", "atomic-window", 2, "operator", "schedule-correction", randomid.Generator{}.New()); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("stale version", err)
	}
	_, err = repo.CancelAvailability(ctx, tenant, "store", "atomic-window", 1, "operator", "schedule-correction", creationEvent)
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != "23505" {
		t.Fatal("expected duplicate audit event", err)
	}
	var state, actor string
	var version, events int
	if err = pool.QueryRow(ctx, `select state,version,coalesce(cancelled_by_subject,'') from crm.availability_entry where tenant_id=$1 and availability_id='atomic-window'`, tenant).Scan(&state, &version, &actor); err != nil || state != "active" || version != 1 || actor != "" {
		t.Fatal("partial cancellation", state, version, actor, err)
	}
	results := make(chan error, 2)
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, e := repo.CancelAvailability(ctx, tenant, "store", "atomic-window", 1, "operator", "schedule-correction", randomid.Generator{}.New())
			results <- e
		}()
	}
	wg.Wait()
	close(results)
	success, conflicts := 0, 0
	for e := range results {
		if e == nil {
			success++
		} else if errors.Is(e, franchisejourney.ErrConflict) {
			conflicts++
		} else {
			t.Fatal(e)
		}
	}
	if success != 1 || conflicts != 1 {
		t.Fatal("race", success, conflicts)
	}
	pool.Close()
	pool, err = pgxpool.NewWithConfig(ctx, cfg.Copy())
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	repo = NewFranchiseJourney(pool)
	items, err := repo.Availability(ctx, tenant, "store", "", start.Add(-time.Minute), start.Add(2*time.Hour))
	if err != nil || len(items) != 1 || items[0].ID != "atomic-window" || items[0].State != "cancelled" || items[0].Version != 2 {
		t.Fatal("reconnected read", items, err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id='atomic-window' and event_type='availability-entry.cancelled' and payload->>'actor_subject'='operator'`, tenant).Scan(&events); err != nil || events != 1 {
		t.Fatal("audit", events, err)
	}
	// Existing cancellation policy: an active appointment prevents removing its working interval.
	start = start.Add(24 * time.Hour)
	if _, err = repo.CreateAvailability(ctx, tenant, "creator", franchisejourney.AvailabilityEntry{ID: "booked-window", OrganizationID: "store", EntryType: "working", StartsAt: start, EndsAt: start.Add(2 * time.Hour)}, randomid.Generator{}.New()); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `insert into crm.appointment(tenant_id,appointment_id,organization_id,lead_id,appointment_kind,starts_at,state,version)values($1,'booked','store','lead','consultation',$2,'requested',1)`, tenant, start.Add(30*time.Minute)); err != nil {
		t.Fatal(err)
	}
	if _, err = repo.CancelAvailability(ctx, tenant, "store", "booked-window", 1, "operator", "schedule-correction", randomid.Generator{}.New()); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("active appointment guard", err)
	}
	if err = pool.QueryRow(ctx, `select state,version,(select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id='booked-window' and event_type='availability-entry.cancelled') from crm.availability_entry where tenant_id=$1 and availability_id='booked-window'`, tenant).Scan(&state, &version, &events); err != nil || state != "active" || version != 1 || events != 0 {
		t.Fatal("booked interval changed", state, version, events, err)
	}
}

func TestLeadAuditedAtomicity(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires disposable loopback database")
	}
	ctx := context.Background()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := randomid.Generator{}.New()
	eventOne, eventTwo := randomid.Generator{}.New(), randomid.Generator{}.New()
	for _, q := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'lead-'||replace(($1::uuid)::text,'-',''),'Synthetic','Synthetic')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic','store')`,
		`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized)values($1,'customer','Synthetic','fixture@example.invalid')`,
		`insert into crm.lead(tenant_id,lead_id,organization_id,customer_principal_id,lifecycle_state,source_code,contact_payload)values($1,'atomic-lead','store','customer','new','fixture','{}')`,
	} {
		if _, err := pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewFranchiseJourney(pool)
	if _, err := repo.AssignLeadAs(ctx, tenant, "store", "atomic-lead", "assignee", 1, "empty", ""); !errors.Is(err, franchisejourney.ErrInvalid) {
		t.Fatal("missing actor accepted", err)
	}
	if _, err := repo.TransitionLeadAs(ctx, tenant, "store", "atomic-lead", "new", "contacted", 1, "empty", ""); !errors.Is(err, franchisejourney.ErrInvalid) {
		t.Fatal("missing actor transition accepted", err)
	}
	if _, err := repo.AssignLeadAs(ctx, tenant, "store", "atomic-lead", "assignee", 1, eventOne, "actor-one"); err != nil {
		t.Fatal(err)
	}
	// Duplicate event ID must roll back the preceding update, not merely its audit.
	if _, err := repo.AssignLeadAs(ctx, tenant, "store", "atomic-lead", "other", 2, eventOne, "actor-two"); err == nil {
		t.Fatal("duplicate event assignment accepted")
	}
	if _, err := repo.TransitionLeadAs(ctx, tenant, "store", "atomic-lead", "new", "contacted", 2, eventOne, "actor-two"); err == nil {
		t.Fatal("duplicate event transition accepted")
	}
	var state, assignee string
	var version, events int
	if err := pool.QueryRow(ctx, `select lifecycle_state,assigned_subject,version from crm.lead where tenant_id=$1 and lead_id='atomic-lead'`, tenant).Scan(&state, &assignee, &version); err != nil || state != "new" || assignee != "assignee" || version != 2 {
		t.Fatal("partial update", state, assignee, version, err)
	}
	if _, err := repo.TransitionLeadAs(ctx, tenant, "store", "atomic-lead", "new", "contacted", 2, eventTwo, "actor-two"); err != nil {
		t.Fatal(err)
	}
	pool.Close()
	pool, err = pgxpool.NewWithConfig(ctx, cfg.Copy())
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	if err := pool.QueryRow(ctx, `select lifecycle_state,assigned_subject,version from crm.lead where tenant_id=$1 and lead_id='atomic-lead'`, tenant).Scan(&state, &assignee, &version); err != nil || state != "contacted" || assignee != "assignee" || version != 3 {
		t.Fatal("recovery", state, assignee, version, err)
	}
	if err := pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id='atomic-lead'`, tenant).Scan(&events); err != nil || events != 2 {
		t.Fatal("events", events, err)
	}
	for _, pair := range [][2]string{{eventOne, "actor-one"}, {eventTwo, "actor-two"}} {
		var actor string
		if err := pool.QueryRow(ctx, `select payload->>'actor_subject' from platform.outbox_event where tenant_id=$1 and event_id=$2`, tenant, pair[0]).Scan(&actor); err != nil || actor != pair[1] {
			t.Fatal("actor", actor, err)
		}
	}
}

type journeyReadProbe struct {
	afterFirstRead func() error
	once           sync.Once
	err            error
	inTransaction  bool
	boundedItems   int
}

func (p *journeyReadProbe) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	if strings.Contains(data.SQL, "h.handover_id=any($4::text[])") && len(data.Args) == 4 {
		if ids, ok := data.Args[3].([]string); ok {
			p.boundedItems = len(ids)
		}
	}
	return context.WithValue(ctx, journeyReadProbeKey{}, data.SQL)
}
func (p *journeyReadProbe) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, _ pgx.TraceQueryEndData) {
	query, _ := ctx.Value(journeyReadProbeKey{}).(string)
	if strings.HasPrefix(query, "select appointment_id,organization_id") {
		p.once.Do(func() { p.inTransaction = conn.PgConn().TxStatus() == 'T'; p.err = p.afterFirstRead() })
	}
}

func TestFranchiseJourneyPersistenceIsolationAndReplay(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c29b01"
	cleanup := func() {
		tx, cleanupErr := pool.Begin(ctx)
		if cleanupErr != nil {
			return
		}
		defer tx.Rollback(ctx)
		_, _ = tx.Exec(ctx, `set local session_replication_role=replica`)
		for _, q := range []string{`delete from sales.return_effect_attempt where tenant_id=$1`, `delete from sales.return_effect_resume where tenant_id=$1`, `delete from sales.return_effect_execution where tenant_id=$1`, `delete from sales.return_effect_request where tenant_id=$1`, `delete from sales.return_disposition where tenant_id=$1`, `delete from sales.return_receipt where tenant_id=$1`} {
			_, _ = tx.Exec(ctx, q, tenant)
		}
		for _, q := range []string{`delete from sales.delivery_exception_resolution where tenant_id=$1`, `delete from sales.return_authorization where tenant_id=$1`, `delete from sales.delivery_exception where tenant_id=$1`} {
			_, _ = tx.Exec(ctx, q, tenant)
		}
		for _, q := range []string{`delete from platform.outbox_event where tenant_id=$1`, `delete from platform.idempotency_record where tenant_id=$1`, `delete from sales.delivery_checklist_response where tenant_id=$1`, `delete from sales.delivery_handover where tenant_id=$1`, `delete from sales.delivery_checklist_item where tenant_id=$1`, `delete from sales.delivery_checklist_template where tenant_id=$1`, `delete from sales.quotation_acceptance where tenant_id=$1`, `delete from sales.quotation where tenant_id=$1`, `delete from sales.customer_order_line where tenant_id=$1`, `delete from sales.customer_order where tenant_id=$1`, `delete from crm.appointment_transition where tenant_id=$1`, `delete from crm.appointment_resource where tenant_id=$1`, `delete from crm.appointment where tenant_id=$1`, `delete from crm.appointment_slot where tenant_id=$1`, `delete from crm.availability_entry where tenant_id=$1`, `delete from crm.resource_skill where tenant_id=$1`, `delete from crm.service_resource where tenant_id=$1`, `delete from crm.consent_evidence where tenant_id=$1`, `delete from crm.lead where tenant_id=$1`, `delete from crm.customer_profile where tenant_id=$1`, `delete from org.public_location where tenant_id=$1`, `delete from inventory.stock_unit where tenant_id=$1`, `delete from pricing.price_book_entry where tenant_id=$1`, `delete from pricing.price_book where tenant_id=$1`, `delete from catalog.vehicle_variant where tenant_id=$1`, `delete from catalog.vehicle_model where tenant_id=$1`, `delete from org.organization where tenant_id=$1`, `delete from platform.tenant where tenant_id=$1`} {
			_, _ = tx.Exec(ctx, q, tenant)
		}
		_ = tx.Commit(ctx)
	}
	cleanup()
	defer cleanup()
	fixtures := []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'journey-api','Journey','Journey')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Store','store'),($1,'other','other','Other','store')`,
		`insert into org.public_location(tenant_id,organization_id,city,region,country,published)values($1,'store','Cordoba','Cordoba','AR',true)`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values($1,'model','model','Model','bicycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'variant','model','variant','Variant','{}','active')`,
		`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized)values($1,'customer','Customer','customer@example.test')`,
		`insert into crm.lead(tenant_id,lead_id,organization_id,customer_principal_id,model_id,lifecycle_state,source_code,contact_payload)values($1,'lead','store','customer','model','new','public-web','{}')`,
		`insert into pricing.price_book(tenant_id,price_book_id,market,currency,valid_from,status)values($1,'retail','AR','ARS',clock_timestamp()-interval '1 day','active')`,
		`insert into pricing.price_book_entry(tenant_id,price_book_id,variant_id,amount_minor_units,tax_mode)values($1,'retail','variant',123456,'inclusive')`,
		`insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,state,version,received_at)values($1,'stock','store','variant','SERIAL-JOURNEY','sold',1,clock_timestamp())`,
		`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'order','store','customer','delivered','ARS',123456,1)`,
		`insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version)values($1,'handover','store','order','customer','stock','prepared',1)`,
		`insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version)values($1,'handover-reject','store','order','customer','stock','prepared',1)`,
		`insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version)values($1,'handover-return','store','order','customer','stock','prepared',1)`,
	}
	for _, q := range fixtures {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewFranchiseJourney(pool)
	locations, err := repo.PublicLocations(ctx, "journey-api")
	if err != nil || len(locations) != 1 || locations[0].OrganizationID != "store" {
		t.Fatalf("locations=%+v err=%v", locations, err)
	}
	slotStart := time.Now().Add(time.Hour).Truncate(time.Second)
	slotEnd := slotStart.Add(45 * time.Minute)
	workingEnd := slotStart.Add(8 * time.Hour)
	organizationWorking, err := repo.CreateAvailability(ctx, tenant, "scheduler", franchisejourney.AvailabilityEntry{ID: "store-working", OrganizationID: "store", EntryType: "working", StartsAt: slotStart.Add(-time.Hour), EndsAt: workingEnd, State: "active", Version: 1}, "018f4d4a-7b36-7a21-8d10-2f4c54c29b20")
	if err != nil || organizationWorking.State != "active" {
		t.Fatalf("organization working window=%+v err=%v", organizationWorking, err)
	}
	slot, err := repo.CreateAppointmentSlot(ctx, tenant, franchisejourney.AppointmentSlot{ID: "test-drive-slot", OrganizationID: "store", Kind: "test-drive", StartsAt: slotStart, EndsAt: slotEnd, Capacity: 1, State: "open", Version: 1}, "018f4d4a-7b36-7a21-8d10-2f4c54c29b00")
	if err != nil || slot.Booked != 0 || slot.Capacity != 1 {
		t.Fatalf("slot=%+v err=%v", slot, err)
	}
	available, err := repo.PublicAppointmentSlots(ctx, "journey-api", "store", "test-drive", slotStart.Add(-time.Minute), slotEnd.Add(time.Minute))
	if err != nil || len(available) != 1 || available[0].ID != "test-drive-slot" {
		t.Fatalf("available=%+v err=%v", available, err)
	}
	if _, err = repo.CreateAppointmentSlot(ctx, tenant, franchisejourney.AppointmentSlot{ID: "overlap-slot", OrganizationID: "store", Kind: "test-drive", StartsAt: slotStart.Add(15 * time.Minute), EndsAt: slotEnd.Add(15 * time.Minute), Capacity: 1, State: "open", Version: 1}, "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatalf("overlapping slot accepted: %v", err)
	}
	appointment := franchisejourney.Appointment{ID: "appointment", LeadID: "lead", ModelID: "model", Kind: "test-drive", StartsAt: slotStart, State: "requested", Version: 1}
	created, replayed, err := repo.RequestAppointment(ctx, "journey-api", "store", "appointment-key-1", appointment, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "018f4d4a-7b36-7a21-8d10-2f4c54c29b02")
	if err != nil || replayed || created.OrganizationID != "store" {
		t.Fatalf("appointment=%+v replay=%v err=%v", created, replayed, err)
	}
	if created.SlotID != "test-drive-slot" || !created.EndsAt.Equal(slotEnd) {
		t.Fatalf("appointment was not server-bound to slot: %+v", created)
	}
	resource, err := repo.CreateServiceResource(ctx, tenant, franchisejourney.ServiceResource{ID: "technician", OrganizationID: "store", PrincipalSubject: "technician-subject", DisplayName: "Technician", Kind: "employee", Status: "active", Version: 1, Skills: []string{"test-drive", "service"}}, "018f4d4a-7b36-7a21-8d10-2f4c54c29b16")
	if err != nil || resource.Version != 1 {
		t.Fatalf("resource=%+v err=%v", resource, err)
	}
	if _, err = repo.CreateServiceResource(ctx, tenant, franchisejourney.ServiceResource{ID: "wrong-skill", OrganizationID: "store", DisplayName: "Service Bay", Kind: "service-bay", Status: "active", Version: 1, Skills: []string{"service"}}, "018f4d4a-7b36-7a21-8d10-2f4c54c29b17"); err != nil {
		t.Fatal(err)
	}
	resourceWorking, err := repo.CreateAvailability(ctx, tenant, "scheduler", franchisejourney.AvailabilityEntry{ID: "technician-working", OrganizationID: "store", ResourceID: "technician", EntryType: "working", StartsAt: slotStart.Add(-time.Hour), EndsAt: workingEnd, State: "active", Version: 1}, "018f4d4a-7b36-7a21-8d10-2f4c54c29b21")
	if err != nil || resourceWorking.ResourceID != "technician" {
		t.Fatalf("resource working window=%+v err=%v", resourceWorking, err)
	}
	if _, err = repo.AssignAppointmentResource(ctx, tenant, "store", "appointment", "wrong-skill", 1, "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatalf("wrong skill assignment accepted: %v", err)
	}
	assignedAppointment, err := repo.AssignAppointmentResource(ctx, tenant, "store", "appointment", "technician", 1, "018f4d4a-7b36-7a21-8d10-2f4c54c29b18")
	if err != nil || assignedAppointment.ResourceID != "technician" || assignedAppointment.Version != 2 {
		t.Fatalf("assigned appointment=%+v err=%v", assignedAppointment, err)
	}
	if _, err = pool.Exec(ctx, `insert into crm.appointment(tenant_id,appointment_id,organization_id,lead_id,customer_principal_id,model_id,appointment_kind,starts_at,ends_at,slot_id,state,version) values($1,'manual-overlap','store','lead','customer','model','test-drive',$2,$3,'test-drive-slot','requested',1)`, tenant, slotStart, slotEnd); err != nil {
		t.Fatal(err)
	}
	if _, err = repo.AssignAppointmentResource(ctx, tenant, "store", "manual-overlap", "technician", 1, "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatalf("overlapping resource assignment accepted: %v", err)
	}
	confirmed, err := repo.TransitionAppointment(ctx, tenant, "store", "appointment", "requested", "confirmed", 2, "operator", "", "018f4d4a-7b36-7a21-8d10-2f4c54c29b19")
	if err != nil || confirmed.State != "confirmed" || confirmed.ResourceID != "technician" || confirmed.Version != 3 {
		t.Fatalf("confirmed appointment=%+v err=%v", confirmed, err)
	}
	if _, err = repo.TransitionAppointment(ctx, tenant, "store", "appointment", "requested", "cancelled", 2, "operator", "customer-request", "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("stale appointment transition accepted")
	}
	if _, err = repo.TransitionAppointment(ctx, tenant, "store", "appointment", "confirmed", "completed", 3, "operator", "", "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("future appointment completion accepted")
	}
	if _, err = repo.CreateAvailability(ctx, tenant, "scheduler", franchisejourney.AvailabilityEntry{ID: "conflicting-leave", OrganizationID: "store", ResourceID: "technician", EntryType: "unavailable", ReasonCode: "annual-leave", StartsAt: slotStart, EndsAt: slotEnd, State: "active", Version: 1}, "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatalf("absence overlapping active appointment accepted: %v", err)
	}
	_, replayed, err = repo.RequestAppointment(ctx, "journey-api", "store", "appointment-key-1", appointment, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "unused")
	if err != nil || !replayed {
		t.Fatalf("replay=%v err=%v", replayed, err)
	}
	if _, _, err = repo.RequestAppointment(ctx, "journey-api", "store", "appointment-key-1", appointment, "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("changed replay accepted")
	}
	available, err = repo.PublicAppointmentSlots(ctx, "journey-api", "store", "test-drive", slotStart.Add(-time.Minute), slotEnd.Add(time.Minute))
	if err != nil || len(available) != 0 {
		t.Fatalf("full slot remained public: %+v err=%v", available, err)
	}
	second := appointment
	second.ID = "appointment-over-capacity"
	if _, _, err = repo.RequestAppointment(ctx, "journey-api", "store", "appointment-key-over-capacity", second, strings.Repeat("c", 64), "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatalf("over-capacity booking accepted: %v", err)
	}
	var partialIdempotency int
	if err = pool.QueryRow(ctx, `select count(*) from platform.idempotency_record where tenant_id=$1 and scope='public-appointment' and idempotency_key='appointment-key-over-capacity'`, tenant).Scan(&partialIdempotency); err != nil || partialIdempotency != 0 {
		t.Fatalf("capacity conflict left idempotency residue=%d err=%v", partialIdempotency, err)
	}
	cancelled, err := repo.TransitionAppointment(ctx, tenant, "store", "appointment", "confirmed", "cancelled", 3, "operator", "customer-request", "018f4d4a-7b36-7a21-8d10-2f4c54c29b1a")
	if err != nil || cancelled.State != "cancelled" || cancelled.Version != 4 {
		t.Fatalf("cancelled appointment=%+v err=%v", cancelled, err)
	}
	var transitionActor, transitionReason string
	if err = pool.QueryRow(ctx, `select actor_subject,reason_code from crm.appointment_transition where tenant_id=$1 and appointment_id='appointment' and to_state='cancelled'`, tenant).Scan(&transitionActor, &transitionReason); err != nil || transitionActor != "operator" || transitionReason != "customer-request" {
		t.Fatalf("transition actor=%s reason=%s err=%v", transitionActor, transitionReason, err)
	}
	concurrentStart := slotStart.Add(3 * time.Hour)
	if _, err = repo.CreateAppointmentSlot(ctx, tenant, franchisejourney.AppointmentSlot{ID: "concurrent-slot", OrganizationID: "store", Kind: "service", StartsAt: concurrentStart, EndsAt: concurrentStart.Add(time.Hour), Capacity: 1, State: "open", Version: 1}, "018f4d4a-7b36-7a21-8d10-2f4c54c29b13"); err != nil {
		t.Fatal(err)
	}
	type bookingResult struct {
		created franchisejourney.Appointment
		err     error
	}
	startBookings := make(chan struct{})
	results := make(chan bookingResult, 2)
	for _, input := range []struct{ appointment, key, hash, event string }{
		{"concurrent-a", "appointment-concurrent-a", strings.Repeat("1", 64), "018f4d4a-7b36-7a21-8d10-2f4c54c29b14"},
		{"concurrent-b", "appointment-concurrent-b", strings.Repeat("2", 64), "018f4d4a-7b36-7a21-8d10-2f4c54c29b15"},
	} {
		go func(input struct{ appointment, key, hash, event string }) {
			<-startBookings
			created, _, bookingErr := repo.RequestAppointment(ctx, "journey-api", "store", input.key, franchisejourney.Appointment{ID: input.appointment, LeadID: "lead", ModelID: "model", Kind: "service", StartsAt: concurrentStart, State: "requested", Version: 1}, input.hash, input.event)
			results <- bookingResult{created: created, err: bookingErr}
		}(input)
	}
	close(startBookings)
	successes, conflicts := 0, 0
	for range 2 {
		result := <-results
		if result.err == nil && result.created.SlotID == "concurrent-slot" {
			successes++
		} else if errors.Is(result.err, franchisejourney.ErrConflict) {
			conflicts++
		} else {
			t.Fatalf("unexpected concurrent booking result: %+v", result)
		}
	}
	if successes != 1 || conflicts != 1 {
		t.Fatalf("concurrent capacity successes=%d conflicts=%d", successes, conflicts)
	}
	customerCancelStart := slotStart.Add(6 * time.Hour)
	if _, err = repo.CreateAppointmentSlot(ctx, tenant, franchisejourney.AppointmentSlot{ID: "customer-cancel-slot", OrganizationID: "store", Kind: "service", StartsAt: customerCancelStart, EndsAt: customerCancelStart.Add(time.Hour), Capacity: 1, State: "open", Version: 1}, "018f4d4a-7b36-7a21-8d10-2f4c54c29b22"); err != nil {
		t.Fatal(err)
	}
	customerAppointment := franchisejourney.Appointment{ID: "customer-cancel-appointment", LeadID: "lead", ModelID: "model", Kind: "service", StartsAt: customerCancelStart, State: "requested", Version: 1}
	createdForCustomer, _, err := repo.RequestAppointment(ctx, "journey-api", "store", "appointment-customer-cancel", customerAppointment, strings.Repeat("4", 64), "018f4d4a-7b36-7a21-8d10-2f4c54c29b23")
	if err != nil || createdForCustomer.Version != 1 {
		t.Fatalf("customer appointment=%+v err=%v", createdForCustomer, err)
	}
	customerCancelled, err := repo.CancelCustomerAppointment(ctx, tenant, "store", "customer", "customer-cancel-appointment", 1, "customer-request", "018f4d4a-7b36-7a21-8d10-2f4c54c29b24", "018f4d4a-7b36-7a21-8d10-2f4c54c29b25")
	if err != nil || customerCancelled.State != "cancelled" || customerCancelled.Version != 2 {
		t.Fatalf("customer cancellation=%+v err=%v", customerCancelled, err)
	}
	if _, err = repo.CancelCustomerAppointment(ctx, tenant, "store", "other-customer", "customer-cancel-appointment", 2, "customer-request", "unused", "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatalf("cross-customer cancellation accepted: %v", err)
	}
	var concurrentBookings int
	if err = pool.QueryRow(ctx, `select count(*) from crm.appointment where tenant_id=$1 and slot_id='concurrent-slot' and state in ('requested','confirmed')`, tenant).Scan(&concurrentBookings); err != nil || concurrentBookings != 1 {
		t.Fatalf("concurrent persisted bookings=%d err=%v", concurrentBookings, err)
	}
	assigned, err := repo.AssignLead(ctx, tenant, "store", "lead", "sales", 1, "018f4d4a-7b36-7a21-8d10-2f4c54c29b03")
	if err != nil || assigned.Version != 2 {
		t.Fatalf("assigned=%+v err=%v", assigned, err)
	}
	if _, err = repo.AssignLead(ctx, tenant, "other", "lead", "intruder", 2, "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("cross organization assignment succeeded")
	}
	contacted, err := repo.TransitionLead(ctx, tenant, "store", "lead", "new", "contacted", 2, "018f4d4a-7b36-7a21-8d10-2f4c54c29b04")
	if err != nil || contacted.Version != 3 {
		t.Fatalf("contacted=%+v err=%v", contacted, err)
	}
	if _, err = repo.TransitionLead(ctx, tenant, "store", "lead", "new", "lost", 2, "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("stale transition succeeded")
	}
	quoteHash := strings.Repeat("5", 64)
	quote, replayedQuote, err := repo.CreateQuoteAs(ctx, tenant, "quote-request-0001", franchisejourney.Quote{ID: "quote", OrganizationID: "store", LeadID: "lead", VariantID: "variant", PriceBookID: "retail", ValidUntil: time.Now().Add(24 * time.Hour), State: "issued", Version: 1}, quoteHash, "018f4d4a-7b36-7a21-8d10-2f4c54c29b05", "writer")
	if err != nil || replayedQuote || quote.TotalMinorUnits != 123456 || quote.Currency != "ARS" {
		t.Fatalf("quote=%+v replayed=%v err=%v", quote, replayedQuote, err)
	}

	recoveredQuote, err := repo.QuoteResult(ctx, tenant, "store", "lead", "quote-request-0001")
	if err != nil || recoveredQuote.ID != "quote" || recoveredQuote.TotalMinorUnits != quote.TotalMinorUnits {
		t.Fatalf("quote result=%+v err=%v", recoveredQuote, err)
	}
	for _, scope := range [][3]string{{tenant, "other", "lead"}, {tenant, "store", "other"}, {randomid.Generator{}.New(), "store", "lead"}} {
		if v, e := repo.QuoteResult(ctx, scope[0], scope[1], scope[2], "quote-request-0001"); !errors.Is(e, franchisejourney.ErrConflict) || v.ID != "" {
			t.Fatalf("lookup scope leak=%+v err=%v", v, e)
		}
	}
	if v, e := repo.QuoteResult(ctx, tenant, "store", "lead", "unknown-key-00001"); !errors.Is(e, franchisejourney.ErrConflict) || v.ID != "" {
		t.Fatal("unknown quote identity accepted")
	}
	for _, mutation := range []string{"resource_type='other'", "response_code=200", "response_body='{\"quotation_id\":\"other\"}'"} {
		if _, e := pool.Exec(ctx, "update platform.idempotency_record set "+mutation+" where tenant_id=$1 and scope='franchise-quote' and idempotency_key='quote-request-0001'", tenant); e != nil {
			t.Fatal(e)
		}
		if v, e := repo.QuoteResult(ctx, tenant, "store", "lead", "quote-request-0001"); !errors.Is(e, franchisejourney.ErrConflict) || v.ID != "" {
			t.Fatalf("incoherent result=%+v err=%v", v, e)
		}
		if _, e := pool.Exec(ctx, `update platform.idempotency_record set resource_type='quotation',response_code=201,response_body='{"quotation_id":"quote"}' where tenant_id=$1 and scope='franchise-quote' and idempotency_key='quote-request-0001'`, tenant); e != nil {
			t.Fatal(e)
		}
	}
	replayedValue, replayedQuote, err := repo.CreateQuote(ctx, tenant, "quote-request-0001", franchisejourney.Quote{ID: "unused-replay-id", OrganizationID: "store", LeadID: "lead", VariantID: "variant", PriceBookID: "retail", ValidUntil: quote.ValidUntil, State: "issued", Version: 1}, quoteHash, "unused-replay-event")
	if err != nil || !replayedQuote || replayedValue.ID != quote.ID {
		t.Fatalf("quote replay=%+v replayed=%v err=%v", replayedValue, replayedQuote, err)
	}
	if _, _, err = repo.CreateQuote(ctx, tenant, "quote-request-0001", franchisejourney.Quote{ID: "conflicting-quote", OrganizationID: "store", LeadID: "lead", VariantID: "variant", PriceBookID: "retail", ValidUntil: quote.ValidUntil, State: "issued", Version: 1}, strings.Repeat("6", 64), "unused-conflict-event"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatalf("same key with different quote request accepted: %v", err)
	}

	var originalActor string
	if _, _, e := repo.CreateQuoteAs(ctx, tenant, "quote-request-0001", quote, quoteHash, "unused", "second-operator"); e != nil {
		t.Fatal(e)
	}
	if e := pool.QueryRow(ctx, `select payload->>'actor_subject' from platform.outbox_event where tenant_id=$1 and aggregate_id='quote' and event_type='quotation.issued'`, tenant).Scan(&originalActor); e != nil || originalActor != "writer" {
		t.Fatalf("replay changed actor=%s err=%v", originalActor, e)
	}
	attempted := quote
	attempted.ID = "quote-audit-rollback"
	if _, _, e := repo.CreateQuoteAs(ctx, tenant, "quote-atomicity-key", attempted, strings.Repeat("a", 64), "018f4d4a-7b36-7a21-8d10-2f4c54c29b05", "writer"); e == nil {
		t.Fatal("duplicate audit event accepted")
	}
	var residue int
	if e := pool.QueryRow(ctx, `select (select count(*) from sales.quotation where tenant_id=$1 and quotation_id='quote-audit-rollback')+(select count(*) from platform.idempotency_record where tenant_id=$1 and scope='franchise-quote' and idempotency_key='quote-atomicity-key')`, tenant).Scan(&residue); e != nil || residue != 0 {
		t.Fatalf("audit rollback residue=%d err=%v", residue, e)
	}
	var quoteRows, quoteEvents int
	if err = pool.QueryRow(ctx, `select count(*) from sales.quotation where tenant_id=$1 and quotation_id=$2`, tenant, quote.ID).Scan(&quoteRows); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id=$2 and event_type='quotation.issued'`, tenant, quote.ID).Scan(&quoteEvents); err != nil || quoteRows != 1 || quoteEvents != 1 {
		t.Fatalf("quote rows=%d events=%d err=%v", quoteRows, quoteEvents, err)
	}
	type quoteResult struct {
		value    franchisejourney.Quote
		replayed bool
		err      error
	}
	startQuotes := make(chan struct{})
	quoteResults := make(chan quoteResult, 2)
	concurrentValidUntil := time.Now().UTC().Add(48 * time.Hour)
	for _, input := range []struct{ id, event string }{
		{"quote-concurrent-a", "018f4d4a-7b36-7a21-8d10-2f4c54c29b40"},
		{"quote-concurrent-b", "018f4d4a-7b36-7a21-8d10-2f4c54c29b41"},
	} {
		go func(input struct{ id, event string }) {
			<-startQuotes
			value, replayed, quoteErr := repo.CreateQuote(ctx, tenant, "quote-concurrent-key", franchisejourney.Quote{ID: input.id, OrganizationID: "store", LeadID: "lead", VariantID: "variant", PriceBookID: "retail", ValidUntil: concurrentValidUntil, State: "issued", Version: 1}, strings.Repeat("9", 64), input.event)
			quoteResults <- quoteResult{value: value, replayed: replayed, err: quoteErr}
		}(input)
	}
	close(startQuotes)
	createdQuotes, replayedQuotes := 0, 0
	var concurrentQuoteID string
	for range 2 {
		result := <-quoteResults
		if result.err != nil {
			t.Fatalf("concurrent quote failed: %v", result.err)
		}
		if concurrentQuoteID == "" {
			concurrentQuoteID = result.value.ID
		} else if result.value.ID != concurrentQuoteID {
			t.Fatalf("concurrent replay returned different ids: %q vs %q", concurrentQuoteID, result.value.ID)
		}
		if result.replayed {
			replayedQuotes++
		} else {
			createdQuotes++
		}
	}
	if createdQuotes != 1 || replayedQuotes != 1 {
		t.Fatalf("concurrent quote created=%d replayed=%d", createdQuotes, replayedQuotes)
	}
	if err = pool.QueryRow(ctx, `select count(*) from sales.quotation where tenant_id=$1 and quotation_id=$2`, tenant, concurrentQuoteID).Scan(&quoteRows); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id=$2 and event_type='quotation.issued'`, tenant, concurrentQuoteID).Scan(&quoteEvents); err != nil || quoteRows != 1 || quoteEvents != 1 {
		t.Fatalf("concurrent quote rows=%d events=%d err=%v", quoteRows, quoteEvents, err)
	}
	accepted, err := repo.AcceptQuote(ctx, tenant, "store", "customer", "quote", 1, "dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd", "accepted-order", "accepted-line", "018f4d4a-7b36-7a21-8d10-2f4c54c29b09", "018f4d4a-7b36-7a21-8d10-2f4c54c29b0a")
	if err != nil || accepted.State != "accepted" || accepted.OrderID != "accepted-order" || accepted.Version != 2 {
		t.Fatalf("accepted=%+v err=%v", accepted, err)
	}
	var orderState, currency string
	var total int64
	if err = pool.QueryRow(ctx, `select state,currency,total_minor_units from sales.customer_order where tenant_id=$1 and order_id='accepted-order'`, tenant).Scan(&orderState, &currency, &total); err != nil || orderState != "placed" || currency != "ARS" || total != 123456 {
		t.Fatalf("order state=%s currency=%s total=%d err=%v", orderState, currency, total, err)
	}
	var linePrice int64
	if err = pool.QueryRow(ctx, `select unit_price_minor_units from sales.customer_order_line where tenant_id=$1 and order_id='accepted-order'`, tenant).Scan(&linePrice); err != nil || linePrice != 123456 {
		t.Fatalf("line price=%d err=%v", linePrice, err)
	}
	if _, err = repo.AcceptQuote(ctx, tenant, "store", "customer", "quote", 1, strings.Repeat("d", 64), "stale-order", "stale-line", "018f4d4a-7b36-7a21-8d10-2f4c54c29b0b", "018f4d4a-7b36-7a21-8d10-2f4c54c29b0c"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("stale quote acceptance succeeded")
	}
	if _, err = repo.AcceptQuote(ctx, tenant, "store", "other-customer", "quote", 2, strings.Repeat("d", 64), "cross-order", "cross-line", "018f4d4a-7b36-7a21-8d10-2f4c54c29b0d", "018f4d4a-7b36-7a21-8d10-2f4c54c29b0e"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("cross-customer quote acceptance succeeded")
	}
	if _, err = pool.Exec(ctx, `insert into sales.quotation(tenant_id,quotation_id,organization_id,lead_id,customer_principal_id,variant_id,price_book_id,currency,total_minor_units,valid_until,state,version,created_at) values($1,'expired-quote','store','lead','customer','variant','retail','ARS',123456,clock_timestamp()-interval '1 day','issued',1,clock_timestamp()-interval '2 days')`, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err = repo.AcceptQuote(ctx, tenant, "store", "customer", "expired-quote", 1, strings.Repeat("e", 64), "expired-order", "expired-line", "018f4d4a-7b36-7a21-8d10-2f4c54c29b0f", "018f4d4a-7b36-7a21-8d10-2f4c54c29b10"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("expired quote acceptance succeeded")
	}
	failedQuote, _, err := repo.CreateQuote(ctx, tenant, "quote-request-0002", franchisejourney.Quote{ID: "rollback-quote", OrganizationID: "store", LeadID: "lead", VariantID: "variant", PriceBookID: "retail", ValidUntil: time.Now().Add(24 * time.Hour), State: "issued", Version: 1}, strings.Repeat("7", 64), "018f4d4a-7b36-7a21-8d10-2f4c54c29b11")
	if err != nil || failedQuote.ID == "" {
		t.Fatal(err)
	}
	_, err = repo.AcceptQuote(ctx, tenant, "store", "customer", "rollback-quote", 1, strings.Repeat("f", 64), "rollback-order", "rollback-line", "018f4d4a-7b36-7a21-8d10-2f4c54c29b12", "018f4d4a-7b36-7a21-8d10-2f4c54c29b11")
	if err == nil {
		t.Fatal("duplicate outbox event did not abort acceptance")
	}
	var rollbackOrders, rollbackAcceptances, rollbackAcceptedEvents int
	if err = pool.QueryRow(ctx, `select count(*) from sales.customer_order where tenant_id=$1 and order_id='rollback-order'`, tenant).Scan(&rollbackOrders); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from sales.quotation_acceptance where tenant_id=$1 and quotation_id='rollback-quote'`, tenant).Scan(&rollbackAcceptances); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id='rollback-quote' and event_type='quotation.accepted'`, tenant).Scan(&rollbackAcceptedEvents); err != nil {
		t.Fatal(err)
	}
	var rollbackState string
	if err = pool.QueryRow(ctx, `select state from sales.quotation where tenant_id=$1 and quotation_id='rollback-quote'`, tenant).Scan(&rollbackState); err != nil || rollbackOrders != 0 || rollbackAcceptances != 0 || rollbackAcceptedEvents != 0 || rollbackState != "issued" {
		t.Fatalf("rollback state=%s orders=%d acceptances=%d events=%d err=%v", rollbackState, rollbackOrders, rollbackAcceptances, rollbackAcceptedEvents, err)
	}
	if _, _, err = repo.CreateQuote(ctx, tenant, "quote-request-0003", franchisejourney.Quote{ID: "orphan-quote", OrganizationID: "store", LeadID: "missing", VariantID: "variant", PriceBookID: "retail", ValidUntil: time.Now().Add(24 * time.Hour), State: "issued", Version: 1}, strings.Repeat("8", 64), "018f4d4a-7b36-7a21-8d10-2f4c54c29b08"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatalf("quote without scoped lead was accepted: %v", err)
	}
	var orphanEvents int
	if err = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id='orphan-quote'`, tenant).Scan(&orphanEvents); err != nil || orphanEvents != 0 {
		t.Fatalf("orphan quote outbox events=%d err=%v", orphanEvents, err)
	}
	evidence := "cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"
	checklist, err := repo.PublishDeliveryChecklist(ctx, tenant, "operator", franchisejourney.DeliveryChecklist{ID: "standard-delivery", OrganizationID: "store", Version: 1, Title: "Entrega estándar", State: "published", Items: []franchisejourney.ChecklistItem{{ID: "serial-observed", Ordinal: 1, Prompt: "Verificar serie", ResponseType: "serial", Required: true}, {ID: "asset-condition", Ordinal: 2, Prompt: "Confirmar condición", ResponseType: "confirmation", Required: true}}}, "018f4d4a-7b36-7a21-8d10-2f4c54c29b30")
	if err != nil || checklist.State != "published" {
		t.Fatalf("checklist=%+v err=%v", checklist, err)
	}
	// The existing foreign keys are tenant-scoped, not relational authorization.
	// Corrupt only this synthetic tenant to exercise a stale or inconsistent link.
	for phaseNo, phase := range []string{"complete", "accept", "reject"} {
		for caseNo, mismatch := range []struct{ name, corrupt, restore string }{
			{"order-organization", `update sales.customer_order set organization_id='other' where tenant_id=$1 and order_id='order'`, `update sales.customer_order set organization_id='store' where tenant_id=$1 and order_id='order'`},
			{"order-customer", `update sales.customer_order set customer_principal_id='other-customer' where tenant_id=$1 and order_id='order'`, `update sales.customer_order set customer_principal_id='customer' where tenant_id=$1 and order_id='order'`},
			{"stock-organization", `update inventory.stock_unit set organization_id='other' where tenant_id=$1 and stock_unit_id='stock'`, `update inventory.stock_unit set organization_id='store' where tenant_id=$1 and stock_unit_id='stock'`},
		} {
			t.Run(phase+"-"+mismatch.name, func(t *testing.T) {
				id := "scope-" + phase + "-" + mismatch.name
				if _, err := pool.Exec(ctx, `insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version) values($1,$2,'store','order','customer','stock','prepared',1)`, tenant, id); err != nil {
					t.Fatal(err)
				}
				responses := []franchisejourney.ChecklistResponse{{ItemID: "serial-observed", ResponseText: "SERIAL-JOURNEY"}, {ItemID: "asset-condition", ResponseText: "confirmed"}}
				event := fmt.Sprintf("018f4d4a-7b36-7a21-8d10-2f4c54c2%04d", 9100+phaseNo*10+caseNo)
				if phase != "complete" {
					if _, err := repo.CompleteDeliveryChecklist(ctx, tenant, "store", "operator", id, 1, checklist.ID, 1, responses, event); err != nil {
						t.Fatal(err)
					}
				}
				if _, err := pool.Exec(ctx, mismatch.corrupt, tenant); err != nil {
					t.Fatal(err)
				}
				defer func() {
					if _, err := pool.Exec(ctx, mismatch.restore, tenant); err != nil {
						t.Fatal(err)
					}
				}()
				event = fmt.Sprintf("018f4d4a-7b36-7a21-8d10-2f4c54c2%04d", 9200+phaseNo*10+caseNo)
				if phase == "complete" {
					read, readErr := repo.CustomerJourney(ctx, tenant, "store", "customer")
					if !errors.Is(readErr, franchisejourney.ErrConflict) || len(read.Appointments)+len(read.Quotes)+len(read.Handovers)+len(read.Exceptions) != 0 {
						t.Errorf("inconsistent customer projection returned data: err=%v handovers=%d", readErr, len(read.Handovers))
					}
					unrelated, unrelatedErr := repo.CustomerJourney(ctx, tenant, "store", "unrelated-customer")
					if unrelatedErr != nil || len(unrelated.Handovers) != 0 {
						t.Errorf("unrelated customer affected: %v", unrelatedErr)
					}
					ops, opsErr := NewCommerce(pool).Operations(ctx, tenant, "store")
					if mismatch.name != "order-organization" {
						if opsErr == nil || len(ops.Orders)+len(ops.Stock) != 0 {
							t.Errorf("inconsistent operation projection returned data: err=%v orders=%d", opsErr, len(ops.Orders))
						}
					} else {
						for _, order := range ops.Orders {
							if order.ID == "order" {
								t.Error("foreign organization order leaked")
							}
						}
					}
				}
				var actionErr error
				switch phase {
				case "complete":
					_, actionErr = repo.CompleteDeliveryChecklist(ctx, tenant, "store", "operator", id, 1, checklist.ID, 1, responses, event)
				case "accept":
					_, actionErr = repo.AcceptHandover(ctx, tenant, "store", "customer", id, 2, "SERIAL-JOURNEY", checklist.ID, 1, evidence, event)
				case "reject":
					_, actionErr = repo.RejectHandover(ctx, tenant, "store", "customer", id, 2, "visible-damage", "synthetic", evidence, id+"-exception", event)
				}
				if !errors.Is(actionErr, franchisejourney.ErrConflict) {
					t.Errorf("inconsistent %s allowed %s: %v", mismatch.name, phase, actionErr)
				}
				var state string
				var events, responseCount int
				if err := pool.QueryRow(ctx, `select state,(select count(*) from platform.outbox_event where tenant_id=$1 and event_id=$3),(select count(*) from sales.delivery_checklist_response where tenant_id=$1 and handover_id=$2) from sales.delivery_handover where tenant_id=$1 and handover_id=$2`, tenant, id, event).Scan(&state, &events, &responseCount); err != nil {
					t.Fatal(err)
				}
				wantState, wantResponses := "presented", 2
				if phase == "complete" {
					wantState, wantResponses = "prepared", 0
				}
				if state != wantState || events != 0 || responseCount != wantResponses {
					t.Errorf("partial effect state=%s events=%d responses=%d", state, events, responseCount)
				}
			})
		}
	}
	if _, err = repo.AcceptHandover(ctx, tenant, "store", "customer", "handover", 1, "SERIAL-JOURNEY", checklist.ID, checklist.Version, evidence, "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatalf("handover accepted before checklist completion: %v", err)
	}
	if _, err = repo.CompleteDeliveryChecklist(ctx, tenant, "store", "operator", "handover", 1, checklist.ID, checklist.Version, []franchisejourney.ChecklistResponse{{ItemID: "serial-observed", ResponseText: "SERIAL-JOURNEY"}}, "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatalf("incomplete required checklist accepted: %v", err)
	}
	completed, err := repo.CompleteDeliveryChecklist(ctx, tenant, "store", "operator", "handover", 1, checklist.ID, checklist.Version, []franchisejourney.ChecklistResponse{{ItemID: "serial-observed", ResponseText: "SERIAL-JOURNEY"}, {ItemID: "asset-condition", ResponseText: "confirmed"}}, "018f4d4a-7b36-7a21-8d10-2f4c54c29b31")
	if err != nil || completed.State != "presented" || completed.Version != 2 {
		t.Fatalf("completed=%+v err=%v", completed, err)
	}
	if _, err = repo.AcceptHandover(ctx, tenant, "store", "customer", "handover", 2, "WRONG-SERIAL", checklist.ID, checklist.Version, evidence, "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatalf("handover accepted with wrong serial: %v", err)
	}
	if _, err = repo.AcceptHandover(ctx, tenant, "store", "customer", "handover", 2, "SERIAL-JOURNEY", checklist.ID, 2, evidence, "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatalf("handover accepted against wrong checklist version: %v", err)
	}
	var acceptanceWG sync.WaitGroup
	// Audit failure must roll back acceptance, leaving the same version usable.
	if _, err := repo.AcceptHandover(ctx, tenant, "store", "customer", "handover", 2, "SERIAL-JOURNEY", checklist.ID, checklist.Version, evidence, "018f4d4a-7b36-7a21-8d10-2f4c54c29b31"); err == nil {
		t.Fatal("duplicate outbox event allowed acceptance")
	}
	var recoverState string
	var recoverVersion int64
	if err := pool.QueryRow(ctx, `select state,version from sales.delivery_handover where tenant_id=$1 and handover_id='handover'`, tenant).Scan(&recoverState, &recoverVersion); err != nil || recoverState != "presented" || recoverVersion != 2 {
		t.Fatalf("acceptance rollback state=%s version=%d err=%v", recoverState, recoverVersion, err)
	}
	// Check the actual row locks, not timing of an assumed concurrent writer.
	for _, mutation := range []string{`update sales.customer_order set customer_principal_id='other-customer' where tenant_id=$1 and order_id='order'`, `update inventory.stock_unit set organization_id='other' where tenant_id=$1 and stock_unit_id='stock'`} {
		locked, err := pool.Begin(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if err = lockDeliveryScope(ctx, locked, tenant, "store", "handover"); err != nil {
			_ = locked.Rollback(ctx)
			t.Fatal(err)
		}
		writer, err := pool.Begin(ctx)
		if err != nil {
			_ = locked.Rollback(ctx)
			t.Fatal(err)
		}
		_, err = writer.Exec(ctx, `set local lock_timeout='50ms'`)
		if err != nil {
			_ = writer.Rollback(ctx)
			_ = locked.Rollback(ctx)
			t.Fatal(err)
		}
		_, err = writer.Exec(ctx, mutation, tenant)
		_ = writer.Rollback(ctx)
		_ = locked.Rollback(ctx)
		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) || pgErr.Code != "55P03" {
			t.Fatalf("linked row not fenced by transaction: %v", err)
		}
	}
	acceptanceResults := make(chan error, 16)
	for range 16 {
		acceptanceWG.Add(1)
		go func() {
			defer acceptanceWG.Done()
			value, err := repo.AcceptHandover(ctx, tenant, "store", "customer", "handover", 2, "SERIAL-JOURNEY", checklist.ID, checklist.Version, evidence, "018f4d4a-7b36-7a21-8d10-2f4c54c29b06")
			if err == nil && (value.State != "accepted" || value.Version != 3) {
				err = fmt.Errorf("invalid receipt state/version")
			}
			acceptanceResults <- err
		}()
	}
	acceptanceWG.Wait()
	close(acceptanceResults)
	acceptedCount, conflictCount := 0, 0
	for err := range acceptanceResults {
		if err == nil {
			acceptedCount++
		} else if errors.Is(err, franchisejourney.ErrConflict) {
			conflictCount++
		} else {
			t.Fatal(err)
		}
	}
	if acceptedCount != 1 || conflictCount != 15 {
		t.Fatalf("concurrent acceptance successes=%d conflicts=%d", acceptedCount, conflictCount)
	}
	var acceptedEvents int
	if err := pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id='handover' and event_type='delivery-handover.accepted'`, tenant).Scan(&acceptedEvents); err != nil || acceptedEvents != 1 {
		t.Fatalf("acceptance events=%d err=%v", acceptedEvents, err)
	}
	if _, err = repo.AcceptHandover(ctx, tenant, "store", "other-customer", "handover", 3, "SERIAL-JOURNEY", checklist.ID, checklist.Version, evidence, "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("handover crossed customer scope")
	}
	rejectedReady, err := repo.CompleteDeliveryChecklist(ctx, tenant, "store", "operator", "handover-reject", 1, checklist.ID, checklist.Version, []franchisejourney.ChecklistResponse{{ItemID: "serial-observed", ResponseText: "SERIAL-JOURNEY"}, {ItemID: "asset-condition", ResponseText: "confirmed"}}, "018f4d4a-7b36-7a21-8d10-2f4c54c29b32")
	if err != nil || rejectedReady.State != "presented" {
		t.Fatalf("rejected-ready=%+v err=%v", rejectedReady, err)
	}
	if _, err = repo.RejectHandover(ctx, tenant, "store", "other-customer", "handover-reject", 2, "visible-damage", "Rayón visible", evidence, "unused", "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("handover rejection crossed customer scope")
	}
	exception, err := repo.RejectHandover(ctx, tenant, "store", "customer", "handover-reject", 2, "visible-damage", "Rayón visible", evidence, "delivery-exception", "018f4d4a-7b36-7a21-8d10-2f4c54c29b33")
	if err != nil || exception.State != "open" || exception.Version != 1 {
		t.Fatalf("exception=%+v err=%v", exception, err)
	}
	openExceptions, err := repo.DeliveryExceptions(ctx, tenant, "store", 100)
	if err != nil || len(openExceptions) != 1 || openExceptions[0].ID != exception.ID {
		t.Fatalf("open exceptions=%+v err=%v", openExceptions, err)
	}
	for caseNo, mismatch := range []struct{ name, corrupt, restore string }{
		{"order-organization", `update sales.customer_order set organization_id='other' where tenant_id=$1 and order_id='order'`, `update sales.customer_order set organization_id='store' where tenant_id=$1 and order_id='order'`},
		{"order-customer", `update sales.customer_order set customer_principal_id='other-customer' where tenant_id=$1 and order_id='order'`, `update sales.customer_order set customer_principal_id='customer' where tenant_id=$1 and order_id='order'`},
		{"stock-organization", `update inventory.stock_unit set organization_id='other' where tenant_id=$1 and stock_unit_id='stock'`, `update inventory.stock_unit set organization_id='store' where tenant_id=$1 and stock_unit_id='stock'`},
	} {
		t.Run("resolve-"+mismatch.name, func(t *testing.T) {
			if _, err := pool.Exec(ctx, mismatch.corrupt, tenant); err != nil {
				t.Fatal(err)
			}
			defer func() {
				if _, err := pool.Exec(ctx, mismatch.restore, tenant); err != nil {
					t.Fatal(err)
				}
			}()
			_, err := repo.ResolveDeliveryException(ctx, tenant, "store", "operator", exception.ID, 1, "correct-and-represent", "synthetic", "bad-successor", "unused", "bad-resolution", fmt.Sprintf("018f4d4a-7b36-7a21-8d10-2f4c54c2%04d", 9300+caseNo))
			if !errors.Is(err, franchisejourney.ErrConflict) {
				t.Errorf("inconsistent resolution allowed: %v", err)
			}
			var state string
			var effects int
			if err := pool.QueryRow(ctx, `select state,(select count(*) from sales.delivery_exception_resolution where tenant_id=$1 and exception_id=$2) from sales.delivery_exception where tenant_id=$1 and exception_id=$2`, tenant, exception.ID).Scan(&state, &effects); err != nil || state != "open" || effects != 0 {
				t.Errorf("resolution state=%s effects=%d err=%v", state, effects, err)
			}
		})
	}
	resolution, err := repo.ResolveDeliveryException(ctx, tenant, "store", "operator", exception.ID, 1, "correct-and-represent", "Corregir preparación", "handover-successor", "unused-authorization", "delivery-resolution", "018f4d4a-7b36-7a21-8d10-2f4c54c29b34")
	if err != nil || resolution.Exception.State != "resolved" || resolution.SuccessorHandover == nil || resolution.SuccessorHandover.State != "prepared" || resolution.SuccessorHandover.SupersedesHandoverID != "handover-reject" {
		t.Fatalf("resolution=%+v err=%v", resolution, err)
	}
	if _, err = repo.ResolveDeliveryException(ctx, tenant, "store", "operator", exception.ID, 1, "return", "retry", "unused", "unused", "unused", "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("resolved delivery exception was processed twice")
	}
	returnReady, err := repo.CompleteDeliveryChecklist(ctx, tenant, "store", "operator", "handover-return", 1, checklist.ID, checklist.Version, []franchisejourney.ChecklistResponse{{ItemID: "serial-observed", ResponseText: "SERIAL-JOURNEY"}, {ItemID: "asset-condition", ResponseText: "confirmed"}}, "018f4d4a-7b36-7a21-8d10-2f4c54c29b35")
	if err != nil || returnReady.State != "presented" {
		t.Fatalf("return-ready=%+v err=%v", returnReady, err)
	}
	returnException, err := repo.RejectHandover(ctx, tenant, "store", "customer", "handover-return", 2, "customer-return", "Cliente solicita devolución", evidence, "delivery-return-exception", "018f4d4a-7b36-7a21-8d10-2f4c54c29b36")
	if err != nil || returnException.State != "open" {
		t.Fatalf("return exception=%+v err=%v", returnException, err)
	}
	returnResolution, err := repo.ResolveDeliveryException(ctx, tenant, "store", "operator", returnException.ID, 1, "return", "Autorizar recepción y revisión", "unused-successor", "return-authorization", "delivery-return-resolution", "018f4d4a-7b36-7a21-8d10-2f4c54c29b37")
	if err != nil || returnResolution.Exception.State != "resolved" || returnResolution.ReturnAuthorizationID != "return-authorization" || returnResolution.Disposition != "return" || returnResolution.SuccessorHandover != nil {
		t.Fatalf("return resolution=%+v err=%v", returnResolution, err)
	}
	var authorizationOrder, authorizationStock, exactCostOrder, authorizationState string
	if err = pool.QueryRow(ctx, `select order_id,stock_unit_id,exact_cost_source_order_id,state from sales.return_authorization where tenant_id=$1 and authorization_id='return-authorization'`, tenant).Scan(&authorizationOrder, &authorizationStock, &exactCostOrder, &authorizationState); err != nil {
		t.Fatal(err)
	}
	if authorizationOrder != "order" || authorizationStock != "stock" || exactCostOrder != "order" || authorizationState != "authorized" {
		t.Fatalf("return authorization order=%q stock=%q exact_cost_order=%q state=%q", authorizationOrder, authorizationStock, exactCostOrder, authorizationState)
	}
	if _, err = pool.Exec(ctx, `update sales.return_authorization set state='received' where tenant_id=$1 and authorization_id='return-authorization'`, tenant); err == nil {
		t.Fatal("return authorization accepted an in-place mutation")
	}
	if _, err = repo.ReceiveReturn(ctx, tenant, "store", "operator", "return-authorization", "WRONG-SERIAL", "damaged", "Daño confirmado", evidence, "return-receipt-invalid", "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatalf("return receipt accepted wrong durable serial: %v", err)
	}
	receipt, err := repo.ReceiveReturn(ctx, tenant, "store", "operator", "return-authorization", "SERIAL-JOURNEY", "damaged", "Daño confirmado", evidence, "return-receipt", "018f4d4a-7b36-7a21-8d10-2f4c54c29b38")
	if err != nil || receipt.AuthorizationID != "return-authorization" || receipt.OrderID != "order" || receipt.StockUnitID != "stock" || receipt.ConditionCode != "damaged" {
		t.Fatalf("return receipt=%+v err=%v", receipt, err)
	}
	if _, err = repo.ReceiveReturn(ctx, tenant, "store", "operator", "return-authorization", "SERIAL-JOURNEY", "damaged", "duplicate", evidence, "return-receipt-duplicate", "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("return authorization was received twice")
	}
	disposition, err := repo.DecideReturn(ctx, tenant, "store", "operator", receipt.ID, "quarantine", "Separar y solicitar efectos", "return-disposition", "return-effect-inventory", "return-effect-remedy", "return-effect-accounting", "return-effect-fiscal", "018f4d4a-7b36-7a21-8d10-2f4c54c29b39")
	if err != nil || disposition.CustomerRemedy != "refund" || disposition.InventoryAction != "quarantine" || len(disposition.Effects) != 4 {
		t.Fatalf("return disposition=%+v err=%v", disposition, err)
	}
	owners := map[string]string{}
	for _, effect := range disposition.Effects {
		owners[effect.EffectKind] = effect.OwnerContext
	}
	if owners["inventory"] != "inventory" || owners["refund"] != "payment" || owners["accounting"] != "accounting" || owners["fiscal"] != "fiscal" {
		t.Fatalf("return effect owners=%+v", owners)
	}
	if _, err = repo.DecideReturn(ctx, tenant, "store", "operator", receipt.ID, "restock", "duplicate", "unused", "unused", "unused", "unused", "unused", "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("return receipt received two dispositions")
	}
	returnCases, err := repo.ReturnCases(ctx, tenant, "store", 100)
	if err != nil || len(returnCases) != 1 || returnCases[0].Receipt == nil || returnCases[0].Disposition == nil || len(returnCases[0].Disposition.Effects) != 4 {
		t.Fatalf("return cases=%+v err=%v", returnCases, err)
	}
	if _, err = pool.Exec(ctx, `update sales.return_receipt set condition_code='sealed' where tenant_id=$1 and receipt_id='return-receipt'`, tenant); err == nil {
		t.Fatal("return receipt accepted in-place mutation")
	}
	journey, err := repo.CustomerJourney(ctx, tenant, "store", "customer")
	acceptedQuotes := 0
	for _, item := range journey.Quotes {
		if item.ID == "quote" && item.State == "accepted" && item.OrderID == "accepted-order" {
			acceptedQuotes++
		}
	}
	successors := 0
	for _, item := range journey.Handovers {
		if item.SupersedesHandoverID == "handover-reject" && item.State == "prepared" {
			successors++
		}
	}
	negativeHandovers := 0
	for _, h := range journey.Handovers {
		if strings.HasPrefix(h.ID, "scope-") {
			negativeHandovers++
		}
	}
	if err != nil || len(journey.Appointments) != 4 || len(journey.Quotes) != 4 || acceptedQuotes != 1 || len(journey.Handovers) != 13 || negativeHandovers != 9 || successors != 1 || len(journey.Exceptions) != 2 || journey.Exceptions[0].State != "resolved" || journey.Exceptions[1].State != "resolved" {
		t.Fatalf("journey=%+v err=%v", journey, err)
	}
	fresh, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer fresh.Close()
	var durableState string
	var durableVersion int64
	var durableEvents int
	if err := fresh.QueryRow(ctx, `select state,version,(select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id='handover' and event_type='delivery-handover.accepted') from sales.delivery_handover where tenant_id=$1 and handover_id='handover'`, tenant).Scan(&durableState, &durableVersion, &durableEvents); err != nil || durableState != "accepted" || durableVersion != 3 || durableEvents != 1 {
		t.Fatalf("fresh pool state=%s version=%d events=%d err=%v", durableState, durableVersion, durableEvents, err)
	}
	// Keep immutable history intact; this is a separate, FK-valid bad import.
	probe := &journeyReadProbe{afterFirstRead: func() error {
		_, err := pool.Exec(ctx, `update sales.customer_order set customer_principal_id='other-customer' where tenant_id=$1 and order_id='order'`, tenant)
		return err
	}}
	probeConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		t.Fatal(err)
	}
	probeConfig.ConnConfig.Tracer = probe
	probePool, err := pgxpool.NewWithConfig(ctx, probeConfig)
	if err != nil {
		t.Fatal(err)
	}
	defer probePool.Close()
	consistent, consistentErr := NewFranchiseJourney(probePool).CustomerJourney(ctx, tenant, "store", "customer")
	changed, changedErr := repo.CustomerJourney(ctx, tenant, "store", "customer")
	if _, err := pool.Exec(ctx, `update sales.customer_order set customer_principal_id='customer' where tenant_id=$1 and order_id='order'`, tenant); err != nil {
		t.Fatal(err)
	}
	if consistentErr != nil || len(consistent.Handovers) != 13 || probe.err != nil || !probe.inTransaction || probe.boundedItems != 13 {
		t.Fatalf("snapshot/bound failed read=%v writer=%v handovers=%d tx=%v items=%d", consistentErr, probe.err, len(consistent.Handovers), probe.inTransaction, probe.boundedItems)
	}
	if !errors.Is(changedErr, franchisejourney.ErrConflict) || len(changed.Handovers) != 0 {
		t.Fatalf("new snapshot did not reject changed relation: %v", changedErr)
	}
	recovered, recoveredErr := repo.CustomerJourney(ctx, tenant, "store", "customer")
	if recoveredErr != nil || len(recovered.Handovers) != 13 {
		t.Fatalf("read recovery failed: %v", recoveredErr)
	}
	for _, q := range []string{
		`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized)values($1,'unrelated-customer','Synthetic','unrelated@example.test')`,
		`insert into sales.delivery_exception(tenant_id,exception_id,organization_id,handover_id,customer_principal_id,reason_code,details,rejection_evidence_sha256_hex,state,version) values($1,'bad-import-exception','store','handover-successor','unrelated-customer','synthetic-case','PRIVATE-SYNTHETIC-NOTE',repeat('a',64),'open',1)`,
	} {
		if _, err := pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	wrongJourney, wrongJourneyErr := repo.CustomerJourney(ctx, tenant, "store", "unrelated-customer")
	wrongExceptions, wrongExceptionsErr := repo.DeliveryExceptions(ctx, tenant, "store", 100)
	if !errors.Is(wrongJourneyErr, franchisejourney.ErrConflict) || len(wrongJourney.Exceptions) != 0 {
		t.Errorf("foreign exception exposed to customer: %v", wrongJourneyErr)
	}
	if !errors.Is(wrongExceptionsErr, franchisejourney.ErrConflict) || len(wrongExceptions) != 0 {
		t.Errorf("incoherent exception exposed to operator: %v", wrongExceptionsErr)
	}
}

func TestAvailabilityCreationIdentityAndRecovery(t *testing.T) {
	raw := os.Getenv("TEST_DATABASE_URL")
	if raw == "" {
		t.Skip("disposable database not selected")
	}
	cfg, err := pgxpool.ParseConfig(raw)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires disposable loopback database")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := (randomid.Generator{}).New()
	for _, sql := range []string{`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'create-'||replace(($1::uuid)::text,'-',''),'Synthetic','Synthetic')`, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic','store')`} {
		if _, err = pool.Exec(ctx, sql, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewFranchiseJourney(pool)
	start := time.Now().UTC().Add(72 * time.Hour).Truncate(time.Second)
	input := franchisejourney.AvailabilityEntry{ID: "created-interval", OrganizationID: "store", EntryType: "working", StartsAt: start, EndsAt: start.Add(time.Hour)}
	const key = "availability-concurrent-key"
	hash := strings.Repeat("a", 64)
	type outcome struct {
		v      franchisejourney.AvailabilityEntry
		replay bool
		err    error
	}
	results := make(chan outcome, 8)
	begin := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-begin
			v := input
			v.ID = fmt.Sprintf("created-%d", i)
			got, replay, e := repo.CreateAvailabilityOnce(ctx, tenant, "scheduler", key, hash, v, (randomid.Generator{}).New())
			results <- outcome{got, replay, e}
		}(i)
	}
	close(begin)
	wg.Wait()
	close(results)
	created, replays := 0, 0
	var id string
	for r := range results {
		if r.err != nil {
			t.Fatal(r.err)
		}
		if id != "" && id != r.v.ID {
			t.Fatal("duplicate identity")
		}
		id = r.v.ID
		if r.replay {
			replays++
		} else {
			created++
		}
	}
	if created != 1 || replays != 7 {
		t.Fatalf("created=%d replay=%d", created, replays)
	}
	recovered, e := repo.AvailabilityCreationResult(ctx, tenant, "store", key)
	if e != nil || recovered.ID != id {
		t.Fatalf("recover=%+v err=%v", recovered, e)
	}
	for _, scope := range [][3]string{{tenant, "other", key}, {tenant, "store", "unknown-availability-key"}, {(randomid.Generator{}).New(), "store", key}} {
		if v, e := repo.AvailabilityCreationResult(ctx, scope[0], scope[1], scope[2]); !errors.Is(e, franchisejourney.ErrConflict) || v.ID != "" {
			t.Fatalf("scope leak=%+v %v", v, e)
		}
	}
	if _, _, e := repo.CreateAvailabilityOnce(ctx, tenant, "scheduler", key, strings.Repeat("b", 64), input, (randomid.Generator{}).New()); !errors.Is(e, franchisejourney.ErrConflict) {
		t.Fatal("divergent intent accepted")
	}
	var event, actor string
	var rows, keys, events int
	if e := pool.QueryRow(ctx, `select event_id::text,payload->>'actor_subject' from platform.outbox_event where tenant_id=$1 and aggregate_id=$2 and event_type='availability-entry.created'`, tenant, id).Scan(&event, &actor); e != nil || actor != "scheduler" {
		t.Fatalf("actor=%s err=%v", actor, e)
	}
	if _, _, e := repo.CreateAvailabilityOnce(ctx, tenant, "another-operator", key, hash, input, (randomid.Generator{}).New()); e != nil {
		t.Fatal(e)
	}
	next := input
	next.ID = "rolled-back"
	next.StartsAt = start.Add(48 * time.Hour)
	next.EndsAt = next.StartsAt.Add(time.Hour)
	if v, _, e := repo.CreateAvailabilityOnce(ctx, tenant, "scheduler", "availability-rollback-key", hash, next, event); e == nil || v.ID != "" {
		t.Fatalf("audit failure not atomic=%+v %v", v, e)
	}
	if e := pool.QueryRow(ctx, `select (select count(*) from crm.availability_entry where tenant_id=$1),(select count(*) from platform.idempotency_record where tenant_id=$1 and scope='franchise-availability'),(select count(*) from platform.outbox_event where tenant_id=$1 and event_type='availability-entry.created')`, tenant).Scan(&rows, &keys, &events); e != nil || rows != 1 || keys != 1 || events != 1 {
		t.Fatalf("effects=%d/%d/%d err=%v", rows, keys, events, e)
	}
	if _, e := repo.CancelAvailability(ctx, tenant, "store", id, 1, "scheduler", "schedule-correction", (randomid.Generator{}).New()); e != nil {
		t.Fatal(e)
	}
	after, replay, e := repo.CreateAvailabilityOnce(ctx, tenant, "another-operator", key, hash, input, (randomid.Generator{}).New())
	if e != nil || !replay || after.ID != id || after.State != "cancelled" || after.Version != 2 {
		t.Fatalf("replay resurrected interval=%+v replay=%v err=%v", after, replay, e)
	}
	fresh, e := pgxpool.NewWithConfig(ctx, cfg.Copy())
	if e != nil {
		t.Fatal(e)
	}
	defer fresh.Close()
	after, e = NewFranchiseJourney(fresh).AvailabilityCreationResult(ctx, tenant, "store", key)
	if e != nil || after.State != "cancelled" {
		t.Fatalf("fresh pool recovery=%+v err=%v", after, e)
	}
	for _, mutation := range []string{"resource_type='other'", "response_code=200", "response_body='{\"availability_id\":\"other\"}'"} {
		if _, e := pool.Exec(ctx, "update platform.idempotency_record set "+mutation+" where tenant_id=$1 and scope='franchise-availability'", tenant); e != nil {
			t.Fatal(e)
		}
		if v, e := repo.AvailabilityCreationResult(ctx, tenant, "store", key); !errors.Is(e, franchisejourney.ErrConflict) || v.ID != "" {
			t.Fatal("incoherent receipt accepted")
		}
		if _, e := pool.Exec(ctx, `update platform.idempotency_record set resource_type='availability-entry',response_code=201,response_body=jsonb_build_object('availability_id',$2::text) where tenant_id=$1 and scope='franchise-availability'`, tenant, id); e != nil {
			t.Fatal(e)
		}
	}
	if e := pool.QueryRow(ctx, `select payload->>'actor_subject' from platform.outbox_event where tenant_id=$1 and event_id=$2::uuid`, tenant, event).Scan(&actor); e != nil || actor != "scheduler" {
		t.Fatal("replay overwrote creator")
	}
	t.Log("AVAILABILITY_IDENTITY_PASS create=1 replay=7 actor_preserved cancelled_replay rollback scoped_lookup corruption fresh_pool")
}

func TestResourceCreationIdentityAndRecovery(t *testing.T) {
	raw := os.Getenv("TEST_DATABASE_URL")
	if raw == "" {
		t.Skip("disposable database not selected")
	}
	cfg, err := pgxpool.ParseConfig(raw)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires disposable loopback database")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := (randomid.Generator{}).New()
	for _, sql := range []string{`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'create-'||replace(($1::uuid)::text,'-',''),'Synthetic','Synthetic')`, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic','store')`} {
		if _, err = pool.Exec(ctx, sql, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewFranchiseJourney(pool)
	input := franchisejourney.ServiceResource{ID: "created-resource", OrganizationID: "store", DisplayName: "Synthetic bay", Kind: "service-bay", Status: "active", Version: 1, Skills: []string{"service", "delivery"}}
	const key = "resource-concurrent-key"
	hash := strings.Repeat("a", 64)
	type outcome struct {
		v      franchisejourney.ServiceResource
		replay bool
		err    error
	}
	results := make(chan outcome, 8)
	begin := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-begin
			v := input
			v.ID = fmt.Sprintf("created-%d", i)
			got, replay, e := repo.CreateServiceResourceOnce(ctx, tenant, "resource-manager", key, hash, v, (randomid.Generator{}).New())
			results <- outcome{got, replay, e}
		}(i)
	}
	close(begin)
	wg.Wait()
	close(results)
	created, replays := 0, 0
	var id string
	for r := range results {
		if r.err != nil {
			t.Fatal(r.err)
		}
		if id != "" && id != r.v.ID {
			t.Fatal("duplicate identity")
		}
		id = r.v.ID
		if r.replay {
			replays++
		} else {
			created++
		}
	}
	if created != 1 || replays != 7 {
		t.Fatalf("created=%d replay=%d", created, replays)
	}
	recovered, e := repo.ServiceResourceCreationResult(ctx, tenant, "store", key)
	if e != nil || recovered.ID != id {
		t.Fatalf("recover=%+v err=%v", recovered, e)
	}
	for _, scope := range [][3]string{{tenant, "other", key}, {tenant, "store", "unknown-resource-key"}, {(randomid.Generator{}).New(), "store", key}} {
		if v, e := repo.ServiceResourceCreationResult(ctx, scope[0], scope[1], scope[2]); !errors.Is(e, franchisejourney.ErrConflict) || v.ID != "" {
			t.Fatalf("scope leak=%+v %v", v, e)
		}
	}
	if _, _, e := repo.CreateServiceResourceOnce(ctx, tenant, "resource-manager", key, strings.Repeat("b", 64), input, (randomid.Generator{}).New()); !errors.Is(e, franchisejourney.ErrConflict) {
		t.Fatal("divergent intent accepted")
	}
	var event, actor string
	var rows, keys, events int
	if e := pool.QueryRow(ctx, `select event_id::text,payload->>'actor_subject' from platform.outbox_event where tenant_id=$1 and aggregate_id=$2 and event_type='service-resource.created'`, tenant, id).Scan(&event, &actor); e != nil || actor != "resource-manager" {
		t.Fatalf("actor=%s err=%v", actor, e)
	}
	if _, _, e := repo.CreateServiceResourceOnce(ctx, tenant, "another-operator", key, hash, input, (randomid.Generator{}).New()); e != nil {
		t.Fatal(e)
	}
	next := input
	next.ID = "rolled-back"
	if v, _, e := repo.CreateServiceResourceOnce(ctx, tenant, "resource-manager", "resource-rollback-key", hash, next, event); e == nil || v.ID != "" {
		t.Fatalf("audit failure not atomic=%+v %v", v, e)
	}
	if e := pool.QueryRow(ctx, `select (select count(*) from crm.service_resource where tenant_id=$1),(select count(*) from platform.idempotency_record where tenant_id=$1 and scope='franchise-resource'),(select count(*) from platform.outbox_event where tenant_id=$1 and event_type='service-resource.created')`, tenant).Scan(&rows, &keys, &events); e != nil || rows != 1 || keys != 1 || events != 1 {
		t.Fatalf("effects=%d/%d/%d err=%v", rows, keys, events, e)
	}
	if _, e := pool.Exec(ctx, `update crm.service_resource set status='inactive',version=2 where tenant_id=$1 and resource_id=$2`, tenant, id); e != nil {
		t.Fatal(e)
	}
	after, replay, e := repo.CreateServiceResourceOnce(ctx, tenant, "another-operator", key, hash, input, (randomid.Generator{}).New())
	if e != nil || !replay || after.ID != id || after.Status != "inactive" || after.Version != 2 {
		t.Fatalf("replay reactivated resource=%+v replay=%v err=%v", after, replay, e)
	}
	fresh, e := pgxpool.NewWithConfig(ctx, cfg.Copy())
	if e != nil {
		t.Fatal(e)
	}
	defer fresh.Close()
	after, e = NewFranchiseJourney(fresh).ServiceResourceCreationResult(ctx, tenant, "store", key)
	if e != nil || after.Status != "inactive" {
		t.Fatalf("fresh pool recovery=%+v err=%v", after, e)
	}
	for _, mutation := range []string{"resource_type='other'", "response_code=200", "response_body='{\"resource_id\":\"other\"}'"} {
		if _, e := pool.Exec(ctx, "update platform.idempotency_record set "+mutation+" where tenant_id=$1 and scope='franchise-resource'", tenant); e != nil {
			t.Fatal(e)
		}
		if v, e := repo.ServiceResourceCreationResult(ctx, tenant, "store", key); !errors.Is(e, franchisejourney.ErrConflict) || v.ID != "" {
			t.Fatal("incoherent receipt accepted")
		}
		if _, e := pool.Exec(ctx, `update platform.idempotency_record set resource_type='service-resource',response_code=201,response_body=jsonb_build_object('resource_id',$2::text) where tenant_id=$1 and scope='franchise-resource'`, tenant, id); e != nil {
			t.Fatal(e)
		}
	}
	if e := pool.QueryRow(ctx, `select payload->>'actor_subject' from platform.outbox_event where tenant_id=$1 and event_id=$2::uuid`, tenant, event).Scan(&actor); e != nil || actor != "resource-manager" {
		t.Fatal("replay overwrote creator")
	}
	var skills int
	if err := pool.QueryRow(ctx, `select count(*) from crm.resource_skill where tenant_id=$1`, tenant).Scan(&skills); err != nil || skills != 2 {
		t.Fatalf("skill transaction leaked rows=%d error=%v", skills, err)
	}
	if len(after.Skills) != 2 || after.Skills[0] != "delivery" || after.Skills[1] != "service" {
		t.Fatal("skills not recovered coherently")
	}

	t.Log("RESOURCE_IDENTITY_PASS create=1 replay=7 actor_preserved inactive_replay rollback scoped_lookup corruption fresh_pool")
}

func TestSlotCreationIdentityAndRecovery(t *testing.T) {
	raw := os.Getenv("TEST_DATABASE_URL")
	if raw == "" {
		t.Skip("disposable database not selected")
	}
	cfg, err := pgxpool.ParseConfig(raw)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires disposable loopback database")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := (randomid.Generator{}).New()
	for _, sql := range []string{`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'create-'||replace(($1::uuid)::text,'-',''),'Synthetic','Synthetic')`, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic','store')`} {
		if _, err = pool.Exec(ctx, sql, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewFranchiseJourney(pool)
	start := time.Now().UTC().Add(72 * time.Hour).Truncate(time.Second)
	input := franchisejourney.AppointmentSlot{ID: "created-slot", OrganizationID: "store", Kind: "service", StartsAt: start, EndsAt: start.Add(time.Hour), Capacity: 2, State: "open", Version: 1}
	if _, err := repo.CreateAvailability(ctx, tenant, "slot-manager", franchisejourney.AvailabilityEntry{ID: "working-window", OrganizationID: "store", EntryType: "working", StartsAt: start.Add(-time.Hour), EndsAt: start.Add(72 * time.Hour)}, randomid.Generator{}.New()); err != nil {
		t.Fatal(err)
	}

	const key = "slot-concurrent-key"
	hash := strings.Repeat("a", 64)
	type outcome struct {
		v      franchisejourney.AppointmentSlot
		replay bool
		err    error
	}
	results := make(chan outcome, 8)
	begin := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-begin
			v := input
			v.ID = fmt.Sprintf("created-%d", i)
			got, replay, e := repo.CreateAppointmentSlotOnce(ctx, tenant, "slot-manager", key, hash, v, (randomid.Generator{}).New())
			results <- outcome{got, replay, e}
		}(i)
	}
	close(begin)
	wg.Wait()
	close(results)
	created, replays := 0, 0
	var id string
	for r := range results {
		if r.err != nil {
			t.Fatal(r.err)
		}
		if id != "" && id != r.v.ID {
			t.Fatal("duplicate identity")
		}
		id = r.v.ID
		if r.replay {
			replays++
		} else {
			created++
		}
	}
	if created != 1 || replays != 7 {
		t.Fatalf("created=%d replay=%d", created, replays)
	}
	recovered, e := repo.AppointmentSlotCreationResult(ctx, tenant, "store", key)
	if e != nil || recovered.ID != id {
		t.Fatalf("recover=%+v err=%v", recovered, e)
	}
	for _, scope := range [][3]string{{tenant, "other", key}, {tenant, "store", "unknown-slot-key"}, {(randomid.Generator{}).New(), "store", key}} {
		if v, e := repo.AppointmentSlotCreationResult(ctx, scope[0], scope[1], scope[2]); !errors.Is(e, franchisejourney.ErrConflict) || v.ID != "" {
			t.Fatalf("scope leak=%+v %v", v, e)
		}
	}
	if _, _, e := repo.CreateAppointmentSlotOnce(ctx, tenant, "slot-manager", key, strings.Repeat("b", 64), input, (randomid.Generator{}).New()); !errors.Is(e, franchisejourney.ErrConflict) {
		t.Fatal("divergent intent accepted")
	}
	var event, actor string
	var rows, keys, events int
	if e := pool.QueryRow(ctx, `select event_id::text,payload->>'actor_subject' from platform.outbox_event where tenant_id=$1 and aggregate_id=$2 and event_type='appointment-slot.created'`, tenant, id).Scan(&event, &actor); e != nil || actor != "slot-manager" {
		t.Fatalf("actor=%s err=%v", actor, e)
	}
	if _, _, e := repo.CreateAppointmentSlotOnce(ctx, tenant, "another-operator", key, hash, input, (randomid.Generator{}).New()); e != nil {
		t.Fatal(e)
	}
	next := input
	next.ID = "rolled-back"
	next.StartsAt = start.Add(48 * time.Hour)
	next.EndsAt = next.StartsAt.Add(time.Hour)
	if v, _, e := repo.CreateAppointmentSlotOnce(ctx, tenant, "slot-manager", "slot-rollback-key", hash, next, event); e == nil || v.ID != "" {
		t.Fatalf("audit failure not atomic=%+v %v", v, e)
	}
	if e := pool.QueryRow(ctx, `select (select count(*) from crm.appointment_slot where tenant_id=$1),(select count(*) from platform.idempotency_record where tenant_id=$1 and scope='franchise-slot'),(select count(*) from platform.outbox_event where tenant_id=$1 and event_type='appointment-slot.created')`, tenant).Scan(&rows, &keys, &events); e != nil || rows != 1 || keys != 1 || events != 1 {
		t.Fatalf("effects=%d/%d/%d err=%v", rows, keys, events, e)
	}
	for _, q := range []string{
		`insert into org.public_location(tenant_id,organization_id,city,region,country,published)values($1,'store','Synthetic','Synthetic','AR',true)`,
		`insert into crm.lead(tenant_id,lead_id,organization_id,lifecycle_state,source_code,contact_payload)values($1,'lead','store','new','fixture','{}')`,
	} {
		if _, e := pool.Exec(ctx, q, tenant); e != nil {
			t.Fatal(e)
		}
	}
	tenantCode := "create-" + strings.ReplaceAll(tenant, "-", "")
	for i := 0; i < 2; i++ {
		appointment := franchisejourney.Appointment{ID: fmt.Sprintf("booking-%d", i), LeadID: "lead", Kind: "service", StartsAt: start}
		if _, _, e := repo.RequestAppointment(ctx, tenantCode, "store", fmt.Sprintf("slot-booking-fixture-%d", i), appointment, hash, randomid.Generator{}.New()); e != nil {
			t.Fatal(e)
		}
	}
	full, e := repo.AppointmentSlotCreationResult(ctx, tenant, "store", key)
	if e != nil || full.Booked != 2 || full.Capacity != 2 {
		t.Fatalf("full slot not recovered: %+v %v", full, e)
	}
	public, e := repo.PublicAppointmentSlots(ctx, tenantCode, "store", "service", start.Add(-time.Minute), start.Add(2*time.Hour))
	if e != nil || len(public) != 0 {
		t.Fatal("full slot publicly bookable")
	}
	if _, e := pool.Exec(ctx, `update crm.appointment_slot set state='closed',version=2 where tenant_id=$1 and slot_id=$2`, tenant, id); e != nil {
		t.Fatal(e)
	}

	after, replay, e := repo.CreateAppointmentSlotOnce(ctx, tenant, "another-operator", key, hash, input, (randomid.Generator{}).New())
	if e != nil || !replay || after.ID != id || after.State != "closed" || after.Version != 2 {
		t.Fatalf("replay reopened slot=%+v replay=%v err=%v", after, replay, e)
	}
	fresh, e := pgxpool.NewWithConfig(ctx, cfg.Copy())
	if e != nil {
		t.Fatal(e)
	}
	defer fresh.Close()
	after, e = NewFranchiseJourney(fresh).AppointmentSlotCreationResult(ctx, tenant, "store", key)
	if e != nil || after.State != "closed" {
		t.Fatalf("fresh pool recovery=%+v err=%v", after, e)
	}
	for _, mutation := range []string{"resource_type='other'", "response_code=200", "response_body='{\"slot_id\":\"other\"}'"} {
		if _, e := pool.Exec(ctx, "update platform.idempotency_record set "+mutation+" where tenant_id=$1 and scope='franchise-slot'", tenant); e != nil {
			t.Fatal(e)
		}
		if v, e := repo.AppointmentSlotCreationResult(ctx, tenant, "store", key); !errors.Is(e, franchisejourney.ErrConflict) || v.ID != "" {
			t.Fatal("incoherent receipt accepted")
		}
		if _, e := pool.Exec(ctx, `update platform.idempotency_record set resource_type='appointment-slot',response_code=201,response_body=jsonb_build_object('slot_id',$2::text) where tenant_id=$1 and scope='franchise-slot'`, tenant, id); e != nil {
			t.Fatal(e)
		}
	}
	if e := pool.QueryRow(ctx, `select payload->>'actor_subject' from platform.outbox_event where tenant_id=$1 and event_id=$2::uuid`, tenant, event).Scan(&actor); e != nil || actor != "slot-manager" {
		t.Fatal("replay overwrote creator")
	}
	t.Log("SLOT_IDENTITY_PASS create=1 replay=7 actor_preserved full_closed_replay rollback scoped_lookup corruption fresh_pool")
}

func TestChecklistPublicationIdentityAndRecovery(t *testing.T) {
	raw := os.Getenv("TEST_DATABASE_URL")
	if raw == "" {
		t.Skip("disposable database not selected")
	}
	cfg, err := pgxpool.ParseConfig(raw)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires disposable loopback database")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := (randomid.Generator{}).New()
	for _, sql := range []string{`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'publish-'||replace(($1::uuid)::text,'-',''),'Synthetic','Synthetic')`, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic','store')`} {
		if _, err = pool.Exec(ctx, sql, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewFranchiseJourney(pool)
	input := franchisejourney.DeliveryChecklist{ID: "published-checklist", OrganizationID: "store", Version: 1, Title: "Synthetic checklist", Items: []franchisejourney.ChecklistItem{{ID: "serial", Ordinal: 1, Prompt: "Serie", ResponseType: "serial", Required: true}, {ID: "accept", Ordinal: 2, Prompt: "Confirmación", ResponseType: "confirmation", Required: true}}}
	type outcome struct {
		actor string
		err   error
	}
	results := make(chan outcome, 8)
	begin := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-begin
			actor := fmt.Sprintf("publisher-%d", i)
			_, e := repo.PublishDeliveryChecklist(ctx, tenant, actor, input, randomid.Generator{}.New())
			results <- outcome{actor, e}
		}(i)
	}
	close(begin)
	wg.Wait()
	close(results)
	created, conflicts := 0, 0
	winner := ""
	for result := range results {
		if result.err == nil {
			created++
			winner = result.actor
		} else if errors.Is(result.err, franchisejourney.ErrConflict) {
			conflicts++
		} else {
			t.Fatal(result.err)
		}
	}
	if created != 1 || conflicts != 7 {
		t.Fatal(created, conflicts)
	}
	got, err := repo.PublishedDeliveryChecklist(ctx, tenant, "store", input.ID, 1)
	if err != nil || got.State != "published" || got.Title != input.Title || len(got.Items) != 2 || got.Items[0].ID != "serial" || got.Items[1].Ordinal != 2 {
		t.Fatal(got, err)
	}
	var rows, items, events int
	var creator, actor, eventID string
	err = pool.QueryRow(ctx, `select (select count(*) from sales.delivery_checklist_template where tenant_id=$1),(select count(*) from sales.delivery_checklist_item where tenant_id=$1),(select count(*) from platform.outbox_event where tenant_id=$1),t.created_by_subject,e.payload->>'actor_subject',e.event_id from sales.delivery_checklist_template t join platform.outbox_event e on e.tenant_id=t.tenant_id and e.aggregate_id=t.checklist_id where t.tenant_id=$1`, tenant).Scan(&rows, &items, &events, &creator, &actor, &eventID)
	if err != nil || rows != 1 || items != 2 || events != 1 || creator != winner || actor != winner {
		t.Fatal(rows, items, events, creator, actor, winner, err)
	}
	for _, scope := range [][3]string{{randomid.Generator{}.New(), "store", input.ID}, {tenant, "other", input.ID}, {tenant, "store", "unknown"}} {
		if v, e := repo.PublishedDeliveryChecklist(ctx, scope[0], scope[1], scope[2], 1); !errors.Is(e, franchisejourney.ErrConflict) || v.ID != "" {
			t.Fatal("scope leak", v, e)
		}
	}
	next := input
	next.Version = 2
	next.Title = "Second version"
	if _, e := repo.PublishDeliveryChecklist(ctx, tenant, "later-actor", next, eventID); e == nil {
		t.Fatal("duplicate outbox ID must roll back")
	}
	if v, e := repo.PublishedDeliveryChecklist(ctx, tenant, "store", input.ID, 2); !errors.Is(e, franchisejourney.ErrConflict) || v.ID != "" {
		t.Fatal("rollback lookup", v, e)
	}
	var remnants int
	if e := pool.QueryRow(ctx, `select (select count(*) from sales.delivery_checklist_template where tenant_id=$1 and checklist_version=2)+(select count(*) from sales.delivery_checklist_item where tenant_id=$1 and checklist_version=2)`, tenant).Scan(&remnants); e != nil || remnants != 0 {
		t.Fatal("partial publication", remnants, e)
	}
	if _, e := repo.PublishDeliveryChecklist(ctx, tenant, "later-actor", next, randomid.Generator{}.New()); e != nil {
		t.Fatal(e)
	}
	for _, sql := range []string{
		`update sales.delivery_checklist_template set title='mutated' where tenant_id=$1 and checklist_version=1`,
		`delete from sales.delivery_checklist_template where tenant_id=$1 and checklist_version=1`,
		`update sales.delivery_checklist_item set prompt='mutated' where tenant_id=$1 and checklist_version=1`,
		`delete from sales.delivery_checklist_item where tenant_id=$1 and checklist_version=1`,
		`insert into sales.delivery_checklist_item(tenant_id,organization_id,checklist_id,checklist_version,item_id,ordinal,prompt,response_type,required) values($1,'store','published-checklist',1,'extra',3,'Extra','text',false)`,
	} {
		if _, e := pool.Exec(ctx, sql, tenant); e == nil {
			t.Fatal("published version accepted mutation")
		}
	}
	if _, e := pool.Exec(ctx, `insert into sales.delivery_checklist_template(tenant_id,organization_id,checklist_id,checklist_version,title,state,created_by_subject) values($1,'store','draft-checklist',1,'Draft','draft','draft-author')`, tenant); e != nil {
		t.Fatal(e)
	}
	if v, e := repo.PublishedDeliveryChecklist(ctx, tenant, "store", "draft-checklist", 1); !errors.Is(e, franchisejourney.ErrConflict) || v.ID != "" {
		t.Fatal("draft exposed", v, e)
	}
	fresh, e := pgxpool.NewWithConfig(ctx, cfg.Copy())
	if e != nil {
		t.Fatal(e)
	}
	defer fresh.Close()
	for version, title := range map[int64]string{1: input.Title, 2: next.Title} {
		v, e := NewFranchiseJourney(fresh).PublishedDeliveryChecklist(ctx, tenant, "store", input.ID, version)
		if e != nil || v.Version != version || v.Title != title || len(v.Items) != 2 || v.Items[0].Prompt != "Serie" {
			t.Fatal("fresh lookup", v, e)
		}
	}
	t.Log("CHECKLIST_PUBLICATION_POSTGRES_PASS created=1 conflicts=7 versions=2 rollback=1 immutable=5 scope=3 draft=blocked fresh_pool=pass")
}

func TestChecklistCompletionIdentityAndRecovery(t *testing.T) {
	raw := os.Getenv("TEST_DATABASE_URL")
	if raw == "" {
		t.Skip("disposable database required")
	}
	cfg, err := pgxpool.ParseConfig(raw)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires disposable loopback database")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := randomid.Generator{}.New()
	for _, q := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'complete-'||replace(($1::uuid)::text,'-',''),'Synthetic','Synthetic')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic','store')`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values($1,'model','model','Synthetic','bicycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'variant','model','variant','Synthetic','{}','active')`,
		`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized)values($1,'customer','Synthetic','synthetic@example.test')`,
		`insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,state,version,received_at)values($1,'stock','store','variant','SYNTHETIC-SERIAL','sold',1,clock_timestamp())`,
		`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'order','store','customer','delivered','ARS',123456,1)`,
		`insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version)values($1,'completion-fixture','store','order','customer','stock','prepared',1)`,
	} {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewFranchiseJourney(pool)
	checklist := franchisejourney.DeliveryChecklist{ID: "completion-checklist", OrganizationID: "store", Version: 1, Title: "Synthetic", Items: []franchisejourney.ChecklistItem{{ID: "serial", Ordinal: 1, Prompt: "Serie", ResponseType: "serial", Required: true}, {ID: "confirmed", Ordinal: 2, Prompt: "Confirmación", ResponseType: "confirmation", Required: true}}}
	if _, err = repo.PublishDeliveryChecklist(ctx, tenant, "fixture-author", checklist, randomid.Generator{}.New()); err != nil {
		t.Fatal(err)
	}
	responses := []franchisejourney.ChecklistResponse{{ItemID: "serial", ResponseText: "SYNTHETIC-SERIAL"}, {ItemID: "confirmed", ResponseText: "confirmed"}}
	type outcome struct {
		actor string
		err   error
	}
	results := make(chan outcome, 8)
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			actor := fmt.Sprintf("completer-%d", i)
			_, e := repo.CompleteDeliveryChecklist(ctx, tenant, "store", actor, "completion-fixture", 1, checklist.ID, 1, responses, randomid.Generator{}.New())
			results <- outcome{actor, e}
		}(i)
	}
	close(start)
	wg.Wait()
	close(results)
	created, conflicts := 0, 0
	winner := ""
	for result := range results {
		if result.err == nil {
			created++
			winner = result.actor
		} else if errors.Is(result.err, franchisejourney.ErrConflict) {
			conflicts++
		} else {
			t.Fatal(result.err)
		}
	}
	if created != 1 || conflicts != 7 {
		t.Fatal(created, conflicts)
	}
	got, err := repo.DeliveryChecklistCompletion(ctx, tenant, "store", "completion-fixture")
	if err != nil || got.Version != 2 || got.State != "presented" || got.ActorSubject != winner || len(got.Responses) != 2 || got.Responses[0].ItemID != "confirmed" || got.Responses[1].ResponseText != "SYNTHETIC-SERIAL" {
		t.Fatal(got, err)
	}
	var events, answers int
	var eventID, actor string
	err = pool.QueryRow(ctx, `select (select count(*) from platform.outbox_event where tenant_id=$1 and event_type='delivery-handover.checklist-completed'),(select count(*) from sales.delivery_checklist_response where tenant_id=$1),event_id,payload->>'actor_subject' from platform.outbox_event where tenant_id=$1 and event_type='delivery-handover.checklist-completed'`, tenant).Scan(&events, &answers, &eventID, &actor)
	if err != nil || events != 1 || answers != 2 || actor != winner {
		t.Fatal(events, answers, actor, err)
	}
	for _, scope := range [][3]string{{randomid.Generator{}.New(), "store", "completion-fixture"}, {tenant, "other", "completion-fixture"}, {tenant, "store", "unknown"}} {
		if v, e := repo.DeliveryChecklistCompletion(ctx, scope[0], scope[1], scope[2]); !errors.Is(e, franchisejourney.ErrConflict) || v.HandoverID != "" {
			t.Fatal("scope", v, e)
		}
	}
	if _, e := pool.Exec(ctx, `insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version) values($1,'rollback-fixture','store','order','customer','stock','prepared',1)`, tenant); e != nil {
		t.Fatal(e)
	}
	if _, e := repo.CompleteDeliveryChecklist(ctx, tenant, "store", "later-actor", "rollback-fixture", 1, checklist.ID, 1, responses, eventID); e == nil {
		t.Fatal("duplicate event must roll back")
	}
	var state string
	var version int64
	if e := pool.QueryRow(ctx, `select state,version,(select count(*) from sales.delivery_checklist_response where tenant_id=$1 and handover_id='rollback-fixture') from sales.delivery_handover where tenant_id=$1 and handover_id='rollback-fixture'`, tenant).Scan(&state, &version, &answers); e != nil || state != "prepared" || version != 1 || answers != 0 {
		t.Fatal("partial completion", state, version, answers, e)
	}
	if v, e := repo.DeliveryChecklistCompletion(ctx, tenant, "store", "rollback-fixture"); !errors.Is(e, franchisejourney.ErrConflict) || v.HandoverID != "" {
		t.Fatal("prepared exposed", v, e)
	}
	if _, e := repo.AcceptHandover(ctx, tenant, "store", "customer", "completion-fixture", 2, "SYNTHETIC-SERIAL", checklist.ID, 1, strings.Repeat("a", 64), randomid.Generator{}.New()); e != nil {
		t.Fatal(e)
	}
	if _, e := repo.CompleteDeliveryChecklist(ctx, tenant, "store", "later-actor", "rollback-fixture", 1, checklist.ID, 1, responses, randomid.Generator{}.New()); e != nil {
		t.Fatal(e)
	}
	if _, e := repo.RejectHandover(ctx, tenant, "store", "customer", "rollback-fixture", 2, "serial-mismatch", "Synthetic rejection", strings.Repeat("b", 64), randomid.Generator{}.New(), randomid.Generator{}.New()); e != nil {
		t.Fatal(e)
	}
	for _, sql := range []string{`update sales.delivery_checklist_response set response_text='changed' where tenant_id=$1 and handover_id='completion-fixture'`, `delete from sales.delivery_checklist_response where tenant_id=$1 and handover_id='completion-fixture'`, `update sales.delivery_handover set checklist_completed_by_subject='forged' where tenant_id=$1 and handover_id='completion-fixture'`} {
		if _, e := pool.Exec(ctx, sql, tenant); e == nil {
			t.Fatal("completed binding accepted mutation")
		}
	}
	for _, sql := range []string{`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'other','other','Other','store')`, `insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,vin,battery_serial_number,state,version,received_at)values($1,'foreign-stock','other','variant','FOREIGN-SYNTHETIC','FOREIGN-SYNTHETIC','FOREIGN-SYNTHETIC','sold',1,clock_timestamp())`, `update sales.delivery_handover set stock_unit_id='foreign-stock' where tenant_id=$1 and handover_id='completion-fixture'`} {
		if _, e := pool.Exec(ctx, sql, tenant); e != nil {
			t.Fatal(e)
		}
	}
	if v, e := repo.DeliveryChecklistCompletion(ctx, tenant, "store", "completion-fixture"); !errors.Is(e, franchisejourney.ErrConflict) || v.HandoverID != "" {
		t.Fatal("foreign stock leaked", v, e)
	}
	if _, e := pool.Exec(ctx, `update sales.delivery_handover set stock_unit_id='stock' where tenant_id=$1 and handover_id='completion-fixture'`, tenant); e != nil {
		t.Fatal(e)
	}
	fresh, e := pgxpool.NewWithConfig(ctx, cfg.Copy())
	if e != nil {
		t.Fatal(e)
	}
	defer fresh.Close()
	for id, want := range map[string]string{"completion-fixture": "accepted", "rollback-fixture": "rejected"} {
		v, e := NewFranchiseJourney(fresh).DeliveryChecklistCompletion(ctx, tenant, "store", id)
		if e != nil || v.State != want || v.Version != 3 || len(v.Responses) != 2 || v.Responses[0].ResponseText != "confirmed" || v.Responses[1].ResponseText != "SYNTHETIC-SERIAL" {
			t.Fatal("current state recovery", v, e)
		}
		if id == "completion-fixture" && (v.ActorSubject != winner || !v.CompletedAt.Equal(got.CompletedAt)) {
			t.Fatal("completion audit changed", v)
		}
	}
	t.Log("CHECKLIST_COMPLETION_POSTGRES_PASS completed=1 conflicts=7 rollback=1 immutable=3 scope=3 foreign_stock=blocked accepted=preserved rejected=preserved actor=preserved fresh_pool=pass")
}

func TestReturnScopeBoundaries(t *testing.T) {
	raw := os.Getenv("TEST_DATABASE_URL")
	if raw == "" {
		t.Skip("disposable database required")
	}
	cfg, err := pgxpool.ParseConfig(raw)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires disposable loopback database")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := randomid.Generator{}.New()
	for _, q := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'complete-'||replace(($1::uuid)::text,'-',''),'Synthetic','Synthetic')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic','store')`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values($1,'model','model','Synthetic','bicycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'variant','model','variant','Synthetic','{}','active')`,
		`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized)values($1,'customer','Synthetic','synthetic@example.test')`,
		`insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,state,version,received_at)values($1,'stock','store','variant','SYNTHETIC-SERIAL','sold',1,clock_timestamp())`,
		`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'order','store','customer','delivered','ARS',123456,1)`,
		`insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version)values($1,'completion-fixture','store','order','customer','stock','prepared',1)`,
	} {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewFranchiseJourney(pool)
	checklist := franchisejourney.DeliveryChecklist{ID: "completion-checklist", OrganizationID: "store", Version: 1, Title: "Synthetic", Items: []franchisejourney.ChecklistItem{{ID: "serial", Ordinal: 1, Prompt: "Serie", ResponseType: "serial", Required: true}, {ID: "confirmed", Ordinal: 2, Prompt: "Confirmación", ResponseType: "confirmation", Required: true}}}
	if _, err = repo.PublishDeliveryChecklist(ctx, tenant, "fixture-author", checklist, randomid.Generator{}.New()); err != nil {
		t.Fatal(err)
	}

	responses := []franchisejourney.ChecklistResponse{{ItemID: "serial", ResponseText: "SYNTHETIC-SERIAL"}, {ItemID: "confirmed", ResponseText: "confirmed"}}
	if _, err = repo.CompleteDeliveryChecklist(ctx, tenant, "store", "fixture-operator", "completion-fixture", 1, checklist.ID, 1, responses, randomid.Generator{}.New()); err != nil {
		t.Fatal(err)
	}
	exception, err := repo.RejectHandover(ctx, tenant, "store", "customer", "completion-fixture", 2, "serial-mismatch", "Synthetic rejection", strings.Repeat("a", 64), randomid.Generator{}.New(), randomid.Generator{}.New())
	if err != nil {
		t.Fatal(err)
	}
	resolution, err := repo.ResolveDeliveryException(ctx, tenant, "store", "fixture-manager", exception.ID, 1, "return", "Synthetic authorization", "", "return-authorization", randomid.Generator{}.New(), randomid.Generator{}.New())
	if err != nil || resolution.ReturnAuthorizationID != "return-authorization" {
		t.Fatal(resolution, err)
	}

	if _, err = pool.Exec(ctx, `insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version) select tenant_id,'scope-second',organization_id,order_id,customer_principal_id,stock_unit_id,'prepared',1 from sales.delivery_handover where tenant_id=$1 and handover_id='completion-fixture'`, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err = repo.CompleteDeliveryChecklist(ctx, tenant, "store", "fixture-operator", "scope-second", 1, checklist.ID, 1, responses, randomid.Generator{}.New()); err != nil {
		t.Fatal(err)
	}
	second, err := repo.RejectHandover(ctx, tenant, "store", "customer", "scope-second", 2, "serial-mismatch", "Synthetic rejection", strings.Repeat("a", 64), randomid.Generator{}.New(), randomid.Generator{}.New())
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.ResolveDeliveryException(ctx, tenant, "store", "fixture-manager", second.ID, 1, "return", "Synthetic authorization", "", "scope-second-authorization", randomid.Generator{}.New(), randomid.Generator{}.New()); err != nil {
		t.Fatal(err)
	}
	received, err := repo.ReceiveReturn(ctx, tenant, "store", "valid-receiver", "scope-second-authorization", "SYNTHETIC-SERIAL", "sealed", "Synthetic receipt", strings.Repeat("b", 64), randomid.Generator{}.New(), randomid.Generator{}.New())
	if err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'other','other','Other','store')`, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `update inventory.stock_unit set organization_id='other' where tenant_id=$1 and stock_unit_id='stock'`, tenant); err != nil {
		t.Fatal(err)
	}
	if values, e := repo.ReturnCases(ctx, tenant, "store", 100); !errors.Is(e, franchisejourney.ErrConflict) || len(values) != 0 {
		t.Errorf("foreign graph listed cases=%d error=%v", len(values), e)
	}
	if value, e := repo.ReceiveReturn(ctx, tenant, "store", "scope-receiver", "return-authorization", "SYNTHETIC-SERIAL", "sealed", "Synthetic receipt", strings.Repeat("b", 64), randomid.Generator{}.New(), randomid.Generator{}.New()); !errors.Is(e, franchisejourney.ErrConflict) || value.ID != "" {
		t.Errorf("foreign graph received id=%s error=%v", value.ID, e)
	}
	if value, e := repo.DecideReturn(ctx, tenant, "store", "scope-decider", received.ID, "quarantine", "Synthetic decision", randomid.Generator{}.New(), randomid.Generator{}.New(), randomid.Generator{}.New(), randomid.Generator{}.New(), randomid.Generator{}.New(), randomid.Generator{}.New()); !errors.Is(e, franchisejourney.ErrConflict) || value.ID != "" {
		t.Errorf("foreign graph decided id=%s effects=%d error=%v", value.ID, len(value.Effects), e)
	}
	var receipts, decisions, effects, events int
	err = pool.QueryRow(ctx, `select (select count(*) from sales.return_receipt where tenant_id=$1 and authorization_id='return-authorization'),(select count(*) from sales.return_disposition where tenant_id=$1),(select count(*) from sales.return_effect_request where tenant_id=$1),(select count(*) from platform.outbox_event where tenant_id=$1 and payload->>'actor_subject' in ('scope-receiver','scope-decider'))`, tenant).Scan(&receipts, &decisions, &effects, &events)
	if err != nil || receipts != 0 || decisions != 0 || effects != 0 || events != 0 {
		t.Errorf("foreign graph effects=%d/%d/%d/%d error=%v", receipts, decisions, effects, events, err)
	}
	if !t.Failed() {
		t.Log("RETURN_SCOPE_PASS foreign_stock_blocks_list_receive_decide effects=0")
	}
}

func TestReturnOperationsIdentityAndRecovery(t *testing.T) {
	raw := os.Getenv("TEST_DATABASE_URL")
	if raw == "" {
		t.Skip("disposable database required")
	}
	cfg, err := pgxpool.ParseConfig(raw)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires disposable loopback database")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := randomid.Generator{}.New()
	for _, q := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'complete-'||replace(($1::uuid)::text,'-',''),'Synthetic','Synthetic')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic','store')`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values($1,'model','model','Synthetic','bicycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'variant','model','variant','Synthetic','{}','active')`,
		`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized)values($1,'customer','Synthetic','synthetic@example.test')`,
		`insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,state,version,received_at)values($1,'stock','store','variant','SYNTHETIC-SERIAL','sold',1,clock_timestamp())`,
		`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'order','store','customer','delivered','ARS',123456,1)`,
		`insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version)values($1,'completion-fixture','store','order','customer','stock','prepared',1)`,
	} {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewFranchiseJourney(pool)
	checklist := franchisejourney.DeliveryChecklist{ID: "completion-checklist", OrganizationID: "store", Version: 1, Title: "Synthetic", Items: []franchisejourney.ChecklistItem{{ID: "serial", Ordinal: 1, Prompt: "Serie", ResponseType: "serial", Required: true}, {ID: "confirmed", Ordinal: 2, Prompt: "Confirmación", ResponseType: "confirmation", Required: true}}}
	if _, err = repo.PublishDeliveryChecklist(ctx, tenant, "fixture-author", checklist, randomid.Generator{}.New()); err != nil {
		t.Fatal(err)
	}

	responses := []franchisejourney.ChecklistResponse{{ItemID: "serial", ResponseText: "SYNTHETIC-SERIAL"}, {ItemID: "confirmed", ResponseText: "confirmed"}}
	if _, err = repo.CompleteDeliveryChecklist(ctx, tenant, "store", "fixture-operator", "completion-fixture", 1, checklist.ID, 1, responses, randomid.Generator{}.New()); err != nil {
		t.Fatal(err)
	}
	exception, err := repo.RejectHandover(ctx, tenant, "store", "customer", "completion-fixture", 2, "serial-mismatch", "Synthetic rejection", strings.Repeat("a", 64), randomid.Generator{}.New(), randomid.Generator{}.New())
	if err != nil {
		t.Fatal(err)
	}
	resolution, err := repo.ResolveDeliveryException(ctx, tenant, "store", "fixture-manager", exception.ID, 1, "return", "Synthetic authorization", "", "return-authorization", randomid.Generator{}.New(), randomid.Generator{}.New())
	if err != nil || resolution.ReturnAuthorizationID != "return-authorization" {
		t.Fatal(resolution, err)
	}

	newID := func() string { return randomid.Generator{}.New() }
	authorize := func(id, action string) {
		t.Helper()
		if _, e := pool.Exec(ctx, `insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version)values($1,$2,'store','order','customer','stock','prepared',1)`, tenant, id); e != nil {
			t.Fatal(e)
		}
		if _, e := repo.CompleteDeliveryChecklist(ctx, tenant, "store", "fixture", id, 1, checklist.ID, 1, responses, newID()); e != nil {
			t.Fatal(e)
		}
		x, e := repo.RejectHandover(ctx, tenant, "store", "customer", id, 2, "serial-mismatch", "Synthetic", strings.Repeat("a", 64), newID(), newID())
		if e != nil {
			t.Fatal(e)
		}
		if _, e = repo.ResolveDeliveryException(ctx, tenant, "store", "fixture", x.ID, 1, action, "Synthetic", "", id, newID(), newID()); e != nil {
			t.Fatal(e)
		}
	}
	receive := func(auth, actor, event string) (franchisejourney.ReturnReceipt, error) {
		return repo.ReceiveReturn(ctx, tenant, "store", actor, auth, "SYNTHETIC-SERIAL", "sealed", "Synthetic receipt", strings.Repeat("b", 64), newID(), event)
	}
	decide := func(receipt, actor, event string) (franchisejourney.ReturnDisposition, error) {
		return repo.DecideReturn(ctx, tenant, "store", actor, receipt, "quarantine", "Synthetic decision", newID(), newID(), newID(), newID(), newID(), event)
	}
	lookup := func(auth string) franchisejourney.ReturnCase {
		t.Helper()
		v, e := repo.ReturnCaseResult(ctx, tenant, "store", auth)
		if e != nil {
			t.Fatal(e)
		}
		return v
	}
	if v := lookup("return-authorization"); v.Receipt != nil || v.Disposition != nil {
		t.Fatal("uncreated operations exposed")
	}
	if v, e := repo.ReceiveReturn(ctx, tenant, "store", "wrong-serial", "return-authorization", "WRONG", "sealed", "Synthetic", strings.Repeat("b", 64), newID(), newID()); !errors.Is(e, franchisejourney.ErrConflict) || v.ID != "" {
		t.Fatal(v, e)
	}
	type outcome struct {
		id  string
		err error
	}
	race := func(fn func(int) (string, error)) string {
		t.Helper()
		ch := make(chan outcome, 8)
		start := make(chan struct{})
		var wg sync.WaitGroup
		for i := 0; i < 8; i++ {
			wg.Add(1)
			go func(i int) { defer wg.Done(); <-start; id, e := fn(i); ch <- outcome{id, e} }(i)
		}
		close(start)
		wg.Wait()
		close(ch)
		wins, conflicts, id := 0, 0, ""
		for v := range ch {
			if v.err == nil {
				wins++
				id = v.id
			} else if errors.Is(v.err, franchisejourney.ErrConflict) {
				conflicts++
			} else {
				t.Fatal(v.err)
			}
		}
		if wins != 1 || conflicts != 7 || id == "" {
			t.Fatal(wins, conflicts, id)
		}
		return id
	}
	receiptID := race(func(i int) (string, error) {
		v, e := receive("return-authorization", fmt.Sprintf("receiver-%d", i), newID())
		return v.ID, e
	})
	decisionID := race(func(i int) (string, error) {
		v, e := decide(receiptID, fmt.Sprintf("decider-%d", i), newID())
		return v.ID, e
	})
	first := lookup("return-authorization")
	if first.Receipt.ID != receiptID || first.Disposition.ID != decisionID || len(first.Disposition.Effects) != 4 {
		t.Fatal(first)
	}
	for _, actor := range []string{first.Receipt.ReceivedBySubject, first.Disposition.DecidedBySubject} {
		var n int
		if e := pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and payload->>'actor_subject'=$2`, tenant, actor).Scan(&n); e != nil || n != 1 {
			t.Fatal(actor, n, e)
		}
	}
	authorize("exchange-authorization", "exchange")
	er, e := receive("exchange-authorization", "exchange-receiver", newID())
	if e != nil {
		t.Fatal(e)
	}
	ed, e := decide(er.ID, "exchange-decider", newID())
	if e != nil || len(ed.Effects) != 3 || ed.CustomerRemedy != "exchange" {
		t.Fatal(ed, e)
	}
	for _, effect := range ed.Effects {
		if effect.OwnerContext == "fiscal" {
			t.Fatal("exchange created fiscal request")
		}
	}
	authorize("atomic-authorization", "return")
	var occupiedEvent string
	if e = pool.QueryRow(ctx, `select event_id from platform.outbox_event where tenant_id=$1 limit 1`, tenant).Scan(&occupiedEvent); e != nil {
		t.Fatal(e)
	}
	if v, e := receive("atomic-authorization", "atomic-receiver", occupiedEvent); e == nil || v.ID != "" {
		t.Fatal("outbox collision committed receipt", v, e)
	}
	if v := lookup("atomic-authorization"); v.Receipt != nil {
		t.Fatal("receipt survived rollback")
	}
	ar, e := receive("atomic-authorization", "atomic-receiver", newID())
	if e != nil {
		t.Fatal(e)
	}
	if v, e := decide(ar.ID, "atomic-decider", occupiedEvent); e == nil || v.ID != "" {
		t.Fatal("outbox collision committed decision", v, e)
	}
	if v := lookup("atomic-authorization"); v.Disposition != nil {
		t.Fatal("decision survived rollback")
	}
	var requests int
	if e = pool.QueryRow(ctx, `select count(*) from sales.return_effect_request where tenant_id=$1`, tenant).Scan(&requests); e != nil || requests != 7 {
		t.Fatal("requests survived rollback", requests, e)
	}
	if _, e = decide(ar.ID, "atomic-decider", newID()); e != nil {
		t.Fatal(e)
	}
	for _, q := range []string{`update sales.return_receipt set notes='changed' where tenant_id=$1`, `delete from sales.return_receipt where tenant_id=$1`, `update sales.return_disposition set notes='changed' where tenant_id=$1`, `delete from sales.return_disposition where tenant_id=$1`, `update sales.return_effect_request set state='changed' where tenant_id=$1`, `delete from sales.return_effect_request where tenant_id=$1`} {
		if _, e = pool.Exec(ctx, q, tenant); e == nil {
			t.Fatal("immutable evidence changed", q)
		}
	}
	// More than the historical list limit: absence from a page is not absence of a committed operation.
	for i := 0; i < 101; i++ {
		authorize(fmt.Sprintf("newer-%03d", i), "return")
	}
	list, e := repo.ReturnCases(ctx, tenant, "store", 100)
	if e != nil || len(list) != 100 {
		t.Fatal(len(list), e)
	}
	for _, v := range list {
		if v.AuthorizationID == "return-authorization" {
			t.Fatal("old case unexpectedly in newest page")
		}
	}
	fresh, e := pgxpool.NewWithConfig(ctx, cfg.Copy())
	if e != nil {
		t.Fatal(e)
	}
	defer fresh.Close()
	again, e := NewFranchiseJourney(fresh).ReturnCaseResult(ctx, tenant, "store", "return-authorization")
	if e != nil || again.Receipt.ID != receiptID || again.Disposition.ID != decisionID || again.Receipt.EvidenceSHA256 != first.Receipt.EvidenceSHA256 || again.Receipt.ReceivedBySubject != first.Receipt.ReceivedBySubject {
		t.Fatal("fresh pool lost evidence", again, e)
	}
	for _, scope := range [][3]string{{tenant, "other", "return-authorization"}, {newID(), "store", "return-authorization"}, {tenant, "store", "unknown"}} {
		v, e := repo.ReturnCaseResult(ctx, scope[0], scope[1], scope[2])
		if !errors.Is(e, franchisejourney.ErrConflict) || v.AuthorizationID != "" {
			t.Fatal(v, e)
		}
	}
	// The scope lock blocks concurrent movement of the authorization's handover, order and stock.
	lock, e := pool.Begin(ctx)
	if e != nil {
		t.Fatal(e)
	}
	defer lock.Rollback(ctx)
	if e = lockReturnAuthorizationScope(ctx, lock, tenant, "store", "return-authorization"); e != nil {
		t.Fatal(e)
	}
	for _, q := range []string{`update sales.delivery_handover set version=version+1 where tenant_id=$1 and handover_id='completion-fixture'`, `update sales.customer_order set version=version+1 where tenant_id=$1 and order_id='order'`, `update inventory.stock_unit set version=version+1 where tenant_id=$1 and stock_unit_id='stock'`} {
		writer, e := fresh.Begin(ctx)
		if e != nil {
			t.Fatal(e)
		}
		if _, e = writer.Exec(ctx, `set local lock_timeout='100ms'`); e != nil {
			t.Fatal(e)
		}
		_, e = writer.Exec(ctx, q, tenant)
		var pe *pgconn.PgError
		if !errors.As(e, &pe) || pe.Code != "55P03" {
			t.Fatal("scope lock failed", e)
		}
		writer.Rollback(ctx)
	}
	lock.Rollback(ctx)
	t.Log("RETURN_IDENTITY_RECOVERY_PASS receipt_race=1/8 decision_race=1/8 requests=4+3+4 atomic_rollbacks=2 immutable=6 outside_limit=100 fresh_pool=PASS scope_locks=3")
}
