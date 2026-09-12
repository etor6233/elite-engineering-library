# Official Document Intelligence Profile

Fecha de corte: 2026-08-28.

## 1. Objetivo y límite honesto

Este perfil obliga al agente a reutilizar código y servicios oficiales de compañías líderes antes de escribir integración propia. Su alcance es ingestión, clasificación, separación, OCR, layout, tablas, extracción, evaluación, corrección y entrega controlada de datos de documentos empresariales.

No existe en las fuentes admitidas un motor universal que garantice por sí solo `100%` de exactitud sobre documentos desconocidos. Por eso `READY_FOR_AUTOMATIC_STORAGE` no se obtiene por reputación del proveedor: se obtiene ejecutando el código oficial seleccionado contra documentos reales aprobados y comparando cada campo con ground truth. Un dato no demostrado permanece en revisión; no se inventa ni se completa por intuición del agente.

## 2. Código oficial exacto disponible

Los siguientes IDs están fijados por commit, bytes, SHA-256 y licencia dentro de `OFFICIAL_UPSTREAM_ACQUISITION_CORE.md`:

| Lane | Empresa/fuente oficial | ID exacto | Código aportado | Admisión inicial |
|---|---|---|---|---|
| conversión local | IBM, `docling-project/docling` | `docling-2.122.0` | PDF/Office/imagen a documento estructurado, OCR, layout y tablas | `PINNED_CANDIDATE`; fixtures oficiales ejecutados |
| conversión Microsoft | Microsoft, `microsoft/markitdown` | `microsoft-markitdown-0.1.7` | PDF/Word/Excel/PowerPoint/imagen/data a Markdown; adapters Azure DI y Content Understanding | `REBUILD_VERIFIED / CONDITIONED` para lane local CPython 3.12 Windows: pack V0.2/6 archivos, perfil sin formatos activos, receipt seguro obligatorio, 44 wheels/0 OSV histórico, 7 tests y conversión con wheel oficial PASS; tag unsigned, seguridad/corpus reales, graph completo/audio/URL/plugins y otros OS permanecen fuera |
| orquestación documental | Microsoft, `microsoft/durabletask-python` + Azure Samples Durable Task Skill | `microsoft-durabletask-python-1.9.0` + `azure-samples-durable-task-scheduler-5636c25` | orquestador durable determinista y skill oficial: seguridad, original retenido, extracción, evaluación, review, persistencia idempotente y evidencia final; dos samples Human Interaction exactos prueban event/timer | `REBUILD_VERIFIED / CONDITIONED`; commits firmados, wheel/sdist/source/samples/skill/licencias exactos, 1.034 tests core oficiales y 13 locales PASS, lock 6/0 OSV; los samples no autentican aprobador y la variante Azure es anónima; backend in-memory sólo test, DTS/identidad/task hub/load/recovery/live pendientes |
| integridad del historial | Transparency.dev/Tessera authors: Anthropic PBC, Google LLC e ISRG; recomendado por Google Trillian para logs nuevos | `transparency-dev-tessera-a8f33c5` | log POSIX append-only con checkpoints y pruebas Merkle para receipts canónicos sin datos | `REBUILD_VERIFIED / CONDITIONED`; commit firmado, cinco archivos Apache-2.0 VERBATIM, suite/vet/fault injection/build/fsck/E2E Linux y cuatro contratos PASS; snapshot post-release, cinco OSV no alcanzables, Linux/POSIX/keys/witness/backup/load pendientes; no reemplaza WORM |
| evaluación | IBM, `docling-project/docling-eval` | `docling-eval-1.4.2` | métricas/evaluación de pipelines documentales | `PINNED_CANDIDATE`; runtime separado y fallos upstream preservados |
| OCR/layout local | PaddlePaddle/Baidu, `PaddlePaddle/PaddleOCR` | `paddlepaddle-paddleocr-3.7.0` + `PADDLEPADDLE-PADDLEOCR-LOCAL-RUNTIME 0.1.0` | OCR PP-OCRv6, wheel/sdist y pesos det/rec exactos; runtime local de 7 archivos con lock, modelo, perfil, runner, verifier y tests | `REBUILD_VERIFIED / CONDITIONED`; commit firmado, modelos Apache-2.0 pinned, 220 upstream + 8 wrapper tests, hash-sync/import/SCA y CPU inference PASS; fixture sin ground truth, confianza no es exactitud, extras test con 7 findings y corpus/storage pendientes |
| runtime OCR/layout | PaddlePaddle/Baidu, `PaddlePaddle/PaddleX` | `paddlepaddle-paddlex-3.7.2` | pipeline ejecutable OCR, resultados y arquitectura/config PP-OCRv6 | `PINNED_CANDIDATE_CONDITIONED`; wheel/sdist/source exactos e inferencia ejercitada, pero tag commit unsigned y CI PaddleOCR sigue rama móvil `release/3.7` |
| Azure SDK/samples | Microsoft Azure | `azure-document-intelligence-code-samples` | Read, Layout, facturas, recibos, query fields, modelos custom neural/template/composed y classifier | `SAMPLE_ONLY`; requiere Azure resource |
| pipeline empresarial | Microsoft Azure | `azure-ai-document-processing-samples` | clasificación, límites de documento y extracción a schema | `SAMPLE_ONLY`; requiere Azure resource y configuración |
| accelerator | Microsoft Azure | `azure-document-intelligence-in-a-box` | ingestión PDF, custom extraction y persistencia Azure/Cosmos | `SAMPLE_ONLY`; no se adopta como producto sin threat/cost review |
| Content Understanding GA | Microsoft Azure | `azure-content-understanding-python-ga` | Bicep, clasificación/separación, analyzers por schema, entrenamiento, migración DI→CU y fixtures oficiales | `SAMPLE_ONLY / CONDITIONED`; 70 tests PASS/1 skipped sólo tras agregar dos dependencias omitidas upstream; requiere Azure resource |
| Content Understanding .NET GA | Microsoft Azure SDK | `microsoft-azure-content-understanding-dotnet-1.1.0` | SDK fuertemente tipado, LRO/retry/auth y sample `prebuilt-invoice` para invoice, utility bill, sales order y purchase order con confidence/sources/geometry/spans/line items | `PINNED_CANDIDATE`; commit firmado, NuGet autor+repositorio firmado y grafo 22/0 OSV; requiere Microsoft Foundry/modelos/costo/corpus |
| Content Understanding Python GA | Microsoft Azure SDK | `microsoft-azure-content-understanding-python-1.1.0` | SDK productivo, modelos/LRO, binary analysis y samples prebuilt con confidence/source/spans/line items | `PINNED_CANDIDATE`; tag/commit firmado, archive/licencia/notice/package/assets/CI verificados, adquisición real 52.748 archivos, sdist 114/114 byte-idéntico, 183 tests oficiales aislables y grafo 9/0 OSV; suite integral exige `.git`/test-proxy y live/corpus siguen bloqueados |
| solución configurable CU | Microsoft Azure | `azure-content-understanding-data-extraction-solution` | Functions, schemas JSON, clasificación, extracción, Cosmos DB, citas, config hashes, Key Vault, telemetría y 55 Terraform | `SAMPLE_ONLY / CONDITIONED`; 166 tests PASS en Python 3.12 tras agregar `parameterized` omitido upstream; hardening productivo requerido |
| processors | Google Cloud | `google-document-ai-samples` | OCR, forms, invoice/expense processors, custom extractor/classifier/splitter, batch y BigQuery | `SAMPLE_ONLY`; requiere GCP project/processor |
| IaC Google | Google Cloud | `google-terraform-document-ai-0.0.1` | módulo Terraform oficial para processors, CMEK, outputs, ejemplo e integration tests | `CONDITIONAL_PLATFORM`; requiere Terraform/Go/gcloud, billing, IAM y GCP project |
| workflow Google | Google Cloud | `google-document-intake-accelerator` | clasificación, extracción, validación, profile matching, autoapproval y revisión humana por microservicios | `SAMPLE_ONLY / REJECTED_FOR_IMMEDIATE_ADOPTION`; commit oficial 2022 y 163 Python AST PASS, pero Python 3.7 EOL, requirements corrupto/incompleto, UI tests/build rojos, findings critical/high, tests cloud-coupled, IaC amplio, API de upload sin auth/seguridad de contenido y autoapproval agregado no apto como storage gate |
| tipo real por contenido | Google | `google-magika-cli-1.1.0` | CLI/modelo oficial para identificar contenido independientemente de extensión; usado por Google para enrutar archivos a scanners | `PINNED_CANDIDATE`; binario Windows release exacto verificó PDF, PNG, Python y PDF disfrazado de `.txt`; no reemplaza antimalware |
| antimalware local | Cisco Talos | `cisco-clamav-1.5.4` | motor antivirus, updater de firmas, daemon/CLI, archive scanning y binarios firmados para Windows/Linux/macOS | `PINNED_CANDIDATE / CONDITIONED`; source y ZIP Windows exactos, firma GPG Cisco y update/test de bases PASS; exige firmas frescas, límites, aislamiento, HA y política GPL |
| reglas de detección | VirusTotal/Google | `virustotal-yara-x-1.20.0` | motor y CLI de reglas textuales/binarias; upstream declara uso productivo de miles de reglas sobre miles de millones de archivos | `PINNED_CANDIDATE`; binario exacto compiló y ejecutó regla oficial; complementa, nunca reemplaza, antimalware y requiere reglas versionadas/revisadas |
| archivos comprimidos seguros | Google | `google-safearchive-f7ce9d7` | reemplazos Go para tar/zip con defensas de traversal y symlink | `SAMPLE_ONLY / CONDITIONED / LINUX`; el propio README niega soporte oficial; 3 paquetes PASS en Linux/WSL2 y 15 fallos Windows retenidos |
| malware y cuarentena | Amazon Web Services | `aws-guardduty-malware-protection-s3` | CloudFormation/CDK para S3 protegido, EventBridge, SQS/DLQ, Lambda, bucket limpio y cuarentena | `SAMPLE_ONLY / CONDITIONED`; TypeScript build PASS; lock reportó 13 vulnerabilidades/7 altas y CDK synth falló por CLI/schema; servicio AWS pago requerido |
| acelerador IDP integral | Amazon Web Services | `aws-genai-idp-accelerator-0.6.5` | OCR, clasificación/separación, extracción/confianza, reglas, evaluation/ground truth, corrección, CLI/SDK, UI y despliegue serverless/headless | `PINNED_CANDIDATE / CONDITIONED`; 4.894 tests PASS más lint/IaC/codegen/build, pero locks con findings OSV/npm, typecheck integral con 2 errores y PR oficial #669 reconoce benchmark v0.6.5 no reproducible por suite ausente del commit; prohibida promoción hasta release corregida y gates repetidos |
| procesamiento multiarchivo | Microsoft | `microsoft-content-processing-accelerator-2.1.2` | Extract→Map→Evaluate→Save, upload validation, schemas, scores, Blob/Cosmos, API/UI y workflow cross-document con reglas/gaps | `SAMPLE_ONLY / CONDITIONED`; workflow 268 PASS y web build PASS; processor/API/Jest conservan fallos upstream; Azure AI/OpenAI/Blob/Cosmos/Container Apps pagos |
| procesamiento multiarchivo actual | Microsoft | `microsoft-content-processing-accelerator-main-659eaa1` | rama firmada post-v2.1.2 con Extract→Map→Evaluate→Save, Agent Framework, workflow cross-document, gap analysis, resultados editables, comentarios y cinco fixtures de reclamos | `SAMPLE_ONLY_REJECTED_FOR_IMMEDIATE_ADOPTION`; 554 Python PASS y web build PASS, pero 3/12 suites web FAIL, dependencias con 5+5+3+6 findings y correcciones sin reviewer persistido, reason/history ni ETag/version/lease; rama no publicada y servicios Azure pagos |
| pipeline agentic y prompt-flow | Amazon Web Services | `aws-intelligent-document-processing-agentic` | clasificación, extracción, validación, revisión A2I y agentes Analyzer/Matcher/Extractor/Validator/Fixer/Troubleshooter con historial e instrucciones | `SAMPLE_ONLY / CONDITIONED / LINUX`; 78 Python AST PASS; archive exacto no materializable en Windows, sin lock/test suite integral y con cuenta/servicios/costo AWS requeridos |
| parsing AWS | Amazon Web Services | `aws-textractor-1.10.0` | parseo oficial de forms, tables, geometry, queries y exports | `INTEGRATION_ONLY`; servicio Textract requerido |
| review/metrics | Amazon Web Services | `aws-textract-document-extraction-platform` | plataforma de extracción, corrección y métricas | `SAMPLE_ONLY / REJECTED_FOR_IMMEDIATE_ADOPTION`; commit exacto instala pero reescribe tres locks Python, build/test integral falla por Projen no fijado, `pnpm audit --prod` contiene 5 critical+65 high y OSV detecta 218 IDs únicos; no seleccionar hasta revisión oficial corregida y nueva admisión |
| revisión por checklist | Amazon Web Services | `aws-rapid-human-review-1.25.1` | React, backend Prisma, revisión asistida, decisión humana, workflow y CDK | `SAMPLE_ONLY / REJECTED_FOR_IMMEDIATE_ADOPTION`; frontend/backend/CDK compilan y 32+12 tests pasan parcialmente, pero locks conservan vulnerabilidades altas/crítica, lint falla y processor ejecuta pruebas cloud/costo con dependencia `@latest`; el `main` firmado reduce findings pero no los cierra y no es release |
| plataforma agentic completa | Amazon Web Services | `aws-sample-document-processing-platform` | UI de revisión, APIs, S3/DynamoDB, Step Functions, Bedrock/AgentCore, GuardDuty, VPC y seis stacks CDK | `SAMPLE_ONLY / CONDITIONED`; cinco builds PASS y workflow infra falla por grafo Node sin lock |
| evaluación multi-método | Amazon Web Services | `aws-agentic-idp-evaluation-framework` | 33 skills, catálogo de modelos, comparación de costo/velocidad/exactitud, backend/frontend y Terraform/CDK | `SAMPLE_ONLY / CONDITIONED`; build y CDK PASS, 356 tests PASS/7 FAIL/12 skipped; PII/audio bloqueados |
| evaluación campo por campo | Amazon Web Services | `aws-ocr-evaluation-workbench` | schemas JSON, truth por documento, comparación Textract/Bedrock/BDA, detección de truncamiento, métricas, costos e historial de runs | `SAMPLE_ONLY / CONDITIONED`; 667 PASS y fixtures exactos en Python 3.12; una prueba symlink bloqueada por privilegio Windows, graph/lint no frozen |
| evaluación estructurada | AWS Labs | `awslabs-stickler-eval-0.6.0` | `ExactComparator`, comparación recursiva, Hungarian uno-a-uno, confusion matrix, confidence calibration y validación de instancia JSON Schema | `PINNED_CANDIDATE / CONDITIONED`; commit firmado, wheel/sdist atestados, lock runtime 26/0 OSV y 1.535 PASS/6 FAIL Windows/2 skip; un extra dentro de `List[StructuredModel]` no suma FA/FP, por lo que exige schema cerrado previo y cero FA/FD/FN/FP; `all_fields_matched` nunca autoriza storage |
| queries/adapters | Amazon Web Services | `aws-textract-queries-example` | Queries, clasificación/orquestación y ejemplo IDP | `SAMPLE_ONLY`; adapter entrenado con corpus propio cuando aplique |
| SDK/runtime Textract | Amazon Web Services | `aws-sdk-go-v2-textract-1.45.0` + `GO-AWS-TEXTRACT-DOCUMENT-RUNTIME 0.2.0` | módulo productivo `service/textract v1.45.0`; perfil cerrado, receipt seguro obligatorio y AnalyzeExpense/AnalyzeDocument/Queries/Adapters sólo por clase aprobada | `REBUILD_VERIFIED / CONDITIONED`; 147/147 source, suite oficial, 7+3 locales, offline/vet/build PASS; tag/commit unsigned y cuenta/corpus/adapter real pendientes |
| extracción GPU | NVIDIA | `nvidia-nemo-retriever-main-2026-08-26` | contenido y metadata de documentos, texto, tablas, gráficos e infografías | `CONDITIONAL_PLATFORM`; commit firmado no publicado; Linux/GPU/NVIDIA services y dos advisories del perfil mínimo impiden promoción |
| tablas | Microsoft | `microsoft-table-transformer-1.0.0` | detección/estructura de tablas y evaluación GriTS | `PINNED_CANDIDATE`; componente especializado |
| validación de facturas | SAP | `sap-btp-dox-invoice-validation` | DOX + CAP + UI + roles, validator assignment, corrección y decisión | `SAMPLE_ONLY_REJECTED_FOR_IMMEDIATE_ADOPTION`; sin locks/CI/tests, build API no reproducible, 5 UI + 17 router + 160 Python advisories, sin CAS/lease/binding invoice-project, CSRF/TLS debilitados y snapshot compartido sin auth; productos SAP y acuerdos separados requeridos |
| SDK OCI | Oracle | `oracle-oci-python-sdk-2.185.0` | cliente oficial `oci.ai_document`, analyze_document y modelos/pretrained/custom | `INTEGRATION_ONLY`; commit firmado, wheel oficial/import Python 3.14 y lock de auditoría 23/0 advisories probados; upstream no trae sample ejecutable ni lock runtime cerrado y cuenta OCI/corpus siguen requeridos |
| factura a ERP | Oracle | `oracle-ai-invoice-handling` | proyecto OIC importable con Document Understanding, Object Storage e ingreso ERP | `SAMPLE_ONLY_REJECTED_FOR_IMMEDIATE_ADOPTION`; el workflow reviewer→approver está sólo en narrativa y no en el export, que conserva OCIDs, no tiene tests y exige Oracle Integration/ERP/OCI |

