# Go Channels Core — evidencia V234

## Resultado

```yaml
pack: GO-CHANNELS-CORE
version: 0.3.0
implementation: REBUILD_VERIFIED
admission: CONDITIONED
verified_at: 2026-09-04
materialization: 4/4 PASS
gofmt_idempotence: 4/4 PASS
full_go_packages: 47 PASS
pack_sha256: c4715391223ccde74007f5387cee41a1aef8315a8015d6789755975a3b945659
```

## Entrega demostrada

V232 preservó el envelope ingress. V234 agrega `DeliveryKey` obligatoria al outbound y la deriva como SHA-256 de tenant, canal, provider message ID y propósito reply. La prueba demuestra que repetir el mismo inbound conserva la key y otro message ID produce otra.

Esto no afirma exactly-once de la red: el adapter concreto debe usar la key del proveedor cuando exista o implementar receipt/reconciliación durable. Sí garantiza que el core no inventa una key distinta en cada retry.

## Hashes

| Archivo | SHA-256 |
|---|---|
| `internal/channels/channel.go` | `2bd88fc11511fac919d8313b60784ee17b9d86c9fbe81877b3d49d0661624a0a` |
| `internal/channels/registry.go` | `38dec01943776ea47e30b3e82a8aa3892a5a00a831abf97d8381188b8178edce` |
| `internal/channels/dispatcher.go` | `eebb41f403b53f3dfba4df644e4d391e9856d23b5e731a79320717b86571e857` |
| `internal/channels/channels_test.go` | `5d8f5ff50abc0dca1ac89b66d2d76142d9268e28e1c157d489ab3bbbeab8d85c` |
