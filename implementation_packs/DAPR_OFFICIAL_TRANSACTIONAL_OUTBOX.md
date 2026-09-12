# Dapr Official Transactional Outbox

## 1. Metadata

```yaml
pack_id: "DAPR-OFFICIAL-TRANSACTIONAL-OUTBOX"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa bytes oficiales firmados de Dapr Runtime 1.18.3 y Dapr Docs v1.18 para transactional outbox, junto con locks, adquisición y gates fail-closed; rechaza explícitamente el SDK Go oficial vulnerable y no inventa lógica de negocio."
stacks: ["PowerShell 7", "Go 1.26.6 upstream", "Dapr Runtime 1.18.3"]
compatible_with: ["OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.2.x", "GO-RELIABLE-ASYNC-WORKERS 0.2.x", "POSTGRES-TRANSACTIONAL-FOUNDATION 0.2.x"]
incompatible_with: ["Dapr Go SDK 1.14.2", "entrega exactly-once no demostrada", "store sin transacciones/outbox", "consumidor sin idempotencia", "selección automática de componentes", "producción sin rollback"]
license_expression: "LicenseRef-Workspace-Owner AND Apache-2.0"
upstream_sources: ["https://github.com/dapr/dapr/tree/2dcf2e3548faf3c740b3d29ee300e2132db65cff", "https://github.com/dapr/docs/tree/5958a7e19a04e199326a6c5321bdd7714ee83b4c"]
verified_at: "2026-08-28"
```

## 2. Applicability

Use este pack únicamente cuando el blueprint seleccione Dapr como plataforma de state transaction + outbox y el proyecto pueda demostrar store transaccional compatible, pub/sub real, permisos, entrega duplicada, consumidor idempotente, recuperación y rollback. Conserva código y pruebas publicados por Dapr, no una implementación local equivalente.

No lo use como backend automático ni como promesa exactly-once. Dapr declara at-least-once. El SDK Go 1.14.2 y la rama principal firmada observada quedan excluidos por dependencias vulnerables alcanzables; el runtime no autoriza un SDK rechazado.

## 3. Architecture contract

Trece archivos son VERBATIM: documentación/licencia del commit firmado de dapr/docs y licencia, interfaz, fake, implementación, unit tests y cinco escenarios de integración del release firmado Dapr Runtime 1.18.3. Los siete archivos AUTHORED sólo fijan identidad, adquisición, aprobación y readiness; no contienen reglas de negocio ni reemplazan componentes Dapr.

El runtime publica el evento dentro de la transacción del state store y luego lo entrega por pub/sub con semántica at-least-once. Por ello la admisión exige idempotencia del consumidor y prueba de duplicados. El archive completo se obtiene sólo desde cache aprobada o red doblemente autorizada, se valida por bytes/hash/commit y nunca se superpone sobre un destino existente.

## 4. Exact file manifest

```text
CREATE dapr_official_transactional_outbox/README.md
CREATE dapr_official_transactional_outbox/official-source-lock.json
CREATE dapr_official_transactional_outbox/source-approval.template.json
CREATE dapr_official_transactional_outbox/project-selection.template.json
CREATE dapr_official_transactional_outbox/acquire_official_runtime.ps1
CREATE dapr_official_transactional_outbox/validate_project_selection.ps1
CREATE dapr_official_transactional_outbox/test_embedded_sources.ps1
CREATE dapr_official_transactional_outbox/test_acquisition.ps1
CREATE dapr_official_transactional_outbox/upstream/docs/LICENSE
CREATE dapr_official_transactional_outbox/upstream/docs/howto-outbox.md
CREATE dapr_official_transactional_outbox/upstream/runtime/LICENSE
CREATE dapr_official_transactional_outbox/upstream/runtime/pkg/outbox/outbox.go
CREATE dapr_official_transactional_outbox/upstream/runtime/pkg/outbox/fake/fake.go
CREATE dapr_official_transactional_outbox/upstream/runtime/pkg/outbox/fake/fake_test.go
CREATE dapr_official_transactional_outbox/upstream/runtime/pkg/runtime/pubsub/outbox.go
CREATE dapr_official_transactional_outbox/upstream/runtime/pkg/runtime/pubsub/outbox_test.go
CREATE dapr_official_transactional_outbox/upstream/runtime/tests/integration/suite/daprd/outbox/http/basic.go
CREATE dapr_official_transactional_outbox/upstream/runtime/tests/integration/suite/daprd/outbox/http/delete.go
CREATE dapr_official_transactional_outbox/upstream/runtime/tests/integration/suite/daprd/outbox/http/projection.go
CREATE dapr_official_transactional_outbox/upstream/runtime/tests/integration/suite/daprd/outbox/grpc/basic.go
CREATE dapr_official_transactional_outbox/upstream/runtime/tests/integration/suite/daprd/outbox/grpc/projection.go
```

## 5. Materialization blocks

### FILE: `dapr_official_transactional_outbox/README.md`
```yaml
block_id: "DAPR-OFFICIAL-OUTBOX:01:v1"
operation: CREATE
provenance: AUTHORED
source: "local admission/acquisition code constrained by exact Dapr upstream evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "86e7f77d25f35d7d9869a1c57ef2f4cbabdeff5c2a8585091fdae6c15f410f8c"
variables: []
secrets_allowed: false
```
````markdown
# Dapr official transactional outbox source

This pack preserves exact Apache-2.0 bytes from Dapr Runtime `v1.18.3` and Dapr Docs branch `v1.18`, each pinned to a signed commit. It includes the runtime outbox contract, implementation, unit tests and HTTP/gRPC integration scenarios published by Dapr.

Admission is platform-conditioned. Dapr documents at-least-once delivery; every consumer must therefore prove idempotency and duplicate handling. The selected state store must support transactions and the Dapr outbox feature, and the selected pub/sub component must be tested with the real project account and broker.

The current Dapr Go SDK is deliberately not included. Official release `v1.14.2` and signed `main` commit `9d745f76a46e23ae0d022f1abb098132b67603e5` retain dependency versions affected by reachable vulnerabilities in the dated scan recorded by this pack. Do not locally bump or patch them and call the result official. Re-admit only a later official Dapr SDK release or commit after identity, license, full graph, offline tests and reachability scan pass.

Use order:

1. inspect `official-source-lock.json` and exact upstream files;
2. complete `project-selection.template.json` without secrets;
3. run `validate_project_selection.ps1`; it must fail until real evidence is supplied;
4. optionally acquire the complete official runtime archive with an approved, hash-bound record;
5. run the upstream unit and integration suites from the complete source in a supported environment;
6. configure the project from Dapr component manifests chosen for that project, not from invented defaults;
7. prove duplicates, state-store transactions, broker recovery and rollback before production authorization.

The embedded tests are executable evidence and design reference; the subset is not a standalone fork of Dapr. The complete runtime source is acquired from the official archive lock.
````

### FILE: `dapr_official_transactional_outbox/official-source-lock.json`
```yaml
block_id: "DAPR-OFFICIAL-OUTBOX:02:v1"
operation: CREATE
provenance: AUTHORED
source: "local admission/acquisition code constrained by exact Dapr upstream evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "45ecc9dbfed4bc5262ae11667704b78ccf8310c4209a8229e935a137ac40d950"
variables: []
secrets_allowed: false
```
````json
{
  "schemaVersion": 1,
  "verifiedAt": "2026-08-28",
  "runtime": {
    "repository": "https://github.com/dapr/dapr",
    "tag": "v1.18.3",
    "commit": "2dcf2e3548faf3c740b3d29ee300e2132db65cff",
    "commitSignatureVerified": true,
    "archiveUrl": "https://github.com/dapr/dapr/archive/2dcf2e3548faf3c740b3d29ee300e2132db65cff.zip",
    "archiveBytes": 20269733,
    "archiveSha256": "dd64436b01c06bdf59db6c7c949dce81eff36c293581980af93e12a80bfd6681",
    "license": "Apache-2.0",
    "licenseSha256": "9523eacf5421b65637420755c00b5175f7804b6a6eb736b86cd10ea6b3937dcf",
    "goModSha256": "4f1941ca734128e36ae3c059d862b9d92759dab63de5f860b4e3a525c0998d9d",
    "goSumSha256": "e7c4c2dba67d067423e9cd85a8d13da8a86028cfd81d5b831a74b2ecb8c95e79",
    "goVersion": "1.26.6",
    "moduleGraphCount": 1338,
    "focusedTests": "go test ./pkg/outbox/... and six pkg/runtime/pubsub outbox tests with twenty subtests",
    "focusedTestsPassed": true,
    "offlineModuleVerificationPassed": true,
    "govulncheckVersion": "v1.7.0",
    "vulnerabilityDbUpdatedAt": "2026-08-27T19:52:28Z",
    "focusedReachableVulnerabilities": 0
  },
  "docs": {
    "repository": "https://github.com/dapr/docs",
    "branch": "v1.18",
    "commit": "5958a7e19a04e199326a6c5321bdd7714ee83b4c",
    "commitSignatureVerified": true,
    "outboxPath": "daprdocs/content/en/developing-applications/building-blocks/state-management/howto-outbox.md",
    "outboxSha256": "60487ccdd57bd7c308c0fc87f5e48e7843fb09dbb723e1532a0d8d412672e385",
    "license": "Apache-2.0",
    "licenseSha256": "0b9cab20a5e2ae7e44f40a5ee6b8416f12d2135a547f9fef00e5b61f8d5be99a"
  },
  "rejectedGoSdk": {
    "release": "v1.14.2",
    "releaseCommit": "2523198e6f8a3ab51fd5677dd525710060d4ce2d",
    "releaseCommitSignatureVerified": true,
    "mainCommit": "9d745f76a46e23ae0d022f1abb098132b67603e5",
    "mainCommitSignatureVerified": true,
    "reason": "reachable GO-2026-6061 and GO-2026-5242 in v1.14.2; signed main still pins grpc v1.80.0",
    "admission": "REJECTED_UNTIL_OFFICIAL_FIX_AND_REVALIDATION"
  },
  "deliverySemantics": "AT_LEAST_ONCE",
  "automaticProjectSelection": false,
  "productionAuthorization": false
}
````

### FILE: `dapr_official_transactional_outbox/source-approval.template.json`
```yaml
block_id: "DAPR-OFFICIAL-OUTBOX:03:v1"
operation: CREATE
provenance: AUTHORED
source: "local admission/acquisition code constrained by exact Dapr upstream evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "b1cf530cfd8c7236db1e38de8ea05c3b2ec36750c0620ab994877080cdb1a62d"
variables: []
secrets_allowed: false
```
````json
{
  "schemaVersion": 1,
  "approved": false,
  "approvedBy": "",
  "approvedAt": "",
  "sourceLockSha256": "",
  "sourceId": "dapr-runtime-v1.18.3",
  "licenseAccepted": "Apache-2.0",
  "networkAuthorized": false
}
````

### FILE: `dapr_official_transactional_outbox/project-selection.template.json`
```yaml
block_id: "DAPR-OFFICIAL-OUTBOX:04:v1"
operation: CREATE
provenance: AUTHORED
source: "local admission/acquisition code constrained by exact Dapr upstream evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "c58a95fbf645c06bf1c66b7206880c9b59d3f0e65c727628b6e1b5fe934f095c"
variables: []
secrets_allowed: false
```
````json
{
  "schemaVersion": 1,
  "approved": false,
  "approvedBy": "",
  "approvedAt": "",
  "runtimeTag": "v1.18.3",
  "runtimeCommit": "2dcf2e3548faf3c740b3d29ee300e2132db65cff",
  "stateStoreComponent": "",
  "pubsubComponent": "",
  "outboxPublishTopic": "",
  "stateStoreTransactionsVerified": false,
  "stateStoreOutboxSupportVerified": false,
  "brokerAccountAccessVerified": false,
  "deliverySemanticsAccepted": "AT_LEAST_ONCE",
  "consumerIdempotencyVerified": false,
  "duplicateDeliveryEvidence": "",
  "brokerRecoveryEvidence": "",
  "rollbackEvidence": "",
  "goSdkSelected": false,
  "productionAuthorization": false
}
````

### FILE: `dapr_official_transactional_outbox/acquire_official_runtime.ps1`
```yaml
block_id: "DAPR-OFFICIAL-OUTBOX:05:v1"
operation: CREATE
provenance: AUTHORED
source: "local admission/acquisition code constrained by exact Dapr upstream evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "44bd7536962863c852de3b0eec96b807b0f2385485575730c1accb439a49db33"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0
[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)] [string]$LockPath,
    [Parameter(Mandatory = $true)] [string]$ApprovalPath,
    [Parameter(Mandatory = $true)] [string]$Destination,
    [switch]$AllowNetwork,
    [string]$ApprovedCacheArchive
)

$ErrorActionPreference = 'Stop'

function Fail([string]$Message) { throw "DAPR_RUNTIME_ACQUISITION_REJECTED: $Message" }

foreach ($path in @($LockPath, $ApprovalPath)) {
    if (-not (Test-Path -LiteralPath $path -PathType Leaf)) { Fail "missing file: $path" }
}
if (Test-Path -LiteralPath $Destination) { Fail 'destination already exists' }

$lock = Get-Content -LiteralPath $LockPath -Raw | ConvertFrom-Json -Depth 30
$approval = Get-Content -LiteralPath $ApprovalPath -Raw | ConvertFrom-Json -Depth 20
$lockHash = (Get-FileHash -LiteralPath $LockPath -Algorithm SHA256).Hash.ToLowerInvariant()

if ($approval.approved -ne $true) { Fail 'approval is false' }
if ([string]::IsNullOrWhiteSpace([string]$approval.approvedBy)) { Fail 'approvedBy is empty' }
$approvedAt = [datetimeoffset]::MinValue
if (-not [datetimeoffset]::TryParse([string]$approval.approvedAt, [ref]$approvedAt)) { Fail 'approvedAt is invalid' }
if ([string]$approval.sourceLockSha256 -cne $lockHash) { Fail 'approval is not bound to this source lock' }
if ([string]$approval.sourceId -cne 'dapr-runtime-v1.18.3') { Fail 'sourceId mismatch' }
if ([string]$approval.licenseAccepted -cne 'Apache-2.0') { Fail 'Apache-2.0 was not accepted' }

$usingCache = -not [string]::IsNullOrWhiteSpace($ApprovedCacheArchive)
if (-not $usingCache -and (-not $AllowNetwork -or $approval.networkAuthorized -ne $true)) {
    Fail 'network acquisition requires both -AllowNetwork and networkAuthorized=true'
}
if ($usingCache -and -not (Test-Path -LiteralPath $ApprovedCacheArchive -PathType Leaf)) {
    Fail 'approved cache archive is missing'
}

$stage = Join-Path ([System.IO.Path]::GetTempPath()) ("dapr-runtime-acquire-" + [guid]::NewGuid().ToString('N'))
$archive = Join-Path $stage 'source.zip'
$expanded = Join-Path $stage 'expanded'
try {
    New-Item -ItemType Directory -Path $stage, $expanded | Out-Null
    if ($usingCache) {
        Copy-Item -LiteralPath $ApprovedCacheArchive -Destination $archive
    } else {
        Invoke-WebRequest -UseBasicParsing -Uri ([string]$lock.runtime.archiveUrl) -OutFile $archive
    }
    $item = Get-Item -LiteralPath $archive
    if ($item.Length -ne [long]$lock.runtime.archiveBytes) { Fail 'archive size mismatch' }
    $archiveHash = (Get-FileHash -LiteralPath $archive -Algorithm SHA256).Hash.ToLowerInvariant()
    if ($archiveHash -cne [string]$lock.runtime.archiveSha256) { Fail 'archive SHA-256 mismatch' }

    Expand-Archive -LiteralPath $archive -DestinationPath $expanded
    $roots = @(Get-ChildItem -LiteralPath $expanded -Directory)
    if ($roots.Count -ne 1 -or $roots[0].Name -cne ("dapr-" + [string]$lock.runtime.commit)) {
        Fail 'archive root does not match the locked commit'
    }
    $sourceRoot = $roots[0].FullName
    $checks = [ordered]@{
        'LICENSE' = [string]$lock.runtime.licenseSha256
        'go.mod' = [string]$lock.runtime.goModSha256
        'go.sum' = [string]$lock.runtime.goSumSha256
    }
    foreach ($entry in $checks.GetEnumerator()) {
        $path = Join-Path $sourceRoot $entry.Key
        if (-not (Test-Path -LiteralPath $path -PathType Leaf)) { Fail "archive lacks $($entry.Key)" }
        $actual = (Get-FileHash -LiteralPath $path -Algorithm SHA256).Hash.ToLowerInvariant()
        if ($actual -cne $entry.Value) { Fail "$($entry.Key) SHA-256 mismatch" }
    }
    Move-Item -LiteralPath $sourceRoot -Destination $Destination
    "DAPR_RUNTIME_ACQUISITION_PASS destination=$Destination commit=$($lock.runtime.commit)"
} finally {
    if (Test-Path -LiteralPath $stage) { Remove-Item -LiteralPath $stage -Recurse -Force }
}
````

### FILE: `dapr_official_transactional_outbox/validate_project_selection.ps1`
```yaml
block_id: "DAPR-OFFICIAL-OUTBOX:06:v1"
operation: CREATE
provenance: AUTHORED
source: "local admission/acquisition code constrained by exact Dapr upstream evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "98fc98f6e93a549e3f4ab2590b1a4c2a6857f631bdf0f6d252cdab32e799a3cd"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0
[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)]
    [string]$SelectionPath
)

$ErrorActionPreference = 'Stop'

function Fail([string]$Message) {
    throw "DAPR_OUTBOX_SELECTION_REJECTED: $Message"
}

