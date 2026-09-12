# AWS BDA Purchase Order Blueprint Optimizer — V82

Fecha: 2026-08-28  
Estado: `SOURCE_LOCKED / PINNED_CANDIDATE / CONDITIONED`

## Fuente oficial

- repositorio: [`aws-samples/sample-blueprint-optimizer-for-data-automation`](https://github.com/aws-samples/sample-blueprint-optimizer-for-data-automation/tree/ddb5ae4b3c8d7e33c4ade53301f0086970e32134);
- commit GitHub verificado: `ddb5ae4b3c8d7e33c4ade53301f0086970e32134`, 2026-06-10;
- archive: 1.803.458 bytes, SHA-256 `da53ca0ac2c87f6639c07325c0df53ec84e609f1c753b16860d8dbb7323a1687`;
- licencia MIT-0: 946 bytes, SHA-256 `5025c7bbedbc8da868b6dce2ab689225b3f33c43b3f03368b0a331809b6ada26`;
- autoridad adicional: post técnico AWS [Optimize blueprint extraction accuracy in Amazon Bedrock Data Automation](https://aws.amazon.com/blogs/machine-learning/optimize-blueprint-extraction-accuracy-in-amazon-bedrock-data-automation/), 2026-06-11.

## Código y corpus real

El archive contiene 42 archivos:

- `purchase-order-optimization-workshop.ipynb`, 19.945 bytes/SHA-256 `c899b1e1c0c66ffe48f5710c9a281fdfac1e936637a507e9546381ca0407e184`;
- schema BDA purchase-order, 3.364 bytes/SHA-256 `c76ba133046419317ed7aca10c2ddc960db952878ac479f5ac32bf6fe6be2aa8`;
- 15 PDFs purchase-order y 15 ground truths JSON;
- CloudFormation SageMaker, config, requirements y README exactos.

El notebook usa las APIs oficiales `CreateBlueprint`, `InvokeBlueprintOptimizationAsync`, `GetBlueprintOptimizationStatus`, `GetBlueprint`, `CopyBlueprintStage` y cleanup. El artículo AWS reporta para su muestra mejoras de exact match por archivo hasta 100% y aggregate exact match de 90% a 92%; esos valores no se generalizan a otro negocio.

## Verificación offline

- 17/17 JSON parsean;
- schema: 7 propiedades raíz y 15 propiedades de item;
- 15/15 ground truths cubren exactamente las propiedades observadas;
- 75 ítems, cero faltantes/extras;
- cero diferencias `quantity × unit_price ↔ line_total`;
- cero diferencias `sum(line_total) ↔ order_total`;
- cero secretos AWS embebidos detectados en celdas;
- el path Windows de un ground truth alcanzó 260 caracteres; lectura con prefijo extendido probó que era archivo regular, no ausencia ni symlink.

## Condiciones retenidas

No se promueve a implementación productiva:

- cero tests upstream y cero dependency locks;
- `requirements.txt` usa rangos abiertos;
- CloudFormation contiene cuatro recursos wildcard;
- schema no declara `required` ni `additionalProperties:false`;
- polling, cuenta, IAM, S3, SageMaker, BDA, optimización y promoción LIVE requieren autorización/costo;
- no existe corpus REVESTEX/target, gate de seguridad, review, rollback o prueba live autorizada;
- cubre purchase order, no proforma, packing list, bill of lading ni customs.

`UP-FAIL-171` conserva estas condiciones. `UP-FAIL-170` rechaza el workshop BDA genérico `131ea7b…` para el objetivo logístico porque no contiene clases objetivo, tests ni manifests.

## Integración en Elite

- `OFFICIAL-UPSTREAM-ACQUISITION-CORE` 0.4.64: 111 fuentes exactas;
- profile `document-intelligence-leaders`: 46 fuentes;
- source entry `aws-bda-blueprint-optimizer-purchase-order-ddb5ae4` con archive/licencia/requirements/schema/notebook/CloudFormation hash-locked;
- 22 planes consumidores alineados a source pack 0.4.64;
- `UPSTREAM_ACQUISITION_TEST_PASS negatives=9`;
- `SOURCE_PROFILE_TEST_PASS valid=14 negatives=5 positives=1`;
- materialización source pack: 23/23 archivos;
- `VERIFY_LIBRARY_PASS`: 59 packs, 548 archivos, 392 Markdown antes de esta evidencia;
- `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit`: 59 packs, `upstream_sources=111`, `document_sdk_artifacts=12`.

El agente puede adquirir el código/corpus AWS exacto sin Git y usarlo como base de evaluación de purchase orders. No puede instalarlo, desplegarlo, promover LIVE ni persistir resultados hasta cerrar todas las condiciones del proyecto.
