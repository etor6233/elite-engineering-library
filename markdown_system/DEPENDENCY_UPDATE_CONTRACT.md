# Dependency Update Contract

## 1. Propósito

Una actualización de librería, runtime, SDK, imagen, action, provider, modelo, dataset o código upstream es un cambio de producto y supply chain. El agente no la aplica por novedad, por un bot ni porque “latest” exista. Debe demostrar que el candidato conserva o mejora contratos, seguridad, licencias, rendimiento, operación y rollback.

Este contrato es `AUTHORED` por Elite Engineering Library. Puede operar sin repositorio Git: trabaja sobre copias aisladas, manifests, locks, hashes y evidencia. Git/PRs son un transporte opcional, no una condición de corrección.

## 2. Artefacto obligatorio

Todo proyecto crea y mantiene:

```text
PROJECT_DEPENDENCY_UPDATE_RECORD.md
```

desde `PROJECT_DEPENDENCY_UPDATE_RECORD_TEMPLATE.md`. La fila existe desde la admisión inicial, aunque todavía no haya actualización pendiente.

## 3. Inventario mínimo

Registrar por dependencia directa, toolchain y componente distribuido:

```yaml
component_id: "ecosystem:name"
kind: library|runtime|sdk|cli|container|action|provider|module|model|dataset|upstream-source
scope: build|development|test|runtime|operation
current_version: ""
current_digest_or_commit: ""
source_registry_or_repository: ""
license_expression: ""
notice_paths: []
lock_or_manifest_paths: []
artifact_consumers: []
support_eol: ""
security_advisory_sources: []
owner: ""
update_policy: patch-auto-candidate|minor-auto-candidate|manual-candidate|frozen-with-reason
rollback_artifact: ""
```

“Auto-candidate” sólo autoriza crear y evaluar un candidato aislado. Nunca autoriza promoción automática.

## 4. Triggers

Abrir evaluación cuando ocurra:

- advisory/CVE/KEV o compromiso de maintainer/registry/signing key;
- release estable aplicable, EOL o fin de soporte;
- yank, deprecación, licencia/terms/notices modificados;
- incompatibilidad de runtime/OS/arquitectura/provider;
- bug reproducible, performance regression o riesgo operacional;
- revisión periódica definida por criticidad;
- cambio del contrato oficial de una integración;
- modificación de un source lock, modelo, dataset o toolchain.

Un bot, feed o scanner descubre candidatos; no decide su admisión.

## 5. Respuesta a vulnerabilidad upstream ya incorporada

Una vulnerabilidad nueva reabre automáticamente la admisión del componente aunque antes estuviera aprobado. El agente no modifica código a ciegas ni espera pasivamente: contiene, demuestra qué artefacto está afectado, consulta autoridad oficial y prepara una corrección verificable.

### 5.1 Secuencia obligatoria

```text
alerta/advisory/fallo explotable
→ registrar incidente y preservar evidencia
→ congelar promotion y limitar exposición de forma reversible
→ identificar versión+commit/digest+SBOM+consumers exactos
→ demostrar reachability/exploitability o conservar UNKNOWN
→ consultar advisory/release/changelog/patch/maintainer oficiales
→ seleccionar una response lane
→ construir candidato aislado
→ ejecutar gates y canary
→ promover exactamente el digest probado o mantener BLOCKED
→ actualizar locks, SBOM, notices, runbooks y lecciones
```

La ausencia de evidencia de reachability no equivale a `NOT_REACHABLE`. Si el componente, configuración, call path o artefacto no pueden demostrarse, el estado es `UNKNOWN` y la promotion permanece bloqueada según severidad y exposición.

### 5.2 Fuentes y datos que debe contrastar el agente

- identidad exacta del artefacto observado y todos sus consumers;
- advisory oficial del proyecto/proveedor y, cuando aplique, CVE, GHSA, OSV y KEV;
- rangos afectados/corregidos, release/tag, changelog, migration guide y maintainer statement;
- commit de corrección y diff oficial, firma, provenance, registry metadata y hashes;
- madurez del exploit, exposición runtime, datos/tenants/secretos alcanzables y controles compensatorios;
- EOL, compromiso de maintainer/registry/signing key y alternativas ya admitidas.