if (-not (Test-Path -LiteralPath $SelectionPath -PathType Leaf)) {
    Fail 'selection file is missing'
}

try {
    $selection = Get-Content -LiteralPath $SelectionPath -Raw | ConvertFrom-Json -Depth 20
} catch {
    Fail "selection JSON is invalid: $($_.Exception.Message)"
}

$expected = [ordered]@{
    schemaVersion = 1
    runtimeTag = 'v1.18.3'
    runtimeCommit = '2dcf2e3548faf3c740b3d29ee300e2132db65cff'
    deliverySemanticsAccepted = 'AT_LEAST_ONCE'
    goSdkSelected = $false
}

foreach ($entry in $expected.GetEnumerator()) {
    if ($selection.PSObject.Properties.Name -notcontains $entry.Key) {
        Fail "missing property $($entry.Key)"
    }
    if ($selection.($entry.Key) -ne $entry.Value) {
        Fail "$($entry.Key) must equal $($entry.Value)"
    }
}

foreach ($name in @('approved', 'stateStoreTransactionsVerified', 'stateStoreOutboxSupportVerified', 'brokerAccountAccessVerified', 'consumerIdempotencyVerified', 'productionAuthorization')) {
    if ($selection.PSObject.Properties.Name -notcontains $name -or $selection.$name -ne $true) {
        Fail "$name must be true with project evidence"
    }
}

foreach ($name in @('approvedBy', 'approvedAt', 'stateStoreComponent', 'pubsubComponent', 'outboxPublishTopic', 'duplicateDeliveryEvidence', 'brokerRecoveryEvidence', 'rollbackEvidence')) {
    if ($selection.PSObject.Properties.Name -notcontains $name -or [string]::IsNullOrWhiteSpace([string]$selection.$name)) {
        Fail "$name must be non-empty"
    }
}

$approvedAt = [datetimeoffset]::MinValue
if (-not [datetimeoffset]::TryParse([string]$selection.approvedAt, [ref]$approvedAt)) {
    Fail 'approvedAt must be a valid timestamp'
}

$serialized = $selection | ConvertTo-Json -Depth 20
if ($serialized -match '(?i)secret|password|token|private[_-]?key') {
    Fail 'selection must contain evidence references, never secrets'
}

"DAPR_OUTBOX_SELECTION_PASS runtime=$($selection.runtimeTag) store=$($selection.stateStoreComponent) pubsub=$($selection.pubsubComponent)"
````

### FILE: `dapr_official_transactional_outbox/test_embedded_sources.ps1`
```yaml
block_id: "DAPR-OFFICIAL-OUTBOX:07:v1"
operation: CREATE
provenance: AUTHORED
source: "local admission/acquisition code constrained by exact Dapr upstream evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "b01ec736d11613158cde55b1fb9175e7ad6643533cbfaeb28ef6eae6bbd60afe"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0
[CmdletBinding()]
param()

$ErrorActionPreference = 'Stop'
$root = Split-Path -Parent $PSCommandPath

function Assert([bool]$Condition, [string]$Message) {
    if (-not $Condition) { throw "DAPR_OUTBOX_SOURCE_TEST_FAILED: $Message" }
}

$expected = [ordered]@{
    'upstream/docs/LICENSE' = '0b9cab20a5e2ae7e44f40a5ee6b8416f12d2135a547f9fef00e5b61f8d5be99a'
    'upstream/docs/howto-outbox.md' = '60487ccdd57bd7c308c0fc87f5e48e7843fb09dbb723e1532a0d8d412672e385'
    'upstream/runtime/LICENSE' = '9523eacf5421b65637420755c00b5175f7804b6a6eb736b86cd10ea6b3937dcf'
    'upstream/runtime/pkg/outbox/outbox.go' = 'a0d48867049b7404e51852c8d688d7d9ca918cf2eecd9fc8a686636c783a985c'
    'upstream/runtime/pkg/outbox/fake/fake.go' = 'bf653d08463f8ca3f8a10658cf4bb1d661e7f244ef3cd8852574932e9e0545bd'
    'upstream/runtime/pkg/outbox/fake/fake_test.go' = 'cac9dfe14be633bdf1034c7b24ca133415f57a8c17338023708a8ab6ba1a4f61'
    'upstream/runtime/pkg/runtime/pubsub/outbox.go' = '8a0a060481e107f3f0617fc1d6baf86c602105d3ef29c429a83abbc0c04923fc'
    'upstream/runtime/pkg/runtime/pubsub/outbox_test.go' = 'fac77524bbbbe8e2a3449ba7ec062f8979a1760085cf56e598acd8670cdd3cc4'
    'upstream/runtime/tests/integration/suite/daprd/outbox/http/basic.go' = 'ce9ab635518713343b8e3db3fe98b12357a55086102ae2343fe80d1c34078bfe'
    'upstream/runtime/tests/integration/suite/daprd/outbox/http/delete.go' = '222dd239f765a04accac212563e56b200884cb77246816388ec80521efd2dbfb'
    'upstream/runtime/tests/integration/suite/daprd/outbox/http/projection.go' = '4e49a51289c0a2f809f496ef7a78e72e963a3d2b5211cbcdcfd246101392c216'
    'upstream/runtime/tests/integration/suite/daprd/outbox/grpc/basic.go' = '42e5524ebe243896d893a00a86afe2c9846a8642f6bf6b9c514467613d05eab3'
    'upstream/runtime/tests/integration/suite/daprd/outbox/grpc/projection.go' = '45849b9bc17678967886a3249b861097d702a14d5e99aa6bbb468d055b2d31bf'
}

foreach ($entry in $expected.GetEnumerator()) {
    $path = Join-Path $root $entry.Key
    Assert (Test-Path -LiteralPath $path -PathType Leaf) "missing $($entry.Key)"
    $actual = (Get-FileHash -LiteralPath $path -Algorithm SHA256).Hash.ToLowerInvariant()
    Assert ($actual -ceq $entry.Value) "SHA-256 mismatch for $($entry.Key)"
}

$runtime = Get-Content -LiteralPath (Join-Path $root 'upstream/runtime/pkg/runtime/pubsub/outbox.go') -Raw
$tests = Get-Content -LiteralPath (Join-Path $root 'upstream/runtime/pkg/runtime/pubsub/outbox_test.go') -Raw
$docs = Get-Content -LiteralPath (Join-Path $root 'upstream/docs/howto-outbox.md') -Raw
$lock = Get-Content -LiteralPath (Join-Path $root 'official-source-lock.json') -Raw | ConvertFrom-Json -Depth 20

Assert ($runtime.Contains('outbox.projection')) 'runtime projection implementation marker missing'
Assert ($runtime.Contains('PublishInternal')) 'runtime internal publish marker missing'
Assert ($tests.Contains('TestNewOutbox')) 'official constructor test missing'
Assert ($tests.Contains('TestPublishInternal')) 'official publish test missing'
Assert ($docs.Contains('at-least-once')) 'official delivery semantics missing'
Assert ($lock.runtime.focusedReachableVulnerabilities -eq 0) 'runtime focal scan is not clean'
Assert ($lock.rejectedGoSdk.admission -eq 'REJECTED_UNTIL_OFFICIAL_FIX_AND_REVALIDATION') 'vulnerable Go SDK is not rejected'

$validator = Join-Path $root 'validate_project_selection.ps1'
$template = Join-Path $root 'project-selection.template.json'
$rejected = $false
try {
    & $validator -SelectionPath $template | Out-Null
} catch {
    $rejected = $_.Exception.Message -match 'DAPR_OUTBOX_SELECTION_REJECTED'
}
Assert $rejected 'incomplete template was not rejected'

"DAPR_OUTBOX_SOURCE_TEST_PASS verbatim_files=$($expected.Count) sdk=REJECTED"
````

### FILE: `dapr_official_transactional_outbox/test_acquisition.ps1`
```yaml
block_id: "DAPR-OFFICIAL-OUTBOX:08:v1"
operation: CREATE
provenance: AUTHORED
source: "local admission/acquisition code constrained by exact Dapr upstream evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "b0fe396e7a64eca7d645dacdf326e04114f9e6aabaa64c2b499c2299b2c96d8a"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0
[CmdletBinding()]
param()

$ErrorActionPreference = 'Stop'
$root = Split-Path -Parent $PSCommandPath
$acquirer = Join-Path $root 'acquire_official_runtime.ps1'
$lockPath = Join-Path $root 'official-source-lock.json'
$templatePath = Join-Path $root 'source-approval.template.json'
$scratch = Join-Path ([System.IO.Path]::GetTempPath()) ("dapr-acquisition-test-" + [guid]::NewGuid().ToString('N'))

function Assert-Rejected([scriptblock]$Action, [string]$Pattern) {
    $rejected = $false
    try { & $Action } catch { $rejected = $_.Exception.Message -match $Pattern }
    if (-not $rejected) { throw "DAPR_ACQUISITION_TEST_FAILED: expected rejection matching $Pattern" }
}

try {
    New-Item -ItemType Directory -Path $scratch | Out-Null
    $destination = Join-Path $scratch 'source'

    Assert-Rejected { & $acquirer -LockPath $lockPath -ApprovalPath $templatePath -Destination $destination } 'approval is false'

    $approval = Get-Content -LiteralPath $templatePath -Raw | ConvertFrom-Json -Depth 10
    $approval.approved = $true
    $approval.approvedBy = 'offline-test'
    $approval.approvedAt = '2026-08-28T00:00:00Z'
    $approval.sourceLockSha256 = ('0' * 64)
    $badApproval = Join-Path $scratch 'bad-hash.json'
    $approval | ConvertTo-Json -Depth 10 | Set-Content -LiteralPath $badApproval -Encoding utf8NoBOM
    Assert-Rejected { & $acquirer -LockPath $lockPath -ApprovalPath $badApproval -Destination $destination } 'not bound'

    $approval.sourceLockSha256 = (Get-FileHash -LiteralPath $lockPath -Algorithm SHA256).Hash.ToLowerInvariant()
    $validApproval = Join-Path $scratch 'valid.json'
    $approval | ConvertTo-Json -Depth 10 | Set-Content -LiteralPath $validApproval -Encoding utf8NoBOM
    Assert-Rejected { & $acquirer -LockPath $lockPath -ApprovalPath $validApproval -Destination $destination } 'network acquisition requires'

    New-Item -ItemType Directory -Path $destination | Out-Null
    Assert-Rejected { & $acquirer -LockPath $lockPath -ApprovalPath $validApproval -Destination $destination -AllowNetwork } 'destination already exists'

    'DAPR_RUNTIME_ACQUISITION_TEST_PASS negatives=4'
} finally {
    if (Test-Path -LiteralPath $scratch) { Remove-Item -LiteralPath $scratch -Recurse -Force }
}
````

### FILE: `dapr_official_transactional_outbox/upstream/docs/LICENSE`
```yaml
block_id: "DAPR-OFFICIAL-OUTBOX:09:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/dapr/docs/5958a7e19a04e199326a6c5321bdd7714ee83b4c/LICENSE"
license: "Apache-2.0"
sha256: "0b9cab20a5e2ae7e44f40a5ee6b8416f12d2135a547f9fef00e5b61f8d5be99a"
variables: []
secrets_allowed: false
```
````text

                                 Apache License
                           Version 2.0, January 2004
                        http://www.apache.org/licenses/

   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION

   1. Definitions.

      "License" shall mean the terms and conditions for use, reproduction,
      and distribution as defined by Sections 1 through 9 of this document.

      "Licensor" shall mean the copyright owner or entity authorized by
      the copyright owner that is granting the License.

      "Legal Entity" shall mean the union of the acting entity and all
      other entities that control, are controlled by, or are under common
      control with that entity. For the purposes of this definition,
      "control" means (i) the power, direct or indirect, to cause the
      direction or management of such entity, whether by contract or
      otherwise, or (ii) ownership of fifty percent (50%) or more of the
      outstanding shares, or (iii) beneficial ownership of such entity.

      "You" (or "Your") shall mean an individual or Legal Entity
      exercising permissions granted by this License.

      "Source" form shall mean the preferred form for making modifications,
      including but not limited to software source code, documentation
      source, and configuration files.

      "Object" form shall mean any form resulting from mechanical
      transformation or translation of a Source form, including but
      not limited to compiled object code, generated documentation,
      and conversions to other media types.

      "Work" shall mean the work of authorship, whether in Source or
      Object form, made available under the License, as indicated by a
      copyright notice that is included in or attached to the work
      (an example is provided in the Appendix below).

      "Derivative Works" shall mean any work, whether in Source or Object
      form, that is based on (or derived from) the Work and for which the
      editorial revisions, annotations, elaborations, or other modifications
      represent, as a whole, an original work of authorship. For the purposes
      of this License, Derivative Works shall not include works that remain
      separable from, or merely link (or bind by name) to the interfaces of,
      the Work and Derivative Works thereof.

      "Contribution" shall mean any work of authorship, including
      the original version of the Work and any modifications or additions
      to that Work or Derivative Works thereof, that is intentionally
      submitted to Licensor for inclusion in the Work by the copyright owner
      or by an individual or Legal Entity authorized to submit on behalf of
      the copyright owner. For the purposes of this definition, "submitted"
      means any form of electronic, verbal, or written communication sent
      to the Licensor or its representatives, including but not limited to
      communication on electronic mailing lists, source code control systems,
      and issue tracking systems that are managed by, or on behalf of, the
      Licensor for the purpose of discussing and improving the Work, but
      excluding communication that is conspicuously marked or otherwise
      designated in writing by the copyright owner as "Not a Contribution."

      "Contributor" shall mean Licensor and any individual or Legal Entity
      on behalf of whom a Contribution has been received by Licensor and
      subsequently incorporated within the Work.

   2. Grant of Copyright License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      copyright license to reproduce, prepare Derivative Works of,
      publicly display, publicly perform, sublicense, and distribute the
      Work and such Derivative Works in Source or Object form.

   3. Grant of Patent License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      (except as stated in this section) patent license to make, have made,
      use, offer to sell, sell, import, and otherwise transfer the Work,
      where such license applies only to those patent claims licensable
      by such Contributor that are necessarily infringed by their
      Contribution(s) alone or by combination of their Contribution(s)
      with the Work to which such Contribution(s) was submitted. If You
      institute patent litigation against any entity (including a
      cross-claim or counterclaim in a lawsuit) alleging that the Work
      or a Contribution incorporated within the Work constitutes direct
      or contributory patent infringement, then any patent licenses
      granted to You under this License for that Work shall terminate
      as of the date such litigation is filed.

   4. Redistribution. You may reproduce and distribute copies of the
      Work or Derivative Works thereof in any medium, with or without
      modifications, and in Source or Object form, provided that You
      meet the following conditions:

      (a) You must give any other recipients of the Work or
          Derivative Works a copy of this License; and

      (b) You must cause any modified files to carry prominent notices
          stating that You changed the files; and

      (c) You must retain, in the Source form of any Derivative Works
          that You distribute, all copyright, patent, trademark, and
          attribution notices from the Source form of the Work,
          excluding those notices that do not pertain to any part of
          the Derivative Works; and

      (d) If the Work includes a "NOTICE" text file as part of its
          distribution, then any Derivative Works that You distribute must
          include a readable copy of the attribution notices contained
          within such NOTICE file, excluding those notices that do not
          pertain to any part of the Derivative Works, in at least one
          of the following places: within a NOTICE text file distributed
          as part of the Derivative Works; within the Source form or
          documentation, if provided along with the Derivative Works; or,
          within a display generated by the Derivative Works, if and
          wherever such third-party notices normally appear. The contents
          of the NOTICE file are for informational purposes only and
          do not modify the License. You may add Your own attribution
          notices within Derivative Works that You distribute, alongside
          or as an addendum to the NOTICE text from the Work, provided
          that such additional attribution notices cannot be construed
          as modifying the License.

      You may add Your own copyright statement to Your modifications and
      may provide additional or different license terms and conditions
      for use, reproduction, or distribution of Your modifications, or
      for any such Derivative Works as a whole, provided Your use,
      reproduction, and distribution of the Work otherwise complies with
      the conditions stated in this License.

   5. Submission of Contributions. Unless You explicitly state otherwise,
      any Contribution intentionally submitted for inclusion in the Work
      by You to the Licensor shall be under the terms and conditions of
      this License, without any additional terms or conditions.
      Notwithstanding the above, nothing herein shall supersede or modify
      the terms of any separate license agreement you may have executed
      with Licensor regarding such Contributions.

   6. Trademarks. This License does not grant permission to use the trade
      names, trademarks, service marks, or product names of the Licensor,
      except as required for reasonable and customary use in describing the
      origin of the Work and reproducing the content of the NOTICE file.

   7. Disclaimer of Warranty. Unless required by applicable law or
      agreed to in writing, Licensor provides the Work (and each
      Contributor provides its Contributions) on an "AS IS" BASIS,
      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
      implied, including, without limitation, any warranties or conditions
      of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
      PARTICULAR PURPOSE. You are solely responsible for determining the
      appropriateness of using or redistributing the Work and assume any
      risks associated with Your exercise of permissions under this License.

   8. Limitation of Liability. In no event and under no legal theory,
      whether in tort (including negligence), contract, or otherwise,
      unless required by applicable law (such as deliberate and grossly
      negligent acts) or agreed to in writing, shall any Contributor be
      liable to You for damages, including any direct, indirect, special,
      incidental, or consequential damages of any character arising as a
      result of this License or out of the use or inability to use the
      Work (including but not limited to damages for loss of goodwill,
      work stoppage, computer failure or malfunction, or any and all
      other commercial damages or losses), even if such Contributor
      has been advised of the possibility of such damages.

   9. Accepting Warranty or Additional Liability. While redistributing
      the Work or Derivative Works thereof, You may choose to offer,
      and charge a fee for, acceptance of support, warranty, indemnity,
      or other liability obligations and/or rights consistent with this
      License. However, in accepting such obligations, You may act only
      on Your own behalf and on Your sole responsibility, not on behalf
      of any other Contributor, and only if You agree to indemnify,
      defend, and hold each Contributor harmless for any liability
      incurred by, or claims asserted against, such Contributor by reason
      of your accepting any such warranty or additional liability.

   END OF TERMS AND CONDITIONS

   APPENDIX: How to apply the Apache License to your work.

      To apply the Apache License to your work, attach the following
      boilerplate notice, with the fields enclosed by brackets "[]"
      replaced with your own identifying information. (Don't include
      the brackets!)  The text should be enclosed in the appropriate
      comment syntax for the file format. We also recommend that a
      file or class name and description of purpose be included on the
      same "printed page" as the copyright notice for easier
      identification within third-party archives.

   Copyright 2025 The Dapr Authors.

   and others that have contributed code to the public domain.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
