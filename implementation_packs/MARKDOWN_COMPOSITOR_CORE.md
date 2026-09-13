# Markdown Project Compositor Core

## 1. Metadata

```yaml
pack_id: "MARKDOWN-COMPOSITOR"
pack_version: "0.3.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Compositor stack-neutral que valida manifests y bloques exactos, preflighta selecciones y materializa múltiples implementation packs con hashes, unicidad, variables y registro trazable."
stacks: ["PowerShell 7+"]
compatible_with: ["PACK_CONTRACT 1.0 + optional final_newline", "COMPOSITION_PROTOCOL 1.0"]
incompatible_with: ["packs SPEC_ONLY/SNIPPET", "admisión DISCOVERED/LICENSE_VERIFIED/CANDIDATE para composición productiva", "operaciones PATCH/DELETE sin merge explícito"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://learn.chatgpt.com/docs/agent-configuration/agents-md"]
verified_at: "2026-09-07"
```

## 2. Applicability

- Usar después de crear el blueprint y el plan de packs.
- Materializa cualquier stack; PowerShell es sólo el tooling de composición, no la base del producto.
- Exige `acknowledgeConditions: true` para packs `CONDITIONED`.
- Rechaza packs candidatos o inferiores: pueden materializarse aisladamente con el bootstrap simple para completar su admisión, pero no entrar silenciosamente en un proyecto.
- Esta versión admite únicamente bloques `CREATE`; `PATCH`/`DELETE` requieren un merge plan posterior y fallan cerrado.

## 3. Architecture contract

```text
PROJECT_PACK_PLAN.md/JSON + library root
→ parse y validar plan
→ parse metadata y bloques de cada pack
→ comprobar estado, versión, paths, hashes y variables
→ detectar colisiones antes de escribir
→ materializar en target nuevo/vacío
→ MATERIALIZATION_RECORD.md con lineage y hashes finales
```

Invariantes: no paths absolutos/`..`; `MATERIALIZATION_RECORD.md` queda reservado; manifest y bloques coinciden exactamente como conjuntos de paths únicos; pack IDs, block IDs y output paths son únicos; no secrets en variables; no contenido con hash inválido; ninguna escritura antes de terminar preflight; ningún pack condicionado sin aceptación explícita; el registro no contiene valores de variables.

## 4. Exact file manifest

```text
CREATE tools/compose-markdown-project.ps1
CREATE tools/test-compose-markdown-project.ps1
CREATE examples/PROJECT_PACK_PLAN.example.md
```

## 5. Materialization blocks

### FILE: `tools/compose-markdown-project.ps1`

```yaml
block_id: "MARKDOWN-COMPOSITOR:compose:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "b9c9473f6bba428a9f1defa0a2cdb4480680c6934bff8adc2df8f328aa36d2f7"
variables: []
secrets_allowed: false
```

````powershell
#requires -Version 7.0

[CmdletBinding()]
param(
  [Parameter(Mandatory = $true)][string] $PlanFile,
  [Parameter(Mandatory = $true)][string] $LibraryRoot,
  [Parameter(Mandatory = $true)][string] $Destination
)

$ErrorActionPreference = 'Stop'
$utf8 = [Text.UTF8Encoding]::new($false)

function Get-Sha256([byte[]] $Bytes) {
  [Convert]::ToHexString([Security.Cryptography.SHA256]::HashData($Bytes)).ToLowerInvariant()
}

function Normalize-RelativePath([string] $Value, [string] $Label, [bool] $ReserveRecord = $false) {
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
  if ($ReserveRecord -and $normalized -ieq 'MATERIALIZATION_RECORD.md') {
    throw "$Label is reserved: $Value"
  }
  $normalized
}

function Resolve-ContainedPath([string] $Root, [string] $Relative, [string] $Label) {
  $normalized = Normalize-RelativePath $Relative $Label
  $rootFull = [IO.Path]::GetFullPath($Root).TrimEnd([IO.Path]::DirectorySeparatorChar, [IO.Path]::AltDirectorySeparatorChar)
  $candidate = [IO.Path]::GetFullPath((Join-Path $rootFull $normalized))
  if (-not $candidate.StartsWith($rootFull + [IO.Path]::DirectorySeparatorChar, [StringComparison]::OrdinalIgnoreCase)) {
    throw "$Label escaped its root: $Relative"
  }
  $candidate
}

