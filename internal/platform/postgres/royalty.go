package postgres

import (
	"context"
	"elite.local/enterprise/internal/royalty"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AUTHORED transaction-starting interface permits an outer receipt transaction.
// Existing domain methods and SQL remain unchanged.
type royaltyTransactionStarter interface {
	Begin(context.Context) (pgx.Tx, error)
	BeginTx(context.Context, pgx.TxOptions) (pgx.Tx, error)
}
type Royalty struct{ pool royaltyTransactionStarter }

func NewRoyalty(pool *pgxpool.Pool) *Royalty { return &Royalty{pool: pool} }

func royaltyConflict(err error) error {
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && (pgErr.Code == "23503" || pgErr.Code == "23505" || pgErr.Code == "23514" || pgErr.Code == "40001") {
		return royalty.ErrConflict
	}
	return err
}

func (r *Royalty) CreatePolicy(ctx context.Context, tenant, eventID string, value royalty.Policy) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1,0))`, tenant+":"+value.AgreementID+":"+value.Currency); err != nil {
		return err
	}
	var agreement string
	err = tx.QueryRow(ctx, `select agreement_id from franchise.agreement where tenant_id=$1 and agreement_id=$2 and franchise_organization_id=$3 and status='active' and starts_on <= ($4::timestamptz at time zone 'UTC')::date and (ends_on is null or ends_on > ($4::timestamptz at time zone 'UTC')::date) for update`, tenant, value.AgreementID, value.OrganizationID, value.ValidFrom).Scan(&agreement)
	if errors.Is(err, pgx.ErrNoRows) {
		return royalty.ErrConflict
	}
	if err != nil {
		return err
	}
	var overlap bool
	err = tx.QueryRow(ctx, `select exists(select 1 from royalty.policy where tenant_id=$1 and agreement_id=$2 and currency=$3 and valid_from < coalesce($5::timestamptz,'infinity'::timestamptz) and coalesce(valid_until,'infinity'::timestamptz) > $4)`, tenant, value.AgreementID, value.Currency, value.ValidFrom, value.ValidUntil).Scan(&overlap)
	if err != nil {
		return err
	}
	if overlap {
		return royalty.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into royalty.policy(tenant_id,policy_id,agreement_id,franchise_organization_id,currency,rate_basis_points,valid_from,valid_until) values($1,$2,$3,$4,$5,$6,$7,$8)`, tenant, value.ID, value.AgreementID, value.OrganizationID, value.Currency, value.RateBasisPoints, value.ValidFrom, value.ValidUntil)
	if err != nil {
		return royaltyConflict(err)
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'royalty-policy',$3,1,'royalty-policy.created',1,clock_timestamp(),jsonb_build_object('agreement_id',$4::text,'organization_id',$5::text,'currency',$6::text,'rate_basis_points',$7::integer))`, tenant, eventID, value.ID, value.AgreementID, value.OrganizationID, value.Currency, value.RateBasisPoints)
	if err != nil {
		return royaltyConflict(err)
	}
	return royaltyConflict(tx.Commit(ctx))
}

func (r *Royalty) AccruePayment(ctx context.Context, tenant, sourceEventKey string, value royalty.PaymentEvent, eventID string) (royalty.Accrual, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return royalty.Accrual{}, err
	}
	defer tx.Rollback(ctx)
	var amount int64
	var currency string
	err = tx.QueryRow(ctx, `select p.amount_minor_units,p.currency from payment.payment_attempt p join sales.customer_order o on o.tenant_id=p.tenant_id and o.order_id=p.order_id where p.tenant_id=$1 and p.payment_attempt_id=$2 and p.state=$3 and p.version=$4 and o.organization_id=$5 for update of p`, tenant, value.PaymentAttemptID, value.State, value.PaymentExpectedVersion, value.OrganizationID).Scan(&amount, &currency)
	if errors.Is(err, pgx.ErrNoRows) {
		return royalty.Accrual{}, royalty.ErrConflict
	}
	if err != nil {
		return royalty.Accrual{}, err
	}
	rows, err := tx.Query(ctx, `select p.policy_id,p.agreement_id,p.rate_basis_points from royalty.policy p join franchise.agreement a on a.tenant_id=p.tenant_id and a.agreement_id=p.agreement_id where p.tenant_id=$1 and p.franchise_organization_id=$2 and p.currency=$3 and p.valid_from <= $4 and (p.valid_until is null or p.valid_until > $4) and a.status='active' and a.starts_on <= ($4::timestamptz at time zone 'UTC')::date and (a.ends_on is null or a.ends_on > ($4::timestamptz at time zone 'UTC')::date) order by p.valid_from desc limit 2`, tenant, value.OrganizationID, currency, value.OccurredAt)
	if err != nil {
		return royalty.Accrual{}, err
	}
	defer rows.Close()
	type resolved struct {
		policy, agreement string
		rate              int
	}
	found := []resolved{}
	for rows.Next() {
		var v resolved
		if err = rows.Scan(&v.policy, &v.agreement, &v.rate); err != nil {
			return royalty.Accrual{}, err
		}
		found = append(found, v)
	}
	if err = rows.Err(); err != nil {
		return royalty.Accrual{}, err
	}
	if len(found) != 1 {
		return royalty.Accrual{}, royalty.ErrConflict
	}
	var royaltyAmount int64
	err = tx.QueryRow(ctx, `select round(($1::numeric*$2::numeric)/10000)::bigint`, amount, found[0].rate).Scan(&royaltyAmount)
	if err != nil {
		return royalty.Accrual{}, err
	}
	if royaltyAmount == 0 {
		return royalty.Accrual{}, royalty.ErrConflict
	}
	if value.State == "refunded" {
		royaltyAmount = -royaltyAmount
	}
	result := royalty.Accrual{ID: value.ID, PolicyID: found[0].policy, AgreementID: found[0].agreement, OrganizationID: value.OrganizationID, PaymentAttemptID: value.PaymentAttemptID, SourceEventKey: sourceEventKey, SourceState: value.State, Currency: currency, BasisMinorUnits: amount, RoyaltyMinorUnits: royaltyAmount, OccurredAt: value.OccurredAt}
	_, err = tx.Exec(ctx, `insert into royalty.accrual(tenant_id,accrual_id,policy_id,agreement_id,franchise_organization_id,payment_attempt_id,source_event_key,source_state,currency,basis_minor_units,royalty_minor_units,occurred_at) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`, tenant, result.ID, result.PolicyID, result.AgreementID, result.OrganizationID, result.PaymentAttemptID, result.SourceEventKey, result.SourceState, result.Currency, result.BasisMinorUnits, result.RoyaltyMinorUnits, result.OccurredAt)
	if err != nil {
		return royalty.Accrual{}, royaltyConflict(err)
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'royalty-accrual',$3,1,'royalty-accrual.posted',1,clock_timestamp(),jsonb_build_object('payment_attempt_id',$4::text,'organization_id',$5::text,'currency',$6::text,'royalty_minor_units',$7::bigint,'source_state',$8::text))`, tenant, eventID, result.ID, result.PaymentAttemptID, result.OrganizationID, result.Currency, result.RoyaltyMinorUnits, result.SourceState)
	if err != nil {
		return royalty.Accrual{}, royaltyConflict(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return royalty.Accrual{}, royaltyConflict(err)
	}
	return result, nil
}

