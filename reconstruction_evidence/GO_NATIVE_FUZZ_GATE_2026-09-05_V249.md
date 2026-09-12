# Go Native Fuzz Gate — V249

## Decisión

`ELITE_REFERENCE / REBUILD_VERIFIED / CONDITIONED` para ejecutar fuzzing nativo del toolchain Go sobre targets reales aportados por el proyecto. Los cuatro archivos son `AUTHORED`; no se presentan como source copiado de Google o del proyecto Go.

## Autoridad y prueba

La documentación oficial de Go define targets `FuzzXxx`, corpus semilla, ejecución como test normal, fuzzing coverage-guided con `go test -fuzz`, presupuesto `-fuzztime`, persistencia de inputs que fallan y necesidad de targets rápidos, deterministas y sin estado global persistente.

Windows x64, Go 1.26.7 y PowerShell 7:

```text
PowerShell parser PASS: 2 scripts
materialization roundtrip: 4/4 archivos
baseline go test ./... PASS
seed corpus FuzzRoundTrip PASS
coverage-guided fuzz 1s PASS: 1.526.076 execs observados
profile disabled REJECTED
package traversal ../x REJECTED
GO_NATIVE_FUZZ_PACK_PASS positives=1 negatives=2
```

El número de ejecuciones no es un umbral portable. Cada proyecto define targets, semillas, presupuesto y evidencia propios; un PASS no demuestra exhaustividad, DAST, seguridad ofensiva del deployment ni ausencia universal de vulnerabilidades.

Fuentes oficiales: <https://go.dev/doc/security/fuzz/>, <https://go.dev/doc/tutorial/fuzz> y <https://go.dev/doc/security/best-practices>.
