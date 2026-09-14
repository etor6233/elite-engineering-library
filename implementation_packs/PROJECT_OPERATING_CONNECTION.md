# Project operating connection — complementary bridge

## 1. Metadata

```yaml
pack_id: "PROJECT-OPERATING-CONNECTION"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: RECONSTRUCTIBLE
  admission: CONDITIONED
claim: "Complementary instruction routing, pinned authority integrity, single-writer checkpoints around the existing execution validator, and bounded connection smoke; no product or vendor-platform admission"
stacks: ["PowerShell 7+", "Python 3.14+ standard library"]
compatible_with: ["INSTALL_AGENT_BRIDGE SHA fe32bf216ea48e9f218dbcf3075161bef445a82bd498bea3f14429d172a12a4b", "ENGINEERING_EXECUTION_VALIDATOR 1.3.1", "FRANCHISE_PROJECT_OPERATING_PROTOCOL 1.0.0"]
incompatible_with: ["editing V402 sealed artifacts", "silent library revision upgrade", "product readiness inferred from connection test", "unverified platform support"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["No upstream product code copied; official methodology references are in FRANCHISE_PROJECT_OPERATING_PROTOCOL.md"]
verified_at: "2026-09-13"
```

## 2. Applicability

Materialize into an absent consumer path `.elite/connection-tools`. The complementary installer invokes the original immutable bridge and appends an independent instruction block. Conditions: separate library/project trees, original bridge SHA, declared protocol, Python/PowerShell available, no conflicting unmanaged destinations. It configures the existing Codex/Claude entries; platform execution remains NOT_TESTED until separately observed. Generic/Grok loading uses the included template.

## 3. Architecture contract

AUTHORED orchestration and tests only. Single protocol owns the workflow; the project retains existing blueprint, authority map, pack plan, execution state/events, tasks and failure lessons. Connection receipt owns paths/hashes only. Standard-library file/process operations do not add a product dependency. The original bridge and execution writer are called unchanged; no corporate authorship is attributed to this glue. Lock is local process coordination, not multi-host consensus. Receipt hashes detect drift; they are not a security signature.

## 4. Exact file manifest

```text
CREATE operating_connection/connection.py
CREATE operating_connection/fixture_work.py
CREATE operating_connection/INSTALL_PROJECT_OPERATING_BRIDGE.ps1
CREATE operating_connection/PLATFORM_ADAPTER_TEMPLATE.md
CREATE operating_connection/smoke_connection.py
CREATE operating_connection/verify_smoke_receipt.py
```

## 5. Materialization blocks

### FILE: `operating_connection/connection.py`

```yaml
block_id: "PROJECT-OPERATING-CONNECTION:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local operating bridge, integrity and test orchestration; not upstream company code"
license: "LicenseRef-Workspace-Owner"
sha256: "af175b406a3e28a050f554f9f0ad15b22a2e978d41afef1018b27b8e080ff67b"
variables: []
secrets_allowed: false
```

````python
#!/usr/bin/env python3
"""AUTHORED operating glue. No vendor authorship or product admission claim."""
from __future__ import annotations

import argparse
from contextlib import contextmanager
import hashlib
import importlib
import json
import os
from pathlib import Path
import subprocess
import sys

VERSION = '1.0.0'
BASE_SHA = 'fe32bf216ea48e9f218dbcf3075161bef445a82bd498bea3f14429d172a12a4b'
PROTOCOL = 'markdown_system/FRANCHISE_PROJECT_OPERATING_PROTOCOL.md'
RECEIPT = '.elite/operating-connection.json'
ENTRY = 'PROJECT_AGENT_ENTRY.md'
BEGIN = '<!-- ELITE-PROJECT-OPERATING:BEGIN -->'
END = '<!-- ELITE-PROJECT-OPERATING:END -->'
BASE_FILES = ['AGENTS.md', 'CLAUDE.md', '.agents/skills/elite-engineering-library/SKILL.md', '.claude/skills/elite-engineering-library/SKILL.md']
AUTHORITY_PATHS = [PROTOCOL, 'AGENTS.md', 'AGENT_SYSTEM_START.md', 'INSTALL_AGENT_BRIDGE.ps1', 'ENGINEERING_EXECUTION_PLAYBOOK.md', 'markdown_system/PROJECT_START_READINESS_GATE.md', 'implementation_packs/ENGINEERING_EXECUTION_VALIDATOR.md']
# Generated from the unchanged, admitted EXECUTION-VALIDATOR pack by pack_source.py.
EXECUTION_TOOL_HASHES = {'validate_execution_state.py': '63b8404611790978318c25008d109895224ea1dc7510f4df32d84f284a682a03', 'checkpoint_execution_state.py': 'a935f94fa3902539f88bf0c3097949977ac5e4091293a5666f985f73c241dc08'}


class Blocked(ValueError):
    pass


def require(condition, message):
    if not condition:
        raise Blocked(message)


def digest(raw):
    return hashlib.sha256(raw).hexdigest()


def checked(path, *, exists=True):
    path = Path(os.path.abspath(path))
    for part in [path, *path.parents]:
        require(not part.is_symlink() and not part.is_junction(), 'symlink/junction refused: '+str(part))
    if exists:
        require(path.exists(), 'missing path: '+str(path))
    return path


def contained(root, rel, *, exists=True):
    require(not Path(rel).is_absolute() and '..' not in Path(rel).parts, 'unsafe relative path: '+rel)
    path = checked(root/rel, exists=exists)
    require(path.is_relative_to(root), 'path escapes project')
    return path


