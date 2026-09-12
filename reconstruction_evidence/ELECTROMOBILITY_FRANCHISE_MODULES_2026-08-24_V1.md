# Electromobility Franchise Modules — Reconstruction Evidence V1

```yaml
evidence_id: "ELECTROMOBILITY-MODULES-20260824-V1"
pack_id: "ELECTROMOBILITY-FRANCHISE-MODULES"
pack_version: "0.1.1"
pack_sha256: "ea45379cd977461f034c412842734850e65802f68529e17d635984a53af76935"
materialized_file_count: 3
ordered_file_hash_aggregate_sha256: "b97695149341a2f5f9738b06ed93c91ef624dfe6178c4aaf20d08a67639ebe04"
database: "PostgreSQL 18.6 windows/amd64"
result: "PASS_CONDITIONED"
verified_at: "2026-08-24"
```

## Cobertura material

La migration creó 14 schemas, 23 tablas y 7 índices para catálogo/variantes, cliente/lead/consentimiento, proveedores/compras, fábrica, stock serial/VIN/batería, pricing, líneas de venta, pagos, logística, garantía/service/recall, franquicias, mapping/reconciliation, marketing y mensajes.

```text
0001 platform + 0002 organization/order       PASS
0003 up                                        PASS
fixture vertical end-to-end de datos           PASS
rechazo de serial duplicado                    PASS
0003 down                                      PASS
0003 up + test nuevamente                      PASS
```

La primera versión del down dejó `sales.customer_order_line` porque vive en un schema compartido; el siguiente up falló cerrado. Se corrigió el ownership del down, se repitió desde estado conocido y recién entonces se emitió esta evidencia 0.1.1.

## Límite

No hay aún APIs/UI/workers para las 23 tablas, ni reglas legales/fiscales/homologación/garantía por mercado. Faltan carreras de stock, reconciliación de providers, roles DB, volumen, performance y restore. El pack es una foundation vertical verificable, no un sistema terminado.

