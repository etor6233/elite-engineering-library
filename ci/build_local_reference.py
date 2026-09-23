"""AUTHORED build glue for the admitted Windows local fixture composition.

No acquisition, live credential use, container build or production admission.
The externally supplied inventory digest is the authority for the source tree.
"""
from __future__ import annotations
import argparse,base64,hashlib,json,os,re,shutil,subprocess,sys,time
from pathlib import Path

GO_SHA='21761eceb9302062c9623fb699f332c8c7fe000f15f70efe8da01a2cfbbc16b9'
NODE_SHA='5c976096e04e5c2c1f091938926234cc9fbebfe9787ddd149351b3b0ecc707b5'
NODE_LICENSE_SHA='ed34dd8e3f0a78dbaf00d0444ce8e285b015b765379c2e17880455f70370f8e9'
HEX=re.compile(r'[0-9a-f]{64}')
def require(ok,why):
    if not ok:raise ValueError(why)
def long_path(p):
    p=os.path.abspath(p)
    return Path(p if p.startswith('\\\\?\\')or os.name!='nt'else '\\\\?\\'+p)
def sha(p):return hashlib.sha256(long_path(p).read_bytes()).hexdigest()
def raw_json(value):return (json.dumps(value,sort_keys=True,indent=2)+'\n').encode()
def read_json(path,expected):
    raw=long_path(path).read_bytes()
    require(HEX.fullmatch(expected) and hashlib.sha256(raw).hexdigest()==expected,'external digest mismatch')
    return json.loads(raw)
def plain(path):
    p=Path(os.path.abspath(path))
    for part in [p,*p.parents]:
        if part.exists():require(not part.is_symlink() and not part.is_junction(),'reparse path rejected')
    return p
def relative(root,rel):
    require(isinstance(rel,str) and '\\'not in rel and ':'not in rel and not rel.startswith('/'),'invalid relative path')
    require(all(x and x not in ('.','..')for x in rel.split('/')),'invalid relative path')
    p=plain(root/rel);p.relative_to(root);return p
def inventory(root,files):
    require(isinstance(files,dict)and files,'empty source inventory')
    seen=set()
    for rel,h in files.items():
        require(rel.casefold()not in seen,'case alias');seen.add(rel.casefold())
        p=relative(root,rel);require(HEX.fullmatch(h)and p.is_file()and sha(p)==h,'source changed: '+rel)
    return hashlib.sha256(raw_json(files)).hexdigest()
def environment(node,work):
    return {'SYSTEMROOT':os.environ['SYSTEMROOT'],'WINDIR':os.environ['SYSTEMROOT'],
      'PATHEXT':'.COM;.EXE;.BAT;.CMD','COMSPEC':str(Path(os.environ['SYSTEMROOT'])/'System32/cmd.exe'),
      'TEMP':str(work),'TMP':str(work),'PATH':str(node.parent)+os.pathsep+str(Path(os.environ['SYSTEMROOT'])/'System32'),
      'CI':'true','NEXT_TELEMETRY_DISABLED':'1','TZ':'UTC','GOTOOLCHAIN':'local','GOPROXY':'off',
      'GOSUMDB':'off','GOWORK':'off','GOFLAGS':'-mod=readonly','GOMAXPROCS':'4','CGO_ENABLED':'0',
      'GOOS':'windows','GOARCH':'amd64','BUSINESS_CONFIG_FILE':'business.example.json',
      'APP_BASE_URL':'https://127.0.0.1:4443','ENTERPRISE_API_BASE_URL':'http://127.0.0.1:9',
      'AUTH_SESSION_SECRET':'synthetic-build-only-not-production-00000000',
      'CATALOG_RELEASE_ENABLED':'false','PUBLIC_INDEXING_ENABLED':'0',
      # Public fixture encryption material makes LOCAL builds comparable.
      # It is forbidden as a live deployment credential.
      'NEXT_SERVER_ACTIONS_ENCRYPTION_KEY':base64.b64encode(hashlib.sha256(b'elite-local-reference-only-actions-v1').digest()).decode()}
def command(label,argv,cwd,env,logs,seconds=300):
    p=logs/(label+'.log');start=time.monotonic()
    with p.open('xb')as f:
        q=subprocess.run(list(map(str,argv)),cwd=cwd,env=env,stdin=subprocess.DEVNULL,stdout=f,stderr=subprocess.STDOUT,
                         timeout=seconds,shell=False,creationflags=subprocess.CREATE_NO_WINDOW)
    require(q.returncode==0,label+' failed; see '+str(p))
    return {'step':label,'exit_code':q.returncode,'log_sha256':sha(p),'seconds':round(time.monotonic()-start,3)}
