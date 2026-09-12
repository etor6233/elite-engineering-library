# Node Official Runtime Advisory Gate

## 1. Metadata

```yaml
pack_id: "NODE-OFFICIAL-RUNTIME-ADVISORY-GATE"
pack_version: "0.1.0"
status:
  authority: ELITE_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Ejecuta offline el motor oficial Node fijado con integridad/freshness/límites y receipt sobre Node24.20.0 Windows x64 observado; falla cerrado ante corpus vacío/alterado, fuente vencida y error."
stacks: ["CPython 3.14", "Node 24.20.0 Windows x64", "semver 7.8.5"]
compatible_with: ["inputs exactos admitidos, licencia íntegra y snapshot vigente"]
incompatible_with: ["otras versiones/plataformas sin reauditoría", "OSV nativo por alias npm/node", "monitor persistente", "SCA integral del runtime o despliegue"]
license_expression: "LicenseRef-Workspace-Owner AND MIT"
upstream_sources: ["https://github.com/nodejs/is-my-node-vulnerable/tree/c37a56bad56e34fe5223ddd3cb223cc4158136ae", "https://github.com/npm/node-semver/tree/6e05b7637396ac66522cff8731f07cfe0ef49a29", "https://nodejs.org/dist/v24.20.0/SHASUMS256.txt"]
verified_at: "2026-09-08"
```

## 2. Applicability

Usar como gate puntual cuando la capability requerida es consultar advisories
Node core/EOL del ejecutable exacto. El proyecto demuestra inputs/licencias,
freshness, retención/triage y ejecuta primero la suite. No activa monitoring,
no declara limpio npm/pnpm/OS/otros componentes ni admite producción.

## 3. Architecture contract

- Motor oficial ADAPTED sólo en acceso local a snapshots; matching sin cambios.
- Wrapping AUTHORED: versión/hash/licencias, perfil deshabilitado por defecto,
  snapshots exactos vigentes, archivos y árboles limitados, entorno mínimo.
- Bridge y motor vinculados al source-lock; bytes ya verificados se copian a
  staging aislado. Sin acceso HTTP del motor, sin matcher propio ni auto-refresh.
- Child finito y output acotado; errores estáticos, no causas privadas.
- Receipt nuevo con scope/identidades; rollback no oculta los findings previos.

## 4. Exact file manifest

```text
CREATE node_runtime_advisory_gate/profile.example.json
CREATE node_runtime_advisory_gate/source-lock.json
CREATE node_runtime_advisory_gate/run_gate.py
CREATE node_runtime_advisory_gate/bridge.cjs
CREATE node_runtime_advisory_gate/verify_pack.py
CREATE node_runtime_advisory_gate/README.md
CREATE node_runtime_advisory_gate/THIRD_PARTY_NOTICES.md
CREATE node_runtime_advisory_gate/official/is-vulnerable.js
CREATE node_runtime_advisory_gate/official/ascii.js
CREATE node_runtime_advisory_gate/official/LICENSE
```

## 5. Materialization blocks

### FILE: `node_runtime_advisory_gate/profile.example.json`
```yaml
block_id: "NODE-OFFICIAL-RUNTIME-ADVISORY-GATE:file-1:v1"
operation: CREATE
provenance: AUTHORED
source: "local adapter governed by the pinned official Node API and library admission/freshness contracts"
license: "LicenseRef-Workspace-Owner"
sha256: "29cb59f281cd9fa40d8d994a951d43935eb39ff31cffea373edcba11e019d881"
variables: []
secrets_allowed: false
```
````json
{
  "schema_version": 1,
  "enabled": false,
  "max_age_days": 7
}
````