func (r *Royalty) OpenSettlement(ctx context.Context, tenant, eventID string, value royalty.Settlement) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `insert into royalty.settlement_run(tenant_id,settlement_id,franchise_organization_id,currency,period_start,period_end,status,expected_minor_units,version) select $1,$2,$3,$4,$5,$6,'draft',0,1 where exists(select 1 from franchise.agreement where tenant_id=$1 and franchise_organization_id=$3 and status='active')`, tenant, value.ID, value.OrganizationID, value.Currency, value.PeriodStart, value.PeriodEnd)
	if err != nil {
		return royaltyConflict(err)
	}
	var exists bool
	if err = tx.QueryRow(ctx, `select exists(select 1 from royalty.settlement_run where tenant_id=$1 and settlement_id=$2)`, tenant, value.ID).Scan(&exists); err != nil {
		return err
	}
	if !exists {
		return royalty.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'royalty-settlement',$3,1,'royalty-settlement.opened',1,clock_timestamp(),jsonb_build_object('organization_id',$4::text,'currency',$5::text))`, tenant, eventID, value.ID, value.OrganizationID, value.Currency)
	if err != nil {
		return royaltyConflict(err)
	}
	return royaltyConflict(tx.Commit(ctx))
}

func (r *Royalty) CloseSettlement(ctx context.Context, tenant, organization, id, eventID string, expectedVersion int64) (royalty.Settlement, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return royalty.Settlement{}, err
	}
	defer tx.Rollback(ctx)
	var value royalty.Settlement
	err = tx.QueryRow(ctx, `select settlement_id,franchise_organization_id,currency,period_start,period_end,status,expected_minor_units,version,coalesce(reversal_of,'') from royalty.settlement_run where tenant_id=$1 and settlement_id=$2 and franchise_organization_id=$3 and status='draft' and version=$4 for update`, tenant, id, organization, expectedVersion).Scan(&value.ID, &value.OrganizationID, &value.Currency, &value.PeriodStart, &value.PeriodEnd, &value.Status, &value.ExpectedMinorUnits, &value.Version, &value.ReversalOf)
	if errors.Is(err, pgx.ErrNoRows) {
		return royalty.Settlement{}, royalty.ErrConflict
	}
	if err != nil {
		return royalty.Settlement{}, err
	}
	if _, err = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1,0))`, tenant+":"+organization+":"+value.Currency); err != nil {
		return royalty.Settlement{}, err
	}
	result, err := tx.Exec(ctx, `insert into royalty.settlement_line(tenant_id,settlement_id,line_id,accrual_id,amount_minor_units) select a.tenant_id,$2,a.accrual_id,a.accrual_id,a.royalty_minor_units from royalty.accrual a where a.tenant_id=$1 and a.franchise_organization_id=$3 and a.currency=$4 and a.occurred_at >= $5 and a.occurred_at < $6 and not exists(select 1 from royalty.settlement_line l where l.tenant_id=a.tenant_id and l.accrual_id=a.accrual_id and l.original_line_id is null) order by a.occurred_at,a.accrual_id`, tenant, id, organization, value.Currency, value.PeriodStart, value.PeriodEnd)
	if err != nil {
		return royalty.Settlement{}, royaltyConflict(err)
	}
	if result.RowsAffected() == 0 {
		return royalty.Settlement{}, royalty.ErrConflict
	}
	err = tx.QueryRow(ctx, `update royalty.settlement_run s set status='closed',expected_minor_units=x.total,version=version+1,closed_at=clock_timestamp() from (select sum(amount_minor_units)::bigint total from royalty.settlement_line where tenant_id=$1 and settlement_id=$2) x where s.tenant_id=$1 and s.settlement_id=$2 and s.status='draft' and s.version=$3 returning s.expected_minor_units,s.version,s.status`, tenant, id, expectedVersion).Scan(&value.ExpectedMinorUnits, &value.Version, &value.Status)
	if errors.Is(err, pgx.ErrNoRows) {
		return royalty.Settlement{}, royalty.ErrConflict
	}
	if err != nil {
		return royalty.Settlement{}, err
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'royalty-settlement',$3,$4,'royalty-settlement.closed',1,clock_timestamp(),jsonb_build_object('organization_id',$5::text,'currency',$6::text,'expected_minor_units',$7::bigint))`, tenant, eventID, id, value.Version, organization, value.Currency, value.ExpectedMinorUnits)
	if err != nil {
		return royalty.Settlement{}, royaltyConflict(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return royalty.Settlement{}, royaltyConflict(err)
	}
	return value, nil
}

func (r *Royalty) ReverseSettlement(ctx context.Context, tenant, organization, id, reversalID, eventID string, expectedVersion int64, reason string) (royalty.Settlement, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return royalty.Settlement{}, err
	}
	defer tx.Rollback(ctx)
	var value royalty.Settlement
	err = tx.QueryRow(ctx, `select settlement_id,franchise_organization_id,currency,period_start,period_end,status,expected_minor_units,version from royalty.settlement_run where tenant_id=$1 and settlement_id=$2 and franchise_organization_id=$3 and status='closed' and version=$4 and reversal_of is null for update`, tenant, id, organization, expectedVersion).Scan(&value.ID, &value.OrganizationID, &value.Currency, &value.PeriodStart, &value.PeriodEnd, &value.Status, &value.ExpectedMinorUnits, &value.Version)
	if errors.Is(err, pgx.ErrNoRows) {
		return royalty.Settlement{}, royalty.ErrConflict
	}
	if err != nil {
		return royalty.Settlement{}, err
	}
	_, err = tx.Exec(ctx, `insert into royalty.settlement_run(tenant_id,settlement_id,franchise_organization_id,currency,period_start,period_end,status,expected_minor_units,version,reversal_of,reversal_reason,closed_at) values($1,$2,$3,$4,$5,$6,'closed',$7,1,$8,$9,clock_timestamp())`, tenant, reversalID, organization, value.Currency, value.PeriodStart, value.PeriodEnd, -value.ExpectedMinorUnits, id, reason)
	if err != nil {
		return royalty.Settlement{}, royaltyConflict(err)
	}
	result, err := tx.Exec(ctx, `insert into royalty.settlement_line(tenant_id,settlement_id,line_id,accrual_id,amount_minor_units,original_line_id) select tenant_id,$3,$3||':'||line_id,accrual_id,-amount_minor_units,$2||':'||line_id from royalty.settlement_line where tenant_id=$1 and settlement_id=$2`, tenant, id, reversalID)
	if err != nil {
		return royalty.Settlement{}, royaltyConflict(err)
	}
	if result.RowsAffected() == 0 {
		return royalty.Settlement{}, royalty.ErrConflict
	}
	result, err = tx.Exec(ctx, `update royalty.settlement_run set status='reversed',version=version+1 where tenant_id=$1 and settlement_id=$2 and status='closed' and version=$3`, tenant, id, expectedVersion)
	if err != nil {
		return royalty.Settlement{}, err
	}
	if result.RowsAffected() != 1 {
		return royalty.Settlement{}, royalty.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'royalty-settlement',$3,1,'royalty-settlement.reversal-posted',1,clock_timestamp(),jsonb_build_object('reversal_of',$4::text,'organization_id',$5::text,'expected_minor_units',$6::bigint,'reason',$7::text))`, tenant, eventID, reversalID, id, organization, -value.ExpectedMinorUnits, reason)
	if err != nil {
		return royalty.Settlement{}, royaltyConflict(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return royalty.Settlement{}, royaltyConflict(err)
	}
	return royalty.Settlement{ID: reversalID, OrganizationID: organization, Currency: value.Currency, PeriodStart: value.PeriodStart, PeriodEnd: value.PeriodEnd, Status: "closed", ExpectedMinorUnits: -value.ExpectedMinorUnits, Version: 1, ReversalOf: id}, nil
}

