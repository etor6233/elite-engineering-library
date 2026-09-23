package postgres

import (
	"context"
	"elite.local/enterprise/internal/fulfillment"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"sync"
	"testing"
	"time"
)

func TestFulfillmentServiceFranchiseFlow(t *testing.T) {
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
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c28701"
	cleanup := func() {
		for _, q := range []string{`delete from platform.outbox_event where tenant_id=$1`, `delete from communication.message where tenant_id=$1`, `delete from franchise.agreement where tenant_id=$1`, `delete from service_ops.recall_unit where tenant_id=$1`, `delete from service_ops.recall where tenant_id=$1`, `delete from service_ops.service_case where tenant_id=$1`, `delete from logistics.shipment where tenant_id=$1`, `delete from inventory.stock_unit where tenant_id=$1`, `delete from catalog.vehicle_variant where tenant_id=$1`, `delete from catalog.vehicle_model where tenant_id=$1`, `delete from crm.customer_profile where tenant_id=$1`, `delete from org.organization where tenant_id=$1`, `delete from platform.tenant where tenant_id=$1`} {
			_, _ = pool.Exec(ctx, q, tenant)
		}
	}
	cleanup()
	defer cleanup()
	fixtures := []string{`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'fulfillment-api','Fulfillment','Fulfillment')`, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'org-a','org-a','A','franchisee'),($1,'org-b','org-b','B','service_center')`, `insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,status)values($1,'customer','Customer','active')`, `insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values($1,'model','model','Model','bicycle','active')`, `insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'variant','model','variant','Variant','{}','active')`, `insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,state,version,received_at)values($1,'stock','org-b','variant','FULFILLMENT-SERIAL','available',1,clock_timestamp())`}
	for _, q := range fixtures {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	r := NewFulfillment(pool)
	n := 0
	event := func() string { n++; return "018f4d4a-7b36-7a21-8d10-2f4c54c28" + fmt.Sprintf("%03d", 700+n) }
	if err = r.CreateShipment(ctx, tenant, event(), fulfillment.Shipment{ID: "shipment", ProviderCode: "carrier", OriginOrganizationID: "org-a", DestinationOrganizationID: "org-b"}); err != nil {
		t.Fatal(err)
	}
	if err = r.TransitionShipment(ctx, tenant, "org-other", "org-b", "shipment", "planned", "booked", "carrier-ref", event()); err == nil {
		t.Fatal("shipment crossed organization scope")
	}
	if err = r.TransitionShipment(ctx, tenant, "org-a", "org-b", "shipment", "planned", "booked", "carrier-ref", event()); err != nil {
		t.Fatal(err)
	}
	c := fulfillment.ServiceCase{ID: "case", StockUnitID: "stock", OrganizationID: "org-b", Severity: "safety", Description: "brake inspection"}
	if err = r.OpenServiceCase(ctx, tenant, event(), c); err != nil {
		t.Fatal(err)
	}
	if err = r.TransitionServiceCase(ctx, tenant, "org-a", "case", "opened", "diagnosis", 1, event()); err == nil {
		t.Fatal("service case crossed organization scope")
	}
	if err = r.TransitionServiceCase(ctx, tenant, "org-b", "case", "opened", "diagnosis", 1, event()); err != nil {
		t.Fatal(err)
	}
	if err = r.TransitionServiceCase(ctx, tenant, "org-b", "case", "opened", "cancelled", 1, event()); err == nil {
		t.Fatal("stale service transition succeeded")
	}
	if err = r.CreateRecall(ctx, tenant, event(), fulfillment.Recall{ID: "recall", Code: "REC-1", Title: "Safety inspection", Severity: "safety"}); err != nil {
		t.Fatal(err)
	}
	if err = r.ActivateRecall(ctx, tenant, "recall", event()); err != nil {
		t.Fatal(err)
	}
	if err = r.AddRecallUnit(ctx, tenant, "org-a", "recall", "stock", event()); err == nil {
		t.Fatal("recall unit crossed organization scope")
	}
	if err = r.AddRecallUnit(ctx, tenant, "org-b", "recall", "stock", event()); err != nil {
		t.Fatal(err)
	}
	if err = r.AddRecallUnit(ctx, tenant, "org-b", "recall", "stock", event()); err == nil {
		t.Fatal("duplicate recall unit succeeded")
	}
	if err = r.CreateOrganization(ctx, tenant, event(), fulfillment.Organization{ID: "store-child", ParentOrganizationID: "org-a", Code: "store-child", DisplayName: "Store Child", Type: "store", Status: "provisioning", Version: 1}); err != nil {
		t.Fatal(err)
	}
	if err = r.TransitionOrganization(ctx, tenant, "store-child", "provisioning", "active", 1, event()); err != nil {
		t.Fatal(err)
	}
	if err = r.TransitionOrganization(ctx, tenant, "store-child", "provisioning", "closed", 1, event()); err == nil {
		t.Fatal("stale organization transition succeeded")
	}
	if err = r.TransitionOrganization(ctx, tenant, "org-a", "active", "closed", 1, event()); err == nil {
		t.Fatal("parent with active child was closed")
	}
	starts := time.Now().UTC()
	if err = r.CreateAgreement(ctx, tenant, event(), fulfillment.Agreement{ID: "agreement", OrganizationID: "org-a", TerritoryCode: "AR-CBA", TermsVersion: "v1", StartsOn: starts}); err != nil {
		t.Fatal(err)
	}
	if err = r.TransitionAgreement(ctx, tenant, "org-b", "agreement", "draft", "active", 1, event()); err == nil {
		t.Fatal("agreement crossed organization scope")
	}
	if err = r.TransitionAgreement(ctx, tenant, "org-a", "agreement", "draft", "active", 1, event()); err != nil {
		t.Fatal(err)
	}
	if err = r.TransitionAgreement(ctx, tenant, "org-a", "agreement", "draft", "terminated", 1, event()); err == nil {
		t.Fatal("stale agreement transition succeeded")
	}
	if err = r.CreateAgreement(ctx, tenant, event(), fulfillment.Agreement{ID: "overlap-agreement", OrganizationID: "org-a", TerritoryCode: "AR-CBA", TermsVersion: "v2", StartsOn: starts.AddDate(0, 1, 0)}); err != nil {
		t.Fatal(err)
	}
	if err = r.TransitionAgreement(ctx, tenant, "org-a", "overlap-agreement", "draft", "active", 1, event()); err == nil {
		t.Fatal("overlapping active territory succeeded")
	}
	for _, agreement := range []fulfillment.Agreement{
		{ID: "concurrent-agreement-a", OrganizationID: "org-a", TerritoryCode: "AR-SFE", TermsVersion: "v1", StartsOn: starts},
		{ID: "concurrent-agreement-b", OrganizationID: "org-a", TerritoryCode: "AR-SFE", TermsVersion: "v2", StartsOn: starts.AddDate(0, 1, 0)},
	} {
		if err = r.CreateAgreement(ctx, tenant, event(), agreement); err != nil {
			t.Fatal(err)
		}
	}
	events := []string{event(), event()}
	start := make(chan struct{})
	results := make(chan error, 2)
	var workers sync.WaitGroup
	for index, agreementID := range []string{"concurrent-agreement-a", "concurrent-agreement-b"} {
		workers.Add(1)
		go func(id, eventID string) {
			defer workers.Done()
			<-start
			results <- r.TransitionAgreement(ctx, tenant, "org-a", id, "draft", "active", 1, eventID)
		}(agreementID, events[index])
	}
	close(start)
	workers.Wait()
	close(results)
	succeeded, conflicted := 0, 0
	for result := range results {
		if result == nil {
			succeeded++
		} else if errors.Is(result, fulfillment.ErrConflict) {
			conflicted++
		} else {
			t.Fatal(result)
		}
	}
	if succeeded != 1 || conflicted != 1 {
		t.Fatalf("concurrent territory activation success=%d conflict=%d", succeeded, conflicted)
	}
	if err = r.QueueMessage(ctx, tenant, event(), fulfillment.Message{ID: "message", RecipientPrincipalID: "customer", Channel: "email", TemplateCode: "recall", TemplateVersion: "v1"}); err != nil {
		t.Fatal(err)
	}
	if err = r.TransitionMessage(ctx, tenant, "message", "queued", "sent", "provider-ref", "", event()); err != nil {
		t.Fatal(err)
	}
	if err = r.TransitionMessage(ctx, tenant, "message", "sent", "delivered", "provider-ref", "", event()); err != nil {
		t.Fatal(err)
	}
	var outboxCount int
	if err = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1`, tenant).Scan(&outboxCount); err != nil || outboxCount != 18 {
		t.Fatalf("outbox=%d err=%v", outboxCount, err)
	}
}
