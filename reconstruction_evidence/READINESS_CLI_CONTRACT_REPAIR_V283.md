# Readiness CLI contract repair — V283

Fecha: 2026-09-07. Scope: T2801, mantenimiento de biblioteca.
Procedencia: cambio y pruebas AUTHORED; no código copiado de GitHub, NASA o Google.

## Resultado

PROJECT-START-READINESS-VALIDATOR 0.6.1 corrige una contradicción verificable:
el contrato y el schema aceptaban CLI/AUTOMATION, pero otra condición exigía
WEB/MOBILE/DESKTOP/API. Un proyecto CLI válido no podía pasar sin declarar algo inexistente.

- Reproducción: dos tests fallan exclusivamente por "at least one delivery surface must be REQUIRED".
- Corrección: primeros journeys seleccionados CLI/AUTOMATION pueden constituir la interfaz.
- Exigen RUNTIME, CONTRACTS, DOMAIN-MODULES y AUTHORIZATION requeridos y enlazados.
- Se conservan todos los refs de interfaz, autorización, contratos, efectos, auditoría,
  respuesta, ayuda, capacitación, soporte, actualización, E2E, negativos y recovery.
- Las superficies web/portal/mobile/desktop/API deben enlazar su capability correspondiente;
  añadir CLI no permite ocultar una web no clasificada.
- 70 tests PASS en candidato y 70 PASS después de reconstruir desde Markdown.
  Incluyen una prueba con 14 subcasos que elimina cada frontera de evidencia.
- 8/8 archivos reconstruidos idénticos al candidato, Python 3.14.4, PowerShell 7.6.5.
- Los fixtures son sintéticos y prueban el validador, no la franquicia ni el readiness local.

Los otros negativos cubren auth retirada, journey ajeno al slice, interfaz mixta
no clasificada, versión de soporte divergente y status distinto de PROVEN.
No se cambia schema v2 ni se agregan capabilities a las 48 normativas.

## Expediente de mantenimiento reconciliado

Se crearon constitution, spec/plan/tasks, PROJECT_BLUEPRINT, PROJECT_AUTHORITY_MAP,
PROJECT_EXTERNAL_SOURCE_LOCK y PROJECT_PACK_PLAN. Son registros locales excluidos
del payload, no otros módulos de negocio. El tasks owner continúa siendo el roadmap.

Se recuperaron respuestas previas del usuario para ronda A: resultado, autoridad de
tarea, urgencia, reutilización NEW/EXISTING, procedencia, no Git obligatorio,
no backend TypeScript impuesto y no gastos externos automáticos.
ANSWERED significa petición registrada; no evidencia de cumplimiento ni aprobación nueva.
Ronda B fue generada y explicada; permanece IN_PROGRESS mientras se vinculan permisos,
versiones y pruebas de los recorridos. C–H no se rellenaron por inferencia.

Readiness real ejecutado con 0.6.1: exit 2, BLOCKED, 87 observaciones.
Antes: 102. El cambio es reconciliación documental, no porcentaje de proyecto,
no eliminación de 15 features ni prueba de calidad productiva.
Se conserva el JSON anterior más abajo; copia original adicional en staging.
El reporte local actual contiene la salida íntegra vigente y su hash se liga al checkpoint.

Observaciones anteriores que ya no aparecen:
- at least one delivery surface must be REQUIRED
- first_vertical_slice.acceptance_test_refs must be a non-empty string array
- first_vertical_slice.id is required
- first_vertical_slice.invariants must be a non-empty string array
- first_vertical_slice.journeys must be a non-empty string array
- first_vertical_slice.rollback_recovery_ref: empty evidence path
- pack_plan.path: evidence file missing: PROJECT_PACK_PLAN.md
- project.owner is required
- required_artifact: at least one non-empty specs/<feature>/spec.md is required
- required_artifact: evidence file missing: .specify/memory/constitution.md
- required_artifact: evidence file missing: PROJECT_AUTHORITY_MAP.md
- required_artifact: evidence file missing: PROJECT_BLUEPRINT.md
- required_artifact: evidence file missing: PROJECT_EXTERNAL_SOURCE_LOCK.md
- required_artifact: evidence file missing: PROJECT_PACK_PLAN.md
- rounds.A.status must be ANSWERED or PROVEN
- sources.lock_path: evidence file missing: PROJECT_EXTERNAL_SOURCE_LOCK.md

