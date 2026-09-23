#requires -Version 7.0

[CmdletBinding()]
param(
  [string]$LockPath = (Join-Path $PSScriptRoot 'fixture-lock.json'),
  [Parameter(Mandatory = $true)]
  [string[]]$FixtureId,
  [string]$DestinationRoot,
  [string]$ApprovalPath,
  [string]$ImportDirectory,
  [switch]$AllowNetwork,
  [switch]$ValidateOnly
)

$ErrorActionPreference = 'Stop'
$utf8 = [Text.UTF8Encoding]::new($false)

function Fail([string]$Message) { throw "AZURE_DI_OFFICIAL_FIXTURE_FAILED: $Message" }

function Assert-ExactProperties($Object, [string[]]$Expected, [string]$Label) {
  $actual = @($Object.PSObject.Properties.Name | Sort-Object)
  $wanted = @($Expected | Sort-Object)
  if (($actual -join '|') -cne ($wanted -join '|')) { Fail "$Label properties are not exact" }
}

function Assert-SafeRelativePath([string]$Value, [string]$Label) {
  $normalized = $Value.Replace('\', '/')
  $parts = $normalized.Split('/', [StringSplitOptions]::None)
  if ([string]::IsNullOrWhiteSpace($Value) -or [IO.Path]::IsPathRooted($Value) -or $parts -contains '' -or $parts -contains '.' -or $parts -contains '..') {
    Fail "$Label is unsafe"
  }
  foreach ($part in $parts) {
    if ($part.EndsWith('.') -or $part.EndsWith(' ')) { Fail "$Label is non-canonical" }
  }
  $normalized
}

if (-not (Test-Path -LiteralPath $LockPath -PathType Leaf)) { Fail 'lock not found' }
$lockBytes = [IO.File]::ReadAllBytes((Resolve-Path -LiteralPath $LockPath).Path)
$lockHash = [Convert]::ToHexString([Security.Cryptography.SHA256]::HashData($lockBytes)).ToLowerInvariant()
$lock = [Text.Encoding]::UTF8.GetString($lockBytes) | ConvertFrom-Json -Depth 20
Assert-ExactProperties $lock @('schema','verified_at','source','fixtures') 'lock'
if ($lock.schema -cne 'elite-azure-documentintelligence-official-fixture-lock/v1') { Fail 'unsupported lock schema' }
if ($lock.verified_at -cnotmatch '^[0-9]{4}-[0-9]{2}-[0-9]{2}$') { Fail 'invalid verified_at' }
Assert-ExactProperties $lock.source @('owner','repository','release','commit','commit_verified','license_expression','license_sha256') 'source'
if ($lock.source.owner -cne 'Microsoft' -or $lock.source.repository -cne 'Azure/azure-sdk-for-python' -or $lock.source.release -cne 'azure-ai-documentintelligence_1.0.2') { Fail 'unexpected source identity' }
if ($lock.source.commit -cnotmatch '^[0-9a-f]{40}$' -or $lock.source.commit_verified -cne $true) { Fail 'unverified source commit' }
if ($lock.source.license_expression -cne 'MIT' -or $lock.source.license_sha256 -cnotmatch '^[0-9a-f]{64}$') { Fail 'invalid license identity' }

$fixtureIds = @($lock.fixtures | ForEach-Object { $_.id })
if ($fixtureIds.Count -eq 0 -or ($fixtureIds | Sort-Object -Unique).Count -ne $fixtureIds.Count) { Fail 'fixture IDs must be non-empty and unique' }
$selectedIds = @($FixtureId)
if ($selectedIds.Count -eq 0 -or ($selectedIds | Sort-Object -Unique).Count -ne $selectedIds.Count) { Fail 'selected fixture IDs must be non-empty and unique' }
$selected = @()
foreach ($id in $selectedIds) {
  $fixture = @($lock.fixtures | Where-Object id -ceq $id)
  if ($fixture.Count -ne 1) { Fail "unknown fixture id: $id" }
  Assert-ExactProperties $fixture[0] @('id','url','filename','destination_path','bytes','sha256') "fixture $id"
  $uri = [Uri]$fixture[0].url
  if ($uri.Scheme -cne 'https' -or $uri.Host -cne 'raw.githubusercontent.com') { Fail "untrusted fixture host: $id" }
  if (-not $uri.AbsolutePath.Contains("/$($lock.source.commit)/", [StringComparison]::Ordinal)) { Fail "fixture URL is not commit pinned: $id" }
  if ($fixture[0].filename -cne [IO.Path]::GetFileName($uri.AbsolutePath)) { Fail "fixture filename mismatch: $id" }
  $null = Assert-SafeRelativePath $fixture[0].destination_path "fixture destination $id"
  if ([long]$fixture[0].bytes -le 0 -or $fixture[0].sha256 -cnotmatch '^[0-9a-f]{64}$') { Fail "invalid fixture integrity: $id" }
  $selected += $fixture[0]
}

if ($ValidateOnly) {
  Write-Output "AZURE_DI_OFFICIAL_FIXTURE_LOCK_VALID fixtures=$($lock.fixtures.Count) selected=$($selected.Count)"
  exit 0
}
if ([string]::IsNullOrWhiteSpace($DestinationRoot)) { Fail 'DestinationRoot is required' }
if (Test-Path -LiteralPath $DestinationRoot) { Fail 'destination already exists' }
if (-not $AllowNetwork -and [string]::IsNullOrWhiteSpace($ImportDirectory)) { Fail 'offline import or explicit network is required' }
if ($ImportDirectory -and -not (Test-Path -LiteralPath $ImportDirectory -PathType Container)) { Fail 'import directory not found' }
if (-not $ApprovalPath -or -not (Test-Path -LiteralPath $ApprovalPath -PathType Leaf)) { Fail 'approval not found' }
$approval = Get-Content -Raw -LiteralPath $ApprovalPath | ConvertFrom-Json -Depth 10
Assert-ExactProperties $approval @('schema','lock_sha256','approved_fixture_ids','approved_by','approved_at','accept_microsoft_mit_license') 'approval'
if ($approval.schema -cne 'elite-azure-documentintelligence-official-fixture-approval/v1') { Fail 'unsupported approval schema' }
if ($approval.lock_sha256 -cne $lockHash) { Fail 'approval lock hash mismatch' }
if ([string]::IsNullOrWhiteSpace($approval.approved_by) -or $approval.approved_at -cnotmatch '^[0-9]{4}-[0-9]{2}-[0-9]{2}$') { Fail 'approval identity/date missing' }
if ($approval.accept_microsoft_mit_license -cne $true) { Fail 'Microsoft MIT license not accepted' }
$approvedIds = @($approval.approved_fixture_ids)
$approvedSet = (($approvedIds | Sort-Object) -join '|')
$selectedSet = (($selectedIds | Sort-Object) -join '|')
if ($approvedSet -cne $selectedSet) { Fail 'approval fixture set mismatch' }

$destination = [IO.Path]::GetFullPath($DestinationRoot)
$parent = Split-Path -Parent $destination
if (-not (Test-Path -LiteralPath $parent -PathType Container)) { Fail 'destination parent not found' }
$staging = Join-Path $parent ('.azure-di-fixture-stage-' + [guid]::NewGuid().ToString('N'))
New-Item -ItemType Directory -Path $staging | Out-Null
try {
  $receiptItems = @()
  foreach ($fixture in $selected) {
    $download = Join-Path $staging $fixture.filename
    if ($ImportDirectory) {
      $cached = Join-Path ([IO.Path]::GetFullPath($ImportDirectory)) $fixture.filename
      if (-not (Test-Path -LiteralPath $cached -PathType Leaf)) { Fail "cached fixture missing: $($fixture.id)" }
      Copy-Item -LiteralPath $cached -Destination $download
    } else {
      Invoke-WebRequest -UseBasicParsing -Uri $fixture.url -OutFile $download
    }
    $item = Get-Item -LiteralPath $download
    $hash = (Get-FileHash -LiteralPath $download -Algorithm SHA256).Hash.ToLowerInvariant()
    if ($item.Length -ne [long]$fixture.bytes -or $hash -cne $fixture.sha256) { Fail "fixture bytes/hash mismatch: $($fixture.id)" }
    $relative = Assert-SafeRelativePath $fixture.destination_path "fixture destination $($fixture.id)"
    $target = [IO.Path]::GetFullPath((Join-Path $staging $relative))
    if (-not $target.StartsWith($staging + [IO.Path]::DirectorySeparatorChar, [StringComparison]::OrdinalIgnoreCase)) { Fail 'fixture escaped staging' }
    $targetParent = Split-Path -Parent $target
    New-Item -ItemType Directory -Path $targetParent -Force | Out-Null
    Move-Item -LiteralPath $download -Destination $target
    $receiptItems += [ordered]@{ id=$fixture.id; destination_path=$relative; bytes=[long]$fixture.bytes; sha256=$fixture.sha256; source=$fixture.url }
  }
  $receipt = [ordered]@{ schema='elite-azure-documentintelligence-official-fixture-receipt/v1'; acquired_at=[DateTimeOffset]::UtcNow.ToString('o'); lock_sha256=$lockHash; approved_by=$approval.approved_by; fixtures=$receiptItems }
  [IO.File]::WriteAllText((Join-Path $staging 'FIXTURE_RECEIPT.json'), ($receipt | ConvertTo-Json -Depth 10) + [Environment]::NewLine, $utf8)
  Move-Item -LiteralPath $staging -Destination $destination
  Write-Output "AZURE_DI_OFFICIAL_FIXTURE_ACQUISITION_PASS fixtures=$($selected.Count)"
} catch {
  if (Test-Path -LiteralPath $staging) { Remove-Item -LiteralPath $staging -Recurse -Force }
  throw
}
