"""AUTHORED packaging glue: bind actual runtime manifests to retained notices."""
from pathlib import Path
import hashlib,json,re

HEX=re.compile(r'[0-9a-f]{64}')
LEGAL=re.compile(r'(?:license|licence|notice|copying)(?:[-.].*)?|.*\.(?:LEGAL|LICENSE)\.txt',re.I)

def require(condition,message):
    if not condition:raise ValueError(message)

def sha(path):return hashlib.sha256(path.read_bytes()).hexdigest()

def no_duplicates(pairs):
    result={}
    for key,value in pairs:
        require(key not in result,'duplicate JSON field');result[key]=value
    return result

def read(path):
    require(path.is_file()and path.stat().st_size<=2000000,'bounded JSON required')
    return json.loads(path.read_bytes(),object_pairs_hook=no_duplicates)

def relative(root,value):
    require(isinstance(value,str)and value and '\\'not in value and ':'not in value,'invalid notice path')
    parts=value.split('/');require(all(x and x not in('.','..')for x in parts),'invalid notice path')
    p=root.joinpath(*parts)
    for item in [p,*p.parents]:
        if item.exists():require(not item.is_symlink()and not item.is_junction(),'reparse notice path')
        if item==root:break
    return p

def collect(source,runtime,output,lock_path):
    source=Path(source);runtime=Path(runtime);output=Path(output)
    for root in [source,runtime,output]:
        require(root.is_dir(),'existing notice root required')
        for p in [root,*root.parents]:require(not p.is_symlink()and not p.is_junction(),'reparse notice root')
    lock=read(Path(lock_path));require(set(lock)=={'schema','source_scope','manifests','texts','build_only_exclusions'},'notice lock fields')
    require(lock['schema']=='elite-reference-runtime-notices/v1','notice schema')
    require(isinstance(lock['manifests'],dict)and 0<len(lock['manifests'])<=1000,'manifest lock budget')
    staged={};total=0
    for rel,expected in lock['texts'].items():
        p=relative(source,rel);require(HEX.fullmatch(expected)and p.is_file()and sha(p)==expected,'notice text changed: '+rel)
        total+=p.stat().st_size;require(total<=16777216,'notice text budget')
        staged[rel]=p
    require(staged,'empty notice texts')
    actual={};seen=set()
    for p in (runtime/'node_modules').rglob('*'):
        require(not p.is_symlink()and not p.is_junction(),'runtime reparse entry')
        if not p.is_file():continue
        require(p.suffix.lower()not in('.node','.dll','.wasm'),'unadmitted native runtime payload')
        if p.name!='package.json':continue
        rel=p.relative_to(runtime).as_posix();require(rel.casefold()not in seen,'runtime case alias');seen.add(rel.casefold());actual[rel]=p
        require(len(actual)<=1000,'runtime manifest budget')
    require(set(actual)==set(lock['manifests']),'unknown or missing runtime manifest')
    rows=[]
    for rel,p in sorted(actual.items()):
        record=lock['manifests'][rel];require(set(record)=={'sha256','name','version','license','notice_refs','source_basis'},'manifest entry fields')
        require(sha(p)==record['sha256'],'runtime manifest changed: '+rel)
        original=relative(source,rel);require(original.is_file()and sha(original)==record['sha256'],'installed manifest correspondence')
        v=read(p);require((v.get('name'),v.get('version'),v.get('license'))==(record['name'],record['version'],record['license']),'runtime metadata differs')
        refs=record['notice_refs'];require(isinstance(refs,list)and refs and len(refs)==len(set(refs))and all(x in staged for x in refs),'runtime notice missing')
        require(any(staged[x].stat().st_size>100 for x in refs),'runtime notice empty')
        rows.append({'manifest':rel,**record,'texts':{x:lock['texts'][x]for x in refs}})
    identities={str(x['name'])+'@'+str(x['version'])for x in rows if x['name']and x['version']}
    require(not identities.intersection(lock['build_only_exclusions']),'excluded build tool was redistributed')
    installed=[]
    for scope in sorted((source/'node_modules/.pnpm').glob('*/node_modules/*')):
        roots=list(scope.iterdir())if scope.name.startswith('@')and scope.is_dir()else[scope]
        for root in roots:
            if root.is_symlink()or root.is_junction()or not(root/'package.json').is_file():continue
            v=read(root/'package.json');identity=str(v.get('name'))+'@'+str(v.get('version'));texts={}
            for p in sorted(root.iterdir()):
                if p.is_file()and LEGAL.fullmatch(p.name):
                    rel=p.relative_to(source).as_posix();require(rel in staged,'installed legal text not locked: '+rel);texts[p.name]=lock['texts'][rel]
            in_runtime=identity in identities
            if in_runtime and not texts:
                refs={r for row in rows if row['name']==v.get('name')and row['version']==v.get('version')for r in row['notice_refs']}
                texts={r:lock['texts'][r]for r in sorted(refs)}
            installed.append({'package':identity,'package_json_sha256':sha(root/'package.json'),'declared_license':v.get('license'),'texts':texts,'included_in_runtime':in_runtime,'scope':'REDISTRIBUTED_RUNTIME'if in_runtime else'BUILD_INPUT_NOT_REDISTRIBUTED'})
    # Validate the full graph before any output. Publication remains create-only.
    names=['dependency-notices.json','runtime-dependency-notices.json']
    require(all(not(output/n).exists()for n in names),'notice output already exists')
    dest=output/'dependency-notices';require(not dest.exists()or dest.is_dir()and not dest.is_symlink()and not dest.is_junction(),'invalid legal destination');dest.mkdir(exist_ok=True)
    for rel,p in sorted(staged.items()):
        target=dest/lock['texts'][rel]
        if target.exists():require(target.is_file()and not target.is_symlink()and sha(target)==lock['texts'][rel],'notice collision')
        else:
            with target.open('xb')as f:f.write(p.read_bytes())
    result={'schema':'elite-runtime-legal-catalogue/v1','state':'PASS_EXACT_NOTICE_CORRESPONDENCE','scope':lock['source_scope'],'runtime_manifests':rows,'notice_texts':lock['texts'],'build_only_exclusions':lock['build_only_exclusions'],'native_payloads':0}
    for name,value in zip(names,[installed,result]):
        with(output/name).open('x',encoding='utf-8',newline='\n')as f:f.write(json.dumps(value,sort_keys=True,indent=2,ensure_ascii=False)+'\n')
    return {'runtime_manifests':len(rows),'installed_identities':len(installed),'unique_notice_texts':len(set(lock['texts'].values())),'excluded_build_tools':len(lock['build_only_exclusions'])}
