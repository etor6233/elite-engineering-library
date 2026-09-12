#requires -Version 7.0
<# AUTHORED library utility. Windows DPAPI uses Microsoft's Export/Import-Clixml;
   this is not a Microsoft-authored vault or a production secret manager.
   See PROJECT_SECRETS_TEMPLATE.md for limits, recovery and provider gates. #>
[CmdletBinding()]
param(
  [Parameter(Mandatory)][ValidateSet('Set','Generate','Check','Run')][string]$Mode,
  [Parameter(Mandatory)][ValidatePattern('^[a-z][a-z0-9-]{2,63}$')][string]$ProjectId,
  [Parameter(Mandatory)][string]$ProjectRoot,
  [string[]]$Names = @(),
  [string]$Executable,
  [string[]]$Arguments = @(),
  [string]$PublicConfigFile
)
Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
if (-not $IsWindows) { throw 'LOCAL_SECRETS_WINDOWS_ONLY' }
$project = (Resolve-Path -LiteralPath $ProjectRoot).Path.TrimEnd('\')
$vault = Join-Path ([Environment]::GetFolderPath('LocalApplicationData')) "EliteEngineeringSecrets/$ProjectId/dev"
$vault = [IO.Path]::GetFullPath($vault)
if ($vault.StartsWith($project + '\', [StringComparison]::OrdinalIgnoreCase) -or $vault -eq $project) { throw 'VAULT_MUST_BE_OUTSIDE_PROJECT' }
if ($Names.Count -eq 0 -or @($Names | Sort-Object -Unique).Count -ne $Names.Count) { throw 'NAMES_REQUIRED_UNIQUE' }
foreach ($name in $Names) { if ($name -cnotmatch '^[A-Z][A-Z0-9_]{2,127}$') { throw 'INVALID_SECRET_NAME' } }
foreach ($name in $Names) {
  if ($name -in @('PATH','PATHEXT','COMSPEC','SYSTEMROOT','WINDIR','TEMP','TMP','USERPROFILE','LOCALAPPDATA','NODE_OPTIONS','PYTHONPATH','PSMODULEPATH','LD_PRELOAD','HTTP_PROXY','HTTPS_PROXY','ALL_PROXY')) { throw 'PROCESS_CONTROL_IS_NOT_A_SECRET_BINDING' }
}
if ($Mode -in @('Set','Generate') -and $Names.Count -ne 1) { throw 'ONE_SECRET_PER_WRITE' }

function Assert-NoLink([string]$Path) {
  $cursor = $Path
  while ($cursor) {
    if (Test-Path -LiteralPath $cursor) {
      if (((Get-Item -LiteralPath $cursor -Force).Attributes -band [IO.FileAttributes]::ReparsePoint) -ne 0) { throw 'REPARSE_POINT_REJECTED' }
    }
    $parent = Split-Path -Path $cursor -Parent
    if ($parent -eq $cursor) { break }
    $cursor = $parent
  }
}
Assert-NoLink $vault
$sid = [Security.Principal.WindowsIdentity]::GetCurrent().User
if ($Mode -in @('Set','Generate')) {
  if (-not (Test-Path -LiteralPath $vault)) { $null = New-Item -ItemType Directory -Path $vault }
  $acl = Get-Acl -LiteralPath $vault
  $acl.SetAccessRuleProtection($true, $false)
  foreach ($existing in @($acl.Access)) { $null = $acl.RemoveAccessRuleSpecific($existing) }
  foreach ($principal in @($sid, [Security.Principal.SecurityIdentifier]::new('S-1-5-18'))) {
    $rule = [Security.AccessControl.FileSystemAccessRule]::new($principal, 'FullControl', 'ContainerInherit,ObjectInherit', 'None', 'Allow')
    $acl.AddAccessRule($rule)
  }
  [IO.FileSystemAclExtensions]::SetAccessControl([IO.DirectoryInfo]::new($vault), $acl)
}

function Read-Protected([string]$Name) {
  $path = Join-Path $vault "$Name.clixml"
  Assert-NoLink $path
  if (-not (Test-Path -LiteralPath $path -PathType Leaf)) { throw 'SECRET_MISSING' }
  $item = Get-Item -LiteralPath $path
  if ($item.Length -gt 65536) { throw 'SECRET_FILE_TOO_LARGE' }
  $acl = Get-Acl -LiteralPath $path
  foreach ($rule in $acl.Access) {
    if ($rule.AccessControlType -eq 'Allow' -and $rule.IdentityReference.Translate([Security.Principal.SecurityIdentifier]).Value -notin @($sid.Value, 'S-1-5-18')) { throw 'SECRET_ACL_TOO_BROAD' }
  }
  try { $credential = Import-Clixml -LiteralPath $path } catch { throw 'SECRET_UNREADABLE_REENTER_OR_REISSUE' }
  if ($credential -isnot [PSCredential] -or $credential.UserName -cne $Name -or $credential.Password.Length -eq 0) { throw 'SECRET_INVALID' }
  return $credential
}

if ($Mode -in @('Set','Generate')) {
  $name = $Names[0]
  $path = Join-Path $vault "$name.clixml"
  Assert-NoLink $path
  if (Test-Path -LiteralPath $path) { throw 'SECRET_EXISTS_USE_NEW_VERSIONED_NAME' }
  $secure = $null
  try {
    if ($Mode -eq 'Set') { $secure = Read-Host "Ingresar $name (no se muestra)" -AsSecureString }
    else {
      $bytes = [Security.Cryptography.RandomNumberGenerator]::GetBytes(32)
      try { $secure = ConvertTo-SecureString ([Convert]::ToBase64String($bytes)) -AsPlainText -Force }
      finally { [Array]::Clear($bytes) }
    }
    if ($secure.Length -eq 0) { throw 'EMPTY_SECRET_REJECTED' }
    [PSCredential]::new($name, $secure) | Export-Clixml -LiteralPath $path -NoClobber
    $null = Read-Protected $name
    Write-Output "$name STORED_DEV_DPAPI"
  } finally { if ($secure) { $secure.Dispose() } }
  return
}
if ($Mode -eq 'Check') {
  $failed = $false
  foreach ($name in $Names) {
    try { $credential = Read-Protected $name; $credential.Password.Dispose(); Write-Output "$name AVAILABLE_NOT_PROVIDER_VALIDATED" }
    catch { $failed = $true; Write-Output "$name UNAVAILABLE" }
  }
  if ($failed) { throw 'LOCAL_SECRET_CHECK_FAILED' }
  return
}

# Run only a reviewed local executable. No shell command evaluation or secret arguments.
if (-not [IO.Path]::IsPathFullyQualified($Executable) -or -not (Test-Path -LiteralPath $Executable -PathType Leaf)) { throw 'EXECUTABLE_ABSOLUTE_FILE_REQUIRED' }
Assert-NoLink $Executable
$start = [Diagnostics.ProcessStartInfo]::new()
$start.FileName = $Executable
$start.WorkingDirectory = $project
$start.UseShellExecute = $false
$start.CreateNoWindow = $true
$start.RedirectStandardOutput = $true
$start.RedirectStandardError = $true
$start.Environment.Clear()
foreach ($key in @('PATH','SystemRoot','WINDIR','TEMP','TMP','LOCALAPPDATA','USERPROFILE','PATHEXT')) {
  $value = [Environment]::GetEnvironmentVariable($key, 'Process')
  if ($null -ne $value) { $start.Environment[$key] = $value }
}
$publicNames = @('NODE_ENV','BUSINESS_CONFIG_FILE','APP_BASE_URL','OIDC_ISSUER','OIDC_AUDIENCE','OIDC_CLIENT_ID','ENTERPRISE_API_BASE_URL','ENTERPRISE_TENANT_CODE','ENTERPRISE_ORGANIZATION_CODE','HTTP_ADDRESS','ARCA_ENABLED','TENANT_ID','ORGANIZATION_ID','LISTEN_ADDR','META_APP_ID','GOOGLE_ADS_CLIENT_ID','GOOGLE_ADS_LOGIN_CUSTOMER_ID','GOOGLE_ADS_USE_PROTO_PLUS','GOOGLE_ADS_USE_APPLICATION_DEFAULT_CREDENTIALS')
if ($PublicConfigFile) {
  $config = Get-Content -LiteralPath $PublicConfigFile -Raw | ConvertFrom-Json -AsHashtable
  if ($config -isnot [System.Collections.IDictionary]) { throw 'PUBLIC_CONFIG_OBJECT_REQUIRED' }
  foreach ($key in $config.Keys) {
    if ($key -cnotin $publicNames -or $key -cin $Names -or $config[$key] -isnot [string]) { throw 'PUBLIC_CONFIG_UNKNOWN_OR_SECRET_FIELD' }
    if ($key -ceq 'NODE_ENV' -and $config[$key] -ceq 'production') { throw 'DEV_VAULT_NOT_FOR_PRODUCTION' }
    $start.Environment[$key] = $config[$key]
  }
}
foreach ($argument in $Arguments) { $start.ArgumentList.Add($argument) }
try {
  foreach ($name in $Names) {
    $credential = Read-Protected $name
    try { $start.Environment[$name] = $credential.GetNetworkCredential().Password }
    finally { $credential.Password.Dispose() }
  }
  $process = [Diagnostics.Process]::Start($start)
  $start.Environment.Clear()
  # Drain without retaining or publishing child output, which may contain secrets.
  $stdout = $process.StandardOutput.BaseStream.CopyToAsync([IO.Stream]::Null)
  $stderr = $process.StandardError.BaseStream.CopyToAsync([IO.Stream]::Null)
  $process.WaitForExit()
  $stdout.GetAwaiter().GetResult()
  $stderr.GetAwaiter().GetResult()
  if ($process.ExitCode -ne 0) { throw 'CHILD_PROCESS_FAILED' }
  Write-Output 'CHILD_EXIT_SUCCESS_NOT_PROVIDER_VALIDATED'
} finally {
  $start.Environment.Clear()
  if (Get-Variable process -ErrorAction SilentlyContinue) { if ($process) { $process.Dispose() } }
}
