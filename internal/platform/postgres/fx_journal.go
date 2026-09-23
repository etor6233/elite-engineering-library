package postgres

// AUTHORED atomic linkage to the existing draft writer. Conversion, posting and
// reversal algorithms remain with their admitted existing owners.
import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"elite.local/enterprise/internal/accounting"
	"elite.local/enterprise/internal/bcamounts"
	"github.com/jackc/pgx/v5"
)

func readFXJournal(ctx context.Context, q fxReader, tenant, organization, actor, key string) (accounting.FXJournalResult, string, error) {
	var out accounting.FXJournalResult
	var raw []byte
	var hash, requestHash, journal, conversion, currency, period, sourceType, sourceID string
	var amount, credit int64
	var date time.Time
	err := q.QueryRow(ctx, `select r.receipt_raw,r.receipt_sha256_hex,r.request_sha256_hex,r.journal_id,r.conversion_id,j.status,j.version,j.currency,j.period_id,j.posting_date,j.total_debit_minor_units,j.total_credit_minor_units,j.source_type,j.source_id from accounting.fx_journal_receipt r join accounting.journal j on j.tenant_id=r.tenant_id and j.journal_id=r.journal_id and j.organization_id=r.organization_id where r.tenant_id=$1 and r.organization_id=$2 and r.requested_by_subject=$3 and r.request_key=$4`, tenant, organization, actor, key).Scan(&raw, &hash, &requestHash, &journal, &conversion, &out.CurrentStatus, &out.CurrentVersion, &currency, &period, &date, &amount, &credit, &sourceType, &sourceID)
	if errors.Is(err, pgx.ErrNoRows) {
		return out, "", accounting.ErrFXNotFound
	}
	if err != nil {
		return out, "", err
	}
	if fxHash(raw) != hash || json.Unmarshal(raw, &out.Receipt) != nil {
		return accounting.FXJournalResult{}, "", accounting.ErrConflict
	}
	r := out.Receipt
	if r.JournalID != journal || r.ConversionID != conversion || r.OrganizationID != organization || r.RequestedBy != actor || r.RequestKey != key || r.Effect != "FX_JOURNAL_PREPARED" || r.AmountMinor != amount || amount != credit || r.Currency != currency || r.PeriodID != period || r.PostingDate != date.Format("2006-01-02") || sourceType != "FX_CONVERSION" || sourceID != conversion {
		return accounting.FXJournalResult{}, "", accounting.ErrConflict
	}
	return out, requestHash, nil
}
func (r *Accounting) FXJournalResult(ctx context.Context, tenant, organization, actor, key string) (accounting.FXJournalResult, error) {
	out, _, err := readFXJournal(ctx, r.pool, tenant, organization, actor, key)
	return out, err
}

