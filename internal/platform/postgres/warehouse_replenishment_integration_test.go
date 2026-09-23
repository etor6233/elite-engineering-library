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
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestWarehouseReplenishmentFEFORegisterAndCancel(t *testing.T) {
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
	if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Replenishment V165','Replenishment V165')`, tenant, "rp-"+tenant[:8]); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'warehouse','warehouse','Warehouse','warehouse')`, tenant); err != nil {
		t.Fatal(err)
	}
	defer cleanupWarehouseReplenishment(t, pool, tenant)

	repo := NewInventoryControl(pool)
	if _, err = repo.CreateBulkItem(ctx, tenant, bulkTestUUID(t), inventorycontrol.BulkItem{ID: "part", Code: "PART-REPLENISH", Description: "Replenished part", BaseUOM: "EA", BaseRoundingPrecision: "0.000001", TrackingMode: "lot", CostingMethod: "fifo", Version: 1}); err != nil {
		t.Fatal(err)
	}
	bins := []inventorycontrol.WarehouseBin{
		{ID: "source-early", OrganizationID: "warehouse", Code: "SOURCE-EARLY", Type: "putpick", Ranking: 30, Version: 1},
		{ID: "source-late", OrganizationID: "warehouse", Code: "SOURCE-LATE", Type: "putpick", Ranking: 40, Version: 1},
		{ID: "source-blocked", OrganizationID: "warehouse", Code: "SOURCE-BLOCKED", Type: "putpick", Ranking: 50, MovementBlocked: true, Version: 1},
		{ID: "source-receive", OrganizationID: "warehouse", Code: "SOURCE-RECEIVE", Type: "receive", Ranking: 60, Version: 1},
		{ID: "target", OrganizationID: "warehouse", Code: "TARGET", Type: "pick", Ranking: 100, Version: 1},
		{ID: "target-cancel", OrganizationID: "warehouse", Code: "TARGET-CANCEL", Type: "pick", Ranking: 110, Version: 1},
		{ID: "target-race", OrganizationID: "warehouse", Code: "TARGET-RACE", Type: "pick", Ranking: 120, Version: 1},
	}
	for _, bin := range bins {
		if _, err = repo.CreateWarehouseBin(ctx, tenant, bulkTestUUID(t), bin); err != nil {
			t.Fatalf("create bin %s: %v", bin.ID, err)
		}
		minimum, maximum := "0", "100"
		if bin.ID == "target" || bin.ID == "target-cancel" || bin.ID == "target-race" {
			minimum, maximum = "5", "10"
		}
		if _, err = repo.ConfigureItemBin(ctx, tenant, bulkTestUUID(t), inventorycontrol.ItemBinPolicy{OrganizationID: "warehouse", ItemID: "part", BinID: bin.ID, Fixed: true, MinQuantity: minimum, MaxQuantity: maximum, Version: 1}); err != nil {
			t.Fatalf("configure bin %s: %v", bin.ID, err)
		}
	}

	if _, err = pool.Exec(ctx, `
insert into inventory.inventory_lot(tenant_id,lot_id,item_id,lot_no,expiration_date,blocked) values
 ($1,'lot-target','part','LOT-TARGET','2036-01-01',false),
 ($1,'lot-early','part','LOT-EARLY','2034-01-01',false),
 ($1,'lot-late','part','LOT-LATE','2035-01-01',false),
		 ($1,'lot-blocked','part','LOT-BLOCKED','2033-01-01',true),
		 ($1,'lot-receive','part','LOT-RECEIVE','2032-01-01',false)`, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `insert into inventory.bulk_balance(tenant_id,organization_id,bin_id,item_id,lot_id,quantity,reserved_quantity) values
		 ($1,'warehouse','target','part','lot-target',2,0),
		 ($1,'warehouse','source-early','part','lot-early',5,0),
		 ($1,'warehouse','source-late','part','lot-late',6,0),
 ($1,'warehouse','source-blocked','part','lot-blocked',20,0),
 ($1,'warehouse','source-receive','part','lot-receive',20,0)`, tenant); err != nil {
		t.Fatal(err)
	}

	ids := warehouseTestIDs(t, "replenish")
	command := inventorycontrol.WarehouseReplenishmentCommand{RequestID: "replenish-request", OrganizationID: "warehouse", ToBinID: "target", ItemID: "part", UseFEFO: true, AssignedTo: "worker-1"}
	activity, err := repo.CreateWarehouseReplenishment(ctx, tenant, ids, command)
	if err != nil {
		t.Fatal(err)
	}
	if activity.Type != "movement" || activity.Status != "open" || len(activity.Lines) != 2 || activity.Lines[0].LotNo != "LOT-EARLY" || activity.Lines[0].Quantity != "5.000000" || activity.Lines[1].LotNo != "LOT-LATE" || activity.Lines[1].Quantity != "3.000000" {
		t.Fatalf("unexpected FEFO activity: %+v", activity)
	}
	replay, err := repo.CreateWarehouseReplenishment(ctx, tenant, warehouseTestIDs(t, "replay"), command)
	if err != nil || replay.ID != activity.ID || len(replay.Lines) != 2 {
		t.Fatalf("idempotent replay=%+v err=%v", replay, err)
	}
	divergent := command
	divergent.ToBinID = "target-cancel"
	if _, err = repo.CreateWarehouseReplenishment(ctx, tenant, warehouseTestIDs(t, "divergent"), divergent); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("divergent replay error=%v", err)
	}

	registered, err := repo.RegisterWarehouseActivity(ctx, tenant, "warehouse", activity.ID, 1, time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC), bulkTestUUID(t))
	if err != nil || registered.Status != "registered" || registered.Version != 2 {
		t.Fatalf("register=%+v err=%v", registered, err)
	}
	var targetQuantity, sourceEarly, sourceLate string
	if err = pool.QueryRow(ctx, `select (select sum(quantity)::text from inventory.bulk_balance where tenant_id=$1 and bin_id='target'),(select quantity::text from inventory.bulk_balance where tenant_id=$1 and bin_id='source-early'),(select quantity::text from inventory.bulk_balance where tenant_id=$1 and bin_id='source-late')`, tenant).Scan(&targetQuantity, &sourceEarly, &sourceLate); err != nil || targetQuantity != "10.000000" || sourceEarly != "0.000000" || sourceLate != "3.000000" {
		t.Fatalf("balances target=%s early=%s late=%s err=%v", targetQuantity, sourceEarly, sourceLate, err)
	}
	var consumed int
	if err = pool.QueryRow(ctx, `select count(*) from inventory.bulk_reservation where tenant_id=$1 and demand_kind='warehouse-replenishment' and status='consumed'`, tenant).Scan(&consumed); err != nil || consumed != 2 {
		t.Fatalf("consumed reservations=%d err=%v", consumed, err)
	}
	var movementEntries int
	if err = pool.QueryRow(ctx, `select count(*) from inventory.bulk_inventory_entry where tenant_id=$1 and source_kind='warehouse-movement'`, tenant).Scan(&movementEntries); err != nil || movementEntries != 4 {
		t.Fatalf("movement entries=%d err=%v", movementEntries, err)
	}
	if _, err = repo.CreateWarehouseReplenishment(ctx, tenant, warehouseTestIDs(t, "full"), inventorycontrol.WarehouseReplenishmentCommand{RequestID: "full-request", OrganizationID: "warehouse", ToBinID: "target", ItemID: "part", UseFEFO: true}); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("full target accepted: %v", err)
	}

	cancelActivity, err := repo.CreateWarehouseReplenishment(ctx, tenant, warehouseTestIDs(t, "cancel"), inventorycontrol.WarehouseReplenishmentCommand{RequestID: "cancel-request", OrganizationID: "warehouse", ToBinID: "target-cancel", ItemID: "part", UseFEFO: true})
	if err != nil || len(cancelActivity.Lines) != 3 || cancelActivity.Lines[0].Quantity != "5.000000" || cancelActivity.Lines[1].Quantity != "3.000000" || cancelActivity.Lines[2].Quantity != "2.000000" {
		t.Fatalf("cancel activity=%+v err=%v", cancelActivity, err)
	}
	cancelled, err := repo.CancelWarehouseReplenishment(ctx, tenant, "warehouse", cancelActivity.ID, 1, bulkTestUUID(t))
	if err != nil || cancelled.Status != "cancelled" || cancelled.Version != 2 {
		t.Fatalf("cancelled=%+v err=%v", cancelled, err)
	}
	var released int
	var reserved string
	if err = pool.QueryRow(ctx, `select (select count(*) from inventory.bulk_reservation where tenant_id=$1 and demand_id=$2 and status='released'),(select coalesce(sum(reserved_quantity),0)::text from inventory.bulk_balance where tenant_id=$1)`, tenant, cancelActivity.ID).Scan(&released, &reserved); err != nil || released != 3 || reserved != "0.000000" {
		t.Fatalf("cancel released=%d reserved=%s err=%v", released, reserved, err)
	}

	activities := make([]inventorycontrol.WarehouseActivity, 2)
	errorsByCall := make([]error, 2)
	raceIDs := []inventorycontrol.WarehouseIDs{warehouseTestIDs(t, "race-0"), warehouseTestIDs(t, "race-1")}
	var group sync.WaitGroup
	for index := range activities {
		group.Add(1)
		go func(i int) {
			defer group.Done()
			activities[i], errorsByCall[i] = repo.CreateWarehouseReplenishment(ctx, tenant, raceIDs[i], inventorycontrol.WarehouseReplenishmentCommand{RequestID: fmt.Sprintf("race-request-%d", i), OrganizationID: "warehouse", ToBinID: "target-race", ItemID: "part", UseFEFO: true})
		}(index)
	}
	group.Wait()
	winner, successes := inventorycontrol.WarehouseActivity{}, 0
	for index, callErr := range errorsByCall {
		if callErr == nil {
			winner, successes = activities[index], successes+1
		} else if !errors.Is(callErr, inventorycontrol.ErrConflict) {
			t.Fatalf("unexpected concurrent error: %v", callErr)
		}
	}
	if successes != 1 || len(winner.Lines) == 0 {
		t.Fatalf("concurrent successes=%d errors=%v winner=%+v", successes, errorsByCall, winner)
	}
	if _, err = repo.CancelWarehouseReplenishment(ctx, tenant, "warehouse", winner.ID, 1, bulkTestUUID(t)); err != nil {
		t.Fatalf("cancel concurrent winner: %v", err)
	}
}

