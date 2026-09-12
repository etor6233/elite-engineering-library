# Official upstream acquisition approval hardening — evidence V1

Audit date: 2026-08-26.

## Defect

`apply_source_profile.ps1` correctly required a hash-linked user approval before acquisition. However, its lower-level `acquire_upstream_sources.ps1` accepted a direct `-SourceId` acquisition without `ProfilePath` or `ApprovalPath`. A caller could therefore bypass the user-input gate while still downloading a correctly locked archive. This violated the library's pre-start rule even though source integrity remained enforced.

## Correction

`OFFICIAL-UPSTREAM-ACQUISITION-CORE` 0.4.7 now requires both paths for every non-`ValidateOnly` execution. The acquirer independently revalidates:

- profile and approval schemas;
- profile, lock and approval SHA linkage;
- approver identity and parseable timestamp;
- exact selected source set equals the approved profile set;
- one closed, non-empty answer for every declared required user input and no extra answers;
- exact production-blocker acknowledgement.

`apply_source_profile.ps1` passes the same already-validated profile and approval to the acquirer. A direct call without both now fails before destination creation or network access.

## Regression

The exact reconstructed core produced:

```text
UPSTREAM_ACQUISITION_TEST_PASS negatives=4
UPSTREAM_LOCK_VALID sources=56 selected=20
SOURCE_PROFILE_VALID profile=document-intelligence-leaders selected=20
UPSTREAM_LOCK_VALID sources=56 selected=24
SOURCE_PROFILE_VALID profile=commerce-communications-leaders selected=24
UPSTREAM_LOCK_VALID sources=56 selected=10
SOURCE_PROFILE_VALID profile=enterprise-platform-leaders selected=10
SOURCE_PROFILE_TEST_PASS valid=3 negatives=5
```

The added negative invokes the direct acquirer with a real locked source but without approval and requires `ProfilePath and ApprovalPath are required for acquisition`. Network is never contacted in that case. This evidence proves fail-closed orchestration, not that any provider account, source use or production deployment has been approved.
