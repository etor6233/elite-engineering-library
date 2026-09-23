# Meta Lead Reconciliation Adapter

Read-only recovery and reconciliation over the official Meta Business SDK for Python 26.0.1 and Graph API v26.0. It calls only `LeadgenForm.get_leads` or `LeadgenForm.get_test_leads`, preserves the exact exported provider rows and emits a separate candidate batch plus hash-linked receipt. The wrapper is local `ADAPTED` glue; it is not Meta-authored code.

The template starts blocked. Before use, prove the app and business registration, lead-access authorization, Page/form ownership, an official test-lead contract, PII storage controls, retention, quota/cost and reconciliation. Secrets belong only in the named environment variables. A retrieved lead is neither permission to contact nor an automatic CRM write. Form-field names not explicitly approved, malformed values and multi-value fields are retained only in the provider evidence and the corresponding candidate is rejected; a versioned project mapping is required before multi-select data can be promoted.

`TEST` calls the SDK's official `/test_leads` edge and marks every candidate `is_test=true`. `LIVE` calls `/leads`. The candidate batch carries tenant and organization scope independently so rejected raw evidence remains isolated even when `candidate=null`. The locally enforced row bound is recorded; a truncated retrieval is not a complete reconciliation and must be continued or rerun under the project's provider procedure.

The frozen lock is for CPython 3.14 on Windows x86-64. Meta's platform license restricts use of the SDK to Facebook services/APIs. Preserve the exact official LICENSE and the source/artifact lock.

```powershell
python -m venv .venv
.venv\Scripts\python -m pip install --require-hashes -r .\meta_lead_reconciliation\requirements-windows-py314.lock
Push-Location .\meta_lead_reconciliation
..\.venv\Scripts\python -m unittest -v test_fetch_leads.py
Pop-Location
```

No fixture proves live access, completeness, provider retention, rate limits, cost, consent or downstream delivery. The output remains candidate evidence until the central ingress and explicit promotion gates accept it.
