"""Synthetic provider used only by isolated behavioral tests. No network."""
from __future__ import annotations
import json
from pathlib import Path
import sys

def main():
    state=Path(sys.argv[1]);args=sys.argv[2:]
    if args[:3]==["gcloud","run","deploy"]:
        service=args[3]; suffix=args[args.index("--revision-suffix")+1]; image=args[args.index("--image")+1]
        old=json.loads(state.read_text()) if state.exists() else {"revisions":{},"traffic":{}}
        old["revisions"][service+"-"+suffix]=image
        state.write_text(json.dumps(old));print(json.dumps({"created":service+"-"+suffix,"traffic":old["traffic"]}));return 0
    if args[:4]==["gcloud","run","services","update-traffic"]:
        old=json.loads(state.read_text());revision,percent=args[args.index("--to-revisions")+1].split("=")
        if revision not in old["revisions"]:return 7
        old["traffic"]={revision:int(percent)};state.write_text(json.dumps(old));print(json.dumps(old));return 0
    if args[:4]==["gcloud","run","services","describe"]:
        print(state.read_text());return 0
    if args[:4]==["gcloud","run","revisions","describe"]:
        old=json.loads(state.read_text());rev=args[4];print(json.dumps({"revision":rev,"image":old["revisions"][rev]}));return 0
    return 9

if __name__=="__main__":sys.exit(main())
