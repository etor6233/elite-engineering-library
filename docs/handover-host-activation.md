# Initial handover host activation

The API composes the existing `InitialHandoverModule` only after its materialized profile passes the owner's exact-byte loader and matches the active payment runtime. This wiring is AUTHORED configuration/orchestration glue. The owner retains its supported algorithm, transactional locks, stock/payment checks, idempotency and release evaluation.

`HANDOVER_ENABLED` absent or `false` omits the module and does not read a profile file or require payments. `true` requires an enabled, successfully prepared payment runtime and every `HANDOVER_*` input listed in `.env.example`. Any incomplete, mismatched or unrecognized activation stops startup before HTTP serving and before starting the payment worker. An enabled owner never falls back to the historical `LOCAL_FIXTURES` literal.

Materialize the selected profile using `tools/materialize_handover_profile.py` as described in `docs/handover-profile.md`. Preserve its `profile.json`, activation lock, decision and receipt. Set `HANDOVER_PROFILE_FILE` to that file and copy its exact profile ID, revision and SHA-256 into `HANDOVER_PROFILE_ID`, `HANDOVER_PROFILE_REVISION` and `HANDOVER_PROFILE_SHA256`. These values are configuration, not secrets. A document change requires a new reviewed revision and matching lock.

The host derives the activation's tenant, organization, provider, connection, account and expected mode directly from the prepared payment runtime. The loader must match those six fields to the hash-locked document before returning a valid contract. There are no duplicate `HANDOVER_*` scope overrides. Provider/account are connection bindings, not a new handover policy; payment storage and observations remain responsible for authenticating and reconciling provider facts. The host does not alter the contract after loading it.

After successful activation, these existing authenticated endpoints are registered:

- `POST /v1/franchise/orders/{id}/handover` prepares the initial handover using its existing command contract and idempotency key.
- `GET /v1/franchise/orders/{id}/handover-result` recovers the committed result using the same scope/key.
- `GET /v1/franchise/handovers/{id}/release-check` evaluates current eligibility against the expected observation hash.

All three require the existing `handover:manage` authorization and organization checks. Preparation and eligibility evaluation do not post a shipment or certify live production.

Selecting `--release-effect commercial-receipt` when materializing the profile
binds algorithm revision 2. The host then registers the commercial-release POST,
GET result recovery and GET current-validity endpoints documented in
`docs/commercial-release.md`, using the same profile-bound service instance.
Revision 1 keeps its existing read-only effect and does not mount these routes.
The immutable receipt records the committed commercial checkpoint; its historical
recovery is distinct from current eligibility after refunds, expiry or changed
observations. No fiscal rule, physical shipment or additional credential is
introduced by this host composition.
