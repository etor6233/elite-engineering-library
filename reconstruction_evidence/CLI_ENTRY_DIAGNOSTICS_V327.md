# Entrada CLI y errores de continuidad — V327

Fecha: 2026-09-08. Scope: mantenimiento T2801/T2808; sin ampliación de producto.
Resultado focal: PASS. TEST-46 aprobado; TEST-08 permanece planificado para
usabilidad sobre el artefacto de release exacto. No sesión de usuarios, agentes
independientes, ahorro de tokens, aceptación de producto ni cierre de macrofrente.

## Defecto y cambio canónico

La guía ejecutada en NEW encontró FileNotFoundError/traceback/exit1 al validar
antes del primer checkpoint. El CLI manejaba StateError, pero errores OS y
Unicode escapaban. FAIL-20260908-537 y LIB-FAIL-2254 conservan la observación.
EXECUTION-VALIDATOR1.3.0→1.3.1 añade OSError/UnicodeError al catch de ambos CLI;
la API interna sigue propagando errores. No retry, append adicional ni reparación
automática. Schemas y formato de eventos no cambian; sólo3bloques de15 cambian.

Tres tests nuevos produjeron10fallos/4errores en1.3.0. Después:27/27tests PASS
(16de estado y11del contrato), con16subcasos subprocess y4errores de I/O
inyectados. Se verifica exit2, prefijo BLOCKED, ausencia de traceback y cero
mutación de archivos en cada negativo. El fixture del test no es una aprobación.

Reconstrucción limpia15/15 y composición del plan local23archivos/2packs PASS;
los15archivos de ejecución son byte-idénticos en ambos destinos. Versiones y
pins activos se sincronizan a1.3.1; reportes históricos1.3.0/24tests se conservan.
El runner mantiene un conteo exacto, ahora27, sin retirar aserciones.

## Guía comprobada

Owner: markdown_system/PROJECT_ENTRY_CLI_GUIDE.md0.1.0, enlazado desde README.
Prepara compositor/readiness/execution/bridge; deja el template BLOCKED y exige
respuestas reales. Explica recuperación y que staging no es una transacción
global. README distingue Python de tooling y actualiza inventarios contra V326
y V322. No se genera authority map ni checkpoint ficticio.

Se ejecutaron literalmente los dos bloques PowerShell extraídos de la guía:
13comandos (2preparaciones,2reentradas,2resume sin checkpoint,4negativos de
preparación,3resume válido/alterado/restaurado) y8consultas --help sobre4CLI en
NEW/EXISTING. Rutas con espacios/acentos; archivo previo BOM/CRLF byte-exacto e
instrucciones AGENTS previas preservadas como prefijo. Sin Git creado, red,
instalaciones, secretos o efectos de proveedor. Ante ausencia de checkpoint se
verifica BLOCKED, no se crea uno.

La reanudación positiva copia un fixture sintético previo de V326 de5archivos,
con2eventos válidos; jamás altera el original. Se conserva una copia de los bytes
corruptos antes de restaurar sólo el owner del fixture desde su snapshot.
Límite90s por comando y15s por help, sin retries. Primera ejecución conserva el
fallo real; su oracle de texto esperaba erróneamente STATE_FAIL y se corrigió a
STATE_BLOCKED, sin aceptar exit1 ni traceback. Una pasada verde intermedia se
retuvo; la final usa las fuentes cuyos hashes se listan abajo. Repeticiones no
se contabilizan como casos independientes adicionales.

Toolchains: PowerShell7.6.5 y Python3.14.4 del entorno local, identidad SHA en
environment.json. No runtimes nuevos. Autoridad de interacción:
FRONTEND_PRODUCT_ENGINEERING_UX.md §2; owners de ejecución: COMPOSITION_PROTOCOL,
PROJECT_START_READINESS_GATE y ENGINEERING_EXECUTION_VALIDATOR. No fuente pública
copiada ni nuevo claim atribuido a MIT/Microsoft/Google/NASA.

## Evidencia conservada

Stage identificado por elite-v327-bf63a9a8a8944726b0fa9cd942f1bbc1 en el directorio
temporal de esta sesión; el puntero elite-v327-current.txt conserva su ruta local.
No incorporar staging ni los expedientes locales al payload portable.

