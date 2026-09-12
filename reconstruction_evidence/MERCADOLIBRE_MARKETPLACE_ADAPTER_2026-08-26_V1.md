# Mercado Libre Marketplace Adapter — Reconstruction Evidence V1

Date: 2026-08-26  
Result: `REBUILD_VERIFIED / CONDITIONED`  
Pack SHA-256: `a895f1356116803390b2922c67a2392a411257f639023ebe6e840e6a4270c168`

## Official authority and honest provenance

Mercado Libre's official Go, Node, PHP, Java and .NET marketplace SDK repositories are archived; the Go repository explicitly says it is no longer maintained and not functional. It is therefore rejected as a runtime dependency. Current official developer documentation instead governs OAuth bearer headers, `GET /items/{id}`, `GET /orders/{id}`, `POST /items/validate`, read/write scopes and notification queue/fetch behavior. The provider's official hosted MCP server was also reviewed; its public repository currently documents a hosted documentation-assistance surface rather than seller-operation code.

All ten materialization blocks are `AUTHORED` integration code. No community SDK and no archived provider SDK was copied or relabeled as current official code.

## Clean reconstruction

- Official Go 1.26.7 Windows amd64 ZIP: 74,955,002 bytes; SHA-256 `f4f534a486e4bc3387fa18f08208f2f854b7aaea8a08f2a2d829a914a05abb11`.
- Isolated source and reconstruction roots under `%LOCALAPPDATA%/Temp/meli-marketplace-audit-20260826`.
- Ten files materialized and matched the source tree byte-for-byte.

| Path | Bytes | SHA-256 |
|---|---:|---|
| `mercadolibre_marketplace/authority.lock.json` | 1,659 | `82cc1e9f15372d1359cb038f4a6547d5bd8dda801a94e1d5c87b923d858ea44d` |
| `mercadolibre_marketplace/cmd/meli-marketplace/main.go` | 1,131 | `1d090bcde11d76d5f96f9160e664fdf19f2f5575a3ee9a8e9c3e17caad757b12` |
| `mercadolibre_marketplace/go.mod` | 54 | `c942e5ce737c8c1e783b3f908bf61aaabc73470bbb86eeedaf5ceee7e38114fa` |
| `mercadolibre_marketplace/marketplace_test.go` | 4,902 | `ea42258463afe30b7fd9c2898681b99ceb54f16e79ec38d2888db87bf8c7d89e` |
| `mercadolibre_marketplace/marketplace.go` | 9,042 | `62bb3523752450b46eeb60e26cf981ff4f35ce1deb71064a358e61a45ce1ca58` |
| `mercadolibre_marketplace/notifications_test.go` | 2,101 | `8fd7cf2fb73c6725117bb274963ee0f3e6aad5e7c64ce5f6569b14681c74f2f4` |
| `mercadolibre_marketplace/notifications.go` | 3,760 | `d2e266841434771bb9744de6305c2547944d0f26dfef65b2b8421d7573458ed5` |
| `mercadolibre_marketplace/provider-profile.template.json` | 898 | `5d4ce698b2a3300d76a749f39af772b084d31de181f4bb95888bb065b89bd62e` |
| `mercadolibre_marketplace/README.md` | 1,403 | `4f9f0320da3164b43105579640796d8460948d1adbd66de9f6b0167baab5dee4` |
| `mercadolibre_marketplace/request.template.json` | 64 | `8b8b4003ff34b19fd9ec996bedfb07f55460a158feb5770d8edcc98434270ce3` |

## Executed gates

```text
go fmt ./...                                  PASS
go mod tidy                                   PASS; no third-party module
go mod verify                                 PASS: all modules verified
go test ./...                                 PASS: 8 tests
go vet ./...                                  PASS
go build ./cmd/meli-marketplace               PASS
materialize_markdown_pack.ps1                 PASS: 10 files
source vs reconstruction SHA-256              PASS: 10/10 identical
VERIFY_EXECUTABLE_LIBRARY -Mode Foundation     PASS in V10
```

Tests enforce the exact production base URL and Bearer header, item/order paths, validator POST with 204/400 semantics, rejection of unimplemented publication, profile/write/token/ID/JSON/status failures, atomic provider failures, notification seller/application/topic/resource/origin allowlists, durable fetch jobs and identifier redaction.

## Conditions and remaining work

No Mercado Libre account, seller grant, token, test user, live validator, callback or missed-feed endpoint was contacted. The provider documents no preproduction marketplace environment, so real changes demand additional care. This V1 deliberately excludes publish/update/delete, stock/price, OAuth refresh concurrency, questions, claims, shipment writes and every Mercado Pago operation. Those are separate conditioned capabilities, not silently implied by this pack.
