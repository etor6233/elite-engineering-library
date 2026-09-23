package postgres

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/inventorycontrol"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestWarehouseHandlingUnitConservationIdempotencyAndConcurrency(t *testing.T) {
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
	if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Packaging V168','Packaging V168')`, tenant, "pack-"+tenant[:8]); err != nil {
		t.Fatal(err)
	}
	defer func() {
		tx, cleanupErr := pool.Begin(ctx)
		if cleanupErr != nil {
			t.Errorf("packaging cleanup begin: %v", cleanupErr)
			return
		}
		defer tx.Rollback(ctx)
		commands := []struct {
			sql  string
			args []any
		}{
			{`alter table inventory.bulk_uom_conversion_line disable trigger bulk_uom_conversion_line_immutable`, nil},
			{`alter table inventory.bulk_uom_conversion disable trigger bulk_uom_conversion_immutable`, nil},
			{`alter table inventory.bulk_inventory_entry disable trigger bulk_inventory_entry_immutable`, nil},
			{`alter table inventory.item_unit_of_measure disable trigger item_unit_of_measure_immutable`, nil},
			{`delete from inventory.bulk_uom_conversion_line where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.bulk_uom_conversion where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.bulk_cost_application where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.bulk_cost_layer where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.bulk_inventory_entry where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.bulk_balance where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.item_bin_policy where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.warehouse_bin where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.item_unit_of_measure where tenant_id=$1`, []any{tenant}},
			{`delete from inventory.stock_item where tenant_id=$1`, []any{tenant}},
			{`delete from platform.outbox_event where tenant_id=$1`, []any{tenant}},
			{`delete from org.organization where tenant_id=$1`, []any{tenant}},
			{`delete from platform.tenant where tenant_id=$1`, []any{tenant}},
			{`alter table inventory.bulk_uom_conversion_line enable trigger bulk_uom_conversion_line_immutable`, nil},
			{`alter table inventory.bulk_uom_conversion enable trigger bulk_uom_conversion_immutable`, nil},
			{`alter table inventory.bulk_inventory_entry enable trigger bulk_inventory_entry_immutable`, nil},
			{`alter table inventory.item_unit_of_measure enable trigger item_unit_of_measure_immutable`, nil},
		}
		for _, command := range commands {
			if _, cleanupErr = tx.Exec(ctx, command.sql, command.args...); cleanupErr != nil {
				t.Errorf("packaging cleanup: %v", cleanupErr)
				return
			}
		}
		if cleanupErr = tx.Commit(ctx); cleanupErr != nil {
			t.Errorf("packaging cleanup commit: %v", cleanupErr)
		}
	}()
	if _, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'warehouse','warehouse','Warehouse','warehouse')`, tenant); err != nil {
		t.Fatal(err)
	}
	repo := NewInventoryControl(pool)
	if _, err = repo.CreateBulkItem(ctx, tenant, bulkTestUUID(t), inventorycontrol.BulkItem{ID: "part", Code: "PACK-PART", Description: "Packaged part", BaseUOM: "EA", BaseRoundingPrecision: "1", TrackingMode: "none", CostingMethod: "fifo", Version: 1}); err != nil {
		t.Fatal(err)
	}
	for _, unit := range []inventorycontrol.ItemUnitOfMeasure{{ItemID: "part", Code: "BOX", QuantityPerUnit: "12", RoundingPrecision: "1", Version: 1}, {ItemID: "part", Code: "PALLET", QuantityPerUnit: "120", RoundingPrecision: "1", Version: 1}} {
		if _, err = repo.ConfigureItemUnitOfMeasure(ctx, tenant, bulkTestUUID(t), unit); err != nil {
			t.Fatal(err)
		}
	}
	for _, bin := range []inventorycontrol.WarehouseBin{{ID: "pick", OrganizationID: "warehouse", Code: "PICK", Type: "pick", Version: 1}, {ID: "target", OrganizationID: "warehouse", Code: "TARGET", Type: "putpick", Version: 1}} {
		if _, err = repo.CreateWarehouseBin(ctx, tenant, bulkTestUUID(t), bin); err != nil {
			t.Fatal(err)
		}
		if _, err = repo.ConfigureItemBin(ctx, tenant, bulkTestUUID(t), inventorycontrol.ItemBinPolicy{OrganizationID: "warehouse", ItemID: "part", BinID: bin.ID, Fixed: true, Default: bin.ID == "pick", MinQuantity: "0", MaxQuantity: "1000", Version: 1}); err != nil {
			t.Fatal(err)
		}
	}
	if _, err = repo.ReceiveBulk(ctx, tenant, bulkTestUUID(t), bulkTestUUID(t), bulkTestUUID(t), bulkTestUUID(t), inventorycontrol.BulkReceipt{OrganizationID: "warehouse", BinID: "pick", ItemID: "part", Quantity: "24", UnitCost: "1", PostingDate: time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC), SourceKind: "purchase-receipt", SourceID: "receipt-1"}); err != nil {
		t.Fatal(err)
	}
	ids := func() inventorycontrol.PackagingConversionIDs {
		return inventorycontrol.PackagingConversionIDs{ConversionID: bulkTestUUID(t), TakeLineID: bulkTestUUID(t), PlaceLineID: bulkTestUUID(t), EventID: bulkTestUUID(t)}
	}
	gather := inventorycontrol.PackagingConversionCommand{RequestID: "gather-1", OrganizationID: "warehouse", BinID: "pick", ItemID: "part", Operation: "gather", FromUOM: "EA", ToUOM: "BOX", FromQuantity: "24"}
	first, err := repo.ConvertBulkHandlingUnits(ctx, tenant, ids(), gather)
	if err != nil || first.ToQuantity != "2.000000" || first.BaseQuantity != "24.000000" || len(first.Lines) != 2 || first.Lines[0].Action != "take" || first.Lines[1].Action != "place" {
		t.Fatalf("gather=%+v err=%v", first, err)
	}
	replay, err := repo.ConvertBulkHandlingUnits(ctx, tenant, ids(), gather)
	if err != nil || replay.ID != first.ID {
		t.Fatalf("replay=%+v err=%v", replay, err)
	}
	divergent := gather
	divergent.FromQuantity = "12"
	if _, err = repo.ConvertBulkHandlingUnits(ctx, tenant, ids(), divergent); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("divergent replay error=%v", err)
	}
	breakOne := inventorycontrol.PackagingConversionCommand{RequestID: "break-1", OrganizationID: "warehouse", BinID: "pick", ItemID: "part", Operation: "breakbulk", FromUOM: "BOX", ToUOM: "EA", FromQuantity: "1"}
	if _, err = repo.ConvertBulkHandlingUnits(ctx, tenant, ids(), breakOne); err != nil {
		t.Fatal(err)
	}

	commands := []inventorycontrol.PackagingConversionCommand{breakOne, breakOne}
	commands[0].RequestID = "race-1"
	commands[1].RequestID = "race-2"
	errs := make([]error, 2)
	var wait sync.WaitGroup
	for index := range commands {
		wait.Add(1)
		go func(i int) {
			defer wait.Done()
			_, errs[i] = repo.ConvertBulkHandlingUnits(ctx, tenant, ids(), commands[i])
		}(index)
	}
	wait.Wait()
	successes, conflicts := 0, 0
	for _, raceErr := range errs {
		if raceErr == nil {
			successes++
		} else if errors.Is(raceErr, inventorycontrol.ErrConflict) {
			conflicts++
		} else {
			t.Fatalf("unexpected race error=%v", raceErr)
		}
	}
	if successes != 1 || conflicts != 1 {
		t.Fatalf("race successes=%d conflicts=%d errors=%v", successes, conflicts, errs)
	}
	var physical, composed string
	var conversions, lines, outbox int
	if err = pool.QueryRow(ctx, `select b.quantity::text,coalesce(sum(u.quantity_base),0)::text,(select count(*) from inventory.bulk_uom_conversion where tenant_id=$1),(select count(*) from inventory.bulk_uom_conversion_line where tenant_id=$1),(select count(*) from platform.outbox_event where tenant_id=$1 and event_type='bulk-uom-conversion.completed') from inventory.bulk_balance b left join inventory.bulk_uom_balance u using(balance_id) where b.tenant_id=$1 and b.organization_id='warehouse' and b.bin_id='pick' and b.item_id='part' group by b.balance_id`, tenant).Scan(&physical, &composed, &conversions, &lines, &outbox); err != nil || physical != "24.000000" || composed != "24.000000" || conversions != 3 || lines != 6 || outbox != 3 {
		t.Fatalf("physical=%s composed=%s conversions=%d lines=%d outbox=%d err=%v", physical, composed, conversions, lines, outbox, err)
	}
	if _, err = pool.Exec(ctx, `update inventory.bulk_uom_conversion set operation='gather' where tenant_id=$1 and conversion_id=$2`, tenant, first.ID); err == nil {
		t.Fatal("immutable conversion mutated")
	}
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Exec(ctx, `insert into inventory.bulk_uom_balance(balance_id,uom_code,quantity,qty_per_uom,quantity_base) select balance_id,'BOX',1,12,12 from inventory.bulk_balance where tenant_id=$1 and bin_id='pick' and item_id='part'`, tenant)
	if err != nil {
		tx.Rollback(ctx)
		t.Fatal(err)
	}
	err = tx.Commit(ctx)
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != "55000" {
		t.Fatalf("conservation commit error=%v", err)
	}

	if _, err = repo.ConvertBulkHandlingUnits(ctx, tenant, ids(), inventorycontrol.PackagingConversionCommand{RequestID: "gather-2", OrganizationID: "warehouse", BinID: "pick", ItemID: "part", Operation: "gather", FromUOM: "EA", ToUOM: "BOX", FromQuantity: "24"}); err != nil {
		t.Fatal(err)
	}
	err = repo.MoveBulk(ctx, tenant, bulkTestUUID(t), bulkTestUUID(t), inventorycontrol.BulkMovement{OrganizationID: "warehouse", FromBinID: "pick", ToBinID: "target", ItemID: "part", Quantity: "1", PostingDate: time.Date(2030, 1, 2, 0, 0, 0, 0, time.UTC), SourceID: "move-packaged"})
	if !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("implicit breakbulk error=%v", err)
	}

	if _, err = repo.CreateBulkItem(ctx, tenant, bulkTestUUID(t), inventorycontrol.BulkItem{ID: "weight", Code: "WEIGHT", Description: "Fractional material", BaseUOM: "KG", BaseRoundingPrecision: "0.001", TrackingMode: "none", CostingMethod: "fifo", Version: 1}); err != nil {
		t.Fatal(err)
	}
	for _, bin := range []string{"pick", "target"} {
		if _, err = repo.ConfigureItemBin(ctx, tenant, bulkTestUUID(t), inventorycontrol.ItemBinPolicy{OrganizationID: "warehouse", ItemID: "weight", BinID: bin, Fixed: true, Default: bin == "pick", MinQuantity: "0", MaxQuantity: "1000", Version: 1}); err != nil {
			t.Fatal(err)
		}
	}
	if _, err = repo.ReceiveBulk(ctx, tenant, bulkTestUUID(t), bulkTestUUID(t), bulkTestUUID(t), bulkTestUUID(t), inventorycontrol.BulkReceipt{OrganizationID: "warehouse", BinID: "pick", ItemID: "weight", Quantity: "0.125", UnitCost: "1", PostingDate: time.Date(2030, 1, 3, 0, 0, 0, 0, time.UTC), SourceKind: "purchase-receipt", SourceID: "weight-1"}); err != nil {
		t.Fatal(err)
	}
	if err = repo.MoveBulk(ctx, tenant, bulkTestUUID(t), bulkTestUUID(t), inventorycontrol.BulkMovement{OrganizationID: "warehouse", FromBinID: "pick", ToBinID: "target", ItemID: "weight", Quantity: "0.025", PostingDate: time.Date(2030, 1, 4, 0, 0, 0, 0, time.UTC), SourceID: "weight-move"}); err != nil {
		t.Fatalf("fractional base move: %v", err)
	}
}
