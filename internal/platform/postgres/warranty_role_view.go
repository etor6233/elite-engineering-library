// AUTHORED bounded projection of already authorized immutable evidence, read in
// the Claim repeatable-read transaction. No write or new domain algorithm.
package postgres

import (
	"context"
	"elite.local/enterprise/internal/platform/identity"
	wc "elite.local/enterprise/internal/warrantyclaim"
	"encoding/json"
	"errors"
)

func readWarrantyRoleContext(ctx context.Context, q warrantyRowReader, tenant, id string) (*wc.RoleContext, error) {
	out := &wc.RoleContext{}
	opened, e := readWarrantyStep(ctx, q, tenant, id, "", "opened")
	if e != nil {
		return nil, e
	}
	var opening struct {
		Request wc.OpenClaim `json:"request"`
	}
	if json.Unmarshal(opened.Payload, &opening) != nil || opening.Request.CaseID != id {
		return nil, wc.ErrConflict
	}
	out.Description = opening.Request.Description
	out.Severity = opening.Request.Severity
	for _, kind := range []string{"diagnosed", "work", "quality", "accepted", "reconciled"} {
		step, e := readWarrantyStep(ctx, q, tenant, id, "", kind)
		if errors.Is(e, wc.ErrNotFound) {
			continue
		}
		if e != nil {
			return nil, e
		}
		switch kind {
		case "diagnosed":
			var v struct {
				Diagnosis wc.Diagnose `json:"diagnosis"`
				Excluded  bool        `json:"fault_excluded"`
			}
			if json.Unmarshal(step.Payload, &v) != nil || v.Diagnosis.CaseID != id {
				return nil, wc.ErrConflict
			}
			out.Diagnosis = &wc.RoleDiagnosis{FaultCode: v.Diagnosis.FaultCode, Description: v.Diagnosis.Description, Excluded: v.Excluded, EvidenceSHA256: v.Diagnosis.EvidenceSHA256}
		case "work":
			var v wc.WorkReceipt
			if json.Unmarshal(step.Payload, &v) != nil || v.Request.CaseID != id {
				return nil, wc.ErrConflict
			}
			out.Work = &wc.RoleWork{Actor: step.Actor, PayloadSHA256: step.PayloadSHA256, EvidenceSHA256: v.Request.EvidenceSHA256}
		case "quality":
			var v wc.Quality
			if json.Unmarshal(step.Payload, &v) != nil || v.CaseID != id {
				return nil, wc.ErrConflict
			}
			out.Quality = &wc.RoleQuality{Passed: v.Passed, PayloadSHA256: step.PayloadSHA256, EvidenceSHA256: v.EvidenceSHA256}
		case "accepted":
			var v wc.AcceptRepair
			if json.Unmarshal(step.Payload, &v) != nil || v.CaseID != id {
				return nil, wc.ErrConflict
			}
			out.Acceptance = &wc.RoleAcceptance{Actor: step.Actor, PayloadSHA256: step.PayloadSHA256}
		case "reconciled":
			var v wc.RoleReconciliation
			if json.Unmarshal(step.Payload, &v) != nil || v.Settlement != "INTERNAL_WARRANTY_SERVICE_ACKNOWLEDGEMENT" || v.RecordedInventoryCost == "" {
				return nil, wc.ErrConflict
			}
			out.Reconciliation = &v
		}
	}
	plan, approvalID, hash, e := readWarrantyPlan(ctx, q, tenant, id)
	if e == nil {
		var state string
		if e = q.QueryRow(ctx, `select state from approval.request where tenant_id=$1 and request_id=$2 and evidence_sha=$3`, tenant, approvalID, hash).Scan(&state); e != nil {
			return nil, e
		}
		out.Plan = &wc.RolePlan{Request: plan.Request, Requester: plan.Requester, PayloadSHA256: hash, Excluded: plan.Excluded, ExpiresAt: plan.ExpiresAt, ApprovalState: state}
	} else if !errors.Is(e, wc.ErrNotFound) {
		return nil, e
	}
	return out, nil
}

func (s *Warranty) QuoteForOffer(ctx context.Context, p identity.Principal, id string) (wc.QuoteForOffer, error) {
	var out wc.QuoteForOffer
	if !s.allowed(p, "warranty:offer") || !wc.ValidID(id) {
		return out, wc.ErrNotFound
	}
	tenant, org := s.profile.Scope()
	e := s.pool.QueryRow(ctx, `select quotation_id,organization_id,version,coalesce(customer_principal_id,''),state,valid_until>clock_timestamp(),currency,total_minor_units from sales.quotation where tenant_id=$1 and organization_id=$2 and quotation_id=$3`, tenant, org, id).Scan(&out.QuoteID, &out.OrganizationID, &out.QuoteVersion, &out.CustomerSubject, &out.State, &out.Current, &out.Currency, &out.TotalMinorUnits)
	return out, warrantyError(e)
}