def read_json(path):
    return json.loads(path.read_text(encoding='utf-8'))


def write_json(path, data):
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(data, ensure_ascii=False, indent=2)+'\n', encoding='utf-8', newline='\n')


@contextmanager
def writer_lease(project):
    path = contained(project, '.elite/operating-write.lock', exists=False)
    path.parent.mkdir(parents=True, exist_ok=True)
    try:
        fd = os.open(path, os.O_WRONLY | os.O_CREAT | os.O_EXCL, 0o600)
    except FileExistsError as exc:
        raise Blocked('writer already owns .elite/operating-write.lock; do independent read-only work') from exc
    try:
        with os.fdopen(fd, 'w') as stream:
            stream.write(str(os.getpid())+'\n')
        yield
    finally:
        path.unlink()


def block(raw):
    text = raw.decode('utf-8-sig')
    require(text.count(BEGIN)==1 and text.count(END)==1, 'operating entry markers missing or duplicated')
    start, end = text.index(BEGIN), text.index(END)
    require(start < end, 'operating entry markers reversed')
    return text[start:end+len(END)].encode('utf-8')


def verify(project):
    project = checked(project)
    receipt = read_json(contained(project, RECEIPT))
    require(receipt['schema']=='elite-operating-connection/v1', 'unknown receipt schema')
    library = checked(project/receipt['library_relative'])
    require(library!=project and not library.is_relative_to(project) and not project.is_relative_to(library), 'library/project must be separate trees')
    require(set(receipt['authorities'])==set(AUTHORITY_PATHS), 'authority inventory differs')
    for rel, expected in receipt['authorities'].items():
        require(digest(contained(library, rel).read_bytes())==expected, 'authority hash mismatch: '+rel)
    require(receipt['authorities']['INSTALL_AGENT_BRIDGE.ps1']==BASE_SHA, 'unsupported original bridge revision')
    require(set(receipt['entry_blocks'])=={'AGENTS.md','CLAUDE.md'}, 'entry inventory differs')
    for rel, expected in receipt['entry_blocks'].items():
        require(digest(block(contained(project, rel).read_bytes()))==expected, 'entry block mismatch: '+rel)
    required_owned={ENTRY, '.elite/platform-adapter.template.md', *BASE_FILES[2:]}
    require(required_owned.issubset(receipt['owned_files']), 'owned inventory missing required files')
    for rel, expected in receipt['owned_files'].items():
        require(digest(contained(project, rel).read_bytes())==expected, 'installed file hash mismatch: '+rel)
    return receipt, library


