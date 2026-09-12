# Microsoft Azure Content Understanding official invoice sample — 2026-08-28 V1

## Resultado

Se materializó un pack nuevo con el sample Python síncrono oficial Microsoft `sample_analyze_invoice.py` de `azure-ai-contentunderstanding 1.1.0`. El código productivo del sample y la licencia son upstream; README, lock y contrato offline son packaging Elite explícitamente `AUTHORED`.

## Identidad oficial

- repositorio: `Azure/azure-sdk-for-python`;
- release: `azure-ai-contentunderstanding_1.1.0`;
- commit firmado: `129a5cbb06fba44e11ffc92f6bcdca14a2f5b6eb`;
- sdist: 230.330 bytes, SHA-256 `00f39c7cf2ba50c8d4586315f4295a45fad0881daaff07d13d33aea2bafc7ec8`;
- wheel: 101.987 bytes, SHA-256 `d1d6bdeffe02f5c8cc5309f0e120a1338f51af04638605f9617b7b501e2d1dd6`;
- sample: 10.539 bytes, SHA-256 `0cb9d7b0e183cd1c3f677c79a9d8d5dd9a1ed30bf397f84b753eedd7819c53fd`;
- el sample del raw commit y el del sdist son byte-idénticos;
- licencia raíz del commit: 1.074 bytes, SHA-256 `7c77a44a8acd9b41fdc209864a8016b3d430b5d0e09309818d5b7444336df744`.

El sdist incluye una variante de la licencia de 1.073 bytes/SHA-256 `fd532481d828e13a0b13ccb598e02338a3617740675a862ee6bdc1541b68e93d` que sólo omite el LF final. El pack usa el archivo raíz exacto del commit, por lo que permanece `VERBATIM`; ambos hashes quedan en `source-lock.json`.

## Código y claim estrecho

El source Microsoft crea `ContentUnderstandingClient`, llama `begin_analyze` con el literal `prebuilt-invoice` y muestra:

- `CustomerName` e `InvoiceDate`;
- `TotalAmount` y moneda;
- `LineItems`, descripción y cantidad;
- confidence, source, spans y usage/costo.

Microsoft documenta `prebuilt-invoice` para invoices, utility bills, sales orders y purchase orders. Este sample no llama `prebuilt-procurement` ni `prebuilt-purchaseOrder`; tampoco cubre proformas, packing lists, bills of lading, customs ni schemas privados.

## Reconstrucción

Destino nuevo: `C:\Temp\elite-v80-azure-cu-official-pack-materialized-v2`.

```text
Materialized 5 files
5/5 source-tree hashes equal
Ran 3 tests in 0.007s — OK
compileall — PASS
```

Los contratos verifican bytes/licencia/commit, AST con un único analyzer `prebuilt-invoice` y una ejecución offline del `main` contra un fake SDK que conserva endpoint, input y una sola llamada. No simulan calidad del modelo ni una llamada real.

## Condiciones conservadas

- el sample usa una URL demo móvil bajo `main`; no se usa para documentos del proyecto;
- `dev_requirements.txt` del commit, 128 bytes/SHA-256 `c00259b708ba76145d2e545fc3b7e915443ececaaa190d295dac2753e6e42215`, enumera `python-dotenv` y `azure-identity` sin versiones; el target debe resolver/pinnear/escanear su grafo;
- recurso Foundry, región, modelos, identidad, cuotas y costo siguen externos;
- corpus, ground truth, evaluación por campo, review, reconciliación y autorización de persistencia siguen bloqueantes.

No se realizó llamada Azure ni se atribuye precisión productiva.

