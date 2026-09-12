# V301 — integridad del histograma candidato

> **Errata de trazabilidad V317 (2026-09-08), prevalece sobre las menciones históricas de Go1.26.7 debajo.** Se retira la atribución histórica exacta Go1.26.7 no vinculada a un receipt de identidad por ejecución. No se sustituye por1.26.8 sin prueba del artefacto concreto. La ruta current-toolchain de V281 contiene1.26.8; official-toolchain contiene1.26.7. El cuerpo original queda conservado como historia. V317 revalida únicamente la composición actual y el scope allí enumerado con1.26.7 observado, no todos los snapshots previos. Ver [TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md](TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md).

2026-09-07. Desde checkpoint32/V300, mantenimiento correctivo aislado T2809.
La creación ordinaria de entrega sigue pendiente de la condición comercial ya
consultada; este paso independiente no la decide ni abre un roadmap paralelo.
V295–V300 y ARCA diferida conservados. No cuentas, gasto ni efectos externos.

## Defecto y corrección real

GO-OBSERVABILITY-CORE0.2.0 mantenía el histograma sin validar desde V292.
Cuatro regresiones nuevas fallaron sobre sus bytes exactos:

- observación200 con último límite100 devolvía p100=100, ocultando el overflow;
- p0 elegía el primer bucket aunque estuviera vacío;
- NaN contaminaba suma y count;
- count MaxInt64 pasaba a negativo al observar otro valor.

El cambio0.3.0 conserva dos archivos AUTHORED / LicenseRef-Workspace-Owner y
admisión CANDIDATE. NewCheckedHistogram comprueba configuración y la copia;
TryObserve rechaza datos no finitos/negativos y overflow sin mutación parcial.
Existe un bucket adicional para valores superiores al último límite finito.
Snapshot retorna copias coherentes de límites/counts/suma/total bajo el mismo lock.
Counts son no acumulados y su suma coincide con Count. No es mensaje OTLP.

Percentile devuelve límite superior del bucket, no valor exacto ni interpolado;
overflow retorna +Inf, p0 el primer bucket ocupado, entrada/configuración inválida
NaN. Vacío válido conserva0: no interpretarlo sin comprobar Count. math/big del
toolchain Go fijado calcula el rango exacto del valor binario de p, evitando
pérdida de precisión de count/int64; se usa sólo al consultar. La suma float64
sigue teniendo redondeo normal: no usarla como ledger financiero o inventario.

Compatibilidad explícita: NewHistogram mantiene firma pero una configuración
inválida queda inerte; preferir NewCheckedHistogram y TryObserve para observar
rechazos. El antiguo test que esperaba100 para overflow ahora exige +Inf: se
corrige un oráculo equivocado, no se retira la aserción. Un futuro exporter debe
representar ausencia/NaN/+Inf según su contrato, no enviarlos como JSON numérico.
El límite1024 es decisión local AUTHORED, no requisito de Google u OTel.

## Fuentes y límites de autoridad

Consultadas el2026-09-07:

