# AWS GuardDuty → Google Magika → AWS IDP dispatch — evidencia V91

Fecha de corte: 2026-08-28. Esta evidencia gobierna el estado V91 y no reemplaza los gates live del proyecto.

## Fuentes oficiales fijadas

- AWS GenAI IDP Accelerator: tag `v0.6.5`, commit `1b5fd74454e593de233342a02ee911af8ee38359`, repositorio oficial `aws-solutions-library-samples/accelerated-intelligent-document-processing-on-aws`, licencia MIT-0.
- Google Magika CLI: release `cli/v1.1.0`, commit `5e2f437fb7b7452368c8c1fa9354858f5487a5c4`, licencia Apache-2.0.
- Asset Linux x86-64: `magika-cli-x86_64-unknown-linux-gnu.tar.xz`, 8.625.112 bytes, SHA-256 `6b4c1010c84d1f4f06205ccef4597f1690bcd7744f46d841eee26426bc100485`.
- Sidecar oficial: 110 bytes, SHA-256 `4028a7c1fe789dea5f84e879a6dec00733cf042eff3ee7089bee3962fd923e19`.
- Binario extraído único: 29.225.280 bytes, SHA-256 `03a56971f9e121666d702db611fd08f1c8bae5a3c31538f1b45b3ef99e6df84d`.
- GitHub API asset `404531370`, digest de release, descarga, sidecar y `gh attestation verify --repo google/magika` coincidieron. El hash distinto indexado por un buscador fue rechazado y registrado como fallo aprendido.

## Resultado materializado

`AWS-GUARDDUTY-MAGIKA-IDP-DISPATCH-GATE` 0.1.0 contiene 11 archivos. El pack canónico final SHA-256 `d192efe903a02d0c5db5d55dfbc7329cbd6142213f87d2755033c7a99ef45012` reconstruyó 11/11 archivos sin diferencias y pasó:

- 14/14 pruebas Python;
- verifier de contratos y allowlist fail-closed;
- cfn-lint 1.55.1, exit 0 sobre `template.yaml` SHA-256 `54d494d1f59e238129e00eced566f910733bd9fed82b1c8112087f6b86147997`;
- build con siete wheels Boto3/Botocore 1.43.83 fijados por hash y el asset oficial Magika;
- build aceptado: 2.193 archivos de contenido, 2.194 contando `BUILD_RECEIPT.json`, cero `.pyc`/`__pycache__`; receipt SHA-256 `0119a6c4a981237d38378965ee850a22c89d42d001580db74af2f9544e138498`;
- updater regression: nested fences, replacements literales, hash reconstruction, atomic failure y traversal PASS.

La primera construcción técnicamente exitosa contenía bytecode accidental y fue rechazada como release. El build se corrigió con `python -B`, se rematerializó desde Markdown y sólo el output limpio precedente se acepta como evidencia.

## Contrato demostrado

La ruta implementada es: recibo GuardDuty limpio e inmutable → validación exacta de cuenta/región/recurso/tenant/esquema/VersionId/eTag/KMS/hash/tamaño → identificación con binario Google Magika exacto → allowlist label/MIME/extensión/score configurada por proyecto → copia de la misma versión S3 con `CopySourceIfMatch`, `IfNoneMatch: *`, checksum, clave content-addressed y `config-version` IDP fijada → receipt Object-Locked confirmado último.

DynamoDB aporta lease con owner fencing; SQS usa partial batch response; duplicados y copia previa al receipt se reconcilian; logs exponen sólo hashes/tipos de error. `handler.py` es glue `ADAPTED`, no código verbatim atribuido a AWS o Google. Las licencias Google/AWS son `ADAPTED` únicamente porque el contrato Markdown agrega el LF final ausente/normaliza newline; los otros ocho bloques son `AUTHORED`.

## Estado global y no-afirmaciones

`VERIFY_LIBRARY_PASS`: 65 packs, 614 archivos materializables y 35 perfiles; email AWS compone 6 packs/63 archivos. `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit` reconoce exactamente un dispatch gate. Procedencia exterior real: 501 `AUTHORED` + 14 `ADAPTED` + 99 `VERBATIM` = 614.

No hubo deploy AWS, llamada IDP, costo cloud ni documento empresarial. El acelerador oficial v0.6.5 conserva intake at-least-once y las limitaciones `UP-FAIL-177`; por eso todo receipt fija `automatic_business_persistence_authorized=false`. Exactitud por clase/campo, revisión humana, idempotencia empresarial, corpus target, canaries, carga, recuperación y rollback siguen siendo condiciones live, no claims inferidos del PASS offline.