function Read-Plan([string] $Path) {
  $text = [IO.File]::ReadAllText((Resolve-Path -LiteralPath $Path), $utf8)
  if ([IO.Path]::GetExtension($Path) -ieq '.md') {
    $matches = [regex]::Matches($text, '(?s)```json\s*(?<json>.*?)\s*```')
    if ($matches.Count -ne 1) { throw 'Markdown plan must contain exactly one fenced json object' }
    $text = $matches[0].Groups['json'].Value
  }
  $plan = $text | ConvertFrom-Json -Depth 100
  if ($plan.compositionVersion -ne '1.0') { throw 'compositionVersion must be 1.0' }
  if ($null -eq $plan.packs -or @($plan.packs).Count -eq 0) { throw 'plan requires at least one pack' }
  if ($null -ne $plan.secretVariables -and @($plan.secretVariables.PSObject.Properties).Count -gt 0) {
    throw 'secrets are forbidden in the composition plan; inject them from the target environment'
  }
  $plan
}

function Assert-PackMetadata([string] $Text, [string] $Path) {
  $match = [regex]::Match($Text, '(?ms)^```yaml\r?\n(?<body>.*?)\r?\n```')
  if (-not $match.Success) { throw "pack metadata YAML is missing: $Path" }
  $metadata = $match.Groups['body'].Value
  foreach ($field in @(
    'pack_id', 'pack_version', 'status', 'claim', 'stacks', 'compatible_with',
    'incompatible_with', 'license_expression', 'upstream_sources', 'verified_at'
  )) {
    if (-not [regex]::IsMatch($metadata, '(?m)^' + [regex]::Escape($field) + ':')) {
      throw "pack metadata field $field is missing: $Path"
    }
  }
  foreach ($field in @('authority', 'implementation', 'admission')) {
    if (-not [regex]::IsMatch($metadata, '(?m)^\s{2}' + $field + ':\s*\S+\s*$')) {
      throw "pack status field $field is missing: $Path"
    }
  }
  $metadata
}

function Read-Manifest([string[]] $Lines, [string] $Path) {
  $heading = -1
  for ($i = 0; $i -lt $Lines.Length; $i++) {
    if ($Lines[$i] -match '(?i)^## .*manifest.*$') { $heading = $i; break }
    if ($Lines[$i] -match '^### FILE: ') { break }
  }
  if ($heading -lt 0) { throw "exact file manifest is missing: $Path" }
  $fence = -1
  for ($i = $heading + 1; $i -lt $Lines.Length; $i++) {
    if ($Lines[$i] -match '^## ') { break }
    if ($Lines[$i] -eq '```text') { $fence = $i; break }
  }
  if ($fence -lt 0) { throw "manifest text fence is missing: $Path" }

  $paths = [Collections.Generic.List[string]]::new()
  $seen = [Collections.Generic.HashSet[string]]::new([StringComparer]::OrdinalIgnoreCase)
  $closed = $false
  for ($i = $fence + 1; $i -lt $Lines.Length; $i++) {
    if ($Lines[$i] -eq '```') { $closed = $true; break }
    if ([string]::IsNullOrWhiteSpace($Lines[$i])) { continue }
    $match = [regex]::Match($Lines[$i], '^CREATE (?<path>\S+)$')
    if (-not $match.Success) { throw "manifest permits only exact CREATE entries: $($Lines[$i])" }
    $relative = Normalize-RelativePath $match.Groups['path'].Value 'manifest path' $true
    if (-not $seen.Add($relative)) { throw "duplicate manifest path: $relative" }
    $paths.Add($relative)
  }
  if (-not $closed) { throw "manifest fence is unclosed: $Path" }
  if ($paths.Count -eq 0) { throw "manifest is empty: $Path" }
  $paths
}