- [Google Cloud Monitoring Distribution](https://docs.cloud.google.com/monitoring/api/ref_v3/rest/v3/TypedValue#Distribution):
  población, count igual a suma de buckets, advertencia sobre valores no finitos
  y overflow separado. Google usa límite inferior inclusivo/superior exclusivo.
- [OpenTelemetry data-model](https://opentelemetry.io/docs/specs/otel/metrics/data-model/):
  aquí se conserva el límite superior inclusivo. No afirmar compatibilidad
  directa con el formato Google ni convertir este componente en un SDK OTel.
- [OpenTelemetry metrics/api](https://opentelemetry.io/docs/specs/otel/metrics/api/):
  las observaciones del histograma se esperan no negativas; nuestra validación
  y API son locales, no implementación atribuida a la especificación.

Los manuales SECURITY_SRE_CLOUD_INFRASTRUCTURE sección15 y
COMPUTER_SYSTEMS_PERFORMANCE_LOW_LATENCY sección12 gobiernan distribuciones,
límites y costo de instrumentación. No se adquirió ni actualizó ningún upstream.
Se reutiliza Go1.26.7 verificado V281; math/big es biblioteca estándar, no nueva
dependencia. Citar las fuentes no implica aprobación del código AUTHORED.

## Evidencia

Staging: `%LOCALAPPDATA%/Temp/elite-v301-33aaeda4aa534c09bb127dd854f5c43c`.
Baseline0.2.0 reutilizado desde reconstrucción V292 con ambos SHA coincidentes.
El candidato se devolvió al pack canónico; rebuild de investigación a destino
ausente extrajo sólo sus dos bloques y validó sus SHA. Esto NO habilita composición
de producto: el compositor real rechazó0.3.0 aun con acknowledgeConditions=true
y dejó el destino ausente. No se modificó ese guard ni el perfil integral67/745.

- Candidato:24 tests y17 semillas explícitas PASS.
- Reconstrucción: los mismos24 tests y17 semillas, repetidos tres veces PASS
  (72 ejecuciones de tests y51 de semillas; no72 casos diferentes).
- Configuración: vacía, desordenada, duplicada, negativa, NaN, infinitos y límite
  excesivo rechazados; límite1024 válido, copia de configuración y snapshot
  sin alias, instrumento zero-value cerrado.
- Recuperación: rechazo atómico de datos y suma/count overflow, retorno a
  observación válida; consultas inválidas no se presentan como medida cero.
- Concurrencia:16 goroutines,2048 observaciones válidas y2048 rechazadas;
  snapshots coherentes y agregado final2048/4096. No equivale a race detector.
- Percentiles: límites exactos, overflow, vacío, p0/25/50/100 y MaxInt64.
- Go test ./... -count=3, vet ./... y build ./... PASS, GOPROXY=off,
  GOTOOLCHAIN=local. Vet/build exit0 con logs vacíos.
- GO_NATIVE_FUZZ_GATE0.1.0 materializado4/4, ambos scripts parseados y
  verify_pack.ps1 PASS: positivo real y negativos de perfil/traversal.
- Gate sobre reconstrucción: tres targets,10s cada uno, GOMAXPROCS=2, PASS.
  FuzzHistogramAggregate ejecutó586292 entradas; Counter1554913. Son ejecuciones,
  no cobertura exhaustiva ni casos únicos. El modelo de histograma compara
  count/suma/buckets y rechazos con las entradas sintéticas conocidas.
- Logs y receipt hash-linked conservados; ningún corpus fallido fue descartado.

No se ejecutó -race: CGO_ENABLED=0 y CC=gcc sin gcc/clang resuelto en PATH.
No se instaló un compilador global ni se presentó concurrencia como sustituto.
No se ejecutaron carga, SCA nuevo, exporter, alertas, UI o PostgreSQL: este
candidato no está en la composición; la evidencia V300 no cambia ni se recuenta.

| Archivo/receipt | SHA256 |
|---|---|
| observability.go | 4ac21741bab93711a40bf5287772a266e40218592ac203af655c9af600f8ed5e |
| observability_test.go | 65fc115d5bb54c8e1ec7c67b6855506c01dd1ef4202d8c6878ae75a421334a22 |
| histogram-red.log | 24e39c6c76b7a16f446c1f94044c2e592313f4361d8a6415437e7c7a39c15228 |
| candidate-test.log | 7ae61d12bab3488700969174912861147494db5d6d9aca879ecb58b3fa92d2eb |
| rebuilt-test.log | 1fc77199e907d3d615315ff546ed39f70932a83a5414c7e323bd1e58e661fed7 |
| rebuilt-vet.log y rebuilt-build.log | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |
| fuzz-runner.log | 1c7ebb147b6da6e0707ba1791d419d7937ecb3a0fd2f9fba26c2f4ba84ec4f27 |
| rebuilt-fuzz.log | b4baaf2eabfd0703e37868c483eef286189beeb3daf33b88dd8c7a49e80e5228 |
| rebuilt-fuzz-receipt.json | 9911a31284b483d768cce66354386b54f6526bacce0842053a715f35696bbe6d |

Reproducir sólo en módulo de investigación Go1.26.7: extraer dos bloques
exactos verificando SHA, go test ./... -count=3, go vet ./..., go build ./...;
materializar GO_NATIVE_FUZZ_GATE y ejecutar su runner con FuzzPolicyLoggerClosedValues,
FuzzCounterMonotonic y FuzzHistogramAggregate, presupuesto10s por target.
El perfil habilitado y el receipt permanecen en staging; los targets/semillas
reutilizables sí están en el Markdown canónico.

## Continuidad sin promoción falsa

FAIL453 queda corregido en este scope; FAIL432 registra errores de diagnóstico.
No quedan abiertos esos cuatro defectos del histograma en0.3.0. Aún faltan
política target, race, montaje con host/exportador, alertas, retención, overhead,
carga, recuperación y gates integrados para readmitir observabilidad/T2809.
No introducir este candidato como sustituto de instrumentos oficiales admitidos.

La próxima dependencia operativa sigue siendo entrega ordinaria: condición de
liberación pendiente, comando/frontend iniciales y enlace con cobro/stock/pedido/
transporte. La pregunta se volvió a presentar sin asumir respuesta. Trabajo
independiente del roadmap puede continuar, pero no hay READY_TO_BUILD integral,
producción, ZIP final ni porcentaje/plazo demostrado.
