begin;
-- Fails atomically if multiple pending references now exist. Never delete,
-- merge or manufacture payment references to make a downgrade pass.
alter table payment.payment_attempt
  drop constraint payment_attempt_tenant_id_provider_code_provider_reference_key,
  add constraint payment_attempt_tenant_id_provider_code_provider_reference_key
    unique nulls not distinct (tenant_id,provider_code,provider_reference);
commit;
