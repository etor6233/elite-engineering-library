# IDENTITY / SECURITY / AUTHZ compose slice

**Scope:** `LIBRARY_INFRASTRUCTURE / LOCAL_FIXTURES` at tip `477358d0d7090f2c39467abe5634d81f38afae7f`.  
**Non-claim:** not production READY, not live IdP certification, not REVESTEX product admission.

Nightly/domain health marks **identity/authz** as **PARCIAL** (fixture-only, `CONDITIONED` admission, live IdP gates open) while **security lint/supply-chain** paths are **HECHO** locally. This slice makes the compose path tangible without inventing packs.

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
| **QR payload identity** | `PARCIAL` | `implementation_packs/GO_QR_CORE.md`, `internal/qr/` | `markdown_system/LIBRARY_HEALTH_CHECK.md` domain matrix |
| **Security delivery / SAST** | `HECHO` (local) | `implementation_packs/SECURE_OPERATIONS_DELIVERY_CORE.md`, `implementation_packs/MICROSOFT_DEVSKIM_ADAPTED_SAST_GATE.md`, `markdown_system/ZERO_COST_VULNERABILITY_MONITORING_PROFILE.md` | `markdown_system/LIBRARY_HEALTH_CHECK.md` Security row |
| **Production identity admission schema** | wired (semantic, not live) | `production_admission_gate/validate_identity_authorization_admission.py`, `production_admission_gate/identity-authorization-admission.template.json` | `reconstruction_evidence/OPENID_IDENTITY_AUTHORIZATION_PRODUCTION_ADMISSION_2026-08-29_V110.md` |
| **Candidate architecture (not admitted)** | `CANDIDATE_PACK` | `architecture_packs/IDENTITY_AUTHORIZATION_KEYCLOAK_OPENFGA.md` | `PUBLIC_CODE_ARCHITECTURE_MASTER_MAP.md` |
| **Local OIDC fixture** | `LOCAL_FIXTURES` | `ci/local_identity_fixture.cjs` | `docs/LOCAL_REFERENCE_DELIVERY.md` |
| **Go principal / authz helpers** | wired | `internal/platform/identity/oidc.go`, `principal_test.go`, `service_token_broker_test.go` | unit tests |
| **Compose plans (selected doors)** | wired | `markdown_system/ENTERPRISE_BACKEND_PACK_PLAN.md`, `markdown_system/ENTERPRISE_WEB_PACK_PLAN.md`, `markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md` | JSON `packId` entries for OIDC packs |

### Audit surfaces (no dedicated `*AUDIT*` pack)

Audit obligations are distributed across admitted owners:

| Concern | Owner path |
|---|---|
| Portal session revocation / retention counters | `GO_OIDC_PORTAL_SESSION` + `docs/PORTAL_IDENTITY_LIFECYCLE.md` |
| Identity admission receipt binding | `production_admission_gate/validate_identity_authorization_admission.py` |
| Security gate evidence | `SECURE_OPERATIONS_DELIVERY_CORE` + `reconstruction_evidence/COMPOSITION_SECURITY_RELEASE_V402.md` |
| ASVS / security suites | `SECURITY_SRE_CLOUD_INFRASTRUCTURE.md` §23 (`security_authz`, `security_identity`) |

## 2. HECHO vs PARCIAL (honest @ tip)

| Concern | Label | Why |
|---|---|---|
| Pack filenames + code paths for OIDC portal, service token, contact identity | **HECHO** | admitted packs + findable Go/TS owners |
| PKCE/state/nonce/JWE negatives, tenant/permission rejections, CAS session owner | **HECHO** (local) | `IDENTITY_PORTAL_RELEASE_V402` connected proof |
| Service-token authority validation, wildcard/extra-permission rejection | **HECHO** (local) | `service_token_broker_test.go` |
| DevSkim-adapted SAST gate, secure-ops delivery core, vulnerability monitoring profile | **HECHO** (local lint path) | `LIBRARY_HEALTH_CHECK` Security row |
| Identity admission **semantic** validator (5 observation kinds) | **HECHO** (schema) | `validate_identity_authorization_admission.py` + tests |
| Live IdP, HTTPS issuer, OIDC OP conformance execution, role journeys on target URL | **PARCIAL** | `CAPABILITY_CATALOG` row 46 pending block |
| J5 IdP account administration, MFA, break-glass on real tenant | **PARCIAL** | `J5_IDP_ACCESS_CONTRACT.md` — fixture-only proof |
| Keycloak/OpenFGA architecture pack | **PARCIAL** (`CANDIDATE_PACK`) | not `REUSABLE_PACK`; deep G3–G7 open |
| Composition-wide source/SCA/security release (T2803) | **PARCIAL** | `IDENTITY_PORTAL_RELEASE_V402` explicitly does not close T2803 |
| Production / cloud / offensive security | **CONDITIONED** | all domains; no `.github/workflows` @ tip |

