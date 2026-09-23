begin;
-- AUTHORED immutable consent and source-date persistence glue. No statutory
-- duration, fault exclusion, reimbursement or coverage percentage is inferred.
create table service_ops.warranty_offer (
 tenant_id uuid not null,
 quotation_id text not null,
 organization_id text not null,
 customer_principal_id text not null,
 quotation_version bigint not null check(quotation_version>1),
 profile_sha256 text not null check(profile_sha256 ~ '^[0-9a-f]{64}$'),
 profile_bytes bytea not null check(octet_length(profile_bytes) between 1 and 32768),
 offered_by text not null check(length(offered_by) between 1 and 128),
 offered_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,quotation_id),
 foreign key(tenant_id,quotation_id) references sales.quotation(tenant_id,quotation_id),
 foreign key(tenant_id,organization_id) references org.organization(tenant_id,organization_id),
 foreign key(tenant_id,customer_principal_id) references crm.customer_profile(tenant_id,customer_principal_id)
);
create table service_ops.warranty_acknowledgement (
 tenant_id uuid not null,
 quotation_id text not null,
 quotation_version bigint not null,
 customer_principal_id text not null,
 profile_sha256 text not null check(profile_sha256 ~ '^[0-9a-f]{64}$'),
 evidence_sha256 text not null check(evidence_sha256 ~ '^[0-9a-f]{64}$'),
 acknowledged_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,quotation_id),
 foreign key(tenant_id,quotation_id) references service_ops.warranty_offer(tenant_id,quotation_id)
);
create table service_ops.warranty_activation (
 tenant_id uuid not null,
 warranty_id text not null,
 handover_id text not null,
 quotation_id text not null,
 order_id text not null,
 organization_id text not null,
 customer_principal_id text not null,
 stock_unit_id text not null,
 profile_sha256 text not null check(profile_sha256 ~ '^[0-9a-f]{64}$'),
 terms_version text not null,
 parts_start date not null, parts_end date not null,
 labor_start date not null, labor_end date not null,
 handover_accepted_at timestamptz not null,
 activated_by text not null check(length(activated_by) between 1 and 128),
 activated_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,warranty_id),
 unique(tenant_id,handover_id),
 unique(tenant_id,stock_unit_id),
 foreign key(tenant_id,warranty_id) references service_ops.warranty(tenant_id,warranty_id),
 foreign key(tenant_id,quotation_id) references service_ops.warranty_acknowledgement(tenant_id,quotation_id),
 foreign key(tenant_id,handover_id) references sales.delivery_handover(tenant_id,handover_id),
 foreign key(tenant_id,order_id) references sales.customer_order(tenant_id,order_id),
 check(parts_start<=parts_end and labor_start<=labor_end)
);
create function service_ops.warranty_fact_immutable() returns trigger language plpgsql as $$
begin raise exception 'warranty fact is immutable'; end;
$$;
create trigger warranty_offer_immutable before update or delete on service_ops.warranty_offer
 for each row execute function service_ops.warranty_fact_immutable();
create trigger warranty_acknowledgement_immutable before update or delete on service_ops.warranty_acknowledgement
 for each row execute function service_ops.warranty_fact_immutable();
create trigger warranty_activation_immutable before update or delete on service_ops.warranty_activation
 for each row execute function service_ops.warranty_fact_immutable();

create function service_ops.warranty_offer_binding() returns trigger language plpgsql as $$
begin
 perform 1 from sales.quotation q where q.tenant_id=new.tenant_id and q.quotation_id=new.quotation_id
 and q.organization_id=new.organization_id and q.customer_principal_id=new.customer_principal_id
 and q.version=new.quotation_version and q.state='issued' and q.valid_until>clock_timestamp() for update;
 if not found then raise exception 'warranty offer requires current issued quotation'; end if;
 return new;
end;
$$;
create trigger warranty_offer_binding before insert on service_ops.warranty_offer
 for each row execute function service_ops.warranty_offer_binding();
