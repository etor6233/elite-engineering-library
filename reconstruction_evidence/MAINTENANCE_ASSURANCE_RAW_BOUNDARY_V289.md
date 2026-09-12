# V289 — contrato de assurance y originales separados de derivados

2026-09-07. Mantenimiento de biblioteca, T2801 y LIB-R10/T2805/T2807.
No rediseño, importación privada, envío, gasto ni promoción de producto.

## Cambio real y evidencia

PROJECT_ENGINEERING_CONTRACT.json, antes ausente, ahora enlaza cinco requisitos,
nueve pruebas y un benchmark pendiente con las nueve dimensiones de assurance,
incluyendo usabilidad CLI. AUTHORED, schema 1.1.0; no código de Microsoft/Google.
SHA-256 inicial: `89e6f88febdc6b400936f507b140c976c4bfe3c7150bb23f1d896cb1676516ff`.
Es trazabilidad sobre el roadmap existente, no otro plan. Un test conserva passed
por evidencia histórica estrecha V287; el resto sigue planned/blocked. EVID-01 es
captured, no aprobación del release actual. Assurance sigue blocked.

- Contrato real: plan PASS; evidence rechazo esperado, exit 1.
- ENGINEERING_EXECUTION_VALIDATOR 1.3.0 reconstruido dos veces: 15/15 archivos
  byte-idénticos. Pack/versiones/licencias sin cambios.
- Suite del kit: 11 tests PASS con raíz de autoridades explícita y nuevamente
  11 PASS sobre reconstrucción independiente. No son 22 tests distintos.
- Ocho probes PASS: plan válido; evidence incompleto rechazado; dimensión ausente,
  referencia inexistente, riesgo rebajado y passed sin evidencia rechazados;
  proven ficticio sin artefacto/release rechazado; contrato original sin mutación.
- Ronda D generada/enlazada y explicada, PENDING por respuestas faltantes.
  Readiness regenerado: BLOCKED, 43 observaciones. No porcentaje de cierre.

Evidencia local en Temp/elite-v289-ad41218f69b848bba366dd0944fba52e:
contract-plan.log, contract-evidence-expected-block.log, contract-regression.log
(fallido), contract-regression-fixed.log, materialize.log,
materialize-independent.log, rebuilt-plan.log, rebuilt-regression.log,
check_contract.py, contract-probes.log, readiness-before.json y readiness.log.
Probe AUTHORED SHA `478811af8a0924e58ce765fcba55b8f0b9b0280ecb0b38d97f39f61a49d5040f`;
salida SHA `1c0b5a6bcad3241d5aa827d4c318e32d9a4896d4e789f337c1226cd3a4209fe1`.
Los expedientes locales no se distribuyen ni heredan como aprobaciones.

Reproducir: materializar el kit en destino nuevo; ejecutar validate_project.py
con PROJECT_ENGINEERING_CONTRACT.json y --level plan. Para unittest del kit
aislado, fijar ELITE_AUTHORITY_ROOT al workspace sólo durante la prueba y
restaurar el valor previo. Evidence debe seguir rechazando los pendientes.

## Requisito raw, sin inventar importador

LIB-R10/TEST-09: preservar originales autorizados antes de transformar;
normalización, deduplicación, redacción e interpretación son derivados trazables.
Un hash prueba bytes, no verdad semántica ni exhaustividad de la API. No completar
ambigüedades por intuición. Retención/acceso/borrado deben definirse: raw no significa
PII permanente ni copiarla íntegra a prompts. No reencolar historia en Handle/live.

[Microsoft, arquitectura medallion](https://learn.microsoft.com/en-us/fabric/onelake/onelake-medallion-lakehouse-architecture)
separa originales preservados de capas transformadas. Consultada 2026-09-07.
Su aplicación a chats es decisión local AUTHORED, no importador Microsoft ni
exigencia de Fabric, otra base o gasto. [NASA](https://www.nasa.gov/reference/system-engineering-handbook-appendix/)
y [Google SRE](https://sre.google/sre-book/release-engineering/), consultados en la
misma fecha, sustentan separar verificación/validación e integración y comprobar
el artefacto de release; no certifican este diseño. Se conservan las autoridades
del kit sin afirmar reauditoría de todos los upstreams.

## Fallos propios y pendiente

FAIL-404: salida extensa, recuperación por rangos sin afirmar lectura de partes
omitidas. FAIL-406: faltó raíz de autoridades; cinco fallos, luego 11 PASS sin
modificar tests. FAIL-407: caller comprobó LASTEXITCODE después de script PowerShell;
proceso hijo y comparación 15/15 resolvieron el diagnóstico. FAIL-408: patch
rechazado por ancla parcial; sin cambios parciales, corregido con líneas completas.

No cierra T2801: falta readiness, authority map y evidencia integral. Siguen los
21 cores, observabilidad/privacidad, adapters, documentos, IA/evals, UX/soporte,
seguridad, recovery y release/ZIP de T2802–T2810. Importador histórico pendiente
de canal/contrato oficial admisible y corpus autorizado; no se inventan respuestas.
Inventario de código sin cambio: 160 packs/1433 archivos; franquicia 67/742.
