# V402 policy profile and host integration evidence

PASS for the narrow local configuration/glue delta. D is integrated; three packs
are unpublished candidates. This is not TEST02/global source closure, production
readiness or vendor-derived business policy. All new code/specification is
AUTHORED. No upstream source was acquired, copied or falsely attributed.

## Delivered scope

10 new files and 5 existing files: 8 prior policy files, 2 host helpers, 3 exact
policy callers, main and .env.example. All 11 original policy before hashes were
validated before any D write. Main before SHA was
e2d8321dc8fbeb5ccaca48adcd25f5375e984ae72b7df57a1ef36c4fecb3de44;
main after SHA is
0b726aab72cc2f2e033c658b699e5cf70b381e1499379980eead6f15f2b0235a.
ProviderObservedPayments:true is preserved. This agent did not change payment.go,
handover.go, handover files or DDL. Parent's subsequent commercial-release changes
in handover.go are independent and are not included in this delta or its tests.

BUSINESS_POLICY_PROFILE_FILE and BUSINESS_POLICY_PROFILE_SHA256 jointly select
an immutable configuration revision; both absent select the exact embedded
reference. Incomplete, unreadable, hash-mismatched or incompatible input fails
before OIDC discovery, DB work and serving. The host passes the same *Profile to
both persistence owners and to the Journey service. Public listing and booking
share the configured lead time; persistent effects and custom receipts bind its
hash. Reference receipt hashes retain historical identity.

## G0–G8, local authored configuration/glue scope

| Gate | Result and evidence |
|---|---|
| G0 | PASS: local AUTHORED identity is explicit for every added block; before/after hashes and exact D files are recorded. No BC algorithm, vendor provenance or new upstream source is claimed. |
| G1 | PASS for authorized internal library use: no upstream code or new dependency/license is added. LicenseRef-Workspace-Owner is preserved, not represented as an OSI license or public redistribution grant. Existing notices retain their owners. |
| G2 | PASS narrow scope: PBC-CORE's tested-extension boundary and existing migration0007/0009/0015 constraints are documented. The base schema is unchanged; incompatible modes are rejected. Remaining lifecycle/resource/consent rules are explicitly excluded, never reclassified as glue. |
| G3 | PASS: bounded loader -> immutable Profile -> shared host owners -> existing transactions, quota/interval primitives and outbox. No new data owner, migration, background worker or control plane. The main patch is explicit. |
| G4 | PASS: focused behavior/invalid/binding/replay tests; prior PG18.6 delta with54 unchanged migrations and no skips; current two host tests plus six real executable startup cases; host vet/build. PG implementation bytes are identical to the prior receipt; no PG repetition. |
| G5 | PASS narrow delta: loader rejects missing/unknown/duplicate/case-aliased keys, malformed UTF-8/JSON, invalid bounds and hash mismatch. Configuration is operator-selected, never request-selected; errors omit file contents/credentials. Same-profile construction and changed-profile replay rejection tested. Tenant/auth/payment code is not replaced. This does not close broader TEST03. |
| G6 | PASS bounded local scope: configuration max16384bytes/depth16, one-time startup parsing, immutable small profile, no new queue/network worker. Finite Go1.26.8 parser fuzz53,626executions/3seconds PASS reused by exact hash. No throughput, tail-latency or production-load claim. |
| G7 | PASS narrow scope: logged startup failure before external work; policy hash on new effects; previous profile bytes/hash and receipts retained for manual configuration rollback. Profiles cannot silently change modes or rewrite old rows. Different-profile retries conflict; existing scoped result reads remain. No automated multi-instance rollout claim. |
| G8 | PASS reconstruction: policy0.1.0 + Commerce0.6.4 + Journey0.10.21 reconstruct51/51 blocks without duplicate owners; all13 changed/new candidate blocks match D. Main/environment have exact separate patches for their existing owners. Publication/global profile promotion and the parent's independent payment/commercial overlays remain root responsibilities. |

The new pack retains SUPPORTED_REFERENCE / REBUILD_VERIFIED / CONDITIONED until
its explicit composition/publication conditions are reflected in the canonical
profile. G0–G8 here audit this local candidate and do not promote all inherited
AUTHORED owner code. License and live deployment conditions are not hidden.

## Results and remaining integration

The six executable cases ran with no live account or connection: absent profile
and a valid file passed policy selection and reached the existing missing general
connection configuration check; path-only/hash-only/wrong-hash/hash-matching but
incompatible data failed at policy selection first. No tests skipped. Host tests
also cover shared Profile, nil bindings, missing file and no silent fallback.
The parent has independently read/built the integrated main and is integrating
commercial release without modifying that file.

At publication, merge the root's pending Commerce HTTP changes into the canonical-
based Commerce candidate. Publish main in GO-ELECTROMOBILITY-APPLICATION and .env
in the selected TS-GO-API-WEB-BRIDGE owner. The two new host helper files belong to
GO-BUSINESS-POLICY-PROFILE and must not be duplicated in the host pack. Do not
replace D from the old snapshot or erase the payment/handover overlays.

## Exact artifacts

| Artifact | SHA256 |
|---|---|
| `policy-host-manifest.json` | `9dfa4903894d0c8f74aa7b502889539de235e150480dbf9c18dd5c5fbdb3f884` |
| `policy-host.patch` | `ad107d773d8a33c5124e172facd4a09ebed25c1e38b7550aa013920853968bdf` |
| `policy-host-main.patch` | `7fac599929dac64b81a28f5b3821a5323291e4970c69881062b7c333d07328c3` |
| `policy-host-tests/result.json` | `8a95a17262fc84b595317e0412c92a61db0e787d2e51a3738401b2fd226ff706` |
| `policy-roundtrip-receipt.json` | `7a5531ff161ae697f1ab8b302f7d1d94be908dade8feb7469dac7a01175ba915` |
| `policy-candidate-pack-identities.json` | `e5eb7cb3c103f8c791eb03f6ee79413763cd17911d4c982486356f129a94e948` |
| `policy-config-pg-2/result.json` | `11823cda4a3038662c47296efe1af36ddd8fcb25b0320a690f82a72a630429b4` |
| `policy-fuzz-receipt.json` | `423e4f4ea56ca8d5186e771a8a30496707d25450fbefc435b57c62047f3de983` |
| `policy-candidate-packs/GO_BUSINESS_POLICY_PROFILE.md` | `e7150c66a1f055d249f45d10b32e14336ae54115a912b1c4e1908bc4f96a0711` |
| `policy-candidate-packs/GO_COMMERCE_PRICING_PAYMENT_API.md` | `c8d87f07ed843cd62103db31761340ab3fbf609b7af22a79bea40c632e36f36e` |
| `policy-candidate-packs/GO_FRANCHISE_CUSTOMER_JOURNEY_API.md` | `116535f5785fe4f00424a3047389d8bccfb25b6d42d2d897851c11d6842c6f01` |
