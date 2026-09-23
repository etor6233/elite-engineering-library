# Governed AI reference

AUTHORED integration guide. This is library infrastructure for local fixtures
and optional provider activation; it does not certify a live model's quality.

The selected WhatsApp host assembles the existing Responses adapter, current
contact resolver, PostgreSQL turn store, lifetime tenant token reservation and
domain gateway. ConversationDomain resolves order status from the same tenant,
organization and lead's customer. A shared service subject is never substituted
for that customer's identity. The host keeps model replies as proposals; existing
human reply approval and the outbound fence govern actual sending.

Apply the selected migrations through0086. ConversationConfig is bound to the
host profile SHA; app wiring supplies the exact configured model and token cap.
Instructions and structured user/assistant roles accompany each initial request;
continuation repeats the trusted instructions. Model text is untrusted and the
fixture checks do not establish universal resistance to prompt injection.

Three offered tools use closed schemas: one-vehicle nonbinding quotation,
contact-scoped order status, and appointment request. The quotation owner
represents one vehicle; quantity other than1 is handed to the operator rather
than silently discarded. Prices, stock and permissions remain domain decisions.
Appointment requests follow the selected approval policy; a pending request
produces a terminal, replayable handoff with the exact proposed arguments in
conversation_turn. The operator continues the existing appointment workflow;
approving a different reply does not silently resume the model's old tool.

Each attempt reserves capacity before the provider request. Reservations survive
process restart and uncertain responses. Replaying the same reservation does not
charge twice, but a new attempt does. The lifetime cap has no scheduled refill;
a different configured cap fails closed and requires an explicit budget migration.
Reservations are conservative token capacity, not provider invoices or a promise
about cost. Actual usage and unresolved requests remain separate evidence.

Completion, retry and intent pinning require the active generation and live DB
lease. The first tool/arguments/contact/config intent cannot be reinterpreted on
retry. The existing domain owner provides effect idempotency; losing the second
model response can repeat the same command but creates only one quotation.
Remote model requests can be charged again after uncertainty. There is no claim
of exactly-once remote execution. Expired history/replay is denied; operational
purge and host scheduling belong to the same reference operations configuration.

## Executable evaluation

Run `python tools/verify_ai_reference.py --go GO_EXE --database-url OWNED_LOOPBACK_DB --receipt ABSENT_RECEIPT_JSON`
after migrations. It executes the real PostgreSQL and HTTP domain owners with
explicit synthetic model responses. Every case is required; average performance
cannot hide a failed isolation, authorization, recovery or handoff case. The
portable runner refuses SKIP and missing success markers. No provider key is used.

Real model promotion requires its own frozen cases and thresholds; synthetic
model decisions here only verify infrastructure and effects. The reusable
aifoundation EvalSuite supports caller cancellation and per-case decisions.

## Historical model adaptation is a separate opt-in lane

The library's `markdown_system/HISTORY_MODEL_TRAINING_PACK_PLAN.md` materializes
the admitted21-file HISTORY-MODEL-TRAINING-PIPELINE reference plus its supervisor.
Read its README and PROJECT_HISTORY_MODEL_TRAINING_CONTRACT before enabling it.
The default NEW profile imports and trains nothing. EXISTING requires a complete
purpose/source/permissions/private storage/retention/model/runtime/budget profile
before reading originals. Raw bytes, reviewed derivatives, isolated partitions,
explicit SFT job, independent evaluation, owner acceptance and rollback remain
separate stages. No historical message is sent to Handle or sales/payment tools.

V373 executed real official SFT on4592 random fixture parameters, independent
quality/privacy/abuse cases,13negative cases and real timeout/empty-process-tree
recovery. Current correspondence binds the nine executable training/evaluation
and supervisor files to the later successful job's code hash; the rebuilt21-file
profile also passes the current27 policy tests. The earlier manifest predated
four final audit corrections and is preserved as historical evidence. Those
random fixture weights are not a useful consumer model. No private
chats, derived datasets or trained weights are distributed with this product.
The optional CPU candidate pointer is separate from this Responses provider;
selecting a different/private model or serving target requires its own admission.
Do not replace training with RAG or conversation memory, infer consumer consent
from fixtures, or automatically activate a trained candidate.
