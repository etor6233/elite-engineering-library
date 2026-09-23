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

func warehousePackagingIDs(t *testing.T, prefix string) inventorycontrol.WarehouseIDs {
	t.Helper()
	return inventorycontrol.WarehouseIDs{
		ActivityID:        prefix + "-activity",
		ReceiptID:         prefix + "-receipt",
		LineID:            prefix + "-line",
		EntryID:           prefix + "-entry",
		LayerID:           prefix + "-layer",
		LotID:             prefix + "-lot",
		EventID:           bulkTestUUID(t),
		ConversionID:      prefix + "-conversion",
		TakeLineID:        prefix + "-take",
		PlaceLineID:       prefix + "-place",
		ConversionEventID: bulkTestUUID(t),
	}
}

func exactDecimalEqual(left, right string) bool {
	leftValue, leftOK := parseRat(left)
	rightValue, rightOK := parseRat(right)
	return leftOK && rightOK && leftValue.Cmp(rightValue) == 0
}

func TestWarehouseReceiptPackagingPreservationBreakbulkAndCancellation(t *testing.T) {
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
	if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Warehouse Packaging V169','Warehouse Packaging V169')`, tenant, "whp-"+tenant[:8]); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'warehouse','warehouse','Warehouse','warehouse')`, tenant); err != nil {
		t.Fatal(err)
	}
	defer func() {
		tx, cleanupErr := pool.Begin(ctx)
		if cleanupErr != nil {
			t.Errorf("warehouse packaging cleanup begin: %v", cleanupErr)
			return
		}
		defer tx.Rollback(ctx)
		commands := []string{
			`alter table inventory.warehouse_pick_request disable trigger warehouse_pick_request_immutable`,
			`alter table inventory.warehouse_activity_uom_reservation disable trigger warehouse_uom_reservation_delete_guard`,
			`alter table inventory.bulk_uom_conversion_line disable trigger bulk_uom_conversion_line_immutable`,
			`alter table inventory.bulk_uom_conversion disable trigger bulk_uom_conversion_immutable`,
			`alter table inventory.warehouse_activity_line disable trigger warehouse_activity_line_immutable`,
			`alter table inventory.warehouse_receipt_line disable trigger warehouse_receipt_line_immutable`,
			`alter table inventory.warehouse_receipt disable trigger warehouse_receipt_immutable`,
			`alter table inventory.bulk_inventory_entry disable trigger bulk_inventory_entry_immutable`,
			`alter table inventory.item_unit_of_measure disable trigger item_unit_of_measure_immutable`,
			`delete from inventory.warehouse_replenishment_request where tenant_id=$1`,
			`delete from inventory.warehouse_activity_uom_reservation where tenant_id=$1`,
			`delete from inventory.bulk_uom_conversion_line where tenant_id=$1`,
			`delete from inventory.bulk_uom_conversion where tenant_id=$1`,
			`delete from inventory.warehouse_pick_request where tenant_id=$1`,
			`delete from inventory.warehouse_activity_line where tenant_id=$1`,
			`delete from inventory.warehouse_activity where tenant_id=$1`,
			`delete from inventory.warehouse_receipt_line where tenant_id=$1`,
			`delete from inventory.warehouse_receipt where tenant_id=$1`,
			`delete from inventory.bulk_cost_application where tenant_id=$1`,
			`delete from inventory.bulk_cost_layer where tenant_id=$1`,
			`delete from inventory.bulk_inventory_entry where tenant_id=$1`,
			`delete from inventory.bulk_reservation where tenant_id=$1`,
			`delete from inventory.bulk_balance where tenant_id=$1`,
			`delete from inventory.item_bin_policy where tenant_id=$1`,
			`delete from inventory.warehouse_bin where tenant_id=$1`,
			`delete from inventory.item_unit_of_measure where tenant_id=$1`,
			`delete from inventory.stock_item where tenant_id=$1`,
			`delete from platform.outbox_event where tenant_id=$1`,
			`delete from org.organization where tenant_id=$1`,
			`delete from platform.tenant where tenant_id=$1`,
			`alter table inventory.bulk_uom_conversion_line enable trigger bulk_uom_conversion_line_immutable`,
			`alter table inventory.bulk_uom_conversion enable trigger bulk_uom_conversion_immutable`,
			`alter table inventory.warehouse_activity_line enable trigger warehouse_activity_line_immutable`,
			`alter table inventory.warehouse_receipt_line enable trigger warehouse_receipt_line_immutable`,
			`alter table inventory.warehouse_receipt enable trigger warehouse_receipt_immutable`,
			`alter table inventory.bulk_inventory_entry enable trigger bulk_inventory_entry_immutable`,
			`alter table inventory.item_unit_of_measure enable trigger item_unit_of_measure_immutable`,
			`alter table inventory.warehouse_activity_uom_reservation enable trigger warehouse_uom_reservation_delete_guard`,
			`alter table inventory.warehouse_pick_request enable trigger warehouse_pick_request_immutable`,
		}
		for _, command := range commands {
			arguments := []any{}
			if strings.Contains(command, "$1") {
				arguments = append(arguments, tenant)
			}
			if _, cleanupErr = tx.Exec(ctx, command, arguments...); cleanupErr != nil {
				t.Errorf("warehouse packaging cleanup: %v", cleanupErr)
				return
			}
		}
		if cleanupErr = tx.Commit(ctx); cleanupErr != nil {
			t.Errorf("warehouse packaging cleanup commit: %v", cleanupErr)
		}
	}()

	repo := NewInventoryControl(pool)
	for _, item := range []string{"preserve", "split", "cancel", "pick-split", "replenish-split"} {
		if _, err = repo.CreateBulkItem(ctx, tenant, bulkTestUUID(t), inventorycontrol.BulkItem{ID: item, Code: "PACK-" + item, Description: item, BaseUOM: "EA", BaseRoundingPrecision: "1", TrackingMode: "none", CostingMethod: "fifo", Version: 1}); err != nil {
			t.Fatal(err)
		}
		if _, err = repo.ConfigureItemUnitOfMeasure(ctx, tenant, bulkTestUUID(t), inventorycontrol.ItemUnitOfMeasure{ItemID: item, Code: "BOX", QuantityPerUnit: "12", RoundingPrecision: "1", Version: 1}); err != nil {
			t.Fatal(err)
		}
	}
	bins := []inventorycontrol.WarehouseBin{
		{ID: "receive", OrganizationID: "warehouse", Code: "RECEIVE", Type: "receive", Version: 1},
		{ID: "target-a", OrganizationID: "warehouse", Code: "TARGET-A", Type: "putpick", Ranking: 20, Version: 1},
		{ID: "target-b", OrganizationID: "warehouse", Code: "TARGET-B", Type: "putpick", Ranking: 10, Version: 1},
		{ID: "half-a", OrganizationID: "warehouse", Code: "HALF-A", Type: "put-away", Ranking: 20, Version: 1},
		{ID: "half-b", OrganizationID: "warehouse", Code: "HALF-B", Type: "put-away", Ranking: 10, Version: 1},
		{ID: "target-c", OrganizationID: "warehouse", Code: "TARGET-C", Type: "put-away", Ranking: 20, Version: 1},
		{ID: "pick-source", OrganizationID: "warehouse", Code: "PICK-SOURCE", Type: "putpick", Ranking: 20, Version: 1},
		{ID: "replenish-source", OrganizationID: "warehouse", Code: "REPLENISH-SOURCE", Type: "putpick", Ranking: 30, Version: 1},
		{ID: "replenish-target", OrganizationID: "warehouse", Code: "REPLENISH-TARGET", Type: "pick", Ranking: 100, Version: 1},
		{ID: "ship", OrganizationID: "warehouse", Code: "SHIP", Type: "ship", Version: 1},
	}
	for _, bin := range bins {
		if _, err = repo.CreateWarehouseBin(ctx, tenant, bulkTestUUID(t), bin); err != nil {
			t.Fatal(err)
		}
	}
	policies := []struct{ item, bin, maximum string }{
		{"preserve", "receive", "100"}, {"preserve", "target-a", "12"}, {"preserve", "target-b", "12"}, {"preserve", "ship", "100"},
		{"split", "receive", "100"}, {"split", "half-a", "6"}, {"split", "half-b", "6"},
		{"cancel", "receive", "100"}, {"cancel", "target-c", "12"},
		{"pick-split", "receive", "100"}, {"pick-split", "pick-source", "12"}, {"pick-split", "ship", "100"},
		{"replenish-split", "receive", "100"}, {"replenish-split", "replenish-source", "12"}, {"replenish-split", "replenish-target", "6"},
	}
	for _, policy := range policies {
		isDefault := policy.bin == "target-a" || policy.bin == "half-a" || policy.bin == "target-c" || policy.bin == "pick-source" || policy.bin == "replenish-source"
		minimum := "0"
		if policy.bin == "replenish-target" {
			minimum = "5"
		}
		if _, err = repo.ConfigureItemBin(ctx, tenant, bulkTestUUID(t), inventorycontrol.ItemBinPolicy{OrganizationID: "warehouse", ItemID: policy.item, BinID: policy.bin, Fixed: true, Default: isDefault, MinQuantity: minimum, MaxQuantity: policy.maximum, Version: 1}); err != nil {
			t.Fatal(err)
		}
	}
	posting := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	preserve, err := repo.PostWarehouseReceipt(ctx, tenant, warehousePackagingIDs(t, "preserve"), inventorycontrol.WarehouseReceiptCommand{RequestID: "preserve-request", OrganizationID: "warehouse", ReceiveBinID: "receive", ItemID: "preserve", HandlingUOM: "BOX", HandlingQuantity: "2", UnitCost: "10", PostingDate: posting, SourceKind: "purchase", SourceID: "preserve-po"})
	if err != nil || !exactDecimalEqual(preserve.Quantity, "24") || preserve.HandlingUOM != "BOX" || len(preserve.PutAway.Lines) != 2 {
		t.Fatalf("preserve receipt=%+v err=%v", preserve, err)
	}
	for _, line := range preserve.PutAway.Lines {
		if !exactDecimalEqual(line.Quantity, "12") || line.UOMCode != "BOX" || !exactDecimalEqual(line.UOMQuantity, "1") || !exactDecimalEqual(line.QuantityPerUOM, "12") {
			t.Fatalf("preserved line=%+v", line)
		}
	}
	var physical, physicalReserved, composed, composedReserved string
	if err = pool.QueryRow(ctx, `select b.quantity::text,b.reserved_quantity::text,u.quantity_base::text,u.reserved_quantity::text from inventory.bulk_balance b join inventory.bulk_uom_balance u using(balance_id) where b.tenant_id=$1 and b.item_id='preserve' and b.bin_id='receive' and u.uom_code='BOX'`, tenant).Scan(&physical, &physicalReserved, &composed, &composedReserved); err != nil || physical != "24.000000" || physicalReserved != "24.000000" || composed != "24.000000" || composedReserved != "2.000000" {
		t.Fatalf("preserve reserve physical=%s/%s composed=%s/%s err=%v", physical, physicalReserved, composed, composedReserved, err)
	}
	_, err = repo.ConvertBulkHandlingUnits(ctx, tenant, inventorycontrol.PackagingConversionIDs{ConversionID: "reserved-conversion", TakeLineID: "reserved-take", PlaceLineID: "reserved-place", EventID: bulkTestUUID(t)}, inventorycontrol.PackagingConversionCommand{RequestID: "reserved-conversion-request", OrganizationID: "warehouse", BinID: "receive", ItemID: "preserve", Operation: "breakbulk", FromUOM: "BOX", ToUOM: "EA", FromQuantity: "1"})
	if !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("reserved package conversion error=%v", err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "warehouse", preserve.PutAway.ID, 1, posting.AddDate(0, 0, 1), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	var targetPackages, targetReserved string
	if err = pool.QueryRow(ctx, `select sum(u.quantity)::text,sum(u.reserved_quantity)::text from inventory.bulk_balance b join inventory.bulk_uom_balance u using(balance_id) where b.tenant_id=$1 and b.item_id='preserve' and b.bin_id in ('target-a','target-b') and u.uom_code='BOX'`, tenant).Scan(&targetPackages, &targetReserved); err != nil || targetPackages != "2.000000" || targetReserved != "0.000000" {
		t.Fatalf("target packages=%s reserved=%s err=%v", targetPackages, targetReserved, err)
	}

	splitCommand := inventorycontrol.WarehouseReceiptCommand{RequestID: "split-rejected-request", OrganizationID: "warehouse", ReceiveBinID: "receive", ItemID: "split", HandlingUOM: "BOX", HandlingQuantity: "1", UnitCost: "10", PostingDate: posting, SourceKind: "purchase", SourceID: "split-rejected-po"}
	if _, err = repo.PostWarehouseReceipt(ctx, tenant, warehousePackagingIDs(t, "split-rejected"), splitCommand); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("breakbulk without authorization error=%v", err)
	}
	var rejectedResidue int
	if err = pool.QueryRow(ctx, `select count(*) from inventory.warehouse_receipt where tenant_id=$1 and receipt_id='split-rejected-receipt'`, tenant).Scan(&rejectedResidue); err != nil || rejectedResidue != 0 {
		t.Fatalf("rejected receipt residue=%d err=%v", rejectedResidue, err)
	}
	splitCommand.RequestID, splitCommand.SourceID, splitCommand.AllowBreakbulk = "split-allowed-request", "split-allowed-po", true
	split, err := repo.PostWarehouseReceipt(ctx, tenant, warehousePackagingIDs(t, "split-allowed"), splitCommand)
	if err != nil || len(split.PutAway.Lines) != 2 {
		t.Fatalf("authorized split=%+v err=%v", split, err)
	}
	for _, line := range split.PutAway.Lines {
		if !exactDecimalEqual(line.Quantity, "6") || line.UOMCode != "EA" || !exactDecimalEqual(line.UOMQuantity, "6") {
			t.Fatalf("breakbulk line=%+v", line)
		}
	}
	var sourceKind, sourceID string
	var conversionLines, conversionEvents int
	if err = pool.QueryRow(ctx, `select c.source_kind,c.source_id,(select count(*) from inventory.bulk_uom_conversion_line l where l.tenant_id=c.tenant_id and l.conversion_id=c.conversion_id),(select count(*) from platform.outbox_event e where e.tenant_id=c.tenant_id and e.aggregate_id=c.conversion_id and e.event_type='bulk-uom-conversion.completed') from inventory.bulk_uom_conversion c where c.tenant_id=$1 and c.conversion_id='split-allowed-conversion'`, tenant).Scan(&sourceKind, &sourceID, &conversionLines, &conversionEvents); err != nil || sourceKind != "warehouse-receipt" || sourceID != "split-allowed-receipt" || conversionLines != 2 || conversionEvents != 1 {
		t.Fatalf("automatic conversion source=%s/%s lines=%d events=%d err=%v", sourceKind, sourceID, conversionLines, conversionEvents, err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "warehouse", split.PutAway.ID, 1, posting.AddDate(0, 0, 2), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}

	cancelReceipt, err := repo.PostWarehouseReceipt(ctx, tenant, warehousePackagingIDs(t, "cancel"), inventorycontrol.WarehouseReceiptCommand{RequestID: "cancel-request", OrganizationID: "warehouse", ReceiveBinID: "receive", ItemID: "cancel", HandlingUOM: "BOX", HandlingQuantity: "1", UnitCost: "10", PostingDate: posting, SourceKind: "purchase", SourceID: "cancel-po"})
	if err != nil {
		t.Fatal(err)
	}
	cancelled, err := repo.CancelWarehousePutAway(ctx, tenant, "warehouse", cancelReceipt.PutAway.ID, 1, bulkTestUUID(t))
	if err != nil || cancelled.Status != "cancelled" || cancelled.Version != 2 {
		t.Fatalf("cancelled put-away=%+v err=%v", cancelled, err)
	}
	if err = pool.QueryRow(ctx, `select b.quantity::text,b.reserved_quantity::text,u.quantity::text,u.reserved_quantity::text from inventory.bulk_balance b join inventory.bulk_uom_balance u using(balance_id) where b.tenant_id=$1 and b.item_id='cancel' and b.bin_id='receive' and u.uom_code='BOX'`, tenant).Scan(&physical, &physicalReserved, &composed, &composedReserved); err != nil || physical != "12.000000" || physicalReserved != "0.000000" || composed != "1.000000" || composedReserved != "0.000000" {
		t.Fatalf("cancel release physical=%s/%s composed=%s/%s err=%v", physical, physicalReserved, composed, composedReserved, err)
	}
	if _, err = repo.ConvertBulkHandlingUnits(ctx, tenant, inventorycontrol.PackagingConversionIDs{ConversionID: "cancel-conversion", TakeLineID: "cancel-take", PlaceLineID: "cancel-place", EventID: bulkTestUUID(t)}, inventorycontrol.PackagingConversionCommand{RequestID: "cancel-conversion-request", OrganizationID: "warehouse", BinID: "receive", ItemID: "cancel", Operation: "breakbulk", FromUOM: "BOX", ToUOM: "EA", FromQuantity: "1"}); err != nil {
		t.Fatalf("conversion after cancellation: %v", err)
	}

	pickReceipt, err := repo.PostWarehouseReceipt(ctx, tenant, warehousePackagingIDs(t, "pick-split"), inventorycontrol.WarehouseReceiptCommand{RequestID: "pick-split-receipt-request", OrganizationID: "warehouse", ReceiveBinID: "receive", ItemID: "pick-split", HandlingUOM: "BOX", HandlingQuantity: "1", UnitCost: "10", PostingDate: posting, SourceKind: "purchase", SourceID: "pick-split-po"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "warehouse", pickReceipt.PutAway.ID, 1, posting.AddDate(0, 0, 3), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	pickCommand := inventorycontrol.WarehousePickCommand{RequestID: "pick-no-breakbulk", OrganizationID: "warehouse", ShipBinID: "ship", ItemID: "pick-split", DemandKind: "manual", DemandID: "order-pack", DemandLineID: "line-pack", Quantity: "6", HandlingUOM: "EA"}
	if _, err = repo.CreateWarehousePick(ctx, tenant, warehouseTestIDs(t, "pick-no-breakbulk"), pickCommand); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("partial package pick without authorization error=%v", err)
	}
	pickCommand.RequestID, pickCommand.AllowBreakbulk = "pick-cancel-request", true
	pickToCancel, err := repo.CreateWarehousePick(ctx, tenant, warehouseTestIDs(t, "pick-cancel"), pickCommand)
	if err != nil || len(pickToCancel.Lines) != 1 || pickToCancel.Lines[0].FromUOMCode != "BOX" || !exactDecimalEqual(pickToCancel.Lines[0].FromUOMQuantity, "1") || pickToCancel.Lines[0].UOMCode != "EA" || !exactDecimalEqual(pickToCancel.Lines[0].UOMQuantity, "6") {
		t.Fatalf("packaging-aware pick=%+v err=%v", pickToCancel, err)
	}
	if err = pool.QueryRow(ctx, `select b.reserved_quantity::text,u.reserved_quantity::text from inventory.bulk_balance b join inventory.bulk_uom_balance u using(balance_id) where b.tenant_id=$1 and b.item_id='pick-split' and b.bin_id='pick-source' and u.uom_code='BOX'`, tenant).Scan(&physicalReserved, &composedReserved); err != nil || physicalReserved != "6.000000" || composedReserved != "1.000000" {
		t.Fatalf("pick reservations physical=%s package=%s err=%v", physicalReserved, composedReserved, err)
	}
	if _, err = repo.ConvertBulkHandlingUnits(ctx, tenant, inventorycontrol.PackagingConversionIDs{ConversionID: "pick-reserved-conversion", TakeLineID: "pick-reserved-take", PlaceLineID: "pick-reserved-place", EventID: bulkTestUUID(t)}, inventorycontrol.PackagingConversionCommand{RequestID: "pick-reserved-conversion-request", OrganizationID: "warehouse", BinID: "pick-source", ItemID: "pick-split", Operation: "breakbulk", FromUOM: "BOX", ToUOM: "EA", FromQuantity: "1"}); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("reserved pick package conversion error=%v", err)
	}
	if _, err = repo.CancelWarehousePick(ctx, tenant, "warehouse", pickToCancel.ID, 1, bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select b.reserved_quantity::text,u.reserved_quantity::text from inventory.bulk_balance b join inventory.bulk_uom_balance u using(balance_id) where b.tenant_id=$1 and b.item_id='pick-split' and b.bin_id='pick-source' and u.uom_code='BOX'`, tenant).Scan(&physicalReserved, &composedReserved); err != nil || physicalReserved != "0.000000" || composedReserved != "0.000000" {
		t.Fatalf("cancelled pick reservations physical=%s package=%s err=%v", physicalReserved, composedReserved, err)
	}
	pickCommand.RequestID = "pick-register-request"
	pick, err := repo.CreateWarehousePick(ctx, tenant, warehouseTestIDs(t, "pick-register"), pickCommand)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "warehouse", pick.ID, 1, posting.AddDate(0, 0, 4), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	var pickSourceBase, pickSourceEA, pickShipBase, pickShipEA string
	if err = pool.QueryRow(ctx, `select (select quantity::text from inventory.bulk_balance where tenant_id=$1 and item_id='pick-split' and bin_id='pick-source'),(select u.quantity::text from inventory.bulk_uom_balance u join inventory.bulk_balance b using(balance_id) where b.tenant_id=$1 and b.item_id='pick-split' and b.bin_id='pick-source' and u.uom_code='EA'),(select quantity::text from inventory.bulk_balance where tenant_id=$1 and item_id='pick-split' and bin_id='ship'),(select u.quantity::text from inventory.bulk_uom_balance u join inventory.bulk_balance b using(balance_id) where b.tenant_id=$1 and b.item_id='pick-split' and b.bin_id='ship' and u.uom_code='EA')`, tenant).Scan(&pickSourceBase, &pickSourceEA, &pickShipBase, &pickShipEA); err != nil || pickSourceBase != "6.000000" || pickSourceEA != "6.000000" || pickShipBase != "6.000000" || pickShipEA != "6.000000" {
		t.Fatalf("registered pick source=%s/%s ship=%s/%s err=%v", pickSourceBase, pickSourceEA, pickShipBase, pickShipEA, err)
	}
	if err = pool.QueryRow(ctx, `select source_kind,source_id from inventory.bulk_uom_conversion where tenant_id=$1 and source_kind='warehouse-pick' and source_id=$2`, tenant, pick.ID).Scan(&sourceKind, &sourceID); err != nil || sourceKind != "warehouse-pick" || sourceID != pick.ID {
		t.Fatalf("pick conversion source=%s/%s err=%v", sourceKind, sourceID, err)
	}

	replenishReceipt, err := repo.PostWarehouseReceipt(ctx, tenant, warehousePackagingIDs(t, "replenish-split"), inventorycontrol.WarehouseReceiptCommand{RequestID: "replenish-split-receipt-request", OrganizationID: "warehouse", ReceiveBinID: "receive", ItemID: "replenish-split", HandlingUOM: "BOX", HandlingQuantity: "1", UnitCost: "10", PostingDate: posting, SourceKind: "purchase", SourceID: "replenish-split-po"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "warehouse", replenishReceipt.PutAway.ID, 1, posting.AddDate(0, 0, 5), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	replenishCommand := inventorycontrol.WarehouseReplenishmentCommand{RequestID: "replenish-no-breakbulk", OrganizationID: "warehouse", ToBinID: "replenish-target", ItemID: "replenish-split", TargetUOM: "EA"}
	if _, err = repo.CreateWarehouseReplenishment(ctx, tenant, warehouseTestIDs(t, "replenish-no-breakbulk"), replenishCommand); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("partial replenishment without authorization error=%v", err)
	}
	replenishCommand.RequestID, replenishCommand.AllowBreakbulk = "replenish-breakbulk", true
	replenishment, err := repo.CreateWarehouseReplenishment(ctx, tenant, warehouseTestIDs(t, "replenish-breakbulk"), replenishCommand)
	if err != nil || len(replenishment.Lines) != 1 || replenishment.Lines[0].FromUOMCode != "BOX" || replenishment.Lines[0].UOMCode != "EA" || !exactDecimalEqual(replenishment.Lines[0].Quantity, "6") {
		t.Fatalf("packaging-aware replenishment=%+v err=%v", replenishment, err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "warehouse", replenishment.ID, 1, posting.AddDate(0, 0, 6), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	var replenished string
	if err = pool.QueryRow(ctx, `select u.quantity::text from inventory.bulk_uom_balance u join inventory.bulk_balance b using(balance_id) where b.tenant_id=$1 and b.item_id='replenish-split' and b.bin_id='replenish-target' and u.uom_code='EA'`, tenant).Scan(&replenished); err != nil || replenished != "6.000000" {
		t.Fatalf("replenished EA=%s err=%v", replenished, err)
	}

	exactPick, err := repo.CreateWarehousePick(ctx, tenant, warehouseTestIDs(t, "pick-exact-box"), inventorycontrol.WarehousePickCommand{RequestID: "pick-exact-box-request", OrganizationID: "warehouse", ShipBinID: "ship", ItemID: "preserve", DemandKind: "manual", DemandID: "order-box", DemandLineID: "line-box", Quantity: "12", HandlingUOM: "BOX"})
	if err != nil || len(exactPick.Lines) != 1 || exactPick.Lines[0].FromUOMCode != "BOX" || exactPick.Lines[0].UOMCode != "BOX" || !exactDecimalEqual(exactPick.Lines[0].FromUOMQuantity, "1") || !exactDecimalEqual(exactPick.Lines[0].UOMQuantity, "1") {
		t.Fatalf("exact-package pick=%+v err=%v", exactPick, err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "warehouse", exactPick.ID, 1, posting.AddDate(0, 0, 7), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	var exactShipPackages string
	if err = pool.QueryRow(ctx, `select u.quantity::text from inventory.bulk_uom_balance u join inventory.bulk_balance b using(balance_id) where b.tenant_id=$1 and b.item_id='preserve' and b.bin_id='ship' and u.uom_code='BOX'`, tenant).Scan(&exactShipPackages); err != nil || exactShipPackages != "1.000000" {
		t.Fatalf("exact package at ship=%s err=%v", exactShipPackages, err)
	}
}
