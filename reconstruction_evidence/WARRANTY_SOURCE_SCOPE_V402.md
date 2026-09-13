# V402 — garantía requerida y límites de selección del resto de T2802

Checkpoint284, investigación activa; no promoción nueva ni cierre T2802.
Franquicia vigente84packs/1112archivos, biblioteca180packs/1877bloques.

ENTERPRISE_FULL_STACK_BLUEPRINT.md exige J1 activación de garantía tras entrega y
J4 cita/diagnóstico/trabajo/repuestos/decisión/quality/aceptación/claim. El owner
fulfillment/service_case existente no demuestra por sí solo todo ese recorrido.
Se reutilizarán entrega, casos, stock, aprobación y outbox; no se inventa una
política legal de garantía ni otro ledger.

Fuente seleccionada para inspección: BCApps MIT, commit ya fijado
2eae56d704a1fd035d104f333602aea7091b7749. Once archivos completos de Service y tests,
1261348bytes, coinciden con Git blob/size/SHA. Perfil y expediente previo escritos
antes de adquisición; sin ejecutar AL, comprar runtime ni adoptar BC Platform.
La función CheckWarranty, ServiceItemLine.Table.al1884–1953, evalúa intervalos
inclusivos de piezas/mano de obra y propaga sólo a líneas no facturadas Item/Resource.
La configuración de duración/porcentajes y las exclusiones no se inventan.
Siguiente: adaptación estrecha con oráculo/fixtures fijados, G0–G8 y luego integración
al recorrido existente. Fuente identificada no significa código admitido.

Docs oficiales consultados: [procesos de servicio](https://learn.microsoft.com/en-us/dynamics365/business-central/service-setup-service-processes),
[servicio](https://learn.microsoft.com/en-us/dynamics365/business-central/service-service),
[fault reporting](https://learn.microsoft.com/en-us/dynamics365/business-central/service-how-setup-fault-reporting).
Se usan para ubicar decisiones y fuentes; no sustituyen la implementación fijada.

La tabla del roadmap distingue inventario de cores de requisitos por blueprint.
El perfil no selecciona payroll, POS/caja offline, cupones, recompensa por referido,
reseña pública ni waitlist. Se preparó una conciliación candidata contra J1–J5:
no se han marcado NONE_WITH_REASON ni se removió ninguna REQUIRED para fabricar
el cierre. Garantía conserva REQUIRED porque el blueprint la exige expresamente.
La decisión por cada subcapability debe indicar razón, límites y trigger, conservando
los cores antiguos y sus brechas del catálogo sin prometerlos como listos.

## Verificación general: incidencias nuevas preservadas

FAIL830: script exact-EOF ausente de allowlists de distribución. Se registró por
nombre exacto en ambos owners; el siguiente run superó esa frontera. FAIL831:
selector previo rechaza el JSON COMMUNICATIONS_RUNTIME_V402.json ya existente y
Assert-FranchiseProfileSelection sigue fijando77packs. Requiere conciliación de
recibos públicos/locales y guards del perfil actual bajo T2810. No permitir JSON
arbitrario, omitir desconocidos ni repetir el verificador global sin ese delta.
Las cinco reconstrucciones exactas del gift/loyalty siguen siendo evidencia propia;
no se afirma verify PASS general. Las suites del producto no se reejecutaron.

## Recibos

- `warranty-source-audit/source-inspection-profile.json` SHA256 `f5bea2f32b3c356290d91e507e3dd2bdbe1f6875897b5ac02af8af17e04aebd3`
- `warranty-source-audit/PROJECT_OFFICIAL_SOURCE_PROFILE_RECORD.md` SHA256 `5f0e6b7234dba22dbfbfc41a1a8cd6554b27e7cf37649e8ca309b5817ee9f56b`
- `warranty-source-audit/source-receipts.json` SHA256 `f1890abfa33be3a7f65d567e159769f9e86059dd4fc4ea334d987928d46ae969`
- `warranty-source-audit/check-warranty-segment.json` SHA256 `df316bf431c0667d803ac86afffacaee8dc373a45dea0f3c2ee165d6bd4029fb`
- `franchise-reference-scope-candidate.json` SHA256 `ee5348bb9c00b743bef2ed77b765ffbdc2348d5ea0bd50a2f24bb95c3340103c`
- `gift-tender-library-verify.log` SHA256 `3227dbbb791516b1492cb90aceafb26fe14eba585f36beb9ad0958ff459855a2`
- `gift-tender-library-verify-registered.log` SHA256 `c7bb56e05c75dd700a620d5b0bce66829128896d34fd30a81332955bc22bd5bc`