def install(project, library, pwsh):
    project, library = checked(project), checked(library)
    require(project.is_dir() and library.is_dir(), 'roots must be directories')
    require(project!=library and not project.is_relative_to(library) and not library.is_relative_to(project), 'library/project must be separate trees')
    require(not any(c in str(project)+str(library) for c in '\n\r`'), 'Markdown control characters in root')
    tool = checked(Path(__file__).parent)
    require(tool.is_relative_to(project), 'materialize this pack inside the consumer project first')
    rel_library = os.path.relpath(library, project).replace('\\','/')
    rel_tool = tool.relative_to(project).as_posix()
    authorities = {rel:digest(contained(library, rel).read_bytes()) for rel in AUTHORITY_PATHS}
    require(authorities['INSTALL_AGENT_BRIDGE.ps1']==BASE_SHA, 'original bridge SHA differs; do not silently upgrade')
    with writer_lease(project):
        if (project/RECEIPT).exists():
            old, bound = verify(project)
            require(bound==library and old['version']==VERSION, 'existing connection targets another revision')
            return {'status':'UNCHANGED', 'version':VERSION}
        paths = BASE_FILES + [ENTRY, '.elite/platform-adapter.template.md', RECEIPT]
        originals = {}
        for rel in paths:
            path = contained(project, rel, exists=False)
            require(not path.exists() or path.is_file(), 'destination must be a regular file: '+rel)
            originals[rel] = path.read_bytes() if path.exists() else None
        for rel in [ENTRY, '.elite/platform-adapter.template.md']:
            require(originals[rel] is None, 'unmanaged file exists: '+rel)
        for rel in ['AGENTS.md','CLAUDE.md']:
            require(not originals[rel] or (BEGIN.encode() not in originals[rel] and END.encode() not in originals[rel]), 'unreceipted operating markers: '+rel)
        try:
            cmd = [pwsh, '-NoProfile', '-File', str(library/'INSTALL_AGENT_BRIDGE.ps1'), '-LibraryRoot', str(library), '-ProjectRoot', str(project), '-Agent', 'Both']
            run = subprocess.run(cmd, capture_output=True, text=True, encoding='utf-8')
            require(run.returncode==0, 'original bridge failed: '+run.stderr[-1600:])
            base_receipt = json.loads(run.stdout)
            require(base_receipt['status']=='ELITE_AGENT_BRIDGE_READY', 'original bridge did not finish')
            supplement = (BEGIN+'\n## Operating protocol entry\n'
                'For every start, resume or handoff, read PROJECT_AGENT_ENTRY.md first. '
                'It verifies the pinned library and routes to '+rel_library+'/'+PROTOCOL+'. '
                'Coordinate the existing project state, events and failure ledger; never inherit library maintenance approvals.\n'+END)
            for rel in ['AGENTS.md','CLAUDE.md']:
                path = project/rel
                raw = path.read_bytes()
                require(BEGIN.encode() not in raw and END.encode() not in raw, 'unexpected marker after base install')
                updated = raw + (b'\n' if raw.endswith(b'\n') else b'\n\n') + supplement.encode('utf-8') + b'\n'
                require(rel!='AGENTS.md' or len(updated)<=32768, 'AGENTS.md entry exceeds 32 KiB; reduce unrelated instructions explicitly')
                path.write_bytes(updated)
            command_tool = (rel_tool+'/connection.py').replace("'", "''")
            entry = f'''# Project operating entry

This file routes instructions; it is not another plan or execution state.
Library: `{rel_library}`. Protocol version: {VERSION}.

1. From THIS project root run `python '{command_tool}' verify --project-root .` (Python 3.14+).
2. If verification fails, record the exact missing/mismatched reference; do not silently rebind it. Continue unrelated work only when its own authorities remain valid.
3. Read [{PROTOCOL}]({rel_library}/{PROTOCOL}) and follow its numbered sequence. Source receipts bind the library files. This receipt is an integrity record, not a security signature or production approval.
4. If PROJECT_EXECUTION_STATE.json and PROJECT_EXECUTION_EVENTS.jsonl exist, validate them with the admitted execution kit, then read context.must_read_refs and next_action. Never read the library's maintenance state as this project's state. If absent, follow the protocol NEW/EXISTING bootstrap and create the existing templates in this project.
5. Follow the project's explicit scope. A connection fixture grants no product implementation, live access or production authorization.

Neutral/platform bootstrap template: [.elite/platform-adapter.template.md](.elite/platform-adapter.template.md).
The installed Codex/Claude entry files are instruction configuration, not proof of either host's agent execution. Grok and other hosts require explicit loading and their own test. All platform execution statuses start NOT_TESTED.
'''
            (project/ENTRY).write_text(entry, encoding='utf-8', newline='\n')
            adapter = (tool/'PLATFORM_ADAPTER_TEMPLATE.md').read_text(encoding='utf-8').replace('{{PROJECT_ENTRY}}', ENTRY)
            (project/'.elite/platform-adapter.template.md').write_text(adapter, encoding='utf-8', newline='\n')
            owned = [ENTRY, '.elite/platform-adapter.template.md', *BASE_FILES[2:]]
            owned += [(tool/name).relative_to(project).as_posix() for name in ['connection.py','INSTALL_PROJECT_OPERATING_BRIDGE.ps1','PLATFORM_ADAPTER_TEMPLATE.md']]
            receipt = {'schema':'elite-operating-connection/v1', 'version':VERSION, 'library_relative':rel_library, 'authorities':authorities,
                'original_bridge':base_receipt, 'entry_blocks':{r:digest(block((project/r).read_bytes())) for r in ['AGENTS.md','CLAUDE.md']},
                'owned_files':{r:digest((project/r).read_bytes()) for r in owned},
                'platform_execution':{p:'NOT_TESTED' for p in ['Codex','Claude','Grok','Other']},
                'scope':'OPERATING_CONNECTION_ONLY', 'product_authorized':False}
            write_json(project/RECEIPT, receipt)
            verify(project)
        except Exception:
            for rel, raw in originals.items():
                path = contained(project, rel, exists=False)
                if raw is None:
                    if path.is_file(): path.unlink()
                else:
                    path.write_bytes(raw)
            raise
    return {'status':'CONNECTION_INSTALLED', 'version':VERSION, 'receipt':RECEIPT, 'product_authorized':False}


def execution_modules(project, kit_relative):
    kit = contained(project, kit_relative)
    require(EXECUTION_TOOL_HASHES, 'execution tool hashes missing')
    for name, expected in EXECUTION_TOOL_HASHES.items():
        require(digest(contained(kit, name).read_bytes())==expected, 'execution tool hash mismatch: '+name)
    sys.path.insert(0, str(kit))
    for name in ['checkpoint_execution_state','validate_execution_state']:
        sys.modules.pop(name, None)
    return importlib.import_module('validate_execution_state'), importlib.import_module('checkpoint_execution_state')


