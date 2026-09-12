# Reconstruction Evidence — TS Enterprise Web V4

```yaml
evidence_id: "TS-EW-20260824-V4"
pack_sha256: "cd8c57f7aed606edc635395144101f4e5850d512af7c2e3985771293b680259d"
materialized_file_count: 64
ordered_file_hash_aggregate_sha256: "886ea80e0219e8bb611027f3ace707a233c388172bcc83a7f8dd09e080eb7a20"
clean_target: "rebuild_verification/ts-enterprise-web-v4"
node: "v24.14.1"
pnpm: "11.19.0"
result: "PASS"
```

Materialization verified every declared file hash. Frozen install produced 89 packages. Configuration, baseline readiness, TypeScript, 9 test files/22 tests, Next.js production build and production dependency audit passed. The added CI workflow uses least-privilege permissions and pins checkout, Node setup and pnpm setup actions by full commit SHA.

