#requires -Version 7.0
param(
  [Parameter(Mandatory=$true)][string]$ProjectRoot,
  [Parameter(Mandatory=$true)][string]$Profile,
  [Parameter(Mandatory=$true)][string]$ReleaseRoot,
  [Parameter(Mandatory=$true)][string]$OsvScanner,
  [Parameter(Mandatory=$true)][string]$SigningKey,
  [string]$SshKeygen = 'ssh-keygen'
)
$ErrorActionPreference='Stop'
$Utf8=[Text.UTF8Encoding]::new($false)
$OsvSha='25e42f5ef6711fd8c0fb45390972205891dd44c6bd02ac93f0f63e8e98d9bfb6'
$Namespace='elite-release-v1'
function Fail([string]$Message){throw "SIGNED_RELEASE_GATE: $Message"}
function GetVulnerabilityIds($Report){
  # Missing/empty vulnerabilities is the official clean-package representation.
  # foreach over $null has zero iterations; @($null) would fabricate one item.
  $ids=[Collections.Generic.List[string]]::new()
  foreach($result in $Report.results){foreach($pkg in $result.packages){foreach($v in $pkg.vulnerabilities){
    if($null-eq$v-or[string]::IsNullOrWhiteSpace([string]$v.id)){Fail 'malformed vulnerability record'}
    $ids.Add([string]$v.id)
  }}}
  $ids.ToArray()
}
function PackageKey([string]$Type,[string]$Name,[string]$Version){
  $Type+"`0"+$Name+"`0"+$Version
}
function AssertSbomCoverage($Report,$Sbom){
  if($Sbom.spdxVersion-cne'SPDX-2.3'){Fail 'SBOM must be SPDX 2.3'}
  $types=@{Go='golang';npm='npm';NuGet='nuget';PyPI='pypi'}
  $actual=[Collections.Generic.HashSet[string]]::new([StringComparer]::Ordinal)
  foreach($p in $Sbom.packages){
    foreach($r in $p.externalRefs){
      if($r.referenceType-ceq'purl'){
        $locator=[uri]::UnescapeDataString([string]$r.referenceLocator)
        if($locator-cmatch'^pkg:([^/]+)/(.+)@([^@?#]+)$'){
          [void]$actual.Add((PackageKey $Matches[1] $Matches[2] $Matches[3]))
        }
      }
    }
    if([string]$p.sourceInfo-cmatch'^OSV-UNVERSIONED:(Go|npm|NuGet|PyPI)$'-and[string]::IsNullOrEmpty([string]$p.versionInfo)){
      [void]$actual.Add((PackageKey $types[$Matches[1]] ([string]$p.name) ''))
    }
  }
  $count=0
  foreach($r in $Report.results){foreach($row in $r.packages){
    $p=$row.package;$count++
    if(-not$types.ContainsKey([string]$p.ecosystem)-or[string]::IsNullOrWhiteSpace([string]$p.name)){Fail 'unsupported or malformed OSV package identity'}
    if(-not$actual.Contains((PackageKey $types[[string]$p.ecosystem] ([string]$p.name) ([string]$p.version)))){Fail ('SBOM missing OSV package: '+$p.ecosystem+'/'+$p.name+'@'+$p.version)}
  }}
  if($count-eq0){Fail 'OSV dependency inventory is empty'}
}
function AddUnversionedSbomRecords($Report,$Sbom){
  # OSV SPDX omits unversioned source records. Preserve their exact reported
  # identity as AUTHORED metadata glue; do not invent versions or provenance.
  $seen=[Collections.Generic.HashSet[string]]::new([StringComparer]::Ordinal)
  foreach($r in $Report.results){foreach($row in $r.packages){
    $p=$row.package
    if([string]::IsNullOrEmpty([string]$p.version)){
      $key=PackageKey ([string]$p.ecosystem) ([string]$p.name) ''
      if($seen.Add($key)){
        $hash=[Convert]::ToHexString([Security.Cryptography.SHA256]::HashData([Text.Encoding]::UTF8.GetBytes($key))).ToLowerInvariant()
        $id='SPDXRef-OSV-Unversioned-'+$hash
        $Sbom.packages+=@([pscustomobject]@{name=[string]$p.name;SPDXID=$id;downloadLocation='NOASSERTION';filesAnalyzed=$false;sourceInfo=('OSV-UNVERSIONED:'+ $p.ecosystem);comment='Unversioned record from the signed OSV source report; no third-party version, license, or origin inferred.'})
        $Sbom.relationships+=@([pscustomobject]@{spdxElementId=$Sbom.SPDXID;relationshipType='DESCRIBES';relatedSpdxElement=$id})
      }
    }
  }}
}

