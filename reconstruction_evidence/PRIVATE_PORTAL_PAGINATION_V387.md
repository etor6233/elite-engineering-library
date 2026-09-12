# V387 — private portal cursor traversal

Library maintenance / TEST02 owner integration. Entry checkpoint235; V386 security
research remains ACCESS_BLOCKED under the user's explicit instruction. This work
does not resume vulnerability reproduction or change its admission status.

## Fixed behavior

Customer orders/service cases and factory units previously ignored next_cursor,
leaving records beyond the first25unreachable. Two rendering regressions demonstrated
the defect before product changes. Existing backend keyset pagination now drives
ordinary GET navigation. Customer cursors stay independent; reload preserves them,
first-page navigation clears only the chosen list, final pages omit next links,
and empty pages state their condition. The server retains session scope and25row
limit; user query parameters cannot override these. Opaque values are encoded,
bounded and checked for ambiguity; existing fixed read recovery handles rejection.

No API, migration, currency, permissions, warranty policy, provider or runtime
dependency changed. New helper/test code is AUTHORED, not attributed to Next/Go.
Next16.3.4 page searchParams contract was checked against its official documentation:
https://nextjs.org/docs/app/api-reference/file-conventions/page . The local backend
is the authority for after/next_cursor, not a newly invented pagination service.

## Validation

272web tests passed; one inherited opt-in test remains skipped in that command.
The actual connected gate ran separately: four browser projects,27records per
list/81per project, no omissions/duplicates, independent cursor retention, reload,
return-to-first and current read recovery. Five inherited direct HTTP negatives,
zero browser/API writes and identical before/after contents of six durable tables.
Mobile final-page screenshot reviewed; no horizontal overflow. Owned PostgreSQL
and Next processes stopped. Typecheck/build passed.6canonical packs reconstructed;
full69pack/805file reference equals the tested source with explicit generated
next-env imports and JSON-equivalent pnpm manifest normalization. Parent68/797.
Go HTTP executable matches V375 byte-for-byte.8589dependency files retain their
authenticated hashes; this identity check does not close native security admission.

## Failures retained and repaired

FAIL767 is the pagination defect. FAIL768 is the interrupted long-path dependency
copy; partial files remain isolated,803canonical inputs plus offline frozen install
and8589hash checks repaired the environment. FAIL769 is invalid fixture data:
26factory units omitted identifiers under UNIQUE NULLS NOT DISTINCT. The fixture
now supplies unique synthetic VIN/battery values; constraints remain intact and
the full connected suite passed in a fresh database. Original failed receipts stay.

## Limits and continuation

TEST02 is not closed by this narrow correction; FAIL385 still requires remaining
functional owner integrations. Admin pagination and the separate journey arrays
are not covered by these three lists. No cross-page transaction snapshot or new
business rules are claimed. TEST03/07 remain blocked as recorded in V386.45/48.
Full Preflight against the new checkpoint is pending below; no final release.

## Evidence

Raw receipts resolve through elite-v387-current.txt under local Temp.

| Receipt | SHA-256 |
|---|---|
| baseline.json | a5ac5770e972f237d0975f9dc5970474734fe97d668dc95801352909c9f1d15e |
| artifact-parity.json | 7563f34b2a058f7ce4c0f799cb29ad7b97d517326579e335c400497c8723cb43 |
| web-pagination-qualified-results.json | 1c4e71a01633f9df8acf5f0911f22cf1abd5b5d2dff55fc6e3f8ac2f8faa0380 |
| canonical-pending-changes.json | 794789757af3bc2ec92a82250700559dff1c33ae6786fdfe2e9674e62d8b0594 |
| final-parity.json | f393c5cb7bc925dda7f241f68c54ed9384d49cdbe7385dbe236154cb037208a2 |
| connected result.json | 88dde720c117a4eab353212aef4f6027567235e62877b69b8fb9578b9ba852e6 |
| connected admin-reads.log | b7c911df8783e551a70dc2e8785d55a1fffc04cf1d0f89a652ba855202ff5b7c |

## Combined final validation — V388 / checkpoint239

Preflight236 rejected only the stale canonical inventory (FAIL771); it is not reported as passed. After its correction and the adjacent admin paging extension, Preflight238 passed all164executed steps and56compositions;276web tests and4connected browsers include these original three lists. Source V387 and failed236 receipts remain preserved; V388 is the final current composition. See ADMIN_PORTAL_PAGINATION_V388.md.
