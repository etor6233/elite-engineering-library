# V279 — continuidad local y frontera de distribución

Fecha: 2026-09-07. Clasificación: mantenimiento de biblioteca.
Baseline: V278, 160 packs / 1433 archivos materializables; perfil de franquicia
67 packs / 742 archivos / 52 migraciones. Ningún bloque de esos packs cambia
en esta ronda. El nuevo documento eleva el inventario a 712 Markdown.

## Hallazgo y alcance

No se encontraron PROJECT_EXECUTION_STATE.json ni PROJECT_EXECUTION_EVENTS.jsonl
en la raíz al comenzar. El protocolo los exige, pero las dos allowlists de
distribución los rechazaban como archivos inesperados. Los 24 tests históricos
del execution validator no demostraban un checkpoint de este mantenimiento.
Se registra FAIL-20260907-384; no se deduce corrupción o mezcla de V278.

## Corrección de controles, no ampliación de producto

- VERIFY_LIBRARY.ps1 y CREATE_PORTABLE_ARCHIVE.ps1 reconocen sólo el par exacto
  de archivos locales regulares; ausencia de ambos permite distribuir la
  biblioteca, presencia de uno solo o reparse/directory no.
- Ambos selectores excluyen el par del payload y conservan el rechazo de
  cualquier archivo inesperado. El archive runner invoca el verificador antes
  de seleccionar/copiar archivos. No se creó ni publicó un ZIP.
- Cuando el par existe, VERIFY_LIBRARY reconstruye EXECUTION-VALIDATOR 1.3.0
  y valida hashes y cadena con --level resume. Un PASS de estado BLOCKED
  no abre readiness ni promoción.
- test_library_maintenance_state.ps1 ejecuta las funciones reales extraídas
  de los dos scripts: ausencia, par, sólo state, sólo events, directorio,
  nombre casi igual y reparse. No reimplementa los selectores como oráculo.
- El primer checkpoint captura EXISTING / DISCOVERY / BLOCKED, con enlaces
  a owners observados. No retrofecha eventos ni inventa pasos completados.

El cambio es control local AUTHORED, no nuevo código publicado por Google,
Meta o Microsoft. Reutiliza el validator canónico existente sin modificarlo.
Google SRE Release Engineering respalda el claim estrecho de identificar
artefactos, conservar registros y probar el contenido realmente distribuido;
no prescribe estos nombres, scripts ni garantiza el resultado de Elite.
Fuente oficial consultada el 2026-09-07:
<https://sre.google/sre-book/release-engineering/>.

## Verificación ejecutada y límites

Staging conservado: elite-v279-7fc43761727a4ab393d489f95bd482ba.

- Materialización limpia de EXECUTION-VALIDATOR: 15 archivos, hashes canónicos.
- Python: 24/24 tests PASS con ELITE_AUTHORITY_ROOT real; execution-tests.log.
- Distribución: 14/14 casos PASS sobre selectores de ambos scripts, incluyendo
  junction/reparse Windows. Fixtures sintéticos conservados; no datos reales.
- Checkpoint real: primera revisión y evento BASELINE_CAPTURED pasan --level
  resume en EXISTING / DISCOVERY / BLOCKED, seis owners por hash. El segundo
  checkpoint conserva el primero y actualiza el expediente al cierre.
- VERIFY_LIBRARY_PASS: 160 packs / 1433 archivos / 712 Markdown / 51 perfiles;
  valida además el checkpoint real. verify-library.log. Es materialización y
  trazabilidad, no ejecución de todos los builds, dependencias ni journeys.
- Preflight: BLOCKED, cuatro herramientas no detectadas (Go, dotnet, Docker,
  psql); exit 0 no sustituye el JSON. preflight.log/json. Go y psql usados en
  V278 permanecen en rutas absolutas observadas, fuera del descubrimiento;
  no se instalaron ni ejecutaron servicios para convertir el reporte en verde.
- No se repitieron los tests Go/PostgreSQL/browser de producto V278. Su evidencia
  permanece histórica; no se suma como nueva ni se extiende a live.

No se publica una release y no se modifica ningún pack de producto. El estado
local se valida con el kit reconstruido en execution/engineering_execution_kit.
Para reanudar en otra máquina hay que materializar el mismo pack; su ruta de
staging no forma parte del contrato portable ni de los hashes de evidencia.

## Continuación exacta

El roadmap canónico sigue siendo REUSABLE_CODE_READINESS_ROADMAP.md, sección 9.
No se crea otro plan ni se cambia el orden funcional.

La ausencia del expediente de readiness/engineering contract/implementation_assurance
de mantenimiento permanece visible. Reconstruirlo desde decisiones y evidencia
existentes, registrar sólo incertidumbres reales y validar antes de ampliar
producto. No enlazar una nota, un capstone o un test del validator como si fuera
una aprobación. Este pendiente no se convierte en una exigencia de credenciales
productivas para trabajos de biblioteca que no las requieren.

Después del gate, el tramo funcional pendiente continúa siendo montaje del host,
identidad de servicio, reportes/alertas y recorridos de referencia completos,
reutilizando sus owners. No se afirma que estén implementados por existir
interfaces. Cuentas/corpus/negocio/seguridad/operación del target conservan sus
condiciones. No hay porcentaje ni plazo universal demostrado.