`SAMPLE_ONLY` significa código real oficial reutilizable como punto de integración y prueba, no arquitectura productiva certificada. El agente conserva esta clasificación en todos los documentos derivados.

Los tres motores locales ya no quedan sólo como referencias independientes: `SECURE_LOCAL_FILE_INGESTION_GATE.md` materializa seis archivos de coordinación `AUTHORED` y el perfil `SECURE_LOCAL_FILE_INGESTION_PACK_PLAN.md` los compone con el adquiridor oficial. Once tests y un recorrido real verifican hash/version de binarios, base ClamAV fresca, tipo Magika, ruleset YARA-X compilado y receipt atómico. El template nace `AWAITING_USER`, no contiene reglas inventadas y mantiene `business_storage_authorized=false`.

`MICROSOFT_DURABLE_DOCUMENT_ORCHESTRATION.md` cierra el handoff entre esos gates: materializa 13 archivos y sólo acepta referencias opacas/hashes en history; review y persistencia quedan ligadas al resultado evaluado, y todo efecto ambiguo exige reconciliación en lugar de retry ciego. `DURABLE_DOCUMENT_PIPELINE_PACK_PLAN.md` compone el control plane 5/51; las lanes Azure 10/85, Google 9/102 y AWS 7/75 ya lo incluyen. Azure añade byte-verbatim los samples Python oficiales Microsoft `prebuilt-invoice`, análisis de bytes locales `prebuilt-documentSearch` y same-resource analyzer copy, con sus tests/contratos; documentSearch conserva contenido/layout pero no es schema de negocio y el copy muta recursos y nunca se autoejecuta. Google añade byte-verbatim el sample oficial Custom Document Extractor con schema variable por solicitud y el lifecycle processor/dataset/import/train/evaluate/deploy/default/undeploy/batch GCS/response handling. AWS añade el wheel/runtime oficial Amazon Textractor 1.10.0 y sus pruebas Queries exactas. Ninguno convierte un fixture o sample en ground truth. Esto no convierte el backend in-memory en producción ni sustituye acceso/corpus/ground truth/identidad humana.
Para recepción SFTP sobre Azure existe una lane opt-in separada: `MICROSOFT_AVM_SECURE_SFTP_INTAKE_PACK_PLAN.md` compone 1/8 con Microsoft AVM Storage 0.33.0, Bicep 0.46.1 hash-locked, SSH-only, permisos `cw`, Private Endpoint y red/Shared Key públicos deshabilitados. No reemplaza el retained-original: Microsoft documenta que versioning no funciona con HNS y WORM no funciona mientras SFTP está habilitado; por eso bytes+SHA-256 deben copiarse y reconciliarse en el storage retenido elegido antes de extracción.

