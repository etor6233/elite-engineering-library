begin;
do $$ begin
  if to_regclass('crm.availability_entry') is null then raise exception 'availability_entry missing'; end if;
  if to_regclass('crm.appointment_transition') is null then raise exception 'appointment_transition missing'; end if;
  if to_regprocedure('crm.enforce_appointment_slot_availability()') is null then raise exception 'slot availability guard missing'; end if;
  if to_regprocedure('crm.enforce_appointment_resource_assignment()') is null then raise exception 'resource availability guard missing'; end if;
  if not exists(select 1 from pg_trigger where tgname='appointment_transition_immutable') then raise exception 'transition immutability missing'; end if;
end $$;
rollback;
