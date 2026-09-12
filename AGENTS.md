# Elite Engineering Library — instrucciones para agentes de código

## Misión

Usa esta biblioteca para acelerar proyectos reales sin rebajar calidad, seguridad, rendimiento, operabilidad ni cumplimiento de licencias. No trates la fama de una empresa, la popularidad de un repositorio o una licencia pública como prueba suficiente de arquitectura de élite.

## Inicio obligatorio de cada tarea

1. Lee `AGENT_SYSTEM_START.md` como protocolo principal, `markdown_system/PROJECT_START_READINESS_GATE.md` como gate normativo previo y `CODEX_ELITE_PROJECT_BOOTSTRAP.md` como intake histórico compatible.
2. Clasifica la tarea como mantenimiento de biblioteca, discovery de proyecto, diseño, implementación, revisión o incidente.
3. Consulta `AI_ENGINEERING_MASTER_MAP.md` y `SYSTEMS_ENGINEERING_MASTER_MAP.md`; carga sólo los manuales autoridad materiales para la decisión.
4. Si se considera código o arquitectura pública, aplica `PUBLIC_CODE_ARCHITECTURE_ADMISSION_STANDARD.md` y consulta `PUBLIC_CODE_ARCHITECTURE_MASTER_MAP.md`.
5. Para ejecución, gates y evidencia usa `ENGINEERING_EXECUTION_PLAYBOOK.md`; materializa `implementation_packs/ENGINEERING_EXECUTION_VALIDATOR.md` 1.3.1, completa `implementation_assurance` para todo código de producto y no promociones `AUTHORED|ADAPTED|VERBATIM|MIXED` por compilación o reputación. Copia su template como `PROJECT_EXECUTION_STATE.json`, conserva `PROJECT_EXECUTION_EVENTS.jsonl` y valida el checkpoint antes de reanudar. El estado es un cursor compacto sobre los owners existentes, no un segundo plan.
6. Ante un fallo aplica `markdown_system/AGENT_ERROR_RECOVERY_PROTOCOL.md` y `markdown_system/FAILURE_LEARNING_CONTRACT.md`; abre o actualiza inmediatamente `PROJECT_FAILURE_LESSONS.md`, consulta `markdown_system/LIBRARY_FAILURE_LEARNING_LEDGER.md`, no abandones ni ocultes el gate y devuelve cualquier fix a la fuente canónica antes de continuar.
7. Si el proyecto ingiere, clasifica, extrae o almacena datos de archivos empresariales, lee y aplica obligatoriamente `markdown_system/OFFICIAL_DOCUMENT_INTELLIGENCE_PROFILE.md`, `markdown_system/OFFICIAL_DOCUMENT_SDK_ARTIFACT_LOCK.md` y `markdown_system/OFFICIAL_DOCUMENT_FIXTURE_CATALOG.md`; ninguna factura, proforma, packing list, orden, bill of lading u otro documento pasa a almacenamiento automático sin su expediente por clase/campo.
8. Si se considera `BUSINESS_CENTRAL_PLATFORM`, lee `markdown_system/MICROSOFT_BUSINESS_CENTRAL_CAPABILITY_PROFILE.md`, adquiere exactamente el snapshot fijado y bloquea implementación hasta demostrar runtime/licencia, country view, build y suites de las áreas elegidas.
9. Para admitir o actualizar dependencias, runtimes, SDKs, imágenes, providers, modelos o upstreams aplica `markdown_system/DEPENDENCY_UPDATE_CONTRACT.md` y mantiene `PROJECT_DEPENDENCY_UPDATE_RECORD.md`; ningún bot ni versión nueva se promueve sin candidate gates y rollback. Una vulnerabilidad upstream reabre la admisión, bloquea promotion, obliga a consultar primero advisory/release/patch oficiales y mantiene provenance propia para cualquier backport local.
10. Si un hecho, versión, documento o afirmación previa puede estar obsoleto o contradicho, aplica `markdown_system/AUTHORITY_FRESHNESS_AND_SELF_CORRECTION_CONTRACT.md`, conserva before/after y corrige mapas/locks/packs/evidencia sin reescribir historia.
11. Para vigilancia sin gasto incremental aplica `markdown_system/ZERO_COST_VULNERABILITY_MONITORING_PROFILE.md` y mantiene `PROJECT_VULNERABILITY_MONITORING_RECORD.md`: Dependabot incluido, OSV oficial local, cero scheduled OSV Actions y producción bloqueada hasta probar monitoring runtime.
12. Antes de adquirir un grupo de upstreams oficiales crea `PROJECT_OFFICIAL_SOURCE_PROFILE_RECORD.md` desde su template, responde todos los `required_user_inputs`, valida el perfil materializado y conserva receipts; adquirir código no cierra sus `production_blockers`.
13. Si el proyecto genera clientes Go desde OpenAPI, materializa `implementation_packs/MICROSOFT_KIOTA_OPENAPI_CLIENT_GATE.md` mediante `markdown_system/MICROSOFT_KIOTA_OPENAPI_CLIENT_PACK_PLAN.md`; fija contrato y SHA-256, genera sólo en destino ausente, conserva receipt/diff y no promueve sin compilación, SCA, auth, idempotencia, reconciliación y contract tests del proveedor.
14. Si el proyecto contiene Go y la superficie admite fuzzing, materializa `implementation_packs/GO_NATIVE_FUZZ_GATE.md` mediante `markdown_system/GO_NATIVE_FUZZ_GATE_PACK_PLAN.md`; aporta invariantes y semillas reales del dominio, fija un presupuesto finito, conserva cualquier input que falle como regresión y no presenta el PASS como sustituto de SAST, DAST, carga o seguridad del despliegue.
15. Para repositorios privados Go/TypeScript sin entitlement GitHub Code Security, materializa `implementation_packs/MICROSOFT_DEVSKIM_ADAPTED_SAST_GATE.md` mediante `markdown_system/MICROSOFT_DEVSKIM_ADAPTED_SAST_PACK_PLAN.md`; reconstruye el commit Microsoft firmado, verifica árbol canónico, adaptación SharpCompress declarada, 300 tests y SCA antes de usar el CLI. Trátalo sólo como security linting `ADAPTED / CONDITIONED`: no sustituye análisis interprocedural, fuzzing, DAST, threat model ni seguridad ofensiva.
16. Si una capability `REQUIRED` no tiene pack compatible admitido, o su fuente quedó obsoleta/rechazada, materializa `implementation_packs/CAPABILITY_GAP_RESOLUTION_GATE.md` mediante `markdown_system/CAPABILITY_GAP_RESOLUTION_PACK_PLAN.md`; crea un expediente por capability, investiga autoridades/docs/repos/releases/advisories oficiales actuales, fija evidencia y candidatos, ejecuta G0–G8 y no programa hasta `USE_REUSABLE_PACK`. `RESEARCH_INCOMPLETE` obliga a continuar; `NO_ADMISSIBLE_SOURCE`, `ACCESS_BLOCKED` o `HUMAN_DECISION_REQUIRED` deben informarse con evidencia y trigger de reapertura, nunca cubrirse con procedencia inventada.
17. Antes de publicar un candidate release local sin Git ni CI hospedado, materializa `implementation_packs/PORTABLE_SIGNED_RELEASE_EVIDENCE_GATE.md` mediante `markdown_system/PORTABLE_SIGNED_RELEASE_EVIDENCE_PACK_PLAN.md`; exige dos builds byte-idénticos, ZIP determinista, OSV 2.5.1 exacto con cero hallazgos, SPDX 2.3, in-toto/SLSA, firma Ed25519 OpenSSH y verificación independiente. La clave vive protegida fuera del proyecto y este PASS no sustituye deploy, restore, rollback, seguridad ofensiva ni aceptación del target.

