# Franchise Pre-Flight Gap — qué falta antes de arrancar

Documento honesto de cierre: qué está materializado, qué código reusable aún
falta y qué sólo puede demostrarse en el proyecto target antes de promover una
franquicia. Inventario actualizado V402335, 2026-09-13. Inventario raíz: 212 packs, 2.578 archivos materializables, 1002 Markdown y 56 perfiles. El cierre V402335 conserva su inventario histórico. Audit V258: 160 packs,
121 fuentes oficiales, un resolver de capabilities ausentes, DevSkim adaptado,
OpenGrep/GitLab preservado pero bloqueado por el SCA actual de Cosign, release
portable OSV/SPDX/in-toto/SLSA/OpenSSH, assurance de implementación y 16 adapters. Estos PASS no
sustituyen los gates live.

El perfil integral vigente selecciona 1653 archivos desde 116 packs sin duplicar
paths; ese conteo también queda comprobado por `VERIFY_LIBRARY.ps1`.

V306: emisión de cotización con clave retenida, consulta scoped y actor atómico.
68 fases/4 proyectos, agenda4, PostgreSQL3, Go/web, rebuild11/11 y snapshot
idéntico tras reinicio PASS. QUOTE_CREATION_RECOVERY_V306.md. FAIL457 sigue OPEN
para creación de intervalos/recursos/checklist/recepción/disposición; sin promoción.

V305: cancelación de disponibilidad recuperable y coherencia con cupos/reservas.
64 fases/4 proyectos más agenda4, PostgreSQL3/focales10, Go/web y rebuild6/6 PASS.
Ver AVAILABILITY_CANCELLATION_RECOVERY_V305.md. FAIL457 sigue OPEN; siguiente
emisión de cotización. Sin nuevos módulos ni dependencias; no cierre productivo.

V304: lead asignación/transición con recuperación y actor autenticado durable.
60 fases/4 proyectos más agenda4, PostgreSQL3, Go/web y rebuild8/8 PASS.
Ver LEAD_COMMAND_RECOVERY_V304.md. FAIL457 no cierra para los demás comandos;
sin promoción ni ZIP ni habilitación externa. No se estima porcentaje global.

V303: resolución de discrepancias con tres decisiones, fence/marker, recibo
contrastado y GET; pérdida/JSON/consulta/storage/concurrencia/permisos probados.
56 fases/4 proyectos y agenda4 PASS, rebuild4/4. No creador inicial ni efectos
downstream concluidos; el resto de FAIL457 sigue OPEN. Evidencia
DELIVERY_EXCEPTION_RECOVERY_V303.md. No se promueve producto ni ZIP.

V302: portal operativo respeta permisos por sección y conserva pedidos ante
fallo503 de discrepancias; recuperación GET, ayuda y cero comandos automáticos.
52 fases/4 proyectos y agenda4 PASS, rebuild4/4; sin nuevos archivos de producto.
Ver OPERATOR_SECTIONS_RECOVERY_V302.md. FAIL457 (respuesta incierta del sender
operativo) permanece abierto y es el próximo cierre independiente T2804.

V300: lecturas relacionales fail-closed y recuperación de recepción/discrepancia
tras respuesta perdida, con48 fases/4 proyectos browser y PG3 veces PASS.
Rebuild7/7, sin archivos nuevos de producto. Ver DELIVERY_READ_RECOVERY_V300.md.
No confundir fixtures de preparación con el comando inicial aún ausente ni
aceptación local con liberación financiera, cobro o operación completa.

V299: defensa relacional del handover en Journey0.10.5, sin nuevos archivos de
producto ni migraciones. No confundir aceptación/checklist ya existentes con
preparación inicial, pago verificado o cierre completo de entrega. La regla
comercial de liberación sigue sin definir; se pidió esa decisión al usuario.
Ver DELIVERY_RELATIONAL_SCOPE_V299.md y el mismo roadmap T2802/T2803/T2804.

V296 stock por UI/BFF/HTTP/PG probado; V297 recibo idempotente y referencias NULL;
V298 añade solicitud de pago UI/BFF/HTTP/OIDC/PG, clave estable, exclusión de otro
intento inicial en ambas rutas y allowlist/configuración ausente cerrada.36 fases
en4 navegadores PASS con ayuda, permisos, concurrencia y recuperación. Aún faltan
elección del usuario para target, adapter/cuenta/sandbox/callback/reconciliación,
política financiera y preparación inicial de entrega. Ver PAYMENT_REQUEST_PORTAL_V298.md.

V295: guía de credenciales vinculada a código y custodia local Windows; ARCA
diferida explícitamente, API1.9.1 sin módulo fiscal por defecto. Recetas de target/
IdP no elegidos, accesos oficiales incompletos y adapters pendientes no se cierran
con una credencial. Evidencia: CREDENTIAL_PREPARATION_V295.md. Continúa el mismo
roadmap operativo de V294, sin sustituirlo por documentación.

