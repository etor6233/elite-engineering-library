"""AUTHORED offline signature/path fixtures, never a real release/SCA result."""
from pathlib import Path
import argparse,tempfile,subprocess,hashlib,json,os,shutil,zipfile,time

def main():
    parser=argparse.ArgumentParser();parser.add_argument('--pwsh',required=True);parser.add_argument('--ssh-keygen',required=True);parser.add_argument('--evidence-root',required=True);args=parser.parse_args()
    root=Path(args.evidence_root).absolute();assert not root.exists()and root.parent.is_dir();root.mkdir();gate=Path(__file__).parent;sha=lambda p:hashlib.sha256(p.read_bytes()).hexdigest()
    def write(path,value):path.write_text(json.dumps(value,sort_keys=True,indent=2)+'\n',encoding='utf-8',newline='\n')
    def run(label,argv,expected=0):
        with(root/(label+'.log')).open('xb')as out:q=subprocess.run(list(map(str,argv)),stdout=out,stderr=subprocess.STDOUT,stdin=subprocess.DEVNULL,timeout=60,creationflags=subprocess.CREATE_NO_WINDOW)
        assert(q.returncode==0)==(expected==0),(label,q.returncode);return(root/(label+'.log')).read_text(encoding='utf-8',errors='replace')
    probe=root/'probe.ps1';probe.write_text(r'''param([string]$Gate,[string]$Root)
$ErrorActionPreference='Stop'
$ast=[Management.Automation.Language.Parser]::ParseFile((Join-Path $Gate 'build_release.ps1'),[ref]$null,[ref]$null)
foreach($name in @('Fail','Sha','ResolveProjectFile','AssertSourceUnchanged','WriteDeterministicZip','GetVulnerabilityIds','PackageKey','AssertSbomCoverage','AddUnversionedSbomRecords')){
 $nodes=@($ast.FindAll({param($n)$n-is[Management.Automation.Language.FunctionDefinitionAst]-and$n.Name-ceq$name},$false));if($nodes.Count-ne1){throw 'missing actual gate helper'};. ([scriptblock]::Create($nodes[0].Extent.Text))
}
$clean='{"results":[{"packages":[{"package":{"name":"SYNTHETIC-ONLY","version":"0"}}]}]}'|ConvertFrom-Json
if(@(GetVulnerabilityIds $clean).Count-ne0){throw 'clean package counted as vulnerability'}
$clean.results[0].packages[0]|Add-Member -NotePropertyName vulnerabilities -NotePropertyValue @([pscustomobject]@{id='TEST-ONLY-1'})
if(@(GetVulnerabilityIds $clean).Count-ne1-or@(GetVulnerabilityIds $clean)[0]-cne'TEST-ONLY-1'){throw 'real vulnerability record omitted'}
$clean.results[0].packages[0].vulnerabilities=@([pscustomobject]@{id=''})
$rejected=$false;try{GetVulnerabilityIds $clean|Out-Null}catch{$rejected=$true};if(-not$rejected){throw 'malformed vulnerability accepted'}
$report='{"results":[{"packages":[{"package":{"name":"@fixture/example","version":"1","ecosystem":"npm"}},{"package":{"name":"./local","version":"","ecosystem":"Go"}}]}]}'|ConvertFrom-Json
$sbom='{"spdxVersion":"SPDX-2.3","SPDXID":"SPDXRef-DOCUMENT","packages":[{"name":"example","versionInfo":"1","externalRefs":[{"referenceType":"purl","referenceLocator":"pkg:npm/%40fixture/example@1"}]}],"relationships":[]}'|ConvertFrom-Json
$rejected=$false;try{AssertSbomCoverage $report $sbom}catch{$rejected=$true};if(-not$rejected){throw 'missing unversioned source accepted'}
AddUnversionedSbomRecords $report $sbom;AssertSbomCoverage $report $sbom
$sbom.packages[0].externalRefs[0].referenceLocator='pkg:npm/%40fixture/example@2'
$rejected=$false;try{AssertSbomCoverage $report $sbom}catch{$rejected=$true};if(-not$rejected){throw 'wrong version accepted'}
$rejected=$false;try{AssertSbomCoverage ([pscustomobject]@{results=@()}) $sbom}catch{$rejected=$true};if(-not$rejected){throw 'empty inventory accepted'}
$tree=Join-Path $Root 'tree';New-Item -ItemType Directory -Path (Join-Path $tree 'z')|Out-Null
[IO.File]::WriteAllText((Join-Path $tree 'a.txt'),"synthetic-a`n",[Text.UTF8Encoding]::new($false))
[IO.File]::WriteAllText((Join-Path $tree 'z/data.txt'),"synthetic-b`n",[Text.UTF8Encoding]::new($false))
$entry=@{path='a.txt';sha256=(Sha (Join-Path $tree 'a.txt'));size=(Get-Item -LiteralPath (Join-Path $tree 'a.txt')).Length}
AssertSourceUnchanged $tree @($entry)
$cfg=@{artifact_zip_timestamp='1980-01-01T00:00:00'}
$assignment=@($ast.FindAll({param($n)$n-is[Management.Automation.Language.AssignmentStatementAst]-and$n.Left-is[Management.Automation.Language.VariableExpressionAst]-and$n.Left.VariablePath.UserPath-ceq'stamp'},$true))
if($assignment.Count-ne1){throw 'missing actual timestamp assignment'}
. ([scriptblock]::Create($assignment[0].Extent.Text))
if($stamp.Offset-ne[timespan]::Zero-or$stamp.ToUnixTimeSeconds()-ne315532800){throw 'release epoch depends on host timezone'}
WriteDeterministicZip $tree (Join-Path $Root 'original.zip') $stamp
WriteDeterministicZip $tree (Join-Path $Root 'second.zip') $stamp
if((Sha (Join-Path $Root 'original.zip'))-cne(Sha (Join-Path $Root 'second.zip'))){throw 'non-deterministic fixture ZIP'}
[IO.File]::AppendAllText((Join-Path $tree 'a.txt'),'changed')
$rejected=$false;try{AssertSourceUnchanged $tree @($entry)}catch{$rejected=$true};if(-not$rejected){throw 'source drift accepted'}
[IO.File]::WriteAllText((Join-Path $tree 'a.txt'),"synthetic-a`n",[Text.UTF8Encoding]::new($false))
$outside=Join-Path $Root 'outside';New-Item -ItemType Directory -Path $outside|Out-Null
[IO.File]::WriteAllText((Join-Path $outside 'data.txt'),'synthetic outside')
New-Item -ItemType Junction -Path (Join-Path $tree 'linked') -Target $outside|Out-Null
$rejected=$false;try{ResolveProjectFile $tree 'linked/data.txt' 'fixture'|Out-Null}catch{$rejected=$true};if(-not$rejected){throw 'parent reparse accepted'}
$rejected=$false;try{WriteDeterministicZip $tree (Join-Path $Root 'reparse.zip') $stamp}catch{$rejected=$true};if(-not$rejected){throw 'directory reparse accepted'}
Write-Output 'BUILDER_HELPERS_PASS source_stability=1 deterministic_zip=1 parent_reparse=1 tree_reparse=1 utc_epoch=1'
''',encoding='utf-8',newline='\n')
    run('builder-helpers',[args.pwsh,'-NoProfile','-File',probe,'-Gate',gate,'-Root',root])
    with zipfile.ZipFile(root/'original.zip')as z:
        assert z.namelist()==['a.txt','z/data.txt'];assert all(x.compress_type==zipfile.ZIP_STORED for x in z.infolist())
    trusted=root/'trusted';trusted.mkdir();policies=[];keys=[]
    for i in [1,2]:
        key=trusted/('fixture-key-'+str(i));run('keygen-'+str(i),[args.ssh_keygen,'-q','-t','ed25519','-f',key,'-N','','-C','TEST_ONLY_NOT_RELEASE_IDENTITY']);keys.append(key)
        policy=trusted/('allowed-'+str(i));parts=Path(str(key)+'.pub').read_text().split();policy.write_text('fixture-signer '+parts[0]+' '+parts[1]+'\n');policies.append(policy)
    counter=0
    def fixture(name,key,policy,zip_path):
        nonlocal counter
        folder=root/name;folder.mkdir();shutil.copyfile(zip_path,folder/'artifact.zip');artifact_hash=sha(folder/'artifact.zip');entries=[]
        with zipfile.ZipFile(zip_path)as z:
            for e in z.infolist():entries.append({'path':e.orig_filename,'sha256':hashlib.sha256(z.read(e)).hexdigest(),'size':e.file_size})
        write(folder/'artifact.manifest.json',{'schema_version':1,'artifact_sha256':artifact_hash,'files':entries})
        write(folder/'source.manifest.json',{'schema_version':1,'files':[],'fixture_only':True})
        write(folder/'osv.json',{'results':[{'packages':[{'package':{'name':'SYNTHETIC-ONLY','version':'0','ecosystem':'npm'}}]}],'fixture_only':True,'note':'Synthetic verifier input, no vulnerability scan performed.'})
        write(folder/'SBOM.spdx.json',{'spdxVersion':'SPDX-2.3','SPDXID':'SPDXRef-DOCUMENT','dataLicense':'CC0-1.0','name':'SYNTHETIC-VERIFIER-FIXTURE','documentNamespace':'https://example.invalid/synthetic-verifier','creationInfo':{'created':'2026-09-13T00:00:00Z','creators':['Tool: local-fixture']},'comment':'Synthetic verifier input, not a real dependency SBOM.','packages':[{'name':'SYNTHETIC-ONLY','versionInfo':'0','externalRefs':[{'referenceType':'purl','referenceLocator':'pkg:npm/SYNTHETIC-ONLY@0'}]}]})
        (folder/'THIRD_PARTY_NOTICES.txt').write_text('Synthetic fixture data authored for verifier tests only. Not a product release.\n');shutil.copyfile(policy,folder/'allowed_signers')
        mapping={'artifact_manifest_sha256':'artifact.manifest.json','source_manifest_sha256':'source.manifest.json','sbom_sha256':'SBOM.spdx.json','sca_sha256':'osv.json','notices_sha256':'THIRD_PARTY_NOTICES.txt','allowed_signers_sha256':'allowed_signers'};evidence={k:sha(folder/v)for k,v in mapping.items()};identity={'product_name':'TEST-ONLY-fixture','product_version':'0','source_revision':'SYNTHETIC','invocation_id':name};internal={**evidence,'artifact_zip_timestamp':'1980-01-01T00:00:00'}
        statement={'_type':'https://in-toto.io/Statement/v1','subject':[{'name':'TEST-ONLY-fixture.zip','digest':{'sha256':artifact_hash}}],'predicateType':'https://slsa.dev/provenance/v1','predicate':{'buildDefinition':{'externalParameters':identity,'internalParameters':internal},'runDetails':{'metadata':{'invocationId':name}}}}
        write(folder/'provenance.intoto.json',statement);counter+=1;run('sign-'+str(counter),[args.ssh_keygen,'-Y','sign','-f',key,'-n','elite-release-v1',folder/'provenance.intoto.json']);evidence.update(provenance_sha256=sha(folder/'provenance.intoto.json'),signature_sha256=sha(folder/'provenance.intoto.json.sig'))
        write(folder/'RELEASE_RECEIPT.json',{'schema_version':1,'gate':'PORTABLE-SIGNED-RELEASE-EVIDENCE-GATE','status':'PASS','reproducibility_runs':2,'product':{'name':identity['product_name'],'version':'0'},'source_revision':'SYNTHETIC','invocation_id':name,'artifact_sha256':artifact_hash,'artifact_zip_timestamp':'1980-01-01T00:00:00','release_identity':'fixture-signer','namespace':'elite-release-v1','tools':{'ssh_keygen':{'sha256':sha(Path(args.ssh_keygen))}},'evidence':evidence,'fixture_only':True})
        return folder
    def verify(label,folder,policy,expected,reason=''):
        result=run(label,[args.pwsh,'-NoProfile','-File',gate/'verify_release.ps1','-ReleaseRoot',folder,'-TrustedAllowedSigners',policy,'-SshKeygen',args.ssh_keygen],expected)
        assert not reason or reason in result,(label,result[-700:])
    good=fixture('valid-fixture',keys[0],policies[0],root/'original.zip');verify('valid-external-policy',good,policies[0],0,'PORTABLE_SIGNED_RELEASE_VERIFY_PASS')
    verify('reject-internal-policy',good,good/'allowed_signers',1,'outside the received release')
    swapped=fixture('self-consistent-other-key',keys[1],policies[1],root/'original.zip');verify('reject-swapped-signer',swapped,policies[0],1,'externally trusted policy')
    receipt=good/'RELEASE_RECEIPT.json';original=receipt.read_bytes();v=json.loads(original);v['namespace']='other-protocol';write(receipt,v);verify('reject-namespace',good,policies[0],1,'unexpected signature namespace');receipt.write_bytes(original)
    for label,path in [('backslash','a\\b.txt'),('drive','C:entry.txt'),('dot-part','a/./b.txt')]:
        badzip=root/(label+'.zip')
        info=zipfile.ZipInfo('fixture.txt',(1980,1,1,0,0,0));info.filename=path;info.orig_filename=path
        with zipfile.ZipFile(badzip,'w',compression=zipfile.ZIP_STORED)as z:z.writestr(info,b'synthetic')
        with zipfile.ZipFile(badzip)as z:assert [e.orig_filename for e in z.infolist()]==[path], 'negative ZIP path was normalized by fixture writer'
        bad=fixture('signed-'+label,keys[0],policies[0],badzip);verify('reject-'+label,bad,policies[0],1,'unsafe or duplicate ZIP entry')
    with(good/'artifact.zip').open('ab')as f:f.write(b'changed')
    verify('reject-artifact-change',good,policies[0],1,'artifact hash mismatch')
    result={'state':'PASS_VERIFIER_FIXTURES_ONLY','builder_helper_cases':12,'zip_separator_case':1,'signed_verifier_cases':8,'real_release_built':False,'real_sca_performed':False,'user_signing_key_used':False,'fixture_private_keys_outside_product':True,'gate_files':{p.name:sha(p)for p in gate.glob('*.ps1')},'logs':{p.name:sha(p)for p in root.glob('*.log')}};write(root/'result.json',result);print(json.dumps(result))

if __name__=='__main__':main()
