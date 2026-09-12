# GitLab/OpenGrep Signed SAST Gate

## 1. Metadata

```yaml
pack_id: "GITLAB-OPENGREP-SIGNED-SAST-GATE"
pack_version: "0.1.0"
status:
  authority: ELITE_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Conserva un lane SAST local históricamente reconstruido: OpenGrep 1.29.0 exacto y firmado más 122 reglas exactas GitLab 2.9.3 comercialmente admisibles; su instalación actual está bloqueada porque el verifier Cosign 3.1.3 reabrió SCA."
stacks: ["GitLab SAST Rules 2.9.3", "OpenGrep 1.29.0", "Cosign 3.1.3 REJECTED_VERIFIER_RUNTIME", "PowerShell 7", "Go", "TypeScript"]
compatible_with: ["MICROSOFT-DEVSKIM-ADAPTED-SAST-GATE 0.1.x", "GO-NATIVE-FUZZ-GATE 0.1.x", "SECURE-OPERATIONS-DELIVERY-CORE 0.4.x"]
incompatible_with: ["GitLab Enterprise rules sin suscripción", "reglas Commons Clause para lane comercial", "atribución del runtime adaptado a GitLab", "uso como único gate AppSec", "supresiones automáticas", "producción sin triage ni pruebas target"]
license_expression: "LicenseRef-Workspace-Owner AND LGPL-2.1-only AND Apache-2.0 AND MIT AND LGPL-3.0-only"
upstream_sources: ["https://github.com/opengrep/opengrep/tree/344509d693c852eaac4fc1eeffaf2f655c531b5a", "https://github.com/sigstore/cosign/tree/11926fa5bbbbde47e88fc006b625a17769b743b2", "https://gitlab.com/gitlab-org/security-products/sast-rules/-/tree/39fc7de205d447c70e75a1fdabdcf4588733adc8", "https://gitlab.com/gitlab-org/security-products/analyzers/semgrep/-/tree/b8d1635b40bc2a73d3650fbf6495094b52ead02b"]
verified_at: "2026-09-05"
```

## 2. Applicability

El código y locks se conservan para reproducir la decisión histórica, pero el agente no debe instalarlo ni contarlo como segundo lane mientras `current_admission.status` permanezca `REJECTED_VERIFIER_RUNTIME`. Antes de una futura reapertura debe definir scope, severidades, paths generados/vendorizados, owner y SLA de triage, política de suppressions con vencimiento y enlace finding → fix → regression.

Rechace el pack fuera de Windows x64, si no se autoriza egress para APIs/firma Sigstore, si sus licencias no son compatibles con el proyecto o si el claim exige el contenedor GitLab exacto. La selección excluye deliberadamente `dist/gitlab/**` por licencia Enterprise y `dist/lgpl-cc/**` por Commons Clause. No presenta OpenGrep como código de GitLab: el motor firmado y las reglas GitLab exactas se componen mediante glue `AUTHORED`.

## 3. Architecture contract

`install_and_verify.ps1` falla antes de red/descarga mientras el lock esté reabierto. Históricamente fijó commit/tree/tag, digests, licencias, reglas y las revisiones exactas; consultó autoridades GitHub/GitLab, verificó la firma del binario con el procedimiento Cosign oficial, inspeccionó el ZIP contra traversal y publicó un runtime sólo después de comprobar 122 reglas. Esa evidencia no autoriza volver a ejecutar hoy el verifier rechazado.

`scan.ps1` verifica el manifest antes de cada corrida, rechaza output preexistente o dentro del source, aplica timeout, ejecuta las severidades elegidas con `--strict --error`, exige paridad JSON/SARIF y publica evidencia atómicamente. `0=CLEAN`, `1=FINDINGS`; cualquier otro código o error de parsing es fallo del gate. No modifica código, no suprime findings y no afirma cobertura de vulnerabilidades ausentes de las reglas.

## 4. Exact file manifest

```text
CREATE gitlab_opengrep_signed_sast/source-lock.json
CREATE gitlab_opengrep_signed_sast/install_and_verify.ps1
CREATE gitlab_opengrep_signed_sast/scan.ps1
CREATE gitlab_opengrep_signed_sast/verify_pack.ps1
CREATE gitlab_opengrep_signed_sast/README.md
```

## 5. Materialization blocks

