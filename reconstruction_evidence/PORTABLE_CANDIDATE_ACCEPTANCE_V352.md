# V352 — aceptación del candidato portable antes de publicar

Mantenimiento de biblioteca, 2026-09-09T15:06:45.442685Z. Entrada147 validada; stage <LOCALAPPDATA>/Temp/elite-v352-32e0c00909af4e31a62ae41f95568773.
Usuario solicita destrabar el progreso y mantener controles correctos.

## Corrección del orden, sin reducción de criterios

FAIL613/LIB2329: TEST06 ya pide un candidato; TEST08 pide artefacto exacto.
La guía0.1.1 y anteriores informes aplazaron su aceptación hasta el release final.
Guía0.1.2 permite ensayar una copia portable completa, congelada y NO PUBLICADA.
Los estados sólo cambiarán tras los ensayos; TEST07 conserva firma, SCA/notices,
reconstrucción independiente y gates materiales. Una publicación debe usar el
payload aceptado o revalidar sus cambios. Original147 conservado sin sobrescribir.

No prueba de usuario humano independiente, producción, pérdida de energía,
proveedores, corpus ni cierre de readiness se infiere de este ensayo local.
El alcance cubre al agente consumidor en NEW/EXISTING, comandos de la guía,
fallo parcial de instalación, restauración exacta y cadena de reanudación.

## Estado

Preparación148; ensayos de candidato completo pendientes. Contrato41/48passed,
5blocked/2planned. La corrección del procedimiento por sí sola no aumenta el avance.

## Qualification source: accept_candidate352.py (AUTHORED)

````python
from pathlib import Path,PurePosixPath
import hashlib,json,os,re,shutil,subprocess,sys,time,zipfile
S=Path(__file__).resolve().parent;py=Path(sys.executable);pwsh=Path('<USERPROFILE>/.cache/codex-runtimes/codex-primary-runtime/dependencies/native/powershell/pwsh.exe');sha=lambda p:hashlib.sha256(p.read_bytes()).hexdigest()
A=Path(os.environ.get('ELITE_ACCEPTANCE_ARTIFACT_ROOT',str(S)));archive=A/'library-candidate148-not-for-release.zip';expected=archive.with_suffix('.zip.sha256').read_text().split()[0];assert sha(archive)==expected
extract=S/'accepted-extraction';extract.mkdir();prefix='Elite Engineering Library/'
with zipfile.ZipFile(archive) as z:
 names=z.namelist();assert len(names)==len(set(n.casefold() for n in names))
 assert all(n.startswith(prefix) and '\\' not in n and not any(x in ('..','.') for x in PurePosixPath(n).parts) for n in names)
 manifest_name=prefix+'DISTRIBUTION_SHA256SUMS.txt';manifest={}
 for line in z.read(manifest_name).decode('utf-8').splitlines():
  m=re.fullmatch(r'([a-f0-9]{64})  (.+)',line);assert m
  assert m[2] not in manifest;manifest[m[2]]=m[1]
 assert set(names)=={prefix+n for n in manifest}|{manifest_name}
 for n,h in manifest.items():assert hashlib.sha256(z.read(prefix+n)).hexdigest()==h;assert sha(A/'source148'/n)==h
 forbidden=['PROJECT_EXECUTION_STATE.json','PROJECT_EXECUTION_EVENTS.jsonl','PROJECT_READINESS_RECORD.md','PROJECT_READINESS_GATE.json','PROJECT_READINESS_REPORT.json','PROJECT_ENGINEERING_CONTRACT.json','PROJECT_BLUEPRINT.md','PROJECT_AUTHORITY_MAP.md','PROJECT_PACK_PLAN.md']
 assert not any(n in manifest for n in forbidden)
 assert not any(n.startswith(('specs/','.specify/','PROJECT_ADVISORY_')) for n in manifest)
 for n in names:
  dest=extract.joinpath(*PurePosixPath(n).parts);assert dest.is_relative_to(extract)
  dest.parent.mkdir(parents=True,exist_ok=True)
  with dest.open('xb') as f:f.write(z.read(n))
