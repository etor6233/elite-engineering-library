# Microsoft BC warehouse cross-dock — reconstruction evidence V166

## Scope and provenance

`GO-SUPPLY-FACTORY-INVENTORY-API` 0.10.0 adds a narrow warehouse cross-dock lane adapted from Microsoft BCApps. The Go, SQL, HTTP and tests are local `ADAPTED`/`AUTHORED` implementation and are not Microsoft-authored code. This evidence does not claim complete Business Central WMS equivalence or project production readiness.

Authority is fixed to `microsoft/BCApps` commit `2eae56d704a1fd035d104f333602aea7091b7749`, tree `f9846fb1254c6311985c6131bb2af11c7f169e1b`, MIT, plus Microsoft Learn [Cross-dock items](https://learn.microsoft.com/en-gb/dynamics365/business-central/warehouse-how-to-cross-dock-items).

| Exact official file | Bytes | SHA-256 |
|---|---:|---|
| `src/Layers/W1/BaseApp/Warehouse/Activity/WhseCrossDockManagement.Codeunit.al` | 47,115 | `85a1eede4b2132f8ef22c853be1bd66a23a29b5176875a0432f500be150609fd` |
| `src/Layers/W1/BaseApp/Warehouse/Activity/WhseCrossDockOpportunity.Table.al` | 18,740 | `778cc8833ea7c28262a9ae695136281c66002d190845785e6f41e100db274ecb` |
| `src/Layers/W1/BaseApp/Warehouse/Document/WarehouseReceiptHeader.Table.al` | 24,887 | `d02b7f1fa0be0847e10e688d98125fdf34e1f157089a99d4632024d16fd2c0fe` |
| `src/Layers/W1/Tests/SCM-Warehouse/SCMWarehouseV.Codeunit.al` | 306,725 | `faee2bf1ef5bae41211f1374b8e1b7d49e9e7e75f2672e70ed0c5b181136f839` |
| `src/Layers/W1/Tests/SCM-Warehouse/SCMWarehouseOrders.Codeunit.al` | 238,014 | `79f2e97c278763aa7c140dcc6dffaa53f7e6e6fe2d55100b7f42bad2ccd0a1da` |
| `src/Layers/W1/Tests/SCM-Warehouse/SCMWarehouseUnitTests.Codeunit.al` | 170,444 | `9af357c3eddaa6f6acaa50623aec43753d212efe7d412e6a948f35c87521ac3f` |

## Executable contract admitted

- explicit fixed, unblocked PICK/PUTPICK cross-dock bin and organization/item due-date policy;
- released same-organization/item bulk-transfer demand inside the horizon; specific-cost lines and late demand are excluded;
- remaining demand subtracts shipped quantity, existing cross-dock allocations and exact generic picks still open or registered-but-unshipped;
- receipt placement is `min(receipt, demand, capacity)` and ordinary put-away receives the exact remainder;
- a second receipt cannot duplicate an existing opportunity; an allocation is not pickable before its put-away line is registered;
- generic picks exclude cross-dock bins; the matching transfer uses its allocation first and cannot exceed outstanding demand;
- reservation identity includes source line, activity and sequence, preserving the active-demand unique index while allowing multiple partial picks;
- cancellation releases balance, allocation and pick link; registration converts reserved allocation into picked evidence and moves exact stock to SHIP;
- serializable transactions plus an organization/item advisory lock close concurrent double allocation/reservation.

Sales/service/production demand, base-calendar working days, breakbulk/UOM conversion, cubage/weight, warehouse class, handheld/barcode execution, external WMS reconciliation and production policies remain conditioned.

## Reconstruction and executable gates

- pack SHA-256 after promotion metadata: `155b9ab75d9affe0f355f49703957d5917feb6b1157315cfdf4268ca18261cb4`; 427,616 bytes;
- supply materialized 58/58 files; a full author-tree versus Markdown round trip reported zero path/hash differences;
- backend composed 350 files from 29 packs from Markdown only;
- official Go 1.26.7 archive baseline: 74,955,002 bytes, SHA-256 `f4f534a486e4bc3387fa18f08208f2f854b7aaea8a08f2a2d829a914a05abb11`;
- official PostgreSQL 18.6 archive: 343,808,005 bytes, SHA-256 `fbe23da234ee31547bf8a36d29dfd81e82b849df2d2b78d2eecb43d360252f8c`; isolated loopback cluster, locale C and data checksums;
- migrations 0001–0032 applied cleanly; 0032 down/up passed after cleanup and the focal flow passed again;
- the focal corpus passed repeatedly and proves invalid bin/policy rejection, generic partial pick 2, receipt split 3+4, non-duplication on a second receipt, unavailable-before-put-away, over-demand rejection, cancellation/retry, picked 3, reserved 0, cross-dock balance 0 and exact ship balance 3;
- cross-dock plus the existing bulk-transfer shipment/transit/receipt/cost flow passed together twice;
- canonical Markdown composition passed full `go test ./...` against PostgreSQL, `go vet ./...` and `go build ./cmd/...`.

## Failure memory and final closure

V166 recorded every acquisition, path, fixture, model and assertion failure through `LIB-FAIL-1590`; the current ledger contains 1,590 local failures plus 202 upstream conditions, 1,792 unique IDs, zero local `OPEN` rows and zero duplicates. After evidence and governing-count promotion, `VERIFY_LIBRARY_PASS` reported 94 packs, 1,067 materialized files, 522 Markdown files and backend 350/web 84. `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit` re-ran the root verifier, 54 readiness regressions, all 121 governed sources/16 source profiles and every available offline executable lane with official Go 1.26.7 supplied explicitly.
