# Return Exchange Fulfillment Execution — V143

Date: 2026-08-31  
Result: `REBUILD_VERIFIED / CONDITIONED`

## Scope and truthful provenance

`GO-RETURN-EXCHANGE-FULFILLMENT-WORKER` 0.1.0 adds eight `AUTHORED` files. No local line is represented as copied Microsoft, PostgreSQL or AWS source. Microsoft Dynamics 365 governs the use of an independent replacement sales order linked to the RMA/original sale and the separation of return, replacement, inventory and financial processing. PostgreSQL governs the serializable transaction, row locks and `SKIP LOCKED`; AWS governs durable retry/idempotency expectations.

The automatic lane is intentionally exact: one delivered order, one quantity-one line, original allocated stock, identical replacement variant at the same organization, original total equal to that single line, succeeded physical inventory effect and a separate accounting request. Different variants, multi-line/split orders, price difference, discounts/taxes/shipping allocation, cross-organization sourcing and remote carrier delivery remain blocked until a project supplies approved policy and an executable owner.

## Reconstructed artifact

- pack round trip: 8/8 paths and SHA-256 byte-identical;
- library after this evidence: 86 packs / 929 materialization blocks / 485 Markdown / 36 profiles;
- provenance: 787 `AUTHORED`, 37 `ADAPTED`, 105 `VERBATIM`;
- backend profile: 27 packs / 277 files composed from Markdown without collision;
- materialization record: 277 source/output hashes equal and recomputed against all outputs;
- schema: PostgreSQL migrations 0001–0021 on PostgreSQL 18.6;
- root runtime: full `go test ./...`, `go vet ./...` and six actual main binaries build from Markdown.

## Official architecture evidence

Current Microsoft documentation establishes that:

1. an RMA/return order formally authorizes and tracks return processing;
2. replacement can occur after receipt/disposition or, under explicit policy, before receipt;
3. the replacement is an independent sales order associated with the return/RMA rather than a mutation of the original sale;
4. return lines and new sales lines require explicit fulfillment attributes and financial settlement is the net of both sides;
5. exact cost reversal links back to the original sales entry and remains an accounting concern.

Sources: [Sales returns](https://learn.microsoft.com/en-us/dynamics365/supply-chain/sales-marketing/sales-returns), [Create an item replacement order](https://learn.microsoft.com/en-us/dynamics365/supply-chain/sales-marketing/create-item-replacement-order), [Configure and process an exchange](https://learn.microsoft.com/en-us/dynamics365/commerce/orderexchanges), and [Process sales returns or cancellations](https://learn.microsoft.com/en-us/dynamics365/business-central/sales-how-process-sales-returns-cancellations).

## Real PostgreSQL evidence

The clean-database integration proves:

1. the worker claims only the fulfillment-owned exchange effect with lease and fencing token;
2. the original sale, exact line/stock, disposition, succeeded inventory effect and accounting request are proven before mutation;
3. one available identical-variant stock unit is locked with `FOR UPDATE SKIP LOCKED` and moves atomically from `available` version 7 to `reserved` version 8;
4. a separate zero-balance replacement order/line, prepared handover and immutable exchange link are written in the same transaction;
5. three outbox events, one immutable successful attempt and terminal execution are committed together;
6. a second exchange cannot allocate the same stock and becomes delayed `REPLACEMENT_STOCK_UNAVAILABLE` retry;
7. the prepared exchange record rejects mutation and terminal work cannot be reclaimed;
8. migration 0021 down removes its objects and up/SQL gates pass again.

`prepared` means only that local stock and handover are prepared. It never means delivered, accepted by the customer, posted to the ledger, fiscally accepted or settled by a provider.

## Failures converted into memory

`LIB-FAIL-1418` through `LIB-FAIL-1424` record assumed filenames/layouts, an incomplete resolved-exception fixture, compositor entrypoint drift, process-ending script chaining, a historical record-row selector and the temporary confusion of metadata URLs with source-lock artifacts. PostgreSQL rejected the invalid fixture before any exchange side effect; the corrected suite recreated the database and passed. The compositor self-test, profile composition and hash audit now run as separate observable stages, and Audit governs the 121-source lock count.

## Conditions not converted into PASS

No live customer, carrier, accounting ledger, ARCA credit note, payment/credit application, cross-franchise transfer, CDN/WAF, IdP, provider account, load, offensive security, backup/PITR, deployment, rollback or business acceptance was executed against a selected target. The pack is immediately materializable code for the admitted exchange lane, not proof that a concrete franchise system is production-ready.
