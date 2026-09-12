# AWS GenAI IDP — benchmark reproducibility 2026-08-27 V1

## 1. Fuente oficial y estado

La consulta se hizo contra la API oficial GitHub del repositorio `aws-solutions-library-samples/accelerated-intelligent-document-processing-on-aws`, sin Git y sin modificar source.

| Propiedad | Resultado observado |
|---|---|
| latest release | `v0.6.5`, publicada 2026-08-24; commit fijado `1b5fd74454e593de233342a02ee911af8ee38359` |
| main | `720f3052a57c22ae0d1e50886965c065ee86d60b`; commit no firmado según API |
| repositorio | activo, no archivado; MIT-0; pushed 2026-08-26 |
| PR oficial | `#669`, `test(benchmarks): commit the advverify suite the published results were run with` |
| PR head | `53ed38c2d8fc203fedd4fb0cae243f1727e5467b`, branch oficial `fix/commit-advverify-suite` |
| PR base | `develop` en `87962176d2217e2c7a1cb0200d3d8f24cb71e6a9` |
| estado al corte | abierta, no draft, mergeable/clean; no pertenece a v0.6.5 |

## 2. Hallazgo reconocido upstream

El cuerpo de la PR declara que `benchmarks/results/v0.6.5-advverify/` llegó a `develop` citando `--suite advverify`, pero la suite sólo existía en el working tree del autor. Por eso la evidencia committed no era reproducible. La PR agrega ocho líneas a `benchmarks/matrices/config_matrix.yaml` y una línea a la documentación; no cambia runtime.

También explica que el fallo medido es silencioso y no determinista: schema válido, scalar accuracy 1.000 y estado `COMPLETED`. Propone cuatro repeticiones y un control sobre el mismo documento. Esta explicación es evidencia de una brecha de reproducción, no prueba de que el release haya sido corregido.

## 3. Decisión

La biblioteca no descarga ni promueve una PR abierta como reemplazo de un release. `v0.6.5` conserva su identidad y sus gates previamente ejecutados, pero queda más condicionado:

1. sus resultados `advverify` no pueden invocarse como prueba reproducible;
2. sus findings OSV/npm y dos errores basedpyright integrales siguen abiertos;
3. una futura revisión debe incluir la suite e inputs completos, publicarse oficialmente y repetir benchmarks, dependencias, suites, seguridad y corpus target desde cero.

La condición queda en `UP-FAIL-069`; `OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.22` la expone antes de adquisición en el lock y perfiles AWS/document-intelligence.

Fuentes oficiales: <https://github.com/aws-solutions-library-samples/accelerated-intelligent-document-processing-on-aws/releases/tag/v0.6.5> y <https://github.com/aws-solutions-library-samples/accelerated-intelligent-document-processing-on-aws/pull/669>.
