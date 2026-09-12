# Mapa de autoridades — mantenimiento local

Estado: seleccionado para el delta T2801; no aprobación general de todos los upstreams.
Usar AI_ENGINEERING_MASTER_MAP.md y SYSTEMS_ENGINEERING_MASTER_MAP.md como índices.

| Decisión | Owner interno | Autoridad primaria / claim limitado |
|---|---|---|
| continuidad, intake y contexto mínimo | AGENT_SYSTEM_START.md; CODEX_ELITE_PROJECT_BOOTSTRAP.md | protocolos vigentes del workspace, sin releer todo el corpus |
| fuentes, licencias y calidad | PUBLIC_CODE_ARCHITECTURE_ADMISSION_STANDARD.md; PUBLIC_CODE_ARCHITECTURE_MASTER_MAP.md | G0–G8; ninguna marca empresarial reemplaza evidencia |
| especificación y ejecución | ENGINEERING_EXECUTION_PLAYBOOK.md | GitHub Spec Kit, spec → plan → tasks, no garante de producción |
| interfaz CLI y conexión completa | markdown_system/TOTAL_SYSTEM_CAPABILITY_CONTRACT.md | admite CLI/AUTOMATION con roles, efectos, pruebas y recovery; V283 corrige contradicción del gate |
| lanzamiento y dependencias | SECURITY_SRE_CLOUD_INFRASTRUCTURE.md; markdown_system/DEPENDENCY_UPDATE_CONTRACT.md | Google SRE: controles prácticos y proporcionales al lanzamiento |
| verificación frente a validación | ENGINEERING_EXECUTION_PLAYBOOK.md | NASA Handbook: evidencias distintas; un check estructural no sustituye aceptación |

