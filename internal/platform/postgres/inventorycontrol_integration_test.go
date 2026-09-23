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

func TestInventoryReservationATPTransferAndConcurrency(t *testing.T) {
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
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c29901"
	cleanup := func() {
		for _, q := range []string{`delete from platform.outbox_event where tenant_id=$1`, `delete from inventory.serial_transfer_unit where tenant_id=$1`, `delete from inventory.serial_transfer where tenant_id=$1`, `delete from inventory.serial_reservation where tenant_id=$1`, `delete from inventory.stock_unit where tenant_id=$1`, `delete from catalog.vehicle_variant where tenant_id=$1`, `delete from catalog.vehicle_model where tenant_id=$1`, `delete from org.organization where tenant_id=$1`, `delete from platform.tenant where tenant_id=$1`} {
			_, _ = pool.Exec(ctx, q, tenant)
		}
	}
	cleanup()
	defer cleanup()
	fixtures := []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'inventory-v159','Inventory','Inventory')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'a','a','A','warehouse'),($1,'b','b','B','franchisee')`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state) values($1,'model','model','Model','bicycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state) values($1,'variant','model','variant','Variant','{}','active')`,
		`insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,vin,battery_serial_number,state,version,received_at) values($1,'s1','a','variant','SERIAL-1','VIN-1','BATTERY-1','available',1,clock_timestamp()),($1,'s2','a','variant','SERIAL-2','VIN-2','BATTERY-2','available',1,clock_timestamp()),($1,'s3','a','variant','SERIAL-3','VIN-3','BATTERY-3','available',1,clock_timestamp())`,
	}
	for _, q := range fixtures {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewInventoryControl(pool)
	horizon := time.Now().Add(24 * time.Hour).UTC()
	atp, err := repo.AvailableToPromise(ctx, tenant, "a", "variant", horizon)
	if err != nil || atp.AvailableInventory != 3 || atp.AvailableToPromise != 3 {
		t.Fatalf("initial ATP=%+v err=%v", atp, err)
	}
	values := []inventorycontrol.Reservation{
		{ID: "r1", OrganizationID: "a", StockUnitID: "s1", VariantID: "variant", DemandKind: "service", DemandID: "service-1", Status: "reservation", Version: 1},
		{ID: "r2", OrganizationID: "a", StockUnitID: "s1", VariantID: "variant", DemandKind: "service", DemandID: "service-2", Status: "reservation", Version: 1},
	}
	errs := make([]error, 2)
	var wg sync.WaitGroup
	for index := range values {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, errs[i] = repo.Reserve(ctx, tenant, "018f4d4a-7b36-7a21-8d10-2f4c54c2991"+string(rune('0'+i)), values[i].ID, values[i], 1)
		}(index)
	}
	wg.Wait()
	winner, success := "", 0
	for i, callErr := range errs {
		if callErr == nil {
			winner, success = values[i].ID, success+1
		} else if !errors.Is(callErr, inventorycontrol.ErrConflict) {
			t.Fatalf("unexpected concurrent error: %v", callErr)
		}
	}
	if success != 1 {
		t.Fatalf("concurrent reservation successes=%d errors=%v", success, errs)
	}
	if err = repo.Release(ctx, tenant, "a", winner, 1, "018f4d4a-7b36-7a21-8d10-2f4c54c29912"); err != nil {
		t.Fatal(err)
	}
	transfer := inventorycontrol.Transfer{ID: "t1", FromOrganizationID: "a", ToOrganizationID: "b", StockUnitIDs: []string{"s2"}, ExpectedReceiptAt: time.Now().Add(time.Hour).UTC(), State: "draft", Version: 1}
	if _, err = repo.CreateTransfer(ctx, tenant, transfer.ID, transfer, map[string]int64{"s2": 1}, "018f4d4a-7b36-7a21-8d10-2f4c54c29913"); err != nil {
		t.Fatal(err)
	}
	if err = repo.TransitionTransfer(ctx, tenant, "a", "t1", "draft", 1, "released", "018f4d4a-7b36-7a21-8d10-2f4c54c29914"); err != nil {
		t.Fatal(err)
	}
	if err = repo.TransitionTransfer(ctx, tenant, "a", "t1", "released", 2, "in-transit", "018f4d4a-7b36-7a21-8d10-2f4c54c29915"); err != nil {
		t.Fatal(err)
	}
	atp, err = repo.AvailableToPromise(ctx, tenant, "b", "variant", horizon)
	if err != nil || atp.ScheduledReceipt != 1 || atp.AvailableToPromise != 1 {
		t.Fatalf("inbound ATP=%+v err=%v", atp, err)
	}
	if err = repo.TransitionTransfer(ctx, tenant, "a", "t1", "in-transit", 3, "received", "018f4d4a-7b36-7a21-8d10-2f4c54c29916"); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("wrong receiver accepted: %v", err)
	}
	if err = repo.TransitionTransfer(ctx, tenant, "b", "t1", "in-transit", 3, "received", "018f4d4a-7b36-7a21-8d10-2f4c54c29917"); err != nil {
		t.Fatal(err)
	}
	var stockVersion int64
	if err = pool.QueryRow(ctx, `select version from inventory.stock_unit where tenant_id=$1 and stock_unit_id='s2' and organization_id='b' and state='available'`, tenant).Scan(&stockVersion); err != nil || stockVersion != 4 {
		t.Fatalf("received stock version=%d err=%v", stockVersion, err)
	}
	returnTransfer := inventorycontrol.Transfer{ID: "t2", FromOrganizationID: "b", ToOrganizationID: "a", StockUnitIDs: []string{"s2"}, ExpectedReceiptAt: time.Now().Add(2 * time.Hour).UTC(), State: "draft", Version: 1}
	if _, err = repo.CreateTransfer(ctx, tenant, returnTransfer.ID, returnTransfer, map[string]int64{"s2": stockVersion}, "018f4d4a-7b36-7a21-8d10-2f4c54c29918"); err != nil {
		t.Fatalf("received unit was not transferable again: %v", err)
	}
	if err = repo.TransitionTransfer(ctx, tenant, "b", "t2", "draft", 1, "cancelled", "018f4d4a-7b36-7a21-8d10-2f4c54c29919"); err != nil {
		t.Fatal(err)
	}
}
