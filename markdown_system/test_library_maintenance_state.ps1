#requires -Version 7.0
# AUTHORED: regression of the actual distribution selectors, not a provider sample.
$ErrorActionPreference = 'Stop'
$sourceRoot = Split-Path -Parent $PSScriptRoot
$testRoot = Join-Path ([IO.Path]::GetTempPath()) ('elite-state-distribution-' + [guid]::NewGuid().ToString('N'))
[void](New-Item -ItemType Directory -Path $testRoot)
$checks = 0
foreach ($scriptName in @('VERIFY_LIBRARY.ps1', 'CREATE_PORTABLE_ARCHIVE.ps1')) {
  & {
    param($scriptName, $sourceRoot, $testRoot)
    $errors = $null
    $tokens = $null
    $ast = [Management.Automation.Language.Parser]::ParseFile((Join-Path $sourceRoot $scriptName), [ref]$tokens, [ref]$errors)
    if ($errors.Count) { throw "Parse error in $scriptName" }
    # Load only declarations from the scripts under test; never run their main body.
    foreach ($name in @('releaseDirectories','releaseRootMarkdown','releaseRootScripts','releaseRootMetadata','releaseMarkdownSystemScripts','distributionManifestName')) {
      $nodes = @($ast.FindAll({ param($node) $node -is [Management.Automation.Language.AssignmentStatementAst] -and $node.Left.Extent.Text -ceq ('$' + $name) }, $false))
      if ($nodes.Count -gt 1) { throw "Ambiguous declaration: $name" }
      if ($nodes.Count -eq 1) { . ([scriptblock]::Create($nodes[0].Extent.Text)) }
    }
    foreach ($name in @('Assert-True','Relative','Relative-ToLibrary','Test-LocalMaintenanceEntry','Get-ReleaseInputPolicy','Get-ReleaseSourceFiles')) {
      $nodes = @($ast.FindAll({ param($node) $node -is [Management.Automation.Language.FunctionDefinitionAst] -and $node.Name -ceq $name }, $false))
      if ($nodes.Count -gt 1) { throw "Ambiguous function: $name" }
      if ($nodes.Count -eq 1) { . ([scriptblock]::Create($nodes[0].Extent.Text)) }
    }
    foreach ($case in @('absent','pair','state-only','events-only','directory','near-name','reparse','intake','intake-no-pair','intake-directory','intake-near-name','spec','spec-unknown','spec-reparse','spec-root-file','spec-near-name')) {
      $libraryRoot = Join-Path $testRoot ($scriptName + '-' + $case)
      [void](New-Item -ItemType Directory -Path $libraryRoot)
      [IO.File]::WriteAllText((Join-Path $libraryRoot 'README.md'), "synthetic selector fixture`n")
      if ($case -in @('pair','state-only','directory')) {
        [IO.File]::WriteAllText((Join-Path $libraryRoot 'PROJECT_EXECUTION_STATE.json'), '{}')
      }
      if ($case -in @('pair','events-only')) {
        [IO.File]::WriteAllText((Join-Path $libraryRoot 'PROJECT_EXECUTION_EVENTS.jsonl'), '{}')
      }
      if ($case -eq 'directory') { [void](New-Item -ItemType Directory -Path (Join-Path $libraryRoot 'PROJECT_EXECUTION_EVENTS.jsonl')) }
      if ($case -eq 'near-name') { [IO.File]::WriteAllText((Join-Path $libraryRoot 'PROJECT_EXECUTION_STATE.json.backup'), '{}') }
      if ($case.StartsWith('intake') -or $case.StartsWith('spec')) {
        if ($case -ne 'intake-no-pair') {
          [IO.File]::WriteAllText((Join-Path $libraryRoot 'PROJECT_EXECUTION_STATE.json'), '{}')
          [IO.File]::WriteAllText((Join-Path $libraryRoot 'PROJECT_EXECUTION_EVENTS.jsonl'), '{}')
        }
        switch ($case) {
          'intake' {
            foreach ($name in @('PROJECT_READINESS_GATE.json','PROJECT_LIBRARY_READINESS_GATE.json','PROJECT_READINESS_RECORD.md','PROJECT_ENGINEERING_CONTRACT.json','PROJECT_ADVISORY_A.md','PROJECT_ADVISORY_H.md')) {
              [IO.File]::WriteAllText((Join-Path $libraryRoot $name), 'synthetic, never distribute')
            }
          }
          'intake-no-pair' { [IO.File]::WriteAllText((Join-Path $libraryRoot 'PROJECT_READINESS_GATE.json'), '{}') }
          'intake-directory' { [void](New-Item -ItemType Directory -Path (Join-Path $libraryRoot 'PROJECT_READINESS_GATE.json')) }
          'intake-near-name' { [IO.File]::WriteAllText((Join-Path $libraryRoot 'PROJECT_ADVISORY_I.md'), 'not allowed') }
          'spec' {
            foreach ($path in @('.specify/memory/constitution.md','specs/library-maintenance/spec.md','specs/library-maintenance/plan.md','specs/library-maintenance/tasks.md')) {
              $target = Join-Path $libraryRoot $path
              [void](New-Item -ItemType Directory -Path (Split-Path $target -Parent) -Force)
              [IO.File]::WriteAllText($target, 'synthetic, never distribute')
            }
          }
          'spec-unknown' {
            $target = Join-Path $libraryRoot 'specs/library-maintenance'
            [void](New-Item -ItemType Directory -Path $target -Force)
            [IO.File]::WriteAllText((Join-Path $target 'payload.ps1'), 'unexpected')
          }
          'spec-reparse' {
            [void](New-Item -ItemType Directory -Path (Join-Path $libraryRoot 'specs'))
            $outside = Join-Path $testRoot ($scriptName + '-intake-outside')
            [void](New-Item -ItemType Directory -Path $outside)
            $linkType = if ($IsWindows) { 'Junction' } else { 'SymbolicLink' }
            [void](New-Item -ItemType $linkType -Path (Join-Path $libraryRoot 'specs/library-maintenance') -Target $outside)
          }
          'spec-root-file' { [IO.File]::WriteAllText((Join-Path $libraryRoot 'specs'), 'not a directory') }
          'spec-near-name' { [void](New-Item -ItemType Directory -Path (Join-Path $libraryRoot 'specs.backup')) }
        }
      }
      if ($case -eq 'reparse') {
        [IO.File]::WriteAllText((Join-Path $libraryRoot 'PROJECT_EXECUTION_STATE.json'), '{}')
        $junctionTarget = Join-Path $testRoot ($scriptName + '-junction-target')
        [void](New-Item -ItemType Directory -Path $junctionTarget)
        $linkType = if ($IsWindows) { 'Junction' } else { 'SymbolicLink' }
        [void](New-Item -ItemType $linkType -Path (Join-Path $libraryRoot 'PROJECT_EXECUTION_EVENTS.jsonl') -Target $junctionTarget)
      }
      $rejection = $null
      $selected = @()
      try { $selected = @(Get-ReleaseSourceFiles) } catch { $rejection = $_.Exception.Message }
      if ($case -in @('absent','pair','intake','spec')) {
        if ($rejection -or $selected.Count -ne 1 -or $selected[0].Name -cne 'README.md') { throw "$scriptName/$case did not select only README: $rejection" }
      } else {
        $expected = switch ($case) {
          'state-only' { 'must exist together' }
          'events-only' { 'must exist together' }
          'directory' { 'must be a regular file' }
          'near-name' { 'unknown top-level file' }
          'reparse' { 'symlink/reparse point' }
          'intake-no-pair' { 'requires the execution checkpoint pair' }
          'intake-directory' { 'must be a regular file' }
          'intake-near-name' { 'unknown top-level file' }
          'spec-unknown' { 'unknown local maintenance file' }
          'spec-reparse' { 'symlink/reparse point' }
          'spec-root-file' { 'must be a directory' }
          'spec-near-name' { 'unknown top-level directory' }
        }
        if (-not $rejection -or $rejection -inotmatch [regex]::Escape($expected)) { throw "$scriptName/$case did not reject for the expected reason: $rejection" }
      }
    }
  } $scriptName $sourceRoot $testRoot
  $checks += 16
}
# Load the canonical publisher in isolation; these .bin fixtures are not releases.
$publicationAst = [Management.Automation.Language.Parser]::ParseFile((Join-Path $sourceRoot 'CREATE_PORTABLE_ARCHIVE.ps1'), [ref]$null, [ref]$null)
$publisherNodes = @($publicationAst.FindAll({ param($node) $node -is [Management.Automation.Language.FunctionDefinitionAst] -and $node.Name -ceq 'Publish-VerifiedArchive' }, $false))
if ($publisherNodes.Count -ne 1) { throw 'Expected one canonical publisher' }
$publisherSource = $publisherNodes[0].Extent.Text
. ([scriptblock]::Create($publisherSource))
foreach ($case in @('success','archive-collision','checksum-collision','archive-directory','checksum-directory','missing-candidate','copy-failure','checksum-write-failure')) {
  $fixture = Join-Path $testRoot ('publication-' + $case)
  [void](New-Item -ItemType Directory -Path $fixture)
  $candidate = Join-Path $fixture 'candidate.bin'
  $destination = Join-Path $fixture 'output.bin'
  $checksum = $destination + '.sha256'
  $bytes = [Text.Encoding]::UTF8.GetBytes('synthetic publication payload')
  if ($case -ne 'missing-candidate') { [IO.File]::WriteAllBytes($candidate, $bytes) }
  if ($case -eq 'archive-collision') { [IO.File]::WriteAllText($destination, 'KEEP_ARCHIVE') }
  if ($case -eq 'checksum-collision') { [IO.File]::WriteAllText($checksum, 'KEEP_CHECKSUM') }
  if ($case -eq 'archive-directory') { [void](New-Item -ItemType Directory -Path $destination) }
  if ($case -eq 'checksum-directory') { [void](New-Item -ItemType Directory -Path $checksum) }
  $testedSource = $publisherSource
  if ($case -in @('copy-failure','checksum-write-failure')) {
    # Inject exactly one throw at an actual write boundary; all reservation and
    # rollback code is verbatim from the canonical publisher.
    $anchor = if ($case -eq 'copy-failure') { '$inputStream.CopyTo($archiveStream)' } else { '$checksumStream.Write($checksumBytes, 0, $checksumBytes.Length)' }
    if (($testedSource.Split($anchor).Count - 1) -ne 1) { throw "Ambiguous fault boundary $case" }
    $testedSource = $testedSource.Replace($anchor, ('throw "synthetic-' + $case + '"'))
  }
  . ([scriptblock]::Create($testedSource))
  $failure = $null
  try { $publishedHash = Publish-VerifiedArchive $candidate $destination } catch { $failure = $_.Exception.Message }
  if ($case -eq 'success') {
    $actual = (Get-FileHash -LiteralPath $destination -Algorithm SHA256).Hash.ToLowerInvariant()
    if ($failure -or $publishedHash -cne $actual -or [IO.File]::ReadAllText($checksum) -cne "$actual  output.bin`n") { throw 'Published pair is not exact' }
    if ([Convert]::ToHexString([IO.File]::ReadAllBytes($destination)) -cne [Convert]::ToHexString($bytes)) { throw 'Published bytes changed' }
  } else {
    if (-not $failure) { throw "Expected publication rejection: $case" }
    if ($case -eq 'archive-collision') {
      if ([IO.File]::ReadAllText($destination) -cne 'KEEP_ARCHIVE' -or (Test-Path -LiteralPath $checksum)) { throw 'Archive collision damaged originals or leaked reservation' }
    } elseif ($case -eq 'checksum-collision') {
      if ([IO.File]::ReadAllText($checksum) -cne 'KEEP_CHECKSUM' -or (Test-Path -LiteralPath $destination)) { throw 'Checksum collision damaged originals' }
    } elseif ($case -eq 'archive-directory') {
      if (-not [IO.Directory]::Exists($destination) -or (Test-Path -LiteralPath $checksum)) { throw 'Archive directory collision cleanup failed' }
    } elseif ($case -eq 'checksum-directory') {
      if (-not [IO.Directory]::Exists($checksum) -or (Test-Path -LiteralPath $destination)) { throw 'Checksum directory collision damaged originals' }
    } else {
      if ((Test-Path -LiteralPath $destination) -or (Test-Path -LiteralPath $checksum)) { throw "Owned partial outputs remain: $case" }
      if ($case -ne 'missing-candidate' -and $failure -notlike "*synthetic-$case*") { throw "Unexpected fault: $failure" }
    }
  }
  if ($case -ne 'missing-candidate' -and [Convert]::ToHexString([IO.File]::ReadAllBytes($candidate)) -cne [Convert]::ToHexString($bytes)) { throw "Candidate changed: $case" }
  $checks++
}

