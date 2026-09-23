// AUTHORED durable request receipt. Existing reservation writer runs in the same transaction.
package postgres

import (
	"context"
	"elite.local/enterprise/internal/inventorycontrol"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func readWarehouseReservationRequest(ctx context.Context, tx pgx.Tx, tenant, org, actor, request string) (inventorycontrol.WarehouseReservationReceipt, error) {
	v := inventorycontrol.WarehouseReservationReceipt{Schema: "warehouse-reservation-receipt/v1", RequestID: request}
	r := &v.Reservation
	e := tx.QueryRow(ctx, `select j.payload_sha256,r.reservation_id,r.organization_id,r.bin_id,r.item_id,coalesce(r.lot_id,''),r.demand_kind,r.demand_id,r.demand_line_id,r.quantity::text,r.cancellation_disallowed,r.expires_at,r.status,r.version from inventory.workspace_reservation_request j join inventory.bulk_reservation r on r.tenant_id=j.tenant_id and r.organization_id=j.organization_id and r.reservation_id=j.reservation_id where j.tenant_id=$1 and j.organization_id=$2 and j.actor_id=$3 and j.request_id=$4`, tenant, org, actor, request).Scan(&v.PayloadSHA256, &r.ID, &r.OrganizationID, &r.BinID, &r.ItemID, &r.LotID, &r.DemandKind, &r.DemandID, &r.DemandLineID, &r.Quantity, &r.CancellationDisallowed, &r.ExpiresAt, &r.Status, &r.Version)
	return v, e
}
func (r *InventoryControl) ReadWarehouseReservationRequest(ctx context.Context, tenant, org, actor, request string) (inventorycontrol.WarehouseReservationReceipt, error) {
	tx, e := r.pool.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadOnly})
	if e != nil {
		return inventorycontrol.WarehouseReservationReceipt{}, e
	}
	defer tx.Rollback(ctx)
	v, e := readWarehouseReservationRequest(ctx, tx, tenant, org, actor, request)
	if errors.Is(e, pgx.ErrNoRows) {
		return v, inventorycontrol.ErrWarehouseRequestNotFound
	}
	if e != nil {
		return v, e
	}
	return v, tx.Commit(ctx)
}
func (r *InventoryControl) RequestWarehouseReservation(ctx context.Context, tenant, actor string, c inventorycontrol.WarehouseReservationRequest) (inventorycontrol.WarehouseReservationReceipt, error) {
	empty := inventorycontrol.WarehouseReservationReceipt{}
	if tenant == "" || actor == "" || c.Validate() != nil {
		return empty, inventorycontrol.ErrWarehouseWorkspaceQuery
	}
	tx, e := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if e != nil {
		return empty, e
	}
	defer tx.Rollback(ctx)
	// Locks the request namespace before the existing stock writer. Collisions only serialize.
	if _, e = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1,0))`, tenant+"|"+c.OrganizationID+"|warehouse-reserve|"+c.RequestID); e != nil {
		return empty, e
	}
	var savedHash, savedActor string
	e = tx.QueryRow(ctx, `select payload_sha256,actor_id from inventory.workspace_reservation_request where tenant_id=$1 and organization_id=$2 and request_id=$3`, tenant, c.OrganizationID, c.RequestID).Scan(&savedHash, &savedActor)
	if e == nil {
		if savedHash != c.Hash() || savedActor != actor {
			return empty, inventorycontrol.ErrConflict
		}
		v, e := readWarehouseReservationRequest(ctx, tx, tenant, c.OrganizationID, actor, c.RequestID)
		if e != nil {
			return empty, e
		}
		return v, tx.Commit(ctx)
	}
	if !errors.Is(e, pgx.ErrNoRows) {
		return empty, e
	}
	value := inventorycontrol.BulkReservation{ID: uuid.NewString(), OrganizationID: c.OrganizationID, BinID: c.BinID, ItemID: c.ItemID, LotID: c.LotID, DemandKind: c.DemandKind, DemandID: c.DemandID, DemandLineID: c.DemandLineID, Quantity: c.Quantity, CancellationDisallowed: c.CancellationDisallowed, ExpiresAt: c.ExpiresAt, Status: "reservation", Version: 1}
	if _, e = reserveBulkInTx(ctx, tx, tenant, uuid.NewString(), value); e != nil {
		return empty, e
	}
	if _, e = tx.Exec(ctx, `insert into inventory.workspace_reservation_request(tenant_id,organization_id,request_id,actor_id,payload_sha256,reservation_id)values($1,$2,$3,$4,$5,$6)`, tenant, c.OrganizationID, c.RequestID, actor, c.Hash(), value.ID); e != nil {
		return empty, bulkConflict(e)
	}
	v, e := readWarehouseReservationRequest(ctx, tx, tenant, c.OrganizationID, actor, c.RequestID)
	if e != nil {
		return empty, e
	}
	return v, tx.Commit(ctx)
}
