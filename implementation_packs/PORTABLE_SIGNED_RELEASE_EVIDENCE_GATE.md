# Portable Signed Release Evidence Gate

## 1. Metadata

```yaml
pack_id: "PORTABLE-SIGNED-RELEASE-EVIDENCE-GATE"
pack_version: "0.1.0"
status:
  authority: ELITE_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un candidate release local sin Git ni CI hospedado: dos builds byte-idénticos, ZIP determinista, manifiestos, OSV 2.5.1 sin hallazgos, SPDX 2.3, provenance in-toto/SLSA y firma Ed25519 OpenSSH verificada contra una política pública explícita."
stacks: ["PowerShell 7+", "Google OSV-Scanner 2.5.1", "Microsoft Windows OpenSSH sshsig", "SPDX 2.3", "in-toto Statement v1", "SLSA provenance v1"]
compatible_with: ["proyecto con build PowerShell hash-locked", "Windows con OpenSSH firmado por Microsoft", "OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x"]
incompatible_with: ["clave dentro del proyecto", "release dentro del proyecto", "build no determinista", "dependencias vulnerables", "firmante implícito", "producción sin gates del target"]
license_expression: "LicenseRef-Workspace-Owner AND LicenseRef-OpenSSH"
upstream_sources: ["https://github.com/slsa-framework/slsa/tree/19e4e2f005f871270c4f555fc47afecfb37f3efe/docs/spec/v1.2", "https://github.com/in-toto/attestation/tree/df02077bf97218a8860a5c534eff1f1381f56984/spec/v1.2", "https://github.com/openssh/openssh-portable/tree/d01efaa1c9ed84fd9011201dbc3c7cb0a82bcee3", "https://github.com/google/osv-scanner/tree/c84fa4568f2526d0333e9a914ea8a0a5f74ad68b", "https://learn.microsoft.com/windows-server/administration/openssh/openssh-overview"]
verified_at: "2026-09-05"
```

## 2. Applicability

Usar al cerrar un artefacto distribuible en Windows sin depender de Git, GitHub Actions, una cuenta cloud ni minutos pagos. El proyecto debe aportar una clave Ed25519 preexistente y protegida fuera del árbol, una política pública de firmantes, un inventario explícito de inputs, manifiestos de dependencias, notices y un build script determinista. No usar el test fixture como identidad productiva ni interpretar este PASS como deployment, recuperación o aceptación empresarial.

## 3. Architecture contract

- El build script está fijado por SHA-256 y se ejecuta dos veces con el mismo `SOURCE_DATE_EPOCH` explícito.
- El artefacto usa entradas ordenadas, timestamp ZIP sin zona y `NoCompression`; ambos ZIP deben ser byte-idénticos.
- Google OSV-Scanner es el binario exacto 2.5.1 ya fijado por hash en la biblioteca; cualquier finding bloquea.
- La salida conserva SCA JSON, SPDX 2.3, source/artifact manifests, notices y allowed signers.
- La provenance es un in-toto Statement v1 con predicate SLSA provenance v1 y hashes de toda la evidencia material.
- OpenSSH `sshsig` firma la provenance en namespace `elite-release-v1`; la política permite exactamente una identidad Ed25519 y se verifica antes de publicar.
- La clave privada nunca se copia ni se registra; debe vivir fuera del proyecto. `ssh-keygen` debe portar Authenticode Microsoft válido en este perfil Windows.
- El release se construye en staging hermano y se mueve al destino sólo después de todos los gates.
- El verificador independiente rechaza cualquier cambio de artefacto, entrada ZIP, timestamp, manifest, SBOM, SCA, provenance, firma, signer policy o receipt.

## 4. Exact file manifest

```text
CREATE portable_release_evidence_gate/README.md
CREATE portable_release_evidence_gate/release-profile.template.json
CREATE portable_release_evidence_gate/build_release.ps1
CREATE portable_release_evidence_gate/verify_release.ps1
CREATE portable_release_evidence_gate/verify_pack.ps1
CREATE portable_release_evidence_gate/fixtures/openssh-ed25519.pub.b64
CREATE portable_release_evidence_gate/fixtures/openssh-ed25519.sig.b64
CREATE portable_release_evidence_gate/fixtures/openssh-namespace.b64
CREATE portable_release_evidence_gate/fixtures/openssh-signed-data.b64
CREATE portable_release_evidence_gate/OPENSSH_FIXTURE_NOTICE.md
```

## 5. Materialization blocks

### FILE: `portable_release_evidence_gate/README.md`
```yaml
block_id: "PORTABLE-SIGNED-RELEASE-EVIDENCE-GATE:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "local, governed by SLSA v1.2, in-toto v1.2, OpenSSH V_10_2_P1, Microsoft Windows OpenSSH and Google OSV 2.5.1"
license: "LicenseRef-Workspace-Owner"
sha256: "dc84cab074dcf39641eb5a98d6c6a74b4ddf273640e30714065b8018930766c0"
variables: []
secrets_allowed: false
```
````markdown
# Portable signed release evidence gate

This gate produces a release candidate without Git, a hosted CI service, or paid signing infrastructure. It executes one hash-locked project build twice, requires byte-identical artifacts, generates a source manifest, scans the source with the exact Google OSV-Scanner 2.5.1, emits SPDX 2.3, creates an in-toto Statement v1 with the SLSA provenance v1 predicate, signs that statement through OpenSSH `sshsig`, immediately verifies the signature against an explicit allowed-signers file, and writes a hash-linked receipt.

The project supplies the business code, an existing release signing key, the public allowed-signers policy, an explicit input manifest, dependency manifests, notices, and the build script. The private key is never copied into the release directory or receipt. A pass proves only the recorded local candidate; it does not prove cloud deployment, runtime security, legal acceptance, production data recovery, or business acceptance.

