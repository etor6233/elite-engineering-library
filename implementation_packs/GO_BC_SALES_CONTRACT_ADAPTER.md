# Go BC Sales Contract Adapter

## 1. Metadata

```yaml
pack_id: "GO-BC-SALES-CONTRACT-ADAPTER"
pack_version: "0.1.0"
status:
  authority: ELITE_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: REUSABLE_PACK
claim: "Narrow fixed-source active price/time/currency eligibility and quote header/line copy with document retargeting, preserving the existing portable representation and caller policies."
stacks: ["Go 1.26.8", "PostgreSQL 18.6"]
compatible_with: ["GO-BC-EXACT-AMOUNT-ADAPTER 0.1.0", "GO-COMMERCE-PRICING-PAYMENT-API 0.6.3", "GO-FRANCHISE-CUSTOMER-JOURNEY-API 0.10.20"]
incompatible_with: ["missing Microsoft-BCApps MIT notice dependency", "blank currency fallback or automatic FX", "fractional quantity/UOM conversion", "calendar Date equivalence without the documented timestamp mapping", "full BC quote/order or scheduling equivalence"]
license_expression: "MIT AND LicenseRef-Workspace-Owner"
upstream_sources: ["https://github.com/microsoft/BCApps/tree/2eae56d704a1fd035d104f333602aea7091b7749"]
verified_at: "2026-09-11"
```

This is an unpublished staging candidate until its owners/profile are integrated.
Only the two explicitly mapped source files are ADAPTED; four files are local
clock/test/document glue. Existing callers retain AUTHORED provenance. Microsoft
did not author or approve this Go/PostgreSQL code. The license file is already
owned by required dependency GO-BC-EXACT-AMOUNT-ADAPTER0.1.0; no duplicate path.

## 2. Applicability

Use with the existing explicit-variant sales price table, nonblank currency,
PostgreSQL microsecond half-open validity and existing one-line quote contract.
BC ending Date inclusive is mapped from exclusive timestamp end minus1microsecond;
the exact order currency excludes BC's blank fallback by schema. Quote type,
source-line filter and field copy are derived; caller consent, expiry, retention,
states, taxes, discounts, number generation and reservation policy are unchanged.

## 3. Architecture contract

All three price readers compose the same fixed SQL predicate with one database
clock context. This avoids a second price owner and fetch/filter races. The quote
owner passes an already scoped header/line, then persists transferred fields in
its existing transaction/outbox. No schema, module dependency, runtime, secret,
provider, background worker or public API is added. Costs are constant per price
row and linear in quote lines; current integration supplies one line.

## 4. Exact file manifest

```text
CREATE internal/platform/postgres/bc_price_sql.go
CREATE internal/platform/postgres/bc_price_clock.go
CREATE internal/bcsales/quote.go
CREATE internal/bcsales/quote_test.go
CREATE internal/platform/postgres/bc_price_integration_test.go
CREATE docs/provenance/BC_SALES_DERIVATION.md
```

## 5. Materialization blocks

### FILE: `internal/platform/postgres/bc_price_sql.go`
```yaml
block_id: "GO-BC-SALES-CONTRACT-ADAPTER:file1:v1"
operation: CREATE
provenance: ADAPTED
source: "Microsoft BCApps 2eae56d704a1fd035d104f333602aea7091b7749; exact source/line/SHA and representation mapping in docs/provenance/BC_SALES_DERIVATION.md"
license: "MIT"
sha256: "5fa8cc2ddda831196909627ebda7be9b463b152fcb98c5a41b335e3060d4b3ad"
variables: []
secrets_allowed: false
```
````go
package postgres

// SPDX-License-Identifier: MIT
// ADAPTED from Microsoft BCApps PriceCalculationBufferMgt.SetFiltersOnPriceListLine
// at 2eae56d704a1fd035d104f333602aea7091b7749, lines 258-269.
// Mapping: the local price table fixes sales/unit-price and exact variant/UOM;
// tenant, market, book ID and exact currency remain caller bindings. SQL NULL
// ending maps to BC's unbounded 0D. PostgreSQL timestamps have microsecond
// precision: local exclusive ending -> inclusive ending minus one microsecond.
// No fallback currency, UOM conversion, tax, discount or best-price policy.
const bcPriceEligibilitySQL = `b.status='active'
 and b.valid_from<=price_context.at
 and (b.valid_until is null or b.valid_until-interval '1 microsecond'>=price_context.at)`