### FILE: `node_runtime_advisory_gate/source-lock.json`
```yaml
block_id: "NODE-OFFICIAL-RUNTIME-ADVISORY-GATE:file-2:v1"
operation: CREATE
provenance: AUTHORED
source: "local adapter governed by the pinned official Node API and library admission/freshness contracts"
license: "LicenseRef-Workspace-Owner"
sha256: "98ee8520ec08f0d6dd988a5d477402dfdb99f9b13ae0b1c2a27c6767a70db247"
variables: []
secrets_allowed: false
```
````json
{
  "schema_version": 1,
  "source": {
    "repository": "https://github.com/nodejs/is-my-node-vulnerable",
    "version": "1.6.1",
    "commit": "c37a56bad56e34fe5223ddd3cb223cc4158136ae",
    "license": "MIT",
    "original_sha256": "0531884fd5f7693f01e499b6a95ab8893b357e2975c06c1ea2531ded59a256c4",
    "adaptation": "getJson reads verified isolated local snapshot only; official matching/platform/EOL unchanged",
    "files": {
      "is-vulnerable.js": "769e22cd956780430210535f562d1d31cb6fda5cdfe22790980b0e9595678c46",
      "ascii.js": "a22ca24d5ca9dd023219f32ed668359a15493ac6e351d3640538067ceb4761bc",
      "LICENSE": "bbf7df376bad39f8be96f9044ad58b77f0a48cbd1b482ecc5893916f7129b91c"
    }
  },
  "semver": {
    "version": "7.8.5",
    "license": "ISC",
    "files": {
      "bin/semver.js": "bd6c871026985937dc945011fc54b74b47f5998154216d7b8d40d5ed782e4402",
      "classes/comparator.js": "054202956430d63d5ff4599fae09760ce465b489e4f0b5ef5ce7cc7ac21157ac",
      "classes/index.js": "3bb69280c2a788d0eb16f915bb9df4dbe812075182024c753dca2283bcea1b17",
      "classes/range.js": "a6544485aa8575aa9854c67e2caa67726815de18a482d0ceb8a1003244de3bc1",
      "classes/semver.js": "813b2c185512d30c9b931c2fed8140889d9d7490ceeeedc370da30337dc8ca57",
      "functions/clean.js": "4eadb0892844cf3ae295121a86163a66c73f89acd1b7f0b114ec115b4539512f",
      "functions/cmp.js": "a63d74e87b73788e78e9ce0a4892b5333d6b809c0de88b31e4ed76cbf17f94b3",
      "functions/coerce.js": "28a251c5ab210ddf9e97551b9f37a53329fcce91f0c3943dcbc02de1a1de915a",
      "functions/compare-build.js": "5ab651d5b40af289bd85c645a92b6d8cfe1a986dc413c797cfcc8d623d7c844c",
      "functions/compare-loose.js": "07b6a3a1db0a5210ceb784c1708bd4679f3a94fc73a9c9eb349349e7070a6f78",
      "functions/compare.js": "d404b5aa48aaddc8a654c5da8fb7d4443404b7948589b21ac4b045d1cee4e34c",
      "functions/diff.js": "7a11fd39b987cdf06c65e928cb1aca49bec583feb86fb5c8fe47a6fd61d7de31",
      "functions/eq.js": "b6e30a7168e52723216fc163d300e2bbabf92ec0251f9ac5438bb6ccf57c8936",
      "functions/gt.js": "135523704aa48cd98834dd170ee9f74f0e68043b379f32d021db11e6304c5c93",
      "functions/gte.js": "991c5bbe48ecb210a562646872f05862ad9fc0d42186d85aa60bdc6fa323eb9d",
      "functions/inc.js": "952069fc8690b7d3af0fe9d55f7c54fe2ac067b48c5e74f6a54f9ce19a334493",
      "functions/lt.js": "1c897a9bc849320e2e9dc0f6c09555c01ee3ddf30734515d717b88ad7740ea25",
      "functions/lte.js": "3a8d0b1d00423f60cc7cb810b36ab77b61330831ad237dbe73eb5cecfa412800",
      "functions/major.js": "5c678480d882f511200fed2c16ec3847dfedb08a1d70328dc2f031d35d825276",
      "functions/minor.js": "1b051794f1713adec2a236517196691687c35f82e0b596ff3316a78b3cc10ae6",
      "functions/neq.js": "a662883751918822c162183b46b9e20d09489132f82686c92ab78bee67f3a127",
      "functions/parse.js": "29a69e15b6d02fe381d573f861881a89590e9d0f0f0ca740c5f85eaf0234c4ad",
      "functions/patch.js": "f02d3c1b059fe3d96ce124886f7eef321d381a95638fd3c4a8d5ccd8e76ffadd",
      "functions/prerelease.js": "baedbf503d5610ad041bfb56071efa48feb331ca278295c399537d35d3ffb593",
      "functions/rcompare.js": "84fa5e88adf08d15c993cac8cdec6d1a65045b7e95a9c55184230a7f807f4dc0",
      "functions/rsort.js": "6ec659ce3b6c2b173c719286caab04409adba046c0917874ec3b5e36ddfbc7e3",
      "functions/satisfies.js": "8cf5e122b757251671ed6c9d9680904b71cd375845853f05312e608cf2cc2946",
      "functions/sort.js": "c2fe2d3ed0be8a4e9de8f02abdfd5d9c0d3bcc510d88e87af84185592882c4b8",
      "functions/truncate.js": "e145833a250311927c92f424f56498c52f1b7bb1271ce6afbf006db87773ff43",
      "functions/valid.js": "0de7ea736cb7807179d46dfac09c830a9338e02b6a07db12d7040cde2def6025",
      "index.js": "4b3e57d3d40e29e0706002eba113d09f35aea593578376bbeec83b777b9912ab",
      "internal/constants.js": "38a112baf27ceca0260082ff26ac2fd7a9861cab1af12dd65e720277f68e6ce9",
      "internal/debug.js": "8a9f420572260f3cf944463b5090d62a60f0730589dc23a7ec4ca25e2ee41bb3",
      "internal/identifiers.js": "b4916b09dc7869ae0eb05e71c855c942a3a7365f7e6e89185d59a9e45a2451d7",
      "internal/lrucache.js": "14d087c87da87b6f5c36fc4cdd7d2d14077874b14a68e20fce5b6138fa2ca34f",
      "internal/parse-options.js": "fdf51d0de8d5442c35a997ef58cd530d239ce206f961d14c5121354451b01d01",
      "internal/re.js": "5833262888e2b5d843a69193f83c05e374818dbe55379b497819b5bf58e48cd8",
      "LICENSE": "4ec3d4c66cd87f5c8d8ad911b10f99bf27cb00cdfcff82621956e379186b016b",
      "package.json": "7c94cb7f2a53c27b20d76386ec144c062894dbcc909cfabd0f728c37874b1776",
      "preload.js": "edb6808911bebcb324b2df57e5c9935149e56984ff083b74c6cfe215f5b710ba",
      "range.bnf": "15d0baf9b7b98e6d862c0cb9e822d2533d2bc23e136bb75314d802ca1fcb0392",
      "ranges/gtr.js": "8fadad28e36d28e93d498ac7ac20badba2a407312845eabc18e82e90a0732b19",
      "ranges/intersects.js": "fe87ac5d3020010ad3ec00636dadbf0c669ff07d0f57e0a8165a8264f79a676f",
      "ranges/ltr.js": "e5186fcc03018acf9be6d968755d4c49727aeb0d981d179eb568ae5fbe983038",
      "ranges/max-satisfying.js": "e1a2c0d6144cc772cd20bbc8ecb9e8a3a4074e9172a3d8e794838b591cdeb416",
      "ranges/min-satisfying.js": "2681abf54098aa670f12826b76a6ec77a2441186ae4243afde3be8ae4908f7ca",
      "ranges/min-version.js": "fefba0a88c2bf74d5cede504b5ae50a8dc3edfd69cf0174a491e2cf3e442614b",
      "ranges/outside.js": "3a2b0b23593d2f49419c06af2af75450cab103b0c25d665d48fe5bca495a21ca",
      "ranges/simplify.js": "7b78581c13322bc68ece2088685386b2a9b51c15b94d0a2063bdf2546bd41934",
      "ranges/subset.js": "ceb9eb6ed5bfc6c7ac7af7f47e5c3535444a4e66a0e2b58cdad4bd02a35d454d",
      "ranges/to-comparators.js": "6c5e966210cff270fa2850668aaf8460fac7759f8d99f282521ef7a78f4564e9",
      "ranges/valid.js": "5ef6f995af801868925940cf8df5735d565ebabb090b068695cae65218bcd3ac",
      "README.md": "f1a789dcec285150be24db2ea04dd3175031554fa9834ec92fab83fb5e025a57"
    },
    "source_url": "https://registry.npmjs.org/semver/-/semver-7.8.5.tgz",
    "artifact_sha256": "d85045d4300d7d57c891336b95df532e73f34c22ffcd222452b6d08b9d127d5d",
    "artifact_integrity": "sha512-Y7/KDsb8LjooZpwaqGyulO6DQlksgCncchHGk+sZIY4SBvUocMBEFH5Ur1fI4dV+Jvl0w6cjvucaIi40puRioA==",
    "source_commit": "6e05b7637396ac66522cff8731f07cfe0ef49a29"
  },
  "runtime": {
    "version": "v24.20.0",
    "platform": "win32",
    "sha256": "5c976096e04e5c2c1f091938926234cc9fbebfe9787ddd149351b3b0ecc707b5",
    "license_sha256": "ed34dd8e3f0a78dbaf00d0444ce8e285b015b765379c2e17880455f70370f8e9"
  },
  "snapshots": {
    "observed_at": "2026-09-08T03:16:29.470083-03:00",
    "max_age_days": 7,
    "core_commit": "21bf6214b10d3a70250828c05b5e7026f4e27f50",
    "schedule_commit": "72fdab20216c5f04e0a0fe72a225c2504e9f2b42",
    "files": {
      "security.json": "b4e6c24e61c9c0b4923e2c79256cfb938c1835fe86901eec60a395d50c9d56f2",
      "schedule.json": "1cf0432ceb9dfde7f1fd4cce43206519942cfdfad5a26039c2f4bb20fde8549c"
    },
    "core_records": 194
  },
  "bridge_sha256": "2e6a524bf2ae602a02f5fdf2ef59310e174490540d37afc56b559f90e91b0fb6"
}
````

