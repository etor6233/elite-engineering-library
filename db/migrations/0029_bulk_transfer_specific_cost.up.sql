alter table inventory.bulk_transfer_line
  add column specific_receipt_entry_id text;

alter table inventory.bulk_transfer_line
  add constraint bulk_transfer_line_specific_receipt_fk
  foreign key (tenant_id, specific_receipt_entry_id)
  references inventory.bulk_inventory_entry (tenant_id, entry_id);
