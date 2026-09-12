# All Implementation Packs Materialization — 2026-08-27 V31

Status: `NOT_READY_UNDER_EXPANDED_USER_STANDARD`  
This snapshot supersedes V30 for current library structure and execution evidence. Historical evidence remains immutable context.

## Current reproducible inventory

- implementation packs: 41;
- materializable files across packs: 366;
- Markdown files after this snapshot: 265;
- composition profiles: 20;
- exact public upstream source archives in lock: 72;
- official document SDK artifacts: 5;
- Business Central executable artifact profile: 1;
- provider adapters counted by executable audit: 7.

The structural verifier passes. Key compositions include initialization 2/24, Azure document 5/41, Google document 5/41, AWS Textract 4/36, MarkItDown local 4/35, secure local file 2/24, strict field evaluation 2/23, backend 16/95 and web/BFF 3/44.

## V31 delta

`MICROSOFT-MARKITDOWN-LOCAL-RUNTIME` is now V0.2.0. The previous runtime could convert based on extension without a configured profile or secure local-file receipt. The corrected local `AUTHORED` integration requires a zero-enabled-by-default V2 profile and a complete input-matching receipt from ClamAV/Magika/YARA-X before invoking official Microsoft MarkItDown 0.1.7. Policy values cannot be overridden by CLI; the V2 receipt hash-links profile, security, input and output and keeps storage false.

Microsoft MarkItDown 0.1.7 remains the current PyPI version checked on 2026-08-27. Source/tag/commit, wheel, sdist and MIT license remain exact. The tag is unsigned and that condition is retained. A clean CPython 3.12 Windows x86-64 environment installed 44/44 hash-locked wheels, passed `pip check`, exact environment verification, seven tests, the official import and a dated OSV batch with 44 queried/zero findings. One real local fixture converted through the official wheel; it used an explicitly non-production security fixture and therefore does not prove real malware scanning or business accuracy.

The MarkItDown composition now contains four packs/35 files: acquisition, secure local-file gate, runtime and strict field evaluation. The executable audit checks its lock/profile/tests every run and exposes a network-authorized exact-runtime installation gate.

`VERIFY_LIBRARY.ps1` passes 41 packs/366 files and all 20 profiles. `VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Audit` passes with 72 upstream sources, five document SDK artifacts, one Business Central artifact and seven provider adapters.

## Honest boundary

No pack is promoted to unconditional production reuse. MarkItDown conversion is not semantic field truth; its business document corpus, schemas and real security receipts do not exist here. Google, AWS, Azure, Business Central and other provider/platform conditions also remain. REVESTEX-specific tenancy, fiscal rules, real access, integrations, field accuracy, production security, load, restore and rollback still require project inputs and target evidence.

The library is stronger and generates a safer conditioned local document lane, but the expanded user goal remains unfinished. The next review must continue closing missing end-to-end enterprise capabilities with exact official sources and executable evidence; it must not infer completion from this snapshot.
