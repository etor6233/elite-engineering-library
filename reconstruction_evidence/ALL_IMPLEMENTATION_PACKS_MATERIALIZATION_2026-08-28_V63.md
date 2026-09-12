# All implementation packs materialization — 2026-08-28 V63

## Cambio gobernante

- Se añadió `PADDLEPADDLE-PADDLEOCR-LOCAL-RUNTIME` 0.1.0: siete archivos materializables para ejecutar los paquetes oficiales PaddleOCR 3.7.0, PaddleX 3.7.2 y PaddlePaddle 3.1.0 con pesos PP-OCRv6 medium locales exactos.
- El código Paddle continúa siendo el wheel/source oficial Apache-2.0. El lock, modelo, perfil, wrapper, verificador, tests y README son `AUTHORED`; no se atribuyen a PaddlePaddle/Baidu.
- El template nace `enabled=false`, cero formatos/clases aprobados, modelos sólo locales y `automatic_business_storage=false`. El runner exige receipt seguro de los mismos bytes/MIME, clase aprobada y ocho archivos de modelo por longitud/SHA antes de construir el predictor.
- `PADDLEOCR_LOCAL_RUNTIME_PACK_PLAN.md` compone adquisición + ingesta segura + runtime + evaluación estricta: 4 packs/41 archivos. El compositor y el gate global lo materializan realmente.

## Runtime reproducido

- uv 0.12.7 oficial: ZIP Windows x86-64 de 16.979.508 bytes; el asset `.sha256` oficial y el ZIP local coincidieron en `bf1518af459a3915511a11fdc6e2f43ef9a2afa138b9d498eeb9642fe9d85218`.
- Lock portable: 104 paquetes únicos, 239.199 bytes, SHA-256 `50f83b052423bea9daf0b331480c0e642440251afdf99bfc04931af7fd5a5a9b`.
- CPython 3.12.14 Windows x86-64 limpio: sync `--require-hashes`, `pip check`, imports/versiones exactos y CPU mode PASS; sello `PADDLEOCR_ENVIRONMENT_PASS`.
- `pip-audit==2.10.1`: 104 dependencias, cero vulnerabilidades conocidas en la consulta fechada; el reporte auditado previo al cambio neutral de header tuvo 5.892 bytes y el contenido de paquetes no cambió.
- `py_compile` y 8/8 tests del wrapper PASS: perfil deshabilitado, storage, clase, receipt hash, no-result, output existente, plataforma y success/atomicidad.

## Reconstrucción y fallos

- Los siete archivos reconstruidos coincidieron con staging en longitud/SHA. Hashes focales: model lock `9a30661c…fd5d`, profile `ff2e18ca…dc91`, runner `6364d5d3…73b2`, verifier `343f3751…d939`, tests `b7db13c6…7175`, README `115e6efa…f59c`.
- `LIB-FAIL-695` a `LIB-FAIL-702` preservan paths históricos supuestos, ubicación del contrato, `foreach` no agrupado, checksum uv inventado desde abreviación, path privado del lock/ledger, dependencia circular de filas abiertas y el perfil omitido de la lista ejecutable.
- El ledger vigente contiene 702 fallos propios más 141 condiciones upstream: 843 IDs esperados y cero filas abiertas.

## Gate global

Antes de añadir esta evidencia, `VERIFY_LIBRARY.ps1` terminó `VERIFY_LIBRARY_PASS` sobre 51 packs, 489 archivos materializables, 349 Markdown y 28 perfiles de composición. Mostró explícitamente `profile=PADDLEOCR_LOCAL_RUNTIME_PACK_PLAN.md implementation_files=41`. La repetición posterior debe gobernar el conteo final de 350 Markdown.

El pack aporta ejecución inmediata y fail-closed del motor oficial, pero permanece `CONDITIONED`: la inferencia real previa usó un fixture oficial sin ground truth; los tests nuevos ejercitan el contrato del wrapper con predictor controlado. Ninguna de ambas pruebas demuestra exactitud de facturas, proformas, packing lists u otros documentos del negocio ni habilita persistencia automática.