func cleanupWarehouseReplenishment(t *testing.T, pool *pgxpool.Pool, tenant string) {
	t.Helper()
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Errorf("replenishment cleanup begin: %v", err)
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
		{`alter table inventory.bulk_inventory_entry disable trigger bulk_inventory_entry_immutable`, nil},
		{`alter table inventory.item_unit_of_measure disable trigger item_unit_of_measure_immutable`, nil},
		{`delete from inventory.warehouse_replenishment_request where tenant_id=$1`, []any{tenant}},
		{`delete from inventory.warehouse_activity_uom_reservation where tenant_id=$1`, []any{tenant}},
		{`delete from inventory.warehouse_pick_request where tenant_id=$1`, []any{tenant}},
		{`delete from inventory.warehouse_activity_line where tenant_id=$1`, []any{tenant}},
		{`delete from inventory.warehouse_activity where tenant_id=$1`, []any{tenant}},
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
		{`alter table inventory.bulk_inventory_entry enable trigger bulk_inventory_entry_immutable`, nil},
		{`alter table inventory.item_unit_of_measure enable trigger item_unit_of_measure_immutable`, nil},
		{`alter table inventory.warehouse_activity_uom_reservation enable trigger warehouse_uom_reservation_delete_guard`, nil},
		{`alter table inventory.warehouse_pick_request enable trigger warehouse_pick_request_immutable`, nil},
	}
	for _, command := range commands {
		if _, err = tx.Exec(ctx, command.sql, command.args...); err != nil {
			t.Errorf("replenishment cleanup failed: %v", err)
			return
		}
	}
	if err = tx.Commit(ctx); err != nil {
		t.Errorf("replenishment cleanup commit: %v", err)
	}
}
