# V402 — J2 conectado, backend de referencia

PROVEN_LOCAL acotado al backend en staging, no publicación canónica ni T2802 completo.
Base canónica86/1148 intacta;14archivos de delta congelados y hasheados en el JSON.

PostgreSQL18.6 con69migraciones, Go1.26.8 exacto, dependencias offline: PASS.
Plan real sobre PO, dos variantes/tres unidades; se registra una cuarta como
reemplazo de una rechazada. Confirmación de proveedor es registro autorizado
con evidencia de fixture; no autenticación de API proveedor inventada.

El owner existente prepara estado/ID/version y escribe PO/fábrica/stock. Seis
helpers conservan su SQL y wrappers; el nuevo glue liga cantidades, series,
dos ASN y recepciones parciales. Se mantiene un solo stock y ATP/reserva originales.
No se calcula precio/costo nuevo ni se modifica el total comprado.

Cuatro registros de fábrica, uno rechazado; tres manifestados/recibidos/disponibles.
Cuarentena no cuenta en ATP; liberación requiere una persona diferente mediante
la aprobación compartida serial_quality, nunca automática. Rechazo conserva
cuarentena; reinspección exige evidencia nueva y otra decisión inmutable.

PASS: sobrecantidad, envío prematuro, organización/rol ajenos, autoaprobación,
bypass de API genérica, ASN incorrecto, actor/payload de replay divergentes.
Fallo de outbox revierte estados, stock, aprobaciones y toda evidencia comparada.
Ocho llamadas idénticas: un efecto nuevo y siete replays. Un pool nuevo recupera
actor/hash/version30 y evidencia original; no reenvío escondido. Fuente/tabla
de cantidades inmutable. PostgreSQL detenido; cero SKIP.

La política strict-serial-reference/v1 fija cero sobrerecepción y cuarentena
hasta revisión. Es fixture declarado, no legislación ni homologación de producto.
El código nuevo es AUTHORED binding/transport/configuration glue; source owners
BC derivados existentes conservan sus atribuciones estrechas. Ninguna empresa
se atribuye el workflow local.

Faltan HTTP/host, negativos de esas fronteras, downgrade/fuzz y publicación exacta
de este mismo código; luego J3. No repetir la prueba del defecto VIN de0071 ni
warranty/FX/pagos sin delta. Receipts y14hashes:
SERIAL_SUPPLY_CONNECTED_BACKEND_V402.json.