````

### FILE: `dapr_official_transactional_outbox/upstream/docs/howto-outbox.md`
```yaml
block_id: "DAPR-OFFICIAL-OUTBOX:10:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/dapr/docs/5958a7e19a04e199326a6c5321bdd7714ee83b4c/daprdocs/content/en/developing-applications/building-blocks/state-management/howto-outbox.md"
license: "Apache-2.0"
sha256: "60487ccdd57bd7c308c0fc87f5e48e7843fb09dbb723e1532a0d8d412672e385"
variables: []
secrets_allowed: false
```
````markdown
---
type: docs
title: "How-To: Enable the transactional outbox pattern"
linkTitle: "How-To: Enable the transactional outbox pattern"
weight: 400
description: "Commit a single transaction across a state store and pub/sub message broker"
---

The transactional outbox pattern is a well known design pattern for sending notifications regarding changes in an application's state. The transactional outbox pattern uses a single transaction that spans across the database and the message broker delivering the notification. 

Developers are faced with many difficult technical challenges when trying to implement this pattern on their own, which often involves writing error-prone central coordination managers that, at most, support a combination of one or two databases and message brokers.

For example, you can use the outbox pattern to:  
1. Write a new user record to an account database. 
1. Send a notification message that the account was successfully created. 

With Dapr's outbox support, you can notify subscribers when an application's state is created or updated when calling Dapr's [transactions API]({{% ref "state_api.md#state-transactions" %}}).

The diagram below is an overview of how the outbox feature works at a high level:

1) Service A saves/updates state to the state store using a transaction.
2) A message is written to the broker under the same transaction. When the message is successfully delivered to the message broker, the transaction completes, ensuring the state and message are transacted together.
3) The message broker delivers the message topic to any subscribers - in this case, Service B.

<img src="/images/state-management-outbox.png" width=800 alt="Diagram showing the overview of outbox pattern">

## How outbox works under the hood

Dapr outbox processes requests in two flows: the user request flow and the background message flow. Together, they guarantee that state and events stay consistent.

<img src="/images/state-management-outbox-steps.png" width=800 alt="Diagram showing the steps of the outbox pattern">

This is the sequence of interactions:

1. An application calls the Dapr State Management API to write state transactionally using the transactional methods.  
   This is the entry point where business data, such as an order or profile update, is submitted for persistence.

2. Dapr publishes an intent message with a unique transaction ID to an internal outbox topic.  
   This durable record ensures the event intent exists before any database commit happens.

3. The state and a transaction marker are written atomically in the same state store.  
   Both the business data and the marker are committed in the same transaction, preventing partial writes.

4. The application receives a success response after the transaction commits.  
   At this point, the application can continue, knowing state is saved and the event intent is guaranteed.

5. A background subscriber reads the intent message.  
   When outbox is enabled, Dapr starts consumers that process the internal outbox topic.

6. The subscriber verifies the transaction marker in the state store.  
   This check confirms that the database commit was successful before publishing externally.

7. Verified business event is published to the external pub/sub topic.  
   The event is sent to the configured broker (Kafka, RabbitMQ, etc.) where other services can consume it.

8. The marker is cleaned up (deleted) from the state store.  
   This prevents unbounded growth in the database once the event has been successfully delivered.

9. Message is acknowledged and removed from internal topic  
   If publishing or cleanup fails, Dapr retries, ensuring reliable at-least-once delivery.
  
## Requirements

1. The outbox feature requires a [transactional state store]({{% ref supported-state-stores %}}) supported by Dapr.  
   [Learn more about the transactional methods you can use.]({{% ref "howto-get-save-state.md#perform-state-transactions" %}})

2. Any [pub/sub broker]({{% ref supported-pubsub %}}) supported by Dapr can be used with the outbox feature.

   {{% alert title="Note" color="primary" %}}
   Message brokers that support the competing consumer pattern (for example, [Apache Kafka]({{% ref setup-apache-kafka%}})) are recommended to reduce the chance of duplicate events.
   {{% /alert %}}

3. Internal outbox topic  
   When outbox is enabled, Dapr creates an internal topic using the following naming convention: `{namespace}{appID}{topic}outbox`, where:

   - `namespace`: the Dapr application namespace (if configured)
   - `appID`: the Dapr application identifier
   - `topic`: the value specified in the `outboxPublishTopic` metadata

   This way each outbox topic is uniquely identified per application and external topic, preventing routing conflicts in multi-tenant environments.

   You can override the auto-generated name by setting the `outboxInternalTopic` metadata field. When set, this value is used as the complete internal topic name with no namespace, appID, or suffix appended.

   {{% alert title="Note" color="primary" %}}
   Ensure that the topic is created in advance, or Dapr has sufficient permissions to create the topic at startup time.
   {{% /alert %}}

## Enable the outbox pattern

To enable the outbox feature, add the following required and optional fields on a state store component:

```yaml
apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
  name: mysql-outbox
spec:
  type: state.mysql
  version: v1
  metadata:
  - name: connectionString
    value: "<CONNECTION STRING>"
  - name: outboxPublishPubsub # Required
    value: "mypubsub"
  - name: outboxPublishTopic # Required
    value: "newOrder"
  - name: outboxPubsub # Optional
    value: "myOutboxPubsub"
  - name: outboxInternalTopic # Optional
    value: "myapp-neworder-outbox"
  - name: outboxDiscardWhenMissingState #Optional. Defaults to false
    value: false
```

### Metadata fields

| Name                | Required    | Default Value | Description                                            |
| --------------------|-------------|---------------|------------------------------------------------------- |
| outboxPublishPubsub | Yes         | N/A           | Sets the name of the pub/sub component to deliver the notifications when publishing state changes
| outboxPublishTopic  | Yes         | N/A           | Sets the topic that receives the state changes on the pub/sub configured with `outboxPublishPubsub`. The message body will be a state transaction item for an `insert` or `update` operation
| outboxPubsub        | No          | `outboxPublishPubsub`           | Sets the pub/sub component used by Dapr to coordinate the state and pub/sub transactions. If not set, the pub/sub component configured with `outboxPublishPubsub` is used. This is useful if you want to separate the pub/sub component used to send the notification state changes from the one used to coordinate the transaction
| outboxInternalTopic | No          | Auto-generated | Sets the internal outbox topic name. When set, this value is used as the complete topic name with no namespace, appID, or suffix appended. When empty, the default naming convention `{namespace}{appID}{topic}outbox` is used
| outboxDiscardWhenMissingState  | No         | `false`           | By setting `outboxDiscardWhenMissingState` to `true`, Dapr discards the transaction if it cannot find the state in the database and does not retry. This setting can be useful if the state store data has been deleted for any reason before Dapr was able to deliver the message and you would like Dapr to drop the items from the pub/sub and stop retrying to fetch the state

## Additional configurations

### Combining outbox and non-outbox messages on the same state store

If you want to use the same state store for sending both outbox and non-outbox messages, simply define two state store components that connect to the same state store, where one has the outbox feature and the other does not.

#### MySQL state store without outbox

```yaml
apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
  name: mysql
spec:
  type: state.mysql
  version: v1
  metadata:
  - name: connectionString
    value: "<CONNECTION STRING>"
```

#### MySQL state store with outbox

```yaml
apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
  name: mysql-outbox
spec:
  type: state.mysql
  version: v1
  metadata:
  - name: connectionString
    value: "<CONNECTION STRING>"
  - name: outboxPublishPubsub # Required
    value: "mypubsub"
  - name: outboxPublishTopic # Required
    value: "newOrder"
```

### Shape the outbox pattern message

You can override the outbox pattern message published to the pub/sub broker by setting another transaction that is not be saved to the database and is explicitly mentioned as a projection. This transaction is added a metadata key named `outbox.projection` with a value set to `true`. When added to the state array saved in a transaction, this payload is ignored when the state is written and the data is used as the payload sent to the upstream subscriber.

To use correctly, the `key` values must match between the operation on the state store and the message projection. If the keys do not match, the whole transaction fails.

If you have two or more `outbox.projection` enabled state items for the same key, the first one defined is used and the others are ignored. 

[Learn more about default and custom CloudEvent messages.]({{% ref pubsub-cloudevents.md %}})

{{< tabpane text=true >}}

{{% tab "Python" %}}

<!--python-->

In the following Python SDK example of a state transaction, the value of `"2"` is saved to the database, but the value of `"3"` is published to the end-user topic.

```python
DAPR_STORE_NAME = "statestore"

async def main():
    client = DaprClient()

    client.execute_state_transaction(
       store_name=DAPR_STORE_NAME,
       operations=[
          # Define the first state operation to save the value "2"
          TransactionalStateOperation(
             key='key1', data='2', metadata={'outbox.projection': 'false'}
          ),
          # Define the second state operation to publish the value "3" with metadata
          TransactionalStateOperation(
             key='key1', data='3', metadata={'outbox.projection': 'true'}
          ),
       ],
    )

    print("State transaction executed.")
```

By setting the metadata item `"outbox.projection"` to `"true"` and making sure the `key` values match (`key1`):
- The first operation is written to the state store and no message is written to the message broker.
- The second operation value is published to the configured pub/sub topic.   

{{% /tab %}}

{{% tab "JavaScript" %}}

<!--javascript-->

In the following JavaScript SDK example of a state transaction, the value of `"2"` is saved to the database, but the value of `"3"` is published to the end-user topic.

```javascript
const { DaprClient, StateOperationType } = require('@dapr/dapr');

const DAPR_STORE_NAME = "statestore";

async function main() {
  const client = new DaprClient();

  // Define the first state operation to save the value "2"
  const op1 = {
    operation: StateOperationType.UPSERT,
    request: {
      key: "key1",
      value: "2"
    }
  };

  // Define the second state operation to publish the value "3" with metadata
  const op2 = {
    operation: StateOperationType.UPSERT,
    request: {
      key: "key1",
      value: "3",
      metadata: {
        "outbox.projection": "true"
      }
    }
  };

  // Create the list of state operations
  const ops = [op1, op2];

  // Execute the state transaction
  await client.state.transaction(DAPR_STORE_NAME, ops);
  console.log("State transaction executed.");
}

main().catch(err => {
  console.error(err);
});
```

By setting the metadata item `"outbox.projection"` to `"true"` and making sure the `key` values match (`key1`):
- The first operation is written to the state store and no message is written to the message broker.
- The second operation value is published to the configured pub/sub topic.   


{{% /tab %}}

{{% tab ".NET" %}}

<!--dotnet-->

In the following .NET SDK example of a state transaction, the value of `"2"` is saved to the database, but the value of `"3"` is published to the end-user topic.

```csharp
public class Program
{
    private const string DAPR_STORE_NAME = "statestore";

    public static async Task Main(string[] args)
    {
        var client = new DaprClientBuilder().Build();

        // Define the first state operation to save the value "2"
        var op1 = new StateTransactionRequest(
            key: "key1",
            value: Encoding.UTF8.GetBytes("2"),
            operationType: StateOperationType.Upsert
        );

        // Define the second state operation to publish the value "3" with metadata
        var metadata = new Dictionary<string, string>
        {
            { "outbox.projection", "true" }
        };
        var op2 = new StateTransactionRequest(
            key: "key1",
            value: Encoding.UTF8.GetBytes("3"),
            operationType: StateOperationType.Upsert,
            metadata: metadata
        );

        // Create the list of state operations
        var ops = new List<StateTransactionRequest> { op1, op2 };

        // Execute the state transaction
        await client.ExecuteStateTransactionAsync(DAPR_STORE_NAME, ops);
        Console.WriteLine("State transaction executed.");
    }
}
```

By setting the metadata item `"outbox.projection"` to `"true"` and making sure the `key` values match (`key1`):
- The first operation is written to the state store and no message is written to the message broker.
- The second operation value is published to the configured pub/sub topic.    

{{% /tab %}}

{{% tab "Java" %}}

<!--java-->

In the following Java SDK example of a state transaction, the value of `"2"` is saved to the database, but the value of `"3"` is published to the end-user topic.

```java
public class Main {
    private static final String DAPR_STORE_NAME = "statestore";

    public static void main(String[] args) {
        try (DaprClient client = new DaprClientBuilder().build()) {
            // Define the first state operation to save the value "2"
            State<String> state1 = new State<>(
                    "key1",
                    "2",
                    null, // etag
                    null // concurrency and consistency options
            );

            // Define the second state operation to publish the value "3" with metadata
            Map<String, String> metadata = new HashMap<>();
            metadata.put("outbox.projection", "true");

            State<String> state2 = new State<>(
                    "key1",
                    "3",
                    null, // etag
                    metadata, 
                    null // concurrency and consistency options
            );
            
            TransactionalStateOperation<String> op1 = new TransactionalStateOperation<>(
                TransactionalStateOperation.OperationType.UPSERT, state1
            );

            TransactionalStateOperation<String> op2 = new TransactionalStateOperation<>(
                TransactionalStateOperation.OperationType.UPSERT, state2
            );

            // Create the list of transaction state operations
            List<TransactionalStateOperation<?>> ops = new ArrayList<>();
            ops.add(op1);
            ops.add(op2);

            // Configure transaction request setting the state store
            ExecuteStateTransactionRequest transactionRequest = new ExecuteStateTransactionRequest(DAPR_STORE_NAME);
            
            transactionRequest.setOperations(ops);

            // Execute the state transaction
            client.executeStateTransaction(transactionRequest).block();
            System.out.println("State transaction executed.");
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}
```

By setting the metadata item `"outbox.projection"` to `"true"` and making sure the `key` values match (`key1`):
- The first operation is written to the state store and no message is written to the message broker.
- The second operation value is published to the configured pub/sub topic.   


{{% /tab %}}

{{% tab "Go" %}}

<!--go-->

In the following Go SDK example of a state transaction, the value of `"2"` is saved to the database, but the value of `"3"` is published to the end-user topic.

```go
ops := make([]*dapr.StateOperation, 0)

op1 := &dapr.StateOperation{
    Type: dapr.StateOperationTypeUpsert,
    Item: &dapr.SetStateItem{
        Key:   "key1",
        Value: []byte("2"),
    },
}
op2 := &dapr.StateOperation{
    Type: dapr.StateOperationTypeUpsert,
    Item: &dapr.SetStateItem{
        Key:   "key1",
				Value: []byte("3"),
         // Override the data payload saved to the database 
				Metadata: map[string]string{
					"outbox.projection": "true",
        },
    },
}
ops = append(ops, op1, op2)
meta := map[string]string{}
err := testClient.ExecuteStateTransaction(ctx, store, meta, ops)
```

By setting the metadata item `"outbox.projection"` to `"true"` and making sure the `key` values match (`key1`):
- The first operation is written to the state store and no message is written to the message broker.
- The second operation value is published to the configured pub/sub topic.   

{{% /tab %}}

{{% tab "HTTP" %}}

<!--http-->

You can pass the message override using the following HTTP request:

```bash
curl -X POST http://localhost:3500/v1.0/state/starwars/transaction \
  -H "Content-Type: application/json" \
  -d '{
  "operations": [
    {
      "operation": "upsert",
      "request": {
        "key": "order1",
        "value": {
            "orderId": "7hf8374s",
            "type": "book",
            "name": "The name of the wind"
        }
      }
    },
    {
      "operation": "upsert",
      "request": {
        "key": "order1",
        "value": {
            "orderId": "7hf8374s"
        },
        "metadata": {
           "outbox.projection": "true"
        },
        "contentType": "application/json"
      }
    }
  ]
}'
```

By setting the metadata item `"outbox.projection"` to `"true"` and making sure the `key` values match (`key1`):
- The first operation is written to the state store and no message is written to the message broker.
- The second operation value is published to the configured pub/sub topic.   

{{% /tab %}}

{{< /tabpane >}}

### Override Dapr-generated CloudEvent fields

You can override the [Dapr-generated CloudEvent fields]({{% ref "pubsub-cloudevents.md#dapr-generated-cloudevents-example" %}}) on the published outbox event with custom CloudEvent metadata.

{{< tabpane text=true >}}

{{% tab "Python" %}}

<!--python-->

```python
async def execute_state_transaction():
    async with DaprClient() as client:
        # Define state operations
        ops = []

        op1 = {
            'operation': 'upsert',
            'request': {
                'key': 'key1',
                'value': b'2',  # Convert string to byte array
                'metadata': {
                    'cloudevent.id': 'unique-business-process-id',
                    'cloudevent.source': 'CustomersApp',
                    'cloudevent.type': 'CustomerCreated',
                    'cloudevent.subject': '123',
                    'my-custom-ce-field': 'abc'
                }
            }
        }

        ops.append(op1)

        # Execute state transaction
        store_name = 'your-state-store-name'
        try:
            await client.execute_state_transaction(store_name, ops)
            print('State transaction executed.')
        except Exception as e:
            print('Error executing state transaction:', e)

# Run the async function
if __name__ == "__main__":
    asyncio.run(execute_state_transaction())
```
{{% /tab %}}

