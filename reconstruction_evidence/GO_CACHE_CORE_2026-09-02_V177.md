# Go Cache Core — V177

## Resultado estrecho

V177 materializa `GO-CACHE-CORE 0.1.0`, un borde de cache gobernado: keys tenant-scoped, store en memoria con TTL y evicción LRU, y carga cache-aside single-flight que evita stampede y no envenena la cache con errores. Sin dependencias de terceros.

No afirma durabilidad ni HA. El adapter Redis (RESP, `redis.io`) es la autoridad del cache compartido/HA y permanece `CONDITIONED` (cliente no disponible sin red/cache en este host).

## Autoridad y procedencia

- Go stdlib (BSD-3-Clause) — core `AUTHORED` (LRU/TTL/single-flight son patrones estándar, no copia de upstream).
- Redis (RESP, `redis.io`) — autoridad del adapter compartido/HA, no materializado.

## Archivos materializados (6)

| Archivo | SHA-256 |
|---|---|
| internal/cache/cache.go | 14fffefe553e98acc6488694109106bf5180e9803cb4f09d720d3712d3b6ecca |
| internal/cache/memory.go | fca7d8abeb45a3ba2729241fdd26abc6f4327c6d36956df10dcf0055cfc1302d |
| internal/cache/load.go | c50ba4fdbfc7174328b66f16357241c530d7480752722002f0ddb19926347eb8 |
| internal/cache/cache_test.go | 181529d42b172df778f0ed112b98f09631f23a205d06ef9b79ff1a85e88e8952 |
| internal/cache/memory_test.go | 4b8fdb2be7de30193ee657324233bfcf8bbb06b2aecddaac0f7236ed9ae41927 |
| internal/cache/load_test.go | 8391defccf6d63e9c5d1b41479bcb74ed398f65ad83c349b407bbbb6741a01aa |

SHA-256 del pack: `f533013a21cf78e092c714aee71e13fa9f635d906939a67ac3e06fac7454fef0` (18.698 bytes).

## Toolchain fijado

- Go 1.26.7 windows/amd64 (runtime local verificado).
- Sin PostgreSQL para este pack (en memoria, sin migración).

## Verificación ejecutada

- `go test ./internal/cache/ -count=1`: **11/11 PASS** (keys tenant-scoped con negativos, validación de key, set/get con copia, expiración TTL determinista con reloj inyectado, evicción LRU con toque MRU, delete idempotente, TTL/oversize inválidos, stampede single-flight 20→1 llamada, populate en hit, error no cacheado).
- `go test ./... -count=1` (aifoundation + search + cache): PASS.
- `go vet ./...`: exit 0.
- Round-trip de hashes: 6/6 bloques reproducen byte a byte los archivos probados; `go test` sobre el árbol materializado PASS.

## Cobertura de invariantes probada

1. `tenant_id` obligatorio en toda key (`Scope.Key` niega vacío).
2. TTL no positivo rechazado (`ErrInvalidTTL`).
3. Valor sobredimensionado rechazado (`ErrTooLarge`), nunca descartado en silencio.
4. Error de carga no se cachea (sin envenenamiento).
5. Stampede: 20 lecturas concurrentes del mismo miss → exactamente 1 llamada a la fuente.
6. Valores copiados en lectura (el llamador no muta la cache).

## Condiciones residuales

- Adapter Redis (RESP) con cliente fijado y target real: pendiente (CONDITIONED).
- Evicción por prioridad/size con TTL por key en modo compartido: se resuelve con Redis.
- Instrumentación/observabilidad de hit-rate: paso de proyecto.
- RAG-AGENTS (consumidor principal de sesiones/contexto): aún no materializado.

V177 añade la superficie `CACHE` ejecutable; no declara durabilidad ni HA.
