"""AUTHORED routing contracts. Synthetic records do not prove OCR/access/corpus."""
import os,sys,json,unittest,tempfile
from pathlib import Path
sys.path.insert(0,os.environ.get('ROUTING_OWNER',str(Path(__file__).resolve().parent)))
from validate_document_routing import validate_configuration,write_receipt,PROVIDER_PLANS
from test_validate_document_routing import complete_configuration
PLAN='markdown_system/PADDLEOCR_LOCAL_RUNTIME_PACK_PLAN.md'
class PaddleRoutingContracts(unittest.TestCase):
 def setUp(self):
  self.tmp=tempfile.TemporaryDirectory();self.root=Path(self.tmp.name);self.path=self.root/'routing.json'
 def tearDown(self):self.tmp.cleanup()
 def fixture(self):
  value=complete_configuration();active=next(x for x in value['classes'] if x['decision']=='REQUIRED')
  active['candidate_lanes']=[{'provider':'PADDLEOCR_LOCAL','role':'PRIMARY','pack_plan':PLAN,'provider_profile_path':'config/paddle-ocr.json'}]
  active['access']={'probe_status':'PROVEN','identity_reference':'local-runtime://synthetic-paddle-fixture','region':'local','evidence_path':'evidence/synthetic-runtime-access.json'}
  return value,active
 def put(self,value):self.path.write_text(json.dumps(value)+'\n',encoding='utf8')
 def test_local_primary_emits_exact_plan_and_no_storage_authority(self):
  value,_=self.fixture();self.put(value);receipt=write_receipt(self.path,self.root/'receipt')
  self.assertEqual(receipt['selected_pack_plans'],[PLAN]);self.assertFalse(receipt['automatic_storage_authorized'])
  self.assertEqual(receipt['selected_routes'][0]['provider'],'PADDLEOCR_LOCAL')
  self.assertEqual(receipt['selected_routes'][0]['provider_profile_path'],'config/paddle-ocr.json')
 def test_local_evaluation_lane_preserves_primary_and_role(self):
  value,active=self.fixture();active['candidate_lanes'][0]['role']='EVALUATION'
  active['candidate_lanes'].append({'provider':'AWS','role':'PRIMARY','pack_plan':PROVIDER_PLANS['AWS'],'provider_profile_path':'config/aws.json'})
  self.put(value);receipt=write_receipt(self.path,self.root/'receipt')
  self.assertEqual(set(receipt['selected_pack_plans']),{PLAN,PROVIDER_PLANS['AWS']})
  self.assertEqual({(x['provider'],x['role']) for x in receipt['selected_routes']},{('AWS','PRIMARY'),('PADDLEOCR_LOCAL','EVALUATION')})
 def test_invented_provider_rejected(self):
  value,active=self.fixture();active['candidate_lanes'][0]['provider']='PADDLE_CLOUD';self.put(value)
  with self.assertRaisesRegex(ValueError,'provider is unsupported'):validate_configuration(self.path)
 def test_provider_cannot_select_different_plan(self):
  value,active=self.fixture();active['candidate_lanes'][0]['pack_plan']=PROVIDER_PLANS['AWS'];self.put(value)
  with self.assertRaisesRegex(ValueError,'does not match provider'):validate_configuration(self.path)
 def test_profile_traversal_rejected(self):
  value,active=self.fixture();active['candidate_lanes'][0]['provider_profile_path']='../unauthorized.json';self.put(value)
  with self.assertRaisesRegex(ValueError,'canonical relative path'):validate_configuration(self.path)
 def test_duplicate_local_provider_rejected(self):
  value,active=self.fixture();active['candidate_lanes'].append(active['candidate_lanes'][0].copy());self.put(value)
  with self.assertRaisesRegex(ValueError,'duplicate providers'):validate_configuration(self.path)
 def test_local_lane_does_not_bypass_access_probe(self):
  value,active=self.fixture();active['access']['probe_status']='NOT_RUN';self.put(value)
  with self.assertRaisesRegex(ValueError,'probe_status must be PROVEN'):validate_configuration(self.path)
 def test_local_lane_does_not_bypass_corpus(self):
  value,active=self.fixture();active['corpus']['authorized']=False;self.put(value)
  with self.assertRaisesRegex(ValueError,'corpus must be explicitly authorized'):validate_configuration(self.path)
 def test_local_lane_cannot_authorize_storage(self):
  value,active=self.fixture();active['automatic_storage']=True;self.put(value)
  with self.assertRaisesRegex(ValueError,'automatic_storage must be false'):validate_configuration(self.path)
 def test_no_output_commit_on_invalid_local_route(self):
  value,active=self.fixture();active['candidate_lanes'][0]['pack_plan']=PROVIDER_PLANS['GOOGLE'];self.put(value);out=self.root/'receipt'
  with self.assertRaises(ValueError):write_receipt(self.path,out)
  self.assertFalse(out.exists());self.assertEqual(list(self.root.glob('.document-routing-stage-*')),[])
if __name__=='__main__':unittest.main()
