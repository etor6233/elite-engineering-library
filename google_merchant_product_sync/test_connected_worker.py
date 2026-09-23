from pathlib import Path
import json,subprocess,threading,http.server,urllib.parse,hashlib,os,sys
import argparse
parser=argparse.ArgumentParser();parser.add_argument("--python",type=Path,required=True);parser.add_argument("--receipt",type=Path,required=True)
args=parser.parse_args();require_receipt=args.receipt
if require_receipt.exists():raise ValueError("absent receipt required")
sha=lambda p:hashlib.sha256(p.read_bytes()).hexdigest()
root=Path(__file__).resolve().parent;python=args.python;worker=root/'connected_worker.py';owner=root/'sync_product.py';lock=root/'connected-runtime.lock.json'
state={'posts':0,'gets':0,'item':None,'mode':'normal','requests':[]}
class Handler(http.server.BaseHTTPRequestHandler):
 def log_message(self,*args):pass
 def respond(self,status,data):
  raw=json.dumps(data).encode();self.send_response(status);self.send_header('Content-Type','application/json');self.send_header('Content-Length',str(len(raw)));self.end_headers();self.wfile.write(raw)
 def do_POST(self):
  state['posts']+=1
  u=urllib.parse.urlsplit(self.path);q=urllib.parse.parse_qs(u.query)
  assert u.path=='/products/v1/accounts/123456/productInputs:insert'and q.get('dataSource')==['accounts/123456/dataSources/789012'],self.path
  assert self.headers.get('Authorization')=='Bearer synthetic-merchant-sdk-fixture'
  raw=self.rfile.read(int(self.headers['Content-Length']));product=json.loads(raw);state['requests'].append(product)
  assert product['offerId']=='reference-bicycle'and product['versionNumber']=='3'
  assert product['productAttributes']['price']['amountMicros']=='12345678901230000'
  if state['mode']=='unauthorized':return self.respond(401,{'error':{'code':401,'message':'synthetic authorization rejection','status':'UNAUTHENTICATED'}})
  if state['mode']=='redirect':
   self.send_response(307);self.send_header('Location',f'http://127.0.0.1:{self.server.server_port}/must-not-follow');self.end_headers();return
  if state['mode']=='oversize':return self.respond(200,{'unknown':'x'*32769})
  if state['mode']=='reject':return self.respond(400,{'error':{'code':400,'message':'synthetic rejection','status':'INVALID_ARGUMENT'}})
  item=dict(product,name='accounts/123456/productInputs/es~AR~reference-bicycle',product='accounts/123456/products/es~AR~reference-bicycle')
  state['item']=item
  if state['mode']=='lost':
   self.close_connection=True;self.connection.close();return
  return self.respond(200,item)
 def do_GET(self):
  state['gets']+=1
  u=urllib.parse.urlsplit(self.path)
  assert u.path=='/products/v1/accounts/123456/products/es~AR~reference-bicycle' and urllib.parse.parse_qs(u.query)=={'$alt':['json;enum-encoding=int']},self.path
  if state['item']is None:return self.respond(404,{'error':{'code':404,'message':'pending','status':'NOT_FOUND'}})
  product=dict(state['item']);product['name']=product.pop('product');product['dataSource']='accounts/123456/dataSources/789012'
  product['productStatus']={'destinationStatuses':[{'reportingContext':'SHOPPING_ADS','pendingCountries':['AR']}],'lastUpdateDate':'2026-09-13T00:00:00Z'}
  return self.respond(200,product)
server=http.server.ThreadingHTTPServer(('127.0.0.1',0),Handler);thread=threading.Thread(target=server.serve_forever,daemon=True);thread.start()
intent={'account_id':'123456','data_source':'accounts/123456/dataSources/789012','generation':'3','resource':'accounts/123456/products/es~AR~reference-bicycle','product':{'offer_id':'reference-bicycle','content_language':'es','feed_label':'AR','title':'Bicicleta de referencia','description':'Descripción sintética revisada del catálogo.','link':'https://catalog.example.invalid/catalog/bicycle','image_link':'https://catalog.example.invalid/api/public/catalog/media/'+'a'*64,'availability':'OUT_OF_STOCK','condition':'NEW','price':{'amount_micros':'12345678901230000','currency_code':'ARS'},'brand':'','gtins':[]}}
env={k:v for k,v in os.environ.items()if k in ['SYSTEMROOT','WINDIR','TEMP','TMP']};receipts=[]
def call(operation,value=None):
 command={'operation':operation,'intent':value or intent,'mode':'LOCAL_FIXTURES','fixture_origin':f'http://127.0.0.1:{server.server_port}'}
 q=subprocess.run([python,'-I','-S','-B',worker,sha(owner),sha(lock)],input=json.dumps(command).encode(),capture_output=True,env=env,timeout=25,creationflags=subprocess.CREATE_NO_WINDOW)
 assert q.returncode==0 and not q.stderr,q.stderr.decode()
 out=json.loads(q.stdout);receipts.append({'operation':operation,'state':out.get('state'),'result':out})
 return out
try:
 assert call('GET')['state']=='NOT_FOUND'
 out=call('INSERT');assert out['state']=='OBSERVED',out
 assert out['response']['version_number']=='3'
 assert out['response']['product_attributes']['price']['amount_micros']=='12345678901230000'
 out=call('GET');assert out['state']=='OBSERVED'and out['response']['data_source']==intent['data_source'],out
 assert out['response']['product_status']['destination_statuses'][0]['pending_countries']==['AR']
 state['mode']='lost';assert call('INSERT')['state']=='UNKNOWN'
 assert state['posts']==2
 assert call('GET')['state']=='OBSERVED'and state['posts']==2
 state['mode']='reject';assert call('INSERT')['state']=='REJECTED'
 tampered=json.loads(json.dumps(intent));tampered['resource']='accounts/999999/products/es~AR~reference-bicycle'
 assert call('INSERT',tampered)['state']=='UNKNOWN'and state['posts']==3
 for mode,want in [('unauthorized','REJECTED'),('redirect','UNKNOWN'),('oversize','UNKNOWN')]:
  before=state['posts'];state['mode']=mode;out=call('INSERT')
  assert out['state']==want and state['posts']==before+1,(mode,out,state['posts'])
 result={'state':'PASS','real_sdk':'1.8.0','posts':state['posts'],'gets':state['gets'],'requests':state['requests'],'receipts':receipts,'sources':{str(p):sha(p)for p in [worker,owner,lock]},'live_effects':False}
except BaseException as e:
 result={'state':'FAIL','failure':str(e),'type':type(e).__name__,'posts':state['posts'],'gets':state['gets'],'receipts':receipts};raise
finally:
 server.shutdown();server.server_close();thread.join(3)
 require_receipt.write_text(json.dumps(result,indent=2)+'\n',encoding='utf-8',newline='\n')
 print(json.dumps({k:result[k]for k in ['state','posts','gets']}))
