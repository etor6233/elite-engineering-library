package postgres

// AUTHORED immutable receipt, shared idempotency and transactional outbox glue.
import (
	"context"
	"crypto/sha256"
	"elite.local/enterprise/internal/accounting"
	"elite.local/enterprise/internal/bcfx"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"time"
)

type fxReader interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}

func fxHash(raw []byte) string { v := sha256.Sum256(raw); return hex.EncodeToString(v[:]) }
func readFXReceipt(ctx context.Context, q fxReader, tenant, organization, actor, key string) (accounting.FXReceipt, string, error) {
	var v accounting.FXReceipt
	var raw []byte
	var digest, requestHash, id string
	err := q.QueryRow(ctx, `select conversion_id,receipt_raw,receipt_sha256_hex,request_sha256_hex from accounting.fx_conversion_receipt where tenant_id=$1 and organization_id=$2 and requested_by_subject=$3 and request_key=$4`, tenant, organization, actor, key).Scan(&id, &raw, &digest, &requestHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return v, "", accounting.ErrFXNotFound
	}
	if err != nil {
		return v, "", err
	}
	if fxHash(raw) != digest || json.Unmarshal(raw, &v) != nil || v.ID != id || v.RequestedBy != actor || v.RequestKey != key || v.Snapshot.TenantID != tenant || v.Snapshot.OrganizationID != organization || v.Effect != "CONVERSION_RECEIPT_ONLY" {
		return accounting.FXReceipt{}, "", accounting.ErrConflict
	}
	return v, requestHash, nil
}
func (r *Accounting) FXConversionResult(ctx context.Context, tenant, organization, actor, key string) (accounting.FXReceipt, error) {
	v, _, err := readFXReceipt(ctx, r.pool, tenant, organization, actor, key)
	return v, err
}
func (r *Accounting) RecordFXConversion(ctx context.Context, tenant, actor, id, event string, c accounting.FXCommand, s *bcfx.Snapshot, hash string) (accounting.FXReceipt, bool, error) {
	var empty accounting.FXReceipt
	if s == nil || c.AmountMinor == nil || !s.Allows(tenant, c.OrganizationID) || len(hash) != 64 {
		return empty, false, accounting.ErrInvalid
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return empty, false, err
	}
	defer tx.Rollback(ctx)
	// A durable receipt still prevents duplicate conversion if shared transient
	// idempotency rows have been retained for a shorter period by the host.
	if prior, saved, e := readFXReceipt(ctx, tx, tenant, c.OrganizationID, actor, c.IdempotencyKey); e == nil {
		if saved != hash {
			return empty, false, accounting.ErrConflict
		}
		return prior, true, tx.Commit(ctx)
	} else if !errors.Is(e, accounting.ErrFXNotFound) {
		return empty, false, e
	}
	claim, err := tx.Exec(ctx, `insert into platform.idempotency_record(tenant_id,scope,idempotency_key,request_sha256_hex,status,locked_until,expires_at) values($1,'accounting-fx-conversion',$2,$3,'processing',clock_timestamp()+interval '30 seconds',clock_timestamp()+interval '24 hours') on conflict do nothing`, tenant, c.IdempotencyKey, hash)
	if err != nil {
		return empty, false, accountingConflict(err)
	}
	if claim.RowsAffected() == 0 {
		var saved string
		if e := tx.QueryRow(ctx, `select request_sha256_hex from platform.idempotency_record where tenant_id=$1 and scope='accounting-fx-conversion' and idempotency_key=$2`, tenant, c.IdempotencyKey).Scan(&saved); e != nil || saved != hash {
			return empty, false, accounting.ErrConflict
		}
		v, receiptHash, e := readFXReceipt(ctx, tx, tenant, c.OrganizationID, actor, c.IdempotencyKey)
		if e != nil || receiptHash != hash {
			return empty, false, accounting.ErrConflict
		}
		return v, true, tx.Commit(ctx)
	}
	var organization string
	if err = tx.QueryRow(ctx, `select organization_id from org.organization where tenant_id=$1 and organization_id=$2 and status='active' for share`, tenant, c.OrganizationID).Scan(&organization); err != nil {
		return empty, false, accounting.ErrConflict
	}
	identity := s.Identity()
	profileRaw, sourceRaw := s.Bytes()
	_, err = tx.Exec(ctx, `insert into accounting.fx_rate_snapshot(tenant_id,profile_id,profile_revision,organization_id,profile_sha256_hex,source_sha256_hex,profile_raw,source_raw,valid_from,valid_until) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) on conflict do nothing`, tenant, identity.ProfileID, identity.Revision, c.OrganizationID, identity.ProfileSHA256, identity.SourceSHA256, profileRaw, sourceRaw, identity.ValidFrom, identity.ValidUntil)
	if err != nil {
		return empty, false, accountingConflict(err)
	}
	var savedProfile, savedSource, savedOrg string
	if err = tx.QueryRow(ctx, `select profile_sha256_hex,source_sha256_hex,organization_id from accounting.fx_rate_snapshot where tenant_id=$1 and profile_id=$2 and profile_revision=$3`, tenant, identity.ProfileID, identity.Revision).Scan(&savedProfile, &savedSource, &savedOrg); err != nil || savedProfile != identity.ProfileSHA256 || savedSource != identity.SourceSHA256 || savedOrg != c.OrganizationID {
		return empty, false, accounting.ErrConflict
	}
	var now time.Time
	if err = tx.QueryRow(ctx, `select clock_timestamp()`).Scan(&now); err != nil {
		return empty, false, err
	}
	value, err := s.Convert(c.FromCurrency, c.ToCurrency, c.ConversionDate, *c.AmountMinor, now)
	if err != nil {
		return empty, false, accounting.ErrConflict
	}
	receipt := accounting.FXReceipt{ID: id, RequestedBy: actor, RequestKey: c.IdempotencyKey, RecordedAt: now.UTC(), Snapshot: identity, Conversion: value, Effect: "CONVERSION_RECEIPT_ONLY"}
	raw, err := json.Marshal(receipt)
	if err != nil || len(raw) > 16384 {
		return empty, false, accounting.ErrInvalid
	}
	_, err = tx.Exec(ctx, `insert into accounting.fx_conversion_receipt(tenant_id,conversion_id,organization_id,requested_by_subject,request_key,request_sha256_hex,profile_id,profile_revision,profile_sha256_hex,receipt_raw,receipt_sha256_hex,recorded_at) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`, tenant, id, c.OrganizationID, actor, c.IdempotencyKey, hash, identity.ProfileID, identity.Revision, identity.ProfileSHA256, raw, fxHash(raw), now)
	if err != nil {
		return empty, false, accountingConflict(err)
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'fx-conversion',$3,1,'fx-conversion.recorded',1,$4,$5)`, tenant, event, id, now, raw)
	if err != nil {
		return empty, false, accountingConflict(err)
	}
	// Expiry is rechecked after blocking inserts/outbox, before any commit.
	if err = tx.QueryRow(ctx, `select clock_timestamp()`).Scan(&now); err != nil || !s.Current(now) {
		return empty, false, accounting.ErrConflict
	}
	done, err := tx.Exec(ctx, `update platform.idempotency_record set status='completed',response_code=201,response_body=jsonb_build_object('conversion_id',$3::text),resource_type='fx-conversion',resource_id=$3,locked_until=null where tenant_id=$1 and scope='accounting-fx-conversion' and idempotency_key=$2 and status='processing'`, tenant, c.IdempotencyKey, id)
	if err != nil {
		return empty, false, accountingConflict(err)
	}
	if done.RowsAffected() != 1 {
		return empty, false, accounting.ErrConflict
	}
	if err = tx.Commit(ctx); err != nil {
		return empty, false, accountingConflict(err)
	}
	return receipt, false, nil
}