function Read-Pack([string] $Path) {
  $bytes = [IO.File]::ReadAllBytes($Path)
  $lines = [IO.File]::ReadAllLines($Path, $utf8)
  $text = $utf8.GetString($bytes)
  $metadata = Assert-PackMetadata $text $Path
  $id = [regex]::Match($metadata, '(?m)^pack_id:\s*"(?<v>[^"]+)"\s*$').Groups['v'].Value
  $version = [regex]::Match($metadata, '(?m)^pack_version:\s*"(?<v>[^"]+)"\s*$').Groups['v'].Value
  $implementation = [regex]::Match($metadata, '(?m)^\s{2}implementation:\s*(?<v>[A-Z_]+)\s*$').Groups['v'].Value
  $admission = [regex]::Match($metadata, '(?m)^\s{2}admission:\s*(?<v>[A-Z_]+)\s*$').Groups['v'].Value
  if (-not $id -or $version -notmatch '^\d+\.\d+\.\d+$') { throw "pack identity/version is invalid: $Path" }
  $manifest = @(Read-Manifest $lines $Path)

  $blocks = [Collections.Generic.List[object]]::new()
  $blockIds = [Collections.Generic.HashSet[string]]::new([StringComparer]::Ordinal)
  $blockPaths = [Collections.Generic.HashSet[string]]::new([StringComparer]::OrdinalIgnoreCase)
  for ($index = 0; $index -lt $lines.Length; $index++) {
    $fileMatch = [regex]::Match($lines[$index], '^### FILE: `(?<path>[^`]+)`$')
    if (-not $fileMatch.Success) { continue }
    $relative = Normalize-RelativePath $fileMatch.Groups['path'].Value 'block path' $true
    if (-not $blockPaths.Add($relative)) { throw "duplicate block path: $relative" }

    $metadataLines = [Collections.Generic.List[string]]::new()
    while ($index -lt $lines.Length -and -not $lines[$index].StartsWith('````')) {
      $metadataLines.Add($lines[$index])
      $index++
    }
    if ($index -ge $lines.Length) { throw "missing content fence for $relative" }
    $blockMetadata = [string]::Join("`n", $metadataLines)
    foreach ($field in @('block_id', 'operation', 'provenance', 'source', 'license', 'sha256', 'variables', 'secrets_allowed')) {
      if (-not [regex]::IsMatch($blockMetadata, '(?m)^' + [regex]::Escape($field) + ':')) {
        throw "block metadata field $field is missing: $relative"
      }
    }
    $blockId = [regex]::Match($blockMetadata, '(?m)^block_id:\s*"(?<v>[^"]+)"\s*$').Groups['v'].Value
    $operation = [regex]::Match($blockMetadata, '(?m)^operation:\s*(?<v>[A-Z]+)\s*$').Groups['v'].Value
    $provenance = [regex]::Match($blockMetadata, '(?m)^provenance:\s*(?<v>[A-Z]+)\s*$').Groups['v'].Value
    $declaredHash = [regex]::Match($blockMetadata, '(?m)^sha256:\s*"(?<v>[0-9a-f]{64})"\s*$').Groups['v'].Value
    $secrets = [regex]::Match($blockMetadata, '(?m)^secrets_allowed:\s*(?<v>true|false)\s*$').Groups['v'].Value
    $variablesMatch = [regex]::Match($blockMetadata, '(?m)^variables:\s*\[(?<v>.*)\]\s*$')
    if (-not $blockId -or -not $blockIds.Add($blockId)) { throw "duplicate or invalid block_id: $blockId" }
    if ($operation -ne 'CREATE') { throw "only CREATE blocks are supported: $relative declares $operation" }
    if ($provenance -notin @('AUTHORED', 'ADAPTED', 'VERBATIM')) { throw "invalid provenance for ${relative}: $provenance" }
    if (-not $declaredHash) { throw "invalid SHA-256 for $relative" }
    if (-not $secrets) { throw "secrets_allowed must be true or false: $relative" }
    if (-not $variablesMatch.Success) { throw "variables must be an inline list: $relative" }
    $variables = @()
    if ($variablesMatch.Groups['v'].Value.Trim()) {
      $variables = @($variablesMatch.Groups['v'].Value.Split(',') | ForEach-Object { $_.Trim().Trim('"', "'") })
      if (@($variables | Sort-Object -Unique).Count -ne $variables.Count) { throw "duplicate block variable: $relative" }
    }

    $index++
    $contentLines = [Collections.Generic.List[string]]::new()
    while ($index -lt $lines.Length -and $lines[$index] -ne '````') { $contentLines.Add($lines[$index]); $index++ }
    if ($index -ge $lines.Length) { throw "unclosed content fence for $relative" }
    $newlineFields = [regex]::Matches($blockMetadata, '(?m)^final_newline:[^\n]*$')
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
    $actualHash = Get-Sha256 $utf8.GetBytes($content)
    if ($actualHash -ne $declaredHash) { throw "SHA-256 mismatch for ${relative}: expected $declaredHash, got $actualHash" }
    $blocks.Add([pscustomobject]@{
      Path = $relative; BlockId = $blockId; Variables = $variables
      SecretsAllowed = ($secrets -eq 'true'); SourceHash = $actualHash; Content = $content
    })
  }
  if ($blocks.Count -eq 0) { throw "pack has no materialization blocks: $Path" }
  if ($manifest.Count -ne $blocks.Count) { throw "manifest/FILE count mismatch in ${Path}: manifest=$($manifest.Count), blocks=$($blocks.Count)" }
  foreach ($manifestPath in $manifest) {
    if (-not @($blocks | Where-Object { $_.Path -ceq $manifestPath })) {
      throw "manifest path has no exact FILE block in ${Path}: $manifestPath"
    }
  }
  [pscustomobject]@{
    Id = $id; Version = $version; Implementation = $implementation; Admission = $admission
    Sha256 = (Get-Sha256 $bytes); Blocks = $blocks
  }
}