No cargues todos los markdowns íntegros de forma indiscriminada. Usa los mapas como índice y conserva contexto para el dominio, el código y la verificación del proyecto.

## Regla de admisión

- `DISCOVERED`, `LICENSE_VERIFIED`, `CANDIDATE` y `CONDITIONED` no autorizan incorporación automática.
- Sólo `ELITE_REFERENCE` puede gobernar el claim estrecho auditado.
- `REUSABLE_PACK` compatible puede incorporarse directamente. Un pack `REBUILD_VERIFIED / CONDITIONED` también puede incorporarse autónomamente cuando el agente demuestre y registre que satisface todas sus condiciones en el proyecto; no necesita frenar por una aprobación ceremonial.
- Las fuentes externas oficiales se adquieren desde `markdown_system/PROJECT_INITIALIZATION_PACK_PLAN.md` y su lock de commits/SHA-256. No usar branches móviles ni exigir un repositorio Git para iniciar.
- Revalida revisión, licencia, notices, dependencias, seguridad, mantenimiento y compatibilidad antes de adoptar código externo.
- Nunca copies desde un sample, tutorial, branch principal inestable o fuente no oficial como si fuera producción.
- Registra procedencia, revisión fijada, cambios locales y obligaciones de redistribución.

## Cobertura obligatoria de sistema

Antes de declarar completa una arquitectura, evalúa explícitamente:

La lista normativa completa es `markdown_system/TOTAL_SYSTEM_CAPABILITY_CONTRACT.md`; cada una de sus 48 superficies debe clasificarse como `REQUIRED`, `OPTIONAL`, `NONE_WITH_REASON` o `BLOCKED`.

- journeys, usuarios, operadores, administradores y modelo de negocio;
- web pública, SEO, accesibilidad, performance y captación;
- portales de cliente, franquicia, administración y operaciones;
- aplicaciones móviles o de escritorio cuando el journey las requiera;
- API/BFF, dominio, workflows, consistencia e idempotencia;
- identidad, autorización por objeto/organización, auditoría y privacidad;
- datos, migraciones, índices, backup, restore y retención;
- fábricas, proveedores, inventario, logística, pagos y postventa;
- adaptadores de marketplaces, publicidad, CRM, ERP y mensajería;
- observabilidad, SLO, seguridad de supply chain e incident response;
- infraestructura, CI/CD, entornos, rollout, rollback y recuperación;
- tests unitarios, integración, contrato, E2E, carga, seguridad y resiliencia.

