# V402 — preparación explícita de infraestructura de biblioteca

Estado de esta evidencia: PREPARATION_PROVEN, únicamente para iniciar/corregir
el trabajo local definido por el usuario. No cierra T2801 final, TEST02/03/07,
T2802–T2810 ni ARCA. No declara READY_FOR_LIBRARY_USE ni producción live.

## Decisión y autoridad

La decisión es reconstruction_evidence/LIBRARY_INFRA_SCOPE_V402.md. El usuario
autorizó infraestructura completa con fixtures y sin secretos; AUTHORED sólo
glue inevitable. PROJECT continúa con sus rondas A–H y estado histórico propio.
El gate de biblioteca se conserva en PROJECT_LIBRARY_READINESS_GATE.json para
no cambiar ni heredar el gate de producción/consumer.

## Preparación demostrada

PROJECT-START-READINESS-VALIDATOR 0.7.0: 102 pruebas (75 anteriores/27 nuevas)
PASS y 11 archivos reconstruidos desde la fuente canónica con bytes idénticos
al candidato probado. Se prueban niveles, deferrals exactas, evidencia SHA,
inventario completo/rebuild independiente, G0–G8 y rechazo de AUTHORED de negocio.
Este código es AUTHORED glue de validación. No se atribuye a empresa externa.
Los logs y el inventario exacto están en LIBRARY_READINESS_TESTS_V402.log y
LIBRARY_READINESS_RECONSTRUCTION_V402.json, junto a este documento.

Tooling usado aquí: CPython local y materializador PowerShell del workspace.
No se ejecutó pnpm, provider, red, proyecto de negocio ni acceso live. El estado
pnpm V401 se referencia para conservar la contención y sus límites, no para
admitir el runtime general o inferir SCA completa por 475 identidades.

## Plan de resolución explícito — blockers OPEN

- SOURCE_PROVENANCE: inspeccionar los owners Commerce, Accounting, Journey,
  Surveys y el resto del perfil. El AUTHORED de reglas de negocio no satisface
  el criterio9. Se mantiene esa clasificación y se resuelve cada claim mediante
  fuentes oficiales fijadas/admisión G0–G8 o una decisión visible del owner;
  nunca cambiar a ADAPTED por cita o reputación. Afecta T2801/T2802.
- DEPENDENCY_ADMISSION: conservar el pnpm original bloqueado y la contención
  aislada. Cerrar la identificación real de dependencias/source/nativos fuera
  de Daybreak, licencias, notices y SCA sobre la composición final. No ejecutar
  un artefacto rechazado ni tomar los hashes como aprobación. Afecta T2803/T2810.
- MISSING_EVIDENCE: cerrar integración y pruebas faltantes en el orden del
  usuario: T2802/TEST02, T2804, T2805 infraestructura, T2803 sin Daybreak,
  T2806/T2807, T2808, T2809 local, T2810 portable; ARCA infraestructura penúltima.
  Reusar evidencia previa sólo con bytes/claim compatibles. Los controles
  pendientes no quedan completos por este validator ni por tests sintéticos.

## Daybreak último diferido

Exclusión única daybreak-libxml2 / T2803: ACCESS_BLOCKED. No se investigó,
ejecutó ni intentó desbloquear. Trigger: nuevo encargo explícito del usuario
para esa rama y acceso permitido; sin ambos se conserva excluida y no se
transfiere la exclusión a otras dependencias o fallos de seguridad.

## Uso y límite

Preparación permite corregir los blockers anteriores con el plan visible.
Release los rechaza mientras sigan OPEN. Todos los adapters deben aportar
contratos, código, fixtures, wiring y reconciliación completos antes de recibir
USER_CREDENTIALS. Las decisiones/corpus/target no se disfrazan de credenciales.
Las pruebas del validator prueban el gate; no hacen verdadera una afirmación
de negocio. El agente debe revisar las evidencias semánticas y la procedencia.
