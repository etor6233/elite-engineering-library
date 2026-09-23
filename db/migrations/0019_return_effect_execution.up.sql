begin;

create table sales.return_effect_execution (
  tenant_id uuid not null,
  request_id text not null,
  status text not null check (status in ('requested','claimed','retry','blocked','failed','succeeded')),
  attempt_count integer not null default 0 check (attempt_count >= 0),
  available_at timestamptz not null default clock_timestamp(),
  claimed_at timestamptz,
  claimed_by text,
  claim_token text,
  claimed_until timestamptz,
  last_error_code text,
  provider_reference text,
  result_sha256_hex text,
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,request_id),
  foreign key (tenant_id,request_id) references sales.return_effect_request(tenant_id,request_id),
  unique (tenant_id,claim_token),
  check ((status='claimed')=(claimed_at is not null and claimed_by is not null and claim_token is not null and claimed_until is not null)),
  check (status<>'claimed' or claimed_until>claimed_at),
  check (result_sha256_hex is null or result_sha256_hex ~ '^[0-9a-f]{64}$'),
  check (last_error_code is null or last_error_code ~ '^[A-Z][A-Z0-9_]{1,63}$'),
  check (status<>'succeeded' or (result_sha256_hex is not null and last_error_code is null)),
  check (status not in ('retry','blocked','failed') or last_error_code is not null)
);

create table sales.return_effect_attempt (
  tenant_id uuid not null,
  attempt_id text not null,
  request_id text not null,
  attempt_no integer not null check (attempt_no>0),
  worker_id text not null check (worker_id ~ '^[A-Za-z0-9][A-Za-z0-9._:-]{0,63}$'),
  outcome text not null check (outcome in ('retry','blocked','failed','succeeded')),
  error_code text,
  provider_reference text,
  result_sha256_hex text,
  started_at timestamptz not null,
  completed_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,attempt_id),
  unique (tenant_id,request_id,attempt_no),
  foreign key (tenant_id,request_id) references sales.return_effect_request(tenant_id,request_id),
  check (completed_at>=started_at),
  check (error_code is null or error_code ~ '^[A-Z][A-Z0-9_]{1,63}$'),
  check (result_sha256_hex is null or result_sha256_hex ~ '^[0-9a-f]{64}$'),
  check ((outcome='succeeded' and error_code is null and result_sha256_hex is not null) or (outcome<>'succeeded' and error_code is not null and result_sha256_hex is null))
);

create table sales.return_effect_resume (
  tenant_id uuid not null,
  resume_id text not null,
  request_id text not null,
  reason_code text not null check (reason_code ~ '^[A-Z][A-Z0-9_]{1,63}$'),
  requested_by_subject text not null check (length(requested_by_subject) between 1 and 255),
  requested_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,resume_id),
  foreign key (tenant_id,request_id) references sales.return_effect_request(tenant_id,request_id)
);

create or replace function sales.enqueue_return_effect_execution() returns trigger language plpgsql as $$
begin
  insert into sales.return_effect_execution(tenant_id,request_id,status)
  values(new.tenant_id,new.request_id,'requested');
  return new;
end $$;

create trigger return_effect_execution_enqueue
after insert on sales.return_effect_request
for each row execute function sales.enqueue_return_effect_execution();

insert into sales.return_effect_execution(tenant_id,request_id,status)
select tenant_id,request_id,'requested' from sales.return_effect_request
on conflict do nothing;

create or replace function sales.prevent_return_effect_audit_mutation() returns trigger language plpgsql as $$
begin
  raise exception using errcode='23514',message='return effect audit evidence is immutable';
end $$;

create trigger return_effect_attempt_immutable before update or delete on sales.return_effect_attempt for each row execute function sales.prevent_return_effect_audit_mutation();
create trigger return_effect_resume_immutable before update or delete on sales.return_effect_resume for each row execute function sales.prevent_return_effect_audit_mutation();

create index return_effect_execution_claim_idx on sales.return_effect_execution(status,available_at,claimed_until,tenant_id,request_id);
create index return_effect_attempt_request_idx on sales.return_effect_attempt(tenant_id,request_id,attempt_no desc);

commit;
