"""AUTHORED local-only Next qualification; real BFF and owner URL supplied by Go."""
from pathlib import Path
import json,os,sys,socket,time,subprocess,urllib.request
root=Path(__file__).resolve().parents[1]
kind,api,identity=sys.argv[1:]
out=Path(os.environ['ELITE_FUNCTIONS_BROWSER_EVIDENCE'])/kind
out.mkdir(parents=True,exist_ok=False)
node=Path(os.environ['ELITE_NODE_BIN'])
with socket.socket() as s:s.bind(('127.0.0.1',0));port=s.getsockname()[1]
origin=f'http://127.0.0.1:{port}'
env={k:v for k,v in os.environ.items() if not k.startswith(('OIDC_','DOCUMENT_','NEXT_PUBLIC_'))}
env.update(ENTERPRISE_API_BASE_URL=api,APP_BASE_URL=origin,AUTH_SESSION_SECRET='synthetic-functions-local-session-32plus',NEXT_TELEMETRY_DISABLED='1',BUSINESS_CONFIG_FILE='functions-browser.fixture.json',CAMPAIGN_WORKSPACE_ENABLED='true',CATALOG_RELEASE_ENABLED='false',PUBLIC_SITE_ORIGIN='https://catalog.example.invalid',PUBLIC_INDEXING_ENABLED='0',ENTERPRISE_TENANT_CODE='functions-local',ENTERPRISE_ORGANIZATION_CODE='functions-local',FUNCTIONS_IDENTITY=identity,FUNCTIONS_BROWSER_ORIGIN=origin,FUNCTIONS_BROWSER_KIND=kind,FUNCTIONS_BROWSER_OUTPUT=str(out))
config=json.loads((root/'config/business.example.json').read_text('utf-8-sig'));config['features']['ai_activity']=True
(root/'config/functions-browser.fixture.json').write_text(json.dumps(config),encoding='utf-8')
flags=getattr(subprocess,'CREATE_NO_WINDOW',0)
server=None
with (out/'next.log').open('wb') as log:
 try:
  server=subprocess.Popen([str(node),str(root/'node_modules/next/dist/bin/next'),'dev','--webpack','--hostname','127.0.0.1','--port',str(port)],cwd=root,env=env,stdout=log,stderr=subprocess.STDOUT,creationflags=flags)
  deadline=time.monotonic()+100
  while True:
   assert server.poll() is None,'Next exited'
   try:
    with urllib.request.urlopen(origin+'/icon.svg',timeout=3) as r:break
   except Exception:
    assert time.monotonic()<deadline,'Next readiness timeout';time.sleep(.3)
  with (out/'browser.log').open('wb') as b:
   result=subprocess.run([str(node),str(root/'tools/final-functions-browser.mjs')],cwd=root,env=env,stdout=b,stderr=subprocess.STDOUT,timeout=150,creationflags=flags)
  if result.returncode:print((out/'browser.log').read_text('utf-8'));raise SystemExit(result.returncode)
 finally:
  if server and server.poll() is None:
   # Kill only this owned server tree; Next dev creates its own child.
   subprocess.run(['taskkill','/PID',str(server.pid),'/T','/F'],capture_output=True,creationflags=flags)
   server.wait(timeout=20)
print(out/'RESULT.json')
