# Pricebook exclusivity and appointment capacity: provenance boundary

Read-only assessment, 2026-09-11. No source reclassification, code change or broad
admission follows from this document.

## Pricebook exclusivity — CONDITIONED contract constraint

`Commerce.ActivatePriceBook` in `internal/platform/postgres/commerce.go` combines
two distinguishable parts. Its advisory lock and atomic status/outbox writes are
technical glue. Refusing overlapping active price sources protects the existing
single-valued lookup contract; it can be documented as fail-closed ownership glue
for the selected single-pricebook profile, because it never chooses a cheaper
price or creates a priority rule. The scope tuple tenant/market/currency and the
single-active-source profile must remain explicit constraints, not claims of BC
business semantics.

There is no demonstrated BC source for this exact exclusivity policy in the
inspected sources. BC PriceCalculationV16.PickBestLine244/IsBetterLine293 allows
multiple candidates and ranks them; treating that as equivalent would change
behavior. Hence no ADAPTED label is justified and no vendor attribution is made.
The surrounding pricebook lifecycle is not automatically all glue.

Trigger: if a target requires multiple concurrent pricebooks, customer-specific
price precedence, ranking or promotion interaction, reopen the source capability
before changing the single-source constraint. If a generalized exclusivity
algorithm is required as source-derived domain, status is RESEARCH_INCOMPLETE;
acquire an exact source only through the governed source profile. Current
mechanical anti-ambiguity guard does not itself require inventing such an engine.

## Appointment capacity — split technical quota from scheduling policy

`RequestAppointment`/`PublicAppointmentSlots` consume an explicitly supplied slot
capacity. Atomic `count(reservations) < configured capacity` plus locks/CAS is a
technical bounded-resource invariant, and can remain AUTHORED glue within that
declared reservation contract. It does not compute staffing, overbooking or
commercial capacity. Range containment/overlap operators are also technical
primitives, but the following choices are substantive existing policy:

- requested and confirmed appointments both consume capacity;
- working must fully contain a slot; any unavailable overlap excludes it;
- `prepareAppointmentSlot` hardcodes30minutes lead time, maximum8hours duration
  and capacity1..100; `RequestAppointment` also applies30minutes lead time;
- appointment kinds, resource skills and lifecycle transitions define which
  reservations are meaningful.

Those choices cannot all be reclassified as inevitable glue solely because the
SQL is short. No equivalent source algorithm has been demonstrated among the
selected BC pricing/quote/inventory files; the admitted inventory derivation
explicitly excludes calendars. Status for generalized scheduling is
RESEARCH_INCOMPLETE, not NO_ADMISSIBLE_SOURCE and not credential-conditioned.

Trigger: admit source plus tests for the required scheduling semantics, or bind
each policy as an explicit target configuration/contract and audit only the
mechanical enforcement as glue. Do not silently pick a staffing/timezone/legal
rule and do not expand implementation before that distinction is resolved. This
subtask leaves these owners unchanged and does not claim criterion9 globally.
