# V402 — J3 catálogo conectado e incorporado

PROVEN_LOCAL para LIBRARY_INFRASTRUCTURE: publicación aprobada→precio efectivo→
search/caché→storefront real→feed de referencia con fence/reconciliación.
GO-CONNECTED-CATALOG-PUBLICATION0.1.0 (28archivos) y
TYPESCRIPT-PUBLISHED-CATALOG-STOREFRONT0.1.0 (8archivos);11archivos de cinco
owners existentes revisados. Biblioteca185packs/1968blocks
(AUTHORED1668/ADAPTED159/VERBATIM141). Franquicia89/1202exacta.

Tres snapshots inmutables, cuatro revisiones humanas compartidas por snapshot,
cinco publicaciones incluyendo rollback. Cada publicación crea un libro efectivo
nuevo por el writer Commerce original y lo activa una sola vez; conserva
created1/activated2, referencias de cotizaciones y snapshot de importes/vigencia.
La primera reactivación falló correctamente: FAIL845 conservado y corregido,
sin cambiar unicidad de eventos ni atribuir una versión inventada al upstream.

PG18.6/Go1.26.8/71migraciones: fuente editada no muta publicación; precio, búsqueda,
storefront y feed coinciden. Ocho publicaciones concurrentes producen1efecto/
7replays; outbox fallido revierte13tablas; rollback vencido rechazado; receipt
recuperado con otro pool. HTTP pierde respuesta confirmada y recupera por GET
actor/hash sin segundo POST aceptado. Tres ETags obsoletos cambian; headers
200/304 mantienen revalidación y privados no-store. PNG estructural por Go
image/png fijado, limitado y normalizado; no antivirus ni nativos/libxml2.

Next16.3.4/Node24.20.0:33tests de esquema/SEO y build con tipos PASS. Chromium
navega cuatro veces contra esas mismas publicaciones1/2/3/5: detalle real,
nombre/precio/imagen exactos, canonical, robots, sitemap y404 para modelo/medio
no publicado. React conserva un snapshot por request y fetch no-store evita
retener una versión anterior. Modelo actual ligado al tenant configurado.
Reutiliza LeadForm y presentación exacta de importes;0dependencias añadidas.

Feed referencia usa OutboundDelivery existente:4POST, rechazo terminal,
aceptación con respuesta perdida, fence sin otro POST y1GET de reconciliación.
Intención inmutable enlaza actor, perfil, generación, source/payload SHA.
Es un protocolo de receptor explícito, no SDK live Google/Meta/ML; ese mapping
permanece T2805. No se etiqueta credencial-only mientras falte ese código.

Host exige12guards de publicación más guard feed, perfiles SHA,32bytes de key
y token; apagado no lee perfiles y activación no hace llamadas externas.
Down73/74 poblado rechazado conservando3drafts/5publicaciones/12reviews/4intentos;
vacío74→73down y73→74up PASS. Fuzz finito2s/7semillas/121059ejecuciones PASS.
Cinco reconstrucciones exactas; narrow Go/web compilan, sin repetir suites
runtime anteriores. FAIL843–849 corregidos, registros rojos conservados.
FAIL842 es preparación contenida. Historial V293 no se reescribe.

G0: Existing admitted Catalog, Commerce/BC predicate, Approval, Search, OutboundDelivery and Next/React/Zod/Playwright selected; no new upstream acquisition.

G1: 36new and11changed AUTHORED glue blocks. Existing source-derived price/inventory/calculation owners and third-party notices unchanged.

G2: Required blueprint J3 binds one explicit tenant/organization/market/currency/origin and4shared review stages; no implicit fiscal/legal policy.

G3: Immutable snapshot; original Commerce transaction writers create fresh effective books; approval/search/outbox and existing send fence remain single owners.

G4: Actual PG18.6/71migrations,3snapshots/12reviews/5publications; source isolation; price/storefront/search/feed match; rollback/expiry/receipt recovery.

G5: Scoped HTTP,1new7replays,13table outbox rollback, generic bypass denied; lost commit response recovered by actor/hash GET; terminal/unknown feed never automatically reposted.

G6: 33schema/SEO tests, Next build and actual Chromium at generations1/2/3/5; canonical/sitemap/robots/PNG/404. Finite2s fuzz7seeds/121059execs.

G7: Hash-bound optional host;12base plus feed immutable guard, exact32byte key/token required; no provider startup call; populated73/74down refuses intact, empty reverse/forward order passes.

G8: Five exact reconstructions89/1202,40/570,9/171,58/804,90/1210. Narrow Go/type closure compiles; unchanged broad runtime suites not repeated.

Hashes, receipts y delta: CATALOG_CONNECTED_RELEASE_V402.json.
Guías materializadas: docs/CATALOG_PUBLICATION_REFERENCE.md y
docs/CATALOG_STOREFRONT_REFERENCE.md. No producción ni cierre global por J3.
