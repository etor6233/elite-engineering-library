begin;

create table sales.return_receipt (
  tenant_id uuid not null,
  receipt_id text not null,
  authorization_id text not null,
  organization_id text not null,
  order_id text not null,
  stock_unit_id text not null,
  customer_principal_id text not null,
  received_serial_number text not null check (length(received_serial_number) between 1 and 128),
  condition_code text not null check (condition_code in ('sealed','opened','damaged','incomplete')),
  notes text not null check (length(notes) between 1 and 1000),
  evidence_sha256_hex text not null check (evidence_sha256_hex ~ '^[0-9a-f]{64}$'),
  received_by_subject text not null check (length(received_by_subject) between 1 and 255),
  received_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,receipt_id),
  unique (tenant_id,authorization_id),
  foreign key (tenant_id,authorization_id) references sales.return_authorization(tenant_id,authorization_id),
  foreign key (tenant_id,organization_id) references org.organization(tenant_id,organization_id),
  foreign key (tenant_id,order_id) references sales.customer_order(tenant_id,order_id),
  foreign key (tenant_id,stock_unit_id) references inventory.stock_unit(tenant_id,stock_unit_id),
  foreign key (tenant_id,customer_principal_id) references crm.customer_profile(tenant_id,customer_principal_id)
);

create or replace function sales.validate_return_receipt_insert() returns trigger language plpgsql as $$
declare expected record;
begin
  select a.organization_id,a.order_id,a.stock_unit_id,a.customer_principal_id,s.serial_number
    into expected
    from sales.return_authorization a
    join inventory.stock_unit s on s.tenant_id=a.tenant_id and s.stock_unit_id=a.stock_unit_id
   where a.tenant_id=new.tenant_id and a.authorization_id=new.authorization_id and a.state='authorized';
  if not found or (new.organization_id,new.order_id,new.stock_unit_id,new.customer_principal_id,new.received_serial_number)
      is distinct from (expected.organization_id,expected.order_id,expected.stock_unit_id,expected.customer_principal_id,expected.serial_number) then
    raise exception using errcode='23514',message='return receipt does not match authorization and durable serial';
  end if;
  return new;
end $$;
create trigger return_receipt_validate before insert on sales.return_receipt for each row execute function sales.validate_return_receipt_insert();

create table sales.return_disposition (
  tenant_id uuid not null,
  disposition_id text not null,
  receipt_id text not null,
  inventory_action text not null check (inventory_action in ('quarantine','restock','repair','scrap')),
  customer_remedy text not null check (customer_remedy in ('refund','exchange')),
  notes text not null check (length(notes) between 1 and 1000),
  decided_by_subject text not null check (length(decided_by_subject) between 1 and 255),
  decided_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,disposition_id),
  unique (tenant_id,receipt_id),
  foreign key (tenant_id,receipt_id) references sales.return_receipt(tenant_id,receipt_id)
);

create or replace function sales.validate_return_disposition_insert() returns trigger language plpgsql as $$
declare authorized_disposition text;
begin
  select a.disposition into authorized_disposition
    from sales.return_receipt r
    join sales.return_authorization a on a.tenant_id=r.tenant_id and a.authorization_id=r.authorization_id
   where r.tenant_id=new.tenant_id and r.receipt_id=new.receipt_id;
  if not found or (authorized_disposition='return' and new.customer_remedy<>'refund') or (authorized_disposition='exchange' and new.customer_remedy<>'exchange') then
    raise exception using errcode='23514',message='return disposition contradicts authorization';
  end if;
  return new;
end $$;
create trigger return_disposition_validate before insert on sales.return_disposition for each row execute function sales.validate_return_disposition_insert();

create table sales.return_effect_request (
  tenant_id uuid not null,
  request_id text not null,
  disposition_id text not null,
  effect_kind text not null check (effect_kind in ('inventory','refund','exchange','accounting','fiscal')),
  owner_context text not null check (owner_context in ('inventory','payment','fulfillment','accounting','fiscal')),
  state text not null check (state='requested'),
  idempotency_key text not null check (length(idempotency_key) between 16 and 128),
  requested_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,request_id),
  unique (tenant_id,disposition_id,effect_kind),
  unique (tenant_id,idempotency_key),
  foreign key (tenant_id,disposition_id) references sales.return_disposition(tenant_id,disposition_id),
  check ((effect_kind='inventory' and owner_context='inventory') or (effect_kind='refund' and owner_context='payment') or (effect_kind='exchange' and owner_context='fulfillment') or (effect_kind='accounting' and owner_context='accounting') or (effect_kind='fiscal' and owner_context='fiscal'))
);

create or replace function sales.prevent_return_processing_mutation() returns trigger language plpgsql as $$
begin
  raise exception using errcode='23514',message='return processing evidence is immutable';
end $$;
create trigger return_receipt_immutable before update or delete on sales.return_receipt for each row execute function sales.prevent_return_processing_mutation();
create trigger return_disposition_immutable before update or delete on sales.return_disposition for each row execute function sales.prevent_return_processing_mutation();
create trigger return_effect_request_immutable before update or delete on sales.return_effect_request for each row execute function sales.prevent_return_processing_mutation();

create index return_receipt_org_time_idx on sales.return_receipt(tenant_id,organization_id,received_at desc,receipt_id);
create index return_effect_request_state_idx on sales.return_effect_request(tenant_id,owner_context,state,requested_at,request_id);

commit;
