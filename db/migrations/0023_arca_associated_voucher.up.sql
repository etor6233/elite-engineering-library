create table fiscal.invoice_associated_voucher (
  tenant_id uuid not null,
  invoice_id text not null,
  associated_invoice_id text not null,
  taxpayer_cuit text not null check (taxpayer_cuit ~ '^[0-9]{11}$'),
  voucher_type integer not null check (voucher_type in (1,6,11)),
  point_of_sale_number integer not null check (point_of_sale_number between 1 and 99999),
  voucher_number bigint not null check (voucher_number > 0),
  issued_on date not null,
  primary key (tenant_id,invoice_id),
  foreign key (tenant_id,invoice_id) references fiscal.invoice(tenant_id,invoice_id),
  foreign key (tenant_id,associated_invoice_id) references fiscal.invoice(tenant_id,invoice_id),
  unique (tenant_id,associated_invoice_id),
  check (invoice_id<>associated_invoice_id)
);

create or replace function fiscal.validate_associated_voucher() returns trigger language plpgsql as $$
declare
  target_tenant uuid := case when tg_table_name='invoice' then new.tenant_id else case when tg_op='DELETE' then old.tenant_id else new.tenant_id end end;
  target_invoice text := case when tg_table_name='invoice' then new.invoice_id else case when tg_op='DELETE' then old.invoice_id else new.invoice_id end end;
  child_type integer;
  association_count bigint;
  valid_count bigint;
begin
  select voucher_type into child_type from fiscal.invoice where tenant_id=target_tenant and invoice_id=target_invoice;
  if not found then return null; end if;
  select count(*) into association_count from fiscal.invoice_associated_voucher where tenant_id=target_tenant and invoice_id=target_invoice;
  if child_type in (3,8,13) then
    if association_count<>1 then raise exception 'credit note requires exactly one associated voucher'; end if;
    select count(*) into valid_count
    from fiscal.invoice_associated_voucher a
    join fiscal.invoice child on child.tenant_id=a.tenant_id and child.invoice_id=a.invoice_id
    join fiscal.invoice original on original.tenant_id=a.tenant_id and original.invoice_id=a.associated_invoice_id
    join fiscal.point_of_sale original_pos on original_pos.tenant_id=original.tenant_id and original_pos.point_of_sale_id=original.point_of_sale_id
    join fiscal.point_of_sale child_pos on child_pos.tenant_id=child.tenant_id and child_pos.point_of_sale_id=child.point_of_sale_id
    where a.tenant_id=target_tenant and a.invoice_id=target_invoice
      and original.status='authorized' and original.organization_id=child.organization_id
      and original.voucher_type=case child.voucher_type when 3 then 1 when 8 then 6 when 13 then 11 end
      and original.voucher_number=a.voucher_number and original.issued_on=a.issued_on
      and original_pos.taxpayer_cuit=a.taxpayer_cuit and original_pos.taxpayer_cuit=child_pos.taxpayer_cuit
      and original_pos.point_of_sale_number=a.point_of_sale_number and a.issued_on<=child.issued_on;
    if valid_count<>1 then raise exception 'associated voucher identity is invalid'; end if;
  elsif child_type in (1,6,11) then
    if association_count<>0 then raise exception 'invoice cannot contain an associated voucher'; end if;
  else
    raise exception 'voucher type outside admitted invoice and credit-note lane';
  end if;
  return null;
end; $$;

create constraint trigger fiscal_invoice_association_valid after insert or update on fiscal.invoice deferrable initially deferred for each row execute function fiscal.validate_associated_voucher();
create constraint trigger fiscal_association_valid after insert or update or delete on fiscal.invoice_associated_voucher deferrable initially deferred for each row execute function fiscal.validate_associated_voucher();
create trigger fiscal_association_immutable before update or delete on fiscal.invoice_associated_voucher for each row execute function fiscal.prevent_fiscal_detail_mutation();