$planPath = (Resolve-Path -LiteralPath $PlanFile).Path
$libraryPath = (Resolve-Path -LiteralPath $LibraryRoot).Path
$destinationPath = [IO.Path]::GetFullPath($Destination)
$plan = Read-Plan $planPath
$outputs = [ordered]@{}
$outputPaths = [Collections.Generic.HashSet[string]]::new([StringComparer]::OrdinalIgnoreCase)
$selectedPackIds = [Collections.Generic.HashSet[string]]::new([StringComparer]::OrdinalIgnoreCase)
$selectedBlockIds = [Collections.Generic.HashSet[string]]::new([StringComparer]::Ordinal)
$packRecords = [Collections.Generic.List[object]]::new()

foreach ($selection in @($plan.packs)) {
  if (-not $selection.path -or -not $selection.packId -or -not $selection.version) { throw 'each selection requires path, packId and version' }
  if (-not $selectedPackIds.Add([string]$selection.packId)) { throw "duplicate packId selection: $($selection.packId)" }
  $packPath = Resolve-ContainedPath $libraryPath ([string]$selection.path) 'pack path'
  if (-not (Test-Path -LiteralPath $packPath -PathType Leaf)) { throw "pack not found: $($selection.path)" }
  $pack = Read-Pack $packPath
  if ($pack.Id -ne $selection.packId -or $pack.Version -ne $selection.version) { throw "pack identity/version mismatch for $($selection.path)" }
  if ($pack.Implementation -notin @('RECONSTRUCTIBLE', 'REBUILD_VERIFIED')) { throw "pack is not materializable: $($pack.Id)" }
  if ($pack.Admission -eq 'CONDITIONED') {
    if ($selection.acknowledgeConditions -ne $true) { throw "conditioned pack requires acknowledgeConditions=true: $($pack.Id)" }
  } elseif ($pack.Admission -ne 'REUSABLE_PACK') {
    throw "pack admission does not permit project composition: $($pack.Id) is $($pack.Admission)"
  }

  [string[]] $selectedFiles = if ($null -eq $selection.files -or @($selection.files).Count -eq 0) { @('*') } else { @($selection.files) }
  if ($selectedFiles -contains '*' -and $selectedFiles.Count -ne 1) { throw "files '*' cannot be combined with explicit paths: $($pack.Id)" }
  $selectedSet = [Collections.Generic.HashSet[string]]::new([StringComparer]::Ordinal)
  if ($selectedFiles -notcontains '*') {
    foreach ($requested in $selectedFiles) {
      $requestedPath = Normalize-RelativePath ([string]$requested) 'selected file' $true
      if (-not $selectedSet.Add($requestedPath)) { throw "duplicate selected file for $($pack.Id): $requestedPath" }
      if (-not @($pack.Blocks | Where-Object { $_.Path -ceq $requestedPath })) {
        throw "selected file does not exist in $($pack.Id): $requestedPath"
      }
    }
  }
  $variables = if ($null -eq $selection.variables) { [pscustomobject]@{} } else { $selection.variables }
  $requiredVariables = [Collections.Generic.HashSet[string]]::new([StringComparer]::Ordinal)

  foreach ($block in $pack.Blocks) {
    if ($selectedFiles -notcontains '*' -and -not $selectedSet.Contains($block.Path)) { continue }
    if (-not $selectedBlockIds.Add($block.BlockId)) { throw "duplicate selected block_id: $($block.BlockId)" }
    foreach ($name in $block.Variables) { [void]$requiredVariables.Add($name) }
    $content = $block.Content
    foreach ($name in $block.Variables) {
      $property = $variables.PSObject.Properties[$name]
      if ($null -eq $property) { throw "missing variable $name for $($pack.Id):$($block.Path)" }
      if ($block.SecretsAllowed -ne $true -and $name -match '(?i)secret|password|token|private[_-]?key') {
        throw "secret-like variable is forbidden in block $($block.BlockId)"
      }
      $content = $content.Replace("{{$name}}", [string]$property.Value)
    }
    if ([regex]::IsMatch($content, '\{\{[A-Z][A-Z0-9_]*\}\}')) { throw "unresolved variable in $($block.Path)" }
    if (-not $outputPaths.Add($block.Path)) { throw "duplicate output path requires an explicit merge plan: $($block.Path)" }
    [void](Resolve-ContainedPath $destinationPath $block.Path 'output path')
    $outputHash = Get-Sha256 $utf8.GetBytes($content)
    $outputs[$block.Path] = [pscustomobject]@{
      Content = $content; OutputHash = $outputHash; SourceHash = $block.SourceHash
      PackId = $pack.Id; PackVersion = $pack.Version; BlockId = $block.BlockId
    }
  }
  foreach ($property in $variables.PSObject.Properties) {
    if (-not $requiredVariables.Contains($property.Name)) { throw "unexpected variable for $($pack.Id): $($property.Name)" }
  }
  $packRecords.Add([pscustomobject]@{
    Id = $pack.Id; Version = $pack.Version; Sha256 = $pack.Sha256
    Implementation = $pack.Implementation; Admission = $pack.Admission
  })
}

