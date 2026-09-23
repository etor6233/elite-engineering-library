// AUTHORED transaction and evidence composition. Existing service cases, human
// approval registry and inventory costing writers retain their ownership.
package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/fulfillment"
	"elite.local/enterprise/internal/inventorycontrol"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/randomid"
	wc "elite.local/enterprise/internal/warrantyclaim"
	"github.com/jackc/pgx/v5"
)

type warrantyRowReader interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}

func readWarrantyClaim(ctx context.Context, q warrantyRowReader, tenant, org, id string, lock bool) (wc.Claim, error) {
	var c wc.Claim
	sql := `select c.service_case_id,c.warranty_id,c.handover_id,c.appointment_id,c.organization_id,c.factory_organization_id,c.customer_principal_id,c.stock_unit_id,c.service_date::text,s.state,s.version,c.profile_sha256,c.parts_covered,c.labor_covered
 from service_ops.warranty_claim c join service_ops.service_case s using(tenant_id,service_case_id)
 where c.tenant_id=$1 and c.organization_id=$2 and c.service_case_id=$3`
	if lock {
		sql += " for update of s"
	}
	err := q.QueryRow(ctx, sql, tenant, org, id).Scan(&c.CaseID, &c.WarrantyID, &c.HandoverID, &c.AppointmentID, &c.OrganizationID, &c.FactoryOrganizationID, &c.CustomerSubject, &c.StockUnitID, &c.ServiceDate, &c.State, &c.Version, &c.ProfileSHA256, &c.PartsCovered, &c.LaborCovered)
	return c, warrantyError(err)
}
func readWarrantyStep(ctx context.Context, q warrantyRowReader, tenant, id, command, kind string) (wc.Step, error) {
	var s wc.Step
	err := q.QueryRow(ctx, `select service_case_id,command_id,version,kind,to_state,actor,request_sha256,payload_sha256,payload,recorded_at
 from service_ops.warranty_claim_step where tenant_id=$1 and service_case_id=$2 and ($3::text='' or command_id=$3) and ($4::text='' or kind=$4) order by version desc limit 1`, tenant, id, command, kind).Scan(&s.CaseID, &s.CommandID, &s.Version, &s.Kind, &s.State, &s.Actor, &s.RequestSHA256, &s.PayloadSHA256, &s.Payload, &s.RecordedAt)
	if err != nil {
		return s, warrantyError(err)
	}
	raw, hash, err := approval.CanonicalPayload(s.Payload)
	if err != nil || hash != s.PayloadSHA256 {
		return wc.Step{}, wc.ErrConflict
	}
	s.Payload = raw
	return s, nil
}
func (s *Warranty) claimAllowed(p identity.Principal, c wc.Claim, permission string) bool {
	tenant, org := s.profile.Scope()
	if c.OrganizationID != org {
		return false
	}
	switch permission {
	case "warranty:factory-read", "warranty:reconcile":
		return approvalPrincipal(p, tenant, c.FactoryOrganizationID, permission)
	case "warranty:self":
		return p.Subject == c.CustomerSubject && approvalPrincipal(p, tenant, org, permission)
	default:
		return approvalPrincipal(p, tenant, org, permission)
	}
}
func (s *Warranty) Claim(ctx context.Context, p identity.Principal, id string) (wc.Claim, error) {
	if s == nil || !s.profile.Valid() || !wc.ValidID(id) {
		return wc.Claim{}, wc.ErrInvalid
	}
	tenant, org := s.profile.Scope()
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly})
	if err != nil {
		return wc.Claim{}, err
	}
	defer tx.Rollback(ctx)
	c, err := readWarrantyClaim(ctx, tx, tenant, org, id, false)
	if err != nil {
		return wc.Claim{}, err
	}
	if !(s.claimAllowed(p, c, "warranty:read") || s.claimAllowed(p, c, "warranty:self") || s.claimAllowed(p, c, "warranty:factory-read")) {
		return wc.Claim{}, wc.ErrNotFound
	}
	latest, err := readWarrantyStep(ctx, tx, tenant, id, "", "")
	if err != nil {
		return wc.Claim{}, err
	}
	c.Latest = &latest
	c.Context, err = readWarrantyRoleContext(ctx, tx, tenant, id)
	if err != nil {
		return wc.Claim{}, err
	}
	return c, tx.Commit(ctx)
}
func (s *Warranty) ClaimCommand(ctx context.Context, p identity.Principal, id, command string) (wc.Step, error) {
	if !wc.ValidID(command) {
		return wc.Step{}, wc.ErrInvalid
	}
	if _, err := s.Claim(ctx, p, id); err != nil {
		return wc.Step{}, err
	}
	tenant, _ := s.profile.Scope()
	return readWarrantyStep(ctx, s.pool, tenant, id, command, "")
}
func claimSoldProfile(ctx context.Context, q warrantyRowReader, tenant string, c wc.Claim) (wc.Profile, error) {
	var raw []byte
	var hash string
	err := q.QueryRow(ctx, `select o.profile_bytes,o.profile_sha256 from service_ops.warranty_activation a
 join service_ops.warranty_offer o on o.tenant_id=a.tenant_id and o.quotation_id=a.quotation_id
 where a.tenant_id=$1 and a.warranty_id=$2 and a.organization_id=$3 and a.profile_sha256=$4`, tenant, c.WarrantyID, c.OrganizationID, c.ProfileSHA256).Scan(&raw, &hash)
	if err != nil {
		return wc.Profile{}, warrantyError(err)
	}
	p, err := wc.LoadProfile(raw, hash)
	if err != nil {
		return wc.Profile{}, err
	}
	pt, po := p.Scope()
	if pt != tenant || po != c.OrganizationID || p.Document().FactoryOrganizationID != c.FactoryOrganizationID {
		return wc.Profile{}, wc.ErrConflict
	}
	return p, nil
}
func writeWarrantyStep(ctx context.Context, tx pgx.Tx, tenant string, c wc.Claim, command, kind, target, actor, requestSHA string, payload any, opening bool) (wc.Step, error) {
	raw, hash, err := wc.Canonical(payload)
	if err != nil {
		return wc.Step{}, err
	}
	version := c.Version + 1
	if opening {
		version = 1
	}
	_, err = tx.Exec(ctx, `insert into service_ops.warranty_claim_step(tenant_id,service_case_id,command_id,version,kind,from_state,to_state,actor,request_sha256,payload_sha256,payload)
 values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`, tenant, c.CaseID, command, version, kind, c.State, target, actor, requestSHA, hash, raw)
	if err != nil {
		return wc.Step{}, warrantyError(err)
	}
	if !opening {
		tag, err := tx.Exec(ctx, `update service_ops.service_case set state=$4,version=version+1,closed_at=case when $4='closed' then clock_timestamp() else null end
 where tenant_id=$1 and service_case_id=$2 and version=$3`, tenant, c.CaseID, c.Version, target)
		if err != nil {
			return wc.Step{}, err
		}
		if tag.RowsAffected() != 1 {
			return wc.Step{}, wc.ErrConflict
		}
	}
	result, err := readWarrantyStep(ctx, tx, tenant, c.CaseID, command, "")
	if err != nil {
		return wc.Step{}, err
	}
	if err = warrantyEvent(ctx, tx, tenant, c.CaseID, "warranty.claim-"+kind, version, result); err != nil {
		return wc.Step{}, err
	}
	return result, nil
}