V280: 21 packs citados como cobertura no están seleccionados en ese perfil;
comparar con owners alternativos, no incorporarlos indiscriminadamente.
GO-OBSERVABILITY-CORE 0.3.1 sigue CANDIDATE: contador V291, rechazo cerrado de tres exposiciones V292 e histograma V301 probados; política target, race, integración, carga y operación siguen pendientes. No compone en producto. Evidencia: HISTOGRAM_BOUNDARY_REPAIR_V301.md.
El cierre de las 48 superficies está desglosado en diez tareas verificables
en la sección 10 del roadmap, sin porcentaje de completitud inventado.

V275 fortalece la cola de jobs existente antes del worker automático: generación exacta, vigencia después del row lock y presupuesto por handler. Corrige cuatro fallos reales y reconstruye 4/4 archivos con tests PostgreSQL y vet/build; no suma módulos. El worker WhatsApp sigue pendiente, al igual que verificar fencing del publisher/outbox y recuperación operativa de intentos agotados antes de activar esos consumidores. Ver `reconstruction_evidence/JOB_CLAIM_GENERATION_AND_EXPIRY_V275.md`.

## Medición real V258 — sin mezclar biblioteca y producción

Actualización vigente V278: el host opt-in, historial autorizado, vista de agenda y ayuda de esa vista están implementados y probados localmente; configuración existente desactiva la función por defecto. Cuatro navegadores por cada uno de los dos recorridos conservan efectos únicos y recuperación. Quedan montaje/supervisor/identidad/reportes y alertas reales, retención, redrive, inbound/autorespuestas y los demás frentes del roadmap. Las menciones siguientes de host/UI pendientes son históricas: quedan superadas sólo para esta implementación local, no para su despliegue ni para todo el frontend. Ver `reconstruction_evidence/WHATSAPP_HOST_OPERATOR_HISTORY_V278.md`.

Actualización vigente V277: el worker de statuses ya está implementado sobre los owners existentes. Claim scoped y cierre job/inbox/audit atómico, retry acotado, terminal/quarantine, crash inyectado y takeover pasan en PostgreSQL; seis hashes idénticos, 34 Go/WhatsApp, siete cola/processors y 42 Python PASS. Las notas V274–V276 siguientes describen sus límites históricos, superados sólo en ese tramo por V277. Aún faltan host/alertas/retención y redrive operativo probados, mensajes inbound, frontend/live y los demás frentes enumerados; no se cierra toda mensajería ni el publisher/outbox. Evidencia: `reconstruction_evidence/WHATSAPP_SCOPED_JOB_COMPLETION_V277.md`.

V276 añade routing automático de statuses retenidos a envíos/citas/receipts autorizados; ya no requiere proporcionar esos IDs desde el fixture. HTTP→PG→verificador→router→observer→GET y recuperación parcial de dos mensajes pasan. Quedan claim/finalización coordinada del worker, retry/DLQ, mensajes inbound, UI y acceso/contratos live; no se confunden con el cierre de este resolver. Perfil 67/728, cuatro nuevos archivos AUTHORED e índice 0052, sin otra cola/dependencia. Evidencia: `reconstruction_evidence/WHATSAPP_RETAINED_STATUS_ROUTING_V276.md`.

El perfil de referencia enumera las 48 superficies: 41 marcadas `R`, cinco
`O`, una `O/N` y una `N`. Entre las 41 marcadas `R`, 25 filas comienzan con
"listo", diez con "input" y seis tienen descripciones condicionadas.
Son etiquetas del documento, no una auditoría que demuestre cierre de cada
superficie. `O/N` debe resolverse a una clasificación normativa en el proyecto.
La afirmación anterior "85 % real" queda retirada: sumar filas no mide esfuerzo,
calidad ni recorridos completos. El avance global permanece **NO_ESTIMADO** hasta
fijar alcance, evidencia por entregable y trabajo restante. Tampoco corresponde
interpretar 25/41 como porcentaje de cierre verificado.

La brecha intrínseca que sí corresponde cerrar antes de llamar óptimo al
acelerador es:

