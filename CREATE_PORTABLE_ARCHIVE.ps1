#requires -Version 7.0

[CmdletBinding()]
param(
  [Parameter(Mandatory = $true)]
  [string] $OutputPath,
  [ValidateRange(315532800, 4354819198)]
  [long] $SourceDateEpoch = 946684800
)

$ErrorActionPreference = 'Stop'
$libraryRoot = [IO.Path]::GetFullPath($PSScriptRoot).TrimEnd([IO.Path]::DirectorySeparatorChar)
$archivePath = [IO.Path]::GetFullPath($OutputPath)
$archiveParent = Split-Path -Parent $archivePath
$archiveName = Split-Path -Leaf $archivePath
$sidecar = $archivePath + '.sha256'
$tempRoot = Join-Path ([IO.Path]::GetTempPath()) ('elite-library-archive-' + [guid]::NewGuid().ToString('N'))
$payload = Join-Path $tempRoot 'Elite Engineering Library'
$candidateArchive = Join-Path $tempRoot 'candidate.zip'
$utf8 = [Text.UTF8Encoding]::new($false)
$strictUtf8 = [Text.UTF8Encoding]::new($false, $true)
$pathComparison = if ($IsWindows) { [StringComparison]::OrdinalIgnoreCase } else { [StringComparison]::Ordinal }
$archivePrefix = 'Elite Engineering Library/'
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
  'PROJECT_FAILURE_LESSONS.md',
  'PUBLIC_CODE_ARCHITECTURE_ADMISSION_STANDARD.md',
  'PUBLIC_CODE_ARCHITECTURE_MASTER_MAP.md',
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

