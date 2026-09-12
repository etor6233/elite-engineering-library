# Project Authority Freshness Record — mantenimiento

## Corrección sucesora V319 — Node y cobertura de distribución

Before: V317 no cubría Node; interpretar npm/node y reputación/firma de un ZIP
como seguridad integral sería incorrecto. After: Node24.14.1 presenta22CVEs
oficiales; OSV alias0 es cobertura inválida. Node24.20.0 core0 no elimina9
advisories del npm que incluye su ZIP. ZIP completo rechazado; selección mínima
node.exe+LICENSE probada con paquete pnpm separado. Checker oficial con corpus
vacío/deadline falla negativos y permaneceCONDITIONED. Firma gpgv, identidad de
procesos y --all-packages corregidos con evidencia before/after preservada.
Fuentes fechadas/snapshots exactos, criterio y límites en NODE_RUNTIME_SECURITY_V319.md.
Triggers: advisory nuevo, cambio de archivo/componente, pre-gate/pre-release,
corpus/schedule vencidos o resolución del adapter portable y su operación.

## Corrección sucesora V317 — identidad Go

- previous_statement: gates recientes atribuidos a Go1.26.7 desde una ruta recordada.
- contradictory_evidence: current-toolchain de V281 es1.26.8; versión, compiler,
  VERSION, 15036 archivos idénticos al ZIP oficial y cuatro binarios V313 lo prueban.
- corrected_statement: official-toolchain es1.26.7, 15034 archivos idénticos al ZIP;
  V313 compiló1.26.8. Para las24 entradas afectadas en conjunto, conservar errata
  y retirar atribución exacta sin receipt individual; no inventar versión sustituta.
- blast_radius: 24 EVID/informes, gates Go, metadata de builds y SCA del runtime.
- files_or_decisions_changed: 24 erratas con before/after, contrato de assurance,
  perfiles de freshness/monitoring y registros. No promoción automática de1.26.8.
- gates_rerun: composición actual67/745 con1.26.7 comprobado; detalle y resultado
  efectivo en reconstruction_evidence/TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md.
- related_failure_ids: FAIL-20260908-506; instalación507 y sincronización508.
- trigger: antes de cada gate/build, cambio de ejecutable, raíz, versión o advisory.

Fuentes oficiales actuales preservadas: https://go.dev/dl/?mode=json&include=all,
https://go.dev/doc/security/vuln/database y https://vuln.go.dev/index/db.json.
No se reescriben logs. V281 ya distinguía ambos árboles; no hay evidencia de una
mutación externa reciente. Una revalidación del head no prueba cada snapshot viejo.


Derivado de markdown_system/PROJECT_AUTHORITY_FRESHNESS_TEMPLATE.md.
Aplicar markdown_system/AUTHORITY_FRESHNESS_AND_SELF_CORRECTION_CONTRACT.md.

## 1. Control

```yaml
record_version: "1.0"
project_id: elite-library-maintenance
owner: agente_de_mantenimiento
updated_at: "2026-09-07"
review_due_items: [runtime_advisories, official_release_artifacts, full_profile_claims]
stale_blockers: []
unresolved_conflicts: []
last_corpus_wide_consistency_gate: "V284 VERIFY_LIBRARY_PASS; no equivale a revisión semántica de todo el corpus"
```

## 2. Authorities

V305, 2026-09-07: se reutilizan AWS/React consultados en V303 para recuperación
de respuestas inciertas. Se consultan PostgreSQL18 Explicit Locking y Transaction
Isolation en documentación oficial para lock transaccional cooperativo y snapshot
Read Committed por sentencia. La implementación AUTHORED aplica la clave de
advisory lock ya usada por los triggers, antes de la sentencia de reserva/cambio;
no atribuirla a PostgreSQL ni a AWS. LEAD_COMMAND_RECOVERY_V304.md conserva su
historia; AVAILABILITY_CANCELLATION_RECOVERY_V305.md documenta el delta y los gates.
No nueva dependencia ni renovación global, SCA o prueba de carga.

V304, 2026-09-07: se reutilizan las consultas oficiales AWS/React de V303 para
el mismo claim estrecho de respuesta incierta y fence de UI. No se renueva la
vigencia de todo el corpus ni se adquiere código upstream nuevo. LeadActions y
los métodos auditados son AUTHORED; actor deriva del principal HTTP existente,
CAS y outbox usan la transacción del owner PostgreSQL. Sin nuevas dependencias,
migraciones ni permisos del proveedor. Ver LEAD_COMMAND_RECOVERY_V304.md para
los gates locales y sus límites; no representa aprobación de esas empresas.

| ID / claim estrecho | Fuente / identidad | Applicability y estado | Próximo trigger |
|---|---|---|---|
| MS-PWSH-SUPPORT | Microsoft Support Lifecycle, consultado 2026-09-07; 7.6.5 LTS / EOL rama 2028-11-14 | CURRENT_CONDITIONED: versión local coincide; grafo/OS no certificados | cambio de host, release o advisory |
| PSF-SUPPORT | Python Developer Guide versions, consultado 2026-09-07; rama 3.14 bugfix / EOL previsto 2030-10 | CURRENT_CONDITIONED: sólo rama de soporte; 3.14.4 observado, no afirmado último | antes de release/update |
| MS-FILES | FileAttributes / Select-Object oficiales; V284 | CURRENT_VERIFIED para semántica ReparsePoint y ausencia de -Reverse; no certificación del bridge | cambio de API o fallo |
| GOOGLE-LAUNCH | Google SRE Reliable Product Launches, reconsultado 2026-09-07 | CURRENT_CONDITIONED: método proporcional, no sello de producción | cambio de riesgo/alcance |
| NASA-VALIDATION | NASA Appendix, consultado V283 | CURRENT_CONDITIONED: verificación separada de validación e integración | decisión de aceptación |
| LOCAL-CONTRACTS | AGENT_SYSTEM_START, TOTAL_SYSTEM_CAPABILITY_CONTRACT, source lock y checkpoint | CURRENT_CONDITIONED: owners locales por hash; pendientes explícitos | cambio de hash/contrato |

## 3. Eventos de contradicción/autocorrección

V295, 2026-09-07: PROJECT_SECRETS_TEMPLATE distinguió entorno de custodia,
constructores de loaders y providers de devolución/cobro. Microsoft Export-Clixml
confirma DPAPI sólo mismo usuario/equipo Windows; se usa primitiva oficial, helper
AUTHORED, no bóveda productiva. OpenAI Projects documenta spend_alerts y
spend_limit duro con enforcement: no repetir que todos sus controles son sólo
alertas ni prometer disponibilidad sin cuenta. Google Ads SDK confirma variables
exactas y precedencia de GOOGLE_ADS_CONFIGURATION_FILE_PATH. Meta leads, TikTok
subscription y ML auth no recuperados íntegramente: mantener ACCESS_REQUIRED,
no inventar menús/permisos/renovación. URLs consultadas y límites en la guía y
reconstruction_evidence/CREDENTIAL_PREPARATION_V295.md; no cambio de SDK/runtime.

DISCOVERY-V288: al investigar uso de chats históricos, las fuentes oficiales
OpenAI de model-optimization y deprecations señalan Evals read-only 2026-10-31,
retiro API/dashboard 2026-11-30 y restricciones de acceso a nuevo fine-tuning.
No seleccionar esos servicios como requisito nuevo por ejemplos antiguos; método
de evaluación y endpoint comercial son claims distintos. Presidio ya está fijado
con owner Data Privacy Stack (origen Microsoft); el sitio Microsoft redirige al
mantenedor actual y no justifica reatribuirlo. Meta history docs no recuperadas:
FAIL-405, sin afirmar scopes/ventana de historia. Fuentes exactas y alcance en
reconstruction_evidence/CHAT_HISTORY_LEARNING_AUDIT_V288.md. No se actualizó un SDK
ni se afirma haber reauditorado todo el corpus histórico de referencias.

CORRECTION-V283: antes, readiness admitía CLI pero exigía otra interfaz; después,
0.6.1 conserva autorización/evidencia y acepta CLI probado. FAIL-389; 70 regresiones.

CORRECTION-V284: antes, el bridge afirmaba instalación y recuperación seguras sin
probar junctions, orden de marcadores y fallo de escritura; después, controles y
23 checks específicos. FAIL-390..392. No defensa garantizada ante escritor hostil concurrente.

CORRECTION-V285: antes, se suponía que el compositor 0.2.1 funcionaba desde un
caller estricto; prueba real y suite lo refutaron antes del primer write. Después,
0.2.2 tipa files como array y repite la suite con StrictMode. FAIL-394; source lock,
pack, tests y evidencia actualizados. Ningún claim se atribuye a Microsoft por usar PowerShell.

## 4. Consistency matrix

V291: OpenTelemetry Metrics API/error handling, Go spec integer overflow y Google
SRE monitoring consultados 2026-09-07. Claim limitado: no-negatividad del contador,
telemetría que no interrumpe negocio y overflow Go no automático. Corrección
AUTHORED TryInc, no implementación OTel ni nuevo código Google. Se corrige el
claim excesivo del pack sin readmitir privacidad; fuente canónica 0.1.1 y tests.

V289: consulta Microsoft Learn medallion (originales preservados y derivados),
NASA Appendix y Google SRE Release Engineering, 2026-09-07. Uso estrecho en
MAINTENANCE_ASSURANCE_RAW_BOUNDARY_V289.md y requisito raw del usuario. No nueva
dependencia ni autoridad de importación Meta; FAIL-405 sigue pendiente. El
contrato AUTHORED no es certificación de esas empresas ni actualización de todo
el corpus; EVID-01 conserva evidencia V287 sin convertirla en release actual.

NEW/EXISTING: spec LIB-R01/R08 → compositor/bridge/validators → V284/V285;
pruebas locales parciales, no readiness integral. Documentos, leads y franquicia:
LIB-R03..07 → T2802..T2809, no cerrados por excluirlos del primer slice de tooling.

