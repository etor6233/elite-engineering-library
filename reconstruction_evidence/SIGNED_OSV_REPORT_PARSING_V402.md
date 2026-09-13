# V402329 — lectura correcta del informe OSV

FAIL977 RESOLVED_LOCAL_PARSER. Dos builds reales328 produjeron artefactos idénticos y dos ZIPs byte-idénticos. El gate no firmó porque contó454supuestos hallazgos. La inspección del JSON oficial demuestra454paquetes y0vulnerabilidades: @($pkg.vulnerabilities) convertía una propiedad ausente en un elemento null y el contador sumaba un string vacío por paquete.

Gate y verifier0.2.2 enumeran la colección nullable directamente. Conservan hallazgos positivos y rechazan IDs vacíos. Ambos procesaron el informe real454/0 con PASS;17aserciones de contrato cubren timestamp, fuente, rutas, firma externa, informe limpio, positivo y malformado. No hubo cambio de dependencia ni supresión de vulnerabilidades.

Tres bloques AUTHORED actualizados,1650archivos intactos; reconstrucción1653exacta.3archivos lint/6hallazgos de hashes de integridad y fixtures revisados,0sin resolver. El source tree incluye sus herramientas de release: por eso la nueva cohorte329 se reconstruye y firma con la corrección incorporada. Las builds y ZIPs328 se conservan como evidencia anterior, sin venderlas como329.

Los insumos de reconstrucción están en una ubicación durable fuera de Temp. La firma final, activación, preflight, archivo portable y aceptación final de biblioteca permanecen en curso. Daybreak no se investigó.
