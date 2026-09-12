# Entrada CLI y recuperación de un proyecto

Versión de guía: 0.1.2. Scope: mantenimiento y preparación local NEW/EXISTING.
Procedencia del procedimiento: AUTHORED; invoca los owners canónicos, sin código
de producto ni atribución a un proveedor. No completa respuestas, authority map,
engineering contract o checkpoint por el usuario. TEST-08 exige recorrer esta guía
desde un candidato portable completo y exacto antes de publicar el release. El
procedimiento aislado no cierra ese gate; conservar SHA del ZIP y de las fuentes
consumidas. TEST-07 exige reutilizar ese payload validado o revalidar la aceptación
si el candidato cambia; un PASS de instalación no concede firma, SCA ni admisión.

## Preparación

Usar una copia verificada de la biblioteca. Se necesitan PowerShell 7 y Python
para los validadores, incluso cuando todavía no hay producto. Pasar rutas exactas
a los ejecutables disponibles; no instalar ni actualizar runtimes por este paso.
Los packs fijan sus requisitos; consultar el inventario de toolchains con
`VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Toolchain` si falta una herramienta.

Crear previamente el directorio del proyecto y un directorio de staging separado.
El agente inspecciona el proyecto antes de escribir: NEW requiere baseline vacío
o scaffold declarado; EXISTING requiere inventario y delta autorizado. Conservar
ese inventario fuera de los archivos que se van a crear. El staging no es un
segundo repositorio fuente: contiene sólo tooling reconstruible y sus receipts.

Guardar el siguiente bloque como script temporal fuera de la biblioteca e
invocarlo con `-LibraryRoot`, `-ProjectRoot`, `-StageRoot`, `-PowerShellExecutable`
y `-PythonExecutable`. Las cinco rutas deben existir. El procedimiento exige
destinos nuevos y rechaza una reentrada antes de escribir; para un proyecto ya
inicializado usar la sección de reanudación.

```powershell
# ELITE_CLI_STEP: prepare
#requires -Version 7.0
param(
  [Parameter(Mandatory)][string]$LibraryRoot,
  [Parameter(Mandatory)][string]$ProjectRoot,
  [Parameter(Mandatory)][string]$StageRoot,
  [Parameter(Mandatory)][string]$PowerShellExecutable,
  [Parameter(Mandatory)][string]$PythonExecutable
)
$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest
foreach ($directory in @($LibraryRoot, $ProjectRoot, $StageRoot)) {
  if (-not (Test-Path -LiteralPath $directory -PathType Container)) {
    throw "ENTRY_DIRECTORY_MISSING: crear o corregir la ruta: $directory"
  }
}
foreach ($executable in @($PowerShellExecutable, $PythonExecutable)) {
  if (-not (Test-Path -LiteralPath $executable -PathType Leaf)) {
    throw "ENTRY_TOOL_MISSING: indicar ejecutable disponible: $executable"
  }
}
$LibraryRoot = (Resolve-Path -LiteralPath $LibraryRoot).Path
$ProjectRoot = (Resolve-Path -LiteralPath $ProjectRoot).Path
$StageRoot = (Resolve-Path -LiteralPath $StageRoot).Path
$compositor = Join-Path $StageRoot 'compositor'
$execution = Join-Path $StageRoot 'execution'
$readiness = Join-Path $StageRoot 'readiness'
$gate = Join-Path $readiness 'project_readiness_gate'
foreach ($destination in @($compositor, $execution, $readiness,
    (Join-Path $ProjectRoot 'PROJECT_READINESS_GATE.json'),
    (Join-Path $ProjectRoot 'PROJECT_ADVISORY_A.md'),
    (Join-Path $ProjectRoot 'PROJECT_READINESS_REPORT.json'),
    (Join-Path $ProjectRoot 'PROJECT_EXECUTION_STATE.json'),
    (Join-Path $ProjectRoot 'PROJECT_EXECUTION_EVENTS.jsonl'))) {
  if (Test-Path -LiteralPath $destination) {
    throw "ENTRY_DESTINATION_OCCUPIED: conservar $destination; usar reanudacion o staging nuevo"
  }
}
function Invoke-EntryTool([string]$Tool, [string[]]$Arguments) {
  & $Tool @Arguments
  if ($LASTEXITCODE -ne 0) { throw "ENTRY_TOOL_FAILED: exit=$LASTEXITCODE; conservar salidas y corregir el diagnostico anterior" }
}
$materializer = Join-Path $LibraryRoot 'materialize_markdown_pack.ps1'
Invoke-EntryTool $PowerShellExecutable @('-NoProfile','-File',$materializer,'-PackFile',(Join-Path $LibraryRoot 'implementation_packs/MARKDOWN_COMPOSITOR_CORE.md'),'-Destination',$compositor)
Invoke-EntryTool $PowerShellExecutable @('-NoProfile','-File',(Join-Path $compositor 'tools/compose-markdown-project.ps1'),'-LibraryRoot',$LibraryRoot,'-PlanFile',(Join-Path $LibraryRoot 'markdown_system/PROJECT_READINESS_GATE_PACK_PLAN.md'),'-Destination',$readiness)
Invoke-EntryTool $PowerShellExecutable @('-NoProfile','-File',$materializer,'-PackFile',(Join-Path $LibraryRoot 'implementation_packs/ENGINEERING_EXECUTION_VALIDATOR.md'),'-Destination',$execution)
Invoke-EntryTool $PowerShellExecutable @('-NoProfile','-File',(Join-Path $LibraryRoot 'INSTALL_AGENT_BRIDGE.ps1'),'-LibraryRoot',$LibraryRoot,'-ProjectRoot',$ProjectRoot,'-Agent','Both')
Copy-Item -LiteralPath (Join-Path $gate 'project-readiness.template.json') -Destination (Join-Path $ProjectRoot 'PROJECT_READINESS_GATE.json')
Invoke-EntryTool $PythonExecutable @('-X','utf8','-B',(Join-Path $gate 'render_project_advisory.py'),'--project-root',$ProjectRoot,'--round','A','--output','PROJECT_ADVISORY_A.md')
& $PythonExecutable -X utf8 -B (Join-Path $gate 'validate_project_readiness.py') --project-root $ProjectRoot --report PROJECT_READINESS_REPORT.json
if ($LASTEXITCODE -ne 2) { throw "ENTRY_UNEXPECTED_READINESS: el template intacto debe devolver exit 2/BLOCKED" }
$report = Get-Content -LiteralPath (Join-Path $ProjectRoot 'PROJECT_READINESS_REPORT.json') -Raw -Encoding utf8 | ConvertFrom-Json
if ($report.status -cne 'BLOCKED' -or $report.errors.Count -eq 0) { throw 'ENTRY_INVALID_REPORT: conservar el reporte y revisar el gate' }
Write-Output 'ENTRY_PREPARED_READINESS_BLOCKED: abrir el agente en ProjectRoot y resolver la ronda A real'
```

