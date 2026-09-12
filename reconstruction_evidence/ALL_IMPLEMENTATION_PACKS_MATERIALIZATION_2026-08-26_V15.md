# All implementation packs materialization — V15

Fecha: 2026-08-26.

## Alcance

Esta evidencia reemplaza V14 como snapshot gobernante de estructura y adquisición oficial. No reemplaza las pruebas ejecutables Foundation V12 ni convierte samples/accelerators en producto.

`OFFICIAL-UPSTREAM-ACQUISITION-CORE` 0.4.13 conserva el lock exacto de 66 fuentes y materializa 17 archivos. Añade tres perfiles separados y fail-closed para adquirir con un comando los candidatos oficiales de pipeline documental AWS, Google Cloud o Microsoft/Azure, después de completar una aprobación enlazada a los hashes exactos del perfil y el lock:

- `aws-secure-document-pipeline`: 7 fuentes oficiales AWS, incluyendo Document Processing, IDP, Textract, evaluación y GuardDuty malware;
- `google-secure-document-pipeline`: 6 fuentes, incluyendo Google Document Intake/Document AI/Terraform/Magika y los motores ClamAV/YARA-X fijados;
- `microsoft-secure-document-pipeline`: 11 fuentes, incluyendo Microsoft Content Processing 2.1.2, MarkItDown, Presidio, Azure Content Understanding/Document Intelligence/Table Transformer y ClamAV/YARA-X.

Los perfiles no mezclan proveedores automáticamente y no ocultan sus condiciones. Exigen cuentas/región/identidades, storage inmutable y cuarentena, malware/rules, schemas/corpus/ground truth, revisión humana, lifecycle/idempotencia/replay/reconciliación, telemetría, costos y autoridad de deploy/rollback. Todos conservan como blockers las pruebas productivas no realizadas y la naturaleza de sample/accelerator cuando corresponda.

Fuentes principales ya fijadas en el lock:

- <https://github.com/aws-samples/sample-document-processing/tree/9ca2eb1bf7b94f25c33926222817e0cd2d7748b4>
- <https://github.com/aws-samples/aws-ai-intelligent-document-processing/tree/31d671ef54c6694fea6431ee4102e3de85b982c1>
- <https://github.com/GoogleCloudPlatform/document-intake-accelerator/tree/956e3cc7900338d1146d82c5ee2bd4b2551f70bb>
- <https://github.com/microsoft/content-processing-solution-accelerator/tree/a23c24d581a886f91de10cfda7fe131571c14b22>

## Gates ejecutados

```text
Materialized 17 files
UPSTREAM_ACQUISITION_TEST_PASS negatives=5
SOURCE_PROFILE_VALID profile=secure-file-ingestion-leaders selected=6
SOURCE_PROFILE_VALID profile=secure-archive-linux selected=1
SOURCE_PROFILE_VALID profile=aws-secure-document-pipeline selected=7
SOURCE_PROFILE_VALID profile=google-secure-document-pipeline selected=6
SOURCE_PROFILE_VALID profile=microsoft-secure-document-pipeline selected=11
SOURCE_PROFILE_VALID profile=document-intelligence-leaders selected=27
SOURCE_PROFILE_VALID profile=commerce-communications-leaders selected=25
SOURCE_PROFILE_VALID profile=enterprise-platform-leaders selected=10
SOURCE_PROFILE_TEST_PASS valid=8 negatives=5
VERIFY_LIBRARY_PASS
packs=37 materialized_files=340 markdown_files=222
```

La composición de inicialización produce 23 archivos. Los perfiles que incluyen el pack completo de adquisición aumentan tres archivos: AWS Textract 24, AWS enterprise 26, Google Ads 22, Meta Ads 25, TikTok Ads 25, Meta WhatsApp 28, Firebase 26, Amazon SP-API 24 y Google Merchant 25. Azure document 11, Google document 11, Mercado Libre 10, payment 7, backend 95 y web/BFF 44 permanecen sin cambio.

## Fallos convertidos en gates

- `LIB-FAIL-092`: nombre histórico de parámetro del materializador rechazado.
- `LIB-FAIL-093`: `output.yarc` residual detectado por allowlist y eliminado de manera exacta.
- `LIB-FAIL-094`: limpieza compuesta rechazada por política; mutación y gate quedaron separados.
- `LIB-FAIL-095`: ID descriptivo AWS no pertenecía al lock; el perfil falló antes de red/escritura y fue corregido contra el identificador canónico.

## Decisión de admisión

Los perfiles están `REBUILD_VERIFIED / CONDITIONED` como mecanismo de selección y adquisición exacta. El código upstream queda inmediatamente adquirible después de aprobación, pero ningún perfil está promovido como pipeline productivo de un proyecto. Esa promoción requiere ejecutar el código elegido en la cuenta real, cerrar sus fallos upstream, probar corpus/ground truth por campo, seguridad de archivos hostiles, aislamiento, revisión, persistencia/linaje, reconciliación, carga, restore y rollback.