Cuando el blueprint requiera detectar manipulación del historial de receipts, `TRANSPARENCY_DEV_TESSERA_POSIX_EVIDENCE_LOG.md` y `TESSERA_POSIX_EVIDENCE_LOG_PACK_PLAN.md` aportan la lane oficial condicionada. Se despliega aparte sobre Linux/POSIX y recibe sólo `document_sha256`, `evidence_sha256` y referencias opacas aprobadas después del flujo durable. No recibe documentos ni campos. Requiere clave bajo custodia, publicación independiente de checkpoints o witness, backup/restore, control de acceso, retención y operación. Un log local es tamper-evident; AWS S3 Object Lock, GCS Object Retention o Azure Blob WORM siguen siendo las opciones de retención inmutable verificable cuando el proyecto las exige.

Los wheels productivos exactos de Azure Content Understanding `1.1.0`, Azure Document Intelligence `1.0.2`, Google Document AI `3.15.0`, Google Cloud Storage `3.13.1` y Oracle OCI `2.185.0`, sus hashes, imports verificados y condiciones transitivas están en el artifact core y sus evidencias. `OFFICIAL_DOCUMENT_SDK_ARTIFACT_CORE.md` los convierte en adquisición ejecutable: lock vigente de catorce artefactos, approval enlazado al SHA, cache/network, verificación de bytes/hash y receipt. Azure Document Intelligence, Azure Content Understanding, Google Document AI, Google Cloud Storage, OpenAI y Oracle OCI incluyen wheel+source oficial; el sdist Azure DI aporta 56 samples —cuatro de factura desde URL/bytes— y 20 tests. OCI expone cliente AI Document y modelos invoice/KV/confidence/geometry, pero no se presenta como schema universal. El wheel+sdist Google Document AI Toolbox `0.17.3` se mantiene sólo como `PINNED_CANDIDATE`: Alpha/experimental y bloqueado por PyArrow CVE-2026-25087 sin versión corregida compatible. Para Azure CU, el sdist PyPI y el source Microsoft firmado se enlazaron por `114/114` archivos byte-idénticos. Los wrappers REST y soluciones completas permanecen condicionados.

