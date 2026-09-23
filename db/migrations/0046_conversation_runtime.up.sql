begin;

create table communication.conversation_turn (
  tenant_id uuid not null references platform.tenant(tenant_id),
  channel_code text not null check (channel_code ~ '^[a-z][a-z0-9_]{0,31}$'),
  provider_message_id text not null check (length(provider_message_id) between 1 and 256),
  external_id text not null check (length(external_id) between 1 and 256),
  thread_id text,
  occurred_at timestamptz not null,
  request_sha256_hex text not null check (request_sha256_hex ~ '^[0-9a-f]{64}$'),
  state text not null check (state in ('processing','retryable','completed','handed_off','failed_terminal')),
  attempt_count integer not null check (attempt_count between 1 and 100),
  locked_until timestamptz,
  user_text text,
  assistant_text text,
  response_sha256_hex text check (response_sha256_hex is null or response_sha256_hex ~ '^[0-9a-f]{64}$'),
  llm_response_id text,
  tool_name text,
  tool_arguments jsonb,
  tool_result text,
  input_tokens bigint not null default 0 check (input_tokens >= 0),
  output_tokens bigint not null default 0 check (output_tokens >= 0),
  reserved_tokens bigint not null default 0 check (reserved_tokens >= 0),
  failure_code text,
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  completed_at timestamptz,
  expires_at timestamptz not null,
  primary key (tenant_id, channel_code, provider_message_id),
  check (thread_id is null or length(thread_id) between 1 and 256),
  check (assistant_text is null or length(assistant_text) between 1 and 65536),
  check (user_text is null or length(user_text) between 1 and 65536),
  check (failure_code is null or failure_code ~ '^[A-Z][A-Z0-9_]{0,127}$'),
  check ((state in ('processing','retryable') and completed_at is null) or
         (state in ('completed','handed_off','failed_terminal') and completed_at is not null)),
  check ((state='processing' and locked_until is not null) or state<>'processing'),
  check ((state in ('completed','handed_off') and assistant_text is not null and response_sha256_hex is not null) or state not in ('completed','handed_off'))
);

create index conversation_turn_history_idx on communication.conversation_turn
  (tenant_id, channel_code, external_id, thread_id, completed_at desc)
  where state='completed';

create index conversation_turn_retry_idx on communication.conversation_turn
  (state, locked_until, updated_at)
  where state in ('processing','retryable');

create function communication.guard_conversation_turn_identity()
returns trigger language plpgsql as $function$
begin
  if new.tenant_id<>old.tenant_id
     or new.channel_code<>old.channel_code
     or new.provider_message_id<>old.provider_message_id
     or new.external_id<>old.external_id
     or new.thread_id is distinct from old.thread_id
     or new.occurred_at<>old.occurred_at
     or new.request_sha256_hex<>old.request_sha256_hex
     or new.created_at<>old.created_at then
    raise exception using errcode='55000', message='conversation turn identity is immutable';
  end if;
  if old.state in ('completed','handed_off','failed_terminal') then
    raise exception using errcode='55000', message='terminal conversation turn is immutable';
  end if;
  return new;
end;
$function$;

create trigger conversation_turn_identity_guard before update on communication.conversation_turn
for each row execute function communication.guard_conversation_turn_identity();

commit;
