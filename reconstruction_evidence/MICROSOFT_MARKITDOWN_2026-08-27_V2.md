# Microsoft MarkItDown Local Runtime — Reconstruction Evidence V2

Date: 2026-08-27  
Result: `REBUILD_VERIFIED / CONDITIONED`  
Production/admission result: **not promoted**; no real malware tools, business document, business schema, representative corpus or automatic storage authority was used or proven.

## Official Microsoft identity

- repository: `microsoft/markitdown`;
- current PyPI release checked on 2026-08-27: `0.1.7`, published 2026-07-29;
- annotated tag: `v0.1.7`; tag object `63714e47f11428189e2597cfd37a35068f82831d`;
- commit: `fd239d5d2be43d9b68329730206b9312c7d5a388`;
- GitHub verification: `verified=false`, reason `unsigned`; retained as `UP-FAIL-077` and never described as signed;
- exact source archive: 4,021,911 bytes, SHA-256 `13b3d78a0ba3807df85208f662e89a2155b0137b1d3bc32f9635d03f329e1c70`;
- MIT `LICENSE`: SHA-256 `c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383`;
- PyPI wheel: 71,093 bytes, SHA-256 `4eca912c87c6aa6897284a7f4bf6769a23bccf8544530f5d8b175fbe3797c916`;
- PyPI sdist: 51,767 bytes, SHA-256 `4d1f3c69cd43b82288fdc3653686d759dcf355ee7c681aa6a855aed98a1e4f44`.

Primary sources: https://pypi.org/project/markitdown/0.1.7/ and https://github.com/microsoft/markitdown/tree/fd239d5d2be43d9b68329730206b9312c7d5a388. Microsoft explicitly warns that MarkItDown performs I/O with the current process privileges and recommends sanitized inputs plus the narrowest conversion function. The local wrapper invokes only `MarkItDown(enable_plugins=False).convert_local()`.

## Gap proved, not hidden

The previous V0.1 wrapper accepted a path and limits, trusted the file extension and converted immediately. Its positive test passed without a configured profile or receipt from the local security gate. That proves the old architecture allowed security-admission bypass; it is retained as `LIB-FAIL-248`.

The V0.2 correction is local `AUTHORED` integration, not Microsoft code. It requires:

- profile schema `elite-markitdown-local-conversion-profile/v2`;
- exact provider/runtime/package identity;
- format, MIME and limits selected only by the profile;
- zero `REQUIRED` formats in the distributed template;
- URL, plugin, LLM, archive expansion and automatic business storage explicitly false;
- `elite-secure-local-file-receipt/v1` with `ADMITTED / ALL_SELECTED_GATES_PASSED`;
- matching input SHA-256, bytes and MIME;
- approval/policy plus ClamAV, Magika and YARA-X binary/version/result evidence;
- ClamAV exit zero and zero YARA-X matches.

The CLI now exposes only input, output, profile and security-receipt paths. It cannot override formats, MIME, limits or policy. Before calling Microsoft code, every profile/security invariant must pass. The V2 conversion receipt hash-links profile, security receipt, input and output and always declares `business_storage_authorized=false`.

## Exact environment and execution

The unchanged target lock contains 44 exact CPython 3.12 Windows x86-64 wheels with SHA-256. A clean venv was rebuilt and installed with `--require-hashes --only-binary=:all:`. The first install remained active after an apparent process close and left only 8/44 packages; `pip check` alone was falsely green. The exact environment verifier caught it. After stopping only the three audit PIDs, recreating the venv and waiting for real process completion:

- `pip check`: PASS;
- exact environment verifier: `MARKITDOWN_WIN_PY312_ENVIRONMENT_PASS distributions=44`;
- official package import: `markitdown==0.1.7` PASS;
- official `MarkItDown(enable_plugins=False).convert_local` callable: PASS;
- official OSV batch at `2026-08-27T16:33:25Z`: 44 packages queried, zero findings returned;
- seven wrapper tests: PASS;
- clean Markdown materialization: 6/6 files with matching hashes;
- seven tests repeated from canonical materialization: PASS.

A real local `.txt` fixture was converted through the installed official wheel, not through the fake test converter. The input was 72 bytes/SHA-256 `7a98ad5afc70821524da5f632c80296a5c3824f8109ddd34eca6f0873155062f`; output was 72 bytes with the same hash. The receipt contained exact profile and security-fixture hashes and storage false. The security receipt was an explicitly labeled structural audit fixture, not evidence that ClamAV/Magika/YARA-X ran on a business document.

## Negative coverage

The seven-test suite proves that the converter is not invoked and no output is committed for a blocked template/format, package mismatch, MIME mismatch, duplicate extension, storage-enabled profile, absent/rejected/mismatched/incomplete/storage-enabled security receipt, disguised content, undeclared archive, empty/oversized provider output, occupied output or wrong installed distribution. CLI introspection proves no policy override exists.

`VERIFY_LIBRARY.ps1` passes 41 packs/366 files and composes the MarkItDown lane as 4 packs/35 files: exact acquisition, secure local-file gate, V0.2 runtime and strict field evaluation. `VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Audit` passes and always materializes the pack, checks the 44-entry lock, zero-enabled V2 profile and seven regressions. With explicit `-AllowNetwork` and exact CPython 3.12 Windows, it now installs the hash lock and requires `pip check`, 44/44, tests and official import.

## Failures converted into regression

- `LIB-FAIL-246`: wrong materializer parameter; corrected by reading the canonical signature and rebuilding 6/6.
- `LIB-FAIL-247`: Windows glob passed as a path to ripgrep; corrected with `-g`.
- `LIB-FAIL-248`: security/profile bypass; corrected by mandatory V2 profile and hash-linked security receipt.
- `LIB-FAIL-249`: remembered Google pack filename did not exist; corrected by enumerating the canonical path.
- `LIB-FAIL-250`: orphaned installer and misleading partial `pip check`; corrected by exact process inventory, clean venv and 44/44 verifier.
- recurrent `LIB-FAIL-130`: null OSV results were again counted as empty findings; corrected by requiring non-null vulnerabilities with non-empty IDs, yielding 44 queried/0 findings.
- recurrent `LIB-FAIL-233`: delete+add same path in one patch was rejected; the unchanged files were verified, operations separated and the full reconstruction repeated.

All are now `REGRESSION_PROVEN` in the governing ledger.

## Honest boundary

MarkItDown converts content to Markdown; it does not prove invoice, proforma, packing-list, bill-of-lading, customs, product or financial field truth. Real project use still requires owner-approved formats/MIME/limits, real security assets and receipts, least-privilege sandboxing, an approved semantic extractor/schema, representative ground truth, strict zero-false acceptance, review/rejection routing, privacy/retention, load, recovery and rollback. No automatic business persistence is authorized by this lane.
