# HTTP host shutdown — V314, 2026-09-08

> **Errata de trazabilidad V317 (2026-09-08), prevalece sobre las menciones históricas de Go1.26.7 debajo.** Se retira la atribución histórica exacta Go1.26.7 no vinculada a un receipt de identidad por ejecución. No se sustituye por1.26.8 sin prueba del artefacto concreto. La ruta current-toolchain de V281 contiene1.26.8; official-toolchain contiene1.26.7. El cuerpo original queda conservado como historia. V317 revalida únicamente la composición actual y el scope allí enumerado con1.26.7 observado, no todos los snapshots previos. Ver [TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md](TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md).

Mantenimiento T2802/T2809, FAIL-20260908-502. Código AUTHORED, no atribuido a Go.
Core0.4.3 y aplicación1.9.3 conservan15s de drenaje y comparten el owner HTTP.
Cuatro archivos reconstruidos desde los dos Markdown canónicos, identidad4/4;
composición67/745. No dependencia, regla de negocio, dato real ni efecto externo.

## Fallo y corrección

Ambos main lanzaban Shutdown en goroutine y retornaban al cerrarse el listener.
La solicitud seguía activa: red TCP real falla «host returned before the active
request drained». La extracción mecánica conserva el baseline en staging/red.py
y baseline/. Go1.26.7 net/http/server.go documenta que Serve/ListenAndServe retorna
antes de finalizar Shutdown; SOFTWARE_BACKEND_API_ENGINEERING11.2 exige detener
admisión, drenar con plazo y cancelar al excederlo.

ServeUntilShutdown ejecuta Shutdown en la ruta que espera el host; normaliza
ErrServerClosed sólo cuando hay cancelación solicitada. Un error del listener
también drena antes de retornar y conserva su causa. Al vencer el plazo, Close
cierra conexiones y se devuelve DeadlineExceeded; no se informa éxito. El
context de solicitudes activas no se ata prematuramente al de la señal.
Los dos main llaman al mismo owner antes del retorno normal/pool.Close.
«starting» describe el log previo al bind sin afirmar que ya está escuchando.

## Verificación y oráculos

- Cuatro casos TCP loopback: solicitud completa tras cancelar; admisión cerrada;
  deadline cancela conexión y el cliente falla; puerto ocupado retorna causa de
  listen sin esperar señal; cierre inesperado de listener preserva respuesta y
  error. La primera prueba cubre respuesta y admisión juntas.
- Un test adicional rechaza seis ciclos inválidos (ctx/server/callback ausentes,
  timeout cero/negativo, ctx ya cancelado), más nil server en la API pública.
- Cinco tests raíz ×3 en candidato y ×3 en reconstrucción:15PASS en cada árbol;
  seis subcasos por ejecución. No SKIP en este scope. El caso rojo se preserva.
- Go1.26.7 test/vet/build -mod=readonly ./... en ambos árboles PASS; mod verify
  candidato PASS. Los tests generales con gate PostgreSQL opt-in no prueban DB
  en este corte. No se atribuyen los92 recorridos V313 como ejecutados en V314.
- Inspección de ambos hosts: una llamada al owner con15s, ningún ciclo anterior.
  x/mod0.40 security floor y locks de V313 intactos; no nuevo claim OSV.

## Límites

Pruebas controlan cancelación de context sobre net/http real, no señales OS al
proceso completo ni un supervisor/target. No race detector, HTTP/2, WebSockets,
hijacked connections o handlers que ignoren cancelación demostrado. Close no
espera handlers arbitrarios: el plazo acota drenaje cooperativo de conexiones,
no garantiza commit/rollback del negocio. Startup/shutdown failure sigue exit1
y política redacted existente. No despliegue, monitor persistente, alertas,
carga, SAST/DAST, backup/restore, aceptación ni cierre integral T2809.

## Receipts y reproducción

Staging: `%LOCALAPPDATA%/Temp/elite-v314-d16ca9c1519e483ba477d9ae8863e622`.
Toolchain1.26.7 admitido por V313, GOTOOLCHAIN=local y TEST_DATABASE_URL ausente.
Reproducir: componer FRANCHISE_COMPLETE_PACK_PLAN en destino ausente; ejecutar
`go test -mod=readonly ./internal/platform/httpapi -run TestHTTPShutdown -count=3 -v`,
`go test -mod=readonly ./...`, `go vet -mod=readonly ./...` y
`go build -mod=readonly ./...`. Canonical.ps1 conserva baseline-packs, sincroniza
cuatro bloques y compara SHA de candidato/reconstrucción. No repetir mutaciones
sobre el mismo staging. Logs y scripts originales se conservan localmente.

| Receipt | SHA-256 |
|---|---|
| red.py | 66038d9ae43bbd214ef1f6f182edbf5bd9f2f196307f7bbbf4415dbf266adda3 |
| shutdown-red.log | 72dd5d48db22fb274f5b1f2747ba9ee57876a0093b272f9af4fe3fb0e285026a |
| shutdown-green.log | e9be7d5b5c2b99aa6ff5bf954c8982c633cb99ca5e94b9499826f0f564a5efde |
| rebuilt-shutdown.log | 2779a212bb3164c8edb96cd00d2383efa2f5a42d9850e8207010e50032914b5d |
| go-test.log | d5004a990ed8d516a2e3caa85f2457c4cc4220d4ca44fa928070c09bee6dec3f |
| go-vet.log | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |
| go-build.log | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |
| go-verify.log | acdb6a6a98dc2c297af31bb3538319778be8c8ef263cd0a8c9799de9c4998533 |
| rebuilt-go-test.log | 00faa33274e9a87ec45b64168e296b39f1734cba284dc7d9006915cc374f84d6 |
| rebuilt-go-vet.log | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |
| rebuilt-go-build.log | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |
| canonical.ps1 | 7d07d7ac90bf169be3c3e4be9c57730175370dbbe7a5c931ba68a48b7f49bed7 |
| compose-rebuilt.log | a2ed51adfb2302f8fa5c3352692224af968c0c85c80c47eb2aa1e51d9b032527 |

| Archivo reconstruido | SHA-256 |
|---|---|
| cmd/api/main.go | 98f2277d4ddafb998c3500c663856f7af524d61d4c9bdd9b2ee60db55954fbae |
| internal/platform/httpapi/server.go | 66f182b2e02a6ef076f27e39a06d4cf540f22eb230920eafcc5a35b572fdb6e1 |
| internal/platform/httpapi/server_test.go | 27ffaf9b6fcb4366ae7152a3950888484a2c10e89d93db96674a0d175d3237ce |
| cmd/electromobility-api/main.go | 5e0d2a1043666d16de1fb7956734547fa2c813baf47f47e081c5a31009708b59 |

TEST33/EVID24 documentan este corte; assurance de plan PASS y evidence integral BLOCKED. Ronda D/canal histórico, decisión comercial, target/IdP/providers V295 y ARCA diferida intactos. Sin ZIP/promoción ni porcentaje inferido.

## Cierre general

VERIFY_LIBRARY_PASS160 packs/1436 archivos/747 Markdown y51 perfiles; franquicia67/745. Checkpoint55 validado antes del verificador; TEST33/EVID24 plan PASS y evidence integral BLOCKED. Log library-final.log SHA256: b48f14a7c161a5ca7336e314bf0c2e3ddd670ca13d164986b4c8e8473f5d28f8.
