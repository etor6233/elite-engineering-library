package postgres

import (
	"context"
	"elite.local/enterprise/internal/approval"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/randomid"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"net/url"
	"sort"
)

func (s *CatalogRelease) UploadPNG(ctx context.Context, p identity.Principal, key string, raw []byte) (cr.Receipt, error) {
	if !cr.ValidID(key) || !s.authorized(p, "catalog:draft") {
		return cr.Receipt{}, cr.ErrInvalid
	}
	media, clean, err := cr.NormalizePNG(raw)
	if err != nil {
		return cr.Receipt{}, err
	}
	_, hash, err := cr.Canonical(map[string]string{"command_id": key, "original_sha256": media.OriginalSHA256})
	if err != nil {
		return cr.Receipt{}, err
	}
	tx, err := s.begin(ctx, p, "catalog:draft")
	if err != nil {
		return cr.Receipt{}, err
	}
	defer tx.Rollback(ctx)
	if r, yes, err := s.replay(ctx, tx, p, key, "media", hash); err != nil || yes {
		return r, err
	}
	media.ID = (randomid.Generator{}).New()
	_, err = tx.Exec(ctx, `insert into catalog.release_media(tenant_id,media_id,maker,original_sha256,sha256,width,height,png) values($1,$2,$3,$4,$5,$6,$7,$8)`,
		p.TenantID, media.ID, p.Subject, media.OriginalSHA256, media.SHA256, media.Width, media.Height, clean)
	if err != nil {
		return cr.Receipt{}, err
	}
	r := cr.Receipt{CommandID: key, Actor: p.Subject, Kind: "media", ResourceID: media.ID, RequestSHA256: hash}
	if err = s.save(ctx, tx, p, r); err != nil {
		return cr.Receipt{}, err
	}
	return r, tx.Commit(ctx)
}
func (s *CatalogRelease) CreateDraft(ctx context.Context, p identity.Principal, input cr.DraftRequest) (cr.Receipt, error) {
	if !input.Valid() {
		return cr.Receipt{}, cr.ErrInvalid
	}
	_, hash, err := cr.Canonical(input)
	if err != nil {
		return cr.Receipt{}, err
	}
	tx, err := s.begin(ctx, p, "catalog:draft")
	if err != nil {
		return cr.Receipt{}, err
	}
	defer tx.Rollback(ctx)
	if r, yes, err := s.replay(ctx, tx, p, input.CommandID, "draft", hash); err != nil || yes {
		return r, err
	}
	snap := cr.Snapshot{Schema: "elite-catalog-snapshot/v1", Profile: s.profile, PriceBookID: input.PriceBookID, Models: []cr.Model{}, Variants: []cr.Variant{}}
	if err = tx.QueryRow(ctx, `select tenant_code from platform.tenant where tenant_id=$1`, p.TenantID).Scan(&snap.TenantCode); err != nil {
		return cr.Receipt{}, err
	}
	err = tx.QueryRow(ctx, `select valid_from,valid_until from pricing.price_book where tenant_id=$1 and price_book_id=$2 and market=$3 and currency=$4 for update`, p.TenantID, input.PriceBookID, s.profile.Market, s.profile.Currency).Scan(&snap.ValidFrom, &snap.ValidUntil)
	if errors.Is(err, pgx.ErrNoRows) {
		return cr.Receipt{}, cr.ErrNotFound
	}
	if err != nil {
		return cr.Receipt{}, err
	}
	selected := append([]cr.ModelReference(nil), input.Models...)
	sort.Slice(selected, func(i, j int) bool { return selected[i].ModelID < selected[j].ModelID })
	modelIDs := map[string]bool{}
	for _, ref := range selected {
		var m cr.Model
		err = tx.QueryRow(ctx, `select model_id,model_code,display_name,vehicle_class,specification from catalog.vehicle_model
    where tenant_id=$1 and model_id=$2 and lifecycle_state in('draft','active') and octet_length(display_name)<=160 and pg_column_size(specification)<=32768 for share`, p.TenantID, ref.ModelID).Scan(&m.ID, &m.Code, &m.DisplayName, &m.VehicleClass, &m.Specification)
		if errors.Is(err, pgx.ErrNoRows) {
			return cr.Receipt{}, cr.ErrNotFound
		}
		if err != nil {
			return cr.Receipt{}, err
		}
		err = tx.QueryRow(ctx, `select media_id,original_sha256,sha256,width,height from catalog.release_media where tenant_id=$1 and media_id=$2`, p.TenantID, ref.MediaID).Scan(&m.Media.ID, &m.Media.OriginalSHA256, &m.Media.SHA256, &m.Media.Width, &m.Media.Height)
		if errors.Is(err, pgx.ErrNoRows) {
			return cr.Receipt{}, cr.ErrNotFound
		}
		if err != nil {
			return cr.Receipt{}, err
		}
		m.CanonicalURL = s.profile.Origin + "/models/" + url.PathEscape(m.Code)
		snap.Models = append(snap.Models, m)
		modelIDs[m.ID] = false
	}
	rows, err := tx.Query(ctx, `select v.variant_id,v.model_id,v.variant_code,v.display_name,v.battery_specification,v.homologation_state,e.amount_minor_units,e.tax_mode
 from pricing.price_book_entry e join catalog.vehicle_variant v using(tenant_id,variant_id)
 where e.tenant_id=$1 and e.price_book_id=$2 and v.lifecycle_state in('draft','active') and pg_column_size(v.battery_specification)<=32768
 order by v.variant_id limit 129 for share of e,v`, p.TenantID, input.PriceBookID)
	if err != nil {
		return cr.Receipt{}, err
	}
	for rows.Next() {
		var v cr.Variant
		if err = rows.Scan(&v.ID, &v.ModelID, &v.Code, &v.DisplayName, &v.BatterySpecification, &v.HomologationState, &v.AmountMinorUnits, &v.TaxMode); err != nil {
			rows.Close()
			return cr.Receipt{}, err
		}
		if _, ok := modelIDs[v.ModelID]; !ok {
			rows.Close()
			return cr.Receipt{}, cr.ErrInvalid
		}
		modelIDs[v.ModelID] = true
		snap.Variants = append(snap.Variants, v)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return cr.Receipt{}, err
	}
	if len(snap.Variants) == 0 || len(snap.Variants) > 128 {
		return cr.Receipt{}, cr.ErrInvalid
	}
	for _, present := range modelIDs {
		if !present {
			return cr.Receipt{}, cr.ErrInvalid
		}
	}
	var entryCount int
	if err = tx.QueryRow(ctx, `select count(*) from pricing.price_book_entry where tenant_id=$1 and price_book_id=$2`, p.TenantID, input.PriceBookID).Scan(&entryCount); err != nil {
		return cr.Receipt{}, err
	}
	if entryCount != len(snap.Variants) {
		return cr.Receipt{}, cr.ErrInvalid
	}
	raw, sha, err := cr.Canonical(snap)
	if err != nil || len(raw) > 262144 {
		return cr.Receipt{}, cr.ErrInvalid
	}
	id := (randomid.Generator{}).New()
	_, err = tx.Exec(ctx, `insert into catalog.release_draft(tenant_id,draft_id,maker,price_book_id,snapshot,snapshot_canonical,snapshot_sha256) values($1,$2,$3,$4,$5,$6,$7)`, p.TenantID, id, p.Subject, input.PriceBookID, raw, string(raw), sha)
	if err != nil {
		return cr.Receipt{}, err
	}
	approvals := NewHumanApprovals(s.pool)
	for _, stage := range []string{"legal", "technical", "media", "publication"} {
		payload, h, err := approval.CanonicalPayload([]byte(`{"draft_id":"` + id + `","snapshot_sha256":"` + sha + `","stage":"` + stage + `"}`))
		if err != nil {
			return cr.Receipt{}, err
		}
		aid := (randomid.Generator{}).New()
		_, err = tx.Exec(ctx, `insert into catalog.release_review(tenant_id,draft_id,stage,approval_id,payload_sha256) values($1,$2,$3,$4,$5)`, p.TenantID, id, stage, aid, h)
		if err != nil {
			return cr.Receipt{}, err
		}
		spec := HumanApprovalSpec{Request: approval.Request{TenantID: p.TenantID, ID: aid, Kind: approval.KindCatalogReview, SubjectID: id, Requester: p.Subject, EvidenceSHA: h}, OrganizationID: s.profile.OrganizationID, Payload: payload}
		if _, err = approvals.submitTx(ctx, tx, p, spec, "catalog:draft", nil); err != nil {
			return cr.Receipt{}, err
		}
	}
	r := cr.Receipt{CommandID: input.CommandID, Actor: p.Subject, Kind: "draft", ResourceID: id, RequestSHA256: hash, SnapshotSHA256: sha}
	if err = s.save(ctx, tx, p, r); err != nil {
		return cr.Receipt{}, err
	}
	return r, tx.Commit(ctx)
}
func (s *CatalogRelease) Draft(ctx context.Context, p identity.Principal, id string) (cr.Draft, error) {
	if !s.authorized(p, "catalog:read") || !cr.ValidID(id) {
		return cr.Draft{}, cr.ErrNotFound
	}
	var d cr.Draft
	err := s.pool.QueryRow(ctx, `select draft_id,maker,snapshot_sha256,snapshot_canonical from catalog.release_draft where tenant_id=$1 and draft_id=$2`, p.TenantID, id).Scan(&d.ID, &d.Maker, &d.SHA256, &d.Snapshot)
	if errors.Is(err, pgx.ErrNoRows) {
		return d, cr.ErrNotFound
	}
	if err != nil {
		return d, err
	}
	var raw []byte
	err = s.pool.QueryRow(ctx, `select jsonb_object_agg(r.stage,q.state) from catalog.release_review r join approval.request q on q.tenant_id=r.tenant_id and q.request_id=r.approval_id where r.tenant_id=$1 and r.draft_id=$2`, p.TenantID, id).Scan(&raw)
	if err != nil {
		return d, err
	}
	err = json.Unmarshal(raw, &d.Reviews)
	return d, err
}