Run `verify_pack.ps1` first. Then copy and complete `release-profile.template.json` and run:

```powershell
pwsh ./build_release.ps1 -ProjectRoot <project> -Profile <profile.json> -ReleaseRoot <new-directory> -OsvScanner <exact-osv.exe> -SigningKey <existing-key>
pwsh ./verify_release.ps1 -ReleaseRoot <new-directory>
```

The release directory is create-only. Altering the artifact, manifest, provenance, signature, signer policy, SBOM, SCA report, or receipt makes verification fail. `ssh-keygen` must be an Authenticode-valid Microsoft Windows component when `require_microsoft_authenticode` is enabled. OSV must match the locked 2.5.1 executable hash.

Authorities are pinned in the pack metadata: SLSA v1.2 provenance, in-toto attestation v1.2.0 and OpenSSH portable `V_10_2_P1`. The four files under `fixtures/` are Base64 transport wrappers that reconstruct the exact public OpenSSH regression bytes; no private fixture key is included.
````

### FILE: `portable_release_evidence_gate/release-profile.template.json`
```yaml
block_id: "PORTABLE-SIGNED-RELEASE-EVIDENCE-GATE:profile:v1"
operation: CREATE
provenance: AUTHORED
source: "local, governed by the pinned release evidence authorities"
license: "LicenseRef-Workspace-Owner"
sha256: "71a005aef877421e43db25b6a5d9e9000e784ed239022336a34c6f306074822b"
variables: []
secrets_allowed: false
```
````json
{
  "schema_version": 1,
  "enabled": false,
  "product_name": "",
  "product_version": "",
  "source_uri": "",
  "source_revision": "",
  "invocation_id": "",
  "builder_id": "",
  "build_script_ref": "",
  "build_script_sha256": "",
  "source_inputs_ref": "",
  "dependency_manifest_refs": [],
  "notices_ref": "",
  "release_identity": "",
  "allowed_signers_ref": "",
  "artifact_zip_timestamp": "1980-01-01T00:00:00",
  "require_microsoft_authenticode": true
}
````

### FILE: `portable_release_evidence_gate/build_release.ps1`
```yaml
block_id: "PORTABLE-SIGNED-RELEASE-EVIDENCE-GATE:builder:v1"
operation: CREATE
provenance: AUTHORED
source: "local glue implementing pinned SLSA/in-toto/OpenSSH/OSV contracts"
license: "LicenseRef-Workspace-Owner"
sha256: "86929e0e03ca0307a8e8aa3c9a9799c6175b5e35abfc37936899fd0a02c64019"
variables: []
secrets_allowed: true
```
````powershell
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
function Sha([string]$Path){(Get-FileHash -LiteralPath $Path -Algorithm SHA256).Hash.ToLowerInvariant()}
function WriteJson([string]$Path,$Value){[IO.File]::WriteAllText($Path,(ConvertTo-Json $Value -Depth 30 -Compress)+"`n",$Utf8)}
function ResolveProjectFile([string]$Root,[string]$Ref,[string]$Label){
  if([string]::IsNullOrWhiteSpace($Ref) -or [IO.Path]::IsPathRooted($Ref) -or $Ref -match '(^|[\\/])\.\.([\\/]|$)'){Fail "$Label must be a safe project-relative file"}
  $p=[IO.Path]::GetFullPath((Join-Path $Root $Ref))
  if(-not $p.StartsWith($Root+[IO.Path]::DirectorySeparatorChar,[StringComparison]::OrdinalIgnoreCase) -or -not(Test-Path -LiteralPath $p -PathType Leaf)){Fail "$Label missing or escaped project"}
  $p
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
  $files=@(Get-ChildItem -LiteralPath $InputRoot -File -Recurse | Sort-Object {$_.FullName.Substring($InputRoot.Length+1).Replace('\\','/')})
  if($files.Count -lt 1){Fail 'build produced no files'}
  $stream=[IO.File]::Open($ZipPath,[IO.FileMode]::CreateNew,[IO.FileAccess]::ReadWrite,[IO.FileShare]::None)
  try{$zip=[IO.Compression.ZipArchive]::new($stream,[IO.Compression.ZipArchiveMode]::Create,$true)
    try{foreach($f in $files){if(($f.Attributes -band [IO.FileAttributes]::ReparsePoint) -ne 0){Fail 'build output contains a reparse point'};$rel=$f.FullName.Substring($InputRoot.Length+1).Replace('\\','/');$e=$zip.CreateEntry($rel,[IO.Compression.CompressionLevel]::NoCompression);$e.LastWriteTime=$Stamp;$src=[IO.File]::OpenRead($f.FullName);try{$dst=$e.Open();try{$src.CopyTo($dst)}finally{$dst.Dispose()}}finally{$src.Dispose()}}}finally{$zip.Dispose()}
  }finally{$stream.Dispose()}
}

