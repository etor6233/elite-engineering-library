# Referencia documental conectada — V402 / execution313

T2806 PROVEN_LOCAL para LIBRARY_INFRASTRUCTURE en fixtures y REVIEW_ONLY.
Biblioteca205packs/2270blocks; franquicia113packs/1518files. No READY global,
producción, OCR live ni eficacia antivirus inferida de fixtures.

El mismo host conecta PUT autenticado/limitado, original inmutable y job atómicos,
gate de recepción, SDK oficial Textract1.45.0, respuesta completa tipada y receipts
SHA-bound, propuesta corregida de cuatro campos y revisión de otra persona.
Commit documental, decisión y outbox se confirman juntos. GET devuelve original
/evidencia exactos y recupera resultado tras perder una respuesta. El texto total
no se transforma en dinero/asiento fiscal. La clase supplier-invoice es REQUIRED;
las otras20clases están clasificadas individualmente fuera de esta referencia.

Pruebas en PostgreSQL18.6 y HTTP: recepción/replay, scope tenant/org/rol, SDK,
rechazo/mismatch, corrección, separación, atomicidad bajo falla de outbox,
response loss, fencing de generación vencida, retry finito, cuarentena terminal,
crash del intento final, perfil histórico y supervivencia tras reiniciar host.
No SKIP aceptado.83migraciones seleccionadas; down/up vacío PASS y down poblado
rechazado. Verificador materializado ejecutado;11tests originales del contrato
secure-file PASS; fuzz finito3s; vet y build de cuatro perfiles Go PASS.

FIXTURE restringe el input al JPEG público Microsoft del commit8555d145...,
184686bytes/SHA489f0c63...: el SDK sigue siendo original, pero tanto su output
como detectores están simulados y declarados en original/receipt/commit. No hay
egress ni scanner ejecutado en este modo. PROVIDER conecta el gate original
ClamAV/Magika/YARA y la credential chain AWS, sin fallback a simulación. Su
admisión nativa sigue separada: no se investiga Daybreak/libxml2 aquí, ni se
oculta como una credencial faltante. Documentos privados requerirán expediente
por clase/corpus/política del target; no se pide ningún dato en esta sesión.

Root graph37 =22anteriores +15AWS fijados; ningún módulo anterior cambia versión.
Go mod verify y unión exacta de dos scans OSV2.5.1 sin hallazgos PASS. El floor
x/mod0.40 y sus checksums previos se conservan.24avisos del lint del delta fueron
revisados con hash/posición; los archivos intactos conservan su evidencia312.
El finding1539 fiscal permanece para ARCA penúltimo. Lint no sustituye seguridad
interprocedural/ofensiva ni aceptación productiva.

Se añaden15archivos de glue documental, cinco de licencias/notices del SDK
(uno AUTHORED y cuatro VERBATIM), y se reutilizan17archivos de owners admitidos.
Siete archivos anteriores cambian;1474permanecen iguales. Seis perfiles finales
reconstruidos exactamente. Licencias originales AWS/Smithy Apache-2.0 y Microsoft
MIT retenidas, con fuentes/versiones/SHA; ninguna atribución por reputación.
FAIL927–930 corregidos, conservando RED. Sigue T2807: runtimeIA, evals y aprendizaje
histórico gobernado. Los controles generales TEST02/03/07 se consolidan al final.
