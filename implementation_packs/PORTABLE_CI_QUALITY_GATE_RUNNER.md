# Portable CI Quality Gate Runner

## 1. Metadata

```yaml
pack_id: "PORTABLE-CI-GATE-RUNNER"
pack_version: "0.1.9"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un runner CI/CD stack-neutral, sin shell, con timeouts, redacción y evidencia hash-linked."
stacks: ["Python 3.14+ stdlib"]
compatible_with: ["SECURE-OPS-DELIVERY-CORE >=1.1.4 <2.0.0"]
incompatible_with: []
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources:
  - "https://csrc.nist.gov/pubs/sp/800/218/final"
  - "https://slsa.dev/spec/v1.2/"
verified_at: "2026-09-13"
```

Los bloques de coordinación son `AUTHORED`; V402322 añade16textos legales `VERBATIM` y3avisos `ADAPTED`, declarados individualmente. Este runner coordina herramientas elegidas por el proyecto; no sustituye compiladores, scanners, un build service aislado ni la verificación de provenance.

## 2. Applicability

Use as a stack-neutral orchestrator for already selected project tools when CI needs bounded, redactable, hash-linked evidence. Reject it as a sandbox or scanner: the target CI runner must provide isolation, resource limits, trusted tool pins and secret injection.

## 3. Architecture contract

- ejecuta `argv` directamente con `shell=False`;
- exige IDs únicos, categorías reconocidas, timeout acotado y directorios dentro del workspace;
- categorías requeridas se declaran en el plan y ninguna puede quedar sin gate;
- cada salida se limita y redacta antes de persistirse;
- fail-fast por defecto; un gate obligatorio fallido produce exit 2;
- evidencia JSON incluye hash del plan, tiempos, código de salida y resultado por gate;
- secretos se inyectan externamente; el plan no contiene `env`, passwords ni tokens.

## 4. Exact file manifest

```text
CREATE ci/arca-build-inputs.template.json
CREATE ci/build_arca_reference.py
CREATE ci/build_signed_reference.ps1
CREATE ci/publish_reference_build.py
CREATE ci/test_signed_reference.py
CREATE ci/local-build-inputs.template.json
CREATE docs/SIGNED_LOCAL_REFERENCE.md
CREATE licenses/reference/next-MIT.txt
CREATE licenses/reference/constants-browserify-README.md.txt
CREATE licenses/reference/data-uri-to-buffer-README.md.txt
CREATE licenses/reference/http-proxy-agent-README.md.txt
CREATE licenses/reference/https-proxy-agent-README.md.txt
CREATE licenses/reference/punycode-LICENSE-MIT.txt.txt
CREATE licenses/reference/setimmediate-LICENSE.txt.txt
CREATE licenses/reference/string-hash-README.md.txt
CREATE licenses/reference/unistore-README.md.txt
CREATE licenses/reference/edge-runtime-cookies-LICENSE.md.txt
CREATE licenses/reference/edge-runtime-ponyfill-LICENSE.md.txt
CREATE licenses/reference/edge-runtime-primitives-LICENSE.md.txt
CREATE licenses/reference/babel-core-LICENSE.txt
CREATE licenses/reference/dotenv-LICENSE.txt
CREATE licenses/reference/dotenv-expand-LICENSE.txt
CREATE licenses/reference/unistore-MIT-terms.txt
CREATE licenses/reference/client-only-MIT-terms.txt
CREATE licenses/reference/client-only-package.json.txt
CREATE licenses/reference/babel-packages-inline-notices.txt
CREATE licenses/reference/runtime-notice-lock.json
CREATE licenses/reference/SOURCE_NOTICES.json
CREATE ci/collect_runtime_notices.py
CREATE ci/test_runtime_notices.py
CREATE docs/REFERENCE_RUNTIME_NOTICES.md
CREATE ci/run_observed_reference.py
CREATE docs/LOCAL_REFERENCE_OPERATIONS.md
CREATE ci/build_local_reference.py
CREATE ci/local_identity_fixture.cjs
CREATE ci/local_release.py
CREATE ci/migrate_local_reference.py
CREATE ci/run_local_reference.py
CREATE ci/serve_local_reference.py
CREATE ci/test_local_release.py
CREATE ci/web_local_stop.cjs
CREATE docs/LOCAL_REFERENCE_DELIVERY.md
CREATE ci/gates.example.json
CREATE ci/run_quality_gates.py
CREATE ci/test_run_quality_gates.py
```

## 5. Materialization blocks

### FILE: `ci/gates.example.json`

```yaml
block_id: "PORTABLE-CI:gates-example:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "0373edc1b533fd759e1d1803df006ef1743ae7c55e98f594807266ea06404fda"
variables: []
secrets_allowed: false
```

````json
{
  "schema_version": "elite.ci-gates.v1",
  "required_categories": ["format", "unit", "static", "build", "security", "supply-chain", "integration", "recovery"],
  "fail_fast": true,
  "max_output_bytes": 131072,
  "gates": [
    {"id": "replace-format", "category": "format", "argv": ["replace-me"], "cwd": ".", "timeout_seconds": 300, "required": true}
  ]
}
````

### FILE: `ci/run_quality_gates.py`

```yaml
block_id: "PORTABLE-CI:runner:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "aab21d67cbb9549c2637430a2b4486039b279e52c2e11a4d99d2c5b8e0d72b55"
variables: []
secrets_allowed: false
```

````python
from __future__ import annotations

import argparse
import hashlib
import json
import re
import subprocess
import sys
import time
from datetime import UTC, datetime
from pathlib import Path
from typing import Any

CATEGORIES = {"format", "unit", "property", "static", "build", "security", "supply-chain", "integration", "e2e", "performance", "recovery", "license"}
SECRET_KEY = re.compile(r"(^|_)(password|passwd|token|secret|private_key|api_key|client_secret)($|_)", re.I)
REDACTIONS = [
    re.compile(r"(?i)(authorization:\s*bearer\s+)[^\s]+"),
    re.compile(r"(?i)((?:password|passwd|token|api[_-]?key|client[_-]?secret)\s*[=:]\s*)[^\s]+"),
    re.compile(r"-----BEGIN (?:RSA )?PRIVATE KEY-----.*?-----END (?:RSA )?PRIVATE KEY-----", re.S),
]


def utc_now() -> str:
    return datetime.now(UTC).isoformat().replace("+00:00", "Z")


def redact(text: str, limit: int) -> str:
    value = text[:limit]
    for pattern in REDACTIONS:
        value = pattern.sub(r"\1[REDACTED]" if pattern.groups else "[REDACTED]", value)
    return value + ("\n[TRUNCATED]" if len(text) > limit else "")


def validate(plan: dict[str, Any], root: Path) -> list[str]:
    errors: list[str] = []
    if plan.get("schema_version") != "elite.ci-gates.v1": errors.append("schema_version: unsupported")
    required = plan.get("required_categories", [])
    if not isinstance(required, list) or not required or not set(required) <= CATEGORIES: errors.append("required_categories: invalid")
    limit = plan.get("max_output_bytes")
    if not isinstance(limit, int) or not 4096 <= limit <= 1048576: errors.append("max_output_bytes: must be 4096..1048576")
    gates = plan.get("gates", [])
    ids: set[str] = set()
    covered: set[str] = set()
    for index, gate in enumerate(gates if isinstance(gates, list) else []):
        prefix = f"gates[{index}]"
        gate_id = gate.get("id")
        if not isinstance(gate_id, str) or not re.fullmatch(r"[a-z][a-z0-9-]{1,63}", gate_id) or "replace" in gate_id: errors.append(f"{prefix}.id: invalid")
        elif gate_id in ids: errors.append(f"{prefix}.id: duplicate")
        else: ids.add(gate_id)
        category = gate.get("category")
        if category not in CATEGORIES: errors.append(f"{prefix}.category: invalid")
        elif gate.get("required") is True: covered.add(category)
        argv = gate.get("argv")
        if not isinstance(argv, list) or not argv or not all(isinstance(v, str) and v and "replace-me" not in v for v in argv): errors.append(f"{prefix}.argv: non-empty resolved string array required")
        if "env" in gate or any(SECRET_KEY.search(str(key)) for key in gate): errors.append(f"{prefix}: inline environment/secrets forbidden")
        timeout = gate.get("timeout_seconds")
        if not isinstance(timeout, int) or not 1 <= timeout <= 7200: errors.append(f"{prefix}.timeout_seconds: must be 1..7200")
        cwd = gate.get("cwd", ".")
        try:
            (root / cwd).resolve().relative_to(root.resolve())
        except (ValueError, TypeError):
            errors.append(f"{prefix}.cwd: escapes workspace")
    missing = set(required if isinstance(required, list) else []) - covered
    if missing: errors.append("required_categories: missing required gates: " + ", ".join(sorted(missing)))
    return sorted(errors)


def execute(plan: dict[str, Any], root: Path, plan_hash: str) -> dict[str, Any]:
    started = time.monotonic()
    report: dict[str, Any] = {"schema_version": "elite.ci-evidence.v1", "plan_sha256": plan_hash, "started_at_utc": utc_now(), "workspace": str(root), "gates": []}
    overall = "PASS"
    for gate in plan["gates"]:
        gate_started = time.monotonic()
        result: dict[str, Any] = {"id": gate["id"], "category": gate["category"], "required": gate["required"], "started_at_utc": utc_now()}
        try:
            completed = subprocess.run(gate["argv"], cwd=(root / gate.get("cwd", ".")).resolve(), shell=False, capture_output=True, text=True, encoding="utf-8", errors="replace", timeout=gate["timeout_seconds"], check=False)
            result.update(status="PASS" if completed.returncode == 0 else "FAIL", exit_code=completed.returncode, stdout=redact(completed.stdout, plan["max_output_bytes"]), stderr=redact(completed.stderr, plan["max_output_bytes"]))
        except subprocess.TimeoutExpired as error:
            result.update(status="TIMEOUT", exit_code=None, stdout=redact((error.stdout or "") if isinstance(error.stdout, str) else "", plan["max_output_bytes"]), stderr="gate timed out")
        except OSError as error:
            result.update(status="ERROR", exit_code=None, stdout="", stderr=redact(str(error), plan["max_output_bytes"]))
        result["duration_seconds"] = round(time.monotonic() - gate_started, 3)
        result["finished_at_utc"] = utc_now()
        report["gates"].append(result)
        if gate["required"] and result["status"] != "PASS":
            overall = "FAIL"
            if plan.get("fail_fast", True): break
    report.update(status=overall, duration_seconds=round(time.monotonic() - started, 3), finished_at_utc=utc_now())
    return report


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("plan", type=Path)
    parser.add_argument("--workspace", type=Path, default=Path.cwd())
    parser.add_argument("--evidence", type=Path, required=True)
    args = parser.parse_args()
    raw = args.plan.read_bytes()
    plan = json.loads(raw)
    root = args.workspace.resolve()
    errors = validate(plan, root)
    report = {"schema_version": "elite.ci-evidence.v1", "status": "INVALID", "plan_sha256": hashlib.sha256(raw).hexdigest(), "errors": errors} if errors else execute(plan, root, hashlib.sha256(raw).hexdigest())
    args.evidence.parent.mkdir(parents=True, exist_ok=True)
    args.evidence.write_text(json.dumps(report, indent=2) + "\n", encoding="utf-8")
    print(json.dumps({"status": report["status"], "evidence": str(args.evidence)}))
    return 0 if report["status"] == "PASS" else 2


if __name__ == "__main__":
    sys.exit(main())
````

### FILE: `ci/test_run_quality_gates.py`

```yaml
block_id: "PORTABLE-CI:runner-tests:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "342cc3403e899b584765e7981ffce116f263235d3ac77f9caae2b73ec3ab52eb"
variables: []
secrets_allowed: false
```

````python
from __future__ import annotations

import sys
import tempfile
import unittest
from pathlib import Path

from run_quality_gates import execute, redact, validate


class GateRunnerTests(unittest.TestCase):
    def plan(self, argv: list[str] | None = None) -> dict:
        return {"schema_version": "elite.ci-gates.v1", "required_categories": ["unit"], "fail_fast": True, "max_output_bytes": 4096, "gates": [{"id": "unit-test", "category": "unit", "argv": argv or [sys.executable, "-c", "print('ok')"], "cwd": ".", "timeout_seconds": 10, "required": True}]}

    def test_valid_plan_and_pass(self) -> None:
        with tempfile.TemporaryDirectory() as value:
            root = Path(value)
            plan = self.plan()
            self.assertEqual([], validate(plan, root))
            self.assertEqual("PASS", execute(plan, root, "a" * 64)["status"])

    def test_failure_is_recorded(self) -> None:
        with tempfile.TemporaryDirectory() as value:
            report = execute(self.plan([sys.executable, "-c", "raise SystemExit(7)"]), Path(value), "b" * 64)
            self.assertEqual("FAIL", report["status"])
            self.assertEqual(7, report["gates"][0]["exit_code"])

    def test_timeout_is_recorded(self) -> None:
        with tempfile.TemporaryDirectory() as value:
            plan = self.plan([sys.executable, "-c", "import time; time.sleep(2)"])
            plan["gates"][0]["timeout_seconds"] = 1
            self.assertEqual("TIMEOUT", execute(plan, Path(value), "c" * 64)["gates"][0]["status"])

    def test_missing_category_fails(self) -> None:
        with tempfile.TemporaryDirectory() as value:
            plan = self.plan()
            plan["required_categories"].append("security")
            self.assertTrue(any("missing required gates" in error for error in validate(plan, Path(value))))

    def test_shell_and_inline_env_are_not_available(self) -> None:
        with tempfile.TemporaryDirectory() as value:
            plan = self.plan()
            plan["gates"][0]["env"] = {"TOKEN": "forbidden"}
            self.assertTrue(any("inline environment" in error for error in validate(plan, Path(value))))

    def test_redaction_and_truncation(self) -> None:
        value = redact("Authorization: Bearer abc\npassword=hunter2\n" + "x" * 5000, 4096)
        self.assertNotIn("abc", value)
        self.assertNotIn("hunter2", value)
        self.assertIn("[TRUNCATED]", value)


if __name__ == "__main__":
    unittest.main()
````

## 6. Configuration surface

| Field | Type/default | Secret | Validation/effect |
|---|---|---|---|
| `required_categories` | known string list / none | no | every category must have a gate |
| gate `argv` | non-empty string array | none | no | executed directly; shell syntax is not interpreted |
| gate `cwd` | workspace-relative path / `.` | no | traversal/outside workspace rejected |
| `timeout_seconds` | bounded positive integer | required per gate | no | timeout produces failure evidence |
| `max_output_bytes` | bounded integer / 131072 example | no | output is redacted then truncated |
| `fail_fast` / `required` | boolean | true | no | controls continuation, never converts failure to pass |

## 7. Dependency bill

| Tool | Pin | Use | License | Scope | Source |
|---|---|---|---|---|---|
| Python | `3.14+` standard library | runner/tests | PSF-2.0 | CI | `python.org` |
| project compilers/scanners | exact pins in target CI | gates selected by project | respective licenses | CI | official sources |

## 8. Apply order

Materialize after the project defines its build/test/security/recovery commands. Replace the example plan, pin every tool, run unit/negative tests, execute in an isolated CI worker and retain plan/result hashes with the release. Existing CI should call this runner as one controlled stage. Rollback restores the runner and plan together; historical evidence remains immutable.

## 9. Verification

Clean materialization, six unit tests, intentional invalid-example rejection and an eight-category CLI execution passed; see `reconstruction_evidence/PORTABLE_CI_QUALITY_GATE_RUNNER_2026-08-24_V1.md`. Admission stays conditioned until the target CI uses an isolated runner, injects secrets externally, stores immutable evidence and generates/verifies actual SBOM, provenance and signatures.

## 10. Reconstruction evidence

Toolchain, hashes, positive execution and negative cases are recorded in `reconstruction_evidence/PORTABLE_CI_QUALITY_GATE_RUNNER_2026-08-24_V1.md`; the final library audit rechecks version 0.1.1.

### FILE: `ci/build_local_reference.py`

```yaml
block_id: "PORTABLE-CI-GATE-RUNNER-LOCAL-DELIVERY:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "5b063aaf9c2c673ac5e851b0ee804cc570b1a83b01d98b2673af9f405e52b2a4"
variables: []
secrets_allowed: false
```

````python
"""AUTHORED build glue for the admitted Windows local fixture composition.

No acquisition, live credential use, container build or production admission.
The externally supplied inventory digest is the authority for the source tree.
"""
from __future__ import annotations
import argparse,base64,hashlib,json,os,re,shutil,subprocess,sys,time
from pathlib import Path

GO_SHA='21761eceb9302062c9623fb699f332c8c7fe000f15f70efe8da01a2cfbbc16b9'
NODE_SHA='5c976096e04e5c2c1f091938926234cc9fbebfe9787ddd149351b3b0ecc707b5'
NODE_LICENSE_SHA='ed34dd8e3f0a78dbaf00d0444ce8e285b015b765379c2e17880455f70370f8e9'
HEX=re.compile(r'[0-9a-f]{64}')
def require(ok,why):
    if not ok:raise ValueError(why)
def long_path(p):
    p=os.path.abspath(p)
    return Path(p if p.startswith('\\\\?\\')or os.name!='nt'else '\\\\?\\'+p)
def sha(p):return hashlib.sha256(long_path(p).read_bytes()).hexdigest()
def raw_json(value):return (json.dumps(value,sort_keys=True,indent=2)+'\n').encode()
def read_json(path,expected):
    raw=long_path(path).read_bytes()
    require(HEX.fullmatch(expected) and hashlib.sha256(raw).hexdigest()==expected,'external digest mismatch')
    return json.loads(raw)
def plain(path):
    p=Path(os.path.abspath(path))
    for part in [p,*p.parents]:
        if part.exists():require(not part.is_symlink() and not part.is_junction(),'reparse path rejected')
    return p
def relative(root,rel):
    require(isinstance(rel,str) and '\\'not in rel and ':'not in rel and not rel.startswith('/'),'invalid relative path')
    require(all(x and x not in ('.','..')for x in rel.split('/')),'invalid relative path')
    p=plain(root/rel);p.relative_to(root);return p
def inventory(root,files):
    require(isinstance(files,dict)and files,'empty source inventory')
    seen=set()
    for rel,h in files.items():
        require(rel.casefold()not in seen,'case alias');seen.add(rel.casefold())
        p=relative(root,rel);require(HEX.fullmatch(h)and p.is_file()and sha(p)==h,'source changed: '+rel)
    return hashlib.sha256(raw_json(files)).hexdigest()
def environment(node,work):
    return {'SYSTEMROOT':os.environ['SYSTEMROOT'],'WINDIR':os.environ['SYSTEMROOT'],
      'PATHEXT':'.COM;.EXE;.BAT;.CMD','COMSPEC':str(Path(os.environ['SYSTEMROOT'])/'System32/cmd.exe'),
      'TEMP':str(work),'TMP':str(work),'PATH':str(node.parent)+os.pathsep+str(Path(os.environ['SYSTEMROOT'])/'System32'),
      'CI':'true','NEXT_TELEMETRY_DISABLED':'1','TZ':'UTC','GOTOOLCHAIN':'local','GOPROXY':'off',
      'GOSUMDB':'off','GOWORK':'off','GOFLAGS':'-mod=readonly','GOMAXPROCS':'4','CGO_ENABLED':'0',
      'GOOS':'windows','GOARCH':'amd64','BUSINESS_CONFIG_FILE':'business.example.json',
      'APP_BASE_URL':'https://127.0.0.1:4443','ENTERPRISE_API_BASE_URL':'http://127.0.0.1:9',
      'AUTH_SESSION_SECRET':'synthetic-build-only-not-production-00000000',
      'CATALOG_RELEASE_ENABLED':'false','PUBLIC_INDEXING_ENABLED':'0',
      # Public fixture encryption material makes LOCAL builds comparable.
      # It is forbidden as a live deployment credential.
      'NEXT_SERVER_ACTIONS_ENCRYPTION_KEY':base64.b64encode(hashlib.sha256(b'elite-local-reference-only-actions-v1').digest()).decode()}
def command(label,argv,cwd,env,logs,seconds=300):
    p=logs/(label+'.log');start=time.monotonic()
    with p.open('xb')as f:
        q=subprocess.run(list(map(str,argv)),cwd=cwd,env=env,stdin=subprocess.DEVNULL,stdout=f,stderr=subprocess.STDOUT,
                         timeout=seconds,shell=False,creationflags=subprocess.CREATE_NO_WINDOW)
    require(q.returncode==0,label+' failed; see '+str(p))
    return {'step':label,'exit_code':q.returncode,'log_sha256':sha(p),'seconds':round(time.monotonic()-start,3)}
def tree(root):
    root=long_path(root)
    result={}
    for p in sorted(root.rglob('*')):
        require(not p.is_symlink()and not p.is_junction(),'artifact link rejected')
        if p.is_file():result[p.relative_to(root).as_posix()]=sha(p)
    return result
def copy_tree(source,target):
    # Runtime copies contain regular files only; junctions are created afterward
    # from the exact framework trace, inside the private activation directory.
    source=long_path(source);target=long_path(target);count=0;total=0
    def copy(a,b):
        nonlocal count,total
        b.mkdir(parents=True,exist_ok=True)
        with os.scandir(a)as entries:
            for entry in entries:
                info=entry.stat(follow_symlinks=False)
                require(not getattr(info,'st_file_attributes',0)&0x400,'runtime copy contains a link')
                if entry.is_dir(follow_symlinks=False):copy(Path(entry.path),b/entry.name)
                else:
                    count+=1;total+=info.st_size
                    require(count<=50000 and total<=536870912,'runtime copy budget exceeded')
                    shutil.copyfile(entry.path,b/entry.name)
    copy(source,target)
def project_standalone(src,target):
    source=long_path(src/'.next/standalone');target=long_path(target);installed=long_path(src/'node_modules').resolve();links={};files=0
    def copy(a,b):
        nonlocal files
        if a.is_symlink()or a.is_junction():
            actual=a.resolve(strict=True)
            if source in actual.parents:rel=actual.relative_to(source).as_posix()
            else:
                require(installed in actual.parents,'framework link escapes qualified consumer')
                rel='node_modules/'+actual.relative_to(installed).as_posix()
            require(actual.is_dir(),'only directory junction projection supported')
            links[a.relative_to(source).as_posix()]=rel;return
        if a.is_dir():
            b.mkdir(parents=True,exist_ok=True)
            for child in sorted(a.iterdir()):copy(child,b/child.name)
        else:
            files+=1;require(files<=20000,'standalone file budget');shutil.copyfile(a,b)
    copy(source,target)
    for rel,dest in links.items():require((target/dest).is_dir(),'traced link target absent: '+rel)
    (target/'elite-runtime-links.json').write_bytes(raw_json(links))
def go_notices(src,dist,go,env,logs):
    step=command('go-modules',[go,'list','-m','-json','all'],src,env,logs)
    raw=(logs/'go-modules.log').read_text(encoding='utf-8');decoder=json.JSONDecoder();modules=[]
    while raw.strip():
        value,n=decoder.raw_decode(raw.lstrip());raw=raw.lstrip()[n:];modules.append(value)
    result=[]
    for value in modules:
        if value.get('Main'):continue
        spec=value.get('Replace',value);root=Path(spec.get('Dir',''));texts={}
        require(root.is_dir(),'module source absent')
        for item in root.iterdir():
            if item.is_file()and re.fullmatch(r'(?:license|licence|notice|copying)(?:[.-].*)?',item.name,re.I):
                h=sha(item);p=dist/'dependency-notices'/h;p.parent.mkdir(exist_ok=True)
                if not p.exists():shutil.copyfile(item,p)
                texts[item.name]=h
        result.append({'module':value['Path'],'version':spec.get('Version'),'sum':spec.get('Sum'),'texts':texts,
          'local_corresponding_source':bool(value.get('Replace'))})
    p=go.parent.parent/'LICENSE';h=sha(p);shutil.copyfile(p,dist/'GO_LICENSE.txt')
    (dist/'go-dependency-notices.json').write_bytes(raw_json({'modules':result,'go_license_sha256':h}))
    return step
def corresponding_sources(source,built,dist,files):
    # Keep exact admitted inputs. Next writes generated type declarations in
    # its workspace; retain that one checked derivation separately.
    inventory(source,files);derived={}
    for rel,h in files.items():
        actual=sha(built/rel)
        if actual!=h:
            require(rel=='next-env.d.ts','unexpected build source mutation: '+rel)
            original=(source/rel).read_bytes();anchor=b'/// <reference types="next/image-types/global" />\n'
            addition=b'import "./.next/types/routes.d.ts";\nimport "./.next/types/root-params.d.ts";\n'
            require(original.count(anchor)==1 and (built/rel).read_bytes()==original.replace(anchor,anchor+addition),
                    'unknown Next generated declaration')
            dest=dist/'build-derived-source'/rel;dest.parent.mkdir(parents=True,exist_ok=True);shutil.copyfile(built/rel,dest)
            require(sha(dest)==actual,'derived source copy drift')
            derived[rel]={'input_sha256':h,'generated_sha256':actual,'generator':'Next16.3.4 canonical route type declaration'}
        dest=dist/'source'/rel;dest.parent.mkdir(parents=True,exist_ok=True);shutil.copyfile(source/rel,dest)
        require(sha(dest)==h,'corresponding source copy drift')
    if derived:(dist/'build-derived-source/manifest.json').write_bytes(raw_json(derived))
    inventory(source,files)
    return derived

def package_artifact(src,dist,node,source_id,files,source):
    standalone=src/'.next/standalone';require((standalone/'server.js').is_file(),'standalone frontend missing')
    project_standalone(src,dist/'web')
    for rel in ['.next/static','public','config']:
        if(src/rel).is_dir():copy_tree(src/rel,dist/'web'/rel)
    shutil.copytree(src/'db/migrations',dist/'migrations')
    corresponding_sources(source,src,dist,files)
    for rel in files:
        if rel.startswith('licenses/')or rel.startswith('docs/provenance/')or rel=='THIRD_PARTY_NOTICES.md':
            p=dist/rel;p.parent.mkdir(parents=True,exist_ok=True);shutil.copyfile(source/rel,p)
    require(sha(node.parent/'LICENSE')==NODE_LICENSE_SHA,'Node license pin changed')
    shutil.copyfile(node,dist/'node.exe');shutil.copyfile(node.parent/'LICENSE',dist/'NODE_LICENSE.txt')
    from collect_runtime_notices import collect
    collect(src,dist/'web',dist,src/'licenses/reference/runtime-notice-lock.json')
    # Framework-generated metadata contains build-machine absolute paths.
    # Only these documented non-serving locator fields are normalized;
    # receipts keep original hashes. Never rewrite executable app bundles.
    normalization=[]
    for nft in sorted(long_path(dist/'web/.next').rglob('*.nft.json')):
        value=json.loads(nft.read_text());require(set(value)=={'version','files'}and value['version']==1,'unknown Next trace format')
        require(isinstance(value['files'],list)and all(isinstance(x,str)for x in value['files']),'invalid Next trace files')
        normalization.append({'path':nft.relative_to(long_path(dist)).as_posix(),'before_sha256':sha(nft),'field':'files sort order'})
        value['files']=sorted(value['files']);nft.write_bytes(raw_json(value))
    p=dist/'web/server.js';text=p.read_text();m=re.search(r'^const nextConfig = (.+)$',text,re.M)
    require(m is not None,'unknown Next standalone server format')
    config=json.loads(m[1]);config['outputFileTracingRoot']='.';config['repoRoot']='.';config['turbopack']['root']='.'
    normalization.append({'path':'web/server.js','before_sha256':sha(p),'fields':['nextConfig.outputFileTracingRoot','nextConfig.repoRoot','nextConfig.turbopack.root']})
    p.write_text(text[:m.start(1)]+json.dumps(config,separators=(',',':'))+text[m.end(1):],encoding='utf-8',newline='\n')
    p=dist/'web/.next/required-server-files.json'
    if p.exists():
        v=json.loads(p.read_text());v['appDir']='.';v['config']['outputFileTracingRoot']='.';v['config']['repoRoot']='.';v['config']['turbopack']['root']='.'
        normalization.append({'path':'web/.next/required-server-files.json','before_sha256':sha(p),'fields':['appDir','config.outputFileTracingRoot','config.repoRoot','config.turbopack.root']})
        p.write_bytes(raw_json(v))
    (dist/'LOCAL_REFERENCE_ONLY.json').write_bytes(raw_json({'schema':'elite-local-reference-release/v1',
      'source_sha256':source_id,'scope':'LOCAL_FIXTURES','production_admitted':False,
      'fixture_credentials':True,'go_sha256':GO_SHA,'node_sha256':NODE_SHA}))
    return normalization
def build(args):
    source=plain(args.source);target=plain(args.target);workspace=plain(args.workspace);go=plain(args.go);node=plain(args.node)
    require(not target.exists()and target.parent.is_dir(),'absent target required')
    require(source!=target and source not in target.parents and target not in source.parents,'overlapping source/output')
    require(not workspace.exists()and workspace.parent.is_dir(),'absent stable build workspace required')
    for other in [source,target]:require(workspace!=other and workspace not in other.parents and other not in workspace.parents,'build workspace overlaps input/output')
    files=read_json(args.inventory,args.inventory_sha256);source_id=inventory(source,files)
    require(sha(go)==GO_SHA and sha(node)==NODE_SHA,'toolchain pin mismatch')
    plan_inputs=read_json(args.install_inputs,args.install_inputs_sha256)
    require(set(plan_inputs)=={'projection','contained','store','cache','acquisition_receipt','acquisition_sha256'},'exact install input fields required')
    runtime=source/'pnpm_artifact_selection';require(runtime.is_dir(),'materialize PNPM_ARTIFACT_SELECTION_GATE first')
    target.mkdir();work=target/'work';work.mkdir();logs=target/'logs';logs.mkdir();src=workspace;src.mkdir()
    result={'schema':'elite-local-reference-build/v1','scope':'LOCAL_FIXTURES','state':'FAIL','source_sha256':source_id,'build_workspace':str(workspace),'steps':[]}
    try:
        for rel in files:
            p=src/rel;p.parent.mkdir(parents=True,exist_ok=True);shutil.copyfile(source/rel,p)
        # Execute the existing exact three-consumer admission, never raw pnpm.
        sys.path.insert(0,str(runtime));import local_runtime
        inputs={**plan_inputs,'consumer':'enterprise-web','project':str(src),'node':str(node),'target':str(target/'install')}
        plan_sha=local_runtime.prepare(**inputs);installed=local_runtime.execute(target/'install',plan_sha)
        require(installed['state']=='PASS','restricted install failed')
        result['install']={'plan_sha256':plan_sha,'receipt_sha256':sha(target/'install/execution-result.json')}
        env=environment(node,work);env.update(GOMODCACHE=str(plain(args.go_mod_cache)),GOCACHE=str(target/'go-cache'),
            USERPROFILE=str(work),ELITE_SOURCE_SHA256=source_id)
        dist=target/'artifact';dist.mkdir();(dist/'api').mkdir()
        from build_arca_reference import build_arca
        result['arca']=build_arca(src,dist,work,logs,go,env,args.arca_inputs,args.arca_inputs_sha256)
        result['steps'].append(command('api',[go,'build','-trimpath','-buildvcs=false','-ldflags=-s -w -buildid=',
                    '-o',dist/'api/electromobility-api.exe','./cmd/electromobility-api'],src,env,logs))
        # Next16.3.4 preserves its preview cache during build. These public
        # fixture values are never production credentials or runtime admission.
        preview={name:hashlib.sha256(('elite-local-only:'+name).encode()).hexdigest()[:size]
            for name,size in [('previewModeId',32),('previewModeSigningKey',64),('previewModeEncryptionKey',64)]}
        preview['expireAt']=4102444800000
        cache=src/'.next/cache';cache.mkdir(parents=True)
        (cache/'.previewinfo').write_bytes(raw_json(preview))
        result['steps'].append(command('web',[node,src/'node_modules/next/dist/bin/next','build','--webpack'],src,env,logs))
        actual=json.loads((src/'.next/prerender-manifest.json').read_text())['preview']
        require(actual=={k:v for k,v in preview.items()if k!='expireAt'},'Next preview fixture contract changed')
        normalization=package_artifact(src,dist,node,source_id,files,source)
        result['steps'].append(go_notices(src,dist,go,env,logs))
        result['normalization']=normalization;result['files']=tree(dist)
        result['artifact_sha256']=hashlib.sha256(raw_json(result['files'])).hexdigest()
        inventory(source,files);result['state']='PASS'
        (target/'artifact-inventory.json').write_bytes(raw_json(result['files']))
    finally:(target/'build-result.json').write_bytes(raw_json(result))
    return result
def main():
    p=argparse.ArgumentParser(description=__doc__)
    for name in ['source','target','workspace','inventory','inventory-sha256','go','node','go-mod-cache','install-inputs','install-inputs-sha256','arca-inputs','arca-inputs-sha256']:p.add_argument('--'+name,required=True)
    a=p.parse_args()
    try:r=build(a);print('LOCAL_REFERENCE_BUILD_PASS '+r['artifact_sha256']);return 0
    except (ValueError,OSError,subprocess.SubprocessError)as e:print('LOCAL_REFERENCE_BUILD_FAIL '+str(e),file=sys.stderr);return 2