{{% tab "JavaScript" %}}

<!--javascript-->

```javascript
const { DaprClient } = require('dapr-client');

async function executeStateTransaction() {
    // Initialize Dapr client
    const daprClient = new DaprClient();

    // Define state operations
    const ops = [];

    const op1 = {
        operationType: 'upsert',
        request: {
            key: 'key1',
            value: Buffer.from('2'),
            metadata: {
                'id': 'unique-business-process-id',
                'source': 'CustomersApp',
                'type': 'CustomerCreated',
                'subject': '123',
                'my-custom-ce-field': 'abc'
            }
        }
    };

    ops.push(op1);

    // Execute state transaction
    const storeName = 'your-state-store-name';
    const metadata = {};
}

executeStateTransaction();
```
{{% /tab %}}

{{% tab ".NET" %}}

<!--csharp-->

```csharp
public class StateOperationExample
{
    public async Task ExecuteStateTransactionAsync()
    {
        var daprClient = new DaprClientBuilder().Build();

        // Define the value "2" as a string and serialize it to a byte array
        var value = "2";
        var valueBytes = JsonSerializer.SerializeToUtf8Bytes(value);

        // Define the first state operation to save the value "2" with metadata
       // Override Cloudevent metadata
        var metadata = new Dictionary<string, string>
        {
            { "cloudevent.id", "unique-business-process-id" },
            { "cloudevent.source", "CustomersApp" },
            { "cloudevent.type", "CustomerCreated" },
            { "cloudevent.subject", "123" },
            { "my-custom-ce-field", "abc" }
        };

        var op1 = new StateTransactionRequest(
            key: "key1",
            value: valueBytes,
            operationType: StateOperationType.Upsert,
            metadata: metadata
        );

        // Create the list of state operations
        var ops = new List<StateTransactionRequest> { op1 };

        // Execute the state transaction
        var storeName = "your-state-store-name";
        await daprClient.ExecuteStateTransactionAsync(storeName, ops);
        Console.WriteLine("State transaction executed.");
    }

    public static async Task Main(string[] args)
    {
        var example = new StateOperationExample();
        await example.ExecuteStateTransactionAsync();
    }
}
```
{{% /tab %}}

{{% tab "Java" %}}

<!--java-->

```java
public class StateOperationExample {

    public static void main(String[] args) {
        executeStateTransaction();
    }

  public static void executeStateTransaction() {
    // Build Dapr client
    try (DaprClient daprClient = new DaprClientBuilder().build()) {

      // Override CloudEvent metadata
      Map<String, String> metadata = new HashMap<>();
      metadata.put("cloudevent.id", "unique-business-process-id");
      metadata.put("cloudevent.source", "CustomersApp");
      metadata.put("cloudevent.type", "CustomerCreated");
      metadata.put("cloudevent.subject", "123");
      metadata.put("my-custom-ce-field", "abc");

      State<String> state = new State<>(
          "key1", // Define the key "key1"
          "value1", // Define the value "value1"
          null, // etag
          metadata,
          null // concurrency and consistency options
      );

      // Define state operations
      List<TransactionalStateOperation<?>> ops = new ArrayList<>();
      TransactionalStateOperation<String> op1 = new TransactionalStateOperation<>(
          TransactionalStateOperation.OperationType.UPSERT,
          state
      );
      ops.add(op1);

      // Execute state transaction
      String storeName = "your-state-store-name";
      daprClient.executeStateTransaction(storeName, ops).block();
      System.out.println("State transaction executed.");
    } catch (Exception e) {
      e.printStackTrace();
    }
  }
}
```
{{% /tab %}}

{{% tab "Go" %}}

<!--go-->

```go
func main() {
	// Create a Dapr client
	client, err := dapr.NewClient()
	if err != nil {
		log.Fatalf("failed to create Dapr client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	store := "your-state-store-name"

	// Define state operations
	ops := make([]*dapr.StateOperation, 0)
	op1 := &dapr.StateOperation{
		Type: dapr.StateOperationTypeUpsert,
		Item: &dapr.SetStateItem{
			Key:   "key1",
			Value: []byte("2"),
			// Override Cloudevent metadata
			Metadata: map[string]string{
				"cloudevent.id":                "unique-business-process-id",
				"cloudevent.source":            "CustomersApp",
				"cloudevent.type":              "CustomerCreated",
				"cloudevent.subject":           "123",
				"my-custom-ce-field":           "abc",
			},
		},
	}
	ops = append(ops, op1)

	// Metadata for the transaction (if any)
	meta := map[string]string{}

	// Execute state transaction
	err = client.ExecuteStateTransaction(ctx, store, meta, ops)
	if err != nil {
		log.Fatalf("failed to execute state transaction: %v", err)
	}

	log.Println("State transaction executed.")
}
```
{{% /tab %}}

{{% tab "HTTP" %}}

<!--http-->

```bash
curl -X POST http://localhost:3500/v1.0/state/starwars/transaction \
  -H "Content-Type: application/json" \
  -d '{
        "operations": [
          {
            "operation": "upsert",
            "request": {
              "key": "key1",
              "value": "2"
            }
          },
        ],
        "metadata": {
          "id": "unique-business-process-id",
          "source": "CustomersApp",
          "type": "CustomerCreated",
          "subject": "123",
          "my-custom-ce-field": "abc",
        }
      }'
```

{{% /tab %}}

{{< /tabpane >}}


{{% alert title="Note" color="primary" %}}
The `data` CloudEvent field is reserved for Dapr's use only, and is non-customizable.

{{% /alert %}}

## Demo