Si una capa no aplica, documenta por qué. No confundas “full-stack” con elegir un framework de frontend y otro de backend.

## Inicio de un proyecto nuevo

Para armar una franquicia o inyectar una capability en un proyecto existente, leer primero `markdown_system/POST_DEEPSEEK_FRANCHISE_REAUDIT_2026-09-04.md` y después usar `FRANCHISE_ACCELERATOR.md`. `markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md` selecciona actualmente 82 packs/1032 archivos y un único runtime conversacional conectado con su núcleo determinista, FinOps, resolver PostgreSQL de identidad/contacto, fence outbound durable y leads gobernados Google/Meta/TikTok/Mercado Libre Questions v4; incluye DevSkim adaptado y release firmado portable, pero no constituye certificación del proyecto productivo. OpenGrep/GitLab se conserva como evidencia fuera del perfil mientras el verifier Cosign esté rechazado por SCA. No prometer una franquicia completa ni un plazo menor a una semana sin cerrar en el target los gates live, regulatorios, de seguridad, carga, recovery, despliegue y aceptación.

Antes de expandir código:

1. si esta biblioteca no es la raíz del proyecto, exige que `INSTALL_AGENT_BRIDGE.ps1` haya instalado el bridge/Skill y que el agente se haya abierto en la raíz del proyecto cuando no hay Git; luego ejecuta `markdown_system/PROJECT_START_READINESS_GATE.md`, materializa su pack, genera y explica `PROJECT_ADVISORY_A.md`…`H.md` una ronda por vez con `render_project_advisory.py`, crea `PROJECT_READINESS_RECORD.md`, `PROJECT_FAILURE_LESSONS.md`, `PROJECT_DEPENDENCY_UPDATE_RECORD.md`, `PROJECT_AUTHORITY_FRESHNESS_RECORD.md`, `PROJECT_VULNERABILITY_MONITORING_RECORD.md` y `PROJECT_OFFICIAL_SOURCE_PROFILE_RECORD.md` desde sus templates, consulta sus autoridades y no implementes hasta `READY_TO_BUILD` con cada `advisory_prompt_ref` exacto;
2. ejecuta `VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Preflight`, registra toolchains/entornos ausentes y no confundas el PASS de biblioteca con acceso de proyecto;
3. captura el negocio y las restricciones en un `PROJECT_BLUEPRINT.md` conforme a `markdown_system/PROJECT_BLUEPRINT_CONTRACT.md`;
4. crea `PROJECT_AUTHORITY_MAP.md`, `PROJECT_EXTERNAL_SOURCE_LOCK.md` y `PROJECT_PACK_PLAN.md`;
   Si hay documentos, crea además `PROJECT_DOCUMENT_INTELLIGENCE_DECISION.md` desde `markdown_system/PROJECT_DOCUMENT_INTELLIGENCE_DECISION_TEMPLATE.md`, aplica `markdown_system/OFFICIAL_DOCUMENT_INTELLIGENCE_PROFILE.md` y mantén cada clase en `CORPUS_BLOCKED`, `ACCESS_BLOCKED`, `REVIEW_ONLY` o el estado demostrado correspondiente;
   Si elige Business Central, crea además `PROJECT_BUSINESS_CENTRAL_PLAN.md` conforme a `markdown_system/MICROSOFT_BUSINESS_CENTRAL_CAPABILITY_PROFILE.md` y no sustituye sus probes con el PASS estructural de Elite;
5. define journeys, invariantes, trust boundaries y failure costs;
6. elige la arquitectura mínima suficiente y registra ADRs materiales;
7. materializa `ENGINEERING_EXECUTION_VALIDATOR.md`, crea o valida el manifest generado en `engineering_execution_kit/`, detecta `NEW|EXISTING`, captura baseline/delta y crea el primer checkpoint de `PROJECT_EXECUTION_STATE.json`;
8. selecciona sólo packs con código y estado suficientes para el claim; un pack descriptivo no se materializa como implementación. Para toda capability requerida no cubierta, ejecuta `CAPABILITY_GAP_RESOLUTION_GATE` y devuelve la resolución a mapas, locks, catálogo, plan o ledger canónicos antes de programar;
9. compone archivos aplicando `markdown_system/COMPOSITION_PROTOCOL.md`;
10. implementa y verifica un vertical slice ejecutable de extremo a extremo;
11. captura evidencia reproducible antes de ampliar superficie.

Avanza de forma autónoma con supuestos reversibles claramente marcados, materializa verticales y ejecuta gates sin pedir permiso ceremonial. Detente sólo cuando falte una decisión humana que cambie materialmente el negocio, la regulación, el gasto, los datos sensibles, una acción destructiva o una garantía irreversible. No inventes reglas comerciales, regulatorias, financieras o de seguridad irreversibles.

## Resultado esperado

Entrega decisiones trazables y código verificable. Indica con claridad qué está listo, qué permanece condicionado, qué fue rechazado y qué evidencia falta. “Élite” significa adecuado al contexto y demostrado por gates, no máxima complejidad.
