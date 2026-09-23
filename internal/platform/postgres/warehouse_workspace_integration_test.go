package postgres

import (
	"context"
"errors"
	"elite.local/enterprise/internal/inventorycontrol"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"testing"
	"time"
)

func TestWarehouseWorkspace(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("owned PostgreSQL required")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := bulkTestUUID(t)
	if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Workspace fixture','Workspace fixture')`, tenant, "workspace-"+tenant[:8]); err != nil {
		t.Fatal(err)
	}
	for _, org := range []string{"warehouse-a", "warehouse-b", "hidden"} {
		if _, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,$2,$2,$2,'warehouse')`, tenant, org); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewInventoryControl(pool)
	scope := inventorycontrol.WarehouseWorkspaceScope{Tenant: tenant, Organization: "warehouse-a", Organizations: []string{"warehouse-a", "warehouse-b"}}
	read := func(dataset string, q inventorycontrol.WarehouseWorkspaceQuery) inventorycontrol.WarehouseWorkspacePage {
		t.Helper()
		q.Dataset = dataset
		q.Organization = scope.Organization
		if q.Limit == 0 {
			q.Limit = 50
		}
		p, e := repo.WarehouseWorkspace(ctx, scope, q)
		if e != nil {
			t.Fatalf("read %s: %v", dataset, e)
		}
		return p
	}
	for dataset := range inventorycontrol.WarehouseDatasets {
		t.Run("query_contract_"+dataset, func(t *testing.T) { read(dataset, inventorycontrol.WarehouseWorkspaceQuery{ParentID: "absent"}) })
	}
	for _, id := range []string{"part-a", "part-b", "part-c"} {
		_, e := repo.CreateBulkItem(ctx, tenant, bulkTestUUID(t), inventorycontrol.BulkItem{ID: id, Code: id, Description: "Repuesto de prueba", BaseUOM: "EA", BaseRoundingPrecision: "0.000001", TrackingMode: "lot", CostingMethod: "fifo", Version: 1})
		if e != nil {
			t.Fatal(e)
		}
	}
	t.Run("stable_cursor_scope", func(t *testing.T) {
		p := read("items", inventorycontrol.WarehouseWorkspaceQuery{Limit: 2})
		if len(p.Rows) != 2 || p.NextCursor == nil || p.Rows[0].ID != "part-a" {
			t.Fatal(p)
		}
		q := inventorycontrol.WarehouseWorkspaceQuery{Dataset: "items", Organization: "warehouse-a", Cursor: *p.NextCursor, Limit: 2}
		next := read("items", q)
		if len(next.Rows) != 1 || next.Rows[0].ID != "part-c" || next.NextCursor != nil {
			t.Fatal(next)
		}
		q.Search = "changed"
		if _, e := repo.WarehouseWorkspace(ctx, scope, q); e == nil {
			t.Fatal("filter-swapped cursor accepted")
		}
		q.Search = ""
		wrong := scope
		wrong.Tenant = bulkTestUUID(t)
		if _, e := repo.WarehouseWorkspace(ctx, wrong, q); e == nil {
			t.Fatal("tenant-swapped cursor accepted")
		}
	})
	for _, org := range []string{"warehouse-a", "warehouse-b", "hidden"} {
		for _, bin := range []struct {
			id, kind string
			rank     int
		}{{"receive", "receive", 0}, {"pick", "putpick", 100}, {"ship", "ship", 0}} {
			_, e := repo.CreateWarehouseBin(ctx, tenant, bulkTestUUID(t), inventorycontrol.WarehouseBin{ID: func()string{if org=="warehouse-a"{return bin.id};return org+"-"+bin.id}(), OrganizationID: org, Code: bin.id, Type: bin.kind, Ranking: bin.rank, Version: 1})
			if e != nil {
				t.Fatal(e)
			}
			_, e = repo.ConfigureItemBin(ctx, tenant, bulkTestUUID(t), inventorycontrol.ItemBinPolicy{OrganizationID: org, ItemID: "part-a", BinID: func()string{if org=="warehouse-a"{return bin.id};return org+"-"+bin.id}(), Fixed: true, Default: bin.id == "pick", MinQuantity: "0", MaxQuantity: "100", Version: 1})
			if e != nil {
				t.Fatal(e)
			}
		}
	}
	exp := time.Date(2035, 1, 1, 0, 0, 0, 0, time.UTC)
	posting := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	receipt, e := repo.PostWarehouseReceipt(ctx, tenant, warehouseTestIDs(t, "workspace"), inventorycontrol.WarehouseReceiptCommand{RequestID: "receive-exact-request", OrganizationID: "warehouse-a", ReceiveBinID: "receive", ItemID: "part-a", LotNo: "LOT-A", ExpirationDate: &exp, Quantity: "10", UnitCost: "2.5", PostingDate: posting, SourceKind: "purchase", SourceID: "receipt-reference"})
	if e != nil {
		t.Fatal(e)
	}
	t.Run("receipt_state_recovery_no_repeat", func(t *testing.T) {
		p := read("activities", inventorycontrol.WarehouseWorkspaceQuery{RequestID: "receive-exact-request"})
		if len(p.Rows) != 1 || p.Rows[0].State != "open" {
			t.Fatal(p)
		}
		d := read("activity-lines", inventorycontrol.WarehouseWorkspaceQuery{ParentID: receipt.PutAway.ID})
		if len(d.Rows) != 1 {
			t.Fatal(d)
		}
		var row map[string]any
		if json.Unmarshal(d.Rows[0].Data, &row) != nil || row["quantity"] != "10.000000" {
			t.Fatal(string(d.Rows[0].Data))
		}
	})
	if _, e = repo.RegisterWarehouseActivity(ctx, tenant, "warehouse-a", receipt.PutAway.ID, 1, posting, bulkTestUUID(t)); e != nil {
		t.Fatal(e)
	}
	t.Run("registered_recovery", func(t *testing.T) {
		p := read("activities", inventorycontrol.WarehouseWorkspaceQuery{ID: receipt.PutAway.ID})
		if len(p.Rows) != 1 || p.Rows[0].State != "registered" || p.Rows[0].Version != "2" {
			t.Fatal(p)
		}
	})
	t.Run("positive_records_and_exact_decimals", func(t *testing.T) {
		for _, ds := range []string{"organizations", "items", "bins", "policies", "uoms", "balances", "activities", "entries"} {
			p := read(ds, inventorycontrol.WarehouseWorkspaceQuery{})
			if len(p.Rows) == 0 {
				t.Fatal(ds)
			}
		}
		p := read("balances", inventorycontrol.WarehouseWorkspaceQuery{ItemID: "part-a"})
		found := false
		for _, v := range p.Rows {
			var d map[string]any
			json.Unmarshal(v.Data, &d)
			if d["quantity"] == "10.000000" {
				found = true
			}
		}
		if !found {
			t.Fatal(p)
		}
	})
	r, e := repo.ReserveBulk(ctx, tenant, bulkTestUUID(t), inventorycontrol.BulkReservation{ID: "manual-reservation", OrganizationID: "warehouse-a", BinID: "pick", ItemID: "part-a", LotID: receipt.LotID, DemandKind: "manual", DemandID: "reservation-request", DemandLineID: "line", Quantity: "2", Status: "reservation", Version: 1})
	if e != nil {
		t.Fatal(e)
	}
	t.Run("reservation_release_recovery", func(t *testing.T) {
		p := read("reservations", inventorycontrol.WarehouseWorkspaceQuery{RequestID: "reservation-request"})
		if len(p.Rows) != 1 || p.Rows[0].ID != r.ID {
			t.Fatal(p)
		}
		if e := repo.ReleaseBulk(ctx, tenant, "warehouse-a", r.ID, 1, bulkTestUUID(t)); e != nil {
			t.Fatal(e)
		}
		p = read("reservations", inventorycontrol.WarehouseWorkspaceQuery{ID: r.ID})
		if p.Rows[0].State != "released" || p.Rows[0].Version != "2" {
			t.Fatal(p)
		}
	})
	for _, to := range []string{"warehouse-b", "hidden"} {
		_, e := repo.CreateBulkTransfer(ctx, tenant, "transfer-"+to, "line-"+to, bulkTestUUID(t), inventorycontrol.BulkTransferCommand{RequestID: "request-" + to, FromOrganizationID: "warehouse-a", ToOrganizationID: to, ReceiveBinID: to+"-receive", InTransitCode: "ROAD", ItemID: "part-a", Quantity: "1", PostingDate: posting, SourceKind: "transfer", SourceID: "reference-" + to})
		if e != nil {
			t.Fatal(e)
		}
	}

	t.Run("durable_reservation_terminal_replay",func(t *testing.T){
	 c:=inventorycontrol.WarehouseReservationRequest{RequestID:"durable-reserve",OrganizationID:"warehouse-a",BinID:"pick",ItemID:"part-a",LotID:receipt.LotID,DemandKind:"manual",DemandID:"durable-demand",DemandLineID:"line",Quantity:"1"}
	 first,e:=repo.RequestWarehouseReservation(ctx,tenant,"actor-a",c);if e!=nil{t.Fatal(e)}
	 _,e=repo.IssueBulk(ctx,tenant,bulkTestUUID(t),"issue-durable-reservation",inventorycontrol.BulkIssue{OrganizationID:"warehouse-a",ReservationID:first.Reservation.ID,ReservationVersion:1,PostingDate:posting,SourceKind:"service",SourceID:"durable-consume"});if e!=nil{t.Fatal(e)}
	 if os.Getenv("WAREHOUSE_LEGACY_COUNTEREXAMPLE")=="true" {v:=first.Reservation;v.ID="late-legacy-duplicate";v.Status="reservation";v.Version=1;_,e:=repo.ReserveBulk(ctx,tenant,bulkTestUUID(t),v);if e==nil{t.Fatal("LEGACY_REPLAY_RECREATES_CONSUMED_RESERVATION: old owner has no durable request fence")};return}
	 replay,e:=repo.RequestWarehouseReservation(ctx,tenant,"actor-a",c);if e!=nil||replay.Reservation.ID!=first.Reservation.ID||replay.Reservation.Status!="consumed"{t.Fatalf("terminal replay=%+v %v",replay,e)}
	 recovered,e:=NewInventoryControl(pool).ReadWarehouseReservationRequest(ctx,tenant,"warehouse-a","actor-a",c.RequestID);if e!=nil||recovered.Reservation.ID!=first.Reservation.ID||recovered.Reservation.Status!="consumed"{t.Fatal(recovered,e)}
	 changed:=c;changed.Quantity="2";if _,e=repo.RequestWarehouseReservation(ctx,tenant,"actor-a",changed);!errors.Is(e,inventorycontrol.ErrConflict){t.Fatal("changed replay",e)}
	 if _,e=repo.RequestWarehouseReservation(ctx,tenant,"actor-b",c);!errors.Is(e,inventorycontrol.ErrConflict){t.Fatal("actor replay",e)}
	 if _,e=repo.ReadWarehouseReservationRequest(ctx,tenant,"warehouse-a","actor-b",c.RequestID);!errors.Is(e,inventorycontrol.ErrWarehouseRequestNotFound){t.Fatal("actor leak",e)}
	 if _,e=repo.ReadWarehouseReservationRequest(ctx,tenant,"warehouse-b","actor-a",c.RequestID);!errors.Is(e,inventorycontrol.ErrWarehouseRequestNotFound){t.Fatal("org leak",e)}
	 if _,e=repo.ReadWarehouseReservationRequest(ctx,bulkTestUUID(t),"warehouse-a","actor-a",c.RequestID);!errors.Is(e,inventorycontrol.ErrWarehouseRequestNotFound){t.Fatal("tenant leak",e)}
	 var count int;if e=pool.QueryRow(ctx,`select count(*) from inventory.bulk_reservation where tenant_id=$1 and demand_id='durable-demand'`,tenant).Scan(&count);e!=nil||count!=1{t.Fatal("duplicate",count,e)}
	 _,e=pool.Exec(ctx,`update inventory.workspace_reservation_request set actor_id='changed' where tenant_id=$1`,tenant);if e==nil{t.Fatal("journal mutable")}
	 // A journal failure must roll back both the business quantity and outbox.
	 rollback:=c;rollback.RequestID="rollback-journal";rollback.DemandID="rollback-demand"
	 _,e=pool.Exec(ctx,`create function inventory.workspace_test_reject() returns trigger language plpgsql as $$begin if new.request_id='rollback-journal' then raise exception 'fixture rejection';end if;return new;end$$;create trigger workspace_test_reject before insert on inventory.workspace_reservation_request for each row execute function inventory.workspace_test_reject()`);if e!=nil{t.Fatal(e)}
	 if _,e=repo.RequestWarehouseReservation(ctx,tenant,"actor-a",rollback);e==nil{t.Fatal("journal failure ignored")}
	 if e=pool.QueryRow(ctx,`select count(*) from inventory.bulk_reservation where tenant_id=$1 and demand_id='rollback-demand'`,tenant).Scan(&count);e!=nil||count!=0{t.Fatal("partial reservation committed",count,e)}
	 if _,e=pool.Exec(ctx,`drop trigger workspace_test_reject on inventory.workspace_reservation_request;drop function inventory.workspace_test_reject()`);e!=nil{t.Fatal(e)}
	})
	t.Run("both_transfer_organizations_required", func(t *testing.T) {
		p := read("transfers", inventorycontrol.WarehouseWorkspaceQuery{})
		if len(p.Rows) != 1 || p.Rows[0].ID != "transfer-warehouse-b" {
			t.Fatal(p)
		}
		one := scope
		one.Organizations = []string{"warehouse-a"}
		p, e := repo.WarehouseWorkspace(ctx, one, inventorycontrol.WarehouseWorkspaceQuery{Dataset: "transfers", Organization: "warehouse-a", Limit: 25})
		if e != nil || len(p.Rows) != 0 {
			t.Fatal(p, e)
		}
	})
	t.Run("cross_tenant_and_organization_reads", func(t *testing.T) {
		wrong := scope
		wrong.Tenant = bulkTestUUID(t)
		p, e := repo.WarehouseWorkspace(ctx, wrong, inventorycontrol.WarehouseWorkspaceQuery{Dataset: "items", Organization: "warehouse-a", Limit: 25})
		if e != nil || len(p.Rows) != 0 {
			t.Fatal(p, e)
		}
		wrong = scope
		wrong.Organization = "warehouse-b"
		p, e = repo.WarehouseWorkspace(ctx, wrong, inventorycontrol.WarehouseWorkspaceQuery{Dataset: "activities", Organization: "warehouse-b", ID: receipt.PutAway.ID, Limit: 25})
		if e != nil || len(p.Rows) != 0 {
			t.Fatal(p, e)
		}
	})
	t.Run("database_failure_not_empty_pass", func(t *testing.T) {
		closed, e := pgxpool.New(ctx, dsn)
		if e != nil {
			t.Fatal(e)
		}
		closed.Close()
		_, e = NewInventoryControl(closed).WarehouseWorkspace(ctx, scope, inventorycontrol.WarehouseWorkspaceQuery{Dataset: "items", Organization: "warehouse-a", Limit: 25})
		if e == nil {
			t.Fatal("database failure hidden")
		}
	})
}
