# Reconstruction Evidence — TS Enterprise Web V3

```yaml
evidence_id: "TS-EW-20260824-V3"
pack_id: "TS-ENTERPRISE-WEB"
pack_version: "0.1.0"
pack_sha256: "1019d71f5c5b8ada20bb9bf24cdbedac0439dd27b6b8c51440be775a89f3f719"
materialized_file_count: 63
ordered_file_hash_aggregate_sha256: "b6be34b5f3d73fcf2904ac771aaf0680a573c09bae90ad5569377b683d099a91"
environment: "Windows / clean target directory rebuild_verification/ts-enterprise-web-v3"
node: "v24.14.1"
pnpm: "11.19.0"
verified_at: "2026-08-24"
result: "PASS"
```

## Resultados

| Gate | Resultado |
|---|---|
| materialización + 63 hashes | PASS |
| frozen lockfile / 89 packages | PASS |
| configuración y readiness baseline | VALID / READY |
| TypeScript | PASS |
| tests | 9 archivos / 22 tests / PASS |
| Next.js production build | PASS |
| production dependency audit | no known vulnerabilities |
| páginas pública/modelos/admin/customer/factory | HTTP 200 |
| APIs health/platform | HTTP 200 |
| lead creation / idempotent replay | HTTP 201 / 200 |

## Cobertura agregada respecto de V2

- supplier y purchase order con líneas y totales seguros;
- shipment/recepción de fábrica e inventario serializado;
- cliente, pedido, transiciones y reserva atómica de stock;
- audit y outbox en las operaciones empresariales;
- verificación OIDC de firma, issuer, audience, edad y claims organizacionales;
- adapter HTTP para proveedores con HTTPS, allowlist por base URL, timeout, límites, no redirects, idempotencia y clasificación de retry.

## Condiciones no simuladas

No se declara producción universal. Cada proyecto debe aportar configuración de negocio real, IdP y JWKS reales, PostgreSQL/roles, credential provider, contratos sandbox por proveedor, observabilidad exportada, deployment y recovery demostrado.