Fuentes primarias consultadas el 2026-09-07:
[GitHub Spec Kit](https://github.github.com/spec-kit/),
[Google SRE](https://sre.google/sre-book/reliable-product-launches/),
[NASA Appendix](https://www.nasa.gov/reference/system-engineering-handbook-appendix/).
Estas fuentes de método no convierten el código local en código publicado por esas empresas.
Las revisiones ya fijadas en cada pack se conservan; no se actualizaron por una página móvil.

Actualización V290: para composición del tooling, la autoridad material es
markdown_system/COMPOSITION_PROTOCOL.md; bootstrap del compositor separado y
plan ejecutable de validators probado en TOOLING_COMPOSITION_CLOSURE_V290.md.
Google SRE Release Engineering y Reliable Product Launches se consultaron otra
vez el 2026-09-07; sus principios no convierten estas herramientas AUTHORED en
código Google. Pack plan probado no significa authorities integral PROVEN.

No cargados para este cambio: manuales fiscales, modelos OCR, GPU, ads y firmware;
no cambian la corrección de un validador de journeys CLI. Siguen obligatorios cuando
se trabaje en sus capabilities. Freshness integral y assurance permanecen pendientes.

## V322 — governed maintenance transport

OFFICIAL_UPSTREAM_ACQUISITION_CORE0.4.78:124entries/17profiles,31files; initialization2packs/37files. Three opaque artifacts Node24.20.0 win-x64, pnpm11.25.0 and Sigstore4.1.1, exact digests and byte-preserving license envelopes. Local AUTHORED transport with ADAPTED base64 sidecars; acquisition only, no execution or production promotion. FAIL525 remains open until fresh materialized receipts. Evidence: reconstruction_evidence/GOVERNED_MAINTENANCE_ACQUISITION_V322.md.

## V322 — adquisición de mantenimiento comprobada

Core0.4.78 / initialization37files; narrow G0–G8 USE_REUSABLE_PACK con owner global CONDITIONED. Perfil materializado maintenance-runtime-artifacts, lock124 y3sources adquiridos por HTTPS con approval actual del agente bajo autorización previa; no aprobación retrospectiva. Node24.20.0 win-x64, pnpm11.25.0 y Sigstore4.1.1: bytes/licencias exactos y10outputs inmutables tras PRESENT. FAIL525/527/528/529/530/531 cerrados focalmente, TEST41/EVID32. Biblioteca161packs/1453files/757Markdown/52profiles y151preflight steps PASS; preflight global BLOCKED por3tools no resueltos y readiness/target separados. Ver GOVERNED_MAINTENANCE_ACQUISITION_V322.md y PROJECT_OFFICIAL_SOURCE_PROFILE_RECORD.md. Continuar6binarios nativos de pnpm y T2809; no scheduler, ZIP, install global, ejecución de artefactos adquiridos ni producción.

## V327 — entrada CLI y diagnóstico

FRONTEND_PRODUCT_ENGINEERING_UX.md §2 aporta visibility/recovery/progressive disclosure para esta guía local; COMPOSITION_PROTOCOL y EXECUTION-VALIDATOR1.3.1 son los owners ejecutables. Prueba automatizada de comandos y errores, no estudio de usuarios ni generación automática del authority map. No se cargaron autoridades fiscales/documentales/proveedores para corregir esta frontera de tooling.

## V329 — selección manual para el gap de automatización

| Decisión / requisito | Autoridad seleccionada | Motivo y frontera | Evidencia / estado |
|---|---|---|---|
| Entrada y contexto compacto / LIB-R01,R08 | CODEX_ELITE_PROJECT_BOOTSTRAP §5; ENGINEERING_EXECUTION_PLAYBOOK §3.1 | una autoridad por decisión, secciones/justificación y delta NEW/EXISTING; índice no es selección | AUTHORITY_MAP_GAP_V329.md; diseño documentado, generador ausente |
| Adopción / LIB-R02 | PUBLIC_CODE_ARCHITECTURE_ADMISSION_STANDARD G0–G8 | método oficial no equivale a código admitido ni implementación local | expediente authority-gap-v329; NO_ADMISSIBLE_SOURCE acotado |
| Recuperación / LIB-R06 | AGENT_ERROR_RECOVERY_PROTOCOL; FAILURE_LEARNING_CONTRACT | conservar originales, registrar fallos y verificar cadena | FAIL473 recurrencia de diagnóstico; checkpoint anterior conservado |
| Vigencia / LIB-R07 | AUTHORITY_FRESHNESS_AND_SELF_CORRECTION_CONTRACT | diferenciar docs actuales y revisión fijada; hash de nota no es hash upstream | tres notas con URLs/revisiones/licencias; sin upgrade |

AI_ENGINEERING_MASTER_MAP y SYSTEMS_ENGINEERING_MASTER_MAP se consultan como
índices. Manuales fiscales, OCR, GPU y ads evaluados pero no cargados: no afectan
la selección de autoridades del tooling local; sus capabilities no se dan por
cerradas. Incertidumbre: falta candidato compatible y evaluación de pertinencia.
No nuevas dependencias/licencias redistribuidas. Esta tabla es manual y limitada
al delta V329, no salida de un generador ni cierre de assurance integral.

## V330 — verificación multi-tab de invariantes existentes

| Decisión | Autoridad material | Aplicación y límite |
|---|---|---|
| UI ante concurrencia y fallo | FRONTEND_PRODUCT_ENGINEERING_UX §2.2/4.1 | referencias, bloqueo visible y recuperación; no mutex cross-tab ni nueva regla |
| Oráculo conectado | ENGINEERING_EXECUTION_PLAYBOOK §11.1–11.3; V312 | comparar intención/durable/actores y contar efectos; compilación no demuestra el resultado |
| Procedencia del delta | PUBLIC_CODE_ARCHITECTURE_ADMISSION_STANDARD | dos tests AUTHORED de packs existentes; no nuevo upstream ni promoción |
| Recuperación de pruebas | AGENT_ERROR_RECOVERY_PROTOCOL; FAILURE_LEARNING_CONTRACT | FAIL542 preservado y selectores corregidos en fuente canónica/reconstrucción |

Mapas AI/SYSTEMS usados como índices. MDN/Playwright documentan aislamiento,
opener y contexto de pruebas; no se adquiere código. Fiscal/documentos/ERP
no intervienen en este delta de verificación; sus gates permanecen pendientes.

## V331 — autorización durante cambio de sesión

SECURITY_SRE_CLOUD_INFRASTRUCTURE §4.4 gobierna subject/action/resource/tenant y negativos de permisos; ENGINEERING_EXECUTION_PLAYBOOK §11.1–11.3 exige oráculo exacto y corrección del fixture. FRONTEND_PRODUCT_ENGINEERING_UX/V312 gobiernan referencia pendiente y recuperación. Mapas AI/SYSTEMS son índices; sólo dos tests AUTHORED existentes, sin código upstream nuevo. No claims de revocación IdP o borrado de datos ya entregados.

## V332 — empaquetado y recuperación de copia portable

ENGINEERING_EXECUTION_PLAYBOOK §15 exige artifact identity/rollback. PORTABLE_SIGNED_RELEASE_EVIDENCE_GATE0.1.0 aporta patrón local de ZIP determinista, sin heredar firma ni admisión. AGENT_ERROR_RECOVERY_PROTOCOL y FAILURE_LEARNING_CONTRACT gobiernan FAIL544/545, candidatos preservados y reconstrucción. Blueprint/roadmap mantienen capacidades genéricas y condiciones de proyectos futuros; no se inventa su cierre.

## V334 — proyección de artefacto adquirido

COMPOSITION_PROTOCOL/DEPENDENCY_UPDATE_CONTRACT/admission standard y execution playbook15 gobiernan tooling AUTHORED; pnpm11.25 release y Python tarfile son fuentes primarias inspeccionadas. Ninguno aporta la convención local455miembros/receipt. G0–G8 y evidencia por nueve dimensiones en PNPM_SELECTION_GATE_V334.md; no heredar runtime ni licencia de redistribución.

V335 mantiene autoridades V334; inspección del bundle11.25 fijado y probes reales gobiernan dialecto CLI/cache/política, sin importar defaults12.x. AUTHORED routing no equivale a runtime admisión. G0–G8/nueve dimensiones en PNPM_CONSUMER_ROUTING_V335.md.

## V349 — publicación de resultados y recuperación local

DATABASE_STORAGE_INTERNALS.md distingue escritura, flush, publicación, confirmación y recuperación; SECURITY_SRE_CLOUD_INFRASTRUCTURE.md gobierna privacidad/operación. APIs oficiales MoveFileExW/FlushFileBuffers y CPython os.fsync explican el mecanismo limitado. Los tests prueban caída del proceso y reconciliación sin repetir el intento; no corte eléctrico, firmware, WORM, réplica ni éxito de negocio. Evidencia: reconstruction_evidence/NATIVE_RESULT_PUBLICATION_V349.md.

V350: SECURITY_SRE_CLOUD_INFRASTRUCTURE.md y DATABASE_STORAGE_INTERNALS.md gobiernan apagado/commit/acknowledgement. WaitForSingleObject/DuplicateHandle oficiales y Go context consultados2026-09-09; puente AUTHORED, no autenticación del tipo de handle ni confirmación del negocio. Evidencia reconstruction_evidence/NATIVE_WORKER_LIFECYCLE_V350.md.

V351: DATABASE_STORAGE_INTERNALS.md gobierna commit/ack/restart y SECURITY_SRE_CLOUD_INFRASTRUCTURE.md la frontera nativa. PostgreSQL18 runtime-config-wal e initdb oficiales consultados2026-09-09. Durabilidad observada sólo en reinicio local controlado; efectos externos simulados. Evidencia reconstruction_evidence/NATIVE_POSTGRES_LIFECYCLE_V351.md.

V352: ENGINEERING_EXECUTION_PLAYBOOK.md y TOOLCHAINS_BUILDS_PACKAGING_FFI.md gobiernan source→artifact→acceptance→publish. Se corrige la dependencia procedimental final-first, conservando criterios TEST06/08 y misma identidad/revalidación para TEST07; fuentes AUTHORED y tooling de copia148 verificados por manifest. No fuente externa nueva ni promoción de producto.

V353: Go specification floating-point types y pkg.go.dev/math gobiernan semántica numérica, no admisión financiera. FX AUTHORED0.1.1 usa stdlib existente y oráculo racional local. Autoridad de aceptación sigue el contrato original y candidato148; el delta FX obliga nueva identidad/revalidación TEST07.

V354: algoritmos/invariantes y Backend gobiernan límites numéricos, claves de identidad y rechazo sin efecto. Go spec Integer_overflow/Map_types consultada2026-09-09; G0/G1 son AUTHORED/LicenseRef-Workspace-Owner, no upstream financiero. V293 sigue owner de equivalencias. Ver CREDIT_CORE_ISOLATION_V354 para nueve dimensiones y límites.

V355: Backend ADT/invariantes y ejecución/evidencia gobiernan tenant explícito y migración. No autoridad comercial ni código upstream nuevo. Recordatorio local AUTHORED no equivale a scheduler durable ni exactly-once; reutilizar owners V293. Supersede claim tenant-scoped de API0.1 y compatibilidad waitlist, preservando fuente y evidencia históricas.

V356: Algorithms/Data Structures y Backend gobiernan invariante numérico y error sin efecto; ejecución/fuzz/admisión gobiernan evidencia. Descomposición aritmética AUTHORED; oráculo math/big stdlib existente. No fuente externa adquirida ni autoridad fiscal/comercial nueva; owners V293 conservados.

V357: Algorithms/Data Structures y Backend/API gobiernan igualdad estructural, invariantes y errores; playbook/fuzz/failure contract gobiernan prueba. Referidos AUTHORED, sin fuente externa adquirida; V293 conserva CRM/lead como owner y brecha de recompensa durable.

V358: Algorithms/Data Structures y Backend/API gobiernan igualdad tuple, moderación y FIFO; ejecución/fuzz/failure contract gobiernan pruebas. Cores AUTHORED, sin fuente externa adquirida o nueva política. Gaps/owners V293 se mantienen.

V359: Algorithms/Data Structures y Backend/API gobiernan identidad tuple, checklist y máquina de estados; ejecución/fuzz/failure contract gobiernan prueba. AUTHORED, sin nueva adquisición ni autoridad legal. Owners V293 conservados.

V360: Algorithms/Data Structures y Backend/API gobiernan igualdad, ownership y alcance de caller; ejecución/fuzz/failure contract gobiernan prueba. Cores AUTHORED, sin adquisición/autoridad externa nueva ni política de marketing inferida. V293 conserva owners lead/ads/notificación.

V361: Algorithms/Data Structures y Backend/API gobiernan aritmética exacta/errores/ownership; ejecución/fuzz/failure contract gobiernan prueba. Docker se explica desde código del verifier y receipt nativo V351 SHA-verificado, sin nueva adquisición o política financiera/SLO.

V362: Backend/API scope e invariantes y Algorithms/Data Structures gobiernan validación de mes, tuplas y límite de versión. Execution/fuzz/failure contracts gobiernan prueba. No fuente externa nueva: Search se verifica como substring en código local; KPI no verifica pertenencia de Input. No se carga dominio fiscal/publicación para inventar políticas.

V363: Backend/API y Algorithms gobiernan scope/validación/modelos; Go fuzz/execution/failure contracts gobiernan pruebas. Google Search title-link/snippet y Unicode CLDR plural-rules actuales refutan atribución60/160 y equivalencia binaria. Se conserva política local AUTHORED, sin motor/provider nuevo. V293 gobierna no duplicar owners y límites de equivalencia; observability sigueCANDIDATE.

## V372 — observabilidad integrada del target de referencia

SECURITY_SRE_CLOUD_INFRASTRUCTURE.md y ENGINEERING_EXECUTION_PLAYBOOK.md gobiernan señales operables, fallos, evidencia y límites; PUBLIC_CODE_ARCHITECTURE_ADMISSION_STANDARD.md G0–G8 gobierna fuentes/licencias/adaptaciones. OTLP/Collector/Prometheus y Win32 son autoridades contractuales, no autores del runner/reporter. El registro V372 conserva búsqueda oficial, versiones exactas, pruebas reales y gates estrechos; el alcance es referencia local sin proveedores.

V373: opt-in HISTORY-MODEL-TRAINING-PIPELINE0.1.0/17files, HISTORY_MODEL_TRAINING_PACK_PLAN2packs/21files;26policy tests and actual earlier synthetic SFT/evaluation/rollback. Final exact reconstruction/bootstrap/E2E and TEST09 decision are recorded in reconstruction_evidence/HISTORY_MODEL_TRAINING_CONTROL_V373.md. NEW requires no history; EXISTING keeps data/model private with its own authority/quality/runtime gates. No automatic training, provider call or deployment. SDK component selection precedes AUTHORED glue; see HISTORY_TRAINING_SDK_QUALIFICATION_V373.md and HISTORY_TRAINING_SDK_GAP_V373.md. Integral franchise remains67/754.

V373 cierre201: TEST09 PASS por pack17files/composición21,27policy, bootstrap58wheels/23713file hashes, SFT real, evaluación independiente,13negativos, timeout/no replay, decisiones exactas y rollback local con baseline previo. Preflight200160pasos ejecutados PASS/55planes;164/1506/818/55, integral67/754.45/48, tres pendientes TEST02/03/07; no cambia oráculos ni hereda corpus/modelo/seguridad/deploy del consumer. Ver HISTORY_MODEL_TRAINING_CONTROL_V373.md.

V374: auditoría por claim21núcleos,228tests/306semillas/23campañas (39811771ejecuciones),6comparaciones nuevas de fuentes,42archivos reconstruidos,20packs/40archivos compuestos y candidato rechazado.20CONDITIONED/1CANDIDATE intactos; cero equivalencias empresariales/promociones inventadas.13negativos de enlace de evidencia integrados en VERIFY_LIBRARY. Preflight/cierre del control aún pendientes;45/48sin cambio. Ver reconstruction_evidence/CORE_CLAIM_ADMISSION_V374.md; la integración de negocio más amplia de FAIL385 mantiene sus owners.

V374 cierre203: TEST02 PASS como auditoría por claim21núcleos;228tests/306semillas,23targets finales (6389808ejecuciones),20packs/40files compuestos, candidato rechazado,13negativos de evidencia y Preflight202160pasos PASS.46/48controles;20CONDITIONED/1CANDIDATE y FAIL385/integración más amplia conservados.164/1506/819/55, integral67/754. No cambia48oráculos ni promueve producto. Ver CORE_CLAIM_ADMISSION_V374.md; continuar TEST03 y luego TEST07 según gates materiales.

V374 rectificación204 / FAIL710: el cierre203 de TEST02 se retracta. Fuente/tests/composición son válidos, pero FAIL385 aún exige equivalencia/integración material.45/48PASS; TEST02/03/07BLOCKED. Se preservan historia y mejoras; no se cambia el oracle ni se reduce el roadmap. Ver CORE_CLAIM_ADMISSION_V374.md, encabezado vigente.

V374 successor205: SECURE-OPS1.1.4 corrects the error-rate denominator that hid low-traffic outages. Canonical Prometheus3.14.0 rules9scenarios/13assertions and5wrapper rejections PASS; exact pin/no-engine distinction retained. Three consumer plans now select1.1.4.164packs/1507files/819Markdown/55plans;1260AUTHORED/140ADAPTED/107VERBATIM; integral67/755. No dependency acquisition, core promotion or whole-control closure;45/48 unchanged. Evidence: reconstruction_evidence/CORE_CLAIM_ADMISSION_V374.md. Full successor Preflight pending.

V374 checkpoint207: fullPreflight206 passed162executed steps and all55compositions, including exact Prometheus engine regression. SECUREOPS1.1.4 correction complete; integral67/755, library164/1507/819/55. Current45/48unchanged; TEST02/03/07 remain blocked. Retain the203retraction and actual FAIL385 integration work. See reconstruction_evidence/CORE_CLAIM_ADMISSION_V374.md.

V375: HTTP_SLI_ALERT_INTEGRATION_V375.md proves real order HTTP/PostgreSQL → official closed metrics → unchanged canonical alert → repair/replay, and a byte-identical full-profile rebuild. GO-HTTP-METRICS-REFERENCE0.1.0/8files, opt-in68packs/763files; ordinary67/755 unchanged. 380official tests and2authored tests PASS;62public versions/0OSV findings in this reference graph,30compiled upstream license inventories. Trusted synthetic reference only; TEST02/03/07 and FAIL385 remain open,45/48.

V375 checkpoint212: fullPreflight212164executed steps PASS,56profiles composed;165/1515/823/56 and1268AUTHORED/140ADAPTED/107VERBATIM. HTTP reference508requests,four actual failures,canonical ten-minute alert and recovery,one durable order/idempotency/outbox,exact rebuilt executable;380official+2local tests PASS. Only Docker unavailable in that receipt.45/48 unchanged; TEST02/03/07 remain blocked. Next action: FAIL732 official Next/Sharp patched-release qualification, preserving TEST02 owner integration work.

V387: existing keyset API owns after/next_cursor. Correct missing UI traversal using independent per-list GET cursors, fixed size/session scope and full27row boundary fixtures. UNIQUE NULLS NOT DISTINCT requires unique fixture VIN/battery values. Reconstruct canonical source and install offline instead of copying pnpm trees across Windows long paths. Reuse unchanged artifact identity; do not resume V386 deferred security research. Evidence reconstruction_evidence/PRIVATE_PORTAL_PAGINATION_V387.md.

V388: seis listados cliente/fabrica/admin paginados con APIs existentes;163registros de lista por navegador,4navegadores PASS,0escrituras/snapshot6tablas intacto.276tests/1skip heredado,typecheck/build;6packs reconstruidos,805files exactos,68/797integral.45/48;TEST02integral pendiente,TEST03libxml2 ACCESS_BLOCKED,TEST07dependiente.
FAIL770 REGRESSION_PROVEN in reconstruction_evidence/ADMIN_PORTAL_PAGINATION_V388.md; FAIL771 full Preflight rerun pending. V386 native research remains ACCESS_BLOCKED; no promotion.

V389 checkpoint241: horas de turnos cliente coherentes entre cuenta/gestion con locale/zona del negocio, labels SSR preservados al hidratar.282tests/1skip previo,typecheck/build;2configuraciones x4navegadores,2instantes con limite de dia/DST,snapshot7tablas intacto y0escrituras.6packs reconstruidos/805files,HTTPbinario idéntico;68/797.45/48; TEST02integral pendiente, V386native diferido,TEST07dependiente. Preflight241 pendiente.
FAIL772/773/774 REGRESSION_PROVEN in reconstruction_evidence/CUSTOMER_APPOINTMENT_TIMEZONE_V389.md. No native/security or whole control promotion.

V390 checkpoint244: cancelacion cliente conectada y recuperable: receipt ligado aID/org/estado/version,guard sincrono,GETexplicito tras resultado incierto.294tests/1skip previo,typecheck/build;4navegadores,16cancelaciones/16audits/16outbox con28POSTincluidos rechazos;8regresiones lectura/fechas sin escrituras.6packs/805files reconstruidos,8589artefactos yHTTPbinario idénticos;68/797.45/48;V386diferido,TEST02integral yTEST07pendientes. Preflight244 pendiente.
FAIL776/777/778 REGRESSION_PROVEN in reconstruction_evidence/CUSTOMER_CANCELLATION_RECOVERY_V390.md. No native/security or whole control promotion.
