# ARCA WSAA Official Client Samples — Reconstruction Evidence V122

## Scope and provenance

V122 adds the official ARCA/AFIP WSAA client samples for PowerShell and C# as protocol authorities with their redistribution agreement preserved. It does not present legacy sample code as a secure production adapter.

The exact official archives are:

- PowerShell: `https://arca.gob.ar/ws/WSAA/ejemplos/dev-wsaa-cliente-powershell.zip`, 48,293 bytes, SHA-256 `6d2a36105dafef967cc3e92a001e9f0f8746135406bbd99e0e3195ca7f162514`;
- C#: `https://arca.gob.ar/ws/WSAA/ejemplos/dev-wsaa-cliente-dotnet-cs.zip`, 64,195 bytes, SHA-256 `f46c8ba58968abf91a96f617d23de256b3c30a29f8f55c3aa9072c88865fbb73`.

Selected original hashes are PowerShell README `2303fe9b36fb2fa5ffc1fddcef1422fea1e76d43730049f0b25b4627c3b96715`, PowerShell source `edd4b7d840cf873d1913177fbb8e5beaac266f53770cf5d893a2e554f5d2db6f`, C# README `c985c5e1e195fe193e1b47605eeb93b24536d49c463e235618fe6f1aa3d3f436`, C# source `6813ae36b017624eca84acca09f7d2dccbb0eacb4be935efbcf6264e2bffea0b` and `app.config` `7f84a6c529f2a134e2d9d350f735008548376b312018b3c0192d88bcf0a7965e`.

## Materialized result

- `ARCA-WSAA-OFFICIAL-CLIENT-SAMPLES 0.1.0`, pack SHA-256 `28b09321cbb839fa3d24c229cdb377d8a4005c00d0cc89c08013ab612bc8677d`, materializes seven files.
- Five official text files are declared `ADAPTED` solely because the Markdown contract normalizes CRLF to LF. Their LF hashes are separately fixed; no vendor authorship is claimed for the normalization.
- Two local `AUTHORED` scripts acquire and verify the exact original archives/files, reject path escape, preserve the redistribution terms, support an offline cache and emit a receipt with `production_admitted=false`.
- The PowerShell sample still uses `New-WebServiceProxy`, external OpenSSL and local ticket files. The C# sample targets legacy .NET Framework and accepts a password argument. Those conditions are explicit blockers, not silently repaired or ignored.
- Current inventory after close: 78 packs / 798 materializable files; provenance 656 `AUTHORED`, 37 `ADAPTED`, 105 `VERBATIM`. Backend remains 19/165 and web 6/80.

## Exact executable proof

- Clean materialization reconstructed 7/7 files and verified every block SHA-256.
- The materialized sample regression verified all five LF hashes, required redistribution text, PowerShell parser syntax and known legacy conditions.
- Offline acquisition verified both archive byte lengths/SHA-256 values and all five original file hashes, then produced the non-secret source receipt.
- The global library verifier is recorded after removal of staging directories and canonical inventory alignment.

## Failure memory

The first gate disproved the byte-exact claim because Markdown normalizes line endings. V122 corrected provenance to `ADAPTED`, separated materialized and original hashes, and retained exact binary acquisition. A rejected combined patch and one wrong parameter invocation are also recorded in the persistent failure ledger.

## Production boundary

This pack provides official protocol samples and acquisition evidence only. A real Argentine fiscal integration still requires the current WSAA/WSFE specifications, modern certificate/key custody, ticket caching without secret leakage, homologation credentials, CUIT/points of sale/document types/tax rules, idempotency, reconciliation, observability, negative contract tests and successful target homologation. No project is production-ready until its CDN/WAF, IdP, providers, PostgreSQL/recovery, load, offensive security, deployment/rollback and business acceptance are demonstrated.
