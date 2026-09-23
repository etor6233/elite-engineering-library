begin;

drop index if exists logistics.shipment_provider_report_time_idx;
drop trigger if exists shipment_provider_report_immutable on logistics.shipment_provider_report;
drop function if exists logistics.reject_shipment_provider_report_mutation();
drop table if exists logistics.shipment_provider_report;

drop index if exists logistics.shipment_customer_shipment_uq;
drop index if exists logistics.shipment_request_uq;
drop index if exists logistics.shipment_provider_reference_uq;

alter table logistics.shipment
  drop constraint if exists shipment_provider_service_code_ck,
  drop constraint if exists shipment_request_id_ck,
  drop constraint if exists shipment_customer_order_fk,
  drop constraint if exists shipment_customer_shipment_fk,
  drop constraint if exists shipment_customer_destination_fk,
  drop constraint if exists shipment_customer_source_pair_ck,
  drop constraint if exists shipment_destination_exactly_one_ck,
  drop column if exists version,
  drop column if exists order_id,
  drop column if exists customer_shipment_id,
  drop column if exists request_id,
  drop column if exists provider_service_code,
  drop column if exists destination_customer_principal_id,
  alter column destination_organization_id set not null,
  add unique nulls not distinct (tenant_id,provider_code,provider_reference);

commit;
