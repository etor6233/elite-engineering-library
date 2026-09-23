begin;

alter table inventory.warehouse_bin
  drop constraint warehouse_bin_bin_type_check;

alter table inventory.warehouse_bin
  add column bin_rank integer not null default 0 check (bin_rank between 0 and 1000000),
  add constraint warehouse_bin_bin_type_check
    check (bin_type in ('receive','ship','put-away','pick','putpick','qc'));

create table inventory.warehouse_receipt (
  tenant_id uuid not null,
  receipt_id text not null,
  organization_id text not null,
  receive_bin_id text not null,
  source_kind text not null,
  source_id text not null,
  posting_date date not null,
  status text not null check (status = 'posted'),
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, receipt_id),
  foreign key (tenant_id, organization_id, receive_bin_id)
    references inventory.warehouse_bin (tenant_id, organization_id, bin_id),
  unique (tenant_id, source_kind, source_id),
  check (length(receipt_id) between 1 and 128),
  check (length(source_kind) between 1 and 64),
  check (length(source_id) between 1 and 128)
);

create table inventory.warehouse_receipt_line (
  tenant_id uuid not null,
  receipt_id text not null,
  line_id text not null,
  item_id text not null,
  lot_id text,
  quantity numeric(20,6) not null check (quantity > 0),
  unit_cost numeric(20,4) not null check (unit_cost >= 0),
  inventory_entry_id text not null,
  primary key (tenant_id, receipt_id, line_id),
  foreign key (tenant_id, receipt_id)
    references inventory.warehouse_receipt (tenant_id, receipt_id),
  foreign key (tenant_id, item_id)
    references inventory.stock_item (tenant_id, item_id),
  foreign key (tenant_id, lot_id)
    references inventory.inventory_lot (tenant_id, lot_id),
  foreign key (tenant_id, inventory_entry_id)
    references inventory.bulk_inventory_entry (tenant_id, entry_id),
  unique (tenant_id, inventory_entry_id)
);

create table inventory.warehouse_activity (
  tenant_id uuid not null,
  activity_id text not null,
  organization_id text not null,
  activity_type text not null check (activity_type in ('put-away','pick')),
  status text not null check (status in ('open','registered','cancelled')),
  source_kind text not null,
  source_id text not null,
  source_line_id text not null default '',
  request_id text not null,
  assigned_to text not null default '',
  version bigint not null check (version > 0),
  created_at timestamptz not null default clock_timestamp(),
  registered_at timestamptz,
  primary key (tenant_id, activity_id),
  foreign key (tenant_id, organization_id)
    references org.organization (tenant_id, organization_id),
  unique (tenant_id, request_id),
  unique (tenant_id, activity_id, organization_id),
  check (length(activity_id) between 1 and 128),
  check (length(request_id) between 1 and 128),
  check ((status = 'registered' and registered_at is not null) or
         (status <> 'registered' and registered_at is null))
);

create table inventory.warehouse_activity_line (
  tenant_id uuid not null,
  activity_id text not null,
  organization_id text not null,
  line_id text not null,
  sequence_no integer not null check (sequence_no > 0),
  from_bin_id text not null,
  to_bin_id text not null,
  item_id text not null,
  lot_id text,
  reservation_id text,
  quantity numeric(20,6) not null check (quantity > 0),
  expiration_date date,
  source_entry_id text,
  primary key (tenant_id, activity_id, line_id),
  foreign key (tenant_id, activity_id)
    references inventory.warehouse_activity (tenant_id, activity_id),
  foreign key (tenant_id, activity_id, organization_id)
    references inventory.warehouse_activity (tenant_id, activity_id, organization_id),
  foreign key (tenant_id, organization_id, from_bin_id)
    references inventory.warehouse_bin (tenant_id, organization_id, bin_id),
  foreign key (tenant_id, organization_id, to_bin_id)
    references inventory.warehouse_bin (tenant_id, organization_id, bin_id),
  foreign key (tenant_id, item_id)
    references inventory.stock_item (tenant_id, item_id),
  foreign key (tenant_id, lot_id)
    references inventory.inventory_lot (tenant_id, lot_id),
  foreign key (tenant_id, reservation_id)
    references inventory.bulk_reservation (tenant_id, reservation_id),
  foreign key (tenant_id, source_entry_id)
    references inventory.bulk_inventory_entry (tenant_id, entry_id),
  unique (tenant_id, activity_id, sequence_no),
  check (from_bin_id <> to_bin_id)
);

create index warehouse_activity_open_scope_idx
  on inventory.warehouse_activity(tenant_id, organization_id, activity_type, created_at)
  where status = 'open';

create function inventory.reject_warehouse_evidence_mutation()
returns trigger language plpgsql as $function$
begin
  raise exception using errcode='55000', message='posted warehouse evidence is append-only';
end;
$function$;

create trigger warehouse_receipt_immutable
before update or delete on inventory.warehouse_receipt
for each row execute function inventory.reject_warehouse_evidence_mutation();

create trigger warehouse_receipt_line_immutable
before update or delete on inventory.warehouse_receipt_line
for each row execute function inventory.reject_warehouse_evidence_mutation();

create trigger warehouse_activity_line_immutable
before update or delete on inventory.warehouse_activity_line
for each row execute function inventory.reject_warehouse_evidence_mutation();

create function inventory.enforce_warehouse_activity_transition()
returns trigger language plpgsql as $function$
begin
  if old.status <> 'open' or new.status not in ('registered','cancelled') or
     new.version <> old.version + 1 or
     new.activity_id <> old.activity_id or new.organization_id <> old.organization_id or
     new.activity_type <> old.activity_type or new.source_kind <> old.source_kind or
     new.source_id <> old.source_id or new.source_line_id <> old.source_line_id or
     new.request_id <> old.request_id or new.created_at <> old.created_at then
    raise exception using errcode='55000', message='invalid warehouse activity transition';
  end if;
  return new;
end;
$function$;

create trigger warehouse_activity_transition_guard
before update on inventory.warehouse_activity
for each row execute function inventory.enforce_warehouse_activity_transition();

commit;
