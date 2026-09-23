"""AUTHORED final projection; only a fresh PASS receipt authorizes artifact copy."""
import argparse,hashlib,json
from pathlib import Path
from build_local_reference import plain,read_json,raw_json,tree,copy_tree,require

def publish(build,output,inventory,inventory_sha256):
    build=plain(build);output=plain(output)
    require(build!=output and build not in output.parents and output not in build.parents,'overlapping publication')
    require(output.is_dir() and not list(output.iterdir()),'empty publication directory required')
    source=read_json(inventory,inventory_sha256)
    source_id=hashlib.sha256(raw_json(source)).hexdigest()
    receipt=json.loads((build/'build-result.json').read_text(encoding='utf-8'))
    require(receipt['schema']=='elite-local-reference-build/v1' and receipt['scope']=='LOCAL_FIXTURES' and receipt['state']=='PASS','actual successful local build required')
    require(receipt['source_sha256']==source_id,'receipt source mismatch')
    actual=tree(build/'artifact')
    require(actual and actual==receipt['files'] and hashlib.sha256(raw_json(actual)).hexdigest()==receipt['artifact_sha256'],'artifact changed after build')
    require(all(actual.get('source/'+rel)==digest for rel,digest in source.items()),'corresponding source mismatch')
    require(json.loads((build/'artifact-inventory.json').read_text(encoding='utf-8'))==actual,'build inventory mismatch')
    copy_tree(build/'artifact',output)
    require(tree(output)==actual,'published artifact mismatch')
    print('REFERENCE_ARTIFACT_PROJECTED '+receipt['artifact_sha256'])

if __name__=='__main__':
    p=argparse.ArgumentParser(description=__doc__)
    for name in ['build','output','inventory','inventory-sha256']:p.add_argument('--'+name,required=True)
    args=p.parse_args()
    publish(args.build,args.output,args.inventory,args.inventory_sha256)
