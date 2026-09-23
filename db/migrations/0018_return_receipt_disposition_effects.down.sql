begin;

drop index if exists sales.return_effect_request_state_idx;
drop index if exists sales.return_receipt_org_time_idx;
drop trigger if exists return_effect_request_immutable on sales.return_effect_request;
drop trigger if exists return_disposition_immutable on sales.return_disposition;
drop trigger if exists return_receipt_immutable on sales.return_receipt;
drop function if exists sales.prevent_return_processing_mutation();
drop table if exists sales.return_effect_request;
drop trigger if exists return_disposition_validate on sales.return_disposition;
drop function if exists sales.validate_return_disposition_insert();
drop table if exists sales.return_disposition;
drop trigger if exists return_receipt_validate on sales.return_receipt;
drop function if exists sales.validate_return_receipt_insert();
drop table if exists sales.return_receipt;

commit;