if ($outputs.Count -eq 0) { throw 'composition selected no files' }
if (Test-Path -LiteralPath $destinationPath) {
  if ((Get-ChildItem -LiteralPath $destinationPath -Force | Measure-Object).Count -gt 0) { throw "destination must be empty: $destinationPath" }
} else {
  [void](New-Item -ItemType Directory -Path $destinationPath)
}

foreach ($entry in $outputs.GetEnumerator()) {
  $target = Resolve-ContainedPath $destinationPath $entry.Key 'output path'
  $parent = Split-Path -Parent $target
  if (-not (Test-Path -LiteralPath $parent)) { [void](New-Item -ItemType Directory -Path $parent) }
  [IO.File]::WriteAllBytes($target, $utf8.GetBytes($entry.Value.Content))
}

$record = [Collections.Generic.List[string]]::new()
$record.Add('# Materialization Record')
$record.Add('')
$record.Add('- composition version: 1.0')
$record.Add("- plan SHA-256: $(Get-Sha256 ([IO.File]::ReadAllBytes($planPath)))")
$record.Add("- generated UTC: $([DateTimeOffset]::UtcNow.ToString('O'))")
$record.Add('')
$record.Add('## Packs')
$record.Add('')
$record.Add('| Pack | Version | Pack SHA-256 | Implementation | Admission |')
$record.Add('|---|---:|---|---|---|')
foreach ($item in $packRecords) { $record.Add("| $($item.Id) | $($item.Version) | $($item.Sha256) | $($item.Implementation) | $($item.Admission) |") }
$record.Add('')
$record.Add('## Files')
$record.Add('')
$record.Add('| Path | Pack/block | Source SHA-256 | Output SHA-256 |')
$record.Add('|---|---|---|---|')
foreach ($entry in $outputs.GetEnumerator()) {
  $value = $entry.Value
  $record.Add("| $($entry.Key) | $($value.PackId)/$($value.BlockId) | $($value.SourceHash) | $($value.OutputHash) |")
}
$recordText = [string]::Join("`n", $record) + "`n"
[IO.File]::WriteAllBytes((Join-Path $destinationPath 'MATERIALIZATION_RECORD.md'), $utf8.GetBytes($recordText))
Write-Output "Materialized $($outputs.Count) files from $($packRecords.Count) pack(s) into $destinationPath"
````

### FILE: `tools/test-compose-markdown-project.ps1`

```yaml
block_id: "MARKDOWN-COMPOSITOR:test:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "7e57209c5ccb8ad4cf5c4d337329855e409eef7df1633d3b88a85df6e0a49dcb"
variables: []
secrets_allowed: false
```

````powershell
#requires -Version 7.0

[CmdletBinding()]
param([string] $Composer = (Join-Path $PSScriptRoot 'compose-markdown-project.ps1'))

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
$utf8 = [Text.UTF8Encoding]::new($false)
$testRoot = Join-Path ([IO.Path]::GetTempPath()) ('markdown-compositor-test-' + [guid]::NewGuid().ToString('N'))
[void](New-Item -ItemType Directory -Path $testRoot)

function Hash([string] $Value) {
  [Convert]::ToHexString([Security.Cryptography.SHA256]::HashData($utf8.GetBytes($Value))).ToLowerInvariant()
}

function Write-Utf8([string] $Path, [string] $Value) {
  $parent = Split-Path -Parent $Path
  if (-not (Test-Path -LiteralPath $parent)) { [void](New-Item -ItemType Directory -Path $parent) }
  [IO.File]::WriteAllBytes($Path, $utf8.GetBytes($Value))
}

function Assert([bool] $Condition, [string] $Message) {
  if (-not $Condition) { throw "ASSERTION FAILED: $Message" }
}

function Expect-Failure([scriptblock] $Action, [string] $Pattern) {
  try {
    & $Action | Out-Null
    throw 'expected failure did not occur'
  } catch {
    if ($_.Exception.Message -notmatch $Pattern) { throw "unexpected failure: $($_.Exception.Message)" }
  }
}