// VerifySelectedLine:280-285 accepts the requested currency or BC's blank
// fallback. The existing local currency constraint excludes blank values, so
// that branch is unreachable. $3 is the caller's locked order currency.
const bcPriceCurrencyBindingSQL = `b.currency=$3`
````

### FILE: `internal/platform/postgres/bc_price_clock.go`
```yaml
block_id: "GO-BC-SALES-CONTRACT-ADAPTER:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local clock, verification or documentation glue; no vendor source attribution"
license: "LicenseRef-Workspace-Owner"
sha256: "7809684cfa3b6eda465b20f3eee196757384d82f3459c17f561354f3097a19ea"
variables: []
secrets_allowed: false
```
````go
package postgres

// AUTHORED query composition: one database-owned observation point for the
// source-derived eligibility predicate; no local/client clock or timezone rule.
const bcPriceTimeContextSQL = `cross join (select clock_timestamp() as at) as price_context`
````

### FILE: `internal/bcsales/quote.go`
```yaml
block_id: "GO-BC-SALES-CONTRACT-ADAPTER:file3:v1"
operation: CREATE
provenance: ADAPTED
source: "Microsoft BCApps 2eae56d704a1fd035d104f333602aea7091b7749; exact source/line/SHA and representation mapping in docs/provenance/BC_SALES_DERIVATION.md"
license: "MIT"
sha256: "dfaf75c731fe8b6f7acdda8b9f3c6f18a4ce55862448caaa72cf902bcf921790"
variables: []
secrets_allowed: false
```
````go
// SPDX-License-Identifier: MIT
// ADAPTED from Microsoft BCApps SalesQuotetoOrder.Codeunit.al at
// 2eae56d704a1fd035d104f333602aea7091b7749. See docs/provenance/BC_SALES_DERIVATION.md.
package bcsales

import "errors"

type DocumentType string

const (
	Quote DocumentType = "Quote"
	Order DocumentType = "Order"
)

var ErrQuoteType = errors.New("source document must be a quote")

// Header/Line represent only fields already present in the portable owner's
// quote/order contract. IDs are supplied by the caller's existing ID generator.
// No posting, VAT, prepayment, reservation, order state or consent is selected.
type Header struct {
	DocumentType DocumentType
	ID           string
	QuoteID      string
	CustomerID   string
	Currency     string
}

type Line struct {
	DocumentType   DocumentType
	DocumentID     string
	VariantID      string
	Quantity       int64
	UnitPriceMinor int64
}

func TransferQuoteToOrder(source Header, lines []Line, orderID string) (Header, []Line, error) {
	// OnRun:37 TestField(Document Type, Quote).
	if source.DocumentType != Quote {
		return Header{}, nil, ErrQuoteType
	}
	// CreateSalesHeader:126-132 copies the header, retargets type/number and
	// preserves the source quote number. Number-series allocation stays caller glue.
	order := source
	order.DocumentType = Order
	order.ID = orderID
	order.QuoteID = source.ID
	result := make([]Line, 0, len(lines))
	for _, sourceLine := range lines {
		// TransferQuoteToOrderLines:331-333 scopes the source records; :340-342
		// copies a matching line and retargets only its document type/number.
		if sourceLine.DocumentType != source.DocumentType || sourceLine.DocumentID != source.ID {
			continue
		}
		line := sourceLine
		line.DocumentType = Order
		line.DocumentID = order.ID
		result = append(result, line)
	}
	return order, result, nil
}
````

### FILE: `internal/bcsales/quote_test.go`
```yaml
block_id: "GO-BC-SALES-CONTRACT-ADAPTER:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local clock, verification or documentation glue; no vendor source attribution"
license: "LicenseRef-Workspace-Owner"
sha256: "b91da0bf7565cf9f9b4fda47515fbf7fd0b7dc8611a7169de8fc247b2f37d711"
variables: []
secrets_allowed: false
```
````go
package bcsales

// AUTHORED source-conservation fixtures, not execution of Microsoft AL suites.
import (
	"errors"
	"reflect"
	"testing"
)

