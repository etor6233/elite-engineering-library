-- AUTHORED local fixture for CRM customer timeline + append-only history.
-- Exercises customer-scoped projection invariants without a live IdP or provider.
begin;

insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29b88','crm-timeline-fixture','CRM Timeline Fixture','CRM Timeline Fixture');
insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29b88','store','store','Store','store'),
       ('018f4d4a-7b36-7a21-8d10-2f4c54c29b88','other','other','Other','store');
insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29b88','model','model','Model','scooter','active');
insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29b88','variant','model','variant','Variant','{}','active');
insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29b88','customer-a','Customer A','customer-a@example.test'),
       ('018f4d4a-7b36-7a21-8d10-2f4c54c29b88','customer-b','Customer B','customer-b@example.test');
insert into crm.lead(tenant_id,lead_id,organization_id,customer_principal_id,model_id,lifecycle_state,source_code,contact_payload)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29b88','lead-a','store','customer-a','model','contacted','public-web','{"email":"customer-a@example.test"}');
insert into crm.consent_evidence(tenant_id,consent_id,lead_id,purpose_code,policy_version,decision,occurred_at,evidence_sha256_hex)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29b88','consent-a','lead-a','sales-contact','policy-1','granted',clock_timestamp(),repeat('a',64));
insert into crm.appointment(tenant_id,appointment_id,organization_id,lead_id,customer_principal_id,model_id,appointment_kind,starts_at,ends_at,state,version)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29b88','appt-a','store','lead-a','customer-a','model','consultation',clock_timestamp()+interval '2 hours',clock_timestamp()+interval '3 hours','confirmed',3);
insert into crm.appointment_transition(tenant_id,transition_id,appointment_id,from_state,to_state,actor_subject,occurred_at)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29b88','018f4d4a-7b36-7a21-8d10-2f4c54c29b01','appt-a','requested','confirmed','operator',clock_timestamp());
insert into pricing.price_book(tenant_id,price_book_id,market,currency,valid_from,status)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29b88','retail','AR','ARS',clock_timestamp()-interval '1 day','active');
insert into pricing.price_book_entry(tenant_id,price_book_id,variant_id,amount_minor_units,tax_mode)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29b88','retail','variant',250000,'inclusive');
insert into sales.quotation(tenant_id,quotation_id,organization_id,lead_id,customer_principal_id,variant_id,price_book_id,currency,total_minor_units,valid_until,state,version)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29b88','quote-a','store','lead-a','customer-a','variant','retail','ARS',250000,clock_timestamp()+interval '7 days','issued',1);

do $timeline$
declare
  owner_appointments integer;
  owner_quotes integer;
  stranger_appointments integer;
  other_org_appointments integer;
begin
  select count(*) into owner_appointments
  from crm.appointment
  where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c29b88'
    and organization_id='store'
    and customer_principal_id='customer-a';
  select count(*) into owner_quotes
  from sales.quotation
  where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c29b88'
    and organization_id='store'
    and customer_principal_id='customer-a';
  select count(*) into stranger_appointments
  from crm.appointment
  where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c29b88'
    and organization_id='store'
    and customer_principal_id='customer-b';
  select count(*) into other_org_appointments
  from crm.appointment
  where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c29b88'
    and organization_id='other'
    and customer_principal_id='customer-a';
  if owner_appointments <> 1 or owner_quotes <> 1 then
    raise exception 'owner timeline counts unexpected: appointments=% quotes=%', owner_appointments, owner_quotes;
  end if;
  if stranger_appointments <> 0 or other_org_appointments <> 0 then
    raise exception 'cross-customer/org leak: stranger=% other_org=%', stranger_appointments, other_org_appointments;
  end if;
end;
$timeline$;

do $history$
begin
  begin
    update crm.appointment_transition
    set actor_subject='tampered'
    where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c29b88'
      and transition_id='018f4d4a-7b36-7a21-8d10-2f4c54c29b01';
    raise exception 'appointment_transition update unexpectedly succeeded';
  exception when sqlstate 'P0001' then
    if sqlerrm not like '%immutable appointment transition%' then
      raise;
    end if;
  end;
  begin
    delete from crm.appointment_transition
    where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c29b88'
      and transition_id='018f4d4a-7b36-7a21-8d10-2f4c54c29b01';
    raise exception 'appointment_transition delete unexpectedly succeeded';
  exception when sqlstate 'P0001' then
    if sqlerrm not like '%immutable appointment transition%' then
      raise;
    end if;
  end;
end;
$history$;

select 1 / case when (
  select count(*)
  from crm.appointment_transition
  where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c29b88'
    and appointment_id='appt-a'
    and from_state='requested'
    and to_state='confirmed'
    and actor_subject='operator'
) = 1 then 1 else 0 end;

rollback;
