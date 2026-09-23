begin;

create table payment.return_refund (
  tenant_id uuid not null,
  request_id text not null,
  payment_attempt_id text not null,
  order_id text not null,
  line_id text not null,
  stock_unit_id text not null,
  provider_code text not null check (provider_code in ('stripe','mercado_pago')),
  provider_payment_reference text not null check (length(provider_payment_reference) between 1 and 255),
  idempotency_key text not null check (length(idempotency_key) between 16 and 128),
  currency text not null check (currency ~ '^[A-Z]{3}$'),
  amount_minor_units bigint not null check (amount_minor_units>0),
  provider_refund_reference text,
  provider_status text,
  state text not null check (state in ('prepared','pending','requires_action','blocked','succeeded','failed','canceled')),
  response_sha256_hex text,
  version bigint not null default 1 check (version>0),
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,request_id),
  foreign key (tenant_id,request_id) references sales.return_effect_request(tenant_id,request_id),
  foreign key (tenant_id,payment_attempt_id) references payment.payment_attempt(tenant_id,payment_attempt_id),
  foreign key (tenant_id,order_id,line_id) references sales.customer_order_line(tenant_id,order_id,line_id),
  foreign key (tenant_id,stock_unit_id) references inventory.stock_unit(tenant_id,stock_unit_id),
  unique (tenant_id,idempotency_key),
  check (provider_refund_reference is null or length(provider_refund_reference) between 1 and 255),
  check (response_sha256_hex is null or response_sha256_hex ~ '^[0-9a-f]{64}$'),
  check ((state='prepared' and provider_refund_reference is null and provider_status is null and response_sha256_hex is null) or
         (state<>'prepared' and provider_refund_reference is not null and provider_status is not null and response_sha256_hex is not null))
);

create unique index return_refund_provider_reference_uq
  on payment.return_refund(tenant_id,provider_code,provider_refund_reference)
  where provider_refund_reference is not null;

create table payment.return_refund_observation (
  tenant_id uuid not null,
  observation_id text not null,
  request_id text not null,
  provider_refund_reference text not null,
  provider_status text not null,
  amount_minor_units bigint not null check (amount_minor_units>0),
  currency text not null check (currency ~ '^[A-Z]{3}$'),
  response_sha256_hex text not null check (response_sha256_hex ~ '^[0-9a-f]{64}$'),
  observed_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,observation_id),
  foreign key (tenant_id,request_id) references payment.return_refund(tenant_id,request_id)
);

create or replace function payment.enforce_return_refund_lifecycle() returns trigger language plpgsql as $$
begin
  if tg_op='DELETE' then
    raise exception using errcode='23514',message='return refund projection cannot be deleted';
  end if;
  if (new.tenant_id,new.request_id,new.payment_attempt_id,new.order_id,new.line_id,new.stock_unit_id,new.provider_code,new.provider_payment_reference,new.idempotency_key,new.currency,new.amount_minor_units,new.created_at)
       is distinct from
     (old.tenant_id,old.request_id,old.payment_attempt_id,old.order_id,old.line_id,old.stock_unit_id,old.provider_code,old.provider_payment_reference,old.idempotency_key,old.currency,old.amount_minor_units,old.created_at)
     or new.version<>old.version+1 or new.updated_at<old.updated_at
     or old.state in ('succeeded','failed','canceled')
     or (old.state='prepared' and new.state not in ('pending','requires_action','blocked','succeeded','failed','canceled'))
     or (old.state='pending' and new.state not in ('pending','requires_action','blocked','succeeded','failed','canceled'))
     or (old.state='requires_action' and new.state not in ('requires_action','pending','blocked','succeeded','failed','canceled'))
     or (old.state='blocked' and new.state not in ('blocked','pending','requires_action','succeeded','failed','canceled')) then
    raise exception using errcode='23514',message='invalid return refund transition';
  end if;
  return new;
end $$;

create trigger return_refund_lifecycle before update or delete on payment.return_refund
for each row execute function payment.enforce_return_refund_lifecycle();

create or replace function payment.prevent_return_refund_observation_mutation() returns trigger language plpgsql as $$
begin
  raise exception using errcode='23514',message='return refund observation is immutable';
end $$;

create trigger return_refund_observation_immutable before update or delete on payment.return_refund_observation
for each row execute function payment.prevent_return_refund_observation_mutation();

create index return_refund_payment_idx on payment.return_refund(tenant_id,payment_attempt_id,state,request_id);
create index return_refund_observation_request_idx on payment.return_refund_observation(tenant_id,request_id,observed_at,observation_id);

commit;
