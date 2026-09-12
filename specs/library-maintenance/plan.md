# Plan de mantenimiento — vista sobre el roadmap existente

Versión: 0.1.0. Estado: diseño del delta, no READY_TO_BUILD.

## Baseline y arquitectura conservados

La fuente de verdad permanece en los Markdown canónicos; materialización y toolchains
viven en staging separado. El baseline histórico de este plan fue67 packs/742
archivos; la composición vigente se consulta en FRANCHISE_COMPLETE_PACK_PLAN,
sin usar ese recuento histórico como inventario actual.
No cambiar stack, proveedor ni arquitectura por completar este expediente.

## Orden y dependencias

1. T2801: reconciliar respuestas y pruebas en readiness, blueprint, mapas y assurance.
2. T2802–T2809: ejecutar por verticales y dependencias del roadmap, sin subsistemas paralelos.
3. T2810: integración, reconstrucción independiente y distribución exacta después de cierres.

Primera unidad de verificación de mantenimiento: agente → selección de scope →
compositor/validator → archivos y reportes → checkpoint → diagnóstico y recuperación.
Es un journey CLI, no una pantalla de franquicia. Sus pruebas existentes pueden reutilizarse
por hash. V287 ya demuestra la entrada CLI con 55 checks y tiempos locales de
procesos; faltan el expediente completo, comparación controlada de eficiencia,
tokens y verificación del artefacto final. No repetir esa suite salvo delta o gate
integrado requerido. PROJECT_ENGINEERING_CONTRACT.json enlaza nueve dimensiones
de assurance con pruebas y evidencia; es trazabilidad, no otro roadmap ni aprobación.

## Verificación proporcional

Cambios de selectors: test_library_maintenance_state.ps1 y VERIFY_LIBRARY.ps1.
Cambios de readiness: reconstruir PROJECT_START_READINESS_VALIDATOR y repetir sus tests.
Cambios de ejecución: reconstruir ENGINEERING_EXECUTION_VALIDATOR, validar contrato y checkpoint.
Cambios de producto: gates del perfil y journey exactos, incluyendo negativos y dependencias afectadas.
Ninguna prueba de selector sustituye SCA, E2E, carga, privacidad o recuperación de datos.

## Recuperación

Preservar evidencia fallida y snapshot anterior antes de cambiar el owner canónico.
No usar reset destructivo ni borrar eventos para recuperar. Restaurar sólo el cambio propio
desde el snapshot verificado, reconstruir y repetir los gates afectados.
Los archivos de un proyecto EXISTING requieren inventario/diff antes de cualquier write.
Restauración de PostgreSQL/PITR sigue siendo un gate independiente del target y T2809.

## Autoridad y costo

Trabajo local y consultas oficiales públicas dentro de la solicitud. Ningún pago, instalación
global, cuenta cloud o efecto externo queda autorizado por este plan. Las decisiones ya
expresadas se registran; las nuevas decisiones materiales se explican y solicitan.

## Fin de esta fase

Termina al demostrar readiness y assurance del primer slice real. El estado actual sigue
DISCOVERY/BLOCKED; no generar un PASS rellenando paths con documentos vacíos.

## Evidencia parcial de medición V326

BENCH02 mide el harness y sus etapas por ubicación con2warmups/10muestras; copia13archivos idénticos y55checks por ejecución. TEST45 aporta cuatro fallos de escritura/interrupción en checkpoint con recuperación sintética; TEST04 conserva el scope adicional pendiente. BENCH01 sigue abierto para comparación del trabajo/contexto efectivo del agente; bytes de archivos no son tokens. Ver ENTRY_LATENCY_AND_INTERRUPTION_V326.md, sin segundo backlog ni cierre de producto.

Continuación V326: TEST04 pasó el scope directo materializer/checkpoint con12fronteras. No prueba release distribuido, writers concurrentes ni atomicidad de bridge; TEST06/07/08 permanecen según sus criterios originales.

### Continuidad V327 — guía y diagnóstico CLI

T2801/T2808: EXECUTION-VALIDATOR1.3.1 corrige OS/Unicode no capturados;27tests y reconstrucción15/15/composición23 PASS. Guía ejecutada13comandos/8helps en NEW/EXISTING; instrucciones/BOM/eventos preservados. FAIL537 cerrado focalmente, TEST46 PASS. TEST08 requiere todavía usabilidad del artefacto final; no cierre de macrofrentes ni authority-map automático. Ver reconstruction_evidence/CLI_ENTRY_DIAGNOSTICS_V327.md.

Cierre V327: Preflight85 con151/151pasos PASS, biblioteca161/1453/763/52 y franquicia67/746; Docker mantiene BLOCKED. Checkpoint86 conserva85eventos previos. Sin cambios de código después del gate integrado ni cierre de TEST08/readiness/roadmap.

### Continuidad V328 — readiness y evidencia sin reemplazo

T2801/T2808: readiness0.6.2 corrige diagnóstico y reemplazo tardío de reporte.75tests,8file parity,6archivos/validate_record intactos y guía13/8 PASS. Gate real conserva42bloqueos. FAIL538/539 cerrados focalmente, TEST47 PASS; TEST08/release/readiness y10macrofrentes abiertos. Ver READINESS_REPORT_PUBLICATION_V328.md.

Cierre V328: Preflight87,151pasos PASS y biblioteca161/1453/764/52; readiness75tests. Docker mantiene BLOCKED;42observaciones reales sin cambios. Checkpoint88 conserva87eventos previos. Sin cambios de código después del gate ni cierre de macrofrentes.

### V329 — gap de generación de autoridades delimitado

T2801:4búsquedas/3candidatos oficiales evaluados con G0–G8; expediente validado
NO_ADMISSIBLE_SOURCE/ready=false para selección explicada de autoridades Elite.
No es inexistencia universal ni rechazo global de Spec Kit/AI-DLC/Conductor.
Bootstrap §5.1 precisa aceptación NEW/EXISTING, trazabilidad y pertinencia;
mapa manual actualizado sin atribuir generación automática. Sin código/adopción,
sin cambio de47tests ni42bloqueos de readiness. Evidence: reconstruction_evidence/AUTHORITY_MAP_GAP_V329.md.
Reabrir con candidato exacto compatible, cambio material de fuentes o vencimiento
antes de adopción; no programar este gap hasta USE_REUSABLE_PACK.

Cierre V329: VERIFY_LIBRARY161/1453/765/52 PASS; FAIL540 de empaquetado cerrado sin cambio de selectores. Checkpoint92 conserva91eventos. Gap de authority-map NO_ADMISSIBLE_SOURCE/ready=false; no cierre T2801. Próxima frontera independiente T2804: coordinación multi-tab del recorrido existente, límite explícito en RETURN_OPERATIONS_RECOVERY_V312.md.

