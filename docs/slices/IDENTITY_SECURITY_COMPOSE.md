# IDENTITY / SECURITY / AUTHZ compose slice

**Scope:** `LIBRARY_INFRASTRUCTURE / LOCAL_FIXTURES` at tip `227c6e0033eaab96cb42fa809ac6c37b758f84f4`.  
**Non-claim:** not production READY, not live IdP certification, not REVESTEX product admission.

Nightly marks **Identity/security** **PARCIAL** in the domain matrix while security lint/supply-chain paths are **HECHO** locally. This slice makes the compose path tangible without inventing packs. Enterprise audit obligations are composed via **`markdown_system/AUDIT_EVENT_PACK_PLAN.md`** (plan only — not a HECHO pack).

## Architecture contract (read-only cites @ `227c6e0`)

Do **not** edit these files from this slice — cite only:

| Authority | Path | Relevant claim |
|---|---|---|
| Nightly domain matrix | [`docs/ROADMAP.md`](../ROADMAP.md) | Identity/security · authz · audit = **PARCIAL**; selected OIDC + SAST + secure-ops pack filenames |
| Full domain map | [`docs/FRANCHISE_ARCHITECTURE.md`](../FRANCHISE_ARCHITECTURE.md) §1 *Identity / security / authorization / audit* | Selected **116** core paths; J5 PROVEN_LOCAL; live IdP/MFA = consumer gate; disk-only PCI/GDPR/OTP packs |

Selected **116** JSON `packId` → filename map (from `markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md` — no shorthand aliases):

| Selected `packId` | Concrete path |
|---|---|
| `GO-OIDC-SERVICE-TOKEN-BROKER` | `implementation_packs/GO_OIDC_SERVICE_TOKEN_BROKER.md` |
| `GO-OIDC-PORTAL-SESSION` | `implementation_packs/GO_OIDC_PORTAL_SESSION.md` |
| `TS-OIDC-PORTAL-ADAPTER` | `implementation_packs/TYPESCRIPT_OIDC_PORTAL_ADAPTER.md` |
| `GO-PG-CONTACT-CHANNEL-IDENTITY` | `implementation_packs/GO_PG_CONTACT_CHANNEL_IDENTITY.md` |
| `GO-HUMAN-APPROVAL-CORE` | `implementation_packs/GO_HUMAN_APPROVAL_CORE.md` |
| `SECURE-OPS-DELIVERY-CORE` | `implementation_packs/SECURE_OPERATIONS_DELIVERY_CORE.md` |
| `MICROSOFT-DEVSKIM-ADAPTED-SAST-GATE` | `implementation_packs/MICROSOFT_DEVSKIM_ADAPTED_SAST_GATE.md` |

Disk-only (catalog ≠ selected; gap gate if REQUIRED): `implementation_packs/GO_PCI_DSS_SCOPE_CORE.md`, `implementation_packs/GO_GDPR_CONSENT_ERASURE_CORE.md`, `implementation_packs/GO_OTP_VERIFICATION_CORE.md`.

**Adjacent (not ∈ selected 116 JSON):** compose bootstrap tooling only — `pack_id` `MARKDOWN-COMPOSITOR` in `implementation_packs/MARKDOWN_COMPOSITOR_CORE.md` (`markdown_system/PACK_PER_CLAIM_INDEX.md` row *composición Markdown*). Materialize before running `tools/compose-markdown-project.ps1`; not an identity/security selected member.

## 1. Door map (existing only)

