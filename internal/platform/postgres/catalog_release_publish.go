package postgres

import (
	"context"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/commerce"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/randomid"
	"elite.local/enterprise/internal/search"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"math"
)

func (s *CatalogRelease) Review(ctx context.Context, p identity.Principal, input cr.ReviewRequest) (cr.Receipt, error) {
	if !input.Valid() {
		return cr.Receipt{}, cr.ErrInvalid
	}
	permission := "catalog:review:" + input.Stage
	tx, err := s.begin(ctx, p, permission)
	if err != nil {
		return cr.Receipt{}, err
	}
	defer tx.Rollback(ctx)
	_, hash, err := cr.Canonical(input)
	if err != nil {
		return cr.Receipt{}, err
	}
	if r, yes, err := s.replay(ctx, tx, p, input.CommandID, "review", hash); err != nil || yes {
		return r, err
	}
	var aid, payloadSHA, maker string
	err = tx.QueryRow(ctx, `select r.approval_id,r.payload_sha256,d.maker from catalog.release_review r join catalog.release_draft d using(tenant_id,draft_id)
 where r.tenant_id=$1 and r.draft_id=$2 and r.stage=$3 and d.snapshot_sha256=$4`, p.TenantID, input.DraftID, input.Stage, input.SnapshotSHA256).Scan(&aid, &payloadSHA, &maker)
	if errors.Is(err, pgx.ErrNoRows) {
		return cr.Receipt{}, cr.ErrNotFound
	}
	if err != nil {
		return cr.Receipt{}, err
	}
	if maker == p.Subject {
		return cr.Receipt{}, cr.ErrConflict
	}
	if input.Stage == "publication" && input.Approved {
		var count int
		err = tx.QueryRow(ctx, `select count(*) from catalog.release_review r join approval.request q on q.tenant_id=r.tenant_id and q.request_id=r.approval_id
   where r.tenant_id=$1 and r.draft_id=$2 and r.stage<>'publication' and q.state='approved'`, p.TenantID, input.DraftID).Scan(&count)
		if err != nil {
			return cr.Receipt{}, err
		}
		if count != 3 {
			return cr.Receipt{}, cr.ErrConflict
		}
	}
	_, err = NewHumanApprovals(s.pool).decideTx(ctx, tx, p, p.TenantID, aid, s.profile.OrganizationID, payloadSHA, input.Approved, input.Reason, permission, nil)
	if err != nil {
		return cr.Receipt{}, err
	}
	_, err = tx.Exec(ctx, `insert into catalog.release_review_action(tenant_id,command_id,approval_id,actor,approved,evidence_sha256) values($1,$2,$3,$4,$5,$6)`,
		p.TenantID, input.CommandID, aid, p.Subject, input.Approved, input.EvidenceSHA256)
	if err != nil {
		return cr.Receipt{}, err
	}
	r := cr.Receipt{CommandID: input.CommandID, Actor: p.Subject, Kind: "review", ResourceID: input.DraftID, RequestSHA256: hash, SnapshotSHA256: input.SnapshotSHA256}
	if err = s.save(ctx, tx, p, r); err != nil {
		return cr.Receipt{}, err
	}
	return r, tx.Commit(ctx)
}
func (s *CatalogRelease) Publish(ctx context.Context, p identity.Principal, input cr.PublishRequest) (cr.Receipt, error) {
	if !input.Valid() {
		return cr.Receipt{}, cr.ErrInvalid
	}
	tx, err := s.begin(ctx, p, "catalog:publish")
	if err != nil {
		return cr.Receipt{}, err
	}
	defer tx.Rollback(ctx)
	_, hash, err := cr.Canonical(input)
	if err != nil {
		return cr.Receipt{}, err
	}
	if r, yes, err := s.replay(ctx, tx, p, input.CommandID, "publish", hash); err != nil || yes {
		return r, err
	}
	var current int64
	if err = tx.QueryRow(ctx, `select coalesce(max(generation),0) from catalog.release_publication where tenant_id=$1`, p.TenantID).Scan(&current); err != nil {
		return cr.Receipt{}, err
	}
	if current != input.ExpectedGeneration || current == math.MaxInt64 {
		return cr.Receipt{}, cr.ErrConflict
	}
	var raw []byte
	var maker, book string
	err = tx.QueryRow(ctx, `select snapshot_canonical,maker,price_book_id from catalog.release_draft where tenant_id=$1 and draft_id=$2 and snapshot_sha256=$3`, p.TenantID, input.DraftID, input.SnapshotSHA256).Scan(&raw, &maker, &book)
	if errors.Is(err, pgx.ErrNoRows) {
		return cr.Receipt{}, cr.ErrNotFound
	}
	if err != nil {
		return cr.Receipt{}, err
	}
	if maker == p.Subject {
		return cr.Receipt{}, cr.ErrConflict
	}
	var approved int
	err = tx.QueryRow(ctx, `select count(*) from catalog.release_review r join approval.request q on q.tenant_id=r.tenant_id and q.request_id=r.approval_id
 where r.tenant_id=$1 and r.draft_id=$2 and q.state='approved'`, p.TenantID, input.DraftID).Scan(&approved)
	if err != nil {
		return cr.Receipt{}, err
	}
	if approved != 4 {
		return cr.Receipt{}, cr.ErrConflict
	}
	var eligible bool
	err = tx.QueryRow(ctx, `select valid_from<=clock_timestamp() and (valid_until is null or valid_until>clock_timestamp())
 from pricing.price_book where tenant_id=$1 and price_book_id=$2 for update`, p.TenantID, book).Scan(&eligible)
	if err != nil {
		return cr.Receipt{}, err
	}
	if !eligible {
		return cr.Receipt{}, cr.ErrConflict
	}
	var snapshot cr.Snapshot
	if err = json.Unmarshal(raw, &snapshot); err != nil {
		return cr.Receipt{}, err
	}
	// Every publication gets a new immutable effective book. The source owner
	// permits one activation per book; rollback clones the approved historical
	// amounts and validity without rewriting either history or source events.
	effective := commerce.PriceBook{ID: (randomid.Generator{}).New(), Market: snapshot.Profile.Market, Currency: snapshot.Profile.Currency, ValidFrom: snapshot.ValidFrom, ValidUntil: snapshot.ValidUntil}
	for _, v := range snapshot.Variants {
		effective.Entries = append(effective.Entries, commerce.PriceEntry{VariantID: v.ID, AmountMinorUnits: v.AmountMinorUnits, TaxMode: v.TaxMode})
	}
	if err = s.price.createPriceBookTx(ctx, tx, p.TenantID, (randomid.Generator{}).New(), effective); err != nil {
		return cr.Receipt{}, err
	}
	r := cr.Receipt{CommandID: input.CommandID, Actor: p.Subject, Kind: "publish", EffectivePriceBookID: effective.ID, ResourceID: input.DraftID, RequestSHA256: hash, SnapshotSHA256: input.SnapshotSHA256, Generation: current + 1}
	_, err = tx.Exec(ctx, `insert into catalog.release_publication(tenant_id,generation,draft_id,command_id,actor,reason,effective_price_book_id) values($1,$2,$3,$4,$5,$6,$7)`, p.TenantID, r.Generation, input.DraftID, input.CommandID, p.Subject, input.Reason, effective.ID)
	if err != nil {
		return cr.Receipt{}, err
	}
	// Status selection is publication glue; the original Commerce activation
	// still owns market/currency overlap, activation SQL and its outbox event.
	_, err = tx.Exec(ctx, `update pricing.price_book set status='retired' where tenant_id=$1 and market=$2 and currency=$3 and status='active'`, p.TenantID, s.profile.Market, s.profile.Currency)
	if err != nil {
		return cr.Receipt{}, err
	}
	if err = s.price.activatePriceBookTx(ctx, tx, p.TenantID, effective.ID, (randomid.Generator{}).New()); err != nil {
		return cr.Receipt{}, err
	}
	if _, err = tx.Exec(ctx, `delete from search.document where tenant_id=$1 and kind='catalog-model'`, p.TenantID); err != nil {
		return cr.Receipt{}, err
	}
	for _, m := range snapshot.Models {
		doc := search.Document{TenantID: p.TenantID, ID: "catalog:" + m.ID, Kind: "catalog-model", ExternalID: m.ID, Title: m.DisplayName, Body: m.CanonicalURL, ContentSHA256: input.SnapshotSHA256}
		if err = doc.Validate(); err != nil {
			return cr.Receipt{}, err
		}
		_, err = tx.Exec(ctx, `insert into search.document(tenant_id,document_id,kind,external_id,title,body,content_sha256) values($1,$2,$3,$4,$5,$6,$7)`,
			doc.TenantID, doc.ID, doc.Kind, doc.ExternalID, doc.Title, doc.Body, doc.ContentSHA256)
		if err != nil {
			return cr.Receipt{}, err
		}
	}
	if err = s.save(ctx, tx, p, r); err != nil {
		return cr.Receipt{}, err
	}
	return r, tx.Commit(ctx)
}
func (s *CatalogRelease) Public(ctx context.Context) (cr.Publication, error) {
	var out cr.Publication
	err := s.pool.QueryRow(ctx, `select generation,draft_id,snapshot_sha256,snapshot_canonical,effective_price_book_id from catalog.release_public_current where tenant_id=$1`, s.profile.TenantID).Scan(&out.Generation, &out.DraftID, &out.SHA256, &out.Snapshot, &out.EffectivePriceBookID)
	if errors.Is(err, pgx.ErrNoRows) {
		return out, cr.ErrNotFound
	}
	return out, err
}