if __name__=='__main__':raise SystemExit(main())
````

### FILE: `ci/local_identity_fixture.cjs`

```yaml
block_id: "PORTABLE-CI-GATE-RUNNER-LOCAL-DELIVERY:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "aa0e624fb5c69b1f6301c17c790e6186f4edee4e4c38f729db02d3316b192168"
variables: []
secrets_allowed: false
```

````text
// AUTHORED test fixture using Node's existing crypto/http standard library.
// Synthetic loopback identity only; never a deployable identity provider.
'use strict';
const http = require('node:http');
const crypto = require('node:crypto');
const args = process.argv.slice(2);
const ttl = args.length === 0 ? 600 : args.length === 2 && args[0] === '--ttl-seconds' && /^(?:[1-9][0-9]{0,3})$/.test(args[1]) ? Number(args[1]) : 0;
if (ttl < 1 || ttl > 1200) throw new Error('bounded fixture lifetime required');
const { publicKey, privateKey } = crypto.generateKeyPairSync('rsa', { modulusLength: 2048 });
const key = { ...publicKey.export({ format: 'jwk' }), kid: 'local-reference', use: 'sig', alg: 'RS256' };
let issuer;
const server = http.createServer((req, res) => {
  res.setHeader('Content-Type', 'application/json');
  if (req.method !== 'GET') { res.writeHead(405).end('{}'); return; }
  if (req.url === '/.well-known/openid-configuration') {
    res.end(JSON.stringify({ issuer, jwks_uri: issuer + '/jwks', authorization_endpoint: issuer + '/authorize',
      token_endpoint: issuer + '/token', response_types_supported: ['code'], subject_types_supported: ['public'],
      id_token_signing_alg_values_supported: ['RS256'] }));
  } else if (req.url === '/jwks') res.end(JSON.stringify({ keys: [key] }));
  else if (req.url === '/fixture-token') {
    const now = Math.floor(Date.now() / 1000);
    const claims = { iss: issuer, aud: 'elite-local-reference', sub: 'local-delivery-customer', iat: now, exp: now + 120,
      tenant_id: 'b0000000-0000-4000-8000-000000000001', organization_ids: ['b0000000-0000-4000-8000-000000000002'],
      permissions: ['order:create'] };
    const part = obj => Buffer.from(JSON.stringify(obj)).toString('base64url');
    const input = part({ alg: 'RS256', kid: key.kid, typ: 'JWT' }) + '.' + part(claims);
    res.end(JSON.stringify({ token: input + '.' + crypto.sign('RSA-SHA256', Buffer.from(input), privateKey).toString('base64url') }));
  } else res.writeHead(404).end('{}');
});
server.listen(0, '127.0.0.1', () => {
  issuer = 'http://127.0.0.1:' + server.address().port;
  process.stdout.write(JSON.stringify({ issuer }) + '\n');
});
process.stdin.on('data', () => server.close(() => process.exit(0)));
process.stdin.on('end', () => server.close(() => process.exit(0)));
setTimeout(() => server.close(() => process.exit(0)), ttl * 1000).unref();
````

### FILE: `ci/local_release.py`

```yaml
block_id: "PORTABLE-CI-GATE-RUNNER-LOCAL-DELIVERY:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "26120750ef65f3d7e39eb11fa6aa62f7928b8dd82ed217f808dd1165cd1a94f6"
variables: []
secrets_allowed: false
```

````python
"""AUTHORED local activation journal around the existing Windows job launcher.

Single trusted local operator. No SCM registration, cloud rollout, hostile-user
isolation or secret storage. Failed candidate health restores the prior process
and pointer; no automatic database downgrade is ever issued.
"""
from __future__ import annotations
import ctypes,hashlib,json,msvcrt,os,re,subprocess,sys,threading,time,uuid
from pathlib import Path
from urllib.parse import urlsplit
from urllib.request import build_opener,ProxyHandler
from build_local_reference import plain,require,read_json,relative,sha,raw_json,tree

def loopback(value,scheme):
    u=urlsplit(value)
    require(u.scheme==scheme and u.hostname=='127.0.0.1'and u.port and 1024<=u.port<=65535 and not u.password and not u.fragment,'owned loopback endpoint required')
    return u
def runtime_config(value):
    require(set(value) in ({'database_url','issuer','api_port','web_port'},{'database_url','issuer','api_port','web_port','metrics_port'}),'runtime config fields')
    u=loopback(value['database_url'],'postgres')
    require(u.username=='postgres'and re.fullmatch(r'/elite_payment_connected_[0-9a-f]{32}',u.path)and u.query=='sslmode=disable','synthetic reference database required')
    u=loopback(value['issuer'],'http');require(not u.username and u.path in ('','/'),'fixture issuer required')
    ports=[value[x]for x in ['api_port','web_port','metrics_port']if x in value]
    require(all(type(x)is int and 1024<=x<=65535 for x in ports)and len(set(ports))==len(ports),'distinct loopback ports required')
    return value
def verify_release(root,inventory_path,digest):
    root=plain(root);files=read_json(inventory_path,digest)
    require(files and tree(root)==files,'artifact inventory mismatch')
    profile=json.loads((root/'LOCAL_REFERENCE_ONLY.json').read_text())
    require(profile['scope']=='LOCAL_FIXTURES'and profile['production_admitted']is False,'local fixture artifact required')
    require(sha(root/'node.exe')==profile['node_sha256'],'node binding changed')
    return {'root':str(root),'inventory':str(plain(inventory_path)),'inventory_sha256':digest,'source_sha256':profile['source_sha256']}
def atomic(path,value):
    temp=path.with_name(path.name+'.'+uuid.uuid4().hex+'.tmp')
    with temp.open('xb')as f:f.write(raw_json(value));f.flush();os.fsync(f.fileno())
    os.replace(temp,path)
def probe(config,seconds=60,alive=lambda:True):
    client=build_opener(ProxyHandler({}));end=time.monotonic()+seconds;checks={}
    while time.monotonic()<end:
        require(alive(),'candidate process exited before health')
        try:
            for name,url in [('api',f"http://127.0.0.1:{config['api_port']}/health/live"),('web',f"http://127.0.0.1:{config['web_port']}/")]:
                with client.open(url,timeout=2)as r:
                    body=r.read(1048576);require(r.status==200,'health status')
                    require(json.loads(body).get('status')=='ok'if name=='api'else b'<html'in body.lower(),'health body')
                    checks[name]={'status':r.status,'body_sha256':hashlib.sha256(body).hexdigest()}
            return checks
        except (OSError,ValueError):time.sleep(.1)
    raise ValueError('candidate health timeout')
class Instance:
    def __init__(self,release,config,root):
        self.cancel=threading.Event();self.result=None;self.config=config;self.root=root
        source=Path(__file__).resolve().parent.parent;sys.path.insert(0,str(source/'reference_telemetry'));import governed_launch
        script=source/'ci/serve_local_reference.py'
        job={'release':release,'runtime':config,'worker_sha256':sha(script)}
        config_path=root/'worker.json';atomic(config_path,job)
        exe=plain(sys.executable);runtime_env={k:str(Path(os.environ['SYSTEMROOT']))for k in ['SYSTEMROOT','WINDIR']}
        runtime_env.update(TEMP=str(root),TMP=str(root))
        profile={'schema':'elite-native-launch-qualification/v2','id':'local-franchise',
          'executable':{'path':str(exe),'sha256':sha(exe),'bytes':exe.stat().st_size},
          'arguments':['-X','utf8','-B',str(script),'--profile',str(config_path),'--sha256',sha(config_path)],
          'cwd':str(root),'environment':runtime_env,
          'budgets':{'timeout':1200 if 'metrics_port'in config else 600,'output_bytes':2097152,'processes':12,'commit_bytes':1073741824},
          'shutdown':{'protocol':'win32-inherited-event/v1','grace_seconds':25}}
        raw=raw_json(profile);atomic(root/'launch.json',profile)
        def run():self.result=governed_launch.launch(raw,hashlib.sha256(raw).hexdigest(),cancel=self.cancel)
        self.thread=threading.Thread(target=run);self.thread.start()
    def stop(self):
        self.cancel.set();self.thread.join(35);require(not self.thread.is_alive(),'launcher did not finish')
        r=self.result;require(r and r.status=='CAPTURED'and r.outcome and r.outcome.tree_empty,'process tree not empty')
        out=r.outcome
        (self.root/'stdout.log').write_bytes(out.stdout);(self.root/'stderr.log').write_bytes(out.stderr)
        receipt={'status':out.status,'exit_code':out.exit_code,'tree_empty':out.tree_empty,
          'shutdown_state':out.shutdown_state,'observed_peak_job_bytes':out.observed_peak_job_bytes,
          'stdout_sha256':sha(self.root/'stdout.log'),'stderr_sha256':sha(self.root/'stderr.log')}
        atomic(self.root/'stop.json',receipt);return receipt
class Deployment:
    def __init__(self,root,config):
        self.root=plain(root);self.root.mkdir(exist_ok=True);self.config=runtime_config(config);self.instance=None;self.serial=0
        self.lock=(self.root/'operator.lock').open('a+b');self.lock.seek(0)
        try:msvcrt.locking(self.lock.fileno(),msvcrt.LK_NBLCK,1)
        except OSError:self.lock.close();raise ValueError('another local operator owns deployment')
        if os.fstat(self.lock.fileno()).st_size==0:self.lock.write(b'0');self.lock.flush()
    def _start(self,release):
        verify_release(release['root'],release['inventory'],release['inventory_sha256'])
        self.serial+=1;work=self.root/('run-'+uuid.uuid4().hex);work.mkdir()
        self.instance=Instance(release,self.config,work)
        return probe(self.config,alive=self.instance.thread.is_alive)
    def _stop(self):
        if self.instance:
            instance=self.instance;self.instance=None;return instance.stop()
    def current(self):
        p=self.root/'current.json';return json.loads(p.read_text())if p.exists()else None
    def activate(self,artifact,inventory_path,digest):
        release=verify_release(artifact,inventory_path,digest);prior=self.current()
        if prior:verify_release(prior['root'],prior['inventory'],prior['inventory_sha256'])
        receipt={'schema':'elite-local-activation/v1','state':'STARTING','previous':prior,'candidate':release,'migration_action':'NONE'}
        journal=self.root/('activation-'+uuid.uuid4().hex+'.json');atomic(journal,receipt)
        try:
            receipt['stopped_previous']=self._stop();receipt['health']=self._start(release)
            atomic(self.root/'current.json',release);receipt['state']='ACTIVE'
        except BaseException:
            receipt['stopped_candidate']=self._stop()
            if prior:
                receipt['rollback_health']=self._start(prior);receipt['state']='ROLLED_BACK'
                require(self.current()==prior,'stable pointer changed on failed candidate')
            else:receipt['state']='FAILED_NO_PREVIOUS'
            raise
        finally:atomic(journal,receipt)
        return receipt
    def recover(self):
        prior=self.current();require(prior is not None,'no committed release to recover')
        self._stop();return self._start(prior)
    def close(self):
        try:return self._stop()
        finally:
            self.lock.seek(0);msvcrt.locking(self.lock.fileno(),msvcrt.LK_UNLCK,1);self.lock.close()
````

### FILE: `ci/migrate_local_reference.py`

```yaml
block_id: "PORTABLE-CI-GATE-RUNNER-LOCAL-DELIVERY:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "ee7e6e1789398428d7069737eadbdc1ecbedae8b37152d1e63795631b4174fd9"
variables: []
secrets_allowed: false
```

````python
"""AUTHORED psql orchestration: hash ledger, session lock, atomic migration steps.

Only an owned loopback fixture database is accepted. A pre-existing unregistered
schema requires baseline admission; the runner never guesses it or downgrades.
"""
from pathlib import Path
import argparse,hashlib,json,os,re,subprocess,sys
from build_local_reference import plain,sha,raw_json,require,read_json,relative
from local_release import loopback

def program(root,files):
    names=sorted(files);require(names and all(re.fullmatch(r'[0-9]{4}_[a-z0-9_]+\.up\.sql',x)for x in names),'migration inventory required')
    prefix="""\\set ON_ERROR_STOP on
select pg_advisory_lock(2808,402);
do $guard$ begin
 if to_regclass('library_delivery.migration') is null and to_regclass('platform.tenant') is not null then
  raise exception 'EXISTING_SCHEMA_REQUIRES_BASELINE';
 end if;
end $guard$;
create schema if not exists library_delivery;
create table if not exists library_delivery.migration(name text primary key, ordinal integer unique not null, sha256 text not null, applied_at timestamptz not null default clock_timestamp());
"""
    allowed=','.join("'"+n+"'"for n in names)
    prefix+=f"select count(*)=0 as no_unknown from library_delivery.migration where name not in ({allowed}) \\gset\n\\if :no_unknown\n\\else\n\\quit 3\n\\endif\n"
    for i,name in enumerate(names,1):
        p=relative(root,name);h=files[name];raw=p.read_bytes();require(hashlib.sha256(raw).hexdigest()==h,'migration changed: '+name);body=raw.decode('utf-8')
        begin=re.match(r'\s*(?:(?:--[^\n]*\n)\s*)*begin\s*;',body,re.I);end=re.search(r'commit\s*;\s*$',body,re.I)
        require(bool(begin)==bool(end),'asymmetric transaction wrapper')
        if begin:body=body[begin.end():end.start()]
        prefix+=f"select not exists(select 1 from library_delivery.migration where name='{name}' and (sha256<>'{h}' or ordinal<>{i})) as matches \\gset\n\\if :matches\n\\else\n\\quit 4\n\\endif\n"
        prefix+=f"select exists(select 1 from library_delivery.migration where name='{name}') as applied \\gset\n\\if :applied\n\\else\nbegin;\n{body}\ninsert into library_delivery.migration(name,ordinal,sha256) values('{name}',{i},'{h}');\ncommit;\n\\endif\n"
    prefix+=f"select count(*)={len(names)} as complete from library_delivery.migration \\gset\n\\if :complete\n\\echo LOCAL_MIGRATIONS_PASS\n\\else\n\\quit 5\n\\endif\n"
    return prefix.encode()
def migrate(root,files,database,psql,psql_sha256,output):
    u=loopback(database,'postgres');require(u.username=='postgres'and re.fullmatch(r'/elite_payment_connected_[0-9a-f]{32}',u.path)and u.query=='sslmode=disable','owned fixture database required')
    root=plain(root);output=plain(output);psql=plain(psql)
    require(not output.exists()and output.parent.is_dir(),'absent receipt directory required');require(sha(psql)==psql_sha256,'psql pin changed')
    source=program(root,files);output.mkdir();(output/'apply.sql').write_bytes(source)
    env={'SYSTEMROOT':os.environ['SYSTEMROOT'],'WINDIR':os.environ['SYSTEMROOT'],'PGCONNECT_TIMEOUT':'3','PGCLIENTENCODING':'UTF8'}
    result={'schema':'elite-local-migration/v1','scope':'LOCAL_FIXTURES','state':'FAIL','database':database,'files':files,'psql_sha256':psql_sha256,'program_sha256':sha(output/'apply.sql')}
    try:
        with(output/'apply.log').open('xb')as log:
            q=subprocess.run([psql,'-X','-w','--dbname',database,'--file',output/'apply.sql'],env=env,stdout=log,stderr=subprocess.STDOUT,shell=False,timeout=120,creationflags=subprocess.CREATE_NO_WINDOW)
        result['exit_code']=q.returncode;result['log_sha256']=sha(output/'apply.log')
        require(q.returncode==0 and b'LOCAL_MIGRATIONS_PASS'in(output/'apply.log').read_bytes(),'migration failed; preserve ledger/log and resume after canonical correction')
        result['state']='PASS';return result
    finally:(output/'receipt.json').write_bytes(raw_json(result))
def main():
    p=argparse.ArgumentParser(description=__doc__)
    for name in ['root','inventory','inventory-sha256','database-url','psql','psql-sha256','receipt-directory']:p.add_argument('--'+name,required=True)
    a=p.parse_args()
    try:migrate(a.root,read_json(a.inventory,a.inventory_sha256),a.database_url,a.psql,a.psql_sha256,a.receipt_directory);print('LOCAL_MIGRATIONS_PASS');return 0
    except (ValueError,OSError,subprocess.SubprocessError)as e:print('LOCAL_MIGRATIONS_FAIL '+str(e),file=sys.stderr);return 2
if __name__=='__main__':raise SystemExit(main())
````

### FILE: `ci/run_local_reference.py`

```yaml
block_id: "PORTABLE-CI-GATE-RUNNER-LOCAL-DELIVERY:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "6a79f9f70c65dfc88e5039b90759ea1fdb6383ee0b45052b9a1b9b21eccab8a4"
variables: []
secrets_allowed: false
```

````python
"""Finite local operator command. No provider activation or production effects."""
from pathlib import Path
import argparse,json,time
from local_release import Deployment,read_json,require,sha,verify_release
def main():
    p=argparse.ArgumentParser(description=__doc__)
    for name in ['artifact','inventory','inventory-sha256','runtime','runtime-sha256','migration-receipt','migration-receipt-sha256','state-directory']:p.add_argument('--'+name,required=True)
    p.add_argument('--seconds',type=int,default=300);a=p.parse_args();require(1<=a.seconds<=480,'finite duration 1..480 required')
    cfg=read_json(a.runtime,a.runtime_sha256);migration=read_json(a.migration_receipt,a.migration_receipt_sha256)
    require(migration['state']=='PASS'and migration['scope']=='LOCAL_FIXTURES'and migration['database']==cfg['database_url'],'migration receipt mismatch')
    files={p.name:sha(p)for p in(Path(a.artifact)/'migrations').glob('*.up.sql')};require(files==migration['files'],'migration receipt does not cover release')
    dep=Deployment(a.state_directory,cfg)
    try:
        dep.activate(a.artifact,a.inventory,a.inventory_sha256)
        print(f"LOCAL_REFERENCE_ACTIVE http://127.0.0.1:{cfg['web_port']} duration={a.seconds}",flush=True)
        deadline=time.monotonic()+a.seconds
        while time.monotonic()<deadline:
            require(dep.instance.thread.is_alive(),'application exited');time.sleep(.2)
    finally:
        receipt=dep.close();require(receipt and receipt['exit_code']==0,'local shutdown failed')
if __name__=='__main__':main()
````

### FILE: `ci/serve_local_reference.py`

```yaml
block_id: "PORTABLE-CI-GATE-RUNNER-LOCAL-DELIVERY:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "37e6952d1026dd12e9786dbe3ee12eee8f0a2020fa9abac24a662625329c161e"
variables: []
secrets_allowed: false
```

````python
"""AUTHORED child of the existing finite Windows Job Object launcher."""
from pathlib import Path
import argparse,ctypes,os,shutil,subprocess,sys,time
from local_release import read_json,runtime_config,verify_release,require,sha
from build_local_reference import environment,copy_tree

def main():
    p=argparse.ArgumentParser();p.add_argument('--profile',required=True);p.add_argument('--sha256',required=True);a=p.parse_args()
    job=read_json(a.profile,a.sha256);require(set(job)=={'release','runtime','worker_sha256'},'worker profile fields')
    require(sha(Path(__file__))==job['worker_sha256'],'worker source changed')
    r=job['release'];verify_release(r['root'],r['inventory'],r['inventory_sha256']);root=Path(r['root']);cfg=runtime_config(job['runtime'])
    raw=os.environ.get('ELITE_STOP_EVENT_HANDLE','');require(raw.isdecimal()and int(raw)>0,'trusted inherited stop event required');handle=int(raw)
    wait=ctypes.WinDLL('kernel32',use_last_error=True).WaitForSingleObject
    wait.argtypes=[ctypes.c_void_p,ctypes.c_uint32];wait.restype=ctypes.c_uint32
    env=environment(root/'node.exe',Path.cwd());env.update(DATABASE_URL=cfg['database_url'],OIDC_ISSUER=cfg['issuer'],OIDC_AUDIENCE='elite-local-reference',
       HTTP_ADDRESS=f"127.0.0.1:{cfg['api_port']}",HOSTNAME='127.0.0.1',PORT=str(cfg['web_port']),NODE_ENV='production',
       APP_BASE_URL=f"http://127.0.0.1:{cfg['web_port']}",ENTERPRISE_API_BASE_URL=f"http://127.0.0.1:{cfg['api_port']}")
    require(wait(handle,0)==258,'stop requested before launch')
    if 'metrics_port'in cfg:
        env.update(HTTP_METRICS_ENABLED='true',HTTP_METRICS_ADDRESS=f"127.0.0.1:{cfg['metrics_port']}")
    procs=[]
    try:
        # Restrict extra inheritance to the selected stop handle for the Go host.
        os.set_handle_inheritable(handle,True);startup=subprocess.STARTUPINFO();startup.lpAttributeList={'handle_list':[handle]}
        api_env={**env,'ELITE_STOP_EVENT_HANDLE':raw}
        api=subprocess.Popen([root/'api/electromobility-api.exe'],cwd=root,env=api_env,stdin=subprocess.DEVNULL,
            stdout=sys.stdout,stderr=sys.stderr,startupinfo=startup,creationflags=subprocess.CREATE_NO_WINDOW)
        procs.append(api);os.set_handle_inheritable(handle,False)
        webroot=Path.cwd()/'web';copy_tree(root/'web',webroot)
        script=Path(__file__).with_name('web_local_stop.cjs')
        web=subprocess.Popen([root/'node.exe',script,webroot/'server.js'],cwd=webroot,env=env,stdin=subprocess.PIPE,
            stdout=sys.stdout,stderr=sys.stderr,creationflags=subprocess.CREATE_NO_WINDOW)
        procs.append(web)
        while wait(handle,50)==258:
            require(all(x.poll()is None for x in procs),'application exited unexpectedly')
        if web.poll()is None:web.stdin.write(b'STOP\n');web.stdin.flush();web.stdin.close()
        # Go uses its inherited event; Next handles its own SIGTERM listeners.
        deadline=time.monotonic()+20
        for proc,expected in [(api,0),(web,143)]:
            proc.wait(max(.1,deadline-time.monotonic()));require(proc.returncode==expected,'application shutdown failed')
        print('LOCAL_REFERENCE_STOPPED',flush=True);return 0
    finally:
        for proc in procs:
            if proc.poll()is None:proc.terminate()
        for proc in procs:proc.wait(5)
if __name__=='__main__':raise SystemExit(main())
````

### FILE: `ci/test_local_release.py`

```yaml
block_id: "PORTABLE-CI-GATE-RUNNER-LOCAL-DELIVERY:file7:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "2518f7ea15761580823854b93070adaf8a7eba9a55126aa5eb30a3327af9b821"
variables: []
secrets_allowed: false
```

````python
"""Mutation tests for local activation integrity; native connected test is separate."""
from pathlib import Path
import hashlib,json,tempfile,unittest
from local_release import Deployment,runtime_config,verify_release,raw_json,sha,tree
class LocalReleaseTests(unittest.TestCase):
    def test_external_digest_required(self):
        with tempfile.TemporaryDirectory()as d:
            p=Path(d);f=p/'manifest.json';f.write_text('{}')
            with self.assertRaisesRegex(ValueError,'digest'):verify_release(p,f,'0'*64)
    def test_remote_database_and_unknown_env_rejected(self):
        base={'database_url':'postgres://postgres@127.0.0.1:6543/elite_payment_connected_'+'a'*32+'?sslmode=disable','issuer':'http://127.0.0.1:6544','api_port':6545,'web_port':6546}
        self.assertEqual(runtime_config(base),base)
        self.assertEqual(runtime_config(base|{'metrics_port':6547})['metrics_port'],6547)
        for port in [True,0,6545,6546,65536,'6547']:
            with self.subTest(metrics_port=port),self.assertRaises(ValueError):runtime_config(base|{'metrics_port':port})
        for patch in [{'database_url':base['database_url'].replace('127.0.0.1','example.com')},{'issuer':'http://example.com:6544'},{'api_port':6546},{'provider_key':'unexpected'}]:
            with self.subTest(patch=patch),self.assertRaises(ValueError):runtime_config(base|patch)
    def test_exclusive_operator(self):
        with tempfile.TemporaryDirectory()as d:
            cfg={'database_url':'postgres://postgres@127.0.0.1:6543/elite_payment_connected_'+'a'*32+'?sslmode=disable','issuer':'http://127.0.0.1:6544','api_port':6545,'web_port':6546}
            first=Deployment(Path(d),cfg)
            try:
                with self.assertRaisesRegex(ValueError,'another local operator'):Deployment(Path(d),cfg)
            finally:first.close()
if __name__=='__main__':unittest.main()
````

### FILE: `ci/web_local_stop.cjs`

```yaml
block_id: "PORTABLE-CI-GATE-RUNNER-LOCAL-DELIVERY:file8:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "4d6f649fafa5357d0ca7519488c5b830f2c44d9e97e1ed396d76ac3b11463155"
variables: []
secrets_allowed: false
```

````text
// AUTHORED local trusted-stdin shutdown bridge to Next's existing SIGTERM owner.
// The Windows Job Object still bounds the whole descendant tree on failure.
'use strict';
const path = require('node:path');
const fs = require('node:fs');
if (process.argv.length !== 3 || !path.isAbsolute(process.argv[2])) throw new Error('absolute server entry required');
const root = path.dirname(process.argv[2]);
const links = JSON.parse(fs.readFileSync(path.join(root, 'elite-runtime-links.json'), 'utf8'));
function local(rel) {
  if (typeof rel !== 'string' || !rel.startsWith('node_modules/') || rel.includes('\\') || rel.includes(':') || rel.split('/').some(p => !p || p === '.' || p === '..')) throw new Error('invalid runtime link');
  return path.join(root, ...rel.split('/'));
}
for (const [name, dest] of Object.entries(links)) {
  const link = local(name), target = local(dest);
  if (fs.existsSync(link) || !fs.statSync(target).isDirectory()) throw new Error('occupied link or absent target');
  // ZIPs preserve files, not the empty scope directories that precede links.
  fs.mkdirSync(path.dirname(link), { recursive: true });
  fs.symlinkSync(target, link, 'junction');
}
require(process.argv[2]);
let requested = false;
function stop() {
  if (requested) return;
  requested = true;
  process.emit('SIGTERM', 'SIGTERM');
}
let input = '';
process.stdin.setEncoding('utf8');
process.stdin.on('data', chunk => {
  input += chunk;
  if (input.length > 16) throw new Error('invalid local shutdown command');
  if (input === 'STOP\n') stop();
});
process.stdin.on('end', stop);
````

### FILE: `docs/LOCAL_REFERENCE_DELIVERY.md`

```yaml
block_id: "PORTABLE-CI-GATE-RUNNER-LOCAL-DELIVERY:file9:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "61e707b31d0c851d36e7233d43d5218defef775d9b9dc5006689c50f081aa96d"
variables: []
secrets_allowed: false
```

````markdown
# Entrega local de referencia

Scope: `LIBRARY_INFRASTRUCTURE / LOCAL_FIXTURES`, Windows x64. Estos comandos
son glue AUTHORED sobre los toolchains, contratos y owners ya fijados. No son
código de Microsoft, Google, xAI ni otra empresa. El artefacto incluye material
sintético público de sesión/Server Actions y rechaza activación productiva.

## Construcción

Materializar el perfil completo y los owners `PNPM-ARTIFACT-SELECTION-GATE` y
`WINDOWS-REFERENCE-TELEMETRY-RUNTIME` seleccionados por ese plan. Crear un JSON
`ruta relativa -> SHA-256` de las salidas exactas del compositor. El hash externo
de ese JSON es la autorización de la revisión; no se confía en un hash declarado
por el propio artefacto. Conservarlo con el receipt de composición.

`python ci/build_local_reference.py --help` enumera las rutas requeridas:
source, target ausente, inventory y su hash, Go, Node, caché de módulos y
install-inputs con su hash. `install-inputs` tiene exactamente projection,
contained, store, cache, acquisition_receipt y acquisition_sha256 del receipt
admitido de pnpm. Usa `pnpm_artifact_selection/local_runtime.py`: instalación
offline/frozen/copy, sin scripts, bajo Node jitless/no-addons. El payload local
adaptado de pnpm se adquiere/verifica fuera del artefacto y nunca se redistribuye
con él. Una caché ausente falla; no se descarga ni se promueve otra versión.

`--workspace` es una ruta local estable y ausente para el build. Next16.3.4 usa
rutas absolutas como identificadores RSC: los dos builds independientes usan esa
misma ruta lógica, cada vez con archivos nuevos e instalación nueva. Archivar el
workspace propio terminado antes de la segunda ejecución; el builder nunca lo
borra ni reutiliza un checkout/caché de compilación. El recibo conserva la ruta.
La cache de preview de Next recibe valores públicos de fixture fijados, y se
comprueba que el framework los consumió. El claim byte-idéntico es para esta
receta y toolchains del host, no para cualquier ruta o plataforma arbitraria.

Go se compila con trimpath, sin VCS ni build ID; se incluyen todos los módulos
locales de la composición. Next usa webpack, su salida standalone y el hash del
inventario como build ID. Sólo se normalizan sus localizadores de máquina en
server.js y required-server-files.json; los hashes previos quedan en el receipt.
Los arrays de archivos de los traces NFT se ordenan como metadata de packaging.
No se reescriben bundles de aplicación. El framework conserva enlaces a paquetes
en Windows: el artefacto almacena únicamente archivos normales y el mapa de
enlaces relativos; cada activación recrea junctions dentro de su copia privada.

El artefacto contiene API, Node fijado y su licencia original, frontend,
migraciones, fuentes correspondientes y los notices de los owners y dependencias
observadas. La igualdad de dos artefactos debe comprobarse sobre el inventario
completo de archivos. La firma y aceptación portable pertenecen a TEST07/T2810.

## Base vacía y reanudación

Preparar un PostgreSQL local admitido y una base de fixtures propia cuyo nombre
cumpla `elite_payment_connected_<32 hex>`. El runner exige loopback, usuario
postgres sin contraseña en el perfil de fixture, puerto explícito y sslmode
disable. Esto no prescribe autenticación para un despliegue productivo.

`python ci/migrate_local_reference.py --help` acepta el directorio migrations,
su inventario/hash, DSN local, psql/hash y directorio de receipt ausente. Aplica
las migraciones fijadas en orden con lock de sesión, transacción por migración y
ledger atómico de nombre/ordinal/hash. Reanudar no repite migraciones; un hash
distinto o una base preexistente sin baseline registrado falla cerrado. Una
migración fallida no autoriza tráfico y no ejecuta down migrations.

## Arranque y rollback

`ci/local_identity_fixture.cjs` es exclusivamente un emisor OIDC sintético
loopback para el ensayo. No sirve como IdP productivo. El runtime JSON local
contiene exactamente database_url, issuer, api_port y web_port. No acepta keys,
env arbitrario ni endpoints remotos. Pago/WhatsApp/IA/ARCA live no se habilitan.

`python ci/run_local_reference.py --help` recibe artefacto, inventario/hash,
runtime/hash, receipt/hash de migraciones y state-directory. Comprobará que la
base y las migraciones corresponden al artefacto. Arranca por un período finito
de hasta 480 segundos; conserva current.json y journals de activación. El owner
local exclusivo inicia API y frontend bajo el launcher de Job Objects existente.
El API adopta el mismo evento nativo ya probado por el worker de devolución.
Next conserva su propio cierre SIGTERM (exit 143); el wrapper termina con 0
sólo tras los cierres esperados y el supervisor comprueba árbol vacío.

`ci/local_release.py` expone Deployment.activate/recover/close para un operador
local. Verifica la revisión antes de detener la anterior; si falla la salud de
la candidata, reinicia el artefacto previo y conserva su puntero. Una activación
satisfactoria sustituye el puntero atómicamente. Los journals fallidos se
preservan. Recovery vuelve al último puntero confirmado; no adivina aceptación
desde un proceso vivo. Los handlers y PostgreSQL conservan idempotencia y datos.

Los probes locales prueban arranque/routing; los journeys completos conservan
sus suites por owner. La operación con carga/alertas/retención se cierra en
T2809. No se infiere HA, DR, SCM, aislamiento frente a otro usuario hostil ni
capacidad productiva por un Job Object o un health check.

## Otras plataformas

El blueprint selecciona esta referencia local. Contenedores/cloud/Kubernetes
siguen opt-in por target; el Dockerfile no aporta defaults de tags móviles y
copia el contexto para conservar los módulos locales. Sus imágenes, runtime,
TLS, permisos, registry, proveedor y ejecución requieren admisión propia. No
se etiqueta una ruta sin ejecutar como CONDITIONED sólo por credenciales.