root=extract/'Elite Engineering Library'
for n,h in manifest.items():assert sha(root/n)==h
receipt={'archive_path':str(archive),'archive_sha256':expected,'manifest_sha256':sha(root/'DISTRIBUTION_SHA256SUMS.txt'),'payload_files':len(manifest),'entries':len(names),'source_revision':148,'all_entries_match_source':True,'local_approvals_excluded':True,'sources':manifest,'published':False,'signed':False}
(S/'archive-verification.json').write_text(json.dumps(receipt,indent=2)+'\n');print('ARCHIVE_INDEPENDENT_PASS',len(manifest),'payload files',flush=True)
env=dict(os.environ);env['PATH']=str(py.parent)+os.pathsep+str(pwsh.parent)+os.pathsep+env['PATH'];env['PYTHONDONTWRITEBYTECODE']='1'
calls=[]
def run(name,args,code=0,marker=None,timeout=120,custom_env=None):
 start=time.perf_counter();p=subprocess.run([str(x) for x in args],env=custom_env or env,stdout=subprocess.PIPE,stderr=subprocess.PIPE,timeout=timeout);(S/(name+'.stdout')).write_bytes(p.stdout);(S/(name+'.stderr')).write_bytes(p.stderr)
 text=(p.stdout+p.stderr).decode('utf-8','replace');calls.append({'name':name,'code':p.returncode,'expected':code,'elapsed_ms':round((time.perf_counter()-start)*1000,3),'stdout_sha256':sha(S/(name+'.stdout')),'stderr_sha256':sha(S/(name+'.stderr'))});assert p.returncode==code,(name,text[-1800:]);assert marker is None or marker in text,(name,text[-1800:]);print(name,'PASS',flush=True);return p
run('extracted-structural',[pwsh,'-NoProfile','-File',root/'VERIFY_LIBRARY.ps1'],marker='VERIFY_LIBRARY_PASS',timeout=600)
run('extracted-execution-kit',[pwsh,'-NoProfile','-File',root/'materialize_markdown_pack.ps1','-PackFile',root/'implementation_packs/ENGINEERING_EXECUTION_VALIDATOR.md','-Destination',S/'archive-execution'])
kit=S/'archive-execution/engineering_execution_kit';sys.path.insert(0,str(kit));from test_execution_state import ExecutionStateTests
from checkpoint_execution_state import append_checkpoint

def dump(p,v):p.write_text(json.dumps(v,ensure_ascii=False,indent=2)+'\n',encoding='utf-8')
def inventory(p):return {x.relative_to(p).as_posix():sha(x) for x in p.rglob('*') if x.is_file()}
def fixture(dest,mode):
 dest.mkdir(parents=True);case=ExecutionStateTests();case.root=dest;case.events=dest/'PROJECT_EXECUTION_EVENTS.jsonl';case.files={};case._write('failure','PROJECT_FAILURE_LESSONS.md','# Synthetic acceptance fixture; no production approval\n')
 if mode=='EXISTING':
  (dest/'user-original.txt').write_bytes(b'\xef\xbb\xbfPreserve original\r\n');(dest/'AGENTS.md').write_text('# Existing user instructions\nPreserve my code.\n');initial=inventory(dest)
 state=case._base(mode);state['phase']='DISCOVERY';state['status']='BLOCKED';state['open']['blockers']=['Synthetic fixture intake intentionally incomplete; no product authorization'];state['context']['summary']='Synthetic artifact-bound acceptance fixture; no approval';state['next_action']='Complete real readiness; this fixture proves only tooling acceptance.'
 if mode=='EXISTING':
  case.files['inventory'].write_text(json.dumps(initial,sort_keys=True)+'\n');case.files['delta'].write_text('Authorized synthetic bridge installation only; preserve original user bytes/events.\n');state['evidence']=[case._evidence(n,'baseline_inventory' if n=='inventory' else 'delta_scope' if n=='delta' else 'failure_ledger') for n in case.files]
 dump(dest/'PROJECT_EXECUTION_STATE.json',state);append_checkpoint(dest/'PROJECT_EXECUTION_STATE.json',dest,case.events,'EVT-0001','BASELINE_CAPTURED','acceptance-fixture','Synthetic baseline; no release approval');return state
