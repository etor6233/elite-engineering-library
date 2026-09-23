package postgres

import (
	"context"
	"errors"
	"os"
	"testing"

	"elite.local/enterprise/internal/inventorycontrol"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestItemUnitOfMeasureExactConversionAndImmutability(t *testing.T) {
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
	if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'UOM V167','UOM V167')`, tenant, "uom-"+tenant[:8]); err != nil {
		t.Fatal(err)
	}
	defer func() {
		tx, cleanupErr := pool.Begin(ctx)
		if cleanupErr != nil {
			t.Errorf("UOM cleanup begin: %v", cleanupErr)
			return
		}
		defer tx.Rollback(ctx)
		for _, command := range []string{
			`alter table inventory.item_unit_of_measure disable trigger item_unit_of_measure_immutable`,
			`delete from inventory.item_unit_of_measure where tenant_id=$1`,
			`delete from platform.outbox_event where tenant_id=$1`,
			`delete from inventory.stock_item where tenant_id=$1`,
			`delete from platform.tenant where tenant_id=$1`,
			`alter table inventory.item_unit_of_measure enable trigger item_unit_of_measure_immutable`,
		} {
			args := []any{}
			if command[0:6] == "delete" {
				args = append(args, tenant)
			}
			if _, cleanupErr = tx.Exec(ctx, command, args...); cleanupErr != nil {
				t.Errorf("UOM cleanup: %v", cleanupErr)
				return
			}
		}
		if cleanupErr = tx.Commit(ctx); cleanupErr != nil {
			t.Errorf("UOM cleanup commit: %v", cleanupErr)
		}
	}()

	repo := NewInventoryControl(pool)
	item := inventorycontrol.BulkItem{ID: "part", Code: "UOM-PART", Description: "Discrete part", BaseUOM: "EA", BaseRoundingPrecision: "1", TrackingMode: "none", CostingMethod: "fifo", Version: 1}
	if _, err = repo.CreateBulkItem(ctx, tenant, bulkTestUUID(t), item); err != nil {
		t.Fatal(err)
	}
	box := inventorycontrol.ItemUnitOfMeasure{ItemID: "part", Code: "BOX", QuantityPerUnit: "12", RoundingPrecision: "1", Version: 1}
	if _, err = repo.ConfigureItemUnitOfMeasure(ctx, tenant, bulkTestUUID(t), box); err != nil {
		t.Fatal(err)
	}
	converted, err := repo.ConvertItemUnitOfMeasure(ctx, tenant, "part", "BOX", "2.5")
	if err != nil || converted.BaseCode != "EA" || converted.BaseQuantity != "30.000000" || converted.QuantityPerUnit != "12.000000" || converted.RoundingPrecision != "1.000000" {
		t.Fatalf("conversion=%+v err=%v", converted, err)
	}
	if _, err = repo.ConvertItemUnitOfMeasure(ctx, tenant, "part", "BOX", "0.1"); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("residual conversion error=%v", err)
	}
	bad := inventorycontrol.ItemUnitOfMeasure{ItemID: "part", Code: "INNER", QuantityPerUnit: "2.5", RoundingPrecision: "1", Version: 1}
	if _, err = repo.ConfigureItemUnitOfMeasure(ctx, tenant, bulkTestUUID(t), bad); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("misaligned factor error=%v", err)
	}
	if _, err = repo.ConfigureItemUnitOfMeasure(ctx, tenant, bulkTestUUID(t), box); !errors.Is(err, inventorycontrol.ErrConflict) {
		t.Fatalf("duplicate UOM error=%v", err)
	}
	_, err = pool.Exec(ctx, `update inventory.item_unit_of_measure set qty_per_uom=24 where tenant_id=$1 and item_id='part' and uom_code='BOX'`, tenant)
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != "55000" {
		t.Fatalf("immutable factor error=%v", err)
	}
}
