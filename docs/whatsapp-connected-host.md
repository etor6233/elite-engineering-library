# Connected WhatsApp host

The electromobility API mounts the existing WhatsApp inbox, appointment notification and human reply modules when `WHATSAPP_ENABLED=true`. The same host starts the existing scoped job worker and constructs `app.New` once. Incoming text reaches `Runtime.Handle`, becomes a durable proposal and requires a different authenticated human reviewer before sending through the PostgreSQL outbound fence. The dispatcher is not run by this host.

Keep activation disabled while accounts are absent. Copy `docs/whatsapp-host-profile.template.json` to a private deployment configuration directory and bind its exact SHA-256 through `WHATSAPP_HOST_PROFILE_SHA256`. `WHATSAPP_HOST_PROFILE_FILE` must be absolute. The template is deliberately incomplete; the reference fixture is synthetic and is never evidence of real account access, terms or consent.

The provider profile is the existing `whatsapp_cloud/provider-profile.template.json` contract, with its own absolute path and exact hash. The fixed Python owner validates its original bytes at startup through a bounded, network-free command. A blocked provider profile prevents activation. Python executable, `whatsapp_cloud.py` and `status_reconciliation.py` also require their exact hashes. `Process` JSON property names follow the existing exported Go contract: `PythonExecutable`, `PythonSHA256`, `AdapterDirectory`, `AdapterSHA256`, `EvidenceDirectory`.

Tenant, organization, provider connection, contact consent purpose/policy and raw-payload retention approval must refer to the same materialized deployment records. The host never derives these from webhook input. Use one durable HMAC key across contact identity, outbound delivery and WhatsApp approvals. Create the evidence directory with private deployment permissions. The ordinary provider webhook registry and this receiver must not activate competing consumers for the same connection.

The domain gateway gets tenant and product/service/price-book configuration from this profile. Its legacy organization/lead defaults are both empty: `ScopedContactResolver` supplies the authorized contact organization and lead for each message. The configured model, instructions, finite token budget, retention and handoff text are passed to the existing app/runtime owners. This host does not make the current process-local economy/approval objects durable; the separate AI assurance gate retains that scope until its owner is completed.

Account secrets remain in externally protected files. The following environment values are absolute file paths, never secret contents:

```text
WHATSAPP_ACCESS_TOKEN_FILE=
WHATSAPP_APP_SECRET_FILE=
WHATSAPP_VERIFY_TOKEN_FILE=
WHATSAPP_SERVICE_TOKEN_FILE=
WHATSAPP_SERVICE_IDENTITY_PROFILE_FILE=
WHATSAPP_SERVICE_IDENTITY_PROFILE_SHA256=
WHATSAPP_SERVICE_CLIENT_SECRET_FILE=
WHATSAPP_CONTACT_HMAC_KEY_HEX_FILE=
LLM_API_KEY_FILE=
WHATSAPP_COLLECTOR_CLIENT_CERT_FILE=
WHATSAPP_COLLECTOR_CLIENT_KEY_FILE=
```

Select exactly one service identity owner: `WHATSAPP_SERVICE_TOKEN_FILE` supplied by an external token agent, or all three `WHATSAPP_SERVICE_IDENTITY_PROFILE_FILE`, `WHATSAPP_SERVICE_IDENTITY_PROFILE_SHA256`, `WHATSAPP_SERVICE_CLIENT_SECRET_FILE` values. Mixed or incomplete configuration fails closed. The OIDC profile uses the existing `elite.oidc.service-token.v1` contract. Its tenant, organization and required permissions must match the WhatsApp host before discovery or a grant is requested. The pinned OAuth/OIDC SDKs obtain and renew client-credentials tokens, verify RS256 and exact configured claims, and serialize concurrent renewal. The secret file is reread for each new grant. A failed renewal does not return a cached expired token.

WhatsApp token/app-secret/verify-token and externally supplied service tokens are reread when used. The actual OIDC verifier checks the selected service token on every poll and domain request; the selected tenant/organization and `whatsapp:process` plus `appointment:manage` grants are mandatory. Expired, malformed or wrong-scope tokens block work. JWT verification does not promise immediate remote revocation of an otherwise valid unexpired token; use the provider's admitted short token lifetime and disable the durable connection for local emergency denial. HMAC and LLM configuration are bound at startup and require a controlled restart to change. HMAC rotation also requires the existing identity/receipt migration procedure; changing the file alone cannot migrate stored identities.

The report collector uses TLS1.3 with exact configured CA bytes, server-name verification and a client certificate/key. The existing bounded JSON reporter writes only its six fixed operational fields. If reporting fails, the host cancels its worker and stops serving; reporting loss is not silently ignored. This validates authenticated transport and local writes, not collector retention/acknowledgement. Retention, supervision and restart policy have separate operations gates.

Public callback: `/v1/providers/whatsapp/webhook`. Human endpoints are the existing `/v1/franchise/whatsapp/replies` and appointment-notification routes, with their OIDC scope checks. Do not log callback subscription query tokens. The deployment supplies TLS/edge policy for the API itself.

Verification: three focused activation/file/identity cases, real mutual-TLS report with wrong-host/missing-key negatives, and full constructor with actual fixed Python validation and blocked-profile rejection. The constructor gate uses a lazy pool and performs no database statements or provider/model calls. Connected PostgreSQL/SDK/conversation behavior has its separate exact-source receipt. New host code is AUTHORED integration glue, with no upstream company authorship claim.
