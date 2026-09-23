# Decisión documental de la composición de referencia

Derivada de PROJECT_DOCUMENT_INTELLIGENCE_DECISION_TEMPLATE.md. Alcance único:
LIBRARY_INFRASTRUCTURE / LOCAL_FIXTURES / REVIEW_ONLY. No habilita documentos
privados, exactitud OCR, contabilización, reglas fiscales ni almacenamiento automático.
La autorización del usuario es completar infraestructura sin cuentas ni secretos.
Owner de este perfil: mantenimiento de biblioteca; operación/revisión/incidentes:
agente durante el ensayo, usuario cuando materialice su propio target.

## Ingesta y conservación

Canal REQUIRED: API autenticada por el verifier OIDC existente, permiso
`documents:write`, tenant y organización. Un original por PUT y UUID; bytes
SHA-256 y profile SHA obligatorios. Mismo ID+actor+nombre+bytes+perfil devuelve
el mismo registro; diferencias rechazan. JPEG/PDF, máximo2MiB, una página para
el processor. FIXTURE sólo admite el JPEG público exacto, no cualquier PDF.
El cuerpo se recibe completo y acotado antes de abrir la transacción; cargas
parciales no crean original/job. Timeout de proceso2min; lease3min;3intentos;
una extracción por job reclamado. El operador posee `documents:process`.

El original se conserva como bytea con hash calculado también por PostgreSQL,
FK tenant/org, timestamp de la base y trigger que impide update/delete. Es
inmutabilidad de aplicación/base bajo roles acotados, no WORM frente al DBA.
Un superusuario sigue pudiendo cambiar la infraestructura. FIXTURE no contiene
cuentas/secretos del usuario. En este ensayo local, acceso del SO restringido al
owner; cifrado/KMS, residencia, plazo legal, legal hold y borrado de documentos
privados son decisiones del futuro target, sin defaults regulatorios inventados.
La retención del fixture dura el ensayo y se conserva en su evidencia; el módulo
no borra evidencia ni ejecuta down con documentos presentes. Temporales por
intento se eliminan al retornar; huérfanos por caída del host los gestiona el
supervisor local/operación T2809. No se promete eliminación segura del disco.

Estados proyectados: QUARANTINED, REVIEW_REQUIRED, REVIEW_PENDING, REJECTED,
PERSISTED y QUARANTINE_TERMINAL. El hash original, perfil, security receipt,
respuesta completa tipada y analysis receipt se enlazan antes de revisión.
GET de original y tres partes de evidencia está autenticado, devuelve bytes
exactos con SHA, no-store, attachment y nosniff. No se registran campos privados
ni salidas de detectores en logs. El outbox sólo lleva IDs/hash/mode.

## Seguridad y ejecución

PROVIDER llama el gate existente SECURE_LOCAL_FILE_INGESTION_GATE: intérprete,
script y policy con SHA fijados; ClamAV1.5.4, Magika1.1.0 y YARA-X1.20.0 conservan
sus propios hashes/contratos. Firmas oficiales frescas, type/MIME concordante,
reglas YARA fijadas, límites y rechazo permanecen en ese owner. Archives,
macros, archivos cifrados, multipágina y clases desconocidas no se habilitan.
No hay override por receipt aportado desde HTTP. El SDK se llama después del
gate, con una sola tentativa interna. Los reintentos de la cola son explícitos;
una respuesta perdida puede repetir una extracción remota de sólo lectura y
su costo, nunca autoriza duplicar el commit. Fallo de scanner/provider/schema
conserva código seguro por intento y cuarentena; agotamiento permanece terminal.

FIXTURE usa transporte HTTP en memoria del SDK original y contrato de detector
simulado, limitado al SHA público. `security_claim`, mode y commit lo declaran.
No prueba eficacia de antivirus, firmas reales ni libxml2. Esa rama nativa
ACCESS_BLOCKED/Daybreak queda para el expediente final, sin investigación aquí.
No se presenta como una credencial faltante ni se oculta código incompleto.
El runtime PROVIDER está conectado al gate; sólo el ensayo fixture está probado.
No hay datos sensibles a redactar antes del proveedor fixture: no existe egress.
Para datos privados, clasificación/redacción/retención y scanners del target
son condiciones de activación, no hechos inferidos del PASS de este ensayo.

## Clases, una decisión por clase

