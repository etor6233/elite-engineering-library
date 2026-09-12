# Dependency License Evidence — reconstruction evidence V1

Date: 2026-08-24  
Status: `REBUILD_VERIFIED / PROJECT_CONDITIONED`

## Tooling admission

- Google `go-licenses` repository and release `v2.0.1` were inspected; the tool is Apache-2.0 and explicitly described upstream as not an officially supported Google product.
- pnpm 11/12 official documentation confirms `licenses list`, `--json`, `--long` and `--prod` for installed packages.
- CycloneDX 1.7 distinguishes declared, concluded and observed license evidence; a public repository or package declaration is not treated as concluded legal approval.

## Go executable result

`go-licenses v2.0.1` was installed at a temporary pinned tool path and executed against `./cmd/electromobility-api` from the clean 16-pack backend profile.

- report: 10 dependency/license records;
- expressions observed: `Apache-2.0`, `BSD-3-Clause`, `MIT`;
- allowlist check for `Apache-2.0,BSD-3-Clause,MIT,ISC`: PASS;
- `go-licenses save`: PASS and produced 11 LICENSE/NOTICE artifacts;
- versioned license URLs were produced for go-oidc 3.20.0, go-jose 4.1.4, pgx 5.10.0 and transitive runtime modules.

## Web artifact result

`pnpm licenses list --prod --json --long` passed against the clean frozen web profile. Observed expressions were:

- `0BSD`;
- `Apache-2.0`;
- `BSD-3-Clause`;
- `CC-BY-4.0`;
- `ISC`;
- `MIT`;
- `Apache-2.0 AND LGPL-3.0-or-later` for the optional Windows sharp/libvips binary package.

The last expression is intentionally not pre-acknowledged by the example policy. The selected deployment must either satisfy and record its redistribution obligations, use an artifact topology that does not distribute it, or select a different admitted image path. The library does not invent legal approval.

## Reusable gate

`DEPENDENCY-LICENSE-EVIDENCE-CORE 0.1.0` materializes three files. Its two regression tests pass and prove:

- deterministic JSON evidence with SHA-256 hashes;
- pnpm license artifact discovery;
- Go report plus saved-notice admission;
- denied-expression rejection;
- unreviewed-expression rejection;
- missing-license-artifact rejection.

The backend profile now materializes 16 packs/95 implementation files plus its record. The web profile materializes 4 packs/83 implementation files plus its record; frozen offline install, typecheck, 31 tests, Next build and license-gate tests pass.

## Condition

This is reproducible engineering evidence, not legal advice. Final distribution still requires an organization-approved policy, actual artifact inventory for each OS/architecture, complete notice inclusion and any source/offer obligations attached to the selected licenses.
