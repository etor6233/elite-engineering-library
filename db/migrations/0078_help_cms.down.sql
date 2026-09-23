begin;
do $$begin if exists(select 1 from help.revision)then raise exception 'cannot discard help content history';end if;end$$;
alter table help.article drop constraint help_current_revision_fk;
drop table help.revision;
drop table help.article;
drop function help.guard_revision();
drop function help.guard_head();
drop schema help;
commit;