## 5. Checklist

V301, 2026-09-07: Google Cloud Monitoring Distribution y OTel data-model/API
contrastados para histograma. Inclusividad de Google y OTel difiere: el código
local conserva superior inclusivo, no adapter Google. Before0.2.0 ocultaba
overflow/permitía NaN y contador negativo; after0.3.0 corrige con regresiones,
sin nueva dependencia ni atribución upstream. CANDIDATE permanece; referencias
y evidencia en HISTOGRAM_BOUNDARY_REPAIR_V301.md. No renovación global del corpus.

V302, 2026-09-07: React conditional-rendering consultado para JSX, no autoridad
de autorización ni aprobación de código local. Permisos cotejados con handlers
HTTP canónicos; React/Next/Playwright/Go conservan versiones y locks fijados.
Cambios AUTHORED, sin upstream nuevo ni renovación global. Evidencia y límites
en OPERATOR_SECTIONS_RECOVERY_V302.md; FAIL457 continúa abierto.

V303, 2026-09-07: AWS Builders' Library / making-retries-safe-with-idempotent-APIs
y React / referencing-values-with-refs consultados para incertidumbre, ausencia
de reenvío ciego y referencias síncronas entre renders. Son autoridades de método,
no origen ni aval de DeliveryExceptionResolution AUTHORED. El contrato real de
ResolveDeliveryException conserva versión/CAS/transacción/outbox y sus tres
decisiones existentes; no se modifica por la documentación. Sin upstream nuevo,
dependencias nuevas ni renovación global. Evidencia DELIVERY_EXCEPTION_RECOVERY_V303.md.

V300, 2026-09-07: PostgreSQL18 transaction-iso consultado para snapshot estable,
sin adquisición o actualización de dependencias. Before: V299 protegía comandos
pero las lecturas podían devolver vínculos incoherentes; el portal decía que no
se registró una acción incluso después de commit. After: proyección fail-closed,
snapshot común y recuperación GET con incertidumbre explícita. Fuente, hashes y
pruebas en DELIVERY_READ_RECOVERY_V300.md. Cambios AUTHORED, no atribuidos al
upstream; no reescribir reportes históricos ni afirmar validación productiva.

V299, 2026-09-07: se consultó PostgreSQL18 explicit-locking y Microsoft Business
Central release-reopen-documents para el delta de entrega. No se actualizan
dependencias ni se incorpora código desde esas páginas. Before: gap map V296
describía pago/entrega sólo lectura. After: mapa V299 enlaza solicitud local V298
y defensa relacional V299, sin afirmar cobro o creación inicial. Límites y URLs
en DELIVERY_RELATIONAL_SCOPE_V299.md; los reportes históricos no se reescriben.

Before/after, blast radius, fallos y referencias preservados. Pendiente reauditoría
semántica del perfil completo y de sus APIs/advisories al adoptarlo. Fuentes móviles
sirven para consulta de método/soporte, nunca como identidad de bytes incorporados.

