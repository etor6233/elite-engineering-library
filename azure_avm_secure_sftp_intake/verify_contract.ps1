#requires -Version 7.0
[CmdletBinding()]
param(
  [string]$SourceRoot = $PSScriptRoot,
  [string]$BicepExe = '',
  [switch]$AllowNetwork
)

$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest

function Assert-Present([string]$Text, [string]$Literal) {
  $count = ([regex]::Matches($Text, [regex]::Escape($Literal))).Count
  if ($count -lt 1) { throw "Required contract token missing: [$Literal]" }
}

$mainPath = Join-Path $SourceRoot 'main.bicep'
$lockPath = Join-Path $SourceRoot 'UPSTREAM.lock.json'
if (-not (Test-Path -LiteralPath $mainPath -PathType Leaf)) { throw 'main.bicep missing' }
if (-not (Test-Path -LiteralPath $lockPath -PathType Leaf)) { throw 'UPSTREAM.lock.json missing' }

$main = Get-Content -Raw -LiteralPath $mainPath
$lock = Get-Content -Raw -LiteralPath $lockPath | ConvertFrom-Json -Depth 20
if ($lock.azure_verified_modules.storage_account_module -cne 'br/public:avm/res/storage/storage-account:0.33.0') { throw 'AVM storage module lock drifted' }
if ($lock.bicep_cli.version -cne '0.46.1') { throw 'Bicep version lock drifted' }
if ($lock.bicep_cli.sha256 -cne '441d3d6094513acaa8a9be0b96b754d174bc0c19064c7be04e35b87e67b66250') { throw 'Bicep hash lock drifted' }

$required = @(
  "br/public:avm/res/storage/storage-account:0.33.0",
  'enableHierarchicalNamespace: true',
  'enableSftp: true',
  'isLocalUserEnabled: true',
  'allowSharedKeyAccess: false',
  "publicNetworkAccess: 'Disabled'",
  'allowBlobPublicAccess: false',
  'requireInfrastructureEncryption: true',
  "minimumTlsVersion: 'TLS1_2'",
  "bypass: 'None'",
  "defaultAction: 'Deny'",
  "service: 'blob'",
  'privateDnsZoneResourceId: privateDnsZoneResourceId',
  'containerDeleteRetentionPolicyEnabled: true',
  'deleteRetentionPolicyEnabled: true',
  'isVersioningEnabled: false',
  "publicAccess: 'None'",
  "{ category: 'StorageRead' }",
  "{ category: 'StorageWrite' }",
  "{ category: 'StorageDelete' }",
  'hasSharedKey: false',
  'hasSshKey: true',
  'hasSshPassword: false',
  "permissions: 'cw'",
  'key: sshPublicKey'
)
foreach ($literal in $required) { Assert-Present $main $literal }
foreach ($literal in @("br/public:avm/res/storage/storage-account:0.33.0", 'hasSshPassword: false', "permissions: 'cw'")) {
  $count = ([regex]::Matches($main, [regex]::Escape($literal))).Count
  if ($count -ne 1) { throw "Expected exactly one occurrence of [$literal], found $count" }
}

$forbidden = @(
  'hasSshPassword: true',
  'hasSharedKey: true',
  'allowSharedKeyAccess: true',
  "publicNetworkAccess: 'Enabled'",
  "permissions: 'r'",
  "permissions: 'rcwdl'",
  'isVersioningEnabled: true'
)
foreach ($literal in $forbidden) {
  if ($main.Contains($literal, [System.StringComparison]::Ordinal)) { throw "Forbidden contract token present: $literal" }
}

$compile = 'SKIPPED_NO_PINNED_BICEP'
if ([string]::IsNullOrWhiteSpace($BicepExe)) {
  $candidate = Get-Command bicep -ErrorAction SilentlyContinue
  if ($null -ne $candidate) { $BicepExe = $candidate.Source }
}
if ([string]::IsNullOrWhiteSpace($BicepExe) -and $AllowNetwork) {
  $toolRoot = Join-Path ([IO.Path]::GetTempPath()) 'elite-bicep-cli-0.46.1-441d3d60'
  if (-not (Test-Path -LiteralPath $toolRoot)) { New-Item -ItemType Directory -Path $toolRoot | Out-Null }
  $BicepExe = Join-Path $toolRoot 'bicep-win-x64.exe'
  if (-not (Test-Path -LiteralPath $BicepExe -PathType Leaf)) {
    Invoke-WebRequest -Uri $lock.bicep_cli.url -OutFile $BicepExe
  }
}

if (-not [string]::IsNullOrWhiteSpace($BicepExe)) {
  if (-not (Test-Path -LiteralPath $BicepExe -PathType Leaf)) { throw 'Specified Bicep executable missing' }
  $actualSize = (Get-Item -LiteralPath $BicepExe).Length
  $actualHash = (Get-FileHash -LiteralPath $BicepExe -Algorithm SHA256).Hash.ToLowerInvariant()
  if ($actualSize -ne [long]$lock.bicep_cli.size) { throw "Bicep size mismatch: $actualSize" }
  if ($actualHash -cne $lock.bicep_cli.sha256) { throw "Bicep SHA-256 mismatch: $actualHash" }
  $versionText = (& $BicepExe --version 2>&1) -join "`n"
  if ($LASTEXITCODE -ne 0 -or $versionText -notmatch '0\.46\.1') { throw "Unexpected Bicep version: $versionText" }
  $outputRoot = Join-Path ([IO.Path]::GetTempPath()) ('elite-sftp-arm-' + [guid]::NewGuid().ToString('N'))
  New-Item -ItemType Directory -Path $outputRoot | Out-Null
  $armPath = Join-Path $outputRoot 'main.json'
  & $BicepExe build $mainPath --outfile $armPath
  if ($LASTEXITCODE -ne 0 -or -not (Test-Path -LiteralPath $armPath -PathType Leaf)) { throw 'Bicep compilation failed' }
  $armText = Get-Content -Raw -LiteralPath $armPath
  foreach ($token in @('isSftpEnabled', 'hasSshPassword', 'allowSharedKeyAccess', 'publicNetworkAccess', 'StorageWrite', 'privateDnsZoneResourceId')) {
    if (-not $armText.Contains($token, [System.StringComparison]::Ordinal)) { throw "Compiled ARM missing token: $token" }
  }
  $compile = 'PASS_PINNED_0.46.1'
}

Write-Output "MICROSOFT_AVM_SECURE_SFTP_INTAKE_PASS static=1 compile=$compile password=0 shared_key=0 public_network=0 upload_permissions=cw retained_store=SEPARATE_REQUIRED"
