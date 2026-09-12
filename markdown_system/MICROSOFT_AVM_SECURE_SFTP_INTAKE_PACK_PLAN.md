# Microsoft AVM Secure SFTP Intake Pack Plan

Este perfil opt-in materializa ocho archivos: una composición privada de Azure Blob SFTP sobre Microsoft Azure Verified Modules, un ejemplo de parámetros, lock y verificador, más tres fuentes Microsoft byte-exactas. No crea recursos y no sustituye la cuenta retained-original separada.

```json
{
  "compositionVersion": "1.0",
  "secretVariables": {},
  "packs": [
    {
      "path": "implementation_packs/MICROSOFT_AVM_SECURE_SFTP_INTAKE.md",
      "packId": "MICROSOFT-AVM-SECURE-SFTP-INTAKE",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": ["*"],
      "variables": {}
    }
  ]
}
```

Después de componer, ejecutar `pwsh ./azure_avm_secure_sftp_intake/verify_contract.ps1 -AllowNetwork`. El gate descarga sólo el Bicep Microsoft 0.46.1 hash-locked y compila; no despliega. Mantener bloqueado hasta cuenta/costo/red/DNS/Log Analytics/clave SSH, transferencia DEV real, SHA-256, malware, retained-original separado y reconciliación durable.
