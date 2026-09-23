# TikTok Lead durable import

This integration consumes only the three hash-linked artifacts emitted by `PYTHON-TIKTOK-LEAD-ADAPTER` 0.2.0. It verifies SDK/wheel/API/contract identity, hashes, retrieval mode, account/page identity hashes, provider metadata, candidate fields and outcome counts before the first write.

The import reuses the provider-neutral `leadstream.Store` and PostgreSQL migration 0044. It creates no parallel CRM, raw-event table, candidate table or outbox. A valid candidate is stored with `contact_eligibility=pending_policy`; a rejected dynamic field remains evidence with a normalization code. Replay is idempotent and a reused provider lead ID with different source bytes fails closed.

This pack does not authenticate TikTok webhook calls. The webhook remains an untrusted hint until the project demonstrates a provider-supported or independently approved ingress control. The authenticated retrieval response—not the webhook—is the input to durable candidate import. Live account, terms/consent, callback, delivery, reconciliation, retention and cost still require project evidence.