| Artefacto relativo al stage | SHA-256 |
|---|---|
| red-cli-tests.log | 39f3808f7c6445df6fc9039cfc48587ebd0fc67e2fca24385d1e138b1f1a3616 |
| green-kit-tests.log | d9db1ae3308f77283e9254a54db29ed2997eebab96594f14e2b3ee7d2ae3c5f9 |
| block-delta.txt | 301262c89a20433560fbaf823fbcfab2e28e45a839f7ab37921a251f65c3bee0 |
| parity.json | a727ee93570196128a6f8cab27516f98af252762c4322fe8cac46b8f85c2e6c7 |
| guide-final/results.json | e8138f81799b91b027a4c31113cec84d1106c7422e47d8b7ec63bb193d787f80 |
| guide-final/environment.json | 333880a498a0420fb6fd986ce331d2f72726167a2d3f07ab397562245752f584 |
| guide-final/probe_guide.py | eb572f298ddd53925804ecee0cae0b508f7f9dfbfa662c3c0f82ed4b296a4b65 |

| Fuente del ensayo final | SHA-256 |
|---|---|
| markdown_system/PROJECT_ENTRY_CLI_GUIDE.md | 671fa99c1bdac8d11caca0690fb1a0072d4bb751330c17b641b044b2aa3e9e65 |
| INSTALL_AGENT_BRIDGE.ps1 | fe32bf216ea48e9f218dbcf3075161bef445a82bd498bea3f14429d172a12a4b |
| materialize_markdown_pack.ps1 | 7ae045ff9b03ea22242d376d16b2a114f00e2f095217789eb92ca1f19e126356 |
| implementation_packs/MARKDOWN_COMPOSITOR_CORE.md | a8505384bb3f9085d36f21f1b6471fe24cbd17b054fa87f82f1d24b7007ce9d5 |
| implementation_packs/PROJECT_START_READINESS_VALIDATOR.md | f0e132178236937e676fb829050745dadcb8ece60f778bb6b3e5c9a1a3e8fe79 |
| implementation_packs/ENGINEERING_EXECUTION_VALIDATOR.md | 368c4f9d0c328353634e3dc0738e0d300f225013b9c7a769628b971d0c854808 |
| markdown_system/PROJECT_READINESS_GATE_PACK_PLAN.md | c4d69ba7893d399c5c697fc2335675df6c3629f7601350d856705167be5360a2 |

## Reproducción

Materializar el pack1.3.1 en destino ausente, fijar ELITE_AUTHORITY_ROOT a una
copia de esta biblioteca y ejecutar Python -X utf8 -B -m unittest discover -v
desde engineering_execution_kit. El pack contiene los tres tests nuevos.
Para el recorrido, guardar el programa siguiente en un staging nuevo, definir
ELITE_LIBRARY_ROOT, ELITE_PWSH_EXECUTABLE y ELITE_RESUME_FIXTURE. Esta última ruta
debe contener sólo un fixture sintético con checkpoint válido; puede crearse con
ExecutionStateTests.setUp/_base('EXISTING') y append_checkpoint del kit,
preservándolo en un destino exclusivo antes del ensayo. No usar datos ni estado
de un proyecto real. Los SHA del ensayo histórico identifican sus bytes, no
prometen timestamps o tiempos idénticos en otro host.

