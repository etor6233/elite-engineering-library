# Pipeline documental de referencia

Infraestructura local/fixtures: recepción → SDK Textract → revisión → commit.
AUTHORED glue; SDK AWS original fijado; fixture Microsoft con licencia MIT exacta.
FIXTURE queda explícito en cada original, evidencia y commit. No es OCR live.

1. Aplicar las migraciones seleccionadas a la DB de referencia, incluida0085.
2. Adquirir el JPEG con `azure_document_intelligence_official_invoice/acquire_official_fixture.ps1` y su lock/approval. Puede importarse offline desde caché con `-ImportDirectory`; el acquirer verifica184686bytes/SHA. Retener `upstream/LICENSE.txt`. El template no habilita archivos privados.
3. En el host, habilitar `DOCUMENTS_ENABLED=true`, `DOCUMENTS_MODE=FIXTURE`, `DOCUMENTS_PROFILE_FILE=config/documents/reference-profile.json`, su SHA256 exacto, `DOCUMENTS_WORK_ROOT` directorio de trabajo existente, `DOCUMENTS_TENANT_ID` y `DOCUMENTS_ORGANIZATION_ID`. El IdP fixture existente entrega permisos explícitos; ninguna variable de entorno otorga roles.
4. PUT `/v1/documents/{uuid}/original`, `Content-Type: application/octet-stream`, headers `X-Document-Name`, `X-Document-SHA256`, `X-Document-Profile-SHA256` y Bearer del uploader. Body: bytes del JPEG fijado.
5. POST `/v1/documents/{uuid}/process` con `{}` y `documents:process`. GET del recurso recupera estado. GET `/original` y `/evidence/security|provider|analysis` conserva bytes/hash para revisión.
6. POST `/review` con `evidence_sha256` y `fields:{invoice_number,vendor,total,currency}`. POST `/decision` por otra persona, con `payload_sha256`, `approved` y `reason`. GET/repetición exacta recuperan resultado después de perder una respuesta.

Verificación en DB sintética loopback ya migrada:

```text
DOCUMENT_CONNECTED_DB_URL=postgres://.../elite_document_reference_<id>
python tools/verify_document_reference.py --go <Go1.26.8-admitido> --fixture <sample_invoice.jpg> --receipt <nuevo-receipt.json>
```

El runner usa DB URL sólo por entorno, nunca en el receipt ni argumentos. Fixture
público SHA `489f0c63b6a05e0ae0fcfbf29121299fc2ccacab35d5cd40a5d63d42764078cb`.
Registra cada prueba y falla si aparece SKIP. El proveedor fixture no abre red.

Para activar PROVIDER en tu futuro target: default credential chain AWS; perfil
por clase/region y runtime de seguridad existente mediante `DOCUMENTS_PYTHON`,
`DOCUMENTS_PYTHON_SHA256`, `DOCUMENTS_SECURITY_SCRIPT`,
`DOCUMENTS_SECURITY_SCRIPT_SHA256`, `DOCUMENTS_SECURITY_POLICY`,
`DOCUMENTS_SECURITY_POLICY_SHA256`, `DOCUMENTS_MAGIKA`, `DOCUMENTS_CLAMSCAN`,
`DOCUMENTS_CLAM_DATABASE`, `DOCUMENTS_YARA`, `DOCUMENTS_YARA_RULES`. Todos los paths
deben ser absolutos y los hashes reales. No pongas secretos en estos archivos.
Lee DOCUMENT_REFERENCE_DECISION.md antes de admitir corpus privado. El rechazo
por falta de scanner no habilita una simulación automática como fallback.