$root=(Resolve-Path -LiteralPath $ProjectRoot).Path
$profilePath=(Resolve-Path -LiteralPath $Profile).Path
if(Test-Path -LiteralPath $ReleaseRoot){Fail 'release root already exists'}
$cfg=Get-Content -LiteralPath $profilePath -Raw | ConvertFrom-Json -Depth 30 -DateKind String
if($cfg.schema_version -ne 1 -or $cfg.enabled -ne $true){Fail 'profile is not explicitly enabled'}
foreach($field in @('product_name','product_version','source_uri','source_revision','invocation_id','builder_id','build_script_ref','build_script_sha256','source_inputs_ref','notices_ref','release_identity','allowed_signers_ref')){if([string]::IsNullOrWhiteSpace([string]$cfg.$field)){Fail "$field is required"}}
if($cfg.product_name -notmatch '^[A-Za-z0-9][A-Za-z0-9._-]{0,127}$' -or $cfg.product_version -notmatch '^[A-Za-z0-9][A-Za-z0-9._+-]{0,127}$'){Fail 'unsafe product identity'}
if($cfg.release_identity -notmatch '^[A-Za-z0-9][A-Za-z0-9@._+-]{0,127}$'){Fail 'unsafe release identity'}
foreach($uriField in @('source_uri','builder_id')){try{if(-not([uri]$cfg.$uriField).IsAbsoluteUri){Fail "$uriField must be absolute"}}catch{Fail "$uriField must be an absolute URI"}}
try{$stamp=[datetimeoffset]::ParseExact([string]$cfg.artifact_zip_timestamp,'yyyy-MM-ddTHH:mm:ss',[Globalization.CultureInfo]::InvariantCulture,[Globalization.DateTimeStyles]::None)}catch{Fail 'invalid artifact ZIP timestamp'}
if($stamp.Year -lt 1980){Fail 'ZIP timestamp must be 1980 or later'}
$buildScript=ResolveProjectFile $root ([string]$cfg.build_script_ref) 'build_script_ref'
if($cfg.build_script_sha256 -notmatch '^[0-9a-f]{64}$' -or (Sha $buildScript) -cne [string]$cfg.build_script_sha256){Fail 'build script hash mismatch'}
$inputsFile=ResolveProjectFile $root ([string]$cfg.source_inputs_ref) 'source_inputs_ref'
$notices=ResolveProjectFile $root ([string]$cfg.notices_ref) 'notices_ref'
$allowed=ResolveProjectFile $root ([string]$cfg.allowed_signers_ref) 'allowed_signers_ref'
$inputRefs=@(Get-Content -LiteralPath $inputsFile | ForEach-Object {$_.Trim()} | Where-Object {$_ -and -not $_.StartsWith('#')})
if($inputRefs.Count -lt 1 -or @($inputRefs|Sort-Object -Unique).Count -ne $inputRefs.Count){Fail 'source input manifest is empty or duplicated'}
$sourceEntries=@();foreach($ref in ($inputRefs|Sort-Object)){ $p=ResolveProjectFile $root $ref 'source input';$sourceEntries += [ordered]@{path=$ref.Replace('\\','/');sha256=Sha $p;size=(Get-Item -LiteralPath $p).Length} }
$inputSet=@{};$inputRefs|ForEach-Object{$inputSet[$_.Replace('\\','/')]=1}
foreach($required in @([string]$cfg.build_script_ref,[string]$cfg.source_inputs_ref,[string]$cfg.notices_ref,[string]$cfg.allowed_signers_ref)+@($cfg.dependency_manifest_refs)){if(-not $inputSet.ContainsKey($required.Replace('\\','/'))){Fail "required source input not listed: $required"};[void](ResolveProjectFile $root $required 'required source input')}
if(@($cfg.dependency_manifest_refs).Count -lt 1){Fail 'at least one dependency manifest is required'}
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
try{
  $out1=Join-Path $work 'build1';$out2=Join-Path $work 'build2';New-Item -ItemType Directory -Path $out1,$out2|Out-Null
  $started=[datetimeoffset]::UtcNow
  foreach($out in @($out1,$out2)){InvokeChecked -Exe 'pwsh' -Arguments @('-NoProfile','-File',$buildScript,'-ProjectRoot',$root,'-OutputRoot',$out,'-SourceDateEpoch',[string]$stamp.ToUnixTimeSeconds()) -Label 'project build'}
  $zip1=Join-Path $work 'artifact1.zip';$zip2=Join-Path $work 'artifact2.zip';WriteDeterministicZip $out1 $zip1 $stamp;WriteDeterministicZip $out2 $zip2 $stamp
  $artifactSha=Sha $zip1;if((Sha $zip2)-cne $artifactSha){Fail 'two clean builds are not byte-identical'}
  $artifactEntries=@(Get-ChildItem -LiteralPath $out1 -File -Recurse|Sort-Object {$_.FullName.Substring($out1.Length+1).Replace('\\','/')}|ForEach-Object{[ordered]@{path=$_.FullName.Substring($out1.Length+1).Replace('\\','/');sha256=Sha $_.FullName;size=$_.Length}})
  $sourceManifest=Join-Path $work 'source.manifest.json';WriteJson $sourceManifest ([ordered]@{schema_version=1;files=$sourceEntries})
  $artifactManifest=Join-Path $work 'artifact.manifest.json';WriteJson $artifactManifest ([ordered]@{schema_version=1;artifact_sha256=$artifactSha;files=$artifactEntries})
  $sca=Join-Path $work 'osv.json';& $osv scan source --format json --output-file $sca $root;$scanExit=$LASTEXITCODE
  if(-not(Test-Path -LiteralPath $sca)){Fail 'OSV JSON report missing'};$scan=Get-Content -LiteralPath $sca -Raw|ConvertFrom-Json -Depth 100;$findings=@();foreach($result in @($scan.results)){foreach($pkg in @($result.packages)){foreach($v in @($pkg.vulnerabilities)){$findings += [string]$v.id}}};if($scanExit -ne 0 -or $findings.Count -ne 0){Fail "OSV found $($findings.Count) vulnerability records"}
  $sbom=Join-Path $work 'SBOM.spdx.json';InvokeChecked -Exe $osv -Arguments @('scan','source','--format','spdx-2-3','--output-file',$sbom,$root) -Label 'SPDX generation'
  $release=Join-Path $work 'release';New-Item -ItemType Directory -Path $release|Out-Null
  Copy-Item -LiteralPath $zip1 -Destination (Join-Path $release 'artifact.zip');Copy-Item -LiteralPath $artifactManifest -Destination $release;Copy-Item -LiteralPath $sourceManifest -Destination $release;Copy-Item -LiteralPath $sca -Destination $release;Copy-Item -LiteralPath $sbom -Destination $release;Copy-Item -LiteralPath $notices -Destination (Join-Path $release 'THIRD_PARTY_NOTICES.txt');Copy-Item -LiteralPath $allowed -Destination (Join-Path $release 'allowed_signers')
  $finished=[datetimeoffset]::UtcNow
  $statement=[ordered]@{_type='https://in-toto.io/Statement/v1';subject=@([ordered]@{name="$($cfg.product_name)-$($cfg.product_version).zip";digest=[ordered]@{sha256=$artifactSha}});predicateType='https://slsa.dev/provenance/v1';predicate=[ordered]@{buildDefinition=[ordered]@{buildType='https://slsa.dev/provenance/v1';externalParameters=[ordered]@{product_name=$cfg.product_name;product_version=$cfg.product_version;source_uri=$cfg.source_uri;source_revision=$cfg.source_revision;invocation_id=$cfg.invocation_id};internalParameters=[ordered]@{profile_sha256=Sha $profilePath;source_manifest_sha256=Sha $sourceManifest;artifact_manifest_sha256=Sha $artifactManifest;sbom_sha256=Sha $sbom;sca_sha256=Sha $sca;notices_sha256=Sha $notices;allowed_signers_sha256=Sha $allowed;artifact_zip_timestamp=$stamp.ToString('yyyy-MM-ddTHH:mm:ss')};resolvedDependencies=@($sourceEntries|ForEach-Object{[ordered]@{uri=([uri]::new([uri]$cfg.source_uri,$_.path)).AbsoluteUri;digest=[ordered]@{sha256=$_.sha256}}})};runDetails=[ordered]@{builder=[ordered]@{id=$cfg.builder_id};metadata=[ordered]@{invocationId=$cfg.invocation_id;startedOn=$started.ToString('O');finishedOn=$finished.ToString('O')}}}}
  $statementPath=Join-Path $release 'provenance.intoto.json';WriteJson $statementPath $statement
  InvokeChecked -Exe $ssh -Arguments @('-Y','sign','-f',$key,'-n',$Namespace,$statementPath) -Label 'provenance signing'
  $signature=$statementPath+'.sig';if(-not(Test-Path -LiteralPath $signature)){Fail 'signature file missing'}
  InvokeSshVerify $ssh (Join-Path $release 'allowed_signers') ([string]$cfg.release_identity) $signature $statementPath
  $authStatus=Get-AuthenticodeSignature -LiteralPath $ssh
  $receipt=[ordered]@{schema_version=1;gate='PORTABLE-SIGNED-RELEASE-EVIDENCE-GATE';status='PASS';product=[ordered]@{name=$cfg.product_name;version=$cfg.product_version};artifact_sha256=$artifactSha;artifact_zip_timestamp=$stamp.ToString('yyyy-MM-ddTHH:mm:ss');reproducibility_runs=2;source_revision=$cfg.source_revision;invocation_id=$cfg.invocation_id;release_identity=$cfg.release_identity;namespace=$Namespace;tools=[ordered]@{osv=[ordered]@{version=$osvVersion;sha256=Sha $osv};ssh_keygen=[ordered]@{version=(Get-Item -LiteralPath $ssh).VersionInfo.ProductVersion;sha256=Sha $ssh;authenticode_status=[string]$authStatus.Status;signer_subject=$authStatus.SignerCertificate.Subject}};evidence=[ordered]@{artifact_manifest_sha256=Sha (Join-Path $release 'artifact.manifest.json');source_manifest_sha256=Sha (Join-Path $release 'source.manifest.json');sbom_sha256=Sha (Join-Path $release 'SBOM.spdx.json');sca_sha256=Sha (Join-Path $release 'osv.json');notices_sha256=Sha (Join-Path $release 'THIRD_PARTY_NOTICES.txt');allowed_signers_sha256=Sha (Join-Path $release 'allowed_signers');provenance_sha256=Sha $statementPath;signature_sha256=Sha $signature}}
  WriteJson (Join-Path $release 'RELEASE_RECEIPT.json') $receipt
  Move-Item -LiteralPath $release -Destination $target
  Write-Output "PORTABLE_SIGNED_RELEASE_PASS artifact_sha256=$artifactSha findings=0 reproducibility_runs=2"
}finally{if(Test-Path -LiteralPath $work){$resolved=[IO.Path]::GetFullPath($work);if($resolved.StartsWith([IO.Path]::GetFullPath($parent)+[IO.Path]::DirectorySeparatorChar+'.elite-release-staging-',[StringComparison]::OrdinalIgnoreCase)){Remove-Item -LiteralPath $resolved -Recurse -Force}}}
````