1. navegador conectado a backend+PostgreSQL para el journey
   público→lead→respuesta→cita/cotización/pedido, con negativos y recovery;
   V260 demuestra el primer tramo público→lead con cuatro proyectos Microsoft
   Playwright, BFF y Go reales, consentimiento/outbox/receipt en PostgreSQL,
   replay sin duplicados y rechazos de divergencia/consentimiento/clave.
   V261 extiende ese recorrido a solicitud de turno durable con capacidad,
   receipt visible y recuperación de respuesta perdida después del commit:
   cuatro navegadores, replay sin duplicación, horario lleno y payload inválido.
   V262 prueba recurso/jornada→cliente BFF→Go→PG con tokens sintéticos RS256
   verificados, una sola confirmación frente a dos intentos y lectura aislada por
   cliente. V263 añade agenda seleccionable y formulario autenticado sobre sesión
   sintética: 4 navegadores→TLS local→Next→Go→PG, sin copiar IDs/versiones,
   confirmación visible tras recarga y recuperación de respuesta perdida.
   No demuestra login/IdP real, PKI desplegada ni operación multizona/gran volumen.
   Respuesta externa, cotización/pedido y recovery completo siguen pendientes.
   Evidencia V263: `reconstruction_evidence/OPERATOR_AGENDA_BROWSER_POSTGRES_V263.md`.
   V269 corrige el transporte compartido del BFF: presupuesto por bytes durante
   lectura, cancelación y errores sin payload; 70 web PASS/1 SKIP y tres HTTP
   nativos. No cuenta como cierre de otro journey ni nuevo PASS navegador/PG.
   Evidencia: `reconstruction_evidence/BFF_STREAM_RESPONSE_BUDGET_V269.md`.
   Evidencia adicional: `reconstruction_evidence/APPOINTMENT_CONFIRMATION_BFF_POSTGRES_V262.md`.
   Evidencias anteriores:
   `reconstruction_evidence/PUBLIC_LEAD_BROWSER_POSTGRES_V260.md` y
   `reconstruction_evidence/PUBLIC_APPOINTMENT_RECOVERY_V261.md`.
2. ayuda, capacitación, soporte y runbooks ligados a la misma release;
3. telemetría del journey demostrada sin PII;
4. generación compacta de authority map/manifest por el bootstrap;
5. mutaciones de proveedores realmente seleccionadas por el blueprint —hoy
   Mercado Libre listing/stock/precio/postventa y receipts/reconciliación de
   notificaciones son brechas explícitas—. V264 corrige el normalizador WhatsApp:
   perfil WABA/teléfono exacto, status con destinatario/scope y keys estables,
   17 tests locales y reconstrucción idéntica. No cierra todavía el puente
   appointment→template→fence→inbox/status→reconciliación; evidencia
   `reconstruction_evidence/WHATSAPP_SCOPED_WEBHOOK_V264.md`. V265 añade
   `whatsappbridge.Sender` hacia el owner Python, con proceso hash-locked,
   autorización ligada al request/perfil y receipt validado. Cinco casos
   prueban replay/unknown en PostgreSQL real con proveedor sintético.
   V266 implementa el resolver durable ligado a confirmación comprometida,
   contacto y consentimiento actuales: concurrencia, revocación y fence pasan
   con fixtures SQL/proveedor sintético. V267 añade API explícita de notificación:
   HTTP→RS256/JWKS→aprobación→fence PG→child sintético, con replay y respuesta
   perdida sin duplicar invocación. V268 añade GET de resultado local histórico
   con scope actual, sesión SQL read-only y accepted separado de entrega real.
   V270 correlaciona el comprobante de envío anclado con statuses firmados:
   mensaje/destinatario exactos, deduplicación, timestamps, empate explícito y
   snapshot atómico. 41 Python y bridge Go focal/vet/build PASS; la correlación
   de artifacts no escribe el fence ni reconoce un inbox como procesado.
   Evidencia: `reconstruction_evidence/WHATSAPP_ANCHORED_STATUS_RECONCILIATION_V270.md`.
   V271 añade StatusObserver durable: anchor obtenido del fence, verificador
   Python real aislado, transacción con scope revalidado, eventos únicos y GET
   autenticado sobre sesión read-only. Concurrencia y conflictos no alteran el
   intento de envío. 42 Python, package Go/PG completo, vet/build, migración 0051
   down/up y reconstrucción 10/10 PASS.
   Evidencia: `reconstruction_evidence/WHATSAPP_DURABLE_STATUS_OBSERVER_V271.md`.
   V274 implementa receptor HTTP opt-in -> verificador Python -> inbox/job
   compartido, con original retenido sólo bajo política explícita. ACK perdido
   tras commit y cuatro replays no duplican; relectura -> StatusObserver ->
   observed_delivered funciona con routing conocido del fixture, no worker
   automático. 14 negativos HTTP, Go/PG/Python/vet/build y rebuild PASS.
   Falta routing/worker automático con retry/DLQ y retención operativos, UI y
   cuenta/contrato/consentimiento/recepción/entrega reales. Evidencia:
   `reconstruction_evidence/WHATSAPP_HTTP_DURABLE_INGRESS_V274.md`.
   V273 corrige antes del montaje la bandeja genérica existente: conexión activa
   comprobada dentro de la transacción y replay con proveedor/tipo/payload
   equivalentes, incluida desactivación concurrente. No añade otro inbox ni
   convierte el HMAC de referencia en protocolo Meta. Evidencia:
   `reconstruction_evidence/PROVIDER_INBOX_ADMISSION_SCOPE_V273.md`.
   Evidencias: `reconstruction_evidence/WHATSAPP_GO_DURABLE_BRIDGE_V265.md` y
   `reconstruction_evidence/WHATSAPP_APPOINTMENT_APPROVAL_V266.md` y
   `reconstruction_evidence/WHATSAPP_AUTHENTICATED_DISPATCH_V267.md` y
   `reconstruction_evidence/WHATSAPP_READ_ONLY_OUTCOME_V268.md`;
