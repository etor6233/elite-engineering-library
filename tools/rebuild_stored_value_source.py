"""AUTHORED rebuild/receipt glue for a user's separately modified LGPL module.
No source changes are made here and the result is not automatically admitted.
"""
import argparse
import hashlib
import json
from pathlib import Path
import re


def sha(raw):return hashlib.sha256(raw).hexdigest()


def rebuild(source,expected,output):
    source=source.resolve();output=output.resolve()
    if output.exists() or output.is_relative_to(source) or not re.fullmatch('[0-9a-f]{64}',expected):
        raise ValueError('an absent destination outside the source module is required')
    raw=(source/'engine-lock.json').read_bytes()
    if len(raw)>65536 or sha(raw)!=expected:raise ValueError('baseline source manifest mismatch')
    lock=json.loads(raw)
    if lock['schema']!='elite.odoo-derived-source-lock.v1' or not 1<=len(lock['files'])<=128:raise ValueError('source lock contract')
    snapshots={};changes=[];total=0
    for row in lock['files']:
        rel=row['path'];part=Path(rel)
        if not isinstance(rel,str) or part.is_absolute() or '..' in part.parts or rel in snapshots:raise ValueError('source path contract')
        p=(source/part).resolve()
        if not p.is_relative_to(source):raise ValueError('source path escapes module')
        data=p.read_bytes();total+=len(data)
        if total>8*1024*1024:raise ValueError('source module exceeds rebuild budget')
        after=sha(data)
        if after!=row['sha256'] or len(data)!=row['bytes']:changes.append({'path':rel,'before_sha256':row['sha256'],'after_sha256':after})
        snapshots[rel]=data;row.update(bytes=len(data),sha256=after)
    if not {'LICENSE','COPYRIGHT','run.py','engine.py','orm_contract.py','protocol.py'}<=snapshots.keys():raise ValueError('incomplete licensed module')
    # Preserve license texts/attribution when making a derivative snapshot.
    if any(r['path'] in ('LICENSE','COPYRIGHT','upstream/LICENSE','upstream/COPYRIGHT') for r in changes):raise ValueError('preserve original license and attribution texts')
    output.mkdir(parents=True,exist_ok=False)
    for rel,data in snapshots.items():
        p=output/rel;p.parent.mkdir(parents=True,exist_ok=True);p.write_bytes(data)
    encoded=(json.dumps(lock,indent=2)+'\n').encode('utf-8');(output/'engine-lock.json').write_bytes(encoded)
    receipt={'schema':'elite.user-source-replacement.v1','classification':'USER_MODIFIED_SOURCE_NOT_ADMITTED' if changes else 'EXACT_REBUILD','baseline_lock_sha256':expected,'replacement_lock_sha256':sha(encoded),'changes':changes,'next_step':'Retain LGPL sources/notices. Review and test the change, then explicitly bind its script/manifest hashes in the deployment. No upstream authorship or prior gate status is inherited.'}
    (output/'source-replacement-receipt.json').write_text(json.dumps(receipt,indent=2)+'\n',encoding='utf-8',newline='\n')
    return receipt


if __name__=='__main__':
    p=argparse.ArgumentParser(description=__doc__)
    p.add_argument('--module-dir',type=Path,required=True);p.add_argument('--expected-lock-sha256',required=True);p.add_argument('--destination',type=Path,required=True)
    a=p.parse_args()
    try:print(json.dumps(rebuild(a.module_dir,a.expected_lock_sha256,a.destination),indent=2))
    except (OSError,ValueError,KeyError,TypeError) as e:p.exit(2,str(e)+'\n')
