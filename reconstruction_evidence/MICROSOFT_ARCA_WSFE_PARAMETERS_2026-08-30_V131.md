# Microsoft-generated ARCA WSFE parameter bridge — reconstruction evidence V131

## Scope and authority

V131 extends `MICROSOFT-ARCA-WSFE-SOAP-ADAPTER` from 0.1.0 to 0.2.0 without creating a second fiscal adapter. The governing public ARCA WSFEv1 page and developer manual 4.6 define the generated operations used here: voucher types, concepts, document types, VAT rates, other taxes, electronic points of sale and recipient VAT conditions. Microsoft `dotnet-svcutil` remains the sole proxy generator. No parameter value, WSDL, generated proxy, provider message, credential, certificate, token or signature is embedded in the pack.

- ARCA WSFEv1: <https://arca.gob.ar/ws/documentacion/ws-factura-electronica.asp>
- ARCA developer manual 4.6: <https://arca.gob.ar/ws/documentacion/manuales/manual-desarrollador-ARCA-COMPG-v4-0.pdf>
- Microsoft dotnet-svcutil: <https://learn.microsoft.com/dotnet/core/additional-tools/dotnet-svcutil-guide>

The fourteen files are local `AUTHORED` integration and tests governed by those contracts. They are not presented as copied ARCA or Microsoft source.

## Exact artifacts

- Adapter pack: `MICROSOFT-ARCA-WSFE-SOAP-ADAPTER 0.2.0`, 64,110 bytes, SHA-256 `74dd0b1e74972fee4005fbf2fea6d84f79b769581358c1637449d821daff6fe4`.
- Compatible UDS worker: `ARCA-WSFE-UDS-WORKER 0.1.1`, SHA-256 `3fc79c481edba89c5583a7f63c77917c70a3e504d6fc4873e9e81c188f6f0d46`.
- Backend plan: SHA-256 `7c4f069336873cb1c6cafabe6fe387d50b0d897d58debb2cad82e84f18896e59`.
- Clean backend record: SHA-256 `0e8042339ea3329415dac8e6b4229958133259ad1de442bf207198b0b7abae5f`.

## Executed gates

1. A staging tree built both .NET test projects with SDK 10.0.400, locked packages, an empty NuGet source set and `TreatWarningsAsErrors`: zero warnings and zero errors.
2. The integrated runner returned `ARCA_WSFE_ADAPTER_GATE_PASS environment=Homologation invoice_cases=12 parameter_cases=15 production_admitted=false`.
3. The canonical Markdown materialized 14/14 files into a new directory. Running the materialized runner repeated both locked offline builds and all 27 contracts with the same marker.
4. The fifteen parameter cases exercise the seven generated operation families, deterministic hashing, exact safe provider-code propagation, CUIT/class/date/duplicate rejection and absence of credentials or provider messages in returned DTOs.
5. The corrected .NET 10 command queried `https://api.nuget.org/v3/index.json` with `--vulnerable --include-transitive --format json`; the current graph returned no vulnerable package entries.
6. The backend composed from Markdown into a new directory: `Materialized 217 files from 24 pack(s)`. Go 1.26.7 completed the full suite twice, `go vet ./...`, and builds of `cmd/api` plus `cmd/electromobility-api`.
7. `LIB-FAIL-1328` preserves the rejected positional .NET 10 package-list syntax and its `--project` correction. `LIB-FAIL-1329` preserves the composite-wrapper timeout and the split gate that produced complete evidence. `LIB-FAIL-1330` preserves the Windows glob mistake and its portable `rg -g` correction. `LIB-FAIL-1331` prevents a partial executable audit from being mistaken for its final result.
8. The governing close returned `VERIFY_LIBRARY_PASS packs=83 materialized_files=862 markdown_files=470`, including backend 217 and web 80. A persistent PTY run then returned exit 0 and `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit packs=83 upstream_sources=121`; account/network/target-dependent gates remained explicitly skipped rather than simulated.

## Admission boundary

This closes reconstructible parameter acquisition and safe normalization, not production fiscal admission. The pack does not claim that a parameter observed once is timeless, approved business policy or safe to persist automatically. Real certificate custody, service association, live homologation, refresh/persistence ownership, current target responses and independent tax/accounting approval remain mandatory. No project may claim production readiness without the real CDN/WAF, IdP, providers, PostgreSQL recovery, load, offensive security, deploy/rollback and business-acceptance receipts required by the library.