| Clase | Decisión referencia | Variante / motivo | Owner |
|---|---|---|---|
| invoice / supplier invoice | REQUIRED / REVIEW_ONLY | JPEG público Microsoft fijado; respuesta AWS simulada, sin ground truth de exactitud | mantenimiento de biblioteca |
| receipt / ticket | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |
| proforma invoice | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |
| commercial invoice | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |
| packing list | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |
| purchase order | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |
| purchase order confirmation | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |
| bill of lading / airway bill | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |
| delivery note / remito | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |
| certificate of origin | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |
| quality/inspection certificate | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |
| customs declaration/clearance | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |
| insurance/freight document | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |
| price list/quotation | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |
| product specification/catalog | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |
| contract/addendum | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |
| bank/payment statement | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |
| warranty/service/claim document | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |
| identity or regulatory document | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |
| multi-document package | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |
| other named class | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |

La fila commercial invoice no amplía supplier-invoice a comercio exterior; cada
clase excluida conserva su gate propio. Idioma/país/reglas tributarias no se
infieren de nombres simulados. Volumen: ensayos acotados; concurrencia probada:
dos generaciones del mismo job. No SLO ni carga de documentos privados declarados.

## Campos y revisión

Schema `supplier-invoice/v1`, una entidad por campo, sin line items ni enlaces
PO/invoice/packing list. Proyección directa de los tipos oficiales Textract:

| Campo | Campo AWS | Tipo / límite | Política |
|---|---|---|---|
| invoice_number | INVOICE_RECEIPT_ID.ValueDetection.Text | texto128bytes | REQUIRED / REVIEW_ONLY |
| vendor | VENDOR_NAME.ValueDetection.Text | texto512bytes | REQUIRED / REVIEW_ONLY |
| total | TOTAL.ValueDetection.Text | texto64bytes; no se convierte a dinero | REQUIRED / REVIEW_ONLY |
| currency | TOTAL.Currency.Code | tres letras ASCII mayúsculas | REQUIRED / REVIEW_ONLY |

No se reemplaza una ausencia por cero ni confidence por exactitud. Faltantes
quedan visibles para corrección; la propuesta exige los cuatro valores. La
revisión consulta original y respuesta oficial completa (incluidos campos de
confidence/provenance que el SDK conserve). Una propuesta es inmutable y está
ligada a original/perfil/extracción; otra persona con `documents:review` decide
sobre su SHA exacto. No hay autoapprove aunque exista umbral global positivo.
El commit, decisión y outbox se confirman en la misma transacción. Rechazar no
crea registro committed. Para corregir una propuesta rechazada se crea otro
expediente/document ID, conservando el original previo; no se reescribe historia.
Los cuatro campos siguen siendo un registro documental revisado, no un asiento
contable, pago o comprobante fiscal. Esas mutaciones pertenecen a otros owners.

## Corpus y autoridad

Un fixture público adquirido por commit/tamaño/SHA y licencia MIT original.
Ground truth de OCR:0; entrenamiento:0; test de exactitud:0. El output simulado
no está etiquetado como lectura de la imagen. Las pruebas demuestran transporte,
persistencia, aislamiento, review, recovery y contratos; no métricas de extracción.
Se selecciona AWS Go SDK Textract1.45.0 (15módulos fijados) por compatibilidad con
la composición. Jobs PostgreSQL evita un segundo runtime/in-memory workflow;
DurableTask, Stickler y CEL Java no son necesarios para este recorrido manual.
No se atribuye el glue a Microsoft/AWS/Google ni a personas famosas.

PROVIDER requiere cuenta/IAM/credenciales, región/egress/cuotas/costo admitidos y
scanners preparados. El perfil público sólo habilita supplier-invoice y
`automatic_storage=false`. Ampliar a corpus privado requiere ground truth,
política por clase/campo y jurisdicción del usuario, sin solicitarla en esta sesión.
Backup/restore, alertas, retención del supervisor y carga general se consolidan
en T2808/T2809/T2810; este documento no los promueve por inferencia.

## Evidencia y reapertura

`tools/verify_document_reference.py` exige fixture/Go exactos, DB loopback de
ensayo y receipt nuevo; rechaza SKIP. Tests cubren transacción, negativos de
scope/hash/review, replay, falla de outbox, fencing, cuarentena y cambio de perfil.
Fuzzing finito3s sobre límites de original; no sustituye SAST/carga/seguridad live.
La evidencia pública T2806 se enlaza en DOCUMENT_REFERENCE_RELEASE_V402.md/json.
Cualquier cambio de owner/perfil/input/schema/permiso/gate invalida el claim
relacionado. Cambios de corpus/modelo/país reabren su expediente antes de automatizar.
