#requires -Version 7.0
[CmdletBinding()]
param(
    [Parameter(Mandatory)][string]$BackendRoot,
    [Parameter(Mandatory)][string]$GeneratedClientRoot,
    [Parameter(Mandatory)][string]$DotnetPath,
    [Parameter(Mandatory)][string]$GoPath,
    [Parameter(Mandatory)][string]$WorkDirectory
)
$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest

function Require-Absolute([string]$PathValue, [string]$Name) {
    if (-not [IO.Path]::IsPathFullyQualified($PathValue)) { throw "$Name must be absolute" }
}
foreach ($entry in @(@($BackendRoot,'BackendRoot'), @($GeneratedClientRoot,'GeneratedClientRoot'), @($DotnetPath,'DotnetPath'), @($GoPath,'GoPath'), @($WorkDirectory,'WorkDirectory'))) {
    Require-Absolute $entry[0] $entry[1]
}
foreach ($file in @($DotnetPath,$GoPath,(Join-Path $BackendRoot 'go.mod'),(Join-Path $GeneratedClientRoot 'Elite.Arca.Wsfe.Generated.csproj'),(Join-Path $GeneratedClientRoot 'packages.lock.json'),(Join-Path $GeneratedClientRoot 'Generated/WsaaReference.cs'),(Join-Path $GeneratedClientRoot 'Generated/WsfeV1Reference.cs'))) {
    if (-not (Test-Path -LiteralPath $file -PathType Leaf)) { throw "Required file missing: $file" }
}
if (-not (Test-Path -LiteralPath $WorkDirectory)) { New-Item -ItemType Directory -Path $WorkDirectory | Out-Null }
if ((Get-ChildItem -LiteralPath $WorkDirectory -Force | Measure-Object).Count -ne 0) { throw 'WorkDirectory must be empty' }

$dotnetVersion = (& $DotnetPath --version).Trim()
$goVersion = (& $GoPath version).Trim()
if ($dotnetVersion -ne '10.0.400') { throw "Unexpected .NET SDK: $dotnetVersion" }
if ($goVersion -notmatch '\bgo1\.26\.7\b') { throw "Unexpected Go toolchain: $goVersion" }

Copy-Item -Path (Join-Path $BackendRoot '*') -Destination $WorkDirectory -Recurse -Force
$generatedDestination = Join-Path $WorkDirectory 'arca/fiscal/generated'
New-Item -ItemType Directory -Path (Join-Path $generatedDestination 'Generated') -Force | Out-Null
Copy-Item -LiteralPath (Join-Path $GeneratedClientRoot 'Elite.Arca.Wsfe.Generated.csproj') -Destination $generatedDestination
Copy-Item -LiteralPath (Join-Path $GeneratedClientRoot 'packages.lock.json') -Destination $generatedDestination
Copy-Item -LiteralPath (Join-Path $GeneratedClientRoot 'Generated/WsaaReference.cs') -Destination (Join-Path $generatedDestination 'Generated')
Copy-Item -LiteralPath (Join-Path $GeneratedClientRoot 'Generated/WsfeV1Reference.cs') -Destination (Join-Path $generatedDestination 'Generated')

$previousToolchain = $env:GOTOOLCHAIN
$env:GOTOOLCHAIN = 'local'
Push-Location $WorkDirectory
try {
    & $GoPath test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'Go tests failed' }
    & $GoPath vet ./...
    if ($LASTEXITCODE -ne 0) { throw 'Go vet failed' }
    & $GoPath build ./cmd/electromobility-api ./cmd/arca-fiscal-worker
    if ($LASTEXITCODE -ne 0) { throw 'Go build failed' }
    $testProject = Join-Path $WorkDirectory 'arca/fiscal/worker/tests/Elite.Arca.Wsfe.Worker.Tests/Elite.Arca.Wsfe.Worker.Tests.csproj'
    & $DotnetPath restore $testProject --locked-mode
    if ($LASTEXITCODE -ne 0) { throw '.NET locked restore failed' }
    & $DotnetPath build $testProject -c Release --no-restore
    if ($LASTEXITCODE -ne 0) { throw '.NET build failed' }
    & $DotnetPath run --project $testProject -c Release --no-build
    if ($LASTEXITCODE -ne 0) { throw '.NET worker contract tests failed' }
}
finally {
    Pop-Location
    $env:GOTOOLCHAIN = $previousToolchain
}
Write-Output 'ARCA_WSFE_UDS_WORKER_PASS go=1.26.7 dotnet=10.0.400 uds_tests=12 parameter_bridge=1 production_admitted=false'
