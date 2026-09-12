# Go Human Approval Core — V182

## Resultado estrecho

V182 materializa `GO-HUMAN-APPROVAL-CORE 0.1.0`, la cola de aprobación humana anti-estafa: separación de deberes (revisor distinto del solicitante), doble control para montos altos, velocidad por sujeto, dinero (pago/devolución) nunca auto-aprobado y decisiones inmutables.

No ejecuta el efecto: decide si el efecto puede ejecutarse. La ejecución del pago/devolución sigue en los workers de dominio; esta cola es la compuerta anti-estafa.

## Autoridad y procedencia

- Go stdlib (BSD-3-Clause) — core `AUTHORED`.
- PostgreSQL 18.6 (PostgreSQL License) — migración `0042` `AUTHORED`.

## Archivos materializados (6)

| Archivo | SHA-256 |
|---|---|
| internal/approval/approval.go | 5a71a34c49fb26bc3e611a055750a1cd80423410bd33cc8c63f77c87cd50b71d |
| internal/approval/registry.go | 43afe89cf81bd4e97ee4829f0805499fc36df9606442ae2ccebd24507d0f274b |
| internal/approval/approval_test.go | a42ae28920611105f5c9b6aa13b248511d0112fcb0fa3ffdb3b756d715696eea |
| db/migrations/0042_human_approval.up.sql | 4b412bb3cbab53c54b3a23aa8f9f8e2d7b60424d293cd5f814c8e9bee21af31b |
| db/migrations/0042_human_approval.down.sql | d456ea67ca5311540603caf6def7b7f5b44be8847e7dfb0ea27445fde8e01355 |
| db/tests/0042_human_approval.test.sql | 24e3a6741e7d8241e1eff57839294cf771fa226e9326db55bcd0c0b3a04458ce |

SHA-256 del pack: `2b903e816de8e4c52ec2043ed373699de82ea4fee1da320fc3c54ecb8f4f47d7`.

## Toolchain fijado

- Go 1.26.7 windows/amd64 (runtime local verificado).
- PostgreSQL 18.6 x86-64 (cluster aislado, puerto 55443).

## Verificación ejecutada

- `go test ./internal/approval/ -count=1`: **7/7 PASS** (auto-aprobar bajo umbral no-dinero, dinero nunca auto-aprueba, velocidad, duplicado, separación de deberes, doble control con dos revisores, rechazo, no-pendiente).
- `go test ./... -count=1` (7 paquetes): PASS.
- `go vet ./...`: exit 0.
- Round-trip de hashes: 6/6 bloques reproducen byte a byte; `go test` sobre el árbol materializado PASS.
- PostgreSQL 18.6 real (initdb → up → test → down): `ON_ERROR_STOP=1` exit 0; aserciones DO: 2 pendientes con `decided_at` nulo, decisión inmutable ligada a request, dinero permanece pendiente; down deja `approval_ns=true`.

## Cobertura de invariantes probada

1. Dinero (`payment`/`refund`) nunca auto-aprueba.
2. Revisor distinto del solicitante (separación de deberes).
3. Montos >= umbral requieren dos revisores distintos (doble control).
4. Velocidad por sujeto acotada (`ErrSubjectFlagged`).
5. Duplicado rechazado; decisión inmutable y ligada a request.

## Condiciones residuales

- Ejecución del efecto aprobado (pago/devolución/reserva) en los workers de dominio: bindear al state machine.
- Identidad real de revisores (OIDC) y doble aprobación por roles: integración del proyecto.
- UI/cola de revisión humana: pendiente (frontend multi-rol).

V182 añade la compuerta anti-estafa; no declara la ejecución productiva del efecto.
