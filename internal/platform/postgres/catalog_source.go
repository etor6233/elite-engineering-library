package postgres

// AUTHORED transaction adapters; model and pricing validation/SQL remain in
// their original services. No alternate pricing or publication algorithm.
import (
	"context"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/commerce"
	"elite.local/enterprise/internal/electromobility"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/randomid"
	"github.com/jackc/pgx/v5"
)

type catalogModelTx struct {
	*Electromobility
	tx pgx.Tx
}

func (r catalogModelTx) CreateModel(ctx context.Context, tenant, event string, m electromobility.Model) error {
	return r.createModelTx(ctx, r.tx, tenant, event, m)
}

type catalogPriceTx struct {
	*Commerce
	tx pgx.Tx
}

func (r catalogPriceTx) CreatePriceBook(ctx context.Context, tenant, event string, b commerce.PriceBook) error {
	return r.createPriceBookTx(ctx, r.tx, tenant, event, b)
}

func (s *CatalogRelease) CreateSource(ctx context.Context, p identity.Principal, in cr.SourceRequest) (cr.Receipt, error) {
	if !in.Valid() {
		return cr.Receipt{}, cr.ErrInvalid
	}
	_, hash, e := cr.Canonical(in)
	if e != nil {
		return cr.Receipt{}, e
	}
	tx, e := s.begin(ctx, p, "catalog:draft")
	if e != nil {
		return cr.Receipt{}, e
	}
	defer tx.Rollback(ctx)
	if r, yes, e := s.replay(ctx, tx, p, in.CommandID, "source", hash); e != nil || yes {
		return r, e
	}
	ids := randomid.Generator{}
	model, e := electromobility.NewService(catalogModelTx{NewElectromobility(s.pool), tx}, ids).CreateModel(ctx, p.TenantID, in.Model)
	if e != nil {
		return cr.Receipt{}, e
	}
	entries := []commerce.PriceEntry{}
	variantIDs := []string{}
	for _, v := range in.Variants {
		id := ids.New()
		_, e = tx.Exec(ctx, `insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)
  values($1,$2,$3,$4,$5,$6,'draft')`, p.TenantID, id, model.ID, v.Code, v.DisplayName, v.BatterySpecification)
		if e != nil {
			return cr.Receipt{}, e
		}
		variantIDs = append(variantIDs, id)
		entries = append(entries, commerce.PriceEntry{VariantID: id, AmountMinorUnits: v.AmountMinorUnits, TaxMode: v.TaxMode})
	}
	book, e := commerce.NewService(catalogPriceTx{s.price, tx}, ids).CreatePriceBook(ctx, p.TenantID, commerce.PriceBook{Market: s.profile.Market, Currency: s.profile.Currency, ValidFrom: in.ValidFrom, ValidUntil: in.ValidUntil, Entries: entries})
	if e != nil {
		return cr.Receipt{}, e
	}
	r := cr.Receipt{CommandID: in.CommandID, Actor: p.Subject, Kind: "source", ResourceID: model.ID, RequestSHA256: hash, SourcePriceBookID: book.ID, SourceVariantIDs: variantIDs}
	if e = s.save(ctx, tx, p, r); e != nil {
		return cr.Receipt{}, e
	}
	return r, tx.Commit(ctx)
}

func (s *CatalogRelease) Current(ctx context.Context, p identity.Principal) (cr.Publication, error) {
	if !s.authorized(p, "catalog:read") {
		return cr.Publication{}, cr.ErrInvalid
	}
	return s.Public(ctx)
}
