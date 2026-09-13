# Indicadores conectados por rol — V402

PROVEN_LOCAL para este delta T2804. **No cierra T2804 ni declara READY global**:
resta idioma privado y los siguientes controles en el orden del usuario.
Canonical198packs/2107blocks; franquicia103packs/1343archivos exactos.

Doce familias leen los owners actuales: pedidos por estado/moneda, leads,
stock, casos, envíos, turnos, vistas propias del cliente, producción por fábrica
o destino y compras con plan seriado. No contabilidad de ingresos: registros
operativos actuales con fecha por consulta. Puntos y NPS usan los owners
StoredValue/CustomerFeedback y sus perfiles/retención/mínimos existentes.
GO-DASHBOARDS-CORE conserva su estado aislado; reviews/POS NONE_WITH_REASON.

AUTHORED glue de consulta, transporte, panel y pruebas, no código corporativo.
G0–G8 y SHA de fuentes/packs/salidas/receipts en JSON.19archivos nuevos,
7revisados;159ADAPTED/141VERBATIM y todas las dependencias originales intactos.
El nuevo pack contiene sólo el ensayo compuesto, para perfiles con sus owners;
las interfaces nuevas se agregan al mismo pack que su implementación.

Pruebas nuevas:
- Importes exactos superiores a JS safe integer:18014398509481986 totalARS,
 9007199254740993 propio;USD separado. Tenant/org/subject/permisos;100grupos
 como máximo, grupo101409 y error de servicio503, nunca un cero inventado.
- Emisión aprobada, pendiente excluida, replay sin incremento, reverso separado
 y perfiles/programas aislados. NPS con0/1/2respuestas guardadas por Submit,
 mínimo2 y definición vencida no disponible.
- Plan de supply creado por sus writers originales, submit/confirm/start y
 registro de1unidad. Vista fábrica/destino separada y compra sin plan excluida.
- Navegador realNext/BFF/Go/PG/Chromium conJWE/RS256/JWKS;4.21s. Selecciones
 por permiso, cambio de organización borra resultado, vacío y unavailable
 diferenciados;2respuestas reales actualizan NPS.0POST de métricas.
 Capturas390px inspeccionadas;no overflow del documento.
-9testsweb ybuild/tiposPASS. Fuzz2s/6semillas/227971ejecuciones.
-5rebuildsexactos:103/1343full,15/232web,40/576backend,64/871serverless,
 104/1351HTTP. Tipos web y compilación backend/serverlessPASS.

FAIL868: Overview excluía won mientras el estado terminal real es converted.
Regresión RED→PASS y corrección mínima preservadas. FAIL869–873 fueron
precondiciones de mutador/fixtures: código de tenant único, cancelación previa
al reverso, encuesta inmutable y plan de supply con evidencia transaccional.
Los gates rechazaron correctamente los intentos; no se deshabilitó ningún guard.
Algunos logs de lotes conservan FAIL de una fixture posterior: el JSON cita
el PASS nominal de cada prueba terminada y su corrección restante independiente.

El body del BFF genérico se autentica antes de consumirlo;65536bytes,UTF8fatal,
deadline2.5s,cancelación y autorización original poracción. Baseline antiguo
devolvió400de parseo antes de401/403; la prueba corregidaPASS. Esta es FAIL874.
El cursor había reutilizado FAIL457 como abreviatura: su cierre histórico de
sender/recovery V312/V325 permanece cerrado; no se reescribe la historia.
En notas de preparación del delta, «FAIL457 body» se refiere a esta FAIL874.

Sin cuentaslive, nuevas dependencias, investigación Daybreak ni producción.
Continuar T2804 i18n privado desde la composición exacta103/1343.
