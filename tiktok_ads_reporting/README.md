# TikTok Ads Reporting Adapter

Read-only wrapper over TikTok's official PyPI distribution `tiktok-business-api-sdk-official==1.1.3`. TikTok's provider-owned GitHub README recommends this distribution and PyPI identifies its author as TikTok Pte. Ltd. The exact universal wheel and four runtime dependencies are fixed by URL and SHA-256; no community substitute or locally invented SDK code is used.

The profile starts blocked. Prove Business/developer app access, reporting permission, test advertiser, exact dimensions/metrics/date semantics, retention, quota/cost and reconciliation. The access token belongs only in its named environment variable. This pack calls only official `ReportingApi.report_integrated_get`, which issues `GET /open_api/v1.3/report/integrated/get/`, preserves every page and emits a hash-linked receipt without advertiser ID or token.

```powershell
python -m venv .venv
.venv\Scripts\python -m pip install --require-hashes --no-deps -r .\tiktok_ads_reporting\requirements.lock
.venv\Scripts\python -m pip check
Push-Location .\tiktok_ads_reporting
..\.venv\Scripts\python -m unittest -v test_reporting.py
Pop-Location
```

Fixtures prove package/API import, pagination, provider errors, atomic evidence and official SDK endpoint/header contract. They do not prove account access, field semantics, provider totals, attribution, rate limits, cost or production reconciliation.