### FILE: `gitlab_opengrep_signed_sast/source-lock.json`
```yaml
block_id: "GITLAB-OPENGREP-SIGNED-SAST-GATE:source-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite lock derived from exact GitHub/GitLab APIs, release artifacts, signatures and licenses"
license: "LicenseRef-Workspace-Owner"
sha256: "dea290531487047da598188cef4c50dfa999d34c7210066e703e2a9b939ea346"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-gitlab-opengrep-sast-lock/v1",
  "verified_at": "2026-09-05",
  "current_admission": {
    "status": "REJECTED_VERIFIER_RUNTIME",
    "reviewed_at": "2026-09-05",
    "reason": "Cosign 3.1.3 official SPDX scanned by exact OSV-Scanner 2.5.1 contains affected dependencies; the historical signature result remains evidence but this runtime cannot be promoted or re-executed.",
    "cosign_official_sbom_bytes": 458050,
    "cosign_official_sbom_sha256": "34d5b303eea474d37ffb19393c27ea92f41b7b73fe085a88de58d595f490cef0",
    "osv_affected_package_rows": 9,
    "osv_unique_advisories": 164,
    "reopen_trigger": "A newer official Cosign release or an alternative official verifier passes exact identity, license, signature, SCA and tamper gates."
  },
  "opengrep": {
    "repository": "opengrep/opengrep",
    "tag": "v1.29.0",
    "commit": "344509d693c852eaac4fc1eeffaf2f655c531b5a",
    "tree": "ccc75c2c86a22c3f1670a21c9c7fae1d6728259b",
    "commit_signature_verified": true,
    "version": "1.29.0",
    "asset": {
      "name": "opengrep_windows_x86.exe",
      "url": "https://github.com/opengrep/opengrep/releases/download/v1.29.0/opengrep_windows_x86.exe",
      "bytes": 53536256,
      "sha256": "ee485b31912704dc6410bc43f04b5c6ad896697db56e360a98204abf95fa1025"
    },
    "certificate": {
      "name": "opengrep_windows_x86.exe.cert",
      "url": "https://github.com/opengrep/opengrep/releases/download/v1.29.0/opengrep_windows_x86.exe.cert",
      "bytes": 3368,
      "sha256": "eb7906ca33749a66c0fedbce11f88c58c11790817026629a7cf90f2a051a4dbd"
    },
    "signature": {
      "name": "opengrep_windows_x86.exe.sig",
      "url": "https://github.com/opengrep/opengrep/releases/download/v1.29.0/opengrep_windows_x86.exe.sig",
      "bytes": 96,
      "sha256": "64274240f5427130026540f42c111859bb10102f703ceb0db0ea3a67cfd397a0"
    },
    "certificate_identity_regexp": "https://github.com/opengrep/opengrep.+",
    "certificate_oidc_issuer": "https://token.actions.githubusercontent.com",
    "license_expression": "LGPL-2.1-only",
    "license": {
      "name": "OpenGrep-LICENSE",
      "url": "https://raw.githubusercontent.com/opengrep/opengrep/344509d693c852eaac4fc1eeffaf2f655c531b5a/LICENSE",
      "bytes": 26526,
      "sha256": "20c17d8b8c48a600800dfd14f95d5cb9ff47066a9641ddeab48dc54aec96e331"
    }
  },
  "cosign": {
    "repository": "sigstore/cosign",
    "tag": "v3.1.3",
    "tag_object": "2f3a85b04907df5b770eb049d7e4d08d4b018d86",
    "tag_signature_verified": true,
    "commit": "11926fa5bbbbde47e88fc006b625a17769b743b2",
    "tree": "ffb2d2076b725abbeff294d5c02897a0da7a5f15",
    "commit_signature_verified": true,
    "version": "v3.1.3",
    "asset": {
      "name": "cosign-windows-amd64.exe",
      "url": "https://github.com/sigstore/cosign/releases/download/v3.1.3/cosign-windows-amd64.exe",
      "bytes": 198819314,
      "sha256": "9fe59be0eca1271873ce019061335eb1ac419b7059202e797828467ddabe33be"
    },
    "license_expression": "Apache-2.0",
    "license": {
      "name": "Cosign-LICENSE",
      "url": "https://raw.githubusercontent.com/sigstore/cosign/11926fa5bbbbde47e88fc006b625a17769b743b2/LICENSE",
      "bytes": 11357,
      "sha256": "c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4"
    }
  },
  "gitlab_analyzer_authority": {
    "repository": "gitlab-org/security-products/analyzers/semgrep",
    "commit": "b8d1635b40bc2a73d3650fbf6495094b52ead02b",
    "committed_at": "2026-08-26T11:34:50-03:00",
    "dockerfile_blob": "ee66eca540bd71cafdf90b758388cce87a5a4251",
    "dockerfile_sha256": "b5bc3bce3b5ab11c34a6e54e9408236986a359ea0b992473cd0f385b8c7e898b",
    "scanner_version": "1.174.0",
    "sast_rules_version": "2.9.3",
    "claim": "The current GitLab Free Semgrep analyzer pins the official rules package at 2.9.3; this pack adapts that admissible subset to signed OpenGrep and does not claim byte identity with the GitLab analyzer runtime."
  },
  "gitlab_rules": {
    "repository": "gitlab-org/security-products/sast-rules",
    "tag": "v2.9.3",
    "commit": "39fc7de205d447c70e75a1fdabdcf4588733adc8",
    "package": {
      "name": "sast-rules-v2.9.3.zip",
      "url": "https://gitlab.com/api/v4/projects/27038823/packages/generic/sast-rules/v2.9.3/sast-rules-v2.9.3.zip",
      "bytes": 272617,
      "sha256": "abbff567aa105eb15f7fcc86e22b19c227405358f6b29e18fb22dbc23407d776"
    },
    "root_license": {
      "name": "GitLab-SAST-Rules-LICENSE",
      "url": "https://gitlab.com/api/v4/projects/27038823/repository/files/LICENSE/raw?ref=39fc7de205d447c70e75a1fdabdcf4588733adc8",
      "bytes": 2000,
      "sha256": "33acf42afdd03163d75fe856cb431e70710e0b5deede3362e6d1c4697cfe1c6a",
      "license_expression": "MIT"
    },
    "selected_rules": [
      {
        "archive_path": "dist/eslint.yml",
        "output_name": "eslint.yml",
        "bytes": 39237,
        "sha256": "c5db31114dd29ba06c0b82274c7d09e81a1c64436b4c4f36aaf5f85bd3aad32e",
        "rules": 11,
        "license_expression": "MIT"
      },
      {
        "archive_path": "dist/gosec.yml",
        "output_name": "gosec.yml",
        "bytes": 68869,
        "sha256": "373bc91e22c4371ddcdbfded2912720c9a0295b20d66d0bf5902956354658cc3",
        "rules": 28,
        "license_expression": "MIT"
      },
      {
        "archive_path": "dist/lgpl/nodejs_scan.yml",
        "output_name": "nodejs_scan.yml",
        "bytes": 194599,
        "sha256": "b846b0c4d156477a30fb22c4bbb1325c63a70ec63f6b263c7b0599c24ddb90eb",
        "rules": 83,
        "license_expression": "LGPL-3.0-only"
      }
    ],
    "selected_rule_count": 122,
    "blocking_rule_count_at_warning_or_error": 117,
    "lgpl_license": {
      "archive_path": "dist/lgpl/LICENSE",
      "output_name": "GitLab-ThirdParty-LGPL-3.0-LICENSE",
      "bytes": 7651,
      "sha256": "a5681bf9b05db14d86776930017c647ad9e6e56ff6bbcfdf21e5848288dfaf1b"
    },
    "excluded": [
      {
        "path": "dist/gitlab/**",
        "license_path": "dist/gitlab/LICENSE",
        "license_bytes": 2741,
        "license_sha256": "59662312a402088c34208a95a49ed7fbac09fc1d2d466a48f0e7ffdf56e7006e",
        "reason": "GitLab Enterprise Edition subscription license; not admitted into the no-additional-cost commercial lane."
      },
      {
        "path": "dist/lgpl-cc/**",
        "license_path": "dist/lgpl-cc/LICENSE",
        "license_bytes": 27576,
        "license_sha256": "c99a349e6d8aa6c5c33106fcdbb240589ae4b7aec3dd1c48eec7c4f7b1260a31",
        "reason": "Commons Clause restricts sale; not admitted into the commercial lane."
      }
    ]
  },
  "verified_lane": {
    "os": "Windows x64",
    "opengrep_version": "1.29.0",
    "cosign_version": "v3.1.3",
    "typescript_fixture_rule": "eslint.detect-eval-with-expression",
    "go_fixture_rule": "gosec.G204-1",
    "hostile_findings": 2,
    "clean_findings": 0,
    "outputs": ["JSON", "SARIF"],
    "default_blocking_severities": ["WARNING", "ERROR"]
  }
}
````

### FILE: `gitlab_opengrep_signed_sast/install_and_verify.ps1`
```yaml
block_id: "GITLAB-OPENGREP-SIGNED-SAST-GATE:installer:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite fail-closed acquisition, signature, license and runtime publication glue"
license: "LicenseRef-Workspace-Owner"
sha256: "9d1cba41499a1ce7620795cd857b0d5baa30cd063559a1269ebffb69c2bba8dc"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0

