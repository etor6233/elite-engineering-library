begin;
alter table catalog.release_command drop constraint release_command_kind_check;
alter table catalog.release_command add constraint release_command_kind_check
 check(kind in('media','draft','review','publish','source'));
commit;
