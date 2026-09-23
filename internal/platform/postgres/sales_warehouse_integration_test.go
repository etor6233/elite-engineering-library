package postgres

import (
	"context"
	"errors"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/inventorycontrol"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestSalesOrderWarehouseDemandIsBoundOutstandingIdempotentAndConcurrent(t *testing.T) {
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
	if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Connected Sales Warehouse','Connected Sales Warehouse')`, tenant, "swh-"+tenant[:8]); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'warehouse','warehouse','Warehouse','warehouse')`, tenant); err != nil {
		t.Fatal(err)
	}
	defer cleanupSalesWarehouse(t, pool, tenant)
	for _, command := range []string{
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state) values($1,'model','MODEL','Model','other','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,homologation_state,lifecycle_state) values($1,'variant','model','VARIANT','Variant','{}','approved','active')`,
		`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version) values($1,'order','warehouse','customer','draft','USD',300,1)`,
		`insert into sales.customer_order_line(tenant_id,order_id,line_id,variant_id,quantity,unit_price_minor_units) values($1,'order','line','variant',3,100)`,
	} {
		if _, err = pool.Exec(ctx, command, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewInventoryControl(pool)
	if _, err = repo.CreateBulkItem(ctx, tenant, bulkTestUUID(t), inventorycontrol.BulkItem{ID: "part", Code: "SALE-PART", Description: "Sale part", BaseUOM: "EA", BaseRoundingPrecision: "1", TrackingMode: "none", CostingMethod: "fifo", Version: 1}); err != nil {
		t.Fatal(err)
	}
	if _, err = repo.ConfigureItemUnitOfMeasure(ctx, tenant, bulkTestUUID(t), inventorycontrol.ItemUnitOfMeasure{ItemID: "part", Code: "BOX", QuantityPerUnit: "12", RoundingPrecision: "1", Version: 1}); err != nil {
		t.Fatal(err)
	}
	for _, bin := range []inventorycontrol.WarehouseBin{{ID: "receive", OrganizationID: "warehouse", Code: "RECEIVE", Type: "receive", Version: 1}, {ID: "pick", OrganizationID: "warehouse", Code: "PICK", Type: "putpick", Ranking: 100, Version: 1}, {ID: "ship", OrganizationID: "warehouse", Code: "SHIP", Type: "ship", Version: 1}} {
		if _, err = repo.CreateWarehouseBin(ctx, tenant, bulkTestUUID(t), bin); err != nil {
			t.Fatal(err)
		}
	}
	for _, policy := range []inventorycontrol.ItemBinPolicy{{OrganizationID: "warehouse", ItemID: "part", BinID: "receive", Fixed: true, MinQuantity: "0", MaxQuantity: "36", Version: 1}, {OrganizationID: "warehouse", ItemID: "part", BinID: "pick", Fixed: true, Default: true, MinQuantity: "0", MaxQuantity: "36", Version: 1}, {OrganizationID: "warehouse", ItemID: "part", BinID: "ship", Fixed: true, MinQuantity: "0", MaxQuantity: "36", Version: 1}} {
		if _, err = repo.ConfigureItemBin(ctx, tenant, bulkTestUUID(t), policy); err != nil {
			t.Fatal(err)
		}
	}
	posting := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	receipt, err := repo.PostWarehouseReceipt(ctx, tenant, warehousePackagingIDs(t, "sales-source"), inventorycontrol.WarehouseReceiptCommand{RequestID: "sales-source-receipt", OrganizationID: "warehouse", ReceiveBinID: "receive", ItemID: "part", HandlingUOM: "BOX", HandlingQuantity: "3", UnitCost: "10", PostingDate: posting, SourceKind: "purchase", SourceID: "sales-source-po"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "warehouse", receipt.PutAway.ID, 1, posting.AddDate(0, 0, 1), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	binding := inventorycontrol.SalesWarehouseBinding{RequestID: "binding-request", VariantID: "variant", ItemID: "part", SalesUOMCode: "BOX", Version: 1}
	configured, err := repo.ConfigureSalesWarehouseBinding(ctx, tenant, bulkTestUUID(t), binding)
	if err != nil || configured != binding {
		t.Fatalf("binding=%+v err=%v", configured, err)
	}
	replayedBinding, err := repo.ConfigureSalesWarehouseBinding(ctx, tenant, bulkTestUUID(t), binding)
	if err != nil || replayedBinding != binding {
		t.Fatalf("binding replay=%+v err=%v", replayedBinding, err)
	}
	divergentBinding := binding
	divergentBinding.SalesUOMCode = "EA"
	if _, err = repo.ConfigureSalesWarehouseBinding(ctx, tenant, bulkTestUUID(t), divergentBinding); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("divergent binding error=%v", err)
	}
	draftPick := inventorycontrol.WarehousePickCommand{RequestID: "draft-pick", OrganizationID: "warehouse", ShipBinID: "ship", ItemID: "part", DemandKind: "customer-order", DemandID: "order", DemandLineID: "line", Quantity: "12", HandlingUOM: "BOX"}
	if _, err = repo.CreateWarehousePick(ctx, tenant, warehousePackagingIDs(t, "draft-pick"), draftPick); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("draft order pick error=%v", err)
	}
	if _, err = pool.Exec(ctx, `update sales.customer_order set state='placed',version=2,updated_at=clock_timestamp() where tenant_id=$1 and order_id='order'`, tenant); err != nil {
		t.Fatal(err)
	}
	cancelCommand := draftPick
	cancelCommand.RequestID = "cancel-pick"
	cancellable, err := repo.CreateWarehousePick(ctx, tenant, warehousePackagingIDs(t, "cancel-pick"), cancelCommand)
	if err != nil || len(cancellable.Lines) != 1 || cancellable.Lines[0].FromUOMCode != "BOX" || cancellable.Lines[0].UOMCode != "BOX" {
		t.Fatalf("cancellable=%+v err=%v", cancellable, err)
	}
	replayed, err := repo.CreateWarehousePick(ctx, tenant, warehousePackagingIDs(t, "ignored-replay"), cancelCommand)
	if err != nil || replayed.ID != cancellable.ID || replayed.Status != "open" {
		t.Fatalf("pick replay=%+v err=%v", replayed, err)
	}
	divergentPick := cancelCommand
	divergentPick.Quantity = "6"
	if _, err = repo.CreateWarehousePick(ctx, tenant, warehousePackagingIDs(t, "divergent-replay"), divergentPick); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("divergent pick replay error=%v", err)
	}
	cancelled, err := repo.CancelWarehousePick(ctx, tenant, "warehouse", cancellable.ID, 1, bulkTestUUID(t))
	if err != nil || cancelled.Status != "cancelled" {
		t.Fatalf("cancelled=%+v err=%v", cancelled, err)
	}
	registeredCommand := cancelCommand
	registeredCommand.RequestID = "registered-pick"
	registeredPick, err := repo.CreateWarehousePick(ctx, tenant, warehousePackagingIDs(t, "registered-pick"), registeredCommand)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "warehouse", registeredPick.ID, 1, posting.AddDate(0, 0, 2), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	firstShipmentCommand := inventorycontrol.CustomerShipmentCommand{RequestID: "shipment-first-request", OrganizationID: "warehouse", OrderID: "order", WarehouseActivityID: registeredPick.ID, PostingDate: posting.AddDate(0, 0, 3)}
	firstShipment, err := repo.PostCustomerShipment(ctx, tenant, inventorycontrol.CustomerShipmentIDs{ShipmentID: "shipment-first", LineID: "shipment-first-line", EventID: bulkTestUUID(t)}, firstShipmentCommand)
	if err != nil || firstShipment.ID != "shipment-first" || firstShipment.FulfillmentState != "partially-shipped" || firstShipment.OrderVersion != 3 || firstShipment.Line.SalesUOMCode != "BOX" || !exactDecimalEqual(firstShipment.Line.Quantity, "1") || !exactDecimalEqual(firstShipment.Line.QuantityBase, "12") || !exactDecimalEqual(firstShipment.Line.CostAmount, "120") || firstShipment.Line.AllocationCount != 1 {
		t.Fatalf("first shipment=%+v err=%v", firstShipment, err)
	}
	firstReplay, err := repo.PostCustomerShipment(ctx, tenant, inventorycontrol.CustomerShipmentIDs{ShipmentID: "ignored-replay", LineID: "ignored-replay-line", EventID: bulkTestUUID(t)}, firstShipmentCommand)
	if err != nil || firstReplay.ID != firstShipment.ID || firstReplay.Line.ID != firstShipment.Line.ID || firstReplay.OrderVersion != firstShipment.OrderVersion {
		t.Fatalf("first shipment replay=%+v err=%v", firstReplay, err)
	}
	divergentShipment := firstShipmentCommand
	divergentShipment.PostingDate = divergentShipment.PostingDate.AddDate(0, 0, 1)
	if _, err = repo.PostCustomerShipment(ctx, tenant, inventorycontrol.CustomerShipmentIDs{ShipmentID: "ignored-divergent", LineID: "ignored-divergent-line", EventID: bulkTestUUID(t)}, divergentShipment); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("divergent shipment replay error=%v", err)
	}
	commands := []inventorycontrol.WarehousePickCommand{registeredCommand, registeredCommand}
	commands[0].RequestID, commands[0].Quantity = "concurrent-a", "24"
	commands[1].RequestID, commands[1].Quantity = "concurrent-b", "24"
	results := make([]inventorycontrol.WarehouseActivity, 2)
	errs := make([]error, 2)
	var wg sync.WaitGroup
	for index := range commands {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			results[i], errs[i] = repo.CreateWarehousePick(ctx, tenant, warehousePackagingIDs(t, commands[i].RequestID), commands[i])
		}(index)
	}
	wg.Wait()
	winners := 0
	for index := range errs {
		if errs[index] == nil {
			winners++
			if len(results[index].Lines) != 1 || results[index].Lines[0].UOMCode != "BOX" || !exactDecimalEqual(results[index].Lines[0].UOMQuantity, "2") {
				t.Fatalf("concurrent winner=%+v", results[index])
			}
		} else if !errors.Is(errs[index], inventorycontrol.ErrConflict) {
			t.Fatalf("concurrent error[%d]=%v", index, errs[index])
		}
	}
	if winners != 1 {
		t.Fatalf("concurrent winners=%d errors=%v", winners, errs)
	}
	winningPick := inventorycontrol.WarehouseActivity{}
	for index := range errs {
		if errs[index] == nil {
			winningPick = results[index]
		}
	}
	registeredRemainder, err := repo.RegisterWarehouseActivity(ctx, tenant, "warehouse", winningPick.ID, 1, posting.AddDate(0, 0, 4), bulkTestUUID(t))
	if err != nil || registeredRemainder.Status != "registered" {
		t.Fatalf("registered remainder=%+v err=%v", registeredRemainder, err)
	}
	shipmentCommands := []inventorycontrol.CustomerShipmentCommand{
		{RequestID: "shipment-race-a", OrganizationID: "warehouse", OrderID: "order", WarehouseActivityID: winningPick.ID, PostingDate: posting.AddDate(0, 0, 5)},
		{RequestID: "shipment-race-b", OrganizationID: "warehouse", OrderID: "order", WarehouseActivityID: winningPick.ID, PostingDate: posting.AddDate(0, 0, 5)},
	}
	shipmentResults := make([]inventorycontrol.CustomerShipment, 2)
	shipmentErrors := make([]error, 2)
	for index := range shipmentCommands {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			shipmentResults[i], shipmentErrors[i] = repo.PostCustomerShipment(ctx, tenant, inventorycontrol.CustomerShipmentIDs{ShipmentID: "shipment-race-" + string(rune('a'+i)), LineID: "shipment-race-line-" + string(rune('a'+i)), EventID: bulkTestUUID(t)}, shipmentCommands[i])
		}(index)
	}
	wg.Wait()
	shipmentWinners := 0
	winningShipmentIndex := -1
	for index := range shipmentErrors {
		if shipmentErrors[index] == nil {
			shipmentWinners++
			winningShipmentIndex = index
			if shipmentResults[index].FulfillmentState != "shipped" || shipmentResults[index].OrderVersion != 4 || !exactDecimalEqual(shipmentResults[index].Line.Quantity, "2") || !exactDecimalEqual(shipmentResults[index].Line.QuantityBase, "24") || !exactDecimalEqual(shipmentResults[index].Line.CostAmount, "240") {
				t.Fatalf("shipment winner=%+v", shipmentResults[index])
			}
		} else if !errors.Is(shipmentErrors[index], inventorycontrol.ErrConflict) {
			t.Fatalf("shipment error[%d]=%v", index, shipmentErrors[index])
		}
	}
	if shipmentWinners != 1 {
		t.Fatalf("shipment winners=%d errors=%v", shipmentWinners, shipmentErrors)
	}
	winningShipmentReplay, err := repo.PostCustomerShipment(ctx, tenant, inventorycontrol.CustomerShipmentIDs{ShipmentID: "ignored-winning-shipment", LineID: "ignored-winning-line", EventID: bulkTestUUID(t)}, shipmentCommands[winningShipmentIndex])
	if err != nil || winningShipmentReplay.ID != shipmentResults[winningShipmentIndex].ID || winningShipmentReplay.FulfillmentState != "shipped" {
		t.Fatalf("winning shipment replay=%+v err=%v", winningShipmentReplay, err)
	}
	var handled string
	if err = pool.QueryRow(ctx, `select coalesce(sum(al.quantity),0)::text from inventory.warehouse_activity a join inventory.warehouse_activity_line al using(tenant_id,activity_id) where a.tenant_id=$1 and a.source_kind='customer-order' and a.source_id='order' and a.source_line_id='line' and a.status in ('open','registered')`, tenant).Scan(&handled); err != nil || !exactDecimalEqual(handled, "36") {
		t.Fatalf("handled=%s err=%v", handled, err)
	}
	blocked := registeredCommand
	blocked.RequestID = "beyond-outstanding"
	if _, err = repo.CreateWarehousePick(ctx, tenant, warehousePackagingIDs(t, "beyond-outstanding"), blocked); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("beyond outstanding error=%v", err)
	}
	var orderFulfillment, shippedQuantity, physicalQuantity, composedQuantity, remainingCost, shipmentCost string
	var orderVersionAfter, shipmentCount, consumedReservations int64
	if err = pool.QueryRow(ctx, `select o.fulfillment_state,o.version,l.shipped_quantity::text from sales.customer_order o join sales.customer_order_line l using(tenant_id,order_id) where o.tenant_id=$1 and o.order_id='order'`, tenant).Scan(&orderFulfillment, &orderVersionAfter, &shippedQuantity); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select coalesce(sum(b.quantity),0)::text,coalesce(sum(c.quantity_base),0)::text from inventory.bulk_balance b left join inventory.bulk_uom_balance c on c.balance_id=b.balance_id where b.tenant_id=$1 and b.item_id='part'`, tenant).Scan(&physicalQuantity, &composedQuantity); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select coalesce(sum(remaining_quantity),0)::text from inventory.bulk_cost_layer where tenant_id=$1 and item_id='part'`, tenant).Scan(&remainingCost); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*),coalesce(sum(total_cost_amount),0)::text from sales.customer_shipment where tenant_id=$1 and order_id='order'`, tenant).Scan(&shipmentCount, &shipmentCost); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from inventory.bulk_reservation where tenant_id=$1 and demand_kind='customer-order' and status='consumed'`, tenant).Scan(&consumedReservations); err != nil {
		t.Fatal(err)
	}
	if orderFulfillment != "shipped" || orderVersionAfter != 4 || !exactDecimalEqual(shippedQuantity, "3") || !exactDecimalEqual(physicalQuantity, "0") || !exactDecimalEqual(composedQuantity, "0") || !exactDecimalEqual(remainingCost, "0") || shipmentCount != 2 || !exactDecimalEqual(shipmentCost, "360") || consumedReservations != 2 {
		t.Fatalf("fulfillment=%s version=%d shipped=%s physical=%s composed=%s cost-remaining=%s shipments=%d shipment-cost=%s consumed=%d", orderFulfillment, orderVersionAfter, shippedQuantity, physicalQuantity, composedQuantity, remainingCost, shipmentCount, shipmentCost, consumedReservations)
	}
}

func cleanupSalesWarehouse(t *testing.T, pool *pgxpool.Pool, tenant string) {
	t.Helper()
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Errorf("sales warehouse cleanup begin: %v", err)
		return
	}
	defer tx.Rollback(ctx)
	commands := []string{
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
	}
	for _, command := range commands {
		arguments := []any{}
		if strings.Contains(command, "$1") {
			arguments = append(arguments, tenant)
		}
		if _, err = tx.Exec(ctx, command, arguments...); err != nil {
			t.Errorf("sales warehouse cleanup: %v", err)
			return
		}
	}
	if err = tx.Commit(ctx); err != nil {
		t.Errorf("sales warehouse cleanup commit: %v", err)
	}
}
