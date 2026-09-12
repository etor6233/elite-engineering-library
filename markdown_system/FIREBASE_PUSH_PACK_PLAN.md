# Firebase Push Pack Plan

Materializa el registro oficial de fuentes y el adaptador Go para Firebase Cloud Messaging. No contiene credenciales ni tokens, inicia bloqueado y usa `SendDryRun` salvo que el proyecto pruebe y apruebe explícitamente un envío real.

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
      "path": "implementation_packs/GO_FIREBASE_CLOUD_MESSAGING_ADAPTER.md",
      "packId": "GO-FIREBASE-CLOUD-MESSAGING-ADAPTER",
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

Antes de ejecutar: licencia oficial, proyecto Firebase, ADC, FCM API, dispositivo consentido, contenido/data, TTL, cuota/costo, retención y reconciliación deben estar `PROVEN`. La primera operación autorizada es `DRY_RUN`; un `SEND` real exige aprobación adicional y nunca escribe estado de negocio automáticamente.