### FILE: `portable_release_evidence_gate/verify_release.ps1`
```yaml
block_id: "PORTABLE-SIGNED-RELEASE-EVIDENCE-GATE:release-verifier:v1"
operation: CREATE
provenance: AUTHORED
source: "local glue implementing pinned SLSA/in-toto/OpenSSH/OSV contracts"
license: "LicenseRef-Workspace-Owner"
sha256: "0caa93dbe4ded971d723df3df780a62cd871a22a2a0c04c9a14f5b443a6667c5"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0
param([Parameter(Mandatory=$true)][string]$ReleaseRoot,[string]$SshKeygen='ssh-keygen')
$ErrorActionPreference='Stop'
function Fail([string]$Message){throw "SIGNED_RELEASE_VERIFY: $Message"}
function Sha([string]$Path){(Get-FileHash -LiteralPath $Path -Algorithm SHA256).Hash.ToLowerInvariant()}
$root=(Resolve-Path -LiteralPath $ReleaseRoot).Path
$required=@('artifact.zip','artifact.manifest.json','source.manifest.json','osv.json','SBOM.spdx.json','THIRD_PARTY_NOTICES.txt','allowed_signers','provenance.intoto.json','provenance.intoto.json.sig','RELEASE_RECEIPT.json')
foreach($name in $required){if(-not(Test-Path -LiteralPath (Join-Path $root $name) -PathType Leaf)){Fail "missing $name"}}
$receipt=Get-Content -LiteralPath (Join-Path $root 'RELEASE_RECEIPT.json') -Raw|ConvertFrom-Json -Depth 30 -DateKind String
if($receipt.schema_version -ne 1 -or $receipt.gate -ne 'PORTABLE-SIGNED-RELEASE-EVIDENCE-GATE' -or $receipt.status -ne 'PASS' -or $receipt.reproducibility_runs -ne 2){Fail 'invalid receipt contract'}
if((Sha (Join-Path $root 'artifact.zip')) -cne [string]$receipt.artifact_sha256){Fail 'artifact hash mismatch'}
$map=[ordered]@{artifact_manifest_sha256='artifact.manifest.json';source_manifest_sha256='source.manifest.json';sbom_sha256='SBOM.spdx.json';sca_sha256='osv.json';notices_sha256='THIRD_PARTY_NOTICES.txt';allowed_signers_sha256='allowed_signers';provenance_sha256='provenance.intoto.json';signature_sha256='provenance.intoto.json.sig'}
foreach($e in $map.GetEnumerator()){if((Sha (Join-Path $root $e.Value))-cne[string]$receipt.evidence.($e.Key)){Fail "$($e.Value) hash mismatch"}}
$statement=Get-Content -LiteralPath (Join-Path $root 'provenance.intoto.json') -Raw|ConvertFrom-Json -Depth 30 -DateKind String
if($statement._type -ne 'https://in-toto.io/Statement/v1' -or $statement.predicateType -ne 'https://slsa.dev/provenance/v1' -or @($statement.subject).Count -ne 1 -or $statement.subject[0].digest.sha256 -cne [string]$receipt.artifact_sha256){Fail 'provenance subject mismatch'}
$parameters=$statement.predicate.buildDefinition;if($parameters.externalParameters.product_name-cne[string]$receipt.product.name-or$parameters.externalParameters.product_version-cne[string]$receipt.product.version-or$parameters.externalParameters.source_revision-cne[string]$receipt.source_revision-or$parameters.externalParameters.invocation_id-cne[string]$receipt.invocation_id-or$statement.predicate.runDetails.metadata.invocationId-cne[string]$receipt.invocation_id){Fail 'provenance/receipt identity mismatch'}
$internal=$parameters.internalParameters;foreach($pair in @(@('artifact_manifest_sha256','artifact_manifest_sha256'),@('source_manifest_sha256','source_manifest_sha256'),@('sbom_sha256','sbom_sha256'),@('sca_sha256','sca_sha256'),@('notices_sha256','notices_sha256'),@('allowed_signers_sha256','allowed_signers_sha256'))){if([string]$internal.($pair[0])-cne[string]$receipt.evidence.($pair[1])){Fail "signed internal evidence mismatch: $($pair[0])"}};if([string]$internal.artifact_zip_timestamp-cne[string]$receipt.artifact_zip_timestamp){Fail 'signed artifact timestamp mismatch'}
$manifest=Get-Content -LiteralPath (Join-Path $root 'artifact.manifest.json') -Raw|ConvertFrom-Json -Depth 30
if($manifest.artifact_sha256 -cne [string]$receipt.artifact_sha256){Fail 'artifact manifest subject mismatch'}
$sca=Get-Content -LiteralPath (Join-Path $root 'osv.json') -Raw|ConvertFrom-Json -Depth 100;$findings=@();foreach($result in @($sca.results)){foreach($pkg in @($result.packages)){foreach($v in @($pkg.vulnerabilities)){$findings += [string]$v.id}}};if($findings.Count-ne0){Fail 'stored OSV report contains findings'}
$sbom=Get-Content -LiteralPath (Join-Path $root 'SBOM.spdx.json') -Raw|ConvertFrom-Json -Depth 100;if($sbom.spdxVersion-ne'SPDX-2.3'){Fail 'stored SBOM is not SPDX 2.3'}
try{$stamp=[datetime]::ParseExact([string]$receipt.artifact_zip_timestamp,'yyyy-MM-ddTHH:mm:ss',[Globalization.CultureInfo]::InvariantCulture,[Globalization.DateTimeStyles]::None)}catch{Fail 'receipt artifact ZIP timestamp invalid'}
Add-Type -AssemblyName System.IO.Compression;$s=[IO.File]::OpenRead((Join-Path $root 'artifact.zip'));try{$z=[IO.Compression.ZipArchive]::new($s,[IO.Compression.ZipArchiveMode]::Read,$true);try{$entries=@($z.Entries|Where-Object{$_.Name -ne ''});if($entries.Count-ne@($manifest.files).Count){Fail 'ZIP entry count mismatch'};$seen=@{};$names=@();foreach($e in $entries){$names+=$e.FullName;if($e.FullName.Contains('\\')-or $e.FullName.StartsWith('/')-or $e.FullName.Split('/')-contains '..'-or$seen.ContainsKey($e.FullName)){Fail 'unsafe or duplicate ZIP entry'};if($e.LastWriteTime.ToString('yyyy-MM-ddTHH:mm:ss')-cne$stamp.ToString('yyyy-MM-ddTHH:mm:ss')){Fail 'ZIP timestamp mismatch'};if($e.CompressedLength-ne$e.Length){Fail 'ZIP entry is not stored deterministically'};$seen[$e.FullName]=1;$expected=@($manifest.files|Where-Object path -ceq $e.FullName);if($expected.Count-ne1){Fail 'ZIP entry absent from manifest'};$stream=$e.Open();try{$hash=[Security.Cryptography.SHA256]::HashData($stream);$actual=[Convert]::ToHexString($hash).ToLowerInvariant()}finally{$stream.Dispose()};if($actual-cne[string]$expected[0].sha256-or$e.Length-ne[long]$expected[0].size){Fail 'ZIP entry digest or size mismatch'}};$sorted=@($names|Sort-Object);if(($names-join"`n")-cne($sorted-join"`n")){Fail 'ZIP entries are not sorted'}}finally{$z.Dispose()}}finally{$s.Dispose()}
$ssh=(Get-Command $SshKeygen -ErrorAction Stop).Source;if((Sha $ssh)-cne[string]$receipt.tools.ssh_keygen.sha256){Fail 'ssh-keygen identity changed'};$auth=Get-AuthenticodeSignature -LiteralPath $ssh;if($auth.Status-ne'Valid'-or$auth.SignerCertificate.Subject-notmatch'O=Microsoft Corporation'){Fail 'ssh-keygen Authenticode invalid'}
$psi=[Diagnostics.ProcessStartInfo]::new();$psi.FileName=$ssh;foreach($a in @('-Y','verify','-f',(Join-Path $root 'allowed_signers'),'-I',[string]$receipt.release_identity,'-n',[string]$receipt.namespace,'-s',(Join-Path $root 'provenance.intoto.json.sig'))){[void]$psi.ArgumentList.Add($a)};$psi.UseShellExecute=$false;$psi.RedirectStandardInput=$true;$psi.RedirectStandardOutput=$true;$psi.RedirectStandardError=$true;$p=[Diagnostics.Process]::Start($psi);$bytes=[IO.File]::ReadAllBytes((Join-Path $root 'provenance.intoto.json'));$p.StandardInput.BaseStream.Write($bytes,0,$bytes.Length);$p.StandardInput.Close();$p.WaitForExit();if($p.ExitCode-ne0){Fail 'OpenSSH signature rejected'}
Write-Output "PORTABLE_SIGNED_RELEASE_VERIFY_PASS artifact_sha256=$($receipt.artifact_sha256)"
````

### FILE: `portable_release_evidence_gate/verify_pack.ps1`
```yaml
block_id: "PORTABLE-SIGNED-RELEASE-EVIDENCE-GATE:pack-verifier:v1"
operation: CREATE
provenance: AUTHORED
source: "local verifier plus exact OpenSSH public regression vector"
license: "LicenseRef-Workspace-Owner"
sha256: "9cb879f54092fc6d88673fc88d4dbe2d69d2d2201c428c7709729f03cde1d225"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0
param([string]$SshKeygen='ssh-keygen')
$ErrorActionPreference='Stop'
function Fail([string]$Message){throw "SIGNED_RELEASE_PACK: $Message"}
function Sha([string]$Path){(Get-FileHash -LiteralPath $Path -Algorithm SHA256).Hash.ToLowerInvariant()}
$expected=[ordered]@{'openssh-ed25519.pub'='1b53846ff8047b6e038c315c4351c35547f5e25756c340bdabe863b9882617e9';'openssh-ed25519.sig'='7e2222d980c5fdcea433cc2b9a9d00dcc28ed5b4a8b9d1fde7822802f6872e9c';'openssh-namespace'='df08227309c92f2f3f030b423009e697445371283d89352a26d064d39fb73c36';'openssh-signed-data'='3d69dd3e7115cb2f5e44836852a8a8dcf0c3642a15faee1b34bbf1bf9f9980f6'}
$ssh=(Get-Command $SshKeygen -ErrorAction Stop).Source;$auth=Get-AuthenticodeSignature -LiteralPath $ssh;if($auth.Status-ne'Valid'-or$auth.SignerCertificate.Subject-notmatch'O=Microsoft Corporation'){Fail 'ssh-keygen is not Microsoft Authenticode-valid'}
$work=Join-Path $env:TEMP ('elite-sshsig-public-'+[guid]::NewGuid().ToString('N'));New-Item -ItemType Directory -Path $work|Out-Null
try{foreach($e in $expected.GetEnumerator()){$bytes=[Convert]::FromBase64String((Get-Content -LiteralPath (Join-Path $PSScriptRoot ('fixtures/'+$e.Key+'.b64')) -Raw).Trim());$path=Join-Path $work $e.Key;[IO.File]::WriteAllBytes($path,$bytes);if((Sha $path)-cne$e.Value){Fail "fixture hash mismatch: $($e.Key)"}};$pub=(Get-Content -LiteralPath (Join-Path $work 'openssh-ed25519.pub') -Raw).Trim();$allowed=Join-Path $work 'allowed_signers';[IO.File]::WriteAllText($allowed,"openssh-test $pub`n",[Text.UTF8Encoding]::new($false));function TestSignature([byte[]]$Bytes){$psi=[Diagnostics.ProcessStartInfo]::new();$psi.FileName=$ssh;foreach($a in @('-Y','verify','-f',$allowed,'-I','openssh-test','-n',(Get-Content -LiteralPath (Join-Path $work 'openssh-namespace') -Raw).Trim(),'-s',(Join-Path $work 'openssh-ed25519.sig'))){[void]$psi.ArgumentList.Add($a)};$psi.UseShellExecute=$false;$psi.RedirectStandardInput=$true;$psi.RedirectStandardOutput=$true;$psi.RedirectStandardError=$true;$p=[Diagnostics.Process]::Start($psi);$p.StandardInput.BaseStream.Write($Bytes,0,$Bytes.Length);$p.StandardInput.Close();$p.WaitForExit();$p.ExitCode};$bytes=[IO.File]::ReadAllBytes((Join-Path $work 'openssh-signed-data'));if((TestSignature $bytes)-ne0){Fail 'official positive signature rejected'};$bad=[byte[]]$bytes.Clone();$bad[0]=[byte]($bad[0]-bxor1);if((TestSignature $bad)-eq0){Fail 'tampered payload accepted'};foreach($script in @('build_release.ps1','verify_release.ps1','verify_pack.ps1')){$errors=$null;[void][Management.Automation.Language.Parser]::ParseFile((Join-Path $PSScriptRoot $script),[ref]$null,[ref]$errors);if($errors.Count){Fail "$script parse failed"}};Write-Output "PORTABLE_SIGNED_RELEASE_PACK_PASS openssh_vector=1 tamper_rejected=1 authenticode=Valid"}finally{if(Test-Path -LiteralPath $work){$resolved=[IO.Path]::GetFullPath($work);$temp=[IO.Path]::GetFullPath([IO.Path]::GetTempPath());if($resolved.StartsWith($temp+'elite-sshsig-public-',[StringComparison]::OrdinalIgnoreCase)){Remove-Item -LiteralPath $resolved -Recurse -Force}}}
````

