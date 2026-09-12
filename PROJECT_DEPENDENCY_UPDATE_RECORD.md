# Project Dependency Update Record — mantenimiento

## Resolución V319 — runtime independiente CONDITIONED, ZIP rechazado

Candidato local24.20.0 con identidad/firma/licencia y gates funcionales probados,
sin promoción. node.exe independiente oficial SHA256
5c976096e04e5c2c1f091938926234cc9fbebfe9787ddd149351b3b0ecc707b5,
con LICENSE íntegra, sólo2 archivos; npm del ZIP no incorporado. ZIP completo
REJECTED:145 paquetes/4 afectados/9 advisories. Consulta Node core puntual con
hashes194 avisos:24.14.1→22CVEs,24.20.0→0; no cobertura OSV por alias npm/node.
Dos árboles115PASS/1SKIP+build+92browser+agenda4,746/746 idénticos; smoke256GET,
8arranques/retorno local PASS. pnpm11.19.0 explícito:891 archivos iguales al tarball;
SCA de bundle compilado pendiente. Gates de Node usarán candidato explícito.
Global24.14.1 retenido, reachability targetUNKNOWN; ningún despliegue ni aceptación.
El checker oficial adicional es evaluaciónADAPTED/CONDITIONED, no pack admitido
ni grupo de SDKs de producto; PROJECT_OFFICIAL_SOURCE_PROFILE_RECORD no habilita
fuentes productivas nuevas. Monitor/freshness/admisión reusable pendientes.
Expediente reconstruction_evidence/NODE_RUNTIME_SECURITY_V319.md. La evaluación
inicial siguiente se preserva como before; no representa el resultado final.

## V319 — Node24 runtime candidate, OFFICIAL_FIXED_RELEASE

FAIL511: Node instalado24.14.1, SHA256
58e74bf02fc5bbacc41dcb8bef089961cd5bddd37830b87784e4fc624d145d1f,
es anterior a fixes oficiales de junio/julio2026. Reachability: UNKNOWN;
promotion bloqueada. El scan npm previo no cubría Node ni sus componentes nativos.

Candidato aislado24.20.0 Windows x64, commit Node
71b8b174857e25106d39b61a9e6f30d927da8b01; ZIP oficial SHA256
6cac9ffbca8f6a47091e4b5c772e0606049c3871cb67d900c0cedde630e545ba.
URL https://nodejs.org/dist/v24.20.0/node-v24.20.0-win-x64.zip.
Source/consumer: TS-GO-API-WEB-BRIDGE Node24, OIDC, portales y gates Node
seleccionados; inventario completo de manifests/engines pendiente. Owner: agente
de mantenimiento. Política: minor-auto-candidate por vulnerabilidad upstream,
no promoción. MIT y notices incorporados a contrastar antes del PASS.

Se evalúa un runtime de reemplazo mediante DEPENDENCY_UPDATE_CONTRACT, sin
seleccionar otro stack ni un grupo de repositorios/SDKs de producto. Metadata de
release, clave pública y advisories son datos de verificación, no código incorporado.
Inputs técnicos: Windows x64 local, Node24 ya requerido, sin cuenta/token/gasto,
sin instalación global ni despliegue. No requiere decisión comercial nueva.
Rollback de prueba: conservar ejecutable24.14.1 y artefactos V317, cambiar sólo
ruta del proceso aislado; nunca promover el runtime vulnerable como rollback seguro.
Gates pendientes: firmas/identidad, componentes/advisories/licencias, frozen install,
web/tests/crypto/4browsers, performance y rebuild/retorno al artefacto exacto.
Staging elite-v319-952f1100edd1424f842acd958e3e1f75; estado ASSESSING.


## V317 — identidad y alcance del runtime Go

Baseline Go1.26.7 en official-toolchain V281, ejecutable SHA256
5463fe58fa999d74420f00ee1b36d31c3da90a57ff204159f159859567ad61fc;
candidate1.26.8 en current-toolchain, SHA256
21761eceb9302062c9623fb699f332c8c7fe000f15f70efe8da01a2cfbbc16b9.
Ambos árboles íntegros contra ZIP oficial y metadata de descarga fresca.
No nueva adquisición ni promoción1.26.8. Retiradas atribuciones históricas exactas
no probadas; V313 buildinfo demuestra1.26.8. Gates del head reejecutados con1.26.7
observado; resultado/receipts en TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md.
Piso explícito x/mod0.40.0 y locks V313 permanecen; no tidy ni nuevos upstreams.
OSV2.5.1 cubre stdlib/toolchain2 paquetes por versión y0matches conocidos;
no completa SCA de todos los runtimes ni autoriza release/producción.


Derivado de markdown_system/PROJECT_DEPENDENCY_UPDATE_RECORD_TEMPLATE.md.
Contrato: markdown_system/DEPENDENCY_UPDATE_CONTRACT.md. Scope: T2801 y herramientas
locales; no inventario aprobado de todos los consumers del perfil de franquicia.

## 1. Control

```yaml
record_version: "1.0"
project_id: elite-library-maintenance
owner: agente_de_mantenimiento
updated_at: "2026-09-07"
baseline_artifact_digest: "PROJECT_EXTERNAL_SOURCE_LOCK.md"
baseline_sbom: PENDING
baseline_license_report: "LICENSE.md; THIRD_PARTY_NOTICES.md"
open_security_updates: ["UP-FAIL-203..209: consultar ledger, no cerrar por este expediente"]
open_eol_updates: []
blocked_updates: ["SCA completo de runtimes y consumers no demostrado"]
active_candidate: MARKDOWN-COMPOSITOR-0.2.2
last_full_review: NOT_COMPLETED
next_review: "antes de cambiar un runtime/lock o promover release; mensual para tooling sin cambios"
```

## 2. Inventario observado

| Componente | Scope / identidad | Fuente/licencia | Consumers / soporte | Policy |
|---|---|---|---|---|
| PowerShell | tooling 7.6.5; pwsh.exe SHA 362a356ce7f0940ec74f73a8fc2c990a2cc24a38a11c90bbd8eca947110ad139 | Microsoft; LICENSE.txt MIT y ThirdPartyNotices.txt locales | materializador/compositor/bridge/verificadores; rama 7.6 hasta 2028-11-14 según fuente oficial consultada | patch-auto-candidate, sin modificar runtime compartido |
| .NET | host de PowerShell, observado 10.0.11 | Microsoft; incluido en avisos del host | no se ha validado aquí todo su grafo o instalación reproducible | candidato aislado junto al host |
| CPython | 3.14.4; python.exe SHA 7ca24f26d6e3f463419ee4f537ddd3acd312c38fe45e678cce08572f26a8bd1a | PSF; LICENSE.txt local, términos históricos preservados | readiness/execution validators; rama 3.14 con EOL programado 2030-10 | patch-auto-candidate; no actualización global |
| CPython runtime DLL | python314.dll SHA a07f7d09c3121492bb066535c6d0811df5fbc2090cbca7031a97bb47ce1480c9 | misma instalación observada | no confundir hash del launcher con integridad de stdlib y extensiones | mismo candidato causal |
| PROJECT-START-READINESS-VALIDATOR | 0.6.1; hash canónico en source lock | AUTHORED / LicenseRef-Workspace-Owner | ocho archivos; 70 tests V283 y V285 | candidato por bug/contrato |
| EXECUTION-VALIDATOR | 1.3.0; hash canónico en source lock | AUTHORED / LicenseRef-Workspace-Owner | estado/checkpoints y assurance; condiciones siguen abiertas | candidato por bug/contrato |
| MARKDOWN-COMPOSITOR | 0.2.1 → 0.2.2; source lock actualizado | AUTHORED / LicenseRef-Workspace-Owner | composición stack-neutral; fallo StrictMode V285 corregido | candidato por bug/contrato |
| Node.js | v24.14.1 detectado por preflight final | OpenJS; identidad de distribución/licencias/SCA pendientes en este record | tooling web de preflight, no stack de producto impuesto | sin upgrade automático |
| pnpm | 11.25.0 detectado por resolución general; algunos paquetes fijan 11.19.0 | pnpm; licencia/grafo/pin aplicable pendientes en este record | available por versión mínima no autoriza adoptar 11.25.0 contra un lock 11.19.0 | resolver versión exacta por perfil |

Licencias observadas por hash: PowerShell LICENSE.txt
296604f12546653ed97726dc5fc2f3945d4cca1f07abd2589f531647177b06fc;
ThirdPartyNotices.txt 5a2d792aa9dae45da7579398ca08c0fbb47d78c563cce492e311b9d5eb9a2142;
Python LICENSE.txt 935cf13e19f8c31b497d20b05d73623431a226b230c3599bc30fa3348979bc68.
No se distribuyen esos runtimes por este cambio ni se certifica su grafo por sus hashes.

## 3. Evaluación del candidato

UPDATE-20260907-COMPOSITOR: trigger bug; FAIL-394; baseline 0.2.1 condicionado.
Cambio: tipar la colección de selección de archivos; no cambios de proveedor,
dependencias externas, privilegios, formato de plan ni admisión. Snapshot previo:
identidad preservada en V283/source lock anterior y V285.

## 4. Gates baseline versus candidate

Baseline: el perfil real y la suite bajo StrictMode fallaron antes de escribir.
Candidate: suite estricta completa PASS, reconstrucción Markdown PASS y tres
entradas reales readiness PASS con BLOCKED esperado, sin overwrites en reentrada.
Licencia propia sin cambio; no SCA nuevo de terceros, no canary ni deploy.

## 5. Decisión

Corrección canónica 0.2.2 REBUILD_VERIFIED / CONDITIONED. No promoción a
REUSABLE_PACK ni cambio de runtimes compartidos. Rollback de mantenimiento:
preservar baseline, revertir sólo delta propio con sus hashes; no volver a usar
el claim StrictMode de 0.2.1, que queda refutado. Revalidar todos los perfiles.

