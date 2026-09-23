begin;
-- AUTHORED service/approval/inventory binding; upstream owners remain authoritative.
alter table approval.request drop constraint request_kind_check;
alter table approval.request add constraint request_kind_check check(kind in
 ('reservation','sale','refund','payment','whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair'));
alter table approval.request drop constraint request_payload_binding_check;
alter table approval.request add constraint request_payload_binding_check check
 (kind not in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair') or
  (organization_id is not null and length(organization_id) between 1 and 128 and
   payload is not null and jsonb_typeof(payload)='object' and pg_column_size(payload)<=65536 and amount_minor_units=0));
create or replace function approval.guard_bound_request() returns trigger language plpgsql as $$
begin
 if tg_op='DELETE' then
  if old.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair') then raise exception 'bound approval request is immutable'; end if;
  return old;
 end if;
 if old.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair') or new.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair') then
  if row(new.tenant_id,new.request_id,new.kind,new.subject_id,new.amount_minor_units,new.requester,new.evidence_sha,new.organization_id,new.payload,new.created_at)
   is distinct from row(old.tenant_id,old.request_id,old.kind,old.subject_id,old.amount_minor_units,old.requester,old.evidence_sha,old.organization_id,old.payload,old.created_at)
   or (old.state<>'pending' and row(new.state,new.decided_at) is distinct from row(old.state,old.decided_at))
  then raise exception 'bound approval request is immutable'; end if;
 end if;
 return new;
end $$;
create or replace function approval.guard_bound_decision() returns trigger language plpgsql as $$
begin
 if exists(select 1 from approval.request where tenant_id=old.tenant_id and request_id=old.request_id and kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair')) then
  raise exception 'bound approval decision is immutable';
 end if;
 if tg_op='DELETE' then return old;end if;return new;
end $$;

create table service_ops.warranty_claim (
 tenant_id uuid not null, service_case_id text not null, warranty_id text not null,
 handover_id text not null, appointment_id text not null,
 organization_id text not null, factory_organization_id text not null,
 customer_principal_id text not null, stock_unit_id text not null,
 service_date date not null, profile_sha256 text not null check(profile_sha256 ~ '^[0-9a-f]{64}$'),
 parts_covered boolean not null, labor_covered boolean not null,
 opened_by text not null, request_sha256 text not null check(request_sha256 ~ '^[0-9a-f]{64}$'),
 opened_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,service_case_id),
 unique(tenant_id,appointment_id,stock_unit_id),
 foreign key(tenant_id,service_case_id) references service_ops.service_case(tenant_id,service_case_id),
 foreign key(tenant_id,warranty_id) references service_ops.warranty_activation(tenant_id,warranty_id),
 foreign key(tenant_id,appointment_id) references crm.appointment(tenant_id,appointment_id),
 foreign key(tenant_id,factory_organization_id) references org.organization(tenant_id,organization_id)
);
create trigger warranty_claim_immutable before update or delete on service_ops.warranty_claim
 for each row execute function service_ops.warranty_fact_immutable();
create table service_ops.warranty_claim_step (
 tenant_id uuid not null, service_case_id text not null, command_id text not null,
 version bigint not null check(version>0), kind text not null check(kind in ('opened','diagnosed','planned','decision','work','quality','accepted','reconciled','cancelled')),
 from_state text not null, to_state text not null, actor text not null check(length(actor) between 1 and 128),
 request_sha256 text not null check(request_sha256 ~ '^[0-9a-f]{64}$'),
 payload_sha256 text not null check(payload_sha256 ~ '^[0-9a-f]{64}$'),
 payload jsonb not null check(jsonb_typeof(payload)='object' and pg_column_size(payload)<=65536),
 recorded_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,service_case_id,version),
 unique(tenant_id,service_case_id,command_id),
 foreign key(tenant_id,service_case_id) references service_ops.warranty_claim(tenant_id,service_case_id)
);
create unique index warranty_single_step_kind on service_ops.warranty_claim_step(tenant_id,service_case_id,kind) where kind<>'quality';
create trigger warranty_claim_step_immutable before update or delete on service_ops.warranty_claim_step
 for each row execute function service_ops.warranty_fact_immutable();