6. promoción por claim de más foundations que superen gates target-agnostic.

Cuentas, credenciales, corpus/ground truth, jurisdicción, IdP, CDN/WAF, nube,
sandbox, homologación, carga, seguridad ofensiva, restore/canary y aceptación
real son evidencia del proyecto. No pueden cerrarse honestamente dentro de una
biblioteca portable ni deben descontarse como si faltara escribir otro pack.

## Estado del alta anti-fake (aclaración)

El requerimiento es **básico**: email real + teléfono real + nombre real. Ya
cubierto por:

- `GO-OTP-VERIFICATION-CORE` (V222) — verifica email/teléfono por código de un
  solo uso (hash SHA-256, TTL, intentos acotados).
- `GO-APPLICANT-ONBOARDING-CORE` (V223) — `submitted → contact_verified →
  approved → activated | rejected` con deduplicación de contacto, velocidad
  anti-falsos, doble control y activación sólo desde aprobado+verificado.

**No se necesita KYC biométrico/documental estricto.**

## Bloque 1 — Código reusable y brecha exacta

| # | Estado vigente | Superficie (48) | Brecha exacta |
|---|---|---|---|
| 1 | 🟠 contratos versionados; OpenAPI→Go reusable presente | `CONTRACTS` | Kiota 1.35.0 genera/compila clientes reproducibles cuando el blueprint lo requiere; cada proyecto aún debe fijar su OpenAPI, revisar diff y demostrar auth, provider, contract/E2E y compatibilidad live |
| 2 | 🟠 Collector OTLP y reglas Prometheus materializados | `OBSERVABILITY` | falta seleccionar/probar el backend telemetry real, credenciales, red, alertas y retención del target; no es un PASS portable |
| 3 | 🟠 runner oficial + validador semántico k6 materializados | `PERFORMANCE` | cada proyecto debe aportar scripts de workload por journey y ejecutar baseline, soak y saturation/recovery contra el target real |
| 4 | ⚠️ DevSkim adaptado, OSV/SCA, fuzz nativo Go y runner/admisión ZAP materializados | `SECURITY-APPSEC` | OpenGrep/reglas GitLab permanece como evidencia fuera del perfil ejecutable porque Cosign 3.1.3 reabrió SCA. DevSkim V252 y fuzz V249 siguen disponibles; no equivalen a análisis interprocedural/cross-file. Scope/triage, targets de fuzz, ZAP, carga y aceptación siguen pendientes del target |
| 5 | ⚠️ idempotencia implementada en APIs, outbox, providers y workers | `API-BACKEND` | no existe un pack standalone universal; sólo crearlo si una capability no puede reutilizar los owners existentes |
| 6 | ⚠️ Debezium PostgreSQL outbox/inbox y fronteras Kafka materializados | `BROKER`/`ANALYTICS` | Kafka live, schema registry, TLS, replay/load y lakehouse siguen condicionados; lakehouse es `OPTIONAL` salvo necesidad del blueprint |
| 7 | 🟢 pack de candidate release portable reconstruido | `BUILD-RELEASE` | V254 ejecuta dos builds, exige ZIP byte-idéntico, OSV 0, SPDX 2.3, in-toto/SLSA y firma OpenSSH Microsoft; cada proyecto aporta su build determinista, inputs, notices y clave Ed25519 protegida, y aún debe demostrar deploy/rollback/restore del target |

## Bloque 2 — Inputs del proyecto (runtime, no código)

Obligatorios para arrancar (no los inventa el agente):

- Credenciales/cuentas y términos: LLM, WhatsApp/email/SMS, pagos,
  marketplaces, Ads/Leads, cloud, secret manager e identidad.
- Reglas aprobadas del país y negocio: fiscales/ARCA, legales, privacidad,
  consentimiento/contacto, tasas FX, territorios, precios y autoridad humana.
- Documentos: inventario de clases, corpus real, schemas/mappings por campo,
  evaluación, revisión y política de persistencia; un sample nunca concede
  exactitud automática.
- Infraestructura target: CDN/WAF, IdP, PostgreSQL/recovery, telemetry backend,
  despliegue/canary/rollback y restore reales.
