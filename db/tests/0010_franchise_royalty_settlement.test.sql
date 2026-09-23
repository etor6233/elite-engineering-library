begin;

insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29000','royalty-test','Royalty Test','Royalty Test');
insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29000','franchise-1','franchise-1','Franchise One','franchisee');
insert into franchise.agreement(tenant_id,agreement_id,franchise_organization_id,territory_code,terms_version,starts_on,status)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29000','agreement-1','franchise-1','AR-BUE-001','terms-v1','2026-01-01','active');
insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29000','order-1','franchise-1','customer-1','confirmed','ARS',10001,1);
insert into payment.payment_attempt(tenant_id,payment_attempt_id,order_id,provider_code,provider_reference,idempotency_key,state,currency,amount_minor_units,version)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29000','payment-1','order-1','provider','provider-1','idem-1','captured','ARS',10001,3);
insert into royalty.policy(tenant_id,policy_id,agreement_id,franchise_organization_id,currency,rate_basis_points,valid_from)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29000','policy-1','agreement-1','franchise-1','ARS',650,'2026-01-01');
insert into royalty.accrual(tenant_id,accrual_id,policy_id,agreement_id,franchise_organization_id,payment_attempt_id,source_event_key,source_state,currency,basis_minor_units,royalty_minor_units,occurred_at)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29000','accrual-1','policy-1','agreement-1','franchise-1','payment-1','provider:capture-1','captured','ARS',10001,650,'2026-08-01');
insert into royalty.settlement_run(tenant_id,settlement_id,franchise_organization_id,currency,period_start,period_end,status,expected_minor_units,version,closed_at)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29000','settlement-1','franchise-1','ARS','2026-08-01','2026-09-01','closed',650,2,clock_timestamp());
insert into royalty.settlement_line(tenant_id,settlement_id,line_id,accrual_id,amount_minor_units)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29000','settlement-1','line-1','accrual-1',650);
insert into royalty.reconciliation(tenant_id,reconciliation_id,settlement_id,external_reference,expected_minor_units,actual_minor_units,difference_minor_units,status,recorded_by)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29000','reconciliation-1','settlement-1','bank-statement-1',650,649,-1,'mismatch','controller-1');

do $$
declare total bigint; difference bigint;
begin
  select sum(amount_minor_units) into total from royalty.settlement_line where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c29000' and settlement_id='settlement-1';
  if total <> 650 then raise exception 'settlement total mismatch'; end if;
  select difference_minor_units into difference from royalty.reconciliation where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c29000' and settlement_id='settlement-1';
  if difference <> -1 then raise exception 'reconciliation difference mismatch'; end if;
  begin
    update royalty.accrual set royalty_minor_units=1 where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c29000' and accrual_id='accrual-1';
    raise exception 'immutable accrual accepted mutation';
  exception when raise_exception then
    if sqlerrm <> 'immutable royalty history' then raise; end if;
  end;
end $$;

rollback;
