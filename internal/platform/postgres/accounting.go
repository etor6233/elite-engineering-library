package postgres

import (
	"context"
	"elite.local/enterprise/internal/accounting"
	"elite.local/enterprise/internal/bcamounts"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AUTHORED transaction reader boundary; permits existing owner execution
// inside one outer transaction with its immutable command receipt.
type accountingTransactionStarter interface {
	royaltyTransactionStarter
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}
type Accounting struct{ pool accountingTransactionStarter }

func NewAccounting(pool *pgxpool.Pool) *Accounting { return &Accounting{pool: pool} }
func accountingConflict(err error) error {
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && (pgErr.Code == "23503" || pgErr.Code == "23505" || pgErr.Code == "23514" || pgErr.Code == "40001") {
		return accounting.ErrConflict
	}
	return err
}

func (r *Accounting) CreateAccount(ctx context.Context, tenant, eventID string, value accounting.Account) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `insert into accounting.account(tenant_id,account_code,display_name,account_type)values($1,$2,$3,$4)`, tenant, value.Code, value.Name, value.Type)
	if err != nil {
		return accountingConflict(err)
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'account',$3,1,'account.created',1,clock_timestamp(),jsonb_build_object('account_type',$4::text))`, tenant, eventID, value.Code, value.Type)
	if err != nil {
		return accountingConflict(err)
	}
	return accountingConflict(tx.Commit(ctx))
}

func (r *Accounting) OpenPeriod(ctx context.Context, tenant, eventID string, value accounting.Period) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1,0))`, tenant+":accounting-period"); err != nil {
		return err
	}
	var overlap bool
	err = tx.QueryRow(ctx, `select exists(select 1 from accounting.period where tenant_id=$1 and starts_on < ($3::timestamptz at time zone 'UTC')::date and ends_on > ($2::timestamptz at time zone 'UTC')::date)`, tenant, value.StartsOn, value.EndsOn).Scan(&overlap)
	if err != nil {
		return err
	}
	if overlap {
		return accounting.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into accounting.period(tenant_id,period_id,starts_on,ends_on,status,version)values($1,$2,($3::timestamptz at time zone 'UTC')::date,($4::timestamptz at time zone 'UTC')::date,'open',1)`, tenant, value.ID, value.StartsOn, value.EndsOn)
	if err != nil {
		return accountingConflict(err)
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'accounting-period',$3,1,'accounting-period.opened',1,clock_timestamp(),'{}')`, tenant, eventID, value.ID)
	if err != nil {
		return accountingConflict(err)
	}
	return accountingConflict(tx.Commit(ctx))
}

func (r *Accounting) CreateJournal(ctx context.Context, tenant, eventID string, value accounting.Journal) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err = createJournalInTx(ctx, tx, tenant, eventID, value); err != nil {
		return err
	}
	return accountingConflict(tx.Commit(ctx))
}

