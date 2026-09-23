package postgres

// AUTHORED persistence/transaction glue around approval.Registry. The typed
// caller supplies its permission and same-transaction business binding guard.
import (
	"context"
	"encoding/json"
	"errors"

	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/platform/identity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HumanApprovalSpec struct {
	Request        approval.Request
	OrganizationID string
	Payload        json.RawMessage
}
type HumanApprovalGuard func(context.Context, pgx.Tx) error
type HumanApprovals struct{ pool *pgxpool.Pool }

func NewHumanApprovals(pool *pgxpool.Pool) *HumanApprovals { return &HumanApprovals{pool} }
func boundApprovalKind(k approval.Kind) bool {
	return k == approval.KindWhatsAppReply || k == approval.KindSocialPublish || k == approval.KindSocialRevoke || k == approval.KindStoredValueOperation || k == approval.KindWarrantyRepair || k == approval.KindSerialQuality || k == approval.KindCatalogReview || k == approval.KindTrainingAssessment || k == approval.KindMarketplaceMutation || k == approval.KindWhatsAppSchedule || k == approval.KindDocumentReview
}
func approvalPrincipal(p identity.Principal, tenant, org, permission string) bool {
	return permission != "" && p.Subject != "" && len(p.Subject) <= 128 && p.TenantID == tenant && p.Allowed(permission) && p.AllowedOrganization(org)
}
func (s *HumanApprovals) Submit(ctx context.Context, p identity.Principal, v HumanApprovalSpec, permission string, guard HumanApprovalGuard) (bool, error) {

	if s == nil || s.pool == nil {
		return false, approval.ErrInvalidRequest
	}
	tx, e := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if e != nil {
		return false, e
	}
	defer tx.Rollback(ctx)
	replay, e := s.submitTx(ctx, tx, p, v, permission, guard)
	if e != nil {
		return false, e
	}
	return replay, tx.Commit(ctx)
}

// submitTx reuses the same approval admission under a caller-owned local transaction.
// It does not commit, approve, or run external provider effects.
func (s *HumanApprovals) submitTx(ctx context.Context, tx pgx.Tx, p identity.Principal, v HumanApprovalSpec, permission string, guard HumanApprovalGuard) (bool, error) {
	canonical, hash, e := approval.CanonicalPayload(v.Payload)
	if s == nil || s.pool == nil || e != nil || v.Request.Validate() != nil || !boundApprovalKind(v.Request.Kind) || v.Request.AmountMinorUnits != 0 || v.Request.Requester != p.Subject || v.Request.EvidenceSHA != hash || !approvalPrincipal(p, v.Request.TenantID, v.OrganizationID, permission) {
		return false, approval.ErrInvalidRequest
	}
	if guard != nil {
		if e = guard(ctx, tx); e != nil {
			return false, e
		}
	}
	tag, e := tx.Exec(ctx, `insert into approval.request(tenant_id,request_id,kind,subject_id,amount_minor_units,requester,evidence_sha,organization_id,payload)
 values($1,$2,$3,$4,0,$5,$6,$7,$8) on conflict do nothing`, v.Request.TenantID, v.Request.ID, v.Request.Kind, v.Request.SubjectID, v.Request.Requester, hash, v.OrganizationID, canonical)
	if e != nil {
		return false, e
	}
	var same bool
	e = tx.QueryRow(ctx, `select kind=$3 and subject_id=$4 and amount_minor_units=0 and requester=$5 and evidence_sha=$6 and organization_id=$7 and payload=$8::jsonb from approval.request where tenant_id=$1 and request_id=$2`, v.Request.TenantID, v.Request.ID, v.Request.Kind, v.Request.SubjectID, v.Request.Requester, hash, v.OrganizationID, canonical).Scan(&same)
	if e != nil || !same {
		return false, approval.ErrDuplicate
	}
	return tag.RowsAffected() == 0, nil
}
func (s *HumanApprovals) Decide(ctx context.Context, p identity.Principal, tenant, id, org, expectedSHA string, approved bool, reason, permission string, guard HumanApprovalGuard) (approval.State, error) {
	if s == nil || s.pool == nil || len(reason) > 2048 || !approvalPrincipal(p, tenant, org, permission) {
		return "", approval.ErrInvalidRequest
	}
	tx, e := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if e != nil {
		return "", e
	}
	defer tx.Rollback(ctx)
	state, e := s.decideTx(ctx, tx, p, tenant, id, org, expectedSHA, approved, reason, permission, guard)
	if e != nil {
		return "", e
	}
	return state, tx.Commit(ctx)
}

// AUTHORED caller-controlled transaction composition; registry decision owner retained.
func (s *HumanApprovals) decideTx(ctx context.Context, tx pgx.Tx, p identity.Principal, tenant, id, org, expectedSHA string, approved bool, reason, permission string, guard HumanApprovalGuard) (approval.State, error) {
	if s == nil || s.pool == nil || len(reason) > 2048 || !approvalPrincipal(p, tenant, org, permission) {
		return "", approval.ErrInvalidRequest
	}
	v, state, e := readHumanApproval(ctx, tx, tenant, id, true)
	if e != nil {
		return "", e
	}
	if v.OrganizationID != org || v.Request.EvidenceSHA != expectedSHA || !boundApprovalKind(v.Request.Kind) {
		return "", approval.ErrInvalidRequest
	}
	if state != approval.StatePending {
		return state, approval.ErrNotPending
	}
	// Registry owns separation and the one-human decision transition. Zero policy
	// intentionally disables automatic/threshold/velocity decisions for this glue.
	registry := approval.NewRegistry(approval.Policy{})
	if _, e = registry.Submit(v.Request); e != nil {
		return "", e
	}
	if approved {
		state, e = registry.Approve(tenant, id, p.Subject, reason)
	} else {
		state, e = registry.Reject(tenant, id, p.Subject, reason)
	}
	if e != nil {
		return "", e
	}
	if guard != nil {
		if e = guard(ctx, tx); e != nil {
			return "", e
		}
	}
	_, e = tx.Exec(ctx, `insert into approval.decision(tenant_id,request_id,reviewer,approved,reason) values($1,$2,$3,$4,$5)`, tenant, id, p.Subject, approved, reason)
	if e != nil {
		return "", e
	}
	_, e = tx.Exec(ctx, `update approval.request set state=$3,decided_at=clock_timestamp() where tenant_id=$1 and request_id=$2 and state='pending'`, tenant, id, state)
	if e != nil {
		return "", e
	}
	return state, nil
}
func readHumanApproval(ctx context.Context, tx pgx.Tx, tenant, id string, lock bool) (HumanApprovalSpec, approval.State, error) {
	var v HumanApprovalSpec
	var state approval.State
	q := `select tenant_id::text,request_id,kind,subject_id,amount_minor_units,requester,evidence_sha,organization_id,payload,state from approval.request where tenant_id=$1 and request_id=$2`
	if lock {
		q += ` for update`
	}
	e := tx.QueryRow(ctx, q, tenant, id).Scan(&v.Request.TenantID, &v.Request.ID, &v.Request.Kind, &v.Request.SubjectID, &v.Request.AmountMinorUnits, &v.Request.Requester, &v.Request.EvidenceSHA, &v.OrganizationID, &v.Payload, &state)
	if errors.Is(e, pgx.ErrNoRows) {
		e = approval.ErrNotFound
	}
	return v, state, e
}

// LookupHumanApproval holds the request row through the caller's fence commit.
// It verifies durable approval only; domain/expiry/provider checks stay typed.
func LookupHumanApproval(ctx context.Context, tx pgx.Tx, tenant, id, org string, kind approval.Kind, expectedSHA string) (HumanApprovalSpec, error) {
	v, state, e := readHumanApproval(ctx, tx, tenant, id, true)
	if e != nil {
		return v, e
	}
	if state != approval.StateApproved || v.OrganizationID != org || v.Request.Kind != kind || v.Request.EvidenceSHA != expectedSHA {
		return HumanApprovalSpec{}, approval.ErrInvalidRequest
	}
	var reviewed bool
	e = tx.QueryRow(ctx, `select count(*)=1 and coalesce(bool_and(approved and reviewer<>$3),false) from approval.decision where tenant_id=$1 and request_id=$2`, tenant, id, v.Request.Requester).Scan(&reviewed)
	if e != nil {
		return HumanApprovalSpec{}, e
	}
	if !reviewed {
		return HumanApprovalSpec{}, approval.ErrSeparation
	}
	return v, nil
}