func (r *Accounting) PrepareFXJournal(ctx context.Context, tenant, actor, journalID, journalEvent, bindingEvent string, c accounting.FXJournalCommand, hash string) (accounting.FXJournalResult, bool, error) {
	var empty accounting.FXJournalResult
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return empty, false, err
	}
	defer tx.Rollback(ctx)
	if prior, saved, e := readFXJournal(ctx, tx, tenant, c.OrganizationID, actor, c.IdempotencyKey); e == nil {
		if saved != hash {
			return empty, false, accounting.ErrConflict
		}
		return prior, true, tx.Commit(ctx)
	} else if !errors.Is(e, accounting.ErrFXNotFound) {
		return empty, false, e
	}
	claim, err := tx.Exec(ctx, `insert into platform.idempotency_record(tenant_id,scope,idempotency_key,request_sha256_hex,status,locked_until,expires_at) values($1,'accounting-fx-journal',$2,$3,'processing',clock_timestamp()+interval '30 seconds',clock_timestamp()+interval '24 hours') on conflict do nothing`, tenant, c.IdempotencyKey, hash)
	if err != nil {
		return empty, false, accountingConflict(err)
	}
	if claim.RowsAffected() == 0 {
		var saved string
		if e := tx.QueryRow(ctx, `select request_sha256_hex from platform.idempotency_record where tenant_id=$1 and scope='accounting-fx-journal' and idempotency_key=$2`, tenant, c.IdempotencyKey).Scan(&saved); e != nil || saved != hash {
			return empty, false, accounting.ErrConflict
		}
		prior, saved, e := readFXJournal(ctx, tx, tenant, c.OrganizationID, actor, c.IdempotencyKey)
		if e != nil || saved != hash {
			return empty, false, accounting.ErrConflict
		}
		return prior, true, tx.Commit(ctx)
	}
	var conversionID string
	if err = tx.QueryRow(ctx, `select c.conversion_id from accounting.fx_conversion_receipt c join org.organization o on o.tenant_id=c.tenant_id and o.organization_id=c.organization_id where c.tenant_id=$1 and c.organization_id=$2 and c.requested_by_subject=$3 and c.request_key=$4 and c.conversion_id=$5 and o.status='active' for update of c for share of o`, tenant, c.OrganizationID, actor, c.ConversionRequestKey, c.ConversionID).Scan(&conversionID); err != nil {
		return empty, false, accounting.ErrConflict
	}
	conversion, _, err := readFXReceipt(ctx, tx, tenant, c.OrganizationID, actor, c.ConversionRequestKey)
	if err != nil {
		return empty, false, err
	}
	if conversion.ID != conversionID || conversion.Conversion.ToCurrency != conversion.Snapshot.LocalCurrency || conversion.Conversion.InputMinor <= 0 || bcamounts.CheckBalance(conversion.Conversion.OutputMinor, conversion.Conversion.OutputMinor) != nil {
		return empty, false, accounting.ErrConflict
	}
	var used bool
	if err = tx.QueryRow(ctx, `select exists(select 1 from accounting.fx_journal_receipt where tenant_id=$1 and conversion_id=$2)`, tenant, conversionID).Scan(&used); err != nil {
		return empty, false, err
	}
	if used {
		return empty, false, accounting.ErrConflict
	}
	date, err := time.Parse("2006-01-02", conversion.Conversion.ConversionDate)
	if err != nil {
		return empty, false, accounting.ErrConflict
	}
	amount := conversion.Conversion.OutputMinor
	journal := accounting.Journal{ID: journalID, OrganizationID: c.OrganizationID, PeriodID: c.PeriodID, SourceType: "FX_CONVERSION", SourceID: conversionID, Currency: conversion.Conversion.ToCurrency, PostingDate: date, Status: "draft", Version: 1, TotalDebitMinorUnits: amount, TotalCreditMinorUnits: amount, Lines: []accounting.Line{{LineNo: 1, AccountCode: c.DebitAccount, Description: c.Description, DebitMinorUnits: amount}, {LineNo: 2, AccountCode: c.CreditAccount, Description: c.Description, CreditMinorUnits: amount}}}
	if err = createJournalInTx(ctx, tx, tenant, journalEvent, journal); err != nil {
		return empty, false, err
	}
	var now time.Time
	if err = tx.QueryRow(ctx, `select clock_timestamp()`).Scan(&now); err != nil {
		return empty, false, err
	}
	receipt := accounting.FXJournalReceipt{JournalID: journalID, ConversionID: conversionID, OrganizationID: c.OrganizationID, RequestedBy: actor, RequestKey: c.IdempotencyKey, PeriodID: c.PeriodID, PostingDate: conversion.Conversion.ConversionDate, Currency: journal.Currency, AmountMinor: amount, DebitAccount: c.DebitAccount, CreditAccount: c.CreditAccount, Description: c.Description, RecordedAt: now.UTC(), Effect: "FX_JOURNAL_PREPARED"}
	raw, err := json.Marshal(receipt)
	if err != nil || len(raw) > 16384 {
		return empty, false, accounting.ErrInvalid
	}
	_, err = tx.Exec(ctx, `insert into accounting.fx_journal_receipt(tenant_id,organization_id,conversion_id,journal_id,requested_by_subject,request_key,request_sha256_hex,receipt_raw,receipt_sha256_hex,recorded_at) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`, tenant, c.OrganizationID, conversionID, journalID, actor, c.IdempotencyKey, hash, raw, fxHash(raw), now)
	if err != nil {
		return empty, false, accountingConflict(err)
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'journal',$3,1,'fx-journal.prepared',1,$4,$5)`, tenant, bindingEvent, journalID, now, raw)
	if err != nil {
		return empty, false, accountingConflict(err)
	}
	done, err := tx.Exec(ctx, `update platform.idempotency_record set status='completed',response_code=201,response_body=jsonb_build_object('journal_id',$3::text),resource_type='fx-journal',resource_id=$3,locked_until=null where tenant_id=$1 and scope='accounting-fx-journal' and idempotency_key=$2 and status='processing'`, tenant, c.IdempotencyKey, journalID)
	if err != nil {
		return empty, false, accountingConflict(err)
	}
	if done.RowsAffected() != 1 {
		return empty, false, accounting.ErrConflict
	}
	if err = tx.Commit(ctx); err != nil {
		return empty, false, accountingConflict(err)
	}
	return accounting.FXJournalResult{Receipt: receipt, CurrentStatus: "draft", CurrentVersion: 1}, false, nil
}
