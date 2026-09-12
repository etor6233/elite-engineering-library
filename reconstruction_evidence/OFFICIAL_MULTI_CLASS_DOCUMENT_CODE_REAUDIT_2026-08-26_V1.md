# Official Multi-Class Document Code Reaudit — 2026-08-26 V1

## 1. Pregunta y criterio

Se reabrió la admisión para responder si existía código público oficial adicional de compañías líderes que permitiera procesar de inmediato facturas, proformas, purchase orders, packing lists, bills of lading, documentos aduaneros y certificados sin inventar schemas, modelos ni garantías.

Un nombre de empresa no fue suficiente. Para entrar a la biblioteca como runtime reusable el candidato debía aportar identidad inmutable, licencia, grafo reproducible, código ejecutable, pruebas relevantes y un claim compatible con la documentación vigente del proveedor. Samples y componentes Pre-GA se conservaron como evidencia o se rechazaron; no se relabelaron como producción.

## 2. Frontera oficial vigente por clase

Las documentaciones oficiales convergen en el mismo límite:

- Microsoft Content Understanding permite analyzers prebuilt y analyzers custom definidos por schema; `prebuilt-invoice` y `prebuilt-purchaseOrder` son candidatos para factura y orden de compra, mientras las clases propias requieren analyzer/schema y evaluación del proyecto: <https://learn.microsoft.com/en-us/azure/ai-services/content-understanding/tutorial/create-custom-analyzer>.
- Google Document AI Custom Extractor permite definir un schema para documentos propios; classifier/splitter resuelve clases y límites, no crea el schema comercial del usuario: <https://docs.cloud.google.com/document-ai/docs/ce-with-genai>.
- Amazon Bedrock Data Automation custom output usa un blueprint que enumera los campos esperados y un project que se referencia al invocar la API: <https://docs.aws.amazon.com/bedrock/latest/userguide/bda-how-it-works.html>.
- Oracle Document Understanding publica modelos preentrenados para un conjunto acotado, incluida factura, y ofrece custom key-value/classification para casos propios: <https://docs.oracle.com/en-us/iaas/Content/document-understanding/using/pre_trained_doc_keyvalue_invoice.htm>.

No se halló un modelo oficial público universal y versionado que pruebe exactitud para todas las variantes de proforma, packing list, bill of lading, aduana y certificados. Incorporar un schema genérico inventado habría violado el contrato del usuario y el de los proveedores.

## 3. Código que ya estaba admitido y sirve a modelos custom

### Microsoft Azure Content Understanding

`AZURE-CONTENT-UNDERSTANDING-DOCUMENT-RUNTIME 0.1.0` materializa cinco archivos alrededor del wheel GA oficial `azure-ai-contentunderstanding==1.1.0`. El runner no fija un modelo inventado: recibe `--analyzer-id`, llama `ContentUnderstandingClient.begin_analyze_binary`, preserva `fields`, `confidence`, `source`, geometría y resultado completo, y emite hashes/receipt atómicos. Por lo tanto, el mismo código ejecuta un analyzer prebuilt o un analyzer custom oficial una vez que el proyecto aporta su ID versionado.

La admisión superior permanece fijada por `microsoft-azure-content-understanding-dotnet-1.1.0`: commit firmado `6af7db6572976f4b6cc65bb628c45456114e88a0`, archive SHA-256 `a88e5c78…`, NuGet SHA-256 `c026c61…`, firmas de autor Microsoft y repositorio NuGet verificadas, grafo de 22 paquetes y cero findings OSV en la consulta registrada. El sample oficial de factura/utility bill/sales order/purchase order conserva confidence, sources, geometry, spans y line items.

### Google Cloud Document AI

`GOOGLE-DOCUMENT-AI-RUNTIME 0.1.0` materializa cinco archivos alrededor del wheel oficial `google-cloud-documentai==3.15.0`. Exige el recurso inmutable `projects/{project}/locations/{location}/processors/{processor}/processorVersions/{version}`, llama `DocumentProcessorServiceClient.process_document`, preserva todo `ProcessResponse` y emite hashes/receipt atómicos. El processor puede ser invoice, Custom Extractor, classifier o splitter; el runner no necesita código distinto ni rellena campos ausentes.

La fuente oficial `GoogleCloudPlatform/document-ai-samples` ya estaba fijada en su revisión vigente `001ba391…`: archive de 155.461.201 bytes, SHA-256 `3131ff0604685967932cad1b689a64f7e50af1dc6b45524d99c5389342d34e8f`, Apache-2.0. Incluye flujos de clasificación/separación/extracción y un fixture público de packing list/barcode. Ese fixture prueba ingestión/layout/barcode, no exactitud de campos de packing lists reales.

### AWS y otras lanes oficiales ya fijadas

