# CATALOG / STOCK / PI-PL — franchise compose slice

Tip pin: `42d9f83ee6fa1536bcb59d66e6c7850b01ddd9b0`

**Scope:** tangible local slice for catalog master data → serialized stock → pricing (PI) and procurement document classes proforma-invoice / packing-list (PL). Library infrastructure / local fixtures only — not production, not REVESTEX product.

**PI-PL in this slice:** `pricing.price_book` / `pricing.price_book_entry` (variant-bound list prices) plus document classes `proforma-invoice` and `packing-list` from the typed document bridge. Not profit-and-loss accounting.

---

## 1. Door order (compose prerequisites)

| # | Door | Path |
|---|------|------|
| 1 | Library vs product | `README.md` |
| 2 | Agent router | `AGENTS.md` |
| 3 | Use / search / resume | `markdown_system/USE_LIBRARY_V403.md` |
| 4 | Health / GAPs | `markdown_system/LIBRARY_HEALTH_CHECK.md` |
| 5 | Franchise consumer | `START_FRANCHISE.md` |
| 6 | Library vs product gate | `reconstruction_evidence/LIBRARY_VS_PRODUCT_GATE_V402.md` |
| 7 | V402 receipt (evidence only) | `qualification/FINAL_LIBRARY_READY_V402.json` |

Full selected composition owner: `markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md` (116 packs). Extension: `markdown_system/FRANCHISE_COMPLETE_EXTENSION_PACK_PLAN_V403.md` (+4 → 120).

---

## 2. Elite packs / modules for this slice

### Schema foundation (catalog + stock + pricing tables)

| packId | version | path |
|--------|---------|------|
| `PG-TX-FOUNDATION` | 0.1.0 | `implementation_packs/POSTGRES_TRANSACTIONAL_FOUNDATION.md` |
| `ELECTROMOBILITY-FRANCHISE-MODULES` | 0.1.2 | `implementation_packs/ELECTROMOBILITY_FRANCHISE_MODULES.md` |

Creates `catalog.*`, `inventory.stock_unit`, `pricing.*`, FK `inventory.stock_unit.variant_id → catalog.vehicle_variant` and `pricing.price_book_entry.variant_id → catalog.vehicle_variant` (`db/migrations/0003_electromobility_franchise_modules.up.sql`).

### Stock / supply API owners

| packId | version | path |
|--------|---------|------|
| `GO-ENTERPRISE-BACKEND` | 0.4.8 | `implementation_packs/GO_ENTERPRISE_BACKEND_CORE.md` |
| `GO-SUPPLY-FACTORY-INVENTORY-API` | 0.18.0 | `implementation_packs/GO_SUPPLY_FACTORY_INVENTORY_API.md` |
| `GO-CONNECTED-SERIAL-SUPPLY` | 0.2.0 | `implementation_packs/GO_CONNECTED_SERIAL_SUPPLY.md` |

### Pricing / commerce (PI)

| packId | version | path |
|--------|---------|------|
| `GO-BC-EXACT-AMOUNT-ADAPTER` | 0.1.0 | `implementation_packs/GO_BC_EXACT_AMOUNT_ADAPTER.md` |
| `GO-BC-SALES-CONTRACT-ADAPTER` | 0.1.0 | `implementation_packs/GO_BC_SALES_CONTRACT_ADAPTER.md` |
| `GO-COMMERCE-PRICING-PAYMENT-API` | 0.8.0 | `implementation_packs/GO_COMMERCE_PRICING_PAYMENT_API.md` |

### Catalog authoring / publication

| packId | version | path |
|--------|---------|------|
| `GO-CONNECTED-CATALOG-AUTHORING` | 0.1.0 | `implementation_packs/GO_CONNECTED_CATALOG_AUTHORING.md` |
| `TS-CATALOG-AUTHORING-PORTAL` | 0.1.2 | `implementation_packs/TYPESCRIPT_CATALOG_AUTHORING_PORTAL.md` |
| `GO-CONNECTED-CATALOG-PUBLICATION` | 0.2.0 | `implementation_packs/GO_CONNECTED_CATALOG_PUBLICATION.md` |
| `TYPESCRIPT-PUBLISHED-CATALOG-STOREFRONT` | 0.1.0 | `implementation_packs/TYPESCRIPT_PUBLISHED_CATALOG_STOREFRONT.md` |

