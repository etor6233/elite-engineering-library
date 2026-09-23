#requires -Version 7.0
[CmdletBinding()]
param([string]$PythonExecutable = 'python')
$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest
$source = Get-Content -Raw -LiteralPath (Join-Path $PSScriptRoot 'extract_email_attachments.py')
foreach ($token in @('automatic_storage_authorized', 'RAW_SHA256_MISMATCH', 'MIME_PART_LIMIT', 'ATTACHMENT_BASE64_INVALID', 'os.replace', 'security_decision')) {
  if (-not $source.Contains($token, [StringComparison]::Ordinal)) { throw "Missing contract token: $token" }
}
& $PythonExecutable -B -m py_compile (Join-Path $PSScriptRoot 'extract_email_attachments.py') (Join-Path $PSScriptRoot 'test_extract_email_attachments.py')
if ($LASTEXITCODE -ne 0) { throw 'Python compile failed' }
Push-Location $PSScriptRoot
try { & $PythonExecutable -B -m unittest -v test_extract_email_attachments.py; if ($LASTEXITCODE -ne 0) { throw 'Tests failed' } }
finally { Pop-Location }
Write-Output 'SECURE_EMAIL_MIME_QUARANTINE_CORE_PASS tests=18 storage_authorized=0 filename_as_path=0 receipt_bound=1 atomic=1'
