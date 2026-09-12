# V285 — entrada medida, compositor estricto y alcance de mantenimiento

Fecha: 2026-09-07. T2801, sin segundo roadmap. Código de control y probe AUTHORED;
no nuevo código upstream adquirido ni promoción de la franquicia.

## Cambio ejecutable

MARKDOWN-COMPOSITOR 0.2.1 desempaquetaba una selección de un elemento al salir de
un if; Count sobre el string fallaba bajo StrictMode heredado. El perfil readiness
real y la suite reprodujeron PropertyNotFoundException antes de escribir.
0.2.2 fija el tipo string[] de selectedFiles. No desactiva StrictMode ni cambia
admisión, archivos seleccionados, checksums, licencias, defaults o permisos.

La suite canónica corre ahora con StrictMode Latest e incluye wildcard, vacío,
null y un path explícito, además de integridad, condiciones, unicidad, colisiones,
operación CREATE-only y rutas reservadas. PASS tras reconstruir desde Markdown;
3/3 archivos idénticos al candidato. Runtime observado: PowerShell 7.6.5,
.NET 10.0.11, Python 3.14.4, Windows 11 build 26200. Latest aquí identifica el modo
estricto de ese runtime observado, no autoriza actualizar el runtime.

Autoridad estrecha: [Microsoft Set-StrictMode](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.core/set-strictmode?view=powershell-7.6),
consultada 2026-09-07, documenta herencia a scopes hijos y rechazo de propiedades
inexistentes. La corrección es nuestra, no una implementación copiada de Microsoft.

## Entrada real y límites de la medición

| Run | Componer perfil (ms) | Componer + copiar template + advisory A + validar (ms) | Archivos conservados tras reentrada rechazada |
|---|---:|---:|---:|
| 1 | 190.91 | 422.03 | 13 |
| 2 | 69.20 | 296.72 | 13 |
| 3 | 83.54 | 315.24 | 13 |

Tres directorios inicialmente ausentes; compositor reconstruido; ocho archivos
del gate y receipt; sin Git, cuentas, llamadas cloud o pagos. El gate inicial
devuelve exit 2 / BLOCKED porque no hay respuestas/aprobaciones. Es la respuesta
correcta, no un fallo que deba maquillarse. La reentrada se rechaza antes de
cambiar inventario o SHA-256. El filesystem de prueba contiene también pycache;
sus 13 archivos no se llaman 13 fuentes materializables.

70 tests de readiness PASS en el último árbol (2.497 s). Estas tres mediciones
usan proceso local y caches existentes; no son p95/p99, ahorro de tokens probado,
latencia de lectura humana ni tiempo de ensamblaje de un sistema completo.
V284 aporta prueba complementaria del bridge con preservación y rollback.

## Preparación reconciliada, no aprobación ficticia

PROJECT_READINESS_GATE.json clasifica las 48 superficies para mantenimiento:
21 REQUIRED y 27 NONE_WITH_REASON, con motivos individuales. IDENTITY y SECRETS-PKI
incluyen operador y futuro firmante/custodia del release, aún no demostrados.
Las 21 REQUIRED conservan condiciones de evidencia abiertas. Cada referencia
general es evidencia parcial del proceso; no afirma que los tests del compositor
prueben identidad, claves, SCA, recuperación o cada capability.

Los 27 NONE_WITH_REASON excluyen aplicaciones/servicios no ejecutados en este
slice local, no capacidades de franquicia. El perfil integral permanece 67/742
y T2802–T2810 siguen vigentes. No se borró frontend, datos, documentos, leads, IA,
operaciones ni sus gates para obtener un porcentaje favorable.

Ronda B queda ANSWERED usando el mandato existente y matriz de actores/permisos,
no PROVEN. Cuatro registros locales derivados de sus templates ya contienen
identidades, consumers, soporte, licencias, políticas, before/after y pendientes:
dependencias, freshness, monitoreo y adquisición oficial. No se adquirieron
fuentes de producto, se ejecutó OSV, creó Git o registró un scheduler por crearlos.
El monitoreo permanece ZERO_COST_BASELINE_BLOCKED; soporte de rama no prueba SCA.

## Identidades y errores conservados

