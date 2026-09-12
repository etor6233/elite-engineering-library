# Reconstruction evidence — V106 web XSS primitives and rejected SAST candidates

Date: 2026-08-29. Scope: exact public-source admission for a zero-additional-cost private project, followed by clean Markdown reconstruction. No repository, provider account, paid service, global tool installation or GitHub Actions minutes were used.

## Admitted Google runtime

Google `safevalues` 1.2.0 is Apache-2.0 and has zero runtime dependencies. Exact identity:

- repository/tag commit `0a900a2dc1d3dce6f28ab28aca8d0a5fd63b1c5b`, tree `6dbe6f600961c45d3be95a5c5bed4bb6e28ea91e`; the GitHub API reports the commit `unsigned`;
- source ZIP 230,696 bytes, SHA-256 `a32720089e3809bb5925e92c463ac7656b95c07c6ce4ee376b4d7a708a34471e`;
- source `LICENSE` SHA-256 `cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30` and `package.json` SHA-256 `7a1721957830789eec16ef1324fa2101f3dddae5146147ffdc153bc3afa0f31d`;
- npm tarball 90,639 bytes, SHA-256 `eff631d170cdcda898ce50064e868e1252ed923638655333a77a1bae58308f15`, registry integrity `sha512-zIsuhjYvJCjfsfjoim2ab6gLKFYAnTiDSJGh0cC3T44L/4kNLL90hBG2BzrXPrHA3f8Ms8FSJ1mljKH5dVR1cw==`.

`TS-GO-API-WEB-BRIDGE` advances to 0.2.1. The generated project exposes the admitted SafeValues builders/DOM wrappers from `src/platform/security/google-safevalues.ts` and contains three attack regressions: HTML escaping, refusal of a `javascript:` anchor without overwriting its safe value, and preservation of an HTTPS URL. Those two files are transparently `AUTHORED` integration/test code; the imported runtime remains Google's exact npm package and is never presented as copied local Google source.

The acquisition core advances to 0.4.68: 114 exact sources, 14 fail-closed profiles and 15 sources in `enterprise-platform-leaders`. SafeValues remains `PINNED_CANDIDATE_CONDITIONED`: it protects only paths that use its builders/wrappers. React escaping remains preferred, direct sinks and raw HTML are not retroactively safe, and final Trusted Types/CSP policy plus whole-application XSS testing remain project gates.

## Candidates correctly rejected

- GitHub CodeQL CLI 2.25.5 is not the universal private, zero-cost baseline requested: official terms require public code or the applicable GitHub Code Security/Advanced Security entitlement. GitHub Team alone is not treated as that entitlement.
- Microsoft DevSkim 1.0.90 executes from its signed NuGet package, but ships `SharpCompress 0.40.0` inside recursive extraction. OSV reports GHSA-6c8g-7p36-r338/CVE-2026-44788, fixed in 0.48.0. It is rejected, not locally patched and renamed official.
- Microsoft Application Inspector 1.10.1 is an immutable release at GitHub-verified commit `b1bbf91373488fb8e42573a417baad11de8d3890`. Its Windows ZIP exactly matches Microsoft's 8,366,289-byte/SHA-256 `935420bc6ee9f89fa6eaed01773f490bea3c48039e59f36e82e217cf492ae5ff`, yet its graph also contains `SharpCompress 0.40.0` and `Microsoft.CST.RecursiveExtractor 1.2.45`; it is rejected. It characterizes features rather than proving code safe in any case.
- Google `safety-web` is excluded because its official README calls the current `0.4.1-alpha.15` project unsupported and not production-ready.

These decisions are recorded as `UP-FAIL-187` through `UP-FAIL-191`. Tooling and process recoveries are `LIB-FAIL-1161` through `LIB-FAIL-1178`. Rejection is a result, not missing work hidden as PASS.

## Clean executable reconstruction

The governing structural verifier passed after every update:

```text
VERIFY_LIBRARY_PASS
packs=73 materialized_files=704 markdown_files=435
profile=ENTERPRISE_WEB_PACK_PLAN.md implementation_files=65
```

The final Foundation run acquired no global Go installation. It used the exact official `go1.26.7.windows-amd64.zip`, 74,955,002 bytes/SHA-256 `f4f534a486e4bc3387fa18f08208f2f854b7aaea8a08f2a2d829a914a05abb11`, passed explicitly to the runner. From clean Markdown composition:

- Go backend tests/build gates passed;
- frozen web install resolved 74 packages and included `safevalues 1.2.0`;
- Next type generation and strict TypeScript check passed;
- six Vitest files / 17 tests passed;
- Next 16.3.2 production build passed with eleven generated pages/routes;
- production `pnpm audit` reported zero findings for the inspected graph;
- Microsoft Playwright 1.62.1 passed runtime 4/4 and enterprise target 4/4 across Chromium desktop/mobile, Firefox and WebKit;
- Google Lighthouse 13.4.1 passed five governed target runs: performance median 0.99, accessibility 1.00, best practices 1.00 and SEO 1.00;
- final result: `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Foundation`.

The subsequent final `Audit` run also completed with exit code 0 after the acquisition core was advanced to 114/114 exact sources and all 36 composition profiles were realigned. Thus the source-lock verification, all-pack materialization, profile composition and bounded executable gates were rerun over the final V106 state rather than inherited from V105.

The portable V106 distribution was then expanded into a new empty temporary root and its embedded `VERIFY_LIBRARY.ps1` independently returned 73 packs, 704 materialized files, 435 Markdown files, backend 95, web 65 and `PORTABLE_SMOKE_CLEAN_PASS`. The archive carries a per-entry internal SHA-256 manifest plus an external archive checksum; temporary extraction roots were removed after verification.

This is executable evidence for the bounded foundation. Provider effects, PostgreSQL integration/PITR/restore, authenticated role journeys, assistive-technology accessibility, offensive security, representative load, deployment/canary/rollback and business acceptance still require the selected project's real environments and are not converted into local PASS.
