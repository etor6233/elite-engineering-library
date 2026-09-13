# V402 / checkpoint303 — mutaciones Mercado Libre conectadas

PROVEN_LOCAL para PRICE/STOCK/PAUSE/RESUME de un ítem User Products ya existente
y perteneciente al vendedor configurado. T2805 sigue en progreso: publicación
inicial, media/contenido, otros mappings y automatización pendiente no se
declaran completos. Tampoco TEST02/03/07, READY global ni producción live.

Perfil completo105packs/1379archivos; biblioteca200packs/2143blocks.
16archivos AUTHORED nuevos y4salidas revisadas.1359salidas previas permanecen
exactas; no cambia ningún frontend ni owner de pago/handover/entrega anterior.
GO_HUMAN_APPROVAL_CORE0.8.0 agrega un tipo manual; aplicación1.20.0 agrega su
montaje opcional. GO_CONNECTED_MARKETPLACE_MUTATION0.1.0 contiene el binding,
contratos, migración0079, host, fixture/configuración y pruebas reproducibles.

Una publicación real local se crea con los owners de catálogo existentes:
modelo/variante/libro, PNG, snapshot, cuatro revisiones y publicación inmutable.
Sus importes y el ATP existente alimentan el payload que revisa otra persona.
El cliente no suministra otro precio/stock/regla; el guard reconstruye el
payload y exige revisión válida antes de reclamar el fence compartido.

Prueba PG+HTTP:6aprobaciones/6fences,12envíos concurrentes→1PUT;5intentos de
proveedor y4efectos. Precio9007199254740991minor exacto. ATP0 retira el stock
sintético4; los tres modos de stock y cantidades positivas se prueban en
contrato. Pausa/reactivación reales del fixture; pérdida de respuesta del
proveedor y del API con recuperación sin duplicados. Rechazo y drift de
x-version son terminales. Estado/aprobación/evidencia conservados e inmutables.

JSONB se canonicaliza antes de comparar/enviar; no se debilita el hash aprobado.
Un User Product distinto del ítem observado no puede recibir la escritura.
HTTP200 con precio ignorado no confirma el efecto; automatización de precio,
seller/categoría/currency/location inesperados y precios ambiguos fallan cerrado.
La API de escritura prices/standard aún no disponible no fue seleccionada.

Host con perfil/hash/catálogo/token/HMAC/migración obligatorios PASS. Revertir
una base con historial se rechaza y mantiene6/6/6registros; base vacía down/up
PASS. Fuzzing focal2s, vet y build PASS. Cuatro perfiles reconstruidos exactos;
compilación de los dos perfiles reducidos PASS sin ejecutar suites anteriores.
Los RED de selección de fixtures/owners, SQL, representación JSONB y packaging
quedan conservados y corregidos. G0–G8 y todos los receipts están en el JSON.

Fuentes: siete snapshots HTTPS oficiales, con URL/SHA. Ningún SDK archivado
se adopta ni código local se atribuye a una empresa. Los pins/dependencias y
159bloques ADAPTED/141VERBATIM previos permanecen intactos. Esta evidencia
confirma estado observado; no promete atomicidad distribuida ni equivalencia
de imágenes remotas. Ver docs/MARKETPLACE_MUTATION_REFERENCE.md materializado.