function New-Pack {
  param(
    [string] $Id = 'TEST-PACK',
    [string] $Path = 'app.txt',
    [string] $Content = "project={{PROJECT_NAME}}`n",
    [string] $BlockId = 'TEST-PACK:app:v1',
    [string] $ManifestOperation = 'CREATE',
    [string] $BlockOperation = 'CREATE'
  )
  $template = @'
# Test pack

## 1. Metadata

```yaml
pack_id: "__ID__"
pack_version: "1.0.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: RECONSTRUCTIBLE
  admission: CONDITIONED
claim: "test fixture"
stacks: ["text"]
compatible_with: []
incompatible_with: []
license_expression: "LicenseRef-Test"
upstream_sources: []
verified_at: "2026-08-25"
```

## 2. Manifest

```text
__MANIFEST_OPERATION__ __PATH__
```

## 3. Materialization blocks

### FILE: `__PATH__`

```yaml
block_id: "__BLOCK_ID__"
operation: __BLOCK_OPERATION__
provenance: AUTHORED
source: "local test fixture"
license: "LicenseRef-Test"
sha256: "__HASH__"
variables: [PROJECT_NAME]
secrets_allowed: false
```

__CONTENT_FENCE__text
__CONTENT____CONTENT_FENCE__
'@
  $template.TrimStart("`n").Replace('__ID__', $Id).Replace('__PATH__', $Path).
    Replace('__BLOCK_ID__', $BlockId).Replace('__MANIFEST_OPERATION__', $ManifestOperation).
    Replace('__BLOCK_OPERATION__', $BlockOperation).Replace('__HASH__', (Hash $Content)).
    Replace('__CONTENT_FENCE__', ('`' * 4)).Replace('__CONTENT__', $Content)
}

function New-Selection {
  param(
    [string] $Path,
    [string] $Id,
    [bool] $Acknowledge = $true,
    [string[]] $Files = @('*'),
    [hashtable] $Variables = @{ PROJECT_NAME = 'elite' }
  )
  [pscustomobject]@{
    path = $Path
    packId = $Id
    version = '1.0.0'
    acknowledgeConditions = $Acknowledge
    files = $Files
    variables = [pscustomobject]$Variables
  }
}

function Write-Plan([string] $Path, [object[]] $Selections) {
  $plan = [ordered]@{ compositionVersion = '1.0'; secretVariables = [pscustomobject]@{}; packs = $Selections }
  Write-Utf8 $Path (($plan | ConvertTo-Json -Depth 20) + "`n")
}

