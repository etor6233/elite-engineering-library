# Implementation Pack Contract

Un pack que pretenda generar código debe seguir este contrato. Los encabezados son obligatorios incluso cuando una sección declare “no aplica”.

## 1. Metadata

```yaml
pack_id: ""
pack_version: ""
status:
  authority: NOTE|SUPPORTED_REFERENCE|ELITE_REFERENCE
  implementation: SPEC_ONLY|SNIPPET|RECONSTRUCTIBLE|REBUILD_VERIFIED
  admission: DISCOVERED|LICENSE_VERIFIED|CANDIDATE|CONDITIONED|REUSABLE_PACK
claim: ""
stacks: []
compatible_with: []
incompatible_with: []
license_expression: ""
upstream_sources: []
verified_at: "YYYY-MM-DD"
```

## 2. Applicability

- problema que resuelve;
- señales de adopción;
- señales de rechazo;
- requisitos asumidos;
- límites del claim.

## 3. Architecture contract

- ownership y boundaries;
- interfaces y data flow;
- invariantes;
- failure/degraded modes;
- seguridad y privacidad;
- performance budget;
- operación, migración y rollback.

## 4. Exact file manifest

```text
CREATE path/to/file.ext
PATCH  existing/file.ext anchor="exact stable anchor"
DELETE path/to/file.ext reason="only when explicitly authorized"
```

No usar “crear los archivos necesarios”. Cada archivo debe enumerarse.

## 5. Materialization blocks

Formato obligatorio por archivo:

````markdown
### FILE: `path/to/file.ext`

```yaml
block_id: "PACK-ID:file:v1"
operation: CREATE|PATCH
provenance: AUTHORED|ADAPTED|VERBATIM
source: "local" # o URL + commit + path + líneas auditadas
license: "Project license" # o SPDX exacto del upstream
variables: [PROJECT_NAME]
secrets_allowed: false
```

```language
contenido completo o patch inequívoco
```
````

`VERBATIM` exige que el texto pueda redistribuirse bajo la licencia declarada y preserve notices. `ADAPTED` registra qué cambió. `AUTHORED` indica composición propia gobernada por referencias, no copia de upstream.

## 6. Configuration surface

Tabla de variables con tipo, default seguro, validación, secreto/no secreto, mutabilidad y efecto. Toda combinación inválida debe fallar antes de servir tráfico.

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|

No usar `latest`, branches móviles ni rangos abiertos para reconstrucción. El proyecto puede actualizar después mediante expediente separado.

## 8. Apply order

Pasos deterministas para un workspace vacío y para uno existente. Declarar precondiciones, conflictos, migrations y rollback.

## 9. Verification

- syntax/type/build;
- unit/property/integration/contract/E2E según riesgo;
- negative authorization y abuse;
- migration/version skew;
- performance workload;
- observabilidad y fallos inducidos;
- backup/restore cuando sea source of truth;
- license/SBOM/notices.

Cada comando declara éxito esperado y artifacts de evidencia.

## 10. Reconstruction evidence

Para `REBUILD_VERIFIED`:

- identificador de entorno limpio;
- SHA-256 del pack;
- toolchain fijado;
- lista y SHA-256 de archivos materializados;
- comandos y resultados;
- smoke journeys;
- divergencias;
- fecha y revisor/agente.

Si falta una pieza, el estado máximo es `RECONSTRUCTIBLE` o inferior.
