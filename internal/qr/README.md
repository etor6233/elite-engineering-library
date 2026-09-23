# QR reference owner 0.2.0

AUTHORED successor to GO-QR-CORE 0.1.0. The checksum is unkeyed SHA-256: anyone can recompute it. It detects accidental field changes; it is not a signature, MAC, issuer proof or access token. Decode grants no authority.

Verify(ctx, encoded, expectedTenant, resolver) now requires the tenant from the authenticated SERVER session. Do not pass payload.TenantID, a user-selected tenant, body or query. There is intentionally no permissive three-argument compatibility wrapper. The resolver must be constructed server-side with actor and intended operation and reject unauthorized existing objects as well as absent ones. An existence-only resolver violates the contract. Resolver failures deny access. Cross-tenant references are rejected before resolution.

After resolving, the domain command owner must recheck object authorization, tenant, version, consent and idempotency at the actual effect. A successful reference lookup does not authorize delivery, payment, inventory change or document commit.

No camera, image decoder/encoder, 1D barcode, GS1, DataMatrix, OCR, vendor SDK or HMAC key is included. The old skip2/go-qrcode mention never supplied a pinned admitted encoder; it is not inherited as an implemented capability. Camera and WeChat candidates remain separate admission/runtime/device work. Use authorized record selection/manual entry while capture is unavailable.

Local tests use two real in-memory tenant records, recompute the forged checksum and assert denial under the other trusted tenant before any resolver call; they additionally check actor/operation-bound object permissions. They are not tests of a deployed identity service or physical scanner. Consumer migration is a required compile-time change to the new Verify signature; no production approval follows.
