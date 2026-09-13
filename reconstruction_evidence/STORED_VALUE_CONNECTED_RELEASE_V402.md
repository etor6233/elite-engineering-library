# V402 — gift/loyalty conectado, empaquetado y reconstruido

T2802, claim PROVEN_LOCAL acotado: gift parcial o descuento loyalty reducen el cobro
real del SDK contra fixture; cobertura total produce un recibo local inmutable y
conecta handover/aceptación/recibo comercial con cero intentos de pago proveedor.
Operador y revisor distintos actúan desde el portal real. No cierra todo T2802,
TEST02 ni READY global. Canonical180packs/1877bloques; franquicia84packs/1112archivos.

Fuentes: Odoo Community19.0, commit99edb6dd82b7b560930c00b03b694ba700785370.
Veinticuatro archivos completos VERBATIM y las dos copias de LICENSE/COPYRIGHT
mantienen bytes/Gitblob/SHA. engine.py y expectativas de tests son ADAPTED LGPL3;
los segmentos exactos y cambios Decimal/representación están en DERIVATION*.json.
El runtime aislado y su glue Python se distribuyen con fuente correspondiente LGPL;
Go/SQL/HTTP/TS son AUTHORED glue sobre los owners existentes, sin atribución ajena.
La licencia oficial43529bytes conserva su SHA sin LF final gracias al compositor0.3.0.

El perfil de referencia ARS fija reglas de gift/loyalty y sus límites explícitamente.
Se trata de cuentas internas asistidas: no se afirma emisión de códigos bearer,
eWallet público, fiscalidad de gift, nómina, tax mapping ni ERP completo. El origen
confirmado de emisión es un fixture declarado; la compra de destino usa el recorrido
quote/accept/order/stock. No se inventa una venta pública de gift cards.

Los efectos de aprobación son inmutables y ligan actor/tenant/org/pedido/versión,
perfil, fuente, cuentas, cálculo y lease. Locks/reservas evitan doble gasto; la
expiración se comprueba al COMMIT después de esperas. Una inversa usa el pedido
cancelado y estado seguro del pago; conserva historia y el comportamiento negativo
de saldo derivado cuando la emisión ya fue gastada. El funding total toma exactamente
las contribuciones aprobadas, produce un receipt auditable e idempotente y reserva
un vínculo físico para entrega. No crea dinero ni un pago proveedor falso.

## Pruebas y reconstrucción

- Comparación focal fuente/oráculo y expectativas oficiales adaptadas:11tests y240
  comparaciones previas reutilizadas por identidad de cuatro AST; no suite ERP.
- PostgreSQL65migraciones: gift5000 sobre123456 cobra118456; loyalty6173 cobra117283.
  SDK/callback/reconciliación y commercial receipt conectados; legado perfil1/2
  rechaza monto ajustado y perfil3 debe seleccionarse expresamente.
- Cobertura total:12concurrentes producen un funding/efecto; recovery idempotente,
  actor/payload divergentes y timing/expiración se rechazan sin efecto parcial.
- HTTP/BFF valida scope/permiso y representación int64;9007199254740993 viaja exacto.
 49tests focales, typecheck/build Next y navegador conectado PASS. El navegador
  usa OIDC/JWKS y sesiones cifradas locales, revisores separados y respuestas perdidas
  recuperadas desde request keys hasta entrega aceptada/recibo comercial.
- Loader: dos RED reales de initializer/cache corregidos; fuente alterada rechazada,
  reemplazo explícito soportado y fuzz local con presupuestos finitos PASS.
- Downgrade con historia se rechaza sin pérdida; dependencia externa no se borra;
  base vacía down/up PASS sin CASCADE. Host activation, vet y build PASS.
- Cinco perfiles reconstruidos byte por byte:84/1112,40/567,8/163,57/793 y85/1120.
  Backend/serverless compilan owners/test packages; cero tests de runtime ejecutados
  por ese chequeo. Reconstrucción final del perfil completo liga metadata admitida.
  No se repitieron suites sin delta ni se sumó compilación como prueba de negocio.

## G0–G8 del claim local y límites

| Gate | Evidencia / condición de uso |
|---|---|
| G0 identidad | Perfil previo y24recibos oficiales con commit/blob/SHA; snapshots completos de archivos seleccionados, no firma/full archive inventados |
| G1 licencia | Fuente LGPL correspondiente, textos exactos, modificaciones y derechos de reemplazo; notices materializados. IPC no decide por sí solo el alcance legal |
| G2 autoridades | Cálculo comercial deriva de métodos identificados; owners de precisión/PG/aprobación existentes gobiernan el glue |
| G3 arquitectura | Un owner de aprobación/outbox; proyección funding común y FKs reales; provider y funding XOR; tenant/org/permisos en cada frontera |
| G4 correctitud local | Oráculo y escenarios PG/HTTP/browser de arriba, expiración al commit, recovery y negativos materiales |
| G5 seguridad del delta | Source closure/hash, sin imports extra; autorización objeto/acción, revisión separada, límites de cuerpo/proceso y JSON exacto probados. No sustituye SCA/seguridad de composición T2803 |
| G6 resiliencia local | Locks/reservas, concurrencia12→1, respuesta perdida, rollback/timeout y presupuesto de fuente acotado. No claim p95/capacidad target; carga compuesta en T2809 |
| G7 evolución local | Runbook, source replacement, perfil versionado, historia inmutable y down/up seguro; ops/retención/backup comunes conservan sus owners |
| G8 reconstrucción | Dos packs nuevos + siete owners actualizados,73archivos nuevos; cinco perfiles exactos y fuente final idéntica a la probada |

Estado de los packs: REBUILD_VERIFIED / CONDITIONED para el perfil asistido y
hash-bound descrito. El claim local ya tiene código y evidencia; no se disfraza
trabajo global de seguridad/ops como falta de credenciales. T2803/T2809 y target
business/live permanecen explícitos. No se declara ELITE_REFERENCE por reputación.
El módulo no promociona los antiguos cores aislados GO-GIFT-CARDS-CORE/GO-LOYALTY-CORE.

Evidencia detallada previa: STORED_VALUE_TENDER_V402.md,
STORED_VALUE_OPERATOR_FLOW_V402.md y STORED_VALUE_SOURCE_CLOSURE_V402.md.
Manifest fuente/delta/recibos y hashes: STORED_VALUE_CONNECTED_RELEASE_V402.json.
Referencia externa: C:/Users/NL/AppData/Local/Temp/elite-v402-library-infra/gift-tender-admitted-reference.
El destino durable y release firmado continúan en T2810; no se presentan como hechos.
