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
