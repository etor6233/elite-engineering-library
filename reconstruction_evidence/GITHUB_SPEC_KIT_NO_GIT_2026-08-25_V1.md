# GitHub Spec Kit no-Git reconstruction evidence V1

Verified at: 2026-08-25.

## Upstream identity

- repository: `github/spec-kit`;
- release: `v1.0.1`;
- commit: `9118ed15a0ba65053469a94c560ea5d233f75884`;
- source archive: `https://github.com/github/spec-kit/archive/9118ed15a0ba65053469a94c560ea5d233f75884.zip`;
- archive bytes: `4566121`;
- archive SHA-256: `ff8674156011f4d551ced8327b6bac0ad0d3378eecbf60f98c0edac04baea981`;
- license: MIT;
- upstream `LICENSE` SHA-256: `2510b446bc1f0cf9702453075d20cd88631e20e5642658edb7325d9c1eb534f7`.

## Environment

- Microsoft Windows 11 Pro `10.0.26200`, 64-bit;
- PowerShell `7.6.4`;
- Git `2.53.0.windows.2` was present for the upstream test suite, but was not initialized in either pilot target;
- Python `3.12.13` for the isolated installation;
- package identity observed: `specify-cli 1.0.1`.

## Real empty-directory pilots

The exact upstream CLI was invoked twice, once for Codex and once for Claude:

```text
specify init --here --force --non-interactive --ignore-agent-tools --integration codex --script ps
specify init --here --force --non-interactive --ignore-agent-tools --integration claude --script ps
```

Both commands exited `0`. Each pilot produced 30 files. Neither target contained a `.git` directory. The Codex target contained the upstream `.agents/skills` integration and the Claude target contained the upstream `.claude/skills` integration; both contained upstream `.specify` scripts, templates and workflow assets.

This proves that the exact official release can initialize its agent harness without making Git a project prerequisite. It does not prove the quality of a later generated product.

## Upstream tests

The release collected 7,318 tests. A focused Windows suite covering CLI initialization, integrations, events and security paths returned:

```text
240 passed
14 failed
6 skipped
```

All 14 focused failures occurred while their fixtures attempted to create Windows symbolic links and received `WinError 1314` before the code under test ran. They are retained as failures; they were not relabeled as passes. Core tests observed passing included no automatic Git extension installation and Codex/Claude initialization.

## Admission decision

`REBUILD_VERIFIED / CONDITIONED` for the narrow claim: official no-Git Codex/Claude harness initialization. Conditions are Python, the exact upstream release, preservation of the MIT notice, and a platform-specific rerun of symlink-security tests when Windows Developer Mode or equivalent privilege is available. No claim is made that Spec Kit supplies REVESTEX domain code, production security or business correctness.
