# V300 — lectura coherente y recuperación de recepción/discrepancia

> **Errata de trazabilidad V317 (2026-09-08), prevalece sobre las menciones históricas de Go1.26.7 debajo.** Se retira la atribución histórica exacta Go1.26.7 no vinculada a un receipt de identidad por ejecución. No se sustituye por1.26.8 sin prueba del artefacto concreto. La ruta current-toolchain de V281 contiene1.26.8; official-toolchain contiene1.26.7. El cuerpo original queda conservado como historia. V317 revalida únicamente la composición actual y el scope allí enumerado con1.26.7 observado, no todos los snapshots previos. Ver [TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md](TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md).

2026-09-07. Mantenimiento correctivo desde checkpoint31; T2802/T2803/T2804
del roadmap existente. V295–V299 permanecen registrados, ARCA diferida.
Este cierre corrige lectura y respuesta incierta de operaciones existentes;
no sustituye creación inicial de entrega, liberación financiera ni cobro.

## Cambio canónico y procedencia

Journey0.10.6, Commerce0.6.1, portales0.13.1 y browser gate0.1.11.
Siete archivos existentes modificados, cero dependencias, migraciones o módulos
nuevos; composición67 packs/745 archivos. Cambios AUTHORED / CONDITIONED.
No son código copiado o aprobado por Microsoft, Google ni PostgreSQL.
El runtime oficial Playwright y las dependencias fijadas conservan sus locks.

