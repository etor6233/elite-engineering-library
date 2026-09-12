# Maintenance readiness boundary — V282

Fecha: 2026-09-07. Scope: T2801, mantenimiento de biblioteca, no producto.
Procedencia: control y pruebas AUTHORED. No son código publicado por Microsoft o GitHub.

## Cambio ejecutado

VERIFY_LIBRARY.ps1 y CREATE_PORTABLE_ARCHIVE.ps1 reconocen los mismos veinte
nombres locales exactos de readiness, advisory y assurance, y cuatro paths
Spec Kit exclusivos de library-maintenance. Exigen el par de checkpoint.
Excluyen esos registros del payload, conservan el rechazo de archivos
inesperados y comprueban reparse points antes de descender por cada directorio.
No usan un wildcard PROJECT_*, no excluyen otros proyectos y no modifican packs.

La frontera se apoya en la semántica oficial de
[Microsoft FileAttributes.ReparsePoint](https://learn.microsoft.com/en-us/dotnet/api/system.io.fileattributes?view=net-10.0).
El proceso de spec/plan/tasks mantiene el flujo de
[GitHub Spec Kit](https://github.github.com/spec-kit/); no se actualizó ni ejecutó
una nueva versión de Spec Kit por consultar su documentación.
La política de distribución es propia de Elite, no una recomendación atribuida
a esas empresas ni una defensa contra mutaciones concurrentes maliciosas del workspace.

## Verificación ejecutada

- 32 checks PASS sobre las funciones reales de ambos scripts, cargadas mediante AST.
- Casos: distribución sin estado, par correcto, par incompleto, archivo convertido
  en directorio, nombre aproximado, reparse raíz, intake válido/sin checkpoint,
  intake como directorio, advisory fuera de A–H, cuatro archivos spec válidos,
  archivo inesperado dentro de specs, junction anidada, specs como archivo y
  carpeta de nombre aproximado.
- Los positivos seleccionan exclusivamente README sintético: ningún registro local
  entra en la lista de archivos a comprimir. No se creó un ZIP final.
- PROJECT-START-READINESS-VALIDATOR 0.6.0 materializado en staging: ocho archivos,
  62 regresiones PASS. No se modificó su código ni se debilitó su gate.

## Expediente real, no fixture

Se crearon desde los templates canónicos PROJECT_READINESS_RECORD.md y
PROJECT_READINESS_GATE.json para elite-library-maintenance. Se generó la ronda A
con el renderer exacto; no se marcaron rondas como contestadas por presencia.
La ejecución real generó PROJECT_READINESS_REPORT.json: BLOCKED, 102 observaciones.
Son observaciones del baseline todavía incompleto, no 102 defectos nuevos de producto
ni una estimación de esfuerzo. Incluyen clasificación de 48 capacidades, owners,
fuentes/contratos y pruebas enlazadas; deben reconciliarse con la evidencia existente.

PROJECT_ADVISORY_A.md tiene SHA-256
2a04dd2e301ead57c7c7593fa056419b7706026b1f171a1d777b638260163b8e.
El baseline es local y excluido de distribución. PROJECT_FAILURE_LESSONS.md conserva
FAIL-384 abierto: esta corrección elimina la incompatibilidad del empaquetador,
pero aún faltan el contenido de readiness y implementation_assurance demostrados.

## Continuación exacta

Completar T2801 desde las decisiones ya expresadas y los artefactos observados,
sin pedir cuentas de una franquicia imaginaria para probar herramientas locales.
La cobertura requerida del perfil integral no se excluye por ese motivo.
Generar y explicar las rondas pendientes una por vez, enlazar fuentes y aceptación,
conservar el roadmap como único backlog y no expandir producto hasta el gate real.
No repetir V281: su comparación de logging, tests y límites permanecen vigentes.

Staging de herramientas/logs: elite-v282-294334f1e07f40dda21adbf118e779e3,
bajo el temporal del sistema. Regresión del selector:
elite-state-distribution-774e335986804d79a9680113ba152bec.

Comandos reproducibles: markdown_system/test_library_maintenance_state.ps1;
materialize_markdown_pack.ps1 sobre PROJECT_START_READINESS_VALIDATOR.md;
python -m unittest discover sobre el gate materializado; render_project_advisory.py
para A; validate_project_readiness.py sobre la raíz real. Los registros locales
no heredan un PASS del test harness y la distribución no hereda sus decisiones.
