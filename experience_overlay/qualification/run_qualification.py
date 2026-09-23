"""AUTHORED bounded loopback browser qualification; no cloud/live/physical-device claim."""
from pathlib import Path
import argparse, base64, hashlib, json, os, secrets, shutil, subprocess, sys, time, urllib.request
p=argparse.ArgumentParser();p.add_argument('--target',required=True);p.add_argument('--evidence',required=True);p.add_argument('--node',default='node');p.add_argument('--openssl',default='openssl');p.add_argument('--candidate',action='store_true');p.add_argument('--attach-existing',action='store_true');p.add_argument('--grep');p.add_argument('--grep-invert');a=p.parse_args()
r=Path(__file__).resolve().parent;t=Path(a.target).resolve();e=Path(a.evidence).resolve();e.mkdir(parents=True,exist_ok=True)
if (t/'AGENT_SYSTEM_START.md').exists():raise SystemExit('Canonical library is not an executable fixture target')
artifact=t/'.next/standalone'
if not (artifact/'server.js').is_file():raise SystemExit('Build standalone artifact before qualification')
if a.attach_existing and not (e/'ephemeral-session.key').is_file():raise SystemExit('Attach requires the existing ephemeral session file in evidence directory')
env=os.environ.copy();env.update(UI_TARGET=str(t),UI_EVIDENCE=str(e),UI_BROWSER_REPORT=str(e/'browser-tests.json'))
processes=[];handles=[];backup={};lock=e/'.qualification.lock';lock.open('x').close();result={'schema_version':'1.0.0','revision':'V403-UI-0.3.0','scope':'LOCAL_LOOPBACK_AUTHORED_BACKEND_FIXTURE_REAL_NEXT_BFF','status':'FAIL','visual_baseline_approval':'PENDING_USER_APPROVAL','test_filter':{'grep':a.grep,'grep_invert':a.grep_invert},'production_authorized':False,'build_entry_sha256':hashlib.sha256((artifact/'server.js').read_bytes()).hexdigest(),'limits':['Not Go/PostgreSQL business-owner acceptance','Not physical phone, screen reader, browser zoom or user acceptance','No cloud execution or authorization']}
try:
 if not a.attach_existing:
  shutil.copytree(t/'.next/static',artifact/'.next/static',dirs_exist_ok=True)
  for name in ['experience-checklists.json','experience-options.json','experience.browser.fixture.json']:
   dest=artifact/'config'/name;backup[dest]=dest.read_bytes() if dest.exists() else None
  (artifact/'config').mkdir(exist_ok=True)
  subprocess.run([sys.executable,str(r/'fixture_config.py')],env=env,check=True)
  subprocess.run([a.openssl,'req','-x509','-newkey','rsa:2048','-nodes','-keyout',str(e/'fixture.key'),'-out',str(e/'fixture.crt'),'-days','2','-subj','/CN=localhost','-addext','subjectAltName=IP:127.0.0.1,DNS:localhost'],stdout=subprocess.DEVNULL,stderr=subprocess.DEVNULL,check=True)
  key=secrets.token_hex(48);(e/'ephemeral-session.key').write_text(key);os.chmod(e/'ephemeral-session.key',0o600)
  env.update(PORT='4311',HOSTNAME='127.0.0.1',ELITE_EXPERIENCE_CATALOGUE='1',DOCUMENTS_ENABLED='true',DOCUMENTS_PROFILE_SHA256='a'*64,DOCUMENTS_TENANT_ID='tenant',DOCUMENTS_ORGANIZATION_ID='store',DOCUMENTS_MODE='FIXTURE',BUSINESS_CONFIG_FILE='experience.browser.fixture.json',APP_BASE_URL='https://127.0.0.1:4313',ENTERPRISE_API_BASE_URL='http://127.0.0.1:4314',ENTERPRISE_TENANT_CODE='elite-mobility',ENTERPRISE_ORGANIZATION_CODE='store',AUTH_SESSION_SECRET=key)
  for name,command,cwd in [('backend',[a.node,str(r/'fixture-server.mjs')],r),('frontend',[a.node,'server.js'],artifact)]:
   output=(e/(name+'.log')).open('w');handles.append(output);processes.append(subprocess.Popen(command,cwd=cwd,env=env,stdout=output,stderr=subprocess.STDOUT))
  deadline=time.monotonic()+20
  while True:
   try:
    urllib.request.urlopen('http://127.0.0.1:4311/experience',timeout=1).close();urllib.request.urlopen('http://127.0.0.1:4314/fixture/stats',timeout=1).close();break
   except Exception:
    if any(x.poll() is not None for x in processes) or time.monotonic()>deadline:raise RuntimeError('Loopback fixture processes did not become ready; see logs')
    time.sleep(.2)
 urllib.request.urlopen(urllib.request.Request('http://127.0.0.1:4314/fixture/reset',data=b'',method='POST'),timeout=2).close()
 command=[a.node,os.environ.get('PLAYWRIGHT_CLI') or str(t/'microsoft_playwright_browser_gate/node_modules/@playwright/test/cli.js'),'test','--config',str(r/'playwright.v403.config.mjs')]
 if a.candidate:command.append('--update-snapshots')
 if a.grep:command.extend(['--grep',a.grep])
 if a.grep_invert:command.extend(['--grep-invert',a.grep_invert])
 completed=subprocess.run(command,cwd=r,env=env)
 result['exit_code']=completed.returncode;result['status']='PASS' if completed.returncode==0 else 'FAIL'
 report=e/'browser-tests.json';result['browser_report_sha256']=hashlib.sha256(report.read_bytes()).hexdigest() if report.exists() else None
 result['harness_files']=[{'path':str(x.relative_to(r)).replace('\\','/'),'sha256':hashlib.sha256(x.read_bytes()).hexdigest()}for x in sorted(r.rglob('*')) if x.is_file() and x.suffix in ['.py','.mjs','.json'] and not x.resolve().is_relative_to(e) and 'evidence' not in x.parts]
 result['candidate_screenshots']=[{'path':str(x.relative_to(r)).replace('\\','/'),'sha256':hashlib.sha256(x.read_bytes()).hexdigest()}for x in sorted((r/'tests/visual-candidates').glob('*.png'))]
finally:
 for process in processes:
  if process.poll() is None:process.terminate()
 for process in processes:
  try:process.wait(timeout=5)
  except subprocess.TimeoutExpired:process.kill();process.wait()
 for handle in handles:handle.close()
 for path,data in backup.items():
  if data is None:path.unlink(missing_ok=True)
  else:path.write_bytes(data)
 lock.unlink(missing_ok=True);(e/'qualification-receipt.json').write_text(json.dumps(result,indent=2)+'\n')
print(json.dumps({'status':result['status'],'receipt':str(e/'qualification-receipt.json')}))
raise SystemExit(0 if result['status']=='PASS' else 1)
