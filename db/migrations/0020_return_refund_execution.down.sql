begin;

drop index if exists payment.return_refund_observation_request_idx;
drop index if exists payment.return_refund_payment_idx;
drop trigger if exists return_refund_observation_immutable on payment.return_refund_observation;
drop function if exists payment.prevent_return_refund_observation_mutation();
drop trigger if exists return_refund_lifecycle on payment.return_refund;
drop function if exists payment.enforce_return_refund_lifecycle();
drop table if exists payment.return_refund_observation;
drop index if exists payment.return_refund_provider_reference_uq;
drop table if exists payment.return_refund;

commit;
