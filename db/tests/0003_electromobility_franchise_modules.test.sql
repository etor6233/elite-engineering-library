begin;

insert into platform.tenant (tenant_id, tenant_code, legal_name, display_name)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c28c01', 'mobility-test', 'Mobility Test S.A.', 'Mobility Test');
insert into org.organization (tenant_id, organization_id, organization_code, display_name, organization_type)
values
 ('018f4d4a-7b36-7a21-8d10-2f4c54c28c01', 'hq', 'hq', 'Headquarters', 'franchisor'),
 ('018f4d4a-7b36-7a21-8d10-2f4c54c28c01', 'franchise-1', 'franchise-1', 'Franchise 1', 'franchisee');
insert into catalog.vehicle_model values (
 '018f4d4a-7b36-7a21-8d10-2f4c54c28c01','model-1','moto-one','Moto One','motorcycle','active',true,
 '{"rangeKm":120}',clock_timestamp(),clock_timestamp()
);
insert into catalog.vehicle_variant values (
 '018f4d4a-7b36-7a21-8d10-2f4c54c28c01','variant-1','model-1','moto-one-lfp','Moto One LFP',
 '{"chemistry":"LFP","capacityWh":4000}','approved','active'
);
insert into crm.customer_profile (tenant_id,customer_principal_id,display_name,email_normalized)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c28c01','customer-1','Customer One','customer@example.com');
insert into crm.lead (tenant_id,lead_id,organization_id,customer_principal_id,model_id,lifecycle_state,source_code,contact_payload)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c28c01','lead-1','franchise-1','customer-1','model-1','qualified','google-ads','{"email":"customer@example.com"}');
insert into crm.consent_evidence values (
 '018f4d4a-7b36-7a21-8d10-2f4c54c28c01','consent-1','lead-1','customer-1','sales-contact','v1','granted',clock_timestamp(),repeat('a',64)
);
insert into partner.supplier (tenant_id,supplier_id,supplier_code,legal_name,status)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c28c01','supplier-1','factory-one','Factory One','active');
insert into procurement.purchase_order (tenant_id,purchase_order_id,supplier_id,destination_organization_id,state,currency,total_minor_units,version)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c28c01','po-1','supplier-1','hq','received','USD',100000,1);
insert into factory.production_unit (tenant_id,production_unit_id,purchase_order_id,variant_id,serial_number,vin,battery_serial_number,state)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c28c01','production-1','po-1','variant-1','SERIAL-1','VIN-1','BATTERY-1','received');
insert into inventory.stock_unit (tenant_id,stock_unit_id,organization_id,variant_id,production_unit_id,serial_number,vin,battery_serial_number,state,version,received_at)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c28c01','stock-1','franchise-1','variant-1','production-1','SERIAL-1','VIN-1','BATTERY-1','available',1,clock_timestamp());
insert into sales.customer_order (tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c28c01','order-1','franchise-1','customer-1','allocated','USD',120000,1);
insert into sales.customer_order_line values (
 '018f4d4a-7b36-7a21-8d10-2f4c54c28c01','order-1','line-1','variant-1',1,120000,'stock-1'
);
insert into payment.payment_attempt (tenant_id,payment_attempt_id,order_id,provider_code,provider_reference,idempotency_key,state,currency,amount_minor_units,version)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c28c01','payment-1','order-1','sandbox','provider-payment-1','idem-payment-1','captured','USD',120000,1);
insert into service_ops.warranty values (
 '018f4d4a-7b36-7a21-8d10-2f4c54c28c01','warranty-1','stock-1','customer-1',clock_timestamp(),clock_timestamp()+interval '2 years','terms-v1','active'
);
insert into franchise.agreement (tenant_id,agreement_id,franchise_organization_id,territory_code,terms_version,starts_on,status)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c28c01','agreement-1','franchise-1','AR-CBA','v1',current_date,'active');

do $test$
begin
  begin
    insert into inventory.stock_unit (tenant_id,stock_unit_id,organization_id,variant_id,serial_number,state,version,received_at)
    values ('018f4d4a-7b36-7a21-8d10-2f4c54c28c01','stock-duplicate','franchise-1','variant-1','SERIAL-1','available',1,clock_timestamp());
    raise exception 'expected duplicate serial rejection';
  exception when unique_violation then null;
  end;
end;
$test$;

do $test$
declare
  model_count bigint;
  paid_total bigint;
begin
  select count(*) into model_count from catalog.vehicle_model where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c28c01' and publicly_visible;
  select sum(amount_minor_units) into paid_total from payment.payment_attempt where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c28c01' and state='captured';
  if model_count <> 1 or paid_total <> 120000 then
    raise exception 'vertical fixture mismatch';
  end if;
end;
$test$;

rollback;
