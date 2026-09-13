# V402 / checkpoint304 — publicación inicial Mercado Libre conectada

PROVEN_LOCAL para catálogo aprobado→PNG normalizado→aprobación distinta→carga
multipart→nueva aprobación→publicación User Products→confirmación/reconciliación.
T2805 permanece abierto para contenido existente, otros feeds y automatización.
No se cierra TEST02/03/07, READY global ni producción por esta evidencia.

Perfil105/1387; biblioteca200packs/2151blocks. El mismo pack de Mercado Libre
sube a0.2.0:8archivos nuevos,9revisados,1370salidas anteriores exactas.
Ningún frontend, amount/ATP/catalog writer, aprobación o fence original cambia.
La migración0080 agrega observaciones inmutables y tipo de efecto. Los dos
perfiles que lo seleccionan reconstruyen exactamente105/1387 y106/1395.
Los perfiles reducidos no cambiaron ni se vuelven a probar.

La prueba usa el catálogo real local: modelo/variante/libro, PNG normalizado,
snapshot, cuatro revisiones y publicación. La API prepara desde ese snapshot
y el ATP original. Una persona distinta aprueba el upload y después CREATE,
que exige un picture_id observado duramente y aceptado para la misma variante
y generación. El cliente no aporta imagen, precio ni stock alternativos.

PG/HTTP:5solicitudes,4efectos/fences,4intentos totales,2observaciones durables;
12clics concurrentes hacen1upload. En total3multipart y1POST/items: la primera
imagen se confirma, otra pierde su respuesta y quedaUNKNOWN, y una tercera
requiere aprobación independiente sin borrar el incidente anterior.
CREATE pierde la respuesta del proveedor; GET/search→item/prices/stock confirma
exactamente precio9007199254740991minor, ATP0/paused e imagen. Repetir send no
crea otro ítem. Incluso con búsqueda temporalmente vacía, el histórico local
aceptado impide un segundo CREATE del SKU.48POST del API incluyen preparación,
aprobación, negativos y concurrencia; no son48efectos del proveedor.

Negativos: tenant/org/permiso, SHA de PNG alterado, autoaprobación, envío sin
aprobación, media inexistente, respuesta con imagen incorrecta, búsqueda
ambigua, provider rejection, persistencia fallida y cambios/borrado de evidencia.
Siete escenarios de contrato; build/vet/fuzz2s PASS. Host rechaza guards de
efecto/observación ausentes. Down0080 con datos falla y conserva5/4/4registros;
base vacía down/up PASS. PG se detuvo al finalizar.

La respuesta de carga se fija al hash exacto del PNG enviado. Mercado Libre
transforma imágenes: no se afirma igualdad de bytes con el JPEG remoto.
Sin respuesta ni ID de upload durable no existe recuperación por SHA admitida:
UNKNOWN permanece, no se reenvía ni se inventa equivalencia. Otro upload requiere
otra aprobación; un PNG huérfano no crea una publicación comercial.

Ocho documentos oficiales fijados por URL/SHA. El nuevo contrato de imágenes:
https://developers.mercadolibre.com.ar/es_ar/saldo-de-la-cuenta/trabajar-con-imagenes
SHA2562e85a485a52fa1311b7e8aaa1391fb18df140dc16b880dd723b5705f43fbe695.
G0–G8/receipts en el JSON. Todo el delta es AUTHORED glue, sin atribución a
Mercado Libre ni otra empresa.159ADAPTED/141VERBATIM y dependencias intactos.
FAIL892 corrige el límite MIME antes de cualquier envío; RED preservado.
