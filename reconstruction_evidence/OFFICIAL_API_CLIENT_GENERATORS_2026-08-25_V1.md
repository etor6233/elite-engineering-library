# Official API Client Generators — 2026-08-25 V1

> Evidencia histórica. Kiota 1.34.1 fue sustituido para uso vigente por Kiota 1.35.0 y el pack autocontenido Go demostrado en `MICROSOFT_KIOTA_OPENAPI_CLIENT_GATE_2026-09-05_V247.md`. El bloqueo .NET de esta revisión histórica ya no aplica a esa lane Kiota; sí continúa para WCF.

## Purpose

This evidence fixes two Microsoft/.NET source releases that can generate typed clients from official API contracts. Generated clients remain project artifacts; the tools do not create business invariants, provider authorization, idempotency, reconciliation or production evidence.

## Microsoft Kiota 1.34.1

- repository: `microsoft/kiota`;
- release: `v1.34.1`;
- commit: `9f9cfb3b1cb9b5311a214ea6ce0f69943c523005`;
- archive: 3,509,773 bytes;
- archive SHA-256: `2b9a25704a274d4f5eac64c826b3cd2a891415f28f4d1912ace32d03f13669b2`;
- license: MIT, SHA-256 `c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383`;
- `global.json` SHA-256: `697c274885ee41438dc9306a384574e64365f9a22aa0f0b7819d52607410ce38`;
- acquired tree: 1,044 files, 8,478,194 bytes, 239 files under test paths.

The official README states that Kiota generates strongly typed API clients from OpenAPI and lists stable generation lanes including Go, Java, Python and C#. TypeScript is not selected as a backend foundation by this admission.

## .NET WCF 10.0.0

- repository: `dotnet/wcf`;
- release: `v10.0.0-rtm`;
- commit: `e9d8c1c2a051618689bc22ab263f6ff0f2493d64`;
- archive: 11,391,314 bytes;
- archive SHA-256: `cd3fb37b61278f894edeeba981ee5ae0535ce58d1d52f580b553a9b0de76102d`;
- license: MIT, SHA-256 `cfc21f5e8bd655ae997eec916138b707b1d290b83272c02a95c9f821b8c87310`;
- third-party notices SHA-256: `9afe4c3e0254e03e96859c999c685deb27519a12baceaf44521424da5fa2e827`;
- `global.json` SHA-256: `fe5dcc20c591152ae2254bacb2a40b35d31269aa4e8aa24496d00e84268228b2`;
- acquired tree: 4,060 files, 57,229,578 bytes, 832 files under test paths.

The official source contains the WCF client libraries and `dotnet-svcutil` solution/tooling used to generate a .NET client from trusted WSDL/XSD metadata. This can support a project lane for an official SOAP contract such as a selected fiscal service, but it is not an ARCA adapter and does not prove credentials, certificates, homologation or legal correctness.

## Actual acquisition

```text
UPSTREAM_ACQUISITION_TEST_PASS
UPSTREAM_SOURCE_ACQUIRED id=microsoft-kiota-1.34.1
UPSTREAM_SOURCE_ACQUIRED id=dotnet-wcf-10.0.0
```

The acquisition runner verified exact archive length/digest, license and `global.json`; WCF notices were also verified.

## Retained blocker

`dotnet --version` failed because this audit host has no .NET SDK. This is `LIB-FAIL-029` and remains `BLOCKED_EXTERNAL`. No build, generator invocation or generated-client contract test is claimed. `VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Preflight` now reports .NET 10+ explicitly for these lanes.

The final project preflight remained correctly blocked and reported:

```text
TOOLCHAIN id=pwsh-7 available=true
TOOLCHAIN id=python-3.12+ available=true
TOOLCHAIN id=go-1.26.7 available=false
TOOLCHAIN id=dotnet-10+ available=false
TOOLCHAIN id=node-24+ available=true
TOOLCHAIN id=pnpm-11.19.0 available=true
TOOLCHAIN id=docker available=false
TOOLCHAIN id=psql available=false
EXECUTABLE_PREFLIGHT_BLOCKED missing=4
PROJECT_ACCESS_REMAINS_REQUIRED provider sandboxes, identities/scopes, secrets manager references, document corpus/ground truth, deployment and recovery targets
```

Go, Docker, .NET and PostgreSQL CLI are indexed as external blockers. The blocked preflight is the intended fail-closed result and is not relabeled as a library failure.

## Admission

Both sources are `PINNED_CANDIDATE`. A project may select them only after fixing:

- an immutable, trusted OpenAPI or WSDL/XSD contract;
- compatible .NET SDK/tool artifact and dependency graph;
- generated output language/runtime and exact generator configuration;
- diff review, compile and official sandbox contract tests;
- authentication, secrets, deadlines, retry budget, idempotency, webhooks and reconciliation;
- regeneration/rollback policy when the provider contract changes.

This path avoids hand-writing transport clients while preserving the truth that provider-specific operational behavior still requires project evidence.
