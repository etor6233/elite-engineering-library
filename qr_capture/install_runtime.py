"""Install only exact admitted local-scope wheels into an absent runtime folder.
No system Python modification. Runtime paths are trusted deployment config.
"""
from pathlib import Path
import json,hashlib,urllib.request,subprocess,sys,argparse
def install(destination,wheel_directory):
 here=Path(__file__).resolve().parent;destination=destination.resolve();wheel_directory=wheel_directory.resolve()
 if destination.exists():raise FileExistsError('runtime destination must be absent')
 wheel_directory.mkdir(parents=True,exist_ok=True)
 lock=json.loads((here/'runtime-package-lock.json').read_text(encoding='utf8'))
 for row in lock['packages']:
  target=wheel_directory/row['filename']
  if not target.exists():
   if not row['url'].startswith('https://files.pythonhosted.org/'):raise ValueError('unapproved wheel origin')
   with urllib.request.urlopen(row['url'],timeout=60)as response:data=response.read(row['bytes']+1)
   if len(data)!=row['bytes'] or hashlib.sha256(data).hexdigest()!=row['sha256']:raise ValueError('wheel bytes mismatch')
   target.write_bytes(data)
  if hashlib.sha256(target.read_bytes()).hexdigest()!=row['sha256']:raise ValueError('cached wheel SHA mismatch')
 args=[sys.executable,'-m','pip','install','--no-compile','--no-cache-dir','--no-index','--find-links',str(wheel_directory),'--require-hashes','--only-binary=:all:','--target',str(destination),'-r',str(here/'requirements.windows-py314.lock')]
 subprocess.run(args,check=True,timeout=180)
 from verify_runtime import verify
 receipt=verify(destination);(destination/'ELITE_RUNTIME_INSTALL_RECEIPT.json').write_text(json.dumps(receipt,indent=2)+'\n',encoding='utf8');return receipt
if __name__=='__main__':
 p=argparse.ArgumentParser();p.add_argument('--destination',type=Path,required=True);p.add_argument('--wheel-directory',type=Path,required=True);a=p.parse_args();print(json.dumps(install(a.destination,a.wheel_directory)))