try {
  $library = Join-Path $testRoot 'library'
  [void](New-Item -ItemType Directory -Path $library)
  $pack = New-Pack
  Write-Utf8 (Join-Path $library 'pack.md') $pack
  $planPath = Join-Path $testRoot 'plan.json'
  Write-Plan $planPath @((New-Selection -Path 'pack.md' -Id 'TEST-PACK'))
  $target = Join-Path $testRoot 'target'
  & $Composer -PlanFile $planPath -LibraryRoot $library -Destination $target | Out-Null
  Assert ((Get-Content -Raw -LiteralPath (Join-Path $target 'app.txt')) -eq "project=elite`n") 'variable materialization'
  Assert (Test-Path -LiteralPath (Join-Path $target 'MATERIALIZATION_RECORD.md')) 'record exists'

  foreach ($caseName in @('empty', 'null', 'single-explicit')) {
    $selection = New-Selection -Path 'pack.md' -Id 'TEST-PACK'
    if ($caseName -eq 'empty') { $selection.files = @() }
    elseif ($caseName -eq 'null') { $selection.files = $null }
    else { $selection.files = @('app.txt') }
    $strictPlan = Join-Path $testRoot "$caseName.json"
    $strictTarget = Join-Path $testRoot "$caseName-target"
    Write-Plan $strictPlan @($selection)
    & $Composer -PlanFile $strictPlan -LibraryRoot $library -Destination $strictTarget | Out-Null
    Assert ((Get-Content -Raw -LiteralPath (Join-Path $strictTarget 'app.txt')) -ceq "project=elite`n") "StrictMode selection $caseName"
  }

  Write-Utf8 (Join-Path $library 'crlf.md') $pack.Replace("`n", "`r`n")
  $crlfPlan = Join-Path $testRoot 'crlf.json'
  Write-Plan $crlfPlan @((New-Selection -Path 'crlf.md' -Id 'TEST-PACK'))
  & $Composer -PlanFile $crlfPlan -LibraryRoot $library -Destination (Join-Path $testRoot 'crlf-target') | Out-Null

  $noAckPlan = Join-Path $testRoot 'no-ack.json'
  Write-Plan $noAckPlan @((New-Selection -Path 'pack.md' -Id 'TEST-PACK' -Acknowledge $false))
  Expect-Failure { & $Composer -PlanFile $noAckPlan -LibraryRoot $library -Destination (Join-Path $testRoot 'no-ack-target') } 'acknowledgeConditions'

  $badHash = [regex]::Replace($pack, '(?m)^sha256: "[0-9a-f]{64}"$', 'sha256: "' + ('0' * 64) + '"')
  Write-Utf8 (Join-Path $library 'bad-hash.md') $badHash
  $badHashPlan = Join-Path $testRoot 'bad-hash.json'
  Write-Plan $badHashPlan @((New-Selection -Path 'bad-hash.md' -Id 'TEST-PACK'))
  Expect-Failure { & $Composer -PlanFile $badHashPlan -LibraryRoot $library -Destination (Join-Path $testRoot 'bad-hash-target') } 'SHA-256 mismatch'

  Write-Utf8 (Join-Path $library 'manifest-mismatch.md') $pack.Replace('CREATE app.txt', 'CREATE absent.txt')
  $manifestPlan = Join-Path $testRoot 'manifest-mismatch.json'
  Write-Plan $manifestPlan @((New-Selection -Path 'manifest-mismatch.md' -Id 'TEST-PACK'))
  Expect-Failure { & $Composer -PlanFile $manifestPlan -LibraryRoot $library -Destination (Join-Path $testRoot 'manifest-target') } 'no exact FILE block'

  Write-Utf8 (Join-Path $library 'duplicate-manifest.md') $pack.Replace('CREATE app.txt', "CREATE app.txt`nCREATE app.txt")
  $duplicateManifestPlan = Join-Path $testRoot 'duplicate-manifest.json'
  Write-Plan $duplicateManifestPlan @((New-Selection -Path 'duplicate-manifest.md' -Id 'TEST-PACK'))
  Expect-Failure { & $Composer -PlanFile $duplicateManifestPlan -LibraryRoot $library -Destination (Join-Path $testRoot 'duplicate-manifest-target') } 'duplicate manifest path'

  Write-Utf8 (Join-Path $library 'patch.md') (New-Pack -ManifestOperation PATCH -BlockOperation PATCH)
  $patchPlan = Join-Path $testRoot 'patch.json'
  Write-Plan $patchPlan @((New-Selection -Path 'patch.md' -Id 'TEST-PACK'))
  Expect-Failure { & $Composer -PlanFile $patchPlan -LibraryRoot $library -Destination (Join-Path $testRoot 'patch-target') } 'only exact CREATE'

  Write-Utf8 (Join-Path $library 'missing-field.md') $pack.Replace("provenance: AUTHORED`n", '')
  $missingFieldPlan = Join-Path $testRoot 'missing-field.json'
  Write-Plan $missingFieldPlan @((New-Selection -Path 'missing-field.md' -Id 'TEST-PACK'))
  Expect-Failure { & $Composer -PlanFile $missingFieldPlan -LibraryRoot $library -Destination (Join-Path $testRoot 'missing-field-target') } 'provenance is missing'

  $missingSelectionPlan = Join-Path $testRoot 'missing-selection.json'
  Write-Plan $missingSelectionPlan @((New-Selection -Path 'pack.md' -Id 'TEST-PACK' -Files @('missing.txt')))
  Expect-Failure { & $Composer -PlanFile $missingSelectionPlan -LibraryRoot $library -Destination (Join-Path $testRoot 'missing-selection-target') } 'selected file does not exist'

  $aliasSelectionPlan = Join-Path $testRoot 'alias-selection.json'
  Write-Plan $aliasSelectionPlan @((New-Selection -Path 'pack.md' -Id 'TEST-PACK' -Files @('./app.txt')))
  Expect-Failure { & $Composer -PlanFile $aliasSelectionPlan -LibraryRoot $library -Destination (Join-Path $testRoot 'alias-selection-target') } 'unsafe'

  $duplicatePackPlan = Join-Path $testRoot 'duplicate-pack.json'
  Write-Plan $duplicatePackPlan @(
    (New-Selection -Path 'pack.md' -Id 'TEST-PACK'),
    (New-Selection -Path 'pack.md' -Id 'TEST-PACK')
  )
  Expect-Failure { & $Composer -PlanFile $duplicatePackPlan -LibraryRoot $library -Destination (Join-Path $testRoot 'duplicate-pack-target') } 'duplicate packId selection'

  Write-Utf8 (Join-Path $library 'same-path.md') (New-Pack -Id 'OTHER-PACK' -BlockId 'OTHER-PACK:app:v1')
  $collisionPlan = Join-Path $testRoot 'collision.json'
  Write-Plan $collisionPlan @(
    (New-Selection -Path 'pack.md' -Id 'TEST-PACK'),
    (New-Selection -Path 'same-path.md' -Id 'OTHER-PACK')
  )
  Expect-Failure { & $Composer -PlanFile $collisionPlan -LibraryRoot $library -Destination (Join-Path $testRoot 'collision-target') } 'duplicate output path'

  Write-Utf8 (Join-Path $library 'duplicate-block-id.md') (New-Pack -Id 'OTHER-PACK' -Path 'other.txt' -BlockId 'TEST-PACK:app:v1')
  $duplicateBlockPlan = Join-Path $testRoot 'duplicate-block.json'
  Write-Plan $duplicateBlockPlan @(
    (New-Selection -Path 'pack.md' -Id 'TEST-PACK'),
    (New-Selection -Path 'duplicate-block-id.md' -Id 'OTHER-PACK')
  )
  Expect-Failure { & $Composer -PlanFile $duplicateBlockPlan -LibraryRoot $library -Destination (Join-Path $testRoot 'duplicate-block-target') } 'duplicate selected block_id'

  Write-Utf8 (Join-Path $library 'reserved.md') (New-Pack -Path 'MATERIALIZATION_RECORD.md')
  $reservedPlan = Join-Path $testRoot 'reserved.json'
  Write-Plan $reservedPlan @((New-Selection -Path 'reserved.md' -Id 'TEST-PACK'))
  Expect-Failure { & $Composer -PlanFile $reservedPlan -LibraryRoot $library -Destination (Join-Path $testRoot 'reserved-target') } 'reserved'

  Write-Output 'PASS: success, CRLF, condition, integrity, manifest, CREATE-only, uniqueness, selection and reserved-path tests'
} finally {
  $resolved = [IO.Path]::GetFullPath($testRoot)
  $temp = [IO.Path]::GetFullPath([IO.Path]::GetTempPath())
  if ($resolved.StartsWith($temp, [StringComparison]::OrdinalIgnoreCase) -and (Split-Path -Leaf $resolved) -like 'markdown-compositor-test-*') {
    Remove-Item -LiteralPath $resolved -Recurse -Force
  }
}
````

