# Amazon SP-API External Fulfillment Inventory Execution Inventory — V157

## Resultado focal

`REBUILD_VERIFIED / CONDITIONED` para el adapter local Amazon External Fulfillment Inventory v2024-09-11.

- Pack 0.1.0: 56.078 bytes; SHA-256 `ae7fddafa65abdc45d23b5df4cc8a0cee9dbb15e73deafcd21d01483ad949f97`.
- Perfil: 2 packs / 34 archivos; 1.329 bytes; SHA-256 `00d5920d107958b6c08a6a70d3b428b63ea2d481afb44356e836bfc5163a461a`.
- Materialización: 9/9 archivos con igualdad SHA-256 entre árbol autor y reconstruido.
- Runtime: wheel oficial `amzn-sp-api==1.11.1`; `pip check` PASS.
- Pruebas: 12/12 PASS en árbol autor y 12/12 PASS reconstruidas desde Markdown.

## Autoridad, licencia y procedencia

Amazon declara que External Fulfillment Inventory permite publicar inventario absoluto location-level desde POS, ERP y WMS para múltiples canales. La operación oficial `batchInventory` admite hasta diez requests y responde por cada ubicación/SKU. El SDK fijado contiene `BatchInventoryApi.batch_inventory` y los modelos exactos usados.

- SDK source: `amzn/selling-partner-api-sdk` commit `8e792ae345a8d334ccdbdd03181f05f040e6a4fc`, Apache-2.0.
- Model source: `amzn/selling-partner-api-models` commit `8e429486005c4ebdce5099e48cc48515a65359bb`.
- El OpenAPI External Fulfillment Inventory declara Amazon Software License; el pack conserva `LicenseRef-Amazon-Software-License` y no copia el modelo.
- Wheel: 3.557.108 bytes; SHA-256 `a88e4059c9fbfd454bfef67b82288a926f73cbdc608a8831f8383565736c70b0`.
- Los nueve archivos son integración/configuración `AUTHORED`; no se presentan como producto escrito por Amazon.

## Contratos demostrados

1. Perfil distribuido bloqueado: programa, rol/scope, canales, mapping, sequence authority, sandbox, cuota, privacidad y owner son obligatorios.
2. Batch estricto de 1–10 identidades únicas; sólo canales oficiales aprobados `FBA`, `MFN`, `DF` y cantidades/sequences enteras acotadas.
3. URI `/inventory/fetch|update` codifica location y SKU; los requests usan `BatchInventoryRequest`, `InventoryRequest`, `InventoryRequestParams`, `MarketplaceAttributes` y `HttpMethod.POST` oficiales.
4. FETCH persiste hashes de location/SKU/marketplace, canal, sequence, sellable y reserved; nunca IDs crudos.
5. UPDATE exige sequence por identidad exacta, approval sobre request+receipt y nuevo FETCH sin drift antes del efecto.
6. Attempt se escribe antes de `batch_inventory`; rechazo total se registra sanitizado, sin reason phrase ni IDs.
7. Batch parcial, excepción o sequence no reconciliada produce `UNKNOWN_EFFECT`; retry automático queda prohibido.
8. Reconciliación por FETCH demuestra sequence provider. La cantidad absoluta solicitada no se promueve automáticamente a hecho ERP/POS/WMS.

## Fallos convertidos en memoria

- `LIB-FAIL-1484`: glob POSIX literal rechazado en Windows; repetición correcta con directorio + `-g '*.py'`.
- `LIB-FAIL-1485`: el fake aplicó fallos de UPDATE al FETCH previo; se restringió la inyección al URI/número de llamada exactos y la suite completa pasó.

## Estado global confirmado

- `VERIFY_LIBRARY_PASS`: 94 packs / 1015 archivos materializables / 513 Markdown / 42 perfiles.
- `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit`: 94 packs / 121 fuentes oficiales / 14 adapters de proveedor.

## Límites honestos

No se ejecutó cuenta, programa, sandbox o inventario real. No demuestra semántica live de sequence, mapeo ERP/POS/WMS, reservas/transfers internos, concurrencia, carga, seguridad o producción. Los gates de red/credenciales/cuenta/target quedan `SKIPPED`, no simulados. CDN/WAF, IdP, proveedores, PostgreSQL/recovery, carga, seguridad ofensiva, deploy/rollback y aceptación empresarial siguen siendo gates del proyecto concreto.