- Evidencia final: journeys en navegador por rol, accesibilidad/AT, responsive,
  carga, seguridad ofensiva, resiliencia y aceptación empresarial.
- TikTok Lead: cuenta/permisos y retrieval autenticado pueden probarse live; la
  autenticidad del callback no se promueve hasta obtener un mecanismo oficial
  verificable. Mientras tanto el webhook sólo es señal y exige reconciliación.

## Ya completo (resumen)

La biblioteca materializa backend transaccional, web/BFF, dominio empresarial,
identidad/autorización, persistencia, documentos, adapters, IA gobernada,
captación/conversación, observabilidad, resiliencia y maquinaria de admisión.
“Materializado y probado localmente” no significa “proveedor, regulación o
producción aceptados”: esa promoción pertenece al proyecto y a su release
inmutable.

## Referencias

- Contrato de 48 superficies: `TOTAL_SYSTEM_CAPABILITY_CONTRACT.md`.
- Mapa de huecos: `FRANCHISE_GAP_MAP.md`.
- Estado de reconstrucción: `MARKDOWN_SYSTEM_READINESS.md`.

## Verificación actual V325

Preflight real con rutas explícitas:151steps PASS; estado global BLOCKED, único tool sin resolver Docker. .NET10.0.400 y psql18.6 disponibles, no admisión de sus grafos/servicios. Perfil67/746 y raíz161/1453/760/52 comprobados. Firefox423 continúa DIAGNOSED tras68probes sin reproducción; FAIL457 ya cerrado por V312 se retira del cursor actual, conservando los cortes históricos de este documento. ReadinessD–H/FAIL532/T2809 y decisiones faltantes permanecen abiertas.

V327: guía CLI y patch de continuidad1.3.1 probados;27tests/13comandos/8helps. Ver CLI_ENTRY_DIAGNOSTICS_V327.md; no cierre de release/usabilidad integrada ni de10macrofrentes.

Preflight V327 observado:151pasos PASS,7/8tools resueltas, Docker ausente; statusBLOCKED y acceso/corpus/target requeridos. El exit0 informativo no habilita producto ni release.

V328: readiness0.6.2/75tests y reporte exclusivo probados; gate real42bloqueos intactos. Preflight87 pendiente de nueva ejecución.

Preflight V328 observado:151pasos PASS,7/8tools,Docker ausente,estadoBLOCKED. Readiness0.6.2 conserva42bloqueos. No se reinterpretan los gates externos omitidos por falta de red/target como PASS.

V342: logger count validation corrects false delivery success; real HTTP/file reference32requests/64records, closed-sink diagnosis,28tests x3 and3finite fuzz targets PASS. Candidate remains outside composition; actual runtime host/export/retention/alerting still pending. See OBSERVABILITY_DELIVERY_INTEGRATION_V342.md.

V343 / T2809: WhatsApp0.13.0 now supplies a concrete bounded JSON reporter to the existing status host. Real PostgreSQL→TCP and mid-record loss after commit/restart prove no job/observation/audit replay. Final14host tests x3,11seeds,10s fuzz298870executions,full Go baseline/vet/build and4/4canonical parity PASS;67/746unchanged. Reporter implementation gap closed locally; real service identity/mounting/supervisor/collector/retention/alerts remain. See reconstruction_evidence/WHATSAPP_STATUS_REPORTER_V343.md.

V344 / T2809: refund host0.1.4 now stops before another token/step on failed/incomplete reporting; fixed errors/exit1 and healthy/idle policy preserved. Six baseline red modes →green; actual loop/closed-file/subprocess,17tests x3,14seeds,fuzz327445,vet/build and15file parity PASS. Three changed files,12unchanged including business/provider/DB/locks;67/746consumer. One opt-in PG test skipped, no live/refund-policy change or supervisor activation. See reconstruction_evidence/REFUND_HOST_REPORT_DELIVERY_V344.md.

V345 / T2809: SECURE-OPS1.1.3 corrects finite process capture/deadline/privacy; real-child cap/timeout/EOF tests and actual refund binary verified. Service-supervisor implementation remains missing and CONDITIONED; qualification continues through the capability-gap gate. See reconstruction_evidence/OPERATIONAL_PROCESS_BOUNDARY_V345.md.

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

Alcance Docker V361: ausencia en inventario amplio de herramientas, no causa única de los5gates abiertos. PostgreSQL nativo ya tiene V351 verificado por receipt SHA:11casos y11snapshots iguales tras reinicio controlado. No equivale a prueba de contenedores ni servicio productivo; no se relajan los checks y no hace falta Docker para estas correcciones locales.

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

V369: fuentes TRL/PEFT adquiridas e inspeccionadas; transporte core0.4.80. Ver reconstruction_evidence/TRAINING_SOURCE_ACQUISITION_V369.md. Admisión del entrenamiento sigue pendiente.

