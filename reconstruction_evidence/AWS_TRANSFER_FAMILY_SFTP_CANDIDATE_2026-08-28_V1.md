# AWS Transfer Family SFTP Candidate — 2026-08-28 V1

## Resultado

`aws-samples/sample-secure-transfer-family-code` queda fijado como `SAMPLE_ONLY_REJECTED_FOR_IMMEDIATE_ADOPTION`.

Es código oficial y licenciado de AWS que demuestra una arquitectura SFTP→S3→GuardDuty→routing. No es un componente listo para ingreso documental productivo: carece de pruebas, locks, aislamiento por usuario, idempotencia y delivery fail-closed. No autoriza liberar archivos al bucket limpio ni persistir información de negocio.

## Identidad exacta

- repositorio: <https://github.com/aws-samples/sample-secure-transfer-family-code>;
- commit: `474ca07cbc9a4937f09e6b5f4f51723f04a10583`;
- fecha del commit: `2025-08-21T14:35:56Z`;
- verificación GitHub: `unsigned`;
- archive: <https://github.com/aws-samples/sample-secure-transfer-family-code/archive/474ca07cbc9a4937f09e6b5f4f51723f04a10583.zip>;
- archive bytes: `149349`;
- archive SHA-256: `5ec5b157b285a5cc505cb5b8a2a508ce299d02a6bf71359f53cb0594c6d2b9e2`;
- archivos/directorios: `7/0`;
- tags/releases: `0/0`;
- licencia: `MIT-0`, `LICENSE`, 947 bytes, SHA-256 `7cb750713252efd1d578837ba8785b61319109906f0c10b87adab7cf4badfc42`.

Artefactos focales:

- `README.md`: 26.977 bytes, SHA-256 `c1ee6eb6fcc2bec6ed876849ca986c89bee3d825a3b3b23cf4df375207b47567`;
- `secure-transfer-family-template.yml`: 60.966 bytes, 1.716 líneas, SHA-256 `0357b88863665a5204f10dcff347df07fb2ac0903edef47e377930bc1aeebb9c`.

## Código realmente presente

El template crea VPC, subnets, endpoints, Transfer Family SFTP, Cognito, KMS, buckets de upload/clean/malware/error, access logs, EventBridge, SNS y tres Lambdas inline.

Los tres bloques Python parsearon con Python 3.14.4:

| Bloque | Líneas | SHA-256 | AST |
|---|---:|---|---|
| reglas de security group | 180 | `c36d76084dd5d2817b4eea1801b5bf36ec4788a0bfe60c250d6e35d57d036f2b` | PASS |
| autenticación Cognito | 56 | `0229029c6cc106e7c341bd37690672b845bdf74d498db2fc959be72d0268d5ee` | PASS |
| routing por malware | 66 | `df86570313ea551025909abc1b459c042ea6d5932a6224a23d60ce3d56f0f6dd` | PASS |

La validación aislada usó `aws-cloudformation/cfn-lint` oficial `1.53.3`, 27 distribuciones y `pip check` PASS. Resultado: exit 4, cero errores y 20 warnings (`W8001`×2, `W1020`×1, `W3005`×17). Los dos `W8001` demuestran conditions declaradas pero no usadas.

El repositorio publica cero test files, dependency manifests, dependency locks y workflows.

## Brechas verificadas

1. El template no crea `AWS::GuardDuty::MalwareProtectionPlan`, no consume `GuardDutyMalwareScanStatus` y no aplica TBAC. El README exige habilitar GuardDuty externamente.
2. AWS documenta que Malware Protection para S3 entrega resultados at-least-once y recomienda manejar duplicados: <https://docs.aws.amazon.com/guardduty/latest/ug/how-malware-protection-for-s3-gdu-works.html>.
3. El handler no usa `versionId`, ETag, sequencer, checksum ni clave de idempotencia. Ejecuta `copy_object` y luego `delete_object`.
4. Un harness aislado demostró `copy→delete` en success. Al forzar `copy_object` a fallar, el handler capturó la excepción y retornó `{"statusCode":500}` normalmente; no propagó error. EventBridge no recibe una invocación fallida para reintentar.
5. No existen `RetryPolicy`, DLQ, Lambda destination, reconciliation ni receipt inmutable.
6. Todo usuario autenticado recibe el mismo `HomeDirectory` del upload bucket y el mismo role puede listar, leer, escribir y borrar todos sus objetos. No hay directorio lógico por usuario/tenant.
7. `MfaConfiguration` está `OFF`; no hay public-key auth implementada.
8. El server fija `TransferSecurityPolicy-2020-06`. AWS recomienda actualizar a la política más reciente y documenta políticas posteriores, incluida 2025-03: <https://docs.aws.amazon.com/transfer/latest/userguide/security-policies.html>.
9. No hay límites de tamaño, checksum, object-version binding, retención inmutable ni prueba de outage/replay/load.

## Admisión

El source exacto se añade a `OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.56` y al perfil AWS únicamente como referencia externa rechazada. El source lock verifica 104 IDs únicos; el perfil AWS selecciona 19. Las suites canónicas pasan `9` negativas de adquisición y `14` perfiles válidos, `5` negativos y `1` positivo.

Para reevaluar se exige una revisión oficial con, como mínimo:

- aislamiento por usuario/tenant y role/prefix;
- política criptográfica vigente y MFA o key auth aprobada;
- plan GuardDuty y TBAC definidos;
- límites, checksums y binding a versión exacta;
- idempotencia para delivery duplicado;
- error propagado, retry/DLQ/reconciliation y receipt durable;
- tests, locks, SCA y sandbox AWS real con malware/outage/replay/load.

Hasta entonces el canal SFTP productivo continúa abierto. La biblioteca no modifica este código ni atribuye una corrección local a AWS.