// AUTHORED transaction composition seam. This is the same existing draft writer.
func createJournalInTx(ctx context.Context, tx pgx.Tx, tenant, eventID string, value accounting.Journal) error {
	var period string
	err := tx.QueryRow(ctx, `select period_id from accounting.period where tenant_id=$1 and period_id=$2 and status='open' and starts_on <= ($3::timestamptz at time zone 'UTC')::date and ends_on > ($3::timestamptz at time zone 'UTC')::date for share`, tenant, value.PeriodID, value.PostingDate).Scan(&period)
	if errors.Is(err, pgx.ErrNoRows) {
		return accounting.ErrConflict
	}
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `insert into accounting.journal(tenant_id,journal_id,organization_id,period_id,source_type,source_id,currency,posting_date,status,total_debit_minor_units,total_credit_minor_units,version)values($1,$2,$3,$4,$5,$6,$7,($8::timestamptz at time zone 'UTC')::date,'draft',$9,$10,1)`, tenant, value.ID, value.OrganizationID, value.PeriodID, value.SourceType, value.SourceID, value.Currency, value.PostingDate, value.TotalDebitMinorUnits, value.TotalCreditMinorUnits)
	if err != nil {
		return accountingConflict(err)
	}
	for _, line := range value.Lines {
		result, lineErr := tx.Exec(ctx, `insert into accounting.journal_line(tenant_id,journal_id,line_no,account_code,description,debit_minor_units,credit_minor_units)select $1,$2,$3,$4,$5,$6,$7 where exists(select 1 from accounting.account where tenant_id=$1 and account_code=$4 and active)`, tenant, value.ID, line.LineNo, line.AccountCode, line.Description, line.DebitMinorUnits, line.CreditMinorUnits)
		if lineErr != nil {
			return accountingConflict(lineErr)
		}
		if result.RowsAffected() != 1 {
			return accounting.ErrConflict
		}
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'journal',$3,1,'journal.created',1,clock_timestamp(),jsonb_build_object('organization_id',$4::text,'source_type',$5::text,'source_id',$6::text,'currency',$7::text))`, tenant, eventID, value.ID, value.OrganizationID, value.SourceType, value.SourceID, value.Currency)
	if err != nil {
		return accountingConflict(err)
	}
	return nil
}

func (r *Accounting) PostJournal(ctx context.Context, tenant, organization, id, actor string, expectedVersion int64, registerID, eventID string) (accounting.Journal, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return accounting.Journal{}, err
	}
	defer tx.Rollback(ctx)
	var value accounting.Journal
	err = tx.QueryRow(ctx, `select j.journal_id,j.organization_id,j.period_id,j.source_type,j.source_id,j.currency,j.posting_date,j.status,j.total_debit_minor_units,j.total_credit_minor_units,j.version from accounting.journal j join accounting.period p on p.tenant_id=j.tenant_id and p.period_id=j.period_id where j.tenant_id=$1 and j.journal_id=$2 and j.organization_id=$3 and j.status='draft' and j.version=$4 and p.status='open' for update of j,p`, tenant, id, organization, expectedVersion).Scan(&value.ID, &value.OrganizationID, &value.PeriodID, &value.SourceType, &value.SourceID, &value.Currency, &value.PostingDate, &value.Status, &value.TotalDebitMinorUnits, &value.TotalCreditMinorUnits, &value.Version)
	if errors.Is(err, pgx.ErrNoRows) {
		return value, accounting.ErrConflict
	}
	if err != nil {
		return value, err
	}
	var debit, credit int64
	err = tx.QueryRow(ctx, `select coalesce(sum(debit_minor_units),0),coalesce(sum(credit_minor_units),0) from accounting.journal_line where tenant_id=$1 and journal_id=$2`, tenant, id).Scan(&debit, &credit)
	if err != nil {
		return value, err
	}
	if bcamounts.CheckBalance(debit, credit) != nil || debit != value.TotalDebitMinorUnits {
		return value, accounting.ErrConflict
	}
	_, err = tx.Exec(ctx, `update accounting.journal set status='posted',version=version+1,posted_by=$4,posted_at=clock_timestamp() where tenant_id=$1 and journal_id=$2 and version=$3`, tenant, id, expectedVersion, actor)
	if err != nil {
		return value, err
	}
	_, err = tx.Exec(ctx, `insert into accounting.register(tenant_id,register_id,journal_id,organization_id,period_id,posted_by,total_debit_minor_units,total_credit_minor_units)values($1,$2,$3,$4,$5,$6,$7,$7)`, tenant, registerID, id, organization, value.PeriodID, actor, debit)
	if err != nil {
		return value, accountingConflict(err)
	}
	result, err := tx.Exec(ctx, `insert into accounting.entry(tenant_id,entry_id,register_id,journal_id,line_no,organization_id,period_id,account_code,currency,posting_date,description,debit_minor_units,credit_minor_units)select l.tenant_id,$3||':'||l.line_no::text,$3,l.journal_id,l.line_no,$4,$5,l.account_code,$6,$7,l.description,l.debit_minor_units,l.credit_minor_units from accounting.journal_line l where l.tenant_id=$1 and l.journal_id=$2 order by l.line_no`, tenant, id, registerID, organization, value.PeriodID, value.Currency, value.PostingDate)
	if err != nil {
		return value, accountingConflict(err)
	}
	if result.RowsAffected() < 2 {
		return value, accounting.ErrConflict
	}
	value.Status = "posted"
	value.Version = expectedVersion + 1
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'journal',$3,$4,'journal.posted',1,clock_timestamp(),jsonb_build_object('organization_id',$5::text,'register_id',$6::text,'currency',$7::text,'total_minor_units',$8::bigint))`, tenant, eventID, id, value.Version, organization, registerID, value.Currency, debit)
	if err != nil {
		return value, accountingConflict(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return value, accountingConflict(err)
	}
	return value, nil
}