Fuentes de configuración consultadas: [Next standalone](https://nextjs.org/docs/app/api-reference/config/next-config-js/output)
y [build ID](https://nextjs.org/docs/app/api-reference/config/next-config-js/generateBuildId).
La implementación consumida sigue fijada en Next16.3.4 y los locks vigentes;
consultar documentación actual no actualiza una dependencia automáticamente.

Firma local: docs/SIGNED_LOCAL_REFERENCE.md conecta este builder con el gate firmado. Cada invocación es limpia y archiva su workspace; la publicación verifica el recibo y el inventario real.
````


V402316: optional Windows local reference commands require PNPM-ARTIFACT-SELECTION-GATE0.9.0 and four exact WINDOW​​S-REFERENCE-TELEMETRY-RUNTIME0.1.0 supervisor files. Original portable runner unchanged. Read docs/LOCAL_REFERENCE_DELIVERY.md. No live production or corporate attribution.

V402 composed delta: V402317 connected local API/Next telemetry, current OIDC, fixed official middleware, finite supervised alert/fault/load/WAL recovery. Historical lock kept separate; docs/LOCAL_REFERENCE_OPERATIONS.md. No production admission.

### FILE: `ci/run_observed_reference.py`

```yaml
block_id: "PORTABLE-CI-GATE-RUNNER-LOCAL-OPERATIONS:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "ed307f01885921d0b19c0a2b9f0069661df27be6acdbb8cb1ffbd6b2e07e259b"
variables: []
secrets_allowed: false
```

````python
"""AUTHORED finite local API/Next -> official Prometheus -> fault/recovery probe.

Only a caller-selected, owned synthetic PostgreSQL reference is accepted. Every
binary, artifact, rule and configuration is bound to an external SHA. No live
provider, notification recipient, production identity or SLO is inferred.
"""
from pathlib import Path
import argparse,csv,hashlib,json,os,queue,re,socket,ssl,subprocess,threading,time,uuid
from concurrent.futures import ThreadPoolExecutor
from urllib.request import Request,build_opener,ProxyHandler,HTTPSHandler
from urllib.error import HTTPError
from urllib.parse import quote,urlsplit
from local_release import Deployment,verify_release,runtime_config
from build_local_reference import plain,read_json,require,sha,raw_json,environment

def main():
    p=argparse.ArgumentParser()
    for n in ['artifact','inventory','inventory-sha256','runtime-lock','runtime-lock-sha256','database-config','database-config-sha256','rules','rules-sha256','work']:
        p.add_argument('--'+n,required=True)
    a=p.parse_args();W=plain(a.work);require(not W.exists(),'work directory must be absent');W.mkdir()
    release=verify_release(a.artifact,a.inventory,a.inventory_sha256);A=Path(release['root']);source=Path(__file__).resolve().parent
    lock=read_json(a.runtime_lock,a.runtime_lock_sha256);bins={}
    for name in ['postgres','pg_ctl','psql','prometheus','promtool','telemetry-reference']:
        pin=lock['binaries'][name];v=plain(pin['path']);require(v.is_file()and v.stat().st_size==pin['bytes']and sha(v)==pin['sha256'],'runtime changed: '+name);bins[name]=v
    conf=read_json(a.database_config,a.database_config_sha256);require(set(conf)=={'database_url','data_path'},'database configuration fields')
    data=plain(conf['data_path']);u=urlsplit(conf['database_url']);db=u.path.removeprefix('/')
    runtime_config({'database_url':conf['database_url'],'issuer':'http://127.0.0.1:6544','api_port':6545,'web_port':6546})
    require(data.is_dir()and not(data/'postmaster.pid').exists(),'owned database must be stopped')
    rules=plain(a.rules);require(sha(rules)==a.rules_sha256,'rules changed')
    system=Path(os.environ['SYSTEMROOT'])/'System32'
    sid=list(csv.reader(subprocess.check_output([str(system/'whoami.exe'),'/user','/fo','csv','/nh'],text=True).splitlines()))[0][1]
    require(re.fullmatch(r'S-[0-9-]+',sid),'local identity unresolved')
    subprocess.run([str(system/'icacls.exe'),str(W),'/inheritance:r','/grant:r','*'+sid+':(OI)(CI)F','*S-1-5-18:(OI)(CI)F'],check=True,stdout=subprocess.DEVNULL,creationflags=subprocess.CREATE_NO_WINDOW)
    env=environment(A/'node.exe',W);env.update(PGCONNECT_TIMEOUT='3',COMSPEC=str(system/'cmd.exe'))
    procs={};handles=[];dep=None;renamed=False;events=[];start=time.monotonic();deadline=start+1170
    result={'state':'FAIL','scope':'LOCAL_FIXTURES','accounts_used':False,'live_effects':False,'requests':{},'source_sha256':release['source_sha256'],'host_sha256':sha(A/'api/electromobility-api.exe'),'rules_sha256':a.rules_sha256,'runtime_lock_sha256':a.runtime_lock_sha256}
    def log(kind,**values):
        events.append({'seconds':round(time.monotonic()-start,3),'kind':kind,**values});(W/'progress.json').write_bytes(raw_json(events));print(json.dumps(events[-1]),flush=True)
    def command(name,args,seconds=30):
        q=subprocess.run(list(map(str,args)),env=env,stdin=subprocess.DEVNULL,capture_output=True,timeout=seconds,creationflags=subprocess.CREATE_NO_WINDOW)
        (W/(name+'.log')).write_bytes(q.stdout+q.stderr);require(q.returncode==0,'reference command failed: '+name);return q.stdout.decode().strip()
    def sql(query,database=db):
        return command('sql-'+uuid.uuid4().hex,[bins['psql'],'-X','-w','-qAt','-h','127.0.0.1','-p',u.port,'-U','postgres','-d',database,'-v','ON_ERROR_STOP=1','-c',query],15)
    def launch(name,args):
        f=(W/(name+'-'+uuid.uuid4().hex+'.log')).open('xb');handles.append(f)
        q=subprocess.Popen(list(map(str,args)),cwd=W,env=env,stdin=subprocess.DEVNULL,stdout=f,stderr=subprocess.STDOUT,creationflags=subprocess.CREATE_NO_WINDOW);procs[name]=q;return q
    def wait(fn,label,seconds=30):
        end=min(time.monotonic()+seconds,deadline)
        while time.monotonic()<end:
            try:
                value=fn()
                if value:return value
            except (OSError,ValueError):pass
            time.sleep(.25)
        raise TimeoutError(label)
    def port():
        with socket.socket()as s:s.bind(('127.0.0.1',0));return s.getsockname()[1]
    plain_client=build_opener(ProxyHandler({}))
    def get(url,client=plain_client):
        with client.open(url,timeout=5)as r:
            raw=r.read(2_000_001);require(len(raw)<=2_000_000,'response budget exceeded');return raw
    try:
        launch('postgres',[bins['postgres'],'-D',data])
        def pg_ready():
            require(procs['postgres'].poll()is None,'owned PostgreSQL exited')
            return Path(sql('show data_directory','postgres')).resolve()==data.resolve()
        wait(pg_ready,'owned PostgreSQL')
        require(sql('select count(*) from library_delivery.migration')=='84','current migration ledger required')
        sql("insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values('b0000000-0000-4000-8000-000000000001','delivery-local','Delivery Fixture','Delivery Fixture')on conflict do nothing; insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values('b0000000-0000-4000-8000-000000000001','b0000000-0000-4000-8000-000000000002','reference','Reference','store')on conflict do nothing")
        f=(W/'fixture-stderr.log').open('xb');handles.append(f)
        fixture=subprocess.Popen([A/'node.exe',source/'local_identity_fixture.cjs','--ttl-seconds','1200'],cwd=W,env=env,stdin=subprocess.PIPE,stdout=subprocess.PIPE,stderr=f,creationflags=subprocess.CREATE_NO_WINDOW);procs['fixture']=fixture
        ready=queue.Queue();threading.Thread(target=lambda:ready.put(fixture.stdout.readline()),daemon=True).start();issuer=json.loads(ready.get(timeout=10))['issuer']
        cfg={'database_url':conf['database_url'],'issuer':issuer,'api_port':port(),'web_port':port(),'metrics_port':port()};runtime_config(cfg)
        dep=Deployment(W/'host',cfg);result['activation']=dep.activate(A,a.inventory,a.inventory_sha256)
        key='local-observability-'+uuid.uuid4().hex;token=''
        body=raw_json({'OrganizationID':'b0000000-0000-4000-8000-000000000002','Currency':'ARS','TotalMinorUnits':1900})
        def refresh():
            nonlocal token
            token=json.loads(get(issuer+'/fixture-token'))['token']
        def request(k=key,credential=None,payload=body):
            req=Request(f"http://127.0.0.1:{cfg['api_port']}/v1/orders?contact=PRIVATE_QUERY@example.invalid",data=payload,
                headers={'Authorization':'Bearer '+(token if credential is None else credential),'Content-Type':'application/json','Idempotency-Key':k,'Cookie':'PRIVATE_COOKIE','X-Private':'PRIVATE_HEADER'},method='POST')
            then=time.perf_counter()
            try:
                with build_opener(ProxyHandler({})).open(req,timeout=8)as r:status=r.status;raw=r.read(16385)
            except HTTPError as e:status=e.code;raw=e.read(16385)
            require(len(raw)<=16384 and b'PRIVATE_'not in raw,'response privacy/budget');return status,raw,time.perf_counter()-then
        def record(r):
            status,raw,seconds=r;result['requests'][str(status)]=result['requests'].get(str(status),0)+1;return r
        refresh();first=record(request());require(first[0]==201,'first order');order=json.loads(first[1])
        again=record(request());require(again[0]==200 and json.loads(again[1])==order,'order replay')
        require(record(request(key+'-auth','PRIVATE_TOKEN'))[0]==401,'identity fail-closed')
        require(record(request(key+'-body',payload=b'{"PRIVATE_BODY":true}'))[0]==400,'invalid body')
        def metrics(label):
            raw=get(f"http://127.0.0.1:{cfg['metrics_port']}/metrics")
            for forbidden in [b'PRIVATE_',b'b0000000',b'local-delivery-customer',b'customer_order',b'http_route',b'server_address',b'url_',b'target_info',b'otel_scope',token.encode()]:require(forbidden not in raw,'metric attribute escaped closed profile')
            (W/(label+'.prom')).write_bytes(raw);return raw
        metrics('healthy');command('certificates',[bins['telemetry-reference'],'--certificates',W])
        tls=ssl.create_default_context(cafile=str(W/'ca.pem'));tls.minimum_version=ssl.TLSVersion.TLSv1_3;tls.load_cert_chain(str(W/'client.pem'),str(W/'client.key'))
        client=build_opener(ProxyHandler({}),HTTPSHandler(context=tls));promport=port()
        def api(path):return json.loads(get('https://127.0.0.1:'+str(promport)+path,client))
        def query(expr,at=None):return api('/api/v1/query?query='+quote(expr)+(''if at is None else '&time='+str(at)))
        def alerts():return [v for v in api('/api/v1/alerts')['data']['alerts']if v['labels']['alertname']=='PlatformHighErrorRate']
        (W/'prometheus.json').write_bytes(raw_json({'global':{'scrape_interval':'1s','evaluation_interval':'30s'},'rule_files':[str(rules)],'scrape_configs':[{'job_name':'reference-http','static_configs':[{'targets':[f"127.0.0.1:{cfg['metrics_port']}"]}]}]}))
        (W/'web.json').write_bytes(raw_json({'tls_server_config':{'cert_file':str(W/'server.pem'),'key_file':str(W/'server.key'),'client_ca_file':str(W/'ca.pem'),'client_auth_type':'RequireAndVerifyClientCert','min_version':'TLS13'}}))
        command('rules-check',[bins['promtool'],'check','config',W/'prometheus.json'])
        promargs=[bins['prometheus'],'--config.file='+str(W/'prometheus.json'),'--web.listen-address=127.0.0.1:'+str(promport),'--web.config.file='+str(W/'web.json'),'--storage.tsdb.path='+str(W/'tsdb'),'--storage.tsdb.retention.time=1h','--storage.tsdb.retention.size=64MB','--log.level=error']
        launch('prometheus',promargs);wait(lambda:query('up{job="reference-http"}')['data']['result'],'first scrape');time.sleep(3)
        sql('alter table sales.customer_order rename to customer_order_reference_fault');renamed=True
        pending=fired=None;last=None;count=0;next_error=time.monotonic();log('fault_injected',rule_hold_seconds=600)
        while time.monotonic()<deadline-130:
            require(dep.instance.thread.is_alive()and procs['prometheus'].poll()is None,'observed host stopped')
            if time.monotonic()>=next_error:
                refresh();require(record(request(key+'-failure-'+str(count)))[0]==500,'actual database fault did not fail');count+=1;next_error+=240;metrics('fault-'+str(count));log('actual_http_500',count=count)
            active=alerts();state=active[0]['state']if active else 'inactive'
            if state!=last:log('alert_state',state=state);last=state
            if state=='pending'and pending is None:pending=time.monotonic()
            if state=='firing':
                fired=time.monotonic();require(pending is not None and fired-pending>=590,'hold shortened');result['firing_alert']=active;break
            time.sleep(1)
        require(fired is not None,'alert did not fire within finite budget')
        sql('alter table sales.customer_order_reference_fault rename to customer_order');renamed=False;refresh()
        started=time.monotonic()
        with ThreadPoolExecutor(max_workers=8)as pool:load=list(pool.map(lambda _:request(),range(1000)))
        for r in load:record(r);require(r[0]==200 and json.loads(r[1])==order,'concurrent replay diverged')
        timings=sorted(r[2]for r in load);elapsed=time.monotonic()-started;require(elapsed<90,'load exceeded finite local budget')
        result['load']={'requests':1000,'concurrency':8,'seconds':round(elapsed,3),'p95_ms':round(timings[949]*1000,3),'max_ms':round(timings[-1]*1000,3),'all_same_order':True,'production_slo_claim':False}
        wait(lambda:not any(x['state']=='firing'for x in alerts()),'alert recovery',70);metrics('recovered')
        fixed=time.time()-5;expr='sum(http_server_request_duration_seconds_count)';before=query(expr,fixed)
        procs['prometheus'].terminate();procs['prometheus'].wait(10);launch('prometheus',promargs)
        after=wait(lambda:query(expr,fixed),'Prometheus WAL restart');require(before['data']['result']==after['data']['result']and before['data']['result'],'retained series changed')
        result['retention']={'restart_fixed_timestamp_query_identical':True,'time':'1h','block_size':'64MB','disk_hard_cap_claim':False,'total_tsdb_bytes':sum(v.stat().st_size for v in(W/'tsdb').rglob('*')if v.is_file())}
        result['host_recovery']=dep.recover();refresh();replay=record(request());require(replay[0]==200 and json.loads(replay[1])==order,'host recovery replay')
        count=sql("select count(*) from sales.customer_order where tenant_id='b0000000-0000-4000-8000-000000000001'and order_id='"+order['id']+"'");require(count=='1','duplicate durable order')
        result['stop']=dep.close();dep=None;require(result['stop']['tree_empty']and result['stop']['exit_code']==0,'supervised stop failed')
        result.update(state='PASS',durable_order_count=1,pending_to_firing_seconds=round(fired-pending,3));log('recovered',load=result['load'],order_count=1)
    except BaseException as e:
        result.update(error_type=type(e).__name__,error=str(e));raise
    finally:
        if renamed:
            try:sql('alter table sales.customer_order_reference_fault rename to customer_order');result['fault_restored_after_failure']=True
            except BaseException as e:result['fault_restore_failure']=str(e)
        if dep:
            try:result['failure_stop']=dep.close()
            except BaseException as e:result['failure_stop_error']=str(e)
        for name in ['fixture','prometheus','postgres']:
            q=procs.get(name)
            if q is None:continue
            if q.poll()is None:
                try:
                    if name=='fixture':q.stdin.write(b'STOP\n');q.stdin.flush();q.wait(5)
                    elif name=='postgres':command('postgres-stop',[bins['pg_ctl'],'-D',data,'-m','fast','-w','stop']);q.wait(10)
                    else:q.terminate();q.wait(10)
                except BaseException:q.terminate();q.wait(5)
            if name=='fixture':q.stdin.close();q.stdout.close()
        for f in handles:f.close()
        result['children_exited']=all(q.poll()is not None for q in procs.values());result['seconds']=round(time.monotonic()-start,3)
        (W/'result.json').write_bytes(raw_json(result));print('OBSERVED_REFERENCE '+json.dumps(result),flush=True)
    return 0
if __name__=='__main__':raise SystemExit(main())
````

### FILE: `docs/LOCAL_REFERENCE_OPERATIONS.md`

```yaml
block_id: "PORTABLE-CI-GATE-RUNNER-LOCAL-OPERATIONS:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "896b718fc54f6c933fafaf78be73ad8bab46de967067625e024cb4a680e353ce"
variables: []
secrets_allowed: false
```

````markdown
# Operación de la referencia local

AUTHORED glue alrededor de Go/OTel/Prometheus y del supervisor Windows existentes.
Scope LIBRARY_INFRASTRUCTURE_LOCAL_FIXTURES, un operador local de confianza.

## Arranque y ensayo

Primero construir la referencia y aplicar ci/migrate_local_reference.py en una
base vacía propia. El contrato del ensayo espera las 84 migraciones de esta
composición, no una base productiva. Su configuración JSON contiene solamente
database_url del patrón elite_payment_connected_<32hex>, loopback sin contraseña,
y data_path de ese clúster detenido. No acepta cuentas ni proveedores live.

Ejecutar con Python fijado:

```text
python -X utf8 -B ci/run_observed_reference.py --artifact <artefacto> --inventory <inventario> --inventory-sha256 <sha> --runtime-lock <lock> --runtime-lock-sha256 <sha> --database-config <fixture.json> --database-config-sha256 <sha> --rules ops/prometheus/platform.rules.yml --rules-sha256 <sha> --work <directorio-ausente>
```

El lock externo fija postgres/pg_ctl/psql, prometheus/promtool y
telemetry-reference. Sus identidades/licencias se conservan en los expedientes
V372/V375; nunca se descargan versiones móviles desde este comando. El artefacto
fija API, Node y Next. La identidad RSA efímera sólo vive en el fixture loopback.

El ensayo arranca API+Next bajo el Job Object ya admitido: doce procesos como
máximo, 1 GiB de commit, salida capturada de 2 MiB, 1200 segundos con métricas y
25 segundos de cierre. Fuera de este ensayo, el controlador conserva 600 segundos
y el entrypoint de uso local sigue limitado a 480. Un puerto de métricas opcional
debe ser entero, loopback y distinto de los puertos API/web. No hay SCM ni HA.

## Señal, alerta y recuperación

Sólo se exporta el histograma HTTP de método/status, con cardinalidad limitada y
sin query, ruta, tenant, tokens, cookies, trazas ni exemplars. Se prueba contra
el handler real, su OIDC, PostgreSQL y un pedido durable. Prometheus consulta con
mTLS efímero y retiene WAL/series entre reinicios. La regla HTTP conserva ventana
5m, umbral 1%, hold 10m y evaluación 30s. No se inyectan muestras ni se acelera
el reloj. El propio ensayo renombra y restaura exclusivamente la tabla de su
base sintética; no copiar esa operación a una base real.

Ante PlatformHighErrorRate: consultar progress.json/result.json, comprobar salud
del API y estado del clúster propietario; restaurar la tabla sólo si el receipt
declara fault_restore_failure. El finally intenta restaurar antes del apagado.
En operación fuera del ensayo, no modificar tablas: recuperar la última release
confirmada con Deployment.recover y conservar el fallo. Los runbook_url históricos
replace.invalid no son enlaces de operación: esta guía gobierna el caso local.
OutboxBacklog y BackupStale no tienen señales conectadas por este cambio y no se
presentan como alertas probadas del host.

Después del firing real: 1000 replay requests, concurrencia ocho, identidad y
pedido constantes; se observan p95/máximo y se exige la misma fila durable.
La reparación limpia la alerta; el reinicio de Prometheus conserva una consulta
con timestamp fijo y el recovery del controlador conserva el pedido. Cada
instancia termina con recibo del árbol vacío, Go0/NextSIGTERM143/wrapper0.

## Retención y límites del claim

Prometheus: 1h y 64MB para bloques persistentes. Eso NO es una cuota dura de disco:
WAL, head y compacción pueden excederla. La documentación oficial lo especifica:
https://prometheus.io/docs/prometheus/latest/storage/ . Se registra tamaño real y
replay de WAL; no se infiere limpieza de bloques expirados de una ejecución de
menos de una hora. La captura del proceso sí tiene un límite duro de 2 MiB.

El owner WINDOWS_REFERENCE_TELEMETRY_RUNTIME conserva el ensayo separado V372:
rotación de archivo del Collector a 1MiB, dos backups, un día, 2000 observaciones
y reinicio. Se conservan código/hashes/binarios de esa evidencia, sin atribuirle
logs del API nuevo. El despliegue productivo permanente, routing de notificaciones,
respondedor humano, SLOs comerciales, cuotas target, PITR y fallos de host quedan
CANDIDATE_TARGET, no CONDITIONED sólo por una credencial. El límite de T2809
admitido aquí es el host y ensayo local finito explícito.
````


V402317: current host profile is documented in LOCAL_REFERENCE_OPERATIONS.md; source-lock separates historical synthetic principal from current OIDC host. Narrow local proof only.

V402 composed delta: V402322 exact runtime notice collection;129traced manifests, original source/legal texts and declared notice renderings. Build-only inputs remain distinct and unshipped. No dependency/runtime version or business code change.

### FILE: `licenses/reference/next-MIT.txt`

```yaml
block_id: "REFERENCE-NOTICE-BLOCKS:file1:v1"
operation: CREATE
provenance: VERBATIM
source: "Exact official provenance and declared rendering in licenses/reference/SOURCE_NOTICES.json sources[licenses/reference/next-MIT.txt], lock SHA256 18b20c9a287eb5ad12217db27f58384812d2d36a7a74c8355390db94f9bdd47b; payload SHA256 ee765244e2d59f5234d474f62e0766fa0c8b99af967fdd4c0cb8dcb0c76ea224"
license: "MIT"
sha256: "ee765244e2d59f5234d474f62e0766fa0c8b99af967fdd4c0cb8dcb0c76ea224"
variables: []
secrets_allowed: false
```

````text
The MIT License (MIT)

Copyright (c) 2025 Vercel, Inc.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
````

### FILE: `licenses/reference/constants-browserify-README.md.txt`

```yaml
block_id: "REFERENCE-NOTICE-BLOCKS:file2:v1"
operation: CREATE
provenance: VERBATIM
source: "Exact official provenance and declared rendering in licenses/reference/SOURCE_NOTICES.json sources[licenses/reference/constants-browserify-README.md.txt], lock SHA256 18b20c9a287eb5ad12217db27f58384812d2d36a7a74c8355390db94f9bdd47b; payload SHA256 b06569bf170a015c56cbad9b45756519d0b8a682770fe30d66bc1e235e26196a"
license: "MIT"
sha256: "b06569bf170a015c56cbad9b45756519d0b8a682770fe30d66bc1e235e26196a"
variables: []
secrets_allowed: false
```

````text

# constants-browserify

Node's `constants` module for the browser.

[![downloads](https://img.shields.io/npm/dm/constants-browserify.svg)](https://www.npmjs.org/package/constants-browserify)

## Usage

To use with browserify cli:

```bash
$ browserify -r constants:constants-browserify script.js
```

To use with browserify api:

```js
browserify()
  .require('constants-browserify', { expose: 'constants' })
  .add(__dirname + '/script.js')
  .bundle()
  // ...
```

## Installation

With [npm](http://npmjs.org) do

```bash
$ npm install constants-browserify
```

## License

Copyright (c) 2013 Julian Gruber &lt;julian@juliangruber.com&gt;

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
````

### FILE: `licenses/reference/data-uri-to-buffer-README.md.txt`

```yaml
block_id: "REFERENCE-NOTICE-BLOCKS:file3:v1"
operation: CREATE
provenance: VERBATIM
source: "Exact official provenance and declared rendering in licenses/reference/SOURCE_NOTICES.json sources[licenses/reference/data-uri-to-buffer-README.md.txt], lock SHA256 18b20c9a287eb5ad12217db27f58384812d2d36a7a74c8355390db94f9bdd47b; payload SHA256 6a14de4d552c0c88bdbbaf57f72f625835685eea97de5047f87a463089b8e9f6"
license: "MIT"
sha256: "6a14de4d552c0c88bdbbaf57f72f625835685eea97de5047f87a463089b8e9f6"
variables: []
secrets_allowed: false
```

````text
data-uri-to-buffer
==================
### Generate a Buffer instance from a [Data URI][rfc] string
[![Build Status](https://travis-ci.org/TooTallNate/node-data-uri-to-buffer.svg?branch=master)](https://travis-ci.org/TooTallNate/node-data-uri-to-buffer)

This module accepts a ["data" URI][rfc] String of data, and returns a
node.js `Buffer` instance with the decoded data.


Installation
------------

Install with `npm`:

``` bash
$ npm install data-uri-to-buffer
```


Example
-------

``` js
var dataUriToBuffer = require('data-uri-to-buffer');

// plain-text data is supported
var uri = 'data:,Hello%2C%20World!';
var decoded = dataUriToBuffer(uri);
console.log(decoded.toString());
// 'Hello, World!'

// base64-encoded data is supported
uri = 'data:text/plain;base64,SGVsbG8sIFdvcmxkIQ%3D%3D';
decoded = dataUriToBuffer(uri);
console.log(decoded.toString());
// 'Hello, World!'
```


API
---

### dataUriToBuffer(String uri) → Buffer

The `type` property on the Buffer instance gets set to the main type portion of
the "mediatype" portion of the "data" URI, or defaults to `"text/plain"` if not
specified.

The `typeFull` property on the Buffer instance gets set to the entire
"mediatype" portion of the "data" URI (including all parameters), or defaults
to `"text/plain;charset=US-ASCII"` if not specified.

The `charset` property on the Buffer instance gets set to the Charset portion of
the "mediatype" portion of the "data" URI, or defaults to `"US-ASCII"` if the
entire type is not specified, or defaults to `""` otherwise.

*Note*: If the only the main type is specified but not the charset, e.g.
`"data:text/plain,abc"`, the charset is set to the empty string. The spec only
defaults to US-ASCII as charset if the entire type is not specified.


License
-------

(The MIT License)

Copyright (c) 2014 Nathan Rajlich &lt;nathan@tootallnate.net&gt;

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
'Software'), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

[rfc]: http://tools.ietf.org/html/rfc2397
````

### FILE: `licenses/reference/http-proxy-agent-README.md.txt`

```yaml
block_id: "REFERENCE-NOTICE-BLOCKS:file4:v1"
operation: CREATE
provenance: VERBATIM
source: "Exact official provenance and declared rendering in licenses/reference/SOURCE_NOTICES.json sources[licenses/reference/http-proxy-agent-README.md.txt], lock SHA256 18b20c9a287eb5ad12217db27f58384812d2d36a7a74c8355390db94f9bdd47b; payload SHA256 ca0b65367e78e9c255a1c9883eb7ec30aae491d76e1c6815b865563106d21dff"
license: "MIT"
sha256: "ca0b65367e78e9c255a1c9883eb7ec30aae491d76e1c6815b865563106d21dff"
variables: []
secrets_allowed: false
```

````text
http-proxy-agent
================
### An HTTP(s) proxy `http.Agent` implementation for HTTP
[![Build Status](https://github.com/TooTallNate/node-http-proxy-agent/workflows/Node%20CI/badge.svg)](https://github.com/TooTallNate/node-http-proxy-agent/actions?workflow=Node+CI)

This module provides an `http.Agent` implementation that connects to a specified
HTTP or HTTPS proxy server, and can be used with the built-in `http` module.

__Note:__ For HTTP proxy usage with the `https` module, check out
[`node-https-proxy-agent`](https://github.com/TooTallNate/node-https-proxy-agent).

Installation
------------

Install with `npm`:

``` bash
$ npm install http-proxy-agent
```


Example
-------

``` js
var url = require('url');
var http = require('http');
var HttpProxyAgent = require('http-proxy-agent');

// HTTP/HTTPS proxy to connect to
var proxy = process.env.http_proxy || 'http://168.63.76.32:3128';
console.log('using proxy server %j', proxy);

// HTTP endpoint for the proxy to connect to
var endpoint = process.argv[2] || 'http://nodejs.org/api/';
console.log('attempting to GET %j', endpoint);
var opts = url.parse(endpoint);

// create an instance of the `HttpProxyAgent` class with the proxy server information
var agent = new HttpProxyAgent(proxy);
opts.agent = agent;

http.get(opts, function (res) {
  console.log('"response" event!', res.headers);
  res.pipe(process.stdout);
});
```


License
-------

(The MIT License)

Copyright (c) 2013 Nathan Rajlich &lt;nathan@tootallnate.net&gt;

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
'Software'), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
````

### FILE: `licenses/reference/https-proxy-agent-README.md.txt`

```yaml
block_id: "REFERENCE-NOTICE-BLOCKS:file5:v1"
operation: CREATE
provenance: VERBATIM
source: "Exact official provenance and declared rendering in licenses/reference/SOURCE_NOTICES.json sources[licenses/reference/https-proxy-agent-README.md.txt], lock SHA256 18b20c9a287eb5ad12217db27f58384812d2d36a7a74c8355390db94f9bdd47b; payload SHA256 32f0856d2c43df7d05cca960fdee84e1e38ab545bd7b2186433dfa41aa90a712"
license: "MIT"
sha256: "32f0856d2c43df7d05cca960fdee84e1e38ab545bd7b2186433dfa41aa90a712"
variables: []
secrets_allowed: false
```

````text
https-proxy-agent
================
### An HTTP(s) proxy `http.Agent` implementation for HTTPS
[![Build Status](https://github.com/TooTallNate/node-https-proxy-agent/workflows/Node%20CI/badge.svg)](https://github.com/TooTallNate/node-https-proxy-agent/actions?workflow=Node+CI)

This module provides an `http.Agent` implementation that connects to a specified
HTTP or HTTPS proxy server, and can be used with the built-in `https` module.

Specifically, this `Agent` implementation connects to an intermediary "proxy"
server and issues the [CONNECT HTTP method][CONNECT], which tells the proxy to
open a direct TCP connection to the destination server.

Since this agent implements the CONNECT HTTP method, it also works with other
protocols that use this method when connecting over proxies (i.e. WebSockets).
See the "Examples" section below for more.


Installation
------------

Install with `npm`:

``` bash
$ npm install https-proxy-agent
```


Examples
--------

#### `https` module example

``` js
var url = require('url');
var https = require('https');
var HttpsProxyAgent = require('https-proxy-agent');

// HTTP/HTTPS proxy to connect to
var proxy = process.env.http_proxy || 'http://168.63.76.32:3128';
console.log('using proxy server %j', proxy);

// HTTPS endpoint for the proxy to connect to
var endpoint = process.argv[2] || 'https://graph.facebook.com/tootallnate';
console.log('attempting to GET %j', endpoint);
var options = url.parse(endpoint);

// create an instance of the `HttpsProxyAgent` class with the proxy server information
var agent = new HttpsProxyAgent(proxy);
options.agent = agent;

https.get(options, function (res) {
  console.log('"response" event!', res.headers);
  res.pipe(process.stdout);
});
```

#### `ws` WebSocket connection example

``` js
var url = require('url');
var WebSocket = require('ws');
var HttpsProxyAgent = require('https-proxy-agent');

// HTTP/HTTPS proxy to connect to
var proxy = process.env.http_proxy || 'http://168.63.76.32:3128';
console.log('using proxy server %j', proxy);

// WebSocket endpoint for the proxy to connect to
var endpoint = process.argv[2] || 'ws://echo.websocket.org';
var parsed = url.parse(endpoint);
console.log('attempting to connect to WebSocket %j', endpoint);

// create an instance of the `HttpsProxyAgent` class with the proxy server information
var options = url.parse(proxy);

var agent = new HttpsProxyAgent(options);

// finally, initiate the WebSocket connection
var socket = new WebSocket(endpoint, { agent: agent });

socket.on('open', function () {
  console.log('"open" event!');
  socket.send('hello world');
});

socket.on('message', function (data, flags) {
  console.log('"message" event! %j %j', data, flags);
  socket.close();
});
```

API
---

### new HttpsProxyAgent(Object options)

The `HttpsProxyAgent` class implements an `http.Agent` subclass that connects
to the specified "HTTP(s) proxy server" in order to proxy HTTPS and/or WebSocket
requests. This is achieved by using the [HTTP `CONNECT` method][CONNECT].

The `options` argument may either be a string URI of the proxy server to use, or an
"options" object with more specific properties:

  * `host` - String - Proxy host to connect to (may use `hostname` as well). Required.
  * `port` - Number - Proxy port to connect to. Required.
  * `protocol` - String - If `https:`, then use TLS to connect to the proxy.
  * `headers` - Object - Additional HTTP headers to be sent on the HTTP CONNECT method.
  * Any other options given are passed to the `net.connect()`/`tls.connect()` functions.


License
-------

(The MIT License)

Copyright (c) 2013 Nathan Rajlich &lt;nathan@tootallnate.net&gt;

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
'Software'), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

[CONNECT]: http://en.wikipedia.org/wiki/HTTP_tunnel#HTTP_CONNECT_Tunneling
````

### FILE: `licenses/reference/punycode-LICENSE-MIT.txt.txt`

```yaml
block_id: "REFERENCE-NOTICE-BLOCKS:file6:v1"
operation: CREATE
provenance: VERBATIM
source: "Exact official provenance and declared rendering in licenses/reference/SOURCE_NOTICES.json sources[licenses/reference/punycode-LICENSE-MIT.txt.txt], lock SHA256 18b20c9a287eb5ad12217db27f58384812d2d36a7a74c8355390db94f9bdd47b; payload SHA256 483acb265f182907d1caf6cff9c16c96f31325ed23792832cc5d8b12d5f88c8a"
license: "MIT"
sha256: "483acb265f182907d1caf6cff9c16c96f31325ed23792832cc5d8b12d5f88c8a"
variables: []
secrets_allowed: false
```

````text
Copyright Mathias Bynens <https://mathiasbynens.be/>

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
````

### FILE: `licenses/reference/setimmediate-LICENSE.txt.txt`

```yaml
block_id: "REFERENCE-NOTICE-BLOCKS:file7:v1"
operation: CREATE
provenance: VERBATIM
source: "Exact official provenance and declared rendering in licenses/reference/SOURCE_NOTICES.json sources[licenses/reference/setimmediate-LICENSE.txt.txt], lock SHA256 18b20c9a287eb5ad12217db27f58384812d2d36a7a74c8355390db94f9bdd47b; payload SHA256 c4b4ad3a5746f1f5249a6dd90396ec519264e1bb02e01e48a6522c48a3a97cb4"
license: "MIT"
sha256: "c4b4ad3a5746f1f5249a6dd90396ec519264e1bb02e01e48a6522c48a3a97cb4"
variables: []
secrets_allowed: false
```

````text
Copyright (c) 2012 Barnesandnoble.com, llc, Donavon West, and Domenic Denicola

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
````

### FILE: `licenses/reference/string-hash-README.md.txt`

```yaml
block_id: "REFERENCE-NOTICE-BLOCKS:file8:v1"
operation: CREATE
provenance: VERBATIM
source: "Exact official provenance and declared rendering in licenses/reference/SOURCE_NOTICES.json sources[licenses/reference/string-hash-README.md.txt], lock SHA256 18b20c9a287eb5ad12217db27f58384812d2d36a7a74c8355390db94f9bdd47b; payload SHA256 62a7a3b8dba4cc818974ab2c00cb7fb7cb20afd7d3c73ed049d7f2eb5621a46e"
license: "CC0-1.0"
sha256: "62a7a3b8dba4cc818974ab2c00cb7fb7cb20afd7d3c73ed049d7f2eb5621a46e"
variables: []
secrets_allowed: false
```

````text
string-hash
===========

A fast string hashing function for Node.JS. The particular algorithm is quite
similar to `djb2`, by Dan Bernstein and available
[here](http://www.cse.yorku.ca/~oz/hash.html). Differences include iterating
over the string *backwards* (as that is faster in JavaScript) and using the XOR
operator instead of the addition operator (as described at that page and
because it obviates the need for modular arithmetic in JavaScript).

The hashing function returns a number between 0 and 4294967295 (inclusive).

Thanks to [cscott](https://github.com/cscott) for reminding us how integers
work in JavaScript.

License
-------

To the extend possible by law, The Dark Sky Company, LLC has [waived all
copyright and related or neighboring rights][cc0] to this library.

[cc0]: http://creativecommons.org/publicdomain/zero/1.0/
````

### FILE: `licenses/reference/unistore-README.md.txt`

```yaml
block_id: "REFERENCE-NOTICE-BLOCKS:file9:v1"
operation: CREATE
provenance: VERBATIM
source: "Exact official provenance and declared rendering in licenses/reference/SOURCE_NOTICES.json sources[licenses/reference/unistore-README.md.txt], lock SHA256 18b20c9a287eb5ad12217db27f58384812d2d36a7a74c8355390db94f9bdd47b; payload SHA256 9e24224113e53a805fe76ed8084e13697a070061191f627c2053c0ea67f5c364"
license: "MIT"
sha256: "9e24224113e53a805fe76ed8084e13697a070061191f627c2053c0ea67f5c364"
variables: []
secrets_allowed: false
```

````text
<p align="center">
  <img src="https://i.imgur.com/o0u6dto.png" width="300" height="300" alt="unistore">
  <br>
  <a href="https://www.npmjs.org/package/unistore"><img src="https://img.shields.io/npm/v/unistore.svg?style=flat" alt="npm"></a> <a href="https://travis-ci.org/developit/unistore"><img src="https://travis-ci.org/developit/unistore.svg?branch=master" alt="travis"></a>
</p>

# unistore

> A tiny 350b centralized state container with component bindings for [Preact] & [React].

- **Small** footprint complements Preact nicely _(unistore + unistore/preact is ~650b)_
- **Familiar** names and ideas from Redux-like libraries
- **Useful** data selectors to extract properties from state
- **Portable** actions can be moved into a common place and imported
- **Functional** actions are just reducers
- **NEW**: seamlessly run Unistore in a worker via [Stockroom](https://github.com/developit/stockroom)

## Table of Contents

- [Install](#install)
- [Usage](#usage)
- [Examples](#examples)
- [API](#api)
- [License](#license)

## Install

This project uses [node](http://nodejs.org) and [npm](https://npmjs.com). Go check them out if you don't have them locally installed.

```sh
npm install --save unistore
```

Then with a module bundler like [webpack](https://webpack.js.org) or [rollup](http://rollupjs.org), use as you would anything else:

```js
// The store:
import createStore from 'unistore'

// Preact integration
import { Provider, connect } from 'unistore/preact'

// React integration
import { Provider, connect } from 'unistore/react'
```

Alternatively, you can import the "full" build for each, which includes both `createStore` and the integration for your library of choice:

```js
import { createStore, Provider, connect } from 'unistore/full/preact'
```

The [UMD](https://github.com/umdjs/umd) build is also available on [unpkg](https://unpkg.com):

```html
<!-- just unistore(): -->
<script src="https://unpkg.com/unistore/dist/unistore.umd.js"></script>
<!-- for preact -->
<script src="https://unpkg.com/unistore/full/preact.umd.js"></script>
<!-- for react -->
<script src="https://unpkg.com/unistore/full/react.umd.js"></script>
```

You can find the library on `window.unistore`.

### Usage

```js
import createStore from 'unistore'
import { Provider, connect } from 'unistore/preact'

let store = createStore({ count: 0 })

// If actions is a function, it gets passed the store:
let actions = store => ({
  // Actions can just return a state update:
  increment(state) {
    return { count: state.count+1 }
  },

  // The above example as an Arrow Function:
  increment2: ({ count }) => ({ count: count+1 }),

  //Actions receive current state as first parameter and any other params next
  //check this function as <button onClick={incrementAndLog}>
  incrementAndLog: ({ count }, event) => {
    console.info(event)
    return { count: count+1 }
  },

  // Async actions can be pure async/promise functions:
  async getStuff(state) {
    let res = await fetch('/foo.json')
    return { stuff: await res.json() }
  },

  // ... or just actions that call store.setState() later:
  incrementAsync(state) {
    setTimeout( () => {
      store.setState({ count: state.count+1 })
    }, 100)
  }
})

const App = connect('count', actions)(
  ({ count, increment }) => (
    <div>
      <p>Count: {count}</p>
      <button onClick={increment}>Increment</button>
    </div>
  )
)

export default () => (
  <Provider store={store}>
    <App />
  </Provider>
)
```

### Debug

Make sure to have [Redux devtools extension](https://github.com/zalmoxisus/redux-devtools-extension) previously installed.

```js
import createStore from 'unistore'
import devtools    from 'unistore/devtools'

let initialState = { count: 0 };
let store = process.env.NODE_ENV === 'production' ?  createStore(initialState) : devtools(createStore(initialState));

// ...
```

### Examples

[README Example on CodeSandbox](https://codesandbox.io/s/l7y7w5qkz9)

### API

<!-- Generated by documentation.js. Update this documentation by updating the source code. -->

#### createStore

Creates a new store, which is a tiny evented state container.

**Parameters**

- `state` **[Object](https://developer.mozilla.org/docs/Web/JavaScript/Reference/Global_Objects/Object)** Optional initial state (optional, default `{}`)

**Examples**

```javascript
let store = createStore();
store.subscribe( state => console.log(state) );
store.setState({ a: 'b' });   // logs { a: 'b' }
store.setState({ c: 'd' });   // logs { a: 'b', c: 'd' }
```

Returns **[store](#store)** 

#### store

An observable state container, returned from [createStore](#createstore)

##### action

Create a bound copy of the given action function.
The bound returned function invokes action() and persists the result back to the store.
If the return value of `action` is a Promise, the resolved value will be used as state.

**Parameters**

- `action` **[Function](https://developer.mozilla.org/docs/Web/JavaScript/Reference/Statements/function)** An action of the form `action(state, ...args) -> stateUpdate`

Returns **[Function](https://developer.mozilla.org/docs/Web/JavaScript/Reference/Statements/function)** boundAction()

##### setState

Apply a partial state object to the current state, invoking registered listeners.

**Parameters**

- `update` **[Object](https://developer.mozilla.org/docs/Web/JavaScript/Reference/Global_Objects/Object)** An object with properties to be merged into state
- `overwrite` **[Boolean](https://developer.mozilla.org/docs/Web/JavaScript/Reference/Global_Objects/Boolean)** If `true`, update will replace state instead of being merged into it (optional, default `false`)

##### subscribe

Register a listener function to be called whenever state is changed. Returns an `unsubscribe()` function.

**Parameters**

- `listener` **[Function](https://developer.mozilla.org/docs/Web/JavaScript/Reference/Statements/function)** A function to call when state changes. Gets passed the new state.

Returns **[Function](https://developer.mozilla.org/docs/Web/JavaScript/Reference/Statements/function)** unsubscribe()

##### unsubscribe

Remove a previously-registered listener function.

**Parameters**

- `listener` **[Function](https://developer.mozilla.org/docs/Web/JavaScript/Reference/Statements/function)** The callback previously passed to `subscribe()` that should be removed.

##### getState

Retrieve the current state object.

Returns **[Object](https://developer.mozilla.org/docs/Web/JavaScript/Reference/Global_Objects/Object)** state

#### connect

Wire a component up to the store. Passes state as props, re-renders on change.

**Parameters**

- `mapStateToProps` **([Function](https://developer.mozilla.org/docs/Web/JavaScript/Reference/Statements/function) \| [Array](https://developer.mozilla.org/docs/Web/JavaScript/Reference/Global_Objects/Array) \| [String](https://developer.mozilla.org/docs/Web/JavaScript/Reference/Global_Objects/String))** A function mapping of store state to prop values, or an array/CSV of properties to map.
- `actions` **([Function](https://developer.mozilla.org/docs/Web/JavaScript/Reference/Statements/function) \| [Object](https://developer.mozilla.org/docs/Web/JavaScript/Reference/Global_Objects/Object))?** Action functions (pure state mappings), or a factory returning them. Every action function gets current state as the first parameter and any other params next

**Examples**

```javascript
const Foo = connect('foo,bar')( ({ foo, bar }) => <div /> )
```

```javascript
const actions = { someAction }
const Foo = connect('foo,bar', actions)( ({ foo, bar, someAction }) => <div /> )
```

Returns **Component** ConnectedComponent

#### Provider

**Extends Component**

Provider exposes a store (passed as `props.store`) into context.

Generally, an entire application is wrapped in a single `<Provider>` at the root.

**Parameters**

- `props` **[Object](https://developer.mozilla.org/docs/Web/JavaScript/Reference/Global_Objects/Object)** 
    -   `props.store` **Store** A {Store} instance to expose via context.

### Reporting Issues

Found a problem? Want a new feature? First of all, see if your issue or idea has [already been reported](../../issues).
If not, just open a [new clear and descriptive issue](../../issues/new).

### License

[MIT License](https://oss.ninja/mit/developit) © [Jason Miller](https://jasonformat.com)

[preact]: https://github.com/developit/preact

[react]: https://github.com/facebook/react
````

### FILE: `licenses/reference/edge-runtime-cookies-LICENSE.md.txt`

```yaml
block_id: "REFERENCE-NOTICE-BLOCKS:file10:v1"
operation: CREATE
provenance: VERBATIM
source: "Exact official provenance and declared rendering in licenses/reference/SOURCE_NOTICES.json sources[licenses/reference/edge-runtime-cookies-LICENSE.md.txt], lock SHA256 18b20c9a287eb5ad12217db27f58384812d2d36a7a74c8355390db94f9bdd47b; payload SHA256 e4f76a7a19ef2989dd79339bd3abf8afcf1d6f065e5a10c76c19415ffd727eb3"
license: "MIT"
sha256: "e4f76a7a19ef2989dd79339bd3abf8afcf1d6f065e5a10c76c19415ffd727eb3"
variables: []
secrets_allowed: false
```

````text
The MIT License (MIT)

Copyright (c) 2024 Vercel, Inc.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
````

### FILE: `licenses/reference/edge-runtime-ponyfill-LICENSE.md.txt`

```yaml
block_id: "REFERENCE-NOTICE-BLOCKS:file11:v1"
operation: CREATE
provenance: VERBATIM
source: "Exact official provenance and declared rendering in licenses/reference/SOURCE_NOTICES.json sources[licenses/reference/edge-runtime-ponyfill-LICENSE.md.txt], lock SHA256 18b20c9a287eb5ad12217db27f58384812d2d36a7a74c8355390db94f9bdd47b; payload SHA256 e4f76a7a19ef2989dd79339bd3abf8afcf1d6f065e5a10c76c19415ffd727eb3"
license: "MIT"
sha256: "e4f76a7a19ef2989dd79339bd3abf8afcf1d6f065e5a10c76c19415ffd727eb3"
variables: []
secrets_allowed: false
```

````text
The MIT License (MIT)

Copyright (c) 2024 Vercel, Inc.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
````

### FILE: `licenses/reference/edge-runtime-primitives-LICENSE.md.txt`

```yaml
block_id: "REFERENCE-NOTICE-BLOCKS:file12:v1"
operation: CREATE
provenance: VERBATIM
source: "Exact official provenance and declared rendering in licenses/reference/SOURCE_NOTICES.json sources[licenses/reference/edge-runtime-primitives-LICENSE.md.txt], lock SHA256 18b20c9a287eb5ad12217db27f58384812d2d36a7a74c8355390db94f9bdd47b; payload SHA256 e4f76a7a19ef2989dd79339bd3abf8afcf1d6f065e5a10c76c19415ffd727eb3"
license: "MIT"
sha256: "e4f76a7a19ef2989dd79339bd3abf8afcf1d6f065e5a10c76c19415ffd727eb3"
variables: []
secrets_allowed: false
```

````text
The MIT License (MIT)

Copyright (c) 2024 Vercel, Inc.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
````

### FILE: `licenses/reference/babel-core-LICENSE.txt`

```yaml
block_id: "REFERENCE-NOTICE-BLOCKS:file13:v1"
operation: CREATE
provenance: VERBATIM
source: "Exact official provenance and declared rendering in licenses/reference/SOURCE_NOTICES.json sources[licenses/reference/babel-core-LICENSE.txt], lock SHA256 18b20c9a287eb5ad12217db27f58384812d2d36a7a74c8355390db94f9bdd47b; payload SHA256 117da2af0d4ce0fe1c8e19b5cff9dcd806adf973d328d27b11d4448c4ff24f76"
license: "MIT"
sha256: "117da2af0d4ce0fe1c8e19b5cff9dcd806adf973d328d27b11d4448c4ff24f76"
variables: []
secrets_allowed: false
```

````text
MIT License

Copyright (c) 2014-present Sebastian McKenzie and other contributors

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
````

### FILE: `licenses/reference/dotenv-LICENSE.txt`

```yaml
block_id: "REFERENCE-NOTICE-BLOCKS:file14:v1"
operation: CREATE
provenance: VERBATIM
source: "Exact official provenance and declared rendering in licenses/reference/SOURCE_NOTICES.json sources[licenses/reference/dotenv-LICENSE.txt], lock SHA256 18b20c9a287eb5ad12217db27f58384812d2d36a7a74c8355390db94f9bdd47b; payload SHA256 74b629b24865e1e83c5277ee84590b7937644d6fd959d0c7bdce758676cd2ced"
license: "BSD-2-Clause"
sha256: "74b629b24865e1e83c5277ee84590b7937644d6fd959d0c7bdce758676cd2ced"
variables: []
secrets_allowed: false
```

````text
Copyright (c) 2015, Scott Motte
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### FILE: `licenses/reference/dotenv-expand-LICENSE.txt`

```yaml
block_id: "REFERENCE-NOTICE-BLOCKS:file15:v1"
operation: CREATE
provenance: VERBATIM
source: "Exact official provenance and declared rendering in licenses/reference/SOURCE_NOTICES.json sources[licenses/reference/dotenv-expand-LICENSE.txt], lock SHA256 18b20c9a287eb5ad12217db27f58384812d2d36a7a74c8355390db94f9bdd47b; payload SHA256 3726b9470c3a6b54e1ebcc1d802d37089b5a5fc97273b52b2cee578a4421ec45"
license: "BSD-2-Clause"
sha256: "3726b9470c3a6b54e1ebcc1d802d37089b5a5fc97273b52b2cee578a4421ec45"
variables: []
secrets_allowed: false
```

````text
Copyright (c) 2016, Scott Motte
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

````

### FILE: `licenses/reference/unistore-MIT-terms.txt`

```yaml
block_id: "REFERENCE-NOTICE-BLOCKS:file16:v1"
operation: CREATE
provenance: ADAPTED
source: "Exact official provenance and declared rendering in licenses/reference/SOURCE_NOTICES.json sources[licenses/reference/unistore-MIT-terms.txt], lock SHA256 18b20c9a287eb5ad12217db27f58384812d2d36a7a74c8355390db94f9bdd47b; payload SHA256 15c8d6c00935406bf494329437b74317b2926adcaa479b47bf73f6a060d7244d"
license: "MIT"
sha256: "15c8d6c00935406bf494329437b74317b2926adcaa479b47bf73f6a060d7244d"
variables: []
secrets_allowed: false
```

````text
MIT License

Copyright (c) Jason Miller

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
````

### FILE: `licenses/reference/client-only-MIT-terms.txt`

```yaml
block_id: "REFERENCE-NOTICE-BLOCKS:file17:v1"
operation: CREATE
provenance: ADAPTED
source: "Exact official provenance and declared rendering in licenses/reference/SOURCE_NOTICES.json sources[licenses/reference/client-only-MIT-terms.txt], lock SHA256 18b20c9a287eb5ad12217db27f58384812d2d36a7a74c8355390db94f9bdd47b; payload SHA256 508a77d2e7b51d98adeed32648ad124b7b30241a8e70b2e72c99f92d8e5874d1"
license: "MIT"
sha256: "508a77d2e7b51d98adeed32648ad124b7b30241a8e70b2e72c99f92d8e5874d1"
variables: []
secrets_allowed: false
```

````text
MIT License

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
````

### FILE: `licenses/reference/client-only-package.json.txt`

```yaml
block_id: "REFERENCE-NOTICE-BLOCKS:file18:v1"
operation: CREATE
provenance: VERBATIM
source: "Exact official provenance and declared rendering in licenses/reference/SOURCE_NOTICES.json sources[licenses/reference/client-only-package.json.txt], lock SHA256 18b20c9a287eb5ad12217db27f58384812d2d36a7a74c8355390db94f9bdd47b; payload SHA256 4d6342705767832f299b9a59c28e4275bcf02db19472732f93f67d979441df8f"
license: "MIT"
sha256: "4d6342705767832f299b9a59c28e4275bcf02db19472732f93f67d979441df8f"
variables: []
secrets_allowed: false
```

````text
{
  "name": "client-only",
  "description": "This is a marker package to indicate that a module can only be used in Client Components.",
  "keywords": [
    "react"
  ],
  "version": "0.0.1",
  "homepage": "https://reactjs.org/",
  "bugs": "https://github.com/facebook/react/issues",
  "license": "MIT",
  "files": ["index.js", "error.js"],
  "main": "index.js",
  "exports": {
    ".": {
      "react-server": "./error.js",
      "default": "./index.js"
    }
  }
}
````

### FILE: `licenses/reference/babel-packages-inline-notices.txt`

```yaml
block_id: "REFERENCE-NOTICE-BLOCKS:file19:v1"
operation: CREATE
provenance: ADAPTED
source: "Exact official provenance and declared rendering in licenses/reference/SOURCE_NOTICES.json sources[licenses/reference/babel-packages-inline-notices.txt], lock SHA256 18b20c9a287eb5ad12217db27f58384812d2d36a7a74c8355390db94f9bdd47b; payload SHA256 d08e6b38f8be371bbd045b3defa253d451207a5936aed6dee5d62afdd62e9582"
license: "MIT"
sha256: "d08e6b38f8be371bbd045b3defa253d451207a5936aed6dee5d62afdd62e9582"
variables: []
secrets_allowed: false
```

````text
/*])))))|(0[xX][\da-fA-F]+|0[oO][0-7]+|0[bB][01]+|(?:\d*\.\d+|\d+\.?)(?:[eE][+-]?\d+)?)|((?!\d)(?:(?!\s)[$\w\u0080-\uFFFF]|\\u[\da-fA-F]{4}|\\u\{[\da-fA-F]+\})+)|(--|\+\+|&&|\|\||=>|\.{3}|(?:[+\-\/%&|^]|\*{1,2}|<{1,2}|>{1,3}|!=?|={1,2})=?|[?~.,:;[\](){}])|(\s+)|(^$|[\s\S])/g;r.matchToToken=function(e){var r={type:"invalid",value:e[0],closed:undefined};if(e[1])r.type="string",r.closed=!!(e[3]||e[4]);else if(e[5])r.type="comment";else if(e[6])r.type="comment",r.closed=!!e[7];else if(e[8])r.type="regex";else if(e[9])r.type="number";else if(e[10])r.type="name";else if(e[11])r.type="punctuator";else if(e[12])r.type="whitespace";return r}},8995:e=>{var r="Expected a function";var t=0/0;var s="[object Symbol]";var a=/^\s+|\s+$/g;var n=/^[-+]0x[0-9a-f]+$/i;var o=/^0b[01]+$/i;var i=/^0o[0-7]+$/i;var l=parseInt;var c=typeof global=="object"&&global&&global.Object===Object&&global;var d=typeof self=="object"&&self&&self.Object===Object&&self;var u=c||d||Function("return this")();var p=Object.prototype;var f=p.toString;var y=Math.max,g=Math.min;var now=function(){return u.Date.now()};function debounce(e,t,s){var a,n,o,i,l,c,d=0,u=false,p=false,f=true;if(typeof e!="function"){throw new TypeError(r)}t=toNumber(t)||0;if(isObject(s)){u=!!s.leading;p="maxWait"in s;o=p?y(toNumber(s.maxWait)||0,t):o;f="trailing"in s?!!s.trailing:f}function invokeFunc(r){var t=a,s=n;a=n=undefined;d=r;i=e.apply(s,t);return i}function leadingEdge(e){d=e;l=setTimeout(timerExpired,t);return u?invokeFunc(e):i}function remainingWait(e){var r=e-c,s=e-d,a=t-r;return p?g(a,o-s):a}function shouldInvoke(e){var r=e-c,s=e-d;return c===undefined||r>=t||r<0||p&&s>=o}function timerExpired(){var e=now();if(shouldInvoke(e)){return trailingEdge(e)}l=setTimeout(timerExpired,remainingWait(e))}function trailingEdge(e){l=undefined;if(f&&a){return invokeFunc(e)}a=n=undefined;return i}function cancel(){if(l!==undefined){clearTimeout(l)}d=0;a=c=n=l=undefined}function flush(){return l===undefined?i:trailingEdge(now())}function debounced(){var e=now(),r=shouldInvoke(e);a=arguments;n=this;c=e;if(r){if(l===undefined){return leadingEdge(c)}if(p){l=setTimeout(timerExpired,t);return invokeFunc(c)}}if(l===undefined){l=setTimeout(timerExpired,t)}return i}debounced.cancel=cancel;debounced.flush=flush;return debounced}function isObject(e){var r=typeof e;return!!e&&(r=="object"||r=="function")}function isObjectLike(e){return!!e&&typeof e=="object"}function isSymbol(e){return typeof e=="symbol"||isObjectLike(e)&&f.call(e)==s}function toNumber(e){if(typeof e=="number"){return e}if(isSymbol(e)){return t}if(isObject(e)){var r=typeof e.valueOf=="function"?e.valueOf():e;e=isObject(r)?r+"":r}if(typeof e!="string"){return e===0?e:+e}e=e.replace(a,"");var s=o.test(e);return s||i.test(e)?l(e.slice(2),s?2:8):n.test(e)?t:+e}e.exports=debounce},4464:function(e,r,t){e=t.nmd(e);
/**
 * @license
 * Lodash <https://lodash.com/>
 * Copyright OpenJS Foundation and other contributors <https://openjsf.org/>
 * Released under MIT license <https://lodash.com/license>
 * Based on Underscore.js 1.8.3 <http://underscorejs.org/LICENSE>
 * Copyright Jeremy Ashkenas, DocumentCloud and Investigative Reporters & Editors
 */

/*! https://mths.be/regenerate v1.4.2 by @mathias | MIT license */

/*!
 * regjsgen 0.8.0
 * Copyright 2014-2023 Benjamin Tan <https://ofcr.se/>
 * Available under the MIT license <https://github.com/bnjmnt4n/regjsgen/blob/main/LICENSE-MIT.txt>
 */
````

### FILE: `licenses/reference/runtime-notice-lock.json`

```yaml
block_id: "REFERENCE-NOTICE-BLOCKS:file20:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "7a830609937f34472e01f8f1de7fc3cd26a9177958984f89b64e5af60095d17a"
variables: []
secrets_allowed: false
```

````json
{
  "build_only_exclusions": [
    "@next/swc-win32-x64-msvc@16.3.4",
    "@rolldown/binding-win32-x64-msvc@1.2.5",
    "server-only@0.0.1",
    "stackback@0.0.2"
  ],
  "manifests": {
    "node_modules/.pnpm/@next+env@16.3.4/node_modules/@next/env/package.json": {
      "license": "MIT",
      "name": "@next/env",
      "notice_refs": [
        "licenses/reference/next-MIT.txt",
        "licenses/reference/dotenv-LICENSE.txt",
        "licenses/reference/dotenv-expand-LICENSE.txt"
      ],
      "sha256": "229045587b8f7898b12126f99fadd1021c8c928f207f92f22146908d46899b5b",
      "source_basis": "Next fixed source package16.3.4 MIT and declared embedded dotenv16.3.1/dotenv-expand10.0.0 BSD notices.",
      "version": "16.3.4"
    },
    "node_modules/.pnpm/@swc+helpers@0.5.23/node_modules/@swc/helpers/_/_interop_require_default/package.json": {
      "license": null,
      "name": null,
      "notice_refs": [
        "node_modules/.pnpm/@swc+helpers@0.5.23/node_modules/@swc/helpers/LICENSE"
      ],
      "sha256": "9cf5cdb2fda7dff60167b91a89a618c78d63f6e044e7681075efdb03a3834ed5",
      "source_basis": "Unnamed package export metadata within existing package; nearest explicit parent license, no independent package identity invented.",
      "version": null
    },
    "node_modules/.pnpm/@swc+helpers@0.5.23/node_modules/@swc/helpers/_/_interop_require_wildcard/package.json": {
      "license": null,
      "name": null,
      "notice_refs": [
        "node_modules/.pnpm/@swc+helpers@0.5.23/node_modules/@swc/helpers/LICENSE"
      ],
      "sha256": "0771ab9c96c3d022cfabd5f7aebfeb43192893b24b025b5b0c941c0b5b450c1f",
      "source_basis": "Unnamed package export metadata within existing package; nearest explicit parent license, no independent package identity invented.",
      "version": null
    },
    "node_modules/.pnpm/@swc+helpers@0.5.23/node_modules/@swc/helpers/package.json": {
      "license": "Apache-2.0",
      "name": "@swc/helpers",
      "notice_refs": [
        "node_modules/.pnpm/@swc+helpers@0.5.23/node_modules/@swc/helpers/LICENSE"
      ],
      "sha256": "10ee79b75759afbf8cfa2b8e5a2eff946c487143e935fc31b2433aa68be84596",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": "0.5.23"
    },
    "node_modules/.pnpm/baseline-browser-mapping@2.11.18/node_modules/baseline-browser-mapping/package.json": {
      "license": "Apache-2.0",
      "name": "baseline-browser-mapping",
      "notice_refs": [
        "node_modules/.pnpm/baseline-browser-mapping@2.11.18/node_modules/baseline-browser-mapping/LICENSE.txt"
      ],
      "sha256": "d1e5409452277a501b376a99b6da53a5084042bea082e51dc476aea3ff2f6715",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": "2.11.18"
    },
    "node_modules/.pnpm/caniuse-lite@1.0.30001809/node_modules/caniuse-lite/package.json": {
      "license": "CC-BY-4.0",
      "name": "caniuse-lite",
      "notice_refs": [
        "node_modules/.pnpm/caniuse-lite@1.0.30001809/node_modules/caniuse-lite/LICENSE"
      ],
      "sha256": "fd9f8880647d958c2711e30850d1638d262fa3a477c1cf94ef9e4ea46047d9b8",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": "1.0.30001809"
    },
    "node_modules/.pnpm/client-only@0.0.1/node_modules/client-only/package.json": {
      "license": "MIT",
      "name": "client-only",
      "notice_refs": [
        "licenses/reference/client-only-MIT-terms.txt",
        "licenses/reference/client-only-package.json.txt"
      ],
      "sha256": "4d6342705767832f299b9a59c28e4275bcf02db19472732f93f67d979441df8f",
      "source_basis": "Fixed official source/tarball supplement declared in licenses/reference/SOURCE_NOTICES.json; expression/notice scope preserved.",
      "version": "0.0.1"
    },
    "node_modules/.pnpm/nanoid@3.3.18/node_modules/nanoid/non-secure/package.json": {
      "license": null,
      "name": null,
      "notice_refs": [
        "node_modules/.pnpm/nanoid@3.3.18/node_modules/nanoid/LICENSE"
      ],
      "sha256": "2f9b8c4bee9312538f388b49c13987f4565ac25d6756fabdbbb771d339945e3e",
      "source_basis": "Unnamed package export metadata within existing package; nearest explicit parent license, no independent package identity invented.",
      "version": null
    },
    "node_modules/.pnpm/nanoid@3.3.18/node_modules/nanoid/package.json": {
      "license": "MIT",
      "name": "nanoid",
      "notice_refs": [
        "node_modules/.pnpm/nanoid@3.3.18/node_modules/nanoid/LICENSE"
      ],
      "sha256": "18673eb0d7ed43dd871d59dc22a5670303f5d4eb9b17783b1f698f5840ce1fff",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": "3.3.18"
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@edge-runtime/cookies/package.json": {
      "license": "MIT",
      "name": "@edge-runtime/cookies",
      "notice_refs": [
        "licenses/reference/edge-runtime-cookies-LICENSE.md.txt"
      ],
      "sha256": "84dec8086a6d6a1176ff52a54dbc4542aeca7241214a0cf7d5af2870ee422db2",
      "source_basis": "Fixed official source/tarball supplement declared in licenses/reference/SOURCE_NOTICES.json; expression/notice scope preserved.",
      "version": "6.0.0"
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@edge-runtime/ponyfill/package.json": {
      "license": "MIT",
      "name": "@edge-runtime/ponyfill",
      "notice_refs": [
        "licenses/reference/edge-runtime-ponyfill-LICENSE.md.txt"
      ],
      "sha256": "04d31884423b345e3abddc712adc85e10347ccf057630b23a930f1ba9af56708",
      "source_basis": "Fixed official source/tarball supplement declared in licenses/reference/SOURCE_NOTICES.json; expression/notice scope preserved.",
      "version": "4.0.0"
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@edge-runtime/primitives/package.json": {
      "license": "MIT",
      "name": "@edge-runtime/primitives",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@edge-runtime/primitives/fetch.js.LEGAL.txt",
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@edge-runtime/primitives/load.js.LEGAL.txt",
        "licenses/reference/edge-runtime-primitives-LICENSE.md.txt"
      ],
      "sha256": "ee5dfbcb78d0c753908a0643fc23d7486cf4478872789bca8722a88f71d5c213",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": "6.0.0"
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@hapi/accept/package.json": {
      "license": "BSD-3-Clause",
      "name": "@hapi/accept",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@hapi/accept/LICENSE"
      ],
      "sha256": "0e264b353f98a5047f47a4af7cd718888f33345d78a8bbea30eaaac585059abc",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@mswjs/interceptors/ClientRequest/package.json": {
      "license": "MIT",
      "name": "@mswjs/interceptors",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@mswjs/interceptors/ClientRequest/LICENSE"
      ],
      "sha256": "d01c0050c9f635ad97c5e49368f456db68dc474bb7f8ada9b78121c90ed9f287",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@napi-rs/triples/package.json": {
      "license": "MIT",
      "name": "@napi-rs/triples",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@napi-rs/triples/LICENSE"
      ],
      "sha256": "828a978674962e997f0f5d2920124bc17a4e9f1ba7a55d7f8ea72e25675bcef8",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@next/font/package.json": {
      "license": "MIT",
      "name": "@next/font",
      "notice_refs": [
        "licenses/reference/next-MIT.txt"
      ],
      "sha256": "c99740a4da3ab17fdbfed0b601e2e8084a7a9463f62a13e5a7ed106648976586",
      "source_basis": "Next-owned font workspace manifest and copy recipe at fixed official Next commit.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@opentelemetry/api/package.json": {
      "license": "Apache-2.0",
      "name": "@opentelemetry/api",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@opentelemetry/api/LICENSE"
      ],
      "sha256": "a7352d5f47b67b86f38d634360f761fa2356a4256f42a5239a12892b6b08074a",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@vercel/detect-agent/package.json": {
      "license": "Apache-2.0",
      "name": "@vercel/detect-agent",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@vercel/detect-agent/LICENSE"
      ],
      "sha256": "c8ed9a25f3b85a9c400e05f23e998bdfa820e83665c71fd9164eb584f7684afd",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@vercel/nft/package.json": {
      "license": "MIT",
      "name": "@vercel/nft",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@vercel/nft/LICENSE"
      ],
      "sha256": "3d26054595f0ea7f5f26e0b5a269521f77fa4fd030dd7c964c492c7a018e4d22",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/acorn/package.json": {
      "license": "MIT",
      "name": "acorn",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/acorn/LICENSE"
      ],
      "sha256": "aa6eff9593179bc7379017992eda314ce97bd7c6576869d82d514139ce612cca",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/assert/package.json": {
      "license": "MIT",
      "name": "assert",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/assert/LICENSE"
      ],
      "sha256": "cee48f1c2973cbe0cc9075fa1ad46540571e39f82368b90680e3f08385acc67d",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/async-retry/package.json": {
      "license": "MIT",
      "name": "async-retry",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/async-retry/LICENSE"
      ],
      "sha256": "b17113b29f3ee160e536bacf5332fb1cc91f81f5dba02db749491a469108b9e5",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/async-sema/package.json": {
      "license": "MIT",
      "name": "async-sema",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/async-sema/LICENSE"
      ],
      "sha256": "327a56df6511702ea28014e18425bd50300890fb433d1e8e8cf200ab33ac70c5",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/babel-packages/package.json": {
      "license": null,
      "name": "babel-packages",
      "notice_refs": [
        "licenses/reference/next-MIT.txt",
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/babel/LICENSE",
        "licenses/reference/babel-packages-inline-notices.txt",
        "licenses/reference/babel-core-LICENSE.txt"
      ],
      "sha256": "828b12b4e497121b9a04735e1defeeb93dfde3c7531ce7b89b568ed520c86b4e",
      "source_basis": "Next fixed Babel packages-bundle import/build recipe; original Babel MIT text plus exact preserved inline copyright/license comments. No original package versions invented for compiled bundle.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/babel/package.json": {
      "license": "MIT",
      "name": "@babel/core",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/babel/LICENSE",
        "licenses/reference/babel-core-LICENSE.txt"
      ],
      "sha256": "74a866ce237c245a6657979b9255e4e702711ba2f83a72a9d1d9d5647037f322",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/browserify-zlib/package.json": {
      "license": "MIT",
      "name": "browserify-zlib",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/browserify-zlib/LICENSE"
      ],
      "sha256": "fbfda5f3dd28400e7e43080c996f39ceb97ff45dc94237b218bdb9785afa163d",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/browserslist/package.json": {
      "license": "MIT",
      "name": "browserslist",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/browserslist/LICENSE"
      ],
      "sha256": "5c9be14e5f45d01c0f11ce740c55b61008e3878d9003121989a2d5700b1e7ff1",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/buffer/package.json": {
      "license": "MIT",
      "name": "buffer",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/buffer/LICENSE"
      ],
      "sha256": "beec4c77cd84bb8bbcbd9ad3ddc435bf9b4369884a56cb2c4cfb9e8805eb87cb",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/busboy/package.json": {
      "license": null,
      "name": "busboy",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/busboy/LICENSE"
      ],
      "sha256": "7ec808358d01f7e153eec8ace30496d69f066e589abd3b1deaaadaa24a2c9e78",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/bytes/package.json": {
      "license": "MIT",
      "name": "bytes",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/bytes/LICENSE"
      ],
      "sha256": "f67cda9a406b0173dd8bce3165b246a7461f0b6aa8b1db5d773a8d858e618e0d",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/ci-info/package.json": {
      "license": "MIT",
      "name": "ci-info",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/ci-info/LICENSE"
      ],
      "sha256": "5f266b005be5715f7971df5255242ce2684ecb7d073ee597ad7a01c5049d9889",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/commander/package.json": {
      "license": "MIT",
      "name": "commander",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/commander/LICENSE"
      ],
      "sha256": "23bca2c540bfa9a86bfc07d3d3322a5821c11d9d1c8d592bc5bfbc8b0641745d",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/comment-json/package.json": {
      "license": "MIT",
      "name": "comment-json",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/comment-json/LICENSE"
      ],
      "sha256": "a3f86bf44c9a56dd6c9724a32938c8aa3876412ebb17400b613563f79017bf5d",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/compression/package.json": {
      "license": "MIT",
      "name": "compression",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/compression/LICENSE"
      ],
      "sha256": "bb26d52201279816e928ccd1a58ecdd0fab49d7feb3136b8c74e82b09a0f2bb9",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/conf/package.json": {
      "license": "MIT",
      "name": "conf",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/conf/LICENSE"
      ],
      "sha256": "0a761c63beca5354c035de3f8a0844a1a64f62e94d08ef8c57e2571701cb0d90",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/constants-browserify/package.json": {
      "license": null,
      "name": "constants-browserify",
      "notice_refs": [
        "licenses/reference/constants-browserify-README.md.txt"
      ],
      "sha256": "57ff101ce1ee4f9ff1cd717424bfc9ba13db28c0ed99e7ce06ff8e4f080a56f2",
      "source_basis": "Fixed official source/tarball supplement declared in licenses/reference/SOURCE_NOTICES.json; expression/notice scope preserved.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/content-disposition/package.json": {
      "license": "MIT",
      "name": "content-disposition",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/content-disposition/LICENSE"
      ],
      "sha256": "deee6f8ca16ef491cfcc902176392f0c06fb7bea87338f521c384a174553c7bb",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/content-type/package.json": {
      "license": "MIT",
      "name": "content-type",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/content-type/LICENSE"
      ],
      "sha256": "0c481f9c3a287b38fe0ccaa61173cd46ae210e1f9a61bca6bc6abacef66c4e3c",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/cookie/package.json": {
      "license": "MIT",
      "name": "cookie",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/cookie/LICENSE"
      ],
      "sha256": "202616a73f282bc2c6243db8781e044d0fa1090d15b0a35eea88431ebcb8ac84",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/cross-spawn/package.json": {
      "license": "MIT",
      "name": "cross-spawn",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/cross-spawn/LICENSE"
      ],
      "sha256": "427908598a096747431c250f457a5440b2b4828d6243db1c98d89ac8467b84a3",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/crypto-browserify/package.json": {
      "license": "MIT",
      "name": "crypto-browserify",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/crypto-browserify/LICENSE"
      ],
      "sha256": "231080cafc3e64b7bcc175551a70ed4b3b15396d721417c56f227fa4aeda45d4",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/data-uri-to-buffer/package.json": {
      "license": "MIT",
      "name": "data-uri-to-buffer",
      "notice_refs": [
        "licenses/reference/data-uri-to-buffer-README.md.txt"
      ],
      "sha256": "7884ade4aa4ccd6fa4c15a37591172fccc6b543e3cc6c1a0fef7888f39d1b39e",
      "source_basis": "Fixed official source/tarball supplement declared in licenses/reference/SOURCE_NOTICES.json; expression/notice scope preserved.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/debug/package.json": {
      "license": "MIT",
      "name": "debug",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/debug/LICENSE"
      ],
      "sha256": "ead56d5fb4a68fe23ce148b366356f9a56c8f12e14891db6bd8425b339cb360a",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/devalue/package.json": {
      "license": "MIT",
      "name": "devalue",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/devalue/LICENSE"
      ],
      "sha256": "02cc0546b6f93a1d07c2f10db73edc761a7b6353130982b90113bb9b051c8b4f",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/domain-browser/package.json": {
      "license": "MIT",
      "name": "domain-browser",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/domain-browser/LICENSE"
      ],
      "sha256": "102f56b1f979828dc82944853a4df4d728d44f61c34e729ecd07b82dca0ba56a",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/edge-runtime/package.json": {
      "license": "MIT",
      "name": "edge-runtime",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/edge-runtime/LICENSE"
      ],
      "sha256": "9887ce733e2d0abdb8a5c2307e1ddd952f144bd8df2fd67d53595732a7590b1e",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/events/package.json": {
      "license": "MIT",
      "name": "events",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/events/LICENSE"
      ],
      "sha256": "d1b2b2fc773bb2dc793372c848c00f60c721137ba0762370a4516f603670f143",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/find-up/package.json": {
      "license": "MIT",
      "name": "find-up",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/find-up/LICENSE"
      ],
      "sha256": "7658d63408818236df3c4fd7107c5712d73106f7939649043b81d2fd1c18f705",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/fresh/package.json": {
      "license": "MIT",
      "name": "fresh",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/fresh/LICENSE"
      ],
      "sha256": "0f8f39469d55646a4a359833c464529df8720059c1db552f9e94036c06ac4f47",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/gzip-size/package.json": {
      "license": "MIT",
      "name": "gzip-size",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/gzip-size/LICENSE"
      ],
      "sha256": "88fa90031b294d01b1ccd41443ac213710cd3aeb3abae96f2d9dd25b31c697ee",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/http-proxy-agent/package.json": {
      "license": "MIT",
      "name": "http-proxy-agent",
      "notice_refs": [
        "licenses/reference/http-proxy-agent-README.md.txt"
      ],
      "sha256": "00884de8effb02adaab8ea8fce8d0ac9240dee19df517eb87dc35cdf5ba8ed79",
      "source_basis": "Fixed official source/tarball supplement declared in licenses/reference/SOURCE_NOTICES.json; expression/notice scope preserved.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/https-browserify/package.json": {
      "license": "MIT",
      "name": "https-browserify",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/https-browserify/LICENSE"
      ],
      "sha256": "f0e319d9ff18e39dae7cd7a7cba7e22e8faae1e7b149dd86279fa5ae7730fa7b",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/https-proxy-agent/package.json": {
      "license": "MIT",
      "name": "https-proxy-agent",
      "notice_refs": [
        "licenses/reference/https-proxy-agent-README.md.txt"
      ],
      "sha256": "144ceffc497b5f331153c9d291f682f5264595c7b8711baa1d93de5f53897c4e",
      "source_basis": "Fixed official source/tarball supplement declared in licenses/reference/SOURCE_NOTICES.json; expression/notice scope preserved.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/httpxy/package.json": {
      "license": "MIT",
      "name": "httpxy",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/httpxy/LICENSE"
      ],
      "sha256": "1a24e02af08692a30d9a29580d267d14c70e307df1d1ba765fee94787fb95afa",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/icss-utils/package.json": {
      "license": "ISC",
      "name": "icss-utils",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/icss-utils/LICENSE"
      ],
      "sha256": "a28db1d8a722f0f40b351bbca8c8e8559bbe09a4021448f954c8b9d4860f92d5",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/ignore-loader/package.json": {
      "license": null,
      "name": "ignore-loader",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/ignore-loader/LICENSE"
      ],
      "sha256": "570b3e062c5757cb39e8f81060b4884ac32fdc8cdcd5a7e7d99385607239bf73",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/image-size/package.json": {
      "license": "MIT",
      "name": "image-size",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/image-size/LICENSE"
      ],
      "sha256": "0f08e3455955e889e24ae141610e910ff606c4ff779f1cb67407ad5036a15bb9",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/ipaddr.js/package.json": {
      "license": "MIT",
      "name": "ipaddr.js",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/ipaddr.js/LICENSE"
      ],
      "sha256": "28f0b70ca1c217b65158ff0d05247025b2e8810351902946d3c5ff3dad209507",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/is-animated/package.json": {
      "license": "MIT",
      "name": "is-animated",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/is-animated/LICENSE"
      ],
      "sha256": "61fa46c5a1bb37239f12272b650527742595eb63a5104f7e625469370562585a",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/is-docker/package.json": {
      "license": "MIT",
      "name": "is-docker",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/is-docker/LICENSE"
      ],
      "sha256": "aac5594fc4ae98f4b297f52c91fede06bc9eca25915a216d192a6bd3952609bb",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/is-wsl/package.json": {
      "license": "MIT",
      "name": "is-wsl",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/is-wsl/LICENSE"
      ],
      "sha256": "189b610be03f7aab4bf124bf31fcf75a03f6504b115648991c2c10dc8a11263e",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/jest-worker/package.json": {
      "license": "MIT",
      "name": "jest-worker",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/jest-worker/LICENSE"
      ],
      "sha256": "9b89fefb8dad25d6f083de6c34d156be2ced3cab898354e6b1360260d6f843d9",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/json5/package.json": {
      "license": "MIT",
      "name": "json5",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/json5/LICENSE"
      ],
      "sha256": "5662b636583985d168a45c9c61b1a874ad5dd0cfc2e95984b1a218e73f93c18f",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/jsonwebtoken/package.json": {
      "license": "MIT",
      "name": "jsonwebtoken",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/jsonwebtoken/LICENSE"
      ],
      "sha256": "f35673a4e7c1dcaa11adc4ea0a6d0272692877bbf7e82a7831a5bbbcd3bb4400",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/loader-utils2/package.json": {
      "license": "MIT",
      "name": "loader-utils",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/loader-utils2/LICENSE"
      ],
      "sha256": "997729fb489c230d7ee272c392237d90b6d2aea6006e66c664374caa0af55db0",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/loader-utils3/package.json": {
      "license": "MIT",
      "name": "loader-utils",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/loader-utils3/LICENSE"
      ],
      "sha256": "997729fb489c230d7ee272c392237d90b6d2aea6006e66c664374caa0af55db0",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/lodash.curry/package.json": {
      "license": "MIT",
      "name": "lodash.curry",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/lodash.curry/LICENSE"
      ],
      "sha256": "423851cd0529ddf3b16f399efa82ef71aa54e03faa275c680e46468c07ab0ed5",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/lru-cache/package.json": {
      "license": "ISC",
      "name": "lru-cache",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/lru-cache/LICENSE"
      ],
      "sha256": "12b2e366f83a8cf697dec2cc400f135bfd07fb62d02df0df5fcab4d7cb98f631",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/mini-css-extract-plugin/package.json": {
      "license": "MIT",
      "name": "mini-css-extract-plugin",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/mini-css-extract-plugin/LICENSE"
      ],
      "sha256": "a940820e002afcfc41db9eb8856fd709af8fe6002ac378e02206ff4dec9bfc62",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/nanoid/package.json": {
      "license": "MIT",
      "name": "nanoid",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/nanoid/LICENSE"
      ],
      "sha256": "9383f87d5189f73399e60c5b9666fe4930f35ba304355e3555c74f5d9d2ff006",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/native-url/package.json": {
      "license": "Apache-2.0",
      "name": "native-url",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/native-url/LICENSE"
      ],
      "sha256": "8056830cb93b8f4904d1e63c7e40abdcac725623e19aa6cb0a07964757570402",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/neo-async/package.json": {
      "license": "MIT",
      "name": "neo-async",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/neo-async/LICENSE"
      ],
      "sha256": "ab030dfc79c2b9b39e329d2fb518891cef7479cc62c95cf1a691e19e2364e0e1",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/os-browserify/package.json": {
      "license": "MIT",
      "name": "os-browserify",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/os-browserify/LICENSE"
      ],
      "sha256": "b27bdc906022bfe4361aeeed67b6874df02c7b324db6eac870f8f89a5113af5a",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/p-limit/package.json": {
      "license": "MIT",
      "name": "p-limit",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/p-limit/LICENSE"
      ],
      "sha256": "40495ef05f0ea39ba8d64eb98e492f2d461985e18a937b87b523b3322dbae5ea",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/p-queue/package.json": {
      "license": "MIT",
      "name": "p-queue",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/p-queue/LICENSE"
      ],
      "sha256": "50f88936a3223c57edfde9e3844e932e20501f49403a761ec210b39d565bc2a6",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/path-browserify/package.json": {
      "license": "MIT",
      "name": "path-browserify",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/path-browserify/LICENSE"
      ],
      "sha256": "6be7adc0449e83f7f2295aad331b7cde780dda2c97c54ba51f51de74114b02a1",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/path-to-regexp/package.json": {
      "license": "MIT",
      "name": "path-to-regexp",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/path-to-regexp/LICENSE"
      ],
      "sha256": "5bf6a22265a08b61b2ab9e05a77cea0dd0f078dd086d17769a2c5b4f93a7426a",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/picomatch/package.json": {
      "license": "MIT",
      "name": "picomatch",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/picomatch/LICENSE"
      ],
      "sha256": "bcb6e75087cfe22755400218f098b6b36805627e5f54171abf58454f6484ac47",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/postcss-flexbugs-fixes/package.json": {
      "license": "MIT",
      "name": "postcss-flexbugs-fixes",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/postcss-flexbugs-fixes/LICENSE"
      ],
      "sha256": "233739e87858b972ff4100d4f7f7a9fa5adab75b7f3202ecab808892ad0d7439",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/postcss-modules-extract-imports/package.json": {
      "license": "ISC",
      "name": "postcss-modules-extract-imports",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/postcss-modules-extract-imports/LICENSE"
      ],
      "sha256": "37b6b10887f8d084500ab94e349f3db0f0ef02dcb6439a4e6aa849b13a81ff8f",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/postcss-modules-local-by-default/package.json": {
      "license": "MIT",
      "name": "postcss-modules-local-by-default",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/postcss-modules-local-by-default/LICENSE"
      ],
      "sha256": "2b99f876f88e1ff8bbd2bcabec5120d9fb12015c1050a8b2e0cd6aed902c2a09",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/postcss-modules-scope/package.json": {
      "license": "ISC",
      "name": "postcss-modules-scope",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/postcss-modules-scope/LICENSE"
      ],
      "sha256": "0594c88976370138023c222480d65eaec2591b859366c7dcd33cf5f2576a1d99",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/postcss-modules-values/package.json": {
      "license": "ISC",
      "name": "postcss-modules-values",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/postcss-modules-values/LICENSE"
      ],
      "sha256": "75686122910108295fd244fcb2e7ed2cd8f2785e571152c15838febaeb5f9e99",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/postcss-preset-env/package.json": {
      "license": "CC0-1.0",
      "name": "postcss-preset-env",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/postcss-preset-env/LICENSE"
      ],
      "sha256": "22056897f5fa4de9e64c5ca87d14008b3ddc01be85fe33c1f3eeae513dadea22",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/postcss-scss/package.json": {
      "license": "MIT",
      "name": "postcss-scss",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/postcss-scss/LICENSE"
      ],
      "sha256": "d054cb1fce9ea5852729c829bab799126eae2c07ca3f33d761760cae829cb10f",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/postcss-value-parser/package.json": {
      "license": "MIT",
      "name": "postcss-value-parser",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/postcss-value-parser/LICENSE"
      ],
      "sha256": "096375657ecf115fc78281f8c43c0cec36cd6396e2801d463f700fd987d9a8e8",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/process/package.json": {
      "license": "MIT",
      "name": "process",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/process/LICENSE"
      ],
      "sha256": "f1c6ff9dc88ded25db7f9e314b18124c8e5648c8fa4de324968d6c83914a5cfd",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/punycode/package.json": {
      "license": "MIT",
      "name": "punycode",
      "notice_refs": [
        "licenses/reference/punycode-LICENSE-MIT.txt.txt"
      ],
      "sha256": "4d69f3d2928226148fd39c2b9724d60c560ca13f01729937964ea3bfb4e2e19b",
      "source_basis": "Fixed official source/tarball supplement declared in licenses/reference/SOURCE_NOTICES.json; expression/notice scope preserved.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/querystring-es3/package.json": {
      "license": null,
      "name": "querystring-es3",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/querystring-es3/LICENSE"
      ],
      "sha256": "deef9957fd78ab6495db1c8ac121d9bd7b0e65af73f50a554efac5b5a7144f08",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/raw-body/package.json": {
      "license": "MIT",
      "name": "raw-body",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/raw-body/LICENSE"
      ],
      "sha256": "ebc1e2456bc3c948d587d9cbbfc845a4cd4a839732d2072afdbf23b918792ef0",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/react-is/package.json": {
      "license": "MIT",
      "name": "react-is",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/react-is/LICENSE"
      ],
      "sha256": "e58e18130d2a543b485afacc83456e161c54c8e3e63ee3174d7160167d9ddc61",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": "19.3.0-canary-cbb046ab-20260731"
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/react-refresh/package.json": {
      "license": "MIT",
      "name": "react-refresh",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/react-refresh/LICENSE"
      ],
      "sha256": "7634851c499058e71f98f6333982d2ed2bcd109bfe0e86a71d4134281a7c497c",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": "0.12.0"
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/regenerator-runtime/package.json": {
      "license": "MIT",
      "name": "regenerator-runtime",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/regenerator-runtime/LICENSE"
      ],
      "sha256": "19a1b69415c9661789cb9c87e822a6fbd976d78ac8d31d7d6ce30acc86a491cf",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": "0.13.4"
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/safe-stable-stringify/package.json": {
      "license": "MIT",
      "name": "safe-stable-stringify",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/safe-stable-stringify/LICENSE"
      ],
      "sha256": "5fd2e6d5ad1018bfb55b5af6bb2bdd61b1314b7993586953e5c6841eb4b936d6",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/sass-loader/package.json": {
      "license": "MIT",
      "name": "sass-loader",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/sass-loader/LICENSE"
      ],
      "sha256": "021cd8eefc24b305b58e1379f83368f54083baa3a8d06f529bbe79d3380ddec8",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/schema-utils3/package.json": {
      "license": "MIT",
      "name": "schema-utils",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/schema-utils3/LICENSE"
      ],
      "sha256": "1213b565bcb89186df47543861eab7c2b262a71c9a3805cfb82d539f1d73d6bc",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/semver/package.json": {
      "license": "ISC",
      "name": "semver",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/semver/LICENSE"
      ],
      "sha256": "3b7c2da335e6f7d178a13d32917fc6ad85d62e6ca2cda2513e52815aa7ca95e9",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/send/package.json": {
      "license": "MIT",
      "name": "send",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/send/LICENSE"
      ],
      "sha256": "9ee32160a660db32561814157aec6fa4054e79ff838c80981808d839bfec6938",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/setimmediate/package.json": {
      "license": "MIT",
      "name": "setimmediate",
      "notice_refs": [
        "licenses/reference/setimmediate-LICENSE.txt.txt"
      ],
      "sha256": "07c28ab21f75001cd3c5f0e49a2171bc5c55536b55296b4fc9425810a01bea8a",
      "source_basis": "Fixed official source/tarball supplement declared in licenses/reference/SOURCE_NOTICES.json; expression/notice scope preserved.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/shell-quote/package.json": {
      "license": "MIT",
      "name": "shell-quote",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/shell-quote/LICENSE"
      ],
      "sha256": "907fda3052f7e1ece852b1cf025642aa3d90cb96ccccf4d568017d974002e9af",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/source-map/package.json": {
      "license": "BSD-3-Clause",
      "name": "source-map",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/source-map/LICENSE"
      ],
      "sha256": "64b17ef1036027151303e0d09a47a0507901878accd7f6d8276a71c06a2e4312",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/source-map08/package.json": {
      "license": "BSD-3-Clause",
      "name": "source-map08",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/source-map08/LICENSE"
      ],
      "sha256": "3c7de96c15e46a376f876ea51b70daac14988f5a8bd261565e21a0e17b14a8c1",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/stacktrace-parser/package.json": {
      "license": "MIT",
      "name": "stacktrace-parser",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/stacktrace-parser/LICENSE"
      ],
      "sha256": "72e32b3357e6fab245dc57c202a274cdd75020cd23ce6a1a9916bfa6127d80cf",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/stream-browserify/package.json": {
      "license": "MIT",
      "name": "stream-browserify",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/stream-browserify/LICENSE"
      ],
      "sha256": "fe7da53b4104f1ef65c7f28d8ba2e0e46f1ab61e71dd11f0bf374dd1742ddcf9",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/stream-http/package.json": {
      "license": "MIT",
      "name": "stream-http",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/stream-http/LICENSE"
      ],
      "sha256": "e5ec05a5ff5a2bf75d194542d26dc7e5fe46a22a60f020e08729ab70bffba947",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/string-hash/package.json": {
      "license": "CC0-1.0",
      "name": "string-hash",
      "notice_refs": [
        "licenses/reference/string-hash-README.md.txt"
      ],
      "sha256": "afb2563fe191f4e7eecc9e6bf1131f07da3f855e6e51b2d517206ae737b343b0",
      "source_basis": "Fixed official source/tarball supplement declared in licenses/reference/SOURCE_NOTICES.json; expression/notice scope preserved.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/string_decoder/package.json": {
      "license": "MIT",
      "name": "string_decoder",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/string_decoder/LICENSE"
      ],
      "sha256": "0ff89407b42b08cb0cfe9aaa590b593693e3b07ebb16902f42297c00453f25ab",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/strip-ansi/package.json": {
      "license": "MIT",
      "name": "strip-ansi",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/strip-ansi/LICENSE"
      ],
      "sha256": "a3005e66c096be0df65ecba459ee8cd1aabb39970615eed0e73f5b7cb8ae5672",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/superstruct/package.json": {
      "license": "MIT",
      "name": "superstruct",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/superstruct/LICENSE"
      ],
      "sha256": "30e83d602cf8471740fa94fe38c3bedbdf400aabe94160154557f558c5ee7fd8",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/tar/package.json": {
      "license": "BlueOak-1.0.0",
      "name": "tar",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/tar/LICENSE"
      ],
      "sha256": "d118552ac665732fb7a1bb0e08dba0ab1d21a0edbc0eb3601a094ef821198f21",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/text-table/package.json": {
      "license": "MIT",
      "name": "text-table",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/text-table/LICENSE"
      ],
      "sha256": "b4175217ffbbbc394e1914b32227416464e8d04ee3fff315c7dea0fc1c478961",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/timers-browserify/package.json": {
      "license": "MIT",
      "name": "timers-browserify",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/timers-browserify/LICENSE"
      ],
      "sha256": "8ca98b86666056d8a2096f7cbf391e67f7db6564a0869736bee34c9503558abd",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/tty-browserify/package.json": {
      "license": "MIT",
      "name": "tty-browserify",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/tty-browserify/LICENSE"
      ],
      "sha256": "14c14e5fb2bb63d8d59ac7a356cd0744e47f7901571835d931e2c92be79b60cb",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/unistore/package.json": {
      "license": "MIT",
      "name": "unistore",
      "notice_refs": [
        "licenses/reference/unistore-README.md.txt",
        "licenses/reference/unistore-MIT-terms.txt"
      ],
      "sha256": "e87b50fb2eb256e59fa61d95ff0fb4e93529edcbc261eafca05d8b78abd8d612",
      "source_basis": "Fixed official source/tarball supplement declared in licenses/reference/SOURCE_NOTICES.json; expression/notice scope preserved.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/util/package.json": {
      "license": "MIT",
      "name": "util",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/util/LICENSE"
      ],
      "sha256": "7572a0bd8de81dec020332fdb8ef9648775e778d45c8a4412f23af72fcbe1112",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/vm-browserify/package.json": {
      "license": "MIT",
      "name": "vm-browserify",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/vm-browserify/LICENSE"
      ],
      "sha256": "962e75a1f1639a60d74753649cc10f087ece7154f97a087e7930011debed6cb1",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/watchpack/package.json": {
      "license": "MIT",
      "name": "watchpack",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/watchpack/LICENSE"
      ],
      "sha256": "e906fea196fd89a89690c506e75ea3daea427ca00f64daeda491f95bd829335c",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/webpack-sources3/package.json": {
      "license": "MIT",
      "name": "webpack-sources",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/webpack-sources3/LICENSE"
      ],
      "sha256": "85853d186a1b3d943c185b0e45c51ea65116505cd62818f975f0494bc607104f",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/write-file-atomic/package.json": {
      "license": "ISC",
      "name": "write-file-atomic",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/write-file-atomic/LICENSE"
      ],
      "sha256": "4f98c17cb9cdd527588d48a536808594d323e4afc459c53e3eeb7961bd1eff70",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/ws/package.json": {
      "license": "MIT",
      "name": "ws",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/ws/LICENSE"
      ],
      "sha256": "40f71b1fd703434f93cf2d4f34a68ccdd884de8b70b5ab05b2ee197726e19b26",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/zod-validation-error/package.json": {
      "license": "MIT",
      "name": "zod-validation-error",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/zod-validation-error/LICENSE"
      ],
      "sha256": "0e378c7dede8f0e0077e25bf6895728c3fb855350c7315fa5175f783fe0a8d14",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/zod/package.json": {
      "license": "MIT",
      "name": "zod",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/zod/LICENSE"
      ],
      "sha256": "53c461e9ad93deaa50e7d07a81b6b772dfb0a99f8a5b30414deb895fa95de30c",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": null
    },
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/package.json": {
      "license": "MIT",
      "name": "next",
      "notice_refs": [
        "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/license.md"
      ],
      "sha256": "c4ef8a1c8da320a162485f96d35df7745da90e6becc36882d6e81be9df66eec1",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": "16.3.4"
    },
    "node_modules/.pnpm/picocolors@1.1.1/node_modules/picocolors/package.json": {
      "license": "ISC",
      "name": "picocolors",
      "notice_refs": [
        "node_modules/.pnpm/picocolors@1.1.1/node_modules/picocolors/LICENSE"
      ],
      "sha256": "ef2ac226c4811d312e12c64214c453878653e834482125ae475f27cea60de737",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": "1.1.1"
    },
    "node_modules/.pnpm/postcss@8.5.23/node_modules/postcss/package.json": {
      "license": "MIT",
      "name": "postcss",
      "notice_refs": [
        "node_modules/.pnpm/postcss@8.5.23/node_modules/postcss/LICENSE"
      ],
      "sha256": "df43e035f054943454c49556c89cea8a288d91f5d4b1a6f4e56d48b85c7ebd75",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": "8.5.23"
    },
    "node_modules/.pnpm/react-dom@19.2.8_react@19.2.8/node_modules/react-dom/package.json": {
      "license": "MIT",
      "name": "react-dom",
      "notice_refs": [
        "node_modules/.pnpm/react-dom@19.2.8_react@19.2.8/node_modules/react-dom/LICENSE"
      ],
      "sha256": "7c82f0967114d4b5d38b1f07e543af87c9919e8d3243f9604004819319cb2200",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": "19.2.8"
    },
    "node_modules/.pnpm/react@19.2.8/node_modules/react/package.json": {
      "license": "MIT",
      "name": "react",
      "notice_refs": [
        "node_modules/.pnpm/react@19.2.8/node_modules/react/LICENSE"
      ],
      "sha256": "e40c3ed9b633c9ecf188e09a4be309780a1ea0a9c068278a30ff6a27e7b2dbb6",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": "19.2.8"
    },
    "node_modules/.pnpm/source-map-js@1.2.1/node_modules/source-map-js/package.json": {
      "license": "BSD-3-Clause",
      "name": "source-map-js",
      "notice_refs": [
        "node_modules/.pnpm/source-map-js@1.2.1/node_modules/source-map-js/LICENSE"
      ],
      "sha256": "58f8e94a461a5e461e0cb2bb2be681cbd89b9181989fd26fceb858e67711469a",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": "1.2.1"
    },
    "node_modules/.pnpm/styled-jsx@5.1.6_react@19.2.8/node_modules/styled-jsx/package.json": {
      "license": "MIT",
      "name": "styled-jsx",
      "notice_refs": [
        "node_modules/.pnpm/styled-jsx@5.1.6_react@19.2.8/node_modules/styled-jsx/license.md"
      ],
      "sha256": "0a11b14d7bbeab650019db18883a14cb67daa3615f498a5a625d45828be7eef9",
      "source_basis": "Exact original installed-package legal text; frozen package/lock/manifest identity. No version inferred for framework bundles.",
      "version": "5.1.6"
    }
  },
  "schema": "elite-reference-runtime-notices/v1",
  "source_scope": "Windows local reference; exact installed Next16.3.4 and traced runtime. License evidence/notice reproduction only; no production acceptance.",
  "texts": {
    "licenses/reference/babel-core-LICENSE.txt": "117da2af0d4ce0fe1c8e19b5cff9dcd806adf973d328d27b11d4448c4ff24f76",
    "licenses/reference/babel-packages-inline-notices.txt": "d08e6b38f8be371bbd045b3defa253d451207a5936aed6dee5d62afdd62e9582",
    "licenses/reference/client-only-MIT-terms.txt": "508a77d2e7b51d98adeed32648ad124b7b30241a8e70b2e72c99f92d8e5874d1",
    "licenses/reference/client-only-package.json.txt": "4d6342705767832f299b9a59c28e4275bcf02db19472732f93f67d979441df8f",
    "licenses/reference/constants-browserify-README.md.txt": "b06569bf170a015c56cbad9b45756519d0b8a682770fe30d66bc1e235e26196a",
    "licenses/reference/data-uri-to-buffer-README.md.txt": "6a14de4d552c0c88bdbbaf57f72f625835685eea97de5047f87a463089b8e9f6",
    "licenses/reference/dotenv-LICENSE.txt": "74b629b24865e1e83c5277ee84590b7937644d6fd959d0c7bdce758676cd2ced",
    "licenses/reference/dotenv-expand-LICENSE.txt": "3726b9470c3a6b54e1ebcc1d802d37089b5a5fc97273b52b2cee578a4421ec45",
    "licenses/reference/edge-runtime-cookies-LICENSE.md.txt": "e4f76a7a19ef2989dd79339bd3abf8afcf1d6f065e5a10c76c19415ffd727eb3",
    "licenses/reference/edge-runtime-ponyfill-LICENSE.md.txt": "e4f76a7a19ef2989dd79339bd3abf8afcf1d6f065e5a10c76c19415ffd727eb3",
    "licenses/reference/edge-runtime-primitives-LICENSE.md.txt": "e4f76a7a19ef2989dd79339bd3abf8afcf1d6f065e5a10c76c19415ffd727eb3",
    "licenses/reference/http-proxy-agent-README.md.txt": "ca0b65367e78e9c255a1c9883eb7ec30aae491d76e1c6815b865563106d21dff",
    "licenses/reference/https-proxy-agent-README.md.txt": "32f0856d2c43df7d05cca960fdee84e1e38ab545bd7b2186433dfa41aa90a712",
    "licenses/reference/next-MIT.txt": "ee765244e2d59f5234d474f62e0766fa0c8b99af967fdd4c0cb8dcb0c76ea224",
    "licenses/reference/punycode-LICENSE-MIT.txt.txt": "483acb265f182907d1caf6cff9c16c96f31325ed23792832cc5d8b12d5f88c8a",
    "licenses/reference/setimmediate-LICENSE.txt.txt": "c4b4ad3a5746f1f5249a6dd90396ec519264e1bb02e01e48a6522c48a3a97cb4",
    "licenses/reference/string-hash-README.md.txt": "62a7a3b8dba4cc818974ab2c00cb7fb7cb20afd7d3c73ed049d7f2eb5621a46e",
    "licenses/reference/unistore-MIT-terms.txt": "15c8d6c00935406bf494329437b74317b2926adcaa479b47bf73f6a060d7244d",
    "licenses/reference/unistore-README.md.txt": "9e24224113e53a805fe76ed8084e13697a070061191f627c2053c0ea67f5c364",
    "node_modules/.pnpm/@jridgewell+sourcemap-codec@1.5.5/node_modules/@jridgewell/sourcemap-codec/LICENSE": "769d154fbde32a915af110b1123650bc79f4cbe675acc66e005265bf069c6c6c",
    "node_modules/.pnpm/@oxc-project+types@0.146.0/node_modules/@oxc-project/types/LICENSE": "95ced5ecf1133fbf41d409b5555c86c344f83f3b019926057ddbc07cfdcc27b3",
    "node_modules/.pnpm/@rolldown+pluginutils@1.0.1/node_modules/@rolldown/pluginutils/LICENSE": "e1919b3b98b8bf6c1b99fcc60a0ed46027c97ac947432e1278446c9606f50161",
    "node_modules/.pnpm/@standard-schema+spec@1.1.0/node_modules/@standard-schema/spec/LICENSE": "653b779005a3a4d64a7288c940f7b9a0e8f0b1e0375f6aa6af9473caf131e564",
    "node_modules/.pnpm/@swc+helpers@0.5.23/node_modules/@swc/helpers/LICENSE": "d96eba1f1881334d50ffc144276dfa98e147a415e7c591ecdfcb81680361dc91",
    "node_modules/.pnpm/@types+chai@5.2.3/node_modules/@types/chai/LICENSE": "c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383",
    "node_modules/.pnpm/@types+deep-eql@4.0.2/node_modules/@types/deep-eql/LICENSE": "c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383",
    "node_modules/.pnpm/@types+estree@1.0.9/node_modules/@types/estree/LICENSE": "c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383",
    "node_modules/.pnpm/@types+node@26.2.0/node_modules/@types/node/LICENSE": "c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383",
    "node_modules/.pnpm/@types+react-dom@19.2.5_@types+react@19.2.18/node_modules/@types/react-dom/LICENSE": "c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383",
    "node_modules/.pnpm/@types+react@19.2.18/node_modules/@types/react/LICENSE": "c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383",
    "node_modules/.pnpm/@typescript+typescript-win32-x64@7.0.2/node_modules/@typescript/typescript-win32-x64/LICENSE": "a7d00bfd54525bc694b6e32f64c7ebcf5e6b7ae3657be5cc12767bce74654a47",
    "node_modules/.pnpm/@typescript+typescript-win32-x64@7.0.2/node_modules/@typescript/typescript-win32-x64/NOTICE.txt": "f5c708b59114507b8b27b48181b6883d106bbca0c1634bbee45b5e344237b66b",
    "node_modules/.pnpm/@vitest+expect@4.1.11/node_modules/@vitest/expect/LICENSE": "04575fc5bfae19a9b63100f8691f901a2e15a399011ac02db435a9ef27295551",
    "node_modules/.pnpm/@vitest+mocker@4.1.11_vite@8.2.2_@types+node@26.2.0_/node_modules/@vitest/mocker/LICENSE": "04575fc5bfae19a9b63100f8691f901a2e15a399011ac02db435a9ef27295551",
    "node_modules/.pnpm/@vitest+pretty-format@4.1.11/node_modules/@vitest/pretty-format/LICENSE": "04575fc5bfae19a9b63100f8691f901a2e15a399011ac02db435a9ef27295551",
    "node_modules/.pnpm/@vitest+runner@4.1.11/node_modules/@vitest/runner/LICENSE": "04575fc5bfae19a9b63100f8691f901a2e15a399011ac02db435a9ef27295551",
    "node_modules/.pnpm/@vitest+snapshot@4.1.11/node_modules/@vitest/snapshot/LICENSE": "04575fc5bfae19a9b63100f8691f901a2e15a399011ac02db435a9ef27295551",
    "node_modules/.pnpm/@vitest+spy@4.1.11/node_modules/@vitest/spy/LICENSE": "04575fc5bfae19a9b63100f8691f901a2e15a399011ac02db435a9ef27295551",
    "node_modules/.pnpm/@vitest+utils@4.1.11/node_modules/@vitest/utils/LICENSE": "04575fc5bfae19a9b63100f8691f901a2e15a399011ac02db435a9ef27295551",
    "node_modules/.pnpm/assertion-error@2.0.1/node_modules/assertion-error/LICENSE": "2130216d5ab4c02134f8247f32999feda0abe1eaf028218baeac445c19ce4cea",
    "node_modules/.pnpm/baseline-browser-mapping@2.11.18/node_modules/baseline-browser-mapping/LICENSE.txt": "c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4",
    "node_modules/.pnpm/caniuse-lite@1.0.30001809/node_modules/caniuse-lite/LICENSE": "fd3a263fe19ed8faa9068b43abaebafc02c77897b0c6fc09abc04bb592e5f16e",
    "node_modules/.pnpm/chai@6.2.2/node_modules/chai/LICENSE": "b181da80336ff9dd1043fc8be1a764d7382363433319aa872e4d2cb5ce2a3066",
    "node_modules/.pnpm/convert-source-map@2.0.0/node_modules/convert-source-map/LICENSE": "1fa6ee8bb95a81ae3d73a5bd074a3ac380ffec13697051063ca1a601921b91db",
    "node_modules/.pnpm/csp_evaluator@1.1.8/node_modules/csp_evaluator/LICENSE": "58d1e17ffe5109a7ae296caafcadfdbe6a7d176f0bc4ab01e12a689b0499d8bd",
    "node_modules/.pnpm/csstype@3.2.3/node_modules/csstype/LICENSE": "11d55bd4541c75ee7879547ac49089c489163dae49551389713c3d026cab383e",
    "node_modules/.pnpm/detect-libc@2.1.2/node_modules/detect-libc/LICENSE": "b40930bbcf80744c86c46a12bc9da056641d722716c378f5659b9e555ef833e1",
    "node_modules/.pnpm/es-module-lexer@2.3.2/node_modules/es-module-lexer/LICENSE": "8a4b6c44eebfb026d23719a348145a661a555568dbfdc11618ff2d0dd9306b00",
    "node_modules/.pnpm/estree-walker@3.0.3/node_modules/estree-walker/LICENSE": "8a6dcbabe7179f9c8489c08a7d874d0f1e093ae7449e03756ccf87cb9d0e296e",
    "node_modules/.pnpm/expect-type@1.4.0/node_modules/expect-type/LICENSE": "7c6cc83c84eaa249a85bf12fe3eedd831bb50c696e8890b5248223ae68c7b408",
    "node_modules/.pnpm/fdir@6.5.0_picomatch@4.0.5/node_modules/fdir/LICENSE": "9a39f2aadab11a3697edd668ff2d8ad885b649737b7ab4d3bf12b34e5ada0c86",
    "node_modules/.pnpm/jose@6.2.10/node_modules/jose/LICENSE.md": "8078b0829d6c3e9ecde2f003e1966e4cad3efc6e7a640efbab0313cc166b1af1",
    "node_modules/.pnpm/lightningcss-win32-x64-msvc@1.33.0/node_modules/lightningcss-win32-x64-msvc/LICENSE": "5eba353fe5076ac3432177f8ab1cf75e3afcd0584251e37c3bfead5f447d040e",
    "node_modules/.pnpm/lightningcss@1.33.0/node_modules/lightningcss/LICENSE": "5eba353fe5076ac3432177f8ab1cf75e3afcd0584251e37c3bfead5f447d040e",
    "node_modules/.pnpm/magic-string@0.30.21/node_modules/magic-string/LICENSE": "1cbe51b907662f6cb1492b16c359384a595180bf0e4d101603ed525e75c4e484",
    "node_modules/.pnpm/nanoid@3.3.18/node_modules/nanoid/LICENSE": "da4db1480d9beea3483a2eda5c53b22238d0827d57da162b48f122e04d2d9987",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@babel/runtime/LICENSE": "117da2af0d4ce0fe1c8e19b5cff9dcd806adf973d328d27b11d4448c4ff24f76",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@edge-runtime/primitives/abort-controller.js.LEGAL.txt": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@edge-runtime/primitives/console.js.LEGAL.txt": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@edge-runtime/primitives/crypto.js.LEGAL.txt": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@edge-runtime/primitives/events.js.LEGAL.txt": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@edge-runtime/primitives/fetch.js.LEGAL.txt": "da2a4ed932064e8eace052f31df1c34d8222a9241a38f8cbb97536a04a9936d9",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@edge-runtime/primitives/index.js.LEGAL.txt": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@edge-runtime/primitives/load.js.LEGAL.txt": "da2a4ed932064e8eace052f31df1c34d8222a9241a38f8cbb97536a04a9936d9",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@edge-runtime/primitives/stream.js.LEGAL.txt": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@edge-runtime/primitives/timers.js.LEGAL.txt": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@edge-runtime/primitives/url.js.LEGAL.txt": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@hapi/accept/LICENSE": "ff96bbd6a7408c166f16dbed7f2c87a8061e5a21289ec8f2d3069f8d3b15fc0a",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@mswjs/interceptors/ClientRequest/LICENSE": "bb3a5cb8fb224f632970e31425348f556f7a83b165d5384ea5e142575f5c412e",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@napi-rs/triples/LICENSE": "7b34e99bc43641767fded49f7f03563ffd50c689b1f3b0d50700b4727bc4ae24",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@opentelemetry/api/LICENSE": "c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@vercel/blob/LICENSE": "2b6fc3f4a995e6ad29f4e2ed8881a32979058b8d6557efaa6d86b7e985849f40",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@vercel/detect-agent/LICENSE": "b070d77bfb2c52a1dd6996de0ce5f64c49a0ca55c889b163a963ddf5cb001ee2",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@vercel/nft/LICENSE": "0b31f3a6ab66b272b7de451df39835f310a309a06b1a8aafb3a2fdab20adc166",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@vercel/og/LICENSE": "1f256ecad192880510e84ad60474eab7589218784b9a50bc7ceee34c2b91f1d5",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@vercel/og/satori/LICENSE": "1f256ecad192880510e84ad60474eab7589218784b9a50bc7ceee34c2b91f1d5",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/@vercel/routing-utils/LICENSE": "b070d77bfb2c52a1dd6996de0ce5f64c49a0ca55c889b163a963ddf5cb001ee2",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/acorn/LICENSE": "76a876cf886ff9be2a8b5e2e86514fed06223c8c9f0c1e9ee9606e93841e00b7",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/anser/LICENSE": "6d2abc90a426cf65bbf84f7cedf9487f7ff86a3daf0e130b335a45c112f040d0",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/assert/LICENSE": "6239c6144c31e58cf925c34483606969c555574d64ffa96518ab5d7f45c75d43",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/async-retry/LICENSE": "8cec82f7beeab8459fa6165c973f5cd6f89aceea7a07de0aa8fc30a9dd339a05",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/async-sema/LICENSE": "291ff3e2bd22891edef553930da53ba64adf79d83f0bd1c5b7f6fe9727526b37",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/babel/LICENSE": "117da2af0d4ce0fe1c8e19b5cff9dcd806adf973d328d27b11d4448c4ff24f76",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/browserify-zlib/LICENSE": "ddd770379e07bf0574dfaa4485be80a23b3248b36d09f33ec79276c09b829daf",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/browserslist/LICENSE": "f25bf9bf3ae8984bcd43bf7fb8f78e7eec8d577081fb8d0989cfa7c67ecebb8e",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/buffer/LICENSE": "06bafa45fdad2579ba0e43b0c9b2c6290287c99c4203c300254a462b38a307f6",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/busboy/LICENSE": "d06b5d27bbbbe22c36b1fd88406b1208876e2d37d795f5b8eaed951a459a3111",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/bytes/LICENSE": "e3b44af066615de2ea48d18d852d0762f18c0b2efcea714fa48a6f729d405b85",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/ci-info/LICENSE": "25bd22427e27b0d3f1c3da2670bb44b622bcced22ebdf4777d5488ac129370a1",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/cli-select/LICENSE": "a8e260b5b60eb27746ff2c5967489e6878c5a6ba99400cbec84c89a9b49a46ca",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/commander/LICENSE": "04512a63dce4d2d506ad612dc0bd7681ccf6e3655f7b6eaef7dfac8323d1ec0b",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/comment-json/LICENSE": "9dbd86998d3cdafc1d8c298aea636146568d05d429b81386e60cfce11e6cedb6",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/compression/LICENSE": "23d466b1b134b415b66fa50c6526b4cf3e7b9258554da88d3abb371721e7ce68",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/conf/LICENSE": "48da2f39e100d4085767e94966b43f4fa95ff6a0698fba57ed460914e35f94a0",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/content-disposition/LICENSE": "bd47ce7b88c7759630d1e2b9fcfa170a0f1fde522be09e13fb1581a79d090400",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/content-type/LICENSE": "257aed98914108e91a337912727b6a802eef218248507f74b76faffaff517a38",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/cookie/LICENSE": "c02110eedc16c7114f1a9bdc026c65626ce1d9c7e27fd51a8e0feee8a48a6858",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/cross-spawn/LICENSE": "aaa78451b6fecd1b9c4594c796c133c0e90cad100372ff8bc6de615e9ef9adf1",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/crypto-browserify/LICENSE": "6134c69bc22c8289252e70de3af20bd67071233459055be74d83acfcc4865e7e",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/debug/LICENSE": "98c970de440dcfc77471610aec2377c9d9b0db2b3be6d1add524a586e1d7f422",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/devalue/LICENSE": "e399db7e952e45cae84f1b53d72ba8dbad2daaa7ab8cd1f0ab72fea5ab73b225",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/domain-browser/LICENSE": "fc06db8dd5e7840db84ab328d9f58bdb91a56fbbf045bb977c11fcd07e2d8e3e",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/edge-runtime/LICENSE": "e4f76a7a19ef2989dd79339bd3abf8afcf1d6f065e5a10c76c19415ffd727eb3",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/events/LICENSE": "631987b7616a325a5b97566c232418481ddf7dbb5ecadefb991e791876cc2599",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/find-up/LICENSE": "48da2f39e100d4085767e94966b43f4fa95ff6a0698fba57ed460914e35f94a0",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/fresh/LICENSE": "a0ec0dab16b3666f24950f86d257930ac2ad475557b4bfe245620e0817d8a45d",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/glob/LICENSE": "6236fa0b88a4a0cce3dda0367979491b2052b3c8d6b1c10b3668de083e86a7f0",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/gzip-size/LICENSE": "48da2f39e100d4085767e94966b43f4fa95ff6a0698fba57ed460914e35f94a0",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/https-browserify/LICENSE": "ff151c00207c908581639851dd8504ce4255be0650b2b236edec2aa90342b0cd",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/httpxy/LICENSE": "4249bd020d1b5c503d20a196c1b686b3bdebcd1c7ac06909c546d9624c64f79b",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/icss-utils/LICENSE": "b2b54c738525e98cc42d4a1c239c9d8421507816acce98950e81d0943daae68e",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/ignore-loader/LICENSE": "62ba766e9b50f9b175b0229f4ea03e93730c9180408acec50dd617e9ae1fa059",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/image-size/LICENSE": "3109783dad7527b490aabc49c3cb148573c19b41b843e70cf9388a93f3d7fccd",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/ipaddr.js/LICENSE": "62568a2d1337b77171ecca9db10579163446de1ba6151678e81f06cdc199971b",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/is-animated/LICENSE": "e08c030e1274fadc495c71e0489612f0c74673b69ae8e41460526c233a94a728",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/is-docker/LICENSE": "48da2f39e100d4085767e94966b43f4fa95ff6a0698fba57ed460914e35f94a0",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/is-wsl/LICENSE": "48da2f39e100d4085767e94966b43f4fa95ff6a0698fba57ed460914e35f94a0",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/jest-worker/LICENSE": "52412d7bc7ce4157ea628bbaacb8829e0a9cb3c58f57f99176126bc8cf2bfc85",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/json5/LICENSE": "53e59feb13058722d977c699eb0407c7bce2f93c949b681bbd2ff31698535927",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/jsonwebtoken/LICENSE": "2144eb6894cb440fde6b2b3aaae3b617c8f1dc9bb19813079dcf62ec6c517042",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/loader-runner/LICENSE": "498d39c83164cf8de195090f5564fc0308ab271344ca2f95484a7f1020952d00",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/loader-utils2/LICENSE": "9068a8782d2fb4c6e432cfa25334efa56f722822180570802bf86e71b6003b1e",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/loader-utils3/LICENSE": "9068a8782d2fb4c6e432cfa25334efa56f722822180570802bf86e71b6003b1e",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/lodash.curry/LICENSE": "ffd8b33b354585f4ce119f19c53728281e48a97b074491eb6bf6d5c5ff305272",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/lru-cache/LICENSE": "451ec07eeb9c4e1b86de9abdaa426462a8be48f887ec7421cf0bbb9c769555ab",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/mini-css-extract-plugin/LICENSE": "9068a8782d2fb4c6e432cfa25334efa56f722822180570802bf86e71b6003b1e",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/nanoid/LICENSE": "da4db1480d9beea3483a2eda5c53b22238d0827d57da162b48f122e04d2d9987",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/native-url/LICENSE": "cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/neo-async/LICENSE": "811238ba7d85f6fe6b820703a32f92705bcf77bc352ddc3476783491c64a129a",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/node-html-parser/LICENSE": "ea75e9321945b7f0bc4644da1e404ad76ce55e9bc05cae3b2f5f2d973aede83b",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/ora/LICENSE": "5c932d88256b4ab958f64a856fa48e8bd1f55bc1d96b8149c65689e0c61789d3",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/os-browserify/LICENSE": "d25d1d6d28c35cb6f358e2833e405c4e53fec2fa24d156323ddea5cd438d3407",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/p-limit/LICENSE": "5c932d88256b4ab958f64a856fa48e8bd1f55bc1d96b8149c65689e0c61789d3",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/p-queue/LICENSE": "5c932d88256b4ab958f64a856fa48e8bd1f55bc1d96b8149c65689e0c61789d3",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/path-browserify/LICENSE": "a22b9d5763f574e5db347c30acc0b33eaf4846767c03d2e27d012e864e79a824",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/path-to-regexp/LICENSE": "4eeb3271453a891df609e5a9f4ee79a68307f730c13417a3bfeffa604ac8cf25",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/picomatch/LICENSE": "d0cd141b0c322fded5dfad1d4645bb2fedfc05b7321fe1009469638190d59ef9",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/postcss-flexbugs-fixes/LICENSE": "f26c92c1c16e78c6b1f7474792f7195b7abc851caccb4a73dcc5b3c31efc7321",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/postcss-modules-extract-imports/LICENSE": "9e534cfb7f7ae1300e16fe1ee7d203fdfe707dc932429c37246f747c3699ee14",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/postcss-modules-local-by-default/LICENSE": "dd40f37957ddcb6a6157f40b677d6f43697583f956144fbce25443e35630453e",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/postcss-modules-scope/LICENSE": "7774c8c4a3f750b294b416818062201a404fc2d6c8756d8b89a1688e24b5d5a7",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/postcss-modules-values/LICENSE": "655287aefbf8e292057318c385d77767b1dd9819c87709fd26730b7090efe4d9",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/postcss-preset-env/LICENSE": "597756adcb51f243ef4fb386920377f61d012ace0904364e1a8ee9aaec6afc84",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/postcss-safe-parser/LICENSE": "c4630ac8b89cb317ac5bdd60ac5e4e185eab9bd5151a0c7b3afa41aa83d7ec9b",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/postcss-scss/LICENSE": "c4630ac8b89cb317ac5bdd60ac5e4e185eab9bd5151a0c7b3afa41aa83d7ec9b",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/postcss-value-parser/LICENSE": "3687447039151857a6ba378db062172c7f33d4aa70a615c87a43a9c50e990485",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/process/LICENSE": "59a400d04c5078579acc27ddd6452c1bdf763f9506e01364700935fbb1a7c91b",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/querystring-es3/LICENSE": "cb72d9714ddc21e758d63f423be0caddf909d23ccbb10a2f5201a870818e4f57",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/raw-body/LICENSE": "c8e6bca7230689d536a3bd7158f66e9c4f89f95d0748743a0370ac229e9023ad",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/react-dom-experimental/LICENSE": "da6d3703ed11cbe42bd212c725957c98da23cbff1998c05fa4b3d976d1a58e93",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/react-dom/LICENSE": "da6d3703ed11cbe42bd212c725957c98da23cbff1998c05fa4b3d976d1a58e93",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/react-experimental/LICENSE": "da6d3703ed11cbe42bd212c725957c98da23cbff1998c05fa4b3d976d1a58e93",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/react-is/LICENSE": "da6d3703ed11cbe42bd212c725957c98da23cbff1998c05fa4b3d976d1a58e93",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/react-refresh/LICENSE": "52412d7bc7ce4157ea628bbaacb8829e0a9cb3c58f57f99176126bc8cf2bfc85",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/react-server-dom-turbopack-experimental/LICENSE": "da6d3703ed11cbe42bd212c725957c98da23cbff1998c05fa4b3d976d1a58e93",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/react-server-dom-turbopack/LICENSE": "da6d3703ed11cbe42bd212c725957c98da23cbff1998c05fa4b3d976d1a58e93",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/react-server-dom-webpack-experimental/LICENSE": "da6d3703ed11cbe42bd212c725957c98da23cbff1998c05fa4b3d976d1a58e93",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/react-server-dom-webpack/LICENSE": "da6d3703ed11cbe42bd212c725957c98da23cbff1998c05fa4b3d976d1a58e93",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/react/LICENSE": "da6d3703ed11cbe42bd212c725957c98da23cbff1998c05fa4b3d976d1a58e93",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/regenerator-runtime/LICENSE": "51887a3d47051ac2fce1210562e5b9fe0830a8a8fabeb272c2d586eeb18a05fd",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/safe-stable-stringify/LICENSE": "3a6d6de49d284a36c5d271ef0e62eca1838cc460dd7630e3f2c347c8f5b2232a",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/sass-loader/LICENSE": "2af12ca8b612a22e6659ae68e0e59ee7ac792718c45524a989b512c49c5f5879",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/scheduler-experimental/LICENSE": "da6d3703ed11cbe42bd212c725957c98da23cbff1998c05fa4b3d976d1a58e93",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/scheduler/LICENSE": "da6d3703ed11cbe42bd212c725957c98da23cbff1998c05fa4b3d976d1a58e93",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/schema-utils2/LICENSE": "9068a8782d2fb4c6e432cfa25334efa56f722822180570802bf86e71b6003b1e",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/schema-utils3/LICENSE": "9068a8782d2fb4c6e432cfa25334efa56f722822180570802bf86e71b6003b1e",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/semver/LICENSE": "4ec3d4c66cd87f5c8d8ad911b10f99bf27cb00cdfcff82621956e379186b016b",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/send/LICENSE": "25a328069fe771f8ed5b6f983ed4b0e6c84b3312ac0f69b28c3d52dc277962c9",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/serve-handler/LICENSE": "be8ba42db3b7c3dbb064bcdc1bbbb191c0a39c8851c2a53564e97c7eea301bca",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/shell-quote/LICENSE": "8bb16db1b047019e4395965f2cf3611b06c34bf86dc2d0210b3c3f91b53c21fe",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/source-map/LICENSE": "6cb0631f71c7749763fd3dd1d5bee52dd1070ec17f2edc1710079ad070bd2fbd",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/source-map08/LICENSE": "6cb0631f71c7749763fd3dd1d5bee52dd1070ec17f2edc1710079ad070bd2fbd",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/stacktrace-parser/LICENSE": "e1898814043cb5500483d9c8517edda824b352f696d35e5043fcfb537c221531",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/stream-browserify/LICENSE": "96e5193183476a03da6667d217ce14110191d26ee105338086e94339edf40d5a",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/stream-http/LICENSE": "a0e6357a5e8ea65827addeb383e0948a1874d2f46bc7feaf6349b7a376ed6e98",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/string_decoder/LICENSE": "11f2aafb37d06b3ee5bdaf06e9811141d0da05263c316f3d627f45c20d43261b",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/strip-ansi/LICENSE": "48da2f39e100d4085767e94966b43f4fa95ff6a0698fba57ed460914e35f94a0",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/superstruct/LICENSE": "0f644bc01d85e4b1958bdfe5a851704bd886404c78188d4d6fba43dc8293bc91",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/tar/LICENSE": "8a1af140fdfbf5afd3df27f7e662f989c5b963a300020dfafce42033cae9e004",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/terser/LICENSE": "901d0a7fcaccaae9917d165039c0473c354e5da2e0e94140dde67aad9ffb19da",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/text-table/LICENSE": "435a6722c786b0a56fbe7387028f1d9d3f3a2d0fb615bb8fee118727c3f59b7b",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/timers-browserify/LICENSE": "d5f14c3258420dfe5a3b641a143d6e6dd90eabb5962244d937e25699c3a45ec9",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/tty-browserify/LICENSE": "435a6722c786b0a56fbe7387028f1d9d3f3a2d0fb615bb8fee118727c3f59b7b",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/ua-parser-js/LICENSE": "fb9548e36708f09ae0cc1047ebd8153928bfe5bed67e4051f0ace552bdcbf0b3",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/util/LICENSE": "6239c6144c31e58cf925c34483606969c555574d64ffa96518ab5d7f45c75d43",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/vm-browserify/LICENSE": "3a62887f73cd035ac1068bc0127ce94291f7a03843391f56df3520e2ca2037f8",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/watchpack/LICENSE": "9068a8782d2fb4c6e432cfa25334efa56f722822180570802bf86e71b6003b1e",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/web-vitals-attribution/LICENSE": "bfdeded4040e05da31ca9b6239dc83bd23fa26ac8db87342a7ca4363f68916ff",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/web-vitals/LICENSE": "bfdeded4040e05da31ca9b6239dc83bd23fa26ac8db87342a7ca4363f68916ff",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/webpack-sources1/LICENSE": "214d0ac13edb475ec426df1979d96fd46b443debc466c0fb592e32152ae5a940",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/webpack-sources3/LICENSE": "214d0ac13edb475ec426df1979d96fd46b443debc466c0fb592e32152ae5a940",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/webpack/LICENSE": "9068a8782d2fb4c6e432cfa25334efa56f722822180570802bf86e71b6003b1e",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/write-file-atomic/LICENSE": "ea7f376fe7a1fc28572b83ac8f806d92effb31852b9981bc9ba9d5266caa6b28",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/ws/LICENSE": "dcd0a622a3a2c94fe3aa6238b76363f87df442afe788cbad966afe943ed62281",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/zod-validation-error/LICENSE": "03cc9a630489a2544b799548a0e8b134f382fe09fa06dc22b339091b28b8f4a8",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/dist/compiled/zod/LICENSE": "3f1189b28e3866e0d979968d466b78f813f76827cfdca1fbb124cc0a5c8841f8",
    "node_modules/.pnpm/next@16.3.4_@types+node@26._8e51c0b8893f067905547ccb40d36b56/node_modules/next/license.md": "ee765244e2d59f5234d474f62e0766fa0c8b99af967fdd4c0cb8dcb0c76ea224",
    "node_modules/.pnpm/oauth4webapi@3.8.7/node_modules/oauth4webapi/LICENSE.md": "81e7b4b0196e8caa43cf243131512e5e4ab04fc312623335e783abae6ff46d79",
    "node_modules/.pnpm/obug@2.1.4/node_modules/obug/LICENSE": "ee48679d379ca6b4493d5e231094d85f818dd9be40b1a4f234fb7ff657ee35d9",
    "node_modules/.pnpm/openid-client@6.8.5/node_modules/openid-client/LICENSE.md": "14c5cc0dc21f44add6d88a5621c65813d40bf382550c978430096db3df0cf68c",
    "node_modules/.pnpm/pathe@2.0.3/node_modules/pathe/LICENSE": "52e92576851154bad7737e90cc72818936f43665cb0e3f7428ed8edc8cc5709b",
    "node_modules/.pnpm/picocolors@1.1.1/node_modules/picocolors/LICENSE": "6582629e2979466878f6014313dcc2f3756c9616148682227ce3063dde310750",
    "node_modules/.pnpm/picomatch@4.0.5/node_modules/picomatch/LICENSE": "d0cd141b0c322fded5dfad1d4645bb2fedfc05b7321fe1009469638190d59ef9",
    "node_modules/.pnpm/postcss@8.5.23/node_modules/postcss/LICENSE": "5be1f3465bba68a626777f984878814aaf35e7ef8e9fd314d469bcf887050fb8",
    "node_modules/.pnpm/postcss@8.5.26/node_modules/postcss/LICENSE": "5be1f3465bba68a626777f984878814aaf35e7ef8e9fd314d469bcf887050fb8",
    "node_modules/.pnpm/react-dom@19.2.8_react@19.2.8/node_modules/react-dom/LICENSE": "da6d3703ed11cbe42bd212c725957c98da23cbff1998c05fa4b3d976d1a58e93",
    "node_modules/.pnpm/react@19.2.8/node_modules/react/LICENSE": "da6d3703ed11cbe42bd212c725957c98da23cbff1998c05fa4b3d976d1a58e93",
    "node_modules/.pnpm/rolldown@1.2.5/node_modules/rolldown/LICENSE": "23ecfff35a5a2e80d92142f75228912c3b1abc4b5a8337a821ff4397e2f9f734",
    "node_modules/.pnpm/safevalues@1.2.0/node_modules/safevalues/LICENSE": "cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30",
    "node_modules/.pnpm/scheduler@0.27.0/node_modules/scheduler/LICENSE": "da6d3703ed11cbe42bd212c725957c98da23cbff1998c05fa4b3d976d1a58e93",
    "node_modules/.pnpm/siginfo@2.0.0/node_modules/siginfo/LICENSE": "3bdddf0b9b08aaa2fe4365b803fc958ed32c0597ca0bd727be8d9fcb8b248d18",
    "node_modules/.pnpm/source-map-js@1.2.1/node_modules/source-map-js/LICENSE": "6cb0631f71c7749763fd3dd1d5bee52dd1070ec17f2edc1710079ad070bd2fbd",
    "node_modules/.pnpm/std-env@4.2.0/node_modules/std-env/LICENCE": "a6f36438e46fb911859f3b9c4cad045ba64e1af3d8f4512c60258fa2d7552d28",
    "node_modules/.pnpm/styled-jsx@5.1.6_react@19.2.8/node_modules/styled-jsx/license.md": "24d4a4580df855c345705b3016404f0c9da7e7674dc3b21a3ce7d0e5d0c37e2d",
    "node_modules/.pnpm/tinybench@2.9.0/node_modules/tinybench/LICENSE": "cebc084d54e6dd99e53292ddb4bc1cdb63d9bfcd9ce438fe3aa6eb106d79e2ea",
    "node_modules/.pnpm/tinyexec@1.3.0/node_modules/tinyexec/LICENSE": "f95f668fe64081ddb4153b322e34fdd719b991285ed08177d5ad7133b7988d92",
    "node_modules/.pnpm/tinyglobby@0.2.17/node_modules/tinyglobby/LICENSE": "22c68811e174cbbfb3813d4135918df4f540959c14d872e601d9abe83d3cde8f",
    "node_modules/.pnpm/tinyrainbow@3.1.1/node_modules/tinyrainbow/LICENCE": "cebc084d54e6dd99e53292ddb4bc1cdb63d9bfcd9ce438fe3aa6eb106d79e2ea",
    "node_modules/.pnpm/tslib@2.8.1/node_modules/tslib/LICENSE.txt": "210b19e543130388c68654b7497e967119ce17145f66ab7d85688fbd70f08751",
    "node_modules/.pnpm/typescript@7.0.2/node_modules/typescript/LICENSE": "a7d00bfd54525bc694b6e32f64c7ebcf5e6b7ae3657be5cc12767bce74654a47",
    "node_modules/.pnpm/typescript@7.0.2/node_modules/typescript/NOTICE.txt": "f5c708b59114507b8b27b48181b6883d106bbca0c1634bbee45b5e344237b66b",
    "node_modules/.pnpm/undici-types@8.3.0/node_modules/undici-types/LICENSE": "a6db8096b2707bc0102d256917d4d33f298ba36d8c3f25de067a2b5bb379db27",
    "node_modules/.pnpm/vite@8.2.2_@types+node@26.2.0/node_modules/vite/LICENSE.md": "387dd7baa307083401a27c58c362c30832f5ba1dba84f10cc22c33401523f45c",
    "node_modules/.pnpm/vitest@4.1.11_@types+node@2_295149386df9925e43b8f6e72d9400b5/node_modules/vitest/LICENSE.md": "881d660c26831481b697e39724d4a35c9f86e07b67156d4aeb693a0b39910435",
    "node_modules/.pnpm/why-is-node-running@2.3.0/node_modules/why-is-node-running/LICENSE": "6a134e51aa31496c15a741592fc5d782b2dbf8d8b6b8524e15c8520ae0cd6374",
    "node_modules/.pnpm/zod@4.4.3/node_modules/zod/LICENSE": "3f1189b28e3866e0d979968d466b78f813f76827cfdca1fbb124cc0a5c8841f8"
  }
}
````

