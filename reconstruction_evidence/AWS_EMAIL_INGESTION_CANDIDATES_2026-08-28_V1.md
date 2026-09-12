# AWS Email Ingestion Candidates — 2026-08-28 V1

## Resultado

Se auditaron tres repositorios públicos oficiales AWS adicionales para cerrar recepción documental por email. Ninguno queda admitido como implementación inmediata. Dos fuentes MIT-0 se fijan por commit/archive/hashes como referencias rechazadas; la tercera no se adquiere porque su archive no publica licencia.

Esto no crea un pack productivo de email ni autoriza persistencia. Evita que un agente confunda procedencia oficial, tests verdes o una arquitectura amplia con un contrato empresarial completo.

## Fuentes exactas

### aws-samples/serverless-mail

- repository: <https://github.com/aws-samples/serverless-mail>;
- commit head firmado/verificado: `70ac9310a313cc95f1092083810659961afb88ba`;
- canonical archive: 1.582.762 bytes;
- archive SHA-256: `b7b3df3d11865b4f419c827cace41cf28cd666608704d4f69b5cc3a2da801cba`;
- root `LICENSE`: MIT-0, SHA-256 `ef47b4ae2a1a8d38ef87b38dc5927967e7b472ffa3f7d125e89dbc8353a1e7af`;
- archive: 158 files;
- focal parser: `lambda-email-parser/lambda_function.py`, 6.458 bytes, SHA-256 `93f93ceee35207d83ab989dbb6303229e57dec11b09bc9a3e01b989544ff0da1`;
- focal README SHA-256: `ee82b5304d5a4c83c062669654eff4a17f9a203d33af78b17edad2c2583650ee`;
- tags/releases: 0/0.

Los 53 Python del archive pasaron parseo AST. El componente focal tiene cero tests, cero manifests de dependencias y cero workflows. Sólo procesa `Records[0]`; ante destino ausente o escritura recursiva retorna sin producir un fallo reintentable; conserva TODOs para email inválido, error de decode y filenames tainted; no limita bytes, partes ni profundidad, ni genera checksums, idempotencia, tenant binding o quarantine. Python actual emitió dos `SyntaxWarning` por escapes `\s` no raw.

Decisión: `SAMPLE_ONLY_REJECTED_FOR_IMMEDIATE_ADOPTION`.

### aws-samples/sample-aws-security-incident-response-email-integration

- repository: <https://github.com/aws-samples/sample-aws-security-incident-response-email-integration>;
- head firmado/verificado: `08f11d40029d6189905b41b7f8b7372ff94803cf`;
- el último commit focal del handler es unsigned: `b6252fc1d38dddcc1c51eac6109f66e719f2671f`;
- canonical archive: 213.459 bytes;
- archive SHA-256: `dd74fb1925f1610043b1ef6b30b45733c4030ca5c9d8ac4925eae3f53334cd9f`;
- `LICENSE`: MIT-0, SHA-256 `7cb750713252efd1d578837ba8785b61319109906f0c10b87adab7cf4badfc42`;
- archive: 16 files;
- handler SHA-256: `f3614cfdf4da2b2536588fbbc5ba7ed25596131dba10502e0073fcbb4e681b6f`;
- template SHA-256: `899bd1f4e68275734f6c1e32d240adab5db4bdc2b7634d37692df29077159f71`;
- tests SHA-256: `b3dc842d38e725b4e04ae97036ac3e86db899084602c6e8a32ca40ae3ec42a4d`;
- tags/releases: 0/0.

El repositorio fija sólo ranges `boto3>=1.35.0`, `pytest>=8.0.0` y `moto>=5.0.0`. El entorno aislado de esta auditoría resolvió 26 distribuciones y `pip check` pasó. La suite abortó inicialmente en collection por crear clientes boto3 globales sin región; con `AWS_DEFAULT_REGION=us-east-1`, metadata deshabilitada y credenciales ficticias ejecutó 42/42 tests en 0,44 s. Google OSV-Scanner 2.5.1 exacto —binario SHA-256 `25e42f5ef6711fd8c0fb45390972205891dd44c6bd02ac93f0f63e8e98d9bfb6`— consultó el freeze y devolvió `results: []`; reporte SHA-256 `54345d5194d238b6017d0f3181f392109ce35023cc6b4bdfa8c134943dcb578e`. Es un resultado fechado, no garantía futura.

Los tests no convierten el sample en ingesta documental:

- `parse_email` toma el primer `text/plain` y nunca extrae adjuntos;
- la autorización compara el header MIME `From` con watchers/equipo, pero el handler no consume verdicts SPF/DKIM/DMARC del receipt;
- S3 missing/failure, case missing y sender rechazado retornan normalmente y el handler finaliza 200;
- `template.yaml` crea `DeadLetterQueue` y su alarma, pero no hay `RedrivePolicy`, `DestinationConfig`, `EventInvokeConfig` ni binding de la queue al evento SNS/Lambda;
- email storage usa SSE-S3; faltan idempotencia, tenant binding, límites, checksum, quarantine/malware y reconciliación.

Decisión: `SAMPLE_ONLY_REJECTED_FOR_IMMEDIATE_ADOPTION`.

### aws-samples/sample-bda-redaction

- repository: <https://github.com/aws-samples/sample-bda-redaction>;
- head firmado/verificado: `2c89832896d21671591e71ff33c07c7ae32f84be`;
- canonical archive: 33.638.911 bytes;
- archive SHA-256: `6f0f345df09111e570d1b582cf5b8bc1f977457c1a44e74f80c07de56f5e211f`;
- archive: 93 files, 44 directories;
- `LICENSE`/`NOTICE`: cero;
- tags/releases: 0/0;
- Python AST: 17/17.

Aunque implementa body, adjuntos, BDA/Guardrails, DynamoDB, retención y portal, visibilidad pública sin licencia no concede permiso de copia o derivación. Además, `infra/tests` no contiene tests ejecutables; no hay locks; `generate_case_id` usa seis dígitos aleatorios, consulta DynamoDB con string pero persiste la key como integer y no usa conditional write; filenames de adjuntos se reutilizan; storage usa SSE-S3 y varios `RemovalPolicy.DESTROY`; no hay límites, checksums, idempotencia ni malware gate; README declara el acceso del portal opcional.

Decisión: `REJECTED_NO_LICENSE`. No se agrega al source lock adquirible.

## Integración en la biblioteca

- `OFFICIAL-UPSTREAM-ACQUISITION-CORE` avanza a 0.4.55;
- lock exacto: 103 sources;
- el perfil `aws-secure-document-pipeline` selecciona 18 sources y conserva ambos candidatos MIT-0 como referencias rechazadas;
- el lock y perfil ejecutan 9 negativos, 14 perfiles válidos, 5 negativos y 1 positivo offline;
- BDA sólo queda en el ledger de admisión y fallos por ausencia de licencia;
- `UP-FAIL-143`, `UP-FAIL-144` y `UP-FAIL-145` impiden promoción silenciosa.

## Conclusión honesta

La búsqueda produjo código oficial útil para estudiar MIME, SES/S3/SNS/Lambda, KMS/logging y un pipeline BDA amplio, pero no una implementación de recepción documental por email que satisfaga licencia, adjuntos, autenticación, idempotencia, delivery fail-closed, tenant isolation, límites/checksums, quarantine, retención y pruebas live al mismo tiempo.

Por lo tanto, recepción email productiva sigue abierta. Portal/API hacia cuarentena S3 continúa siendo la única entrada AWS ejecutable condicionada ya materializada; no se presenta ninguno de estos tres repositorios como código listo para almacenar documentos empresariales.
