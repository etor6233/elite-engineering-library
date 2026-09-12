# Microsoft Content Understanding GA Source — 2026-08-26 V1

## Autoridad observada

Microsoft recomienda oficialmente el SDK GA `azure-ai-contentunderstanding` para producción en lugar de los wrappers REST livianos de `Azure-Samples/azure-ai-content-understanding-python`. La release estable vigente auditada es `1.1.0`, vinculada al API GA `2025-11-01`; las releases `1.2.0b*` son preview y no reemplazan este pin.

## Artefacto real verificado

```yaml
owner: Microsoft
package: azure-ai-contentunderstanding
version: 1.1.0
artifact: azure_ai_contentunderstanding-1.1.0.tar.gz
bytes: 230330
sha256: 00f39c7cf2ba50c8d4586315f4295a45fad0881daaff07d13d33aea2bafc7ec8
license_expression: MIT
upload: 2026-04-20T23:29:08.315545Z
trusted_publishing_reported: false
```

La descarga desde `files.pythonhosted.org` coincidió exactamente en bytes y SHA-256 con el JSON oficial de PyPI. El archive contiene 138 entries: 105 Python, 68 samples y 43 tests. Incluye runtime sync/async, models, long-running operations, analyzers, defaults, classification, labeled training, result management y diagnostics.

## Separación de claims

- el sdist y wheel son runtime oficial Microsoft `Production/Stable`;
- el repositorio tutorial GA sigue siendo sample y recomienda el SDK;
- la solución `data-extraction-using-azure-content-understanding` incluye Functions, Cosmos DB, Key Vault, Storage, telemetría y 55 Terraform, pero su propio README exige adaptación productiva;
- ningún artefacto aporta cuenta, región, modelos desplegados, corpus, schema empresarial ni prueba de exactitud REVESTEX.

## Decisión

El perfil documental debe depender del SDK GA exacto para ejecución Azure. Los samples oficiales gobiernan contract probes; la solución completa gobierna arquitectura/configuración condicionada. No se copiará un wrapper REST antiguo como si fuera el runtime recomendado ni se elevará el accelerator a producto terminado.