`MICROSOFT_AZURE_DOCUMENT_INTELLIGENCE_OFFICIAL_INVOICE_SAMPLE.md` convierte uno de esos samples en una unidad inmediatamente reconstruible sin modificarlo: fija el archivo Microsoft sync y la licencia MIT byte-a-byte, adquiere el `sample_invoice.jpg` oficial del mismo commit por approval/tamaño/SHA-256 y ejecuta offline la función upstream para demostrar `prebuilt-invoice`, bytes, locale y endpoint. `AZURE_DOCUMENT_INTELLIGENCE_OFFICIAL_INVOICE_SAMPLE_PACK_PLAN.md` lo compone con el artifact core en 14 archivos. Esto resuelve la ausencia de fixture para ese camino exacto; no hereda esa evidencia a los otros 55 samples ni autoriza precisión o almacenamiento productivo.

`GOOGLE_CLOUD_DOCUMENT_AI_OFFICIAL_PROCESS_SAMPLE.md` materializa sin modificar los samples y tests live `process_document` y Custom Document Extractor de `GoogleCloudPlatform/python-docs-samples` en commits firmados. El sample actual `handle_response_sample.py` construye `DocumentSchema` y lo entrega por `ProcessOptions.schema_override`, permitiendo campos distintos por solicitud; su contrato offline pasa con el SDK oficial 3.15.0. El pack adquiere además invoice, packing-list y JSON procesado oficiales. Su perfil process compone 2/19. `GOOGLE_DOCUMENT_AI_OFFICIAL_LIFECYCLE.md` añade 16 source/tests Google exactos y el bloque oficial de dataset/migración para create/import/train/evaluate/deploy/default/undeploy, batch de un objeto o prefijo GCS, estado por documento, JSON shards y response handling OCR/form/table/entity/split/layout/custom; su perfil lifecycle compone 3/40 y la lane Google completa 9/102. El recorrido packing-list verificó siete entidades y barcode `E-58702865`; las confidencias publicadas se preservan como evidencia, no como autorización automática.

