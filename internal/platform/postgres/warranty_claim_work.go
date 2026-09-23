package postgres

// AUTHORED same-transaction service/approval/inventory composition. Repair
// acknowledgements never represent a paid reimbursement or a tax document.
import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/inventorycontrol"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/randomid"
	wc "elite.local/enterprise/internal/warrantyclaim"
	"github.com/jackc/pgx/v5"
)

func (s *Warranty) DecideRepair(ctx context.Context, p identity.Principal, r wc.DecideRepair) (wc.Step, error) {
	if !wc.ValidSHA(r.PayloadSHA256) || len(strings.TrimSpace(r.Reason)) < 3 || len(r.Reason) > 2048 {
		return wc.Step{}, wc.ErrInvalid
	}
	return s.claimCommand(ctx, p, r.Command, "warranty:approve", r, func(ctx context.Context, tx pgx.Tx, tenant string, c wc.Claim) (string, string, any, error) {
		if c.State != "diagnosis" {
			return "", "", nil, wc.ErrConflict
		}
		plan, id, hash, err := readWarrantyPlan(ctx, tx, tenant, c.CaseID)
		if err != nil {
			return "", "", nil, err
		}
		if hash != r.PayloadSHA256 || plan.ProfileSHA256 != c.ProfileSHA256 {
			return "", "", nil, wc.ErrConflict
		}
		if r.Approved && (plan.Excluded || !plan.PartsCovered && !plan.LaborCovered) {
			return "", "", nil, wc.ErrConflict
		}
		if _, err = NewHumanApprovals(s.pool).decideTx(ctx, tx, p, tenant, id, c.OrganizationID, hash, r.Approved, r.Reason, "warranty:approve", nil); err != nil {
			return "", "", nil, err
		}
		ids := randomid.Generator{}
		for _, reservation := range plan.Reservations {
			if r.Approved {
				// Approval converts the pending lease to an explicit durable repair hold.
				// Only metadata/version changes; the existing stock reservation remains.
				tag, err := tx.Exec(ctx, `update inventory.bulk_reservation set expires_at=null,version=version+1,updated_at=clock_timestamp()
 where tenant_id=$1 and organization_id=$2 and reservation_id=$3 and version=$4 and status='reservation' and expires_at=$5 and expires_at>clock_timestamp()`, tenant, c.OrganizationID, reservation.ID, reservation.Version, plan.ExpiresAt)
				if err != nil {
					return "", "", nil, err
				}
				if tag.RowsAffected() != 1 {
					return "", "", nil, wc.ErrConflict
				}
			} else if err = releaseBulkInTx(ctx, tx, tenant, c.OrganizationID, reservation.ID, reservation.Version, ids.New()); err != nil {
				return "", "", nil, err
			}
		}
		target := "cancelled"
		if r.Approved {
			target = "repair"
		}
		return "decision", target, map[string]any{"decision": r, "approval_id": id, "approved_sha256": hash, "requester": plan.Requester, "reviewer": p.Subject}, nil
	})
}
func (s *Warranty) CompleteRepairWork(ctx context.Context, p identity.Principal, r wc.CompleteWork) (wc.Step, error) {
	if !wc.ValidSHA(r.EvidenceSHA256) {
		return wc.Step{}, wc.ErrInvalid
	}
	return s.claimCommand(ctx, p, r.Command, "warranty:work", r, func(ctx context.Context, tx pgx.Tx, tenant string, c wc.Claim) (string, string, any, error) {
		if c.State != "repair" {
			return "", "", nil, wc.ErrConflict
		}
		plan, id, hash, err := readWarrantyPlan(ctx, tx, tenant, c.CaseID)
		if err != nil {
			return "", "", nil, err
		}
		if _, err = LookupHumanApproval(ctx, tx, tenant, id, c.OrganizationID, approval.KindWarrantyRepair, hash); err != nil {
			return "", "", nil, err
		}
		if plan.Excluded || !plan.PartsCovered && !plan.LaborCovered || len(plan.Reservations) != len(plan.Request.Parts) {
			return "", "", nil, wc.ErrConflict
		}
		receipt := wc.WorkReceipt{Request: r, ApprovalID: id, ApprovedSHA256: hash, Parts: []wc.IssuedPart{}}
		sold, err := claimSoldProfile(ctx, tx, tenant, c)
		if err != nil {
			return "", "", nil, err
		}
		zone, err := time.LoadLocation(sold.Document().BusinessTimeZone)
		if err != nil {
			return "", "", nil, err
		}
		var now time.Time
		if err = tx.QueryRow(ctx, `select clock_timestamp()`).Scan(&now); err != nil {
			return "", "", nil, err
		}
		local := now.In(zone)
		posting := time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, time.UTC)
		ids := randomid.Generator{}
		for i, part := range plan.Request.Parts {
			reservation := plan.Reservations[i]
			if reservation.DemandLineID != part.LineID || reservation.ItemID != part.ItemID || reservation.BinID != part.BinID {
				return "", "", nil, wc.ErrConflict
			}
			result, err := issueBulkInTx(ctx, tx, tenant, ids.New(), ids.New(), inventorycontrol.BulkIssue{OrganizationID: c.OrganizationID, ReservationID: reservation.ID, ReservationVersion: reservation.Version + 1, SpecificReceiptEntry: part.SpecificReceiptEntry, PostingDate: posting, SourceKind: "warranty-repair", SourceID: c.CaseID})
			if err != nil {
				return "", "", nil, err
			}
			receipt.Parts = append(receipt.Parts, wc.IssuedPart{LineID: part.LineID, ReservationID: reservation.ID, Issue: result})
		}
		return "work", "quality", receipt, nil
	})
}
func (s *Warranty) RecordRepairQuality(ctx context.Context, p identity.Principal, r wc.Quality) (wc.Step, error) {
	if !wc.ValidSHA(r.EvidenceSHA256) || !wc.ValidSHA(r.WorkEvidenceSHA256) || r.CorrectionEvidenceSHA256 != "" && !wc.ValidSHA(r.CorrectionEvidenceSHA256) {
		return wc.Step{}, wc.ErrInvalid
	}
	return s.claimCommand(ctx, p, r.Command, "warranty:quality", r, func(ctx context.Context, tx pgx.Tx, tenant string, c wc.Claim) (string, string, any, error) {
		if c.State != "quality" {
			return "", "", nil, wc.ErrConflict
		}
		work, err := readWarrantyStep(ctx, tx, tenant, c.CaseID, "", "work")
		if err != nil {
			return "", "", nil, err
		}
		if work.PayloadSHA256 != r.WorkEvidenceSHA256 || work.Actor == p.Subject {
			return "", "", nil, wc.ErrConflict
		}
		if _, err = readWarrantyStep(ctx, tx, tenant, c.CaseID, "", "accepted"); err == nil {
			return "", "", nil, wc.ErrConflict
		} else if !errors.Is(err, wc.ErrNotFound) {
			return "", "", nil, err
		}
		previous, err := readWarrantyStep(ctx, tx, tenant, c.CaseID, "", "quality")
		if err == nil {
			var q wc.Quality
			if json.Unmarshal(previous.Payload, &q) != nil || q.Passed || !wc.ValidSHA(r.CorrectionEvidenceSHA256) {
				return "", "", nil, wc.ErrConflict
			}
		} else if !errors.Is(err, wc.ErrNotFound) {
			return "", "", nil, err
		}
		return "quality", "quality", r, nil
	})
}
func (s *Warranty) AcceptRepair(ctx context.Context, p identity.Principal, r wc.AcceptRepair) (wc.Step, error) {
	if !wc.ValidSHA(r.QualitySHA256) || !wc.ValidSHA(r.EvidenceSHA256) {
		return wc.Step{}, wc.ErrInvalid
	}
	return s.claimCommand(ctx, p, r.Command, "warranty:self", r, func(ctx context.Context, tx pgx.Tx, tenant string, c wc.Claim) (string, string, any, error) {
		if c.State != "quality" {
			return "", "", nil, wc.ErrConflict
		}
		quality, err := readWarrantyStep(ctx, tx, tenant, c.CaseID, "", "quality")
		if err != nil {
			return "", "", nil, err
		}
		var q wc.Quality
		if json.Unmarshal(quality.Payload, &q) != nil || !q.Passed || quality.PayloadSHA256 != r.QualitySHA256 {
			return "", "", nil, wc.ErrConflict
		}
		return "accepted", "quality", r, nil
	})
}
func (s *Warranty) ReconcileRepair(ctx context.Context, p identity.Principal, r wc.ReconcileRepair) (wc.Step, error) {
	if !wc.ValidSHA(r.AcceptanceSHA256) || !wc.ValidSHA(r.EvidenceSHA256) {
		return wc.Step{}, wc.ErrInvalid
	}
	return s.claimCommand(ctx, p, r.Command, "warranty:reconcile", r, func(ctx context.Context, tx pgx.Tx, tenant string, c wc.Claim) (string, string, any, error) {
		if c.State != "quality" {
			return "", "", nil, wc.ErrConflict
		}
		accepted, err := readWarrantyStep(ctx, tx, tenant, c.CaseID, "", "accepted")
		if err != nil {
			return "", "", nil, err
		}
		if accepted.PayloadSHA256 != r.AcceptanceSHA256 || accepted.Actor != c.CustomerSubject {
			return "", "", nil, wc.ErrConflict
		}
		quality, err := readWarrantyStep(ctx, tx, tenant, c.CaseID, "", "quality")
		if err != nil {
			return "", "", nil, err
		}
		var ack wc.AcceptRepair
		var q wc.Quality
		if json.Unmarshal(accepted.Payload, &ack) != nil || json.Unmarshal(quality.Payload, &q) != nil || !q.Passed || ack.QualitySHA256 != quality.PayloadSHA256 {
			return "", "", nil, wc.ErrConflict
		}
		work, err := readWarrantyStep(ctx, tx, tenant, c.CaseID, "", "work")
		if err != nil {
			return "", "", nil, err
		}
		var performed wc.WorkReceipt
		if json.Unmarshal(work.Payload, &performed) != nil || q.WorkEvidenceSHA256 != work.PayloadSHA256 {
			return "", "", nil, wc.ErrConflict
		}
		plan, id, hash, err := readWarrantyPlan(ctx, tx, tenant, c.CaseID)
		if err != nil {
			return "", "", nil, err
		}
		if _, err = LookupHumanApproval(ctx, tx, tenant, id, c.OrganizationID, approval.KindWarrantyRepair, hash); err != nil {
			return "", "", nil, err
		}
		if performed.ApprovalID != id || performed.ApprovedSHA256 != hash || len(performed.Parts) != len(plan.Reservations) {
			return "", "", nil, wc.ErrConflict
		}
		entries := []string{}
		for i, part := range performed.Parts {
			if part.LineID != plan.Request.Parts[i].LineID || part.ReservationID != plan.Reservations[i].ID {
				return "", "", nil, wc.ErrConflict
			}
			var valid bool
			err = tx.QueryRow(ctx, `select exists(select 1 from inventory.bulk_inventory_entry e join inventory.bulk_reservation b on b.tenant_id=e.tenant_id and b.reservation_id=$3
 where e.tenant_id=$1 and e.entry_id=$2 and e.organization_id=$4 and e.entry_type='issue' and e.source_kind='warranty-repair' and e.source_id=$5
 and e.item_id=b.item_id and e.bin_id=b.bin_id and e.lot_id is not distinct from b.lot_id and e.quantity=-b.quantity
 and -e.cost_amount=$6::numeric and b.status='consumed' and b.demand_kind='service' and b.demand_id=$5 and b.demand_line_id=$7)`, tenant, part.Issue.EntryID, part.ReservationID, c.OrganizationID, c.CaseID, part.Issue.CostAmount, part.LineID).Scan(&valid)
			if err != nil {
				return "", "", nil, err
			}
			if !valid {
				return "", "", nil, wc.ErrConflict
			}
			entries = append(entries, part.Issue.EntryID)
		}
		var recordedCost string
		err = tx.QueryRow(ctx, `select coalesce(sum(-cost_amount),0)::numeric(38,4)::text from inventory.bulk_inventory_entry where tenant_id=$1 and entry_id=any($2::text[])`, tenant, entries).Scan(&recordedCost)
		if err != nil {
			return "", "", nil, err
		}
		sold, err := claimSoldProfile(ctx, tx, tenant, c)
		if err != nil {
			return "", "", nil, err
		}
		return "reconciled", "closed", map[string]any{"request": r, "settlement": sold.Document().Settlement, "factory_organization_id": c.FactoryOrganizationID, "work_sha256": work.PayloadSHA256, "quality_sha256": quality.PayloadSHA256, "acceptance_sha256": accepted.PayloadSHA256, "approved_sha256": hash, "recorded_inventory_cost": recordedCost, "inventory_entry_ids": entries, "payment_created": false}, nil
	})
}
func (s *Warranty) CancelRepair(ctx context.Context, p identity.Principal, r wc.CancelRepair) (wc.Step, error) {
	if !wc.ValidSHA(r.EvidenceSHA256) || len(strings.TrimSpace(r.Reason)) < 3 || len(r.Reason) > 2048 {
		return wc.Step{}, wc.ErrInvalid
	}
	return s.claimCommand(ctx, p, r.Command, "warranty:cancel", r, func(ctx context.Context, tx pgx.Tx, tenant string, c wc.Claim) (string, string, any, error) {
		if c.State == "closed" || c.State == "cancelled" || c.State == "quality" {
			return "", "", nil, wc.ErrConflict
		}
		if _, err := readWarrantyStep(ctx, tx, tenant, c.CaseID, "", "work"); err == nil {
			return "", "", nil, wc.ErrConflict
		} else if !errors.Is(err, wc.ErrNotFound) {
			return "", "", nil, err
		}
		plan, id, hash, err := readWarrantyPlan(ctx, tx, tenant, c.CaseID)
		if err == nil {
			_, state, err := readHumanApproval(ctx, tx, tenant, id, true)
			if err != nil {
				return "", "", nil, err
			}
			versionDelta := int64(0)
			if state == approval.StatePending {
				if _, err = NewHumanApprovals(s.pool).decideTx(ctx, tx, p, tenant, id, c.OrganizationID, hash, false, r.Reason, "warranty:cancel", nil); err != nil {
					return "", "", nil, err
				}
			} else if state == approval.StateApproved {
				versionDelta = 1
			} else {
				return "", "", nil, wc.ErrConflict
			}
			ids := randomid.Generator{}
			for _, reservation := range plan.Reservations {
				if err = releaseBulkInTx(ctx, tx, tenant, c.OrganizationID, reservation.ID, reservation.Version+versionDelta, ids.New()); err != nil {
					return "", "", nil, err
				}
			}
		} else if !errors.Is(err, wc.ErrNotFound) {
			return "", "", nil, err
		}
		return "cancelled", "cancelled", map[string]any{"request": r, "approval_id": id, "approved_sha256": hash, "inventory_issue_reversed": false}, nil
	})
}