function Sha([string]$Path){(Get-FileHash -LiteralPath $Path -Algorithm SHA256).Hash.ToLowerInvariant()}
function WriteJson([string]$Path,$Value){[IO.File]::WriteAllText($Path,(ConvertTo-Json $Value -Depth 30 -Compress)+"`n",$Utf8)}
function ResolveProjectFile([string]$Root,[string]$Ref,[string]$Label){
  if([string]::IsNullOrWhiteSpace($Ref) -or [IO.Path]::IsPathRooted($Ref) -or $Ref -match '(^|[\\/])\.\.([\\/]|$)'){Fail "$Label must be a safe project-relative file"}
  $p=[IO.Path]::GetFullPath((Join-Path $Root $Ref))
  if(-not $p.StartsWith($Root+[IO.Path]::DirectorySeparatorChar,[StringComparison]::OrdinalIgnoreCase) -or -not(Test-Path -LiteralPath $p -PathType Leaf)){Fail "$Label missing or escaped project"}
  $item=Get-Item -LiteralPath $p
  while($null-ne$item){
    if(($item.Attributes-band[IO.FileAttributes]::ReparsePoint)-ne0){Fail "$Label contains a reparse point"}
    if($item.FullName-ceq$Root){break}
    $item=if($item-is[IO.DirectoryInfo]){$item.Parent}else{$item.Directory}
  }
  $p
}
function AssertSourceUnchanged([string]$Root,[object[]]$Entries){
  foreach($entry in $Entries){
    $p=ResolveProjectFile $Root ([string]$entry.path) 'frozen source input'
    if((Sha $p)-cne[string]$entry.sha256-or(Get-Item -LiteralPath $p).Length-ne[long]$entry.size){Fail ('source input changed: '+$entry.path)}
  }
}
function InvokeChecked([string]$Exe,[string[]]$Arguments,[string]$Label){
  & $Exe @Arguments
  if($LASTEXITCODE -ne 0){Fail "$Label failed with exit $LASTEXITCODE"}
}
function InvokeSshVerify([string]$Exe,[string]$Allowed,[string]$Identity,[string]$Signature,[string]$Payload){
  $psi=[Diagnostics.ProcessStartInfo]::new();$psi.FileName=$Exe
  foreach($a in @('-Y','verify','-f',$Allowed,'-I',$Identity,'-n',$Namespace,'-s',$Signature)){[void]$psi.ArgumentList.Add($a)}
  $psi.UseShellExecute=$false;$psi.RedirectStandardInput=$true;$psi.RedirectStandardOutput=$true;$psi.RedirectStandardError=$true
  $p=[Diagnostics.Process]::Start($psi);$bytes=[IO.File]::ReadAllBytes($Payload);$p.StandardInput.BaseStream.Write($bytes,0,$bytes.Length);$p.StandardInput.Close();$p.WaitForExit()
  if($p.ExitCode -ne 0){Fail ("signature verification failed: "+$p.StandardError.ReadToEnd().Trim())}
}
function WriteDeterministicZip([string]$InputRoot,[string]$ZipPath,[datetimeoffset]$Stamp){
  Add-Type -AssemblyName System.IO.Compression
  foreach($node in @(Get-Item -LiteralPath $InputRoot)+@(Get-ChildItem -LiteralPath $InputRoot -Recurse -Force)){
    if(($node.Attributes-band[IO.FileAttributes]::ReparsePoint)-ne0){Fail 'build tree contains a reparse point'}
  }
  $files=@(Get-ChildItem -LiteralPath $InputRoot -File -Recurse | Sort-Object {$_.FullName.Substring($InputRoot.Length+1).Replace('\','/')})
  if($files.Count -lt 1){Fail 'build produced no files'}
  $stream=[IO.File]::Open($ZipPath,[IO.FileMode]::CreateNew,[IO.FileAccess]::ReadWrite,[IO.FileShare]::None)
  try{$zip=[IO.Compression.ZipArchive]::new($stream,[IO.Compression.ZipArchiveMode]::Create,$true)
    try{foreach($f in $files){if(($f.Attributes -band [IO.FileAttributes]::ReparsePoint) -ne 0){Fail 'build output contains a reparse point'};$rel=$f.FullName.Substring($InputRoot.Length+1).Replace('\','/');$e=$zip.CreateEntry($rel,[IO.Compression.CompressionLevel]::NoCompression);$e.LastWriteTime=$Stamp;$src=[IO.File]::OpenRead($f.FullName);try{$dst=$e.Open();try{$src.CopyTo($dst)}finally{$dst.Dispose()}}finally{$src.Dispose()}}}finally{$zip.Dispose()}
  }finally{$stream.Dispose()}
}

$root=(Resolve-Path -LiteralPath $ProjectRoot).Path
$profilePath=(Resolve-Path -LiteralPath $Profile).Path
if(Test-Path -LiteralPath $ReleaseRoot){Fail 'release root already exists'}
$profileSha=Sha $profilePath
$cfg=Get-Content -LiteralPath $profilePath -Raw | ConvertFrom-Json -Depth 30 -DateKind String
if($cfg.schema_version -ne 1 -or $cfg.enabled -ne $true){Fail 'profile is not explicitly enabled'}
foreach($field in @('product_name','product_version','source_uri','source_revision','invocation_id','builder_id','build_script_ref','build_script_sha256','source_inputs_ref','notices_ref','release_identity','allowed_signers_ref')){if([string]::IsNullOrWhiteSpace([string]$cfg.$field)){Fail "$field is required"}}
if($cfg.product_name -notmatch '^[A-Za-z0-9][A-Za-z0-9._-]{0,127}$' -or $cfg.product_version -notmatch '^[A-Za-z0-9][A-Za-z0-9._+-]{0,127}$'){Fail 'unsafe product identity'}
if($cfg.release_identity -notmatch '^[A-Za-z0-9][A-Za-z0-9@._+-]{0,127}$'){Fail 'unsafe release identity'}
foreach($uriField in @('source_uri','builder_id')){try{if(-not([uri]$cfg.$uriField).IsAbsoluteUri){Fail "$uriField must be absolute"}}catch{Fail "$uriField must be an absolute URI"}}
try{$stamp=[datetimeoffset]::ParseExact([string]$cfg.artifact_zip_timestamp,'yyyy-MM-ddTHH:mm:ss',[Globalization.CultureInfo]::InvariantCulture,[Globalization.DateTimeStyles]::AssumeUniversal)}catch{Fail 'invalid artifact ZIP timestamp'}
if($stamp.Year -lt 1980){Fail 'ZIP timestamp must be 1980 or later'}
$buildScript=ResolveProjectFile $root ([string]$cfg.build_script_ref) 'build_script_ref'
if($cfg.build_script_sha256 -notmatch '^[0-9a-f]{64}$' -or (Sha $buildScript) -cne [string]$cfg.build_script_sha256){Fail 'build script hash mismatch'}
$inputsFile=ResolveProjectFile $root ([string]$cfg.source_inputs_ref) 'source_inputs_ref'
$notices=ResolveProjectFile $root ([string]$cfg.notices_ref) 'notices_ref'
$allowed=ResolveProjectFile $root ([string]$cfg.allowed_signers_ref) 'allowed_signers_ref'
$inputRefs=@(Get-Content -LiteralPath $inputsFile | ForEach-Object {$_.Trim()} | Where-Object {$_ -and -not $_.StartsWith('#')})
if($inputRefs.Count -lt 1 -or @($inputRefs|Sort-Object -Unique).Count -ne $inputRefs.Count){Fail 'source input manifest is empty or duplicated'}
$sourceEntries=@();foreach($ref in ($inputRefs|Sort-Object)){ $p=ResolveProjectFile $root $ref 'source input';$sourceEntries += [ordered]@{path=$ref.Replace('\','/');sha256=Sha $p;size=(Get-Item -LiteralPath $p).Length} }
$inputSet=@{};$inputRefs|ForEach-Object{$inputSet[$_.Replace('\','/')]=1}
foreach($required in @([string]$cfg.build_script_ref,[string]$cfg.source_inputs_ref,[string]$cfg.notices_ref,[string]$cfg.allowed_signers_ref)+@($cfg.dependency_manifest_refs)){if(-not $inputSet.ContainsKey($required.Replace('\','/'))){Fail "required source input not listed: $required"};[void](ResolveProjectFile $root $required 'required source input')}
$dependencyRefs=@($cfg.dependency_manifest_refs)
if($dependencyRefs.Count-lt1-or@($dependencyRefs|Sort-Object -Unique).Count-ne$dependencyRefs.Count){Fail 'explicit unique dependency manifests required'}
$scanArguments=@('scan','source')
foreach($ref in $dependencyRefs){$scanArguments+=@('--lockfile',(ResolveProjectFile $root ([string]$ref) 'dependency manifest'))}
$osv=(Resolve-Path -LiteralPath $OsvScanner).Path;if((Sha $osv)-cne $OsvSha){Fail 'OSV Scanner hash mismatch'}
$osvVersion=(& $osv --version 2>&1|Out-String).Trim();if($LASTEXITCODE -ne 0 -or $osvVersion -notmatch '2\.5\.1'){Fail 'OSV Scanner version mismatch'}
$ssh=(Get-Command $SshKeygen -ErrorAction Stop).Source
if($cfg.require_microsoft_authenticode -eq $true){$auth=Get-AuthenticodeSignature -LiteralPath $ssh;if($auth.Status -ne 'Valid' -or $auth.SignerCertificate.Subject -notmatch 'O=Microsoft Corporation'){Fail 'ssh-keygen is not Authenticode-valid Microsoft code'}}
$key=(Resolve-Path -LiteralPath $SigningKey).Path
$target=[IO.Path]::GetFullPath($ReleaseRoot);$parent=Split-Path -Parent $target;if(-not(Test-Path -LiteralPath $parent)){Fail 'release parent missing'}
if($target.StartsWith($root+[IO.Path]::DirectorySeparatorChar,[StringComparison]::OrdinalIgnoreCase)){Fail 'release root must be outside project root'}
if($key.StartsWith($root+[IO.Path]::DirectorySeparatorChar,[StringComparison]::OrdinalIgnoreCase)){Fail 'signing key must be outside project root'}
$allowedLines=@(Get-Content -LiteralPath $allowed|Where-Object{$_.Trim() -and -not $_.Trim().StartsWith('#')});if($allowedLines.Count-ne1-or$allowedLines[0]-notmatch ('^'+[regex]::Escape([string]$cfg.release_identity)+'\s+ssh-ed25519\s+[A-Za-z0-9+/=]+(?:\s+.*)?$')){Fail 'allowed signers must contain exactly the release identity and one Ed25519 public key'}
$work=Join-Path $parent ('.elite-release-staging-'+[guid]::NewGuid().ToString('N'));New-Item -ItemType Directory -Path $work|Out-Null
$completed=$false
try{
  $out1=Join-Path $work 'build1';$out2=Join-Path $work 'build2';New-Item -ItemType Directory -Path $out1,$out2|Out-Null
  $started=[datetimeoffset]::UtcNow
  foreach($out in @($out1,$out2)){AssertSourceUnchanged $root $sourceEntries;InvokeChecked -Exe 'pwsh' -Arguments @('-NoProfile','-File',$buildScript,'-ProjectRoot',$root,'-OutputRoot',$out,'-SourceDateEpoch',[string]$stamp.ToUnixTimeSeconds()) -Label 'project build';AssertSourceUnchanged $root $sourceEntries}
  $zip1=Join-Path $work 'artifact1.zip';$zip2=Join-Path $work 'artifact2.zip';WriteDeterministicZip $out1 $zip1 $stamp;WriteDeterministicZip $out2 $zip2 $stamp
  $artifactSha=Sha $zip1;if((Sha $zip2)-cne $artifactSha){Fail 'two clean builds are not byte-identical'}
  $artifactEntries=@(Get-ChildItem -LiteralPath $out1 -File -Recurse|Sort-Object {$_.FullName.Substring($out1.Length+1).Replace('\','/')}|ForEach-Object{[ordered]@{path=$_.FullName.Substring($out1.Length+1).Replace('\','/');sha256=Sha $_.FullName;size=$_.Length}})
  $sourceManifest=Join-Path $work 'source.manifest.json';WriteJson $sourceManifest ([ordered]@{schema_version=1;files=$sourceEntries})
  $artifactManifest=Join-Path $work 'artifact.manifest.json';WriteJson $artifactManifest ([ordered]@{schema_version=1;artifact_sha256=$artifactSha;files=$artifactEntries})
  $sca=Join-Path $work 'osv.json';& $osv @scanArguments --all-packages --format json --output-file $sca;$scanExit=$LASTEXITCODE
  if(-not(Test-Path -LiteralPath $sca)){Fail 'OSV JSON report missing'};$scan=Get-Content -LiteralPath $sca -Raw|ConvertFrom-Json -Depth 100;$findings=@(GetVulnerabilityIds $scan);if($scanExit -ne 0 -or $findings.Count -ne 0){Fail "OSV found $($findings.Count) vulnerability records"}
  $sbom=Join-Path $work 'SBOM.spdx.json';InvokeChecked -Exe $osv -Arguments ($scanArguments+@('--all-packages','--format','spdx-2-3','--output-file',$sbom)) -Label 'SPDX generation'
  $document=Get-Content -LiteralPath $sbom -Raw|ConvertFrom-Json -Depth 100;AddUnversionedSbomRecords $scan $document;AssertSbomCoverage $scan $document;WriteJson $sbom $document
  $release=Join-Path $work 'release';New-Item -ItemType Directory -Path $release|Out-Null
  Copy-Item -LiteralPath $zip1 -Destination (Join-Path $release 'artifact.zip');Copy-Item -LiteralPath $artifactManifest -Destination $release;Copy-Item -LiteralPath $sourceManifest -Destination $release;Copy-Item -LiteralPath $sca -Destination $release;Copy-Item -LiteralPath $sbom -Destination $release;Copy-Item -LiteralPath $notices -Destination (Join-Path $release 'THIRD_PARTY_NOTICES.txt');Copy-Item -LiteralPath $allowed -Destination (Join-Path $release 'allowed_signers')
  AssertSourceUnchanged $root $sourceEntries
  if((Sha $profilePath)-cne$profileSha){Fail 'release profile changed during build'}
  $finished=[datetimeoffset]::UtcNow
  $statement=[ordered]@{_type='https://in-toto.io/Statement/v1';subject=@([ordered]@{name="$($cfg.product_name)-$($cfg.product_version).zip";digest=[ordered]@{sha256=$artifactSha}});predicateType='https://slsa.dev/provenance/v1';predicate=[ordered]@{buildDefinition=[ordered]@{buildType='https://slsa.dev/provenance/v1';externalParameters=[ordered]@{product_name=$cfg.product_name;product_version=$cfg.product_version;source_uri=$cfg.source_uri;source_revision=$cfg.source_revision;invocation_id=$cfg.invocation_id};internalParameters=[ordered]@{profile_sha256=$profileSha;source_manifest_sha256=Sha $sourceManifest;artifact_manifest_sha256=Sha $artifactManifest;sbom_sha256=Sha $sbom;sca_sha256=Sha $sca;notices_sha256=Sha $notices;allowed_signers_sha256=Sha $allowed;artifact_zip_timestamp=$stamp.ToString('yyyy-MM-ddTHH:mm:ss')};resolvedDependencies=@($sourceEntries|ForEach-Object{[ordered]@{uri=([uri]::new([uri]$cfg.source_uri,$_.path)).AbsoluteUri;digest=[ordered]@{sha256=$_.sha256}}})};runDetails=[ordered]@{builder=[ordered]@{id=$cfg.builder_id};metadata=[ordered]@{invocationId=$cfg.invocation_id;startedOn=$started.ToString('O');finishedOn=$finished.ToString('O')}}}}
  $statementPath=Join-Path $release 'provenance.intoto.json';WriteJson $statementPath $statement
  InvokeChecked -Exe $ssh -Arguments @('-Y','sign','-f',$key,'-n',$Namespace,$statementPath) -Label 'provenance signing'
  $signature=$statementPath+'.sig';if(-not(Test-Path -LiteralPath $signature)){Fail 'signature file missing'}
  InvokeSshVerify $ssh (Join-Path $release 'allowed_signers') ([string]$cfg.release_identity) $signature $statementPath
  $authStatus=Get-AuthenticodeSignature -LiteralPath $ssh
  $receipt=[ordered]@{schema_version=1;gate='PORTABLE-SIGNED-RELEASE-EVIDENCE-GATE';status='PASS';product=[ordered]@{name=$cfg.product_name;version=$cfg.product_version};artifact_sha256=$artifactSha;artifact_zip_timestamp=$stamp.ToString('yyyy-MM-ddTHH:mm:ss');reproducibility_runs=2;source_revision=$cfg.source_revision;invocation_id=$cfg.invocation_id;release_identity=$cfg.release_identity;namespace=$Namespace;tools=[ordered]@{osv=[ordered]@{version=$osvVersion;sha256=Sha $osv};ssh_keygen=[ordered]@{version=(Get-Item -LiteralPath $ssh).VersionInfo.ProductVersion;sha256=Sha $ssh;authenticode_status=[string]$authStatus.Status;signer_subject=$authStatus.SignerCertificate.Subject}};evidence=[ordered]@{artifact_manifest_sha256=Sha (Join-Path $release 'artifact.manifest.json');source_manifest_sha256=Sha (Join-Path $release 'source.manifest.json');sbom_sha256=Sha (Join-Path $release 'SBOM.spdx.json');sca_sha256=Sha (Join-Path $release 'osv.json');notices_sha256=Sha (Join-Path $release 'THIRD_PARTY_NOTICES.txt');allowed_signers_sha256=Sha (Join-Path $release 'allowed_signers');provenance_sha256=Sha $statementPath;signature_sha256=Sha $signature}}
  WriteJson (Join-Path $release 'RELEASE_RECEIPT.json') $receipt
  Move-Item -LiteralPath $release -Destination $target
  $completed=$true
  Write-Output "PORTABLE_SIGNED_RELEASE_PASS artifact_sha256=$artifactSha findings=0 reproducibility_runs=2"
}finally{if(-not$completed){Write-Warning ("FAILED_RELEASE_EVIDENCE_RETAINED="+$work)}elseif(Test-Path -LiteralPath $work){$resolved=[IO.Path]::GetFullPath($work);if($resolved.StartsWith([IO.Path]::GetFullPath($parent)+[IO.Path]::DirectorySeparatorChar+'.elite-release-staging-',[StringComparison]::OrdinalIgnoreCase)){Remove-Item -LiteralPath $resolved -Recurse -Force}}}