Ese lock también fija `openai==3.3.1` como proveedor agentic opcional para Responses con file input, structured output, file search/code interpreter, MCP y function tools. Es `INTEGRATION_ONLY`: no reemplaza un processor documental con evidencia por campo ni autoriza almacenamiento automático.

Los fixtures públicos oficiales para invoices, receipts, una packing list/barcode, Queries, forms y tablas están fijados por path/hash en `markdown_system/OFFICIAL_DOCUMENT_FIXTURE_CATALOG.md`; sólo autorizan smoke tests.

### Frescura oficial verificada hasta el 2026-08-27

La API oficial de GitHub confirmó que los commits fijados de `Azure-Samples/azure-ai-content-understanding-python`, `Azure-Samples/data-extraction-using-azure-content-understanding`, `GoogleCloudPlatform/document-ai-samples`, `GoogleCloudPlatform/document-intake-accelerator` y `aws-samples/aws-textract-document-data-extraction-platform` todavía son exactamente la cabeza de sus ramas oficiales y que los repositorios permanecen activos con la licencia esperada. Para `amazon-textract-textractor`, `v1.10.0` sigue siendo el tag más reciente aunque `master` tenga commits posteriores. Para AWS GenAI IDP, `v0.6.5` sigue siendo el tag más reciente aunque `main` haya avanzado; la PR oficial abierta #669 no pertenece al release y reconoce que una suite benchmark faltaba del source committed. AWS RAPID mantiene `v1.25.1` como release estable y un `main` firmado posterior; ambos fueron probados por separado y siguen rechazados. AWS Assess Workbench no publica releases: el `main` exacto `79a57b5…` fue auditado y rechazado por fallos de suite y dependencias vulnerables. IBM Document Extraction Toolkit tampoco publica releases; su commit unsigned `d9144ec…` fue rechazado porque el propio contrato es prototipo con roles placeholder, sus locks son vulnerables y test/build fallan. Los samples Microsoft Expense Submission MCP y Approvals Box del commit firmado `8c2cb6e…` conservan patrones oficiales de MCP Apps, Entra OBO/Graph y UI de aprobación, pero fueron rechazados para adopción inmediata: no fijan su grafo npm ni publican tests; Expense no demuestra persistencia durable, aislamiento ni idempotencia, y Approvals Box no compila, usa actor privilegiado de fallback y separa decisión/auditoría sin transacción. Una rama o PR posterior nunca reemplaza automáticamente un release fijado ni hereda sus pruebas.

Evidencia reproducible: `reconstruction_evidence/OFFICIAL_DOCUMENT_UPSTREAM_FRESHNESS_2026-08-26_V1.md`. Esta comprobación valida identidad y actualidad del source, no corrige los blockers de dependencias, corpus, cuenta, costo o producción ya registrados. Las reconstrucciones profundas de los samples Textract, Google Intake, Azure Data Extraction, AWS RAPID, AWS Assess Workbench, IBM Document Extraction Toolkit y Microsoft Expense/Approvals confirman precisamente esa separación: siguen siendo código oficial, pero están rechazados para adopción inmediata en `reconstruction_evidence/AWS_TEXTRACT_EXTRACTION_PLATFORM_REBUILD_2026-08-26_V1.md`, `reconstruction_evidence/GOOGLE_DOCUMENT_INTAKE_ACCELERATOR_REBUILD_2026-08-26_V1.md`, `reconstruction_evidence/AZURE_CONTENT_UNDERSTANDING_DATA_EXTRACTION_REBUILD_2026-08-27_V1.md`, `reconstruction_evidence/AWS_RAPID_HUMAN_REVIEW_REAUDIT_2026-08-27_V1.md`, `reconstruction_evidence/AWS_ASSESS_WORKBENCH_REAUDIT_2026-08-27_V1.md`, `reconstruction_evidence/IBM_DOCUMENT_EXTRACTION_TOOLKIT_REAUDIT_2026-08-27_V1.md`, `reconstruction_evidence/MICROSOFT_EXPENSE_SUBMISSION_MCP_REAUDIT_2026-08-27_V1.md` y `reconstruction_evidence/MICROSOFT_APPROVALS_BOX_REAUDIT_2026-08-27_V1.md`. La última auditoría Azure preserva el SDK GA 1.1.0 separado, no adopta el cliente REST preview ni su frontera anónima/persistencia sin gate.

