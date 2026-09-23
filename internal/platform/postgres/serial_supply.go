// AUTHORED physical/actor/evidence binding and recovery around existing owners.
package postgres

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"

	"elite.local/enterprise/internal/platform/identity"
	sc "elite.local/enterprise/internal/serialsupply"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SerialSupply struct{ pool *pgxpool.Pool }

func NewSerialSupply(pool *pgxpool.Pool) (*SerialSupply, error) {
	if pool == nil {
		return nil, sc.ErrInvalid
	}
	return &SerialSupply{pool: pool}, nil
}
func supplyError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return sc.ErrNotFound
	}
	var p *pgconn.PgError
	if postgresConflict(err) || errors.As(err, &p) && p.Code == "P0001" {
		return sc.ErrConflict
	}
	return err
}

type supplyQuery interface {
	QueryRow(context.Context, string, ...any) pgx.Row
	Query(context.Context, string, ...any) (pgx.Rows, error)
}

func readSupplyPlan(ctx context.Context, q supplyQuery, tenant, po string, lock bool) (sc.Plan, error) {
	var p sc.Plan
	sql := `select s.purchase_order_id,s.destination_organization_id,s.factory_organization_id,s.supplier_id,p.state,p.version,
 coalesce((select max(version) from procurement.serial_supply_step where tenant_id=s.tenant_id and purchase_order_id=s.purchase_order_id),0),
 s.policy_code,s.demand_reference,p.currency,p.total_minor_units from procurement.serial_supply_plan s join procurement.purchase_order p using(tenant_id,purchase_order_id)
 where s.tenant_id=$1 and s.purchase_order_id=$2`
	if lock {
		sql += " for update of p"
	}
	err := q.QueryRow(ctx, sql, tenant, po).Scan(&p.PurchaseOrderID, &p.DestinationOrganizationID, &p.FactoryOrganizationID, &p.SupplierID, &p.State, &p.PurchaseVersion, &p.Version, &p.PolicyCode, &p.DemandReference, &p.Currency, &p.TotalMinorUnits)
	return p, supplyError(err)
}
func readSupplyReceipt(ctx context.Context, q supplyQuery, tenant, po, command string) (sc.Receipt, error) {
	var r sc.Receipt
	err := q.QueryRow(ctx, `select purchase_order_id,command_id,version,kind,actor,request_sha256,payload_sha256,payload,recorded_at
 from procurement.serial_supply_step where tenant_id=$1 and purchase_order_id=$2 and ($3::text='' or command_id=$3) order by version desc limit 1`, tenant, po, command).Scan(&r.PurchaseOrderID, &r.CommandID, &r.Version, &r.Kind, &r.Actor, &r.RequestSHA256, &r.PayloadSHA256, &r.Payload, &r.RecordedAt)
	if err != nil {
		return r, supplyError(err)
	}
	// CanonicalPayload already enforces exact JSON numbers; decode through that
	// shared owner directly, not via float64 in an interface.
	raw, hash, err := sc.Canonical(r.Payload)
	if err != nil || hash != r.PayloadSHA256 {
		return sc.Receipt{}, sc.ErrConflict
	}
	r.Payload = raw
	return r, nil
}
func supplyReadAllowed(p identity.Principal, plan sc.Plan) bool {
	return approvalPrincipal(p, p.TenantID, plan.DestinationOrganizationID, "supply:read") ||
		approvalPrincipal(p, p.TenantID, plan.FactoryOrganizationID, "supply:factory-read")
}
func (s *SerialSupply) Plan(ctx context.Context, p identity.Principal, id, afterUnit string) (sc.Plan, error) {
	if s == nil || !sc.ValidID(id) || afterUnit != "" && !sc.ValidID(afterUnit) || p.TenantID == "" {
		return sc.Plan{}, sc.ErrInvalid
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly})
	if err != nil {
		return sc.Plan{}, err
	}
	defer tx.Rollback(ctx)
	plan, err := readSupplyPlan(ctx, tx, p.TenantID, id, false)
	if err != nil {
		return sc.Plan{}, err
	}
	if !supplyReadAllowed(p, plan) {
		return sc.Plan{}, sc.ErrNotFound
	}
	rows, err := tx.Query(ctx, `select line_id,variant_id,quantity from procurement.serial_supply_line where tenant_id=$1 and purchase_order_id=$2 order by line_id`, p.TenantID, id)
	if err != nil {
		return sc.Plan{}, err
	}
	plan.Lines = []sc.Line{}
	for rows.Next() {
		var l sc.Line
		if err = rows.Scan(&l.ID, &l.VariantID, &l.Quantity); err != nil {
			rows.Close()
			return sc.Plan{}, err
		}
		plan.Lines = append(plan.Lines, l)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return sc.Plan{}, err
	}
	rows, err = tx.Query(ctx, `select u.production_unit_id,b.line_id,u.state,u.serial_number,coalesce(m.stock_unit_id,''),coalesce(i.state,''),coalesce(m.shipment_id,''),
 coalesce((select a.state from procurement.serial_supply_quality q join approval.request a on a.tenant_id=q.tenant_id and a.request_id=q.approval_id where q.tenant_id=u.tenant_id and q.production_unit_id=u.production_unit_id and q.stage='factory' order by q.attempt desc limit 1),''),
 coalesce((select a.requester from procurement.serial_supply_quality q join approval.request a on a.tenant_id=q.tenant_id and a.request_id=q.approval_id where q.tenant_id=u.tenant_id and q.production_unit_id=u.production_unit_id and q.stage='factory' order by q.attempt desc limit 1),''),
 coalesce((select a.state from procurement.serial_supply_quality q join approval.request a on a.tenant_id=q.tenant_id and a.request_id=q.approval_id where q.tenant_id=u.tenant_id and q.production_unit_id=u.production_unit_id and q.stage='receipt' order by q.attempt desc limit 1),''),
 coalesce((select a.requester from procurement.serial_supply_quality q join approval.request a on a.tenant_id=q.tenant_id and a.request_id=q.approval_id where q.tenant_id=u.tenant_id and q.production_unit_id=u.production_unit_id and q.stage='receipt' order by q.attempt desc limit 1),'')
 from procurement.serial_supply_unit b join factory.production_unit u using(tenant_id,production_unit_id)
 left join procurement.serial_supply_manifest m using(tenant_id,production_unit_id)
 left join inventory.stock_unit i on i.tenant_id=m.tenant_id and i.stock_unit_id=m.stock_unit_id
 where b.tenant_id=$1 and b.purchase_order_id=$2 and u.production_unit_id>$3 order by u.production_unit_id limit 101`, p.TenantID, id, afterUnit)
	if err != nil {
		return sc.Plan{}, err
	}
	plan.Units = []sc.Unit{}
	for rows.Next() {
		var u sc.Unit
		if err = rows.Scan(&u.ID, &u.LineID, &u.State, &u.SerialNumber, &u.StockUnitID, &u.StockState, &u.ShipmentID, &u.FactoryReviewState, &u.FactoryRequester, &u.ReceiptReviewState, &u.ReceiptRequester); err != nil {
			rows.Close()
			return sc.Plan{}, err
		}
		plan.Units = append(plan.Units, u)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return sc.Plan{}, err
	}
	if len(plan.Units) > 100 {
		plan.Units = plan.Units[:100]
		plan.NextUnitID = plan.Units[99].ID
	}
	latest, err := readSupplyReceipt(ctx, tx, p.TenantID, id, "")
	if err != nil {
		return sc.Plan{}, err
	}
	plan.Latest = &latest
	if err = tx.Commit(ctx); err != nil {
		return sc.Plan{}, err
	}
	return plan, nil
}
func (s *SerialSupply) CommandReceipt(ctx context.Context, p identity.Principal, id, command string) (sc.Receipt, error) {
	if s == nil || !sc.ValidID(id) || !sc.ValidID(command) || p.TenantID == "" {
		return sc.Receipt{}, sc.ErrInvalid
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly})
	if err != nil {
		return sc.Receipt{}, err
	}
	defer tx.Rollback(ctx)
	plan, err := readSupplyPlan(ctx, tx, p.TenantID, id, false)
	if err != nil {
		return sc.Receipt{}, err
	}
	if !supplyReadAllowed(p, plan) {
		return sc.Receipt{}, sc.ErrNotFound
	}
	r, err := readSupplyReceipt(ctx, tx, p.TenantID, id, command)
	if err != nil {
		return sc.Receipt{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return sc.Receipt{}, err
	}
	return r, nil
}
func writeSupplyReceipt(ctx context.Context, tx pgx.Tx, tenant, po, command, kind, actor, requestSHA, evidence string, version int64, payload any) (sc.Receipt, error) {
	raw, hash, err := sc.Canonical(payload)
	if err != nil {
		return sc.Receipt{}, err
	}
	_, err = tx.Exec(ctx, `insert into procurement.serial_supply_step(tenant_id,purchase_order_id,command_id,version,kind,actor,request_sha256,evidence_sha256,payload_sha256,payload) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`, tenant, po, command, version, kind, actor, requestSHA, evidence, hash, raw)
	if err != nil {
		return sc.Receipt{}, supplyError(err)
	}
	r, err := readSupplyReceipt(ctx, tx, tenant, po, command)
	if err != nil {
		return sc.Receipt{}, err
	}
	event, err := json.Marshal(r)
	if err != nil {
		return sc.Receipt{}, err
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)
 values($1,gen_random_uuid(),'serial-supply',$2,$3,$4,1,clock_timestamp(),$5)`, tenant, po, version, "serial-supply."+kind, event)
	return r, err
}
func supplyEffect(ctx context.Context, tx pgx.Tx, tenant, po, command, entity, id, from, to string, fromVersion, toVersion int64) error {
	_, err := tx.Exec(ctx, `insert into procurement.serial_supply_effect(tenant_id,purchase_order_id,command_id,entity,entity_id,from_state,to_state,from_version,to_version)
 values($1,$2,$3,$4,$5,$6,$7,$8,$9)`, tenant, po, command, entity, id, from, to, fromVersion, toVersion)
	return err
}
func (s *SerialSupply) BindPlan(ctx context.Context, p identity.Principal, r sc.PlanRequest) (sc.Receipt, error) {
	if s == nil || !r.Valid() || p.TenantID == "" {
		return sc.Receipt{}, sc.ErrInvalid
	}
	_, requestSHA, err := sc.Canonical(r)
	if err != nil {
		return sc.Receipt{}, err
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return sc.Receipt{}, err
	}
	defer tx.Rollback(ctx)
	receipt, err := bindSupplyPlanInTx(ctx, tx, p, r, requestSHA)
	if err != nil {
		return sc.Receipt{}, err
	}
	return receipt, tx.Commit(ctx)
}

// AUTHORED extraction: original BindPlan SQL/authorization/replay are retained.
func bindSupplyPlanInTx(ctx context.Context, tx pgx.Tx, p identity.Principal, r sc.PlanRequest, requestSHA string) (sc.Receipt, error) {
	var org, supplier, state string
	err := tx.QueryRow(ctx, `select destination_organization_id,supplier_id,state from procurement.purchase_order where tenant_id=$1 and purchase_order_id=$2 for update`, p.TenantID, r.PurchaseOrderID).Scan(&org, &supplier, &state)
	if err != nil {
		return sc.Receipt{}, supplyError(err)
	}
	if !approvalPrincipal(p, p.TenantID, org, "supply:plan") {
		return sc.Receipt{}, sc.ErrNotFound
	}
	existing, err := readSupplyReceipt(ctx, tx, p.TenantID, r.PurchaseOrderID, r.CommandID)
	if err == nil {
		if existing.Kind != "planned" || existing.Actor != p.Subject || existing.RequestSHA256 != requestSHA {
			return sc.Receipt{}, sc.ErrConflict
		}
		existing.Replay = true
		return existing, nil
	}
	if !errors.Is(err, sc.ErrNotFound) {
		return sc.Receipt{}, err
	}
	if state != "draft" {
		return sc.Receipt{}, sc.ErrConflict
	}
	var active bool
	err = tx.QueryRow(ctx, `select exists(select 1 from partner.supplier where tenant_id=$1 and supplier_id=$2 and status='active')
 and exists(select 1 from org.organization where tenant_id=$1 and organization_id=$3 and status='active')
 and exists(select 1 from org.organization where tenant_id=$1 and organization_id=$4 and status='active' and organization_type='factory') and not exists(select 1 from factory.production_unit where tenant_id=$1 and purchase_order_id=$5)`, p.TenantID, supplier, org, r.FactoryOrganizationID, r.PurchaseOrderID).Scan(&active)
	if err != nil {
		return sc.Receipt{}, err
	}
	if !active {
		return sc.Receipt{}, sc.ErrConflict
	}
	policyHash := sha256.Sum256([]byte(sc.PolicyJSON))
	_, err = tx.Exec(ctx, `insert into procurement.serial_supply_plan(tenant_id,purchase_order_id,destination_organization_id,factory_organization_id,supplier_id,demand_reference,policy_code,policy_sha256,created_by)
 values($1,$2,$3,$4,$5,$6,$7,$8,$9)`, p.TenantID, r.PurchaseOrderID, org, r.FactoryOrganizationID, supplier, r.DemandReference, sc.PolicyCode, hex.EncodeToString(policyHash[:]), p.Subject)
	if err != nil {
		return sc.Receipt{}, supplyError(err)
	}
	for _, line := range r.Lines {
		tag, err := tx.Exec(ctx, `insert into procurement.serial_supply_line(tenant_id,purchase_order_id,line_id,variant_id,quantity)
 select $1,$2,$3,variant_id,$5 from catalog.vehicle_variant where tenant_id=$1 and variant_id=$4 and lifecycle_state='active'`, p.TenantID, r.PurchaseOrderID, line.ID, line.VariantID, line.Quantity)
		if err != nil {
			return sc.Receipt{}, supplyError(err)
		}
		if tag.RowsAffected() != 1 {
			return sc.Receipt{}, sc.ErrConflict
		}
	}
	receipt, err := writeSupplyReceipt(ctx, tx, p.TenantID, r.PurchaseOrderID, r.CommandID, "planned", p.Subject, requestSHA, r.EvidenceSHA256, 1, map[string]any{"plan": r, "destination_organization_id": org, "supplier_id": supplier, "policy_sha256": hex.EncodeToString(policyHash[:])})
	if err != nil {
		return sc.Receipt{}, err
	}
	return receipt, nil
}
