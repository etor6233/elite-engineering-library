#requires -Version 7.0

[CmdletBinding()]
param()

$ErrorActionPreference = 'Stop'
$libraryRoot = [IO.Path]::GetFullPath($PSScriptRoot).TrimEnd([IO.Path]::DirectorySeparatorChar)
$materializer = Join-Path $libraryRoot 'materialize_markdown_pack.ps1'
$tempRoot = Join-Path ([IO.Path]::GetTempPath()) ('elite-library-verification-' + [guid]::NewGuid().ToString('N'))
$utf8 = [Text.UTF8Encoding]::new($false)
$strictUtf8 = [Text.UTF8Encoding]::new($false, $true)
$distributionManifestName = 'DISTRIBUTION_SHA256SUMS.txt'
$releaseDirectories = @{
  'architecture_packs'      = @('.md')
  'implementation_packs'    = @('.md')
  'markdown_system'         = @('.md', '.ps1')
  'reconstruction_evidence' = @('.md')
}
$releaseRootMarkdown = @(
  'AGENT_SYSTEM_START.md',
  'AGENTIC_AI_SOFTWARE_ENGINEERING_CODEX.md',
  'AGENTS.md',
  'AI_ENGINEERING_MASTER_MAP.md',
  'AI_INFERENCE_PERFORMANCE_HARDWARE.md',
  'AI_SECURITY_GOVERNANCE_PRIVACY.md',
  'ALGORITHMS_DATA_STRUCTURES_PROBLEM_SOLVING.md',
  'CLAUDE.md',
  'CODEX_ELITE_PROJECT_BOOTSTRAP.md',
  'COMPUTER_SYSTEMS_PERFORMANCE_LOW_LATENCY.md',
  'DATA_ENGINEERING_ANALYTICS.md',
  'DATABASE_STORAGE_INTERNALS.md',
  'DEEP_LEARNING_ANDREW_NG.md',
  'ENGINEERING_EXECUTION_PLAYBOOK.md',
  'ENTERPRISE_FULL_STACK_BLUEPRINT.md',
  'FRANCHISE_ACCELERATOR.md',
  'FRONTEND_PRODUCT_ENGINEERING_UX.md',
  'GPU_ACCELERATED_COMPUTING.md',
  'LEADING_COMPANY_PUBLIC_CODE_MATRIX.md',
  'LICENSE.md',
  'MACHINE_LEARNING_SPECIALIZATION_AND_MATH_FOUNDATIONS.md',
  'MARKET_MICROSTRUCTURE_EXCHANGE_SYSTEMS.md',
  'ML_PRODUCTION_LLMOPS_EVALUATION.md',
  'NATIVE_MOBILE_DESKTOP_ENGINEERING.md',
  'NETWORKING_DISTRIBUTED_STREAMING.md',
  'NLP_RAG_RETRIEVAL_DATA.md',
  'PUBLIC_CODE_ARCHITECTURE_ADMISSION_STANDARD.md',
  'PUBLIC_CODE_ARCHITECTURE_MASTER_MAP.md',
  'PROJECT_FAILURE_LESSONS.md',
  'PYTORCH_TRANSFORMERS_GENERATIVE_MODELS.md',
  'README.md',
  'REUSABLE_CODE_READINESS_ROADMAP.md',
  'SECURITY_SRE_CLOUD_INFRASTRUCTURE.md',
  'SOFTWARE_ARCHITECTURE_SYSTEM_DESIGN.md',
  'SOFTWARE_BACKEND_API_ENGINEERING.md',
  'START_ANY_PROJECT.md',
  'START_FRANCHISE.md',
  'SYSTEMS_ENGINEERING_MASTER_MAP.md',
  'THIRD_PARTY_NOTICES.md',
  'TOOLCHAINS_BUILDS_PACKAGING_FFI.md'
)
$releaseRootScripts = @('materialize_markdown_pack.ps1', 'VERIFY_LIBRARY.ps1', 'VERIFY_EXECUTABLE_LIBRARY.ps1', 'CREATE_PORTABLE_ARCHIVE.ps1', 'INSTALL_AGENT_BRIDGE.ps1')
$releaseRootMetadata = @('.gitignore','.gitattributes')
$releaseMarkdownSystemScripts = @('update_pack_from_tree.ps1', 'test_update_pack_from_tree.ps1', 'test_pack_exact_eof.ps1', 'test_public_release_policy.ps1', 'test_install_agent_bridge.ps1', 'test_library_maintenance_state.ps1', 'test_agent_bridge_safety.ps1', 'test_toolchain_resolution.ps1', 'test_agent_entry_lifecycle.ps1', 'project_local_secrets.ps1', 'test_project_local_secrets.ps1')

function Assert-True([bool] $Condition, [string] $Message) {
  if (-not $Condition) { throw "VERIFY_LIBRARY_FAILED: $Message" }
}

function Assert-FranchiseProfileSelection([string[]] $PackIds) {
  Assert-True ($PackIds -contains 'MICROSOFT-DEVSKIM-ADAPTED-SAST-GATE') 'franchise profile omitted the admitted DevSkim baseline'
  Assert-True ($PackIds -contains 'PORTABLE-SIGNED-RELEASE-EVIDENCE-GATE') 'franchise profile omitted portable signed release evidence'
  Assert-True ($PackIds -contains 'GO-MERCADOLIBRE-MARKETPLACE-ADAPTER') 'franchise profile omitted Mercado Libre questions ingress'
  Assert-True ($PackIds -contains 'GO-MERCADOLIBRE-QUESTION-OUTBOUND') 'franchise profile omitted Mercado Libre question outbound'
  Assert-True ($PackIds -notcontains 'GITLAB-OPENGREP-SIGNED-SAST-GATE') 'franchise profile selected OpenGrep while its Cosign verifier is REJECTED_VERIFIER_RUNTIME'
  Assert-True ($PackIds -contains 'GO-CUSTOMER-SURVEY-API') 'franchise profile omitted customer survey API'
  Assert-True ($PackIds -contains 'TS-CUSTOMER-SURVEY-PORTAL') 'franchise profile omitted customer survey portal'
  Assert-True ($PackIds -contains 'GO-BC-EXACT-AMOUNT-ADAPTER') 'franchise profile omitted GO-BC-EXACT-AMOUNT-ADAPTER'
  Assert-True ($PackIds -contains 'GO-BC-SALES-CONTRACT-ADAPTER') 'franchise profile omitted GO-BC-SALES-CONTRACT-ADAPTER'
  Assert-True ($PackIds -contains 'GO-OFFICIAL-PAYMENT-WEBHOOK-ADAPTERS') 'franchise profile omitted GO-OFFICIAL-PAYMENT-WEBHOOK-ADAPTERS'
  Assert-True ($PackIds -contains 'GO-PAYMENT-CHECKOUT-RUNTIME') 'franchise profile omitted GO-PAYMENT-CHECKOUT-RUNTIME'
  Assert-True ($PackIds -contains 'GO-INITIAL-HANDOVER-API') 'franchise profile omitted GO-INITIAL-HANDOVER-API'
  Assert-True ($PackIds -contains 'GO-BUSINESS-POLICY-PROFILE') 'franchise profile omitted GO-BUSINESS-POLICY-PROFILE'
  Assert-True ($PackIds -contains 'TS-PAYMENT-CHECKOUT-PORTAL') 'franchise profile omitted TS-PAYMENT-CHECKOUT-PORTAL'
  $policy = Get-ReleaseInputPolicy
  Assert-True ($null -ne $policy) 'franchise release selection lock missing'
  Assert-True ($PackIds.Count -eq $policy.franchise_pack_ids.Count) "franchise profile pack count drifted: expected=$($policy.franchise_pack_ids.Count) actual=$($PackIds.Count)"
  Assert-True (($PackIds | Sort-Object -Unique).Count -eq $PackIds.Count) 'franchise profile duplicate ID'
  foreach ($id in $policy.franchise_pack_ids) { Assert-True ($PackIds -ccontains $id) "franchise profile omitted locked source: $id" }
}

function Assert-FranchiseProfileSelectionRegression {
  $valid = @((Get-ReleaseInputPolicy).franchise_pack_ids)
  Assert-FranchiseProfileSelection $valid
  foreach ($negative in @(
    [pscustomobject]@{ Name = 'unknown-replacement'; Ids = @($valid | Where-Object { $_ -ne 'GO-CONNECTED-DOCUMENT-REFERENCE' }) + 'SYNTHETIC-UNKNOWN'; Expected = 'omitted locked source' },
    [pscustomobject]@{ Name = 'duplicate'; Ids = @($valid | Where-Object { $_ -ne 'GO-CONNECTED-DOCUMENT-REFERENCE' }) + 'GO-CUSTOMER-SURVEY-API'; Expected = 'duplicate ID' },
    [pscustomobject]@{ Name = 'missing-count'; Ids = @($valid | Where-Object { $_ -ne 'GO-CONNECTED-DOCUMENT-REFERENCE' }); Expected = 'pack count drifted' },
    [pscustomobject]@{ Name = 'missing-survey-api'; Ids = @($valid | Where-Object { $_ -ne 'GO-CUSTOMER-SURVEY-API' }) + 'SYNTHETIC-SURVEY-A'; Expected = 'omitted customer survey API' },
    [pscustomobject]@{ Name = 'missing-survey-portal'; Ids = @($valid | Where-Object { $_ -ne 'TS-CUSTOMER-SURVEY-PORTAL' }) + 'SYNTHETIC-SURVEY-B'; Expected = 'omitted customer survey portal' },
    [pscustomobject]@{ Name = 'rejected'; Ids = @($valid + 'GITLAB-OPENGREP-SIGNED-SAST-GATE'); Expected = 'selected OpenGrep' },
    [pscustomobject]@{ Name = 'missing-devskim'; Ids = @($valid | Where-Object { $_ -ne 'MICROSOFT-DEVSKIM-ADAPTED-SAST-GATE' }) + 'SYNTHETIC-REPLACEMENT-A'; Expected = 'omitted the admitted DevSkim' },
    [pscustomobject]@{ Name = 'missing-release'; Ids = @($valid | Where-Object { $_ -ne 'PORTABLE-SIGNED-RELEASE-EVIDENCE-GATE' }) + 'SYNTHETIC-REPLACEMENT-B'; Expected = 'omitted portable signed release' },
    [pscustomobject]@{ Name = 'missing-mercadolibre'; Ids = @($valid | Where-Object { $_ -ne 'GO-MERCADOLIBRE-MARKETPLACE-ADAPTER' }) + 'SYNTHETIC-REPLACEMENT-C'; Expected = 'omitted Mercado Libre questions ingress' },
    [pscustomobject]@{ Name = 'missing-mercadolibre-outbound'; Ids = @($valid | Where-Object { $_ -ne 'GO-MERCADOLIBRE-QUESTION-OUTBOUND' }) + 'SYNTHETIC-REPLACEMENT-D'; Expected = 'omitted Mercado Libre question outbound' }
  )) {
    $rejected = $false
    try { Assert-FranchiseProfileSelection $negative.Ids } catch {
      if ($_.Exception.Message.Contains($negative.Expected, [StringComparison]::Ordinal)) { $rejected = $true } else { throw }
    }
    Assert-True $rejected "franchise profile negative regression was accepted: $($negative.Name)"
  }
}

