import copy,json,unittest
from pathlib import Path
from validate_task_mapping import validate
class MappingTests(unittest.TestCase):
 def setUp(self):self.value=json.loads((Path(__file__).parent/'document-task-mapping.json').read_text(encoding='utf8'))
 def test_all_classes_explicit_without_ocr_claim(self):
  got=validate(self.value);self.assertEqual(got['class_count'],21);self.assertFalse(got['ocr21_acceptance'])
 def test_missing_class_rejected(self):
  self.value['classes'].pop()
  with self.assertRaises(ValueError):validate(self.value)
 def test_duplicate_class_rejected(self):
  self.value['classes'][1]=copy.deepcopy(self.value['classes'][0])
  with self.assertRaises(ValueError):validate(self.value)
 def test_false_ocr_pass_rejected(self):
  self.value['classes'][0]['states']['ocr_field_quality']='PASS'
  with self.assertRaises(ValueError):validate(self.value)
 def test_implicit_invoice_ui_for_packing_rejected(self):
  self.value['classes'][4]['current_ui_route']='/experience/documents'
  with self.assertRaises(ValueError):validate(self.value)
 def test_reviewer_edit_without_new_proposal_rejected(self):
  self.value['common_review_contract']['reviewer_mutates_proposal']=True
  with self.assertRaises(ValueError):validate(self.value)
 def test_self_approval_rejected(self):
  self.value['common_review_contract']['same_actor_cannot_approve_own_proposal']=False
  with self.assertRaises(ValueError):validate(self.value)
 def test_inline_untrusted_preview_rejected(self):
  self.value['download_contract']['inline_untrusted_document']=True
  with self.assertRaises(ValueError):validate(self.value)
 def test_storage_authorization_rejected(self):
  self.value['classes'][0]['automatic_storage_authorized']=True
  with self.assertRaises(ValueError):validate(self.value)
 def test_duplicate_task_mapping_rejected(self):
  self.value['classes'][1]['proposed_task_id']=self.value['classes'][0]['proposed_task_id']
  with self.assertRaises(ValueError):validate(self.value)
if __name__=='__main__':unittest.main()