def checkpoint(args):
    project = checked(args.project_root)
    verify(project)
    with writer_lease(project):
        validator, writer = execution_modules(project, args.kit)
        state_path, events_path = project/'PROJECT_EXECUTION_STATE.json', project/'PROJECT_EXECUTION_EVENTS.jsonl'
        candidate = contained(project, args.candidate)
        require(candidate not in [state_path, events_path], 'candidate must be a separate project file')
        before_state = state_path.read_bytes() if state_path.exists() else None
        before_events = events_path.read_bytes() if events_path.exists() else None
        require((before_state is None)==(before_events is None), 'state/events pair is incomplete; recover before advancing')
        revision = read_json(state_path)['revision'] if before_state else 0
        require(revision==args.expected_revision, 'stale expected revision')
        if before_state:
            # The candidate may intentionally update evidence owners. Validate the previous
            # committed state's event binding here; the writer validates all current evidence.
            validator.validate_event_log(events_path, validator.state_sha256(read_json(state_path)))
        next_state = read_json(candidate)
        require(next_state['revision']==revision+1, 'candidate must advance exactly one revision')
        validator.validate_state_data(next_state, project, 'resume')
        try:
            state_path.write_bytes(candidate.read_bytes())
            event = writer.append_checkpoint(state_path, project, events_path, args.event_id, args.event_type, args.actor, args.summary)
        except Exception:
            for path, raw in [(state_path,before_state),(events_path,before_events)]:
                if raw is None:
                    if path.exists(): path.unlink()
                else: path.write_bytes(raw)
            raise
        return {'status':'CHECKPOINT_PASS','revision':next_state['revision'],'event_sha256':event['event_sha256']}


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    commands = parser.add_subparsers(dest='command', required=True)
    install_parser = commands.add_parser('install')
    install_parser.add_argument('--project-root', type=Path, required=True)
    install_parser.add_argument('--library-root', type=Path, required=True)
    install_parser.add_argument('--pwsh', required=True)
    verify_parser = commands.add_parser('verify')
    verify_parser.add_argument('--project-root', type=Path, required=True)
    cp = commands.add_parser('checkpoint')
    cp.add_argument('--project-root', type=Path, required=True)
    cp.add_argument('--kit', default='.elite/execution/engineering_execution_kit')
    cp.add_argument('--candidate', required=True)
    cp.add_argument('--expected-revision', type=int, required=True)
    for name in ['event-id','event-type','actor','summary']: cp.add_argument('--'+name, required=True)
    args = parser.parse_args()
    try:
        if args.command=='install': result=install(args.project_root,args.library_root,args.pwsh)
        elif args.command=='checkpoint': result=checkpoint(args)
        else:
            verify(args.project_root)
            result={'status':'CONNECTION_INTEGRITY_PASS','product_authorized':False}
        print(json.dumps(result))
        return 0
    except (ValueError, KeyError, OSError, UnicodeError, ImportError) as exc:
        print(json.dumps({'status':'CONNECTION_BLOCKED','error':str(exc)}))
        return 2


if __name__=='__main__':
    raise SystemExit(main())
````

### FILE: `operating_connection/fixture_work.py`

```yaml
block_id: "PROJECT-OPERATING-CONNECTION:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local operating bridge, integrity and test orchestration; not upstream company code"
license: "LicenseRef-Workspace-Owner"
sha256: "05dd1c39fa14b52f4733b1d0c8f71fe3cc3a51e2b4ee777736d589eb0f8879e4"
variables: []
secrets_allowed: false
```

````python
#!/usr/bin/env python3
"""AUTHORED connection-test fixture only. Never bootstrap a real product with this."""
from pathlib import Path
import argparse
from datetime import datetime, timezone
import hashlib
import json
import subprocess
import sys

from connection import checked, verify, write_json


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def fixture(project):
    assert (project/'CONNECTION_SMOKE_ONLY').read_text().strip()=='NO_REAL_PRODUCT'
    return json.loads((project/'fixtures/records.json').read_text())


def produce(project):
    rows=fixture(project)
    assert isinstance(rows,list) and all(isinstance(x,str) and x for x in rows)
    result={'count':len(rows),'labels':sorted(rows),'input_sha256':sha(project/'fixtures/records.json')}
    write_json(project/'evidence/result.json',result)
    return result


