#requires -Version 7.0

[CmdletBinding()]
param(
  [string] $DotNetExecutable = '',
  [switch] $AllowNetwork
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
$lock = Get-Content -LiteralPath (Join-Path $PSScriptRoot 'source-lock.json') -Raw | ConvertFrom-Json -Depth 30
if ($lock.source.owner -cne 'Microsoft' -or $lock.source.commit -notmatch '^[0-9a-f]{40}$' -or $lock.source.tree -notmatch '^[0-9a-f]{40}$' -or $lock.source.archive_url -match '/(main|master|latest)(/|\.|$)' -or $lock.source.canonical_file_tree_sha256 -notmatch '^[0-9a-f]{64}$') { throw 'Source identity is not exact Microsoft source' }
if ($lock.source.license_expression -cne 'MIT' -or $lock.adaptation.classification -cne 'ADAPTED') { throw 'License or adaptation classification drifted' }
if ($lock.adaptation.sharpcompress_version -cne '0.48.0' -or $lock.adaptation.license_expression -cne 'MIT') { throw 'SharpCompress lock drifted' }
if ($lock.verified_lane.official_tests_passed -ne 300 -or $lock.verified_lane.nuget_vulnerability_findings -ne 0 -or $lock.verified_lane.typescript_fixture_rule -cne 'DS189424' -or $lock.verified_lane.go_fixture_rule -cne 'DS112852') { throw 'Verified lane drifted' }
foreach ($name in @('build_and_verify.ps1','verify_pack.ps1')) {
  $tokens=$null; $errors=$null
  [void][Management.Automation.Language.Parser]::ParseFile((Join-Path $PSScriptRoot $name),[ref]$tokens,[ref]$errors)
  if (@($errors).Count -ne 0) { throw "PowerShell syntax failed: $name" }
}
$runner = Join-Path $PSScriptRoot 'build_and_verify.ps1'
$negativeRoot = Join-Path ([IO.Path]::GetTempPath()) ('elite-devskim-negative-' + [Guid]::NewGuid().ToString('N'))
try {
  [void](New-Item -ItemType Directory -Path $negativeRoot)
  $existingDestination = Join-Path $negativeRoot 'existing-destination'
  [void](New-Item -ItemType Directory -Path $existingDestination)
  $failed = $false
  try { & $runner -Destination $existingDestination -ReceiptPath (Join-Path $negativeRoot 'receipt.json') -DotNetExecutable 'C:\missing\dotnet.exe' -AllowNetwork } catch { $failed = $_.Exception.Message -like 'Destination already exists:*' }
  if (-not $failed) { throw 'Existing-destination negative did not fail closed' }
  $existingReceipt = Join-Path $negativeRoot 'existing-receipt.json'
  [IO.File]::WriteAllText($existingReceipt, '{}')
  $failed = $false
  try { & $runner -Destination (Join-Path $negativeRoot 'new-destination') -ReceiptPath $existingReceipt -DotNetExecutable 'C:\missing\dotnet.exe' -AllowNetwork } catch { $failed = $_.Exception.Message -like 'Receipt already exists:*' }
  if (-not $failed) { throw 'Existing-receipt negative did not fail closed' }
  if ($AllowNetwork) {
    if (-not $DotNetExecutable) { throw 'DotNetExecutable is required with AllowNetwork' }
    & $runner -Destination (Join-Path $negativeRoot 'runtime') -ReceiptPath (Join-Path $negativeRoot 'runtime-receipt.json') -DotNetExecutable $DotNetExecutable -AllowNetwork
    if ($LASTEXITCODE -ne 0) { throw "Runtime verification exited $LASTEXITCODE" }
  }
  Write-Output "MICROSOFT_DEVSKIM_ADAPTED_PACK_PASS static=1 negatives=2 runtime=$([int][bool]$AllowNetwork)"
} finally {
  if (Test-Path -LiteralPath $negativeRoot) { [IO.Directory]::Delete($negativeRoot,$true) }
}