Watch [this video for an overview of the outbox pattern](https://youtu.be/rTovKpG0rhY?t=1338):

{{< youtube id=rTovKpG0rhY start=1338 >}}

## Next Steps

[How Dapr Outbox Eliminates Dual Writes in Distributed Applications](https://www.diagrid.io/blog/how-dapr-outbox-eliminates-dual-writes-in-distributed-applications)
````

### FILE: `dapr_official_transactional_outbox/upstream/runtime/LICENSE`
```yaml
block_id: "DAPR-OFFICIAL-OUTBOX:11:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/dapr/dapr/2dcf2e3548faf3c740b3d29ee300e2132db65cff/LICENSE"
license: "Apache-2.0"
sha256: "9523eacf5421b65637420755c00b5175f7804b6a6eb736b86cd10ea6b3937dcf"
variables: []
secrets_allowed: false
```
````text

                                 Apache License
                           Version 2.0, January 2004
                        http://www.apache.org/licenses/

   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION

   1. Definitions.

      "License" shall mean the terms and conditions for use, reproduction,
      and distribution as defined by Sections 1 through 9 of this document.

      "Licensor" shall mean the copyright owner or entity authorized by
      the copyright owner that is granting the License.

      "Legal Entity" shall mean the union of the acting entity and all
      other entities that control, are controlled by, or are under common
      control with that entity. For the purposes of this definition,
      "control" means (i) the power, direct or indirect, to cause the
      direction or management of such entity, whether by contract or
      otherwise, or (ii) ownership of fifty percent (50%) or more of the
      outstanding shares, or (iii) beneficial ownership of such entity.

      "You" (or "Your") shall mean an individual or Legal Entity
      exercising permissions granted by this License.

      "Source" form shall mean the preferred form for making modifications,
      including but not limited to software source code, documentation
      source, and configuration files.

      "Object" form shall mean any form resulting from mechanical
      transformation or translation of a Source form, including but
      not limited to compiled object code, generated documentation,
      and conversions to other media types.

      "Work" shall mean the work of authorship, whether in Source or
      Object form, made available under the License, as indicated by a
      copyright notice that is included in or attached to the work
      (an example is provided in the Appendix below).

      "Derivative Works" shall mean any work, whether in Source or Object
      form, that is based on (or derived from) the Work and for which the
      editorial revisions, annotations, elaborations, or other modifications
      represent, as a whole, an original work of authorship. For the purposes
      of this License, Derivative Works shall not include works that remain
      separable from, or merely link (or bind by name) to the interfaces of,
      the Work and Derivative Works thereof.

      "Contribution" shall mean any work of authorship, including
      the original version of the Work and any modifications or additions
      to that Work or Derivative Works thereof, that is intentionally
      submitted to Licensor for inclusion in the Work by the copyright owner
      or by an individual or Legal Entity authorized to submit on behalf of
      the copyright owner. For the purposes of this definition, "submitted"
      means any form of electronic, verbal, or written communication sent
      to the Licensor or its representatives, including but not limited to
      communication on electronic mailing lists, source code control systems,
      and issue tracking systems that are managed by, or on behalf of, the
      Licensor for the purpose of discussing and improving the Work, but
      excluding communication that is conspicuously marked or otherwise
      designated in writing by the copyright owner as "Not a Contribution."

      "Contributor" shall mean Licensor and any individual or Legal Entity
      on behalf of whom a Contribution has been received by Licensor and
      subsequently incorporated within the Work.

   2. Grant of Copyright License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      copyright license to reproduce, prepare Derivative Works of,
      publicly display, publicly perform, sublicense, and distribute the
      Work and such Derivative Works in Source or Object form.

   3. Grant of Patent License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      (except as stated in this section) patent license to make, have made,
      use, offer to sell, sell, import, and otherwise transfer the Work,
      where such license applies only to those patent claims licensable
      by such Contributor that are necessarily infringed by their
      Contribution(s) alone or by combination of their Contribution(s)
      with the Work to which such Contribution(s) was submitted. If You
      institute patent litigation against any entity (including a
      cross-claim or counterclaim in a lawsuit) alleging that the Work
      or a Contribution incorporated within the Work constitutes direct
      or contributory patent infringement, then any patent licenses
      granted to You under this License for that Work shall terminate
      as of the date such litigation is filed.

   4. Redistribution. You may reproduce and distribute copies of the
      Work or Derivative Works thereof in any medium, with or without
      modifications, and in Source or Object form, provided that You
      meet the following conditions:

      (a) You must give any other recipients of the Work or
          Derivative Works a copy of this License; and

      (b) You must cause any modified files to carry prominent notices
          stating that You changed the files; and

      (c) You must retain, in the Source form of any Derivative Works
          that You distribute, all copyright, patent, trademark, and
          attribution notices from the Source form of the Work,
          excluding those notices that do not pertain to any part of
          the Derivative Works; and

      (d) If the Work includes a "NOTICE" text file as part of its
          distribution, then any Derivative Works that You distribute must
          include a readable copy of the attribution notices contained
          within such NOTICE file, excluding those notices that do not
          pertain to any part of the Derivative Works, in at least one
          of the following places: within a NOTICE text file distributed
          as part of the Derivative Works; within the Source form or
          documentation, if provided along with the Derivative Works; or,
          within a display generated by the Derivative Works, if and
          wherever such third-party notices normally appear. The contents
          of the NOTICE file are for informational purposes only and
          do not modify the License. You may add Your own attribution
          notices within Derivative Works that You distribute, alongside
          or as an addendum to the NOTICE text from the Work, provided
          that such additional attribution notices cannot be construed
          as modifying the License.

      You may add Your own copyright statement to Your modifications and
      may provide additional or different license terms and conditions
      for use, reproduction, or distribution of Your modifications, or
      for any such Derivative Works as a whole, provided Your use,
      reproduction, and distribution of the Work otherwise complies with
      the conditions stated in this License.

   5. Submission of Contributions. Unless You explicitly state otherwise,
      any Contribution intentionally submitted for inclusion in the Work
      by You to the Licensor shall be under the terms and conditions of
      this License, without any additional terms or conditions.
      Notwithstanding the above, nothing herein shall supersede or modify
      the terms of any separate license agreement you may have executed
      with Licensor regarding such Contributions.

   6. Trademarks. This License does not grant permission to use the trade
      names, trademarks, service marks, or product names of the Licensor,
      except as required for reasonable and customary use in describing the
      origin of the Work and reproducing the content of the NOTICE file.

   7. Disclaimer of Warranty. Unless required by applicable law or
      agreed to in writing, Licensor provides the Work (and each
      Contributor provides its Contributions) on an "AS IS" BASIS,
      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
      implied, including, without limitation, any warranties or conditions
      of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
      PARTICULAR PURPOSE. You are solely responsible for determining the
      appropriateness of using or redistributing the Work and assume any
      risks associated with Your exercise of permissions under this License.

   8. Limitation of Liability. In no event and under no legal theory,
      whether in tort (including negligence), contract, or otherwise,
      unless required by applicable law (such as deliberate and grossly
      negligent acts) or agreed to in writing, shall any Contributor be
      liable to You for damages, including any direct, indirect, special,
      incidental, or consequential damages of any character arising as a
      result of this License or out of the use or inability to use the
      Work (including but not limited to damages for loss of goodwill,
      work stoppage, computer failure or malfunction, or any and all
      other commercial damages or losses), even if such Contributor
      has been advised of the possibility of such damages.

   9. Accepting Warranty or Additional Liability. While redistributing
      the Work or Derivative Works thereof, You may choose to offer,
      and charge a fee for, acceptance of support, warranty, indemnity,
      or other liability obligations and/or rights consistent with this
      License. However, in accepting such obligations, You may act only
      on Your own behalf and on Your sole responsibility, not on behalf
      of any other Contributor, and only if You agree to indemnify,
      defend, and hold each Contributor harmless for any liability
      incurred by, or claims asserted against, such Contributor by reason
      of your accepting any such warranty or additional liability.

   END OF TERMS AND CONDITIONS

   APPENDIX: How to apply the Apache License to your work.

      To apply the Apache License to your work, attach the following
      boilerplate notice, with the fields enclosed by brackets "[]"
      replaced with your own identifying information. (Don't include
      the brackets!)  The text should be enclosed in the appropriate
      comment syntax for the file format. We also recommend that a
      file or class name and description of purpose be included on the
      same "printed page" as the copyright notice for easier
      identification within third-party archives.

   Copyright 2021 The Dapr Authors.

   and others that have contributed code to the public domain.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
````

### FILE: `dapr_official_transactional_outbox/upstream/runtime/pkg/outbox/outbox.go`
```yaml
block_id: "DAPR-OFFICIAL-OUTBOX:12:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/dapr/dapr/2dcf2e3548faf3c740b3d29ee300e2132db65cff/pkg/outbox/outbox.go"
license: "Apache-2.0"
sha256: "a0d48867049b7404e51852c8d688d7d9ca918cf2eecd9fc8a686636c783a985c"
variables: []
secrets_allowed: false
```
````go
/*
Copyright 2023 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package outbox

import (
	"context"

	"github.com/dapr/components-contrib/state"
	"github.com/dapr/dapr/pkg/apis/components/v1alpha1"
)

// Outbox defines the interface for all Outbox pattern operations combining state and pubsub.
type Outbox interface {
	AddOrUpdateOutbox(stateStore v1alpha1.Component)
	Enabled(stateStore string) bool
	PublishInternal(ctx context.Context, stateStore string, states []state.TransactionalStateOperation, source, traceID, traceState string) ([]state.TransactionalStateOperation, error)
	SubscribeToInternalTopics(ctx context.Context, appID string) error
}
````

### FILE: `dapr_official_transactional_outbox/upstream/runtime/pkg/outbox/fake/fake.go`
```yaml
block_id: "DAPR-OFFICIAL-OUTBOX:13:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/dapr/dapr/2dcf2e3548faf3c740b3d29ee300e2132db65cff/pkg/outbox/fake/fake.go"
license: "Apache-2.0"
sha256: "bf653d08463f8ca3f8a10658cf4bb1d661e7f244ef3cd8852574932e9e0545bd"
variables: []
secrets_allowed: false
```
````go
/*
Copyright 2023 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package fake

import (
	"context"

	"github.com/dapr/components-contrib/state"
	"github.com/dapr/dapr/pkg/apis/components/v1alpha1"
)

type Fake struct {
	addOrUpdateOutboxFn         func(stateStore v1alpha1.Component)
	enabledFn                   func(stateStore string) bool
	publishInternalFn           func(ctx context.Context, stateStore string, states []state.TransactionalStateOperation, source, traceID, traceState string) ([]state.TransactionalStateOperation, error)
	subscribeToInternalTopicsFn func(ctx context.Context, appID string) error
}

func New() *Fake {
	return &Fake{
		addOrUpdateOutboxFn: func(stateStore v1alpha1.Component) {},
		enabledFn:           func(stateStore string) bool { return false },
		publishInternalFn: func(ctx context.Context, stateStore string, states []state.TransactionalStateOperation, source, traceID, traceState string) ([]state.TransactionalStateOperation, error) {
			return nil, nil
		},
		subscribeToInternalTopicsFn: func(ctx context.Context, appID string) error { return nil },
	}
}

func (f *Fake) WithAddOrUpdateOutbox(fn func(stateStore v1alpha1.Component)) *Fake {
	f.addOrUpdateOutboxFn = fn
	return f
}

func (f *Fake) WithEnabled(fn func(stateStore string) bool) *Fake {
	f.enabledFn = fn
	return f
}

func (f *Fake) WithPublishInternal(fn func(ctx context.Context, stateStore string, states []state.TransactionalStateOperation, source, traceID, traceState string) ([]state.TransactionalStateOperation, error)) *Fake {
	f.publishInternalFn = fn
	return f
}

func (f *Fake) WithSubscribeToInternalTopics(fn func(ctx context.Context, appID string) error) *Fake {
	f.subscribeToInternalTopicsFn = fn
	return f
}

func (f *Fake) AddOrUpdateOutbox(stateStore v1alpha1.Component) {
	f.addOrUpdateOutboxFn(stateStore)
}

func (f *Fake) Enabled(stateStore string) bool {
	return f.enabledFn(stateStore)
}

func (f *Fake) PublishInternal(ctx context.Context, stateStore string, states []state.TransactionalStateOperation, source, traceID, traceState string) ([]state.TransactionalStateOperation, error) {
	return f.publishInternalFn(ctx, stateStore, states, source, traceID, traceState)
}

func (f *Fake) SubscribeToInternalTopics(ctx context.Context, appID string) error {
	return f.subscribeToInternalTopicsFn(ctx, appID)
}
````

### FILE: `dapr_official_transactional_outbox/upstream/runtime/pkg/outbox/fake/fake_test.go`
```yaml
block_id: "DAPR-OFFICIAL-OUTBOX:14:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/dapr/dapr/2dcf2e3548faf3c740b3d29ee300e2132db65cff/pkg/outbox/fake/fake_test.go"
license: "Apache-2.0"
sha256: "cac9dfe14be633bdf1034c7b24ca133415f57a8c17338023708a8ab6ba1a4f61"
variables: []
secrets_allowed: false
```
````go
/*
Copyright 2023 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package fake

import (
	"testing"

	"github.com/dapr/dapr/pkg/outbox"
)

func Test_Fake(t *testing.T) {
	var _ outbox.Outbox = New()
}
````

### FILE: `dapr_official_transactional_outbox/upstream/runtime/pkg/runtime/pubsub/outbox.go`
```yaml
block_id: "DAPR-OFFICIAL-OUTBOX:15:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/dapr/dapr/2dcf2e3548faf3c740b3d29ee300e2132db65cff/pkg/runtime/pubsub/outbox.go"
license: "Apache-2.0"
sha256: "8a0a060481e107f3f0617fc1d6baf86c602105d3ef29c429a83abbc0c04923fc"
variables: []
secrets_allowed: false
```
````go
/*
Copyright 2023 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/google/uuid"

	"github.com/dapr/components-contrib/metadata"
	contribPubsub "github.com/dapr/components-contrib/pubsub"
	"github.com/dapr/components-contrib/state"
	"github.com/dapr/dapr/pkg/apis/components/v1alpha1"
	"github.com/dapr/dapr/pkg/outbox"
	"github.com/dapr/kit/logger"
	kitstrings "github.com/dapr/kit/strings"
)

const (
	outboxPublishPubsubKey           = "outboxPublishPubsub"
	outboxPublishTopicKey            = "outboxPublishTopic"
	outboxPubsubKey                  = "outboxPubsub"
	outboxDiscardWhenMissingStateKey = "outboxDiscardWhenMissingState"
	outboxStatePrefix                = "outbox"
	defaultStateScanDelay            = time.Second * 1
)

var outboxLogger = logger.NewLogger("dapr.outbox")

type outboxConfig struct {
	publishPubSub                 string
	publishTopic                  string
	outboxPubsub                  string
	outboxDiscardWhenMissingState bool
}

type outboxImpl struct {
	cloudEventExtractorFn func(map[string]any, string) string
	getPubsubFn           func(string) (contribPubsub.PubSub, bool)
	getStateFn            func(string) (state.Store, bool)
	publisher             Adapter
	outboxStores          map[string]outboxConfig
	lock                  sync.RWMutex
	namespace             string

	// componentCtxFn decorates the context used for the internal outbox
	// subscription and its state operations. Dapr wires this to attach the
	// workload's SPIFFE identity so the pubsub and state components can
	// authenticate to their backing infrastructure service.
	componentCtxFn func(context.Context) context.Context
}

type OptionsOutbox struct {
	Publisher             Adapter
	GetPubsubFn           func(string) (contribPubsub.PubSub, bool)
	GetStateFn            func(string) (state.Store, bool)
	CloudEventExtractorFn func(map[string]any, string) string
	Namespace             string
	ComponentContextFn    func(context.Context) context.Context
}

// NewOutbox returns an instance of an Outbox.
func NewOutbox(opts OptionsOutbox) outbox.Outbox {
	return &outboxImpl{
		cloudEventExtractorFn: opts.CloudEventExtractorFn,
		getPubsubFn:           opts.GetPubsubFn,
		getStateFn:            opts.GetStateFn,
		publisher:             opts.Publisher,
		outboxStores:          make(map[string]outboxConfig),
		namespace:             opts.Namespace,
		componentCtxFn:        opts.ComponentContextFn,
	}
}

// componentContext decorates ctx with the workload's SPIFFE identity for
// component operations. It is a no-op when no decorator is configured.
func (o *outboxImpl) componentContext(ctx context.Context) context.Context {
	if o.componentCtxFn == nil {
		return ctx
	}
	return o.componentCtxFn(ctx)
}

// AddOrUpdateOutbox examines a statestore for outbox properties and saves it for later usage in outbox operations.
func (o *outboxImpl) AddOrUpdateOutbox(stateStore v1alpha1.Component) {
	var (
		publishPubSub, publishTopicKey, outboxPubsub string
		outboxDiscardWhenMissingState                bool
	)

	for _, v := range stateStore.Spec.Metadata {
		switch v.Name {
		case outboxPublishPubsubKey:
			publishPubSub = v.Value.String()
		case outboxPublishTopicKey:
			publishTopicKey = v.Value.String()
		case outboxPubsubKey:
			outboxPubsub = v.Value.String()
		case outboxDiscardWhenMissingStateKey:
			outboxDiscardWhenMissingState = kitstrings.IsTruthy(v.Value.String())
		}
	}

	if publishPubSub != "" && publishTopicKey != "" {
		o.lock.Lock()
		defer o.lock.Unlock()

		if outboxPubsub == "" {
			outboxPubsub = publishPubSub
		}

		o.outboxStores[stateStore.Name] = outboxConfig{
			publishPubSub:                 publishPubSub,
			publishTopic:                  publishTopicKey,
			outboxPubsub:                  outboxPubsub,
			outboxDiscardWhenMissingState: outboxDiscardWhenMissingState,
		}
	}
}

// Enabled returns a bool to indicate if a state store has outbox configured
func (o *outboxImpl) Enabled(stateStore string) bool {
	o.lock.RLock()
	defer o.lock.RUnlock()

	_, ok := o.outboxStores[stateStore]

	return ok
}

func transaction() (state.TransactionalStateOperation, error) {
	uid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return state.SetRequest{
		Key:   outboxStatePrefix + "-" + uid.String(),
		Value: "0",
	}, nil
}

// PublishInternal publishes the state to an internal topic for outbox processing and returns the updated list of transactions
func (o *outboxImpl) PublishInternal(ctx context.Context, stateStore string, operations []state.TransactionalStateOperation, source, traceID, traceState string) ([]state.TransactionalStateOperation, error) {
	o.lock.RLock()
	c, ok := o.outboxStores[stateStore]
	o.lock.RUnlock()

	if !ok {
		return nil, fmt.Errorf("error publishing internal outbox message: could not find outbox configuration on state store %s", stateStore)
	}

	projections := map[string]state.SetRequest{}

	for i, op := range operations {
		sr, ok := op.(state.SetRequest)

		if ok {
			for k, v := range sr.Metadata {
				if k == "outbox.projection" && kitstrings.IsTruthy(v) {
					projections[sr.Key] = sr

					operations = append(operations[:i], operations[i+1:]...)
				}
			}
		}
	}

	for _, op := range operations {
		sr, ok := op.(state.SetRequest)
		if ok {
			tr, err := transaction()
			if err != nil {
				return nil, err
			}

			var (
				payload     any
				contentType string
			)

			if proj, ok := projections[sr.Key]; ok {
				payload = proj.Value

				if proj.ContentType != nil {
					contentType = *proj.ContentType
				} else if ct, ok := proj.Metadata[metadata.ContentType]; ok {
					contentType = ct
				}
			} else {
				payload = sr.Value

				if sr.ContentType != nil {
					contentType = *sr.ContentType
				} else if ct, ok := sr.Metadata[metadata.ContentType]; ok {
					contentType = ct
				}
			}

			var ceData []byte

			bt, ok := payload.([]byte)
			if ok {
				ceData = bt
			} else if contentType != "" && strings.EqualFold(contentType, "application/json") {
				b, sErr := json.Marshal(payload)
				if sErr != nil {
					return nil, sErr
				}

				ceData = b
			} else {
				ceData = fmt.Appendf(nil, "%v", payload)
			}

			var dataContentType string
			if contentType != "" {
				dataContentType = contentType
			}

			ce := contribPubsub.NewCloudEventsEnvelope(tr.GetKey(), source, "", "", "", c.outboxPubsub, dataContentType, ceData, "", traceState)
			ce[contribPubsub.TraceIDField] = traceID

			for k, v := range op.GetMetadata() {
				if k == contribPubsub.DataField || k == contribPubsub.IDField {
					continue
				}

				ce[k] = v
			}

			data, err := json.Marshal(ce)
			if err != nil {
				return nil, err
			}

			err = o.publisher.Publish(ctx, &contribPubsub.PublishRequest{
				PubsubName: c.outboxPubsub,
				Data:       data,
				Topic:      outboxTopic(source, c.publishTopic, o.namespace),
			}, TransportModeGRPC)
			if err != nil {
				return nil, err
			}

			operations = append(operations, tr)
		}
	}

	return operations, nil
}

func outboxTopic(appID, topic, namespace string) string {
	return namespace + appID + topic + "outbox"
}

func (o *outboxImpl) SubscribeToInternalTopics(ctx context.Context, appID string) error {
	o.lock.RLock()
	defer o.lock.RUnlock()

	for stateStore, c := range o.outboxStores {
		outboxPubsub, ok := o.getPubsubFn(c.outboxPubsub)
		if !ok {
			outboxLogger.Warnf("could not subscribe to internal outbox topic: outbox pubsub %s not loaded", c.outboxPubsub)
			continue
		}

		outboxPubsub.Subscribe(o.componentContext(ctx), contribPubsub.SubscribeRequest{
			Topic: outboxTopic(appID, c.publishTopic, o.namespace),
		}, func(ctx context.Context, msg *contribPubsub.NewMessage) error {
			// The per-message context originates from the pubsub component, so
			// re-attach the workload's SPIFFE identity for the state operations
			// (Get/Delete) this handler performs directly.
			ctx = o.componentContext(ctx)

			var cloudEvent map[string]any

			err := json.Unmarshal(msg.Data, &cloudEvent)
			if err != nil {
				return err
			}

			stateKey := o.cloudEventExtractorFn(cloudEvent, contribPubsub.IDField)

			store, ok := o.getStateFn(stateStore)
			if !ok {
				return fmt.Errorf("cannot get outbox state: state store %s not found", stateStore)
			}

			time.Sleep(defaultStateScanDelay)

			bo := &backoff.ExponentialBackOff{
				InitialInterval:     time.Millisecond * 500,
				MaxInterval:         time.Second * 3,
				MaxElapsedTime:      time.Second * 10,
				Multiplier:          3,
				Clock:               backoff.SystemClock,
				RandomizationFactor: 0.1,
			}

			err = backoff.Retry(func() error {
				resp, sErr := store.Get(ctx, &state.GetRequest{
					Key: stateKey,
				})
				if sErr != nil {
					return sErr
				}

				if resp != nil && len(resp.Data) > 0 {
					return nil
				}

				return fmt.Errorf("cannot publish outbox message to topic %s with pubsub %s: outbox state not found", c.publishTopic, c.publishPubSub)
			}, bo)
			if err != nil {
				if c.outboxDiscardWhenMissingState {
					outboxLogger.Errorf("failed to publish outbox topic to pubsub %s: %s, discarding message", c.publishPubSub, err)
					//lint:ignore nilerr dropping message
					return nil
				}

				outboxLogger.Errorf("failed to publish outbox topic to pubsub %s: %s, rejecting for later processing", c.publishPubSub, err)

				return err
			}

			cloudEvent[contribPubsub.TopicField] = c.publishTopic
			cloudEvent[contribPubsub.PubsubField] = c.publishPubSub

			b, err := json.Marshal(cloudEvent)
			if err != nil {
				return err
			}

			contentType := cloudEvent[contribPubsub.DataContentTypeField].(string)

			err = o.publisher.Publish(ctx, &contribPubsub.PublishRequest{
				PubsubName:  c.publishPubSub,
				Data:        b,
				Topic:       c.publishTopic,
				ContentType: &contentType,
			}, TransportModeGRPC)
			if err != nil {
				return err
			}

			err = backoff.Retry(func() error {
				err = store.Delete(ctx, &state.DeleteRequest{
					Key: stateKey,
				})
				if err != nil {
					return err
				}

				return nil
			}, bo)

			return err
		})
	}

	return nil
}
````

### FILE: `dapr_official_transactional_outbox/upstream/runtime/pkg/runtime/pubsub/outbox_test.go`
```yaml
block_id: "DAPR-OFFICIAL-OUTBOX:16:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/dapr/dapr/2dcf2e3548faf3c740b3d29ee300e2132db65cff/pkg/runtime/pubsub/outbox_test.go"
license: "Apache-2.0"
sha256: "fac77524bbbbe8e2a3449ba7ec062f8979a1760085cf56e598acd8670cdd3cc4"
variables: []
secrets_allowed: false
```
````go
/*
Copyright 2023 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pubsub

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	contribPubsub "github.com/dapr/components-contrib/pubsub"
	"github.com/dapr/components-contrib/state"
	"github.com/dapr/dapr/pkg/apis/common"
	"github.com/dapr/dapr/pkg/apis/components/v1alpha1"
	"github.com/dapr/dapr/pkg/outbox"
	"github.com/dapr/dapr/pkg/runtime/pubsub/publisher/fake"
)

func newTestOutbox(publishFn func(context.Context, *contribPubsub.PublishRequest) error) outbox.Outbox {
	p := fake.New()
	if publishFn != nil {
		p.WithPublishFn(publishFn)
	}

	return NewOutbox(OptionsOutbox{
		Publisher:             p,
		CloudEventExtractorFn: extractCloudEventProperty,
	})
}

func TestNewOutbox(t *testing.T) {
	o := newTestOutbox(nil)
	assert.NotNil(t, o)
}

func TestEnabled(t *testing.T) {
	t.Run("required config", func(t *testing.T) {
		o := newTestOutbox(nil).(*outboxImpl)
		o.AddOrUpdateOutbox(v1alpha1.Component{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test",
			},
			Spec: v1alpha1.ComponentSpec{
				Metadata: []common.NameValuePair{
					{
						Name: outboxPublishPubsubKey,
						Value: common.DynamicValue{
							JSON: v1.JSON{
								Raw: []byte("a"),
							},
						},
					},
					{
						Name: outboxPublishTopicKey,
						Value: common.DynamicValue{
							JSON: v1.JSON{
								Raw: []byte("1"),
							},
						},
					},
				},
			},
		})

		assert.True(t, o.Enabled("test"))
		assert.False(t, o.Enabled("test1"))
	})

	t.Run("missing pubsub config", func(t *testing.T) {
		o := newTestOutbox(nil).(*outboxImpl)
		o.AddOrUpdateOutbox(v1alpha1.Component{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test",
			},
			Spec: v1alpha1.ComponentSpec{
				Metadata: []common.NameValuePair{
					{
						Name: outboxPublishTopicKey,
						Value: common.DynamicValue{
							JSON: v1.JSON{
								Raw: []byte("1"),
							},
						},
					},
				},
			},
		})

		assert.False(t, o.Enabled("test"))
		assert.False(t, o.Enabled("test1"))
	})

	t.Run("missing topic config", func(t *testing.T) {
		o := newTestOutbox(nil).(*outboxImpl)
		o.AddOrUpdateOutbox(v1alpha1.Component{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test",
			},
			Spec: v1alpha1.ComponentSpec{
				Metadata: []common.NameValuePair{
					{
						Name: outboxPublishPubsubKey,
						Value: common.DynamicValue{
							JSON: v1.JSON{
								Raw: []byte("a"),
							},
						},
					},
				},
			},
		})

		assert.False(t, o.Enabled("test"))
		assert.False(t, o.Enabled("test1"))
	})
}

func TestAddOrUpdateOutbox(t *testing.T) {
	t.Run("config values correct", func(t *testing.T) {
		o := newTestOutbox(nil).(*outboxImpl)
		o.AddOrUpdateOutbox(v1alpha1.Component{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test",
			},
			Spec: v1alpha1.ComponentSpec{
				Metadata: []common.NameValuePair{
					{
						Name: outboxPublishPubsubKey,
						Value: common.DynamicValue{
							JSON: v1.JSON{
								Raw: []byte("a"),
							},
						},
					},
					{
						Name: outboxPublishTopicKey,
						Value: common.DynamicValue{
							JSON: v1.JSON{
								Raw: []byte("1"),
							},
						},
					},
					{
						Name: outboxPubsubKey,
						Value: common.DynamicValue{
							JSON: v1.JSON{
								Raw: []byte("2"),
							},
						},
					},
				},
			},
		})

		c := o.outboxStores["test"]
		assert.Equal(t, "2", c.outboxPubsub)
		assert.Equal(t, "a", c.publishPubSub)
		assert.Equal(t, "1", c.publishTopic)
	})

	t.Run("config default values correct", func(t *testing.T) {
		o := newTestOutbox(nil).(*outboxImpl)
		o.AddOrUpdateOutbox(v1alpha1.Component{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test",
			},
			Spec: v1alpha1.ComponentSpec{
				Metadata: []common.NameValuePair{
					{
						Name: outboxPublishPubsubKey,
						Value: common.DynamicValue{
							JSON: v1.JSON{
								Raw: []byte("a"),
							},
						},
					},
					{
						Name: outboxPublishTopicKey,
						Value: common.DynamicValue{
							JSON: v1.JSON{
								Raw: []byte("1"),
							},
						},
					},
				},
			},
		})

		c := o.outboxStores["test"]
		assert.Equal(t, "a", c.outboxPubsub)
		assert.Equal(t, "a", c.publishPubSub)
		assert.Equal(t, "1", c.publishTopic)
	})
}

func TestPublishInternal(t *testing.T) {
	t.Run("valid operation, correct default parameters", func(t *testing.T) {
		o := newTestOutbox(func(ctx context.Context, pr *contribPubsub.PublishRequest) error {
			var cloudEvent map[string]any

			err := json.Unmarshal(pr.Data, &cloudEvent)
			require.NoError(t, err)

			assert.Equal(t, "test", cloudEvent["data"])
			assert.Equal(t, "a", pr.PubsubName)
			assert.Equal(t, "testapp1outbox", pr.Topic)
			assert.Equal(t, "testapp", cloudEvent["source"])
			assert.Equal(t, "text/plain", cloudEvent["datacontenttype"])
			assert.Equal(t, "a", cloudEvent["pubsubname"])

			return nil
		}).(*outboxImpl)

		o.AddOrUpdateOutbox(v1alpha1.Component{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test",
			},
			Spec: v1alpha1.ComponentSpec{
				Metadata: []common.NameValuePair{
					{
						Name: outboxPublishPubsubKey,
						Value: common.DynamicValue{
							JSON: v1.JSON{
								Raw: []byte("a"),
							},
						},
					},
					{
						Name: outboxPublishTopicKey,
						Value: common.DynamicValue{
							JSON: v1.JSON{
								Raw: []byte("1"),
							},
						},
					},
				},
			},
		})

		_, err := o.PublishInternal(t.Context(), "test", []state.TransactionalStateOperation{
			state.SetRequest{
				Key:   "key",
				Value: "test",
			},
		}, "testapp", "", "")

		require.NoError(t, err)
	})

	t.Run("valid operation, correct overridden parameters", func(t *testing.T) {
		o := newTestOutbox(func(ctx context.Context, pr *contribPubsub.PublishRequest) error {
			var cloudEvent map[string]any

			err := json.Unmarshal(pr.Data, &cloudEvent)
			require.NoError(t, err)

			assert.Equal(t, "test", cloudEvent["data"])
			assert.Equal(t, "a", pr.PubsubName)
			assert.Equal(t, "testapp1outbox", pr.Topic)
			assert.Equal(t, "testsource", cloudEvent["source"])
			assert.Equal(t, "text/plain", cloudEvent["datacontenttype"])
			assert.Equal(t, "a", cloudEvent["pubsubname"])

			return nil
		}).(*outboxImpl)

		o.AddOrUpdateOutbox(v1alpha1.Component{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test",
			},
			Spec: v1alpha1.ComponentSpec{
				Metadata: []common.NameValuePair{
					{
						Name: outboxPublishPubsubKey,
						Value: common.DynamicValue{
							JSON: v1.JSON{
								Raw: []byte("a"),
							},
						},
					},
					{
						Name: outboxPublishTopicKey,
						Value: common.DynamicValue{
							JSON: v1.JSON{
								Raw: []byte("1"),
							},
						},
					},
				},
			},
		})

		_, err := o.PublishInternal(t.Context(), "test", []state.TransactionalStateOperation{
			state.SetRequest{
				Key:      "key",
				Value:    "test",
				Metadata: map[string]string{"source": "testsource"},
			},
		}, "testapp", "", "")

		require.NoError(t, err)
	})

	t.Run("valid operation, no datacontenttype", func(t *testing.T) {
		o := newTestOutbox(func(ctx context.Context, pr *contribPubsub.PublishRequest) error {
			var cloudEvent map[string]any

			err := json.Unmarshal(pr.Data, &cloudEvent)
			require.NoError(t, err)

			assert.Equal(t, "test", cloudEvent["data"])
			assert.Equal(t, "a", pr.PubsubName)
			assert.Equal(t, "testapp1outbox", pr.Topic)
			assert.Equal(t, "testapp", cloudEvent["source"])
			assert.Equal(t, "text/plain", cloudEvent["datacontenttype"])
			assert.Equal(t, "a", cloudEvent["pubsubname"])

			return nil
		}).(*outboxImpl)

		o.AddOrUpdateOutbox(v1alpha1.Component{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test",
			},
			Spec: v1alpha1.ComponentSpec{
				Metadata: []common.NameValuePair{
					{
						Name: outboxPublishPubsubKey,
						Value: common.DynamicValue{
							JSON: v1.JSON{
								Raw: []byte("a"),
							},
						},
					},
					{
						Name: outboxPublishTopicKey,
						Value: common.DynamicValue{
							JSON: v1.JSON{
								Raw: []byte("1"),
							},
						},
					},
				},
			},
		})

		contentType := ""
		_, err := o.PublishInternal(t.Context(), "test", []state.TransactionalStateOperation{
			state.SetRequest{
				Key:         "key",
				Value:       "test",
				ContentType: &contentType,
			},
		}, "testapp", "", "")

		require.NoError(t, err)
	})

	type customData struct {
		Name string `json:"name"`
	}

	t.Run("valid operation, application/json datacontenttype", func(t *testing.T) {
		o := newTestOutbox(func(ctx context.Context, pr *contribPubsub.PublishRequest) error {
			var cloudEvent map[string]any

			err := json.Unmarshal(pr.Data, &cloudEvent)
			require.NoError(t, err)

			data := cloudEvent["data"]
			j := customData{}

			err = json.Unmarshal([]byte(data.(string)), &j)
			require.NoError(t, err)

			assert.Equal(t, "test", j.Name)
			assert.Equal(t, "a", pr.PubsubName)
			assert.Equal(t, "testapp1outbox", pr.Topic)
			assert.Equal(t, "testapp", cloudEvent["source"])
			assert.Equal(t, "application/json", cloudEvent["datacontenttype"])
			assert.Equal(t, "a", cloudEvent["pubsubname"])

			return nil
		}).(*outboxImpl)

		o.AddOrUpdateOutbox(v1alpha1.Component{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test",
			},
			Spec: v1alpha1.ComponentSpec{
				Metadata: []common.NameValuePair{
					{
						Name: outboxPublishPubsubKey,
						Value: common.DynamicValue{
							JSON: v1.JSON{
								Raw: []byte("a"),
							},
						},
					},
					{
						Name: outboxPublishTopicKey,
						Value: common.DynamicValue{
							JSON: v1.JSON{
								Raw: []byte("1"),
							},
						},
					},
				},
			},
		})

		j := customData{
			Name: "test",
		}
		b, err := json.Marshal(&j)
		require.NoError(t, err)

		contentType := "application/json"
		_, err = o.PublishInternal(t.Context(), "test", []state.TransactionalStateOperation{
			state.SetRequest{
				Key:         "key",
				Value:       string(b),
				ContentType: &contentType,
			},
		}, "testapp", "", "")

		require.NoError(t, err)
	})

	t.Run("valid operation, application/json contenttype metadata", func(t *testing.T) {
		o := newTestOutbox(func(ctx context.Context, pr *contribPubsub.PublishRequest) error {
			var cloudEvent map[string]any

			err := json.Unmarshal(pr.Data, &cloudEvent)
			require.NoError(t, err)

			data := cloudEvent["data"]
			j := customData{}

			err = json.Unmarshal([]byte(data.(string)), &j)
			require.NoError(t, err)

			assert.Equal(t, "test", j.Name)
			assert.Equal(t, "a", pr.PubsubName)
			assert.Equal(t, "testapp1outbox", pr.Topic)
			assert.Equal(t, "testapp", cloudEvent["source"])
			assert.Equal(t, "application/json", cloudEvent["datacontenttype"])
			assert.Equal(t, "a", cloudEvent["pubsubname"])

			return nil
		}).(*outboxImpl)

		o.AddOrUpdateOutbox(v1alpha1.Component{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test",
			},
			Spec: v1alpha1.ComponentSpec{
				Metadata: []common.NameValuePair{
					{
						Name: outboxPublishPubsubKey,
						Value: common.DynamicValue{
							JSON: v1.JSON{
								Raw: []byte("a"),
							},
						},
					},
					{
						Name: outboxPublishTopicKey,
						Value: common.DynamicValue{
							JSON: v1.JSON{
								Raw: []byte("1"),
							},
						},
					},
				},
			},
		})

		j := customData{
			Name: "test",
		}
		b, err := json.Marshal(&j)
		require.NoError(t, err)

		_, err = o.PublishInternal(t.Context(), "test", []state.TransactionalStateOperation{
			state.SetRequest{
				Key:      "key",
				Value:    string(b),
				Metadata: map[string]string{"contentType": "application/json"},
			},
		}, "testapp", "", "")

		require.NoError(t, err)
	})

	t.Run("valid operation, application/json contenttype metadata outbox projection", func(t *testing.T) {
		o := newTestOutbox(func(ctx context.Context, pr *contribPubsub.PublishRequest) error {
			var cloudEvent map[string]any

			err := json.Unmarshal(pr.Data, &cloudEvent)
			require.NoError(t, err)

			data := cloudEvent["data"]
			j := customData{}

			err = json.Unmarshal([]byte(data.(string)), &j)
			require.NoError(t, err)

			assert.Equal(t, "projection", j.Name)
			assert.Equal(t, "a", pr.PubsubName)
			assert.Equal(t, "testapp1outbox", pr.Topic)
			assert.Equal(t, "testapp", cloudEvent["source"])
			assert.Equal(t, "application/json", cloudEvent["datacontenttype"])
			assert.Equal(t, "a", cloudEvent["pubsubname"])

			return nil
		}).(*outboxImpl)

		o.AddOrUpdateOutbox(v1alpha1.Component{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test",
			},
			Spec: v1alpha1.ComponentSpec{
				Metadata: []common.NameValuePair{
					{
						Name: outboxPublishPubsubKey,
						Value: common.DynamicValue{
							JSON: v1.JSON{
								Raw: []byte("a"),
							},
						},
					},
					{
						Name: outboxPublishTopicKey,
						Value: common.DynamicValue{
							JSON: v1.JSON{
								Raw: []byte("1"),
							},
						},
					},
				},
			},
		})

		j := customData{
			Name: "test",
		}

		projection := customData{
			Name: "projection",
		}

		b, err := json.Marshal(&j)
		require.NoError(t, err)

		jp, err := json.Marshal(&projection)
		require.NoError(t, err)

		_, err = o.PublishInternal(t.Context(), "test", []state.TransactionalStateOperation{
			state.SetRequest{
				Key:   "key",
				Value: string(b),
			},
			state.SetRequest{
				Key:      "key",
				Value:    string(jp),
				Metadata: map[string]string{"contentType": "application/json", "outbox.projection": "true"},
			},
		}, "testapp", "", "")

		require.NoError(t, err)
	})

	t.Run("missing state store", func(t *testing.T) {
		o := newTestOutbox(nil).(*outboxImpl)

		_, err := o.PublishInternal(t.Context(), "test", []state.TransactionalStateOperation{
			state.SetRequest{
				Key:   "key",
				Value: "test",
			},
		}, "testapp", "", "")
		require.Error(t, err)
	})

	t.Run("no op when no transactions", func(t *testing.T) {
		o := newTestOutbox(func(ctx context.Context, pr *contribPubsub.PublishRequest) error {
			assert.Fail(t, "unexptected message received")
			return nil
		}).(*outboxImpl)

		o.AddOrUpdateOutbox(v1alpha1.Component{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test",
			},
			Spec: v1alpha1.ComponentSpec{
				Metadata: []common.NameValuePair{
					{
						Name: outboxPublishPubsubKey,
						Value: common.DynamicValue{
							JSON: v1.JSON{
								Raw: []byte("a"),
							},
						},
					},
					{
						Name: outboxPublishTopicKey,
						Value: common.DynamicValue{
							JSON: v1.JSON{
								Raw: []byte("1"),
							},
						},
					},
				},
			},
		})

		_, err := o.PublishInternal(t.Context(), "test", []state.TransactionalStateOperation{}, "testapp", "", "")

		require.NoError(t, err)
	})

	t.Run("error when pubsub fails", func(t *testing.T) {
		o := newTestOutbox(func(ctx context.Context, pr *contribPubsub.PublishRequest) error {
			return errors.New("")
		}).(*outboxImpl)

		o.AddOrUpdateOutbox(v1alpha1.Component{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test",
			},
			Spec: v1alpha1.ComponentSpec{
				Metadata: []common.NameValuePair{
					{
						Name: outboxPublishPubsubKey,
						Value: common.DynamicValue{
							JSON: v1.JSON{
								Raw: []byte("a"),
							},
						},
					},
					{
						Name: outboxPublishTopicKey,
						Value: common.DynamicValue{
							JSON: v1.JSON{
								Raw: []byte("1"),
							},
						},
					},
				},
			},
		})

		_, err := o.PublishInternal(t.Context(), "test", []state.TransactionalStateOperation{
			state.SetRequest{
				Key:   "1",
				Value: "hello",
			},
		}, "testapp", "", "")

		require.Error(t, err)
	})
}

func TestSubscribeToInternalTopics(t *testing.T) {
	t.Run("correct configuration with trace, custom field and nonoverridable fields", func(t *testing.T) {
		const outboxTopic = "test1outbox"

		psMock := &outboxPubsubMock{
			expectedOutboxTopic: outboxTopic,
			t:                   t,
		}
		stateMock := &outboxStateMock{
			receivedKey: make(chan string, 1),
		}

		internalCalledCh := make(chan struct{})
		externalCalledCh := make(chan struct{})

		var closed bool

		o := newTestOutbox(func(ctx context.Context, pr *contribPubsub.PublishRequest) error {
			switch pr.Topic {
			case outboxTopic:
				close(internalCalledCh)
			case "1":
				if !closed {
					close(externalCalledCh)

					closed = true
				}
			}

			ce := map[string]string{}
			json.Unmarshal(pr.Data, &ce)

			traceID := ce[contribPubsub.TraceIDField]
			traceState := ce[contribPubsub.TraceStateField]
			customField := ce["outbox.cloudevent.customfield"]
			data := ce[contribPubsub.DataField]
			id := ce[contribPubsub.IDField]

			assert.Equal(t, "00-ecdf5aaa79bff09b62b201442c0f3061-d2597ed7bfd029e4-01", traceID)
			assert.Equal(t, "00-ecdf5aaa79bff09b62b201442c0f3061-d2597ed7bfd029e4-01", traceState)
			assert.Equal(t, "a", customField)
			assert.Equal(t, "hello", data)
			assert.Contains(t, id, "outbox-")

			return psMock.Publish(ctx, pr)
		}).(*outboxImpl)
		o.cloudEventExtractorFn = extractCloudEventProperty

		o.getPubsubFn = func(s string) (contribPubsub.PubSub, bool) {
			return psMock, true
		}
		o.getStateFn = func(s string) (state.Store, bool) {
			return stateMock, true
		}

		stateScan := "1s"

		o.AddOrUpdateOutbox(v1alpha1.Component{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test",
			},
			Spec: v1alpha1.ComponentSpec{
				Metadata: []common.NameValuePair{
					{
						Name: outboxPublishPubsubKey,
						Value: common.DynamicValue{
							JSON: v1.JSON{
								Raw: []byte("a"),
							},
						},
					},
					{
						Name: outboxPublishTopicKey,
						Value: common.DynamicValue{
							JSON: v1.JSON{
								Raw: []byte("1"),
							},
						},
					},
				},
			},
		})

		const appID = "test"

		err := o.SubscribeToInternalTopics(t.Context(), appID)
		require.NoError(t, err)

		errCh := make(chan error, 1)

		go func() {
			trs, pErr := o.PublishInternal(t.Context(), "test", []state.TransactionalStateOperation{
				state.SetRequest{
					Key:      "1",
					Value:    "hello",
					Metadata: map[string]string{"outbox.cloudevent.customfield": "a", "data": "a", "id": "b"},
				},
			}, appID, "00-ecdf5aaa79bff09b62b201442c0f3061-d2597ed7bfd029e4-01", "00-ecdf5aaa79bff09b62b201442c0f3061-d2597ed7bfd029e4-01")

			trs = append(trs[:0], trs[0+1:]...)

			if pErr != nil {
				errCh <- pErr
				return
			}

			if len(trs) != 1 {
				errCh <- fmt.Errorf("expected trs to have len(1), but got: %d", len(trs))
				return
			}

			errCh <- nil

			stateMock.expectedKey.Store(new(trs[0].GetKey()))
		}()

		d, err := time.ParseDuration(stateScan)
		require.NoError(t, err)

		start := time.Now()
		doneCh := make(chan error, 2)
		timeout := time.After(5 * time.Second)

		go func() {
			select {
			case <-internalCalledCh:
				doneCh <- nil
			case <-timeout:
				doneCh <- errors.New("timeout waiting for internalCalledCh")
			}
		}()
		go func() {
			select {
			case <-externalCalledCh:
				doneCh <- nil
			case <-timeout:
				doneCh <- errors.New("timeout waiting for externalCalledCh")
			}
		}()

		for range 2 {
			require.NoError(t, <-doneCh)
		}

		require.GreaterOrEqual(t, time.Since(start), d)

		// Publishing should not have errored
		require.NoError(t, <-errCh)

		expected := stateMock.expectedKey.Load()
		require.NotNil(t, expected)
		assert.Equal(t, *expected, <-stateMock.receivedKey)
	})

	t.Run("state store not present", func(t *testing.T) {
		const outboxTopic = "test1outbox"

		psMock := &outboxPubsubMock{
			expectedOutboxTopic: outboxTopic,
			t:                   t,
		}

		o := newTestOutbox(func(ctx context.Context, pr *contribPubsub.PublishRequest) error {
			return psMock.Publish(ctx, pr)
		}).(*outboxImpl)

		o.getPubsubFn = func(s string) (contribPubsub.PubSub, bool) {
			return psMock, true
		}
		o.getStateFn = func(s string) (state.Store, bool) {
			return nil, false
		}

		const appID = "test"

		err := o.SubscribeToInternalTopics(t.Context(), appID)
		require.NoError(t, err)

		trs, pErr := o.PublishInternal(t.Context(), "test", []state.TransactionalStateOperation{
			state.SetRequest{
				Key:   "1",
				Value: "hello",
			},
		}, appID, "", "")

		require.Error(t, pErr)
		assert.Empty(t, trs)
	})

	t.Run("outbox state not present", func(t *testing.T) {
		const outboxTopic = "test1outbox"

		psMock := &outboxPubsubMock{
			expectedOutboxTopic: outboxTopic,
			t:                   t,
		}
		stateMock := &outboxStateMock{}

		internalCalledCh := make(chan struct{})
		externalCalledCh := make(chan struct{})
		o := newTestOutbox(func(ctx context.Context, pr *contribPubsub.PublishRequest) error {
			switch pr.Topic {
			case outboxTopic:
				close(internalCalledCh)
			case "1":
				close(externalCalledCh)
			}

			return psMock.Publish(ctx, pr)
		}).(*outboxImpl)

		o.getPubsubFn = func(s string) (contribPubsub.PubSub, bool) {
			return psMock, true
		}
		o.getStateFn = func(s string) (state.Store, bool) {
			return stateMock, true
		}

		const stateScan = "1s"

		o.AddOrUpdateOutbox(v1alpha1.Component{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test",
			},
			Spec: v1alpha1.ComponentSpec{
				Metadata: []common.NameValuePair{
					{
						Name: outboxPublishPubsubKey,
						Value: common.DynamicValue{
							JSON: v1.JSON{
								Raw: []byte("a"),
							},
						},
					},
					{
						Name: outboxPublishTopicKey,
						Value: common.DynamicValue{
							JSON: v1.JSON{
								Raw: []byte("1"),
							},
						},
					},
				},
			},
		})

		const appID = "test"

		err := o.SubscribeToInternalTopics(t.Context(), appID)
		require.NoError(t, err)

		errCh := make(chan error, 1)

		go func() {
			trs, pErr := o.PublishInternal(t.Context(), "test", []state.TransactionalStateOperation{
				state.SetRequest{
					Key:   "1",
					Value: "hello",
				},
			}, appID, "", "")

			trs = append(trs[:0], trs[0+1:]...)

			if pErr != nil {
				errCh <- pErr
				return
			}

			if len(trs) != 1 {
				errCh <- fmt.Errorf("expected trs to have len(1), but got: %d", len(trs))
				return
			}

			errCh <- nil
		}()

		d, err := time.ParseDuration(stateScan)
		require.NoError(t, err)

		start := time.Now()
		doneCh := make(chan error, 2)
		timeout := time.After(2 * time.Second)

		go func() {
			select {
			case <-internalCalledCh:
				doneCh <- nil
			case <-timeout:
				doneCh <- errors.New("timeout waiting for internalCalledCh")
			}
		}()
		go func() {
			// Here we expect no signal
			select {
			case <-externalCalledCh:
				doneCh <- errors.New("received unexpected signal on externalCalledCh")
			case <-timeout:
				doneCh <- nil
			}
		}()

		for range 2 {
			require.NoError(t, <-doneCh)
		}

		require.GreaterOrEqual(t, time.Since(start), d)

		// Publishing should not have errored
		require.NoError(t, <-errCh)
	})

	t.Run("outbox state not present with discard", func(t *testing.T) {
		const outboxTopic = "test1outbox"

		psMock := &outboxPubsubMock{
			expectedOutboxTopic: outboxTopic,
			t:                   t,
			validateNoError:     true,
		}
		stateMock := &outboxStateMock{
			returnEmptyOnGet: true,
		}

		internalCalledCh := make(chan struct{})
		externalCalledCh := make(chan struct{})

		o := newTestOutbox(func(ctx context.Context, pr *contribPubsub.PublishRequest) error {
			switch pr.Topic {
			case outboxTopic:
				close(internalCalledCh)
			case "1":
				close(externalCalledCh)
			}

			return psMock.Publish(ctx, pr)
		}).(*outboxImpl)

		o.getPubsubFn = func(s string) (contribPubsub.PubSub, bool) {
			return psMock, true
		}
		o.getStateFn = func(s string) (state.Store, bool) {
			return stateMock, true
		}

		const stateScan = "1s"

		o.AddOrUpdateOutbox(v1alpha1.Component{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test",
			},
			Spec: v1alpha1.ComponentSpec{
				Metadata: []common.NameValuePair{
					{
						Name: outboxPublishPubsubKey,
						Value: common.DynamicValue{
							JSON: v1.JSON{
								Raw: []byte("a"),
							},
						},
					},
					{
						Name: outboxPublishTopicKey,
						Value: common.DynamicValue{
							JSON: v1.JSON{
								Raw: []byte("1"),
							},
						},
					},
					{
						Name: outboxDiscardWhenMissingStateKey,
						Value: common.DynamicValue{
							JSON: v1.JSON{
								Raw: []byte("true"),
							},
						},
					},
				},
			},
		})

		const appID = "test"

		err := o.SubscribeToInternalTopics(t.Context(), appID)
		require.NoError(t, err)

		errCh := make(chan error, 1)

		go func() {
			trs, pErr := o.PublishInternal(t.Context(), "test", []state.TransactionalStateOperation{
				state.SetRequest{
					Key:   "1",
					Value: "hello",
				},
			}, appID, "", "")

			trs = append(trs[:0], trs[0+1:]...)

			if pErr != nil {
				errCh <- pErr
				return
			}

			if len(trs) != 1 {
				errCh <- fmt.Errorf("expected trs to have len(1), but got: %d", len(trs))
				return
			}

			errCh <- nil
		}()

		d, err := time.ParseDuration(stateScan)
		require.NoError(t, err)

		start := time.Now()
		doneCh := make(chan error, 2)

		// account for max retry time
		timeout := time.After(11 * time.Second)

		go func() {
			select {
			case <-internalCalledCh:
				doneCh <- nil
			case <-timeout:
				doneCh <- errors.New("timeout waiting for internalCalledCh")
			}
		}()
		go func() {
			// Here we expect no signal
			select {
			case <-externalCalledCh:
				doneCh <- errors.New("received unexpected signal on externalCalledCh")
			case <-timeout:
				doneCh <- nil
			}
		}()

		for range 2 {
			require.NoError(t, <-doneCh)
		}

		require.GreaterOrEqual(t, time.Since(start), d)

		// Publishing should not have errored
		require.NoError(t, <-errCh)
	})
}

type outboxPubsubMock struct {
	expectedOutboxTopic string
	t                   *testing.T
	handler             contribPubsub.Handler
	validateNoError     bool
}

func (o *outboxPubsubMock) Init(ctx context.Context, metadata contribPubsub.Metadata) error {
	return nil
}

func (o *outboxPubsubMock) Features() []contribPubsub.Feature {
	return nil
}

func (o *outboxPubsubMock) Publish(ctx context.Context, req *contribPubsub.PublishRequest) error {
	go func() {
		err := o.handler(context.Background(), &contribPubsub.NewMessage{
			Data:  req.Data,
			Topic: req.Topic,
		})

		if o.validateNoError {
			require.NoError(o.t, err)
			return
		}
	}()

	return nil
}

func (o *outboxPubsubMock) Subscribe(ctx context.Context, req contribPubsub.SubscribeRequest, handler contribPubsub.Handler) error {
	if req.Topic != o.expectedOutboxTopic {
		assert.Fail(o.t, fmt.Sprintf("expected outbox topic %s, got %s", o.expectedOutboxTopic, req.Topic))
	}

	o.handler = handler

	return nil
}

func (o *outboxPubsubMock) Close() error {
	return nil
}

type outboxStateMock struct {
	expectedKey      atomic.Pointer[string]
	receivedKey      chan string
	returnEmptyOnGet bool
}

func (o *outboxStateMock) Init(ctx context.Context, metadata state.Metadata) error {
	return nil
}

func (o *outboxStateMock) Features() []state.Feature {
	return nil
}

func (o *outboxStateMock) Delete(ctx context.Context, req *state.DeleteRequest) error {
	return nil
}

func (o *outboxStateMock) Get(ctx context.Context, req *state.GetRequest) (*state.GetResponse, error) {
	if o.returnEmptyOnGet {
		return &state.GetResponse{}, nil
	}

	if o.receivedKey != nil {
		o.receivedKey <- req.Key
	}

	expected := o.expectedKey.Load()

	if expected != nil && *expected != "" && *expected == req.Key {
		return &state.GetResponse{
			Data: []byte("0"),
		}, nil
	}

	return nil, nil
}

func TestOutboxTopic(t *testing.T) {
	t.Run("not namespaced", func(t *testing.T) {
		o := newTestOutbox(nil).(*outboxImpl)
		topic := outboxTopic("a", "b", o.namespace)

		assert.Equal(t, "aboutbox", topic)
	})

	t.Run("namespaced", func(t *testing.T) {
		o := newTestOutbox(nil).(*outboxImpl)
		o.namespace = "default"

		topic := outboxTopic("a", "b", o.namespace)

		assert.Equal(t, "defaultaboutbox", topic)
	})
}

func (o *outboxStateMock) Set(ctx context.Context, req *state.SetRequest) error {
	return nil
}

func (o *outboxStateMock) BulkGet(ctx context.Context, req []state.GetRequest, opts state.BulkGetOpts) ([]state.BulkGetResponse, error) {
	return nil, nil
}

func (o *outboxStateMock) BulkSet(ctx context.Context, req []state.SetRequest, opts state.BulkStoreOpts) error {
	return nil
}

func (o *outboxStateMock) BulkDelete(ctx context.Context, req []state.DeleteRequest, opts state.BulkStoreOpts) error {
	return nil
}

func (o *outboxStateMock) Close() error {
	return nil
}

func extractCloudEventProperty(cloudEvent map[string]any, property string) string {
	if cloudEvent == nil {
		return ""
	}

	iValue, ok := cloudEvent[property]
	if ok {
		if value, ok := iValue.(string); ok {
			return value
		}
	}

	return ""
}
````

### FILE: `dapr_official_transactional_outbox/upstream/runtime/tests/integration/suite/daprd/outbox/http/basic.go`
```yaml
block_id: "DAPR-OFFICIAL-OUTBOX:17:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/dapr/dapr/2dcf2e3548faf3c740b3d29ee300e2132db65cff/tests/integration/suite/daprd/outbox/http/basic.go"
license: "Apache-2.0"
sha256: "ce9ab635518713343b8e3db3fe98b12357a55086102ae2343fe80d1c34078bfe"
variables: []
secrets_allowed: false
```
````go
/*
Copyright 2023 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dapr/components-contrib/state"
	"github.com/dapr/dapr/tests/integration/framework"
	"github.com/dapr/dapr/tests/integration/framework/client"
	procdaprd "github.com/dapr/dapr/tests/integration/framework/process/daprd"
	prochttp "github.com/dapr/dapr/tests/integration/framework/process/http"
	"github.com/dapr/dapr/tests/integration/suite"
)

func init() {
	suite.Register(new(basic))
}

type basic struct {
	appTestCalled atomic.Int32
	daprd         *procdaprd.Daprd
}

func (o *basic) Setup(t *testing.T) []framework.Option {
	newHTTPServer := func() *prochttp.HTTP {
		handler := http.NewServeMux()
		var msg atomic.Value

		handler.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
			o.appTestCalled.Add(1)
			defer r.Body.Close()
			b, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatal(err)
			}

			msg.Store(b)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		})

		handler.HandleFunc("/getValue", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			m := msg.Load()
			if m == nil {
				return
			}
			w.Write(msg.Load().([]byte))
		})

		return prochttp.New(t, prochttp.WithHandler(handler))
	}
	srv1 := newHTTPServer()

	o.daprd = procdaprd.New(t, procdaprd.WithAppID("outboxtest"), procdaprd.WithAppPort(srv1.Port()), procdaprd.WithResourceFiles(`
apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
  name: mystore
spec:
  type: state.in-memory
  version: v1
  metadata:
  - name: outboxPublishPubsub
    value: "mypubsub"
  - name: outboxPublishTopic
    value: "test"
`,
		`
apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
  name: 'mypubsub'
spec:
  type: pubsub.in-memory
  version: v1
`,
		`
apiVersion: dapr.io/v2alpha1
kind: Subscription
metadata:
  name: 'order'
spec:
  topic: 'test'
  routes:
    default: '/test'
  pubsubname: 'mypubsub'
scopes:
- outboxtest
`))

	return []framework.Option{
		framework.WithProcesses(srv1, o.daprd),
	}
}

type stateTransactionRequestBody struct {
	Operations []stateTransactionRequestBodyOperation `json:"operations"`
	Metadata   map[string]string                      `json:"metadata,omitempty"`
}

type stateTransactionRequestBodyOperation struct {
	Operation string `json:"operation"`
	Request   any    `json:"request"`
}

func (o *basic) Run(t *testing.T, ctx context.Context) {
	o.daprd.WaitUntilRunning(t, ctx)

	postURL := fmt.Sprintf("http://localhost:%d/v1.0/state/mystore/transaction", o.daprd.HTTPPort())
	stateReq := state.SetRequest{
		Key:      "1",
		Value:    "2",
		Metadata: map[string]string{"outbox.cloudevent.myapp": "myapp1", "data": "a", "id": "b"},
	}

	tr := stateTransactionRequestBody{
		Operations: []stateTransactionRequestBodyOperation{
			{
				Operation: "upsert",
				Request:   stateReq,
			},
		},
	}

	b, err := json.Marshal(&tr)
	require.NoError(t, err)

	httpClient := client.HTTP(t)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, postURL, bytes.NewReader(b))
	require.NoError(t, err)
	resp, err := httpClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())
	assert.Empty(t, string(body))

	assert.EventuallyWithT(t, func(c *assert.CollectT) {
		req, err = http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("http://localhost:%v/getValue", o.daprd.AppPort(t)), nil)
		assert.NoError(c, err)
		resp, err = httpClient.Do(req)
		assert.NoError(c, err)
		t.Cleanup(func() {
			assert.NoError(t, resp.Body.Close())
		})
		body, err = io.ReadAll(resp.Body)
		assert.NoError(c, err)

		var ce map[string]string
		err = json.Unmarshal(body, &ce)
		assert.NoError(c, err)
		assert.Equal(c, "2", ce["data"])
		assert.Equal(c, "myapp1", ce["outbox.cloudevent.myapp"])
		assert.Contains(c, ce["id"], "outbox-")
	}, time.Second*10, time.Millisecond*10)

	assert.Equal(t, int32(1), o.appTestCalled.Load())
}
````

### FILE: `dapr_official_transactional_outbox/upstream/runtime/tests/integration/suite/daprd/outbox/http/delete.go`
```yaml
block_id: "DAPR-OFFICIAL-OUTBOX:18:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/dapr/dapr/2dcf2e3548faf3c740b3d29ee300e2132db65cff/tests/integration/suite/daprd/outbox/http/delete.go"
license: "Apache-2.0"
sha256: "222dd239f765a04accac212563e56b200884cb77246816388ec80521efd2dbfb"
variables: []
secrets_allowed: false
```
````go
/*
Copyright 2024 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dapr/components-contrib/state"
	"github.com/dapr/dapr/tests/integration/framework"
	"github.com/dapr/dapr/tests/integration/framework/client"
	"github.com/dapr/dapr/tests/integration/framework/process/daprd"
	"github.com/dapr/dapr/tests/integration/framework/process/http/app"
	"github.com/dapr/dapr/tests/integration/suite"
)

func init() {
	suite.Register(new(delete))
}

type delete struct {
	appTestCalled atomic.Int32
	daprd         *daprd.Daprd
	msg           atomic.Value
}

func (o *delete) Setup(t *testing.T) []framework.Option {
	app := app.New(t,
		app.WithHandlerFunc("/test", func(w http.ResponseWriter, r *http.Request) {
			defer o.appTestCalled.Add(1)
			var ce map[string]string
			assert.NoError(t, json.NewDecoder(r.Body).Decode(&ce))
			o.msg.Store(ce)
		}),
	)
	o.daprd = daprd.New(t,
		daprd.WithAppID("outboxtest"),
		daprd.WithAppPort(app.Port()),
		daprd.WithResourceFiles(`
apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
  name: mystore
spec:
  type: state.in-memory
  version: v1
  metadata:
  - name: outboxPublishPubsub
    value: "mypubsub"
  - name: outboxPublishTopic
    value: "test"
`,
			`
apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
  name: 'mypubsub'
spec:
  type: pubsub.in-memory
  version: v1
`,
			`
apiVersion: dapr.io/v2alpha1
kind: Subscription
metadata:
  name: 'order'
spec:
  topic: 'test'
  routes:
    default: '/test'
  pubsubname: 'mypubsub'
scopes:
- outboxtest
`))
	return []framework.Option{
		framework.WithProcesses(app, o.daprd),
	}
}

func (o *delete) Run(t *testing.T, ctx context.Context) {
	o.daprd.WaitUntilRunning(t, ctx)
	transactionRequest, err := json.Marshal(&stateTransactionRequestBody{
		Operations: []stateTransactionRequestBodyOperation{
			{
				Operation: "upsert",
				Request: state.SetRequest{
					Key:      "1",
					Value:    "2",
					Metadata: map[string]string{"outbox.cloudevent.myapp": "myapp1", "data": "a", "id": "b"},
				},
			},
		},
	})
	require.NoError(t, err)
	client := client.HTTP(t)
	postURL := fmt.Sprintf("http://%s/v1.0/state/mystore/transaction", o.daprd.HTTPAddress())
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, postURL, bytes.NewReader(transactionRequest))
	require.NoError(t, err)
	resp, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())
	assert.Empty(t, string(body))
	assert.EventuallyWithT(t, func(c *assert.CollectT) {
		ce, ok := o.msg.Load().(map[string]string)
		if !assert.True(c, ok) {
			return
		}
		assert.Equal(c, "2", ce["data"])
		assert.Equal(c, "myapp1", ce["outbox.cloudevent.myapp"])
		assert.Contains(c, ce["id"], "outbox-")
	}, time.Second*10, time.Millisecond*10)
	assert.Equal(t, int32(1), o.appTestCalled.Load())
	transactionRequest, err = json.Marshal(&stateTransactionRequestBody{
		Operations: []stateTransactionRequestBodyOperation{
			{
				Operation: "delete",
				Request: state.DeleteRequest{
					Key:      "1",
					Metadata: map[string]string{"outbox.cloudevent.myapp": "myapp2", "data": "a", "id": "b"},
				},
			},
		},
	})
	require.NoError(t, err)
	req, err = http.NewRequestWithContext(ctx, http.MethodPost, postURL, bytes.NewReader(transactionRequest))
	require.NoError(t, err)
	resp, err = client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	body, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())
	assert.Empty(t, string(body))
	//
	assert.EventuallyWithT(t, func(c *assert.CollectT) {
		ce, ok := o.msg.Load().(map[string]string)
		if !assert.True(c, ok) {
			return
		}
		assert.Equal(c, "2", ce["data"])
		assert.Equal(c, "myapp1", ce["outbox.cloudevent.myapp"])
		assert.Contains(c, ce["id"], "outbox-")
	}, time.Second*10, time.Millisecond*10)
	// expecting the delete to not call the test app topic and calls counter stay the same
	assert.Equal(t, int32(1), o.appTestCalled.Load())
}
````

### FILE: `dapr_official_transactional_outbox/upstream/runtime/tests/integration/suite/daprd/outbox/http/projection.go`
```yaml
block_id: "DAPR-OFFICIAL-OUTBOX:19:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/dapr/dapr/2dcf2e3548faf3c740b3d29ee300e2132db65cff/tests/integration/suite/daprd/outbox/http/projection.go"
license: "Apache-2.0"
sha256: "4e49a51289c0a2f809f496ef7a78e72e963a3d2b5211cbcdcfd246101392c216"
variables: []
secrets_allowed: false
```
````go
/*
Copyright 2023 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dapr/components-contrib/state"
	"github.com/dapr/dapr/tests/integration/framework"
	"github.com/dapr/dapr/tests/integration/framework/client"
	procdaprd "github.com/dapr/dapr/tests/integration/framework/process/daprd"
	prochttp "github.com/dapr/dapr/tests/integration/framework/process/http"
	"github.com/dapr/dapr/tests/integration/suite"
)

func init() {
	suite.Register(new(projection))
}

type projection struct {
	daprd *procdaprd.Daprd
}

func (o *projection) Setup(t *testing.T) []framework.Option {
	newHTTPServer := func() *prochttp.HTTP {
		handler := http.NewServeMux()
		var msg atomic.Value

		handler.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			b, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatal(err)
			}

			msg.Store(b)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		})

		handler.HandleFunc("/getValue", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			m := msg.Load()
			if m == nil {
				return
			}
			w.Write(msg.Load().([]byte))
		})

		return prochttp.New(t, prochttp.WithHandler(handler))
	}
	srv1 := newHTTPServer()

	o.daprd = procdaprd.New(t, procdaprd.WithAppID("outboxtest"), procdaprd.WithAppPort(srv1.Port()), procdaprd.WithResourceFiles(`
apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
  name: mystore
spec:
  type: state.in-memory
  version: v1
  metadata:
  - name: outboxPublishPubsub
    value: "mypubsub"
  - name: outboxPublishTopic
    value: "test"
`,
		`
apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
  name: 'mypubsub'
spec:
  type: pubsub.in-memory
  version: v1
`,
		`
apiVersion: dapr.io/v2alpha1
kind: Subscription
metadata:
  name: 'order'
spec:
  topic: 'test'
  routes:
    default: '/test'
  pubsubname: 'mypubsub'
scopes:
- outboxtest
`))

	return []framework.Option{
		framework.WithProcesses(srv1, o.daprd),
	}
}

func (o *projection) Run(t *testing.T, ctx context.Context) {
	o.daprd.WaitUntilRunning(t, ctx)

	postURL := fmt.Sprintf("http://localhost:%d/v1.0/state/mystore/transaction", o.daprd.HTTPPort())
	stateReq := state.SetRequest{
		Key:   "1",
		Value: "2",
	}
	projectionRequest := state.SetRequest{
		Key:      "1",
		Value:    "3",
		Metadata: map[string]string{"outbox.projection": "true"},
	}

	tr := stateTransactionRequestBody{
		Operations: []stateTransactionRequestBodyOperation{
			{
				Operation: "upsert",
				Request:   stateReq,
			},
			{
				Operation: "upsert",
				Request:   projectionRequest,
			},
		},
	}

	b, err := json.Marshal(&tr)
	require.NoError(t, err)

	httpClient := client.HTTP(t)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, postURL, bytes.NewReader(b))
	require.NoError(t, err)
	resp, err := httpClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())
	assert.Empty(t, string(body))

	assert.EventuallyWithT(t, func(c *assert.CollectT) {
		// validate projection data is reflected in final publish
		req, err = http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("http://localhost:%v/getValue", o.daprd.AppPort(t)), nil)
		require.NoError(c, err)
		resp, err = httpClient.Do(req)
		require.NoError(c, err)
		t.Cleanup(func() {
			require.NoError(t, resp.Body.Close())
		})
		body, err = io.ReadAll(resp.Body)
		require.NoError(c, err)

		var ce map[string]string
		err = json.Unmarshal(body, &ce)
		assert.NoError(c, err)
		assert.Equal(c, "3", ce["data"])
	}, time.Second*10, time.Millisecond*10)

	assert.EventuallyWithT(t, func(c *assert.CollectT) {
		// validate correct state is in the db and not the projection data
		req, err = http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("http://localhost:%d/v1.0/state/mystore/1", o.daprd.HTTPPort()), nil)
		require.NoError(c, err)
		resp, err = httpClient.Do(req)
		require.NoError(c, err)
		t.Cleanup(func() {
			require.NoError(t, resp.Body.Close())
		})
		body, err = io.ReadAll(resp.Body)
		require.NoError(c, err)

		val, err := strconv.Unquote(string(body))
		require.NoError(c, err)

		require.NoError(c, err)
		assert.Equal(c, "2", val)
	}, time.Second*10, time.Millisecond*10)
}
````

### FILE: `dapr_official_transactional_outbox/upstream/runtime/tests/integration/suite/daprd/outbox/grpc/basic.go`
```yaml
block_id: "DAPR-OFFICIAL-OUTBOX:20:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/dapr/dapr/2dcf2e3548faf3c740b3d29ee300e2132db65cff/tests/integration/suite/daprd/outbox/grpc/basic.go"
license: "Apache-2.0"
sha256: "42e5524ebe243896d893a00a86afe2c9846a8642f6bf6b9c514467613d05eab3"
variables: []
secrets_allowed: false
```
````go
/*
Copyright 2023 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package grpc

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/dapr/dapr/pkg/proto/common/v1"
	runtimev1pb "github.com/dapr/dapr/pkg/proto/runtime/v1"
	"github.com/dapr/dapr/tests/integration/framework"
	"github.com/dapr/dapr/tests/integration/framework/process/daprd"
	"github.com/dapr/dapr/tests/integration/framework/process/grpc/app"
	"github.com/dapr/dapr/tests/integration/suite"
)

func init() {
	suite.Register(new(basic))
}

type basic struct {
	eventCalled atomic.Int32
	daprd       *daprd.Daprd
	lock        sync.Mutex
	msg         []byte
}

func (o *basic) Setup(t *testing.T) []framework.Option {
	onTopicEvent := func(ctx context.Context, in *runtimev1pb.TopicEventRequest) (*runtimev1pb.TopicEventResponse, error) {
		o.lock.Lock()
		defer o.lock.Unlock()
		o.eventCalled.Add(1)
		o.msg = in.GetData()
		return &runtimev1pb.TopicEventResponse{
			Status: runtimev1pb.TopicEventResponse_SUCCESS,
		}, nil
	}

	srv1 := app.New(t, app.WithOnTopicEventFn(onTopicEvent))
	o.daprd = daprd.New(t, daprd.WithAppID("outboxtest"), daprd.WithAppPort(srv1.Port(t)), daprd.WithAppProtocol("grpc"), daprd.WithResourceFiles(`
apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
  name: mystore
spec:
  type: state.in-memory
  version: v1
  metadata:
  - name: outboxPublishPubsub
    value: "mypubsub"
  - name: outboxPublishTopic
    value: "test"
`,
		`
apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
  name: 'mypubsub'
spec:
  type: pubsub.in-memory
  version: v1
`,
		`
apiVersion: dapr.io/v2alpha1
kind: Subscription
metadata:
  name: 'order'
spec:
  topic: 'test'
  routes:
    default: '/test'
  pubsubname: 'mypubsub'
scopes:
- outboxtest
`))

	return []framework.Option{
		framework.WithProcesses(srv1, o.daprd),
	}
}

func (o *basic) Run(t *testing.T, ctx context.Context) {
	o.daprd.WaitUntilRunning(t, ctx)

	conn, err := grpc.DialContext(ctx, o.daprd.GRPCAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock()) //nolint:staticcheck
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close()) })

	_, err = runtimev1pb.NewDaprClient(conn).ExecuteStateTransaction(ctx, &runtimev1pb.ExecuteStateTransactionRequest{
		StoreName: "mystore",
		Operations: []*runtimev1pb.TransactionalStateOperation{
			{
				OperationType: "upsert",
				Request: &common.StateItem{
					Key:   "1",
					Value: []byte("2"),
				},
			},
		},
	})
	require.NoError(t, err)

	assert.Eventually(t, func() bool {
		o.lock.Lock()
		defer o.lock.Unlock()
		return string(o.msg) == "2"
	}, time.Second*10, time.Millisecond*10, "failed to receive message in time")

	assert.Equal(t, int32(1), o.eventCalled.Load())
}
````

### FILE: `dapr_official_transactional_outbox/upstream/runtime/tests/integration/suite/daprd/outbox/grpc/projection.go`
```yaml
block_id: "DAPR-OFFICIAL-OUTBOX:21:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/dapr/dapr/2dcf2e3548faf3c740b3d29ee300e2132db65cff/tests/integration/suite/daprd/outbox/grpc/projection.go"
license: "Apache-2.0"
sha256: "45849b9bc17678967886a3249b861097d702a14d5e99aa6bbb468d055b2d31bf"
variables: []
secrets_allowed: false
```
````go
/*
Copyright 2023 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package grpc

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/dapr/dapr/pkg/proto/common/v1"
	runtimev1pb "github.com/dapr/dapr/pkg/proto/runtime/v1"
	"github.com/dapr/dapr/tests/integration/framework"
	"github.com/dapr/dapr/tests/integration/framework/process/daprd"
	"github.com/dapr/dapr/tests/integration/framework/process/grpc/app"
	"github.com/dapr/dapr/tests/integration/suite"
)

func init() {
	suite.Register(new(projection))
}

type projection struct {
	daprd *daprd.Daprd
	lock  sync.Mutex
	msg   []byte
}

func (o *projection) Setup(t *testing.T) []framework.Option {
	onTopicEvent := func(ctx context.Context, in *runtimev1pb.TopicEventRequest) (*runtimev1pb.TopicEventResponse, error) {
		o.lock.Lock()
		defer o.lock.Unlock()
		o.msg = in.GetData()
		return &runtimev1pb.TopicEventResponse{
			Status: runtimev1pb.TopicEventResponse_SUCCESS,
		}, nil
	}

	srv1 := app.New(t, app.WithOnTopicEventFn(onTopicEvent))
	o.daprd = daprd.New(t, daprd.WithAppID("outboxtest"), daprd.WithAppPort(srv1.Port(t)), daprd.WithAppProtocol("grpc"), daprd.WithResourceFiles(`
apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
  name: mystore
spec:
  type: state.in-memory
  version: v1
  metadata:
  - name: outboxPublishPubsub
    value: "mypubsub"
  - name: outboxPublishTopic
    value: "test"
`,
		`
apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
  name: 'mypubsub'
spec:
  type: pubsub.in-memory
  version: v1
`,
		`
apiVersion: dapr.io/v2alpha1
kind: Subscription
metadata:
  name: 'order'
spec:
  topic: 'test'
  routes:
    default: '/test'
  pubsubname: 'mypubsub'
scopes:
- outboxtest
`))

	return []framework.Option{
		framework.WithProcesses(srv1, o.daprd),
	}
}

func (o *projection) Run(t *testing.T, ctx context.Context) {
	o.daprd.WaitUntilRunning(t, ctx)
	//nolint:staticcheck
	conn, err := grpc.DialContext(ctx, o.daprd.GRPCAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close()) })

	_, err = runtimev1pb.NewDaprClient(conn).ExecuteStateTransaction(ctx, &runtimev1pb.ExecuteStateTransactionRequest{
		StoreName: "mystore",
		Operations: []*runtimev1pb.TransactionalStateOperation{
			{
				OperationType: "upsert",
				Request: &common.StateItem{
					Key:   "1",
					Value: []byte("2"),
				},
			},
			{
				OperationType: "upsert",
				Request: &common.StateItem{
					Key:   "1",
					Value: []byte("3"),
					Metadata: map[string]string{
						"outbox.projection": "true",
					},
				},
			},
		},
	})
	require.NoError(t, err)

	assert.Eventually(t, func() bool {
		o.lock.Lock()
		defer o.lock.Unlock()
		return string(o.msg) == "3"
	}, time.Second*5, time.Millisecond*10, "failed to receive message in time")

	assert.Eventually(t, func() bool {
		o.lock.Lock()
		defer o.lock.Unlock()

		resp, err := runtimev1pb.NewDaprClient(conn).GetState(ctx, &runtimev1pb.GetStateRequest{
			Key:       "1",
			StoreName: "mystore",
		})
		require.NoError(t, err)
		return string(resp.GetData()) == "2"
	}, time.Second*5, time.Millisecond*10, "failed to receive message in time")
}
````

## 6. Configuration surface

- `project-selection.template.json`: decisión humana/proyecto; no admite secretos.
- `stateStoreComponent`, `pubsubComponent`, `outboxPublishTopic`: nombres reales definidos por el proyecto.
- Evidencias obligatorias: transacciones/outbox del store, cuenta broker, duplicados, idempotencia, recuperación y rollback.
- `goSdkSelected` debe permanecer `false` hasta una corrección oficial revalidada.

## 7. Dependency bill

| Dependencia | Identidad | Licencia | Condición |
|---|---|---|---|
| Dapr Runtime | v1.18.3 / `2dcf2e3548faf3c740b3d29ee300e2132db65cff` | Apache-2.0 | archive/hash/grafo/tests/scan fechados; plataforma opt-in |
| Dapr Docs | v1.18 / `5958a7e19a04e199326a6c5321bdd7714ee83b4c` | Apache-2.0 | documentación outbox exacta |
| PowerShell | 7+ | runtime externo | gates y adquisición local |
| Dapr Go SDK | 1.14.2 y main observado | Apache-2.0 | RECHAZADO por vulnerabilidades alcanzables; no materializado |

## 8. Apply order

1. Materializar y ejecutar ambos tests offline.
2. Completar y validar selección del proyecto; un template intacto debe fallar.
3. Adquirir el runtime completo desde cache aprobada o red autorizada.
4. Repetir tests oficiales y scan con la identidad exacta.
5. Configurar componentes oficiales elegidos por el proyecto y probar duplicados/recuperación.
6. Autorizar producción sólo con rollback demostrado. Para revertir, detener selección Dapr y restaurar el camino previo probado; no borrar evidencia.

## 9. Verification

```powershell
pwsh ./dapr_official_transactional_outbox/test_embedded_sources.ps1
pwsh ./dapr_official_transactional_outbox/test_acquisition.ps1
pwsh ./dapr_official_transactional_outbox/validate_project_selection.ps1 -SelectionPath <completed-selection.json>
```

En la evidencia V1, el runtime completo pasó `go mod verify`, fake tests y seis tests principales/20 subtests outbox; govulncheck focal no encontró vulnerabilidades alcanzables. La integración real y la selección siguen condicionadas al proyecto. El verificador debe rechazar cambios de bytes, template incompleto, SDK Go vulnerable, red no autorizada y destino existente.

## 10. Reconstruction evidence

- Source archive, commits, firmas, hashes, grafo, tests y scan: `reconstruction_evidence/DAPR_OFFICIAL_TRANSACTIONAL_OUTBOX_GO_2026-08-28_V1.md`.
- Auditoría global posterior a materialización: evidencia total V53.
- Las condiciones `UP-FAIL-125/126` y los fallos locales V56 permanecen en el ledger; un release nuevo no hereda PASS.