| Layer | Status @ tip | Primary doors / packs | Evidence |
|---|---|---|---|
| **Router** | wired | `AGENTS.md`, `markdown_system/PACK_PER_CLAIM_INDEX.md` (claim → one pack), `markdown_system/CAPABILITY_CATALOG.md` (row *identidad y autorización*) | claim index row 22; catalog row 46 |
| **Authority manuals** | wired | `SECURITY_SRE_CLOUD_INFRASTRUCTURE.md` (identity/authz gates §21, suites §23), `SOFTWARE_BACKEND_API_ENGINEERING.md` (OIDC/API authz), `AI_SECURITY_GOVERNANCE_PRIVACY.md` (agent/tool authz) | search index |
| **48-surface contract** | wired | `markdown_system/TOTAL_SYSTEM_CAPABILITY_CONTRACT.md` (`IDENTITY`, `AUTHORIZATION`, `WEB-PORTALS` coupling) | surfaces 43, 49, 88 |
| **OIDC portal session** | `CONDITIONED` | `implementation_packs/GO_OIDC_PORTAL_SESSION.md`, `implementation_packs/TYPESCRIPT_OIDC_PORTAL_ADAPTER.md`, `docs/PORTAL_IDENTITY_LIFECYCLE.md` | `reconstruction_evidence/IDENTITY_PORTAL_RELEASE_V402.md` |
| **OIDC service token** | `REBUILD_VERIFIED / CONDITIONED` | `implementation_packs/GO_OIDC_SERVICE_TOKEN_BROKER.md`, `docs/oidc-service-token-runtime.md` | broker tests in `internal/platform/identity/` |
| **J5 / IdP administration contract** | `CONDITIONED` | `docs/J5_IDP_ACCESS_CONTRACT.md`, `internal/platform/postgres/identity_j5_integration_test.go` | `reconstruction_evidence/IDENTITY_J5_RELEASE_V402.md` |
| **Channel/contact identity** | `REBUILD_VERIFIED / CONDITIONED` | `implementation_packs/GO_PG_CONTACT_CHANNEL_IDENTITY.md`, `internal/contactidentity/` | `reconstruction_evidence/GO_PG_CONTACT_CHANNEL_IDENTITY_2026-09-04_V237.md` |
| **Human approval (selected 116)** | `CONDITIONED` | `implementation_packs/GO_HUMAN_APPROVAL_CORE.md` | `docs/FRANCHISE_ARCHITECTURE.md` §1 |
| **QR payload identity** | `PARCIAL` | `implementation_packs/GO_QR_CORE.md`, `internal/qr/` | `docs/ROADMAP.md` QR row; not in selected 116 JSON |
| **Security delivery / SAST** | `HECHO` (local) | `implementation_packs/SECURE_OPERATIONS_DELIVERY_CORE.md`, `implementation_packs/MICROSOFT_DEVSKIM_ADAPTED_SAST_GATE.md`, `markdown_system/ZERO_COST_VULNERABILITY_MONITORING_PROFILE.md` | `markdown_system/LIBRARY_HEALTH_CHECK.md` Security row |
| **Production identity admission schema** | wired (semantic, not live) | `production_admission_gate/validate_identity_authorization_admission.py`, `production_admission_gate/identity-authorization-admission.template.json` | `reconstruction_evidence/OPENID_IDENTITY_AUTHORIZATION_PRODUCTION_ADMISSION_2026-08-29_V110.md` |
| **Candidate architecture (not admitted)** | `CANDIDATE_PACK` | `architecture_packs/IDENTITY_AUTHORIZATION_KEYCLOAK_OPENFGA.md` | `PUBLIC_CODE_ARCHITECTURE_MASTER_MAP.md` |
| **Local OIDC fixture** | `LOCAL_FIXTURES` | `ci/local_identity_fixture.cjs` | `docs/LOCAL_REFERENCE_DELIVERY.md` |
| **Go principal / authz helpers** | wired | `internal/platform/identity/oidc.go`, `principal_test.go`, `service_token_broker_test.go` | unit tests |
| **Compose plans (selected doors)** | wired | `markdown_system/ENTERPRISE_BACKEND_PACK_PLAN.md`, `markdown_system/ENTERPRISE_WEB_PACK_PLAN.md`, `markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md` | JSON `packId` + `path` per row above |

### Audit surfaces (plan composes existing owners — no `*AUDIT*` HECHO pack)

Compose authority: **`markdown_system/AUDIT_EVENT_PACK_PLAN.md`** (`PLAN_ONLY`; **not** selected-116 `packId`). Canonical ledger: `audit.event` from `PG-TX-FOUNDATION` / `db/migrations/0001_platform_foundation.up.sql`.

