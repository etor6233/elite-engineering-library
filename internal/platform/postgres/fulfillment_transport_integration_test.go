package postgres

import (
	"context"
	"elite.local/enterprise/internal/fulfillment"
	"elite.local/enterprise/internal/inventorycontrol"
	"errors"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestConnectedCarrierDeliveryIsIdempotentConcurrentAndNeverAcceptsHandover(t *testing.T) {
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
	tenant := bulkTestUUID(t)
	defer cleanupConnectedTransport(t, pool, tenant)
	for _, command := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Connected Carrier','Connected Carrier')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'warehouse','warehouse','Warehouse','warehouse')`,
		`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,status) values($1,'customer','Customer','active')`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state) values($1,'model','MODEL','Model','other','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,homologation_state,lifecycle_state) values($1,'variant','model','VARIANT','Variant','{}','approved','active')`,
		`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version) values($1,'order','warehouse','customer','placed','USD',100,1)`,
		`insert into sales.customer_order_line(tenant_id,order_id,line_id,variant_id,quantity,unit_price_minor_units) values($1,'order','line','variant',1,100)`,
	} {
		args := []any{tenant}
		if strings.Contains(command, "$2") {
			args = append(args, "carrier-"+tenant[:8])
		}
		if _, err = pool.Exec(ctx, command, args...); err != nil {
			t.Fatal(err)
		}
	}
	inventory := NewInventoryControl(pool)
	if _, err = inventory.CreateBulkItem(ctx, tenant, bulkTestUUID(t), inventorycontrol.BulkItem{ID: "item", Code: "ITEM", Description: "Item", BaseUOM: "EA", BaseRoundingPrecision: "1", TrackingMode: "none", CostingMethod: "fifo", Version: 1}); err != nil {
		t.Fatal(err)
	}
	for _, bin := range []inventorycontrol.WarehouseBin{{ID: "receive", OrganizationID: "warehouse", Code: "RECEIVE", Type: "receive", Version: 1}, {ID: "pick", OrganizationID: "warehouse", Code: "PICK", Type: "putpick", Ranking: 100, Version: 1}, {ID: "ship", OrganizationID: "warehouse", Code: "SHIP", Type: "ship", Version: 1}} {
		if _, err = inventory.CreateWarehouseBin(ctx, tenant, bulkTestUUID(t), bin); err != nil {
			t.Fatal(err)
		}
	}
	for _, policy := range []inventorycontrol.ItemBinPolicy{{OrganizationID: "warehouse", ItemID: "item", BinID: "receive", Fixed: true, MinQuantity: "0", MaxQuantity: "1", Version: 1}, {OrganizationID: "warehouse", ItemID: "item", BinID: "pick", Fixed: true, Default: true, MinQuantity: "0", MaxQuantity: "1", Version: 1}, {OrganizationID: "warehouse", ItemID: "item", BinID: "ship", Fixed: true, MinQuantity: "0", MaxQuantity: "1", Version: 1}} {
		if _, err = inventory.ConfigureItemBin(ctx, tenant, bulkTestUUID(t), policy); err != nil {
			t.Fatal(err)
		}
	}
	posting := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	receipt, err := inventory.PostWarehouseReceipt(ctx, tenant, warehousePackagingIDs(t, "carrier-source"), inventorycontrol.WarehouseReceiptCommand{RequestID: "carrier-source-receipt", OrganizationID: "warehouse", ReceiveBinID: "receive", ItemID: "item", HandlingUOM: "EA", HandlingQuantity: "1", UnitCost: "10", PostingDate: posting, SourceKind: "purchase", SourceID: "carrier-source-po"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = inventory.RegisterWarehouseActivity(ctx, tenant, "warehouse", receipt.PutAway.ID, 1, posting.AddDate(0, 0, 1), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	if _, err = inventory.ConfigureSalesWarehouseBinding(ctx, tenant, bulkTestUUID(t), inventorycontrol.SalesWarehouseBinding{RequestID: "carrier-binding", VariantID: "variant", ItemID: "item", SalesUOMCode: "EA", Version: 1}); err != nil {
		t.Fatal(err)
	}
	pick, err := inventory.CreateWarehousePick(ctx, tenant, warehousePackagingIDs(t, "carrier-pick"), inventorycontrol.WarehousePickCommand{RequestID: "carrier-pick", OrganizationID: "warehouse", ShipBinID: "ship", ItemID: "item", DemandKind: "customer-order", DemandID: "order", DemandLineID: "line", Quantity: "1", HandlingUOM: "EA"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = inventory.RegisterWarehouseActivity(ctx, tenant, "warehouse", pick.ID, 1, posting.AddDate(0, 0, 2), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	if _, err = inventory.PostCustomerShipment(ctx, tenant, inventorycontrol.CustomerShipmentIDs{ShipmentID: "sales-shipment", LineID: "sales-shipment-line", EventID: bulkTestUUID(t)}, inventorycontrol.CustomerShipmentCommand{RequestID: "sales-shipment-request", OrganizationID: "warehouse", OrderID: "order", WarehouseActivityID: pick.ID, PostingDate: posting.AddDate(0, 0, 3)}); err != nil {
		t.Fatal(err)
	}
	transportRepo := NewFulfillment(pool)
	command := fulfillment.CustomerTransportCommand{RequestID: "transport-request", CustomerShipmentID: "sales-shipment", OriginOrganizationID: "warehouse", ProviderCode: "amazon-easyship", ProviderServiceCode: "standard", ProviderReference: "tracking-1"}
	transport, err := transportRepo.CreateCustomerTransport(ctx, tenant, fulfillment.TransportIDs{ShipmentID: "transport", EventID: bulkTestUUID(t)}, command)
	if err != nil || transport.State != "booked" || transport.OrderFulfillmentState != "shipped" {
		t.Fatalf("transport=%+v err=%v", transport, err)
	}
	replay, err := transportRepo.CreateCustomerTransport(ctx, tenant, fulfillment.TransportIDs{ShipmentID: "ignored", EventID: bulkTestUUID(t)}, command)
	if err != nil || replay.ID != transport.ID {
		t.Fatalf("transport replay=%+v err=%v", replay, err)
	}
	divergent := command
	divergent.ProviderReference = "tracking-other"
	if _, err = transportRepo.CreateCustomerTransport(ctx, tenant, fulfillment.TransportIDs{ShipmentID: "ignored-divergent", EventID: bulkTestUUID(t)}, divergent); !errors.Is(err, fulfillment.ErrConflict) {
		t.Fatalf("divergent transport error=%v", err)
	}
	baseReport := fulfillment.ProviderReportCommand{OriginOrganizationID: "warehouse", ProviderEventID: "pickup-event", ProviderStatus: "PickedUp", ReportSchema: fulfillment.AmazonEasyShipReconciliationSchema, EvidenceSHA256: strings.Repeat("a", 64), OccurredAt: posting.AddDate(0, 0, 4), Version: 1}
	picked, err := transportRepo.RecordProviderReport(ctx, tenant, transport.ID, bulkTestUUID(t), baseReport)
	if err != nil || picked.State != "picked-up" || picked.Version != 2 {
		t.Fatalf("picked=%+v err=%v", picked, err)
	}
	start := make(chan struct{})
	results := make(chan error, 2)
	var workers sync.WaitGroup
	for index, status := range []string{"AtOriginFC", "AtDestinationFC"} {
		workers.Add(1)
		go func(index int, status string) {
			defer workers.Done()
			<-start
			report := fulfillment.ProviderReportCommand{OriginOrganizationID: "warehouse", ProviderEventID: "transit-event-" + string(rune('a'+index)), ProviderStatus: status, ReportSchema: fulfillment.AmazonEasyShipReconciliationSchema, EvidenceSHA256: strings.Repeat(string(rune('b'+index)), 64), OccurredAt: posting.AddDate(0, 0, 5), Version: 2}
			_, reportErr := transportRepo.RecordProviderReport(ctx, tenant, transport.ID, bulkTestUUID(t), report)
			results <- reportErr
		}(index, status)
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
		t.Fatalf("provider concurrency success=%d conflict=%d", succeeded, conflicted)
	}
	outForDelivery := fulfillment.ProviderReportCommand{OriginOrganizationID: "warehouse", ProviderEventID: "out-for-delivery-event", ProviderStatus: "OutForDelivery", ReportSchema: fulfillment.AmazonEasyShipReconciliationSchema, EvidenceSHA256: strings.Repeat("d", 64), OccurredAt: posting.AddDate(0, 0, 6), Version: 3}
	if _, err = transportRepo.RecordProviderReport(ctx, tenant, transport.ID, bulkTestUUID(t), outForDelivery); err != nil {
		t.Fatal(err)
	}
	delivery := fulfillment.ProviderReportCommand{OriginOrganizationID: "warehouse", ProviderEventID: "delivery-event", ProviderStatus: "Delivered", ReportSchema: fulfillment.AmazonEasyShipReconciliationSchema, EvidenceSHA256: strings.Repeat("e", 64), OccurredAt: posting.AddDate(0, 0, 7), Version: 4}
	delivered, err := transportRepo.RecordProviderReport(ctx, tenant, transport.ID, bulkTestUUID(t), delivery)
	if err != nil || delivered.State != "delivered" || delivered.Version != 5 || delivered.OrderFulfillmentState != "delivered" {
		t.Fatalf("delivered=%+v err=%v", delivered, err)
	}
	replayedDelivery, err := transportRepo.RecordProviderReport(ctx, tenant, transport.ID, bulkTestUUID(t), delivery)
	if err != nil || replayedDelivery.Version != 5 {
		t.Fatalf("delivery replay=%+v err=%v", replayedDelivery, err)
	}
	divergentDelivery := delivery
	divergentDelivery.EvidenceSHA256 = strings.Repeat("f", 64)
	if _, err = transportRepo.RecordProviderReport(ctx, tenant, transport.ID, bulkTestUUID(t), divergentDelivery); !errors.Is(err, fulfillment.ErrConflict) {
		t.Fatalf("divergent delivery error=%v", err)
	}
	var handovers, reports int
	if err = pool.QueryRow(ctx, `select count(*) from sales.delivery_handover where tenant_id=$1`, tenant).Scan(&handovers); err != nil || handovers != 0 {
		t.Fatalf("provider created business handover=%d err=%v", handovers, err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from logistics.shipment_provider_report where tenant_id=$1`, tenant).Scan(&reports); err != nil || reports != 4 {
		t.Fatalf("provider reports=%d err=%v", reports, err)
	}
}

func cleanupConnectedTransport(t *testing.T, pool *pgxpool.Pool, tenant string) {
	t.Helper()
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Errorf("connected transport cleanup begin: %v", err)
		return
	}
	defer tx.Rollback(ctx)
	commands := []string{
		`alter table logistics.shipment_provider_report disable trigger shipment_provider_report_immutable`,
		`alter table inventory.customer_shipment_allocation disable trigger customer_shipment_allocation_immutable`,
		`alter table sales.customer_shipment_line disable trigger customer_shipment_line_immutable`,
		`alter table sales.customer_shipment disable trigger customer_shipment_immutable`,
		`alter table inventory.warehouse_pick_request disable trigger warehouse_pick_request_immutable`,
		`alter table inventory.sales_warehouse_binding disable trigger sales_warehouse_binding_immutable`,
		`alter table inventory.warehouse_activity_uom_reservation disable trigger warehouse_uom_reservation_delete_guard`,
		`alter table inventory.warehouse_activity_line disable trigger warehouse_activity_line_immutable`,
		`alter table inventory.warehouse_receipt_line disable trigger warehouse_receipt_line_immutable`,
		`alter table inventory.warehouse_receipt disable trigger warehouse_receipt_immutable`,
		`alter table inventory.bulk_inventory_entry disable trigger bulk_inventory_entry_immutable`,
		`alter table inventory.bulk_cost_application disable trigger bulk_cost_application_immutable`,
		`alter table inventory.item_unit_of_measure disable trigger item_unit_of_measure_immutable`,
		`delete from logistics.shipment_provider_report where tenant_id=$1`,
		`delete from logistics.shipment where tenant_id=$1`,
		`delete from inventory.warehouse_pick_request where tenant_id=$1`,
		`delete from inventory.customer_shipment_allocation where tenant_id=$1`,
		`delete from sales.customer_shipment_line where tenant_id=$1`,
		`delete from sales.customer_shipment where tenant_id=$1`,
		`delete from inventory.warehouse_activity_uom_reservation where tenant_id=$1`,
		`delete from inventory.warehouse_activity_line where tenant_id=$1`,
		`delete from inventory.warehouse_activity where tenant_id=$1`,
		`delete from inventory.warehouse_receipt_line where tenant_id=$1`,
		`delete from inventory.warehouse_receipt where tenant_id=$1`,
		`delete from inventory.bulk_uom_conversion_line where tenant_id=$1`,
		`delete from inventory.bulk_uom_conversion where tenant_id=$1`,
		`delete from inventory.bulk_cost_application where tenant_id=$1`,
		`delete from inventory.bulk_cost_layer where tenant_id=$1`,
		`delete from inventory.bulk_inventory_entry where tenant_id=$1`,
		`delete from inventory.bulk_reservation where tenant_id=$1`,
		`delete from inventory.bulk_balance where tenant_id=$1`,
		`delete from inventory.sales_warehouse_binding where tenant_id=$1`,
		`delete from inventory.item_bin_policy where tenant_id=$1`,
		`delete from inventory.warehouse_bin where tenant_id=$1`,
		`delete from inventory.item_unit_of_measure where tenant_id=$1`,
		`delete from inventory.stock_item where tenant_id=$1`,
		`delete from sales.customer_order_line where tenant_id=$1`,
		`delete from sales.customer_order where tenant_id=$1`,
		`delete from catalog.vehicle_variant where tenant_id=$1`,
		`delete from catalog.vehicle_model where tenant_id=$1`,
		`delete from crm.customer_profile where tenant_id=$1`,
		`delete from platform.outbox_event where tenant_id=$1`,
		`delete from org.organization where tenant_id=$1`,
		`delete from platform.tenant where tenant_id=$1`,
		`alter table inventory.item_unit_of_measure enable trigger item_unit_of_measure_immutable`,
		`alter table inventory.bulk_inventory_entry enable trigger bulk_inventory_entry_immutable`,
		`alter table inventory.bulk_cost_application enable trigger bulk_cost_application_immutable`,
		`alter table inventory.warehouse_receipt enable trigger warehouse_receipt_immutable`,
		`alter table inventory.warehouse_receipt_line enable trigger warehouse_receipt_line_immutable`,
		`alter table inventory.warehouse_activity_line enable trigger warehouse_activity_line_immutable`,
		`alter table inventory.warehouse_activity_uom_reservation enable trigger warehouse_uom_reservation_delete_guard`,
		`alter table inventory.sales_warehouse_binding enable trigger sales_warehouse_binding_immutable`,
		`alter table inventory.warehouse_pick_request enable trigger warehouse_pick_request_immutable`,
		`alter table sales.customer_shipment enable trigger customer_shipment_immutable`,
		`alter table sales.customer_shipment_line enable trigger customer_shipment_line_immutable`,
		`alter table inventory.customer_shipment_allocation enable trigger customer_shipment_allocation_immutable`,
		`alter table logistics.shipment_provider_report enable trigger shipment_provider_report_immutable`,
	}
	for _, command := range commands {
		args := []any{}
		if strings.Contains(command, "$1") {
			args = append(args, tenant)
		}
		if _, err = tx.Exec(ctx, command, args...); err != nil {
			t.Errorf("connected transport cleanup: %v", err)
			return
		}
	}
	if err = tx.Commit(ctx); err != nil {
		t.Errorf("connected transport cleanup commit: %v", err)
	}
}
