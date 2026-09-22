# V403 exact source publication

[Start](../markdown_system/START_V403_LOCAL.md) · [Plan](../markdown_system/UNIFIED_REFERENCE_PACK_PLAN_V403_R4.md)

This index fixes source and publication identities. It is not a signature or executed qualification.

```json
{
  "schema": "elite-v403-publication-index/v1",
  "source_revision": "NEXT-UNIFIED-R4-RECOVERED-CANDIDATE",
  "source_manifest_sha256": "91db99c7b00d30b06c34669628b9050d6bc55983d52f1c6a8d06da11742940ef",
  "source_files": 2147,
  "carrier_files": 2113,
  "base64_objects": 90,
  "files": {
    "implementation_packs/UNIFIED_REFERENCE_V403_R4.md": "a48ba1f0755cdb2fd505cf04b81444d76b8fd4415e2980888d59078f71286280",
    "markdown_system/UNIFIED_REFERENCE_PACK_PLAN_V403_R4.md": "2095f2689101cff83229387a54c704e28cbb011eb3024bc3401b3b9b46ba2574",
    "markdown_system/START_V403_LOCAL.md": "c05a529ccbba6a7c97fa07f970a7d4be67578d949e1f034c5ea9ba811298e550"
  },
  "generator_sha256": "f3649a2eeedfdfd9b84522a70b3b916d3035b9ac1019f2661664d51fc498cbe4",
  "production_authorized": false
}
```

Append an immutable qualification receipt and its SHA after running the generated tool. Do not put that receipt into the source manifest or carrier that it verifies: the dependency is source → carrier → qualification receipt. Preserve failures separately. The evidence document does not hash itself; Git/release manifest pins its final bytes. V402 historical VERIFY/ZIPs remain separate.
