package postgres

import (
	"bytes"
	"context"
	"crypto/sha256"
	"elite.local/enterprise/internal/platform/identity"
	"github.com/jackc/pgx/v5"
)

// Sweep leases a bounded batch without needing a browser cookie. Retention
// counts unconfirmed provider revocation durably; it never labels it success.
func (s *PortalSessions) Sweep(ctx context.Context, p identity.PortalProfile, operation string) (identity.PortalSessionSweep, error) {
	result := identity.PortalSessionSweep{Items: []identity.PortalSessionSweepItem{}}
	if p.SHA256() == "" || !portalOpaque(operation) || s == nil || s.Pool == nil {
		return result, identity.ErrPortalSessionConflict
	}
	tx, err := s.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return result, err
	}
	defer tx.Rollback(ctx)
	err = tx.QueryRow(ctx, `with picked as (select profile_sha256,session_id_sha256 from platform.portal_session where profile_sha256=$1 and absolute_expires_at<clock_timestamp()-make_interval(secs=>$2) order by absolute_expires_at for update skip locked limit 100), removed as (delete from platform.portal_session s using picked p where s.profile_sha256=p.profile_sha256 and s.session_id_sha256=p.session_id_sha256 returning s.revocation_pending,s.state) select count(*),count(*) filter(where revocation_pending or state<>'revoked') from removed`, p.SHA256(), p.Document().RetentionSeconds).Scan(&result.Purged, &result.UnconfirmedPurged)
	if err != nil {
		return result, err
	}
	if result.Purged > 0 {
		_, err = tx.Exec(ctx, `insert into platform.portal_session_retention_summary(profile_sha256,purged_count,unconfirmed_provider_revocation_count) values($1,$2,$3) on conflict(profile_sha256) do update set purged_count=platform.portal_session_retention_summary.purged_count+excluded.purged_count,unconfirmed_provider_revocation_count=platform.portal_session_retention_summary.unconfirmed_provider_revocation_count+excluded.unconfirmed_provider_revocation_count,updated_at=clock_timestamp()`, p.SHA256(), result.Purged, result.UnconfirmedPurged)
		if err != nil {
			return result, err
		}
	}
	rows, err := tx.Query(ctx, `with picked as (select profile_sha256,session_id_sha256 from platform.portal_session where profile_sha256=$1 and ciphertext<>'' and revocation_lease_until<clock_timestamp() and (state in ('reauth_required','revoked') or absolute_expires_at<=clock_timestamp() or (state='refreshing' and updated_at<clock_timestamp()-interval '120 seconds')) order by updated_at for update skip locked limit 2) update platform.portal_session s set state='revoked',revocation_pending=true,revocation_lease_until=clock_timestamp()+interval '120 seconds',revocation_operation_id=$2,updated_at=clock_timestamp() from picked p where s.profile_sha256=p.profile_sha256 and s.session_id_sha256=p.session_id_sha256 returning encode(s.session_id_sha256,'hex'),s.version,s.state,s.operation_id,s.ciphertext,s.ciphertext_sha256,s.access_expires_at,s.absolute_expires_at,s.revocation_pending`, p.SHA256(), operation)
	if err != nil {
		return result, err
	}
	for rows.Next() {
		var item identity.PortalSessionSweepItem
		var hash []byte
		if err = rows.Scan(&item.SessionIDSHA256, &item.Version, &item.State, &item.OperationID, &item.Ciphertext, &hash, &item.AccessExpiresAt, &item.AbsoluteExpiresAt, &item.RevocationPending); err != nil {
			rows.Close()
			return result, err
		}
		current := sha256.Sum256([]byte(item.Ciphertext))
		if !bytes.Equal(current[:], hash) {
			rows.Close()
			return result, identity.ErrPortalSessionConflict
		}
		result.Items = append(result.Items, item)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return result, err
	}
	if err = tx.Commit(ctx); err != nil {
		return result, err
	}
	return result, nil
}
