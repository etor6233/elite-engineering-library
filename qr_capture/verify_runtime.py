from pathlib import Path
import json,hashlib
def verify(runtime):
 runtime=runtime.resolve(strict=True);path=Path(__file__).with_name('runtime-file-lock.json');lock=json.loads(path.read_text(encoding='utf8'));expected={x['path']:x['sha256'] for x in lock['files']}
 for rel,digest in expected.items():
  p=runtime.joinpath(*rel.split('/'))
  if not p.resolve().is_relative_to(runtime) or p.is_symlink() or not p.is_file() or hashlib.sha256(p.read_bytes()).hexdigest()!=digest:raise ValueError('runtime file mismatch: '+rel)
 for p in runtime.rglob('*'):
  if p.is_file() and p.suffix.lower() in {'.py','.pyc','.pyd','.dll','.so','.pth'} and p.relative_to(runtime).as_posix() not in expected:raise ValueError('unexpected runtime executable source: '+p.relative_to(runtime).as_posix())
 return {'schema':'elite-qr-runtime-verification/v1','status':'PASS_EXACT_WHEEL_FILES','files_verified':len(expected),'lock_sha256':hashlib.sha256(path.read_bytes()).hexdigest(),'runtime':str(runtime),'platform_scope':'Windows x64 CPython3.14','production_authorized':False}
if __name__=='__main__':
 import argparse
 p=argparse.ArgumentParser();p.add_argument('runtime',type=Path);a=p.parse_args();print(json.dumps(verify(a.runtime)))
