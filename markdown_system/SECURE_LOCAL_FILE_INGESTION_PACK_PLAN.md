# Secure Local File Ingestion Pack Plan

Este perfil materializa adquisición oficial exacta y el gate local Magika/ClamAV/YARA-X. No descarga automáticamente, no incorpora reglas inventadas, no aprueba tipos ni autoriza almacenamiento de negocio. El template queda deliberadamente bloqueado hasta que el usuario complete la política del proyecto.

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
      "path": "implementation_packs/SECURE_LOCAL_FILE_INGESTION_GATE.md",
      "packId": "SECURE-LOCAL-FILE-INGESTION-GATE",
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

Antes de ejecutar: aprobar el source profile, aceptar licencias, completar todos los trust boundaries/tipos/límites/storage/incident owners, actualizar bases ClamAV, revisar/compilar reglas YARA y demostrar aislamiento. Antes de convertir o persistir: el receipt de seguridad sólo admite el siguiente gate; jamás afirma exactitud semántica ni habilita tablas de negocio.
