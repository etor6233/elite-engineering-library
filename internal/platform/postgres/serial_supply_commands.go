// AUTHORED transaction composition. Domain preparation/transitions remain in
// operations.Service and its original PostgreSQL writers.
package postgres

import (
	"context"
	"elite.local/enterprise/internal/operations"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/randomid"
	sc "elite.local/enterprise/internal/serialsupply"
	"errors"
	"github.com/jackc/pgx/v5"
)

func supplyCommandPermission(r sc.Command, plan sc.Plan) (string, string) {
	switch r.Kind {
	case "confirm", "start", "register", "milestone", "ship":
		return plan.FactoryOrganizationID, "supply:factory"
	case "receive":
		return plan.DestinationOrganizationID, "supply:receive"
	case "quality", "quality-reject":
		return plan.DestinationOrganizationID, "supply:release"
	case "reinspect":
		return plan.DestinationOrganizationID, "supply:inspect"
	default:
		return plan.DestinationOrganizationID, "supply:plan"
	}
}
func moveSupplyPurchase(ctx context.Context, tx pgx.Tx, tenant string, p sc.Plan, r sc.Command, target string) error {
	svc := operations.NewService(operationsTxRepository{tx: tx}, randomid.Generator{})
	if err := svc.TransitionPurchaseOrder(ctx, tenant, p.DestinationOrganizationID, p.PurchaseOrderID, p.State, target, p.PurchaseVersion); err != nil {
		return supplyError(err)
	}
	return supplyEffect(ctx, tx, tenant, p.PurchaseOrderID, r.CommandID, "purchase", p.PurchaseOrderID, p.State, target, p.PurchaseVersion, p.PurchaseVersion+1)
}
func moveSupplyFactory(ctx context.Context, tx pgx.Tx, tenant string, p sc.Plan, r sc.Command, id, from, to string) error {
	svc := operations.NewService(operationsTxRepository{tx: tx}, randomid.Generator{})
	if err := svc.TransitionProductionUnit(ctx, tenant, p.DestinationOrganizationID, id, from, to); err != nil {
		return supplyError(err)
	}
	return supplyEffect(ctx, tx, tenant, p.PurchaseOrderID, r.CommandID, "factory", id, from, to, 0, 0)
}
func moveSupplyStock(ctx context.Context, tx pgx.Tx, tenant string, p sc.Plan, r sc.Command, id, from, to string, version int64) error {
	svc := operations.NewService(operationsTxRepository{tx: tx}, randomid.Generator{})
	if err := svc.TransitionStockUnit(ctx, tenant, p.DestinationOrganizationID, id, from, to, version); err != nil {
		return supplyError(err)
	}
	return supplyEffect(ctx, tx, tenant, p.PurchaseOrderID, r.CommandID, "stock", id, from, to, version, version+1)
}

type supplyUnitState struct {
	ID, LineID, VariantID, State, Serial, VIN, Battery string
	StockID, StockState, ShipmentID                    string
	StockVersion                                       int64
}

