# Go Native Fuzz Gate

## 1. Metadata

```yaml
pack_id: "GO-NATIVE-FUZZ-GATE"
pack_version: "0.1.0"
status:
  authority: ELITE_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un runner fail-closed para ejecutar regresiones, corpus semilla y fuzzing coverage-guided nativo de Go con presupuesto explícito y receipt hash-linked."
stacks: ["Go 1.26.7+", "PowerShell 7+"]
compatible_with: ["módulo Go con targets FuzzXxx deterministas"]
incompatible_with: ["targets con estado global", "efectos externos", "producción", "presupuesto ilimitado"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev/doc/security/fuzz/", "https://go.dev/doc/tutorial/fuzz", "https://go.dev/doc/security/best-practices"]
verified_at: "2026-09-05"
```

## 2. Applicability

Use cuando un proyecto Go aporta invariantes reales como targets `FuzzXxx`. Rechace si se pretende inventar el oracle, mutar sistemas externos, usar tiempo ilimitado o sustituir DAST/autorización/carga del target.

## 3. Architecture contract

- El proyecto posee código, semillas e invariantes; el gate sólo valida y ejecuta el toolchain oficial.
- Configuración deshabilitada, vacía, duplicada, insegura o sin `go.mod` falla antes del fuzz.
- Primero pasa `go test ./...`, luego el corpus del target y finalmente `go test -fuzz` con tiempo finito.
- El receipt no contiene corpus ni payload y nunca sobrescribe evidencia.
- Un PASS no prueba exhaustividad, ausencia de vulnerabilidades ni seguridad del deployment.

## 4. Exact file manifest

```text
CREATE go_fuzz_gate/profile.example.json
CREATE go_fuzz_gate/run_go_fuzz_gate.ps1
CREATE go_fuzz_gate/verify_pack.ps1
CREATE go_fuzz_gate/README.md
```

## 5. Materialization blocks

### FILE: `go_fuzz_gate/profile.example.json`
```yaml
block_id: "GO-NATIVE-FUZZ-GATE:profile:v1"
operation: CREATE
provenance: AUTHORED
source: "local, governed by official Go fuzzing documentation"
license: "LicenseRef-Workspace-Owner"
sha256: "f7188663149f50c720c54b84d7173533d91711c9baf85444153c354fedfc0fc2"
variables: []
secrets_allowed: false
```
````json
{
  "schema_version": 1,
  "enabled": false,
  "fuzz_time": "10s",
  "targets": []
}
````

