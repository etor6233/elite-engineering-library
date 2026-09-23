# Scheduled WhatsApp reference infrastructure

Scope: local/fixture proof of appointment reminders on the same CRM, identity,
approval, jobs, provider adapter, delivery fence and status host. No live account
is configured by this composition. Campaign audience/drip attribution is a
separate connected T2805 claim.

Materialize the selected profile and apply migrations through
0083_scheduled_whatsapp.up.sql. Configure the existing WhatsApp host from its
hash-locked profile, then set WHATSAPP_ENABLED=true and
WHATSAPP_SCHEDULE_ENABLED=true. The schedule flag is effective only while the
WhatsApp host is enabled. The service identity additionally needs
notification:dispatch in its configured tenant and organization. No new secret
file or provider credential mechanism is introduced.

The host refuses a missing selected schedule implementation, a nonexact flag,
an unprivileged worker or missing/disabled immutable schema guards. It starts
both the original status worker and the scheduled worker, serializes reports
over the original authenticated mTLS connection and joins both on shutdown.

The operator API under /v1/franchise/notifications/scheduled provides:

- POST /prepare with request_id, appointment_id, recipient, template_name,
  language_code, body_parameters, not_before and expires_at.
- POST /requests with the exact prepared context, followed by
  POST /requests/{delivery_key}/decision by a distinct authorized reviewer.
- GET /requests/{delivery_key} for approval, durable job and provider status.
- POST /requests/{delivery_key}/cancel for an unsent scheduled job.
- POST /requests/{delivery_key}/reconcile for a preserved acceptance receipt.

Use notification:request, notification:approve, notification:read,
notification:cancel and notification:reconcile for their respective operations;
reconciliation also requires read permission. IDs, hashes, tenant, organization,
appointment version/time, contact binding, latest consent, approved template and
provider profile remain bound to the reviewed context. JSON is limited to 32 KiB
with duplicate/case-collision and unknown-field rejection. Body reads have a
5-second deadline and requests a 25-second deadline.

Approval atomically records one durable job, due in at most 30 days. Expiry must
follow the due time by at most 24 hours and cannot exceed the appointment start.
Dispatch uses the original lease-generation worker with three attempts and one
job per poll. Source drift, cancellation or withdrawn consent suppress an
unclaimed provider effect. A rescheduled appointment requires a new prepared
context and distinct review; the old approval is never edited.

The original fence serializes the first effect. A source change after its claim
is treated conservatively as uncertain. No cancellation or consent update can
recall a provider effect already accepted. An uncertain effect requires review;
it is never retried as a fresh message. Exact text and template receipts use the
same original payload builders and recover acceptance without network traffic.
The immutable historical worker result and the current delivery-fence state
are shown separately. Provider acceptance does not mean delivery; signed
observations remain observed_delivered or observed_read.

Expired final leases are quarantined as LEASE_EXHAUSTED with terminal job
evidence, not a false completion. Approval, schedule, cancellation and result
records are immutable. Retention policy and archival for this data must be
included in the deployment's operational policy; these tests do not authorize
deleting evidence. A populated downgrade is refused. Empty down/up is verified.

Local verification: TestScheduledWhatsAppConnected uses PostgreSQL 18.6, the
actual pinned Python adapter and a loopback provider fixture; seven approvals
and jobs, three provider POSTs, six immutable results, twelve concurrent
dispatchers, cancellation, CRM reschedule, consent withdrawal, signed status
and lost-stdout recovery. TestScheduledWhatsAppHost proves actual activation,
both loops, mTLS and shutdown. Boundary regressions and a finite fuzz profile
exercise the new input parser. Original official Meta bytes and source lock are
unchanged; the new scheduling and recovery composition is AUTHORED glue.

Later user setup uses the existing WhatsApp account/phone, approved templates,
provider credentials and identity configuration. The materialized code contains
no account values, access tokens or private keys.
