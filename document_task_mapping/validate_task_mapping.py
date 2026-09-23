"""Validate visibility and separation of states; never authorize OCR/storage."""
import json
from pathlib import Path
BASELINE_IDS=frozenset('supplier-invoice receipt proforma-invoice commercial-invoice packing-list purchase-order purchase-order-confirmation bill-of-lading-airway-bill delivery-note certificate-of-origin quality-inspection-certificate customs-declaration-clearance insurance-freight-document price-list-quotation product-specification-catalog contract-addendum bank-payment-statement warranty-service-claim identity-regulatory-document multi-document-package other-named-class'.split())
def validate(value):
 if value.get('schema')!='elite-document-task-mapping/v1':raise ValueError('schema mismatch')
 rows=value.get('classes',[]);ids=[x['class_id'] for x in rows]
 if len(ids)!=21 or set(ids)!=BASELINE_IDS:raise ValueError('21 unique baseline classes required')
 tasks=[x.get('proposed_task_id') for x in rows]
 if any(not isinstance(x,str) or not x.strip() for x in tasks) or len(set(tasks))!=21:raise ValueError('unique proposed task mapping required')
 review=value['common_review_contract']
 if review.get('reviewer_mutates_proposal') is not False or review.get('same_actor_cannot_approve_own_proposal') is not True or review.get('changed_fields_require_new_proposal') is not True or review.get('decision_binds')!='payload_sha256' or review.get('reason_required') is not True:raise ValueError('review separation invariant')
 if value.get('download_contract')!={'disposition':'attachment','content_type':'application/octet-stream','cache_control':'no-store','x_content_type_options':'nosniff','inline_untrusted_document':False}:raise ValueError('safe original download contract required')
 for row in rows:
  if row.get('automatic_storage_authorized') is not False or row['states'].get('production_authorized') is not False:raise ValueError('planning contract cannot authorize effects')
  if row.get('task_mapping_status')!='PROPOSED_NOT_IMPLEMENTED' or row.get('business_scope_decision')!='REQUIRES_PROJECT_BLUEPRINT':raise ValueError('proposed tasks are not implemented or business scope decisions')
  if row['states'].get('runtime_this_revision')!='NOT_RUN' or row['states'].get('ocr_field_quality')!='NOT_EVALUATED':raise ValueError('no runtime or OCR acceptance in inventory')
  if row['class_id']!='supplier-invoice' and (row.get('current_ui_route') is not None or row.get('implementation_owner') is not None):raise ValueError('unimplemented class cannot inherit invoice owner or UI')
 return {'schema':'elite-document-task-map-validation/v1','class_count':21,'contract_validation':'PASS','runtime_executed':False,'ocr21_acceptance':False,'storage_authorized':False,'production_authorized':False}
if __name__=='__main__':
 import argparse
 p=argparse.ArgumentParser();p.add_argument('mapping',type=Path);a=p.parse_args();print(json.dumps(validate(json.loads(a.mapping.read_text(encoding='utf8'))),sort_keys=True))
