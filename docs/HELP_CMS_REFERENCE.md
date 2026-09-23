# CMS y guías de la misma revisión

Infraestructura local de biblioteca. HELP_CMS_ENABLED=true activa el host sólo
con migración0078,5guards y PostgreSQL18/UTF8/pg_unicode_fast. Deshabilitado lee
sólo la bandera. features.help_cms habilita /help/library y su BFF. Sin nuevas
cuentas/dependencias. help:read, help:write y help:publish exigen organización
explícita; wildcard conserva la autoridad ya existente, sin permisos heredados.

GO-HELP-CENTER-CORE conserva Store y12pruebas/16semillas. Se extraen funciones
puras de validación/version/draft→published→archived para el adapter durable.
El perfil CMS restringe edición a borradores. Cada transición guarda snapshot
inmutable, actor/comando/hash/version e incorpora outbox en una sola transacción.
Una respuesta perdida se recupera por GET del mismo actor/scope/hash; la UI
conserva sólo referencia/hash/versión bajo tenant+subject, nunca el texto.
El resultado histórico no habilita la siguiente transición sin lectura actual.

Crear/editar/publicar/archivar requiere acción explícita. Publicar no vuelve
público el artículo: sólo lectores de esa organización acceden al publicado.
Archivar retira también consultas a revisiones publicadas anteriores para
lectores; editores conservan historial. No edición silenciosa de publicados:
se prepara otro artículo y se revisa el reemplazo. Texto literal, sin HTML
ejecutable ni ingestión/clasificación de archivos empresariales.
Idioma es/en inmutable por artículo. Búsqueda substring con lower Unicode
pg_unicode_fast de PostgreSQL18.6 ya fijado; no afirma equivalencia exacta con
Go Search ni eliminación de acentos/normalización Unicode. No ICU/dependencia nueva.
Referencia oficial: https://www.postgresql.org/docs/18/collation.html
Contenido≤16KiB, comando codificado≤24KiB; bodyHTTP≤32KiB. Listas/historial50,
cursor por ID/versión, scope y búsqueda acotados; int64 como texto exacto.

Las6guías nuevas de catálogo/supply/garantía/red/CMS/capacitación comparten source
entre UI, /help y training_content/help.bundle.json. Las15anteriores son intactas.
Perfil reference-onboarding revisión2 fija SHA de21guías y5cursos, cada uno≤4
lecciones, mismos límites. Nuevas guías son AUTHORED instrucciones de operación
y extracción exacta de párrafos previamente publicados, no fuente corporativa.
No incorporación automática de artículos privados al curso ni a instrucciones
de IA; nueva versión de curso requiere revisión explícita del contenido.
Evaluación humana por otra persona y cero grants. Una activación futura conserva
el intento/assessment con su revisión original, sin mutación.

Pruebas: CMS1nuevo+3replay, rollback3tablas por outbox tardío, actor/org/tenant,
stale, permisos de publicación, inmutabilidad y retiro, paginación52artículos.
Browser realNext/BFF/Go/PG/Chromium conJWE/RS256/JWKS:6POST para2artículos/6
versiones/6eventos; crear y archivar pierden respuesta y recuperanGET trasreload,
0POST adicional. Reader no ve borradores/historial/retirados, textoHTML literal,
inglés/español por contenido, capturasdesktop390px.
5guards faltantes rechazan host, downgradevacío/reapply y rechazo poblado.
4goldens Go/TS,11webtests y7source/bundle/rolechecks; HTTP niega antesdelbody,
duplicados/oversize/versions. Fuzz2s3seeds328160ejecuciones. Next/tipos PASS.
Capacitación nueva:3guías→3lecturas→respuestas→revisor distinto;6facts6events,
1request1decision0resources/grants; futura revisión3 preserva historial2.

Owned PG con75migraciones; PAYMENT_CONNECTED_DB_URL prefixelite_payment_connected_
paraCMS; ELITE_TRAINING_DATABASE_URL prefixelite_training_connected_ para
TestTrainingReleaseGuidesConnected. Browser requiereELITE_HELP_CMS_BROWSER=1,
ELITE_WEB_ROOT/ELITE_NODE_BIN absolutos yNextbuild--webpack. Usa issuer sintético
de catálogo existente. Perfil web estrecho prueba imports/tipos, no Go remoto.
Usar node tools/export-training-content.mjs con outputausente bajo consumidor
para resolver la dependencia existente; conservar output/SHA después.

FAIL864 búsqueda C-locale corrigida con collationUnicode explícita;865output
del exporter fuera de árbol de módulos;866acceso índice testJSON;867currículo
superó4lecciones y se dividió en5cursos sin ampliar límites. REDs conservados.
HELP_CMS_RELEASE_V402.md/json contiene evidencia y reconstrucciones exactas.
T2804 aún requiereKPIs e i18n privado; resto del orden continúa. No READY global.
