# All Implementation Packs Materialization — V78

Fecha: 2026-08-28  
Estado: `VERIFY_LIBRARY_PASS` + `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit`

## Resultado gobernante

- 57 implementation packs y 536 archivos materializables desde Markdown;
- 387 archivos Markdown contando esta evidencia final;
- 32 perfiles de composición; Azure document compone 8 packs/73 archivos, Google document 7/68 y AWS document 7/75;
- source lock `OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.63`: 110 fuentes oficiales exactas y 14 perfiles fail-closed;
- ledger: 878 fallos locales + 160 condiciones upstream = 1.038 IDs únicos, cero duplicados y cero lecciones locales abiertas.

## Código oficial Microsoft incorporado

`MICROSOFT_AZURE_CONTENT_UNDERSTANDING_OFFICIAL_INVOICE_SAMPLE 0.1.0` materializa cinco archivos. Dos son `VERBATIM` Microsoft MIT:

- `sample_analyze_invoice.py`: 10.539 bytes, SHA-256 `0cb9d7b0e183cd1c3f677c79a9d8d5dd9a1ed30bf397f84b753eedd7819c53fd`;
- `LICENSE.txt`: 1.074 bytes, SHA-256 `7c77a44a8acd9b41fdc209864a8016b3d430b5d0e09309818d5b7444336df744`.

El sample proviene de `Azure/azure-sdk-for-python@129a5cbb06fba44e11ffc92f6bcdca14a2f5b6eb`, path `sdk/contentunderstanding/azure-ai-contentunderstanding/samples/sample_analyze_invoice.py`, y es byte-idéntico al contenido del sdist oficial `azure-ai-contentunderstanding 1.1.0` de 230.330 bytes/SHA-256 `00f39c7cf2ba50c8d4586315f4295a45fad0881daaff07d13d33aea2bafc7ec8`. El wheel oficial permanece fijado en 101.987 bytes/SHA-256 `d1d6bdeffe02f5c8cc5309f0e120a1338f51af04638605f9617b7b501e2d1dd6`.

El código llama exactamente `prebuilt-invoice` y proyecta los campos demostrativos, confidence, source, spans, line items y usage documentados por Microsoft para invoice, utility bill, sales order y purchase order. Los otros tres archivos del pack son packaging/lock/test `AUTHORED` y no se atribuyen a Microsoft.

## Verificación ejecutable

- pack reconstruido desde Markdown: 5/5 archivos y hashes PASS;
- contratos offline: 3/3 PASS;
- sintaxis del sample y del contrato: dos bytecodes explícitos fuera del árbol source PASS;
- composición Azure limpia después del bump del source lock: 8 packs/73 archivos PASS;
- `VERIFY_LIBRARY_PASS`: 57 packs/536 archivos y los 32 perfiles con sus conteos exactos;
- `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit`: `official_invoice_samples=1`, `official_azure_cu_samples=1`, `official_google_process_samples=1`, diez artifacts SDK, 110 fuentes y los demás componentes focales en uno.

No se realizó una llamada live a Azure, no se usaron credenciales y no se atribuyó precisión sobre documentos del proyecto. El upstream demo usa una URL móvil bajo `main`; el perfil empresarial continúa exigiendo archivo admitido local, identidad/costo/región/modelos, corpus/ground truth y evaluación estricta antes de persistir.

## Reauditoría Google classify/split

El archive oficial `GoogleCloudPlatform/document-ai-samples@001ba391ab4a2f40d001cc0387618cb3c3699523` permanece fijado en 155.461.201 bytes/SHA-256 `3131ff0604685967932cad1b689a64f7e50af1dc6b45524d99c5389342d34e8f`. El source lock ahora conserva nueve artifacts exactos: el notebook HITL ya auditado, cinco archivos del workflow classify/split y tres del splitter histórico.

El repositorio completo conserva `SAMPLE_ONLY` porque contiene samples estrechos útiles, pero tres componentes quedan rechazados para adopción inmediata:

- `hitl-custom-review`: usa un servicio Document AI HITL deprecado y no cierra identidad/resultado/conflicto/evidencia;
- `classify-split-extract-workflow`: helper manual en lugar de suite automatizada, excepciones absorbidas, fallback de baja confianza a `generic_form`, IAM admin, base `python:3.10-buster`, CMD a archivo inexistente y trabajo futuro explícito en HITL/cuotas/>15 páginas/MIME/logging/splitting;
- `pdf-splitter-python`: Google lo marca deprecated y dirige a Document AI Toolbox.

Document AI Toolbox 0.17.3 sigue condicionado: sus 164/164 tests oficiales pasaron en Linux, pero el metadata exige `pyarrow<23` y la corrección de CVE-2026-25087 está en PyArrow 23.0.1. No se forzó una combinación que Google no publica como compatible.

## Recuperación de fallos

Las lecciones 862–878 conservaron y corrigieron: globs Windows inválidos, entornos de test incompletos, paths supuestos, URIs GitHub mal compuestas, binding de arrays PowerShell, uso incorrecto de `LASTEXITCODE`, diferencia LF de licencia, dos conteos narrativos obsoletos, consulta de regresión que confundía snapshots históricos, cobertura faltante del auditor y escritura inestable de bytecode dentro del source. Cada reparación tiene regresión focal y los gates globales finales son verdes.

## Límites vigentes

Esto no completa ingestión universal ni extracción perfecta. El código oficial agregado resuelve una integración estrecha de `prebuilt-invoice`; no aporta schemas ni exactitud demostrada para proforma invoice, packing list, bill of lading, aduana u otras clases. Esas superficies continúan `CUSTOM_MODEL_REQUIRED` o `CORPUS_BLOCKED` hasta encontrar/admitir source oficial suficiente y demostrarlo con documentos autorizados, ground truth por campo, seguridad, revisión, persistencia idempotente, operación, recuperación y rollback del proyecto.