Siguen pendientes: rondas y mapa completo de capabilities, connected journeys con
evidencia real, fuentes/dependencias/freshness/monitoring, plataforma y
implementation_assurance. T2801 no se declara completo.
T2802–T2810 mantienen sus pendientes originales; no se repite V281 ni se modifica el
perfil de franquicia 67/742.

## Método y límites

El contrato normativo material es markdown_system/TOTAL_SYSTEM_CAPABILITY_CONTRACT.md:
un journey puede ser CLI o automatización, con persona/rol y todos sus eslabones.
[GitHub Spec Kit](https://github.github.com/spec-kit/) respalda trabajar con spec/plan/tasks;
[Google SRE](https://sre.google/sre-book/reliable-product-launches/) describe controles
adaptables al lanzamiento; [NASA](https://www.nasa.gov/reference/system-engineering-handbook-appendix/)
distingue verificación, validación e integración. Esas fuentes se consultaron hoy;
no publicaron este fix ni se les atribuye su implementación.

## Evidencia reproducible

Staging: elite-v283-601e98ab11404598b576727935bd63c5 bajo el temporal del sistema.
Materializar PROJECT_START_READINESS_VALIDATOR.md 0.6.1 en destino vacío y ejecutar
python -m unittest discover -s project_readiness_gate -p 'test_*.py'.
El validador real se ejecuta con --project-root y un --report nuevo;
no sustituir sus registros por los fixtures de unittest.

| Archivo | SHA-256 |
|---|---|
| red-cli.log | 7646b75aa0ee57729ca5fc6c1c15e7136bb0fc36239b58fed7d7420778895657 |
| green-cli.log | 8057ccee95310702744c8477fcd44f54d5a7d7ba4c788f3f1321ad10fc10268a |
| rebuilt-tests.log | 9f35357e2b98dd0defc6b5f1d2d7d1773934489ec501df2c5217402c1bbe104a |
| readiness-before-v283.json | 1f4bbf6bf6405cd54f434dec087bdad525e1c04f976b357376fe32e7e222fbd7 |
| readiness-after-v283.log | 53f4ca2460d3da345aa95cc10dca36246c964ff57757198a0a1c52d6d03c289c |

## Historial: reporte anterior conservado

Contenido del reporte V282, sin sobrescribir el fallo inicial:
```json
{
  "capabilities_expected": 48,
  "errors": [
    "at least one delivery surface must be REQUIRED",
    "authorities.status must be PROVEN",
    "blockers.critical_unknowns must be empty",
    "blockers.open_high_or_critical_failures must be empty",
    "capabilities.ADS-ATTRIBUTION.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.ANALYTICS-BI.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.API-BACKEND.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.ARCH-DOMAIN must be REQUIRED for a buildable vertical slice",
    "capabilities.ARCH-DOMAIN.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.AUTHORIZATION.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.BACKUP-DR.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.BROKER-STREAMING.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.CACHE.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.CD-RELEASE.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.CI must be REQUIRED for a buildable vertical slice",
    "capabilities.CI.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.CONTAINERS.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.CONTRACTS must be REQUIRED for a buildable vertical slice",
    "capabilities.CONTRACTS.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.COST-FINOPS.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.DATA-INGEST.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.DESKTOP.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.DOCS-OPS must be REQUIRED for a buildable vertical slice",
    "capabilities.DOCS-OPS.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.DOMAIN-MODULES must be REQUIRED for a buildable vertical slice",
    "capabilities.DOMAIN-MODULES.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.EMBEDDED-IOT.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.GPU-ACCEL.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.IAC-CLOUD.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.IDENTITY.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.INTEGRATIONS.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.JOBS-WORKFLOWS.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.MARKETPLACES.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.ML-AI.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.MOBILE.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.NETWORK-EDGE.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.NOTIFICATIONS.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.OBJECT-STORAGE.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.OBSERVABILITY.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.ORCHESTRATION.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.OUTBOX-INBOX.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.PAYMENTS.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.PERFORMANCE.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.PRD-INTAKE must be REQUIRED for a buildable vertical slice",
    "capabilities.PRD-INTAKE.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.PRIVACY-COMPLIANCE.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.RAG-AGENTS.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.REPO-SCM must be REQUIRED for a buildable vertical slice",
    "capabilities.REPO-SCM.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.RUNTIME must be REQUIRED for a buildable vertical slice",
    "capabilities.RUNTIME.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.SEARCH.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.SECRETS-PKI.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.SECURITY-APPSEC must be REQUIRED for a buildable vertical slice",
    "capabilities.SECURITY-APPSEC.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.SLO-INCIDENT.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.SUPPLY-CHAIN must be REQUIRED for a buildable vertical slice",
    "capabilities.SUPPLY-CHAIN.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.TEST-PLATFORM must be REQUIRED for a buildable vertical slice",
    "capabilities.TEST-PLATFORM.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.TX-DATABASE.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.WEB-PORTALS.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "capabilities.WEB-PUBLIC.classification must be REQUIRED, OPTIONAL or NONE_WITH_REASON",
    "connected_journeys must be a non-empty array",
    "dependencies.status must be PROVEN",
    "first_vertical_slice.acceptance_test_refs must be a non-empty string array",
    "first_vertical_slice.id is required",
    "first_vertical_slice.invariants must be a non-empty string array",
    "first_vertical_slice.journeys must be a non-empty string array",
    "first_vertical_slice.rollback_recovery_ref: empty evidence path",
    "operations.backup_restore_dr.status must be PROVEN",
    "operations.deployment_canary_rollback.status must be PROVEN",
    "operations.monitoring_incident_response.status must be PROVEN",
    "operations.security_privacy.status must be PROVEN",
    "operations.slo_performance_capacity_cost.status must be PROVEN",
    "pack_plan must be PROVEN with collision_check and rollback_defined true",
    "pack_plan.evidence: at least one evidence path is required",
    "pack_plan.path: evidence file missing: PROJECT_PACK_PLAN.md",
    "project.owner is required",
    "project.platform_mode must be an implemented non-BLOCK choice",
    "project.status must be READY_TO_BUILD",
    "required_artifact: at least one non-empty specs/<feature>/spec.md is required",
    "required_artifact: evidence file missing: .specify/memory/constitution.md",
    "required_artifact: evidence file missing: PROJECT_AUTHORITY_FRESHNESS_RECORD.md",
    "required_artifact: evidence file missing: PROJECT_AUTHORITY_MAP.md",
    "required_artifact: evidence file missing: PROJECT_BLUEPRINT.md",
    "required_artifact: evidence file missing: PROJECT_DEPENDENCY_UPDATE_RECORD.md",
    "required_artifact: evidence file missing: PROJECT_EXTERNAL_SOURCE_LOCK.md",
    "required_artifact: evidence file missing: PROJECT_OFFICIAL_SOURCE_PROFILE_RECORD.md",
    "required_artifact: evidence file missing: PROJECT_PACK_PLAN.md",
    "required_artifact: evidence file missing: PROJECT_VULNERABILITY_MONITORING_RECORD.md",
    "rounds.A.status must be ANSWERED or PROVEN",
    "rounds.B.status must be ANSWERED or PROVEN",
    "rounds.C.status must be ANSWERED or PROVEN",
    "rounds.D.status must be ANSWERED or PROVEN",
    "rounds.E.status must be ANSWERED or PROVEN",
    "rounds.F.status must be ANSWERED or PROVEN",
    "rounds.G.status must be ANSWERED or PROVEN",
    "rounds.H.status must be ANSWERED or PROVEN",
    "sources must be PROVEN with exact revisions and resolved licenses/notices",
    "sources.lock_path: evidence file missing: PROJECT_EXTERNAL_SOURCE_LOCK.md",
    "sources.product_provenance_modes must be non-empty"
  ],
  "record": "PROJECT_READINESS_GATE.json",
  "schema": "elite-project-readiness-report/v1",
  "status": "BLOCKED"
}
```