### FILE: `licenses/reference/SOURCE_NOTICES.json`

```yaml
block_id: "REFERENCE-NOTICE-BLOCKS:file21:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "18b20c9a287eb5ad12217db27f58384812d2d36a7a74c8355390db94f9bdd47b"
variables: []
secrets_allowed: false
```

````json
{
  "schema": "elite-reference-notice-sources/v1",
  "state": "SOURCE_VERIFIED_RUNTIME_GATE_PENDING",
  "claim": "Notice correspondence; no new dependency code or runtime version. Three adapted notice renderings are explicitly not upstream verbatim files.",
  "sources": {
    "licenses/reference/next-MIT.txt": {
      "repo": "vercel/next.js",
      "commit": "299180d3315c7ebd7b199d2b1a265b5986c5fc7d",
      "path": "license.md",
      "url": "https://api.github.com/repos/vercel/next.js/contents/license.md?ref=299180d3315c7ebd7b199d2b1a265b5986c5fc7d",
      "status": 200,
      "response_path": "fixed-source/300cdf6a66c6b88d2fa4.response.json",
      "response_sha256": "afc0acd6bcb231cf43b2dceec2a33d7b58aa66e255fafb874110f84bde7eff58",
      "local_path": "fixed-source/300cdf6a66c6b88d2fa4.txt",
      "sha256": "ee765244e2d59f5234d474f62e0766fa0c8b99af967fdd4c0cb8dcb0c76ea224",
      "bytes": 1079,
      "git_blob_sha1": "5948ee9bd0de5064423688a2967ab3111c2658ed"
    },
    "licenses/reference/constants-browserify-README.md.txt": {
      "package": "constants-browserify",
      "version": "1.0.0",
      "artifact_sha256": "ea9307fd609ddce9497e3ca6b49e6102d7326bd337f66c10143b7b326f931f59",
      "integrity": "sha512-xFxOwqIzR/e1k1gLiWEophSCMqXcwVHIH7akf7b/vxcUeGunlj3hvZaaqxwHsTgn+IndtkQJgSztIDWeumWJDQ==",
      "member": "README.md",
      "path": "nested-sources/constants-browserify-README.md.txt",
      "sha256": "b06569bf170a015c56cbad9b45756519d0b8a682770fe30d66bc1e235e26196a",
      "bytes": 1668
    },
    "licenses/reference/data-uri-to-buffer-README.md.txt": {
      "package": "data-uri-to-buffer",
      "version": "3.0.1",
      "artifact_sha256": "3787e95764b16659c50f96ab4889f9a354ba3c273620b1f0aea3fdad9fbd4dcc",
      "integrity": "sha512-WboRycPNsVw3B3TL559F7kuBUM4d8CgMEvk6xEJlOp7OBPjt6G7z8WMWlD2rOFZLk6OYfFIUGsCOWzcQH9K2og==",
      "member": "README.md",
      "path": "nested-sources/data-uri-to-buffer-README.md.txt",
      "sha256": "6a14de4d552c0c88bdbbaf57f72f625835685eea97de5047f87a463089b8e9f6",
      "bytes": 2929
    },
    "licenses/reference/http-proxy-agent-README.md.txt": {
      "package": "http-proxy-agent",
      "version": "5.0.0",
      "artifact_sha256": "492fe73645f2dfdb62ecd3360c8215a2253b6f87dc7f5e44149893b009c624db",
      "integrity": "sha512-n2hY8YdoRE1i7r6M0w9DIw5GgZN0G25P8zLCRQ8rjXtTU3vsNFBI/vWK/UIeE6g5MUUz6avwAPXmL6Fy9D/90w==",
      "member": "README.md",
      "path": "nested-sources/http-proxy-agent-README.md.txt",
      "sha256": "ca0b65367e78e9c255a1c9883eb7ec30aae491d76e1c6815b865563106d21dff",
      "bytes": 2527
    },
    "licenses/reference/https-proxy-agent-README.md.txt": {
      "package": "https-proxy-agent",
      "version": "5.0.1",
      "artifact_sha256": "6da16fb44331f2e5d30bd21bf880aa934c1ad4fe7da7187910ef2b2509712019",
      "integrity": "sha512-dFcAjpTQFgoLMzC2VwU+C/CbS7uRL0lWmxDITmqm7C+7F0Odmj6s9l6alZc6AELXhrnggM2CeWSXHGOdX2YtwA==",
      "member": "README.md",
      "path": "nested-sources/https-proxy-agent-README.md.txt",
      "sha256": "32f0856d2c43df7d05cca960fdee84e1e38ab545bd7b2186433dfa41aa90a712",
      "bytes": 4761
    },
    "licenses/reference/punycode-LICENSE-MIT.txt.txt": {
      "package": "punycode",
      "version": "2.1.1",
      "artifact_sha256": "e5ccd079bf02390e34cbb51b75f73a8be42795146a2e2e23cbb204dddb04ff3c",
      "integrity": "sha512-XRsRjdf+j5ml+y/6GKHPZbrF/8p2Yga0JPtdqTIY2Xe5ohJPD9saDJJLPvp9+NSBprVvevdXZybnj2cv8OEd0A==",
      "member": "LICENSE-MIT.txt",
      "path": "nested-sources/punycode-LICENSE-MIT.txt.txt",
      "sha256": "483acb265f182907d1caf6cff9c16c96f31325ed23792832cc5d8b12d5f88c8a",
      "bytes": 1077
    },
    "licenses/reference/setimmediate-LICENSE.txt.txt": {
      "package": "setimmediate",
      "version": "1.0.5",
      "artifact_sha256": "5cb9fc22698364ed42c02d6aa3dc50ffeafa68452ae84699672e3dfd74922c9e",
      "integrity": "sha512-MATJdZp8sLqDl/68LfQmbP8zKPLQNV6BIZoIgrscFDQ+RsvK/BxeDQOgyxKKoh0y/8h3BqVFnCqQ/gd+reiIXA==",
      "member": "LICENSE.txt",
      "path": "nested-sources/setimmediate-LICENSE.txt.txt",
      "sha256": "c4b4ad3a5746f1f5249a6dd90396ec519264e1bb02e01e48a6522c48a3a97cb4",
      "bytes": 1103
    },
    "licenses/reference/string-hash-README.md.txt": {
      "package": "string-hash",
      "version": "1.1.3",
      "artifact_sha256": "264f97224f48229c23c8cac1f045868a4bdc8dd95996e595e448d0be622ba1f4",
      "integrity": "sha512-kJUvRUFK49aub+a7T1nNE66EJbZBMnBgoC1UbCZ5n6bsZKBRga4KgBRTMn/pFkeCZSYtNeSyMxPDM0AXWELk2A==",
      "member": "README.md",
      "path": "nested-sources/string-hash-README.md.txt",
      "sha256": "62a7a3b8dba4cc818974ab2c00cb7fb7cb20afd7d3c73ed049d7f2eb5621a46e",
      "bytes": 841
    },
    "licenses/reference/unistore-README.md.txt": {
      "package": "unistore",
      "version": "3.4.1",
      "artifact_sha256": "894e04593ba98d04f8225eb3f3a9505e232853b27a03ad7d84318a2cf2fbee83",
      "integrity": "sha512-p2Ej8qqrqcD10Ah0ZUKUU/mhRB8pM4q6gzjxq9kZpgxa8dks7oHT8jDP4CqLhoRof3RXOZLKB9EBV1DTzHiJRw==",
      "member": "README.md",
      "path": "nested-sources/unistore-README.md.txt",
      "sha256": "9e24224113e53a805fe76ed8084e13697a070061191f627c2053c0ea67f5c364",
      "bytes": 8474
    },
    "licenses/reference/edge-runtime-cookies-LICENSE.md.txt": {
      "package": "@edge-runtime/cookies",
      "version": "6.0.0",
      "artifact_sha256": "6ed10bffb7f722669cad0ef21ae28c7ced140785a587ddf2aa4cb8b98eadda4a",
      "integrity": "sha512-VVO/8AwC2qVbygLb2IOkX1zWFx2yWIHzFv4D602CTnoRffd/+cdcXqpSydKaedFrk7a1dRYXbWwjzfV/gwZ2Gw==",
      "member": "LICENSE.md",
      "path": "nested-sources/edge-runtime-cookies-LICENSE.md.txt",
      "sha256": "e4f76a7a19ef2989dd79339bd3abf8afcf1d6f065e5a10c76c19415ffd727eb3",
      "bytes": 1079
    },
    "licenses/reference/edge-runtime-ponyfill-LICENSE.md.txt": {
      "package": "@edge-runtime/ponyfill",
      "version": "4.0.0",
      "artifact_sha256": "8c5af3358af8889d2fa65b9bf9432b0f36a88486cce5702b14f18f8986fb55c1",
      "integrity": "sha512-JZdDyYUqweSbhMgSxxL8QuquNuanVw+gMK8RfFq1LHMQHxzXPVvAf9aw6M7MSuyhnCBtFhqtSTr682rf8U5eyA==",
      "member": "LICENSE.md",
      "path": "nested-sources/edge-runtime-ponyfill-LICENSE.md.txt",
      "sha256": "e4f76a7a19ef2989dd79339bd3abf8afcf1d6f065e5a10c76c19415ffd727eb3",
      "bytes": 1079
    },
    "licenses/reference/edge-runtime-primitives-LICENSE.md.txt": {
      "package": "@edge-runtime/primitives",
      "version": "6.0.0",
      "artifact_sha256": "92878e3cfa4d42010c760449210ab4ca6d71435d71db8ee97b7ecbf4fd75c6e2",
      "integrity": "sha512-FqoxaBT+prPBHBwE1WXS1ocnu/VLTQyZ6NMUBAdbP7N2hsFTTxMC/jMu2D/8GAlMQfxeuppcPuCUk/HO3fpIvA==",
      "member": "LICENSE.md",
      "path": "nested-sources/edge-runtime-primitives-LICENSE.md.txt",
      "sha256": "e4f76a7a19ef2989dd79339bd3abf8afcf1d6f065e5a10c76c19415ffd727eb3",
      "bytes": 1079
    },
    "licenses/reference/babel-core-LICENSE.txt": {
      "package": "@babel/core",
      "version": "7.26.10",
      "artifact_sha256": "03de5fbdd939f02afd159797f488d56995250f596a42cea0468b784c8d7fb351",
      "integrity": "sha512-vMqyb7XCDMPvJFFOaT9kxtiRh42GwlZEg1/uIgtZshS5a/8OaduUfCi7kynKgc3Tw/6Uo2D+db9qBttghhmxwQ==",
      "member": "LICENSE",
      "path": "nested-sources/babel-core-LICENSE.txt",
      "sha256": "117da2af0d4ce0fe1c8e19b5cff9dcd806adf973d328d27b11d4448c4ff24f76",
      "bytes": 1106
    },
    "licenses/reference/dotenv-LICENSE.txt": {
      "package": "dotenv",
      "version": "16.3.1",
      "integrity": "sha512-IPzF4w4/Rd94bA9imS68tZBaYyBWSCE47V1RGuMrB94iyTOIEwRmVL2x/4An+6mETpLrKJ5hQkB8W4kFAadeIQ==",
      "path": "env-notices/dotenv-LICENSE.txt",
      "sha256": "74b629b24865e1e83c5277ee84590b7937644d6fd959d0c7bdce758676cd2ced",
      "bytes": 1294
    },
    "licenses/reference/dotenv-expand-LICENSE.txt": {
      "package": "dotenv-expand",
      "version": "10.0.0",
      "integrity": "sha512-GopVGCpVS1UKH75VKHGuQFqS1Gusej0z4FyQkPdwjil2gNIv+LNsqBlboOzpJFZKVT95GkCyWJbBSdFEFUWI2A==",
      "path": "env-notices/dotenv-expand-LICENSE.txt",
      "sha256": "3726b9470c3a6b54e1ebcc1d802d37089b5a5fc97273b52b2cee578a4421ec45",
      "bytes": 1295
    },
    "licenses/reference/unistore-MIT-terms.txt": {
      "basis": "Exact unistore3.4.1 README grants MIT and names Jason Miller; retain original README alongside this declared rendering. Standard permission terms from fixed Next MIT license; no Vercel attribution.",
      "standard_terms_source": {
        "repo": "vercel/next.js",
        "commit": "299180d3315c7ebd7b199d2b1a265b5986c5fc7d",
        "path": "license.md",
        "url": "https://api.github.com/repos/vercel/next.js/contents/license.md?ref=299180d3315c7ebd7b199d2b1a265b5986c5fc7d",
        "status": 200,
        "response_path": "fixed-source/300cdf6a66c6b88d2fa4.response.json",
        "response_sha256": "afc0acd6bcb231cf43b2dceec2a33d7b58aa66e255fafb874110f84bde7eff58",
        "local_path": "fixed-source/300cdf6a66c6b88d2fa4.txt",
        "sha256": "ee765244e2d59f5234d474f62e0766fa0c8b99af967fdd4c0cb8dcb0c76ea224",
        "bytes": 1079,
        "git_blob_sha1": "5948ee9bd0de5064423688a2967ab3111c2658ed"
      },
      "adaptation": "Replace/remove only copyright heading; standard permission paragraphs byte-identical."
    },
    "licenses/reference/client-only-MIT-terms.txt": {
      "basis": "Exact client-only0.0.1 package.json grants MIT without a copyright notice. Standard terms reproduced without inventing an owner/year; actual three files match official Next fixed source. No Vercel/Meta authorship is inferred.",
      "standard_terms_source": {
        "repo": "vercel/next.js",
        "commit": "299180d3315c7ebd7b199d2b1a265b5986c5fc7d",
        "path": "license.md",
        "url": "https://api.github.com/repos/vercel/next.js/contents/license.md?ref=299180d3315c7ebd7b199d2b1a265b5986c5fc7d",
        "status": 200,
        "response_path": "fixed-source/300cdf6a66c6b88d2fa4.response.json",
        "response_sha256": "afc0acd6bcb231cf43b2dceec2a33d7b58aa66e255fafb874110f84bde7eff58",
        "local_path": "fixed-source/300cdf6a66c6b88d2fa4.txt",
        "sha256": "ee765244e2d59f5234d474f62e0766fa0c8b99af967fdd4c0cb8dcb0c76ea224",
        "bytes": 1079,
        "git_blob_sha1": "5948ee9bd0de5064423688a2967ab3111c2658ed"
      },
      "adaptation": "Replace/remove only copyright heading; standard permission paragraphs byte-identical."
    },
    "licenses/reference/client-only-package.json.txt": {
      "repo": "vercel/next.js",
      "commit": "299180d3315c7ebd7b199d2b1a265b5986c5fc7d",
      "path": "packages/next/src/compiled/client-only/package.json",
      "url": "https://api.github.com/repos/vercel/next.js/contents/packages/next/src/compiled/client-only/package.json?ref=299180d3315c7ebd7b199d2b1a265b5986c5fc7d",
      "status": 200,
      "response_path": "fixed-source/98ba32f000673c1ae5b2.response.json",
      "response_sha256": "25eab18db37b2e9dfaee9f6fb1d55faad9ad6fe0a7615a4550cd3d9bd75a564b",
      "local_path": "fixed-source/98ba32f000673c1ae5b2.txt",
      "sha256": "4d6342705767832f299b9a59c28e4275bcf02db19472732f93f67d979441df8f",
      "bytes": 467,
      "git_blob_sha1": "9af6d8dfe9949064ba23a8f0a46134e14bb616d1"
    },
    "licenses/reference/babel-packages-inline-notices.txt": {
      "basis": "Exact copyright/license comments extracted in source order from unchanged official Next bundled code, retaining each original fragment. Standard MIT terms and Babel copyright also supplied.",
      "bundle_sha256": "1a20ee41e0ac0a702c4fa96a1b85530a14f8bc35a58c38df5cd0e3e738cc3c8e",
      "comments": 3,
      "derivation": "Extract matching complete block comments, concatenate with two LF; no replacement of original copyright/terms."
    }
  },
  "modes": {
    "licenses/reference/next-MIT.txt": "VERBATIM",
    "licenses/reference/constants-browserify-README.md.txt": "VERBATIM",
    "licenses/reference/data-uri-to-buffer-README.md.txt": "VERBATIM",
    "licenses/reference/http-proxy-agent-README.md.txt": "VERBATIM",
    "licenses/reference/https-proxy-agent-README.md.txt": "VERBATIM",
    "licenses/reference/punycode-LICENSE-MIT.txt.txt": "VERBATIM",
    "licenses/reference/setimmediate-LICENSE.txt.txt": "VERBATIM",
    "licenses/reference/string-hash-README.md.txt": "VERBATIM",
    "licenses/reference/unistore-README.md.txt": "VERBATIM",
    "licenses/reference/edge-runtime-cookies-LICENSE.md.txt": "VERBATIM",
    "licenses/reference/edge-runtime-ponyfill-LICENSE.md.txt": "VERBATIM",
    "licenses/reference/edge-runtime-primitives-LICENSE.md.txt": "VERBATIM",
    "licenses/reference/babel-core-LICENSE.txt": "VERBATIM",
    "licenses/reference/dotenv-LICENSE.txt": "VERBATIM",
    "licenses/reference/dotenv-expand-LICENSE.txt": "VERBATIM",
    "licenses/reference/unistore-MIT-terms.txt": "ADAPTED",
    "licenses/reference/client-only-MIT-terms.txt": "ADAPTED",
    "licenses/reference/client-only-package.json.txt": "VERBATIM",
    "licenses/reference/babel-packages-inline-notices.txt": "ADAPTED",
    "licenses/reference/runtime-notice-lock.json": "AUTHORED"
  },
  "upstream_build_metadata_sha256": "d764d69e2940b30afd879baf0df67bd480d6df78ce55f8a7ad284398eddcafe1",
  "no_original_copyright_invented": true
}
````

