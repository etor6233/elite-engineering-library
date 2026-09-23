begin;

alter table sales.delivery_handover add column supersedes_handover_id text;
alter table sales.delivery_handover add constraint delivery_handover_supersedes_fk foreign key (tenant_id,supersedes_handover_id) references sales.delivery_handover(tenant_id,handover_id);
alter table sales.delivery_handover add constraint delivery_handover_not_self_superseding_ck check (supersedes_handover_id is null or supersedes_handover_id<>handover_id);
create unique index delivery_handover_one_successor_uq on sales.delivery_handover(tenant_id,supersedes_handover_id) where supersedes_handover_id is not null;

create table sales.delivery_exception (
  tenant_id uuid not null,
  exception_id text not null,
  organization_id text not null,
  handover_id text not null,
  customer_principal_id text not null,
  reason_code text not null,
  details text not null check (length(details) between 1 and 1000),
  rejection_evidence_sha256_hex text not null check (rejection_evidence_sha256_hex ~ '^[0-9a-f]{64}$'),
  state text not null check (state in ('open','resolved')),
  version bigint not null check (version>0),
  created_at timestamptz not null default clock_timestamp(),
  resolved_at timestamptz,
  resolved_by_subject text,
  primary key (tenant_id,exception_id),
  unique (tenant_id,handover_id),
  foreign key (tenant_id,handover_id) references sales.delivery_handover(tenant_id,handover_id),
  foreign key (tenant_id,organization_id) references org.organization(tenant_id,organization_id),
  foreign key (tenant_id,customer_principal_id) references crm.customer_profile(tenant_id,customer_principal_id),
  check (reason_code ~ '^[a-z][a-z0-9]*(-[a-z0-9]+)*$'),
  check ((state='open' and resolved_at is null and resolved_by_subject is null) or (state='resolved' and resolved_at is not null and length(resolved_by_subject) between 1 and 255))
);

create table sales.return_authorization (
  tenant_id uuid not null,
  authorization_id text not null,
  organization_id text not null,
  exception_id text not null,
  handover_id text not null,
  order_id text not null,
  stock_unit_id text not null,
  customer_principal_id text not null,
  disposition text not null check (disposition in ('return','exchange')),
  exact_cost_source_order_id text not null,
  state text not null check (state='authorized'),
  authorized_by_subject text not null check (length(authorized_by_subject) between 1 and 255),
  authorized_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,authorization_id),
  unique (tenant_id,exception_id),
  foreign key (tenant_id,exception_id) references sales.delivery_exception(tenant_id,exception_id),
  foreign key (tenant_id,handover_id) references sales.delivery_handover(tenant_id,handover_id),
  foreign key (tenant_id,order_id) references sales.customer_order(tenant_id,order_id),
  foreign key (tenant_id,stock_unit_id) references inventory.stock_unit(tenant_id,stock_unit_id),
  foreign key (tenant_id,customer_principal_id) references crm.customer_profile(tenant_id,customer_principal_id),
  check (exact_cost_source_order_id=order_id)
);

create table sales.delivery_exception_resolution (
  tenant_id uuid not null,
  resolution_id text not null,
  exception_id text not null,
  action text not null check (action in ('correct-and-represent','return','exchange')),
  notes text not null check (length(notes) between 1 and 1000),
  resolved_by_subject text not null check (length(resolved_by_subject) between 1 and 255),
  resolved_at timestamptz not null default clock_timestamp(),
  successor_handover_id text,
  return_authorization_id text,
  primary key (tenant_id,resolution_id),
  unique (tenant_id,exception_id),
  foreign key (tenant_id,exception_id) references sales.delivery_exception(tenant_id,exception_id),
  foreign key (tenant_id,successor_handover_id) references sales.delivery_handover(tenant_id,handover_id),
  foreign key (tenant_id,return_authorization_id) references sales.return_authorization(tenant_id,authorization_id),
  check ((action='correct-and-represent' and successor_handover_id is not null and return_authorization_id is null) or (action in ('return','exchange') and successor_handover_id is null and return_authorization_id is not null))
);

create or replace function sales.enforce_delivery_exception_lifecycle() returns trigger language plpgsql as $$
begin
  if tg_op='DELETE' then raise exception using errcode='23514',message='delivery exception is immutable'; end if;
  if old.state<>'open' or new.state<>'resolved' or (new.tenant_id,new.exception_id,new.organization_id,new.handover_id,new.customer_principal_id,new.reason_code,new.details,new.rejection_evidence_sha256_hex,new.created_at) is distinct from (old.tenant_id,old.exception_id,old.organization_id,old.handover_id,old.customer_principal_id,old.reason_code,old.details,old.rejection_evidence_sha256_hex,old.created_at) or new.version<>old.version+1 then
    raise exception using errcode='23514',message='invalid delivery exception transition';
  end if;
  return new;
end $$;
create trigger delivery_exception_lifecycle before update or delete on sales.delivery_exception for each row execute function sales.enforce_delivery_exception_lifecycle();

create or replace function sales.prevent_delivery_resolution_mutation() returns trigger language plpgsql as $$ begin raise exception using errcode='23514',message='delivery resolution is immutable'; end $$;
create trigger delivery_exception_resolution_immutable before update or delete on sales.delivery_exception_resolution for each row execute function sales.prevent_delivery_resolution_mutation();
create trigger return_authorization_immutable before update or delete on sales.return_authorization for each row execute function sales.prevent_delivery_resolution_mutation();

create index delivery_exception_org_state_idx on sales.delivery_exception(tenant_id,organization_id,state,created_at,exception_id);
create index delivery_exception_customer_idx on sales.delivery_exception(tenant_id,customer_principal_id,created_at desc,exception_id);

commit;
