#requires -Version 7.0

[CmdletBinding()]
param(
  [string]$ComponentRoot = $PSScriptRoot,
  [switch]$RunEnterpriseWeb,
  [string]$WebRoot
)

$ErrorActionPreference = 'Stop'
$root = (Resolve-Path -LiteralPath $ComponentRoot).Path
$lock = Get-Content -Raw -LiteralPath (Join-Path $root 'source-lock.json') | ConvertFrom-Json -Depth 30
$package = Get-Content -Raw -LiteralPath (Join-Path $root 'package.json') | ConvertFrom-Json -Depth 20
if ($lock.source.repository -cne 'microsoft/playwright' -or $lock.source.release -cne 'v1.62.1' -or $lock.source.commit -cne '26a9e470a7b3c7822084b09fb7f13902c5f37b51' -or $lock.source.commitSignatureVerified -ne $true) { throw 'Microsoft source identity drifted' }
if ($lock.npm.package -cne '@playwright/test' -or $lock.npm.version -cne '1.62.1' -or $lock.npm.graphPackages -ne 4 -or $lock.npm.datedOsvFindings -ne 0) { throw 'Playwright artifact/runtime evidence drifted' }
$expectedBrowsers = @(
  @{ name='Chromium'; revision='1234'; version='151.0.7922.34'; projects=2 },
  @{ name='Firefox'; revision='1538'; version='153.0'; projects=1 },
  @{ name='WebKit'; revision='2336'; version='26.5'; projects=1 }
)
if (@($lock.browsers).Count -ne 3) { throw 'Playwright browser inventory drifted' }
for ($index = 0; $index -lt $expectedBrowsers.Count; $index++) {
  $actual = $lock.browsers[$index]; $expected = $expectedBrowsers[$index]
  if ($actual.name -cne $expected.name -or $actual.revision -cne $expected.revision -or $actual.version -cne $expected.version -or $actual.projects -ne $expected.projects) { throw "Playwright browser identity drifted: $($expected.name)" }
}
if ($package.devDependencies.'@playwright/test' -cne '1.62.1' -or $package.packageManager -cne 'pnpm@11.25.0') { throw 'Package pins drifted' }
$license = Join-Path $root 'LICENSE.playwright.txt'
if ((Get-Item -LiteralPath $license).Length -ne 11399 -or (Get-FileHash -LiteralPath $license -Algorithm SHA256).Hash.ToLowerInvariant() -cne '7fab1461b41970ff376f1c9303a637076bfaaeb71cd12dd3a1c44aaf59a1a2b9' -or $lock.source.packagedLicenseNormalization -cne 'CRLF_TO_LF_ONLY') { throw 'Packaged Microsoft license normalization drifted' }
$lockHash = (Get-FileHash -LiteralPath (Join-Path $root 'pnpm-lock.yaml') -Algorithm SHA256).Hash.ToLowerInvariant()
if ($lockHash -cne '63eec1e3d5bac29965e850bf751269c7a8b317929b570205a1e223544f80a470') { throw 'pnpm lock drifted' }
if (-not (Test-Path -LiteralPath (Join-Path $root 'node_modules') -PathType Container)) { throw 'Install with pnpm --ignore-workspace --frozen-lockfile before runtime verification' }
Push-Location $root
try {
  $version = (& pnpm --ignore-workspace exec playwright --version | Out-String).Trim()
  if ($LASTEXITCODE -ne 0 -or $version -cne 'Version 1.62.1') { throw "Unexpected Playwright runtime: $version" }
  $env:ELITE_RUNTIME_ONLY = '1'
  & pnpm --ignore-workspace exec playwright test tests/runtime.spec.mjs
  Remove-Item Env:ELITE_RUNTIME_ONLY -ErrorAction SilentlyContinue
  if ($LASTEXITCODE -ne 0) { throw 'Microsoft Playwright runtime smoke failed' }
  if ($RunEnterpriseWeb) {
    if ([string]::IsNullOrWhiteSpace($WebRoot)) { throw 'WebRoot is required with RunEnterpriseWeb' }
    $resolvedWeb = (Resolve-Path -LiteralPath $WebRoot).Path
    foreach ($required in @('package.json','.next/BUILD_ID','config/business.example.json')) {
      if (-not (Test-Path -LiteralPath (Join-Path $resolvedWeb $required) -PathType Leaf)) { throw "Enterprise web target missing: $required" }
    }
    $env:ELITE_WEB_ROOT = $resolvedWeb
    Remove-Item Env:ELITE_BASE_URL -ErrorAction SilentlyContinue
    & pnpm --ignore-workspace exec playwright test tests/enterprise-web.spec.mjs tests/nonce-csp.spec.mjs
    if ($LASTEXITCODE -ne 0) { throw 'Enterprise web browser gate failed' }
  }
} finally {
  Pop-Location
}
'MICROSOFT_PLAYWRIGHT_BROWSER_GATE_PASS runtime=4 target=' + ($(if ($RunEnterpriseWeb) {'8'} else {'SKIPPED'}))