### FILE: `ci/collect_runtime_notices.py`

```yaml
block_id: "REFERENCE-NOTICE-BLOCKS:file22:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "535230ec76013c14fb7329c7d14e2b19dafcaaa4746819ac40a1e6a101d9f1a2"
variables: []
secrets_allowed: false
```

````python
"""AUTHORED packaging glue: bind actual runtime manifests to retained notices."""
from pathlib import Path
import hashlib,json,re

HEX=re.compile(r'[0-9a-f]{64}')
LEGAL=re.compile(r'(?:license|licence|notice|copying)(?:[-.].*)?|.*\.(?:LEGAL|LICENSE)\.txt',re.I)

def require(condition,message):
    if not condition:raise ValueError(message)

def sha(path):return hashlib.sha256(path.read_bytes()).hexdigest()

def no_duplicates(pairs):
    result={}
    for key,value in pairs:
        require(key not in result,'duplicate JSON field');result[key]=value
    return result

def read(path):
    require(path.is_file()and path.stat().st_size<=2000000,'bounded JSON required')
    return json.loads(path.read_bytes(),object_pairs_hook=no_duplicates)

def relative(root,value):
    require(isinstance(value,str)and value and '\\'not in value and ':'not in value,'invalid notice path')
    parts=value.split('/');require(all(x and x not in('.','..')for x in parts),'invalid notice path')
    p=root.joinpath(*parts)
    for item in [p,*p.parents]:
        if item.exists():require(not item.is_symlink()and not item.is_junction(),'reparse notice path')
        if item==root:break
    return p

