# AWS SES Mail Manager Attachment Pipeline — exact admission evidence V1

Fecha de corte: 2026-08-28.

## Fuente, licencia e integridad

- repositorio oficial: `aws-samples/sample-amazon-ses-mail-manager-attachment-pipeline`;
- commit exacto sin firma: `79314a93bda431bd03a4c1236bc4b22e75e2f41f`;
- archive: 709.657 bytes, SHA-256 `e1e241c4b8a0765f913bff04cf1eb7b715dfbeaace580b398c181eda5cd30fd5`;
- 23 archivos;
- licencia `MIT-0`, `LICENSE` 946 bytes/SHA-256 `5025c7bbedbc8da868b6dce2ab689225b3f33c43b3f03368b0a331809b6ada26`;
- `NOTICE` 126 bytes/SHA-256 `46ac37221a2328550b6963048bce2a787cabea518e6e2419e19e02ec3e2c923b`.

El source exacto no se copia ni se presenta como código authored por Elite. `OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.54` conserva su URL inmutable, archive, licencia, artefactos focales y condiciones para adquirirlo sin Git después de aprobación hash-linked.

## Código oficial observado

El sample aporta código real de AWS para:

- SES Mail Manager traffic policy y rule set;
- allowlist de destinatarios;
- Abusix y Trend Micro add-ons;
- archive, escritura de MIME a S3 e invocaciones Lambda;
- buckets versionados, TLS obligatorio, bloqueo público y `RETAIN`;
- tablas DynamoDB con PITR;
- extracción MIME de adjuntos y routing por destinatario;
- categorización advisory con Amazon Bedrock y validación de valores permitidos.

Once archivos Python pasaron compilación AST. El repositorio no contiene tests automatizados ni lock transitive.

## Dependencias y synth reproducible

El `cdk/requirements.txt` oficial fija `aws-cdk-lib==2.251.0` y `constructs==10.6.0`. Un entorno aislado resolvió 13 paquetes. Google OSV-Scanner 2.5.1, binario oficial SHA-256 `25e42f5ef6711fd8c0fb45390972205891dd44c6bd02ac93f0f63e8e98d9bfb6`, detectó `GHSA-464c-974j-9xm6` en `aws-cdk-lib 2.251.0`; la versión corregida oficial es 2.253.0.

Además, `setup.sh` solicita `aws-cdk@2.251.0`, versión que npm no publica. La combinación separada y explícita `aws-cdk-lib 2.253.0` + CLI oficial `aws-cdk 2.1139.0`:

- resolvió 13 paquetes y `pip check` pasó;
- devolvió cero coincidencias conocidas en el scan OSV fechado;
- sintetizó el source exacto en 32 recursos CloudFormation;
- emitió deprecaciones DynamoDB/Lambda y reportó 79 feature flags sin configurar.

Esto demuestra que una actualización oficial compatible puede sintetizar el stack; no convierte el commit original ni su `setup.sh` en un despliegue reproducible o seguro.

## Brechas materiales que impiden adopción inmediata

- AWS declara que es un sample de demostración y prescribe hardening productivo;
- el commit es unsigned, no tiene release, tests ni dependency lock;
- el pin exacto tiene un advisory conocido y el CLI solicitado no existe;
- add-ons Abusix/Trend Micro y Bedrock implican términos/costos y cuenta real;
- usa SSE-S3/AWS-managed encryption, no customer-managed KMS;
- acciones críticas usan `CONTINUE` ante fallos;
- toma sólo el primer recipient;
- deriva prefijo sólo del local-part, por lo que dominios distintos pueden colisionar;
- nombres iguales de adjuntos pueden sobrescribirse dentro del mismo prefijo;
- no hay contrato de tamaño/cantidad/checksum/idempotencia por adjunto;
- elimina el MIME de recepción después de copiar, aunque exista un archive separado;
- la clasificación Bedrock es advisory y no autoriza almacenamiento de campos empresariales.

## Decisión

Estado: `SAMPLE_ONLY_REJECTED_FOR_IMMEDIATE_ADOPTION`.

El agente puede adquirir el commit exacto para reutilizar sus contratos Mail Manager y comparar una implementación, pero no puede desplegarlo, persistir datos confiables ni llamarlo pack elite listo. Antes de reconsiderar se exige un revisionado oficial o adaptación explícitamente `AUTHORED` con dependencies hash-locked/clean, tests, retención inmutable, binding tenant no colisionable, claves únicas, límites/checksums/idempotencia, SSE-KMS, fallo cerrado, sandbox malware/outage/replay/load, costos aprobados y composición con la cuarentena segura ya materializada.

## Regresión canónica

- `OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.54` materializó 23/23 archivos;
- lock: 101 fuentes únicas;
- adquisición: nueve negativos PASS;
- catorce perfiles source validaron, cinco negativos y una adquisición positiva offline PASS;
- perfil AWS selecciona 16 fuentes e incluye esta revisión rechazada antes de pedir aprobación.

Fuentes oficiales: <https://github.com/aws-samples/sample-amazon-ses-mail-manager-attachment-pipeline/tree/79314a93bda431bd03a4c1236bc4b22e75e2f41f>, <https://github.com/aws/aws-cdk/security/advisories/GHSA-464c-974j-9xm6>, <https://github.com/aws/aws-cdk/releases/tag/v2.253.0> y <https://github.com/google/osv-scanner/releases/tag/v2.5.1>.
