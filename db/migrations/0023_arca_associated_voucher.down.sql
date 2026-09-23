drop trigger if exists fiscal_association_immutable on fiscal.invoice_associated_voucher;
drop trigger if exists fiscal_association_valid on fiscal.invoice_associated_voucher;
drop trigger if exists fiscal_invoice_association_valid on fiscal.invoice;
drop function if exists fiscal.validate_associated_voucher();
drop table if exists fiscal.invoice_associated_voucher;
