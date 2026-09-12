# All implementation packs materialization — 2026-08-28 V62

## Cambio gobernante

- El source lock pasó de 98 a 100 fuentes e incorporó código público oficial exacto de PaddlePaddle/Baidu: `PaddleOCR` v3.7.0 en commit firmado `b03f46425e8ff4442b268ce449e3eef758146cd4` y `PaddleX` v3.7.2 en commit no verificado `ffb64904d23708863ff5b8da312a5cbd52a7f462`.
- Ambos sources quedan `PINNED_CANDIDATE_CONDITIONED`, no `REUSABLE_PACK`: archives, licencias Apache-2.0, wheel/sdist, locks y pesos det/rec están fijados por bytes/SHA-256, pero el corpus empresarial, ground truth, exactitud por campo y autorización de persistencia siguen bloqueados.
- PaddleOCR ejecutó su comando oficial con 220 tests PASS, 5 deselected y 1 warning. El runtime exacto de 104 paquetes pasó sync hash-locked, imports, compatibilidad y SCA sin findings; el grafo ampliado de pruebas conserva siete registros de vulnerabilidad en cuatro paquetes LangChain opcionales.
- Una inferencia CPU local usó exclusivamente pesos PP-OCRv6 det/rec en revisiones Hugging Face verificadas: ocho archivos/139.157.714 bytes. El fixture oficial produjo doce textos/scores, pero no publica ground truth; por contrato, score de confianza no equivale a corrección ni habilita almacenamiento automático.
- La identidad completa, artefactos, modelos, locks, tests, SCA y condiciones están en `PADDLEPADDLE_PADDLEOCR_3_7_0_REAUDIT_2026-08-28_V1.md`.

## Reconstrucción focal

- `OFFICIAL-UPSTREAM-ACQUISITION-CORE` 0.4.53 materializó 23 archivos desde Markdown en un destino nuevo.
- Los cuatro bloques modificados coincidieron en longitud y SHA-256 con el árbol auditado:
  - `upstream-source-lock.json`: 180.577 bytes, SHA-256 `7289d0764d62af054e0d36c647b2539b98ab44aec24f436aed71ca5bf5ba8224`;
  - `test_upstream_acquisition.ps1`: 46.598 bytes, SHA-256 `5a06297a69bba6b6f57e799639a0c956d79c11c03ee09a638ffef97707daa720`;
  - `source-profiles/document-intelligence.json`: 12.289 bytes, SHA-256 `27b27a98d87f760dc2eea8924cfa860dfdc891ca068ba6456b1c1f979217db8f`;
  - `test_source_profiles.ps1`: 8.174 bytes, SHA-256 `f325ec13bdb8a191f19e79e4892c433e2c69333a32e69bc3f450761dd3817756`.
- El acquirer pasó sus 9 negativas; el contrato de perfiles pasó 14 perfiles válidos, 5 negativas y 1 positiva. El perfil documental selecciona 45 fuentes únicas y el lock valida 100/100.
- Las 18 referencias de composición consumidoras del pack quedaron alineadas en 0.4.53; cero referencias consumidoras permanecen en 0.4.52.

## Fallos y cierre fail-closed

- Las lecciones `LIB-FAIL-676` a `LIB-FAIL-694` preservan cada fallo de esta auditoría: locks demasiado verbosos, entorno de tests incompleto, comandos uv inválidos, búsquedas truncadas, acceso/modelos, parsing del resultado, paths asumidos, comparación PowerShell no terminante, dos defectos del gate documental y un cleanup compuesto rechazado preventivamente.
- Cada fila cerró sólo después de una prueba acotada. El ledger final contiene 694 fallos propios más 141 condiciones upstream: 835 IDs esperados y cero filas abiertas.
- El gate documental final comprobó pack 0.4.53, 18 planes consumidores, 100 sources únicos, perfil documental de 45 sources únicos, catálogo y ledger de admisión con PaddleOCR/PaddleX condicionados.

## Gate global

Antes de añadir esta evidencia, `VERIFY_LIBRARY.ps1` terminó con `VERIFY_LIBRARY_PASS` sobre 50 packs, 482 archivos materializables, 346 Markdown y 27 perfiles de composición. La repetición posterior debe gobernar el conteo final de 347 Markdown; backend/web permanecen en 95/44 archivos.

Este PASS demuestra consistencia y reconstrucción de la biblioteca. No promueve PaddleOCR/PaddleX a implementación productiva ni demuestra exactitud de facturas, proformas, packing lists u otras clases reales: esas afirmaciones requieren corpus autorizado, ground truth, evaluación por campo, revisión/rechazo, seguridad, carga, recovery y autorización de persistencia del proyecto.
