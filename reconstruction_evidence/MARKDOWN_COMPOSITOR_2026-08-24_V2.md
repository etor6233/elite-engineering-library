# Markdown Compositor — Reconstruction Evidence V2

```yaml
evidence_id: "MD-COMPOSITOR-20260824-V2"
pack_id: "MARKDOWN-COMPOSITOR"
pack_version: "0.2.0"
pack_sha256: "e0a097c9b78b844b6f0d872fbd5b05d4cb98da519553edfd51d8128adbc0ddcb"
materialized_file_count: 3
ordered_file_hash_aggregate_sha256: "eb404d966dfb8837a36d35e307d24115338f857e38e726cd6c0d4a5c360b48d3"
toolchain: "PowerShell 7.6.4"
result: "PASS_CONDITIONED"
verified_at: "2026-08-24"
```

V2 corrige el parser de metadata para packs con LF o CRLF y añade una regresión específica. Success, acknowledgement de condiciones, SHA corrupto, colisión y CRLF: PASS.

El gate que faltaba en V1 también pasó: cuatro packs reales condicionados produjeron una composición de 35 archivos más `MATERIALIZATION_RECORD.md`; sus gates integrados están en `FOUNDATION_COMPOSITION_2026-08-24_V1.md`.

Límite: la versión 0.2 sólo admite `CREATE`; cualquier `PATCH`, `DELETE` o merge requiere una implementación explícita y falla cerrado.

