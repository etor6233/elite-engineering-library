"""AUTHORED adverse publication/pre-build contracts, not full build evidence."""
import argparse,hashlib,json,subprocess,sys
from pathlib import Path
from publish_reference_build import publish
from build_local_reference import tree,raw_json,corresponding_sources

def test_runtime_junctions(node,target):
    """Actual Node loader regression: ZIP projection contains no empty scope dir."""
    target=Path(target).resolve();target.mkdir()
    package=target/'node_modules/.store/fixture';package.mkdir(parents=True)
    (package/'index.js').write_text('module.exports = 42;\n',encoding='utf-8')
    (target/'elite-runtime-links.json').write_bytes(raw_json({'node_modules/@fixture/module':'node_modules/.store/fixture'}))
    server=target/'server.js';server.write_text("if(require('@fixture/module')!==42)throw new Error('wrong module');process.on('SIGTERM',()=>process.exit(143));console.log('ZIP_JUNCTION_PASS');\n",encoding='utf-8')
    assert not(target/'node_modules/@fixture').exists()
    q=subprocess.run([node,str(Path(__file__).with_name('web_local_stop.cjs')),str(server)],input=b'STOP\n',stdout=subprocess.PIPE,stderr=subprocess.STDOUT,timeout=15)
    (target/'node.log').write_bytes(q.stdout)
    assert q.returncode==143 and b'ZIP_JUNCTION_PASS'in q.stdout,q.stdout.decode('utf-8','replace')
    assert(target/'node_modules/@fixture/module/index.js').read_bytes()==(package/'index.js').read_bytes()
    return {'state':'PASS','empty_parent_created':True,'real_node_module_resolution':True,'graceful_stop':True}

