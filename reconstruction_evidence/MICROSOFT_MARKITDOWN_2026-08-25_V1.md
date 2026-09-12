# Microsoft MarkItDown Evidence — 2026-08-25 V1

## Claim

This evidence proves exact acquisition and a locally executed document-focused test lane of Microsoft MarkItDown 0.1.7. It is admitted as a conversion component and as an official bridge to Azure Document Intelligence/Content Understanding. It is not an invoice truth engine and does not authorize automatic storage of extracted business fields.

## Exact source

```text
repository: microsoft/markitdown
release: v0.1.7
commit: fd239d5d2be43d9b68329730206b9312c7d5a388
archive URL: https://github.com/microsoft/markitdown/archive/fd239d5d2be43d9b68329730206b9312c7d5a388.zip
archive bytes: 4021911
archive SHA-256: 13b3d78a0ba3807df85208f662e89a2155b0137b1d3bc32f9635d03f329e1c70
license: MIT
LICENSE SHA-256: c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383
canonical acquisition: UPSTREAM_SOURCE_ACQUIRED id=microsoft-markitdown-0.1.7
files: 165
Python files: 72
test-tree files: 85
```

## Functionality inspected

The official package supports PDF, Word, Excel, PowerPoint, images, HTML, text/data formats, archives and optional media converters. Its release includes explicit integrations for Azure Document Intelligence and Azure Content Understanding. The latter can expose analyzer-defined structured fields; those fields inherit the selected Azure analyzer, account, schema and evaluation conditions.

The upstream security notice says conversion performs I/O with the current process privileges and recommends sanitizing untrusted inputs and using the narrowest local/stream conversion function. Elite therefore does not authorize arbitrary URL conversion or plugin enablement in a document-ingestion boundary.

## Test execution

The host has Python 3.14. The upstream `[all]` extra could not resolve because its `youtube-transcript-api~=1.0.0` constraint excludes Python 3.14. A document-specific environment was built with PDF, DOCX, XLS/XLSX, PPTX, Outlook and Azure DI/Content Understanding dependencies plus the official tests.

```text
document-focused suite:
273 passed, 5 skipped, 62 deselected

full-suite observations after compatible extras:
275 passed, 5 skipped, 60 failed
```

The 60 full-suite failures are retained. Most CLI cases spawn the global `C:\Python314\python.exe` rather than the isolated interpreter on this Windows host; media needs FFmpeg; file-URI behavior and the omitted incompatible YouTube graph also remain outside the admitted lane. This is not presented as a green full upstream suite.

The canonically acquired source then converted its official `packages/markitdown/tests/test_files/test.pdf` through the isolated module CLI. Output was 5.265 bytes, SHA-256 `a1305ace476784a37257a541262e3788886a87459f5402139e4cc28b448bdf55`, and contained the fixture's expected text. This proves the selected local PDF path, not field extraction accuracy.

## Admission

```text
archive/license identity: PASS
document conversion lane: PASS_WITH_EXPLICIT_SELECTION
full optional dependency graph: CONDITIONED_NOT_FROZEN
live Azure calls: NOT_RUN_MISSING_ACCOUNT
REVESTEX field accuracy: NOT_RUN_MISSING_CORPUS
admission: PINNED_CANDIDATE
```

Before production the project must lock every selected wheel, disable unused URL/plugins/media surfaces, malware-scan inputs, execute only local/stream conversions under least privilege, run Azure provider probes when selected, and evaluate fields through `PROJECT_DOCUMENT_INTELLIGENCE_DECISION.md`.
