package postgres

// AUTHORED transaction/binding glue. Commercial points and reward arithmetic
// are executed by the separately licensed, hash-locked Odoo derivation.
import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"sort"
	"time"

	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/bcamounts"
	"elite.local/enterprise/internal/platform/identity"
	sv "elite.local/enterprise/internal/storedvaluebridge"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StoredValue struct {
	pool      *pgxpool.Pool
	profile   *sv.Profile
	process   sv.Process
	approvals *HumanApprovals
}

func NewStoredValue(pool *pgxpool.Pool, profile *sv.Profile, process sv.Process) (*StoredValue, error) {
	if pool == nil || !profile.Valid() || process.Validate() != nil {
		return nil, sv.ErrBinding
	}
	return &StoredValue{pool, profile, process, NewHumanApprovals(pool)}, nil
}
func (s *StoredValue) allowed(p identity.Principal, permission string) bool {
	if s == nil || s.pool == nil || !s.profile.Valid() {
		return false
	}
	tenant, org := s.profile.Scope()
	return approvalPrincipal(p, tenant, org, permission)
}
func (s *StoredValue) ProfileSHA256() string {
	if s == nil {
		return ""
	}
	return s.profile.SHA256()
}
func points(v string) (*big.Rat, error) {
	n, e := sv.Minor(v, 6)
	if e != nil {
		return nil, e
	}
	return new(big.Rat).SetFrac(big.NewInt(n), big.NewInt(1000000)), nil
}
func decimalPoints(v *big.Rat) string { return v.FloatString(6) }

