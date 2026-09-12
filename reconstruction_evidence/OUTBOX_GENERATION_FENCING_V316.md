# Outbox generation and lease fencing — V316, 2026-09-08

> **Errata de trazabilidad V317 (2026-09-08), prevalece sobre las menciones históricas de Go1.26.7 debajo.** Se retira la atribución histórica exacta Go1.26.7 no vinculada a un receipt de identidad por ejecución. No se sustituye por1.26.8 sin prueba del artefacto concreto. La ruta current-toolchain de V281 contiene1.26.8; official-toolchain contiene1.26.7. El cuerpo original queda conservado como historia. V317 revalida únicamente la composición actual y el scope allí enumerado con1.26.7 observado, no todos los snapshots previos. Ver [TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md](TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md).

Mantenimiento T2802/T2809. Core0.4.4 + GO-RELIABLE-ASYNC-WORKERS0.3.1,
cuatro archivos AUTHORED reconstruidos idénticos. Composición67/745, sin nueva
migración/dependencia. Reutiliza el método de generación/plazo ya demostrado en
Jobs del mismo owner; esa evidencia previa no sustituye los rojos nuevos.

## Defecto y oráculos

MarkPublished/Release cotejaban sólo worker ID y published_at. Seis rojos PostgreSQL
demuestran confirmación/liberación con lease vencido, intento viejo tras reclaim
del mismo worker y espera de row lock que supera el lease. Tres rojos del processor
demuestran publisher llamado sin presupuesto, context superior al lease y segunda
entrada tardía publicada. FAIL504; logs before preservados.

El primer observador de row lock consultaba pg_stat_activity en el tx que bloquea;
su snapshot estadístico no se refrescaba y no observó al waiter. Se corrigió a
consultas autocommit del pool, sin cambiar el lock ni el oráculo del reloj DB.
El rojo refinado confirmó las seis fallas. FAIL505 conserva before/after; no se
confunden los dos timeouts iniciales de diagnóstico con prueba del defecto SQL.

## Corrección canónica

Claim devuelve attempts como generación monótona; nunca debe reiniciarse. Las
firmas internas de MarkPublished/Release requieren esa generación. Un CTE
materializado bloquea la fila antes de comparar worker, generación y reloj DB;
un lease perdido devuelve ErrOutboxClaimLost sin modificar estado. No basta con
agregar clock_timestamp a un WHERE evaluado antes de esperar el row lock.

RemainingLease requiere dueño, generación, estado pendiente y plazo vigente.
El processor resta conservadoramente todo el roundtrip de esa consulta, limita
el context del publisher y vuelve a consultar cada entrada del lote. Confirmación
y reprogramación vuelven a comprobar el fence después del callback. La clave
de idempotencia sigue siendo EventID, nunca la generación variable.

Core0.4.4 y workers0.3.1 deben componerse juntos; el worker declara compatibilidad
exacta demostrada. Los consumers con firmas anteriores requieren actualización.
No se publica una API alternativa que permita confirmar sin generación.

## Pruebas ejecutadas

- Dos DB locales nuevas separan packages que Go ejecuta concurrentemente. Cada
  una usa sólo0001_platform_foundation; no hay bypass de triggers ni acceso al
  audit compartido V296. Tests anteriores pasan a OUTBOX_TEST_DATABASE_URL y el
  processor a OUTBOX_PROCESSOR_TEST_DATABASE_URL; exigen host127.0.0.1, prefijo
  elite_outbox_test_ y sufijo>=16, rechazan fallback remoto. Processor rechaza
  mismo DB que store tests. TEST_DATABASE_URL genérica no habilita este scope.
- Seis tests raíz PostgreSQL (incluido guard con seis negativos) y cinco de
  processor ×3:33PASS por árbol, candidato y reconstrucción, sin SKIP focal.
  Casos SQL incluyen dos operaciones ×tres fallos, RemainingLease vigente/
  expirado/reasignado/publicado/dueño incorrecto y claim disjunto/control válido.
- Processor real + store PostgreSQL + transporte sintético: primera tanda
  confirma un evento y conserva retry con PUBLISH_FAILED; siguiente tanda usa
  mismo EventID, generación1→2 y confirma. Dos filas publicadas, claims/error
  liberados y tres invocaciones observadas. Tres receipts por árbol.
- Go1.26.7 test/vet/build -mod=readonly ./... en ambos árboles PASS y mod verify
  candidato PASS. Otros tests opt-in con TEST_DATABASE_URL ausente no prueban
  sus DB aquí. No se transfieren92 browser phases como ejecuciones V316.
- Reconstrucción4/4 desde dos packs; locks y fundación idénticos a V315. No scan
  SCA nuevo; piso x/mod0.40 y evidencia V313 conservados.

## Límites y continuidad

Fencing de estado local y context cooperativo no deshacen un efecto remoto ya
emitido. Entrega sigue al menos una vez: provider debe honrar idempotencia y
requiere reconciliación tras ambigüedad. No publisher real, exactly-once externo,
cancelación de adapter no cooperativo, host/supervisor/alerta, dead-letter policy,
retención/carga/race/reinicio o operación productiva demostrado. No se inventan
máximos de retry ni decisiones de negocio para cerrar T2809.