El éxito de este procedimiento significa tooling preparado y readiness bloqueado
correctamente. Leer `PROJECT_ADVISORY_A.md`, explicar la ronda y registrar la
respuesta real antes de avanzar a B. Seguir `PROJECT_START_READINESS_GATE.md` y
`ENGINEERING_EXECUTION_PLAYBOOK.md` para los expedientes restantes; no cargar los
manuales de proveedores que el alcance no seleccionó.

## Reanudación y diagnóstico

Conservar las rutas usadas. `$ExecutionKit` debe apuntar a
`StageRoot/execution/engineering_execution_kit`; `$ProjectRoot` es el proyecto
original y `$PythonExecutable` el ejecutable exacto. Abrir el agente en esa raíz.
Si todavía no existe un checkpoint, completar primero su template con baseline,
delta, owners y evidencia reales siguiendo `EXECUTION_STATE_PROTOCOL.md` del kit.
No crear un checkpoint sintético para superar el gate.

```powershell
# ELITE_CLI_STEP: resume
& $PythonExecutable -X utf8 -B (Join-Path $ExecutionKit 'validate_execution_state.py') (Join-Path $ProjectRoot 'PROJECT_EXECUTION_STATE.json') --project-root $ProjectRoot --events (Join-Path $ProjectRoot 'PROJECT_EXECUTION_EVENTS.jsonl') --level resume
```

Exit 0/`EXECUTION_STATE_PASS` permite leer `next_action` y los
`context.must_read_refs`; no convierte DISCOVERY/BLOCKED en READY_TO_BUILD.
Exit 2 exige conservar estado, eventos y diagnóstico antes de corregir. Ante
hash distinto, localizar el cambio real y su owner; nunca recalcular hashes para
ocultar un cambio desconocido. Ante una respuesta perdida al guardar, validar
primero si el evento ya está completo; no repetir el append a ciegas.

| Diagnóstico | Acción que conserva evidencia |
|---|---|
| herramienta o directorio ausente | corregir la ruta exacta y volver a intentar; no instalar automáticamente |
| destino de tooling ocupado o parcial | conservarlo; reconstruir en un staging nuevo y comparar con el pack canónico |
| readiness BLOCKED | resolver la siguiente ronda y sus pruebas reales; emitir el siguiente reporte en una ruta nueva |
| bridge rechazado | leer el error de markers, Skill ajena, tamaño o ruta; preservar instrucciones existentes y corregir sólo el conflicto propio |
| checkpoint/hash/cadena inválidos | conservar bytes; aplicar AGENT_ERROR_RECOVERY_PROTOCOL y FAILURE_LEARNING_CONTRACT antes de mutar |

La preparación no es una transacción global del árbol. Un fallo puede dejar
tooling o bridge completos y un expediente incompleto; conservar esas salidas y
reanudar el paso faltante después de inspeccionarlas. Los writers concurrentes y
el corte eléctrico no están cubiertos. No borrar ni mover un proyecto para
reutilizar el comando inicial. La recuperación del release distribuido y de
datos productivos conserva sus gates separados.

Readiness0.6.2: inputs inválidos y errores de I/O devuelven exit2 y
PROJECT_READINESS_BLOCKED en stderr, sin fabricar reporte. La publicación de un
reporte completo exige hard links; si el filesystem no los permite, conservar
el diagnóstico y usar un destino compatible autorizado. Nunca reemplazar un
reporte existente como fallback. La publicación de reportes tiene prueba de
writers concurrentes; la preparación completa del proyecto y el corte eléctrico
conservan los límites anteriores.
