# AWS Agentic Intelligent Document Processing — 2026-08-26 V1

## Resultado

Se fijó y auditó el repositorio oficial `aws-samples/aws-ai-intelligent-document-processing`. Aporta código real de AWS para las cinco fases declaradas por el propio proyecto: clasificación, extracción, enriquecimiento, validación y revisión humana. Incluye además una implementación agentic con Analyzer, Matcher, Extractor, Validator, Instructions Fixer, Troubleshooter y Save Instructions. Se admite como `SAMPLE_ONLY_CONDITIONED_LINUX`: no se presenta como motor universal, servicio gratuito ni precisión demostrada sobre documentos del usuario.

## Identidad inmutable

- repositorio oficial: `https://github.com/aws-samples/aws-ai-intelligent-document-processing`;
- commit: `31d671ef54c6694fea6431ee4102e3de85b982c1`, fechado 2026-02-06;
- archive ZIP: 127.378.048 bytes;
- SHA-256 archive: `23326e62e8cd2cdd87854070907628615ac8597af9cedbe4dfcc236a4491fd6c`;
- licencia MIT-0, 927 bytes;
- SHA-256 licencia: `ef47b4ae2a1a8d38ef87b38dc5927967e7b472ffa3f7d125e89dbc8353a1e7af`;
- 537 entradas ZIP; 78 archivos Python; 19 `requirements*.txt`; cero lockfiles integrales observados.

El source lock y el runner de adquisición conservan commit, bytes, hashes, licencia y plataformas soportadas. El código sigue perteneciendo a AWS; Elite sólo fija su identidad, clasificación, condiciones y evidencia.

## Código oficial observado

En `guidance/agentic-orchestration` existen implementaciones para:

- grafo multiagente con Analyzer, Matcher, Extractor y Validator;
- refinamiento de instrucciones mediante Instructions Fixer;
- diagnóstico de fallos mediante Troubleshooter;
- guardado/versionado de instrucciones;
- schemas de extracción y reglas de validación de negocio;
- historial de jobs, resultados y artefactos en servicios AWS;
- ejemplo de órdenes de compra y comparación con datos de referencia;
- UI y acciones de ejecución/inspección del procesamiento.

En `guidance/prompt-flow-orchestration` existen implementaciones para:

- evento S3, Textract asíncrono, SQS/SNS y tracking DynamoDB;
- clasificación y separación por páginas;
- extracción específica por clase documental;
- validación JSON Schema y controles contra información de referencia;
- resultado explícito `needs_manual_review`;
- integración Amazon A2I para revisión humana cuando falla schema, contenido, clase o sistema;
- persistencia de resultados y estados de validación.

La descripción anterior proviene del código y README exactos del commit. No implica que las reglas de facturas, proformas, packing lists, órdenes, embarques o aduana del proyecto ya estén definidas: esas reglas y sus documentos reales deben aportarse y probarse en el proyecto.

## Verificación local

- se extrajo un subconjunto portable de 406 archivos sin modificar sus bytes;
- los 78 Python parsearon con Python 3.12: 78 PASS, 0 fallos AST;
- no se observó una suite unitaria ejecutable integral ni un lock completo, por lo que no se inventó un PASS de tests/build;
- no se desplegó AWS ni se realizaron llamadas pagas;
- no se usaron credenciales, cuentas ni documentos del usuario.

## Defecto de portabilidad retenido

El ZIP oficial contiene 15 entradas cuyos segmentos incluyen literalmente `s3:`. Windows no puede materializar `:` como parte de un nombre de archivo. Tres son archivos reales:

| Archivo final | Bytes | SHA-256 |
|---|---:|---|
| `extraction_prompt.md` | 5.206 | `cbc1701bd80144404658a308da6fd2589e746f4c2d116441787c224ad2cd2647` |
| `extracted_data.json` | 2.642 | `758f6b568230066ccd0bc489117631d117ce3817ae04f8b95532933929e44246` |
| `extracted_text.md` | 844 | `d286821c430fd8efc7515474a74225e52378315abd5b19dde4b5d90de5df806d` |

No se renombraron ni se ocultaron dentro del source adquirido. `OFFICIAL-UPSTREAM-ACQUISITION-CORE` 0.4.10 declara `supported_platforms=[linux,macos]` y rechaza Windows antes de iniciar la descarga. La regresión produjo:

```text
PLATFORM_GATE_PASS windows_rejected_before_network
```

## Condiciones reales

La lane agentic requiere, según la propia solución, cuenta AWS, identidad/permisos, Bedrock AgentCore, S3 Vectors, Aurora DSQL, Lambda y DynamoDB. La lane prompt-flow incorpora Bedrock, Textract, A2I, SQS, SNS y DynamoDB. Son servicios cloud que pueden generar costo. También faltan antes de cualquier promoción: lock de dependencias del target, threat/privacy/cost review, despliegue en Linux, tests oficiales o propios reproducibles, corpus autorizado, ground truth, schemas/reglas del negocio, carga, seguridad, recovery y comprobación de escritura idempotente.

Los datos incluidos por AWS son sintéticos y el propio repositorio exige evaluación independiente. Por eso esta evidencia no autoriza `READY_FOR_AUTOMATIC_STORAGE`.

## Fallos preservados

- `LIB-FAIL-080`: el extractor ZIP estándar de Windows no pudo materializar los segmentos `s3:`; recuperación: admisión fail-closed por plataforma antes de red.
- `LIB-FAIL-081`: el primer test auxiliar confundió `stderr` estructurado y luego el salto visual de PowerShell; recuperación: captura cruda de proceso y regex multilinea, conservando el rechazo original.
- `UP-FAIL-031`: archive oficial no portable a Windows por 15 entradas `s3:`.
- `UP-FAIL-032`: no hay lock integral ni suite unitaria integral observada; servicios AWS, cuenta y costo permanecen condicionados.

Fuente oficial gobernante: https://github.com/aws-samples/aws-ai-intelligent-document-processing/tree/31d671ef54c6694fea6431ee4102e3de85b982c1.
