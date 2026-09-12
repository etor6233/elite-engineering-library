# Sustainability Pillar Guidance

## 1. Metadata

```yaml
pack_id: "SUSTAINABILITY-PILLAR-GUIDANCE"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa el checklist del pilar Sustainability de AWS Well-Architected y su mapeo a controles de la biblioteca (FinOps, serverless, outbox, residencia). Información AUTHORED; no es medición de huella certificada."
stacks: ["Markdown"]
compatible_with: ["GO-FINOPS-CORE 0.1.x", "GO-DATA-RESIDENCY-CORE 0.1.x"]
incompatible_with: []
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://docs.aws.amazon.com/wellarchitected/latest/sustainability-pillar/sustainability-pillar.html"]
verified_at: "2026-09-02"
```

## 2. Applicability

Use para evaluar la sostenibilidad de un workload en nube antes de desplegar. No sustituye la medición de huella real con telemetría del proveedor.

## 3. Architecture contract

- **Ownership**: el checklist es información de mapeo; la telemetría y los factores de emisión son runtime del proveedor.
- **Invariantes**: (1) todo recurso etiquetado. (2) escalar a cero sin demanda. (3) residencia respetada. (4) teardown sin recursos huérfanos.
- **Data flow**: lectura estática (checklist + mapeo).
- **Failure modes**: no aplica (documento de información).
- **Seguridad/privacidad**: no procesa datos.
- **Performance budget**: no aplica.

## 4. Exact file manifest

```text
CREATE compliance/SUSTAINABILITY_PILLAR_CHECKLIST.md
```

## 5. Materialization blocks

### FILE: `compliance/SUSTAINABILITY_PILLAR_CHECKLIST.md`
```yaml
block_id: "SUSTAINABILITY-PILLAR-GUIDANCE:compliance/SUSTAINABILITY_PILLAR_CHECKLIST.md:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "b65f1daf621f86a57157c57355cf69117a05134f2ba77da1a8a3181120b27af5"
variables: []
secrets_allowed: false
```
````markdown
# Sustainability Pillar — Checklist (AWS Well-Architected)

Gobernado por el pilar Sustainability de AWS Well-Architected. Información `AUTHORED`: checklist de mejores prácticas y mapeo a controles de la biblioteca. No es una medición de huella certificada.

## 1. Principios de diseño

1. Entender el impacto (medir y atribuir el uso de recursos).
2. Establecer metas de sostenibilidad (KPIs por workload).
3. Maximizar la utilización (eliminar recursos ociosos).
4. Adoptar hardware/software más eficiente al madurar.
5. Usar servicios gestionados (reducir mantenimiento y energía).
6. Reducir el impacto downstream (retirar, no dejar recursos huérfanos).

## 2. Áreas de mejores prácticas → control de la biblioteca

| Área | Mejor práctica | Control de la biblioteca |
|---|---|---|
| Selección de región | Elegir regiones cerca de energía renovable y de los usuarios | Residencia de datos (V211) + `COST-FINOPS` (V190) |
| Alineación con la demanda | Right-sizing y escalar a cero | Serverless plan (`FRANCHISE_SERVERLESS_PACK_PLAN`) + `COST-FINOPS` |
| Software y arquitectura | Event-driven, retirar activos no usados, batch vs stream | `GO-RELIABLE-ASYNC-WORKERS` + transactional outbox |
| Datos | Clasificar, tiering, políticas de ciclo de vida | `GO-DATA-ANALYTICS-CORE` (V179) + `POSTGRES-BACKUP-RESTORE-CORE` |
| Hardware y servicios | Servicios gestionados, evitar sobre-aprovisionamiento | `COST-FINOPS` (V190) + serverless |
| Proceso y cultura | Medir KPIs, rastrear mejora | `SECURE-OPS` + `ENGINEERING_EXECUTION_PLAYBOOK` |

## 3. Checklist operativo

- [ ] Todo recurso lleva tags `tenant`/`environment`/`owner` (FinOps).
- [ ] No hay recursos untagged ni huérfanos (audit FinOps).
- [ ] Los workloads escalan a cero cuando no hay demanda.
- [ ] Los datos se clasifican y tiering por acceso.
- [ ] El teardown contabiliza cada recurso (no se deja nada).
- [ ] Las regiones respetan residencia de datos y cercanía a los usuarios.

## 4. Límite honesto

La medición de huella de carbono real (kgCO₂e) requiere factores de emisión por región y telemetría del proveedor (runtime). Este pack entrega el checklist y el mapeo, no el cálculo de emisiones.
````


## 6. Configuration surface

Sin variables ni secretos.

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Markdown | — | checklist sostenibilidad | LicenseRef-Workspace-Owner | none | https://docs.aws.amazon.com/wellarchitected |

## 8. Apply order

1. Componer el checklist en el proyecto bajo `compliance/`.
2. Completar el checklist operativo por workload.
3. Rollback: eliminar `compliance/SUSTAINABILITY_PILLAR_CHECKLIST.md`.

## 9. Verification

- Round-trip de hashes: 1/1 bloque reproduce byte a byte.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/SUSTAINABILITY_PILLAR_GUIDANCE_2026-09-02_V212.md`.