TEST35/EVID26: plan PASS, evidence integral BLOCKED. D/canal histórico, decisión
comercial, target/IdP/providers V295 y ARCA diferida intactos. Sin ZIP/promoción
ni porcentaje inferido. Continuar por owners existentes y sus gates.

## Reproducción y receipts

Staging `%LOCALAPPDATA%/Temp/elite-v316-2b9a0e49cb264ee3b0272801e0b1d62e`.
PostgreSQL18.6 loopback, Go1.26.7, GOTOOLCHAIN=local. Dos bases retenidas:
`elite_outbox_test_2b9a0e49cb264ee3b0272801e0b1d62e` y
`elite_outbox_test_processor_2b9a0e49cb264ee3b0272801e0b1d62e`.
No DROP. Para reproducir, crear otras dos DB exclusivas con fundación0001,
configurar los dos env vars y ejecutar en composición ausente/reconstruida:
`go test -mod=readonly ./internal/platform/postgres ./internal/platform/workers -run 'TestOutbox|TestProcessorsSeparateSuccessAndRetry' -count=3 -v`.
Luego `go test`, `go vet`, `go build -mod=readonly ./...` y `go mod verify`.
Red.py/red-refined.py registran los baselines; no repetir sus mutaciones. Los
scripts de baseline, sync, rebuilt-gates y receipts se conservan en staging.

| Receipt | SHA-256 |
|---|---|
| red.py | 1a2daf9fd284b3f5ba407d5dceac039c9e6ae31e416244cd2e23974efe3641e4 |
| red-refined.py | 6245d3035981265de674ec42e63c5ac0d5d099a14337e621d943c29f9f6ac44d |
| database.log | 39f0fc2c99aedbbdb53157c993739da64f0885fb6174c52321dab1d89953e692 |
| processor-database.log | 39f0fc2c99aedbbdb53157c993739da64f0885fb6174c52321dab1d89953e692 |
| outbox-red.log | 2c440a9e68b228a5698699ab1b0f1c66ffc5a6e4518b1203c1f48bb304e2e07e |
| outbox-red-refined.log | 8e4aab2eaf9a58447a9e263a506260b60cb20512fae05451d5cc7adf177a7334 |
| publisher-red.log | 7a8ee236282bedde6e5c4a839e3853f89ef0625d3ea5289d86143baa8933f3d4 |
| outbox-connected.log | 8faf781006daed7e4350eb177df2c6faabbe1b617aadfbc6a8c335378afa3916 |
| rebuilt-connected.log | 2e6897441afdedc8e98deb1696cad6fc46008e95d0338c135bba586cdb35ad6c |
| go-test.log | 34ae2d41551a663bf393b4c9cc4f413cb3dc0a46cadc977610d9fe32bfed0d97 |
| go-vet.log | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |
| go-build.log | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |
| go-verify.log | acdb6a6a98dc2c297af31bb3538319778be8c8ef263cd0a8c9799de9c4998533 |
| rebuilt-go-test.log | f4f7f498cc5145542f811e289e2954e8a2cfcad84130f1989ac05fcf3f3c1e0c |
| rebuilt-go-vet.log | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |
| rebuilt-go-build.log | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |
| canonical.ps1 | 89991002928860d482a24650281fee69314c85159d7942f2dbd8a68008cbe023 |
| rebuilt-gates.ps1 | de1e88ef9fb21b6a33d11946ba97f608d82cd698c0fef2f64b7c8e85ad76637d |
| compose-rebuilt.log | 6b03edbe93bc391490d0f98c0422fc2195852c1858632f5ac39bbc055fb41ed0 |
| rebuilt/internal/platform/postgres/outbox.go | ac68c30fdc76b1b9d9fb2a3c6014a9d805604ffa3bf6625ee8b939b399d72328 |
| rebuilt/internal/platform/postgres/outbox_integration_test.go | f4f8ea1921f2981e5ea93a9e0b22b6dc504978e7abd996808a37647c97e129aa |
| rebuilt/internal/platform/workers/processors.go | 7bff6315aae3f6c0a2b04253860c9056649ee6527208c5191b9f84a560c62314 |
| rebuilt/internal/platform/workers/processors_test.go | c07d2a932b516fda912226bc44ff3dd55604e6dd6f641ce5ba0ebb34ed10e00e |

## Cierre general

VERIFY_LIBRARY_PASS160 packs/1436 archivos/749 Markdown y51 perfiles; franquicia67/745. Checkpoint59 validado antes del verificador; TEST35/EVID26 plan PASS y evidence integral BLOCKED. Log library-final.log SHA256: 634b5d0b47b67cdddb65d2096fd528d0f3c791cedbdd21b485405cdf8a5ba39b. El cursor siguiente conserva los87 artefactos y la cadena histórica; reduce must_read a los12 refs materiales del siguiente paso y mueve historia estable a reuse_without_reload.
