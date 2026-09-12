# Microsoft Kiota OpenAPI Client Gate

## 1. Metadata

```yaml
pack_id: "MICROSOFT-KIOTA-OPENAPI-CLIENT-GATE"
pack_version: "0.1.0"
status:
  authority: ELITE_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un gate Go para adquirir Microsoft Kiota 1.35.0 autocontenido por identidad exacta, generar dos veces desde un OpenAPI local hash-locked sin sobrescribir, probar determinismo, compilar, vetear y ejecutar SCA; no afirma comportamiento de proveedor ni producción."
stacks: ["Microsoft Kiota 1.35.0", "Go 1.26.7", "Google OSV-Scanner 2.5.1", "PowerShell 7"]
compatible_with: ["GO-PROVIDER-INTEGRATION-CORE 0.1.x", "DEPENDENCY-LICENSE-EVIDENCE-CORE 0.1.x", "PORTABLE-CI-GATE-RUNNER 0.1.x"]
incompatible_with: ["OpenAPI remoto o sin hash aprobado", "destino preexistente", "autenticación o reglas de negocio inferidas", "branches o releases móviles"]
license_expression: "LicenseRef-Workspace-Owner AND MIT"
upstream_sources: ["https://github.com/microsoft/kiota/tree/114aa7ee609262d892fd9ceb02b2d9f7ecb84190", "https://github.com/microsoft/kiota/releases/tag/v1.35.0"]
verified_at: "2026-09-05"
```

## 2. Applicability

Use cuando un proyecto tenga un contrato OpenAPI confiable e inmutable y necesite eliminar escritura manual del transporte Go. El contrato y su SHA-256 deben pertenecer al lock del proyecto. La salida se genera siempre en un destino ausente, produce receipt por archivo y sólo puede incorporarse tras revisión, grafo congelado, compilación, contract tests contra sandbox y los gates de seguridad del proyecto.

No use el fixture oficial de Microsoft como modelo comercial. No presente el glue `AUTHORED` como código de Microsoft. Otros lenguajes expuestos por Kiota permanecen fuera del claim verificado de este pack.

## 3. Architecture contract

El acquisition runner descarga el asset oficial Windows x64, valida longitud/SHA-256, inventario ZIP seguro, ejecutable interno y versión `release+commit`. El generator runner exige un archivo local y su hash, deshabilita telemetría, rechaza destinos/receipts existentes, invoca sólo la lane Go verificada y registra manifest y hash del árbol resultante.

El verifier usa el `ToDoApi.yaml` del propio commit Microsoft únicamente como fixture de generador. Ejecuta dos generaciones independientes, exige el árbol exacto fijado, prueba dos casos negativos, aplica `go.mod`/`go.sum` congelados, compila con `go test`, ejecuta `go vet` y consulta OSV. Un finding futuro falla el gate y reabre la admisión. Ninguno de estos pasos sustituye auth, idempotencia, reconciliación o pruebas live del proveedor.

## 4. Exact file manifest

```text
CREATE microsoft_kiota_openapi_gate/source-lock.json
CREATE microsoft_kiota_openapi_gate/acquire_kiota.ps1
CREATE microsoft_kiota_openapi_gate/generate_go_client.ps1
CREATE microsoft_kiota_openapi_gate/fixture-go.mod
CREATE microsoft_kiota_openapi_gate/fixture-go.sum
CREATE microsoft_kiota_openapi_gate/verify_contract.ps1
CREATE microsoft_kiota_openapi_gate/README.md
```

## 5. Materialization blocks