## 6. Checklist y pendientes

- V286 localiza Go 1.26.7, SDK .NET 10.0.400 y psql 18.6 fuera del PATH;
  versiones, rutas y SHA-256 observados en
  `reconstruction_evidence/ISOLATED_TOOLCHAIN_RESOLUTION_V286.md`.
  El runner admite rutas explícitas uniformes y modo Toolchain sin Audit.
  Probe real: 7/8 disponibles, Docker no resuelto; no SCA, pin universal ni
  disponibilidad transitiva aprobados. No reinstalar ni cambiar PATH global.

- Preflight V285 completó su recorrido; BLOCKED con Go, SDK dotnet, Docker y psql
  no detectados mediante resolución por defecto. No demuestra ausencia física:
  V281 conserva Go aislado y otros expedientes conservan toolchains fuera del PATH.
  Resolver cada ejecutable/pin del perfil elegido, sin instalación global automática.
  El host .NET de PowerShell no sustituye el SDK dotnet para compilar SOAP.
- Identidades locales, consumers y política registrados; sin upgrades silenciosos.
- Pendientes: instalación reproducible y SCA del conjunto de runtimes, SBOM completo,
  soporte del OS exacto, inventario transitivo por perfil y rollback de release.
- Advisory/release/security oficiales siempre antes de un parche; un backport local
  se declara ADAPTED_PATCH/AUTHORED_PATCH. Nunca ignorar un hallazgo por no conocer reachability.