def record(project, stage):
    verify(project)
    assert (project/'CONNECTION_SMOKE_ONLY').read_text().strip()=='NO_REAL_PRODUCT'
    template=project/'.elite/execution/engineering_execution_kit/execution_state.template.json'
    state_path=project/'PROJECT_EXECUTION_STATE.json'
    previous=json.loads(state_path.read_text()) if state_path.exists() else None
    state=previous or json.loads(template.read_text())
    revision=previous['revision'] if previous else 0
    now=datetime.now(timezone.utc).isoformat().replace('+00:00','Z')
    completed={'baseline':[], 'progress':['T101'], 'failure':['T101'], 'recovery':['T101','T102'], 'handoff':['T101','T102','T103']}[stage]
    blocked=['T102'] if stage=='failure' else []
    pending=[x for x in ['T101','T102','T103'] if x not in completed+blocked]
    state.update({'schema_version':'1.0.0','revision':revision+1,'project':{'id':'operating-connection-fixture','mode':'NEW','rigor':'LIGHT'},'phase':'READINESS','status':'BLOCKED' if blocked else 'ACTIVE'})
    state['baseline']={'kind':'SCAFFOLD','captured_at':now if not previous else previous['baseline']['captured_at'],'inventory_ref':'inventory','delta_scope_ref':None}
    titles={'T101':'Produce sorted synthetic labels and bind the input SHA.','T102':'Observe malformed input, record the actual failure, recover and verify.','T103':'Fresh-session file-only handoff: inspect evidence, independently reproduce result, then record.'}
    (project/'specs/connection-smoke/tasks.md').write_text('# Connection smoke tasks only\n\n'+''.join(f'- [{"x" if t in completed else " "}] {t} {titles[t]}\n' for t in titles),encoding='utf-8',newline='\n')
    ledger=project/'PROJECT_FAILURE_LESSONS.md'
    if stage=='failure':
        with ledger.open('a',encoding='utf-8') as f:f.write('\n## FAIL-20260913-001 — OPEN\nMalformed synthetic records caused a real JSON decoding failure. Evidence: evidence/invalid-input.log. Recovery: restore the declared fixture and rerun only its processor. No live or product system involved.\n')
    if stage=='recovery':
        with ledger.open('a',encoding='utf-8') as f:f.write('\n### FAIL-20260913-001 — CLOSED_FIXTURE\nRestored the declared fixture; rerun exited 0 and result SHA/input/count checked. Evidence: evidence/recovery.log and evidence/result.json. Reopen if the same fixture produces a mismatch.\n')
    paths={'failure':'PROJECT_FAILURE_LESSONS.md','tasks':'specs/connection-smoke/tasks.md','blueprint':'PROJECT_BLUEPRINT.md','authority':'PROJECT_AUTHORITY_MAP.md','packs':'PROJECT_PACK_PLAN.md','inventory':'evidence/initial-inventory.json'}
    if (project/'evidence/result.json').exists():paths['result']='evidence/result.json'
    for name in ['invalid-input','recovery','agent-handoff']:
        path='evidence/'+name+('.json' if name=='agent-handoff' else '.log')
        if (project/path).exists():paths[name]=path
    state['evidence']=[{'id':key,'kind':'failure_ledger' if key=='failure' else 'artifact','path':rel,'sha256':sha(project/rel),'produced_at':now,'producer':'connection-smoke-fixture'} for key,rel in paths.items()]
    state['links']['failure_ledger']='failure';state['links']['tasks']='tasks'
    state['task_progress']={'pending':pending,'in_progress':[],'blocked':blocked,'completed':completed}
    state['context']={'summary':'Connection smoke only; phase READINESS means no product implementation or production authorization. Stage: '+stage,'must_read_refs':[x for x in ['tasks','failure','blueprint','result'] if x in paths],'reuse_without_reload_refs':['authority','packs','inventory']}
    state['open']={'blockers':['Malformed fixture must be recovered.'] if blocked else [],'failure_ids':['FAIL-20260913-001'] if blocked else []}
    state['next_action']=('T103: read PROJECT_AGENT_ENTRY.md and the linked protocol, validate connection and state/events; inspect specs/connection-smoke/tasks.md and evidence/result.json. Independently recompute sorted labels/count/input SHA from fixtures/records.json; write evidence/agent-handoff.json with verified=true, input_sha256, count, labels, authority_path, observed_revision, files_read and limitations. Then run python -B .elite/connection-tools/operating_connection/fixture_work.py record --project-root . --stage handoff. This is a fixture; do not build a product.' if stage=='recovery' else
        'Connection fixture completed. Stop; wait for explicit user authorization before any real project.' if stage=='handoff' else
        'T102: restore fixtures/records.json from the declared synthetic inputs, rerun fixture_work.py run, retain real outputs and record recovery.' if stage=='failure' else
        'Follow the next unchecked task in specs/connection-smoke/tasks.md; this fixture is not a real product.')
    candidate=project/'.elite/state-candidate.json';write_json(candidate,state)
    cmd=[sys.executable,'-B',str(Path(__file__).with_name('connection.py')),'checkpoint','--project-root',str(project),'--candidate','.elite/state-candidate.json','--expected-revision',str(revision),'--event-id','CONNECTION-'+str(revision+1).zfill(3),'--event-type',stage.upper(),'--actor','fresh-session-agent' if stage=='handoff' else 'fixture-harness','--summary','Observed connection fixture '+stage]
    p=subprocess.run(cmd,capture_output=True,text=True,encoding='utf-8')
    (project/'evidence'/('checkpoint-'+str(revision+1)+'.log')).write_text(p.stdout+p.stderr,encoding='utf-8')
    assert p.returncode==0,p.stdout+p.stderr
    return json.loads(p.stdout)


def main():
    parser=argparse.ArgumentParser(description=__doc__)
    parser.add_argument('command',choices=['run','record'])
    parser.add_argument('--project-root',type=Path,required=True)
    parser.add_argument('--stage',choices=['baseline','progress','failure','recovery','handoff'])
    args=parser.parse_args();project=checked(args.project_root)
    print(json.dumps(produce(project) if args.command=='run' else record(project,args.stage)))


if __name__=='__main__':main()
````

### FILE: `operating_connection/INSTALL_PROJECT_OPERATING_BRIDGE.ps1`

```yaml
block_id: "PROJECT-OPERATING-CONNECTION:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local operating bridge, integrity and test orchestration; not upstream company code"
license: "LicenseRef-Workspace-Owner"
sha256: "bcc47601a19c6726d08845347da4e489b3e315cb3dcc572e66a76b39e67f778d"
variables: []
secrets_allowed: false
```

````powershell
#requires -Version 7.0
<# AUTHORED orchestration; invokes the original immutable bridge. #>
[CmdletBinding()]
param(
  [Parameter(Mandatory=$true)][string] $ProjectRoot,
  [Parameter(Mandatory=$true)][string] $LibraryRoot,
  [string] $Python = 'python'
)
$ErrorActionPreference = 'Stop'
$pwsh = Join-Path $PSHOME $(if ($IsWindows) { 'pwsh.exe' } else { 'pwsh' })
& $Python -X utf8 -B (Join-Path $PSScriptRoot 'connection.py') install --project-root $ProjectRoot --library-root $LibraryRoot --pwsh $pwsh
if ($LASTEXITCODE -ne 0) { throw 'Operating bridge installation blocked; inspect the JSON error above.' }
````

### FILE: `operating_connection/PLATFORM_ADAPTER_TEMPLATE.md`

```yaml
block_id: "PROJECT-OPERATING-CONNECTION:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local operating bridge, integrity and test orchestration; not upstream company code"
license: "LicenseRef-Workspace-Owner"
sha256: "038eb2fff61e452862becc74d30512e475559da09f0639e71112dd87a7b2eac7"
variables: []
secrets_allowed: false
```

````markdown
# Platform entry template — NOT_TESTED

AUTHORED portable instructions. Installing this template does not demonstrate any vendor's instruction discovery, filesystem access or agent execution.