- Owner 0.2.1 previo SHA: de83abe3257ba0cbfeb6a4c6461dc4ad2aba91602bd6ac0dbe044171e4a1444e.
- Owner 0.2.2 SHA: a8505384bb3f9085d36f21f1b6471fe24cbd17b054fa87f82f1d24b7007ce9d5.
- Compositor ejecutable SHA: 1e286714c7122116835264c7e28fc97cde1568eac59f6e03d8e262bc71c40fc4.
- Test canónico SHA: 7e57209c5ccb8ad4cf5c4d337329855e409eef7df1633d3b88a85df6e0a49dcb.

FAIL-393: comando diagnóstico con foreach seguido de pipe rechazado por parser;
se corrigió capturando primero la colección, sin ejecutar mutaciones parciales.
Otra lectura diagnóstica supuso un MATERIALIZATION_RECORD del bootstrap simple
que no existe; se corrigió comparando el inventario observado de tres archivos.
No se inventa un receipt del bootstrap: el receipt sí lo produce el compositor.
Dos URLs documentales propuestas devolvieron error; se resolvieron autoridades
desde repositorios/búsqueda oficial y ninguna URL fallida gobierna un claim.
FAIL-394: defecto funcional StrictMode cerrado por suite y reconstrucción.
La salida recortada de algunos comandos diagnósticos no se usó como inventario
completo ni como evidencia de PASS.

Staging local: elite-v285-d22866c1c7f34df1bfd5a95f3a319798. Conserva probe,
logs de red/green, tres entradas, reconstrucciones y snapshot previo del reporte.
No se emitió ZIP final. El preflight integral y cualquier prueba posterior deben
conservar su propio resultado; estos tests focales no se atribuyen a ese gate.

Reporte actualizado del proyecto: BLOCKED, 44 observaciones (antes 87). Elimina
clasificaciones/archivos ausentes como tales; conserva condiciones de las 21
REQUIRED, operaciones, rondas C–H y assurance no demostrado. No es un porcentaje.
La lectura diagnóstica de README.md también falló: el tercer archivo real es
examples/PROJECT_PACK_PLAN.example.md. La comparación 3/3 usó el inventario
observado, no ese nombre supuesto; el fallo no altera las pruebas ejecutadas.

## Probe reproducible AUTHORED

### Verificación transversal posterior

VERIFY_LIBRARY_PASS: 160 packs / 1433 archivos / 718 Markdown / 51 perfiles.
El primer preflight falló porque VERIFY_EXECUTABLE_LIBRARY.ps1 aún exigía 62
regresiones en vez de las 70 reales. FAIL-395; se corrigieron nombre de paso y
aserción exacta, sin quitar detección de drift. Rerun completo: paso
project-readiness-syntax-and-70-regressions PASS, 2955 ms; exit del runner 0,
estado de disponibilidad EXECUTABLE_PREFLIGHT_BLOCKED missing=4.

Herramientas detectadas: PowerShell 7.6.5, Python 3.14.4, Node v24.14.1 y pnpm
11.25.0. Go, SDK dotnet, Docker y psql no se detectaron con resolución por defecto;
eso no borra evidencia de instalaciones aisladas anteriores. Available es un
check mínimo, no promoción de pnpm 11.25.0 sobre perfiles con pin 11.19.0.
Gates runtime/live omitidos o bloqueados permanecen explícitos en el log: no son
PASS de SAST, cloud, proveedor, corpus, PostgreSQL ni producción.

Reporte preflight-final.json local SHA-256:
88af78b5a4526c05ba6d2bab422eae595de4ff892e0cbcc3f95fbe286591be4b.
El snapshot del reporte JSON contiene paths locales y permanece en staging;
se preservan aquí resultado, identidades no sensibles, alcance y hash.

### Ejecución del probe

Al sellar este informe, el cursor next_action de 815 caracteres excedió su límite
de 800 y fue rechazado (FAIL-396). El intento no agregó EVT-0011; una verificación
iniciada después sin comprobar exit también falló y no cuenta como PASS. Se
resumió el cursor sin elevar el límite, manteniendo el detalle en este expediente;
la prevalidación semántica pasó y el checkpoint posterior debe verificar la cadena.

