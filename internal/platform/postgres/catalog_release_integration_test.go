package postgres_test

// AUTHORED actual PostgreSQL owner composition; synthetic principals and media.
import (
	"bytes"
	"context"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/commerce"
	"elite.local/enterprise/internal/electromobility"
	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"image"
	"image/png"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

type catalogReleaseReference interface {
	UploadPNG(context.Context, identity.Principal, string, []byte) (cr.Receipt, error)
	CreateDraft(context.Context, identity.Principal, cr.DraftRequest) (cr.Receipt, error)
	Draft(context.Context, identity.Principal, string) (cr.Draft, error)
	Review(context.Context, identity.Principal, cr.ReviewRequest) (cr.Receipt, error)
	Publish(context.Context, identity.Principal, cr.PublishRequest) (cr.Receipt, error)
	CommandReceipt(context.Context, identity.Principal, string) (cr.Receipt, error)
	Public(context.Context) (cr.Publication, error)
	PublicMedia(context.Context, string) ([]byte, error)
	PublicSearch(context.Context, string) ([]cr.SearchHit, error)
}

func TestCatalogReleaseConnectedReference(t *testing.T) { testCatalogReleaseConnected(t, nil) }
func testCatalogReleaseConnected(t *testing.T, hook func(*testing.T, *db.CatalogRelease) catalogReleaseReference) {

	rawURL := os.Getenv("PAYMENT_CONNECTED_DB_URL")
	if rawURL == "" {
		t.Skip("owned fixture DB not configured")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, rawURL)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := "ca6c9cb7-e548-4f19-a2e6-27d2aa2d4893"
	for _, q := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'catalog-j3','Fixture','Fixture')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'j3-store','j3-store','Fixture','store')`,
	} {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	profile := cr.Profile{Schema: "elite-catalog-publication/v1", TenantID: tenant, OrganizationID: "j3-store", Market: "AR", Currency: "ARS", Origin: "https://catalog.example.invalid", PolicyCode: cr.PolicyCode}
	price := db.NewCommerce(pool)
	backend, err := db.NewCatalogRelease(pool, profile, price)
	if err != nil {
		t.Fatal(err)
	}
	var store catalogReleaseReference = backend
	if hook != nil {
		store = hook(t, backend)
	}
	person := func(subject string) identity.Principal {
		return identity.Principal{TenantID: tenant, Subject: subject, Permissions: map[string]struct{}{"catalog:feed": {}, "catalog:draft": {}, "catalog:read": {}, "catalog:publish": {}, "catalog:review:legal": {}, "catalog:review:technical": {}, "catalog:review:media": {}, "catalog:review:publication": {}}, Organizations: map[string]struct{}{"j3-store": {}}}
	}
	maker, reviewer := person("j3-maker"), person("j3-reviewer")
	var pngBytes bytes.Buffer
	if err = png.Encode(&pngBytes, image.NewNRGBA(image.Rect(0, 0, 2, 2))); err != nil {
		t.Fatal(err)
	}
	media, err := store.UploadPNG(ctx, maker, "upload-png", append(pngBytes.Bytes(), []byte("<script>untrusted trailer</script>")...))
	if err != nil {
		t.Fatal("upload", err)
	}
	model, err := electromobility.NewService(db.NewElectromobility(pool), randomid.Generator{}).CreateModel(ctx, tenant, electromobility.Model{Code: "bicycle-one", DisplayName: "Version One", VehicleClass: "bicycle", Specification: json.RawMessage(`{"range_km":20}`)})
	if err != nil {
		t.Fatal(err)
	}
	_, err = pool.Exec(ctx, `insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state) values($1,'j3-variant',$2,'j3-variant','Fixture','{}','active')`, tenant, model.ID)
	if err != nil {
		t.Fatal(err)
	}
	sales := commerce.NewService(price, randomid.Generator{})
	book := func(amount int64, until *time.Time) commerce.PriceBook {
		t.Helper()
		var now time.Time
		if err := pool.QueryRow(ctx, `select clock_timestamp()`).Scan(&now); err != nil {
			t.Fatal(err)
		}
		out, err := sales.CreatePriceBook(ctx, tenant, commerce.PriceBook{Market: "AR", Currency: "ARS", ValidFrom: now.Add(-time.Hour), ValidUntil: until, Entries: []commerce.PriceEntry{{VariantID: "j3-variant", AmountMinorUnits: amount, TaxMode: "not-applicable"}}})
		if err != nil {
			t.Fatal("price book", err)
		}
		return out
	}
	firstBook := book(120000, nil)
	draft := func(key, bookID string) cr.Receipt {
		t.Helper()
		r, err := store.CreateDraft(ctx, maker, cr.DraftRequest{CommandID: key, PriceBookID: bookID, Models: []cr.ModelReference{{ModelID: model.ID, MediaID: media.ResourceID}}})
		if err != nil {
			t.Fatal("draft", err)
		}
		return r
	}
	one := draft("draft-one", firstBook.ID)
	snapshot := func() string {
		t.Helper()
		var out string
		for _, table := range []string{"catalog.release_stream", "catalog.release_media", "catalog.release_draft", "catalog.release_review", "catalog.release_review_action", "catalog.release_command", "catalog.release_publication", "pricing.price_book", "pricing.price_book_entry", "search.document", "approval.request", "approval.decision", "platform.outbox_event"} {
			var text string
			if err := pool.QueryRow(ctx, fmt.Sprintf("select coalesce(jsonb_agg(to_jsonb(x) order by to_jsonb(x)::text),'[]'::jsonb)::text from %s x where tenant_id=$1", table), tenant).Scan(&text); err != nil {
				t.Fatal(table, err)
			}
			out += text
		}
		return out
	}
	expectDenied := func(action func() error, name string) {
		t.Helper()
		before := snapshot()
		if err := action(); err == nil {
			t.Fatal(name + " accepted")
		}
		if snapshot() != before {
			t.Fatal(name + " leaked effects")
		}
	}
	requestReview := func(d cr.Receipt, stage, key string) cr.ReviewRequest {
		return cr.ReviewRequest{CommandID: key, DraftID: d.ResourceID, Stage: stage, SnapshotSHA256: d.SnapshotSHA256, Approved: true, Reason: "synthetic reviewer evidence; no legal certification", EvidenceSHA256: strings.Repeat("a", 64)}
	}
	expectDenied(func() error { _, e := store.Review(ctx, maker, requestReview(one, "legal", "self-review")); return e }, "self review")
	expectDenied(func() error {
		_, e := store.Review(ctx, reviewer, requestReview(one, "publication", "early-final"))
		return e
	}, "early final review")
	var aid, approvalHash string
	if err = pool.QueryRow(ctx, `select approval_id,payload_sha256 from catalog.release_review where tenant_id=$1 and draft_id=$2 and stage='legal'`, tenant, one.ResourceID).Scan(&aid, &approvalHash); err != nil {
		t.Fatal(err)
	}
	expectDenied(func() error {
		_, e := db.NewHumanApprovals(pool).Decide(ctx, reviewer, tenant, aid, "j3-store", approvalHash, true, "bypass", "catalog:review:legal", nil)
		return e
	}, "generic approval bypass")
	approve := func(d cr.Receipt) {
		t.Helper()
		for _, stage := range []string{"legal", "technical", "media", "publication"} {
			if _, err := store.Review(ctx, reviewer, requestReview(d, stage, d.CommandID+"-"+stage)); err != nil {
				t.Fatal(stage, err)
			}
		}
	}
	if public, err := db.NewElectromobility(pool).ListPublicModels(ctx, "catalog-j3"); err != nil || len(public) != 0 {
		t.Fatal("draft visible", public, err)
	}
	if _, err = pool.Exec(ctx, `update catalog.vehicle_model set display_name='Version Two' where tenant_id=$1 and model_id=$2`, tenant, model.ID); err != nil {
		t.Fatal(err)
	}
	captured, err := store.Draft(ctx, maker, one.ResourceID)
	if err != nil || !bytes.Contains(captured.Snapshot, []byte("Version One")) || bytes.Contains(captured.Snapshot, []byte("Version Two")) {
		t.Fatal("captured source mutated", err)
	}
	expectDenied(func() error { return sales.ActivatePriceBook(ctx, tenant, firstBook.ID) }, "generic price activation")
	expectDenied(func() error {
		_, e := pool.Exec(ctx, `update pricing.price_book_entry set amount_minor_units=1 where tenant_id=$1 and price_book_id=$2`, tenant, firstBook.ID)
		return e
	}, "captured price update")
	approve(one)
	pub := cr.PublishRequest{CommandID: "publish-one", DraftID: one.ResourceID, SnapshotSHA256: one.SnapshotSHA256, ExpectedGeneration: 0, Reason: "approved synthetic release"}
	expectDenied(func() error { _, e := store.Publish(ctx, maker, pub); return e }, "maker publish")
	var wg sync.WaitGroup
	var mu sync.Mutex
	newCount, replays := 0, 0
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r, e := store.Publish(ctx, reviewer, pub)
			mu.Lock()
			defer mu.Unlock()
			if e != nil {
				t.Error(e)
				return
			}
			if r.Replay {
				replays++
			} else {
				newCount++
			}
		}()
	}
	wg.Wait()
	if newCount != 1 || replays != 7 {
		t.Fatal("concurrency", newCount, replays)
	}
	checkPublic := func(name string, amount int64, generation int64, sha string) {
		t.Helper()
		models, e := db.NewElectromobility(pool).ListPublicModels(ctx, "catalog-j3")
		if e != nil || len(models) != 1 || models[0].DisplayName != name {
			t.Fatal("public owner", models, e)
		}
		priced, e := sales.PublicPrice(ctx, "catalog-j3", "AR", "j3-variant")
		if e != nil || priced.AmountMinorUnits != amount {
			t.Fatal("price projection", priced, e)
		}
		release, e := store.Public(ctx)
		if e != nil || release.Generation != generation || release.SHA256 != sha {
			t.Fatal("public release", release, e)
		}
		hits, e := store.PublicSearch(ctx, "Version")
		if e != nil || len(hits) != 1 || hits[0].Title != name || hits[0].Generation != generation || hits[0].SourceSHA256 != sha {
			t.Fatal("public search projection", hits, e)
		}
		var assetHash string
		if e = pool.QueryRow(ctx, `select sha256 from catalog.release_media where tenant_id=$1 and media_id=$2`, tenant, media.ResourceID).Scan(&assetHash); e != nil {
			t.Fatal(e)
		}
		asset, e := store.PublicMedia(ctx, assetHash)
		if e != nil || bytes.Contains(asset, []byte("script")) {
			t.Fatal("public normalized media", e)
		}
		var count int
		e = pool.QueryRow(ctx, `select count(*) from search.document where tenant_id=$1 and kind='catalog-model' and content_sha256=$2 and title=$3`, tenant, sha, name).Scan(&count)
		if e != nil || count != 1 {
			t.Fatal("search source mismatch", e, count)
		}
	}
	checkPublic("Version One", 120000, 1, one.SnapshotSHA256)
	secondBook := book(130000, nil)
	two := draft("draft-two", secondBook.ID)
	approve(two)
	next := cr.PublishRequest{CommandID: "publish-two", DraftID: two.ResourceID, SnapshotSHA256: two.SnapshotSHA256, ExpectedGeneration: 1, Reason: "approved second version"}
	_, err = pool.Exec(ctx, `create function catalog.reject_fixture_outbox() returns trigger language plpgsql as $$begin if new.event_type='catalog-release.publish' then raise exception 'fixture outbox unavailable';end if;return new;end$$;create trigger reject_catalog_outbox before insert on platform.outbox_event for each row execute function catalog.reject_fixture_outbox()`)
	if err != nil {
		t.Fatal(err)
	}
	expectDenied(func() error { _, e := store.Publish(ctx, reviewer, next); return e }, "publication outbox rollback")
	if _, err = pool.Exec(ctx, `drop trigger reject_catalog_outbox on platform.outbox_event;drop function catalog.reject_fixture_outbox()`); err != nil {
		t.Fatal(err)
	}
	if _, err = store.Publish(ctx, reviewer, next); err != nil {
		t.Fatal("second publish", err)
	}
	checkPublic("Version Two", 130000, 2, two.SnapshotSHA256)
	rollback := cr.PublishRequest{CommandID: "rollback-one", DraftID: one.ResourceID, SnapshotSHA256: one.SnapshotSHA256, ExpectedGeneration: 2, Reason: "explicit rollback to prior approved snapshot"}
	if _, err = store.Publish(ctx, reviewer, rollback); err != nil {
		t.Fatal("rollback", err)
	}
	checkPublic("Version One", 120000, 3, one.SnapshotSHA256)
	expectDenied(func() error {
		_, e := pool.Exec(ctx, `update search.document set title='bypass' where tenant_id=$1 and kind='catalog-model'`, tenant)
		return e
	}, "search bypass")
	// Publish a briefly valid version, leave it, then reject rollback after expiry.
	var now time.Time
	if err = pool.QueryRow(ctx, `select clock_timestamp()`).Scan(&now); err != nil {
		t.Fatal(err)
	}
	until := now.Add(2 * time.Second)
	shortBook := book(140000, &until)
	three := draft("draft-short", shortBook.ID)
	approve(three)
	short := cr.PublishRequest{CommandID: "publish-short", DraftID: three.ResourceID, SnapshotSHA256: three.SnapshotSHA256, ExpectedGeneration: 3, Reason: "synthetic expiry boundary"}
	if _, err = store.Publish(ctx, reviewer, short); err != nil {
		t.Fatal(err)
	}
	next.CommandID = "return-two"
	next.ExpectedGeneration = 4
	if _, err = store.Publish(ctx, reviewer, next); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `select pg_sleep(greatest(0,extract(epoch from $1::timestamptz-clock_timestamp()))+0.02)`, until); err != nil {
		t.Fatal(err)
	}
	short.CommandID = "expired-rollback"
	short.ExpectedGeneration = 5
	expectDenied(func() error { _, e := store.Publish(ctx, reviewer, short); return e }, "expired rollback")
	checkPublic("Version Two", 130000, 5, two.SnapshotSHA256)
	otherPool, err := pgxpool.New(ctx, rawURL)
	if err != nil {
		t.Fatal(err)
	}
	defer otherPool.Close()
	reopened, err := db.NewCatalogRelease(otherPool, profile, db.NewCommerce(otherPool))
	if err != nil {
		t.Fatal(err)
	}
	recovered, err := reopened.CommandReceipt(ctx, reviewer, "rollback-one")
	if err != nil || recovered.Generation != 3 || recovered.SnapshotSHA256 != one.SnapshotSHA256 {
		t.Fatal("durable recovery", recovered, err)
	}
	t.Log("CATALOG_RELEASE_CONNECTED_PASS three immutable snapshots/four shared reviews each; five publications including rollback; source edit isolated; storefront/price/search match; concurrency1new7replay;13table rollback; expired rollback denied; durable recovery")
}