### FILE: `go_fuzz_gate/run_go_fuzz_gate.ps1`
```yaml
block_id: "GO-NATIVE-FUZZ-GATE:runner:v1"
operation: CREATE
provenance: AUTHORED
source: "local, governed by official Go fuzzing documentation"
license: "LicenseRef-Workspace-Owner"
sha256: "13513186b96ac0fc3e692a502c92397e1190c41b2ddc44a0eec73dab38f2139a"
variables: []
secrets_allowed: false
```
````powershell
param(
  [Parameter(Mandatory=$true)][string]$ProjectRoot,
  [Parameter(Mandatory=$true)][string]$Profile,
  [Parameter(Mandatory=$true)][string]$Receipt,
  [string]$GoExecutable = 'go'
)
$ErrorActionPreference='Stop'
function Fail([string]$m){ throw "GO_FUZZ_GATE: $m" }
$root=(Resolve-Path -LiteralPath $ProjectRoot).Path
$profilePath=(Resolve-Path -LiteralPath $Profile).Path
if(Test-Path -LiteralPath $Receipt){Fail 'receipt already exists'}
if(-not (Test-Path -LiteralPath (Join-Path $root 'go.mod') -PathType Leaf)){Fail 'go.mod missing'}
$cfg=Get-Content -LiteralPath $profilePath -Raw | ConvertFrom-Json -Depth 20
if($cfg.schema_version -ne 1 -or $cfg.enabled -ne $true){Fail 'profile not explicitly enabled'}
if($cfg.fuzz_time -notmatch '^(?:[1-9][0-9]*)(?:s|m)$'){Fail 'fuzz_time must be positive seconds or minutes'}
$targets=@($cfg.targets)
if($targets.Count -lt 1){Fail 'at least one target required'}
$seen=@{}
foreach($t in $targets){
  if($t.package -notmatch '^\.(?:\/[A-Za-z0-9_-]+)*(?:\/\.\.\.)?$'){Fail 'unsafe package path'}
  if($t.name -notmatch '^Fuzz[A-Za-z0-9_]+$'){Fail 'invalid fuzz target'}
  $key="$($t.package)|$($t.name)"; if($seen.ContainsKey($key)){Fail 'duplicate target'}; $seen[$key]=$true
}
$version=(& $GoExecutable version 2>&1 | Out-String).Trim()
if($LASTEXITCODE -ne 0 -or $version -notmatch '^go version go([0-9]+\.[0-9]+(?:\.[0-9]+)?) '){Fail 'Go unavailable or version unreadable'}
$results=@()
Push-Location $root
try{
  & $GoExecutable test ./...
  if($LASTEXITCODE -ne 0){Fail 'baseline tests failed'}
  foreach($t in $targets){
    $exact='^'+[regex]::Escape([string]$t.name)+'$'
    & $GoExecutable test ([string]$t.package) "-run=$exact"
    if($LASTEXITCODE -ne 0){Fail "seed corpus failed: $key"}
    & $GoExecutable test ([string]$t.package) "-run=$exact" "-fuzz=$exact" "-fuzztime=$($cfg.fuzz_time)"
    if($LASTEXITCODE -ne 0){Fail "fuzz failed: $key"}
    $results += [ordered]@{package=[string]$t.package;name=[string]$t.name;status='PASS'}
  }
} finally {Pop-Location}
$hashes=[ordered]@{profile=(Get-FileHash -LiteralPath $profilePath -Algorithm SHA256).Hash.ToLowerInvariant();go_mod=(Get-FileHash -LiteralPath (Join-Path $root 'go.mod') -Algorithm SHA256).Hash.ToLowerInvariant();go_sum=$null}
$sum=Join-Path $root 'go.sum'; if(Test-Path -LiteralPath $sum){$hashes.go_sum=(Get-FileHash -LiteralPath $sum -Algorithm SHA256).Hash.ToLowerInvariant()}
$body=[ordered]@{schema_version=1;gate='GO-NATIVE-FUZZ-GATE';go_version=$version;fuzz_time=[string]$cfg.fuzz_time;hashes=$hashes;targets=$results}
$parent=Split-Path -Parent $Receipt; if($parent -and -not (Test-Path -LiteralPath $parent)){[void](New-Item -ItemType Directory -Path $parent)}
$tmp="$Receipt.tmp-$([guid]::NewGuid().ToString('N'))"
try{[IO.File]::WriteAllText($tmp,(ConvertTo-Json $body -Depth 20)+"`n",[Text.UTF8Encoding]::new($false)); Move-Item -LiteralPath $tmp -Destination $Receipt}catch{if(Test-Path $tmp){Remove-Item -LiteralPath $tmp -Force};throw}
Write-Output "GO_NATIVE_FUZZ_GATE_PASS targets=$($results.Count)"
````

