# Microsoft BCQuality Agent Knowledge Evidence — 2026-08-25 V1

## Claim

This evidence proves exact acquisition and local validation of Microsoft BCQuality 1.5, an official machine-readable knowledge and skills library for Business Central development agents. It is directly relevant to making an agent review and generate AL with Microsoft-published Business Central-specific rules. It does not claim functional ERP coverage or replace BCApps product tests.

## Exact source

```text
repository: microsoft/BCQuality
tag: v1.5
commit: ad8ccde5953fd6ce072e3b59049406783066ed60
archive URL: https://github.com/microsoft/BCQuality/archive/ad8ccde5953fd6ce072e3b59049406783066ed60.zip
archive bytes: 714396
archive SHA-256: 7db8ccddfad1b4a48a45c9dfa4d9b88519f795201fd7fd0f47e52cf0bdd37e7c
license: MIT
LICENSE SHA-256: 9515a04280cf50ce154eb3ab8fa375c7cc5e6b45ea5cb23b7ff3830737020c15
canonical acquisition: UPSTREAM_SOURCE_ACQUIRED id=microsoft-bcquality-1.5
files: 700
Markdown files: 280
AL fixtures/samples: 392
```

## Inspected structure

```text
microsoft/knowledge Markdown: 248
microsoft/skills Markdown: 17
community/knowledge Markdown: 2
global skills Markdown: 6
```

The upstream README defines an entry skill, READ/DO/WRITE meta-skill contracts, Microsoft/community/custom layers, atomic YAML-frontmatter knowledge files, structured agent output and precedence/suppression behavior. It explicitly states that current coverage is technical AL review and that functional Business Central areas remain future scope. Elite preserves that limitation.

## Official validators executed

The exact archive was tested with the commands used by its official GitHub workflows:

```text
python .github/scripts/validate_frontmatter.py --root .
Validator: 0 error(s), 0 warning(s)

pwsh .github/scripts/Test-KnowledgeIndex.ps1 -Root .
Knowledge-index check PASSED: 250 articles, deterministic, full coverage, selection inputs intact.

pwsh tools/Test-ReviewFixtures.ps1 -Root . -PrepareDirectory <isolated-temp>
Review fixture validation PASSED: 32 cases cover 16 leaf domains.
```

The Python validator's official workflow installs unpinned `pyyaml`; the isolated audit environment resolved PyYAML 6.0.3. This is audit-tool evidence, not a production dependency lock.

## Admission

```text
archive/license identity: PASS
frontmatter validation: PASS
deterministic knowledge index: PASS
neutral fixture coverage: PASS
admission: PINNED_CANDIDATE
use: invoke skills/entry.md when BUSINESS_CENTRAL_PLATFORM is selected
non-claim: no functional Finance/SCM/Manufacturing/Service knowledge coverage
```

BCQuality improves agent behavior for AL because the guidance is published and structured by Microsoft. The project must still run BCApps builds/tests in a licensed runtime and must label all custom REVESTEX knowledge in the custom layer or project provenance.
