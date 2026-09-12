# V305 — cancelación de disponibilidad y coherencia de reservas

> **Errata de trazabilidad V317 (2026-09-08), prevalece sobre las menciones históricas de Go1.26.7 debajo.** Se retira la atribución histórica exacta Go1.26.7 no vinculada a un receipt de identidad por ejecución. No se sustituye por1.26.8 sin prueba del artefacto concreto. La ruta current-toolchain de V281 contiene1.26.8; official-toolchain contiene1.26.7. El cuerpo original queda conservado como historia. V317 revalida únicamente la composición actual y el scope allí enumerado con1.26.7 observado, no todos los snapshots previos. Ver [TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md](TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md).

2026-09-07. Checkpoint36/V304 → tramo T2803/T2804 del roadmap vigente.
Mantenimiento local, no proyecto productivo. V295, accesos pendientes, decisiones
comerciales y ARCA diferida se conservan. Sin cuentas, compras ni efectos externos.

## Código y alcance

Journey0.10.10, Portal0.14.3 y BrowserGate0.1.15. Seis archivos existentes
modificados, ninguna migración ni dependencia nueva; 67 packs/745 archivos.
Planes backend/web/integral/serverless sincronizados. Reconstrucción a destino
ausente; los seis archivos coinciden por SHA256 entre candidate y rebuilt-r3.

AvailabilityCancellation sustituye sólo el botón de cancelación del sender
genérico. Conserva ID/versión y motivo schedule-correction existentes. Controles
cerrados antes de hidratación, fence síncrono y referencia de sesión mínima:
sólo versión, bajo hash de tenant/subject/organización e ID opaco de intervalo.
No contiene horarios, motivos, contactos ni credenciales. Falla al guardar:
cero POST. Respuesta perdida o no demostrada: incertidumbre, no «No se aplicó».
Recibo exige ID/organización/versión/cancelled/tipo/recurso/fechas coincidentes.
Timeout10s es límite local conservador, no presupuesto de latencia medido.

Consultar intervalo hace GET/recarga, no reenvía cancelaciones. Sólo cancelled
con versión posterior resuelve la referencia. Si no aparece en la ventana de90
días, falla la consulta, permanece active o hay referencia inválida, no deduce
ausencia del efecto: conserva bloqueo y requiere revisión autorizada.
Una consulta no atribuye el cambio a esta solicitud si pudo actuar otro operador.
El permiso de lectura no muestra botón de cancelación. Ayuda
availability-cancel-view/1.0.0 conecta acción, incertidumbre, restricción por
citas activas, práctica sintética y soporte. No es capacitación evaluada.

## Brecha adicional corregida: cupos públicos y reserva

FAIL462: cancelar una jornada sin citas no cerraba el slot previamente abierto.
PublicAppointmentSlots y RequestAppointment no consultaban disponibilidad actual.
Se reprodujo publicación y reserva indebidas tanto tras cancelación como ante
indisponibilidad; dos carreras aceptaron simultáneamente reserva y cambio.

La corrección aplica la misma política que ya exigía CreateAppointmentSlot:
una jornada organizacional active cubre todo el cupo y no existe indisponibilidad
organizacional active solapada. Publicación y reserva comprueban ese predicado.
No cierra/modifica slots automáticamente, no inventa horarios ni cancela citas.

RequestAppointment y cambios de disponibilidad organizacional comparten la clave
de advisory lock transaccional ya usada por los triggers. El lock se adquiere
antes de la sentencia decisiva, con Read Committed explícito; cancelación también
lo toma antes del UPDATE. El replay idempotente previamente completado conserva
su resultado. Reserva/cancelación/bloqueo concurrentes no pueden aprobar efectos
contradictorios en los owners probados. Restricción existente de citas activas,
CAS, autorización y outbox permanecen; no se elimina ningún guard.

## Procedencia y fuentes

Delta AUTHORED / LicenseRef-Workspace-Owner. No código copiado de AWS,
Microsoft o PostgreSQL. Dependencias y revisiones fijadas se conservan; citar
una fuente no implica que esa empresa apruebe el código local.

