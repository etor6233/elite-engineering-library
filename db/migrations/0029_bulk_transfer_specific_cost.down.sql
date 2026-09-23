alter table inventory.bulk_transfer_line
  drop constraint if exists bulk_transfer_line_specific_receipt_fk;

alter table inventory.bulk_transfer_line
  drop column if exists specific_receipt_entry_id;