### FILE: `portable_release_evidence_gate/fixtures/openssh-ed25519.pub.b64`
```yaml
block_id: "PORTABLE-SIGNED-RELEASE-EVIDENCE-GATE:openssh-pub:v1"
operation: CREATE
provenance: ADAPTED
source: "OpenSSH blob b078e4516fbe44d2cf03e32a98e624bcb2ea6875, Base64 transport without byte changes after decode"
license: "LicenseRef-OpenSSH"
sha256: "65606d6bf81a6ebb0f5e6af2768cfd5b1ab350d22e4546ac28a1cc819b2b1cd8"
variables: []
secrets_allowed: false
```
````text
c3NoLWVkMjU1MTkgQUFBQUMzTnphQzFsWkRJMU5URTVBQUFBSUlsaXpTSU4zRFFWNzhWUE5qVnZ2d2pnZitQNUhxYlBZQ1l1M0JPTWRqQUUgRUQyNTUxOSB0ZXN0IGtleQo=
````

### FILE: `portable_release_evidence_gate/fixtures/openssh-ed25519.sig.b64`
```yaml
block_id: "PORTABLE-SIGNED-RELEASE-EVIDENCE-GATE:openssh-signature:v1"
operation: CREATE
provenance: ADAPTED
source: "OpenSSH blob 8e8ff2a8ac197df506ee00eef1d398a8090609c5, Base64 transport without byte changes after decode"
license: "LicenseRef-OpenSSH"
sha256: "238dc6513b90392612b7c821772720d90e34f384d529fb15ff8c2be6a11b27a0"
variables: []
secrets_allowed: false
```
````text
LS0tLS1CRUdJTiBTU0ggU0lHTkFUVVJFLS0tLS0KVTFOSVUwbEhBQUFBQVFBQUFETUFBQUFMYzNOb0xXVmtNalUxTVRrQUFBQWdpV0xOSWczY05CWHZ4VTgyTlcrL0NPQi80LwprZXBzOWdKaTdjRTR4Mk1BUUFBQUFJZFc1cGRIUmxjM1FBQUFBQUFBQUFCbk5vWVRVeE1nQUFBRk1BQUFBTGMzTm9MV1ZrCk1qVTFNVGtBQUFCQWloUXNiVXp1TkVGZmxrNVR3MStIOWFMUzd0WlFrMFJHOEtXMUR0T21EWVluV2UzRDNVS2lHM2ZjSmEKRE5nNHZCV3AxajFnTFJpQk1PRitnd1lOZWdEZz09Ci0tLS0tRU5EIFNTSCBTSUdOQVRVUkUtLS0tLQo=
````