| Concern | Owner path |
|---|---|
| Central audit DDL + append-only trigger | `implementation_packs/POSTGRES_TRANSACTIONAL_FOUNDATION.md` → `db/migrations/0001_platform_foundation.up.sql` |
| Transactional outbox (at-least-once, ≠ audit ledger) | `platform.outbox_event` in same migration; CDC optional via disk-only `implementation_packs/DEBEZIUM_POSTGRES_OUTBOX_RUNTIME.md` |
| Portal session revocation / retention counters | `implementation_packs/GO_OIDC_PORTAL_SESSION.md` + `docs/PORTAL_IDENTITY_LIFECYCLE.md` |
| Identity admission receipt binding | `production_admission_gate/validate_identity_authorization_admission.py` (`audit_event_observed`) |
| Security gate evidence | `implementation_packs/SECURE_OPERATIONS_DELIVERY_CORE.md` + `reconstruction_evidence/COMPOSITION_SECURITY_RELEASE_V402.md` |
| ASVS / security suites | `SECURITY_SRE_CLOUD_INFRASTRUCTURE.md` §23 (`security_authz`, `security_identity`) |
| Human approval audit trail | `implementation_packs/GO_HUMAN_APPROVAL_CORE.md` |

**Disk-only gap gate (catalog ≠ selected 116; cite only):** `implementation_packs/GO_PCI_DSS_SCOPE_CORE.md`, `implementation_packs/GO_GDPR_CONSENT_ERASURE_CORE.md`, `implementation_packs/GO_OTP_VERIFICATION_CORE.md` — open via `implementation_packs/CAPABILITY_GAP_RESOLUTION_GATE.md` when REQUIRED.

## 2. HECHO vs PARCIAL (honest @ `227c6e0`)

| Concern | Label | Why |
|---|---|---|
| Selected 116 pack filenames + code paths for OIDC portal, service token, contact identity, human approval | **HECHO** | `docs/ROADMAP.md` + `FRANCHISE_COMPLETE_PACK_PLAN.md` JSON rows |
| PKCE/state/nonce/JWE negatives, tenant/permission rejections, CAS session owner | **HECHO** (local) | `reconstruction_evidence/IDENTITY_PORTAL_RELEASE_V402.md` |
| Service-token authority validation, wildcard/extra-permission rejection | **HECHO** (local) | `internal/platform/identity/service_token_broker_test.go` |
| DevSkim-adapted SAST gate, secure-ops delivery core | **HECHO** (local lint path) | `docs/ROADMAP.md` nightly matrix |
| Identity admission **semantic** validator (5 observation kinds) | **HECHO** (schema) | `production_admission_gate/validate_identity_authorization_admission.py` + tests |
| Live IdP, HTTPS issuer, OIDC OP conformance execution, role journeys on target URL | **PARCIAL** | `markdown_system/CAPABILITY_CATALOG.md` row 46 pending block |
| J5 IdP account administration, MFA, break-glass on real tenant | **PARCIAL** | `docs/J5_IDP_ACCESS_CONTRACT.md` — fixture-only proof |
| Central enterprise audit schema (composed) | **PARCIAL** (plan) | `markdown_system/AUDIT_EVENT_PACK_PLAN.md` composes `audit.event` + owners; **FALTA** dedicated `*AUDIT*` HECHO pack in selected 116 |
| Keycloak/OpenFGA architecture pack | **PARCIAL** (`CANDIDATE_PACK`) | `architecture_packs/IDENTITY_AUTHORIZATION_KEYCLOAK_OPENFGA.md` — not `REUSABLE_PACK` |
| Composition-wide source/SCA/security release (T2803) | **PARCIAL** | `IDENTITY_PORTAL_RELEASE_V402` does not close T2803 |
| Production / cloud / offensive security | **CONDITIONED** | `qualification/FINAL_LIBRARY_READY_V402.json` `production_authorized: false` |

## 3. Compose recipe (tangible, no new packs)

Materialize **one claim at a time** via existing plans. Compositor entry (after adjacent `MARKDOWN-COMPOSITOR` / `implementation_packs/MARKDOWN_COMPOSITOR_CORE.md` is materialized — **not** a selected-116 member): `tools/compose-markdown-project.ps1`.

### Minimal identity/authz vertical (backend + web)

