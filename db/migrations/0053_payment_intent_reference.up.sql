begin;
-- AUTHORED: absence of a provider receipt is not a shared provider identity.
-- Keep the original migration immutable and preserve non-null uniqueness.
alter table payment.payment_attempt
  drop constraint payment_attempt_tenant_id_provider_code_provider_reference_key,
  add constraint payment_attempt_tenant_id_provider_code_provider_reference_key
    unique nulls distinct (tenant_id,provider_code,provider_reference);
commit;
