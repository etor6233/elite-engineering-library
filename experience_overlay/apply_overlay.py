"""AUTHORED hash-guarded V403 integration glue; no dependency install or deployment."""
from pathlib import Path
import argparse, hashlib, json, os, tempfile

def digest(data):
    return hashlib.sha256(data).hexdigest()

def safe(root, relative):
    rel = Path(relative)
    if rel.is_absolute() or '..' in rel.parts:
        raise ValueError('Unsafe relative path: ' + relative)
    target = root / rel
    if not target.resolve().is_relative_to(root.resolve()):
        raise ValueError('Path escapes target: ' + relative)
    for parent in [target, *target.parents]:
        if parent == root:
            break
        if parent.is_symlink() or (hasattr(parent, 'is_junction') and parent.is_junction()):
            raise ValueError('Link in affected path: ' + relative)
    return target

def atomic_write(destination, data):
    # Stage on the same filesystem: a failed write never truncates the destination.
    descriptor, temporary = tempfile.mkstemp(prefix='.ui-overlay-', dir=destination.parent)
    temp = Path(temporary)
    try:
        with os.fdopen(descriptor, 'wb') as stream:
            stream.write(data)
            stream.flush()
            os.fsync(stream.fileno())
        os.replace(temp, destination)
    finally:
        temp.unlink(missing_ok=True)

def apply(bundle, target, verify=False):
    target = target.resolve()
    if (target / 'AGENT_SYSTEM_START.md').exists():
        raise ValueError('Refusing canonical library root')
    manifest = json.loads((bundle / 'overlay-manifest.json').read_text(encoding='utf8'))
    if digest((target / 'package.json').read_bytes()) != manifest['base_package_sha256']:
        raise ValueError('Incompatible package.json; explicit compatibility review required')
    pending = []
    for row in manifest['files']:
        source = safe(bundle, row['bundle_path'])
        destination = safe(target, row['target'])
        data = source.read_bytes()
        if digest(data) != row['after_sha256']:
            raise ValueError('Bundle source hash mismatch: ' + row['target'])
        before = destination.read_bytes() if destination.exists() else None
        observed = digest(before) if before is not None else None
        if observed not in (row['before_sha256'], row['after_sha256']):
            raise ValueError('Target revision conflict: ' + row['target'])
        if verify and observed != row['after_sha256']:
            raise ValueError('Overlay missing: ' + row['target'])
        pending.append((destination, before, data, row))
    written = []
    if not verify:
        lock = target / '.experience-overlay.lock'
        descriptor = os.open(lock, os.O_CREAT | os.O_EXCL | os.O_WRONLY, 0o600)
        try:
            os.write(descriptor, str(os.getpid()).encode())
            for destination, before, data, row in pending:
                if before == data:
                    continue
                # Recheck immediately before mutation; no concurrent writers are allowed.
                now = destination.read_bytes() if destination.exists() else None
                if now != before:
                    raise ValueError('Concurrent modification: ' + row['target'])
                destination.parent.mkdir(parents=True, exist_ok=True)
                written.append((destination, before))
                atomic_write(destination, data)
        except Exception:
            for destination, before in reversed(written):
                if before is None:
                    destination.unlink(missing_ok=True)
                elif destination.read_bytes() != before:
                    atomic_write(destination, before)
            raise
        finally:
            os.close(descriptor)
            lock.unlink()
    for destination, _, data, row in pending:
        if destination.read_bytes() != data:
            raise ValueError('Postcondition mismatch: ' + row['target'])
    return {'schema_version': '1.0.0', 'revision': manifest['revision'],
            'claim': 'HASH_LOCKED_UI_OVERLAY', 'status': 'PASS',
            'production_authorized': False, 'target': str(target),
            'files_verified': len(pending), 'files_changed': len(written),
            'manifest_sha256': digest((bundle / 'overlay-manifest.json').read_bytes()),
            'files': [{'path': r['target'], 'sha256': r['after_sha256']} for _, _, _, r in pending]}

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--target', required=True, type=Path)
    parser.add_argument('--report', required=True, type=Path)
    parser.add_argument('--verify', action='store_true')
    args = parser.parse_args()
    result = apply(Path(__file__).resolve().parent, args.target, args.verify)
    args.report.parent.mkdir(parents=True, exist_ok=True)
    args.report.write_text(json.dumps(result, indent=2) + '\n', encoding='utf8')
    print(json.dumps({'status': result['status'], 'files_verified': result['files_verified']}))
