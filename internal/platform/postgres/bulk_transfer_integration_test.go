package postgres

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/inventorycontrol"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestBulkTransferShipTransitReceiveConservesQuantityLotAndCost(t *testing.T) {
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
	if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Transfer V164','Transfer V164')`, tenant, "tr-"+tenant[:8]); err != nil {
		t.Fatal(err)
	}
	defer func() {
		tx, cleanupErr := pool.Begin(ctx)
		if cleanupErr != nil {
			t.Errorf("cleanup begin: %v", cleanupErr)
			return
		}
		defer tx.Rollback(ctx)
		commands := []string{
			`alter table inventory.warehouse_pick_request disable trigger warehouse_pick_request_immutable`,
			`alter table inventory.warehouse_activity_uom_reservation disable trigger warehouse_uom_reservation_delete_guard`,
			`alter table inventory.warehouse_activity_line disable trigger warehouse_activity_line_immutable`,
			`alter table inventory.warehouse_receipt_line disable trigger warehouse_receipt_line_immutable`,
			`alter table inventory.warehouse_receipt disable trigger warehouse_receipt_immutable`,
			`alter table inventory.bulk_cost_application disable trigger bulk_cost_application_immutable`,
			`alter table inventory.bulk_inventory_entry disable trigger bulk_inventory_entry_immutable`,
			`alter table inventory.item_unit_of_measure disable trigger item_unit_of_measure_immutable`,
			`alter table inventory.bulk_uom_conversion disable trigger bulk_uom_conversion_immutable`,
			`alter table inventory.bulk_uom_conversion_line disable trigger bulk_uom_conversion_line_immutable`,
			`delete from inventory.warehouse_activity_uom_reservation where tenant_id=$1`,
			`delete from inventory.bulk_transfer_receipt_component where tenant_id=$1`,
			`delete from inventory.bulk_transfer_cost_component where tenant_id=$1`,
			`delete from inventory.bulk_transfer_allocation where tenant_id=$1`,
			`delete from inventory.bulk_transfer_receipt where tenant_id=$1`,
			`delete from inventory.bulk_transfer_shipment where tenant_id=$1`,
			`delete from inventory.bulk_transfer_line where tenant_id=$1`,
			`delete from inventory.bulk_transfer where tenant_id=$1`,
			`delete from inventory.warehouse_pick_request where tenant_id=$1`,
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
			`alter table inventory.bulk_cost_application enable trigger bulk_cost_application_immutable`,
			`alter table inventory.bulk_inventory_entry enable trigger bulk_inventory_entry_immutable`,
			`alter table inventory.item_unit_of_measure enable trigger item_unit_of_measure_immutable`,
			`alter table inventory.bulk_uom_conversion_line enable trigger bulk_uom_conversion_line_immutable`,
			`alter table inventory.bulk_uom_conversion enable trigger bulk_uom_conversion_immutable`,
			`alter table inventory.warehouse_activity_uom_reservation enable trigger warehouse_uom_reservation_delete_guard`,
			`alter table inventory.warehouse_pick_request enable trigger warehouse_pick_request_immutable`,
		}
		for _, command := range commands {
			args := []any{}
			if len(command) >= 6 && command[:6] == "delete" {
				args = []any{tenant}
			}
			if _, cleanupErr = tx.Exec(ctx, command, args...); cleanupErr != nil {
				t.Errorf("cleanup %s: %v", command, cleanupErr)
				return
			}
		}
		if cleanupErr = tx.Commit(ctx); cleanupErr != nil {
			t.Errorf("cleanup commit: %v", cleanupErr)
		}
	}()
	if _, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'source','source','Source','warehouse'),($1,'destination','destination','Destination','warehouse')`, tenant); err != nil {
		t.Fatal(err)
	}
	repo := NewInventoryControl(pool)
	if _, err = repo.CreateBulkItem(ctx, tenant, bulkTestUUID(t), inventorycontrol.BulkItem{ID: "part", Code: "TRANSFER-PART", Description: "Transfer part", BaseUOM: "EA", BaseRoundingPrecision: "0.000001", TrackingMode: "lot", CostingMethod: "fifo", Version: 1}); err != nil {
		t.Fatal(err)
	}
	if _, err = repo.ConfigureItemUnitOfMeasure(ctx, tenant, bulkTestUUID(t), inventorycontrol.ItemUnitOfMeasure{ItemID: "part", Code: "BOX", QuantityPerUnit: "2", RoundingPrecision: "1", Version: 1}); err != nil {
		t.Fatal(err)
	}
	bins := []inventorycontrol.WarehouseBin{{ID: "source-receive", OrganizationID: "source", Code: "RECEIVE", Type: "receive"}, {ID: "source-pick", OrganizationID: "source", Code: "PICK", Type: "putpick", Ranking: 100}, {ID: "source-ship", OrganizationID: "source", Code: "SHIP", Type: "ship"}, {ID: "destination-receive", OrganizationID: "destination", Code: "RECEIVE", Type: "receive"}, {ID: "destination-pick", OrganizationID: "destination", Code: "PICK", Type: "putpick", Ranking: 100}}
	for _, bin := range bins {
		if _, err = repo.CreateWarehouseBin(ctx, tenant, bulkTestUUID(t), bin); err != nil {
			t.Fatal(err)
		}
		if _, err = repo.ConfigureItemBin(ctx, tenant, bulkTestUUID(t), inventorycontrol.ItemBinPolicy{OrganizationID: bin.OrganizationID, ItemID: "part", BinID: bin.ID, Fixed: true, Default: bin.Type == "putpick", MinQuantity: "0", MaxQuantity: "100", Version: 1}); err != nil {
			t.Fatal(err)
		}
	}
	expiry := time.Date(2035, 1, 1, 0, 0, 0, 0, time.UTC)
	receipt, err := repo.PostWarehouseReceipt(ctx, tenant, warehouseTestIDs(t, "transfer-source"), inventorycontrol.WarehouseReceiptCommand{RequestID: "source-receipt", OrganizationID: "source", ReceiveBinID: "source-receive", ItemID: "part", LotNo: "LOT-TRANSFER", ExpirationDate: &expiry, HandlingUOM: "BOX", HandlingQuantity: "3", UnitCost: "10", PostingDate: time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC), SourceKind: "purchase", SourceID: "po-transfer"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "source", receipt.PutAway.ID, 1, time.Date(2030, 1, 2, 0, 0, 0, 0, time.UTC), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	transfer, err := repo.CreateBulkTransfer(ctx, tenant, "transfer-1", "transfer-line-1", bulkTestUUID(t), inventorycontrol.BulkTransferCommand{RequestID: "transfer-request", FromOrganizationID: "source", ToOrganizationID: "destination", ReceiveBinID: "destination-receive", InTransitCode: "OWN-LOG", ItemID: "part", Quantity: "3", PostingDate: time.Date(2030, 1, 3, 0, 0, 0, 0, time.UTC), SourceKind: "replenishment", SourceID: "plan-1"})
	if err != nil || transfer.Status != "released" {
		t.Fatalf("create=%+v err=%v", transfer, err)
	}
	replay, err := repo.CreateBulkTransfer(ctx, tenant, "ignored-transfer", "ignored-line", bulkTestUUID(t), inventorycontrol.BulkTransferCommand{RequestID: "transfer-request", FromOrganizationID: "source", ToOrganizationID: "destination", ReceiveBinID: "destination-receive", InTransitCode: "OWN-LOG", ItemID: "part", Quantity: "3", PostingDate: time.Date(2030, 1, 3, 0, 0, 0, 0, time.UTC), SourceKind: "replenishment", SourceID: "plan-1"})
	if err != nil || replay.ID != transfer.ID {
		t.Fatalf("idempotent replay=%+v err=%v", replay, err)
	}
	firstPick, err := repo.CreateWarehousePick(ctx, tenant, warehouseTestIDs(t, "transfer-pick-first"), inventorycontrol.WarehousePickCommand{RequestID: "transfer-pick-first-request", OrganizationID: "source", ShipBinID: "source-ship", ItemID: "part", DemandKind: "transfer-outbound", DemandID: transfer.ID, DemandLineID: transfer.LineID, Quantity: "1", AllowBreakbulk: true, UseFEFO: true})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "source", firstPick.ID, 1, time.Date(2030, 1, 4, 0, 0, 0, 0, time.UTC), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	firstShipCommand := inventorycontrol.BulkTransferPostingCommand{RequestID: "ship-first-request", Quantity: "1", WarehouseActivityID: firstPick.ID, PostingDate: time.Date(2030, 1, 5, 0, 0, 0, 0, time.UTC)}
	firstShipped, err := repo.ShipBulkTransfer(ctx, tenant, transfer.ID, "source", "destination", 1, "shipment-first", bulkTestUUID(t), firstShipCommand)
	if err != nil || firstShipped.Status != "partially-shipped" || firstShipped.Version != 2 || firstShipped.ShippedQuantity != "1.000000" || firstShipped.CostAmount != "10.0000" || firstShipped.PostingID != "shipment-first" {
		t.Fatalf("first partial shipment=%+v err=%v", firstShipped, err)
	}
	firstShipReplay, err := repo.ShipBulkTransfer(ctx, tenant, transfer.ID, "source", "destination", 1, "ignored-shipment", bulkTestUUID(t), firstShipCommand)
	if err != nil || firstShipReplay.PostingID != firstShipped.PostingID || firstShipReplay.Version != 2 {
		t.Fatalf("shipment replay=%+v err=%v", firstShipReplay, err)
	}
	divergentShip := firstShipCommand
	divergentShip.Quantity = "2"
	if _, err = repo.ShipBulkTransfer(ctx, tenant, transfer.ID, "source", "destination", 1, "ignored-divergent", bulkTestUUID(t), divergentShip); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("divergent shipment replay: %v", err)
	}
	firstReceiptCommand := inventorycontrol.BulkTransferPostingCommand{RequestID: "receive-first-request", Quantity: "1", PostingDate: time.Date(2030, 1, 6, 0, 0, 0, 0, time.UTC)}
	firstReceived, err := repo.ReceiveBulkTransfer(ctx, tenant, transfer.ID, "source", "destination", 2, "receipt-first", "put-away-first", bulkTestUUID(t), firstReceiptCommand)
	if err != nil || firstReceived.Status != "partially-received" || firstReceived.Version != 3 || firstReceived.ReceivedQuantity != "1.000000" || firstReceived.PostingID != "receipt-first" || firstReceived.PutAwayID != "put-away-first" {
		t.Fatalf("first partial receipt=%+v err=%v", firstReceived, err)
	}
	firstReceiptReplay, err := repo.ReceiveBulkTransfer(ctx, tenant, transfer.ID, "source", "destination", 2, "ignored-receipt", "ignored-put-away", bulkTestUUID(t), firstReceiptCommand)
	if err != nil || firstReceiptReplay.PostingID != firstReceived.PostingID || firstReceiptReplay.PutAwayID != firstReceived.PutAwayID {
		t.Fatalf("receipt replay=%+v err=%v", firstReceiptReplay, err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "destination", firstReceived.PutAwayID, 1, time.Date(2030, 1, 7, 0, 0, 0, 0, time.UTC), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	secondPick, err := repo.CreateWarehousePick(ctx, tenant, warehouseTestIDs(t, "transfer-pick-second"), inventorycontrol.WarehousePickCommand{RequestID: "transfer-pick-second-request", OrganizationID: "source", ShipBinID: "source-ship", ItemID: "part", DemandKind: "transfer-outbound", DemandID: transfer.ID, DemandLineID: transfer.LineID, Quantity: "2", HandlingUOM: "BOX", UseFEFO: true})
	if err != nil || len(secondPick.Lines) != 1 || secondPick.Lines[0].UOMCode != "BOX" || !exactDecimalEqual(secondPick.Lines[0].UOMQuantity, "1") {
		t.Fatalf("second packaged pick=%+v err=%v", secondPick, err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "source", secondPick.ID, 1, time.Date(2030, 1, 8, 0, 0, 0, 0, time.UTC), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	results := make([]inventorycontrol.BulkTransfer, 2)
	failures := make([]error, 2)
	shipmentIDs := []string{"shipment-second-a", "shipment-second-b"}
	eventIDs := []string{bulkTestUUID(t), bulkTestUUID(t)}
	commands := []inventorycontrol.BulkTransferPostingCommand{{RequestID: "ship-second-a", Quantity: "2", WarehouseActivityID: secondPick.ID, PostingDate: time.Date(2030, 1, 9, 0, 0, 0, 0, time.UTC)}, {RequestID: "ship-second-b", Quantity: "2", WarehouseActivityID: secondPick.ID, PostingDate: time.Date(2030, 1, 9, 0, 0, 0, 0, time.UTC)}}
	var group sync.WaitGroup
	for index := range results {
		group.Add(1)
		go func(i int) {
			defer group.Done()
			results[i], failures[i] = repo.ShipBulkTransfer(ctx, tenant, transfer.ID, "source", "destination", 3, shipmentIDs[i], eventIDs[i], commands[i])
		}(index)
	}
	group.Wait()
	successes := 0
	for index, failure := range failures {
		if failure == nil {
			successes++
			if results[index].Status != "partially-received" || results[index].Version != 4 || results[index].CostAmount != "30.0000" {
				t.Fatalf("shipment=%+v", results[index])
			}
		} else if !errors.Is(failure, inventorycontrol.ErrConflict) {
			t.Fatalf("shipment error=%v", failure)
		}
	}
	if successes != 1 {
		t.Fatalf("shipment successes=%d errors=%v", successes, failures)
	}
	var winningShipment inventorycontrol.BulkTransfer
	var winningCommand inventorycontrol.BulkTransferPostingCommand
	for index := range results {
		if failures[index] == nil {
			winningShipment, winningCommand = results[index], commands[index]
		}
	}
	winningReplay, err := repo.ShipBulkTransfer(ctx, tenant, transfer.ID, "source", "destination", 3, "ignored-winning-replay", bulkTestUUID(t), winningCommand)
	if err != nil || winningReplay.PostingID != winningShipment.PostingID || winningReplay.Version != 4 {
		t.Fatalf("winning shipment replay=%+v err=%v", winningReplay, err)
	}
	var transitQuantity, transitCost, transitLot string
	if err = pool.QueryRow(ctx, `select quantity::text,cost_amount::text,lot_id from inventory.bulk_inventory_in_transit where tenant_id=$1 and transfer_id=$2`, tenant, transfer.ID).Scan(&transitQuantity, &transitCost, &transitLot); err != nil || transitQuantity != "2.000000" || transitCost != "20.0000" || transitLot != "transfer-source-lot" {
		t.Fatalf("transit quantity=%s cost=%s lot=%s err=%v", transitQuantity, transitCost, transitLot, err)
	}
	secondReceiptCommand := inventorycontrol.BulkTransferPostingCommand{RequestID: "receive-second-request", Quantity: "1", PostingDate: time.Date(2030, 1, 10, 0, 0, 0, 0, time.UTC)}
	secondReceived, err := repo.ReceiveBulkTransfer(ctx, tenant, transfer.ID, "source", "destination", 4, "receipt-second", "put-away-second", bulkTestUUID(t), secondReceiptCommand)
	if err != nil || secondReceived.Status != "partially-received" || secondReceived.Version != 5 || secondReceived.ReceivedQuantity != "2.000000" {
		t.Fatalf("second partial receipt=%+v err=%v", secondReceived, err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "destination", secondReceived.PutAwayID, 1, time.Date(2030, 1, 11, 0, 0, 0, 0, time.UTC), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select quantity::text,cost_amount::text,lot_id from inventory.bulk_inventory_in_transit where tenant_id=$1 and transfer_id=$2`, tenant, transfer.ID).Scan(&transitQuantity, &transitCost, &transitLot); err != nil || transitQuantity != "1.000000" || transitCost != "10.0000" {
		t.Fatalf("partial transit quantity=%s cost=%s err=%v", transitQuantity, transitCost, err)
	}
	finalReceiptCommand := inventorycontrol.BulkTransferPostingCommand{RequestID: "receive-final-request", Quantity: "1", PostingDate: time.Date(2030, 1, 12, 0, 0, 0, 0, time.UTC)}
	received, err := repo.ReceiveBulkTransfer(ctx, tenant, transfer.ID, "source", "destination", 5, "receipt-final", "put-away-final", bulkTestUUID(t), finalReceiptCommand)
	if err != nil || received.Status != "received" || received.Version != 6 || received.ReceivedQuantity != "3.000000" || received.CostAmount != "30.0000" || received.PutAwayID != "put-away-final" {
		t.Fatalf("final receive=%+v err=%v", received, err)
	}
	registeredPutAway, err := repo.RegisterWarehouseActivity(ctx, tenant, "destination", received.PutAwayID, 1, time.Date(2030, 1, 13, 0, 0, 0, 0, time.UTC), bulkTestUUID(t))
	if err != nil || registeredPutAway.Status != "registered" || len(registeredPutAway.Lines) != 1 || registeredPutAway.Lines[0].ToBinID != "destination-pick" {
		t.Fatalf("destination put-away=%+v err=%v", registeredPutAway, err)
	}
	var transitRows, shipmentRows, receiptRows int
	if err = pool.QueryRow(ctx, `select (select count(*) from inventory.bulk_inventory_in_transit where tenant_id=$1 and transfer_id=$2),(select count(*) from inventory.bulk_transfer_shipment where tenant_id=$1 and transfer_id=$2),(select count(*) from inventory.bulk_transfer_receipt where tenant_id=$1 and transfer_id=$2)`, tenant, transfer.ID).Scan(&transitRows, &shipmentRows, &receiptRows); err != nil || transitRows != 0 || shipmentRows != 2 || receiptRows != 3 {
		t.Fatalf("posting evidence transit=%d shipments=%d receipts=%d err=%v", transitRows, shipmentRows, receiptRows, err)
	}
	var sourceQuantity, destinationQuantity, destinationCost string
	if err = pool.QueryRow(ctx, `select (select sum(quantity)::text from inventory.bulk_balance where tenant_id=$1 and organization_id='source'),(select sum(quantity)::text from inventory.bulk_balance where tenant_id=$1 and organization_id='destination'),(select sum(cost_amount)::text from inventory.bulk_inventory_entry where tenant_id=$1 and organization_id='destination' and entry_type='receipt' and source_kind='bulk-transfer-receipt')`, tenant).Scan(&sourceQuantity, &destinationQuantity, &destinationCost); err != nil || sourceQuantity != "3.000000" || destinationQuantity != "3.000000" || destinationCost != "30.0000" {
		t.Fatalf("conservation source=%s destination=%s cost=%s err=%v", sourceQuantity, destinationQuantity, destinationCost, err)
	}
	cancellable, err := repo.CreateBulkTransfer(ctx, tenant, "transfer-cancel", "transfer-cancel-line", bulkTestUUID(t), inventorycontrol.BulkTransferCommand{RequestID: "transfer-cancel-request", FromOrganizationID: "source", ToOrganizationID: "destination", ReceiveBinID: "destination-receive", InTransitCode: "OWN-LOG", ItemID: "part", Quantity: "1", PostingDate: time.Date(2030, 1, 8, 0, 0, 0, 0, time.UTC), SourceKind: "replenishment", SourceID: "plan-cancel"})
	if err != nil {
		t.Fatal(err)
	}
	cancelPick, err := repo.CreateWarehousePick(ctx, tenant, warehouseTestIDs(t, "transfer-cancel-pick"), inventorycontrol.WarehousePickCommand{RequestID: "transfer-cancel-pick-request", OrganizationID: "source", ShipBinID: "source-ship", ItemID: "part", DemandKind: "transfer-outbound", DemandID: cancellable.ID, DemandLineID: cancellable.LineID, Quantity: "1", UseFEFO: true})
	if err != nil {
		t.Fatal(err)
	}
	cancelled, err := repo.CancelBulkTransfer(ctx, tenant, cancellable.ID, "source", "destination", 1, bulkTestUUID(t))
	if err != nil || cancelled.Status != "cancelled" || cancelled.Version != 2 {
		t.Fatalf("cancelled transfer=%+v err=%v", cancelled, err)
	}
	var cancelledActivity, releasedReservation, reservedAfterCancel string
	if err = pool.QueryRow(ctx, `select a.status,r.status,b.reserved_quantity::text from inventory.warehouse_activity a join inventory.warehouse_activity_line al using(tenant_id,activity_id) join inventory.bulk_reservation r on r.tenant_id=al.tenant_id and r.reservation_id=al.reservation_id join inventory.bulk_balance b on b.tenant_id=r.tenant_id and b.organization_id=r.organization_id and b.bin_id=r.bin_id and b.item_id=r.item_id and b.lot_id is not distinct from r.lot_id where a.tenant_id=$1 and a.activity_id=$2`, tenant, cancelPick.ID).Scan(&cancelledActivity, &releasedReservation, &reservedAfterCancel); err != nil || cancelledActivity != "cancelled" || releasedReservation != "released" || reservedAfterCancel != "0.000000" {
		t.Fatalf("cancel evidence activity=%s reservation=%s reserved=%s err=%v", cancelledActivity, releasedReservation, reservedAfterCancel, err)
	}
	if _, err = repo.ShipBulkTransfer(ctx, tenant, cancellable.ID, "source", "destination", 1, "cancelled-shipment", bulkTestUUID(t), inventorycontrol.BulkTransferPostingCommand{RequestID: "cancelled-ship-request", Quantity: "1", WarehouseActivityID: cancelPick.ID, PostingDate: time.Date(2030, 1, 14, 0, 0, 0, 0, time.UTC)}); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("cancelled transfer shipped: %v", err)
	}

	if _, err = repo.CreateBulkItem(ctx, tenant, bulkTestUUID(t), inventorycontrol.BulkItem{ID: "specific-part", Code: "TRANSFER-SPECIFIC", Description: "Specific-cost transfer part", BaseUOM: "EA", BaseRoundingPrecision: "0.000001", TrackingMode: "none", CostingMethod: "specific", Version: 1}); err != nil {
		t.Fatal(err)
	}
	for _, bin := range bins {
		if _, err = repo.ConfigureItemBin(ctx, tenant, bulkTestUUID(t), inventorycontrol.ItemBinPolicy{OrganizationID: bin.OrganizationID, ItemID: "specific-part", BinID: bin.ID, Fixed: true, Default: bin.Type == "putpick", MinQuantity: "0", MaxQuantity: "100", Version: 1}); err != nil {
			t.Fatal(err)
		}
	}
	firstSpecific, err := repo.PostWarehouseReceipt(ctx, tenant, warehouseTestIDs(t, "specific-first"), inventorycontrol.WarehouseReceiptCommand{RequestID: "specific-first-request", OrganizationID: "source", ReceiveBinID: "source-receive", ItemID: "specific-part", Quantity: "1", UnitCost: "11", PostingDate: time.Date(2030, 2, 1, 0, 0, 0, 0, time.UTC), SourceKind: "purchase", SourceID: "specific-po-first"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "source", firstSpecific.PutAway.ID, 1, time.Date(2030, 2, 2, 0, 0, 0, 0, time.UTC), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	secondSpecific, err := repo.PostWarehouseReceipt(ctx, tenant, warehouseTestIDs(t, "specific-second"), inventorycontrol.WarehouseReceiptCommand{RequestID: "specific-second-request", OrganizationID: "source", ReceiveBinID: "source-receive", ItemID: "specific-part", Quantity: "1", UnitCost: "17", PostingDate: time.Date(2030, 2, 3, 0, 0, 0, 0, time.UTC), SourceKind: "purchase", SourceID: "specific-po-second"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "source", secondSpecific.PutAway.ID, 1, time.Date(2030, 2, 4, 0, 0, 0, 0, time.UTC), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	specificCommand := inventorycontrol.BulkTransferCommand{RequestID: "specific-transfer-request", FromOrganizationID: "source", ToOrganizationID: "destination", ReceiveBinID: "destination-receive", InTransitCode: "OWN-LOG", ItemID: "specific-part", Quantity: "1", PostingDate: time.Date(2030, 2, 5, 0, 0, 0, 0, time.UTC), SourceKind: "replenishment", SourceID: "specific-plan"}
	if _, err = repo.CreateBulkTransfer(ctx, tenant, "specific-missing", "specific-missing-line", bulkTestUUID(t), specificCommand); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("specific transfer without receipt selector: %v", err)
	}
	specificCommand.SpecificReceiptEntry = secondSpecific.EntryID
	specificTransfer, err := repo.CreateBulkTransfer(ctx, tenant, "specific-transfer", "specific-transfer-line", bulkTestUUID(t), specificCommand)
	if err != nil || specificTransfer.SpecificReceiptEntry != secondSpecific.EntryID {
		t.Fatalf("specific create=%+v err=%v", specificTransfer, err)
	}
	mismatchedReplay := specificCommand
	mismatchedReplay.SpecificReceiptEntry = firstSpecific.EntryID
	if _, err = repo.CreateBulkTransfer(ctx, tenant, "ignored-specific", "ignored-specific-line", bulkTestUUID(t), mismatchedReplay); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("specific replay changed selector: %v", err)
	}
	specificPick, err := repo.CreateWarehousePick(ctx, tenant, warehouseTestIDs(t, "specific-transfer-pick"), inventorycontrol.WarehousePickCommand{RequestID: "specific-transfer-pick-request", OrganizationID: "source", ShipBinID: "source-ship", ItemID: "specific-part", DemandKind: "transfer-outbound", DemandID: specificTransfer.ID, DemandLineID: specificTransfer.LineID, Quantity: "1"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "source", specificPick.ID, 1, time.Date(2030, 2, 6, 0, 0, 0, 0, time.UTC), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	specificShipped, err := repo.ShipBulkTransfer(ctx, tenant, specificTransfer.ID, "source", "destination", 1, "specific-shipment", bulkTestUUID(t), inventorycontrol.BulkTransferPostingCommand{RequestID: "specific-ship-request", Quantity: "1", WarehouseActivityID: specificPick.ID, PostingDate: time.Date(2030, 2, 7, 0, 0, 0, 0, time.UTC)})
	if err != nil || specificShipped.Status != "shipped" || specificShipped.CostAmount != "17.0000" {
		t.Fatalf("specific shipment=%+v err=%v", specificShipped, err)
	}
	specificReceived, err := repo.ReceiveBulkTransfer(ctx, tenant, specificTransfer.ID, "source", "destination", 2, "specific-receipt-transfer", "specific-putaway", bulkTestUUID(t), inventorycontrol.BulkTransferPostingCommand{RequestID: "specific-receive-request", Quantity: "1", PostingDate: time.Date(2030, 2, 8, 0, 0, 0, 0, time.UTC)})
	if err != nil || specificReceived.Status != "received" || specificReceived.CostAmount != "17.0000" {
		t.Fatalf("specific receipt=%+v err=%v", specificReceived, err)
	}
	if _, err = repo.RegisterWarehouseActivity(ctx, tenant, "destination", specificReceived.PutAwayID, 1, time.Date(2030, 2, 9, 0, 0, 0, 0, time.UTC), bulkTestUUID(t)); err != nil {
		t.Fatal(err)
	}
	var firstRemaining, secondRemaining, transferredCost string
	if err = pool.QueryRow(ctx, `select (select remaining_quantity::text from inventory.bulk_cost_layer where tenant_id=$1 and receipt_entry_id=$2),(select remaining_quantity::text from inventory.bulk_cost_layer where tenant_id=$1 and receipt_entry_id=$3),(select sum(e.cost_amount)::text from inventory.bulk_transfer_receipt_component rc join inventory.bulk_inventory_entry e on e.tenant_id=rc.tenant_id and e.entry_id=rc.receipt_entry_id where rc.tenant_id=$1 and rc.receipt_id=$4)`, tenant, firstSpecific.EntryID, secondSpecific.EntryID, specificReceived.PostingID).Scan(&firstRemaining, &secondRemaining, &transferredCost); err != nil || firstRemaining != "1.000000" || secondRemaining != "0.000000" || transferredCost != "17.0000" {
		t.Fatalf("specific cost evidence first=%s second=%s transferred=%s err=%v", firstRemaining, secondRemaining, transferredCost, err)
	}
}