### V330 — verificación concurrente de devoluciones existentes

T2804:24carreras entre pestañas en4proyectos,48POST con24éxitos/24conflictos,
8respuestas perdidas tras commit y4copias opener recuperadas. Estado durable
único, digest divergente bloqueado y referencias independientes. Sólo2tests
AUTHORED cambiados;744/746fuentes intactas y reconstrucción exacta. TEST48/EVID40
PASS acotado;48tests=41passed/5blocked/2planned. No coordinación frontend añadida,
no efectos reales downstream ni cierre de T2804/TEST08/readiness42/10macrofrentes.
Evidencia: reconstruction_evidence/RETURN_MULTITAB_VERIFICATION_V330.md.

### V331 — aislamiento de recuperación entre sesiones/roles

T2803/T2804:4proyectos PASS con32transiciones,48rechazos401/403,8lecturas
permitidas a segundo operador y8restauraciones de referencia. Producto intacto;
2tests AUTHORED reconstruidos,744archivos sin cambios. TEST48 ampliado conEVID41;
se conservan48tests/7pendientes:175/12%=14,583333…% sólo de filas de prueba,
sin porcentaje global de trabajo inferido.10macrofrentes/readiness42 abiertos.
Evidencia: reconstruction_evidence/RETURN_SESSION_BOUNDARIES_V331.md.

### V332 — dos bloqueos del empaquetador corregidos

T2808/T2810: reserva exclusiva de ZIP/checksum y rollback de excepción normal;
orden ordinal/NoCompression/epoch estable.50checks canónicos y50en copia3/3;
dos ejecuciones completas sobre mini-fixture producen ZIP idénticos. README
actualizado; no final release, firma ni cierre de TEST06/07/08. TEST04/EVID42
ampliado;48tests/7pendientes y42readiness observaciones conservadas. Evidencia:
reconstruction_evidence/PORTABLE_ARCHIVE_PUBLICATION_V332.md.

### V333 — candidato para resolver la dependencia nativa

Selección pnpm11.25.0:442/455archivos idénticos,13excluidos con6binarios nativos,22notices preservados.4instalaciones copy y8comparativas hardlink pasan sin intentos nativos;115tests/1SKIP,build e instalación Playwright PASS.473identidades/0avisos en scan nuevo. NOT_ADMITTED: falta selector canónico/routing/gates completos/licencias por scope; original FAIL532 BLOCKED. Evidencia PNPM_NATIVE_SELECTION_V333.md. No nuevos consumers, release ni cierre global;48tests/7pendientes.

### V334 — proyección pnpm canónica completada

PNPM-ARTIFACT-SELECTION-GATE0.1.0,4files/17tests dos veces y2proyecciones reales443files idénticas; G0–G8/USE_REUSABLE_PACK sólo tooling en este contexto. Perfil separado1/4 y2pasos de verificación integrados. No consumers/runtime/redistribución admitidos; condiciones explícitas y FAIL532 original conservado. Evidencia reconstruction_evidence/PNPM_SELECTION_GATE_V334.md.

Cierre V334:153pasos PASS,162packs/1457files/772Markdown/53profiles; Docker mantiene disponibilidad BLOCKED.4navegadores conectados PASS y5mediciones Lighthouse reales (performance99mediana;accessibility/best-practices/SEO100).745fuentes intactas y mismo next-env generadoV333.17regresiones/rebuild y proyección443idéntica conservados. FAIL552/553 cerrados; tooling listo dentro de sus condiciones, runtime/redistribución y48tests/7pendientes no se promueven. Siguiente owner: routing canónico y admisión de consumo acotado. Ver PNPM_SELECTION_GATE_V334.md.

### V335 — routing de consumidores pnpm fijados

PNPM-ARTIFACT-SELECTION-GATE0.2.0,6files/43tests fuente+rebuild. Prepare/verify produce recetas offline con configuración aislada, store/cache explícitos y perfiles exactos; nunca ejecuta ni admite runtime.3instalaciones desde rebuild PASS sin descargas/cambios de lock;115tests/1skip ybuild web PASS. FAIL555/557 corregidos sin omitir política; G0–G8 acotado. Perfil separado1/6; producto67/746 inalterado. Evidence PNPM_CONSUMER_ROUTING_V335.md.

Cierre V335:154/154pasos Preflight PASS;162packs/1459files/773Markdown/53profiles; disponibilidad Docker BLOCKED. Pack0.2.0 incorpora routing1/6,43tests/rebuild,3instalaciones offline finales,115tests/1skip ybuild.473declaraciones licenciarias inventariadas, sin admisión total. Checkpoint108 preserva107eventos; runtime/redistribución,48tests/7pendientes y10macrofrentes abiertos.

### V336 — renovación verificable de cache pnpm

Tres checks online desde cache vacía,17metadata oficiales completas con bytes/receipts,3revalidaciones offline sin verdicts históricos y3instalaciones finales offline PASS.275entradas entre3locks;278metadata inventariadas. FAIL561 corregido; no cambio de ejecutables/versiones/política. Procedimiento y límites: reconstruction_evidence/PNPM_CACHE_FRESHNESS_V336.md. Validez acotada a esta ejecución; renovación al reanudar, sin TTL o seguridad continua inferidos. Licencias/runtime/redistribución y10macrofrentes permanecen abiertos.

Cierre V336: VERIFY_LIBRARY PASS162/1459/774/53;3cold policy replays y3instalaciones offline finales sin verdicts históricos;17metadata completas/278mirrors.4textos de licencia/fuente preservados con SHA, sin admisión total. FAIL563 orden de checkpoint corregido; cierre110 conserva109eventos. No código/materialización cambiado; reuse154steps V335.48tests/7pendientes,readiness42 y10macrofrentes abiertos.

### V337 — licencia exacta npm-lifecycle y transporte gobernado

OFFICIAL-UPSTREAM-ACQUISITION-CORE0.4.79:33files/125sources,78opaque checks/rebuild y18perfiles de adquisición. Nuevo perfil separado adquiere sólo npm-lifecycle1100.1.0; firma registry e integridad verificadas. LICENSE del tar exacto coincide con commit/sidecar,7de8archivos con Git blobs; único delta packageManager omitido.29planes de composición actualizados. Fuente licencia verificada, runtime/redistribución no admitidos. Ver reconstruction_evidence/PNPM_LIFECYCLE_LICENSE_V337.md.

