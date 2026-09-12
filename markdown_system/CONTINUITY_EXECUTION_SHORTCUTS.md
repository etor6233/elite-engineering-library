# Continuidad: avanzar con evidencia y evitar trabajo repetido

Procedimiento AUTHORED de mantenimiento, consolidado en V375 por instrucción del
usuario. No crea otro roadmap ni reduce sus criterios de aceptación.

1. Validar el checkpoint y los hashes antes de reanudar. Definir el resultado
   observable del siguiente paso y el control al que aporta. Un avance parcial
   no cierra el control completo; conservar sus blockers y oráculos.
2. Reutilizar fuentes, ejecutables, caches, entornos y suites ya admitidos cuando
   sus hashes y condiciones siguen siendo compatibles. Repetir sólo por delta
   material, fallo, cambio de entorno o duda concreta. Un binario encontrado en
   una carpeta temporal no queda admitido por existir.
3. Enumerar rutas y leer el manifest antes de construir un comando. No deducir
   nombres de locks, migraciones, owners o archivos internos de otra versión.
   Procesar salida estructurada; anclar los parsers de diagnósticos a su campo.
4. Usar apply_patch para scripts con varias capas de comillas, JSON o heredocs.
   Mantener texto UTF-8/LF explícito. Al copiar un cache Go, preservar el original
   y hacer escribibles únicamente los manifests de la copia que se adaptará.
5. Calcular primero los límites de una prueba costosa. En una alerta, la tasa
   media no acota cada ventana deslizante; contar el máximo de eventos por
   ventana, considerar extrapolación y preservar threshold, hold y evaluación.
   Ejecutar la validación rápida del formato y los casos negativos antes del
   ensayo largo. No sustituir un tiempo contractual por un reloj acelerado.
6. Reutilizar el harness anterior y su entorno exacto, incluidos COMSPEC,
   toolchain, aislamiento, ownership y limpieza. Cambiar una frontera por vez.
   Registrar exit code y salida por fase, también cuando la salida esté vacía.
7. Durante una espera material, avanzar con verificaciones independientes:
   licencias, grafo exacto, SCA, paridad, documentación del alcance y revisión.
   No modificar el artefacto bajo prueba ni registrar PASS antes del resultado.
8. Devolver el arreglo a los bloques canónicos, reconstruir en destino vacío y
   comprobar paridad. Ejecutar los gates afectados y un Preflight integral sobre
   la revisión final. Conservar fallos originales; portar evidencia sin paths
   personales, secretos ni artefactos privados.
9. Mantener un encabezado actual y un próximo paso único en el cursor. Las
   entradas anteriores son historia fechada, no instrucciones vigentes. No
   reiniciar trabajo ya cerrado porque un párrafo antiguo decía “pendiente”.

Casos que sustentan la secuencia: V372 (runtime y telemetría reales), V373
(pipeline y rollback), V374 (regla de poco tráfico y promoción retractada),
V375 (instrumentación HTTP, reutilización de fuentes y límites de ventana).
Los resultados pendientes de V375 sólo se acreditan cuando su informe conserve
la ejecución y sus hashes; este procedimiento no los anticipa.

FAIL726: reconstruir los paquetes completos del perfil usado para compilar. La paridad de un subconjunto de archivos no demuestra igualdad de entradas del compilador. La reconstrucción integral V375 produjo el mismo SHA del binario ensayado.

Antes del primer Preflight de un pack nuevo, verificar también los diez encabezados del template, no sólo sus bloques. VERIFY_LIBRARY valida el cursor local: actualizar checkpoint después de cualquier cambio de owner o ledger y antes de ese gate. Mantener el reporte de una ejecución fallida y usar otro destino para su sucesor.

V376: distinguir el grafo declarado de paquetes del contenido vendorizado y nativo. Usar perfiles/hashes existentes y repetir sólo los gates afectados; para cambiar dependencias del BFF, ejecutar los recorridos conectados sobre una base descartable nueva. Derivar versiones y presupuestos de los planes que seleccionan realmente el owner. Un arreglo de advisories no elimina otras condiciones de licencia, SBOM o admisión.

