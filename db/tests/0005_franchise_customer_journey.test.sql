begin;

insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29a01','journey-test','Journey Test','Journey Test');
insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29a01','store','store','Store','store');
insert into org.public_location(tenant_id,organization_id,city,region,country,published)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29a01','store','Cordoba','Cordoba','AR',true);
insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29a01','model','model','Model','bicycle','active');
insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29a01','variant','model','variant','Variant','{}','active');
insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29a01','customer','Customer','customer@example.test');
insert into crm.lead(tenant_id,lead_id,organization_id,customer_principal_id,model_id,lifecycle_state,source_code,contact_payload)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29a01','lead','store','customer','model','new','public-web','{}');
insert into crm.appointment(tenant_id,appointment_id,organization_id,lead_id,customer_principal_id,model_id,appointment_kind,starts_at,state,version)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29a01','appointment','store','lead','customer','model','test-drive',clock_timestamp()+interval '1 hour','requested',1);
insert into pricing.price_book(tenant_id,price_book_id,market,currency,valid_from,status)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29a01','retail','AR','ARS',clock_timestamp()-interval '1 day','active');
insert into pricing.price_book_entry(tenant_id,price_book_id,variant_id,amount_minor_units,tax_mode)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29a01','retail','variant',100000,'inclusive');
insert into sales.quotation(tenant_id,quotation_id,organization_id,lead_id,customer_principal_id,variant_id,price_book_id,currency,total_minor_units,valid_until,state,version)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29a01','quote','store','lead','customer','variant','retail','ARS',100000,clock_timestamp()+interval '7 days','issued',1);
insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,state,version,received_at)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29a01','stock','store','variant','SERIAL','sold',1,clock_timestamp());
insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29a01','order','store','customer','delivered','ARS',100000,1);
insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29a01','handover','store','order','customer','stock','prepared',1);

do $$
begin
  begin
    update crm.lead set lifecycle_state='converted',version=version+1 where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c29a01' and lead_id='lead';
    raise exception 'test failed: database alone allowed a skipped lead state';
  exception when check_violation then null;
  end;
end $$;

select 1 / case when (select count(*) from org.public_location where published)=1 then 1 else 0 end;
select 1 / case when (select count(*) from crm.appointment)=1 then 1 else 0 end;
select 1 / case when (select count(*) from sales.quotation)=1 then 1 else 0 end;
select 1 / case when (select count(*) from sales.delivery_handover)=1 then 1 else 0 end;

rollback;