### FILE: `examples/PROJECT_PACK_PLAN.example.md`

```yaml
block_id: "MARKDOWN-COMPOSITOR:plan-example:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "21c80235fafddd59b887775f2d26b168bc975fc50100dc91870217386f51b87f"
variables: []
secrets_allowed: false
```

````markdown
# Project Pack Plan

El plan selecciona adapters después del blueprint. No contiene secretos.

```json
{
  "compositionVersion": "1.0",
  "secretVariables": {},
  "packs": [
    {
      "path": "implementation_packs/GO_ENTERPRISE_BACKEND_CORE.md",
      "packId": "GO-ENTERPRISE-BACKEND",
      "version": "0.4.0",
      "acknowledgeConditions": true,
      "files": ["*"],
      "variables": {}
    }
  ]
}
```

Las condiciones aceptadas deben detallarse fuera del JSON con evidencia y gates del proyecto; el booleano sólo evita adopción accidental, no demuestra que la condición esté satisfecha.
````

## 6. Configuration surface

| Entrada | Tipo | Secreto | Regla |
|---|---|---:|---|
| `PlanFile` | JSON o Markdown con JSON fenced | no | composición 1.0 |
| `LibraryRoot` | path existente | no | todos los packs quedan contenidos |
| `Destination` | path nuevo o vacío | no | ninguna escritura antes del preflight |

## 7. Dependency bill

| Tool | Pin | Uso | Licencia |
|---|---|---|---|
| PowerShell | 7+ | parser, hashes y materialización | MIT |

## 8. Apply order

1. Materializar este pack con `materialize_markdown_pack.ps1`.
2. Ejecutar su test.
3. Crear `PROJECT_PACK_PLAN.md` desde blueprint/authority map.
4. Ejecutar el compositor con library root y target vacío.
5. Conservar `MATERIALIZATION_RECORD.md` y ejecutar gates de cada pack.

## 9. Verification

```powershell
pwsh -NoProfile -File .\tools\test-compose-markdown-project.ps1
```

Debe aprobar success path y CRLF/LF, más rechazos de condición omitida, SHA inválido, manifest divergente/duplicado, operación no `CREATE`, metadata incompleta, selección inexistente/no canónica, pack ID/block ID/output path duplicado y path de registro reservado. La composición real de los perfiles backend/web también debe conservar lineage y pasar sus gates integrados.

## 10. Reconstruction evidence

Estado: `REBUILD_VERIFIED / CONDITIONED`. La revisión 0.2.1 pasó su suite negativa ampliada y recompuso los perfiles backend/web desde targets vacíos el 2026-08-25. V285 / 0.2.2 corrige el desempaquetado escalar de files bajo StrictMode Latest; toda la suite se ejecuta ahora con ese modo, incluyendo wildcard, array vacío, null y selección explícita única. Sigue condicionado porque sólo soporta `CREATE`, no merges `PATCH`/`DELETE`, y la compatibilidad de negocio/plataforma continúa siendo responsabilidad del plan específico.

V402 composed delta: Optional final_newline false preserves exact source bytes; legacy final LF default remains. Malformed/duplicate flags and contradictory trailing empty lines fail before writes. Forty-two EOF cases and existing compositor/updater regressions PASS; source bytes unchanged after tests.