### FILE: `node_runtime_advisory_gate/run_gate.py`
```yaml
block_id: "NODE-OFFICIAL-RUNTIME-ADVISORY-GATE:file-3:v1"
operation: CREATE
provenance: AUTHORED
source: "local adapter governed by the pinned official Node API and library admission/freshness contracts"
license: "LicenseRef-Workspace-Owner"
sha256: "510cfe27a052ba803d8748091772f2727212fd0ae20a55e6a7fb27db4b587123"
variables: []
secrets_allowed: false
```
````python
"""AUTHORED fail-closed adapter around the pinned official Node advisory engine.

This is an offline one-shot admission gate, not a scanner engine or scheduler.
"""
from __future__ import annotations
import argparse
import datetime as dt
import hashlib
import json
import os
from pathlib import Path, PurePosixPath
import subprocess
import tempfile
import threading

LOCK_SHA256 = '98ee8520ec08f0d6dd988a5d477402dfdb99f9b13ae0b1c2a27c6767a70db247'
ROOT = Path(__file__).resolve().parent
MAX_FILE = 1_048_576
MAX_BUNDLE = 4_194_304
MAX_OUTPUT = 65_536


class GateError(Exception):
    pass


def reject(code: str) -> None:
    raise GateError(code)


def local_path(path: Path) -> Path:
    path = path.absolute()
    if str(path).startswith(('\\\\', '//')):
        reject('NONLOCAL_PATH')
    for item in (path, *path.parents):
        if item.is_symlink() or item.is_junction():
            reject('LINK_PATH')
    return path


def read_limited(path: Path, maximum: int = MAX_FILE) -> bytes:
    path = local_path(path)
    if not path.is_file():
        reject('MISSING_FILE')
    with path.open('rb') as handle:
        data = handle.read(maximum + 1)
    if len(data) > maximum:
        reject('INPUT_LIMIT')
    return data


def parse_json(data: bytes):
    def unique(pairs):
        value = {}
        for key, item in pairs:
            if key in value:
                reject('DUPLICATE_JSON_KEY')
            value[key] = item
        return value
    try:
        return json.loads(data, object_pairs_hook=unique)
    except (UnicodeError, ValueError):
        reject('INVALID_JSON')


def checked_file(path: Path, expected: str) -> bytes:
    data = read_limited(path)
    if hashlib.sha256(data).hexdigest() != expected:
        reject('INTEGRITY_MISMATCH')
    return data


def source_lock(root: Path = ROOT) -> dict:
    data = checked_file(root / 'source-lock.json', LOCK_SHA256)
    lock = parse_json(data)
    if not isinstance(lock, dict) or lock.get('schema_version') != 1:
        reject('INVALID_LOCK')
    return lock


def profile_value(path: Path) -> tuple[dict, str]:
    raw = read_limited(path, 4096)
    profile = parse_json(raw)
    if not isinstance(profile, dict) or set(profile) != {'schema_version', 'enabled', 'max_age_days'}:
        reject('INVALID_PROFILE')
    if type(profile['schema_version']) is not int or profile['schema_version'] != 1:
        reject('INVALID_PROFILE')
    if profile['enabled'] is not True:
        reject('DISABLED')
    if type(profile['max_age_days']) is not int or not 1 <= profile['max_age_days'] <= 7:
        reject('INVALID_FRESHNESS_WINDOW')
    return profile, hashlib.sha256(raw).hexdigest()


def validate_snapshots(directory: Path, lock: dict, max_age_days: int,
                       now: dt.datetime | None = None) -> dict[str, bytes]:
    current = now or dt.datetime.now(dt.timezone.utc)
    observed = dt.datetime.fromisoformat(lock['snapshots']['observed_at'])
    age = current - observed
    if age < -dt.timedelta(minutes=5):
        reject('FUTURE_SNAPSHOT')
    if age > dt.timedelta(days=min(max_age_days, lock['snapshots']['max_age_days'])):
        reject('STALE_SNAPSHOT')
    values = {name: checked_file(directory / name, digest)
              for name, digest in lock['snapshots']['files'].items()}
    core = parse_json(values['security.json'])
    schedule = parse_json(values['schedule.json'])
    if not isinstance(core, dict) or len(core) != lock['snapshots']['core_records'] or not core:
        reject('INVALID_CORPUS')
    if any(not isinstance(v, dict) or not isinstance(v.get('vulnerable'), str)
           or not isinstance(v.get('patched'), str) for v in core.values()):
        reject('INVALID_CORPUS')
    if not isinstance(schedule, dict) or not schedule.get('v24', {}).get('end'):
        reject('INVALID_SCHEDULE')
    return values


def validate_semver(directory: Path, lock: dict) -> dict[str, bytes]:
    directory = local_path(directory)
    if not directory.is_dir():
        reject('MISSING_SEMVER')
    expected = lock['semver']['files']
    actual = set()
    for count, child in enumerate(directory.rglob('*'), start=1):
        if count > 256:
            reject('SEMVER_TREE_LIMIT')
        local_path(child)
        if child.is_file():
            actual.add(child.relative_to(directory).as_posix())
            if len(actual) > len(expected):
                reject('SEMVER_FILE_SET')
    if actual != set(expected):
        reject('SEMVER_FILE_SET')
    values = {}
    total = 0
    for name, digest in expected.items():
        rel = PurePosixPath(name)
        if rel.is_absolute() or '..' in rel.parts or ':' in name or '\\' in name:
            reject('INVALID_SEMVER_PATH')
        value = checked_file(directory / name, digest)
        total += len(value)
        if total > MAX_BUNDLE:
            reject('BUNDLE_LIMIT')
        values[name] = value
    manifest = parse_json(values['package.json'])
    if manifest.get('version') != lock['semver']['version'] or not values.get('LICENSE'):
        reject('SEMVER_IDENTITY')
    return values


def clean_environment() -> dict[str, str]:
    allowed = {'SYSTEMROOT', 'WINDIR', 'TEMP', 'TMP'}
    return {key: value for key, value in os.environ.items() if key.upper() in allowed}