Tras materializar MARKDOWN_COMPOSITOR_CORE 0.2.2 en rebuilt-compositor junto al
script, ejecutar con LibraryRoot apuntando a la copia canónica. Exige destinos
new-1..3 ausentes; no borrar un proyecto para satisfacer esa condición.

```powershell
# AUTHORED diagnostic probe; not a product readiness fixture or upstream code.
param([Parameter(Mandatory)][string]$LibraryRoot)
Set-StrictMode -Version Latest
$ErrorActionPreference='Stop'
$compositor=Join-Path $PSScriptRoot 'rebuilt-compositor/tools/compose-markdown-project.ps1'
$runs=@()
foreach($index in 1..3){
  $target=Join-Path $PSScriptRoot "new-$index"
  if(Test-Path -LiteralPath $target){throw 'probe requires absent destination'}
  $watch=[Diagnostics.Stopwatch]::StartNew()
  & $compositor -PlanFile (Join-Path $LibraryRoot 'markdown_system/PROJECT_READINESS_GATE_PACK_PLAN.md') -LibraryRoot $LibraryRoot -Destination $target | Out-Null
  $compositionMs=$watch.Elapsed.TotalMilliseconds
  Copy-Item -LiteralPath (Join-Path $target 'project_readiness_gate/project-readiness.template.json') -Destination (Join-Path $target 'PROJECT_READINESS_GATE.json')
  & python (Join-Path $target 'project_readiness_gate/render_project_advisory.py') --project-root $target --round A --output PROJECT_ADVISORY_A.md | Out-Null
  if($LASTEXITCODE -ne 0){throw 'advisory render failed'}
  & python (Join-Path $target 'project_readiness_gate/validate_project_readiness.py') --project-root $target --report PROJECT_READINESS_REPORT.json *> (Join-Path $PSScriptRoot "initial-$index.log")
  if($LASTEXITCODE -ne 2){throw 'fresh template did not fail closed'}
  $watch.Stop()
  if((Get-Content -LiteralPath (Join-Path $target 'PROJECT_READINESS_REPORT.json') -Raw | ConvertFrom-Json).status -ne 'BLOCKED'){throw 'unexpected readiness status'}
  if(Test-Path -LiteralPath (Join-Path $target '.git')){throw 'unexpected Git requirement'}
  $snapshot=@{}
  foreach($file in Get-ChildItem -LiteralPath $target -File -Recurse){$snapshot[$file.FullName]=(Get-FileHash -LiteralPath $file.FullName).Hash}
  $rejected=$false
  try{ & $compositor -PlanFile (Join-Path $LibraryRoot 'markdown_system/PROJECT_READINESS_GATE_PACK_PLAN.md') -LibraryRoot $LibraryRoot -Destination $target | Out-Null }catch{$rejected=$true}
  if(-not $rejected){throw 'existing destination unexpectedly overwritten'}
  foreach($path in $snapshot.Keys){if((Get-FileHash -LiteralPath $path).Hash -cne $snapshot[$path]){throw "existing file changed: $path"}}
  if(@(Get-ChildItem -LiteralPath $target -File -Recurse).Count -ne $snapshot.Count){throw 'existing file inventory changed'}
  $runs+=[ordered]@{run=$index;composition_ms=[math]::Round($compositionMs,2);entry_to_blocked_ms=[math]::Round($watch.Elapsed.TotalMilliseconds,2);existing_preserved_files=$snapshot.Count;initial_status='BLOCKED'}
}
& python -B -m unittest discover -s (Join-Path $PSScriptRoot 'new-3/project_readiness_gate') -p 'test_*.py' *> (Join-Path $PSScriptRoot 'readiness-tests.log')
if($LASTEXITCODE -ne 0){throw 'fresh readiness tests failed'}
$runs | ConvertTo-Json
```

## Snapshot del reporte anterior V283/V284

Se conserva íntegro antes de emitir el reporte actualizado; sus 87 observaciones
eran rechazos del gate, no 87 features ni un porcentaje de avance.

