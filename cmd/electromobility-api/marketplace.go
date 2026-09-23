package main

// AUTHORED activation glue. Credentials stay in the process environment and
// provider calls use the fixed official origin with no redirect or proxy.
import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"strings"

	cr "elite.local/enterprise/internal/catalogrelease"
	mb "elite.local/enterprise/internal/marketplacebridge"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

var errMarketplaceConfiguration = errors.New("marketplace activation invalid")

type marketplaceToken struct{ lookup func(string) string }

func (t marketplaceToken) MercadoLibreAccessToken(context.Context) (string, error) {
	return t.lookup("MERCADOLIBRE_ACCESS_TOKEN"), nil
}
func init() { marketplaceModuleFactory = selectedMarketplaceModule }
func selectedMarketplaceModule(ctx context.Context, pool *pgxpool.Pool, price *postgres.Commerce, lookup func(string) string) (httpapi.EnterpriseModule, error) {
	if lookup == nil {
		return nil, errMarketplaceConfiguration
	}
	switch lookup("MARKETPLACE_ENABLED") {
	case "", "false":
		return nil, nil
	case "true":
	default:
		return nil, errMarketplaceConfiguration
	}
	if lookup("CATALOG_RELEASE_ENABLED") != "true" || pool == nil || price == nil {
		return nil, errMarketplaceConfiguration
	}
	f, e := os.Open(lookup("MARKETPLACE_PROFILE_FILE"))
	if e != nil {
		return nil, errMarketplaceConfiguration
	}
	defer f.Close()
	raw, e := io.ReadAll(io.LimitReader(f, 32769))
	if e != nil || len(raw) > 32768 {
		return nil, errMarketplaceConfiguration
	}
	sum := sha256.Sum256(raw)
	if hex.EncodeToString(sum[:]) != lookup("MARKETPLACE_PROFILE_SHA256") {
		return nil, errMarketplaceConfiguration
	}
	var profile mb.Profile
	if mb.Decode(raw, &profile) != nil || profile.Validate() != nil {
		return nil, errMarketplaceConfiguration
	}
	cp, e := cr.LoadProfile(lookup("CATALOG_RELEASE_PROFILE_FILE"), lookup("CATALOG_RELEASE_PROFILE_SHA256"))
	if e != nil {
		return nil, errMarketplaceConfiguration
	}
	catalog, e := postgres.NewCatalogRelease(pool, cp, price)
	if e != nil {
		return nil, errMarketplaceConfiguration
	}
	key, e := base64.StdEncoding.DecodeString(lookup("MARKETPLACE_HMAC_KEY"))
	if e != nil || len(key) != 32 {
		return nil, errMarketplaceConfiguration
	}
	token := lookup("MERCADOLIBRE_ACCESS_TOKEN")
	if len(token) < 20 || len(token) > 4096 || strings.ContainsAny(token, "\r\n") {
		return nil, errMarketplaceConfiguration
	}
	var ready bool
	e = pool.QueryRow(ctx, `select exists(select 1 from pg_trigger where tgrelid=to_regclass('catalog.marketplace_effect') and tgname='marketplace_effect_immutable' and tgenabled in('O','A'))
 and exists(select 1 from pg_trigger where tgrelid=to_regclass('catalog.marketplace_observation') and tgname='marketplace_observation_immutable' and tgenabled in('O','A'))
 and exists(select 1 from pg_attribute where attrelid=to_regclass('catalog.marketplace_effect') and attname='operation' and not attisdropped)
 and exists(select 1 from pg_constraint where conrelid=to_regclass('catalog.marketplace_effect') and conname='marketplace_effect_operation_check' and pg_get_constraintdef(oid) like '%CONTENT%')
 and exists(select 1 from pg_constraint where conrelid=to_regclass('approval.request') and conname='request_kind_check' and pg_get_constraintdef(oid) like '%marketplace_mutation%')
 and (select count(*) from pg_trigger where not tgisinternal and tgenabled in('O','A') and
 (tgrelid=to_regclass('approval.request') and tgname='bound_request_immutable' or tgrelid=to_regclass('approval.decision') and tgname='bound_decision_immutable'))=2
 and pg_get_functiondef(to_regprocedure('approval.guard_bound_request()')) like '%marketplace_mutation%'
 and pg_get_functiondef(to_regprocedure('approval.guard_bound_decision()')) like '%marketplace_mutation%'`).Scan(&ready)
	if e != nil || !ready {
		return nil, errMarketplaceConfiguration
	}
	service, e := postgres.NewMarketplacePublication(catalog, profile, key, marketplaceToken{lookup}, mb.NewHTTPClient())
	if e != nil {
		return nil, errMarketplaceConfiguration
	}
	return httpapi.MarketplaceModule{Service: service, TenantID: profile.TenantID, OrganizationID: profile.OrganizationID}, nil
}
