begin;
do $$begin
 if exists(select 1 from platform.outbox_event where aggregate_type='social_approval') then
  raise exception 'social provider evidence exists; retain schema or use reviewed forward migration';
 end if;
end $$;
drop trigger social_observation_immutable on platform.outbox_event;
drop function communication.guard_social_observation();
commit;