- [AWS Builders' Library — retries seguros](https://aws.amazon.com/builders-library/making-retries-safe-with-idempotent-APIs/): se reutiliza la consulta2026-09-07 de V303 para no interpretar respuesta perdida como ausencia de efecto.
- [React — refs](https://react.dev/learn/referencing-values-with-refs): misma consulta V303; fence local, no exclusión durable.
- [PostgreSQL18 — explicit locking](https://www.postgresql.org/docs/18/explicit-locking.html): consultado2026-09-07; locks transaccionales y disciplina cooperativa.
- [PostgreSQL18 — transaction isolation](https://www.postgresql.org/docs/18/transaction-iso.html): consultado2026-09-07; Read Committed obtiene snapshot por sentencia.

FRONTEND_PRODUCT_ENGINEERING_UX2.2 y SOFTWARE_BACKEND_API_ENGINEERING8 mantienen
la autoridad de UX/recuperación/reintentos. No nueva admisión de upstream,
actualización global, SCA ni promesa de rendimiento por usar locks.

## Evidencia

Staging: `%LOCALAPPDATA%/Temp/elite-v305-76e6baedc09645d48a1c7f0ab1add8d8`.
Baseline V304 original restaurado tras ensayo. Web construida en candidate;
backend final y tests PostgreSQL desde rebuilt-r3.

- availability-red.log: HTTP200/cancelled/versión2 antes de abortar respuesta;
  baseline mostró «No se aplicó: Failed to fetch» y botón disponible. Regresión
  falló por status de recuperación ausente; trace/error-context preservados.
- availability-candidate.log: Chromium16 fases PASS del primer candidato.
  browser-final.log:64 fases de UI previas al descubrimiento/corrección FAIL462;
  no usar ese PASS para acreditar coherencia de reservas.
- availability-booking-red.log: dos cupos indebidos publicados, pero la prueba
  de reserva usaba clave demasiado corta y falló23514 (FAIL463).
  availability-booking-red-valid.log conserva el baseline con fixture válido:
  dos casos publican1 cupo y crean reserva/idempotencia/outbox (3 efectos);
  dos carreras arrojan2 éxitos/0 conflictos.
- availability-booking-fixed.log:10 repeticiones de cuatro escenarios
  (cancelled/unavailable/race-cancel/race-unavailable) y atomicidad PASS.
  Donde reserva gana, replay devuelve el mismo registro sin nuevo efecto.
- Atomicidad: evento de auditoría duplicado23505 revierte cancelación/actor;
  concurrencia produce un éxito/un conflicto; nuevo pool lee cancelled/version2
  y un evento con actor. Jornada con cita activa rechaza cancelación y conserva
  active/version1/cero eventos de cancelación.
- browser-final-r3.log:64 fases (16×4) PASS, Chromium desktop/mobile,
  Firefox desktop y WebKit desktop, sobre la revisión final.
  Por proyecto:4 cancelaciones únicas (working/unavailable), cada una versión2
  y evento con actor scheduler. Respuesta perdida, JSON200 vacío, GET503, storage
  fallido con cero POST, carrera200/409, replay409, organización403,
  actor JSON400 y permiso HTTP403. El lector no recibe botón de cancelación.
- Una quinta jornada tiene cita activa: HTTP409, permanece active/version1,
  cero evento de cancelación, GET conserva referencia y bloqueo.
  Fallo de disponibilidad no oculta pedidos ni filtra detalle interno.
- En la misma prueba conectada, dos slots antes publicados desaparecen de la
  consulta pública después de cancelar la jornada; dos reservas HTTP tardías
  reciben409 y SQL confirma cero reservas/idempotencias/eventos nuevos.
  Prueba de endpoint público desde runner, no recorrido completo del formulario
  de reservas público ni entrenamiento de usuarios.
- Agenda4 PASS, web105 PASS/1 SKIP y build PASS; Go test/vet/build PASS.
  Suite PostgreSQL completa más tres regresiones repetida3 PASS desde la
  reconstrucción final. Los opt-in no ejecutados por go test genérico sólo se
  acreditan por sus logs conectados.
- Go1.26.7/PostgreSQL18.6/Node24.14.1/pnpm11.19.0/Next16.3.2/Playwright1.62.1.
  Instalación offline/frozen-lockfile; gate anidado --ignore-workspace.

Reproducir como V304, agregando ELITE_AVAILABILITY_RECOVERY_E2E=1 al runner
TestQuoteAcceptanceBrowserPostgres; dejar ELITE_QUOTE_PROJECT vacío para4.
Usar TEST_DATABASE_URL loopback elite_confirmation_* descartable y53 migraciones,
ELITE_WEB_ROOT absoluto y los flags V304. Agenda: ELITE_CONFIRMATION_E2E=1.
PostgreSQL: go test ./internal/platform/postgres -run
'^(TestFranchiseJourneyPersistenceIsolationAndReplay|TestLeadAuditedAtomicity|TestAvailabilityCancellationAtomicity|TestBookingRequiresCurrentAvailability)$'
-count=3 -v. No ejecutar fixtures contra datos de negocio.

## Límites y continuación

FAIL464: browser-final-r2.log conserva cuatro rechazos400 del probe de cupos;
se corrigió sólo el fixture de90 días a from futuro y ventana14 días, dentro
del contrato público de31. Go/SQL/web y harness idénticos por hash permiten
reutilizar las suites r2; sólo cambia el spec Playwright y se repite toda su matriz.

FAIL461: import fmt omitido en harness, compilación fallida antes de ejecutar;
corregido y compilado en gates posteriores. FAIL463 conserva fixture fallida,
no se debilita constraint. FAIL462 se cierra sólo para los owners/condiciones
probados; advisory locks requieren cooperación, no protegen escrituras SQL
arbitrarias fuera de esos owners. No prueba integral de asignación de recursos,
cambios de duración, cierre de slots, todas las intercalaciones, carga ni deadlocks.

FAIL457 sigue OPEN para creación de intervalos, cotización, recursos, checklist,
recepción/disposición y otros comandos. Siguiente tramo: create-quote del sender
genérico, preservando el contrato Idempotency-Key existente y resolviendo su
identidad/resultado ante incertidumbre antes de permitir otro envío.
No exigir nuevas decisiones comerciales para corregir ese defecto de transporte.

No equivalencia entre cancelación de disponibilidad y cancelación/reprogramación
de citas. No invariantes empresariales verificadas tras reinicio/PITR (sí reconexión de pool), a11y integral,
multi-tab, pérdida de storage posterior a envío, soporte completo, proveedores,
cobro, liberación/entrega inicial, despliegue/rollback, aceptación o ZIP.
La serialización por organización debe medirse con carga del target; no se
presenta como mejora de latencia. T2804 y la biblioteca no quedan terminados.

## Integridad

FAIL465: launcher con Start-Process -Wait mantuvo pendiente el pg_isready
encadenado hasta apagar PostgreSQL (exit1/no response). No fue el probe utilizado
para habilitar los tests: un pg_isready independiente pasó mientras estaba activo.
postgres-control.ps1 corrige espera acotada sólo del PID pg_ctl; ciclo real
start/readiness/stop PASS en0.80s del mismo cluster sintético, sin borrar datos.
No acredita invariantes del negocio después de ese reinicio. Cluster detenido.

| Archivo/log | SHA256 |
|---|---|
| availability-red.log | 1974118a00f3aca6231b98744f663801df688a38ce7e48f31af6bf7ddabecf03 |
| availability-candidate.log | c02af09642df97f803a17731b1d859eeecbf9abbbe06d8af15f886edae1d08fa |
| availability-atomicity.log | 20537f8ce5e249b8a913352c56310f987062527d1572ac5eb32a09cbe2b5f3b2 |
| availability-booking-red.log | 90323ab65193e254ec58aaef8efebb42a8672733c4df8b59b3d5260e576b5ab8 |
| availability-booking-red-valid.log | c6c8147a88df543356f91a12573808920bb5fed4ce1e68b5e853e61a5df6cb5f |
| availability-booking-fixed.log | 335988f7c080d6761e24c851f55a82333a9d0739851c27ded633139d1f9a8ca4 |
| browser-final.log | 2385cbc136159f7f5bfb9be902ac6880c2b55c173ac65d7ec374d6c89fbca169 |
| browser-final-r2.log | 715beeff43fd12c70f04f72d6f1fe854d41a3dd6b31dc47aa724ff67ee5ecbbd |
| go-test-r2.log | 066427b993c2f8758b2d9d12d75ef1ab47d938422b864d3c3c9fcb4685969207 |
| go-vet-r2.log | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |
| go-build-r2.log | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |
| postgres-final-r2.log | 48503e404a8cb696582b31c41e459265ace36fbdf2423c3607572348e895bebf |
| agenda-final-r2.log | e6cff36325f8654097b00c8c4b85cde421e32407b00d39320eb950c372c9bba2 |
| web-test.log | 0269458b165b52d1913b8ff7e385964830e0b631d47ea4bd914cbe90bf486b00 |
| web-build.log | 074c94a2eec2c1b6c8d953301393528f47febcc01821610d1526341022fe44af |
| src/components/franchise-command-panel.tsx | 16c1ff9508b851854499b22fa57d26b7be0b01c98886568c25804ebbea9389c9 |
| microsoft_playwright_browser_gate/tests/enterprise-web.spec.mjs | 1ecc437696c2b994e35672d1acc93264daf5d84ed99f00ffb25ab68bf9648a8b |
| internal/platform/httpapi/franchisejourney_test.go | 03c3bdc5585870316f2c09e6f53d0fb88e7834ab63dc488e68087c13bd444113 |
| internal/platform/postgres/franchisejourney_integration_test.go | 1dca4b7f6077454ede3ad0605197f69d1570b6ea1a51aa3f544e09b8052eaf85 |
| internal/platform/postgres/franchisejourney.go | ac09ab23f3c7da2ca3e1181312cf68e85bfe5f286681e3547530e56587b5050d |
| internal/platform/postgres/franchisejourney_availability.go | 897aaf433e9f7679a3448d851afa336e462ce99853ac0399ea146577cd26a869 |
| browser-final-r3.log | b0c0f8e15cfd9d0147c37f5e95fee5716f79cc8ad98f6c9b75eba526aa54310c |
