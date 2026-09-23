// AUTHORED serial ASN/partial receiving composition over existing state writers.
package postgres

import (
	"context"
	"elite.local/enterprise/internal/operations"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/randomid"
	sc "elite.local/enterprise/internal/serialsupply"
	"github.com/jackc/pgx/v5"
	"sort"
)

func (s *SerialSupply) supplyShip(ctx context.Context, tx pgx.Tx, p identity.Principal, plan sc.Plan, r sc.Command) (any, error) {
	if plan.State != "in-production" {
		return nil, sc.ErrConflict
	}
	_, err := tx.Exec(ctx, `insert into procurement.serial_supply_shipment(tenant_id,purchase_order_id,shipment_id,command_id) values($1,$2,$3,$4)`, p.TenantID, plan.PurchaseOrderID, r.ShipmentID, r.CommandID)
	if err != nil {
		return nil, supplyError(err)
	}
	ids := append([]string(nil), r.Units...)
	sort.Strings(ids)
	lines := []map[string]string{}
	svc := operations.NewService(operationsTxRepository{tx: tx}, randomid.Generator{})
	for _, id := range ids {
		u, err := readSupplyUnit(ctx, tx, p.TenantID, plan.PurchaseOrderID, id)
		if err != nil {
			return nil, err
		}
		if u.State != "released" || u.StockID != "" || u.ShipmentID != "" {
			return nil, sc.ErrConflict
		}
		q, err := readSupplyQuality(ctx, tx, p.TenantID, u.ID, "factory")
		if err != nil {
			return nil, err
		}
		if q.State != "approved" || q.SourceState != "quality" {
			return nil, sc.ErrConflict
		}
		if err = moveSupplyFactory(ctx, tx, p.TenantID, plan, r, u.ID, u.State, "shipped"); err != nil {
			return nil, err
		}
		stock, err := svc.ReceiveStockUnit(ctx, p.TenantID, operations.StockUnit{OrganizationID: plan.DestinationOrganizationID, VariantID: u.VariantID, ProductionUnitID: u.ID, SerialNumber: u.Serial, VIN: u.VIN, BatterySerialNumber: u.Battery})
		if err != nil {
			return nil, supplyError(err)
		}
		_, err = tx.Exec(ctx, `insert into procurement.serial_supply_manifest(tenant_id,purchase_order_id,shipment_id,production_unit_id,stock_unit_id) values($1,$2,$3,$4,$5)`, p.TenantID, plan.PurchaseOrderID, r.ShipmentID, u.ID, stock.ID)
		if err != nil {
			return nil, supplyError(err)
		}
		if err = supplyEffect(ctx, tx, p.TenantID, plan.PurchaseOrderID, r.CommandID, "stock", stock.ID, "", "in-transit", 0, 1); err != nil {
			return nil, err
		}
		lines = append(lines, map[string]string{"unit_id": u.ID, "stock_unit_id": stock.ID, "serial_number": u.Serial, "variant_id": u.VariantID, "line_id": u.LineID})
	}
	var complete bool
	err = tx.QueryRow(ctx, `select not exists(select 1 from procurement.serial_supply_line l where tenant_id=$1 and purchase_order_id=$2
 and l.quantity<>(select count(*) from procurement.serial_supply_unit b join procurement.serial_supply_manifest m using(tenant_id,production_unit_id)
 where b.tenant_id=l.tenant_id and b.purchase_order_id=l.purchase_order_id and b.line_id=l.line_id))`, p.TenantID, plan.PurchaseOrderID).Scan(&complete)
	if err != nil {
		return nil, err
	}
	if complete {
		if err = moveSupplyPurchase(ctx, tx, p.TenantID, plan, r, "shipped"); err != nil {
			return nil, err
		}
	}
	return map[string]any{"shipment_id": r.ShipmentID, "manifest": lines, "all_declared_units_shipped": complete}, nil
}
func (s *SerialSupply) supplyReceive(ctx context.Context, tx pgx.Tx, p identity.Principal, plan sc.Plan, r sc.Command) (any, error) {
	if plan.State != "in-production" && plan.State != "shipped" {
		return nil, sc.ErrConflict
	}
	ids := append([]string(nil), r.Units...)
	sort.Strings(ids)
	received := []map[string]string{}
	for _, id := range ids {
		u, err := readSupplyUnit(ctx, tx, p.TenantID, plan.PurchaseOrderID, id)
		if err != nil {
			return nil, err
		}
		if u.State != "shipped" || u.StockState != "in-transit" || u.ShipmentID != r.ShipmentID || u.StockID == "" {
			return nil, sc.ErrConflict
		}
		_, err = tx.Exec(ctx, `insert into procurement.serial_supply_receipt(tenant_id,purchase_order_id,shipment_id,production_unit_id,command_id) values($1,$2,$3,$4,$5)`, p.TenantID, plan.PurchaseOrderID, r.ShipmentID, id, r.CommandID)
		if err != nil {
			return nil, supplyError(err)
		}
		if err = moveSupplyFactory(ctx, tx, p.TenantID, plan, r, id, u.State, "received"); err != nil {
			return nil, err
		}
		if err = moveSupplyStock(ctx, tx, p.TenantID, plan, r, u.StockID, u.StockState, "quarantine", u.StockVersion); err != nil {
			return nil, err
		}
		u.State = "received"
		u.StockState = "quarantine"
		u.StockVersion++
		approvalID, err := s.proposeSupplyQuality(ctx, tx, p, plan, r, u, "receipt")
		if err != nil {
			return nil, err
		}
		received = append(received, map[string]string{"unit_id": id, "stock_unit_id": u.StockID, "approval_id": approvalID, "state": "quarantine"})
	}
	var complete bool
	err := tx.QueryRow(ctx, `select not exists(select 1 from procurement.serial_supply_line l where tenant_id=$1 and purchase_order_id=$2
 and l.quantity<>(select count(*) from procurement.serial_supply_unit b join procurement.serial_supply_receipt r using(tenant_id,production_unit_id)
 where b.tenant_id=l.tenant_id and b.purchase_order_id=l.purchase_order_id and b.line_id=l.line_id))`, p.TenantID, plan.PurchaseOrderID).Scan(&complete)
	if err != nil {
		return nil, err
	}
	if complete {
		if plan.State != "shipped" {
			return nil, sc.ErrConflict
		}
		if err = moveSupplyPurchase(ctx, tx, p.TenantID, plan, r, "received"); err != nil {
			return nil, err
		}
	}
	return map[string]any{"shipment_id": r.ShipmentID, "received": received, "all_declared_units_received": complete, "payment_created": false}, nil
}
func (s *SerialSupply) Apply(ctx context.Context, p identity.Principal, r sc.Command) (sc.Receipt, error) {
	if !r.Valid() {
		return sc.Receipt{}, sc.ErrInvalid
	}
	var action supplyAction
	switch r.Kind {
	case "submit", "confirm", "start", "cancel":
		action = func(ctx context.Context, tx pgx.Tx, p identity.Principal, plan sc.Plan, r sc.Command) (any, error) {
			return supplyLifecycle(ctx, tx, p.TenantID, plan, r)
		}
	case "register":
		action = func(ctx context.Context, tx pgx.Tx, p identity.Principal, plan sc.Plan, r sc.Command) (any, error) {
			return registerSupplyUnit(ctx, tx, p.TenantID, plan, r)
		}
	case "milestone":
		action = s.supplyMilestone
	case "ship":
		action = s.supplyShip
	case "receive":
		action = s.supplyReceive
	case "quality", "quality-reject", "reinspect":
		action = s.supplyReceiptQuality
	default:
		return sc.Receipt{}, sc.ErrInvalid
	}
	return s.executeSupplyCommand(ctx, p, r, action)
}
