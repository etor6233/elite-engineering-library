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
  if (-not $content.EndsWith("`n")) { $content += "`n" }
  $bytes = [Text.UTF8Encoding]::new($false).GetBytes($content)
  $hash = [Convert]::ToHexString([Security.Cryptography.SHA256]::HashData($bytes)).ToLowerInvariant()
  $escaped = [regex]::Escape($relative)
  $pattern = '(?ms)(### FILE: `' + $escaped + '`\r?\n.*?sha256: ")[0-9a-f]{64}(".*?\r?\n````[^\r\n]*\r?\n)(.*?)(\r?\n````(?=\r?\n|$))'
  $matches = [regex]::Matches($markdown, $pattern)
  if ($matches.Count -ne 1) {
    throw "Expected exactly one materialization block for $relative; found $($matches.Count)"
  }
  # The Markdown fence contributes exactly one trailing newline when the
  # materializer joins content lines. Remove only that one newline here so
  # source files ending in multiple newlines remain byte-identical.
  $escapedContent = $content.Substring(0, $content.Length - 1)
  $regex = [regex]::new($pattern)
  $markdown = $regex.Replace($markdown, [Text.RegularExpressions.MatchEvaluator]{
    param($match)
    $match.Groups[1].Value + $hash + $match.Groups[2].Value + $escapedContent + $match.Groups[4].Value
  }, 1)
}

[IO.File]::WriteAllText($packPath, $markdown, [Text.UTF8Encoding]::new($false))
Write-Output "Updated $($Files.Count) blocks in $packPath"