```json
{
  "capabilities_expected": 48,
  "errors": [
    "at least one delivery surface must be REQUIRED or the selected slice must have proven CLI/AUTOMATION journeys",
    "authorities.status must be PROVEN",
    "blockers.critical_unknowns must be empty",
    "blockers.open_high_or_critical_failures must be empty",
    "capabilities.ADS-ATTRIBUTION.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.ANALYTICS-BI.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.API-BACKEND.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.ARCH-DOMAIN must be REQUIRED for a buildable vertical slice",
    "capabilities.ARCH-DOMAIN.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.AUTHORIZATION.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.BACKUP-DR.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.BROKER-STREAMING.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.CACHE.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.CD-RELEASE.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.CI must be REQUIRED for a buildable vertical slice",
    "capabilities.CI.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.CONTAINERS.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.CONTRACTS must be REQUIRED for a buildable vertical slice",
    "capabilities.CONTRACTS.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.COST-FINOPS.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.DATA-INGEST.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.DESKTOP.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.DOCS-OPS must be REQUIRED for a buildable vertical slice",
    "capabilities.DOCS-OPS.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.DOMAIN-MODULES must be REQUIRED for a buildable vertical slice",
    "capabilities.DOMAIN-MODULES.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.EMBEDDED-IOT.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.GPU-ACCEL.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.IAC-CLOUD.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.IDENTITY.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.INTEGRATIONS.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.JOBS-WORKFLOWS.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.MARKETPLACES.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.ML-AI.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.MOBILE.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.NETWORK-EDGE.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.NOTIFICATIONS.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.OBJECT-STORAGE.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.OBSERVABILITY.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.ORCHESTRATION.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.OUTBOX-INBOX.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.PAYMENTS.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.PERFORMANCE.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.PRD-INTAKE must be REQUIRED for a buildable vertical slice",
    "capabilities.PRD-INTAKE.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.PRIVACY-COMPLIANCE.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.RAG-AGENTS.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.REPO-SCM must be REQUIRED for a buildable vertical slice",
    "capabilities.REPO-SCM.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.RUNTIME must be REQUIRED for a buildable vertical slice",
    "capabilities.RUNTIME.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.SEARCH.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.SECRETS-PKI.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.SECURITY-APPSEC must be REQUIRED for a buildable vertical slice",
    "capabilities.SECURITY-APPSEC.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.SLO-INCIDENT.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.SUPPLY-CHAIN must be REQUIRED for a buildable vertical slice",
    "capabilities.SUPPLY-CHAIN.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.TEST-PLATFORM must be REQUIRED for a buildable vertical slice",
    "capabilities.TEST-PLATFORM.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.TX-DATABASE.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.WEB-PORTALS.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.WEB-PUBLIC.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "connected_journeys must be a non-empty array",
    "dependencies.status must be PROVEN",
    "operations.backup_restore_dr.status must be PROVEN",
    "operations.deployment_canary_rollback.status must be PROVEN",
    "operations.monitoring_incident_response.status must be PROVEN",
    "operations.security_privacy.status must be PROVEN",
    "operations.slo_performance_capacity_cost.status must be PROVEN",
    "pack_plan must be PROVEN with collision_check and rollback_defined true",
    "pack_plan.evidence: at least one evidence path is required",
    "project.platform_mode must be an implemented non-BLOCK choice",
    "project.status must be READY_TO_BUILD",
    "required_artifact: evidence file missing: PROJECT_AUTHORITY_FRESHNESS_RECORD.md",
    "required_artifact: evidence file missing: PROJECT_DEPENDENCY_UPDATE_RECORD.md",
    "required_artifact: evidence file missing: PROJECT_OFFICIAL_SOURCE_PROFILE_RECORD.md",
    "required_artifact: evidence file missing: PROJECT_VULNERABILITY_MONITORING_RECORD.md",
    "rounds.B.status must be ANSWERED or PROVEN",
    "rounds.C.status must be ANSWERED or PROVEN",
    "rounds.D.status must be ANSWERED or PROVEN",
    "rounds.E.status must be ANSWERED or PROVEN",
    "rounds.F.status must be ANSWERED or PROVEN",
    "rounds.G.status must be ANSWERED or PROVEN",
    "rounds.H.status must be ANSWERED or PROVEN",
    "sources must be PROVEN with exact revisions and resolved licenses/notices",
    "sources.product_provenance_modes must be non-empty"
  ],
  "record": "PROJECT_READINESS_GATE.json",
  "schema": "elite-project-readiness-report/v1",
  "status": "BLOCKED"
}
```
