# Connected campaign reference infrastructure

Selected scope: explicit permitted CRM lead audiences, exact reviewed WhatsApp
template steps, current consent/source/stop/conversion suppression, durable
provider evidence and observed quotation/order conversion. The immutable
campaign binds the existing CRM, contact identity, approval, jobs, adapter,
outbound fence and host. All new composition code is AUTHORED.

Apply migrations through0084_whatsapp_campaign.up.sql. Keep the existing
WhatsApp host and scheduled worker enabled. In its exact host profile, set
campaign_policy_file and campaign_policy_sha256 to an absolute policy file and
its SHA-256. A complete reference shape is in
config/whatsapp.campaign.reference.json; its reference purpose/policy values are
configuration examples, never evidence of real consent. The selected campaign
pack must be present. Missing/disabled schema guards prevent activation.

The policy must declare a marketing consent purpose distinct from the host's
appointment-notification purpose, and the actual contact/consent policy version.
The profile caps an audience at20members and a campaign at3steps. No new provider
credential path is introduced. The same existing account, exact approved
templates, identity and credential providers are used when the user configures
them. Local fixtures use only synthetic values and loopback provider traffic.

The authenticated operator API is /v1/franchise/marketing/campaigns:

- POST /prepare: campaign_id, explicit sources/states, members containing
  lead_id/recipient, and exact template steps with nonoverlapping due/expiry
  windows. Returns the current source/consent snapshot.
- POST the exact snapshot to the collection to store it and create the original
  pending schedule approvals. This operation can return complete:false with
  per-item codes; no failure is concealed as a fully prepared campaign.
- POST /{id}/resume with campaign_sha256 to resume preparation as its creator.
- POST /{id}/decision with campaign_sha256, approve and reason as a distinct
  reviewer. It explicitly reviews the stored member/step batch. Partial review
  remains visible and can be resumed.
- GET /{id}: immutable audience, stop state, each original approval/job/fence/
  signed notification status and first observed downstream quote/order.
- POST /{id}/stop with campaign_sha256 and reason stops future eligible work.

Use marketing:request plus notification:request for create/resume;
marketing:approve plus notification:approve for review; marketing:read plus
notification:read for state; marketing:cancel for stop. The service retains
notification:dispatch. Receipt reconciliation uses the original scheduled
notification endpoint and its existing notification:reconcile/read permissions.
Nothing manufactures the creator's or reviewer's identity.

Each member must match the explicitly selected source and one of new/contacted/
qualified, an active contact binding, current marketing consent and configured
organization. Source snapshots are immutable. A lifecycle/source/binding/consent
change suppresses the old reviewed effect and requires a new campaign snapshot.
No recipients are dynamically added after review. Requests/snapshots are bounded
to32KiB; template arity/body limits come from the existing approved profile.
Each due time is within30days, each expiry window within24hours, and steps do not
overlap. A follow-up requires an accepted predecessor; uncertain/rejected/
suppressed predecessors never authorize another provider attempt.

The same original worker claims typed jobs. Campaign jobs use
whatsapp.campaign.scheduled.v1 so a host without the campaign extension leaves
them untouched. When enabled, the existing worker alternates the two admitted
types. It retains the original finite lease generations, one job per poll,
durable fence and exact receipt recovery. A stop or opt-out cannot recall an
accepted/in-flight effect. Unknown effects are not retried as new messages.

An interrupted batch retry reads an already-stored decision only when request
hash, reviewer, reason, decision and its unique immutable record match. It
creates neither another approval nor another job. A changed decision/reviewer/
reason remains unresolved rather than rewriting historical approval.

Conversion is the first original quotation_acceptance/order for this member
after an accepted campaign message. The API reports its original IDs, acceptance
time, order state and evidence hash. This is observed sequence, not causal
attribution, payment proof, attributed revenue or a provider conversion upload.
Existing acceptance suppresses subsequent campaign steps.

Reference proofs: six campaigns,10jobs,4actual Python adapter POSTs,10durable
results and one original quote-to-order conversion;12concurrent workers give
one first effect. Further proof covers partial create/review/resume with2unique
requests/jobs/decisions and zero sends. Four decision-replay negatives are
read-only against those exact receipts. The host mounts both capabilities with
two existing worker loops, mTLS reports and joined shutdown. Populated rollback
is refused; empty down/up is tested. Original appointment behavior has its
specific compatibility regression after this source extension.

Retention/archival for campaign and recipient evidence belongs in the deployment
operational policy and T2809. These fixtures do not authorize deleting immutable
history or certify live delivery, production security or jurisdictional consent.