def bounded_process(args: list[str], cwd: Path, timeout: float = 20) -> bytes:
    """Bound wall time and output while keeping child stderr out of gate receipts."""
    flags = subprocess.CREATE_NO_WINDOW if os.name == 'nt' else 0
    process = subprocess.Popen(args, cwd=cwd, env=clean_environment(), stdout=subprocess.PIPE,
                               stderr=subprocess.STDOUT, creationflags=flags)
    results = []
    def read_output():
        data = process.stdout.read(MAX_OUTPUT + 1)
        results.append(data)
        if len(data) > MAX_OUTPUT and process.poll() is None:
            process.kill()
    reader = threading.Thread(target=read_output, daemon=True)
    reader.start()
    timed_out = False
    try:
        process.wait(timeout=timeout)
    except subprocess.TimeoutExpired:
        timed_out = True
        process.kill()
        process.wait(timeout=5)
    finally:
        reader.join(timeout=5)
        process.stdout.close()
    if reader.is_alive() or not results:
        reject('PROCESS_OUTPUT_UNAVAILABLE')
    if timed_out:
        reject('PROCESS_TIMEOUT')
    if len(results[0]) > MAX_OUTPUT:
        reject('PROCESS_OUTPUT_LIMIT')
    if process.returncode != 0:
        reject('ENGINE_FAILED')
    return results[0]


def runtime_identity(node: Path, lock: dict) -> dict:
    node = local_path(node)
    if not node.is_file():
        reject('MISSING_RUNTIME')
    digest = hashlib.sha256()
    with node.open('rb') as stream:
        total = 0
        while chunk := stream.read(MAX_FILE):
            total += len(chunk)
            if total > 150_000_000:
                reject('RUNTIME_LIMIT')
            digest.update(chunk)
    if digest.hexdigest() != lock['runtime']['sha256']:
        reject('RUNTIME_IDENTITY')
    checked_file(node.parent / 'LICENSE', lock['runtime']['license_sha256'])
    result = parse_json(bounded_process([str(node), '-e',
        'process.stdout.write(JSON.stringify({version:process.version,platform:process.platform}))'],
        node.parent, timeout=5))
    if result != {key: lock['runtime'][key] for key in ('version', 'platform')}:
        reject('RUNTIME_IDENTITY')
    return result


def evaluate(node: Path, version: str, platform: str, lock: dict,
             snapshots: dict[str, bytes], semver: dict[str, bytes], root: Path = ROOT) -> dict:
    # Copy only already-verified bytes, so changes to input files do not affect this run.
    sources = {name: checked_file(root / 'official' / name, digest)
               for name, digest in lock['source']['files'].items()}
    bridge = checked_file(root / 'bridge.cjs', lock['bridge_sha256'])
    with tempfile.TemporaryDirectory(prefix='elite-node-advisory-') as scratch:
        work = Path(scratch)
        (work / 'official').mkdir()
        for name, value in sources.items():
            (work / 'official' / name).write_bytes(value)
        for name, value in snapshots.items():
            (work / 'official' / name).write_bytes(value)
        for name, value in semver.items():
            target = work / 'official/node_modules/semver' / name
            target.parent.mkdir(parents=True, exist_ok=True)
            target.write_bytes(value)
        (work / 'bridge.cjs').write_bytes(bridge)
        result = parse_json(bounded_process([str(node), '--max-old-space-size=64',
                            str(work / 'bridge.cjs'), version, platform], work))
    if not isinstance(result, dict) or set(result) != {'version', 'platform', 'vulnerable'}:
        reject('INVALID_ENGINE_RESULT')
    if result['version'] != version or result['platform'] != platform or type(result['vulnerable']) is not bool:
        reject('INVALID_ENGINE_RESULT')
    return result


def run(profile: Path, node: Path, semver_root: Path, snapshot_dir: Path,
        receipt: Path, root: Path = ROOT) -> dict:
    receipt = local_path(receipt)
    if receipt.exists() or not receipt.parent.is_dir():
        reject('RECEIPT_DESTINATION')
    config, profile_sha = profile_value(profile)
    lock = source_lock(root)
    snapshots = validate_snapshots(snapshot_dir, lock, config['max_age_days'])
    semver = validate_semver(semver_root, lock)
    identity = runtime_identity(node, lock)
    result = evaluate(node.absolute(), identity['version'].removeprefix('v'), identity['platform'],
                      lock, snapshots, semver, root)
    output = {'schema_version': 1, 'gate': 'NODE-OFFICIAL-RUNTIME-ADVISORY-GATE',
              'status': 'FINDINGS' if result['vulnerable'] else 'PASS',
              'runtime_version': identity['version'], 'runtime_platform': identity['platform'],
              'runtime_sha256': lock['runtime']['sha256'], 'source_lock_sha256': LOCK_SHA256,
              'profile_sha256': profile_sha, 'max_age_days': config['max_age_days'],
              'snapshot_sha256': lock['snapshots']['files'],
              'snapshot_observed_at': lock['snapshots']['observed_at'],
              'checked_at': dt.datetime.now(dt.timezone.utc).isoformat(),
              'scope': 'known official Node core advisories and EOL only; not npm, embedded component SCA or deployment monitoring',
              'production_admitted': False}
    with receipt.open('x', encoding='utf-8', newline='\n') as stream:
        json.dump(output, stream, indent=2)
        stream.write('\n')
    return output


def main() -> int:
    parser = argparse.ArgumentParser()
    for name in ('profile', 'node', 'semver-root', 'snapshot-dir', 'receipt'):
        parser.add_argument('--' + name, type=Path, required=True)
    args = parser.parse_args()
    try:
        result = run(args.profile, args.node, args.semver_root, args.snapshot_dir, args.receipt)
        print('NODE_RUNTIME_ADVISORY_' + result['status'])
        return 1 if result['status'] == 'FINDINGS' else 0
    except GateError as error:
        print('NODE_RUNTIME_ADVISORY_ERROR ' + str(error))
        return 2
    except (OSError, ValueError, KeyError, TypeError, subprocess.SubprocessError):
        print('NODE_RUNTIME_ADVISORY_ERROR INPUT_OR_PROCESS_FAILURE')
        return 2


if __name__ == '__main__':
    raise SystemExit(main())
````

### FILE: `node_runtime_advisory_gate/bridge.cjs`
```yaml
block_id: "NODE-OFFICIAL-RUNTIME-ADVISORY-GATE:file-4:v1"
operation: CREATE
provenance: AUTHORED
source: "local adapter governed by the pinned official Node API and library admission/freshness contracts"
license: "LicenseRef-Workspace-Owner"
sha256: "2e6a524bf2ae602a02f5fdf2ef59310e174490540d37afc56b559f90e91b0fb6"
variables: []
secrets_allowed: false
```
````javascript
// AUTHORED adapter. Matching and EOL logic are the pinned official implementation.
'use strict';
require('node:https').request = () => { throw new Error('NETWORK_DISABLED'); };
require('node:http').request = () => { throw new Error('NETWORK_DISABLED'); };
const {isNodeVulnerable} = require('./official/is-vulnerable');
const [version, platform] = process.argv.slice(2);
isNodeVulnerable(version, platform).then(vulnerable => {
  if (typeof vulnerable !== 'boolean') throw new Error('INVALID_RESULT');
  process.stdout.write(JSON.stringify({version, platform, vulnerable}));
}).catch(() => {
  process.stderr.write('NODE_ADVISORY_ENGINE_FAILED');
  process.exitCode = 2;
});
````