### FILE: `portable_release_evidence_gate/fixtures/openssh-namespace.b64`
```yaml
block_id: "PORTABLE-SIGNED-RELEASE-EVIDENCE-GATE:openssh-namespace:v1"
operation: CREATE
provenance: ADAPTED
source: "OpenSSH blob 1570cd548baa9998c08cd500e88daa254dbbe66c, Base64 transport without byte changes after decode"
license: "LicenseRef-OpenSSH"
sha256: "64df7eca384bf6c35d9ba9d6681a1ba303a2d8434592557949db17b070974065"
variables: []
secrets_allowed: false
```
````text
dW5pdHRlc3Q=
````

### FILE: `portable_release_evidence_gate/fixtures/openssh-signed-data.b64`
```yaml
block_id: "PORTABLE-SIGNED-RELEASE-EVIDENCE-GATE:openssh-data:v1"
operation: CREATE
provenance: ADAPTED
source: "OpenSSH blob 7df4bedd135c385ed6e5cf99c901e3ad6e44a043, Base64 transport without byte changes after decode"
license: "LicenseRef-OpenSSH"
sha256: "c0c81e7e27b86b6456396b2d1c6877c19e7a69e164053041e30b7c175a662a0e"
variables: []
secrets_allowed: false
```
````text
VGhpcyBpcyBhIHRlc3QsIHRoaXMgaXMgb25seSBhIHRlc3Q=
````