Fuentes de soporte consultadas 2026-09-07:
[Microsoft](https://learn.microsoft.com/en-us/powershell/scripting/install/powershell-support-lifecycle?view=powershell-7.6),
[PSF](https://devguide.python.org/versions/). Vigencia de rama no demuestra que cada
patch esté libre de vulnerabilidades; no se afirmó que Python 3.14.4 sea el último patch.

## V313-XTEXT-CANDIDATE — 2026-09-08

- component_id: Go:golang.org/x/text; library, runtime indirect through selected pgx consumers.
- current_version: v0.29.0; candidate_version: v0.39.0, first fixed release named by official GO-2026-5970/CVE-2026-56852, not latest by preference.
- consumers: GO_ENTERPRISE_BACKEND_CORE go.mod/go.sum; GO_OFFICIAL_RETURN_REFUND_WORKER return_refund_worker/go.mod/go.sum. Scan composition67/745 proves both occurrences; other pins of this exact old version not found in executable packs.
- source_registry_or_repository: proxy.golang.org and sum.golang.org; official golang.org/x/text repository go.googlesource.com/text.
- provenance/license: download receipt commit/checksums and BSD-3-Clause LICENSE must be verified before canonical change; no local backport or code patch planned.
- trigger: official OSV2.5.1 exit1 with GO-2026-5970. Advisory https://pkg.go.dev/vuln/GO-2026-5970, JSON https://vuln.go.dev/ID/GO-2026-5970.json, patch https://go.dev/cl/794100 consulted first. Browser redirect of patch failed; canonical Gerrit/source alternative inspected before adoption.
- status: CANDIDATE; promotion blocked. No use of osv fix, dependency bot promotion or floating version.
- candidate_gates: preserve baseline, authenticate module artifact, minimal lock delta, Go tests/vet/build in both consumers, upstream norm regression and package tests, PostgreSQL and connected journey regression, fresh OSV with explicit Python coverage, canonical reconstruction and checkpoint.
- rollback_artifact: staging V313 baseline/ retains four original lock files; canonical source snapshot retained before edits. Rollback to vulnerable version is evidence only, not production authorization.
- open_limits: source SCA does not prove deployed runtime, host exporters, DAST/SAST/call reachability, persistent monitoring or human acceptance; ARCA remains deferred.

- V313 candidate transitive requirement: x/text v0.39.0 go.mod fija Go>=1.25.0 y golang.org/x/sync v0.21.0. Toolchain1.26.7 cumple. x/sync0.17.0→0.21.0 debe declararse si MVS lo selecciona; preservar receipts/licencia/SCA/tests de ambos. No incluir x/tools/x/mod sólo por pertenecer al módulo upstream si no aparecen en el grafo del consumer.

### V313 successor — full selected module graph

Candidate0.39.0 direct scan passed, but full go list -m all selected x/mod0.37.0,
affected by GO-2026-6179/CVE-2026-56865 and GO-2026-6180/CVE-2026-56864.
Official Go advisories fix x/mod at0.40.0 and cmd/go at1.26.6; the actual
verification toolchain1.26.7 is outside that affected range. No cache compromise
is inferred and no destructive cache/go.sum cleanup is authorized or performed.

Keep x/text0.39.0 and x/sync0.21.0, add explicit indirect x/mod0.40.0 security
constraint to both consumer go.mod files. x/text0.40.0/0.41.0 metadata still
selects affected x/mod0.37.0/0.38.0; upgrading by novelty does not solve closure.
x/mod is unused in compiled source/test graphs, but selected by upstream tooling.
This intentional security pin must survive maintenance; plain go mod tidy may
remove an unused constraint and reopen the full-graph scan. Do not suppress it.
Artifacts, commit, BSD-3-Clause and official sumdb package tests required;
candidate1 locks/graphs/logs remain retained. Canonical successor and full SCA
must pass. Functional evidence may transfer only with proven unchanged compiled
graphs and byte-identical trimmed binaries; otherwise rerun affected suites.

### V313 candidate gates completed — local scope

Canonical backend0.4.2/refund0.1.2 fixes the three advisories with official
x/text0.39.0,x/sync0.21.0 and explicit x/mod0.40.0 floor. License/commit/sum/ZIP
receipts, upstream tests, finite norm regression, Go test/vet/build/verify,
PG3,92 phases/4+agenda4 and fuzz PASS. Reconstruction5/5; successor source and
compiled graphs unchanged with2 byte-identical application binaries. OSV2.5.1
final21 sources/465 occurrences/347 unique packages/0known findings; wheel
closure90 pins/47 hashes/134 active requirements verified. Preserve security
floor on future tidy. Rollback artifacts retained, no vulnerable promotion.
Status: CONDITIONED / local dependency gates passed; full runtime/monitoring,
SAST/DAST/load and acceptance still blocked. See V313 evidence, TEST32/EVID23.

## V320 — candidate adapter del checker Node existente

Mismo upstream evaluado enV319, commitc37a56bad56e34fe5223ddd3cb223cc4158136ae
MIT y semver7.8.5ISC fijados; no nuevo matcher, provider o SDK de producto.
Lane de mantenimiento: adaptar únicamente acceso de snapshots a lectura local
verificada y envolver ejecución con identidad/freshness/límites/receipt; candidato
fuera del perfil hasta G0–G8. Source-lock/before originales preservados en staging
elite-v320-513682952df14bd2b01959d4ff04e3d8. Node24.20.0 independiente exacto
deV319; no instalar globalmente ni crear scanner/daemon/scheduler. Reversión:
retirar candidate y conservar V319 CONDITIONED, nunca usar el CLI vacío como gate.

## Resolución V320 — candidato de gate reusable bajo condiciones

Adapter offline30tests/4CLI y rebuild10/10 PASS; sólogetJson difiere del motor
oficial. Bridge SHA requerido después de rojoFAIL519. Semver cambiado de copia
NodeZIP a tarball npm7.8.5 exacto, no de versión:53diferencias de formato preservadas,
licenciaISC y hash de artefacto/53files fijados; tests30+SCA CLI2/0. Condiciones de
instancia Node24.20.0/licencias/snapshots vigentes demostradas; resolverG0–G8
USE_REUSABLE_PACK. No producto/SDKs/nuevo proveedor ni GITHUB/TaskScheduler.
El pack nuevo se materializa aparte; no se actualiza global, pnpm compilado,
CPython u otro runtime ni se elimina el ZIP rechazado. Refresh/renovación abre
otra revisión, nunca cambia sólofecha. Expediente NODE_OFFICIAL_ADVISORY_GATE_V320.md.

## V321 — reapertura pnpm y candidato oficial 11.25.0

SCA del bundle publicado11.19.0:891 archivos verificados,471 identidades npm,
brace-expansion5.0.8 afectado GHSA-rgw5-rvv9-x895. FAIL523 OPEN; alcance de
ejecución desconocido, promoción bloqueada. Advisory oficial fija5.0.9; no ignore
ni parche binario. Evaluar release oficial pnpm11.25.0 del29-agosto mediante
tag/commit fijos, metadata npm, integridad tarball, licencia e inventario compilado
corroborado con lock upstream. La versión nueva no demuestra por sí sola el fix.
Staging elite-v321-184b4d5926e44dc5a3cd1ecbb05ff077, adquisición local aislada.
Candidate gates: SCA exacto2.5.1/all-packages, identidad Node24.20.0 independiente,
install frozen/offline, tests/build y navegador sobre composición nueva; comparar
lock/artefactos y ejercicio de reversión local antes de admitir cambio canónico.
Rollback: conservar baseline/evidencia y retirar candidato; baseline vulnerable
no puede promocionarse. Sin instalación global, nuevos providers ni producción.

V321 verificación de provenance: sigstore4.1.1 ya forma parte del grafo auditado
de pnpm candidato. La copia del ZIP Node requiere45 paquetes y arrastra
brace-expansion5.0.7/ip-address10.2.0; se rechaza ejecutar esa selección completa.
Preparar herramienta aislada con sigstore4.1.1 oficial Apache-2.0, commit
c1dc7d4778a450787fc72b083f2490ad02b714c6, lock transitivo nuevo, scripts desactivados,
SCA antes de ejecución. Sólo verifica attestations públicas del tarball pnpm;
sin claves privadas, firma propia, instalación global ni capability de producto.
Trust roots y resultados deben quedar en receipt; no atribuir verificación
Sigstore a una comparación de campos JSON ni a la firma del tag GitHub.

## V321 focal gate result

OFFICIAL_FIXED_RELEASE pnpm11.25.0:473npm packages0known findings, exact artifact455files; ECDSA/publish/SLSA2bundles and4negative controls; isolated verifier50/0. Three active consumers updated,8files/3unchanged locks,746/746rebuilt;115tests+1existingSKIP pertree/build,92phases/agenda4/runtime4/policy5 and8alternating installs PASS. Historical golden path remains unselected with reuse block. Frozen policy metadata cache fixed without bypass. CONDITIONED: acquisition profile FAIL525, unchanged6native binaries transitive SCA, unsupported uid-number, realtarget monitoring/rollout remain open. Next review before each release and on official advisory; no scheduled automation added.

## V322 — adquisición de mantenimiento comprobada

Core0.4.78 / initialization37files; narrow G0–G8 USE_REUSABLE_PACK con owner global CONDITIONED. Perfil materializado maintenance-runtime-artifacts, lock124 y3sources adquiridos por HTTPS con approval actual del agente bajo autorización previa; no aprobación retrospectiva. Node24.20.0 win-x64, pnpm11.25.0 y Sigstore4.1.1: bytes/licencias exactos y10outputs inmutables tras PRESENT. FAIL525/527/528/529/530/531 cerrados focalmente, TEST41/EVID32. Biblioteca161packs/1453files/757Markdown/52profiles y151preflight steps PASS; preflight global BLOCKED por3tools no resueltos y readiness/target separados. Ver GOVERNED_MAINTENANCE_ACQUISITION_V322.md y PROJECT_OFFICIAL_SOURCE_PROFILE_RECORD.md. Continuar6binarios nativos de pnpm y T2809; no scheduler, ZIP, install global, ejecución de artefactos adquiridos ni producción.

## V323 — cobertura nativa explícita

Reflink0.1.19 commit2223f06f:75Cargo identities exactas/0avisos OSV2.5.1;4addons contienen rustc f6e511ee→1.82.0.8statements npm/SLSA verificados/9negativos;18RustSec std y4avisos oficiales std revisados sin afectar1.82.0. No tarball nativo nuevo, build, ejecución ni SBOM del binario. Fastlist0.3.0 CRT/build no fijados; full license text ausente en ambos trees inspeccionados. FAIL532/LIB-FAIL2248 BLOCKED_EXTERNAL para admisión integral nativa y redistribución; no vulnerabilidad inventada. TEST42/EVID33, PNPM_NATIVE_COVERAGE_V323.md. Continuar T2809 por owners y conservar target/readiness pendientes.

### V325 — diagnóstico de runtime existente

Firefox153.0/rev1538 con Playwright1.62.1:68probes locales PASS no cierran FAIL423 ni autorizan upgrade/patch. Runtime/source hashes y límites en FIREFOX_LIFECYCLE_AND_CONTINUITY_V325.md. Rutas aisladas .NET10.0.400/psql18.6 disponibles para probes; disponibilidad no acredita SCA/licencias/servicio. Node24.20 y pnpm11.25 exactos; FAIL532 nativo sigue abierto.

- **Preflight V325 comprobado:**151steps PASS;7/8tools disponibles por rutas explícitas; Docker sin resolver. Receipt/hash y diagnóstico completo en FIREFOX_LIFECYCLE_AND_CONTINUITY_V325.md. El exit0 informativo no equivale a estado PASS ni a admisión productiva.

## V327 — patch local del tooling, sin upstream nuevo

EXECUTION-VALIDATOR1.3.0→1.3.1 AUTHORED: catch OS/Unicode en los dos CLI y tres tests nuevos; schemas/eventos inalterados.27tests PASS y13comandos/8helps de guía PASS. Fuente anterior en stageV327/pack-before.md y hashes de3bloques en block-delta.txt. Sin dependencias añadidas, runtimes sustituidos ni descarga; rollback restaura el pack anterior y pins como un conjunto, sin reescribir eventos. CONDITIONED y admisión productiva sin cambios.

## V328 — patch local readiness0.6.1→0.6.2

Sin dependencias/runtime/upstream nuevos: dos bloques AUTHORED cambian para
diagnóstico y publicación exclusiva de reportes completos.75tests PASS;8files
reconstruidos y6sin cambios, validate_record idéntico. Requiere hard links del
filesystem, falla cerrado si no soportados. Rollback conjunto pack/pin/lock desde
pack-before.md y before-propagation del stageV328; jamás reemplazar reportes
para recuperar. Fuentes de método y selección empresarial sin cambios.

## V333 — alternativa aislada de pnpm en evaluación

FAIL532 permanece BLOCKED. Se estudia una selección ADAPTED de pnpm11.25.0 a partir de los455archivos exactos ya adquiridos/verificados en V321–322: excluir toda la familia física @reflink y los dos fastlist, conservar sin modificar el resto y exigir packageImportMethod=copy. No nuevo upstream, descarga, instalación global ni cambio de consumers canónicos. Esto no es release oficial ni promoción. Primero inventario, paridad, cerco de carga nativa, instalación frozen/offline sobre fixtures y contrato de fallo; después decidir gates de consumo completos. Rollback: descartar candidato temporal, conservar baseline y receipts. El SCA previo de473identidades es sólo evidencia fechada conservadora; cualquier admisión requiere cobertura y licencias del conjunto seleccionado.

### V333 — candidato para resolver la dependencia nativa

Selección pnpm11.25.0:442/455archivos idénticos,13excluidos con6binarios nativos,22notices preservados.4instalaciones copy y8comparativas hardlink pasan sin intentos nativos;115tests/1SKIP,build e instalación Playwright PASS.473identidades/0avisos en scan nuevo. NOT_ADMITTED: falta selector canónico/routing/gates completos/licencias por scope; original FAIL532 BLOCKED. Evidencia PNPM_NATIVE_SELECTION_V333.md. No nuevos consumers, release ni cierre global;48tests/7pendientes.

## V334 — materialización del candidato existente

Mantenimiento AUTHORED del transporte/proyección, autorizado por continuación del usuario. Selección de pnpm11.25.0 desde bytes ya adquiridos y autenticados V322/V333; sin actualización ni nueva descarga. Owner propuesto PNPM-ARTIFACT-SELECTION-GATE separado del transporte opaco. El selector conserva442archivos/22notices, excluye13, genera receipt y verifica bytes; nunca instala/ejecuta ni admite runtime. Candidato aislado, tests adversariales y reconstrucción antes de incorporar al catálogo. Consumers pnpm y perfil de adquisición inalterados; rollback es retirar sólo tooling candidato preservando receipts.

### V334 — proyección pnpm canónica completada

PNPM-ARTIFACT-SELECTION-GATE0.1.0,4files/17tests dos veces y2proyecciones reales443files idénticas; G0–G8/USE_REUSABLE_PACK sólo tooling en este contexto. Perfil separado1/4 y2pasos de verificación integrados. No consumers/runtime/redistribución admitidos; condiciones explícitas y FAIL532 original conservado. Evidencia reconstruction_evidence/PNPM_SELECTION_GATE_V334.md.

Cierre V334:153pasos PASS,162packs/1457files/772Markdown/53profiles; Docker mantiene disponibilidad BLOCKED.4navegadores conectados PASS y5mediciones Lighthouse reales (performance99mediana;accessibility/best-practices/SEO100).745fuentes intactas y mismo next-env generadoV333.17regresiones/rebuild y proyección443idéntica conservados. FAIL552/553 cerrados; tooling listo dentro de sus condiciones, runtime/redistribución y48tests/7pendientes no se promueven. Siguiente owner: routing canónico y admisión de consumo acotado. Ver PNPM_SELECTION_GATE_V334.md.

### V335 — routing de consumidores pnpm fijados

PNPM-ARTIFACT-SELECTION-GATE0.2.0,6files/43tests fuente+rebuild. Prepare/verify produce recetas offline con configuración aislada, store/cache explícitos y perfiles exactos; nunca ejecuta ni admite runtime.3instalaciones desde rebuild PASS sin descargas/cambios de lock;115tests/1skip ybuild web PASS. FAIL555/557 corregidos sin omitir política; G0–G8 acotado. Perfil separado1/6; producto67/746 inalterado. Evidence PNPM_CONSUMER_ROUTING_V335.md.

Cierre V335:154/154pasos Preflight PASS;162packs/1459files/773Markdown/53profiles; disponibilidad Docker BLOCKED. Pack0.2.0 incorpora routing1/6,43tests/rebuild,3instalaciones offline finales,115tests/1skip ybuild.473declaraciones licenciarias inventariadas, sin admisión total. Checkpoint108 preserva107eventos; runtime/redistribución,48tests/7pendientes y10macrofrentes abiertos.

### V336 — renovación verificable de cache pnpm

Tres checks online desde cache vacía,17metadata oficiales completas con bytes/receipts,3revalidaciones offline sin verdicts históricos y3instalaciones finales offline PASS.275entradas entre3locks;278metadata inventariadas. FAIL561 corregido; no cambio de ejecutables/versiones/política. Procedimiento y límites: reconstruction_evidence/PNPM_CACHE_FRESHNESS_V336.md. Validez acotada a esta ejecución; renovación al reanudar, sin TTL o seguridad continua inferidos. Licencias/runtime/redistribución y10macrofrentes permanecen abiertos.

Cierre V336: VERIFY_LIBRARY PASS162/1459/774/53;3cold policy replays y3instalaciones offline finales sin verdicts históricos;17metadata completas/278mirrors.4textos de licencia/fuente preservados con SHA, sin admisión total. FAIL563 orden de checkpoint corregido; cierre110 conserva109eventos. No código/materialización cambiado; reuse154steps V335.48tests/7pendientes,readiness42 y10macrofrentes abiertos.

V337 candidate evaluation opened: fixed npm-lifecycle1100.1.0 for opaque license inspection only, no version update/install. Registry signature+SHA512 and observed signed release commit retained; current opaque roster expanded with bounded four-digit majors in isolated0.4.79,78tests and33file rebuild. Before/source-lock/profile/rollback evidence in V337; full runtime licence/SCA admission stays open.

### V337 — licencia exacta npm-lifecycle y transporte gobernado

OFFICIAL-UPSTREAM-ACQUISITION-CORE0.4.79:33files/125sources,78opaque checks/rebuild y18perfiles de adquisición. Nuevo perfil separado adquiere sólo npm-lifecycle1100.1.0; firma registry e integridad verificadas. LICENSE del tar exacto coincide con commit/sidecar,7de8archivos con Git blobs; único delta packageManager omitido.29planes de composición actualizados. Fuente licencia verificada, runtime/redistribución no admitidos. Ver reconstruction_evidence/PNPM_LIFECYCLE_LICENSE_V337.md.

V337 closure114: core0.4.79 integrated154step gate PASS, 53composition counts verified; npm-lifecycle1100.1.0 exact licence comparison remains qualified only for observed text/applicability. Additional two fixed-source legacy licence texts retained, not package archive admission. Whole473graph use/notice compatibility and original full-native FAIL532 remain open. Docker absent; no production promotion.

V338 notice coverage:473metadata hashes verified;468identities observed from456bundle identities+21external manifests+root (with overlap).22notice files and21embedded attribution paths inventoried;32exact draft texts preserved. BlueOak five external packages and QRCode vendored MIT require explicit notice treatment; nested works, semver-utils exact text and whole-bundle source/use obligations remain open. All442payload files unchanged; no runtime/redistribution admission. Evidence: reconstruction_evidence/PNPM_NOTICE_COVERAGE_V338.md.

V339:8fixed source manifests and7licence texts verified; four BlueOak source texts,3Noble locked build-input identities outside the original473graph,70node-gyp vendored files identical to fixed tree. Draft39preserves32previous texts. chownr declaration-only, semver-utils exact text and source/bundle/redistribution conditions remain open. All442selected files unchanged. See reconstruction_evidence/PNPM_NESTED_LICENSE_SOURCES_V339.md.

V340:103external Undici6files match fixed source; Undici7licences/headers bound separately. Yarn4.1.7 registry gitHead returns404 but official tag/source version exists at a distinct commit; preserve both, no published-byte equivalence. Root BSD and3MIT source comments retained; draft47preserves39previous texts,442selected files unchanged. FAIL575 scoped provenance remains unresolved. See reconstruction_evidence/PNPM_YARN_UNDICI_NOTICES_V340.md.

V341 research instrument: existing Next16.3.2 compiled Acorn8.14.0, SHA256 758cead0e9764f94320f938ac169fb95c7eeea30dc675a6e5b1133fae239fa19, MIT text preserved; Node24.20.0. No install/update or source acquisition roster change. Parser only reads source; package source is never evaluated. Five AST units/explicit renames do not admit complete Undici7/pnpm runtime or change the473identity inventory.

V342: no upstream/runtime/dependency update. Observed Go1.26.7 executable and io/slog sources are hash-bound in the report; standard-library-only candidate0.3.1 repair. Actual loopback HTTP/local-file tests, offline test/vet/build and finite fuzz gates pass; no new SCA/race or production admission claim.

V343: existing WhatsApp owner0.12.0→0.13.0 is AUTHORED standard-library reporter completion; no upstream, SDK, runtime or external dependency update. Go1.26.7 remains fixed. No new SCA or licence clearance inferred. Final verification interrupted and being resumed.

V343 recovery changes only local GOCACHE after unusable export data; no module/toolchain update or global cache deletion. Exact Go1.26.7 SHA,go.mod/go.sum recorded; final full baseline/vet/build and finite fuzz PASS. No new SCA/redistribution admission claim.

V344: no dependency/SDK/runtime update. Go1.26.7 SHA and unchanged nested-module go.mod/go.sum recorded. Three AUTHORED host/doc/test files change; all12other pack files including provider/processor/DB/migrations/locks unchanged. No new SCA or financial/provider promotion.

V345: no dependency/runtime/SDK change. CPython3.14.4 executable hash retained in runtime-receipt.json; three AUTHORED operational runner/doc/test blocks and three profile pins1.1.3 changed. All Go/provider/DB/lock consumer files identical. No new SCA claim; generic service supervisor still requires separate candidate admission.

V346: no dependency, SDK, runtime or product pack update. Two AUTHORED qualification source blocks in reconstruction evidence reuse exact V345 base e33dbd14b0cd09262897c98263de0d572ac0b66673df6bce2f62c3bf9ade316c; CPython3.14.4 and built-in Windows APIs only. No upstream code/source artifact acquired, no new SCA/entitlement claim.

V347: no dependency/runtime/provider/model update or new upstream source acquisition. AUTHORED qualification uses pinned CPython3.14.4, built-in Windows APIs and existing hashed refund fixture. Exact main-image hash/size and read-shared handles do not certify DLL/script dependencies, runtime SCA/license or trust of the profile signer/selector. Existing dependency blockers remain.

V348 uses the same pinned CPython3.14.4 and built-in Windows event/handle APIs for AUTHORED qualification; no dependency/runtime/source artifact update, product admission or license/SCA closure. V2 is an isolated launch-profile protocol, not a product pack version.

V349: dos archivos AUTHORED de calificación agregan publicación/reconciliación local sobre el launcher existente; seis archivos de supervisión reutilizados sin cambios. Sin nuevas dependencias, runtimes, SDKs, providers, licencias upstream ni promoción de producto. Store local no sustituye admisión WORM/DB/replicación ni política de retención.

V350: adapter/tests AUTHORED de calificación agregados fuera de producto; real host/processor reutilizados byte-idénticos. Go1.26.7, CPython3.14.4 y dependencias go.mod/go.sum sin cambios. Sin nuevo upstream, instalación global, runtime/proveedor ni promoción; binario construido es fixture con Store/Provider sintéticos.

V351: sin nuevas dependencias/runtime ni cambios en packs producto. Se reutilizan Go1.26.7, PostgreSQL18.6, CPython3.14.4, PowerShell7.6.5; cuatro fuentes AUTHORED de calificación y binarios con SHA propios. Instancia PostgreSQL efímera propia con datos sintéticos; no cuentas/proveedores/instalación global. Config loopback está acotada por gramática, tamaño y SHA; no secret manager productivo.

V352: ningún runtime/dependencia/upstream nuevo. PowerShell7.6.5 y Python3.14.4 existentes ejecutan tooling exacto del candidato. Guía0.1.1→0.1.2 sólo corrige orden de aceptación; los9owners consumidos concuerdan con archive148. El ZIP interno no incluye runtimes ni constituye firma/SCA/admisión de release.

V353: corrección AUTHORED GO-FX-CORE0.1.0→0.1.1, sólo math/math-big stdlib existente; ningún runtime, dependencia externa o upstream actualizado. Candidate gates locales: reconstrucción2/2,9tests x3, vet/build, fuzz10s PASS. Rollback conserva pack0.1.0 y fuentes en before151/fx-red, pero reabriría FAIL618; no promocionar rollback inseguro ni aceptar automáticamente otro release. CONDITIONED persiste.

V354: dos revisiones AUTHORED0.1.0→0.1.1, sin nueva dependencia/runtime. math y math/big pertenecen a stdlib existente. Gates candidate+rebuild4/4,19tests x3,vet/build y dos2target fuzz gates4workers PASS; primera corrida24workers falla al deadline, conservada. Rollback a fuentes before153 posible para diagnóstico pero reabre FAIL621/622. CONDITIONED y gates del proyecto/release persisten.

V355: GO-REMINDERS-CORE0.1.0→0.2.0 con ruptura de firmas explícita; GO-WAITLIST-CORE0.1.0→0.1.1 sólo metadata de compatibilidad. Sin nuevos runtimes/dependencias. Gates:2/2fuentes,11tests x3,vet/build, caller legado compile-negative y fuzz acotado PASS. Rollback conserva0.1 para diagnóstico, reabre FAIL624 y no se promueve. Ambos siguen fuera del perfil/CONDITIONED.

V356: promotions/payroll0.1.0→0.1.1 AUTHORED, sin dependencias/runtime nuevos. Gates4/4fuentes,18tests x3,vet/build,37semillas/fuzz4workers PASS. Rollback a before157 reabriría FAIL625/626; no promoverlo. math/big es stdlib de pruebas; no ledger externo ni política/reglas incorporadas. CONDITIONED/fuera del perfil persisten.

V357: GO-REFERRALS-CORE0.1.0→0.1.1 AUTHORED, sin API/dependencia/runtime nuevos.2fuentes,10tests x3,vet/build/modelo fuzz PASS. Rollback before159 reabre FAIL628, no promoverlo. CONDITIONED/fuera del perfil; compatibilidad histórica no integra CRM/loyalty.

V358: reviews0.1.0→0.1.1 y waitlist0.1.1→0.1.2 AUTHORED. Sin firmas/dependencias/runtime nuevos;4fuentes,19tests x3,vet/build/modelos fuzz PASS. Rollback before161 reabre FAIL629/630. FAIL631 de probes reparado antes de calificar. CONDITIONED/fuera de perfil, sin adapter de publicación/envío nuevo.

V359: onboarding/warranty0.1.0→0.1.1 AUTHORED, sin API/dependencias/runtime nuevos.4fuentes,17tests x3,vet/build y modelos fuzz PASS. Rollback before163 reabre FAIL632/633 y no se promueve. CONDITIONED/fuera del perfil; sin activar franquicia ni efectos de garantía/devolución.

V360: surveys0.1.0→0.1.1 y marketing0.1.0→0.2.0 AUTHORED; sin nuevas dependencias/runtime. Marketing requiere tenant en3firmas, retira compatibilidad no integrada; no shim global.4fuentes,20tests x3,vet/build,legacy compile-negative y fuzz PASS. Rollback before165 reabre FAIL635/636. No dispatch ni consentimiento inferido; CONDITIONED/fuera de perfil.

V361: POS0.1.0→0.2.0 rompe Total/ChangeDue a (int64,error); SLO0.1.0→0.1.1 rechaza NaN.4fuentes,20tests x3,vet/build,2legacy negativas y fuzz PASS. Sin dependencia/runtime nuevos; Docker no instalado/adoptado, sólo probe y aclaración del alcance respecto a V351 nativo. Rollback before167 reabre FAIL638/639. CONDITIONED/fuera del perfil, sin cobro/alerta/servicio productivos.

V362: dashboards/helpcenter0.1.0→0.1.1 AUTHORED; firmas conservadas, nuevo ErrVersionExhausted en extremo int64. Sólo stdlib/runtime existente;4fuentes,19tests x3,vet/build/modelos fuzz. Rollback before169 reabre FAIL641/642/643. No integración KPI/search adoptada ni efecto externo, runtime/SDK nuevo o release final. CONDITIONED/fuera del perfil.

V363: i18n/SEO0.1.0→0.1.1,social0.1.0→0.2.0 AUTHORED. Social exige tenant en3firmas; legacy rechazada sin shim.6fuentes,32tests x3,vet/build/fuzz; observability0.3.1 idéntico28tests/vet/build aislados. Ningún runtime/dependencia/motor CLDR nuevo; docs oficiales consultadas sin copiar código. Rollback before171 reabreFAIL646/647/648. No SCA integrada ni release final/servicio admitidos.

V364: reconstruction_evidence/YARN_EMBEDDED_ASSET_PROVENANCE_V364.md fija un asset Yarn ESM completo (72463bytes decodificados) y5funciones AST contra source fijado;14controles negativos y extracción independiente. Fuente de paquetes no ejecutada.163Markdown de packs,442payload y47notices intactos. FAIL575/publicación completa,FAIL532/nativos y obligaciones de licencia siguen;43/48sin promoción.

V365:21cores/42fuentes reconstruidos,222tests ordinarios/23targets con306semillas,vet/build PASS;95,4%statements observados no equivalen a readiness. Se corrigen6compatibilidades ausentes/desfasadas mediante metadata sucesora,12fuentes byte-idénticas.20CONDITIONED/1CANDIDATE,fuera de53planes;43/48intacto,FAIL385 abierto. Evidencia reconstruction_evidence/CORE_COMPATIBILITY_AUDIT_V365.md.

V367: investigación oficial de entrenamiento prioriza TRL/SFT y PEFT opcional como DISCOVERED; torchtune/torchforge fuera de primera línea por mantenimiento, torchtitan diferido por scope/runtime. Identidades/artefactos sólo observados en metadata, no adquiridos/admitidos. RESEARCH_INCOMPLETE/FAIL663; TEST09 sigue bloqueado. Evidencia reconstruction_evidence/HISTORY_TRAINING_SOURCE_RESEARCH_V367.md y training_gap_v367/record.json. No pedir corpus privado para esta preparación de biblioteca.

V367 distribución181: los registros training_gap_v367 se conservan íntegros, con SHA-256, como secciones del reporte HISTORY_TRAINING_SOURCE_RESEARCH_V367.md; materializarlos sólo en raíz aislada de investigación. No son nuevos archivos sueltos del release ni un pack admitido.

V368: identidad Git TRL1.11.0/1.12.0 contrastada:580entradas por árbol no truncado,579idénticas y sólo VERSION cambia; commit1.12.0 exacto59c4a8e104413fa9f4ca1a54eaf2ff93c0f299be. No equivale a bytes de paquetes ni admisión/runtime. Spec0.1.2 precisa cierre funcional con capacidades comprometidas ejecutadas y cero bloqueos contradictorios; no cambia48oráculos/estados. Evidencia reconstruction_evidence/TRAINING_SOURCE_IDENTITY_AND_CLOSURE_V368.md.

V369: core de adquisición0.4.80,36files,127sources y19perfiles internos. Perfil aislado adquirió2sdists TRL1.12.0/PEFT0.20.0 y conserva7outputs inmutables.559archivos inspeccionados:544Git-equal,15metadata de packaging revisados,licencias raíz iguales.98checks transporte; no runtime/training admission. FAIL663 sigue investigación pendiente del grafo y pipeline. Evidencia reconstruction_evidence/TRAINING_SOURCE_ACQUISITION_V369.md.

V370: candidato59distribuciones/103relaciones,58baseTRL+PEFTopcional.59METADATA hash-verified/equivalentes; resolver pip confirma59y2negativos sin wheels/install/framework execution. OSV59versiones0hallazgos declarados;30provenance subjects coinciden,firmas no verificadas. Licencias de artefactos/nativos,adquisición gobernada y runtime/gates pendientes. DISCOVERED/FAIL663 RESEARCH_INCOMPLETE. Ver reconstruction_evidence/TRAINING_DEPENDENCY_GRAPH_V370.md.

V371 FAIL675: safetensors0.8.0 and tokenizers0.23.2 acquisition does not authorize use. Embedded crate query returned6records (aliases/severity pending triage); quarantine retained. Exact wheels/SBOM/query hashes in V371. Prefer official fixed artifacts; any source rebuild must declare its own provenance, dependency delta, tests and license/SCA evidence.

V371: core0.4.81/39files;55quarantine checks+98opaque PASS.59wheels/214055762bytes adquiridos;23854files/23795RECORD hashes inspeccionados.226native files/3SBOMs;371identidades consultadas,6records=4avisos distintos (3security+1maintenance). FAIL675 mantiene cuarentena. Safetensors sourceb7c0f38 corrige pyo3/memmap2,45registry crates0OSV sólo metadata; build pendiente.43/48sin cambio. Ver reconstruction_evidence/TRAINING_WHEEL_QUARANTINE_V371.md.

V371 closure193: Rust/MSVC located outside PATH (inventory, not new installation or build). Narrow official safetensors b3d8d72 lock rejected unchanged:67registry crates retain memmap2/anyhow advisories. Current b7c0f38 lock45/0 remains unreleased-source research candidate with extra features, not admitted. Tokenizers current exact tree has no Cargo.lock and still declares paste; no fixed resolved graph inferred. See V371 portable receipts.

V372 TEST05: isolated official OCB0.160.0 builder source acquired at Collector cd3455c. Evaluate only a minimal component selection (OTLP receiver, bounded file exporter, Prometheus exporter, memory limiter and health check), not the full rejected distributions. Go1.26.8 executable rehashed against prior admitted toolchain receipt; no toolchain update. Official Go proxy and sum.golang.org resolution retained in isolated cache, source go.mod/go.sum exact and no automatic source test/build until graph query reviewed. Prometheus3.14.0 source acquired for rule evaluation; no external target or paid service.

V372 OCB build candidate: official Go advisory GO-2026-5024/CL770080 fixes x/sys NewNTUnicodeString overflow in0.44.0. Preserve upstream0.41.0 go.mod/go.sum; select exactly0.44.0 as declared local dependency-only ADAPTED_PATCH candidate. Source Go files unchanged; build/tests/SCA must pass before use. No scanner ignores or upstream release relabeling.

V372 final qualification: governed two-source acquisition PASS, selected Go component ZIP origins fixed; OCB32queries/0findings, runtime387queries retain3package-absent advisory records. All four module graphs verify; Collector rebuild after preserving upstream test output is byte-identical. Actual runtime uses Go1.26.8 and preserved dependency-only corrections. Review/rebuild/rollback uses isolated owned directories and prior source snapshots, never replaces a running project. See V372 licenses, build info, security binding and tests; integral TEST-03 remains blocked.

V373 official safetensors0.9.0rc0 candidate observed Sep10, Git tag v0.9.0-rc.0 commit77a5d1d5234918bc9da7637ad178dac57f0e2834. Prefer inspection of exact published artifact over unreleased-source rebuild; not yet runtime admitted. RustSec paste advisory is INFO/unmaintained, with explicit upstream tokenizers maintainer revert of pastey replacement. Determine exact compile/runtime scope before any nonblocking disposition; do not call it fixed or suppress it.

V373 isolated candidate qualification:58exact wheels (base TRL lane; optional PEFT excluded), safetensors0.9.0rc0 replaces quarantined0.8.0 only in the new candidate.423PyPI/declared vendored-SBOM identities return only RUSTSEC-2024-0436 INFO/unmaintained paste1.0.15, retained and under explicit compile-only/source review, never called fixed. All wheel/RECORD bytes verified; no private data/model/networked job. Install only into fresh owned candidate venv for finite synthetic runtime tests; current project/global Python unchanged. Native C/C++ component completeness is not claimed from these SBOMs; must qualify actual exercised components and preserve limitations before reusable admission.

V373: opt-in HISTORY-MODEL-TRAINING-PIPELINE0.1.0/17files, HISTORY_MODEL_TRAINING_PACK_PLAN2packs/21files;26policy tests and actual earlier synthetic SFT/evaluation/rollback. Final exact reconstruction/bootstrap/E2E and TEST09 decision are recorded in reconstruction_evidence/HISTORY_MODEL_TRAINING_CONTROL_V373.md. NEW requires no history; EXISTING keeps data/model private with its own authority/quality/runtime gates. No automatic training, provider call or deployment. SDK component selection precedes AUTHORED glue; see HISTORY_TRAINING_SDK_QUALIFICATION_V373.md and HISTORY_TRAINING_SDK_GAP_V373.md. Integral franchise remains67/754.

V373 cierre201: TEST09 PASS por pack17files/composición21,27policy, bootstrap58wheels/23713file hashes, SFT real, evaluación independiente,13negativos, timeout/no replay, decisiones exactas y rollback local con baseline previo. Preflight200160pasos ejecutados PASS/55planes;164/1506/818/55, integral67/754.45/48, tres pendientes TEST02/03/07; no cambia oráculos ni hereda corpus/modelo/seguridad/deploy del consumer. Ver HISTORY_MODEL_TRAINING_CONTROL_V373.md.

V374: auditoría por claim21núcleos,228tests/306semillas/23campañas (39811771ejecuciones),6comparaciones nuevas de fuentes,42archivos reconstruidos,20packs/40archivos compuestos y candidato rechazado.20CONDITIONED/1CANDIDATE intactos; cero equivalencias empresariales/promociones inventadas.13negativos de enlace de evidencia integrados en VERIFY_LIBRARY. Preflight/cierre del control aún pendientes;45/48sin cambio. Ver reconstruction_evidence/CORE_CLAIM_ADMISSION_V374.md; la integración de negocio más amplia de FAIL385 mantiene sus owners.

V374 cierre203: TEST02 PASS como auditoría por claim21núcleos;228tests/306semillas,23targets finales (6389808ejecuciones),20packs/40files compuestos, candidato rechazado,13negativos de evidencia y Preflight202160pasos PASS.46/48controles;20CONDITIONED/1CANDIDATE y FAIL385/integración más amplia conservados.164/1506/819/55, integral67/754. No cambia48oráculos ni promueve producto. Ver CORE_CLAIM_ADMISSION_V374.md; continuar TEST03 y luego TEST07 según gates materiales.

V374 rectificación204 / FAIL710: el cierre203 de TEST02 se retracta. Fuente/tests/composición son válidos, pero FAIL385 aún exige equivalencia/integración material.45/48PASS; TEST02/03/07BLOCKED. Se preservan historia y mejoras; no se cambia el oracle ni se reduce el roadmap. Ver CORE_CLAIM_ADMISSION_V374.md, encabezado vigente.

V374 successor205: SECURE-OPS1.1.4 corrects the error-rate denominator that hid low-traffic outages. Canonical Prometheus3.14.0 rules9scenarios/13assertions and5wrapper rejections PASS; exact pin/no-engine distinction retained. Three consumer plans now select1.1.4.164packs/1507files/819Markdown/55plans;1260AUTHORED/140ADAPTED/107VERBATIM; integral67/755. No dependency acquisition, core promotion or whole-control closure;45/48 unchanged. Evidence: reconstruction_evidence/CORE_CLAIM_ADMISSION_V374.md. Full successor Preflight pending.

V374 checkpoint207: fullPreflight206 passed162executed steps and all55compositions, including exact Prometheus engine regression. SECUREOPS1.1.4 correction complete; integral67/755, library164/1507/819/55. Current45/48unchanged; TEST02/03/07 remain blocked. Retain the203retraction and actual FAIL385 integration work. See reconstruction_evidence/CORE_CLAIM_ADMISSION_V374.md.

V375: HTTP_SLI_ALERT_INTEGRATION_V375.md proves real order HTTP/PostgreSQL → official closed metrics → unchanged canonical alert → repair/replay, and a byte-identical full-profile rebuild. GO-HTTP-METRICS-REFERENCE0.1.0/8files, opt-in68packs/763files; ordinary67/755 unchanged. 380official tests and2authored tests PASS;62public versions/0OSV findings in this reference graph,30compiled upstream license inventories. Trusted synthetic reference only; TEST02/03/07 and FAIL385 remain open,45/48.

V375 checkpoint212: fullPreflight212164executed steps PASS,56profiles composed;165/1515/823/56 and1268AUTHORED/140ADAPTED/107VERBATIM. HTTP reference508requests,four actual failures,canonical ten-minute alert and recovery,one durable order/idempotency/outbox,exact rebuilt executable;380official+2local tests PASS. Only Docker unavailable in that receipt.45/48 unchanged; TEST02/03/07 remain blocked. Next action: FAIL732 official Next/Sharp patched-release qualification, preserving TEST02 owner integration work.

V376: Next16.3.4/Sharp0.35.4 replaces the affected Next16.3.2/Sharp0.35.3 baseline in TS-GO-API-WEB-BRIDGE0.5.5. Five exact public artifacts, ten verified provenance attestations, 8589 installed files identical to inspected archives,152 declared npm versions/0OSV findings versus3baseline advisories,115web tests plus1explicit integration skip, typecheck/build,92connected browser phases and4agenda projects with53exact migrations PASS. Multitab/session gates enabled on all four browser projects. Native AVIF/PNG/WebP smoke confirms libheif1.23.2. This is scoped advisory remediation; full vendored/native security and redistribution obligations remain blocked. Acquisition core0.4.84/46files,193IDs/22profiles and all four transport suites PASS;30dependent source plans updated. HTTP reference0.1.1 records the changed web-parent plan; Go/SQL/runtime evidence is reused only with byte parity.45/48 unchanged. See reconstruction_evidence/NEXT_SHARP_SECURITY_REPAIR_V376.md. Final canonical composition, HTTP binary parity and full Preflight remain pending in this entry.

V376 closure: Preflight214 passed all164executed steps and all56compositions; only Docker unavailable. Canonical165packs/1521files/824Markdown/56plans,1269AUTHORED/145ADAPTED/107VERBATIM; ordinary67/755 and HTTP68/763. Exact Next/Sharp package-advisory repair, signed artifacts and92connected phases/4agendas complete. Current45/48; TEST02/03/07 and native/vendor license/security gates remain blocked. Final canonical consumer has754parent source files byte-identical to the connected candidate plus the two proven Next-generated type imports. HTTP executable matches V375 byte-for-byte, so its unchanged15minute runtime evidence is reused without repetition. Next work returns to TEST02 owner integration, including the public-web SEO/i18n boundary; do not repeat completed browser creation/recovery work or the fixed dependency qualification.

V377 integrates actual Next canonical/robots/sitemap/WebSite JSON-LD under the existing BFF owner0.5.6/43files:5new AUTHORED files and3page changes; no dependency update. Public indexing is opt-in, limited to the real home/model-list pages and public render assets; private pages default noindex.137web tests/1preserved explicit integration skip, strict typecheck/build and3actual HTTP configurations PASS. The same build proves disabled/enabled/changed origin,Host poisoning isolation,mandatory cache revalidation,private noindex and JSON-LD HTML escaping/CSP nonce. Connected qualification: three browser projects passed initially, Firefox failed teardown with known FAIL423/742; a separate exact Firefox23phase rerun plus4agenda projects passed. These results cover4projects/92distinct phases across two runs, not a flawless92phase initial run or a Firefox lifecycle fix. Known upstream incident stays open.45/48 unchanged. See reconstruction_evidence/PUBLIC_WEB_METADATA_INTEGRATION_V377.md; final canonical composition and Preflight pending.

V377 cierre217: integración pública canonical/robots/sitemap/JSON-LD en BFF0.5.6/43archivos comprobada con137tests,3configuraciones HTTP reales y composición final67/760. Regresión conectada:69fases de tres perfiles PASS en el primer run y23Firefox+4agendas PASS en una repetición aislada; primer fallo de teardown conservado y FAIL423/742 abierto. Preflight216:164pasos ejecutados PASS,56planes compuestos;165packs/1526files/825Markdown/56planes,1274AUTHORED/145ADAPTED/107VERBATIM. Docker sigue ausente. Binario HTTP idéntico a V375; fuentes Go/SQL/regla/dependencias sin cambios. Los48criterios/status son idénticos:45/48, TEST02/03/07 bloqueados. Esta frontera SEO de referencia está integrada; faltan i18n/localización y otras equivalencias funcionales, seguridad/admisión nativa y release integral. No contar el replay Firefox como reparación upstream ni repetir creación/recuperación ya probadas.

V378: no dependency update; exact8589installed artifact files and frozen lock match V376. Node/Intl runtime identity recorded. Native/vendor/security/redistribution admissions unchanged.

V379: no dependency update; offline frozen install and8589exact artifact files match V376, unchanged lock. Help is AUTHORED read-only integration; Go stdlib TLS proxy is test-only. All existing native/vendor/security/redistribution admissions remain conditioned.

V380: no dependency changes. Frozen offline install and8589exact authenticated Next/Sharp files match V376; unchanged lock. Three new AUTHORED files under the portal owner. No new runtime, provider, data ingestion or training.

V381: frozen offline installs;8589Next/Sharp files match V376, identical pnpm lock. Two new AUTHORED files, no dependency/runtime/provider changes. Browser manifest CRLF normalization preserves exact parsed JSON values.

V382: no runtime/dependency/provider changes. Frozen offline installs and8589exact Next/Sharp file hashes preserve V376 artifacts and unchanged pnpm lock. Five new AUTHORED helper/test/document files; existing quote formatter moved to shared owner with reversible source parity. Native/vendor licence/security and original full-native pnpm conditions remain open.

V383: no dependency/runtime/provider updates.8589authenticated Next/Sharp files and pnpm lock unchanged; five authored recovery component/boundary/test files. No native/vendor licence/security admission inferred.

V383 continuation FAIL735/532: exact Next MIT companion captured and byte-matched to the verified repository and three signed packages. Sharp native29notice rows/28versions reconciled; libnsgif owner identified as13vendored blobs at verified libvips426af3f44246fce9cfa8dd51a353aa4dfd48c553, four document/blob matches. No fabricated NetSurf version, binary correspondence, distribution or full native/SCA admission. Raw receipts in V383 notice-analysis; production blocks retained.

V384 demuestra correspondencia exacta del insumo nativo Sharp: libvips-42.dll18614784bytes, versions.json y avisos coinciden con @img/sharp-libvips-win32-x641.3.3 firmado. Perfil previo/2attestations/tamper, cuarentena sin ejecución y4suites del transporte PASS; core0.4.85/48files/194IDs/23perfiles. DLL C++ y build MXE/notices/SCA completos pendientes. Runtime y perfil68/795 intactos;45/48 sin promoción.

FAIL760 corrected from actual source profile selections; original unpublished preparation retained. Runtime parent unchanged. Preserve moving-notice input and C++/MXE/SCA/source obligations; no substituted or executed native runtime.

V385: ZIP MXE8.18.6 adquirido con perfil exacto;488archivos inspeccionados, DLL18614784bytes/versiones iguales al input Sharp firmado y licencia igual a fuente fijada.28/28versiones y28SHA256 de componentes vinculados a recetas. Core0.4.86/50files/195IDs/24perfiles,4suites PASS. Runtime68/795 y69/803 conservado; TEST02/03/07 siguen bloqueados,45/48.

Evidence reconstruction_evidence/MXE_NATIVE_ORIGIN_V385.md. FAIL761 keeps historical base/build and full nested native source/licenses/SCA open. Scoped libvips source LICENSE text is not a combined payload license.

V385 cierre233: origen ZIP-DLL-version-license exacto,488files inspeccionados y28versiones/SHA256 de fuentes fijados; core0.4.86/50files/195IDs/24perfiles,4suites PASS. Preflight232164pasos/56composiciones PASS, Docker ausente. Runtime68/795 y69/803 idéntico aV384.45/48. FAIL763: libxml22.15.4 corrige seguridad frente al2.15.3deSharp; fix oficial fijado, sin release Sharp/MXE posterior observado ni reparación binaria demostrada.

FAIL763 official fix c94eb021→96498992,39commits/47files captured through verified GNOME mirror. No repaired binary/runtime admission. Next work follows fixed source graph and controlled candidate rebuild/containment. FAIL762 original404 and exact checksum correction retained.

V386: usuario informa restriccion de contenido de ciberseguridad con acceso Daybreak y pide registrar/bypassear en orden de trabajo, no el control. Investigacion libxml2 suspendida ACCESS_BLOCKED; TEST03/07 siguen BLOCKED. Ambos make check completos ASan/UBSan Linux PASS; ninguna DLL Windows reparada o promovida. Siguiente TEST02: integracion funcional pendiente de FAIL385, sin repetir auditoria21yahecha.45/48.

FAIL764/765/766 REGRESSION_PROVEN for acquisition/tooling/module-suite scope only. FAIL763 remains ACCESS_BLOCKED; user reported platform restriction and directed proceeding to independent TEST02. Evidence reconstruction_evidence/NATIVE_QUALIFICATION_ACCESS_BLOCK_V386.md.

V387: paginacion real en pedidos/casos cliente y unidades fabrica;3listas/81registros por navegador,4navegadores PASS,0escrituras/snapshot6tablas intacto.272tests/1skip heredado,typecheck/build;6packs reconstruidos y805files exactos,franquicia68/797. TEST02 sigue pendiente integral;TEST03libxml2 ACCESS_BLOCKED por indicacion del usuario,TEST07dependiente.45/48.
FAIL767/768/769 REGRESSION_PROVEN for the bounded pagination/environment/fixture claims in reconstruction_evidence/PRIVATE_PORTAL_PAGINATION_V387.md. Native security scope remains ACCESS_BLOCKED; no test status promotion.

V388: seis listados cliente/fabrica/admin paginados con APIs existentes;163registros de lista por navegador,4navegadores PASS,0escrituras/snapshot6tablas intacto.276tests/1skip heredado,typecheck/build;6packs reconstruidos,805files exactos,68/797integral.45/48;TEST02integral pendiente,TEST03libxml2 ACCESS_BLOCKED,TEST07dependiente.
FAIL770 REGRESSION_PROVEN in reconstruction_evidence/ADMIN_PORTAL_PAGINATION_V388.md; FAIL771 full Preflight rerun pending. V386 native research remains ACCESS_BLOCKED; no promotion.

V388 cierre239: seis listados cliente/fabrica/admin conectados al cursor real;163registros de lista por navegador,4navegadores PASS,0escrituras/snapshot6tablas idéntico.276tests/1skip previo,typecheck/build;6packs y805files reconstruidos,HTTPbinario idéntico. Preflight238164pasos/56perfiles PASS; Docker ausente.165packs/1561files/836Markdown; integral68/797.45/48,TEST02integral pendiente,TEST03 ACCESS_BLOCKED por restriccion reportada/usuario,TEST07dependiente.
FAIL771 REGRESSION_PROVEN; original236 failure preserved, full238 verifies canonical inventory and both pagination repairs.

V389 checkpoint241: horas de turnos cliente coherentes entre cuenta/gestion con locale/zona del negocio, labels SSR preservados al hidratar.282tests/1skip previo,typecheck/build;2configuraciones x4navegadores,2instantes con limite de dia/DST,snapshot7tablas intacto y0escrituras.6packs reconstruidos/805files,HTTPbinario idéntico;68/797.45/48; TEST02integral pendiente, V386native diferido,TEST07dependiente. Preflight241 pendiente.
FAIL772/773/774 REGRESSION_PROVEN in reconstruction_evidence/CUSTOMER_APPOINTMENT_TIMEZONE_V389.md. No native/security or whole control promotion.

V389 cierre242: turnos cliente muestran fecha/hora/zona del negocio coherentes en cuenta/gestion/recarga, sin depender del dispositivo.282tests/1skip previo,8runs navegador(2config x4),0errores en segmento de fechas,0escrituras/snapshot7tablas idéntico.6packs/805files reconstruidos,8589artefactos yHTTPbinario idénticos. Preflight241164pasos/56composiciones PASS; Docker ausente.165packs/1561files/837Markdown;68/797.45/48; TEST02integral pendiente, V386native diferido yTEST07dependiente.
FAIL772–774 repaired; FAIL775 qualified fixture sequencing only, not an upstream browser fix. Temporary error classifier rejected and absent from canonical code.

V390 checkpoint244: cancelacion cliente conectada y recuperable: receipt ligado aID/org/estado/version,guard sincrono,GETexplicito tras resultado incierto.294tests/1skip previo,typecheck/build;4navegadores,16cancelaciones/16audits/16outbox con28POSTincluidos rechazos;8regresiones lectura/fechas sin escrituras.6packs/805files reconstruidos,8589artefactos yHTTPbinario idénticos;68/797.45/48;V386diferido,TEST02integral yTEST07pendientes. Preflight244 pendiente.
FAIL776/777/778 REGRESSION_PROVEN in reconstruction_evidence/CUSTOMER_CANCELLATION_RECOVERY_V390.md. No native/security or whole control promotion.

V390 cierre245: recorrido cliente cuenta→turnos→cancelacion→recuperacion GET cerrado en referencia.294tests/1skip previo,4runs cancelacion con16cambios/audits/outbox unicos y28POSTincluidos rechazos;8runs lectura/fechas siguen sin escrituras.6packs/805files reconstruidos,8589artefactos yHTTPbinario idénticos. Preflight244164pasos/56composiciones PASS; Docker ausente.165packs/1561files/838Markdown;68/797.45/48; TEST02integral pendiente, V386native diferido yTEST07dependiente.
FAIL776/777/778 proven with final canonical source and unchanged full verifier. No production/security or21core-equivalence claim.

V393 candidate lane: official latest metadata observed2026-09-11 stillSharp0.35.4/libvips1.3.3/MXE8.18.6; no fixed release inferred. ExactNext16.3.4 package declaresSharp optional. Current805source reference has no next/image, image imports or transformation call in application code. Isolated candidate will omit only this optional dependency via pinnedpnpm11.25.0 and explicitly disable unused runtime image optimization; no component-version upgrade or known finding suppressed. Promotion requires unchanged domain journeys, build/typecheck/tests, authenticated remaining artifacts and disabled image endpoint with Sharp absent. Source and prior runtime retained; security/fullrelease remain blocked pending evidence.

V393: perfil BFF0.5.16 omite Sharp opcional y desactiva optimización runtime; lock152→122sin versiones nuevas,122stanzas y12571archivos restantes exactos,8534artefactos firmados exactos;294tests/1skip previo,92fases/4agendas/4lecturas/4cancelaciones PASS,OSV122/0.4packs reconstruidos/805files,68/797; HTTPbinario idéntico. No Sharp fix ni reapertura V386; native/tooling restantes yTEST02/03/07 abiertos,45/48. Ver reconstruction_evidence/OPTIONAL_IMAGE_DEPENDENCY_CONTAINMENT_V393.md.

V393 cierre251: cuatro packs incorporados/805files reconstruidos; lock122versiones/0avisos OSV,12571archivos y8534firmados exactos en candidato e instalación limpia.294tests/1skip heredado,92fases/4agendas/4lecturas/4cancelaciones PASS. Preflight250:164pasos PASS/56perfiles; disponibilidad global BLOCKED sólo Docker. Inventario165/1561/841/56,integral68/797.45/48; TEST02integral,TEST03tooling/native restante yTEST07siguen pendientes. Sharp no se instala en este perfil; V386investigación diferida.

V394: PNPM-ARTIFACT-SELECTION-GATE0.3.0/6files corrige FAIL782 (pins empresariales y Playwright obsoletos), entrega BLUEOAK-NOTICE.md ligado a5declaraciones exactas y verifica3recetas actuales;51tests y10mutaciones reales rechazadas. Payload442/22notices intactos; sin ejecución ni admisión runtime/redistribución. Evidencia: reconstruction_evidence/PNPM_CURRENT_ROUTING_NOTICES_V394.md. Los pendientes BlueOak se reducen a entrega downstream/revisión restante; ausencia de LICENSE original chownr no se falsea. Otros permisos/source/publicación de pnpm continúan abiertos.

V394 cierre253: FAIL782 reparado con selector0.3.0/6files;51tests PASS,3recetas actuales,5declaraciones BlueOak exactas con aviso local y10mutaciones reales rechazadas. Preflight252:164pasos PASS/56perfiles, nuevo control de compatibilidad de3consumers; disponibilidad BLOCKED sólo Docker.165/1561/842/56;45/48sin promoción. Payload442 y22notices intactos; pnpm no ejecutado, runtime/redistribución pendientes. Ver reconstruction_evidence/PNPM_CURRENT_ROUTING_NOTICES_V394.md.

V395: PNPM-ARTIFACT-SELECTION-GATE0.4.0/6files entrega QRCODE-NOTICE.md con encabezado de autor/modificación exacto y texto MIT completo, además del aviso BlueOak intacto.59tests PASS,3recetas reales y12negativos rechazados;10source blobs/10regiones de bundle fijados por separado, sin claim de equivalencia completa. FAIL783 entrega local resuelto, FAIL784 metadata ADAPTED reconstruida;payload442/22notices intactos. No ejecución ni admisión runtime/redistribución;45/48. Ver reconstruction_evidence/QRCODE_VENDOR_NOTICE_DELIVERY_V395.md.

V395 cierre257: aviso QRCode con copyright/permiso completos entregado y ligado a3recetas;59tests/12negativos PASS,6fuentes reconstruidas;BlueOak ypayload442 intactos. Preflight256164pasos/56perfiles PASS;Docker ausente.165/1561/843/56;procedencia1306/148/107. FAIL783/784/785 reparados;45/48sin promoción. semver-utils fuente gitHead en GitHub404, investigación pendiente sin bypass ni licencia inventada. Ver reconstruction_evidence/QRCODE_VENDOR_NOTICE_DELIVERY_V395.md.

V396 candidate inspection only: C:\Users\NL\AppData\Local\Temp\elite-v396-e5a797060309431d8e37863e7f5eeea3\rebuilt2\elite_sources\source-profiles\maintenance-semver-license-quarantine.json; approved scope hash c45d05378960de4be7424266e4a3b539ac60fb273f851d811e0f5dd2e4368ab8. semver-utils1.1.4 exact4193bytes/SHA512; no installed dependency update or production admission. Registry key expired; original license inspection pending. Rollback is retaining isolated evidence, never overwriting original payload. See PROJECT_OFFICIAL_SOURCE_PROFILE_RECORD.md.

V396: cerrado descubrimiento y entrega local de licencia original semver-utils1.1.4: MIT OR Apache-2.0 explícito, opción MIT completa y1839bytes originales retenidos. Core0.4.88/51files/196sources/25profiles añade sólo cuarentena npm exact4193bytes/SHA512;140checks y gates anteriores PASS. Planner0.5.0/6files/67tests;3recetas reales/12negativos PASS, BlueOak/QRCode y442payloadfiles intactos. Clave registry vencida documentada; no firma vigente, equivalencia build, ejecución de recetas ni admisión global pnpm.45/48sin promoción; V386/Daybreak sigue diferido. Ver reconstruction_evidence/SEMVER_ORIGINAL_LICENSE_DELIVERY_V396.md.

V396 cierre262: licencia original semver-utils1.1.4 localizada y entregada bajo opción MIT completa;3recetas/12negativos/67tests PASS. Core0.4.88/51files/196sources/25profiles,140checks; planner0.5.0/6files. Preflight261164pasos/56perfiles PASS;Docker ausente.165/1562/844/56;procedencia1307/148/107.45/48sin promoción, Daybreak/V386 diferido. Fuente/licencia semver resuelta en alcance local; firma registry vencida y admisión pnpm restante condicionadas. Ver reconstruction_evidence/SEMVER_ORIGINAL_LICENSE_DELIVERY_V396.md.

V397: PNPM-ARTIFACT-SELECTION-GATE0.6.0/7files entrega el conjunto completo de47textos retenidos/151253bytes en3recetas:141copias verificadas independientemente,12negativos reales y82tests PASS.22notices originales comparados contra payload y25evidencias conservan alcance; BlueOak/QRCode/semver y442payloadfiles intactos. Catálogo ADAPTED con términos por texto, sin algoritmos ni promoción de licencias/source/relinking/runtime.45/48; Daybreak diferido. Ver reconstruction_evidence/PNPM_RETAINED_NOTICE_DELIVERY_V397.md.

V397 cierre264: conjunto47textos/151253bytes entregado en3recetas,141copias exactas/12negativos/82tests PASS. Preflight263164pasos/56perfiles PASS;165/1563/845/56 y1307/149/107. Docker ausente,45/48sin promoción. README/roadmap vigentes sincronizados y cortes previos preservados (FAIL794). Licencias/source/relinking/publicación restante y Daybreak diferido siguen explícitos. Ver reconstruction_evidence/PNPM_RETAINED_NOTICE_DELIVERY_V397.md.

V398: PNPM-ARTIFACT-SELECTION-GATE0.7.0/8files entrega fuente original next-path1.0.0, manifiesto y MPL completa en3recetas:9copias exactas/12negativos reales/90tests PASS. Commit oficial y3Git blobs verificados;4sentencias comparadas bajo adaptadores explícitos,10mutaciones rechazadas. No equivalencia runtime ni reproducibilidad pnpm. Colección47, suplementos anteriores y442payloadfiles intactos.45/48; Daybreak diferido. Ver reconstruction_evidence/NEXT_PATH_MPL_SOURCE_DELIVERY_V398.md.

V398 cierre266: next-path fuente/manifiesto/MPL entregados,9copias exactas/12negativos/90tests/8files PASS. Preflight265164pasos/56perfiles PASS;165/1564/846/56 y1307/150/107. Docker ausente,45/48 sin promoción. Entradas vigentes sincronizadas con historia conservada. Daybreak diferido; integración funcional,source/relinking/publicación restante yrelease pendientes. Ver reconstruction_evidence/NEXT_PATH_MPL_SOURCE_DELIVERY_V398.md.

V400: encuestas conectadas con PostgreSQL/OIDC/Next, recuperación GET,4navegadores/8respuestas/8POST y retención CLI1/1/0.23archivos nuevos AUTHORED, packs GO-CUSTOMER-SURVEY-API0.1.0 y TS-CUSTOMER-SURVEY-PORTAL0.1.0, CONDITIONED.45/48sin promoción; TEST02/03/07 siguen abiertos. Ver reconstruction_evidence/CONNECTED_CUSTOMER_SURVEYS_V400.md.

V400 cierre270: encuestas conectadas y retención transaccional incorporadas. Go/TSsurvey0.1.1, app1.9.4, HTTPref0.1.15;23archivos nuevos,828/828 reconstruidos.4navegadores/8respuestas/8POST,24concurrentes,restart,CLI1/1/0,fuzz6219;correcciónSQL2→0tablas parciales PASS. Preflight269164pasos/56perfiles PASS más deltaSQL focal;167/1587/850 y1330/150/107,franquicia70/820,HTTP71/828,web7/134.45/48sin promoción; Daybreak diferido. Ver reconstruction_evidence/CONNECTED_CUSTOMER_SURVEYS_V400.md.


V401 / FAIL804 — 2026-09-11: adm-zip0.6.0 bundled in exact pnpm11.25.0 is affected by GHSA-vwc7-r8mq-g2x9 / CVE-2026-76845. Pinned OSV2.5.1 identifies one affected package among476 identified components; this is not a complete recursive native SBOM. The historical V321473/0 observation remains historical and cannot authorize current promotion. pnpm11.26.0 also contains adm-zip0.6.0; rejected as a fix without executing it. pnpm12.4.1 is an isolated inspection candidate with native optional executables, not a clean22-package runtime. Current pnpm promotion and new use are frozen; prior payload retained. Advisory lists no patched release. Reachability UNKNOWN pending call-site examination; no exploit or destructive probe executed. Response lane NO_UPSTREAM_FIX pending official patch review. Daybreak/V386 remains deferred. Evidence: Temp/elite-v401-20260911-closure/nested-scan-result.json and candidate-acquisition-receipts.json. Authority: https://github.com/advisories/GHSA-vwc7-r8mq-g2x9 . Reopen on demonstrated removal/containment or admissible fixed replacement; no risk acceptance.

V401: contención ZIP terminada en candidato pnpm aislado;16módulos adm-zip retirados,441archivos preservados/1bundle cambiado,4rechazos sinIO y3instalaciones offline reales PASS.475identidades conocidas/0hallazgos; SBOMrecursivo no cerrado;12571archivos node_modules idénticos. Pack0.8.0/9files,167/1588; pnpm general BLOCKED, TEST02/03/07 siguen abiertos,45/48; Daybreak diferido. Ver reconstruction_evidence/PNPM_ZIP_CONTAINMENT_V401.md.

V401 FAIL807: the three explicit frozen offline installs succeeded without downloads, but the later pnpm exec Next probe omitted the matching store/offline settings. pnpm automatically recreated node_modules and downloaded71packages before reporting Next16.3.4. This invalidates zero-network for that dispatch, not the initial install receipts. Prior logs preserved. Fix the diagnostic invocation with explicit offline/store/cache/no-script/no-hook settings; no global runtime promotion.

V401 reconciliación final: comparación12571/12571idéntica completada antes de la solicitud de detenerla; no se detuvo proceso. FAIL807 separa las3instalaciones iniciales offline de un exec Next que descargó71paquetes por configuración omitida. Exec corregido con store/offline explícitos PASS sin descargas. Historial/log anterior retenido; no prueba global de ausencia de red. reconstruction_evidence/PNPM_ZIP_CONTAINMENT_V401.md

V402 communications: no external dependency version update. Reused exact Go OAuth/OIDC, Meta SDK/Python, Node/Next and PG graph; added materializable license notices and runtime/adapter preflight. Canonical normalization preserves Python AST/JavaScript output with receipts. Current global SCA and pnpm runtime admission remain pending; prior scans are not promoted to this revision.
