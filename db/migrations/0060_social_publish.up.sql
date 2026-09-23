begin;
-- Provider references and receipts remain public outbox evidence. Protect the
-- fields which authorize GET-only recovery while allowing normal outbox leases.
create function communication.guard_social_observation() returns trigger language plpgsql as $$
begin
 if tg_op='DELETE' then
  if old.aggregate_type='social_approval' then raise exception 'social evidence is immutable';end if;return old;
 end if;
 if (old.aggregate_type='social_approval' or new.aggregate_type='social_approval') and
 row(new.tenant_id,new.event_id,new.aggregate_type,new.aggregate_id,new.aggregate_version,new.event_type,new.schema_version,new.occurred_at,new.payload)
 is distinct from row(old.tenant_id,old.event_id,old.aggregate_type,old.aggregate_id,old.aggregate_version,old.event_type,old.schema_version,old.occurred_at,old.payload)
 then raise exception 'social evidence is immutable';end if;return new;
end $$;
create trigger social_observation_immutable before update or delete on platform.outbox_event
 for each row execute function communication.guard_social_observation();
commit;
