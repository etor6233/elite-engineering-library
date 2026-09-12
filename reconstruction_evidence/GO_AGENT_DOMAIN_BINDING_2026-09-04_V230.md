# Go Agent Domain Binding — evidencia V230

## Resultado

```yaml
pack: GO-AGENT-DOMAIN-BINDING
version: 0.2.0
implementation: REBUILD_VERIFIED
admission: CONDITIONED
verified_at: 2026-09-04
tests: PASS
materialization: 3/3 PASS
gofmt_idempotence: 3/3 PASS
full_go_packages: 47 PASS
go_vet: PASS
go_build: PASS
pack_sha256: 77569144728a48ec21739471d6d79d0607a54764c2d82018d517102baf2204d7
```

## Corrección demostrada

`CommandIdentity` obliga a que el ingress —nunca el LLM— entregue clave estable y timestamp del mensaje. Así, fechas relativas de citas y vigencia de cotizaciones conservan el mismo body en cada retry. El gateway propaga `Idempotency-Key`, rechaza tenant cruzado y obtiene el Bearer actual en cada request protegido mediante `AccessTokenProvider`; token ausente o inválido falla cerrado.

Las pruebas HTTP inspeccionan ruta, body, clave, Bearer y rotación de token. Las firmas heredadas sin identidad estable se conservan sólo para compatibilidad de compilación y devuelven `ErrIdempotencyRequired`; no son el recorrido permitido del runtime conectado.

## Hashes

| Archivo | SHA-256 |
|---|---|
| `internal/domainbind/config.go` | `ee62027dac5f728385c0d8314d4f7d81d751a6568e9df8f782a3783f141f5a8c` |
| `internal/domainbind/gateway.go` | `cd4b8c4ac13a810ef0a0f27b3904e18721d27d67d1a92828234561d066bcbcc9` |
| `internal/domainbind/gateway_test.go` | `f7ceb8fb5e5114e2d0c8f5907f2f54c2a33c73a5d16924094f0dee755e5c9598` |

## Límite

El pack es código local `AUTHORED`. Demuestra el borde HTTP, no adquisición OAuth/OIDC del proyecto ni un journey live. El caller debe aportar identidad estable y un proveedor de credencial aprobado.