### FILE: `node_runtime_advisory_gate/verify_pack.py`
```yaml
block_id: "NODE-OFFICIAL-RUNTIME-ADVISORY-GATE:file-5:v1"
operation: CREATE
provenance: AUTHORED
source: "local adapter governed by the pinned official Node API and library admission/freshness contracts"
license: "LicenseRef-Workspace-Owner"
sha256: "6368e937709d430fc3a7993de68b57d97146645a8af602d10ea63c9227f5debe"
variables: []
secrets_allowed: false
```
````python
"""AUTHORED regression suite; requires the exact admitted external inputs."""
import argparse
import copy
import datetime as dt
import json
import os
from pathlib import Path
import shutil
import tempfile
import time
import unittest
from unittest import mock
import run_gate as gate

INPUTS = None


class NodeGateTests(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory(prefix='elite-node-gate-test-')
        self.addCleanup(self.temp.cleanup)
        self.work = Path(self.temp.name)
        self.root = self.work / 'gate'
        shutil.copytree(Path(__file__).parent, self.root, ignore=shutil.ignore_patterns('__pycache__'))
        self.snapshots = self.work / 'snapshots'
        shutil.copytree(INPUTS.snapshot_dir, self.snapshots)
        self.semver = self.work / 'semver'
        shutil.copytree(INPUTS.semver_root, self.semver)
        self.profile = self.work / 'profile.json'
        self.profile.write_text(json.dumps({'schema_version': 1, 'enabled': True, 'max_age_days': 7}), encoding='utf-8')
        self.node = INPUTS.node
        self.receipt = self.work / 'receipt.json'
        self.lock = gate.source_lock(self.root)

    def run_gate(self):
        return gate.run(self.profile, self.node, self.semver, self.snapshots, self.receipt, self.root)

    def error(self, code, action=None):
        with self.assertRaisesRegex(gate.GateError, '^' + code + '$'):
            (action or self.run_gate)()
        self.assertFalse(self.receipt.exists())

    def test_current_runtime_receipt(self):
        result = self.run_gate()
        self.assertEqual(result['status'], 'PASS')
        self.assertEqual(result['runtime_version'], 'v24.20.0')
        self.assertFalse(result['production_admitted'])
        self.assertEqual(result, json.loads(self.receipt.read_text()))
        self.assertEqual(result['source_lock_sha256'], gate.LOCK_SHA256)

    def test_official_matching_vulnerable_control(self):
        snapshots = gate.validate_snapshots(self.snapshots, self.lock, 7)
        semver = gate.validate_semver(self.semver, self.lock)
        result = gate.evaluate(self.node, '24.14.1', 'win32', self.lock, snapshots, semver, self.root)
        self.assertTrue(result['vulnerable'])

    def test_official_invalid_platform(self):
        snapshots = gate.validate_snapshots(self.snapshots, self.lock, 7)
        semver = gate.validate_semver(self.semver, self.lock)
        self.error('ENGINE_FAILED', lambda: gate.evaluate(self.node, '24.20.0', 'invalid', self.lock, snapshots, semver, self.root))

    def test_disabled_profile(self):
        shutil.copyfile(self.root / 'profile.example.json', self.profile)
        self.error('DISABLED')

    def test_profile_unknown_key(self):
        self.profile.write_text('{"schema_version":1,"enabled":true,"max_age_days":7,"ignore":true}')
        self.error('INVALID_PROFILE')

    def test_profile_duplicate_key(self):
        self.profile.write_text('{"schema_version":1,"enabled":true,"enabled":false,"max_age_days":7}')
        self.error('DUPLICATE_JSON_KEY')

    def test_profile_window_not_boolean(self):
        self.profile.write_text('{"schema_version":1,"enabled":true,"max_age_days":true}')
        self.error('INVALID_FRESHNESS_WINDOW')

    def test_profile_no_freshness_extension(self):
        self.profile.write_text('{"schema_version":1,"enabled":true,"max_age_days":365}')
        self.error('INVALID_FRESHNESS_WINDOW')

    def test_empty_corpus_never_clean(self):
        (self.snapshots / 'security.json').write_text('{}')
        self.error('INTEGRITY_MISMATCH')

    def test_truncated_corpus_never_clean(self):
        (self.snapshots / 'security.json').write_text('{')
        self.error('INTEGRITY_MISMATCH')

    def test_missing_snapshot(self):
        (self.snapshots / 'schedule.json').unlink()
        self.error('MISSING_FILE')

    def test_oversized_snapshot(self):
        (self.snapshots / 'security.json').write_bytes(b'x' * (gate.MAX_FILE + 1))
        self.error('INPUT_LIMIT')

    def test_expired_source_snapshot(self):
        when = dt.datetime.fromisoformat(self.lock['snapshots']['observed_at']) + dt.timedelta(days=8)
        self.error('STALE_SNAPSHOT', lambda: gate.validate_snapshots(self.snapshots, self.lock, 7, when))

    def test_future_source_snapshot(self):
        when = dt.datetime.fromisoformat(self.lock['snapshots']['observed_at']) - dt.timedelta(minutes=6)
        self.error('FUTURE_SNAPSHOT', lambda: gate.validate_snapshots(self.snapshots, self.lock, 7, when))

    def test_operator_can_shorten_window(self):
        when = dt.datetime.fromisoformat(self.lock['snapshots']['observed_at']) + dt.timedelta(days=2)
        self.error('STALE_SNAPSHOT', lambda: gate.validate_snapshots(self.snapshots, self.lock, 1, when))

    def test_semver_tamper(self):
        (self.semver / 'functions/satisfies.js').write_text('module.exports=()=>false')
        self.error('INTEGRITY_MISMATCH')

    def test_semver_extra_file(self):
        (self.semver / 'unexpected.js').write_text('')
        self.error('SEMVER_FILE_SET')

    def test_semver_missing_notice(self):
        (self.semver / 'LICENSE').unlink()
        self.error('SEMVER_FILE_SET')

    def test_semver_directory_enumeration_bound(self):
        for i in range(257):
            (self.semver / f'extra-empty-{i}').mkdir()
        self.error('SEMVER_TREE_LIMIT')

    def test_official_engine_tamper(self):
        (self.root / 'official/is-vulnerable.js').write_text('module.exports={isNodeVulnerable:async()=>false}')
        self.error('INTEGRITY_MISMATCH')

    def test_bridge_tamper(self):
        (self.root / 'bridge.cjs').write_text("process.stdout.write(JSON.stringify({version:'24.20.0',platform:'win32',vulnerable:false}))")
        self.error('INTEGRITY_MISMATCH')

    def test_lock_tamper(self):
        lock = copy.deepcopy(self.lock)
        lock['snapshots']['observed_at'] = dt.datetime.now(dt.timezone.utc).isoformat()
        (self.root / 'source-lock.json').write_text(json.dumps(lock))
        self.error('INTEGRITY_MISMATCH')

    def test_runtime_identity_mismatch(self):
        self.node = self.work / 'node.exe'
        self.node.write_bytes(b'not-the-verified-executable')
        self.error('RUNTIME_IDENTITY')

    def test_receipt_never_overwritten(self):
        self.receipt.write_text('existing-evidence')
        with self.assertRaisesRegex(gate.GateError, '^RECEIPT_DESTINATION$'):
            self.run_gate()
        self.assertEqual(self.receipt.read_text(), 'existing-evidence')

    def test_runtime_license_required(self):
        runtime = self.work / 'runtime'
        runtime.mkdir()
        self.node = runtime / 'node.exe'
        shutil.copyfile(INPUTS.node, self.node)
        self.error('MISSING_FILE')

    def test_environment_options_are_not_inherited(self):
        with mock.patch.dict(os.environ, {'NODE_OPTIONS': '--require=PRIVATE_NOT_EXECUTED', 'NODE_PATH': 'PRIVATE_NOT_USED', 'DEBUG': '1'}):
            result = self.run_gate()
        self.assertEqual(result['status'], 'PASS')
        self.assertNotIn('PRIVATE', self.receipt.read_text())

    def test_junction_boundary_rejected_before_read(self):
        target = self.snapshots / 'security.json'
        original = Path.is_junction
        with mock.patch.object(Path, 'is_junction', lambda path: path == target or original(path)):
            self.error('LINK_PATH')

    def test_real_child_timeout(self):
        started = time.monotonic()
        self.error('PROCESS_TIMEOUT', lambda: gate.bounded_process([str(self.node), '-e', 'setInterval(()=>{},1000)'], self.work, .3))
        self.assertLess(time.monotonic() - started, 4)

    def test_real_child_output_limit(self):
        self.error('PROCESS_OUTPUT_LIMIT', lambda: gate.bounded_process([str(self.node), '-e', "process.stdout.write('x'.repeat(70000))"], self.work, 5))

    def test_real_child_private_failure_not_returned(self):
        with self.assertRaises(gate.GateError) as caught:
            gate.bounded_process([str(self.node), '-e', "process.stderr.write('PRIVATE_CAUSE');process.exit(1)"], self.work)
        self.assertEqual(str(caught.exception), 'ENGINE_FAILED')


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('--node', type=Path, required=True)
    parser.add_argument('--semver-root', type=Path, required=True)
    parser.add_argument('--snapshot-dir', type=Path, required=True)
    INPUTS = parser.parse_args()
    suite = unittest.defaultTestLoader.loadTestsFromTestCase(NodeGateTests)
    result = unittest.TextTestRunner(verbosity=2).run(suite)
    raise SystemExit(0 if result.wasSuccessful() and not result.skipped else 1)
````

### FILE: `node_runtime_advisory_gate/README.md`
```yaml
block_id: "NODE-OFFICIAL-RUNTIME-ADVISORY-GATE:file-6:v1"
operation: CREATE
provenance: AUTHORED
source: "local adapter governed by the pinned official Node API and library admission/freshness contracts"
license: "LicenseRef-Workspace-Owner"
sha256: "b5c83e4e571f364ad8f64412dd62da69d33c533f1abfbb1a5b04f505883e1d99"
variables: []
secrets_allowed: false
```
````markdown
# Official Node runtime advisory gate

This one-shot offline gate uses the pinned official Node vulnerability matcher.
It checks known Node core advisories and EOL for the observed, exact Windows x64
Node24.20.0 executable. It does not scan npm, pnpm, all embedded components, the
OS, first-party code, malware or deployment configuration. PASS is not production
admission. The official upstream network CLI remains unsuitable without guards.

Inputs must already have an acquisition/admission record. Obtain the standalone
Node executable from https://nodejs.org/dist/v24.20.0/win-x64/node.exe and keep its
unmodified LICENSE alongside it. Verify the official signed checksums and trust
root in the library maintenance-tool lock; this runner rechecks the executable,
version/platform and license digests. Do not replace the global installation.
The full Node ZIP was rejected for bundled npm advisories; this gate does not
remediate that archive or make it safe to redistribute.

Acquire semver7.8.5 from the exact npm tarball in source-lock.json, preserve its
ISC LICENSE and extract its 53 files unchanged into a dedicated directory.
Do not install packages or run lifecycle scripts to supply this input. The runner
checks every file digest and rejects missing/extra files and links. The copy
bundled in the Node ZIP has different bytes (line endings and README whitespace)
and will be rejected; do not normalize bytes to bypass the lock.

Provide two local files in a snapshot directory:

- security.json: https://raw.githubusercontent.com/nodejs/security-wg/21bf6214b10d3a70250828c05b5e7026f4e27f50/vuln/core/index.json
- schedule.json: https://raw.githubusercontent.com/nodejs/Release/72fdab20216c5f04e0a0fe72a225c2504e9f2b42/schedule.json

The source-lock contains their digests and the recorded observation time. It is
bound into the runner. The maximum age is seven days; a profile may shorten it.
Expired data, future observation time, empty or changed data fail before matching.
Refreshing a snapshot requires current official research, preserving before/after,
updating the source-lock and runner binding, repeating gates and publishing a new
pack revision. Changing only a timestamp is not a legitimate refresh.

Copy profile.example.json to a project-owned profile, set enabled=true after its
inputs are present, and invoke with explicit absolute paths:

```text
python verify_pack.py --node <NODE_EXE> --semver-root <SEMVER_DIR> --snapshot-dir <SNAPSHOT_DIR>
python run_gate.py --profile <PROFILE_JSON> --node <NODE_EXE> --semver-root <SEMVER_DIR> --snapshot-dir <SNAPSHOT_DIR> --receipt <NEW_RECEIPT_JSON>
```

Use CPython3.14 on Windows and the exact Node executable. No network runs inside
the gate. Official source, bridge and semver are checked before being copied into
an isolated temporary directory. Only verified bytes are executed; user NODE_OPTIONS,
NODE_PATH and DEBUG are not inherited. The matching child has a64MiB V8 heap
budget,20s wall deadline and64KiB combined-output limit. Identity probing has5s.
Each input file is limited to1MiB, semver to4MiB and256 enumerated entries. These
are execution guards, not OS sandboxing or full process RSS limits.

Exit0 means the scoped known-advisory/EOL check passed. Exit1 means findings;
exit2 means disabled, invalid, stale, missing or failed input/execution. Neither
findings nor errors may be translated into a clean result. Receipts contain
version/platform, hashes, scope and time; no causes, paths, tokens or payloads.
A receipt path must be new and its parent must exist; never overwrite evidence.
Retain error logs/receipts under the project's evidence and retention policy.
No scheduler, daemon, notification destination or recurring scan is created.

The project owns triage, renewal, retained evidence and target acceptance. On
failure, open/update its failure lesson, inspect the official advisory and block
promotion. Reversal removes this gate selection while preserving receipts;
it never falls back to the unguarded CLI or labels a vulnerable old runtime safe.
````

### FILE: `node_runtime_advisory_gate/THIRD_PARTY_NOTICES.md`
```yaml
block_id: "NODE-OFFICIAL-RUNTIME-ADVISORY-GATE:file-7:v1"
operation: CREATE
provenance: AUTHORED
source: "local adapter governed by the pinned official Node API and library admission/freshness contracts"
license: "LicenseRef-Workspace-Owner"
sha256: "ecf555b3fb96298493d253287132691844ce592335909e5d2ccd2e6f0092d586"
variables: []
secrets_allowed: false
```
````markdown
# Source and license boundary

- official/is-vulnerable.js: ADAPTED from nodejs/is-my-node-vulnerable1.6.1,
  commit c37a56bad56e34fe5223ddd3cb223cc4158136ae, MIT. Only getJson was changed
  to read the verified isolated snapshot. Official matching, platform filtering
  and EOL evaluation are preserved; the adaptation is not an upstream release.
- official/ascii.js and official/LICENSE: VERBATIM from that same revision.
  The original copyright and MIT terms remain in official/LICENSE.
- semver7.8.5 is an external dependency under ISC, not embedded in this pack.
  Preserve its unmodified LICENSE and53 files. The exact npm artifact, source
  commit6e05b7637396ac66522cff8731f07cfe0ef49a29, SHA/integrity and file hashes
  are in source-lock.json. The tarball is not interchangeable byte for byte
  with the copy included by Node.
- Node24.20.0 and its full LICENSE are external runtime inputs. Node and bundled
  component notices remain authoritative; neither the binary nor npm is embedded.
- JSON snapshots are external official data inputs fixed by commit and hash.
  This pack includes their identities, not their full contents.
- Runner, bridge, tests, profile, lock and explanatory files are AUTHORED by the
  workspace owner, under LicenseRef-Workspace-Owner. They are not Node project
  code and do not imply endorsement by Node.js, OpenJS or npm.
````

### FILE: `node_runtime_advisory_gate/official/is-vulnerable.js`
```yaml
block_id: "NODE-OFFICIAL-RUNTIME-ADVISORY-GATE:file-8:v1"
operation: CREATE
provenance: ADAPTED
source: "https://raw.githubusercontent.com/nodejs/is-my-node-vulnerable/c37a56bad56e34fe5223ddd3cb223cc4158136ae/is-vulnerable.js"
license: "MIT"
sha256: "769e22cd956780430210535f562d1d31cb6fda5cdfe22790980b0e9595678c46"
variables: []
secrets_allowed: false
```
````javascript
const { danger, allGood, bold, vulnerableWarning, separator } = require('./ascii')
const { request } = require('https')
const fs = require('fs')
const path = require('path')
const satisfies = require('semver/functions/satisfies')

const STORE = {
  security: {
    url: 'https://raw.githubusercontent.com/nodejs/security-wg/main/vuln/core/index.json',
    jsonFile: path.join(__dirname, 'security.json'),
    etagFile: path.join(__dirname, 'security.etag'),
    etagValue: ''
  },
  schedule: {
    url: 'https://raw.githubusercontent.com/nodejs/Release/main/schedule.json',
    jsonFile: path.join(__dirname, 'schedule.json'),
    etagFile: path.join(__dirname, 'schedule.etag'),
    etagValue: ''
  }
}

async function readLocal (file) {
  return require(file)
}

function debug (msg) {
  if (process.env.DEBUG) {
    console.debug(msg)
  }
}

function loadETag () {
  for (const [key, obj] of Object.entries(STORE)) {
    if (fs.existsSync(obj.etagFile)) {
      debug(`Loading local ETag for '${key}'`)
      obj.etagValue = fs.readFileSync(obj.etagFile).toString()
    }
  }
}

async function fetchJson (obj) {
  await new Promise((resolve) => {
    request(obj.url, (res) => {
      if (res.statusCode !== 200) {
        console.error(`Request to Github returned http status ${res.statusCode}. Aborting...`)
        process.nextTick(() => { process.exit(1) })
      }

      const fileStream = fs.createWriteStream(obj.jsonFile)
      res.pipe(fileStream)

      fileStream.on('finish', () => {
        fileStream.close()
        resolve()
      })

      fileStream.on('error', (err) => {
        console.error(`Error ${err.message} while writing to '${obj.jsonFile}'. Aborting...`)
        process.nextTick(() => { process.exit(1) })
      })
    }).on('error', (err) => {
      console.error(`Request to Github returned error ${err.message}. Aborting...`)
      process.nextTick(() => { process.exit(1) })
    }).end()
  })
  return readLocal(obj.jsonFile)
}

// ADAPTED: the caller supplies hash-verified, fresh local snapshots.
// Matching, platform filtering and EOL evaluation remain the official code.
async function getJson (obj) {
  return readLocal(obj.jsonFile)
}

const checkPlatform = platform => {
  const availablePlatforms = ['aix', 'darwin', 'freebsd', 'linux', 'openbsd', 'sunos', 'win32', 'android']
  if (platform && !availablePlatforms.includes(platform)) {
    throw new Error(`platform ${platform} is not valid. Please use ${availablePlatforms.join(',')}.`)
  }
}

const isSystemAffected = (platform, affectedEnvironments) => {
  // No platform specified (legacy mode)
  if (!platform || !Array.isArray(affectedEnvironments)) {
    return true
  }
  // If the environment is matching or all the environments are affected
  if (affectedEnvironments.includes(platform) || affectedEnvironments.includes('all')) {
    return true
  }
  // Default to false
  return false
}

function getVulnerabilityList (currentVersion, data, platform) {
  const list = []
  for (const key in data) {
    const vuln = data[key]

    if (
      (
        satisfies(currentVersion, vuln.vulnerable) &&
        !satisfies(currentVersion, vuln.patched)
      ) && isSystemAffected(platform, vuln.affectedEnvironments)
    ) {
      const severity = vuln.severity === 'unknown' ? '' : `(${vuln.severity})`
      list.push(`${bold(vuln.cve)}${severity}: ${vuln.overview}\n${bold('Patched versions')}: ${vuln.patched}`)
    }
  }
  return list
}

async function cli (currentVersion, platform) {
  checkPlatform(platform)

  const isEOL = await isNodeEOL(currentVersion)
  if (isEOL) {
    console.error(danger)
    console.error(`${currentVersion} is end-of-life. There are high chances of being vulnerable. Please upgrade it.`)
    process.exit(1)
  }

  const securityJson = await getJson(STORE.security)
  const list = getVulnerabilityList(currentVersion, securityJson, platform)

  if (list.length) {
    console.error(danger)
    console.error(vulnerableWarning + '\n')
    console.error(`${list.join(`\n${separator}\n\n`)}\n${separator}`)
    process.exit(1)
  } else {
    console.info(allGood)
  }
}

async function getVersionInfo (version) {
  const scheduleJson = await getJson(STORE.schedule)

  if (scheduleJson[version.toLowerCase()]) {
    return scheduleJson[version.toLowerCase()]
  }

  for (const [key, value] of Object.entries(scheduleJson)) {
    if (satisfies(version, key)) {
      return value
    }
  }

  return null
}

/**
 * @param {string} version
 * @returns {Promise<boolean>} true if the version is end-of-life
 */
async function isNodeEOL (version) {
  const myVersionInfo = await getVersionInfo(version)

  if (!myVersionInfo) {
    // i.e. isNodeEOL('abcd') or isNodeEOL('lts') or isNodeEOL('99')
    throw Error(`Could not fetch version information for ${version}`)
  } else if (!myVersionInfo.end) {
    // We got a record, but..
    // v0.12.18 etc does not have an EOL date, which probably means too old.
    return true
  }

  const now = new Date()
  const end = new Date(myVersionInfo.end)
  return now > end
}

async function isNodeVulnerable (version, platform) {
  checkPlatform(platform)
  const isEOL = await isNodeEOL(version)
  if (isEOL) {
    return true
  }

  const coreIndex = await getJson(STORE.security)
  const list = getVulnerabilityList(version, coreIndex, platform)
  return list.length > 0
}

if (process.argv[2] !== '-r') {
  loadETag()
}

module.exports = {
  isNodeVulnerable,
  cli
}
````

### FILE: `node_runtime_advisory_gate/official/ascii.js`
```yaml
block_id: "NODE-OFFICIAL-RUNTIME-ADVISORY-GATE:file-9:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/nodejs/is-my-node-vulnerable/c37a56bad56e34fe5223ddd3cb223cc4158136ae/ascii.js"
license: "MIT"
sha256: "a22ca24d5ca9dd023219f32ed668359a15493ac6e351d3640538067ceb4761bc"
variables: []
secrets_allowed: false
```
````javascript
const util = require('util')

const danger = '\n' +
'\n' +
'██████   █████  ███    ██  ██████  ███████ ███████\n' +
'██   ██ ██   ██ ████   ██ ██       ██      ██   ██\n' +
'██   ██ ███████ ██ ██  ██ ██   ███ █████   ███████\n' +
'██   ██ ██   ██ ██  ██ ██ ██    ██ ██      ██   ██\n' +
'██████  ██   ██ ██   ████  ██████  ███████ ██   ██\n' +
'\n'

const allGood = '\n' +
'\n' +
' █████  ██      ██           ██████   ██████   ██████  ██████         ██\n' +
'██   ██ ██      ██          ██       ██    ██ ██    ██ ██   ██     ██  ██\n' +
'███████ ██      ██          ██   ███ ██    ██ ██    ██ ██   ██         ██\n' +
'██   ██ ██      ██          ██    ██ ██    ██ ██    ██ ██   ██     ██  ██\n' +
'██   ██ ███████ ███████      ██████   ██████   ██████  ██████         ██\n' +
'\n'

function escapeStyleCode (code) {
  return '\u001b[' + code + 'm'
}

function bold (text) {
  var left = ''
  var right = ''
  const formatCodes = util.inspect.colors.bold
  left += escapeStyleCode(formatCodes[0])
  right = escapeStyleCode(formatCodes[1]) + right
  return left + text + right
}

const vulnerableWarning = bold('The current Node.js version (' + process.version + ') is vulnerable to the following CVEs:')

var separator = '='
for (var i = 0; i < process.stdout.columns; ++i) {
  separator = separator + '='
}

module.exports.danger = danger
module.exports.allGood = allGood
module.exports.bold = bold
module.exports.vulnerableWarning = vulnerableWarning
module.exports.separator = separator
````

### FILE: `node_runtime_advisory_gate/official/LICENSE`
```yaml
block_id: "NODE-OFFICIAL-RUNTIME-ADVISORY-GATE:file-10:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/nodejs/is-my-node-vulnerable/c37a56bad56e34fe5223ddd3cb223cc4158136ae/LICENSE"
license: "MIT"
sha256: "bbf7df376bad39f8be96f9044ad58b77f0a48cbd1b482ecc5893916f7129b91c"
variables: []
secrets_allowed: false
```
````text
MIT License

Copyright (c) 2023 RafaelGSS

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

## 6. Configuration surface

| Campo | Tipo | Default | Validación | Secreto | Efecto |
|---|---|---|---|---|---|
| enabled | bool | false | true explícito | no | activa consulta puntual |
| max_age_days | int | 7 | 1..7, no bool; no extiende lock | no | ventana de freshness |
| node/semver-root/snapshot-dir | paths | ausentes | archivo/tree/hash/licencia local exactos | no | inputs declarados |
| receipt | path | ausente | padre existe, archivo nuevo | no | evidencia sin sobrescritura |

## 7. Dependency bill

| Dependencia | Pin | Licencia | Scope | Identidad |
|---|---|---|---|---|
| Node checker | 1.6.1/c37a56bad56e34fe5223ddd3cb223cc4158136ae | MIT | motor core | source-lock, adaptación getJson |
| semver | 7.8.5/6e05b7637396ac66522cff8731f07cfe0ef49a29 | ISC | comparación oficial | tarball y53 archivos externos exactos |
| Node | 24.20.0 Windows x64 | MIT + notices oficiales | runtime | exe+LICENSE SHA y process identity |
| CPython | 3.14.4 probado | PSF-2.0 | stdlib runner | runtime externo, no SCA global implícito |

## 8. Apply order

Adquirir inputs fijados con sus receipts/licencias, materializar en destino ausente,
ejecutar verify_pack.py y después run_gate.py con perfil propio habilitado.
El README contiene las URLs/versiones y todas las opciones. No sobrescribir un
proyecto existente. Retener los recibos al retirar el directorio. Actualizar datos
exige el contrato de freshness y un nuevo lock/runner/pack probado, no sólo fecha.

## 9. Verification

Suite30 pruebas: actual Node0, control vulnerable del motor, plataforma inválida,
vacío/truncado/oversize/future/stale, perfil inválido, alteración de motor/bridge/
semver/lock, licencias ausentes, receipt existente, env hostil y child real con
timeout/exceso de salida/causa privada. Frontera junction inyectada como test unitario,
no certificación OS de sandbox. SCA del grafo CLI2 paquetes/0matches conocidos.
No repetir browser/DB por este gate aislado: V319 retiene su evidencia por snapshot.

## 10. Reconstruction evidence

Expediente y restricciones: reconstruction_evidence/NODE_OFFICIAL_ADVISORY_GATE_V320.md.
Diez archivos,7 AUTHORED/1 ADAPTED/2 VERBATIM; plan focal1/10 fuera del perfil
de producto67/746. Verificar30 pruebas sobre árbol reconstruido e identidad10/10
antes de usar; G0–G8 sólo se cierra para el claim estrecho e inputs actuales.
