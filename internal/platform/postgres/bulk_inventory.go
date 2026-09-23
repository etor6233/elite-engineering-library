package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"strings"

	"elite.local/enterprise/internal/inventorycontrol"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func bulkConflict(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && (pgErr.Code == "23505" || pgErr.Code == "23514" || pgErr.Code == "23503" || pgErr.Code == "55000" || pgErr.Code == "40001") {
		return inventorycontrol.ErrConflict
	}
	return err
}

func recordBulkEvent(ctx context.Context, tx pgx.Tx, tenant, eventID, aggregateType, aggregateID, eventType string, version int64, payload any) error {
	encoded, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,$3,$4,$5,$6,1,clock_timestamp(),$7)`, tenant, eventID, aggregateType, aggregateID, version, eventType, encoded)
	return err
}

func (r *InventoryControl) CreateBulkItem(ctx context.Context, tenant, eventID string, value inventorycontrol.BulkItem) (inventorycontrol.BulkItem, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return value, err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `insert into inventory.stock_item(tenant_id,item_id,item_code,description,base_uom,tracking_mode,costing_method,version) values($1,$2,$3,$4,$5,$6,$7,1)`, tenant, value.ID, value.Code, value.Description, value.BaseUOM, value.TrackingMode, value.CostingMethod)
	if err != nil {
		return value, bulkConflict(err)
	}
	_, err = tx.Exec(ctx, `insert into inventory.item_unit_of_measure(tenant_id,item_id,uom_code,qty_per_uom,rounding_precision,is_base,version) values($1,$2,$3,1,$4::numeric,true,1)`, tenant, value.ID, value.BaseUOM, value.BaseRoundingPrecision)
	if err != nil {
		return value, bulkConflict(err)
	}
	if err = recordBulkEvent(ctx, tx, tenant, eventID, "bulk-item", value.ID, "bulk-item.created", 1, map[string]string{"item_code": value.Code, "base_uom": value.BaseUOM, "base_rounding_precision": value.BaseRoundingPrecision, "tracking_mode": value.TrackingMode, "costing_method": value.CostingMethod}); err != nil {
		return value, err
	}
	return value, tx.Commit(ctx)
}

func (r *InventoryControl) CreateWarehouseBin(ctx context.Context, tenant, eventID string, value inventorycontrol.WarehouseBin) (inventorycontrol.WarehouseBin, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return value, err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `insert into inventory.warehouse_bin(tenant_id,organization_id,bin_id,bin_code,bin_type,bin_rank,movement_blocked,cross_dock,version) values($1,$2,$3,$4,$5,$6,$7,$8,1)`, tenant, value.OrganizationID, value.ID, value.Code, value.Type, value.Ranking, value.MovementBlocked, value.CrossDock)
	if err != nil {
		return value, bulkConflict(err)
	}
	if err = recordBulkEvent(ctx, tx, tenant, eventID, "warehouse-bin", value.ID, "warehouse-bin.created", 1, map[string]any{"organization_id": value.OrganizationID, "bin_code": value.Code, "bin_type": value.Type, "bin_rank": value.Ranking, "cross_dock": value.CrossDock}); err != nil {
		return value, err
	}
	return value, tx.Commit(ctx)
}

func (r *InventoryControl) ConfigureItemBin(ctx context.Context, tenant, eventID string, value inventorycontrol.ItemBinPolicy) (inventorycontrol.ItemBinPolicy, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return value, err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `insert into inventory.item_bin_policy(tenant_id,organization_id,item_id,bin_id,fixed,dedicated,is_default,min_quantity,max_quantity,version) values($1,$2,$3,$4,$5,$6,$7,$8::numeric,nullif($9,'')::numeric,1)`, tenant, value.OrganizationID, value.ItemID, value.BinID, value.Fixed, value.Dedicated, value.Default, value.MinQuantity, value.MaxQuantity)
	if err != nil {
		return value, bulkConflict(err)
	}
	aggregate := value.OrganizationID + ":" + value.ItemID + ":" + value.BinID
	if err = recordBulkEvent(ctx, tx, tenant, eventID, "item-bin-policy", aggregate, "item-bin-policy.created", 1, map[string]any{"organization_id": value.OrganizationID, "item_id": value.ItemID, "bin_id": value.BinID, "fixed": value.Fixed, "dedicated": value.Dedicated, "default": value.Default}); err != nil {
		return value, err
	}
	return value, tx.Commit(ctx)
}

func (r *InventoryControl) ReceiveBulk(ctx context.Context, tenant, eventID, entryID, layerID, proposedLotID string, value inventorycontrol.BulkReceipt) (inventorycontrol.BulkReceiptResult, error) {
	result := inventorycontrol.BulkReceiptResult{EntryID: entryID, Quantity: value.Quantity, UnitCost: value.UnitCost}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return result, err
	}
	defer tx.Rollback(ctx)
	var tracking string
	var movementBlocked bool
	err = tx.QueryRow(ctx, `select i.tracking_mode,w.movement_blocked from inventory.stock_item i join inventory.item_bin_policy p on p.tenant_id=i.tenant_id and p.item_id=i.item_id join inventory.warehouse_bin w on w.tenant_id=p.tenant_id and w.organization_id=p.organization_id and w.bin_id=p.bin_id where i.tenant_id=$1 and i.item_id=$2 and p.organization_id=$3 and p.bin_id=$4 for share of i,w`, tenant, value.ItemID, value.OrganizationID, value.BinID).Scan(&tracking, &movementBlocked)
	if errors.Is(err, pgx.ErrNoRows) || movementBlocked {
		return result, inventorycontrol.ErrConflict
	}
	if err != nil {
		return result, err
	}
	lotID := ""
	if tracking == "lot" {
		if value.LotNo == "" {
			return result, inventorycontrol.ErrConflict
		}
		var existingExpiration, existingWarranty string
		err = tx.QueryRow(ctx, `select lot_id,coalesce(expiration_date::text,''),coalesce(warranty_date::text,'') from inventory.inventory_lot where tenant_id=$1 and item_id=$2 and lot_no=$3 for update`, tenant, value.ItemID, value.LotNo).Scan(&lotID, &existingExpiration, &existingWarranty)
		expiration, warranty := "", ""
		if value.ExpirationDate != nil {
			expiration = value.ExpirationDate.UTC().Format("2006-01-02")
		}
		if value.WarrantyDate != nil {
			warranty = value.WarrantyDate.UTC().Format("2006-01-02")
		}
		if errors.Is(err, pgx.ErrNoRows) {
			lotID = proposedLotID
			_, err = tx.Exec(ctx, `insert into inventory.inventory_lot(tenant_id,lot_id,item_id,lot_no,expiration_date,warranty_date) values($1,$2,$3,$4,nullif($5,'')::date,nullif($6,'')::date)`, tenant, lotID, value.ItemID, value.LotNo, expiration, warranty)
		} else if err == nil && (existingExpiration != expiration || existingWarranty != warranty) {
			return result, inventorycontrol.ErrConflict
		}
		if err != nil {
			return result, bulkConflict(err)
		}
	} else if value.LotNo != "" || value.ExpirationDate != nil || value.WarrantyDate != nil {
		return result, inventorycontrol.ErrConflict
	}
	result.LotID = lotID
	var balanceVersion int64
	err = tx.QueryRow(ctx, `insert into inventory.bulk_balance(tenant_id,organization_id,bin_id,item_id,lot_id,quantity,reserved_quantity,version) values($1,$2,$3,$4,nullif($5,''),$6::numeric,0,1) on conflict (tenant_id,organization_id,bin_id,item_id,lot_id) do update set quantity=inventory.bulk_balance.quantity+excluded.quantity,version=inventory.bulk_balance.version+1,updated_at=clock_timestamp() where (select max_quantity is null or inventory.bulk_balance.quantity+excluded.quantity<=max_quantity from inventory.item_bin_policy where tenant_id=$1 and organization_id=$2 and item_id=$4 and bin_id=$3) returning version`, tenant, value.OrganizationID, value.BinID, value.ItemID, lotID, value.Quantity).Scan(&balanceVersion)
	if errors.Is(err, pgx.ErrNoRows) {
		return result, inventorycontrol.ErrConflict
	}
	if err != nil {
		return result, bulkConflict(err)
	}
	err = tx.QueryRow(ctx, `insert into inventory.bulk_inventory_entry(tenant_id,entry_id,organization_id,bin_id,item_id,lot_id,entry_type,quantity,unit_cost,cost_amount,posting_date,source_kind,source_id) values($1,$2,$3,$4,$5,nullif($6,''),'receipt',$7::numeric,$8::numeric,round($7::numeric*$8::numeric,4),$9::date,$10,$11) returning cost_amount::text`, tenant, entryID, value.OrganizationID, value.BinID, value.ItemID, lotID, value.Quantity, value.UnitCost, value.PostingDate, value.SourceKind, value.SourceID).Scan(&result.CostAmount)
	if err != nil {
		return result, bulkConflict(err)
	}
	_, err = tx.Exec(ctx, `insert into inventory.bulk_cost_layer(tenant_id,layer_id,receipt_entry_id,organization_id,item_id,lot_id,posting_date,original_quantity,remaining_quantity,unit_cost) values($1,$2,$3,$4,$5,nullif($6,''),$7::date,$8::numeric,$8::numeric,$9::numeric)`, tenant, layerID, entryID, value.OrganizationID, value.ItemID, lotID, value.PostingDate, value.Quantity, value.UnitCost)
	if err != nil {
		return result, bulkConflict(err)
	}
	if err = recordBulkEvent(ctx, tx, tenant, eventID, "bulk-inventory-entry", entryID, "bulk-inventory.received", 1, map[string]string{"organization_id": value.OrganizationID, "bin_id": value.BinID, "item_id": value.ItemID, "lot_id": lotID, "quantity": value.Quantity, "cost_amount": result.CostAmount}); err != nil {
		return result, err
	}
	return result, tx.Commit(ctx)
}

func (r *InventoryControl) BulkAvailability(ctx context.Context, tenant, organization, item string) ([]inventorycontrol.BulkAvailability, error) {
	rows, err := r.pool.Query(ctx, `select organization_id,bin_id,bin_code,item_id,coalesce(lot_id,''),coalesce(lot_no,''),expiration_date,quantity::text,reserved_quantity::text,available_quantity::text from inventory.bulk_available where tenant_id=$1 and organization_id=$2 and item_id=$3 and available_quantity>0 order by expiration_date nulls last,lot_no,bin_code`, tenant, organization, item)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	values := []inventorycontrol.BulkAvailability{}
	for rows.Next() {
		var value inventorycontrol.BulkAvailability
		if err = rows.Scan(&value.OrganizationID, &value.BinID, &value.BinCode, &value.ItemID, &value.LotID, &value.LotNo, &value.ExpirationDate, &value.Quantity, &value.ReservedQuantity, &value.AvailableQuantity); err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return values, rows.Err()
}

func (r *InventoryControl) ReserveBulk(ctx context.Context, tenant, eventID string, value inventorycontrol.BulkReservation) (inventorycontrol.BulkReservation, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return value, err
	}
	defer tx.Rollback(ctx)
	value, err = reserveBulkInTx(ctx, tx, tenant, eventID, value)
	if err != nil {
		return value, err
	}
	return value, tx.Commit(ctx)
}

// AUTHORED transaction composition only: original writer SQL/calculation preserved.
func reserveBulkInTx(ctx context.Context, tx pgx.Tx, tenant, eventID string, value inventorycontrol.BulkReservation) (inventorycontrol.BulkReservation, error) {
	var err error
	var balanceVersion int64
	err = tx.QueryRow(ctx, `update inventory.bulk_balance b set reserved_quantity=reserved_quantity+$6::numeric,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and bin_id=$3 and item_id=$4 and lot_id is not distinct from nullif($5,'') and quantity-reserved_quantity >= $6::numeric and exists(select 1 from inventory.item_bin_policy p join inventory.warehouse_bin w using(tenant_id,organization_id,bin_id) left join inventory.inventory_lot l on l.tenant_id=b.tenant_id and l.lot_id=b.lot_id where p.tenant_id=b.tenant_id and p.organization_id=b.organization_id and p.item_id=b.item_id and p.bin_id=b.bin_id and not p.dedicated and not w.movement_blocked and coalesce(l.blocked,false)=false and (l.expiration_date is null or l.expiration_date>=current_date)) returning version`, tenant, value.OrganizationID, value.BinID, value.ItemID, value.LotID, value.Quantity).Scan(&balanceVersion)
	if errors.Is(err, pgx.ErrNoRows) {
		return value, inventorycontrol.ErrConflict
	}
	if err != nil {
		return value, err
	}
	_, err = tx.Exec(ctx, `insert into inventory.bulk_reservation(tenant_id,reservation_id,organization_id,bin_id,item_id,lot_id,demand_kind,demand_id,demand_line_id,quantity,status,cancellation_disallowed,expires_at,version) values($1,$2,$3,$4,$5,nullif($6,''),$7,$8,$9,$10::numeric,'reservation',$11,$12,1)`, tenant, value.ID, value.OrganizationID, value.BinID, value.ItemID, value.LotID, value.DemandKind, value.DemandID, value.DemandLineID, value.Quantity, value.CancellationDisallowed, value.ExpiresAt)
	if err != nil {
		return value, bulkConflict(err)
	}
	if err = recordBulkEvent(ctx, tx, tenant, eventID, "bulk-reservation", value.ID, "bulk-reservation.created", 1, map[string]string{"organization_id": value.OrganizationID, "bin_id": value.BinID, "item_id": value.ItemID, "lot_id": value.LotID, "quantity": value.Quantity}); err != nil {
		return value, err
	}
	return value, nil
}

func (r *InventoryControl) ReleaseBulk(ctx context.Context, tenant, organization, reservationID string, version int64, eventID string) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err = releaseBulkInTx(ctx, tx, tenant, organization, reservationID, version, eventID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// AUTHORED transaction composition: original release writer and SQL retained.
func releaseBulkInTx(ctx context.Context, tx pgx.Tx, tenant, organization, reservationID string, version int64, eventID string) error {
	var err error
	var binID, itemID, lotID, quantity string
	err = tx.QueryRow(ctx, `update inventory.bulk_reservation set status='released',version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and reservation_id=$3 and status='reservation' and version=$4 and not cancellation_disallowed returning bin_id,item_id,coalesce(lot_id,''),quantity::text`, tenant, organization, reservationID, version).Scan(&binID, &itemID, &lotID, &quantity)
	if errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.ErrConflict
	}
	if err != nil {
		return err
	}
	result, err := tx.Exec(ctx, `update inventory.bulk_balance set reserved_quantity=reserved_quantity-$6::numeric,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and bin_id=$3 and item_id=$4 and lot_id is not distinct from nullif($5,'') and reserved_quantity >= $6::numeric`, tenant, organization, binID, itemID, lotID, quantity)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return inventorycontrol.ErrConflict
	}
	if err = recordBulkEvent(ctx, tx, tenant, eventID, "bulk-reservation", reservationID, "bulk-reservation.released", version+1, map[string]string{"quantity": quantity}); err != nil {
		return err
	}
	return nil
}

func (r *InventoryControl) MoveBulk(ctx context.Context, tenant, eventOutID, eventInID string, value inventorycontrol.BulkMovement) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `update inventory.bulk_balance set quantity=quantity-$6::numeric,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and bin_id=$3 and item_id=$4 and lot_id is not distinct from nullif($5,'') and quantity-reserved_quantity >= $6::numeric and exists(select 1 from inventory.warehouse_bin where tenant_id=$1 and organization_id=$2 and bin_id=$3 and not movement_blocked)`, tenant, value.OrganizationID, value.FromBinID, value.ItemID, value.LotID, value.Quantity)
	if err != nil {
		return bulkConflict(err)
	}
	if result.RowsAffected() != 1 {
		return inventorycontrol.ErrConflict
	}
	var targetVersion int64
	err = tx.QueryRow(ctx, `insert into inventory.bulk_balance(tenant_id,organization_id,bin_id,item_id,lot_id,quantity,reserved_quantity,version) values($1,$2,$3,$4,nullif($5,''),$6::numeric,0,1) on conflict (tenant_id,organization_id,bin_id,item_id,lot_id) do update set quantity=inventory.bulk_balance.quantity+excluded.quantity,version=inventory.bulk_balance.version+1,updated_at=clock_timestamp() where exists(select 1 from inventory.item_bin_policy p join inventory.warehouse_bin w using(tenant_id,organization_id,bin_id) where p.tenant_id=$1 and p.organization_id=$2 and p.item_id=$4 and p.bin_id=$3 and not w.movement_blocked and (p.max_quantity is null or inventory.bulk_balance.quantity+excluded.quantity<=p.max_quantity)) returning version`, tenant, value.OrganizationID, value.ToBinID, value.ItemID, value.LotID, value.Quantity).Scan(&targetVersion)
	if errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.ErrConflict
	}
	if err != nil {
		return bulkConflict(err)
	}
	_, err = tx.Exec(ctx, `insert into inventory.bulk_inventory_entry(tenant_id,entry_id,organization_id,bin_id,item_id,lot_id,entry_type,quantity,unit_cost,cost_amount,posting_date,source_kind,source_id) values($1,$2,$3,$4,$5,nullif($6,''),'movement-out',-$7::numeric,0,0,$8::date,'bin-movement',$9),($1,$10,$3,$11,$5,nullif($6,''),'movement-in',$7::numeric,0,0,$8::date,'bin-movement',$9)`, tenant, eventOutID, value.OrganizationID, value.FromBinID, value.ItemID, value.LotID, value.Quantity, value.PostingDate, value.SourceID, eventInID, value.ToBinID)
	if err != nil {
		return bulkConflict(err)
	}
	if err = recordBulkEvent(ctx, tx, tenant, eventOutID, "bulk-bin-movement", value.SourceID, "bulk-inventory.moved", 1, map[string]string{"organization_id": value.OrganizationID, "from_bin_id": value.FromBinID, "to_bin_id": value.ToBinID, "item_id": value.ItemID, "lot_id": value.LotID, "quantity": value.Quantity}); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

type bulkAllocation struct {
	layerID, inboundEntryID, quantity, unitCost, costAmount string
}

func parseRat(value string) (*big.Rat, bool) { return new(big.Rat).SetString(value) }

func minRat(a, b *big.Rat) *big.Rat {
	if a.Cmp(b) <= 0 {
		return new(big.Rat).Set(a)
	}
	return new(big.Rat).Set(b)
}

func formatRat(value *big.Rat, scale int) string {
	factor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(scale)), nil)
	numerator := new(big.Int).Mul(value.Num(), factor)
	quotient, remainder := new(big.Int), new(big.Int)
	quotient.QuoRem(numerator, value.Denom(), remainder)
	if new(big.Int).Mul(remainder, big.NewInt(2)).Cmp(value.Denom()) >= 0 {
		quotient.Add(quotient, big.NewInt(1))
	}
	digits := quotient.String()
	if scale == 0 {
		return digits
	}
	for len(digits) <= scale {
		digits = "0" + digits
	}
	whole, fraction := digits[:len(digits)-scale], strings.TrimRight(digits[len(digits)-scale:], "0")
	if fraction == "" {
		return whole
	}
	return whole + "." + fraction
}

func (r *InventoryControl) IssueBulk(ctx context.Context, tenant, eventID, outboundEntryID string, value inventorycontrol.BulkIssue) (inventorycontrol.BulkIssueResult, error) {
	result := inventorycontrol.BulkIssueResult{EntryID: outboundEntryID}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return result, err
	}
	defer tx.Rollback(ctx)
	result, err = issueBulkInTx(ctx, tx, tenant, eventID, outboundEntryID, value)
	if err != nil {
		return result, err
	}
	return result, tx.Commit(ctx)
}

// AUTHORED transaction composition only: original FIFO/specific-cost writer preserved.
func issueBulkInTx(ctx context.Context, tx pgx.Tx, tenant, eventID, outboundEntryID string, value inventorycontrol.BulkIssue) (inventorycontrol.BulkIssueResult, error) {
	var err error
	result := inventorycontrol.BulkIssueResult{EntryID: outboundEntryID}
	var binID, itemID, lotID, quantity, costing string
	err = tx.QueryRow(ctx, `select r.bin_id,r.item_id,coalesce(r.lot_id,''),r.quantity::text,i.costing_method from inventory.bulk_reservation r join inventory.stock_item i using(tenant_id,item_id) left join inventory.inventory_lot l on l.tenant_id=r.tenant_id and l.lot_id=r.lot_id where r.tenant_id=$1 and r.organization_id=$2 and r.reservation_id=$3 and r.status='reservation' and r.version=$4 and (r.expires_at is null or r.expires_at>clock_timestamp()) and coalesce(l.blocked,false)=false and (l.expiration_date is null or l.expiration_date>=current_date) for update of r`, tenant, value.OrganizationID, value.ReservationID, value.ReservationVersion).Scan(&binID, &itemID, &lotID, &quantity, &costing)
	if errors.Is(err, pgx.ErrNoRows) {
		return result, inventorycontrol.ErrConflict
	}
	if err != nil {
		return result, err
	}
	if (costing == "specific") != (value.SpecificReceiptEntry != "") {
		return result, inventorycontrol.ErrConflict
	}
	wanted, ok := parseRat(quantity)
	if !ok {
		return result, inventorycontrol.ErrConflict
	}
	query := `select layer_id,receipt_entry_id,remaining_quantity::text,unit_cost::text from inventory.bulk_cost_layer where tenant_id=$1 and organization_id=$2 and item_id=$3 and lot_id is not distinct from nullif($4,'') and remaining_quantity>0`
	args := []any{tenant, value.OrganizationID, itemID, lotID}
	if costing == "specific" {
		query += ` and receipt_entry_id=$5`
		args = append(args, value.SpecificReceiptEntry)
	}
	query += ` order by posting_date,receipt_entry_id for update`
	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return result, err
	}
	remaining := new(big.Rat).Set(wanted)
	applications := []bulkAllocation{}
	totalCost := new(big.Rat)
	for rows.Next() && remaining.Sign() > 0 {
		var layerID, inbound, availableText, unitCost string
		if err = rows.Scan(&layerID, &inbound, &availableText, &unitCost); err != nil {
			rows.Close()
			return result, err
		}
		available, validAvailable := parseRat(availableText)
		unit, validUnit := parseRat(unitCost)
		if !validAvailable || !validUnit {
			rows.Close()
			return result, inventorycontrol.ErrConflict
		}
		take := minRat(available, remaining)
		cost := new(big.Rat).Mul(take, unit)
		costText := formatRat(cost, 4)
		rounded, _ := parseRat(costText)
		totalCost.Add(totalCost, rounded)
		applications = append(applications, bulkAllocation{layerID: layerID, inboundEntryID: inbound, quantity: formatRat(take, 6), unitCost: unitCost, costAmount: costText})
		remaining.Sub(remaining, take)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return result, err
	}
	if remaining.Sign() != 0 || len(applications) == 0 {
		return result, inventorycontrol.ErrConflict
	}
	result.Quantity, result.CostAmount, result.Applications = quantity, formatRat(totalCost, 4), len(applications)
	unitCost := formatRat(new(big.Rat).Quo(totalCost, wanted), 4)
	_, err = tx.Exec(ctx, `insert into inventory.bulk_inventory_entry(tenant_id,entry_id,organization_id,bin_id,item_id,lot_id,entry_type,quantity,unit_cost,cost_amount,posting_date,source_kind,source_id) values($1,$2,$3,$4,$5,nullif($6,''),'issue',-$7::numeric,$8::numeric,-$9::numeric,$10::date,$11,$12)`, tenant, outboundEntryID, value.OrganizationID, binID, itemID, lotID, quantity, unitCost, result.CostAmount, value.PostingDate, value.SourceKind, value.SourceID)
	if err != nil {
		return result, bulkConflict(err)
	}
	for _, application := range applications {
		updated, updateErr := tx.Exec(ctx, `update inventory.bulk_cost_layer set remaining_quantity=remaining_quantity-$3::numeric where tenant_id=$1 and layer_id=$2 and remaining_quantity >= $3::numeric`, tenant, application.layerID, application.quantity)
		if updateErr != nil {
			return result, updateErr
		}
		if updated.RowsAffected() != 1 {
			return result, inventorycontrol.ErrConflict
		}
		_, err = tx.Exec(ctx, `insert into inventory.bulk_cost_application(tenant_id,outbound_entry_id,inbound_entry_id,quantity,cost_amount) values($1,$2,$3,$4::numeric,$5::numeric)`, tenant, outboundEntryID, application.inboundEntryID, application.quantity, application.costAmount)
		if err != nil {
			return result, bulkConflict(err)
		}
	}
	updated, err := tx.Exec(ctx, `update inventory.bulk_balance set quantity=quantity-$6::numeric,reserved_quantity=reserved_quantity-$6::numeric,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and bin_id=$3 and item_id=$4 and lot_id is not distinct from nullif($5,'') and quantity >= $6::numeric and reserved_quantity >= $6::numeric`, tenant, value.OrganizationID, binID, itemID, lotID, quantity)
	if err != nil {
		return result, err
	}
	if updated.RowsAffected() != 1 {
		return result, inventorycontrol.ErrConflict
	}
	updated, err = tx.Exec(ctx, `update inventory.bulk_reservation set status='consumed',version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and reservation_id=$2 and status='reservation' and version=$3`, tenant, value.ReservationID, value.ReservationVersion)
	if err != nil {
		return result, err
	}
	if updated.RowsAffected() != 1 {
		return result, inventorycontrol.ErrConflict
	}
	if err = recordBulkEvent(ctx, tx, tenant, eventID, "bulk-inventory-entry", outboundEntryID, "bulk-inventory.issued", 1, map[string]string{"organization_id": value.OrganizationID, "item_id": itemID, "lot_id": lotID, "quantity": quantity, "cost_amount": result.CostAmount, "reservation_id": value.ReservationID}); err != nil {
		return result, err
	}
	return result, nil
}