def collect(source,runtime,output,lock_path):
    source=Path(source);runtime=Path(runtime);output=Path(output)
    for root in [source,runtime,output]:
        require(root.is_dir(),'existing notice root required')
        for p in [root,*root.parents]:require(not p.is_symlink()and not p.is_junction(),'reparse notice root')
    lock=read(Path(lock_path));require(set(lock)=={'schema','source_scope','manifests','texts','build_only_exclusions'},'notice lock fields')
    require(lock['schema']=='elite-reference-runtime-notices/v1','notice schema')
    require(isinstance(lock['manifests'],dict)and 0<len(lock['manifests'])<=1000,'manifest lock budget')
    staged={};total=0
    for rel,expected in lock['texts'].items():
        p=relative(source,rel);require(HEX.fullmatch(expected)and p.is_file()and sha(p)==expected,'notice text changed: '+rel)
        total+=p.stat().st_size;require(total<=16777216,'notice text budget')
        staged[rel]=p
    require(staged,'empty notice texts')
    actual={};seen=set()
    for p in (runtime/'node_modules').rglob('*'):
        require(not p.is_symlink()and not p.is_junction(),'runtime reparse entry')
        if not p.is_file():continue
        require(p.suffix.lower()not in('.node','.dll','.wasm'),'unadmitted native runtime payload')
        if p.name!='package.json':continue
        rel=p.relative_to(runtime).as_posix();require(rel.casefold()not in seen,'runtime case alias');seen.add(rel.casefold());actual[rel]=p
        require(len(actual)<=1000,'runtime manifest budget')
    require(set(actual)==set(lock['manifests']),'unknown or missing runtime manifest')
    rows=[]
    for rel,p in sorted(actual.items()):
        record=lock['manifests'][rel];require(set(record)=={'sha256','name','version','license','notice_refs','source_basis'},'manifest entry fields')
        require(sha(p)==record['sha256'],'runtime manifest changed: '+rel)
        original=relative(source,rel);require(original.is_file()and sha(original)==record['sha256'],'installed manifest correspondence')
        v=read(p);require((v.get('name'),v.get('version'),v.get('license'))==(record['name'],record['version'],record['license']),'runtime metadata differs')
        refs=record['notice_refs'];require(isinstance(refs,list)and refs and len(refs)==len(set(refs))and all(x in staged for x in refs),'runtime notice missing')
        require(any(staged[x].stat().st_size>100 for x in refs),'runtime notice empty')
        rows.append({'manifest':rel,**record,'texts':{x:lock['texts'][x]for x in refs}})
    identities={str(x['name'])+'@'+str(x['version'])for x in rows if x['name']and x['version']}
    require(not identities.intersection(lock['build_only_exclusions']),'excluded build tool was redistributed')
    installed=[]
    for scope in sorted((source/'node_modules/.pnpm').glob('*/node_modules/*')):
        roots=list(scope.iterdir())if scope.name.startswith('@')and scope.is_dir()else[scope]
        for root in roots:
            if root.is_symlink()or root.is_junction()or not(root/'package.json').is_file():continue
            v=read(root/'package.json');identity=str(v.get('name'))+'@'+str(v.get('version'));texts={}
            for p in sorted(root.iterdir()):
                if p.is_file()and LEGAL.fullmatch(p.name):
                    rel=p.relative_to(source).as_posix();require(rel in staged,'installed legal text not locked: '+rel);texts[p.name]=lock['texts'][rel]
            in_runtime=identity in identities
            if in_runtime and not texts:
                refs={r for row in rows if row['name']==v.get('name')and row['version']==v.get('version')for r in row['notice_refs']}
                texts={r:lock['texts'][r]for r in sorted(refs)}
            installed.append({'package':identity,'package_json_sha256':sha(root/'package.json'),'declared_license':v.get('license'),'texts':texts,'included_in_runtime':in_runtime,'scope':'REDISTRIBUTED_RUNTIME'if in_runtime else'BUILD_INPUT_NOT_REDISTRIBUTED'})
    # Validate the full graph before any output. Publication remains create-only.
    names=['dependency-notices.json','runtime-dependency-notices.json']
    require(all(not(output/n).exists()for n in names),'notice output already exists')
    dest=output/'dependency-notices';require(not dest.exists()or dest.is_dir()and not dest.is_symlink()and not dest.is_junction(),'invalid legal destination');dest.mkdir(exist_ok=True)
    for rel,p in sorted(staged.items()):
        target=dest/lock['texts'][rel]
        if target.exists():require(target.is_file()and not target.is_symlink()and sha(target)==lock['texts'][rel],'notice collision')
        else:
            with target.open('xb')as f:f.write(p.read_bytes())
    result={'schema':'elite-runtime-legal-catalogue/v1','state':'PASS_EXACT_NOTICE_CORRESPONDENCE','scope':lock['source_scope'],'runtime_manifests':rows,'notice_texts':lock['texts'],'build_only_exclusions':lock['build_only_exclusions'],'native_payloads':0}
    for name,value in zip(names,[installed,result]):
        with(output/name).open('x',encoding='utf-8',newline='\n')as f:f.write(json.dumps(value,sort_keys=True,indent=2,ensure_ascii=False)+'\n')
    return {'runtime_manifests':len(rows),'installed_identities':len(installed),'unique_notice_texts':len(set(lock['texts'].values())),'excluded_build_tools':len(lock['build_only_exclusions'])}
