create schema if not exists accounting;

create table accounting.account (
  tenant_id uuid not null,
  account_code text not null check (account_code ~ '^[A-Z0-9][A-Z0-9._-]{0,31}$'),
  display_name text not null check (length(display_name) between 2 and 120),
  account_type text not null check (account_type in ('asset','liability','equity','revenue','expense')),
  active boolean not null default true,
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, account_code),
  foreign key (tenant_id) references platform.tenant (tenant_id)
);

create table accounting.period (
  tenant_id uuid not null,
  period_id text not null,
  starts_on date not null,
  ends_on date not null,
  status text not null check (status in ('open','closed')),
  version bigint not null check (version > 0),
  closed_at timestamptz,
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, period_id),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  unique (tenant_id, starts_on, ends_on),
  check (ends_on > starts_on),
  check ((status='open' and closed_at is null) or (status='closed' and closed_at is not null))
);

create table accounting.journal (
  tenant_id uuid not null,
  journal_id text not null,
  organization_id text not null,
  period_id text not null,
  source_type text not null check (source_type ~ '^[A-Z0-9][A-Z0-9._-]{0,31}$'),
  source_id text not null,
  currency text not null check (currency ~ '^[A-Z]{3}$'),
  posting_date date not null,
  status text not null check (status in ('draft','posted','reversed')),
  total_debit_minor_units bigint not null check (total_debit_minor_units > 0),
  total_credit_minor_units bigint not null check (total_credit_minor_units > 0),
  version bigint not null check (version > 0),
  reversal_of text,
  reversal_reason text,
  posted_by text,
  posted_at timestamptz,
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, journal_id),
  foreign key (tenant_id, organization_id) references org.organization (tenant_id, organization_id),
  foreign key (tenant_id, period_id) references accounting.period (tenant_id, period_id),
  foreign key (tenant_id, reversal_of) references accounting.journal (tenant_id, journal_id),
  unique (tenant_id, source_type, source_id),
  check (total_debit_minor_units = total_credit_minor_units),
  check ((status='draft' and posted_at is null and posted_by is null) or (status in ('posted','reversed') and posted_at is not null and posted_by is not null)),
  check ((reversal_of is null and reversal_reason is null) or (reversal_of is not null and length(reversal_reason) between 3 and 500))
);

create unique index accounting_one_reversal_idx on accounting.journal(tenant_id,reversal_of) where reversal_of is not null;

create table accounting.journal_line (
  tenant_id uuid not null,
  journal_id text not null,
  line_no integer not null check (line_no > 0),
  account_code text not null,
  description text not null default '' check (length(description) <= 250),
  debit_minor_units bigint not null default 0 check (debit_minor_units >= 0),
  credit_minor_units bigint not null default 0 check (credit_minor_units >= 0),
  primary key (tenant_id,journal_id,line_no),
  foreign key (tenant_id,journal_id) references accounting.journal(tenant_id,journal_id),
  foreign key (tenant_id,account_code) references accounting.account(tenant_id,account_code),
  check ((debit_minor_units > 0) <> (credit_minor_units > 0))
);

create table accounting.register (
  tenant_id uuid not null,
  register_id text not null,
  journal_id text not null,
  organization_id text not null,
  period_id text not null,
  posted_by text not null,
  total_debit_minor_units bigint not null check (total_debit_minor_units > 0),
  total_credit_minor_units bigint not null check (total_credit_minor_units > 0),
  posted_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,register_id),
  foreign key (tenant_id,journal_id) references accounting.journal(tenant_id,journal_id),
  unique (tenant_id,journal_id),
  check (total_debit_minor_units=total_credit_minor_units)
);

create table accounting.entry (
  tenant_id uuid not null,
  entry_id text not null,
  register_id text not null,
  journal_id text not null,
  line_no integer not null,
  organization_id text not null,
  period_id text not null,
  account_code text not null,
  currency text not null check (currency ~ '^[A-Z]{3}$'),
  posting_date date not null,
  description text not null,
  debit_minor_units bigint not null check (debit_minor_units >= 0),
  credit_minor_units bigint not null check (credit_minor_units >= 0),
  original_entry_id text,
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,entry_id),
  foreign key (tenant_id,register_id) references accounting.register(tenant_id,register_id),
  foreign key (tenant_id,journal_id,line_no) references accounting.journal_line(tenant_id,journal_id,line_no),
  foreign key (tenant_id,account_code) references accounting.account(tenant_id,account_code),
  unique (tenant_id,journal_id,line_no),
  check ((debit_minor_units > 0) <> (credit_minor_units > 0))
);

create unique index accounting_entry_one_reversal_idx on accounting.entry(tenant_id,original_entry_id) where original_entry_id is not null;

create index accounting_trial_balance_idx on accounting.entry(tenant_id,organization_id,period_id,account_code);

create or replace function accounting.prevent_posted_history_mutation() returns trigger language plpgsql as $$ begin raise exception 'immutable accounting history'; end; $$;
create trigger accounting_line_immutable before update or delete on accounting.journal_line for each row execute function accounting.prevent_posted_history_mutation();
create trigger accounting_register_immutable before update or delete on accounting.register for each row execute function accounting.prevent_posted_history_mutation();
create trigger accounting_entry_immutable before update or delete on accounting.entry for each row execute function accounting.prevent_posted_history_mutation();
