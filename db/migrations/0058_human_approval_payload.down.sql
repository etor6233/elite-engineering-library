begin;
do $$begin
 if exists(select 1 from approval.request where kind in ('whatsapp_reply','social_publish','social_revoke')) then
  raise exception 'rollback requires preserving bound approval evidence; no destructive downgrade';
 end if;
end $$;
drop trigger bound_decision_immutable on approval.decision;
drop function approval.guard_bound_decision();
drop trigger bound_request_immutable on approval.request;
drop function approval.guard_bound_request();
alter table approval.request drop constraint request_payload_binding_check;
alter table approval.request drop constraint request_kind_check;
alter table approval.request add constraint request_kind_check check(kind in ('reservation','sale','refund','payment'));
alter table approval.request drop column organization_id;
alter table approval.request drop column payload;
commit;