func (r *Royalty) Reconcile(ctx context.Context, tenant, organization string, value royalty.Reconciliation, actor, eventID string) (royalty.Reconciliation, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return value, err
	}
	defer tx.Rollback(ctx)
	err = tx.QueryRow(ctx, `select expected_minor_units from royalty.settlement_run where tenant_id=$1 and settlement_id=$2 and franchise_organization_id=$3 and status='closed' for update`, tenant, value.SettlementID, organization).Scan(&value.ExpectedMinorUnits)
	if errors.Is(err, pgx.ErrNoRows) {
		return value, royalty.ErrConflict
	}
	if err != nil {
		return value, err
	}
	value.DifferenceMinorUnits = value.ActualMinorUnits - value.ExpectedMinorUnits
	value.Status = "mismatch"
	if value.DifferenceMinorUnits == 0 {
		value.Status = "matched"
	}
	_, err = tx.Exec(ctx, `insert into royalty.reconciliation(tenant_id,reconciliation_id,settlement_id,external_reference,expected_minor_units,actual_minor_units,difference_minor_units,status,recorded_by) values($1,$2,$3,$4,$5,$6,$7,$8,$9)`, tenant, value.ID, value.SettlementID, value.ExternalReference, value.ExpectedMinorUnits, value.ActualMinorUnits, value.DifferenceMinorUnits, value.Status, actor)
	if err != nil {
		return value, royaltyConflict(err)
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'royalty-reconciliation',$3,1,'royalty-reconciliation.recorded',1,clock_timestamp(),jsonb_build_object('settlement_id',$4::text,'status',$5::text,'difference_minor_units',$6::bigint))`, tenant, eventID, value.ID, value.SettlementID, value.Status, value.DifferenceMinorUnits)
	if err != nil {
		return value, royaltyConflict(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return value, royaltyConflict(err)
	}
	return value, nil
}
