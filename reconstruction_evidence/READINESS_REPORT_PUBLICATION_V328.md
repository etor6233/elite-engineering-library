# Diagnóstico y publicación de readiness — V328

Mantenimiento de biblioteca T2801/T2808, 2026-09-08. Resultado focal PASS;
TEST-47, sin completar TEST-08 del release final ni ninguna ronda pendiente.

## Defectos preservados

FAIL538: readiness0.6.1 resolvía la raíz fuera de manejo controlado: root ausente
producía FileNotFoundError/traceback/exit1. CLI inputs/reportes rechazados tampoco
tenían diagnóstico uniforme. Tres regresiones dieron9fallos/2errores.
FAIL539: después del chequeo exists, atomic_write usaba replace y podía
sobrescribir un reporte aparecido antes de publicar. Probe sintético directo
confirmó original_preserved=false; dos tests adicionales fallaron, incluyendo
competencia de procesos. Ningún reporte real se perdió ni se restauró.

## Corrección canónica

PROJECT-START-READINESS-VALIDATOR0.6.2, AUTHORED, sólo2bloques cambiados:
validate_project_readiness.py y test_validate_project_readiness.py.
Inputs/OS/Unicode producen PROJECT_READINESS_BLOCKED en stderr/exit2, sin
fabricar JSON de aprobación. Argparse conserva exit2 y validate_record no cambia.
El reporte se escribe en un temporal, flush/fsync y se publica con os.link:
la creación exclusiva rechaza un destino ya presente; finalmente elimina sólo
su temporal. Nunca hay fallback a replace si el filesystem no soporta hard links.
Un error después de publicar exige inspeccionar el destino antes de repetir;
la salida incierta no equivale a ausencia del reporte.

Esto no prueba durabilidad de la entrada de directorio ante corte eléctrico,
filesystems no ensayados ni defensa contra sustitución adversarial concurrente
de directorios. El claim de publicación exclusiva se ensayó en este Windows
local. Preparación/bridge/otros writers y release distribuido tienen gates propios.

## Verificación observada

-75/75tests PASS:70previos y5nuevos. Casos nuevos:6entradas inválidas,
3errores de destino,2errores I/O inyectados en la frontera CLI, un destino tardío
y3carreras con2procesos. Un solo ganador por carrera, payload completo y cero
temporales restantes. No contar6workers como6pruebas independientes adicionales.
-8/8archivos canónicos byte-idénticos entre reconstrucción y composición del
plan local de23archivos/2packs.6archivos (schema/template, catálogo, renderers,
README y pruebas de advisory) y el cuerpo exacto de validate_record permanecen
idénticos al baseline. Los dos pyc creados por tests históricos se registran
separados; el primer inventario de todos los archivos fue rechazado, no ocultado.
-El expediente real ejecutado sin publicar reporte devuelve BLOCKED/42 y un
objeto JSON idéntico al reporte anterior. No se rellenan respuestas/evidencias.
-Guía0.1.1 ejecutada desde sus bloques literales:13comandos/8helps NEW/EXISTING
PASS; originales, BOM, instrucciones y cadena sintética previa preservados.
Reutiliza el harness portable V327 y un fixture sintético V326, jamás el estado
real como objeto de tamper. No prueba usabilidad humana ni agentes independientes.

PowerShell7.6.5/Python3.14.4 disponibles; no runtime, paquete o fuente upstream
nuevos. Frontend UX §2, COMPOSITION_PROTOCOL, protocolo de recuperación y owners
readiness/execution gobiernan el método local; no se atribuye este código a
Microsoft/Google/NASA. Pins/catálogo/runner sincronizados a0.6.2/75tests;
snapshots V2830.6.1/70tests permanecen como evidencia histórica.

## Evidencia

Stage elite-v328-49d9267da9af45a89072a17605c82260, localizado por el puntero
elite-v328-current.txt del directorio temporal de esta sesión. Los expedientes
locales y el staging no entran en la distribución portable.