### FILE: `go_fuzz_gate/verify_pack.ps1`
```yaml
block_id: "GO-NATIVE-FUZZ-GATE:verifier:v1"
operation: CREATE
provenance: AUTHORED
source: "local, governed by official Go fuzzing documentation"
license: "LicenseRef-Workspace-Owner"
sha256: "c1fcc95390f838434c031a270be7bce76e55835379a285331b3b8554547a282e"
variables: []
secrets_allowed: false
```
````powershell
param([Parameter(Mandatory=$true)][string]$GoExecutable)
$ErrorActionPreference='Stop'
$work=Join-Path $env:TEMP ("elite-go-fuzz-"+[guid]::NewGuid().ToString('N')); New-Item -ItemType Directory -Path $work | Out-Null
$module=Join-Path $work 'module'; New-Item -ItemType Directory -Path $module | Out-Null
[IO.File]::WriteAllText((Join-Path $module 'go.mod'),"module example.com/fuzzgate`n`ngo 1.26.0`n",[Text.UTF8Encoding]::new($false))
[IO.File]::WriteAllText((Join-Path $module 'roundtrip_test.go'),@'
package fuzzgate
import "testing"
func FuzzRoundTrip(f *testing.F){f.Add("elite");f.Fuzz(func(t *testing.T,s string){if string([]byte(s))!=s{t.Fatal("roundtrip")}})}
'@,[Text.UTF8Encoding]::new($false))
$profile=Join-Path $work 'profile.json'; [IO.File]::WriteAllText($profile,'{"schema_version":1,"enabled":true,"fuzz_time":"1s","targets":[{"package":"./...","name":"FuzzRoundTrip"}]}',[Text.UTF8Encoding]::new($false))
& (Join-Path $PSScriptRoot 'run_go_fuzz_gate.ps1') -ProjectRoot $module -Profile $profile -Receipt (Join-Path $work 'receipt.json') -GoExecutable $GoExecutable
if($LASTEXITCODE -ne 0){throw 'positive failed'}
$bad=Join-Path $work 'bad.json'; [IO.File]::WriteAllText($bad,'{"schema_version":1,"enabled":false,"fuzz_time":"1s","targets":[]}',[Text.UTF8Encoding]::new($false))
$failed=$false; try{& (Join-Path $PSScriptRoot 'run_go_fuzz_gate.ps1') -ProjectRoot $module -Profile $bad -Receipt (Join-Path $work 'bad-receipt.json') -GoExecutable $GoExecutable}catch{$failed=$true}
if(-not $failed){throw 'disabled negative accepted'}
$unsafe=Join-Path $work 'unsafe.json'; [IO.File]::WriteAllText($unsafe,'{"schema_version":1,"enabled":true,"fuzz_time":"1s","targets":[{"package":"../x","name":"FuzzRoundTrip"}]}',[Text.UTF8Encoding]::new($false))
$failed=$false; try{& (Join-Path $PSScriptRoot 'run_go_fuzz_gate.ps1') -ProjectRoot $module -Profile $unsafe -Receipt (Join-Path $work 'unsafe-receipt.json') -GoExecutable $GoExecutable}catch{$failed=$true}
if(-not $failed){throw 'unsafe path negative accepted'}
Write-Output 'GO_NATIVE_FUZZ_PACK_PASS positives=1 negatives=2'
````

### FILE: `go_fuzz_gate/README.md`
```yaml
block_id: "GO-NATIVE-FUZZ-GATE:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "local, governed by official Go fuzzing documentation"
license: "LicenseRef-Workspace-Owner"
sha256: "2d4c77a7562f439c4c923670291b48ce28c65bd938953259d9846dfaee1313f9"
variables: []
secrets_allowed: false
```
````markdown
# Go native fuzz gate

Configure at least one real project `FuzzXxx` target, explicitly enable the profile and run `verify_pack.ps1` first. The gate runs baseline tests, seed corpora and coverage-guided fuzzing with a finite budget, then writes a hash-linked receipt without overwriting prior evidence. It does not invent business invariants, prove API authorization or replace DAST against the deployed target.
````

## 6. Configuration surface

| Variable | Tipo | Default seguro | Validación | Secreto | Mutabilidad | Efecto |
|---|---|---|---|---|---|---|
| `enabled` | bool | `false` | debe ser `true` explícito | no | por corrida | habilita gate |
| `fuzz_time` | duración | `10s` | entero positivo + `s|m` | no | por corrida | presupuesto por target |
| `targets` | lista | vacía/rechazada | package relativo seguro + `FuzzXxx` único | no | por revisión | oracles ejecutados |

## 7. Dependency bill

| Tool | Pin | Uso | Licencia | Scope | Fuente |
|---|---|---|---|---|---|
| Go | 1.26.7+ observado | test/fuzz coverage-guided | BSD-3-Clause | build/test | <https://go.dev/doc/security/fuzz/> |
| PowerShell | 7+ | validación/receipt | MIT | tooling | <https://github.com/PowerShell/PowerShell> |

## 8. Apply order

Materializar en destino dedicado, copiar el ejemplo, declarar targets reales, habilitarlo y ejecutar primero `verify_pack.ps1`. En proyecto existente no sobrescribir archivos. Rollback: retirar el directorio si aún no hay receipts requeridos por release.

## 9. Verification

Parsear ambos scripts; ejecutar `verify_pack.ps1 -GoExecutable <path>`. Éxito: baseline, seed y fuzz real PASS; perfil deshabilitado y traversal rechazados; receipt nuevo. En el proyecto ejecutar además cada corpus de regresión, `-race`, SCA, DAST y gates de seguridad correspondientes.

## 10. Reconstruction evidence

Windows x64, Go 1.26.7 y PowerShell 7: cuatro archivos materializados; parser PASS; positivo ejecutó 956.687 inputs aproximados en un segundo; dos negativos PASS. Evidencia: `reconstruction_evidence/GO_NATIVE_FUZZ_GATE_2026-09-05_V249.md`.
