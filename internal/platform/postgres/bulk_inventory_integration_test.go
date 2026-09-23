package postgres

import (
	"context"
	"crypto/rand"
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

func bulkTestUUID(t *testing.T) string {
	t.Helper()
	var value [16]byte
	if _, err := rand.Read(value[:]); err != nil {
		t.Fatal(err)
	}
	value[6] = (value[6] & 0x0f) | 0x40
	value[8] = (value[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", value[0:4], value[4:6], value[6:8], value[8:10], value[10:16])
}

func TestBulkLotBinFIFOConcurrencyAndImmutableApplications(t *testing.T) {
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
	if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Bulk Inventory','Bulk Inventory')`, tenant, "bulk-"+tenant[:8]); err != nil {
		t.Fatal(err)
	}
	defer func() {
		tx, cleanupErr := pool.Begin(ctx)
		if cleanupErr != nil {
			t.Errorf("bulk fixture cleanup begin failed: %v", cleanupErr)
			return
		}
		defer tx.Rollback(ctx)
		commands := []struct {
			sql  string
			args []any
		}{
			{`alter table inventory.bulk_cost_application disable trigger bulk_cost_application_immutable`, nil},
			{`alter table inventory.bulk_inventory_entry disable trigger bulk_inventory_entry_immutable`, nil},
			{`alter table inventory.item_unit_of_measure disable trigger item_unit_of_measure_immutable`, nil},
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
			{`alter table inventory.bulk_cost_application enable trigger bulk_cost_application_immutable`, nil},
			{`alter table inventory.bulk_inventory_entry enable trigger bulk_inventory_entry_immutable`, nil},
			{`alter table inventory.item_unit_of_measure enable trigger item_unit_of_measure_immutable`, nil},
		}
		for _, command := range commands {
			if _, cleanupErr = tx.Exec(ctx, command.sql, command.args...); cleanupErr != nil {
				t.Errorf("bulk fixture cleanup failed: %v", cleanupErr)
				return
			}
		}
		if cleanupErr = tx.Commit(ctx); cleanupErr != nil {
			t.Errorf("bulk fixture cleanup commit failed: %v", cleanupErr)
		}
	}()
	if _, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'warehouse','warehouse','Warehouse','warehouse')`, tenant); err != nil {
		t.Fatal(err)
	}
	repo := NewInventoryControl(pool)
	item := inventorycontrol.BulkItem{ID: "part", Code: "PART-1", Description: "Brake pad", BaseUOM: "EA", BaseRoundingPrecision: "0.000001", TrackingMode: "lot", CostingMethod: "fifo", Version: 1}
	if _, err = repo.CreateBulkItem(ctx, tenant, bulkTestUUID(t), item); err != nil {
		t.Fatal(err)
	}
	for _, bin := range []inventorycontrol.WarehouseBin{
		{ID: "pick", OrganizationID: "warehouse", Code: "PICK-01", Type: "pick", Version: 1},
		{ID: "ship", OrganizationID: "warehouse", Code: "SHIP-01", Type: "ship", Version: 1},
	} {
		if _, err = repo.CreateWarehouseBin(ctx, tenant, bulkTestUUID(t), bin); err != nil {
			t.Fatal(err)
		}
		policy := inventorycontrol.ItemBinPolicy{OrganizationID: "warehouse", ItemID: "part", BinID: bin.ID, Fixed: true, Default: bin.ID == "pick", MinQuantity: "0", MaxQuantity: "100", Version: 1}
		if _, err = repo.ConfigureItemBin(ctx, tenant, bulkTestUUID(t), policy); err != nil {
			t.Fatal(err)
		}
	}
	expires := time.Date(2035, 1, 1, 0, 0, 0, 0, time.UTC)
	receipts := []struct {
		entry, layer, lot, quantity, cost string
		date                              time.Time
	}{
		{"receipt-1", "layer-1", "lot-1", "5", "100", time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"receipt-2", "layer-2", "ignored-existing-lot", "5", "120", time.Date(2030, 1, 2, 0, 0, 0, 0, time.UTC)},
	}
	for index, receipt := range receipts {
		value := inventorycontrol.BulkReceipt{OrganizationID: "warehouse", BinID: "pick", ItemID: "part", LotNo: "LOT-2035", ExpirationDate: &expires, Quantity: receipt.quantity, UnitCost: receipt.cost, PostingDate: receipt.date, SourceKind: "purchase-receipt", SourceID: fmt.Sprintf("purchase-%d", index+1)}
		got, receiveErr := repo.ReceiveBulk(ctx, tenant, bulkTestUUID(t), receipt.entry, receipt.layer, receipt.lot, value)
		if receiveErr != nil || got.LotID != "lot-1" {
			t.Fatalf("receipt %d result=%+v err=%v", index, got, receiveErr)
		}
	}
	availability, err := repo.BulkAvailability(ctx, tenant, "warehouse", "part")
	if err != nil || len(availability) != 1 || availability[0].AvailableQuantity != "10.000000" {
		t.Fatalf("availability=%+v err=%v", availability, err)
	}
	reservations := []inventorycontrol.BulkReservation{
		{ID: "reserve-a", OrganizationID: "warehouse", BinID: "pick", ItemID: "part", LotID: "lot-1", DemandKind: "service", DemandID: "service-a", DemandLineID: "line", Quantity: "7", Status: "reservation", Version: 1},
		{ID: "reserve-b", OrganizationID: "warehouse", BinID: "pick", ItemID: "part", LotID: "lot-1", DemandKind: "service", DemandID: "service-b", DemandLineID: "line", Quantity: "7", Status: "reservation", Version: 1},
	}
	errs := make([]error, 2)
	var wg sync.WaitGroup
	for index := range reservations {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, errs[i] = repo.ReserveBulk(ctx, tenant, bulkTestUUID(t), reservations[i])
		}(index)
	}
	wg.Wait()
	winner, successes := "", 0
	for index, reserveErr := range errs {
		if reserveErr == nil {
			winner, successes = reservations[index].ID, successes+1
		} else if !errors.Is(reserveErr, inventorycontrol.ErrConflict) {
			t.Fatalf("unexpected reserve error: %v", reserveErr)
		}
	}
	if successes != 1 {
		t.Fatalf("reservation successes=%d errors=%v", successes, errs)
	}
	issue, err := repo.IssueBulk(ctx, tenant, bulkTestUUID(t), "issue-1", inventorycontrol.BulkIssue{OrganizationID: "warehouse", ReservationID: winner, ReservationVersion: 1, PostingDate: time.Date(2030, 1, 3, 0, 0, 0, 0, time.UTC), SourceKind: "service-use", SourceID: "service-use-1"})
	if err != nil || issue.Quantity != "7.000000" || issue.CostAmount != "740" || issue.Applications != 2 {
		t.Fatalf("issue=%+v err=%v", issue, err)
	}
	var applicationQuantity, applicationCost string
	if err = pool.QueryRow(ctx, `select sum(quantity)::text,sum(cost_amount)::text from inventory.bulk_cost_application where tenant_id=$1 and outbound_entry_id='issue-1'`, tenant).Scan(&applicationQuantity, &applicationCost); err != nil || applicationQuantity != "7.000000" || applicationCost != "740.0000" {
		t.Fatalf("applications quantity=%s cost=%s err=%v", applicationQuantity, applicationCost, err)
	}
	if _, err = pool.Exec(ctx, `update inventory.bulk_inventory_entry set cost_amount=0 where tenant_id=$1 and entry_id='issue-1'`, tenant); err == nil {
		t.Fatal("immutable inventory entry accepted update")
	} else {
		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) || pgErr.Code != "55000" {
			t.Fatalf("unexpected immutable error: %v", err)
		}
	}
	if err = repo.MoveBulk(ctx, tenant, bulkTestUUID(t), bulkTestUUID(t), inventorycontrol.BulkMovement{OrganizationID: "warehouse", FromBinID: "pick", ToBinID: "ship", ItemID: "part", LotID: "lot-1", Quantity: "2", PostingDate: time.Date(2030, 1, 4, 0, 0, 0, 0, time.UTC), SourceID: "movement-1"}); err != nil {
		t.Fatal(err)
	}
	var pickQuantity, shipQuantity string
	if err = pool.QueryRow(ctx, `select quantity::text from inventory.bulk_balance where tenant_id=$1 and bin_id='pick'`, tenant).Scan(&pickQuantity); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select quantity::text from inventory.bulk_balance where tenant_id=$1 and bin_id='ship'`, tenant).Scan(&shipQuantity); err != nil {
		t.Fatal(err)
	}
	if pickQuantity != "1.000000" || shipQuantity != "2.000000" {
		t.Fatalf("movement pick=%s ship=%s", pickQuantity, shipQuantity)
	}

	specific := inventorycontrol.BulkItem{ID: "specific-part", Code: "SPECIFIC-1", Description: "Specific assembly", BaseUOM: "EA", BaseRoundingPrecision: "0.000001", TrackingMode: "none", CostingMethod: "specific", Version: 1}
	if _, err = repo.CreateBulkItem(ctx, tenant, bulkTestUUID(t), specific); err != nil {
		t.Fatal(err)
	}
	if _, err = repo.ConfigureItemBin(ctx, tenant, bulkTestUUID(t), inventorycontrol.ItemBinPolicy{OrganizationID: "warehouse", ItemID: "specific-part", BinID: "pick", Fixed: true, MinQuantity: "0", MaxQuantity: "10", Version: 1}); err != nil {
		t.Fatal(err)
	}
	if _, err = repo.ReceiveBulk(ctx, tenant, bulkTestUUID(t), "specific-receipt", "specific-layer", "unused-lot", inventorycontrol.BulkReceipt{OrganizationID: "warehouse", BinID: "pick", ItemID: "specific-part", Quantity: "1", UnitCost: "999.99", PostingDate: time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC), SourceKind: "purchase-receipt", SourceID: "specific-purchase"}); err != nil {
		t.Fatal(err)
	}
	specificReservation := inventorycontrol.BulkReservation{ID: "specific-reservation", OrganizationID: "warehouse", BinID: "pick", ItemID: "specific-part", DemandKind: "service", DemandID: "specific-service", DemandLineID: "line", Quantity: "1", Status: "reservation", Version: 1}
	if _, err = repo.ReserveBulk(ctx, tenant, bulkTestUUID(t), specificReservation); err != nil {
		t.Fatal(err)
	}
	withoutSpecific := inventorycontrol.BulkIssue{OrganizationID: "warehouse", ReservationID: specificReservation.ID, ReservationVersion: 1, PostingDate: time.Date(2030, 1, 2, 0, 0, 0, 0, time.UTC), SourceKind: "service-use", SourceID: "specific-use"}
	if _, err = repo.IssueBulk(ctx, tenant, bulkTestUUID(t), "specific-issue-rejected", withoutSpecific); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("specific issue without source=%v", err)
	}
	withoutSpecific.SpecificReceiptEntry = "specific-receipt"
	specificIssue, err := repo.IssueBulk(ctx, tenant, bulkTestUUID(t), "specific-issue", withoutSpecific)
	if err != nil || specificIssue.CostAmount != "999.99" || specificIssue.Applications != 1 {
		t.Fatalf("specific issue=%+v err=%v", specificIssue, err)
	}
}
