# Google Ads Reporting Pack Plan

Materializa el lock oficial y el adapter read-only Google Ads. No contiene credenciales ni queries reales y no habilita llamadas por sí solo.

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
      "path": "implementation_packs/PYTHON_GOOGLE_ADS_REPORTING_ADAPTER.md",
      "packId": "PYTHON-GOOGLE-ADS-REPORTING-ADAPTER",
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

Antes de ejecutar: developer token/API access, OAuth identity reference, customer/login hierarchy, test account, GAQL allowlist, fields/data policy, quota/cost, retention y reconciliation deben estar `PROVEN`. El profile sigue `BLOCKED_ACCESS_AND_QUERY_APPROVAL_REQUIRED` y `automatic_business_write=false`.
