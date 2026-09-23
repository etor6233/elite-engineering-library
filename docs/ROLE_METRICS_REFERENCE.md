# Indicadores conectados por rol

Infraestructura de biblioteca, local/fixtures; AUTHORED read-model/transport/UI
glue sobre los owners admitidos, sin nueva dependencia ni autoría corporativa.
GO-DASHBOARDS-CORE permanece alternativa aislada, no se promueve por estos tests.

El host electromobility integra RoleMetrics en EnterpriseQuery. El panel
/dashboard requiere features.role_workspace=true y sesión. Sus accesos a las
operaciones siguen visibleSections del owner por rol. El BFF sólo recibe kind,
organization_id y, para puntos/NPS, identificador de programa/encuesta.
Nunca recibe tenant, actor, SQL ni valores calculados por el cliente.

| Familia | Fuente | Permiso / alcance |
|---|---|---|
| orders, stock, cases, shipments | sales.customer_order, inventory.stock_unit, service_ops.service_case, logistics.shipment | admin:read + organización; envío origen o destino |
| own-orders, own-cases, own-appointments | pedido, caso vinculado a warranty propia, crm.appointment | customer:self + organización + subject |
| leads / appointments | crm.lead / crm.appointment | lead:read / appointment:manage + organización |
| factory-destination | factory.production_unit + procurement.purchase_order | factory:read + destino |
| factory-owned | factory.production_unit + serial_supply_plan | supply:factory-read + fábrica |
| supply | purchase_order + serial_supply_plan | supply:read + destino; excluye compras sin plan |
| stored-value | immutable operation + entry + account | stored_value:read + tenant/organización/programa del perfil activo |
| survey | Summary original de crm.survey_definition/response | surveys:read + organización; mínimo y retención originales |

Las métricas operativas son registros actuales por estado, no ventas mensuales
ni ingresos cobrados, utilidad o contabilidad. Importes se suman como numeric y
se transmiten como strings exactos por moneda y estado, nunca como JS number.
La suma supera int64 sin redondeo (wire hasta38dígitos); un resultado fuera de
la capacidad del contrato se declara no disponible. Sin conversión de monedas.
Los puntos se separan por programa/moneda y operación: issue/accrue/redeem/reverse.
Sólo se cuentan operaciones con entradas de ledger confirmadas, nunca propuestas
pendientes. No se convierten puntos a dinero ni se suman programas distintos.
NPS preserva el cálculo y mínimo existentes; cero respuestas no es NPS0.

Lecturas de operación/puntos usan transacción read-only repeatable-read,
statement_timeout3s; endpoint4s; NPS usa el SELECT original con deadline3s.
Cada consulta tiene su propio observed_at, no hay snapshot global entre familias.
Máximo100grupos operativos; el grupo101 produce409 y no truncamiento silencioso.
Wire puntos≤4operaciones. Autorización/servicio/perfil no disponible jamás se
presentan como cero. Cambiar selector borra resultados previos; generaciones
evitan que una respuesta tardía restaure datos de otra selección.
No-store en API/BFF. El BFF retenido usa transporte con deadline y bytes acotados.

FAIL868 Overview excluía won, estado ajeno al dominio. El cambio mínimo excluye
converted y lost; regresión RED→PASS preservada. FAIL457 del BFF genérico:
sesión/permisos antes del body, límite65536bytes, UTF8 fatal, deadline2.5s,
cancelación y permisos originales por acción. No nuevos grants ni cambios
de transición de negocio. No prueba de seguridad de producción.

Fixture config/role.metrics.fixture.json habilita el panel/encuestas para el
ensayo local; los flags de usuario siguen explícitos. Browser Chromium real
usa JWE+JWT/JWKS sintéticos, Go/PG y BFF; dos respuestas guardadas actualizan
NPS, distintas monedas/estado/importes exactos, roles, vacío y unavailable.
TestRoleMetricsConnectedScopeAndPrecision, StoredValueAndSurveyOwners,
FactorySourceScope y Browser prueban el delta; no se repiten recorridos previos.
FuzzMetricBoundary tiene seis semillas y presupuesto finito2s, no reemplaza SAST.
ROLE_METRICS_RELEASE_V402.md/json registra hashes/evidencia/reconstrucciones.
T2804 continúa con i18n privado. No READY global ni producción live.
