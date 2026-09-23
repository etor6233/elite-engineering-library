# Meta Ads Reporting Adapter

Read-only wrapper over the official Meta Business SDK for Python 26.0.1 and Graph API v26.0. It executes only `AdAccount.get_insights`, preserves the provider rows and creates a hash-linked receipt. The local wrapper is `AUTHORED`; it is not Meta code.

The template is blocked. Before use, prove app/business registration, `ads_read`, a test ad account, exact account ownership, field allowlist, retention, quota/cost and reconciliation. Put app ID, app secret and access token only in the named environment variables. Never place credentials in JSON, Markdown, arguments or output.

The frozen lock is specific to CPython 3.14 on Windows x86-64. For another target, resolve and admit a separate hash-complete wheel graph. The official Meta license restricts the SDK and `capi-param-builder` to Facebook web services/APIs. The `capi-param-builder-python` 1.3.0 wheel omits its license file, so preserve the exact root LICENSE from official source commit `66a7eb84e29a65aa999089366c1f7c3c0e49bdb7`.

```powershell
python -m venv .venv
.venv\Scripts\python -m pip install --require-hashes -r .\meta_ads_reporting\requirements-windows-py314.lock
Push-Location .\meta_ads_reporting
..\.venv\Scripts\python -m unittest -v test_insights.py
Pop-Location
```

No fixture proves real permissions, field availability, provider totals, attribution, rate limits or billing. Admit production only after sandbox/test-account comparison and reconciliation against Meta UI/export for the selected dates and fields.
