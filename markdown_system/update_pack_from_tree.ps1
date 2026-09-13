#requires -Version 7.0

[CmdletBinding()]
param(
  [Parameter(Mandatory = $true)]
  [string] $PackFile,
  [Parameter(Mandatory = $true)]
  [string] $SourceRoot,
  [Parameter(Mandatory = $true)]
  [string[]] $Files
)

$ErrorActionPreference = 'Stop'
$packPath = (Resolve-Path -LiteralPath $PackFile).Path
$sourcePath = (Resolve-Path -LiteralPath $SourceRoot).Path
$markdown = [IO.File]::ReadAllText($packPath)

foreach ($file in $Files) {
  $relative = $file.Replace('\', '/')
  if ([IO.Path]::IsPathRooted($relative) -or $relative.Split('/') -contains '..') {
    throw "Unsafe source path: $relative"
  }
  $sourceFile = [IO.Path]::GetFullPath((Join-Path $sourcePath $relative))
  if (-not $sourceFile.StartsWith($sourcePath + [IO.Path]::DirectorySeparatorChar, [StringComparison]::OrdinalIgnoreCase)) {
    throw "Source escaped root: $relative"
  }
  if (-not (Test-Path -LiteralPath $sourceFile -PathType Leaf)) {
    throw "Source file missing: $relative"
  }
  $content = [IO.File]::ReadAllText($sourceFile).Replace("`r`n", "`n").Replace("`r", "`n")
  $finalNewline = $content.EndsWith("`n")
  $bytes = [Text.UTF8Encoding]::new($false).GetBytes($content)
  $hash = [Convert]::ToHexString([Security.Cryptography.SHA256]::HashData($bytes)).ToLowerInvariant()
  $escaped = [regex]::Escape($relative)
  $pattern = '(?ms)(### FILE: `' + $escaped + '`\r?\n)(.*?)(^````[^\r\n]*\r?\n)(.*?)(^````(?=\r?$))'
  $matches = [regex]::Matches($markdown, $pattern)
  if ($matches.Count -ne 1) {
    throw "Expected exactly one materialization block for $relative; found $($matches.Count)"
  }
  $metadata = $matches[0].Groups[2].Value
  $hashFields = [regex]::Matches($metadata, '(?m)^sha256: "[0-9a-f]{64}"[ \t]*\r?$')
  if ($hashFields.Count -ne 1) { throw "Expected exactly one sha256 for $relative" }
  $newlineFields = [regex]::Matches($metadata, '(?m)^final_newline:[^\n]*$')
  if ($newlineFields.Count -gt 1 -or ($newlineFields.Count -eq 1 -and $newlineFields[0].Value -cnotmatch '^final_newline:[ \t]*(true|false)[ \t]*\r?$')) {
    throw "Invalid or duplicate final_newline for $relative"
  }
  $metadata = [regex]::Replace($metadata, '(?m)^sha256: "[0-9a-f]{64}"[ \t]*\r?$', 'sha256: "' + $hash + '"')
  $metadata = [regex]::Replace($metadata, '(?m)^final_newline:[^\n]*\n', '')
  if (-not $finalNewline) {
    $metadata = [regex]::Replace($metadata, '(?m)^sha256:', "final_newline: false`nsha256:")
  }
  # The separator before the closing fence is not part of a no-LF payload.
  # Empty payloads have adjacent opening/closing lines; LF payloads retain all LFs.
  $payload = $content
  if ($content.Length -gt 0 -and -not $finalNewline) { $payload += "`n" }
  $regex = [regex]::new($pattern)
  $markdown = $regex.Replace($markdown, [Text.RegularExpressions.MatchEvaluator]{
    param($match)
    $match.Groups[1].Value + $metadata + $match.Groups[3].Value + $payload + $match.Groups[5].Value
  }, 1)
}

[IO.File]::WriteAllText($packPath, $markdown, [Text.UTF8Encoding]::new($false))
Write-Output "Updated $($Files.Count) blocks in $packPath"
