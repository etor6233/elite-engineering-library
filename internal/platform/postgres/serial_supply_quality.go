// AUTHORED quality payload binding. Decisions use the existing shared registry.
package postgres

import (
	"context"
	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/randomid"
	sc "elite.local/enterprise/internal/serialsupply"
	"github.com/jackc/pgx/v5"
	"strconv"
)

type supplyQualityState struct {
	ID, Hash, State, SourceState, Evidence string
	Version                                int64
	Attempt                                int
}

func readSupplyQuality(ctx context.Context, tx pgx.Tx, tenant, unit, stage string) (supplyQualityState, error) {
	var q supplyQualityState
	err := tx.QueryRow(ctx, `select q.approval_id,q.payload_sha256,a.state,q.source_state,q.source_version,q.attempt,coalesce(a.payload->>'basis_evidence_sha256','')
 from procurement.serial_supply_quality q join approval.request a on a.tenant_id=q.tenant_id and a.request_id=q.approval_id
 where q.tenant_id=$1 and q.production_unit_id=$2 and q.stage=$3 and a.kind='serial_quality' and a.evidence_sha=q.payload_sha256
 order by q.attempt desc limit 1 for update of a`, tenant, unit, stage).Scan(&q.ID, &q.Hash, &q.State, &q.SourceState, &q.Version, &q.Attempt, &q.Evidence)
	return q, supplyError(err)
}
func (s *SerialSupply) proposeSupplyQuality(ctx context.Context, tx pgx.Tx, p identity.Principal, plan sc.Plan, r sc.Command, u supplyUnitState, stage string) (string, error) {
	org, permission := supplyCommandPermission(r, plan)
	state, version := u.State, int64(0)
	if stage == "receipt" {
		state = u.StockState
		version = u.StockVersion
	}
	raw, hash, err := sc.Canonical(map[string]any{"purchase_order_id": plan.PurchaseOrderID, "unit_id": u.ID, "stock_unit_id": u.StockID, "serial_number": u.Serial, "variant_id": u.VariantID, "stage": stage, "source_state": state, "source_version": strconv.FormatInt(version, 10), "policy_code": sc.PolicyCode, "organization_id": org, "basis_evidence_sha256": r.EvidenceSHA256, "requester": p.Subject})
	if err != nil {
		return "", err
	}
	id := randomid.Generator{}.New()
	spec := HumanApprovalSpec{Request: approval.Request{TenantID: p.TenantID, ID: id, Kind: approval.KindSerialQuality, SubjectID: u.ID, Requester: p.Subject, EvidenceSHA: hash}, OrganizationID: org, Payload: raw}
	if _, err = NewHumanApprovals(s.pool).submitTx(ctx, tx, p, spec, permission, nil); err != nil {
		return "", err
	}
	_, err = tx.Exec(ctx, `insert into procurement.serial_supply_quality(tenant_id,purchase_order_id,production_unit_id,stage,attempt,approval_id,command_id,source_state,source_version,payload_sha256)
 select $1,$2,$3,$4,coalesce(max(attempt),0)+1,$5,$6,$7,$8,$9 from procurement.serial_supply_quality where tenant_id=$1 and production_unit_id=$3 and stage=$4`, p.TenantID, plan.PurchaseOrderID, u.ID, stage, id, r.CommandID, state, version, hash)
	return id, supplyError(err)
}
func (s *SerialSupply) decideSupplyQuality(ctx context.Context, tx pgx.Tx, p identity.Principal, plan sc.Plan, r sc.Command, u supplyUnitState, stage string, approved bool) (string, error) {
	q, err := readSupplyQuality(ctx, tx, p.TenantID, u.ID, stage)
	if err != nil {
		return "", err
	}
	state, version := u.State, int64(0)
	if stage == "receipt" {
		state = u.StockState
		version = u.StockVersion
	}
	if q.State != "pending" || q.SourceState != state || q.Version != version {
		return "", sc.ErrConflict
	}
	org, permission := supplyCommandPermission(r, plan)
	if _, err = NewHumanApprovals(s.pool).decideTx(ctx, tx, p, p.TenantID, q.ID, org, q.Hash, approved, r.Reason, permission, nil); err != nil {
		return "", err
	}
	return q.ID, nil
}
func (s *SerialSupply) supplyMilestone(ctx context.Context, tx pgx.Tx, p identity.Principal, plan sc.Plan, r sc.Command) (any, error) {
	if plan.State != "in-production" {
		return nil, sc.ErrConflict
	}
	u, err := readSupplyUnit(ctx, tx, p.TenantID, plan.PurchaseOrderID, r.UnitID)
	if err != nil {
		return nil, err
	}
	var approvalID string
	if u.State == "quality" && (r.TargetState == "released" || r.TargetState == "rejected") {
		approvalID, err = s.decideSupplyQuality(ctx, tx, p, plan, r, u, "factory", r.TargetState == "released")
		if err != nil {
			return nil, err
		}
	}
	if err = moveSupplyFactory(ctx, tx, p.TenantID, plan, r, u.ID, u.State, r.TargetState); err != nil {
		return nil, err
	}
	if r.TargetState == "quality" {
		u.State = "quality"
		approvalID, err = s.proposeSupplyQuality(ctx, tx, p, plan, r, u, "factory")
		if err != nil {
			return nil, err
		}
	}
	return map[string]any{"unit_id": u.ID, "to_state": r.TargetState, "approval_id": approvalID}, nil
}
func (s *SerialSupply) supplyReceiptQuality(ctx context.Context, tx pgx.Tx, p identity.Principal, plan sc.Plan, r sc.Command) (any, error) {
	u, err := readSupplyUnit(ctx, tx, p.TenantID, plan.PurchaseOrderID, r.UnitID)
	if err != nil {
		return nil, err
	}
	if u.State != "received" || u.StockState != "quarantine" {
		return nil, sc.ErrConflict
	}
	if r.Kind == "reinspect" {
		q, err := readSupplyQuality(ctx, tx, p.TenantID, u.ID, "receipt")
		if err != nil {
			return nil, err
		}
		if q.State != "rejected" || q.Evidence == r.EvidenceSHA256 {
			return nil, sc.ErrConflict
		}
		id, err := s.proposeSupplyQuality(ctx, tx, p, plan, r, u, "receipt")
		return map[string]any{"unit_id": u.ID, "stock_unit_id": u.StockID, "approval_id": id, "stock_state": "quarantine", "inspection_reason": r.Reason}, err
	}
	approved := r.Kind == "quality"
	id, err := s.decideSupplyQuality(ctx, tx, p, plan, r, u, "receipt", approved)
	if err != nil {
		return nil, err
	}
	target := "quarantine"
	if approved {
		target = "available"
		if err = moveSupplyStock(ctx, tx, p.TenantID, plan, r, u.StockID, u.StockState, target, u.StockVersion); err != nil {
			return nil, err
		}
	}
	return map[string]any{"unit_id": u.ID, "stock_unit_id": u.StockID, "approval_id": id, "approved": approved, "stock_state": target}, nil
}
