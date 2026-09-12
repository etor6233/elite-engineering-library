# Amazon SP-API Shipping claim — execution inventory V152

Fecha: 2026-08-31. Resultado: `REBUILD_VERIFIED / CONDITIONED`.

## Autoridad

Amazon SDK Python 1.11.1, commit `8e792ae345a8d334ccdbdd03181f05f040e6a4fc`, modelos `8e429486005c4ebdce5099e48cc48515a65359bb`, Apache-2.0; wheel oficial SHA-256 `a88e4059c9fbfd454bfef67b82288a926f73cbdc608a8831f8383565736c70b0`.

Se inspeccionaron `CreateClaimRequest`, `CreateClaimResponse`, `ClaimReason`, `SettlementType`, `Currency` y `ShippingApi.create_claim(body, x_amzn_shipping_business_id=...)`. La API no expone idempotency key. Referencia: https://developer-docs.amazon.com/sp-api/reference/createclaim.

Los tres archivos V152 son `AUTHORED`; no se atribuyen a Amazon.

## Inventario

- Pack 0.6.0: 134.089 bytes/SHA-256 `a1afe97159d1fc16b582d6c14b12519f4ea388c0ffaab40f79a26570ae28ab99`; 23/23 archivos.
- Approval 359/`7dca1a6a4b04bda9714939b3e5a4273e7f262f7eb9cf292516ae76ad9a163260`.
- Runner 10.379/`0f2211c6084dfab5da0bf1fd13a92031aaef8f7c4e6814512b7d377840e968e3`; tests 5.396/`0737b5f327363a5fa068b6834298ea122dcbd7148916ed4ede7e7824d2b3b6a8`.
- Perfil 2/48; record 14.678/`285134cca64fe43559284f19235b624844c86dc08ce50c652333981c4d7b67b0`.

## Contrato demostrado

Claim requiere purchase receipt/respuesta, quote receipt/respuesta, rate request y purchase selection mutuamente hash-bound. Package/tracking salen del purchase; insured value/moneda salen del rate request. Approval selecciona razón/settlement/valor/proofs y el valor no puede exceder el asegurado. Proofs son HTTPS, sin credenciales/query/fragment y con host exacto aprobado.

Antes del POST se persiste attempt. Éxito exige claim ID oficial y lo conserva sólo hasheado en receipt. Como no hay idempotency key oficial, excepción o respuesta inválida deja `CLAIM_UNKNOWN_EFFECT`, `automatic_retry_authorized=false` y reconciliación obligatoria.

## Gates

33/33 pruebas PASS desde autor y Markdown; `compileall`, `pip check`, 23/23 materialización y perfil 2/48 PASS. Negativos: tamper, valor excesivo, proof host inválido, output existente, timeout ambiguo.

## Condiciones abiertas

Cuenta/sandbox, eligibility/plazos legales, evidencia real, privacidad/retención, settlement/reconciliación live, rate headers, carga, alarmas y aceptación siguen bloqueados. Claim no prueba resolución, pago, handover ni entrega. Pickup/escaneo/firma y carrier universal permanecen fuera de V152.
