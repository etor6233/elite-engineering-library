begin;
do $$ begin
 if exists(select 1 from service_ops.warranty_offer) or exists(select 1 from service_ops.warranty_activation)
 then raise exception 'populated warranty terms cannot be downgraded without an admitted archival plan'; end if;
end $$;
drop trigger warranty_quote_acceptance_guard on sales.quotation;
drop function service_ops.warranty_quote_acceptance_guard();
drop table service_ops.warranty_activation;
drop function service_ops.warranty_activation_binding();
drop table service_ops.warranty_acknowledgement;
drop function service_ops.warranty_acknowledgement_binding();
drop table service_ops.warranty_offer;
drop function service_ops.warranty_offer_binding();
drop function service_ops.warranty_fact_immutable();
commit;