// ApprovalView checks tenant, organization and the complete approved payload.
func (s *StoredValue) Read(ctx context.Context, p identity.Principal, id string) (sv.ApprovalView, error) {
	var out sv.ApprovalView
	if !s.allowed(p, "stored_value:read") || !sv.ValidID(id) {
		return out, sv.ErrBinding
	}
	tx, e := s.pool.Begin(ctx)
	if e != nil {
		return out, e
	}
	defer tx.Rollback(ctx)
	out, e = s.readTx(ctx, tx, p, id)
	if e != nil {
		return out, e
	}
	return out, tx.Commit(ctx)
}
func (s *StoredValue) readTx(ctx context.Context, tx pgx.Tx, p identity.Principal, id string) (sv.ApprovalView, error) {
	var out sv.ApprovalView
	tenant, org := s.profile.Scope()
	a, state, e := readHumanApproval(ctx, tx, tenant, id, false)
	if e != nil {
		return out, e
	}
	if a.OrganizationID != org || a.Request.Kind != approval.KindStoredValueOperation || sv.Decode(a.Payload, &out.Payload) != nil {
		return out, sv.ErrBinding
	}
	_, hash, e := sv.Canonical(out.Payload)
	if e != nil || hash != a.Request.EvidenceSHA {
		return out, sv.ErrBinding
	}
	out.ApprovalID = id
	out.PayloadSHA256 = hash
	out.State = string(state)
	out.Receipt = json.RawMessage("null")
	var storedSHA string
	e = tx.QueryRow(ctx, `select receipt,receipt_sha256 from stored_value.operation where tenant_id=$1 and approval_id=$2`, tenant, id).Scan(&out.Receipt, &storedSHA)
	if e == nil {
		raw, hash, err := approval.CanonicalPayload(out.Receipt)
		if err != nil || hash != storedSHA {
			return out, sv.ErrBinding
		}
		out.Receipt = raw
	}
	if errors.Is(e, pgx.ErrNoRows) {
		e = nil
	}
	return out, e
}
func (s *StoredValue) Propose(ctx context.Context, p identity.Principal, r sv.Request) (sv.ApprovalView, error) {
	var out sv.ApprovalView
	if !s.allowed(p, "stored_value:request") || r.Validate(s.profile) != nil {
		return out, sv.ErrBinding
	}
	tenant, _ := s.profile.Scope()
	id := sv.StableID("svapproval", r.OperationID)
	tx, e := s.pool.Begin(ctx)
	if e != nil {
		return out, e
	}
	defer tx.Rollback(ctx)
	// Serializes same request before order/account locks; an exact replay does not
	// consult a later order revision or regenerate the immutable lease/payload.
	if _, e = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1,0))`, tenant+":stored-value-request:"+r.OperationID); e != nil {
		return out, e
	}
	existing, e := s.readTx(ctx, tx, p, id)
	if e == nil {
		a, _, _ := sv.Canonical(r)
		b, _, _ := sv.Canonical(existing.Payload.Request)
		if string(a) != string(b) || existing.Payload.Requester != p.Subject {
			return out, approval.ErrDuplicate
		}
		existing.Replay = true
		return existing, tx.Commit(ctx)
	}
	if !errors.Is(e, approval.ErrNotFound) {
		return out, e
	}
	var expires time.Time
	if e = tx.QueryRow(ctx, `select date_trunc('microseconds',clock_timestamp())+make_interval(secs=>$1)`, s.profile.ReservationSeconds()).Scan(&expires); e != nil {
		return out, e
	}
	bound, e := s.calculate(ctx, tx, r, p.Subject, expires.UTC().Format(time.RFC3339Nano), false)
	if e != nil {
		return out, e
	}
	raw, hash, e := sv.Canonical(bound)
	if e != nil {
		return out, e
	}
	spec := HumanApprovalSpec{Request: approval.Request{TenantID: tenant, ID: id, Kind: approval.KindStoredValueOperation, SubjectID: r.OrderID, Requester: p.Subject, EvidenceSHA: hash}, OrganizationID: r.OrganizationID, Payload: raw}
	if _, e = s.approvals.submitTx(ctx, tx, p, spec, "stored_value:request", nil); e != nil {
		return out, e
	}
	if r.Operation == "redeem" {
		for _, entry := range bound.Entries {
			delta, e := points(entry.PointsDelta)
			if e != nil {
				return out, e
			}
			delta.Neg(delta)
			_, e = tx.Exec(ctx, `insert into stored_value.reservation(tenant_id,operation_id,account_id,order_id,approval_id,points,state,expires_at) values($1,$2,$3,$4,$5,$6::numeric,'reserved',$7)`, tenant, r.OperationID, entry.AccountID, r.OrderID, id, decimalPoints(delta), expires)
			if e != nil {
				return out, e
			}
		}
	}
	out = sv.ApprovalView{ApprovalID: id, PayloadSHA256: hash, State: string(approval.StatePending), Payload: bound, Receipt: json.RawMessage("null")}
	return out, tx.Commit(ctx)
}

func (s *StoredValue) Decide(ctx context.Context, p identity.Principal, d sv.Decision) (sv.ApprovalView, error) {
	var out sv.ApprovalView
	if !s.allowed(p, "stored_value:approve") || !sv.ValidID(d.ApprovalID) {
		return out, sv.ErrBinding
	}
	tenant, org := s.profile.Scope()
	if d.OrganizationID != org {
		return out, sv.ErrBinding
	}
	_, e := s.approvals.Decide(ctx, p, tenant, d.ApprovalID, org, d.PayloadSHA256, d.Approved, d.Reason, "stored_value:approve", func(ctx context.Context, tx pgx.Tx) error {
		a, _, e := readHumanApproval(ctx, tx, tenant, d.ApprovalID, false)
		if e != nil {
			return e
		}
		var approved sv.BoundOperation
		if a.Request.Kind != approval.KindStoredValueOperation || sv.Decode(a.Payload, &approved) != nil {
			return sv.ErrBinding
		}
		if !d.Approved {
			_, e = tx.Exec(ctx, `update stored_value.reservation set state='rejected' where tenant_id=$1 and approval_id=$2 and state='reserved'`, tenant, d.ApprovalID)
			return e
		}
		actual, e := s.calculate(ctx, tx, approved.Request, approved.Requester, approved.ExpiresAt, true)
		if e != nil {
			return e
		}
		_, actualSHA, e := sv.Canonical(actual)
		if e != nil || actualSHA != d.PayloadSHA256 {
			return sv.ErrBinding
		}
		return s.commit(ctx, tx, d.ApprovalID, d.PayloadSHA256, p.Subject, actual)
	})
	if e != nil {
		return out, e
	}
	// The authorization for this response is the decision permission, not a
	// separate read grant that could hide an already committed local result.
	tx, e := s.pool.Begin(ctx)
	if e != nil {
		return out, e
	}
	defer tx.Rollback(ctx)
	out, e = s.readTx(ctx, tx, p, d.ApprovalID)
	if e != nil {
		return out, e
	}
	return out, tx.Commit(ctx)
}

func (s *StoredValue) order(ctx context.Context, tx pgx.Tx, r sv.Request) (sv.Order, string, sv.Allocation, error) {
	var o sv.Order
	var state string
	var a sv.Allocation
	tenant, _ := s.profile.Scope()
	var subject *string
	e := tx.QueryRow(ctx, `select o.order_id,o.organization_id,o.customer_principal_id,o.state,o.currency,o.version,o.total_minor_units from sales.customer_order o where tenant_id=$1 and order_id=$2 and organization_id=$3 for update`, tenant, r.OrderID, r.OrganizationID).Scan(&o.OrderID, &o.OrganizationID, &subject, &state, &o.Currency, &a.OrderVersion, &a.GrossMinor)
	if e != nil {
		return o, state, a, e
	}
	if a.OrderVersion != r.ExpectedOrderVersion || subject == nil || !sv.ValidID(*subject) {
		return o, state, a, sv.ErrBinding
	}
	o.SubjectID = *subject
	o.State = "draft"
	o.EnabledRuleIDs = []string{}
	o.Lines = []sv.Line{}
	e = tx.QueryRow(ctx, `select gift_minor_units,discount_minor_units,provider_due_minor_units from stored_value.order_allocation where tenant_id=$1 and order_id=$2`, tenant, r.OrderID).Scan(&a.GiftMinor, &a.DiscountMinor, &a.ProviderDueMinor)
	if e != nil {
		return o, state, a, e
	}
	if a.GiftMinor < 0 || a.DiscountMinor < 0 || a.ProviderDueMinor < 0 {
		return o, state, a, sv.ErrBinding
	}
	a.OrderID = o.OrderID
	a.Currency = o.Currency
	rows, e := tx.Query(ctx, `select line_id,variant_id,quantity,unit_price_minor_units from sales.customer_order_line where tenant_id=$1 and order_id=$2 order by line_id`, tenant, r.OrderID)
	if e != nil {
		return o, state, a, e
	}
	total := new(big.Int)
	for rows.Next() {
		var l sv.Line
		var qty, price int64
		if e = rows.Scan(&l.ID, &l.ProductID, &qty, &price); e != nil {
			rows.Close()
			return o, state, a, e
		}
		amount, e := bcamounts.LineAmount(qty, price, 0)
		if e != nil {
			rows.Close()
			return o, state, a, e
		}
		total.Add(total, big.NewInt(amount))
		l.Quantity = sv.Number(qty)
		l.Total = sv.Major(amount, 2)
		l.Subtotal = l.Total
		l.Tax = "0"
		o.Lines = append(o.Lines, l)
	}
	e = rows.Err()
	rows.Close()
	if e != nil {
		return o, state, a, e
	}
	if len(o.Lines) < 1 || len(o.Lines) > 40 || !total.IsInt64() || total.Int64() != a.GrossMinor {
		return o, state, a, sv.ErrBinding
	}
	// Upstream reward lines are represented separately from immutable gross
	// sales lines. These typed negative lines preserve source points filtering.
	rows, e = tx.Query(ctx, `select op.operation_id,op.program_id,ac.kind,sum(en.applied_minor_units)::bigint from stored_value.operation op join stored_value.entry en using(tenant_id,operation_id) join stored_value.account ac using(tenant_id,account_id) where op.tenant_id=$1 and op.order_id=$2 group by op.operation_id,op.program_id,ac.kind order by op.operation_id`, tenant, r.OrderID)
	if e != nil {
		return o, state, a, e
	}
	for rows.Next() {
		var id, program, kind string
		var amount int64
		if e = rows.Scan(&id, &program, &kind, &amount); e != nil {
			rows.Close()
			return o, state, a, e
		}
		if amount == 0 {
			continue
		}
		val := sv.Major(-amount, 2)
		o.Lines = append(o.Lines, sv.Line{ID: id, ProductID: "stored-value-reward", Quantity: "1", Subtotal: val, Tax: "0", Total: val, RewardProgramType: kind, RewardProgramID: program, RewardTrigger: "auto"})
	}
	e = rows.Err()
	rows.Close()
	o.Total = sv.Major(a.ProviderDueMinor, 2)
	return o, state, a, e
}

func (s *StoredValue) account(ctx context.Context, tx pgx.Tx, r sv.Request, subject, id string) (string, error) {
	tenant, _ := s.profile.Scope()
	if _, e := tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1,0))`, tenant+":stored-value-account:"+id); e != nil {
		return "", e
	}
	var kind, currency, program, org string
	var owner *string
	e := tx.QueryRow(ctx, `select kind,currency,program_id,organization_id,customer_subject from stored_value.account where tenant_id=$1 and account_id=$2 for update`, tenant, id).Scan(&kind, &currency, &program, &org, &owner)
	if e != nil {
		return "", e
	}
	pp, e := s.profile.Program(r.ProgramID)
	if e != nil || kind != pp.Calculation.Kind || currency != pp.Calculation.Currency || program != r.ProgramID || org != r.OrganizationID || (owner != nil && *owner != subject) {
		return "", sv.ErrBinding
	}
	var balance string
	e = tx.QueryRow(ctx, `select (coalesce((select sum(points_delta) from stored_value.entry where tenant_id=$1 and account_id=$2),0)-coalesce((select sum(points) from stored_value.reservation where tenant_id=$1 and account_id=$2 and operation_id<>$3 and state='reserved' and expires_at>clock_timestamp()),0))::numeric(38,6)::text`, tenant, id, r.OperationID).Scan(&balance)
	return balance, e
}

