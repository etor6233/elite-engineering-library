#requires -Version 7.0
[CmdletBinding()]
param([Parameter(Mandatory)][string] $Materializer,[Parameter(Mandatory)][string] $Composer,[Parameter(Mandatory)][string] $Updater,[string] $OfficialLicense,[Parameter(Mandatory)][string] $EvidenceRoot)
Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
$utf8 = [Text.UTF8Encoding]::new($false)
if (Test-Path -LiteralPath $EvidenceRoot) { throw 'Evidence destination must be absent' }
[void][IO.Directory]::CreateDirectory($EvidenceRoot)
function Hash([string] $Text) { [Convert]::ToHexString([Security.Cryptography.SHA256]::HashData($utf8.GetBytes($Text))).ToLowerInvariant() }
function Put([string] $Path, [string] $Text) { [IO.File]::WriteAllText($Path,$Text,$utf8) }
function Assert([bool] $Condition,[string] $Message) { if (-not $Condition) { throw $Message } }
function Pack([string] $Text,[string] $Flag) {
  $payload = $Text
  if ($payload.Length -gt 0 -and -not $payload.EndsWith("`n")) { $payload += "`n" }
  $fence = '`' * 4
  $template = @'
# Exact EOF fixture
## 1. Metadata
```yaml
pack_id: "EXACT-EOF"
pack_version: "1.0.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: RECONSTRUCTIBLE
  admission: CONDITIONED
claim: "exact bytes fixture"
stacks: ["text"]
compatible_with: []
incompatible_with: []
license_expression: "LicenseRef-Test"
upstream_sources: []
verified_at: "2026-09-12"
```
## 2. Exact file manifest
```text
CREATE artifact.txt
```
## 3. Materialization blocks
### FILE: `artifact.txt`
```yaml
block_id: "EXACT-EOF:artifact:v1"
operation: CREATE
provenance: AUTHORED
source: "synthetic test; optional official license inspected separately"
license: "LicenseRef-Test"
sha256: "__HASH__"
variables: []
secrets_allowed: false
__FLAG__```
__FENCE__text
__PAYLOAD____FENCE__
'@
  $template.Replace('__HASH__',(Hash $Text)).Replace('__FLAG__',$Flag).Replace('__FENCE__',$fence).Replace('__PAYLOAD__',$payload)
}
$results = [Collections.Generic.List[object]]::new()
function Run-Case([string] $Name,[string] $Text,[string] $Markdown,[string] $Failure = '') {
  $caseRoot = Join-Path $EvidenceRoot $Name
  [void][IO.Directory]::CreateDirectory($caseRoot)
  Put (Join-Path $caseRoot 'pack.md') $Markdown
  $plan = @{ compositionVersion='1.0';secretVariables=@{};packs=@(@{path='pack.md';packId='EXACT-EOF';version='1.0.0';files=@('*');variables=@{};acknowledgeConditions=$true}) }
  Put (Join-Path $caseRoot 'plan.json') ($plan | ConvertTo-Json -Depth 8)
  foreach ($tool in @('materializer','composer')) {
    $destination = Join-Path $caseRoot $tool
    $caught = ''
    try {
      if ($tool -eq 'materializer') { & $Materializer -PackFile (Join-Path $caseRoot 'pack.md') -Destination $destination | Out-Null }
      else { & $Composer -PlanFile (Join-Path $caseRoot 'plan.json') -LibraryRoot $caseRoot -Destination $destination | Out-Null }
    } catch { $caught = $_.Exception.Message }
    if ($Failure) {
      Assert ($caught -match $Failure) "$Name/$tool unexpected error: $caught"
      Assert (-not (Test-Path -LiteralPath $destination)) "$Name/$tool wrote before validation"
    } else {
      Assert (-not $caught) "$Name/$tool failed: $caught"
      $bytes = [IO.File]::ReadAllBytes((Join-Path $destination 'artifact.txt'))
      Assert ([Convert]::ToHexString($bytes) -ceq [Convert]::ToHexString($utf8.GetBytes($Text))) "$Name/$tool changed exact bytes"
    }
    $results.Add(@{case=$Name;tool=$tool;status='PASS';expected_failure=$Failure;sha256=(Hash $Text)})
  }
}
foreach ($case in @(
  @{name='legacy-lf';text="alpha`n";flag=''},
  @{name='multiple-lf';text="alpha`n`n";flag=''},
  @{name='explicit-true';text="alpha`n";flag="final_newline: true`n"},
  @{name='no-lf';text="alpha`nbeta";flag="final_newline: false`n"},
  @{name='empty-legacy';text='';flag=''},
  @{name='empty-false';text='';flag="final_newline: false`n"},
  @{name='single-lf';text="`n";flag=''}
)) { Run-Case $case.name $case.text (Pack $case.text $case.flag) }
$noLF = Pack 'alpha' "final_newline: false`n"
Run-Case 'no-lf-crlf-envelope' 'alpha' ($noLF.Replace("`n","`r`n"))
Run-Case 'tamper' 'alpha' ($noLF.Replace("alpha`n", "changed`n")) 'SHA-256 mismatch'
Run-Case 'missing-flag' 'alpha' (Pack 'alpha' '') 'SHA-256 mismatch'
Run-Case 'invalid-flag' 'alpha' (Pack 'alpha' "final_newline: yes`n") 'final_newline must'
Run-Case 'quoted-flag' 'alpha' (Pack 'alpha' "final_newline: ''false''`n") 'final_newline must'
Run-Case 'duplicate-flag' 'alpha' (Pack 'alpha' "final_newline: false`nfinal_newline: true`n") 'Duplicate final_newline'
Run-Case 'contradictory-false' "alpha`n`n" (Pack "alpha`n`n" "final_newline: false`n") 'conflicts with trailing'
$updateRoot = Join-Path $EvidenceRoot 'updater'
[void][IO.Directory]::CreateDirectory($updateRoot)
$packPath = Join-Path $updateRoot 'pack.md'
Put $packPath (Pack "old`n" '')
$i=0
foreach ($text in @('no final LF', '', "one`n", "many`n`n", 'second no LF', '')) {
  Put (Join-Path $updateRoot 'artifact.txt') $text
  & $Updater -PackFile $packPath -SourceRoot $updateRoot -Files @('artifact.txt') | Out-Null
  Run-Case "updater-$i" $text ([IO.File]::ReadAllText($packPath,$utf8))
  $i++
}
if ($OfficialLicense) {
  $license = [IO.File]::ReadAllText((Resolve-Path -LiteralPath $OfficialLicense).Path,$utf8)
  Assert ((Hash $license) -ceq 'abc09dad5f84a76e1b0279237053cae16c03228ab27d8d467677054c2bd17eeb') 'official license hash changed'
  Assert ($utf8.GetByteCount($license) -eq 43529) 'official license size changed'
  Run-Case 'official-lgpl' $license (Pack $license "final_newline: false`n")
}
Put (Join-Path $EvidenceRoot 'result.json') (@{status='PASS';results=@($results.ToArray())} | ConvertTo-Json -Depth 8)
Write-Output "PASS: $($results.Count) exact EOF reconstruction/negative/updater cases"
