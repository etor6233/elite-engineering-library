begin;

alter table logistics.shipment
  drop constraint shipment_tenant_id_provider_code_provider_reference_key,
  alter column destination_organization_id drop not null,
  add column destination_customer_principal_id text,
  add column provider_service_code text,
  add column request_id text,
  add column customer_shipment_id text,
  add column order_id text,
  add column version bigint not null default 1 check (version > 0),
  add constraint shipment_destination_exactly_one_ck check (
    (destination_organization_id is not null and destination_customer_principal_id is null)
    or (destination_organization_id is null and destination_customer_principal_id is not null)
  ),
  add constraint shipment_customer_source_pair_ck check ((customer_shipment_id is null) = (order_id is null)),
  add constraint shipment_customer_destination_fk foreign key (tenant_id,destination_customer_principal_id)
    references crm.customer_profile(tenant_id,customer_principal_id),
  add constraint shipment_customer_shipment_fk foreign key (tenant_id,customer_shipment_id)
    references sales.customer_shipment(tenant_id,shipment_id),
  add constraint shipment_customer_order_fk foreign key (tenant_id,order_id)
    references sales.customer_order(tenant_id,order_id),
  add constraint shipment_request_id_ck check (request_id is null or length(request_id) between 1 and 128),
  add constraint shipment_provider_service_code_ck check (provider_service_code is null or length(provider_service_code) between 1 and 128);

create unique index shipment_provider_reference_uq
  on logistics.shipment(tenant_id,provider_code,provider_reference)
  where provider_reference is not null;

create unique index shipment_request_uq
  on logistics.shipment(tenant_id,request_id)
  where request_id is not null;

create unique index shipment_customer_shipment_uq
  on logistics.shipment(tenant_id,customer_shipment_id)
  where customer_shipment_id is not null;

create table logistics.shipment_provider_report (
  tenant_id uuid not null,
  shipment_id text not null,
  provider_event_id text not null,
  provider_status text not null check (length(provider_status) between 1 and 64),
  normalized_state text not null check (normalized_state in ('picked-up','in-transit','delivered','exception')),
  report_schema text not null check (report_schema='elite-amazon-spapi-easyship-reconciliation-receipt/v1'),
  evidence_sha256_hex text not null check (evidence_sha256_hex ~ '^[0-9a-f]{64}$'),
  occurred_at timestamptz not null,
  received_at timestamptz not null default clock_timestamp(),
  automatic_business_delivery_acceptance boolean not null default false check (not automatic_business_delivery_acceptance),
  shipment_version bigint not null check (shipment_version > 1),
  primary key (tenant_id,shipment_id,provider_event_id),
  foreign key (tenant_id,shipment_id) references logistics.shipment(tenant_id,shipment_id),
  unique (tenant_id,shipment_id,shipment_version)
);

create function logistics.reject_shipment_provider_report_mutation()
returns trigger language plpgsql as $function$
begin
  raise exception using errcode='55000', message='shipment provider report evidence is immutable';
end;
$function$;

create trigger shipment_provider_report_immutable
before update or delete on logistics.shipment_provider_report
for each row execute function logistics.reject_shipment_provider_report_mutation();

create index shipment_provider_report_time_idx
  on logistics.shipment_provider_report(tenant_id,shipment_id,occurred_at,provider_event_id);

commit;
