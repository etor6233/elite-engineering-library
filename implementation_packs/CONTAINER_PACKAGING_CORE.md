# Container Packaging Core

## 1. Metadata

```yaml
pack_id: "CONTAINER-PACKAGING-CORE"
pack_version: "0.1.2"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa packaging OCI multi-stage, runtime no-root y una topología Compose local endurecida con validación ejecutable."
stacks: ["Dockerfile 1.18", "Compose Specification", "Python 3.14 stdlib"]
compatible_with: ["GO-ELECTROMOBILITY-APPLICATION >=1.22.2 <2.0.0", "PG-TX-FOUNDATION 0.1.x", "ELECTROMOBILITY-FRANCHISE-MODULES 0.1.x"]
incompatible_with: []
license_expression: "LicenseRef-Workspace-Owner AND Apache-2.0/MIT image dependencies"
upstream_sources: ["https://docs.docker.com/build/", "https://docs.docker.com/compose/compose-file/", "https://hub.docker.com/_/golang", "https://hub.docker.com/_/debian"]
verified_at: "2026-08-25"
```

Image defaults are intentionally absent. Container execution is outside the Windows local reference admission; choose and admit exact image digests, target and runtime before using this template. Docker/Podman execution remains unqualified. The source context now includes all selected local Go modules.

## 2. Applicability

Use this pack when a composed Go/PostgreSQL service needs reproducible OCI packaging and a disposable local topology. Reject it for production until image digests, the target secret mechanism, TLS and the selected runtime have been admitted. It does not choose a cloud, orchestrator or registry.

## 3. Architecture contract

The build is multi-stage; the runtime is non-root and contains only the service artifact and required trust/runtime files. Database migration is a separate, fail-fast step. Compose is a local integration topology, not a production control plane. Configuration enters through environment references; literal credentials are rejected. A failed migration prevents traffic, and rollback reuses a previously admitted immutable image rather than rebuilding it.

## 4. Exact file manifest

```text
CREATE .dockerignore
CREATE Dockerfile.api
CREATE deploy/compose.yaml
CREATE deploy/postgres-migrate.sh
CREATE packaging/validate_packaging.py
CREATE packaging/test_validate_packaging.py
```

## 5. Materialization blocks

### FILE: `.dockerignore`

```yaml
block_id: "CONTAINER-PACKAGING:dockerignore:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "a19f31cc3915069c5298cfe8d517432adcadaba362a0c526aee00d8982ce0712"
variables: []
secrets_allowed: false
```