create table service_ops.warranty_work_plan (
 tenant_id uuid not null, service_case_id text not null, approval_id text not null,
 payload_sha256 text not null check(payload_sha256 ~ '^[0-9a-f]{64}$'),
 payload jsonb not null check(jsonb_typeof(payload)='object' and pg_column_size(payload)<=65536),
 expires_at timestamptz not null, created_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,service_case_id),unique(tenant_id,approval_id),
 foreign key(tenant_id,service_case_id) references service_ops.warranty_claim(tenant_id,service_case_id),
 foreign key(tenant_id,approval_id) references approval.request(tenant_id,request_id),
 check(expires_at>created_at)
);
create trigger warranty_work_plan_immutable before update or delete on service_ops.warranty_work_plan
 for each row execute function service_ops.warranty_fact_immutable();
create table service_ops.warranty_part_reservation (
 tenant_id uuid not null, service_case_id text not null, line_id text not null, reservation_id text not null,
 primary key(tenant_id,service_case_id,line_id),unique(tenant_id,reservation_id),
 foreign key(tenant_id,service_case_id) references service_ops.warranty_work_plan(tenant_id,service_case_id),
 foreign key(tenant_id,reservation_id) references inventory.bulk_reservation(tenant_id,reservation_id)
);
create trigger warranty_part_reservation_immutable before update or delete on service_ops.warranty_part_reservation
 for each row execute function service_ops.warranty_fact_immutable();

create function service_ops.warranty_claim_case_guard() returns trigger language plpgsql as $$
begin
 if not exists(select 1 from service_ops.warranty_claim c where c.tenant_id=old.tenant_id and c.service_case_id=old.service_case_id) then
  if tg_op='DELETE' then return old;end if;return new;
 end if;
 if tg_op='DELETE' then raise exception 'warranty claim case cannot be deleted';end if;
 if row(new.tenant_id,new.service_case_id,new.organization_id,new.stock_unit_id,new.severity,new.description,new.opened_at)
 is distinct from row(old.tenant_id,old.service_case_id,old.organization_id,old.stock_unit_id,old.severity,old.description,old.opened_at)
 then raise exception 'warranty case identity is immutable';end if;
 if new.version<>old.version+1 or not exists(select 1 from service_ops.warranty_claim_step s
 where s.tenant_id=old.tenant_id and s.service_case_id=old.service_case_id and s.version=new.version
 and s.from_state=old.state and s.to_state=new.state)
 then raise exception 'warranty case requires its connected command receipt';end if;
 return new;
end $$;
create trigger warranty_claim_case_guard before update or delete on service_ops.service_case
 for each row execute function service_ops.warranty_claim_case_guard();

-- A pending reservation may expire while a transaction is waiting. Approval
-- binds it durably to the repair; it is no longer a pending reservation lease.
create function service_ops.warranty_pending_commit_guard() returns trigger language plpgsql as $$
declare expiry timestamptz;
begin
 if tg_table_name='warranty_work_plan' then expiry:=new.expires_at;
 elsif new.kind='decision' and new.to_state='repair' then
  select expires_at into expiry from service_ops.warranty_work_plan where tenant_id=new.tenant_id and service_case_id=new.service_case_id;
 else return null;end if;
 if expiry is null or expiry<=clock_timestamp() then raise exception 'warranty pending reservation expired before commit';end if;
 return null;
end $$;
create constraint trigger warranty_plan_commit_expiry after insert on service_ops.warranty_work_plan
 deferrable initially deferred for each row execute function service_ops.warranty_pending_commit_guard();
create constraint trigger warranty_approval_commit_expiry after insert on service_ops.warranty_claim_step
 deferrable initially deferred for each row execute function service_ops.warranty_pending_commit_guard();
commit;
