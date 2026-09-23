#requires -Version 7.0

[CmdletBinding()]
param(
  [switch] $RunTarget,
  [string] $TargetUrl = '',
  [string] $ChromePath = '',
  [string] $OutputDirectory = ''
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
$root = [IO.Path]::GetFullPath($PSScriptRoot)

function Fail([string] $Message) { throw "GOOGLE_LIGHTHOUSE_WEB_QUALITY_GATE_FAILED: $Message" }

$nodeVersion = (& node --version | Out-String).Trim()
if ($LASTEXITCODE -ne 0 -or $nodeVersion -notmatch '^v(?<major>\d+)\.') { Fail 'Node.js is unavailable' }
if ([int]$Matches.major -lt 22) { Fail "Node.js >=22.19 is required; actual=$nodeVersion" }
$pnpmVersion = (& pnpm --version | Out-String).Trim()
if ($LASTEXITCODE -ne 0 -or $pnpmVersion -ne '11.25.0') { Fail "pnpm 11.25.0 is required; actual=$pnpmVersion" }
$lighthouseVersion = (& pnpm --dir $root --ignore-workspace exec lighthouse --version | Out-String).Trim()
if ($LASTEXITCODE -ne 0 -or $lighthouseVersion -ne '13.4.1') { Fail "Lighthouse 13.4.1 is required; actual=$lighthouseVersion" }

$sourceLock = Get-Content -Raw -LiteralPath (Join-Path $root 'source-lock.json') | ConvertFrom-Json -Depth 30
if ($sourceLock.source.commit -cne '1d58f5b06d28e3419b38817a6c7488ec4413c67d' -or $sourceLock.npm.tarballSha256 -cne '110759ba9e863c024e214e9b08ed2b0344d89b492286227235d4dfb990dc3e54') { Fail 'source identity drifted' }
$licensePath = Join-Path $root 'LICENSE.lighthouse.txt'
if ((Get-Item -LiteralPath $licensePath).Length -ne 11358 -or (Get-FileHash -Algorithm SHA256 -LiteralPath $licensePath).Hash.ToLowerInvariant() -cne 'cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30') { Fail 'Apache-2.0 license bytes drifted' }

& node --test (Join-Path $root 'tests/validate-reports.test.mjs')
if ($LASTEXITCODE -ne 0) { Fail 'negative policy suite failed' }

if ($RunTarget) {
  if ([string]::IsNullOrWhiteSpace($TargetUrl)) { Fail 'TargetUrl is required with RunTarget' }
  if ([string]::IsNullOrWhiteSpace($ChromePath)) { Fail 'ChromePath is required with RunTarget' }
  $resolvedChrome = [IO.Path]::GetFullPath($ChromePath)
  if (-not (Test-Path -LiteralPath $resolvedChrome -PathType Leaf)) { Fail "Chrome executable is missing: $resolvedChrome" }
  if ([string]::IsNullOrWhiteSpace($OutputDirectory)) { $OutputDirectory = Join-Path $root 'lighthouse-results' }
  $resolvedOutput = [IO.Path]::GetFullPath($OutputDirectory)
  if (Test-Path -LiteralPath $resolvedOutput) { Fail "output directory must not exist: $resolvedOutput" }
  $executionRoot = $root
  $capsule = $null
  try {
    if ($IsWindows -and $root.Length -gt 110) {
      $capsule = Join-Path ([IO.Path]::GetTempPath()) ('elite-lh-run-' + [guid]::NewGuid().ToString('N'))
      [void](New-Item -ItemType Directory -Path $capsule)
      [void](New-Item -ItemType Directory -Path (Join-Path $capsule 'lib'))
      foreach ($relative in @('package.json', 'pnpm-lock.yaml', 'quality-policy.json', 'run-quality-gate.mjs', 'lib/validate-reports.mjs')) {
        [IO.File]::Copy((Join-Path $root $relative), (Join-Path $capsule $relative), $false)
      }
      & pnpm --dir $capsule --ignore-workspace install --frozen-lockfile --offline
      if ($LASTEXITCODE -ne 0) { Fail 'failed to build the Windows short-path execution capsule' }
      $executionRoot = $capsule
    }
    & node (Join-Path $executionRoot 'run-quality-gate.mjs') "--url=$TargetUrl" "--chrome-path=$resolvedChrome" "--output-dir=$resolvedOutput" "--policy=$(Join-Path $executionRoot 'quality-policy.json')"
    if ($LASTEXITCODE -ne 0) { Fail 'target quality gate failed' }
  } finally {
    if ($capsule) {
      $temp = [IO.Path]::GetFullPath([IO.Path]::GetTempPath()).TrimEnd('\')
      $resolvedCapsule = [IO.Path]::GetFullPath($capsule)
      if (-not $resolvedCapsule.StartsWith($temp + '\', [StringComparison]::OrdinalIgnoreCase) -or (Split-Path -Leaf $resolvedCapsule) -notlike 'elite-lh-run-*') {
        Fail 'refused unsafe Lighthouse execution-capsule cleanup'
      }
      $cleanupError = $null
      for ($attempt = 1; $attempt -le 5; $attempt++) {
        try {
          Remove-Item -LiteralPath $resolvedCapsule -Recurse -Force -ErrorAction Stop
          $cleanupError = $null
          break
        } catch {
          $cleanupError = $_
          Start-Sleep -Milliseconds 200
        }
      }
      if ($cleanupError -or (Test-Path -LiteralPath $resolvedCapsule)) { Fail 'failed to remove the Lighthouse execution capsule' }
    }
  }
}

Write-Output ('GOOGLE_LIGHTHOUSE_WEB_QUALITY_GATE_CONTRACT_PASS runtime=13.4.1 tests=5 target=' + $(if ($RunTarget) {'5-runs'} else {'SKIPPED'}))
