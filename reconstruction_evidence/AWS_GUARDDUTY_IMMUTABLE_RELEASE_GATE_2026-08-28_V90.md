# AWS GuardDuty Immutable Release Gate — V90

## Resultado

El corte V90 agrega un pack materializable de 9 archivos y extiende el perfil email a 5 packs/52 archivos. El componente liga la decisión GuardDuty al plan, cuenta, región, bucket, key, `VersionId`, eTag y tag de la versión exacta. Sólo `COMPLETED/NO_THREATS_FOUND` puede copiar bytes a un key content-addressed retenido; los demás estados producen rechazo inmutable y alerta. `automatic_business_persistence_authorized=false` permanece fijo.

Estado global antes de esta evidencia: `VERIFY_LIBRARY_PASS` con 64 packs, 603 archivos materializables, 35 perfiles y 408 Markdown. `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit` confirmó `aws_guardduty_immutable_release_gates=1`, 111 fuentes upstream y el resto del corpus ejecutable. Esta evidencia eleva el conteo Markdown a 409 sin cambiar archivos materializables.

## Autoridad oficial y divergencias honestas

- Repositorio AWS oficial fijado: `aws-samples/guardduty-malware-protection`, commit `fd9d2bdf46eb52bb336a7e7a837bd38ef73002e2`, MIT-0.
- AWS documenta evento `GuardDuty Malware Protection Object Scan Result` con `versionId`, eTag, statuses y amenazas; entrega al menos una vez y recomienda idempotencia.
- AWS documenta el tag `GuardDutyMalwareScanStatus`, TBAC y los estados `NO_THREATS_FOUND`, `THREATS_FOUND`, `UNSUPPORTED`, `ACCESS_DENIED`, `FAILED`.
- El sample upstream exacto no se promueve: copia por key, no fija versión/eTag/tag dentro del handler, imprime la key, usa recursos eliminables, carece de recibo durable y su script de prueba no devuelve error al imprimir FAIL. cfn-lint 1.55.1 sobre su template terminó exit 4 con nueve warnings.
- `handler.py` y `template.yaml` se declaran `ADAPTED`, no `VERBATIM`; seis archivos son `AUTHORED` y el notice MIT-0 es `ADAPTED` sólo por normalización de newline. El pack suma 6 `AUTHORED` + 3 `ADAPTED`, dejando el corpus en 493/11/99 = 603.

Referencias oficiales verificadas:

- https://github.com/aws-samples/guardduty-malware-protection/tree/fd9d2bdf46eb52bb336a7e7a837bd38ef73002e2
- https://docs.aws.amazon.com/guardduty/latest/ug/monitor-with-eventbridge-s3-malware-protection.html
- https://docs.aws.amazon.com/guardduty/latest/ug/how-malware-protection-for-s3-gdu-works.html
- https://docs.aws.amazon.com/guardduty/latest/ug/tag-based-access-s3-malware-protection.html
- https://aws.amazon.com/blogs/security/using-amazon-guardduty-malware-protection-to-scan-uploads-to-amazon-s3/

## Invariantes probados

1. Autoridad exacta: EventBridge version 0, source AWS, detail type, cuenta, región y único plan ARN.
2. Schema 1.0, `S3_OBJECT`, status mapping cerrado y shape de threat/reasons.
3. Sólo keys emitidos por el receptor, tenant exacto, versión obligatoria y eTag válido.
4. Tag GuardDuty de la versión exacta debe coincidir con el evento.
5. La lectura source sólo ocurre para clean y exige SSE-KMS/key, tamaño, tenant y hashes receiver.
6. Copia por `VersionId` + `CopySourceIfMatch` + `IfNoneMatch`, destino content-addressed, KMS y Object Lock verificado.
7. Malware no se lee/copia y el nombre se persiste sólo como hash.
8. Unsupported/access denied/failed nunca se transforman en clean.
9. Lease DynamoDB reclaimable, owner fence, durable receipt, idempotencia post-TTL y reconciliación de copia parcial.
10. Batch SQS devuelve únicamente IDs fallidos; logs no incluyen key, filename, threat name ni contenido.
11. Recibo se escribe último y nunca habilita persistencia de negocio.
12. IaC: GuardDuty plan/tag/prefix, EventBridge exacto, SQS/DLQ cifradas, visibility 6× timeout, partial batch, PITR/TTL/deletion protection, buckets KMS/version/Object Lock/TLS/public-block/Retain y alarmas.

## Ejecuciones

- `python -m unittest -v test_handler.py`: 16/16 PASS.
- `verify_contract.ps1`: tests, compile, notice SHA, tokens y profile fail-closed PASS.
- cfn-lint 1.55.1 sobre template adaptado: exit 0, cero findings.
- materialización desde Markdown: 9/9; verifier repetido sobre reconstrucción limpia PASS.
- build con `--require-hashes --only-binary=:all:`: Boto3/Botocore 1.43.83 y siete wheels; receipt con 2.191 archivos PASS.
- `VERIFY_LIBRARY.ps1`: 64/603/408 y perfil email 52 PASS.
- `VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Audit`: PASS con el nuevo contador en uno.

No se ejecutó AWS LIVE: faltan autoridad de cuenta, costo aprobado, quarantine/KMS target, plan activo, confirmación SNS, canaries clean/EICAR/unsupported, aislamiento global, DLQ/redrive y rollback. El pack permanece `REBUILD_VERIFIED / CONDITIONED`, no demuestra ausencia universal de malware ni precisión documental.

## Identidades

Pack SHA-256: `8f221fcbdb0ccc58d3b39dad837d64580625c8dc0fcc0b3612ef0ab1fe60a052`.

| Archivo | Bytes | SHA-256 |
|---|---:|---|
| `AWS_MIT_NO_ATTRIBUTION_LICENSE.txt` | 946 | `5025c7bbedbc8da868b6dce2ab689225b3f33c43b3f03368b0a331809b6ada26` |
| `README.md` | 4285 | `388739cbd2feb0f1cf5b707e48479918c07a124703e59dab7f2926d36c37daa9` |
| `build.ps1` | 3371 | `65c866fba913939f9da5c9a0e9f4cbf36659ca7d6e77026c2aa3e5803bb9fa62` |
| `deploy.ps1` | 5157 | `dbcc40302f6db2927406bab8b2ade069f597dd0d8156dd7117b5fc63d24f771e` |
| `handler.py` | 24536 | `02d8127e4fe78710681bc473249b7fa491b1418cf470256bc798fa7d6881c636` |
| `provider-profile.template.json` | 735 | `b723a187807051daca364d337e518be78cf305f3f24727d0ee13b3843fe1a2ce` |
| `template.yaml` | 18223 | `ec58aa530de0b4a6d0d04c8d5d644100f25288b00ccf4a5f0b586badb53399ea` |
| `test_handler.py` | 18773 | `212630aedb720da587ff9ab41ded3982c0b9a12d5ee1b703db7849063eddd480` |
| `verify_contract.ps1` | 3369 | `73b12d935d9568950a33e5c6cad797225bc4e648c1271b50eba34d9a2ee142c7` |

## Memoria

El corte incluye 1.011 fallos locales y 176 condiciones upstream, 1.187 IDs únicos, cero abiertos locales y cero duplicados. Los fallos de tooling, propiedades CloudFormation, normalización de licencia, path temporal, parámetros, patch ancho e interpolación de arrays quedaron registrados antes del cierre.