- Platform and version: {{PLATFORM_VERSION}}
- Project root accessible to this agent: {{PROJECT_ROOT}}
- Library location in that environment: {{LIBRARY_ROOT}}
- Role and allowed actions: {{ROLE_AND_AUTHORIZED_SCOPE}}
- Tools and identities actually available: {{TOOLS_AND_IDENTITY_REFERENCES}}
- Authorized external actions and approval boundaries: {{EXTERNAL_ACTION_POLICY}}
- Execution evidence: NOT_TESTED until a real fresh session supplies its transcript and artifacts.

## Initial message / persistent role instruction

Work in {{PROJECT_ROOT}}. Before acting, read {{PROJECT_ENTRY}}. Verify its connection receipt and follow the linked operating protocol. Resume from PROJECT_EXECUTION_STATE.json and PROJECT_EXECUTION_EVENTS.jsonl when present; use context.must_read_refs and next_action. Follow the project's authorization and existing gates. Write progress and failures in the project owners; preserve the library. Record missing access and continue independently authorized work. Do not infer approvals from library readiness, chat memory or this template. Never substitute a different library revision silently.

## Platform-specific loading to verify

| Platform | Configuration starting point | Actual agent execution |
|---|---|---|
| Codex | Original bridge AGENTS.md and project Skill, plus operating entry block | NOT_TESTED |
| Claude | Original bridge CLAUDE.md and project Skill, plus operating entry block | NOT_TESTED |
| Grok | Explicit role description and saved Skill linking this entry on its accessible computer | NOT_TESTED |
| Other | Explicit initial/persistent instruction linking this entry | NOT_TESTED |

Grok's cloud filesystem is separate from the user's desktop. Establish real repository/file access and the intended revision first; shared bot roles do not imply security isolation. No connector, login, Skill registration or routine is installed by this template. Official mechanism: https://docs.x.ai/grok-bot/skills-routines-and-automations and https://docs.x.ai/grok-bot/computer-and-apps .

To qualify a platform: use a fresh session without the prior conversation; demonstrate authority selection, bounded task, actual failure and recovery, evidence and a subsequent file-only resume. Preserve the observed platform/version and limitations in the test receipt. A local harness PASS cannot fill these cells automatically.
````

### FILE: `operating_connection/smoke_connection.py`

```yaml
block_id: "PROJECT-OPERATING-CONNECTION:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local operating bridge, integrity and test orchestration; not upstream company code"
license: "LicenseRef-Workspace-Owner"
sha256: "341facaadfbcd6c4c3c37b451c8916e98ffe8b9ab11aa3560d348a9f08840bef"
variables: []
secrets_allowed: false
```

````python
#!/usr/bin/env python3
"""AUTHORED bounded smoke: instruction connection, fixtures, failures and resume."""
from pathlib import Path
import argparse
import hashlib
import json
import shutil
import subprocess
import sys
import time

from connection import verify, write_json, checked


def sha(path):return hashlib.sha256(path.read_bytes()).hexdigest()


def run(command, cwd, expected=0):
    start=time.monotonic();p=subprocess.run(command,cwd=cwd,capture_output=True,text=True,encoding='utf-8')
    result={'command':command,'exit_code':p.returncode,'seconds':round(time.monotonic()-start,3),'stdout':p.stdout,'stderr':p.stderr}
    if expected is not None:assert p.returncode==expected,result
    return result


