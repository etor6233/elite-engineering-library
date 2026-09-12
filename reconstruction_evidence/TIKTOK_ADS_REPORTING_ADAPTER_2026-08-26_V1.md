# TikTok Ads Reporting Adapter — 2026-08-26 V1

## Authority and identity

Official repository `tiktok/tiktok-business-api-sdk` was observed active and MIT-licensed. TikTok publishes no GitHub release/tag for the admitted snapshot. The exact head commit `f809c396520df2d7b201a9ccc5378d822b728ed3` was reported by GitHub as signature verified/valid. Its official source archive was 3,515,705 bytes with SHA-256 `d77e75ba8f02d04e53d03f3b3fc4c90a57288bdbceeffdd96cb4ebdd5f8834cb`; `LICENSE.md` was 1,070 bytes with SHA-256 `991c218ea2332d85698a9f2bae9e15a035492db91e7b4b2ffba18c8c8863ca03`.

Five critical source files were independently hashed. The official `ReportingApi.report_integrated_get` issues `GET /open_api/v1.3/report/integrated/get/` and places the credential in `Access-Token`. The source `setup.py` declares distribution metadata version 1.0.0 while `Changelog.md` declares 0.1.8. `UP-FAIL-023` preserves the missing release/tag/wheel and version conflict; runtime identity is the signed commit, archive and critical-file hashes.

## Executable reconstruction

The authored pack contains eight files. A clean CPython 3.14 Windows x86-64 virtual environment installed four exact universal wheels with `--require-hashes`; `pip check` passed. Six adapter tests passed against the official source, including endpoint/header inspection, bounded pagination, provider rejection, 100-page safety limit, atomic cleanup, approval allowlists and secret/account redaction. A research receipt matching `elite-official-source-receipt/v1` plus the exact archive/commit allowed the wrapper to verify all five critical files. The project still requires its own profile approval and acquisition receipt; this research receipt is not project authorization.

The complete upstream model tree emitted `FileNotFoundError` while `compileall` attempted to create very long `__pycache__` paths on Windows/Python 3.14, even though the command returned zero and the six scoped tests passed. This is retained as `UP-FAIL-024`; no full-tree portability claim is made. The adapter's own eight materialized files matched source hashes 1:1.

## Failure learning

- `LIB-FAIL-058..060`: two malformed PowerShell inventory pipelines and one guessed upstream filename were rejected; direct enumeration established exact wheels/files.
- `LIB-FAIL-061`: the first SDK test used the parent expansion path instead of the archive's single top directory; the official module could not import until the exact source root was enumerated.
- `LIB-FAIL-062`: treating mass `compileall` as a clean gate would have hidden its non-fatal Windows path errors; the condition is now separately retained upstream.
- `LIB-FAIL-063`: test bytecode under the research tree caused a recursive source comparison to report an extra file; verification was repeated against the exact eight-file manifest and materialized output count.
- `LIB-FAIL-064`: a diagnostic command guessed a standalone upstream failure ledger although failures live in `LIBRARY_FAILURE_LEARNING_LEDGER.md`; canonical routing was re-read and no second ledger was created.
- `LIB-FAIL-065`: the first ledger update duplicated a canonical failure token inside a narrative cross-reference; the uniqueness gate rejected it and the declaration/reference syntax was separated.

This evidence proves reconstruction and local official-SDK contract behavior. It does not prove TikTok account access, current production metric semantics, attribution, provider totals, rate limits, cost or reconciliation.