## 3. Compose recipe (tangible, no new packs)

Materialize **one claim at a time** via existing plans. Compositor entry (after `MARKDOWN_COMPOSITOR_CORE` is materialized): `tools/compose-markdown-project.ps1`.

### Minimal identity/authz vertical (backend + web)

```powershell
# 1) Materialize compositor (if not already present in destination tree)
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

**Packs touched by those plans (identity/security subset only):**

| packId | File |
|---|---|
| `GO-OIDC-PORTAL-SESSION` | `implementation_packs/GO_OIDC_PORTAL_SESSION.md` |
| `GO-OIDC-SERVICE-TOKEN-BROKER` | `implementation_packs/GO_OIDC_SERVICE_TOKEN_BROKER.md` |
| `GO-PG-CONTACT-CHANNEL-IDENTITY` | `implementation_packs/GO_PG_CONTACT_CHANNEL_IDENTITY.md` |
| `TS-OIDC-PORTAL-ADAPTER` | `implementation_packs/TYPESCRIPT_OIDC_PORTAL_ADAPTER.md` |
| `SECURE-OPS-DELIVERY-CORE` | `implementation_packs/SECURE_OPERATIONS_DELIVERY_CORE.md` |
| `MICROSOFT-DEVSKIM-ADAPTED-SAST-GATE` | `implementation_packs/MICROSOFT_DEVSKIM_ADAPTED_SAST_GATE.md` |

Franchise-wide selection adds the same OIDC packIds via `markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md` (116-pack profile). Do not open `implementation_packs/UNIFIED_REFERENCE_V403_R4.md` in context — materialize per `markdown_system/START_V403_LOCAL.md`.

### Runtime wiring after compose

1. Apply portal migration `db/migrations/0062_portal_session.up.sql` (and dependencies) per `PORTAL_IDENTITY_LIFECYCLE.md`.
2. Bind profile `config/identity/portal-lifecycle.example.json` + SHA-256; enable `OIDC_PORTAL_LIFECYCLE_ENABLED=true` on BFF **and** Go together.
3. For service workloads, bind `config/identity/service-token.synthetic.json` per `docs/oidc-service-token-runtime.md`.
4. Local loopback only: start `node ci/local_identity_fixture.cjs --ttl-seconds 600` and set issuer from stdout; requires `OIDC_PORTAL_ALLOW_LOOPBACK_FIXTURE=true` on Go (BFF rejects loopback in production mode).

## 4. Local verify (`LOCAL_FIXTURES`, fail-closed)

### 4.1 Slice verifier (this PR)

```bash
python ci/verify_identity_security_fixture.py --workspace .
```

Checks:

- required door paths exist;
- `ci/fixtures/identity_authz_scenarios.json` schema + mandatory negative matrix;
- optional `go test ./internal/platform/identity/...` when `go` is on PATH.

Exit non-zero on any missing door or scenario gap — does **not** imply production admission.

### 4.2 Go identity unit tests

```bash
go test ./internal/platform/identity/... -count=1
```

### 4.3 Production admission schema (synthetic profile only)

```bash
python -m unittest production_admission_gate.test_validate_identity_authorization_admission -v
```

Proves the **receipt adapter** rejects tampered OIDC/tenant/role evidence; still requires live execution to populate observations.

### 4.4 Library preflight (broader)

```powershell
pwsh -File ./markdown_system/VERIFY_LIBRARY_V403.ps1 -LibraryRoot .
```

### 4.5 Local reference harness (synthetic OIDC + PostgreSQL)

See `docs/LOCAL_REFERENCE_DELIVERY.md` and `markdown_system/START_V403_LOCAL.md` — qualification run, not a persistent IdP.

## 5. Local checklist (library-side)

| # | Check | Pass criterion |
|---|---|---|
| C1 | Claim routed | `PACK_PER_CLAIM_INDEX.md` row *identidad y autorización* → `GO_OIDC_SERVICE_TOKEN_BROKER.md` |
| C2 | Portal lifecycle docs | `docs/PORTAL_IDENTITY_LIFECYCLE.md` + `docs/J5_IDP_ACCESS_CONTRACT.md` present |
| C3 | Admission validator | `production_admission_gate/validate_identity_authorization_admission.py` importable |
| C4 | Fixture scenarios | `ci/fixtures/identity_authz_scenarios.json` covers 6 unauthorized + cross-tenant matrix |
| C5 | Loopback fixture bounded | `ci/local_identity_fixture.cjs` TTL 1–1200s; stdout issuer only |
| C6 | Wildcard rejection | `service_token_broker_test.go` rejects `permissions: ["*"]` |
| C7 | Security lint path | `MICROSOFT_DEVSKIM_ADAPTED_SAST_GATE.md` + `SECURE_OPERATIONS_DELIVERY_CORE.md` exist |
| C8 | No production claim | No `READY` / live IdP / conformance PASS in this slice's receipts |

## 6. CONDITIONED IdP gates (project, not library)

These remain **open until the target IdP and URLs are selected and evidenced**:

| Gate | Owner | Blocker |
|---|---|---|
| G-IdP-01 | Target IdP operator | HTTPS issuer, client registration, exact redirect/logout URIs |
| G-IdP-02 | `GO_OIDC_PORTAL_SESSION` | PKCE S256, offline_access, refresh rotation, RS256 JWKS |
| G-IdP-03 | `production_admission_gate` | OpenID Foundation Conformance Suite 5.2.4 execution on target issuer |
| G-IdP-04 | Playwright gate | Role positive + 6 unauthorized + tenant matrix on `application_url` |
| G-IdP-05 | `J5_IDP_ACCESS_CONTRACT` | MFA, break-glass, access review, deprovisioning at IdP |
| G-IdP-06 | Edge/TLS | Loopback fixtures disabled; `APP_BASE_URL` / proxy alignment |
| G-IdP-07 | T2803 composition | Source/SCA/security release receipt for full franchise compose |

Official conformance reference (not executed by this slice): [OpenID Certification — Conformance Suite](https://openid.net/certification/about-conformance-suite/), [release v5.2.4](https://gitlab.com/openid/conformance-suite/-/releases/release-v5.2.4).

## 7. Public pattern research → elite **FALTA** rows only

Patterns cited for gap identification only — **no new packs invented**.

| Public authority | URL | Elite **FALTA** row (honest gap) |
|---|---|---|
| [NIST SP 800-63B (Digital Identity Guidelines — Authentication)](https://pages.nist.gov/800-63-3/sp800-63b.html) | Memorized secret strength, replay-resistant authentication, verifier compromise detection | **FALTA:** target IdP password/MFA policy receipt binding (library delegates to external IdP; fixtures do not prove AAL/IAL) |
| [NIST SP 800-63B — Session management](https://pages.nist.gov/800-63-3/sp800-63b.html#sec5) | Session binding, re-authentication, termination | **FALTA:** step-up / re-auth policy per journey on production `application_url` (portal session proves refresh/logout locally only) |
| [OIDC Core 1.0 — Authentication](https://openid.net/specs/openid-connect-core-1_0.html#Authentication) | `code` + PKCE, `state`, `nonce` | **FALTA:** live OP conformance run bound to project issuer (`OIDC_CONFORMANCE` observation kind) |
| [OIDC Core 1.0 — Token validation](https://openid.net/specs/openid-connect-core-1_0.html#TokenValidation) | Issuer, audience, signing key rotation | **FALTA:** signing-key rotation drill on target IdP with session eviction evidence (`SESSION_REVOCATION_ROTATION`) |
| [OIDC Core 1.0 — Standard Claims](https://openid.net/specs/openid-connect-core-1_0.html#StandardClaims) | `sub`, `aud`, `exp` | **HECHO** locally for synthetic tokens; **FALTA** cross-tenant matrix on deployed tenants |
| [OIDC Core 1.0 — Refresh Tokens](https://openid.net/specs/openid-connect-core-1_0.html#RefreshTokens) | Rotation, reuse detection | **PARCIAL:** CAS owner proven in portal pack; **FALTA** provider-side rotation policy attestation |

## 8. What this slice does not do

- Does not edit `docs/ROADMAP.md` or franchise architecture docs (PR #13 scope).
- Does not invent SDR/GTM packs or workflows.
- Does not certify production or live IdP compatibility.
- Does not assert REVESTEX product admission.

## 9. Citations quick index

```text
markdown_system/PACK_PER_CLAIM_INDEX.md:22
markdown_system/CAPABILITY_CATALOG.md:46
markdown_system/LIBRARY_HEALTH_CHECK.md:74-75
markdown_system/TOTAL_SYSTEM_CAPABILITY_CONTRACT.md:43,49,88
implementation_packs/GO_OIDC_PORTAL_SESSION.md
implementation_packs/GO_OIDC_SERVICE_TOKEN_BROKER.md
implementation_packs/TYPESCRIPT_OIDC_PORTAL_ADAPTER.md
production_admission_gate/validate_identity_authorization_admission.py
reconstruction_evidence/IDENTITY_PORTAL_RELEASE_V402.md
reconstruction_evidence/OPENID_IDENTITY_AUTHORIZATION_PRODUCTION_ADMISSION_2026-08-29_V110.md
ci/local_identity_fixture.cjs
ci/verify_identity_security_fixture.py
```