### PI-PL document classes (PL = packing-list; not admitted end-to-end procurement pack)

| artifact | path | note |
|----------|------|------|
| Typed fixture registry | `internal/documentbridge/typed-fixtures.json` | `proforma-invoice`, `packing-list` classes |
| Document intelligence profile | `markdown_system/OFFICIAL_DOCUMENT_INTELLIGENCE_PROFILE.md` | PO/PI/PL/BOL relations |
| Fixture catalog | `markdown_system/OFFICIAL_DOCUMENT_FIXTURE_CATALOG.md` | official sample inventory |
| PDF fixture | `config/documents/typed-fixtures/product-specification-catalog.pdf` | catalog-adjacent typed doc |

### Optional external catalog read adapters (not required for local slice)

| packId | path |
|--------|------|
| `PYTHON-AMAZON-SPAPI-CATALOG-ADAPTER` | `implementation_packs/PYTHON_AMAZON_SPAPI_CATALOG_ADAPTER.md` |
| plan | `markdown_system/AMAZON_SPAPI_CATALOG_PACK_PLAN.md` |

---

## 3. Fixture paths

| Purpose | Path |
|---------|------|
| Franchise module vertical fixture (catalog→factory→stock→order) | `db/tests/0003_electromobility_franchise_modules.test.sql` |
| Catalog→stock FK consistency (this slice) | `db/tests/0088_catalog_stock_consistency.test.sql` |
| Catalog authoring profile | `config/catalog.authoring.fixture.json` |
| Commerce price→order→stock→payment flow | `internal/platform/postgres/commerce_integration_test.go` (`TestCommercePriceOrderAllocationPaymentFlow`) |
| Catalog–stock join / QR resolution | `internal/platform/postgres/capture_integration_test.go`, `internal/platform/postgres/capture.go` |
| Order–stock connected evidence | `reconstruction_evidence/ORDER_STOCK_CONNECTED_V296.md` |
| Catalog connected release | `reconstruction_evidence/CATALOG_CONNECTED_RELEASE_V402.md` |
| Proforma / packing-list typed fixtures | `internal/documentbridge/typed-fixtures.json` |
| Supply role fixture | `config/supply.role.fixture.json` |

---

## 4. Compose order (minimal tangibility)

1. Materialize foundation: `PG-TX-FOUNDATION` → `ELECTROMOBILITY-FRANCHISE-MODULES` (migrations through `0003`).
2. Apply enterprise + inventory packs: `GO-ENTERPRISE-BACKEND`, `GO-SUPPLY-FACTORY-INVENTORY-API` (migrations `0025`+ per pack manifest).
3. Wire commerce pricing: `GO-BC-EXACT-AMOUNT-ADAPTER`, `GO-BC-SALES-CONTRACT-ADAPTER`, `GO-COMMERCE-PRICING-PAYMENT-API`.
4. Optional connected surfaces: catalog authoring/publication packs + `config/catalog.authoring.fixture.json`.
5. PI-PL document lane: load `internal/documentbridge/typed-fixtures.json` only; do not claim PI/PL→catalog SKU reconciliation without an admitted pack.

Full 116-pack materialization: `markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md` + `materialize_markdown_pack.ps1` / compositor per `START_FRANCHISE.md`.

---

## 5. Verify locally

**Preflight (library):**

```powershell
./VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Preflight
./markdown_system/VERIFY_LIBRARY_V403.ps1
```

**Slice-specific (PostgreSQL 18.x disposable loopback DB with all migrations applied):**

```powershell
./markdown_system/verify_catalog_stock_slice.ps1
```

Or manually:

