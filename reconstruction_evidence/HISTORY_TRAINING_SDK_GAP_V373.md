# V373 — pre-implementation SDK selection receipt

Original record and actual resolver receipt, exact JSON bytes and SHA-256. Scope is the component qualification, not the whole pipeline.

## HISTORY_TRAINING_SDK_GAP_V373.json

SHA-256 `faac78bf998fc511cb17dd6de29102f4b372b743253161174179cc1882a91705`

```json
{
  "schema_version": 1,
  "capability_id": "TEST09-OFFLINE-SFT-COMPONENT",
  "requirement_ref": "markdown_system/PROJECT_HISTORY_MODEL_TRAINING_CONTRACT.md",
  "observed_at": "2026-09-10T14:09:40.212113+00:00",
  "max_age_days": 7,
  "authority_refs": [
    "markdown_system/PROJECT_HISTORY_MODEL_TRAINING_CONTRACT.md",
    "PUBLIC_CODE_ARCHITECTURE_ADMISSION_STANDARD.md"
  ],
  "source_searches": [
    {
      "search_id": "official-training-search",
      "query": "Current official training implementations, release/source identity, SDK contracts and native advisories, V367 through V373",
      "official_domains": [
        "github.com",
        "pypi.org",
        "huggingface.co",
        "rustsec.org",
        "docs.rs"
      ],
      "executed_at": "2026-09-10T14:09:40.212113+00:00",
      "evidence_refs": [
        "reconstruction_evidence/HISTORY_TRAINING_SOURCE_RESEARCH_V367.md",
        "reconstruction_evidence/HISTORY_TRAINING_SDK_QUALIFICATION_V373.md"
      ]
    }
  ],
  "candidates": [
    {
      "candidate_id": "qualified-hf-sft-component",
      "origin_kind": "LOCAL_AUTHORED",
      "provenance": "AUTHORED",
      "source_url": null,
      "source_ref": "reconstruction_evidence/HISTORY_TRAINING_SDK_QUALIFICATION_V373.md",
      "claim": "Reusable qualification harness selects exact official SFT dependencies for trusted synthetic Windows CPU reference; pipeline glue not yet implemented.",
      "non_claims": [
        "HF authorship of local harness",
        "complete native security",
        "consumer corpus or model quality",
        "whole pipeline completion"
      ],
      "status": "REUSABLE_PACK",
      "immutable_revision": "V373-sdk-qualified",
      "artifact_sha256": "01ced1c242f852de2e314c446c8ed00dcac58ae393fa243f928fbcd48db1336d",
      "artifact_ref": "reconstruction_evidence/HISTORY_TRAINING_SDK_QUALIFICATION_V373.md",
      "license_expression": "LicenseRef-Workspace-Owner",
      "license_evidence_ref": "AGENTS.md",
      "gates": {
        "G0": "PASS",
        "G1": "PASS",
        "G2": "PASS",
        "G3": "PASS",
        "G4": "PASS",
        "G5": "PASS",
        "G6": "PASS",
        "G7": "PASS",
        "G8": "PASS"
      },
      "governing_authorities": [
        "markdown_system/PROJECT_HISTORY_MODEL_TRAINING_CONTRACT.md",
        "PUBLIC_CODE_ARCHITECTURE_ADMISSION_STANDARD.md"
      ],
      "no_source_search_ids": [
        "official-training-search"
      ]
    }
  ],
  "decision": {
    "state": "USE_REUSABLE_PACK",
    "selected_candidate_id": "qualified-hf-sft-component",
    "reason": "Conditions of narrow SDK qualification demonstrated on the actual reference. Authorize necessary governed integration, retaining provenance and consumer blockers.",
    "canonical_updates": [
      "reconstruction_evidence/HISTORY_TRAINING_SDK_QUALIFICATION_V373.md"
    ],
    "blockers": []
  }
}
```

## HISTORY_TRAINING_SDK_GAP_RECEIPT_V373.json

SHA-256 `0220f93d9f70474b2668863b879324da9a95c4c55cf054257b20da63cbcafdb7`

```json
{
  "authority_count": 2,
  "candidate_count": 1,
  "capability_id": "TEST09-OFFLINE-SFT-COMPONENT",
  "decision_state": "USE_REUSABLE_PACK",
  "gate": "CAPABILITY-GAP-RESOLUTION-GATE",
  "implementation_ready": true,
  "observed_at": "2026-09-10T14:09:40.212113Z",
  "official_domain_count": 5,
  "outcome": "READY",
  "record_sha256": "faac78bf998fc511cb17dd6de29102f4b372b743253161174179cc1882a91705",
  "requirement_ref": "markdown_system/PROJECT_HISTORY_MODEL_TRAINING_CONTRACT.md",
  "schema_version": 1,
  "search_count": 1,
  "selected_candidate_id": "qualified-hf-sft-component"
}
```


