create table fiscal.parameter_snapshot (
  tenant_id uuid not null,
  snapshot_id text not null,
  organization_id text not null,
  taxpayer_cuit text not null check (taxpayer_cuit ~ '^[0-9]{11}$'),
  parameter_kind text not null check (parameter_kind in ('voucher_type','concept','document_type','vat_rate','other_tax','point_of_sale','recipient_vat_condition')),
  voucher_class text,
  response_hash text not null check (response_hash ~ '^[a-f0-9]{64}$'),
  provider_codes jsonb not null check (jsonb_typeof(provider_codes)='array'),
  fetched_at timestamptz not null,
  recorded_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,snapshot_id),
  foreign key (tenant_id,organization_id) references org.organization(tenant_id,organization_id),
  check ((parameter_kind='recipient_vat_condition' and voucher_class in ('A','B','C','M')) or (parameter_kind<>'recipient_vat_condition' and voucher_class is null)),
  unique nulls not distinct (tenant_id,organization_id,taxpayer_cuit,parameter_kind,voucher_class,response_hash)
);

create table fiscal.parameter_item (
  tenant_id uuid not null,
  snapshot_id text not null,
  parameter_code text not null check (length(parameter_code) between 1 and 16),
  description text,
  valid_from date,
  valid_until date,
  voucher_class text,
  emission_type text,
  blocked text,
  deregistered_on date,
  primary key (tenant_id,snapshot_id,parameter_code),
  foreign key (tenant_id,snapshot_id) references fiscal.parameter_snapshot(tenant_id,snapshot_id),
  check (valid_until is null or valid_from is null or valid_until>=valid_from)
);

create table fiscal.parameter_decision (
  tenant_id uuid not null,
  snapshot_id text not null,
  decision text not null check (decision in ('approved','rejected')),
  decided_by_subject text not null check (length(decided_by_subject) between 1 and 255),
  decided_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,snapshot_id),
  foreign key (tenant_id,snapshot_id) references fiscal.parameter_snapshot(tenant_id,snapshot_id)
);

create or replace function fiscal.prevent_parameter_history_mutation() returns trigger language plpgsql as $$ begin raise exception 'immutable fiscal parameter history'; end; $$;
create trigger fiscal_parameter_snapshot_immutable before update or delete on fiscal.parameter_snapshot for each row execute function fiscal.prevent_parameter_history_mutation();
create trigger fiscal_parameter_item_immutable before update or delete on fiscal.parameter_item for each row execute function fiscal.prevent_parameter_history_mutation();
create trigger fiscal_parameter_decision_immutable before update or delete on fiscal.parameter_decision for each row execute function fiscal.prevent_parameter_history_mutation();

create view fiscal.approved_parameter_snapshot as
select s.* from fiscal.parameter_snapshot s join fiscal.parameter_decision d using(tenant_id,snapshot_id) where d.decision='approved';
