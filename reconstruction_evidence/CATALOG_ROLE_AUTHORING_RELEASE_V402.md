# V402 / execution296 — autoría de catálogo por rol

PROVEN_LOCAL en LIBRARY_INFRASTRUCTURE para alta inicial→PNG→snapshot→cuatro
revisiones humanas→publicación→storefront y rollback.93packs/1248archivos exactos;
189packs/2014blocks en biblioteca. GO-CONNECTED-CATALOG-AUTHORING0.1.0 y
TS-CATALOG-AUTHORING-PORTAL0.1.0:10archivos cada uno, seis owners revisados.

Se cerró la dependencia de SQL fixture para la variante inicial: el formulario
crea un modelo con1–16variantes y precios explícitos usando los services y SQL
originales de Electromobility/Commerce en una transacción. La variante usa la
tabla existente y el comando se registra con actor/hash/referencias/outbox.
El alta permanece draft/no pública/homologation unknown; no otorga aprobación.
El formulario de snapshot compone un modelo/libro; no se afirma editor masivo
multi-modelo ni edición de snapshots aprobados.

La prueba real inicia sólo tenant/org y opera mediante Next/BFF/Go/PG y Chromium
con JWE/RS256/JWKS. Produce2modelos/2variantes/5libros/1media/2snapshots/
8revisiones/3publicaciones/16comandos. La tercera publicación restaura la primera
versión con otro libro efectivo; nunca reactiva el antiguo. Se pierden respuestas
de alta y publicación y se recuperan por GET, incluso recargando la página,
sin otro POST. La consulta vigente privada rechaza otra organización; el creador
no revisa ni publica su snapshot.

Cuatro llamadas simultáneas a la misma alta producen1nuevo3replays. Fallo real
de outbox revierte seis tablas; actor/hash/org y referencias sobreviven reinicio.
14pruebas frontend focales y4goldens Go/TS comprueban comandos, int64 como texto,
Unicode/escapes y recuperación ligada a actor/hash/snapshot. Bodies JSON32KiB,
PNG1MiB y deadline/cancelación, sin cambiar los límites textuales anteriores.
Next --webpack/strict types PASS. Fuzz2s/3semillas/112321ejecuciones.
Escritorio y390px inspeccionados, sin desborde horizontal.

El host requiere migration0076 además de sus guards previos. Down con
16comandos/3publicaciones/2modelos se rechaza conservándolos;73migraciones en
base vacía y76down/up PASS. Cinco perfiles reconstruidos exactos:
93/1248franquicia,40/570backend,11/190web,60/823serverless,94/1256HTTP.
Los perfiles menores compilan su cierre de imports/tipos.

FAIL853: inferencia UUID/text de fixture, antes de escrituras; corregida.
FAIL854: tipos RequestInit/duplex del test, con14checks ya PASS; corregidos
sin relajar TypeScript. FAIL855: label de textarea ya poblada; rol/nombre del
árbol accesible conserva la aserción. Los RED y fixes están fijados en el JSON.
La prueba final de navegador completa pasó en6.95segundos.

20archivos nuevos y10revisados son AUTHORED glue declarado; no se agrega una
dependencia, fuente corporativa o algoritmo de negocio alternativo.
CATALOG_ROLE_AUTHORING_RELEASE_V402.json conserva G0–G8, hashes y receipts.
Guía materializada: docs/CATALOG_ROLE_AUTHORING_REFERENCE.md.

T2804 continúa con supply/warranty/J5/CMS/guías/KPIs/i18n privado.
No cierre de TEST02/03/07 ni promoción45/48 por este recorrido.
