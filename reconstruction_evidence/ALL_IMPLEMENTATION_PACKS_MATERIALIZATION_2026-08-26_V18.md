# All implementation packs materialization — 2026-08-26 V18

## 1. Snapshot gobernante

Este snapshot sucede a V17 y conserva `38` packs con `346` archivos materializables. El cambio no añade glue local ni aumenta el número de fuentes: corrige la selección NVIDIA dentro del lock oficial y deja de dirigir al agente a la release pública vulnerable `26.5.0`.

`OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.16` fija `67` sources. `nvidia-nemo-retriever-main-2026-08-26` apunta al commit oficial firmado `6a05f67e4934264eebb346b153111944593af80a`, con archive, licencia Apache-2.0, `THIRD_PARTY_LICENSES.md` y `uv.lock` exactos. La comparación ejecutable, hashes, checks GitHub y findings OSV están en `NVIDIA_NEMO_RETRIEVER_CURRENT_2026-08-26_V1.md`.

## 2. Resultado verificable

Desde el Markdown canónico se materializaron los 17 archivos del pack de adquisición y se ejecutaron:

```text
UPSTREAM_ACQUISITION_TEST_PASS negatives=5
SOURCE_PROFILE_TEST_PASS valid=8 negatives=5 positives=1
```

El perfil documental seleccionó `28` fuentes contra el lock de `67`. El commit NVIDIA mínimo se probó aparte con Python 3.12.14 y uv 0.12.6: grafo sin GPU, 189 distribuciones, `pip check PASS`, `compileall PASS`, `1 passed` en el test oficial slim sin Triton. Dos advisories de dependencias mantienen la promoción bloqueada.

## 3. Verificadores globales

```text
VERIFY_LIBRARY_PASS
packs=38 materialized_files=346 markdown_files=231

VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit packs=38 upstream_sources=67 document_sdk_artifacts=5 business_central_artifacts=1 provider_adapters=7
```

Los 17 perfiles compuestos conservaron exactamente `23/11/11/24/26/22/25/25/28/26/10/24/25/7/6/95/44` archivos de implementación.

## 4. Estado y límite

La evidencia amplía la actualidad y reduce riesgo, pero no cambia la conclusión global: `NOT_READY_UNDER_EXPANDED_USER_STANDARD`. NeMo no es `REUSABLE_PACK`; es `PINNED_CANDIDATE / CONDITIONAL_PLATFORM / UPSTREAM_OPEN`. El agente debe esperar una corrección oficial verificable o una nueva release, reauditar el lock y probar el entorno real del proyecto antes de usarlo en producción.

Los fallos propios retenidos alcanzan `118` y las condiciones upstream `46`. Ninguno fue ocultado; la selección accidental del grupo GPU, el rechazo web, el parámetro incorrecto, los globs Windows, la limpieza rechazada, la recurrencia PowerShell y la expansión cruzada WSL fueron incorporados al ledger con su corrección o bloqueo.