```powershell
# 1) Materialize compositor (adjacent pack_id MARKDOWN-COMPOSITOR; not ∈ selected 116)
pwsh -File ./materialize_markdown_pack.ps1 `
  -PackFile ./implementation_packs/MARKDOWN_COMPOSITOR_CORE.md `
  -Destination <compositor-dest>

# 2) Compose backend slice with OIDC owners
pwsh -File <compositor-dest>/tools/compose-markdown-project.ps1 `
  -LibraryRoot . `
  -PlanFile ./markdown_system/ENTERPRISE_BACKEND_PACK_PLAN.md `
  -Destination <backend-dest>

# 3) Compose web/BFF slice with OIDC portal adapter
pwsh -File <compositor-dest>/tools/compose-markdown-project.ps1 `
  -LibraryRoot . `
  -PlanFile ./markdown_system/ENTERPRISE_WEB_PACK_PLAN.md `
  -Destination <web-dest>
```

**Identity/security subset in selected 116 JSON** — use the `packId` → path table in §Architecture contract (from `markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md`).

Franchise-wide selection uses the same rows via `markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md` (116-pack profile). Do not open `implementation_packs/UNIFIED_REFERENCE_V403_R4.md` in context — materialize per `markdown_system/START_V403_LOCAL.md`.

### Runtime wiring after compose

1. Apply portal migration `db/migrations/0062_portal_session.up.sql` (and dependencies) per `docs/PORTAL_IDENTITY_LIFECYCLE.md`.
2. Bind profile `config/identity/portal-lifecycle.example.json` + SHA-256; enable `OIDC_PORTAL_LIFECYCLE_ENABLED=true` on BFF **and** Go together.
3. For service workloads, bind `config/identity/service-token.synthetic.json` per `docs/oidc-service-token-runtime.md`.
4. Local loopback only: start `node ci/local_identity_fixture.cjs --ttl-seconds 600` and set issuer from stdout; requires `OIDC_PORTAL_ALLOW_LOOPBACK_FIXTURE=true` on Go (BFF rejects loopback in production mode).

## 4. Local verify (`LOCAL_FIXTURES`, fail-closed)

### 4.1 Slice verifier (this PR)

```bash
python3 ci/verify_identity_security_fixture.py --workspace .
```

Checks:

- required door paths exist;
- `ci/fixtures/identity_authz_scenarios.json` fields **1:1** with verifier (see §5);
- tenant-scoped `audit_obligation_fields` match `audit.event` contract in `0001_platform_foundation.up.sql`;
- `audit_owners` doors exist (including `markdown_system/AUDIT_EVENT_PACK_PLAN.md`);
- optional `go test ./internal/platform/identity/...` when `go` is on PATH.

Exit non-zero on any missing door, extra fixture field, audit schema gap, or scenario gap — does **not** imply production admission.

Evidence index: `docs/slices/IDENTITY_SECURITY_AUDIT_FIXTURE_EVIDENCE.json`.

### 4.2 Go identity unit tests

```bash
go test ./internal/platform/identity/... -count=1
```

### 4.3 Production admission schema (synthetic profile only)

```bash
python3 -m unittest production_admission_gate.test_validate_identity_authorization_admission -v
```

Proves the **receipt adapter** rejects tampered OIDC/tenant/role evidence; still requires live execution to populate observations.

### 4.4 Library preflight (broader)

```powershell
pwsh -File ./markdown_system/VERIFY_LIBRARY_V403.ps1 -LibraryRoot .
```

### 4.5 Local reference harness (synthetic OIDC + PostgreSQL)

See `docs/LOCAL_REFERENCE_DELIVERY.md` and `markdown_system/START_V403_LOCAL.md` — qualification run, not a persistent IdP.

## 5. Fixture ↔ verifier field checklist (1:1)

`ci/fixtures/identity_authz_scenarios.json` must contain **exactly** these keys — no more, no less:

