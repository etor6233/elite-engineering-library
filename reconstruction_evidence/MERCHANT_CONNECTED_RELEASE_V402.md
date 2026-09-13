# V402 / checkpoint306 — Google Merchant conectado

PROVEN_LOCAL para publicación/upsert y renovación manual del catálogo aprobado
en una fuente API primaria dedicada.106packs/1418archivos exactos. T2805 sigue
abierto para sus mappings realmente seleccionados y comunicaciones; TEST02/03/07
y READY global no se cierran por este tramo.

Fuente revisada y precio entero original/ATP→solicitud tipada→aprobación distinta
→SDK Google1.8.0 fijado→recibo raw PostgreSQL→GET de procesamiento/reconciliación.
Version_number se liga a generación local; oferta/cuenta/fuente/idioma/feed label,
imagen y descripción se fijan por perfil/snapshot. No redondeo ni reglas nuevas
de precio/impuesto/inventario. Tres aprobaciones/efectos/intentos;cinco observaciones.
Tres INSERT reales del SDK ycuatro GET;12 replays no repiten INSERT.
Pérdida de respuesta del proveedor y del API, pending404 y data_source ajeno
probados. Un GET exacto prueba el estado observado, nunca renueva por inferencia
el plazo de frescura. Product_status pendiente se conserva, sin vender aprobación.

La cola autenticada/scoped deriva aprobación/envío/reconciliación y refresh_due_at
de las evidencias originales. Sólo un acknowledgement INSERT fija ese plazo.
No reenvía solicitudes vencidas ni crea trabajos por consultar la cola.
Montaje host, hashes, migraciones, límites, vet/build y fuzz2s PASS.
Downgrade con evidencia se rechaza y conserva3/3/3;emptydown/up PASS.
Los PostgreSQL propios terminaron. Los otros journeys conservan1393salidas exactas.

Instalador portable offline:20wheels originales,1041payloads instalados verificados,
pip check y runtime config hash. Prueba SDK sobre ese segundo entorno:
6POST/3GET,401 sin reintento,307 bloqueado y respuesta grande rechazada.
OSV2.5.1 cubre20/20identidades PyPI, cero hallazgos. Su normalización PEP440
produjo el falso desacuerdo del comparador FAIL896; se corrigió reutilizando
el scan original. FAIL895 fixture del query $alt y FAIL897 UTF8 del generador
de receipts quedan preservados y corregidos; no se modificó el SDK.

Tres snapshots oficiales URL/SHA, commit Google97d7b42cd74b41211f5ec8871cc0dd15debdb1a0
y wheelSHA722ef095eca35126255c129964a61818f0d7bc9eeef7f1da12732d81c1dcd505.
SDK Google DEPENDENCY_PIN Apache-2.0; su grafo conserva MIT/MPL/BSD/PSF y
los archivos de licencia de los wheels. Los24archivos nuevos son glue AUTHORED.
El perfil ilustrativo no aporta cuentas. Infra local completa para el claim;
credencial/account access del usuario se aportan cuando decida activar.
