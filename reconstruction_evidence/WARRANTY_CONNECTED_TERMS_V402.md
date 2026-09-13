# V402 — garantía ligada a los términos vendidos

Checkpoint286. Evidencia local estrecha en warranty-connected-candidate; el perfil
canónico84/1112 no cambia todavía. T2802 permanece abierto por el reclamo J4,
superficie HTTP/host/UI y reconciliación de las demás filas.

PASS: un recorrido PostgreSQL sin skips crea una cotización real, fija términos
por hash/version, exige consentimiento del cliente, convierte a pedido, asigna
stock, ejecuta SDK Stripe contra fixture, callback/reconciliación, checklist,
aceptación del handover, recibo comercial y activación durable de garantía.
Dos términos configurados diferentes demuestran que la activación conserva
la versión vendida. El reembolso posterior no inventa anulación legal.

La falta de consentimiento revierte también el pedido; una falla inyectada en
outbox revierte garantía y activación. Dos consentimientos concurrentes producen
un hecho; replays exactos no duplican eventos. Otro cliente, tenant, organización,
versión/hash incorrectos, precio alterado y mutación de recibo son rechazados.
Dos tests de contrato prueban inmutabilidad, rechazo de timezone Local/valores
implícitos/duplicados y calendario DST con bordes inclusivos del owner admitido.
Go vet PASS; compilación de paquete de tests tiene cero crédito funcional.

Se reutilizan writers de quote/order/payment/handover; sólo el predicado de fecha
es ADAPTED de BCApps. Consentimiento, calendario explícito y composición SQL son
AUTHORED declarados. No plazo legal, porcentaje, indemnización ni cobertura por
reputación. Los writers de stock/caso extraídos conservan sus SQL; el uso conjunto
en J4 aún requiere test de la nueva transacción. Sin nueva dependencia externa.

Preparación: dos errores del script de extracción (función final sin sucesor y
alias SQL r. confundido con receiver Go) se corrigieron conservando las salidas
parciales; la compilación final está separada de la previa. No fueron fallos de
runtime ni se acreditan como pruebas funcionales. FAIL832 documenta esta lección.

Manifest y hashes de once archivos y seis recibos: WARRANTY_CONNECTED_TERMS_V402.json.
PG66migraciones aplicadas y detenido al finalizar; no cuentas live. El adaptador de
fechas mantiene su expediente G0–G8 estrecho anterior. No verify global repetido.
