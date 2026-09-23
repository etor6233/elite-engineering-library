create schema if not exists royalty;

create table royalty.policy (
  tenant_id uuid not null,
  policy_id text not null,
  agreement_id text not null,
  franchise_organization_id text not null,
  currency text not null check (currency ~ '^[A-Z]{3}$'),
  rate_basis_points integer not null check (rate_basis_points between 1 and 10000),
  basis text not null default 'captured-payment' check (basis = 'captured-payment'),
  valid_from timestamptz not null,
  valid_until timestamptz,
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, policy_id),
  foreign key (tenant_id, agreement_id) references franchise.agreement (tenant_id, agreement_id),
  foreign key (tenant_id, franchise_organization_id) references org.organization (tenant_id, organization_id),
  check (valid_until is null or valid_until > valid_from)
);

create index royalty_policy_resolution_idx
  on royalty.policy (tenant_id, franchise_organization_id, currency, valid_from desc);

create table royalty.accrual (
  tenant_id uuid not null,
  accrual_id text not null,
  policy_id text not null,
  agreement_id text not null,
  franchise_organization_id text not null,
  payment_attempt_id text not null,
  source_event_key text not null,
  source_state text not null check (source_state in ('captured','refunded')),
  currency text not null check (currency ~ '^[A-Z]{3}$'),
  basis_minor_units bigint not null check (basis_minor_units > 0),
  royalty_minor_units bigint not null check (royalty_minor_units <> 0),
  occurred_at timestamptz not null,
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, accrual_id),
  foreign key (tenant_id, policy_id) references royalty.policy (tenant_id, policy_id),
  foreign key (tenant_id, agreement_id) references franchise.agreement (tenant_id, agreement_id),
  foreign key (tenant_id, franchise_organization_id) references org.organization (tenant_id, organization_id),
  foreign key (tenant_id, payment_attempt_id) references payment.payment_attempt (tenant_id, payment_attempt_id),
  unique (tenant_id, source_event_key),
  unique (tenant_id, payment_attempt_id, source_state),
  check ((source_state = 'captured' and royalty_minor_units > 0) or (source_state = 'refunded' and royalty_minor_units < 0))
);

create index royalty_accrual_unsettled_idx
  on royalty.accrual (tenant_id, franchise_organization_id, currency, occurred_at, accrual_id);

create table royalty.settlement_run (
  tenant_id uuid not null,
  settlement_id text not null,
  franchise_organization_id text not null,
  currency text not null check (currency ~ '^[A-Z]{3}$'),
  period_start timestamptz not null,
  period_end timestamptz not null,
  status text not null check (status in ('draft','closed','reversed')),
  expected_minor_units bigint not null default 0,
  version bigint not null check (version > 0),
  reversal_of text,
  reversal_reason text,
  created_at timestamptz not null default clock_timestamp(),
  closed_at timestamptz,
  primary key (tenant_id, settlement_id),
  foreign key (tenant_id, franchise_organization_id) references org.organization (tenant_id, organization_id),
  foreign key (tenant_id, reversal_of) references royalty.settlement_run (tenant_id, settlement_id),
  check (period_end > period_start),
  check ((status = 'draft' and closed_at is null) or (status in ('closed','reversed') and closed_at is not null)),
  check ((reversal_of is null and reversal_reason is null) or (reversal_of is not null and length(reversal_reason) between 3 and 500))
);

create unique index royalty_settlement_one_reversal_idx
  on royalty.settlement_run (tenant_id, reversal_of)
  where reversal_of is not null;

create unique index royalty_settlement_one_draft_period_idx
  on royalty.settlement_run (tenant_id, franchise_organization_id, currency, period_start, period_end)
  where reversal_of is null and status = 'draft';

create table royalty.settlement_line (
  tenant_id uuid not null,
  settlement_id text not null,
  line_id text not null,
  accrual_id text not null,
  amount_minor_units bigint not null check (amount_minor_units <> 0),
  original_line_id text,
  primary key (tenant_id, settlement_id, line_id),
  foreign key (tenant_id, settlement_id) references royalty.settlement_run (tenant_id, settlement_id),
  foreign key (tenant_id, accrual_id) references royalty.accrual (tenant_id, accrual_id)
);

create unique index royalty_settlement_line_once_idx
  on royalty.settlement_line (tenant_id, accrual_id)
  where original_line_id is null;

create unique index royalty_settlement_reversal_once_idx
  on royalty.settlement_line (tenant_id, original_line_id)
  where original_line_id is not null;

create table royalty.reconciliation (
  tenant_id uuid not null,
  reconciliation_id text not null,
  settlement_id text not null,
  external_reference text not null,
  expected_minor_units bigint not null,
  actual_minor_units bigint not null,
  difference_minor_units bigint not null,
  status text not null check (status in ('matched','mismatch')),
  recorded_by text not null,
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, reconciliation_id),
  foreign key (tenant_id, settlement_id) references royalty.settlement_run (tenant_id, settlement_id),
  unique (tenant_id, settlement_id),
  unique (tenant_id, external_reference),
  check (difference_minor_units = actual_minor_units - expected_minor_units),
  check ((status = 'matched') = (difference_minor_units = 0))
);

create or replace function royalty.prevent_financial_history_mutation()
returns trigger language plpgsql as $$
begin
  raise exception 'immutable royalty history';
end;
$$;

create trigger royalty_policy_immutable before update or delete on royalty.policy
for each row execute function royalty.prevent_financial_history_mutation();
create trigger royalty_accrual_immutable before update or delete on royalty.accrual
for each row execute function royalty.prevent_financial_history_mutation();
create trigger royalty_line_immutable before update or delete on royalty.settlement_line
for each row execute function royalty.prevent_financial_history_mutation();
create trigger royalty_reconciliation_immutable before update or delete on royalty.reconciliation
for each row execute function royalty.prevent_financial_history_mutation();
