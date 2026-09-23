#!/usr/bin/env python3
"""AUTHORED glue: exact, one-shot local offline installation with pinned Node/pnpm."""
from __future__ import annotations
import argparse
import base64
import os
from pathlib import Path
import subprocess
import sys
import time
import plan_install as legacy
from select_artifact import SelectionError, decode, digest, json_bytes, no_reparse, read_bounded, require

POLICY_SHA256 = '77956707bda8be688bbbf01b72b9182d255f882c70f99a47a68eb64204452b1d'
PLAN = 'local-runtime-plan.json'
CATALOG = 'local-tool-notices.json'
NOTICE = 'THIRD_PARTY_NOTICES.md'
POLICY = 'local-runtime-policy.json'
STARTED = 'execution-started.json'
RESULT = 'execution-result.json'
LOG = 'execution.log'
MAX_SECONDS = 300


def policy():
    raw = read_bounded(Path(__file__).with_name(POLICY), 262144)
    require(digest(raw) == POLICY_SHA256, 'local runtime policy changed')
    return decode(raw), raw


def verify_payload(contained, expected):
    contained = no_reparse(contained)
    require(contained.is_dir(), 'contained payload directory missing')
    seen = set()
    allowed_dirs = {str(p).replace('\\', '/') for rel in expected for p in Path(rel).parents if str(p) != '.'}
    for root, dirs, files in os.walk(contained, followlinks=False):
        for name in dirs:
            p = no_reparse(Path(root) / name)
            require(p.relative_to(contained).as_posix() in allowed_dirs, 'unexpected payload directory')
        for name in files:
            p = no_reparse(Path(root) / name); rel = p.relative_to(contained).as_posix()
            require(rel in expected, 'unexpected payload file')
            identity = expected[rel]; raw = read_bounded(p, identity['bytes'])
            require(len(raw) == identity['bytes'] and digest(raw) == identity['sha256'], 'contained payload drift: ' + rel)
            seen.add(rel)
    require(seen == set(expected), 'contained payload incomplete')
    return contained


def notice_bytes(raw):
    catalog = decode(raw)
    chunks = [b'# pnpm local tool notices\n\n',
              b'Keep the accompanying JSON catalogue, original 47 retained notices, source delivery and adaptation notice. '
              b'This tool is acquired and adapted locally. Redistribution of the altered pnpm payload is outside this admission.\n\n']
    for entry in catalog['entries']:
        chunks.append((entry['package'] + '@' + entry['version'] + '\n' + entry['artifact_url'] + '\n' + entry['artifact_sha256'] + '\n').encode())
        for ref in entry['texts']:
            chunks.append((ref['sha256'] + ' | ' + ref['locator'] + ' | ' + ref['scope'] + '\n').encode())
        chunks.append(b'\n')
    for h, item in sorted(catalog['texts'].items()):
        content = base64.b64decode(item['content_base64'], validate=True)
        require(len(content) == item['bytes'] and digest(content) == h, 'notice text identity mismatch')
        chunks.extend([('\nBEGIN ORIGINAL TEXT SHA256 ' + h + '\n').encode(), content, b'\nEND ORIGINAL TEXT\n'])
    return b''.join(chunks)


def expected_plan(consumer, project, projection, contained, node, store, cache, target, acquisition_receipt, acquisition_sha256):
    p, policy_raw = policy(); target = no_reparse(target)
    contained = verify_payload(contained, p['files'])
    for other in (project, projection, node, store, cache, target, acquisition_receipt):
        other = no_reparse(other)
        require(contained != other and contained not in other.parents and other not in contained.parents, 'contained payload overlaps an input/output')
    base = legacy.recipe(consumer, project, projection, node, store, cache, target / 'legacy', acquisition_receipt, acquisition_sha256)
    require(base['node_sha256'] == p['node_sha256'], 'Node policy mismatch')
    catalog = read_bounded(Path(__file__).with_name(CATALOG), 4194304)
    require(digest(catalog) == p['notices_sha256'], 'catalogue changed')
    notices = notice_bytes(catalog)
    argv = [base['argv'][0], *p['node_flags'], str(contained / 'payload/bin/pnpm.mjs'), *base['argv'][2:]]
    env = dict(base['environment'])
    env['npm_config_userconfig'] = str(target / 'legacy' / legacy.EMPTY_FILES[0])
    env['npm_config_globalconfig'] = str(target / 'legacy' / legacy.EMPTY_FILES[1])
    value = {'schema': 'elite-pnpm-local-runtime-plan/v1', 'scope': p['scope'],
             'consumer': consumer, 'inputs': {**base['inputs'], 'contained': str(contained)},
             'legacy_plan_sha256': digest(json_bytes(base)), 'policy_sha256': POLICY_SHA256,
             'notices_sha256': digest(catalog), 'readable_notices_sha256': digest(notices),
             'consumer_files': base['consumer_files'], 'cwd': base['cwd'], 'argv': argv,
             'environment': env, 'environment_mode': 'REPLACE_NOT_MERGE', 'shell': False,
             'timeout_seconds': MAX_SECONDS, 'runtime_admitted': True, 'redistribution_admitted': False,
             'executed': False, 'conditions': p['conditions']}
    return value, base, {CATALOG: catalog, NOTICE: notices, POLICY: policy_raw}


