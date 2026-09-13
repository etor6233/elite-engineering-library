# Meta WhatsApp Cloud Pack Plan

Materializa el source lock oficial y la adaptación explícita de WhatsApp Cloud API. Incluye tres archivos de código Meta byte-verbatim, LICENSE con sólo newline final normalizado, provenance exacta y una ruta segura template/webhook. No contiene credenciales ni habilita llamadas por sí solo.

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
      "path": "implementation_packs/PYTHON_META_WHATSAPP_CLOUD_ADAPTER.md",
      "packId": "PYTHON-META-WHATSAPP-CLOUD-ADAPTER",
      "version": "0.15.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    }
  ]
}
```

Antes de ejecutar: licencia/Platform Policy, Business/app/phone, Graph API version vigente, template/language/arity, recipient consent/opt-out, webhook subscription, secret references, retención, quota/costo e idempotencia/reconciliación deben estar `PROVEN`. Los tres ejemplos oficiales son referencia, no la ruta productiva; el SDK Node archivado permanece excluido.

V402 composition update: exact FX snapshot receipt owner selected with its host; connected WhatsApp/OIDC/social owners selected only where their full application dependencies are present. The shared BC MIT license has one composed owner. Local fixture and source gates are scoped in COMMUNICATIONS_RUNTIME_V402.md; this is not a production claim.

V402 scheduled WhatsApp local claim; SCHEDULED_COMMUNICATIONS_RELEASE_V402.md/json. T2805 campaign work remains open.
