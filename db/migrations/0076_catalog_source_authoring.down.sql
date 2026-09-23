begin;
do $$begin
 if exists(select 1 from catalog.release_command where kind='source') then
  raise exception 'catalog source authoring history exists; downgrade refused';
 end if;
end$$;
alter table catalog.release_command drop constraint release_command_kind_check;
alter table catalog.release_command add constraint release_command_kind_check
 check(kind in('media','draft','review','publish'));
commit;