func TestQuoteTransferPreservesFieldsAndFiltersForeignDocuments(t *testing.T) {
	source := Header{DocumentType: Quote, ID: "quote-a", CustomerID: "customer-a", Currency: "ARS"}
	lines := []Line{
		{DocumentType: Quote, DocumentID: "quote-b", VariantID: "foreign", Quantity: 1, UnitPriceMinor: 9},
		{DocumentType: Order, DocumentID: "quote-a", VariantID: "wrong-type", Quantity: 1, UnitPriceMinor: 9},
		{DocumentType: Quote, DocumentID: "quote-a", VariantID: "one", Quantity: 1, UnitPriceMinor: 1001},
		{DocumentType: Quote, DocumentID: "quote-a", VariantID: "two", Quantity: 3, UnitPriceMinor: 2003},
	}
	before := append([]Line(nil), lines...)
	order, transferred, err := TransferQuoteToOrder(source, lines, "order-a")
	if err != nil || order != (Header{DocumentType: Order, ID: "order-a", QuoteID: "quote-a", CustomerID: "customer-a", Currency: "ARS"}) {
		t.Fatalf("header copy/retarget failed: %+v %v", order, err)
	}
	want := []Line{{DocumentType: Order, DocumentID: "order-a", VariantID: "one", Quantity: 1, UnitPriceMinor: 1001}, {DocumentType: Order, DocumentID: "order-a", VariantID: "two", Quantity: 3, UnitPriceMinor: 2003}}
	if !reflect.DeepEqual(transferred, want) || !reflect.DeepEqual(lines, before) {
		t.Fatalf("line conservation/scope failed: %+v", transferred)
	}
	transferred[0].UnitPriceMinor = 999
	if !reflect.DeepEqual(lines, before) {
		t.Fatal("output aliases the caller's immutable quote lines")
	}
	if _, _, err = TransferQuoteToOrder(Header{DocumentType: Order}, lines, "new"); !errors.Is(err, ErrQuoteType) {
		t.Fatal("non-quote source accepted")
	}
	_, empty, err := TransferQuoteToOrder(source, lines[:2], "order-empty")
	if err != nil || len(empty) != 0 {
		t.Fatal("foreign lines copied or a new empty-order policy invented")
	}
}

func FuzzQuoteTransferConservation(f *testing.F) {
	f.Add("quote-a", "order-a", "variant-a", int64(1), int64(1000), true)
	f.Add("q", "o", "v", int64(0), int64(-1), false)
	f.Add("q", "o", "v", int64(9223372036854775807), int64(9223372036854775807), true)
	f.Fuzz(func(t *testing.T, quoteID, orderID, variant string, quantity, price int64, belongs bool) {
		if len(quoteID)+len(orderID)+len(variant) > 1024 {
			t.Skip()
		}
		lineID := quoteID
		if !belongs {
			lineID += "different"
		}
		source := Header{DocumentType: Quote, ID: quoteID, Currency: "ARS", CustomerID: "customer"}
		line := Line{DocumentType: Quote, DocumentID: lineID, VariantID: variant, Quantity: quantity, UnitPriceMinor: price}
		order, lines, err := TransferQuoteToOrder(source, []Line{line}, orderID)
		if err != nil || order.ID != orderID || order.QuoteID != quoteID || order.Currency != source.Currency || order.CustomerID != source.CustomerID || order.DocumentType != Order {
			t.Fatal("header conservation failed")
		}
		if !belongs {
			if len(lines) != 0 {
				t.Fatal("foreign quote copied")
			}
			return
		}
		if len(lines) != 1 || lines[0].DocumentID != orderID || lines[0].DocumentType != Order || lines[0].Quantity != quantity || lines[0].UnitPriceMinor != price || lines[0].VariantID != variant {
			t.Fatal("line field changed beyond type/document")
		}
	})
}
````

### FILE: `internal/platform/postgres/bc_price_integration_test.go`
```yaml
block_id: "GO-BC-SALES-CONTRACT-ADAPTER:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local clock, verification or documentation glue; no vendor source attribution"
license: "LicenseRef-Workspace-Owner"
sha256: "f4a8aa5c2b21d3b2d06201460668a08e644b320f09ffd779d3d18c47e42d5e12"
variables: []
secrets_allowed: false
```
````go
package postgres

