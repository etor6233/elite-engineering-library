// AUTHORED transaction/identity binding for existing catalog/commerce/approval.
package postgres

import (
	"context"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/randomid"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CatalogRelease struct {
	pool    *pgxpool.Pool
	profile cr.Profile
	price   *Commerce
}

func NewCatalogRelease(pool *pgxpool.Pool, profile cr.Profile, price *Commerce) (*CatalogRelease, error) {
	if pool == nil || !profile.Valid() || price == nil || price.pool != pool {
		return nil, cr.ErrInvalid
	}
	return &CatalogRelease{pool, profile, price}, nil
}
func (s *CatalogRelease) authorized(p identity.Principal, permission string) bool {
	return s != nil && approvalPrincipal(p, s.profile.TenantID, s.profile.OrganizationID, permission)
}
func (s *CatalogRelease) begin(ctx context.Context, p identity.Principal, permission string) (pgx.Tx, error) {
	if !s.authorized(p, permission) {
		return nil, cr.ErrNotFound
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	ok := false
	defer func() {
		if !ok {
			_ = tx.Rollback(ctx)
		}
	}()
	if _, err = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1,0))`, "catalog-release:"+s.profile.TenantID); err != nil {
		return nil, err
	}
	raw, hash, err := cr.Canonical(s.profile)
	if err != nil {
		return nil, err
	}
	var exists bool
	err = tx.QueryRow(ctx, `select exists(select 1 from catalog.release_stream where tenant_id=$1)`, s.profile.TenantID).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if !exists {
		var active bool
		err = tx.QueryRow(ctx, `select exists(select 1 from pricing.price_book where tenant_id=$1 and market=$2 and currency=$3 and status='active')`, s.profile.TenantID, s.profile.Market, s.profile.Currency).Scan(&active)
		if err != nil {
			return nil, err
		}
		if active {
			return nil, cr.ErrConflict
		}
		_, err = tx.Exec(ctx, `insert into catalog.release_stream(tenant_id,organization_id,market,currency,profile,profile_canonical,profile_sha256)
   select $1,$2,$3,$4,$5,$6,$7 where exists(select 1 from org.organization where tenant_id=$1 and organization_id=$2 and status='active')`,
			s.profile.TenantID, s.profile.OrganizationID, s.profile.Market, s.profile.Currency, raw, string(raw), hash)
		if err != nil {
			return nil, err
		}
	}
	var same bool
	err = tx.QueryRow(ctx, `select s.profile_sha256=$2 and o.status='active' and t.status='active' from catalog.release_stream s
  join org.organization o using(tenant_id,organization_id) join platform.tenant t using(tenant_id) where s.tenant_id=$1`, s.profile.TenantID, hash).Scan(&same)
	if errors.Is(err, pgx.ErrNoRows) || err == nil && !same {
		return nil, cr.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	ok = true
	return tx, nil
}
func (s *CatalogRelease) replay(ctx context.Context, tx pgx.Tx, p identity.Principal, key, kind, hash string) (cr.Receipt, bool, error) {
	var raw []byte
	var actor, oldKind, oldHash string
	err := tx.QueryRow(ctx, `select actor,kind,request_sha256,receipt from catalog.release_command where tenant_id=$1 and command_id=$2`, p.TenantID, key).Scan(&actor, &oldKind, &oldHash, &raw)
	if errors.Is(err, pgx.ErrNoRows) {
		return cr.Receipt{}, false, nil
	}
	if err != nil {
		return cr.Receipt{}, false, err
	}
	if actor != p.Subject || oldKind != kind || oldHash != hash {
		return cr.Receipt{}, false, cr.ErrConflict
	}
	var r cr.Receipt
	err = json.Unmarshal(raw, &r)
	r.Replay = true
	return r, true, err
}
func (s *CatalogRelease) save(ctx context.Context, tx pgx.Tx, p identity.Principal, r cr.Receipt) error {
	raw, _, err := cr.Canonical(r)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `insert into catalog.release_command(tenant_id,command_id,actor,kind,resource_id,request_sha256,receipt) values($1,$2,$3,$4,$5,$6,$7)`,
		p.TenantID, r.CommandID, p.Subject, r.Kind, r.ResourceID, r.RequestSHA256, raw)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)
 values($1,$2,'catalog-release-command',$3,$4,$5,1,clock_timestamp(),$6)`, p.TenantID, (randomid.Generator{}).New(), r.CommandID, 1, "catalog-release."+r.Kind, raw)
	return err
}
func (s *CatalogRelease) CommandReceipt(ctx context.Context, p identity.Principal, key string) (cr.Receipt, error) {
	if !s.authorized(p, "catalog:read") || !cr.ValidID(key) {
		return cr.Receipt{}, cr.ErrNotFound
	}
	var raw []byte
	err := s.pool.QueryRow(ctx, `select receipt from catalog.release_command where tenant_id=$1 and command_id=$2 and actor=$3`, p.TenantID, key, p.Subject).Scan(&raw)
	if errors.Is(err, pgx.ErrNoRows) {
		return cr.Receipt{}, cr.ErrNotFound
	}
	if err != nil {
		return cr.Receipt{}, err
	}
	var r cr.Receipt
	err = json.Unmarshal(raw, &r)
	return r, err
}
