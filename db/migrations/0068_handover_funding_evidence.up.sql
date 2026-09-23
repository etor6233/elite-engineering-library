begin;
-- AUTHORED sum-type relationship: real provider payment OR local full funding.
alter table sales.delivery_handover_preparation alter column payment_attempt_id drop not null;
alter table sales.delivery_handover_preparation add column funding_receipt_id text;
alter table sales.delivery_handover_preparation add constraint initial_handover_funding_kind check ((payment_attempt_id is not null) <> (funding_receipt_id is not null));
alter table sales.delivery_handover_preparation add foreign key(tenant_id,funding_receipt_id) references payment.local_funding_receipt(tenant_id,funding_id);
alter table sales.commercial_release_receipt alter column payment_attempt_id drop not null;
alter table sales.commercial_release_receipt add column funding_receipt_id text;
alter table sales.commercial_release_receipt add constraint commercial_release_funding_kind check ((payment_attempt_id is not null) <> (funding_receipt_id is not null));
alter table sales.commercial_release_receipt add constraint commercial_local_funding_generation check(funding_receipt_id is null or observation_generation=1);
alter table sales.commercial_release_receipt add foreign key(tenant_id,funding_receipt_id) references payment.local_funding_receipt(tenant_id,funding_id);
commit;