### FILE: `microsoft_kiota_openapi_gate/source-lock.json`
```yaml
block_id: "MICROSOFT-KIOTA-OPENAPI-CLIENT-GATE:source-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite receipt derived from the exact Microsoft release, commit API, source archive and binary asset"
license: "LicenseRef-Workspace-Owner"
sha256: "4de730734fa2be27cb47a2e834783e336162b35ed77eee7ae9e02417a8017490"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-microsoft-kiota-openapi-gate-lock/v1",
  "verified_at": "2026-09-05",
  "source": {
    "owner": "Microsoft",
    "repository": "microsoft/kiota",
    "release": "v1.35.0",
    "commit": "114aa7ee609262d892fd9ceb02b2d9f7ecb84190",
    "commit_signature_verified": true,
    "archive_url": "https://github.com/microsoft/kiota/archive/114aa7ee609262d892fd9ceb02b2d9f7ecb84190.zip",
    "archive_bytes": 3580264,
    "archive_sha256": "9a9856310e11e5c0877b8b3a732dee78e051d67e5809875c560b46e6c5541248",
    "tree_files": 1073,
    "tree_bytes": 8798877,
    "test_files": 267,
    "global_json_sha256": "8d0d9990a2bfd2c1cc9c6552ef6eb179e397c0ff873f1b46dfae515984b1a468",
    "license_expression": "MIT",
    "license_path": "LICENSE",
    "license_sha256": "c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383"
  },
  "windows_x64": {
    "asset_id": 543519503,
    "asset_name": "win-x64.zip",
    "url": "https://github.com/microsoft/kiota/releases/download/v1.35.0/win-x64.zip",
    "bytes": 37140722,
    "sha256": "66b5547b948f7be724fa5e0ddddebda8f7b662574ae58eb9d3b8c51c609c6271",
    "executable": "kiota.exe",
    "executable_bytes": 89879984,
    "executable_sha256": "dffa90d51f5068c0fdd51dd9b539bce2a4c028335c09f449d9b2aa1f5f3d5d71",
    "version_output": "1.35.0+114aa7ee609262d892fd9ceb02b2d9f7ecb84190",
    "archive_entries": ["appsettings.json", "Kiota.Builder.pdb", "kiota.exe", "kiota.pdb"]
  },
  "official_fixture": {
    "path": "tests/Kiota.Builder.IntegrationTests/ToDoApi.yaml",
    "url": "https://raw.githubusercontent.com/microsoft/kiota/114aa7ee609262d892fd9ceb02b2d9f7ecb84190/tests/Kiota.Builder.IntegrationTests/ToDoApi.yaml",
    "bytes": 1655,
    "sha256": "9141520b53b6f9fcf544006290e912f7e02ffee9eab79b4e9907d694b8b18b49"
  },
  "verified_lane": {
    "language": "Go",
    "class_name": "TodoClient",
    "namespace_name": "github.com/elite/kiota_todo_fixture",
    "exclude_backward_compatible": true,
    "generated_files_including_lock": 7,
    "generated_code_files": 6,
    "generated_code_bytes": 28750,
    "generated_code_tree_sha256": "d6dcb4092c1b83ecd84ae864729f9a0df4ab1114aca8b3f54508d1a4d3cd2a53",
    "kiota_description_hash": "05867D3B5A2BC51B99A158371D67447B7C0605773D59DC20D0F7B7FE07D404DE6599D2F1BD22F500A6799B8B53189AEF32318C0B716B22F360D8FF53C3FDEE0C",
    "go_version": "1.26.7",
    "go_mod_sha256": "f2469af1de3e7b16a4e532c555e2189549844b0757c07fa991e124cfeb21125b",
    "go_sum_sha256": "4f03b975309f179ffff292cec6df43e37aaa8cb21351d8fc7c691763ebd7c444",
    "osv_scanner_version": "2.5.1",
    "osv_packages": 14,
    "osv_findings_at_verification": 0
  }
}
````

### FILE: `microsoft_kiota_openapi_gate/acquire_kiota.ps1`
```yaml
block_id: "MICROSOFT-KIOTA-OPENAPI-CLIENT-GATE:acquire:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite fail-closed acquisition runner for the exact Microsoft asset"
license: "LicenseRef-Workspace-Owner"
sha256: "21cfec4cdcd2c43dd48a9858c083caf26eff28f232f56303dcc32924232ba266"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0

[CmdletBinding()]
param(
  [Parameter(Mandatory = $true)]
  [string] $Destination
)

$ErrorActionPreference = 'Stop'
$lockPath = Join-Path $PSScriptRoot 'source-lock.json'
$lock = Get-Content -LiteralPath $lockPath -Raw | ConvertFrom-Json -Depth 30
$destinationPath = [IO.Path]::GetFullPath($Destination)
$parent = Split-Path -Parent $destinationPath
if (-not $parent) { throw 'Destination requires a parent directory' }
if (Test-Path -LiteralPath $destinationPath) {
  $exe = Join-Path $destinationPath $lock.windows_x64.executable
  if (-not (Test-Path -LiteralPath $exe -PathType Leaf)) { throw "Existing Kiota destination is incomplete: $destinationPath" }
  $hash = (Get-FileHash -LiteralPath $exe -Algorithm SHA256).Hash.ToLowerInvariant()
  $version = (& $exe --version 2>&1 | Out-String).Trim()
  if ($LASTEXITCODE -ne 0 -or $hash -cne $lock.windows_x64.executable_sha256 -or $version -cne $lock.windows_x64.version_output) {
    throw "Existing Kiota destination does not match the lock: $destinationPath"
  }
  Write-Output "KIOTA_ACQUISITION_REUSED version=$version sha256=$hash"
  exit 0
}
if (-not (Test-Path -LiteralPath $parent)) { [void](New-Item -ItemType Directory -Path $parent) }

$stage = Join-Path ([IO.Path]::GetTempPath()) ('elite-kiota-' + [Guid]::NewGuid().ToString('N'))
$archive = Join-Path $stage $lock.windows_x64.asset_name
$expanded = Join-Path $stage 'expanded'
[void](New-Item -ItemType Directory -Path $stage)
Invoke-WebRequest -Uri $lock.windows_x64.url -OutFile $archive
$archiveItem = Get-Item -LiteralPath $archive
$archiveHash = (Get-FileHash -LiteralPath $archive -Algorithm SHA256).Hash.ToLowerInvariant()
if ($archiveItem.Length -ne [int64]$lock.windows_x64.bytes -or $archiveHash -cne $lock.windows_x64.sha256) {
  throw "Kiota archive identity mismatch: bytes=$($archiveItem.Length) sha256=$archiveHash"
}