## 3. Matriz por clase documental

| Clase | Ruta oficial primaria | Ruta alternativa | Estado antes de datos/acceso |
|---|---|---|---|
| factura comercial o de proveedor | Azure prebuilt invoice; Google Invoice Parser; AWS AnalyzeExpense | Docling + extractor custom | `PREBUILT_CANDIDATE` |
| recibo/ticket | Azure prebuilt receipt; Google Expense Parser; AWS AnalyzeExpense | Docling | `PREBUILT_CANDIDATE` |
| proforma invoice | Azure Content Understanding/custom neural; Google Custom Extractor; AWS Queries/Bedrock schema | Docling/Table Transformer para estructura | `CUSTOM_MODEL_REQUIRED` |
| packing list | Azure Content Understanding/custom model; Google Custom Extractor; AWS Queries/Bedrock schema | Docling/Table Transformer | `CUSTOM_MODEL_REQUIRED` |
| purchase order | Azure Content Understanding `prebuilt-procurement`/`prebuilt-purchaseOrder`; Google Custom Extractor; AWS Queries/Bedrock schema | Docling/Table Transformer | `PREBUILT_CANDIDATE`; evaluar por país/proveedor |
| bill of lading / conocimiento de embarque | Google Custom Extractor; Azure Content Understanding/custom model; AWS Queries/Bedrock schema | Docling/NVIDIA para estructura | `CUSTOM_MODEL_REQUIRED` |
| commercial invoice de importación | prebuilt invoice sólo como baseline; luego custom schema | multi-provider evaluation | `CUSTOM_MODEL_REQUIRED` |
| delivery note / remito | Azure/Google/AWS custom extraction | Docling/Table Transformer | `CUSTOM_MODEL_REQUIRED` |
| certificado de origen, calidad o inspección | custom classifier + custom extractor | Docling/NVIDIA | `CUSTOM_MODEL_REQUIRED` |
| despacho, aduana, seguro, carta de porte | custom classifier/splitter/extractor | Docling/NVIDIA | `CUSTOM_MODEL_REQUIRED` |
| contrato, lista de precios, ficha técnica | custom classifier/extractor; layout/table engines | Docling/NVIDIA/Table Transformer | `CUSTOM_MODEL_REQUIRED` |
| paquete PDF con varios documentos | Azure Content Understanding classifier/splitter; Google splitter/classifier; pipelines AWS | Docling como conversor por unidad | `CLASSIFIER_AND_SPLITTER_REQUIRED` |
| manuscrito, escaneo degradado o formato desconocido | comparación OCR/layout multi-provider | revisión obligatoria | `EVALUATION_REQUIRED` |

El nombre de una clase no habilita almacenamiento. El proyecto debe fijar schema y reglas por país, proveedor, idioma y versión documental.

## 4. Adquisición sin Git

Después de materializar `PROJECT_INITIALIZATION_PACK_PLAN.md`, el agente valida primero el perfil local-first de seis fuentes. La adquisición real exige el approval hash-linked completo:

Para obtener conjuntamente el adquiridor y el gate ejecutable se compone `markdown_system/SECURE_LOCAL_FILE_INGESTION_PACK_PLAN.md`. Esto crea código/config/tests, pero no descarga ni aprueba nada; las respuestas del usuario y el approval siguen siendo bloqueantes.

```powershell
pwsh -NoProfile -File elite_sources/apply_source_profile.ps1 `
  -ProfilePath elite_sources/source-profiles/secure-file-ingestion.json `
  -LockPath elite_sources/upstream-source-lock.json `
  -Destination elite_upstreams `
  -ValidateOnly
```

Una vez contestados todos sus `required_user_inputs` y reconocidos los blockers, ejecuta la adquisición real con el approval exacto:

```powershell
pwsh -NoProfile -File elite_sources/apply_source_profile.ps1 `
  -ProfilePath elite_sources/source-profiles/secure-file-ingestion.json `
  -LockPath elite_sources/upstream-source-lock.json `
  -Destination elite_upstreams `
  -ApprovalPath PROJECT_SECURE_FILE_SOURCE_PROFILE_APPROVAL.json
```

Para el perfil documental base, el agente ejecuta:

```powershell
pwsh -NoProfile -File elite_sources/apply_source_profile.ps1 `
  -ProfilePath elite_sources/source-profiles/document-intelligence.json `
  -LockPath elite_sources/upstream-source-lock.json `
  -Destination elite_upstreams `
  -ApprovalPath PROJECT_DOCUMENT_SOURCE_PROFILE_APPROVAL.json
```

Cuando el objetivo es un pipeline cloud completo, el agente no adquiere el perfil integral de 46 fuentes por defecto. Debe elegir exactamente uno de estos perfiles después de completar todas sus respuestas y generar el approval hash-linked correspondiente:

```powershell
# AWS: 9 fuentes oficiales para intake, GuardDuty, Textract/IDP, evaluación y revisión.
pwsh -NoProfile -File elite_sources/apply_source_profile.ps1 `
  -ProfilePath elite_sources/source-profiles/aws-secure-document-pipeline.json `
  -LockPath elite_sources/upstream-source-lock.json `
  -Destination elite_upstreams `
  -ValidateOnly

# Google Cloud: 6 fuentes para Document Intake/AI/Terraform/Magika + ClamAV/YARA-X.
pwsh -NoProfile -File elite_sources/apply_source_profile.ps1 `
  -ProfilePath elite_sources/source-profiles/google-secure-document-pipeline.json `
  -LockPath elite_sources/upstream-source-lock.json `
  -Destination elite_upstreams `
  -ValidateOnly

# Microsoft/Azure: 13 fuentes para Content Processing/Understanding/DI, SDK .NET/Python GA, privacidad y seguridad local.
pwsh -NoProfile -File elite_sources/apply_source_profile.ps1 `
  -ProfilePath elite_sources/source-profiles/microsoft-secure-document-pipeline.json `
  -LockPath elite_sources/upstream-source-lock.json `
  -Destination elite_upstreams `
  -ValidateOnly
```

`-ValidateOnly` no descarga ni aprueba. Para adquirir, el agente repite sólo el perfil elegido con su `-ApprovalPath`; una aprobación de AWS no es válida para Google o Microsoft porque cada archivo y su hash son diferentes.

El runner verifica archive, tamaño, SHA-256, licencia y lockfile cuando existe, y genera receipt. El código descargado no se atribuye a Elite ni se mezcla con código local sin registrar `VERBATIM`, `DEPENDENCY_PIN` o `ADAPTED` por path.

Los artefactos runtime se adquieren con el segundo runner materializado por `PROJECT_INITIALIZATION_PACK_PLAN.md`, después de completar `official_document_sdks/artifact-approval.template.json` con el SHA exacto y sólo los IDs elegidos:

```powershell
pwsh -NoProfile -File official_document_sdks/acquire_document_sdk_artifacts.ps1 `
  -ArtifactId azure-ai-contentunderstanding-1.1.0-wheel `
  -Destination elite_artifacts/document-sdk `
  -ApprovalPath PROJECT_DOCUMENT_SDK_ARTIFACT_APPROVAL.json
```

Descargar todos los wheels por disponibilidad no está permitido: el proyecto selecciona un lane, genera su lock transitivo por Python/OS y recién entonces instala/proba.

Cuando el target aprobado sea Linux y el proyecto permita archivos comprimidos, la defensa de extracción se adquiere por separado para que Windows nunca reciba una falsa admisión:

```powershell
pwsh -NoProfile -File elite_sources/apply_source_profile.ps1 `
  -ProfilePath elite_sources/source-profiles/secure-archive-linux.json `
  -LockPath elite_sources/upstream-source-lock.json `
  -Destination elite_upstreams `
  -ApprovalPath PROJECT_SECURE_ARCHIVE_SOURCE_PROFILE_APPROVAL.json