def run(target,pwsh):
    target=Path(target).resolve();target.mkdir()
    results=[];fixture_source={'a.txt':hashlib.sha256(b'original').hexdigest()}
    inventory=target/'inventory.json';inventory.write_bytes(raw_json(fixture_source))
    inventory_hash=hashlib.sha256(inventory.read_bytes()).hexdigest()
    def fixture(name):
        base=target/name;base.mkdir();build=base/'build';build.mkdir();artifact=build/'artifact';artifact.mkdir()
        (artifact/'source').mkdir();(artifact/'source/a.txt').write_bytes(b'original')
        output=base/'output';output.mkdir()
        files=tree(artifact);receipt={'schema':'elite-local-reference-build/v1','scope':'LOCAL_FIXTURES','state':'PASS','source_sha256':inventory_hash,'files':files,'artifact_sha256':hashlib.sha256(raw_json(files)).hexdigest()}
        (build/'build-result.json').write_bytes(raw_json(receipt));(build/'artifact-inventory.json').write_bytes(raw_json(files))
        return build,output,receipt
    for name in ['valid','changed-artifact','wrong-source','failed-receipt','occupied-output','changed-inventory','coherent-wrong-source']:
        build,output,receipt=fixture(name)
        if name=='changed-artifact':(build/'artifact/source/a.txt').write_bytes(b'modified')
        if name=='wrong-source':receipt['source_sha256']='0'*64
        if name=='failed-receipt':receipt['state']='FAIL'
        if name=='occupied-output':(output/'existing').write_text('keep')
        if name=='changed-inventory':(build/'artifact-inventory.json').write_text('{}')
        if name=='coherent-wrong-source':
            (build/'artifact/source/a.txt').write_bytes(b'modified');receipt['files']=tree(build/'artifact');receipt['artifact_sha256']=hashlib.sha256(raw_json(receipt['files'])).hexdigest();(build/'artifact-inventory.json').write_bytes(raw_json(receipt['files']))
        (build/'build-result.json').write_bytes(raw_json(receipt))
        try:publish(build,output,inventory,inventory_hash);accepted=True
        except ValueError:accepted=False
        assert accepted==(name=='valid'),name
        results.append({'case':name,'state':'PASS','accepted':accepted})
    for name in ['exact-source','generated-next-declaration','changed-domain-source','unknown-next-declaration']:
        base=target/name;base.mkdir();source=base/'input';built=base/'built';out=base/'artifact'
        for path in [source,built,out]:path.mkdir()
        declaration=b'/// <reference types="next/image-types/global" />\n'
        original={'next-env.d.ts':declaration,'domain.txt':b'original-domain'}
        for rel,data in original.items():(source/rel).write_bytes(data);(built/rel).write_bytes(data)
        files={rel:hashlib.sha256(data).hexdigest()for rel,data in original.items()}
        if name=='generated-next-declaration':(built/'next-env.d.ts').write_bytes(declaration+b'import "./.next/types/routes.d.ts";\nimport "./.next/types/root-params.d.ts";\n')
        if name=='changed-domain-source':(built/'domain.txt').write_bytes(b'changed')
        if name=='unknown-next-declaration':(built/'next-env.d.ts').write_bytes(b'unknown-generation')
        try:derived=corresponding_sources(source,built,out,files);accepted=True
        except ValueError:accepted=False
        assert accepted==(name in ['exact-source','generated-next-declaration']),name
        if accepted:
            assert tree(out/'source')==files
            assert bool(derived)==(name=='generated-next-declaration')
        results.append({'case':name,'state':'PASS','accepted':accepted})
    script=Path(__file__).with_name('build_signed_reference.ps1')
    for name in ['python-pin','source-pin','overlap','partial-build','already-complete']:
        base=target/name;base.mkdir();root=base/'project';root.mkdir();output=base/'output';output.mkdir();runroot=base/'session'
        inv=root/'inventory.json';inv.write_bytes(inventory.read_bytes())
        cfg={'schema_version':1,'scope':'LOCAL_FIXTURES','python':{'path':sys.executable,'sha256':hashlib.sha256(Path(sys.executable).read_bytes()).hexdigest()},'source_inventory_ref':'inventory.json','source_inventory_sha256':inventory_hash,'go':'unused-before-validation','node':'unused-before-validation','go_mod_cache':'unused-before-validation','install_inputs':'unused-before-validation','install_inputs_sha256':'0'*64,'arca_inputs':'unused-before-validation','arca_inputs_sha256':'0'*64,'run_root':str(runroot)}
        if name=='python-pin':cfg['python']['sha256']='0'*64
        if name=='source-pin':cfg['source_inventory_sha256']='0'*64
        if name=='overlap':cfg['run_root']=str(root/'overlap')
        config=root/'.elite-local-build-inputs.json';config.write_bytes(raw_json(cfg))
        if name in ['partial-build','already-complete']:
            runroot.mkdir();(runroot/'session.json').write_bytes(raw_json({'schema_version':1,'configuration_sha256':hashlib.sha256(config.read_bytes()).hexdigest(),'source_inventory_sha256':inventory_hash,'completed_builds':2 if name=='already-complete' else 0}))
            if name=='partial-build':(runroot/'build-1').mkdir()
        q=subprocess.run([pwsh,'-NoProfile','-File',str(script),'-ProjectRoot',str(root),'-OutputRoot',str(output),'-SourceDateEpoch','315532800'],stdout=subprocess.PIPE,stderr=subprocess.STDOUT,timeout=30)
        (base/'rejection.log').write_bytes(q.stdout)
        expected={'python-pin':'Python pin changed','source-pin':'Source inventory pin changed','overlap':'Overlapping build paths','partial-build':'Previous or partial build retained','already-complete':'already completed both builds'}[name]
        assert q.returncode!=0 and expected in q.stdout.decode('utf-8','replace'),(name,q.stdout.decode('utf-8','replace'))
        assert not list(output.iterdir()),name
        results.append({'case':name,'state':'PASS','accepted':False})
    result={'state':'PASS_CONTRACT_FIXTURES_ONLY','cases':results,'actual_product_build_performed':False}
    (target/'result.json').write_bytes(raw_json(result));print(json.dumps(result))

if __name__=='__main__':
    p=argparse.ArgumentParser(description=__doc__);p.add_argument('--target',required=True);p.add_argument('--pwsh',required=True);p.add_argument('--node');args=p.parse_args();run(args.target,args.pwsh)
    if args.node:print(json.dumps(test_runtime_junctions(args.node,Path(args.target)/'zip-runtime')))