V370: candidato59distribuciones/103relaciones,58baseTRL+PEFTopcional.59METADATA hash-verified/equivalentes; resolver pip confirma59y2negativos sin wheels/install/framework execution. OSV59versiones0hallazgos declarados;30provenance subjects coinciden,firmas no verificadas. Licencias de artefactos/nativos,adquisición gobernada y runtime/gates pendientes. DISCOVERED/FAIL663 RESEARCH_INCOMPLETE. Ver reconstruction_evidence/TRAINING_DEPENDENCY_GRAPH_V370.md.

V371: core0.4.81/39files;55quarantine checks+98opaque PASS.59wheels/214055762bytes adquiridos;23854files/23795RECORD hashes inspeccionados.226native files/3SBOMs;371identidades consultadas,6records=4avisos distintos (3security+1maintenance). FAIL675 mantiene cuarentena. Safetensors sourceb7c0f38 corrige pyo3/memmap2,45registry crates0OSV sólo metadata; build pendiente.43/48sin cambio. Ver reconstruction_evidence/TRAINING_WHEEL_QUARANTINE_V371.md.

V372: opt-in WINDOWS-REFERENCE-TELEMETRY-RUNTIME0.1.0/13files connects actual GO-OFFICIAL-RETURN-REFUND-WORKER0.1.5 main and CLI to authenticated official-source Collector/Prometheus and a fresh owned PostgreSQL reference. See WINDOWS_REFERENCE_TELEMETRY_PACK_PLAN.md and reconstruction_evidence/INTEGRATED_TELEMETRY_CONTROL_V372.md. GO-OBSERVABILITY-CORE remains CANDIDATE and is not silently selected. Production target gates, other21core equivalence and training remain separate. Final integrated closure is recorded only after fresh reconstruction.

V373: opt-in HISTORY-MODEL-TRAINING-PIPELINE0.1.0/17files, HISTORY_MODEL_TRAINING_PACK_PLAN2packs/21files;26policy tests and actual earlier synthetic SFT/evaluation/rollback. Final exact reconstruction/bootstrap/E2E and TEST09 decision are recorded in reconstruction_evidence/HISTORY_MODEL_TRAINING_CONTROL_V373.md. NEW requires no history; EXISTING keeps data/model private with its own authority/quality/runtime gates. No automatic training, provider call or deployment. SDK component selection precedes AUTHORED glue; see HISTORY_TRAINING_SDK_QUALIFICATION_V373.md and HISTORY_TRAINING_SDK_GAP_V373.md. Integral franchise remains67/754.

V373 cierre201: TEST09 PASS por pack17files/composición21,27policy, bootstrap58wheels/23713file hashes, SFT real, evaluación independiente,13negativos, timeout/no replay, decisiones exactas y rollback local con baseline previo. Preflight200160pasos ejecutados PASS/55planes;164/1506/818/55, integral67/754.45/48, tres pendientes TEST02/03/07; no cambia oráculos ni hereda corpus/modelo/seguridad/deploy del consumer. Ver HISTORY_MODEL_TRAINING_CONTROL_V373.md.

V374 cierre203: TEST02 PASS como auditoría por claim21núcleos;228tests/306semillas,23targets finales (6389808ejecuciones),20packs/40files compuestos, candidato rechazado,13negativos de evidencia y Preflight202160pasos PASS.46/48controles;20CONDITIONED/1CANDIDATE y FAIL385/integración más amplia conservados.164/1506/819/55, integral67/754. No cambia48oráculos ni promueve producto. Ver CORE_CLAIM_ADMISSION_V374.md; continuar TEST03 y luego TEST07 según gates materiales.

V374 rectificación204 / FAIL710: el cierre203 de TEST02 se retracta. Fuente/tests/composición son válidos, pero FAIL385 aún exige equivalencia/integración material.45/48PASS; TEST02/03/07BLOCKED. Se preservan historia y mejoras; no se cambia el oracle ni se reduce el roadmap. Ver CORE_CLAIM_ADMISSION_V374.md, encabezado vigente.

V374 successor205: SECURE-OPS1.1.4 corrects the error-rate denominator that hid low-traffic outages. Canonical Prometheus3.14.0 rules9scenarios/13assertions and5wrapper rejections PASS; exact pin/no-engine distinction retained. Three consumer plans now select1.1.4.164packs/1507files/819Markdown/55plans;1260AUTHORED/140ADAPTED/107VERBATIM; integral67/755. No dependency acquisition, core promotion or whole-control closure;45/48 unchanged. Evidence: reconstruction_evidence/CORE_CLAIM_ADMISSION_V374.md. Full successor Preflight pending.

