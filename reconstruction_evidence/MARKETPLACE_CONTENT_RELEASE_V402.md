# V402 / checkpoint305 — contenido Mercado Libre conectado

PROVEN_LOCAL: catálogo aprobado y PNG aceptado→CONTENT exacto→aprobación
distinta→un PUT→GET de confirmación/reconciliación. Precio/stock/status no
forman parte del payload. T2805 queda abierto para otros feeds/comunicaciones;
TEST02/03/07 y READY global no se cierran por esta evidencia.

Mismo owner GO_CONNECTED_MARKETPLACE_MUTATION0.3.0,31archivos.7nuevos y7revisados
AUTHORED;1380salidas previas exactas. Biblioteca200packs/2158blocks; perfil
105/1394 y perfilHTTP106/1402 reconstruidos exactos.159ADAPTED/141VERBATIM y
todos los pins/dependencias originales intactos. Ningún frontend ni algoritmo
de catálogo/price/ATP/aprobación/fence original cambió.

El contrato oficial User Products propaga family_name a sus condiciones de
venta y exige ausencia de ventas. La composición exige un único ítem asociado
observado y sold_quantity explícito cero antes de PUT. Ventas, campo ausente,
varios ítems, owner/UP distinto o cambios de precondición se rechazan.
No se asume permiso para afectar una familia ajena o crear política comercial.

La prueba real source/API/PG conserva el catálogo normalizado/revisado y una
carga MEDIA aceptada.3aprobaciones/3fences/3intentos;1multipart y1CONTENT PUT.
Un caso multi-item se vuelve failed_terminal sin emitir PUT. Otro pierde la
respuesta después del efecto y se recupera sólo por GET;12replays no escriben.
El nombre e imagen anteriores cambian al snapshot aprobado y picture_id
aceptado. El precio9007199254740993minor y stock4 del proveedor no cambian.
36POST del API incluyen preparación/aprobación/negativos/replays.

Ocho escenarios de contrato: efecto, respuesta perdida, multi-item, ventas,
ventas desconocidas, imagen incorrecta en reconciliación, rechazo y payload
monetario extra. Tenant/org/permiso y título inventado también se rechazan
desde el API. No igualdad ficticia entre PNG original e imagen transformada.

Build/vet/fuzz2s/host PASS. Migración0081 amplía el constraint inmutable;
down con CONTENT histórico falla y conserva3/3/3registros, empty down/up PASS.
Los dos procesos PostgreSQL usados terminaron. FAIL894 fue una expectativa
errónea del fixture (failed frente a failed_terminal); el owner de fence
actuó correctamente y no se modificó. RED y prueba corregida preservados.

Fuentes oficiales: los mismos ocho URL/SHA de la revisión304. G0–G8 y receipts
están en el JSON. El delta es AUTHORED glue sin atribución empresarial.
El efecto observado no promete una transacción distribuida con cambios
concurrentes hechos fuera de esta aplicación.