V337 closure114: npm-lifecycle1100.1.0 exact licence evidence and core0.4.79 scoped acquisition extension integrated. 78opaque tests,18source profiles,33file rebuild and154step integrated Preflight PASS;53composition counts exact. Overall availability BLOCKED only Docker. Fixed-source licence texts for individual/qrcode-terminal retained; next work resolves remaining473graph notice/use evidence and semver-utils source text without promoting original full-native FAIL532. Macrofronts and production blockers unchanged. Report: reconstruction_evidence/PNPM_LIFECYCLE_LICENSE_V337.md.

V338 notice coverage:473metadata hashes verified;468identities observed from456bundle identities+21external manifests+root (with overlap).22notice files and21embedded attribution paths inventoried;32exact draft texts preserved. BlueOak five external packages and QRCode vendored MIT require explicit notice treatment; nested works, semver-utils exact text and whole-bundle source/use obligations remain open. All442payload files unchanged; no runtime/redistribution admission. Evidence: reconstruction_evidence/PNPM_NOTICE_COVERAGE_V338.md.

V338 closure117: notice matrix and32text draft verified; structural PASS162packs/1461files/776Markdown,53profiles unchanged. Ledger closure references corrected and uniqueness rechecked. Continue exact semver-utils acquisition and nested-work/source/notice qualification; runtime,redistribution,original FAIL532 and production gates remain open.

V339:8fixed source manifests and7licence texts verified; four BlueOak source texts,3Noble locked build-input identities outside the original473graph,70node-gyp vendored files identical to fixed tree. Draft39preserves32previous texts. chownr declaration-only, semver-utils exact text and source/bundle/redistribution conditions remain open. All442selected files unchanged. See reconstruction_evidence/PNPM_NESTED_LICENSE_SOURCES_V339.md.

V339 closure119: eight source manifests,seven source licence texts,70gyp file matches and draft39 verified; structural PASS162/1461/777 with53profiles unchanged. Next independent work: per-version Yarn/undici nested-source notices; semver/chownr and source/relinking remain open. Original473inventory preserved and all442selected payload files unchanged.

V340:103external Undici6files match fixed source; Undici7licences/headers bound separately. Yarn4.1.7 registry gitHead returns404 but official tag/source version exists at a distinct commit; preserve both, no published-byte equivalence. Root BSD and3MIT source comments retained; draft47preserves39previous texts,442selected files unchanged. FAIL575 scoped provenance remains unresolved. See reconstruction_evidence/PNPM_YARN_UNDICI_NOTICES_V340.md.

V340 closure121: 103 external Undici6 file matches, separate Undici7 notices, Yarn tagged BSD/three MIT comments and draft47 verified; structural PASS 162/1461/778 with 53 profiles unchanged. FAIL575 published Yarn binding stays open. Next: exact formdata/ws source identities and transformed-source mapping, remaining licence applicability and source/relinking/delivery. All 442 selected payload files unchanged; no runtime/redistribution promotion.

V341: five selected Undici7 AST units in two bundled modules match fixed-source syntax under declared identifier renames; nine initial controls plus eight actual-source mutations and two exact byte-range checks pass. No whole-module/runtime/redistribution claim; draft47/payload442 unchanged. Contract remains41passed/7pending of48 (approximately15% of rows, never global effort); ten macrofronts remain open. Evidence: reconstruction_evidence/PNPM_SCOPED_SOURCE_CORRESPONDENCE_V341.md.

V341 closure124: structural123 PASS162/1461/779/53;47draft texts and442selected payload files unchanged. Five scoped Undici7 AST matches qualified,9initial controls plus8actual-source negative mutations. Local parser-path/portable-reference failures recovered; Yarn provenance and other source/licence/relinking gates stay open.41passed/7pending of48 remains a row count only.

V342 / T2809: GO-OBSERVABILITY-CORE0.3.1 corrects four false-success writer cases (FAIL581); two canonical files rebuilt exactly. Real HTTP/file reference32requests/64correlated records plus closed-file failure/replacement,28tests x3,vet/build and3fuzz targets10s PASS. Actual runtime host/reporter/export/retention/alerts and admission remain pending; isolated test integration is not production wiring. See reconstruction_evidence/OBSERVABILITY_DELIVERY_INTEGRATION_V342.md.

V342 closure126: writer delivery defect581 fixed in canonical0.3.1;28tests x3,17seed cases,three10s fuzz targets and real32request/64record HTTP-file reference PASS. Structural125 PASS162/1461/780/53. Actual selected host/supervisor/reporter/export/retention/alerts remain the next T2809 boundary; no global percentage or test-row closure inferred.

V343 in progress: concrete bounded JSONStatusReporter replaces the abstract-only reporting gap in the existing WhatsApp host; canonical0.13.0, four blocks rebuilt exactly,67/746composition. Local real-worker/PostgreSQL→TCP and cancellation/failure/concurrency tests pass x3. Final gates interrupted by computer shutdown are pending rerun, not PASS. Actual service mounting/supervisor/collector/retention/alerts remain T2809; no percentage. Evidence: reconstruction_evidence/WHATSAPP_STATUS_REPORTER_V343.md.

V343 / T2809: WhatsApp0.13.0 now supplies a concrete bounded JSON reporter to the existing status host. Real PostgreSQL→TCP and mid-record loss after commit/restart prove no job/observation/audit replay. Final14host tests x3,11seeds,10s fuzz298870executions,full Go baseline/vet/build and4/4canonical parity PASS;67/746unchanged. Reporter implementation gap closed locally; real service identity/mounting/supervisor/collector/retention/alerts remain. See reconstruction_evidence/WHATSAPP_STATUS_REPORTER_V343.md.

V343 closure130: bounded status reporter and post-commit report-loss recovery canonical0.13.0;14tests x3/11seeds/fuzz298870/full Go gates PASS. Embedded structural129 plus154Preflight checks PASS,67/746and53profile counts preserved; Docker absent. Controlled restart4tables/412rows identical, owned PG stopped. Next T2809 existing service/supervisor signal path; no false total-completion claim.

V344 / T2809: refund host0.1.4 now stops before another token/step on failed/incomplete reporting; fixed errors/exit1 and healthy/idle policy preserved. Six baseline red modes →green; actual loop/closed-file/subprocess,17tests x3,14seeds,fuzz327445,vet/build and15file parity PASS. Three changed files,12unchanged including business/provider/DB/locks;67/746consumer. One opt-in PG test skipped, no live/refund-policy change or supervisor activation. See reconstruction_evidence/REFUND_HOST_REPORT_DELIVERY_V344.md.

