"""Build an isolated pnpm11.25.0 candidate with binary ZIP support removed.

AUTHORED containment, not an official pnpm fix or general runtime admission.
Only exact published input bytes accepted. No network or execution operation.
"""
import argparse
import difflib
import os
from pathlib import Path
import re
import shutil
import sys
import tempfile
import select_artifact as base

BUNDLE = 'dist/pnpm.mjs'
ORIGINAL_SHA = 'ddc64218bc85fb88d28b5def06eb01fafb39cb67f7a9465ce736456658cc7f11'
NOTICE = b'''Local adaptation: Elite pnpm11.25.0 ZIP-disabled candidate, revision1.
The published pnpm artifact is retained separately. This is not an official release.
All16 bundled adm-zip0.6.0 modules have been removed. Binary ZIP fetch/extraction
fails closed before download or creation of the extraction directory. TAR package
installation is unchanged. Automatic runtime acquisition through ZIP is unavailable.
Reason: GHSA-vwc7-r8mq-g2x9. No upstream fixed release was identified2026-09-11.
Original license and notice files remain. Additional source/license/relinking and
runtime admission conditions remain open. No production or redistribution grant.
'''
DENIAL = 'throw new PnpmError("ZIP_DISABLED_BY_ELITE_CONTAINMENT", "Binary ZIP acquisition is disabled in this local pnpm adaptation; use a separately admitted runtime.");'

def transform(raw):
    base.require(base.digest(raw) == ORIGINAL_SHA, 'published bundle identity mismatch')
    text = raw.decode('utf-8')
    chunks = list(re.finditer(r'^// .*?/node_modules/adm-zip/[^\n]+\n', text, re.M))
    base.require(len(chunks) == 16, 'adm-zip module inventory changed')
    start = chunks[0].start()
    after = text.index('\n// ', chunks[-1].end()) + 1
    removed = text[start:after]
    symbols = re.findall(r'^var (\w+) = __commonJS\(', removed, re.M)
    base.require(len(symbols) == 16, 'adm-zip wrapper inventory changed')
    text = text[:start] + '// Elite containment: adm-zip implementation removed.\n' + text[after:]
    for name, next_name in [('downloadAndUnpackZip', 'downloadWithIntegrityCheck'), ('extractZipToTarget', 'toStatelessTester')]:
        pattern = r'async function ' + name + r'\([^\n]+\) \{\n[\s\S]*?\n\}\n(?=(?:async )?function ' + next_name + r'\()'
        signature = re.search(pattern, text)
        base.require(signature is not None, 'extraction function boundary changed')
        replacement = signature.group().split('\n', 1)[0] + '\n  ' + DENIAL + '\n}\n'
        text, count = re.subn(pattern, lambda _: replacement, text)
        base.require(count == 1, 'ambiguous extraction function')
    old = '      case "zip": {\n        const tempLocation = await cafs.tempDir();'
    base.require(text.count(old) == 1, 'binary ZIP dispatch boundary changed')
    # Reject before even creating the extraction directory.
    end = text.index('      default: {', text.index(old))
    text = text[:text.index(old)] + '      case "zip": {\n        ' + DENIAL + '\n      }\n' + text[end:]
    old_init = '    import_adm_zip = __toESM(require_adm_zip(), 1);\n'
    base.require(text.count(old_init) == 1, 'ZIP initialization changed')
    text = text.replace(old_init, '')
    base.require(text.count('var import_adm_zip, import_ssri3;') == 1, 'ZIP import binding changed')
    text = text.replace('var import_adm_zip, import_ssri3;', 'var import_ssri3;')
    for symbol in symbols:
        base.require(not re.search(r'\b' + re.escape(symbol) + r'\b', text), 'removed wrapper still referenced: ' + symbol)
    base.require('/node_modules/adm-zip/' not in text and 'import_adm_zip' not in text, 'ZIP implementation still present')
    return text.encode('utf-8')

def expected_files(source):
    source = base.no_reparse(source)
    policy = base.load_policy()
    entries = {x['path']: x for x in policy['files']}
    seen = set()
    for directory, dirs, files in os.walk(source, followlinks=False):
        for name in dirs: base.no_reparse(Path(directory) / name)
        for name in files:
            path = Path(directory) / name
            rel = path.relative_to(source).as_posix()
            base.require(rel in entries, 'unexpected source file')
            entry = entries[rel]
            raw = base.read_bounded(path, entry['bytes'])
            base.require(len(raw) == entry['bytes'] and base.digest(raw) == entry['sha256'], 'source file identity mismatch: ' + rel)
            seen.add(rel)
    base.require(seen == set(entries), 'published source incomplete')
    original = base.read_bounded(source / BUNDLE, entries[BUNDLE]['bytes'])
    return derived_files(entries, original, lambda rel, size: base.read_bounded(source / rel, size))


