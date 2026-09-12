# TikTok Ads Reporting Pack Plan

Materializa el wheel lock oficial y el adapter read-only `ReportingApi.report_integrated_get` sobre `tiktok-business-api-sdk-official==1.1.3`. No contiene credenciales y no habilita llamadas por sí solo.

```json
{
  "compositionVersion": "1.0",
  "secretVariables": {},
  "packs": [
    {
      "path": "implementation_packs/OFFICIAL_UPSTREAM_ACQUISITION_CORE.md",
      "packId": "OFFICIAL-UPSTREAM-ACQUISITION-CORE",
      "version": "0.4.88",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/PYTHON_TIKTOK_ADS_REPORTING_ADAPTER.md",
      "packId": "PYTHON-TIKTOK-ADS-REPORTING-ADAPTER",
      "version": "0.2.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    }
  ]
}
```

Antes de ejecutar: Business Center/developer app, reporting scope, test advertiser, dimensions/metrics, fechas, retención, quota/costo, términos y reconciliación deben estar `PROVEN`. La plantilla conserva `automatic_business_write=false`; el wheel oficial 1.1.3 se instala sólo por su URL/SHA-256 exactos.