### FILE: `portable_release_evidence_gate/OPENSSH_FIXTURE_NOTICE.md`
```yaml
block_id: "PORTABLE-SIGNED-RELEASE-EVIDENCE-GATE:notice:v1"
operation: CREATE
provenance: AUTHORED
source: "local notice preserving exact OpenSSH revision, blobs and upstream license identity"
license: "LicenseRef-Workspace-Owner"
sha256: "b6a1420564a25b9f0a652404138a8fcf594c334bb9b029d1b4c28637bcc1926c"
variables: []
secrets_allowed: false
```
````markdown
# OpenSSH public regression fixture notice

The four Base64 files under `fixtures/` losslessly reconstruct bytes copied verbatim from the OpenSSH portable `V_10_2_P1` commit `d01efaa1c9ed84fd9011201dbc3c7cb0a82bcee3`, tree `c91beec59da56218afced35b14087d362b5bba8e`:

- `regress/unittests/sshsig/testdata/ed25519.pub`, blob `b078e4516fbe44d2cf03e32a98e624bcb2ea6875`;
- `regress/unittests/sshsig/testdata/ed25519.sig`, blob `8e8ff2a8ac197df506ee00eef1d398a8090609c5`;
- `regress/unittests/sshsig/testdata/namespace`, blob `1570cd548baa9998c08cd500e88daa254dbbe66c`;
- `regress/unittests/sshsig/testdata/signed-data`, blob `7df4bedd135c385ed6e5cf99c901e3ad6e44a043`.

