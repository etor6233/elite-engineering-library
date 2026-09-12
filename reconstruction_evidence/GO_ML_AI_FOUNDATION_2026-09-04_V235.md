# Go ML/AI Foundation — evidencia V235

## Resultado

```yaml
pack: GO-ML-AI-FOUNDATION
version: 0.2.0
implementation: REBUILD_VERIFIED
admission: CONDITIONED
verified_at: 2026-09-04
materialization: 10/10 PASS
gofmt_idempotence: 10/10 PASS
full_go_packages: 47 PASS
pack_sha256: 63dfcebeafd03e318e5931756261059ac1dc81657c2583025abd399454a14bd3
```

## Corrección de seguridad

La allowlist de PII dejó de ser un bypass global. Pre-flight evalúa siempre las frases de inyección antes de aplicar la excepción PII; post-flight evalúa siempre frases de fuga aun cuando la salida esté allowlisted. Las regresiones cubren ambos bypasses. El pack sigue declarando que estas heurísticas no sustituyen clasificadores/recognizers admitidos en producción.

## Hashes modificados

| Archivo | SHA-256 |
|---|---|
| `internal/aifoundation/safety.go` | `d3ebfbb6c2cf3a6aa36676ac3a9e1ff5e2cfbc6abc82da80135a2e6f7227e753` |
| `internal/aifoundation/safety_test.go` | `c4f353d2715ce754353c48f3d53ad70df21ccb7cfdae357e91fe500c968d09b2` |
| `internal/aifoundation/eval_test.go` | `d28c3b2b033735b01c2d97037284dda08342789813ef38f1807292951975bf41` |

`eval_test.go` recibió sólo normalización `gofmt`, registrada separadamente como `LIB-FAIL-1826`.
