# Start the V403 local reference

[Plan](UNIFIED_REFERENCE_PACK_PLAN_V403_R4.md) · [Pack](../implementation_packs/UNIFIED_REFERENCE_V403_R4.md) · [Evidence](../reconstruction_evidence/UNIFIED_REFERENCE_V403_R4.md)

This is a library reference, not a new franchise or production approval. Keep the library checkout separate from bundle, package, runtime data and reference destinations. These outputs contain source/evidence; do not delete them as caches.

From the library root, use three absent directories you choose:

```powershell
$Bundle = "<bundle-new-directory>"
$Package = "<package-new-directory>"
$Reference = "<reference-new-directory>"
$Proof = "<new-source-receipt.json>"
if (Test-Path -LiteralPath $Bundle) { throw "Bundle destination must not exist" }
pwsh -File ./materialize_markdown_pack.ps1 -PackFile ./implementation_packs/UNIFIED_REFERENCE_V403_R4.md -Destination $Bundle
if ($LASTEXITCODE -ne 0) { throw "Carrier materialization failed" }
python "$Bundle/unified_reference_v403/prepare_bundle.py" --destination $Package
if ($LASTEXITCODE -ne 0) { throw "Portable package preparation failed" }
python "$Package/materialize.py" --destination $Reference
if ($LASTEXITCODE -ne 0) { throw "Reference materialization failed" }
python "$Bundle/unified_reference_v403/VERIFY_SUCCESSOR.py" --package $Package --reference $Reference --library-root . --receipt $Proof
if ($LASTEXITCODE -ne 0) { throw "Source/publication verification failed" }
```

The original Markdown tool can overwrite existing files: the explicit absent-bundle guard above is mandatory, and that directory must have one writer. The new preparation/materialization tools reject existing outputs themselves.

The receipt verifies sources and publication routing, not runtime execution. Use the selected reference's runtime entry and its separate evidence for startup/fixtures. Product credentials, hardware, hosted CI/cloud, human visual acceptance and production remain distinct.

Any agent resumes through the existing project bridge and operating protocol, then the consumer's own state/events and exact next action. Search one authority/owner at a time. Use functional XERJ if available, otherwise bounded rg; token savings are not measured. A platform instruction template is not proof every vendor executed it. Never treat library-maintenance state as a consumer's state.

## Reproduce the local Windows fixture qualification

Prerequisites: Windows x64, PowerShell 7, Python 3.14, Node 24 and pnpm compatible with the selected package. Go and PostgreSQL use pinned official archives. Network access is required for missing dependencies; no live account or credential is required. A failed prerequisite is a failure, not a PASS. Keep these directories outside the library; retain all evidence.

Fill the paths below with your chosen absolute paths. `$Package` and `$Reference` are the two outputs already materialized above; `$Tools`, `$PG` and `$Evidence` must not exist yet. `$Cache` is a tool cache only, never the reference or evidence directory.

```powershell
$Carrier = "$Bundle/unified_reference_v403"
$Package = "<package>"
$Reference = "<reference>"
$Tools = "<new-runtime-tools>"
$Cache = "<tool-cache>"
$PG = "<new-postgres-runtime>"
$Evidence = "<new-qualification-evidence>"
python "$Carrier/prepare_runtime.py" --destination $Tools
if ($LASTEXITCODE -ne 0) { throw "Runtime tool reconstruction failed" }
python "$Tools/bootstrap_go.py" --cache "$Cache/go"
if ($LASTEXITCODE -ne 0) { throw "Go bootstrap failed" }
python "$Tools/runtime-support/prepare_postgres.py" --archive "$Cache/postgresql.zip" --runtime $PG --download
if ($LASTEXITCODE -ne 0) { throw "PostgreSQL preparation failed" }
# If that archive already exists, omit --download; it is rechecked before extraction.
$BrowserGate = "$Reference/microsoft_playwright_browser_gate"
pnpm --dir $BrowserGate install --frozen-lockfile --ignore-scripts
if ($LASTEXITCODE -ne 0) { throw "Browser gate dependencies failed" }
pnpm --dir $BrowserGate exec playwright install chromium
if ($LASTEXITCODE -ne 0) { throw "Chromium installation failed" }
python "$Tools/run_local_reference.py" --target $Reference --source-package $Package --evidence $Evidence --go "$Cache/go/go1.26.8/go/bin/go.exe" --postgres-runtime $PG --playwright "$BrowserGate/node_modules/@playwright/test/index.mjs"
if ($LASTEXITCODE -ne 0) { throw "Fixture qualification failed; preserve its evidence" }
python "$Carrier/VERIFY_SUCCESSOR.py" --package $Package --reference $Reference --library-root . --receipt "$Evidence/SOURCE_AFTER_RECEIPT.json"
if ($LASTEXITCODE -ne 0) { throw "Post-runtime source verification failed" }
```

This harness builds and starts the reference with synthetic OIDC/provider fixtures, exercises four journeys and resumes from its saved state, then stops its own processes. It is a qualification run, not a persistent development server. Its own receipts determine PASS/FAIL; installing the carrier never generates a runtime PASS. Reuse of a PostgreSQL evidence directory is deliberately rejected. The source checker permits declared generated caches only and restores the exact known Next type seed after stopping; it does not accept source drift.

For an authorized future consumer or a disposable bridge test (this command does not create a franchise), use the original library bridge through the admitted complementary installer:

```powershell
$Consumer = "<consumer-or-disposable-directory>"
pwsh -File "$Tools/successor-tools/operating_connection/INSTALL_PROJECT_OPERATING_BRIDGE.ps1" -LibraryRoot $PWD -ProjectRoot $Consumer
if ($LASTEXITCODE -ne 0) { throw "Operating bridge failed" }
python "$Tools/successor-tools/operating_connection/connection.py" verify --project-root $Consumer
if ($LASTEXITCODE -ne 0) { throw "Operating bridge verification failed" }
```

The complementary bridge requires the exact library bridge revision pinned in its compatibility contract and the existing operating protocol. It never rewrites the frozen V402 bridge. A different checkout must be reconciled explicitly; vendor templates alone do not establish support.
