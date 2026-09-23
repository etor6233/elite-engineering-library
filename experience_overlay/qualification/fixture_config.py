"""AUTHORED fixture configuration. No credentials, accounts or business policies inferred."""
from pathlib import Path
import json,sys,os
r=Path(__file__).resolve().parent;w=Path(os.environ['UI_TARGET'])/'.next/standalone/config';e=Path(os.environ['UI_EVIDENCE'])
defaults={'experience-checklists.json':{'schema_version':'1','tenants':{}},'experience-options.json':{'schema_version':'1','tenants':{}}}
options={k:[] for k in ['variants','priceBooks','assignees','resources','availabilityReasons','appointmentReasons','evidence']}
for group,value,label in [('variants','variant-v403','Modelo de referencia'),('priceBooks','retail-v403','Venta minorista'),('assignees','operator','Operador de referencia'),('resources','resource-v403','Equipo de entrega'),('availabilityReasons','maintenance','Mantenimiento programado'),('appointmentReasons','customer-request','Solicitud del cliente'),('evidence','a'*64,'Inspección técnica verificada')]:
 options[group].append({'value':value,'label':label,'allowed_subjects':['operator']})
for group in json.loads((r/'additional-guided-controls.json').read_text(encoding='utf8'))['groups']:
 options[group]=[{'value':group+'-v403','label':'Referencia aprobada de '+group,'allowed_subjects':['operator']}]
options['variants'].append({'value':'private','label':'No autorizado','allowed_subjects':['another-subject']})
configured={'experience-checklists.json':{'schema_version':'1','tenants':{'tenant':{'organizations':{'store':[{'id':'standard-delivery','version':1}]}}}},'experience-options.json':{'schema_version':'1','tenants':{'tenant':{'organization_labels':{'store':'Sucursal Centro'},'organizations':{'store':options}}}}}
for name,data in (defaults if '--restore' in sys.argv else configured).items():
 (w/name).write_text(json.dumps(data,indent=2)+'\n',encoding='utf8')
 if '--restore' not in sys.argv:(e/('runtime-'+name)).write_text(json.dumps(data,indent=2)+'\n',encoding='utf8')
print('Default configuration restored' if '--restore' in sys.argv else 'Synthetic runtime configuration applied')

# Dedicated non-production config; the qualification launcher removes/restores it.
fixture_business=json.loads((w/"business.example.json").read_text(encoding="utf-8-sig"))
fixture_business["features"]["supply_portal"]=True
(w/"experience.browser.fixture.json").write_text(json.dumps(fixture_business,indent=2)+"\n",encoding="utf8")
