# All Implementation Packs Materialization — V79

Fecha: 2026-08-28  
Estado: `VERIFY_LIBRARY_PASS` + `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit`

## Resultado gobernante

- 58 implementation packs y 542 archivos materializables desde Markdown;
- 389 archivos Markdown contando esta evidencia;
- 32 perfiles de composición; Azure document compone 9 packs/79 archivos;
- 110 fuentes oficiales exactas y 14 perfiles source permanecen fijados por `OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.63`;
- ledger: 879 fallos locales + 160 condiciones upstream = 1.039 IDs únicos, cero duplicados y cero lecciones locales abiertas.

## Código oficial Microsoft agregado

`MICROSOFT-AZURE-CONTENT-UNDERSTANDING-OFFICIAL-ANALYZER-COPY-SAMPLE 0.1.0` materializa seis archivos. Tres son bytes `VERBATIM` Microsoft MIT del commit firmado `Azure/azure-sdk-for-python@129a5cbb06fba44e11ffc92f6bcdca14a2f5b6eb`:

- `sample_copy_analyzer.py`: 6.071 bytes/SHA-256 `4b8aa70c941c9b0b8fcd54ddac588ef0d1a6dbd4cefe5e2c21d11fd8e44bdb71`;
- `test_sample_copy_analyzer.py`: 9.468 bytes/SHA-256 `227d59a71c0e1ecaa67a9eb406654c79bc6a1294334a1d1aca6eceeb58c18f8e`;
- `LICENSE.txt`: 1.074 bytes/SHA-256 `7c77a44a8acd9b41fdc209864a8016b3d430b5d0e09309818d5b7444336df744`.

Sample y test son byte-idénticos al sdist oficial `azure-ai-contentunderstanding 1.1.0` ya fijado. Los otros tres archivos son packaging, lock y contratos `AUTHORED`; no se atribuyen a Microsoft.

## Lifecycle demostrado

El sample oficial ejecuta create → get → same-resource copy → get → update → get → delete mediante `ContentUnderstandingClient`. Crea un analyzer custom basado en `prebuilt-document`, habilita layout/OCR y estimación de source/confidence, define dos campos de demostración y actualiza el tag del target.

El test oficial usa Azure SDK test-proxy y un recurso Azure; no es standalone. La biblioteca lo preserva completo y no lo presenta como suite offline. Los contratos locales sólo comprueban identidad, sintaxis, llamadas y non-claims.

## Gates

- materialización y hashes: 6/6 PASS;
- contratos offline del pack: 3/3 PASS;
- sintaxis: sample, test upstream y contrato, 3/3 PASS con bytecode fuera del source;
- composición Azure: 9 packs/79 archivos PASS;
- `VERIFY_LIBRARY_PASS`: 58 packs/542 archivos/388 Markdown antes de esta evidencia;
- `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit`: `official_invoice_samples=1`, `official_azure_cu_samples=1`, `official_azure_cu_copy_samples=1`, `official_google_process_samples=1` y todos los demás componentes focales verdes.

## Condiciones y non-claims

Este código muta recursos externos y no se autoejecuta. Requiere cuenta/región/endpoint/identidad/modelos/costo, autoridad explícita de create/copy/update/delete, IDs inmutables, receipt de versión, inventario posterior y reconciliación. El sample absorbe errores de cleanup; no se interpreta un delete intentado como rollback demostrado.

No usa `prebuilt-procurement` ni `prebuilt-purchaseOrder` y no demuestra proforma, packing list, bill of lading, aduana, corpus, exactitud o promoción productiva. Microsoft documenta esos analyzers como capacidades del servicio; su mera existencia no sustituye código/sample específico ni evaluación por clase/campo.

## Recuperación

`LIB-FAIL-879` conserva el intento fallido de comparar arrays con `.AsSpan()` en PowerShell. La repetición tipó ambos `byte[]` y usó `Enumerable.SequenceEqual[byte]`, demostrando dos filas byte-idénticas entre commit y sdist antes de construir el pack.
