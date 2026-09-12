# Meta Ads Reporting Pack Plan

Materializa el source lock oficial y el adapter read-only `AdAccount.get_insights` sobre Meta Business SDK 26.0.1/Graph API v26.0. No contiene credenciales ni habilita llamadas por sí solo.

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
      "path": "implementation_packs/PYTHON_META_ADS_REPORTING_ADAPTER.md",
      "packId": "PYTHON-META-ADS-REPORTING-ADAPTER",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    }
  ]
}
```

Antes de ejecutar: app/business registration, `ads_read`, test ad account, ownership, field allowlist, dates, retention, quota/cost, Meta terms y reconciliation deben estar `PROVEN`. La plantilla conserva `automatic_business_write=false`. La licencia restringe los SDKs a servicios/APIs Facebook y el LICENSE source de CAPI debe acompañar cualquier redistribución.
