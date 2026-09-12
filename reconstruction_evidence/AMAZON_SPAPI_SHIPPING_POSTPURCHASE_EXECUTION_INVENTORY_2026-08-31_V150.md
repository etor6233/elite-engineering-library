# Amazon SP-API Shipping post-purchase — execution inventory V150

Fecha: 2026-08-31. Resultado: `REBUILD_VERIFIED / CONDITIONED`.

## Autoridad oficial exacta

Amazon SDK Python 1.11.1, commit `8e792ae345a8d334ccdbdd03181f05f040e6a4fc`, y modelos commit `8e429486005c4ebdce5099e48cc48515a65359bb`, ambos Apache-2.0. Wheel oficial SHA-256 `a88e4059c9fbfd454bfef67b82288a926f73cbdc608a8831f8383565736c70b0`.

Firmas generadas inspeccionadas: `ShippingApi.cancel_shipment(shipment_id, x_amzn_shipping_business_id=...)` realiza PUT y no admite idempotency key; `ShippingApi.get_shipment_documents(shipment_id, package_client_reference_id, format=..., dpi=..., x_amzn_shipping_business_id=...)` devuelve `GetShipmentDocumentsResult`, cuyo `PackageDocument.contents` es Base64. Referencias: https://developer-docs.amazon.com/sp-api/reference/cancelshipment y https://developer-docs.amazon.com/sp-api/reference/getshipmentdocuments.

Los tres archivos V150 son `AUTHORED`, no código Amazon. Ejecutan exclusivamente los métodos/modelos oficiales fijados.

## Inventario reconstruido

- Pack 0.4.0: 92.353 bytes/SHA-256 `ad12f6e77d4819ecd516dddb9b5f0aafe2a6c82ae2794adc01bbc9a83bbadbf6`; 16/16 archivos.
- `cancellation-approval.template.json`: 183/`dc87f254ec001397fcf07b7383a5c34595b9cefa69f26bffbd1e0e0c1fc2a214`.
- `run_shipping_postpurchase.py`: 14.943/`79c2aafae920fb362619e19918496e13c65885060df672ab731d52d04e11fa5d`.
- `test_shipping_postpurchase.py`: 7.527/`9c743d6eb239b36b0031f2d48e5c095649fa397db068487cdcd89f4194c102d5`.
- Perfil: 2 packs/41 payload files; record 12.560 bytes/SHA-256 `b4d46edb3a5aed12269d74fa687393fd5bbf15d01cb7b8e82e4fb27015c30390`.

## Comportamiento demostrado

Cancel no acepta shipment ID desde CLI: verifica receipt+respuesta V149 y extrae el ID hash-bound. La aprobación liga el receipt exacto y razón. El attempt queda durable antes del PUT; éxito exige el payload vacío documentado. Toda excepción posterior produce `CANCELLATION_UNKNOWN_EFFECT`, prohíbe retry automático y exige reconciliación.

Documents verifica el mismo purchase y la selección V149 exacta; package ID debe existir una vez en la respuesta comprada. Format/dpi provienen de esa selección. Respuesta exige shipment/package exactos; Base64, type, format, cardinalidad y límites se validan antes de commit. Los binarios se nombran por hash y se publican junto a receipt de forma atómica.

## Gates

22/22 tests PASS desde árbol autor y Markdown materializado; `compileall`, `pip check`, 16/16 materialización y perfil 2/41 PASS. Negativos cubren tamper, output repetido, timeout ambiguo, selección/package divergentes y Base64 inválido sin residuo.

## Condiciones abiertas

No hubo cuenta/sandbox real. Acceso/rol/región/business/carrier, costo, privacidad, comportamiento live de cancelación, reconciliación, retención/impresión de etiquetas, rate headers, carga, alarmas y aceptación siguen bloqueados. Handover/manifest físico, pickup, cross-location, claims y delivery proof no quedan cubiertos por V150.
