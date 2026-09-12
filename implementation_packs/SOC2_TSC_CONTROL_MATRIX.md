# SOC 2 TSC Control Matrix

## 1. Metadata

```yaml
pack_id: "SOC2-TSC-CONTROL-MATRIX"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa el mapeo de controles SOC 2 (AICPA TSC 2017) a los controles existentes de la biblioteca, con plantilla de evidencia por criterio. Información AUTHORED; no es la certificación."
stacks: ["Markdown"]
compatible_with: ["GO-HUMAN-APPROVAL-CORE 0.1.x", "POSTGRES-BACKUP-RESTORE-CORE 0.1.x", "OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.1.x"]
incompatible_with: []
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://www.aicpa-cima.com/topic/audit-assurance/audit-and-assurance-greater-than-soc-2"]
verified_at: "2026-09-02"
```

## 2. Applicability

Use como punto de partida del mapeo de evidencia SOC 2 Type II. No sustituye la opinión del auditor ni el periodo de observación.

## 3. Architecture contract

- **Ownership**: la matriz es información de mapeo; la evidencia y la opinión son del operador/auditor.
- **Invariantes**: (1) Security (CC) aplica siempre. (2) A/PI/C/P aplican sólo si se comprometen. (3) cada criterio comprometido tiene owner, artefacto, frecuencia y retención.
- **Data flow**: lectura estática (matriz + plantilla).
- **Failure modes**: no aplica (documento de información).
- **Seguridad/privacidad**: no procesa datos.
- **Performance budget**: no aplica.

## 4. Exact file manifest

```text
CREATE compliance/SOC2_TSC_CONTROL_MATRIX.md
```

## 5. Materialization blocks

### FILE: `compliance/SOC2_TSC_CONTROL_MATRIX.md`
```yaml
block_id: "SOC2-TSC-CONTROL-MATRIX:compliance/SOC2_TSC_CONTROL_MATRIX.md:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "d9b119a36a3d3e541fdeb299aae5a352bce75409b482f72b71bfa16d3dc6a992"
variables: []
secrets_allowed: false
```
````markdown
# SOC 2 Type II — Control Matrix (TSC 2017)

Gobernado por AICPA Trust Services Criteria (TSC 2017). Documento de información `AUTHORED`: mapea las categorías de criterios a controles ya existentes en la biblioteca. No sustituye la opinión de un auditor; es el punto de partida del mapeo de evidencia.

## 1. Categorías TSC 2017

| Categoría | Aplica | Nota |
|---|---|---|
| Seguridad (Common Criteria, CC) | Siempre | Base obligatoria de toda revisión SOC 2 |
| Disponibilidad (A) | Opcional | Sólo si el compromiso cubre availability |
| Integridad de procesamiento (PI) | Opcional | Sólo si se compromete PI |
| Confidencialidad (C) | Opcional | Sólo si se compromete C |
| Privacidad (P) | Opcional | Sólo si se procesan datos personales |

## 2. Mapeo de criterios comunes a controles de la biblioteca

| Serie TSC (categoría) | Control de la biblioteca (pack/área) |
|---|---|
| CC1 — Entorno de control | ADRs + `AGENTS.md` + estándar de admisión `PUBLIC_CODE_ARCHITECTURE_ADMISSION_STANDARD` |
| CC2 — Comunicación e información | `SECURE-OPS` runbooks + `ENGINEERING_EXECUTION_PLAYBOOK` |
| CC3 — Evaluación de riesgos | `ZERO_COST_VULNERABILITY_MONITORING` + `FAILURE_LEARNING_CONTRACT` + `LIBRARY_FAILURE_LEARNING_LEDGER` |
| CC4 — Monitoreo | Auditoría append-only + FinOps audit (V190) + observabilidad |
| CC5 — Actividades de control | `GO-HUMAN-APPROVAL-CORE` (doble control, dinero nunca auto-aprobado) |
| CC6 — Acceso lógico/físico | Identidad + permisos por objeto/organización + OIDC + aislamiento tenant |
| CC7 — Operaciones del sistema | `POSTGRES-BACKUP-RESTORE-CORE` + idempotencia + transactional outbox |
| CC8 — Gestión de cambios | `DEPENDENCY_UPDATE_CONTRACT` + pin exacto de versiones + migraciones up/down |
| CC9 — Mitigación de riesgos | `OFFICIAL_UPSTREAM_ACQUISITION_CORE` + `DEPENDENCY_LICENSE_EVIDENCE_CORE` + pin de SHA-256 |

## 3. Plantilla de evidencia (por criterio)

Para cada criterio comprometido, registrar:

1. **Control**: qué control de la biblioteca lo cubre.
2. **Owner**: rol responsable.
3. **Artefacto**: evidencia reproducible (log, migración, test, pin).
4. **Frecuencia**: con qué periodicidad se genera.
5. **Retención**: dónde y por cuánto se conserva.

## 4. Límite honesto

El mapeo de evidencia y la opinión Type II (periodo de observación + pruebas del auditor independiente) son proceso externo, no código. Este pack entrega el mapeo, no la certificación.
````


## 6. Configuration surface

Sin variables ni secretos.

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Markdown | — | matriz de controles | LicenseRef-Workspace-Owner | none | https://www.aicpa-cima.com |

## 8. Apply order

1. Componer la matriz en el proyecto bajo `compliance/`.
2. Completar la plantilla de evidencia por criterio comprometido.
3. Rollback: eliminar `compliance/SOC2_TSC_CONTROL_MATRIX.md`.

## 9. Verification

- Round-trip de hashes: 1/1 bloque reproduce byte a byte.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/SOC2_TSC_CONTROL_MATRIX_2026-09-02_V210.md`.
