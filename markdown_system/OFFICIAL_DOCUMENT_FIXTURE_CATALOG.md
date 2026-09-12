# Official Document Fixture Catalog

Fecha de corte: 2026-08-28.

Fixtures públicos contenidos dentro de los archives oficiales exactos. Sirven para smoke/contract tests del pipeline; no reemplazan corpus representativo ni ground truth del proyecto. La licencia es la del source archive y puede contener términos adicionales por path.

## Facturas Microsoft Azure

Source ID `azure-ai-document-processing-samples`, raíz `samples/assets/invoices/`:

| Path | Bytes | SHA-256 |
|---|---:|---|
| `invoice_1.pdf` | 272,192 | `ab8ec59d22803a922b69ef412adba9360d77fea835d35a75ed8f66f181483c8e` |
| `invoice_2.pdf` | 82,628 | `13e8c5bc3040e74eab031aebd74bf58ed203bb85131c5f41b608c14d4c023fc5` |
| `invoice_3.pdf` | 567,893 | `1a4e016c5300ee398c564e94a78da1c97501f056953cef892c1ffaed760baac1` |
| `invoice_4.pdf` | 570,699 | `a19115e7ce1bfa6ab8382abc074036aa9414aaa53187f782a9e536bc489c16d9` |
| `invoice_5.pdf` | 358,218 | `fc18bf8fffc2c4e63a82aee18c9e921034bc70239b09fffe48791af8f15a80d3` |
| `invoice_6.pdf` | 2,279,694 | `00ad2587ff68d2231213ca738b961553113582165e230a6a5b9a572c3dcbe988` |

El source incluye notebooks Python/.NET de invoice extraction y un schema común. Estos PDFs prueban conexión/prebuilt/boundary, no los formatos de un proveedor del usuario.

El pack `MICROSOFT-AZURE-DOCUMENT-INTELLIGENCE-OFFICIAL-INVOICE-SAMPLE` fija además el fixture publicado con el SDK Azure Document Intelligence 1.0.2 y lo enlaza al sample oficial desde bytes:

| Source/path | Bytes | SHA-256 | Uso estrecho |
|---|---:|---|---|
| `Azure/azure-sdk-for-python@8555d145…/samples/sample_forms/forms/sample_invoice.jpg` | 184.686 | `489f0c63b6a05e0ae0fcfbf29121299fc2ccacab35d5cd40a5d63d42764078cb` | ejecución de integración `prebuilt-invoice`; no ground truth productivo |

El JPEG no se redistribuye en la biblioteca: se adquiere desde el raw commit-pinned con aprobación MIT, tamaño/hash, staging y receipt.

El pack `MICROSOFT-AZURE-CONTENT-UNDERSTANDING-OFFICIAL-INVOICE-SAMPLE` conserva además byte-verbatim el sample Python oficial Microsoft Azure Content Understanding 1.1.0 que llama `prebuilt-invoice` y proyecta confidence, source, spans, line items y usage. Su URL de demostración apunta a un asset móvil bajo `main`; por eso no se incorpora como fixture fijado ni se usa con documentos del proyecto. El pack demuestra identidad e integración del código, no exactitud ni persistencia.

## Reclamos multiarchivo Microsoft

Source ID `microsoft-content-processing-accelerator-main-659eaa1`, commit firmado `659eaa1f503dd08b1e1aea1c72eab11c7c191d00`, raíz `src/ContentProcessorAPI/samples/`:

| Path | Bytes | SHA-256 | Uso estrecho |
|---|---:|---|---|
| `claim_date_of_loss/claim_form.pdf` | 3.576 | `db2ae7da7cecff1b26e3f0038ccfc777d63ed1177738f72c61459d1d7581059e` | formulario de reclamo multiarchivo |
| `claim_date_of_loss/police_report.pdf` | 111.743 | `d23f28dfa066d16d2767687e5e6da73792a29afca70a3bb23e231600da918f2b` | reporte policial asociado |
| `claim_date_of_loss/repair_estimate.pdf` | 3.137 | `7790f79f349794afd7fa1735f9de4292b8cadb79370bf1dc6893e6f186935b6f` | presupuesto de reparación asociado |
| `claim_hail/claim_form.pdf` | 8.606 | `1e009d06685e9c40ee90897bf7115a2a2268442fccdfc77cda5324d91afd7567` | formulario de reclamo por granizo |
| `claim_hail/repair_estimate.pdf` | 123.752 | `50e683bdd13b5e848090852e967b58ef211ba34b35d16dc6a0a672c709c64330` | presupuesto de reparación por granizo |

Tres copias E2E de `claim_date_of_loss` son byte-idénticas y no se cuentan como fixtures nuevos. Estos archivos permiten probar clasificación y extracción cruzada entre documentos; no constituyen ground truth empresarial ni corrigen el rechazo de la rama para adopción inmediata.

## Facturas y workflow SAP

Source ID `sap-btp-dox-invoice-validation`, raíz canónica `api/srv/samples/`:

| Path | Bytes | SHA-256 |
|---|---:|---|
| `3420987413543.pdf` | 173,407 | `a04f1db97667750006fbce8a341afecbbda42ea791ee006042f1c5683646fcc1` |
| `5435569865439.pdf` | 193,822 | `a69c5a4ff56ca03d2a32240f7f12de4f4fe3921edf7e5c5006b0d7f6d91a5e21` |
| `6632559877890.pdf` | 177,475 | `cd6ecd62d5de83739d80c6c31dbca428151dd25adda61df1f50de7ed3fe64328` |

