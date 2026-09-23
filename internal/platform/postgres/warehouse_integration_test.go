package postgres

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/inventorycontrol"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func warehouseTestIDs(t *testing.T, prefix string) inventorycontrol.WarehouseIDs {
	t.Helper()
	return inventorycontrol.WarehouseIDs{ActivityID: prefix + "-activity", ReceiptID: prefix + "-receipt", LineID: prefix + "-line", EntryID: prefix + "-entry", LayerID: prefix + "-layer", LotID: prefix + "-lot", EventID: bulkTestUUID(t)}
}

func TestWarehouseReceiptPutAwayFEFOPickAndCancellation(t *testing.T) {
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
	if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Warehouse V161','Warehouse V161')`, tenant, "wh-"+tenant[:8]); err != nil {
		t.Fatal(err)
	}
	defer func() {
		tx, cleanupErr := pool.Begin(ctx)
		if cleanupErr != nil {
			t.Errorf("warehouse cleanup begin: %v", cleanupErr)
			return
		}
		defer tx.Rollback(ctx)
		commands := []struct {
			sql  string
			args []any
		}{
			{`alter table inventory.warehouse_pick_request disable trigger warehouse_pick_request_immutable`, nil},
			{`alter table inventory.warehouse_activity_uom_reservation disable trigger warehouse_uom_reservation_delete_guard`, nil},
			{`alter table inventory.warehouse_activity_line disable trigger warehouse_activity_line_immutable`, nil},
			{`alter table inventory.warehouse_receipt_line disable trigger warehouse_receipt_line_immutable`, nil},
			{`alter table inventory.warehouse_receipt disable trigger warehouse_receipt_immutable`, nil},
			{`alter table inventory.bulk_cost_application disable trigger bulk_cost_application_immutable`, nil},
			{`alter table inventory.bulk_inventory_entry disable trigger bulk_inventory_entry_immutable`, nil},
			{`alter table inventory.item_unit_of_measure disable trigger item_unit_of_measure_immutable`, nil},
			{`delete from inventory.warehouse_activity_uom_reservation where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.warehouse_pick_request where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.warehouse_activity_line where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.warehouse_activity where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.warehouse_receipt_line where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.warehouse_receipt where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.bulk_cost_application where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.bulk_cost_layer where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.bulk_inventory_entry where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.bulk_reservation where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.bulk_balance where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.inventory_lot where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.item_bin_policy where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.warehouse_bin where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.item_unit_of_measure where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.stock_item where tenant_id=$1`, []any{tenant}},
			{`delete from platform.outbox_event where tenant_id=$1`, []any{tenant}},
			{`delete from org.organization where tenant_id=$1`, []any{tenant}},
			{`delete from platform.tenant where tenant_id=$1`, []any{tenant}},
			{`alter table inventory.warehouse_activity_line enable trigger warehouse_activity_line_immutable`, nil},
			{`alter table inventory.warehouse_receipt_line enable trigger warehouse_receipt_line_immutable`, nil},
			{`alter table inventory.warehouse_receipt enable trigger warehouse_receipt_immutable`, nil},
			{`alter table inventory.bulk_cost_application enable trigger bulk_cost_application_immutable`, nil},
			{`alter table inventory.bulk_inventory_entry enable trigger bulk_inventory_entry_immutable`, nil},
			{`alter table inventory.item_unit_of_measure enable trigger item_unit_of_measure_immutable`, nil},
			{`alter table inventory.warehouse_activity_uom_reservation enable trigger warehouse_uom_reservation_delete_guard`, nil},
			{`alter table inventory.warehouse_pick_request enable trigger warehouse_pick_request_immutable`, nil},
		}
		for _, command := range commands {
			if _, cleanupErr = tx.Exec(ctx, command.sql, command.args...); cleanupErr != nil {
				t.Errorf("warehouse cleanup failed: %v", cleanupErr)
				return
			}
		}
		if cleanupErr = tx.Commit(ctx); cleanupErr != nil {
			t.Errorf("warehouse cleanup commit: %v", cleanupErr)
		}
	}()
	if _, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'warehouse','warehouse','Warehouse','warehouse')`, tenant); err != nil {
		t.Fatal(err)
	}
	repo := NewInventoryControl(pool)
	if _, err = repo.CreateBulkItem(ctx, tenant, bulkTestUUID(t), inventorycontrol.BulkItem{ID: "part", Code: "PART-FEFO", Description: "Tracked part", BaseUOM: "EA", BaseRoundingPrecision: "0.000001", TrackingMode: "lot", CostingMethod: "fifo", Version: 1}); err != nil {
		t.Fatal(err)
	}
	bins := []inventorycontrol.WarehouseBin{
		{ID: "receive", OrganizationID: "warehouse", Code: "RECEIVE", Type: "receive", Ranking: 0, Version: 1},
		{ID: "bulk", OrganizationID: "warehouse", Code: "BULK", Type: "putpick", Ranking: 10, Version: 1},
		{ID: "forward", OrganizationID: "warehouse", Code: "FORWARD", Type: "putpick", Ranking: 100, Version: 1},
		{ID: "ship", OrganizationID: "warehouse", Code: "SHIP", Type: "ship", Ranking: 0, Version: 1},
		{ID: "qc", OrganizationID: "warehouse", Code: "QC", Type: "qc", Ranking: 1000, Version: 1},
	}
	for _, bin := range bins {
		if _, err = repo.CreateWarehouseBin(ctx, tenant, bulkTestUUID(t), bin); err != nil {
			t.Fatal(err)
		}
		maximum := "100"
		if bin.ID == "forward" {
			maximum = "6"
		}
		if _, err = repo.ConfigureItemBin(ctx, tenant, bulkTestUUID(t), inventorycontrol.ItemBinPolicy{OrganizationID: "warehouse", ItemID: "part", BinID: bin.ID, Fixed: true, Default: bin.ID == "forward", MinQuantity: "0", MaxQuantity: maximum, Version: 1}); err != nil {
			t.Fatal(err)
		}
	}
	var configuredTracking, configuredBinType string
	if err = pool.QueryRow(ctx, `select i.tracking_mode,w.bin_type from inventory.stock_item i join inventory.item_bin_policy p on p.tenant_id=i.tenant_id and p.item_id=i.item_id join inventory.warehouse_bin w on w.tenant_id=p.tenant_id and w.organization_id=p.organization_id and w.bin_id=p.bin_id where i.tenant_id=$1 and i.item_id='part' and p.organization_id='warehouse' and p.bin_id='receive'`, tenant).Scan(&configuredTracking, &configuredBinType); err != nil || configuredTracking != "lot" || configuredBinType != "receive" {
		t.Fatalf("configured receipt boundary tracking=%s bin=%s err=%v", configuredTracking, configuredBinType, err)
	}
	late := time.Date(2035, 12, 31, 0, 0, 0, 0, time.UTC)
	early := time.Date(2034, 1, 1, 0, 0, 0, 0, time.UTC)
	receipts := []struct {
		prefix, lot, quantity, cost string
		expiry                      *time.Time
	}{
		{"late", "LOT-LATE", "6", "100", &late},
		{"early", "LOT-EARLY", "4", "120", &early},
	}
	activities := make([]inventorycontrol.WarehouseActivity, 0, len(receipts))
	for index, receipt := range receipts {
		result, receiveErr := repo.PostWarehouseReceipt(ctx, tenant, warehouseTestIDs(t, receipt.prefix), inventorycontrol.WarehouseReceiptCommand{RequestID: receipt.prefix + "-request", OrganizationID: "warehouse", ReceiveBinID: "receive", ItemID: "part", LotNo: receipt.lot, ExpirationDate: receipt.expiry, Quantity: receipt.quantity, UnitCost: receipt.cost, PostingDate: time.Date(2030, 1, index+1, 0, 0, 0, 0, time.UTC), SourceKind: "purchase", SourceID: receipt.prefix + "-po"})
		if receiveErr != nil {
			t.Fatalf("receipt %s: %v", receipt.prefix, receiveErr)
		}
		if len(result.PutAway.Lines) != 1 {
			t.Fatalf("put-away lines=%+v", result.PutAway.Lines)
		}
		activities = append(activities, result.PutAway)
	}
	if activities[0].Lines[0].ToBinID != "forward" || activities[1].Lines[0].ToBinID != "bulk" {
		t.Fatalf("ranking/capacity placements=%s,%s", activities[0].Lines[0].ToBinID, activities[1].Lines[0].ToBinID)
	}
	for index, activity := range activities {
		registered, registerErr := repo.RegisterWarehouseActivity(ctx, tenant, "warehouse", activity.ID, 1, time.Date(2030, 1, index+3, 0, 0, 0, 0, time.UTC), bulkTestUUID(t))
		if registerErr != nil || registered.Status != "registered" || registered.Version != 2 {
			t.Fatalf("register=%+v err=%v", registered, registerErr)
		}
	}
	var pickableQuantity string
	if err = pool.QueryRow(ctx, `select coalesce(sum(b.quantity-b.reserved_quantity),0)::text from inventory.bulk_balance b join inventory.item_bin_policy p using(tenant_id,organization_id,item_id,bin_id) join inventory.warehouse_bin w using(tenant_id,organization_id,bin_id) left join inventory.inventory_lot l on l.tenant_id=b.tenant_id and l.lot_id=b.lot_id where b.tenant_id=$1 and b.organization_id='warehouse' and b.item_id='part' and w.bin_type in ('pick','putpick') and not w.movement_blocked and not p.dedicated and coalesce(l.blocked,false)=false and (l.expiration_date is null or l.expiration_date>=current_date)`, tenant).Scan(&pickableQuantity); err != nil || pickableQuantity != "10.000000" {
		t.Fatalf("pickable quantity=%s err=%v", pickableQuantity, err)
	}
	pick, err := repo.CreateWarehousePick(ctx, tenant, warehouseTestIDs(t, "pick"), inventorycontrol.WarehousePickCommand{RequestID: "pick-request", OrganizationID: "warehouse", ShipBinID: "ship", ItemID: "part", DemandKind: "manual", DemandID: "order-1", DemandLineID: "line-1", Quantity: "7", UseFEFO: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(pick.Lines) != 2 || pick.Lines[0].LotNo != "LOT-EARLY" || pick.Lines[0].Quantity != "4" || pick.Lines[1].LotNo != "LOT-LATE" || pick.Lines[1].Quantity != "3" {
		t.Fatalf("FEFO lines=%+v", pick.Lines)
	}
	registeredPick, err := repo.RegisterWarehouseActivity(ctx, tenant, "warehouse", pick.ID, 1, time.Date(2030, 1, 5, 0, 0, 0, 0, time.UTC), bulkTestUUID(t))
	if err != nil || registeredPick.Status != "registered" || registeredPick.Lines[0].ReservationVersion != 2 {
		t.Fatalf("registered pick=%+v err=%v", registeredPick, err)
	}
	var shipQuantity, shipReserved string
	if err = pool.QueryRow(ctx, `select sum(quantity)::text,sum(reserved_quantity)::text from inventory.bulk_balance where tenant_id=$1 and organization_id='warehouse' and bin_id='ship'`, tenant).Scan(&shipQuantity, &shipReserved); err != nil || shipQuantity != "7.000000" || shipReserved != "7.000000" {
		t.Fatalf("ship quantity=%s reserved=%s err=%v", shipQuantity, shipReserved, err)
	}
	concurrentIDs := []inventorycontrol.WarehouseIDs{warehouseTestIDs(t, "race-a"), warehouseTestIDs(t, "race-b")}
	concurrentPicks := make([]inventorycontrol.WarehouseActivity, 2)
	concurrentErrors := make([]error, 2)
	var group sync.WaitGroup
	for index := range concurrentIDs {
		group.Add(1)
		go func(i int) {
			defer group.Done()
			concurrentPicks[i], concurrentErrors[i] = repo.CreateWarehousePick(ctx, tenant, concurrentIDs[i], inventorycontrol.WarehousePickCommand{RequestID: fmt.Sprintf("race-request-%d", i), OrganizationID: "warehouse", ShipBinID: "ship", ItemID: "part", DemandKind: "service", DemandID: fmt.Sprintf("service-%d", i), DemandLineID: "line-1", Quantity: "3", UseFEFO: true})
		}(index)
	}
	group.Wait()
	winner, successes := inventorycontrol.WarehouseActivity{}, 0
	for index, concurrentErr := range concurrentErrors {
		if concurrentErr == nil {
			winner, successes = concurrentPicks[index], successes+1
		} else if !errors.Is(concurrentErr, inventorycontrol.ErrConflict) {
			t.Fatalf("unexpected concurrent pick error: %v", concurrentErr)
		}
	}
	if successes != 1 {
		t.Fatalf("concurrent pick successes=%d errors=%v", successes, concurrentErrors)
	}
	cancelled, err := repo.CancelWarehousePick(ctx, tenant, "warehouse", winner.ID, 1, bulkTestUUID(t))
	if err != nil || cancelled.Status != "cancelled" || cancelled.Version != 2 {
		t.Fatalf("cancelled=%+v err=%v", cancelled, err)
	}
	var remainingReserved string
	if err = pool.QueryRow(ctx, `select reserved_quantity::text from inventory.bulk_balance where tenant_id=$1 and bin_id='forward' and lot_id='late-lot'`, tenant).Scan(&remainingReserved); err != nil || remainingReserved != "0.000000" {
		t.Fatalf("cancel released=%s err=%v", remainingReserved, err)
	}
	if _, err = pool.Exec(ctx, `update inventory.warehouse_receipt set source_id='tampered' where tenant_id=$1 and receipt_id='late-receipt'`, tenant); err == nil {
		t.Fatal("immutable receipt accepted update")
	} else {
		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) || pgErr.Code != "55000" {
			t.Fatalf("unexpected immutability error: %v", err)
		}
	}
}