| Artefacto relativo al stage | SHA-256 |
|---|---|
| missing-root-before.log | 3ac342618d8d3627fa15716d77e756687c562215f61d6b535765021949293aaa |
| red-tests.log | 6087f70f98ffba883c7cc6f7c07dcfc6a2f0fb37710a6509099e6bf0aed03633 |
| red-publication-tests.log | 6dbebaad337eaba6b24006681efac3549e6300a718b4e7f5abf22cb3737522b3 |
| late-report-before.json | e24bc7f1fd6d80764b14409f4c873156db2a0c5b9370be7e8874d87d528ac04e |
| green-tests.log | c5fdbc6b68aeeee2c516af8940126bf4264416213465a5b41f4a7d5eb2caab34 |
| block-delta.txt | facb8a482c6af476208b8c08420d17db5809387a3f5d5a24f2558eb9f3674e44 |
| semantic-parity.json | 01365fdf297c71d706e13c40dd1d5565876b73dd01143f29a45fef6cf47f2579 |
| parity.json | 6b197371cb81c3a9eb2b936f0ea0424cf6418b6059f51491b6514ae75ed2d470 |
| root-readiness.json | 3cc10ad433e1411e38374bec500c7b4bd7a760dbb7dd638e04a018d817be21eb |
| guide-final/results.json | 5e4caf3bd31b9ddd80a7adbf0bbeb300c02f843173d3dfe6b39acb73ba2a18cb |
| guide-final/environment.json | 2284c7c07ecb6539c0828cf203b2ddd49f3a95a40acb4a7301b9084427d2ab10 |

| Fuente de la guía final | SHA-256 |
|---|---|
| markdown_system/PROJECT_ENTRY_CLI_GUIDE.md | fe8d625335febdc530e465ff82f369816a8a7ab4826447d77956aa0e3d13c7d7 |
| INSTALL_AGENT_BRIDGE.ps1 | fe32bf216ea48e9f218dbcf3075161bef445a82bd498bea3f14429d172a12a4b |
| materialize_markdown_pack.ps1 | 7ae045ff9b03ea22242d376d16b2a114f00e2f095217789eb92ca1f19e126356 |
| implementation_packs/MARKDOWN_COMPOSITOR_CORE.md | a8505384bb3f9085d36f21f1b6471fe24cbd17b054fa87f82f1d24b7007ce9d5 |
| implementation_packs/PROJECT_START_READINESS_VALIDATOR.md | 22c31cbc5ef0369767faac230ad5d03fb158b283434aeeb771fad3e13c547b9f |
| implementation_packs/ENGINEERING_EXECUTION_VALIDATOR.md | 368c4f9d0c328353634e3dc0738e0d300f225013b9c7a769628b971d0c854808 |
| markdown_system/PROJECT_READINESS_GATE_PACK_PLAN.md | f83d7a5133be17aa3e8a77a78de789085f6e8913db46ae98e6bf573bc6448adc |

## Reproducción

Materializar implementation_packs/PROJECT_START_READINESS_VALIDATOR.md en destino
ausente y ejecutar Python -X utf8 -B -m unittest discover -v desde
project_readiness_gate. Las cinco regresiones y sus subprocesos están en el pack
canónico; límite15s por proceso y10s de espera de barrera, sin retries. El test
concurrente usa3pares de hijos propios y conserva exactamente el payload ganador.
Comparar los8archivos del manifest contra una composición independiente;
separar caches de ejecución. Para la guía, usar el programa portable conservado
en CLI_ENTRY_DIAGNOSTICS_V327.md con fuentes vigentes y staging nuevo.
Para el expediente real, ejecutar el CLI sin --report y capturar stdout fuera
del proyecto; comparar el objeto JSON con el reporte previo, esperando exit2.

## Integración pendiente

Después de checkpoint87 ejecutar Preflight con rutas aisladas ya verificadas y
sin AllowNetwork. Guardar preflight87.json/log. No heredar un PASS anterior sobre
los nuevos bytes; la disponibilidad Docker y los42bloqueos reales siguen aparte.

## Integración observada y cierre de mantenimiento

Preflight87 terminó exit0 con151pasos PASS, incluida la suite readiness75 y
VERIFY_LIBRARY161packs/1453archivos/764Markdown/52perfiles; franquicia67/746.
Disponibilidad7/8, Docker ausente: statusBLOCKED. AllowNetwork=false conserva
fuera de esta corrida los gates de runtimes/servicios externos que requieren
red o target; sus omisiones explícitas no se convierten en pruebas superadas.
Los42bloqueos reales y los10macrofrentes permanecen.

SHA-256 preflight87.json: `8f2292285f4af2e54ca7b26cde8aa875ef358897b94e5308513389f5876d8e8b`.

SHA-256 preflight87.log: `24cc22dd5397d4455bdeaec663202b2a6814749b69c040b4574bdd0157acdd66`.

Checkpoint88 añade únicamente evidencia posterior y continuidad, sin cambios
de código/pins/guía/tests después de este preflight. Los87eventos anteriores
permanecen byte a byte. Se revalidan contrato y cadena sin repetir suites estables.
