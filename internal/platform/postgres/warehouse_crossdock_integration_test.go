package postgres

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/inventorycontrol"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestWarehouseCrossDockTransferDemandEndToEnd(t *testing.T) {
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
	if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Cross Dock V166','Cross Dock V166')`, tenant, "xd-"+tenant[:8]); err != nil {
		t.Fatal(err)
	}
	defer cleanupWarehouseCrossDock(t, pool, tenant)
	if _, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'warehouse','warehouse','Warehouse','warehouse')`, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'franchise','franchise','Franchise','franchisee')`, tenant); err != nil {
		t.Fatal(err)
	}

	repo := NewInventoryControl(pool)
	if _, err = repo.CreateBulkItem(ctx, tenant, bulkTestUUID(t), inventorycontrol.BulkItem{ID: "part", Code: "PART-XD", Description: "Cross-docked part", BaseUOM: "EA", BaseRoundingPrecision: "0.000001", TrackingMode: "lot", CostingMethod: "fifo", Version: 1}); err != nil {
		t.Fatal(err)
	}
	if _, err = repo.CreateWarehouseBin(ctx, tenant, bulkTestUUID(t), inventorycontrol.WarehouseBin{ID: "invalid-cross-dock", OrganizationID: "warehouse", Code: "INVALID-XD", Type: "receive", CrossDock: true, Version: 1}); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("cross-dock receive bin accepted: %v", err)
	}
	bins := []inventorycontrol.WarehouseBin{
		{ID: "receive", OrganizationID: "warehouse", Code: "RECEIVE", Type: "receive", Version: 1},
		{ID: "storage", OrganizationID: "warehouse", Code: "STORAGE", Type: "putpick", Ranking: 20, Version: 1},
		{ID: "cross-dock", OrganizationID: "warehouse", Code: "CROSS-DOCK", Type: "putpick", Ranking: 100, CrossDock: true, Version: 1},
		{ID: "ship", OrganizationID: "warehouse", Code: "SHIP", Type: "ship", Version: 1},
		{ID: "destination-receive", OrganizationID: "franchise", Code: "RECEIVE", Type: "receive", Version: 1},
	}
	for _, bin := range bins {
		if _, err = repo.CreateWarehouseBin(ctx, tenant, bulkTestUUID(t), bin); err != nil {
			t.Fatalf("create bin %s: %v", bin.ID, err)
		}
		maximum := "100"
		if bin.CrossDock {
			maximum = "5"
		}
		if _, err = repo.ConfigureItemBin(ctx, tenant, bulkTestUUID(t), inventorycontrol.ItemBinPolicy{OrganizationID: bin.OrganizationID, ItemID: "part", BinID: bin.ID, Fixed: true, Default: bin.ID == "storage", MinQuantity: "0", MaxQuantity: maximum, Version: 1}); err != nil {
			t.Fatalf("configure bin %s: %v", bin.ID, err)
		}
	}
	policy := inventorycontrol.WarehouseCrossDockPolicy{OrganizationID: "warehouse", ItemID: "part", BinID: "cross-dock", DueDateDays: 5, Enabled: true, Version: 1}
	invalidPolicy := policy
	invalidPolicy.BinID = "storage"
	if _, err = repo.ConfigureWarehouseCrossDock(ctx, tenant, bulkTestUUID(t), invalidPolicy); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("ordinary storage accepted as cross-dock policy: %v", err)
	}
	if _, err = repo.ConfigureWarehouseCrossDock(ctx, tenant, bulkTestUUID(t), policy); err != nil {
		t.Fatal(err)
	}

	eligible, err := repo.CreateBulkTransfer(ctx, tenant, "eligible-transfer", "eligible-line", bulkTestUUID(t), inventorycontrol.BulkTransferCommand{RequestID: "eligible-request", FromOrganizationID: "warehouse", ToOrganizationID: "franchise", ReceiveBinID: "destination-receive", InTransitCode: "ROAD", ItemID: "part", Quantity: "5", PostingDate: time.Date(2030, 1, 3, 0, 0, 0, 0, time.UTC), SourceKind: "replenishment", SourceID: "eligible-source"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.CreateBulkTransfer(ctx, tenant, "late-transfer", "late-line", bulkTestUUID(t), inventorycontrol.BulkTransferCommand{RequestID: "late-request", FromOrganizationID: "warehouse", ToOrganizationID: "franchise", ReceiveBinID: "destination-receive", InTransitCode: "ROAD", ItemID: "part", Quantity: "10", PostingDate: time.Date(2030, 1, 20, 0, 0, 0, 0, time.UTC), SourceKind: "replenishment", SourceID: "late-source"}); err != nil {
		t.Fatal(err)
	}
	expires := time.Date(2035, 1, 1, 0, 0, 0, 0, time.UTC)
	if _, err = pool.Exec(ctx, `insert into inventory.inventory_lot(tenant_id,lot_id,item_id,lot_no,expiration_date) values($1,'initial-lot','part','LOT-INITIAL',$2::date)`, tenant, expires); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `insert into inventory.bulk_balance(tenant_id,organization_id,bin_id,item_id,lot_id,quantity,reserved_quantity) values($1,'warehouse','storage','part','initial-lot',2,0)`, tenant); err != nil {
		t.Fatal(err)
	}
	genericPick, err := repo.CreateWarehousePick(ctx, tenant, warehouseTestIDs(t, "generic-transfer-pick"), inventorycontrol.WarehousePickCommand{RequestID: "generic-transfer-pick-request", OrganizationID: "warehouse", ShipBinID: "ship", ItemID: "part", DemandKind: "transfer-outbound", DemandID: eligible.ID, DemandLineID: eligible.LineID, Quantity: "2", UseFEFO: true})
	if err != nil || len(genericPick.Lines) != 1 || genericPick.Lines[0].FromBinID != "storage" {
		t.Fatalf("generic transfer pick=%+v err=%v", genericPick, err)
	}
	first, err := repo.PostWarehouseReceipt(ctx, tenant, warehouseTestIDs(t, "cross-first"), inventorycontrol.WarehouseReceiptCommand{RequestID: "cross-first-request", OrganizationID: "warehouse", ReceiveBinID: "receive", ItemID: "part", LotNo: "LOT-XD-1", ExpirationDate: &expires, Quantity: "7", UnitCost: "10", PostingDate: time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC), SourceKind: "purchase", SourceID: "po-cross-1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(first.PutAway.Lines) != 2 || first.PutAway.Lines[0].ToBinID != "cross-dock" || first.PutAway.Lines[0].Quantity != "3" || first.PutAway.Lines[1].ToBinID != "storage" || first.PutAway.Lines[1].Quantity != "4" {
		t.Fatalf("partial cross-dock placement=%+v", first.PutAway.Lines)
	}
	second, err := repo.PostWarehouseReceipt(ctx, tenant, warehouseTestIDs(t, "cross-second"), inventorycontrol.WarehouseReceiptCommand{RequestID: "cross-second-request", OrganizationID: "warehouse", ReceiveBinID: "receive", ItemID: "part", LotNo: "LOT-XD-2", ExpirationDate: &expires, Quantity: "3", UnitCost: "11", PostingDate: time.Date(2030, 1, 2, 0, 0, 0, 0, time.UTC), SourceKind: "purchase", SourceID: "po-cross-2"})
	if err != nil {
		t.Fatal(err)
	}
	if len(second.PutAway.Lines) != 1 || second.PutAway.Lines[0].ToBinID != "storage" || second.PutAway.Lines[0].Quantity != "3" {
		t.Fatalf("existing opportunity was duplicated: %+v", second.PutAway.Lines)
	}
	var allocations int
	var allocated string
	if err = pool.QueryRow(ctx, `select count(*),sum(quantity)::text from inventory.warehouse_crossdock_allocation where tenant_id=$1`, tenant).Scan(&allocations, &allocated); err != nil || allocations != 1 || allocated != "3.000000" {
		t.Fatalf("allocations=%d quantity=%s err=%v", allocations, allocated, err)
	}
	if _, err = repo.CreateWarehousePick(ctx, tenant, warehouseTestIDs(t, "unavailable-pick"), inventorycontrol.WarehousePickCommand{RequestID: "unavailable-pick-request", OrganizationID: "warehouse", ShipBinID: "ship", ItemID: "part", DemandKind: "transfer-outbound", DemandID: eligible.ID, DemandLineID: eligible.LineID, Quantity: "1", UseFEFO: true}); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("unregistered cross-dock allocation was pickable: %v", err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "warehouse", first.PutAway.ID, 1, time.Date(2030, 1, 2, 0, 0, 0, 0, time.UTC), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	if _, err = repo.CreateWarehousePick(ctx, tenant, warehouseTestIDs(t, "cross-over-demand"), inventorycontrol.WarehousePickCommand{RequestID: "cross-over-demand-request", OrganizationID: "warehouse", ShipBinID: "ship", ItemID: "part", DemandKind: "transfer-outbound", DemandID: eligible.ID, DemandLineID: eligible.LineID, Quantity: "4", UseFEFO: true}); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("over-demand transfer pick accepted: %v", err)
	}

	cancelPick, err := repo.CreateWarehousePick(ctx, tenant, warehouseTestIDs(t, "cross-cancel"), inventorycontrol.WarehousePickCommand{RequestID: "cross-cancel-request", OrganizationID: "warehouse", ShipBinID: "ship", ItemID: "part", DemandKind: "transfer-outbound", DemandID: eligible.ID, DemandLineID: eligible.LineID, Quantity: "1", UseFEFO: true})
	if err != nil || len(cancelPick.Lines) != 1 || cancelPick.Lines[0].FromBinID != "cross-dock" {
		t.Fatalf("cross-dock cancellation pick=%+v err=%v", cancelPick, err)
	}
	if _, err = repo.CancelWarehousePick(ctx, tenant, "warehouse", cancelPick.ID, 1, bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	var reserved string
	var releasedLinks int
	if err = pool.QueryRow(ctx, `select (select reserved_quantity::text from inventory.warehouse_crossdock_allocation where tenant_id=$1),(select count(*) from inventory.warehouse_crossdock_pick_link where tenant_id=$1 and status='released')`, tenant).Scan(&reserved, &releasedLinks); err != nil || reserved != "0.000000" || releasedLinks != 1 {
		t.Fatalf("cancel reserved=%s released_links=%d err=%v", reserved, releasedLinks, err)
	}

	pick, err := repo.CreateWarehousePick(ctx, tenant, warehouseTestIDs(t, "cross-pick"), inventorycontrol.WarehousePickCommand{RequestID: "cross-pick-request", OrganizationID: "warehouse", ShipBinID: "ship", ItemID: "part", DemandKind: "transfer-outbound", DemandID: eligible.ID, DemandLineID: eligible.LineID, Quantity: "3", UseFEFO: true})
	if err != nil || len(pick.Lines) != 1 || pick.Lines[0].FromBinID != "cross-dock" || pick.Lines[0].Quantity != "3" {
		t.Fatalf("cross-dock pick=%+v err=%v", pick, err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "warehouse", pick.ID, 1, time.Date(2030, 1, 3, 0, 0, 0, 0, time.UTC), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	var picked, crossQuantity, storageQuantity, shipQuantity string
	var pickedLinks, lateAllocations int
	if err = pool.QueryRow(ctx, `select x.reserved_quantity::text,x.picked_quantity::text,(select sum(quantity)::text from inventory.bulk_balance where tenant_id=$1 and organization_id='warehouse' and bin_id='cross-dock'),(select sum(quantity)::text from inventory.bulk_balance where tenant_id=$1 and organization_id='warehouse' and bin_id='storage'),(select sum(quantity)::text from inventory.bulk_balance where tenant_id=$1 and organization_id='warehouse' and bin_id='ship'),(select count(*) from inventory.warehouse_crossdock_pick_link where tenant_id=$1 and status='picked'),(select count(*) from inventory.warehouse_crossdock_allocation where tenant_id=$1 and transfer_id='late-transfer') from inventory.warehouse_crossdock_allocation x where x.tenant_id=$1`, tenant).Scan(&reserved, &picked, &crossQuantity, &storageQuantity, &shipQuantity, &pickedLinks, &lateAllocations); err != nil {
		t.Fatal(err)
	}
	if reserved != "0.000000" || picked != "3.000000" || crossQuantity != "0.000000" || storageQuantity != "6.000000" || shipQuantity != "3.000000" || pickedLinks != 1 || lateAllocations != 0 {
		t.Fatalf("final reserved=%s picked=%s cross=%s storage=%s ship=%s picked_links=%d late=%d", reserved, picked, crossQuantity, storageQuantity, shipQuantity, pickedLinks, lateAllocations)
	}
}

func cleanupWarehouseCrossDock(t *testing.T, pool *pgxpool.Pool, tenant string) {
	t.Helper()
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Errorf("cross-dock cleanup begin: %v", err)
		return
	}
	defer tx.Rollback(ctx)
	commands := []string{
		`alter table inventory.warehouse_pick_request disable trigger warehouse_pick_request_immutable`,
		`alter table inventory.warehouse_activity_uom_reservation disable trigger warehouse_uom_reservation_delete_guard`,
		`alter table inventory.warehouse_activity_line disable trigger warehouse_activity_line_immutable`,
		`alter table inventory.warehouse_receipt_line disable trigger warehouse_receipt_line_immutable`,
		`alter table inventory.warehouse_receipt disable trigger warehouse_receipt_immutable`,
		`alter table inventory.bulk_inventory_entry disable trigger bulk_inventory_entry_immutable`,
		`alter table inventory.item_unit_of_measure disable trigger item_unit_of_measure_immutable`,
		`delete from inventory.warehouse_activity_uom_reservation where tenant_id=$1`,
		`delete from inventory.warehouse_crossdock_pick_link where tenant_id=$1`,
		`delete from inventory.warehouse_crossdock_allocation where tenant_id=$1`,
		`delete from inventory.warehouse_crossdock_policy where tenant_id=$1`,
		`delete from inventory.warehouse_pick_request where tenant_id=$1`,
		`delete from inventory.warehouse_activity_line where tenant_id=$1`,
		`delete from inventory.warehouse_activity where tenant_id=$1`,
		`delete from inventory.warehouse_receipt_line where tenant_id=$1`,
		`delete from inventory.warehouse_receipt where tenant_id=$1`,
		`delete from inventory.bulk_transfer_line where tenant_id=$1`,
		`delete from inventory.bulk_transfer where tenant_id=$1`,
		`delete from inventory.bulk_cost_layer where tenant_id=$1`,
		`delete from inventory.bulk_inventory_entry where tenant_id=$1`,
		`delete from inventory.bulk_reservation where tenant_id=$1`,
		`delete from inventory.bulk_balance where tenant_id=$1`,
		`delete from inventory.inventory_lot where tenant_id=$1`,
		`delete from inventory.item_bin_policy where tenant_id=$1`,
		`delete from inventory.warehouse_bin where tenant_id=$1`,
		`delete from inventory.item_unit_of_measure where tenant_id=$1`,
		`delete from inventory.stock_item where tenant_id=$1`,
		`delete from platform.outbox_event where tenant_id=$1`,
		`delete from org.organization where tenant_id=$1`,
		`delete from platform.tenant where tenant_id=$1`,
		`alter table inventory.warehouse_activity_line enable trigger warehouse_activity_line_immutable`,
		`alter table inventory.warehouse_receipt_line enable trigger warehouse_receipt_line_immutable`,
		`alter table inventory.warehouse_receipt enable trigger warehouse_receipt_immutable`,
		`alter table inventory.bulk_inventory_entry enable trigger bulk_inventory_entry_immutable`,
		`alter table inventory.item_unit_of_measure enable trigger item_unit_of_measure_immutable`,
		`alter table inventory.warehouse_activity_uom_reservation enable trigger warehouse_uom_reservation_delete_guard`,
		`alter table inventory.warehouse_pick_request enable trigger warehouse_pick_request_immutable`,
	}
	for _, command := range commands {
		args := []any{}
		if strings.Contains(command, "$1") {
			args = append(args, tenant)
		}
		if _, err = tx.Exec(ctx, command, args...); err != nil {
			t.Errorf("cross-dock cleanup failed: %v", err)
			return
		}
	}
	if err = tx.Commit(ctx); err != nil {
		t.Errorf("cross-dock cleanup commit: %v", err)
	}
}
