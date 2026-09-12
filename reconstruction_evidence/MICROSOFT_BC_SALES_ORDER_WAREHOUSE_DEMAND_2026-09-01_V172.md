# Microsoft BC sales-order to warehouse demand — reconstruction evidence V172

## Result and claim boundary

`GO-SUPPLY-FACTORY-INVENTORY-API` 0.15.0 materializes 83 files and the enterprise backend composes 29 packs / 375 implementation files. V172 adds four `ADAPTED` Go/SQL/test blocks and one `AUTHORED` provenance record, and updates ten existing domain, adapter, HTTP and regression-test owners. None of the portable Go or SQL is represented as Microsoft-authored code.

The admitted claim is narrow: an operational customer-order line can govern packaging-aware warehouse picking only after an exact immutable variant→item/sales-UOM binding exists. Draft, missing, mismatched, over-picked or divergent requests fail closed. Shipment, invoicing, service-part demand, production-component demand and target production readiness are not claimed.

## Fixed official authority

- Repository: `microsoft/BCApps`
- Signed commit: `2eae56d704a1fd035d104f333602aea7091b7749`
- Tree: `f9846fb1254c6311985c6131bb2af11c7f169e1b`
- Root source license: MIT
- Snapshot ZIP SHA-256: `f7e984f2a1e9784a351068f42f0a0cefdc8317704b309f2a94fce5f0f91bc5b6`

| Exact Microsoft owner/test | Bytes | SHA-256 |
|---|---:|---|
| `Warehouse/Request/WhseCreateSourceDocument.Codeunit.al` | 7,814 | `df2a5f4767b58d4d307a223dcfe1c3d0176c9d336d6f9ebdd25ec414f40672f8` |
| `Sales/SalesWhsePostShipment.Codeunit.al` | 46,608 | `84121bd38555f3a7cba59dd97f985ae41af2da349359a60217bcab14ee7c3747` |
| `Sales/Document/SalesLine.Table.al` | 708,062 | `ca65615dc05bc555a878ba2635c955897930332ef171e313a00103074fbc7278` |
| `Tests/SCM-Warehouse/SCMWarehousePick.Codeunit.al` | 202,517 | `f856e941d86a78a09d233aba6dfe5470dacbf08c853c26bc5e26512c3d6881c5` |
| `Tests/SCM-Warehouse/SalesOrderWhseValidateLine.Codeunit.al` | 15,092 | `12b4a3ce4fdef9592aad3054f5f174c2c4493d8377216cd4bea6890289d66deb` |

Microsoft initializes warehouse outstanding quantities from quantity/base quantity, checks the source document line and bin, resolves sales document/line identity and validates variant/UOM coupling. Its fixed tests cover released sales orders, partial and multiple picks and registration. Elite transfers only those demonstrated invariants into the portable boundary.

## Executable behavior

- `sales_warehouse_binding` fixes one sellable variant to one inventory item and sales UOM. The UOM FK uses the immutable item-UOM owner; retrying the same request is idempotent and a changed replay fails.
- `customer-order` is no longer a decorative demand label. The adapter requires the exact tenant, organization, order, line, variant binding, item and an order state in `placed`, `confirmed`, `paid` or `allocated`.
- Sold quantity is converted through the fixed sales UOM factor into base quantity. Every open or registered pick for the exact order line reduces the outstanding quantity. Cancelled picks do not.
- Pick request evidence fixes organization, bin, item, demand identity, base quantity, requested/resolved UOM and behavior flags. Exact replay returns the existing open/registered/cancelled activity; divergent replay fails.
- Request identity and item work are serialized with PostgreSQL transaction advisory locks. Two requests racing for the same remaining quantity produce exactly one winner.
- Registration remains warehouse handling evidence at SHIP; V172 does not falsely mark the order shipped, delivered, invoiced or fiscally issued.

## Executable proof

Pinned runtimes:

- Go `1.26.7` Windows/amd64, `go.exe` SHA-256 `5463fe58fa999d74420f00ee1b36d31c3da90a57ff204159f159859567ad61fc`;
- PostgreSQL `18.6` admitted archive SHA-256 `fbe23da234ee31547bf8a36d29dfd81e82b849df2d2b78d2eecb43d360252f8c`;
- `pgx` `5.10.0` from the composed module lock.

The following passed on 2026-09-01:

1. official exact-source acquisition in memory with byte and SHA-256 identity;
2. isolated Go formatting and focal package compilation;
3. migrations 0001–0037 on a new PostgreSQL cluster/database;
4. focal binding, draft rejection, placed order, 3 BOX→36 EA base, exact replay, divergent replay, cancellation, registration and one-winner concurrency test;
5. migration 0037 down/up and focal repetition;
6. all uncached Go/PostgreSQL tests on a separate clean database, then `go vet ./...` and `go build ./...`;
7. updater regression suite and atomic synchronization of fifteen materialization blocks;
8. Markdown-only supply reconstruction 83/83 and exact SHA equality for all fifteen changed/new owners;
9. Markdown-only enterprise composition 29 packs / 375 files with exact owner equality;
10. canonical fresh-database migrations 0001–0037, full uncached suite, vet and build;
11. canonical 0037 down/up and focal repetition.
12. root `VERIFY_LIBRARY_PASS` with 94 packs, 1,092 materialized files, 529 Markdown files and backend 375; executable Audit exit 0 with 94 packs / 121 upstream sources and all governed component regressions.

The final focal starts with one order line of 3 BOX. A draft-order pick is refused. After placement, a 1 BOX pick can be replayed exactly, a changed replay is refused, and cancellation releases its demand. A registered 1 BOX pick consumes 12 EA base. Two independent requests then race for the remaining 2 BOX/24 EA and exactly one succeeds. The authoritative handled total is 36 EA; every subsequent pick is refused.

## Failure learning

V172 records local failures 1690 through 1704. Truncated discovery, commit/tree identity confusion, PowerShell interpolation, assumed paths, broad helper output, incorrect patch anchoring, absent PATH toolchain, retired helper alias, multi-command fixture SQL, new cleanup dependencies, synthetic demand labels, contaminated shared-database output and non-literal documentation hunks were each invalidated before canonical promotion. No failed output is counted as PASS. Library memory at this point is 1,704 local failures plus 202 upstream conditions: 1,906 unique IDs.

## Artifact identity

- Pack: `implementation_packs/GO_SUPPLY_FACTORY_INVENTORY_API.md`
- Version: `0.15.0`
- Bytes before root-gate evidence/readiness linkage: `617882`
- SHA-256 before root-gate evidence/readiness linkage: `9a041b78ef815cc7c347ee2e27ea7af3423345e8c44081b31f8ac5370e204ebe`
- Materialized files: `83`
- Enterprise backend composition: `29 packs / 375 files`

Production remains conditioned on shipment/fiscal completion, service/production demand owners, project infrastructure/providers, load/security/recovery evidence, live monitoring and business acceptance.