Fuentes: [Microsoft soporte](https://learn.microsoft.com/en-us/powershell/scripting/install/powershell-support-lifecycle?view=powershell-7.6),
[PSF](https://devguide.python.org/versions/),
[Google SRE](https://sre.google/sre-book/reliable-product-launches/),
[NASA](https://www.nasa.gov/reference/system-engineering-handbook-appendix/).

## Autocorrección V308 de trazabilidad

Contrato antes: SHA256 `a59946cf4db7a2a43dd9aa3b891eaed84928ba2b772c78035f94e51dc80717b6`.
La versión anterior queda en checkpoint42 y snapshot fuera de distribución.
V306–V308 tenían evidencia pero faltaban sus refs propias en assurance;
TEST25–27/EVID16–18 las enlazan, con límites. Plan PASS no sustituye el
BLOCKED de evidencia integral. Sin nueva autoridad externa, dependencia ni
renovación global de fuentes; los métodos estrechos V303/V305 siguen aplicando.

## V312-CONTINUITY-CORRECTION — 2026-09-08

- previous_statement: borrador local de próxima acción ubicaba histograma pendiente después de V292.
- contradictory_evidence: GO_OBSERVABILITY_CORE0.3.0 y HISTOGRAM_BOUNDARY_REPAIR_V301.md ya prueban el cierre focal; mapa roadmap V301 lo registra. Autoridad es fuente local y evidencia reproducida, no memoria resumida.
- corrected_statement: logger/counter/histograma reparados; candidato fuera del perfil, faltan montaje host, exportación/alertas/retención/carga/race y gates integrados.
- blast_radius: borrador checkpoint.py aún no ejecutado; no cambió producto ni se repitieron pruebas de histograma.
- files_or_decisions_changed: siguiente acción de V312 pasa a auditar host/supervisor/worker seleccionados; no se añade el candidato. Borrador before conservado en staging V312 checkpoint-before-freshness.py.
- gates_rerun: checkpoint y resume V312 al materializar el cursor corregido; no runtime gate nuevo por cambio de planificación.
- related_failure_ids: FAIL473, pérdida de contexto; no nuevo defecto del histograma.
- status: CURRENT_VERIFIED para este hecho local; no admisión de observabilidad integral.

## V313-XTEXT — 2026-09-08

- previous_statement: v0.29.0 remained pinned in two selected Go consumers; V285 monitoring record had no current composition scan.
- contradictory_evidence: OSV2.5.1 exact engine reports GO-2026-5970 twice. Official Go advisory published2026-07-14, modified2026-08-10, affected beforev0.39.0.
- corrected_statement: x/text v0.29.0 admission reopened; v0.39.0 is an isolated candidate pending all gates. No claim of exploitable input reachability without separate evidence.
- blast_radius: backend and refund worker module locks; no migration or product feature.
- files_or_decisions_changed: dependency update record and failure494; canonical locks remain unchanged at candidate opening.
- gates_rerun: current source scan, official advisory retrieval; candidate tests/SCA/reconstruction pending.
- related_failure_ids: FAIL494.
- status: CURRENT_CONDITIONED; production remains blocked.

## V313 final successor — 2026-09-08

Initial minimum x/text candidate superseded for full selected-graph admission:
x/mod0.37.0 remained affected even though not compiled. Final floor0.40.0 is
fixed by both official advisories; Go verifier1.26.7 exceeds corrected1.26.6.
Before/after locks, graphs, patches and OSV JSON preserved. Canonical backend
0.4.2/refund0.1.2; scan347 unique packages/0findings,5/5 rebuild,2 binary hashes
identical across unused floor; functional gates recorded in V313. Claims remain
conditioned to identified sources, not runtime/deploy/complete library readiness.
related_failure_ids: FAIL494,FAIL495,FAIL500. status: CURRENT_CONDITIONED.

## Sucesor V320 — guards del checker y precisión de bytes

V319 permanece como before del CLI raw; V320 demuestra adapter offline que
preserva el motor oficial y cierra vacío/deadline con hash/freshness/límites.
La condición se resuelve para la instancia exacta mediante G0–G8, no mediante
reputación. semver7.8.5 del Node ZIP difiere de npm por terminadores/whitespace;
se elige el tarball exacto y se conserva before/after, sin normalizar una prueba
de identidad. Expiración≤7días en lock vinculado por SHA al runner, sin auto-refresh.
El binding es integridad, no una firma digital; la firma PGP de Node pertenece
a la adquisición separada V319. Before de esta corrección editorial conservado
en V320/freshness-before-binding-wording.md; ninguna firma fue emitida por el gate.
Se corrige texto de assurance heredado: observabilidad actual0.3 e histograma
V301 ya reparado; cores21 conciliadosV293, con integración/operación target aún
pendientes. Ni las evidencias antiguas ni sus logs se reescriben. Triggers:
archivo/advisory/snapshot/runtime distinto o vencido, antes de reutilización/release.

## V321 — corrected scope of pnpm identity and profile record

Before:V319 verified891pnpm11.19.0 published bytes, no bundled dependency SCA. After:471npm identities/1advisory; official11.25.0 verified455files/473identities/0findings, source-lock472transitives, ECDSA+publish/SLSA exactcommit. Reports/locks keep before and do not transfer rawarchive identity to SCA. Before official-source-profile V285 remains in staging; new discrepancy FAIL525 blocks ACQUIRED until canonical profile is materialized/executed. No signature claim from SHAalone and no fabricated historical approval.

## V322 — adquisición de mantenimiento comprobada

Core0.4.78 / initialization37files; narrow G0–G8 USE_REUSABLE_PACK con owner global CONDITIONED. Perfil materializado maintenance-runtime-artifacts, lock124 y3sources adquiridos por HTTPS con approval actual del agente bajo autorización previa; no aprobación retrospectiva. Node24.20.0 win-x64, pnpm11.25.0 y Sigstore4.1.1: bytes/licencias exactos y10outputs inmutables tras PRESENT. FAIL525/527/528/529/530/531 cerrados focalmente, TEST41/EVID32. Biblioteca161packs/1453files/757Markdown/52profiles y151preflight steps PASS; preflight global BLOCKED por3tools no resueltos y readiness/target separados. Ver GOVERNED_MAINTENANCE_ACQUISITION_V322.md y PROJECT_OFFICIAL_SOURCE_PROFILE_RECORD.md. Continuar6binarios nativos de pnpm y T2809; no scheduler, ZIP, install global, ejecución de artefactos adquiridos ni producción.

## V323 — cobertura nativa explícita

Reflink0.1.19 commit2223f06f:75Cargo identities exactas/0avisos OSV2.5.1;4addons contienen rustc f6e511ee→1.82.0.8statements npm/SLSA verificados/9negativos;18RustSec std y4avisos oficiales std revisados sin afectar1.82.0. No tarball nativo nuevo, build, ejecución ni SBOM del binario. Fastlist0.3.0 CRT/build no fijados; full license text ausente en ambos trees inspeccionados. FAIL532/LIB-FAIL2248 BLOCKED_EXTERNAL para admisión integral nativa y redistribución; no vulnerabilidad inventada. TEST42/EVID33, PNPM_NATIVE_COVERAGE_V323.md. Continuar T2809 por owners y conservar target/readiness pendientes.

## V327 — diagnóstico y guía contrastados con ejecución local

El CLI1.3.0 exponía traceback/exit1 ante paths ausentes y UTF8 inválido; el patch1.3.1 devuelve BLOCKED/exit2. Pack anterior y red10fallos/4errores conservados;27tests verdes posteriores. README ahora explicita Python de tooling, enlaza la guía y refleja inventario V326161/1453/52 y source registry V322124/17; no se revisan hechos externos por memoria ni se reescriben reportes históricos. Evidencia CLI_ENTRY_DIAGNOSTICS_V327.md.

## V328 — evidencia local corrige diagnóstico y publicación de readiness

0.6.1 exponía traceback ante root ausente y replace podía perder un reporte
tardío.0.6.2 conserva funciones semánticas/catálogo y reemplaza esa publicación
por hard link exclusivo;75tests PASS, historial y resultados previos preservados.
No revisión de autoridad externa ni cambio de respuestas:42bloqueos reales
idénticos. Ver READINESS_REPORT_PUBLICATION_V328.md.

## V329 — ajuste de métodos a generación de authority map

Antes: el registro describe métodos de contexto y flujo; no había expediente
específico del generador. Después: tres revisiones documentales contrastadas con
docs/releases/security actuales, rechazo acotado por falta de compatibilidad.
Spec Kit9118ed1 y AI-DLCe49341d no cambian de pin; Conductor99ba10e sólo investigado.
No se adjudican las cinco fases de docs actuales a AI-DLC1.0.1.30d de vigencia
para la observación; nota hasheada no representa source bytes ni admisión.
Before conservado en stage V329/before; criterio sucesor bootstrap §5.1 y
AUTHORITY_MAP_GAP_V329.md. Sin fuente móvil promovida ni historia reescrita.

## V333 — documentación11.x redirige a12.x

Antes: se intentó consultar pnpm.io/11.x/settings/node-modules para confirmar defaults. Después: redirige a la página12.x, cuyo orden auto por plataforma difiere del bundle11.25.0 fijado. No usar ese contenido como autoridad histórica; comportamiento de esta observación se verifica en bytes fijados y probes. Sin upgrade ni inferencia de cambio de licencia sobre pins. El URL GitHub supuesto no dio fuente; no se adquirió ni promovió. Reporte PNPM_NATIVE_SELECTION_V333.md conserva la corrección y límites.

## V334 — alcance de licencias por ruta observado

Antes: la revisión de selección no había distinguido la declaración del cliente pnpr de la licencia del servidor padre. Después: package.json de pnpr/client3.0.0 en v11.25.0 declara MIT; pnpr/LICENSE.md contiene PolyForm Shield1.0.0. Son evidencias distintas: no se deduce licencia del bundle desde un enum ni se impone automáticamente la del servidor al cliente. LICENSE raíz del repositorio muestra2026, el artefacto conservado2025; no se sustituyen bytes del artefacto. Observaciones y URLs en stageV334/license-scope-observations.md SHA256 581a1302618c79ac695c410dbf6f3da51d175163f29b19a2ad36ac4c2d531d0a. Investigación documental sin nuevo código adquirido, grant legal inventado ni admisión de runtime/redistribución.

## V335 — declaraciones de licencia y owner redirect

Antes: cuatro familias de licencia pendientes se mencionaban agregadas. Después: censo473metadata sin etiquetas ausentes y cuatro paquetes exactos, paths conservados en bundle. next-path/openpgp/SPDX con revisiones de metadata observadas; SPDX redirige kemitchell→jslicense al mismo commit. npm-lifecycle sin gitHead: no se reemplaza por branch actual ni se declara licencia total verificada. Inventario/observaciones hasheados en PNPM_CONSUMER_ROUTING_V335.md. No fuente/lock histórico reescrito, nuevo runtime adquirido ni permiso de redistribución inferido.

## V336 — cache y attestations

previous_statement:3offline installs V335 prueban ruta funcional, pero cache/frescura seguían pendientes. contradictory_evidence: cache reutiliza políticas sin expiry independiente; primer cold replay V336 falla por full metadata ausente pese a online PASS. corrected_statement:17full JSON oficiales más renovación sin verdicts copiados permiten3checks offline y3instalaciones finales, sólo para este snapshot. blast_radius: procedimiento de mantenimiento pnpm, no producto/versiones/licencias. Canonical owners: pack/plan/roadmap/dependency record y PNPM_CACHE_FRESHNESS_V336.md. gates_rerun:3online,3cold offline,3final installs,auditor y verifier estructural; reuse43tests y154steps anteriores con fecha. related_failure_ids:FAIL560/561/562. No garantía de SCA/no-downgrade ni frescura indefinida.

V336 source-license observation successor:4full text bodies/receipts now retained in stage license-texts (three package-linked fixed refs and official GNU GPL3). Previously only URLs/declarations were inventoried. This improves preservation, not whole-bundle admission. npm-lifecycle1100.1.0 exact tarball/SHA512 known from V321metadata; gitHead absence does not justify moving-branch substitution and does not preclude governed inspection of its exact package licence. Source links/hashes in PNPM_CACHE_FRESHNESS_V336.md.

## V337 — npm-lifecycle sin gitHead, licencia resuelta por evidencia

previous_statement: V335/V336 no tenían commit/texto aplicable al paquete; no sustituir rama actual. new_evidence: commit de release da8e0e802ab2e50860fe237c705ca6fb3092e910 observado signed/verified; npm firma+SHA512; tar gobernado10780bytes y LICENSE idéntica. corrected_statement: licencia concreta verificada por comparación; registry gitHead sigue ausente y packageManager es delta declarado.7archivos son Git-blob iguales. blast_radius: source lock/perfil/transporte0.4.79/29planes ynoticeADAPTED, sin cambio de consumidores. Gates78/18/9/rebuild/acquisition, Preflight posterior; related_failures564-568. No admisión del conjunto473 ni alcance de redistribución inferido.

V337 FAIL570 profile-count self-correction: six explicit verifier expectations missed+2 core-file delta; ten intros were already below V336verified counts. Before/verified/currentexpected mapping retained in profile-count-correction.json. Corrected owners are VERIFY_LIBRARY.ps1 and ten composition-plan introductions; no selection, runtime or archive altered. All53expected counts will be compared to integrated output, not accepted merely because assertions were edited.

V337 closure114: all53actual composition counts now independently match expected; Preflight113154steps PASS. Fixed-source individual3.0.0 LICENCE and qrcode-terminal0.12.0 LICENSE retained with Gitblob/SHA256 receipts, without normalizing legacy metadata. semver-utils declared-source discovery limitation is not proof of source absence; its exact licence remains unresolved. See V337 closure and legacy-source-license-receipts.json.

## V338 — selected notices and stale summary correction

473registry declarations remain historical hash-verified metadata, not a complete notice inventory. New fixed-artifact evidence separates468observed identities,22notice files,21embedded paths and nested works. QRCode file-level MIT supplements root Apache declaration; five BlueOak external packages lack direct text and exact licence URL in selected payload. Full official BlueOak text and source-header Gitblob proof retained; no automatic legal clearance. README previous161/1453/52 and email142 corrected to verified162/1461/53 and150; before snapshot retained. Documentary inventory775→776 from one report only. FAIL571/572 recorded; structural gate follows115. Semver-utils discovery errors preserve uncertainty, not evidence of absent licence.

V338 post-gate closure correction: Preflight113 PASS predates duplicate canonical ledger references added during closure114. Hash-valid resume is not a structural rerun. Successor formatting corrected without removing historical content; exact root uniqueness regex PASS and full structural replay after116. FAIL573/574 preserve before/after and separate162packs from163directoryMarkdown entries.

V338 closure117: structural replay PASS162/1461/776;53profile counts equal V337. Final closure rechecks exact ledger-token uniqueness and execution hashes, preserving failed115log and prior closure formatting evidence. No implementation-pack change or new154step execution claim.

## V339 — source/notice scope refined

V338473registry identities are not a complete recursive embedded SBOM: OpenPGP fixed lock reveals three exact Noble build inputs absent from that list. Preserve original graph and add conditioned overlay rather than infer476runtime components. Source LICENSE for ciphers includes Thomas Pornin beyond the abbreviated parent attribution. All70gyp bytes equal fixed node-gyp tree; embedded packaging declares23.3.dev0 despite external requirement>=24.0, so never substitute the constraint for vendored identity. yallist root scope preamble preserved; chownr fixed tree lacks named licence text. Seven texts added to draft39, no automatic licence normalization/admission. Source receipts and bounded reads govern; semver tool access limitation remains uncertainty.

V339 closure119: structural PASS162/1461/777;53profiles and163implementation-directoryMarkdown hashes unchanged. Next-source map separates external undici6.28 from bundled7.29 and Yarn4.1.7. Source notices remain conditioned; final resume/plan and unique ledger-token checks follow without replaying unchanged executable suites.

## V340 — Yarn registry/tag divergence and versioned Undici proof

Undici6external103files now equal its fixed tree; Undici7source notices remain a distinct transformed-bundle claim. Yarn current npm gitHead5761b03f returns404 from tree/directory/direct-commit; official tag objectbbe8151e resolves to4fe4d4bf with4.1.7manifest, but is not signature-verified. No declared metadata overwritten; alternate tagged source remains conditioned. Node/Blake/Joyent MIT comments preserved under parent BSD and origin README records adapted Node18.9.0. FAIL575 remains scoped; draft47 is evidence, not cleared distribution.

V340 closure121: structural PASS 162/1461/778; 53 profile counts and 163 implementation-directory Markdown hashes unchanged. Independent draft47 and selected payload442 verification passed. Yarn registry/tag divergence remains unresolved; no new full Preflight or production admission claim.

## V341 — scoped transformed-source evidence and honest remaining scope

Undici7 source-marker/notice evidence is supplemented by five matched AST units from two modules. Preserve V340 narrower before-state: this does not retroactively prove all112modules or original borrowed revisions. Eight actual-source mutations reject; exact bundle bytes retained. No global percentage is inferred from7pending of48test rows. See PNPM_SCOPED_SOURCE_CORRESPONDENCE_V341.md.

V341 gate122 rejected a machine-local evidence reference (FAIL577). Preserve failed log/snapshot; use temp-relative notation per existing lesson2081. Exact AST evidence and package bytes are unchanged. Structural replay follows123.

V341 closure124: structural123 PASS162/1461/779 with53profiles and163implementation-directory hashes unchanged. New nonportable reference corrected with failed122 evidence preserved. Five scoped AST matches/8actual-source mutations remain narrow evidence; no complete source/runtime admission.

## V342 — explicit delivery claim corrected

Local0.3.0 returned successful delivery for four short/invalid writer counts despite its explicit failure contract.0.3.1 validates count before slog and keeps ErrDeliveryUnknown/no retry. Source before-state/red log preserved;2/2 rebuilt identity and actual HTTP/file reference pass. No Go vulnerability or newer-toolchain adoption inferred; CANDIDATE and target gaps remain.

V342 closure126: structural125 PASS162/1461/780;53composition counts unchanged. Only observability candidate pack changed, two tested/rebuilt blocks exact. Writer delivery fix is locally regression-proven; actual runtime operational gaps and CANDIDATE admission remain.

V343: prior abstract-only StatusReporter now has an explicit JSON/net.Conn implementation in existing owner0.13.0, locally connected to StatusWorker.Run tests. Actual service mounting and collector/alert/retention claims remain open. Go context/net methods checked; no documentation-displayed toolchain upgrade.

V343 final code evidence: canonical0.13.0 report implementation and post-commit failure/restart are verified. Reporter is no longer abstract-only; deployment/supervisor/collector/retention/alerts remain open. Preserved interrupted logs do not count as PASS. Toolchain SHA matches the prior1.26.7 receipt.

V343 closure130: all154Preflight checks and embedded structural129 PASS162/1461/781/53; Docker unavailable only. Exact4table/412row synthetic snapshot survives controlled PostgreSQL restart. Five local failures recovered; no product-code change after green gates. Actual operational target and overall roadmap remain open.

V344: former refund host private logging claim did not include delivery failure;0.1.4 makes failed step reports observable to actual run/main and preserves privacy. Before0.1.3/red probe and after hashes/tests retained. Go log/io docs govern method only; execution stays Go1.26.7.

V344 closure132: actual executable report-failure handling verified in0.1.4; full154Preflight checks plus structural162/1461/782/53 PASS, Docker absent. Only3consumer files change and743remain exact. Checkpoint context routes closed history to reuse without deleting evidence; target supervision/retention/alerts and overall roadmap remain open.

V345: bounded-output claim contradicted by real child and fixed1.1.2→1.1.3; before/after pack and red evidence retained. README catalogue still showed1.1.1 and101tests; refreshed from canonical1.1.3/observed JSON with exact before/after receipt. Current Python3.14 docs advertise3.14.7; execution remains hashed3.14.4, no runtime upgrade/admission. Microsoft/systemd docs are research, no source pin or runtime promotion.

V345 native follow-up: twelve successful lifecycle observations were insufficient; injected pre-assignment owner crash left a suspended child. Preserve rejected source; official JOB_LIST/CreateProcess semantics govern correction, fifteen scenarios pass before/after canonical reconstruction. Microsoft moving header used only to inspect constants, no imported source pin or product admission.

V346: capture/handle/cleanup/resource qualification follows current Microsoft HANDLE_LIST/job-limit/accounting and CPython CRT-handle docs. Observed memory peak contradicted the draft ceiling assertion; preserve before/failure and paired native allocation evidence, rename to observed_peak_job_bytes, no clamping/RSS claim. Runtime remains CPython3.14.4; docs3.14.7 do not promote an update.

V347: current official Microsoft CreateFileW sharing/reparse semantics, GetFinalPathNameByHandleW and CreateProcessW were consulted for an isolated AUTHORED launch apparatus. Native denial32 and CRT errno EACCES are separately observed after FAIL598. No signature authority, DLL/runtime identity, sandbox or service admission is inferred from the main-image SHA-256. Canonical evidence: reconstruction_evidence/GOVERNED_NATIVE_LAUNCH_V347.md.

V348 official DuplicateHandle semantics and a real child contradict the broad wait-only security wording: some object types allow greater access; observed event duplication can signal. Preserve before/after, rename COOPERATIVE to EXITED_DURING_GRACE and require trusted cooperative workers. Direct SYNCHRONIZE handle denial remains valid but is not sandbox evidence. Descendant accounting vs filesystem release is separately tested after FAIL600. Evidence: reconstruction_evidence/NATIVE_GRACEFUL_SHUTDOWN_V348.md.

V349: Microsoft MoveFileExW/FlushFileBuffers y CPython os.fsync actuales consultados para semántica, sin adquirir código externo. WRITE_THROUGH y fsync no convierten pruebas de proceso en certificación de pérdida de energía/disco. Documentación Python3.14.7 no modifica ejecución fijada3.14.4. Se conservan límites de metadatos hash y filesystem confiable.

V350: APIs oficiales Microsoft WaitForSingleObject/DuplicateHandle y Go context actuales consultados sin adquirir código/dependencias. Wait finito y join antes de CloseHandle; pre-signaled event se observa sincrónicamente antes del primer claim. Scheduler y suspensión excluyen un SLO garantizado20ms. Runtime Go1.26.7 fijado, sin elevar automáticamente versión por docs.

V351: docs oficiales PostgreSQL18 WAL/initdb contrastadas con runtime local18.6 ya fijado. Su instalación carece del alias timezone/UTC; initdb America/Buenos_Aires conservado y limitación registrada, sin modificar ni adquirir runtime. Assertions de WAL verificadas en instancia propia; no prueba de power-loss/PITR ni target.

V352 autocorrección de orden: TEST06 del contrato solicita un release candidato, TEST08 un artefacto exacto; guide0.1.1 y narrativa histórica añadieron una dependencia de publicación final que aplazaba su ejecución. Before congelado en source147/before148, after guide0.1.2 permite aceptación interna previa, misma identidad de payload o revalidación para TEST07. Se conservan criterios, firma/SCA/notices y estados pendientes; no se declara obsoleta evidencia de código sólo por el cambio documental.

V353 before/after: claim FX ISO4217 implicaba catálogo validado; inspección demuestra sólo regex de tres letras ASCII. Se corrige a ese alcance en fuente canónica sin reescribir informe V217. Consultados https://go.dev/ref/spec#Floating_point_types y https://pkg.go.dev/math#IsNaN el2026-09-09. pkg.go.dev muestra versión posterior, no adoptada; Go1.26.7 fijado permanece. Evidencia: ACCEPTANCE_REAUDIT_FX_BOUNDARY_V353.

V354 before/after: Giftcards decía Issue/top-up idempotente pero rechaza duplicados; Loyalty decía reconstruible desde history pero no tiene restore durable y History recorre toda la historia. Claims canónicos estrechados sin reescribir V194/V198. Go spec actual1.27 no se adopta como toolchain; runtime1.26.7 probado. Deadline24workers no se disfraza de PASS ni de bug upstream corregido.

V355 before/after: reminders0.1 declaraba tenant-scoped con Due/MarkSent globales y envío exactamente una vez; tres probes contradicen aislamiento y doble lectura demuestra ausencia de claim de despacho. API0.2 exige tenant; promesa acotada a flag local. Metadata decía4tests pero fuente contiene5: todos se preservaron. Waitlist retira compatibilidad0.1; V199 permanece historia, no aprobación actual. No consulta/adquisición ni actualización de upstream nueva.

V356 before/after: percent y net fail-closed contradichos por3probes.0.1.1 corrige rangos sin alterar truncado ni prioridad de errores. Nota histórica errónea V209 corregida a V216 desde footer observado (FAIL627); ambos antes preservados. Sin reescribir V195/V216 ni afirmar suites históricas reejecutadas.

V357 before/after: aislamiento tenant-scoped de referrals0.1.0 contradicho por3probes.0.1.1 corrige clave ambigua con tuple; original/V196 preservados. Reward se acota a transición local, no emisión real. No reescritura histórica ni claim productivo nuevo.

V358 before/after: reviews/waitlist tenant-scoped contradicho por4probes. Reviews0.1.1/waitlist0.1.2 corrigen claves, conservan validación y lifecycle. V197/V213 históricos preservados; Notify local y Seat desde pending explícitos, sin notificación integrada inferida.

V359 before/after: onboarding/warranty tenant-scoped contradicho por cuatro grupos públicos.0.1.1 corrige claves, conserva API/errores; Complete es O(N) y orden de Steps no es prerequisite. Diagrama de garantías explicita Resolve sólo desde approved y Close desde rejected/resolved. V206/V214 y fuentes previas conservados.

V360 before/after: tenant-scoped/idempotent-send de marketing0.1.0 contradicho;0.2.0 exige tenant, protege snapshots y limita Send a acuse local. Surveys0.1.1 corrige duplicado de tuplas distintas sin atribuir fuga de Responses. V200/V201 anteriores preservados; complejidad Send O(R) corregida.

V361 before/after: la mención genérica de Docker ausente podía sugerir bloqueo de todo PostgreSQL/avance. Verifier inventaría herramientas amplias; V351 ya prueba PostgreSQL nativo11casos/reinicio11bases, receipt confirmado por SHA sin reejecutarlo. No cambia criterio ni se admite servicio. POS falla overflow/tender/ownership;0.2.0 corrige y obliga errores explícitos. SLO0.1.1 rechaza NaN. El footer POS4/4 es histórico; baselineobservado5tests, sin inventar reejecución V202/V220.

V362 before/after: Help0.1.0 metadata afirmaba delegación GO-SEARCH-CORE, pero código local sólo usa substring;0.1.1 corrige comentario/claim/compatibilidad. Dashboard etiqueta tenant sin comprobar origen de Input: se retira garantía amplia y compatibilidad no integrada. Nuevas pruebas refutan mes00/13/99 y separación por concatenación; before169 conserva bytes. No releases/advisories externos adquiridos ni afirmados. V203/V207 históricos intactos.

V363 before/after: SEO0.1.0 atribuía60/160 a Google; guías oficiales leídas2026-09-09 no fijan ese máximo. i18n regla universal n==1 no implementa CLDR según plural-rules oficial. Owners0.1.1 retiran atribuciones y preservan políticas locales. authority-review.json contiene paráfrasis/URLs/fecha, no snapshot upstream inventado. before171 conserva historia; pruebas/refactor no promueven equivalencia ni observabilityCANDIDATE.

V364: reconstruction_evidence/YARN_EMBEDDED_ASSET_PROVENANCE_V364.md fija un asset Yarn ESM completo (72463bytes decodificados) y5funciones AST contra source fijado;14controles negativos y extracción independiente. Fuente de paquetes no ejecutada.163Markdown de packs,442payload y47notices intactos. FAIL575/publicación completa,FAIL532/nativos y obligaciones de licencia siguen;43/48sin promoción.

V364 cierre175: estructural174 PASS162/1461/802/53 tras corregir cursor de inventario; fallido173 preservado. Evidencia de asset Yarn72463bytes y5funciones parcial;43/48criterios intactos,163Markdown de packs sin cambios.

V365:21cores/42fuentes reconstruidos,222tests ordinarios/23targets con306semillas,vet/build PASS;95,4%statements observados no equivalen a readiness. Se corrigen6compatibilidades ausentes/desfasadas mediante metadata sucesora,12fuentes byte-idénticas.20CONDITIONED/1CANDIDATE,fuera de53planes;43/48intacto,FAIL385 abierto. Evidencia reconstruction_evidence/CORE_COMPATIBILITY_AUDIT_V365.md.

V365 cierre177: structural176 PASS162/1461/803/53, perfiles idénticos;42SHA de fuentes conservados.6sucesiones de metadata no readmiten sus owners.48criterios intactos. Fuente de método SDK1.26.7 fijada; docs online1.27 no sustituyen ese pin.

V366: usuario aclara NEW sin historial, EXISTING con historial y objetivo de entrenar un modelo. PROJECT_HISTORY_MODEL_TRAINING_CONTRACT separa biblioteca, datos privados del consumer, entrenamiento explícito, evaluación y runtime. No pedir corpus real para preparar esta biblioteca ni sustituir entrenamiento por RAG/memoria. Pipeline reusable pendiente; TEST09 sigue BLOCKED, no training/upload/admisión inferidos. Evidencia HISTORY_TRAINING_SCOPE_CORRECTION_V366.md.

V366 cierre179: alcance de entrenamiento corregido por respuesta del usuario y gate178 PASS162/1461/805/53. REQ05/TEST09 cambian definición, ningún estado se promueve;47otras pruebas y163packs Markdown intactos. No dataset, subida, entrenamiento o pipeline nuevo.

V367: investigación oficial de entrenamiento prioriza TRL/SFT y PEFT opcional como DISCOVERED; torchtune/torchforge fuera de primera línea por mantenimiento, torchtitan diferido por scope/runtime. Identidades/artefactos sólo observados en metadata, no adquiridos/admitidos. RESEARCH_INCOMPLETE/FAIL663; TEST09 sigue bloqueado. Evidencia reconstruction_evidence/HISTORY_TRAINING_SOURCE_RESEARCH_V367.md y training_gap_v367/record.json. No pedir corpus privado para esta preparación de biblioteca.

V367 distribución181: los registros training_gap_v367 se conservan íntegros, con SHA-256, como secciones del reporte HISTORY_TRAINING_SOURCE_RESEARCH_V367.md; materializarlos sólo en raíz aislada de investigación. No son nuevos archivos sueltos del release ni un pack admitido.

V367 cierre182: gate estructural181 PASS162/1461/806/53;8registros embebidos recuperables byte-idénticos. Las5superficies investigadas no se admiten por metadata ni por consultas OSV directas.43/48controles sin cambio; siguiente acción adquisición/admisión exacta, sin corpus privado de consumer para mantenimiento.

V368: identidad Git TRL1.11.0/1.12.0 contrastada:580entradas por árbol no truncado,579idénticas y sólo VERSION cambia; commit1.12.0 exacto59c4a8e104413fa9f4ca1a54eaf2ff93c0f299be. No equivale a bytes de paquetes ni admisión/runtime. Spec0.1.2 precisa cierre funcional con capacidades comprometidas ejecutadas y cero bloqueos contradictorios; no cambia48oráculos/estados. Evidencia reconstruction_evidence/TRAINING_SOURCE_IDENTITY_AND_CLOSURE_V368.md.

V368 cierre184: metadata de dos árboles completos580entradas cada uno;579identidades iguales y sólo VERSION cambia. Gate estructural183 PASS162/1461/807/53; código y48oráculos/estados intactos. Cierre funcional explicitado sin reducir alcance. Source/package binding y admisión de entrenamiento siguen pendientes.

V369: core de adquisición0.4.80,36files,127sources y19perfiles internos. Perfil aislado adquirió2sdists TRL1.12.0/PEFT0.20.0 y conserva7outputs inmutables.559archivos inspeccionados:544Git-equal,15metadata de packaging revisados,licencias raíz iguales.98checks transporte; no runtime/training admission. FAIL663 sigue investigación pendiente del grafo y pipeline. Evidencia reconstruction_evidence/TRAINING_SOURCE_ACQUISITION_V369.md.

V370: candidato59distribuciones/103relaciones,58baseTRL+PEFTopcional.59METADATA hash-verified/equivalentes; resolver pip confirma59y2negativos sin wheels/install/framework execution. OSV59versiones0hallazgos declarados;30provenance subjects coinciden,firmas no verificadas. Licencias de artefactos/nativos,adquisición gobernada y runtime/gates pendientes. DISCOVERED/FAIL663 RESEARCH_INCOMPLETE. Ver reconstruction_evidence/TRAINING_DEPENDENCY_GRAPH_V370.md.

V370 FAIL670 before/after: tabla actual1461(Preflight149)→1464(Preflight185). before187/freshness670.json conservan el original; Preflight149 histórico no cambia.

V371: core0.4.81/39files;55quarantine checks+98opaque PASS.59wheels/214055762bytes adquiridos;23854files/23795RECORD hashes inspeccionados.226native files/3SBOMs;371identidades consultadas,6records=4avisos distintos (3security+1maintenance). FAIL675 mantiene cuarentena. Safetensors sourceb7c0f38 corrige pyo3/memmap2,45registry crates0OSV sólo metadata; build pendiente.43/48sin cambio. Ver reconstruction_evidence/TRAINING_WHEEL_QUARANTINE_V371.md.

V371 correction: PATH-only probe did not find cargo/rustc/cl; bounded installed-tool inventory located Rust1.98.0 and MSVC14.50.35717. Preserve before/after; no claim that tools were globally absent. Fixed-source candidate remains unbuilt. Tokenizers guessed lock paths corrected against complete exact tree; no lock present.

V372 maintenance: reference telemetry implementation is returned to canonical worker0.1.5 and WINDOWS_REFERENCE_TELEMETRY_RUNTIME.md/0.1.0, with opt-in WINDOWS_REFERENCE_TELEMETRY_PACK_PLAN.md. Source-only acquirer0.4.82 remains separate from runtime admission. TEST05 closure requires the fresh integrated replay and gates in reconstruction_evidence/INTEGRATED_TELEMETRY_CONTROL_V372.md; no production acceptance or private training data inferred.

V372 cierre196: conteos activos de AGENTS745 y FRANCHISE_ACCELERATOR704 estaban obsoletos. Preflight195 demuestra67packs/754files; se corrigen sólo esos contadores y se conservan before-images en V372. Fuentes/criterios/garantías sin cambios. TEST05 cerrado con evidencia actual, otros4controles bloqueados.

V373: opt-in HISTORY-MODEL-TRAINING-PIPELINE0.1.0/17files, HISTORY_MODEL_TRAINING_PACK_PLAN2packs/21files;26policy tests and actual earlier synthetic SFT/evaluation/rollback. Final exact reconstruction/bootstrap/E2E and TEST09 decision are recorded in reconstruction_evidence/HISTORY_MODEL_TRAINING_CONTROL_V373.md. NEW requires no history; EXISTING keeps data/model private with its own authority/quality/runtime gates. No automatic training, provider call or deployment. SDK component selection precedes AUTHORED glue; see HISTORY_TRAINING_SDK_QUALIFICATION_V373.md and HISTORY_TRAINING_SDK_GAP_V373.md. Integral franchise remains67/754.

V373 cierre201: TEST09 PASS por pack17files/composición21,27policy, bootstrap58wheels/23713file hashes, SFT real, evaluación independiente,13negativos, timeout/no replay, decisiones exactas y rollback local con baseline previo. Preflight200160pasos ejecutados PASS/55planes;164/1506/818/55, integral67/754.45/48, tres pendientes TEST02/03/07; no cambia oráculos ni hereda corpus/modelo/seguridad/deploy del consumer. Ver HISTORY_MODEL_TRAINING_CONTROL_V373.md.

V374: auditoría por claim21núcleos,228tests/306semillas/23campañas (39811771ejecuciones),6comparaciones nuevas de fuentes,42archivos reconstruidos,20packs/40archivos compuestos y candidato rechazado.20CONDITIONED/1CANDIDATE intactos; cero equivalencias empresariales/promociones inventadas.13negativos de enlace de evidencia integrados en VERIFY_LIBRARY. Preflight/cierre del control aún pendientes;45/48sin cambio. Ver reconstruction_evidence/CORE_CLAIM_ADMISSION_V374.md; la integración de negocio más amplia de FAIL385 mantiene sus owners.

V374 cierre203: TEST02 PASS como auditoría por claim21núcleos;228tests/306semillas,23targets finales (6389808ejecuciones),20packs/40files compuestos, candidato rechazado,13negativos de evidencia y Preflight202160pasos PASS.46/48controles;20CONDITIONED/1CANDIDATE y FAIL385/integración más amplia conservados.164/1506/819/55, integral67/754. No cambia48oráculos ni promueve producto. Ver CORE_CLAIM_ADMISSION_V374.md; continuar TEST03 y luego TEST07 según gates materiales.

V374 rectificación204 / FAIL710: el cierre203 de TEST02 se retracta. Fuente/tests/composición son válidos, pero FAIL385 aún exige equivalencia/integración material.45/48PASS; TEST02/03/07BLOCKED. Se preservan historia y mejoras; no se cambia el oracle ni se reduce el roadmap. Ver CORE_CLAIM_ADMISSION_V374.md, encabezado vigente.

V374 successor205: SECURE-OPS1.1.4 corrects the error-rate denominator that hid low-traffic outages. Canonical Prometheus3.14.0 rules9scenarios/13assertions and5wrapper rejections PASS; exact pin/no-engine distinction retained. Three consumer plans now select1.1.4.164packs/1507files/819Markdown/55plans;1260AUTHORED/140ADAPTED/107VERBATIM; integral67/755. No dependency acquisition, core promotion or whole-control closure;45/48 unchanged. Evidence: reconstruction_evidence/CORE_CLAIM_ADMISSION_V374.md. Full successor Preflight pending.

V374 checkpoint207: fullPreflight206 passed162executed steps and all55compositions, including exact Prometheus engine regression. SECUREOPS1.1.4 correction complete; integral67/755, library164/1507/819/55. Current45/48unchanged; TEST02/03/07 remain blocked. Retain the203retraction and actual FAIL385 integration work. See reconstruction_evidence/CORE_CLAIM_ADMISSION_V374.md.

V374 successor208: FRANCHISE_GAP_MAP current header V305 and67/743 were stale. Before SHAa432c3ca065c0dcc0254e73f97fe0ee4aa72194f9611ea3e0b691b131f305fee; after SHAc406a0d3fc165098937df6eb804017f8a36335bd97385dffed2b77d37791b0c2. Current45/48and67/755 plus V372 worker evidence replace the lead-in; dated history preserved. No control promotion. CORE_CLAIM_ADMISSION_V374.md records exact before/after.

V375: HTTP_SLI_ALERT_INTEGRATION_V375.md proves real order HTTP/PostgreSQL → official closed metrics → unchanged canonical alert → repair/replay, and a byte-identical full-profile rebuild. GO-HTTP-METRICS-REFERENCE0.1.0/8files, opt-in68packs/763files; ordinary67/755 unchanged. 380official tests and2authored tests PASS;62public versions/0OSV findings in this reference graph,30compiled upstream license inventories. Trusted synthetic reference only; TEST02/03/07 and FAIL385 remain open,45/48.

V375 checkpoint212: fullPreflight212164executed steps PASS,56profiles composed;165/1515/823/56 and1268AUTHORED/140ADAPTED/107VERBATIM. HTTP reference508requests,four actual failures,canonical ten-minute alert and recovery,one durable order/idempotency/outbox,exact rebuilt executable;380official+2local tests PASS. Only Docker unavailable in that receipt.45/48 unchanged; TEST02/03/07 remain blocked. Next action: FAIL732 official Next/Sharp patched-release qualification, preserving TEST02 owner integration work.

V376: Next16.3.4/Sharp0.35.4 replaces the affected Next16.3.2/Sharp0.35.3 baseline in TS-GO-API-WEB-BRIDGE0.5.5. Five exact public artifacts, ten verified provenance attestations, 8589 installed files identical to inspected archives,152 declared npm versions/0OSV findings versus3baseline advisories,115web tests plus1explicit integration skip, typecheck/build,92connected browser phases and4agenda projects with53exact migrations PASS. Multitab/session gates enabled on all four browser projects. Native AVIF/PNG/WebP smoke confirms libheif1.23.2. This is scoped advisory remediation; full vendored/native security and redistribution obligations remain blocked. Acquisition core0.4.84/46files,193IDs/22profiles and all four transport suites PASS;30dependent source plans updated. HTTP reference0.1.1 records the changed web-parent plan; Go/SQL/runtime evidence is reused only with byte parity.45/48 unchanged. See reconstruction_evidence/NEXT_SHARP_SECURITY_REPAIR_V376.md. Final canonical composition, HTTP binary parity and full Preflight remain pending in this entry.

V376 closure: Preflight214 passed all164executed steps and all56compositions; only Docker unavailable. Canonical165packs/1521files/824Markdown/56plans,1269AUTHORED/145ADAPTED/107VERBATIM; ordinary67/755 and HTTP68/763. Exact Next/Sharp package-advisory repair, signed artifacts and92connected phases/4agendas complete. Current45/48; TEST02/03/07 and native/vendor license/security gates remain blocked. Final canonical consumer has754parent source files byte-identical to the connected candidate plus the two proven Next-generated type imports. HTTP executable matches V375 byte-for-byte, so its unchanged15minute runtime evidence is reused without repetition. Next work returns to TEST02 owner integration, including the public-web SEO/i18n boundary; do not repeat completed browser creation/recovery work or the fixed dependency qualification.

V377 integrates actual Next canonical/robots/sitemap/WebSite JSON-LD under the existing BFF owner0.5.6/43files:5new AUTHORED files and3page changes; no dependency update. Public indexing is opt-in, limited to the real home/model-list pages and public render assets; private pages default noindex.137web tests/1preserved explicit integration skip, strict typecheck/build and3actual HTTP configurations PASS. The same build proves disabled/enabled/changed origin,Host poisoning isolation,mandatory cache revalidation,private noindex and JSON-LD HTML escaping/CSP nonce. Connected qualification: three browser projects passed initially, Firefox failed teardown with known FAIL423/742; a separate exact Firefox23phase rerun plus4agenda projects passed. These results cover4projects/92distinct phases across two runs, not a flawless92phase initial run or a Firefox lifecycle fix. Known upstream incident stays open.45/48 unchanged. See reconstruction_evidence/PUBLIC_WEB_METADATA_INTEGRATION_V377.md; final canonical composition and Preflight pending.

V377 cierre217: integración pública canonical/robots/sitemap/JSON-LD en BFF0.5.6/43archivos comprobada con137tests,3configuraciones HTTP reales y composición final67/760. Regresión conectada:69fases de tres perfiles PASS en el primer run y23Firefox+4agendas PASS en una repetición aislada; primer fallo de teardown conservado y FAIL423/742 abierto. Preflight216:164pasos ejecutados PASS,56planes compuestos;165packs/1526files/825Markdown/56planes,1274AUTHORED/145ADAPTED/107VERBATIM. Docker sigue ausente. Binario HTTP idéntico a V375; fuentes Go/SQL/regla/dependencias sin cambios. Los48criterios/status son idénticos:45/48, TEST02/03/07 bloqueados. Esta frontera SEO de referencia está integrada; faltan i18n/localización y otras equivalencias funcionales, seguridad/admisión nativa y release integral. No contar el replay Firefox como reparación upstream ni repetir creación/recuperación ya probadas.

V378 integra el recorrido público home→modelos→consulta→ubicaciones→solicitud de turno en español/inglés, con fallback español explícito y zona horaria del mercado.12/12recorridos navegador/configuración PASS,12leads y12turnos durables sin duplicados/sobrecupo; pérdida de respuesta y replay probados.164tests/1skip explícito conservado, typecheck/build PASS. El probe HTTP reproduce y corrige labels españoles bajo en-US y offset fijo de3horas bajo UTC. Cuatro archivos nuevos y ocho existentes modificados bajo owners BFF/portales/browser, sin nuevas dependencias ni reglas de dominio.45/48 permanece; no cerrar TEST02 por una sola frontera.

V378 cierre219: recorrido público localizado es/en/fallback y zona horaria configurada,12/12casos navegador/configuración,12leads y12turnos durables sin duplicados/sobrecupo;164tests y typecheck/build PASS. Composición final67/764 y HTTP68/772; Preflight218164pasos ejecutados/56planes PASS, Docker ausente. Sin nuevas dependencias;8589archivos autenticados idénticos y ejecutable HTTP idéntico a V375.45/48, TEST02/03/07 siguen bloqueados. La brecha pública de idiomas queda demostrada; portales privados y demás integraciones mantienen su alcance. Se conservan los fallos iniciales y FAIL423/742 abierto.

V379 conecta búsqueda/lectura de ayuda operativa por permiso y versión exacta: cotizaciones y estado WhatsApp reutilizan sus textos1.0.0.186tests/1skip explícito heredado, typecheck/build y4navegadores PASS sin reintentos; sesión vencida/adulterada, cambio de permiso, versión ausente, lecturas tardías y recuperación503/malformed comprobados. Sin operaciones ni nuevas dependencias. Sólo cierra esta frontera de referencia;45/48 y TEST02/03/07 bloqueados, CMS/capacitación integral y target pendientes.

Portal current metadata corrected from Node24.14.0/Next16.3.2 to retained Node24.20.0/Next16.3.4 (V321/V376); actual dependency bytes unchanged, before/after in canonical-pending-changes and canonical-before. No external assertion refreshed by assumption.

V379 cierre221: búsqueda/lectura de ayuda de cotizaciones y WhatsApp por permiso y versión exacta, compartiendo textos inline;186tests, typecheck/build y4/4navegadores PASS sin reintentos. Composición final67/774 y HTTP68/782; Preflight220164pasos ejecutados/56planes PASS, Docker ausente.8589archivos autenticados y ejecutable HTTP idénticos a evidencia previa. Esta frontera de lectura queda cerrada en referencia;45/48 permanece con TEST02/03/07 bloqueados. CMS, capacitación integral, otras guías y target conservan sus pendientes. No cambios de dependencia, reglas comerciales ni operaciones.

V380 completa el índice de las15/15guías versionadas existentes:13extracciones exactas, permisos por pantalla/sección y enlaces fijados.214tests,13regresiones de HTML inline, typecheck/build y8casos en4navegadores PASS;15perfiles ×15guías ×4navegadores=900decisiones de acceso, más60navegaciones exactas. Tres archivos de operación conservan todos sus handlers idénticos fuera de la ayuda. Sin nuevas dependencias ni reglas.45/48 permanece: no equivale a CMS, guías para tareas que antes carecían de ellas, capacitación integral ni target.

Current inventory corrected from2indexed guides to15existing guides, using preserved source/version hashes and actual permission predicates. No new external authority claim or dependency version.

V380 cierre223: índice15/15de las guías versionadas existentes, permisos por pantalla/sección y enlaces exactos;13extracciones preservan texto/versión/HTML y handlers operativos intactos.214tests, typecheck/build y8casos en4navegadores PASS:900decisiones de acceso y60enlaces UI, sin reintentos. Composición67/777 y HTTP68/785; Preflight222164pasos ejecutados/56planes PASS, Docker ausente. FAIL747/748 corregidos con evidencia;45/48 mantiene TEST02/03/07 bloqueados. El inventario de guías existentes queda cerrado; CMS, tareas sin guía previa, capacitación integral, seguridad/release y target permanecen separados.

V381 integra el pack de panel por rol que estaba fuera del perfil y fallaba2de5tests: navegación compartida con BFF, permisos efectivos y cuatro vistas sin concesión por etiqueta.22tests de rol,236web/1skip previo, typecheck/build y8casos en4navegadores PASS;15perfiles ×5páginas ×4navegadores=300estados de página, además de flag desactivado, sesión inválida,404 y acceso admin denegado. Sin reglas financieras nuevas, operaciones, datos privados ni entrenamiento.45/48 permanece; cierre acotado del panel, no capacitación o integración integral.

Stale role pack0.1.0 metadata/claims preserved in canonical-before; current Next16.3.4 and exact BFF/portal compatibility derived from tested locked composition. No newly asserted upstream facts. Role names no longer imply authority or unsupported financial instructions.

V381 cierre225: panel y cuatro vistas por rol integrados con navegación única BFF y permisos efectivos; pack antes fuera del perfil y2tests fallidos ahora reparado.22tests de rol,236web/1skip previo, typecheck/build,8casos en4navegadores y300estados de página PASS sin reintentos. Composición68/785 y HTTP69/793,4packs reconstruidos,8589artefactos idénticos y binario HTTP sin cambios. Preflight224:164pasos ejecutados/56planes PASS; Docker ausente. Panel cerrado en referencia, opt-in, sin conceder permisos por etiqueta ni inventar reglas/entrenamiento.45/48 mantiene TEST02/03/07 pendientes.

V382 corrige dos defectos reales de lectura: admin:read ya no consulta leads sin lead:read y los portales dejan de mostrar unidades menores como importes mayores. Se comparte exactamente el formatter de cotizaciones, sin alterar precio/estado/aceptación.24regresiones focales,260tests web/1skip previo, typecheck/build y4navegadores contra Go/JWKS/PostgreSQL PASS;5negativos HTTP, consultas indebidas0, escrituras0 y snapshot de6tablas idéntico.45/48 permanece, sin promoción global.

Updated current compatibility metadata from actually selected versions: preserve original dated claims/before hashes. A private route name does not imply its aggregate subrequests share the same permission. Actual screenshot exposed raw minor units; canonical quote algorithm reused exactly, not replaced by invented pricing logic.

V382 cierre227: admin:read funciona sin consultar leads no autorizados; admin/customer muestran importes correctos mediante el formatter exacto de cotizaciones compartido.260tests web/1skip previo,24focales,typecheck/build y4navegadores con Go/JWKS/PostgreSQL PASS;5negativos HTTP,0consultas indebidas,0escrituras y snapshot6tablas idéntico.6packs reconstruidos,composición68/790 y HTTP69/798;8589artefactos y binario HTTP idénticos. Preflight226164pasos ejecutados/56planes PASS; Docker ausente. FAIL754/755/756 corregidos con originales preservados.45/48 conserva TEST02/03/07 pendientes; no cerrar integración,seguridad o release completos por estos fixes.

V383 completa recuperación de lecturas admin/customer/factory: error sin detalles internos, rechazo persistente sin permisos y nueva consulta GET tras corregir sesión; ninguna repetición de operaciones.263tests web/1skip previo,27focales,typecheck/build y4navegadores Go/JWKS/PostgreSQL PASS;5negativos HTTP,0escrituras y snapshot6tablas idéntico. Cinco archivos AUTHORED nuevos; producción Go/SQL/dependencias intactas.45/48 mantiene TEST02/03/07 pendientes.
Fixed local Next16.3.4 error component interface inspected; compatibility versions updated from exact selected composition. Preserve old dated evidence and failed oracles. A framework streaming status is not a substitute for backend authorization/data-disclosure assertions.

V383 cierre229: recuperación admin/customer/factory integrada, error fijo sin datos internos y GET explícito; bearer insuficiente permanece denegado y sesión válida recupera datos.263tests/1skip previo,27focales,typecheck/build y4navegadores Go/JWKS/PostgreSQL PASS;0escrituras/snapshot6tablas idéntico.5packs reconstruidos,68/795 y HTTP69/803,8589artefactos/binario HTTP idénticos. Preflight228164pasos/56planes PASS, Docker ausente.45/48; TEST02/03/07 pendientes, sin cambiar oráculos.

V383 continuation FAIL735/532: exact Next MIT companion captured and byte-matched to the verified repository and three signed packages. Sharp native29notice rows/28versions reconciled; libnsgif owner identified as13vendored blobs at verified libvips426af3f44246fce9cfa8dd51a353aa4dfd48c553, four document/blob matches. No fabricated NetSurf version, binary correspondence, distribution or full native/SCA admission. Raw receipts in V383 notice-analysis; production blocks retained.

V384 demuestra correspondencia exacta del insumo nativo Sharp: libvips-42.dll18614784bytes, versions.json y avisos coinciden con @img/sharp-libvips-win32-x641.3.3 firmado. Perfil previo/2attestations/tamper, cuarentena sin ejecución y4suites del transporte PASS; core0.4.85/48files/194IDs/23perfiles. DLL C++ y build MXE/notices/SCA completos pendientes. Runtime y perfil68/795 intactos;45/48 sin promoción.

FAIL760 corrected from actual source profile selections; original unpublished preparation retained. Runtime parent unchanged. Preserve moving-notice input and C++/MXE/SCA/source obligations; no substituted or executed native runtime.

V384 cierre231: input nativo Sharp1.3.3 firmado y adquirido con perfil; libvips-42.dll18614784bytes, versiones y avisos byte-idénticos. Core0.4.85/48files/194IDs/23perfiles y4suites PASS;30planes actualizados. Preflight230164pasos/56composiciones PASS, Docker ausente. Runtime68/795 y69/803 idéntico aV383.45/48; integración completa, seguridad nativa y release pendientes.

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

V393 FAIL781 REGRESSION_PROVEN: extended-root enumeration gives12571files in both candidate and fresh canonical install, complete maps identical and every hash matches the earlier qualified installation. Ordinary-path counts12555/12544 undercounted16/27files respectively; prior artifacts retained. Signed8534file identity was independently correct throughout. No package content or runtime code changed by this verifier repair.

V393: perfil BFF0.5.16 omite Sharp opcional y desactiva optimización runtime; lock152→122sin versiones nuevas,122stanzas y12571archivos restantes exactos,8534artefactos firmados exactos;294tests/1skip previo,92fases/4agendas/4lecturas/4cancelaciones PASS,OSV122/0.4packs reconstruidos/805files,68/797; HTTPbinario idéntico. No Sharp fix ni reapertura V386; native/tooling restantes yTEST02/03/07 abiertos,45/48. Ver reconstruction_evidence/OPTIONAL_IMAGE_DEPENDENCY_CONTAINMENT_V393.md.

V393 cierre251: cuatro packs incorporados/805files reconstruidos; lock122versiones/0avisos OSV,12571archivos y8534firmados exactos en candidato e instalación limpia.294tests/1skip heredado,92fases/4agendas/4lecturas/4cancelaciones PASS. Preflight250:164pasos PASS/56perfiles; disponibilidad global BLOCKED sólo Docker. Inventario165/1561/841/56,integral68/797.45/48; TEST02integral,TEST03tooling/native restante yTEST07siguen pendientes. Sharp no se instala en este perfil; V386investigación diferida.

V394: PNPM-ARTIFACT-SELECTION-GATE0.3.0/6files corrige FAIL782 (pins empresariales y Playwright obsoletos), entrega BLUEOAK-NOTICE.md ligado a5declaraciones exactas y verifica3recetas actuales;51tests y10mutaciones reales rechazadas. Payload442/22notices intactos; sin ejecución ni admisión runtime/redistribución. Evidencia: reconstruction_evidence/PNPM_CURRENT_ROUTING_NOTICES_V394.md. Los pendientes BlueOak se reducen a entrega downstream/revisión restante; ausencia de LICENSE original chownr no se falsea. Otros permisos/source/publicación de pnpm continúan abiertos.

V394 cierre253: FAIL782 reparado con selector0.3.0/6files;51tests PASS,3recetas actuales,5declaraciones BlueOak exactas con aviso local y10mutaciones reales rechazadas. Preflight252:164pasos PASS/56perfiles, nuevo control de compatibilidad de3consumers; disponibilidad BLOCKED sólo Docker.165/1561/842/56;45/48sin promoción. Payload442 y22notices intactos; pnpm no ejecutado, runtime/redistribución pendientes. Ver reconstruction_evidence/PNPM_CURRENT_ROUTING_NOTICES_V394.md.

V395: PNPM-ARTIFACT-SELECTION-GATE0.4.0/6files entrega QRCODE-NOTICE.md con encabezado de autor/modificación exacto y texto MIT completo, además del aviso BlueOak intacto.59tests PASS,3recetas reales y12negativos rechazados;10source blobs/10regiones de bundle fijados por separado, sin claim de equivalencia completa. FAIL783 entrega local resuelto, FAIL784 metadata ADAPTED reconstruida;payload442/22notices intactos. No ejecución ni admisión runtime/redistribución;45/48. Ver reconstruction_evidence/QRCODE_VENDOR_NOTICE_DELIVERY_V395.md.

V395 cierre257: aviso QRCode con copyright/permiso completos entregado y ligado a3recetas;59tests/12negativos PASS,6fuentes reconstruidas;BlueOak ypayload442 intactos. Preflight256164pasos/56perfiles PASS;Docker ausente.165/1561/843/56;procedencia1306/148/107. FAIL783/784/785 reparados;45/48sin promoción. semver-utils fuente gitHead en GitHub404, investigación pendiente sin bypass ni licencia inventada. Ver reconstruction_evidence/QRCODE_VENDOR_NOTICE_DELIVERY_V395.md.

V395 FAIL786 / OVERBROAD_NONEXECUTION_CLAIM: V394 report incorrectly described the entire Preflight252 as not executing pnpm; preserved logs show its existing offline frozen Playwright/Lighthouse installs. Only the three new recipes were never executed. Original report retained; corrected canonical wording distinguishes planner qualification from existing general runtime gates. No hidden execution, new authorization or changed runtime-admission verdict inferred; no code change. Repeat final structural/state validation after this documentary correction, not unchanged runtime suites.

V396: cerrado descubrimiento y entrega local de licencia original semver-utils1.1.4: MIT OR Apache-2.0 explícito, opción MIT completa y1839bytes originales retenidos. Core0.4.88/51files/196sources/25profiles añade sólo cuarentena npm exact4193bytes/SHA512;140checks y gates anteriores PASS. Planner0.5.0/6files/67tests;3recetas reales/12negativos PASS, BlueOak/QRCode y442payloadfiles intactos. Clave registry vencida documentada; no firma vigente, equivalencia build, ejecución de recetas ni admisión global pnpm.45/48sin promoción; V386/Daybreak sigue diferido. Ver reconstruction_evidence/SEMVER_ORIGINAL_LICENSE_DELIVERY_V396.md.

V396 cierre262: licencia original semver-utils1.1.4 localizada y entregada bajo opción MIT completa;3recetas/12negativos/67tests PASS. Core0.4.88/51files/196sources/25profiles,140checks; planner0.5.0/6files. Preflight261164pasos/56perfiles PASS;Docker ausente.165/1562/844/56;procedencia1307/148/107.45/48sin promoción, Daybreak/V386 diferido. Fuente/licencia semver resuelta en alcance local; firma registry vencida y admisión pnpm restante condicionadas. Ver reconstruction_evidence/SEMVER_ORIGINAL_LICENSE_DELIVERY_V396.md.

V397: PNPM-ARTIFACT-SELECTION-GATE0.6.0/7files entrega el conjunto completo de47textos retenidos/151253bytes en3recetas:141copias verificadas independientemente,12negativos reales y82tests PASS.22notices originales comparados contra payload y25evidencias conservan alcance; BlueOak/QRCode/semver y442payloadfiles intactos. Catálogo ADAPTED con términos por texto, sin algoritmos ni promoción de licencias/source/relinking/runtime.45/48; Daybreak diferido. Ver reconstruction_evidence/PNPM_RETAINED_NOTICE_DELIVERY_V397.md.

V397 cierre264: conjunto47textos/151253bytes entregado en3recetas,141copias exactas/12negativos/82tests PASS. Preflight263164pasos/56perfiles PASS;165/1563/845/56 y1307/149/107. Docker ausente,45/48sin promoción. README/roadmap vigentes sincronizados y cortes previos preservados (FAIL794). Licencias/source/relinking/publicación restante y Daybreak diferido siguen explícitos. Ver reconstruction_evidence/PNPM_RETAINED_NOTICE_DELIVERY_V397.md.

V398: PNPM-ARTIFACT-SELECTION-GATE0.7.0/8files entrega fuente original next-path1.0.0, manifiesto y MPL completa en3recetas:9copias exactas/12negativos reales/90tests PASS. Commit oficial y3Git blobs verificados;4sentencias comparadas bajo adaptadores explícitos,10mutaciones rechazadas. No equivalencia runtime ni reproducibilidad pnpm. Colección47, suplementos anteriores y442payloadfiles intactos.45/48; Daybreak diferido. Ver reconstruction_evidence/NEXT_PATH_MPL_SOURCE_DELIVERY_V398.md.

V398 cierre266: next-path fuente/manifiesto/MPL entregados,9copias exactas/12negativos/90tests/8files PASS. Preflight265164pasos/56perfiles PASS;165/1564/846/56 y1307/150/107. Docker ausente,45/48 sin promoción. Entradas vigentes sincronizadas con historia conservada. Daybreak diferido; integración funcional,source/relinking/publicación restante yrelease pendientes. Ver reconstruction_evidence/NEXT_PATH_MPL_SOURCE_DELIVERY_V398.md.

V399: continuidad de cierre conciliada con V372/V393 y21packs/42fuentes byte-idénticos aV374. No repetir esas suites sin delta; agrupar trabajo por bloqueo material.45/48 intacto; TEST02integración,TEST03grafo/tooling yTEST07release siguen abiertos. Observabilidad candidata0.3.3 separada del runtime probado; Sharp ausente del perfil web no reabre Daybreak. Ver reconstruction_evidence/CURRENT_CLOSURE_BOUNDARIES_V399.md.

V400: encuestas conectadas con PostgreSQL/OIDC/Next, recuperación GET,4navegadores/8respuestas/8POST y retención CLI1/1/0.23archivos nuevos AUTHORED, packs GO-CUSTOMER-SURVEY-API0.1.0 y TS-CUSTOMER-SURVEY-PORTAL0.1.0, CONDITIONED.45/48sin promoción; TEST02/03/07 siguen abiertos. Ver reconstruction_evidence/CONNECTED_CUSTOMER_SURVEYS_V400.md.

V400 cierre270: encuestas conectadas y retención transaccional incorporadas. Go/TSsurvey0.1.1, app1.9.4, HTTPref0.1.15;23archivos nuevos,828/828 reconstruidos.4navegadores/8respuestas/8POST,24concurrentes,restart,CLI1/1/0,fuzz6219;correcciónSQL2→0tablas parciales PASS. Preflight269164pasos/56perfiles PASS más deltaSQL focal;167/1587/850 y1330/150/107,franquicia70/820,HTTP71/828,web7/134.45/48sin promoción; Daybreak diferido. Ver reconstruction_evidence/CONNECTED_CUSTOMER_SURVEYS_V400.md.


V401 / FAIL804 — 2026-09-11: adm-zip0.6.0 bundled in exact pnpm11.25.0 is affected by GHSA-vwc7-r8mq-g2x9 / CVE-2026-76845. Pinned OSV2.5.1 identifies one affected package among476 identified components; this is not a complete recursive native SBOM. The historical V321473/0 observation remains historical and cannot authorize current promotion. pnpm11.26.0 also contains adm-zip0.6.0; rejected as a fix without executing it. pnpm12.4.1 is an isolated inspection candidate with native optional executables, not a clean22-package runtime. Current pnpm promotion and new use are frozen; prior payload retained. Advisory lists no patched release. Reachability UNKNOWN pending call-site examination; no exploit or destructive probe executed. Response lane NO_UPSTREAM_FIX pending official patch review. Daybreak/V386 remains deferred. Evidence: Temp/elite-v401-20260911-closure/nested-scan-result.json and candidate-acquisition-receipts.json. Authority: https://github.com/advisories/GHSA-vwc7-r8mq-g2x9 . Reopen on demonstrated removal/containment or admissible fixed replacement; no risk acceptance.

V401: contención ZIP terminada en candidato pnpm aislado;16módulos adm-zip retirados,441archivos preservados/1bundle cambiado,4rechazos sinIO y3instalaciones offline reales PASS.475identidades conocidas/0hallazgos; SBOMrecursivo no cerrado;12571archivos node_modules idénticos. Pack0.8.0/9files,167/1588; pnpm general BLOCKED, TEST02/03/07 siguen abiertos,45/48; Daybreak diferido. Ver reconstruction_evidence/PNPM_ZIP_CONTAINMENT_V401.md.

V401 reconciliación final: comparación12571/12571idéntica completada antes de la solicitud de detenerla; no se detuvo proceso. FAIL807 separa las3instalaciones iniciales offline de un exec Next que descargó71paquetes por configuración omitida. Exec corregido con store/offline explícitos PASS sin descargas. Historial/log anterior retenido; no prueba global de ausencia de red. reconstruction_evidence/PNPM_ZIP_CONTAINMENT_V401.md


## Decisión sucesora V402 — alcance de infraestructura local

La instrucción explícita del usuario en reconstruction_evidence/LIBRARY_INFRA_SCOPE_V402.md
rige el cierre desde checkpoint272. READY significa biblioteca/infra lista para
usar con referencia materializada fuera de la fuente, contratos y pruebas locales
completos. Producción live permanece separada. No solicitar secretos; ARCA se
completa penúltimo como infraestructura, Daybreak último queda diferido y no impide
el cierre infra. No confundir credenciales pendientes con código incompleto.
AUTHORED queda limitado a glue inevitable: cualquier algoritmo de negocio local
sin source ADAPTED/VERBATIM/DEPENDENCY_PIN admisible requiere resolución real.
Los estados/evidencias previos se conservan; no se marca ningún pendiente PASS por
registrar el nuevo alcance. Orden exacto y nueve criterios en el expediente V402.
