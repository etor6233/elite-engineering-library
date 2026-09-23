begin;
-- AUTHORED service/approval/inventory binding; upstream owners remain authoritative.
alter table approval.request drop constraint request_kind_check;
alter table approval.request add constraint request_kind_check check(kind in
 ('reservation','sale','refund','payment','whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review'));
alter table approval.request drop constraint request_payload_binding_check;
alter table approval.request add constraint request_payload_binding_check check
 (kind not in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review') or
  (organization_id is not null and length(organization_id) between 1 and 128 and
   payload is not null and jsonb_typeof(payload)='object' and pg_column_size(payload)<=65536 and amount_minor_units=0));
create or replace function approval.guard_bound_request() returns trigger language plpgsql as $$
begin
 if tg_op='DELETE' then
  if old.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review') then raise exception 'bound approval request is immutable'; end if;
  return old;
 end if;
 if old.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review') or new.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review') then
  if row(new.tenant_id,new.request_id,new.kind,new.subject_id,new.amount_minor_units,new.requester,new.evidence_sha,new.organization_id,new.payload,new.created_at)
   is distinct from row(old.tenant_id,old.request_id,old.kind,old.subject_id,old.amount_minor_units,old.requester,old.evidence_sha,old.organization_id,old.payload,old.created_at)
   or (old.state<>'pending' and row(new.state,new.decided_at) is distinct from row(old.state,old.decided_at))
  then raise exception 'bound approval request is immutable'; end if;
 end if;
 return new;
