# Unified V403 source delivery plan

[Start](START_V403_LOCAL.md) · [Source pack](../implementation_packs/UNIFIED_REFERENCE_V403_R4.md)

One exact carrier; does not mutate original public/admin packs or confer their visual approval.

```json
{
  "compositionVersion": "1.0",
  "secretVariables": {},
  "packs": [
    {
      "path": "implementation_packs/UNIFIED_REFERENCE_V403_R4.md",
      "packId": "UNIFIED-REFERENCE-V403-R4",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    }
  ]
}
```

Compose with the existing Markdown compositor or materialize the selected pack directly. Then follow START_V403_LOCAL. The successor verifier checks this explicit plan independently of the historical profile-count gate.