// AUTHORED differential fixtures: source inclusive-ending predicate versus the
// established local half-open timestamp contract. This does not run AL suites.
import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestBCPriceEligibilityInclusiveSourceHalfOpenStorage(t *testing.T) {
	raw := os.Getenv("TEST_DATABASE_URL")
	if raw == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	cfg, err := pgxpool.ParseConfig(raw)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires isolated loopback confirmation database")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	at := time.Date(2026, 9, 11, 12, 0, 0, 123456000, time.UTC)
	count := 0
	for _, status := range []string{"draft", "active", "retired"} {
		for _, fromDelta := range []time.Duration{-24 * time.Hour, -time.Microsecond, 0, time.Microsecond, 24 * time.Hour} {
			from := at.Add(fromDelta)
			for _, untilDelta := range []time.Duration{-24 * time.Hour, -time.Microsecond, 0, time.Microsecond, 24 * time.Hour, 48 * time.Hour} {
				var until *time.Time
				if untilDelta != 48*time.Hour {
					v := at.Add(untilDelta)
					until = &v
				}
				for _, currency := range []string{"ARS", "USD"} {
					// $3 matches the exact caller binding used by AddOrderLine.
					query := `select ` + bcPriceEligibilitySQL + ` and ` + bcPriceCurrencyBindingSQL + `
					from (select $1::text status,$2::text currency,$4::timestamptz valid_from,$5::timestamptz valid_until) b
					cross join (select $6::timestamptz at) price_context`
					var got bool
					if err := pool.QueryRow(ctx, query, status, currency, "ARS", from, until, at).Scan(&got); err != nil {
						t.Fatal(err)
					}
					// Independent original local contract; includes exact ending,
					// next microsecond and T176's starting-after-observation refusal.
					want := status == "active" && !from.After(at) && (until == nil || at.Before(*until)) && currency == "ARS"
					if got != want {
						t.Fatalf("eligibility drift: status=%s from=%s until=%v currency=%s got=%t want=%t", status, from, until, currency, got, want)
					}
					count++
				}
			}
		}
	}
	if count != 180 {
		t.Fatalf("fixture matrix incomplete: %d", count)
	}
}
````

### FILE: `docs/provenance/BC_SALES_DERIVATION.md`
```yaml
block_id: "GO-BC-SALES-CONTRACT-ADAPTER:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local clock, verification or documentation glue; no vendor source attribution"
license: "LicenseRef-Workspace-Owner"
sha256: "87ae95d84c8c334d802dd77c734adfc816290d95b53206ce087858863cb3ae40"
variables: []
secrets_allowed: false
```
````markdown
# BC sales predicate and quote-copy derivation

Two files are ADAPTED from Microsoft BCApps fixed commit
`2eae56d704a1fd035d104f333602aea7091b7749`, tree
`f9846fb1254c6311985c6131bb2af11c7f169e1b`. Local Go/PostgreSQL code is not written,
endorsed or executed by Microsoft. Existing caller files remain AUTHORED.

| Source | SHA-256 | Exact source lines | Destination |
|---|---|---|---|
| src/Layers/W1/BaseApp/Pricing/Calculation/PriceCalculationBufferMgt.Codeunit.al | 6bd6c18dbdea9d0af11f132c4c5e50b975c8c2053838c0fe43fc668f1998297b | SetFiltersOnPriceListLine258-269: active260, ending263, starting268; VerifySelectedLine280-285: currency282 | internal/platform/postgres/bc_price_sql.go |
| src/Layers/W1/BaseApp/Sales/Document/SalesQuotetoOrder.Codeunit.al | 5ba30dbcb3d5bb9fb250df6691ad7fc461e1e023665d1572652b3561a736270d | OnRun37 quote type; CreateSalesHeader126-132 clone/type/quote reference; TransferQuoteToOrderLines332-333 filter and340-342 clone/type/document reference | internal/bcsales/quote.go |
| src/Layers/W1/Tests/ERM/TestPriceCalculationV16.Codeunit.al | 4b199f235bb527c8853157187bfa5fa35c72a7816ced0833f0ff3bc043c211bb | T1762661 rejects start after order date; currency refusal is the VerifySelectedLine source branch | TestBCPriceEligibilityInclusiveSourceHalfOpenStorage |

