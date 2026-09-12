# All implementation packs materialization — 2026-08-26 V20

## 1. Snapshot gobernante

Este snapshot sucede a V19. La biblioteca contiene `39` implementation packs y `352` archivos materializables. Mantiene `68` upstreams oficiales fijados y añade una lane ejecutable local sobre la fuente Microsoft MarkItDown ya admitida; no agrega un upstream redundante.

`MICROSOFT-MARKITDOWN-LOCAL-RUNTIME 0.1.0` materializa seis archivos: lock de 44 wheels con hashes, profile local fail-closed, wrapper `convert_local`, environment verifier, cinco regresiones y README operativo. El engine/converters son el wheel oficial Microsoft 0.1.7; el glue permanece `AUTHORED` y no se atribuye a Microsoft.

## 2. Ingesta documental ampliada

La nueva lane convierte localmente PDF, DOCX, XLS/XLSX, PPTX, Outlook MSG y formatos textuales después de type/antimalware/quarantine. URL, plugins, archives, audio/video, YouTube, LLM y almacenamiento de negocio están deshabilitados. Se probó sobre seis fixtures oficiales upstream y conserva hashes input/output en receipt atómico.

El grafo CPython 3.12 Windows x86-64 contiene 44 wheels/82.672.112 bytes, instala offline, pasa `pip check`, environment verifier y consulta OSV 44/0 findings conocidos. No se creó ni afirmó un lock Linux/macOS. Conversión no equivale a extracción semántica: proformas, packing lists, BOL, aduana y certificados continúan requiriendo analyzer/schema/corpus oficial del provider seleccionado.

La reauditoría adicional rechazó como runtime reusable el workshop AWS BDA por ausencia de lock/tests y mantuvo Google Document AI Toolbox fuera del perfil estable porque Google lo clasifica Experimental/Pre-GA. Las decisiones están en `OFFICIAL_MULTI_CLASS_DOCUMENT_CODE_REAUDIT_2026-08-26_V1.md`.

## 3. Gates globales

```text
VERIFY_LIBRARY_PASS
packs=39 materialized_files=352 markdown_files=238

UPSTREAM_ACQUISITION_TEST_PASS negatives=5
SOURCE_PROFILE_TEST_PASS valid=8 negatives=5 positives=1
UPSTREAM_LOCK_VALID sources=68 selected=68
MICROSOFT_MARKITDOWN local wrapper: 5 tests PASS

VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit packs=39 upstream_sources=68 document_sdk_artifacts=5 business_central_artifacts=1 provider_adapters=7
```

Los 18 perfiles compuestos conservan exactamente `23/11/11/24/23/26/22/25/25/28/26/10/24/25/7/6/95/44` archivos: inicialización, Azure, Google, AWS Textract, MarkItDown local, AWS enterprise, Google Ads, Meta Ads, TikTok Ads, WhatsApp, Firebase, Mercado Libre, Amazon Catalog, Google Merchant, payments, Business Central, backend y web.

## 4. Fallos y estado

El ledger retiene `132` fallos propios y `48` condiciones upstream. En este ciclo quedaron registrados paths Windows profundos, dependencia asumida del arnés, helper temporal, toolchain no comprobado, URL PyPI inferida, parser OSV nulo y sesión de instalación larga; todos se repitieron con corrección o fueron rechazados por gates.

El estado global continúa `NOT_READY_UNDER_EXPANDED_USER_STANDARD`. La biblioteca entrega conversión local y runtimes cloud materializables de inmediato, pero ningún proveedor público aporta un modelo universal que elimine la definición real de schemas, corpus, cuentas, reglas comerciales, acceso, revisión, carga, seguridad, rollback y recuperación del proyecto.