````

### FILE: `ci/test_runtime_notices.py`

```yaml
block_id: "REFERENCE-NOTICE-BLOCKS:file23:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "a37f1c31eff464669dcc14f8c0a88c3a3fcddc2a3da1f1e4ae69e58d11b59a1f"
variables: []
secrets_allowed: false
```

````python
"""Local packaging regressions; all fixture data is synthetic."""
from pathlib import Path
import tempfile,json,hashlib,unittest
import collect_runtime_notices as n

class NoticeBoundaryTests(unittest.TestCase):
    def setUp(self):
        self.temp=tempfile.TemporaryDirectory();self.addCleanup(self.temp.cleanup);self.root=Path(self.temp.name)
        self.source=self.root/'source';self.runtime=self.root/'runtime';self.out=self.root/'out'
        for p in [self.source,self.runtime,self.out]:p.mkdir()
        self.rel='node_modules/.pnpm/example@1.0.0/node_modules/example/package.json'
        data=b'{"name":"example","version":"1.0.0","license":"MIT"}\n'
        for root in [self.source,self.runtime]:
            p=root/self.rel;p.parent.mkdir(parents=True);p.write_bytes(data)
        self.notice='licenses/example.txt';p=self.source/self.notice;p.parent.mkdir();p.write_text('Synthetic license fixture only. '*8)
        self.lock={'schema':'elite-reference-runtime-notices/v1','source_scope':'SYNTHETIC_FIXTURE',
          'manifests':{self.rel:{'sha256':hashlib.sha256(data).hexdigest(),'name':'example','version':'1.0.0','license':'MIT','notice_refs':[self.notice],'source_basis':'synthetic'}},
          'texts':{self.notice:n.sha(p)},'build_only_exclusions':[]}
        self.path=self.source/'lock.json'

    def run_collect(self):
        self.path.write_text(json.dumps(self.lock));return n.collect(self.source,self.runtime,self.out,self.path)

    def rejected(self):
        with self.assertRaises(ValueError):self.run_collect()
        self.assertFalse((self.out/'runtime-dependency-notices.json').exists())

    def test_unknown_runtime_manifest_blocks_release(self):
        p=self.runtime/'node_modules/injected/package.json';p.parent.mkdir();p.write_text('{}');self.rejected()

    def test_missing_runtime_dependency_blocks_release(self):
        (self.runtime/self.rel).unlink();self.rejected()

    def test_modified_runtime_manifest_blocks_release(self):
        (self.runtime/self.rel).write_text('{"name":"different"}');self.rejected()

    def test_changed_installed_source_blocks_release(self):
        (self.source/self.rel).write_text('{}');self.rejected()

    def test_changed_or_missing_notice_blocks_release(self):
        (self.source/self.notice).write_text('changed');self.rejected()
        (self.source/self.notice).unlink();self.rejected()

    def test_empty_notice_mapping_blocks_release(self):
        self.lock['manifests'][self.rel]['notice_refs']=[];self.rejected()

    def test_empty_notice_file_is_not_legal_evidence(self):
        p=self.source/self.notice;p.write_bytes(b'');self.lock['texts'][self.notice]=n.sha(p);self.rejected()

    def test_path_escape_blocks_before_publication(self):
        self.lock['texts']['../outside.txt']='0'*64;self.rejected()

    def test_unadmitted_native_payload_blocks_release(self):
        p=self.runtime/'node_modules/payload.node';p.write_bytes(b'synthetic-not-a-binary');self.rejected()

    def test_excluded_build_tool_cannot_be_shipped(self):
        self.lock['build_only_exclusions']=['example@1.0.0'];self.rejected()

    def test_existing_catalogue_is_never_overwritten(self):
        self.run_collect();before=(self.out/'runtime-dependency-notices.json').read_bytes()
        with self.assertRaises(ValueError):self.run_collect()
        self.assertEqual(before,(self.out/'runtime-dependency-notices.json').read_bytes())

    def test_duplicate_json_fields_rejected(self):
        self.path.write_text('{"schema":1,"schema":2}')
        with self.assertRaisesRegex(ValueError,'duplicate'):n.collect(self.source,self.runtime,self.out,self.path)

if __name__=='__main__':unittest.main()
````

### FILE: `docs/REFERENCE_RUNTIME_NOTICES.md`

```yaml
block_id: "REFERENCE-NOTICE-BLOCKS:file24:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "909299def7e6623781db3226cf9cd594055d7f82500abd61e8e60f0f60d83ac1"
variables: []
secrets_allowed: false
```

````markdown
# Exact local reference runtime notices

The portable builder checks the actual traced runtime against
`licenses/reference/runtime-notice-lock.json` before producing its catalogue.
All129observed package manifests are bound by SHA256, including framework bundles
without published version fields and unnamed export metadata. Versions are not
invented.71installed package identities are listed separately; an installed build
input is not automatically redistributed. The exact reference contains no native
addon, DLL or WASM payload under its web runtime.

The collector retains direct and nested original license/notice/LEGAL files and
the fixed source supplements listed in SOURCE_NOTICES.json. It validates all
inputs before writing the catalogue. Missing or changed runtime manifests,
unknown dependencies, missing/changed/empty notices, duplicate JSON fields,
unsafe paths, reparses and an excluded build tool in the runtime reject. Existing
catalogues cannot be overwritten. Notice collection neither installs dependencies
nor executes their scripts, and does not replace the separate SCA/source gates.

Next16.3.4 official source commit299180d3315c7ebd7b199d2b1a265b5986c5fc7d binds
the package manifests and bundling recipe. Full original README files are retained
where the publisher places its license there. The embedded dotenv16.3.1 and
dotenv-expand10.0.0 carry their original BSD-2-Clause texts. The three edge-runtime
packages carry their original npm LICENSE.md files alongside embedded notices.
String-hash retains its publisher's original CC0 waiver. Local glue is AUTHORED;
original notice files are VERBATIM, not an attribution of library code.

