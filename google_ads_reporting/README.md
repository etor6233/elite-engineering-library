# Google Ads reporting adapter

This read-only adapter uses the exact official `google-ads==31.4.0` client and API `v25`. It loads credentials and developer-token configuration only through `GoogleAdsClient.load_from_env`, executes one bounded GAQL `SELECT` with `GoogleAdsService.search_stream`, preserves every modeled result row, and writes an atomic response/receipt. The receipt hashes the customer ID and query instead of exposing them and never authorizes a business write.

It does not create or mutate campaigns, upload conversions, choose attribution policy, or claim that report data is reconciled with orders. Before a real query the project must prove developer token/API access, OAuth identity, manager/login customer hierarchy, test account, query allowlist, field data policy, quota/cost, retention, logging redaction, reconciliation and rollback.

```text
python -m unittest -v test_reporting_query.py
python run_reporting_query.py --customer-id <digits> --query-file <approved.gaql> --output <new-directory>
```

The direct wheel and the complete CPython 3.14/Windows x86-64 dependency graph are hash-pinned. A different interpreter or platform must resolve, lock, scan and prove its own graph before deployment.
