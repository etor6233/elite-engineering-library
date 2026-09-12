# All Implementation Packs Materialization — V77

Fecha: 2026-08-28  
Estado: `VERIFY_LIBRARY_PASS` + `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit`

## Resultado gobernante

- 56 implementation packs y 531 archivos materializables desde Markdown;
- 382 archivos Markdown después de incorporar esta evidencia;
- 32 perfiles de composición; `AWS_TEXTRACT_DOCUMENT_RUNTIME_PACK_PLAN.md` compone 7 packs/75 archivos;
- source lock global permanece en 110 fuentes oficiales/14 perfiles: Textractor 1.10.0 ya estaba adquirido y fijado, por lo que no se inventó una fuente adicional;
- ledger: 861 fallos locales + 158 condiciones upstream = 1.019 IDs únicos esperados, cero duplicados y cero lecciones locales pendientes.

## Código oficial incorporado

Fuente: `aws-samples/amazon-textract-textractor`, release `v1.10.0`, commit `8ea5f9ae65dbcbb2b83b409e0007c395e59b0e92`, Apache-2.0 con NOTICE. Archive: 79.080.629 bytes/SHA-256 `4fa38999e355d977d2b46816b76c9d32ab6c5d8607987e086045470d101c7219`; inventario expandido: 382 archivos/112.777.769 bytes. Wheel oficial PyPI `amazon_textract_textractor-1.10.0-py3-none-any.whl`: 311.287 bytes/SHA-256 `8524224f07a776ca2959e0d80d5e031a30bf53938c74e8352848f3d0b6b3762c`.

`PYTHON_AWS_TEXTRACTOR_OFFICIAL_COMPONENT.md` materializa 12 archivos: cinco de packaging/locks/contratos Elite `AUTHORED` y siete archivos oficiales AWS —LICENSE, NOTICE, el test de Queries y cuatro archivos auxiliares/fixtures. Tres son `VERBATIM`; cuatro son `ADAPTED` sólo para representar un LF final bajo el contrato de fences. `source-lock.json` conserva para cada adaptación los bytes y SHA-256 originales, el cambio exacto y el SHA-256 empaquetado. No se atribuye a AWS ningún wrapper, schema empresarial ni permiso de persistencia.

El componente se compone después de `GO_AWS_TEXTRACT_DOCUMENT_RUNTIME.md`: el runtime Go gobierna el llamado y los receipts; Textractor aporta el parser/objetos y pruebas oficiales reutilizables. Esa composición no convierte respuestas OCR en verdad empresarial ni autoriza almacenamiento automático.

## Ejecución y seguridad de dependencias

- reconstrucción limpia desde Markdown: 12/12 archivos y hashes PASS;
- contratos offline del componente: 4/4 PASS;
- prueba oficial Queries materializada: 2 PASS/1 skip; el skip es la llamada live que requiere AWS;
- suite core determinista source-exact en Linux/Python 3.14: 69 PASS/16 skips, excluyendo únicamente el fuzz que upstream describe como aleatorio/flaky pero no marca como skip;
- Windows/Python 3.14: 67 PASS/16 skips/2 FAIL; los dos nodeids fallidos pasan 2/2 sobre Linux sin modificar el source;
- runtime lock: 15 paquetes exactos, `pip check` limpio y OSV-Scanner 2.5.1 con cero vulnerabilidades conocidas; reporte 3.130 bytes/SHA-256 `3083672b5fcd9f1d529fa9256f75adfc6ba1a9802d8e4049899d73d2860b7787`;
- test graph: 26 paquetes exactos y cero vulnerabilidades conocidas; reporte 5.257 bytes/SHA-256 `1b63bfab44bbad8970d06b52267aa65c4bf36ce1e2e6ca7af1d95e7dbd58475d`.
- auditor global: `aws_textractor_official_components=1`; los 56 packs, 110 fuentes y componentes offline admitidos terminaron `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit`.

## Condiciones upstream conservadas

- 155: los snippets vigentes AWS para intermediary dispatcher expresan la arquitectura, pero contienen variables/campos inconsistentes y no son una implementación copiable verificada; se rechazaron;
- 156: `Page.__repr__()` difiere CRLF/LF en Windows; el nodeid pasa en Linux;
- 157: el test de orden de palabras en una celda falla en Windows y pasa en Linux; por eso no se promete orden de lectura exacto de cualquier documento;
- 158: el fuzz `test_parse_no_fail.py` se autodescribe como deshabilitado/flaky, pero el source/workflow lo recolecta sin marker; se excluyó explícitamente del gate determinista y no se contó como PASS.

Las cuatro condiciones permanecen `UPSTREAM_OPEN`: una reproducción Linux distingue portabilidad, pero no corrige ni reescribe el upstream.

## Recuperación de fallos

- 851 corrigió un nombre de pack Textract supuesto mediante inventario literal;
- 852 rechazó pytest 9.0.2 y reconstruyó el grafo con pytest 9.0.3 antes de ejecutar;
- 853 limitó pytest al scope core oficial y evitó atribuir imports de subproyectos independientes;
- 854 identificó y separó el fuzz no acotado antes de repetir la suite determinista;
- 855 reutilizó la lección de WSL sin pip y ejecutó desde wheels exactos expandidos, sin `apt`, `sudo` ni instalación global;
- 856 derivó el perfil AWS desde el compositor y corrigió 6/58 obsoleto a 7/75 real;
- 857 y 858 corrigieron la firma del materializador y el enum de provenance sin escribir un destino parcial;
- 859 corrigió la invocación de `unittest` y demostró cuatro contratos reales;
- 860 cerró las filas probadas antes de repetir el gate estructural;
- 861 eliminó la suposición de un ledger upstream separado y reafirmó el ledger canónico único.

## Límites vigentes

Este componente no demuestra extracción perfecta, exactitud de negocio, almacenamiento automático, schemas para todas las clases, revisión humana completa ni operación AWS productiva. Antes de persistir un campo siguen siendo obligatorios: clase y schema cerrados, corpus/ground truth representativo, evidencia por campo, comparación exacta, política de confianza/rechazo, revisión cuando el proyecto la exige, binding a original/tenant/version, seguridad de archivo, cuenta/región/IAM/costo, datos live, carga, recovery y rollback. También continúan abiertos el dispatcher durable empresarial, orden/fencing de eventos S3, DLQ/redrive y reconciliación del efecto final.
