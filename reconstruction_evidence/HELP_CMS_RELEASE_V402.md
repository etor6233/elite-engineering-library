# V402 / execution300 — CMS y guías conectados

PROVEN_LOCAL para CMS de ayuda y capacitación de la misma revisión bajoT2804.
102packs/1324archivos exactos; biblioteca197packs/2088blocks. T2804 continúa
conKPIs/i18n privado; todavía no se declara infraestructuraREADY ni TEST02/03/07.

Se reutiliza el kernel AUTHORED GO-HELP-CENTER-CORE, extrayendo funciones puras
y conservando Store,12tests y16semillas. El adapter PostgreSQL publica snapshots
inmutables con actor/comando/hash/versión y outbox atómicos. Edita sólo borradores;
archivar retira también versiones publicadas anteriores a lectores. Historial
permanece para editores. Permisos help:read/write/publish y organización exacta,
sin permisos automáticos ni publicación pública. Texto literal, no ingestión
de documentos empresariales ni ejecución HTML.

CMS:1nuevo+3replay concurrentes, rollback3tablas por fallo tardío, límites,
tenant/org/actor/versión/permisos/historial y paginación52artículos probados.
Búsqueda Unicode explícita de PostgreSQL18.6 fijado, es/en por artículo.
Browser realNext/BFF/Go/PG/Chromium:6POST2artículos6versiones6eventos.
Crear y archivar pierden respuesta y recuperanGET trasreload sinPOSTadicional.
Reader/foreignorg yHTMLliteral comprobados.5.53s ycapturasdesktop390px inspeccionadas.

Las6guías nuevas de catálogo/supply/garantía/red/CMS/capacitación comparten source
enUI, /help ybundle. Las15anteriores permanecen iguales.21guías/5cursos fijados
porSHA enreference-onboarding revisión2, máximo4lecciones porcurso. Nuevo curso
de3guías recorrió3lecturas→respuestas→revisión por otra persona,6facts6events,
1request1decision0grants/resources. Activación futura preserva intento/assessment
de revisión2. CMS privado no ingresa automáticamente a capacitación ni a IA.

Host:5guards,UTF8/pg_unicode_fast; cada ausencia bloquea. Downgradevacío/reapply
PASS y poblado rechaza pérdida de historia. Deshabilitado sólo lee bandera.
11testsweb/7alignment checks/4goldensGoTS/HTTPboundaries/Next/tipos y perfil
materializado PASS. Fuzz2s3semillas328160ejecuciones. Sin nuevas dependencias.

FAIL864 encontró lower() C-locale para búsqueda conacentos: fix explícitoUnicode.
FAIL865 exportaba fuera del árbol de módulos: outputausente dentro delconsumidor.
FAIL866 acceso sin guard a índices de goldenJSON; FAIL867 currículoadmin excedía
4lecciones: se dividió en2cursos manteniendo límites. REDs conservados; no se
relajaron asserts/guards ni repitieron recorridos anteriores sin delta.

Dos companionpacks nuevos,10owners revisados:23blocks nuevosAUTHORED y2archivos
existentes de kernel ahora seleccionados. Cinco rebuildsexactos:102/1324full,
15/231web,40/570backend,64/864serverless,103/1332HTTP; imports/tiposestrechosPASS.
Fuentes y receipts detallados enHELP_CMS_RELEASE_V402.json ydocs/HELP_CMS_REFERENCE.
No reputación corporativa como prueba, no producción ni cuentaslive.
