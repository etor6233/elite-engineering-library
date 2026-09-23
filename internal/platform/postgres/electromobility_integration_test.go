package postgres

import (
	"context"
	"elite.local/enterprise/internal/electromobility"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"testing"
)

func TestElectromobilityRepositoryCatalogAndLead(t *testing.T) {
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
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c28e01"
	org := "em-api-org"
	cleanup := func() {
		_, _ = pool.Exec(ctx, `delete from platform.outbox_event where tenant_id=$1`, tenant)
		_, _ = pool.Exec(ctx, `delete from platform.idempotency_record where tenant_id=$1`, tenant)
		_, _ = pool.Exec(ctx, `delete from crm.consent_evidence where tenant_id=$1`, tenant)
		_, _ = pool.Exec(ctx, `delete from crm.lead where tenant_id=$1`, tenant)
		_, _ = pool.Exec(ctx, `delete from catalog.vehicle_model where tenant_id=$1`, tenant)
		_, _ = pool.Exec(ctx, `delete from org.organization where tenant_id=$1`, tenant)
		_, _ = pool.Exec(ctx, `delete from platform.tenant where tenant_id=$1`, tenant)
	}
	cleanup()
	defer cleanup()
	_, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'em-api','EM API','EM')`, tenant)
	if err != nil {
		t.Fatal(err)
	}
	_, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,$2,'central-store','Central','store')`, tenant, org)
	if err != nil {
		t.Fatal(err)
	}
	repo := NewElectromobility(pool)
	model := electromobility.Model{ID: "model-1", Code: "urban-one", DisplayName: "Urban One", VehicleClass: "bicycle", Specification: json.RawMessage(`{"range_km":80}`)}
	if err := repo.CreateModel(ctx, tenant, "018f4d4a-7b36-7a21-8d10-2f4c54c28e02", model); err != nil {
		t.Fatal(err)
	}
	_, err = pool.Exec(ctx, `update catalog.vehicle_model set lifecycle_state='active',publicly_visible=true where tenant_id=$1 and model_id=$2`, tenant, model.ID)
	if err != nil {
		t.Fatal(err)
	}
	models, err := repo.ListPublicModels(ctx, "em-api")
	if err != nil || len(models) != 1 {
		t.Fatalf("models=%+v err=%v", models, err)
	}
	tenantID, orgID, err := repo.ResolvePublicOrganization(ctx, "em-api", "central-store")
	if err != nil || tenantID != tenant || orgID != org {
		t.Fatalf("resolve %s %s %v", tenantID, orgID, err)
	}
	lead := electromobility.Lead{ID: "lead-1", TenantID: tenant, OrganizationID: org, ModelID: model.ID, SourceCode: "public-web", Contact: json.RawMessage(`{"email":"person@example.test"}`), ConsentID: "consent-1", ConsentEvidenceHash: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}
	resourceID, replayed, err := repo.CreateLead(ctx, lead, "018f4d4a-7b36-7a21-8d10-2f4c54c28e03", "lead-request-00000001")
	if err != nil || replayed || resourceID != lead.ID {
		t.Fatal(err)
	}
	resourceID, replayed, err = repo.CreateLead(ctx, lead, "018f4d4a-7b36-7a21-8d10-2f4c54c28e04", "lead-request-00000001")
	if err != nil || !replayed || resourceID != lead.ID {
		t.Fatalf("replay id=%s replayed=%t err=%v", resourceID, replayed, err)
	}
	changed := lead
	changed.ConsentEvidenceHash = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	if _, _, err = repo.CreateLead(ctx, changed, "018f4d4a-7b36-7a21-8d10-2f4c54c28e05", "lead-request-00000001"); !errors.Is(err, electromobility.ErrConflict) {
		t.Fatalf("changed request err=%v", err)
	}
	var leads, consents, events int
	_ = pool.QueryRow(ctx, `select count(*) from crm.lead where tenant_id=$1`, tenant).Scan(&leads)
	_ = pool.QueryRow(ctx, `select count(*) from crm.consent_evidence where tenant_id=$1`, tenant).Scan(&consents)
	_ = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1`, tenant).Scan(&events)
	if leads != 1 || consents != 1 || events != 2 {
		t.Fatalf("lead=%d consent=%d events=%d", leads, consents, events)
	}
}
