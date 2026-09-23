begin;
-- AUTHORED durable representation, reservation and receipt glue. No financial
-- general ledger or external-provider payment is synthesized here.
alter table approval.request drop constraint request_kind_check;
alter table approval.request add constraint request_kind_check check(kind in
 ('reservation','sale','refund','payment','whatsapp_reply','social_publish','social_revoke','stored_value_operation'));
alter table approval.request drop constraint request_payload_binding_check;
alter table approval.request add constraint request_payload_binding_check check
 (kind not in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation') or
  (organization_id is not null and length(organization_id) between 1 and 128 and
   payload is not null and jsonb_typeof(payload)='object' and pg_column_size(payload)<=65536 and amount_minor_units=0));
create or replace function approval.guard_bound_request() returns trigger language plpgsql as $$
begin
 if tg_op='DELETE' then
  if old.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation') then raise exception 'bound approval request is immutable'; end if;
  return old;
 end if;
 if old.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation') or new.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation') then
  if row(new.tenant_id,new.request_id,new.kind,new.subject_id,new.amount_minor_units,new.requester,new.evidence_sha,new.organization_id,new.payload,new.created_at)
   is distinct from row(old.tenant_id,old.request_id,old.kind,old.subject_id,old.amount_minor_units,old.requester,old.evidence_sha,old.organization_id,old.payload,old.created_at)
   or (old.state<>'pending' and row(new.state,new.decided_at) is distinct from row(old.state,old.decided_at))
  then raise exception 'bound approval request is immutable'; end if;
 end if;
 return new;
end $$;
create or replace function approval.guard_bound_decision() returns trigger language plpgsql as $$
begin
 if exists(select 1 from approval.request where tenant_id=old.tenant_id and request_id=old.request_id and kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation')) then
  raise exception 'bound approval decision is immutable';
 end if;
 if tg_op='DELETE' then return old;end if;return new;
end $$;

create schema stored_value;
create table stored_value.account (
 tenant_id uuid not null,
 account_id text not null check(length(account_id) between 1 and 128),
 organization_id text not null,
 program_id text not null,
 kind text not null check(kind in ('gift_card','loyalty')),
 customer_subject text,
 currency text not null check(currency ~ '^[A-Z]{3}$'),
 created_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,account_id),
 foreign key(tenant_id,organization_id) references org.organization(tenant_id,organization_id),
 check((kind='loyalty' and customer_subject is not null) or kind='gift_card')
);
create table stored_value.operation (
 tenant_id uuid not null,
 operation_id text not null check(length(operation_id) between 1 and 128),
 organization_id text not null,
 order_id text not null,
 order_version bigint not null check(order_version>0),
 program_id text not null,
 operation text not null check(operation in ('issue','accrue','redeem','reverse')),
 original_operation_id text,
 approval_id text not null,
 profile_sha256 text not null check(profile_sha256 ~ '^[0-9a-f]{64}$'),
 request_sha256 text not null check(request_sha256 ~ '^[0-9a-f]{64}$'),
 calculation_sha256 text not null check(calculation_sha256 ~ '^[0-9a-f]{64}$'),
 receipt_sha256 text not null check(receipt_sha256 ~ '^[0-9a-f]{64}$'),
 receipt jsonb not null check(jsonb_typeof(receipt)='object' and pg_column_size(receipt)<=65536),
 created_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,operation_id),
 foreign key(tenant_id,organization_id) references org.organization(tenant_id,organization_id),
 foreign key(tenant_id,order_id) references sales.customer_order(tenant_id,order_id),
 foreign key(tenant_id,approval_id) references approval.request(tenant_id,request_id),
 foreign key(tenant_id,original_operation_id) references stored_value.operation(tenant_id,operation_id),
 unique(tenant_id,approval_id),
 check((operation='reverse')=(original_operation_id is not null))
);
create unique index one_stored_value_reversal on stored_value.operation(tenant_id,original_operation_id) where original_operation_id is not null;
create index stored_value_order_operations on stored_value.operation(tenant_id,order_id,program_id,created_at);
create table stored_value.entry (
 tenant_id uuid not null,
 operation_id text not null,
 account_id text not null,
 points_delta numeric(38,6) not null check(points_delta::text not in ('NaN','Infinity','-Infinity')),
 applied_minor_units bigint not null default 0,
 primary key(tenant_id,operation_id,account_id),
 foreign key(tenant_id,operation_id) references stored_value.operation(tenant_id,operation_id),
 foreign key(tenant_id,account_id) references stored_value.account(tenant_id,account_id),
 check(points_delta<>0 or applied_minor_units<>0)
);
create index stored_value_account_entries on stored_value.entry(tenant_id,account_id,operation_id);
create table stored_value.reservation (
 tenant_id uuid not null,
 operation_id text not null,
 account_id text not null,
 order_id text not null,
 approval_id text not null,
 points numeric(38,6) not null check(points>0 and points::text not in ('NaN','Infinity','-Infinity')),
 state text not null check(state in ('reserved','consumed','rejected','expired')),
 expires_at timestamptz not null,
 created_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,operation_id,account_id),
 foreign key(tenant_id,account_id) references stored_value.account(tenant_id,account_id),
 foreign key(tenant_id,order_id) references sales.customer_order(tenant_id,order_id),
 foreign key(tenant_id,approval_id) references approval.request(tenant_id,request_id) deferrable initially deferred,
 check(expires_at>created_at)
);
create index stored_value_live_reservations on stored_value.reservation(tenant_id,account_id,expires_at) where state='reserved';
create view stored_value.order_allocation as
 select o.tenant_id,o.order_id,o.organization_id,o.currency,o.version as order_version,o.total_minor_units as gross_minor_units,
 coalesce(a.gift_minor,0)::bigint as gift_minor_units,coalesce(a.discount_minor,0)::bigint as discount_minor_units,
 (o.total_minor_units-coalesce(a.gift_minor,0)-coalesce(a.discount_minor,0))::bigint as provider_due_minor_units
 from sales.customer_order o left join (
  select op.tenant_id,op.order_id,
  sum(e.applied_minor_units) filter(where a.kind='gift_card') as gift_minor,
  sum(e.applied_minor_units) filter(where a.kind='loyalty') as discount_minor
  from stored_value.operation op join stored_value.entry e using(tenant_id,operation_id)
  join stored_value.account a using(tenant_id,account_id)
  group by op.tenant_id,op.order_id
 ) a using(tenant_id,order_id);
