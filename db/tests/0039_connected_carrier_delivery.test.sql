begin;

do $$
begin
  if to_regclass('logistics.shipment_provider_report') is null then raise exception 'provider report table missing'; end if;
  if not exists(select 1 from information_schema.columns where table_schema='logistics' and table_name='shipment' and column_name='customer_shipment_id') then raise exception 'customer shipment link missing'; end if;
  if not exists(select 1 from pg_indexes where schemaname='logistics' and indexname='shipment_provider_reference_uq' and indexdef like '%WHERE (provider_reference IS NOT NULL)%') then raise exception 'partial provider reference identity missing'; end if;
  if not exists(select 1 from pg_trigger where tgname='shipment_provider_report_immutable') then raise exception 'provider report immutability missing'; end if;
end;
$$;

insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)
values('018f4d4a-7b36-7a21-8d10-2f4c54c29939','carrier-v174','Carrier V174','Carrier V174');
insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)
values
('018f4d4a-7b36-7a21-8d10-2f4c54c29939','origin','origin','Origin','store'),
('018f4d4a-7b36-7a21-8d10-2f4c54c29939','destination','destination','Destination','store');

insert into logistics.shipment(tenant_id,shipment_id,provider_code,origin_organization_id,destination_organization_id,state)
values
('018f4d4a-7b36-7a21-8d10-2f4c54c29939','null-ref-a','carrier','origin','destination','planned'),
('018f4d4a-7b36-7a21-8d10-2f4c54c29939','null-ref-b','carrier','origin','destination','planned');

insert into logistics.shipment(tenant_id,shipment_id,provider_code,provider_reference,origin_organization_id,destination_organization_id,state)
values('018f4d4a-7b36-7a21-8d10-2f4c54c29939','ref-a','carrier','provider-1','origin','destination','booked');

do $$
begin
  begin
    insert into logistics.shipment(tenant_id,shipment_id,provider_code,provider_reference,origin_organization_id,destination_organization_id,state)
    values('018f4d4a-7b36-7a21-8d10-2f4c54c29939','ref-b','carrier','provider-1','origin','destination','booked');
    raise exception 'duplicate provider reference accepted';
  exception when unique_violation then null;
  end;
end;
$$;

rollback;
