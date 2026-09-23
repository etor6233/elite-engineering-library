begin;
do $$ begin
  if to_regclass('communication.whatsapp_appointment_approval') is null then
    raise exception 'approval table is required';
  end if;
  if not exists(select 1 from pg_trigger where tgrelid='communication.whatsapp_appointment_approval'::regclass and tgname='whatsapp_approval_immutable' and tgenabled='O') then
    raise exception 'approval immutability trigger is required';
  end if;
end $$;
rollback;