create or replace view payment.order_funding as
 select a.tenant_id,a.order_id,a.organization_id,a.currency,a.order_version,a.gross_minor_units,a.gift_minor_units,a.discount_minor_units,a.provider_due_minor_units,
 encode(sha256(convert_to(coalesce((select jsonb_agg(jsonb_build_array(o.operation_id,o.receipt_sha256) order by o.operation_id)
 from stored_value.operation o where o.tenant_id=a.tenant_id and o.order_id=a.order_id
 and exists(select 1 from stored_value.entry e where e.tenant_id=o.tenant_id and e.operation_id=o.operation_id and e.applied_minor_units<>0)),'[]'::jsonb)::text,'UTF8')),'hex') as contributions_sha256
 from stored_value.order_allocation a;
create function stored_value.immutable_record() returns trigger language plpgsql as $$
begin raise exception 'stored value evidence is immutable; record a bound reversal';end $$;
create trigger stored_value_account_immutable before update or delete on stored_value.account for each row execute function stored_value.immutable_record();
create trigger stored_value_operation_immutable before update or delete on stored_value.operation for each row execute function stored_value.immutable_record();
create trigger stored_value_entry_immutable before update or delete on stored_value.entry for each row execute function stored_value.immutable_record();
create function stored_value.guard_reservation() returns trigger language plpgsql as $$
begin
 if tg_op='DELETE' then raise exception 'reservation evidence is immutable';end if;
 if row(new.tenant_id,new.operation_id,new.account_id,new.order_id,new.approval_id,new.points,new.expires_at,new.created_at)
  is distinct from row(old.tenant_id,old.operation_id,old.account_id,old.order_id,old.approval_id,old.points,old.expires_at,old.created_at)
  or (old.state<>'reserved' and new.state<>old.state)
 then raise exception 'reservation identity or terminal state is immutable';end if;
 return new;
end $$;
create trigger stored_value_reservation_immutable before update or delete on stored_value.reservation for each row execute function stored_value.guard_reservation();
-- A lock wait, slow IPC or outbox write must not turn an expired reviewed
-- snapshot into a committed effect. Deferred constraints execute at COMMIT,
-- after every statement, using the database clock rather than transaction time.
create function stored_value.require_live_proposal() returns trigger language plpgsql as $$
begin
 if new.kind='stored_value_operation' and
  not coalesce((new.payload->>'expires_at')::timestamptz>clock_timestamp(),false) then
  raise exception using errcode='23514', constraint='stored_value_lease_valid', message='stored-value proposal expired before commit';
 end if;
 return null;
end $$;
create constraint trigger stored_value_proposal_lease after insert on approval.request
 deferrable initially deferred for each row execute function stored_value.require_live_proposal();
create function stored_value.require_live_operation() returns trigger language plpgsql as $$
begin
 if not coalesce((new.receipt->'operation'->>'expires_at')::timestamptz>clock_timestamp(),false) then
  raise exception using errcode='23514', constraint='stored_value_lease_valid', message='stored-value operation expired before commit';
 end if;
 return null;
end $$;
create constraint trigger stored_value_operation_lease after insert on stored_value.operation
 deferrable initially deferred for each row execute function stored_value.require_live_operation();
commit;