## License evidence correction before final qualification

The first record pointed to AGENTS.md for the local license; instructions are not a license grant. The successor points to the actual LICENSE.md, retaining the original record/receipt above. No dependency or control promotion is inferred from that reference correction.

### sdk-gap-license-correction.json

SHA-256 `62dbe2a2295e5b60cef12b25e08734f9091863d8bfd02506625d46ae63d322c1`

```json
{
  "schema_version": 1,
  "capability_id": "TEST09-OFFLINE-SFT-COMPONENT",
  "requirement_ref": "markdown_system/PROJECT_HISTORY_MODEL_TRAINING_CONTRACT.md",
  "observed_at": "2026-09-10T14:09:40.212113+00:00",
  "max_age_days": 7,
  "authority_refs": [
    "markdown_system/PROJECT_HISTORY_MODEL_TRAINING_CONTRACT.md",
    "PUBLIC_CODE_ARCHITECTURE_ADMISSION_STANDARD.md"
  ],
  "source_searches": [
    {
      "search_id": "official-training-search",
      "query": "Current official training implementations, release/source identity, SDK contracts and native advisories, V367 through V373",
      "official_domains": [
        "github.com",
        "pypi.org",
        "huggingface.co",
        "rustsec.org",
        "docs.rs"
      ],
      "executed_at": "2026-09-10T14:09:40.212113+00:00",
      "evidence_refs": [
        "reconstruction_evidence/HISTORY_TRAINING_SOURCE_RESEARCH_V367.md",
        "reconstruction_evidence/HISTORY_TRAINING_SDK_QUALIFICATION_V373.md"
      ]
    }
  ],
  "candidates": [
    {
      "candidate_id": "qualified-hf-sft-component",
      "origin_kind": "LOCAL_AUTHORED",
      "provenance": "AUTHORED",
      "source_url": null,
      "source_ref": "reconstruction_evidence/HISTORY_TRAINING_SDK_QUALIFICATION_V373.md",
      "claim": "Reusable qualification harness selects exact official SFT dependencies for trusted synthetic Windows CPU reference; pipeline glue not yet implemented.",
      "non_claims": [
        "HF authorship of local harness",
        "complete native security",
        "consumer corpus or model quality",
        "whole pipeline completion"
      ],
      "status": "REUSABLE_PACK",
      "immutable_revision": "V373-sdk-qualified",
      "artifact_sha256": "01ced1c242f852de2e314c446c8ed00dcac58ae393fa243f928fbcd48db1336d",
      "artifact_ref": "reconstruction_evidence/HISTORY_TRAINING_SDK_QUALIFICATION_V373.md",
      "license_expression": "LicenseRef-Workspace-Owner",
      "license_evidence_ref": "LICENSE.md",
      "gates": {
        "G0": "PASS",
        "G1": "PASS",
        "G2": "PASS",
        "G3": "PASS",
        "G4": "PASS",
        "G5": "PASS",
        "G6": "PASS",
        "G7": "PASS",
        "G8": "PASS"
      },
      "governing_authorities": [
        "markdown_system/PROJECT_HISTORY_MODEL_TRAINING_CONTRACT.md",
        "PUBLIC_CODE_ARCHITECTURE_ADMISSION_STANDARD.md"
      ],
      "no_source_search_ids": [
        "official-training-search"
      ]
    }
  ],
  "decision": {
    "state": "USE_REUSABLE_PACK",
    "selected_candidate_id": "qualified-hf-sft-component",
    "reason": "Conditions of narrow SDK qualification demonstrated on the actual reference. Authorize necessary governed integration, retaining provenance and consumer blockers.",
    "canonical_updates": [
      "reconstruction_evidence/HISTORY_TRAINING_SDK_QUALIFICATION_V373.md"
    ],
    "blockers": []
  }
}
```

### sdk-gap-license-correction-receipt.json

SHA-256 `b60288803a00d5d410a2affcee73e925575dbde9ef78a929c76f8c619609a577`

```json
{
  "authority_count": 2,
  "candidate_count": 1,
  "capability_id": "TEST09-OFFLINE-SFT-COMPONENT",
  "decision_state": "USE_REUSABLE_PACK",
  "gate": "CAPABILITY-GAP-RESOLUTION-GATE",
  "implementation_ready": true,
  "observed_at": "2026-09-10T14:09:40.212113Z",
  "official_domain_count": 5,
  "outcome": "READY",
  "record_sha256": "62dbe2a2295e5b60cef12b25e08734f9091863d8bfd02506625d46ae63d322c1",
  "requirement_ref": "markdown_system/PROJECT_HISTORY_MODEL_TRAINING_CONTRACT.md",
  "schema_version": 1,
  "search_count": 1,
  "selected_candidate_id": "qualified-hf-sft-component"
}
```

