begin;

do $$
begin
  if to_regclass('payment.return_refund') is null or to_regclass('payment.return_refund_observation') is null then
    raise exception 'return refund tables missing';
  end if;
  if (select count(*) from pg_trigger where tgname in ('return_refund_lifecycle','return_refund_observation_immutable') and not tgisinternal)<>2 then
    raise exception 'return refund lifecycle triggers missing';
  end if;
  if not exists (select 1 from pg_indexes where schemaname='payment' and indexname='return_refund_provider_reference_uq') then
    raise exception 'return refund provider identity index missing';
  end if;
end $$;

rollback;
