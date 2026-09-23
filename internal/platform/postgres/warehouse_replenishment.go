package postgres

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"elite.local/enterprise/internal/inventorycontrol"
	"github.com/jackc/pgx/v5"
)

func (r *InventoryControl) CreateWarehouseReplenishment(ctx context.Context, tenant string, ids inventorycontrol.WarehouseIDs, command inventorycontrol.WarehouseReplenishmentCommand) (inventorycontrol.WarehouseActivity, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return inventorycontrol.WarehouseActivity{}, err
	}
	defer tx.Rollback(ctx)

	var replayActivity, replayOrganization, replayBin, replayItem, replayAssigned, replayTargetUOM, replayBaseUOM string
	var replayFEFO, replayBreakbulk bool
	err = tx.QueryRow(ctx, `select r.activity_id,r.organization_id,r.to_bin_id,r.item_id,r.target_uom,r.allow_breakbulk,r.use_fefo,a.assigned_to,(select u.uom_code from inventory.item_unit_of_measure u where u.tenant_id=r.tenant_id and u.item_id=r.item_id and u.is_base) from inventory.warehouse_replenishment_request r join inventory.warehouse_activity a using(tenant_id,activity_id) where r.tenant_id=$1 and r.request_id=$2`, tenant, command.RequestID).Scan(&replayActivity, &replayOrganization, &replayBin, &replayItem, &replayTargetUOM, &replayBreakbulk, &replayFEFO, &replayAssigned, &replayBaseUOM)
	if err == nil {
		requestedTarget := command.TargetUOM
		if requestedTarget == "" {
			requestedTarget = replayBaseUOM
		}
		if replayOrganization != command.OrganizationID || replayBin != command.ToBinID || replayItem != command.ItemID || replayTargetUOM != requestedTarget || replayBreakbulk != command.AllowBreakbulk || replayFEFO != command.UseFEFO || replayAssigned != command.AssignedTo {
			return inventorycontrol.WarehouseActivity{}, inventorycontrol.ErrConflict
		}
		value, readErr := readWarehouseActivity(ctx, tx, tenant, command.OrganizationID, replayActivity)
		if readErr != nil {
			return value, readErr
		}
		return value, tx.Commit(ctx)
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.WarehouseActivity{}, err
	}
	if err = warehouseLock(ctx, tx, tenant, command.OrganizationID, command.ItemID); err != nil {
		return inventorycontrol.WarehouseActivity{}, err
	}
	targetUOM, err := resolveWarehouseUOM(ctx, tx, tenant, command.ItemID, command.TargetUOM)
	if err != nil {
		return inventorycontrol.WarehouseActivity{}, err
	}

	var minimumText, maximumText, currentText, pendingText, targetType string
	var targetRank int
	var targetFixed, targetBlocked bool
	err = tx.QueryRow(ctx, `select p.min_quantity::text,coalesce(p.max_quantity::text,''),w.bin_rank,w.bin_type,p.fixed,w.movement_blocked,coalesce((select sum(b.quantity)::text from inventory.bulk_balance b where b.tenant_id=p.tenant_id and b.organization_id=p.organization_id and b.bin_id=p.bin_id and b.item_id=p.item_id),'0'),coalesce((select sum(al.quantity)::text from inventory.warehouse_activity a join inventory.warehouse_activity_line al using(tenant_id,activity_id) where a.tenant_id=p.tenant_id and a.organization_id=p.organization_id and a.activity_type='movement' and a.status='open' and al.to_bin_id=p.bin_id and al.item_id=p.item_id),'0') from inventory.item_bin_policy p join inventory.warehouse_bin w using(tenant_id,organization_id,bin_id) where p.tenant_id=$1 and p.organization_id=$2 and p.bin_id=$3 and p.item_id=$4 for share of p,w`, tenant, command.OrganizationID, command.ToBinID, command.ItemID).Scan(&minimumText, &maximumText, &targetRank, &targetType, &targetFixed, &targetBlocked, &currentText, &pendingText)
	if errors.Is(err, pgx.ErrNoRows) || !targetFixed || targetBlocked || (targetType != "pick" && targetType != "putpick") || maximumText == "" {
		return inventorycontrol.WarehouseActivity{}, inventorycontrol.ErrConflict
	}
	if err != nil {
		return inventorycontrol.WarehouseActivity{}, err
	}
	minimum, minimumOK := parseRat(minimumText)
	maximum, maximumOK := parseRat(maximumText)
	current, currentOK := parseRat(currentText)
	pending, pendingOK := parseRat(pendingText)
	if !minimumOK || !maximumOK || !currentOK || !pendingOK {
		return inventorycontrol.WarehouseActivity{}, inventorycontrol.ErrConflict
	}
	effective := new(big.Rat).Add(current, pending)
	if effective.Cmp(minimum) >= 0 {
		return inventorycontrol.WarehouseActivity{}, inventorycontrol.ErrConflict
	}
	remaining := new(big.Rat).Sub(maximum, effective)
	if remaining.Sign() <= 0 {
		return inventorycontrol.WarehouseActivity{}, inventorycontrol.ErrConflict
	}

	order := `w.bin_rank desc,w.bin_code,coalesce(l.expiration_date,'infinity'::date),coalesce(b.lot_id,'')`
	if command.UseFEFO {
		order = `coalesce(l.expiration_date,'infinity'::date),w.bin_rank desc,w.bin_code,coalesce(b.lot_id,'')`
	}
	query := `select b.balance_id,b.bin_id,coalesce(b.lot_id,''),(b.quantity-b.reserved_quantity)::text,l.expiration_date from inventory.bulk_balance b join inventory.warehouse_bin w using(tenant_id,organization_id,bin_id) join inventory.item_bin_policy p on p.tenant_id=b.tenant_id and p.organization_id=b.organization_id and p.bin_id=b.bin_id and p.item_id=b.item_id left join inventory.inventory_lot l on l.tenant_id=b.tenant_id and l.lot_id=b.lot_id where b.tenant_id=$1 and b.organization_id=$2 and b.item_id=$3 and b.bin_id<>$4 and w.bin_rank<$5 and w.bin_type not in ('receive','ship') and not w.movement_blocked and coalesce(l.blocked,false)=false and (l.expiration_date is null or l.expiration_date>=current_date) and b.quantity>b.reserved_quantity order by ` + order + ` for update of b`
	rows, err := tx.Query(ctx, query, tenant, command.OrganizationID, command.ItemID, command.ToBinID, targetRank)
	if err != nil {
		return inventorycontrol.WarehouseActivity{}, err
	}
	type candidate struct {
		balanceID           int64
		bin, lot, available string
		expiration          *time.Time
		plan                warehouseUOMPlan
	}
	rawCandidates := []candidate{}
	for rows.Next() {
		var candidate candidate
		if err = rows.Scan(&candidate.balanceID, &candidate.bin, &candidate.lot, &candidate.available, &candidate.expiration); err != nil {
			rows.Close()
			return inventorycontrol.WarehouseActivity{}, err
		}
		rawCandidates = append(rawCandidates, candidate)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return inventorycontrol.WarehouseActivity{}, err
	}
	candidates := []candidate{}
	for _, candidate := range rawCandidates {
		if remaining.Sign() <= 0 {
			break
		}
		available, ok := parseRat(candidate.available)
		if !ok || available.Sign() <= 0 {
			return inventorycontrol.WarehouseActivity{}, inventorycontrol.ErrConflict
		}
		limit := minRat(available, remaining)
		plans, leftover, planErr := planWarehousePackaging(ctx, tx, candidate.balanceID, targetUOM, command.AllowBreakbulk, limit)
		if planErr != nil {
			return inventorycontrol.WarehouseActivity{}, planErr
		}
		moved := new(big.Rat).Sub(limit, leftover)
		for _, plan := range plans {
			plannedCandidate := candidate
			plannedCandidate.plan = plan
			candidates = append(candidates, plannedCandidate)
		}
		remaining.Sub(remaining, moved)
	}
	if len(candidates) == 0 {
		return inventorycontrol.WarehouseActivity{}, inventorycontrol.ErrConflict
	}
	planned := new(big.Rat).Sub(new(big.Rat).Sub(maximum, effective), remaining)

	_, err = tx.Exec(ctx, `insert into inventory.warehouse_activity(tenant_id,activity_id,organization_id,activity_type,status,source_kind,source_id,source_line_id,request_id,assigned_to,version) values($1,$2,$3,'movement','open','bin-replenishment',$4,$5,$6,$7,1)`, tenant, ids.ActivityID, command.OrganizationID, command.ToBinID, command.ItemID, command.RequestID, command.AssignedTo)
	if err != nil {
		return inventorycontrol.WarehouseActivity{}, bulkConflict(err)
	}
	_, err = tx.Exec(ctx, `insert into inventory.warehouse_replenishment_request(tenant_id,request_id,activity_id,organization_id,to_bin_id,item_id,target_uom,allow_breakbulk,use_fefo,planned_quantity) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10::numeric)`, tenant, command.RequestID, ids.ActivityID, command.OrganizationID, command.ToBinID, command.ItemID, targetUOM.code, command.AllowBreakbulk, command.UseFEFO, formatRat(planned, 6))
	if err != nil {
		return inventorycontrol.WarehouseActivity{}, bulkConflict(err)
	}
	for index, candidate := range candidates {
		lineID := fmt.Sprintf("%s-%06d", ids.ActivityID, index+1)
		reservationID := fmt.Sprintf("%s-r-%06d", ids.ActivityID, index+1)
		_, err = tx.Exec(ctx, `insert into inventory.bulk_reservation(tenant_id,reservation_id,organization_id,bin_id,item_id,lot_id,demand_kind,demand_id,demand_line_id,quantity,status,cancellation_disallowed,version) values($1,$2,$3,$4,$5,nullif($6,''),'warehouse-replenishment',$7,$8,$9::numeric,'reservation',true,1)`, tenant, reservationID, command.OrganizationID, candidate.bin, command.ItemID, candidate.lot, ids.ActivityID, lineID, candidate.plan.movedBase)
		if err != nil {
			return inventorycontrol.WarehouseActivity{}, bulkConflict(err)
		}
		_, err = tx.Exec(ctx, `insert into inventory.warehouse_activity_line(tenant_id,activity_id,organization_id,line_id,sequence_no,from_bin_id,to_bin_id,item_id,lot_id,reservation_id,quantity,expiration_date,uom_code,uom_quantity,qty_per_uom,from_uom_code,from_uom_quantity,from_qty_per_uom) values($1,$2,$3,$4,$5,$6,$7,$8,nullif($9,''),$10,$11::numeric,$12,$13,$14::numeric,$15::numeric,$16,$17::numeric,$18::numeric)`, tenant, ids.ActivityID, command.OrganizationID, lineID, index+1, candidate.bin, command.ToBinID, command.ItemID, candidate.lot, reservationID, candidate.plan.movedBase, candidate.expiration, candidate.plan.targetUOM, candidate.plan.targetQuantity, candidate.plan.targetFactor, candidate.plan.sourceUOM, candidate.plan.sourceQuantity, candidate.plan.sourceFactor)
		if err != nil {
			return inventorycontrol.WarehouseActivity{}, bulkConflict(err)
		}
		line := inventorycontrol.WarehouseActivityLine{ID: lineID, Sequence: index + 1, FromBinID: candidate.bin, ToBinID: command.ToBinID, ItemID: command.ItemID, LotID: candidate.lot, ReservationID: reservationID, ReservationVersion: 1, Quantity: candidate.plan.movedBase, FromUOMCode: candidate.plan.sourceUOM, FromUOMQuantity: candidate.plan.sourceQuantity, FromQuantityPerUOM: candidate.plan.sourceFactor, UOMCode: candidate.plan.targetUOM, UOMQuantity: candidate.plan.targetQuantity, QuantityPerUOM: candidate.plan.targetFactor, ExpirationDate: candidate.expiration}
		if err = reserveWarehousePackaging(ctx, tx, tenant, ids.ActivityID, command.OrganizationID, candidate.plan, line); err != nil {
			return inventorycontrol.WarehouseActivity{}, err
		}
	}
	if err = recordBulkEvent(ctx, tx, tenant, ids.EventID, "warehouse-replenishment", ids.ActivityID, "warehouse-replenishment.created", 1, map[string]any{"organization_id": command.OrganizationID, "item_id": command.ItemID, "to_bin_id": command.ToBinID, "quantity": formatRat(planned, 6), "target_uom": targetUOM.code, "allow_breakbulk": command.AllowBreakbulk, "use_fefo": command.UseFEFO, "lines": len(candidates)}); err != nil {
		return inventorycontrol.WarehouseActivity{}, err
	}
	value, err := readWarehouseActivity(ctx, tx, tenant, command.OrganizationID, ids.ActivityID)
	if err != nil {
		return value, err
	}
	return value, tx.Commit(ctx)
}

