package postgres

// AUTHORED PostgreSQL CAS adapter; OAuth algorithms remain in official SDKs.
import (
	"bytes"
	"context"
	"crypto/sha256"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/base64"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
	"time"
)

type PortalSessions struct{ Pool *pgxpool.Pool }

func portalOpaque(s string) bool {
	raw, err := base64.RawURLEncoding.DecodeString(s)
	return err == nil && len(raw) == 32 && base64.RawURLEncoding.EncodeToString(raw) == s
}
func (s *PortalSessions) Execute(ctx context.Context, p identity.PortalProfile, action string, c identity.PortalSessionCommand) (identity.PortalSessionRecord, error) {
	var out identity.PortalSessionRecord
	bad := identity.ErrPortalSessionConflict
	if s == nil || s.Pool == nil || p.SHA256() == "" || !portalOpaque(c.SessionID) || len(c.Ciphertext) > 32768 {
		return out, bad
	}
	key := sha256.Sum256([]byte(c.SessionID))
	cipherHash := sha256.Sum256([]byte(c.Ciphertext))
	tx, err := s.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return out, err
	}
	defer tx.Rollback(ctx)
	var now time.Time
	if err = tx.QueryRow(ctx, `select clock_timestamp()`).Scan(&now); err != nil {
		return out, err
	}
	if action == "create" {
		if c.Version != 0 || c.OperationID != "" || len(c.Ciphertext) < 32 || strings.Count(c.Ciphertext, ".") != 4 || !c.AccessExpiresAt.After(now) || c.AccessExpiresAt.After(c.AbsoluteExpiresAt) || !c.AbsoluteExpiresAt.After(now) || c.AbsoluteExpiresAt.After(now.Add(time.Duration(p.Document().MaximumSessionSeconds)*time.Second)) {
			return out, bad
		}
		_, err = tx.Exec(ctx, `insert into platform.portal_session(profile_sha256,session_id_sha256,state,ciphertext,ciphertext_sha256,access_expires_at,absolute_expires_at) values($1,$2,'active',$3,$4,$5,$6) on conflict do nothing`, p.SHA256(), key[:], c.Ciphertext, cipherHash[:], c.AccessExpiresAt, c.AbsoluteExpiresAt)
		if err != nil {
			return out, err
		}
	}
	var storedHash []byte
	err = tx.QueryRow(ctx, `select version,state,operation_id,ciphertext,ciphertext_sha256,access_expires_at,absolute_expires_at,revocation_pending from platform.portal_session where profile_sha256=$1 and session_id_sha256=$2 for update`, p.SHA256(), key[:]).Scan(&out.Version, &out.State, &out.OperationID, &out.Ciphertext, &storedHash, &out.AccessExpiresAt, &out.AbsoluteExpiresAt, &out.RevocationPending)
	if errors.Is(err, pgx.ErrNoRows) {
		return out, identity.ErrPortalSessionNotFound
	}
	if err != nil {
		return out, err
	}
	currentHash := sha256.Sum256([]byte(out.Ciphertext))
	if !bytes.Equal(currentHash[:], storedHash) {
		return out, bad
	}
	// Re-read time after acquiring the row lock: waiting must not extend a lease/session.
	if err = tx.QueryRow(ctx, `select clock_timestamp()`).Scan(&now); err != nil {
		return out, err
	}
	if !out.AbsoluteExpiresAt.After(now) && out.State != "revoked" {
		_, err = tx.Exec(ctx, `update platform.portal_session set state='reauth_required',revocation_pending=true,updated_at=clock_timestamp() where profile_sha256=$1 and session_id_sha256=$2`, p.SHA256(), key[:])
		if err != nil {
			return out, err
		}
		out.State = "reauth_required"
		out.RevocationPending = true
	}
	switch action {
	case "create":
		if out.Version != 1 || out.State != "active" || out.Ciphertext != c.Ciphertext || !out.AccessExpiresAt.Equal(c.AccessExpiresAt) || !out.AbsoluteExpiresAt.Equal(c.AbsoluteExpiresAt) {
			return out, bad
		}
	case "read":
	case "claim":
		if !portalOpaque(c.OperationID) || c.Version != out.Version || !out.AbsoluteExpiresAt.After(now) {
			return out, bad
		}
		if out.State == "refreshing" && out.OperationID == c.OperationID {
			break
		}
		if out.State != "active" {
			return out, bad
		}
		out.State = "refreshing"
		out.OperationID = c.OperationID
	case "commit":
		if !portalOpaque(c.OperationID) || out.OperationID != c.OperationID || len(c.Ciphertext) < 32 || strings.Count(c.Ciphertext, ".") != 4 || !c.AccessExpiresAt.After(now) || c.AccessExpiresAt.After(out.AbsoluteExpiresAt) || !out.AbsoluteExpiresAt.After(now) {
			return out, bad
		}
		if out.Version == c.Version+1 && out.Ciphertext == c.Ciphertext && out.AccessExpiresAt.Equal(c.AccessExpiresAt) {
			break
		}
		if out.Version != c.Version || (out.State != "refreshing" && out.State != "revoked") {
			return out, bad
		}
		out.Version++
		out.Ciphertext = c.Ciphertext
		out.AccessExpiresAt = c.AccessExpiresAt
		if out.State == "refreshing" {
			out.State = "active"
		} else {
			out.RevocationPending = true
		}
	case "abort":
		if out.OperationID != c.OperationID || out.Version != c.Version {
			return out, bad
		}
		if out.State == "refreshing" {
			out.State = "reauth_required"
			out.RevocationPending = true
		}
	case "revoke":
		out.State = "revoked"
		out.RevocationPending = out.Ciphertext != ""
	case "ack-revocation":
		if out.State != "revoked" || out.Version != c.Version || out.OperationID != c.OperationID {
			return out, bad
		}
		out.RevocationPending = false
		out.Ciphertext = ""
	default:
		return out, bad
	}
	finalHash := sha256.Sum256([]byte(out.Ciphertext))
	_, err = tx.Exec(ctx, `update platform.portal_session set version=$3,state=$4,operation_id=$5,ciphertext=$6,ciphertext_sha256=$7,access_expires_at=$8,revocation_pending=$9,updated_at=clock_timestamp() where profile_sha256=$1 and session_id_sha256=$2`, p.SHA256(), key[:], out.Version, out.State, out.OperationID, out.Ciphertext, finalHash[:], out.AccessExpiresAt, out.RevocationPending)
	if err != nil {
		return out, err
	}
	if err = tx.Commit(ctx); err != nil {
		return identity.PortalSessionRecord{}, err
	}
	return out, nil
}
