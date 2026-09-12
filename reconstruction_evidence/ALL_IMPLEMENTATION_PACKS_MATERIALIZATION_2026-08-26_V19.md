# All implementation packs materialization — 2026-08-26 V19

## 1. Snapshot gobernante

Este snapshot sucede a V18. Conserva `38` packs y `346` archivos materializables, y eleva el source lock de `67` a `68` fuentes oficiales exactas mediante Microsoft Azure Content Understanding .NET GA `1.1.0`.

`OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.17` fija el tag `Azure.AI.ContentUnderstanding_1.1.0`, commit firmado `6af7db6572976f4b6cc65bb628c45456114e88a0`, archive SHA-256, licencia MIT y source sample oficial de invoice. El paquete NuGet exacto fue verificado con firmas de Microsoft y NuGet.org; su grafo restaurado `net8.0` tiene 22 paquetes y cero findings conocidos en la consulta OSV fechada. La evidencia detallada es `MICROSOFT_AZURE_CONTENT_UNDERSTANDING_DOTNET_2026-08-26_V1.md`.

## 2. Perfil documental actualizado

El pipeline Microsoft selecciona `12` fuentes y el perfil documental integral `29`. El sample oficial `prebuilt-invoice` cubre invoices, utility bills, sales orders y purchase orders con campos tipados, confidence, sources, geometry, spans e ítems. Proformas, packing lists, aduana y schemas particulares siguen requiriendo analyzer/model/schema elegido y corpus del proyecto; no se inventó cobertura.

## 3. Gates globales

```text
VERIFY_LIBRARY_PASS
packs=38 materialized_files=346 markdown_files=233

UPSTREAM_ACQUISITION_TEST_PASS negatives=5
SOURCE_PROFILE_TEST_PASS valid=8 negatives=5 positives=1
UPSTREAM_LOCK_VALID sources=68 selected=68

VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit packs=38 upstream_sources=68 document_sdk_artifacts=5 business_central_artifacts=1 provider_adapters=7
```

Los 17 perfiles compuestos conservan exactamente `23/11/11/24/26/22/25/25/28/26/10/24/25/7/6/95/44` archivos.

## 4. Estado

La ampliación es código oficial inmediatamente adquirible, no un runtime productivo preaprobado. El estado global continúa `NOT_READY_UNDER_EXPANDED_USER_STANDARD`: ejecutar Content Understanding exige Microsoft Foundry, identidad, región, cuotas, model deployments, costo autorizado, corpus/ground truth y gates de seguridad/revisión/rollback.

El ledger retiene `123` fallos propios y `46` condiciones upstream. Los fallos de toolchain, array PowerShell, first-time experience .NET, acceso a propiedad JSON y quoting del ZIP extraído fueron preservados con su corrección.