| Field | Verifier check |
|---|---|
| `schema` | must equal `elite.identity-authz-fixture/v1` |
| `scope` | must equal `LOCAL_FIXTURES` |
| `live_effects` | must be `false` |
| `tenants` | ≥2 unique strings |
| `unauthorized_scenarios` | exact set of 6 mandatory cases |
| `tenant_actions` | exact set `read`, `write`, `list` |
| `session_cases` | exact set of 5 revocation cases |
| `invariants.data_disclosed` | must be `false` |
| `invariants.mutation_observed` | must be `false` |
| `invariants.old_credential_accepted` | must be `false` |
| `audit_obligation_fields` | exact set: `tenant_id`, `event_id`, `actor_subject`, `action`, `resource_type`, `resource_id`, `decision` |
| `audit_owners` | non-empty list; each path must exist (includes `AUDIT_EVENT_PACK_PLAN.md`) |
| `audit_invariants.append_only` | must be `true` |
| `audit_invariants.tenant_scoped` | must be `true` |
| `audit_invariants.live_effects` | must be `false` |
| `doors` | non-empty list; each path must exist on disk |

Extra fixture keys (e.g. `description`, `roles`, `required_status_codes`) are rejected by `ci/verify_identity_security_fixture.py`.

## 6. Local checklist (library-side)

| # | Check | Pass criterion |
|---|---|---|
| C1 | Claim routed | `markdown_system/PACK_PER_CLAIM_INDEX.md` row *identidad y autorización* → `implementation_packs/GO_OIDC_SERVICE_TOKEN_BROKER.md` |
| C2 | Portal lifecycle docs | `docs/PORTAL_IDENTITY_LIFECYCLE.md` + `docs/J5_IDP_ACCESS_CONTRACT.md` present |
| C3 | Admission validator | `production_admission_gate/validate_identity_authorization_admission.py` importable |
| C4 | Fixture 1:1 | `ci/fixtures/identity_authz_scenarios.json` keys match §5 exactly |
| C5 | Loopback fixture bounded | `ci/local_identity_fixture.cjs` TTL 1–1200s; stdout issuer only |
| C6 | Wildcard rejection | `internal/platform/identity/service_token_broker_test.go` rejects `permissions: ["*"]` |
| C7 | Security lint path | `implementation_packs/MICROSOFT_DEVSKIM_ADAPTED_SAST_GATE.md` + `implementation_packs/SECURE_OPERATIONS_DELIVERY_CORE.md` exist |
| C8 | Architecture contract cited | `docs/ROADMAP.md` + `docs/FRANCHISE_ARCHITECTURE.md` §1 referenced; not edited |
| C9 | Audit plan wired | `markdown_system/AUDIT_EVENT_PACK_PLAN.md` present; `PACK_PER_CLAIM_INDEX.md` row *enterprise audit stream* |
| C10 | Audit fixture 1:1 | `audit_obligation_fields` + `audit_owners` + `audit_invariants` per §5 |
| C11 | No production claim | No `READY` / live IdP / conformance PASS in this slice's receipts |

## 7. CONDITIONED IdP gates (project, not library)

These remain **open until the target IdP and URLs are selected and evidenced**:

| Gate | Owner path | Blocker |
|---|---|---|
| G-IdP-01 | Target IdP operator | HTTPS issuer, client registration, exact redirect/logout URIs |
| G-IdP-02 | `implementation_packs/GO_OIDC_PORTAL_SESSION.md` | PKCE S256, offline_access, refresh rotation, RS256 JWKS |
| G-IdP-03 | `production_admission_gate/validate_identity_authorization_admission.py` | OpenID Foundation Conformance Suite 5.2.4 execution on target issuer |
| G-IdP-04 | `implementation_packs/MICROSOFT_PLAYWRIGHT_BROWSER_GATE.md` | Role positive + 6 unauthorized + tenant matrix on `application_url` |
| G-IdP-05 | `docs/J5_IDP_ACCESS_CONTRACT.md` | MFA, break-glass, access review, deprovisioning at IdP |
| G-IdP-06 | Edge/TLS | Loopback fixtures disabled; `APP_BASE_URL` / proxy alignment |
| G-IdP-07 | T2803 composition | Source/SCA/security release receipt for full franchise compose |

