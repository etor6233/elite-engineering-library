#requires -Version 7.0

[CmdletBinding()]
param()

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
$libraryRoot = [IO.Path]::GetFullPath((Join-Path $PSScriptRoot '..')).TrimEnd([IO.Path]::DirectorySeparatorChar)
$installer = Join-Path $libraryRoot 'INSTALL_AGENT_BRIDGE.ps1'
$tempRoot = Join-Path ([IO.Path]::GetTempPath()) ('elite-agent-bridge-test-' + [guid]::NewGuid().ToString('N'))
$utf8 = [Text.UTF8Encoding]::new($false)

function Assert-True([bool] $Condition, [string] $Message) {
  if (-not $Condition) { throw "TEST_INSTALL_AGENT_BRIDGE_FAILED: $Message" }
}

function New-TestProject([string] $Name) {
  $path = Join-Path $tempRoot $Name
  [void](New-Item -ItemType Directory -Path $path)
  $path
}

function Invoke-Installer([string] $Project, [string] $Agent = 'Both') {
  $output = @(& $installer -ProjectRoot $Project -LibraryRoot $libraryRoot -Agent $Agent)
  ($output -join "`n") | ConvertFrom-Json -Depth 10
}

function Expect-Failure([scriptblock] $Action, [string] $Pattern) {
  $failed = $false
  try { & $Action | Out-Null } catch {
    $failed = $true
    Assert-True ($_.Exception.Message -match $Pattern) "unexpected failure: $($_.Exception.Message)"
  }
  Assert-True $failed "expected failure matching $Pattern"
}

