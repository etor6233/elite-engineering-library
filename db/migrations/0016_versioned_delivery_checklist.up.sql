begin;

create table sales.delivery_checklist_template (
  tenant_id uuid not null,
  organization_id text not null,
  checklist_id text not null,
  checklist_version bigint not null check (checklist_version > 0),
  title text not null check (length(title) between 1 and 160),
  state text not null check (state in ('draft','published','retired')),
  created_by_subject text not null check (length(created_by_subject) between 1 and 255),
  created_at timestamptz not null default clock_timestamp(),
  published_at timestamptz,
  primary key (tenant_id,organization_id,checklist_id,checklist_version),
  foreign key (tenant_id,organization_id) references org.organization(tenant_id,organization_id),
  check (checklist_id ~ '^[a-z][a-z0-9]*(-[a-z0-9]+)*$'),
  check ((state='draft' and published_at is null) or (state in ('published','retired') and published_at is not null))
);

create table sales.delivery_checklist_item (
  tenant_id uuid not null,
  organization_id text not null,
  checklist_id text not null,
  checklist_version bigint not null,
  item_id text not null,
  ordinal integer not null check (ordinal between 1 and 1000),
  prompt text not null check (length(prompt) between 1 and 500),
  response_type text not null check (response_type in ('confirmation','text','serial','evidence')),
  required boolean not null,
  primary key (tenant_id,organization_id,checklist_id,checklist_version,item_id),
  unique (tenant_id,organization_id,checklist_id,checklist_version,ordinal),
  foreign key (tenant_id,organization_id,checklist_id,checklist_version) references sales.delivery_checklist_template(tenant_id,organization_id,checklist_id,checklist_version),
  check (item_id ~ '^[a-z][a-z0-9]*(-[a-z0-9]+)*$')
);

alter table sales.delivery_handover
  add column checklist_id text,
  add column checklist_version bigint,
  add column checklist_completed_at timestamptz,
  add column checklist_completed_by_subject text,
  add constraint delivery_handover_checklist_fk foreign key (tenant_id,organization_id,checklist_id,checklist_version) references sales.delivery_checklist_template(tenant_id,organization_id,checklist_id,checklist_version),
  add constraint delivery_handover_checklist_state_ck check (
    (checklist_id is null and checklist_version is null and checklist_completed_at is null and checklist_completed_by_subject is null)
    or
    (checklist_id is not null and checklist_version > 0 and checklist_completed_at is not null and length(checklist_completed_by_subject) between 1 and 255)
  );

create table sales.delivery_checklist_response (
  tenant_id uuid not null,
  organization_id text not null,
  handover_id text not null,
  checklist_id text not null,
  checklist_version bigint not null,
  item_id text not null,
  response_text text not null check (length(response_text) between 1 and 2048),
  evidence_sha256_hex text,
  answered_by_subject text not null check (length(answered_by_subject) between 1 and 255),
  answered_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,handover_id,item_id),
  foreign key (tenant_id,handover_id) references sales.delivery_handover(tenant_id,handover_id),
  foreign key (tenant_id,organization_id,checklist_id,checklist_version,item_id) references sales.delivery_checklist_item(tenant_id,organization_id,checklist_id,checklist_version,item_id),
  check (evidence_sha256_hex is null or evidence_sha256_hex ~ '^[0-9a-f]{64}$')
);

create or replace function sales.prevent_published_delivery_checklist_mutation() returns trigger language plpgsql as $$
declare parent_state text;
begin
  if tg_table_name='delivery_checklist_template' then
    if old.state in ('published','retired') then raise exception using errcode='23514',message='published delivery checklist is immutable'; end if;
    if tg_op='DELETE' then return old; end if;
    return new;
  end if;
  select state into parent_state from sales.delivery_checklist_template
    where tenant_id=coalesce(new.tenant_id,old.tenant_id) and organization_id=coalesce(new.organization_id,old.organization_id)
      and checklist_id=coalesce(new.checklist_id,old.checklist_id) and checklist_version=coalesce(new.checklist_version,old.checklist_version);
  if parent_state in ('published','retired') then raise exception using errcode='23514',message='published delivery checklist items are immutable'; end if;
  if tg_op='DELETE' then return old; end if;
  return new;
