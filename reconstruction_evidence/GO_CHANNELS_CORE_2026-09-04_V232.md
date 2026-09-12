# Go Channels Core — evidencia V232

## Resultado

```yaml
pack: GO-CHANNELS-CORE
version: 0.2.0
implementation: REBUILD_VERIFIED
admission: CONDITIONED
verified_at: 2026-09-04
materialization: 4/4 PASS
gofmt_idempotence: 4/4 PASS
full_go_packages: 47 PASS
pack_sha256: d5bc0f13fe3b371cb3852c9c03a5e1f501e90247f547d9540b725af1a9354f9c
```

## Contrato demostrado

Un inbound ahora exige `ProviderMessageID` y `OccurredAt`. El dispatcher entrega el `channels.Message` completo al responder; tenant, canal, contacto y thread se conservan al responder. La regresión rechaza identidad faltante e inspecciona que el responder recibió el envelope exacto.

## Hashes

| Archivo | SHA-256 |
|---|---|
| `internal/channels/channel.go` | `e2f0fb39deaf9e34eece4786585609f2c279df107c9a4641cedb480d8f84c068` |
| `internal/channels/registry.go` | `38dec01943776ea47e30b3e82a8aa3892a5a00a831abf97d8381188b8178edce` |
| `internal/channels/dispatcher.go` | `1c6305806edb58878dd5c053ec8fe1b39970a78173a27e3102895d4295160ee0` |
| `internal/channels/channels_test.go` | `12106a7d3cb42ecf59c63bfbf819463b671857214797060363337a543b9b5d60` |

## Límite

Es código local `AUTHORED`. Preserva identidad pero no persiste conversación ni implementa el webhook de cada proveedor; esos claims pertenecen al runtime conectado y al adapter concreto.
