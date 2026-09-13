# Microsoft DevSkim Adapted SAST Pack Plan

Use este perfil cuando el proyecto privado Go/TypeScript necesita un baseline local de security linting sin costo y no dispone de GitHub Code Security. No lo selecciones como sustituto de análisis interprocedural, threat modeling, revisión, fuzzing, DAST o seguridad ofensiva.

```json
{
  "compositionVersion": "1.0",
  "secretVariables": {},
  "packs": [
    {
      "path": "implementation_packs/MICROSOFT_DEVSKIM_ADAPTED_SAST_GATE.md",
      "packId": "MICROSOFT-DEVSKIM-ADAPTED-SAST-GATE",
      "version": "0.1.1",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    }
  ]
}
```

Antes de ejecutar, el agente debe registrar: lenguajes/rutas reales; archivos generados, vendorizados y secretos que se excluyen; severidades que bloquean; owner y SLA de triage; suppressions justificadas y con vencimiento; baseline inicial; momento de ejecución local/CI; retención del SARIF sin código sensible; y cómo se enlazan findings con fixes/regresiones. El pack requiere red autorizada y .NET SDK 10.0.400 exacto para reconstruirse.

```text
materializar 4/4
→ verificar commit Microsoft firmado + árbol canónico
→ aplicar adaptación SharpCompress 0.48.0 declarada
→ restore público sin warnings
→ build + 300 pruebas oficiales + SCA=0
→ publicar herramienta con licencias/provenance/manifest
→ ejecutar DevSkim sobre el scope real
→ triar findings y convertir fixes en regresiones
→ ejecutar fuzzing + DAST + gates de proyecto separados
```

La salida sigue `REBUILD_VERIFIED / CONDITIONED`. Sólo el receipt del proyecto, el scope real y la aceptación de findings permiten contar el lint como ejecutado; nunca permiten afirmar que el producto completo es seguro.

V402 composition security delta: fixed Go graph floor and bounded DevSkim failure diagnostics. COMPOSITION_SECURITY_RELEASE_V402.md/json.
