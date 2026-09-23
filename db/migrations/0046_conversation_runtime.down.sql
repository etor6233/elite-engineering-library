begin;

drop trigger if exists conversation_turn_identity_guard on communication.conversation_turn;
drop function if exists communication.guard_conversation_turn_identity();
drop table if exists communication.conversation_turn;

commit;
