#requires -Version 7.0
param([Parameter(Mandatory=$true)][string]$ReleaseRoot,[Parameter(Mandatory=$true)][string]$TrustedAllowedSigners,[string]$SshKeygen='ssh-keygen')
$ErrorActionPreference='Stop'
function Fail([string]$Message){throw "SIGNED_RELEASE_VERIFY: $Message"}
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
$root=(Resolve-Path -LiteralPath $ReleaseRoot).Path
$trusted=(Resolve-Path -LiteralPath $TrustedAllowedSigners).Path
if($trusted.StartsWith($root+[IO.Path]::DirectorySeparatorChar,[StringComparison]::OrdinalIgnoreCase)){Fail 'trusted signer policy must be outside the received release'}
foreach($file in @($root,$trusted)){
  $item=Get-Item -LiteralPath $file
  while($null-ne$item){if(($item.Attributes-band[IO.FileAttributes]::ReparsePoint)-ne0){Fail 'release/trusted policy contains a reparse point'};$item=if($item-is[IO.DirectoryInfo]){$item.Parent}else{$item.Directory}}
}
if(-not(Test-Path -LiteralPath $trusted -PathType Leaf)-or(Get-Item -LiteralPath $trusted).Length-gt8192){Fail 'bounded external signer policy required'}
$required=@('artifact.zip','artifact.manifest.json','source.manifest.json','osv.json','SBOM.spdx.json','THIRD_PARTY_NOTICES.txt','allowed_signers','provenance.intoto.json','provenance.intoto.json.sig','RELEASE_RECEIPT.json')
foreach($name in $required){if(-not(Test-Path -LiteralPath (Join-Path $root $name) -PathType Leaf)){Fail "missing $name"}}
$receipt=Get-Content -LiteralPath (Join-Path $root 'RELEASE_RECEIPT.json') -Raw|ConvertFrom-Json -Depth 30 -DateKind String
if($receipt.schema_version -ne 1 -or $receipt.gate -ne 'PORTABLE-SIGNED-RELEASE-EVIDENCE-GATE' -or $receipt.status -ne 'PASS' -or $receipt.reproducibility_runs -ne 2){Fail 'invalid receipt contract'}
if((Sha (Join-Path $root 'artifact.zip')) -cne [string]$receipt.artifact_sha256){Fail 'artifact hash mismatch'}
$map=[ordered]@{artifact_manifest_sha256='artifact.manifest.json';source_manifest_sha256='source.manifest.json';sbom_sha256='SBOM.spdx.json';sca_sha256='osv.json';notices_sha256='THIRD_PARTY_NOTICES.txt';allowed_signers_sha256='allowed_signers';provenance_sha256='provenance.intoto.json';signature_sha256='provenance.intoto.json.sig'}
foreach($e in $map.GetEnumerator()){if((Sha (Join-Path $root $e.Value))-cne[string]$receipt.evidence.($e.Key)){Fail "$($e.Value) hash mismatch"}}
if((Sha $trusted)-cne[string]$receipt.evidence.allowed_signers_sha256){Fail 'release signer differs from externally trusted policy'}
if([string]$receipt.namespace-cne'elite-release-v1'){Fail 'unexpected signature namespace'}
$trustedLines=@(Get-Content -LiteralPath $trusted|Where-Object{$_.Trim()-and-not$_.Trim().StartsWith('#')})
if($trustedLines.Count-ne1-or$trustedLines[0]-notmatch('^'+[regex]::Escape([string]$receipt.release_identity)+'\s+ssh-ed25519\s+[A-Za-z0-9+/=]+(?:\s+.*)?$')){Fail 'external signer identity/policy mismatch'}
$statement=Get-Content -LiteralPath (Join-Path $root 'provenance.intoto.json') -Raw|ConvertFrom-Json -Depth 30 -DateKind String
if($statement._type -ne 'https://in-toto.io/Statement/v1' -or $statement.predicateType -ne 'https://slsa.dev/provenance/v1' -or @($statement.subject).Count -ne 1 -or $statement.subject[0].digest.sha256 -cne [string]$receipt.artifact_sha256){Fail 'provenance subject mismatch'}
$parameters=$statement.predicate.buildDefinition;if($parameters.externalParameters.product_name-cne[string]$receipt.product.name-or$parameters.externalParameters.product_version-cne[string]$receipt.product.version-or$parameters.externalParameters.source_revision-cne[string]$receipt.source_revision-or$parameters.externalParameters.invocation_id-cne[string]$receipt.invocation_id-or$statement.predicate.runDetails.metadata.invocationId-cne[string]$receipt.invocation_id){Fail 'provenance/receipt identity mismatch'}
$internal=$parameters.internalParameters;foreach($pair in @(@('artifact_manifest_sha256','artifact_manifest_sha256'),@('source_manifest_sha256','source_manifest_sha256'),@('sbom_sha256','sbom_sha256'),@('sca_sha256','sca_sha256'),@('notices_sha256','notices_sha256'),@('allowed_signers_sha256','allowed_signers_sha256'))){if([string]$internal.($pair[0])-cne[string]$receipt.evidence.($pair[1])){Fail "signed internal evidence mismatch: $($pair[0])"}};if([string]$internal.artifact_zip_timestamp-cne[string]$receipt.artifact_zip_timestamp){Fail 'signed artifact timestamp mismatch'}
$manifest=Get-Content -LiteralPath (Join-Path $root 'artifact.manifest.json') -Raw|ConvertFrom-Json -Depth 30
if($manifest.artifact_sha256 -cne [string]$receipt.artifact_sha256){Fail 'artifact manifest subject mismatch'}
$sca=Get-Content -LiteralPath (Join-Path $root 'osv.json') -Raw|ConvertFrom-Json -Depth 100;$findings=@(GetVulnerabilityIds $sca);if($findings.Count-ne0){Fail 'stored OSV report contains findings'}
$sbom=Get-Content -LiteralPath (Join-Path $root 'SBOM.spdx.json') -Raw|ConvertFrom-Json -Depth 100;AssertSbomCoverage $sca $sbom
try{$stamp=[datetime]::ParseExact([string]$receipt.artifact_zip_timestamp,'yyyy-MM-ddTHH:mm:ss',[Globalization.CultureInfo]::InvariantCulture,[Globalization.DateTimeStyles]::None)}catch{Fail 'receipt artifact ZIP timestamp invalid'}
Add-Type -AssemblyName System.IO.Compression;$s=[IO.File]::OpenRead((Join-Path $root 'artifact.zip'));try{$z=[IO.Compression.ZipArchive]::new($s,[IO.Compression.ZipArchiveMode]::Read,$true);try{$entries=@($z.Entries|Where-Object{$_.Name -ne ''});if($entries.Count-ne@($manifest.files).Count){Fail 'ZIP entry count mismatch'};$seen=@{};$names=@();foreach($e in $entries){$names+=$e.FullName;if($e.FullName.Contains('\')-or $e.FullName.StartsWith('/')-or $e.FullName.Contains(':')-or @($e.FullName.Split('/')|Where-Object{$_-in@('','.', '..')}).Count-gt0-or$seen.ContainsKey($e.FullName)){Fail 'unsafe or duplicate ZIP entry'};if($e.LastWriteTime.ToString('yyyy-MM-ddTHH:mm:ss')-cne$stamp.ToString('yyyy-MM-ddTHH:mm:ss')){Fail 'ZIP timestamp mismatch'};if($e.CompressedLength-ne$e.Length){Fail 'ZIP entry is not stored deterministically'};$seen[$e.FullName]=1;$expected=@($manifest.files|Where-Object path -ceq $e.FullName);if($expected.Count-ne1){Fail 'ZIP entry absent from manifest'};$stream=$e.Open();try{$hash=[Security.Cryptography.SHA256]::HashData($stream);$actual=[Convert]::ToHexString($hash).ToLowerInvariant()}finally{$stream.Dispose()};if($actual-cne[string]$expected[0].sha256-or$e.Length-ne[long]$expected[0].size){Fail 'ZIP entry digest or size mismatch'}};$sorted=@($names|Sort-Object);if(($names-join"`n")-cne($sorted-join"`n")){Fail 'ZIP entries are not sorted'}}finally{$z.Dispose()}}finally{$s.Dispose()}
$ssh=(Get-Command $SshKeygen -ErrorAction Stop).Source;if((Sha $ssh)-cne[string]$receipt.tools.ssh_keygen.sha256){Fail 'ssh-keygen identity changed'};$auth=Get-AuthenticodeSignature -LiteralPath $ssh;if($auth.Status-ne'Valid'-or$auth.SignerCertificate.Subject-notmatch'O=Microsoft Corporation'){Fail 'ssh-keygen Authenticode invalid'}
$psi=[Diagnostics.ProcessStartInfo]::new();$psi.FileName=$ssh;foreach($a in @('-Y','verify','-f',$trusted,'-I',[string]$receipt.release_identity,'-n','elite-release-v1','-s',(Join-Path $root 'provenance.intoto.json.sig'))){[void]$psi.ArgumentList.Add($a)};$psi.UseShellExecute=$false;$psi.RedirectStandardInput=$true;$psi.RedirectStandardOutput=$true;$psi.RedirectStandardError=$true;$p=[Diagnostics.Process]::Start($psi);$bytes=[IO.File]::ReadAllBytes((Join-Path $root 'provenance.intoto.json'));$p.StandardInput.BaseStream.Write($bytes,0,$bytes.Length);$p.StandardInput.Close();$p.WaitForExit();if($p.ExitCode-ne0){Fail 'OpenSSH signature rejected'}
Write-Output "PORTABLE_SIGNED_RELEASE_VERIFY_PASS artifact_sha256=$($receipt.artifact_sha256)"
