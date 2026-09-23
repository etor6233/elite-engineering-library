// AUTHORED atomic consent, historical scope and outbox composition. Coverage
// predicates come from the separately declared BC adaptation, not this glue.
package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"elite.local/enterprise/internal/platform/identity"
	wc "elite.local/enterprise/internal/warrantyclaim"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Warranty struct {
	pool    *pgxpool.Pool
	profile wc.Profile
}

func (s *Warranty) Scope() (string, string) {
	if s == nil {
		return "", ""
	}
	return s.profile.Scope()
}
func (s *Warranty) ProfileDocument(ctx context.Context, p identity.Principal) (json.RawMessage, error) {
	if !s.allowed(p, "warranty:read") {
		return nil, wc.ErrNotFound
	}
	return s.profile.Bytes(), nil
}

func NewWarranty(pool *pgxpool.Pool, profile wc.Profile) (*Warranty, error) {
	if pool == nil || !profile.Valid() {
		return nil, wc.ErrInvalid
	}
	return &Warranty{pool: pool, profile: profile}, nil
}
func (s *Warranty) ProfileSHA256() string {
	if s == nil {
		return ""
	}
	return s.profile.Hash()
}
func (s *Warranty) allowed(p identity.Principal, permission string) bool {
	if s == nil || s.pool == nil || !s.profile.Valid() {
		return false
	}
	tenant, org := s.profile.Scope()
	return approvalPrincipal(p, tenant, org, permission)
}
func warrantyError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return wc.ErrNotFound
	}
	if postgresConflict(err) {
		return wc.ErrConflict
	}
	return err
}
func warrantyEvent(ctx context.Context, tx pgx.Tx, tenant, id, event string, version int64, payload any) error {
	raw, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)
 values($1,gen_random_uuid(),'warranty',$2,$3,$4,1,clock_timestamp(),$5)`, tenant, id, version, event, raw)
	return err
}
func readWarrantyOffer(ctx context.Context, q interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}, tenant, org, id string) (wc.Offer, error) {
	var out wc.Offer
	err := q.QueryRow(ctx, `select o.quotation_id,o.quotation_version,o.customer_principal_id,o.profile_sha256,o.profile_bytes,o.offered_by,o.offered_at,
 a.quotation_id is not null,coalesce(a.evidence_sha256,'') from service_ops.warranty_offer o
 left join service_ops.warranty_acknowledgement a on a.tenant_id=o.tenant_id and a.quotation_id=o.quotation_id
 where o.tenant_id=$1 and o.organization_id=$2 and o.quotation_id=$3`, tenant, org, id).Scan(&out.QuoteID, &out.QuoteVersion, &out.CustomerSubject, &out.ProfileSHA256, &out.Profile, &out.OfferedBy, &out.OfferedAt, &out.Acknowledged, &out.EvidenceSHA256)
	if err != nil {
		return out, warrantyError(err)
	}
	profile, err := wc.LoadProfile(out.Profile, out.ProfileSHA256)
	if err != nil {
		return wc.Offer{}, err
	}
	pt, po := profile.Scope()
	if pt != tenant || po != org {
		return wc.Offer{}, wc.ErrConflict
	}
	return out, nil
}
func (s *Warranty) Offer(ctx context.Context, p identity.Principal, id string) (wc.Offer, error) {
	if !wc.ValidID(id) || !(s.allowed(p, "warranty:read") || s.allowed(p, "warranty:self")) {
		return wc.Offer{}, wc.ErrInvalid
	}
	tenant, org := s.profile.Scope()
	out, err := readWarrantyOffer(ctx, s.pool, tenant, org, id)
	if err == nil && !s.allowed(p, "warranty:read") && out.CustomerSubject != p.Subject {
		return wc.Offer{}, wc.ErrNotFound
	}
	return out, err
}
func (s *Warranty) BindOffer(ctx context.Context, p identity.Principal, r wc.OfferRequest) (wc.Offer, error) {
	if !s.allowed(p, "warranty:offer") || !wc.ValidID(r.QuoteID) || r.QuoteVersion < 1 || r.ProfileSHA256 != s.profile.Hash() {
		return wc.Offer{}, wc.ErrInvalid
	}
	tenant, org := s.profile.Scope()
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return wc.Offer{}, err
	}
	defer tx.Rollback(ctx)
	var customer, state string
	var version int64
	var current bool
	err = tx.QueryRow(ctx, `select coalesce(customer_principal_id,''),state,version,valid_until>clock_timestamp() from sales.quotation
 where tenant_id=$1 and organization_id=$2 and quotation_id=$3 for update`, tenant, org, r.QuoteID).Scan(&customer, &state, &version, &current)
	if err != nil {
		return wc.Offer{}, warrantyError(err)
	}
	existing, err := readWarrantyOffer(ctx, tx, tenant, org, r.QuoteID)
	if err == nil {
		if existing.OfferedBy != p.Subject || existing.ProfileSHA256 != r.ProfileSHA256 || existing.QuoteVersion-1 != r.QuoteVersion {
			return wc.Offer{}, wc.ErrConflict
		}
		existing.Replay = true
		return existing, tx.Commit(ctx)
	}
	if !errors.Is(err, wc.ErrNotFound) {
		return wc.Offer{}, err
	}
	if customer == "" || state != "issued" || version != r.QuoteVersion || !current {
		return wc.Offer{}, wc.ErrConflict
	}
	if _, err = tx.Exec(ctx, `update sales.quotation set version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and quotation_id=$2`, tenant, r.QuoteID); err != nil {
		return wc.Offer{}, err
	}
	_, err = tx.Exec(ctx, `insert into service_ops.warranty_offer(tenant_id,quotation_id,organization_id,customer_principal_id,quotation_version,profile_sha256,profile_bytes,offered_by)
 values($1,$2,$3,$4,$5,$6,$7,$8)`, tenant, r.QuoteID, org, customer, version+1, r.ProfileSHA256, s.profile.Bytes(), p.Subject)
	if err != nil {
		return wc.Offer{}, warrantyError(err)
	}
	out, err := readWarrantyOffer(ctx, tx, tenant, org, r.QuoteID)
	if err != nil {
		return wc.Offer{}, err
	}
	if err = warrantyEvent(ctx, tx, tenant, r.QuoteID, "warranty.terms-offered", version+1, map[string]any{"organization_id": org, "quote_id": r.QuoteID, "quote_version": out.QuoteVersion, "profile_sha256": r.ProfileSHA256, "actor": p.Subject}); err != nil {
		return wc.Offer{}, err
	}
	return out, tx.Commit(ctx)
}
func (s *Warranty) Acknowledge(ctx context.Context, p identity.Principal, r wc.AcknowledgeRequest) (wc.Offer, error) {
	if !s.allowed(p, "warranty:self") || !wc.ValidID(r.QuoteID) || r.QuoteVersion < 2 || !wc.ValidSHA(r.ProfileSHA256) || !wc.ValidSHA(r.EvidenceSHA256) {
		return wc.Offer{}, wc.ErrInvalid
	}
	tenant, org := s.profile.Scope()
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return wc.Offer{}, err
	}
	defer tx.Rollback(ctx)
	var customer, state string
	var version int64
	var current bool
	err = tx.QueryRow(ctx, `select coalesce(customer_principal_id,''),state,version,valid_until>clock_timestamp() from sales.quotation
 where tenant_id=$1 and organization_id=$2 and quotation_id=$3 and customer_principal_id=$4 for update`, tenant, org, r.QuoteID, p.Subject).Scan(&customer, &state, &version, &current)
	if err != nil {
		return wc.Offer{}, warrantyError(err)
	}
	offer, err := readWarrantyOffer(ctx, tx, tenant, org, r.QuoteID)
	if err != nil {
		return wc.Offer{}, err
	}
	if offer.CustomerSubject != p.Subject || offer.ProfileSHA256 != r.ProfileSHA256 || offer.QuoteVersion != r.QuoteVersion {
		return wc.Offer{}, wc.ErrConflict
	}
	if offer.Acknowledged {
		if offer.EvidenceSHA256 != r.EvidenceSHA256 {
			return wc.Offer{}, wc.ErrConflict
		}
		offer.Replay = true
		return offer, tx.Commit(ctx)
	}
	if state != "issued" || !current || version != r.QuoteVersion {
		return wc.Offer{}, wc.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into service_ops.warranty_acknowledgement(tenant_id,quotation_id,quotation_version,customer_principal_id,profile_sha256,evidence_sha256)
 values($1,$2,$3,$4,$5,$6)`, tenant, r.QuoteID, r.QuoteVersion, p.Subject, r.ProfileSHA256, r.EvidenceSHA256)
	if err != nil {
		return wc.Offer{}, warrantyError(err)
	}
	offer.Acknowledged = true
	offer.EvidenceSHA256 = r.EvidenceSHA256
	if err = warrantyEvent(ctx, tx, tenant, r.QuoteID, "warranty.terms-acknowledged", version, map[string]any{"organization_id": org, "quote_id": r.QuoteID, "quote_version": version, "profile_sha256": r.ProfileSHA256, "evidence_sha256": r.EvidenceSHA256, "actor": p.Subject}); err != nil {
		return wc.Offer{}, err
	}
	return offer, tx.Commit(ctx)
}
func readWarrantyActivation(ctx context.Context, q interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}, tenant, org, handover string) (wc.Activation, error) {
	var a wc.Activation
	err := q.QueryRow(ctx, `select warranty_id,handover_id,quotation_id,order_id,stock_unit_id,organization_id,customer_principal_id,profile_sha256,terms_version,
 parts_start::text,parts_end::text,labor_start::text,labor_end::text,handover_accepted_at,activated_by,activated_at from service_ops.warranty_activation
 where tenant_id=$1 and organization_id=$2 and handover_id=$3`, tenant, org, handover).Scan(&a.WarrantyID, &a.HandoverID, &a.QuoteID, &a.OrderID, &a.StockUnitID, &a.OrganizationID, &a.CustomerSubject, &a.ProfileSHA256, &a.TermsVersion, &a.Dates.PartsStart, &a.Dates.PartsEnd, &a.Dates.LaborStart, &a.Dates.LaborEnd, &a.AcceptedAt, &a.ActivatedBy, &a.ActivatedAt)
	return a, warrantyError(err)
}
func (s *Warranty) Activation(ctx context.Context, p identity.Principal, handover string) (wc.Activation, error) {
	if !wc.ValidID(handover) || !(s.allowed(p, "warranty:read") || s.allowed(p, "warranty:self")) {
		return wc.Activation{}, wc.ErrInvalid
	}
	tenant, org := s.profile.Scope()
	out, err := readWarrantyActivation(ctx, s.pool, tenant, org, handover)
	if err == nil && !s.allowed(p, "warranty:read") && out.CustomerSubject != p.Subject {
		return wc.Activation{}, wc.ErrNotFound
	}
	return out, err
}
func (s *Warranty) Activate(ctx context.Context, p identity.Principal, handover string) (wc.Activation, error) {
	var out wc.Activation
	if !wc.ValidID(handover) || !s.allowed(p, "warranty:activate") {
		return out, wc.ErrInvalid
	}
	tenant, org := s.profile.Scope()
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return out, err
	}
	defer tx.Rollback(ctx)
	var state string
	err = tx.QueryRow(ctx, `select handover_id,order_id,stock_unit_id,customer_principal_id,state,customer_accepted_at from sales.delivery_handover
 where tenant_id=$1 and organization_id=$2 and handover_id=$3 and state='accepted' for update`, tenant, org, handover).Scan(&out.HandoverID, &out.OrderID, &out.StockUnitID, &out.CustomerSubject, &state, &out.AcceptedAt)
	if err != nil {
		return out, warrantyError(err)
	}
	existing, err := readWarrantyActivation(ctx, tx, tenant, org, handover)
	if err == nil {
		if existing.ActivatedBy != p.Subject {
			return out, wc.ErrConflict
		}
		existing.Replay = true
		return existing, tx.Commit(ctx)
	}
	if !errors.Is(err, wc.ErrNotFound) {
		return out, err
	}
	err = tx.QueryRow(ctx, `select qa.quotation_id from sales.quotation_acceptance qa where qa.tenant_id=$1 and qa.order_id=$2 and qa.customer_principal_id=$3`, tenant, out.OrderID, out.CustomerSubject).Scan(&out.QuoteID)
	if err != nil {
		return out, warrantyError(err)
	}
	offer, err := readWarrantyOffer(ctx, tx, tenant, org, out.QuoteID)
	if err != nil {
		return out, err
	}
	if !offer.Acknowledged || offer.CustomerSubject != out.CustomerSubject {
		return out, wc.ErrConflict
	}
	sold, err := wc.LoadProfile(offer.Profile, offer.ProfileSHA256)
	if err != nil {
		return out, err
	}
	dates, err := sold.DatesAt(out.AcceptedAt)
	if err != nil {
		return out, err
	}
	// service_ops.warranty uses instants; the declared date intervals are inclusive.
	// Store the following local midnight as the exclusive outer bound.
	doc := sold.Document()
	zone, err := time.LoadLocation(doc.BusinessTimeZone)
	if err != nil {
		return out, err
	}
	start, err := time.ParseInLocation(time.DateOnly, dates.PartsStart, zone)
	if err != nil {
		return out, err
	}
	last := dates.PartsEnd
	if dates.LaborEnd > last {
		last = dates.LaborEnd
	}
	end, err := time.ParseInLocation(time.DateOnly, last, zone)
	if err != nil {
		return out, err
	}
	end = end.AddDate(0, 0, 1)
	if end.Year() > 9999 {
		return out, wc.ErrInvalid
	}
	out.WarrantyID = wc.StableID("warranty", tenant, handover)
	_, err = tx.Exec(ctx, `insert into service_ops.warranty(tenant_id,warranty_id,stock_unit_id,customer_principal_id,starts_at,ends_at,terms_version,status)
 values($1,$2,$3,$4,$5,$6,$7,'active')`, tenant, out.WarrantyID, out.StockUnitID, out.CustomerSubject, start, end, doc.TermsVersion)
	if err != nil {
		return out, warrantyError(err)
	}
	_, err = tx.Exec(ctx, `insert into service_ops.warranty_activation(tenant_id,warranty_id,handover_id,quotation_id,order_id,organization_id,customer_principal_id,stock_unit_id,profile_sha256,terms_version,
 parts_start,parts_end,labor_start,labor_end,handover_accepted_at,activated_by)
 values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11::date,$12::date,$13::date,$14::date,$15,$16)`, tenant, out.WarrantyID, handover, out.QuoteID, out.OrderID, org, out.CustomerSubject, out.StockUnitID, sold.Hash(), doc.TermsVersion, dates.PartsStart, dates.PartsEnd, dates.LaborStart, dates.LaborEnd, out.AcceptedAt, p.Subject)
	if err != nil {
		return out, warrantyError(err)
	}
	out, err = readWarrantyActivation(ctx, tx, tenant, org, handover)
	if err != nil {
		return out, err
	}
	if err = warrantyEvent(ctx, tx, tenant, out.WarrantyID, "warranty.activated", 1, out); err != nil {
		return out, err
	}
	return out, tx.Commit(ctx)
}
