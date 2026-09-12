# Microsoft Azure Content Understanding .NET GA — 2026-08-26 V1

## 1. Resultado de admisión

Se admitió como `PINNED_CANDIDATE` el código fuente oficial de `Azure.AI.ContentUnderstanding 1.1.0`, SDK GA de Microsoft para la API estable `2025-11-01`. No es glue de Elite y no se declaró deployable: el código requiere Microsoft Foundry, identidad, región, cuotas, model deployments, documentos autorizados y consumo potencialmente pago.

## 2. Fuente y firma exactas

- repositorio oficial: `https://github.com/Azure/azure-sdk-for-net`;
- tag oficial: `Azure.AI.ContentUnderstanding_1.1.0`;
- commit exacto: `6af7db6572976f4b6cc65bb628c45456114e88a0`;
- fecha: `2026-04-22T00:44:47Z`;
- verificación GitHub del commit: `verified=true`, `reason=valid`;
- archive: `https://github.com/Azure/azure-sdk-for-net/archive/6af7db6572976f4b6cc65bb628c45456114e88a0.zip`;
- bytes archive: `251,185,929`;
- SHA-256 archive: `a88e5c78ae1c881d565141867decaf0dba8a0ef63cf02a16833d0ba51bc16860`;
- licencia raíz `LICENSE.txt`: `1,076` bytes, SHA-256 `9b45236978bb5cd5de992021769e1eeaa79f0116d3b02cf8ba065a1ed603d5fa` (`MIT`).

El paquete oficial NuGet `Azure.AI.ContentUnderstanding 1.1.0` se descargó desde `api.nuget.org`:

- bytes: `678,712`;
- SHA-256: `c026c61c1464d96723e55ff91dd46deae45f741a83d4a0ce68f7cc212d5aea37`;
- nuspec: Microsoft, MIT, repository commit `6af7db6572976f4b6cc65bb628c45456114e88a0`, dependencia directa `Azure.Core 1.53.0`;
- firma de autor válida: `Microsoft Corporation`, certificate SHA-256 `566A31882BE208BE4422F7CFD66ED09F5D4524A5994F50CCC8B05EC0528C1353`;
- countersignature de repositorio válida: `NuGet.org Repository by Microsoft`, SHA-256 `1F4B311D9ACC115C8DC8018B5A49E00FCE6DA8E2855F9F014CA6F34570BC482D`.

La verificación se ejecutó con .NET SDK oficial `8.0.424`; su ZIP Windows x64 de `285,090,820` bytes coincidió con el SHA-512 publicado por Microsoft `1787ab90635c2950672ed7c6507b000e1b212ea7d9a22fcef37061344d37c64d4c4eda12b8742601eff5b45c8736485b31c55613892f240c300190e4e88a58b0`.

## 3. Código oficial relevante

El sample exacto `sdk/contentunderstanding/Azure.AI.ContentUnderstanding/samples/Sample03_AnalyzeInvoice.md` tiene `18,991` bytes y SHA-256 `19eef824a33a8023559f85f26057ecd74bcda6461d674c350dd36f7d06084bdf`.

Microsoft declara que el analyzer `prebuilt-invoice` está optimizado para:

- invoices;
- utility bills;
- sales orders;
- purchase orders.

El sample oficial usa el cliente fuertemente tipado, long-running operation `AnalyzeAsync`, `prebuilt-invoice`, campos con confidence, sources, page, polygon/bounding box, spans, fechas, importes, moneda y line items. Es código reutilizable del fabricante para esa frontera; no demuestra por sí solo proformas, packing lists, documentos aduaneros ni los schemas particulares del proyecto.

## 4. Grafo restaurado y seguridad

Un proyecto de auditoría mínimo fijó el paquete `1.1.0`, generó `packages.lock.json` y repitió `restore --locked-mode`. El lock tiene SHA-256 `c590b41bcb9a397ee58560d6cdb24969e127a52287a8dc01f2e23fd3c4ec4c80`.

El grafo `net8.0` resolvió `22` paquetes: un top-level y `21` transitivos. Una consulta fechada a la API oficial de OSV para el ecosistema NuGet devolvió `0` paquetes afectados y `0` findings conocidos. Esto es evidencia temporal, no garantía de ausencia de vulnerabilidades futuras.

## 5. Integración y límite

`OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.17` eleva el lock a `68` fuentes. El perfil Microsoft selecciona `12` y el perfil documental integral `29`. Ambos incluyen `microsoft-azure-content-understanding-dotnet-1.1.0`; el perfil Microsoft registra explícitamente cuentas, costo, modelo y pruebas de proyecto como blockers.

La biblioteca todavía no contiene un runtime productivo .NET para documentos ni credenciales. La incorporación actual permite adquirir el código oficial exacto al instante; ejecutar el analyzer requiere completar el intake y demostrarlo contra corpus/ground truth del negocio.
