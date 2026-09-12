# Tessera POSIX Evidence Log Pack Plan

Este perfil materializa, sin Git y desde Markdown, la adquisición oficial fail-closed y un runner verificable para el código POSIX de Transparency.dev Tessera. Produce 51 archivos desde dos packs. Tessera es el sucesor que Google Trillian recomienda a nuevos operadores; el source exacto reconoce como autores a Anthropic PBC, Google LLC e Internet Security Research Group. No se presenta como producto exclusivo de ninguna de esas empresas.

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
      "path": "implementation_packs/TRANSPARENCY_DEV_TESSERA_POSIX_EVIDENCE_LOG.md",
      "packId": "TRANSPARENCY-DEV-TESSERA-POSIX-EVIDENCE-LOG",
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

Este perfil no descarga, instala ni ejecuta automáticamente nada. Antes de adquirir o habilitar el log, el agente debe completar el source profile `tamper-evident-local-log`, obtener aprobación hash-linked y demostrar Linux/POSIX, Go exacto, custodia de claves, publicación independiente de checkpoints o witness, backup/restore, retención, observabilidad, carga y respuesta a incidentes. El log sólo admite receipts canónicos con hashes; prohíbe bytes documentales, valores extraídos, PII, credenciales y secretos. La evidencia de inmutabilidad regulatoria sigue perteneciendo a un storage/WORM y política aprobados: un log local aislado es tamper-evident, no deletion-proof.