end $$;
create or replace function approval.guard_bound_decision() returns trigger language plpgsql as $$
begin
 if exists(select 1 from approval.request where tenant_id=old.tenant_id and request_id=old.request_id and kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review')) then
  raise exception 'bound approval decision is immutable';
 end if;
 if tg_op='DELETE' then return old;end if;return new;
end $$;

create table catalog.release_stream(
 tenant_id uuid primary key references platform.tenant(tenant_id),
 organization_id text not null,
 market text not null check(market~'^[A-Z]{2}$'),
 currency text not null check(currency~'^[A-Z]{3}$'),
 profile jsonb not null,
 profile_canonical text not null,
 profile_sha256 text not null,
 foreign key(tenant_id,organization_id) references org.organization(tenant_id,organization_id),
 check(profile=profile_canonical::jsonb and profile_sha256=encode(sha256(convert_to(profile_canonical,'UTF8')),'hex'))
);
create table catalog.release_media(
 tenant_id uuid not null references catalog.release_stream(tenant_id),
 media_id text not null, maker text not null, original_sha256 text not null check(original_sha256~'^[0-9a-f]{64}$'),
 sha256 text not null, width integer not null, height integer not null, png bytea not null,
 primary key(tenant_id,media_id),check(width between 1 and 2048 and height between 1 and 2048 and width*height<=1048576),
 check(octet_length(png)<=4194304 and sha256=encode(sha256(png),'hex'))
);
create table catalog.release_draft(
 tenant_id uuid not null references catalog.release_stream(tenant_id),
 draft_id text not null, maker text not null, price_book_id text not null,
 snapshot jsonb not null, snapshot_canonical text not null, snapshot_sha256 text not null,
 created_at timestamptz not null default clock_timestamp(),transaction_id xid8 not null default pg_current_xact_id(),
 primary key(tenant_id,draft_id),foreign key(tenant_id,price_book_id) references pricing.price_book(tenant_id,price_book_id),
 check(octet_length(snapshot_canonical)<=262144 and snapshot=snapshot_canonical::jsonb and snapshot_sha256=encode(sha256(convert_to(snapshot_canonical,'UTF8')),'hex')),
 check(jsonb_typeof(snapshot->'models')='array' and jsonb_array_length(snapshot->'models') between 1 and 32),
 check(jsonb_typeof(snapshot->'variants')='array' and jsonb_array_length(snapshot->'variants') between 1 and 128)
);
create table catalog.release_review(
 tenant_id uuid not null,draft_id text not null,stage text not null check(stage in('legal','technical','media','publication')),
 approval_id text not null,payload_sha256 text not null,
 primary key(tenant_id,draft_id,stage),unique(tenant_id,approval_id),
 foreign key(tenant_id,draft_id) references catalog.release_draft(tenant_id,draft_id),
 foreign key(tenant_id,approval_id) references approval.request(tenant_id,request_id) deferrable initially deferred
);
create table catalog.release_command(
 tenant_id uuid not null references catalog.release_stream(tenant_id),
 command_id text not null,actor text not null,kind text not null check(kind in('media','draft','review','publish')),
 resource_id text not null,request_sha256 text not null check(request_sha256~'^[0-9a-f]{64}$'),
 receipt jsonb not null,transaction_id xid8 not null default pg_current_xact_id(),
 primary key(tenant_id,command_id),check(jsonb_typeof(receipt)='object')
);
create table catalog.release_review_action(
 tenant_id uuid not null,command_id text not null,approval_id text not null,actor text not null,
 approved boolean not null,evidence_sha256 text not null check(evidence_sha256~'^[0-9a-f]{64}$'),
 transaction_id xid8 not null default pg_current_xact_id(),
 primary key(tenant_id,approval_id),unique(tenant_id,command_id),
 foreign key(tenant_id,command_id) references catalog.release_command(tenant_id,command_id) deferrable initially deferred,
 foreign key(tenant_id,approval_id) references catalog.release_review(tenant_id,approval_id)
);
create table catalog.release_publication(
 tenant_id uuid not null,generation bigint not null check(generation>0),draft_id text not null,command_id text not null,
 actor text not null,reason text not null,effective_price_book_id text not null,created_at timestamptz not null default clock_timestamp(),
 transaction_id xid8 not null default pg_current_xact_id(),
 primary key(tenant_id,generation),unique(tenant_id,command_id),
 foreign key(tenant_id,effective_price_book_id) references pricing.price_book(tenant_id,price_book_id),
 foreign key(tenant_id,draft_id) references catalog.release_draft(tenant_id,draft_id),
 foreign key(tenant_id,command_id) references catalog.release_command(tenant_id,command_id) deferrable initially deferred
);
create index catalog_release_draft_book on catalog.release_draft(tenant_id,price_book_id);
create unique index catalog_release_effective_book on catalog.release_publication(tenant_id,effective_price_book_id);
create index catalog_release_media_digest on catalog.release_media(tenant_id,sha256);
create function catalog.release_immutable() returns trigger language plpgsql as $$
begin raise exception 'catalog release facts are immutable';end $$;
do $$ declare n text;begin
 foreach n in array array['release_stream','release_media','release_draft','release_review','release_command','release_review_action','release_publication'] loop
 execute format('create trigger catalog_release_immutable before update or delete on catalog.%I for each row execute function catalog.release_immutable()',n);
 end loop;
end $$;
create view catalog.release_current as
 select p.*,d.snapshot,d.snapshot_canonical,d.snapshot_sha256,d.price_book_id as source_price_book_id
 from catalog.release_publication p join catalog.release_draft d using(tenant_id,draft_id)
 where p.generation=(select max(x.generation) from catalog.release_publication x where x.tenant_id=p.tenant_id);
create view catalog.release_public_current as
 select c.* from catalog.release_current c join pricing.price_book b on b.tenant_id=c.tenant_id and b.price_book_id=c.effective_price_book_id
 join catalog.release_stream st on st.tenant_id=c.tenant_id join platform.tenant tenant on tenant.tenant_id=c.tenant_id join org.organization org on org.tenant_id=st.tenant_id and org.organization_id=st.organization_id
 where tenant.status='active' and org.status='active' and b.status='active' and b.valid_from<=statement_timestamp() and (b.valid_until is null or b.valid_until>statement_timestamp());
create view catalog.release_public_models as
 select c.tenant_id,m->>'id' as model_id,m->>'code' as model_code,m->>'displayName' as display_name,
 m->>'vehicleClass' as vehicle_class,m->'specification' as specification
 from catalog.release_public_current c cross join lateral jsonb_array_elements(c.snapshot->'models') m
 union all
 select tenant_id,model_id,model_code,display_name,vehicle_class,specification from catalog.vehicle_model m
 where m.lifecycle_state='active' and m.publicly_visible and not exists(select 1 from catalog.release_stream s where s.tenant_id=m.tenant_id);
create function catalog.release_guard_price() returns trigger language plpgsql as $$
declare t uuid;bid text;market_ text;currency_ text;begin
 if tg_op='DELETE' then t:=old.tenant_id;bid:=old.price_book_id;else t:=new.tenant_id;bid:=new.price_book_id;end if;
 if tg_table_name='price_book_entry' then
  if (exists(select 1 from catalog.release_draft where tenant_id=t and price_book_id=bid) or exists(select 1 from catalog.release_publication where tenant_id=t and effective_price_book_id=bid)) then raise exception 'captured price entries immutable';end if;
 else
  if tg_op='DELETE' then
   if (exists(select 1 from catalog.release_draft where tenant_id=t and price_book_id=bid) or exists(select 1 from catalog.release_publication where tenant_id=t and effective_price_book_id=bid)) then raise exception 'captured price book immutable';end if;
  elsif tg_op='UPDATE' then
   if (exists(select 1 from catalog.release_draft where tenant_id=old.tenant_id and price_book_id=old.price_book_id) or exists(select 1 from catalog.release_publication where tenant_id=old.tenant_id and effective_price_book_id=old.price_book_id))
   and row(new.tenant_id,new.price_book_id,new.market,new.currency,new.valid_from,new.valid_until)
    is distinct from row(old.tenant_id,old.price_book_id,old.market,old.currency,old.valid_from,old.valid_until)
   then raise exception 'captured price header immutable';end if;
  end if;
  if tg_op<>'DELETE' then
   market_:=new.market;currency_:=new.currency;
   if exists(select 1 from catalog.release_stream where tenant_id=t and market=market_ and currency=currency_)
   and ((tg_op='INSERT' and new.status='active') or (tg_op='UPDATE' and new.status is distinct from old.status))
   and not exists(select 1 from catalog.release_publication where tenant_id=t and transaction_id=pg_current_xact_id())
   then raise exception 'price publication command required';end if;
  end if;
 end if;
 if tg_op='DELETE' then return old;end if;return new;
end $$;
create trigger catalog_release_price_guard before insert or update or delete on pricing.price_book for each row execute function catalog.release_guard_price();
create trigger catalog_release_entry_guard before insert or update or delete on pricing.price_book_entry for each row execute function catalog.release_guard_price();
create function catalog.release_guard_approval() returns trigger language plpgsql as $$
declare binding record;action record;begin
 if new.kind<>'catalog_review' then return new;end if;
 select r.*,d.maker,d.snapshot_sha256,d.transaction_id,s.organization_id into binding
 from catalog.release_review r join catalog.release_draft d using(tenant_id,draft_id)
 join catalog.release_stream s using(tenant_id)
 where r.tenant_id=new.tenant_id and r.approval_id=new.request_id;
 if not found or new.requester<>binding.maker or new.organization_id<>binding.organization_id
  or new.evidence_sha<>binding.payload_sha256 or new.amount_minor_units<>0
  or new.payload->>'draft_id' is distinct from binding.draft_id or new.payload->>'stage' is distinct from binding.stage
  or new.payload->>'snapshot_sha256' is distinct from binding.snapshot_sha256
 then raise exception 'catalog approval binding mismatch';end if;
 if tg_op='INSERT' then
  if binding.transaction_id<>pg_current_xact_id() then raise exception 'catalog draft command required';end if;
 elsif new.state<>old.state then
  select a.* into action from catalog.release_review_action a
  join catalog.release_command c using(tenant_id,command_id)
  join approval.decision v on v.tenant_id=a.tenant_id and v.request_id=a.approval_id
  where a.tenant_id=new.tenant_id and a.approval_id=new.request_id and a.transaction_id=pg_current_xact_id()
   and c.transaction_id=pg_current_xact_id() and c.kind='review' and c.actor=a.actor
   and a.actor=v.reviewer and a.actor<>binding.maker and a.approved=v.approved
   and new.state=case when a.approved then 'approved' else 'rejected' end;
  if not found then raise exception 'catalog review command required';end if;
  if binding.stage='publication' and action.approved and
   (select count(*) from catalog.release_review r join approval.request q on q.tenant_id=r.tenant_id and q.request_id=r.approval_id
    where r.tenant_id=new.tenant_id and r.draft_id=binding.draft_id and r.stage<>'publication' and q.state='approved')<>3
  then raise exception 'prior catalog reviews required';end if;
 end if;
 return new;
end $$;
create constraint trigger catalog_release_approval_guard after insert or update on approval.request deferrable initially deferred for each row execute function catalog.release_guard_approval();
create function catalog.release_guard_publication() returns trigger language plpgsql as $$
declare d record;s record;c record;expected bigint;begin
 select * into d from catalog.release_draft where tenant_id=new.tenant_id and draft_id=new.draft_id;
 select * into s from catalog.release_stream where tenant_id=new.tenant_id;
 select * into c from catalog.release_command where tenant_id=new.tenant_id and command_id=new.command_id;
 select coalesce(max(generation),0)+1 into expected from catalog.release_publication where tenant_id=new.tenant_id and generation<new.generation;
 if new.transaction_id<>pg_current_xact_id() or c.transaction_id<>pg_current_xact_id() or c.kind<>'publish' or c.actor<>new.actor
  or c.resource_id<>new.draft_id or new.generation<>expected or new.actor=d.maker
  or (c.receipt->>'generation')::bigint is distinct from new.generation or c.receipt->>'snapshot_sha256' is distinct from d.snapshot_sha256
 then raise exception 'catalog publication receipt mismatch';end if;
 if (select count(*) from catalog.release_review r join approval.request q on q.tenant_id=r.tenant_id and q.request_id=r.approval_id
  where r.tenant_id=new.tenant_id and r.draft_id=new.draft_id and q.kind='catalog_review' and q.state='approved')<>4
 then raise exception 'catalog publication needs four approved reviews';end if;
 if not exists(select 1 from pricing.price_book b where b.tenant_id=new.tenant_id and b.price_book_id=new.effective_price_book_id
  and b.market=d.snapshot->'profile'->>'market' and b.currency=d.snapshot->'profile'->>'currency'
  and b.valid_from=(d.snapshot->>'valid_from')::timestamptz and b.valid_until is not distinct from (d.snapshot->>'valid_until')::timestamptz)
 or c.receipt->>'effective_price_book_id' is distinct from new.effective_price_book_id
 or (select count(*) from pricing.price_book_entry e where e.tenant_id=new.tenant_id and e.price_book_id=new.effective_price_book_id)<>jsonb_array_length(d.snapshot->'variants')
 or exists(select 1 from jsonb_array_elements(d.snapshot->'variants') v where not exists(
  select 1 from pricing.price_book_entry e where e.tenant_id=new.tenant_id and e.price_book_id=new.effective_price_book_id and e.variant_id=v->>'id'
   and e.amount_minor_units=(v->>'amount_minor_units')::bigint and e.tax_mode=v->>'tax_mode'))
 then raise exception 'effective price book differs from approved snapshot';end if;

 if not exists(select 1 from pricing.price_book b where b.tenant_id=new.tenant_id and b.price_book_id=new.effective_price_book_id and b.status='active'
  and b.valid_from<=clock_timestamp() and (b.valid_until is null or b.valid_until>clock_timestamp()))
 then raise exception 'catalog price expired or not active at commit';end if;
 if exists(select 1 from pricing.price_book b where b.tenant_id=new.tenant_id and b.market=s.market and b.currency=s.currency and b.status='active' and b.price_book_id<>new.effective_price_book_id)
 then raise exception 'catalog competing active price';end if;
 if (select count(*) from search.document x where x.tenant_id=new.tenant_id and x.kind='catalog-model')<>jsonb_array_length(d.snapshot->'models')
 or exists(select 1 from jsonb_array_elements(d.snapshot->'models') m where not exists(
  select 1 from search.document x where x.tenant_id=new.tenant_id and x.document_id='catalog:'||(m->>'id') and x.kind='catalog-model'
   and x.title=m->>'displayName' and x.body=m->>'canonical_url' and x.external_id=m->>'id' and x.content_sha256=d.snapshot_sha256))
 then raise exception 'catalog search projection mismatch';end if;
 return new;
end $$;
create constraint trigger catalog_release_publication_guard after insert on catalog.release_publication deferrable initially deferred for each row execute function catalog.release_guard_publication();
create function catalog.release_guard_search() returns trigger language plpgsql as $$
declare t uuid;begin
 if tg_op='DELETE' then t:=old.tenant_id;else t:=new.tenant_id;end if;
 if ((tg_op<>'INSERT' and old.kind='catalog-model') or (tg_op<>'DELETE' and new.kind='catalog-model'))
 and exists(select 1 from catalog.release_stream where tenant_id=t)
 and not exists(select 1 from catalog.release_publication where tenant_id=t and transaction_id=pg_current_xact_id())
 then raise exception 'catalog projection requires publication';end if;
 if tg_op='DELETE' then return old;end if;return new;
end $$;
create trigger catalog_release_search_guard before insert or update or delete on search.document for each row execute function catalog.release_guard_search();

commit;