try {
  [void](New-Item -ItemType Directory -Path $tempRoot)

  $empty = New-TestProject 'empty-no-git'
  $first = Invoke-Installer $empty
  Assert-True ($first.status -eq 'ELITE_AGENT_BRIDGE_READY') 'empty project receipt status mismatch'
  Assert-True (-not (Test-Path -LiteralPath (Join-Path $empty '.git'))) 'installer created a Git repository'
  foreach ($path in @('AGENTS.md','CLAUDE.md','.agents/skills/elite-engineering-library/SKILL.md','.claude/skills/elite-engineering-library/SKILL.md')) {
    Assert-True (Test-Path -LiteralPath (Join-Path $empty $path) -PathType Leaf) "missing installed file: $path"
  }
  $agents = [IO.File]::ReadAllText((Join-Path $empty 'AGENTS.md'), $utf8)
  $claude = [IO.File]::ReadAllText((Join-Path $empty 'CLAUDE.md'), $utf8)
  $skill = [IO.File]::ReadAllText((Join-Path $empty '.agents/skills/elite-engineering-library/SKILL.md'), $utf8)
  $claudeSkill = [IO.File]::ReadAllText((Join-Path $empty '.claude/skills/elite-engineering-library/SKILL.md'), $utf8)
  $expectedRelative = [IO.Path]::GetRelativePath($empty, $libraryRoot).Replace('\','/')
  Assert-True ($agents.Contains("``$expectedRelative``", [StringComparison]::Ordinal)) 'AGENTS bridge has wrong relative library path'
  Assert-True ($claude.Contains("@$expectedRelative/AGENTS.md", [StringComparison]::Ordinal)) 'CLAUDE import has wrong relative library path'
  Assert-True ($skill.StartsWith("---`nname: elite-engineering-library`n", [StringComparison]::Ordinal)) 'skill frontmatter invalid'
  Assert-True ($skill.Contains('PROJECT_START_READINESS_GATE.md', [StringComparison]::Ordinal)) 'skill does not route readiness gate'
  Assert-True ($skill.Contains('render_project_advisory.py', [StringComparison]::Ordinal)) 'skill does not route executable advisory prompts'
  Assert-True ($skill.Contains('AUTHORED', [StringComparison]::Ordinal)) 'skill omits provenance boundary'
  Assert-True ($claudeSkill -ceq $skill) 'Codex and Claude skills diverged'
  Assert-True ([Text.Encoding]::UTF8.GetByteCount($agents) -le 32768) 'installed AGENTS exceeds 32 KiB'

  $before = @{}
  foreach ($path in @('AGENTS.md','CLAUDE.md','.agents/skills/elite-engineering-library/SKILL.md','.claude/skills/elite-engineering-library/SKILL.md')) {
    $before[$path] = (Get-FileHash -LiteralPath (Join-Path $empty $path) -Algorithm SHA256).Hash
  }
  [void](Invoke-Installer $empty)
  foreach ($path in $before.Keys) {
    Assert-True ((Get-FileHash -LiteralPath (Join-Path $empty $path) -Algorithm SHA256).Hash -eq $before[$path]) "idempotency failed: $path"
  }

  $existing = New-TestProject 'existing-instructions'
  [IO.File]::WriteAllText((Join-Path $existing 'AGENTS.md'), "# Existing`n`n- preserve-me`n", $utf8)
  [IO.File]::WriteAllText((Join-Path $existing 'CLAUDE.md'), "# Existing Claude`n`nKeep this.`n", $utf8)
  [void](Invoke-Installer $existing)
  Assert-True ([IO.File]::ReadAllText((Join-Path $existing 'AGENTS.md'), $utf8).Contains('- preserve-me', [StringComparison]::Ordinal)) 'existing AGENTS content was lost'
  Assert-True ([IO.File]::ReadAllText((Join-Path $existing 'CLAUDE.md'), $utf8).Contains('Keep this.', [StringComparison]::Ordinal)) 'existing CLAUDE content was lost'

  $codexOnly = New-TestProject 'codex-only'
  [void](Invoke-Installer $codexOnly 'Codex')
  Assert-True (Test-Path -LiteralPath (Join-Path $codexOnly 'AGENTS.md')) 'Codex-only did not create AGENTS'
  Assert-True (-not (Test-Path -LiteralPath (Join-Path $codexOnly 'CLAUDE.md'))) 'Codex-only created CLAUDE'
  Assert-True (-not (Test-Path -LiteralPath (Join-Path $codexOnly '.claude'))) 'Codex-only created Claude skill'

  $claudeOnly = New-TestProject 'claude-only'
  [void](Invoke-Installer $claudeOnly 'Claude')
  Assert-True (Test-Path -LiteralPath (Join-Path $claudeOnly 'CLAUDE.md')) 'Claude-only did not create CLAUDE'
  Assert-True (Test-Path -LiteralPath (Join-Path $claudeOnly '.claude/skills/elite-engineering-library/SKILL.md')) 'Claude-only did not create Claude skill'
  Assert-True (-not (Test-Path -LiteralPath (Join-Path $claudeOnly 'AGENTS.md'))) 'Claude-only created AGENTS'
  Assert-True (-not (Test-Path -LiteralPath (Join-Path $claudeOnly '.agents'))) 'Claude-only created Codex skill'

  $grokOnly = New-TestProject 'grok-only'
  [void](Invoke-Installer $grokOnly 'Grok')
  Assert-True (Test-Path -LiteralPath (Join-Path $grokOnly 'AGENTS.md')) 'Grok-only did not create AGENTS'
  Assert-True (Test-Path -LiteralPath (Join-Path $grokOnly '.grok/skills/elite-engineering-library/SKILL.md')) 'Grok-only did not create Grok skill'
  $grokSkillOnly = [IO.File]::ReadAllText((Join-Path $grokOnly '.grok/skills/elite-engineering-library/SKILL.md'), $utf8)
  Assert-True ($grokSkillOnly.Contains('PACK_PER_CLAIM_INDEX.md', [StringComparison]::Ordinal)) 'Grok skill does not route pack-per-claim index'
  Assert-True ($grokSkillOnly.Contains('Do not load wholesale', [StringComparison]::Ordinal)) 'Grok skill omits wholesale guard'
  Assert-True (-not (Test-Path -LiteralPath (Join-Path $grokOnly 'CLAUDE.md'))) 'Grok-only created CLAUDE'
  Assert-True (-not (Test-Path -LiteralPath (Join-Path $grokOnly '.agents'))) 'Grok-only created Codex skill'

  $allAgents = New-TestProject 'all-agents'
  [void](Invoke-Installer $allAgents 'All')
  foreach ($path in @('AGENTS.md','CLAUDE.md','.agents/skills/elite-engineering-library/SKILL.md','.claude/skills/elite-engineering-library/SKILL.md','.grok/skills/elite-engineering-library/SKILL.md')) {
    Assert-True (Test-Path -LiteralPath (Join-Path $allAgents $path) -PathType Leaf) "All-agents missing: $path"
  }

  $vendored = New-TestProject 'vendored-library'
  $vendoredLibrary = Join-Path $vendored 'tools/elite-engineering-library'
  foreach ($required in @('AGENTS.md','GROK_AGENT_ENTRY.md','AGENT_SYSTEM_START.md','START_ANY_PROJECT.md','VERIFY_EXECUTABLE_LIBRARY.ps1','markdown_system/PROJECT_START_READINESS_GATE.md','markdown_system/CAPABILITY_CATALOG.md','markdown_system/PACK_PER_CLAIM_INDEX.md')) {
    $target = Join-Path $vendoredLibrary $required
    [void](New-Item -ItemType Directory -Path (Split-Path -Parent $target) -Force)
    [IO.File]::WriteAllText($target, "fixture for $required`n", $utf8)
  }
  $vendoredReceipt = @(& $installer -ProjectRoot $vendored -LibraryRoot $vendoredLibrary -Agent Codex) -join "`n" | ConvertFrom-Json -Depth 10
  Assert-True ($vendoredReceipt.library_root_relative -eq 'tools/elite-engineering-library') 'vendored library relative path mismatch'
  Assert-True ([IO.File]::ReadAllText((Join-Path $vendored 'AGENTS.md'), $utf8).Contains('`tools/elite-engineering-library`', [StringComparison]::Ordinal)) 'vendored AGENTS bridge path mismatch'

  $damaged = New-TestProject 'damaged-marker'
  $damagedAgents = "# Existing`n`n<!-- ELITE-ENGINEERING-LIBRARY:BEGIN -->`nbroken`n"
  [IO.File]::WriteAllText((Join-Path $damaged 'AGENTS.md'), $damagedAgents, $utf8)
  Expect-Failure { & $installer -ProjectRoot $damaged -LibraryRoot $libraryRoot -Agent Codex } 'markers are damaged'
  Assert-True ([IO.File]::ReadAllText((Join-Path $damaged 'AGENTS.md'), $utf8) -ceq $damagedAgents) 'damaged marker failure modified AGENTS'
  Assert-True (-not (Test-Path -LiteralPath (Join-Path $damaged '.agents'))) 'damaged marker failure created skill directory'

  $unmanaged = New-TestProject 'unmanaged-skill'
  $unmanagedSkillPath = Join-Path $unmanaged '.agents/skills/elite-engineering-library/SKILL.md'
  [void](New-Item -ItemType Directory -Path (Split-Path -Parent $unmanagedSkillPath) -Force)
  [IO.File]::WriteAllText($unmanagedSkillPath, "---`nname: elite-engineering-library`ndescription: foreign`n---`n", $utf8)
  Expect-Failure { & $installer -ProjectRoot $unmanaged -LibraryRoot $libraryRoot -Agent Codex } 'unmanaged skill'
  Assert-True (-not (Test-Path -LiteralPath (Join-Path $unmanaged 'AGENTS.md'))) 'unmanaged skill failure partially wrote AGENTS'

  $oversized = New-TestProject 'oversized-agents'
  $oversizedText = '# Existing' + "`n" + ('x' * 33000)
  [IO.File]::WriteAllText((Join-Path $oversized 'AGENTS.md'), $oversizedText, $utf8)
  Expect-Failure { & $installer -ProjectRoot $oversized -LibraryRoot $libraryRoot -Agent Codex } 'exceeds the official 32 KiB'
  Assert-True ([IO.File]::ReadAllText((Join-Path $oversized 'AGENTS.md'), $utf8) -ceq $oversizedText) 'oversize failure modified AGENTS'
  Assert-True (-not (Test-Path -LiteralPath (Join-Path $oversized '.agents'))) 'oversize failure created skill directory'

  $whatIf = New-TestProject 'what-if'
  & $installer -ProjectRoot $whatIf -LibraryRoot $libraryRoot -Agent Both -WhatIf -InformationAction Ignore *>$null
  Assert-True (@(Get-ChildItem -LiteralPath $whatIf -Force).Count -eq 0) 'WhatIf mutated project'

  Expect-Failure { & $installer -ProjectRoot $libraryRoot -LibraryRoot $libraryRoot -Agent Codex } 'must be different directories'

  if ($IsWindows) {
    & (Join-Path $PSScriptRoot 'test_agent_bridge_safety.ps1') -LibraryRoot $libraryRoot
  } else { Write-Warning 'Windows file-sharing rollback safety suite not executed on this host' }
  Write-Output 'PASS: install_agent_bridge preserves existing instructions, installs Codex/Claude/Grok progressive context without Git, is idempotent, enforces 32 KiB, rejects conflicts before writing and restores file bytes on tested write failures'
} finally {
  if (Test-Path -LiteralPath $tempRoot) {
    $resolved = [IO.Path]::GetFullPath($tempRoot)
    $systemTemp = [IO.Path]::GetFullPath([IO.Path]::GetTempPath())
    if ($resolved.StartsWith($systemTemp, [StringComparison]::OrdinalIgnoreCase) -and (Split-Path -Leaf $resolved) -like 'elite-agent-bridge-test-*') {
      Remove-Item -LiteralPath $resolved -Recurse -Force
    }
  }
}