base= S/'guide-resume-source';fixture(base,'EXISTING')
guide_dir=S/'guide-run';guide_dir.mkdir();source=(A/'guide_journey_v327.py').read_text();before=source
source=source.replace("'no release archive'","'exact V352 candidate archive independently verified by outer harness; no final release approval'")
(guide_dir/'guide_journey.py').write_text(source,encoding='utf-8',newline='\n');(S/'guide-harness-delta.json').write_text(json.dumps({'original_sha256':sha(A/'guide_journey_v327.py'),'adapted_sha256':sha(guide_dir/'guide_journey.py'),'delta':'Only receipt non-claim clarifies outer exact-archive verification; assertions and executed guide blocks unchanged'},indent=2)+'\n')
genv=dict(env,ELITE_LIBRARY_ROOT=str(root),ELITE_PWSH_EXECUTABLE=str(pwsh),ELITE_RESUME_FIXTURE=str(base))
run('actual-guide',[py,'-X','utf8','-B',guide_dir/'guide_journey.py'],marker='GUIDE_CLI_PASS',timeout=240,custom_env=genv)
# Real installer from the extracted archive; only an OS sharing lock injects failure.
lockscript=S/'locked-install.ps1'
lockscript.write_text('''param([string]$Library,[string]$Project,[string]$Skill)
$ErrorActionPreference='Stop'
$handle=[IO.File]::Open($Skill,[IO.FileMode]::Open,[IO.FileAccess]::Read,[IO.FileShare]::Read)
$observed=$null
try {
 try { & (Join-Path $Library 'INSTALL_AGENT_BRIDGE.ps1') -LibraryRoot $Library -ProjectRoot $Project -Agent Both | Out-Null }
 catch {
  $codes=@();$cause=$_.Exception
  while($null -ne $cause){$codes+=([long]$cause.HResult -band 4294967295L);$cause=$cause.InnerException}
  $observed=@{blocked=$true;hresults=$codes;message=$_.Exception.Message;stack=$_.ScriptStackTrace}
 }
} finally {$handle.Dispose()}
if($null -eq $observed){throw 'Locked fixture unexpectedly installed'}
$observed | ConvertTo-Json -Depth 6
''',encoding='utf-8',newline='\n')
faults=[]
for mode in ['NEW','EXISTING']:
 dest=S/('fault-'+mode);state=fixture(dest,mode);skill=dest/'.agents/skills/elite-engineering-library/SKILL.md';skill.parent.mkdir(parents=True);skill.write_text('<!-- ELITE-ENGINEERING-LIBRARY:MANAGED-SKILL -->')
 before=inventory(dest);events=(dest/'PROJECT_EXECUTION_EVENTS.jsonl').read_bytes();dump(S/(mode+'-before.json'),before)
 p=run(mode+'-partial-install',[pwsh,'-NoProfile','-File',lockscript,'-Library',root,'-Project',dest,'-Skill',skill]);locked=json.loads(p.stdout);assert any(h in (0x80070020,0x80070005) for h in locked['hresults']) and 'Invoke-PlannedWrites' in locked['stack'],locked
 assert inventory(dest)==before,'rollback did not restore exact files';assert (dest/'PROJECT_EXECUTION_EVENTS.jsonl').read_bytes()==events
 resume=[py,'-X','utf8','-B',kit/'validate_execution_state.py',dest/'PROJECT_EXECUTION_STATE.json','--project-root',dest,'--events',dest/'PROJECT_EXECUTION_EVENTS.jsonl','--level','resume']
 run(mode+'-resume-after-failure',resume,marker='EXECUTION_STATE_PASS')
 run(mode+'-retry-install',[pwsh,'-NoProfile','-File',root/'INSTALL_AGENT_BRIDGE.ps1','-LibraryRoot',root,'-ProjectRoot',dest,'-Agent','Both'],marker='ELITE_AGENT_BRIDGE_READY')
 for n,h in before.items():
  if n not in ['AGENTS.md','.agents/skills/elite-engineering-library/SKILL.md']:assert sha(dest/n)==h,n
 assert (dest/'PROJECT_EXECUTION_EVENTS.jsonl').read_bytes()==events
 if mode=='EXISTING':assert (dest/'AGENTS.md').read_text().startswith('# Existing user instructions\nPreserve my code.\n')
 run(mode+'-resume-after-recovery',resume,marker='EXECUTION_STATE_PASS')
 state['revision']=2;state['next_action']='Synthetic partial bridge installation recovered; actual readiness remains BLOCKED.';dump(dest/'PROJECT_EXECUTION_STATE.json',state)
 checkpoint=[py,'-X','utf8','-B',kit/'checkpoint_execution_state.py',dest/'PROJECT_EXECUTION_STATE.json','--project-root',dest,'--events',dest/'PROJECT_EXECUTION_EVENTS.jsonl','--event-id','EVT-0002','--event-type','STATE_ADVANCED','--actor','acceptance-fixture','--summary','Exact candidate bridge recovered; synthetic fixture remains BLOCKED']
 run(mode+'-checkpoint-recovery',checkpoint,marker='EXECUTION_CHECKPOINT_PASS');advanced=(dest/'PROJECT_EXECUTION_EVENTS.jsonl').read_bytes();assert advanced.startswith(events) and len(advanced.splitlines())==2
 run(mode+'-duplicate-checkpoint',checkpoint,2,'already');assert (dest/'PROJECT_EXECUTION_EVENTS.jsonl').read_bytes()==advanced
 run(mode+'-final-resume',resume,marker='EXECUTION_STATE_PASS');faults.append({'mode':mode,'rollback_exact':True,'old_events_preserved':True,'recovery_events':2,'duplicate_append_rejected':True,'sharing_violation':locked,'fixture':str(dest)})