```python
import hashlib
import json
import os
from pathlib import Path
import re
import shutil
import subprocess
import sys
import time

root = Path(os.environ['ELITE_LIBRARY_ROOT']).resolve()
stage = Path(__file__).resolve().parent
pwsh = Path(os.environ['ELITE_PWSH_EXECUTABLE']).resolve()
python = Path(sys.executable).resolve()
guide = root / 'markdown_system/PROJECT_ENTRY_CLI_GUIDE.md'
blocks = re.findall(r'```powershell\r?\n(.*?)```', guide.read_text(encoding='utf-8'), re.S)
assert len(blocks) == 2
prepare = stage / 'prepare-from-guide.ps1'
prepare.write_text(blocks[0], encoding='utf-8', newline='\n')
resume = stage / 'resume-from-guide.ps1'
resume.write_text('param([string]$PythonExecutable,[string]$ExecutionKit,[string]$ProjectRoot)\n' + blocks[1] + '\nexit $LASTEXITCODE\n', encoding='utf-8', newline='\n')
results = []

def digest(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()

def inventory(path):
    return {p.relative_to(path).as_posix(): digest(p) for p in sorted(path.rglob('*')) if p.is_file()}

def run(name, args, expected, marker):
    start = time.perf_counter()
    completed = subprocess.run([str(pwsh), '-NoProfile', '-File', *map(str, args)], capture_output=True, timeout=90)
    (stage / (name + '.stdout.txt')).write_bytes(completed.stdout)
    (stage / (name + '.stderr.txt')).write_bytes(completed.stderr)
    text = (completed.stdout + completed.stderr).decode('utf-8', errors='replace')
    row = {'name': name, 'exit': completed.returncode, 'expected': expected, 'marker': marker, 'elapsed_ms': round((time.perf_counter()-start)*1000, 3)}
    results.append(row)
    (stage / 'progress.json').write_text(json.dumps(results, indent=2), encoding='utf-8')
    assert completed.returncode == expected, (name, text[-3000:])
    assert marker in text, (name, text[-3000:])
    return row

source_paths = ['markdown_system/PROJECT_ENTRY_CLI_GUIDE.md', 'INSTALL_AGENT_BRIDGE.ps1', 'materialize_markdown_pack.ps1', 'implementation_packs/MARKDOWN_COMPOSITOR_CORE.md', 'implementation_packs/PROJECT_START_READINESS_VALIDATOR.md', 'implementation_packs/ENGINEERING_EXECUTION_VALIDATOR.md', 'markdown_system/PROJECT_READINESS_GATE_PACK_PLAN.md']
source_hashes = {name: digest(root/name) for name in source_paths}
runtime_hashes = {str(p): digest(p) for p in (pwsh, python)}
env = {'python': sys.version, 'runtimes': runtime_hashes, 'source_hashes': source_hashes, 'guide_steps_extracted': 2, 'timeout_seconds_per_command': 90, 'network': 'none', 'agents_or_human_usability_sessions': 0}
(stage/'environment.json').write_text(json.dumps(env, indent=2), encoding='utf-8')

def prep_args(project, tooling, **overrides):
    values = {'LibraryRoot': root, 'ProjectRoot': project, 'StageRoot': tooling, 'PowerShellExecutable': pwsh, 'PythonExecutable': python}
    values.update(overrides)
    return [prepare, *[item for k, v in values.items() for item in ('-'+k, v)]]

for mode in ('NEW', 'EXISTING'):
    case = stage / ('caso ' + mode + ' ñ')
    project = case/'proyecto á'
    tooling = case/'herramientas ü'
    project.mkdir(parents=True)
    tooling.mkdir()
    if mode == 'EXISTING':
        (project/'usuario.txt').write_bytes(b'\xef\xbb\xbfOriginal sin cambios\r\n')
        (project/'AGENTS.md').write_text('# Instrucciones previas\nConservar el proyecto del usuario.\n', encoding='utf-8')
    baseline = inventory(project)
    original_agents = (project/'AGENTS.md').read_bytes() if mode == 'EXISTING' else b''
    (case/'baseline.json').write_text(json.dumps(baseline, indent=2), encoding='utf-8')
    run(mode+'-prepare', prep_args(project, tooling), 0, 'ENTRY_PREPARED_READINESS_BLOCKED')
    assert not (project/'PROJECT_EXECUTION_STATE.json').exists()
    assert not (project/'PROJECT_EXECUTION_EVENTS.jsonl').exists()
    report = json.loads((project/'PROJECT_READINESS_REPORT.json').read_text(encoding='utf-8'))
    assert report['status'] == 'BLOCKED' and report['errors']
    if mode == 'EXISTING':
        assert digest(project/'usuario.txt') == baseline['usuario.txt']
        assert (project/'AGENTS.md').read_bytes().startswith(original_agents)
    before = inventory(case)
    run(mode+'-repeat', prep_args(project, tooling), 1, 'ENTRY_DESTINATION_OCCUPIED')
    assert inventory(case) == before
    kit = tooling/'execution/engineering_execution_kit'
    run(mode+'-resume-before-checkpoint', [resume, '-PythonExecutable', python, '-ExecutionKit', kit, '-ProjectRoot', project], 2, 'EXECUTION_STATE_BLOCKED')
    assert inventory(case) == before
    gate = tooling/'readiness/project_readiness_gate'
    for script in (gate/'render_project_advisory.py', gate/'validate_project_readiness.py', kit/'validate_execution_state.py', kit/'checkpoint_execution_state.py'):
        help_result = subprocess.run([str(python), '-X', 'utf8', '-B', str(script), '--help'], capture_output=True, timeout=15)
        assert help_result.returncode == 0 and b'usage:' in help_result.stdout
        (case/(script.stem+'.help.txt')).write_bytes(help_result.stdout)
    results.append({'name': mode+'-help', 'scripts': 4, 'status': 'PASS'})

negative = stage/'negative'
negative.mkdir()
np = negative/'project'
nt = negative/'tooling'
np.mkdir()
nt.mkdir()
for name, override, marker in [
    ('missing-project', {'ProjectRoot': negative/'absent-project'}, 'ENTRY_DIRECTORY_MISSING'),
    ('missing-python', {'PythonExecutable': negative/'absent-python.exe'}, 'ENTRY_TOOL_MISSING'),
    ('missing-pwsh', {'PowerShellExecutable': negative/'absent-pwsh.exe'}, 'ENTRY_TOOL_MISSING'),
]:
    before = inventory(negative)
    run(name, prep_args(np, nt, **override), 1, marker)
    assert inventory(negative) == before
(np/'PROJECT_ADVISORY_A.md').write_bytes(b'\xef\xbb\xbfUser-owned advisory\r\n')
before = inventory(negative)
run('existing-advisory', prep_args(np, nt), 1, 'ENTRY_DESTINATION_OCCUPIED')
assert inventory(negative) == before

fixture = stage/'resume fixture'
shutil.copytree(Path(os.environ['ELITE_RESUME_FIXTURE']), fixture)
before = inventory(fixture)
kit = stage/'caso NEW ñ/herramientas ü/execution/engineering_execution_kit'
run('resume-valid', [resume, '-PythonExecutable', python, '-ExecutionKit', kit, '-ProjectRoot', fixture], 0, 'EXECUTION_STATE_PASS')
assert inventory(fixture) == before
owner = fixture/'PROJECT_FAILURE_LESSONS.md'
original = owner.read_bytes()
owner.write_bytes(original + b'\nSynthetic tamper\n')
tampered = inventory(fixture)
run('resume-tampered', [resume, '-PythonExecutable', python, '-ExecutionKit', kit, '-ProjectRoot', fixture], 2, 'evidence hash mismatch')
assert inventory(fixture) == tampered
(stage/'tampered-owner-preserved.bin').write_bytes(owner.read_bytes())
owner.write_bytes(original)
run('resume-restored', [resume, '-PythonExecutable', python, '-ExecutionKit', kit, '-ProjectRoot', fixture], 0, 'EXECUTION_STATE_PASS')
assert inventory(fixture) == before
assert source_hashes == {name: digest(root/name) for name in source_paths}
(stage/'results.json').write_text(json.dumps({'status':'PASS','cases':results,'source_hashes':source_hashes,'limits':['CLI preparation only','no release archive','no human or independent agent usability session','no actual intake answers or new checkpoint created','no concurrency or power-loss proof']}, indent=2), encoding='utf-8')
print('GUIDE_CLI_PASS commands=13 help_checks=8 source_files=7 NEW_EXISTING_preservation=true')
```