V344 closure132: actual refund host0.1.4 stops on report failure with static error/exit1;17tests x3,14seeds,fuzz327445,binary and15file parity PASS. Preflight131154steps/structural162/1461/782/53 PASS; Docker absent. Consumer delta3/746,743unchanged. Compact checkpoint keeps current operational owners primary, closed history reusable. Continue T2809 reference supervision/collection/retention/alerts; no global completion claim.

V345 / T2809: SECURE-OPS1.1.3 corrects finite process capture/deadline/privacy; real-child cap/timeout/EOF tests and actual refund binary verified. Service-supervisor implementation remains missing and CONDITIONED; qualification continues through the capability-gap gate. See reconstruction_evidence/OPERATIONAL_PROCESS_BOUNDARY_V345.md.

V345 closure134: finite operational runner1.1.3 rebuilt; all154Preflight checks PASS and inventory162/1461/783/53, Docker alone absent. Native Windows qualification adds creation-bound job membership, root/descendant cleanup, invalid-job and breakaway refusal, abrupt/pre-resume owner death:5scenarios x3 plus exact rebuild repeat. Prior launch-window defect retained/rejected; no service supervisor admission. Next operational I/O/handle inheritance and resource/privilege/lifecycle qualification, then G0–G8.

V346 / T2809: native candidate now combines creation-bound JOB_LIST, explicit inherited HANDLE_LIST, bounded separate pipe capture, EOF+tree completion, deadline/cancel/cleanup and process/commit budget probes.23tests+15base lifecycle observations PASS after canonical rebuild. Actual refund fixture retains exit1. Product SECURE-OPS1.1.3 and profiles unchanged; candidate remains CONDITIONED until identity/launch governance, shutdown/lifecycle, result delivery and service gates. See reconstruction_evidence/NATIVE_SUPERVISOR_IO_QUALIFICATION_V346.md.

V346 closure 136: native I/O/resource qualification has 23 tests and 15 lifecycle observations PASS after exact reconstruction. Preflight 135: 154 checks PASS; inventory 162/1461/784/53; Docker unavailable. Product packs/profiles unchanged. Full service-supervisor remains CONDITIONED; next governed executable identity/launch, lifecycle/shutdown and durable operational results before G0–G8 admission.

V347 / T2809: governed native launch qualification adds an externally digest-bound strict profile, pinned executable/cwd path handles, held-handle file hashing, reparse/alias rejection and explicit argv/environment/budgets. Canonical reconstruction:20 launch tests +23 prior capture/resource tests PASS without skips;16 concurrent launches,24 bad-hash acquisitions without handle growth. Product owner/profiles unchanged; candidate CONDITIONED. Next lifecycle/cooperative shutdown and durable result delivery, then complete target identity/security/runtime G0–G8. Evidence: reconstruction_evidence/GOVERNED_NATIVE_LAUNCH_V347.md.

V347 closure 138: governed main-image/profile qualification closed with20+23 tests PASS/0 skips and exact canonical reconstruction; Preflight137 all 154 executed checks PASS, inventory162/1461/785/53, Docker absent. Product owner/profiles unchanged. Next cooperative shutdown/lifecycle and durable operational results; trusted profile authority/runtime/security and target admission remain conditions. No global percentage.

V348 / T2809: native candidate adds opt-in strict v2 stop-event protocol, finite grace, continued bounded capture and owned-job force fallback. Rebuilt19shutdown+20launch+23capture tests PASS. EXITED_DURING_GRACE is lifecycle only; trusted workers required, same-user event amplification remains explicit. Product owner/profiles unchanged. Next durable operational result delivery and actual worker lifecycle mounting; full identity/runtime/security/service/target admission remains conditioned. Evidence: reconstruction_evidence/NATIVE_GRACEFUL_SHUTDOWN_V348.md.

V348 closure140: opt-in trusted-worker shutdown/grace qualification,19+20+23tests PASS/0skips after exact rebuild. Preflight139 all 154 executed checks PASS; inventory162/1461/786/53, Docker absent. Same-user event amplification and lifecycle-only EXITED_DURING_GRACE are explicit; owned-descendant fixture release proved. Product packs/profiles unchanged. Next durable operational results and actual worker mounting; full service/target admission stays CONDITIONED.

V349 / T2809: reserva exclusiva por run ID antes de ejecutar, claim con flush y publicación final sin reemplazo; lectura estricta ligada a run/profile/claim.28tests nuevos +62 previos PASS/0skips tras rebuild exacto. Caídas de proceso y8procesos concurrentes no duplican el efecto sintético; incertidumbre nunca habilita replay automático. Metadatos/digests, sin salida bruta. Product owner/perfiles intactos; no garantía de energía/disco/WORM/negocio. Next montaje del worker real y admisión servicio/target. Evidencia: reconstruction_evidence/NATIVE_RESULT_PUBLICATION_V349.md.

V349 closure142: local result-store qualification,28new+62prior tests PASS/0skips after exact rebuild. Preflight141 all 154 executed checks PASS, inventory162/1461/787/53; Docker absent. Exclusive run reservation, metadata publication and process-crash readback preserve uncertainty without replay. No product owner/profile changes or power-loss/WORM/business-success claim. Next actual worker lifecycle/result mounting and full service/storage/target admission.

V350 / T2809: puente Win32 de apagado montado en el bucle y Processor.Step reales, con Store/Provider sintéticos.20tests nuevos +90supervisor PASS; Go x3, vet/build y fuzz10s PASS. Resultado del proceso separado del acknowledgement de ciclo y del resultado sintético de dominio; cancelación/efecto incierto no habilita replay.15archivos refund intactos; native runtime unchanged, capture fixture strengthened. Sigue PostgreSQL real/startup/config/identidad/seguridad/target; ningún pack de servicio admitido. Evidencia: reconstruction_evidence/NATIVE_WORKER_LIFECYCLE_V350.md.

V350 closure145: native stop connected to unchanged real refund loop/processor through synthetic Store/Provider;20new+90prior tests PASS, Go x3/vet/build/fuzz2922431 PASS. Preflight144 all 154 executed checks PASS;162/1461/788/53, Docker absent. Normal startup/config/PostgreSQL-backed acknowledgement and full service/identity/security/storage/target remain open. No product pack/profile promotion.

V351 / T2809: calificación nativa conecta startup run(), Processor y PostgresStore reales en11casos con proveedor sintético y seed con constraints/triggers activos; snapshots de11bases idénticos tras reinicio controlado, servidor propio detenido.746archivos del perfil exactos;90regresiones supervisor, Go x3/vet/build y fuzz configuración7semillas/3687878ejecuciones PASS. Pendientes identidad/config/secretos de servicio, collector/retención/alertas, seguridad y aceptación target; ningún pack promovido. Evidencia: reconstruction_evidence/NATIVE_POSTGRES_LIFECYCLE_V351.md.