create function service_ops.warranty_acknowledgement_binding() returns trigger language plpgsql as $$
begin
 perform 1 from sales.quotation q join service_ops.warranty_offer o using(tenant_id,quotation_id)
 where q.tenant_id=new.tenant_id and q.quotation_id=new.quotation_id and q.state='issued'
 and q.valid_until>clock_timestamp() and q.version=new.quotation_version and o.quotation_version=q.version
 and q.customer_principal_id=new.customer_principal_id and o.customer_principal_id=new.customer_principal_id
 and o.profile_sha256=new.profile_sha256 for update of q;
 if not found then raise exception 'warranty acknowledgement does not bind offered quotation'; end if;
 return new;
end;
$$;
create trigger warranty_acknowledgement_binding before insert on service_ops.warranty_acknowledgement
 for each row execute function service_ops.warranty_acknowledgement_binding();
create function service_ops.warranty_quote_acceptance_guard() returns trigger language plpgsql as $$
declare o service_ops.warranty_offer%rowtype;
begin
 select * into o from service_ops.warranty_offer where tenant_id=old.tenant_id and quotation_id=old.quotation_id;
 if not found then return new; end if;
 if row(new.organization_id,new.customer_principal_id,new.variant_id,new.price_book_id,new.currency,new.total_minor_units,new.valid_until)
 is distinct from row(old.organization_id,old.customer_principal_id,old.variant_id,old.price_book_id,old.currency,old.total_minor_units,old.valid_until)
 then raise exception 'offered warranty quote terms are immutable; issue a new quotation'; end if;
 if old.state='issued' and new.state='accepted' then
  perform 1 from service_ops.warranty_acknowledgement a where a.tenant_id=old.tenant_id and a.quotation_id=old.quotation_id
  and a.quotation_version=old.version and o.quotation_version=old.version and a.profile_sha256=o.profile_sha256
  and a.customer_principal_id=old.customer_principal_id;
  if not found then raise exception 'warranty terms acknowledgement required before quote acceptance'; end if;
 end if;
 return new;
end;
$$;
create trigger warranty_quote_acceptance_guard before update on sales.quotation
 for each row execute function service_ops.warranty_quote_acceptance_guard();

create function service_ops.warranty_activation_binding() returns trigger language plpgsql as $$
begin
 perform 1 from sales.delivery_handover h
 join sales.quotation_acceptance qa on qa.tenant_id=h.tenant_id and qa.order_id=h.order_id
 join service_ops.warranty_offer o on o.tenant_id=qa.tenant_id and o.quotation_id=qa.quotation_id
 join service_ops.warranty_acknowledgement a on a.tenant_id=o.tenant_id and a.quotation_id=o.quotation_id
 join service_ops.warranty w on w.tenant_id=h.tenant_id and w.warranty_id=new.warranty_id
 where h.tenant_id=new.tenant_id and h.handover_id=new.handover_id and h.state='accepted'
 and h.customer_accepted_at=new.handover_accepted_at and h.organization_id=new.organization_id
 and h.order_id=new.order_id and h.stock_unit_id=new.stock_unit_id and h.customer_principal_id=new.customer_principal_id
 and qa.quotation_id=new.quotation_id and qa.customer_principal_id=h.customer_principal_id
 and o.organization_id=h.organization_id and o.customer_principal_id=h.customer_principal_id
 and o.profile_sha256=new.profile_sha256 and a.profile_sha256=o.profile_sha256
 and a.quotation_version=o.quotation_version and a.customer_principal_id=h.customer_principal_id
 and w.stock_unit_id=h.stock_unit_id and w.customer_principal_id=h.customer_principal_id and w.terms_version=new.terms_version;
 if not found then raise exception 'warranty activation requires sold terms and accepted handover'; end if;
 return new;
end;
$$;
create trigger warranty_activation_binding before insert on service_ops.warranty_activation
 for each row execute function service_ops.warranty_activation_binding();
commit;
