# Microsoft Business Central Platform Pack Plan

Este perfil materializa exclusivamente el lock, configuración, adquisición y ejecución condicionada del entorno local Microsoft Business Central. No acepta la EULA, no instala Docker, no eleva permisos, no crea el contenedor y no autoriza producción por sí mismo.

```json
{
  "compositionVersion": "1.0",
  "secretVariables": {},
  "packs": [
    {
      "path": "implementation_packs/MICROSOFT_BUSINESS_CENTRAL_EXECUTABLE_PLATFORM.md",
      "packId": "MICROSOFT-BUSINESS-CENTRAL-EXECUTABLE-PLATFORM",
      "version": "0.1.1",
      "acknowledgeConditions": true,
      "files": ["*"],
      "variables": {}
    }
  ]
}
```

Después de componer, el usuario completa `business_central/project-bc-platform-config.template.json`, conserva los blockers exactos y ejecuta primero `test_business_central_runtime.ps1` y el preflight. `invoke_business_central_dev.ps1 -Execute` requiere approval separado ligado a los hashes exactos; nunca se invoca durante composición.