func (r *InventoryControl) CancelWarehouseReplenishment(ctx context.Context, tenant, organization, activityID string, version int64, eventID string) (inventorycontrol.WarehouseActivity, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return inventorycontrol.WarehouseActivity{}, err
	}
	defer tx.Rollback(ctx)
	var itemID string
	err = tx.QueryRow(ctx, `select al.item_id from inventory.warehouse_activity a join inventory.warehouse_activity_line al using(tenant_id,activity_id) where a.tenant_id=$1 and a.organization_id=$2 and a.activity_id=$3 and a.activity_type='movement' and a.source_kind='bin-replenishment' and a.status='open' and a.version=$4 order by al.sequence_no limit 1 for update of a`, tenant, organization, activityID, version).Scan(&itemID)
	if errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.WarehouseActivity{}, inventorycontrol.ErrConflict
	}
	if err != nil {
		return inventorycontrol.WarehouseActivity{}, err
	}
	if err = warehouseLock(ctx, tx, tenant, organization, itemID); err != nil {
		return inventorycontrol.WarehouseActivity{}, err
	}
	activity, err := readWarehouseActivity(ctx, tx, tenant, organization, activityID)
	if err != nil {
		return activity, err
	}
	for _, line := range activity.Lines {
		if err = releaseWarehousePackagingReservation(ctx, tx, tenant, organization, activityID, line); err != nil {
			return activity, err
		}
		updated, err := tx.Exec(ctx, `update inventory.bulk_reservation set status='released',cancellation_disallowed=false,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and reservation_id=$3 and status='reservation' and version=1`, tenant, organization, line.ReservationID)
		if err != nil || updated.RowsAffected() != 1 {
			if err != nil {
				return activity, err
			}
			return activity, inventorycontrol.ErrConflict
		}
	}
	updated, err := tx.Exec(ctx, `update inventory.warehouse_activity set status='cancelled',version=version+1 where tenant_id=$1 and organization_id=$2 and activity_id=$3 and activity_type='movement' and status='open' and version=$4`, tenant, organization, activityID, version)
	if err != nil || updated.RowsAffected() != 1 {
		if err != nil {
			return activity, err
		}
		return activity, inventorycontrol.ErrConflict
	}
	activity.Status, activity.Version = "cancelled", version+1
	if err = recordBulkEvent(ctx, tx, tenant, eventID, "warehouse-replenishment", activityID, "warehouse-replenishment.cancelled", activity.Version, map[string]any{"organization_id": organization, "item_id": itemID, "lines": len(activity.Lines)}); err != nil {
		return activity, err
	}
	return activity, tx.Commit(ctx)
}
