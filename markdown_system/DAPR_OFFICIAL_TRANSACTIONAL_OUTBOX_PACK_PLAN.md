# Dapr Official Transactional Outbox Pack Plan

Este perfil opt-in materializa veintiún archivos: trece archivos Apache-2.0 byte-exactos del runtime y documentación Dapr firmados, más locks, approval, adquisición y gates locales. No selecciona componentes, cuentas ni producción; no promete exactly-once y excluye el SDK Go vulnerable observado.

```json
{
  "compositionVersion": "1.0",
  "secretVariables": {},
  "packs": [
    {
      "path": "implementation_packs/DAPR_OFFICIAL_TRANSACTIONAL_OUTBOX.md",
      "packId": "DAPR-OFFICIAL-TRANSACTIONAL-OUTBOX",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": ["*"],
      "variables": {}
    }
  ]
}
```

Después de componer, el agente ejecuta ambas suites offline. Sólo continúa cuando el usuario/proyecto completa la selección con store transaccional/outbox, pub/sub, acceso real, idempotencia, duplicados, recuperación y rollback demostrados. La adquisición completa usa cache aprobada o autorización de red explícita y verifica bytes, hash y commit. Un SDK Dapr sólo puede añadirse mediante nueva fuente oficial corregida y revalidada.
