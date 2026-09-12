# Microsoft pg_durable Human Handoff Pack Plan

Este perfil opt-in materializa veintiséis archivos PostgreSQL License byte-exactos del release Microsoft pg_durable `v0.2.6`: licencia, NOTICE y 24 archivos del ejemplo `invoice-approval`. `live_smoke_check.sh` queda sólo en el archive hash-bound porque no termina en LF y no puede representarse VERBATIM en el contrato Markdown. No autoriza recursos Azure ni producción y no afirma que el ejemplo complete autenticación humana, correcciones, motivo o control de conflictos.

```json
{
  "compositionVersion": "1.0",
  "secretVariables": {},
  "packs": [
    {
      "path": "implementation_packs/MICROSOFT_PG_DURABLE_HUMAN_HANDOFF.md",
      "packId": "MICROSOFT-PG-DURABLE-HUMAN-HANDOFF",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": ["*"],
      "variables": {}
    }
  ]
}
```

Después de componer, ejecutar primero el smoke offline. La evaluación live requiere pg_durable `v0.2.6` completo, PostgreSQL 17/18, una cuenta Azure explícitamente autorizada y un lock de dependencias aprobado. Para integrarlo en documentos empresariales, el proyecto debe demostrar identidad OIDC/sesión, rol de revisión, evidencia por campo, correcciones, motivo, versión/conflicto, timeout/escalación, audit trail, backup/restore y rollback.
