# Connected campaign composition plan

Historical implementation plan; the resulting local proof and runtime contract are in CAMPAIGN_CONNECTED_REFERENCE.md and CAMPAIGN_CONNECTED_RELEASE_V402.md/json. Reuse the existing
scheduled WhatsApp worker, original approval/job/fence/provider source, CRM
lead/contact/consent and quotation_acceptance/customer_order. No new marketing
engine, provider scheduler, consent authority or attribution formula.

One immutable campaign captures a bounded explicit audience: source/lifecycle
filters, selected lead IDs and recipient bindings, exact approved template
steps/times and source snapshots. At most20members and3steps. No arbitrary
predicate SQL or dynamic unreviewed recipient expansion. Marketing consent uses
an explicit configured purpose/policy distinct from appointment notification
consent. Future policy/account values are configuration, never inferred from a
synthetic fixture or acquisition of an upstream.

Each member/step gets its own original pending WhatsApp schedule approval.
Create/resume is idempotent; partial preparation is visible. The reviewer
explicitly approves the exact stored campaign hash and its pending member/step
hashes, using the original distinct-actor decision owner. Partial review is
visible/resumable; it never manufactures an approval principal. Original
service workers perform only the effects actually approved.

Dispatch rechecks the immutable campaign snapshot, current source/lifecycle,
contact binding, marketing consent, append-only campaign stop, original quote
acceptance and any preceding step. Follow-ups require the preceding step's
accepted fence; an uncertain or rejected predecessor suppresses that follow-up.
Step windows cannot overlap. No source or consent update recalls an accepted
provider effect. All unknown effects use the existing receipt recovery; no resend.

Conversion is a read-only observation from this campaign member's original
quotation_acceptance/order after an accepted campaign delivery. Report quote,
order, acceptance time and original evidence hash. Never claim marketing caused
the sale, charge twice, alter the order or upload a provider conversion.

Mount the campaign API on the existing authenticated role API and host through
an optional exact policy profile. Reuse the existing worker, mTLS reporter,
shutdown and status endpoint. Apply one schema delta; preserve appointment
schedule compatibility and populated rollback refusal.
