begin;

drop index if exists sales.return_effect_attempt_request_idx;
drop index if exists sales.return_effect_execution_claim_idx;
drop trigger if exists return_effect_resume_immutable on sales.return_effect_resume;
drop trigger if exists return_effect_attempt_immutable on sales.return_effect_attempt;
drop function if exists sales.prevent_return_effect_audit_mutation();
drop trigger if exists return_effect_execution_enqueue on sales.return_effect_request;
drop function if exists sales.enqueue_return_effect_execution();
drop table if exists sales.return_effect_resume;
drop table if exists sales.return_effect_attempt;
drop table if exists sales.return_effect_execution;

commit;
