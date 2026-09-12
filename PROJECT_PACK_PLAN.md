# Plan de packs — preparación de mantenimiento

Estado V402: preparación local `READY_FOR_LIBRARY_WORK`; T2801 final y release de biblioteca pendientes. El modo `PROJECT` y la aprobación productiva permanecen separados.

| Orden | Pack | Versión | Destino / condiciones |
|---|---|---|---|
| 1 | PROJECT-START-READINESS-VALIDATOR | 0.7.0 | staging separado; once archivos reconstruidos, 102 tests; modos PROJECT y LIBRARY_INFRASTRUCTURE separados |
| 2 | EXECUTION-VALIDATOR | 1.3.1 | engineering_execution_kit en staging; mantener estado/eventos en raíz; completar implementation_assurance antes de expansión |
| 3 | MARKDOWN-COMPOSITOR | 0.2.2 | tooling en staging; selección estable bajo StrictMode, destino nuevo y rechazo de colisiones; V285 |

No se agregan estos outputs a la fuente portable ni se mezclan con los archivos de negocio.
La selección integral de producto continúa en markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md,
82 packs/1032 archivos al inicio V402, sujeta al inventario ejecutable vigente, al gate de infraestructura y a los gates propios de cada target.

Código canónico: paths y digests en PROJECT_EXTERNAL_SOURCE_LOCK.md.
Reproducción: materialize_markdown_pack.ps1 con PackFile exacto y Destination ausente;
no requiere Git ni un framework nuevo. El rollback de mantenimiento conserva el
snapshot y evidencia anterior; consultar specs/library-maintenance/plan.md.

V287 ya demuestra la entrada CLI NEW/EXISTING; V289 aporta el contrato de assurance
en nivel plan. Siguen pendientes assurance integral, comparación controlada de
eficiencia/tokens y artefacto final. No volver a contar esas pruebas existentes
como ausentes ni confundirlas con aprobación del perfil de producto.

## Selección ejecutable del tooling local

La ejecución usa los tres owners de la tabla: primero materializa el compositor
con materialize_markdown_pack.ps1; después este bloque compone los dos validators.
No incluir el compositor en su propio plan: sus tests contienen placeholders
literales que el gate de variables rechaza (FAIL-410). Acknowledgement permite
el ensayo local condicionado, no promueve los packs ni aprueba una franquicia.
Destino nuevo y separado del producto; ninguna descarga o credencial requerida.

```json
{
  "compositionVersion": "1.0",
  "secretVariables": {},
  "packs": [
    {
      "path": "implementation_packs/PROJECT_START_READINESS_VALIDATOR.md",
      "packId": "PROJECT-START-READINESS-VALIDATOR",
      "version": "0.7.0",
      "acknowledgeConditions": false,
      "files": ["*"],
      "variables": {}
    },
    {
      "path": "implementation_packs/ENGINEERING_EXECUTION_VALIDATOR.md",
      "packId": "EXECUTION-VALIDATOR",
      "version": "1.3.1",
      "acknowledgeConditions": true,
      "files": ["*"],
      "variables": {}
    }
  ]
}
```

Invocar el compositor materializado con -PlanFile PROJECT_PACK_PLAN.md,
-LibraryRoot apuntando a esta biblioteca y -Destination a un directorio ausente.
Conservar MATERIALIZATION_RECORD.md y comparar los archivos reconstruidos.
Rollback definido en specs/library-maintenance/plan.md: no sobrescribir destinos;
conservar snapshot/evidencia y restaurar sólo cambios propios verificados. La
colisión debe rechazar sin alterar bytes. El plan no autoriza despliegue ni PITR.

## V322 — governed maintenance transport

OFFICIAL_UPSTREAM_ACQUISITION_CORE0.4.78:124entries/17profiles,31files; initialization2packs/37files. Three opaque artifacts Node24.20.0 win-x64, pnpm11.25.0 and Sigstore4.1.1, exact digests and byte-preserving license envelopes. Local AUTHORED transport with ADAPTED base64 sidecars; acquisition only, no execution or production promotion. FAIL525 remains open until fresh materialized receipts. Evidence: reconstruction_evidence/GOVERNED_MAINTENANCE_ACQUISITION_V322.md.

## V322 — adquisición de mantenimiento comprobada

Core0.4.78 / initialization37files; narrow G0–G8 USE_REUSABLE_PACK con owner global CONDITIONED. Perfil materializado maintenance-runtime-artifacts, lock124 y3sources adquiridos por HTTPS con approval actual del agente bajo autorización previa; no aprobación retrospectiva. Node24.20.0 win-x64, pnpm11.25.0 y Sigstore4.1.1: bytes/licencias exactos y10outputs inmutables tras PRESENT. FAIL525/527/528/529/530/531 cerrados focalmente, TEST41/EVID32. Biblioteca161packs/1453files/757Markdown/52profiles y151preflight steps PASS; preflight global BLOCKED por3tools no resueltos y readiness/target separados. Ver GOVERNED_MAINTENANCE_ACQUISITION_V322.md y PROJECT_OFFICIAL_SOURCE_PROFILE_RECORD.md. Continuar6binarios nativos de pnpm y T2809; no scheduler, ZIP, install global, ejecución de artefactos adquiridos ni producción.

V325 concilia el conteo vigente con el perfil canónico y Preflight151steps:67/746. Sólo corrige una cifra desactualizada; no aprueba composición productiva ni altera selección. Evidencia FIREFOX_LIFECYCLE_AND_CONTINUITY_V325.md.

### V334 — proyección pnpm canónica completada

