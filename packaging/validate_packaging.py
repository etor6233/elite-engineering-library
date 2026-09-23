from __future__ import annotations
import pathlib,re,sys

def require(text:str,pattern:str,label:str)->None:
    if not re.search(pattern,text,re.MULTILINE): raise ValueError(f"missing {label}")
def reject(text:str,pattern:str,label:str)->None:
    if re.search(pattern,text,re.MULTILINE|re.IGNORECASE): raise ValueError(f"forbidden {label}")
def validate(root:pathlib.Path)->None:
    docker=(root/"Dockerfile.api").read_text(encoding="utf-8")
    compose=(root/"deploy/compose.yaml").read_text(encoding="utf-8")
    migration=(root/"deploy/postgres-migrate.sh").read_text(encoding="utf-8")
    require(docker,r"^FROM \$\{GO_IMAGE\} AS build$","named build stage")
    require(docker,r"^FROM \$\{RUNTIME_IMAGE\} AS runtime$","separate runtime stage")
    require(docker,r"CGO_ENABLED=0","static Go build")
    require(docker,r"-trimpath","trimmed source paths")
    require(docker,r'^USER nonroot:nonroot$',"non-root runtime")
    require(docker,r'^ENTRYPOINT \["/electromobility-api"\]$',"exec entrypoint")
    reject(docker,r"FROM .*:(latest|main|master)(\s|$)","floating image tag")
    reject(docker,r"\b(curl|wget|apt-get|apk add)\b","runtime package download")
    for label in ("read_only: true",'cap_drop: ["ALL"]','no-new-privileges:true',"internal: true","service_completed_successfully"):
        require(compose,re.escape(label),label)
    require(compose,r"127\.0\.0\.1:\$\{API_PORT:-8080\}:8080","loopback-only local publication")
    reject(compose,r"(?m)^\s*POSTGRES_PASSWORD:(?![ \t]*\$)[^\r\n]+","literal database secret")
    require(compose,r"DATABASE_URL:.*\$\{POSTGRES_PASSWORD:\?","database password interpolation")
    require(migration,r"ON_ERROR_STOP=1","fail-closed SQL")
    require(migration,r"--single-transaction","transactional migration invocation")
def main()->int:
    try: validate(pathlib.Path(sys.argv[1] if len(sys.argv)>1 else ".").resolve())
    except (OSError,ValueError) as error: print(f"PACKAGING_INVALID: {error}",file=sys.stderr);return 1
    print("PACKAGING_VALID");return 0
if __name__=="__main__": raise SystemExit(main())