```

La solución agentic AWS antigua permanece fuera del perfil documental base de 46 fuentes porque su archive exacto no es portable a Windows; ya está declarada dentro de `aws-secure-document-pipeline`, cuyo target debe ser Linux o macOS. AWS Labs Stickler sí integra el perfil base y AWS, pero sólo como evaluador condicionado: `STRICT-DOCUMENT-FIELD-EVALUATION-GATE 0.1.1` valida ambos payloads con JSON Schema cerrado, sanitiza rechazos, exige cero FA/FD/FN/FP y mantiene storage false; el perfil `strict-field-evaluation` adquiere sólo esa fuente. El sample AWS BDA purchase-order firmado integra el perfil como corpus/schema/notebook condicionado, nunca como storage automático ni cobertura de otras clases. PaddleOCR/PaddleX integran el perfil únicamente como lane local condicionada: sus dos revisiones de pesos se verifican por SHA y nunca por alias móvil, y ningún score habilita storage. El acelerador AWS 0.6.5 queda fail-closed por sus findings OSV/npm y typecheck. Google safearchive tampoco integra el perfil base: queda disponible como source Linux explícito y la extracción de archives permanece prohibida hasta seleccionarlo o admitir otra implementación. El runner rechaza plataformas incompatibles antes de red; no renombra rutas upstream ni descarga un subconjunto silencioso.

## 5. Intake obligatorio antes de elegir motor

El agente no empieza implementación documental hasta recibir y registrar:

- inventario completo de clases y variantes; país, idioma, emisor y versión;
- ejemplos reales autorizados, incluidos casos buenos, malos, multipágina, duplicados y ambiguos;
- schema esperado por clase: campo, tipo, obligatoriedad, unidad, moneda, zona horaria y cardinalidad;
- relaciones y controles: totales, subtotales, impuestos, moneda, cantidades, pesos, bultos, SKUs, PO/PI/PL/BOL relacionados;
- tratamiento de PII, datos financieros, retención, residencia, cifrado y quién puede revisar;
- cuentas sandbox, región, identidad/scopes y presupuesto de los proveedores candidatos;
- ground truth aprobado por responsable humano, separado del conjunto de evaluación;
- política explícita por campo: autoaceptar, validar, revisar o rechazar;
- volumen, tamaño, latencia, concurrencia, disponibilidad y costo objetivo;
- destino transaccional y estrategia de idempotencia, versionado, corrección y auditoría.

Si falta uno de estos elementos, el estado es `CORPUS_BLOCKED`, `ACCESS_BLOCKED` o `POLICY_BLOCKED`; el agente sigue pudiendo adquirir y probar código público, pero no declara la automatización lista.

## 6. Contrato de evaluación derivado de herramientas oficiales

El expediente se crea desde `markdown_system/PROJECT_DOCUMENT_INTELLIGENCE_DECISION_TEMPLATE.md`. Para cada combinación `{clase, variante, campo, proveedor/model_version}` se conserva:

```yaml
document_class: ""
variant: ""
field: ""
ground_truth_owner: ""
ground_truth_count: 0
exact_match_count: 0
normalized_match_count: 0
missing_count: 0
false_value_count: 0
review_count: 0
model_provider: ""
model_id: ""
model_version: ""
schema_version: ""
evaluated_at: ""
evidence_path: ""
decision: BLOCK|REVIEW_ONLY|LIMITED_AUTOMATION|READY_FOR_AUTOMATIC_STORAGE
```

Las métricas y herramientas se toman del proveedor seleccionado —por ejemplo evaluación de custom extractors/adapters, response confidences, geometry y Docling Eval—. Elite sólo normaliza el expediente; no inventa un modelo ni sustituye los tests oficiales.

`READY_FOR_AUTOMATIC_STORAGE` exige, como mínimo:

1. todos los campos críticos presentes en ground truth representativo;
2. cero `false_value_count` permitido para campos cuya política exige exactitud total;
3. controles aritméticos y referenciales del documento aprobados;
4. documentos fuera de distribución enviados a `REVIEW_ONLY`, nunca aceptados silenciosamente;
5. proveedor/model/schema fijados y re-evaluación obligatoria ante cambios;
6. evidencia de reintento, idempotencia, duplicados, páginas faltantes y documento corrupto;
7. reviewer y procedimiento de corrección con historial inmutable;
8. aprobación del owner de negocio, seguridad y datos.

Estos gates no dicen que el OCR sea infalible; impiden que un error no demostrado entre al sistema como dato definitivo.

## 7. Pipeline permitido

```text
archivo original inmutable + hash
→ malware/type/size gate
→ clasificación y separación oficiales
→ OCR/layout/tablas oficiales
→ extracción prebuilt o custom oficial
→ normalización declarativa versionada
→ validaciones aritméticas/referenciales
→ decisión por campo y documento
→ autoaceptación demostrada o bandeja de revisión
→ escritura idempotente + provenance + auditoría
→ evaluación continua y rollback de modelo/schema
```

El agente puede escribir el glue mínimo que conecte estas APIs sólo cuando el usuario haya elegido `CUSTOM_PLATFORM`; debe marcarlo `AUTHORED` y jamás llamarlo código de Microsoft, Google, AWS, IBM o NVIDIA.

## 8. Gates de ejecución

Antes de empezar el sistema:

- `UPSTREAM_LOCK_VALID sources=84 selected=84` para el lock global vigente;
- adquisición exacta del lane seleccionado y receipt íntegro;
- compilación/tests oficiales locales posibles;
- probe sandbox redacted del servicio elegido;
- corpus/ground truth y split de evaluación presentes;
- benchmark comparativo de al menos dos rutas oficiales cuando el campo sea crítico o el prebuilt no cubra la clase;
- threat/privacy/cost review y owner de revisión;
- `PROJECT_DOCUMENT_INTELLIGENCE_DECISION.md` con estado por clase y campo;
- rollback probado para modelo/schema y conservación del original.

Hasta que todos los `REQUIRED` estén en `READY_FOR_AUTOMATIC_STORAGE` o `REVIEW_ONLY` aceptado, el gate global permanece `NOT_READY_UNDER_EXPANDED_USER_STANDARD`.

## 9. Uso en proyecto nuevo o existente

En un proyecto nuevo, este perfil corre antes del primer vertical documental. En un proyecto existente, el agente descubre parsers, modelos, schemas, tablas, colas y evidencia ya presentes; ejecuta el mismo benchmark y reemplaza sólo el delta no demostrado. No obliga a reescribir una integración que ya produzca evidencia equivalente o superior.

## 10. Fuentes oficiales de comportamiento

- Microsoft Azure Document Intelligence: quickstart y modelos prebuilt/custom de `learn.microsoft.com`.
- Google Cloud Document AI: processor list, Custom Extractor, clasificación/splitting y evaluación de `cloud.google.com`.
- AWS Textract: AnalyzeExpense, Queries y Custom Adapters de `docs.aws.amazon.com`.
- OpenAI Responses: file inputs, structured outputs y tools de `developers.openai.com`.
- IBM Docling/Docling Eval, NVIDIA NeMo Retriever y Microsoft Table Transformer: repositorios oficiales exactos fijados en el source lock.

La documentación online gobierna contrato y capacidades; el archive inmutable gobierna el código exacto auditado.

## 11. Entrada documental por email

Para AWS, `AWS_SECURE_EMAIL_ATTACHMENT_PROCESSING_PACK_PLAN.md` compone 7 packs/71 archivos: acquisition lock, receptor SES 1/9, core MIME 1/5, release GuardDuty 1/9, dispatch Magika→AWS IDP 1/11, handoff postprocessing inmutable 1/8 y seguridad local 1/6. Receiver fija Boto3/Botocore 1.43.83, retiene raw y cuarentena por versión; GuardDuty sólo libera la versión clean con tag coincidente; Magika 1.1.0 identifica y copia esa versión al input IDP configurado; el helper `idp_common` oficial difiere HITL pendiente o congela salida terminal KMS/Object Lock y emite FIFO/receipt sin autoridad de persistencia. Round-trip, 11+18+16+14+14 tests, builds hash-locked y cfn-lint PASS; LIVE exige costo/cuenta/canaries/aislamiento/config/corpus/HITL/SNS/DLQ-redrive/consumer/restore/rollback.

No es un código AWS `VERBATIM`: es integración `AUTHORED` gobernada por contratos oficiales, porque los tres samples email AWS auditados permanecen rechazados. En `OFFICIAL_PLATFORM` estricto la clase queda `NO_SOURCE`; en una plataforma que admita glue auditado puede usarse sólo después de cuenta/SES/MX/KMS/IAM/costo/canary/DLQ-replay/seguridad target. En ambos casos, no entrega storage authority: cada adjunto sigue `PENDING` y pasa routing, provider extraction, evaluación por campo y review antes de persistencia empresarial.
