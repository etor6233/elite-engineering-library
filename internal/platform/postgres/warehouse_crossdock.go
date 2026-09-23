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

func (r *InventoryControl) ConfigureWarehouseCrossDock(ctx context.Context, tenant, eventID string, value inventorycontrol.WarehouseCrossDockPolicy) (inventorycontrol.WarehouseCrossDockPolicy, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return value, err
	}
	defer tx.Rollback(ctx)
	var crossDock, blocked, fixed bool
	var binType string
	err = tx.QueryRow(ctx, `select w.cross_dock,w.movement_blocked,w.bin_type,p.fixed from inventory.warehouse_bin w join inventory.item_bin_policy p using(tenant_id,organization_id,bin_id) where w.tenant_id=$1 and w.organization_id=$2 and w.bin_id=$3 and p.item_id=$4 for share of w,p`, tenant, value.OrganizationID, value.BinID, value.ItemID).Scan(&crossDock, &blocked, &binType, &fixed)
	if errors.Is(err, pgx.ErrNoRows) || !crossDock || blocked || !fixed || (binType != "pick" && binType != "putpick") {
		return value, inventorycontrol.ErrConflict
	}
	if err != nil {
		return value, err
	}
	_, err = tx.Exec(ctx, `insert into inventory.warehouse_crossdock_policy(tenant_id,organization_id,item_id,bin_id,due_date_days,enabled,version) values($1,$2,$3,$4,$5,$6,1)`, tenant, value.OrganizationID, value.ItemID, value.BinID, value.DueDateDays, value.Enabled)
	if err != nil {
		return value, bulkConflict(err)
	}
	aggregate := value.OrganizationID + ":" + value.ItemID
	if err = recordBulkEvent(ctx, tx, tenant, eventID, "warehouse-crossdock-policy", aggregate, "warehouse-crossdock-policy.created", 1, map[string]any{"organization_id": value.OrganizationID, "item_id": value.ItemID, "bin_id": value.BinID, "due_date_days": value.DueDateDays, "enabled": value.Enabled}); err != nil {
		return value, err
	}
	return value, tx.Commit(ctx)
}

func planWarehouseCrossDock(ctx context.Context, tx pgx.Tx, tenant string, ids inventorycontrol.WarehouseIDs, value inventorycontrol.WarehouseReceiptCommand, lotID string) ([]warehousePlacement, *big.Rat, error) {
	receiptQuantity, ok := parseRat(value.Quantity)
	if !ok {
		return nil, nil, inventorycontrol.ErrConflict
	}
	var binID, maximumText, onHandText, plannedText string
	var dueDays int
	err := tx.QueryRow(ctx, `select p.bin_id,p.due_date_days,coalesce(bp.max_quantity::text,''),coalesce((select sum(b.quantity)::text from inventory.bulk_balance b where b.tenant_id=p.tenant_id and b.organization_id=p.organization_id and b.bin_id=p.bin_id and b.item_id=p.item_id),'0'),coalesce((select sum(al.quantity)::text from inventory.warehouse_activity a join inventory.warehouse_activity_line al using(tenant_id,activity_id) where a.tenant_id=p.tenant_id and a.organization_id=p.organization_id and a.activity_type='put-away' and a.status='open' and al.to_bin_id=p.bin_id and al.item_id=p.item_id),'0') from inventory.warehouse_crossdock_policy p join inventory.warehouse_bin w using(tenant_id,organization_id,bin_id) join inventory.item_bin_policy bp on bp.tenant_id=p.tenant_id and bp.organization_id=p.organization_id and bp.item_id=p.item_id and bp.bin_id=p.bin_id where p.tenant_id=$1 and p.organization_id=$2 and p.item_id=$3 and p.enabled and w.cross_dock and not w.movement_blocked and w.bin_type in ('pick','putpick') and bp.fixed for share of p,w,bp`, tenant, value.OrganizationID, value.ItemID).Scan(&binID, &dueDays, &maximumText, &onHandText, &plannedText)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, receiptQuantity, nil
	}
	if err != nil {
		return nil, nil, err
	}
	capacity, valid := decimalCapacity(maximumText, onHandText, plannedText)
	if !valid {
		return nil, nil, inventorycontrol.ErrConflict
	}
	crossDockBudget := new(big.Rat).Set(receiptQuantity)
	if capacity != nil {
		crossDockBudget = minRat(crossDockBudget, capacity)
	}
	if crossDockBudget.Sign() <= 0 {
		return nil, receiptQuantity, nil
	}
	cutoff := value.PostingDate.UTC().AddDate(0, 0, dueDays)
	rows, err := tx.Query(ctx, `select t.transfer_id,l.line_id,t.posting_date,(l.quantity-l.shipped_quantity-coalesce(x.quantity,0)-coalesce(k.quantity,0))::text from inventory.bulk_transfer t join inventory.bulk_transfer_line l using(tenant_id,transfer_id) left join lateral (select sum(a.quantity) quantity from inventory.warehouse_crossdock_allocation a where a.tenant_id=t.tenant_id and a.transfer_id=t.transfer_id and a.transfer_line_id=l.line_id) x on true left join lateral (select sum(al.quantity) quantity from inventory.warehouse_activity wa join inventory.warehouse_activity_line al using(tenant_id,activity_id) left join inventory.warehouse_crossdock_pick_link pl on pl.tenant_id=al.tenant_id and pl.pick_activity_id=al.activity_id and pl.pick_line_id=al.line_id where wa.tenant_id=t.tenant_id and wa.organization_id=t.from_organization_id and wa.activity_type='pick' and wa.source_kind='transfer-outbound' and wa.source_id=t.transfer_id and wa.source_line_id=l.line_id and pl.pick_activity_id is null and (wa.status='open' or (wa.status='registered' and not exists(select 1 from inventory.bulk_transfer_shipment s where s.tenant_id=wa.tenant_id and s.warehouse_activity_id=wa.activity_id)))) k on true where t.tenant_id=$1 and t.from_organization_id=$2 and l.item_id=$3 and t.status='released' and t.posting_date<=$4::date and l.specific_receipt_entry_id is null and l.quantity-l.shipped_quantity-coalesce(x.quantity,0)-coalesce(k.quantity,0)>0 order by t.posting_date,t.created_at,t.transfer_id for update of t,l`, tenant, value.OrganizationID, value.ItemID, cutoff)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	placements := []warehousePlacement{}
	allocated := new(big.Rat)
	for rows.Next() && crossDockBudget.Sign() > 0 {
		var transferID, lineID, neededText string
		var dueDate time.Time
		if err = rows.Scan(&transferID, &lineID, &dueDate, &neededText); err != nil {
			return nil, nil, err
		}
		needed, valid := parseRat(neededText)
		if !valid || needed.Sign() <= 0 {
			return nil, nil, inventorycontrol.ErrConflict
		}
		take := minRat(needed, crossDockBudget)
		index := len(placements) + 1
		placements = append(placements, warehousePlacement{binID: binID, quantity: formatRat(take, 6), crossDockAllocationID: fmt.Sprintf("%s-crossdock-%06d", ids.ActivityID, index), transferID: transferID, transferLineID: lineID, dueDate: dueDate, lotID: lotID})
		crossDockBudget.Sub(crossDockBudget, take)
		allocated.Add(allocated, take)
	}
	if err = rows.Err(); err != nil {
		return nil, nil, err
	}
	if capacity != nil && allocated.Cmp(capacity) > 0 {
		return nil, nil, inventorycontrol.ErrConflict
	}
	remaining := new(big.Rat).Sub(receiptQuantity, allocated)
	return placements, remaining, nil
}
