#requires -Version 7.0

[CmdletBinding()]
param(
  [Parameter(Mandatory = $true)]
  [string] $PackFile,
  [Parameter(Mandatory = $true)]
  [string] $Destination
)

$ErrorActionPreference = 'Stop'
$utf8 = [Text.UTF8Encoding]::new($false)
$packPath = (Resolve-Path -LiteralPath $PackFile).Path
$destinationPath = [IO.Path]::GetFullPath($Destination)

function Assert-SafeRelativePath([string] $Value, [string] $Label) {
  $normalized = $Value.Replace('\', '/')
  $segments = $normalized.Split('/', [StringSplitOptions]::None)
  if (
    [string]::IsNullOrWhiteSpace($Value) -or
    [IO.Path]::IsPathRooted($Value) -or
    $segments.Count -eq 0 -or
    $segments -contains '' -or
    $segments -contains '.' -or
    $segments -contains '..'
  ) {
    throw "$Label is unsafe: $Value"
  }
  foreach ($segment in $segments) {
    if ($segment.EndsWith('.') -or $segment.EndsWith(' ')) { throw "$Label is non-canonical: $Value" }
  }
  if ($normalized -ieq 'MATERIALIZATION_RECORD.md') {
    throw "$Label is reserved: $Value"
  }
  $normalized
}

function Assert-RequiredPackMetadata([string] $Text) {
  $metadataMatch = [regex]::Match($Text, '(?ms)^```yaml\r?\n(?<body>.*?)\r?\n```')
  if (-not $metadataMatch.Success) { throw 'Pack metadata YAML block is missing' }
  $metadata = $metadataMatch.Groups['body'].Value
  foreach ($field in @(
    'pack_id', 'pack_version', 'status', 'claim', 'stacks', 'compatible_with',
    'incompatible_with', 'license_expression', 'upstream_sources', 'verified_at'
  )) {
    if (-not [regex]::IsMatch($metadata, '(?m)^' + [regex]::Escape($field) + ':')) {
      throw "Pack metadata field is missing: $field"
    }
  }
  foreach ($field in @('authority', 'implementation', 'admission')) {
    if (-not [regex]::IsMatch($metadata, '(?m)^\s{2}' + $field + ':\s*\S+\s*$')) {
      throw "Pack status field is missing: $field"
    }
  }
  $id = [regex]::Match($metadata, '(?m)^pack_id:\s*"(?<value>[^"]+)"\s*$').Groups['value'].Value
  $version = [regex]::Match($metadata, '(?m)^pack_version:\s*"(?<value>[^"]+)"\s*$').Groups['value'].Value
  if (-not $id) { throw 'pack_id must be a non-empty quoted value' }
  if ($version -notmatch '^\d+\.\d+\.\d+$') { throw "pack_version is invalid: $version" }
}

function Read-ExactManifest([string[]] $Lines) {
  $heading = -1
  for ($i = 0; $i -lt $Lines.Length; $i++) {
    if ($Lines[$i] -match '(?i)^## .*manifest.*$') { $heading = $i; break }
    if ($Lines[$i] -match '^### FILE: ') { break }
  }
  if ($heading -lt 0) { throw 'Exact file manifest section is missing' }

  $fence = -1
  for ($i = $heading + 1; $i -lt $Lines.Length; $i++) {
    if ($Lines[$i] -match '^## ') { break }
    if ($Lines[$i] -eq '```text') { $fence = $i; break }
  }
  if ($fence -lt 0) { throw 'Exact file manifest text fence is missing' }

  $paths = [Collections.Generic.List[string]]::new()
  $seen = [Collections.Generic.HashSet[string]]::new([StringComparer]::OrdinalIgnoreCase)
  $closed = $false
  for ($i = $fence + 1; $i -lt $Lines.Length; $i++) {
    if ($Lines[$i] -eq '```') { $closed = $true; break }
    if ([string]::IsNullOrWhiteSpace($Lines[$i])) { continue }
    $match = [regex]::Match($Lines[$i], '^CREATE (?<path>\S+)$')
    if (-not $match.Success) { throw "Manifest permits only exact CREATE entries: $($Lines[$i])" }
    $path = Assert-SafeRelativePath $match.Groups['path'].Value 'Manifest path'
    if (-not $seen.Add($path)) { throw "Duplicate manifest path: $path" }
    $paths.Add($path)
  }
  if (-not $closed) { throw 'Exact file manifest fence is unclosed' }
  if ($paths.Count -eq 0) { throw 'Exact file manifest is empty' }
  $paths
}

if (Test-Path -LiteralPath $destinationPath) {
  if ((Get-ChildItem -LiteralPath $destinationPath -Force | Measure-Object).Count -gt 0) {
    throw "Destination must be empty: $destinationPath"
  }
}

$packText = [IO.File]::ReadAllText($packPath, $utf8)
$lines = [IO.File]::ReadAllLines($packPath, $utf8)
Assert-RequiredPackMetadata $packText
$manifestPaths = @(Read-ExactManifest $lines)
$blocks = [Collections.Generic.List[object]]::new()
$blockIds = [Collections.Generic.HashSet[string]]::new([StringComparer]::Ordinal)
$blockPaths = [Collections.Generic.HashSet[string]]::new([StringComparer]::OrdinalIgnoreCase)

for ($index = 0; $index -lt $lines.Length; $index++) {
  $match = [regex]::Match($lines[$index], '^### FILE: `(?<path>[^`]+)`$')
  if (-not $match.Success) { continue }

  $relative = Assert-SafeRelativePath $match.Groups['path'].Value 'Materialization path'
  if (-not $blockPaths.Add($relative)) { throw "Duplicate materialization path: $relative" }

  $metadataLines = [Collections.Generic.List[string]]::new()
  while ($index -lt $lines.Length -and -not $lines[$index].StartsWith('````')) {
    $metadataLines.Add($lines[$index])
    $index++
  }
  if ($index -ge $lines.Length) { throw "Missing content fence for $relative" }
  $metadata = [string]::Join("`n", $metadataLines)
  foreach ($field in @('block_id', 'operation', 'provenance', 'source', 'license', 'sha256', 'variables', 'secrets_allowed')) {
    if (-not [regex]::IsMatch($metadata, '(?m)^' + [regex]::Escape($field) + ':')) {
      throw "Block metadata field is missing for ${relative}: $field"
    }
  }

  $blockId = [regex]::Match($metadata, '(?m)^block_id:\s*"(?<value>[^"]+)"\s*$').Groups['value'].Value
  $operation = [regex]::Match($metadata, '(?m)^operation:\s*(?<value>[A-Z]+)\s*$').Groups['value'].Value
  $provenance = [regex]::Match($metadata, '(?m)^provenance:\s*(?<value>[A-Z]+)\s*$').Groups['value'].Value
  $declaredHash = [regex]::Match($metadata, '(?m)^sha256:\s*"(?<value>[0-9a-f]{64})"\s*$').Groups['value'].Value
  $secretsAllowed = [regex]::Match($metadata, '(?m)^secrets_allowed:\s*(?<value>true|false)\s*$').Groups['value'].Value
  if (-not $blockId) { throw "block_id is invalid for $relative" }
  if (-not $blockIds.Add($blockId)) { throw "Duplicate block_id: $blockId" }
  if ($operation -ne 'CREATE') { throw "Only CREATE blocks are supported: $relative declares $operation" }
  if ($provenance -notin @('AUTHORED', 'ADAPTED', 'VERBATIM')) { throw "Invalid provenance for ${relative}: $provenance" }
  if (-not $declaredHash) { throw "SHA-256 is invalid for $relative" }
  if (-not $secretsAllowed) { throw "secrets_allowed must be true or false for $relative" }

  $index++
  $contentLines = [Collections.Generic.List[string]]::new()
  while ($index -lt $lines.Length -and $lines[$index] -ne '````') {
    $contentLines.Add($lines[$index])
    $index++
  }
  if ($index -ge $lines.Length) { throw "Unclosed content fence for $relative" }

  $newlineFields = [regex]::Matches($metadata, '(?m)^final_newline:[^\n]*$')
  $finalNewline = $true
  if ($newlineFields.Count -gt 1) { throw "Duplicate final_newline for $relative" }
  if ($newlineFields.Count -eq 1) {
    $newlineValue = [regex]::Match($newlineFields[0].Value, '^final_newline:[ \t]*(true|false)[ \t]*$')
    if (-not $newlineValue.Success) { throw "final_newline must be true or false for $relative" }
    $finalNewline = $newlineValue.Groups[1].Value -ceq 'true'
  }
  $content = [string]::Join("`n", $contentLines)
  if ($finalNewline -and $contentLines.Count -gt 0) { $content += "`n" }
  if (-not $finalNewline -and $content.EndsWith("`n")) { throw "final_newline false conflicts with trailing empty content line for $relative" }
  $bytes = $utf8.GetBytes($content)
  $actualHash = [Convert]::ToHexString([Security.Cryptography.SHA256]::HashData($bytes)).ToLowerInvariant()
  if ($actualHash -ne $declaredHash) {
    throw "SHA-256 mismatch for ${relative}: expected $declaredHash, got $actualHash"
  }
  $blocks.Add([pscustomobject]@{ Path = $relative; Bytes = $bytes })
}

if ($blocks.Count -eq 0) { throw 'No materialization blocks found' }
if ($manifestPaths.Count -ne $blocks.Count) {
  throw "Manifest/FILE count mismatch: manifest=$($manifestPaths.Count), blocks=$($blocks.Count)"
}
foreach ($manifestPath in $manifestPaths) {
  if (-not @($blocks | Where-Object { $_.Path -ceq $manifestPath })) {
    throw "Manifest path has no exact FILE block: $manifestPath"
  }
}

foreach ($block in $blocks) {
  $target = [IO.Path]::GetFullPath((Join-Path $destinationPath $block.Path))
  if (-not $target.StartsWith($destinationPath + [IO.Path]::DirectorySeparatorChar, [StringComparison]::OrdinalIgnoreCase)) {
    throw "Materialization escaped destination: $($block.Path)"
  }
}

if (-not (Test-Path -LiteralPath $destinationPath)) {
  [void](New-Item -ItemType Directory -Path $destinationPath)
}
foreach ($block in $blocks) {
  $target = [IO.Path]::GetFullPath((Join-Path $destinationPath $block.Path))
  $parent = Split-Path -Parent $target
  if (-not (Test-Path -LiteralPath $parent)) { [void](New-Item -ItemType Directory -Path $parent) }
  [IO.File]::WriteAllBytes($target, $block.Bytes)
}

Write-Output "Materialized $($blocks.Count) files into $destinationPath"
