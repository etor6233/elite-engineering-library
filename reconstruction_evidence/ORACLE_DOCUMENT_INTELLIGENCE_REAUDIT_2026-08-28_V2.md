# Oracle Document Intelligence Reaudit — 2026-08-28 V2

## Resultado

Se reauditaron las dos fuentes Oracle ya fijadas en Elite. El SDK oficial OCI Python 2.185.0 queda confirmado como código de integración exacto y ejecutable localmente; el sample Oracle AI Invoice Handling queda rechazado para adopción inmediata porque su flujo humano mostrado en la narrativa no forma parte del proyecto exportado. No se inventó un adapter, un workflow ni una garantía Oracle.

## Oracle OCI Python SDK 2.185.0

- repositorio: <https://github.com/oracle/oci-python-sdk>;
- tag `v2.185.0`, commit firmado `e988c91dcc9963718454cb1e215bd44540524881`, tree `02b796fcfaa00cd6325932104b6ff110c5c20932`;
- source archive ya fijado: 50.664.498 bytes, SHA-256 `7584935da31a4af1898b709500bfa17bec96e68eb232785d999f93bd6f73ddda`;
- licencia dual UPL-1.0 OR Apache-2.0: SHA-256 `8922f6a4bf38ae164a077bd7a294d690cb9e60826d1651d3b36cc7358379bf96`;
- release asset oficial: 221.760.851 bytes, SHA-256/API digest `a655ea06e68a955334d3025e1fb592cdda23ea672776f5de1c56bd6a564728c4`;
- wheel dentro del asset y en PyPI: 36.716.800 bytes, SHA-256 idéntico `93556f08c6270e2e174d59d94fcd49fa794ba3f4d60630af173bc9b6204b085c`;
- wheel: 17.459 entradas; 115 archivos `oci/ai_document`, 822.586 bytes; cero examples Document Understanding.

Astral uv 0.12.7 instaló el wheel exacto en Python 3.14.4, `uv pip check` informó compatibilidad y un probe offline importó `AIServiceDocumentClient`, `AnalyzeDocumentDetails`, `InvoiceProcessorConfig` e `InlineDocumentDetails`. Un fixture de auditoría `oci==2.185.0` resolvió 23 paquetes, lock SHA-256 `40e585ff0ad724839de56d253a29dc31d36b28db32b8c4bacc820589fb78cc92`, y `uv audit --locked` devolvió cero advisories al corte.

Esto prueba disponibilidad e integridad del SDK, no exactitud documental. El `requirements.txt` upstream mezcla pins y rangos; no es un runtime lock universal. No se ejecutó ninguna llamada OCI, no existen credenciales/cuenta/endpoint/costo autorizados y no hay corpus/ground truth del proyecto. La clasificación permanece `INTEGRATION_ONLY`.

## Oracle AI Invoice Handling

- repositorio: <https://github.com/oracle-devrel/oci-ai-invoice-handling>;
- commit protegido pero unsigned `f1f52bb99cdb91c4e4e8332caba2c051fccd638f`; cero releases y cero tags;
- archive: 728.626 bytes, SHA-256 `c3875ae121442d5f3972b5c335a42f5c03ace18037b7ce8b0e88cad0e6eab4a5`;
- licencia UPL-1.0: SHA-256 `c48d4b7187285146c57bf0b035658806be91b1ad65397e5e53ffa9f0f5acd6d9`;
- 102 archivos totales y 78 dentro de `project/`;
- `project.yaml`: dos integrations, siete connections y un label; cero assets Process Automation;
- búsqueda fuera de README/transcript: cero hits de reviewer, approver, human-in-the-loop o process automation;
- cero tests automatizados y cero build gate observados;
- `project.yaml` conserva OCIDs de tenancy e Integration Instance de origen;
- Oracle declara expresamente en README que no representa haber realizado una revisión de seguridad habitual.

El transcript oficial muestra conceptualmente reviewer, comentario, approver y retorno, pero esos artefactos no están distribuidos. El código exportado no demuestra identidad humana, razón obligatoria, correcciones versionadas, lease, CAS, dual control, evidencia durable ni aislamiento. Requiere OCI Document Understanding/Object Storage, Oracle Integration y Fusion ERP, con cuentas y costos externos.

## Decisión gobernante

`oracle-ai-invoice-handling` pasa a `SAMPLE_ONLY_REJECTED_FOR_IMMEDIATE_ADOPTION`; se conserva por hash para investigación y no se compone. `oracle-oci-python-sdk-2.185.0` sigue `INTEGRATION_ONLY` y sólo puede adquirirse tras seleccionar Oracle y completar cuenta, IAM, región, modelo/version, schemas, corpus, seguridad, storage, costo y evaluación. La condición persistente está en la fila upstream Oracle del failure ledger.