def tree(root):
    root=long_path(root)
    result={}
    for p in sorted(root.rglob('*')):
        require(not p.is_symlink()and not p.is_junction(),'artifact link rejected')
        if p.is_file():result[p.relative_to(root).as_posix()]=sha(p)
    return result
def copy_tree(source,target):
    # Runtime copies contain regular files only; junctions are created afterward
    # from the exact framework trace, inside the private activation directory.
    source=long_path(source);target=long_path(target);count=0;total=0
    def copy(a,b):
        nonlocal count,total
        b.mkdir(parents=True,exist_ok=True)
        with os.scandir(a)as entries:
            for entry in entries:
                info=entry.stat(follow_symlinks=False)
                require(not getattr(info,'st_file_attributes',0)&0x400,'runtime copy contains a link')
                if entry.is_dir(follow_symlinks=False):copy(Path(entry.path),b/entry.name)
                else:
                    count+=1;total+=info.st_size
                    require(count<=50000 and total<=536870912,'runtime copy budget exceeded')
                    shutil.copyfile(entry.path,b/entry.name)
    copy(source,target)
def project_standalone(src,target):
    source=long_path(src/'.next/standalone');target=long_path(target);installed=long_path(src/'node_modules').resolve();links={};files=0
    def copy(a,b):
        nonlocal files
        if a.is_symlink()or a.is_junction():
            actual=a.resolve(strict=True)
            if source in actual.parents:rel=actual.relative_to(source).as_posix()
            else:
                require(installed in actual.parents,'framework link escapes qualified consumer')
                rel='node_modules/'+actual.relative_to(installed).as_posix()
            require(actual.is_dir(),'only directory junction projection supported')
            links[a.relative_to(source).as_posix()]=rel;return
        if a.is_dir():
            b.mkdir(parents=True,exist_ok=True)
            for child in sorted(a.iterdir()):copy(child,b/child.name)
        else:
            files+=1;require(files<=20000,'standalone file budget');shutil.copyfile(a,b)
    copy(source,target)
    for rel,dest in links.items():require((target/dest).is_dir(),'traced link target absent: '+rel)
    (target/'elite-runtime-links.json').write_bytes(raw_json(links))
def go_notices(src,dist,go,env,logs):
    step=command('go-modules',[go,'list','-m','-json','all'],src,env,logs)
    raw=(logs/'go-modules.log').read_text(encoding='utf-8');decoder=json.JSONDecoder();modules=[]
    while raw.strip():
        value,n=decoder.raw_decode(raw.lstrip());raw=raw.lstrip()[n:];modules.append(value)
    result=[]
    for value in modules:
        if value.get('Main'):continue
        spec=value.get('Replace',value);root=Path(spec.get('Dir',''));texts={}
        require(root.is_dir(),'module source absent')
        for item in root.iterdir():
            if item.is_file()and re.fullmatch(r'(?:license|licence|notice|copying)(?:[.-].*)?',item.name,re.I):
                h=sha(item);p=dist/'dependency-notices'/h;p.parent.mkdir(exist_ok=True)
                if not p.exists():shutil.copyfile(item,p)
                texts[item.name]=h
        result.append({'module':value['Path'],'version':spec.get('Version'),'sum':spec.get('Sum'),'texts':texts,
          'local_corresponding_source':bool(value.get('Replace'))})
    p=go.parent.parent/'LICENSE';h=sha(p);shutil.copyfile(p,dist/'GO_LICENSE.txt')
    (dist/'go-dependency-notices.json').write_bytes(raw_json({'modules':result,'go_license_sha256':h}))
    return step