def prepare(**args):
    target = no_reparse(args['target'])
    require(target.parent.is_dir() and not os.path.lexists(target), 'new runtime plan requires absent target')
    value, base, deliveries = expected_plan(**args)
    target.mkdir()  # Exclusive; preserve every partial output on any subsequent failure.
    inp = base['inputs']
    legacy.create_plan(base['consumer'], inp['project'], inp['projection'], inp['node'], inp['store'], inp['cache'],
                       target / 'legacy', inp['acquisition_receipt'], inp['acquisition_receipt_sha256'])
    for name, raw in deliveries.items():
        with (target / name).open('xb') as stream: stream.write(raw)
    raw = json_bytes(value)
    with (target / PLAN).open('xb') as stream: stream.write(raw)
    verify_plan(target, digest(raw))
    return digest(raw)


def verify_plan(target, sha256):
    target = no_reparse(target); raw = read_bounded(target / PLAN, 65536)
    require(digest(raw) == sha256, 'external runtime plan hash mismatch')
    value = decode(raw); inp = value['inputs']
    expected, _, deliveries = expected_plan(value['consumer'], inp['project'], inp['projection'], inp['contained'],
        inp['node'], inp['store'], inp['cache'], target, inp['acquisition_receipt'], inp['acquisition_receipt_sha256'])
    require(raw == json_bytes(expected), 'runtime plan is not the admitted recipe')
    require({p.name for p in target.iterdir()} == {'legacy', PLAN, *deliveries}, 'occupied or unexpected runtime plan entries')
    legacy.verify_plan(target / 'legacy', value['legacy_plan_sha256'])
    for name, expected_raw in deliveries.items():
        require(read_bounded(target / name, len(expected_raw)) == expected_raw, 'delivered runtime notice/policy changed')
    return value


def execute(target, sha256):
    target = no_reparse(target); value = verify_plan(target, sha256)
    with (target / STARTED).open('xb') as stream:
        stream.write(json_bytes({'schema': 'elite-pnpm-local-execution-start/v1', 'plan_sha256': sha256, 'one_shot': True}))
        stream.flush(); os.fsync(stream.fileno())
    started = time.monotonic(); result = {'schema': 'elite-pnpm-local-execution-result/v1', 'plan_sha256': sha256,
        'scope': value['scope'], 'state': 'FAIL', 'executed': True, 'redistribution_admitted': False}
    try:
        with (target / LOG).open('xb') as log:
            process = subprocess.run(value['argv'], cwd=value['cwd'], env=value['environment'], shell=False,
                stdout=log, stderr=subprocess.STDOUT, timeout=MAX_SECONDS, creationflags=subprocess.CREATE_NO_WINDOW)
        result['exit_code'] = process.returncode
        require(process.returncode == 0, 'restricted offline install failed; preserve plan, log and partial consumer')
        p, _ = policy(); verify_payload(value['inputs']['contained'], p['files'])
        require(digest(read_bounded(value['inputs']['node'], 150000000)) == p['node_sha256'], 'Node changed during execution')
        for name, h in value['consumer_files'].items():
            require(digest(read_bounded(Path(value['cwd']) / name, 2097152)) == h, 'consumer input changed during execution')
        log_raw = read_bounded(target / LOG, 16777216)
        result.update(state='PASS', log_sha256=digest(log_raw), consumer_files_unchanged=True,
                      node_flags=p['node_flags'], payload_unchanged=True)
    except (OSError, subprocess.TimeoutExpired, SelectionError) as exc:
        result['error'] = str(exc)
        raise
    finally:
        result['elapsed_seconds'] = round(time.monotonic() - started, 6)
        with (target / RESULT).open('xb') as stream:
            stream.write(json_bytes(result)); stream.flush(); os.fsync(stream.fileno())
    return result


def main(argv=None):
    parser = argparse.ArgumentParser(description=__doc__); sub = parser.add_subparsers(dest='action', required=True)
    make = sub.add_parser('prepare'); make.add_argument('--consumer', required=True, choices=sorted(legacy.PROFILES))
    for name in ('project', 'projection', 'contained', 'node', 'store', 'cache', 'target', 'acquisition-receipt', 'acquisition-sha256'):
        make.add_argument('--' + name, required=True)
    for action in ('verify', 'execute'):
        command = sub.add_parser(action); command.add_argument('--target', required=True); command.add_argument('--sha256', required=True)
    args = vars(parser.parse_args(argv)); action = args.pop('action')
    try:
        result = prepare(**args) if action == 'prepare' else verify_plan(**args) if action == 'verify' else execute(**args)
        print('PNPM_LOCAL_RUNTIME_PASS ' + (result if isinstance(result, str) else action))
        return 0
    except (SelectionError, OSError, ValueError, KeyError, TypeError, RecursionError, subprocess.TimeoutExpired) as exc:
        print('PNPM_LOCAL_RUNTIME_BLOCKED: ' + str(exc), file=sys.stderr); return 2


if __name__ == '__main__':
    raise SystemExit(main())
