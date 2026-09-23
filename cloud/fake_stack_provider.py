"""AUTHORED fixture provider. Stores argv as synthetic state; no network."""
from pathlib import Path
import json,sys
p=Path(sys.argv[1]);argv=sys.argv[2:]
if not argv or argv[0]!="gcloud":sys.exit(9)
state=json.loads(p.read_text()) if p.exists() else []
state.append(argv);p.write_text(json.dumps(state));print(json.dumps({"method":"SIMULATED_PROVIDER","effects":len(state)}))