V351 closure147: startup/Processor/PostgresStore montados en calificación nativa;11casos PostgreSQL y reinicio11bases idénticas,90regresiones, Go x3/vet/build/fuzz3687878 PASS. Preflight146 154checks PASS;162/1461/789/53, Docker ausente.746archivos producto exactos. Pendientes identidad/secretos/collector/retención/alertas/seguridad/target;41/48controles aprobados no es porcentaje del objetivo.

### V352 — orden de aceptación del candidato

Corrección de procedimiento FAIL613: se permite preparar y probar una copia portable completa, interna y congelada antes del release final. TEST06/08 evalúan los mismos criterios sobre ese archivo exacto; su evidencia debe incluir SHA del ZIP, manifest y fuentes consumidas. TEST07 sigue bloqueado hasta firma, seguridad/licencias, gates materiales y reconstrucción independiente; la publicación final exige identidad con el candidato aceptado o revalidación de los cambios. No se altera el denominador, elimina un test, promueve un pack o sustituye un proveedor real por fixture. Este orden materializa la aceptación antes de la publicación, sin crear otro roadmap.

V352 / T2808-T2810: TEST06 y TEST08 pasan contra candidato portable completo148, SHA0590588c4ab139e86e581b13c1db561f9342462646105c109f190aeefc755085.804fuentes+manifest,13comandos guía/8help, fallo real de instalación y rollback exacto NEW/EXISTING, cadena preservada/reanudada y55checks entrada. Orden corregido: aceptación antes de publicación; TEST07 sigue bloqueado, payload cambiado exige revalidación.43/48checks PASS,5bloqueados; no porcentaje global ni cierre de macrofrente. Evidencia reconstruction_evidence/PORTABLE_CANDIDATE_ACCEPTANCE_V352.md.

V352 closure150: TEST06/08 aceptados contra candidato148 exacto;43/48PASS y5bloqueados. Preflight149 154checks PASS,162/1461/790/53; Docker ausente. Mismos9owners consumidos; cambios de evidencia posteriores al freeze no se heredan como aceptación automática de otro ZIP. TEST07 y restantes2/3/5/9 siguen abiertos; sin release final ni porcentaje de esfuerzo.

V353: auditoría V35231hashes/18recibos y aceptación completa148 repetida PASS. Corrección real GO-FX-CORE0.1.0→0.1.1: NaN/infinito/overflow/underflow-cero rechazados, tasa previa preservada;2fuentes reconstruidas,9tests x3,14semillas y13.047.867fuzz ejecuciones PASS. Sigue CONDITIONED/fuera del perfil; TEST02 parcial,43/48sin cambio. El ZIP148 NO contiene el fix; TEST07 exige nuevo payload y revalidación. Evidencia reconstruction_evidence/ACCEPTANCE_REAUDIT_FX_BOUNDARY_V353.md.

V353 cierre152: Preflight151 154controles ejecutados PASS;162/1461/791/53, Docker ausente.43/48sin cambios de criterios ni estados. FX0.1.1 corregido/reconstruido/fuzz PASS pero CONDITIONED; TEST02 parcial. Candidato148 auditado y repetido conserva su alcance histórico; no contiene este fix. Sin release final ni porcentaje de esfuerzo.

V354 / TEST02 parcial: GO-LOYALTY-CORE y GO-GIFT-CARDS-CORE0.1.1 corrigen overflow int64, alias de cuentas y IDs cross-tenant; giftcards rechaza ID vacío y sincroniza String.11tests originales verdes,5regresiones rojas;4fuentes canónicas,19tests x3,vet/build y2gates fuzz de2targets a4workers PASS. Deadline24workers preservado como límite sin causa demostrada. CONDITIONED/fuera del perfil;43/48sin cambios. Evidencia reconstruction_evidence/CREDIT_CORE_ISOLATION_V354.md.

V354 cierre154: Preflight153 154controles ejecutados PASS,162/1461/792/53; Docker ausente. Dos cores corregidos0.1.1,4fuentes,19tests x3 y5839355fuzz ejecuciones en dos gates4workers PASS. Deadline original24workers preservado, sin claim de -race o causa probada.43/48 y criterios intactos; TEST02/03parciales, ambos packs CONDITIONED/fuera del perfil; sin release final.

V355: GO-REMINDERS-CORE0.2.0 exige Due(tenant)/MarkSent(tenant,id) y claves tenant/ID;3fallos de aislamiento reproducidos,5aserciones originales preservadas con nuevos argumentos.2fuentes exactas,11tests x3,vet/build y fuzz10s/4workers/2112907ejecuciones PASS. Caller0.1.x rechazado al compilar; no shim global ni promesa de envío único. Waitlist0.1.1 sólo retira compatibilidad insegura, fuentes intactas. TEST02/03parciales,43/48sin cambio; ambos CONDITIONED/fuera del perfil. Evidencia reconstruction_evidence/REMINDER_TENANT_API_V355.md.

V355 cierre156: Preflight155 154controles ejecutados PASS,162/1461/793/53; Docker ausente. Reminders0.2 con tenant explícito,2fuentes,11tests x3,fuzz2112907PASS y migración obligatoria de callers. Waitlist0.1.1 sólo metadata.43/48criterios/estados intactos; evidencia parcial, no exactly-once, admisión integral ni release final.

V356: promotions/payroll0.1.1 corrigen overflow de porcentajes y neto negativo que wrappeaba positivo; sin nuevas reglas ni firmas.11tests históricos intactos,3grupos red,4fuentes reconstruidas,18tests x3 y37semillas/4614438fuzz ejecuciones PASS. Truncado, expiración inclusiva y prioridad ErrInvalidRun preservados. TEST02 parcial,43/48intactos, ambos CONDITIONED/fuera del perfil. Evidencia reconstruction_evidence/BASIS_POINT_ARITHMETIC_V356.md.

V356 cierre158: Preflight157 154controles ejecutados PASS,162/1461/794/53; Docker ausente. Promotions/payroll0.1.1 corrigen overflow con reglas/truncado intactos;4fuentes,18tests x3,4614438fuzz ejecuciones PASS.43/48criterios/estados sin cambio, TEST02parcial; CONDITIONED y sin release final.

V357: referrals0.1.1 corrige alias tenant/referee que permitía transiciones de otro referido. API/precedencia intactas;5tests históricos,3regresiones,2fuentes reconstruidas,10tests x3,vet/build,12semillas/1709469fuzz PASS. TEST02/03 parciales,43/48intactos; AUTHORED/CONDITIONED fuera del perfil, Reward sólo estado local. Evidencia reconstruction_evidence/REFERRAL_IDENTITY_ISOLATION_V357.md.