Un scanner es señal de entrada, no autoridad suficiente para corregir ni para ignorar.

### 5.3 Response lanes permitidas

- `OFFICIAL_FIXED_RELEASE`: existe release oficial corregida. Fijar su versión y digest, reconstruir locks/SBOM/notices, ejecutar todos los gates afectados, canary y promover el mismo digest probado.
- `OFFICIAL_PATCH_UNRELEASED`: existe commit/diff oficial, pero no release. Se puede contener y esperar; fijar el commit exacto si licencia, build y provenance lo permiten; o hacer un backport mínimo. Todo backport se etiqueta `ADAPTED_PATCH` o `AUTHORED_PATCH`, enlaza el commit oficial y conserva diff/tests propios: nunca se presenta como release oficial.
- `NO_UPSTREAM_FIX`: no existe corrección oficial. Deshabilitar/aislar la ruta afectada, reemplazar por una alternativa admitida o bloquear producción. Una aceptación de riesgo exige owner, controles y expiración; el agente no la autoaprueba.
- `COMPROMISED_UPSTREAM_OR_SUPPLY_CHAIN`: poner artefactos/registry en cuarentena, detener promotion, rotar secretos potencialmente expuestos, reconstruir desde inputs conocidos y verificar provenance. Activar respuesta a incidentes y sustituir el componente si la confianza no se restablece.
- `FALSE_POSITIVE_OR_NOT_REACHABLE`: sólo con evidencia reproducible del artefacto, configuración y call path. La excepción tiene owner, alcance, expiración y próxima revisión; un ignore del scanner no elimina el hallazgo.

Un WAF, flag, rollback o aislamiento puede ser contención inmediata, pero no se registra como corrección definitiva salvo que elimine de forma demostrable la capacidad vulnerable.

### 5.4 Autonomía y límites

El agente puede investigar fuentes públicas oficiales, inspeccionar código/locks/SBOM, bloquear promotion, aplicar contención reversible en entornos no productivos, preparar candidatos y ejecutar pruebas sin esperar instrucciones adicionales. Debe pedir autoridad cuando la acción implique interrupción o mutación productiva, aceptación de riesgo, credenciales, gasto, cambio legal/comercial, datos reales o una operación destructiva. Ante explotación crítica activa puede fallar cerrado y bloquear un despliegue; debe registrar por qué.

### 5.5 Cierre verificable

No se cierra el incidente hasta demostrar que el componente fue corregido, eliminado o contenido con condición explícita; escanear todos los consumers; promover el digest exacto probado; revocar o retener para forensics el artefacto anterior; actualizar `PROJECT_DEPENDENCY_UPDATE_RECORD.md`, `PROJECT_FAILURE_LESSONS.md`, `PROJECT_AUTHORITY_FRESHNESS_RECORD.md`, SBOM, locks, notices, ADR/runbook e incidente productivo; y fijar próxima vigilancia.

## 6. Intake del candidato

Antes de editar manifests:

1. consultar únicamente registry, release, advisory, changelog, migration guide y repositorio oficiales;
2. fijar versión/tag, commit/digest, tamaño y hashes publicados u observados;
3. verificar firma/provenance cuando exista y registrar su trust root;
4. comparar licencia, notices, terms, ownership, mantenimiento, EOL y path overrides;
5. leer breaking changes, deprecations, defaults, flags, schemas, protocolos, data migrations y security notes;
6. identificar dependencias transitivas nuevas/eliminadas, scripts de instalación y privilegios;
7. definir consumers, blast radius, incompatibilidades y rollback antes de construir;
8. enlazar fallos/condiciones anteriores desde `PROJECT_FAILURE_LESSONS.md` y `LIBRARY_FAILURE_LEARNING_LEDGER.md`.

Si una fuente móvil es lo único disponible, se conserva `BLOCKED` o se fija un commit explícito con la condición de desarrollo; nunca se usa `latest`, `main`, `master` o un rango flotante como identidad productiva.