def corresponding_sources(source,built,dist,files):
    # Keep exact admitted inputs. Next writes generated type declarations in
    # its workspace; retain that one checked derivation separately.
    inventory(source,files);derived={}
    for rel,h in files.items():
        actual=sha(built/rel)
        if actual!=h:
            require(rel=='next-env.d.ts','unexpected build source mutation: '+rel)
            original=(source/rel).read_bytes();anchor=b'/// <reference types="next/image-types/global" />\n'
            addition=b'import "./.next/types/routes.d.ts";\nimport "./.next/types/root-params.d.ts";\n'
            require(original.count(anchor)==1 and (built/rel).read_bytes()==original.replace(anchor,anchor+addition),
                    'unknown Next generated declaration')
            dest=dist/'build-derived-source'/rel;dest.parent.mkdir(parents=True,exist_ok=True);shutil.copyfile(built/rel,dest)
            require(sha(dest)==actual,'derived source copy drift')
            derived[rel]={'input_sha256':h,'generated_sha256':actual,'generator':'Next16.3.4 canonical route type declaration'}
        dest=dist/'source'/rel;dest.parent.mkdir(parents=True,exist_ok=True);shutil.copyfile(source/rel,dest)
        require(sha(dest)==h,'corresponding source copy drift')
    if derived:(dist/'build-derived-source/manifest.json').write_bytes(raw_json(derived))
    inventory(source,files)
    return derived

def package_artifact(src,dist,node,source_id,files,source):
    standalone=src/'.next/standalone';require((standalone/'server.js').is_file(),'standalone frontend missing')
    project_standalone(src,dist/'web')
    for rel in ['.next/static','public','config']:
        if(src/rel).is_dir():copy_tree(src/rel,dist/'web'/rel)
    shutil.copytree(src/'db/migrations',dist/'migrations')
    corresponding_sources(source,src,dist,files)
    for rel in files:
        if rel.startswith('licenses/')or rel.startswith('docs/provenance/')or rel=='THIRD_PARTY_NOTICES.md':
            p=dist/rel;p.parent.mkdir(parents=True,exist_ok=True);shutil.copyfile(source/rel,p)
    require(sha(node.parent/'LICENSE')==NODE_LICENSE_SHA,'Node license pin changed')
    shutil.copyfile(node,dist/'node.exe');shutil.copyfile(node.parent/'LICENSE',dist/'NODE_LICENSE.txt')
    from collect_runtime_notices import collect
    collect(src,dist/'web',dist,src/'licenses/reference/runtime-notice-lock.json')
    # Framework-generated metadata contains build-machine absolute paths.
    # Only these documented non-serving locator fields are normalized;
    # receipts keep original hashes. Never rewrite executable app bundles.
    normalization=[]
    for nft in sorted(long_path(dist/'web/.next').rglob('*.nft.json')):
        value=json.loads(nft.read_text());require(set(value)=={'version','files'}and value['version']==1,'unknown Next trace format')
        require(isinstance(value['files'],list)and all(isinstance(x,str)for x in value['files']),'invalid Next trace files')
        normalization.append({'path':nft.relative_to(long_path(dist)).as_posix(),'before_sha256':sha(nft),'field':'files sort order'})
        value['files']=sorted(value['files']);nft.write_bytes(raw_json(value))
    p=dist/'web/server.js';text=p.read_text();m=re.search(r'^const nextConfig = (.+)$',text,re.M)
    require(m is not None,'unknown Next standalone server format')
    config=json.loads(m[1]);config['outputFileTracingRoot']='.';config['repoRoot']='.';config['turbopack']['root']='.'
    normalization.append({'path':'web/server.js','before_sha256':sha(p),'fields':['nextConfig.outputFileTracingRoot','nextConfig.repoRoot','nextConfig.turbopack.root']})
    p.write_text(text[:m.start(1)]+json.dumps(config,separators=(',',':'))+text[m.end(1):],encoding='utf-8',newline='\n')
    p=dist/'web/.next/required-server-files.json'
    if p.exists():
        v=json.loads(p.read_text());v['appDir']='.';v['config']['outputFileTracingRoot']='.';v['config']['repoRoot']='.';v['config']['turbopack']['root']='.'
        normalization.append({'path':'web/.next/required-server-files.json','before_sha256':sha(p),'fields':['appDir','config.outputFileTracingRoot','config.repoRoot','config.turbopack.root']})
        p.write_bytes(raw_json(v))
    (dist/'LOCAL_REFERENCE_ONLY.json').write_bytes(raw_json({'schema':'elite-local-reference-release/v1',
      'source_sha256':source_id,'scope':'LOCAL_FIXTURES','production_admitted':False,
      'fixture_credentials':True,'go_sha256':GO_SHA,'node_sha256':NODE_SHA}))
    return normalization