OpenSSH portable publishes its full license and copyright notices at `https://github.com/openssh/openssh-portable/blob/d01efaa1c9ed84fd9011201dbc3c7cb0a82bcee3/LICENCE` (blob `aeb3017e76c6ca6a53337dfb4af2f434e1628568`, 21,140 bytes, SHA-256 `5bb5b160726ef5756e4f32fe95b35249c294962419650f48d05134b486d27ccb`). Preserve this notice and the linked upstream license when redistributing the fixtures. No OpenSSH executable, source implementation, or private key is embedded by this pack.
````

## 6. Configuration surface

| Campo | Default seguro | Gate |
|---|---|---|
| `enabled` | `false` | debe habilitarse explícitamente |
| `build_script_ref` + SHA | vacío | archivo relativo incluido en inputs y hash exacto |
| `source_inputs_ref` | vacío | lista explícita, segura, única y no vacía |
| `dependency_manifest_refs` | vacío | al menos uno, existente e incluido |
| `allowed_signers_ref` | vacío | una identidad exacta y una clave pública Ed25519 |
| `SigningKey` | sin default | existente, protegido y fuera del proyecto; nunca se copia |
| `ReleaseRoot` | sin default | nuevo, fuera del proyecto y con padre existente |
| `artifact_zip_timestamp` | `1980-01-01T00:00:00` | timestamp DOS/ZIP explícito, sin afirmar timezone inexistente |
| `require_microsoft_authenticode` | `true` | `ssh-keygen` debe tener firma Microsoft válida |

## 7. Dependency bill

| Componente | Revisión | Uso | Licencia/servicing | Admisión |
|---|---|---|---|---|
| Google OSV-Scanner | 2.5.1 / `c84fa456…`; exe SHA `25e42f5e…d9bfb6` | SCA JSON + SPDX 2.3 | Apache-2.0 | exacto y hash-locked |
| Windows OpenSSH `ssh-keygen` | versión host registrada; Authenticode Microsoft válido | `sshsig` Ed25519 | componente Windows/OpenSSH | condicionado a servicing y firma válida |
| SLSA | spec v1.2 / `19e4e2f…` | predicate provenance v1 | Community Specification License | autoridad de schema |
| in-toto attestation | v1.2.0 / `df02077…` | Statement v1 | CC-BY-4.0 | autoridad de envelope/statement |
| OpenSSH portable | `V_10_2_P1` / `d01efaa…` | vector público de regresión | licencia OpenSSH | cuatro fixtures transportados en Base64 |

Microsoft SBOM Tool 4.1.5, su `main` 2026-04-24, Cosign 3.1.3 e in-toto-golang 0.11.0 fueron probados pero no incorporados: sus grafos actuales presentaron avisos OSV. Esta exclusión es deliberada y no se oculta con overrides locales.

## 8. Apply order

Materializar en un directorio dedicado; adquirir OSV desde el lock existente; verificar pack; crear fuera del proyecto una identidad Ed25519 y su política pública mediante el procedimiento de seguridad del usuario; completar el perfil; ejecutar build; ejecutar verificador independiente; conservar el directorio por digest. Nunca reutilizar la clave pública de test ni guardar una clave privada en el repositorio. Rollback selecciona un release anterior ya verificado por su digest; no reconstruye “lo mismo”.

## 9. Verification

`verify_pack.ps1` valida hashes del vector público oficial, Authenticode Microsoft, firma positiva, rechazo del payload alterado y parseo de los tres scripts. El test end-to-end debe ejecutar `build_release.ps1` con una clave de prueba protegida, obtener dos ZIP idénticos, cero hallazgos, SPDX, provenance y firma, pasar `verify_release.ps1` y rechazar una copia alterada. En cada proyecto repetir con su build, dependencias, clave, notices y inputs reales.

## 10. Reconstruction evidence

Windows 11 / PowerShell 7.6.5 / Microsoft OpenSSH `9.5p2` file version `9.5.6.1` con Authenticode válido / Google OSV-Scanner 2.5.1 exacto. El vector OpenSSH `V_10_2_P1` verificó y el payload alterado fue rechazado. Un proyecto Go sin dependencias produjo dos veces el ZIP SHA-256 `c6a018ae452644877c6f4c37b1acb36f753a05c2d110cef3ff89dac0da199fe3`, OSV 0, SPDX 2.3, provenance firmada y verificación independiente PASS; una copia con un byte alterado fue rechazada. La clave usada fue exclusivamente el fixture privado publicado por OpenSSH y permaneció fuera del proyecto y del pack.
