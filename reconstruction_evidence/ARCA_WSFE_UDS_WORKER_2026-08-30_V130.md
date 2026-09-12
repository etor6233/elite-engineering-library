# ARCA WSFE Unix-socket worker — reconstruction evidence V130

- Date: 2026-08-30
- Agent/reviewer: Codex maintenance run
- Pack: `ARCA-WSFE-UDS-WORKER 0.1.0`
- Pack SHA-256: `0535b2c0130fe96852aef2e0037028f98766526b610e49465c23a4329bdc59aa`
- Updated application pack: `GO-ELECTROMOBILITY-APPLICATION 1.4.0`
- Application pack SHA-256: `e63df1475f7e2e73bb0cd5597ec7e872401dfa07ee4f407386810266668e68c4`
- Clean backend composition: 24 packs / 213 files
- Plan SHA-256: `a353d3da3a5601d9de79a12f0e8cc68d20a18149dd4b04d779659eace36fd417`
- Materialization record SHA-256: `4a242bb9f298a47d1072d9f148e0d58bcf8217a98f31cec5e4363ab07d9efa63`
- Toolchains: Go 1.26.7; .NET SDK 10.0.400/runtime 10.0.11; WCF 10.0.652802; pgx 5.10.0; generated homologation clients from the V124/V129 pinned lane
- Host exercised: Windows with native Unix-domain sockets; no public TCP listener

## Authority and provenance

The implementation is local `AUTHORED` composition, not copied or attributed to Microsoft, Google, Go or ARCA. It applies the official Go `net.Dialer.DialContext` Unix transport, Microsoft Kestrel `ListenUnixSocket` and security requirements, Microsoft-generated WCF interfaces and the current ARCA WSFE owner/adapter admitted in V128/V129. The thirteen new worker/IPC files and two shared-ID files contain no WSDL, generated proxy, upstream source, PFX, password, private key, Token, Sign, account credential or provider message.

Primary authorities:

- `https://pkg.go.dev/net`
- `https://learn.microsoft.com/en-us/dotnet/api/microsoft.aspnetcore.server.kestrel.core.kestrelserveroptions.listenunixsocket?view=aspnetcore-10.0`
- `https://learn.microsoft.com/en-us/aspnet/core/fundamentals/servers/kestrel/security-considerations?view=aspnetcore-10.0`
- `https://github.com/dotnet/wcf`
- `https://www.arca.gob.ar/ws/documentacion/manuales/manual-desarrollador-ARCA-COMPG.pdf`

## Exact new-file hashes

| SHA-256 | File |
|---|---|
| `df53732678241d8b3063bffb0524ec0e6e67c7014f41905bedc00aa768c0bf31` | `cmd/arca-fiscal-worker/main.go` |
| `ce7e559424346f277286bee1657034afc5f684e713865a4ce5f8ad3d3ab7661b` | `internal/fiscal/wsfeipc/provider.go` |
| `c2bd848c93b8faef51f4575d12c32c9286ccbd8c69f00e83eb397dfeaa84ecf3` | `internal/fiscal/wsfeipc/provider_test.go` |
| `7e315b51ef64b3587b5943dd7ff62b0615817dbc69775e3fd64d214784b849ec` | `arca/fiscal/worker/README.md` |
| `68c57bf0bb38a34e90e5e2d8272d17e090660b0e0519ac6ec4b153c12e31942e` | worker `.csproj` |
| `cc28c4c6193a2f96790da7d6b178ffbb6c5788377dfd2a6ed22e57b3f187d82d` | worker `packages.lock.json` |
| `5e43347a37fbd0da3d6044c6c0a4937c7b2e50cb7628332a43ed2a21657baaa8` | worker `Program.cs` |
| `15e088902093d79de348817d8ff27d7099f63cc23b954be4c2426a34d341137e` | worker `SoapTransports.cs` |
| `b8adae81c0867cfc591b360c531b8fc2faf34c02cd296459db5ae056ab80c9bf` | worker `WorkerHost.cs` |
| `f1f2e0664386aec472e683629d97f530e23c1cf3391c1afea706ceae69d0cb1a` | test `.csproj` |
| `55e434a478b434a6a2180f8ab587263657f990b72441b8eff1c38004b29fa431` | test `packages.lock.json` |
| `e6c22b780793ad8161ce31aec8e5e85dccfd99e23e7dc0b023a5618897dbd508` | test `Program.cs` |
| `13c1ea03371f89cff7b09ac0e3ce47294882cbb3544aa5ccfd76a1fc32326e82` | `tools/test-arca-wsfe-uds-worker.ps1` |
| `a5894d48b94c44149b6bd5966d2a9423a23021b52897b2a2c660d1b3ed195120` | shared `randomid/generator.go` |
| `e5675ea7fd971cf35a48f9e31891e5d7d4eca555fb1d3fd7235a5218fca52098` | shared `randomid/generator_test.go` |

## Executed gates

1. `ARCA-WSFE-UDS-WORKER` materialized 13/13 files with exact manifest/block/hash agreement. Application 1.4.0 materialized 4/4 and removed the composition root's private UUID implementation.
2. The backend composed from Markdown into a new directory: `Materialized 213 files from 24 pack(s)`; the record contains 213 unique paths and the two new identities.
3. Go full suite ran twice from clean/canonical trees. Both executed the real `net.Listen("unix", socket)` contract and every package; `go vet ./...` and builds for `cmd/electromobility-api` plus `cmd/arca-fiscal-worker` passed.
4. The Go test proved exact tax fields crossed the UDS, no tenant/internal identity crossed, last/consult/authorize mapped, relative sockets failed, unknown/invalid responses failed, hashes were exact and response size was bounded.
5. .NET locked restore, warnings-as-errors build and eight executable Kestrel contracts passed. The test started the actual server over UDS, connected with `AddressFamily.Unix`, observed only a Unix address, exercised liveness/last/authorize, rejected unknown and oversized requests, and proved provider messages/credentials absent from the failure response.
6. A new source tree restored all .NET projects with `--locked-mode` and a NuGet config containing `<clear />`, then rebuilt and repeated all eight tests. Current NuGet.org vulnerable-package inspection reported no vulnerable packages in the complete worker test graph.
7. The runner produced `ARCA_WSFE_UDS_WORKER_PASS go=1.26.7 dotnet=10.0.400 uds_tests=8 production_admitted=false`.
8. The governing library gates closed after all documentation updates: `VERIFY_LIBRARY_PASS packs=83 materialized_files=858 markdown_files=469`, backend profile 213, and `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit packs=83 upstream_sources=121`. The audit re-executed its offline/browser/document/security/provider gates; live account/cost/target gates remained explicitly skipped and conditioned.

## Failure memory returned to the library

`LIB-FAIL-1322` through `LIB-FAIL-1327` preserve the incorrect literal wildcard, C# top-level ordering, Kestrel UDS address assumption, out-of-module `go fmt` call, PowerShell backtick recurrence and rejected optional temp cleanup. Every implementation failure was discarded, fixed at the source and repeated through its gate; the two auxiliary files retained outside the workspace cannot enter verification or distribution.

## Retained conditions

This evidence closes local Go↔.NET IPC, native UDS behavior, strict serialization, transport bounds, error redaction, shared identifiers and reconstructed builds. It does not claim live ARCA homologation, certificate authorization/association, current project parameter tables, tax/legal approval, production availability or a deployed socket directory/service identity. Production remains blocked until those project-specific gates, PostgreSQL recovery, load, offensive security, rollout/rollback and business acceptance are demonstrated on the selected target.
