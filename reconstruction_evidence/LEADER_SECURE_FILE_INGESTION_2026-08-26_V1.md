# Leader Secure File Ingestion — 2026-08-26 V1

## Resultado

Se fijaron y auditaron tres fuentes oficiales de Google, Amazon Web Services y Microsoft que cubren capas distintas de una ingesta empresarial: identificar el contenido real, escanear/aislar malware y procesar múltiples documentos mediante extracción, transformación, evaluación, persistencia y reglas cruzadas. No se combinan en un producto ficticio ni se etiquetan todas como verdes: cada una conserva exactamente sus pruebas y fallos.

## Google Magika CLI 1.1.0

- repositorio: `google/magika`;
- release: `cli/v1.1.0`;
- commit: `5e2f437fb7b7452368c8c1fa9354858f5487a5c4`;
- source ZIP: 84.532.977 bytes, SHA-256 `54983bcea9c11499aee7162ed11989167918b37b69adce67a4891f2fd3e1d7a0`;
- Apache-2.0, LICENSE SHA-256 `58d1e17ffe5109a7ae296caafcadfdbe6a7d176f0bc4ab01e12a689b0499d8bd`;
- Rust CLI lock: `rust/cli/Cargo.lock`, SHA-256 `6aed136fe7d903f13a2cd8e64ee28632246e8788044631f3aabadf238e13962e`;
- release Windows x64: 10.076.032 bytes, SHA-256 `d4de347c53e2d25f9663780c93fdf26d156284c64d6a5ec4bcdb829f9994f6f3`;
- checksum sidecar: 105 bytes, SHA-256 `2a7da12574b85a50f1a98e24fb8ed4f6e096eefac395c1349a7d103491e5d91d`.

El binario oficial respondió `magika 1.1.0 standard_v3_3`. Sobre fixtures exactos del source y una copia byte-idéntica de PDF renombrada localmente como `invoice.txt`:

```text
MAGIKA_CONTENT_GATE_PASS labels=pdf,png,python,pdf
```

Esto demuestra detección por contenido para esos cuatro casos, no precisión universal ni ausencia de malware. Magika es `PINNED_CANDIDATE` para el gate de tipo y nunca reemplaza un scanner antimalware.

## AWS GuardDuty Malware Protection for S3

- repositorio: `aws-samples/guardduty-malware-protection`;
- commit: `fd9d2bdf46eb52bb336a7e7a837bd38ef73002e2`, fechado 2025-11-20;
- source ZIP: 103.641 bytes, SHA-256 `0598bcdd176cbf33df987be126c6849cc1aa040c6866ec643b28d3b2a9ee17f0`;
- MIT-0, LICENSE SHA-256 `7cb750713252efd1d578837ba8785b61319109906f0c10b87adab7cf4badfc42`;
- `cdk/package-lock.json`: 271.990 bytes, SHA-256 `6a1b223859cdd214f97a451463494f024e0d06e2c00d47a013220a32f245b93b`.

El código oficial aporta CloudFormation y CDK para bucket protegido, GuardDuty Malware Protection Plan, EventBridge, SQS/DLQ, Lambda, bucket limpio y escenarios de copia. `npm ci --ignore-scripts` instaló 307 paquetes; `tsc` pasó. La admisión se detuvo porque:

- `npm audit` reportó 13 vulnerabilidades: 2 bajas, 4 moderadas y 7 altas;
- el repo fija `aws-cdk 2.148.1` con `aws-cdk-lib 2.177.0`; `cdk synth` falló por cloud assembly schema 39 contra máximo CLI 36;
- Node 24 no está soportado por esa release, que declara Node 22;
- GuardDuty, S3, KMS, Lambda, SQS y demás recursos requieren cuenta, IAM, región y pueden generar costo.

Queda `SAMPLE_ONLY_CONDITIONED`; no es un scanner local listo ni puede entrar en producción hasta resolver el grafo, repetir audit/build/synth y probar cuarentena/release/incident response en la cuenta elegida.

## Microsoft Content Processing Solution Accelerator 2.1.2