type warrantyClaimEffect func(context.Context, pgx.Tx, string, wc.Claim) (kind, target string, payload any, err error)

func (s *Warranty) claimCommand(ctx context.Context, p identity.Principal, c wc.Command, permission string, request any, effect warrantyClaimEffect) (wc.Step, error) {
	if s == nil || !s.profile.Valid() || !c.Valid() {
		return wc.Step{}, wc.ErrInvalid
	}
	_, requestSHA, err := wc.Canonical(request)
	if err != nil {
		return wc.Step{}, err
	}
	tenant, org := s.profile.Scope()
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return wc.Step{}, err
	}
	defer tx.Rollback(ctx)
	claim, err := readWarrantyClaim(ctx, tx, tenant, org, c.CaseID, true)
	if err != nil {
		return wc.Step{}, err
	}
	if !s.claimAllowed(p, claim, permission) {
		return wc.Step{}, wc.ErrNotFound
	}
	existing, err := readWarrantyStep(ctx, tx, tenant, c.CaseID, c.CommandID, "")
	if err == nil {
		if existing.Actor != p.Subject || existing.RequestSHA256 != requestSHA {
			return wc.Step{}, wc.ErrConflict
		}
		existing.Replay = true
		return existing, tx.Commit(ctx)
	}
	if !errors.Is(err, wc.ErrNotFound) {
		return wc.Step{}, err
	}
	if claim.Version != c.ExpectedVersion {
		return wc.Step{}, wc.ErrConflict
	}
	kind, target, payload, err := effect(ctx, tx, tenant, claim)
	if err != nil {
		return wc.Step{}, err
	}
	result, err := writeWarrantyStep(ctx, tx, tenant, claim, c.CommandID, kind, target, p.Subject, requestSHA, payload, false)
	if err != nil {
		return wc.Step{}, err
	}
	return result, tx.Commit(ctx)
}
func (s *Warranty) OpenClaim(ctx context.Context, p identity.Principal, r wc.OpenClaim) (wc.Step, error) {
	if !(s.allowed(p, "warranty:request") || s.allowed(p, "warranty:self")) || !wc.ValidID(r.CaseID) || !wc.ValidID(r.HandoverID) || !wc.ValidID(r.AppointmentID) || !wc.ValidSHA(r.EvidenceSHA256) || len(strings.TrimSpace(r.Description)) < 3 || len(r.Description) > 4000 || !map[string]bool{"low": true, "medium": true, "high": true, "safety": true}[r.Severity] {
		return wc.Step{}, wc.ErrInvalid
	}
	_, hash, err := wc.Canonical(r)
	if err != nil {
		return wc.Step{}, err
	}
	tenant, org := s.profile.Scope()
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return wc.Step{}, err
	}
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1,0))`, tenant+":warranty-case:"+r.CaseID); err != nil {
		return wc.Step{}, err
	}
	existing, err := readWarrantyStep(ctx, tx, tenant, r.CaseID, "open", "")
	if err == nil {
		if _, scopeErr := readWarrantyClaim(ctx, tx, tenant, org, r.CaseID, false); scopeErr != nil {
			return wc.Step{}, wc.ErrNotFound
		}
		if existing.Actor != p.Subject || existing.RequestSHA256 != hash {
			return wc.Step{}, wc.ErrConflict
		}
		existing.Replay = true
		return existing, tx.Commit(ctx)
	}
	if !errors.Is(err, wc.ErrNotFound) {
		return wc.Step{}, err
	}
	a, err := readWarrantyActivation(ctx, tx, tenant, org, r.HandoverID)
	if err != nil {
		return wc.Step{}, err
	}
	if !s.allowed(p, "warranty:request") && p.Subject != a.CustomerSubject {
		return wc.Step{}, wc.ErrNotFound
	}
	var warrantyStatus string
	if err = tx.QueryRow(ctx, `select status from service_ops.warranty where tenant_id=$1 and warranty_id=$2 for share`, tenant, a.WarrantyID).Scan(&warrantyStatus); err != nil || warrantyStatus == "void" {
		return wc.Step{}, wc.ErrConflict
	}
	// The appointment and unit must belong to the same actual customer/model.
	var valid bool
	err = tx.QueryRow(ctx, `select true from crm.appointment ap join inventory.stock_unit u on u.tenant_id=ap.tenant_id
 join catalog.vehicle_variant v on v.tenant_id=u.tenant_id and v.variant_id=u.variant_id
 where ap.tenant_id=$1 and ap.organization_id=$2 and ap.appointment_id=$3 and ap.customer_principal_id=$4
 and ap.appointment_kind='service' and ap.state in ('confirmed','completed') and u.stock_unit_id=$5
 and (ap.model_id is null or ap.model_id=v.model_id) for share of ap`, tenant, org, r.AppointmentID, a.CustomerSubject, a.StockUnitID).Scan(&valid)
	if err != nil || !valid {
		return wc.Step{}, wc.ErrConflict
	}
	var raw []byte
	err = tx.QueryRow(ctx, `select profile_bytes from service_ops.warranty_offer where tenant_id=$1 and quotation_id=$2`, tenant, a.QuoteID).Scan(&raw)
	if err != nil {
		return wc.Step{}, err
	}
	sold, err := wc.LoadProfile(raw, a.ProfileSHA256)
	if err != nil {
		return wc.Step{}, err
	}
	doc := sold.Document()
	zone, err := time.LoadLocation(doc.BusinessTimeZone)
	if err != nil {
		return wc.Step{}, err
	}
	var now time.Time
	if err = tx.QueryRow(ctx, `select clock_timestamp()`).Scan(&now); err != nil {
		return wc.Step{}, err
	}
	date := now.In(zone).Format(time.DateOnly)
	coverage, err := wc.CheckCoverage(date, a.Dates)
	if err != nil {
		return wc.Step{}, err
	}
	id := randomid.Generator{}
	if err = openServiceCaseInTx(ctx, tx, tenant, id.New(), fulfillment.ServiceCase{ID: r.CaseID, OrganizationID: org, StockUnitID: a.StockUnitID, Severity: r.Severity, Description: r.Description}); err != nil {
		return wc.Step{}, err
	}
	_, err = tx.Exec(ctx, `insert into service_ops.warranty_claim(tenant_id,service_case_id,warranty_id,handover_id,appointment_id,organization_id,factory_organization_id,customer_principal_id,stock_unit_id,service_date,profile_sha256,parts_covered,labor_covered,opened_by,request_sha256)
 values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10::date,$11,$12,$13,$14,$15)`, tenant, r.CaseID, a.WarrantyID, r.HandoverID, r.AppointmentID, org, doc.FactoryOrganizationID, a.CustomerSubject, a.StockUnitID, date, sold.Hash(), coverage.Parts, coverage.Labor, p.Subject, hash)
	if err != nil {
		return wc.Step{}, warrantyError(err)
	}
	claim, err := readWarrantyClaim(ctx, tx, tenant, org, r.CaseID, false)
	if err != nil {
		return wc.Step{}, err
	}
	step, err := writeWarrantyStep(ctx, tx, tenant, claim, "open", "opened", "opened", p.Subject, hash, map[string]any{"request": r, "coverage": coverage, "profile_sha256": sold.Hash(), "service_date": date}, true)
	if err != nil {
		return wc.Step{}, err
	}
	return step, tx.Commit(ctx)
}
func (s *Warranty) Diagnose(ctx context.Context, p identity.Principal, r wc.Diagnose) (wc.Step, error) {
	if !wc.ValidID(r.FaultCode) || !wc.ValidSHA(r.EvidenceSHA256) || len(strings.TrimSpace(r.Description)) < 3 || len(r.Description) > 4000 {
		return wc.Step{}, wc.ErrInvalid
	}
	return s.claimCommand(ctx, p, r.Command, "warranty:diagnose", r, func(ctx context.Context, tx pgx.Tx, tenant string, c wc.Claim) (string, string, any, error) {
		if c.State != "opened" {
			return "", "", nil, wc.ErrConflict
		}
		sold, err := claimSoldProfile(ctx, tx, tenant, c)
		if err != nil {
			return "", "", nil, err
		}
		return "diagnosed", "diagnosis", map[string]any{"diagnosis": r, "fault_excluded": sold.FaultExcluded(r.FaultCode)}, nil
	})
}
func readWarrantyPlan(ctx context.Context, q warrantyRowReader, tenant, id string) (wc.RepairPlan, string, string, error) {
	var plan wc.RepairPlan
	var raw []byte
	var approvalID, hash string
	err := q.QueryRow(ctx, `select approval_id,payload_sha256,payload from service_ops.warranty_work_plan where tenant_id=$1 and service_case_id=$2`, tenant, id).Scan(&approvalID, &hash, &raw)
	if err != nil {
		return plan, "", "", warrantyError(err)
	}
	_, got, err := approval.CanonicalPayload(raw)
	if err != nil || got != hash || json.Unmarshal(raw, &plan) != nil || plan.ClaimID != id {
		return plan, "", "", wc.ErrConflict
	}
	return plan, approvalID, hash, nil
}
func (s *Warranty) PlanRepair(ctx context.Context, p identity.Principal, r wc.Plan) (wc.Step, error) {
	if !r.Valid() {
		return wc.Step{}, wc.ErrInvalid
	}
	return s.claimCommand(ctx, p, r.Command, "warranty:plan", r, func(ctx context.Context, tx pgx.Tx, tenant string, c wc.Claim) (string, string, any, error) {
		if c.State != "diagnosis" {
			return "", "", nil, wc.ErrConflict
		}
		diagnosis, err := readWarrantyStep(ctx, tx, tenant, c.CaseID, "", "diagnosed")
		if err != nil {
			return "", "", nil, err
		}
		var d struct {
			FaultExcluded bool `json:"fault_excluded"`
		}
		if json.Unmarshal(diagnosis.Payload, &d) != nil {
			return "", "", nil, wc.ErrConflict
		}
		sold, err := claimSoldProfile(ctx, tx, tenant, c)
		if err != nil {
			return "", "", nil, err
		}
		// An uncovered diagnosis can be submitted for an explicit human rejection;
		// it never reserves parts or permits an approved covered-service claim.
		covered := !d.FaultExcluded && (c.PartsCovered || c.LaborCovered)
		if covered && (len(r.Parts) > 0 && !c.PartsCovered || r.LaborWork != "" && !c.LaborCovered) {
			return "", "", nil, wc.ErrConflict
		}
		var expires time.Time
		if err = tx.QueryRow(ctx, `select clock_timestamp()+make_interval(secs=>$1)`, sold.Document().WorkReservationSeconds).Scan(&expires); err != nil {
			return "", "", nil, err
		}
		plan := wc.RepairPlan{Schema: "elite-warranty-repair-plan/v1", ClaimID: c.CaseID, ProfileSHA256: c.ProfileSHA256, DiagnosisSHA256: diagnosis.PayloadSHA256, Requester: p.Subject, Request: r, PartsCovered: c.PartsCovered, LaborCovered: c.LaborCovered, Excluded: d.FaultExcluded, ExpiresAt: expires, Reservations: []inventorycontrol.BulkReservation{}}
		ids := randomid.Generator{}
		if covered {
			for _, part := range r.Parts {
				reservation, err := reserveBulkInTx(ctx, tx, tenant, ids.New(), inventorycontrol.BulkReservation{ID: wc.StableID("warranty-part", tenant, c.CaseID, part.LineID), OrganizationID: c.OrganizationID, BinID: part.BinID, ItemID: part.ItemID, LotID: part.LotID, DemandKind: "service", DemandID: c.CaseID, DemandLineID: part.LineID, Quantity: part.Quantity, Status: "reservation", Version: 1, ExpiresAt: &expires})
				if err != nil {
					return "", "", nil, err
				}
				plan.Reservations = append(plan.Reservations, reservation)
			}
		}
		raw, hash, err := wc.Canonical(plan)
		if err != nil {
			return "", "", nil, err
		}
		approvalID := wc.StableID("warranty-approval", tenant, c.CaseID)
		spec := HumanApprovalSpec{Request: approval.Request{TenantID: tenant, ID: approvalID, Kind: approval.KindWarrantyRepair, SubjectID: c.CaseID, Requester: p.Subject, EvidenceSHA: hash}, OrganizationID: c.OrganizationID, Payload: raw}
		if _, err = NewHumanApprovals(s.pool).submitTx(ctx, tx, p, spec, "warranty:plan", nil); err != nil {
			return "", "", nil, err
		}
		if _, err = tx.Exec(ctx, `insert into service_ops.warranty_work_plan(tenant_id,service_case_id,approval_id,payload_sha256,payload,expires_at) values($1,$2,$3,$4,$5,$6)`, tenant, c.CaseID, approvalID, hash, raw, expires); err != nil {
			return "", "", nil, err
		}
		for _, reservation := range plan.Reservations {
			if _, err = tx.Exec(ctx, `insert into service_ops.warranty_part_reservation(tenant_id,service_case_id,line_id,reservation_id) values($1,$2,$3,$4)`, tenant, c.CaseID, reservation.DemandLineID, reservation.ID); err != nil {
				return "", "", nil, err
			}
		}
		return "planned", "diagnosis", map[string]any{"approval_id": approvalID, "payload_sha256": hash, "plan": plan}, nil
	})
}