func (s *StoredValue) calculate(ctx context.Context, tx pgx.Tx, r sv.Request, requester, expiry string, deciding bool) (sv.BoundOperation, error) {
	out := sv.BoundOperation{Schema: "elite.stored-value-operation.v1", Request: r, Requester: requester, EngineSHA256: s.process.ManifestSHA256, ExpiresAt: expiry, Entries: []sv.Entry{}}
	if r.Validate(s.profile) != nil {
		return out, sv.ErrBinding
	}
	tenant, _ := s.profile.Scope()
	pp, e := s.profile.Program(r.ProgramID)
	if e != nil {
		return out, e
	}
	o, state, a, e := s.order(ctx, tx, r)
	if e != nil {
		return out, e
	}
	out.Order = o
	out.OrderState = state
	out.GrossMinor = a.GrossMinor
	out.GiftMinor = a.GiftMinor
	out.DiscountMinor = a.DiscountMinor
	if o.Currency != pp.Calculation.Currency {
		return out, sv.ErrBinding
	}
	var validTime bool
	if e = tx.QueryRow(ctx, `select $1::timestamptz>clock_timestamp()`, expiry).Scan(&validTime); e != nil {
		return out, e
	}
	if !validTime {
		return out, sv.ErrBinding
	}
	if r.Operation != "reverse" && !s.profile.StateAllowed(r.Operation, state) {
		return out, sv.ErrBinding
	}
	if r.Operation == "redeem" || r.Operation == "reverse" {
		var frozen bool
		// Unknown provider effects remain fenced. Source cancellation may reopen
		// stored value only after the authoritative payment owner reaches refunded.
		e = tx.QueryRow(ctx, `select exists(select 1 from payment.payment_attempt where tenant_id=$1 and order_id=$2 and state not in ('failed','refunded')) or ($3::text<>'cancelled' and exists(select 1 from payment.local_funding_receipt where tenant_id=$1 and order_id=$2))`, tenant, r.OrderID, state).Scan(&frozen)
		if e != nil {
			return out, e
		}
		if frozen {
			return out, sv.ErrBinding
		}
	}
	data := map[string]any{}
	calcOp := "evaluate"
	switch r.Operation {
	case "issue", "accrue":
		if (r.Operation == "issue") != (pp.Calculation.Kind == "gift_card") {
			return out, sv.ErrBinding
		}
		var repeated bool
		e = tx.QueryRow(ctx, `select exists(select 1 from stored_value.operation op where tenant_id=$1 and order_id=$2 and program_id=$3 and operation in ('issue','accrue') and not exists(select 1 from stored_value.operation rev where rev.tenant_id=op.tenant_id and rev.original_operation_id=op.operation_id))`, tenant, r.OrderID, r.ProgramID).Scan(&repeated)
		if e != nil {
			return out, e
		}
		if repeated {
			return out, sv.ErrBinding
		}
	case "redeem":
		if pp.Calculation.Kind == "loyalty" && (a.DiscountMinor != 0 || a.GiftMinor != 0) {
			return out, sv.ErrBinding
		}
		balance, e := s.account(ctx, tx, r, o.SubjectID, r.AccountID)
		if e != nil {
			return out, e
		}
		var futureFromThis bool
		e = tx.QueryRow(ctx, `select exists(select 1 from stored_value.entry en join stored_value.operation op using(tenant_id,operation_id) where en.tenant_id=$1 and en.account_id=$2 and op.order_id=$3 and op.operation='issue')`, tenant, r.AccountID, r.OrderID).Scan(&futureFromThis)
		if e != nil {
			return out, e
		}
		if futureFromThis {
			return out, sv.ErrBinding
		}
		if deciding {
			var reserved bool
			e = tx.QueryRow(ctx, `select exists(select 1 from stored_value.reservation where tenant_id=$1 and operation_id=$2 and account_id=$3 and state='reserved' and expires_at=$4::timestamptz and expires_at>clock_timestamp())`, tenant, r.OperationID, r.AccountID, expiry).Scan(&reserved)
			if e != nil {
				return out, e
			}
			if !reserved {
				return out, sv.ErrBinding
			}
		}
		calcOp = "reward"
		data = map[string]any{"coupon_id": r.AccountID, "balance": balance, "pending_earned": "0", "pending_cost": "0", "discountable": o.Total, "reward": pp.Reward}
	case "reverse":
		if state != "cancelled" {
			return out, sv.ErrBinding
		}
		var originalOrder, originalProgram, operation string
		e = tx.QueryRow(ctx, `select order_id,program_id,operation from stored_value.operation where tenant_id=$1 and operation_id=$2`, tenant, r.OriginalOperationID).Scan(&originalOrder, &originalProgram, &operation)
		if e != nil {
			return out, e
		}
		if originalOrder != r.OrderID || originalProgram != r.ProgramID || operation == "reverse" {
			return out, sv.ErrBinding
		}
		rows, e := tx.Query(ctx, `select account_id,points_delta::text,applied_minor_units from stored_value.entry where tenant_id=$1 and operation_id=$2 order by account_id`, tenant, r.OriginalOperationID)
		if e != nil {
			return out, e
		}
		originals := []sv.Entry{}
		for rows.Next() {
			var entry sv.Entry
			if e = rows.Scan(&entry.AccountID, &entry.PointsDelta, &entry.AppliedMinor); e != nil {
				rows.Close()
				return out, e
			}
			originals = append(originals, entry)
		}
		e = rows.Err()
		rows.Close()
		if e != nil {
			return out, e
		}
		for _, entry := range originals {
			if _, e = s.account(ctx, tx, r, o.SubjectID, entry.AccountID); e != nil {
				return out, e
			}
		}
		calcOp = "reverse"
		data = map[string]any{"entries": originals}
	}
	request, e := sv.NewCalculation(calcOp, pp.Calculation, o, data)
	if e != nil {
		return out, e
	}
	result, e := s.process.Execute(ctx, request)
	if e != nil {
		return out, e
	}
	out.Calculation = request
	out.Result = result
	switch r.Operation {
	case "issue", "accrue":
		var resultData struct {
			Points []string `json:"points"`
			Error  string   `json:"error"`
		}
		if json.Unmarshal(result.Result, &resultData) != nil || resultData.Error != "" {
			return out, sv.ErrBinding
		}
		for i, value := range resultData.Points {
			n, e := points(value)
			if e != nil || n.Sign() < 0 {
				return out, sv.ErrBinding
			}
			if n.Sign() == 0 {
				continue
			}
			id := sv.StableID("gift", tenant, r.OperationID, sv.Number(int64(i)))
			if r.Operation == "accrue" {
				id = sv.StableID("loyalty", tenant, r.OrganizationID, r.ProgramID, o.SubjectID)
			}
			out.Entries = append(out.Entries, sv.Entry{AccountID: id, PointsDelta: decimalPoints(n)})
		}
		if len(out.Entries) == 0 {
			return out, sv.ErrBinding
		}
	case "redeem":
		var v struct {
			Value      string `json:"value"`
			PointsCost string `json:"points_cost"`
		}
		if sv.Decode(result.Result, &v) != nil {
			return out, sv.ErrBinding
		}
		amount, e := sv.Minor(v.Value, 2)
		if e != nil || amount <= 0 || amount > a.ProviderDueMinor {
			return out, sv.ErrBinding
		}
		n, e := points(v.PointsCost)
		if e != nil || n.Sign() <= 0 {
			return out, sv.ErrBinding
		}
		n.Neg(n)
		out.Entries = []sv.Entry{{AccountID: r.AccountID, PointsDelta: decimalPoints(n), AppliedMinor: amount}}
	case "reverse":
		if json.Unmarshal(result.Result, &out.Entries) != nil || len(out.Entries) == 0 {
			return out, sv.ErrBinding
		}
	}
	sort.Slice(out.Entries, func(i, j int) bool { return out.Entries[i].AccountID < out.Entries[j].AccountID })
	return out, nil
}

