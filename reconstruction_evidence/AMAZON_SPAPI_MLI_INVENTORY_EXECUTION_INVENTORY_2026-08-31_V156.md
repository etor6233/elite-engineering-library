# Amazon SP-API Multi-Location Inventory Execution Inventory — V156

## Resultado focal

`REBUILD_VERIFIED / CONDITIONED` para el adapter local Amazon Multi-Location Inventory mediante Listings Items v2021-08-01.

- Pack 0.1.0: 56.116 bytes; SHA-256 `1bc2c81c3b5299e69ac640dc26aa37303c2e174d77c2e18c1d286bb99dd0518c`.
- Perfil: 3 packs / 43 archivos; 1.625 bytes; SHA-256 `4def050d29b67d71291ede2c65c93badf0636e8452e9f04bf41555d658c4f120`.
- Materialización: 9/9 archivos con igualdad SHA-256 entre árbol autor y reconstruido.
- Runtime: wheel oficial fijado `amzn-sp-api==1.11.1`; `pip check` PASS.
- Pruebas: 13/13 PASS en árbol autor y 13/13 PASS desde Markdown.

## Autoridad y procedencia

Amazon documenta que Multi-Location Inventory publica inventario FBM por ubicación mediante Supply Sources y Listings Items/Feeds. Para Listings Items exige `productType=PRODUCT` y path `/attributes/fulfillment_availability`; GET devuelve Supply Source ID en `fulfillmentChannelCode` y cantidad. Amazon declara expresamente que MLI no está disponible para FBA inventory ni fuera de US en la guía observada, por lo que este pack no relabela FBA Inventory como inventario por sucursal.

- SDK: `amzn/selling-partner-api-sdk` commit `8e792ae345a8d334ccdbdd03181f05f040e6a4fc`.
- Modelos: `amzn/selling-partner-api-models` commit `8e429486005c4ebdce5099e48cc48515a65359bb`.
- Wheel: 3.557.108 bytes; SHA-256 `a88e4059c9fbfd454bfef67b82288a926f73cbdc608a8831f8383565736c70b0`.
- Licencia upstream: Apache-2.0.
- Los nueve archivos son integración/configuración `AUTHORED`; no se presentan falsamente como producto escrito por Amazon.

## Contratos demostrados

1. Perfil distribuido bloqueado: registro, app, rol Product Listing, inscripción MLI, receipt Supply Sources, sandbox, cuota, datos y owner son obligatorios.
2. Sólo Supply Sources únicos y `Active` del receipt hash-validado pueden aparecer en un request; cantidades y cardinalidad están acotadas.
3. Snapshot usa `get_listings_item` con `includedData=fulfillmentAvailability`, liga target/request/receipt y no persiste seller, SKU, marketplace ni source IDs crudos.
4. Approval se liga a bytes exactos de request, receipt Supply Sources y snapshot Listing; una alteración o drift live bloquea antes del PATCH.
5. La mutación construye modelos oficiales `ListingsItemPatchRequest` y `PatchOperation(op="replace", path="/attributes/fulfillment_availability")` con `product_type="PRODUCT"`.
6. `MUTATION_ATTEMPT.json` existe antes del efecto; respuesta provider no aceptada queda sanitizada sin mensajes ni submission ID.
7. Aceptación se reconcilia mediante GET acotado; excepción o divergencia crea `UNKNOWN_EFFECT` sin mensaje sensible y sin retry automático.
8. Output existente bloquea antes de llamar al provider; `automatic_internal_inventory_commit=false` queda fijo.

## Fallos convertidos en memoria

- `LIB-FAIL-1481`: el setter del modelo oficial `Issue` exige `categories`; el fixture se corrigió sin tocar el SDK ni relajar el adapter.
- `LIB-FAIL-1482`: un array no sobrevivió `pwsh -File`; el updater se invocó en la sesión actual y actualizó 9/9 bloques.
- `LIB-FAIL-1483`: la primera suite reconstruida se lanzó desde el workspace; se repitió desde el `cwd` materializado exacto y pasó 13/13.

## Estado global

- `VERIFY_LIBRARY_PASS`: 93 packs / 1006 archivos materializables / 510 Markdown / 41 perfiles.
- `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit`: 93 packs / 121 fuentes oficiales / 13 adapters de proveedor.

## Límites honestos

No se ejecutó cuenta, dynamic sandbox, marketplace, inscripción MLI o listing real. No demuestra consistencia temporal provider, ERP/POS/WMS, promesas, transfers, reservas, concurrencia, carga, seguridad o producción. Los gates que requieren red, credenciales, cuenta o target real permanecen `SKIPPED`; no se simulan ni cuentan como evidencia productiva. El proyecto conserva gates live para CDN/WAF, IdP, proveedores, PostgreSQL/recovery, carga, seguridad ofensiva, deploy/rollback y aceptación empresarial.
