begin;
do $$begin
 if exists(select 1 from sales.delivery_handover_preparation where funding_receipt_id is not null) or exists(select 1 from sales.commercial_release_receipt where funding_receipt_id is not null) then raise exception 'preserve local-funded delivery history';end if;
end $$;
alter table sales.commercial_release_receipt drop constraint commercial_release_funding_kind;
alter table sales.commercial_release_receipt drop constraint commercial_local_funding_generation;
alter table sales.commercial_release_receipt drop column funding_receipt_id;
alter table sales.commercial_release_receipt alter column payment_attempt_id set not null;
alter table sales.delivery_handover_preparation drop constraint initial_handover_funding_kind;
alter table sales.delivery_handover_preparation drop column funding_receipt_id;
alter table sales.delivery_handover_preparation alter column payment_attempt_id set not null;
commit;