def expected_projection(projection):
    """Reuse the selector's complete archive verification; excluded bytes are not needed twice."""
    projection = base.no_reparse(projection)
    policy = base.load_policy()
    base.verify(projection, policy)
    entries = {x['path']: x for x in policy['files']}
    original = base.read_bounded(projection / 'payload' / BUNDLE, entries[BUNDLE]['bytes'])
    return derived_files(entries, original, lambda rel, size: base.read_bounded(projection / 'payload' / rel, size))


def derived_files(entries, original, read_original):
    changed = transform(original)
    result = {}
    for rel, entry in entries.items():
        if entry['selection'] == 'unchanged':
            result['payload/' + rel] = changed if rel == BUNDLE else read_original(rel, entry['bytes'])
    diff = ''.join(difflib.unified_diff(original.decode().splitlines(True), changed.decode().splitlines(True), fromfile='official/dist/pnpm.mjs', tofile='elite-zip-disabled/dist/pnpm.mjs'))
    result['adaptation.diff'] = diff.encode('utf-8')
    result['ADAPTATION.txt'] = NOTICE
    result['containment-receipt.json'] = base.json_bytes({
        'schema':'elite-pnpm-zip-disabled/v1', 'provenance':'ADAPTED_AUTHORED_CONTAINMENT',
        'base_policy_sha256':base.POLICY_SHA256, 'base_bundle_sha256':ORIGINAL_SHA,
        'candidate_bundle_sha256':base.digest(changed), 'removed_modules':16,
        'changed_payload_files':[BUNDLE], 'unchanged_payload_files':441,
        'excluded_physical_files':13, 'binary_zip_enabled':False,
        'upstream_advisory':'GHSA-vwc7-r8mq-g2x9', 'runtime_admitted':False,
        'redistribution_admitted':False, 'executed_by_materializer':False,
        'files':{p:{'sha256':base.digest(b),'bytes':len(b)} for p,b in sorted(result.items())}})
    return result

def verify(target, expected):
    target = base.no_reparse(target)
    found = set()
    allowed_dirs = {p.as_posix() for rel in expected for p in Path(rel).parents if p.as_posix() != '.'}
    for directory, dirs, files in os.walk(target, followlinks=False):
        for name in dirs:
            path = base.no_reparse(Path(directory) / name)
            base.require(path.relative_to(target).as_posix() in allowed_dirs, 'unexpected candidate directory')
        for name in files:
            path = Path(directory) / name
            rel = path.relative_to(target).as_posix()
            base.require(rel in expected, 'unexpected candidate file')
            base.require(base.read_bounded(path, len(expected[rel])) == expected[rel], 'candidate bytes changed: ' + rel)
            found.add(rel)
    base.require(found == set(expected), 'candidate files missing')

def create(source, target, from_projection=False):
    base.require(os.name == 'nt', 'Windows publication required')
    target = base.no_reparse(target)
    base.require(target.parent.is_dir() and not os.path.lexists(target), 'target must be absent with existing parent')
    expected = expected_projection(source) if from_projection else expected_files(source)
    staging = Path(tempfile.mkdtemp(prefix='.pnpm-zip-disabled-', dir=target.parent))
    try:
        for rel, raw in expected.items():
            path = staging / rel
            path.parent.mkdir(parents=True, exist_ok=True)
            with path.open('xb') as stream: stream.write(raw)
        verify(staging, expected)
        base.no_reparse(target.parent)
        os.rename(staging, target)
    finally:
        if staging.exists():
            base.require(staging.parent == target.parent and staging.name.startswith('.pnpm-zip-disabled-'), 'cleanup scope rejected')
            shutil.rmtree(staging)
    verify(target, expected)

def main(argv=None):
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('action', choices=['create','verify'])
    origin = parser.add_mutually_exclusive_group(required=True)
    origin.add_argument('--source')
    origin.add_argument('--projection')
    parser.add_argument('--target', required=True)
    args = parser.parse_args(argv)
    try:
        source = args.projection if args.projection is not None else args.source
        if args.action == 'create': create(source, args.target, args.projection is not None)
        else: verify(args.target, expected_projection(source) if args.projection is not None else expected_files(source))
        print('PNPM_ZIP_CONTAINMENT_BYTES_PASS runtime_admitted=false redistribution_admitted=false')
        return 0
    except (base.SelectionError, OSError, ValueError, TypeError, KeyError) as exc:
        print('PNPM_ZIP_CONTAINMENT_BLOCKED: ' + str(exc), file=sys.stderr)
        return 2

if __name__ == '__main__': raise SystemExit(main())
