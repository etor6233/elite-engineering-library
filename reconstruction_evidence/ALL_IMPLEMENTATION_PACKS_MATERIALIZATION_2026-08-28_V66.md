# All implementation packs materialization — V66

Fecha: 2026-08-28.

## Cambio gobernante

`OFFICIAL-UPSTREAM-ACQUISITION-CORE` avanzó a 0.4.54 y el lock pasó de 100 a 101 fuentes oficiales. La nueva revisión exacta es `aws-samples/sample-amazon-ses-mail-manager-attachment-pipeline` commit unsigned `79314a93bda431bd03a4c1236bc4b22e75e2f41f`, MIT-0 + NOTICE, archive 709.657 bytes/SHA-256 `e1e241c4b8a0765f913bff04cf1eb7b715dfbeaace580b398c181eda5cd30fd5`.

La fuente queda `SAMPLE_ONLY_REJECTED_FOR_IMMEDIATE_ADOPTION`, no compuesta ni desplegada. Once Python parsean, pero no hay tests/lock; el pin oficial `aws-cdk-lib 2.251.0` tiene `GHSA-464c-974j-9xm6` y el npm CLI 2.251.0 solicitado no existe. Una evaluación separada de library 2.253.0 + CLI 2.1139.0 sintetizó 32 recursos y un grafo 13/0 OSV, pero no cierra brechas de tenant, key uniqueness, límites/checksum/idempotencia, KMS, retención ni failure policy.

## Reconstrucción del pack de adquisición

- 23/23 archivos desde Markdown;
- pack SHA-256 `700d2c26127953762c5ebfc999fd19be9acda277c394dad7d23ccb6309a0fd09`;
- lock 101/101 fuentes válidas y únicas;
- nueve negativos de lock/acquirer PASS;
- catorce perfiles válidos, cinco negativos y una adquisición positiva offline PASS;
- perfil AWS: 16 fuentes;
- veinte planes consumidores alineados a 0.4.54; cero consumidor vigente en 0.4.53.

## Gates globales

`VERIFY_LIBRARY.ps1`:

```text
VERIFY_LIBRARY_PASS
packs=53 materialized_files=499 markdown_files=357
```

Antes de crear este snapshot, `VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Audit` terminó:

```text
VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit packs=53 upstream_sources=101 document_sdk_artifacts=10 official_invoice_samples=1 official_google_process_samples=1 official_dapr_outboxes=1 official_pg_durable_handoffs=1 business_central_artifacts=1 provider_adapters=9 document_orchestrators=1 evidence_logs=1
```

Los `SKIPPED` por cuentas, autorización, red o toolchain permanecen condiciones, no evidencia productiva transferida. La recepción email productiva continúa abierta: este cambio evita repetir la búsqueda del sample y evita que un agente lo adopte por error; no afirma haber terminado el canal.