````gitignore
.git
.env
.env.*
!.env.example
**/.DS_Store
**/node_modules
**/.next
**/dist
**/coverage
**/*.log
**/*.dump
**/*.pem
**/*.key
**/*secret*
MATERIALIZATION_RECORD.md
reconstruction_evidence
````

### FILE: `Dockerfile.api`

```yaml
block_id: "CONTAINER-PACKAGING:api-image:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "0c6db00f02eb69dcf8f3fd69f23de125ab4654034c60765c2137d44aeb001ce6"
variables: []
secrets_allowed: false
```

````dockerfile
# Container execution is outside the Windows local reference admission.
# Select and admit immutable image digests before using this template.
ARG GO_IMAGE
ARG RUNTIME_IMAGE

FROM ${GO_IMAGE} AS build
WORKDIR /src
COPY . ./
RUN --mount=type=cache,target=/go/pkg/mod,sharing=locked go mod download
ARG TARGETOS=linux
ARG TARGETARCH=amd64
RUN --mount=type=cache,target=/go/pkg/mod,sharing=locked --mount=type=cache,target=/root/.cache/go-build,sharing=locked CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -trimpath -buildvcs=false -ldflags="-s -w -buildid=" -o /out/electromobility-api ./cmd/electromobility-api

FROM ${RUNTIME_IMAGE} AS runtime
WORKDIR /
COPY --from=build --chown=nonroot:nonroot /out/electromobility-api /electromobility-api
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/electromobility-api"]
````

### FILE: `deploy/compose.yaml`

```yaml
block_id: "CONTAINER-PACKAGING:compose:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "610c8641399ee0e5c790e41a75e3a0513bf92281314fe18acd34c1c27e4bd037"
variables: []
secrets_allowed: false
```

````yaml
name: elite-enterprise-local
services:
  postgres:
    image: postgres:18.6-trixie
    environment:
      POSTGRES_DB: elite
      POSTGRES_USER: elite
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD:?set POSTGRES_PASSWORD outside source control}
    volumes:
      - postgres-data:/var/lib/postgresql
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U elite -d elite"]
      interval: 5s
      timeout: 3s
      retries: 20
    restart: unless-stopped
    security_opt: ["no-new-privileges:true"]
    cap_drop: ["ALL"]
    cap_add: ["CHOWN", "DAC_OVERRIDE", "FOWNER", "SETGID", "SETUID"]
    networks: [data]

  migrate:
    image: postgres:18.6-trixie
    depends_on:
      postgres: {condition: service_healthy}
    environment:
      PGHOST: postgres
      PGPORT: "5432"
      PGDATABASE: elite
      PGUSER: elite
      PGPASSWORD: ${POSTGRES_PASSWORD:?set POSTGRES_PASSWORD outside source control}
    volumes:
      - ../db/migrations:/migrations:ro
      - ./postgres-migrate.sh:/postgres-migrate.sh:ro
    entrypoint: ["/bin/sh", "/postgres-migrate.sh"]
    restart: "no"
    security_opt: ["no-new-privileges:true"]
    cap_drop: ["ALL"]
    networks: [data]

  api:
    build:
      context: ..
      dockerfile: Dockerfile.api
    depends_on:
      migrate: {condition: service_completed_successfully}
    environment:
      DATABASE_URL: postgresql://elite:${POSTGRES_PASSWORD:?set POSTGRES_PASSWORD outside source control}@postgres:5432/elite?sslmode=disable
      OIDC_ISSUER: ${OIDC_ISSUER:?set OIDC_ISSUER}
      OIDC_AUDIENCE: ${OIDC_AUDIENCE:?set OIDC_AUDIENCE}
      HTTP_ADDRESS: :8080
    ports: ["127.0.0.1:${API_PORT:-8080}:8080"]
    read_only: true
    tmpfs: ["/tmp:rw,noexec,nosuid,size=16m"]
    init: true
    restart: unless-stopped
    security_opt: ["no-new-privileges:true"]
    cap_drop: ["ALL"]
    networks: [edge, data]

volumes:
  postgres-data: {}
networks:
  edge: {}
  data:
    internal: true
````

### FILE: `deploy/postgres-migrate.sh`

```yaml
block_id: "CONTAINER-PACKAGING:migrate:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "08d868d4aec724135ebacade7fb5bd7d956a34da1811f76f0f16550237959c0e"
variables: []
secrets_allowed: false
```

````sh
#!/bin/sh
set -eu
for migration in /migrations/*.up.sql; do
  test -f "$migration"
  psql -X --set=ON_ERROR_STOP=1 --single-transaction --file="$migration"
done
````

### FILE: `packaging/validate_packaging.py`

```yaml
block_id: "CONTAINER-PACKAGING:validator:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "2eac9384cd79ee1ffe093049de7de1d52a96d48a61bdaad8dd0be13c6833ad3b"
variables: []
secrets_allowed: false
```

````python
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
````

### FILE: `packaging/test_validate_packaging.py`

```yaml
block_id: "CONTAINER-PACKAGING:validator-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "6e62baf2c56066f1423fb3a53f46bd813cc1b00307d27d472e1d07d2bc635343"
variables: []
secrets_allowed: false
```

````python
from __future__ import annotations
import pathlib,tempfile,unittest
from validate_packaging import validate
ROOT=pathlib.Path(__file__).resolve().parents[1]
class PackagingTests(unittest.TestCase):
    def test_canonical_packaging(self)->None: validate(ROOT)
    def test_root_runtime_is_rejected(self)->None:
        with tempfile.TemporaryDirectory() as folder:
            target=pathlib.Path(folder)
            (target/"deploy").mkdir();(target/"packaging").mkdir()
            (target/"Dockerfile.api").write_text((ROOT/"Dockerfile.api").read_text().replace("USER nonroot:nonroot","USER 0"),encoding="utf-8")
            for name in ("compose.yaml","postgres-migrate.sh"):(target/"deploy"/name).write_text((ROOT/"deploy"/name).read_text(),encoding="utf-8")
            with self.assertRaisesRegex(ValueError,"non-root"):validate(target)
    def test_literal_secret_is_rejected(self)->None:
        with tempfile.TemporaryDirectory() as folder:
            target=pathlib.Path(folder);(target/"deploy").mkdir()
            (target/"Dockerfile.api").write_text((ROOT/"Dockerfile.api").read_text(),encoding="utf-8")
            (target/"deploy/postgres-migrate.sh").write_text((ROOT/"deploy/postgres-migrate.sh").read_text(),encoding="utf-8")
            compose=(ROOT/"deploy/compose.yaml").read_text().replace("POSTGRES_PASSWORD: ${POSTGRES_PASSWORD:?set POSTGRES_PASSWORD outside source control}","POSTGRES_PASSWORD: hardcoded")
            (target/"deploy/compose.yaml").write_text(compose,encoding="utf-8")
            with self.assertRaisesRegex(ValueError,"literal database secret"):validate(target)
if __name__=="__main__":unittest.main()
````

## 6. Configuration surface

| Variable/input | Type | Safe default | Secret | Validation/effect |
|---|---|---|---|---|
| `POSTGRES_PASSWORD` | external environment value | none | yes | required by local Compose; never embedded |
| image references | OCI reference/digest | versioned development tags | no | production requires reviewed immutable digests |
| API/OIDC/database variables | environment references | none | mixed | inherited from the composed application and validated before serve |

## 7. Dependency bill

| Package/image/tool | Pin | Use | License | Scope | Official source |
|---|---|---|---|---|---|
| Dockerfile frontend | `1.18` | build grammar | Apache-2.0 | build | `docs.docker.com` |
| Go and Debian runtime images | version family declared in files | build/runtime | image notices apply | build/runtime | Docker Official Images |
| PostgreSQL image | version family declared in Compose | local database | PostgreSQL | runtime | Docker Official Image |
| Python | `3.14+` stdlib | validator/tests | PSF-2.0 | verification | `python.org` |

Production admission replaces every movable image reference with a reviewed digest and retains its SBOM/notices.

## 8. Apply order

Materialize after the Go application and SQL packs, review collisions, set external configuration, validate the pack, build the image, run migrations, start the API and execute smoke/negative tests. In an existing workspace, reconcile Dockerfile and Compose ownership explicitly. Rollback stops the new workload and redeploys the prior digest; schema rollback must obey expand/migrate/contract safety.

## 9. Verification

Run `python packaging/validate_packaging.py .` and `python packaging/test_validate_packaging.py`. On a container-enabled host, additionally resolve builder/runtime/database tags to reviewed digests, build with SBOM/provenance, scan, run Compose with a disposable OIDC issuer, apply migrations, call liveness and a protected negative route, restart the API, verify persistence, then tear down. Production must replace local password interpolation and plaintext database transport with the selected secret manager, workload identity and TLS policy.

## 10. Reconstruction evidence

Clean reconstruction, hashes and static negative gates are recorded in `reconstruction_evidence/CONTAINER_PACKAGING_CORE_2026-08-24_V1.md`; the final library audit rechecks the current pack. Container execution remains an explicit conditioned gate on a host with Docker or Podman.

## 11. Primary references

- Docker multi-stage build documentation: https://docs.docker.com/build/building/multi-stage/
- Go Docker Official Image: https://hub.docker.com/_/golang
- PostgreSQL Docker Official Image: https://hub.docker.com/_/postgres
- Distroless repository and support policy: https://github.com/GoogleContainerTools/distroless

V402 composed delta: Local delivery316: native stop owner reused, Next standalone exact build identity, container template without mutable defaults and complete local module context. Local fixture qualification only; docs/LOCAL_REFERENCE_DELIVERY.md.
