# Meta WhatsApp Cloud Adapter — 2026-08-26 V1

## Official authority

`fbsamples/whatsapp-api-examples` is an active official Meta sample repository. The admitted signed/verified commit is `de70ee908a67026e642aaee3703d20464e2a9466`; its archive is 348,229 bytes with SHA-256 `38d183a0d041d116dcf51356dd840c96690015f34c14c74b20681dd626ca9a7c`. Its 1,034-byte root LICENSE SHA-256 is `ef5c10eeba318e71ebf81a7300fc7c97a115e230dcb4647dc597946976069300` and permits use only with Facebook web services/APIs subject to Platform Policy.

Three exact code files are packaged byte-verbatim: signature validation SHA `6c052c13…12c`, Python message helper SHA `fef1aa03…702` and e-commerce incoming webhook SHA `a92a62bf…30f`. Markdown materialization adds a required final newline to LICENSE only; the packaged 1,035-byte hash `48d97b3c…d01f` and upstream raw hash are both retained, so LICENSE is marked `ADAPTED_FINAL_NEWLINE_ONLY`, not falsely `VERBATIM`.

The official Node SDK is archived. `UP-FAIL-026` keeps it rejected as a new production dependency. The official Python signature sample at this commit has an unconditional `return 'INVALID SIGNATURE HASH', 403` after its comparison, so matching signatures cannot reach the body handler. `UP-FAIL-025` preserves that defect; the file is reference-only and is never executed by the adapted path.

## Adaptation and executable evidence

`whatsapp_cloud.py` is explicitly `ADAPTED` under the Meta license. It retains the official Graph messages URL/Bearer JSON, GET subscription challenge and raw-body HMAC patterns, while adding a fail-closed approval profile, exact template/language/arity policy, environment-only secrets, constant-time full signature comparison, bounded HTTP, injected transport, atomic response/receipt, hashed identifiers and normalized webhook evidence without raw PII. Durable provider inbox/outbox idempotency remains a project composition gate.

Eleven files materialized with exact pack hashes. Eight offline tests passed from the authored tree and again from clean Markdown reconstruction: embedded source integrity, blocked profile/template, exact URL/header, atomic provider failure, subscription challenge, raw-body signature/tamper, webhook normalization/redaction and invalid-envelope atomicity.

## Failure learning

- `LIB-FAIL-066`: two browser-search calls returned no usable body; the audit switched to official GitHub API/raw endpoints and did not infer results.
- `LIB-FAIL-067`: the first reconstruction comparison caught LICENSE's missing upstream final newline; both raw/packaged hashes and adaptation status were recorded.
- `LIB-FAIL-068`: the first newline diagnostic repeated an invalid PowerShell object-to-pipeline form; collection and formatting were separated before measuring all four files.

This evidence proves source identity, license/provenance, reconstruction and local contract behavior. It does not prove current Meta terms, Graph version, account/template approval, recipient consent, delivery, pricing, opt-out, provider retry behavior or production reconciliation.
