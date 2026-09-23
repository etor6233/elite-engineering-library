#requires -Version 7.0
[CmdletBinding()]
param(
  [Parameter(Mandatory=$true)][string]$DotnetExecutable,
  [Parameter(Mandatory=$true)][string]$CacheDirectory,
  [Parameter(Mandatory=$true)][string]$OutputDirectory,
  [Parameter(Mandatory=$true)][string]$SvcutilRuntimeDirectory,
  [Parameter(Mandatory=$true)][string]$OfflinePackageDirectory,
  [ValidateSet('Homologation','Production')][string]$Environment='Homologation',
  [switch]$AllowProductionMetadata,
  [switch]$Offline
)
$ErrorActionPreference='Stop'
function Hash([string]$Path){(Get-FileHash -LiteralPath $Path -Algorithm SHA256).Hash.ToLowerInvariant()}
function Plain([string]$Path){
  $full=[IO.Path]::GetFullPath($Path)
  if($full.StartsWith('\\')){throw 'local paths only'}
  $cursor=$full
  while($cursor){if(Test-Path -LiteralPath $cursor){if((Get-Item -Force -LiteralPath $cursor).Attributes-band[IO.FileAttributes]::ReparsePoint){throw 'reparse paths rejected'}};$cursor=[IO.Path]::GetDirectoryName($cursor)}
  $full
}
function Run([string]$Exe,[string[]]$Arguments){$lines=@(& $Exe @Arguments 2>&1);if($LASTEXITCODE-ne 0){throw "command failed ($LASTEXITCODE): $Exe $($Arguments-join' ')`n$($lines-join[Environment]::NewLine)"};$lines}
$dotnet=Plain $DotnetExecutable
if(-not(Test-Path -LiteralPath $dotnet -PathType Leaf)){throw 'dotnet executable missing'}
if((Hash $dotnet)-ne'ab1b71fd3dd71062e074c9fab8312081a81b7f2b3e0327c48c4d249c8d1a3135'){throw 'dotnet executable pin mismatch'}
$assets=[IO.Path]::GetFullPath((Join-Path $PSScriptRoot '..\arca\fiscal'))
$lock=Get-Content -LiteralPath (Join-Path $assets 'arca-ws-contracts.lock.json') -Raw|ConvertFrom-Json -Depth 20
$adaptation=Get-Content -LiteralPath (Join-Path $assets 'svcutil-adaptation.lock.json') -Raw|ConvertFrom-Json -Depth 30
if(-not$Offline-or$Environment-ne'Homologation'-or$AllowProductionMetadata){throw 'this admitted generator supports only pinned offline homologation WSDL metadata'}
$svcutil=Plain $SvcutilRuntimeDirectory;$feed=Plain $OfflinePackageDirectory
if(-not(Test-Path -LiteralPath $feed -PathType Container)){throw 'offline package directory missing'}
$packageLock=Get-Content -LiteralPath (Join-Path $assets 'nuget-runtime.lock.json') -Raw|ConvertFrom-Json -Depth 20
$feedFiles=@(Get-ChildItem -Force -LiteralPath $feed -File)
if($feedFiles.Count-ne$packageLock.packages.Count-or@(Get-ChildItem -Force -LiteralPath $feed -Directory).Count-ne0){throw 'offline feed contains undeclared entries'}
foreach($package in $packageLock.packages){
  if($package.archive-match'[/\\:]'-or$package.archive-eq'..'){throw 'invalid package filename'}
  $path=Plain (Join-Path $feed $package.archive)
  if((Hash $path)-ne$package.sha256-or(Get-Item -LiteralPath $path).Length-ne$package.bytes){throw 'offline package pin mismatch'}
}
$toolReceipt=Get-Content -LiteralPath (Join-Path $svcutil 'ASSEMBLY_RECEIPT.json') -Raw|ConvertFrom-Json -Depth 30
if($toolReceipt.inventory_sha256-ne$adaptation.expected_inventory_sha256-or$toolReceipt.scope-ne'PINNED_LOCAL_WSDL_ONLY'){throw 'svcutil adaptation receipt mismatch'}
$expected=@{};foreach($row in $toolReceipt.files){
  if($row.path-match'(^/|\\|:|(^|/)\.\.(/|$))'-or$expected.ContainsKey($row.path)){throw 'invalid tool member'}
  $expected[$row.path]=$row.sha256
}
# Bind the receipt itself to the canonical inventory, then check every actual file.
# Receipt rows are independently matched to the fixed inventory stored in the source lock.
if($expected.Count-ne$adaptation.expected_files.Count){throw 'tool inventory count mismatch'}
foreach($row in $adaptation.expected_files){if($expected[$row.path]-ne$row.sha256){throw 'tool receipt file drift'}}
$actual=@(Get-ChildItem -Force -LiteralPath $svcutil -Recurse -File)
if($actual.Count-ne($expected.Count+1)){throw 'tool contains undeclared files'}
foreach($file in $actual){
  $path=Plain $file.FullName;$rel=[IO.Path]::GetRelativePath($svcutil,$path).Replace('\','/')
  if($rel-eq'ASSEMBLY_RECEIPT.json'){continue}
  if(-not$expected.ContainsKey($rel)-or(Hash $path)-ne$expected[$rel]){throw 'tool member hash mismatch'}
}
$actualSdk=(@(Run $dotnet @('--version'))[-1]).Trim()
if($actualSdk-ne$lock.toolchain.dotnetSdk){throw "dotnet SDK mismatch: expected=$($lock.toolchain.dotnetSdk) actual=$actualSdk"}
$wsaaContract=$lock.contracts.wsaaHomologation
$wsfeContract=$lock.contracts.wsfeHomologation
$cache=Plain $CacheDirectory;$output=Plain $OutputDirectory
if($cache-eq$output){throw 'cache and output must differ'}
foreach($inputRoot in @($assets,$cache,$svcutil,$feed)){
  # The selected project owns this one absent generated subtree. Other source
  # descendants and every cache/tool/feed overlap remain forbidden.
  if($inputRoot-eq$assets-and$output-eq(Join-Path $assets 'generated')){continue}
  if($output-eq$inputRoot-or$output.StartsWith($inputRoot+[IO.Path]::DirectorySeparatorChar,[StringComparison]::OrdinalIgnoreCase)-or$inputRoot.StartsWith($output+[IO.Path]::DirectorySeparatorChar,[StringComparison]::OrdinalIgnoreCase)){throw 'overlapping generation input/output'}
}
[IO.Directory]::CreateDirectory($cache)|Out-Null
if(Test-Path -LiteralPath $output){if((Get-ChildItem -Force -LiteralPath $output).Count-ne0){throw 'output must be absent or empty'}}else{[IO.Directory]::CreateDirectory($output)|Out-Null}
function Acquire([string]$Name,$Contract,[string[]]$Operations){
  $path=Join-Path $cache ("{0}-{1}.wsdl" -f $Name,$Environment.ToLowerInvariant())
  if(-not(Test-Path -LiteralPath $path -PathType Leaf)){if($Offline){throw "offline WSDL cache miss: $Name"};Invoke-WebRequest -UseBasicParsing -Uri $Contract.url -OutFile $path -TimeoutSec 30}
  if((Get-Item -LiteralPath $path).Length-ne[long]$Contract.bytes-or(Hash $path)-ne$Contract.sha256){throw "WSDL identity mismatch: $Name"}
  $settings=[Xml.XmlReaderSettings]::new();$settings.DtdProcessing=[Xml.DtdProcessing]::Prohibit;$settings.XmlResolver=$null
  $reader=[Xml.XmlReader]::Create($path,$settings);try{while($reader.Read()){}}finally{$reader.Dispose()}
  $text=[IO.File]::ReadAllText($path);foreach($operation in $Operations){if($text-notmatch[regex]::Escape($operation)){throw "required operation missing: $Name/$operation"}}
  $path
}
$wsaaWsdl=Acquire 'wsaa' $wsaaContract @('loginCms')
$wsfeWsdl=Acquire 'wsfe' $wsfeContract @('FECAESolicitar','FECompUltimoAutorizado','FEParamGetTiposCbte')
$project=Join-Path $output 'Elite.Arca.Wsfe.Generated.csproj';Copy-Item -LiteralPath (Join-Path $assets 'Elite.Arca.Wsfe.Generated.csproj') -Destination $project
$offlineConfig=Join-Path $output 'NuGet.Config'
if($Offline){[IO.File]::WriteAllText($offlineConfig,('<configuration><packageSources><clear /><add key="fixed-local" value="'+[Security.SecurityElement]::Escape($feed)+'" /></packageSources><auditSources><clear /></auditSources></configuration>'),[Text.UTF8Encoding]::new($false))}
$oldPath=$env:PATH;$oldRoll=$env:DOTNET_ROLL_FORWARD;$oldCli=$env:DOTNET_CLI_TELEMETRY_OPTOUT;$oldSvc=$env:DOTNET_SVCUTIL_TELEMETRY_OPTOUT
try{
  $env:PATH="$([IO.Path]::GetDirectoryName($dotnet));$oldPath";$env:DOTNET_ROLL_FORWARD='Major';$env:DOTNET_CLI_TELEMETRY_OPTOUT='1';$env:DOTNET_SVCUTIL_TELEMETRY_OPTOUT='1'
  Copy-Item -LiteralPath (Join-Path $assets 'generated-lock\packages.lock.json') -Destination (Join-Path $output 'packages.lock.json')
  Run $dotnet @('restore',$project,'--locked-mode','--no-cache','--warnaserror','--configfile',$offlineConfig,'-p:NuGetAudit=false')|Out-Null
  function Generate([string]$Service,[string]$Wsdl,[string]$Namespace,[string]$Filename){
    # Use the exact parameter interface emitted by the official svcutil bootstrapper.
    # The already qualified host has the complete dependency graph; do not create
    # another bootstrap executable that drops the explicit NuGet patch references.
    $parameters=[ordered]@{providerId='Microsoft.Tools.ServiceModel.Svcutil-bootstrap';version='8.0.0';options=[ordered]@{
      inputs=@($Wsdl);bootstrapPath=$output;namespaceMappings=@(('*, '+$Namespace));noBootstrapping=$true;
      noLogo=$true;noProjectUpdates=$true;outputDir=(Join-Path $output 'Generated');outputFile=(Join-Path $output ('Generated\'+$Filename));
      projectFile=$project;references=@('Microsoft.Extensions.ObjectPool, {Microsoft.Extensions.ObjectPool, 10.0.0}','System.Security.Cryptography.Pkcs, {System.Security.Cryptography.Pkcs, 10.0.11}','System.Security.Cryptography.Xml, {System.Security.Cryptography.Xml, 10.0.11}');
      serializerMode='XmlSerializer';targetFramework='net10.0';toolContext='Bootstrapper';typeReuseMode='All';verbosity='Minimal'}}
    $parameterPath=Join-Path $output ($Service+'.svcutil.json')
    [IO.File]::WriteAllText($parameterPath,($parameters|ConvertTo-Json -Depth 10),[Text.UTF8Encoding]::new($false))
    Run $dotnet @((Join-Path $svcutil 'dotnet-svcutil.dll'),$parameterPath)|Out-Null
  }
  Push-Location $output
  try{
    Generate 'wsaa' $wsaaWsdl 'Elite.Arca.Wsaa.Generated' 'WsaaReference.cs'
    Generate 'wsfe' $wsfeWsdl 'Elite.Arca.Wsfe.Generated' 'WsfeV1Reference.cs'
  }finally{Pop-Location}
  if((Hash $project)-ne(Hash (Join-Path $assets 'Elite.Arca.Wsfe.Generated.csproj'))){throw 'generator changed the project template'}
  if((Hash (Join-Path $output 'packages.lock.json'))-ne(Hash (Join-Path $assets 'generated-lock\packages.lock.json'))){throw 'locked restore changed the admitted package graph'}
  Run $dotnet @('build',$project,'--configuration','Release','--no-restore','--warnaserror','-p:Deterministic=true','-p:ContinuousIntegrationBuild=true','-p:DebugType=None',('-p:PathMap='+$output+'=/_/arca/fiscal/generated'))|Out-Null

}finally{$env:PATH=$oldPath;$env:DOTNET_ROLL_FORWARD=$oldRoll;$env:DOTNET_CLI_TELEMETRY_OPTOUT=$oldCli;$env:DOTNET_SVCUTIL_TELEMETRY_OPTOUT=$oldSvc}
$wsaaGenerated=Join-Path $output 'Generated\WsaaReference.cs';if((Hash $wsaaGenerated)-ne$wsaaContract.generatedSha256){throw 'generated proxy identity mismatch: WSAA'}
$wsfeGenerated=Join-Path $output 'Generated\WsfeV1Reference.cs';if((Hash $wsfeGenerated)-ne$wsfeContract.generatedSha256){throw 'generated proxy identity mismatch: WSFE'}
$receipt=[ordered]@{schema='elite-arca-generated-client/v3';environment=$Environment;contracts=@([ordered]@{service='wsaa';url=$wsaaContract.url;wsdl_sha256=$wsaaContract.sha256;generated_sha256=$wsaaContract.generatedSha256},[ordered]@{service='wsfev1';url=$wsfeContract.url;wsdl_sha256=$wsfeContract.sha256;generated_sha256=$wsfeContract.generatedSha256});dotnet_sdk=$actualSdk;dotnet_svcutil=$lock.toolchain.dotnetSvcutil;generator_provenance='ADAPTED_RUNTIME_COMPOSITION';generator_inventory_sha256=$adaptation.expected_inventory_sha256;sca_state='NOT_RUN_OFFLINE_RELEASE_SCA_REQUIRED';production_admitted=$false}
[IO.File]::WriteAllText((Join-Path $output 'GENERATION_RECEIPT.json'),($receipt|ConvertTo-Json -Depth 10),[Text.UTF8Encoding]::new($false))
"ARCA_GENERATION_PASS environment=$Environment generated_clients=2 production_admitted=false"
