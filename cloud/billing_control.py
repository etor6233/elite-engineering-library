"""AUTHORED provider-export -> durable journal -> aggregate alert command.

Does not send notifications or enforce a provider spending cap. The output is
consumed by the existing operations/alerting owner after target integration.
"""
from pathlib import Path
import argparse
from cloud_control import BillingJournal,normalize_billing_rows,raw,read,require

def main():
    p=argparse.ArgumentParser();p.add_argument("--journal",required=True);p.add_argument("--export",required=True);p.add_argument("--identity",required=True);p.add_argument("--policy",required=True);p.add_argument("--observed-at",required=True);p.add_argument("--output",required=True);a=p.parse_args()
    out=Path(a.output);require(not out.exists(),"fresh budget observation receipt required")
    normalized=normalize_billing_rows(read(a.export),read(a.identity));journal=BillingJournal(a.journal)
    try:
        effect=journal.import_snapshot(normalized);alert=journal.budget(read(a.policy),a.observed_at)
        result={"schema":"elite-billing-control-receipt/v403.1","import_result":effect,"alert":alert,"notification_sent":False,"provider_restriction":"NOT_EXECUTED","hard_cap":False}
        out.parent.mkdir(parents=True,exist_ok=True);out.write_bytes(raw(result));print(raw(result).decode(),end="")
    finally:journal.close()

if __name__=="__main__":main()
