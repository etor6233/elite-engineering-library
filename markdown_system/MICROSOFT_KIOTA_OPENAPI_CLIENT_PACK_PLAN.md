# Microsoft Kiota OpenAPI Client Pack Plan

Use este perfil sólo cuando el blueprint selecciona un cliente Go generado desde un contrato OpenAPI local, confiable y hash-locked. No forma parte automática de todo proyecto.

```json
{
  "compositionVersion": "1.0",
  "secretVariables": {},
  "packs": [
    {
      "path": "implementation_packs/MICROSOFT_KIOTA_OPENAPI_CLIENT_GATE.md",
      "packId": "MICROSOFT-KIOTA-OPENAPI-CLIENT-GATE",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": ["*"],
      "variables": {}
    }
  ]
}
```

Antes de ejecutar, el agente debe obtener del proyecto: contrato OpenAPI y autoridad propietaria, revisión/licencia y SHA-256; clase/namespace Go; provider/base URL por entorno; esquema de autenticación/secret references; deadlines/retry budget; claves idempotentes; webhooks y reconciliación; sandbox y assertions de contrato; política de regeneración, revisión de diff y rollback. Si falta una decisión irreversible o acceso real, genera sólo evidencia local y conserva la integración como `CONDITIONED`.

La secuencia es:

```text
materializar pack
→ adquirir Kiota exacto
→ verificar contrato local por SHA-256
→ generar en destino ausente + receipt
→ revisar diff y congelar grafo Go
→ compile/test/vet/SCA
→ sandbox contract tests
→ integrar auth/idempotencia/reconciliación
→ E2E y rollback del proyecto
```