func (s *StoredValue) commit(ctx context.Context, tx pgx.Tx, approvalID, approvedHash, reviewer string, b sv.BoundOperation) error {
	tenant, _ := s.profile.Scope()
	r := b.Request
	pp, e := s.profile.Program(r.ProgramID)
	if e != nil {
		return e
	}
	receipt := map[string]any{"schema": "elite.stored-value-receipt.v1", "approval_id": approvalID, "approved_sha256": approvedHash, "reviewer": reviewer, "operation": b}
	raw, receiptSHA, e := sv.Canonical(receipt)
	if e != nil {
		return e
	}
	_, requestSHA, e := sv.Canonical(r)
	if e != nil {
		return e
	}
	_, calcSHA, e := sv.Canonical(b.Result)
	if e != nil {
		return e
	}
	for _, entry := range b.Entries {
		if r.Operation == "issue" || r.Operation == "accrue" {
			if _, e = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1,0))`, tenant+":stored-value-account:"+entry.AccountID); e != nil {
				return e
			}
			var owner any
			if pp.Calculation.Nominative {
				owner = b.Order.SubjectID
			}
			_, e = tx.Exec(ctx, `insert into stored_value.account(tenant_id,account_id,organization_id,program_id,kind,customer_subject,currency) values($1,$2,$3,$4,$5,$6,$7) on conflict do nothing`, tenant, entry.AccountID, r.OrganizationID, r.ProgramID, pp.Calculation.Kind, owner, pp.Calculation.Currency)
			if e != nil {
				return e
			}
			if _, e = s.account(ctx, tx, r, b.Order.SubjectID, entry.AccountID); e != nil {
				return e
			}
		}
	}
	_, e = tx.Exec(ctx, `insert into stored_value.operation(tenant_id,operation_id,organization_id,order_id,order_version,program_id,operation,original_operation_id,approval_id,profile_sha256,request_sha256,calculation_sha256,receipt_sha256,receipt) values($1,$2,$3,$4,$5,$6,$7,nullif($8,''),$9,$10,$11,$12,$13,$14)`, tenant, r.OperationID, r.OrganizationID, r.OrderID, r.ExpectedOrderVersion, r.ProgramID, r.Operation, r.OriginalOperationID, approvalID, r.ProfileSHA256, requestSHA, calcSHA, receiptSHA, raw)
	if e != nil {
		return e
	}
	allocated := false
	for _, entry := range b.Entries {
		_, e = tx.Exec(ctx, `insert into stored_value.entry(tenant_id,operation_id,account_id,points_delta,applied_minor_units) values($1,$2,$3,$4::numeric,$5)`, tenant, r.OperationID, entry.AccountID, entry.PointsDelta, entry.AppliedMinor)
		if e != nil {
			return e
		}
		allocated = allocated || entry.AppliedMinor != 0
	}
	if _, e = tx.Exec(ctx, `update stored_value.reservation set state='consumed' where tenant_id=$1 and operation_id=$2 and state='reserved'`, tenant, r.OperationID); e != nil {
		return e
	}
	if allocated {
		if _, e = tx.Exec(ctx, `update sales.customer_order set version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and order_id=$2`, tenant, r.OrderID); e != nil {
			return e
		}
	}
	_, e = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,gen_random_uuid(),'stored-value',$2,1,'stored-value.committed',1,clock_timestamp(),$3)`, tenant, r.OperationID, raw)
	return e
}

func (s *StoredValue) Allocation(ctx context.Context, p identity.Principal, order string) (sv.Allocation, error) {
	var out sv.Allocation
	if !s.allowed(p, "stored_value:read") || !sv.ValidID(order) {
		return out, sv.ErrBinding
	}
	tenant, org := s.profile.Scope()
	e := s.pool.QueryRow(ctx, `select order_id,order_version,currency,gross_minor_units,gift_minor_units,discount_minor_units,provider_due_minor_units from stored_value.order_allocation where tenant_id=$1 and organization_id=$2 and order_id=$3`, tenant, org, order).Scan(&out.OrderID, &out.OrderVersion, &out.Currency, &out.GrossMinor, &out.GiftMinor, &out.DiscountMinor, &out.ProviderDueMinor)
	return out, e
}