foreach ($round in 1..3) {
  $fixture = Join-Path $testRoot "publication-race-$round"
  [void](New-Item -ItemType Directory -Path $fixture)
  $destination = Join-Path $fixture 'output.bin'
  $gate = [Threading.ManualResetEvent]::new($false)
  $workers = @()
  try {
    foreach ($index in 0..1) {
      $candidate = Join-Path $fixture "candidate-$index.bin"
      [IO.File]::WriteAllText($candidate, ('synthetic-race-' + $index) * 4096)
      $worker = [PowerShell]::Create()
      [void]$worker.AddScript({ param($definition,$source,$target,$gate)
        $ErrorActionPreference='Stop'
        . ([scriptblock]::Create($definition))
        [void]$gate.WaitOne()
        try { $hash=Publish-VerifiedArchive $source $target; [pscustomobject]@{ok=$true;hash=$hash} }
        catch { [pscustomobject]@{ok=$false;error=$_.Exception.Message} }
      }).AddArgument($publisherSource).AddArgument($candidate).AddArgument($destination).AddArgument($gate)
      $workers += @{worker=$worker;handle=$worker.BeginInvoke()}
    }
    [void]$gate.Set()
    $results = @($workers | ForEach-Object { $_.worker.EndInvoke($_.handle) })
    if (@($results | Where-Object ok).Count -ne 1 -or @($results | Where-Object { -not $_.ok }).Count -ne 1) { throw 'Race did not yield exactly one publisher' }
    $hash = (Get-FileHash -LiteralPath $destination -Algorithm SHA256).Hash.ToLowerInvariant()
    if ($hash -cne @($results | Where-Object ok)[0].hash -or [IO.File]::ReadAllText($destination+'.sha256') -cne "$hash  output.bin`n") { throw 'Race published mixed archive/checksum' }
    $checks++
  } finally {
    [void]$gate.Set()
    foreach ($entry in $workers) { $entry.worker.Dispose() }
    $gate.Dispose()
  }
}
$zipNodes = @($publicationAst.FindAll({ param($node) $node -is [Management.Automation.Language.FunctionDefinitionAst] -and $node.Name -ceq 'Write-DeterministicArchive' }, $false))
if ($zipNodes.Count -ne 1) { throw 'Expected one canonical deterministic writer' }
. ([scriptblock]::Create($zipNodes[0].Extent.Text))
$zipHashes = @()
foreach ($year in 2000,2002) {
  $fixture = Join-Path $testRoot "zip-$year"
  [void](New-Item -ItemType Directory -Path (Join-Path $fixture 'nested'))
  foreach ($name in @('z.txt','nested/a.txt','A.txt')) {
    $file = Join-Path $fixture $name
    [IO.File]::WriteAllText($file, "same payload for $name")
    [IO.File]::SetLastWriteTimeUtc($file, [datetime]::new($year,1,1,0,0,0,[DateTimeKind]::Utc))
  }
  $candidate = Join-Path $testRoot "zip-$year.bin"
  Write-DeterministicArchive $fixture $candidate 946684800
  $zipHashes += (Get-FileHash -LiteralPath $candidate -Algorithm SHA256).Hash
  $zip = [IO.Compression.ZipFile]::OpenRead($candidate)
  try {
    $names = @($zip.Entries | ForEach-Object FullName)
    if (($names -join '|') -cne 'Elite Engineering Library/A.txt|Elite Engineering Library/nested/a.txt|Elite Engineering Library/z.txt') { throw 'ZIP ordering/path changed' }
    foreach ($entry in $zip.Entries) {
      if ($entry.LastWriteTime.ToString('yyyy-MM-ddTHH:mm:ss') -cne '2000-01-01T00:00:00' -or $entry.CompressedLength -ne $entry.Length -or $entry.ExternalAttributes -ne 0) { throw 'ZIP metadata/compression not deterministic' }
      $reader=[IO.StreamReader]::new($entry.Open())
      try { $actual=$reader.ReadToEnd() } finally { $reader.Dispose() }
      if ($actual -cne ('same payload for ' + $entry.FullName.Substring('Elite Engineering Library/'.Length))) { throw 'ZIP payload changed' }
    }
  } finally { $zip.Dispose() }
}
if ($zipHashes[0] -cne $zipHashes[1]) { throw 'Equal payload produced different ZIP bytes' }
$checks++
foreach ($case in @('odd-epoch','early-epoch','late-epoch','existing-target','empty','reparse')) {
  $fixture = Join-Path $testRoot "zip-reject-$case"
  [void](New-Item -ItemType Directory -Path $fixture)
  $candidate = Join-Path $testRoot "zip-reject-$case.bin"
  if ($case -ne 'empty') { [IO.File]::WriteAllText((Join-Path $fixture 'README.md'),'synthetic') }
  $epoch=946684800
  if ($case -eq 'odd-epoch') { $epoch++ }
  if ($case -eq 'early-epoch') { $epoch=315532798 }
  if ($case -eq 'late-epoch') { $epoch=4354819200 }
  if ($case -eq 'existing-target') { [IO.File]::WriteAllText($candidate,'KEEP') }
  if ($case -eq 'reparse') {
    $linkType=if($IsWindows){'Junction'}else{'SymbolicLink'}
    [void](New-Item -ItemType $linkType -Path (Join-Path $fixture 'outside') -Target (Join-Path $testRoot 'zip-2000'))
  }
  $failure=$null;try{Write-DeterministicArchive $fixture $candidate $epoch}catch{$failure=$_.Exception.Message}
  if(-not $failure){throw "Expected ZIP rejection: $case"}
  if($case-eq'existing-target') { if([IO.File]::ReadAllText($candidate)-cne'KEEP'){throw 'ZIP target overwritten'} }
  elseif(Test-Path -LiteralPath $candidate){throw "Rejected ZIP left output: $case"}
  $checks++
}
Write-Output "PASS: local maintenance selectors, publication and reproducibility ($checks checks); fixture=$testRoot"
