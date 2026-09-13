#requires -Version 7.0
# AUTHORED regression of the actual selectors against public evidence and local records.
[CmdletBinding()]
param()
$ErrorActionPreference = 'Stop'
$sourceRoot = Split-Path -Parent $PSScriptRoot
$testRoot = Join-Path ([IO.Path]::GetTempPath()) ('elite-public-evidence-' + [guid]::NewGuid().ToString('N'))
[void][IO.Directory]::CreateDirectory($testRoot)
$utf8 = [Text.UTF8Encoding]::new($false)
$results = [Collections.Generic.List[object]]::new()
foreach ($scriptName in @('VERIFY_LIBRARY.ps1','CREATE_PORTABLE_ARCHIVE.ps1')) {
  & {
    param($scriptName,$sourceRoot,$testRoot,$results,$utf8)
    $tokens = $null; $parseErrors = $null
    $ast = [Management.Automation.Language.Parser]::ParseFile((Join-Path $sourceRoot $scriptName),[ref]$tokens,[ref]$parseErrors)
    if ($parseErrors.Count) { throw "Parse error in $scriptName" }
    foreach ($name in @('releaseDirectories','releaseRootMarkdown','releaseRootScripts','releaseRootMetadata','releaseMarkdownSystemScripts','distributionManifestName')) {
      $nodes = @($ast.FindAll({param($n) $n -is [Management.Automation.Language.AssignmentStatementAst] -and $n.Left.Extent.Text -ceq ('$'+$name)},$false))
      if ($nodes.Count -gt 1) { throw "Ambiguous declaration $name" }
      if ($nodes.Count -eq 1) { . ([scriptblock]::Create($nodes[0].Extent.Text)) }
    }
    foreach ($name in @('Assert-True','Relative','Relative-ToLibrary','Get-ReleaseInputPolicy','Find-MachineLocalHomePath','Test-LocalMaintenanceEntry','Get-ReleaseSourceFiles','Assert-FranchiseProfileSelection','Assert-FranchiseProfileSelectionRegression')) {
      $nodes = @($ast.FindAll({param($n) $n -is [Management.Automation.Language.FunctionDefinitionAst] -and $n.Name -ceq $name},$false))
      if ($nodes.Count -gt 1) { throw "Ambiguous function $name" }
      if ($nodes.Count -eq 1) { . ([scriptblock]::Create($nodes[0].Extent.Text)) }
    }
    $cases = @('public-json','public-log','unlisted-json','unlisted-log','unlisted-executable','missing-file','duplicate-key','duplicate-name','unsafe-name','lowercase-name','unknown-field','no-policy','nested-evidence','policy-reparse','local-readiness','root-metadata','metadata-directory','metadata-near-name')
    foreach ($case in $cases) {
      $libraryRoot = Join-Path $testRoot ($scriptName + '-' + $case)
      [void][IO.Directory]::CreateDirectory((Join-Path $libraryRoot 'markdown_system'))
      [void][IO.Directory]::CreateDirectory((Join-Path $libraryRoot 'reconstruction_evidence'))
      [IO.File]::WriteAllText((Join-Path $libraryRoot 'README.md'),"fixture`n",$utf8)
      $policy = [ordered]@{schema='elite-public-release-input-policy/v1';json_files=@('PUBLIC_RECEIPT.json');text_logs=@();franchise_pack_ids=@('PACK-A');franchise_files=1}
      [IO.File]::WriteAllText((Join-Path $libraryRoot 'reconstruction_evidence/PUBLIC_RECEIPT.json'),'{"scope":"synthetic-public-library-evidence"}',$utf8)
      switch ($case) {
        'public-log' { $policy.text_logs=@('PUBLIC_TEST.log'); [IO.File]::WriteAllText((Join-Path $libraryRoot 'reconstruction_evidence/PUBLIC_TEST.log'),'fixture pass',$utf8) }
        'unlisted-json' { [IO.File]::WriteAllText((Join-Path $libraryRoot 'reconstruction_evidence/SECRET.json'),'{}',$utf8) }
        'unlisted-log' { [IO.File]::WriteAllText((Join-Path $libraryRoot 'reconstruction_evidence/UNKNOWN.log'),'not selected',$utf8) }
        'unlisted-executable' { [IO.File]::WriteAllText((Join-Path $libraryRoot 'reconstruction_evidence/PAYLOAD.ps1'),'throw "not public"',$utf8) }
        'missing-file' { $policy.json_files+= 'MISSING.json' }
        'duplicate-name' { $policy.json_files+= 'PUBLIC_RECEIPT.json' }
        'unsafe-name' { $policy.json_files=@('../PUBLIC_RECEIPT.json') }
        'lowercase-name' { $policy.json_files=@('public_receipt.json') }
        'unknown-field' { $policy.extra=@() }
        'nested-evidence' { [void][IO.Directory]::CreateDirectory((Join-Path $libraryRoot 'reconstruction_evidence/nested')); [IO.File]::WriteAllText((Join-Path $libraryRoot 'reconstruction_evidence/nested/PUBLIC_RECEIPT.json'),'{}',$utf8) }
        'local-readiness' {
          foreach ($name in @('PROJECT_EXECUTION_STATE.json','PROJECT_EXECUTION_EVENTS.jsonl','PROJECT_LIBRARY_READINESS_GATE.json')) { [IO.File]::WriteAllText((Join-Path $libraryRoot $name),'{}',$utf8) }
        }
        'root-metadata' { foreach ($name in @('.gitignore','.gitattributes')) { [IO.File]::WriteAllText((Join-Path $libraryRoot $name),'synthetic',$utf8) } }
        'metadata-directory' { [void][IO.Directory]::CreateDirectory((Join-Path $libraryRoot '.gitattributes')) }
        'metadata-near-name' { [IO.File]::WriteAllText((Join-Path $libraryRoot '.gitattributes.extra'),'unexpected',$utf8) }
      }
      $json = $policy | ConvertTo-Json -Depth 6
      if ($case -eq 'duplicate-key') { $json = $json.Replace('"schema":','"schema":"duplicate", "schema":') }
      $policyText = "# Fixture`n`n" + ('`'*3) + "json`n$json`n" + ('`'*3) + "`n"
      $policyPath = Join-Path $libraryRoot 'markdown_system/PUBLIC_RELEASE_INPUT_POLICY.md'
      if ($case -eq 'policy-reparse') {
        $target = Join-Path $testRoot ('outside-'+$scriptName);[void][IO.Directory]::CreateDirectory($target)
        $linkType = if ($IsWindows) { 'Junction' } else { 'SymbolicLink' }
        [void](New-Item -ItemType $linkType -Path $policyPath -Target $target)
      } elseif ($case -ne 'no-policy') { [IO.File]::WriteAllText($policyPath,$policyText,$utf8) }
      $errorText = $null;$selected=@()
      try { $selected=@(Get-ReleaseSourceFiles) } catch { $errorText=$_.Exception.Message }
      if ($case -in @('public-json','public-log','local-readiness','root-metadata')) {
        $expected = switch ($case) { 'public-log' {4} 'root-metadata' {5} default {3} }
        if ($errorText -or $selected.Count -ne $expected) { throw "$scriptName/$case expected$expected public files: $errorText" }
        if (@($selected | Where-Object Name -Like 'PROJECT_*').Count) { throw "$scriptName/$case leaked local approval" }
      } else {
        $expected = switch ($case) {
          {$_ -in @('unlisted-json','unlisted-log','unlisted-executable','no-policy','nested-evidence')} { 'unknown release file';break }
          'missing-file' {'named evidence missing'}
          'duplicate-key' {'duplicate field'}
          {$_ -in @('duplicate-name','unsafe-name','lowercase-name')} {'unsafe/duplicate item';break}
          'unknown-field' {'unknown/missing field'}
          'policy-reparse' {'symlink/reparse point'}
          'metadata-directory' {'unknown top-level directory'}
          'metadata-near-name' {'unknown top-level file'}
        }
        if (-not $errorText -or $errorText -inotmatch [regex]::Escape($expected)) { throw "$scriptName/$case wrong rejection: $errorText (expected $expected)" }
      }
      $results.Add([pscustomobject]@{script=$scriptName;case=$case;state='PASS';selected=$selected.Count;rejection=$errorText})
    }
    if ($scriptName -eq 'VERIFY_LIBRARY.ps1') {
      $libraryRoot=$sourceRoot
      Assert-FranchiseProfileSelectionRegression
      $backslash = [char]92
      foreach ($entry in @(
        @{text=('C:'+$backslash+'Users'+$backslash+'fixture'+$backslash+'private');expected=$true},
        @{text=('c:'+$backslash+'users'+$backslash+'fixture');expected=$true},
        @{text=('/'+'home'+'/fixture/private');expected=$true},
        @{text=('/'+'Users'+'/fixture/private');expected=$true},
        @{text='/users/123456789';expected=$false},
        @{text='https://api.example.test/users/123456789';expected=$false}
      )) {
        $match = Find-MachineLocalHomePath $entry.text
        if ($match.Success -ne $entry.expected) { throw 'Machine home path classification regression' }
        $results.Add([pscustomobject]@{script=$scriptName;case='home-path-semantics';state='PASS';expected_match=$entry.expected})
      }

      $results.Add([pscustomobject]@{script=$scriptName;case='current-franchise-lock-and-negatives';state='PASS'})
    }
  } $scriptName $sourceRoot $testRoot $results $utf8
}
[IO.File]::WriteAllText((Join-Path $testRoot 'result.json'),(@{state='PASS';cases=@($results.ToArray())}|ConvertTo-Json -Depth 8),$utf8)
Write-Output "PASS: public release policy $($results.Count) cases; exact public evidence, locked franchise and local approvals excluded"
Write-Output "EVIDENCE_ROOT=$testRoot"
