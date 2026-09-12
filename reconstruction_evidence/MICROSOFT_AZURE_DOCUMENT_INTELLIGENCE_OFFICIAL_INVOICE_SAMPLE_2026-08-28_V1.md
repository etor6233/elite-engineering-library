# Microsoft Azure Document Intelligence Official Invoice Sample — Evidence V1

Date: 2026-08-28  
Scope: sample invoice desde bytes, licencia y fixture publicados por Microsoft; pack/reconstrucción offline local.  
No cubre: llamada live a Azure, exactitud de campos del negocio, persistencia productiva, SLA, cuota/costo ni otras clases documentales.

## Upstream authority

- Owner/repository: Microsoft, `Azure/azure-sdk-for-python`.
- Tag: `azure-ai-documentintelligence_1.0.2`.
- Commit exacto: `8555d14532a9688b751d8408d822d1dd5feb47f6`.
- Fecha del commit: `2025-03-26T23:01:34Z`.
- Verificación GitHub consultada el 2026-08-28: `verified=true`, `reason=valid`, firma presente.
- El sample sync descargado desde raw GitHub coincide byte-a-byte con el miembro del sdist oficial PyPI `azure-ai-documentintelligence 1.0.2`.

## Exact public artifacts

| Artefacto | Bytes | SHA-256 | Resultado |
|---|---:|---|---|
| `sample_analyze_invoices_from_bytes_source.py` | 13.084 | `ebc32b0fc534e625b15e636c6b55376f01d28a83f35217dd141202073297a293` | `VERBATIM`, Microsoft MIT |
| root `LICENSE` | 1.074 | `7c77a44a8acd9b41fdc209864a8016b3d430b5d0e09309818d5b7444336df744` | `VERBATIM`, MIT |
| `sample_forms/forms/sample_invoice.jpg` | 184.686 | `489f0c63b6a05e0ae0fcfbf29121299fc2ccacab35d5cd40a5d63d42764078cb` | fixture oficial commit-pinned, no redistribuido |

El fixture se adquiere únicamente desde `raw.githubusercontent.com/Azure/azure-sdk-for-python/<commit>/...` mediante approval enlazado al hash del lock, tamaño/hash exactos, staging, receipt y destino inexistente. La red queda deshabilitada por defecto y una cache manipulada falla cerrada.

## Reconstruction and tests

El pack `MICROSOFT-AZURE-DOCUMENT-INTELLIGENCE-OFFICIAL-INVOICE-SAMPLE` 0.1.0 materializó ocho archivos desde Markdown. El perfil mínimo materializó catorce archivos desde dos packs, más `MATERIALIZATION_RECORD.md`, sin colisiones; el sample reconstruido conservó el SHA upstream.

Resultados offline:

```text
python test_official_invoice_sample.py
Ran 2 tests ... OK

AZURE_DI_OFFICIAL_FIXTURE_LOCK_VALID fixtures=1 selected=1
AZURE_DI_OFFICIAL_FIXTURE_ACQUISITION_PASS fixtures=1
AZURE_DI_OFFICIAL_FIXTURE_TEST_PASS positives=2 negatives=3
```

La prueba Python carga la función real del sample Microsoft y demuestra que abre los bytes del fixture, crea `DocumentIntelligenceClient` con endpoint/credential recibidos, llama `begin_analyze_document` con `prebuilt-invoice`, `locale=en-US` y `AnalyzeDocumentRequest(bytes_source=...)`. Los módulos Azure son dobles sólo dentro de la prueba offline; no se reimplementa ni modifica el sample.

Los negativos PowerShell prueban ID desconocido, destino ocupado y cache con hash incorrecto. Las pruebas no hacen una llamada live ni demuestran precisión semántica del proveedor.

## Admission boundary

Clasificación: `SUPPORTED_REFERENCE / REBUILD_VERIFIED / CONDITIONED`.

Este resultado cierra la brecha técnica de tener un sample oficial de invoice realmente materializable junto con su fixture oficial exacto. No convierte los otros 55 samples del sdist en ejecutables autónomos ni permite afirmar que toda ingesta documental está terminada. Para uso productivo siguen siendo obligatorios: decisión de proveedor, cuenta/secretos/costo autorizados, corpus y ground truth por clase, schema y reglas del negocio, evidence envelope, evaluación, revisión de ambigüedad, seguridad de archivo, observabilidad y política de persistencia.

## Failure learning incorporated

La construcción detectó y preservó como regresiones: búsqueda de tag mal invocada, bytes finales extra en sample/licencia, precedencia incorrecta al comparar aprobaciones, contexto de patch supuesto, path temporal supuesto, placeholder SHA no contractual y nombre de parámetro incorrecto. Los gates finales exigen identidad byte-a-byte y rechazos, por lo que estas clases no se silencian.