## Cierre integrado pendiente de ejecución

Después del checkpoint85, ejecutar Preflight sobre este head, que incluye
VERIFY_LIBRARY y las27regresiones. Guardar preflight85.json/preflight85.log;
su estado puede seguir BLOCKED por Docker o gates target pese a pasos locales
PASS. No marcar el roadmap ni readiness completos por este arreglo.

## Resultado integrado observado

Preflight85 terminó exit0,151/151pasos PASS; VERIFY_LIBRARY_PASS161packs/1453archivos/763Markdown/52perfiles y franquicia67/746. La suite del kit1.3.1 aparece como engineering-execution-syntax-and-27-regressions PASS. Disponibilidad7/8; Docker ausente mantiene statusBLOCKED. allow_network=false; no nuevas cuentas, instalaciones o efectos. Estado del proyecto DISCOVERY/BLOCKED y TEST08 siguen intactos.

SHA-256 preflight85.json: `4aeb86e2599304494cbf32c12b0ebe7dce63249ff94f2d03373c7e3e1bcc9a58`.

SHA-256 preflight85.log: `b4519a9bf298e39c87a35d559540d23fe576506a8c86d4ee84afe879575b5952`.

Los apartados anteriores conservan el orden histórico: el cierre ya no está pendiente. Checkpoint86 añade sólo evidencia/continuidad después del gate completo; no cambia código, pins, guías ni tests. Se validan de nuevo contrato y cadena sin repetir suites sin delta. Los85eventos previos se preservan byte a byte.
