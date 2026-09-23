-- CATALOG / STOCK / PI slice: catalog master → stock FK and pricing (PI) FK.
-- Requires migrations through 0003_electromobility_franchise_modules (and platform/org foundation).
-- Rollback-safe; synthetic tenant only.

begin;

insert into platform.tenant (tenant_id, tenant_code, legal_name, display_name)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c28c88', 'catalog-stock-slice', 'Catalog Stock Slice', 'Catalog Stock Slice');

insert into org.organization (tenant_id, organization_id, organization_code, display_name, organization_type)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c28c88', 'store', 'store', 'Store', 'store');

insert into catalog.vehicle_model (
  tenant_id, model_id, model_code, display_name, vehicle_class, lifecycle_state, publicly_visible
) values (
  '018f4d4a-7b36-7a21-8d10-2f4c54c28c88', 'model-slice', 'model-slice', 'Slice Model', 'bicycle', 'active', true
);

insert into catalog.vehicle_variant (
  tenant_id, variant_id, model_id, variant_code, display_name, battery_specification, lifecycle_state
) values (
  '018f4d4a-7b36-7a21-8d10-2f4c54c28c88', 'variant-slice', 'model-slice', 'variant-slice', 'Slice Variant', '{}', 'active'
);

-- Stock must bind to catalog variant (happy path).
insert into inventory.stock_unit (
  tenant_id, stock_unit_id, organization_id, variant_id, serial_number, state, version, received_at
) values (
  '018f4d4a-7b36-7a21-8d10-2f4c54c28c88', 'stock-slice', 'store', 'variant-slice', 'SLICE-SERIAL-1', 'available', 1, clock_timestamp()
);

-- Pricing (PI) entry must bind to the same catalog variant.
insert into pricing.price_book (
  tenant_id, price_book_id, market, currency, valid_from, status
) values (
  '018f4d4a-7b36-7a21-8d10-2f4c54c28c88', 'retail-slice', 'AR', 'ARS', clock_timestamp() - interval '1 day', 'active'
);

insert into pricing.price_book_entry (
  tenant_id, price_book_id, variant_id, amount_minor_units, tax_mode
) values (
  '018f4d4a-7b36-7a21-8d10-2f4c54c28c88', 'retail-slice', 'variant-slice', 250000, 'inclusive'
);

do $test$
begin
  begin
    insert into inventory.stock_unit (
      tenant_id, stock_unit_id, organization_id, variant_id, serial_number, state, version, received_at
    ) values (
      '018f4d4a-7b36-7a21-8d10-2f4c54c28c88', 'stock-bad-variant', 'store', 'missing-variant', 'SLICE-SERIAL-BAD', 'available', 1, clock_timestamp()
    );
    raise exception 'expected stock FK rejection for unknown variant';
  exception when foreign_key_violation then null;
  end;

  begin
    insert into pricing.price_book_entry (
      tenant_id, price_book_id, variant_id, amount_minor_units, tax_mode
    ) values (
      '018f4d4a-7b36-7a21-8d10-2f4c54c28c88', 'retail-slice', 'missing-variant', 100, 'inclusive'
    );
    raise exception 'expected price book entry FK rejection for unknown variant';
  exception when foreign_key_violation then null;
  end;
end;
$test$;

do $test$
declare
  joined_label text;
  priced_amount bigint;
begin
  select v.display_name into joined_label
  from inventory.stock_unit s
  join catalog.vehicle_variant v using (tenant_id, variant_id)
  where s.tenant_id = '018f4d4a-7b36-7a21-8d10-2f4c54c28c88'
    and s.stock_unit_id = 'stock-slice';

  select e.amount_minor_units into priced_amount
  from pricing.price_book_entry e
  where e.tenant_id = '018f4d4a-7b36-7a21-8d10-2f4c54c28c88'
    and e.price_book_id = 'retail-slice'
    and e.variant_id = 'variant-slice';

  if joined_label <> 'Slice Variant' or priced_amount <> 250000 then
    raise exception 'catalog-stock-pricing join mismatch: label=% amount=%', joined_label, priced_amount;
  end if;
end;
$test$;

rollback;