La biblioteca conserva AWS Textract Go v2 para AnalyzeExpense/AnalyzeDocument/Queries/Adapters, AWS GenAI IDP 0.6.5 para evaluación/ground truth/corrección y pipelines AWS adicionales; también OCI SDK, SAP invoice validation, IBM Docling, Microsoft MarkItDown/Table Transformer y NVIDIA NeMo condicionado. Las condiciones de cada revisión permanecen en `LIBRARY_FAILURE_LEARNING_LEDGER.md`; una lane con findings o samples incompletos no fue promovida.

## 4. Candidato AWS BDA inspeccionado y rechazado

Repositorio oficial: <https://github.com/aws-samples/sample-document-processing-with-amazon-bedrock-data-automation>.

- commit vigente y firmado: `131ea7bbac3f84104926731d3fc724d770824fcb`;
- fecha del commit: 2026-04-08;
- archive: 10.427.797 bytes, SHA-256 `eb61f18d674f96a02f0d946a296676c943a02d3b4c30324f2b245a6dd00872f`;
- licencia MIT-0: 948 bytes, SHA-256 `f085c09cb8360edeacaf3aa9d5666655d2ed7dcb1b668167a7758253625bc3ad`;
- inventario observado: 11 `.py`, 6 notebooks, 122 celdas de código, 10 JSON;
- auditoría estática repetida: 11 Python y 110 celdas compilables por AST; 12 celdas de comandos/magics quedaron fuera del claim de sintaxis Python;
- no se hallaron requirements/lock ni suite propia;
- los blueprints públicos observados son médicos/claims/cheques y un schema genérico, no schemas oficiales de proforma/packing/BOL/aduana;
- el README declara que sus wildcard IAM son demostrativos y recomienda mínimo privilegio.

Decisión: `REJECTED_COMPONENT` para el lock reusable. Es un workshop oficial útil para comprender BDA, pero no mejora el acelerador AWS integral ya fijado y añadirlo aumentaría superficie sin un grafo o tests reproducibles. El contrato oficial de blueprint/project/API sí queda como autoridad si un proyecto selecciona BDA.

Durante la auditoría se registraron `LIB-FAIL-124` a `LIB-FAIL-127`: escritura de bytecode en path Windows profundo, primera lectura sin path extendido, dependencia `jsonschema` asumida por el arnés y un helper temporal que el allowlist global rechazó hasta su eliminación exacta. Ninguno fue atribuido falsamente al código upstream.

## 5. Google Document AI Toolbox inspeccionado y no promovido

Google archivó el repositorio histórico sólo porque movió el código al monorepo `googleapis/google-cloud-python`; no se interpretó como abandono. PyPI publicó `google-cloud-documentai-toolbox==0.17.3` mediante Trusted Publishing:

- wheel: 43,7 kB, SHA-256 `46f37a42ab98a1255ea078fe4ce85084cc7ecf127b6f5ccc3c546599262da014`;
- sdist: 24,3 MB, SHA-256 `4e19e04abc9fcb4d3b131697b91c3486d03c38cbf01d06581f2bda43050f0f46`;
- publicación: 2026-08-24/25 con attestations de Google visibles en PyPI.

El paquete implementa merge de shards, acceso a text/pages/tables/forms/entities, CSV/DataFrame/BigQuery, split PDF por salida de classifier/splitter y conversiones. Sin embargo, Google lo clasifica expresamente `Experimental`, sujeto a Pre-GA Offerings Terms y con posibles cambios incompatibles: <https://docs.cloud.google.com/document-ai/docs/handle-response#document_ai_toolbox>.

Decisión: no incorporarlo al runtime estable ni a los perfiles por defecto. Puede reabrirse como `CONDITIONED` sólo si un proyecto necesita exactamente esas utilidades y acepta Pre-GA, pin, compatibilidad y rollback. La ruta default conserva el cliente GA Document AI 3.15.0 ya verificado.

## 6. Resultado verificable

Código inmediato real existente:

1. Azure y Google reciben cualquier analyzer/processor oficial versionado y guardan el resultado completo sin transformar inferencias en hechos de negocio.
2. AWS Textract recibe Queries/Adapters oficiales y conserva output/receipt; los aceleradores AWS aportan evaluación, ground truth y revisión, condicionados por sus findings registrados.
3. Los perfiles obligan a inventariar todas las clases, adquirir el proveedor por hash, aportar acceso/costo aprobado, crear schemas custom donde el propio proveedor los exige y medir cada campo contra ground truth antes de almacenar automáticamente.

Código que no existe públicamente y no se fabricó: un modelo/schema universal para cada documento comercial del mundo. La biblioteca está preparada para ejecutar el código oficial una vez definido el contrato real de datos; no está autorizada a afirmar exactitud productiva antes del corpus del proyecto.
