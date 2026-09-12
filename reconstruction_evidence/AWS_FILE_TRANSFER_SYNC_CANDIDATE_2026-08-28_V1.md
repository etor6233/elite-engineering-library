# AWS File Transfer Sync Candidate — 2026-08-28 V1

## Resultado

`aws-samples/file-transfer-sync-solution` queda fijado como `SAMPLE_ONLY_REJECTED_FOR_IMMEDIATE_ADOPTION`.

Es código oficial AWS para listar un SFTP remoto y disparar Transfer Family Connectors hacia S3. No demuestra finalización por archivo, idempotencia, cuarentena ni almacenamiento confiable. Se conserva como referencia exacta; no se despliega ni se adapta.

## Identidad

- repository: <https://github.com/aws-samples/file-transfer-sync-solution>;
- commit GitHub-verified: `5ee79ee1c222d4b514c97345911ffb30b72f32a7`;
- archive: 215.738 bytes;
- archive SHA-256: `970817ec45289123ff0ddbffc599bd4e16d66f08376242b4b7f9d2a448f1be0c`;
- 22 archivos, 9 directorios, cero tags y cero releases;
- MIT-0 `LICENSE` SHA-256 `5025c7bbedbc8da868b6dce2ab689225b3f33c43b3f03368b0a331809b6ada26`;
- seis archivos Python compilan con Python 3.14.4.

Artefactos focales:

| Path | SHA-256 |
|---|---|
| `README.md` | `e851b585352567b88db80e8333d5e991324c9987efd6dd544a5f8d8f716159ef` |
| `requirements.txt` | `d87f23d0c5125f6097c9d2457a7a145b4f868e1bc24d4a39a61011c7c5ef6855` |
| `configuration/examples/example-sftp-sync.json` | `c12c219063ff07f5771a7df87fe78d6b6b26b2bc87b356e4b696fb4fed5a2a74` |
| `transfer_sync_service/transfer_sync_service_stack.py` | `4a707972ee73e0dc529c07937bb54246be494105893904be93aee3c7fad6b633` |
| `lambda/remote_server_list/remote_server_list.py` | `eb8af7f3e8f35832900867e83e4fe62939fa2ce6df1c76ca5b55750226e3883b` |
| `lambda/get_list_status/get_list_status.py` | `184ddef7746eeb996f42faf96d150185705e71cbd85b78b20efc0fc4f4f11ab0` |
| `lambda/sync_files/sync_files.py` | `53cf8438b79f69d6136c4bf18d86f3b7fa347912f1a968339d4cc22ffa78ca22` |

## Verificación y brechas

- cero test files, locks y workflows;
- dos manifests declaran once dependencias directas mediante ranges; no hay grafo frozen;
- el módulo CDK ejecuta `pip install` hacia el propio árbol source durante import/synth;
- trusted host keys, Secrets Manager, KMS opcional, IAM, schedules, Step Functions y monitoring sí existen;
- la detección de cambio compara `modifiedTimestamp` remoto con `LastModified` S3; no compara contenido, checksum ni versión remota;
- `StartFileTransfer` es asíncrono y entrega `transfer-id`; el código descarta la respuesta y no consume eventos/logs `COMPLETED|FAILED` por archivo. Autoridad AWS: <https://docs.aws.amazon.com/transfer/latest/userguide/events-detail-reference.html>;
- cada llamada admite diez paths y el sample particiona de a diez, pero dispatch no equivale a completion. Autoridad AWS: <https://docs.aws.amazon.com/transfer/latest/userguide/transfer-files-and-track.html>;
- `process_s3_pages` y `transfer_files` capturan errores y continúan;
- `lambda_handler` transforma cualquier excepción restante en `{"status":"error"}` retornado normalmente;
- el first-copy flag se escribe después del dispatch, no después de completion;
- harness aislado: un report corrupto fue capturado, el handler retornó `status=success` y escribió `quarantine/remote.flag`;
- el report bucket usa `DESTROY`, auto-delete y expiración a 30 días;
- no hay cuarentena/malware gate, receipt inmutable, reconciliation de transfer-id, duplicate policy ni autoridad para persistir.

## Admisión

Se incorpora a `OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.57` y al perfil AWS sólo como referencia externa rechazada. La revisión siguiente debe traer frozen dependencies, tests/SCA, checksum/version semantics, errores propagados, seguimiento `COMPLETED|FAILED`, idempotencia durable, receipts, cuarentena/malware y pruebas live remote-server/AWS de outage/replay/load.

El canal SFTP productivo sigue abierto. La presencia de CDK, Step Functions y Transfer Family no convierte dispatch asincrónico en ingesta completa.
