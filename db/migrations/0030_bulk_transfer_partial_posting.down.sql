begin;

do $migration$
begin
  if exists (select 1 from inventory.bulk_transfer_shipment)
     or exists (select 1 from inventory.bulk_transfer_receipt) then
    raise exception using errcode='55000', message='cannot remove partial transfer posting schema after postings exist';
  end if;
end;
$migration$;

drop index inventory.bulk_transfer_in_transit_idx;
drop view inventory.bulk_inventory_in_transit;
drop table inventory.bulk_transfer_receipt_component;

alter table inventory.bulk_transfer_cost_component
  drop constraint bulk_transfer_cost_component_received_check,
  drop column received_quantity;

alter table inventory.bulk_transfer_allocation
  drop constraint bulk_transfer_allocation_shipment_fk,
  drop column shipment_id;

drop table inventory.bulk_transfer_receipt;
drop table inventory.bulk_transfer_shipment;

alter table inventory.bulk_transfer
  drop constraint bulk_transfer_state_v164_check,
  drop constraint bulk_transfer_status_v164_check,
  add constraint bulk_transfer_status_check
    check (status in ('released','shipped','received','cancelled')),
  add constraint bulk_transfer_state_check
    check (
      (status='released' and shipped_at is null and received_at is null and put_away_activity_id is null) or
      (status='shipped' and shipped_at is not null and received_at is null and put_away_activity_id is null) or
      (status='received' and shipped_at is not null and received_at is not null and put_away_activity_id is not null) or
      (status='cancelled' and received_at is null and put_away_activity_id is null)
    );

create index bulk_transfer_in_transit_idx
  on inventory.bulk_transfer(tenant_id,to_organization_id,status,posting_date)
  where status='shipped';

create view inventory.bulk_inventory_in_transit as
select t.tenant_id,t.transfer_id,t.from_organization_id,t.to_organization_id,t.in_transit_code,
       l.line_id,l.item_id,a.lot_id,sum(c.quantity) quantity,sum(c.cost_amount) cost_amount,
       t.shipped_at,t.version
from inventory.bulk_transfer t
join inventory.bulk_transfer_line l using(tenant_id,transfer_id)
join inventory.bulk_transfer_allocation a using(tenant_id,transfer_id,line_id)
join inventory.bulk_transfer_cost_component c using(tenant_id,transfer_id,allocation_id)
where t.status='shipped'
group by t.tenant_id,t.transfer_id,t.from_organization_id,t.to_organization_id,t.in_transit_code,
         l.line_id,l.item_id,a.lot_id,t.shipped_at,t.version;

comment on view inventory.bulk_inventory_in_transit is
  'BC-derived transfer order evidence: shipped demand remains unavailable at source and visible in transit until destination receipt.';

commit;
