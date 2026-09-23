# TikTok Lead Adapter

This pack uses the exact official TikTok Business API wheel `tiktok-business-api-sdk-official==1.1.3` and its generic `ApiClient.call_api` transport. The wheel's distribution metadata says `1.1.3`, while its own `business_api_client.__version__` constant says `1.2.1`; the artifact lock and tests preserve both upstream facts and use filename plus SHA-256 as executable identity. TikTok publishes the Lead retrieval and subscription contracts in its v1.3 portal, but the pinned SDK does not contain provider-generated Lead methods. `adapter.py` is therefore explicitly `ADAPTED`: it follows the official generated-method request shape without claiming TikTok authored this wrapper.

The provider profile starts blocked and is a mandatory constructor argument; client construction fails unless every required proof is `PROVEN`. Before any live call, prove account/app/Lead access, terms and consent, a test account, callback control and TLS, provider delivery, durable idempotency, authenticated retrieval reconciliation, retention/deletion and quota/cost. Secrets belong only in the named environment variables.

TikTok documents webhook delivery as at least once, so duplicates are normal. `webhook.py` parses the documented Lead envelope and emits stable hashes, but marks every inbound signal `UNAUTHENTICATED_PROVIDER_SIGNAL`: the reviewed public documentation did not demonstrate a signature mechanism. The pack never persists a webhook automatically. A project must use durable deduplication, retrieve/reconcile through the authenticated API and keep an authorized audit trail.

```powershell
python -m venv .venv
.venv\Scripts\python -m pip install --require-hashes --no-deps -r .\tiktok_lead_adapter\requirements.lock
.venv\Scripts\python -m pip check
Push-Location .\tiktok_lead_adapter
..\.venv\Scripts\python -m unittest -v test_tiktok_lead.py
Pop-Location
```

`run_retrieval.py` performs one approved authenticated GET and writes a provider SDK response, candidate batch and retrieval receipt atomically. The receipt links both data artifacts by SHA-256 and contains no lead fields. Dynamic values are admitted as candidates only when they are non-empty strings; arrays, objects, nulls or empty values remain intact in provider evidence and produce a rejected outcome requiring an approved mapping.

The sixteen tests prove mandatory profile enforcement, exact endpoint/method/header/body construction, input gates, secret-safe errors, official wheel transport identity, hash-linked retrieval artifacts, atomic no-overwrite output and deterministic parsing/deduplication signals. They do not prove live permissions, provider delivery, webhook authenticity, field semantics, consent, reconciliation, quotas, cost or production persistence.