run('candidate-entry-lifecycle',[pwsh,'-NoProfile','-File',root/'markdown_system/test_agent_entry_lifecycle.ps1','-LibraryRoot',root,'-KeepEvidence'],marker='55 checks',timeout=240)
for n,h in manifest.items():assert sha(root/n)==h,n
assert sha(archive)==expected
(S/'acceptance.json').write_text(json.dumps({'status':'PASS','archive_sha256':expected,'manifest_sha256':receipt['manifest_sha256'],'payload_files':len(manifest),'guide_commands':13,'help_checks':8,'lifecycle_checks':55,'partial_install_modes':faults,'calls':calls,'sources_unchanged':True,'test06':'PASS_ON_EXACT_CANDIDATE','test08':'PASS_ON_EXACT_CANDIDATE','test07':'BLOCKED','limits':['synthetic projects and actual CLI agent execution, not independent human UX study','normal sharing-failure rollback, not power-loss or concurrent installation','no final release/signature/SCA/target approval'],'source_projection':{n:manifest[n] for n in ['INSTALL_AGENT_BRIDGE.ps1','materialize_markdown_pack.ps1','VERIFY_LIBRARY.ps1','markdown_system/PROJECT_ENTRY_CLI_GUIDE.md','markdown_system/test_agent_entry_lifecycle.ps1','implementation_packs/MARKDOWN_COMPOSITOR_CORE.md','implementation_packs/PROJECT_START_READINESS_VALIDATOR.md','implementation_packs/ENGINEERING_EXECUTION_VALIDATOR.md','markdown_system/PROJECT_READINESS_GATE_PACK_PLAN.md']}},indent=2)+'\n')
print('CANDIDATE_ACCEPTANCE_PASS TEST06 TEST08 exact_archive',expected,flush=True)
````

## Candidate acceptance149 — observed 2026-09-09T15:19:50.713560Z

TEST06 and TEST08 pass on complete internal candidate148:
`0590588c4ab139e86e581b13c1db561f9342462646105c109f190aeefc755085` (library-candidate148-not-for-release.zip).
804source files plus distribution manifest =805entries. Every payload entry was
independently hashed against the manifest and frozen source148; no local state,
readiness/engineering approval, advisory response or Spec Kit record is included.
Historical PROJECT_FAILURE_LESSONS.md remains intentionally distributed evidence,
not an inherited approval. The complete extracted verifier passes53profiles and
162packs/1461materializable files/790Markdown, using real canonical code.