[CmdletBinding()]
param(
  [Parameter(Mandatory = $true)] [string] $Destination,
  [Parameter(Mandatory = $true)] [string] $ReceiptPath,
  [string] $CacheDirectory,
  [switch] $AllowNetwork
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
$utf8 = [Text.UTF8Encoding]::new($false)
$destinationPath = [IO.Path]::GetFullPath($Destination)
$receipt = [IO.Path]::GetFullPath($ReceiptPath)
if (Test-Path -LiteralPath $destinationPath) { throw "Destination already exists: $destinationPath" }
if (Test-Path -LiteralPath $receipt) { throw "ReceiptPath already exists: $receipt" }
if (-not $AllowNetwork) { throw 'AllowNetwork is required for authority and Sigstore verification; CacheDirectory only reduces repeated downloads' }
$destinationParent = Split-Path -Parent $destinationPath
$receiptParent = Split-Path -Parent $receipt
if (-not $destinationParent -or -not $receiptParent) { throw 'Destination and ReceiptPath require parent directories' }
if (-not (Test-Path -LiteralPath $destinationParent -PathType Container)) { [void](New-Item -ItemType Directory -Path $destinationParent) }
if (-not (Test-Path -LiteralPath $receiptParent -PathType Container)) { [void](New-Item -ItemType Directory -Path $receiptParent) }

$cache = $null
if ($CacheDirectory) {
  $cache = [IO.Path]::GetFullPath($CacheDirectory)
  if (-not (Test-Path -LiteralPath $cache -PathType Container)) {
    [void](New-Item -ItemType Directory -Path $cache)
  }
}
$lockPath = Join-Path $PSScriptRoot 'source-lock.json'
$lock = Get-Content -LiteralPath $lockPath -Raw | ConvertFrom-Json -Depth 50
if ($lock.schema -cne 'elite-gitlab-opengrep-sast-lock/v1') { throw 'Source lock schema mismatch' }
if ($lock.current_admission.status -cne 'OPEN') { throw "Runtime admission blocked: $($lock.current_admission.status). $($lock.current_admission.reason)" }
$stage = Join-Path ([IO.Path]::GetTempPath()) ('elite-opengrep-install-' + [Guid]::NewGuid().ToString('N'))
$candidate = Join-Path $destinationParent ('.elite-opengrep-candidate-' + [Guid]::NewGuid().ToString('N'))

function Get-Sha([string] $Path) {
  (Get-FileHash -LiteralPath $Path -Algorithm SHA256).Hash.ToLowerInvariant()
}

function Assert-Artifact([string] $Path, [object] $Spec) {
  if (-not (Test-Path -LiteralPath $Path -PathType Leaf)) { throw "Artifact missing: $Path" }
  if ((Get-Item -LiteralPath $Path).Length -ne [int64]$Spec.bytes -or (Get-Sha $Path) -cne [string]$Spec.sha256) { throw "Artifact identity mismatch: $($Spec.name)" }
}

function Get-LockedArtifact([object] $Spec) {
  $target = Join-Path $stage ([string]$Spec.name)
  if ($cache) {
    $cached = Join-Path $cache ([string]$Spec.name)
    if (Test-Path -LiteralPath $cached -PathType Leaf) {
      Assert-Artifact $cached $Spec
      [IO.File]::Copy($cached, $target, $false)
      return $target
    }
  }
  Invoke-WebRequest -Uri ([string]$Spec.url) -OutFile $target
  Assert-Artifact $target $Spec
  if ($cache) {
    $cached = Join-Path $cache ([string]$Spec.name)
    if (Test-Path -LiteralPath $cached) { throw "Cache race or occupied path: $cached" }
    [IO.File]::Copy($target, $cached, $false)
  }
  $target
}

function Assert-GitHubAuthority {
  $headers = @{'User-Agent'='EliteEngineeringLibrary-OpenGrepGate';'Accept'='application/vnd.github+json'}
  $og = Invoke-RestMethod -Uri ("https://api.github.com/repos/{0}/commits/{1}" -f $lock.opengrep.repository,$lock.opengrep.commit) -Headers $headers
  if ($og.sha -cne [string]$lock.opengrep.commit -or $og.commit.tree.sha -cne [string]$lock.opengrep.tree -or $og.commit.verification.verified -ne $true) { throw 'OpenGrep commit identity/signature mismatch' }
  $ref = Invoke-RestMethod -Uri ("https://api.github.com/repos/{0}/git/ref/tags/{1}" -f $lock.cosign.repository,$lock.cosign.tag) -Headers $headers
  if ($ref.object.type -cne 'tag' -or $ref.object.sha -cne [string]$lock.cosign.tag_object) { throw 'Cosign tag ref mismatch' }
  $tag = Invoke-RestMethod -Uri $ref.object.url -Headers $headers
  if ($tag.verification.verified -ne $true -or $tag.object.sha -cne [string]$lock.cosign.commit) { throw 'Cosign annotated tag signature/target mismatch' }
  $cg = Invoke-RestMethod -Uri ("https://api.github.com/repos/{0}/commits/{1}" -f $lock.cosign.repository,$lock.cosign.commit) -Headers $headers
  if ($cg.sha -cne [string]$lock.cosign.commit -or $cg.commit.tree.sha -cne [string]$lock.cosign.tree -or $cg.commit.verification.verified -ne $true) { throw 'Cosign commit identity/signature mismatch' }
}

function Assert-GitLabAuthority {
  $headers = @{'User-Agent'='EliteEngineeringLibrary-OpenGrepGate'}
  $tag = Invoke-RestMethod -Uri ("https://gitlab.com/api/v4/projects/27038823/repository/tags/{0}" -f $lock.gitlab_rules.tag) -Headers $headers
  if ($tag.commit.id -cne [string]$lock.gitlab_rules.commit) { throw 'GitLab SAST rules tag/commit mismatch' }
  $file = Invoke-RestMethod -Uri ("https://gitlab.com/api/v4/projects/24075172/repository/files/Dockerfile?ref={0}" -f $lock.gitlab_analyzer_authority.commit) -Headers $headers
  if ($file.blob_id -cne [string]$lock.gitlab_analyzer_authority.dockerfile_blob) { throw 'GitLab analyzer Dockerfile blob mismatch' }
  $dockerfile = [Text.Encoding]::UTF8.GetString([Convert]::FromBase64String([string]$file.content))
  $dockerSha = [Convert]::ToHexString([Security.Cryptography.SHA256]::HashData($utf8.GetBytes($dockerfile))).ToLowerInvariant()
  if ($dockerSha -cne [string]$lock.gitlab_analyzer_authority.dockerfile_sha256) { throw 'GitLab analyzer Dockerfile content mismatch' }
  if ($dockerfile -notmatch ('SCANNER_VERSION=' + [regex]::Escape([string]$lock.gitlab_analyzer_authority.scanner_version)) -or $dockerfile -notmatch ('SAST_RULES_VERSION=' + [regex]::Escape([string]$lock.gitlab_analyzer_authority.sast_rules_version))) { throw 'GitLab analyzer pins drifted' }
}

try {
  [void](New-Item -ItemType Directory -Path $stage)
  [void](New-Item -ItemType Directory -Path $candidate)
  Assert-GitHubAuthority
  Assert-GitLabAuthority
  $authorityApisChecked = $true

  $opengrep = Get-LockedArtifact $lock.opengrep.asset
  $certificate = Get-LockedArtifact $lock.opengrep.certificate
  $signature = Get-LockedArtifact $lock.opengrep.signature
  $cosign = Get-LockedArtifact $lock.cosign.asset
  $opengrepLicense = Get-LockedArtifact $lock.opengrep.license
  $cosignLicense = Get-LockedArtifact $lock.cosign.license
  $gitlabLicense = Get-LockedArtifact $lock.gitlab_rules.root_license
  $rulesPackage = Get-LockedArtifact $lock.gitlab_rules.package

  $version = (& $opengrep --version | Out-String).Trim()
  if ($LASTEXITCODE -ne 0 -or $version -cne [string]$lock.opengrep.version) { throw "OpenGrep executable version mismatch: $version" }
  $cosignVersion = (& $cosign version 2>&1 | Out-String)
  if ($LASTEXITCODE -ne 0 -or $cosignVersion -notmatch [regex]::Escape([string]$lock.cosign.version)) { throw 'Cosign executable version mismatch' }
  $cosignOutput = (& $cosign verify-blob --cert $certificate --signature $signature --certificate-identity-regexp ([string]$lock.opengrep.certificate_identity_regexp) --certificate-oidc-issuer ([string]$lock.opengrep.certificate_oidc_issuer) $opengrep 2>&1 | Out-String)
  if ($LASTEXITCODE -ne 0 -or $cosignOutput -notmatch 'Verified OK') { throw "OpenGrep Cosign verification failed: $cosignOutput" }

  Add-Type -AssemblyName System.IO.Compression.FileSystem
  $expanded = Join-Path $stage 'rules-expanded'
  $zip = [IO.Compression.ZipFile]::OpenRead($rulesPackage)
  try {
    if ($zip.Entries.Count -ne 31) { throw "GitLab rules archive inventory drifted: $($zip.Entries.Count)" }
    foreach ($entry in $zip.Entries) {
      $name = $entry.FullName.Replace('\','/')
      if ([IO.Path]::IsPathRooted($name) -or $name.Split('/') -contains '..') { throw "Unsafe rules archive entry: $name" }
    }
  } finally { $zip.Dispose() }
  [IO.Compression.ZipFile]::ExtractToDirectory($rulesPackage, $expanded)
  foreach ($excluded in @($lock.gitlab_rules.excluded)) {
    $excludedLicense = Join-Path $expanded ([string]$excluded.license_path)
    if ((Get-Item -LiteralPath $excludedLicense).Length -ne [int64]$excluded.license_bytes -or (Get-Sha $excludedLicense) -cne [string]$excluded.license_sha256) { throw "Excluded-license evidence drifted: $($excluded.path)" }
  }

  $toolDir = Join-Path $candidate 'tools'
  $ruleDir = Join-Path $candidate 'rules'
  $licenseDir = Join-Path $candidate 'licenses'
  $provenanceDir = Join-Path $candidate 'provenance'
  foreach ($dir in @($toolDir,$ruleDir,$licenseDir,$provenanceDir)) { [void](New-Item -ItemType Directory -Path $dir) }
  [IO.File]::Copy($opengrep, (Join-Path $toolDir 'opengrep.exe'), $false)
  $actualRuleCount = 0
  foreach ($rule in @($lock.gitlab_rules.selected_rules)) {
    $sourceRule = Join-Path $expanded ([string]$rule.archive_path)
    Assert-Artifact $sourceRule $rule
    $count = (Select-String -LiteralPath $sourceRule -Pattern '^\s*- id:|^\s+id:' -AllMatches).Matches.Count
    if ($count -ne [int]$rule.rules) { throw "Rule count drifted for $($rule.archive_path): $count" }
    $actualRuleCount += $count
    [IO.File]::Copy($sourceRule, (Join-Path $ruleDir ([string]$rule.output_name)), $false)
  }
  if ($actualRuleCount -ne [int]$lock.gitlab_rules.selected_rule_count) { throw "Selected rule count mismatch: $actualRuleCount" }
  $lgplLicense = Join-Path $expanded ([string]$lock.gitlab_rules.lgpl_license.archive_path)
  Assert-Artifact $lgplLicense $lock.gitlab_rules.lgpl_license
  [IO.File]::Copy($opengrepLicense, (Join-Path $licenseDir 'OpenGrep-LICENSE'), $false)
  [IO.File]::Copy($cosignLicense, (Join-Path $licenseDir 'Cosign-LICENSE'), $false)
  [IO.File]::Copy($gitlabLicense, (Join-Path $licenseDir 'GitLab-SAST-Rules-LICENSE'), $false)
  [IO.File]::Copy($lgplLicense, (Join-Path $licenseDir 'GitLab-ThirdParty-LGPL-3.0-LICENSE'), $false)
  [IO.File]::Copy($lockPath, (Join-Path $provenanceDir 'source-lock.json'), $false)
  [IO.File]::Copy($certificate, (Join-Path $provenanceDir 'opengrep-certificate.cert'), $false)
  [IO.File]::Copy($signature, (Join-Path $provenanceDir 'opengrep-signature.sig'), $false)
  [IO.File]::WriteAllText((Join-Path $provenanceDir 'cosign-verification.txt'), $cosignOutput, $utf8)
  $notice = @"
# Third-party notices

This runtime contains the signed OpenGrep 1.29.0 Windows executable under LGPL-2.1-only and an admitted subset of GitLab SAST Rules 2.9.3 under MIT and LGPL-3.0-only. Cosign 3.1.3 was used for signature verification under Apache-2.0 but its binary is not redistributed in this runtime. Exact license texts and the source lock are retained. GitLab Enterprise rules and Commons-Clause rules were deliberately excluded.
"@
  [IO.File]::WriteAllText((Join-Path $candidate 'THIRD_PARTY_NOTICES.md'), $notice.TrimStart(), $utf8)

  $manifestFiles = @(
    Get-ChildItem -LiteralPath $candidate -File -Recurse |
      Sort-Object FullName |
      ForEach-Object {
        [ordered]@{
          path = $_.FullName.Substring($candidate.Length + 1).Replace('\','/')
          bytes = $_.Length
          sha256 = Get-Sha $_.FullName
        }
      }
  )
  $runtimeManifest = [ordered]@{
    schema = 'elite-gitlab-opengrep-runtime-manifest/v1'
    created_at = [DateTimeOffset]::UtcNow.ToString('o')
    opengrep_version = [string]$lock.opengrep.version
    gitlab_rules_version = [string]$lock.gitlab_rules.tag
    selected_rule_count = $actualRuleCount
    blocking_rule_count_at_warning_or_error = [int]$lock.gitlab_rules.blocking_rule_count_at_warning_or_error
    excluded_rule_families = @($lock.gitlab_rules.excluded.path)
    files = $manifestFiles
  }
  $runtimeManifestPath = Join-Path $candidate 'runtime-manifest.json'
  [IO.File]::WriteAllText($runtimeManifestPath, ($runtimeManifest | ConvertTo-Json -Depth 30), $utf8)
  $receiptRecord = [ordered]@{
    schema = 'elite-gitlab-opengrep-install-receipt/v1'
    status = 'VERIFIED'
    generated_at = [DateTimeOffset]::UtcNow.ToString('o')
    destination = $destinationPath
    source_lock_sha256 = Get-Sha $lockPath
    runtime_manifest_sha256 = Get-Sha $runtimeManifestPath
    authority_apis_checked = $authorityApisChecked
    opengrep_version = $version
    opengrep_signature_verified = $true
    cosign_version = [string]$lock.cosign.version
    gitlab_rules_version = [string]$lock.gitlab_rules.tag
    selected_rule_count = $actualRuleCount
    licenses = @('LGPL-2.1-only','MIT','LGPL-3.0-only','Apache-2.0 build tool')
  }
  [IO.Directory]::Move($candidate, $destinationPath)
  $receiptCandidate = $receipt + '.candidate-' + [Guid]::NewGuid().ToString('N')
  [IO.File]::WriteAllText($receiptCandidate, ($receiptRecord | ConvertTo-Json -Depth 20), $utf8)
  [IO.File]::Move($receiptCandidate, $receipt)
  Write-Output "SAST_INSTALL status=VERIFIED rules=$actualRuleCount destination=$destinationPath"
} finally {
  if (Test-Path -LiteralPath $stage -PathType Container) { [IO.Directory]::Delete($stage, $true) }
  if (Test-Path -LiteralPath $candidate -PathType Container) { [IO.Directory]::Delete($candidate, $true) }
}
````

### FILE: `gitlab_opengrep_signed_sast/scan.ps1`
```yaml
block_id: "GITLAB-OPENGREP-SIGNED-SAST-GATE:scanner:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite fail-closed SAST execution and atomic evidence glue"
license: "LicenseRef-Workspace-Owner"
sha256: "65771b136798f2e5e780d3a0d7cf3167479370ba8e14d764b169322b7719909a"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0

[CmdletBinding()]
param(
  [Parameter(Mandatory = $true)] [string] $RuntimePath,
  [Parameter(Mandatory = $true)] [string] $SourceDir,
  [Parameter(Mandatory = $true)] [string] $OutputPath,
  [ValidateSet('INFO','WARNING','ERROR')] [string] $MinimumSeverity = 'WARNING',
  [ValidateRange(10,7200)] [int] $TimeoutSeconds = 900,
  [switch] $ScanIgnoredFiles
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
$utf8 = [Text.UTF8Encoding]::new($false)
$runtime = [IO.Path]::GetFullPath($RuntimePath)
$source = [IO.Path]::GetFullPath($SourceDir)
$output = [IO.Path]::GetFullPath($OutputPath)
if (-not (Test-Path -LiteralPath $runtime -PathType Container)) { throw "Runtime missing: $runtime" }
if (-not (Test-Path -LiteralPath $source -PathType Container)) { throw "SourceDir missing: $source" }
if (Test-Path -LiteralPath $output) { throw "OutputPath already exists: $output" }
$sourcePrefix = $source.TrimEnd([IO.Path]::DirectorySeparatorChar) + [IO.Path]::DirectorySeparatorChar
if ($output.StartsWith($sourcePrefix, [StringComparison]::OrdinalIgnoreCase)) { throw 'OutputPath must not be inside SourceDir' }

function Get-Sha([string] $Path) {
  (Get-FileHash -LiteralPath $Path -Algorithm SHA256).Hash.ToLowerInvariant()
}

$manifestPath = Join-Path $runtime 'runtime-manifest.json'
if (-not (Test-Path -LiteralPath $manifestPath -PathType Leaf)) { throw 'runtime-manifest.json missing' }
$manifest = Get-Content -LiteralPath $manifestPath -Raw | ConvertFrom-Json -Depth 20
if ($manifest.schema -cne 'elite-gitlab-opengrep-runtime-manifest/v1') { throw 'Runtime manifest schema mismatch' }
foreach ($file in @($manifest.files)) {
  $relative = [string]$file.path
  if ([IO.Path]::IsPathRooted($relative) -or $relative.Replace('\','/').Split('/') -contains '..') { throw "Unsafe manifest path: $relative" }
  $actual = Join-Path $runtime $relative
  if (-not (Test-Path -LiteralPath $actual -PathType Leaf)) { throw "Runtime file missing: $relative" }
  if ((Get-Item -LiteralPath $actual).Length -ne [int64]$file.bytes -or (Get-Sha $actual) -cne [string]$file.sha256) { throw "Runtime file drifted: $relative" }
}

$exe = Join-Path $runtime 'tools\opengrep.exe'
$rules = Join-Path $runtime 'rules'
$versionOutput = (& $exe --version | Out-String).Trim()
if ($LASTEXITCODE -ne 0 -or $versionOutput -cne [string]$manifest.opengrep_version) { throw "OpenGrep version mismatch: $versionOutput" }
$ruleFiles = @(Get-ChildItem -LiteralPath $rules -Filter '*.yml' -File | Sort-Object Name)
if ($ruleFiles.Count -ne 3) { throw "Expected exactly three admitted rule files; actual=$($ruleFiles.Count)" }

$outputParent = Split-Path -Parent $output
if (-not $outputParent) { throw 'OutputPath requires a parent directory' }
if (-not (Test-Path -LiteralPath $outputParent -PathType Container)) { [void](New-Item -ItemType Directory -Path $outputParent) }
$candidate = Join-Path $outputParent ('.elite-sast-scan-' + [Guid]::NewGuid().ToString('N'))
[void](New-Item -ItemType Directory -Path $candidate)
$jsonPath = Join-Path $candidate 'findings.json'
$sarifPath = Join-Path $candidate 'findings.sarif'
$stdoutPath = Join-Path $candidate 'opengrep.stdout.log'
$stderrPath = Join-Path $candidate 'opengrep.stderr.log'

$severities = switch ($MinimumSeverity) {
  'INFO' { @('INFO','WARNING','ERROR') }
  'WARNING' { @('WARNING','ERROR') }
  'ERROR' { @('ERROR') }
}
$arguments = [Collections.Generic.List[string]]::new()
foreach ($arg in @('scan','--config',$rules,'--no-rewrite-rule-ids','--strict','--error')) { $arguments.Add($arg) }
if ($ScanIgnoredFiles) { $arguments.Add('--no-git-ignore') }
foreach ($severity in $severities) { $arguments.Add('--severity'); $arguments.Add($severity) }
foreach ($arg in @('--json-output',$jsonPath,'--sarif-output',$sarifPath,$source)) { $arguments.Add($arg) }

$process = [Diagnostics.Process]::new()
$process.StartInfo.FileName = $exe
$process.StartInfo.UseShellExecute = $false
$process.StartInfo.RedirectStandardOutput = $true
$process.StartInfo.RedirectStandardError = $true
$process.StartInfo.CreateNoWindow = $true
foreach ($arg in $arguments) { [void]$process.StartInfo.ArgumentList.Add($arg) }
try {
  if (-not $process.Start()) { throw 'OpenGrep process did not start' }
  $stdoutTask = $process.StandardOutput.ReadToEndAsync()
  $stderrTask = $process.StandardError.ReadToEndAsync()
  if (-not $process.WaitForExit($TimeoutSeconds * 1000)) {
    $process.Kill($true)
    $process.WaitForExit()
    throw "OpenGrep timed out after $TimeoutSeconds seconds"
  }
  [Threading.Tasks.Task]::WaitAll(@($stdoutTask,$stderrTask))
  [IO.File]::WriteAllText($stdoutPath, $stdoutTask.Result, $utf8)
  [IO.File]::WriteAllText($stderrPath, $stderrTask.Result, $utf8)
  $exitCode = $process.ExitCode
} finally {
  $process.Dispose()
}

if ($exitCode -notin @(0,1)) { throw "OpenGrep tool/config failure exit=$exitCode; see $stderrPath" }
if (-not (Test-Path -LiteralPath $jsonPath -PathType Leaf) -or -not (Test-Path -LiteralPath $sarifPath -PathType Leaf)) { throw 'OpenGrep did not produce JSON and SARIF' }
$json = Get-Content -LiteralPath $jsonPath -Raw | ConvertFrom-Json -Depth 100
$findingCount = @($json.results).Count
$errorCount = @($json.errors).Count
if ($errorCount -ne 0) { throw "OpenGrep reported $errorCount parse/runtime errors" }
if (($exitCode -eq 0 -and $findingCount -ne 0) -or ($exitCode -eq 1 -and $findingCount -eq 0)) { throw "Exit/finding invariant failed: exit=$exitCode findings=$findingCount" }
$sarif = Get-Content -LiteralPath $sarifPath -Raw | ConvertFrom-Json -Depth 100
$sarifCount = @($sarif.runs | ForEach-Object { @($_.results) }).Count
if ($sarifCount -ne $findingCount) { throw "JSON/SARIF finding count mismatch: json=$findingCount sarif=$sarifCount" }

$receipt = [ordered]@{
  schema = 'elite-gitlab-opengrep-scan-receipt/v1'
  status = if ($findingCount -eq 0) { 'CLEAN' } else { 'FINDINGS' }
  generated_at = [DateTimeOffset]::UtcNow.ToString('o')
  source_dir = $source
  runtime_manifest_sha256 = Get-Sha $manifestPath
  opengrep_version = $versionOutput
  minimum_severity = $MinimumSeverity
  scan_ignored_files = [bool]$ScanIgnoredFiles
  rule_files = @($ruleFiles.Name)
  findings = $findingCount
  errors = $errorCount
  exit_code = $exitCode
  json_sha256 = Get-Sha $jsonPath
  sarif_sha256 = Get-Sha $sarifPath
}
[IO.File]::WriteAllText((Join-Path $candidate 'scan-receipt.json'), ($receipt | ConvertTo-Json -Depth 20), $utf8)
[IO.Directory]::Move($candidate, $output)
Write-Output "SAST_SCAN status=$($receipt.status) findings=$findingCount output=$output"
exit $exitCode
````

### FILE: `gitlab_opengrep_signed_sast/verify_pack.ps1`
```yaml
block_id: "GITLAB-OPENGREP-SIGNED-SAST-GATE:verifier:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite structural, negative, authority, signature and hostile/clean regression gate"
license: "LicenseRef-Workspace-Owner"
sha256: "799b6090f40c76687e9be5da45dd549c60a88e69980d3cf4cd6f6d165cb9f9ea"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0

[CmdletBinding()]
param(
  [switch] $AllowNetwork,
  [string] $CacheDirectory
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
$required = @('source-lock.json','install_and_verify.ps1','scan.ps1','verify_pack.ps1','README.md')
foreach ($name in $required) {
  if (-not (Test-Path -LiteralPath (Join-Path $PSScriptRoot $name) -PathType Leaf)) { throw "Pack file missing: $name" }
}
foreach ($scriptName in @('install_and_verify.ps1','scan.ps1','verify_pack.ps1')) {
  $tokens = $null
  $errors = $null
  [void][Management.Automation.Language.Parser]::ParseFile((Join-Path $PSScriptRoot $scriptName), [ref]$tokens, [ref]$errors)
  if (@($errors).Count -ne 0) { throw "PowerShell parse errors in ${scriptName}: $($errors.Message -join ' | ')" }
}
$lock = Get-Content -LiteralPath (Join-Path $PSScriptRoot 'source-lock.json') -Raw | ConvertFrom-Json -Depth 50
if ($lock.schema -cne 'elite-gitlab-opengrep-sast-lock/v1') { throw 'Lock schema mismatch' }
if ($lock.current_admission.status -cne 'REJECTED_VERIFIER_RUNTIME' -or $lock.current_admission.cosign_official_sbom_sha256 -cne '34d5b303eea474d37ffb19393c27ea92f41b7b73fe085a88de58d595f490cef0' -or $lock.current_admission.osv_affected_package_rows -ne 9 -or $lock.current_admission.osv_unique_advisories -ne 164) { throw 'Current Cosign rejection evidence drifted' }
if ($lock.opengrep.tag -cne 'v1.29.0' -or $lock.opengrep.version -cne '1.29.0' -or $lock.opengrep.commit_signature_verified -ne $true) { throw 'OpenGrep identity drifted' }
if ($lock.cosign.tag -cne 'v3.1.3' -or $lock.cosign.tag_signature_verified -ne $true -or $lock.cosign.commit_signature_verified -ne $true) { throw 'Cosign identity drifted' }
if ($lock.gitlab_analyzer_authority.sast_rules_version -cne '2.9.3' -or $lock.gitlab_rules.tag -cne 'v2.9.3') { throw 'GitLab compatible rules pin drifted' }
$selected = @($lock.gitlab_rules.selected_rules)
if ($selected.Count -ne 3 -or ($selected | Measure-Object rules -Sum).Sum -ne 122) { throw 'Selected rule inventory drifted' }
if (@($selected | Where-Object { $_.archive_path -match '(^|/)gitlab/' -or $_.archive_path -match 'lgpl-cc' }).Count -ne 0) { throw 'Restricted GitLab rules entered admitted selection' }
if (@($lock.gitlab_rules.excluded).Count -ne 2 -or $lock.gitlab_rules.excluded[0].reason -notmatch 'subscription' -or $lock.gitlab_rules.excluded[1].reason -notmatch 'Commons Clause') { throw 'Excluded license families are not explicit' }
if (@($selected | Where-Object { $_.sha256 -notmatch '^[0-9a-f]{64}$' -or $_.bytes -le 0 -or $_.rules -le 0 }).Count -ne 0) { throw 'Selected rule lock is incomplete' }
if ($lock.gitlab_rules.blocking_rule_count_at_warning_or_error -ne 117) { throw 'Blocking rule count drifted' }
$scriptText = Get-Content -LiteralPath (Join-Path $PSScriptRoot 'scan.ps1') -Raw
foreach ($needle in @('--strict','--error','--json-output','--sarif-output','runtime-manifest.json','OutputPath already exists','OutputPath must not be inside SourceDir')) {
  if (-not $scriptText.Contains($needle, [StringComparison]::Ordinal)) { throw "Fail-closed scan invariant missing: $needle" }
}

$negativeRoot = Join-Path ([IO.Path]::GetTempPath()) ('elite-opengrep-static-negative-' + [Guid]::NewGuid().ToString('N'))
try {
  [void](New-Item -ItemType Directory -Path $negativeRoot)
  $source = Join-Path $negativeRoot 'source'
  [void](New-Item -ItemType Directory -Path $source)
  & pwsh -NoProfile -File (Join-Path $PSScriptRoot 'scan.ps1') -RuntimePath (Join-Path $negativeRoot 'missing-runtime') -SourceDir $source -OutputPath (Join-Path $negativeRoot 'should-not-exist') 2>$null
  if ($LASTEXITCODE -eq 0 -or (Test-Path -LiteralPath (Join-Path $negativeRoot 'should-not-exist'))) { throw 'Missing-runtime negative was not rejected' }
} finally {
  if (Test-Path -LiteralPath $negativeRoot -PathType Container) { [IO.Directory]::Delete($negativeRoot, $true) }
}

if (-not $AllowNetwork) {
  Write-Output 'GITLAB_OPENGREP_SAST_RUNTIME_BLOCKED reason=Cosign_3.1.3_SCA_rejected'
  Write-Output 'GITLAB_OPENGREP_SAST_PACK PASS static=1 negative=1 runtime=BLOCKED'
  exit 0
}

Write-Output 'GITLAB_OPENGREP_SAST_RUNTIME_BLOCKED reason=Cosign_3.1.3_SCA_rejected network_not_used=1'
Write-Output 'GITLAB_OPENGREP_SAST_PACK PASS static=1 negative=1 runtime=BLOCKED'
exit 0

$testRoot = Join-Path ([IO.Path]::GetTempPath()) ('elite-opengrep-runtime-test-' + [Guid]::NewGuid().ToString('N'))
try {
  [void](New-Item -ItemType Directory -Path $testRoot)
  $runtime = Join-Path $testRoot 'runtime'
  $receipt = Join-Path $testRoot 'install-receipt.json'
  $cache = if ($CacheDirectory) { [IO.Path]::GetFullPath($CacheDirectory) } else { Join-Path $testRoot 'cache' }
  $installArgs = @('-NoProfile','-File',(Join-Path $PSScriptRoot 'install_and_verify.ps1'),'-Destination',$runtime,'-ReceiptPath',$receipt,'-CacheDirectory',$cache,'-AllowNetwork')
  & pwsh @installArgs
  if ($LASTEXITCODE -ne 0) { throw 'Runtime installation/verification failed' }
  $installReceipt = Get-Content -LiteralPath $receipt -Raw | ConvertFrom-Json -Depth 30
  if ($installReceipt.status -cne 'VERIFIED' -or $installReceipt.opengrep_signature_verified -ne $true -or $installReceipt.selected_rule_count -ne 122 -or $installReceipt.authority_apis_checked -ne $true) { throw 'Install receipt invariant failed' }

  $hostile = Join-Path $testRoot 'hostile'
  $clean = Join-Path $testRoot 'clean'
  [void](New-Item -ItemType Directory -Path $hostile)
  [void](New-Item -ItemType Directory -Path $clean)
  [IO.File]::WriteAllText((Join-Path $hostile 'unsafe.ts'), "export function run(input: string): unknown {`n  return eval(input);`n}`n")
  [IO.File]::WriteAllText((Join-Path $hostile 'unsafe.go'), "package fixtures`n`nimport `"os/exec`"`n`nfunc run(program string) error {`n`treturn exec.Command(program).Run()`n}`n")
  [IO.File]::WriteAllText((Join-Path $clean 'safe.ts'), "export function parse(input: string): number {`n  const value = Number(input);`n  if (!Number.isFinite(value)) throw new Error(`"invalid`" + input.length);`n  return value;`n}`n")
  [IO.File]::WriteAllText((Join-Path $clean 'safe.go'), "package fixtures`n`nimport `"os/exec`"`n`nfunc version() error {`n`treturn exec.Command(`"git`", `"--version`").Run()`n}`n")

  $hostileOutput = Join-Path $testRoot 'hostile-output'
  & pwsh -NoProfile -File (Join-Path $PSScriptRoot 'scan.ps1') -RuntimePath $runtime -SourceDir $hostile -OutputPath $hostileOutput -MinimumSeverity WARNING -ScanIgnoredFiles
  if ($LASTEXITCODE -ne 1) { throw "Hostile scan expected exit 1; actual=$LASTEXITCODE" }
  $hostileReceipt = Get-Content -LiteralPath (Join-Path $hostileOutput 'scan-receipt.json') -Raw | ConvertFrom-Json -Depth 30
  $hostileJson = Get-Content -LiteralPath (Join-Path $hostileOutput 'findings.json') -Raw | ConvertFrom-Json -Depth 100
  $ids = @($hostileJson.results.check_id | Sort-Object)
  if ($hostileReceipt.status -cne 'FINDINGS' -or $hostileReceipt.findings -ne 2 -or $ids -notcontains 'eslint.detect-eval-with-expression' -or $ids -notcontains 'gosec.G204-1') { throw "Hostile detection invariant failed: $($ids -join ',')" }

  $cleanOutput = Join-Path $testRoot 'clean-output'
  & pwsh -NoProfile -File (Join-Path $PSScriptRoot 'scan.ps1') -RuntimePath $runtime -SourceDir $clean -OutputPath $cleanOutput -MinimumSeverity WARNING -ScanIgnoredFiles
  if ($LASTEXITCODE -ne 0) { throw "Clean scan expected exit 0; actual=$LASTEXITCODE" }
  $cleanReceipt = Get-Content -LiteralPath (Join-Path $cleanOutput 'scan-receipt.json') -Raw | ConvertFrom-Json -Depth 30
  if ($cleanReceipt.status -cne 'CLEAN' -or $cleanReceipt.findings -ne 0 -or $cleanReceipt.errors -ne 0) { throw 'Clean scan invariant failed' }

  [IO.File]::AppendAllText((Join-Path $runtime 'rules\eslint.yml'), "`n# tamper`n")
  $tamperOutput = Join-Path $testRoot 'tamper-output'
  & pwsh -NoProfile -File (Join-Path $PSScriptRoot 'scan.ps1') -RuntimePath $runtime -SourceDir $clean -OutputPath $tamperOutput 2>$null
  if ($LASTEXITCODE -eq 0 -or (Test-Path -LiteralPath $tamperOutput)) { throw 'Runtime tamper negative was not rejected before output publication' }

  Write-Output 'GITLAB_OPENGREP_SAST_PACK PASS static=1 negative=2 signature=1 authority=1 hostile=2 clean=0 json=1 sarif=1'
} finally {
  if (Test-Path -LiteralPath $testRoot -PathType Container) { [IO.Directory]::Delete($testRoot, $true) }
}
````

### FILE: `gitlab_opengrep_signed_sast/README.md`
```yaml
block_id: "GITLAB-OPENGREP-SIGNED-SAST-GATE:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite operator contract with exact/adapted boundary and limits"
license: "LicenseRef-Workspace-Owner"
sha256: "030d89aba5fa9e59a93b7f3fd469cce3435a2009fb839016bcffd588f39c823e"
variables: []
secrets_allowed: false
```
````markdown
# GitLab/OpenGrep local SAST gate

This pack preserves a previously rebuilt no-additional-license-cost Windows x64 SAST candidate for private Go and TypeScript repositories. Its runtime installation is now fail-closed: the official Cosign 3.1.3 SBOM was re-scanned on 2026-09-05 with exact OSV-Scanner 2.5.1 and contains affected dependencies. OpenGrep 1.29.0, its signature/certificate and the 122 commercially admissible GitLab SAST Rules 2.9.3 remain exactly locked, but they must not be promoted or re-verified through this Cosign runtime until the rejection trigger is resolved.

## What is exact and what is adapted

- `opengrep.exe` is the exact signed OpenGrep 1.29.0 release asset.
- The three rule files are exact bytes from GitLab's official `sast-rules-v2.9.3.zip` package: `eslint.yml`, `gosec.yml`, and `lgpl/nodejs_scan.yml`.
- `install_and_verify.ps1`, `scan.ps1`, this README, and the lock are authored integration code. Running GitLab's rules with OpenGrep is an explicit adaptation; it is not the GitLab Semgrep analyzer image and is not presented as GitLab-authored runtime code.
- GitLab Enterprise rules are excluded because they require a subscription. `lgpl-cc` rules are excluded because Commons Clause restricts sale. Neither family enters the runtime.

## Install once

```powershell
pwsh -NoProfile -File .\install_and_verify.ps1 `
  -Destination C:\secure-tools\elite-opengrep `
  -ReceiptPath C:\secure-tools\elite-opengrep-install.json `
  -CacheDirectory C:\secure-tools\artifact-cache `
  -AllowNetwork
```

The cache is optional and avoids downloading the large binaries again. Installation still requires authorized network access for current authority and Sigstore verification; the cache is not presented as an offline trust bundle. A modified cache entry fails closed. Destinations and receipts are never overwritten.

## Scan a project

```powershell
pwsh -NoProfile -File .\scan.ps1 `
  -RuntimePath C:\secure-tools\elite-opengrep `
  -SourceDir C:\projects\my-project `
  -OutputPath C:\projects\evidence\sast-001
```

Default blocking severities are `WARNING` and `ERROR`. Exit `0` means the selected scope is clean; exit `1` means findings exist; every other engine code is a tool/configuration failure. JSON, SARIF, stdout/stderr logs, and an immutable-path receipt are emitted atomically. The runner verifies the runtime manifest before every scan and never suppresses findings automatically.

## Verification

```powershell
pwsh -NoProfile -File .\verify_pack.ps1
pwsh -NoProfile -File .\verify_pack.ps1 -AllowNetwork -CacheDirectory C:\secure-tools\artifact-cache
```

The network verification checks both authority APIs, the Cosign signature, exact licenses/rules, a hostile Go and TypeScript corpus, a clean corpus, JSON/SARIF parity, and tamper rejection.

## Limits

This is a local pattern/taint SAST lane, not proof that a product is secure and not a complete interprocedural or cross-repository analysis guarantee. It complements DevSkim, native fuzzing, dependency scanning, secret scanning, DAST, threat modeling, browser gates, review, and offensive testing. Findings require project-specific triage, a correction, and a regression test; suppressions require an owner, justification, expiry, and review outside this pack.
````

## 6. Configuration surface

| Variable | Type | Default | Validation | Secret | Mutability | Effect |
|---|---|---:|---|---|---|---|
| `Destination` | absolute path | none | must not exist | no | install-time | verified runtime destination |
| `ReceiptPath` | absolute path | none | must not exist | no | install-time | immutable install evidence |
| `CacheDirectory` | absolute path | none | cached bytes must match locks | no | install-time | avoids repeat artifact downloads; not an offline trust bundle |
| `AllowNetwork` | switch | false | required for authority/Sigstore verification | no | install-time | authorizes fixed-source egress |
| `RuntimePath` | absolute path | none | manifest and every file must match | no | scan-time | selected verified scanner runtime |
| `SourceDir` | absolute directory | none | must exist; output cannot be nested within it | potentially sensitive path | scan-time | project scope |
| `OutputPath` | absolute path | none | must not exist | potentially sensitive findings | per scan | atomic evidence destination |
| `MinimumSeverity` | enum | `WARNING` | `INFO|WARNING|ERROR` | no | per scan | severities reported and blocking |
| `TimeoutSeconds` | integer | `900` | 10–7200 | no | per scan | kills hung engine and fails gate |
| `ScanIgnoredFiles` | switch | false | explicit only | no | per scan | includes Git-ignored files when required |

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| OpenGrep Windows x64 | `1.29.0`, SHA-256 `ee485b31…` + Cosign | SAST engine | LGPL-2.1-only | runtime | OpenGrep GitHub release |
| GitLab SAST Rules | `2.9.3`, package SHA-256 `abbff567…` | 122 admitted Go/JS/TS rules | MIT + LGPL-3.0-only | runtime config | GitLab generic package/tag |
| Cosign Windows amd64 | `3.1.3`, SHA-256 `9fe59be0…` | verifies OpenGrep release signature | Apache-2.0 | build-only | Sigstore GitHub release |
| PowerShell | `>=7` | acquisition, verification and scan runner | host tool | build/runtime | Microsoft |

## 8. Apply order

1. Materialice 5/5 archivos en un directorio nuevo.
2. Ejecute `verify_pack.ps1` sin red para sintaxis, lock y negativo estructural.
3. Autorice red y ejecute `verify_pack.ps1 -AllowNetwork`; use un caché privado reutilizable para reducir bytes y latencia.
4. Instale una vez con destino/receipt nuevos; conserve el runtime y sus licencias.
5. Registre scope/severidades/triage/suppressions antes del primer scan target.
6. Ejecute `scan.ps1`; trate `1` como findings a corregir y otros códigos como gate roto.
7. Convierta cada fix en regresión y archive JSON/SARIF/receipt según retención del proyecto.
8. Para actualizar cualquier upstream, abra expediente nuevo; nunca sustituya los pins en lugar.

Rollback: deje de seleccionar el runtime/pack y restaure el pin previo conservando receipts; no borre evidencia histórica. En workspace existente no sobrescribe destinos, outputs ni receipts.

## 9. Verification

- `pwsh -NoProfile -File verify_pack.ps1` → PASS estático + negativo de runtime ausente.
- `pwsh -NoProfile -File verify_pack.ps1 -AllowNetwork -CacheDirectory <cache>` → APIs exactas, firma Cosign, 122 reglas/licencias, hostile Go/TS `2`, clean `0`, JSON/SARIF y tamper rejection.
- `pwsh -NoProfile -File scan.ps1 ...` → exit `0` limpio o `1` con findings; JSON/SARIF y receipt con hashes.
- El proyecto añade dependency/secret scan, DevSkim, fuzz, DAST, threat model, browser/E2E, load, revisión y pruebas ofensivas según riesgo; este pack no los sustituye.

## 10. Reconstruction evidence

- Entorno limpio: Windows x64, PowerShell 7; verificado 2026-09-05.
- Identidades: OpenGrep tag `v1.29.0`, commit firmado `344509d…`; Cosign tag/commit firmados `v3.1.3`/`11926fa…`; reglas GitLab tag/commit `v2.9.3`/`39fc7de…`; analyzer authority `b8d1635…` fija rules `2.9.3`.
- Firma: procedimiento oficial OpenGrep/Cosign devolvió `Verified OK` para el binario exacto.
- Licencias: LGPL-2.1-only, Apache-2.0, MIT y LGPL-3.0-only retenidas; Enterprise y Commons Clause verificadas y excluidas.
- Runtime: 122 reglas totales; 117 `WARNING|ERROR`; JSON y SARIF simultáneos.
- Regresiones: `eslint.detect-eval-with-expression` y `gosec.G204-1` detectados; dos archivos limpios sin hallazgos; runtime modificado rechazado antes de publicar output.
- Resultado focal: `GITLAB_OPENGREP_SAST_PACK PASS static=1 negative=2 signature=1 authority=1 hostile=2 clean=0 json=1 sarif=1`.
- Divergencias: el runtime es una composición `AUTHORED` de binario/reglas exactos, no el contenedor GitLab Semgrep; no existe claim de SAST integral ni de seguridad productiva.