V357 cierre160: Preflight159 154controles ejecutados PASS,162/1461/795/53; Docker ausente. Referrals0.1.1:2fuentes,10tests x3,1709469fuzz PASS.43/48definiciones/estados sin cambio, TEST02/03parciales; CONDITIONED, sin release final.

V358: reviews0.1.1/waitlist0.1.2 corrigen alias en moderación y lookup/transiciones de cola. API/estados intactos;11tests históricos,4red,4fuentes,19tests x3,vet/build,24semillas/4216603fuzz PASS. TEST02/03 parciales,43/48sin cambio; CONDITIONED/fuera del perfil. Evidencia reconstruction_evidence/MODERATION_WAITLIST_ISOLATION_V358.md.

V358 cierre162: Preflight161 154controles ejecutados PASS,162/1461/796/53; Docker ausente. Reviews/waitlist4fuentes,19tests x3,4216603fuzz PASS. Acumulado V357–358 tres cores/6fuentes/29tests x3/5926072fuzz PASS.43/48sin cambiar criterios/estados; TEST02/03parciales, CONDITIONED/sin release.

V359: onboarding/warranty 0.1.1 corrigen alias de identidad que permitían leer/completar checklist o transicionar un reclamo ajeno. API/errores conservados; 9 tests históricos intactos, 4 grupos red, 4 fuentes reconstruidas, 17 tests x3, vet/build y 30 semillas/4760595 fuzz PASS. TEST02/03 parciales, 43/48 sin cambios; AUTHORED/CONDITIONED fuera del perfil. Evidencia reconstruction_evidence/ONBOARDING_CLAIM_ISOLATION_V359.md.

V359 cierre164: Preflight163 154 controles ejecutados PASS, 162/1461/797/53; Docker ausente. Onboarding/warranty0.1.1: 4 fuentes, 17 tests x3, 4760595 fuzz PASS.43/48 criterios/estados sin cambio; TEST02/03 parciales, CONDITIONED y sin release final.

V360: surveys0.1.1/marketing0.2.0 corrigen identidad y ownership. Marketing exige tenant en Due/Send/Remaining y copia Recipients; no dispatch real.11tests históricos conservan aserciones con migración de tenant/map index;5red,4fuentes,20tests x3,vet/build,3firmas antiguas rechazadas y24semillas/3864824fuzz PASS. TEST02/03 parciales,43/48intactos; CONDITIONED/fuera del perfil. Evidencia reconstruction_evidence/SURVEY_CAMPAIGN_SCOPE_V360.md.

V360 cierre166: Preflight165 154 controles ejecutados PASS,162/1461/798/53; Docker ausente. Surveys0.1.1/marketing0.2.0:4fuentes,20tests x3,3864824fuzz PASS y3firmas antiguas rechazadas.43/48criterios/estados sin cambio; TEST02/03parciales, CONDITIONED/sin dispatch ni release final.

V361: POS0.2.0 calcula importes con errores explícitos, evita overflow/tender negativo, copia Lines y aísla identidad. SLO0.1.1 rechaza NaN.11aserciones históricas conservadas con helpers de error POS;5red,4fuentes,20tests x3,vet/build,2firmas legacy rechazadas y34semillas/3381083fuzz PASS. Docker ausente no bloquea esta ruta; V351 nativo11casos/reinicio11bases confirmado por SHA, sin nueva corrida ni admisión target. TEST02/03parciales,43/48intactos. Evidencia reconstruction_evidence/POS_SLO_BOUNDARIES_AND_DOCKER_SCOPE_V361.md.

V361 cierre168: Preflight167 154controles ejecutados PASS,162/1461/799/53. POS0.2.0/SLO0.1.1:4fuentes,20tests x3,3381083fuzz PASS,2firmas legacy rechazadas. Docker no disponible para contenedores, sin bloquear este trabajo ni la ruta nativa V351 confirmada por receiptSHA. Verifier/48criterios intactos,43/48sin cambio; TEST02/03parciales y sin release final.

V362: dashboards/helpcenter0.1.1 corrigen mes00/13/99, alias de tenant/id y overflow de versión con ErrVersionExhausted. Metadata limita KPI a datos caller y Search a substring local; sin integración inventada.8tests históricos,5red,4fuentes,19tests x3,vet/build,34semillas/4266789fuzz PASS.17cores con fixes acotados no equivalen a21admitidos;43/48intactos. Evidencia reconstruction_evidence/DASHBOARD_HELP_BOUNDARIES_V362.md.

V362 cierre170: Preflight169 154controles ejecutados PASS,162/1461/800/53. Dashboards/helpcenter0.1.1,4fuentes,19tests x3,4266789fuzz PASS. Verifier/48criterios intactos,43/48sin cambio;17cores con fixes acotados no son21admitidos. Sin integración inventada ni release final.

V363: i18n0.1.1/SEO0.1.1/social0.2.0 corrigen formas vacías, scope exacto, reglas inyectadas y API social con tenant.6fuentes,15tests históricos conservados con migración de argumentos,8red,32tests x3,vet/build,3legacy negativas y62semillas/4518771fuzz PASS. Observability0.3.1 sin cambios,28tests/vet/build PASS,CANDIDATE. Censo21owners/20cores con fixes acotados no implica21admitidos.43/48intactos. Evidencia reconstruction_evidence/PUBLIC_CORE_SCOPE_V363.md.

V363 cierre172: Preflight171 154checks ejecutados PASS,162/1461/801/53. i18n/SEO/social6fuentes,32tests x3,4518771fuzz PASS y3legacy negativas. Observability28tests/vet/build aislados,CANDIDATE.21owners por SHA,43/48criterios intactos; admisión/equivalencia/servicio/security/release/histórico siguen abiertos.

V364: reconstruction_evidence/YARN_EMBEDDED_ASSET_PROVENANCE_V364.md fija un asset Yarn ESM completo (72463bytes decodificados) y5funciones AST contra source fijado;14controles negativos y extracción independiente. Fuente de paquetes no ejecutada.163Markdown de packs,442payload y47notices intactos. FAIL575/publicación completa,FAIL532/nativos y obligaciones de licencia siguen;43/48sin promoción.

V364 cierre175: estructural174 PASS162/1461/802/53 tras corregir cursor de inventario; fallido173 preservado. Evidencia de asset Yarn72463bytes y5funciones parcial;43/48criterios intactos,163Markdown de packs sin cambios.

