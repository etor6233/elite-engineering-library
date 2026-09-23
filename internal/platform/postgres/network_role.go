// AUTHORED transaction/recovery glue; all business mutations use existing fulfillment writers.
package postgres

import (
	"context"
	"elite.local/enterprise/internal/fulfillment"
	nr "elite.local/enterprise/internal/networkrole"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/randomid"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
	"time"
)

type NetworkRole struct{ pool *pgxpool.Pool }

func NewNetworkRole(pool *pgxpool.Pool) (*NetworkRole, error) {
	if pool == nil {
		return nil, nr.ErrInvalid
	}
	return &NetworkRole{pool}, nil
}

type networkRowReader interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}
type networkCommandIDs struct{ first string }

func (g *networkCommandIDs) New() string {
	if g.first != "" {
		v := g.first
		g.first = ""
		return v
	}
	return randomid.Generator{}.New()
}
func networkError(e error) error {
	if errors.Is(e, pgx.ErrNoRows) {
		return nr.ErrNotFound
	}
	if errors.Is(e, fulfillment.ErrConflict) {
		return nr.ErrConflict
	}
	if e = fulfillmentConstraint(e); errors.Is(e, fulfillment.ErrConflict) {
		return nr.ErrConflict
	}
	return e
}
func readNetworkReceipt(ctx context.Context, q networkRowReader, tenant, id string) (nr.Receipt, error) {
	var v nr.Receipt
	var raw []byte
	var hash string
	e := q.QueryRow(ctx, `select command_id,action,scope_organization_id,actor_subject,request_sha256,entity_payload,entity_sha256,recorded_at from franchise.network_command_receipt where tenant_id=$1 and command_id=$2`, tenant, id).Scan(&v.CommandID, &v.Action, &v.ScopeOrganizationID, &v.Actor, &v.RequestSHA256, &raw, &hash, &v.RecordedAt)
	if e != nil {
		return v, networkError(e)
	}
	if json.Unmarshal(raw, &v.Entity) != nil {
		return v, nr.ErrConflict
	}
	_, actual, e := nr.Canonical(v.Entity)
	if e != nil || actual != hash {
		return v, nr.ErrConflict
	}
	return v, nil
}
func readNetworkEntity(ctx context.Context, q networkRowReader, tenant, kind, id string) (nr.Entity, error) {
	var v nr.Entity
	var version int64
	var e error
	v.Kind = kind
	switch kind {
	case "organization":
		e = q.QueryRow(ctx, `select organization_id,organization_id,coalesce(parent_organization_id,''),organization_code,display_name,organization_type,status,version from org.organization where tenant_id=$1 and organization_id=$2`, tenant, id).Scan(&v.ID, &v.OrganizationID, &v.ParentOrganizationID, &v.Code, &v.DisplayName, &v.Type, &v.State, &version)
	case "agreement":
		e = q.QueryRow(ctx, `select agreement_id,franchise_organization_id,territory_code,terms_version,starts_on::text,coalesce(ends_on::text,''),status,version from franchise.agreement where tenant_id=$1 and agreement_id=$2`, tenant, id).Scan(&v.ID, &v.OrganizationID, &v.TerritoryCode, &v.TermsVersion, &v.StartsOn, &v.EndsOn, &v.State, &version)
	default:
		return v, nr.ErrInvalid
	}
	v.Version = strconv.FormatInt(version, 10)
	return v, networkError(e)
}
func (s *NetworkRole) Result(ctx context.Context, p identity.Principal, scope, id string) (nr.Receipt, error) {
	if !nr.ID(id) || p.TenantID == "" || p.Subject == "" {
		return nr.Receipt{}, nr.ErrNotFound
	}
	v, e := readNetworkReceipt(ctx, s.pool, p.TenantID, id)
	if e != nil {
		return v, e
	}
	if v.Actor != p.Subject || v.ScopeOrganizationID != scope || !nr.Authorized(p, v.Action, scope) {
		return nr.Receipt{}, nr.ErrNotFound
	}
	v.Replay = true
	return v, nil
}
func (s *NetworkRole) Entity(ctx context.Context, p identity.Principal, kind, scope, id string) (nr.Entity, error) {
	action := "transition-" + kind
	if !nr.ID(id) || !nr.Authorized(p, action, scope) {
		return nr.Entity{}, nr.ErrNotFound
	}
	v, e := readNetworkEntity(ctx, s.pool, p.TenantID, kind, id)
	if e != nil {
		return v, e
	}
	if v.OrganizationID != scope {
		return nr.Entity{}, nr.ErrNotFound
	}
	return v, nil
}
func (s *NetworkRole) Execute(ctx context.Context, p identity.Principal, c nr.Command) (nr.Receipt, error) {
	var empty nr.Receipt
	if !c.Valid() {
		return empty, nr.ErrInvalid
	}
	if !nr.Authorized(p, c.Action, c.ScopeOrganizationID) {
		return empty, nr.ErrNotFound
	}
	_, hash, e := nr.Canonical(c)
	if e != nil {
		return empty, e
	}
	tx, e := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if e != nil {
		return empty, e
	}
	defer tx.Rollback(ctx)
	if _, e = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended('network-command:'||$1::text||':'||$2::text,0))`, p.TenantID, c.CommandID); e != nil {
		return empty, e
	}
	old, e := readNetworkReceipt(ctx, tx, p.TenantID, c.CommandID)
	if e == nil {
		if old.Actor != p.Subject || old.Action != c.Action || old.ScopeOrganizationID != c.ScopeOrganizationID || old.RequestSHA256 != hash {
			return empty, nr.ErrConflict
		}
		old.Replay = true
		return old, tx.Commit(ctx)
	}
	if !errors.Is(e, nr.ErrNotFound) {
		return empty, e
	}
	ids := &networkCommandIDs{}
	if c.Action == "create-organization" || c.Action == "create-agreement" {
		ids.first = c.EntityID
	}
	service := fulfillment.NewService(networkTxRepository{tx: tx}, ids)
	kind := ""
	switch c.Action {
	case "create-organization":
		kind = "organization"
		_, e = service.CreateOrganization(ctx, p.TenantID, fulfillment.Organization{ParentOrganizationID: c.ScopeOrganizationID, Code: c.Code, DisplayName: c.DisplayName, Type: c.Type})
	case "transition-organization":
		kind = "organization"
		version, _ := strconv.ParseInt(c.Version, 10, 64)
		e = service.TransitionOrganization(ctx, p.TenantID, c.EntityID, c.Current, c.Target, version)
	case "create-agreement":
		kind = "agreement"
		start, _ := time.Parse("2006-01-02", c.StartsOn)
		var end *time.Time
		if c.EndsOn != "" {
			v, _ := time.Parse("2006-01-02", c.EndsOn)
			end = &v
		}
		_, e = service.CreateAgreement(ctx, p.TenantID, fulfillment.Agreement{OrganizationID: c.ScopeOrganizationID, TerritoryCode: c.TerritoryCode, TermsVersion: c.TermsVersion, StartsOn: start, EndsOn: end})
	case "transition-agreement":
		kind = "agreement"
		version, _ := strconv.ParseInt(c.Version, 10, 64)
		e = service.TransitionAgreement(ctx, p.TenantID, c.ScopeOrganizationID, c.EntityID, c.Current, c.Target, version)
	}
	if e != nil {
		return empty, networkError(e)
	}
	entity, e := readNetworkEntity(ctx, tx, p.TenantID, kind, c.EntityID)
	if e != nil {
		return empty, e
	}
	raw, entityHash, e := nr.Canonical(entity)
	if e != nil {
		return empty, e
	}
	receipt := nr.Receipt{CommandID: c.CommandID, Action: c.Action, ScopeOrganizationID: c.ScopeOrganizationID, Actor: p.Subject, RequestSHA256: hash, Entity: entity}
	e = tx.QueryRow(ctx, `insert into franchise.network_command_receipt(tenant_id,command_id,action,scope_organization_id,actor_subject,request_sha256,entity_payload,entity_sha256)values($1,$2,$3,$4,$5,$6,$7::jsonb,$8) returning recorded_at`, p.TenantID, c.CommandID, c.Action, c.ScopeOrganizationID, p.Subject, hash, string(raw), entityHash).Scan(&receipt.RecordedAt)
	if e != nil {
		return empty, networkError(e)
	}
	return receipt, networkError(tx.Commit(ctx))
}
