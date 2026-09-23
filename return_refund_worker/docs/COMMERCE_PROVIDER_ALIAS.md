# Commerce to refund provider compatibility

`mercadopago` is the Commerce payment code. `mercado_pago` remains the existing
refund owner's code, provider registration name and durable refund projection.
`PostgresStore.Prepare` maps the former to the latter once, while selecting a new
refund's source payment. Both exact spellings are accepted; no case folding or
other alias is introduced. Stripe is unchanged.

The source payment row is never renamed. Existing refund rows are returned before
this mapping and remain byte-for-byte unchanged in their identity fields. Payment
attempt IDs, source provider references, return request IDs and idempotency keys
are carried through without modification. The existing SQL constraint and SDK
adapter continue to receive `mercado_pago`; no historical migration is needed.

This is AUTHORED compatibility glue. The official Mercado Pago Go SDK remains
dependency-pinned at 1.14.0. No financial rules, rounding, refund policy, SDK source
or attribution changes in this delta.

The integration fixture seeds the same explicitly disposable database as the
legacy refund test, including its documented trigger-bypass setup/cleanup. It
then runs actual store queries/transactions and pinned SDK requests through an
in-process transport for both source spellings. It proves partial and full refund
states, persisted identity, GET reconciliation, exact idempotency headers,
ambiguity rejection and immutable observations. It does not prove originating
return authorization journeys or live provider acceptance.