V365:21cores/42fuentes reconstruidos,222tests ordinarios/23targets con306semillas,vet/build PASS;95,4%statements observados no equivalen a readiness. Se corrigen6compatibilidades ausentes/desfasadas mediante metadata sucesora,12fuentes byte-idénticas.20CONDITIONED/1CANDIDATE,fuera de53planes;43/48intacto,FAIL385 abierto. Evidencia reconstruction_evidence/CORE_COMPATIBILITY_AUDIT_V365.md.

V366: usuario aclara NEW sin historial, EXISTING con historial y objetivo de entrenar un modelo. PROJECT_HISTORY_MODEL_TRAINING_CONTRACT separa biblioteca, datos privados del consumer, entrenamiento explícito, evaluación y runtime. No pedir corpus real para preparar esta biblioteca ni sustituir entrenamiento por RAG/memoria. Pipeline reusable pendiente; TEST09 sigue BLOCKED, no training/upload/admisión inferidos. Evidencia HISTORY_TRAINING_SCOPE_CORRECTION_V366.md.

V367: investigación oficial de entrenamiento prioriza TRL/SFT y PEFT opcional como DISCOVERED; torchtune/torchforge fuera de primera línea por mantenimiento, torchtitan diferido por scope/runtime. Identidades/artefactos sólo observados en metadata, no adquiridos/admitidos. RESEARCH_INCOMPLETE/FAIL663; TEST09 sigue bloqueado. Evidencia reconstruction_evidence/HISTORY_TRAINING_SOURCE_RESEARCH_V367.md y training_gap_v367/record.json. No pedir corpus privado para esta preparación de biblioteca.

V367 distribución181: los registros training_gap_v367 se conservan íntegros, con SHA-256, como secciones del reporte HISTORY_TRAINING_SOURCE_RESEARCH_V367.md; materializarlos sólo en raíz aislada de investigación. No son nuevos archivos sueltos del release ni un pack admitido.

V368: identidad Git TRL1.11.0/1.12.0 contrastada:580entradas por árbol no truncado,579idénticas y sólo VERSION cambia; commit1.12.0 exacto59c4a8e104413fa9f4ca1a54eaf2ff93c0f299be. No equivale a bytes de paquetes ni admisión/runtime. Spec0.1.2 precisa cierre funcional con capacidades comprometidas ejecutadas y cero bloqueos contradictorios; no cambia48oráculos/estados. Evidencia reconstruction_evidence/TRAINING_SOURCE_IDENTITY_AND_CLOSURE_V368.md.

V369: core de adquisición0.4.80,36files,127sources y19perfiles internos. Perfil aislado adquirió2sdists TRL1.12.0/PEFT0.20.0 y conserva7outputs inmutables.559archivos inspeccionados:544Git-equal,15metadata de packaging revisados,licencias raíz iguales.98checks transporte; no runtime/training admission. FAIL663 sigue investigación pendiente del grafo y pipeline. Evidencia reconstruction_evidence/TRAINING_SOURCE_ACQUISITION_V369.md.

V370: candidato59distribuciones/103relaciones,58baseTRL+PEFTopcional.59METADATA hash-verified/equivalentes; resolver pip confirma59y2negativos sin wheels/install/framework execution. OSV59versiones0hallazgos declarados;30provenance subjects coinciden,firmas no verificadas. Licencias de artefactos/nativos,adquisición gobernada y runtime/gates pendientes. DISCOVERED/FAIL663 RESEARCH_INCOMPLETE. Ver reconstruction_evidence/TRAINING_DEPENDENCY_GRAPH_V370.md.

V371: core0.4.81/39files;55quarantine checks+98opaque PASS.59wheels/214055762bytes adquiridos;23854files/23795RECORD hashes inspeccionados.226native files/3SBOMs;371identidades consultadas,6records=4avisos distintos (3security+1maintenance). FAIL675 mantiene cuarentena. Safetensors sourceb7c0f38 corrige pyo3/memmap2,45registry crates0OSV sólo metadata; build pendiente.43/48sin cambio. Ver reconstruction_evidence/TRAINING_WHEEL_QUARANTINE_V371.md.

V372 maintenance: reference telemetry implementation is returned to canonical worker0.1.5 and WINDOWS_REFERENCE_TELEMETRY_RUNTIME.md/0.1.0, with opt-in WINDOWS_REFERENCE_TELEMETRY_PACK_PLAN.md. Source-only acquirer0.4.82 remains separate from runtime admission. TEST05 closure requires the fresh integrated replay and gates in reconstruction_evidence/INTEGRATED_TELEMETRY_CONTROL_V372.md; no production acceptance or private training data inferred.

V373: opt-in HISTORY-MODEL-TRAINING-PIPELINE0.1.0/17files, HISTORY_MODEL_TRAINING_PACK_PLAN2packs/21files;26policy tests and actual earlier synthetic SFT/evaluation/rollback. Final exact reconstruction/bootstrap/E2E and TEST09 decision are recorded in reconstruction_evidence/HISTORY_MODEL_TRAINING_CONTROL_V373.md. NEW requires no history; EXISTING keeps data/model private with its own authority/quality/runtime gates. No automatic training, provider call or deployment. SDK component selection precedes AUTHORED glue; see HISTORY_TRAINING_SDK_QUALIFICATION_V373.md and HISTORY_TRAINING_SDK_GAP_V373.md. Integral franchise remains67/754.

V373 cierre201: TEST09 PASS por pack17files/composición21,27policy, bootstrap58wheels/23713file hashes, SFT real, evaluación independiente,13negativos, timeout/no replay, decisiones exactas y rollback local con baseline previo. Preflight200160pasos ejecutados PASS/55planes;164/1506/818/55, integral67/754.45/48, tres pendientes TEST02/03/07; no cambia oráculos ni hereda corpus/modelo/seguridad/deploy del consumer. Ver HISTORY_MODEL_TRAINING_CONTROL_V373.md.

V374: auditoría por claim21núcleos,228tests/306semillas/23campañas (39811771ejecuciones),6comparaciones nuevas de fuentes,42archivos reconstruidos,20packs/40archivos compuestos y candidato rechazado.20CONDITIONED/1CANDIDATE intactos; cero equivalencias empresariales/promociones inventadas.13negativos de enlace de evidencia integrados en VERIFY_LIBRARY. Preflight/cierre del control aún pendientes;45/48sin cambio. Ver reconstruction_evidence/CORE_CLAIM_ADMISSION_V374.md; la integración de negocio más amplia de FAIL385 mantiene sus owners.

