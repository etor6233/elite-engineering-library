# Materialized initial-handover profile

The supported algorithm is `single-unit-full-observed-payment`, algorithm
revision 1. It requires one allocated serialized unit, an active matching
reservation, full observed order payment with no refunds/disputes/hold, and the
existing customer/checklist acceptance. Its final check is read-only eligibility;
no shipment, fiscal operation or irreversible release is introduced.

Generate a profile in an absent destination using the supplied standard-library
Python materializer. Tenant/organization IDs are configuration, never secrets.
The profile, exact-byte SHA-256 activation lock, documentary decision and receipt
are written together. Omission of `--activate` leaves activation disabled.

```powershell
python tools/materialize_handover_profile.py --output handover-policy --profile-id franchise-initial --revision 1 --tenant-id YOUR_TENANT_UUID --organization-id YOUR_ORGANIZATION_ID --payment-provider stripe --payment-account-ref YOUR_PROVIDER_ACCOUNT_ID --payment-connection-id YOUR_CONNECTION_ID --expected-mode sandbox --activate
```

The host loads `profile.json` with `LoadHandoverProfileFile` and an explicit
`HandoverProfileActivation`. Its tenant, organization, provider, account,
connection and expected provider mode
must also match the server/payment configuration, not request input. A future
`live` configuration uses the same algorithm and loader without Go edits. This
does not make live credentials, deployment or business acceptance PROVEN.

The loader rejects missing/altered documents, duplicate or unknown JSON fields,
unsupported schema/algorithm revisions, unbound profile revisions, scope/mode
mismatches and every option outside the supported algorithm. `policy:true` has
no meaning. A private admitted binding prevents a reconstructed JSON contract
or subsequent field edits from bypassing the loader. Both domain and PostgreSQL
owners enforce the admitted tenant/organization scope. PostgreSQL also requires
the profile's active provider connection, account and exact checkout binding.

Change profile configuration by issuing a new document revision and activation
hash. A preparation records that exact profile ID, revision and SHA-256; another
profile cannot silently reinterpret its release result. Its receipt remains
recoverable as history. Existing exact `LOCAL_FIXTURES` policy and fixtures are
retained for regression, separately from the materialized profile.

Authority/decision references are documentary selection records. They do not
assert signature verification, credentials, live production readiness, legal or
regulatory approval. Requirements beyond this supported algorithm remain a
capability decision rather than an arbitrary override flag.
