#requires -Version 7.0
# CATALOG / STOCK / PI-PL slice verifier — local fixtures only.
# See docs/slices/CATALOG_STOCK_COMPOSE.md

[CmdletBinding()]
param(
  [string] $DatabaseUrl = $env:TEST_DATABASE_URL,
  [string] $PsqlExecutable = '',
  [string] $GoExecutable = ''
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
$libraryRoot = [IO.Path]::GetFullPath((Join-Path $PSScriptRoot '..'))

function Resolve-Tool([string] $Name, [string] $Explicit = '') {
  if ($Explicit) {
    $path = [IO.Path]::GetFullPath($Explicit)
    if (Test-Path -LiteralPath $path -PathType Leaf) { return $path }
    return $null
  }
  $command = Get-Command $Name -ErrorAction SilentlyContinue | Select-Object -First 1
  if ($command) { return $command.Source }
  $null
}

if (-not $DatabaseUrl) {
  throw 'TEST_DATABASE_URL is required (owned loopback disposable database with migrations applied)'
}

$psql = Resolve-Tool 'psql' $PsqlExecutable
if (-not $psql) {
  throw 'psql not found; install PostgreSQL client or pass -PsqlExecutable'
}

$go = Resolve-Tool 'go' $GoExecutable
if (-not $go) {
  throw 'go not found; pass -GoExecutable'
}

$sqlTest = Join-Path $libraryRoot 'db/tests/0088_catalog_stock_consistency.test.sql'
if (-not (Test-Path -LiteralPath $sqlTest -PathType Leaf)) {
  throw "missing slice SQL test: $sqlTest"
}

Write-Host '== CATALOG_STOCK_SLICE: SQL FK + join ==' -ForegroundColor Cyan
& $psql $DatabaseUrl -v ON_ERROR_STOP=1 -f $sqlTest
if (-not $?) { throw 'SQL catalog-stock consistency test failed' }

Write-Host '== CATALOG_STOCK_SLICE: Go integration ==' -ForegroundColor Cyan
$env:TEST_DATABASE_URL = $DatabaseUrl
Push-Location $libraryRoot
try {
  & $go test ./internal/platform/postgres -run '^TestCatalogStockConsistency' -count=1
  if (-not $?) { throw 'Go catalog-stock consistency tests failed' }
} finally {
  Pop-Location
}

Write-Host 'CATALOG_STOCK_SLICE PASS (local fixtures; not production)' -ForegroundColor Green