func readSupplyUnit(ctx context.Context, tx pgx.Tx, tenant, po, id string) (supplyUnitState, error) {
	var u supplyUnitState
	err := tx.QueryRow(ctx, `select u.production_unit_id,b.line_id,u.variant_id,u.state,u.serial_number,coalesce(u.vin,''),coalesce(u.battery_serial_number,''),
 coalesce(m.stock_unit_id,''),coalesce(i.state,''),coalesce(m.shipment_id,''),coalesce(i.version,0)
 from procurement.serial_supply_unit b join factory.production_unit u using(tenant_id,production_unit_id)
 left join procurement.serial_supply_manifest m using(tenant_id,production_unit_id)
 left join inventory.stock_unit i on i.tenant_id=m.tenant_id and i.stock_unit_id=m.stock_unit_id
 where b.tenant_id=$1 and b.purchase_order_id=$2 and b.production_unit_id=$3 for update of u`, tenant, po, id).Scan(&u.ID, &u.LineID, &u.VariantID, &u.State, &u.Serial, &u.VIN, &u.Battery, &u.StockID, &u.StockState, &u.ShipmentID, &u.StockVersion)
	if err != nil {
		return u, supplyError(err)
	}
	if u.StockID != "" {
		// Lock the stock owner too, rather than relying on a stale outer-join snapshot.
		var org, variant, unit, serial, vin, battery string
		err = tx.QueryRow(ctx, `select organization_id,variant_id,coalesce(production_unit_id,''),serial_number,coalesce(vin,''),coalesce(battery_serial_number,''),state,version
  from inventory.stock_unit where tenant_id=$1 and stock_unit_id=$2 for update`, tenant, u.StockID).Scan(&org, &variant, &unit, &serial, &vin, &battery, &u.StockState, &u.StockVersion)
		if err != nil {
			return u, supplyError(err)
		}
		var destination string
		if err = tx.QueryRow(ctx, `select destination_organization_id from procurement.serial_supply_plan where tenant_id=$1 and purchase_order_id=$2`, tenant, po).Scan(&destination); err != nil {
			return u, err
		}
		if org != destination || variant != u.VariantID || unit != u.ID || serial != u.Serial || vin != u.VIN || battery != u.Battery {
			return u, sc.ErrConflict
		}
	}
	return u, nil
}
func registerSupplyUnit(ctx context.Context, tx pgx.Tx, tenant string, p sc.Plan, r sc.Command) (any, error) {
	if p.State != "in-production" {
		return nil, sc.ErrConflict
	}
	var variant string
	var maximum, active int
	err := tx.QueryRow(ctx, `select l.variant_id,l.quantity,
 (select count(*) from procurement.serial_supply_unit b join factory.production_unit u using(tenant_id,production_unit_id)
 where b.tenant_id=l.tenant_id and b.purchase_order_id=l.purchase_order_id and b.line_id=l.line_id and u.state<>'rejected')
 from procurement.serial_supply_line l where l.tenant_id=$1 and l.purchase_order_id=$2 and l.line_id=$3`, tenant, p.PurchaseOrderID, r.LineID).Scan(&variant, &maximum, &active)
	if err != nil {
		return nil, supplyError(err)
	}
	if active >= maximum {
		return nil, sc.ErrConflict
	}
	svc := operations.NewService(operationsTxRepository{tx: tx}, randomid.Generator{})
	unit, err := svc.RegisterProductionUnit(ctx, tenant, operations.ProductionUnit{OrganizationID: p.DestinationOrganizationID, PurchaseOrderID: p.PurchaseOrderID, VariantID: variant, SerialNumber: r.SerialNumber, VIN: r.VIN, BatterySerialNumber: r.BatterySerialNumber})
	if err != nil {
		return nil, supplyError(err)
	}
	_, err = tx.Exec(ctx, `insert into procurement.serial_supply_unit(tenant_id,purchase_order_id,line_id,production_unit_id,command_id) values($1,$2,$3,$4,$5)`, tenant, p.PurchaseOrderID, r.LineID, unit.ID, r.CommandID)
	if err != nil {
		return nil, supplyError(err)
	}
	if err = supplyEffect(ctx, tx, tenant, p.PurchaseOrderID, r.CommandID, "factory", unit.ID, "", "planned", 0, 0); err != nil {
		return nil, err
	}
	return map[string]any{"unit": unit, "line_id": r.LineID}, nil
}
func supplyLifecycle(ctx context.Context, tx pgx.Tx, tenant string, p sc.Plan, r sc.Command) (any, error) {
	target := map[string]string{"submit": "submitted", "confirm": "accepted", "start": "in-production", "cancel": "cancelled"}[r.Kind]
	if target == "" {
		return nil, sc.ErrInvalid
	}
	if r.Kind == "confirm" {
		var active bool
		if err := tx.QueryRow(ctx, `select exists(select 1 from partner.supplier where tenant_id=$1 and supplier_id=$2 and status='active')`, tenant, p.SupplierID).Scan(&active); err != nil {
			return nil, err
		}
		if !active {
			return nil, sc.ErrConflict
		}
	}
	if err := moveSupplyPurchase(ctx, tx, tenant, p, r, target); err != nil {
		return nil, err
	}
	payload := map[string]any{"purchase_order_id": p.PurchaseOrderID, "from_state": p.State, "to_state": target, "supplier_id": p.SupplierID}
	if r.Kind == "confirm" {
		payload["confirmation_kind"] = "authorized-recorded-evidence"
	}
	return payload, nil
}

type supplyAction func(context.Context, pgx.Tx, identity.Principal, sc.Plan, sc.Command) (any, error)

// executeSupplyCommand supplies the single lock/replay/receipt boundary for all
// commands. The action is selected internally by the typed public dispatcher.
func (s *SerialSupply) executeSupplyCommand(ctx context.Context, p identity.Principal, r sc.Command, action supplyAction) (sc.Receipt, error) {
	if s == nil || !r.Valid() || p.TenantID == "" || action == nil {
		return sc.Receipt{}, sc.ErrInvalid
	}
	_, hash, err := sc.Canonical(r)
	if err != nil {
		return sc.Receipt{}, err
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return sc.Receipt{}, err
	}
	defer tx.Rollback(ctx)
	plan, err := readSupplyPlan(ctx, tx, p.TenantID, r.PurchaseOrderID, true)
	if err != nil {
		return sc.Receipt{}, err
	}
	org, permission := supplyCommandPermission(r, plan)
	if !approvalPrincipal(p, p.TenantID, org, permission) {
		return sc.Receipt{}, sc.ErrNotFound
	}
	old, err := readSupplyReceipt(ctx, tx, p.TenantID, r.PurchaseOrderID, r.CommandID)
	if err == nil {
		if old.Kind != r.Kind || old.Actor != p.Subject || old.RequestSHA256 != hash {
			return sc.Receipt{}, sc.ErrConflict
		}
		old.Replay = true
		return old, supplyError(tx.Commit(ctx))
	}
	if !errors.Is(err, sc.ErrNotFound) {
		return sc.Receipt{}, err
	}
	if plan.Version != r.ExpectedVersion {
		return sc.Receipt{}, sc.ErrConflict
	}
	payload, err := action(ctx, tx, p, plan, r)
	if err != nil {
		return sc.Receipt{}, err
	}
	receipt, err := writeSupplyReceipt(ctx, tx, p.TenantID, r.PurchaseOrderID, r.CommandID, r.Kind, p.Subject, hash, r.EvidenceSHA256, plan.Version+1, map[string]any{"effect": payload, "evidence_sha256": r.EvidenceSHA256})
	if err != nil {
		return sc.Receipt{}, err
	}
	return receipt, supplyError(tx.Commit(ctx))
}