## 7. Actualización aislada

```text
snapshot baseline verificable
→ workspace candidato vacío/aislado
→ actualizar un cluster causal mínimo
→ resolver lock de forma determinista
→ inspeccionar diff de manifest/lock/generated code
→ regenerar SBOM, license report y notices
→ build/test candidate
→ comparar baseline/candidate
→ decidir promote|reject|blocked
```

Cluster causal significa runtime+SDKs inseparables o generator+runtimes generados que deban versionarse juntos. No mezclar actualizaciones independientes en un único resultado verde.

Nunca ejecutar install scripts no auditados con secretos, acceso productivo o permisos administrativos. Registry cache y artefactos descargados se consideran input hostil hasta verificar identidad e integridad.

## 8. Gates obligatorios

Seleccionar según blast radius, pero no omitir sin `NONE_WITH_REASON`:

- frozen/locked install desde entorno limpio y source permitido;
- compile, lint/type, unit/property/negative tests;
- contract/integration contra sandbox o fixture oficial;
- schema/data migration forward, mixed-version y rollback/forward-recovery;
- authn/authz, secrets, crypto, serialization y input-hostile regressions;
- SCA/advisories, malware, license/notices, SBOM y provenance;
- API/ABI/protocol/config compatibility y deprecation budget;
- browser/mobile/OS/architecture matrix cuando aplica;
- performance p50/p95/p99, memoria, tamaño, startup, throughput y costo;
- backup/restore, deploy/canary/rollback y observabilidad;
- clean rebuild del artefacto exacto que será promovido.

Una vulnerabilidad “corregida” no habilita promotion si el update introduce una regresión material. Un advisory sin exploitability/reachability se prioriza con evidencia, no se ignora.

## 9. Resultado

Estados:

- `DISCOVERED`: existe candidato, no evaluado;
- `ASSESSING`: expediente en curso;
- `BLOCKED`: falta toolchain, acceso, información oficial o decisión;
- `REJECTED`: candidato no satisface gates; conservar causa y alternativa;
- `CANDIDATE_GREEN`: gates seleccionados pasan, aún no desplegado;
- `CANARY`: desplegado con límites y rollback armado;
- `PROMOTED`: mismo digest/lock aprobado se promovió y observa;
- `ROLLED_BACK`: candidato retirado; datos/config/protocol recovery registrado;
- `SUPERSEDED`: una evaluación posterior reemplaza esta decisión sin borrar historia.

`PROMOTED` exige manifest/lock, SBOM/notices, evidencia y rollback actualizados. La versión anterior y sus símbolos/migrations se conservan durante la ventana definida.

## 10. Fallos y aprendizaje

Cada error, warning material, skip, incompatibilidad o regresión abre/actualiza `PROJECT_FAILURE_LESSONS.md` conforme a `FAILURE_LEARNING_CONTRACT.md`. Un candidato rechazado es conocimiento valioso: registrar fingerprint, versión exacta y gate que falló para que el agente no lo repita.

## 11. Cadencia por riesgo

La cadencia la fija el proyecto, no este contrato. Baseline recomendado para intake, sujeto a jurisdicción y operación:

- advisories críticos/explotados: monitor continuo o diario y triage inmediato;
- runtime/auth/crypto/database/cloud provider: revisión semanal y antes de cada release;
- runtime libraries: semanal/quincenal según exposición;
- dev/test tooling: mensual o por release;
- source locks/modelos/datasets: sólo con expediente de cambio y re-evaluación;
- revisión integral de EOL/licencias/alternativas: trimestral.

No confundir revisión con obligación de actualizar. A veces la decisión segura es fijar, aislar, compensar o reemplazar.

## 12. Continuación de proyectos existentes

El agente descubre manifests/locks/images/toolchains existentes, reconstruye el baseline y crea el registro delta. No actualiza todo al abrir el proyecto. Primero separa urgente, EOL, incompatible, mantenimiento y opcional; después procesa un cluster causal por vez.
