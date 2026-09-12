# V366 — historial y entrenamiento: corrección de alcance

2026-09-09. Mantenimiento de contratos, por aclaración explícita del usuario.
Entrada177 validada por kit1.3.1 reconstruido, resume/plan y169hashes de owners.

## Antes y después

Antes: V288 propuso memoria privada, conocimiento revisado y evaluaciones;
TEST09 exigía corpus/canal real como setup del contrato local. En V365 el agente
pidió ese historial para destrabar la biblioteca. El usuario aclaró NEW sin
historial, EXISTING con historial y finalidad de entrenar un modelo con código
de élite. Pedir los chats aquí y sustituir entrenamiento por otros usos fue una
interpretación incorrecta, registrada en FAIL661.

Ahora: PROJECT_HISTORY_MODEL_TRAINING_CONTRACT.md es el owner normativo portable.
Los datos y los resultados del entrenamiento permanecen privados en el proyecto
consumidor. La biblioteca proporciona capacidad reusable, no un almacén de datos
de clientes ni un modelo entrenado global. NEW no exige historial. EXISTING
inventaría su fuente y prueba permiso/contrato/dataset antes de activar el job.
No hay proveedor, técnica, modelo base o presupuesto elegido por esta aclaración.

Se hace explícito el flujo: importación sin efectos → original privado → derivados
revisados → train/validation/eval separados → entrenamiento/adaptación explícitos
→ candidato → evaluación contra baseline → promoción/rollback → runtime.
RAG, memoria y evaluaciones son usos diferentes; no satisfacen por sí solos el
requisito de entrenamiento. Los originales conservan bytes dentro de la política
de retención; no se destruyen por redactar derivados ni se preservan ilimitadamente.

Código/modelo/datos tienen admisiones distintas. Reutilizar sólo código compatible
admitido, con fuente/revisión/licencia/SCA/pruebas; resolver gaps por G0–G8 antes
de implementar. Ninguna marca o sample oficial acredita la pipeline completa.
No se importó historial, subió información, entrenó un modelo, instaló una
dependencia ni generó código de entrenamiento durante V366.

## Cambios autorizados y alcance de verificación

REQ05/TEST09 y el riesgo/bloqueo asociados se corrigen por esta respuesta real,
con before178 intacto y delta JSON conservado. TEST09 sigue BLOCKED: faltan
implementación/admisión y pruebas reusables; no se presenta el contrato nuevo
como ejecutable. Las otras47filas y los estados de las48filas no cambian.
43PASS/5BLOCKED no es avance medido del esfuerzo. Ninguna tarea se cierra aquí.
Readiness documents.required=false ya correspondía al mantenimiento y no cambió.
La aclaración parcial no responde todas las rondas D–H ni autoriza datos reales.

El router de inicio, el mapa IA, la spec y los owners del roadmap apuntan al
contrato nuevo; V288/V289 quedan históricos. Sólo investigación/contratos de
esta capacidad pueden avanzar sin decisiones de un consumer. No volver a pedir
historial para preparar la biblioteca ni alegar que ya existe un trainer élite.

FAIL660 conserva una ruta de manual supuesta; se recuperó usando el nombre exacto
del mapa. FAIL661 conserva la interpretación previa y su corrección de alcance.
Los163Markdown de implementation_packs quedan byte-idénticos;162packs,
1461archivos materializables y53perfiles intactos. Inventario esperado805Markdown.
No se repite el Preflight171 de154checks para un cambio sin código; gate
estructural después del checkpoint178 y cierre179 con resume/plan.

Stage: %TEMP%/elite-v366-76e9e33167d2491ba5c6d493637ec149.
Autoridades consultadas: ML_PRODUCTION_LLMOPS_EVALUATION.md,
AI_SECURITY_GOVERNANCE_PRIVACY.md y contratos de admisión/documentos existentes.
No se selecciona ni se afirma vigencia de un servicio de entrenamiento.

| Recibo relativo al stage | SHA-256 |
|---|---|
| user-scope-clarification.json | `767006d7450e1896afb50e17e94da91e6c14d0238e24efd88d6114c938c2ab05` |
| contract-scope-delta.json | `278c01c20c2f4bcb547a5b27adcdfcff798f49a8e4884d7d7aaaf5897d74573e` |
| baseline177.json | `ba288454788e1afebfe3ab48dac1890ceb6e74226a9b9b18fcb316516e18077f` |
| resume177.log | `a9da687d8590b56280e66cf59fd503f5cfb0667ba73df4c088244db0d27cccea` |
| plan177.log | `2fd65a7ad27c8b968c64f8079abf8610894d6e7732b0b0f514dba962cf93b565` |

## Cierre179

VERIFY_LIBRARY178 PASS:162packs/1461archivos materializables/805Markdown/53perfiles.163Markdown de packs byte-idénticos y ambos verifiers intactos. Cambio autorizado sólo en REQ05/TEST09 y riesgo/assurance asociados: las otras47filas y todos los48estados intactos. No es un PASS de entrenamiento. Readiness JSON permanece igual;3routers enlazan el contrato. FAIL661 queda SCOPE_CORRECTION_VERIFIED por revisión del delta, referencias y gate estructural. Checkpoint179 se registra y se valida con resume/plan.

No pedir un canal o historial real para completar la preparación de biblioteca; los siguientes pasos son selección/admisión de código y diseño verificable del pipeline reusable bajo el propósito corregido. Modelo/proveedor/método/costo y calidad se decidirán en el consumer.

| Recibo de cierre relativo al stage | SHA-256 |
|---|---|
| verify-library178.log | `6714df83e9fc3fdfd5e9eb2db3b43e26abf00058caabea5fdc4fb33d5faae621` |
| scope-verification178.json | `c94c771075dc814f54d1df8e21a9ba25661c1af1ea0ca17ea800b32b06b1ffe7` |
| checkpoint178.log | `3b49db1083a810529f72c1ce6c5e96b05f54e6adebc00603a5cb13b39b4c14e6` |
| resume178.log | `db780c3c823888d86de1d63f00382e4fa364c0e3d0623c18e77bd357fda380c2` |
| plan178.log | `2fd65a7ad27c8b968c64f8079abf8610894d6e7732b0b0f514dba962cf93b565` |