The actual guide0.1.2 from that archive runs13commands and8help requests across
NEW/EXISTING, including incomplete-input diagnosis and safe reentry. Two real
FileShare.Read replacement failures prove exact rollback; retry succeeds only
after releasing the owned lock. Original user bytes and checkpoint events remain
unchanged; successor append preserves its prefix and duplicate append is rejected.
The existing55check lifecycle regression then passes from the same archive.
All9consumed tooling/guide owners remain byte-identical in the current library.

The recovered harness findings are615 (missing synthetic blocker) and616
(native access-denied5 rather than sharing32); original failed sources remain.
No production readiness, independent human UX study, concurrent installer or
power-loss recovery is inferred. The candidate is unsigned and NOT FOR RELEASE.
TEST07 still requires security/licences/signature/independent release gates.
A changed release payload must be diffed and revalidated; these PASS results are
not automatically transferred to a future ZIP. Metadata-only evidence changes
in the library do not retroactively alter the frozen candidate or its receipts.

The original48test IDs are retained.41→43passed;7→5pending (all5blocked),
89.5833%/10.4167% of registered checks, not total effort. TEST02/03/05/07/09 stay
blocked. Benchmarks and11roadmap tasks are separate; no macro-task is completed
by these two candidate-specific acceptance results. Full Preflight149 pending.

| Acceptance evidence | SHA-256 |
|---|---|
| library-candidate148-not-for-release.zip | `0590588c4ab139e86e581b13c1db561f9342462646105c109f190aeefc755085` |
| library-candidate148-not-for-release.zip.sha256 | `47a5cd4fd3c2f00c1aeb9e21bfe0ca23fc258a1eee618c44a488944452756e99` |
| candidate148-archive.log | `7bc9cf0750f5b3e0955bb09c8a877cdae8d5c31d65f717de2329f1bed9c29493` |
| baseline148.json | `49efb61c88b0efb35883489b2463e6b07e106afdb4001cc9eb6d7b41e7f4bdd8` |
| accept-candidate.log | `720d4b3532febb619607456fee25fd8c8d1a15b7fdd480f9e3afb8160dfd975f` |
| canonical-acceptance.log | `b73704002ecbfba7aabbf2083811f6265f1aab58b6e9b81a973e598b7d5f2fc0` |
| canonical-green.log | `aace1066fd11fe00c78402ab374261aa279adc1399c67cf54530dafa0e301fb9` |
| canonical-green/acceptance.json | `39492902a888a790cbffe31d210221fd5a75757146c4c2894c84b2cb4731b821` |
| canonical-green/archive-verification.json | `efddef232819fd41409cf1420a785891568b314182ce1f434cf558833b6aefd5` |
| canonical-green/guide-harness-delta.json | `0fb0c0855ea33c783a23d678f63fd457aadae5ba9363ee0adc87fc3c9be20d9d` |
| canonical-green/guide-run/results.json | `9fd1008fa3acb3152da69650d8ae3281eecca7d04a6e278f83f313d1d4fa15ca` |
| canonical-green/extracted-structural.stdout | `90dea09d60b5607dab7fb973f31cdeacdb1eb52155170a5194be21e794d9572e` |
| canonical-green/candidate-entry-lifecycle.stdout | `4fd860795b93fc9ec33c613f43679e35de4d8de7e0001bc1cd52880f6fd32a4f` |
| accept_candidate352.py | `4fb0f0fd07742d2a006fd28b76982f4d310ca08fd86a00165bda8a5b7961cf4a` |
| accept_candidate352-before615.py | `cf8c96847bdd9722ef59afef8a013b4d50428c2aa27029c9b92520190bac2e6f` |
| accept_candidate352-before616.py | `1f06938df57cf4ee9d54a9e5b6ca31afd74c8ddc03ecf67fd9550cfbc794a19f` |

### Exact consumed source projection

