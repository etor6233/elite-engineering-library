#requires -Version 7.0

[CmdletBinding(SupportsShouldProcess = $true)]
param(
  [Parameter(Mandatory = $true)]
  [string] $ProjectRoot,
  [string] $LibraryRoot = $PSScriptRoot,
  [ValidateSet('Codex','Claude','Both')]
  [string] $Agent = 'Both'
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
$utf8 = [Text.UTF8Encoding]::new($false)
$strictUtf8 = [Text.UTF8Encoding]::new($false, $true)
$begin = '<!-- ELITE-ENGINEERING-LIBRARY:BEGIN -->'
$end = '<!-- ELITE-ENGINEERING-LIBRARY:END -->'
$skillMarker = '<!-- ELITE-ENGINEERING-LIBRARY:MANAGED-SKILL -->'
$maximumAgentsBytes = 32768

function Fail([string] $Message) { throw "INSTALL_AGENT_BRIDGE_FAILED: $Message" }

function Assert-NoReparsePath([string] $Path) {
  $full = [IO.Path]::GetFullPath($Path)
  $cursor = [IO.Path]::GetPathRoot($full)
  foreach ($part in $full.Substring($cursor.Length).Split([char[]]@([IO.Path]::DirectorySeparatorChar, [IO.Path]::AltDirectorySeparatorChar), [StringSplitOptions]::RemoveEmptyEntries)) {
    $cursor = Join-Path $cursor $part
    $item = Get-Item -LiteralPath $cursor -Force -ErrorAction SilentlyContinue
    if ($null -ne $item -and ($item.Attributes -band [IO.FileAttributes]::ReparsePoint) -ne 0) {
      Fail "path must not traverse a symlink/reparse point: $cursor"
    }
  }
}

function Resolve-ExistingDirectory([string] $Path, [string] $Name) {
  Assert-NoReparsePath $Path
  if (-not (Test-Path -LiteralPath $Path -PathType Container)) { Fail "$Name does not exist: $Path" }
  $resolved = [IO.Path]::GetFullPath((Resolve-Path -LiteralPath $Path).Path).TrimEnd([IO.Path]::DirectorySeparatorChar, [IO.Path]::AltDirectorySeparatorChar)
  $item = Get-Item -LiteralPath $resolved -Force
  if (($item.Attributes -band [IO.FileAttributes]::ReparsePoint) -ne 0) { Fail "$Name must not be a symlink/reparse point: $resolved" }
  $resolved
}

function Read-Utf8([string] $Path) {
  Assert-NoReparsePath $Path
  $bytes = [IO.File]::ReadAllBytes($Path)
  try { $strictUtf8.GetString($bytes) } catch { Fail "file is not valid UTF-8: $Path" }
}

function Merge-ManagedBlock([string] $Existing, [string] $Block, [string] $Path, [string] $NewHeading) {
  $beginCount = ([regex]::Matches($Existing, [regex]::Escape($begin))).Count
  $endCount = ([regex]::Matches($Existing, [regex]::Escape($end))).Count
  if ($beginCount -ne $endCount -or $beginCount -gt 1) { Fail "managed markers are damaged or duplicated: $Path" }
  if ($beginCount -eq 1) {
    if ($Existing.IndexOf($begin, [StringComparison]::Ordinal) -ge $Existing.IndexOf($end, [StringComparison]::Ordinal)) { Fail "managed markers are reversed: $Path" }
    $pattern = '(?s)' + [regex]::Escape($begin) + '.*?' + [regex]::Escape($end)
    return [regex]::Replace($Existing, $pattern, [Text.RegularExpressions.MatchEvaluator]{ param($match) $Block }, 1)
  }
  if (-not $Existing) { return "$NewHeading`n`n$Block`n" }
  $separator = if ($Existing.EndsWith("`n`n", [StringComparison]::Ordinal)) { '' } elseif ($Existing.EndsWith("`n", [StringComparison]::Ordinal)) { "`n" } else { "`n`n" }
  $Existing + $separator + $Block + "`n"
}

function Add-PlannedWrite([Collections.Generic.List[object]] $Plan, [string] $Path, [string] $Content, [bool] $ManagedSkill = $false) {
  Assert-NoReparsePath $Path
  $existing = $null
  $originalBytes = $null
  if (Test-Path -LiteralPath $Path) {
    if (-not (Test-Path -LiteralPath $Path -PathType Leaf)) { Fail "target exists but is not a file: $Path" }
    $item = Get-Item -LiteralPath $Path -Force
    if (($item.Attributes -band [IO.FileAttributes]::ReparsePoint) -ne 0) { Fail "target must not be a symlink/reparse point: $Path" }
    $existing = Read-Utf8 $Path
    $originalBytes = [IO.File]::ReadAllBytes($Path)
    if ($ManagedSkill -and -not $existing.Contains($skillMarker, [StringComparison]::Ordinal)) {
      Fail "refusing to replace an unmanaged skill: $Path"
    }
  }
  $Plan.Add([pscustomobject]@{ Path=$Path; Content=$Content; Existed=($null -ne $existing); OriginalBytes=$originalBytes })
}

function Invoke-PlannedWrites([Collections.Generic.List[object]] $Plan) {
  $completed = [Collections.Generic.List[object]]::new()
  try {
    foreach ($entry in $Plan) {
      $parent = Split-Path -Parent $entry.Path
      if ($PSCmdlet.ShouldProcess($entry.Path, 'install or update Elite agent bridge')) {
        Assert-NoReparsePath $entry.Path
        if (-not (Test-Path -LiteralPath $parent)) { [void](New-Item -ItemType Directory -Path $parent -Force) }
        Assert-NoReparsePath $entry.Path
        $temporary = Join-Path $parent ('.elite-bridge-' + [guid]::NewGuid().ToString('N') + '.tmp')
        try {
          [IO.File]::WriteAllText($temporary, $entry.Content, $utf8)
          [IO.File]::Move($temporary, $entry.Path, $true)
          $completed.Add($entry)
        } finally {
          if (Test-Path -LiteralPath $temporary) { Remove-Item -LiteralPath $temporary -Force }
        }
      }
    }
  } catch {
    for ($index = $completed.Count - 1; $index -ge 0; $index--) {
      $entry = $completed[$index]
      Assert-NoReparsePath $entry.Path
      if ($entry.Existed) { [IO.File]::WriteAllBytes($entry.Path, $entry.OriginalBytes) }
      elseif (Test-Path -LiteralPath $entry.Path -PathType Leaf) { Remove-Item -LiteralPath $entry.Path -Force }
    }
    throw
  }
}

$project = Resolve-ExistingDirectory $ProjectRoot 'ProjectRoot'
$library = Resolve-ExistingDirectory $LibraryRoot 'LibraryRoot'
$comparison = if ($IsWindows) { [StringComparison]::OrdinalIgnoreCase } else { [StringComparison]::Ordinal }
if ($project.Equals($library, $comparison)) { Fail 'ProjectRoot and LibraryRoot must be different directories' }

foreach ($required in @('AGENTS.md','AGENT_SYSTEM_START.md','START_ANY_PROJECT.md','VERIFY_EXECUTABLE_LIBRARY.ps1','markdown_system/PROJECT_START_READINESS_GATE.md','markdown_system/CAPABILITY_CATALOG.md')) {
  if (-not (Test-Path -LiteralPath (Join-Path $library $required) -PathType Leaf)) { Fail "LibraryRoot is incomplete; missing $required" }
}

$relativeLibrary = [IO.Path]::GetRelativePath($project, $library).Replace('\','/')
if (-not $relativeLibrary -or $relativeLibrary -eq '.') { Fail 'could not derive a distinct relative library path' }
if ($relativeLibrary.IndexOfAny([char[]]@("`r","`n",'`')) -ge 0) { Fail 'relative LibraryRoot contains Markdown control characters' }

$agentsBlock = @"
$begin
## Elite Engineering Library

- Biblioteca de autoridad: ``$relativeLibrary``.
- Para iniciar, continuar, diseñar, implementar, revisar o recuperar este sistema, usa la Skill ``elite-engineering-library`` y lee primero ``$relativeLibrary/AGENTS.md`` y ``$relativeLibrary/AGENT_SYSTEM_START.md``.
- Ejecuta su readiness gate y preflight; antes de cada ronda usa el renderer de asesoramiento, explica el prompt y registra la respuesta real. No escribas producto hasta ``READY_TO_BUILD``. Después avanza autónomamente con los packs y fuentes admitidos.
- No cargues toda la biblioteca ni atribuyas al upstream código local ``AUTHORED``. No inventes reglas, accesos, licencias, resultados de tests ni garantías.
- Si este proyecto no tiene raíz Git, inicia Codex desde esta raíz: la detección oficial sólo consulta el directorio actual cuando no encuentra una raíz de proyecto.
$end
"@

$claudeBlock = @"
$begin
@${relativeLibrary}/AGENTS.md

Usa ``${relativeLibrary}/AGENT_SYSTEM_START.md`` como router y ejecuta ``${relativeLibrary}/markdown_system/PROJECT_START_READINESS_GATE.md`` antes de programar.
$end
"@

$skill = @"
---
name: elite-engineering-library
description: Start, continue, design, implement, review, or recover a software project that explicitly uses the Elite Engineering Library. Do not use for unrelated projects.
---

$skillMarker

# Use Elite Engineering Library

The library root is ``$relativeLibrary`` relative to the project root. Treat it as read-only authority; write project artifacts in the project.

## Route the task

1. Read ``$relativeLibrary/AGENTS.md`` and ``$relativeLibrary/AGENT_SYSTEM_START.md`` completely.
2. Read ``$relativeLibrary/markdown_system/PROJECT_START_READINESS_GATE.md`` and execute it for a new project or as a delta audit. Materialize ``PROJECT_READINESS_GATE_PACK_PLAN.md``; before each pending round run ``render_project_advisory.py``, present its exact prompt conversationally, explain unfamiliar terms and record the real answer/evidence plus ``advisory_prompt_ref``. Ask only material follow-ups; never infer answers or request secrets in chat.
3. Run ``pwsh -NoProfile -File \"$relativeLibrary/VERIFY_EXECUTABLE_LIBRARY.ps1\" -Mode Preflight`` and record missing toolchains or access honestly.
4. Use ``$relativeLibrary/markdown_system/CAPABILITY_CATALOG.md`` and the two master maps to select only the relevant contracts, official source profiles and implementation packs. Do not load every Markdown file.
5. Do not write product code until ``PROJECT_READINESS_RECORD.md`` says ``READY_TO_BUILD`` with no applicable blocker. Once ready, implement and verify the smallest end-to-end vertical slice without ceremonial pauses.
6. When a failure occurs, follow the library failure-learning contract immediately: record it, find the root cause, correct the canonical source, add a regression, rebuild cleanly and continue. Never hide or waive a failed gate.

## Provenance boundary

This Skill and bridge are local ``AUTHORED`` orchestration based on OpenAI's official AGENTS.md discovery and progressive-disclosure Skill mechanisms. They are not upstream product code. Product code may come only from a compatible implementation pack or exact official source admitted by the library, with its own revision, hash, license, conditions and target verification.
"@

$plan = [Collections.Generic.List[object]]::new()
$installCodex = $Agent -in @('Codex','Both')
$installClaude = $Agent -in @('Claude','Both')

if ($installCodex) {
  $agentsPath = Join-Path $project 'AGENTS.md'
  $agentsExisting = if (Test-Path -LiteralPath $agentsPath -PathType Leaf) { Read-Utf8 $agentsPath } else { '' }
  $agentsContent = Merge-ManagedBlock $agentsExisting $agentsBlock $agentsPath '# Project agent instructions'
  $agentsBytes = $utf8.GetByteCount($agentsContent)
  if ($agentsBytes -gt $maximumAgentsBytes) { Fail "resulting AGENTS.md exceeds the official 32 KiB default: $agentsBytes bytes" }
  Add-PlannedWrite $plan $agentsPath $agentsContent

  $skillPath = Join-Path $project '.agents/skills/elite-engineering-library/SKILL.md'
  Add-PlannedWrite $plan $skillPath ($skill.TrimEnd() + "`n") $true
}

if ($installClaude) {
  $claudePath = Join-Path $project 'CLAUDE.md'
  $claudeExisting = if (Test-Path -LiteralPath $claudePath -PathType Leaf) { Read-Utf8 $claudePath } else { '' }
  $claudeContent = Merge-ManagedBlock $claudeExisting $claudeBlock $claudePath '# Claude Code instructions'
  Add-PlannedWrite $plan $claudePath $claudeContent

  $claudeSkillPath = Join-Path $project '.claude/skills/elite-engineering-library/SKILL.md'
  Add-PlannedWrite $plan $claudeSkillPath ($skill.TrimEnd() + "`n") $true
}

Invoke-PlannedWrites $plan

$receipt = [ordered]@{
  status = if ($WhatIfPreference) { 'ELITE_AGENT_BRIDGE_PLANNED' } else { 'ELITE_AGENT_BRIDGE_READY' }
  project_root = $project
  library_root = $library
  library_root_relative = $relativeLibrary
  agent = $Agent
  managed_files = @($plan | ForEach-Object { [IO.Path]::GetRelativePath($project, $_.Path).Replace('\','/') })
  agents_md_bytes = if ($installCodex) { $agentsBytes } else { 0 }
  claude_md_lines = if ($installClaude) { @($claudeContent -split '\r?\n').Count } else { 0 }
  codex_start_directory = if ($installCodex) { $project } else { $null }
}
$receipt | ConvertTo-Json -Depth 5