Los duplicados bajo `docs/samples/` tienen los mismos hashes y no se cuentan como fixtures independientes. El sample SAP conecta estos PDFs con DOX, line items, schema custom, correcciones y aprobación.

## Packing list, factura y recibo Google

Source ID `google-document-ai-samples`, raíz `web-app-pix2info-python/pics/`:

| Path | Bytes | SHA-256 | Uso estrecho |
|---|---:|---|---|
| `16a_packing_list_with_barcode.gif` | 153,969 | `5d7da61ef1c1717fc4a4578d7f4f4d9dc7d1146e04d49db8ec7073b13bda88ed` | packing-list/barcode visual smoke fixture |
| `15a_invoice.gif` | 1,131,656 | `eaadb1f7c318ec938c1cf8f34654bc83b569be0f41e415fd66b04d272a79f719` | invoice visual/example |
| `14e_receipt.jpg` | 178,483 | `6fbe1b421dec342d87947172f3feac7d33c25281031800ba6a94f91452bad326` | receipt visual/example |

La packing list es una imagen pública de ejemplo y no un dataset etiquetado. Puede verificar ingestión/barcode/layout; no autoriza un claim de extracción de campos de packing lists reales.

El pack `GOOGLE-CLOUD-DOCUMENT-AI-OFFICIAL-PROCESS-SAMPLE` fija además la pareja input/output ubicada bajo `src/samples/4-INVOICE_PROCESSOR/` del commit firmado `001ba391…`:

| Path | Bytes | SHA-256 | Uso estrecho |
|---|---:|---|---|
| `b. Packing list with barcode.png` | 69.350 | `b45f86a82194a6210ec3f2b149a9c486934b01f6087cf7e6f859f32c5f709e7b` | input oficial |
| `b. Packing list with barcode.json` | 361.762 | `509cdc9e763f4c37253c617fba800452473c7a7e2bc5545ad6f6c5db956b91fb` | output Document AI cacheado oficial |

El output conserva barcode y siete entidades, pero varios campos tienen confianza inferior a 0,56. Se usa para smoke/contrato del formato, nunca como ground truth ni autorización de persistencia.

## Queries AWS

Source ID `aws-textract-queries-example`, raíz `training/`:

| Path | Bytes | SHA-256 | Uso |
|---|---:|---|---|
| `bank-statement.pdf` | 66,842 | `99bc7cebf2fc63d70cc3fedb0d49b3f4338c578e55d86130ce16adc0930b0c9c` | Queries/classification smoke |
| `payslip.pdf` | 72,010 | `96170cc6f14f9d4a8def78fc8fe66edd11184acefa71ec31f2abdc8e5cf0fb37` | Queries/classification smoke |
| `rental-application.pdf` | 73,464 | `df7456c9138f278c8fbded58228bc3ab74fa242e9e6611e46a1d1fb1d32efe3e` | Queries/classification smoke |

## Tablas y forms NVIDIA

Source ID `nvidia-nemo-retriever-main-2026-08-26`, raíz `data/`:

| Path | Bytes | SHA-256 |
|---|---:|---|
| `table_test.pdf` | 26,342 | `b81f79487b9956997c00e3266e5a320a65180f8b8a4c4d8d76fbeabb822025a3` |
| `embedded_table.pdf` | 192,612 | `7eeacd3c1a5a0111dd2a5a7204601d1d30636aa10a3af4fee1bb428a62a93d38` |
| `test-page-form.pdf` | 724,323 | `a358e820550124889ae570a6f384d58b80b3ee837645cdfb983e213ec2f66aa7` |

Sólo se usan dentro del entorno Linux/GPU/NVIDIA admitido para el commit fijado; no prueban precisión productiva.

## Forms/tablas Microsoft accelerator

Source ID `azure-document-intelligence-in-a-box` contiene sets train/test/combined de Contoso Safety 360 y un set específico de tablas. Fixture de smoke recomendado:

| Path | Bytes | SHA-256 |
|---|---:|---|
| `data/samples/test/contoso_test_table/ContosoSafety360-Sample-Table-Test-1-New.pdf` | 76,145 | `488eb13d83b42bc963d97ff7c500d3f89cae496ecb86740329b9f7af540d89cd` |

El archive conserva los demás ejemplos y su propio split; no mezclar train/test al evaluar.

## Gate de uso

1. adquirir el source ID mediante el runner y verificar receipt;
2. recalcular SHA-256 del fixture antes de usarlo;
3. conservar output crudo, provider/model/version, request ID y costo;
4. comparar sólo con labels/expected output oficiales cuando existan;
5. marcar el resultado `SMOKE_ONLY` si no hay ground truth;
6. nunca usar estos ejemplos públicos para decidir autoalmacenamiento de documentos privados;
7. incorporar fixtures reales autorizados y casos difíciles en un dataset separado del proveedor y versionado por el proyecto.

No se encontró dentro de estos archives un corpus oficial suficiente y etiquetado de proformas, packing lists, purchase orders, bills of lading, aduana y certificados que demuestre exactitud universal. Purchase orders disponen ahora del candidato oficial Azure Content Understanding `prebuilt-procurement`/`prebuilt-purchaseOrder`, pero continúan requiriendo evaluación con documentos reales; las demás clases citadas permanecen `CUSTOM_MODEL_REQUIRED` hasta aportar sus documentos y schema.
