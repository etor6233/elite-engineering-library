begin;
do $$ begin
 if exists(select 1 from platform.portal_session) or exists(select 1 from platform.portal_session_retention_summary)
 then raise exception 'cannot remove portal session or revocation evidence';end if;
end $$;
drop table platform.portal_session;
drop table platform.portal_session_retention_summary;
commit;
