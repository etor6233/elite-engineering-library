"""AUTHORED psql orchestration: hash ledger, session lock, atomic migration steps.

Only an owned loopback fixture database is accepted. A pre-existing unregistered
schema requires baseline admission; the runner never guesses it or downgrades.
"""
from pathlib import Path
import argparse,hashlib,json,os,re,subprocess,sys
from build_local_reference import plain,sha,raw_json,require,read_json,relative
from local_release import loopback

def program(root,files):
    names=sorted(files);require(names and all(re.fullmatch(r'[0-9]{4}_[a-z0-9_]+\.up\.sql',x)for x in names),'migration inventory required')
    prefix="""\\set ON_ERROR_STOP on
select pg_advisory_lock(2808,402);
do $guard$ begin
 if to_regclass('library_delivery.migration') is null and to_regclass('platform.tenant') is not null then
  raise exception 'EXISTING_SCHEMA_REQUIRES_BASELINE';
 end if;
end $guard$;
create schema if not exists library_delivery;
create table if not exists library_delivery.migration(name text primary key, ordinal integer unique not null, sha256 text not null, applied_at timestamptz not null default clock_timestamp());
"""
    allowed=','.join("'"+n+"'"for n in names)
    prefix+=f"select count(*)=0 as no_unknown from library_delivery.migration where name not in ({allowed}) \\gset\n\\if :no_unknown\n\\else\n\\quit 3\n\\endif\n"
    for i,name in enumerate(names,1):
        p=relative(root,name);h=files[name];raw=p.read_bytes();require(hashlib.sha256(raw).hexdigest()==h,'migration changed: '+name);body=raw.decode('utf-8')
        begin=re.match(r'\s*(?:(?:--[^\n]*\n)\s*)*begin\s*;',body,re.I);end=re.search(r'commit\s*;\s*$',body,re.I)
        require(bool(begin)==bool(end),'asymmetric transaction wrapper')
        if begin:body=body[begin.end():end.start()]
        prefix+=f"select not exists(select 1 from library_delivery.migration where name='{name}' and (sha256<>'{h}' or ordinal<>{i})) as matches \\gset\n\\if :matches\n\\else\n\\quit 4\n\\endif\n"
        prefix+=f"select exists(select 1 from library_delivery.migration where name='{name}') as applied \\gset\n\\if :applied\n\\else\nbegin;\n{body}\ninsert into library_delivery.migration(name,ordinal,sha256) values('{name}',{i},'{h}');\ncommit;\n\\endif\n"
    prefix+=f"select count(*)={len(names)} as complete from library_delivery.migration \\gset\n\\if :complete\n\\echo LOCAL_MIGRATIONS_PASS\n\\else\n\\quit 5\n\\endif\n"
    return prefix.encode()
def migrate(root,files,database,psql,psql_sha256,output):
    u=loopback(database,'postgres');require(u.username=='postgres'and re.fullmatch(r'/elite_payment_connected_[0-9a-f]{32}',u.path)and u.query=='sslmode=disable','owned fixture database required')
    root=plain(root);output=plain(output);psql=plain(psql)
    require(not output.exists()and output.parent.is_dir(),'absent receipt directory required');require(sha(psql)==psql_sha256,'psql pin changed')
    source=program(root,files);output.mkdir();(output/'apply.sql').write_bytes(source)
    env={'SYSTEMROOT':os.environ['SYSTEMROOT'],'WINDIR':os.environ['SYSTEMROOT'],'PGCONNECT_TIMEOUT':'3','PGCLIENTENCODING':'UTF8'}
    result={'schema':'elite-local-migration/v1','scope':'LOCAL_FIXTURES','state':'FAIL','database':database,'files':files,'psql_sha256':psql_sha256,'program_sha256':sha(output/'apply.sql')}
    try:
        with(output/'apply.log').open('xb')as log:
            q=subprocess.run([psql,'-X','-w','--dbname',database,'--file',output/'apply.sql'],env=env,stdout=log,stderr=subprocess.STDOUT,shell=False,timeout=120,creationflags=subprocess.CREATE_NO_WINDOW)
        result['exit_code']=q.returncode;result['log_sha256']=sha(output/'apply.log')
        require(q.returncode==0 and b'LOCAL_MIGRATIONS_PASS'in(output/'apply.log').read_bytes(),'migration failed; preserve ledger/log and resume after canonical correction')
        result['state']='PASS';return result
    finally:(output/'receipt.json').write_bytes(raw_json(result))
def main():
    p=argparse.ArgumentParser(description=__doc__)
    for name in ['root','inventory','inventory-sha256','database-url','psql','psql-sha256','receipt-directory']:p.add_argument('--'+name,required=True)
    a=p.parse_args()
    try:migrate(a.root,read_json(a.inventory,a.inventory_sha256),a.database_url,a.psql,a.psql_sha256,a.receipt_directory);print('LOCAL_MIGRATIONS_PASS');return 0
    except (ValueError,OSError,subprocess.SubprocessError)as e:print('LOCAL_MIGRATIONS_FAIL '+str(e),file=sys.stderr);return 2
if __name__=='__main__':raise SystemExit(main())
