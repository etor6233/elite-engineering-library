# All implementation packs — materialization evidence V17

Fecha: `2026-08-26`.

Este snapshot sucede a V16 y gobierna el estado fuente actual después de admitir de forma condicionada AWS GenAI IDP Accelerator 0.6.5 y corregir el flujo positivo de adquisición por perfil.

## Resultado global

`VERIFY_LIBRARY.ps1`:

```text
VERIFY_LIBRARY_PASS
packs=38 materialized_files=346 markdown_files=229
```

Los 17 perfiles de composición continúan reconstruyendo realmente sus conteos esperados: inicialización `23`, Azure document `11`, Google document `11`, AWS Textract `24`, AWS enterprise `26`, Google Ads `22`, Meta Ads `25`, TikTok Ads `25`, Meta WhatsApp `28`, Firebase `26`, Mercado Libre `10`, Amazon Catalog `24`, Google Merchant `25`, payments `7`, Business Central `6`, backend `95` y web/BFF `44`.

`VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Audit`:

```text
VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit packs=38 upstream_sources=67 document_sdk_artifacts=5 business_central_artifacts=1 provider_adapters=7
```

## Delta V17

- `OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.15` conserva `17` archivos y eleva el lock a `67` fuentes exactas;
- añade `aws-genai-idp-accelerator-0.6.5`, release/commit/archive/licencia/NOTICE/uv.lock fijados;
- `aws-secure-document-pipeline` selecciona `8` fuentes y `document-intelligence-leaders` selecciona `28`;
- los ocho perfiles pasan validación; suite: `5` negativos y `1` positivo offline;
- la adquisición real del candidato emitió source/profile receipts sin Git;
- la corrección `LIB-FAIL-110` elimina `$LASTEXITCODE` entre scripts PowerShell y queda cubierta por el positivo;
- el candidato conserva `UP-FAIL-042..044`: logs privados de Dependabot, findings OSV/npm y dos errores full-tree basedpyright. No fue promovido ni desplegado.

La evidencia técnica detallada está en `AWS_GENAI_IDP_ACCELERATOR_2026-08-26_V1.md`. El estado global continúa `NOT_READY_UNDER_EXPANDED_USER_STANDARD`: tener código oficial adquirible y miles de tests offline no reemplaza corrección upstream, cuentas/sandbox, corpus/ground truth, seguridad, carga, recuperación ni aceptación del proyecto.