Official conformance reference (not executed by this slice): [OpenID Certification — Conformance Suite](https://openid.net/certification/about-conformance-suite/), [release v5.2.4](https://gitlab.com/openid/conformance-suite/-/releases/release-v5.2.4).

## 8. Public pattern research → elite **FALTA** rows only

Patterns cited for gap identification only — **no new packs invented**.

| Public authority | URL | Elite **FALTA** row (honest gap) |
|---|---|---|
| [NIST SP 800-63B (Digital Identity Guidelines — Authentication)](https://pages.nist.gov/800-63-3/sp800-63b.html) | Memorized secret strength, replay-resistant authentication, verifier compromise detection | **FALTA:** target IdP password/MFA policy receipt binding (library delegates to external IdP; fixtures do not prove AAL/IAL) |
| [NIST SP 800-63B — Session management](https://pages.nist.gov/800-63-3/sp800-63b.html#sec5) | Session binding, re-authentication, termination | **FALTA:** step-up / re-auth policy per journey on production `application_url` (portal session proves refresh/logout locally only) |
| [OIDC Core 1.0 — Authentication](https://openid.net/specs/openid-connect-core-1_0.html#Authentication) | `code` + PKCE, `state`, `nonce` | **FALTA:** live OP conformance run bound to project issuer (`OIDC_CONFORMANCE` observation kind) |
| [OIDC Core 1.0 — Token validation](https://openid.net/specs/openid-connect-core-1_0.html#TokenValidation) | Issuer, audience, signing key rotation | **FALTA:** signing-key rotation drill on target IdP with session eviction evidence (`SESSION_REVOCATION_ROTATION`) |
| [OIDC Core 1.0 — Standard Claims](https://openid.net/specs/openid-connect-core-1_0.html#StandardClaims) | `sub`, `aud`, `exp` | **HECHO** locally for synthetic tokens; **FALTA** cross-tenant matrix on deployed tenants |
| [OIDC Core 1.0 — Refresh Tokens](https://openid.net/specs/openid-connect-core-1_0.html#RefreshTokens) | Rotation, reuse detection | **PARCIAL:** CAS owner proven in portal pack; **FALTA** provider-side rotation policy attestation |
| [AWS SaaS Lens — Multi-tenant microservices](https://docs.aws.amazon.com/wellarchitected/latest/saas-lens/multi-tenant-microservices.html) | Tenant context in logs / observability | **FALTA:** centralized SIEM export with tenant partition on target cloud |

## 9. What this slice does not do

- Does not edit `docs/ROADMAP.md`, `docs/FRANCHISE_ARCHITECTURE.md`, or `docs/FRANCHISE_PLAYBOOK.md`.
- Does not invent SDR/GTM packs or workflows.
- Does not certify production or live IdP compatibility.
- Does not assert REVESTEX product admission.

## 10. Citations quick index

```text
docs/ROADMAP.md @ 227c6e0
docs/FRANCHISE_ARCHITECTURE.md §1 @ 227c6e0
markdown_system/AUDIT_EVENT_PACK_PLAN.md
markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md (selected packId JSON)
markdown_system/PACK_PER_CLAIM_INDEX.md (identidad + enterprise audit stream)
markdown_system/CAPABILITY_CATALOG.md:46
implementation_packs/POSTGRES_TRANSACTIONAL_FOUNDATION.md
implementation_packs/GO_OIDC_PORTAL_SESSION.md
implementation_packs/GO_OIDC_SERVICE_TOKEN_BROKER.md
implementation_packs/TYPESCRIPT_OIDC_PORTAL_ADAPTER.md
implementation_packs/GO_HUMAN_APPROVAL_CORE.md
implementation_packs/SECURE_OPERATIONS_DELIVERY_CORE.md
implementation_packs/MICROSOFT_DEVSKIM_ADAPTED_SAST_GATE.md
db/migrations/0001_platform_foundation.up.sql
production_admission_gate/validate_identity_authorization_admission.py
reconstruction_evidence/IDENTITY_PORTAL_RELEASE_V402.md
reconstruction_evidence/IDENTITY_J5_RELEASE_V402.md
docs/slices/IDENTITY_SECURITY_AUDIT_FIXTURE_EVIDENCE.json
ci/local_identity_fixture.cjs
ci/fixtures/identity_authz_scenarios.json
ci/verify_identity_security_fixture.py
```
