# Python Google Ads Reporting Adapter — reconstruction evidence V1

## Alcance y autoridad

Pack `PYTHON-GOOGLE-ADS-REPORTING-ADAPTER` 0.1.0. Usa el wheel oficial Apache-2.0 `google-ads==31.3.0` de 23.675.722 bytes, SHA-256 `286946249d54dea97ad4c81aa5d466a49060cc726f0846281a4e756d24de9890`, y API v25. El glue/tests es `AUTHORED`; no se atribuye a Google.

## Reconstrucción limpia

```text
Materialized 5 files into a new directory
test_empty_report_is_valid_evidence ... ok
test_fails_closed_and_atomic ... ok
test_preserves_rows_and_redacts_customer ... ok
Ran 3 tests in 0.012s
OK
GOOGLE_ADS_SDK_CONTRACT_PASS version=31.3.0 api=v25
```

Se probaron GAQL SELECT único, `search_stream`, preservación de filas, empty report válido, customer/query/response hashes, customer redacted, `automatic_business_write=false`, output atómico y negativos de query/customer/output. No hubo llamada Google.

## Fallos y condiciones

`LIB-FAIL-042` corrigió el lock inexistente 31.4.0. `LIB-FAIL-043` retiene el proceso combinado de instalación/probe que no terminó; metadata, import y firma pasaron por procesos separados. El grafo transitivo fue instalado pero todavía debe congelarse/SBOM/scanearse como artifact de proyecto.

Faltan developer token/access level, OAuth, customer/login hierarchy, cuenta test, query allowlist real, política de campos, cuota/costo, retención, logging y reconciliation. Mutaciones de campaña y conversion uploads no están implementados. Estado `CONDITIONED`.

## Fuentes oficiales

- https://pypi.org/pypi/google-ads/31.3.0/json
- https://github.com/googleads/google-ads-python/tree/f7bf312d26904ea50ad9ef02e395edccfaa1ebb7
- https://developers.google.com/google-ads/api/docs/client-libs/python
- https://developers.google.com/google-ads/api/docs/query/overview