V374 checkpoint207: fullPreflight206 passed162executed steps and all55compositions, including exact Prometheus engine regression. SECUREOPS1.1.4 correction complete; integral67/755, library164/1507/819/55. Current45/48unchanged; TEST02/03/07 remain blocked. Retain the203retraction and actual FAIL385 integration work. See reconstruction_evidence/CORE_CLAIM_ADMISSION_V374.md.

V393 cierre251: cuatro packs incorporados/805files reconstruidos; lock122versiones/0avisos OSV,12571archivos y8534firmados exactos en candidato e instalación limpia.294tests/1skip heredado,92fases/4agendas/4lecturas/4cancelaciones PASS. Preflight250:164pasos PASS/56perfiles; disponibilidad global BLOCKED sólo Docker. Inventario165/1561/841/56,integral68/797.45/48; TEST02integral,TEST03tooling/native restante yTEST07siguen pendientes. Sharp no se instala en este perfil; V386investigación diferida.

V394: PNPM-ARTIFACT-SELECTION-GATE0.3.0/6files corrige FAIL782 (pins empresariales y Playwright obsoletos), entrega BLUEOAK-NOTICE.md ligado a5declaraciones exactas y verifica3recetas actuales;51tests y10mutaciones reales rechazadas. Payload442/22notices intactos; sin ejecución ni admisión runtime/redistribución. Evidencia: reconstruction_evidence/PNPM_CURRENT_ROUTING_NOTICES_V394.md. Los pendientes BlueOak se reducen a entrega downstream/revisión restante; ausencia de LICENSE original chownr no se falsea. Otros permisos/source/publicación de pnpm continúan abiertos.

V394 cierre253: FAIL782 reparado con selector0.3.0/6files;51tests PASS,3recetas actuales,5declaraciones BlueOak exactas con aviso local y10mutaciones reales rechazadas. Preflight252:164pasos PASS/56perfiles, nuevo control de compatibilidad de3consumers; disponibilidad BLOCKED sólo Docker.165/1561/842/56;45/48sin promoción. Payload442 y22notices intactos; pnpm no ejecutado, runtime/redistribución pendientes. Ver reconstruction_evidence/PNPM_CURRENT_ROUTING_NOTICES_V394.md.

V395: PNPM-ARTIFACT-SELECTION-GATE0.4.0/6files entrega QRCODE-NOTICE.md con encabezado de autor/modificación exacto y texto MIT completo, además del aviso BlueOak intacto.59tests PASS,3recetas reales y12negativos rechazados;10source blobs/10regiones de bundle fijados por separado, sin claim de equivalencia completa. FAIL783 entrega local resuelto, FAIL784 metadata ADAPTED reconstruida;payload442/22notices intactos. No ejecución ni admisión runtime/redistribución;45/48. Ver reconstruction_evidence/QRCODE_VENDOR_NOTICE_DELIVERY_V395.md.

V395 cierre257: aviso QRCode con copyright/permiso completos entregado y ligado a3recetas;59tests/12negativos PASS,6fuentes reconstruidas;BlueOak ypayload442 intactos. Preflight256164pasos/56perfiles PASS;Docker ausente.165/1561/843/56;procedencia1306/148/107. FAIL783/784/785 reparados;45/48sin promoción. semver-utils fuente gitHead en GitHub404, investigación pendiente sin bypass ni licencia inventada. Ver reconstruction_evidence/QRCODE_VENDOR_NOTICE_DELIVERY_V395.md.

V396: cerrado descubrimiento y entrega local de licencia original semver-utils1.1.4: MIT OR Apache-2.0 explícito, opción MIT completa y1839bytes originales retenidos. Core0.4.88/51files/196sources/25profiles añade sólo cuarentena npm exact4193bytes/SHA512;140checks y gates anteriores PASS. Planner0.5.0/6files/67tests;3recetas reales/12negativos PASS, BlueOak/QRCode y442payloadfiles intactos. Clave registry vencida documentada; no firma vigente, equivalencia build, ejecución de recetas ni admisión global pnpm.45/48sin promoción; V386/Daybreak sigue diferido. Ver reconstruction_evidence/SEMVER_ORIGINAL_LICENSE_DELIVERY_V396.md.

V396 cierre262: licencia original semver-utils1.1.4 localizada y entregada bajo opción MIT completa;3recetas/12negativos/67tests PASS. Core0.4.88/51files/196sources/25profiles,140checks; planner0.5.0/6files. Preflight261164pasos/56perfiles PASS;Docker ausente.165/1562/844/56;procedencia1307/148/107.45/48sin promoción, Daybreak/V386 diferido. Fuente/licencia semver resuelta en alcance local; firma registry vencida y admisión pnpm restante condicionadas. Ver reconstruction_evidence/SEMVER_ORIGINAL_LICENSE_DELIVERY_V396.md.

