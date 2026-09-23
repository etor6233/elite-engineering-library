#requires -Version 7.0
# AUTHORED orchestration for the admitted local builder; no provider credentials.
param(
 [Parameter(Mandatory=$true)][string]$ProjectRoot,
 [Parameter(Mandatory=$true)][string]$OutputRoot,
 [Parameter(Mandatory=$true)][long]$SourceDateEpoch
)
Set-StrictMode -Version Latest
$ErrorActionPreference='Stop'
$utf8=[Text.UTF8Encoding]::new($false)
function Fail([string]$message){throw ('LOCAL_SIGNED_BUILD: '+$message)}
function Plain([string]$path){
 $full=[IO.Path]::GetFullPath($path).TrimEnd([IO.Path]::DirectorySeparatorChar)
 $cursor=$full
 while($cursor){
  $item=Get-Item -LiteralPath $cursor -Force -ErrorAction SilentlyContinue
  if($null-ne$item-and($item.Attributes-band[IO.FileAttributes]::ReparsePoint)-ne0){Fail 'Reparse path rejected'}
  $cursor=Split-Path -Parent $cursor
 }
 $full
}
function Sha([string]$path){(Get-FileHash -LiteralPath $path -Algorithm SHA256).Hash.ToLowerInvariant()}
function WriteJson([string]$path,$value){[IO.File]::WriteAllText($path,($value|ConvertTo-Json -Depth 10)+"`n",$utf8)}
function Disjoint([string]$a,[string]$b){
 if($a-eq$b-or$a.StartsWith($b+[IO.Path]::DirectorySeparatorChar,[StringComparison]::OrdinalIgnoreCase)-or$b.StartsWith($a+[IO.Path]::DirectorySeparatorChar,[StringComparison]::OrdinalIgnoreCase)){Fail 'Overlapping build paths'}
}
$root=Plain $ProjectRoot;$output=Plain $OutputRoot
Disjoint $root $output
if(-not(Test-Path -LiteralPath $root -PathType Container)-or-not(Test-Path -LiteralPath $output -PathType Container)-or@(Get-ChildItem -LiteralPath $output -Force).Count-ne0){Fail 'Existing source and empty output directories required'}
$configPath=Join-Path $root '.elite-local-build-inputs.json'
[void](Plain $configPath)
$configSha=Sha $configPath
$cfg=Get-Content -LiteralPath $configPath -Raw|ConvertFrom-Json -AsHashtable
$expected=@('schema_version','scope','python','source_inventory_ref','source_inventory_sha256','go','node','go_mod_cache','install_inputs','install_inputs_sha256','arca_inputs','arca_inputs_sha256','run_root')
if(@(Compare-Object @($cfg.Keys|Sort-Object) @($expected|Sort-Object)).Count-ne0-or$cfg.schema_version-ne1-or$cfg.scope-cne'LOCAL_FIXTURES'){Fail 'Exact local input schema required'}
if($SourceDateEpoch-ne315532800){Fail 'Reference ZIP timestamp must be 1980-01-01 UTC'}
if(@(Compare-Object @($cfg.python.Keys|Sort-Object) @('path','sha256')).Count-ne0){Fail 'Exact Python pin required'}
$python=Plain ([string]$cfg.python.path)
if($cfg.python.sha256-notmatch'^[0-9a-f]{64}$'-or(Sha $python)-cne$cfg.python.sha256){Fail 'Python pin changed'}
$inventoryRef=[string]$cfg.source_inventory_ref
if([IO.Path]::IsPathRooted($inventoryRef)-or$inventoryRef-match'[:\\]' -or @($inventoryRef.Split('/')|Where-Object{$_-in@('','..','.')}).Count){Fail 'Safe relative inventory required'}
$inventory=Plain (Join-Path $root $inventoryRef)
if($cfg.source_inventory_sha256-notmatch'^[0-9a-f]{64}$'-or(Sha $inventory)-cne$cfg.source_inventory_sha256){Fail 'Source inventory pin changed'}
$run=Plain ([string]$cfg.run_root)
Disjoint $run $root;Disjoint $run $output
if(-not(Test-Path -LiteralPath (Split-Path -Parent $run) -PathType Container)){Fail 'Run parent absent'}
$markerPath=Join-Path $run 'session.json'
if(-not(Test-Path -LiteralPath $run)){
 [void](New-Item -ItemType Directory -Path $run)
 WriteJson $markerPath @{schema_version=1;configuration_sha256=$configSha;source_inventory_sha256=$cfg.source_inventory_sha256;completed_builds=0}
}
[void](Plain $markerPath);[void](Plain (Join-Path $run 'session.lock'))
$lock=[IO.File]::Open((Join-Path $run 'session.lock'),[IO.FileMode]::OpenOrCreate,[IO.FileAccess]::ReadWrite,[IO.FileShare]::None)
try{
 $marker=Get-Content -LiteralPath $markerPath -Raw|ConvertFrom-Json -AsHashtable
 if($marker.schema_version-ne1-or$marker.configuration_sha256-cne$configSha-or$marker.source_inventory_sha256-cne$cfg.source_inventory_sha256-or$marker.completed_builds-notin@(0,1)){Fail 'Run belongs to another configuration or already completed both builds'}
 $number=[int]$marker.completed_builds+1
 $workspace=Plain (Join-Path $run 'workspace')
 $archive=Plain (Join-Path $run ('workspace-'+$number))
 $build=Plain (Join-Path $run ('build-'+$number))
 foreach($path in @($workspace,$archive,$build)){if(Test-Path -LiteralPath $path){Fail 'Previous or partial build retained; use a new run_root after diagnosing it'}}
 $arguments=@('-X','utf8','-B',(Join-Path $root 'ci/build_local_reference.py'),'--source',$root,'--target',$build,'--workspace',$workspace,'--inventory',$inventory,'--inventory-sha256',[string]$cfg.source_inventory_sha256,'--go',[string]$cfg.go,'--node',[string]$cfg.node,'--go-mod-cache',[string]$cfg.go_mod_cache,'--install-inputs',[string]$cfg.install_inputs,'--install-inputs-sha256',[string]$cfg.install_inputs_sha256,'--arca-inputs',[string]$cfg.arca_inputs,'--arca-inputs-sha256',[string]$cfg.arca_inputs_sha256)
 & $python @arguments
 if($LASTEXITCODE-ne0){Fail 'Actual independent build failed; all working evidence retained'}
 if((Sha $configPath)-cne$configSha-or(Sha $inventory)-cne$cfg.source_inventory_sha256){Fail 'Build inputs changed'}
 & $python -X utf8 -B (Join-Path $root 'ci/publish_reference_build.py') --build $build --output $output --inventory $inventory --inventory-sha256 ([string]$cfg.source_inventory_sha256)
 if($LASTEXITCODE-ne0){Fail 'Artifact/receipt correspondence failed'}
 # Both names were resolved and constrained to the exact owned run directory.
 if((Split-Path -Parent (Plain $workspace))-cne$run-or(Split-Path -Parent (Plain $archive))-cne$run-or(Test-Path -LiteralPath $archive)){Fail 'Unsafe workspace archive target'}
 Move-Item -LiteralPath $workspace -Destination $archive
 $marker.completed_builds=$number
 WriteJson $markerPath $marker
 Write-Output ('LOCAL_SIGNED_BUILD_PASS independent_build='+$number+' evidence='+$build)
}finally{$lock.Dispose()}
