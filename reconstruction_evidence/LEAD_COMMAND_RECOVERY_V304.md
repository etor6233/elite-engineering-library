# V304 — leads: recuperación y auditoría autenticada

> **Errata de trazabilidad V317 (2026-09-08), prevalece sobre las menciones históricas de Go1.26.7 debajo.** Se retira la atribución histórica exacta Go1.26.7 no vinculada a un receipt de identidad por ejecución. No se sustituye por1.26.8 sin prueba del artefacto concreto. La ruta current-toolchain de V281 contiene1.26.8; official-toolchain contiene1.26.7. El cuerpo original queda conservado como historia. V317 revalida únicamente la composición actual y el scope allí enumerado con1.26.7 observado, no todos los snapshots previos. Ver [TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md](TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md).

2026-09-07. Checkpoint35/V303 → siguiente tramo T2803/T2804/FAIL457.
Mantenimiento de biblioteca; conserva V295, credenciales pendientes, decisiones
comerciales sin resolver y ARCA diferida. Sin cuentas, compras ni efectos externos.

## Implementación canónica

GO-FRANCHISE-CUSTOMER-JOURNEY-API0.10.9,
TS-FRANCHISE-JOURNEY-PORTALS0.14.2 y
MICROSOFT-PLAYWRIGHT-BROWSER-GATE0.1.14. Ocho archivos existentes modificados,
ninguna dependencia, migración ni archivo de producto nuevo. Planes backend,
web, integral y serverless sincronizados; integral67 packs/745 archivos.

LeadActions sustituye sólo asignación/transición del sender genérico:
controles cerrados antes de hidratación, fence síncrono, marker de sesión con
versión/acción exclusivamente, namespace hash de tenant/subject/organización.
Sin contacto, responsable, notas ni credenciales en ese marker. Error de escritura
impide POST. Timeout10s y recibo contrastado con ID/organización/versión y efecto
esperado; respuesta perdida o no demostrada conserva incertidumbre y bloqueo.
«Consultar lead» recarga por GET sin reenviar. Una versión durable posterior
permite revisar responsable/estado y continuar; no atribuye el cambio a esta
solicitud si pudo actuar otro operador. Ausencia, fallo o versión sin avance
conservan la referencia y requieren revisión autorizada. No se inventa rollback
de un comando aceptado ni se cambia la máquina de estados existente.

Ayuda lead-command-view/1.0.0 integrada con práctica sintética, estado actual,
consulta, privacidad del marker y soporte. No equivale a capacitación evaluada
persistente ni a soporte integral.

AssignLeadAs/TransitionLeadAs exigen actor; HTTP lo toma exclusivamente del
Subject autenticado. Repositorio sin métodos auditados falla cerrado, sin IDs
ni fallback silencioso. PostgreSQL añade actor_subject al outbox existente en
la misma transacción CAS del cambio. Los métodos internos antiguos mantienen
compatibilidad explícitamente no auditada: no se afirma que sus llamadores
hereden esta garantía. No nuevo tipo de evento ni contrato externo certificado.

## Procedencia y autoridad

Todo el delta es AUTHORED / LicenseRef-Workspace-Owner, no código copiado de una
empresa externa. Se conservan revisiones/licencias/locks de los packs existentes.
Se reutilizan las consultas oficiales del2026-09-07 de V303, para el mismo claim:

- [AWS Builders' Library — retries seguros](https://aws.amazon.com/builders-library/making-retries-safe-with-idempotent-APIs/): respuesta perdida no demuestra ausencia de efecto; no reenvío ciego.
- [React — refs](https://react.dev/learn/referencing-values-with-refs): fence mutable de UI; no reemplaza la exclusión transaccional.

FRONTEND_PRODUCT_ENGINEERING_UX2.2 y SOFTWARE_BACKEND_API_ENGINEERING8 conservan
su autoridad estrecha. Citar estas fuentes no implica aprobación empresarial
del código local, actualización de todo el corpus ni admisión productiva nueva.

## Evidencia reproducible

Staging: `%LOCALAPPDATA%/Temp/elite-v304-739f947ba60d436eab4539a03e7c1f40`.
candidate: web compilada; rebuilt: composición desde Markdown a destino ausente,
745 archivos/67 packs. Ocho archivos delta coinciden por SHA256 entre ambos.
Backend final ejecutado desde rebuilt; frontend y spec comparados con candidate.

- lead-red.log: baseline V303 confirmó HTTP200/versión2 antes de abortar respuesta;
  el test de UI de recuperación falló. Baseline original restaurado tras ensayo.
  Probe SQL read-only posterior:1 lead.assigned,0 con actor_subject; FAIL459.
- lead-candidate.log:15 fases Chromium desktop PASS antes del test de atomicidad.
- browser-final.log:60 fases (15 por proyecto) PASS: Chromium desktop/mobile,
  Firefox desktop y WebKit desktop. Conserva cotización/pedido/stock/solicitud
  local de pago/entrega/secciones/resolución; no cobra ni prepara entrega inicial.
- En cada nueva fase lead: asignación y tres transiciones, pérdida de respuesta
  después del commit, JSON200 vacío, GET503 y recuperación, storage fallido con
  cero POST, carrera200/409, replay409, organización ajena403, actor JSON400 y
  token sin permiso403. GET no reenvía. Estado final lost/version5/sales-owner;
  exactamente4 eventos con actor lead-operator, sin duplicados. Los requests
  adicionales son inyección de test, no reintentos automáticos del producto.
- Atomicidad: evento repetido provoca error y conserva lead new/version2/assignee;
  transición válida posterior produce contacted/version3 y dos eventos con actores
  distintos. Pool nuevo confirma estado/eventos. Tres repeticiones PASS aisladas
  y tres en suite PostgreSQL final junto al owner completo.
- lead-atomicity.log conserva tres fallos22P02 de fixture: event_id era texto
  en lugar de UUID; FAIL460. Corrección usa randomid existente, conserva el mismo
  UUID en los intentos duplicados y repite sin cambiar esquema/oráculo.
- Agenda cuatro proyectos PASS. Web105 PASS/1 SKIP (12 archivos), build PASS.
  Go test ./..., vet ./... y build ./... PASS desde rebuilt. Los tests opt-in
  omitidos en el comando general sólo se acreditan por sus logs conectados.
- Go1.26.7, PostgreSQL18.6, Node24.14.1, pnpm11.19.0, Next16.3.2,
  Playwright1.62.1; locks existentes, instalaciones offline/frozen-lockfile.
  Gate anidado instalado con --ignore-workspace. Sin nuevo SCA/release de librerías.

Reproducción: componer FRANCHISE_COMPLETE_PACK_PLAN, instalar/build web y gate
con locks como V303. PostgreSQL loopback descartable elite_confirmation_*,
53 migraciones; nunca DB de negocio. Variables TEST_DATABASE_URL, ELITE_WEB_ROOT
absoluto, ELITE_QUOTE_E2E, ELITE_ORDER_E2E, ELITE_PAYMENT_E2E,
ELITE_DELIVERY_READ_E2E, ELITE_DELIVERY_ACTION_E2E, ELITE_OPERATOR_SECTIONS_E2E,
ELITE_RESOLUTION_RECOVERY_E2E y ELITE_LEAD_RECOVERY_E2E=1; ELITE_QUOTE_PROJECT vacío.
Ejecutar `go test ./internal/platform/httpapi -run
'^TestQuoteAcceptanceBrowserPostgres$' -count=1 -v -timeout 15m`.
Agenda: ELITE_CONFIRMATION_E2E=1 y TestAppointmentAgendaBrowserPostgres.
PostgreSQL: `go test ./internal/platform/postgres -run
'^(TestFranchiseJourneyPersistenceIsolationAndReplay|TestLeadAuditedAtomicity)$'
-count=3 -v`. Registros y runners locales preservan las invocaciones.

## Límites y continuidad

FAIL459/460 corregidos en este scope; FAIL457 sigue abierto para cotización,
recursos, disponibilidad, checklist, recepción/disposición y otros comandos.
Próximo paso del mismo T2804: cancelación de intervalo existente, consulta del
owner y recuperación; creación de intervalos permanece separada hasta resolver
su identidad de operación. No nueva política comercial ni segundo roadmap.

No se prueban todas las variantes de marker corrupto, storage perdido después
del envío, múltiples pestañas, accesibilidad integral o todos los estados de lead.
sessionStorage no es bóveda ni log empresarial durable; PostgreSQL es autoridad.
Reabrir pool/GET no demuestra reinicio del cluster/PITR. No prueba de proveedores,
cobro/conciliación reales, entrega inicial/liberación comercial, efectos de retorno
finalizados, carga, seguridad ofensiva, deployment/rollback, aceptación o ZIP.
No cierre T2804 completo ni de la biblioteca por un PASS local.

## Integridad

| Archivo/log | SHA256 |
|---|---|
| lead-red.log | 08646777aec00644adc43c269303153f2b4b268d8630b45a4d8f698a95ace597 |
| lead-atomicity.log | dfc31c692715a63bad7f6b8232d0c05a13466b3204b5b6a6f1aa9ab9b6148219 |
| lead-atomicity-fixed.log | 7df45a875b57b7d809c12e34b794294007c9be60273668808ed7573572e1c5ac |
| lead-candidate.log | 82d5dfd8fab05296593d82c65ee3d2d253a62a6f4e6580a1d5d0d9e028135ec8 |
| web-test.log | f229717fa61990336dd91d86ffd74639ece81c4a9b036c01e1021b5f8fa0b0dd |
| web-build.log | a01113e0737260912e5e6680346f429f09bb867e315f929a1288321bf76262bc |
| go-test.log | f6a96f4b1814c7929ec83593e6744cff16e47a7647556a14c3b5b2847d81e4bd |
| go-vet.log | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |
| go-build.log | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |
| postgres-final.log | 007af44eac144bac2917b5f4033b4387c247c023a11a62bd2fd386dfc629dcdb |
| agenda-final.log | 5b920266e8be7c7ca32de09773ab3704af817ff66be5e8423f019fe31ad04ca9 |
| src/components/franchise-command-panel.tsx | 8901199d61dc73062cccf8f14f4dfcdcbe8cf7e258cd4c9f92f734e745ae6781 |
| microsoft_playwright_browser_gate/tests/enterprise-web.spec.mjs | 139dfa52786e618023db91de1e8c996173e8ef40d0f7021d5a6a5c9ed9ed3ac3 |
| internal/franchisejourney/service.go | 67c260fccdaf323b39ba49b83490470beac2ddbbf5c887f605554013ce2dd0f9 |
| internal/franchisejourney/service_test.go | 6720728a771d1d1823aabf206aeb9cd9e7c9b015eb67dfb4ea270815fd77a557 |
| internal/platform/postgres/franchisejourney.go | d79e2d959b49a64408c534ca5b29daa04381c9d0b150db99e80c137dafaf83d1 |
| internal/platform/postgres/franchisejourney_integration_test.go | 5a05751bb738ec290721d596c32adf11dacac434235c6e2c0d11e1f6fbc8753d |
| internal/platform/httpapi/franchisejourney.go | 49f9030bd41bf48a2576d497288bb2a02523d4696680e62dadc93b9045e9e191 |
| internal/platform/httpapi/franchisejourney_test.go | 69b7362ebc0376a2e2b9741ae72414cc816219835e3b575baaeb6caf4da05a95 |
| browser-final.log | 58229b2a77c0ac34f4d5efc872574a78d18368a565a271979dbe6029a3a121b2 |