PNPM-ARTIFACT-SELECTION-GATE0.1.0,4files/17tests dos veces y2proyecciones reales443files idénticas; G0–G8/USE_REUSABLE_PACK sólo tooling en este contexto. Perfil separado1/4 y2pasos de verificación integrados. No consumers/runtime/redistribución admitidos; condiciones explícitas y FAIL532 original conservado. Evidencia reconstruction_evidence/PNPM_SELECTION_GATE_V334.md.

### V335 — routing de consumidores pnpm fijados

PNPM-ARTIFACT-SELECTION-GATE0.2.0,6files/43tests fuente+rebuild. Prepare/verify produce recetas offline con configuración aislada, store/cache explícitos y perfiles exactos; nunca ejecuta ni admite runtime.3instalaciones desde rebuild PASS sin descargas/cambios de lock;115tests/1skip ybuild web PASS. FAIL555/557 corregidos sin omitir política; G0–G8 acotado. Perfil separado1/6; producto67/746 inalterado. Evidence PNPM_CONSUMER_ROUTING_V335.md.

### V337 — licencia exacta npm-lifecycle y transporte gobernado

OFFICIAL-UPSTREAM-ACQUISITION-CORE0.4.79:33files/125sources,78opaque checks/rebuild y18perfiles de adquisición. Nuevo perfil separado adquiere sólo npm-lifecycle1100.1.0; firma registry e integridad verificadas. LICENSE del tar exacto coincide con commit/sidecar,7de8archivos con Git blobs; único delta packageManager omitido.29planes de composición actualizados. Fuente licencia verificada, runtime/redistribución no admitidos. Ver reconstruction_evidence/PNPM_LIFECYCLE_LICENSE_V337.md.

V337 reference correction:30composition plans (29initial plus1compact JSON correction) now select core0.4.79. composition-reference-audit.json independently parses each composition object. Preflight rerun follows checkpoint112.

V338 notice coverage:473metadata hashes verified;468identities observed from456bundle identities+21external manifests+root (with overlap).22notice files and21embedded attribution paths inventoried;32exact draft texts preserved. BlueOak five external packages and QRCode vendored MIT require explicit notice treatment; nested works, semver-utils exact text and whole-bundle source/use obligations remain open. All442payload files unchanged; no runtime/redistribution admission. Evidence: reconstruction_evidence/PNPM_NOTICE_COVERAGE_V338.md.

V339:8fixed source manifests and7licence texts verified; four BlueOak source texts,3Noble locked build-input identities outside the original473graph,70node-gyp vendored files identical to fixed tree. Draft39preserves32previous texts. chownr declaration-only, semver-utils exact text and source/bundle/redistribution conditions remain open. All442selected files unchanged. See reconstruction_evidence/PNPM_NESTED_LICENSE_SOURCES_V339.md.

V340:103external Undici6files match fixed source; Undici7licences/headers bound separately. Yarn4.1.7 registry gitHead returns404 but official tag/source version exists at a distinct commit; preserve both, no published-byte equivalence. Root BSD and3MIT source comments retained; draft47preserves39previous texts,442selected files unchanged. FAIL575 scoped provenance remains unresolved. See reconstruction_evidence/PNPM_YARN_UNDICI_NOTICES_V340.md.

V341: five selected Undici7 AST units in two bundled modules match fixed-source syntax under declared identifier renames; nine initial controls plus eight actual-source mutations and two exact byte-range checks pass. No whole-module/runtime/redistribution claim; draft47/payload442 unchanged. Contract remains41passed/7pending of48 (approximately15% of rows, never global effort); ten macrofronts remain open. Evidence: reconstruction_evidence/PNPM_SCOPED_SOURCE_CORRESPONDENCE_V341.md.

V342 / T2809: GO-OBSERVABILITY-CORE0.3.1 corrects four false-success writer cases (FAIL581); two canonical files rebuilt exactly. Real HTTP/file reference32requests/64correlated records plus closed-file failure/replacement,28tests x3,vet/build and3fuzz targets10s PASS. Actual runtime host/reporter/export/retention/alerts and admission remain pending; isolated test integration is not production wiring. See reconstruction_evidence/OBSERVABILITY_DELIVERY_INTEGRATION_V342.md.

V343 in progress: concrete bounded JSONStatusReporter replaces the abstract-only reporting gap in the existing WhatsApp host; canonical0.13.0, four blocks rebuilt exactly,67/746composition. Local real-worker/PostgreSQL→TCP and cancellation/failure/concurrency tests pass x3. Final gates interrupted by computer shutdown are pending rerun, not PASS. Actual service mounting/supervisor/collector/retention/alerts remain T2809; no percentage. Evidence: reconstruction_evidence/WHATSAPP_STATUS_REPORTER_V343.md.

V343 / T2809: WhatsApp0.13.0 now supplies a concrete bounded JSON reporter to the existing status host. Real PostgreSQL→TCP and mid-record loss after commit/restart prove no job/observation/audit replay. Final14host tests x3,11seeds,10s fuzz298870executions,full Go baseline/vet/build and4/4canonical parity PASS;67/746unchanged. Reporter implementation gap closed locally; real service identity/mounting/supervisor/collector/retention/alerts remain. See reconstruction_evidence/WHATSAPP_STATUS_REPORTER_V343.md.

V344 / T2809: refund host0.1.4 now stops before another token/step on failed/incomplete reporting; fixed errors/exit1 and healthy/idle policy preserved. Six baseline red modes →green; actual loop/closed-file/subprocess,17tests x3,14seeds,fuzz327445,vet/build and15file parity PASS. Three changed files,12unchanged including business/provider/DB/locks;67/746consumer. One opt-in PG test skipped, no live/refund-policy change or supervisor activation. See reconstruction_evidence/REFUND_HOST_REPORT_DELIVERY_V344.md.

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
