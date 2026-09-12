# Microsoft-generated ARCA WSAA and WSFEv1 Clients — Reconstruction Evidence V124

## Result

V124 extends the existing generator instead of creating a duplicated WSAA toolchain. `MICROSOFT-ARCA-WSFE-GENERATED-CLIENT 0.2.0`, pack SHA-256 `791e3d488534d595c076602c8e5675b6260191cc1ff1540a147acd37270fe645`, still materializes six local `AUTHORED` files and now generates two namespace-isolated clients in one project.

No WSDL or generated proxy is embedded. Public ARCA contracts are acquired by exact URL, byte length and SHA-256; Microsoft `dotnet-svcutil` 8.0.0 generates both clients under .NET SDK 10.0.400 with WCF 10.0.652802 and patched Cryptography.Xml/Pkcs 10.0.11.

## Exact generated identities

| Environment | WSAA generated SHA-256 | WSFEv1 generated SHA-256 |
|---|---|---|
| Homologation | `910a9e267718e663424972d2dd10d1f757e14f617e8a0d1afab8580bd6c450ea` | `e4608d9e113f293945a0fadf7a8c5f018ae0a2e19568ea450772be8bd1ec9d61` |
| Production metadata | `71c36f4420b645365e06b35e76d4f856ac0094aff87aa80c615e26195bf3fb02` | `85c3fa927eec659e952b9a8982e2a21368fee29808c3e32fa1fd5704749fb021` |

After replacing only their official endpoint URLs, homologation and production WSDLs are byte-equal per service, and their generated clients are likewise byte-equal per service. Environment-specific bytes remain fixed independently rather than normalized in the pack.

## Executable proof

- Materialization and structural regression: 6/6 PASS.
- Homologation: WSAA + WSFEv1 generated, exact hashes, restore/build with zero warnings/errors and current NuGet vulnerable scan clean.
- Production metadata: explicit `-AllowProductionMetadata`, both exact hashes and the same compile/security gates PASS.
- The runner restores PATH/telemetry/roll-forward environment in `finally`, writes a two-contract receipt and always reports `production_admitted=false`.
- Global inventory remains 79 packs / 804 materializable files; provenance remains 662 `AUTHORED`, 37 `ADAPTED`, 105 `VERBATIM`.

## Boundary

Generated transport types do not implement certificate/key custody, CMS request policy, ticket caching/fencing, fiscal document rules, durable issuance, CAE reconciliation, authorization or homologation. Those must be a separate application adapter with project credentials and evidence. A concrete project remains unready for production until CDN/WAF, IdP, providers, PostgreSQL/recovery, load, offensive security, deployment/rollback and business acceptance are demonstrated.