def prepare(args):
    destination=args.destination.absolute();assert not destination.exists(),destination
    destination.mkdir(parents=True);project=destination/'Fictional Project';project.mkdir()
    source=Path(__file__).parent
    tool=project/'.elite/connection-tools/operating_connection';shutil.copytree(source,tool,ignore=shutil.ignore_patterns('__pycache__'))
    foreign=b'# Existing consumer instructions\r\n\r\nPreserve this exact foreign text.\r\n'
    (project/'AGENTS.md').write_bytes(foreign)
    (project/'CLAUDE.md').write_bytes(b'# Existing Claude notes\n')
    records=[];checks={}
    def call(command,expected=0):
        row=run(command,project,expected);records.append(row);return row
    base=[args.pwsh,'-NoProfile','-File',str(tool/'INSTALL_PROJECT_OPERATING_BRIDGE.ps1'),'-ProjectRoot',str(project),'-LibraryRoot',str(args.library_root),'-Python',sys.executable]
    call(base);receipt,library=verify(project)
    checks['entry_routes_to_protocol']=(str(library/'markdown_system/FRANCHISE_PROJECT_OPERATING_PROTOCOL.md').replace('\\','/') is not None and 'FRANCHISE_PROJECT_OPERATING_PROTOCOL.md' in (project/'PROJECT_AGENT_ENTRY.md').read_text())
    checks['foreign_instructions_preserved']=(project/'AGENTS.md').read_bytes().startswith(foreign)
    before={p.relative_to(project).as_posix():sha(p) for p in project.rglob('*') if p.is_file()}
    call(base);after={p.relative_to(project).as_posix():sha(p) for p in project.rglob('*') if p.is_file()}
    checks['repeat_install_byte_identical']=before==after
    cli=[sys.executable,'-B',str(tool/'connection.py')]
    entry=project/'PROJECT_AGENT_ENTRY.md';raw=entry.read_bytes();entry.unlink()
    call(cli+['verify','--project-root',str(project)],2);entry.write_bytes(raw);call(cli+['verify','--project-root',str(project)])
    checks['missing_entry_detected']=True
    rp=project/'.elite/operating-connection.json';raw=rp.read_bytes();bad=json.loads(raw);bad['authorities']['markdown_system/FRANCHISE_PROJECT_OPERATING_PROTOCOL.md']='0'*64;write_json(rp,bad)
    call(cli+['verify','--project-root',str(project)],2);rp.write_bytes(raw);call(cli+['verify','--project-root',str(project)])
    checks['altered_authority_reference_detected']=True
    # Unmanaged collisions are rejected before the original bridge can write.
    collision=destination/'Existing Collision';collision.mkdir();ct=collision/'.elite/connection-tools/operating_connection';shutil.copytree(source,ct,ignore=shutil.ignore_patterns('__pycache__'))
    (collision/'PROJECT_AGENT_ENTRY.md').write_text('Unmanaged owner\n');(collision/'AGENTS.md').write_bytes(foreign)
    saved={p.relative_to(collision).as_posix():sha(p) for p in collision.rglob('*') if p.is_file()}
    row=run([sys.executable,'-B',str(ct/'connection.py'),'install','--project-root',str(collision),'--library-root',str(library),'--pwsh',args.pwsh],collision,2);records.append(row)
    checks['unmanaged_collision_preserved']=saved=={p.relative_to(collision).as_posix():sha(p) for p in collision.rglob('*') if p.is_file()}
    call([args.pwsh,'-NoProfile','-File',str(library/'materialize_markdown_pack.ps1'),'-PackFile',str(library/'implementation_packs/ENGINEERING_EXECUTION_VALIDATOR.md'),'-Destination',str(project/'.elite/execution')])
    (project/'evidence').mkdir();(project/'fixtures').mkdir();(project/'specs/connection-smoke').mkdir(parents=True)
    (project/'CONNECTION_SMOKE_ONLY').write_text('NO_REAL_PRODUCT\n')
    (project/'PROJECT_BLUEPRINT.md').write_text('# Fictional connection fixture\nOnly test file-based agent continuity with three synthetic labels. No business, accounts, franchise, deployment or product implementation is authorized. The project stays in READINESS; product admission is not assessed.\n')
    (project/'PROJECT_AUTHORITY_MAP.md').write_text('# Selected authorities\nRead PROJECT_AGENT_ENTRY.md and the exact operating protocol it links. Authority digests live in the connection receipt. Execution methods come from the existing ENGINEERING_EXECUTION_VALIDATOR pack. The fixture tasks govern this test only.\n')
    (project/'PROJECT_PACK_PLAN.md').write_text('# Test-only selection\nPROJECT_OPERATING_CONNECTION 0.1.0 and existing ENGINEERING_EXECUTION_VALIDATOR. No franchise profile or product runtime selected.\n')
    (project/'PROJECT_FAILURE_LESSONS.md').write_text('# Connection fixture failure lessons\nNo failure observed at baseline.\n')
    write_json(project/'fixtures/records.json',['gamma','alpha','beta'])
    write_json(project/'evidence/initial-inventory.json',{'baseline':'SCAFFOLD','files':sorted(str(p.relative_to(project)).replace('\\','/') for p in project.rglob('*') if p.is_file())})
    fixture=[sys.executable,'-B',str(tool/'fixture_work.py')]
    call(fixture+['record','--project-root',str(project),'--stage','baseline'])
    call(fixture+['run','--project-root',str(project)])
    call(fixture+['record','--project-root',str(project),'--stage','progress'])
    expected_input=(project/'fixtures/records.json').read_bytes();(project/'fixtures/records.json').write_text('{broken')
    fail=call(fixture+['run','--project-root',str(project)],expected=None);assert fail['exit_code']!=0
    (project/'evidence/invalid-input.log').write_text(json.dumps(fail,indent=2),encoding='utf-8')
    call(fixture+['record','--project-root',str(project),'--stage','failure'])
    (project/'fixtures/records.json').write_bytes(expected_input)
    recovery=call(fixture+['run','--project-root',str(project)])
    (project/'evidence/recovery.log').write_text(json.dumps(recovery,indent=2),encoding='utf-8')
    call(fixture+['record','--project-root',str(project),'--stage','recovery'])
    checks['material_progress_failure_recovery']=True
    state=project/'PROJECT_EXECUTION_STATE.json';events=project/'PROJECT_EXECUTION_EVENTS.jsonl';before_pair=(sha(state),sha(events))
    checkpoint=cli+['checkpoint','--project-root',str(project),'--candidate','.elite/state-candidate.json','--expected-revision','3','--event-id','STALE-001','--event-type','STALE','--actor','smoke','--summary','Stale write must fail']
    call(checkpoint,2);checks['stale_revision_rejected']=(sha(state),sha(events))==before_pair
    lock=project/'.elite/operating-write.lock';lock.write_text('owned-by-test\n')
    call(checkpoint,2);assert lock.read_text()=='owned-by-test\n';lock.unlink()
    checks['concurrent_writer_rejected']=(sha(state),sha(events))==before_pair
    validator=[sys.executable,'-B',str(project/'.elite/execution/engineering_execution_kit/validate_execution_state.py'),str(state),'--project-root',str(project),'--events',str(events),'--level','resume']
    call(validator);checks['real_execution_validator_pass']=True
    assert all(checks.values()),checks
    write_json(destination/'harness.json',{'status':'AWAITING_FRESH_AGENT','project':str(project),'checks':checks,'commands':records,'platform_execution':receipt['platform_execution']})
    print(json.dumps({'status':'AWAITING_FRESH_AGENT','project':str(project),'harness_checks':len(checks),'next':'Fresh session must use only project files to finish T103.'}))