function Relative-ToLibrary([string] $Path) {
  [IO.Path]::GetRelativePath($libraryRoot, $Path).Replace('\', '/')
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
  $localStateNames = @('PROJECT_EXECUTION_STATE.json', 'PROJECT_EXECUTION_EVENTS.jsonl')
  $localStateCount = @($localStateNames | Where-Object { Test-Path -LiteralPath (Join-Path $libraryRoot $_) }).Count
  if ($localStateCount -notin @(0, 2)) { throw 'Local execution state and events must exist together' }
  foreach ($item in Get-ChildItem -LiteralPath $libraryRoot -Force) {
    if ($item.Name -ieq '.git') { continue }
    if (($item.Attributes -band [IO.FileAttributes]::ReparsePoint) -ne 0) {
      throw "Release entry is a symlink/reparse point: $($item.Name)"
    }

    if ($localStateNames -ccontains $item.Name) {
      if ($item.PSIsContainer) { throw "Local execution artifact must be a regular file: $($item.Name)" }
      continue
    }
    if (Test-LocalMaintenanceEntry $item $localStateCount) { continue }

    if ($item.PSIsContainer) {
      if (-not $releaseDirectories.ContainsKey($item.Name)) { throw "Unknown top-level directory: $($item.Name)" }
      $allowedExtensions = $releaseDirectories[$item.Name]
      foreach ($child in Get-ChildItem -LiteralPath $item.FullName -Recurse -Force) {
        if (($child.Attributes -band [IO.FileAttributes]::ReparsePoint) -ne 0) {
          throw "Release entry is a symlink/reparse point: $(Relative-ToLibrary $child.FullName)"
        }
        if ($child.PSIsContainer) { continue }
        $publicEvidence = $null -ne $releasePolicy -and $item.Name -ceq 'reconstruction_evidence' -and
          $child.DirectoryName -ceq $item.FullName -and (@($releasePolicy.json_files) + @($releasePolicy.text_logs)) -ccontains $child.Name
        if ($allowedExtensions -inotcontains $child.Extension -and -not $publicEvidence) { throw "Unknown release file: $(Relative-ToLibrary $child.FullName)" }
        if ($item.Name -eq 'markdown_system' -and $child.Extension -ieq '.ps1' -and $releaseMarkdownSystemScripts -inotcontains $child.Name) {
          throw "Unknown markdown_system script: $(Relative-ToLibrary $child.FullName)"
        }
        $files.Add($child)
      }
      continue
    }

    if ($releaseRootMarkdown -inotcontains $item.Name -and $releaseRootScripts -inotcontains $item.Name -and $releaseRootMetadata -cnotcontains $item.Name) {
      throw "Unknown top-level file: $($item.Name)"
    }
    $files.Add($item)
  }

  if ($files.Count -eq 0) { throw 'Release file allowlist produced no files' }
  if ($null -ne $releasePolicy) {
    foreach ($name in (@($releasePolicy.json_files) + @($releasePolicy.text_logs))) {
      if (@($files | Where-Object { $_.DirectoryName -ceq (Join-Path $libraryRoot 'reconstruction_evidence') -and $_.Name -ceq $name }).Count -ne 1) { throw "Release policy named evidence missing: $name" }
    }
  }
  $relativePaths = @($files | ForEach-Object { Relative-ToLibrary $_.FullName })
  if (($relativePaths | Sort-Object -Unique).Count -ne $relativePaths.Count) { throw 'Duplicate release path' }
  $files | Sort-Object FullName
}

function Read-ZipText([IO.Compression.ZipArchiveEntry] $Entry) {
  $stream = $Entry.Open()
  $reader = [IO.StreamReader]::new($stream, $strictUtf8, $false)
  try {
    $reader.ReadToEnd()
  } finally {
    $reader.Dispose()
  }
}

function Get-ZipEntrySha256([IO.Compression.ZipArchiveEntry] $Entry) {
  $stream = $Entry.Open()
  $sha = [Security.Cryptography.SHA256]::Create()
  try {
    [Convert]::ToHexString($sha.ComputeHash($stream)).ToLowerInvariant()
  } finally {
    $sha.Dispose()
    $stream.Dispose()
  }
}

# AUTHORED adaptation of the existing PORTABLE-SIGNED-RELEASE-EVIDENCE-GATE
# archive pattern. This helper does not sign or admit the resulting artifact.
function Write-DeterministicArchive([string] $InputRoot, [string] $Candidate, [long] $Epoch) {
  if ($Epoch -lt 315532800 -or $Epoch -gt 4354819198 -or $Epoch % 2 -ne 0) { throw 'SourceDateEpoch must be an even ZIP-representable UTC second (1980-2107)' }
  $stamp = [DateTimeOffset]::FromUnixTimeSeconds($Epoch)
  $root = [IO.Path]::GetFullPath($InputRoot).TrimEnd([IO.Path]::DirectorySeparatorChar)
  $items = @(Get-ChildItem -LiteralPath $root -Recurse -Force)
  if (@($items | Where-Object { ($_.Attributes -band [IO.FileAttributes]::ReparsePoint) -ne 0 }).Count) { throw 'Archive payload contains a reparse point' }
  [string[]]$paths = @($items | Where-Object { -not $_.PSIsContainer } | ForEach-Object FullName)
  if ($paths.Count -eq 0) { throw 'Archive payload is empty' }
  [Array]::Sort($paths, [StringComparer]::Ordinal)
  $stream = [IO.File]::Open($Candidate, [IO.FileMode]::CreateNew, [IO.FileAccess]::ReadWrite, [IO.FileShare]::None)
  try {
    $zip = [IO.Compression.ZipArchive]::new($stream, [IO.Compression.ZipArchiveMode]::Create, $true, [Text.UTF8Encoding]::new($false))
    try {
      foreach ($path in $paths) {
        $relative = [IO.Path]::GetRelativePath($root, $path).Replace('\','/')
        if ([IO.Path]::IsPathRooted($relative) -or ($relative -split '/') -contains '..') { throw 'Archive entry escaped payload' }
        $entry = $zip.CreateEntry(('Elite Engineering Library/' + $relative), [IO.Compression.CompressionLevel]::NoCompression)
        $entry.LastWriteTime = $stamp
        $entry.ExternalAttributes = 0
        $inputFile = [IO.File]::OpenRead($path)
        try {
          $outputFile = $entry.Open()
          try { $inputFile.CopyTo($outputFile) } finally { $outputFile.Dispose() }
        } finally { $inputFile.Dispose() }
      }
    } finally { $zip.Dispose() }
    $stream.Flush($true)
  } finally { $stream.Dispose() }
}

# AUTHORED: exclusive publication of a verified candidate and its checksum.
# Cooperating publishers never replace an existing path. Normal exceptions remove
# only outputs reserved by this invocation. A process/OS crash may leave a partial
# pair: consumers must verify both files; no power-loss or hostile-directory claim.
function Publish-VerifiedArchive([string] $Candidate, [string] $Destination) {
  $checksumPath = $Destination + '.sha256'
  $name = [IO.Path]::GetFileName($Destination)
  $inputStream = $null
  $archiveStream = $null
  $checksumStream = $null
  $archiveOwned = $false
  $checksumOwned = $false
  $complete = $false
  try {
    $inputStream = [IO.File]::Open($Candidate, [IO.FileMode]::Open, [IO.FileAccess]::Read, [IO.FileShare]::Read)
    $hash = [Convert]::ToHexString([Security.Cryptography.SHA256]::HashData($inputStream)).ToLowerInvariant()
    $inputStream.Position = 0
    $checksumBytes = [Text.UTF8Encoding]::new($false).GetBytes("$hash  $name`n")
    # Reserve the checksum first; a late competing checksum must stay untouched.
    $checksumStream = [IO.File]::Open($checksumPath, [IO.FileMode]::CreateNew, [IO.FileAccess]::Write, [IO.FileShare]::None)
    $checksumOwned = $true
    $archiveStream = [IO.File]::Open($Destination, [IO.FileMode]::CreateNew, [IO.FileAccess]::Write, [IO.FileShare]::None)
    $archiveOwned = $true
    $inputStream.CopyTo($archiveStream)
    $archiveStream.Flush($true)
    $checksumStream.Write($checksumBytes, 0, $checksumBytes.Length)
    $checksumStream.Flush($true)
    $complete = $true
    return $hash
  } finally {
    try {
      if ($archiveStream) { $archiveStream.Dispose() }
    } finally {
      try {
        if ($checksumStream) { $checksumStream.Dispose() }
      } finally {
        if ($inputStream) { $inputStream.Dispose() }
        if (-not $complete) {
          # No recursive cleanup or deletion of another publisher's files.
          if ($archiveOwned) { [IO.File]::Delete($Destination) }
          if ($checksumOwned) { [IO.File]::Delete($checksumPath) }
        }
      }
    }
  }
}

if ($SourceDateEpoch % 2 -ne 0) { throw 'SourceDateEpoch must be an even UTC second' }
if ([IO.Path]::GetExtension($archivePath) -ine '.zip') { throw 'OutputPath must end in .zip' }
if ($archivePath.StartsWith($libraryRoot + [IO.Path]::DirectorySeparatorChar, $pathComparison)) { throw 'OutputPath must be outside the library root' }
if (Test-Path -LiteralPath $archivePath) { throw "Archive already exists: $archivePath" }
if (Test-Path -LiteralPath $sidecar) { throw "Archive checksum already exists: $sidecar" }
if (-not (Test-Path -LiteralPath $archiveParent -PathType Container)) { throw "Archive parent does not exist: $archiveParent" }

& (Join-Path $libraryRoot 'VERIFY_LIBRARY.ps1') | Write-Output

try {
  $releaseFiles = @(Get-ReleaseSourceFiles)
  [void](New-Item -ItemType Directory -Path $payload)
  foreach ($file in $releaseFiles) {
    $relative = Relative-ToLibrary $file.FullName
    $target = [IO.Path]::GetFullPath((Join-Path $payload $relative))
    if (-not $target.StartsWith($payload + [IO.Path]::DirectorySeparatorChar, $pathComparison)) {
      throw "Release path escaped payload: $relative"
    }
    $parent = Split-Path -Parent $target
    if (-not (Test-Path -LiteralPath $parent)) { [void](New-Item -ItemType Directory -Path $parent) }
    Copy-Item -LiteralPath $file.FullName -Destination $target
  }

  $manifestPath = Join-Path $payload 'DISTRIBUTION_SHA256SUMS.txt'
  [string[]]$manifestLines = @(Get-ChildItem -LiteralPath $payload -Recurse -File | Where-Object FullName -ne $manifestPath | ForEach-Object {
    $relative = [IO.Path]::GetRelativePath($payload, $_.FullName).Replace('\','/')
    $hash = (Get-FileHash -Algorithm SHA256 -LiteralPath $_.FullName).Hash.ToLowerInvariant()
    "$hash  $relative"
  })
  [Array]::Sort($manifestLines, [StringComparer]::Ordinal)
  $expectedManifest = [string]::Join("`n", $manifestLines) + "`n"
  [IO.File]::WriteAllText($manifestPath, $expectedManifest, $utf8)

  Write-DeterministicArchive $payload $candidateArchive $SourceDateEpoch
  Add-Type -AssemblyName System.IO.Compression.FileSystem
  $zip = [IO.Compression.ZipFile]::OpenRead($candidateArchive)
  try {
    $entries = @($zip.Entries | Where-Object { $_.FullName -ne '' -and -not $_.FullName.EndsWith('/') })
    $manifestEntryName = $archivePrefix + 'DISTRIBUTION_SHA256SUMS.txt'
    $manifestEntries = @($entries | Where-Object FullName -eq $manifestEntryName)
    if ($manifestEntries.Count -ne 1) { throw 'Archive manifest is missing or duplicated' }

    $archiveManifest = Read-ZipText $manifestEntries[0]
    if ($archiveManifest -cne $expectedManifest) { throw 'Archive manifest differs from the staged manifest' }
    $manifest = [Collections.Generic.Dictionary[string,string]]::new([StringComparer]::OrdinalIgnoreCase)
    foreach ($line in $archiveManifest.Split("`n", [StringSplitOptions]::RemoveEmptyEntries)) {
      $match = [regex]::Match($line, '^(?<hash>[0-9a-f]{64})  (?<path>.+)$')
      if (-not $match.Success) { throw "Invalid archive manifest line: $line" }
      $relative = $match.Groups['path'].Value
      if ($relative.Contains('\') -or [IO.Path]::IsPathRooted($relative) -or ($relative -split '/') -contains '..') {
        throw "Unsafe archive manifest path: $relative"
      }
      if ($manifest.ContainsKey($relative)) { throw "Duplicate archive manifest path: $relative" }
      $manifest.Add($relative, $match.Groups['hash'].Value)
    }

    $seen = [Collections.Generic.HashSet[string]]::new([StringComparer]::OrdinalIgnoreCase)
    foreach ($entry in $entries) {
      if ($entry.FullName.Contains('\') -or -not $entry.FullName.StartsWith($archivePrefix, [StringComparison]::Ordinal)) {
        throw "Archive contains an unsafe entry: $($entry.FullName)"
      }
      if ($entry.FullName -eq $manifestEntryName) { continue }
      $relative = $entry.FullName.Substring($archivePrefix.Length)
      if ($relative -eq '' -or ($relative -split '/') -contains '..' -or [IO.Path]::IsPathRooted($relative)) {
        throw "Archive contains an unsafe entry: $($entry.FullName)"
      }
      if (-not $seen.Add($relative)) { throw "Archive contains a duplicate path: $relative" }
      if (-not $manifest.ContainsKey($relative)) { throw "Archive entry is absent from manifest: $relative" }
      $actualHash = Get-ZipEntrySha256 $entry
      if ($actualHash -cne $manifest[$relative]) { throw "Archive entry SHA-256 mismatch: $relative" }
    }
    if ($seen.Count -ne $manifest.Count) { throw "Archive manifest/entry count mismatch: manifest=$($manifest.Count), entries=$($seen.Count)" }
  } finally {
    $zip.Dispose()
  }

  $archiveHash = Publish-VerifiedArchive $candidateArchive $archivePath
  Write-Output "ARCHIVE_READY path=$archivePath"
  Write-Output "ARCHIVE_SHA256=$archiveHash"
  Write-Output "ARCHIVE_FILES=$($manifestLines.Count + 1)"
  Write-Output "ARCHIVE_SOURCE_DATE_EPOCH=$SourceDateEpoch"
} finally {
  if (Test-Path -LiteralPath $tempRoot) {
    $resolved = [IO.Path]::GetFullPath($tempRoot)
    $systemTemp = [IO.Path]::GetFullPath([IO.Path]::GetTempPath())
    if ($resolved.StartsWith($systemTemp, $pathComparison) -and (Split-Path -Leaf $resolved) -like 'elite-library-archive-*') {
      Remove-Item -LiteralPath $resolved -Recurse -Force
    }
  }
}