The root Microsoft MIT notice, SHA
`c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383`, is already owned
verbatim at `licenses/Microsoft-BCApps-MIT.txt` by required dependency
GO-BC-EXACT-AMOUNT-ADAPTER0.1.0. This pack must be composed with that dependency;
it does not duplicate ownership or remove the notice. Each adapted file carries
SPDX MIT and the pinned source pointer. Tests, clock/DTO mapping and this document
remain local LicenseRef-Workspace-Owner glue.

## Representation and semantic limits

BC pricing uses calendar Date with inclusive ending; the local owner already uses
PostgreSQL microsecond timestamps and exclusive valid_until. The adapter maps a
finite local end to inclusive end minus one microsecond and NULL to BC's unbounded
0D before applying the source >= comparison. The existing starting endpoint stays
inclusive. A single database clock reading per query is AUTHORED composition,
separate from the adapted predicate; it is not a new regional calendar policy.

The existing local table represents sales/unit prices of explicit variants, with
nonblank uppercase currency and no selectable UOM. BC's blank-currency fallback
is therefore unreachable; exact currency equality is used for the locked order
currency. No automatic FX, discount, UOM conversion, best-price ranking or price
book overlap policy is introduced. Tenant, market, book and variant are caller
bindings. PublicPrice, AddOrderLine and CreateQuote all call the same predicate.

Quote transfer copies only represented header/line fields: customer, currency,
variant, quantity and unit price; it retargets document type/ID and keeps source
quote ID. Source line scoping excludes a different quote/document type. The
existing caller maps its one-line quote to quantity1 and provides generated IDs;
it keeps its existing placed/accepted states, customer authorization, quote
expiry, acceptance evidence, transaction, outbox and quote retention. None of
those policies is attributed to BC. No VAT, prepayment, posting, dimensions,
deferral, assembly, number series or reservations are newly implemented. This is
not a complete Sales-Quote-to-Order engine or a Business Central runtime.

## Verification scope

The PostgreSQL fixture compares 180 combinations of status, start, end and
currency against the previous half-open local contract, including exact
boundaries and wrong currency. Existing Commerce price/order/stock/request and
Journey persistence/isolation/replay suites passed on54 unchanged migrations,
with zero skips. Quote unit fixtures verify scope, header/line conservation,
wrong source type and independent output storage. The native Go fuzz gate used
3seeds and70,545executions with a3s budget; no failing input. Vet and complete
module build passed. No AL suite, provider account, deploy or live acceptance was
executed. The gate is for these two source transformations only.

No schema, API, dependency pin or durable format changes. Rollback restores the
two previous callers and removes this adapter after no consumer references it.
Existing BC amount/overflow fixes must remain in Commerce during integration.
````

## 6. Configuration surface

No secrets or implicit business configuration. Source timestamp, source records,
order ID and existing price scope are caller inputs. Required license dependency:
licenses/Microsoft-BCApps-MIT.txt SHA c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383.

## 7. Dependency bill

BCApps MIT fixed-source adaptation only; Go1.26.8 and PostgreSQL18.6 retain the
existing composition's locks. No AL runtime, extra Go module or SDK is introduced.

## 8. Apply order

Compose after GO-BC-EXACT-AMOUNT-ADAPTER0.1.0 and with both exact candidate callers.
Never materialize a second owner of the MIT notice. No migration. Rollback the
two callers before removing this adapter; preserve their existing BC amounts fix.

## 9. Verification

180 predicate differential cases; native quote-copy invariants and finite3s fuzz;
fresh54-migration PostgreSQL Commerce and Journey suites, vet/build. G0-G8 source,
license, deltas and receipts are in sales-derivation-evidence.md in the staging
review bundle. This is not AL-suite execution, live payment, production or whole
Commerce/Journey provenance closure.

## 10. Reconstruction evidence

Materialize all three candidate packs and compare exact output hashes to the
tested reference. Preserve the source-map, dependency-owned MIT notice and
before/after caller hashes when integrating into the canonical profile.