V374 cierre203: TEST02 PASS como auditoría por claim21núcleos;228tests/306semillas,23targets finales (6389808ejecuciones),20packs/40files compuestos, candidato rechazado,13negativos de evidencia y Preflight202160pasos PASS.46/48controles;20CONDITIONED/1CANDIDATE y FAIL385/integración más amplia conservados.164/1506/819/55, integral67/754. No cambia48oráculos ni promueve producto. Ver CORE_CLAIM_ADMISSION_V374.md; continuar TEST03 y luego TEST07 según gates materiales.

V374 rectificación204 / FAIL710: el cierre203 de TEST02 se retracta. Fuente/tests/composición son válidos, pero FAIL385 aún exige equivalencia/integración material.45/48PASS; TEST02/03/07BLOCKED. Se preservan historia y mejoras; no se cambia el oracle ni se reduce el roadmap. Ver CORE_CLAIM_ADMISSION_V374.md, encabezado vigente.

V374 successor205: SECURE-OPS1.1.4 corrects the error-rate denominator that hid low-traffic outages. Canonical Prometheus3.14.0 rules9scenarios/13assertions and5wrapper rejections PASS; exact pin/no-engine distinction retained. Three consumer plans now select1.1.4.164packs/1507files/819Markdown/55plans;1260AUTHORED/140ADAPTED/107VERBATIM; integral67/755. No dependency acquisition, core promotion or whole-control closure;45/48 unchanged. Evidence: reconstruction_evidence/CORE_CLAIM_ADMISSION_V374.md. Full successor Preflight pending.

V374 checkpoint207: fullPreflight206 passed162executed steps and all55compositions, including exact Prometheus engine regression. SECUREOPS1.1.4 correction complete; integral67/755, library164/1507/819/55. Current45/48unchanged; TEST02/03/07 remain blocked. Retain the203retraction and actual FAIL385 integration work. See reconstruction_evidence/CORE_CLAIM_ADMISSION_V374.md.

V394: PNPM-ARTIFACT-SELECTION-GATE0.3.0/6files corrige FAIL782 (pins empresariales y Playwright obsoletos), entrega BLUEOAK-NOTICE.md ligado a5declaraciones exactas y verifica3recetas actuales;51tests y10mutaciones reales rechazadas. Payload442/22notices intactos; sin ejecución ni admisión runtime/redistribución. Evidencia: reconstruction_evidence/PNPM_CURRENT_ROUTING_NOTICES_V394.md. Los pendientes BlueOak se reducen a entrega downstream/revisión restante; ausencia de LICENSE original chownr no se falsea. Otros permisos/source/publicación de pnpm continúan abiertos.

V395: PNPM-ARTIFACT-SELECTION-GATE0.4.0/6files entrega QRCODE-NOTICE.md con encabezado de autor/modificación exacto y texto MIT completo, además del aviso BlueOak intacto.59tests PASS,3recetas reales y12negativos rechazados;10source blobs/10regiones de bundle fijados por separado, sin claim de equivalencia completa. FAIL783 entrega local resuelto, FAIL784 metadata ADAPTED reconstruida;payload442/22notices intactos. No ejecución ni admisión runtime/redistribución;45/48. Ver reconstruction_evidence/QRCODE_VENDOR_NOTICE_DELIVERY_V395.md.

V396: cerrado descubrimiento y entrega local de licencia original semver-utils1.1.4: MIT OR Apache-2.0 explícito, opción MIT completa y1839bytes originales retenidos. Core0.4.88/51files/196sources/25profiles añade sólo cuarentena npm exact4193bytes/SHA512;140checks y gates anteriores PASS. Planner0.5.0/6files/67tests;3recetas reales/12negativos PASS, BlueOak/QRCode y442payloadfiles intactos. Clave registry vencida documentada; no firma vigente, equivalencia build, ejecución de recetas ni admisión global pnpm.45/48sin promoción; V386/Daybreak sigue diferido. Ver reconstruction_evidence/SEMVER_ORIGINAL_LICENSE_DELIVERY_V396.md.

V397: PNPM-ARTIFACT-SELECTION-GATE0.6.0/7files entrega el conjunto completo de47textos retenidos/151253bytes en3recetas:141copias verificadas independientemente,12negativos reales y82tests PASS.22notices originales comparados contra payload y25evidencias conservan alcance; BlueOak/QRCode/semver y442payloadfiles intactos. Catálogo ADAPTED con términos por texto, sin algoritmos ni promoción de licencias/source/relinking/runtime.45/48; Daybreak diferido. Ver reconstruction_evidence/PNPM_RETAINED_NOTICE_DELIVERY_V397.md.

V398: PNPM-ARTIFACT-SELECTION-GATE0.7.0/8files entrega fuente original next-path1.0.0, manifiesto y MPL completa en3recetas:9copias exactas/12negativos reales/90tests PASS. Commit oficial y3Git blobs verificados;4sentencias comparadas bajo adaptadores explícitos,10mutaciones rechazadas. No equivalencia runtime ni reproducibilidad pnpm. Colección47, suplementos anteriores y442payloadfiles intactos.45/48; Daybreak diferido. Ver reconstruction_evidence/NEXT_PATH_MPL_SOURCE_DELIVERY_V398.md.

V399: continuidad de cierre conciliada con V372/V393 y21packs/42fuentes byte-idénticos aV374. No repetir esas suites sin delta; agrupar trabajo por bloqueo material.45/48 intacto; TEST02integración,TEST03grafo/tooling yTEST07release siguen abiertos. Observabilidad candidata0.3.3 separada del runtime probado; Sharp ausente del perfil web no reabre Daybreak. Ver reconstruction_evidence/CURRENT_CLOSURE_BOUNDARIES_V399.md.

V400: encuestas conectadas con PostgreSQL/OIDC/Next, recuperación GET,4navegadores/8respuestas/8POST y retención CLI1/1/0.23archivos nuevos AUTHORED, packs GO-CUSTOMER-SURVEY-API0.1.0 y TS-CUSTOMER-SURVEY-PORTAL0.1.0, CONDITIONED.45/48sin promoción; TEST02/03/07 siguen abiertos. Ver reconstruction_evidence/CONNECTED_CUSTOMER_SURVEYS_V400.md.

V400 cierre270: encuestas conectadas y retención transaccional incorporadas. Go/TSsurvey0.1.1, app1.9.4, HTTPref0.1.15;23archivos nuevos,828/828 reconstruidos.4navegadores/8respuestas/8POST,24concurrentes,restart,CLI1/1/0,fuzz6219;correcciónSQL2→0tablas parciales PASS. Preflight269164pasos/56perfiles PASS más deltaSQL focal;167/1587/850 y1330/150/107,franquicia70/820,HTTP71/828,web7/134.45/48sin promoción; Daybreak diferido. Ver reconstruction_evidence/CONNECTED_CUSTOMER_SURVEYS_V400.md.

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