def build(args):
    source=plain(args.source);target=plain(args.target);workspace=plain(args.workspace);go=plain(args.go);node=plain(args.node)
    require(not target.exists()and target.parent.is_dir(),'absent target required')
    require(source!=target and source not in target.parents and target not in source.parents,'overlapping source/output')
    require(not workspace.exists()and workspace.parent.is_dir(),'absent stable build workspace required')
    for other in [source,target]:require(workspace!=other and workspace not in other.parents and other not in workspace.parents,'build workspace overlaps input/output')
    files=read_json(args.inventory,args.inventory_sha256);source_id=inventory(source,files)
    require(sha(go)==GO_SHA and sha(node)==NODE_SHA,'toolchain pin mismatch')
    plan_inputs=read_json(args.install_inputs,args.install_inputs_sha256)
    require(set(plan_inputs)=={'projection','contained','store','cache','acquisition_receipt','acquisition_sha256'},'exact install input fields required')
    runtime=source/'pnpm_artifact_selection';require(runtime.is_dir(),'materialize PNPM_ARTIFACT_SELECTION_GATE first')
    target.mkdir();work=target/'work';work.mkdir();logs=target/'logs';logs.mkdir();src=workspace;src.mkdir()
    result={'schema':'elite-local-reference-build/v1','scope':'LOCAL_FIXTURES','state':'FAIL','source_sha256':source_id,'build_workspace':str(workspace),'steps':[]}
    try:
        for rel in files:
            p=src/rel;p.parent.mkdir(parents=True,exist_ok=True);shutil.copyfile(source/rel,p)
        # Execute the existing exact three-consumer admission, never raw pnpm.
        sys.path.insert(0,str(runtime));import local_runtime
        inputs={**plan_inputs,'consumer':'enterprise-web','project':str(src),'node':str(node),'target':str(target/'install')}
        plan_sha=local_runtime.prepare(**inputs);installed=local_runtime.execute(target/'install',plan_sha)
        require(installed['state']=='PASS','restricted install failed')
        result['install']={'plan_sha256':plan_sha,'receipt_sha256':sha(target/'install/execution-result.json')}
        env=environment(node,work);env.update(GOMODCACHE=str(plain(args.go_mod_cache)),GOCACHE=str(target/'go-cache'),
            USERPROFILE=str(work),ELITE_SOURCE_SHA256=source_id)
        dist=target/'artifact';dist.mkdir();(dist/'api').mkdir()
        from build_arca_reference import build_arca
        result['arca']=build_arca(src,dist,work,logs,go,env,args.arca_inputs,args.arca_inputs_sha256)
        result['steps'].append(command('api',[go,'build','-trimpath','-buildvcs=false','-ldflags=-s -w -buildid=',
                    '-o',dist/'api/electromobility-api.exe','./cmd/electromobility-api'],src,env,logs))
        # Next16.3.4 preserves its preview cache during build. These public
        # fixture values are never production credentials or runtime admission.
        preview={name:hashlib.sha256(('elite-local-only:'+name).encode()).hexdigest()[:size]
            for name,size in [('previewModeId',32),('previewModeSigningKey',64),('previewModeEncryptionKey',64)]}
        preview['expireAt']=4102444800000
        cache=src/'.next/cache';cache.mkdir(parents=True)
        (cache/'.previewinfo').write_bytes(raw_json(preview))
        result['steps'].append(command('web',[node,src/'node_modules/next/dist/bin/next','build','--webpack'],src,env,logs))
        actual=json.loads((src/'.next/prerender-manifest.json').read_text())['preview']
        require(actual=={k:v for k,v in preview.items()if k!='expireAt'},'Next preview fixture contract changed')
        normalization=package_artifact(src,dist,node,source_id,files,source)
        result['steps'].append(go_notices(src,dist,go,env,logs))
        result['normalization']=normalization;result['files']=tree(dist)
        result['artifact_sha256']=hashlib.sha256(raw_json(result['files'])).hexdigest()
        inventory(source,files);result['state']='PASS'
        (target/'artifact-inventory.json').write_bytes(raw_json(result['files']))
    finally:(target/'build-result.json').write_bytes(raw_json(result))
    return result
def main():
    p=argparse.ArgumentParser(description=__doc__)
    for name in ['source','target','workspace','inventory','inventory-sha256','go','node','go-mod-cache','install-inputs','install-inputs-sha256','arca-inputs','arca-inputs-sha256']:p.add_argument('--'+name,required=True)
    a=p.parse_args()
    try:r=build(a);print('LOCAL_REFERENCE_BUILD_PASS '+r['artifact_sha256']);return 0
    except (ValueError,OSError,subprocess.SubprocessError)as e:print('LOCAL_REFERENCE_BUILD_FAIL '+str(e),file=sys.stderr);return 2
if __name__=='__main__':raise SystemExit(main())
