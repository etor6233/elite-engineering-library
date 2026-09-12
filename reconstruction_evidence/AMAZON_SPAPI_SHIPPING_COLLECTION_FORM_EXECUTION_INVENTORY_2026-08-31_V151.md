# Amazon SP-API Shipping collection form — execution inventory V151

Fecha: 2026-08-31. Resultado: `REBUILD_VERIFIED / CONDITIONED`.

## Autoridad

SDK oficial Amazon `amzn-sp-api` 1.11.1, commit `8e792ae345a8d334ccdbdd03181f05f040e6a4fc`, modelos `8e429486005c4ebdce5099e48cc48515a65359bb`, Apache-2.0, wheel SHA-256 `a88e4059c9fbfd454bfef67b82288a926f73cbdc608a8831f8383565736c70b0`.

Se inspeccionaron `GetUnmanifestedShipmentsRequest`, `ClientReferenceDetail`, `GenerateCollectionFormRequest`, `CollectionsFormDocument` y las firmas oficiales `get_unmanifested_shipments`/`generate_collection_form`; la segunda admite `x-amzn-IdempotencyKey`. Referencias: https://developer-docs.amazon.com/sp-api/reference/getunmanifestedshipments y https://developer-docs.amazon.com/sp-api/reference/generatecollectionform.

Los cuatro archivos nuevos son `AUTHORED`, no código Amazon; construyen y verifican las superficies oficiales fijadas.

## Artefacto

- Pack 0.5.0: 116.024 bytes/SHA-256 `98f8786c3872687c8a6d7d5d33c84e1b5cc4b4d44f8992ff1a98e470a05ffb39`; 20/20 archivos.
- Request template 213/`b587b2d58711a18794ba18ff49abd4946f8809279a734530b953f67c915be2b8`.
- Approval template 204/`d330f11b825e02f5e797fd814d93971749322ee9b45060b679116fdad67bfd0e`.
- Runner 14.360/`c17ad986bcf7e78eee7e59d918c78c756d2435c94b8b1799f24d4bed339043c5`; tests 6.133/`5669972af6cd3efd7606949b437cb28b4447041e185ff42d43b11c580e3267d3`.
- Perfil 2/45; record 13.802/`441807af873284c926fb050497f060223737d12356d9936785a271f64d4da5f0`.

## Contrato demostrado

Inventory construye el request oficial con 1–2 referencias tipadas, conserva respuesta/receipt atómicos y no genera form. Generate exige request original, hashes exactos, aprobación y una única combinación carrier/address presente en inventory. Persiste attempt e idempotency key antes del POST. Sólo acepta PDF Base64 acotado, lo nombra por SHA-256 y conserva respuesta+receipt. Excepción o respuesta inválida deja `UNKNOWN_EFFECT`, sin falso éxito.

El receipt fija `physical_handover_completed=false`: collection form no demuestra pickup, escaneo, custodia o aceptación del carrier.

## Gates

28/28 tests PASS en árbol autor y reconstruido; `compileall`, `pip check`, 20/20 materialización y perfil 2/45 PASS. Negativos cubren referencia inválida, output existente, tamper de request/receipt/address, timeout ambiguo y Base64 no admitido.

## Condiciones abiertas

Cuenta/sandbox, referencias reales, carrier/address, idempotency/reconciliación live, impresión/retención, pickup booking, escaneo y firma de handover, claims, cross-location, rate headers/carga/alarmas/aceptación siguen bloqueados. V151 no prueba handover físico ni entrega.