V377: reutilizar resultados independientes por perfil sólo si fuentes, configuración y oráculos conservan identidad. Tras un fallo intermitente, conservar el run fallido y repetir la unidad completa afectada en un entorno nuevo; informar resultados compuestos y fallos abiertos. No llamar a un retry exitoso arreglo upstream ni repetir tres navegadores verdes sin un delta que lo justifique. En probes HTTP conservar HTTPMessage (headers case-insensitive) y leer el contrato de cache real del SDK antes de fijar el oracle.

V378: seleccionar el harness del recorrido afectado desde el índice observado. La consulta→turno reutiliza TestPublicAppointmentBrowserPostgres y sus oráculos SQL; no repetir92fases privadas sin un delta que las afecte. Probar antes/después del mismo defecto y tres configuraciones en un build. No exigir igualdad entre espacios NBSP y ordinarios de distintos ICU, pero conservar números, puntuación, zona y el instante UTC del receipt. Resolver paths largos de evidencia explícitamente y no adivinar nombres de scripts.

V379: para ayuda pública estática, reutilizar sesión real Next/JWE y probar4navegadores sobre el proxy TLS de fixture; no arrancar PostgreSQL ni proveedores sin efectos durables. Extraer los textos inline a una única fuente y fijar enlaces id/versión. Permisos por GET, errores y respuesta tardía se prueban separados de CMS, documentos confidenciales y capacitación; no promover estos últimos por un PASS de lectura.

V380: extraer todas las guías versionadas existentes en un mismo lote y probar una matriz explícita de permisos, no una guía por turno. Conservar hashes de textos/render previo y demostrar que el resto del handler queda idéntico. Comparar HTML renderizado con HTML, no className de JSX. Con exactOptionalPropertyTypes omitir una prop ausente; no pasar undefined ni relajar el compilador. Distinguir15/15guías existentes de CMS o capacitación integral.

V381: probar primero el pack existente fuera del perfil y sus propios tests antes de duplicarlo. Derivar navegación de los predicates reales del destino y probar15perfiles en lote. Para redirects, exigir status/Location con maxRedirects:0 y cookies compartidas; un route callback puede omitir requests de redirección. Conservar artefactos antes del siguiente Playwright run. Normalizar manifiestos a LF antes de fijar SHA, probando igualdad JSON. No contar etiquetas de rol como permisos ni navegación como capacitación.

V382: comprobar cada subconsulta de un portal contra el permiso real del endpoint, no sólo contra el permiso de entrada. Reproducir admin-only además de wildcard. Leer include de Vitest antes de elegir sufijo; un exit1 sin tests no demuestra regresión. Revisar imágenes con valores conocidos:100unidades menores no sonUSD100. Extraer el formatter existente y demostrar reversibilidad de la pantalla original; probar USD/JPY/KWD, límites enteros y rechazos sin inventar redondeo comercial. Reusar issuer/Go/PG/harness y comparar snapshot durable para lecturas, sin repetir92fases de comandos inalterados.

V383: para lecturas privadas, reutilizar el fixture Go/JWKS/PG y sus snapshots. Error boundary con texto fijo, ningún uso de error/reset/retry y anchor de ruta fija para GET explícito. Probar bearer insuficiente antes/después de reconsultar y luego sesión válida, sin escrituras. En Next, scope del alert al main: el anunciador del router también usa alert. No exigir500 de la pantalla si el streaming entrega el fallback: exigir rechazo backend y ausencia de datos. Conservar ambos rojos y no repetir gates de comandos/Go inalterados.

V383 native-notice shortcut: compare authenticated README component rows with versions.json through an explicit alias map. A versionless component can be vendored; resolve the official containing tree, record all blob IDs and verify retrieved documents against Git blob hashes before assigning a source identity. Do not invent a package version or infer DLL correspondence from matching a release number. Next MIT companion is already captured; do not repeat its acquisition.

V384: seguir el script de la plataforma real. Sharp Windows usa @img/sharp-libvips-win32-x641.3.3 y su receta descarga el ZIP MXE8.18.6; no usar las versiones globales de Linux como identidad Windows. libvips-42.dll/versions/avisos ya tienen paridad exacta con ese paquete firmado; no repetirla. Sigue ZIP/receta MXE y wrapper C++ separados. Derivar perfiles afectados por owner; adquisición no está en el parent runtime, por lo que sus dos archivos nuevos no incrementan68/795 ni69/803.

