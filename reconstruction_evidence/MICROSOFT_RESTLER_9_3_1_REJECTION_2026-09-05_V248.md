# Microsoft RESTler 9.3.1 — reconstrucción y rechazo V248

## Decisión

`REJECTED_COMPONENT / UPSTREAM_OPEN`. Microsoft RESTler 9.3.1 es código oficial MIT, reconstruible y funcional, pero no es admisible como pack reusable porque su grafo exacto falla SCA. No se modificaron dependencias ni source para fabricar un PASS.

## Identidad y toolchain

- fuente oficial: `microsoft/restler-fuzzer`;
- objeto: tag `v9.3.1`, no GitHub Release;
- commit: `eaacd92733771d35e90acfb2040fe37d4e5d09b1`, firma GitHub verificada;
- archive: 3.464.211 bytes, SHA-256 `e9c11f54d7172edd2444c8efdb796c64155ca118fc6c8db1fb777f3a4c892a75`;
- licencia MIT SHA-256 `c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383`;
- .NET SDK oficial 8.0.424 Windows x64: 285.090.820 bytes, SHA-512 `1787ab90635c2950672ed7c6507b000e1b212ea7d9a22fcef37061344d37c64d4c4eda12b8742601eff5b45c8736485b31c55613892f240c300190e4e88a58b0`;
- CPython 3.12.14 x64 y dependencias Python aisladas; `pip check` PASS.

El tag no incluye `global.json`; el build real gobierna desde `src/Restler.sln`, cinco proyectos `net8.0` y cinco `packages.lock.json`.

## Evidencia positiva

```text
build-restler.py release PASS
drop: 185 archivos / 27.476.288 bytes / manifest SHA-256 80e39eea…
dotnet restore --locked-mode PASS
dotnet build: 0 warnings / 0 errors
Restler.Compiler.Test: 70/70 PASS
Python engine E2E: 27/27 PASS en 855,98 s
Python restantes: 59/59 PASS en 44,23 s
test_fuzz source-exact PASS
```

El docstring de `test_fuzz` dice tres minutos/100 secuencias, mientras el código fija 0,1 horas y 300 secuencias. La ejecución respetó el código oficial; no se parcheó la divergencia.

## Bloqueo de admisión

OSV Scanner oficial 2.5.1 recorrió 13 resultados y 558 ocurrencias de paquetes. Terminó 1 con 23 hallazgos y 19 IDs únicos. Los componentes afectados observados incluyen `System.Text.RegularExpressions 4.3.0` en los locks .NET y, en el demo server, `h11 0.9.0`, `idna 3.7.0` y `starlette 0.45.3`. El commit firmado actual de `main` (`6d984deedbc54aad957fa3da0c7e9e5df23a2aee`) repite 23/19 sobre los mismos cuatro componentes; fijar ese commit tampoco resuelve el gate.

Un motor AppSec no gobierna seguridad externa si su propio grafo falla la política de vulnerabilidades. Reabrir únicamente ante revisión oficial limpia. No presentar un backport local como código Microsoft.

## Próxima vía oficial

Para el runtime Go de la biblioteca, evaluar el fuzzing coverage-guided nativo del toolchain Go: targets `FuzzXxx`, corpus de regresión, `go test -fuzz` y presupuesto explícito. Es una capacidad oficial del proyecto Go y está soportada por OSS-Fuzz, pero aún requiere un pack runner verificable y fuzz targets específicos del proyecto; no queda promovida por este expediente.

Fuentes oficiales: <https://github.com/microsoft/restler-fuzzer>, <https://github.com/microsoft/restler-fuzzer/blob/main/README.md>, <https://github.com/microsoft/restler-fuzzer/blob/main/docs/user-guide/QuickStart.md>, <https://go.dev/doc/security/fuzz/>.
