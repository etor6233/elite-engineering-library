begin;

create table sales.return_exchange (
  tenant_id uuid not null,
  request_id text not null,
  disposition_id text not null,
  original_order_id text not null,
  original_line_id text not null,
  original_stock_unit_id text not null,
  replacement_order_id text not null,
  replacement_line_id text not null,
  replacement_stock_unit_id text not null,
  replacement_handover_id text not null,
  accounting_request_id text not null,
  currency text not null check (currency ~ '^[A-Z]{3}$'),
  original_unit_price_minor_units bigint not null check (original_unit_price_minor_units > 0),
  settlement_mode text not null check (settlement_mode='even-exchange-zero-balance'),
  status text not null check (status='prepared'),
  result_sha256_hex text not null check (result_sha256_hex ~ '^[0-9a-f]{64}$'),
  prepared_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,request_id),
  unique (tenant_id,disposition_id),
  unique (tenant_id,replacement_order_id),
  unique (tenant_id,replacement_stock_unit_id),
  unique (tenant_id,replacement_handover_id),
  foreign key (tenant_id,request_id) references sales.return_effect_request(tenant_id,request_id),
  foreign key (tenant_id,disposition_id) references sales.return_disposition(tenant_id,disposition_id),
  foreign key (tenant_id,original_order_id,original_line_id) references sales.customer_order_line(tenant_id,order_id,line_id),
  foreign key (tenant_id,original_stock_unit_id) references inventory.stock_unit(tenant_id,stock_unit_id),
  foreign key (tenant_id,replacement_order_id,replacement_line_id) references sales.customer_order_line(tenant_id,order_id,line_id),
  foreign key (tenant_id,replacement_stock_unit_id) references inventory.stock_unit(tenant_id,stock_unit_id),
  foreign key (tenant_id,replacement_handover_id) references sales.delivery_handover(tenant_id,handover_id),
  foreign key (tenant_id,accounting_request_id) references sales.return_effect_request(tenant_id,request_id),
  check (original_order_id<>replacement_order_id),
  check (original_stock_unit_id<>replacement_stock_unit_id)
);

create or replace function sales.prevent_return_exchange_mutation() returns trigger language plpgsql as $$
begin
  raise exception using errcode='23514',message='prepared return exchange is immutable';
end $$;

create trigger return_exchange_immutable before update or delete on sales.return_exchange for each row execute function sales.prevent_return_exchange_mutation();
create index return_exchange_original_order_idx on sales.return_exchange(tenant_id,original_order_id,prepared_at desc,request_id);

commit;