end $$;
create trigger delivery_checklist_template_immutable before update or delete on sales.delivery_checklist_template for each row execute function sales.prevent_published_delivery_checklist_mutation();
create trigger delivery_checklist_item_immutable before insert or update or delete on sales.delivery_checklist_item for each row execute function sales.prevent_published_delivery_checklist_mutation();

create or replace function sales.enforce_delivery_checklist_response() returns trigger language plpgsql as $$
declare expected_type text;
declare completed_at timestamptz;
begin
  select h.checklist_completed_at into completed_at from sales.delivery_handover h where h.tenant_id=new.tenant_id and h.handover_id=new.handover_id for update;
  if completed_at is not null then raise exception using errcode='23514',message='completed delivery checklist responses are immutable'; end if;
  select i.response_type into expected_type from sales.delivery_checklist_item i
    where i.tenant_id=new.tenant_id and i.organization_id=new.organization_id and i.checklist_id=new.checklist_id and i.checklist_version=new.checklist_version and i.item_id=new.item_id;
  if expected_type='confirmation' and new.response_text<>'confirmed' then raise exception using errcode='23514',message='confirmation response must be confirmed'; end if;
  if expected_type='evidence' and new.evidence_sha256_hex is null then raise exception using errcode='23514',message='evidence response requires digest'; end if;
  return new;
end $$;
create trigger delivery_checklist_response_guard before insert on sales.delivery_checklist_response for each row execute function sales.enforce_delivery_checklist_response();

create or replace function sales.prevent_delivery_checklist_response_mutation() returns trigger language plpgsql as $$ begin raise exception using errcode='23514',message='delivery checklist response is immutable'; end $$;
create trigger delivery_checklist_response_immutable before update or delete on sales.delivery_checklist_response for each row execute function sales.prevent_delivery_checklist_response_mutation();

create or replace function sales.enforce_delivery_checklist_completion() returns trigger language plpgsql as $$
begin
  if new.checklist_completed_at is not null and old.checklist_completed_at is null then
    if new.state<>'presented' then raise exception using errcode='23514',message='completed delivery checklist must be presented'; end if;
    if not exists (select 1 from sales.delivery_checklist_template t where t.tenant_id=new.tenant_id and t.organization_id=new.organization_id and t.checklist_id=new.checklist_id and t.checklist_version=new.checklist_version and t.state='published') then
      raise exception using errcode='23514',message='delivery checklist version is not published';
    end if;
    if exists (
      select 1 from sales.delivery_checklist_item i
      left join sales.delivery_checklist_response r on r.tenant_id=i.tenant_id and r.organization_id=i.organization_id and r.checklist_id=i.checklist_id and r.checklist_version=i.checklist_version and r.item_id=i.item_id and r.handover_id=new.handover_id
      where i.tenant_id=new.tenant_id and i.organization_id=new.organization_id and i.checklist_id=new.checklist_id and i.checklist_version=new.checklist_version and i.required and r.item_id is null
    ) then raise exception using errcode='23514',message='required delivery checklist response is missing'; end if;
  elsif old.checklist_completed_at is not null and (new.checklist_id,new.checklist_version,new.checklist_completed_at,new.checklist_completed_by_subject) is distinct from (old.checklist_id,old.checklist_version,old.checklist_completed_at,old.checklist_completed_by_subject) then
    raise exception using errcode='23514',message='completed delivery checklist binding is immutable';
  end if;
  return new;
end $$;
create trigger delivery_checklist_completion_guard before update of state,checklist_id,checklist_version,checklist_completed_at,checklist_completed_by_subject on sales.delivery_handover for each row execute function sales.enforce_delivery_checklist_completion();

create index delivery_checklist_template_state_idx on sales.delivery_checklist_template(tenant_id,organization_id,state,checklist_id,checklist_version desc);
create index delivery_checklist_response_handover_idx on sales.delivery_checklist_response(tenant_id,handover_id,answered_at,item_id);

commit;
