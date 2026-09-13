# V402319 — NPS con fuente real y contrato conectado

PROVEN_LOCAL del delta NPS. T2801 completo sigue abierto. Biblioteca205packs/2308
bloques:1999AUTHORED161ADAPTED148VERBATIM; perfil116packs/1591archivos.

El cálculo seleccionado de PostHog MIT, commit6fafbb9081bd15e79448af5650e02a4f9ea435cc,
se traduce en un owner separado. El resumen PostgreSQL existente lo llama después
de validar conteos y umbral de divulgación. No se reclasifica la API previa como
código de PostHog. La salida numérica conserva su precisión; toFixed(1) de la UI
upstream no se introduce en el contrato actual. Se declaran bins agrupados por SQL,
int64 y rechazo de conteos inconsistentes como fronteras de la adaptación.

Nueve casos TS originales ejecutados sin instalar módulos; seis equivalentes
tipados adaptados, oráculo racional/MaxInt64 y416153inputs de fuzz3s PASS. G0–G8
califican sólo ese cálculo. Ocho archivos reconstruidos exactos antes de habilitar
USE_REUSABLE_PACK. Luego cuatro composiciones completas exactas:1591/615/936/1591.
La regresión sobre fuente reconstruida crea una base fixture nueva, aplica84
migraciones y prueba HTTP/OIDC/PostgreSQL: consentimiento, aislamiento,24peticiones
concurrentes, replay/respuesta perdida, umbral y agregación. Reinicio real de PG
conserva tres respuestas y el resultado. Todos los procesos propietarios cerrados.
Vet y compilación de los imports afectados en perfiles pequeños PASS. Cinco Go
modificados analizados con DevSkim:0hallazgos. Grafo de módulos sin cambios: se
conserva el SCA317 de78módulos, no se repite la misma consulta ni suites ajenas.

Dos nuevos bloques ADAPTED MIT, tres VERBATIM (fuente, tests, licencia completos)
y tres AUTHORED (oráculo, configuración/lock, guía). Tres archivos previos cambian
para delegar/documentar;1580previos exactos. No nueva dependencia de runtime.
Notices y fuente fijada viajan con el pack; no hay atribución empresarial del glue.

FAIL945 RESOLVED_LOCAL. FAIL946: el primer runner quedó esperando EOF heredado tras
arrancar PG; se canceló ese runner y se cerró su cluster antes de crear la base.
La corrección usa logs en archivo y cleanup por el data-dir real; el ensayo
pendiente completo pasó. Se preserva el intento anterior, sin PASS inventado.

Receipts/SHA, casos, reconstrucción y deltas: NPS_SOURCE_ADAPTATION_V402.json y
NPS_SOURCE_ADMISSION_V402.json. No certifica PostHog completo ni producción.
Continúan T2801 G/H, T2810 release/destino durable; ARCA penúltimo y nativo último.
