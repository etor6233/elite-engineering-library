begin;

create schema if not exists approval;

create table approval.request (
  tenant_id uuid not null,
  request_id text not null,
  kind text not null check (kind in ('reservation', 'sale', 'refund', 'payment')),
  subject_id text not null,
  amount_minor_units bigint not null check (amount_minor_units >= 0),
  requester text not null,
  evidence_sha text not null check (evidence_sha ~ '^[0-9a-f]{64}$'),
  state text not null default 'pending' check (state in ('pending', 'approved', 'rejected')),
  created_at timestamptz not null default clock_timestamp(),
  decided_at timestamptz,
  primary key (tenant_id, request_id),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  check (length(request_id) between 1 and 128),
  check (length(subject_id) between 1 and 128),
  check (length(requester) between 1 and 128),
  check ((state <> 'pending') = (decided_at is not null))
);

create table approval.decision (
  tenant_id uuid not null,
  request_id text not null,
  reviewer text not null,
  approved boolean not null,
  reason text not null default '',
  decided_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, request_id, reviewer, decided_at),
  foreign key (tenant_id, request_id) references approval.request (tenant_id, request_id),
  check (length(reviewer) between 1 and 128)
);

create index approval_request_pending_idx
  on approval.request (tenant_id, subject_id) where state = 'pending';

commit;
