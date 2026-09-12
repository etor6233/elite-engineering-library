# NVIDIA NeMo Retriever current signed commit — 2026-08-26 V1

## 1. Resultado honesto

Se reemplazó la identidad bloqueada de la release `26.5.0` por el commit exacto y firmado `6a05f67e4934264eebb346b153111944593af80a` de la rama oficial `NVIDIA/NeMo-Retriever`. No se lo promovió a código productivo reusable: la release pública continúa siendo `26.5.0`, el commit fijado es no publicado y su perfil mínimo de servicio conserva dos dependencias con advisories conocidos.

Estado admisible: `PINNED_CANDIDATE / CONDITIONAL_PLATFORM / UPSTREAM_OPEN`.

## 2. Autoridad e identidad exacta

- repositorio oficial: `https://github.com/NVIDIA/NeMo-Retriever`;
- commit: `6a05f67e4934264eebb346b153111944593af80a`;
- fecha GitHub del commit: `2026-08-26T21:09:28Z`;
- verificación GitHub: `verified=true`, `reason=valid`;
- release oficial más reciente observada: `26.5.0`;
- archive exacto: `https://github.com/NVIDIA/NeMo-Retriever/archive/6a05f67e4934264eebb346b153111944593af80a.zip`;
- bytes archive: `10,263,665`;
- SHA-256 archive: `577ebb21469d61c4c9031ac98a32daf0d09993fa463bc98d7eabeceee23b8cc1`;
- licencia `Apache-2.0`, `LICENSE` de `11,357` bytes, SHA-256 `c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4`;
- notice `THIRD_PARTY_LICENSES.md` de `548` bytes, SHA-256 `5b89b1db6e9a8527eb16dfa54721fe01285483a5c6843deea3484dcde4e290fa`;
- lock `nemo_retriever/uv.lock` de `575,991` bytes, SHA-256 `eab5b6f5ad18d2dd90f7969c184349d5fb4c46b1f07c8123b88d7d7a31a5e1a5`;
- `nemo_retriever/pyproject.toml` de `10,228` bytes, SHA-256 `40cdb335752ce54a6a1137045b1ec18d893029e1caa0b0652ae169cddd5bd00b`;
- `.python-version` fija `3.12`, SHA-256 `7b55f8e67b5623c4bef3fa691288da9437d79d3aba156de48d481db32ac7d16d`.

La consulta oficial de GitHub devolvió siete checks para ese commit, todos terminados y exitosos: dos `build`, `Build & Test Docker (amd64)`, `deploy`, `Pre-commit Checks`, `report-build-status` y `Run Retriever Unit Tests`.

## 3. Comparación contra la release anterior

La release `26.5.0`, commit `0a6cb709b1ee4ca1d3f8c2fe0bbd1a0247f0aa09`, quedó descartada como selección vigente. Su archive exacto tenía `10,209,364` bytes y SHA-256 `b1c6856c83554cb5afda0b5bbe7c65d6f6866e40e365ee83c2dffc3c0c65d43d`. Su `uv.lock` fijaba `pillow==12.2.0`.

OSV-Scanner `2.5.1` sobre la exportación exacta del perfil `service` de esa release encontró `17` paquetes afectados y `130` IDs/aliases. Sobre todo el lock encontró `25` paquetes afectados y `173` ocurrencias (`172` IDs únicos). Esto impide reutilizar `26.5.0` como base productiva aunque sea la última release publicada.

## 4. Prueba independiente del commit firmado

El archive fijado contenía `1,303` entradas y `1,093` archivos, sin traversal, `202` archivos Python de tests, `25` workflows y un único `uv.lock`.

Con Python `3.12.14` y uv `0.12.6` se exportó primero el grafo exacto:

```text
uv export --frozen --no-default-groups --extra service --extra dev --no-emit-project --no-hashes
```

La exportación produjo `197` paquetes declarados y no contenía `torch`, `vllm`, `triton`, `nvidia-*`, `nemotron-*`, `pynvvideocodec` ni `tilelang`. Luego se ejecutó desde un extracto vacío:

```text
uv sync --frozen --no-default-groups --extra service --extra dev
uv pip check
python -m compileall -q src
PYTHONPATH=src python -m pytest tests/test_slim_imports_no_triton.py -q -p no:cacheprovider
```

Resultado: `189` distribuciones instaladas, `pip check PASS`, compilación `PASS` y `1 passed`; se observó una advertencia de deprecación Starlette/httpx. El freeze tiene SHA-256 `fc80af0af88a618884a2eb3d8127471403f57bb5f9d801daff7df061d101359b`.

El primer intento activó por error el grupo uv `dev` predeterminado. Upstream define ese grupo como `nemo_retriever[all]`, por lo que comenzó a resolver CUDA/vLLM. Fue interrumpido, descartado y repetido desde cero con `--no-default-groups`; queda registrado como `LIB-FAIL-112`.

## 5. Seguridad vigente y bloqueo

El `uv.lock` actual cambia Pillow a `12.3.0` y reduce materialmente los hallazgos: todo el lock contiene `7` paquetes afectados y `11` aliases (`10` IDs únicos). La exportación exacta `service` todavía contiene `2` paquetes afectados y `4` aliases:

| Paquete fijado | Condición OSV observada | Primera versión corregida informada | Decisión |
|---|---|---|---|
| `cryptography 49.0.0` | `CVE-2026-69247`, oracle Bleichenbacher en descifrado PKCS7 EnvelopedData | `50.0.0` | bloquear promoción |
| `setuptools 80.10.2` | `CVE-2026-59890`, colisión Unicode de `MANIFEST` al construir sdist en macOS | `83.0.0` | bloquear promoción y evaluar reachability por plataforma |

La exportación `service` auditada tiene SHA-256 `915a481c29a725ececea4b4898dc12d1c40ee6ed03ae29ee063df8349a4a45f8`. El reporte OSV JSON tiene SHA-256 `e8a87b0f5c3eade7b0609325137395f39c09222e6432a6d1513c8954d2f490d1`.

No se alteró el lock oficial, no se inventó un backport y no se interpretaron checks verdes como ausencia de vulnerabilidades. Para promoción se exige una release o commit oficial posterior, lock exacto sin blockers alcanzables, repetición de tests, GPU/NVIDIA services reales cuando aplique, corpus autorizado y métricas del proyecto.

## 6. Integración en Elite

`OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.16` conserva `67` sources pero reemplaza el source ID histórico por `nvidia-nemo-retriever-main-2026-08-26`. El lock incorpora archive, licencia, notice y uv lock exactos. El perfil `document-intelligence-leaders` conserva `28` selecciones y apunta al nuevo ID.

Los tests reconstruidos obtuvieron:

```text
UPSTREAM_ACQUISITION_TEST_PASS negatives=5
SOURCE_PROFILE_TEST_PASS valid=8 negatives=5 positives=1
```

Los tres fixtures NVIDIA ya catalogados (`table_test.pdf`, `embedded_table.pdf`, `test-page-form.pdf`) conservan bytes y SHA-256 idénticos en el commit fijado. Siguen siendo fixtures upstream, no evidencia de exactitud productiva para documentos de un negocio.