def finalize(args):
    root=checked(args.destination);harness=json.loads((root/'harness.json').read_text());project=Path(harness['project'])
    verify(project)
    agent=json.loads((project/'evidence/agent-handoff.json').read_text())
    expected=json.loads((project/'evidence/result.json').read_text())
    assert agent['verified'] is True and agent['observed_revision']==4
    for key in ['count','labels','input_sha256']:assert agent[key]==expected[key],key
    assert 'FRANCHISE_PROJECT_OPERATING_PROTOCOL.md' in agent['authority_path']
    required={'PROJECT_AGENT_ENTRY.md','PROJECT_EXECUTION_STATE.json','PROJECT_EXECUTION_EVENTS.jsonl','PROJECT_FAILURE_LESSONS.md','fixtures/records.json'}
    assert required.issubset(set(agent['files_read']))
    state=json.loads((project/'PROJECT_EXECUTION_STATE.json').read_text());assert state['revision']==5 and state['task_progress']['completed']==['T101','T102','T103'] and state['release']['target']=='NONE'
    v=project/'.elite/execution/engineering_execution_kit/validate_execution_state.py'
    command=run([sys.executable,'-B',str(v),str(project/'PROJECT_EXECUTION_STATE.json'),'--project-root',str(project),'--events',str(project/'PROJECT_EXECUTION_EVENTS.jsonl'),'--level','resume'],project)
    harness['commands'].append(command);harness['checks']['fresh_agent_file_only_resume']=True;harness['checks']['no_product_admission']=True
    harness['status']='PASS_CONNECTION_SMOKE';harness['artifacts_sha256']={p.relative_to(root).as_posix():sha(p) for p in project.rglob('*') if p.is_file() and '__pycache__' not in p.parts}
    write_json(root/'smoke-receipt.json',harness)
    print(json.dumps({'status':harness['status'],'checks':len(harness['checks']),'project_revision':state['revision'],'receipt':str(root/'smoke-receipt.json')}))


if __name__=='__main__':
    p=argparse.ArgumentParser(description=__doc__);p.add_argument('mode',choices=['prepare','finalize']);p.add_argument('--destination',type=Path,required=True);p.add_argument('--library-root',type=Path);p.add_argument('--pwsh');a=p.parse_args()
    prepare(a) if a.mode=='prepare' else finalize(a)
````

### FILE: `operating_connection/verify_smoke_receipt.py`

```yaml
block_id: "PROJECT-OPERATING-CONNECTION:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local operating bridge, integrity and test orchestration; not upstream company code"
license: "LicenseRef-Workspace-Owner"
sha256: "67569984500ad8fdd519b1348324aa8cced1e148cb51f0323e33fd5ea2aa320a"
variables: []
secrets_allowed: false
```

````python
#!/usr/bin/env python3
"""Verify retained smoke artifact hashes; this is not a signature verifier."""
from pathlib import Path
import argparse
import hashlib
import json

from connection import checked, contained


if __name__=='__main__':
    p=argparse.ArgumentParser(description=__doc__);p.add_argument('--root',type=Path,required=True);a=p.parse_args();root=checked(a.root)
    receipt=json.loads((root/'smoke-receipt.json').read_text())
    assert receipt['status']=='PASS_CONNECTION_SMOKE' and all(receipt['checks'].values())
    for rel,expected in receipt['artifacts_sha256'].items():
        assert hashlib.sha256(contained(root,rel).read_bytes()).hexdigest()==expected,rel
    print(json.dumps({'status':'SMOKE_RECEIPT_HASHES_PASS','artifacts':len(receipt['artifacts_sha256'])}))
````

## 6. Configuration surface

Explicit ProjectRoot and LibraryRoot, local Python and PowerShell executables. No network, keys, accounts, credentials or product settings. Templates use placeholders and NOT_TESTED states. Preserve materialization and installation receipts in the project.

## 7. Dependency bill

Python standard library and PowerShell process launch only. Pin and verify the existing library's bridge and execution scripts. Method references: DORA/Google Cloud, Microsoft Well-Architected, official Grok skills documentation in the operating protocol. This pack contains no copied vendor runtime.

## 8. Apply order

Follow START_FRANCHISE: materialize this pack into an absent consumer tooling directory, run operating_connection/INSTALL_PROJECT_OPERATING_BRIDGE.ps1, then open PROJECT_AGENT_ENTRY.md in the consumer. A repeated installer invocation verifies and returns UNCHANGED; it does not silently upgrade a modified installation. Preserve foreign instructions. An explicit upgrade requires a separately reviewed delta.

## 9. Verification

Run connection.py verify against the consumer; run the admitted execution validator when state exists. The bounded harness is smoke_connection.py prepare --destination <absent disposable root> --library-root <library> --pwsh <PowerShell>. It installs the bridge, checks drift/conflicts/repeatability, materializes only the execution tooling and records a synthetic progress/failure/recovery. Then a fresh agent, with no prior conversation, follows project files to finish T103. Run smoke_connection.py finalize --destination <same root> and verify_smoke_receipt.py --root <same root>. Do not replace fresh-agent evidence with a scripted assertion. All product and platform permissions remain separate.

## 10. Reconstruction evidence

See reconstruction_evidence/PROJECT_OPERATING_CONNECTION_REVIEW.md for observed results, method source hashes and limits. This additive connection does not alter V402 source330, metadata331, closure337 or either immutable ZIP. The source remains reconstructible; its operational claim is only the tested connection scope.