Three notice files are explicitly ADAPTED, not claimed as original upstream files:

- Unistore3.4.1 offers MIT with copyright Jason Miller in its original README.
  Its separate notice reproduces that stated owner, without inventing a year,
  and the standard MIT permission paragraphs. The README grant accompanies it.
- Client-only0.0.1 declares MIT in its original package.json without a copyright
  notice. The notice reproduces standard MIT terms without assigning an owner.
  Its three original files match the commit-fixed Next distributed source.
- The Babel bundle notice contains complete original copyright/license comments
  extracted from the exact bundle in their original order. The original Babel
  MIT notice and fixed Next build/import metadata accompany this correspondence.

The standard MIT permission paragraphs above match the fixed original Next MIT
text; its Vercel copyright heading is not assigned to Unistore or client-only.
These are declared notice renderings, not fabricated author/source receipts.

The runtime excludes @next/swc-win32-x64-msvc16.3.4,
@rolldown/binding-win32-x64-msvc1.2.5, server-only0.0.1 and stackback0.0.2.
Their build-input metadata remains visible. This exclusion does not certify a
native binary or authorize redistribution of the restricted package-manager
bundle. The lock must be regenerated and reviewed with any package/trace change.

The tests cover actual historical reference runtime files read-only and12adverse
publication cases. The final post-ARCA artifact must run this same collector as
part of both independent builds. No production, provider or jurisdiction claim
is derived from this local packaging proof.
````


V402322: original license/README sources and notice derivations are declared in SOURCE_NOTICES.json; root pack ownership does not override those per-file licenses. Runtime collection is connected to ci/build_local_reference.py and rejects release on drift.

V402 composed delta: V402324 signed-build driver reuses actual local builder twice at one stable workspace with independent installs/caches, retained evidence and exact artifact/source correspondence; no provider or corporate identity claim.

### FILE: `ci/build_signed_reference.ps1`

```yaml
block_id: "LOCAL-SIGNED-BUILD:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "69cbb6cf001b1b52ade5b8adeb43a06d615fb7b097c70a2f332d48f93e4e6b56"
variables: []
secrets_allowed: false
```

````text
#requires -Version 7.0
# AUTHORED orchestration for the admitted local builder; no provider credentials.
param(
 [Parameter(Mandatory=$true)][string]$ProjectRoot,
 [Parameter(Mandatory=$true)][string]$OutputRoot,
 [Parameter(Mandatory=$true)][long]$SourceDateEpoch
)
Set-StrictMode -Version Latest
$ErrorActionPreference='Stop'
$utf8=[Text.UTF8Encoding]::new($false)
function Fail([string]$message){throw ('LOCAL_SIGNED_BUILD: '+$message)}
function Plain([string]$path){
 $full=[IO.Path]::GetFullPath($path).TrimEnd([IO.Path]::DirectorySeparatorChar)
 $cursor=$full
 while($cursor){
  $item=Get-Item -LiteralPath $cursor -Force -ErrorAction SilentlyContinue
  if($null-ne$item-and($item.Attributes-band[IO.FileAttributes]::ReparsePoint)-ne0){Fail 'Reparse path rejected'}
  $cursor=Split-Path -Parent $cursor
 }
 $full
}
function Sha([string]$path){(Get-FileHash -LiteralPath $path -Algorithm SHA256).Hash.ToLowerInvariant()}
function WriteJson([string]$path,$value){[IO.File]::WriteAllText($path,($value|ConvertTo-Json -Depth 10)+"`n",$utf8)}
function Disjoint([string]$a,[string]$b){
 if($a-eq$b-or$a.StartsWith($b+[IO.Path]::DirectorySeparatorChar,[StringComparison]::OrdinalIgnoreCase)-or$b.StartsWith($a+[IO.Path]::DirectorySeparatorChar,[StringComparison]::OrdinalIgnoreCase)){Fail 'Overlapping build paths'}
}
$root=Plain $ProjectRoot;$output=Plain $OutputRoot
Disjoint $root $output
if(-not(Test-Path -LiteralPath $root -PathType Container)-or-not(Test-Path -LiteralPath $output -PathType Container)-or@(Get-ChildItem -LiteralPath $output -Force).Count-ne0){Fail 'Existing source and empty output directories required'}
$configPath=Join-Path $root '.elite-local-build-inputs.json'
[void](Plain $configPath)
$configSha=Sha $configPath
$cfg=Get-Content -LiteralPath $configPath -Raw|ConvertFrom-Json -AsHashtable
$expected=@('schema_version','scope','python','source_inventory_ref','source_inventory_sha256','go','node','go_mod_cache','install_inputs','install_inputs_sha256','arca_inputs','arca_inputs_sha256','run_root')
if(@(Compare-Object @($cfg.Keys|Sort-Object) @($expected|Sort-Object)).Count-ne0-or$cfg.schema_version-ne1-or$cfg.scope-cne'LOCAL_FIXTURES'){Fail 'Exact local input schema required'}
if($SourceDateEpoch-ne315532800){Fail 'Reference ZIP timestamp must be 1980-01-01 UTC'}
if(@(Compare-Object @($cfg.python.Keys|Sort-Object) @('path','sha256')).Count-ne0){Fail 'Exact Python pin required'}
$python=Plain ([string]$cfg.python.path)
if($cfg.python.sha256-notmatch'^[0-9a-f]{64}$'-or(Sha $python)-cne$cfg.python.sha256){Fail 'Python pin changed'}
$inventoryRef=[string]$cfg.source_inventory_ref
if([IO.Path]::IsPathRooted($inventoryRef)-or$inventoryRef-match'[:\\]' -or @($inventoryRef.Split('/')|Where-Object{$_-in@('','..','.')}).Count){Fail 'Safe relative inventory required'}
$inventory=Plain (Join-Path $root $inventoryRef)
if($cfg.source_inventory_sha256-notmatch'^[0-9a-f]{64}$'-or(Sha $inventory)-cne$cfg.source_inventory_sha256){Fail 'Source inventory pin changed'}
$run=Plain ([string]$cfg.run_root)
Disjoint $run $root;Disjoint $run $output
if(-not(Test-Path -LiteralPath (Split-Path -Parent $run) -PathType Container)){Fail 'Run parent absent'}
$markerPath=Join-Path $run 'session.json'
if(-not(Test-Path -LiteralPath $run)){
 [void](New-Item -ItemType Directory -Path $run)
 WriteJson $markerPath @{schema_version=1;configuration_sha256=$configSha;source_inventory_sha256=$cfg.source_inventory_sha256;completed_builds=0}
}
[void](Plain $markerPath);[void](Plain (Join-Path $run 'session.lock'))
$lock=[IO.File]::Open((Join-Path $run 'session.lock'),[IO.FileMode]::OpenOrCreate,[IO.FileAccess]::ReadWrite,[IO.FileShare]::None)
try{
 $marker=Get-Content -LiteralPath $markerPath -Raw|ConvertFrom-Json -AsHashtable
 if($marker.schema_version-ne1-or$marker.configuration_sha256-cne$configSha-or$marker.source_inventory_sha256-cne$cfg.source_inventory_sha256-or$marker.completed_builds-notin@(0,1)){Fail 'Run belongs to another configuration or already completed both builds'}
 $number=[int]$marker.completed_builds+1
 $workspace=Plain (Join-Path $run 'workspace')
 $archive=Plain (Join-Path $run ('workspace-'+$number))
 $build=Plain (Join-Path $run ('build-'+$number))
 foreach($path in @($workspace,$archive,$build)){if(Test-Path -LiteralPath $path){Fail 'Previous or partial build retained; use a new run_root after diagnosing it'}}
 $arguments=@('-X','utf8','-B',(Join-Path $root 'ci/build_local_reference.py'),'--source',$root,'--target',$build,'--workspace',$workspace,'--inventory',$inventory,'--inventory-sha256',[string]$cfg.source_inventory_sha256,'--go',[string]$cfg.go,'--node',[string]$cfg.node,'--go-mod-cache',[string]$cfg.go_mod_cache,'--install-inputs',[string]$cfg.install_inputs,'--install-inputs-sha256',[string]$cfg.install_inputs_sha256,'--arca-inputs',[string]$cfg.arca_inputs,'--arca-inputs-sha256',[string]$cfg.arca_inputs_sha256)
 & $python @arguments
 if($LASTEXITCODE-ne0){Fail 'Actual independent build failed; all working evidence retained'}
 if((Sha $configPath)-cne$configSha-or(Sha $inventory)-cne$cfg.source_inventory_sha256){Fail 'Build inputs changed'}
 & $python -X utf8 -B (Join-Path $root 'ci/publish_reference_build.py') --build $build --output $output --inventory $inventory --inventory-sha256 ([string]$cfg.source_inventory_sha256)
 if($LASTEXITCODE-ne0){Fail 'Artifact/receipt correspondence failed'}
 # Both names were resolved and constrained to the exact owned run directory.
 if((Split-Path -Parent (Plain $workspace))-cne$run-or(Split-Path -Parent (Plain $archive))-cne$run-or(Test-Path -LiteralPath $archive)){Fail 'Unsafe workspace archive target'}
 Move-Item -LiteralPath $workspace -Destination $archive
 $marker.completed_builds=$number
 WriteJson $markerPath $marker
 Write-Output ('LOCAL_SIGNED_BUILD_PASS independent_build='+$number+' evidence='+$build)
}finally{$lock.Dispose()}
````

### FILE: `ci/publish_reference_build.py`

```yaml
block_id: "LOCAL-SIGNED-BUILD:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "dc9a7fdfc6ebe609eff6b8a78deefa7ddce7b123491797923a91ddb07a3e2758"
variables: []
secrets_allowed: false
```

````python
"""AUTHORED final projection; only a fresh PASS receipt authorizes artifact copy."""
import argparse,hashlib,json
from pathlib import Path
from build_local_reference import plain,read_json,raw_json,tree,copy_tree,require

def publish(build,output,inventory,inventory_sha256):
    build=plain(build);output=plain(output)
    require(build!=output and build not in output.parents and output not in build.parents,'overlapping publication')
    require(output.is_dir() and not list(output.iterdir()),'empty publication directory required')
    source=read_json(inventory,inventory_sha256)
    source_id=hashlib.sha256(raw_json(source)).hexdigest()
    receipt=json.loads((build/'build-result.json').read_text(encoding='utf-8'))
    require(receipt['schema']=='elite-local-reference-build/v1' and receipt['scope']=='LOCAL_FIXTURES' and receipt['state']=='PASS','actual successful local build required')
    require(receipt['source_sha256']==source_id,'receipt source mismatch')
    actual=tree(build/'artifact')
    require(actual and actual==receipt['files'] and hashlib.sha256(raw_json(actual)).hexdigest()==receipt['artifact_sha256'],'artifact changed after build')
    require(all(actual.get('source/'+rel)==digest for rel,digest in source.items()),'corresponding source mismatch')
    require(json.loads((build/'artifact-inventory.json').read_text(encoding='utf-8'))==actual,'build inventory mismatch')
    copy_tree(build/'artifact',output)
    require(tree(output)==actual,'published artifact mismatch')
    print('REFERENCE_ARTIFACT_PROJECTED '+receipt['artifact_sha256'])

if __name__=='__main__':
    p=argparse.ArgumentParser(description=__doc__)
    for name in ['build','output','inventory','inventory-sha256']:p.add_argument('--'+name,required=True)
    args=p.parse_args()
    publish(args.build,args.output,args.inventory,args.inventory_sha256)
````

### FILE: `ci/test_signed_reference.py`

```yaml
block_id: "LOCAL-SIGNED-BUILD:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "78e34aa7692ef51820866e265e530483473e580d31136c5ad3a04d008b972749"
variables: []
secrets_allowed: false
```

````python
"""AUTHORED adverse publication/pre-build contracts, not full build evidence."""
import argparse,hashlib,json,subprocess,sys
from pathlib import Path
from publish_reference_build import publish
from build_local_reference import tree,raw_json,corresponding_sources

def test_runtime_junctions(node,target):
    """Actual Node loader regression: ZIP projection contains no empty scope dir."""
    target=Path(target).resolve();target.mkdir()
    package=target/'node_modules/.store/fixture';package.mkdir(parents=True)
    (package/'index.js').write_text('module.exports = 42;\n',encoding='utf-8')
    (target/'elite-runtime-links.json').write_bytes(raw_json({'node_modules/@fixture/module':'node_modules/.store/fixture'}))
    server=target/'server.js';server.write_text("if(require('@fixture/module')!==42)throw new Error('wrong module');process.on('SIGTERM',()=>process.exit(143));console.log('ZIP_JUNCTION_PASS');\n",encoding='utf-8')
    assert not(target/'node_modules/@fixture').exists()
    q=subprocess.run([node,str(Path(__file__).with_name('web_local_stop.cjs')),str(server)],input=b'STOP\n',stdout=subprocess.PIPE,stderr=subprocess.STDOUT,timeout=15)
    (target/'node.log').write_bytes(q.stdout)
    assert q.returncode==143 and b'ZIP_JUNCTION_PASS'in q.stdout,q.stdout.decode('utf-8','replace')
    assert(target/'node_modules/@fixture/module/index.js').read_bytes()==(package/'index.js').read_bytes()
    return {'state':'PASS','empty_parent_created':True,'real_node_module_resolution':True,'graceful_stop':True}

def run(target,pwsh):
    target=Path(target).resolve();target.mkdir()
    results=[];fixture_source={'a.txt':hashlib.sha256(b'original').hexdigest()}
    inventory=target/'inventory.json';inventory.write_bytes(raw_json(fixture_source))
    inventory_hash=hashlib.sha256(inventory.read_bytes()).hexdigest()
    def fixture(name):
        base=target/name;base.mkdir();build=base/'build';build.mkdir();artifact=build/'artifact';artifact.mkdir()
        (artifact/'source').mkdir();(artifact/'source/a.txt').write_bytes(b'original')
        output=base/'output';output.mkdir()
        files=tree(artifact);receipt={'schema':'elite-local-reference-build/v1','scope':'LOCAL_FIXTURES','state':'PASS','source_sha256':inventory_hash,'files':files,'artifact_sha256':hashlib.sha256(raw_json(files)).hexdigest()}
        (build/'build-result.json').write_bytes(raw_json(receipt));(build/'artifact-inventory.json').write_bytes(raw_json(files))
        return build,output,receipt
    for name in ['valid','changed-artifact','wrong-source','failed-receipt','occupied-output','changed-inventory','coherent-wrong-source']:
        build,output,receipt=fixture(name)
        if name=='changed-artifact':(build/'artifact/source/a.txt').write_bytes(b'modified')
        if name=='wrong-source':receipt['source_sha256']='0'*64
        if name=='failed-receipt':receipt['state']='FAIL'
        if name=='occupied-output':(output/'existing').write_text('keep')
        if name=='changed-inventory':(build/'artifact-inventory.json').write_text('{}')
        if name=='coherent-wrong-source':
            (build/'artifact/source/a.txt').write_bytes(b'modified');receipt['files']=tree(build/'artifact');receipt['artifact_sha256']=hashlib.sha256(raw_json(receipt['files'])).hexdigest();(build/'artifact-inventory.json').write_bytes(raw_json(receipt['files']))
        (build/'build-result.json').write_bytes(raw_json(receipt))
        try:publish(build,output,inventory,inventory_hash);accepted=True
        except ValueError:accepted=False
        assert accepted==(name=='valid'),name
        results.append({'case':name,'state':'PASS','accepted':accepted})
    for name in ['exact-source','generated-next-declaration','changed-domain-source','unknown-next-declaration']:
        base=target/name;base.mkdir();source=base/'input';built=base/'built';out=base/'artifact'
        for path in [source,built,out]:path.mkdir()
        declaration=b'/// <reference types="next/image-types/global" />\n'
        original={'next-env.d.ts':declaration,'domain.txt':b'original-domain'}
        for rel,data in original.items():(source/rel).write_bytes(data);(built/rel).write_bytes(data)
        files={rel:hashlib.sha256(data).hexdigest()for rel,data in original.items()}
        if name=='generated-next-declaration':(built/'next-env.d.ts').write_bytes(declaration+b'import "./.next/types/routes.d.ts";\nimport "./.next/types/root-params.d.ts";\n')
        if name=='changed-domain-source':(built/'domain.txt').write_bytes(b'changed')
        if name=='unknown-next-declaration':(built/'next-env.d.ts').write_bytes(b'unknown-generation')
        try:derived=corresponding_sources(source,built,out,files);accepted=True
        except ValueError:accepted=False
        assert accepted==(name in ['exact-source','generated-next-declaration']),name
        if accepted:
            assert tree(out/'source')==files
            assert bool(derived)==(name=='generated-next-declaration')
        results.append({'case':name,'state':'PASS','accepted':accepted})
    script=Path(__file__).with_name('build_signed_reference.ps1')
    for name in ['python-pin','source-pin','overlap','partial-build','already-complete']:
        base=target/name;base.mkdir();root=base/'project';root.mkdir();output=base/'output';output.mkdir();runroot=base/'session'
        inv=root/'inventory.json';inv.write_bytes(inventory.read_bytes())
        cfg={'schema_version':1,'scope':'LOCAL_FIXTURES','python':{'path':sys.executable,'sha256':hashlib.sha256(Path(sys.executable).read_bytes()).hexdigest()},'source_inventory_ref':'inventory.json','source_inventory_sha256':inventory_hash,'go':'unused-before-validation','node':'unused-before-validation','go_mod_cache':'unused-before-validation','install_inputs':'unused-before-validation','install_inputs_sha256':'0'*64,'arca_inputs':'unused-before-validation','arca_inputs_sha256':'0'*64,'run_root':str(runroot)}
        if name=='python-pin':cfg['python']['sha256']='0'*64
        if name=='source-pin':cfg['source_inventory_sha256']='0'*64
        if name=='overlap':cfg['run_root']=str(root/'overlap')
        config=root/'.elite-local-build-inputs.json';config.write_bytes(raw_json(cfg))
        if name in ['partial-build','already-complete']:
            runroot.mkdir();(runroot/'session.json').write_bytes(raw_json({'schema_version':1,'configuration_sha256':hashlib.sha256(config.read_bytes()).hexdigest(),'source_inventory_sha256':inventory_hash,'completed_builds':2 if name=='already-complete' else 0}))
            if name=='partial-build':(runroot/'build-1').mkdir()
        q=subprocess.run([pwsh,'-NoProfile','-File',str(script),'-ProjectRoot',str(root),'-OutputRoot',str(output),'-SourceDateEpoch','315532800'],stdout=subprocess.PIPE,stderr=subprocess.STDOUT,timeout=30)
        (base/'rejection.log').write_bytes(q.stdout)
        expected={'python-pin':'Python pin changed','source-pin':'Source inventory pin changed','overlap':'Overlapping build paths','partial-build':'Previous or partial build retained','already-complete':'already completed both builds'}[name]
        assert q.returncode!=0 and expected in q.stdout.decode('utf-8','replace'),(name,q.stdout.decode('utf-8','replace'))
        assert not list(output.iterdir()),name
        results.append({'case':name,'state':'PASS','accepted':False})
    result={'state':'PASS_CONTRACT_FIXTURES_ONLY','cases':results,'actual_product_build_performed':False}
    (target/'result.json').write_bytes(raw_json(result));print(json.dumps(result))

if __name__=='__main__':
    p=argparse.ArgumentParser(description=__doc__);p.add_argument('--target',required=True);p.add_argument('--pwsh',required=True);p.add_argument('--node');args=p.parse_args();run(args.target,args.pwsh)
    if args.node:print(json.dumps(test_runtime_junctions(args.node,Path(args.target)/'zip-runtime')))
````

### FILE: `ci/local-build-inputs.template.json`

```yaml
block_id: "LOCAL-SIGNED-BUILD:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "2049f60f1db59918e98b0992b15013409caed555d1c13bce7e36f5ed8cbff4fd"
variables: []
secrets_allowed: false
```

````json
{
  "schema_version": 1,
  "scope": "LOCAL_FIXTURES",
  "python": {"path": "", "sha256": ""},
  "source_inventory_ref": ".elite-source-inventory.json",
  "source_inventory_sha256": "",
  "go": "",
  "node": "",
  "go_mod_cache": "",
  "install_inputs": "",
  "install_inputs_sha256": "",
  "arca_inputs": "",
  "arca_inputs_sha256": "",
  "run_root": ""
}
````

### FILE: `docs/SIGNED_LOCAL_REFERENCE.md`

```yaml
block_id: "LOCAL-SIGNED-BUILD:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "29e0010b1386118c1ccaaf60e15d36f69602dd89dfa9a4ab8d2cca72d1d06a84"
variables: []
secrets_allowed: false
```

````markdown
# Firma de la referencia local

AUTHORED glue sobre los owners admitidos Go, Next, pnpm restringido y OpenSSH.
No representa identidad corporativa, credenciales productivas ni un build service
aislado. La clave privada y el almacén de herramientas quedan fuera del producto.

El caller genera `.elite-local-build-inputs.json` desde el template, fija el
ejecutable Python y el inventario exacto materializado, y declara un `run_root`
externo ausente. Incluye configuración, inventario, driver y todos los manifests
de dependencias elegidos en `source_inputs_ref` del gate firmado. No contiene
secretos. El perfil del gate fija el SHA del driver `ci/build_signed_reference.ps1`.

Cada invocación ejecuta realmente `ci/build_local_reference.py`: copia limpia,
instalación offline restringida y GOCACHE independiente. Usa la misma ruta explícita
de workspace por el contrato de IDs RSC de Next16.3.4; archiva cada workspace terminado
con sus logs e inventario y nunca lo reutiliza como entrada compilada. Un lock
exclusivo permite exactamente dos builds por configuración. Fallos parciales se
conservan; requieren diagnóstico y un nuevo run_root, nunca borrado automático.

La proyección compara el árbol actual contra el recibo PASS y su inventario antes
y después de copiarlo al output del gate. El gate compara ambos ZIP deterministas,
escanea los manifests explícitos con OSV fijado y firma la procedencia. Su verificador
independiente exige `-TrustedAllowedSigners` desde una política confiada externa al
release recibido. El archivo de firmantes incluido en el ZIP no se autoriza solo.

El timestamp de ZIP es 1980-01-01 UTC. El claim de reproducibilidad se limita al
host/toolchain/arquitectura/workspace admitidos; no implica builds idénticos en
rutas arbitrarias o plataformas diferentes. Las fixtures de auth/build son públicas
y sólo locales. No se habilita deploy live con ellas.
````


V402 composed delta: V402325 connected local ARCA fixture, fixed offline generated-tool adaptation, explicit legal/source pins and portable .NET/Go build. See ARCA_CONNECTED_INFRA_V402 evidence; only future homologation credentials conditioned, no productive fiscal claim.

### FILE: `ci/arca-build-inputs.template.json`

```yaml
block_id: "PORTABLE_CI_QUALITY_GATE_RUNNER-V402325:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "b8f767f232b956ba72ad588b455b7eff636bb163528fe5cea17956754b2768dc"
variables: []
secrets_allowed: false
```

````json
{
  "schema_version": 1,
  "scope": "LOCAL_FIXTURES",
  "dotnet_root": "",
  "pwsh": {"path": "", "sha256": ""},
  "svcutil_archive": "",
  "svcutil_packages": "",
  "wsdl_cache": "",
  "runtime_packages": ""
}
````

### FILE: `ci/build_arca_reference.py`

```yaml
block_id: "PORTABLE_CI_QUALITY_GATE_RUNNER-V402325:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "7f6c3339adf91d5abc5f4a347599295070634030a6162a9853f186f931df7349"
variables: []
secrets_allowed: false
```

````python
"""AUTHORED offline build glue for the existing ARCA Go/UDS/.NET owners."""
from __future__ import annotations

import hashlib
import json
import os
from pathlib import Path
import shutil
import sys
import zipfile

from build_local_reference import command, plain, raw_json, read_json, relative, require, sha


def build_arca(src: Path, dist: Path, work: Path, logs: Path, go: Path,
               environment: dict, input_path: str, input_sha256: str) -> dict:
    cfg = read_json(input_path, input_sha256)
    require(set(cfg) == {"schema_version", "scope", "dotnet_root", "pwsh", "svcutil_archive",
                        "svcutil_packages", "wsdl_cache", "runtime_packages"}, "Exact ARCA build inputs required")
    require(cfg["schema_version"] == 1 and cfg["scope"] == "LOCAL_FIXTURES", "ARCA scope mismatch")
    assets = src / "arca/fiscal"
    sdk = plain(cfg["dotnet_root"])
    sdk_lock = json.loads((assets / "dotnet-sdk.lock.json").read_text())
    require(sdk_lock["schema"] == "elite-arca-dotnet-sdk/v1", "Unknown .NET SDK inventory")
    for row in sdk_lock["files"]:
        path = relative(sdk, row["path"])
        require(path.stat().st_size == row["bytes"] and sha(path) == row["sha256"], "SDK input changed")
    runtime_lock = json.loads((assets / "dotnet-runtime.lock.json").read_text())
    require(runtime_lock["schema"] == "elite-arca-dotnet-runtime/v1", "Unknown .NET runtime projection")
    for row in runtime_lock["files"]:
        path = relative(sdk, row["path"])
        require(path.stat().st_size == row["bytes"] and sha(path) == row["sha256"], "Runtime pin changed")
    require(set(cfg["pwsh"]) == {"path", "sha256"}, "Exact PowerShell pin required")
    pwsh = plain(cfg["pwsh"]["path"])
    require(sha(pwsh) == cfg["pwsh"]["sha256"], "PowerShell pin changed")
    for name in ["svcutil_archive", "svcutil_packages", "wsdl_cache", "runtime_packages"]:
        cfg[name] = str(plain(cfg[name]))
    area = work / "arca"
    require(not area.exists(), "Fresh ARCA build area required")
    area.mkdir()
    generated = src / "arca/fiscal/generated"
    require(not generated.exists(), "Generated output must be absent")
    dotnet = sdk / "dotnet.exe"
    env = dict(environment)
    env.update(DOTNET_ROOT=str(sdk), DOTNET_ROOT_X64=str(sdk), DOTNET_CLI_TELEMETRY_OPTOUT="1",
               DOTNET_SVCUTIL_TELEMETRY_OPTOUT="1", DOTNET_NOLOGO="1", DOTNET_ROLL_FORWARD="Major",
               DOTNET_CLI_HOME=str(area / "home"), NUGET_PACKAGES=str(area / "packages"),
               APPDATA=str(area / "roaming"), LOCALAPPDATA=str(area / "local"),
               TEMP=str(area), TMP=str(area), PATH=str(sdk) + os.pathsep + env["PATH"])
    # NuGet's Windows configuration discovery requires the machine directory
    # locators even when all restore sources are explicitly local and cleared.
    for name in ["ProgramFiles", "ProgramFiles(x86)", "ProgramData", "SystemDrive"]:
        require(bool(os.environ.get(name)), "Windows machine directory missing: " + name)
        env[name] = os.environ[name]
    steps = []
    steps.append(command("arca-tool", [sys.executable, "-X", "utf8", "-B", src / "tools/assemble_arca_svcutil.py",
        "--source-archive", cfg["svcutil_archive"], "--package-directory", cfg["svcutil_packages"],
        "--output", area / "svcutil"], src, env, logs))
    steps.append(command("arca-generate", [pwsh, "-NoLogo", "-NoProfile", "-File", src / "tools/generate-arca-wsfe-client.ps1",
        "-DotnetExecutable", dotnet, "-CacheDirectory", cfg["wsdl_cache"], "-OutputDirectory", generated,
        "-SvcutilRuntimeDirectory", area / "svcutil", "-OfflinePackageDirectory", cfg["runtime_packages"], "-Offline"], src, env, logs))
    artifact = dist / "arca"
    require(not artifact.exists(), "ARCA artifact must be absent")
    artifact.mkdir()
    projects = {"worker": "worker/Elite.Arca.Wsfe.Worker/Elite.Arca.Wsfe.Worker.csproj",
                "fixture": "worker/fixtures/Elite.Arca.Wsfe.Fixture/Elite.Arca.Wsfe.Fixture.csproj"}
    for name, rel in projects.items():
        project = assets / rel
        steps.append(command("arca-" + name + "-restore", [dotnet, "restore", project, "--locked-mode", "--no-cache",
            "--warnaserror", "--configfile", generated / "NuGet.Config", "-p:NuGetAudit=false"], src, env, logs))
        steps.append(command("arca-" + name + "-publish", [dotnet, "publish", project, "--configuration", "Release",
            "--no-restore", "--no-self-contained", "--output", artifact / name, "--warnaserror", "-p:UseAppHost=false",
            "-p:Deterministic=true", "-p:ContinuousIntegrationBuild=true", "-p:DebugType=None",
            "-p:PathMap=" + str(src) + "=/_/source", "-p:NuGetAudit=false"], src, env, logs))
    (artifact / "go").mkdir()
    for name in ["arca-fiscal-worker", "arca-parameter-worker"]:
        steps.append(command(name, [go, "build", "-trimpath", "-buildvcs=false", "-ldflags=-s -w -buildid=",
            "-o", artifact / "go" / (name + ".exe"), "./cmd/" + name], src, env, logs))
    for row in runtime_lock["files"]:
        destination = artifact / "dotnet" / row["path"]
        destination.parent.mkdir(parents=True, exist_ok=True)
        shutil.copyfile(relative(sdk, row["path"]), destination)
        require(sha(destination) == row["sha256"], "Runtime copy drift")
    steps.append(command("arca-runtime", [artifact / "dotnet/dotnet.exe", "--list-runtimes"], artifact, env, logs))
    runtime_output = (logs / "arca-runtime.log").read_text()
    require("Microsoft.NETCore.App 10.0.11" in runtime_output and "Microsoft.AspNetCore.App 10.0.11" in runtime_output,
            "Portable .NET runtime incomplete")
    require((artifact / "dotnet/shared").as_posix() in runtime_output.replace("\\", "/"),
            "Runtime executable did not resolve its own shipped frameworks")
    notices = []
    for row in json.loads((assets / "nuget-runtime.lock.json").read_text())["packages"]:
        archive = plain(Path(cfg["runtime_packages"]) / row["archive"])
        require(sha(archive) == row["sha256"], "NuGet archive drift")
        with zipfile.ZipFile(archive) as package:
            for legal in row["legal"]:
                data = package.read(legal["member"])
                require(hashlib.sha256(data).hexdigest() == legal["sha256"], "NuGet notice drift")
                dest = artifact / "dependency-notices" / row["name"] / legal["member"]
                dest.parent.mkdir(parents=True, exist_ok=True)
                dest.write_bytes(data)
        if "supplemental_legal" in row:
            legal = row["supplemental_legal"]
            source = relative(assets, legal["input"])
            require(sha(source) == legal["sha256"], "Supplemental MIT notice drift")
            dest = artifact / "dependency-notices" / row["name"] / legal["file"]
            dest.parent.mkdir(parents=True, exist_ok=True)
            shutil.copyfile(source, dest)
        notices.append(row)
    (artifact / "nuget-dependency-notices.json").write_bytes(raw_json({"packages": notices}))
    source = artifact / "generated-source"
    source.mkdir()
    for rel in ["Generated/WsaaReference.cs", "Generated/WsfeV1Reference.cs", "GENERATION_RECEIPT.json",
                "packages.lock.json", "Elite.Arca.Wsfe.Generated.csproj"]:
        path = source / rel
        path.parent.mkdir(parents=True, exist_ok=True)
        shutil.copyfile(generated / rel, path)
    shutil.copyfile(src / "docs/ARCA_LOCAL_REFERENCE.md", artifact / "START_ARCA.md")
    receipt = {"schema": "elite-arca-reference-artifact/v1", "scope": "LOCAL_FIXTURES",
               "runtime_files": len(runtime_lock["files"]), "runtime_lock_sha256": sha(assets / "dotnet-runtime.lock.json"),
               "generation_receipt_sha256": sha(generated / "GENERATION_RECEIPT.json"),
               "provider_state": "CONDITIONED_USER_CREDENTIALS", "production_admitted": False}
    (artifact / "ARCA_REFERENCE.json").write_bytes(raw_json(receipt))
    return {"state": "PASS", "steps": steps, "receipt": receipt, "fresh_nuget_cache": True}
````


V402 composed delta: V402328 fixes FAIL974: exact original corresponding source plus separately preserved, narrowly checked Next route type derivation. Unexpected source mutation fails closed; no executable app bundle change or dependency update.

V402 composed delta: V402330 fixes signed extraction boundaries978/979: complete OSV SPDX identity coverage and recreate empty parent directories before runtime junctions; actual437package coverage and extracted329activation/rollback qualified. AUTHORED glue, no dependency change.

V402331 metadata correction only:19notice source labels now point to their existing exact official SOURCE_NOTICES receipts, including the three declared ADAPTED renderings. All materialized payload bytes are unchanged; current signed330 source identity remains exact.
