# V402 / checkpoint308 — T2805 conectado

PROVEN_LOCAL para el alcance seleccionado de comunicaciones/adapters:
WA ingreso/propuestas/revisión/envío/estado y Page conservados desde V402;
ML mutaciones/MEDIA/CREATE/CONTENT y Merchant catálogo/refresh/reconciliación;
recordatorios durables y ahora campañas con audiencias permitidas y conversión.
108packs/1451archivos;203packs/2215bloques. No READY global ni producción.

Campaña: snapshot explícito del lead/contacto/consentimiento de marketing,
filtros source/lifecycle, destinatarios exactos y pasos de plantilla revisados.
Se reutiliza el owner de aprobación, jobs y fence/adapter; no segundo scheduler
ni fórmula de atribución.6campañas/10jobs/4POSTreales/10resultados y1cotización
convertida en pedido por el owner original.12workers concurrentes→1primer envío.
Stop, retiro de consentimiento, cambio del lead, antecesor incierto y conversión
suspenden los pasos correspondientes. Recovery valida el recibo sin reenviar.

La preparación/revisión parcial queda visible y se retoma con2solicitudes,
2jobs y2decisiones únicas,0POST. FAIL907 corrigió el retry de decisión mediante
lectura de su recibo inmutable: mismo hash/actor/decisión/reason.4negativos
verificados sin escrituras. No reset ni actor fabricado.

Host:4módulos/2loops existentes, mTLS y shutdown conjunto. Policy/pack/guard
exigidos. Jobs de campaña no son consumidos con la extensión desactivada.
Rollback poblado rechazado,6/10/10intactos; base vacía down/up PASS.
Compatibilidad del recordatorio anterior:7jobs/3POST/6resultados PASS tras
la extensión. Boundaries/vet/build/fuzz2s y2perfiles exactos.

16nuevos AUTHORED y7revisados;1428salidas previas exactas. Fuente Meta,
SDKs y dependencias oficiales anteriores sin upgrade ni atribución local.
El JSON enlaza las seis evidencias T2805 previas y sus hashes; publishers
ML/Merchant seleccionados, otros product writes no elegidos NONE_WITH_REASON.
No se elimina ningún adapter de lectura/captación.

La conversión es una observación de los IDs/evidencia/timestamp originales,
no causalidad de marketing ni prueba de pago. Credenciales/cuentas y policy
configurada quedan para el usuario; toda la implementación seleccionada está
probada en fixtures. SCA/identidad integral es T2803, IA/evals T2807 y operación/
retención T2809; no se ocultan como falta de credenciales. Sigue T2803.
