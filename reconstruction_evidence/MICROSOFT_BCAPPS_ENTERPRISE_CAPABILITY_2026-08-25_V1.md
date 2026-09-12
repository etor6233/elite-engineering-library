# Microsoft BCApps Enterprise Capability Evidence — 2026-08-25 V1

## Claim

This evidence proves exact acquisition and read-only inspection of the official Microsoft Business Central application source selected by Elite. It establishes that real Microsoft product code exists for finance, procurement, inventory, manufacturing, sales, warehouse and service. It does not prove a stable release build, a licensed Business Central runtime, Argentine localization or REVESTEX-specific behavior.

## Source identity

```text
owner/repository: microsoft/BCApps
commit: 31a860b527f0dc72c7a44a255d7e7d403cfa4789
archive URL: https://github.com/microsoft/BCApps/archive/31a860b527f0dc72c7a44a255d7e7d403cfa4789.zip
archive bytes: 225939109
archive SHA-256: e3151b040df39cada83d41fcf5d9e6cdff1d8fddf226934c1ad21a2bdebc7210
root LICENSE SHA-256: c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383
source license: MIT, with third-party terms retained by path
acquisition receipt: microsoft-bcapps-snapshot.json
receipt result: PASS
```

The archive is a commit snapshot from `main`, not a stable integrated product release. The newer `releases/28.4/StrictMode` tag at `cb07eef12935e07dd4258c122d0e7b1920fc1223` was rechecked through the official repository API: `src/Layers/W1/BaseApp` and `src/Layers/W1/Tests` both returned absent, although `src/Apps/W1` exists. It was not substituted as if it contained the integral business source.

## Official upstream statement inspected

The repository README identifies BCApps as the home of Business Central application development and says Microsoft engineers, partners and contributors use the same code, pipelines and tooling. It enumerates System Application, Business Foundation, Base Application, first-party apps, Test Framework, Performance Toolkit and AI Test Toolkit. This supports the product-source claim only; Elite does not extend the statement to runtime licensing or target readiness.

## Static inventory reproduced

```text
all AL under src: 36608
W1 BaseApp AL: 8142
Finance: 1037
Inventory: 950
Service: 637
Sales: 591
Manufacturing: 529
Purchases: 370
CRM: 369
Warehouse: 355
Projects: 338
Integration: 332
Bank: 239
FixedAssets: 213
Foundation: 209
HumanResources: 115
Assembly: 113
OtherCapabilities: 99
CostAccounting: 82
Pricing: 78
eServices: 72
Utilities: 65
System: 563
CashFlow: 49
RoleCenters: 34
Entitlements: 25
Invoicing: 21
FinancialMgt: 1
Permissions: 235
Modules: 287
Removed: 2
```

The count was produced from the extracted archive with recursive file enumeration of `*.al`. Counts are evidence of inspected source surface, not a functionality or quality score.

Relevant Microsoft first-party W1 apps were also enumerated: EDocument 431 AL, EDocumentConnectors 249, PEPPOL 53, ExpenseAgent 421, PayablesAgent 66, BankAccRecWithAI 18, Shopify 589, Subscription Billing 433, Quality Management 253, Subcontracting 226, FieldServiceIntegration 83, APIV1/APIV2 411 combined, external-storage connectors/attachments 76 combined and PowerBIReports 304. These are real product modules inside the same exact archive; they remain subject to Business Central runtime and per-service account/terms.

## Test surface reproduced

The W1 test layer contains concrete suites, including:

```text
ERM 145; ERM-Finance 62; ERM-Purchase 32; ERM-Sales 25
Bank 19; Cost Accounting 12; Physical Inventory 12
Prepayment 10; Reverse 18
SCM 77; SCM-Assembly 37; SCM-Costing 43
SCM-Manufacturing 36; SCM-Planning 26; SCM-Reservation 22
SCM-Service 57; SCM-Warehouse 39; SCM-Workflow 3
```

## Build condition observed

The official `LOCAL_DEV_ENV.md` requires Docker Desktop in Windows-container mode and `BcContainerHelper`; it creates a local Business Central container through `build/scripts/DevEnv/NewDevEnv.ps1`. The audit host did not provide Docker, a Business Central runtime/container or a project license/tenant. Consequently:

```text
archive identity: PASS
license identity: PASS
source inventory: PASS
official test inventory: PASS
local Business Central build: NOT_RUN_MISSING_PLATFORM
official AL suites: NOT_RUN_MISSING_PLATFORM
REVESTEX journeys: NOT_RUN_MISSING_PROJECT_INPUTS
admission: CONDITIONAL_PLATFORM
```

## Capability conclusion

The inspected source is the strongest complete public enterprise-product candidate currently admitted from a highest-value technology company. It materially covers generic finance, purchasing, inventory/costing, manufacturing, sales, warehouse/logistics and service management inside Business Central.

It does not supply, ready and proven for REVESTEX, Argentine customs/fiscal adapters, marketplace/payment accounts, carrier integrations, franchise territory/royalty rules, the private document corpus, exact product schema, or a provider-independent Go/PostgreSQL implementation. Those claims remain blocked or must be implemented as clearly labeled project code after explicit platform authorization.
