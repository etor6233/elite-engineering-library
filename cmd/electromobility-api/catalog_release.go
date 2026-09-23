package main

// AUTHORED optional host: the same selected Commerce owner supplies pricing.
import (
	"context"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/postgres"
	"encoding/base64"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

var errCatalogReleaseConfiguration = errors.New("catalog publication activation invalid")

func init() { catalogReleaseModuleFactory = selectedCatalogReleaseModule }
func selectedCatalogReleaseModule(ctx context.Context, pool *pgxpool.Pool, price *postgres.Commerce, lookup func(string) string) (httpapi.EnterpriseModule, error) {
	if lookup == nil {
		return nil, errCatalogReleaseConfiguration
	}
	switch lookup("CATALOG_RELEASE_ENABLED") {
	case "", "false":
		return nil, nil
	case "true":
	default:
		return nil, errCatalogReleaseConfiguration
	}
	if pool == nil || price == nil {
		return nil, errCatalogReleaseConfiguration
	}
	profile, e := cr.LoadProfile(lookup("CATALOG_RELEASE_PROFILE_FILE"), lookup("CATALOG_RELEASE_PROFILE_SHA256"))
	if e != nil {
		return nil, errCatalogReleaseConfiguration
	}
	var ready bool
	e = pool.QueryRow(ctx, `select (select count(*) from pg_trigger where not tgisinternal and tgenabled in('O','A') and
 (tgname='catalog_release_immutable' and tgrelid in(to_regclass('catalog.release_stream'),to_regclass('catalog.release_media'),
 to_regclass('catalog.release_draft'),to_regclass('catalog.release_review'),to_regclass('catalog.release_command'),
 to_regclass('catalog.release_review_action'),to_regclass('catalog.release_publication'))
 or tgname='catalog_release_approval_guard' and tgrelid=to_regclass('approval.request')
 or tgname='catalog_release_publication_guard' and tgrelid=to_regclass('catalog.release_publication')
 or tgname='catalog_release_price_guard' and tgrelid=to_regclass('pricing.price_book')
 or tgname='catalog_release_entry_guard' and tgrelid=to_regclass('pricing.price_book_entry')
 or tgname='catalog_release_search_guard' and tgrelid=to_regclass('search.document')))=12
 and to_regclass('catalog.release_public_models') is not null
 and exists(select 1 from pg_constraint where conrelid=to_regclass('catalog.release_command') and conname='release_command_kind_check' and pg_get_constraintdef(oid) like '%source%')
 and exists(select 1 from pg_constraint where conrelid=to_regclass('approval.request') and conname='request_kind_check' and pg_get_constraintdef(oid) like '%catalog_review%')`).Scan(&ready)
	if e != nil || !ready {
		return nil, errCatalogReleaseConfiguration
	}
	store, e := postgres.NewCatalogRelease(pool, profile, price)
	if e != nil {
		return nil, errCatalogReleaseConfiguration
	}
	module := httpapi.CatalogReleaseModule{Service: store, DeferPublicMedia: true}
	switch lookup("CATALOG_FEED_ENABLED") {
	case "", "false":
		return module, nil
	case "true":
	default:
		return nil, errCatalogReleaseConfiguration
	}
	feedProfile, e := cr.LoadFeedProfile(lookup("CATALOG_FEED_PROFILE_FILE"), lookup("CATALOG_FEED_PROFILE_SHA256"))
	if e != nil {
		return nil, errCatalogReleaseConfiguration
	}
	key, e := base64.StdEncoding.DecodeString(lookup("CATALOG_FEED_HMAC_KEY"))
	if e != nil || len(key) != 32 {
		return nil, errCatalogReleaseConfiguration
	}
	e = pool.QueryRow(ctx, `select exists(select 1 from pg_trigger where tgrelid=to_regclass('catalog.release_feed_intent') and tgname='catalog_release_feed_immutable' and tgenabled in('O','A'))`).Scan(&ready)
	if e != nil || !ready {
		return nil, errCatalogReleaseConfiguration
	}
	receiver, e := cr.NewFeedHTTP(feedProfile, lookup("CATALOG_FEED_TOKEN"))
	if e != nil {
		return nil, errCatalogReleaseConfiguration
	}
	feeder, e := postgres.NewCatalogFeed(store, key, receiver)
	if e != nil {
		receiver.Close()
		return nil, errCatalogReleaseConfiguration
	}
	module.Feed = feeder
	return module, nil
}
