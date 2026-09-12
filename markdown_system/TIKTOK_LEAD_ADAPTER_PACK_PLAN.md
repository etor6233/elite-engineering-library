# TikTok Lead Adapter Pack Plan

Materializa el wheel oficial TikTok Business API 1.1.3, la adaptación mínima del contrato Lead v1.3 y el importador hash-linked que reutiliza el owner PostgreSQL/outbox omnicanal. No contiene credenciales, no autentica por sí sola el webhook y no autoriza contacto automático.

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
      "path": "implementation_packs/GO_OMNICHANNEL_LEAD_INGRESS.md",
      "packId": "GO-OMNICHANNEL-LEAD-INGRESS",
      "version": "0.2.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/PYTHON_TIKTOK_LEAD_ADAPTER.md",
      "packId": "PYTHON-TIKTOK-LEAD-ADAPTER",
      "version": "0.2.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GO_TIKTOK_LEAD_DURABLE_IMPORT.md",
      "packId": "GO-TIKTOK-LEAD-DURABLE-IMPORT",
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

Antes de ejecutar: cuenta/app/Lead access, términos/consentimiento, test account, callback/TLS, delivery, deduplicación durable, retrieval/reconciliation, retención/borrado y cuota/costo deben estar `PROVEN`. El adapter exige el profile en su constructor; el webhook queda `UNAUTHENTICATED_PROVIDER_SIGNAL`. Sólo el resultado del GET autenticado genera tres artefactos que el importador verifica antes de escribir raw/candidato/outbox; contacto sigue `pending_policy`.
