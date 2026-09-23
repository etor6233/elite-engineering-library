-- AUTHORED local fixture. Pre-existing captured order / authorized return are
-- synthetic starting conditions, NOT an originating sales/approval workflow.
-- All actual constraints/triggers remain enabled. Fresh database only.
begin;
insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','fixture351','Synthetic','Synthetic');
insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','store','store','Synthetic','store');
insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','model','model','Synthetic','bicycle','active');
insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','variant','model','variant','Synthetic','{}','active');
insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','customer','Synthetic','synthetic@example.test');
insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,state,version,received_at)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','stock','store','variant','SYNTHETIC-SERIAL','sold',1,clock_timestamp());
insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','order','store','customer','delivered','ARS',1000,1);
insert into sales.customer_order_line(tenant_id,order_id,line_id,variant_id,quantity,unit_price_minor_units,allocated_stock_unit_id)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','order','line','variant',1,1000,'stock');
insert into payment.payment_attempt(tenant_id,payment_attempt_id,order_id,provider_code,provider_reference,idempotency_key,state,currency,amount_minor_units,version)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','payment','order','stripe','synthetic-pi','fixture-payment-00000001','captured','ARS',1000,3);
insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','handover','store','order','customer','stock','prepared',1);
insert into sales.delivery_exception(tenant_id,exception_id,organization_id,handover_id,customer_principal_id,reason_code,details,rejection_evidence_sha256_hex,state,version)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','exception','store','handover','customer','synthetic-return','Synthetic',repeat('a',64),'open',1);
insert into sales.return_authorization(tenant_id,authorization_id,organization_id,exception_id,handover_id,order_id,stock_unit_id,customer_principal_id,disposition,exact_cost_source_order_id,state,authorized_by_subject)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','authorization','store','exception','handover','order','stock','customer','return','order','authorized','synthetic-operator');
insert into sales.return_receipt(tenant_id,receipt_id,authorization_id,organization_id,order_id,stock_unit_id,customer_principal_id,received_serial_number,condition_code,notes,evidence_sha256_hex,received_by_subject)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','receipt','authorization','store','order','stock','customer','SYNTHETIC-SERIAL','opened','Synthetic',repeat('a',64),'synthetic-operator');
insert into sales.return_disposition(tenant_id,disposition_id,receipt_id,inventory_action,customer_remedy,notes,decided_by_subject)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','disposition','receipt','restock','refund','Synthetic','synthetic-operator');
insert into sales.return_effect_request(tenant_id,request_id,disposition_id,effect_kind,owner_context,state,idempotency_key)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','refund','disposition','refund','payment','requested','synthetic-refund-0000000001');
create schema qualification;
create table qualification.effect(run_id text primary key,idempotency_key text not null unique);
commit;
