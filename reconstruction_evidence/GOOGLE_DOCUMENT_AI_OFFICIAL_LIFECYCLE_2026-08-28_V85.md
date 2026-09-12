# Google Document AI Official Lifecycle — V85

Fecha: 2026-08-28. Estado: `REBUILD_VERIFIED / CONDITIONED`.

## Objetivo

Cerrar con código Google público y reutilizable el lifecycle que rodea la extracción documental: crear processor, inicializar/copiar dataset y schema, importar documentos train/test/unassigned, entrenar, evaluar, desplegar, fijar default, migrar y conservar undeploy como rollback. No se escribió un motor de extracción ni se atribuyó código local a Google.

## Fuentes oficiales fijadas

- `GoogleCloudPlatform/python-docs-samples`, commit firmado/verificado `dc0eecc2187791fea70f48201241288e25fc9620`, Apache-2.0: seis source y seis tests exactos bajo `documentai/snippets` para create/train/evaluate/deploy/set-default/undeploy.
- [Google Cloud: Copy processor versions and datasets across projects](https://docs.cloud.google.com/document-ai/docs/copy-processor-versions?hl=en), `Last-Modified: Wed, 26 Aug 2026 23:31:34 GMT`: un único bloque DOM Python con dataset/schema/import/migración/deploy/default.
- [Google Developers Site Policies](https://developers.google.com/terms/site-policies): code samples Apache-2.0.
- `googleapis/google-cloud-python`, commit firmado `02d1fd863e938aa0dd8aa9fbbe4e532372837cb3`, Apache-2.0: API/samples generados actuales confirmaron `update_dataset`, `import_documents`, train/evaluate/deploy/default; sus templates genéricos no se incorporaron porque declaran requerir modificaciones.

El HTML inglés tiene SHA-256 `21c4fd07e8931e3549a085e51f1294f130dca70ea35b4007f08e4b3c5687a625`. El único code block con `def import_documents`, obtenido quitando markup y decodificando HTML sin editar lógica, tiene 18.662 bytes y SHA-256 `10897a366d7d7201f4ad4901e1b14d16472dc8f7e614b1411eb32b3c9361f232`; se clasifica `ADAPTED`, no `VERBATIM` de Git.

## Materialización

`GOOGLE-DOCUMENT-AI-OFFICIAL-LIFECYCLE` 0.1.0 materializa 17 archivos:

- 12 source/tests Google `VERBATIM`;
- 1 licencia Google Apache-2.0 `VERBATIM`;
- 1 bloque de código Cloud docs `ADAPTED` con transformación registrada;
- 3 archivos locales `AUTHORED` de guía, source lock y verificación, nunca presentados como runtime Google.

La comparación source↔materialized pasó 17/17 hashes. El layout corto `upstream/google/documentai/snippets` preserva imports y evita el límite Windows; `source-lock.json` conserva las rutas Git exactas. Un gate deliberadamente largo pasó con path máximo de 255 caracteres.

## Ejecución observada

- 13/13 Python compilan.
- `verify_official_lifecycle.py --static-only` verifica 14 upstream, hashes, bytes, compilación y diez funciones/operaciones por AST sin instalar SDK.
- El mismo verificador con `google-cloud-documentai==3.15.0` ejecuta la función Google set-default sin modificarla contra un capture client y emite `GOOGLE_DOCUMENT_AI_OFFICIAL_LIFECYCLE_OFFLINE_PASS upstream_files=14 python_files=13 functions=10`.
- Cinco tests Google offline-safe: `5 passed`.
- El sexto test Google set-default, preservado exacto, hizo una llamada live y obtuvo 403 `CONSUMER_INVALID` con proyecto ficticio; queda correctamente condicionado a GCP/ADC/IAM/cuota/costo.
- Perfil mínimo lifecycle: 3 packs / 36 archivos.
- Lane Google completa: 9 packs / 98 archivos.
- `VERIFY_LIBRARY_PASS`: 60 packs / 568 archivos / 397 Markdown antes de esta evidencia / 33 perfiles.
- `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit`: `official_google_lifecycles=1`, 111 fuentes y 12 artefactos SDK; sin red ni gasto cloud.

## Límites retenidos

El script oficial de documentación contiene placeholders, dependencias Storage/`tqdm`, una espera fija y una llamada `main` superior. El agente configura una copia del proyecto y no lo importa como módulo sin efectos. Ningún mock, sample ni fixture prueba exactitud universal. Persistencia automática permanece bloqueada hasta corpus/ground truth, evaluación por campo, revisión/reconciliación, cuenta, IAM, storage/CMEK/retención, cuotas y costo reales.

Las lecciones locales 941–959 y la condición upstream 173 preservan todas las búsquedas truncadas, paths/quoting, test live, caches y portabilidad detectados. Ningún PASS posterior borra esos fallos.
