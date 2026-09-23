begin;
drop table if exists fiscal.return_credit_note_link;
drop function if exists fiscal.prevent_return_credit_link_mutation();
commit;