Se consultó la documentación oficial de [aislamiento de PostgreSQL18](https://www.postgresql.org/docs/18/transaction-iso.html)
para el snapshot estable. FRONTEND_PRODUCT_ENGINEERING_UX gobierna visibilidad
del estado, prevención y recuperación. Estas fuentes no deciden reglas
comerciales ni convierten el código local en una implementación de sus autores.

### Lectura y privacidad

V299 protegía comandos, no todas las proyecciones. El baseline reprodujo siete
lecturas indebidas: tres CustomerJourney y dos Operations con vínculos de
acta/pedido/cliente/stock incoherentes, más dos lecturas de una excepción
asignada a otro cliente con nota sintética. El pedido de organización ajena ya
queda fuera de Operations; no se cuenta como bypass de esa consulta.

CustomerJourney usa ahora una única transacción Repeatable Read / Read Only
para sus cinco consultas. Verifica acta→pedido/cliente/stock y las excepciones
antes de devolver datos; cualquier error retorna resultado vacío, no parcial.
DeliveryExceptions aplica la misma validación relacional. Operations conserva
su snapshot y rechaza handovers incoherentes. Los items de checklist se limitan
a las actas ya seleccionadas; no prueba paginación completa ni SLO de gran volumen.

El portal distingue indisponibilidad de una lista vacía. Ante fallo no renderiza
acciones, ofrece consulta GET y ayuda contextual handover-read-view/1.0.0.
El reintento de lectura no crea ni acepta una entrega. Esto no corrige ni
certifica todas las otras páginas del portal, capacitación o soporte integral.

### Respuesta perdida después de un efecto real

El baseline action-red.log recibió200/accepted del backend y perdió la respuesta;
la UI decía falsamente “No se registró: Failed to fetch”. Esa reproducción roja
alcanzó aceptación solamente; no se afirma un segundo rojo para rechazo.

CustomerHandoverActions usa fence síncrono, timeout de10s y validación del recibo.
Ante resultado no comprobado mantiene los formularios deshabilitados e informa
que la operación pudo haberse registrado. Una navegación GET permite ver el
estado durable antes de actuar. No reenvía automáticamente ni inventa un efecto.
La versión optimista y el outbox continúan en el owner PostgreSQL existente.
El test conectado prueba pérdida post-commit en aceptar y rechazar; no hay una
prueba dedicada de cada recibo malformado ni de expiración completa del timeout.

## Verificación reproducible

Staging local: `%LOCALAPPDATA%/Temp/elite-v300-ef0b29864b3441b3ba76a0859b6332d6`.
`candidate` se devolvió a los cuatro packs; `final-source` se recompuso desde
Markdown en destino ausente y `rebuilt` ejecutó el web. Siete archivos coinciden
byte a byte entre los tres árboles. `verified-source` conserva la etapa anterior
de lectura, sin el posterior fix de incertidumbre; no usarla como cierre final.

- Go1.26.7: test ./..., vet ./..., build ./... PASS desde final-source.
  Los tests opt-in sin entorno en la suite raíz no se cuentan como ejecutados.
- PostgreSQL18.6/schema53: TestFranchiseJourneyPersistenceIsolationAndReplay
  tres veces PASS desde rebuilt. Conserva12 negativos de escritura V299,
  locks,16 aceptaciones concurrentes/1 éxito/15 conflictos, rollback por fallo
  de outbox y lectura con un pool nuevo. Añade los siete negativos de lectura.
- Un tracer sincroniza una modificación concurrente del cliente del pedido
  después del primer SELECT: snapshot actual coherente, siguiente lectura
  rechazada y recuperación tras restaurar la relación. Verifica transacción
  activa y13 IDs para la consulta acotada de checklist. No es benchmark.
- TestCustomerJourneyReadFailureIsRedactedAndRecoverable PASS: RS256/JWKS local,
 401/403/403/500 sin datos parciales privados y recuperación200. Usa repository
  double para esa prueba HTTP; no se presenta como prueba PostgreSQL.
- Frontend:105 tests PASS,1 SKIP declarado; build Next16.3.2 PASS.
- TestQuoteAcceptanceBrowserPostgres:48 fases PASS,127.94s observados;
  Chromium desktop/mobile, Firefox desktop y WebKit desktop.12 fases por proyecto:
  initial/read-failure/restart, operation/operation-read-failure/operation-restart,
  payment/payment-disabled/payment-restart, delivery-read-corrupt,
  delivery-read-restored y delivery-action-loss.
- Las ocho fases corrupt/restored usan PostgreSQL→Go→BFF→Next real: relación
  sintética inválida causa alerta sin datos/acciones; restaurar la relación y
  reiniciar procesos permite consultar prepared sin aceptación ni outbox nuevo.
- Las cuatro fases action-loss preparan dos actas de fixture mediante checklist
  del owner. El navegador acepta una y rechaza otra por HTTP autorizado; intercepta
  y pierde cada respuesta DESPUÉS del200 del backend. Muestra incertidumbre,
  bloquea reenvíos y recupera por GET accepted/rejected. Exactamente dos POSTs,
  una excepción abierta y un evento por cada acta; versiones3 en PostgreSQL.
- Se comprueban ausencia de overflow horizontal en los viewports seleccionados,
  ayuda y estados. No es auditoría completa de accesibilidad/AT ni rendimiento.
- Issuer, cookies, TLS y datos son fixtures locales. Avisos de certificado
  autofirmado quedan en logs; no prueban IdP/PKI de producción. No llamadas de
  pago ni transporte externas. El cluster propio loopback se detuvo al terminar.
  Nuevo pool/reinicio de API/Next no es PITR, restauración de backup ni DR.

La ejecución anterior browser-final.log probó44 fases antes del fix de acciones.
La final browser-closure.log la supersede; no sumar ambas como92 pruebas finales.

| Receipt local | SHA256 |
|---|---|
| read-red.log | e57fb0b00832d3b7e865c538131718ca32687fb92b55f62310746aea6f726fb0 |
| exception-read-red.log (trigger rechazó UPDATE) | ef45716ffb78c9a7c676805af837a6d632da4b44a2211475d7e05027ba17c78d |
| exception-import-red.log | 06ca44fd120e934d5947c1d0ebde01b23c4132393cfeb0a4330eaed0e18a053c |
| action-red.log | da84dd3f85746fe0f47e8181bdf3010f3abf1c88f7efcf91e0ced471921a7877 |
| rebuilt-pg.log | 4db74e005833060ea1e838f818e70519369abe6a43e84e28050ac5c976fdb56a |
| rebuilt-http.log | 0b01ab332b02361999b9fd4615133aeccce060b050aa3ae42cffb97bc67a6ddc |
| go-test-closure.log | 8224001e397650e62bb703b3361fc36e088955843a95fb920a84692452f88c38 |
| web-test-final.log | 60b1eb4ca2dbd9026b94847924899599ef4162315c1a661017a5080087d4798e |
| web-build-final.log | 35bcbae8fe6000a79cb2fc4ba889357fac34f2528c4a358f8cfe212ce38981e1 |
| browser-closure.log | 7cda4e451e88c573f44b770d7653e9f8db5f0981c206f2cdd5c8050d6dd5c983 |

Vet/build terminaron exit0 con logs vacíos; su SHA256 común es
e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855.

| Archivo reconstruido | SHA256 |
|---|---|
| internal/platform/postgres/franchisejourney.go | e232c389373066c9a927e37f6e0ada063d417c096761957e5bfd8df2497d1ee0 |
| internal/platform/postgres/franchisejourney_integration_test.go | d152759014d8a9867b8ab6d76c0a7457bec48a46d9dd5639620a224c8323b860 |
| internal/platform/postgres/commerce.go | f50f64b51242bc71a5d92ccf87a7b2186440d865dc207ccecd8481ad7c00be4a |
| internal/platform/httpapi/franchisejourney_test.go | cd9a4e21340f779312ca3ea46774eaa73bffc35b5dc75e2dd04708fdb790f169 |
| src/app/customer/handovers/page.tsx | 9e889f9da062581429fb7c7e0903036b2967ec461a4ad8588e633b04d5882ea7 |
| src/components/customer-handover-actions.tsx | c1f264d9a83d58f7c34976e8479b107ef8b1d472357123ebd97c9b87f380e5af |
| microsoft_playwright_browser_gate/tests/enterprise-web.spec.mjs | 42f181dad8a9b9f89cf27a0025a88528eb81a2d5bddc14e2a2b7547110b3b49d |

Desde composición limpia, instalar locks de web y gate, construir web y aportar
TEST_DATABASE_URL de una base sintética loopback elite_confirmation_* con schema53:

```text
go test ./internal/platform/postgres -run '^TestFranchiseJourneyPersistenceIsolationAndReplay$' -count=3 -v
go test ./internal/platform/httpapi -run '^TestCustomerJourneyReadFailureIsRedactedAndRecoverable$' -count=1 -v
go test ./...
go vet ./...
go build ./...
go test ./internal/platform/httpapi -run '^TestQuoteAcceptanceBrowserPostgres$' -count=1 -timeout 22m -v
```

Para el último comando: ELITE_WEB_ROOT apunta a la composición web construida;
ELITE_QUOTE_E2E, ELITE_ORDER_E2E, ELITE_PAYMENT_E2E, ELITE_DELIVERY_READ_E2E y
ELITE_DELIVERY_ACTION_E2E valen1. Omitir ELITE_QUOTE_PROJECT ejecuta los cuatro.
El harness conserva artefactos por ejecución y no permite bases remotas.

## Continuidad y pendientes reales

FAIL450/451/452 corregidos y preservados como regresiones; FAIL432 conserva
recurrencia de diagnóstico. El UPDATE de excepción falló correctamente por su
trigger; se cambió sólo el fixture por INSERT FK-válido, sin desactivar guards.
FAIL384/385/405/423 no se cierran por esta evidencia ajena a sus claims.

Las actas creadas por INSERT en el harness son fixtures, NO la preparación
inicial de un pedido ordinario. Falta su comando/frontend y la regla comercial
que autoriza liberar, ya consultada sin respuesta registrada. No asumir pago
total, anticipo ni autorización manual. Falta además cobro/proveedor/callback/
reconciliación y cierre conjunto de stock/pedido/transporte. Conservar owners:
Journey acta/checklist/excepción; Commerce pedido/pago/stock; Fulfillment transporte.

No se cierra T2802/T2804 ni readiness/implementation_assurance integral. Pendientes
V295, capacitación/soporte/privacidad/observabilidad/carga/recovery/despliegue y
gates externos siguen en el mismo roadmap. No ZIP final ni promoción productiva.