function Relative([string] $Path) {
  [IO.Path]::GetRelativePath($libraryRoot, $Path).Replace('\', '/')
}

function Assert-LocalGoImportsResolved([string] $Destination, [string] $Profile) {
  $goModPath = Join-Path $Destination 'go.mod'
  if (-not (Test-Path -LiteralPath $goModPath -PathType Leaf)) { return }
  $goMod = [IO.File]::ReadAllText($goModPath, $strictUtf8)
  $moduleMatch = [regex]::Match($goMod, '(?m)^module\s+(?<module>\S+)\s*$')
  Assert-True $moduleMatch.Success "Go module declaration missing in $Profile"
  $module = $moduleMatch.Groups['module'].Value
  $prefix = $module + '/'
  $imports = [Collections.Generic.HashSet[string]]::new([StringComparer]::Ordinal)
  foreach ($goFile in Get-ChildItem -LiteralPath $Destination -Filter '*.go' -File -Recurse) {
    $source = [IO.File]::ReadAllText($goFile.FullName, $strictUtf8)
    foreach ($importMatch in [regex]::Matches($source, '"(?<path>[^"\r\n]+)"')) {
      $importPath = $importMatch.Groups['path'].Value
      if ($importPath.StartsWith($prefix, [StringComparison]::Ordinal)) { [void]$imports.Add($importPath) }
    }
  }
  foreach ($importPath in $imports) {
    $relativePackage = $importPath.Substring($prefix.Length)
    Assert-True ($relativePackage -ne '' -and -not $relativePackage.Contains('..')) "unsafe local Go import in ${Profile}: $importPath"
    $packageDirectory = Join-Path $Destination $relativePackage.Replace('/', [IO.Path]::DirectorySeparatorChar)
    Assert-True (Test-Path -LiteralPath $packageDirectory -PathType Container) "unresolved local Go import in ${Profile}: $importPath"
    $implementationFiles = @(Get-ChildItem -LiteralPath $packageDirectory -Filter '*.go' -File | Where-Object { $_.Name -notlike '*_test.go' })
    Assert-True ($implementationFiles.Count -gt 0) "local Go import has no implementation files in ${Profile}: $importPath"
  }
}

function Read-MaterializationProvenance([string] $Path) {
  $values = [Collections.Generic.List[string]]::new()
  $lines = [IO.File]::ReadAllLines($Path, $strictUtf8)
  $insidePayload = $false
  for ($index = 0; $index -lt $lines.Length; $index++) {
    if ($insidePayload) {
      if ($lines[$index] -ceq '````') { $insidePayload = $false }
      continue
    }
    if ($lines[$index] -notmatch '^### FILE: `[^`]+`$') { continue }
    $cursor = $index + 1
    while ($cursor -lt $lines.Length -and $lines[$cursor] -ceq '') { $cursor++ }
    Assert-True ($cursor -lt $lines.Length -and $lines[$cursor] -ceq '```yaml') "FILE metadata fence missing in $(Relative $Path)"
    $cursor++
    $metadata = [Collections.Generic.List[string]]::new()
    while ($cursor -lt $lines.Length -and $lines[$cursor] -cne '```') {
      $metadata.Add($lines[$cursor])
      $cursor++
    }
    Assert-True ($cursor -lt $lines.Length) "FILE metadata fence is unterminated in $(Relative $Path)"
    $matches = @($metadata | Where-Object { $_ -match '^provenance:\s*(AUTHORED|ADAPTED|VERBATIM)\s*$' })
    Assert-True ($matches.Count -eq 1) "FILE provenance is missing, duplicated or invalid in $(Relative $Path)"
    [void]($matches[0] -match '^provenance:\s*(AUTHORED|ADAPTED|VERBATIM)\s*$')
    $values.Add($Matches[1])
    $cursor++
    while ($cursor -lt $lines.Length -and $lines[$cursor] -ceq '') { $cursor++ }
    Assert-True ($cursor -lt $lines.Length -and $lines[$cursor] -match '^````[^`]*$') "FILE payload fence missing in $(Relative $Path)"
    $index = $cursor
    $insidePayload = $true
  }
  Assert-True (-not $insidePayload) "FILE payload fence is unterminated in $(Relative $Path)"
  return $values.ToArray()
}

function Test-LocalMaintenanceEntry([IO.FileSystemInfo] $Item, [int] $CheckpointCount) {
  # AUTHORED distribution boundary, not approval of any local record.
  $localFiles = @(
    'PROJECT_READINESS_RECORD.md', 'PROJECT_READINESS_GATE.json', 'PROJECT_READINESS_REPORT.json',
    'PROJECT_LIBRARY_READINESS_GATE.json',
    'PROJECT_DEPENDENCY_UPDATE_RECORD.md', 'PROJECT_AUTHORITY_FRESHNESS_RECORD.md',
    'PROJECT_VULNERABILITY_MONITORING_RECORD.md', 'PROJECT_OFFICIAL_SOURCE_PROFILE_RECORD.md',
    'PROJECT_BLUEPRINT.md', 'PROJECT_AUTHORITY_MAP.md', 'PROJECT_EXTERNAL_SOURCE_LOCK.md',
    'PROJECT_PACK_PLAN.md', 'PROJECT_ENGINEERING_CONTRACT.json',
    'PROJECT_ADVISORY_A.md', 'PROJECT_ADVISORY_B.md', 'PROJECT_ADVISORY_C.md', 'PROJECT_ADVISORY_D.md',
    'PROJECT_ADVISORY_E.md', 'PROJECT_ADVISORY_F.md', 'PROJECT_ADVISORY_G.md', 'PROJECT_ADVISORY_H.md'
  )
  $localDirectories = @('.specify', '.specify/memory', 'specs', 'specs/library-maintenance')
  $localSpecFiles = @('.specify/memory/constitution.md', 'specs/library-maintenance/spec.md',
    'specs/library-maintenance/plan.md', 'specs/library-maintenance/tasks.md')
  if ($localFiles -ccontains $Item.Name) {
    if ($Item.PSIsContainer) { throw "Local maintenance artifact must be a regular file: $($Item.Name)" }
    if ($CheckpointCount -ne 2) { throw 'Local maintenance requires the execution checkpoint pair' }
    return $true
  }
  if ($Item.Name -cnotin @('.specify', 'specs')) { return $false }
  if (-not $Item.PSIsContainer) { throw "Local maintenance root must be a directory: $($Item.Name)" }
  if ($CheckpointCount -ne 2) { throw 'Local maintenance requires the execution checkpoint pair' }
  # Inspect each child before descending; never follow a junction to external data.
  $pending = [Collections.Generic.Stack[IO.DirectoryInfo]]::new()
  $pending.Push($Item)
  while ($pending.Count -gt 0) {
    foreach ($child in Get-ChildItem -LiteralPath $pending.Pop().FullName -Force) {
      if (($child.Attributes -band [IO.FileAttributes]::ReparsePoint) -ne 0) {
        throw "Local maintenance entry is a symlink/reparse point: $($child.Name)"
      }
      $relative = [IO.Path]::GetRelativePath($libraryRoot, $child.FullName).Replace('\', '/')
      if ($child.PSIsContainer) {
        if ($localDirectories -cnotcontains $relative) { throw "Unknown local maintenance directory: $relative" }
        $pending.Push($child)
      } elseif ($localSpecFiles -cnotcontains $relative) {
        throw "Unknown local maintenance file: $relative"
      }
    }
  }
  return $true
}

# AUTHORED strict selection metadata reader, shared byte-for-byte by both entrypoints.
function Find-MachineLocalHomePath([string] $Text) {
  [regex]::Match($Text, '(?:(?i:[A-Z]:\\Users\\[^\\\s`"''<>]+)|(?<![A-Za-z0-9:])/(?:home|Users)/[^/\s`"''<>]+)')
}

function Get-ReleaseInputPolicy {
  $parent = Join-Path $libraryRoot 'markdown_system'
  $path = Join-Path $parent 'PUBLIC_RELEASE_INPUT_POLICY.md'
  if (-not (Test-Path -LiteralPath $path)) { return $null }
  foreach ($part in @($parent,$path)) {
    $item = Get-Item -LiteralPath $part -Force
    if (($item.Attributes -band [IO.FileAttributes]::ReparsePoint) -ne 0) { throw 'Release policy is a symlink/reparse point' }
  }
  $item = Get-Item -LiteralPath $path -Force
  if ($item.PSIsContainer -or $item.Length -gt 65536) { throw 'Release policy must be a bounded regular file' }
  $raw = [IO.File]::ReadAllText($path,[Text.UTF8Encoding]::new($false,$true))
  $matches = [regex]::Matches($raw,'(?s)```json\s*(?<json>.*?)\s*```')
  if ($matches.Count -ne 1) { throw 'Release policy requires exactly one JSON block' }
  $doc = [Text.Json.JsonDocument]::Parse([string]$matches[0].Groups['json'].Value)
  try {
    $root = $doc.RootElement
    if ($root.ValueKind -ne [Text.Json.JsonValueKind]::Object) { throw 'Release policy must be an object' }
    $keys = [Collections.Generic.HashSet[string]]::new([StringComparer]::Ordinal)
    foreach ($prop in $root.EnumerateObject()) {
      if (-not $keys.Add($prop.Name)) { throw 'Release policy duplicate field' }
    }
    $expected = @('schema','json_files','text_logs','franchise_pack_ids','franchise_files')
    if ($keys.Count -ne $expected.Count -or @($expected | Where-Object { -not $keys.Contains($_) }).Count) { throw 'Release policy unknown/missing field' }
    if ($root.GetProperty('schema').GetString() -cne 'elite-public-release-input-policy/v1') { throw 'Release policy schema mismatch' }
    $result = @{ schema = 'elite-public-release-input-policy/v1' }
    foreach ($field in @('json_files','text_logs','franchise_pack_ids')) {
      $array = $root.GetProperty($field)
      if ($array.ValueKind -ne [Text.Json.JsonValueKind]::Array) { throw 'Release policy array required' }
      $list = [Collections.Generic.List[string]]::new()
      $seen = [Collections.Generic.HashSet[string]]::new([StringComparer]::Ordinal)
      foreach ($entry in $array.EnumerateArray()) {
        if ($entry.ValueKind -ne [Text.Json.JsonValueKind]::String) { throw 'Release policy string item required' }
        $name = $entry.GetString()
        $pattern = switch ($field) { 'json_files' { '^[A-Z][A-Z0-9_]*\.json$' } 'text_logs' { '^[A-Z][A-Z0-9_]*\.log$' } 'franchise_pack_ids' { '^[A-Z][A-Z0-9-]+$' } }
        if ($name -cnotmatch $pattern -or -not $seen.Add($name)) { throw 'Release policy unsafe/duplicate item' }
        $list.Add($name)
      }
      $result[$field] = @($list.ToArray())
    }
    $count = $root.GetProperty('franchise_files').GetInt32()
    if ($count -lt 1 -or $result.franchise_pack_ids.Count -lt 1) { throw 'Release policy empty franchise lock' }
    $result.franchise_files = $count
    return $result
  } finally { $doc.Dispose() }
}

function Get-ReleaseSourceFiles {
  $files = [Collections.Generic.List[IO.FileInfo]]::new()
  $releasePolicy = Get-ReleaseInputPolicy
  # Local continuity is not readiness evidence for a consuming project.
  $localStateNames = @('PROJECT_EXECUTION_STATE.json', 'PROJECT_EXECUTION_EVENTS.jsonl')
  $localStateCount = @($localStateNames | Where-Object { Test-Path -LiteralPath (Join-Path $libraryRoot $_) }).Count
  Assert-True ($localStateCount -in @(0, 2)) 'local execution state and events must exist together'
  foreach ($item in Get-ChildItem -LiteralPath $libraryRoot -Force) {
    if ($item.Name -ieq '.git') { continue }
    Assert-True (($item.Attributes -band [IO.FileAttributes]::ReparsePoint) -eq 0) "release entry is a symlink/reparse point: $($item.Name)"

    if ($localStateNames -ccontains $item.Name) {
      Assert-True (-not $item.PSIsContainer) "local execution artifact must be a regular file: $($item.Name)"
      continue
    }
    if (Test-LocalMaintenanceEntry $item $localStateCount) { continue }

    if (-not $item.PSIsContainer -and $item.Name -ceq $distributionManifestName) {
      continue
    }

    if ($item.PSIsContainer) {
      Assert-True ($releaseDirectories.ContainsKey($item.Name)) "unknown top-level directory: $($item.Name)"
      $allowedExtensions = $releaseDirectories[$item.Name]
      foreach ($child in Get-ChildItem -LiteralPath $item.FullName -Recurse -Force) {
        Assert-True (($child.Attributes -band [IO.FileAttributes]::ReparsePoint) -eq 0) "release entry is a symlink/reparse point: $(Relative $child.FullName)"
        if ($child.PSIsContainer) { continue }
        $publicEvidence = $null -ne $releasePolicy -and $item.Name -ceq 'reconstruction_evidence' -and
          $child.DirectoryName -ceq $item.FullName -and (@($releasePolicy.json_files) + @($releasePolicy.text_logs)) -ccontains $child.Name
        Assert-True (($allowedExtensions -icontains $child.Extension) -or $publicEvidence) "unknown release file: $(Relative $child.FullName)"
        if ($item.Name -eq 'markdown_system' -and $child.Extension -ieq '.ps1') {
          Assert-True ($releaseMarkdownSystemScripts -icontains $child.Name) "unknown markdown_system script: $(Relative $child.FullName)"
        }
        $files.Add($child)
      }
      continue
    }

    $allowedRootFile = $releaseRootMarkdown -icontains $item.Name -or $releaseRootScripts -icontains $item.Name -or $releaseRootMetadata -ccontains $item.Name
    Assert-True $allowedRootFile "unknown top-level file: $($item.Name)"
    $files.Add($item)
  }

  Assert-True ($files.Count -gt 0) 'release file allowlist produced no files'
  if ($null -ne $releasePolicy) {
    foreach ($name in (@($releasePolicy.json_files) + @($releasePolicy.text_logs))) {
      if (@($files | Where-Object { $_.DirectoryName -ceq (Join-Path $libraryRoot 'reconstruction_evidence') -and $_.Name -ceq $name }).Count -ne 1) { throw "Release policy named evidence missing: $name" }
    }
  }
  $relativePaths = @($files | ForEach-Object { Relative $_.FullName })
  Assert-True (($relativePaths | Sort-Object -Unique).Count -eq $relativePaths.Count) 'duplicate release path'
  $files | Sort-Object FullName
}

function Assert-DistributionManifest([IO.FileInfo[]] $ReleaseFiles) {
  $manifestPath = Join-Path $libraryRoot $distributionManifestName
  if (-not (Test-Path -LiteralPath $manifestPath)) { return }

  $manifestItem = Get-Item -LiteralPath $manifestPath -Force
  Assert-True (-not $manifestItem.PSIsContainer) "$distributionManifestName must be a file"
  Assert-True (($manifestItem.Attributes -band [IO.FileAttributes]::ReparsePoint) -eq 0) "$distributionManifestName must not be a symlink/reparse point"
  $manifestBytes = [IO.File]::ReadAllBytes($manifestPath)
  try {
    $manifestText = $strictUtf8.GetString($manifestBytes)
  } catch {
    throw "VERIFY_LIBRARY_FAILED: $distributionManifestName is not valid UTF-8"
  }
  Assert-True (-not $manifestText.Contains([char]0)) "NUL byte in $distributionManifestName"
  Assert-True ($manifestText.EndsWith("`n", [StringComparison]::Ordinal)) "$distributionManifestName must end with a newline"

  $lines = @($manifestText -split '\r?\n')
  Assert-True ($lines.Count -gt 1) "$distributionManifestName is empty"
  $lines = @($lines[0..($lines.Count - 2)])
  Assert-True (@($lines | Where-Object { $_ -eq '' }).Count -eq 0) "$distributionManifestName contains a blank line"

  $expected = [Collections.Generic.Dictionary[string,object]]::new([StringComparer]::OrdinalIgnoreCase)
  foreach ($file in $ReleaseFiles) {
    $relative = Relative $file.FullName
    Assert-True (-not $expected.ContainsKey($relative)) "case-insensitive duplicate release path: $relative"
    $expected.Add($relative, [pscustomobject]@{
      Path = $relative
      Hash = (Get-FileHash -Algorithm SHA256 -LiteralPath $file.FullName).Hash.ToLowerInvariant()
    })
  }

  $seen = [Collections.Generic.HashSet[string]]::new([StringComparer]::OrdinalIgnoreCase)
  foreach ($line in $lines) {
    $match = [regex]::Match($line, '^(?<hash>[0-9a-f]{64})  (?<path>.+)$')
    Assert-True $match.Success "invalid $distributionManifestName line: $line"
    $relative = $match.Groups['path'].Value
    $segments = @($relative -split '/')
    $unsafePath = $relative.Contains('\') -or
      $relative.StartsWith('/', [StringComparison]::Ordinal) -or
      $relative.EndsWith('/', [StringComparison]::Ordinal) -or
      $relative -match '^[A-Za-z]:' -or
      @($segments | Where-Object { $_ -eq '' -or $_ -eq '.' -or $_ -eq '..' }).Count -gt 0
    Assert-True (-not $unsafePath) "unsafe $distributionManifestName path: $relative"
    Assert-True ($relative -cne $distributionManifestName) "$distributionManifestName must not hash itself"
    Assert-True ($seen.Add($relative)) "duplicate $distributionManifestName path: $relative"
    Assert-True ($expected.ContainsKey($relative)) "unexpected $distributionManifestName path: $relative"
    $expectedEntry = $expected[$relative]
    Assert-True ($relative -ceq $expectedEntry.Path) "path casing mismatch in ${distributionManifestName}: $relative"
    Assert-True ($match.Groups['hash'].Value -ceq $expectedEntry.Hash) "SHA-256 mismatch in ${distributionManifestName}: $relative"
  }

  Assert-True ($seen.Count -eq $expected.Count) "$distributionManifestName coverage mismatch: manifest=$($seen.Count), release=$($expected.Count)"
}

function Get-StructuralH2Headings([string] $Text) {
  $headings = [Collections.Generic.List[object]]::new()
  $insideFence = $false
  $fenceCharacter = ''
  $fenceLength = 0

  foreach ($lineMatch in [regex]::Matches($Text, '(?m)^(?<line>[^\r\n]*)(?:\r?\n|$)')) {
    $line = $lineMatch.Groups['line'].Value
    if ($insideFence) {
      $closingPattern = '^ {0,3}' + [regex]::Escape($fenceCharacter) + '{' + $fenceLength + ',}[ \t]*$'
      if ($line -match $closingPattern) {
        $insideFence = $false
        $fenceCharacter = ''
        $fenceLength = 0
      }
      continue
    }

    $opening = [regex]::Match($line, '^ {0,3}(?<fence>`{3,}|~{3,}).*$')
    if ($opening.Success) {
      $insideFence = $true
      $fenceCharacter = [string]$opening.Groups['fence'].Value[0]
      $fenceLength = $opening.Groups['fence'].Value.Length
      continue
    }

    if ($line.StartsWith('## ', [StringComparison]::Ordinal)) {
      $headings.Add([pscustomobject]@{
        Text = $line
        Index = $lineMatch.Index
        Length = $line.Length
      })
    }
  }

  $headings
}

function Assert-PackContract([string] $Path, [string] $Text) {
  $requiredSections = @(
    '## 1. Metadata',
    '## 2. Applicability',
    '## 3. Architecture contract',
    '## 4. Exact file manifest',
    '## 5. Materialization blocks',
    '## 6. Configuration surface',
    '## 7. Dependency bill',
    '## 8. Apply order',
    '## 9. Verification',
    '## 10. Reconstruction evidence'
  )
  $structuralHeadings = @(Get-StructuralH2Headings $Text)
  $matches = [Collections.Generic.List[object]]::new()
  $previous = -1
  foreach ($heading in $requiredSections) {
    $headingMatches = @($structuralHeadings | Where-Object { $_.Text -ceq $heading })
    Assert-True ($headingMatches.Count -eq 1) "required pack section must occur exactly once outside code fences in $(Relative $Path): $heading"
    $match = $headingMatches[0]
    Assert-True ($match.Index -gt $previous) "pack sections out of order in $(Relative $Path): $heading"
    $matches.Add($match)
    $previous = $match.Index
  }
  for ($index = 0; $index -lt $matches.Count; $index++) {
    $start = $matches[$index].Index + $matches[$index].Length
    $end = if ($index + 1 -lt $matches.Count) { $matches[$index + 1].Index } else { $Text.Length }
    Assert-True ($Text.Substring($start, $end - $start).Trim().Length -gt 0) "empty pack section in $(Relative $Path): $($requiredSections[$index])"
  }

  $metadataStart = $matches[0].Index + $matches[0].Length
  $metadataLength = $matches[1].Index - $metadataStart
  $metadataSection = $Text.Substring($metadataStart, $metadataLength)
  $metadataMatch = [regex]::Match($metadataSection, '(?ms)^\s*^```yaml\r?$\s*(?<body>.*?)^```\r?$')
  Assert-True $metadataMatch.Success "metadata YAML block missing in $(Relative $Path)"
  $metadata = $metadataMatch.Groups['body'].Value
  foreach ($field in @('pack_id','pack_version','status','claim','stacks','compatible_with','incompatible_with','license_expression','upstream_sources','verified_at')) {
    $fieldMatches = [regex]::Matches($metadata, '(?m)^' + [regex]::Escape($field) + ':.*$')
    Assert-True ($fieldMatches.Count -eq 1) "metadata field $field must occur exactly once in $(Relative $Path)"
  }
  foreach ($statusField in @('authority','implementation','admission')) {
    $statusMatches = [regex]::Matches($metadata, '(?m)^  ' + [regex]::Escape($statusField) + ':\s*(?<value>\S+)\s*$')
    Assert-True ($statusMatches.Count -eq 1) "status field $statusField must occur exactly once in $(Relative $Path)"
  }
  $authority = [regex]::Match($metadata, '(?m)^  authority:\s*(?<value>\S+)\s*$').Groups['value'].Value
  $implementation = [regex]::Match($metadata, '(?m)^  implementation:\s*(?<value>\S+)\s*$').Groups['value'].Value
  $admission = [regex]::Match($metadata, '(?m)^  admission:\s*(?<value>\S+)\s*$').Groups['value'].Value
  Assert-True (@('NOTE','SUPPORTED_REFERENCE','ELITE_REFERENCE') -contains $authority) "invalid authority status in $(Relative $Path): $authority"
  Assert-True (@('SPEC_ONLY','SNIPPET','RECONSTRUCTIBLE','REBUILD_VERIFIED') -contains $implementation) "invalid implementation status in $(Relative $Path): $implementation"
  Assert-True (@('DISCOVERED','LICENSE_VERIFIED','CANDIDATE','CONDITIONED','REUSABLE_PACK') -contains $admission) "invalid admission status in $(Relative $Path): $admission"
  $verifiedAt = [regex]::Match($metadata, '(?m)^verified_at:\s*"(?<value>\d{4}-\d{2}-\d{2})"\s*$')
  Assert-True $verifiedAt.Success "verified_at must be an ISO date in $(Relative $Path)"
  $parsedDate = [datetime]::MinValue
  $validDate = [datetime]::TryParseExact(
    $verifiedAt.Groups['value'].Value,
    'yyyy-MM-dd',
    [Globalization.CultureInfo]::InvariantCulture,
    [Globalization.DateTimeStyles]::None,
    [ref]$parsedDate
  )
  Assert-True $validDate "verified_at is not a real calendar date in $(Relative $Path)"
}


# V374 binds the narrow audit to its sources. Not product admission or live freshness.
function Assert-CoreClaimAudit([object] $Audit) {
  $expectedIds = @('DASHBOARDS','FX','GIFT-CARDS','HELP-CENTER','I18N','LOYALTY','MARKETING','OBSERVABILITY','ONBOARDING','PAYROLL','POS','PROMOTIONS','REFERRALS','REMINDERS','REVIEWS','SEO','SLO','SOCIAL-POSTING','SURVEYS','WAITLIST','WARRANTY-CLAIMS' | ForEach-Object { "GO-$_-CORE" })
  $rows = @($Audit.cores)
  Assert-True ($Audit.schema_version -eq 1 -and $Audit.audit -ceq 'TEST-02') 'core audit identity'
  Assert-True ($rows.Count -eq 21) 'core audit requires 21 rows'
  $seen = @{}
  $sourceIds = @($Audit.sources | ForEach-Object id)
  foreach ($row in $rows) {
    Assert-True ($row.id -cin $expectedIds -and -not $seen.ContainsKey($row.id)) 'unknown or duplicate audited core'
    $seen[$row.id] = $true
    $expectedPath = 'implementation_packs/' + $row.id.Replace('-','_') + '.md'
    Assert-True ($row.path -ceq $expectedPath) "audited core path mismatch: $($row.id)"
    $full = Join-Path $libraryRoot $row.path
    Assert-True ((Get-FileHash -LiteralPath $full -Algorithm SHA256).Hash.ToLowerInvariant() -ceq $row.sha256) "stale core audit: $($row.id)"
    $text = [IO.File]::ReadAllText($full)
    $claim = [regex]::Match($text, '(?m)^claim: "(?<value>.*)"\r?$').Groups['value'].Value
    $admission = [regex]::Match($text, '(?m)^  admission: (?<value>\w+)\r?$').Groups['value'].Value
    $version = [regex]::Match($text, '(?m)^pack_version: "(?<value>[^"]+)"\r?$').Groups['value'].Value
    Assert-True ($claim -ceq $row.claim -and $version -ceq $row.version -and $admission -ceq $row.admission) 'audited claim/version/admission drift'
    Assert-True ($row.license_ref -ceq 'LICENSE.md' -and (Get-FileHash -LiteralPath (Join-Path $libraryRoot 'LICENSE.md') -Algorithm SHA256).Hash.ToLowerInvariant() -ceq $row.license_sha256) 'audited license drift'
    Assert-True ($row.source_is_enterprise_upstream -ceq $false -and $row.candidate_promotion -ceq $false) 'false upstream attribution or promotion'
    Assert-True (-not [string]::IsNullOrWhiteSpace($row.conditions) -and -not [string]::IsNullOrWhiteSpace($row.comparison)) 'missing core conditions/comparison'
    Assert-True (@($row.source_refs).Count -gt 0 -and @($row.source_refs | Where-Object { $_ -cnotin $sourceIds }).Count -eq 0) 'missing core source authority'
    Assert-True (@($row.reference_test_names).Count -gt 0 -and @($row.reference_test_names).Count -eq $row.reference_test_count) 'missing core test mapping'
    Assert-True (@($row.reference_test_names | Sort-Object -Unique).Count -eq $row.reference_test_count) 'duplicate core test mapping'
    foreach ($testName in $row.reference_test_names) {
      Assert-True ($testName -cmatch '^Test[A-Za-z0-9_]+$' -and [regex]::IsMatch($text, '(?m)^func ' + [regex]::Escape($testName) + '\(')) 'audited test is not in canonical source'
    }
    Assert-True (@($row.files).Count -eq 2) 'core source manifest differs'
    foreach ($file in $row.files) {
      Assert-True ($file.provenance -ceq 'AUTHORED' -and $file.license -ceq 'LicenseRef-Workspace-Owner') 'core provenance changed'
      Assert-True ($file.path -match '^internal/[a-z0-9_]+/[a-z0-9_]+\.go$' -and $file.sha256 -cmatch '^[a-f0-9]{64}$') 'invalid core file identity'
      $fileHeading = '### FILE: ' + [char]96 + $file.path + [char]96
      $blockPattern = '(?s)' + [regex]::Escape($fileHeading) + '(?:(?!### FILE:).)*?sha256: "' + $file.sha256 + '"'
      Assert-True ($text.Contains('CREATE ' + $file.path) -and [regex]::IsMatch($text, $blockPattern)) 'core source hash no longer matches its canonical block'
    }
  }
  Assert-True ($seen.Count -eq 21) 'incomplete core census'
}


function Assert-CoreClaimAuditRegression([object] $Audit) {
  $original = ConvertTo-Json -InputObject $Audit -Depth 100
$cases=@(
 @{name='missing row';mutate={param($a) $a.cores=@($a.cores | Select-Object -Skip 1)}},
 @{name='duplicate row';mutate={param($a) $a.cores[1]=$a.cores[0]}},
 @{name='pack tamper';mutate={param($a) $a.cores[0].sha256='0'*64}},
 @{name='claim tamper';mutate={param($a) $a.cores[0].claim='invented authority'}},
 @{name='promoted candidate';mutate={param($a) $a.cores[7].admission='REUSABLE_PACK'}},
 @{name='false upstream';mutate={param($a) $a.cores[0].source_is_enterprise_upstream=$true}},
 @{name='lost condition';mutate={param($a) $a.cores[0].conditions=''}},
 @{name='missing tests';mutate={param($a) $a.cores[0].reference_test_names=@()}},
 @{name='unknown source';mutate={param($a) $a.cores[0].source_refs=@('invented')}},
 @{name='invented test';mutate={param($a) $a.cores[0].reference_test_names[0]='TestInventedNeverExecuted'}},
 @{name='wrong license';mutate={param($a) $a.cores[0].license_sha256='0'*64}},
 @{name='path escape';mutate={param($a) $a.cores[0].path='../outside'}},
 @{name='swapped file hashes';mutate={param($a) $first=$a.cores[0].files[0].sha256;$a.cores[0].files[0].sha256=$a.cores[0].files[1].sha256;$a.cores[0].files[1].sha256=$first}}
)
foreach($case in $cases){
 $a=$original | ConvertFrom-Json -Depth 100
 & $case.mutate $a
 $rejected=$false
 try{Assert-CoreClaimAudit $a}catch{$rejected=$true}
 if(-not $rejected){throw ('negative accepted: '+$case.name)}
}

}

function Read-PackMetadata([string] $Path) {
  $text = [IO.File]::ReadAllText($Path, $utf8)
  Assert-PackContract $Path $text
  $id = [regex]::Match($text, '(?m)^pack_id: "(?<value>[^"]+)"\r?$').Groups['value'].Value
  $version = [regex]::Match($text, '(?m)^pack_version: "(?<value>[^"]+)"\r?$').Groups['value'].Value
  Assert-True ($id -ne '') "pack_id missing in $(Relative $Path)"
  Assert-True ($version -match '^\d+\.\d+\.\d+$') "invalid pack_version in $(Relative $Path)"
  [pscustomobject]@{ Id = $id; Version = $version; Path = Relative $Path }
}

function Read-Manifest([string] $Path) {
  $lines = [IO.File]::ReadAllLines($Path, $utf8)
  $inside = $false
  $paths = [Collections.Generic.List[string]]::new()
  foreach ($line in $lines) {
    if (-not $inside -and $line -match '^## .*manifest') { $inside = $true; continue }
    if ($inside -and $line -match '^## ') { break }
    if ($inside -and $line -match '^CREATE (?<path>\S+)$') { $paths.Add($Matches['path']) }
  }
  Assert-True ($paths.Count -gt 0) "manifest missing or empty in $(Relative $Path)"
  Assert-True (($paths | Sort-Object -Unique).Count -eq $paths.Count) "duplicate manifest path in $(Relative $Path)"
  $paths
}

try {
  [void](New-Item -ItemType Directory -Path $tempRoot)
  Assert-FranchiseProfileSelectionRegression

  foreach ($required in @('README.md','LICENSE.md','THIRD_PARTY_NOTICES.md','AGENTS.md','CLAUDE.md','AGENT_SYSTEM_START.md','START_ANY_PROJECT.md','materialize_markdown_pack.ps1','VERIFY_LIBRARY.ps1','VERIFY_EXECUTABLE_LIBRARY.ps1','CREATE_PORTABLE_ARCHIVE.ps1','INSTALL_AGENT_BRIDGE.ps1','markdown_system/test_update_pack_from_tree.ps1','markdown_system/test_install_agent_bridge.ps1','markdown_system/MARKDOWN_SYSTEM_READINESS.md','markdown_system/FAILURE_LEARNING_CONTRACT.md','markdown_system/PROJECT_FAILURE_LESSONS_TEMPLATE.md','markdown_system/LIBRARY_FAILURE_LEARNING_LEDGER.md','markdown_system/DEPENDENCY_UPDATE_CONTRACT.md','markdown_system/PROJECT_DEPENDENCY_UPDATE_RECORD_TEMPLATE.md','markdown_system/AUTHORITY_FRESHNESS_AND_SELF_CORRECTION_CONTRACT.md','markdown_system/PROJECT_AUTHORITY_FRESHNESS_TEMPLATE.md','markdown_system/ZERO_COST_VULNERABILITY_MONITORING_PROFILE.md','markdown_system/PROJECT_VULNERABILITY_MONITORING_RECORD_TEMPLATE.md','markdown_system/PROJECT_OFFICIAL_SOURCE_PROFILE_RECORD_TEMPLATE.md','markdown_system/CAPABILITY_GAP_RESOLUTION_PACK_PLAN.md','markdown_system/MICROSOFT_DEVSKIM_ADAPTED_SAST_PACK_PLAN.md','markdown_system/PORTABLE_SIGNED_RELEASE_EVIDENCE_PACK_PLAN.md','markdown_system/PROJECT_READINESS_GATE_PACK_PLAN.md','markdown_system/PROJECT_INITIALIZATION_PACK_PLAN.md','markdown_system/DOCUMENT_PIPELINE_ROUTING_PACK_PLAN.md','markdown_system/DURABLE_DOCUMENT_PIPELINE_PACK_PLAN.md','markdown_system/AZURE_DOCUMENT_RUNTIME_PACK_PLAN.md','markdown_system/GOOGLE_DOCUMENT_RUNTIME_PACK_PLAN.md','markdown_system/AWS_TEXTRACT_DOCUMENT_RUNTIME_PACK_PLAN.md','markdown_system/PADDLEOCR_LOCAL_RUNTIME_PACK_PLAN.md','markdown_system/MARKITDOWN_LOCAL_RUNTIME_PACK_PLAN.md','markdown_system/SECURE_LOCAL_FILE_INGESTION_PACK_PLAN.md','markdown_system/STRICT_DOCUMENT_FIELD_EVALUATION_PACK_PLAN.md','markdown_system/TESSERA_POSIX_EVIDENCE_LOG_PACK_PLAN.md','markdown_system/AWS_ENTERPRISE_ADAPTERS_PACK_PLAN.md','markdown_system/AWS_SECURE_DOCUMENT_INTAKE_PACK_PLAN.md','markdown_system/AWS_SECURE_EMAIL_ATTACHMENT_PROCESSING_PACK_PLAN.md','markdown_system/AWS_IDP_POSTGRES_PERSISTENCE_PACK_PLAN.md','markdown_system/GOOGLE_ADS_REPORTING_PACK_PLAN.md','markdown_system/META_ADS_REPORTING_PACK_PLAN.md','markdown_system/TIKTOK_ADS_REPORTING_PACK_PLAN.md','markdown_system/META_WHATSAPP_CLOUD_PACK_PLAN.md','markdown_system/FIREBASE_PUSH_PACK_PLAN.md','markdown_system/MERCADOLIBRE_MARKETPLACE_PACK_PLAN.md','markdown_system/AMAZON_SPAPI_CATALOG_PACK_PLAN.md','markdown_system/GOOGLE_MERCHANT_PRODUCT_SYNC_PACK_PLAN.md','markdown_system/PAYMENT_WEBHOOK_ADAPTERS_PACK_PLAN.md','markdown_system/MICROSOFT_BUSINESS_CENTRAL_PLATFORM_PACK_PLAN.md','markdown_system/MICROSOFT_AVM_SECURE_SFTP_INTAKE_PACK_PLAN.md','markdown_system/ENTERPRISE_BACKEND_PACK_PLAN.md','markdown_system/ENTERPRISE_WEB_PACK_PLAN.md')) {
    Assert-True (Test-Path -LiteralPath (Join-Path $libraryRoot $required) -PathType Leaf) "required file missing: $required"
  }
  $officialInvoicePlan = 'markdown_system/AZURE_DOCUMENT_INTELLIGENCE_OFFICIAL_INVOICE_SAMPLE_PACK_PLAN.md'
  Assert-True (Test-Path -LiteralPath (Join-Path $libraryRoot $officialInvoicePlan) -PathType Leaf) "required file missing: $officialInvoicePlan"
  $amazonShippingPlan = 'markdown_system/AMAZON_SPAPI_SHIPPING_TRACKING_PACK_PLAN.md'
  Assert-True (Test-Path -LiteralPath (Join-Path $libraryRoot $amazonShippingPlan) -PathType Leaf) "required file missing: $amazonShippingPlan"
  $amazonEasyShipPlan = 'markdown_system/AMAZON_SPAPI_EASYSHIP_HANDOVER_PACK_PLAN.md'
  Assert-True (Test-Path -LiteralPath (Join-Path $libraryRoot $amazonEasyShipPlan) -PathType Leaf) "required file missing: $amazonEasyShipPlan"
  $amazonFulfillmentDeliveryEvidencePlan = 'markdown_system/AMAZON_SPAPI_FULFILLMENT_DELIVERY_EVIDENCE_PACK_PLAN.md'
  Assert-True (Test-Path -LiteralPath (Join-Path $libraryRoot $amazonFulfillmentDeliveryEvidencePlan) -PathType Leaf) "required file missing: $amazonFulfillmentDeliveryEvidencePlan"
  $amazonSupplySourcesPlan = 'markdown_system/AMAZON_SPAPI_SUPPLY_SOURCES_PACK_PLAN.md'
  Assert-True (Test-Path -LiteralPath (Join-Path $libraryRoot $amazonSupplySourcesPlan) -PathType Leaf) "required file missing: $amazonSupplySourcesPlan"
  $amazonMliInventoryPlan = 'markdown_system/AMAZON_SPAPI_MLI_INVENTORY_PACK_PLAN.md'
  Assert-True (Test-Path -LiteralPath (Join-Path $libraryRoot $amazonMliInventoryPlan) -PathType Leaf) "required file missing: $amazonMliInventoryPlan"
  $amazonExternalInventoryPlan = 'markdown_system/AMAZON_SPAPI_EXTERNAL_FULFILLMENT_INVENTORY_PACK_PLAN.md'
  Assert-True (Test-Path -LiteralPath (Join-Path $libraryRoot $amazonExternalInventoryPlan) -PathType Leaf) "required file missing: $amazonExternalInventoryPlan"
  $officialGoogleProcessPlan = 'markdown_system/GOOGLE_DOCUMENT_AI_OFFICIAL_PROCESS_SAMPLE_PACK_PLAN.md'
  Assert-True (Test-Path -LiteralPath (Join-Path $libraryRoot $officialGoogleProcessPlan) -PathType Leaf) "required file missing: $officialGoogleProcessPlan"
  $officialGoogleLifecyclePlan = 'markdown_system/GOOGLE_DOCUMENT_AI_OFFICIAL_LIFECYCLE_PACK_PLAN.md'
  Assert-True (Test-Path -LiteralPath (Join-Path $libraryRoot $officialGoogleLifecyclePlan) -PathType Leaf) "required file missing: $officialGoogleLifecyclePlan"
  $daprOutboxPlan = 'markdown_system/DAPR_OFFICIAL_TRANSACTIONAL_OUTBOX_PACK_PLAN.md'
  Assert-True (Test-Path -LiteralPath (Join-Path $libraryRoot $daprOutboxPlan) -PathType Leaf) "required file missing: $daprOutboxPlan"
  $pgDurableHandoffPlan = 'markdown_system/MICROSOFT_PG_DURABLE_HUMAN_HANDOFF_PACK_PLAN.md'
  Assert-True (Test-Path -LiteralPath (Join-Path $libraryRoot $pgDurableHandoffPlan) -PathType Leaf) "required file missing: $pgDurableHandoffPlan"

  $releaseFiles = @(Get-ReleaseSourceFiles)
  Assert-DistributionManifest $releaseFiles
  foreach ($file in $releaseFiles) {
    $bytes = [IO.File]::ReadAllBytes($file.FullName)
    try {
      $text = $strictUtf8.GetString($bytes)
    } catch {
      throw "VERIFY_LIBRARY_FAILED: release file is not valid UTF-8: $(Relative $file.FullName)"
    }
    Assert-True (-not $text.Contains([char]0)) "NUL byte in release file: $(Relative $file.FullName)"
    $localHome = Find-MachineLocalHomePath $text
    Assert-True (-not $localHome.Success) "machine-local home path in release file $(Relative $file.FullName): $($localHome.Value)"
  }

  $thirdPartyText = [IO.File]::ReadAllText((Join-Path $libraryRoot 'THIRD_PARTY_NOTICES.md'), $utf8)
  Assert-True ($thirdPartyText.Contains('f809c396520df2d7b201a9ccc5378d822b728ed3', [StringComparison]::Ordinal)) 'TikTok exact commit missing from third-party notices'
  Assert-True (-not $thirdPartyText.Contains('f809c396520df2d7b2019ccc5378d822b728ed3', [StringComparison]::Ordinal)) 'known malformed TikTok commit remains in third-party notices'
  $admissionLedgerText = [IO.File]::ReadAllText((Join-Path $libraryRoot 'markdown_system/REVESTEX_ELITE_PUBLIC_CODE_ADMISSION_LEDGER.md'), $utf8)
  Assert-True ($admissionLedgerText.Contains('cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30', [StringComparison]::Ordinal)) 'Google Document AI exact license hash missing from admission ledger'
  Assert-True (-not $admissionLedgerText.Contains('cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc5e3d30', [StringComparison]::Ordinal)) 'known malformed Google Document AI license hash remains in admission ledger'

  $failureContractText = [IO.File]::ReadAllText((Join-Path $libraryRoot 'markdown_system/FAILURE_LEARNING_CONTRACT.md'), $utf8)
  foreach ($requiredField in @('failure_id:','fingerprint:','recurrence_count:','status:','severity:','classification:','violated_invariant:','root_cause:','canonical_correction:','regression_test:','clean_rebuild_evidence:','residual_risk:')) {
    Assert-True ($failureContractText.Contains($requiredField, [StringComparison]::Ordinal)) "failure-learning contract field missing: $requiredField"
  }
  $failureTemplateText = [IO.File]::ReadAllText((Join-Path $libraryRoot 'markdown_system/PROJECT_FAILURE_LESSONS_TEMPLATE.md'), $utf8)
  foreach ($requiredField in @('failure_id:','fingerprint:','recurrence_count:','violated_invariant:','root_cause:','canonical_correction:','regression_test:','clean_rebuild_evidence:')) {
    Assert-True ($failureTemplateText.Contains($requiredField, [StringComparison]::Ordinal)) "failure-learning template field missing: $requiredField"
  }
  $failureLedgerText = [IO.File]::ReadAllText((Join-Path $libraryRoot 'markdown_system/LIBRARY_FAILURE_LEARNING_LEDGER.md'), $utf8)
  $failureIds = @([regex]::Matches($failureLedgerText, '(?m)`(?<id>(?:LIB|UP)-FAIL-[0-9]{3,})`') | ForEach-Object { $_.Groups['id'].Value })
  Assert-True ($failureIds.Count -gt 0) 'library failure ledger contains no failure IDs'
  Assert-True (($failureIds | Sort-Object -Unique).Count -eq $failureIds.Count) 'duplicate failure ID in library failure ledger'
  $openLocalFailures = @([regex]::Matches($failureLedgerText, '(?m)^\|\s*`LIB-FAIL-[0-9]{3,}`\s*\|.*\|\s*`OPEN`\s*\|\s*$'))
  Assert-True ($openLocalFailures.Count -eq 0) "library failure ledger contains $($openLocalFailures.Count) open local failure rows"
  foreach ($router in @('AGENTS.md','AGENT_SYSTEM_START.md','START_ANY_PROJECT.md','markdown_system/PROJECT_START_READINESS_GATE.md')) {
    $routerText = [IO.File]::ReadAllText((Join-Path $libraryRoot $router), $utf8)
    Assert-True ($routerText.Contains('PROJECT_FAILURE_LESSONS', [StringComparison]::Ordinal)) "failure ledger not routed by $router"
    Assert-True ($routerText.Contains('PROJECT_DEPENDENCY_UPDATE_RECORD', [StringComparison]::Ordinal)) "dependency update record not routed by $router"
    Assert-True ($routerText.Contains('PROJECT_AUTHORITY_FRESHNESS_RECORD', [StringComparison]::Ordinal)) "authority freshness record not routed by $router"
    Assert-True ($routerText.Contains('PROJECT_VULNERABILITY_MONITORING_RECORD', [StringComparison]::Ordinal)) "vulnerability monitoring record not routed by $router"
    Assert-True ($routerText.Contains('PROJECT_OFFICIAL_SOURCE_PROFILE_RECORD', [StringComparison]::Ordinal)) "official source profile record not routed by $router"
    Assert-True ($routerText.Contains('CAPABILITY_GAP_RESOLUTION', [StringComparison]::Ordinal)) "capability gap resolution not routed by $router"
    Assert-True ($routerText.Contains('MICROSOFT_DEVSKIM_ADAPTED_SAST', [StringComparison]::Ordinal)) "adapted DevSkim lint gate not routed by $router"
    Assert-True ($routerText.Contains('render_project_advisory.py', [StringComparison]::Ordinal)) "executable advisory renderer not routed by $router"
  }
  $readinessGateText = [IO.File]::ReadAllText((Join-Path $libraryRoot 'markdown_system/PROJECT_START_READINESS_GATE.md'), $utf8)
  foreach ($requiredTerm in @('OFFICIAL_PLATFORM','VERBATIM|DEPENDENCY_PIN','NO_SOURCE','ADAPTED|AUTHORED')) {
    Assert-True ($readinessGateText.Contains($requiredTerm, [StringComparison]::Ordinal)) "official-only platform invariant missing from readiness gate: $requiredTerm"
  }
  $readinessTemplateText = [IO.File]::ReadAllText((Join-Path $libraryRoot 'markdown_system/PROJECT_READINESS_RECORD_TEMPLATE.md'), $utf8)
  foreach ($requiredTerm in @('OFFICIAL_PLATFORM','official_platform:','product_provenance_allowed: [VERBATIM, DEPENDENCY_PIN]','authored_or_adapted_product_code_allowed: false','uncovered_capabilities:','advisory_prompt_ref','PROJECT_ADVISORY_A.md')) {
    Assert-True ($readinessTemplateText.Contains($requiredTerm, [StringComparison]::Ordinal)) "official-only platform invariant missing from readiness template: $requiredTerm"
  }
  $agentsBytes = [IO.File]::ReadAllBytes((Join-Path $libraryRoot 'AGENTS.md')).Length
  Assert-True ($agentsBytes -le 32768) "root AGENTS.md exceeds the official 32 KiB default: $agentsBytes bytes"
  $bridgeInstallerText = [IO.File]::ReadAllText((Join-Path $libraryRoot 'INSTALL_AGENT_BRIDGE.ps1'), $utf8)
  foreach ($requiredTerm in @('ELITE-ENGINEERING-LIBRARY:BEGIN','elite-engineering-library','PROJECT_START_READINESS_GATE.md','render_project_advisory.py','32768','refusing to replace an unmanaged skill','ELITE_AGENT_BRIDGE_READY')) {
    Assert-True ($bridgeInstallerText.Contains($requiredTerm, [StringComparison]::Ordinal)) "agent bridge installer term missing: $requiredTerm"
  }
  $dependencyContractText = [IO.File]::ReadAllText((Join-Path $libraryRoot 'markdown_system/DEPENDENCY_UPDATE_CONTRACT.md'), $utf8)
  foreach ($requiredTerm in @('current_digest_or_commit:','license_expression:','update_policy:','CANDIDATE_GREEN','PROMOTED','ROLLED_BACK','PROJECT_FAILURE_LESSONS.md','OFFICIAL_FIXED_RELEASE','OFFICIAL_PATCH_UNRELEASED','NO_UPSTREAM_FIX','ADAPTED_PATCH','FALSE_POSITIVE_OR_NOT_REACHABLE')) {
    Assert-True ($dependencyContractText.Contains($requiredTerm, [StringComparison]::Ordinal)) "dependency update contract term missing: $requiredTerm"
  }
  $freshnessContractText = [IO.File]::ReadAllText((Join-Path $libraryRoot 'markdown_system/AUTHORITY_FRESHNESS_AND_SELF_CORRECTION_CONTRACT.md'), $utf8)
  foreach ($requiredTerm in @('PROJECT_AUTHORITY_FRESHNESS_RECORD.md','STALE_BLOCKED','CONFLICT_UNRESOLVED','SUPERSEDED','previous_statement')) {
    Assert-True ($freshnessContractText.Contains($requiredTerm, [StringComparison]::Ordinal)) "authority freshness contract term missing: $requiredTerm"
  }
  $zeroCostMonitoringText = [IO.File]::ReadAllText((Join-Path $libraryRoot 'markdown_system/ZERO_COST_VULNERABILITY_MONITORING_PROFILE.md'), $utf8)
  foreach ($requiredTerm in @('ZERO_COST_BASELINE_READY','NOT_READY_FOR_PRODUCTION_RUNTIME_MONITORING','PROJECT_VULNERABILITY_MONITORING_RECORD.md','25e42f5ef6711fd8c0fb45390972205891dd44c6bd02ac93f0f63e8e98d9bfb6','StartWhenAvailable','RunOnlyIfNetworkAvailable','scheduled OSV workflow','osv-scanner fix')) {
    Assert-True ($zeroCostMonitoringText.Contains($requiredTerm, [StringComparison]::Ordinal)) "zero-cost monitoring term missing: $requiredTerm"
  }

  $updaterTestOutput = @(& (Join-Path $libraryRoot 'markdown_system/test_update_pack_from_tree.ps1'))
  $localSecretsTestOutput = @(& (Join-Path $libraryRoot 'markdown_system/test_project_local_secrets.ps1'))
  Assert-True (@($updaterTestOutput | Where-Object { $_ -match '^PASS:' }).Count -eq 1) 'update_pack_from_tree regression did not report PASS'
  $bridgeTestOutput = @(& (Join-Path $libraryRoot 'markdown_system/test_install_agent_bridge.ps1') 6>$null)
  Assert-True (@($bridgeTestOutput | Where-Object { $_ -match '^PASS:' }).Count -eq 1) 'install_agent_bridge regression did not report PASS'
  $stateTestOutput = @(& (Join-Path $libraryRoot 'markdown_system/test_library_maintenance_state.ps1'))
  Assert-True (@($stateTestOutput | Where-Object { $_ -match '^PASS:' }).Count -eq 1) 'local execution distribution regression did not report PASS'
  $toolchainTestOutput = @(& (Join-Path $libraryRoot 'markdown_system/test_toolchain_resolution.ps1'))
  Assert-True (@($toolchainTestOutput | Where-Object { $_ -match '^PASS: toolchain resolution 46 checks;' }).Count -eq 1) 'toolchain resolution regression did not report 46 checks PASS'
  $entryTestOutput = @(& (Join-Path $libraryRoot 'markdown_system/test_agent_entry_lifecycle.ps1'))
  Assert-True (@($entryTestOutput | Where-Object { $_ -match '^PASS: agent entry lifecycle 55 checks;' }).Count -eq 1) 'agent entry lifecycle regression did not report 55 checks PASS'

  $packFiles = @(Get-ChildItem -LiteralPath (Join-Path $libraryRoot 'implementation_packs') -Filter '*.md' | Where-Object Name -ne 'README.md' | Sort-Object Name)
  Assert-True ($packFiles.Count -gt 0) 'no implementation packs found'
  $packs = @{}
  $materializedFiles = 0
  $provenanceCounts = @{ AUTHORED = 0; ADAPTED = 0; VERBATIM = 0 }
  foreach ($packFile in $packFiles) {
    $metadata = Read-PackMetadata $packFile.FullName
    Assert-True (-not $packs.ContainsKey($metadata.Id)) "duplicate pack_id: $($metadata.Id)"
    $packs[$metadata.Id] = $metadata
    $manifest = @(Read-Manifest $packFile.FullName)
    $destination = Join-Path $tempRoot ('pack-' + [guid]::NewGuid().ToString('N'))
    & $materializer -PackFile $packFile.FullName -Destination $destination | Out-Null
    if ($metadata.Id -eq 'EXECUTION-VALIDATOR') {
      $executionValidator = Join-Path $destination 'engineering_execution_kit/validate_execution_state.py'
    }
    $actual = @(Get-ChildItem -LiteralPath $destination -Recurse -File | ForEach-Object { [IO.Path]::GetRelativePath($destination, $_.FullName).Replace('\','/') } | Sort-Object)
    $expected = @($manifest | Sort-Object)
    Assert-True ($actual.Count -eq $expected.Count) "manifest/file count differs in $($packFile.Name): manifest=$($expected.Count), materialized=$($actual.Count)"
    Assert-True ((Compare-Object $expected $actual).Count -eq 0) "manifest paths differ from materialized paths in $($packFile.Name)"
    $materializedFiles += $actual.Count
    foreach ($provenance in @(Read-MaterializationProvenance $packFile.FullName)) {
      $provenanceCounts[$provenance]++
    }
  }
  $coreAuditText = [IO.File]::ReadAllText((Join-Path $libraryRoot 'reconstruction_evidence/CORE_CLAIM_CURRENT_V402.md'), $utf8)
  $coreFence = [string][char]96 * 3
  $coreAuditMatch = [regex]::Match($coreAuditText, '(?s)' + $coreFence + 'json\s*(?<json>.*?)\s*' + $coreFence)
  Assert-True $coreAuditMatch.Success 'core audit JSON missing'
  $coreAudit = $coreAuditMatch.Groups['json'].Value | ConvertFrom-Json -Depth 100
  Assert-CoreClaimAudit $coreAudit
  Assert-CoreClaimAuditRegression $coreAudit
  Write-Output 'CORE_CLAIM_AUDIT_21_SOURCE_BINDINGS_VERIFIED_NOT_PRODUCT_ADMISSION'
  $provenanceTotal = $provenanceCounts.AUTHORED + $provenanceCounts.ADAPTED + $provenanceCounts.VERBATIM
  if (Test-Path -LiteralPath (Join-Path $libraryRoot 'PROJECT_EXECUTION_STATE.json')) {
    $pythonCommand = Get-Command python -CommandType Application -ErrorAction Stop | Select-Object -First 1
    Assert-True (Test-Path -LiteralPath $executionValidator -PathType Leaf) 'execution validator was not reconstructed'
    & $pythonCommand.Source $executionValidator (Join-Path $libraryRoot 'PROJECT_EXECUTION_STATE.json') --project-root $libraryRoot --events (Join-Path $libraryRoot 'PROJECT_EXECUTION_EVENTS.jsonl') --level resume | Write-Output
    Assert-True ($LASTEXITCODE -eq 0) 'local execution checkpoint validation failed'
    Write-Output 'LOCAL_EXECUTION_CHECKPOINT_VERIFIED_NOT_PRODUCT_READINESS'
  }
  Assert-True ($provenanceTotal -eq $materializedFiles) "provenance/materialization count differs: provenance=$provenanceTotal materialized=$materializedFiles"
  $notices = [IO.File]::ReadAllText((Join-Path $libraryRoot 'THIRD_PARTY_NOTICES.md'), $utf8)
  $expectedProvenanceLine = "Current-Provenance-Counts: AUTHORED=$($provenanceCounts.AUTHORED); ADAPTED=$($provenanceCounts.ADAPTED); VERBATIM=$($provenanceCounts.VERBATIM); TOTAL=$provenanceTotal"
  Assert-True ($notices.Contains($expectedProvenanceLine, [StringComparison]::Ordinal)) "THIRD_PARTY_NOTICES provenance ledger drifted; expected $expectedProvenanceLine"

  foreach ($planFile in Get-ChildItem -LiteralPath (Join-Path $libraryRoot 'markdown_system') -Filter '*PACK_PLAN.md' | Sort-Object Name) {
    $text = [IO.File]::ReadAllText($planFile.FullName, $utf8)
    $match = [regex]::Match($text, '(?s)```json\s*(?<json>.*?)\s*```')
    Assert-True $match.Success "JSON plan missing in $($planFile.Name)"
    $plan = $match.Groups['json'].Value | ConvertFrom-Json -Depth 100
    Assert-True ($plan.compositionVersion -eq '1.0') "unsupported compositionVersion in $($planFile.Name)"
    $planPackIds = @($plan.packs | ForEach-Object { [string]$_.packId })
    Assert-True ($planPackIds.Count -gt 0) "plan contains no packs in $($planFile.Name)"
    Assert-True (($planPackIds | Sort-Object -Unique).Count -eq $planPackIds.Count) "duplicate packId in $($planFile.Name)"
    if ($planFile.Name -eq 'FRANCHISE_COMPLETE_PACK_PLAN.md') {
      Assert-FranchiseProfileSelection $planPackIds
    }
    foreach ($entry in @($plan.packs)) {
      Assert-True ($packs.ContainsKey([string]$entry.packId)) "unknown packId $($entry.packId) in $($planFile.Name)"
      $known = $packs[[string]$entry.packId]
      Assert-True ($known.Version -eq [string]$entry.version) "version mismatch for $($entry.packId) in $($planFile.Name): plan=$($entry.version), pack=$($known.Version)"
      Assert-True ($known.Path -eq ([string]$entry.path).Replace('\','/')) "path mismatch for $($entry.packId) in $($planFile.Name)"
    }
  }

  $compositorPack = Join-Path $libraryRoot 'implementation_packs/MARKDOWN_COMPOSITOR_CORE.md'
  $compositorRoot = Join-Path $tempRoot 'compositor'
  & $materializer -PackFile $compositorPack -Destination $compositorRoot | Out-Null
  $compositor = Join-Path $compositorRoot 'tools/compose-markdown-project.ps1'
  $eofOutput = @(& (Join-Path $libraryRoot 'markdown_system/test_pack_exact_eof.ps1') -Materializer $materializer -Composer $compositor -Updater (Join-Path $libraryRoot 'markdown_system/update_pack_from_tree.ps1') -EvidenceRoot (Join-Path $tempRoot 'exact-eof-regression'))
  Assert-True (@($eofOutput | Where-Object { $_ -match '^PASS: [0-9]+ exact EOF reconstruction/negative/updater cases$' }).Count -eq 1) 'exact EOF regression did not report PASS'
  $publicPolicyOutput = @(& (Join-Path $libraryRoot 'markdown_system/test_public_release_policy.ps1'))
  Assert-True (@($publicPolicyOutput | Where-Object { $_ -match '^PASS: public release policy ' }).Count -eq 1) 'public release policy regression did not report PASS'

  $profileResults = [Collections.Generic.List[object]]::new()
  foreach ($profile in @('PROJECT_READINESS_GATE_PACK_PLAN.md','CAPABILITY_GAP_RESOLUTION_PACK_PLAN.md','PNPM_ARTIFACT_SELECTION_PACK_PLAN.md','PROJECT_INITIALIZATION_PACK_PLAN.md','DOCUMENT_PIPELINE_ROUTING_PACK_PLAN.md','DURABLE_DOCUMENT_PIPELINE_PACK_PLAN.md','AZURE_DOCUMENT_RUNTIME_PACK_PLAN.md','GOOGLE_DOCUMENT_RUNTIME_PACK_PLAN.md','AWS_TEXTRACT_DOCUMENT_RUNTIME_PACK_PLAN.md','AWS_POWERTOOLS_IDEMPOTENT_SQS_BATCH_PACK_PLAN.md','AWS_DURABLE_OBJECT_EVENT_WORKER_PACK_PLAN.md','PADDLEOCR_LOCAL_RUNTIME_PACK_PLAN.md','MARKITDOWN_LOCAL_RUNTIME_PACK_PLAN.md','SECURE_LOCAL_FILE_INGESTION_PACK_PLAN.md','STRICT_DOCUMENT_FIELD_EVALUATION_PACK_PLAN.md','TESSERA_POSIX_EVIDENCE_LOG_PACK_PLAN.md','AWS_ENTERPRISE_ADAPTERS_PACK_PLAN.md','AWS_SECURE_DOCUMENT_INTAKE_PACK_PLAN.md','AWS_SECURE_EMAIL_ATTACHMENT_PROCESSING_PACK_PLAN.md','AWS_IDP_POSTGRES_PERSISTENCE_PACK_PLAN.md','GOOGLE_ADS_REPORTING_PACK_PLAN.md','META_ADS_REPORTING_PACK_PLAN.md','TIKTOK_ADS_REPORTING_PACK_PLAN.md','TIKTOK_LEAD_ADAPTER_PACK_PLAN.md','META_WHATSAPP_CLOUD_PACK_PLAN.md','FIREBASE_PUSH_PACK_PLAN.md','MERCADOLIBRE_MARKETPLACE_PACK_PLAN.md','AMAZON_SPAPI_CATALOG_PACK_PLAN.md','GOOGLE_MERCHANT_PRODUCT_SYNC_PACK_PLAN.md','PAYMENT_WEBHOOK_ADAPTERS_PACK_PLAN.md','MICROSOFT_BUSINESS_CENTRAL_PLATFORM_PACK_PLAN.md','MICROSOFT_AVM_SECURE_SFTP_INTAKE_PACK_PLAN.md','ENTERPRISE_BACKEND_PACK_PLAN.md','ENTERPRISE_WEB_PACK_PLAN.md','FRANCHISE_COMPLETE_PACK_PLAN.md','FRANCHISE_SERVERLESS_PACK_PLAN.md','WINDOWS_REFERENCE_TELEMETRY_PACK_PLAN.md')) {
    $destination = Join-Path $tempRoot ([IO.Path]::GetFileNameWithoutExtension($profile))
    & $compositor -PlanFile (Join-Path $libraryRoot "markdown_system/$profile") -LibraryRoot $libraryRoot -Destination $destination | Out-Null
    Assert-True (Test-Path -LiteralPath (Join-Path $destination 'MATERIALIZATION_RECORD.md')) "composition record missing for $profile"
    $profileResults.Add([pscustomobject]@{ Profile = $profile; Files = (Get-ChildItem -LiteralPath $destination -Recurse -File | Measure-Object).Count - 1 })
  }
  $profile = 'MICROSOFT_KIOTA_OPENAPI_CLIENT_PACK_PLAN.md'
  $destination = Join-Path $tempRoot ([IO.Path]::GetFileNameWithoutExtension($profile))
  & $compositor -PlanFile (Join-Path $libraryRoot "markdown_system/$profile") -LibraryRoot $libraryRoot -Destination $destination | Out-Null
  Assert-True (Test-Path -LiteralPath (Join-Path $destination 'MATERIALIZATION_RECORD.md')) "composition record missing for $profile"
  $profileFiles = (Get-ChildItem -LiteralPath $destination -Recurse -File | Measure-Object).Count - 1
  Assert-True ($profileFiles -eq 7) "profile file count mismatch for $profile expected=7 actual=$profileFiles"
  $profileResults.Add([pscustomobject]@{ Profile = $profile; Files = $profileFiles })
  $profile = 'PORTABLE_SIGNED_RELEASE_EVIDENCE_PACK_PLAN.md'
  $destination = Join-Path $tempRoot ([IO.Path]::GetFileNameWithoutExtension($profile))
  & $compositor -PlanFile (Join-Path $libraryRoot "markdown_system/$profile") -LibraryRoot $libraryRoot -Destination $destination | Out-Null
  Assert-True (Test-Path -LiteralPath (Join-Path $destination 'MATERIALIZATION_RECORD.md')) "composition record missing for $profile"
  $profileFiles = (Get-ChildItem -LiteralPath $destination -Recurse -File | Measure-Object).Count - 1
  Assert-True ($profileFiles -eq 11) "profile file count mismatch for $profile expected=11 actual=$profileFiles"
  $profileResults.Add([pscustomobject]@{ Profile = $profile; Files = $profileFiles })
  $profile = 'NODE_OFFICIAL_RUNTIME_ADVISORY_PACK_PLAN.md'
  $destination = Join-Path $tempRoot ([IO.Path]::GetFileNameWithoutExtension($profile))
  $nodePlanText = [IO.File]::ReadAllText((Join-Path $libraryRoot "markdown_system/$profile"), $utf8)
  $nodePlanMatch = [regex]::Match($nodePlanText, '(?s)```json\s*(?<json>\{.*?\})\s*```')
  Assert-True $nodePlanMatch.Success "JSON plan missing in $profile"
  $nodePlan = $nodePlanMatch.Groups['json'].Value | ConvertFrom-Json -Depth 100
  Assert-True (@($nodePlan.packs).Count -eq 1) "profile must select exactly one pack: $profile"
  Assert-True ($nodePlan.packs[0].packId -eq 'NODE-OFFICIAL-RUNTIME-ADVISORY-GATE') "unexpected pack selected by $profile"
  & $compositor -PlanFile (Join-Path $libraryRoot "markdown_system/$profile") -LibraryRoot $libraryRoot -Destination $destination | Out-Null
  Assert-True (Test-Path -LiteralPath (Join-Path $destination 'MATERIALIZATION_RECORD.md')) "composition record missing for $profile"
  $profileFiles = (Get-ChildItem -LiteralPath $destination -Recurse -File | Measure-Object).Count - 1
  Assert-True ($profileFiles -eq 10) "profile file count mismatch for $profile expected=10 actual=$profileFiles"
  $profileResults.Add([pscustomobject]@{ Profile = $profile; Files = $profileFiles })
  $profile = 'GO_NATIVE_FUZZ_GATE_PACK_PLAN.md'
  $destination = Join-Path $tempRoot ([IO.Path]::GetFileNameWithoutExtension($profile))
  & $compositor -PlanFile (Join-Path $libraryRoot "markdown_system/$profile") -LibraryRoot $libraryRoot -Destination $destination | Out-Null
  Assert-True (Test-Path -LiteralPath (Join-Path $destination 'MATERIALIZATION_RECORD.md')) "composition record missing for $profile"
  $profileFiles = (Get-ChildItem -LiteralPath $destination -Recurse -File | Measure-Object).Count - 1
  Assert-True ($profileFiles -eq 4) "profile file count mismatch for $profile expected=4 actual=$profileFiles"
  $profileResults.Add([pscustomobject]@{ Profile = $profile; Files = $profileFiles })
  $profile = 'MICROSOFT_DEVSKIM_ADAPTED_SAST_PACK_PLAN.md'
  $destination = Join-Path $tempRoot ([IO.Path]::GetFileNameWithoutExtension($profile))
  & $compositor -PlanFile (Join-Path $libraryRoot "markdown_system/$profile") -LibraryRoot $libraryRoot -Destination $destination | Out-Null
  Assert-True (Test-Path -LiteralPath (Join-Path $destination 'MATERIALIZATION_RECORD.md')) "composition record missing for $profile"
  $profileFiles = (Get-ChildItem -LiteralPath $destination -Recurse -File | Measure-Object).Count - 1
  Assert-True ($profileFiles -eq 4) "profile file count mismatch for $profile expected=4 actual=$profileFiles"
  $profileResults.Add([pscustomobject]@{ Profile = $profile; Files = $profileFiles })
  $profile = 'GITLAB_OPENGREP_SIGNED_SAST_PACK_PLAN.md'
  $destination = Join-Path $tempRoot ([IO.Path]::GetFileNameWithoutExtension($profile))
  & $compositor -PlanFile (Join-Path $libraryRoot "markdown_system/$profile") -LibraryRoot $libraryRoot -Destination $destination | Out-Null
  Assert-True (Test-Path -LiteralPath (Join-Path $destination 'MATERIALIZATION_RECORD.md')) "composition record missing for $profile"
  $profileFiles = (Get-ChildItem -LiteralPath $destination -Recurse -File | Measure-Object).Count - 1
  Assert-True ($profileFiles -eq 5) "profile file count mismatch for $profile expected=5 actual=$profileFiles"
  $profileResults.Add([pscustomobject]@{ Profile = $profile; Files = $profileFiles })
  $profile = 'AMAZON_SPAPI_SHIPPING_TRACKING_PACK_PLAN.md'
  $destination = Join-Path $tempRoot ([IO.Path]::GetFileNameWithoutExtension($profile))
  & $compositor -PlanFile (Join-Path $libraryRoot "markdown_system/$profile") -LibraryRoot $libraryRoot -Destination $destination | Out-Null
  Assert-True (Test-Path -LiteralPath (Join-Path $destination 'MATERIALIZATION_RECORD.md')) "composition record missing for $profile"
  $profileFiles = (Get-ChildItem -LiteralPath $destination -Recurse -File | Measure-Object).Count - 1
  Assert-True ($profileFiles -eq 74) "profile file count mismatch for $profile expected=74 actual=$profileFiles"
  $profileResults.Add([pscustomobject]@{ Profile = $profile; Files = $profileFiles })
  $profile = 'AMAZON_SPAPI_EASYSHIP_HANDOVER_PACK_PLAN.md'
  $destination = Join-Path $tempRoot ([IO.Path]::GetFileNameWithoutExtension($profile))
  & $compositor -PlanFile (Join-Path $libraryRoot "markdown_system/$profile") -LibraryRoot $libraryRoot -Destination $destination | Out-Null
  Assert-True (Test-Path -LiteralPath (Join-Path $destination 'MATERIALIZATION_RECORD.md')) "composition record missing for $profile"
  $profileFiles = (Get-ChildItem -LiteralPath $destination -Recurse -File | Measure-Object).Count - 1
  Assert-True ($profileFiles -eq 60) "profile file count mismatch for $profile expected=60 actual=$profileFiles"
  $profileResults.Add([pscustomobject]@{ Profile = $profile; Files = $profileFiles })
  $profile = 'AMAZON_SPAPI_FULFILLMENT_DELIVERY_EVIDENCE_PACK_PLAN.md'
  $destination = Join-Path $tempRoot ([IO.Path]::GetFileNameWithoutExtension($profile))
  & $compositor -PlanFile (Join-Path $libraryRoot "markdown_system/$profile") -LibraryRoot $libraryRoot -Destination $destination | Out-Null
  Assert-True (Test-Path -LiteralPath (Join-Path $destination 'MATERIALIZATION_RECORD.md')) "composition record missing for $profile"
  $profileFiles = (Get-ChildItem -LiteralPath $destination -Recurse -File | Measure-Object).Count - 1
  Assert-True ($profileFiles -eq 59) "profile file count mismatch for $profile expected=59 actual=$profileFiles"
  $profileResults.Add([pscustomobject]@{ Profile = $profile; Files = $profileFiles })
  $profile = 'AMAZON_SPAPI_SUPPLY_SOURCES_PACK_PLAN.md'
  $destination = Join-Path $tempRoot ([IO.Path]::GetFileNameWithoutExtension($profile))
  & $compositor -PlanFile (Join-Path $libraryRoot "markdown_system/$profile") -LibraryRoot $libraryRoot -Destination $destination | Out-Null
  Assert-True (Test-Path -LiteralPath (Join-Path $destination 'MATERIALIZATION_RECORD.md')) "composition record missing for $profile"
  $profileFiles = (Get-ChildItem -LiteralPath $destination -Recurse -File | Measure-Object).Count - 1
  Assert-True ($profileFiles -eq 60) "profile file count mismatch for $profile expected=60 actual=$profileFiles"
  $profileResults.Add([pscustomobject]@{ Profile = $profile; Files = $profileFiles })
  $profile = 'AMAZON_SPAPI_MLI_INVENTORY_PACK_PLAN.md'
  $destination = Join-Path $tempRoot ([IO.Path]::GetFileNameWithoutExtension($profile))
  & $compositor -PlanFile (Join-Path $libraryRoot "markdown_system/$profile") -LibraryRoot $libraryRoot -Destination $destination | Out-Null
  Assert-True (Test-Path -LiteralPath (Join-Path $destination 'MATERIALIZATION_RECORD.md')) "composition record missing for $profile"
  $profileFiles = (Get-ChildItem -LiteralPath $destination -Recurse -File | Measure-Object).Count - 1
  Assert-True ($profileFiles -eq 69) "profile file count mismatch for $profile expected=69 actual=$profileFiles"
  $profileResults.Add([pscustomobject]@{ Profile = $profile; Files = $profileFiles })
  $profile = 'AMAZON_SPAPI_EXTERNAL_FULFILLMENT_INVENTORY_PACK_PLAN.md'
  $destination = Join-Path $tempRoot ([IO.Path]::GetFileNameWithoutExtension($profile))
  & $compositor -PlanFile (Join-Path $libraryRoot "markdown_system/$profile") -LibraryRoot $libraryRoot -Destination $destination | Out-Null
  Assert-True (Test-Path -LiteralPath (Join-Path $destination 'MATERIALIZATION_RECORD.md')) "composition record missing for $profile"
  $profileFiles = (Get-ChildItem -LiteralPath $destination -Recurse -File | Measure-Object).Count - 1
  Assert-True ($profileFiles -eq 60) "profile file count mismatch for $profile expected=60 actual=$profileFiles"
  $profileResults.Add([pscustomobject]@{ Profile = $profile; Files = $profileFiles })
  $profile = 'AZURE_DOCUMENT_INTELLIGENCE_OFFICIAL_INVOICE_SAMPLE_PACK_PLAN.md'
  $destination = Join-Path $tempRoot ([IO.Path]::GetFileNameWithoutExtension($profile))
  & $compositor -PlanFile (Join-Path $libraryRoot "markdown_system/$profile") -LibraryRoot $libraryRoot -Destination $destination | Out-Null
  Assert-True (Test-Path -LiteralPath (Join-Path $destination 'MATERIALIZATION_RECORD.md')) "composition record missing for $profile"
  $profileFiles = (Get-ChildItem -LiteralPath $destination -Recurse -File | Measure-Object).Count - 1
  Assert-True ($profileFiles -eq 14) "profile file count mismatch for $profile expected=14 actual=$profileFiles"
  $profileResults.Add([pscustomobject]@{ Profile = $profile; Files = $profileFiles })
  $profile = 'GOOGLE_DOCUMENT_AI_OFFICIAL_PROCESS_SAMPLE_PACK_PLAN.md'
  $destination = Join-Path $tempRoot ([IO.Path]::GetFileNameWithoutExtension($profile))
  & $compositor -PlanFile (Join-Path $libraryRoot "markdown_system/$profile") -LibraryRoot $libraryRoot -Destination $destination | Out-Null
  Assert-True (Test-Path -LiteralPath (Join-Path $destination 'MATERIALIZATION_RECORD.md')) "composition record missing for $profile"
  $profileFiles = (Get-ChildItem -LiteralPath $destination -Recurse -File | Measure-Object).Count - 1
  Assert-True ($profileFiles -eq 19) "profile file count mismatch for $profile expected=19 actual=$profileFiles"
  $profileResults.Add([pscustomobject]@{ Profile = $profile; Files = $profileFiles })
  $profile = 'GOOGLE_DOCUMENT_AI_OFFICIAL_LIFECYCLE_PACK_PLAN.md'
  $destination = Join-Path $tempRoot ([IO.Path]::GetFileNameWithoutExtension($profile))
  & $compositor -PlanFile (Join-Path $libraryRoot "markdown_system/$profile") -LibraryRoot $libraryRoot -Destination $destination | Out-Null
  Assert-True (Test-Path -LiteralPath (Join-Path $destination 'MATERIALIZATION_RECORD.md')) "composition record missing for $profile"
  $profileFiles = (Get-ChildItem -LiteralPath $destination -Recurse -File | Measure-Object).Count - 1
  Assert-True ($profileFiles -eq 40) "profile file count mismatch for $profile expected=40 actual=$profileFiles"
  $profileResults.Add([pscustomobject]@{ Profile = $profile; Files = $profileFiles })
  $profile = 'DAPR_OFFICIAL_TRANSACTIONAL_OUTBOX_PACK_PLAN.md'
  $destination = Join-Path $tempRoot ([IO.Path]::GetFileNameWithoutExtension($profile))
  & $compositor -PlanFile (Join-Path $libraryRoot "markdown_system/$profile") -LibraryRoot $libraryRoot -Destination $destination | Out-Null
  Assert-True (Test-Path -LiteralPath (Join-Path $destination 'MATERIALIZATION_RECORD.md')) "composition record missing for $profile"
  $profileFiles = (Get-ChildItem -LiteralPath $destination -Recurse -File | Measure-Object).Count - 1
  Assert-True ($profileFiles -eq 21) "profile file count mismatch for $profile expected=21 actual=$profileFiles"
  $profileResults.Add([pscustomobject]@{ Profile = $profile; Files = $profileFiles })
  $profile = 'MICROSOFT_PG_DURABLE_HUMAN_HANDOFF_PACK_PLAN.md'
  $destination = Join-Path $tempRoot ([IO.Path]::GetFileNameWithoutExtension($profile))
  & $compositor -PlanFile (Join-Path $libraryRoot "markdown_system/$profile") -LibraryRoot $libraryRoot -Destination $destination | Out-Null
  Assert-True (Test-Path -LiteralPath (Join-Path $destination 'MATERIALIZATION_RECORD.md')) "composition record missing for $profile"
  $profileFiles = (Get-ChildItem -LiteralPath $destination -Recurse -File | Measure-Object).Count - 1
  Assert-True ($profileFiles -eq 26) "profile file count mismatch for $profile expected=26 actual=$profileFiles"
  $profileResults.Add([pscustomobject]@{ Profile = $profile; Files = $profileFiles })

  $profile = 'HISTORY_MODEL_TRAINING_PACK_PLAN.md'
  $destination = Join-Path $tempRoot ([IO.Path]::GetFileNameWithoutExtension($profile))
  & $compositor -PlanFile (Join-Path $libraryRoot "markdown_system/$profile") -LibraryRoot $libraryRoot -Destination $destination | Out-Null
  Assert-True (Test-Path -LiteralPath (Join-Path $destination 'MATERIALIZATION_RECORD.md')) "composition record missing for $profile"
  $profileFiles = (Get-ChildItem -LiteralPath $destination -Recurse -File | Measure-Object).Count - 1
  Assert-True ($profileFiles -eq 21) "profile file count mismatch for $profile expected=21 actual=$profileFiles"
  $profileResults.Add([pscustomobject]@{ Profile = $profile; Files = $profileFiles })

  $profile = 'HTTP_METRICS_REFERENCE_PACK_PLAN.md'
  $destination = Join-Path $tempRoot ([IO.Path]::GetFileNameWithoutExtension($profile))
  & $compositor -PlanFile (Join-Path $libraryRoot "markdown_system/$profile") -LibraryRoot $libraryRoot -Destination $destination | Out-Null
  Assert-True (Test-Path -LiteralPath (Join-Path $destination 'MATERIALIZATION_RECORD.md')) "composition record missing for $profile"
  $profileFiles = (Get-ChildItem -LiteralPath $destination -Recurse -File | Measure-Object).Count - 1
  Assert-True ($profileFiles -eq (Get-ReleaseInputPolicy).franchise_files) "full HTTP reference differs from the current locked franchise file count: actual=$profileFiles"
  $profileResults.Add([pscustomobject]@{ Profile = $profile; Files = $profileFiles })

  foreach ($result in $profileResults) {
    $profileDestination = Join-Path $tempRoot ([IO.Path]::GetFileNameWithoutExtension($result.Profile))
    Assert-LocalGoImportsResolved $profileDestination $result.Profile
  }

  $allPlanNames = @(Get-ChildItem -LiteralPath (Join-Path $libraryRoot 'markdown_system') -Filter '*PACK_PLAN.md' | ForEach-Object Name | Sort-Object)
  $composedPlanNames = @($profileResults | ForEach-Object Profile | Sort-Object)
  Assert-True ($allPlanNames.Count -eq $composedPlanNames.Count) "not every pack plan was composed: discovered=$($allPlanNames.Count) composed=$($composedPlanNames.Count)"
  foreach ($planName in $allPlanNames) {
    Assert-True (@($composedPlanNames | Where-Object { $_ -eq $planName }).Count -eq 1) "pack plan was not composed exactly once: $planName"
  }

  $markdownFiles = @($releaseFiles | Where-Object Extension -eq '.md')
  $franchiseProfile = @($profileResults | Where-Object Profile -eq 'FRANCHISE_COMPLETE_PACK_PLAN.md')
  Assert-True ($franchiseProfile.Count -eq 1) 'franchise profile result missing or duplicated'
  Assert-True ($franchiseProfile[0].Files -eq (Get-ReleaseInputPolicy).franchise_files) 'franchise reconstructed file count differs from release lock'
  $preflightSummary = Get-Content -LiteralPath (Join-Path $libraryRoot 'markdown_system/FRANCHISE_PREFLIGHT_GAP.md') -Raw
  $normalizedPreflight = [regex]::Replace($preflightSummary, '\s+', ' ')
  $materializedDisplay = $materializedFiles.ToString('N0', [Globalization.CultureInfo]::GetCultureInfo('es-AR'))
  $expectedInventory = "Inventario raíz: $($packFiles.Count) packs, $materializedDisplay archivos materializables, $($markdownFiles.Count) Markdown y $($profileResults.Count) perfiles."
  Assert-True ($normalizedPreflight.Contains($expectedInventory)) "franchise preflight inventory is stale; expected $($packFiles.Count)/$materializedFiles/$($markdownFiles.Count)/$($profileResults.Count)"
  $expectedFranchise = "El perfil integral vigente selecciona $($franchiseProfile[0].Files) archivos"
  Assert-True ($normalizedPreflight.Contains($expectedFranchise)) "franchise preflight composition count is stale; expected=$($franchiseProfile[0].Files)"
  Write-Output 'VERIFY_LIBRARY_PASS'
  Write-Output "packs=$($packFiles.Count) materialized_files=$materializedFiles markdown_files=$($markdownFiles.Count)"
  foreach ($result in $profileResults) { Write-Output "profile=$($result.Profile) implementation_files=$($result.Files)" }
} finally {
  if (Test-Path -LiteralPath $tempRoot) {
    $resolved = [IO.Path]::GetFullPath($tempRoot)
    $systemTemp = [IO.Path]::GetFullPath([IO.Path]::GetTempPath())
    if ($resolved.StartsWith($systemTemp, [StringComparison]::OrdinalIgnoreCase) -and (Split-Path -Leaf $resolved) -like 'elite-library-verification-*') {
      Remove-Item -LiteralPath $resolved -Recurse -Force
    }
  }
}