| Source | SHA-256 |
|---|---|
| INSTALL_AGENT_BRIDGE.ps1 | `fe32bf216ea48e9f218dbcf3075161bef445a82bd498bea3f14429d172a12a4b` |
| materialize_markdown_pack.ps1 | `7ae045ff9b03ea22242d376d16b2a114f00e2f095217789eb92ca1f19e126356` |
| VERIFY_LIBRARY.ps1 | `addef5ca9f2f1076811ef1e019a6c56aca44857143092af4e95d3973f023d684` |
| markdown_system/PROJECT_ENTRY_CLI_GUIDE.md | `c6ecff9e199d1e67405e3abfc9e640ed5d3c2b77f3f075278e8533f6747e26fc` |
| markdown_system/test_agent_entry_lifecycle.ps1 | `158b9268658736db55a06ae4f22574827478b56cde1211aa1e76636b1feff21f` |
| implementation_packs/MARKDOWN_COMPOSITOR_CORE.md | `a8505384bb3f9085d36f21f1b6471fe24cbd17b054fa87f82f1d24b7007ce9d5` |
| implementation_packs/PROJECT_START_READINESS_VALIDATOR.md | `22c31cbc5ef0369767faac230ad5d03fb158b283434aeeb771fad3e13c547b9f` |
| implementation_packs/ENGINEERING_EXECUTION_VALIDATOR.md | `368c4f9d0c328353634e3dc0738e0d300f225013b9c7a769628b971d0c854808` |
| markdown_system/PROJECT_READINESS_GATE_PACK_PLAN.md | `f83d7a5133be17aa3e8a77a78de789085f6e8913db46ae98e6bf573bc6448adc` |

## Acceptance criteria retained

The delta audit confirms all48original test IDs, actions and oracles, unchanged
requirements, and only TEST06/08planned→passed. All162pack owners plus their
README (163directory files) remain byte-identical. Initial helper617 incorrectly
counted README as a pack; corrected from actual metadata, original error retained.
Evidence acceptance-delta.json SHA-256 `f3fbb96450e2795d1eb3a9438e5aedc3af733ab4e5d50741b4f5a1c1793a275f`.

## Closure150 — verified candidate acceptance and current library

Preflight149 finished: 154 executed checks PASS,162packs/1461files/
790Markdown/53profiles. Docker remains unavailable; skipped/conditional target
checks are not counted as executed. Checkpoint150 records the corrected status:
43/48passed, TEST02/03/05/07/09blocked, zero planned rows.5/48=10.4167% of checks
remain; this is not an estimate of effort. The roadmap still has1complete/10open
macro-tasks and benchmark BENCH01 remains pending independently of test-row count.

The accepted complete candidate148 archive retains SHA:
`0590588c4ab139e86e581b13c1db561f9342462646105c109f190aeefc755085`. All9consumed tooling/guide owners remain exact in the
current library; canonical harness source also reconstructs byte-identically.
The current library adds the audit evidence and successor checkpoint after that
freeze. Consequently it is not silently represented as the same release payload:
TEST07 must compare the final archive and revalidate changed inputs before release.
No final archive is delivered, signed or promoted by this closure.

FAIL613/615/616 recovered by actual acceptance;614 was a read timing error.
Original unsuccessful probes, candidate147, candidate148 and all receipts remain.
The remaining blockers require their original evidence; none is closed by changing
a label, omitting a requirement, accepting a synthetic provider as live or treating
local operational tests as installed production service acceptance.

| Closure evidence | SHA-256 |
|---|---|
| preflight149.json | `27cc3a0dae8a0503691dc311fe04e6fe727048f700fd396f268c8ac38d7d3aa6` |
| preflight149.log | `29fcad3c5f0e1b4ae51c216c7b0aff22ed4900e6b48a68d406900badab4d12b7` |
| closure149.json | `601ff30e84f12942269e8b866c54010f0b9b57e763face2b4c6ce01c3e14066a` |
| ledger149.json | `f6e30854d87fae4645c0fbe4e346a6873b9ab4be5913d79e89ee542e84267c8e` |
| resume149.log | `fc21fc48e58cc610f511cb26846e4873f75649c8937d068178012c42d422a380` |
| plan149.log | `2fd65a7ad27c8b968c64f8079abf8610894d6e7732b0b0f514dba962cf93b565` |