```bash
# SQL fixture (rollback-safe; requires migrated schema)
psql "$TEST_DATABASE_URL" -v ON_ERROR_STOP=1 -f db/tests/0088_catalog_stock_consistency.test.sql

# Go integration (catalog FK + join + pricing FK)
go test ./internal/platform/postgres -run '^TestCatalogStockConsistency(ForeignKeys|Join|PriceBook)$' -count=1 -v

# End-to-end commerce catalog→price→stock (existing owner)
go test ./internal/platform/postgres -run '^TestCommercePriceOrderAllocationPaymentFlow$' -count=1 -v
```

Set `TEST_DATABASE_URL` to an owned loopback database (e.g. `elite_confirmation_*` on `127.0.0.1`). Never use production targets.

---

## 6. Expected receipts

| Check | PASS signal | Evidence path |
|-------|-------------|---------------|
| Schema FK catalog→stock | SQL test completes without exception | `db/tests/0088_catalog_stock_consistency.test.sql` |
| Go FK + join consistency | `TestCatalogStockConsistency*` PASS | `internal/platform/postgres/catalog_stock_consistency_integration_test.go` |
| Price book bound to variant | `TestCatalogStockConsistencyPriceBook` PASS | same |
| Commerce allocation | `TestCommercePriceOrderAllocationPaymentFlow` PASS | `reconstruction_evidence/ORDER_STOCK_CONNECTED_V296.md` |
| Library preflight | `VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Preflight` exit 0 | stdout receipt |
| Franchise module fixture | `0003` test PASS | `db/tests/0003_electromobility_franchise_modules.test.sql` |

Local PASS ≠ cloud, live provider, device, visual approval, or production admission.

---

## 7. FALTA — gaps (cite only; no invented HECHO packs)

| ID | Gap | Public pattern (cite only) | Proposed elite addition |
|----|-----|---------------------------|-------------------------|
| FALTA-CS-01 | No admitted pack links `proforma-invoice` / `packing-list` extracted lines to `catalog.vehicle_variant` or `inventory.stock_unit` | Oracle Fusion [Inventory Management](https://docs.oracle.com/en/cloud/saas/supply-chain-and-manufacturing/25a/faimi/) — item master vs on-hand | `CAPABILITY_GAP_RESOLUTION` row + pack plan under `implementation_packs/CAPABILITY_GAP_RESOLUTION_GATE.md` |
| FALTA-CS-02 | No franchise-selected pack for Amazon Catalog **write** / listing mutation | [Amazon SP-API Catalog Items API](https://developer-docs.amazon.com/sp-api/docs/catalog-items-api-v2022-04-01-reference) — read-oriented adapter exists (`PYTHON-AMAZON-SPAPI-CATALOG-ADAPTER`) | Extend `markdown_system/AMAZON_SPAPI_CATALOG_PACK_PLAN.md`; no new packId until admitted |
| FALTA-CS-03 | PI-PL document corpus lacks production-grade SKU↔variant reconciliation tests | Google Document AI packing-list sample (see `markdown_system/OFFICIAL_DOCUMENT_FIXTURE_CATALOG.md`) | Project-specific corpus + `PROJECT_DOCUMENT_INTELLIGENCE_DECISION_TEMPLATE.md` |
| FALTA-CS-04 | Bulk/lot `inventory.sales_warehouse_binding` variant→item path not covered by this SQL slice | Microsoft BC [Item Tracking](https://learn.microsoft.com/en-us/dynamics365/business-central/design-details-item-tracking) | Covered by `GO-SUPPLY-FACTORY-INVENTORY-API` separate tests (`docs/inventory/MICROSOFT_BC_DERIVATION.md`) |

---

## 8. Invariants (from elite sources)

- `inventory.stock_unit.variant_id` must reference an existing `catalog.vehicle_variant` (`db/migrations/0003_electromobility_franchise_modules.up.sql`).
- QR product capture resolves catalog variant, never silently aliases stock-unit ID (`internal/platform/postgres/capture.go`).
- Catalog publication does not confirm stock or commercial authorization (`src/platform/i18n/private-messages.json` / connected catalog owners).
- Price books bind amounts per `variant_id`; commerce derives line price server-side (`implementation_packs/GO_COMMERCE_PRICING_PAYMENT_API.md`).