func (r *Accounting) ReverseJournal(ctx context.Context, tenant, organization, id, reversalID, registerID, eventID string, expectedVersion int64, actor, reason string) (accounting.Journal, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return accounting.Journal{}, err
	}
	defer tx.Rollback(ctx)
	var original accounting.Journal
	err = tx.QueryRow(ctx, `select j.period_id,j.currency,j.posting_date,j.total_debit_minor_units,j.total_credit_minor_units from accounting.journal j join accounting.period p on p.tenant_id=j.tenant_id and p.period_id=j.period_id where j.tenant_id=$1 and j.journal_id=$2 and j.organization_id=$3 and j.status='posted' and j.version=$4 and j.reversal_of is null and p.status='open' for update of j,p`, tenant, id, organization, expectedVersion).Scan(&original.PeriodID, &original.Currency, &original.PostingDate, &original.TotalDebitMinorUnits, &original.TotalCreditMinorUnits)
	if errors.Is(err, pgx.ErrNoRows) {
		return original, accounting.ErrConflict
	}
	if err != nil {
		return original, err
	}
	_, err = tx.Exec(ctx, `insert into accounting.journal(tenant_id,journal_id,organization_id,period_id,source_type,source_id,currency,posting_date,status,total_debit_minor_units,total_credit_minor_units,version,reversal_of,reversal_reason,posted_by,posted_at)values($1,$2,$3,$4,'REVERSAL',$5,$6,$7,'posted',$8,$8,1,$5,$9,$10,clock_timestamp())`, tenant, reversalID, organization, original.PeriodID, id, original.Currency, original.PostingDate, original.TotalDebitMinorUnits, reason, actor)
	if err != nil {
		return original, accountingConflict(err)
	}
	_, err = tx.Exec(ctx, reverseJournalLinesBCSQL, tenant, id, reversalID)
	if err != nil {
		return original, accountingConflict(err)
	}
	_, err = tx.Exec(ctx, `insert into accounting.register(tenant_id,register_id,journal_id,organization_id,period_id,posted_by,total_debit_minor_units,total_credit_minor_units)values($1,$2,$3,$4,$5,$6,$7,$7)`, tenant, registerID, reversalID, organization, original.PeriodID, actor, original.TotalDebitMinorUnits)
	if err != nil {
		return original, accountingConflict(err)
	}
	result, err := tx.Exec(ctx, `insert into accounting.entry(tenant_id,entry_id,register_id,journal_id,line_no,organization_id,period_id,account_code,currency,posting_date,description,debit_minor_units,credit_minor_units,original_entry_id)select l.tenant_id,$4||':'||l.line_no::text,$4,l.journal_id,l.line_no,$5,$6,l.account_code,$7,$8,l.description,l.debit_minor_units,l.credit_minor_units,e.entry_id from accounting.journal_line l join accounting.entry e on e.tenant_id=l.tenant_id and e.journal_id=$2 and e.line_no=l.line_no where l.tenant_id=$1 and l.journal_id=$3`, tenant, id, reversalID, registerID, organization, original.PeriodID, original.Currency, original.PostingDate)
	if err != nil {
		return original, accountingConflict(err)
	}
	if result.RowsAffected() < 2 {
		return original, accounting.ErrConflict
	}
	result, err = tx.Exec(ctx, `update accounting.journal set status='reversed',version=version+1 where tenant_id=$1 and journal_id=$2 and status='posted' and version=$3`, tenant, id, expectedVersion)
	if err != nil {
		return original, err
	}
	if result.RowsAffected() != 1 {
		return original, accounting.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'journal',$3,1,'journal.reversal-posted',1,clock_timestamp(),jsonb_build_object('reversal_of',$4::text,'organization_id',$5::text,'register_id',$6::text,'reason',$7::text))`, tenant, eventID, reversalID, id, organization, registerID, reason)
	if err != nil {
		return original, accountingConflict(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return original, accountingConflict(err)
	}
	return accounting.Journal{ID: reversalID, OrganizationID: organization, PeriodID: original.PeriodID, SourceType: "REVERSAL", SourceID: id, Currency: original.Currency, PostingDate: original.PostingDate, Status: "posted", TotalDebitMinorUnits: original.TotalDebitMinorUnits, TotalCreditMinorUnits: original.TotalCreditMinorUnits, Version: 1, ReversalOf: id}, nil
}

func (r *Accounting) ClosePeriod(ctx context.Context, tenant, id string, expectedVersion int64, eventID string) (accounting.Period, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return accounting.Period{}, err
	}
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1,0))`, tenant+":accounting-period"); err != nil {
		return accounting.Period{}, err
	}
	var value accounting.Period
	err = tx.QueryRow(ctx, `select period_id,starts_on,ends_on,status,version from accounting.period where tenant_id=$1 and period_id=$2 and status='open' and version=$3 and not exists(select 1 from accounting.journal where tenant_id=$1 and period_id=$2 and status='draft') for update`, tenant, id, expectedVersion).Scan(&value.ID, &value.StartsOn, &value.EndsOn, &value.Status, &value.Version)
	if errors.Is(err, pgx.ErrNoRows) {
		return value, accounting.ErrConflict
	}
	if err != nil {
		return value, err
	}
	_, err = tx.Exec(ctx, `update accounting.period set status='closed',version=version+1,closed_at=clock_timestamp() where tenant_id=$1 and period_id=$2 and version=$3`, tenant, id, expectedVersion)
	if err != nil {
		return value, err
	}
	value.Status = "closed"
	value.Version++
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'accounting-period',$3,$4,'accounting-period.closed',1,clock_timestamp(),'{}')`, tenant, eventID, id, value.Version)
	if err != nil {
		return value, accountingConflict(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return value, accountingConflict(err)
	}
	return value, nil
}

func (r *Accounting) TrialBalance(ctx context.Context, tenant, organization, period string) ([]accounting.Balance, error) {
	rows, err := r.pool.Query(ctx, `select account_code,sum(debit_minor_units)::bigint,sum(credit_minor_units)::bigint,(sum(debit_minor_units)-sum(credit_minor_units))::bigint,(select count(distinct currency) from accounting.entry where tenant_id=$1 and organization_id=$2 and period_id=$3) from accounting.entry where tenant_id=$1 and organization_id=$2 and period_id=$3 group by account_code order by account_code`, tenant, organization, period)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	values := []accounting.Balance{}
	for rows.Next() {
		var v accounting.Balance
		var currencies int
		if err = rows.Scan(&v.AccountCode, &v.DebitMinorUnits, &v.CreditMinorUnits, &v.NetMinorUnits, &currencies); err != nil {
			return nil, err
		}
		if currencies > 1 {
			return nil, accounting.ErrConflict
		} // A legacy balance has no currency field; never combine currencies.
		values = append(values, v)
	}
	return values, rows.Err()
}
