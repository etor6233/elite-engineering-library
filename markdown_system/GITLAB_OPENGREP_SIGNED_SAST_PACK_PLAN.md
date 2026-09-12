# GitLab/OpenGrep Signed SAST Pack Plan

Este perfil conserva el segundo lane OpenGrep/GitLab reconstruido, pero está bloqueado para instalación: el SCA fechado del SBOM oficial Cosign 3.1.3 reabrió el verifier. No ejecutarlo ni contarlo como cobertura actual hasta fijar una revisión oficial limpia; mientras tanto usar DevSkim, SCA, secretos, fuzzing, DAST y revisión sin afirmar equivalencia.

```json
{
  "compositionVersion": "1.0",
  "secretVariables": {},
  "packs": [
    {
      "path": "implementation_packs/GITLAB_OPENGREP_SIGNED_SAST_GATE.md",
      "packId": "GITLAB-OPENGREP-SIGNED-SAST-GATE",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": ["*"],
      "variables": {}
    }
  ]
}
```

Antes de una futura reapertura, registre scope, exclusiones, severidad mínima, owner/SLA de triage, política de suppressions con vencimiento, retención de JSON/SARIF y enlace finding → fix → regression. Hoy la secuencia termina en materialización 5/5 → verifier estático → `runtime=BLOCKED`; no se autoriza red ni descarga de Cosign.

La procedencia es estrecha: OpenGrep y los tres archivos GitLab fueron bytes upstream exactos en la reconstrucción fechada; el pack sólo contiene glue `AUTHORED`. GitLab Enterprise y Commons Clause están excluidos. `REBUILD_VERIFIED / CONDITIONED` conserva evidencia, mientras `current_admission=REJECTED_VERIFIER_RUNTIME` bloquea incorporación y cualquier afirmación de seguridad.
