#requires -Version 7.0
# AUTHORED integration regression, not Google/NASA code or a product readiness grant.
[CmdletBinding()]
param([string] $LibraryRoot = (Join-Path $PSScriptRoot '..'), [switch] $KeepEvidence)
Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
$LibraryRoot = [IO.Path]::GetFullPath($LibraryRoot)
$pwshPath = (Get-Process -Id $PID).Path
$pythonPath = (Get-Command python -CommandType Application -ErrorAction Stop | Select-Object -First 1).Source
$testRoot = Join-Path ([IO.Path]::GetTempPath()) ('elite-entry-lifecycle-' + [guid]::NewGuid().ToString('N'))
[void](New-Item -ItemType Directory -Path $testRoot)
$utf8 = [Text.UTF8Encoding]::new($false)
$checks = 0
$calls = [Collections.Generic.List[object]]::new()
function Check([bool] $Value, [string] $Message) {
  if (-not $Value) { throw $Message }
  $script:checks++
}
function Run([string] $Id, [string] $Executable, [string[]] $Arguments, [int] $Expected = 0, [string] $Pattern = '') {
  $watch = [Diagnostics.Stopwatch]::StartNew()
  $output = @(& $Executable @Arguments 2>&1)
  $code = $LASTEXITCODE
  $watch.Stop()
  $message = $output -join "`n"
  $calls.Add([ordered]@{ id=$Id; exit_code=$code; expected_exit=$Expected; elapsed_ms=$watch.ElapsedMilliseconds; output=$message })
  Check ($code -eq $Expected -and (-not $Pattern -or $message -match $Pattern)) "$Id unexpected exit/output: $code $message"
}
function Save-Json([string] $Path, [object] $Value) {
  [IO.File]::WriteAllText($Path, (($Value | ConvertTo-Json -Depth 30) + "`n"), $utf8)
}
function Digest([string] $Path) { (Get-FileHash -LiteralPath $Path -Algorithm SHA256).Hash.ToLowerInvariant() }
try {
  $materializer = Join-Path $LibraryRoot 'materialize_markdown_pack.ps1'
  $compositorRoot = Join-Path $testRoot 'compositor'
  $executionRoot = Join-Path $testRoot 'execution'
  Run 'materialize-compositor' $pwshPath @('-NoProfile','-File',$materializer,'-PackFile',(Join-Path $LibraryRoot 'implementation_packs/MARKDOWN_COMPOSITOR_CORE.md'),'-Destination',$compositorRoot)
  Run 'materialize-execution' $pwshPath @('-NoProfile','-File',$materializer,'-PackFile',(Join-Path $LibraryRoot 'implementation_packs/ENGINEERING_EXECUTION_VALIDATOR.md'),'-Destination',$executionRoot)
  $compositor = Join-Path $compositorRoot 'tools/compose-markdown-project.ps1'
  $execution = Join-Path $executionRoot 'engineering_execution_kit'
  foreach ($mode in @('NEW','EXISTING')) {
    $project = Join-Path $testRoot $mode
    [void](New-Item -ItemType Directory -Path $project)
    $owned = Join-Path $project 'existing-user-code.txt'
    if ($mode -eq 'EXISTING') { [IO.File]::WriteAllText($owned, "synthetic user-owned bytes`r`n", [Text.UTF8Encoding]::new($true)) }
    $inventory = @(Get-ChildItem -LiteralPath $project -File -Recurse -Force | ForEach-Object { [ordered]@{ path=[IO.Path]::GetRelativePath($project,$_.FullName); sha256=(Digest $_.FullName) } })
    Check ($inventory.Count -eq $(if ($mode -eq 'NEW') { 0 } else { 1 })) "$mode baseline not observed before writes"
    $ownedHash = if ($mode -eq 'EXISTING') { Digest $owned } else { '' }
    # Compose in a fresh tooling subdirectory, never into existing product files.
    $tooling = Join-Path $project '.elite-tools'
    $composeArgs = @('-NoProfile','-File',$compositor,'-PlanFile',(Join-Path $LibraryRoot 'markdown_system/PROJECT_READINESS_GATE_PACK_PLAN.md'),'-LibraryRoot',$LibraryRoot,'-Destination',$tooling)
    Run "$mode-compose" $pwshPath $composeArgs
    Check (Test-Path -LiteralPath (Join-Path $tooling 'MATERIALIZATION_RECORD.md') -PathType Leaf) "$mode missing composition receipt"
    $bridgeArgs = @('-NoProfile','-File',(Join-Path $LibraryRoot 'INSTALL_AGENT_BRIDGE.ps1'),'-ProjectRoot',$project,'-LibraryRoot',$LibraryRoot,'-Agent','Both')
    Run "$mode-bridge" $pwshPath $bridgeArgs
    foreach ($path in @('AGENTS.md','CLAUDE.md','.agents/skills/elite-engineering-library/SKILL.md','.claude/skills/elite-engineering-library/SKILL.md')) {
      Check (Test-Path -LiteralPath (Join-Path $project $path) -PathType Leaf) "$mode missing bridge artifact: $path"
    }
    $gate = Join-Path $tooling 'project_readiness_gate'
    Copy-Item -LiteralPath (Join-Path $gate 'project-readiness.template.json') -Destination (Join-Path $project 'PROJECT_READINESS_GATE.json')
    Run "$mode-advisory" $pythonPath @('-B',(Join-Path $gate 'render_project_advisory.py'),'--project-root',$project,'--round','A','--output','PROJECT_ADVISORY_A.md') 0 'PROJECT_ADVISORY_READY'
    Run "$mode-readiness-blocked" $pythonPath @('-B',(Join-Path $gate 'validate_project_readiness.py'),'--project-root',$project,'--report','PROJECT_READINESS_REPORT.json') 2 'BLOCKED'
    $report = Get-Content -LiteralPath (Join-Path $project 'PROJECT_READINESS_REPORT.json') -Raw | ConvertFrom-Json
    Check ($report.status -ceq 'BLOCKED' -and $report.errors.Count -gt 0) "$mode inherited readiness approval"
    $failures = Join-Path $project 'PROJECT_FAILURE_LESSONS.md'
    [IO.File]::WriteAllText($failures, "# Synthetic lifecycle fixture; no product approvals`n", $utf8)
    Save-Json (Join-Path $project 'baseline.json') $inventory
    [IO.File]::WriteAllText((Join-Path $project 'delta.md'), "# Synthetic authorized delta: add local tooling; preserve existing user bytes; no product or external effects.`n", $utf8)
    $state = Get-Content -LiteralPath (Join-Path $execution 'execution_state.template.json') -Raw | ConvertFrom-Json -AsHashtable
    $state.revision = 1
    $state.project = @{id="entry-$($mode.ToLowerInvariant())"; mode=$mode; rigor='HIGH'}
    $state.baseline = @{kind=$(if($mode -eq 'NEW'){'EMPTY'}else{'EXISTING'}); captured_at='2026-09-07T00:00:00Z'; inventory_ref='inventory'; delta_scope_ref=$(if($mode -eq 'EXISTING'){'delta'}else{$null})}
    $state.context = @{summary='Synthetic CLI entry assembled; real readiness remains BLOCKED.'; must_read_refs=@('report','failures'); reuse_without_reload_refs=@('inventory','delta')}
    $state.links.failure_ledger = 'failures'
    $state.links.readiness_report = 'report'
    $state.next_action = 'Resolve actual intake; never promote this synthetic blocked fixture.'
    $state.evidence = @()
    foreach ($pair in @(@('failures','PROJECT_FAILURE_LESSONS.md'),@('report','PROJECT_READINESS_REPORT.json'),@('inventory','baseline.json'),@('delta','delta.md'))) {
      $state.evidence += @{id=$pair[0];kind='fixture';path=$pair[1];sha256=(Digest (Join-Path $project $pair[1]));produced_at='2026-09-07T00:00:00Z';producer='synthetic-integration-test'}
    }
    $statePath = Join-Path $project 'PROJECT_EXECUTION_STATE.json'
    $eventsPath = Join-Path $project 'PROJECT_EXECUTION_EVENTS.jsonl'
    Save-Json $statePath $state
    $checkpointArgs = @('-B',(Join-Path $execution 'checkpoint_execution_state.py'),$statePath,'--project-root',$project,'--events',$eventsPath,'--event-id','EVT-0001','--event-type','BASELINE_CAPTURED','--actor','test','--summary','Synthetic entry BLOCKED')
    Run "$mode-checkpoint" $pythonPath $checkpointArgs 0 'EXECUTION_CHECKPOINT_PASS'
    $resumeArgs = @('-B',(Join-Path $execution 'validate_execution_state.py'),$statePath,'--project-root',$project,'--events',$eventsPath,'--level','resume')
    Run "$mode-resume" $pythonPath $resumeArgs 0 'EXECUTION_STATE_PASS'
    # Existing destination rejection must not mutate any composed file.
    $before = @(Get-ChildItem -LiteralPath $tooling -File -Recurse -Force | ForEach-Object { $_.FullName + ':' + (Digest $_.FullName) })
    Run "$mode-compose-collision" $pwshPath $composeArgs 1 'already exists|must not exist|must be absent|must be empty'
    $after = @(Get-ChildItem -LiteralPath $tooling -File -Recurse -Force | ForEach-Object { $_.FullName + ':' + (Digest $_.FullName) })
    Check (($before -join "`n") -ceq ($after -join "`n")) "$mode collision changed generated bytes"
    $bridgeHash = Digest (Join-Path $project 'AGENTS.md')
    Run "$mode-bridge-repeat" $pwshPath $bridgeArgs
    Check ((Digest (Join-Path $project 'AGENTS.md')) -ceq $bridgeHash) "$mode bridge not idempotent"
    $eventBytes = [IO.File]::ReadAllBytes($eventsPath)
    Run "$mode-checkpoint-repeat" $pythonPath $checkpointArgs 2 'event_id already exists|already has'
    Check (([Convert]::ToBase64String([IO.File]::ReadAllBytes($eventsPath))) -ceq ([Convert]::ToBase64String($eventBytes))) "$mode rejected checkpoint changed history"
    # Tamper only a test-owned artifact; restoration uses the saved exact bytes.
    $original = [IO.File]::ReadAllBytes($failures)
    [IO.File]::AppendAllText($failures, 'synthetic corruption', $utf8)
    Run "$mode-tamper-rejected" $pythonPath $resumeArgs 2 'evidence hash mismatch'
    [IO.File]::WriteAllBytes($failures, $original)
    Run "$mode-recovered" $pythonPath $resumeArgs 0 'EXECUTION_STATE_PASS'
    $state.revision = 2
    $state.next_action = 'Synthetic recovery verified; actual intake is still BLOCKED.'
    Save-Json $statePath $state
    $nextArgs = $checkpointArgs.Clone()
    $nextArgs[$nextArgs.IndexOf('EVT-0001')] = 'EVT-0002'
    $nextArgs[$nextArgs.IndexOf('BASELINE_CAPTURED')] = 'STATE_ADVANCED'
    Run "$mode-checkpoint-advance" $pythonPath $nextArgs 0 'EXECUTION_CHECKPOINT_PASS'
    Run "$mode-resume-advanced" $pythonPath $resumeArgs 0 'EXECUTION_STATE_PASS'
    Check (@(Get-Content -LiteralPath $eventsPath).Count -eq 2) "$mode event sequence not preserved"
    Check ((Get-Content -LiteralPath $eventsPath -First 1) -ceq ([Text.Encoding]::UTF8.GetString($eventBytes).TrimEnd())) "$mode first event rewritten"
    Check ($state.status -ceq 'BLOCKED' -and $state.release.status -ceq 'NONE') "$mode synthetic test promoted product"
    if ($mode -eq 'EXISTING') { Check ((Digest $owned) -ceq $ownedHash) 'Existing user bytes/BOM changed' }
    Write-Output "LIFECYCLE mode=$mode status=BLOCKED events=2 product_approval=false"
  }
  Save-Json (Join-Path $testRoot 'evidence.json') @{schema='elite-entry-lifecycle-test/v1';status='PASS';checks=$checks;calls=@($calls);product_readiness=$false}
  Write-Output "PASS: agent entry lifecycle $checks checks; NEW/EXISTING; product_readiness=false"
} finally {
  if ($KeepEvidence) { Write-Output "LIFECYCLE_EVIDENCE=$testRoot" }
  else {
    $resolved = [IO.Path]::GetFullPath($testRoot)
    $parent = [IO.Path]::GetFullPath([IO.Path]::GetTempPath()).TrimEnd([IO.Path]::DirectorySeparatorChar)
    if ((Split-Path -Parent $resolved) -cne $parent -or (Split-Path -Leaf $resolved) -notlike 'elite-entry-lifecycle-*') { throw 'Unsafe lifecycle cleanup target' }
    Remove-Item -LiteralPath $resolved -Recurse -Force
  }
}
exit 0