V397: PNPM-ARTIFACT-SELECTION-GATE0.6.0/7files entrega el conjunto completo de47textos retenidos/151253bytes en3recetas:141copias verificadas independientemente,12negativos reales y82tests PASS.22notices originales comparados contra payload y25evidencias conservan alcance; BlueOak/QRCode/semver y442payloadfiles intactos. Catálogo ADAPTED con términos por texto, sin algoritmos ni promoción de licencias/source/relinking/runtime.45/48; Daybreak diferido. Ver reconstruction_evidence/PNPM_RETAINED_NOTICE_DELIVERY_V397.md.

V397 cierre264: conjunto47textos/151253bytes entregado en3recetas,141copias exactas/12negativos/82tests PASS. Preflight263164pasos/56perfiles PASS;165/1563/845/56 y1307/149/107. Docker ausente,45/48sin promoción. README/roadmap vigentes sincronizados y cortes previos preservados (FAIL794). Licencias/source/relinking/publicación restante y Daybreak diferido siguen explícitos. Ver reconstruction_evidence/PNPM_RETAINED_NOTICE_DELIVERY_V397.md.

V398: PNPM-ARTIFACT-SELECTION-GATE0.7.0/8files entrega fuente original next-path1.0.0, manifiesto y MPL completa en3recetas:9copias exactas/12negativos reales/90tests PASS. Commit oficial y3Git blobs verificados;4sentencias comparadas bajo adaptadores explícitos,10mutaciones rechazadas. No equivalencia runtime ni reproducibilidad pnpm. Colección47, suplementos anteriores y442payloadfiles intactos.45/48; Daybreak diferido. Ver reconstruction_evidence/NEXT_PATH_MPL_SOURCE_DELIVERY_V398.md.

V398 cierre266: next-path fuente/manifiesto/MPL entregados,9copias exactas/12negativos/90tests/8files PASS. Preflight265164pasos/56perfiles PASS;165/1564/846/56 y1307/150/107. Docker ausente,45/48 sin promoción. Entradas vigentes sincronizadas con historia conservada. Daybreak diferido; integración funcional,source/relinking/publicación restante yrelease pendientes. Ver reconstruction_evidence/NEXT_PATH_MPL_SOURCE_DELIVERY_V398.md.

V399: continuidad de cierre conciliada con V372/V393 y21packs/42fuentes byte-idénticos aV374. No repetir esas suites sin delta; agrupar trabajo por bloqueo material.45/48 intacto; TEST02integración,TEST03grafo/tooling yTEST07release siguen abiertos. Observabilidad candidata0.3.3 separada del runtime probado; Sharp ausente del perfil web no reabre Daybreak. Ver reconstruction_evidence/CURRENT_CLOSURE_BOUNDARIES_V399.md.

V400: encuestas conectadas con PostgreSQL/OIDC/Next, recuperación GET,4navegadores/8respuestas/8POST y retención CLI1/1/0.23archivos nuevos AUTHORED, packs GO-CUSTOMER-SURVEY-API0.1.0 y TS-CUSTOMER-SURVEY-PORTAL0.1.0, CONDITIONED.45/48sin promoción; TEST02/03/07 siguen abiertos. Ver reconstruction_evidence/CONNECTED_CUSTOMER_SURVEYS_V400.md.

V400 cierre270: encuestas conectadas y retención transaccional incorporadas. Go/TSsurvey0.1.1, app1.9.4, HTTPref0.1.15;23archivos nuevos,828/828 reconstruidos.4navegadores/8respuestas/8POST,24concurrentes,restart,CLI1/1/0,fuzz6219;correcciónSQL2→0tablas parciales PASS. Preflight269164pasos/56perfiles PASS más deltaSQL focal;167/1587/850 y1330/150/107,franquicia70/820,HTTP71/828,web7/134.45/48sin promoción; Daybreak diferido. Ver reconstruction_evidence/CONNECTED_CUSTOMER_SURVEYS_V400.md.

V401: contención ZIP terminada en candidato pnpm aislado;16módulos adm-zip retirados,441archivos preservados/1bundle cambiado,4rechazos sinIO y3instalaciones offline reales PASS.475identidades conocidas/0hallazgos; SBOMrecursivo no cerrado;12571archivos node_modules idénticos. Pack0.8.0/9files,167/1588; pnpm general BLOCKED, TEST02/03/07 siguen abiertos,45/48; Daybreak diferido. Ver reconstruction_evidence/PNPM_ZIP_CONTAINMENT_V401.md.

V401 reconciliación final: comparación12571/12571idéntica completada antes de la solicitud de detenerla; no se detuvo proceso. FAIL807 separa las3instalaciones iniciales offline de un exec Next que descargó71paquetes por configuración omitida. Exec corregido con store/offline explícitos PASS sin descargas. Historial/log anterior retenido; no prueba global de ausencia de red. reconstruction_evidence/PNPM_ZIP_CONTAINMENT_V401.md