- repositorio: `microsoft/content-processing-solution-accelerator`;
- release: `v2.1.2`;
- commit: `b47cec48475cfedd7109debf2f407ef0c5530311`;
- source ZIP: 18.706.401 bytes, SHA-256 `b2f0a20055111420558c1903eaa8ca570d835d309968b30650f0f86989c373fb`;
- MIT, LICENSE SHA-256 `d9a1b1e30d633d5732ea18e3cba9538d293ebc53e1a9e4e96ab739e0c5c4f1cb`;
- 701 archivos y 340 Python source; AST Python 3.12: 340 PASS, 0 fallos;
- locks exactos:
  - `ContentProcessor/uv.lock`: 682.287 bytes, SHA-256 `af76bfc133211a2ca979a2e5ad1e885ebbe73b502ed8f2ac2133e99c4bc038ba`;
  - `ContentProcessorAPI/uv.lock`: 532.434 bytes, SHA-256 `de0994f57cf9d4574f78fd2cf0f682e62e865f5385646280d4b00e25c4f42e1a`;
  - `ContentProcessorWorkflow/uv.lock`: 672.405 bytes, SHA-256 `9b9e51b47a9e5e002e5784a23c43671ef9f5b997a94b0fbb90671df1efb4629e`;
  - `ContentProcessorWeb/pnpm-lock.yaml`: 573.420 bytes, SHA-256 `0441faf082f787edc66f6d1b11a6fa26ee02f276c6020728f508eceb4bd07998`.

El código oficial implementa upload validation, MIME sniffing, límites, Blob/Queue/Cosmos, schemas, scoring, el pipeline `Extract → Map → Evaluate → Save`, API/UI y un workflow multi-documento con retries/DLQ, resumen y reglas declarativas de gaps/discrepancias. Las pruebas exactas produjeron:

- `ContentProcessorWorkflow`: 268 PASS, 3 warnings;
- `ContentProcessor`: instalación frozen de 205 paquetes; suite integral bloqueada por 2 errores de colección relacionados con `agent_framework`; excluyendo sólo esos dos módulos para medir el resto, 208 PASS y 1 fallo upstream de test obsoleto de credenciales;
- `ContentProcessorAPI`: suite integral bloqueada al requerir `APP_CONFIG_ENDPOINT` durante colección; midiendo el resto, 198 PASS y 15 fallos de helper/DI mocks;
- `ContentProcessorWeb`: instalación frozen de 1.506 paquetes con `pnpm 10.28.2`; 12 suites Jest no colectaron tests por transform TypeScript/JSX; build productivo PASS, con `caniuse-lite` siete meses antiguo y deprecaciones PostCSS;
- Azure live/E2E no se ejecutó sin suscripción, cuota, identidad ni costo autorizados.

La primera instalación Microsoft desde una ruta larga de Windows falló al persistir metadata `dist-info`; la repetición desde raíz corta con los mismos locks instaló correctamente. No se modificó código ni lock para convertir fallos upstream en PASS.

Queda `SAMPLE_ONLY_CONDITIONED`. Es la fuente oficial más cercana observada a un flujo empresarial multiarchivo tangible, pero requiere corregir o actualizar sus fallos, fijar el target Azure, ejecutar seguridad/E2E y sustituir los schemas/rules de ejemplo por corpus y ground truth aprobados.

## Integración en Elite

`OFFICIAL-UPSTREAM-ACQUISITION-CORE` 0.4.11 fija 63 fuentes. El perfil `document-intelligence-leaders` selecciona 25, incorpora estas tres y mantiene el AWS agentic incompatible con Windows como opt-in Linux/macOS separado. El runner exige aprobación enlazada a hashes y respuestas sobre tipos/límites, antimalware/cuarentena, cuentas/costos, corpus/schemas y política de datos.

La ruta permitida queda:

```text
original + SHA-256
→ Magika content-type gate
→ malware scan/quarantine del provider seleccionado
→ almacenamiento inmutable
→ clasificación/separación/extracción oficial
→ schema + evaluación + reglas cruzadas
→ revisión o aceptación demostrada
→ persistencia idempotente + provenance
```

Esta secuencia aún es un contrato de composición entre fuentes oficiales, no un pack unificado productivo. No autoriza afirmar ingesta perfecta o almacenamiento automático universal.

Fuentes oficiales gobernantes: https://github.com/google/magika/tree/5e2f437fb7b7452368c8c1fa9354858f5487a5c4, https://github.com/aws-samples/guardduty-malware-protection/tree/fd9d2bdf46eb52bb336a7e7a837bd38ef73002e2 y https://github.com/microsoft/content-processing-solution-accelerator/tree/b47cec48475cfedd7109debf2f407ef0c5530311.
