# Project Readiness Gate Pack Plan

Este perfil materializa únicamente el control previo al producto: catálogo explicativo official-source-routed, renderer determinista de una ronda por vez, template JSON fail-closed, CLI semántica, cierre connected-journey, renderer readiness→consumer, 102 regresiones y guía. Produce 11 archivos desde un pack `REUSABLE_PACK`. Debe componerse antes de `PROJECT_INITIALIZATION_PACK_PLAN.md`; no requiere Git, red, cuenta externa ni secrets.

```json
{
  "compositionVersion": "1.0",
  "secretVariables": {},
  "packs": [
    {
      "path": "implementation_packs/PROJECT_START_READINESS_VALIDATOR.md",
      "packId": "PROJECT-START-READINESS-VALIDATOR",
      "version": "0.7.0",
      "acknowledgeConditions": false,
      "files": ["*"],
      "variables": {}
    }
  ]
}
```

Este perfil conserva un único owner para readiness y no duplica el estado operativo. Inmediatamente después se materializa `implementation_packs/ENGINEERING_EXECUTION_VALIDATOR.md` 1.3.1, se completa `implementation_assurance`, se crea `PROJECT_EXECUTION_STATE.json` y se checkpointa `PROJECT_EXECUTION_EVENTS.jsonl`; el cursor enlaza por SHA-256 el readiness aprobado y los artifacts Spec Kit. En proyecto `EXISTING` también enlaza inventario y delta observados. Un PASS de readiness sin estado de ejecución válido no autoriza comenzar ni reanudar implementación.

El agente copia el template como `PROJECT_READINESS_GATE.json`, genera `PROJECT_ADVISORY_A.md`…`H.md` una ronda por vez, explica sus términos y registra respuestas/evidencia sin inferirlas. El validador rechaza una ronda respondida si falta, cambia, cruza o escapa su prompt exacto. Exit 2/BLOCKED es el estado inicial esperado. Sólo exit 0, reporte nuevo `READY_TO_BUILD` y evidencia consistente permiten componer el primer slice. Cada capability `REQUIRED` debe pertenecer a un journey probado y ligado a una release: persona/rol, interfaz, autorización, contrato, dominio, dato/efecto, auditoría, respuesta/error, ayuda, capacitación, soporte, impacto de actualización, E2E, negativos y recuperación. Si hay documentos, exige inventario único por clase/variante; el almacenamiento automático exige modo, mapping consumer exacto, bytes+SHA, schemas, output keys, evidencias y capabilities. Para outbox, el renderer copia sólo diez campos a un template intacto, sin habilitar approvals/proofs/endpoints/secrets. Modificar negocio, journey, contenido operativo, accesos, sources, autoridades, corpus/mapping, operación o plan reabre el gate.


V402 / modo explícito de biblioteca: el mismo pack 0.7.0 añade un template
separado y `--scope LIBRARY_INFRASTRUCTURE --level preparation|release`. El default
PROJECT conserva su contrato. En mantenimiento, copiar el template de biblioteca
a PROJECT_LIBRARY_READINESS_GATE.json y pasar --record con ese nombre.
READY_FOR_LIBRARY_WORK habilita corregir el alcance local con sus blockers y
planes visibles; no cierra T2801. READY_FOR_LIBRARY_USE exige todos los controles
locales, procedencia admitida y dos árboles exactos fuera de la raíz canónica.
No heredar ese PASS como READY_TO_BUILD de una franquicia ni producción live.
