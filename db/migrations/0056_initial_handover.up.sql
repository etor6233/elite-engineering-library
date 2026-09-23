begin;

-- AUTHORED projection/binding glue. No price, credit, tax, stock posting or
-- shipment rules are defined by this table. Those owners remain authoritative.
create table sales.delivery_handover_preparation (
  tenant_id uuid not null,
  handover_id text not null,
  organization_id text not null,
  order_id text not null,
  order_line_id text not null,
  reservation_id text not null,
  payment_attempt_id text not null,
  observation_sha256_hex text not null check (observation_sha256_hex ~ '^[0-9a-f]{64}$'),
  contract_id text not null,
  contract_sha256_hex text not null check (contract_sha256_hex ~ '^[0-9a-f]{64}$'),
  contract_scope text not null,
  check ((contract_scope='LOCAL_FIXTURES' and contract_id='reference-single-unit-observed-payment-v1')
      or (contract_scope='MATERIALIZED_PROFILE' and contract_id ~ '^[a-z][a-z0-9._-]{0,79}@[1-9][0-9]{0,6}$')),
  maximum_observation_age_ns bigint not null check (maximum_observation_age_ns>0 and maximum_observation_age_ns<=900000000000),
  prepared_by_subject text not null check (length(prepared_by_subject) between 1 and 256),
  prepared_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,handover_id),
  unique (tenant_id,organization_id,order_id,order_line_id),
  foreign key (tenant_id,handover_id) references sales.delivery_handover(tenant_id,handover_id),
  foreign key (tenant_id,organization_id) references org.organization(tenant_id,organization_id),
  foreign key (tenant_id,order_id,order_line_id) references sales.customer_order_line(tenant_id,order_id,line_id),
  foreign key (tenant_id,reservation_id) references inventory.serial_reservation(tenant_id,reservation_id),
  foreign key (tenant_id,payment_attempt_id) references payment.payment_attempt(tenant_id,payment_attempt_id)
);

create function sales.initial_handover_preparation_immutable() returns trigger language plpgsql as $$
begin
  raise exception 'initial handover preparation is immutable';
end;
$$;
create trigger initial_handover_preparation_immutable
  before update or delete on sales.delivery_handover_preparation
  for each row execute function sales.initial_handover_preparation_immutable();

commit;
