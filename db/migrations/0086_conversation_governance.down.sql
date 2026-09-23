begin;
do $f$ begin
 if exists(select 1 from communication.conversation_reservation) or exists(select 1 from communication.conversation_intent) or exists(select 1 from communication.conversation_budget) then
  raise exception using errcode='55000',message='populated conversation governance rollback requires an explicit data migration';
 end if;
end; $f$;
drop table communication.conversation_intent;
drop table communication.conversation_reservation;
drop table communication.conversation_budget;
drop function communication.reject_conversation_governance_update();
commit;
