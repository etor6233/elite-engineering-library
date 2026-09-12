#requires -Version 7.0

[CmdletBinding()]
param(
  [string] $Updater = (Join-Path $PSScriptRoot 'update_pack_from_tree.ps1'),
  [string] $Materializer = (Join-Path (Split-Path -Parent $PSScriptRoot) 'materialize_markdown_pack.ps1')
)

$ErrorActionPreference = 'Stop'
$utf8 = [Text.UTF8Encoding]::new($false)
$testRoot = Join-Path ([IO.Path]::GetTempPath()) ('update-pack-from-tree-test-' + [guid]::NewGuid().ToString('N'))

function Write-Utf8([string] $Path, [string] $Value) {
  $parent = Split-Path -Parent $Path
  if (-not (Test-Path -LiteralPath $parent)) { [void](New-Item -ItemType Directory -Path $parent) }
  [IO.File]::WriteAllBytes($Path, $utf8.GetBytes($Value))
}

function Hash-Text([string] $Value) {
  [Convert]::ToHexString([Security.Cryptography.SHA256]::HashData($utf8.GetBytes($Value))).ToLowerInvariant()
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

try {
  [void](New-Item -ItemType Directory -Path $testRoot)
  $sourceRoot = Join-Path $testRoot 'source'
  [void](New-Item -ItemType Directory -Path $sourceRoot)
  $oldContent = "old`n"
  $pack = @'
# Update fixture

## 1. Metadata

```yaml
pack_id: "UPDATE-FIXTURE"
pack_version: "1.0.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: RECONSTRUCTIBLE
  admission: CONDITIONED
claim: "updater regression"
stacks: ["text"]
compatible_with: []
incompatible_with: []
license_expression: "LicenseRef-Test"
upstream_sources: []
verified_at: "2026-08-25"
```

## 2. Exact file manifest

```text
CREATE artifact.txt
```

## 3. Materialization blocks

### FILE: `artifact.txt`

```yaml
block_id: "UPDATE-FIXTURE:artifact:v1"
operation: CREATE
provenance: AUTHORED
source: "local test"
license: "LicenseRef-Test"
sha256: "__OLD_HASH__"
variables: []
secrets_allowed: false
```

````text
old
````
'@
  $pack = $pack.TrimStart("`n").Replace('__OLD_HASH__', (Hash-Text $oldContent))
  $packPath = Join-Path $testRoot 'pack.md'
  Write-Utf8 $packPath $pack

  $newContent = @'
alpha
````text
$1 $&
__CONTENT__````
omega
'@
  $newContent = $newContent.TrimStart("`n")
  $newContent = $newContent.TrimEnd("`n") + "`n`n"
  Write-Utf8 (Join-Path $sourceRoot 'artifact.txt') $newContent

  & $Updater -PackFile $packPath -SourceRoot $sourceRoot -Files @('artifact.txt') | Out-Null
  $destination = Join-Path $testRoot 'materialized'
  & $Materializer -PackFile $packPath -Destination $destination | Out-Null
  $actual = [IO.File]::ReadAllText((Join-Path $destination 'artifact.txt'), $utf8)
  Assert ($actual -ceq $newContent) 'nested fence and literal replacement content'
  Assert ($actual.Contains('$1') -and $actual.Contains('$&')) 'replacement metacharacters remain literal'

  $beforeAtomicFailure = (Get-FileHash -LiteralPath $packPath -Algorithm SHA256).Hash
  Expect-Failure {
    & $Updater -PackFile $packPath -SourceRoot $sourceRoot -Files @('artifact.txt', 'missing.txt')
  } 'Source file missing'
  $afterAtomicFailure = (Get-FileHash -LiteralPath $packPath -Algorithm SHA256).Hash
  Assert ($beforeAtomicFailure -eq $afterAtomicFailure) 'failed multi-file update must not change pack'

  Expect-Failure {
    & $Updater -PackFile $packPath -SourceRoot $sourceRoot -Files @('../escape.txt')
  } 'Unsafe source path'

  Write-Output 'PASS: nested fence, literal replacement, multiple trailing newlines, hash reconstruction, atomic failure and traversal tests'
} finally {
  if (Test-Path -LiteralPath $testRoot) {
    $resolved = [IO.Path]::GetFullPath($testRoot)
    $systemTemp = [IO.Path]::GetFullPath([IO.Path]::GetTempPath())
    if ($resolved.StartsWith($systemTemp, [StringComparison]::OrdinalIgnoreCase) -and (Split-Path -Leaf $resolved) -like 'update-pack-from-tree-test-*') {
      Remove-Item -LiteralPath $resolved -Recurse -Force
    }
  }
}