$zip = [IO.Compression.ZipFile]::OpenRead($archive)
try {
  $actualEntries = @($zip.Entries | Where-Object { -not [string]::IsNullOrEmpty($_.Name) } | ForEach-Object { $_.FullName.Replace('\', '/') } | Sort-Object)
  $expectedEntries = @($lock.windows_x64.archive_entries | Sort-Object)
  if (($actualEntries -join "`n") -cne ($expectedEntries -join "`n")) { throw 'Kiota archive entry inventory drifted' }
  foreach ($entry in $zip.Entries) {
    $normalized = $entry.FullName.Replace('\', '/')
    if ([IO.Path]::IsPathRooted($normalized) -or $normalized.Split('/') -contains '..') { throw "Unsafe Kiota archive entry: $normalized" }
  }
} finally { $zip.Dispose() }

[IO.Compression.ZipFile]::ExtractToDirectory($archive, $expanded)
$exe = Join-Path $expanded $lock.windows_x64.executable
if ((Get-Item -LiteralPath $exe).Length -ne [int64]$lock.windows_x64.executable_bytes) { throw 'Kiota executable length mismatch' }
$exeHash = (Get-FileHash -LiteralPath $exe -Algorithm SHA256).Hash.ToLowerInvariant()
if ($exeHash -cne $lock.windows_x64.executable_sha256) { throw 'Kiota executable SHA-256 mismatch' }
$env:KIOTA_CLI_TELEMETRY_OPTOUT = 'true'
$version = (& $exe --version 2>&1 | Out-String).Trim()
if ($LASTEXITCODE -ne 0 -or $version -cne $lock.windows_x64.version_output) { throw "Kiota version mismatch: $version" }
Move-Item -LiteralPath $expanded -Destination $destinationPath
Write-Output "KIOTA_ACQUISITION_PASS version=$version archive_sha256=$archiveHash executable_sha256=$exeHash"
````

### FILE: `microsoft_kiota_openapi_gate/generate_go_client.ps1`
```yaml
block_id: "MICROSOFT-KIOTA-OPENAPI-CLIENT-GATE:generate-go:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite hash-locked non-overwriting wrapper around the official Kiota CLI"
license: "LicenseRef-Workspace-Owner"
sha256: "3c1eae1db3f8a26cb70c6bab078839152fb95671148609c79cb947c987f0e842"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0

[CmdletBinding()]
param(
  [Parameter(Mandatory = $true)] [string] $KiotaExecutable,
  [Parameter(Mandatory = $true)] [string] $OpenApi,
  [Parameter(Mandatory = $true)] [ValidatePattern('^[0-9a-f]{64}$')] [string] $ExpectedOpenApiSha256,
  [Parameter(Mandatory = $true)] [string] $Output,
  [Parameter(Mandatory = $true)] [ValidatePattern('^[A-Za-z][A-Za-z0-9]*$')] [string] $ClassName,
  [Parameter(Mandatory = $true)] [ValidatePattern('^[A-Za-z0-9._/-]+$')] [string] $NamespaceName,
  [Parameter(Mandatory = $true)] [string] $ReceiptPath
)

$ErrorActionPreference = 'Stop'
$exe = (Resolve-Path -LiteralPath $KiotaExecutable).Path
$contract = (Resolve-Path -LiteralPath $OpenApi).Path
$outputPath = [IO.Path]::GetFullPath($Output)
$receipt = [IO.Path]::GetFullPath($ReceiptPath)
if (Test-Path -LiteralPath $outputPath) { throw "Output must not exist: $outputPath" }
if (Test-Path -LiteralPath $receipt) { throw "Receipt must not exist: $receipt" }
if ($receipt.StartsWith($outputPath + [IO.Path]::DirectorySeparatorChar, [StringComparison]::OrdinalIgnoreCase)) {
  throw 'Receipt must remain outside generated output'
}
$contractHash = (Get-FileHash -LiteralPath $contract -Algorithm SHA256).Hash.ToLowerInvariant()
if ($contractHash -cne $ExpectedOpenApiSha256) { throw "OpenAPI SHA-256 mismatch: expected=$ExpectedOpenApiSha256 actual=$contractHash" }
$exeHash = (Get-FileHash -LiteralPath $exe -Algorithm SHA256).Hash.ToLowerInvariant()
$env:KIOTA_CLI_TELEMETRY_OPTOUT = 'true'
$version = (& $exe --version 2>&1 | Out-String).Trim()
if ($LASTEXITCODE -ne 0) { throw 'Kiota version probe failed' }
& $exe generate --openapi $contract --language Go --output $outputPath --class-name $ClassName --namespace-name $NamespaceName --exclude-backward-compatible --log-level Warning
if ($LASTEXITCODE -ne 0) { throw "Kiota generation failed: exit=$LASTEXITCODE" }
if (-not (Test-Path -LiteralPath (Join-Path $outputPath 'kiota-lock.json') -PathType Leaf)) { throw 'Kiota generation did not create kiota-lock.json' }
$reparse = @(Get-ChildItem -LiteralPath $outputPath -Recurse -Force | Where-Object { $_.Attributes -band [IO.FileAttributes]::ReparsePoint })
if ($reparse.Count -ne 0) { throw 'Generated output contains reparse points' }
$files = @(Get-ChildItem -LiteralPath $outputPath -Recurse -File | ForEach-Object {
  [pscustomobject]@{
    path = $_.FullName.Substring($outputPath.Length + 1).Replace('\', '/')
    bytes = $_.Length
    sha256 = (Get-FileHash -LiteralPath $_.FullName -Algorithm SHA256).Hash.ToLowerInvariant()
  }
} | Sort-Object path)
if ($files.Count -eq 0) { throw 'Kiota generated no files' }
$manifest = ($files | ForEach-Object { "$($_.sha256) $($_.bytes) $($_.path)" }) -join "`n"
$treeHash = [Convert]::ToHexString([Security.Cryptography.SHA256]::HashData([Text.Encoding]::UTF8.GetBytes($manifest))).ToLowerInvariant()
$codeFiles = @($files | Where-Object path -CNE 'kiota-lock.json')
if ($codeFiles.Count -eq 0) { throw 'Kiota generated no code files' }
$codeManifest = ($codeFiles | ForEach-Object { "$($_.sha256) $($_.bytes) $($_.path)" }) -join "`n"
$codeTreeHash = [Convert]::ToHexString([Security.Cryptography.SHA256]::HashData([Text.Encoding]::UTF8.GetBytes($codeManifest))).ToLowerInvariant()
$kiotaLockPath = Join-Path $outputPath 'kiota-lock.json'
$kiotaLock = Get-Content -LiteralPath $kiotaLockPath -Raw | ConvertFrom-Json -Depth 20
$resolvedDescription = [IO.Path]::GetFullPath((Join-Path $outputPath $kiotaLock.descriptionLocation))
if ($resolvedDescription -cne $contract) { throw 'kiota-lock.json descriptionLocation does not resolve to the approved contract' }
if ($kiotaLock.kiotaVersion -cne ([regex]::Match($version, '^\d+\.\d+\.\d+').Value) -or $kiotaLock.clientClassName -cne $ClassName -or $kiotaLock.clientNamespaceName -cne $NamespaceName -or $kiotaLock.language -cne 'Go' -or $kiotaLock.excludeBackwardCompatible -ne $true -or $kiotaLock.disableSSLValidation -ne $false) {
  throw 'kiota-lock.json semantic configuration differs from the requested safe lane'
}
if ($kiotaLock.descriptionHash -notmatch '^[0-9A-F]{128}$') { throw 'kiota-lock.json descriptionHash is malformed' }
$record = [ordered]@{
  schema = 'elite-kiota-generation-receipt/v1'
  generated_at_utc = [DateTimeOffset]::UtcNow.ToString('O')
  generator = [ordered]@{ version = $version; executable_sha256 = $exeHash }
  input = [ordered]@{ path = $contract; sha256 = $contractHash }
  configuration = [ordered]@{ language = 'Go'; class_name = $ClassName; namespace_name = $NamespaceName; exclude_backward_compatible = $true }
  output = [ordered]@{
    path = $outputPath
    files = $files.Count
    bytes = [int64](($files | Measure-Object bytes -Sum).Sum)
    tree_sha256 = $treeHash
    code_files = $codeFiles.Count
    code_bytes = [int64](($codeFiles | Measure-Object bytes -Sum).Sum)
    code_tree_sha256 = $codeTreeHash
    kiota_lock_sha256 = (Get-FileHash -LiteralPath $kiotaLockPath -Algorithm SHA256).Hash.ToLowerInvariant()
    kiota_description_hash = $kiotaLock.descriptionHash
    manifest = $files
  }
}
$receiptParent = Split-Path -Parent $receipt
if (-not (Test-Path -LiteralPath $receiptParent)) { [void](New-Item -ItemType Directory -Path $receiptParent) }
[IO.File]::WriteAllText($receipt, ($record | ConvertTo-Json -Depth 10), [Text.UTF8Encoding]::new($false))
Write-Output "KIOTA_GENERATION_PASS version=$version files=$($files.Count) code_files=$($codeFiles.Count) code_tree_sha256=$codeTreeHash"
````

### FILE: `microsoft_kiota_openapi_gate/fixture-go.mod`
```yaml
block_id: "MICROSOFT-KIOTA-OPENAPI-CLIENT-GATE:fixture-go-mod:v1"
operation: CREATE
provenance: AUTHORED
source: "Frozen module graph resolved from Kiota 1.35.0 generated Go imports"
license: "LicenseRef-Workspace-Owner"
sha256: "f2469af1de3e7b16a4e532c555e2189549844b0757c07fa991e124cfeb21125b"
variables: []
secrets_allowed: false
```
````go.mod
module github.com/elite/kiota_todo_fixture

go 1.26.7

require (
	github.com/microsoft/kiota-abstractions-go v1.10.1
	github.com/microsoft/kiota-serialization-form-go v1.1.3
	github.com/microsoft/kiota-serialization-json-go v1.1.4
	github.com/microsoft/kiota-serialization-multipart-go v1.1.2
	github.com/microsoft/kiota-serialization-text-go v1.1.3
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-logr/logr v1.4.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/std-uritemplate/std-uritemplate/go/v2 v2.0.12 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel v1.45.0 // indirect
	go.opentelemetry.io/otel/metric v1.45.0 // indirect
	go.opentelemetry.io/otel/trace v1.45.0 // indirect
)
````

### FILE: `microsoft_kiota_openapi_gate/fixture-go.sum`
```yaml
block_id: "MICROSOFT-KIOTA-OPENAPI-CLIENT-GATE:fixture-go-sum:v1"
operation: CREATE
provenance: AUTHORED
source: "Go module checksums for the frozen Kiota-generated fixture graph"
license: "LicenseRef-Workspace-Owner"
sha256: "4f03b975309f179ffff292cec6df43e37aaa8cb21351d8fc7c691763ebd7c444"
variables: []
secrets_allowed: false
```
````text
github.com/cespare/xxhash/v2 v2.3.0 h1:UL815xU9SqsFlibzuggzjXhog7bL6oX9BbNZnL2UFvs=
github.com/cespare/xxhash/v2 v2.3.0/go.mod h1:VGX0DQ3Q6kWi7AoAeZDth3/j3BFtOZR5XLFGgcrjCOs=
github.com/go-logr/logr v1.2.2/go.mod h1:jdQByPbusPIv2/zmleS9BjJVeZ6kBagPoEUsqbVz/1A=
github.com/go-logr/logr v1.4.4 h1:tG4xh9yMsRCAiodLVTxyrkzSZ9+o0L1Kg/+cPVcbP/8=
github.com/go-logr/logr v1.4.4/go.mod h1:9T104GzyrTigFIr8wt5mBrctHMim0Nb2HLGrmQ40KvY=
github.com/go-logr/stdr v1.2.2 h1:hSWxHoqTgW2S2qGc0LTAI563KZ5YKYRhT3MFKZMbjag=
github.com/go-logr/stdr v1.2.2/go.mod h1:mMo/vtBO5dYbehREoey6XUKy/eSumjCCveDpRre4VKE=
github.com/google/go-cmp v0.7.0 h1:wk8382ETsv4JYUZwIsn6YpYiWiBsYLSJiTsyBybVuN8=
github.com/google/go-cmp v0.7.0/go.mod h1:pXiqmnSA92OHEEa9HXL2W4E7lf9JzCmGVUdgjX3N/iU=
github.com/google/uuid v1.6.0 h1:NIvaJDMOsjHA8n1jAhLSgzrAzy1Hgr+hNrb57e+94F0=
github.com/google/uuid v1.6.0/go.mod h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=
github.com/microsoft/kiota-abstractions-go v1.10.1 h1:5WecM4UpH60gpaKDlvsEVaZvyBwmuvjHevvJpRHqswE=
github.com/microsoft/kiota-abstractions-go v1.10.1/go.mod h1:kSLXbAikqjRQdpOYrCq3T91/4CkzmHPr41xu7oj6yzU=
github.com/microsoft/kiota-serialization-form-go v1.1.3 h1:eUY8eHXPFe4ma8cAdx0ya3g4NPlZgbPT+GlFC3xcgGY=
github.com/microsoft/kiota-serialization-form-go v1.1.3/go.mod h1:RMO99zyik+NvZjdVcIeyu6ikyfuKhQtzq2RK0fWJJio=
github.com/microsoft/kiota-serialization-json-go v1.1.4 h1:XAhsBjxiunPgsPLLHEbuvAe81Q3wD1pNOE22Dq2CD3k=
github.com/microsoft/kiota-serialization-json-go v1.1.4/go.mod h1:HUTiYs9llTGLjh9+O+yOkBbNEaZ1kxh3sBPU5tPhmeI=
github.com/microsoft/kiota-serialization-multipart-go v1.1.2 h1:1pUyA1QgIeKslQwbk7/ox1TehjlCUUT3r1f8cNlkvn4=
github.com/microsoft/kiota-serialization-multipart-go v1.1.2/go.mod h1:j2K7ZyYErloDu7Kuuk993DsvfoP7LPWvAo7rfDpdPio=
github.com/microsoft/kiota-serialization-text-go v1.1.3 h1:8z7Cebn0YAAr++xswVgfdxZjnAZ4GOB9O7XP4+r5r/M=
github.com/microsoft/kiota-serialization-text-go v1.1.3/go.mod h1:NDSvz4A3QalGMjNboKKQI9wR+8k+ih8UuagNmzIRgTQ=
github.com/std-uritemplate/std-uritemplate/go/v2 v2.0.12 h1:wxCxe+V5vu01YFrFgKj22gViqn8Yj4i9d9ecvdHvTGU=
github.com/std-uritemplate/std-uritemplate/go/v2 v2.0.12/go.mod h1:Z5KcoM0YLC7INlNhEezeIZ0TZNYf7WSNO0Lvah4DSeQ=
github.com/stretchr/testify v1.12.1 h1:EuwCh5fleGS7H32xRwO3wRGT7DxrDhLAT6FF8MpWDWE=
github.com/stretchr/testify v1.12.1/go.mod h1:MDEgiDPPsNp5cuIrHPPCyornHKgEVbtFUmoNlxoYthg=
go.opentelemetry.io/auto/sdk v1.2.1 h1:jXsnJ4Lmnqd11kwkBV2LgLoFMZKizbCi5fNZ/ipaZ64=
go.opentelemetry.io/auto/sdk v1.2.1/go.mod h1:KRTj+aOaElaLi+wW1kO/DZRXwkF4C5xPbEe3ZiIhN7Y=
go.opentelemetry.io/otel v1.45.0 h1:pdrWmLHofpubmArBv1LgFSv1Z0Ie/ppdZzu+kUN5EeU=
go.opentelemetry.io/otel v1.45.0/go.mod h1:XZxIqPapzEYnhNSScF5DIqXhm/rYi0FzCe2XddAwZfQ=
go.opentelemetry.io/otel/metric v1.45.0 h1:7Eg1uH7CJ5cXv9is6tnBe1FI6rj1nwUdbFypRm3br/M=
go.opentelemetry.io/otel/metric v1.45.0/go.mod h1:HAPbm1nd3p1PmFH7v2dR+6BjXxw+Lq4a2+pndMAm08s=
go.opentelemetry.io/otel/trace v1.45.0 h1:l/mP6Uv7oNO7/TblbhpbgMidxhq1uO/rPsikOyVhxag=
go.opentelemetry.io/otel/trace v1.45.0/go.mod h1:qoJJA2xNMnxRrdISU/kLtfUH2wNeQbiv+jhs/CxI8bc=
go.yaml.in/yaml/v3 v3.0.5 h1:N6y/pJk8buWs9NY5ERU2HSMfm+IuD/OtfdAnq6kESPw=
go.yaml.in/yaml/v3 v3.0.5/go.mod h1:HVTZu1O7/Vkt2N+BFy8Zza+lnLsABggaTM2ZpNIGuKg=
````

### FILE: `microsoft_kiota_openapi_gate/verify_contract.ps1`
```yaml
block_id: "MICROSOFT-KIOTA-OPENAPI-CLIENT-GATE:verify:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite reproducibility, negative-case, compilation, vet and SCA gate"
license: "LicenseRef-Workspace-Owner"
sha256: "70fdf44e0411fdba1c1bf194efd186f61fd2e4b519552a10f952eab20b3f4f23"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0

[CmdletBinding()]
param(
  [Parameter(Mandatory = $true)] [string] $GoExecutable,
  [Parameter(Mandatory = $true)] [string] $OsvScanner,
  [string] $WorkRoot = (Join-Path ([IO.Path]::GetTempPath()) ('elite-kiota-gate-' + [Guid]::NewGuid().ToString('N')))
)

$ErrorActionPreference = 'Stop'
$lock = Get-Content -LiteralPath (Join-Path $PSScriptRoot 'source-lock.json') -Raw | ConvertFrom-Json -Depth 30
$go = (Resolve-Path -LiteralPath $GoExecutable).Path
$osv = (Resolve-Path -LiteralPath $OsvScanner).Path
$work = [IO.Path]::GetFullPath($WorkRoot)
if (Test-Path -LiteralPath $work) { throw "WorkRoot must not exist: $work" }
[void](New-Item -ItemType Directory -Path $work)
$tool = Join-Path $work 'tool'
& (Join-Path $PSScriptRoot 'acquire_kiota.ps1') -Destination $tool
if ($LASTEXITCODE -ne 0) { throw 'Kiota acquisition failed' }
$exe = Join-Path $tool $lock.windows_x64.executable

$fixture = Join-Path $work 'ToDoApi.yaml'
Invoke-WebRequest -Uri $lock.official_fixture.url -OutFile $fixture
$fixtureHash = (Get-FileHash -LiteralPath $fixture -Algorithm SHA256).Hash.ToLowerInvariant()
if ((Get-Item -LiteralPath $fixture).Length -ne [int64]$lock.official_fixture.bytes -or $fixtureHash -cne $lock.official_fixture.sha256) {
  throw 'Official Kiota fixture identity mismatch'
}

$outA = Join-Path $work 'generated-a'
$outB = Join-Path $work 'generated-b'
$receiptA = Join-Path $work 'receipt-a.json'
$receiptB = Join-Path $work 'receipt-b.json'
$generate = Join-Path $PSScriptRoot 'generate_go_client.ps1'
& $generate -KiotaExecutable $exe -OpenApi $fixture -ExpectedOpenApiSha256 $fixtureHash -Output $outA -ClassName $lock.verified_lane.class_name -NamespaceName $lock.verified_lane.namespace_name -ReceiptPath $receiptA
if ($LASTEXITCODE -ne 0) { throw 'First Kiota generation failed' }
& $generate -KiotaExecutable $exe -OpenApi $fixture -ExpectedOpenApiSha256 $fixtureHash -Output $outB -ClassName $lock.verified_lane.class_name -NamespaceName $lock.verified_lane.namespace_name -ReceiptPath $receiptB
if ($LASTEXITCODE -ne 0) { throw 'Second Kiota generation failed' }
$a = Get-Content -LiteralPath $receiptA -Raw | ConvertFrom-Json -Depth 20
$b = Get-Content -LiteralPath $receiptB -Raw | ConvertFrom-Json -Depth 20
if ($a.output.code_tree_sha256 -cne $b.output.code_tree_sha256 -or $a.output.code_tree_sha256 -cne $lock.verified_lane.generated_code_tree_sha256) { throw 'Kiota code output is not reproducible or differs from lock' }
if ($a.output.files -ne $lock.verified_lane.generated_files_including_lock -or $a.output.code_files -ne $lock.verified_lane.generated_code_files -or $a.output.code_bytes -ne $lock.verified_lane.generated_code_bytes) { throw 'Kiota output inventory differs from lock' }
if ($a.output.kiota_description_hash -cne $b.output.kiota_description_hash -or $a.output.kiota_description_hash -cne $lock.verified_lane.kiota_description_hash) { throw 'Kiota semantic description hash differs from lock' }

$tampered = Join-Path $work 'tampered.yaml'
[IO.File]::WriteAllText($tampered, ([IO.File]::ReadAllText($fixture) + "`n"), [Text.UTF8Encoding]::new($false))
$tamperRejected = $false
try {
  & $generate -KiotaExecutable $exe -OpenApi $tampered -ExpectedOpenApiSha256 $fixtureHash -Output (Join-Path $work 'must-not-generate') -ClassName $lock.verified_lane.class_name -NamespaceName $lock.verified_lane.namespace_name -ReceiptPath (Join-Path $work 'must-not-exist.json')
} catch { $tamperRejected = $_.Exception.Message -match 'OpenAPI SHA-256 mismatch' }
if (-not $tamperRejected -or (Test-Path -LiteralPath (Join-Path $work 'must-not-generate'))) { throw 'Tampered OpenAPI was not rejected before generation' }
$occupiedRejected = $false
try {
  & $generate -KiotaExecutable $exe -OpenApi $fixture -ExpectedOpenApiSha256 $fixtureHash -Output $outA -ClassName $lock.verified_lane.class_name -NamespaceName $lock.verified_lane.namespace_name -ReceiptPath (Join-Path $work 'occupied.json')
} catch { $occupiedRejected = $_.Exception.Message -match 'Output must not exist' }
if (-not $occupiedRejected) { throw 'Occupied output was not rejected' }

Copy-Item -LiteralPath (Join-Path $PSScriptRoot 'fixture-go.mod') -Destination (Join-Path $outA 'go.mod')
Copy-Item -LiteralPath (Join-Path $PSScriptRoot 'fixture-go.sum') -Destination (Join-Path $outA 'go.sum')
if ((Get-FileHash -LiteralPath (Join-Path $outA 'go.mod') -Algorithm SHA256).Hash.ToLowerInvariant() -cne $lock.verified_lane.go_mod_sha256) { throw 'Go module lock drifted' }
if ((Get-FileHash -LiteralPath (Join-Path $outA 'go.sum') -Algorithm SHA256).Hash.ToLowerInvariant() -cne $lock.verified_lane.go_sum_sha256) { throw 'Go checksum lock drifted' }
$goVersion = (& $go version 2>&1 | Out-String).Trim()
if ($LASTEXITCODE -ne 0 -or $goVersion -notmatch ('go' + [regex]::Escape($lock.verified_lane.go_version) + '\b')) { throw "Go version mismatch: $goVersion" }
$env:GOTOOLCHAIN = 'local'
$env:PATH = (Split-Path -Parent $go) + [IO.Path]::PathSeparator + $env:PATH
Push-Location -LiteralPath $outA
try {
  & $go mod download
  if ($LASTEXITCODE -ne 0) { throw 'go mod download failed' }
  & $go test ./...
  if ($LASTEXITCODE -ne 0) { throw 'go test failed' }
  & $go vet ./...
  if ($LASTEXITCODE -ne 0) { throw 'go vet failed' }
} finally { Pop-Location }

$osvVersion = (& $osv --version 2>&1 | Out-String)
if ($LASTEXITCODE -ne 0 -or $osvVersion -notmatch ('osv-scanner version: ' + [regex]::Escape($lock.verified_lane.osv_scanner_version))) { throw 'OSV Scanner version mismatch' }
$report = Join-Path $work 'osv.json'
& $osv scan source --recursive --all-packages --format json --output-file $report $outA
$osvExit = $LASTEXITCODE
if (-not (Test-Path -LiteralPath $report -PathType Leaf)) { throw "OSV report missing: exit=$osvExit" }
$osvJson = Get-Content -LiteralPath $report -Raw | ConvertFrom-Json -Depth 100
$packages = @($osvJson.results | ForEach-Object { $_.packages } | Where-Object { $null -ne $_ })
$findings = @($packages | ForEach-Object { $_.vulnerabilities } | Where-Object { $null -ne $_ })
if ($osvExit -ne 0 -or $packages.Count -ne $lock.verified_lane.osv_packages -or $findings.Count -ne 0) {
  throw "OSV gate failed: exit=$osvExit packages=$($packages.Count) findings=$($findings.Count)"
}
Write-Output "MICROSOFT_KIOTA_OPENAPI_GATE_PASS version=$($lock.source.release) code_files=$($a.output.code_files) code_tree_sha256=$($a.output.code_tree_sha256) go=$($lock.verified_lane.go_version) osv_packages=$($packages.Count) findings=0 negative_cases=2"
````

### FILE: `microsoft_kiota_openapi_gate/README.md`
```yaml
block_id: "MICROSOFT-KIOTA-OPENAPI-CLIENT-GATE:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite operating boundary and project use instructions"
license: "LicenseRef-Workspace-Owner"
sha256: "96172b516a838adb5a15b40b15c4a8ac57b508ba324c91f9ecc3dd48093d0cac"
variables: []
secrets_allowed: false
```
````markdown
# Microsoft Kiota OpenAPI Gate

Este módulo adquiere el binario autocontenido oficial de Microsoft Kiota 1.35.0 para Windows x64, verifica bytes, SHA-256, inventario, ejecutable y versión, y genera un cliente Go sólo desde un contrato OpenAPI local cuyo SHA-256 fue aprobado.

## Uso en un proyecto

1. Fije el contrato OpenAPI oficial o propio autorizado y su SHA-256 en `PROJECT_EXTERNAL_SOURCE_LOCK.md`.
2. Ejecute `acquire_kiota.ps1` hacia una carpeta de herramientas nueva.
3. Ejecute `generate_go_client.ps1` con el path local, hash esperado, destino ausente, clase, namespace y receipt externo.
4. Revise el diff generado, congele `go.mod`/`go.sum`, compile, pruebe contratos contra sandbox oficial y ejecute SCA antes de incorporar el resultado.

El runner nunca limpia ni sobrescribe un output. Rechaza un contrato alterado antes de ejecutar Kiota y conserva un manifest con hash de cada archivo generado.

## Límite del claim

Kiota aporta generación tipada desde OpenAPI. No implementa autenticación comercial, autorización, secretos, deadlines, retry budget, idempotencia, webhooks, reconciliación, reglas del negocio ni evidencia productiva. Esas capas deben componerse mediante sus packs y gates específicos. El fixture `ToDoApi.yaml` sólo demuestra el generador y nunca se copia como dominio del proyecto.

## Auditoría reproducible

`verify_contract.ps1` descarga el asset y fixture oficiales fijados, exige dos árboles de código idénticos, valida por separado la semántica de `kiota-lock.json`, prueba rechazo por hash alterado y destino ocupado, compila el cliente con Go 1.26.7, ejecuta `go test`, `go vet` y OSV Scanner 2.5.1. `descriptionLocation` se excluye sólo del hash portable porque Kiota registra allí el path relativo del contrato; ningún archivo de código se excluye. Cero findings significa únicamente que OSV no conocía vulnerabilidades para el grafo detectado en la fecha de ejecución.
````

## 6. Configuration surface

Project inputs are the local OpenAPI path and approved SHA-256, an absent output path, an absent external receipt path, Go class/namespace, exact Kiota tool destination, exact Go executable and exact OSV Scanner. Provider URL, auth, secrets, retry, idempotency, webhook and reconciliation configuration remain outside this generator and must be supplied by their owning packs.

## 7. Dependency bill

- Microsoft Kiota 1.35.0 `win-x64.zip`, acquired at runtime by exact length/SHA-256 under MIT;
- Go 1.26.7 for the demonstrated lane;
- Google OSV-Scanner 2.5.1 for dated dependency SCA;
- PowerShell 7 for orchestration;
- network only for first acquisition/module download and current OSV lookup.

The fixture module pins five direct Microsoft Kiota Go modules and their transitive graph in `fixture-go.mod`/`fixture-go.sum`. Project output must freeze and review its own graph.

## 8. Apply order

1. Materialize this pack using `MICROSOFT_KIOTA_OPENAPI_CLIENT_PACK_PLAN.md`.
2. Fix the approved OpenAPI bytes/SHA-256 in the project source lock.
3. Acquire Kiota into an absent tool destination.
4. Generate into an absent output and preserve the external receipt.
5. Review the diff, freeze dependencies, compile/test/vet/SCA.
6. Compose provider auth, deadlines/retry, idempotency, webhooks and reconciliation.
7. Run sandbox contract tests, E2E, rollout and rollback gates.

## 9. Verification

Run `verify_contract.ps1 -GoExecutable <go-1.26.7> -OsvScanner <osv-2.5.1>`. It must report exact acquisition, two identical code trees, two negative cases, three compiling packages, vet and a current zero-finding SCA. A future advisory must fail and reopen dependency admission.

## 10. Reconstruction evidence

The governing execution evidence is `reconstruction_evidence/MICROSOFT_KIOTA_OPENAPI_CLIENT_GATE_2026-09-05_V247.md`. The pack remains conditioned on the selected project's authoritative OpenAPI, provider sandbox, auth, deadlines, idempotency, reconciliation, diff review and target tests.
