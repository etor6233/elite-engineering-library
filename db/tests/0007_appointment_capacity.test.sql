begin;
do $$
begin
  if not exists (select 1 from information_schema.tables where table_schema='crm' and table_name='appointment_slot') then raise exception 'appointment_slot missing'; end if;
  if not exists (select 1 from pg_constraint where conname='appointment_slot_time_check') then raise exception 'appointment slot binding missing'; end if;
  if not exists (select 1 from pg_indexes where schemaname='crm' and indexname='appointment_slot_public_idx') then raise exception 'appointment slot public index missing'; end if;
  if to_regprocedure('crm.enforce_appointment_slot_non_overlap()') is null then raise exception 'appointment slot overlap guard missing'; end if;
end $$;
rollback;
