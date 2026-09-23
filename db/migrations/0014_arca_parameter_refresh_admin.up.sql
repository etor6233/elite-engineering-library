alter table fiscal.parameter_decision add column decision_reason text check (decision_reason is null or length(decision_reason) between 3 and 500);
create or replace function fiscal.require_parameter_decision_reason() returns trigger language plpgsql as $$ begin if new.decision_reason is null then raise exception 'decision reason required'; end if; return new; end; $$;
create trigger fiscal_parameter_decision_reason before insert on fiscal.parameter_decision for each row execute function fiscal.require_parameter_decision_reason();

create table fiscal.parameter_refresh_schedule (
  tenant_id uuid not null,
  schedule_id text not null,
  organization_id text not null,
  taxpayer_cuit text not null check (taxpayer_cuit ~ '^[0-9]{11}$'),
  parameter_kind text not null check (parameter_kind in ('voucher_type','concept','document_type','vat_rate','other_tax','point_of_sale','recipient_vat_condition')),
  voucher_class text,
  interval_seconds bigint not null check (interval_seconds between 900 and 2592000),
  next_run_at timestamptz not null,
  active boolean not null,
  attempts bigint not null default 0,
  lease_owner text,
  lease_until timestamptz,
  last_snapshot_id text,
  last_succeeded_at timestamptz,
  last_error_code text,
  version bigint not null check (version>0),
  configured_by_subject text not null check (length(configured_by_subject) between 1 and 255),
  configured_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,schedule_id),
  foreign key (tenant_id,organization_id) references org.organization(tenant_id,organization_id),
  foreign key (tenant_id,last_snapshot_id) references fiscal.parameter_snapshot(tenant_id,snapshot_id),
  check ((parameter_kind='recipient_vat_condition' and voucher_class in ('A','B','C','M')) or (parameter_kind<>'recipient_vat_condition' and voucher_class is null)),
  check ((lease_owner is null)=(lease_until is null)),
  unique nulls not distinct (tenant_id,organization_id,taxpayer_cuit,parameter_kind,voucher_class)
);

create index fiscal_parameter_refresh_due_idx on fiscal.parameter_refresh_schedule(next_run_at) where active;
