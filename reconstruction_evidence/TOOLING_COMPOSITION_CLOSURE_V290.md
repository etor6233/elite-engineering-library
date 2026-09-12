# V290 — composición de mantenimiento comprobada y cierre pendiente

2026-09-07. Scope: tooling local T2801/T2808, no implementación comercial.

## Avance ejecutado

PROJECT_PACK_PLAN.md ya contiene selección ejecutable, no sólo una tabla.
Bootstrap canónico del compositor 0.2.2 (3 archivos), seguido por composición de
readiness 0.6.1 (8) y execution validator 1.3.0 (15). Dos destinos independientes
producen 23/23 archivos de validators byte-idénticos. Cada destino agrega un
MATERIALIZATION_RECORD.md; sus timestamps difieren legítimamente y no se afirma
identidad del receipt ni reproducibilidad de un ZIP final.

Pruebas ejecutadas sobre el resultado:

- execution validator: 24 tests PASS;
- readiness validator: 70 tests PASS;
- suite del compositor: PASS, incluyendo selección/condiciones/integridad/paths;
- segunda composición sobre destino ocupado: exit 1 esperado, 24/24 archivos
  conservan bytes, incluyendo receipt;
- preflight de selección inicial rechazada: cero archivos generados.

No se modificaron validadores ni fixtures para lograr PASS. La condición de
colisiones y rollback del plan se apoya además en V284/V287, que verifican
recuperación y conservación del proyecto EXISTING. Esto cierra la evidencia del
plan local seleccionado, no los gates release-bound de todo el perfil.

Staging preservado: Temp/elite-v290-5fcdc9bf739d417babe3691898563e2b:
bootstrap/, composition-b/, composition-c/, compose-b.log, compose-c.log,
collision-expected.log, execution-tests.log, readiness-tests.log, composer-tests.log.
SHA del plan Markdown: `5cfb8eb2e653f00ac1e1536fe89838b9e1fc209e522bfa149ffbb2cd0bfe5fec`.

El validador real reduce observaciones de 43 a 41 al reconocer pack_plan PROVEN
con colisión comprobada, rollback definido y evidencia. Readiness global sigue
BLOCKED; no se cerró otro gate ni se alteraron sus reglas para lograr la reducción.

## Fallo encontrado y corregido

FAIL-410: incluir el compositor en su propio plan hacía rechazar placeholders
literales de sus fixtures. Se corrigió la selección: bootstrap por extractor y
composición de los dos validators. No se debilita el control de variables ni se
borra ningún test. No se afirma que el compositor pueda autoensamblarse con su
propio gate actual. El flujo de dos etapas usa el protocolo canónico existente.
FAIL-404: lectura inicial demasiado extensa; recuperados rangos materiales sin
interpretar como leídos los segmentos omitidos.

## Cuánto falta: denominador honesto

El objetivo es una biblioteca portable, trazable y ejecutable para construir o
ampliar sistemas completos; no una colección de manuales ni una franquicia
productiva preautorizada. Los diez bloques T2801–T2810 siguen abiertos con avances
parciales. T2800 es auditoría terminada, no un undécimo componente de negocio.

| Bloques del roadmap existente | Cierre que sigue pendiente |
|---|---|
| T2801 | readiness/assurance y journeys de mantenimiento; no repetir entrada V287 ni composición V290 como si faltaran |
| T2802 | transacciones conectadas y equivalencia/admisión de 21 cores |
| T2803 | identidad, autorización, privacidad y seguridad de composición |
| T2804 | UX por rol, accesibilidad, ayuda/capacitación/soporte conectados |
| T2805 | adapters y reconciliación; importación histórica aún no implementada |
| T2806 | ingesta documental → evidencia → evaluación → mapping → commit/recovery |
| T2807 | IA/tools/evals representativos, handoff y coste demostrado |
| T2808 | entorno/build/deploy/rollback de referencia; tooling local no lo sustituye |
| T2809 | privacidad/métricas abiertas, operación, carga y recuperación |
| T2810 | integración final, promoción por claim y ZIP exacto independiente |

La unidad de cierre será un journey con sus gates satisfechos, no un Markdown,
un test o un porcentaje por cantidad de archivos. No hay base para prometer
15% restante, cierre hoy/mañana o una franquicia completa en una semana.
No significa empezar desde cero: 160 packs/1433 archivos y perfil 67/742 existen,
pero su inventario no mide trabajo restante ni tiempo de integración.

Lo local independiente puede continuar. Para historia real siguen pendientes
canal/export/API permitido, acceso, corpus autorizado y políticas del proyecto;
para una franquicia productiva, cuentas/configuración empresarial, infraestructura,
regulación y aceptación reales. No exigir esas cuentas para pruebas locales ni
usarlas para justificar dejar inconcluso código que sí pueda verificarse aquí.

[Google SRE Release Engineering](https://sre.google/sre-book/release-engineering/)
y [Reliable Product Launches](https://sre.google/sre-book/reliable-product-launches/),
consultados 2026-09-07: pruebas reproducibles, artefacto de release y comprobaciones
proporcionales al lanzamiento. Aplicación local AUTHORED; no nuevo código Google,
certificación externa ni afirmación de actualidad de todos los SDKs.
