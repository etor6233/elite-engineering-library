#requires -Version 7.0
[CmdletBinding()]
param([string] $LibraryRoot = (Join-Path $PSScriptRoot '..'))
Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
$installer = Join-Path $LibraryRoot 'INSTALL_AGENT_BRIDGE.ps1'
$testRoot = Join-Path ([IO.Path]::GetTempPath()) ('elite-bridge-safety-' + [guid]::NewGuid().ToString('N'))
[void](New-Item -ItemType Directory -Path $testRoot)
$failed = [Collections.Generic.List[string]]::new()
$links = [Collections.Generic.List[string]]::new()
$checks = 0
function Check([bool] $Value, [string] $Name) {
  $script:checks++
  if (-not $Value) { $failed.Add($Name); Write-Output "FAIL: $Name" }
}
function Run-Bridge([string] $Project) {
  & $installer -ProjectRoot $Project -LibraryRoot $LibraryRoot -Agent Both | Out-Null
}
try {
  # Real directory redirection; never link to user data.
  foreach ($relative in @('.agents', '.agents/skills', '.claude/skills/elite-engineering-library', 'ancestor')) {
    $caseRoot = Join-Path $testRoot ([guid]::NewGuid().ToString('N'))
    $project = Join-Path $caseRoot 'project'
    $outside = Join-Path $caseRoot 'outside'
    [void](New-Item -ItemType Directory -Path $project, $outside)
    if ($relative -eq 'ancestor') {
      [void](New-Item -ItemType Directory -Path (Join-Path $outside 'child'))
      $link = Join-Path $caseRoot 'redirect'
      $project = Join-Path $link 'child'
    } else {
      $link = Join-Path $project $relative
      [void](New-Item -ItemType Directory -Path (Split-Path -Parent $link) -Force)
    }
    $linkType = if ($IsWindows) { 'Junction' } else { 'SymbolicLink' }
    [void](New-Item -ItemType $linkType -Path $link -Target $outside)
    $links.Add($link)
    $before = @(Get-ChildItem -LiteralPath $outside -File -Recurse -Force).Count
    $rejected = $false
    try { Run-Bridge $project } catch { $rejected = $_.Exception.Message -match 'reparse|symlink' }
    Check $rejected "reject linked path: $relative"
    Check (@(Get-ChildItem -LiteralPath $outside -File -Recurse -Force).Count -eq $before) "no outside writes: $relative"
    Check (-not (Test-Path -LiteralPath (Join-Path $project 'AGENTS.md'))) "preflight before any write: $relative"
    # Remove only the link, never recurse through it.
    Remove-Item -LiteralPath $link -Force
    [void]$links.Remove($link)
  }
  $reversed = Join-Path $testRoot 'reversed-markers'
  [void](New-Item -ItemType Directory -Path $reversed)
  $path = Join-Path $reversed 'AGENTS.md'
  [IO.File]::WriteAllText($path, "<!-- ELITE-ENGINEERING-LIBRARY:END -->`nkeep`n<!-- ELITE-ENGINEERING-LIBRARY:BEGIN -->")
  $before = (Get-FileHash -LiteralPath $path).Hash
  $rejected = $false
  try { Run-Bridge $reversed } catch { $rejected = $_.Exception.Message -match 'markers' }
  Check $rejected 'reject reversed managed markers'
  Check ((Get-FileHash -LiteralPath $path).Hash -eq $before) 'reversed marker bytes unchanged'
  Check (-not (Test-Path -LiteralPath (Join-Path $reversed '.agents'))) 'reversed markers produce no skill'

  # FileShare.Read permits planning/read but denies the replacement at write time.
  if ($IsWindows) {
    foreach ($existing in @($true, $false)) {
      $project = Join-Path $testRoot "rollback-$existing"
      $skillPath = Join-Path $project '.agents/skills/elite-engineering-library/SKILL.md'
      [void](New-Item -ItemType Directory -Path (Split-Path -Parent $skillPath))
      [IO.File]::WriteAllText($skillPath, '<!-- ELITE-ENGINEERING-LIBRARY:MANAGED-SKILL -->')
      $agentsPath = Join-Path $project 'AGENTS.md'
      if ($existing) { [IO.File]::WriteAllText($agentsPath, "# Original`r`nKeep bytes.`r`n", [Text.UTF8Encoding]::new($true)) }
      $before = if ($existing) { (Get-FileHash -LiteralPath $agentsPath).Hash } else { '' }
      $handle = [IO.File]::Open($skillPath, [IO.FileMode]::Open, [IO.FileAccess]::Read, [IO.FileShare]::Read)
      $errorText = ''
      try { Run-Bridge $project } catch { $errorText = $_.Exception.Message } finally { $handle.Dispose() }
      Check ($errorText.Length -gt 0) "write failure surfaced: existing=$existing"
      Check ($errorText -notmatch 'Reverse|parameter cannot be found') "original failure preserved: existing=$existing"
      $restored = if ($existing) { (Get-FileHash -LiteralPath $agentsPath).Hash -eq $before } else { -not (Test-Path -LiteralPath $agentsPath) }
      Check $restored "rollback restores original bytes or absence: existing=$existing"
      Check (@(Get-ChildItem -LiteralPath $project -Recurse -Force -Filter '.elite-bridge-*.tmp').Count -eq 0) "no temporary payload remains: existing=$existing"
    }
  } else { throw 'Windows file-sharing rollback scenarios require Windows; not counted as PASS' }
  if ($failed.Count) { throw "BRIDGE_SAFETY_FAILED failures=$($failed.Count) checks=$checks" }
  Write-Output "BRIDGE_SAFETY_PASS checks=$checks"
} finally {
  foreach ($link in $links) { Remove-Item -LiteralPath $link -Force }
  $resolved = [IO.Path]::GetFullPath($testRoot)
  $prefix = [IO.Path]::GetFullPath([IO.Path]::GetTempPath()).TrimEnd([IO.Path]::DirectorySeparatorChar) + [IO.Path]::DirectorySeparatorChar
  if (-not $resolved.StartsWith($prefix, [StringComparison]::OrdinalIgnoreCase) -or (Split-Path -Leaf $resolved) -notlike 'elite-bridge-safety-*') { throw 'unsafe test cleanup path' }
  Remove-Item -LiteralPath $resolved -Recurse -Force
}
