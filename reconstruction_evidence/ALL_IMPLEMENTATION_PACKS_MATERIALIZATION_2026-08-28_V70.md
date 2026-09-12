# All Implementation Packs Materialization — 2026-08-28 V70

## Cambio gobernante

`OFFICIAL-UPSTREAM-ACQUISITION-CORE` avanzó de `0.4.57` a `0.4.58` y de 105 a 106 fuentes exactas. Se añadió únicamente como referencia rechazada:

- `aws-samples/sample-ai-receipt-processing-methods` commit unsigned `6eb72bf0060e4fb38074d86fe9242e781db9df97`;
- archive 8.640.517 bytes, SHA-256 `844add2f964397361331d3554cc1ac9bd0ef5d5c3b8b07812b1c693740253674`;
- MIT-0 SHA-256 `7cb750713252efd1d578837ba8785b61319109906f0c10b87adab7cf4badfc42`;
- classification `SAMPLE_ONLY_REJECTED_FOR_IMMEDIATE_ADOPTION`.

No se copió el código al producto ni se escribió una corrección local atribuida a AWS. El source aporta cámara PWA, presigned POST, S3 event bridge, Step Functions y Textract/BDA/LLM, pero no satisface adopción inmediata: suite declarada ausente, SCA productiva 14/13, debug/test pre-auth, log de Authorization, demo sin versionado/PITR, POST sin checksum/content security, evento sin version binding que absorbe `StartExecution` failures y persistencia automática con confianza global mayor a 50 sin CAS/evidencia por campo/review.

Evidencia focal: `AWS_RECEIPT_MOBILE_INGESTION_CANDIDATE_2026-08-28_V1.md`.

## Source pack

- versión: `0.4.58`;
- bytes: `396573`;
- SHA-256: `94b50e49f237bb65eca633af3b08625dd2a95136ffbf4860be1a0210cb4be312`;
- archivos materializados: `23`;
- fuentes: `106/106`, IDs únicos;
- perfil AWS: `21/21`;
- perfiles source: `14` válidos, `5` negativos, `1` positivo;
- adquisición: `9` negativos PASS;
- planes consumidores: `20/20` en `0.4.58`, `0` en `0.4.57`.

La reconstrucción fresca desde Markdown materializó 23 archivos. `test_upstream_acquisition.ps1`, `test_source_profiles.ps1` y `acquire_upstream_sources.ps1 -ValidateOnly` pasaron; el perfil AWS seleccionó exactamente 21 fuentes y el lock completo validó 106/106.

## Ejecución focal del source

- 19 JavaScript y 8 Python parsearon;
- dos installs frozen con scripts deshabilitados completaron;
- frontend Vite/PWA y TypeScript de infraestructura compilaron;
- el test oficial falló porque `tests/integration/` no está publicado;
- SCA productiva: root 14 (6 moderate/8 high), frontend 13 (2 moderate/11 high);
- tags y releases reales: 0/0 después de excluir respuestas nulas.

## Verificación global

Después del cambio:

```text
VERIFY_LIBRARY_PASS
packs=53 materialized_files=499 markdown_files=366

VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit
packs=53 upstream_sources=106 document_sdk_artifacts=10
official_invoice_samples=1 official_google_process_samples=1
official_dapr_outboxes=1 official_pg_durable_handoffs=1
business_central_artifacts=1 provider_adapters=9
document_orchestrators=1 evidence_logs=1
```

Esta evidencia gobierna sólo después de que ambos verificadores produzcan exactamente esos sellos con cero filas abiertas.

## Memoria de fallos

- locales al cierre del snapshot: `776`;
- upstream: `148`;
- total al cierre del snapshot: `924`;
- duplicados: `0`;
- filas `OPEN`: `0`.

`LIB-FAIL-770`–`776` retienen dos truncamientos, test ausente, SCA rechazada, conteo API nulo, parche agrupado y ambigüedad de router. `UP-FAIL-148` conserva las brechas del source. Ningún PASS borra esos hechos.

## Estado honesto

La biblioteca permanece `NOT_READY_UNDER_EXPANDED_USER_STANDARD`. Existe código AWS exacto para cámara/PWA y object-event processing, pero esta revisión demostró que ese sample no puede usarse instantáneamente como ingreso confiable. Portal/API→cuarentena AWS sigue siendo el único ingreso remoto compuesto y ejecutable condicionado; email, SFTP y mobile productivos permanecen abiertos.