V385: no repetir ZIP→DLL/versions ni la licencia raíz:488files inspeccionados, DLL y versiones exactas al input firmado, LICENSE igual a libvips426af3f. component-recipe-map.json fija28versiones/28SHA256 con field_origins; reutilizar para adquirir fuentes completas bajo perfil, expandir GLib gvdb/PCRE y librsvg Cargo. Base MXEd973945 observada no es provenance histórica; preservar FAIL761. Receta80blobs y base27blobs ya verificados. Wrapper C++/sharp.node separados. Runtime68/795 y69/803 no incorpora core de adquisición.

V385 seguridad posterior: source-locations.json tiene27HEAD200 y un404libxml2; usar xml2-location-correction.json para elURLoficial2.15conmismoSHA. La publicación2.15.4introduce correcciones: FAIL763 y xml2-official-fix-binding.json fijan c94eb021→96498992,39commits/47files y textos oficiales. GitLabAPIcodificada404; mirror oficial GNOME verificado funciona. No repetir discovery ni tomarlatestSharp0.35.4como prueba de seguridad. Runtime aún contiene2.15.3; siguientes acciones deben producir reparación/contención verificable, no volver a contar firmas/DLLyaiguales.

V387: existing keyset API owns after/next_cursor. Correct missing UI traversal using independent per-list GET cursors, fixed size/session scope and full27row boundary fixtures. UNIQUE NULLS NOT DISTINCT requires unique fixture VIN/battery values. Reconstruct canonical source and install offline instead of copying pnpm trees across Windows long paths. Reuse unchanged artifact identity; do not resume V386 deferred security research. Evidence reconstruction_evidence/PRIVATE_PORTAL_PAGINATION_V387.md.

V388: seis listas privadas ya tienen paginación real (cliente pedidos/casos, fábrica unidades, admin pedidos/casos/leads).163registros por navegador,4proyectos y276tests; preservar permiso independiente de leads. No reabrir ese defecto sin delta. Al añadir evidencia/fuentes actualizar FRANCHISE_PREFLIGHT_GAP junto con README/notices/planes; el verificador calcula los conteos y rechaza drift. Native researchV386 permanece diferida por usuario.

V389: fechas de turnos cliente ya usan locale/zona del negocio en cuenta/gestion/recarga; no reabrir ese defecto sin delta. El cliente recibe label SSR y conserva datetime de API. En fixture WebKit separar navegación ya asentada antes de observar errores; esperar red inactiva antes de recarga forzada. No filtrar errores ni atribuir una prueba estable a reparación upstream. ISO Z y offset pueden representar el mismo instante; assert epoch exacto y etiqueta de negocio independiente.

V390: cancelacion cliente ya conectada cuenta→gestion→POST real→receipt ligado→GETdurable; perdidas de respuesta son resultado incierto, bloqueo de comandos y GETexplicito.4navegadores/16cancelaciones unicas con16audits/outbox,28intentosPOST incluidos rechazos;8regresiones de lectura/fecha siguen verdes. No reabrir este recorrido ni repetir esos gates sin delta. Registrar espera del documento nuevo antes del click que provoca reload; idle del documento viejo no prueba navegación completada. Próximo owner debe ser una brecha material distinta;45/48 yV386diferido se mantienen.

V391-V392: D/E/F contestadas desde scope previo,39observaciones reales; no volver a pedir corpus/cuentas/branding de negocio para el CLI. Candidato247 tiene TEST06/08 reejecutados sobre su SHA exacto y dos ZIP byte-idénticos. Reusar receipts sólo para ese payload o demostrar delta; no presentar empaquetado como compilación/firma/SCA. TEST02 requiere resolver integración de FAIL385 por owner; las mejoras web V377-V390 y auditoría21V374 ya están probadas. TEST03/V386 diferido; TEST07 no aprobado.

V393: el perfil BFF actual excluye Sharp y desactiva el optimizador; no confundir su investigación histórica V386 con una DLL instalada en el nuevo perfil. Mantener separado el bloqueo de los demás componentes/tooling. Comparar paths largos con representación extendida en ambos lados y delimitar secciones del lock por sus claves de nivel superior. La paridad y los recorridos conectados deben pasar antes de integrar esta selección.

V393 FAIL781: abrir archivos con prefix extendido no alcanza si os.walk empieza en una raíz ordinaria: puede omitir directorios largos sin error. Enumerar desde la raíz extendida y comparar conjuntos completos además de hashes. Conservar el inventario parcial y su corrección; no inventar paquetes faltantes.
